# Modelos y RAM en 8 GB

Target: **Windows 10/11 Intel 8 GB sin GPU dedicada**.

## Presupuesto de RAM

Una máquina típica de 8 GB con Windows + VS Code + browser ya consume ~4–5 GB de baseline. Lo que tienes para el modelo + harness oscila entre **2 GB (laptop con muchas pestañas) y 4 GB (sistema limpio)**. Esa es la cifra real que limita los slots.

Recomendación práctica: **cerrar Chrome / Edge antes de slot 5**. La diferencia entre 1.2 GB de RAM disponible y 3 GB cambia si el modelo carga o swappea.

## Tabla por slot

| Slot | Modelo | Pesos | KV cache (ctx) | RAM total | Tok/s CPU AVX2 | Modo | Para qué |
|------|--------|-------|----------------|-----------|----------------|------|----------|
| 1 | Qwen2.5-Coder-1.5B Q4_K_M | 1.0 GB | ~250 MB (4k) | **~1.3 GB** | 10–15 | agente | Edits triviales / typos. Lane rápido. |
| 2 | Qwen2.5-Coder-3B Q4_K_M | 2.0 GB | ~400 MB (8k) | **~2.4 GB** | 3–5 | agente | **Default** del advisor para tareas de coding. Sweet spot 8 GB. |
| 3 | Qwen2.5-Coder-3B Q4_K_M | 2.0 GB | ~400 MB (8k) | **~2.4 GB** | 3–5 | chat | Explicar código, sin tocar ficheros. |
| 4 | Qwen2.5-Coder-1.5B Q4_K_M | 1.0 GB | ~250 MB (4k) | **~1.3 GB** | 10–15 | chat | Chat ligero. Reusa GGUF de slot 1. |
| 5 | Qwen2.5-Coder-7B Q4_K_M | 4.5 GB | ~700 MB (16k) | **~5.2 GB** | 1.5–3 | agente | Tareas grandes (5 cpp desde MD, refactors profundos). Cierra apps pesadas. |

Beneficio: solo **3 ficheros GGUF distintos** para los 5 slots (slot 1 y 4 comparten 1.5B, slot 2 y 3 comparten 3B).

```
qwen2.5-coder-1.5b-instruct-q4_k_m.gguf   (~1.0 GB)  → slots 1, 4
qwen2.5-coder-3b-instruct-q4_k_m.gguf     (~2.0 GB)  → slots 2, 3
qwen2.5-coder-7b-instruct-q4_k_m.gguf     (~4.7 GB)  → slot 5 (opt-in)
```

Total disco con WITH_LARGE=1: ~7.7 GB. Sin slot 5: ~3.0 GB.

## Por qué Qwen2.5-Coder y no Qwen3 dense

Mirando 2026, hay tres familias relevantes para 8 GB CPU sin GPU:

1. **Qwen2.5-Coder** (Alibaba, octubre 2024): familia coder-específica, fine-tuned sobre repositorios de código. Tamaños 1.5B, 3B, 7B encajan en 8 GB.
2. **Qwen3 dense** (Alibaba, abril 2025): familia generalista con thinking mode dual. 1.7B, 4B, 8B.
3. **Qwen3-Coder**: existe pero sólo en variantes grandes (MoE 80B / 3B-active, 30B-A3B, 480B). **Ninguna cabe en 8 GB**.

**Para nuestro use case (agente con tool-calling XML estricto + escribir código), Qwen2.5-Coder gana**:

- **Tool-calling más fiable**: fine-tuning sobre código real → respeta formatos estrictos sin caracteres extra. En pruebas comparativas, Qwen3-1.7B emitió ocasionalmente sufijos `>` después del JSON; Qwen2.5-Coder-1.5B no.
- **Sin thinking mode**: Qwen3 dual-mode requiere directiva `/no_think` runtime para evitar tokens `<think>...</think>` antes de la respuesta. Qwen-Coder no tiene ese problema.
- **Validado empíricamente**: los smoke tests del proyecto pasaron con Qwen2.5-Coder.
- **Calidad coding-pura competitiva**: Alibaba publicó comparaciones de Qwen2.5-Coder-7B con GPT-4o donde empata o gana en EvalPlus, BigCodeBench, LiveCodeBench. Qwen3-8B generalista no llega a ese nivel en coding puro.

