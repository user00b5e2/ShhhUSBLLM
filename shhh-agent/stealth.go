package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/user"
	"path/filepath"
	"runtime"
	"strings"
)

// errInterrupted signals that the user pressed Ctrl+C / Ctrl+D during input.
var errInterrupted = errors.New("interrupted")

// ANSI sequences shared across platforms (Windows 10+ supports them).
const (
	ansiConcealOn  = "\x1b[8m"
	ansiConcealOff = "\x1b[0m"
	ansiClear      = "\x1b[2J\x1b[H"
)

// ShellKind is what the parent shell appeared to be — drives the fake prompt.
type ShellKind int

const (
	ShellCMD ShellKind = iota
	ShellPowerShell
	ShellZsh
	ShellBash
)

// ResolvePrompt returns the prompt string to display, using the 3-tier cascade:
//  1. SHHH_PROMPT — explicit override; used verbatim. For users with custom themes
//     (oh-my-posh, posh-git) who want pixel-perfect mimicry.
//  2. SHHH_FAKE_PROMPT — captured by update.bat/update.ps1 from the real shell's
//     `prompt` function before launching. Reflects the real PS prompt at launch time.
//  3. Fallback — generic templated prompt based on detected shell + cwd.
//
// A trailing space is enforced to match real shell ergonomics.
func ResolvePrompt(k ShellKind) string {
	if p := os.Getenv("SHHH_PROMPT"); p != "" {
		return ensureTrailingSpace(p)
	}
	if p := os.Getenv("SHHH_FAKE_PROMPT"); p != "" {
		return ensureTrailingSpace(p)
	}
	return FakePrompt(k)
}

func ensureTrailingSpace(s string) string {
	if !strings.HasSuffix(s, " ") {
		return s + " "
	}
	return s
}

// FakePrompt returns a generic prompt that mimics the shell's idle prompt.
// Used as last-resort fallback when no captured/override prompt is available.
func FakePrompt(k ShellKind) string {
	cwd, _ := os.Getwd()
	switch k {
	case ShellCMD:
		// CMD shows the absolute path; on Windows convert to backslashes.
		p := filepath.FromSlash(cwd)
		return p + "> "
	case ShellPowerShell:
		p := filepath.FromSlash(cwd)
		return "PS " + p + "> "
	case ShellZsh, ShellBash:
		u, _ := user.Current()
		host, _ := os.Hostname()
		short := strings.SplitN(host, ".", 2)[0]
		home := ""
		if u != nil {
			home = u.HomeDir
		}
		display := cwd
		if home != "" && strings.HasPrefix(cwd, home) {
			display = "~" + strings.TrimPrefix(cwd, home)
		}
		uname := "user"
		if u != nil {
			uname = u.Username
		}
		if k == ShellZsh {
			return fmt.Sprintf("%s@%s %s %% ", uname, short, display)
		}
		return fmt.Sprintf("%s@%s:%s$ ", uname, short, display)
	}
	return "$ "
}

// DetectShell heuristically guesses the parent shell.
func DetectShell() ShellKind {
	if runtime.GOOS == "windows" {
		// PSModulePath is set by PowerShell host; cmd.exe doesn't set it.
		if os.Getenv("PSModulePath") != "" && os.Getenv("PSExecutionPolicyPreference") != "" {
			return ShellPowerShell
		}
		// VS Code's integrated terminal exposes TERM_PROGRAM.
		if os.Getenv("TERM_PROGRAM") == "vscode" && strings.Contains(strings.ToLower(os.Getenv("PSModulePath")), "powershell") {
			return ShellPowerShell
		}
		return ShellCMD
	}
	sh := os.Getenv("SHELL")
	if strings.HasSuffix(sh, "zsh") {
		return ShellZsh
	}
	return ShellBash
}

// ConcealStart hides everything written to w from now on.
func ConcealStart(w io.Writer) { fmt.Fprint(w, ansiConcealOn) }

// ConcealEnd restores normal rendering.
func ConcealEnd(w io.Writer) { fmt.Fprint(w, ansiConcealOff) }

// ClearScreen wipes the terminal.
func ClearScreen(w io.Writer) { fmt.Fprint(w, ansiClear) }
