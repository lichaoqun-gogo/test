package handlers

import (
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"osm-mngt/internal/models"
	"osm-mngt/pkg/response"
	"path/filepath"
	"strconv"
	"time"

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

func (h *FileHandler) getUserInfo(c *gin.Context) (string, string) {
	userID := c.GetHeader("X-User-ID")
	userName := c.GetHeader("X-User-Name")
	return userID, userName
}

// Upload 上传文件（带容错处理）
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

	// 打开文件
	src, err := file.Open()
	if err != nil {
		response.Error(c, 500, "打开文件失败")
		return
	}
	defer src.Close()

	// 生成MinIO存储路径
	ext := filepath.Ext(file.Filename)
	storePath := fmt.Sprintf("attachments/%s/%d%s", source, time.Now().UnixNano(), ext)

	// 尝试上传到MinIO（容错处理）
	minioUploadSuccess := false
	if h.minioClient != nil {
		_, err = h.minioClient.PutObject(c.Request.Context(), h.bucketName, storePath, src, file.Size, minio.PutObjectOptions{
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

	// 获取文件类型
	fileType := detectFileType(file)

	// 保存数据库记录
	userID, _ := h.getUserInfo(c)
	now := time.Now()
	fileSize := file.Size
	attachment := models.Attachment{
		FileName:    file.Filename,
		FileSize:    &fileSize,
		FileType:    fileType,
		MinioPath:   storePath,
		MinioBucket: h.bucketName,
		RelatedType: source,
		Description: description,
		Uploader:    userID,
		CreateTime:  now,
		UpdateTime:  now,
	}

	if err := h.db.Create(&attachment).Error; err != nil {
		response.Error(c, 500, "保存记录失败")
		return
	}

	// 构建响应
	result := gin.H{
		"id":        attachment.ID,
		"filename":  attachment.FileName,
		"minioPath": attachment.MinioPath,
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

	// 从MinIO获取文件
	obj, err := h.minioClient.GetObject(c.Request.Context(), attachment.MinioBucket, attachment.MinioPath, minio.GetObjectOptions{})
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

	// 从MinIO删除
	err = h.minioClient.RemoveObject(c.Request.Context(), attachment.MinioBucket, attachment.MinioPath, minio.RemoveObjectOptions{})
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

// detectFileType 根据文件扩展名检测文件类型
func detectFileType(file *multipart.FileHeader) string {
	ext := filepath.Ext(file.Filename)
	switch ext {
	case ".pdf":
		return "application/pdf"
	case ".doc":
		return "application/msword"
	case ".docx":
		return "application/vnd.openxmlformats-officedocument.wordprocessingml.document"
	case ".xls":
		return "application/vnd.ms-excel"
	case ".xlsx":
		return "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
	case ".ppt":
		return "application/vnd.ms-powerpoint"
	case ".pptx":
		return "application/vnd.openxmlformats-officedocument.presentationml.presentation"
	case ".txt":
		return "text/plain"
	case ".csv":
		return "text/csv"
	case ".png":
		return "image/png"
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".gif":
		return "image/gif"
	case ".zip":
		return "application/zip"
	case ".rar":
		return "application/x-rar-compressed"
	case ".tar":
		return "application/x-tar"
	case ".gz":
		return "application/gzip"
	default:
		ct := file.Header.Get("Content-Type")
		if ct != "" {
			return ct
		}
		return "application/octet-stream"
	}
}
