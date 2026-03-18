# Shhh USB LLM

Terminal de Inteligencia Artificial portable que se ejecuta desde un USB en Windows. Se camufla como una sesion normal del Simbolo del sistema (`cmd.exe`) o PowerShell. Sin instalacion, sin internet, sin rastro.

Escribes tu pregunta en lenguaje natural, la IA responde en texto plano. Especializado en programacion (C++, Python, etc.) y razonamiento logico.

---

## Que contiene este repositorio

| Archivo | Descripcion |
|---------|-------------|
| `win_host.bat` | Script lanzador para CMD |
| `win_host.ps1` | Script lanzador para PowerShell |
| `README.md` | Esta guia |

Los modelos de IA y el motor de ejecucion NO estan incluidos porque pesan varios GB. Descargalos siguiendo los pasos de abajo.

---

## Requisitos del USB

| Caracteristica | Minimo | Recomendado |
|---------------|--------|-------------|
| Capacidad | 32 GB | 64 GB o 128 GB |
| Velocidad | USB 3.0 | USB 3.1 / 3.2 |
| Formato | exFAT | exFAT |
| Lectura secuencial | 100 MB/s | 200+ MB/s |

Ten en cuenta que el USB tambien llevara tus apuntes, PDFs, proyectos y demas archivos de trabajo. Calculo de espacio aproximado:

| Contenido | Espacio |
|-----------|---------|
| Motor (llama.cpp + DLLs) | ~500 MB |
| Modelo Qwen 3B | 2.0 GB |
| Modelo Qwen 7B | 4.3 GB |
| Modelo DeepSeek 7B | 4.7 GB |
| **Total IA** | **~11.5 GB** |
| Apuntes, PDFs, proyectos, etc. | Variable |

- **32 GB**: cabe la IA completa + ~18 GB para tus archivos.
- **64 GB**: espacio de sobra para todo. Lo mas recomendable.
- **128 GB**: si quieres llevar absolutamente todo encima.
- **USB 3.0 minimo**: con USB 2.0 el modelo tardaria varios minutos en cargar en vez de segundos.
- **exFAT obligatorio**: FAT32 no permite archivos de mas de 4GB (los modelos de 7B pesan ~4.7GB).

---

## Instalacion

### Paso 1: Formatear el USB

1. Conecta el pendrive al PC.
2. **Windows**: Clic derecho sobre la unidad en "Este equipo" > Formatear > Sistema de archivos: **exFAT** > Iniciar.
3. **macOS**: Abre "Utilidad de Discos" > Selecciona el USB > Borrar > Formato: **ExFAT**.

### Paso 2: Crear la carpeta oculta

Crea una carpeta llamada **`.sys_tools`** en la raiz del USB.

El punto al inicio del nombre (`.sys_tools`) la hace invisible automaticamente en macOS y Linux.

### Paso 3: Ocultar la carpeta en Windows

Abre CMD, navega a la raiz del USB y ejecuta:
```cmd
attrib +h +s +r .sys_tools
```
La carpeta desaparece del Explorador de archivos. Para acceder: `cd .sys_tools`.

Para volver a verla:
```cmd
attrib -h -s -r .sys_tools
```

En macOS, ademas del punto, puedes ejecutar:
```bash
chflags hidden /Volumes/TuUSB/.sys_tools
```

### Paso 4: Descargar el motor de IA

