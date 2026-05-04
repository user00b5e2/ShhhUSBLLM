package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
)

const (
	defaultHost = "127.0.0.1"
	defaultPort = 8765
	defaultCtx  = 4096
)

func main() {
	var (
		stopFlag    = flag.Bool("stop", false, "kill background hostcfg and exit")
		unsafeFlag  = flag.Bool("unsafe", false, "disable CWD scoping (DANGEROUS)")
		modelsDir   = flag.String("models", "", "override models directory")
		binDir      = flag.String("bin", "", "override bin directory")
		port        = flag.Int("port", defaultPort, "llama-server port")
		host        = flag.String("host", defaultHost, "llama-server host")
		verbose     = flag.Bool("verbose", os.Getenv("SHHH_VERBOSE") == "1", "show intermediate agent steps")
		once        = flag.String("once", "", "run a single prompt non-interactively (no stealth, prints output)")
	)
	flag.Parse()

	if *stopFlag {
		StopServer()
		return
	}

	defer StopServer()

	if *once != "" {
		runOnce(*once, *unsafeFlag, *modelsDir, *binDir, *host, *port, *verbose)
		return
	}

	// Slot may be passed positionally (mimics `shhh N`).
	var requestedSlot ModelSlot
	if args := flag.Args(); len(args) > 0 {
		if n, err := strconv.Atoi(args[0]); err == nil && n >= 1 && n <= 5 {
			requestedSlot = ModelSlot(n)
		}
	}

	exeDir := executableDir()
	if *modelsDir == "" {
		*modelsDir = filepath.Join(exeDir, "..", "models")
	}
	if *binDir == "" {
		*binDir = exeDir
	}

	// Hide everything from now on (including any error printing during warm-up).
	ConcealStart(os.Stdout)
	defer ConcealEnd(os.Stdout)

	shell := DetectShell()
	prompt := ResolvePrompt(shell)
	showResult := os.Getenv("SHHH_SHOW_RESULT") == "1" || *verbose

	// Print fake prompt immediately so the terminal looks idle while warm-up runs.
	ConcealEnd(os.Stdout)
	fmt.Print(prompt)
	ConcealStart(os.Stdout)

	// Panic button: Ctrl+C clears screen and re-prints fake prompt before exit.
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)
	go func() {
		<-sigCh
		ConcealEnd(os.Stdout)
		ClearScreen(os.Stdout)
		fmt.Print(prompt)
		StopServer()
		os.Exit(0)
	}()

	tools, err := NewTools(*unsafeFlag)
	if err != nil {
		bail(err)
	}

	for {
		// Read input invisibly. The fake prompt is already on screen.
		req, err := readLineHidden()
		if err != nil {
			if err == errInterrupted {
				ConcealEnd(os.Stdout)
				ClearScreen(os.Stdout)
				fmt.Print(prompt)
				return
			}
			bail(err)
		}
		req = strings.TrimSpace(req)
		if req == "" {
			// Empty enter: re-print prompt and continue, mimicking shells.
			ConcealEnd(os.Stdout)
			fmt.Print(prompt)
			ConcealStart(os.Stdout)
			continue
		}
		if req == "exit" || req == "quit" {
			ConcealEnd(os.Stdout)
			ClearScreen(os.Stdout)
			fmt.Print(prompt)
			return
		}

		// Decide the slot.
		slot := requestedSlot
		if slot == 0 {
			slot = AdviseSlot(req)
		}
		info := ModelTable()[slot]

		// Ensure backend.
		modelPath := filepath.Join(*modelsDir, info.File)
		if _, err := os.Stat(modelPath); err != nil {
			reportInfraError(prompt, fmt.Errorf("missing model: %s", modelPath), showResult)
			continue
		}
		ctxSize := info.CtxSize
		if ctxSize == 0 {
			ctxSize = defaultCtx
		}
		cfg := ServerConfig{
			BinaryPath: serverBinaryPath(*binDir),
			ModelPath:  modelPath,
			Host:       *host,
			Port:       *port,
			CtxSize:    ctxSize,
			Threads:    threadCount(),
		}
		if err := EnsureServer(cfg); err != nil {
			reportInfraError(prompt, fmt.Errorf("server: %w", err), showResult)
			continue
		}
		if err := WaitHealthy(*host, *port, 120*time.Second); err != nil {
			reportInfraError(prompt, err, showResult)
			continue
		}

		client := NewClient(*host, *port)
		turnTimeout := info.TurnTimeout
		if turnTimeout == 0 {
			turnTimeout = 10 * time.Minute
		}
		ctx, cancel := context.WithTimeout(context.Background(), turnTimeout)

		var output string
		if info.Mode == ModeAgent {
			ag := &Agent{
				Cli: client, Tools: tools, Out: os.Stdout,
				Verbose: *verbose, MaxIter: info.MaxIter,
				EagerDone: info.EagerDone,
			}
			output, err = ag.Run(ctx, req)
		} else {
			output, err = client.Complete(ctx, []Message{
				{Role: "system", Content: chatSystemPrompt},
				{Role: "user", Content: req},
			}, nil)
		}
		cancel()

		ConcealEnd(os.Stdout)
		// Stealth: by default the user only sees the next prompt. The result of the
		// turn lives in the filesystem; if they want a summary they set SHHH_SHOW_RESULT=1.
		// Errors are reduced to a single discreet '!' so the user knows something went wrong.
		if err != nil {
			fmt.Println("!")
		} else if showResult {
			s := strings.TrimSpace(output)
			if s != "" {
				fmt.Println()
				fmt.Println(s)
			}
		}
		fmt.Print(prompt)
		ConcealStart(os.Stdout)
	}
}

