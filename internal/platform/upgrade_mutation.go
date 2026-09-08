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
