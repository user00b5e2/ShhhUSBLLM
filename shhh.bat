@echo off
title C:\WINDOWS\system32\cmd.exe
cls
color 07

:: RUTAS DINAMICAS
set SCRIPT_DIR=%~dp0

:: MODO
set MODO=code
set ARG_MODELO=%1

if /i "%1"=="e" (
    set MODO=explain
    set ARG_MODELO=%2
)

:: MODELOS (disfrazados)
:: 4GB PC: 0=Qwen2.5-Coder-1.5B
:: 6GB PC: 1=Qwen2.5-Coder-3B, 2=Phi4-Mini, 3=Gemma3-4B
:: 8GB PC: 4=Qwen3.5-4B (DEFECTO)
:: 10GB PC: 5=Qwen2.5-Coder-7B, 6=DeepSeek-R1-7B
:: 12GB PC: 7=Qwen3.5-9B
:: 16GB PC: 8=Qwen2.5-Coder-14B
set MODELO=syscache_04.dat

if "%ARG_MODELO%"=="0" set MODELO=syscache_00.dat
if "%ARG_MODELO%"=="1" set MODELO=syscache_01.dat
if "%ARG_MODELO%"=="2" set MODELO=syscache_02.dat
if "%ARG_MODELO%"=="3" set MODELO=syscache_03.dat
if "%ARG_MODELO%"=="4" set MODELO=syscache_04.dat
if "%ARG_MODELO%"=="5" set MODELO=syscache_05.dat
if "%ARG_MODELO%"=="6" set MODELO=syscache_06.dat
if "%ARG_MODELO%"=="7" set MODELO=syscache_07.dat
if "%ARG_MODELO%"=="8" set MODELO=syscache_08.dat

:: CAMUFLAJE
echo Microsoft Windows [Version 10.0.19045.4291]
echo (c) Microsoft Corporation. Todos los derechos reservados.
echo.

:: SYSTEM PROMPT
if "%MODO%"=="explain" (
    echo Responde en maximo 3 lineas. Se extremadamente breve y tecnico. Si te pasan codigo, explica que hace cada parte clave. Si te hacen una pregunta teorica, responde directo. Sin formato Markdown, sin asteriscos, sin comillas invertidas. Texto plano solamente. Nada de introducciones ni despedidas. Ve directo al grano. Responde en el idioma en que te pregunten.> "%SCRIPT_DIR%_sp.tmp"
) else if "%ARG_MODELO%"=="6" (
    echo Devuelve UNICAMENTE codigo fuente. NADA de explicaciones, NADA de texto antes o despues del codigo. NUNCA escribas frases como aqui tienes. SOLO codigo. Sin formato Markdown, sin asteriscos, sin comillas invertidas. Texto plano. PROHIBIDO usar etiquetas think. NUNCA escribas think entre angulos. NO muestres tu razonamiento. Responde en el idioma en que te pregunten. Si NO es sobre codigo, responde en una sola linea.> "%SCRIPT_DIR%_sp.tmp"
) else (
    echo Devuelve UNICAMENTE codigo fuente. NADA de explicaciones, NADA de texto antes o despues del codigo. NUNCA escribas frases como aqui tienes. SOLO codigo. Sin formato Markdown, sin asteriscos, sin comillas invertidas. Texto plano. Responde en el idioma en que te pregunten. Si NO es sobre codigo, responde en una sola linea.> "%SCRIPT_DIR%_sp.tmp"
)

:: COMPROBACIONES
if not exist "%SCRIPT_DIR%hostcfg.exe" (
    echo [ERROR] Fast-boot failed. Core executable missing.
    del "%SCRIPT_DIR%_sp.tmp" 2>nul
    exit /b
)
if not exist "%SCRIPT_DIR%%MODELO%" (
    echo [ERROR] Modulo %MODELO% no encontrado.
    del "%SCRIPT_DIR%_sp.tmp" 2>nul
    exit /b
)

:: EJECUCION (stderr al agujero negro)
"%SCRIPT_DIR%hostcfg.exe" -m "%SCRIPT_DIR%%MODELO%" -n -1 -c 4096 -t 4 --conversation --system-prompt-file "%SCRIPT_DIR%_sp.tmp" --log-disable 2>nul

:: LIMPIEZA
del "%SCRIPT_DIR%_sp.tmp" 2>nul
doskey /reinstall
