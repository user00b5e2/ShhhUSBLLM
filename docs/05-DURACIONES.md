# Duraciones aproximadas en 8 GB sin GPU

Mediciones esperadas en **CPU AVX2 típica de portátil 2018–2023** (i5/i7 4–8 cores, sin GPU). En CPUs más modernas (i7-13xxx, Ryzen 7xxx) reduce ~30 %. En CPUs viejas (i3 / sin AVX2) **multiplica × 2–3**.

## Por slot

### Slot 1 (1.5B agent)

| Tarea | Tiempo |
|-------|--------|
| Carga inicial del modelo | 20–40 s |
| Edit puntual de un fichero (read + edit + done) | 10–20 s |
| Crear un fichero pequeño | 8–15 s |
| Read de fichero + análisis sin escribir | 5–10 s |

### Slot 2 (3B agent)

| Tarea | Tiempo |
|-------|--------|
| Carga inicial | 30–60 s |
| Edit puntual | 30–60 s |
| Edit en 2–3 ficheros | 1.5–3 min |
| Refactor con compilación + corrección | 3–8 min |

### Slot 3 (3B chat)

| Tarea | Tiempo |
|-------|--------|
| Carga inicial | 30–60 s |
| Respuesta corta (8 líneas) | 15–30 s |
| Respuesta larga (50 líneas) | 1–3 min |

### Slot 5 (7B agent — tareas grandes)

| Tarea | Tiempo |
|-------|--------|
| Carga inicial | 60–120 s |
| Lectura de MD de 5 páginas | 5–10 s (no hay generación) |
| Generar 1 fichero cpp de 200 líneas | 2–5 min |
| Generar 5 ficheros cpp desde un MD | **15–30 min** sin auto-corrección, **30–60 min** con compile + fix |
| Ronda de auto-corrección (g++ compile + leer errores + edit) | 2–5 min/ronda |

## Caso "5 cpp desde MD de 5 páginas" — desglose

Suposiciones: el MD detalla 5 ficheros (~200 líneas/cada uno) que deben compilar entre sí.

| Paso | Tiempo aprox. |
|------|---------------|
| Carga del 7B (1ª vez) | 90 s |
| Lectura del MD | 5 s |
| Generación del 1.er .cpp (~200 líneas, 6 k tokens) | 3–5 min |
| Generación del 2.º | 3–5 min |
| Generación del 3.º | 3–5 min |
| Generación del 4.º | 3–5 min |
| Generación del 5.º | 3–5 min |
| Compilación inicial con `g++` | 10–30 s |
| Si hay errores: lectura, edit, recompilación × 2–3 rondas | 6–15 min |
| Llamada final a `done` | 5 s |
| **Total esperado** | **20–50 min** |

Realismo:

- Si todo va a la primera (poco realista): ~20 min.
- Si hay 2 rondas de corrección (típico): ~30–40 min.
- Si el modelo se atasca en alguna iteración (el 7B es bueno pero no perfecto): hasta 60 min antes de agotar `MaxIter=20`.

## Memory pressure durante la tarea

| Estado | RAM `hostcfg.exe` | RAM `bgupd.exe` | Sistema total |
|--------|-------------------|------------------|---------------|
| Idle (modelo no cargado) | 0 | 30 MB | baseline + 30 MB |
| Slot 1 cargado | 1.3 GB | 30 MB | baseline + 1.3 GB |
| Slot 5 cargado | 5.2 GB | 30 MB | baseline + 5.2 GB |

Si tu baseline es 4 GB y cargas slot 5: total 9.2 GB. **Con 8 GB físicos eso swappea**, latencia se dispara × 3–5 fácilmente. Por eso recomiendo **cerrar Chrome / browser pesado** antes de slot 5.

## Trucos para acortar tiempos

1. **Pre-warm**: lanza `update 5` y deja que cargue mientras preparas el MD. Así la primera tarea ya no paga la carga.
2. **Slot 2 en vez de slot 5** para tareas que no requieran calidad máxima — 3–4× más rápido.
3. **Subir threads** con `-t N` no ayuda más allá de núcleos físicos. En i5 4-core con SMT, `-t 4` (no 8) es óptimo.
4. **Apagar Defender realtime scanning** durante la tarea baja overhead un 10–20 %. **No lo recomiendo de rutina** — sí para una demo controlada.
