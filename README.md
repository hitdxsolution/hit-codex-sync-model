# hit-codex-sync-model

English · [한국어](README.ko.md)

A Go CLI and Codex skill for bulk-updating the model, reasoning effort, and response speed of existing Codex chats. By default, it targets chats used in the last 30 days on the current host and account. It backs up settings before applying changes and reads them back to verify the result.

## Install

1. Download the ZIP for your operating system and CPU from [GitHub Releases](https://github.com/hitdxsolution/hit-codex-sync-model/releases/latest). Choose `arm64` for Apple Silicon or Windows on Arm, and `amd64` for Intel or AMD x64. macOS filenames start with `darwin`.
2. Extract the ZIP. To use `$hit-codex-sync-model` in Codex, place the entire extracted `hit-codex-sync-model` folder in your Codex skills directory. The default is `~/.codex/skills/` on macOS and Linux or `%USERPROFILE%\.codex\skills\` on Windows.
3. Use a Codex CLI version that supports the Codex App Server thread settings API. The tool looks for `codex` on your `PATH` and can find the bundled CLI in the Codex or ChatGPT app on macOS. Use `--codex PATH` for a custom installation.

Release binaries run without Go installed. If you're building from source, install Go 1.26 or later. From the repository root, run `sh scripts/build.sh` on macOS or Linux. On Windows PowerShell, run the following.

```powershell
New-Item -ItemType Directory -Path bin -Force | Out-Null
go -C app build -trimpath -o ../bin/hit-codex-sync-model.exe .
```

## Quick start

Run these commands from the extracted `hit-codex-sync-model` folder.

```sh
./bin/hit-codex-sync-model --models
./bin/hit-codex-sync-model --model gpt-6.2-sol --effort medium --service-tier fast
./bin/hit-codex-sync-model --model gpt-6.2-sol --effort medium --service-tier fast --apply
```

On Windows PowerShell, use `.\bin\hit-codex-sync-model.exe` in place of `./bin/hit-codex-sync-model`. Run without `--apply` to preview the selection. Use `--models` to check which models and speed options your account supports.

To run it as a Codex skill, ask:

> $hit-codex-sync-model Set the chats I used in the last 30 days to gpt-6.2-sol, medium reasoning, and Fast.

| Option | Behavior |
|---|---|
| `--service-tier default` | Use the default response speed |
| `--service-tier fast` | Use Fast |
| `--service-tier ultrafast` | Use Ultrafast |
| Omit `--service-tier` | Keep each chat's current response speed |
| `--all` | Select chats from any date |
| `--thread ID` | Select a specific chat. Repeat to select more than one |
| `--workers 8` | Set the number of concurrent updates. The default is 4 |

## Resume and restore

The tool saves a JSON report in `~/.codex/model-sync-backups/` by default, or under the directory set by `CODEX_HOME`. Use the report path printed by the command to resume or restore.

```sh
./bin/hit-codex-sync-model --resume "/path/to/report.json" --apply
./bin/hit-codex-sync-model --restore "/path/to/report.json"
./bin/hit-codex-sync-model --restore "/path/to/report.json" --apply
```

If another Codex process holds a write lock, close it before resuming. On macOS, `--force-locked` can quit the desktop app and reopen it after a successful retry. Run that option from an external terminal because quitting the app interrupts active Codex work. On Windows and Linux, close the process holding the lock manually.

The tool sends requests to the local Codex App Server and does not start model responses. Reports contain chat IDs and settings before and after changes. Each installation uses the chats and account settings available on that host.

See the [skill instructions](SKILL.md) for details in Korean.
