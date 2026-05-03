package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func setupTools(t *testing.T) (*Tools, string) {
	t.Helper()
	dir := t.TempDir()
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	tt, err := NewTools(false)
	if err != nil {
		t.Fatal(err)
	}
	return tt, dir
}

func TestPathTraversalBlocked(t *testing.T) {
	tt, _ := setupTools(t)
	args, _ := json.Marshal(map[string]string{"path": "../etc/passwd"})
	if _, err := tt.ReadFile(args); err == nil {
		t.Fatal("expected path traversal to fail")
	}
}

func TestProtectedNameBlocked(t *testing.T) {
	tt, dir := setupTools(t)
	_ = os.Mkdir(filepath.Join(dir, ".git"), 0o755)
	_ = os.WriteFile(filepath.Join(dir, ".git", "HEAD"), []byte("x"), 0o644)
	args, _ := json.Marshal(map[string]string{"path": ".git/HEAD"})
	if _, err := tt.ReadFile(args); err == nil {
		t.Fatal("expected .git read to be blocked")
	}
}

func TestEditFileNotUnique(t *testing.T) {
	tt, dir := setupTools(t)
	p := filepath.Join(dir, "x.txt")
	_ = os.WriteFile(p, []byte("foo\nfoo\n"), 0o644)
	args, _ := json.Marshal(map[string]string{"path": "x.txt", "old": "foo", "new": "bar"})
	_, err := tt.EditFile(args)
	if err == nil || !strings.Contains(err.Error(), "not unique") {
		t.Fatalf("expected not-unique error, got %v", err)
	}
}

func TestEditFileSuccess(t *testing.T) {
	tt, dir := setupTools(t)
	p := filepath.Join(dir, "y.txt")
	_ = os.WriteFile(p, []byte("hello world"), 0o644)
	args, _ := json.Marshal(map[string]string{"path": "y.txt", "old": "world", "new": "shhh"})
	if _, err := tt.EditFile(args); err != nil {
		t.Fatal(err)
	}
	b, _ := os.ReadFile(p)
	if string(b) != "hello shhh" {
		t.Fatalf("got %q", string(b))
	}
}

func TestRunCmdBlocksDangerous(t *testing.T) {
	tt, _ := setupTools(t)
	args, _ := json.Marshal(map[string]string{"cmd": "rm -rf /"})
	if _, err := tt.RunCmd(args); err == nil {
		t.Fatal("expected rm -rf to be blocked")
	}
}

func TestParseToolXML(t *testing.T) {
	got, err := parseTool(`<tool>read_file</tool><args>{"path":"a"}</args>`)
	if err != nil || got.Name != "read_file" || string(got.Args) != `{"path":"a"}` {
		t.Fatalf("parse failed: %+v err=%v", got, err)
	}
}
