//go:build !windows

package main

import "os/exec"

func configureProcess(_ *exec.Cmd) {}