Qwen3 dense sería preferible si el use case fuera chat general, multilingüismo o razonamiento abstracto — pero aquí el target es coding y agente.

## Alternativas evaluadas y descartadas

| Modelo | Razón |
|--------|-------|
| Qwen3-Coder-Next | MoE 80B / 3B activos. Necesita ~32 GB para cargar pesos; **fuera de target 8 GB**. |
| Qwen3-Coder-30B-A3B | 30B MoE, ~17 GB Q4. Inviable. |
| Qwen3 dense (1.7B / 4B / 8B) | Generalistas; tool-calling menos fiable que Qwen-Coder para nuestro XML estricto (ver arriba). |
| Codestral-22B | ~13 GB Q4. Inviable. |
| DeepSeek-Coder-V2-Lite (16B MoE) | 5.5 GB Q4. Buena calidad pero requiere build llama.cpp con MoE actualizado y Qwen2.5-Coder-7B le iguala con menos complicación. |
| Phi-4 mini (3.8B) | Decente, pero por debajo de Qwen2.5-Coder-3B en code generation. |
| Llama-3.1-8B / Mistral-Small-3 | Generalistas; coding peor que Qwen-Coder a paridad de tamaño. |

## Cómo afecta el quant

| Quant | Pesos relativos | Calidad | Cuándo usar |
|-------|-----------------|---------|-------------|
| Q8_0 | 100 % | Mejor | Si te sobra RAM |
| Q4_K_M | 56 % | -2–3 % en benchmarks | **Default razonable** |
| Q3_K_S | 42 % | -5 % | Si Q4 no cabe (slot 5 en RAM muy justa) |
| Q2_K | 30 % | -10 %, errores notables | Solo en emergencias |

Por defecto Q4_K_M (no Q8) para que el 7B quepa en 8 GB.

## Contexto (`-c`) — el otro consumidor de RAM

El KV-cache crece linealmente con el contexto. Cada 4 k tokens son ~250–700 MB según el modelo. Por eso:

- 1.5B / 3B chat: 4–8 k es suficiente.
- 7B agente con MDs largos: 16 k. Más sería tirar RAM.

Configurado en `advisor.go` (`CtxSize` por slot).

## Carga inicial del modelo

| Slot | Tiempo de carga (CPU AVX2 desde SSD) | Desde USB (lento) |
|------|---------------------------------------|-------------------|
| 1 / 4 | 5–10 s | 20–30 s |
| 2 / 3 | 10–15 s | 30–45 s |
| 5 | 25–40 s | 60–90 s |

El backend queda residente mientras el REPL esté abierto (`server.go`), así sólo pagas la carga la primera vez de cada sesión. `update --stop` lo libera al instante.

## Default del advisor

Cuando ejecutas `update` sin número, el advisor heurístico (en `advisor.go`) decide:

- Verbo de acción en el prompt (`crea`, `edita`, `arregla`, `refactoriza`, ...) → **slot 2 (3B-Coder)** por defecto.
- Marcadores triviales (`typo`, `one-liner`, `rename de variable`) → slot 1 (1.5B-Coder, más rápido).
- Marcadores de tarea grande (`spec.md`, `5 cpp`, `g++`, `compile`) → slot 5 (7B-Coder).
- Marcadores de explicación (`explica`, `cómo funciona`) → slot 3 (3B-Coder chat).
- Ambiguo → slot 3 (chat — más seguro, no toca ficheros).

Forzar slot manual con `update N` siempre prevalece sobre el advisor.
