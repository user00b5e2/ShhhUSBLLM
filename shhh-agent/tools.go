package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"time"
)

const (
	maxFileBytes = 1024 * 1024 // 1 MB — fits multi-page MD specs
	maxListItems = 200
	cmdTimeout   = 120 * time.Second
)

var protectedNames = []string{
	".git", ".svn", ".hg",
	"node_modules", "vendor", "target", "dist", "build", "__pycache__",
	".venv", "venv", ".tox",
}

var protectedFilePatterns = []string{
	".env", ".env.local", ".env.production", ".env.development",
	"id_rsa", "id_dsa", "id_ecdsa", "id_ed25519",
	".npmrc", ".pypirc", ".netrc",
}

var protectedSuffixes = []string{".key", ".pem", ".pfx", ".p12", ".keystore"}

var dangerousCmds = []string{
	"format", "fdisk", "diskpart",
	"shutdown", "logoff", "reboot", "halt",
	"reg ", "regedit",
	"net user", "net localgroup",
	"rmdir /s", "rd /s", "rm -rf",
	"del /s", "del /q /s",
	"chmod 777 /", "chown -r",
	"mkfs", "dd if=", "dd of=/dev",
	"curl ", "wget ",                // exfil prevention by default
	"powershell ", "pwsh ", "iex ",  // nested PS / invoke-expression
	"cmd /c", "cmd.exe /c",
	"start ", "rundll32",
}

// Tools holds the runtime state shared by every tool call.
type Tools struct {
	cwd    string // canonical absolute root, the only writable scope
	unsafe bool   // disables CWD scoping when true (--unsafe flag)
}

func NewTools(unsafe bool) (*Tools, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	abs, err := filepath.Abs(wd)
	if err != nil {
		return nil, err
	}
	resolved, err := filepath.EvalSymlinks(abs)
	if err != nil {
		resolved = abs
	}
	return &Tools{cwd: resolved, unsafe: unsafe}, nil
}

func (t *Tools) resolve(p string) (string, error) {
	if p == "" {
		return "", errors.New("empty path")
	}
	if !filepath.IsAbs(p) {
		p = filepath.Join(t.cwd, p)
	}
	abs, err := filepath.Abs(p)
	if err != nil {
		return "", err
	}
	// EvalSymlinks fails if the file does not exist yet (write_file case).
	// Walk up until we find a parent that exists, resolve that, re-append.
	resolved, err := filepath.EvalSymlinks(abs)
	if err != nil {
		parent := abs
		var tail []string
		for {
			p2, e2 := filepath.EvalSymlinks(parent)
			if e2 == nil {
				resolved = filepath.Join(append([]string{p2}, tail...)...)
				break
			}
			tail = append([]string{filepath.Base(parent)}, tail...)
			next := filepath.Dir(parent)
			if next == parent {
				resolved = abs
				break
			}
			parent = next
		}
	}

	if !t.unsafe {
		rel, err := filepath.Rel(t.cwd, resolved)
		if err != nil || strings.HasPrefix(rel, "..") || rel == ".." {
			return "", fmt.Errorf("path outside workspace: %s", p)
		}
	}
	return resolved, nil
}

func isProtected(p string) error {
	parts := strings.Split(filepath.ToSlash(p), "/")
	for _, seg := range parts {
		for _, bad := range protectedNames {
			if strings.EqualFold(seg, bad) {
				return fmt.Errorf("protected directory: %s", bad)
			}
		}
	}
	base := filepath.Base(p)
	for _, pat := range protectedFilePatterns {
		if strings.EqualFold(base, pat) || strings.HasPrefix(strings.ToLower(base), strings.ToLower(pat)+".") {
			return fmt.Errorf("protected file: %s", base)
		}
	}
	low := strings.ToLower(base)
	for _, s := range protectedSuffixes {
		if strings.HasSuffix(low, s) {
			return fmt.Errorf("protected suffix: %s", s)
		}
	}
	return nil
}

// --- read_file ---

type readArgs struct{ Path string `json:"path"` }

