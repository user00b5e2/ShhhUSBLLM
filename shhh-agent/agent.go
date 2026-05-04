package main

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"
)

const defaultMaxIterations = 8

const agentSystemPrompt = `You are a code-editing agent. The user types instructions blindly; you act on the workspace.

You answer with EXACTLY ONE tool call per turn, in this XML format and nothing else:

<tool>NAME</tool><args>JSON_OBJECT</args>

Available tools:
- read_file       args: {"path":"relative/path"}
- write_file      args: {"path":"relative/path","content":"..."}
- edit_file       args: {"path":"relative/path","old":"exact text","new":"replacement"}  // 'old' must be unique
- list_dir        args: {"path":"relative/path"}
- run_cmd         args: {"cmd":"shell command"}
- done            args: {"summary":"one short sentence for the user"}

Rules:
- Output ONLY the XML. No markdown, no prose, no explanations outside <args>.
- Paths are relative to the workspace; never use ".." or absolute paths.
- For edit_file, 'old' must be a substring that appears EXACTLY ONCE in the file.
- After ANY error, examine the message and adapt; do not repeat the same call.
- Never repeat the EXACT same call twice in a row.
- For SINGLE-FILE tasks: write the file, then call 'done'.
- For MULTI-FILE tasks (e.g. "create 5 cpp files from spec.md"): read inputs first, then write each file in sequence (one tool call per file), THEN call run_cmd to compile if a build step makes sense, fix errors found, finally call 'done' with a one-line summary.
- A successful observation (anything not starting with "ERROR") means: continue with the next required step or finish.
- If 'run_cmd' shows compilation errors, read the relevant files, edit to fix, and recompile. Iterate until clean or until you've tried 3 times for the same file.`

// Message — OpenAI-compatible.
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type chatReq struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	Stream      bool      `json:"stream"`
	Temperature float64   `json:"temperature"`
	MaxTokens   int       `json:"max_tokens"`
	Stop        []string  `json:"stop,omitempty"`
}

type chatChoice struct {
	Message Message `json:"message"`
}

type chatResp struct {
	Choices []chatChoice `json:"choices"`
}

type streamDelta struct {
	Choices []struct {
		Delta struct {
			Content string `json:"content"`
		} `json:"delta"`
	} `json:"choices"`
}

// Client speaks to the local llama-server.
type Client struct {
	BaseURL string
	HTTP    *http.Client
}

func NewClient(host string, port int) *Client {
	return &Client{
		BaseURL: fmt.Sprintf("http://%s:%d", host, port),
		HTTP:    &http.Client{Timeout: 5 * time.Minute},
	}
}

// Complete runs a non-streaming completion.
func (c *Client) Complete(ctx context.Context, msgs []Message, stop []string) (string, error) {
	body, _ := json.Marshal(chatReq{
		Model: "local", Messages: msgs, Stream: false,
		Temperature: 0.2, MaxTokens: 1024, Stop: stop,
	})
	req, _ := http.NewRequestWithContext(ctx, "POST", c.BaseURL+"/v1/chat/completions", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.HTTP.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		b, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("llama-server %d: %s", resp.StatusCode, string(b))
	}
	var cr chatResp
	if err := json.NewDecoder(resp.Body).Decode(&cr); err != nil {
		return "", err
	}
	if len(cr.Choices) == 0 {
		return "", fmt.Errorf("no choices")
	}
	return cr.Choices[0].Message.Content, nil
}

// Stream streams content; the callback can return true to stop early.
func (c *Client) Stream(ctx context.Context, msgs []Message, stop []string, onDelta func(string) bool) (string, error) {
	body, _ := json.Marshal(chatReq{
		Model: "local", Messages: msgs, Stream: true,
		Temperature: 0.2, MaxTokens: 1024, Stop: stop,
	})
	req, _ := http.NewRequestWithContext(ctx, "POST", c.BaseURL+"/v1/chat/completions", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.HTTP.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		b, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("llama-server %d: %s", resp.StatusCode, string(b))
	}

	var full strings.Builder
	sc := bufio.NewScanner(resp.Body)
	sc.Buffer(make([]byte, 0, 64*1024), 4*1024*1024)
	for sc.Scan() {
		line := sc.Text()
		if !strings.HasPrefix(line, "data: ") {
			continue
		}
		payload := strings.TrimPrefix(line, "data: ")
		if payload == "[DONE]" {
			break
		}
		var d streamDelta
		if err := json.Unmarshal([]byte(payload), &d); err != nil {
			continue
		}
		for _, ch := range d.Choices {
			if ch.Delta.Content != "" {
				full.WriteString(ch.Delta.Content)
				if onDelta != nil && onDelta(full.String()) {
					return full.String(), nil
				}
			}
		}
	}
	return full.String(), nil
}

// --- ReAct loop ---

var toolRe = regexp.MustCompile(`(?s)<tool>\s*([a-z_]+)\s*</tool>\s*<args>(.*?)</args>`)

type parsedCall struct {
	Name string
	Args json.RawMessage
}

