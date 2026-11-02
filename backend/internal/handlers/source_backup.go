package handlers

import (
	"archive/zip"
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net/url"
	"osm-mngt/pkg/response"
	"sort"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
)

// SourceBackupHandler 源码备份处理器
type SourceBackupHandler struct {
	minioClient *minio.Client
	bucketName  string
}

func NewSourceBackupHandler(minioClient *minio.Client, bucketName string) *SourceBackupHandler {
	return &SourceBackupHandler{
		minioClient: minioClient,
		bucketName:  bucketName,
	}
}

// -----------------------------------------------------------------------
// DTO
// -----------------------------------------------------------------------

// ObjectItem 对象项（文件或目录）
type ObjectItem struct {
	Key          string    `json:"key"`          // 目录前缀或文件完整路径
	Name         string    `json:"name"`         // 显示名称（当前层级名称）
	Type         string    `json:"type"`         // "folder" 或 "file"
	Size         int64     `json:"size"`         // 文件大小；目录为子文件总大小
	LastModified time.Time `json:"lastModified"` // 最后修改时间；目录为子文件最新时间
	FileCount    int       `json:"fileCount"`    // 包含的文件数（目录）
}

// ListObjectsResponse 列出对象响应
type ListObjectsResponse struct {
	List   []ObjectItem `json:"list"`
	Bucket string       `json:"bucket"` // bucket 名称
}

// -----------------------------------------------------------------------
// ListObjects 列出对象（递归列出+按下一层级分组）
// 参数：prefix（当前路径）、query（名称前缀过滤）
// 核心逻辑：递归列出 prefix 下所有文件，按下一级路径分量分组
// 目录项返回所有子文件的合计大小和最新修改时间
// -----------------------------------------------------------------------
func (h *SourceBackupHandler) ListObjects(c *gin.Context) {
	if h.minioClient == nil {
		response.Error(c, 503, "MinIO 服务未配置或连接失败")
		return
	}

	prefix := c.Query("prefix")
	query := c.Query("query")

	if prefix != "" && !strings.HasSuffix(prefix, "/") {
		prefix += "/"
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 30*time.Second)
	defer cancel()

	// 递归列出 prefix 下所有对象
	objectCh := h.minioClient.ListObjects(ctx, h.bucketName, minio.ListObjectsOptions{
		Prefix:    prefix,
		Recursive: true,
		MaxKeys:   10000,
	})

	// 按下一级路径分量分组
	type groupInfo struct {
		key          string // 目录前缀或文件完整路径
		name         string // 显示名称
		isFolder     bool
		totalSize    int64
		lastModified time.Time
		fileCount    int
	}
	groups := make(map[string]*groupInfo)

	for obj := range objectCh {
		if obj.Err != nil {
			response.Error(c, 500, "查询 MinIO 对象失败: "+obj.Err.Error())
			return
		}
		// 跳过目录占位符
		if strings.HasSuffix(obj.Key, "/") {
			continue
		}
		// 跳过与 prefix 完全相同的对象
		if obj.Key == prefix {
			continue
		}

		// 取相对于 prefix 的路径
		relPath := strings.TrimPrefix(obj.Key, prefix)
		if relPath == "" {
			continue
		}

		// 按第一级 "/" 分组
		parts := strings.SplitN(relPath, "/", 2)
		firstComponent := parts[0]

		g, exists := groups[firstComponent]
		if !exists {
			if len(parts) == 1 {
				// 直接文件
				g = &groupInfo{
					key:          obj.Key,
					name:         firstComponent,
					isFolder:     false,
					totalSize:    obj.Size,
					lastModified: obj.LastModified,
					fileCount:    1,
				}
			} else {
				// 子目录
				g = &groupInfo{
					key:          prefix + firstComponent + "/",
					name:         firstComponent,
					isFolder:     true,
					totalSize:    obj.Size,
					lastModified: obj.LastModified,
					fileCount:    1,
				}
			}
			groups[firstComponent] = g
		} else {
			// 累加大小和更新时间
			g.totalSize += obj.Size
			if obj.LastModified.After(g.lastModified) {
				g.lastModified = obj.LastModified
			}
			g.fileCount++
		}
	}

	// 转为列表并过滤
	var items []ObjectItem
	for _, g := range groups {
		// query 过滤
		if query != "" {
			if !strings.HasPrefix(strings.ToLower(g.name), strings.ToLower(query)) {
				continue
			}
		}
		itemType := "file"
		if g.isFolder {
			itemType = "folder"
		}
		items = append(items, ObjectItem{
			Key:          g.key,
			Name:         g.name,
			Type:         itemType,
			Size:         g.totalSize,
			LastModified: g.lastModified,
			FileCount:    g.fileCount,
		})
	}

	// 排序：目录在前，再按名称排序
	sort.Slice(items, func(i, j int) bool {
		if items[i].Type != items[j].Type {
			return items[i].Type == "folder"
		}
		return items[i].Name < items[j].Name
	})

	response.Success(c, ListObjectsResponse{
		List:   items,
		Bucket: h.bucketName,
	})
}

