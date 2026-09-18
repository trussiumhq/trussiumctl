package platform

import (
	"reflect"
	"testing"
)

func TestApplyUpgradeValidatesBeforeApplyingAndVerifies(t *testing.T) {
	runner := &mutationRunner{}
	report, err := ApplyUpgrade(runner, "trussium-system", "trussium", "chart", "", [3]string{"1.22.0", "1.3.0", "1.0.2"}, [3]string{"1.23.0", "1.3.0", "1.0.2"}, 0)
	if err != nil || !report.Applied || !report.Verified {
		t.Fatalf("report=%+v err=%v", report, err)
	}
	if len(runner.commands) != 4 || runner.commands[2] != "helm upgrade" {
		t.Fatalf("commands=%v", runner.commands)
	}
}

func TestApplyUpgradeRejectsIncompatibleTargetBeforeRendering(t *testing.T) {
	runner := &mutationRunner{}
	_, err := ApplyUpgrade(runner, "default", "trussium", "chart", "", [3]string{"1.22.0", "1.3.0", "1.0.2"}, [3]string{"2.0.0", "1.3.0", "1.0.2"}, 0)
	if err == nil || len(runner.commands) != 0 {
		t.Fatalf("err=%v commands=%v", err, runner.commands)
	}
}

func TestGuardedUpgradeDiscoversConfirmsAndVerifies(t *testing.T) {
	runner := &guardedUpgradeRunner{}
	report, err := GuardedUpgrade(runner, "trussium-system", "trussium", "trussium-operator", "chart", "", "TRUSSIUM", [3]string{}, [3]string{"1.28.0", "1.3.1", "1.0.3"}, 0)
	if err != nil || !report.Applied || !report.Verified || !report.Validation.Valid {
		t.Fatalf("report=%+v err=%v", report, err)
	}
	want := []string{"helm status", "kubectl get", "helm template", "kubectl apply", "helm upgrade", "helm status"}
	if !reflect.DeepEqual(runner.commands, want) {
		t.Fatalf("commands=%v", runner.commands)
	}
}

func TestGuardedUpgradeRejectsConfirmationBeforeDiscovery(t *testing.T) {
	runner := &guardedUpgradeRunner{}
	_, err := GuardedUpgrade(runner, "trussium-system", "trussium", "trussium-operator", "chart", "", "no", [3]string{}, [3]string{"1.28.0", "1.3.1", "1.0.3"}, 0)
	if err == nil || len(runner.commands) != 0 {
		t.Fatalf("err=%v commands=%v", err, runner.commands)
	}
}

type guardedUpgradeRunner struct {
	commands []string
}

func (r *guardedUpgradeRunner) Run(name string, args ...string) ([]byte, error) {
	r.commands = append(r.commands, name+" "+args[0])
	switch {
	case name == "helm" && args[0] == "status":
		return []byte(`{"info":{"status":"deployed"},"chart":{"metadata":{"name":"trussium","version":"1.3.1"}},"config":{"appVersion":"1.27.0"}}`), nil
	case name == "kubectl" && args[0] == "get":
		return []byte(`{"metadata":{"name":"trussium-operator"},"spec":{"replicas":1,"template":{"spec":{"containers":[{"image":"ghcr.io/trussiumhq/trussium-operator:v1.0.3"}]}}},"status":{"availableReplicas":1}}`), nil
	case name == "helm" && args[0] == "template":
		return []byte("kind: Deployment\n"), nil
	default:
		return nil, nil
	}
}

func (r *guardedUpgradeRunner) RunInput(name string, _ []byte, args ...string) ([]byte, error) {
	r.commands = append(r.commands, name+" "+args[0])
	return nil, nil
}
