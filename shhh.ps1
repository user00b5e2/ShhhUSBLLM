# PowerShell Stealth Launcher

$host.UI.RawUI.WindowTitle = "Windows PowerShell"
Clear-Host

# MODO
$MODO = "code"
$ARG_MODELO = $args[0]

if ($args[0] -eq "e") {
    $MODO = "explain"
    $ARG_MODELO = $args[1]
}

# MODELOS (disfrazados)
$MODELO = "syscache_04.dat"

switch ($ARG_MODELO) {
    "0" { $MODELO = "syscache_00.dat" }
    "1" { $MODELO = "syscache_01.dat" }
    "2" { $MODELO = "syscache_02.dat" }
    "3" { $MODELO = "syscache_03.dat" }
    "4" { $MODELO = "syscache_04.dat" }
    "5" { $MODELO = "syscache_05.dat" }
    "6" { $MODELO = "syscache_06.dat" }
    "7" { $MODELO = "syscache_07.dat" }
    "8" { $MODELO = "syscache_08.dat" }
}

# RUTAS
$SCRIPT_DIR = $PSScriptRoot

# CAMUFLAJE
Write-Host "Windows PowerShell" -ForegroundColor White
Write-Host "Copyright (C) Microsoft Corporation. All rights reserved." -ForegroundColor White
Write-Host ""
Write-Host "Install the latest PowerShell for new features and improvements! https://aka.ms/PSWindows" -ForegroundColor White
Write-Host ""

# SYSTEM PROMPT
$SP_FILE = Join-Path $SCRIPT_DIR "_sp.tmp"

if ($MODO -eq "explain") {
    $SYS = "Responde en maximo 3 lineas. Se extremadamente breve y tecnico. Si te pasan codigo, explica que hace cada parte clave. Si te hacen una pregunta teorica, responde directo. Sin formato Markdown, sin asteriscos, sin comillas invertidas. Texto plano solamente. Nada de introducciones ni despedidas. Ve directo al grano. Responde en el idioma en que te pregunten."
} elseif ($ARG_MODELO -eq "6") {
    $SYS = "Devuelve UNICAMENTE codigo fuente. NADA de explicaciones, NADA de texto antes o despues del codigo. NUNCA escribas frases como aqui tienes. SOLO codigo. Sin formato Markdown, sin asteriscos, sin comillas invertidas. Texto plano. PROHIBIDO usar etiquetas think. NUNCA escribas think entre angulos. NO muestres tu razonamiento. Responde en el idioma en que te pregunten. Si NO es sobre codigo, responde en una sola linea."
} else {
    $SYS = "Devuelve UNICAMENTE codigo fuente. NADA de explicaciones, NADA de texto antes o despues del codigo. NUNCA escribas frases como aqui tienes. SOLO codigo. Sin formato Markdown, sin asteriscos, sin comillas invertidas. Texto plano. Responde en el idioma en que te pregunten. Si NO es sobre codigo, responde en una sola linea."
}
[System.IO.File]::WriteAllText($SP_FILE, $SYS)

# COMPROBACIONES
$EXE = Join-Path $SCRIPT_DIR "hostcfg.exe"
$MODEL_FILE = Join-Path $SCRIPT_DIR $MODELO

if (-not (Test-Path $EXE)) {
    Write-Host "[ERROR] Fast-boot failed. Core executable missing." -ForegroundColor Red
    Remove-Item $SP_FILE -ErrorAction SilentlyContinue
    exit
}
if (-not (Test-Path $MODEL_FILE)) {
    Write-Host "[ERROR] Modulo $MODELO no encontrado." -ForegroundColor Red
    Remove-Item $SP_FILE -ErrorAction SilentlyContinue
    exit
}

# EJECUCION
& $EXE -m $MODEL_FILE -n -1 -c 4096 -t 4 --conversation --system-prompt-file $SP_FILE --log-disable 2>$null

# LIMPIEZA
Remove-Item $SP_FILE -ErrorAction SilentlyContinue
Clear-History -ErrorAction SilentlyContinue
try { Remove-Item (Get-PSReadlineOption).HistorySavePath -ErrorAction SilentlyContinue } catch {}
