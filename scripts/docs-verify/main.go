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
	flags := flag.NewFlagSet("docs-verify", flag.ContinueOnError)
	flags.SetOutput(stderr)
	repo := flags.String("repo", "../..", "repository root")
	if err := flags.Parse(args); err != nil {
		return 2
	}
	if flags.NArg() != 0 {
		fmt.Fprintln(stderr, "docs-verify: unexpected positional arguments")
		return 2
	}

	result, err := verifyRepository(*repo)
	if err != nil {
		fmt.Fprintf(stderr, "docs-verify: %v\n", err)
		return 2
	}
	if len(result.violations) != 0 {
		for _, item := range violationStrings(result.violations) {
			fmt.Fprintln(stderr, item)
		}
		fmt.Fprintf(stderr, "docs-verify: %d violation(s)\n", len(result.violations))
		return 1
	}

	fmt.Fprintf(stdout, "docs-verify: checked %d documents, %d local links, %d registered examples, and %d format rows\n", result.documents, result.localLinks, result.registeredExamples, result.formatRows)
	return 0
}
