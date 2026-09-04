package platform

import (
	"fmt"
	"regexp"
	"strconv"
)

type CompatibilityReport struct {
	Runtime    string   `json:"runtime"`
	Chart      string   `json:"chart"`
	Operator   string   `json:"operator"`
	Compatible bool     `json:"compatible"`
	Reasons    []string `json:"reasons,omitempty"`
}

var semverPattern = regexp.MustCompile(`^v?([0-9]+)\.([0-9]+)\.([0-9]+)(?:[-+].*)?$`)

type version struct{ major, minor, patch int }

func parseVersion(name, value string) (version, error) {
	match := semverPattern.FindStringSubmatch(value)
	if match == nil {
		return version{}, fmt.Errorf("%s version must be semantic version", name)
	}
	parts := [3]int{}
	for i := range parts {
		parsed, err := strconv.Atoi(match[i+1])
		if err != nil {
			return version{}, fmt.Errorf("%s version must be semantic version", name)
		}
		parts[i] = parsed
	}
	return version{parts[0], parts[1], parts[2]}, nil
}

func CheckCompatibility(runtime, chart, operator string) CompatibilityReport {
	report := CompatibilityReport{Runtime: runtime, Chart: chart, Operator: operator, Compatible: true}
	checks := []struct {
		name, value string
		min         version
	}{
		{"runtime", runtime, version{1, 22, 0}},
		{"chart", chart, version{1, 3, 0}},
		{"operator", operator, version{1, 0, 2}},
	}
	for _, check := range checks {
		parsed, err := parseVersion(check.name, check.value)
		if err != nil {
			report.Compatible = false
			report.Reasons = append(report.Reasons, err.Error())
			continue
		}
		if parsed.major != 1 || parsed.major < check.min.major || parsed.major == check.min.major && (parsed.minor < check.min.minor || parsed.minor == check.min.minor && parsed.patch < check.min.patch) {
			report.Compatible = false
			report.Reasons = append(report.Reasons, fmt.Sprintf("%s version is outside the supported 1.x baseline", check.name))
		}
	}
	return report
}
