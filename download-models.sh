#!/usr/bin/env bash
# Download the GGUF models that shhh-agent expects in models/.
# Run from any platform (curl required). Idempotent: skips existing files.
set -euo pipefail

cd "$(dirname "$0")"
mkdir -p models

declare -a MODELS=(
  # slot:filename:url
  # Slots 1 and 4 share Qwen2.5-Coder-1.5B; only one download needed.
  "1:qwen2.5-coder-1.5b-instruct-q4_k_m.gguf:https://huggingface.co/Qwen/Qwen2.5-Coder-1.5B-Instruct-GGUF/resolve/main/qwen2.5-coder-1.5b-instruct-q4_k_m.gguf"
  # Slots 2 and 3 share Qwen2.5-Coder-3B (default coding agent + chat).
  "2:qwen2.5-coder-3b-instruct-q4_k_m.gguf:https://huggingface.co/Qwen/Qwen2.5-Coder-3B-Instruct-GGUF/resolve/main/qwen2.5-coder-3b-instruct-q4_k_m.gguf"
  # Slot 5 (7B) is OPT-IN — adds ~4.7 GB and is only viable on 8 GB RAM if you
  # close Chrome / heavy apps. Enable by setting WITH_LARGE=1.
  "${WITH_LARGE:+5:qwen2.5-coder-7b-instruct-q4_k_m.gguf:https://huggingface.co/Qwen/Qwen2.5-Coder-7B-Instruct-GGUF/resolve/main/qwen2.5-coder-7b-instruct-q4_k_m.gguf}"
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
