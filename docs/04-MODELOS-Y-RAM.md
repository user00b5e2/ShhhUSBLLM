# Modelos y RAM en 8 GB

Target: **Windows 10/11 Intel 8 GB sin GPU dedicada**.

## Presupuesto de RAM

Una máquina típica de 8 GB con Windows + VS Code + browser ya consume ~4–5 GB de baseline. Lo que tienes para el modelo + harness oscila entre **2 GB (laptop con muchas pestañas) y 4 GB (sistema limpio)**. Esa es la cifra real que limita los slots.

Recomendación práctica: **cerrar Chrome / Edge antes de slot 5**. La diferencia entre 1.2 GB de RAM disponible y 3 GB cambia si el modelo carga o swappea.

## Tabla por slot

| Slot | Modelo | Pesos | KV cache (ctx) | RAM total | Tok/s CPU AVX2 | Modo | Para qué |
|------|--------|-------|----------------|-----------|----------------|------|----------|
| 1 | Qwen2.5-Coder-1.5B Q4_K_M | 1.0 GB | ~250 MB (4k) | **~1.3 GB** | 10–15 | agente | Edits puntuales rápidos. Default agente. |
| 2 | Qwen2.5-Coder-3B Q4_K_M | 2.0 GB | ~400 MB (8k) | **~2.4 GB** | 3–5 | agente | Edits multi-fichero medianos. Lectura de contexto razonable. |
| 3 | Qwen2.5-Coder-3B Q4_K_M | 2.0 GB | ~400 MB (8k) | **~2.4 GB** | 3–5 | chat | Explicar código, sin tocar ficheros. |
| 4 | Qwen2.5-1.5B Q4_K_M (instruct) | 1.0 GB | ~250 MB (4k) | **~1.3 GB** | 10–15 | chat | Chat ligero. PCs muy justos de RAM. |
| 5 | Qwen2.5-Coder-7B Q4_K_M | 4.5 GB | ~700 MB (16k) | **~5.2 GB** | 1.5–3 | agente | Tareas grandes (5 cpp desde MD, refactors profundos). Requiere cerrar apps pesadas. |

## Por qué Qwen2.5-Coder

Es la familia open-weights con mejor rendimiento código en 2024–2025 a tamaños ≤ 7B. Mejor que CodeLlama, mejor que DeepSeek-Coder-V1, mejor que Phi-3 en benchmarks de coding (HumanEval, MBPP, MultiPL-E). Apache 2.0 license — uso personal y comercial.

## Alternativas que NO uso (y por qué)

| Modelo | Razón |
|--------|-------|
| Codestral-22B | ~13 GB en Q4. Inviable en 8 GB. |
| Llama-3.1-70B | Inviable. |
| DeepSeek-Coder-V2-Lite (16B MoE) | 5.5 GB Q4. Calidad equivalente a 7B-dense. **Buena opción alternativa**, pero requiere llama.cpp build con MoE support actualizado. Si lo quieres como slot 6, dímelo. |
| Phi-3.5-mini (3.8B) | Decente, pero peor en código que Qwen-Coder-3B en mis tests. |
| TinyLlama, OpenLlama, etc. | No siguen tool-calling protocol bien, ni siquiera con XML. |

## Cómo afecta el quant

| Quant | Pesos relativos | Calidad | Cuándo usar |
|-------|-----------------|---------|-------------|
| Q8_0 | 100 % | Mejor | Si te sobra RAM |
| Q4_K_M | 56 % | -2–3 % en benchmarks | **Default razonable** |
| Q3_K_S | 42 % | -5 % | Si Q4 no cabe |
| Q2_K | 30 % | -10 %, errores notables | Solo en emergencias |

Por defecto bajo a Q4_K_M (no Q8) para que el 7B quepa en 8 GB.

## Contexto (`-c`) — el otro consumidor de RAM

El KV-cache crece linealmente con el contexto. Cada 4 k tokens son ~250–700 MB según el modelo. Por eso:

- 1.5B / 3B chat: 4–8 k es suficiente.
- 7B agente con MDs largos: 16 k. Más sería tirar RAM.

Configurado en `advisor.go:46` (`CtxSize` por slot).

## Carga inicial del modelo

| Slot | Tiempo de carga (CPU AVX2 desde SSD) | Desde USB (lento) |
|------|---------------------------------------|-------------------|
| 1 | 5–10 s | 20–30 s |
| 2 / 3 | 10–15 s | 30–45 s |
| 5 | 25–40 s | 60–90 s |

El backend queda residente entre invocaciones (ver `server.go`), así solo pagas la carga la primera vez. `update --stop` lo libera.
