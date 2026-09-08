package platform

import "testing"

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
