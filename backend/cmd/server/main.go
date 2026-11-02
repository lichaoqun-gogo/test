package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"osm-mngt/internal/config"
	"osm-mngt/internal/database"
	"osm-mngt/internal/minio"
	"osm-mngt/internal/routes"
	"osm-mngt/internal/scheduler"
	"osm-mngt/internal/service/component_cleanup"
	"path/filepath"
)

func main() {
	// 加载配置（根据 APP_ENV 环境变量自动选择 config/config.{env}.yaml）
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 初始化日志文件输出（同时输出到控制台和文件）
	if cfg.Log.FilePath != "" {
		logDir := filepath.Dir(cfg.Log.FilePath)
		if err := os.MkdirAll(logDir, 0755); err != nil {
			log.Printf("Warning: Failed to create log directory %s: %v", logDir, err)
		} else {
			logFile, err := os.OpenFile(cfg.Log.FilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
			if err != nil {
				log.Printf("Warning: Failed to open log file %s: %v (logging to stdout only)", cfg.Log.FilePath, err)
			} else {
				log.SetOutput(io.MultiWriter(os.Stdout, logFile))
				log.Printf("Logging to file: %s", cfg.Log.FilePath)
			}
		}
	}

	// 初始化数据库
	db, err := database.Init(cfg.Database)
	if err != nil {
		log.Fatalf("Failed to init database: %v", err)
	}

	// 初始化MinIO（容错处理：失败时不阻断服务启动）
	minioClient, err := minio.Init(cfg.MinIO)
	if err != nil {
		log.Printf("Warning: Failed to init MinIO: %v (service will continue without file upload)", err)
	} else if minioClient == nil {
		log.Println("Warning: MinIO client is nil, file upload feature will be disabled")
	} else {
		log.Println("Info: MinIO initialized successfully")
	}

	// 创建Gin引擎
	engine := routes.Setup(db, minioClient, cfg)

	// 启动定时任务（组件快照每小时清理，cron 表达式配置在 config）
	cleanupSvc := component_cleanup.New(db, cfg.Scheduler.CleanupBatchSize)
	scheduler.New(cfg.Scheduler, cleanupSvc).Start()

	// 启动服务
	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	log.Printf("Server starting on %s", addr)
	if err := engine.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
