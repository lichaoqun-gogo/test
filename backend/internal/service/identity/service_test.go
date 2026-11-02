package identity

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"osm-mngt/internal/iam"
	"osm-mngt/internal/openapi"
)

// TestOpenAPIService_GetUserInfo 验证 IAMUser = account + org 接口信息拼接：
// 传入的 account 原样保留，组织信息来自 IAM org/list 接口
func TestOpenAPIService_GetUserInfo(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"code":0,"data":[{"orgCode":"TECH","orgName":"技术部","depName":"技术部","children":[{"orgCode":"TECH-B","orgName":"基础组","parentOrgCode":"TECH"}]}]}`))
	}))
	defer srv.Close()

	cli := iam.NewClient(openapi.New(openapi.Config{BaseURL: srv.URL, AppID: "osm", Secret: "test-secret"}), "osm")
	svc := NewOpenAPIService(cli)

	account := &iam.Account{AccountCode: "zhangsan", AccountName: "张三", DepName: "技术部"}
	user, err := svc.GetUserInfo(context.Background(), account)
	if err != nil {
		t.Fatalf("GetUserInfo err: %v", err)
	}
	if user.Account != account {
		t.Error("account 未被原样拼入 IAMUser")
	}
	if len(user.Orgs) != 1 || user.Orgs[0].OrgCode != "TECH" {
		t.Errorf("orgs 解析异常: %+v", user.Orgs)
	}
	if len(user.Orgs[0].Children) != 1 {
		t.Errorf("org 树形 children 解析异常: %+v", user.Orgs[0].Children)
	}
}

// TestOpenAPIService_GetUserInfo_EmptyAccount account 为空时也应正常返回（orgs 照常查询）
func TestOpenAPIService_GetUserInfo_EmptyAccount(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"code":0,"data":[]}`))
	}))
	defer srv.Close()

	cli := iam.NewClient(openapi.New(openapi.Config{BaseURL: srv.URL, AppID: "osm", Secret: "s"}), "osm")
	svc := NewOpenAPIService(cli)
	if _, err := svc.GetUserInfo(context.Background(), nil); err != nil {
		t.Fatalf("nil account 不应报错: %v", err)
	}
}

// TestOpenAPIService_GetUserInfo_GatewayError 网关非 2xx / 业务码非 0 时返回错误
func TestOpenAPIService_GetUserInfo_GatewayError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"code":500,"message":"internal error"}`))
	}))
	defer srv.Close()

	cli := iam.NewClient(openapi.New(openapi.Config{BaseURL: srv.URL, AppID: "osm", Secret: "s"}), "osm")
	svc := NewOpenAPIService(cli)
	if _, err := svc.GetUserInfo(context.Background(), &iam.Account{AccountCode: "zhangsan"}); err == nil {
		t.Fatal("网关错误应返回 err")
	}
}
