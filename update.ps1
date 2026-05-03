# update.ps1 — silent agent launcher (formerly shhh-agent.ps1)
# Captures the real PowerShell prompt and forwards args to bgupd.exe.
# Usage:  .\update.ps1            (advisor picks slot)
#         .\update.ps1 1..5       (force model slot; 5 = 7B for long tasks)
#         .\update.ps1 --stop     (kill background backend)

$Root = Split-Path -Parent $MyInvocation.MyCommand.Path

# Capture the user's actual prompt (oh-my-posh, posh-git, plain — whatever they have).
# `prompt` is a function in PS; calling it returns the rendered string.
try {
    $captured = (& { prompt }) -join ''
    if ($captured) { $env:SHHH_FAKE_PROMPT = $captured.TrimEnd() + ' ' }
} catch {
    # Fallback: classic PS prompt
    $env:SHHH_FAKE_PROMPT = "PS $((Get-Location).Path)> "
}

& (Join-Path $Root "bin\bgupd.exe") $args
