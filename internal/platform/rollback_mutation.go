package platform

import (
	"fmt"
	"time"
)

type RollbackReport struct {
	Release             string              `json:"release"`
	Namespace           string              `json:"namespace"`
	Revision            int                 `json:"revision"`
	TargetCompatibility CompatibilityReport `json:"targetCompatibility"`
	Validation          ServerValidation    `json:"validation"`
	Applied             bool                `json:"applied"`
	Verified            bool                `json:"verified"`
}

func ApplyRollback(runner MutationRunner, namespace, release, chart, values string, revision int, target [3]string, timeout time.Duration) (RollbackReport, error) {
	report := RollbackReport{Release: release, Namespace: namespace, Revision: revision, TargetCompatibility: CheckCompatibility(target[0], target[1], target[2])}
	if !report.TargetCompatibility.Compatible {
		return report, fmt.Errorf("target component versions are incompatible")
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
	args := []string{"rollback", release}
	if revision > 0 {
		args = append(args, fmt.Sprintf("%d", revision))
	}
	args = append(args, "--namespace", namespace, "--wait")
	if timeout > 0 {
		args = append(args, "--timeout", timeout.String())
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
		return report, fmt.Errorf("rolled back release is not deployed")
	}
	return report, nil
}
