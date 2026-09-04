package platform

import (
	"fmt"
	"time"
)

type MutationRunner interface {
	CommandRunner
	InputRunner
}

type InstallReport struct {
	Release    string           `json:"release"`
	Namespace  string           `json:"namespace"`
	Applied    bool             `json:"applied"`
	Verified   bool             `json:"verified"`
	Validation ServerValidation `json:"validation"`
}

func ApplyInstall(runner MutationRunner, namespace, release, chart, values string, timeout time.Duration) (InstallReport, error) {
	report := InstallReport{Release: release, Namespace: namespace}
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
	args := []string{"install", release, chart, "--namespace", namespace, "--create-namespace", "--wait"}
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
		return report, fmt.Errorf("installed release is not deployed")
	}
	return report, nil
}
