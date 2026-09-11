package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type stringList []string

func (values *stringList) String() string { return fmt.Sprint([]string(*values)) }
func (values *stringList) Set(value string) error {
	*values = append(*values, value)
	return nil
}

func main() {
	os.Exit(run(context.Background(), os.Args[1:], os.Stdout, os.Stderr))
}

func run(ctx context.Context, args []string, stdout, stderr io.Writer) int {
	flags := flag.NewFlagSet("release-verify", flag.ContinueOnError)
	flags.SetOutput(io.Discard)
	var config Config
	var forbidden stringList
	flags.StringVar(&config.DistDir, "dist", "", "GoReleaser distribution directory")
	flags.StringVar(&config.RepoDir, "repo", "", "repository root")
	flags.StringVar(&config.Version, "version", "", "expected release version")
	flags.StringVar(&config.Commit, "commit", "", "expected full source revision")
	flags.StringVar(&config.EvidencePath, "evidence", "", "distinct output path for accepted release evidence")
	flags.BoolVar(&config.ExecuteHost, "execute-host", false, "execute the compatible packaged binary")
	flags.BoolVar(&config.Development, "development", false, "verify the evolving canonical schema without requiring an immutable release copy")
	flags.Var(&forbidden, "forbid", "additional local identifier to reject (repeatable)")
	if err := flags.Parse(args); err != nil {
		fmt.Fprintf(stderr, "release verification: invalid arguments: %v\n", err)
		return 2
	}
	if flags.NArg() != 0 {
		fmt.Fprintf(stderr, "release verification: unexpected arguments: %v\n", flags.Args())
		return 2
	}
	config.Forbidden = forbidden
	if config.DistDir == "" || config.RepoDir == "" || config.Version == "" || config.Commit == "" {
		fmt.Fprintln(stderr, "release verification: -dist, -repo, -version, and -commit are required")
		return 2
	}
	evidence, err := Verify(ctx, config)
	if err != nil {
		fmt.Fprintf(stderr, "release verification: %v\n", err)
		return 1
	}
	path := evidenceOutputPath(config)
	if err := writeEvidence(path, evidence); err != nil {
		fmt.Fprintf(stderr, "release verification: write evidence: %v\n", err)
		return 1
	}
	proofKind := "non-publishing release candidate"
	if evidence.Development {
		proofKind = "non-publishing development snapshot"
	}
	fmt.Fprintf(stdout, "verified %s %s at %s\n", proofKind, evidence.Version, evidence.SourceRevision)
	fmt.Fprintf(stdout, "evidence: %s\n", path)
	return 0
}

func evidenceOutputPath(config Config) string {
	if config.EvidencePath != "" {
		return config.EvidencePath
	}
	return filepath.Join(config.DistDir, evidenceFilename)
}

func writeEvidence(path string, evidence ReleaseEvidence) error {
	data, err := json.MarshalIndent(evidence, "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o644)
	if err != nil {
		if errors.Is(err, os.ErrExist) {
			return fmt.Errorf("%s already exists; choose a distinct evidence output path", path)
		}
		return err
	}
	written := false
	defer func() {
		_ = file.Close()
		if !written {
			_ = os.Remove(path)
		}
	}()
	if _, err := file.Write(data); err != nil {
		return err
	}
	if err := file.Sync(); err != nil {
		return err
	}
	if err := file.Close(); err != nil {
		return err
	}
	written = true
	return nil
}
