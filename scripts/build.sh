#!/bin/sh
set -eu
skill_dir=$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)
mkdir -p "$skill_dir/bin"
go -C "$skill_dir/app" build -trimpath -o "$skill_dir/bin/hit-codex-sync-model" .
