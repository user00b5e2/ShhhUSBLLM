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

# Iteración 3 — Refinamiento del default y reverso a Qwen2.5-Coder (mayo 2026)

## Resumen ejecutivo

Dos cambios:

1. **Default del advisor cambia**: ahora cuando `update` (sin slot) detecta una tarea de coding, abre el **slot 2 (3B-Coder)** por defecto, no el slot 1 (1.5B). El 1.5B se reserva para tareas explícitamente triviales (typos, one-liners, renames). Esto pone el sweet-spot calidad/velocidad como default para programar.
2. **Mantenemos Qwen2.5-Coder** en todos los slots tras evaluar la alternativa Qwen3 dense. Razón: tool-calling con XML estricto es más fiable en una familia coder-específica fine-tuned sobre código.

## Tabla final de modelos (idéntica a iter 2 pero con default ajustado)

| Slot | Modelo | RAM | Modo | Default del advisor |
|------|--------|-----|------|---------------------|
| 1 | Qwen2.5-Coder-1.5B Q4 | 1.3 GB | agente | Solo si la petición tiene "typo / one-liner / trivial" |
| 2 | Qwen2.5-Coder-3B Q4 | 2.4 GB | agente | **Sí** — default cuando hay verbo de acción |
| 3 | Qwen2.5-Coder-3B Q4 | 2.4 GB | chat | Sí — para "explica / cómo funciona" |
| 4 | Qwen2.5-Coder-1.5B Q4 | 1.3 GB | chat | No (solo manual) |
| 5 | Qwen2.5-Coder-7B Q4 | 5.2 GB | agente | Sí — para "spec.md / 5 cpp / compile" |

Beneficio: **3 GGUFs distintos** en disco para 5 slots (slot 4 ahora reusa el 1.5B-Coder de slot 1; antes era 1.5B-Instruct generalista).

## Por qué se descartó Qwen3 dense

Probé Qwen3-1.7B / Qwen3-4B-Instruct-2507 / Qwen3-8B en una rama intermedia de iter3. Los datos:

- **Qwen3-1.7B emitió `}>` después del JSON** en `<args>`, requiriendo un parser tolerante (`extractFirstJSONObject`). Qwen2.5-Coder-1.5B no tiene ese problema.
- **Qwen3 dual-mode requería `/no_think` runtime** para evitar tokens `<think>` rompiendo el parser. Complicación innecesaria.
- Los benchmarks oficiales de Alibaba decían "Qwen3-1.7B base ≈ Qwen2.5-Coder-3B" pero esos eran HumanEval/MBPP, **no tool-calling estricto bajo XML**, que es nuestro caso real.
- Qwen2.5-Coder-7B sigue siendo la opción coder-específica más potente que cabe en 8 GB; Qwen3-Coder no tiene variantes pequeñas dense (sólo MoE 80B+, fuera de target).

## Cambios en el código respecto a iter 2

| Fichero | Cambio |
|---------|--------|
| `shhh-agent/advisor.go` | Default heurístico de "agente" cambió de slot 1 a **slot 2**. Slot 4 ahora usa `qwen2.5-coder-1.5b-instruct-q4_k_m.gguf` (mismo fichero que slot 1) en lugar de `qwen2.5-1.5b-instruct-q4_k_m.gguf` (no-coder). Nuevos marcadores `trivialMarkers` para identificar peticiones que sí justifican slot 1. |
| `shhh-agent/agent.go` | Mantengo `extractFirstJSONObject` como parser tolerante general aunque no es necesario con Qwen2.5-Coder — robustez para futuros modelos. |
| `download-models.sh` | URLs apuntan a la familia Qwen2.5-Coder (oficial `Qwen/Qwen2.5-Coder-*-Instruct-GGUF`). |
| `docs/04-MODELOS-Y-RAM.md` | Tabla rehecha + explicación de la decisión Qwen2.5-Coder vs Qwen3 dense. |
| `docs/05-DURACIONES.md` | Latencias revertidas a las de iter 2. |

## Lo que NO se cambió

- `shhh.bat`, `shhh.ps1`, `shhhps.bat` — modo chat clásico del proyecto base, sigue intacto.
- Tabla del README "Descargar modelos" (sección modo chat clásico).
- Lógica de stealth, tools, server lifecycle — independiente del modelo.
- `extractFirstJSONObject` (parser tolerante de iter3-Qwen3) — se mantiene como defensa adicional aunque ya no es estrictamente necesario.

## Validación realizada

Smoke test E2E en macOS con Qwen2.5-Coder-1.5B (los slots 2 y 5 requieren descarga; lógica idéntica):

- Spec: MD pidiendo crear `hello.cpp` con bucle `cout` N veces, N por stdin.
- Resultado: el agente leyó el MD, generó código C++ correcto en 1 turno.
- `g++ hello.cpp -o hello && echo 3 | ./hello` → `holaholahola`. ✅

Validación pendiente en Windows real (sigue válida la nota de "no probado en Windows" de `08-LIMITACIONES.md`).
