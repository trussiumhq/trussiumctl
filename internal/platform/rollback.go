package platform

import "fmt"

type RollbackDryRunReport struct {
	TargetCompatibility CompatibilityReport  `json:"targetCompatibility"`
	Render              *InstallDryRunReport `json:"render,omitempty"`
	Safe                bool                 `json:"safe"`
	Reasons             []string             `json:"reasons,omitempty"`
}

func PlanRollback(runner CommandRunner, namespace, release, chart, values string, target [3]string) RollbackDryRunReport {
	report := RollbackDryRunReport{TargetCompatibility: CheckCompatibility(target[0], target[1], target[2])}
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
	report.Safe = true
	return report
}

func ValidateRollback(report RollbackDryRunReport) error {
	if !report.Safe {
		return fmt.Errorf("rollback preflight is not safe")
	}
	return nil
}
