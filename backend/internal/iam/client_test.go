package iam

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"osm-mngt/internal/openapi"
)

// newTestClient 构造指向 mock 网关的 iam.Client。
// 认证头由 openapi.Client 自动生成，mock 网关端不校验签名（签名算法已在 openapi 包单测覆盖）。
func newTestClient(t *testing.T, handler http.HandlerFunc) *Client {
	t.Helper()
	srv := httptest.NewServer(handler)
	t.Cleanup(srv.Close)
	oc := openapi.New(openapi.Config{BaseURL: srv.URL, AppID: "osm", Secret: "test-secret"})
	return NewClient(oc, "osm")
}

// mockGateway 按路径返回预设 JSON；同时记录收到的原始 query（空格拼接），供参数传递断言
func mockGateway(t *testing.T, routes map[string]string, gotQueries *string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if gotQueries != nil {
			*gotQueries += r.URL.RawQuery + " "
		}
		if body, ok := routes[r.URL.Path]; ok {
			w.Write([]byte(body))
			return
		}
		t.Errorf("unexpected path: %s", r.URL.Path)
		http.NotFound(w, r)
	}
}

// 各接口的返回 JSON 测试数据（3-5 条，覆盖常规/多值/树形/空/参数传递等形态）
const (
	// version 接口：3 条
	versionJSON1 = `{"code":0,"message":"ok","data":"1.0.0"}`
	versionJSON2 = `{"code":0,"data":"v2.3.1"}`
	versionJSON3 = `{"code":0,"data":"2026.09.07"}`

	// org/map 接口：3 条（单组织/多组织/带下级树）
	orgMapJSON1 = `{"code":0,"data":[{"orgCode":"TECH","orgName":"技术部","orgPrincipal":"zhangsan","orgType":"DEPT","orgLevel":"2","depCode":"D001","depName":"技术部"}]}`
	orgMapJSON2 = `{"code":0,"data":[
		{"orgCode":"TECH","orgName":"技术部","depName":"技术部"},
		{"orgCode":"PLAT","orgName":"平台研发部","depName":"平台研发部"}]}`
	orgMapJSON3 = `{"code":0,"data":[{"orgCode":"TECH","orgName":"技术部","depName":"技术部","children":[
		{"orgCode":"TECH-P","orgName":"平台组","parentOrgCode":"TECH","depName":"平台组"}]}]}`

	// org/list 接口：5 条（空/单组织/多组织/树形/多层级）
	orgListJSON1 = `{"code":0,"data":[]}`
	orgListJSON2 = `{"code":0,"data":[{"orgCode":"PLAT","orgName":"平台研发部","orgPrincipal":"liuyang","depName":"平台研发部"}]}`
	orgListJSON3 = `{"code":0,"data":[
		{"orgCode":"TECH","orgName":"技术部","depName":"技术部"},
		{"orgCode":"PLAT","orgName":"平台研发部","depName":"平台研发部"},
		{"orgCode":"OPS","orgName":"运维部","depName":"运维部"}]}`
	orgListJSON4 = `{"code":0,"data":[{"orgCode":"TECH","orgName":"技术部","depName":"技术部","children":[
		{"orgCode":"TECH-A","orgName":"应用组","parentOrgCode":"TECH","depName":"应用组"},
		{"orgCode":"TECH-B","orgName":"基础组","parentOrgCode":"TECH","depName":"基础组"}]}]}`
	orgListJSON5 = `{"code":0,"data":[{"orgCode":"ROOT","orgName":"集团","orgType":"GROUP","orgLevel":"1","children":[
		{"orgCode":"TECH","orgName":"技术部","orgType":"DEPT","orgLevel":"2","parentOrgCode":"ROOT","children":[
			{"orgCode":"TECH-B","orgName":"基础组","orgType":"TEAM","orgLevel":"3","parentOrgCode":"TECH"}]}]}]}`
)

// TestGetLastestVersion version 接口解析（3 条）
func TestGetLastestVersion(t *testing.T) {
	cases := []struct {
		name, body, want string
	}{
		{"语义化版本", versionJSON1, "1.0.0"},
		{"v 前缀", versionJSON2, "v2.3.1"},
		{"日期版本", versionJSON3, "2026.09.07"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var queries string
			cli := newTestClient(t, mockGateway(t, map[string]string{
				"/openapi/iam/sdk/authz/version": c.body,
			}, &queries))
			got, err := cli.GetLastestVersion(context.Background())
			if err != nil {
				t.Fatalf("GetLastestVersion err: %v", err)
			}
			if got != c.want {
				t.Errorf("version = %q, want %q", got, c.want)
			}
			if !strings.Contains(queries, "app_code=osm") {
				t.Errorf("缺少 app_code 参数，实际 query: %q", queries)
			}
		})
	}
}

// TestGetOrgsByAppCode org/map 接口解析（3 条）
func TestGetOrgsByAppCode(t *testing.T) {
	cases := []struct {
		name, body string
		wantOrgs   int
	}{
		{"单个组织", orgMapJSON1, 1},
		{"多个组织", orgMapJSON2, 2},
		{"带下级树", orgMapJSON3, 1},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var queries string
			cli := newTestClient(t, mockGateway(t, map[string]string{
				"/openapi/iam/api/org/map": c.body,
			}, &queries))
			orgs, err := cli.GetOrgsByAppCode(context.Background())
			if err != nil {
				t.Fatalf("GetOrgsByAppCode err: %v", err)
			}
			if len(orgs) != c.wantOrgs {
				t.Errorf("orgs 数量 = %d, want %d", len(orgs), c.wantOrgs)
			}
			if c.name == "带下级树" && len(orgs[0].Children) != 1 {
				t.Errorf("树形组织 children 数量 = %d, want 1", len(orgs[0].Children))
			}
			if !strings.Contains(queries, "app_code=osm") {
				t.Errorf("缺少 app_code 参数，实际 query: %q", queries)
			}
		})
	}
}

// TestGetOrgs org/list 接口解析（5 条：空/单/多/树形/多层级 + 过滤参数传递）
func TestGetOrgs(t *testing.T) {
	cases := []struct {
		name, body string
		wantOrgs   int
	}{
		{"空列表", orgListJSON1, 0},
		{"单个组织", orgListJSON2, 1},
		{"多个组织", orgListJSON3, 3},
		{"带下级树", orgListJSON4, 1},
		{"多层级树", orgListJSON5, 1},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var queries string
			cli := newTestClient(t, mockGateway(t, map[string]string{
				"/openapi/iam/api/org/list": c.body,
			}, &queries))
			orgs, err := cli.GetOrgs(context.Background(), "", "")
			if err != nil {
				t.Fatalf("GetOrgs err: %v", err)
			}
			if len(orgs) != c.wantOrgs {
				t.Errorf("orgs 数量 = %d, want %d", len(orgs), c.wantOrgs)
			}
		})
	}

	t.Run("过滤参数传递", func(t *testing.T) {
		var queries string
		cli := newTestClient(t, mockGateway(t, map[string]string{
			"/openapi/iam/api/org/list": orgListJSON3,
		}, &queries))
		if _, err := cli.GetOrgs(context.Background(), "1", "IAM"); err != nil {
			t.Fatalf("GetOrgs err: %v", err)
		}
		if !strings.Contains(queries, "org_status=1") || !strings.Contains(queries, "org_source_type=IAM") {
			t.Errorf("过滤参数未传递，实际 query: %q", queries)
		}
	})
}
