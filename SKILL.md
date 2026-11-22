---
name: hit-codex-sync-model
description: Go CLI로 최근 30일 이내 사용한 Codex 채팅의 모델·추론 수준·답변 속도(기본, Fast, Ultrafast)를 일괄 변경한다. 명시적으로 전체 기간을 요청하면 --all을 사용한다. 모든 채팅 모델 동기화, 모델 일괄 전환, 잠금 재시도, 중단 후 재개, 변경 기록 복구 요청에 사용한다.
---

# 채팅 모델 동기화

## 설치와 실행 환경

스킬 폴더를 각 PC의 Codex 스킬 폴더에 넣는다. macOS·Windows·Linux의 x64와 ARM64용 배포 압축에는 해당 환경에서 바로 실행할 바이너리가 들어 있다. Codex CLI는 이 PC의 계정과 채팅에 접근할 수 있어야 하며, 앱이 사용하는 App Server의 채팅 설정 API를 지원해야 한다.

- macOS·Linux 실행 파일은 `./bin/hit-codex-sync-model`이다.
- Windows PowerShell 실행 파일은 `.\bin\hit-codex-sync-model.exe`이다. 아래 명령의 실행 파일 경로만 이 경로로 바꾼다.
- Codex CLI는 `PATH`에서 먼저 찾는다. macOS에서는 사용자 또는 시스템 Applications 폴더의 Codex·ChatGPT 앱에 포함된 CLI도 찾는다. 별도 설치 경로를 사용한다면 `--codex`에 이 PC의 실제 실행 파일 경로를 전달한다.
- 채팅과 백업 위치는 현재 `CODEX_HOME`을 따른다. 설정하지 않았으면 현재 사용자의 홈 아래 `.codex`를 사용한다.

실행 파일이 없거나 소스가 더 최신이면 Go 1.26 이상으로 해당 환경에서 빌드한다. macOS·Linux에서는 `sh scripts/build.sh`를 실행한다. Windows에서는 스킬 폴더에서 다음 PowerShell 명령을 실행한다.

```powershell
New-Item -ItemType Directory -Path bin -Force | Out-Null
go -C app build -trimpath -o ../bin/hit-codex-sync-model.exe .
```

## 실행 절차

1. 이 `SKILL.md`가 있는 폴더를 작업 디렉터리로 지정하고 **`./bin/hit-codex-sync-model` Go 실행 파일을 셸에서 실행한다.** 모델 목록 조회, 대상 미리보기, 일괄 적용, 재개, 복구를 모두 같은 실행 파일로 처리한다. 명령 예시만 답하는 대신 요청된 작업을 직접 실행하고 결과를 확인한다.
2. 설치 환경에 맞는 실행 파일을 지정한다. 실행 파일이 없거나 소스가 더 최신이면 위 빌드 절차를 따른다.
3. 새 변경 요청에서는 사용자에게 받은 모델과 추론 수준을 `--model`과 `--effort`로 넘긴다. 목표 모델을 판단할 정보가 없으면 사용자에게 묻고, 모델 목록을 요청하면 `./bin/hit-codex-sync-model --models`를 실행한다. 아래 모델과 추론 수준은 예시다. 추론 수준을 생략하면 앱이 지원되는 기존 값을 유지하고, 지원하지 않는 값에는 대상 모델의 기본값을 적용한다.
4. 먼저 `--apply` 없이 실행해 대상 수와 지정한 모델을 확인한다. 사용자가 변경을 요청했고 대상 범위가 맞으면 같은 인자에 `--apply`를 추가해 실행한다. 조회 요청은 미리보기에서 끝낸다. 실행 결과의 `complete`와 상태별 개수, 결과 JSON 경로를 읽어 보고한다.

```sh
./bin/hit-codex-sync-model --model gpt-6.2-sol --effort medium
./bin/hit-codex-sync-model --model gpt-6.2-sol --effort medium --apply
```

실행 파일은 Codex App Server를 통해 현재 호스트·로그인 계정의 저장된 모델·추론 수준·답변 속도를 변경한다. 기본 대상은 최근 30일 사용한 채팅이며 보관된 채팅에도 같은 기간 기준을 적용한다.

기간 기준은 App Server의 `recencyAt`(최근 사용 순서용 시각)이다. 설정 변경에도 갱신되는 `updatedAt`은 사용하지 않는다. 기준 시각을 포함한 최근 30×24시간을 선택하며, 최근 사용 시각이 없거나 0인 채팅은 기본 대상에서 제외한다. 조회 결과의 `since`는 기준 Unix 초다.

