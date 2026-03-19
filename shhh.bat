@echo off
title C:\WINDOWS\system32\cmd.exe
cls
color 07

:: 1. DEFINIR VARIABLES Y MODELOS
set MODELO=qwen2.5-coder-3b-instruct-q4_k_m.gguf

if "%1"=="1" set MODELO=qwen2.5-coder-3b-instruct-q4_k_m.gguf
if "%1"=="2" set MODELO=Qwen2.5-Coder-7B-Instruct-Q4_K_M.gguf
if "%1"=="3" set MODELO=deepseek-r1-distill-qwen-7b-q4_k_m.gguf
if "%1"=="4" set MODELO=Phi-4-mini-instruct-Q4_K_M.gguf
if "%1"=="5" set MODELO=gemma-3-4b-it-Q4_K_M.gguf

:: 2. CAMUFLAJE
echo Microsoft Windows [Version 10.0.19045.4291]
echo (c) Microsoft Corporation. Todos los derechos reservados.
echo.

:: 3. SYSTEM PROMPT (guardado en archivo temporal)
if "%1"=="3" (
    echo Eres un subsistema de linea de comandos de Windows. No eres un asistente de IA. No uses saludos, no te despidas, no uses formato Markdown, ni des explicaciones innecesarias. PROHIBIDO usar etiquetas think. NUNCA escribas las palabras think entre angulos. NO muestres tu proceso de razonamiento interno. Responde DIRECTAMENTE con el resultado. Si el usuario introduce texto, devuelve unicamente el codigo resultante o la salida tecnica esperada en formato de texto plano. Cero charla. Responde siempre en el idioma en el que te pregunten.> _sys_prompt.txt
) else (
    echo Eres un subsistema de linea de comandos de Windows. No eres un asistente de IA. No uses saludos, no te despidas, no uses formato Markdown, ni des explicaciones innecesarias. Si el usuario introduce texto, devuelve unicamente el codigo resultante o la salida tecnica esperada en formato de texto plano. Cero charla. Responde siempre en el idioma en el que te pregunten.> _sys_prompt.txt
)

:: 4. COMPROBACIONES
if not exist "llama-cli.exe" (
    echo [ERROR] Fast-boot failed. Core executable missing.
    del _sys_prompt.txt 2>nul
    exit /b
)
if not exist "%MODELO%" (
    echo [ERROR] Modulo %MODELO% no encontrado.
    del _sys_prompt.txt 2>nul
    exit /b
)

:: 5. EJECUCION
:: --conversation = modo chat (gestiona turnos automaticamente)
:: --system-prompt-file = carga system prompt desde archivo
:: -t 4 = 4 hilos CPU
:: -c 4096 = contexto
llama-cli.exe -m "%MODELO%" -n -1 -c 4096 -t 4 --conversation --system-prompt-file _sys_prompt.txt --log-disable 2> debug_log.txt

:: 6. LIMPIEZA
del _sys_prompt.txt 2>nul

if %ERRORLEVEL% NEQ 0 (
    echo.
    echo [SYS_ERROR] Process terminated unexpectedly.
    type debug_log.txt
    echo.
)
