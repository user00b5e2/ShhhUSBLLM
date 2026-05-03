# Documentación técnica — shhh-agent

Este directorio NO duplica el `README.md` (instalación + uso). Aquí está la documentación **explicativa**: cómo funciona el harness por dentro, qué técnicas de discreción usa para no llamar la atención en pantalla, qué modelos puede correr en 8 GB y cuánto tarda cada cosa.

| Documento | Para qué sirve |
|-----------|----------------|
| [`01-CONTEXTO.md`](01-CONTEXTO.md) | Por qué existe el harness, qué problema resuelve, qué NO es. |
| [`02-ARQUITECTURA.md`](02-ARQUITECTURA.md) | Mapa de procesos, módulos del harness Go, flujo de un turno completo. |
| [`03-TECNICAS-SIGILO.md`](03-TECNICAS-SIGILO.md) | Catálogo de las técnicas de discreción visual que el binario implementa, con código real. |
| [`04-MODELOS-Y-RAM.md`](04-MODELOS-Y-RAM.md) | Slots 1–5, RAM, tok/s, tipo de tarea. Tabla específica para 8 GB. |
| [`05-DURACIONES.md`](05-DURACIONES.md) | Cuánto tarda cada tipo de tarea en CPU sin GPU 8 GB. |
| [`06-EJECUCION-WINDOWS.md`](06-EJECUCION-WINDOWS.md) | Pasos exactos para correr el USB en un Windows 8 GB sin GPU. |
| [`07-CAMBIOS-ITER2.md`](07-CAMBIOS-ITER2.md) | Qué cambió respecto a la primera versión y por qué. |
| [`08-LIMITACIONES.md`](08-LIMITACIONES.md) | Lo que **NO** funciona o falla, sin maquillaje. |

Lee los documentos en orden si quieres entenderlo todo. Si solo vienes a probar, usa el `README.md` de la raíz.
