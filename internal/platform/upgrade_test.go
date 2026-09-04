package platform

import "testing"

func TestPlanUpgradeRendersOnlyCompatibleTargets(t *testing.T) {
	runner := &fakeRunner{out: []byte("kind: Deployment\n")}
	report := PlanUpgrade(runner, "trussium-system", "trussium", "trussium/trussium", "", [3]string{"1.22.0", "1.3.0", "1.0.2"}, [3]string{"1.23.0", "1.3.0", "1.0.2"})
	if !report.Safe || report.Render == nil || runner.name != "helm" {
		t.Fatalf("report=%+v command=%s", report, runner.name)
	}
}

func TestPlanUpgradeRejectsIncompatibleTargetWithoutHelm(t *testing.T) {
	runner := &fakeRunner{}
	report := PlanUpgrade(runner, "default", "trussium", "chart", "", [3]string{"1.22.0", "1.3.0", "1.0.2"}, [3]string{"2.0.0", "1.3.0", "1.0.2"})
	if report.Safe || report.Render != nil || runner.name != "" {
		t.Fatalf("report=%+v command=%s", report, runner.name)
	}
}
