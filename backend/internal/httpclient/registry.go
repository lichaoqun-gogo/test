package httpclient

import (
	"fmt"
	"time"
)

// ServiceConfig 下游服务调用配置（对应 Feign 中 @FeignClient 的 url/配置）。
// 在 config.yaml 的 services 段以「服务名 -> ServiceConfig」形式声明，静态配置、重启生效。
type ServiceConfig struct {
	URL                string     `yaml:"url"`                // 服务基址，如 https://collector.example.com/api
	Timeout            int        `yaml:"timeout"`            // 请求超时（秒），0=默认 10
	InsecureSkipVerify bool       `yaml:"insecureSkipVerify"` // 跳过 SSL 证书校验（仅自签名/内网开，生产必须 false）
	Auth               AuthConfig `yaml:"auth"`               // 认证配置
}

// AuthConfig 认证配置
type AuthConfig struct {
	Type  string `yaml:"type"`  // none | header；签名等动态认证由调用方用 Auth 闭包自行构造
	Name  string `yaml:"name"`  // header 认证头名称，如 Authorization
	Value string `yaml:"value"` // header 认证头值，如 Bearer xxx
}

// Registry 下游服务客户端注册表：按服务名持有各自配置的 Client，调用方按需取用。
// 对应 Feign 中「服务名 -> 地址/配置」的注册关系。
type Registry struct {
	clients map[string]*Client
}

// NewRegistry 按配置构建服务客户端注册表；某服务配置非法时返回错误（fail-fast 暴露配置问题）
func NewRegistry(services map[string]ServiceConfig) (*Registry, error) {
	r := &Registry{clients: make(map[string]*Client, len(services))}
	for name, sc := range services {
		if sc.URL == "" {
			return nil, fmt.Errorf("httpclient: 服务 %q 缺少 url 配置", name)
		}
		auth, err := buildAuth(sc.Auth)
		if err != nil {
			return nil, fmt.Errorf("httpclient: 服务 %q 认证配置错误: %w", name, err)
		}
		r.clients[name] = New(Config{
			BaseURL:            sc.URL,
			Timeout:            time.Duration(sc.Timeout) * time.Second,
			InsecureSkipVerify: sc.InsecureSkipVerify,
			Auth:               auth,
		})
	}
	return r, nil
}

// Get 按服务名取客户端；未配置返回错误
func (r *Registry) Get(name string) (*Client, error) {
	c, ok := r.clients[name]
	if !ok {
		return nil, fmt.Errorf("httpclient: 未配置服务 %q（检查 config.yaml services 段）", name)
	}
	return c, nil
}

func buildAuth(a AuthConfig) (Auth, error) {
	switch a.Type {
	case "", "none":
		return NoAuth(), nil
	case "header":
		if a.Name == "" {
			return nil, fmt.Errorf("header 认证缺少 name")
		}
		return HeaderAuth(a.Name, a.Value), nil
	default:
		return nil, fmt.Errorf("未知认证类型 %q（支持 none/header）", a.Type)
	}
}
