package platform

import "testing"

func TestRequireConfirmation(t *testing.T) {
	if err := RequireConfirmation("TRUSSIUM"); err != nil {
		t.Fatal(err)
	}
	if err := RequireConfirmation("yes"); err == nil {
		t.Fatal("expected exact confirmation token")
	}
}
