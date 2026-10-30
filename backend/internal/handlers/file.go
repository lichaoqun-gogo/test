package handlers

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"path/filepath"
	"strconv"
	"time"

	"osm-mngt/internal/models"
	"osm-mngt/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
	"gorm.io/gorm"
)

type FileHandler struct {
	db          *gorm.DB
	minioClient *minio.Client
	bucketName  string
}

func NewFileHandler(db *gorm.DB, minioClient *minio.Client, bucketName string) *FileHandler {
	return &FileHandler{
		db:          db,
		minioClient: minioClient,
		bucketName:  bucketName,
	}
}

// Upload 通用文件上传：记录附件信息（新建/更新，复用统一逻辑）+ 上传 MinIO（带容错处理）
// 兼容生命周期治理流程：各流程通过 source（AttachmentSource）复用同一张附件表
func (h *FileHandler) Upload(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		response.Error(c, 400, "请选择文件")
		return
	}

	source := c.PostForm("source")
	description := c.PostForm("description")
	if source == "" {
		response.Error(c, 400, "请指定附件来源")
		return
	}

	src, err := file.Open()
	if err != nil {
		response.Error(c, 500, "打开文件失败")
		return
	}
	defer src.Close()
	data, err := io.ReadAll(src)
	if err != nil {
		response.Error(c, 500, "读取文件失败")
		return
	}

	// 生成 MinIO 存储路径（通用附件目录，按 source 隔离）
	ext := filepath.Ext(file.Filename)
	storePath := fmt.Sprintf("attachments/%s/%d%s", source, time.Now().UnixNano(), ext)

	userID, _ := currentUser(c)

	// 记录附件信息（新建/更新，幂等）
	attID, _, err := saveAttachment(h.db, source, file.Filename, computeMD5(data), storePath, description, userID)
	if err != nil {
		response.Error(c, 500, "保存附件记录失败")
		return
	}

	// 上传到 MinIO（容错处理）
	minioUploadSuccess := false
	if h.minioClient != nil {
		_, err = h.minioClient.PutObject(c.Request.Context(), h.bucketName, storePath, bytes.NewReader(data), int64(len(data)), minio.PutObjectOptions{
			ContentType: file.Header.Get("Content-Type"),
		})
		if err != nil {
			// MinIO上传失败，记录日志但不阻断
			log.Printf("Warning: Failed to upload to MinIO: %v (file will be stored locally)", err)
		} else {
			minioUploadSuccess = true
		}
	} else {
		log.Println("Warning: MinIO client not available, file will be stored locally")
	}

	result := gin.H{
		"id":        attID,
		"filename":  file.Filename,
		"minioPath": storePath,
	}
	if !minioUploadSuccess {
		result["warning"] = "MinIO上传失败，文件已保存到数据库，请检查MinIO配置"
	}
	response.Success(c, result)
}

// Download 下载文件
func (h *FileHandler) Download(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, 400, "无效的文件ID")
		return
	}

	var attachment models.Attachment
	if err := h.db.First(&attachment, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			response.Error(c, 404, "文件不存在")
			return
		}
		response.Error(c, 500, "查询失败")
		return
	}

	// 从 MinIO 获取文件
	obj, err := h.minioClient.GetObject(c.Request.Context(), h.bucketName, attachment.StoreURL, minio.GetObjectOptions{})
	if err != nil {
		response.Error(c, 500, "获取文件失败")
		return
	}
	defer obj.Close()

	// 检查对象是否存在
	stat, err := obj.Stat()
	if err != nil {
		response.Error(c, 404, "文件不存在于存储中")
		return
	}

	// 设置响应头
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", attachment.FileName))
	c.Header("Content-Type", stat.ContentType)
	c.Header("Content-Length", fmt.Sprintf("%d", stat.Size))

	// 返回文件流
	if _, err := io.Copy(c.Writer, obj); err != nil {
		// 连接断开等情况下忽略错误
		return
	}
}

// Delete 删除文件
func (h *FileHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, 400, "无效的文件ID")
		return
	}

	var attachment models.Attachment
	if err := h.db.First(&attachment, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			response.Error(c, 404, "文件不存在")
			return
		}
		response.Error(c, 500, "查询失败")
		return
	}

	// 从 MinIO 删除
	err = h.minioClient.RemoveObject(c.Request.Context(), h.bucketName, attachment.StoreURL, minio.RemoveObjectOptions{})
	if err != nil {
		response.Error(c, 500, "删除文件失败")
		return
	}

	// 删除数据库记录
	if err := h.db.Delete(&attachment).Error; err != nil {
		response.Error(c, 500, "删除记录失败")
		return
	}

	response.Success(c, nil)
}
