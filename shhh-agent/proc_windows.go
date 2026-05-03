//go:build windows

package main

import (
	"os/exec"
	"syscall"
)

const (
	createNoWindow    = 0x08000000
	detachedProcess   = 0x00000008
	createNewProcGrp  = 0x00000200
)

func hideWindow(c *exec.Cmd) {
	c.SysProcAttr = &syscall.SysProcAttr{
		HideWindow:    true,
		CreationFlags: createNoWindow,
	}
}

// pidAliveWindows checks if a PID exists by opening it with PROCESS_QUERY_LIMITED_INFORMATION.
func pidAliveWindows(pid int) bool {
	const PROCESS_QUERY_LIMITED_INFORMATION = 0x1000
	h, err := syscall.OpenProcess(PROCESS_QUERY_LIMITED_INFORMATION, false, uint32(pid))
	if err != nil {
		return false
	}
	defer syscall.CloseHandle(h)
	var code uint32
	if err := syscall.GetExitCodeProcess(h, &code); err != nil {
		return false
	}
	return code == 259 // STILL_ACTIVE
}
