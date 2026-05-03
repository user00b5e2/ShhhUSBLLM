# Ejecución en Windows 8 GB

Pasos exactos para construir el USB y correrlo en un PC Windows. Asume Windows 10/11 64-bit Intel sin GPU como destino.

## Preparación (en una máquina con Go + bash)

Necesitas Go 1.20+ y `curl`. La preparación funciona en cualquier sistema (Linux, WSL, macOS); solo se cross-compila a Windows.

```bash
cd /ruta/al/proyecto

# 1. Compilar el binario Windows (bgupd.exe)
./build.sh

# 2. Descargar llama-server.exe build AVX2 y renombrarlo a hostcfg.exe
./download-llama.sh

# 3. Modelos básicos (slots 1, 2, 4) — ~5 GB total
./download-models.sh

# 4. Si vas a usar slot 5 (7B), añade ~4.7 GB:
WITH_LARGE=1 ./download-models.sh
```

Resultado esperado en `bin/`:

```
bgupd.exe         (~7 MB, harness Windows)
hostcfg.exe       (~50–80 MB, llama-server renombrado)
*.dll             (varios, las que llama.cpp necesite)
```

Resultado en `models/`:
```
qwen2.5-coder-1.5b-instruct-q4_k_m.gguf   (~1.1 GB)
qwen2.5-coder-3b-instruct-q4_k_m.gguf     (~2.0 GB)
qwen2.5-1.5b-instruct-q4_k_m.gguf         (~1.1 GB)
qwen2.5-coder-7b-instruct-q4_k_m.gguf     (~4.7 GB) si usaste WITH_LARGE=1
```

## Estructura final del USB

```
USB:\
├── README.md                ← guía de uso (sin tocar)
├── QUICKSTART.md            ← referencia comandos (sin tocar)
├── docs\                    ← esta documentación técnica
│   ├── 00-INDICE.md
│   ├── 01-CONTEXTO.md
│   ├── 02-ARQUITECTURA.md
│   ├── 03-TECNICAS-SIGILO.md
│   ├── 04-MODELOS-Y-RAM.md
│   ├── 05-DURACIONES.md
│   ├── 06-EJECUCION-WINDOWS.md  ← este
│   ├── 07-CAMBIOS-ITER2.md
│   └── 08-LIMITACIONES.md
├── update.bat               ← entrypoint CMD (renombrado de shhh-agent.bat)
├── update.ps1               ← entrypoint PowerShell
├── shhh.bat / shhh.ps1      ← modo chat clásico (proyecto original, sin tocar)
├── bin\
│   ├── bgupd.exe
│   ├── hostcfg.exe
│   └── *.dll
└── models\
    └── *.gguf
```

## Copia al USB

Linux / WSL:
```bash
USB=/mnt/MIUSB               # o donde monte tu sistema
mkdir -p "$USB/LLM"
cp -R bin models docs *.bat *.ps1 *.md "$USB/LLM/"
sync && umount "$USB"
```

Windows (cuando preparas el USB desde el propio Windows):
```cmd
robocopy . D:\LLM /E /XF *.go *.sh /XD shhh-agent .git
```
(copia todo el árbol a `D:\LLM`, excluyendo fuentes Go y scripts de build).

## Primera ejecución en Windows

### Desde la terminal integrada de VS Code (recomendado)

VS Code en Windows abre PowerShell por defecto. Pulsa `` Ctrl+` `` y:

```powershell
PS C:\Users\demo> D:                       # tu letra de USB
PS D:\> cd LLM
PS D:\LLM> .\update.ps1
```

Verás aparecer un prompt **idéntico al tuyo** (porque `update.ps1` capturó tu `prompt` real antes de pasar control). Tipea a ciegas y enter:

```
crea hola.md con el texto "hola"
```

Espera 10–20 s (modelo 1.5B). Cuando aparezca el prompt nuevo, abre `hola.md` en VS Code → contenido correcto.

### Desde CMD

```cmd
D:
cd \LLM
update.bat
```

### Forzar slot manualmente

```powershell
.\update.ps1 1     # 1.5B agent (rápido)
.\update.ps1 2     # 3B agent (calidad media)
.\update.ps1 3     # 3B chat (sin tocar ficheros)
.\update.ps1 4     # 1.5B chat
.\update.ps1 5     # 7B agent (largo, requiere modelo descargado)
```

### Salir

Tipea `exit` (a ciegas) + Enter, o pulsa `Ctrl+C`. Eso limpia pantalla, mata el backend y vuelves a la shell real.

## Ver lo que pasa la primera vez

Si quieres confirmar que funciona antes de tirar a ciegas:

```powershell
$env:SHHH_VERBOSE = "1"
.\update.ps1 1
```

Verás cada llamada de tool (`[write_file] {...}`) en pantalla. Para volver al sigilo: `Remove-Item Env:\SHHH_VERBOSE`.

## Diagnóstico

```powershell
# ¿hay backend vivo?
Get-Process hostcfg -ErrorAction SilentlyContinue
Test-NetConnection -ComputerName 127.0.0.1 -Port 8765

# matar backend
.\update.ps1 --stop

# si quedó huérfano
Stop-Process -Name hostcfg -Force
```

## Rendimiento esperado (8 GB Intel sin GPU)

Mide tu primera tarea con cronómetro. Compara con [`05-DURACIONES.md`](05-DURACIONES.md). Si va > 2× lo que dice esa tabla:

1. ¿Tu CPU tiene AVX2? `wmic cpu get caption,featureset` o `Get-ComputerInfo`.
2. ¿Estás corriendo el `hostcfg.exe` AVX2 build (no el genérico)? `download-llama.sh` baja el AVX2.
3. ¿RAM al límite? `Get-Counter '\Memory\Available MBytes'`.
4. ¿Antivirus está escaneando cada read? Excepción para la carpeta del USB durante la prueba.

## Lo que NO va a funcionar

- **Slot 5 con Chrome abierto y 8 GB**: el 7B necesita ~5 GB, Chrome típicamente 2–3 GB, Windows ~4 GB → swap.
- **PCs Intel sin AVX2** (pre-2013, Atom, Pentium N): llama.cpp arrancará pero a tok/s inutilizables. No es target.
- **PCs ARM Windows (Surface Pro X)**: necesitarías build ARM de llama.cpp, no el AVX2. Nuestro USB no lo trae por defecto.
