// Command brand-verify proves the retained Cueson brand archive, payload, and
// repository references match the import manifest without writing or fetching.
package main

import (
	"flag"
	"fmt"
	"io"
	"os"
)

func main() {
	os.Exit(runCLI(os.Args[1:], os.Stdout, os.Stderr))
}

func runCLI(args []string, stdout, stderr io.Writer) int {
	flags := flag.NewFlagSet("brand-verify", flag.ContinueOnError)
	flags.SetOutput(stderr)
	repo := flags.String("repo", "../..", "repository root")
	manifest := flags.String("manifest", defaultManifestPath, "repository-relative import manifest path")
	if err := flags.Parse(args); err != nil {
		return 2
	}
	if flags.NArg() != 0 {
		fmt.Fprintln(stderr, "brand-verify: unexpected positional arguments")
		return 2
	}

	result, err := verifyRepository(*repo, *manifest)
	if err != nil {
		fmt.Fprintf(stderr, "brand-verify: %v\n", err)
		return 2
	}
	if len(result.violations) != 0 {
		for _, violation := range result.violations {
			fmt.Fprintln(stderr, violation)
		}
		fmt.Fprintf(stderr, "brand-verify: %d violation(s)\n", len(result.violations))
		return 1
	}

	fmt.Fprintf(stdout, "brand-verify: verified %d entries and %d references\n", result.entries, result.references)
	return 0
}
