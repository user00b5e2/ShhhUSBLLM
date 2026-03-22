# PowerShell Stealth Launcher (la invisibilidad la maneja shhhps.bat)

# MODO: code (defecto), explain (e), think (t)
$MODO = "code"
$ARG_MODELO = $args[0]

if ($args[0] -eq "e") {
    $MODO = "explain"
    $ARG_MODELO = $args[1]
}
if ($args[0] -eq "t") {
    $MODO = "think"
    $ARG_MODELO = $args[1]
}

# MODELOS
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

$SCRIPT_DIR = $PSScriptRoot

# SYSTEM PROMPT + REASONING
$REASONING = "--reasoning-budget", "0"

if ($MODO -eq "explain") {
    $SP = "Responde en maximo 3 lineas. Breve y tecnico. Explica cada parte clave del codigo. Texto plano sin Markdown sin asteriscos. Responde en el idioma en que te pregunten."
} elseif ($MODO -eq "think") {
    $SP = "Muestra tu razonamiento paso a paso y luego devuelve UNICAMENTE el codigo final. Sin Markdown sin asteriscos. Texto plano. Responde en el idioma en que te pregunten."
    $REASONING = @()
} elseif ($ARG_MODELO -eq "6") {
    $SP = "Devuelve UNICAMENTE codigo. NADA de texto extra. Sin Markdown sin asteriscos. Texto plano. Responde en el idioma en que te pregunten. Si NO es codigo responde en una linea."
} else {
    $SP = "Devuelve UNICAMENTE codigo. NADA de texto extra. Sin Markdown sin asteriscos. Texto plano. Responde en el idioma en que te pregunten. Si NO es codigo responde en una linea."
}

# COMPROBACIONES
$EXE = Join-Path $SCRIPT_DIR "hostcfg.exe"
$MODEL_FILE = Join-Path $SCRIPT_DIR $MODELO

if (-not (Test-Path $EXE)) {
    Write-Host "[!] hostcfg.exe not found." -ForegroundColor Red
    exit
}
if (-not (Test-Path $MODEL_FILE)) {
    Write-Host "[!] $MODELO not found." -ForegroundColor Red
    exit
}

# INVISIBILIDAD: PS resetea ANSI al arrancar, asi que ponemos texto invisible aqui
# El shhhps.bat restaura la visibilidad con su proceso en segundo plano
[Console]::Write("$([char]27)[8m")

# EJECUCION DIRECTA EN TERMINAL
$PROMPT = "PS $PWD> "
& $EXE -m $MODEL_FILE -c 2048 -t 8 -cnv --simple-io --color off -r $PROMPT -sys $SP @REASONING --no-show-timings --log-disable 2>$null

# LIMPIEZA
Clear-History -ErrorAction SilentlyContinue
try { Remove-Item (Get-PSReadlineOption).HistorySavePath -ErrorAction SilentlyContinue } catch {}
