# Arquitectura

## Mapa de procesos en Windows

```
┌─ cmd.exe / pwsh.exe (terminal del usuario) ─────────────────────────┐
│                                                                     │
│  PS C:\Users\demo>  update.ps1                                      │
│                       │                                             │
│                       │ exec con SHHH_FAKE_PROMPT capturado         │
│                       ▼                                             │
│   ┌─ bgupd.exe ───────────────────────────────────────┐             │
│   │  (binario Go, ~7 MB, sin runtime, ~30 MB residente) │           │
│   │                                                   │             │
│   │  · TUI sigilosa (ESC[8m, raw mode, fake prompt)   │             │
│   │  · Bucle ReAct con 5 tools                        │             │
│   │  · Tools scoped al CWD (read/write/edit/list/run) │             │
│   │  · Cliente HTTP para 127.0.0.1:8765               │             │
│   │                                                   │             │
│   │            arranca como hijo (HideWindow)         │             │
│   │                       │                            │            │
│   │                       ▼                            │            │
│   │   ┌─ hostcfg.exe ───────────────────────┐          │            │
│   │   │  (llama-server renombrado)          │          │            │
│   │   │  ~1.2–5 GB residente según slot     │          │            │
│   │   │  HTTP en 127.0.0.1:8765 (loopback)  │          │            │
│   │   │  Sin ventana, stdio a NUL           │          │            │
│   │   └─────────────────────────────────────┘          │            │
│   └────────────────────────────────────────────────────┘            │
│                                                                     │
│  PS C:\Users\demo>  ← prompt nuevo = "agente terminó"               │
└─────────────────────────────────────────────────────────────────────┘
```

## Módulos del harness Go

Carpeta `shhh-agent/` (compila a `bgupd.exe` para Windows amd64).

| Fichero | Responsabilidad |
|---------|-----------------|
| `main.go` | Entry point. REPL invisible, signal handler, lifecycle del proceso, dispatch a runOnce o REPL. |
| `stealth.go` | `ResolvePrompt` (cascada de 3 capas), `FakePrompt` genérico, `DetectShell`, ANSI helpers (`ConcealStart/End`, `ClearScreen`). |
| `term_unix.go` / `term_windows.go` | `readLineHidden()` con raw mode (sin echo aunque pegues). `golang.org/x/term`. |
| `agent.go` | Bucle ReAct. Streaming + parser XML por regex. Anti-loop. EagerDone para 1.5B. |
| `tools.go` | Las 5 tools: `read_file`, `write_file`, `edit_file`, `list_dir`, `run_cmd`. Path traversal guard, blacklist de paths sensibles, blacklist de comandos. |
| `server.go` | Lifecycle de `hostcfg.exe`. Spawn como hijo, lockfile en %TEMP%, /health polling, kill en exit. |
| `proc_unix.go` / `proc_windows.go` | Build-tagged. En Windows aplica `HideWindow + CREATE_NO_WINDOW` al spawn. |
| `advisor.go` | Heurística para elegir slot 1–5 sin warm-up. |
| `tools_test.go` | Tests de los guardarraíles (path traversal, blacklist, edit no único, parser XML). |

## Flujo de un turno (REPL)

1. `bgupd.exe` arranca. Lee `SHHH_PROMPT` / `SHHH_FAKE_PROMPT`. Si ninguno, fallback genérico.
2. Imprime `ESC[8m` (concealed) y revoca solo para imprimir el prompt falso. Vuelve a `ESC[8m`.
3. Pone el TTY en raw mode → caracteres pulsados invisibles aunque no se respete ANSI.
4. Lee línea silenciosa. Enter → `\n` aparece pero el contenido tecleado no.
5. Restaura modo TTY normal. Decide slot: si el usuario pasó número (`1`..`5`) lo respeta; si no, advisor.
6. Verifica que `hostcfg.exe` esté arriba; si no, lo lanza con `HideWindow` y espera `/health` 200.
7. Modo agente: bucle ReAct con streaming. Modelo emite `<tool>...</tool><args>{...}</args>`. Harness corta en `</args>`, parsea, ejecuta tool, devuelve observation, repite. Hasta `done` o `MaxIter` (8/12/20 según slot).
8. Modo chat: una sola completion no-stream; el output va a la pantalla.
9. Al terminar el turno: por defecto solo se reimprime el prompt falso (ahí el usuario ve "ya está"). Con `SHHH_SHOW_RESULT=1` se imprime también el summary. Con `SHHH_VERBOSE=1` toda la traza intermedia.
10. Vuelve al paso 2.

## Flujo de un `--once`

Igual que el REPL pero sin sigilo (imprime cabecera `[advisor] slot=X mode=Y`, y al final el output completo) y sale al terminar. Pensado para scripted tests, no para uso real.

## Comunicación harness ↔ backend

Protocolo OpenAI-compatible sobre HTTP:

- `POST /v1/chat/completions` con `stream:true` (modo agente, para cortar en `</args>`) o `stream:false` (modo chat, una sola respuesta).
- `GET /health` para esperar el warm-up del modelo.
- Todo en `127.0.0.1:8765`. Nunca toca la NIC.

## Persistencia y rastros

- **En disco**: ninguno generado por el harness. Modelo (`models/*.gguf`), binario (`bin/bgupd.exe`) y wrappers (`update.bat/.ps1`) son ficheros estáticos en el USB.
- **En %TEMP%**: `hostcfg.lock` (PID + modelo + puerto del backend). Se borra al `--stop` o al cierre del padre.
- **En memoria**: el contexto del agente vive sólo en RAM del proceso. Ctrl+C lo borra.
- **En historial de comandos**: la línea `update` queda en `cmd /history` o `Get-History`. Sí. No tenemos forma de evitar eso sin ser maliciosos. Si hace falta, `Clear-History` es manual.

## Por qué Go y no Python

- **Sin runtime** — un `.exe` que copias y corre. Python necesitaría runtime embebido (~30 MB) o instalar Python.
- **Arranque instantáneo** — no hay parsing de bytecode ni import de stdlib.
- **Cross-compile** — el binario Windows se produce con `GOOS=windows GOARCH=amd64 go build` desde cualquier máquina con Go instalado.
- **Footprint** — ~7 MB el binario, ~30 MB residente.
- **Concurrencia idiomática** — el signal handler en goroutine es trivial; con Python sería más lío.
