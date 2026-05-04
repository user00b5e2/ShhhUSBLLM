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

---

# Iteración 3 — Modelos Qwen3 (mayo 2026)

## Resumen ejecutivo

Reemplazo de la familia de modelos en todos los slots: **Qwen2.5-Coder → Qwen3 / Qwen3-Instruct-2507**. Mejora estimada de calidad: una "talla" — el 1.7B se aproxima al 3B anterior, el 4B al 7B anterior, el 8B al 14B anterior.

| Área | Iter 2 | Iter 3 |
|------|--------|--------|
| Slot 1 | Qwen2.5-Coder-1.5B Q4 | **Qwen3-1.7B Q4** + `/no_think` |
| Slot 2 | Qwen2.5-Coder-3B Q4 | **Qwen3-4B-Instruct-2507 Q4** (non-thinking nativo) |
| Slot 3 | Qwen2.5-Coder-3B Q4 | **Qwen3-4B-Instruct-2507 Q4** |
| Slot 4 | Qwen2.5-1.5B-Instruct Q4 | **Qwen3-1.7B Q4** + `/no_think` |
| Slot 5 | Qwen2.5-Coder-7B Q4 | **Qwen3-8B Q4** + `/no_think` |
| GGUFs distintos en disco | 4 | **3** (slots 1+4 comparten Qwen3-1.7B; 2+3 comparten 4B-Instruct-2507) |
| `Qwen3DualMode` flag | n/a | Nuevo en `ModelInfo`; inyecta `/no_think` al system prompt |

## Por qué este cambio

A fecha mayo 2026, los benchmarks oficiales de Alibaba (BigCodeBench, EvalPlus, LiveCodeBench) y la comunidad muestran que Qwen3 dense pequeño bate a Qwen2.5-Coder a paridad de tamaño:

- Qwen3-1.7B base ≈ Qwen2.5-Coder-3B en code benchmarks.
- Qwen3-4B base ≈ Qwen2.5-Coder-7B.
- Qwen3-8B base ≈ Qwen2.5-Coder-14B.

A igual RAM, ganamos calidad. El coste es ~5–25 % más de latencia por token según el slot, parcialmente compensado por menos rondas de corrección.

## Decisiones técnicas tomadas

- **Variantes Instruct-2507** (julio 2025) para los slots donde existen (4B). Esa rama está fine-tuned **solo en non-thinking**, perfecto para nuestro tool-calling estricto con XML.
- **Qwen3 base + `/no_think` runtime** para los slots donde no hay 2507 (1.7B, 8B). El harness inyecta `/no_think\n\n` al inicio del system prompt cuando `Qwen3DualMode=true`.
- **No subimos a Qwen3-Coder-Next**: es MoE 80B / 3B activos — necesita ~32 GB para cargar pesos, fuera del target 8 GB.
- **Repo elegido**: `unsloth/Qwen3-*-GGUF` por estabilidad y disponibilidad de quants. Fallback no documentado: `Qwen/Qwen3-*-GGUF` oficial.

## Cambios en el código

| Fichero | Cambio |
|---------|--------|
| `shhh-agent/advisor.go` | Nuevo campo `Qwen3DualMode bool` en `ModelInfo`; tabla de modelos rehecha. |
| `shhh-agent/agent.go` | `Agent` acepta `Qwen3DualMode`; prepend `/no_think\n\n` al system prompt si está activo. **Nuevo helper `extractFirstJSONObject`** para tolerar sufijos non-JSON en `<args>` (algunos quants de Qwen3 emiten `>` extra). |
| `shhh-agent/main.go` | Propaga `Qwen3DualMode` al `Agent` y al chat path. |
| `download-models.sh` | URLs nuevas (3 GGUFs en lugar de 4); slot 5 sigue opt-in con `WITH_LARGE=1`. |
| `docs/04-MODELOS-Y-RAM.md` | Tabla rehecha; explicación de `/no_think` y Qwen3 vs Qwen2.5-Coder. |
| `docs/05-DURACIONES.md` | Latencias actualizadas (1.7B ~5 % más lento, 4B ~25 %, 8B ~15 %). |
| `README.md`, `QUICKSTART.md` | Tablas de slots con los nombres nuevos. |

## Validación realizada

Smoke test E2E en macOS M4 con Qwen3-1.7B Q4_K_M:

- Spec: MD pidiendo crear `hello.cpp` con bucle `cout` N veces, N por stdin.
- Resultado: el agente leyó el MD, generó código C++ correcto en 1 turno.
- `g++ hello.cpp -o hello && echo 3 | ./hello` → `holaholahola`. ✅

Detalle observado durante la validación: Qwen3-1.7B emite ocasionalmente caracteres extra después del JSON de `<args>` (visto: `}>` en lugar de `}`). El parser ahora extrae el primer objeto JSON balanceado y descarta el sufijo, igual que un parser XML laxo.

## Lo que NO se hizo

- **No se actualizó la tabla del README "Descargar modelos" del modo chat clásico** (líneas que mencionan Qwen1.5-4B, Phi-4-mini, Llama 3.1, etc.). Esa sección del proyecto base está atada a `shhh.bat`/`shhh.ps1`/`shhhps.bat` con URLs hardcoded; tocarla requiere reescribir esos `.bat` y se pospuso.
- **No se eliminó el modelo Qwen2.5-Coder-1.5B viejo del disco automáticamente** — lo borré manualmente al validar la iter3. Si lo tienes descargado, puedes borrarlo: ya no se referencia.
- **No se validó en Windows real**, sigue siendo válida la nota de "no probado en Windows" de `08-LIMITACIONES.md`.
