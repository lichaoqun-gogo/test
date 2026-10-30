package httpclient

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Config 客户端配置
type Config struct {
	BaseURL            string        // 服务基址，如 https://collector.example.com/api
	Timeout            time.Duration // 单次请求超时，默认 10s
	InsecureSkipVerify bool          // 跳过 SSL 证书校验：仅自签名/内网环境开，生产必须 false
	Auth               Auth          // 认证策略（nil 视为无认证）
}

// Client 通用 HTTP 调用客户端：负责传输（超时/证书/重试不需要的编解码）+ JSON 编解码 + 错误包装。
// 不感知具体业务接口；具体接口方法由上层服务 Client 封装（对应 Feign 的接口定义层）。
type Client struct {
	baseURL string
	hc      *http.Client
	auth    Auth
}

// New 构造客户端
func New(cfg Config) *Client {
	timeout := cfg.Timeout
	if timeout <= 0 {
		timeout = 10 * time.Second
	}
	transport := &http.Transport{}
	if cfg.InsecureSkipVerify {
		// 信任所有证书：收敛为配置项，避免生产环境误开；默认关闭
		transport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	}
	return &Client{
		baseURL: strings.TrimRight(cfg.BaseURL, "/"),
		hc:      &http.Client{Timeout: timeout, Transport: transport},
		auth:    cfg.Auth,
	}
}

// BaseURL 返回客户端基址（供需要自定义请求构造的上层复用）
func (c *Client) BaseURL() string {
	return c.baseURL
}

// Do 发起一次 JSON HTTP 请求。
//   - method: HTTP 方法（GET/POST/PUT/DELETE...）
//   - path:   服务相对路径，如 /defects/sync（不含 query）
//   - query:  URL query 参数（可为 nil），值按字符串处理
//   - body:   JSON 请求体（GET 传 nil）
//   - out:    响应体反序列化目标（传 nil 则忽略响应体）；HTTP 非 2xx 返回 *HTTPError
func (c *Client) Do(ctx context.Context, method, path string, query map[string]string, body any, out any) error {
	fullURL := c.baseURL + path
	if qs := encodeQuery(query); qs != "" {
		fullURL += "?" + qs
	}

	var reqBody io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("httpclient: marshal request body: %w", err)
		}
		reqBody = bytes.NewReader(data)
	}

	req, err := http.NewRequestWithContext(ctx, method, fullURL, reqBody)
	if err != nil {
		return fmt.Errorf("httpclient: build request: %w", err)
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if c.auth != nil {
		if err := c.auth(req); err != nil {
			return fmt.Errorf("httpclient: apply auth: %w", err)
		}
	}
	return c.RoundTrip(req, out)
}

// RoundTrip 执行已构造的请求并统一处理响应（传输层复用点）：
// 非 2xx 返回 *HTTPError；out 非 nil 时反序列化 JSON 响应体。
// 上层需要自定义请求构造（如签名认证）时，构造完 *http.Request 后调用本方法复用传输与响应处理。
func (c *Client) RoundTrip(req *http.Request, out any) error {
	resp, err := c.hc.Do(req)
	if err != nil {
		return fmt.Errorf("httpclient: request %s %s: %w", req.Method, req.URL.String(), err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("httpclient: read response: %w", err)
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return &HTTPError{StatusCode: resp.StatusCode, Body: string(respBody)}
	}
	if out != nil && len(respBody) > 0 {
		if err := json.Unmarshal(respBody, out); err != nil {
			return fmt.Errorf("httpclient: unmarshal response: %w", err)
		}
	}
	return nil
}

// HTTPError 非 2xx 响应错误
type HTTPError struct {
	StatusCode int
	Body       string
}

func (e *HTTPError) Error() string {
	return fmt.Sprintf("httpclient: unexpected status %d: %s", e.StatusCode, e.Body)
}

// encodeQuery 将参数编码为 URL query
func encodeQuery(query map[string]string) string {
	if len(query) == 0 {
		return ""
	}
	values := url.Values{}
	for k, v := range query {
		values.Add(k, v)
	}
	return values.Encode()
}
