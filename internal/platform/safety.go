package platform

import "fmt"

const confirmationToken = "TRUSSIUM"

type RollbackVerification struct {
	Performed bool   `json:"performed"`
	Expected  string `json:"expected"`
	Reason    string `json:"reason"`
}

// RequireConfirmation is reserved for future mutating commands.
func RequireConfirmation(value string) error {
	if value != confirmationToken {
		return fmt.Errorf("confirmation must be %q", confirmationToken)
	}
	return nil
}
