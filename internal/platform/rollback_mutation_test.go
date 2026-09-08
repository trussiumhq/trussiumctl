package platform

import "testing"

func TestApplyRollbackValidatesBeforeApplyingAndVerifies(t *testing.T) {
	runner := &mutationRunner{}
	report, err := ApplyRollback(runner, "trussium-system", "trussium", "chart", "", 4, [3]string{"1.22.0", "1.3.0", "1.0.2"}, 0)
	if err != nil || !report.Applied || !report.Verified || report.Revision != 4 {
		t.Fatalf("report=%+v err=%v", report, err)
	}
	want := []string{"helm template", "kubectl apply", "helm rollback", "helm status"}
	if len(runner.commands) != len(want) {
		t.Fatalf("commands=%v", runner.commands)
	}
	for i := range want {
		if runner.commands[i] != want[i] {
			t.Fatalf("commands=%v", runner.commands)
		}
	}
}

func TestApplyRollbackRejectsIncompatibleTargetBeforeRendering(t *testing.T) {
	runner := &mutationRunner{}
	_, err := ApplyRollback(runner, "default", "trussium", "chart", "", 0, [3]string{"2.0.0", "1.3.0", "1.0.2"}, 0)
	if err == nil || len(runner.commands) != 0 {
		t.Fatalf("err=%v commands=%v", err, runner.commands)
	}
}
