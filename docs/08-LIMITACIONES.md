# Limitaciones honestas

Esto **no funciona al 100 %** y no quiero venderlo como tal. Lista exhaustiva.

## No probado exhaustivamente en Windows real

La validación de la lógica del agente (parsing XML, tool dispatch, guardarraíles, manejo de errores) se hizo en una máquina de desarrollo. La integración específica con `conhost.exe`, Windows Terminal y la terminal integrada de VS Code-on-Windows tiene **riesgo desconocido** hasta una prueba completa en hardware Intel real.

Cosas que pueden fallar específicamente en Windows:

- **`ESC[8m`**: lo respeta Windows Terminal moderno y VS Code (ambos usan ConPTY + xterm.js). En `conhost.exe` clásico depende de la versión de Windows. En Windows 10 < 1809 no lo respeta y los caracteres tecleados sí se ven.
- **`term.MakeRaw` en Windows**: depende del subsistema de consola. `golang.org/x/term` lo soporta pero hay edge cases con redirección.
- **Captura del `prompt` PS** desde `update.ps1`: si tu PS tiene `oh-my-posh` con render asíncrono, lo capturado puede ser un placeholder o cadena vacía. Cubre el fallback.

## Modelos pequeños fallan en formato XML

Qwen-Coder-1.5B se sale del formato `<tool>...</tool><args>...</args>` de vez en cuando — escribe markdown, comenta el plan, etc. El parser detecta y manda un re-prompt automático, pero **no es infalible**. En sesiones largas con 1.5B verás ocasionalmente un turno que termina sin haber hecho nada.

Mitigación: usar slot 2 (3B) para cualquier cosa no trivial.

## Auto-corrección de C++ no es 100 %

El loop "compila → lee errores → edita → recompila" funciona en el caso típico, pero:

- Si el error de `g++` cita líneas que el modelo confunde con líneas de otro fichero, edita lo equivocado.
- Si el `old` que necesita el `edit_file` no es único en el fichero, el agente lo reintenta — pero a veces se atasca.
- Errores semánticos (compila pero no hace lo que pide el MD) **no los detecta**, porque no hay tests automáticos. Esos te toca pillarlos a ti.

Tasa de éxito esperada en Intel 8 GB sin GPU:

- 1 fichero cpp simple desde spec corta: ~85 %.
- 5 ficheros cpp interrelacionados desde MD detallado: ~50–70 %, con 2–3 rondas de auto-corrección.

## Antivirus puede dar falso positivo

`bgupd.exe` y `hostcfg.exe` son binarios sin firma comercial. Windows Defender o algunos AV de terceros pueden marcarlos por heurística (binario sin firmar que abre puerto local, spawn de proceso hijo). Es un falso positivo conocido para llama.cpp; añadir excepción manual a la carpeta del USB lo resuelve.

Si tu PC tiene EDR corporativo, no es el target de este harness. El proyecto está pensado para tu propia máquina personal.

## CPU spike y RAM spike son visibles

No hay forma de ocultar:

- 4 cores al 80–100 % durante inferencia.
- 1.5–5 GB de RAM por proceso `hostcfg.exe`.
- ~50 % del rendimiento del disco si el modelo está cargado desde USB lento.

Cualquier monitor de actividad lo enseña. Eso también es **parte del valor educativo**: enseñar que "running locally" no es invisible.

## El historial de comandos sí queda

`update.bat` o `.\update.ps1` queda en:

- `doskey /history` (CMD)
- `Get-History` y `(Get-PSReadlineOption).HistorySavePath` (PowerShell — esto se persiste a fichero por defecto).

No lo limpiamos automáticamente. Si quieres mantenerlo limpio, ejecuta `Clear-History` antes de cerrar la sesión, o configura `Set-PSReadlineOption -HistorySaveStyle SaveNothing` en tu `$PROFILE` si quieres que no se persista nunca.

## El advisor a veces clasifica mal

Ejemplos de fallos conocidos:

- "haz un md diciendo hola" — empieza por "haz" pero "diciendo" no es action verb fuerte → puede caer en chat.
- "what does this do?" en español dicho como "que pasa con esto" — "que" no está en mis chat markers exactos.

Forzar slot manual (`update 1`/`2`/`3`/`5`) elimina el problema.

## Ctrl+C no siempre limpia

Si el agente está en medio de una llamada HTTP a `hostcfg.exe`, Ctrl+C dispara el handler en goroutine que mata el server y sale, **pero la línea de generación que estaba en curso puede dejar un residuo en pantalla** antes del `cls`. Visible 200–500 ms. Aceptable pero no perfecto.

## El backend no tiene timeout de idle

Si lanzas `update`, haces una tarea, y cierras el REPL con Ctrl+C: `defer StopServer()` mata el backend → bien.

Si lanzas `update`, haces tarea, **dejas la ventana abierta** y te vas a comer: el backend sigue cargado consumiendo 1.5–5 GB hasta que vuelvas y cierres. No hay "idle shutdown after 10 min".

Workaround: `update --stop` desde otra terminal libera al instante.

## Path traversal con symlinks raros

`filepath.EvalSymlinks` cubre el caso normal. Pero hay edge cases:

- En Windows, junctions (`mklink /J`) se comportan distinto que symlinks Unix. No los he probado a fondo.
- Si el CWD es ya un symlink que apunta fuera del "workspace pretendido", el guard considera el target del symlink como CWD. Eso es deliberado y correcto, pero si **pretendías** que el agente quedara confinado al directorio original, no al resuelto, te llevarías una sorpresa.

## TODO razonable (no lo hago si no me lo pides)

- Test unitario del advisor.
- Compilation flag `-trimpath` ya está en build.sh, pero `--ldflags=-buildid=` para no dejar build IDs en el binario es un toque más.
- Empaquetar las DLLs de llama.cpp dentro del `.exe` con MSYS o link estático — daría 1 fichero y nada más en `bin/`.
- Sustituir `hostcfg.exe` por una build "barebones" de llama.cpp sin opciones que el harness no usa — reduce footprint y firma.
- Detectar si `Get-PSReadlineOption` tiene `HistorySaveStyle = SaveNothing` y advertir si no.
