package platform

import (
	"reflect"
	"testing"
)

type mutationRunner struct {
	commands []string
}

func (r *mutationRunner) Run(name string, args ...string) ([]byte, error) {
	r.commands = append(r.commands, name+" "+args[0])
	if name == "helm" && args[0] == "template" {
		return []byte("kind: Deployment\n"), nil
	}
	if name == "helm" && args[0] == "status" {
		return []byte(`{"info":{"status":"deployed"},"chart":{"metadata":{"name":"trussium","version":"1.3.0"}},"config":{"appVersion":"1.22.0"}}`), nil
	}
	return nil, nil
}

func (r *mutationRunner) RunInput(name string, input []byte, args ...string) ([]byte, error) {
	r.commands = append(r.commands, name+" "+args[0])
	return nil, nil
}

func TestApplyInstallValidatesBeforeApplyingAndVerifies(t *testing.T) {
	runner := &mutationRunner{}
	report, err := ApplyInstall(runner, "trussium-system", "trussium", "chart", "", 0)
	if err != nil || !report.Applied || !report.Verified {
		t.Fatalf("report=%+v err=%v", report, err)
	}
	want := []string{"helm template", "kubectl apply", "helm install", "helm status"}
	if !reflect.DeepEqual(runner.commands, want) {
		t.Fatalf("commands=%v", runner.commands)
	}
}
