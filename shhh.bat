@echo off
title C:\WINDOWS\system32\cmd.exe
cls
color 07

:: Generar caracter ESC para codigos ANSI
for /f %%a in ('echo prompt $E ^| cmd') do set "ESC=%%a"

:: RUTAS DINAMICAS
set SCRIPT_DIR=%~dp0

:: MODO: code (defecto), explain (e), think (t)
set MODO=code
set ARG_MODELO=%1

if /i "%1"=="e" (
    set MODO=explain
    set ARG_MODELO=%2
)
if /i "%1"=="t" (
    set MODO=think
    set ARG_MODELO=%2
)

:: MODELOS
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

:: SYSTEM PROMPT + REASONING
set "REASONING=--reasoning-budget 0"

if "%MODO%"=="explain" (
    set "SP=Responde en maximo 3 lineas. Breve y tecnico. Explica cada parte clave del codigo. Texto plano sin Markdown sin asteriscos. Responde en el idioma en que te pregunten."
) else if "%MODO%"=="think" (
    set "SP=Razona paso a paso antes de responder. Muestra tu razonamiento completo. Luego da la respuesta final. Texto plano sin Markdown sin asteriscos. Responde en el idioma en que te pregunten."
    set "REASONING="
) else if "%ARG_MODELO%"=="6" (
    set "SP=Devuelve UNICAMENTE codigo. NADA de texto extra. Sin Markdown sin asteriscos. Texto plano. Responde en el idioma en que te pregunten. Si NO es codigo responde en una linea."
) else (
    set "SP=Devuelve UNICAMENTE codigo. NADA de texto extra. Sin Markdown sin asteriscos. Texto plano. Responde en el idioma en que te pregunten. Si NO es codigo responde en una linea."
)

:: CAMUFLAJE VISUAL (se imprime en color normal)
echo Microsoft Windows [Version 10.0.19045.4291]
echo (c) Microsoft Corporation. Todos los derechos reservados.
echo.

:: COMPROBACIONES
if not exist "%SCRIPT_DIR%hostcfg.exe" (
    echo [!] hostcfg.exe not found.
    exit /b
)
if not exist "%SCRIPT_DIR%%MODELO%" (
    echo [!] %MODELO% not found.
    exit /b
)

:: === TRUCO DE INVISIBILIDAD ===
:: Poner texto FUTURO en negro (invisible sobre fondo negro)
:: El header de Microsoft YA impreso sigue visible
<nul set /p "=%ESC%[30m"

:: Proceso en segundo plano: espera y luego restaura la visibilidad
:: color 07 resetea TODA la consola a blanco, cls limpia, echo reimprime header
start /b cmd /c "ping -n 62 127.0.0.1 >nul & color 07 & cls & echo Microsoft Windows [Version 10.0.19045.4291] & echo (c) Microsoft Corporation. Todos los derechos reservados. & echo."

:: EJECUCION DIRECTA EN TERMINAL (banner en negro = invisible)
"%SCRIPT_DIR%hostcfg.exe" -m "%SCRIPT_DIR%%MODELO%" -c 2048 -t 8 -cnv --simple-io --color off -r "%CD%> " -sys "%SP%" %REASONING% --no-show-timings --log-disable 2>nul

:: RESTAURAR COLOR al salir
<nul set /p "=%ESC%[0m"
doskey /reinstall
