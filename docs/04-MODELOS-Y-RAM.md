# Modelos y RAM en 8 GB

Target: **Windows 10/11 Intel 8 GB sin GPU dedicada**.

## Presupuesto de RAM

Una máquina típica de 8 GB con Windows + VS Code + browser ya consume ~4–5 GB de baseline. Lo que tienes para el modelo + harness oscila entre **2 GB (laptop con muchas pestañas) y 4 GB (sistema limpio)**. Esa es la cifra real que limita los slots.

Recomendación práctica: **cerrar Chrome / Edge antes de slot 5**. La diferencia entre 1.2 GB de RAM disponible y 3 GB cambia si el modelo carga o swappea.

## Tabla por slot (familia Qwen3, mayo 2026)

| Slot | Modelo | Pesos | KV cache (ctx) | RAM total | Tok/s CPU AVX2 | Modo | Para qué |
|------|--------|-------|----------------|-----------|----------------|------|----------|
| 1 | Qwen3-1.7B Q4_K_M | 1.1 GB | ~250 MB (4k) | **~1.4 GB** | 9–14 | agente | Edits puntuales rápidos. Default agente. |
| 2 | Qwen3-4B-Instruct-2507 Q4_K_M | 2.5 GB | ~400 MB (8k) | **~2.9 GB** | 2.5–4 | agente | Edits multi-fichero medianos. Lectura de contexto razonable. |
| 3 | Qwen3-4B-Instruct-2507 Q4_K_M | 2.5 GB | ~400 MB (8k) | **~2.9 GB** | 2.5–4 | chat | Explicar código, sin tocar ficheros. |
| 4 | Qwen3-1.7B Q4_K_M | 1.1 GB | ~250 MB (4k) | **~1.4 GB** | 9–14 | chat | Chat ligero. PCs muy justos de RAM. |
| 5 | Qwen3-8B Q4_K_M | 5.0 GB | ~700 MB (16k) | **~5.7 GB** | 1.2–2.5 | agente | Tareas grandes (5 cpp desde MD, refactors profundos). Requiere cerrar apps pesadas. |

Beneficio de la nueva tabla: **slots 1+4 comparten un solo GGUF** (Qwen3-1.7B), y **slots 2+3 comparten otro** (Qwen3-4B-Instruct-2507). Sólo 3 ficheros de modelo en disco para los 5 slots:

```
qwen3-1.7b-q4_k_m.gguf                   (~1.1 GB)  → slots 1, 4
qwen3-4b-instruct-2507-q4_k_m.gguf       (~2.5 GB)  → slots 2, 3
qwen3-8b-q4_k_m.gguf                     (~5.0 GB)  → slot 5 (opt-in)
```

Total disco con WITH_LARGE=1: ~8.6 GB. Sin slot 5: ~3.6 GB.

## Por qué Qwen3 (vs Qwen2.5-Coder)

La familia Qwen3 (Alibaba, abril 2025) y su refinamiento Qwen3-2507 (julio 2025) son la mejor opción dense pequeña a fecha de hoy:

- **Qwen3-1.7B base ≈ Qwen2.5-Coder-3B** en code benchmarks.
- **Qwen3-4B base ≈ Qwen2.5-Coder-7B**.
- **Qwen3-8B base ≈ Qwen2.5-Coder-14B**.

Es decir, a paridad de RAM ganamos aproximadamente una "talla" de calidad. La variante **Instruct-2507** del 4B es non-thinking nativo (no emite `<think>...</think>` antes de responder), perfecto para nuestro tool-calling con XML estricto. Para 1.7B y 8B, los modelos son dual-mode y nuestro harness inyecta automáticamente la directiva `/no_think` en el system prompt.

## Alternativas evaluadas y descartadas

| Modelo | Razón |
|--------|-------|
| Qwen3-Coder-Next | MoE 80B total / 3B activos. Necesita ~32 GB para cargar; **fuera de target 8 GB**. |
| Qwen3.5-9B | Q4 ~5.7 GB; al límite de slot 5. Si en algún momento queremos un slot 6, candidato natural. |
| Qwen3.6-27B / GLM-5 / Kimi K2.5 | Modelos grandes 2026; fuera de target. |
| Codestral-22B | ~13 GB en Q4. Inviable. |
| Llama-3.3-8B / Mistral-Small-3 | Generalistas; coding peor que Qwen3-8B según benchmarks 2026. |
| Phi-4 mini (3.8B) | Decente, pero por debajo de Qwen3-4B-Instruct-2507 en code generation. |
| DeepSeek-Coder-V2-Lite (16B MoE) | 5.5 GB Q4. Buena opción pero requiere llama.cpp con MoE actualizado y Qwen3-8B le iguala con menos complicación. |
| Qwen2.5-Coder | Familia anterior (oct 2024). Sustituida por Qwen3 en esta iteración. |

## Cómo afecta el quant

| Quant | Pesos relativos | Calidad | Cuándo usar |
|-------|-----------------|---------|-------------|
| Q8_0 | 100 % | Mejor | Si te sobra RAM |
| Q4_K_M | 56 % | -2–3 % en benchmarks | **Default razonable** |
| Q3_K_M | 48 % | -4–5 % | Si Q4 no cabe (slot 5 en RAM muy justa) |
| Q2_K | 30 % | -10 %, errores notables | Solo en emergencias |

Por defecto Q4_K_M (no Q8) para que el 8B quepa en 8 GB.

## Contexto (`-c`) — el otro consumidor de RAM

El KV-cache crece linealmente con el contexto. Cada 4 k tokens son ~250–700 MB según el modelo. Por eso:

- 1.7B chat / agente: 4 k es suficiente.
- 4B agente / chat: 8 k cubre la mayoría de casos multi-fichero.
- 8B agente con MDs largos: 16 k. Más sería tirar RAM.

Configurado en `advisor.go` (`CtxSize` por slot).

## Carga inicial del modelo

| Slot | Tiempo de carga (CPU AVX2 desde SSD) | Desde USB (lento) |
|------|---------------------------------------|-------------------|
| 1 / 4 | 5–10 s | 20–30 s |
| 2 / 3 | 12–18 s | 35–55 s |
| 5 | 30–50 s | 70–110 s |

El backend queda residente mientras el REPL esté abierto (`server.go`), así sólo pagas la carga la primera vez de cada sesión. `update --stop` lo libera al instante.

## Sobre el flag `/no_think`

Qwen3 dual-mode (1.7B y 8B) emite por defecto un bloque `<think>...</think>` antes de la respuesta final, similar a un proceso de cadena de pensamiento. Eso:

1. **Duplica fácilmente la latencia** en CPU sin GPU (genera tokens "internos" que no son útiles para nuestro tool-calling).
2. **Puede romper el parser** si el `<think>` se cuela antes del `<tool>...</tool>`.

El harness inyecta `/no_think` al inicio del system prompt cuando el slot tiene `Qwen3DualMode=true` (advisor.go). El modelo entonces salta directamente a la respuesta.

Para el 4B-Instruct-2507, no hace falta el flag: ese modelo está fine-tuned **solo en non-thinking** desde origen.
