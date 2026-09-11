//go:build !windows

package cli

import "os/exec"

func configureTestCommand(_ *exec.Cmd) {}
