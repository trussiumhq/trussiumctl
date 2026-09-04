package platform

import "fmt"

type UpgradeDryRunReport struct {
	CurrentCompatibility CompatibilityReport  `json:"currentCompatibility"`
	TargetCompatibility  CompatibilityReport  `json:"targetCompatibility"`
	Render               *InstallDryRunReport `json:"render,omitempty"`
	Safe                 bool                 `json:"safe"`
	Reasons              []string             `json:"reasons,omitempty"`
}

func PlanUpgrade(runner CommandRunner, namespace, release, chart, values string, current, target [3]string) UpgradeDryRunReport {
	report := UpgradeDryRunReport{
		CurrentCompatibility: CheckCompatibility(current[0], current[1], current[2]),
		TargetCompatibility:  CheckCompatibility(target[0], target[1], target[2]),
	}
	if !report.CurrentCompatibility.Compatible {
		report.Reasons = append(report.Reasons, "current component versions are outside the supported baseline")
	}
	if !report.TargetCompatibility.Compatible {
		report.Reasons = append(report.Reasons, report.TargetCompatibility.Reasons...)
		return report
	}
	render, err := RenderInstall(runner, namespace, release, chart, values)
	if err != nil {
		report.Reasons = append(report.Reasons, err.Error())
		return report
	}
	report.Render = &render
	report.Safe = len(report.Reasons) == 0
	return report
}

func ValidateUpgrade(report UpgradeDryRunReport) error {
	if !report.Safe {
		return fmt.Errorf("upgrade preflight is not safe")
	}
	return nil
}
