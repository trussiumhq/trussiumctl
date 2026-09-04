package main

import (
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
