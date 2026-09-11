package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
)

func main() {
	root := flag.String("root", "", "private corpus directory to verify")
	jsonOutput := flag.Bool("json", false, "emit one aggregate JSON report")
	flag.Parse()
	if *root == "" || flag.NArg() != 0 {
		fmt.Fprintln(os.Stderr, "usage: corpus-verify -root DIRECTORY [-json]")
		os.Exit(2)
	}
	report, err := Verify(context.Background(), *root)
	if err != nil {
		fmt.Fprintln(os.Stderr, "corpus verification failed:", err)
		os.Exit(1)
	}
	if *jsonOutput {
		if err := json.NewEncoder(os.Stdout).Encode(report); err != nil {
			fmt.Fprintln(os.Stderr, "write corpus report")
			os.Exit(1)
		}
		return
	}
	fmt.Printf("verified files=%d accepted=%d rejected=%d unsupported=%d srt=%d vtt=%d\n", report.Files, report.Accepted, report.Rejected, report.Unsupported, report.Formats.SRT, report.Formats.VTT)
	for _, failure := range report.Failures {
		fmt.Printf("%s: %s\n", failure.RelativeIdentity, failure.Result)
	}
}
