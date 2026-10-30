package httpclient

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestClient_GetJSON(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/health" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"ok"}`))
	}))
	defer srv.Close()

	cli := New(Config{BaseURL: srv.URL, Timeout: 2 * time.Second})
	var out struct{ Status string `json:"status"` }
	if err := cli.Do(context.Background(), "GET", "/health", nil, nil, &out); err != nil {
		t.Fatalf("Do failed: %v", err)
	}
	if out.Status != "ok" {
		t.Fatalf("unexpected response: %+v", out)
	}
}

func TestClient_QueryEncoding(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"qs":"` + r.URL.RawQuery + `"}`))
	}))
	defer srv.Close()

	cli := New(Config{BaseURL: srv.URL, Timeout: 2 * time.Second})
	var out struct {
		QS string `json:"qs"`
	}
	if err := cli.Do(context.Background(), "GET", "/list", map[string]string{"a": "1", "b": "x y"}, nil, &out); err != nil {
		t.Fatalf("Do failed: %v", err)
	}
	if !strings.Contains(out.QS, "a=1") || !strings.Contains(out.QS, "b=x+y") {
		t.Fatalf("unexpected query: %s", out.QS)
	}
}

func TestClient_PostHeaderAuth(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.Header.Get("X-Token"); got != "secret-token" {
			t.Errorf("missing auth header: %q", got)
		}
		var req struct{ Name string `json:"name"` }
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Errorf("decode body: %v", err)
		}
		if req.Name != "redis" {
			t.Errorf("unexpected body: %+v", req)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"id":1}`))
	}))
	defer srv.Close()

	cli := New(Config{
		BaseURL: srv.URL,
		Timeout: 2 * time.Second,
		Auth:    HeaderAuth("X-Token", "secret-token"),
	})
	var out struct{ ID int `json:"id"` }
	if err := cli.Do(context.Background(), "POST", "/sync", nil, map[string]string{"name": "redis"}, &out); err != nil {
		t.Fatalf("Do failed: %v", err)
	}
	if out.ID != 1 {
		t.Fatalf("unexpected response: %+v", out)
	}
}

func TestClient_HTTPError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, `{"msg":"bad request"}`, http.StatusBadRequest)
	}))
	defer srv.Close()

	cli := New(Config{BaseURL: srv.URL, Timeout: 2 * time.Second})
	err := cli.Do(context.Background(), "GET", "/x", nil, nil, nil)
	if err == nil {
		t.Fatal("expected error")
	}
	var hErr *HTTPError
	if !errors.As(err, &hErr) {
		t.Fatalf("expected HTTPError, got: %v", err)
	}
	if hErr.StatusCode != http.StatusBadRequest {
		t.Fatalf("unexpected status: %d", hErr.StatusCode)
	}
}

func TestClient_Timeout(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(500 * time.Millisecond)
	}))
	defer srv.Close()

	cli := New(Config{BaseURL: srv.URL, Timeout: 50 * time.Millisecond})
	start := time.Now()
	if err := cli.Do(context.Background(), "GET", "/slow", nil, nil, nil); err == nil {
		t.Fatal("expected timeout error")
	}
	if time.Since(start) > 300*time.Millisecond {
		t.Fatal("timeout not enforced promptly")
	}
}

func TestClient_InsecureSkipVerify(t *testing.T) {
	// httptest.NewTLSServer 使用自签名证书，默认校验必失败；开 insecureSkipVerify 应可连通
	srv := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`"ok"`)) // 合法 JSON 字符串
	}))
	defer srv.Close()

	cli := New(Config{BaseURL: srv.URL, Timeout: 2 * time.Second, InsecureSkipVerify: true})
	var out string
	if err := cli.Do(context.Background(), "GET", "/", nil, nil, &out); err != nil {
		t.Fatalf("expected success with insecure skip, got: %v", err)
	}
	if out != "ok" {
		t.Fatalf("unexpected body: %q", out)
	}
}

func TestClient_DefaultBaseURLTrim(t *testing.T) {
	cli := New(Config{BaseURL: "https://example.com/api/", Timeout: 2 * time.Second})
	if cli.baseURL != "https://example.com/api" {
		t.Fatalf("baseURL not trimmed: %q", cli.baseURL)
	}
}
