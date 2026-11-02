package scheduler

import (
	"context"
	"log"

	"osm-mngt/internal/config"
	"osm-mngt/internal/service/component_cleanup"

	"github.com/robfig/cron/v3"
)

// Scheduler 定时任务调度器：cron 表达式与开关从 config 配置读取（对应 Java 版在配置文件中配置定时表达式）。
// 结构：本包只负责「按配置注册并触发任务」，任务的具体逻辑在各 service 中实现，便于后续新增定时任务。
type Scheduler struct {
	cron *cron.Cron
}

// New 构造调度器并注册定时任务（Enabled=false 时不注册）
func New(cfg config.SchedulerConfig, cleanup *component_cleanup.CleanupService) *Scheduler {
	c := cron.New()
	if cfg.Enabled {
		expr := cfg.ComponentCleanupCron
		if _, err := c.AddFunc(expr, func() {
			log.Println("[scheduler] 触发组件快照清理任务")
			if err := cleanup.Run(context.Background()); err != nil {
				log.Printf("[scheduler] 组件快照清理任务失败: %v", err)
			}
		}); err != nil {
			log.Printf("Warning: 注册组件快照清理任务失败（cron=%s）: %v", expr, err)
		}
	}
	return &Scheduler{cron: c}
}

// Start 启动定时任务（非阻塞）
func (s *Scheduler) Start() {
	if s.cron != nil {
		s.cron.Start()
		log.Println("Scheduler started")
	}
}

// Stop 停止定时任务（等待已运行任务结束）
func (s *Scheduler) Stop() {
	if s.cron != nil {
		s.cron.Stop()
	}
}
