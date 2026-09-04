package platform

import (
	"reflect"
	"testing"
)

type inputRunner struct {
	name  string
	args  []string
	input []byte
}

func (r *inputRunner) RunInput(name string, input []byte, args ...string) ([]byte, error) {
	r.name, r.args, r.input = name, args, input
	return nil, nil
}

func TestValidateManifestUsesServerDryRun(t *testing.T) {
	runner := &inputRunner{}
	result, err := ValidateManifest(runner, []byte("apiVersion: v1\nkind: ConfigMap\n"))
	if err != nil || !result.Valid || runner.name != "kubectl" {
		t.Fatalf("result=%+v runner=%+v err=%v", result, runner, err)
	}
	want := []string{"apply", "--dry-run=server", "--validate=true", "--filename", "-"}
	if !reflect.DeepEqual(runner.args, want) || len(runner.input) == 0 {
		t.Fatalf("args=%v input=%q", runner.args, runner.input)
	}
}