const chatSystemPrompt = `You are a concise coding assistant. Reply in plain text, no markdown fences. Maximum 8 lines unless the user asks for more.`

// runOnce executes a single prompt without any stealth or REPL. Useful for
// scripted tests and first-time validation when typing blind would be awkward.
func runOnce(req string, unsafe bool, modelsDir, binDir, host string, port int, verbose bool) {
	if modelsDir == "" {
		modelsDir = filepath.Join(executableDir(), "..", "models")
	}
	if binDir == "" {
		binDir = executableDir()
	}
	tools, err := NewTools(unsafe)
	if err != nil {
		bailPlain(err)
	}
	slot := AdviseSlot(req)
	info := ModelTable()[slot]
	fmt.Printf("[advisor] slot=%d (%s) mode=%v\n", info.Slot, info.HumanTag, info.Mode)

	modelPath := filepath.Join(modelsDir, info.File)
	if _, err := os.Stat(modelPath); err != nil {
		bailPlain(fmt.Errorf("missing model: %s", modelPath))
	}
	cfg := ServerConfig{
		BinaryPath: serverBinaryPath(binDir),
		ModelPath:  modelPath,
		Host:       host, Port: port, CtxSize: defaultCtx, Threads: threadCount(),
	}
	if err := EnsureServer(cfg); err != nil {
		bailPlain(fmt.Errorf("server: %w", err))
	}
	if err := WaitHealthy(host, port, 120*time.Second); err != nil {
		bailPlain(err)
	}
	client := NewClient(host, port)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()
	var out string
	if info.Mode == ModeAgent {
		ag := &Agent{
			Cli: client, Tools: tools, Out: os.Stdout,
			Verbose: verbose, MaxIter: info.MaxIter,
			EagerDone: info.EagerDone,
		}
		out, err = ag.Run(ctx, req)
	} else {
		out, err = client.Complete(ctx, []Message{
			{Role: "system", Content: chatSystemPrompt},
			{Role: "user", Content: req},
		}, nil)
	}
	if err != nil {
		bailPlain(err)
	}
	fmt.Println(strings.TrimSpace(out))
}

func bailPlain(err error) {
	fmt.Fprintf(os.Stderr, "%v\n", err)
	os.Exit(1)
}

// reportInfraError surfaces backend/model errors in stealth mode.
// In verbose / show-result mode the full message is printed; otherwise just '!'.
func reportInfraError(prompt string, err error, showResult bool) {
	ConcealEnd(os.Stdout)
	if showResult {
		fmt.Printf("\n[%v]\n", err)
	} else {
		fmt.Println("!")
	}
	fmt.Print(prompt)
	ConcealStart(os.Stdout)
}

func bail(err error) {
	ConcealEnd(os.Stdout)
	fmt.Fprintf(os.Stderr, "\n%v\n", err)
	os.Exit(1)
}

func executableDir() string {
	exe, err := os.Executable()
	if err != nil {
		return "."
	}
	resolved, err := filepath.EvalSymlinks(exe)
	if err == nil {
		exe = resolved
	}
	return filepath.Dir(exe)
}

func serverBinaryPath(binDir string) string {
	name := "hostcfg"
	if runtime.GOOS == "windows" {
		name += ".exe"
	}
	return filepath.Join(binDir, name)
}

func threadCount() int {
	n := runtime.NumCPU() - 1
	if n < 2 {
		n = 2
	}
	if n > 8 {
		n = 8
	}
	return n
}
