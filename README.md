# Shhh USB LLM

Terminal de IA portable desde USB en Windows. Se camufla como CMD o PowerShell. Sin instalacion, sin internet, sin rastro.

---

## Requisitos del USB

| Caracteristica | Minimo | Recomendado |
|---------------|--------|-------------|
| Formato | exFAT | exFAT |
| Capacidad | 32 GB | 64 GB |
| Velocidad | USB 2.0 | USB 3.0+ |

**Espacio necesario segun los modelos que descargues:**

| Tier | Modelos | Espacio total (motor + modelos) |
|------|---------|----------------------------------|
| Solo 4 GB RAM | Motor + Qwen 1.5B | ~2 GB |
| Solo 6 GB RAM | Motor + un modelo de 3-4B | ~3-4 GB |
| Solo 8 GB RAM | Motor + Qwen3.5-4B | ~4 GB |
| Solo 10 GB RAM | Motor + Qwen2.5-Coder 7B | ~6 GB |
| Solo 12 GB RAM | Motor + Qwen3.5-9B | ~7 GB |
| Solo 16 GB RAM | Motor + Qwen2.5-Coder 14B | ~10 GB |
| TODOS los modelos | Motor + los 9 modelos | ~37 GB |

---

## Instalacion

### 1. Formatear en exFAT

- **Windows**: Clic derecho en la unidad > Formatear > exFAT.
- **macOS**: Utilidad de Discos > Borrar > ExFAT.

### 2. Crear la carpeta oculta

Windows crea automaticamente una carpeta llamada `System Volume Information` en cada unidad. Es invisible por defecto y nadie la abre jamas. Vamos a usarla.

**Desde macOS (Terminal):**
```bash
mkdir "/Volumes/TuUSB/System Volume Information"
chflags hidden "/Volumes/TuUSB/System Volume Information"
```
(Sustituye `TuUSB` por el nombre de tu pendrive.)

- `chflags hidden` marca la carpeta como invisible en macOS (Finder) y ademas en exFAT establece el atributo oculto que Windows respeta.
- El nombre `System Volume Information` hace que Windows la trate como carpeta de sistema propia, ocultandola automaticamente.
- En Linux, los gestores de archivos (Nautilus, Dolphin, Thunar) reconocen este nombre y la ocultan.

**Recomendado:** la primera vez que conectes el USB a Windows, abre CMD y ejecuta:
```cmd
D:
attrib +h +s +r "System Volume Information"
```
Esto annade los atributos de sistema + oculto + solo lectura. La carpeta sera invisible incluso con "Mostrar archivos ocultos" activado (porque tiene el atributo sistema ademas del oculto).

Al conectar el USB a Windows, el SO NO borra el contenido de la carpeta. Puede que anada algun archivo pequeno suyo (como `IndexerVolumeGuid`), lo cual es perfecto: hace que los archivos parezcan aun mas legitimos.

Es perfecta porque:
- Windows la oculta automaticamente en el Explorador de archivos.
- Tiene atributos de sistema, oculto y solo lectura.
- Ningun usuario normal la abre.
- Si alguien la ve, asumira que es del propio sistema operativo.

Para acceder a ella desde CMD:
```cmd
D:
cd "System Volume Information"
```

### 3. Descargar el motor

Hay dos versiones. Usa la que corresponda al PC:

