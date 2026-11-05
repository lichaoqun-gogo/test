package config

import (
	"fmt"
	"os"
	"path/filepath"

	"osm-mngt/internal/httpclient"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Server   ServerConfig    `yaml:"server"`
	Database DatabaseConfig  `yaml:"database"`
	MinIO    MinIOConfig     `yaml:"minio"`
	Log      LogConfig       `yaml:"log"`
	Auth     AuthConfig      `yaml:"auth"`
	OpenAPI  OpenAPIConfig   `yaml:"openapi"`
	Services map[string]httpclient.ServiceConfig `yaml:"services"` // 下游 HTTP 服务注册（对应 Feign 服务名->地址配置）
	Scheduler SchedulerConfig `yaml:"scheduler"`                   // 定时任务配置
}

// SchedulerConfig 定时任务配置（对应 Java 版在配置文件中配置定时表达式）
type SchedulerConfig struct {
	Enabled               bool   `yaml:"enabled"`               // 是否启用定时任务
	ComponentCleanupCron  string `yaml:"componentCleanupCron"`  // 组件快照清理 cron 表达式（6 段，如 "0 0 * * * *"=每小时整点）
	CleanupBatchSize      int    `yaml:"cleanupBatchSize"`      // 组件归档批次大小（默认 300，不写死）
}

type ServerConfig struct {
	Port int    `yaml:"port"`
	Mode string `yaml:"mode"`
}

type DatabaseConfig struct {
	Host         string `yaml:"host"`
	Port         int    `yaml:"port"`
	User         string `yaml:"user"`
	Password     string `yaml:"password"`
	DBName       string `yaml:"dbname"`
	Charset      string `yaml:"charset"`
	ParseTime    bool   `yaml:"parseTime"`
	Loc          string `yaml:"loc"`
	MaxIdleConns int    `yaml:"maxIdleConns"`
	MaxOpenConns int    `yaml:"maxOpenConns"`
}

type MinIOConfig struct {
	Endpoint        string `yaml:"endpoint"`
	AccessKeyID     string `yaml:"accessKeyId"`
	SecretAccessKey string `yaml:"secretAccessKey"`
	UseSSL          bool   `yaml:"useSSL"`
	BucketName      string `yaml:"bucketName"`
}

type LogConfig struct {
	Level    string `yaml:"level"`
	Format   string `yaml:"format"`
	FilePath string `yaml:"filePath"` // 日志文件路径，为空则只输出到控制台
}

// AuthConfig 认证配置
type AuthConfig struct {
	Mode       string          `yaml:"mode"`       // 认证模式：local=本地mock账号（开发），cas=CAS单点登录（生产）
	CAS        CASConfig       `yaml:"cas"`
	JWT        JWTConfig       `yaml:"jwt"`
	Permission PermissionConfig `yaml:"permission"` // 统一权限系统（用户档案/菜单/组织均由其下发，后端只转发）
}

// PermissionConfig 统一权限系统配置。
// 设计约定：平台的权限类信息（用户档案/角色/权限点/菜单/组织）统一由外部「统一权限系统」下发，
// 后端只做转发不落库（登录不再写 users 表、菜单不再存表）。接口未对接（baseURL 为空或请求失败）时，
// 后端打印警告日志并返回兜底数据，保证登录与开发调试不被阻塞；对接后只需对齐本文件各接口路径
// 与 remote.go 中承接结构的字段即可。
type PermissionConfig struct {
	BaseURL         string `yaml:"baseURL"`         // 统一权限系统地址，如 http://iam.example.com；为空=未对接
	AppID           string `yaml:"appId"`           // 本系统在权限系统的应用标识（公司模型：每个系统一个 app，有 app 即有默认角色）
	ProfilePath     string `yaml:"profilePath"`     // 用户档案接口路径（GET {baseURL}{path}?username=xxx）
	MenusPath       string `yaml:"menusPath"`       // 菜单接口路径（GET {baseURL}{path}?username=xxx）
	OrgDeptsPath    string `yaml:"orgDeptsPath"`    // 部门-团队-成员 接口路径
	OrgDeptTeamsPath string `yaml:"orgDeptTeamsPath"` // 部门-所属团队 接口路径
	TeamUsersPath   string `yaml:"teamUsersPath"`   // 团队人员接口路径（GET {baseURL}{path}?team=xxx&dept_name=xxx）
	Timeout         int    `yaml:"timeout"`         // 请求超时（秒），默认 5
}