func parseTool(text string) (*parsedCall, error) {
	m := toolRe.FindStringSubmatch(text)
	if m == nil {
		return nil, fmt.Errorf("no <tool>...</tool><args>...</args> found")
	}
	args := strings.TrimSpace(m[2])
	// Tolerate a non-JSON suffix after the closing `}` (some Qwen3 quants emit
	// stray characters like ">" after the args object). Extract the first
	// balanced JSON object and drop anything after it.
	args = extractFirstJSONObject(args)
	return &parsedCall{Name: m[1], Args: json.RawMessage(args)}, nil
}

// extractFirstJSONObject returns the substring from the first `{` to its
// balanced matching `}`, ignoring `{`/`}` inside JSON strings. If parsing
// fails, the original input is returned unchanged.
func extractFirstJSONObject(s string) string {
	start := strings.IndexByte(s, '{')
	if start < 0 {
		return s
	}
	depth := 0
	inStr := false
	escape := false
	for i := start; i < len(s); i++ {
		c := s[i]
		if escape {
			escape = false
			continue
		}
		if inStr {
			switch c {
			case '\\':
				escape = true
			case '"':
				inStr = false
			}
			continue
		}
		switch c {
		case '"':
			inStr = true
		case '{':
			depth++
		case '}':
			depth--
			if depth == 0 {
				return s[start : i+1]
			}
		}
	}
	return s
}

// Agent runs the ReAct loop until 'done' or MaxIter.
type Agent struct {
	Cli       *Client
	Tools     *Tools
	Out       io.Writer // where to print intermediate steps if verbose
	Verbose   bool
	MaxIter   int  // 0 → defaultMaxIterations
	EagerDone bool // tiny-model hack: end turn after first successful mutation
}

// Run executes one user request, blocking until 'done' is called or iterations exhaust.
// Returns the summary string.
func (a *Agent) Run(ctx context.Context, userReq string) (string, error) {
	msgs := []Message{
		{Role: "system", Content: agentSystemPrompt},
		{Role: "user", Content: userReq},
	}
	stop := []string{"</args>"}
	var lastSig string
	repeats := 0

	maxIter := a.MaxIter
	if maxIter <= 0 {
		maxIter = defaultMaxIterations
	}

	for i := 0; i < maxIter; i++ {
		// Stream up to the first </args>; cut early to save tokens.
		raw, err := a.Cli.Stream(ctx, msgs, stop, func(buf string) bool {
			return strings.Contains(buf, "</args>")
		})
		if err != nil {
			return "", err
		}
		// llama-server's stop sequence is consumed before emit; re-add for parser.
		text := raw
		if !strings.Contains(text, "</args>") {
			text += "</args>"
		}

		call, err := parseTool(text)
		if err != nil {
			// One automatic correction round.
			msgs = append(msgs,
				Message{Role: "assistant", Content: text},
				Message{Role: "user", Content: "Output ONLY one <tool>NAME</tool><args>{...}</args> call. No prose."},
			)
			continue
		}

		if a.Verbose {
			fmt.Fprintf(a.Out, "[%s] %s\n", call.Name, string(call.Args))
		}

		if call.Name == "done" {
			var d struct{ Summary string `json:"summary"` }
			_ = json.Unmarshal(call.Args, &d)
			return d.Summary, nil
		}

		// Detect exact-repeat loops (small models often forget to call done).
		sig := call.Name + "::" + string(call.Args)
		if sig == lastSig {
			repeats++
			if repeats >= 1 {
				return "task likely complete (auto-done after repeat)", nil
			}
		} else {
			repeats = 0
		}
		lastSig = sig

		result, toolErr := a.dispatch(call)

		// Tiny models reliably fail to call done(). When EagerDone is set,
		// one successful mutation = task done. Larger models (3B/7B) follow
		// the protocol and can chain multiple write/edit calls before done.
		if a.EagerDone && toolErr == nil && (call.Name == "write_file" || call.Name == "edit_file") {
			return result, nil
		}
		obs := result
		if toolErr != nil {
			obs = "ERROR: " + toolErr.Error()
			if result != "" {
				obs += "\n" + result
			}
		}
		// Cap observation size to avoid context blowup.
		if len(obs) > 8000 {
			obs = obs[:8000] + "\n... (truncated)"
		}

		msgs = append(msgs,
			Message{Role: "assistant", Content: text},
			Message{Role: "user", Content: "<observation>" + obs + "</observation>"},
		)
	}
	return "max iterations reached without done()", nil
}

func (a *Agent) dispatch(call *parsedCall) (string, error) {
	switch call.Name {
	case "read_file":
		return a.Tools.ReadFile(call.Args)
	case "write_file":
		return a.Tools.WriteFile(call.Args)
	case "edit_file":
		return a.Tools.EditFile(call.Args)
	case "list_dir":
		return a.Tools.ListDir(call.Args)
	case "run_cmd":
		return a.Tools.RunCmd(call.Args)
	default:
		return "", fmt.Errorf("unknown tool: %s", call.Name)
	}
}
