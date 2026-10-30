package openapi

import "fmt"

// OpenAPIError 开放平台非 2xx 响应错误：保留 HTTP 状态码与原始响应体，
// 便于上层 adapter 对齐 IAM 业务错误码时解析。
type OpenAPIError struct {
	StatusCode int    // HTTP 状态码
	Body       string // 原始响应体（JSON 文本）
}

func (e *OpenAPIError) Error() string {
	return fmt.Sprintf("openapi: http %d: %s", e.StatusCode, e.Body)
}