// CASConfig CAS 单点登录配置
type CASConfig struct {
	ServerURL    string `yaml:"serverURL"`    // CAS 服务器地址，如 http://localhost:8443
	ServiceURL   string `yaml:"serviceURL"`   // 回调地址，如 http://localhost:8080/api/v1/auth/cas/callback
	ValidateURL  string `yaml:"validateURL"`  // 票据验证地址，如 http://localhost:8443/serviceValidate
	FrontendURL  string `yaml:"frontendURL"`  // 前端首页地址，回调成功后重定向，如 http://localhost:3000
}

// JWTConfig JWT 配置
type JWTConfig struct {
	Secret string `yaml:"secret"` // JWT 签名密钥
	Expiry int    `yaml:"expiry"`  // 过期时间（秒），默认 86400（24小时）
}

// OpenAPIConfig 外部开放平台（IAM 权限系统）OpenAPI 调用配置
type OpenAPIConfig struct {
	BaseURL           string `yaml:"baseURL"`           // 网关基址，如 https://iam.example.com/openapi
	AppID             string `yaml:"appId"`             // IAM 分配给应用的 appId（如 osm）
	Secret            string `yaml:"secret"`            // IAM 分配的应用密钥
	Timeout           int    `yaml:"timeout"`           // 请求超时（秒），默认 10
	InsecureSkipVerify bool  `yaml:"insecureSkipVerify"` // 跳过 SSL 证书校验：仅开发环境连自签名网关时开，生产必须 false
}

// Load 加载配置文件
// 规则：先读取 config.yaml 中的 env 字段，再加载 config.{env}.yaml
// env=dev  -> config/config.dev.yaml
// env=prod -> config/config.prod.yaml
// 环境变量 APP_ENV 优先于配置文件中的 env（保留兼容）
func Load() (*Config, error) {
	configDir := "config"

	// 第一步：读取 config.yaml 获取 env 字段
	basePath := filepath.Join(configDir, "config.yaml")
	baseData, err := os.ReadFile(basePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file %s: %w", basePath, err)
	}

	var base struct {
		Env string `yaml:"env"`
	}
	if err := yaml.Unmarshal(baseData, &base); err != nil {
		return nil, fmt.Errorf("failed to parse config file %s: %w", basePath, err)
	}

	// 环境变量优先（保留兼容），否则用配置文件中的 env
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = base.Env
	}

	// 第二步：加载环境配置文件 config.{env}.yaml
	var configPath string
	if env != "" {
		configPath = filepath.Join(configDir, fmt.Sprintf("config.%s.yaml", env))
	} else {
		// env 为空时直接用 config.yaml 本身
		configPath = basePath
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file %s: %w", configPath, err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file %s: %w", configPath, err)
	}

	// 填充默认值（生产环境无需配置这些项）
	cfg.applyDefaults()

	return &cfg, nil
}

// applyDefaults 填充配置默认值
func (c *Config) applyDefaults() {
	// JWT 默认值：不配置时使用内置默认值
	if c.Auth.JWT.Secret == "" {
		c.Auth.JWT.Secret = "osm-mngt-default-jwt-secret"
	}
	if c.Auth.JWT.Expiry <= 0 {
		c.Auth.JWT.Expiry = 86400 // 24 小时
	}
	// 定时任务默认值：清理 cron 默认每小时整点，批次默认 300
	if c.Scheduler.ComponentCleanupCron == "" {
		c.Scheduler.ComponentCleanupCron = "0 0 * * * *"
	}
	if c.Scheduler.CleanupBatchSize <= 0 {
		c.Scheduler.CleanupBatchSize = 300
	}
	// 统一权限系统默认值：超时默认 5 秒，各接口路径给默认值（未对接时 baseURL 为空即跳过远程调用）
	if c.Auth.Permission.Timeout <= 0 {
		c.Auth.Permission.Timeout = 5
	}
	if c.Auth.Permission.ProfilePath == "" {
		c.Auth.Permission.ProfilePath = "/user/profile"
	}
	if c.Auth.Permission.MenusPath == "" {
		c.Auth.Permission.MenusPath = "/user/menus"
	}
	if c.Auth.Permission.OrgDeptsPath == "" {
		c.Auth.Permission.OrgDeptsPath = "/org/departments"
	}
	if c.Auth.Permission.OrgDeptTeamsPath == "" {
		c.Auth.Permission.OrgDeptTeamsPath = "/org/departments/teams"
	}
	if c.Auth.Permission.TeamUsersPath == "" {
		c.Auth.Permission.TeamUsersPath = "/org/teams/users"
	}
}
