package platform

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCollectClusterDiagnostics(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"status":"ready"}`))
	}))
	defer server.Close()
	runner := &diagnosticsRunner{}
	report := CollectClusterDiagnostics(HTTPRuntimeClient{BaseURL: server.URL}, runner, "trussium-system", "trussium-operator", "trussium", "1.22.0", "1.3.0", "1.0.2")
	if len(report.Errors) != 0 || report.Runtime.Status != "ready" || report.Helm.Status != "deployed" {
		t.Fatalf("report=%+v", report)
	}
}

type diagnosticsRunner struct{}

func (*diagnosticsRunner) Run(name string, args ...string) ([]byte, error) {
	if name == "kubectl" {
		return []byte(`{"metadata":{"name":"trussium-operator"},"spec":{"replicas":1},"status":{"availableReplicas":1}}`), nil
	}
	return []byte(`{"info":{"status":"deployed"},"chart":{"metadata":{"name":"trussium","version":"1.3.0"}},"config":{"appVersion":"1.22.0"}}`), nil
}
