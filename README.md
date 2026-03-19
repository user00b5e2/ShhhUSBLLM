# Shhh USB LLM

Terminal de IA portable desde USB en Windows. Se camufla como CMD o PowerShell. Sin instalacion, sin internet, sin rastro.

---

## Requisitos del USB

| Caracteristica | Minimo | Recomendado |
|---------------|--------|-------------|
| Capacidad | 32 GB | 64-128 GB |
| Velocidad | USB 3.0 (100+ MB/s) | USB 3.1/3.2 (200+ MB/s) |
| Formato | exFAT | exFAT |

Espacio aproximado:

| Contenido | Espacio |
|-----------|---------|
| Motor + DLLs | ~500 MB |
| Todo tier 4 GB (4 modelos) | ~8.4 GB |
| Todo tier 8 GB (3 modelos) | ~15.5 GB |
| Modelo 14B (tier 12 GB) | ~8.7 GB |
| **Todos los modelos** | **~25 GB** |
| Apuntes, PDFs, proyectos | Variable |

Con un USB de 64 GB caben todos los modelos + ~35 GB para tus archivos.

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

**Recomendado:** la primera vez que conectes el USB a Windows, abre CMD y ejecuta esto para blindar los atributos:
```cmd
D:
attrib +h +s +r "System Volume Information"
```
Esto annade los atributos de sistema + oculto + solo lectura. A partir de ahi, la carpeta es invisible en el Explorador de archivos incluso con "Mostrar archivos ocultos" activado (porque tiene el atributo sistema ademas del oculto).

Al conectar el USB a Windows, el SO NO borra el contenido de la carpeta. Puede que anada algun archivo pequeno suyo (como `IndexerVolumeGuid`), lo cual es perfecto: hace que los archivos parezcan aun mas legitimos.

Para acceder a ella:
```cmd
cd "System Volume Information"
```

### 3. Descargar el motor

