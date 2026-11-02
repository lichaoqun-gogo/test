package openapi

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"osm-mngt/internal/httpclient"
)

// Config OpenAPI 调用配置（与 config.OpenAPIConfig 字段对应，调用方负责映射）
type Config struct {
	BaseURL            string        // 网关基址，如 https://iam.example.com/openapi
	AppID              string        // IAM 分配给应用的 appId
	Secret             string        // IAM 分配的应用密钥
	Timeout            time.Duration // 单次请求超时
	InsecureSkipVerify bool          // 跳过 SSL 证书校验：仅开发环境连自签名网关时开，生产必须为 false
}

// Client 外部开放平台 OpenAPI 调用客户端：
// 负责「appId + 签名」认证（请求头 auth-appid/auth-secret/auth-timestamp/auth-sign）与 HTTP 传输，
// 不感知具体业务接口；具体接口方法由上层 adapter 封装。
// 传输复用 httpclient.Client（超时/忽略证书校验/响应处理收敛在通用层），签名作为请求头注入。
type Client struct {
	appID     string
	secret    string
	transport *httpclient.Client
}

// New 构造 OpenAPI 客户端
func New(cfg Config) *Client {
	return &Client{
		appID: cfg.AppID,
		secret: cfg.Secret,
		transport: httpclient.New(httpclient.Config{
			BaseURL:            cfg.BaseURL,
			Timeout:            cfg.Timeout,
			InsecureSkipVerify: cfg.InsecureSkipVerify,
		}),
	}
}

// Do 发起一次带签名认证的 OpenAPI 请求。
//   - method: HTTP 方法（GET/POST/...）
//   - path:   网关相对路径，如 /permission/team/list（不含 query）
//   - params: 请求参数，用于拼接 URL query，并按网关规则参与 auth-sign 签名
//   - body:   JSON 请求体（GET 传 nil）；当前约定不参与签名，若网关要求 body 参与签名，
//     调用方需将 body 中参与签名的字段并入 params
//   - out:    响应体反序列化目标（传 nil 则忽略响应体）；HTTP 非 2xx 返回 *OpenAPIError
func (c *Client) Do(ctx context.Context, method, path string, params map[string]any, body any, out any) error {
	timestamp := strconv.FormatInt(time.Now().UnixMilli(), 10)

	fullURL := c.transport.BaseURL() + path
	if qs := encodeQuery(params); qs != "" {
		fullURL += "?" + qs
	}

	var reqBody io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("openapi: marshal request body: %w", err)
		}
		reqBody = bytes.NewReader(data)
	}

	req, err := http.NewRequestWithContext(ctx, method, fullURL, reqBody)
	if err != nil {
		return fmt.Errorf("openapi: build request: %w", err)
	}
	// appId + 签名认证请求头
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("auth-appid", c.appID)
	req.Header.Set("auth-timestamp", timestamp)
	req.Header.Set("auth-secret", SecretSign(c.secret, timestamp))
	req.Header.Set("auth-sign", MD5Sign(c.appID, c.secret, timestamp, params))

	if err := c.transport.RoundTrip(req, out); err != nil {
		// 通用层错误映射回 OpenAPI 语义（HTTP 状态码 + 原始响应体），保持调用方兼容
		var hErr *httpclient.HTTPError
		if errors.As(err, &hErr) {
			return &OpenAPIError{StatusCode: hErr.StatusCode, Body: hErr.Body}
		}
		return err
	}
	return nil
}

// encodeQuery 将参数编码为 URL query；数组参数展开为多个同名键
func encodeQuery(params map[string]any) string {
	if len(params) == 0 {
		return ""
	}
	values := url.Values{}
	for k, v := range params {
		switch t := v.(type) {
		case []string:
			for _, item := range t {
				values.Add(k, item)
			}
		case []any:
			for _, item := range t {
				values.Add(k, fmt.Sprint(item))
			}
		default:
			values.Add(k, fmt.Sprint(v))
		}
	}
	return values.Encode()
}