사용자가 **전체 기간**을 명시한 경우에만 `--all`을 추가한다. 단순히 "모든 채팅"이라고 하면 기본 기간인 최근 30일을 유지한다. 일부 채팅을 ID로 명시했다면 `--thread ID`를 반복한다. 명시한 ID에는 기간 제한을 적용하지 않는다. `--include-archived=false`로 보관 채팅을 제외한다. 다른 모델 제공자의 채팅은 건너뛴다. 사용자가 해당 제공자까지 명시했다면 `--provider ID --models`로 확인한 뒤 별도로 적용한다.

## 답변 속도

사용자가 답변 속도를 지정하면 모델과 추론 수준 인자에 `--service-tier`를 추가한다.

- 기본 또는 표준 → `--service-tier default`
- Fast → `--service-tier fast`
- Ultrafast → `--service-tier ultrafast`

속도를 지정하지 않으면 이 옵션을 생략해 채팅마다 기존 속도를 유지한다. 모델별 지원 옵션은 `--models` 결과의 `serviceTiers`에서 확인한다. 속도를 지정한 경우 상세 조회로 현재 값을 확인하고, 적용 결과와 백업에 속도도 기록한다. `--workers`는 일괄 변경 작업의 동시 처리 수를 조정한다.

```sh
./bin/hit-codex-sync-model --model gpt-6.2-sol --effort medium --service-tier fast
./bin/hit-codex-sync-model --model gpt-6.2-sol --effort medium --service-tier fast --apply
```

## 데스크톱 잠금 및 재시도

자동 종료·종료 대기·다시 열기는 macOS용이다. Windows·Linux에서 잠금이 남으면 해당 Codex를 직접 종료한 뒤 같은 결과 파일을 `--resume`으로 이어서 실행한다. 두 환경에서 `--retry-locked`는 잠긴 항목만 한 번 즉시 재시도한다.

현재 호스트에 제어 소켓이 있으면 연결하고, 없으면 독립 App Server로 저장된 설정을 변경한다. `--socket PATH`로 직접 지정하거나 `--standalone`으로 독립 서버를 선택할 수 있다.

앱에서 사용 중인 채팅은 `active writer` 오류로 잠길 수 있다. 실행 중인 다른 서버의 잠금을 가져오는 API는 없으므로, 잠금 보유 앱을 종료한 뒤 재시도한다. 앱 안에서 앱을 직접 종료하면 현재 작업이 끊기므로 **앱 외부 터미널**에서 다음 macOS 명령을 실행하고 사용자에게 앱 종료가 필요한 이유를 알린다.

```sh
./bin/hit-codex-sync-model --model gpt-6.2-sol --effort medium --apply --retry-locked --wait-seconds 3600 --reopen-app
```

사용자가 잠긴 채팅도 강제로 처리하도록 요청했다면 외부 터미널에서 `--force-locked`를 사용한다.

```sh
./bin/hit-codex-sync-model --model gpt-6.2-sol --effort medium --apply --force-locked
```

이 옵션은 먼저 가능한 항목을 처리하고, 잠금이 남았을 때만 실행 중인 ChatGPT/Codex 데스크톱 앱에 종료를 요청한다. **진행 중인 다른 채팅 작업도 중단된다.** 종료가 확인되면 잠긴 항목을 한 번 재시도하고 전체 성공 후 종료했던 앱을 다시 연다. `--reopen-app`을 따로 붙일 필요는 없다. `--apply`가 없는 조회에서는 앱을 종료하지 않는다.

앱을 종료하면 현재 채팅의 실행도 끊기므로 `--force-locked --apply`는 앱 외부 터미널에서 실행한다. 외부 터미널용 명령을 준비할 때는 이 스킬의 `bin/hit-codex-sync-model` 절대 경로와 확인한 대상 범위 인자를 사용한다. 사용자가 앱 종료와 실행을 이미 요청했다면 재승인을 요구하지 않는다. 앱 종료가 거부되거나 다른 CLI가 잠금을 보유했다면 실패를 보고한다. 알 수 없는 프로세스를 죽이거나 잠금 파일을 삭제하지 않는다. 일부 적용이나 재시도 실패 시 결과 기록을 남기고 멈추며, 앱 자동 재실행도 완료되지 않았음을 알린다.

