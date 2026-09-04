package platform

import "testing"

func TestCheckCompatibilitySupported(t *testing.T) {
	report := CheckCompatibility("1.22.0", "1.3.0", "1.0.2")
	if !report.Compatible || len(report.Reasons) != 0 {
		t.Fatalf("report=%+v", report)
	}
}

func TestCheckCompatibilityRejectsUnsupportedVersions(t *testing.T) {
	report := CheckCompatibility("2.0.0", "bad", "1.0.1")
	if report.Compatible || len(report.Reasons) != 3 {
		t.Fatalf("report=%+v", report)
	}
}
