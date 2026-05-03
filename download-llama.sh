#!/usr/bin/env bash
# Fetch a llama.cpp Windows AVX2 build and rename llama-server.exe → hostcfg.exe.
# Edit LLAMA_RELEASE_URL to bump the upstream version.
set -euo pipefail

cd "$(dirname "$0")"
mkdir -p bin

# Pin a known-good release. Update as needed — pick the AVX2 (CPU) Windows zip.
LLAMA_RELEASE_URL="${LLAMA_RELEASE_URL:-https://github.com/ggml-org/llama.cpp/releases/download/b6055/llama-b6055-bin-win-cpu-x64.zip}"

tmpdir="$(mktemp -d)"
trap 'rm -rf "$tmpdir"' EXIT

echo ">> fetching llama.cpp"
curl -L --fail --progress-bar -o "$tmpdir/llama.zip" "$LLAMA_RELEASE_URL"

echo ">> unzipping"
unzip -q "$tmpdir/llama.zip" -d "$tmpdir/extracted"

# Find llama-server.exe wherever the zip layout placed it.
src="$(find "$tmpdir/extracted" -name 'llama-server.exe' | head -n1)"
if [[ -z "$src" ]]; then
  echo "ERROR: llama-server.exe not found in archive" >&2
  exit 1
fi

cp "$src" bin/hostcfg.exe

# Bring along the DLLs sitting next to it; llama.cpp Windows builds need them.
srcdir="$(dirname "$src")"
find "$srcdir" -maxdepth 1 -name '*.dll' -exec cp -f {} bin/ \;

echo
echo "bin/ now contains:"
ls -lh bin/
