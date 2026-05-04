# Duraciones aproximadas en 8 GB sin GPU

Mediciones esperadas en **CPU AVX2 típica de portátil 2018–2023** (i5/i7 4–8 cores, sin GPU). En CPUs más modernas (i7-13xxx, Ryzen 7xxx) reduce ~30 %. En CPUs viejas (i3 / sin AVX2) **multiplica × 2–3**.

Cifras actualizadas para los modelos Qwen3 (mayo 2026). Cambios respecto a iter2 (Qwen2.5-Coder):

- Slot 1: 1.5B → **1.7B** dual-mode con `/no_think`. ~5 % más lento.
- Slots 2/3: 3B → **4B**. ~25 % más lento por turno pero produce código de calidad ≈ 7B-anterior.
- Slot 5: 7B → **8B**. ~15 % más lento, calidad ≈ 14B-anterior.

## Por slot

### Slot 1 (1.7B agent)

| Tarea | Tiempo |
|-------|--------|
| Carga inicial del modelo | 20–40 s |
| Edit puntual de un fichero (read + edit + done) | 10–22 s |
| Crear un fichero pequeño | 8–18 s |
| Read de fichero + análisis sin escribir | 5–12 s |

### Slot 2 (4B agent)

| Tarea | Tiempo |
|-------|--------|
| Carga inicial | 35–70 s |
| Edit puntual | 40–80 s |
| Edit en 2–3 ficheros | 2–4 min |
| Refactor con compilación + corrección | 4–10 min |

### Slot 3 (4B chat)

| Tarea | Tiempo |
|-------|--------|
| Carga inicial | 35–70 s |
| Respuesta corta (8 líneas) | 20–40 s |
| Respuesta larga (50 líneas) | 1.5–4 min |

### Slot 4 (1.7B chat)

| Tarea | Tiempo |
|-------|--------|
| Carga inicial | 20–40 s |
| Respuesta corta | 8–20 s |
| Respuesta larga (50 líneas) | 1–2 min |

### Slot 5 (8B agent — tareas grandes)

| Tarea | Tiempo |
|-------|--------|
| Carga inicial | 70–140 s |
| Lectura de MD de 5 páginas | 5–10 s (no hay generación) |
| Generar 1 fichero cpp de 200 líneas | 2.5–6 min |
| Generar 5 ficheros cpp desde un MD | **20–35 min** sin auto-corrección, **35–70 min** con compile + fix |
| Ronda de auto-corrección (g++ compile + leer errores + edit) | 2.5–6 min/ronda |

## Caso "5 cpp desde MD de 5 páginas" — desglose

Suposiciones: el MD detalla 5 ficheros (~200 líneas/cada uno) que deben compilar entre sí. Con Qwen3-8B en CPU 8 GB sin GPU.

| Paso | Tiempo aprox. |
|------|---------------|
| Carga del 8B (1ª vez) | 100 s |
| Lectura del MD | 5 s |
| Generación del 1.er .cpp (~200 líneas, 6 k tokens) | 3.5–6 min |
| Generación del 2.º | 3.5–6 min |
| Generación del 3.º | 3.5–6 min |
| Generación del 4.º | 3.5–6 min |
| Generación del 5.º | 3.5–6 min |
| Compilación inicial con `g++` | 10–30 s |
| Si hay errores: lectura, edit, recompilación × 2–3 rondas | 8–18 min |
| Llamada final a `done` | 5 s |
| **Total esperado** | **25–55 min** |

Realismo:

- Si todo va a la primera (poco realista): ~25 min.
- Si hay 2 rondas de corrección (típico): ~35–45 min.
- Si el modelo se atasca en alguna iteración: hasta 70 min antes de agotar `MaxIter=20`.

Notas sobre Qwen3-8B vs Qwen2.5-Coder-7B en este caso de uso:

- El 8B es ~15 % más lento por token, pero **necesita menos rondas de corrección** porque acerca calidad de un modelo 14B.
- En la práctica, el tiempo total tiende a ser **similar o ligeramente menor** que con el 7B antiguo, con resultado de mejor calidad.

## Memory pressure durante la tarea

| Estado | RAM `hostcfg.exe` | RAM `bgupd.exe` | Sistema total |
|--------|-------------------|------------------|---------------|
| Idle (modelo no cargado) | 0 | 30 MB | baseline + 30 MB |
| Slot 1 / 4 cargado | 1.4 GB | 30 MB | baseline + 1.4 GB |
| Slot 2 / 3 cargado | 2.9 GB | 30 MB | baseline + 2.9 GB |
| Slot 5 cargado | 5.7 GB | 30 MB | baseline + 5.7 GB |

Si tu baseline es 4 GB y cargas slot 5: total 9.7 GB. **Con 8 GB físicos eso swappea**, latencia se dispara × 3–5 fácilmente. Por eso recomiendo **cerrar Chrome / browser pesado** antes de slot 5. Si ni así cabe, fallback a `qwen3-8b-q3_k_m.gguf` (~4.0 GB).

## Trucos para acortar tiempos

1. **Pre-warm**: lanza `update 5` y deja que cargue mientras preparas el MD. Así la primera tarea ya no paga la carga.
2. **Slot 2 en vez de slot 5** para tareas que no requieran calidad máxima — 3–4× más rápido.
3. **Subir threads** con `-t N` no ayuda más allá de núcleos físicos. En i5 4-core con SMT, `-t 4` (no 8) es óptimo.
4. **Apagar Defender realtime scanning** durante la tarea baja overhead un 10–20 %. **No lo recomiendo de rutina** — sí para una sesión de trabajo controlada.
5. **Usar inglés en los prompts** con slot 1: el 1.7B comprende mejor inglés que español, igual que su antecesor.
