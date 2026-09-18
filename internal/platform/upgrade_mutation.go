package platform

import (
	"fmt"
	"time"
)

type UpgradeReport struct {
	Release              string              `json:"release"`
	Namespace            string              `json:"namespace"`
	CurrentCompatibility CompatibilityReport `json:"currentCompatibility"`
	TargetCompatibility  CompatibilityReport `json:"targetCompatibility"`
	Validation           ServerValidation    `json:"validation"`
	Applied              bool                `json:"applied"`
	Verified             bool                `json:"verified"`
}

// GuardedUpgrade resolves the deployed baseline, checks confirmation, and then
// runs the validated upgrade workflow. Discovery and confirmation happen before
// any manifest rendering or mutation.
func GuardedUpgrade(runner MutationRunner, namespace, release, operator, chart, values, confirmation string, explicitCurrent, target [3]string, timeout time.Duration) (UpgradeReport, error) {
	var report UpgradeReport
	if err := RequireConfirmation(confirmation); err != nil {
		return report, err
	}
	current, err := ResolveCurrentVersions(runner, namespace, release, operator, explicitCurrent)
	if err != nil {
		return report, fmt.Errorf("resolve current component versions: %w", err)
	}
	return ApplyUpgrade(runner, namespace, release, chart, values, current, target, timeout)
}

func ApplyUpgrade(runner MutationRunner, namespace, release, chart, values string, current, target [3]string, timeout time.Duration) (UpgradeReport, error) {
	report := UpgradeReport{
		Release:              release,
		Namespace:            namespace,
		CurrentCompatibility: CheckCompatibility(current[0], current[1], current[2]),
		TargetCompatibility:  CheckCompatibility(target[0], target[1], target[2]),
	}
	if !report.CurrentCompatibility.Compatible || !report.TargetCompatibility.Compatible {
		return report, fmt.Errorf("component versions are incompatible")
	}
	manifest, err := RenderInstallManifest(runner, namespace, release, chart, values)
	if err != nil {
		return report, err
	}
	validation, err := ValidateManifest(runner, manifest)
	if err != nil {
		return report, err
	}
	report.Validation = validation
	if !validation.Valid {
		return report, fmt.Errorf("server-side validation failed")
	}
	args := []string{"upgrade", "--install", release, chart, "--namespace", namespace, "--create-namespace", "--wait"}
	if timeout > 0 {
		args = append(args, "--timeout", timeout.String())
	}
	if values != "" {
		args = append(args, "--values", values)
	}
	if _, err := runner.Run("helm", args...); err != nil {
		return report, err
	}
	report.Applied = true
	status, err := HelmStatusFor(runner, namespace, release)
	if err != nil {
		return report, err
	}
	report.Verified = status.Status == "deployed"
	if !report.Verified {
		return report, fmt.Errorf("upgraded release is not deployed")
	}
	return report, nil
}