1. Descarga directa: [llama.cpp para Windows (Vulkan x64)](https://github.com/ggml-org/llama.cpp/releases/download/b8394/llama-b8394-bin-win-vulkan-x64.zip)
2. Abre el `.zip`.
3. Extrae **TODO el contenido** (todos los `.exe` y `.dll`) dentro de `.sys_tools`.
4. Comprueba que `llama-cli.exe` esta dentro.

### Paso 5: Descargar los modelos

Descarga al menos los modelos 1 y 2. Haz clic en "Descargar" para iniciar la descarga directa, luego mueve el archivo `.gguf` a `.sys_tools`.

| Opcion | Modelo | Para que sirve | RAM del PC | Tamaño | Descarga directa |
|--------|--------|---------------|-----------|--------|-------------------|
| 1 | Qwen2.5-Coder 3B | Codigo rapido (C++, Python, JS) | ~4 GB | 2.0 GB | [Descargar](https://huggingface.co/Qwen/Qwen2.5-Coder-3B-Instruct-GGUF/resolve/main/qwen2.5-coder-3b-instruct-q4_k_m.gguf?download=true) |
| 2 | Qwen2.5-Coder 7B | Codigo preciso y avanzado | ~8 GB | 4.3 GB | [Descargar](https://huggingface.co/bartowski/Qwen2.5-Coder-7B-Instruct-GGUF/resolve/main/Qwen2.5-Coder-7B-Instruct-Q4_K_M.gguf?download=true) |
| 3 | DeepSeek-R1 7B | Depuracion, logica, razonamiento | ~6 GB | 4.7 GB | [Descargar](https://huggingface.co/bartowski/DeepSeek-R1-Distill-Qwen-7B-GGUF/resolve/main/DeepSeek-R1-Distill-Qwen-7B-Q4_K_M.gguf?download=true) |
| 4 | Phi-4 Mini 3.8B | Razonamiento + codigo (Microsoft) | ~4 GB | 2.5 GB | [Descargar](https://huggingface.co/bartowski/microsoft_Phi-4-mini-instruct-GGUF/resolve/main/microsoft_Phi-4-mini-instruct-Q4_K_M.gguf?download=true) |
| 5 | Gemma 3 4B | Generalista: resumenes, idiomas | ~4 GB | 2.8 GB | [Descargar](https://huggingface.co/bartowski/google_gemma-3-4b-it-GGUF/resolve/main/google_gemma-3-4b-it-Q4_K_M.gguf?download=true) |

### Paso 6: Copiar los scripts

Descarga `win_host.bat` y `win_host.ps1` de este repositorio y ponlos en `.sys_tools`.

### Estructura final

```
(USB) D:\
 └── .sys_tools/                                        <- Oculta
      ├── llama-cli.exe
      ├── ggml-vulkan.dll
      ├── ggml-cpu-*.dll
      ├── (otros .dll del ZIP)
      ├── qwen2.5-coder-3b-instruct-q4_k_m.gguf        <- Modelo 1
      ├── Qwen2.5-Coder-7B-Instruct-Q4_K_M.gguf        <- Modelo 2
      ├── deepseek-r1-distill-qwen-7b-q4_k_m.gguf       <- Modelo 3
      ├── win_host.bat
      └── win_host.ps1
```

---

## Como usarlo

### Desde CMD

```cmd
D:
cd .sys_tools
win_host
```

### Desde PowerShell

```powershell
D:
cd .sys_tools
powershell -ExecutionPolicy Bypass -File .\win_host.ps1
```

(Sustituye `D:` por la letra de tu USB: `E:`, `F:`, etc.)

---

## Que pasa al ejecutarlo

1. Aparece un texto identico al de una consola de Windows real.
2. El modelo se carga en silencio (unos segundos).
3. Aparece el cursor: `C:\Users\Admin> ` (CMD) o `PS C:\Users\Admin> ` (PowerShell).
4. Escribes tu pregunta y pulsas Enter.
5. La IA responde en texto plano y devuelve el cursor.
6. Para salir: `Ctrl + C`.

Cualquiera que mire tu pantalla vera lo que parece una consola de Windows normal con salida de texto tecnico.

---

## Comandos

### Seleccion de modelo

**CMD:**
| Comando | Modelo |
|---------|--------|
| `win_host` | Qwen 3B (defecto) |
| `win_host 1` | Qwen 3B |
| `win_host 2` | Qwen 7B |
| `win_host 3` | DeepSeek R1 7B |
| `win_host 4` | Phi-4 Mini |
| `win_host 5` | Gemma 3 4B |

**PowerShell:**
| Comando | Modelo |
|---------|--------|
| `powershell -ExecutionPolicy Bypass -File .\win_host.ps1` | Qwen 3B (defecto) |
| `powershell -ExecutionPolicy Bypass -File .\win_host.ps1 2` | Qwen 7B |
| `powershell -ExecutionPolicy Bypass -File .\win_host.ps1 3` | DeepSeek R1 7B |

### Cuando usar cada modelo

| Opcion | Modelo | Situacion | Velocidad en i7 (CPU) |
|--------|--------|-----------|----------------------|
| 1 | **Qwen 3B** | Escribir codigo C++: clases, STL, punteros, ficheros, templates. Tu dia a dia. | Rapido (~10 tok/s) |
| 2 | **Qwen 7B** | Codigo C++ complejo y preciso. Cuando el 3B se queda corto. Necesita RAM. | Medio (~4 tok/s) |
| 3 | **DeepSeek R1** | Cuando algo no compila y no sabes por que. Segfaults, bugs logicos, algoritmos. | Medio (~4 tok/s) |
| 4 | **Phi-4 Mini** | Alternativa al Qwen 3B. Fuerte en matematicas y razonamiento. | Rapido (~9 tok/s) |
| 5 | **Gemma 3** | Preguntas generales, resumenes, traducciones. Menos preciso en codigo. | Rapido (~8 tok/s) |

### Controles dentro de la sesion

| Tecla / Accion | Que hace |
|----------------|----------|
| `Enter` | Enviar pregunta |
| Seleccionar texto + clic derecho | Copiar (CMD clasico) |
| Clic derecho | Pegar (CMD clasico) |
| `Ctrl + Shift + C` | Copiar (Windows Terminal) |
| `Ctrl + Shift + V` | Pegar (Windows Terminal) |
| Flecha arriba | Recuperar ultima pregunta |
| Cerrar la ventana (X) | Salir |

---

## Precision de los modelos

### Qwen2.5-Coder 3B (Opcion 1 - Recomendado)
- Generacion de codigo (HumanEval): **~65-70%** — Comparable a GPT-3.5.
- Domina C++ (STL, clases, herencia, punteros, templates, ficheros), Python, JavaScript.
- Flojea en codigo muy largo (+200 lineas) o librerias especificas (Boost, Qt).

### Qwen2.5-Coder 7B (Opcion 2 - El mas preciso)
- Generacion de codigo (HumanEval): **~83-88%** — Nivel GPT-4 en codigo puro.
- Escribe C++ casi perfecto. El mas preciso de todos.
- Requiere ~8GB de RAM solo para el modelo. Si Windows ya usa 3GB, necesitas que el PC tenga al menos 12GB para que no se congele.

### DeepSeek-R1 7B (Opcion 3 - Para depurar)
- Razonamiento logico (MATH/GSM8K): **~75-80%** — Nivel GPT-4o-mini en razonamiento.
- Su fuerte es ENTENDER por que algo falla, no escribir codigo desde cero.
- El script ya incluye instrucciones para que NO muestre su proceso de razonamiento interno, asi mantiene la apariencia de terminal limpia.

### Phi-4 Mini 3.8B (Opcion 4 - Microsoft)
- Razonamiento (MATH): **~70-75%**.
- Buen hibrido entre razonamiento y codigo. Ligero y rapido.
- Buena opcion si necesitas razonar pero no quieres cargar el DeepSeek de 7B.

### Gemma 3 4B (Opcion 5 - Google)
- Tareas generales (MMLU): **~65%**.
- Modelo generalista. Menos preciso en C++ que Qwen, pero util para redactar textos, resumir o traducir.

---

## Solucion de problemas

### Sale `[SYS_ERROR]` y se cierra
Abre el archivo `debug_log.txt` en `.sys_tools`. Ahi esta el error exacto.

| Error | Causa | Solucion |
|-------|-------|----------|
| `error: invalid argument` | Version de llama.cpp incompatible | Descarga la ultima version del ZIP |
| `failed to load model` | Archivo `.gguf` corrupto o incompleto | Vuelve a descargarlo |
| `out of memory` / `alloc failed` | El PC no tiene RAM suficiente | Usa un modelo mas pequeno (opcion 1 o 4) |
| `vulkan error` | El PC no tiene GPU Vulkan | No pasa nada, sigue funcionando en CPU |

### Al escribir se ejecuta como comando real
La IA ha crasheado y estas en la terminal real. Revisa `debug_log.txt` y prueba con un modelo mas ligero.

---

## Discrecion

- No instala nada en el PC. Todo vive en el USB.
- No necesita permisos de administrador.
- No necesita internet.
- Al desconectar el USB no queda rastro en el disco duro.
- La carpeta `.sys_tools` es invisible en cualquier explorador de archivos.
- Si te preguntan: "Estoy compilando un proyecto" o "Son los tests del compilador".

---

*Proyecto de proposito educativo.*
