package platform

import (
	"fmt"
	"strings"
)

type InstallDryRunReport struct {
	Release       string `json:"release"`
	Namespace     string `json:"namespace"`
	Chart         string `json:"chart"`
	Rendered      bool   `json:"rendered"`
	ManifestBytes int    `json:"manifestBytes"`
	ResourceCount int    `json:"resourceCount"`
}

func RenderInstall(runner CommandRunner, namespace, release, chart, values string) (InstallDryRunReport, error) {
	output, err := RenderInstallManifest(runner, namespace, release, chart, values)
	if err != nil {
		return InstallDryRunReport{}, err
	}
	resources := 0
	for _, line := range strings.Split(string(output), "\n") {
		if strings.HasPrefix(line, "kind:") {
			resources++
		}
	}
	return InstallDryRunReport{Release: release, Namespace: namespace, Chart: chart, Rendered: true, ManifestBytes: len(output), ResourceCount: resources}, nil
}

func RenderInstallManifest(runner CommandRunner, namespace, release, chart, values string) ([]byte, error) {
	if strings.TrimSpace(namespace) == "" || strings.TrimSpace(release) == "" || strings.TrimSpace(chart) == "" {
		return nil, fmt.Errorf("namespace, release, and chart are required")
	}
	args := []string{"template", release, chart, "--namespace", namespace}
	if strings.TrimSpace(values) != "" {
		args = append(args, "--values", values)
	}
	output, err := runner.Run("helm", args...)
	if err != nil {
		return nil, err
	}
	if len(output) > maxCommandOutput {
		return nil, fmt.Errorf("helm returned oversized output")
	}
	return output, nil
}
