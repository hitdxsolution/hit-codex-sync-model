#!/usr/bin/env python3
"""Compatibility entrypoint for commands created before the Go CLI migration."""
import os
from pathlib import Path
import sys

name = "hit-codex-sync-model.exe" if os.name == "nt" else "hit-codex-sync-model"
binary = Path(__file__).resolve().parents[1] / "bin" / name
if not binary.is_file():
    sys.exit("Go 실행 파일이 없습니다. SKILL.md의 운영체제별 빌드 방법을 확인하세요.")
os.execv(str(binary), [str(binary), *sys.argv[1:]])
