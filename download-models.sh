#!/usr/bin/env bash
# Download the GGUF models that shhh-agent expects in models/.
# Run from any platform (curl required). Idempotent: skips existing files.
set -euo pipefail

cd "$(dirname "$0")"
mkdir -p models

declare -a MODELS=(
  # slot:filename:url
  # Slots 1 and 4 share Qwen3-1.7B (dual-mode, /no_think directive); only one download needed.
  "1:qwen3-1.7b-q4_k_m.gguf:https://huggingface.co/unsloth/Qwen3-1.7B-GGUF/resolve/main/Qwen3-1.7B-Q4_K_M.gguf"
  # Slots 2 and 3 share Qwen3-4B-Instruct-2507 (non-thinking native).
  "2:qwen3-4b-instruct-2507-q4_k_m.gguf:https://huggingface.co/unsloth/Qwen3-4B-Instruct-2507-GGUF/resolve/main/Qwen3-4B-Instruct-2507-Q4_K_M.gguf"
  # Slot 5 (8B) is OPT-IN — adds ~5.0 GB and is only viable on 8 GB RAM if you
  # close Chrome / heavy apps. Enable by setting WITH_LARGE=1.
  "${WITH_LARGE:+5:qwen3-8b-q4_k_m.gguf:https://huggingface.co/unsloth/Qwen3-8B-GGUF/resolve/main/Qwen3-8B-Q4_K_M.gguf}"
)

for entry in "${MODELS[@]}"; do
  [[ -z "$entry" ]] && continue
  IFS=":" read -r slot file url <<< "$entry"
  out="models/$file"
  if [[ -f "$out" ]]; then
    echo "[skip slot $slot] $file"
    continue
  fi
  echo "[slot $slot] downloading $file"
  curl -L --fail --progress-bar -o "$out.part" "$url"
  mv "$out.part" "$out"
done

echo
echo "Models in models/:"
ls -lh models/
