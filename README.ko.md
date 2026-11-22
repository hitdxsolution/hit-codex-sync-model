# hit-codex-sync-model

[English](README.md) · 한국어

Codex 채팅의 모델, 추론 수준, 답변 속도를 일괄 변경하는 Go CLI이자 Codex 스킬입니다. 현재 PC와 계정에 저장된 채팅 중 최근 30일 내 사용한 채팅을 기본 대상으로 삼습니다. 변경 전에 설정을 백업하고, 적용 직후 다시 읽어 확인합니다.

## 설치

1. [GitHub Releases](https://github.com/hitdxsolution/hit-codex-sync-model/releases/latest)에서 운영체제와 CPU에 맞는 ZIP을 받습니다. Apple Silicon과 Windows on Arm은 `arm64`, Intel·AMD x64는 `amd64`를 선택합니다. macOS 파일 이름은 `darwin`으로 시작합니다.
2. 압축을 풉니다. Codex에서 `$hit-codex-sync-model`로 사용하려면 압축 안의 `hit-codex-sync-model` 폴더 전체를 해당 PC의 Codex 스킬 폴더에 넣습니다. 기본 위치는 macOS·Linux의 `~/.codex/skills/`, Windows의 `%USERPROFILE%\.codex\skills\`입니다.
3. Codex App Server와 채팅 설정 API를 지원하는 Codex CLI가 필요합니다. CLI는 `PATH`에서 찾습니다. macOS에서는 Codex·ChatGPT 앱에 포함된 CLI도 검색합니다. 별도 설치 경로는 `--codex PATH`로 지정할 수 있습니다.

릴리스의 실행 파일에는 Go 설치가 필요 없습니다. 소스에서 빌드하려면 Go 1.26 이상이 필요합니다. [운영체제별 빌드 절차](SKILL.md#설치와-실행-환경)를 확인하세요.

## 빠른 시작

터미널에서 압축을 푼 `hit-codex-sync-model` 폴더로 이동해 실행합니다.

```sh
./bin/hit-codex-sync-model --models
./bin/hit-codex-sync-model --model gpt-6.2-sol --effort medium --service-tier fast
./bin/hit-codex-sync-model --model gpt-6.2-sol --effort medium --service-tier fast --apply
```

Windows PowerShell에서는 실행 파일 이름을 `.\bin\hit-codex-sync-model.exe`로 바꿉니다. `--apply`가 없는 명령은 미리보기입니다. 지원 모델과 속도는 현재 계정의 `--models` 결과에서 확인하세요.

스킬로 실행할 때는 다음처럼 요청합니다.

> $hit-codex-sync-model 최근 30일 채팅을 gpt-6.2-sol, medium, Fast로 바꿔줘.

| 옵션 | 동작 |
|---|---|
| `--service-tier default` | 기본 답변 속도로 설정 |
| `--service-tier fast` | Fast로 설정 |
| `--service-tier ultrafast` | Ultrafast로 설정 |
| 속도 옵션 생략 | 각 채팅의 기존 답변 속도 유지 |
| `--all` | 전체 기간의 채팅 선택 |
| `--thread ID` | 특정 채팅 선택. 여러 번 사용 가능 |
| `--workers 8` | 일괄 작업의 동시 처리 수 설정. 기본 4 |

## 재개와 복구

결과 JSON은 기본적으로 `~/.codex/model-sync-backups/`에 저장합니다. `CODEX_HOME`을 지정했다면 그 아래에서 찾습니다. 실패한 작업을 이어서 실행하거나 변경 전 설정을 복구할 때 명령 출력에 표시된 경로를 사용합니다.

```sh
./bin/hit-codex-sync-model --resume "/path/to/report.json" --apply
./bin/hit-codex-sync-model --restore "/path/to/report.json"
./bin/hit-codex-sync-model --restore "/path/to/report.json" --apply
```

다른 Codex 프로세스가 쓰기 잠금을 보유하면 그 앱을 종료한 뒤 재개합니다. macOS에서는 외부 터미널에서 `--force-locked`를 쓰면 잠금이 남았을 때 앱을 종료하고, 재시도 후 성공하면 다시 엽니다. 실행 중인 채팅 작업도 중단됩니다. Windows·Linux에서는 잠금 보유 앱을 직접 종료합니다.

도구는 Codex App Server에 로컬 요청을 보내며 모델 답변 생성을 시작하지 않습니다. 결과 JSON에는 채팅 ID와 변경 전후 설정이 들어갑니다. 다른 PC에서도 그 PC에서 접근할 수 있는 채팅과 계정 설정을 사용합니다.

자세한 동작은 [스킬 안내](SKILL.md)를 참고하세요.
