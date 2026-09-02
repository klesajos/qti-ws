> 🌍 Read this in: **English** | [Česky](cheatsheet.cs.md)

# Claude Code cheat-sheet

The numbered guides teach you how to **extend** Claude Code ([skills](01-skills.md),
[hooks](02-hooks.md), [MCP](03-mcp.md), [plugins](04-plugins.md),
[subagents](05-agents.md), [workflows](06-workflows.md)). This page is the other
half: how to **operate**
it day to day — the flags, slash commands, shortcuts, prefixes and hook events
you reach for constantly but that don't each need a full guide.

> ✅ **Verified against Claude Code 2.1.257 (2026-09-03).** The CLI surface
> changes between releases — check `claude --version` and the
> [official docs](https://code.claude.com/docs/en/cli-reference) if something
> here doesn't match.

## Launch flags (`claude ...`)

| Command | What it does |
|---------|--------------|
| `claude` | Start an interactive session in the current folder |
| `claude "fix the lint errors"` | Start a session with the first prompt already queued |
| `claude -c` / `--continue` | Resume the **most recent** session in this folder |
| `claude -r` / `--resume [id]` | Resume a **specific** session (no id → interactive picker); `-n <name>` names a session so you can resume it by name |
| `claude -p "..."` / `--print` | **Headless**: run the prompt, print the result, exit. For scripts / CI |
| `claude --model opus` | Start on a specific model (`fable`, `opus`, `sonnet`, `haiku`, `opusplan`, `best`) |
| `claude --effort high` | Start at a reasoning effort (`low` … `max`, or `ultracode`) |
| `claude --agent go-reviewer` | Run the whole session as a named subagent |
| `claude --add-dir ../lib` | Give the session access to extra folders outside the repo |
| `claude --worktree` | Start in a fresh git worktree so the session can't touch your checkout |
| `claude --permission-mode plan` | Start in a permission mode (`manual`/`default` · `acceptEdits` · `plan` · `auto` · `dontAsk` · `bypassPermissions`) |
| `claude -p "..." --allowedTools "Read,Edit,Bash(git diff *)"` | Allowlist the tools an unattended / CI run may use |
| `claude -p "..." --output-format json` | Headless output as `json` or `stream-json` for scripts |
| `claude --restricted` | Locked-down session for shared machines: no command tools, file access confined to the working dirs |
| `claude --safe-mode` / `--bare` | Troubleshooting: run without any customisations (CLAUDE.md, hooks, skills, plugins, MCP) |

## Background sessions (`claude agents`)

A **background session** is a full conversation that keeps running detached (owned by
a supervisor process) — dispatch work, walk away, check back later.

| Command | What it does |
|---------|--------------|
| `claude --bg "..."` | Start a background session from the shell (long form `--background`) |
| `/bg <prompt>` | Background the current work from inside a session (also `/background`) |
| `claude agents` | Open the agent view — monitor / dispatch sessions (`--json` to script it) |
| `claude attach <id>` | Attach to a running background session |
| `claude logs <id>` | Print a background session's output |
| `claude stop <id>` | Stop a session (alias `claude kill`) |
| `claude respawn <id>` | Restart a stopped session with its conversation intact (`--all` for every one) |
| `claude rm <id>` | Remove a finished session from the list |
| `claude daemon status` | Show the supervisor that owns background sessions (`daemon stop --any` to stop it) |

**Headless (`-p`) vs. background (`--bg`) — when to use which:**

| | `claude -p "..."` (headless) | `claude --bg "..."` (background) |
|---|---|---|
| Runs | Foreground — prints, then exits | Detached — persists under a supervisor |
| Tied to your shell | Yes (close it → gone) | No (keeps running) |
| In `claude agents`? | No | Yes (attach / logs / stop) |
| Reach for it when | Scripting, CI, a one-off you pipe or parse | A long task you dispatch and revisit while you keep working |

Not to be confused with `/agents` (which now just reminds you to edit
`.claude/agents/` — the interactive wizard was removed) or
`claude --agent <name>` (run the whole session *as* a named subagent).

## Slash commands

Type `/` to autocomplete. Grouped by what they're for:

**Session & context**

| Command | What it does |
|---------|--------------|
| `/clear` | Wipe the conversation — clean slate for a new task (aliases `/new`, `/reset`) |
| `/compact [instructions]` | Summarise the conversation to reclaim context; optional focus, e.g. `/compact keep the undo work` |
| `/context` | Show what's filling the context window right now |
| `/btw <question>` | Quick side question — answered *without* adding to the conversation history |
| `/rewind` | Jump back to an earlier checkpoint (code and/or conversation) — also `Esc Esc` |
| `/resume` | Switch to another past session without leaving |
| `/rename <name>` | Rename the current session (easier to find later with `--resume`) |
| `/export` | Save the whole conversation to a Markdown file |
| `/plan <task>` | Enter plan mode straight from the prompt |
| `/subtask <prompt>` | Hand a side task to a subagent; the result comes back into this conversation |
| `/fork` | Copy the conversation into a new background session (its own worktree) |
| `/tasks` | List this session's background work — subagents, shell jobs, workflows |
| `/branch [name]` | Fork the *conversation* to try a different direction |
| `/goal <condition>` | Keep Claude working until a condition is met |
| `/diff` | Interactive viewer for uncommitted changes |
| `/cd <path>` | Move the session to another directory, keeping the conversation |

**Configuration**

| Command | What it does |
|---------|--------------|
| `/config` | Open the settings menu (model, theme, permissions, workflow size…) |
| `/model` | Switch model mid-session (`Option+P` opens the picker) |
| `/effort [level]` | Set reasoning effort — `/effort high`, `/effort auto`, `/effort ultracode` |
| `/fast` | Toggle fast mode (same model, faster output) |
| `/permissions` | View / grant standing tool permissions |
| `/sandbox` | Toggle filesystem / network isolation for Bash (or set `"sandbox.enabled": true`) |
| `/memory` | Edit persistent memory (`CLAUDE.md` files, auto-memory on/off) |
| `/init` | Generate a starter `CLAUDE.md` for this repo |
| `/status`, `/statusline` | Session health dashboard / customise the bottom info bar |
| `/doctor` | Diagnose a broken install (auth, network, dependencies) |
| `/reload-skills`, `/reload-plugins` | Re-read skills, workflows and plugins without restarting |

**Extensions** (the six mechanisms the guides cover)

| Command | What it does |
|---------|--------------|
| `/skill-name` | Run a skill, e.g. `/board-tests` (Example 1) |
| `/2048-dev:build-test` | Run a plugin's bundled skill (Example 4) |
| `/hooks` | View hooks configured for this session |
| `/mcp` | View MCP servers, their status, and authorise them |
| `/plugin` | Install / manage plugins and marketplaces |
| `/agents` | (Wizard removed) — edit `.claude/agents/*.md` directly, see Example 5 |
| `/workflows` | Watch, pause, stop or save running workflows (Example 6) |
| `/workflow-authoring` | Load the workflow script reference before writing one |

**Info**

| Command | What it does |
|---------|--------------|
| `/help` | List all commands, shortcuts and features |
| `/usage` | Remaining capacity on a subscription, spend and cache stats (`/cost` is an alias) |
| `/powerup` | Interactive feature lessons |
| `/insights` | HTML report analysing your recent sessions |

**Bundled skills worth knowing:** `/code-review` (alias `/review`),
`/security-review`, `/verify`, `/debug`, `/loop`, `/deep-research` (a bundled
workflow), `/fewer-permission-prompts`.

## Keyboard shortcuts

| Key | What it does |
|-----|--------------|
| `Ctrl+C` | Interrupt the current action (kill switch); twice to exit |
| `Ctrl+R` | Reverse-search your prompt/command history |
| `Esc` | Interrupt Claude, cancel the current input / close a menu |
| `Esc Esc` | Clear the draft — or, on an empty input, open **Rewind** |
| `Shift+Tab` | Cycle permission modes (Manual → Accept edits → Plan) |
| `Ctrl+G` | Open the current prompt in your `$EDITOR` |
| `Ctrl+O` | Toggle the transcript viewer |
| `Ctrl+T` | Toggle Claude's task checklist |
| `Ctrl+B` | Push the running task to the background |
| `Ctrl+S` | Stash / restore the prompt you're typing |
| `Option+P` / `Option+T` | Switch model / toggle extended thinking (macOS; `Alt+` on Windows/Linux) |
| `?` on an empty prompt | Show the shortcut help panel |

## Input prefixes

| Prefix | What it does |
|--------|--------------|
| `!` | **Bash mode** — run a shell command directly, no model, no tokens (e.g. `!git status`) |
| `@` | **File mention** — pull a specific file/folder into context (autocompletes) |
| `@agent-<name>` | **Delegate** — hand the task to a specific subagent (e.g. `@agent-go-reviewer`) |
| `:` | **Emoji shortcode** — `:tada:` |
| `\` + Enter | **Multiline** — continue the prompt on a new line without sending (also `Ctrl+J`, or `Shift+Enter` in most modern terminals) |

- **Remembering things:** the old `#` shortcut for "add this line to
  memory" is no longer part of the documented surface. Ask Claude to remember
  something (auto-memory) or edit `CLAUDE.md` via `/memory`.
- **Paste an image:** `Ctrl+V` (or drag-and-drop) drops a screenshot / mock-up into
  the prompt — great for UI work. On macOS use `Ctrl+V`, *not* `Cmd+V`.
- **Pipe a file:** `cat error.log | claude -p "explain this"` feeds stdin into a
  headless run (needs `-p`).

## Rewind & checkpoints

Every prompt is a **checkpoint**. Open the menu with `/rewind` or `Esc Esc`:

| Restore | Effect |
|---------|--------|
| **Code + conversation** | Roll both back to that prompt |
| **Conversation only** | Rewind the chat, keep the current code |
| **Code only** | Revert file edits, keep the chat |
| **Summarize from / up to here** | Compress the conversation after / before that point |

- Checkpoints persist across sessions and are kept ~30 days (`cleanupPeriodDays`).
- **Not captured:** files changed via Bash (`rm` / `mv` / `cp`) or edited outside
  Claude Code. **Not a Git replacement** — keep committing at milestones (see
  the [Git basics](00-git-basics.md) guide).

## Hook events

Where a hook can attach in the session lifecycle (see [Example 2](02-hooks.md)):

| Event | Fires | Typical use |
|-------|-------|-------------|
| `SessionStart` | Session begins | Load context, environment checks |
| `UserPromptSubmit` | You submit a prompt | Inject context, validate/rewrite the prompt |
| `PreToolUse` | Before a tool runs | **Block** dangerous actions (e.g. `git push --force`) |
| `PostToolUse` | After a tool succeeds | Format / lint — *this repo's `format-go` hook* |
| `Stop` | Claude finishes its turn | Verify tests — *this repo's `run-tests` hook* |
| `SubagentStart` / `SubagentStop` | A subagent starts / finishes | Set up or tear down what it needs |
| `PreModelSwitch` / `PostModelSwitch` | Model changes mid-session | Audit, re-apply settings |
| `SessionEnd` | Session ends | Cleanup, archiving, reporting |

Handler types: `command` (shell), `http`, `mcp_tool`, `prompt`, `agent`.

## Permission modes (cycle with `Shift+Tab`)

| Mode | Claude may… |
|------|-------------|
| **Manual** (`default`) | Ask before every edit and command — a grey ⏸ badge in the footer shows you're here |
| **Accept edits** (`acceptEdits`) | Edit files freely, still ask before shell commands |
| **Plan** (`plan`) | Read-only — research and draft a plan, change nothing |
| **Auto** (`auto`) | Act, with a background classifier blocking risky actions instead of asking |
| **Bypass permissions** | Do anything without asking ("YOLO" — sandbox only). Not in the `Shift+Tab` cycle — enable explicitly with `--dangerously-skip-permissions` |

One more mode exists — **dontAsk** (only pre-approved tools, everything else
denied) — for unattended runs; see the
[permission-modes docs](https://code.claude.com/docs/en/permission-modes).

## Models & effort

| Model | Best for |
|-------|----------|
| `fable` | Long autonomous sessions, investigation, verification (Fable 5.1, 1M context) |
| `opus` | Hardest reasoning — architecture, tricky bugs (Opus 5) |
| `sonnet` | Daily driver — features, refactors, tests (Sonnet 5) |
| `haiku` | Fast & cheap — boilerplate, quick checks, parallel subagents |
| `opusplan` | Opus to plan, then auto-switches to Sonnet to implement |
| `best` / `default` | The strongest model available to you / your account's default |

**Effort** controls how deeply Claude reasons — higher = better answers, slower,
more tokens. The ladder, low → high: `low` · `medium` · `high` · `xhigh` · `max`.
`ultracode` sits on top: `xhigh` **plus** Claude writes a workflow for every
substantive task (see [Example 6](06-workflows.md)).

| Set effort | How |
|------------|-----|
| In session | `/effort` (slider) · `/effort high` (jump to a level) · `/effort auto` (back to the model default) |
| At launch | `claude --effort xhigh "..."` — this session only |
| Persisted | `"effortLevel": "high"` in `settings.json`, or env `CLAUDE_CODE_EFFORT_LEVEL` (`max`/`ultracode` can't be persisted this way) |
| Per agent / skill | an `effort:` field in the frontmatter |

**One-shot:** drop `ultrathink` anywhere in a prompt to push that single turn
harder *without* changing the session effort — only `ultrathink` is a real
keyword (`think` / `think hard` are just ordinary words).

## Context economics (why `/compact` matters)

- Claude re-reads the **whole** conversation each turn, so cost grows with
  history. `/compact` summarises to stay lean; `/clear` resets for a new task.
- **Prompt caching** gives a large discount on the stable prefix (your
  `CLAUDE.md`, project structure) — it's automatic; keep that prefix stable.
- **Be mindful of extensions:** every active MCP server, skill and subagent
  consumes context. Enable what the task needs, not everything.

## What's legacy in 2026 (the convergence)

Custom slash commands have been **merged into Skills** (Example 1). The old
command files keep working, but new work should use Skills:

| Was its own thing | Now | Status |
|-------------------|-----|--------|
| Custom slash commands (`.claude/commands/*.md`) | **Skills** (`.claude/skills/<name>/SKILL.md`) | Merged — both still create `/name`; legacy files keep working |
| `/agents` wizard | Edit `.claude/agents/*.md` by hand (or ask Claude) | Wizard removed in 2.1.198 |
| `Task` tool (spawning subagents) | `Agent` tool | Renamed; old name still works as an alias |
| `/cost` | `/usage` | `/cost` is an alias now |
| `/review` | `/code-review` | Alias |
| Todo / task-tracking tools | Off by default on newer models | `CLAUDE_CODE_ENABLE_TODO_TOOLS=1` restores them |
| `ultraplan` | Removed (2.1.222) | Use plan mode + `ultrathink` |

(Output styles are a *separate*, still-current feature — they modify the system
prompt, skills load task instructions — see the
[output-styles docs](https://code.claude.com/docs/en/output-styles); only the
standalone `/output-style` command was deprecated, in favour of `/config`.)

A skill is a strict superset of an old command: same `/name` invocation, **plus**
optional autonomous loading, a supporting-files folder, and invocation control.
Want a skill to behave like the old "manual only" command? Add
`disable-model-invocation: true` to its frontmatter.

Source: [official Skills docs](https://code.claude.com/docs/en/skills) — *"Custom
commands have been merged into skills."*

---

**Keep this open while you work.** It's reference, not a tutorial — for the *how*
and *why* of each extension mechanism, see the numbered guides in the
[index](README.md).
