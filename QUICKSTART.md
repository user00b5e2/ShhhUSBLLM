# shhh-agent — Guía rápida (Windows 8 GB)

Referencia de comandos para usar y diagnosticar el agente desde el USB en Windows.

Para la documentación técnica completa, ver `docs/`.

---

## Estructura del USB

```
USB:\LLM\
├── README.md
├── QUICKSTART.md            ← este fichero
├── docs\                    ← documentación técnica
├── update.bat               ← entrypoint CMD
├── update.ps1               ← entrypoint PowerShell
├── shhh.bat / shhh.ps1      ← modo chat clásico (proyecto original)
├── bin\
│   ├── bgupd.exe            ← harness Windows (~7 MB)
│   ├── hostcfg.exe          ← llama-server renombrado
│   └── *.dll                ← DLLs de llama.cpp
└── models\
    └── *.gguf
```

---

## Lanzar el agente

Desde la terminal integrada de VS Code (PowerShell por defecto):

```powershell
PS C:\> D:
PS D:\> cd LLM
PS D:\LLM> .\update.ps1
```

Desde CMD:

```cmd
D:
cd \LLM
update.bat
```

---

## Slots de modelo

| Slot | Modelo | RAM | Modo | Cuándo usarlo |
|------|--------|-----|------|---------------|
| 1 | Qwen3-1.7B Q4 | ~1.4 GB | agente | **Default agente.** Edits puntuales rápidos. |
| 2 | Qwen3-4B-Instruct-2507 Q4 | ~2.9 GB | agente | Edits multi-fichero o lógica compleja. |
| 3 | Qwen3-4B-Instruct-2507 Q4 | ~2.9 GB | chat | Explicar código, sin tocar ficheros. |
| 4 | Qwen3-1.7B Q4 | ~1.4 GB | chat | Chat ligero, PCs muy justos. |
| 5 | Qwen3-8B Q4 | ~5.7 GB | agente | Tareas grandes (5 cpp desde MD). Cierra Chrome. |

---

## Comandos del binario

```cmd
update                  REM advisor decide slot
update 1                REM forzar slot 1..5
update --stop           REM mata el hostcfg en background, libera RAM
update --unsafe         REM SIN guardarraíles (no usar)
update --verbose        REM equivalente a SHHH_VERBOSE=1
update --port 9999      REM cambiar puerto si 8765 está ocupado
update --once "prompt"  REM una sola tarea sin REPL, con output visible
```

Variables de entorno (`set NAME=value` en CMD, `$env:NAME = "value"` en PS):

```
SHHH_VERBOSE=1       imprime cada paso del agente y errores completos
SHHH_SHOW_RESULT=1   solo el resumen final del turno (entre silencio y verbose)
SHHH_PROMPT="..."    override manual del prompt falso (útil con oh-my-posh)
```

Dentro del REPL:
- `Enter` con texto vacío → reimprime prompt y sigue (mimetiza shell).
- `exit` o `quit` (tipeado a ciegas) → sale, limpia pantalla, restaura prompt real.
- `Ctrl+C` → panic button: limpia pantalla y sale.

---

## Diagnóstico (PowerShell)

```powershell
# ¿hay backend vivo?
Get-Process hostcfg -ErrorAction SilentlyContinue
Test-NetConnection -ComputerName 127.0.0.1 -Port 8765

# ¿qué consume?
Get-Process hostcfg | Select-Object Name, WS, CPU

# matar todo
.\update.ps1 --stop
Stop-Process -Name hostcfg -Force -ErrorAction SilentlyContinue
Remove-Item $env:TEMP\hostcfg.lock -ErrorAction SilentlyContinue
```

---

## Permisos y seguridad

Por defecto el agente está **scoped al CWD** (directorio donde lanzas `update`). Bloqueado:

- Salir del CWD: paths con `..`, absolutos, symlinks que apunten fuera.
- Directorios: `.git\`, `node_modules\`, `vendor\`, `target\`, `dist\`, `build\`, `__pycache__\`, `.venv\`, `venv\`.
- Ficheros: `.env*`, `id_rsa*`, `id_dsa*`, `id_ed25519*`, `.npmrc`, `.netrc`.
- Sufijos: `.key`, `.pem`, `.pfx`, `.p12`, `.keystore`.
- Comandos peligrosos: `format`, `shutdown`, `rm -rf`, `del /s`, `reg delete`, `net user`, `curl`, `wget`, `iex`, `cmd /c`, `start`, etc.

Flag `--unsafe` desactiva todo. **No la uses.**

---

## Resolución de problemas

| Síntoma | Causa probable | Solución |
|---------|----------------|----------|
| Aparece `!` tras Enter | Modelo no descargado o backend no arranca | `set SHHH_VERBOSE=1` y reintentar para ver el error |
| `[missing model: ...]` | Falta el GGUF en `models\` | Descargar el modelo del slot que pediste |
| `[server error: ...]` | `hostcfg.exe` no encontrado o crash | Verificar `bin\hostcfg.exe` y DLLs presentes |
| `[health timeout]` | Modelo tarda en cargar (USB lento) | Espera; segunda vez será instantáneo |
| Pega un texto y aparece eco | `term.MakeRaw` falló | Lanza desde Windows Terminal o VS Code, no conhost antiguo |
| El agente hace cosas raras con tu prompt | Modelo 1.5B mal interpretó español | Usa inglés, o slot 2 (3B) |
| `'old' string is not unique` | El edit_file pide texto único | Reformula con más contexto literal |
| RAM al límite | Demasiado modelo para la máquina | Slot 1, cierra Chrome, `update --stop` cuando termines |

---

## Re-construir / actualizar el USB

Esto se hace en una máquina de desarrollo (con Go instalado):

```bash
# en la máquina de dev
cd /ruta/al/proyecto

./build.sh                                    # cross-compile bgupd.exe
./download-llama.sh                           # llama-server.exe → hostcfg.exe + DLLs
./download-models.sh                          # modelos básicos (slots 1, 2, 4)
WITH_LARGE=1 ./download-models.sh             # añade slot 5 (7B, ~4.7 GB)
```

Luego copias el árbol al USB y lo enchufas al Windows objetivo.