[Descargar llama.cpp Windows Vulkan x64](https://github.com/ggml-org/llama.cpp/releases/download/b8394/llama-b8394-bin-win-vulkan-x64.zip) → Extraer TODO dentro de `System Volume Information`.

### 4. Disfrazar el ejecutable

```cmd
ren llama-cli.exe hostcfg.exe
```

### 5. Descargar y disfrazar los modelos

Descarga, mueve a la carpeta y renombra. No necesitas descargarlos todos, solo los que vayas a usar segun la RAM del PC.

**PCs con 4 GB de RAM:**

| Opcion | Modelo | Renombrar a | Peso | Descarga |
|--------|--------|-------------|------|----------|
| 0 | Qwen2.5-Coder 1.5B | `syscache_00.dat` | 1.1 GB | [Descargar](https://huggingface.co/Qwen/Qwen2.5-Coder-1.5B-Instruct-GGUF/resolve/main/qwen2.5-coder-1.5b-instruct-q4_k_m.gguf?download=true) |

**PCs con 6-8 GB de RAM:**

| Opcion | Modelo | Renombrar a | Peso | Descarga |
|--------|--------|-------------|------|----------|
| 1 | Qwen2.5-Coder 3B | `syscache_01.dat` | 2.0 GB | [Descargar](https://huggingface.co/Qwen/Qwen2.5-Coder-3B-Instruct-GGUF/resolve/main/qwen2.5-coder-3b-instruct-q4_k_m.gguf?download=true) |
| 2 | Phi-4 Mini 3.8B | `syscache_02.dat` | 2.5 GB | [Descargar](https://huggingface.co/bartowski/microsoft_Phi-4-mini-instruct-GGUF/resolve/main/microsoft_Phi-4-mini-instruct-Q4_K_M.gguf?download=true) |
| 3 | Gemma 3 4B | `syscache_03.dat` | 2.8 GB | [Descargar](https://huggingface.co/bartowski/google_gemma-3-4b-it-GGUF/resolve/main/google_gemma-3-4b-it-Q4_K_M.gguf?download=true) |
| 4 | **Qwen3.5-4B** ★ | `syscache_04.dat` | 3.1 GB | [Descargar](https://huggingface.co/unsloth/Qwen3.5-4B-GGUF/resolve/main/Qwen3.5-4B-Q4_K_M.gguf?download=true) |

**PCs con 10-12 GB de RAM:**

| Opcion | Modelo | Renombrar a | Peso | Descarga |
|--------|--------|-------------|------|----------|
| 5 | **Qwen2.5-Coder 7B** ★ | `syscache_05.dat` | 4.3 GB | [Descargar](https://huggingface.co/bartowski/Qwen2.5-Coder-7B-Instruct-GGUF/resolve/main/Qwen2.5-Coder-7B-Instruct-Q4_K_M.gguf?download=true) |
| 6 | DeepSeek-R1 7B | `syscache_06.dat` | 4.7 GB | [Descargar](https://huggingface.co/bartowski/DeepSeek-R1-Distill-Qwen-7B-GGUF/resolve/main/DeepSeek-R1-Distill-Qwen-7B-Q4_K_M.gguf?download=true) |
| 7 | **Qwen3.5-9B** ★ | `syscache_07.dat` | 6.5 GB | [Descargar](https://huggingface.co/unsloth/Qwen3.5-9B-GGUF/resolve/main/Qwen3.5-9B-Q4_K_M.gguf?download=true) |

**PCs con 16 GB de RAM:**

| Opcion | Modelo | Renombrar a | Peso | Descarga |
|--------|--------|-------------|------|----------|
| 8 | **Qwen2.5-Coder 14B** ★★ | `syscache_08.dat` | 8.7 GB | [Descargar](https://huggingface.co/bartowski/Qwen2.5-Coder-14B-Instruct-GGUF/resolve/main/Qwen2.5-Coder-14B-Instruct-Q4_K_M.gguf?download=true) |

★ = recomendado en su categoria. ★★ = lo mejor que existe en local.

Ejemplo de renombrado:
```cmd
ren qwen2.5-coder-1.5b-instruct-q4_k_m.gguf syscache_00.dat
```

### 6. Copiar los scripts

Copia `shhh.bat`, `shhh.ps1` y `shhhps.bat` dentro de la carpeta.

### Estructura final

```
(USB) D:\
 └── System Volume Information/    <- Invisible por defecto
      ├── hostcfg.exe              <- Motor disfrazado
      ├── (DLLs del ZIP)
      ├── syscache_00.dat          <- Qwen 1.5B
      ├── syscache_04.dat          <- Qwen3.5-4B (defecto)
      ├── syscache_05.dat          <- Qwen2.5-Coder 7B
      ├── (otros syscache_XX.dat)
      ├── shhh.bat
      ├── shhh.ps1
      └── shhhps.bat
```

---

## Como usarlo

```
D:
cd "System Volume Information"
shhh
```
(Sustituye `D:` por la letra de tu USB.)

Para aspecto PowerShell: `shhhps`

---

## Comandos

### shhh (CMD) / shhhps (PowerShell)

| Comando | Modelo | RAM del PC | HumanEval | Velocidad |
|---------|--------|-----------|-----------|-----------|
| `shhh 0` | Qwen2.5-Coder 1.5B | 4 GB | ~50% | Muy rapido |
| `shhh 1` | Qwen2.5-Coder 3B | 6 GB | ~67% | Rapido |
| `shhh 2` | Phi-4 Mini 3.8B | 6 GB | ~72% | Rapido |
| `shhh 3` | Gemma 3 4B | 6 GB | ~65% | Rapido |
| `shhh` o `shhh 4` | **Qwen3.5-4B** ★ | 8 GB | ~78% | Rapido |
| `shhh 5` | **Qwen2.5-Coder 7B** ★ | 10 GB | ~88% | Medio |
| `shhh 6` | DeepSeek R1 7B | 10 GB | ~78% | Medio |
| `shhh 7` | **Qwen3.5-9B** ★ | 12 GB | ~83% | Medio |
| `shhh 8` | **Qwen2.5-Coder 14B** ★★ | 16 GB | ~92% | Lento |

La columna "RAM del PC" es la RAM TOTAL que necesita el PC (sistema + modelo). Windows con un par de apps abiertas consume ~3 GB.

**Modo explicacion:** Antepon `e`. Ejemplo: `shhh e 5` = Qwen 7B con explicaciones breves.

Para PowerShell, sustituye `shhh` por `shhhps`. Ejemplo: `shhhps 5`, `shhhps e 5`.

### Diferencia shhh vs shhhps

| | `shhh` | `shhhps` |
|--|--------|----------|
| Apariencia | CMD (fondo negro) | PowerShell (fondo azul) |
| Funciona desde | CMD y PowerShell | Solo PowerShell |

Usa el que coincida con la terminal que ya tiene abierta el PC.

---

## Que modelo usar

### PC con 4 GB RAM → `shhh 0`
Qwen2.5-Coder 1.5B. Es limitado pero funciona. Genera snippets cortos de C/C++, Python. Se equivoca mas que los grandes pero es tu unica opcion con tan poca RAM.

### PC con 6-8 GB RAM → `shhh` (defecto = Qwen3.5-4B)
El mejor equilibrio calidad/velocidad. Marzo 2026, ultima generacion. Supera a modelos de 7B de 2024. Escribe C++ correcto la mayor parte del tiempo.

### PC con 10-12 GB RAM → `shhh 5` (Qwen2.5-Coder 7B)
El mejor modelo de codigo en 7B que existe. 88% en HumanEval, casi GPT-4o. Si el PC aguanta, este es el que quieres para codigo puro.

### PC con 16 GB RAM → `shhh 8` (Qwen2.5-Coder 14B)
92% en HumanEval. Nivel GPT-4o. Lo mejor que puedes correr en local. C++ casi perfecto.

### Para depurar → `shhh 6` (DeepSeek R1 7B)
No para escribir codigo, sino para entender POR QUE falla. Segfaults, errores logicos, algoritmos.

---

## Como escribir preguntas

- **Todo en una linea.** Cada Enter envia inmediatamente.
- **Se directo.** "funcion C++ que ordene vector con quicksort"
- **Pega codigo en una sola linea.** La IA lo entiende.
- **Evita acentos** si puedes.

### Controles

| Accion | Que hace |
|--------|----------|
| `Enter` | Enviar pregunta |
| Seleccionar + clic derecho | Copiar (CMD) |
| Clic derecho | Pegar (CMD) |
| `Ctrl + Shift + C/V` | Copiar/Pegar (Windows Terminal) |
| Flecha arriba | Ultima pregunta |
| Cerrar ventana (X) | Salir (borra historial) |

---

## Capas de sigilo

1. **Carpeta de sistema**: Los archivos viven dentro de `System Volume Information`, una carpeta que Windows crea automaticamente en cada unidad y que es invisible por defecto. Nadie la abre jamas. Si alguien la ve, asumira que es del propio Windows.

2. **Archivos disfrazados**: El motor se llama `hostcfg.exe` (parece un servicio de red). Los modelos se llaman `syscache_0X.dat` (parecen caches del sistema).

3. **Limitacion conocida: `llama.dll`**: El ZIP del motor trae DLLs como `ggml-vulkan.dll`, `ggml-cpu-*.dll` y `llama.dll`. No se pueden renombrar porque el motor las busca por nombre exacto. En la practica, `ggml` no significa nada para nadie, y `llama.dll` es un nombre generico que podria ser cualquier libreria. Solo alguien que conozca el proyecto llama.cpp lo reconoceria, y para llegar hasta ahi tendria que: (1) saber que la carpeta existe, (2) desactivar la proteccion de archivos de sistema, (3) entrar y leer nombres de DLLs. Improbable.

4. **Interfaz identica**: El titulo de la ventana, el texto de inicio y el prompt son identicos a una terminal real de Windows. No hay prefijos de "Usuario" ni "Asistente".

5. **Silencio absoluto**: Toda la salida tecnica del motor (carga, memoria, tiempos) se redirige a la nada. Los logs estan desactivados. La pantalla permanece limpia.

6. **Rutas dinamicas**: Los scripts detectan automaticamente la ruta del USB. Funciona sin importar que letra asigne Windows.

7. **Historial borrado**: Al cerrar, el historial de comandos de CMD y PowerShell se borra automaticamente. Nadie puede pulsar flecha arriba para ver que ejecutaste.

---

## Solucion de problemas

| Error | Solucion |
|-------|----------|
| `Core executable missing` | Renombraste `llama-cli.exe` a `hostcfg.exe`? |
| `Modulo no encontrado` | Renombraste el modelo a `syscache_0X.dat`? |
| Se cierra sin mostrar nada | Prueba un modelo mas ligero (0, 1 o 2) |
| Se congela o va muy lento | El PC no tiene RAM suficiente para ese modelo |

---

## Discrecion

- No instala nada. No necesita admin. No necesita internet.
- No deja rastro en disco duro ni en historial.
- Los archivos viven en una carpeta que Windows nunca muestra.
- Si te preguntan: "Estoy comprobando la integridad del disco."

---

*Proyecto de proposito educativo.*
