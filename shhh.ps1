# PowerShell Stealth Launcher
# Uso: .\shhh.ps1 [1|2|3|4|5]

$host.UI.RawUI.WindowTitle = "Windows PowerShell"
Clear-Host

# 1. MODELOS
$MODELO = "qwen2.5-coder-3b-instruct-q4_k_m.gguf"

switch ($args[0]) {
    "1" { $MODELO = "qwen2.5-coder-3b-instruct-q4_k_m.gguf" }
    "2" { $MODELO = "Qwen2.5-Coder-7B-Instruct-Q4_K_M.gguf" }
    "3" { $MODELO = "deepseek-r1-distill-qwen-7b-q4_k_m.gguf" }
    "4" { $MODELO = "Phi-4-mini-instruct-Q4_K_M.gguf" }
    "5" { $MODELO = "gemma-3-4b-it-Q4_K_M.gguf" }
}

# 2. CAMUFLAJE
Write-Host "Windows PowerShell" -ForegroundColor White
Write-Host "Copyright (C) Microsoft Corporation. All rights reserved." -ForegroundColor White
Write-Host ""
Write-Host "Install the latest PowerShell for new features and improvements! https://aka.ms/PSWindows" -ForegroundColor White
Write-Host ""

# 3. SYSTEM PROMPT (guardado en archivo temporal)
if ($args[0] -eq "3") {
    $SYS = "Eres un subsistema de linea de comandos de Windows. No eres un asistente de IA. No uses saludos, no te despidas, no uses formato Markdown, ni des explicaciones innecesarias. PROHIBIDO usar etiquetas think. NUNCA escribas las palabras think entre angulos. NO muestres tu proceso de razonamiento interno. Responde DIRECTAMENTE con el resultado. Si el usuario introduce texto, devuelve unicamente el codigo resultante o la salida tecnica esperada en formato de texto plano. Cero charla. Responde siempre en el idioma en el que te pregunten."
} else {
    $SYS = "Eres un subsistema de linea de comandos de Windows. No eres un asistente de IA. No uses saludos, no te despidas, no uses formato Markdown, ni des explicaciones innecesarias. Si el usuario introduce texto, devuelve unicamente el codigo resultante o la salida tecnica esperada en formato de texto plano. Cero charla. Responde siempre en el idioma en el que te pregunten."
}
[System.IO.File]::WriteAllText("_sys_prompt.txt", $SYS)

# 4. COMPROBACIONES
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

# 5. EJECUCION
# --conversation = modo chat (gestiona turnos automaticamente)
# --system-prompt-file = carga system prompt desde archivo
# -t 4 = 4 hilos CPU
# -c 4096 = contexto
& .\llama-cli.exe -m $MODELO -n -1 -c 4096 -t 4 --conversation --system-prompt-file _sys_prompt.txt --log-disable 2> debug_log.txt

# 6. LIMPIEZA
Remove-Item "_sys_prompt.txt" -ErrorAction SilentlyContinue

if ($LASTEXITCODE -ne 0) {
    Write-Host ""
    Write-Host "[SYS_ERROR] Process terminated unexpectedly." -ForegroundColor Red
    Get-Content debug_log.txt
    Write-Host ""
}
