# Cambios de la iteración 2

Esta iteración refina la discreción visual del harness, añade un slot de modelo más grande (7B) para tareas multi-fichero, y reorganiza la documentación.

## Resumen ejecutivo

| Área | Iteración 1 | Iteración 2 |
|------|-------------|-------------|
| Binario | `shhh-agent.exe` | `bgupd.exe` (nombre neutro tipo utilidad de sistema) |
| Wrappers | `shhh-agent.bat`, `.ps1` | `update.bat`, `update.ps1` |
| Prompt | Genérico (`PS C:\path>`) | Captura literal del prompt real (`update.ps1` ejecuta `& { prompt }` y exporta `SHHH_FAKE_PROMPT`) |
| Sigilo de output | Imprime el output del agente al final | **Por defecto: nada**. Solo el prompt. Activable con `SHHH_SHOW_RESULT=1` |
| Errores | `[error: ...]` visible | `!` discreto (texto completo solo con `SHHH_SHOW_RESULT=1`) |
| Slots | 1–4 | 1–5 (nuevo slot 5 = 7B Q4 para tareas grandes) |
| `MaxIter` agente | Fijo en 8 | Dinámico por slot: 8 (1.5B), 12 (3B), 20 (7B) |
| `EagerDone` | Forzado siempre | Sólo en slot 1 (1.5B). Slots 3B/7B siguen el protocolo |
| Tope `read_file` | 200 KB | 1 MB (para MDs largos de spec) |
| `cmdTimeout` en `run_cmd` | 60 s | 120 s (compilaciones) |
| System prompt agente | Single-file focus | Mejorado para multi-fichero + auto-corrección con `g++` |
| Heurística advisor | "edita / refactoriza" → 1.5B/3B | Detecta "spec.md", "5 cpp", "g++", "compile" → slot 5 (7B) |
| Documentación educativa | No existía | `docs/` con 9 ficheros |

## Lista de ficheros tocados o creados

### Tocados (lógica)

- `shhh-agent/main.go` — `ResolvePrompt`, modo sigilo total, `MaxIter`/`EagerDone`/`TurnTimeout` por slot, errores discretos.
- `shhh-agent/stealth.go` — `ResolvePrompt` con cascada `SHHH_PROMPT` > `SHHH_FAKE_PROMPT` > genérico.
- `shhh-agent/agent.go` — `MaxIter` configurable, `EagerDone` flag, system prompt multi-fichero.
- `shhh-agent/advisor.go` — slot 5, `ModelInfo` extendido con `CtxSize`/`MaxIter`/`EagerDone`/`TurnTimeout`, heurística refinada.
- `shhh-agent/tools.go` — `maxFileBytes` 1 MB, `cmdTimeout` 120 s.
- `build.sh` — output a `bgupd.exe`.
- `download-models.sh` — opt-in para 7B con `WITH_LARGE=1`.

### Eliminados

- `shhh-agent.bat` — reemplazado por `update.bat`.
- `shhh-agent.ps1` — reemplazado por `update.ps1`.
- `bin/shhh-agent.exe` — reemplazado por `bin/bgupd.exe`.

### Creados

- `update.bat` — wrapper CMD que captura CWD y forwardea a `bgupd.exe`.
- `update.ps1` — wrapper PowerShell que captura `prompt` real y lo exporta.
- `bin/bgupd.exe` — binario Windows con nombre nuevo.
- `docs/00-INDICE.md` … `docs/08-LIMITACIONES.md` — esta documentación.

### NO tocados (deliberadamente)

- `README.md` — sigue siendo guía de uso de la iteración 1.
- `QUICKSTART.md` — comandos de referencia primera versión.
- `shhh.bat`, `shhh.ps1`, `shhhps.bat` — proyecto original (modo chat de la versión inicial).
(no aplica — todos los artefactos del USB son ahora Windows-only)

## Decisiones que tomé sin consultarte

- **`EagerDone` solo en slot 1**, no en slot 2/5. Razón: el 1.5B no llama a `done()` correctamente y se queda creando ficheros hasta agotar iteraciones. El 3B y 7B sí lo hacen, así que les dejo seguir el protocolo (importante para multi-fichero).
- **Heurística del advisor refinada**: cualquier `.md` ya **no** dispara slot 5. Solo combinaciones explícitas como "lee spec.md", "5 cpp", "g++", "compile cpp". Antes "crea hola.md" iba a slot 5 por error.
- **Tope de `read_file` a 1 MB sin condicional**: razoné que 1 MB sigue siendo lo bastante pequeño para no explotar el contexto y útil para specs largas. No hay heurística condicional.
- **Sin auto-shutdown por idle**: el plan original lo mencionaba; no lo implementé. `defer StopServer()` cierra el backend al salir del harness, lo cual es suficiente. Si quieres timer de idle, dímelo.

## Lo que NO hice

- **`SHHH_BEEP`** — el usuario rechazó el beep audible.
- **OSC para cambiar título de ventana** — el usuario eligió "sigilo total" en su lugar.
- **Detección de compilador en runtime con prompt adaptativo** — el system prompt asume que `g++` existe (confirmado por el usuario). Si no estuviera, `run_cmd` fallaría limpiamente con "command not found", el agente lo vería y reaccionaría.
- **DeepSeek-Coder-V2-Lite como slot 6** — alternativa al 7B Q4. No la incluí para no abrumar; añadirlo es trivial si quieres.
- **Cambio del nombre de `hostcfg.exe`** a algo más neutro — el actual es suficiente.
- **Test unitario del advisor** — debería existir. No lo escribí. Pendiente.
