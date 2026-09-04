package platform

import (
	"reflect"
	"testing"
)

type fakeRunner struct {
	name string
	args []string
	out  []byte
	err  error
}

func (f *fakeRunner) Run(name string, args ...string) ([]byte, error) {
	f.name, f.args = name, args
	return f.out, f.err
}

func TestOperatorStatusFor(t *testing.T) {
	runner := &fakeRunner{out: []byte(`{"metadata":{"name":"trussium-operator"},"spec":{"replicas":2},"status":{"availableReplicas":1}}`)}
	status, err := OperatorStatusFor(runner, "trussium-system", "trussium-operator")
	if err != nil || status.Ready != 1 || status.Desired != 2 {
		t.Fatalf("status=%+v err=%v", status, err)
	}
	want := []string{"get", "deployment", "trussium-operator", "--namespace", "trussium-system", "--output", "json"}
	if runner.name != "kubectl" || !reflect.DeepEqual(runner.args, want) {
		t.Fatalf("command=%s %v", runner.name, runner.args)
	}
}

func TestHelmStatusFor(t *testing.T) {
	runner := &fakeRunner{out: []byte(`{"info":{"status":"deployed"},"chart":{"metadata":{"name":"trussium","version":"1.3.0"}},"config":{"appVersion":"1.22.0"}}`)}
	status, err := HelmStatusFor(runner, "trussium-system", "trussium")
	if err != nil || status.Status != "deployed" || status.Chart != "trussium-1.3.0" {
		t.Fatalf("status=%+v err=%v", status, err)
	}
}
