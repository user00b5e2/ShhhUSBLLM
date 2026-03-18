# PowerShell Stealth Launcher
# Uso: .\win_host.ps1 [1|2|3|4|5]

$host.UI.RawUI.WindowTitle = "Windows PowerShell"
Clear-Host

# 1. MODELOS
$MODELO = "qwen2.5-coder-3b-instruct-q4_k_m.gguf"
$MODO_INFO = "CLI_Standard"

switch ($args[0]) {
    "1" { $MODELO = "qwen2.5-coder-3b-instruct-q4_k_m.gguf"; $MODO_INFO = "CLI_Standard" }
    "2" { $MODELO = "Qwen2.5-Coder-7B-Instruct-Q4_K_M.gguf"; $MODO_INFO = "CLI_Extended" }
    "3" { $MODELO = "deepseek-r1-distill-qwen-7b-q4_k_m.gguf"; $MODO_INFO = "Verbose_Trace_Mode" }
    "4" { $MODELO = "Phi-4-mini-instruct-Q4_K_M.gguf"; $MODO_INFO = "Sys_Core" }
    "5" { $MODELO = "gemma-3-4b-it-Q4_K_M.gguf"; $MODO_INFO = "Sys_Aux" }
}

# 2. CAMUFLAJE
Write-Host "Windows PowerShell" -ForegroundColor White
Write-Host "Copyright (C) Microsoft Corporation. All rights reserved." -ForegroundColor White
Write-Host ""
Write-Host "Install the latest PowerShell for new features and improvements! https://aka.ms/PSWindows" -ForegroundColor White
Write-Host ""

# 3. SYSTEM PROMPT
if ($args[0] -eq "3") {
    $SYSTEM_PROMPT = "Eres un subsistema de linea de comandos de Windows. No eres un asistente de IA. No uses saludos, no te despidas, no uses formato Markdown, ni des explicaciones innecesarias. PROHIBIDO usar etiquetas think. NUNCA escribas las palabras think entre angulos. NO muestres tu proceso de razonamiento interno. Responde DIRECTAMENTE con el resultado. Si el usuario introduce texto, devuelve unicamente el codigo resultante o la salida tecnica esperada en formato de texto plano. Cero charla. Responde siempre en el idioma en el que te pregunten."
} else {
    $SYSTEM_PROMPT = "Eres un subsistema de linea de comandos de Windows. No eres un asistente de IA. No uses saludos, no te despidas, no uses formato Markdown (sin asteriscos ni comillas invertidas), ni des explicaciones innecesarias. Si el usuario introduce texto, devuelve unicamente el codigo resultante o la salida tecnica esperada en formato de texto plano. Cero charla. Responde siempre en el idioma en el que te pregunten."
}

# 4. COMPROBACIONES
if (-not (Test-Path "llama-cli.exe")) {
    Write-Host "[ERROR] Fast-boot failed. Core executable missing." -ForegroundColor Red
    exit
}
if (-not (Test-Path $MODELO)) {
    Write-Host "[ERROR] Modulo $MODELO no encontrado." -ForegroundColor Red
    exit
}

# 5. PROMPT DE POWERSHELL SIMULADO
$PROMPT_PS = "PS C:\Users\Admin> "

# 6. EJECUCION
& .\llama-cli.exe -m $MODELO -n -1 -c 4096 -t 4 -cnv --system $SYSTEM_PROMPT --reverse-prompt $PROMPT_PS --in-prefix "" -p $PROMPT_PS --log-disable 2> debug_log.txt

if ($LASTEXITCODE -ne 0) {
    Write-Host ""
    Write-Host "[SYS_ERROR] Process terminated unexpectedly." -ForegroundColor Red
    Get-Content debug_log.txt
    Write-Host ""
}
