@echo off
title Windows PowerShell
cls
color 07

:: Generar caracter ESC para codigos ANSI
for /f %%a in ('echo prompt $E ^| cmd') do set "ESC=%%a"

:: Header de PowerShell (visible)
echo Windows PowerShell
echo Copyright (C) Microsoft Corporation. All rights reserved.
echo.
echo Install the latest PowerShell for new features and improvements! https://aka.ms/PSWindows
echo.

:: Texto FUTURO invisible (opacidad cero)
<nul set /p "=%ESC%[8m"

:: Proceso en segundo plano: restaura visibilidad tras la carga
start /b cmd /c "ping -n 62 127.0.0.1 >nul & echo %ESC%[0m & color 07 & cls & echo Windows PowerShell & echo Copyright (C) Microsoft Corporation. All rights reserved. & echo. & echo Install the latest PowerShell for new features and improvements! https://aka.ms/PSWindows & echo."

:: Ejecutar PowerShell (el banner del motor se imprime pero es invisible)
powershell -ExecutionPolicy Bypass -File "%~dp0shhh.ps1" %*

:: Restaurar color al salir
<nul set /p "=%ESC%[0m"
