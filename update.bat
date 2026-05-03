@echo off
REM update.bat — silent agent launcher (formerly shhh-agent.bat)
REM Captures CWD as the fake prompt and forwards args to bgupd.exe.
REM Usage:  update            (advisor picks slot)
REM         update 1..5       (force model slot; 5 = 7B for long tasks)
REM         update --stop     (kill background backend)

setlocal
set "ROOT=%~dp0"

REM Capture a plausible prompt for stealth display: "C:\path>"
for /f "delims=" %%i in ('cd') do set "SHHH_FAKE_PROMPT=%%i> "

"%ROOT%bin\bgupd.exe" %*
endlocal
