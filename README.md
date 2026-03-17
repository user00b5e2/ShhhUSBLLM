# Shhh USB LLM: Terminal AI Encubierta
Este proyecto convierte cualquier unidad USB en una terminal de Inteligencia Artificial ejecutable de forma nativa ("bare-metal") en Windows, sin necesidad de instalación ni internet. Está diseñado para ocultarse como una sesión estándar del "Símbolo del sistema" (`cmd.exe`) y responder **exclusivamente con código plano o respuestas técnicas**, ideal para programadores y depuración (especialmente C++).

La IA se auto-limita mediante la inyección de un "System Prompt" para comportarse como una máquina fría: cero saludos, cero explicaciones, cero interacciones humanas. Entra código, sale código.

---

## Instalación y Preparación del USB

Este repositorio **solo contiene el script lanzador (`win_host.bat`)**. Para que funcione, necesitas descargar el motor de ejecución y los modelos de IA (los cuales no están incluidos aquí porque pesan varios Gigabytes).

Sigue estos pasos para crear tu unidad:

### Paso 1: Configurar la unidad USB
1. Usa un pendrive rápido (USB 3.0 o superior recomendado).
2. Formatéalo en **exFAT**. Esto es crucial porque permite almacenar archivos individuales de más de 4GB (necesario para los modelos) y es compatible nativamente tanto con Windows como con macOS/Linux.
3. En la raíz del USB, crea una carpeta llamada `sys_tools` (o cualquier otro nombre discreto).

### Paso 2: Descargar el motor (`llama.cpp`)
1. [Descargar Directo: Archivo ZIP de llama.cpp (Versión Vulkan HW-Acc)](https://github.com/ggml-org/llama.cpp/releases/download/b8390/llama-b8390-bin-win-vulkan-x64.zip)
2. Descomprime **todo el contenido** del `.zip` dentro de la carpeta `sys_tools` de tu USB. Asegúrate de que el archivo `llama-cli.exe` quede junto al script `.bat`.

### Paso 3: Descargar los Modelos de Lenguaje (GGUF)
Descarga los modelos que quieras usar y guárdalos en la carpeta `sys_tools`. Haz clic en los siguientes enlaces para que comience la **descarga automática y directa**:

*   **(Opción 1) Precisión en Código C++ (Para PCs con ~8GB RAM):**
    *   [Descargar Directo: qwen2.5-coder-7b-instruct-q4_k_m.gguf (4.3 GB)](https://huggingface.co/bartowski/Qwen2.5-Coder-7B-Instruct-GGUF/resolve/main/Qwen2.5-Coder-7B-Instruct-Q4_K_M.gguf?download=true)
*   **(Opción 2) Depuración y Lógica Compleja (Para PCs con ~8GB RAM):**
    *   [Descargar Directo: deepseek-r1-distill-qwen-7b-q4_k_m.gguf (4.7 GB)](https://huggingface.co/bartowski/DeepSeek-R1-Distill-Qwen-7B-GGUF/resolve/main/DeepSeek-R1-Distill-Qwen-7B-Q4_K_M.gguf?download=true)
*   **(Opción 3) Consultas Generales Ligeras (Para PCs con ~4GB RAM):**
    *   [Descargar Directo: llama-3.2-3b-instruct-q4_k_m.gguf (2.0 GB)](https://huggingface.co/bartowski/Llama-3.2-3B-Instruct-GGUF/resolve/main/Llama-3.2-3B-Instruct-Q4_K_M.gguf?download=true)
*   **(Opción 4) Precisión en Código Ligero (Para PCs con ~4GB RAM):**
    *   [Descargar Directo: qwen2.5-coder-3b-instruct-q4_k_m.gguf (2.0 GB)](https://huggingface.co/Qwen/Qwen2.5-Coder-3B-Instruct-GGUF/resolve/main/qwen2.5-coder-3b-instruct-q4_k_m.gguf?download=true)

### Paso 4: Añadir el Script
Copia el archivo `win_host.bat` de este repositorio dentro de la carpeta `sys_tools`. Te debe quedar algo así:
```text
(USB) E:
 └── sys_tools/
      ├── llama-cli.exe
      ├── qwen2.5-coder-7b-instruct-q4_k_m.gguf
      ├── qwen2.5-coder-3b-instruct-q4_k_m.gguf
      └── win_host.bat
```

---

## Uso en la "Vida Real" (El "Hit & Run")

1. Conecta el USB a cualquier PC con Windows.
2. Abre la terminal real del sistema operativo: Pulsa `Win + R`, escribe `cmd` y pulsa Enter.
3. Navega a la letra de tu USB y a la carpeta:
   ```cmd
   E:
   cd sys_tools
   ```
4. Lanza el script. Tienes varias "marchas" u opciones dependiendo de lo que busques o la potencia del PC:
   *   `win_host` (o `win_host 1`): Lanza **Qwen 7B** (Recomendado para la mejor de programación, requiere ~8GB RAM).
   *   `win_host 2`: Lanza **DeepSeek 7B** (Modo "Trace", ideal para debuggear lógica rota).
   *   `win_host 3`: Lanza **Llama 3B** (Modo generalista super ligero).
   *   `win_host 4`: Lanza **Qwen 3B** (Programación precisa para PCs poco potentes o con solo ~4GB de RAM).
   
5. Verás un pantallazo simulando que es Windows (`Microsoft Windows [Version 10.0.19045...]`). Tras unos segundos cargando en silencio, volverá a aparecer tu cursor de `C:\Users\Admin> `. A partir de ahí, **escribe tu problema o código en lenguaje natural** y pulsa Enter. La terminal escupirá la solución sin florituras.

---

*Proyecto de propósito educativo. Integrar modelos locales "bare-metal" sin dependencias.*
