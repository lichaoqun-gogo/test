package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
	MinIO    MinIOConfig    `yaml:"minio"`
	Log      LogConfig      `yaml:"log"`
	Auth     AuthConfig     `yaml:"auth"`
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
	Mode    string   `yaml:"mode"`    // 认证模式：local=本地mock账号（开发），cas=CAS单点登录（生产）
	CAS     CASConfig `yaml:"cas"`
	JWT     JWTConfig `yaml:"jwt"`
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

	return &cfg, nil
}
