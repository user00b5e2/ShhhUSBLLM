#!/usr/bin/env bash
# Cross-compile the Windows build (bgupd.exe) for the USB target.
# Run from any platform with Go installed.
set -euo pipefail

cd "$(dirname "$0")/shhh-agent"

mkdir -p ../bin

echo ">> windows/amd64 (USB target → bgupd.exe)"
GOOS=windows GOARCH=amd64 \
  go build -trimpath -ldflags="-s -w" -o ../bin/bgupd.exe .

echo
echo "Built:"
ls -lh ../bin/bgupd.exe
