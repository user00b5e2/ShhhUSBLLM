//go:build !windows

package main

import "os/exec"

func hideWindow(c *exec.Cmd) {}

func pidAliveWindows(pid int) bool { return false }