// -----------------------------------------------------------------------
// DownloadObject 下载单个对象（流式传输）
// 参数：key（对象完整路径）
// 使用 MinIO GetObject() 直接流式写入 HTTP response
// -----------------------------------------------------------------------
func (h *SourceBackupHandler) DownloadObject(c *gin.Context) {
	if h.minioClient == nil {
		response.Error(c, 503, "MinIO 服务未配置或连接失败")
		return
	}

	key := c.Query("key")
	if key == "" {
		response.Error(c, 400, "请指定要下载的对象路径")
		return
	}

	decodedKey, err := url.QueryUnescape(key)
	if err != nil {
		decodedKey = key
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Minute)
	defer cancel()

	// 校验对象是否存在并获取元信息
	objInfo, err := h.minioClient.StatObject(ctx, h.bucketName, decodedKey, minio.StatObjectOptions{})
	if err != nil {
		response.Error(c, 404, "对象不存在: "+err.Error())
		return
	}

	// 获取对象流
	object, err := h.minioClient.GetObject(ctx, h.bucketName, decodedKey, minio.GetObjectOptions{})
	if err != nil {
		response.Error(c, 500, "获取对象失败: "+err.Error())
		return
	}
	defer object.Close()

	// 设置响应头
	fileName := decodedKey
	if idx := strings.LastIndex(decodedKey, "/"); idx != -1 {
		fileName = decodedKey[idx+1:]
	}
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", fileName))
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Length", fmt.Sprintf("%d", objInfo.Size))

	// 流式传输到 response
	written, err := io.Copy(c.Writer, object)
	if err != nil {
		log.Printf("Error streaming object %s: %v", decodedKey, err)
		return
	}
	log.Printf("Downloaded %s (%d bytes)", decodedKey, written)
}

// -----------------------------------------------------------------------
// DownloadFolder 下载目录下所有文件（打包为 zip 流式传输）
// 参数：prefix（目录前缀）
// 递归列出 prefix 下所有文件，逐个 GetObject 写入 zip
// -----------------------------------------------------------------------
func (h *SourceBackupHandler) DownloadFolder(c *gin.Context) {
	if h.minioClient == nil {
		response.Error(c, 503, "MinIO 服务未配置或连接失败")
		return
	}

	prefix := c.Query("prefix")
	if prefix == "" {
		response.Error(c, 400, "请指定要下载的目录路径")
		return
	}

	decodedPrefix, err := url.QueryUnescape(prefix)
	if err != nil {
		decodedPrefix = prefix
	}
	if !strings.HasSuffix(decodedPrefix, "/") {
		decodedPrefix += "/"
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Minute)
	defer cancel()

	// 先递归列出所有文件
	objectCh := h.minioClient.ListObjects(ctx, h.bucketName, minio.ListObjectsOptions{
		Prefix:    decodedPrefix,
		Recursive: true,
		MaxKeys:   10000,
	})

	var fileKeys []string
	for obj := range objectCh {
		if obj.Err != nil {
			response.Error(c, 500, "查询对象失败: "+obj.Err.Error())
			return
		}
		if strings.HasSuffix(obj.Key, "/") {
			continue
		}
		fileKeys = append(fileKeys, obj.Key)
	}

	if len(fileKeys) == 0 {
		response.Error(c, 404, "该目录下没有文件")
		return
	}

	// zip 文件名取 prefix 最后一段
	folderName := strings.TrimSuffix(decodedPrefix, "/")
	if idx := strings.LastIndex(folderName, "/"); idx != -1 {
		folderName = folderName[idx+1:]
	}
	if folderName == "" {
		folderName = "download"
	}

	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s.zip\"", folderName))
	c.Header("Content-Type", "application/zip")

	zipWriter := zip.NewWriter(c.Writer)
	defer zipWriter.Close()

	for _, key := range fileKeys {
		object, err := h.minioClient.GetObject(ctx, h.bucketName, key, minio.GetObjectOptions{})
		if err != nil {
			log.Printf("Failed to get object %s: %v", key, err)
			continue
		}

		// zip 内文件名用相对路径
		relPath := strings.TrimPrefix(key, decodedPrefix)
		writer, err := zipWriter.Create(relPath)
		if err != nil {
			object.Close()
			log.Printf("Failed to create zip entry %s: %v", relPath, err)
			continue
		}

		_, err = io.Copy(writer, object)
		object.Close()
		if err != nil {
			log.Printf("Failed to copy object %s to zip: %v", key, err)
		}
	}

	log.Printf("Downloaded folder %s (%d files) as zip", decodedPrefix, len(fileKeys))
}

