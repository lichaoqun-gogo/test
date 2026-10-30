package httpclient

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRegistry_GetAndAuth(t *testing.T) {
	var gotAuth string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotAuth = r.Header.Get("X-Token")
		w.Write([]byte(`{}`))
	}))
	defer srv.Close()

	reg, err := NewRegistry(map[string]ServiceConfig{
		"collector": {
			URL:                srv.URL,
			Timeout:            2,
			InsecureSkipVerify: false,
			Auth:               AuthConfig{Type: "header", Name: "X-Token", Value: "abc"},
		},
	})
	if err != nil {
		t.Fatalf("NewRegistry failed: %v", err)
	}

	cli, err := reg.Get("collector")
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}
	if err := cli.Do(context.Background(), "GET", "/", nil, nil, nil); err != nil {
		t.Fatalf("Do failed: %v", err)
	}
	if gotAuth != "abc" {
		t.Fatalf("auth header not applied: %q", gotAuth)
	}
}

func TestRegistry_NoAuthDefault(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{}`))
	}))
	defer srv.Close()

	reg, err := NewRegistry(map[string]ServiceConfig{
		"plain": {URL: srv.URL},
	})
	if err != nil {
		t.Fatalf("NewRegistry failed: %v", err)
	}
	cli, _ := reg.Get("plain")
	if err := cli.Do(context.Background(), "GET", "/", nil, nil, nil); err != nil {
		t.Fatalf("Do failed: %v", err)
	}
}

func TestRegistry_GetMissing(t *testing.T) {
	reg, err := NewRegistry(map[string]ServiceConfig{})
	if err != nil {
		t.Fatalf("NewRegistry failed: %v", err)
	}
	if _, err := reg.Get("not-exist"); err == nil {
		t.Fatal("expected error for missing service")
	}
}

func TestRegistry_MissingURL(t *testing.T) {
	if _, err := NewRegistry(map[string]ServiceConfig{"bad": {}}); err == nil {
		t.Fatal("expected error for missing url")
	}
}

func TestRegistry_UnknownAuthType(t *testing.T) {
	if _, err := NewRegistry(map[string]ServiceConfig{
		"bad": {URL: "http://x", Auth: AuthConfig{Type: "sign"}},
	}); err == nil {
		t.Fatal("expected error for unknown auth type")
	}
}
