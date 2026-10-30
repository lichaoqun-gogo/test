package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"osm-mngt/internal/iam"
	"osm-mngt/internal/openapi"
	"osm-mngt/internal/service/identity"

	"github.com/gin-gonic/gin"
)

// newIamHandlerWithMock 构造指向 mock 网关的 IamHandler（OpenAPI 实现）
func newIamHandlerWithMock(t *testing.T, gatewayBody string) *IamHandler {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(gatewayBody))
	}))
	t.Cleanup(srv.Close)
	cli := iam.NewClient(openapi.New(openapi.Config{BaseURL: srv.URL, AppID: "osm", Secret: "test-secret"}), "osm")
	return NewIamHandler(identity.NewOpenAPIService(cli))
}

// TestIamHandler_GetUserInfo 验证对外接口 /user/info 跑通：
// 身份来自请求头（local 模式 / CAS 模式 JWT 注入的 X-User-*），组织来自 IAM org 接口，
// 响应为 {code:0, data:{account, orgs}} 全量权限资源包
func TestIamHandler_GetUserInfo(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := newIamHandlerWithMock(t, `{"code":0,"data":[{"orgCode":"TECH","orgName":"技术部","depName":"技术部"}]}`)

	router := gin.New()
	router.GET("/user/info", handler.GetUserInfo)

	req := httptest.NewRequest("GET", "/user/info", nil)
	req.Header.Set("X-User-ID", "zhangsan")
	req.Header.Set("X-User-Name", "张三")
	req.Header.Set("X-User-Dept", "技术部")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}
	var resp struct {
		Code int `json:"code"`
		Data struct {
			Account *iam.Account `json:"account"`
			Orgs    []*iam.Org   `json:"orgs"`
		} `json:"data"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("解析响应失败: %v, body=%s", err, rec.Body.String())
	}
	if resp.Code != 0 {
		t.Fatalf("code = %d, want 0, body=%s", resp.Code, rec.Body.String())
	}
	if resp.Data.Account == nil || resp.Data.Account.AccountCode != "zhangsan" || resp.Data.Account.AccountName != "张三" {
		t.Errorf("account 未正确解析: %+v", resp.Data.Account)
	}
	if len(resp.Data.Orgs) != 1 || resp.Data.Orgs[0].OrgCode != "TECH" {
		t.Errorf("orgs 未正确拼接: %+v", resp.Data.Orgs)
	}
}

// TestIamHandler_GetUserInfo_NoUser 未识别到用户时返回业务错误
func TestIamHandler_GetUserInfo_NoUser(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := newIamHandlerWithMock(t, `{"code":0,"data":[]}`)

	router := gin.New()
	router.GET("/user/info", handler.GetUserInfo)

	req := httptest.NewRequest("GET", "/user/info", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	var resp struct {
		Code int `json:"code"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("解析响应失败: %v", err)
	}
	if resp.Code == 0 {
		t.Fatal("未识别用户时应返回业务错误 code!=0")
	}
}
