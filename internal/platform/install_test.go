package platform

import (
	"reflect"
	"testing"
)

func TestRenderInstallUsesHelmTemplateOnly(t *testing.T) {
	runner := &fakeRunner{out: []byte("apiVersion: v1\nkind: ConfigMap\n---\nkind: Service\n")}
	report, err := RenderInstall(runner, "trussium-system", "trussium", "trussium/trussium", "values.yaml")
	if err != nil || !report.Rendered || report.ResourceCount != 2 {
		t.Fatalf("report=%+v err=%v", report, err)
	}
	want := []string{"template", "trussium", "trussium/trussium", "--namespace", "trussium-system", "--values", "values.yaml"}
	if runner.name != "helm" || !reflect.DeepEqual(runner.args, want) {
		t.Fatalf("command=%s %v", runner.name, runner.args)
	}
}
