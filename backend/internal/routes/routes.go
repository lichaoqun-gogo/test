package routes

import (
	"osm-mngt/internal/config"
	"osm-mngt/internal/handlers"
	"osm-mngt/internal/middleware"

	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
	"gorm.io/gorm"
)

func Setup(db *gorm.DB, minioClient *minio.Client, cfg *config.Config) *gin.Engine {
	gin.SetMode(cfg.Server.Mode)

	engine := gin.New()
	engine.Use(middleware.CORS())
	engine.Use(middleware.Logger())
	engine.Use(gin.Recovery())

	// 注册路由
	api := engine.Group("/api/v1")

	// 认证路由（CAS 登录/回调/登出）
	authHandler := handlers.NewAuthHandler(cfg, db)
	auth := api.Group("/auth")
	{
		auth.GET("/mode", authHandler.GetAuthMode)
		auth.GET("/cas/login", authHandler.CasLogin)
		auth.GET("/cas/callback", authHandler.CasCallback)
		auth.POST("/logout", authHandler.Logout)
		auth.GET("/me", authHandler.GetUserInfo)
		auth.GET("/menus", authHandler.GetUserMenus)
	}

	// 业务 API 路由（应用认证中间件）
	// CAS 模式下验证 JWT，local 模式下放行
	protected := api.Group("")
	protected.Use(middleware.AuthMiddleware(cfg))

	// 用户列表（人员选择等场景，当前从 users 表查询，后期对接统一权限系统）
	protected.GET("/users", authHandler.GetUserList)

	// 台账管理
	assetHandler := handlers.NewAssetHandler(db)
	asset := protected.Group("/asset")
	{
		asset.GET("/projects", assetHandler.ListProjects)
		asset.GET("/projects/:id", assetHandler.GetProject)
		asset.POST("/projects", assetHandler.CreateProject)
		asset.PUT("/projects/:id", assetHandler.UpdateProject)
		asset.DELETE("/projects/:id", assetHandler.DeleteProject)
		asset.GET("/projects/:id/components", assetHandler.GetProjectComponents)
		asset.GET("/projects/:id/license-stats", assetHandler.GetProjectLicenseStats)

		asset.GET("/inner-projects", assetHandler.ListInnerProjects)
		asset.GET("/inner-projects/components", assetHandler.GetInnerProjectComponents)
		asset.GET("/inner-projects/license-stats", assetHandler.GetInnerProjectLicenseStats)
		asset.GET("/inner-projects/modules", assetHandler.GetInnerProjectModules)
		asset.GET("/inner-projects/:id", assetHandler.GetInnerProject)
		asset.POST("/inner-projects", assetHandler.CreateInnerProject)
		asset.PUT("/inner-projects/:id", assetHandler.UpdateInnerProject)
		asset.DELETE("/inner-projects/:id", assetHandler.DeleteInnerProject)
		asset.GET("/inner-projects/:id/dependencies", assetHandler.GetProjectDependencies)

		asset.GET("/components", assetHandler.ListComponents)
		asset.GET("/components/:id", assetHandler.GetComponent)
		asset.GET("/components/:id/related-projects", assetHandler.GetRelatedProjects)

		// 核心依赖统计
		asset.GET("/core-dependencies/summary", assetHandler.GetCoreDependencySummary)
		asset.GET("/core-dependencies/assoc", assetHandler.GetComponentProjectAssoc)
		asset.GET("/core-dependencies/ranking", assetHandler.GetComponentRanking)
		asset.GET("/core-dependencies/version-fragmentation", assetHandler.GetVersionFragmentation)
		asset.GET("/core-dependencies/high-risk-projects", assetHandler.GetHighRiskProjects)
		asset.GET("/core-dependencies/blacklist-detail", assetHandler.GetProjectBlacklistDetail)
		asset.GET("/core-dependencies/license-overview", assetHandler.GetLicenseOverview)
		asset.POST("/core-dependencies/seed-license-data", assetHandler.SeedLicenseTestData)
		asset.POST("/core-dependencies/seed-test-data", assetHandler.SeedCoreDependencyTestData)
		asset.POST("/core-dependencies/seed-high-risk-data", assetHandler.SeedHighRiskTestData)

		// 项目-负责人关联（支持开源基础软件和自研项目）
		projectOwnerHandler := handlers.NewProjectOwnerHandler(db)
		asset.GET("/project-owners", projectOwnerHandler.GetOwnersByProject)
		asset.GET("/project-owners/map", projectOwnerHandler.GetOwnersMap)
		asset.POST("/project-owners", projectOwnerHandler.SaveOwners)
	}

	// 风险管理
	riskHandler := handlers.NewRiskHandler(db)
	risk := protected.Group("/risk")
	{
		// 黑名单配置
		risk.GET("/blacklist", riskHandler.ListBlacklist)
		risk.POST("/blacklist", riskHandler.CreateBlacklist)
		risk.PUT("/blacklist/:id", riskHandler.UpdateBlacklist)
		risk.DELETE("/blacklist/:id", riskHandler.DeleteBlacklist)

		// 黑名单阻断记录
		risk.GET("/blacklist/records", riskHandler.ListBlacklistRecords)
		risk.PUT("/blacklist/records/:id/process", riskHandler.ProcessBlacklistRecord)

		// 许可证
		risk.GET("/licenses", riskHandler.ListLicenses)
		risk.GET("/licenses/:id", riskHandler.GetLicenseDetail)
		risk.GET("/license-groups", riskHandler.ListLicenseGroups)

		// 缺陷台账
		risk.GET("/defects", riskHandler.ListDefects)
		risk.GET("/defects/:id", riskHandler.GetDefect)
		risk.POST("/defects", riskHandler.CreateDefect)
		risk.PUT("/defects/:id/status", riskHandler.UpdateDefectStatus)

		// CICD检查
		risk.GET("/cicd-check", riskHandler.ListCicdCheck)
		risk.GET("/cicd-check/detail", riskHandler.GetCicdDetail)
		risk.GET("/cicd-check/block-records", riskHandler.GetProjectBlockRecords)
	}

	// 任务管理
	taskHandler := handlers.NewTaskHandler(db)
	task := protected.Group("/tasks")
	{
		task.GET("", taskHandler.ListTasks)
		task.GET("/:id", taskHandler.GetTask)
		task.GET("/:id/detail", taskHandler.GetTaskDetail)
		task.POST("", taskHandler.CreateTask)
		task.PUT("/:id", taskHandler.UpdateTask)
		task.PUT("/:id/status", taskHandler.UpdateTaskStatus)
		task.PUT("/:id/withdraw", taskHandler.WithdrawTask)
	}

	// 审批管理（旧版，保留兼容）
	approvalHandler := handlers.NewApprovalHandler(db)
	approval := protected.Group("/approvals")
	{
		approval.POST("/:taskId/approve", approvalHandler.Approve)
		approval.POST("/:taskId/reject", approvalHandler.Reject)
		approval.GET("/:taskId/history", approvalHandler.GetApprovalHistory)
	}

	// 审批流管理（新版，支持动态路由和状态机）
	workflowHandler := handlers.NewWorkflowHandler(db)
	workflow := protected.Group("/workflow")
	{
		// 任务提交
		workflow.POST("/tasks/:id/submit", workflowHandler.SubmitTask)
		// 任务撤回
		workflow.POST("/tasks/:id/withdraw", workflowHandler.WithdrawTask)
		// 审批通过
		workflow.POST("/tasks/:id/approve", workflowHandler.ApproveNode)
		// 驳回任务
		workflow.POST("/tasks/:id/reject", workflowHandler.RejectTask)
		// 查询审批历史
		workflow.GET("/tasks/:id/history", workflowHandler.GetApprovalHistory)
		// 查询审批节点列表
		workflow.GET("/tasks/:id/nodes", workflowHandler.GetApprovalNodes)
		// 查询当前待处理节点
		workflow.GET("/tasks/:id/current-node", workflowHandler.GetCurrentNode)
	}

	// 文件管理
	fileHandler := handlers.NewFileHandler(db, minioClient, cfg.MinIO.BucketName)
	file := protected.Group("/files")
	{
		file.POST("/upload", fileHandler.Upload)
		file.GET("/:id/download", fileHandler.Download)
		file.DELETE("/:id", fileHandler.Delete)
	}

	// 源码备份管理
	sourceBackupHandler := handlers.NewSourceBackupHandler(minioClient, cfg.MinIO.BucketName)
	sourceBackup := protected.Group("/source-backups")
	{
		sourceBackup.GET("/objects", sourceBackupHandler.ListObjects)
		sourceBackup.GET("/download", sourceBackupHandler.DownloadObject)
		sourceBackup.GET("/download-folder", sourceBackupHandler.DownloadFolder)
		sourceBackup.POST("/seed", sourceBackupHandler.SeedTestObjects)
	}

	return engine
}
