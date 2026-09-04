package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/trussiumhq/trussiumctl/internal/platform"
	"github.com/trussiumhq/trussiumctl/internal/version"
)

func main() {
	if len(os.Args) > 1 && (os.Args[1] == "version" || os.Args[1] == "--version") {
		fmt.Println(version.Version)
		return
	}
	if len(os.Args) > 1 && os.Args[1] == "runtime" {
		if runRuntime(os.Args[2:]) != nil {
			os.Exit(1)
		}
		return
	}
	if len(os.Args) > 1 && (os.Args[1] == "operator" || os.Args[1] == "helm") {
		if runClusterInspection(os.Args[1], os.Args[2:]) != nil {
			os.Exit(1)
		}
		return
	}
	if len(os.Args) > 1 && os.Args[1] == "compatibility" {
		if runCompatibility(os.Args[2:]) != nil {
			os.Exit(1)
		}
		return
	}

	fs := flag.NewFlagSet("trussiumctl", flag.ExitOnError)
	showVersion := fs.Bool("version", false, "print the client version")
	fs.Usage = func() {
		_, _ = fmt.Fprintln(fs.Output(), "trussiumctl manages Trussium Kubernetes and Helm workflows.")
		_, _ = fmt.Fprintln(fs.Output(), "Runtime execution remains owned by the Trussium Python service.")
		fs.PrintDefaults()
	}
	_ = fs.Parse(os.Args[1:])
	if *showVersion {
		fmt.Println(version.Version)
		return
	}
	fs.Usage()
}

func runCompatibility(args []string) error {
	if len(args) == 0 || args[0] != "check" {
		fmt.Fprintln(os.Stderr, "usage: trussiumctl compatibility check --runtime VERSION --chart VERSION --operator VERSION")
		return fmt.Errorf("invalid compatibility command")
	}
	fs := flag.NewFlagSet("compatibility check", flag.ContinueOnError)
	runtime := fs.String("runtime", "", "runtime semantic version")
	chart := fs.String("chart", "", "Helm chart semantic version")
	operator := fs.String("operator", "", "Operator semantic version")
	if err := fs.Parse(args[1:]); err != nil {
		return err
	}
	report := platform.CheckCompatibility(*runtime, *chart, *operator)
	if err := printJSON(report); err != nil {
		return err
	}
	if !report.Compatible {
		return fmt.Errorf("component versions are incompatible")
	}
	return nil
}

func runClusterInspection(kind string, args []string) error {
	fs := flag.NewFlagSet(kind+" status", flag.ContinueOnError)
	namespace := fs.String("namespace", "default", "Kubernetes namespace")
	name := fs.String("name", "trussium-operator", "deployment name")
	release := fs.String("release", "trussium", "Helm release name")
	if len(args) == 0 || args[0] != "status" {
		fmt.Fprintf(os.Stderr, "usage: trussiumctl %s status [flags]\n", kind)
		return fmt.Errorf("invalid %s command", kind)
	}
	if err := fs.Parse(args[1:]); err != nil {
		return err
	}
	runner := platform.ExecRunner{}
	if kind == "operator" {
		status, err := platform.OperatorStatusFor(runner, *namespace, *name)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return err
		}
		return printJSON(status)
	}
	status, err := platform.HelmStatusFor(runner, *namespace, *release)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return err
	}
	return printJSON(status)
}

func printJSON(value any) error {
	encoded, err := json.Marshal(value)
	if err != nil {
		return err
	}
	_, _ = fmt.Println(string(encoded))
	return nil
}

func runRuntime(args []string) error {
	if len(args) == 0 || args[0] != "status" {
		fmt.Fprintln(os.Stderr, "usage: trussiumctl runtime status [--url URL]")
		return fmt.Errorf("invalid runtime command")
	}
	fs := flag.NewFlagSet("runtime status", flag.ContinueOnError)
	url := fs.String("url", "http://127.0.0.1:9000", "runtime base URL")
	if err := fs.Parse(args[1:]); err != nil {
		return err
	}
	client := platform.HTTPRuntimeClient{BaseURL: *url, Client: &http.Client{Timeout: 5 * time.Second}}
	status, err := client.Ready()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return err
	}
	_, _ = fmt.Printf("{\"status\":%q}\n", status.Status)
	return nil
}