// -----------------------------------------------------------------------
// SeedTestObjects 上传测试数据到 MinIO
// 路径规则：src/{组织}/{软件}/{版本号}
// 版本号本身就是源码包文件名，前缀已区分项目
// -----------------------------------------------------------------------
func (h *SourceBackupHandler) SeedTestObjects(c *gin.Context) {
	if h.minioClient == nil {
		response.Error(c, 503, "MinIO 服务未配置或连接失败")
		return
	}

	testData := []struct {
		Key     string
		Content string
	}{
		// minio 组织
		{"src/minio/minio/v1.0.0", "minio source v1.0.0"},
		{"src/minio/minio/v1.1.0", "minio source v1.1.0"},
		{"src/minio/console/v0.9.0", "minio console v0.9.0"},
		// apache 组织
		{"src/apache/kafka/v3.7.0", "kafka source v3.7.0"},
		{"src/apache/kafka/v3.6.1", "kafka source v3.6.1"},
		{"src/apache/skywalking/v9.7.0", "skywalking source v9.7.0"},
		// spring 组织
		{"src/spring/boot/v3.2.5", "spring boot source v3.2.5"},
		{"src/spring/boot/v3.3.0", "spring boot source v3.3.0"},
		{"src/spring/cloud/v2023.0.1", "spring cloud source"},
		// google 组织
		{"src/google/gson/v2.10.1", "gson source v2.10.1"},
		{"src/google/guava/v33.0.0", "guava source v33.0.0"},
		// alibaba 组织
		{"src/alibaba/fastjson/v2.0.49", "fastjson source v2.0.49"},
		{"src/alibaba/nacos/v2.3.0", "nacos source v2.3.0"},
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 60*time.Second)
	defer cancel()

	// 先清理 src/ 下所有旧对象
	oldCh := h.minioClient.ListObjects(ctx, h.bucketName, minio.ListObjectsOptions{
		Prefix:    "src/",
		Recursive: true,
		MaxKeys:   10000,
	})
	for obj := range oldCh {
		if obj.Err != nil {
			continue
		}
		h.minioClient.RemoveObject(ctx, h.bucketName, obj.Key, minio.RemoveObjectOptions{})
	}

	uploaded := 0
	for _, item := range testData {
		content := bytes.NewReader([]byte(item.Content))
		_, err := h.minioClient.PutObject(ctx, h.bucketName, item.Key, content, int64(len(item.Content)), minio.PutObjectOptions{
			ContentType: "application/octet-stream",
		})
		if err != nil {
			log.Printf("Failed to upload %s: %v", item.Key, err)
			continue
		}
		uploaded++
	}

	response.Success(c, gin.H{
		"total":    len(testData),
		"uploaded": uploaded,
		"message":  fmt.Sprintf("成功上传 %d/%d 个测试对象", uploaded, len(testData)),
	})
}
