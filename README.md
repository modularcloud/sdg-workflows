# sdg-workflows

A collection of [loopx](https://github.com/lzrscg/loopx) workflows used to drive the ADR process for the `loopx` project itself.

> ⚠️ **Not ready for external use.** These workflows hardcode specific ADR filenames (e.g. `adr/0002-run-subcommand.md`, `adr/0004-tmpdir-and-args.md`) and expect a particular project layout (`SPEC.md`, `TEST-SPEC.md`, `adr/`, `PROMPT.md`). To use them on another project you will need to fork and edit the `index.sh` of each workflow. They are published here for reference and personal use, not as a general-purpose tool.

## Workflows

| Workflow | Purpose |
|---|---|
| `ralph` | Runs a prompt in a loop with Claude Code until a readiness check says the work is production-ready. |
| `review-adr` | Asks a reviewer whether a draft ADR can be marked accepted; iterates on Q&A until no questions remain. |
| `apply-adr` | After an ADR is accepted, iterates with a reviewer to apply the ADR's changes to `SPEC.md`. |
| `spec-test-adr` | After `SPEC.md` is updated for an ADR, iterates with a reviewer to update `TEST-SPEC.md`. |

## Install

```bash
# 1. Install loopx globally
npm install -g loopx

# 2. Install these workflows into your project's .loopx/ directory
cd /path/to/your/loopx-project
loopx install lzrscg/sdg-workflows            # installs all four
loopx install -w ralph lzrscg/sdg-workflows   # or install just one
```

If you plan to use the `api` reviewer mode, also run `npm install` inside each workflow directory that needs it:

```bash
cd .loopx/apply-adr && npm install
cd .loopx/review-adr && npm install
cd .loopx/spec-test-adr && npm install
```

## Required tools on PATH

- `claude` — Claude Code CLI (all workflows)
- `codex` — Codex CLI (used by `check-question.sh` in the ADR workflows, and when `LOOPX_REVIEWER=codex`)
- `jq`, `curl` — used by the Telegram reviewer and `ralph`

## Environment variables

| Variable | Used by | Notes |
|---|---|---|
| `TELEGRAM_BOT_TOKEN` | `ralph`, ADR workflows (telegram mode) | Bot token from @BotFather. |
| `TELEGRAM_CHAT_ID` | `ralph`, ADR workflows (telegram mode) | Numeric chat ID that will receive prompts and supply replies. |
| `LOOPX_REVIEWER` | ADR workflows | `telegram` (default), `codex`, or `api`. Selects how review prompts are sent out and answers are collected. |
| `OPENAI_API_KEY` | ADR workflows when `LOOPX_REVIEWER=api` | Required for the GPT batch reviewer. |
| `GPT_PRO_THINKING` | ADR workflows when `LOOPX_REVIEWER=api` | Optional, `medium` (default) / `high` / `xhigh`. |

Set them however you prefer — shell export, a local `.env` file passed via `loopx run -e .env`, etc.

## Required project files

Each workflow assumes these files already exist at the project root:

- `ralph` — `PROMPT.md`
- `review-adr` — `adr/0001-adr-process.md`, `adr/0004-tmpdir-and-args.md`, `SPEC.md`
- `apply-adr` — `adr/0001-adr-process.md`, `adr/0002-run-subcommand.md`, `SPEC.md`
- `spec-test-adr` — `adr/0001-adr-process.md`, `adr/0002-run-subcommand.md`, `SPEC.md`, `TEST-SPEC.md`

If you want to target a different ADR number, edit the `ADR_0002` / `ADR_0004` assignments at the top of the relevant `index.sh`.

## Run

```bash
# Ralph loop — keep running PROMPT.md through Claude until READY
export TELEGRAM_BOT_TOKEN=...
export TELEGRAM_CHAT_ID=...
loopx run ralph

# Review / apply / spec-test an ADR (pick a reviewer)
export LOOPX_REVIEWER=telegram   # or: codex, api
loopx run review-adr
loopx run apply-adr
loopx run spec-test-adr

# Cap iterations
loopx run -n 5 ralph
```

Each workflow halts on its own when done (`stop: true`) or when you hit the `-n` limit. Ctrl-C also exits cleanly.
