package platform

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

const maxCommandOutput = 64 * 1024

// CommandRunner is the narrow boundary around kubectl and helm.
type CommandRunner interface {
	Run(name string, args ...string) ([]byte, error)
}

type InputRunner interface {
	RunInput(name string, input []byte, args ...string) ([]byte, error)
}

type ExecRunner struct{ Timeout time.Duration }

func (r ExecRunner) Run(name string, args ...string) ([]byte, error) {
	timeout := r.Timeout
	if timeout <= 0 {
		timeout = 10 * time.Second
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	output, err := exec.CommandContext(ctx, name, args...).Output()
	if err != nil {
		return nil, fmt.Errorf("%s command failed", name)
	}
	if len(output) > maxCommandOutput {
		return nil, fmt.Errorf("%s returned oversized output", name)
	}
	return output, nil
}

func (r ExecRunner) RunInput(name string, input []byte, args ...string) ([]byte, error) {
	timeout := r.Timeout
	if timeout <= 0 {
		timeout = 10 * time.Second
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	command := exec.CommandContext(ctx, name, args...)
	command.Stdin = bytes.NewReader(input)
	output, err := command.Output()
	if err != nil {
		return nil, fmt.Errorf("%s command failed", name)
	}
	if len(output) > maxCommandOutput {
		return nil, fmt.Errorf("%s returned oversized output", name)
	}
	return output, nil
}

type OperatorStatus struct {
	Namespace string `json:"namespace"`
	Name      string `json:"name"`
	Ready     int32  `json:"readyReplicas"`
	Desired   int32  `json:"desiredReplicas"`
	Version   string `json:"version,omitempty"`
}

type HelmStatus struct {
	Namespace    string `json:"namespace"`
	Release      string `json:"release"`
	Status       string `json:"status"`
	Chart        string `json:"chart,omitempty"`
	ChartVersion string `json:"chartVersion,omitempty"`
	App          string `json:"appVersion,omitempty"`
}

type kubernetesDeployment struct {
	Metadata struct {
		Name string `json:"name"`
	} `json:"metadata"`
	Spec struct {
		Replicas int32 `json:"replicas"`
		Template struct {
			Spec struct {
				Containers []struct {
					Image string `json:"image"`
				} `json:"containers"`
			} `json:"spec"`
		} `json:"template"`
	} `json:"spec"`
	Status struct {
		Available int32 `json:"availableReplicas"`
	} `json:"status"`
}

func OperatorStatusFor(runner CommandRunner, namespace, name string) (OperatorStatus, error) {
	if strings.TrimSpace(namespace) == "" || strings.TrimSpace(name) == "" {
		return OperatorStatus{}, fmt.Errorf("namespace and name are required")
	}
	output, err := runner.Run("kubectl", "get", "deployment", name, "--namespace", namespace, "--output", "json")
	if err != nil {
		return OperatorStatus{}, err
	}
	var deployment kubernetesDeployment
	if err := json.Unmarshal(output, &deployment); err != nil || deployment.Metadata.Name == "" {
		return OperatorStatus{}, fmt.Errorf("kubectl returned invalid deployment")
	}
	return OperatorStatus{Namespace: namespace, Name: deployment.Metadata.Name, Ready: deployment.Status.Available, Desired: deployment.Spec.Replicas, Version: imageVersion(deployment.Spec.Template.Spec.Containers)}, nil
}

type helmRelease struct {
	Info struct {
		Status string `json:"status"`
	} `json:"info"`
	Chart struct {
		Metadata struct {
			Name    string `json:"name"`
			Version string `json:"version"`
		} `json:"metadata"`
	} `json:"chart"`
	Config struct {
		AppVersion string `json:"appVersion"`
	} `json:"config"`
}

func HelmStatusFor(runner CommandRunner, namespace, release string) (HelmStatus, error) {
	if strings.TrimSpace(namespace) == "" || strings.TrimSpace(release) == "" {
		return HelmStatus{}, fmt.Errorf("namespace and release are required")
	}
	output, err := runner.Run("helm", "status", release, "--namespace", namespace, "--output", "json")
	if err != nil {
		return HelmStatus{}, err
	}
	var status helmRelease
	if err := json.Unmarshal(output, &status); err != nil || status.Info.Status == "" {
		return HelmStatus{}, fmt.Errorf("helm returned invalid release")
	}
	chart := status.Chart.Metadata.Name
	if status.Chart.Metadata.Version != "" {
		chart += "-" + status.Chart.Metadata.Version
	}
	return HelmStatus{Namespace: namespace, Release: release, Status: status.Info.Status, Chart: chart, ChartVersion: status.Chart.Metadata.Version, App: status.Config.AppVersion}, nil
}

// ComponentVersions identifies the versions discovered from the running cluster.
type ComponentVersions struct {
	Runtime  string `json:"runtime,omitempty"`
	Chart    string `json:"chart,omitempty"`
	Operator string `json:"operator,omitempty"`
}

func imageVersion(containers []struct {
	Image string `json:"image"`
}) string {
	for _, container := range containers {
		image := strings.SplitN(container.Image, "@", 2)[0]
		if index := strings.LastIndex(image, ":"); index >= 0 && index+1 < len(image) {
			return image[index+1:]
		}
	}
	return ""
}

// DiscoverComponentVersions reads deployed versions without mutating the cluster.
func DiscoverComponentVersions(runner CommandRunner, namespace, release, operator string) (ComponentVersions, error) {
	helm, err := HelmStatusFor(runner, namespace, release)
	if err != nil {
		return ComponentVersions{}, err
	}
	operatorStatus, err := OperatorStatusFor(runner, namespace, operator)
	if err != nil {
		return ComponentVersions{}, err
	}
	versions := ComponentVersions{Runtime: helm.App, Chart: helm.ChartVersion, Operator: operatorStatus.Version}
	if versions.Runtime == "" || versions.Chart == "" || versions.Operator == "" {
		return ComponentVersions{}, fmt.Errorf("deployed component versions are incomplete")
	}
	return versions, nil
}
