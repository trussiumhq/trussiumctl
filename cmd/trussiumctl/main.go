package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/trussiumhq/trussiumctl/internal/version"
)

func main() {
	if len(os.Args) > 1 && os.Args[1] == "version" {
		fmt.Println(version.Version)
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
