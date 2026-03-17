@echo off
title Símbolo del sistema
cls

:: 1. DEFINIR VARIABLES Y MODELOS
:: Por defecto (win_host) carga Qwen 7B (Requiere ~8GB RAM)
set MODELO=qwen2.5-coder-7b-instruct-q4_k_m.gguf
set MODO_INFO=CLI_Standard_8GB

:: Selector oculto: Dependiendo del numero, cambia el modelo
if "%1"=="1" (
    :: Qwen 7B (8GB RAM) - Programación pesada
    set MODELO=qwen2.5-coder-7b-instruct-q4_k_m.gguf
    set MODO_INFO=CLI_Standard_8GB
)
if "%1"=="2" (
    :: DeepSeek R1 7B (8GB RAM) - Razonamiento profundo y bugs lógicos
    set MODELO=deepseek-r1-distill-qwen-7b-q4_k_m.gguf
    set MODO_INFO=Verbose_Trace_Mode
)
if "%1"=="3" (
    :: Llama 3.2 3B (4GB RAM) - Consultas generales ligeras
    set MODELO=llama-3.2-3b-instruct-q4_k_m.gguf
    set MODO_INFO=Sys_Core
)
if "%1"=="4" (
    :: Qwen 3B (Requiere ~4GB RAM) - Programación ligera / PCs poco potentes
    set MODELO=qwen2.5-coder-3b-instruct-q4_k_m.gguf
    set MODO_INFO=CLI_Lite_4GB
)

:: 2. EL CAMUFLAJE VISUAL (Falsa terminal)
echo Microsoft Windows [Version 10.0.19045.4291]
echo (c) Microsoft Corporation. Todos los derechos reservados.
echo.
echo [System check OK. Modulo activo: %MODO_INFO%]
echo.

:: 3. EL PROMPT DEL SISTEMA MUERTO
set SYSTEM_PROMPT="Eres un subsistema de linea de comandos de Windows. No eres un asistente de IA. No uses saludos, no te despidas, no uses formato Markdown (sin asteriscos ni comillas invertidas), ni des explicaciones. Si el usuario introduce texto, devuelve unicamente el codigo resultante o la salida tecnica esperada en formato de texto plano. Cero charla."

:: MAIN CHECK
if not exist "llama-cli.exe" (
    echo [ERROR] Fast-boot failed. Core executable missing.
    exit /b
)
if not exist "%MODELO%" (
    echo [ERROR] Modulo %MODELO% no encontrado. Por favor revisa que el archivo exista en la carpeta.
    exit /b
)

:: 4. EJECUCION SILENCIOSA DE LA IA
llama-cli.exe -m "%MODELO%" -n -1 -c 4096 -i --system %SYSTEM_PROMPT% --reverse-prompt "C:\Users\Admin> " --in-prefix "" -p "C:\Users\Admin> " --log-disable 2>NUL
