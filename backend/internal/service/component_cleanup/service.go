package component_cleanup

import (
	"context"
	"fmt"
	"log"
	"time"

	"osm-mngt/internal/models"

	"gorm.io/gorm"
)

// historyColumns component_snapshot_history 的显式列清单
// = component_snapshot 去掉 is_delete，含 create_time；归档保留原主键 id，故显式列出全部列
var historyColumns = "id, current_id, PARENT_COMPONENT_ID, AUTHOR, PUBLISHER, _group, _name, _version, " +
	"DESCRIPTION, CLASSIFIER, FILENAME, EXTENSION, INTERNAL, MD5, SHA1, SHA_256, SHA_384, SHA_512, " +
	"SHA3_256, SHA3_384, SHA3_512, BLAKE2B_256, BLAKE2B_384, BLAKE2B_512, BLAKE3, NOTES, PURL, " +
	"PURLCOORDINATES, CPE, SWIDTAGID, DIRECT_DEPENDENCIES, EXTERNAL_REFERENCES, COPYRIGHT, LICENSE, " +
	"LICENSE_ID, LAST_RISKSCORE, PROJECT_ID, create_time"

// CleanupService 组件快照每小时清理服务（定时任务逻辑，cron 由 scheduler 层配置触发）
// 流程（对应需求 1.1-1.5）：
//   1.2 查待清理 commit：is_delete=1（未删除）的 current_id 集合
//   1.3 排除不清理：project_inner.current_id（当前最新）+ project_inner_pro_info.is_deleted=0（线上依赖）
//   1.4 待清理 - 不清理 = 要执行的 commitid
//   1.5 逐个 commit：记录 commit_id_history -> 查依赖记录 ids -> 按 batchSize 分批显式归档到 history -> 置逻辑删除
type CleanupService struct {
	db        *gorm.DB
	batchSize int
}

func New(db *gorm.DB, batchSize int) *CleanupService {
	if batchSize <= 0 {
		batchSize = 300
	}
	return &CleanupService{db: db, batchSize: batchSize}
}

// Run 执行一轮组件快照清理
func (s *CleanupService) Run(ctx context.Context) error {
	db := s.db.WithContext(ctx)
	start := time.Now()

	// 1.2 待清理 commit：is_delete=1（未删除）的 current_id
	// 注：MAX(create_time) 在此流程后续仅作归档顺序参考，实际按 current_id 集合处理
	var pending []string
	if err := db.Raw("SELECT current_id FROM component_snapshot WHERE is_delete = 1 GROUP BY current_id").
		Scan(&pending).Error; err != nil {
		return fmt.Errorf("查询待清理 commit 失败: %w", err)
	}
	if len(pending) == 0 {
		return nil
	}

	// 1.3 不清理集合
	keep := map[string]bool{}
	var latest []string // 当前最新（project_inner.current_id）
	if err := db.Table("project_inner").Distinct("current_id").Scan(&latest).Error; err != nil {
		return fmt.Errorf("查询当前最新 commit 失败: %w", err)
	}
	for _, k := range latest {
		keep[k] = true
	}
	var online []string // 线上依赖（project_inner_pro_info.is_deleted=0）
	if err := db.Table("project_inner_pro_info").Where("is_deleted = 0").Distinct("commit_id").Scan(&online).Error; err != nil {
		return fmt.Errorf("查询线上依赖 commit 失败: %w", err)
	}
	for _, k := range online {
		keep[k] = true
	}

	// 1.4 要执行的 commitid = 待清理 - 不清理
	toClean := make([]string, 0, len(pending))
	for _, c := range pending {
		if !keep[c] {
			toClean = append(toClean, c)
		}
	}
	if len(toClean) == 0 {
		return nil
	}

	// 1.5 逐个 commit 归档
	log.Printf("[component-cleanup] 待归档 commit %d 个，批次大小 %d", len(toClean), s.batchSize)
	for _, commitID := range toClean {
		if err := s.cleanupOne(db, commitID); err != nil {
			log.Printf("[component-cleanup] commit %s 归档失败: %v（继续下一个）", commitID, err)
			continue
		}
	}
	log.Printf("[component-cleanup] 本轮清理完成，耗时 %s", time.Since(start))
	return nil
}

// cleanupOne 归档单个 commit：记录 commit_id_history -> 分批归档依赖记录 -> 置逻辑删除
func (s *CleanupService) cleanupOne(db *gorm.DB, commitID string) error {
	// 记录 commit_id_history（唯一键，已存在则跳过，保证不重复归档）
	var exists int64
	if err := db.Model(&models.CommitIDHistory{}).Where("commit_id = ?", commitID).Count(&exists).Error; err != nil {
		return err
	}
	if exists == 0 {
		now := time.Now()
		if err := db.Create(&models.CommitIDHistory{
			CommitID:   commitID,
			CreateBy:   "scheduler",
			CreateTime: now,
			UpdateBy:   "scheduler",
			UpdateTime: now,
		}).Error; err != nil {
			return err
		}
	}

	// 查该 commit 的依赖记录 ids（保留主键归档）
	var ids []int64
	if err := db.Table("component_snapshot").Where("current_id = ? AND is_delete = 1", commitID).
		Pluck("id", &ids).Error; err != nil {
		return err
	}
	if len(ids) == 0 {
		return nil
	}

	// 按 batchSize 分批显式归档（INSERT INTO ... SELECT，保留原主键 id）
	for i := 0; i < len(ids); i += s.batchSize {
		end := i + s.batchSize
		if end > len(ids) {
			end = len(ids)
		}
		batch := ids[i:end]
		if err := db.Exec("INSERT INTO component_snapshot_history ("+historyColumns+") "+
			"SELECT "+historyColumns+" FROM component_snapshot WHERE id IN ?", batch).Error; err != nil {
			return fmt.Errorf("归档批次 %d-%d 失败: %w", i+1, end, err)
		}
	}

	// 归档后更新逻辑删除标识为已删除
	if err := db.Table("component_snapshot").Where("current_id = ? AND is_delete = 1", commitID).
		Update("is_delete", 0).Error; err != nil {
		return fmt.Errorf("更新逻辑删除失败: %w", err)
	}
	log.Printf("[component-cleanup] commit %s 归档完成（%d 条，分 %d 批）", commitID, len(ids), (len(ids)+s.batchSize-1)/s.batchSize)
	return nil
}