func (t *Tools) ReadFile(raw json.RawMessage) (string, error) {
	var a readArgs
	if err := json.Unmarshal(raw, &a); err != nil {
		return "", fmt.Errorf("invalid args: %w", err)
	}
	abs, err := t.resolve(a.Path)
	if err != nil {
		return "", err
	}
	if err := isProtected(abs); err != nil && !t.unsafe {
		return "", err
	}
	st, err := os.Stat(abs)
	if err != nil {
		return "", err
	}
	if st.IsDir() {
		return "", fmt.Errorf("is a directory: %s", a.Path)
	}
	if st.Size() > maxFileBytes {
		return "", fmt.Errorf("file too large (%d bytes, max %d)", st.Size(), maxFileBytes)
	}
	b, err := os.ReadFile(abs)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// --- write_file ---

type writeArgs struct {
	Path    string `json:"path"`
	Content string `json:"content"`
}

func (t *Tools) WriteFile(raw json.RawMessage) (string, error) {
	var a writeArgs
	if err := json.Unmarshal(raw, &a); err != nil {
		return "", fmt.Errorf("invalid args: %w", err)
	}
	abs, err := t.resolve(a.Path)
	if err != nil {
		return "", err
	}
	if err := isProtected(abs); err != nil && !t.unsafe {
		return "", err
	}
	if err := os.MkdirAll(filepath.Dir(abs), 0o755); err != nil {
		return "", err
	}
	if err := os.WriteFile(abs, []byte(a.Content), 0o644); err != nil {
		return "", err
	}
	return fmt.Sprintf("wrote %d bytes to %s", len(a.Content), a.Path), nil
}

// --- edit_file ---

type editArgs struct {
	Path string `json:"path"`
	Old  string `json:"old"`
	New  string `json:"new"`
}

func (t *Tools) EditFile(raw json.RawMessage) (string, error) {
	var a editArgs
	if err := json.Unmarshal(raw, &a); err != nil {
		return "", fmt.Errorf("invalid args: %w", err)
	}
	if a.Old == "" {
		return "", errors.New("'old' must be non-empty")
	}
	abs, err := t.resolve(a.Path)
	if err != nil {
		return "", err
	}
	if err := isProtected(abs); err != nil && !t.unsafe {
		return "", err
	}
	b, err := os.ReadFile(abs)
	if err != nil {
		return "", err
	}
	src := string(b)
	count := strings.Count(src, a.Old)
	if count == 0 {
		return "", fmt.Errorf("'old' string not found in %s — read_file first", a.Path)
	}
	if count > 1 {
		return "", fmt.Errorf("'old' string is not unique (%d matches) — add more context", count)
	}
	out := strings.Replace(src, a.Old, a.New, 1)
	if err := os.WriteFile(abs, []byte(out), 0o644); err != nil {
		return "", err
	}
	return fmt.Sprintf("edited %s (%d→%d bytes)", a.Path, len(src), len(out)), nil
}

// --- list_dir ---

type listArgs struct{ Path string `json:"path"` }

func (t *Tools) ListDir(raw json.RawMessage) (string, error) {
	var a listArgs
	if err := json.Unmarshal(raw, &a); err != nil {
		return "", fmt.Errorf("invalid args: %w", err)
	}
	if a.Path == "" {
		a.Path = "."
	}
	abs, err := t.resolve(a.Path)
	if err != nil {
		return "", err
	}
	entries, err := os.ReadDir(abs)
	if err != nil {
		return "", err
	}
	var names []string
	for _, e := range entries {
		name := e.Name()
		skip := false
		for _, bad := range protectedNames {
			if strings.EqualFold(name, bad) {
				skip = true
				break
			}
		}
		if skip {
			continue
		}
		if e.IsDir() {
			name += "/"
		}
		names = append(names, name)
		if len(names) >= maxListItems {
			names = append(names, "... (truncated)")
			break
		}
	}
	sort.Strings(names)
	return strings.Join(names, "\n"), nil
}

// --- run_cmd ---

type runArgs struct{ Cmd string `json:"cmd"` }

func (t *Tools) RunCmd(raw json.RawMessage) (string, error) {
	var a runArgs
	if err := json.Unmarshal(raw, &a); err != nil {
		return "", fmt.Errorf("invalid args: %w", err)
	}
	cmd := strings.TrimSpace(a.Cmd)
	if cmd == "" {
		return "", errors.New("empty command")
	}
	low := strings.ToLower(cmd)
	if !t.unsafe {
		for _, bad := range dangerousCmds {
			if strings.Contains(low, bad) {
				return "", fmt.Errorf("blocked dangerous pattern: %q", bad)
			}
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), cmdTimeout)
	defer cancel()

	var c *exec.Cmd
	if runtime.GOOS == "windows" {
		c = exec.CommandContext(ctx, "cmd.exe", "/c", cmd)
	} else {
		c = exec.CommandContext(ctx, "sh", "-c", cmd)
	}
	c.Dir = t.cwd
	out, err := c.CombinedOutput()
	result := string(out)
	if len(result) > maxFileBytes {
		result = result[:maxFileBytes] + "\n... (truncated)"
	}
	if ctx.Err() == context.DeadlineExceeded {
		return result, fmt.Errorf("timeout after %s", cmdTimeout)
	}
	if err != nil {
		return result, fmt.Errorf("exit: %w", err)
	}
	if result == "" {
		result = "(no output)"
	}
	return result, nil
}
