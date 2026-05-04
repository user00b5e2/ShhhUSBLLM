package main

import (
	"strings"
	"time"
)

// ModelSlot is the user-facing 1..N selection.
type ModelSlot int

const (
	SlotAgentFast    ModelSlot = 1 // Qwen3-1.7B Q4 — agent default
	SlotAgentPrecise ModelSlot = 2 // Qwen3-4B-Instruct-2507 Q4
	SlotChatCode     ModelSlot = 3 // Qwen3-4B-Instruct-2507 Q4 — chat mode
	SlotChatFallback ModelSlot = 4 // Qwen3-1.7B Q4 — chat fallback
	SlotAgentLarge   ModelSlot = 5 // Qwen3-8B Q4 — multi-file tasks
)

// ModelInfo describes a slot.
type ModelInfo struct {
	Slot          ModelSlot
	File          string        // basename inside models/
	Mode          Mode          // chat or agent
	HumanTag      string
	CtxSize       int           // tokens; 0 = use defaultCtx
	MaxIter       int           // agent loop iterations; 0 = default 8
	EagerDone     bool          // 1.5B-style hack: end turn after first successful mutation
	TurnTimeout   time.Duration // 0 = default 10m
	Qwen3DualMode bool          // true ⇒ inject "/no_think" directive in system prompt
}

// Mode of operation for a slot.
type Mode int

const (
	ModeAgent Mode = iota
	ModeChat
)

func ModelTable() map[ModelSlot]ModelInfo {
	return map[ModelSlot]ModelInfo{
		SlotAgentFast: {
			Slot: SlotAgentFast, File: "qwen3-1.7b-q4_k_m.gguf",
			Mode: ModeAgent, HumanTag: "1.7B agent",
			CtxSize: 4096, MaxIter: 8, EagerDone: true,
			TurnTimeout:   5 * time.Minute,
			Qwen3DualMode: true,
		},
		SlotAgentPrecise: {
			Slot: SlotAgentPrecise, File: "qwen3-4b-instruct-2507-q4_k_m.gguf",
			Mode: ModeAgent, HumanTag: "4B agent",
			CtxSize: 8192, MaxIter: 12, EagerDone: false,
			TurnTimeout: 10 * time.Minute,
			// Qwen3-Instruct-2507 is non-thinking native; no flag needed.
		},
		SlotChatCode: {
			Slot: SlotChatCode, File: "qwen3-4b-instruct-2507-q4_k_m.gguf",
			Mode: ModeChat, HumanTag: "4B chat",
			CtxSize: 8192, TurnTimeout: 5 * time.Minute,
		},
		SlotChatFallback: {
			Slot: SlotChatFallback, File: "qwen3-1.7b-q4_k_m.gguf",
			Mode: ModeChat, HumanTag: "1.7B chat",
			CtxSize: 4096, TurnTimeout: 3 * time.Minute,
			Qwen3DualMode: true,
		},
		SlotAgentLarge: {
			Slot: SlotAgentLarge, File: "qwen3-8b-q4_k_m.gguf",
			Mode: ModeAgent, HumanTag: "8B agent (long)",
			CtxSize: 16384, MaxIter: 20, EagerDone: false,
			TurnTimeout:   45 * time.Minute,
			Qwen3DualMode: true,
		},
	}
}

// AdviseSlot picks a slot from the first user request when no -N was given.
// Pure heuristic — no model call — to avoid an extra warm-up.
func AdviseSlot(req string) ModelSlot {
	r := strings.ToLower(req)

	// Agent triggers (action verbs) → agent.
	agentVerbs := []string{
		"edit", "edita", "modify", "modifica", "fix", "arregla",
		"create", "crea", "add", "añade", "anade", "remove", "elimina", "borra", "delete",
		"refactor", "refactoriza", "rename", "renombra", "replace", "reemplaza",
		"run", "ejecuta", "test", "build", "compila", "install", "instala",
		"write", "escribe", "implement", "implementa",
	}
	hasVerb := false
	for _, v := range agentVerbs {
		if strings.Contains(r, v) {
			hasVerb = true
			break
		}
	}

	// Long task signals: multi-file projects, spec-driven, compile loops.
	// We deliberately avoid generic ".md" markers — creating a single .md is NOT a big task.
	// Triggers must indicate either reading a spec or producing many files.
	largeTaskMarkers := []string{
		"según spec", "segun spec", "from spec.md", "de spec.md", "de specs.md", "lee spec", "read spec",
		"5 cpp", "5 ficheros", "five files", "varios ficheros cpp", "varios .cpp",
		"compila", "compile cpp", "compile c++", "g++", "clang++", "cl /c",
		"test suite", "all of them",
	}
	largeTask := false
	for _, m := range largeTaskMarkers {
		if strings.Contains(r, m) {
			largeTask = true
			break
		}
	}

	multiFile := strings.Contains(r, "all ") || strings.Contains(r, "todos") ||
		strings.Contains(r, "every ") || strings.Contains(r, "across") ||
		strings.Contains(r, "varios") || strings.Contains(r, "multiple") ||
		largeTask

	chatMarkers := []string{
		"explain", "explica", "what does", "qué hace", "que hace",
		"how does", "cómo funciona", "como funciona",
		"why ", "por qué", "porque ",
		"summarize", "resume", "describe",
	}
	for _, m := range chatMarkers {
		if strings.Contains(r, m) {
			return SlotChatCode
		}
	}

	if hasVerb {
		if largeTask {
			return SlotAgentLarge
		}
		if multiFile {
			return SlotAgentPrecise
		}
		return SlotAgentFast
	}
	// Default fallback for ambiguous prompts: chat (safer, no edits).
	return SlotChatCode
}
