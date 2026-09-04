package platform

import "fmt"

type ServerValidation struct {
	Performed bool   `json:"performed"`
	Valid     bool   `json:"valid"`
	Reason    string `json:"reason,omitempty"`
}

func ValidateManifest(runner InputRunner, manifest []byte) (ServerValidation, error) {
	if len(manifest) == 0 {
		return ServerValidation{}, fmt.Errorf("rendered manifest is empty")
	}
	if len(manifest) > maxCommandOutput {
		return ServerValidation{}, fmt.Errorf("rendered manifest is oversized")
	}
	if _, err := runner.RunInput("kubectl", manifest, "apply", "--dry-run=server", "--validate=true", "--filename", "-"); err != nil {
		return ServerValidation{Performed: true, Valid: false, Reason: err.Error()}, nil
	}
	return ServerValidation{Performed: true, Valid: true}, nil
}
