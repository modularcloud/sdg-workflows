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

## 1. Install loopx and the workflows

```bash
# Install loopx globally
npm install -g loopx

# From your project root, install these workflows into .loopx/
loopx install modularcloud/sdg-workflows            # all four
loopx install -w ralph modularcloud/sdg-workflows   # or just one
```

If you plan to use the `api` reviewer (GPT batch), also install its Node dependencies:

```bash
cd .loopx/apply-adr && npm install && cd -
cd .loopx/review-adr && npm install && cd -
cd .loopx/spec-test-adr && npm install && cd -
```

## 2. Required tools on PATH

- `claude` — Claude Code CLI (all workflows)
- `codex` — Codex CLI (used by the `check-question` step in the ADR workflows, and when `LOOPX_REVIEWER=codex`)
- `jq`, `curl` — used by the Telegram reviewer and `ralph`

## 3. Set environment variables with `loopx env`

`loopx env set <NAME> <VALUE>` stores a variable globally so every `loopx run` picks it up. Only set the ones that apply to the scenario you're using.

### Recommended for every scenario (optional)

Raises Claude Code's effort level so it thinks harder on each iteration. Recommended but not required:

```bash
loopx env set CLAUDE_CODE_EFFORT_LEVEL max
```

### Scenario A — `ralph` loop

`ralph` sends Telegram pings at the start of each iteration and when the work is judged ready. Both vars are required:

```bash
loopx env set TELEGRAM_BOT_TOKEN <your-bot-token>   # from @BotFather
loopx env set TELEGRAM_CHAT_ID   <your-chat-id>     # numeric chat ID the bot posts to
```

### Scenario B — ADR workflows with the **Telegram** reviewer (default)

The workflow sends the review prompt to Telegram and waits for your reply. The intended flow is: **copy the prompt out of Telegram, paste it into ChatGPT Pro, then paste ChatGPT's answer back into the Telegram chat.** Reply messages received within a 10s window are concatenated into one answer.

```bash
loopx env set LOOPX_REVIEWER    telegram            # default; setting it explicitly is optional
loopx env set TELEGRAM_BOT_TOKEN <your-bot-token>
loopx env set TELEGRAM_CHAT_ID   <your-chat-id>
```

### Scenario C — ADR workflows with the **Codex** reviewer

Sends the prompt to the local `codex` CLI, no Telegram round-trip. Fully automated:

```bash
loopx env set LOOPX_REVIEWER codex
```

### Scenario D — ADR workflows with the **API** reviewer (GPT-5.4-Pro batch)

Submits the prompt as an OpenAI batch job and polls until completion:

```bash
loopx env set LOOPX_REVIEWER  api
loopx env set OPENAI_API_KEY  <your-openai-key>
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
