package iam

import (
	"context"
	"encoding/json"
	"fmt"

	"osm-mngt/internal/openapi"
)

// Client IAM 开放接口客户端：封装 IAM 各业务接口。
// 对应参考项目「OpenAPISDK 签名拦截器 + Retrofit.create(接口类)」的组合：
//   - openapi.Client 承担签名认证与 HTTP 传输（等价 OpenAPISDK 在 OkHttpClient 上的签名 Interceptor）
//   - 本 Client 每个方法对应一个 Retrofit 接口定义（@GET + 路径 + @Query 参数），直接按方法名调用
type Client struct {
	cli     *openapi.Client
	appCode string // 业务应用编号（接口 app_code 查询参数，如 osm）
}

// NewClient 构造 IAM 接口客户端
func NewClient(cli *openapi.Client, appCode string) *Client {
	return &Client{cli: cli, appCode: appCode}
}

// Response IAM 统一响应包装（对应参考项目 JsonReturn，@ResultAdapAnno 表示响应结果适配/解包 data）。
// 结构按 {code, message, data} 假设实现，字段名与成功码约定待与 IAM 实际返回核对后调整。
type Response[T any] struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    T      `json:"data"`
}

// call 执行一次 IAM GET 请求：签名传输 + 统一响应解包 + 业务码校验
func (c *Client) call(ctx context.Context, path string, params map[string]any, data any) error {
	var resp Response[json.RawMessage]
	if err := c.cli.Do(ctx, "GET", path, params, nil, &resp); err != nil {
		return err
	}
	if resp.Code != 0 { // 成功码约定待与 IAM 核对（常见 0 或 200）
		return fmt.Errorf("iam: %s code %d: %s", path, resp.Code, resp.Message)
	}
	if data != nil && len(resp.Data) > 0 {
		return json.Unmarshal(resp.Data, data)
	}
	return nil
}

// GetLastestVersion 查询 IAM SDK 最新版本（用于版本校验/缓存失效判断）
// GET /openapi/iam/sdk/authz/version?app_code={appCode}
// 返回类型暂按版本字符串承接，具体结构待与 IAM 返回核对
func (c *Client) GetLastestVersion(ctx context.Context) (string, error) {
	var version string
	if err := c.call(ctx, "/openapi/iam/sdk/authz/version", map[string]any{"app_code": c.appCode}, &version); err != nil {
		return "", err
	}
	return version, nil
}

// GetOrgsByAppCode 按应用查询组织映射
// GET /openapi/iam/api/org/map?app_code={appCode}
// 返回承接暂按组织列表，data 具体结构（list/map）待与 IAM 返回核对
func (c *Client) GetOrgsByAppCode(ctx context.Context) ([]*Org, error) {
	var orgs []*Org
	if err := c.call(ctx, "/openapi/iam/api/org/map", map[string]any{"app_code": c.appCode}, &orgs); err != nil {
		return nil, err
	}
	return orgs, nil
}

// GetOrgs 查询组织列表（可过滤）
// GET /openapi/iam/api/org/list?org_status={orgStatus}&org_source_type={orgSourceType}
func (c *Client) GetOrgs(ctx context.Context, orgStatus, orgSourceType string) ([]*Org, error) {
	params := map[string]any{}
	if orgStatus != "" {
		params["org_status"] = orgStatus
	}
	if orgSourceType != "" {
		params["org_source_type"] = orgSourceType
	}
	var orgs []*Org
	if err := c.call(ctx, "/openapi/iam/api/org/list", params, &orgs); err != nil {
		return nil, err
	}
	return orgs, nil
}

// Privs 用户权限（对应参考项目 osm-front 文档：/api/auth/privs 返回的权限信息）。
// 字段按文档约定（roles/menus/buttons），具体结构与 IAM 实际返回核对后调整。
type Privs struct {
	Roles       []string `json:"roles"`       // 角色列表
	Menus       []*Menu  `json:"menus"`       // 菜单权限（含路径与名称）
	Buttons     []string `json:"buttons"`     // 按钮/操作权限码
	Permissions []string `json:"permissions"` // 权限码集合
}

// GetPrivs 查询用户在指定应用/组织下的权限（对应参考项目 IAMAuthUtils 的权限查询能力）
// GET /api/auth/privs?userName={userName}&appCode={appCode}&orgCode={orgCode}
func (c *Client) GetPrivs(ctx context.Context, userName, orgCode string) (*Privs, error) {
	params := map[string]any{"userName": userName, "appCode": c.appCode}
	if orgCode != "" {
		params["orgCode"] = orgCode
	}
	var privs Privs
	if err := c.call(ctx, "/api/auth/privs", params, &privs); err != nil {
		return nil, err
	}
	return &privs, nil
}