| Version | Velocidad | Compatibilidad | Descarga |
|---------|-----------|----------------|----------|
| **CPU** (recomendada) | Normal | Funciona en TODOS los PCs | [Descargar](https://github.com/ggml-org/llama.cpp/releases/download/b8429/llama-b8429-bin-win-cpu-x64.zip) |
| Vulkan (GPU) | Mas rapida | Solo PCs con GPU Vulkan | [Descargar](https://github.com/ggml-org/llama.cpp/releases/download/b8429/llama-b8429-bin-win-vulkan-x64.zip) |

Si la version Vulkan se cierra sin mostrar nada, usa la CPU. Extrae TODO el contenido del ZIP dentro de `System Volume Information`.

### 4. Disfrazar el ejecutable

```cmd
ren llama-cli.exe hostcfg.exe
```

### 5. Descargar y disfrazar los modelos

Descarga, mueve a la carpeta y renombra. Solo necesitas los que vayas a usar segun la RAM del PC.

**PCs con 4 GB de RAM:**

| Comando | Modelo | Precision | Renombrar a | Peso | Descarga |
|---------|--------|-----------|-------------|------|----------|
| `shhh 0` | Qwen2.5-Coder 1.5B | ~45% HumanEval | `syscache_00.dat` | 1.1 GB | [Descargar](https://huggingface.co/Qwen/Qwen2.5-Coder-1.5B-Instruct-GGUF/resolve/main/qwen2.5-coder-1.5b-instruct-q4_k_m.gguf?download=true) |

**PCs con 6 GB de RAM:**

| Comando | Modelo | Precision | Renombrar a | Peso | Descarga |
|---------|--------|-----------|-------------|------|----------|
| `shhh 1` | Qwen2.5-Coder 3B | ~67% HumanEval | `syscache_01.dat` | 2.0 GB | [Descargar](https://huggingface.co/Qwen/Qwen2.5-Coder-3B-Instruct-GGUF/resolve/main/qwen2.5-coder-3b-instruct-q4_k_m.gguf?download=true) |
| `shhh 2` | Phi-4 Mini 3.8B | ~72% HumanEval | `syscache_02.dat` | 2.5 GB | [Descargar](https://huggingface.co/unsloth/Phi-4-mini-instruct-GGUF/resolve/main/Phi-4-mini-instruct-Q4_K_M.gguf?download=true) |
| `shhh 3` | Gemma 3 4B | ~70% HumanEval | `syscache_03.dat` | 3.0 GB | [Descargar](https://huggingface.co/unsloth/gemma-3-4b-it-GGUF/resolve/main/gemma-3-4b-it-Q4_K_M.gguf?download=true) |

**PCs con 8 GB de RAM:**

| Comando | Modelo | Precision | Renombrar a | Peso | Descarga |
|---------|--------|-----------|-------------|------|----------|
| `shhh` / `shhh 4` | Qwen1.5-4B (RECOMENDADO) | ~78% HumanEval | `syscache_04.dat` | 2.6 GB | [Descargar](https://huggingface.co/Qwen/Qwen1.5-4B-Chat-GGUF/resolve/main/qwen1_5-4b-chat-q4_k_m.gguf?download=true) |

**PCs con 10 GB de RAM:**

| Comando | Modelo | Precision | Renombrar a | Peso | Descarga |
|---------|--------|-----------|-------------|------|----------|
| `shhh 5` | Qwen2.5-Coder 7B (RECOMENDADO) | ~84% HumanEval | `syscache_05.dat` | 4.7 GB | [Descargar](https://huggingface.co/Qwen/Qwen2.5-Coder-7B-Instruct-GGUF/resolve/main/qwen2.5-coder-7b-instruct-q4_k_m.gguf?download=true) |
| `shhh 6` | DeepSeek R1 7B | ~75% HumanEval | `syscache_06.dat` | 4.9 GB | [Descargar](https://huggingface.co/unsloth/DeepSeek-R1-Distill-Qwen-7B-GGUF/resolve/main/DeepSeek-R1-Distill-Qwen-7B-Q4_K_M.gguf?download=true) |

**PCs con 12 GB de RAM:**

| Comando | Modelo | Precision | Renombrar a | Peso | Descarga |
|---------|--------|-----------|-------------|------|----------|
| `shhh 7` | Llama 3.1 8B (RECOMENDADO) | ~86% HumanEval | `syscache_07.dat` | 4.9 GB | [Descargar](https://huggingface.co/bartowski/Meta-Llama-3.1-8B-Instruct-GGUF/resolve/main/Meta-Llama-3.1-8B-Instruct-Q4_K_M.gguf?download=true) |

**PCs con 16 GB de RAM:**

| Comando | Modelo | Precision | Renombrar a | Peso | Descarga |
|---------|--------|-----------|-------------|------|----------|
| `shhh 8` | Qwen2.5-Coder 14B (EL MEJOR) | ~92% HumanEval | `syscache_08.dat` | 9.0 GB | [Descargar](https://huggingface.co/Qwen/Qwen2.5-Coder-14B-Instruct-GGUF/resolve/main/qwen2.5-coder-14b-instruct-q4_k_m.gguf?download=true) |

### 6. Copiar los scripts

Copia `shhh.bat`, `shhh.ps1` y `shhhps.bat` dentro de `System Volume Information` junto al motor y los modelos.

---

## Uso

Abre CMD o PowerShell, navega a la carpeta y ejecuta:

```cmd
D:
cd "System Volume Information"
shhh
```

**IMPORTANTE:** Al ejecutar, la pantalla quedara en negro durante unos segundos mientras el modelo carga. Es normal. A los ~60 segundos la pantalla se limpia automaticamente y aparece el prompt listo para escribir. No toques nada durante la carga.

### Modos disponibles

| Modo | Letra | Ejemplo | Que hace |
|------|-------|---------|----------|
| Codigo | *(ninguna)* | `shhh` / `shhh 5` | Devuelve solo codigo, sin explicaciones |
| Explicar | `e` | `shhh e` / `shhh e 5` | Explicacion breve (3 lineas max) |
| Pensar | `t` | `shhh t` / `shhh t 5` | Muestra el razonamiento completo del modelo |

**CMD — modo codigo:**
| Comando | Modelo |
|---------|--------|
| `shhh` | Qwen1.5-4B (defecto) |
| `shhh 0` | Qwen2.5-Coder 1.5B |
| `shhh 1` | Qwen2.5-Coder 3B |
| `shhh 2` | Phi-4 Mini |
| `shhh 3` | Gemma 3 4B |
| `shhh 4` | Qwen1.5-4B |
| `shhh 5` | Qwen2.5-Coder 7B |
| `shhh 6` | DeepSeek R1 7B |
| `shhh 7` | Llama 3.1 8B |
| `shhh 8` | Qwen2.5-Coder 14B |

**CMD — modo explicacion:**
| Comando | Modelo |
|---------|--------|
| `shhh e` | Qwen1.5-4B (defecto) |
| `shhh e 0` a `shhh e 8` | Igual que arriba |

**CMD — modo razonamiento (muestra el pensamiento del modelo):**
| Comando | Modelo |
|---------|--------|
| `shhh t` | Qwen1.5-4B (defecto) |
| `shhh t 0` a `shhh t 8` | Igual que arriba |

**PowerShell — modo codigo (aspecto PS):**
| Comando | Modelo |
|---------|--------|
| `shhhps` | Qwen1.5-4B (defecto) |
| `shhhps 0` a `shhhps 8` | Igual que arriba |

**PowerShell — modo explicacion:**
| Comando | Modelo |
|---------|--------|
| `shhhps e` | Qwen1.5-4B (defecto) |
| `shhhps e 0` a `shhhps e 8` | Igual que arriba |

**PowerShell — modo razonamiento:**
| Comando | Modelo |
|---------|--------|
| `shhhps t` | Qwen1.5-4B (defecto) |
| `shhhps t 0` a `shhhps t 8` | Igual que arriba |

### Cuando usar shhh vs shhhps

| | `shhh` | `shhhps` |
|---|--------|----------|
| Aspecto | CMD (Simbolo del sistema) | PowerShell |
| Texto inicial | `Microsoft Windows [Version 10.0...]` | `Windows PowerShell Copyright (C)...` |
| Prompt | `C:\>` | `PS C:\>` |
| Cuando usarlo | Si el PC tiene CMD abierto | Si el PC tiene PowerShell abierto |

El objetivo es que la ventana parezca lo que YA esta abierto en el PC para no levantar sospechas.

---

## Capas de sigilo

1. **Carpeta de sistema**: Los archivos viven dentro de `System Volume Information`, invisible por defecto. Si alguien la ve, asumira que es de Windows.

2. **Archivos disfrazados**: El motor se llama `hostcfg.exe` (parece servicio de red). Los modelos se llaman `syscache_0X.dat` (parecen caches del sistema).

3. **DLLs del motor**: El ZIP trae DLLs como `ggml-cpu-*.dll` y `llama.dll`. No se pueden renombrar porque el motor las busca por nombre exacto. En la practica solo alguien que conozca llama.cpp las reconoceria, y para eso tendria que desactivar la proteccion de archivos de sistema y leer nombres de DLLs dentro de una carpeta oculta de sistema. Improbable.

4. **Interfaz identica**: El titulo de la ventana, el texto de inicio y el prompt (`D:\ruta>`) son identicos a una terminal real de Windows. El prompt se genera dinamicamente usando `%CD%`.

5. **Banner invisible**: El motor muestra un banner al cargar (ASCII art, info del build). Para ocultarlo, el script utiliza el codigo ANSI de invisibilidad (`ESC[8m`) ANTES de ejecutar el motor, volviendo el texto 100% transparente sin importar el color de fondo de la terminal (funciona perfectamente en VSCode). Un proceso en segundo plano restaura la visibilidad a los 60 segundos, limpia la pantalla y reimprime el header falso de Windows.

6. **Razonamiento oculto**: En modo codigo, el modelo usa `--reasoning-budget 0` para activar su ruta de razonamiento sin generar texto de pensamiento visible. En modo think (`t`), el razonamiento se muestra completo.

7. **Silencio absoluto**: Toda la salida tecnica del motor (carga de modelo, memoria, tensores, tiempos) se redirige a la nada (`2>nul`). Los logs internos estan desactivados (`--log-disable`). Los colores del motor estan desactivados (`--color off`). Los tiempos de respuesta estan ocultos (`--no-show-timings`).

8. **Rutas dinamicas**: Los scripts detectan automaticamente la ruta del USB. Funciona sin importar que letra asigne Windows.

9. **Historial borrado**: Al cerrar, el historial de comandos de CMD y PowerShell se borra automaticamente.

---

## Solucion de problemas

| Problema | Solucion |
|----------|----------|
| `hostcfg.exe not found` | Renombraste `llama-cli.exe` a `hostcfg.exe`? |
| `syscache_0X.dat not found` | Renombraste el modelo a `syscache_0X.dat`? |
| Se cierra sin mostrar nada | Usa la version **CPU** del motor en vez de Vulkan |
| Pantalla negra mucho rato | Normal, el modelo esta cargando. Espera hasta 60 seg |
| Se congela o va muy lento | El PC no tiene RAM suficiente, prueba un modelo mas ligero (0, 1 o 2) |
| Sale texto de pensamiento | Asegurate de NO usar modo `t`. El modo normal ya lo oculta |
| Prompt en verde | Asegurate de que el script tiene `--color off` |
| Texto raro o basura | Prueba otro modelo, no todos son compatibles con todas las versiones del motor |

---

## Discrecion

- No dejes la ventana abierta cuando no la uses.
- Desde la terminal integrada de VSCode es todavia mas discreto: parece que estas compilando.
- Usa el modelo mas ligero que te sirva para que responda rapido.
- Si alguien se acerca, pulsa Ctrl+C para parar la respuesta y escribe `exit` para cerrar.

---

## Modo agente (`update`)

Adicional al modo chat (`shhh`), el USB incluye un **harness de agente** que da al modelo permisos para leer, escribir, editar y ejecutar comandos dentro del workspace donde lo lances. Pensado para no tocar nada con el ratón: tipeas a ciegas qué quieres que pase y el agente edita los ficheros por ti.

Para la documentación técnica completa, ver `docs/`.

### Uso

```cmd
D:\System Volume Information> update
```

o desde PowerShell:

```powershell
PS D:\System Volume Information> .\update.ps1
```

| Forma | Qué hace |
|-------|----------|
| `update` | El advisor elige modelo según la frase que escribas. **Default coding: 3B-Coder.** |
| `update 1` | Fuerza Qwen2.5-Coder-1.5B Q4 (rápido, edits triviales / typos). |
| `update 2` | Fuerza Qwen2.5-Coder-3B Q4 (default coding agent). |
| `update 3` | Modo chat con 3B-Coder (no toca ficheros). |
| `update 4` | Modo chat con 1.5B-Coder (PCs muy justos). |
| `update 5` | Qwen2.5-Coder-7B Q4 (tareas grandes, requiere descargar modelo opcional). |
| `update --stop` | Mata el `hostcfg.exe` que quedó cargado en background. |

### Cómo funciona el sigilo

- Misma técnica visual que `shhh`: prompt falso (capturado del PS real) + `ESC[8m`.
- **Además** la entrada del usuario es invisible (modo "raw" sin echo): aunque alguien mire por encima del hombro o estés compartiendo pantalla, no ven lo que tecleas ni cuando pegas con el portapapeles.
- Por defecto, **cero output del agente** entre prompts: solo aparece el siguiente prompt cuando termina. `set SHHH_SHOW_RESULT=1` para ver el resumen del turno; `set SHHH_VERBOSE=1` para la traza completa.
- `Ctrl+C` durante la ejecución limpia la pantalla y reimprime un prompt limpio (panic button).

### Permisos y guardarraíles

Por defecto el agente está confinado al directorio actual del workspace. Bloqueado:

- Escapar del CWD (rutas con `..`, absolutas o symlinks que apunten fuera).
- Tocar `.git/`, `.env*`, `*.key`, `*.pem`, `id_rsa*`, `node_modules/` y similares.
- Comandos peligrosos (`format`, `shutdown`, `rm -rf`, `reg delete`, `net user`, `curl`, `wget`, `iex`, ...).

Flag `--unsafe` desactiva los guardarraíles. **No se recomienda** para uso real: una alucinación del modelo de 1.5B puede borrar trabajo.

### Requisitos de RAM

Pensado para Windows 8 GB sin GPU. En CPU AVX2 típica:

| Slot | RAM modelo | tok/s aprox. | Latencia 1 edit |
|------|-----------|--------------|-----------------|
| 1 (1.5B) | ~1.3 GB | 10–15 | 10–20 s |
| 2 (3B) | ~2.4 GB | 3–5 | 40–80 s |
| 5 (7B) | ~5.2 GB | 1.5–3 | 2–5 min/fichero (cierra Chrome) |

El backend (`hostcfg.exe`) queda residente entre invocaciones del REPL para no pagar la carga del modelo cada vez (~30–40 s desde USB la primera vez). Se cierra automáticamente al salir del REPL (`exit` o `Ctrl+C`). `update --stop` lo libera al instante desde otra terminal.

### Limitaciones conocidas

- Modelos pequeños (1.5B) ocasionalmente se salen del formato XML; el harness reintenta una vez. Si falla, escribe la petición de otra forma.
- `edit_file` exige que el texto a sustituir aparezca exactamente una vez. Si el modelo lo confunde, te lo dirá y reintentará leyendo el fichero.
- No hay historial entre sesiones (decisión deliberada, igual que `shhh`).
- Un antivirus puede marcar `bgupd.exe` o `hostcfg.exe` por ser binarios sin firmar; añadir excepción manualmente o asumir falso positivo.
