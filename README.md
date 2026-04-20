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
| `shared` | Library workflow — holds the reviewer-dispatch, Telegram Q&A loop, and feedback-done check used by the three ADR workflows. Not run directly. |

## 1. Install loopx and the workflows

```bash
# Install loopx globally
npm install -g loopx

# From your project root, install these workflows into .loopx/
loopx install modularcloud/sdg-workflows              # all five (recommended)
loopx install -w ralph modularcloud/sdg-workflows     # just ralph (standalone)

# Or install one ADR workflow — must also install `shared` (library workflow)
loopx install -w apply-adr modularcloud/sdg-workflows
loopx install -w shared    modularcloud/sdg-workflows
```

If you plan to use the `api` or `batch` reviewer (both call GPT-5.4-Pro), also install the `shared` workflow's Node dependencies:

```bash
cd .loopx/shared && npm install && cd -
```

## 2. Required tools on PATH

- `claude` — Claude Code CLI (all workflows)
- `codex` — Codex CLI (used by the `check-question` step in the ADR workflows, and when `LOOPX_REVIEWER=codex`)
- `jq`, `curl` — used by the Telegram reviewer and `ralph`

## 3. Set environment variables with `loopx env`

`loopx env set <NAME> <VALUE>` stores a variable globally so every `loopx run` picks it up.

### Required for everyone

Telegram is the default reviewer and `ralph` also posts to Telegram, so both vars are always needed:

```bash
loopx env set TELEGRAM_BOT_TOKEN <your-bot-token>   # from @BotFather
loopx env set TELEGRAM_CHAT_ID   <your-chat-id>     # numeric chat ID the bot posts to
```

When a prompt is sent to Telegram, the intended flow is: **copy the prompt, paste it into ChatGPT Pro, then paste ChatGPT's answer back into the Telegram chat.** Reply messages received within a 10s window are concatenated into one answer.

### Optional for everyone

Raises Claude Code's effort level so it thinks harder on each iteration:

```bash
loopx env set CLAUDE_CODE_EFFORT_LEVEL max
```

### Optional: use Codex instead of Telegram for ADR reviews

Sends review prompts to the local `codex` CLI instead of Telegram — fully automated, no copy/paste:

```bash
loopx env set LOOPX_REVIEWER codex
```

### Optional: use the OpenAI Responses API instead of Telegram for ADR reviews

Sends review prompts directly to GPT-5.4-Pro via the Responses API and waits for the reply synchronously — fully automated, no copy/paste, no batch polling:

```bash
loopx env set LOOPX_REVIEWER   api
loopx env set OPENAI_API_KEY   <your-openai-key>
loopx env set GPT_PRO_THINKING medium               # optional: medium (default) | high | xhigh
```

### Optional: use the OpenAI Batch API instead of Telegram for ADR reviews

Same as `api` but submits the prompt as a `/v1/batches` job and polls until completion. Trades latency for batch pricing:

```bash
loopx env set LOOPX_REVIEWER   batch
loopx env set OPENAI_API_KEY   <your-openai-key>
loopx env set GPT_PRO_THINKING medium               # optional: medium (default) | high | xhigh
```

### Inspect / remove

```bash
loopx env list
loopx env remove <NAME>
```

## 4. Required project files

Each workflow assumes these files exist at the project root:

- `ralph` — `PROMPT.md`
- `review-adr` — `adr/0001-adr-process.md`, `adr/0004-tmpdir-and-args.md`, `SPEC.md`
- `apply-adr` — `adr/0001-adr-process.md`, `adr/0002-run-subcommand.md`, `SPEC.md`
- `spec-test-adr` — `adr/0001-adr-process.md`, `adr/0002-run-subcommand.md`, `SPEC.md`, `TEST-SPEC.md`

To target a different ADR number, edit the `ADR_0002` / `ADR_0004` assignments at the top of the relevant `index.sh`.

## 5. Run

```bash
loopx run ralph
loopx run review-adr
loopx run apply-adr
loopx run spec-test-adr

loopx run -n 5 ralph   # cap iterations
```

Each workflow halts on its own when done (`stop: true`) or when `-n` is hit. Ctrl-C also exits cleanly.
