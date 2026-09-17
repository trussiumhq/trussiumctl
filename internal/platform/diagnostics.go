package platform

type ClusterDiagnostics struct {
	Runtime       RuntimeStatus       `json:"runtime,omitempty"`
	Operator      OperatorStatus      `json:"operator,omitempty"`
	Helm          HelmStatus          `json:"helm,omitempty"`
	Versions      ComponentVersions   `json:"versions,omitempty"`
	Compatibility CompatibilityReport `json:"compatibility,omitempty"`
	Errors        []string            `json:"errors,omitempty"`
}

func CollectClusterDiagnostics(runtime RuntimeClient, runner CommandRunner, namespace, operator, release, runtimeVersion, chartVersion, operatorVersion string) ClusterDiagnostics {
	report := ClusterDiagnostics{}
	if status, err := runtime.Ready(); err != nil {
		report.Errors = append(report.Errors, err.Error())
	} else {
		report.Runtime = status
	}
	if status, err := OperatorStatusFor(runner, namespace, operator); err != nil {
		report.Errors = append(report.Errors, err.Error())
	} else {
		report.Operator = status
	}
	if status, err := HelmStatusFor(runner, namespace, release); err != nil {
		report.Errors = append(report.Errors, err.Error())
	} else {
		report.Helm = status
	}
	if runtimeVersion == "" {
		runtimeVersion = report.Helm.App
	}
	if chartVersion == "" {
		chartVersion = report.Helm.ChartVersion
	}
	if operatorVersion == "" {
		operatorVersion = report.Operator.Version
	}
	report.Versions = ComponentVersions{Runtime: runtimeVersion, Chart: chartVersion, Operator: operatorVersion}
	report.Compatibility = CheckCompatibility(runtimeVersion, chartVersion, operatorVersion)
	if !report.Compatibility.Compatible {
		report.Errors = append(report.Errors, report.Compatibility.Reasons...)
	}
	return report
}
