package httpclient

import (
	"fmt"
	"net/http"
)

// Auth 认证策略：在请求发出前向请求注入认证信息。
// 对应 Feign 中按服务配置的拦截器/认证头；以函数形式提供，便于支持
// 无认证、固定头、以及签名等动态计算头的场景。
type Auth func(req *http.Request) error

// NoAuth 无认证（内网直调）
func NoAuth() Auth {
	return func(*http.Request) error { return nil }
}

// HeaderAuth 固定请求头认证（如 Authorization: Bearer xxx / X-Token: xxx）
func HeaderAuth(name, value string) Auth {
	return func(req *http.Request) error {
		if name == "" {
			return fmt.Errorf("httpclient: header auth 缺少认证头名称")
		}
		req.Header.Set(name, value)
		return nil
	}
}
