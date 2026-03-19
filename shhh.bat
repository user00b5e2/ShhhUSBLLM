@echo off
title C:\WINDOWS\system32\cmd.exe
cls
color 07

:: 1. MODO: CODIGO o EXPLICACION
set MODO=code
set ARG_MODELO=%1

if /i "%1"=="e" (
    set MODO=explain
    set ARG_MODELO=%2
)

:: 2. DEFINIR MODELO
set MODELO=qwen2.5-coder-3b-instruct-q4_k_m.gguf

if "%ARG_MODELO%"=="1" set MODELO=qwen2.5-coder-3b-instruct-q4_k_m.gguf
if "%ARG_MODELO%"=="2" set MODELO=Qwen2.5-Coder-7B-Instruct-Q4_K_M.gguf
if "%ARG_MODELO%"=="3" set MODELO=deepseek-r1-distill-qwen-7b-q4_k_m.gguf
if "%ARG_MODELO%"=="4" set MODELO=Phi-4-mini-instruct-Q4_K_M.gguf
if "%ARG_MODELO%"=="5" set MODELO=gemma-3-4b-it-Q4_K_M.gguf

:: 3. CAMUFLAJE
echo Microsoft Windows [Version 10.0.19045.4291]
echo (c) Microsoft Corporation. Todos los derechos reservados.
echo.

:: 4. SYSTEM PROMPT segun modo
if "%MODO%"=="explain" (
    echo Responde en maximo 3 lineas. Se extremadamente breve y tecnico. Si te pasan codigo, explica que hace cada parte clave. Si te hacen una pregunta teorica, responde directo. Sin formato Markdown, sin asteriscos, sin comillas invertidas. Texto plano solamente. Nada de introducciones ni despedidas. Ve directo al grano. Responde en el idioma en que te pregunten.> _sys_prompt.txt
) else if "%ARG_MODELO%"=="3" (
    echo Devuelve UNICAMENTE codigo fuente. NADA de explicaciones, NADA de comentarios, NADA de texto antes o despues del codigo. Si te piden un programa, devuelve solo el codigo. Si te piden corregir, devuelve solo el codigo corregido. NUNCA escribas frases como aqui tienes o este es el codigo. SOLO codigo. Sin formato Markdown, sin comillas invertidas, sin asteriscos. Texto plano solamente. PROHIBIDO usar etiquetas think. NUNCA escribas think entre angulos. NO muestres tu razonamiento interno. Responde en el idioma en que te pregunten. Si la pregunta NO es sobre codigo, responde en una sola linea tecnica.> _sys_prompt.txt
) else (
    echo Devuelve UNICAMENTE codigo fuente. NADA de explicaciones, NADA de comentarios, NADA de texto antes o despues del codigo. Si te piden un programa, devuelve solo el codigo. Si te piden corregir, devuelve solo el codigo corregido. NUNCA escribas frases como aqui tienes o este es el codigo. SOLO codigo. Sin formato Markdown, sin comillas invertidas, sin asteriscos. Texto plano solamente. Responde en el idioma en que te pregunten. Si la pregunta NO es sobre codigo, responde en una sola linea tecnica.> _sys_prompt.txt
)

:: 5. COMPROBACIONES
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

:: 6. EJECUCION
llama-cli.exe -m "%MODELO%" -n -1 -c 4096 -t 4 --conversation --system-prompt-file _sys_prompt.txt --log-disable 2> debug_log.txt

:: 7. LIMPIEZA
del _sys_prompt.txt 2>nul

if %ERRORLEVEL% NEQ 0 (
    echo.
    echo [SYS_ERROR] Process terminated unexpectedly.
    type debug_log.txt
    echo.
)
