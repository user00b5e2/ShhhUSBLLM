package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
)

// ServerConfig describes a llama-server invocation.
type ServerConfig struct {
	BinaryPath string // path to hostcfg(.exe) — llama-server renamed
	ModelPath  string
	Host       string
	Port       int
	CtxSize    int
	Threads    int
}

// runningServer holds the in-process child so we can stop it on exit.
var runningServer *exec.Cmd

// LockInfo persisted in %TEMP%/hostcfg.lock so subsequent invocations
// can find or kill an existing server.
type LockInfo struct {
	PID       int
	ModelPath string
	Port      int
}

func lockPath() string {
	return filepath.Join(os.TempDir(), "hostcfg.lock")
}

func readLock() (*LockInfo, error) {
	b, err := os.ReadFile(lockPath())
	if err != nil {
		return nil, err
	}
	parts := strings.SplitN(strings.TrimSpace(string(b)), "\n", 3)
	if len(parts) < 3 {
		return nil, fmt.Errorf("malformed lock")
	}
	pid, err := strconv.Atoi(parts[0])
	if err != nil {
		return nil, err
	}
	port, err := strconv.Atoi(parts[2])
	if err != nil {
		return nil, err
	}
	return &LockInfo{PID: pid, ModelPath: parts[1], Port: port}, nil
}

func writeLock(li *LockInfo) error {
	body := fmt.Sprintf("%d\n%s\n%d", li.PID, li.ModelPath, li.Port)
	return os.WriteFile(lockPath(), []byte(body), 0o600)
}

func removeLock() { _ = os.Remove(lockPath()) }

// processAlive is best-effort: signal 0 on Unix, FindProcess+exit on Windows.
func processAlive(pid int) bool {
	p, err := os.FindProcess(pid)
	if err != nil {
		return false
	}
	if runtime.GOOS == "windows" {
		// On Windows, FindProcess always succeeds; verify with a no-op signal.
		// The simplest cross-version check: try to read the exit code via a kill -0 emulation.
		// We use OpenProcess implicitly through Signal which is not implemented for Windows.
		// Fall back to listing tasklist would be heavier — accept best-effort here.
		_ = p
		return pidAliveWindows(pid)
	}
	return p.Signal(nil) == nil
}

func killPID(pid int) {
	p, err := os.FindProcess(pid)
	if err != nil {
		return
	}
	_ = p.Kill()
}

// EnsureServer guarantees a llama-server is running for the requested model.
// If another instance (from a previous run or another shhh-agent) is alive
// with the right model, we reuse it. Otherwise we spawn one as our child;
// it will die with us on Stop / exit.
func EnsureServer(cfg ServerConfig) error {
	if li, err := readLock(); err == nil && processAlive(li.PID) {
		if li.ModelPath == cfg.ModelPath && li.Port == cfg.Port && healthy(cfg.Host, cfg.Port) {
			return nil
		}
		killPID(li.PID)
		removeLock()
		for i := 0; i < 20 && portInUse(cfg.Host, cfg.Port); i++ {
			time.Sleep(100 * time.Millisecond)
		}
	} else {
		removeLock()
	}

	args := []string{
		"-m", cfg.ModelPath,
		"-c", strconv.Itoa(cfg.CtxSize),
		"--host", cfg.Host,
		"--port", strconv.Itoa(cfg.Port),
		"--log-disable",
		"-t", strconv.Itoa(cfg.Threads),
		"--no-mmap",
	}
	cmd := exec.Command(cfg.BinaryPath, args...)
	// stdin: nil → connected to os.DevNull internally by exec without RDWR oddities.
	cmd.Stdin = nil
	if os.Getenv("SHHH_DEBUG_SERVER") != "" {
		logf, _ := os.Create(filepath.Join(os.TempDir(), "hostcfg.debug.log"))
		cmd.Stdout = logf
		cmd.Stderr = logf
	} else {
		devnullW, _ := os.OpenFile(os.DevNull, os.O_WRONLY, 0)
		cmd.Stdout = devnullW
		cmd.Stderr = devnullW
	}
	hideWindow(cmd)
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("start hostcfg %s: %w", cfg.BinaryPath, err)
	}
	runningServer = cmd
	if err := writeLock(&LockInfo{PID: cmd.Process.Pid, ModelPath: cfg.ModelPath, Port: cfg.Port}); err != nil {
		return err
	}
	return nil
}

// WaitHealthy polls /health until 200 OK or timeout.
func WaitHealthy(host string, port int, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		if healthy(host, port) {
			return nil
		}
		time.Sleep(500 * time.Millisecond)
	}
	return fmt.Errorf("server not healthy after %s", timeout)
}

func healthy(host string, port int) bool {
	c := http.Client{Timeout: 1 * time.Second}
	resp, err := c.Get(fmt.Sprintf("http://%s:%d/health", host, port))
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	return resp.StatusCode == 200
}

func portInUse(host string, port int) bool {
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", host, port), 200*time.Millisecond)
	if err != nil {
		return false
	}
	_ = conn.Close()
	return true
}

// StopServer kills the running hostcfg if any. Safe to call multiple times.
func StopServer() {
	if runningServer != nil && runningServer.Process != nil {
		_ = runningServer.Process.Kill()
		_, _ = runningServer.Process.Wait()
		runningServer = nil
	}
	if li, err := readLock(); err == nil && processAlive(li.PID) {
		killPID(li.PID)
	}
	removeLock()
}
