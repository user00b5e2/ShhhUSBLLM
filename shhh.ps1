# PowerShell Stealth Launcher
# Uso: .\shhh.ps1 [e] [1|2|3|4|5]

$host.UI.RawUI.WindowTitle = "Windows PowerShell"
Clear-Host

# 1. MODO: CODIGO o EXPLICACION
$MODO = "code"
$ARG_MODELO = $args[0]

if ($args[0] -eq "e") {
    $MODO = "explain"
    $ARG_MODELO = $args[1]
}

# 2. MODELOS
$MODELO = "qwen2.5-coder-3b-instruct-q4_k_m.gguf"

switch ($ARG_MODELO) {
    "1" { $MODELO = "qwen2.5-coder-3b-instruct-q4_k_m.gguf" }
    "2" { $MODELO = "Qwen2.5-Coder-7B-Instruct-Q4_K_M.gguf" }
    "3" { $MODELO = "deepseek-r1-distill-qwen-7b-q4_k_m.gguf" }
    "4" { $MODELO = "Phi-4-mini-instruct-Q4_K_M.gguf" }
    "5" { $MODELO = "gemma-3-4b-it-Q4_K_M.gguf" }
}

# 3. CAMUFLAJE
Write-Host "Windows PowerShell" -ForegroundColor White
Write-Host "Copyright (C) Microsoft Corporation. All rights reserved." -ForegroundColor White
Write-Host ""
Write-Host "Install the latest PowerShell for new features and improvements! https://aka.ms/PSWindows" -ForegroundColor White
Write-Host ""

# 4. SYSTEM PROMPT
if ($MODO -eq "explain") {
    $SYS = "Responde en maximo 3 lineas. Se extremadamente breve y tecnico. Si te pasan codigo, explica que hace cada parte clave. Si te hacen una pregunta teorica, responde directo. Sin formato Markdown, sin asteriscos, sin comillas invertidas. Texto plano solamente. Nada de introducciones ni despedidas. Ve directo al grano. Responde en el idioma en que te pregunten."
} elseif ($ARG_MODELO -eq "3") {
    $SYS = "Devuelve UNICAMENTE codigo fuente. NADA de explicaciones, NADA de comentarios, NADA de texto antes o despues del codigo. Si te piden un programa, devuelve solo el codigo. Si te piden corregir, devuelve solo el codigo corregido. NUNCA escribas frases como aqui tienes o este es el codigo. SOLO codigo. Sin formato Markdown, sin comillas invertidas, sin asteriscos. Texto plano solamente. PROHIBIDO usar etiquetas think. NUNCA escribas think entre angulos. NO muestres tu razonamiento interno. Responde en el idioma en que te pregunten. Si la pregunta NO es sobre codigo, responde en una sola linea tecnica."
} else {
    $SYS = "Devuelve UNICAMENTE codigo fuente. NADA de explicaciones, NADA de comentarios, NADA de texto antes o despues del codigo. Si te piden un programa, devuelve solo el codigo. Si te piden corregir, devuelve solo el codigo corregido. NUNCA escribas frases como aqui tienes o este es el codigo. SOLO codigo. Sin formato Markdown, sin comillas invertidas, sin asteriscos. Texto plano solamente. Responde en el idioma en que te pregunten. Si la pregunta NO es sobre codigo, responde en una sola linea tecnica."
}
[System.IO.File]::WriteAllText("_sys_prompt.txt", $SYS)

# 5. COMPROBACIONES
if (-not (Test-Path "llama-cli.exe")) {
    Write-Host "[ERROR] Fast-boot failed. Core executable missing." -ForegroundColor Red
    Remove-Item "_sys_prompt.txt" -ErrorAction SilentlyContinue
    exit
}
if (-not (Test-Path $MODELO)) {
    Write-Host "[ERROR] Modulo $MODELO no encontrado." -ForegroundColor Red
    Remove-Item "_sys_prompt.txt" -ErrorAction SilentlyContinue
    exit
}

# 6. EJECUCION
& .\llama-cli.exe -m $MODELO -n -1 -c 4096 -t 4 --conversation --system-prompt-file _sys_prompt.txt --log-disable 2> debug_log.txt

# 7. LIMPIEZA
Remove-Item "_sys_prompt.txt" -ErrorAction SilentlyContinue

if ($LASTEXITCODE -ne 0) {
    Write-Host ""
    Write-Host "[SYS_ERROR] Process terminated unexpectedly." -ForegroundColor Red
    Get-Content debug_log.txt
    Write-Host ""
}