답변 속도를 지정하지 않았고 목록에 모델·추론 수준이 모두 제공되며 이미 일치하면 상세 조회 없이 `already-matching`으로 집계한다. 목록 정보가 없거나 불일치한 채팅만 기본 4개씩 동시에 상세 처리한다. `--workers 1`로 순차 처리하거나 `--workers 8`로 늘릴 수 있으며 허용 범위는 1~16이다. 잠금 여부는 별도 조회 없이 상세 처리의 `active writer` 오류로 판정한다. 독립 서버의 목록 상태만 보고 다른 프로세스의 잠금을 추정하지 않는다.

이 명령은 먼저 가능한 채팅을 처리하고, 잠긴 채팅이 있으면 앱 종료를 기다린 뒤 그 항목만 다시 시도한다. 실패한 다른 항목과 원래 백업은 유지한다. 처음부터 앱 종료 후 시작하려면 `--retry-locked` 대신 `--wait-for-desktop-exit`을 사용한다. `--reopen-app`은 모든 대상이 성공한 경우 연결을 정리한 뒤 앱을 다시 연다. 종료 대기와 다시 열기는 macOS에서 지원한다.

독립 서버로 바꾼 저장 설정과 현재 열려 있는 화면의 모델 표시는 다를 수 있다. 전체 반영은 앱 재시작 후 확인한다. 잠금이나 미처리 항목이 남았다면 전체 성공이라고 보고하지 않는다.

## 백업·재개·복구

기본 기록 위치는 `~/.codex/model-sync-backups/`다. `CODEX_HOME`을 지정했다면 그 아래에 저장한다. `--report-dir PATH`로 바꿀 수 있다. 변경 전 설정과 목표 설정을 먼저 원자적으로 기록하고, 실제로 변경한 채팅은 변경 직후 다시 읽어 검증한다. 병렬 처리에서도 백업 파일 저장은 직렬화해 항목 누락을 방지한다. 이미 일치한 채팅은 개별 로그 대신 합계로 표시한다. 내용 전문이나 인증정보는 기록하지 않는다. 같은 Codex 홈에서는 동시 변경을 차단한다.

중단되었거나 일부가 실패했다면 **같은 결과 JSON**으로 재개한다. 재개·복구는 기록에 저장된 대상 범위를 유지하며 날짜가 지났다고 대상을 다시 선정하지 않는다. 이전 전체 기간 기록을 재개하면 그 기록의 전체 범위를 사용하므로 최근 30일로 새로 실행하는 요청에는 사용하지 않는다. 검증이 끝난 항목은 다시 바꾸지 않는다.

```sh
./bin/hit-codex-sync-model --resume /absolute/path/model-sync-REPORT.json --apply
./bin/hit-codex-sync-model --resume /absolute/path/model-sync-REPORT.json --apply --wait-for-desktop-exit --wait-seconds 3600 --reopen-app
```

복구도 조회 후 적용한다. 동기화 이후 사용자가 다시 바꾼 설정은 덮어쓰지 않는다. 원본 기록은 보존하고 별도 복구 결과를 쓴다.

```sh
./bin/hit-codex-sync-model --restore /absolute/path/model-sync-REPORT.json
./bin/hit-codex-sync-model --restore /absolute/path/model-sync-REPORT.json --apply
```

보관된 채팅은 잠시 보관을 해제하고 처리한 다음 원래 보관 상태로 되돌린다. 되돌리기에 실패하면 즉시 중단하고 해당 ID와 결과 파일을 보고한다. `missing`은 현재 계정·호스트 목록에서 찾을 수 없는 대상이다. 임의로 재개해 변경하지 않는다.

실행 결과의 경로, 성공·이미 일치·다른 제공자·잠금·실패·미처리 개수를 구분해 보고한다. `complete:false` 또는 종료 코드 1이면 미완료 또는 오류다. 원인을 확인하고 재개한다.

## 범위 및 검증

새 채팅 기본 모델, 다른 호스트, 일반 ChatGPT 대화, 서브에이전트의 모델은 별도 범위다. 서버가 `thread/settings/update`를 지원하지 않으면 실패를 보고하고 Codex 버전을 확인한다. SQLite·JSONL을 직접 수정해 우회하지 않는다.

Go 코드 검증은 `go -C app test -race ./...`와 `go -C app vet ./...`로 수행한다. 실제 적용 테스트는 반드시 별도 `CODEX_HOME`에서 생성한 테스트 채팅으로 진행한다.

- [공식 App Server 프로토콜](https://github.com/openai/codex/blob/main/codex-rs/app-server/README.md)
- 설치 버전 스키마: `codex app-server generate-ts --experimental --out TEMP_DIRECTORY`
