# sdg-workflows

A collection of [loopx](https://github.com/lzrscg/loopx) workflows used to drive the ADR process for the `loopx` project itself.

> ⚠️ **Not ready for external use.** These workflows expect a particular project layout (`SPEC.md`, `TEST-SPEC.md`, `adr/`, `PROMPT.md`) and the ADR lifecycle defined in `adr/0001-adr-process.md`. They are published here for reference and personal use, not as a general-purpose tool.

## Workflows

| Workflow | Purpose |
|---|---|
| `ralph` | Runs a prompt in a loop with Claude Code until a readiness check says the work is production-ready. |
| `review-adr` | Asks a reviewer whether a draft ADR can be marked accepted; iterates on Q&A until no questions remain. |
| `apply-adr` | After an ADR is accepted, iterates with a reviewer to apply the ADR's changes to `SPEC.md`. |
| `spec-test-adr` | After `SPEC.md` is updated for an ADR, iterates with a reviewer to update `TEST-SPEC.md`. |
| `review-spec` | Asks a reviewer whether `SPEC.md` is implementation-ready; iterates on Q&A until no questions remain. |
| `review-test-spec` | Asks a reviewer whether `TEST-SPEC.md` covers `SPEC.md` correctly and completely; iterates on Q&A until no questions remain. |
| `shared` | Library workflow — holds the reviewer-dispatch, Telegram Q&A loop, and feedback-done check used by the ADR and spec workflows. Not run directly. |

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

If you plan to use the `api` or `batch` reviewer (both call GPT-5.5-Pro), also install the `shared` workflow's Node dependencies:

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

`TELEGRAM_CHAT_ID` must be a **forum-enabled supergroup** (Group Info → Edit → Topics → on), and the bot must be an admin with *Manage Topics* rights. Each run posts into its own topic named `<cwd> / <workflow>[ / ADR-NNNN]`, so concurrent runs in different repos can't step on each other. Topic IDs are cached in `~/.cache/loopx-telegram/`.

When a prompt is sent to Telegram, the intended flow is: **copy the prompt, paste it into ChatGPT Pro, then paste ChatGPT's answer back into the corresponding topic.** Reply messages received within a 10s window are concatenated into one answer.

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

Sends review prompts directly to GPT-5.5-Pro via the Responses API and waits for the reply synchronously — fully automated, no copy/paste, no batch polling:

```bash
loopx env set LOOPX_REVIEWER   api
loopx env set OPENAI_API_KEY   <your-openai-key>
loopx env set GPT_PRO_THINKING medium               # optional: medium (default) | high | xhigh
loopx env set OPENAI_FLEX      true                 # optional: route through the flex service tier (lower cost, slower / capacity-dependent)
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
- `review-adr` — `adr/0001-adr-process.md`, `adr/NNNN-*.md` (target ADR), `SPEC.md`
- `apply-adr` — `adr/0001-adr-process.md`, `adr/NNNN-*.md` (target ADR), `SPEC.md`
- `spec-test-adr` — `adr/0001-adr-process.md`, `adr/NNNN-*.md` (target ADR), `SPEC.md`, `TEST-SPEC.md`
- `review-spec` — `SPEC.md`
- `review-test-spec` — `SPEC.md`, `TEST-SPEC.md`

The target ADR is selected per run via the `ADR` env var (see §5).

### Bootstrap the `adr/` directory

If your project does not yet have an `adr/` directory or an `adr/0001-adr-process.md`, create one and drop in the bundled template:

```bash
mkdir -p adr
curl -fsSL -o adr/0001-adr-process.md \
  https://raw.githubusercontent.com/modularcloud/sdg-workflows/main/adr/0001-adr-process.md
```

Then author your own target ADRs as `adr/NNNN-short-slug.md` (e.g. `adr/0002-run-subcommand.md`).

## 5. Run

The three ADR workflows resolve their target ADR from the `ADR` env var. Either `ADR=4` or `ADR=0004` works — the value is zero-padded to four digits and matched against `adr/NNNN-*.md`.

```bash
loopx run ralph

# ADR workflows — set ADR per run
ADR=4 loopx run review-adr
ADR=2 loopx run apply-adr
ADR=2 loopx run spec-test-adr

# Spec-only review — no ADR required
loopx run review-spec
loopx run review-test-spec

loopx run -n 5 ralph   # cap iterations
```

Each workflow halts on its own when done (`stop: true`) or when `-n` is hit. Ctrl-C also exits cleanly.
