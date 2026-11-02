package handlers

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"path/filepath"
	"strings"
	"time"

	"osm-mngt/internal/httpclient"
	"osm-mngt/internal/models"
	"osm-mngt/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
	"gorm.io/gorm"
)

// bomParseService 外部 BOM 解析服务名（config.yaml services 段配置；暂定地址 url1，后期修改）
const bomParseService = "bom-parser"

// BomHandler BOM 清单上传/解析
type BomHandler struct {
	db          *gorm.DB
	minioClient *minio.Client
	bucketName  string
	registry    *httpclient.Registry // 下游服务注册表（可 nil：外部解析服务未配置时跳过调用）
}

func NewBomHandler(db *gorm.DB, minioClient *minio.Client, bucketName string, registry *httpclient.Registry) *BomHandler {
	return &BomHandler{db: db, minioClient: minioClient, bucketName: bucketName, registry: registry}
}

// UploadBom 按 BOM 解析上传（全局治理-开源基础软件-软件名称-按BOM解析）
// 接口参数：uuid（项目 UUID）、file（multipart 文件，非文件名）
// 流程：
//   1. 按 uuid 查项目详情
//   2. 组装存储地址：bom/{项目组织}/{项目名称}/{项目版本}/{文件名}
//   3. 附件信息记录 attachment（新建/更新，含文件名/md5/来源 BOM_UPLOAD/创建者/更新者）
//   4. 记录成功上传 MinIO
//   5. 上传成功后在 bom 表新建或更新记录（按 project_id），store_url 关联 attachment.id
//   6. 调用外部解析服务 /task/bomUpload，方法 upload(uuid, bomId)
func (h *BomHandler) UploadBom(c *gin.Context) {
	uuid := strings.TrimSpace(c.PostForm("uuid"))
	file, err := c.FormFile("file")
	if err != nil {
		response.Error(c, 400, "请选择文件")
		return
	}
	if uuid == "" {
		response.Error(c, 400, "缺少项目 uuid")
		return
	}
	if !strings.EqualFold(filepath.Ext(file.Filename), ".xml") {
		response.Error(c, 400, "仅支持 .xml 文件（CycloneDX BOM）")
		return
	}

	// 1. 按 uuid 查项目详情
	var prj models.Project
	if err := h.db.Where("uuid = ?", uuid).First(&prj).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			response.Error(c, 404, "项目不存在")
			return
		}
		response.Error(c, 500, "查询项目失败")
		return
	}

	// 2. 组装 MinIO 存储地址：bom/{项目组织}/{项目名称}/{项目版本}/{文件名}
	storeDir := fmt.Sprintf("bom/%s/%s/%s", pathSegment(prj.Group), pathSegment(prj.Name), pathSegment(prj.Version))
	storePath := storeDir + "/" + file.Filename

	// 读取文件 + 计算 md5
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
	md5Val := computeMD5(data)

	userID, _ := currentUser(c)

	// 3. 附件信息记录（新建/更新判断，来源 BOM_UPLOAD）
	attID, _, err := saveAttachment(h.db, string(models.AttachmentSourceBOMUpload), file.Filename, md5Val, storePath, "", userID)
	if err != nil {
		response.Error(c, 500, "保存附件记录失败")
		return
	}

	// 4. 上传 MinIO
	if h.minioClient != nil {
		if _, err = h.minioClient.PutObject(c.Request.Context(), h.bucketName, storePath, bytes.NewReader(data), int64(len(data)), minio.PutObjectOptions{
			ContentType: "application/xml",
		}); err != nil {
			response.Error(c, 500, "上传 MinIO 失败: "+err.Error())
			return
		}
	}

	// 5. bom 表新建或更新记录（按 project_id），store_url 关联 attachment.id
	bomID, err := h.upsertBom(prj.ID, attID, userID)
	if err != nil {
		response.Error(c, 500, "保存 BOM 记录失败: "+err.Error())
		return
	}

	// 6. 调用外部解析服务（未配置地址时跳过，不影响主流程）
	externalErr := h.callBomParse(c.Request.Context(), uuid, bomID)
	if externalErr != nil {
		log.Printf("Warning: call bom parse service failed: %v", externalErr)
	}

	result := gin.H{
		"bomId":        bomID,
		"attachmentId": attID,
		"storeUrl":     storePath,
	}
	if externalErr != nil {
		result["externalWarning"] = "BOM 记录已保存，但外部解析服务调用失败: " + externalErr.Error()
	}
	response.Success(c, result)
}

// upsertBom 按 project_id 新建或更新 BOM 记录（一个项目一条），storeURL 指向 attachment.id
func (h *BomHandler) upsertBom(projectID, attID int64, user string) (int64, error) {
	now := time.Now()
	var bom models.Bom
	err := h.db.Where("project_id = ?", projectID).First(&bom).Error
	if err == nil {
		// 更新
		updates := map[string]interface{}{
			"store_url":   attID,
			"update_by":   user,
			"update_time": now,
		}
		if err := h.db.Model(&bom).Updates(updates).Error; err != nil {
			return 0, err
		}
		return bom.ID, nil
	}
	if err != gorm.ErrRecordNotFound {
		return 0, err
	}
	// 新建（元数据由外部解析服务回填，本地先记录基础信息）
	bom = models.Bom{
		BomFormat:    "cyclonedx",
		SerialNumber: fmt.Sprintf("%d", time.Now().UnixNano()),
		StoreURL:     &attID,
		ProjectID:    projectID,
		CreateBy:     user,
		CreateTime:   now,
		UpdateBy:     user,
		UpdateTime:   now,
	}
	if err := h.db.Create(&bom).Error; err != nil {
		return 0, err
	}
	return bom.ID, nil
}

// callBomParse 调用外部 BOM 解析服务：POST /task/bomUpload，body {uuid, bomId}
func (h *BomHandler) callBomParse(ctx context.Context, uuid string, bomID int64) error {
	if h.registry == nil {
		return fmt.Errorf("下游服务注册表未初始化")
	}
	cli, err := h.registry.Get(bomParseService)
	if err != nil {
		return err
	}
	payload := map[string]any{"uuid": uuid, "bomId": bomID}
	return cli.Do(ctx, "POST", "/task/bomUpload", nil, payload, nil)
}

// pathSegment 清理存储路径中的目录段：去分隔符，空值给占位
func pathSegment(s string) string {
	s = strings.TrimSpace(s)
	s = strings.ReplaceAll(s, "/", "_")
	s = strings.ReplaceAll(s, "\\", "_")
	if s == "" {
		return "unknown"
	}
	return s
}
