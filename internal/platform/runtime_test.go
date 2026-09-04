package platform

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHTTPRuntimeClientReady(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/health/ready" {
			t.Fatalf("path = %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"status":"ready"}`))
	}))
	defer server.Close()

	status, err := (HTTPRuntimeClient{BaseURL: server.URL}).Ready()
	if err != nil || status.Status != "ready" {
		t.Fatalf("status=%+v err=%v", status, err)
	}
}

func TestHTTPRuntimeClientBoundsFailures(t *testing.T) {
	server := httptest.NewServer(http.NotFoundHandler())
	defer server.Close()
	if _, err := (HTTPRuntimeClient{BaseURL: server.URL}).Ready(); err == nil {
		t.Fatal("expected non-success response error")
	}
}
