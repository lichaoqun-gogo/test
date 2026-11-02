package handlers

import (
	"crypto/md5"
	"encoding/hex"
	"time"

	"osm-mngt/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// currentUser 读取当前登录用户（X-User-ID / X-User-Name 请求头）
func currentUser(c *gin.Context) (string, string) {
	return c.GetHeader("X-User-ID"), c.GetHeader("X-User-Name")
}

// computeMD5 计算文件内容 md5（hex 小写），用于 attachment.md5 与去重
func computeMD5(data []byte) string {
	sum := md5.Sum(data)
	return hex.EncodeToString(sum[:])
}

// saveAttachment 记录附件信息：按唯一键 (source, filename, store_url) 判断新建或更新（幂等防重）。
// 兼容生命周期治理流程：不同流程（BOM/升级/退出等）用不同 source 复用同一张附件表。
// 返回：附件 id、是否新建。
func saveAttachment(db *gorm.DB, source, filename, md5Val, storeURL, description, user string) (int64, bool, error) {
	now := time.Now()
	var att models.Attachment
	err := db.Where("source = ? AND filename = ? AND store_url = ?", source, filename, storeURL).First(&att).Error
	if err == nil {
		// 已存在：更新（md5/描述/更新者），不重复建记录
		updates := map[string]interface{}{
			"md5":         md5Val,
			"description": description,
			"update_by":   user,
			"update_time": now,
		}
		if err := db.Model(&att).Updates(updates).Error; err != nil {
			return 0, false, err
		}
		return att.ID, false, nil
	}
	if err != gorm.ErrRecordNotFound {
		return 0, false, err
	}
	att = models.Attachment{
		Source:      source,
		FileName:    filename,
		MD5:         md5Val,
		Description: description,
		StoreURL:    storeURL,
		CreateBy:    user,
		CreateTime:  now,
		UpdateBy:    user,
		UpdateTime:  now,
	}
	if err := db.Create(&att).Error; err != nil {
		return 0, false, err
	}
	return att.ID, true, nil
}
