> 🌍 Read this in: **English** | [Česky](02-hooks.cs.md)

# Example 2: Project-scoped hooks

## What is a hook?

A **hook** is a shell command that Claude Code runs **automatically** when a
certain event happens — for example "after Claude edits a file" or "before
Claude runs a terminal command".

The key difference from a skill: a skill is *advice* the model may follow;
a hook is *enforcement* that runs outside the model, **every single time**.
Use hooks for things that must never be skipped: formatting, linting,
blocking dangerous commands.

## What this example does

Every time Claude edits or creates a file, our first hook checks whether it's
a Go file (`.go`) — and if so, runs `gofmt` on it. Result: Claude's code
always lands formatted the way `gofmt` wants it, even if the model wrote it
messy. (If `goimports` is installed, the hook uses that instead — it is
`gofmt` plus import-block fixing.)

Three files are involved:

```
.claude/settings.json        ← registers WHEN to run each hook
.claude/hooks/format-go.sh   ← the script that says WHAT to do after an edit
.claude/hooks/run-tests.sh   ← a second hook, run when Claude finishes a turn
```

This guide walks through the format hook in full, then the **second hook** on
a different event — a `Stop` hook that runs the tests when Claude finishes a
turn — and explains the five hook **types**.

## Part 1: the registration (`.claude/settings.json`), line by line

```json
{
  "hooks": {
    "PostToolUse": [
      {
        "matcher": "Edit|Write",
        "hooks": [
          {
            "type": "command",
            "command": "${CLAUDE_PROJECT_DIR}/.claude/hooks/format-go.sh"
          }
        ]
      }
    ]
  }
}
```

- `"hooks"` — the top-level section of settings where all hooks live.
- `"PostToolUse"` — the **event**: fire *after* Claude successfully uses a
  tool. Other useful events: `PreToolUse` (before a tool runs — can block
  it), `SessionStart`, `UserPromptSubmit`, `Stop` (when Claude finishes).
- `"matcher": "Edit|Write"` — a filter on the **tool name**. The `|` means
  "or", like in regular expressions: run only when the tool was `Edit` or
  `Write` (the two tools Claude uses to change files). Without a matcher the
  hook would also fire after every `Read`, `Bash`, etc.
- `"type": "command"` — this hook runs a shell command (other types exist —
  see [Hook types](#hook-types-command-http-mcp_tool-prompt-agent) below).
- `"command": "${CLAUDE_PROJECT_DIR}/..."` — the script to run.
  `${CLAUDE_PROJECT_DIR}` is a variable Claude Code replaces with the
  absolute path of the repo root — so the hook works no matter which
  subdirectory Claude is currently in.

## Part 2: the script (`.claude/hooks/format-go.sh`), line by line

```bash
#!/usr/bin/env bash
set -euo pipefail

input=$(cat)
```

- `#!/usr/bin/env bash` — the "shebang": tells the OS to run this file
  with bash. (On Windows that is Git Bash — see the
  [install guide](00-install-claude-code.md).)
- `set -euo pipefail` — safety switches: stop on any error (`-e`), treat
  unset variables as errors (`-u`), fail a pipeline if any step fails
  (`pipefail`). Standard practice for every bash script.
- `input=$(cat)` — **this is how hooks receive data.** Claude Code sends
  the details of the tool call as JSON on **stdin** (standard input).
  `cat` reads all of it; we store it in the variable `input`.

```bash
if command -v jq >/dev/null 2>&1; then
    file_path=$(printf '%s' "$input" | jq -r '.tool_input.file_path // empty')
elif command -v python3 >/dev/null 2>&1; then
    file_path=$(printf '%s' "$input" | python3 -c \
        'import json,sys; print(json.load(sys.stdin).get("tool_input", {}).get("file_path", ""))')
else
    file_path=$(printf '%s' "$input" | sed -n 's/.*"file_path"[[:space:]]*:[[:space:]]*"\([^"]*\)".*/\1/p')
fi
```

- `command -v jq` — checks whether the JSON tool `jq` is installed
  (`>/dev/null 2>&1` just hides the output of the check).
- `jq -r '.tool_input.file_path // empty'` — extracts the edited file's
  path from the JSON. The input looks like
  `{"tool_name": "Edit", "tool_input": {"file_path": "/path/to/board.go", ...}}`,
  so `.tool_input.file_path` navigates to the path. `// empty` means
  "if missing, output nothing instead of the word null".
- The `elif` branch does exactly the same with Python — a fallback for
  machines without `jq`.
- The last branch is a plain `sed` pattern for machines with neither (a
  fresh Windows install with Git Bash has bash and sed, but often no `jq`
  and no `python3`). It handles the ordinary case; install `jq` for anything
  fancier.

```bash
if command -v goimports >/dev/null 2>&1; then
    formatter="goimports -w"
elif command -v gofmt >/dev/null 2>&1; then
    formatter="gofmt -w"
else
    exit 0
fi
```

- Picks the formatter. `gofmt` ships with every Go installation, so unlike
  `clang-format` in the C++ world it is practically always there;
  `goimports` (from `golang.org/x/tools`) is an optional upgrade that also
  adds/removes import lines.
- `exit 0` — if there's no formatter at all, exit **successfully** and do
  nothing. A non-zero exit code would surface as an error to Claude;
  "formatter not installed" shouldn't break the session.

```bash
case "$file_path" in
    *.go)
        if [ -f "$file_path" ]; then
            $formatter "$file_path"
            echo "format-go hook: formatted ${file_path##*/}"
        fi
        ;;
esac

exit 0
```

- `case ... in *.go)` — pattern match on the file extension. Only Go
  sources proceed; a `.md` or `.json` file falls through and the script just
  exits.
- `[ -f "$file_path" ]` — "does the file exist?" (it might have been
  deleted in the meantime).
- `$formatter "$file_path"` — the actual work: `-w` (already inside the
  variable) means "write", i.e. rewrite the file with formatted content.
- `echo ...` — whatever a hook prints is shown in the Claude Code
  transcript, so you can see the hook did its job.
  (`${file_path##*/}` strips the directory part, leaving just the filename.)

## Create your own hook, step by step

1. **Write a script** in `.claude/hooks/`, e.g. `my-hook.sh`. Start from
   the skeleton above: read stdin, extract what you need, act, `exit 0`.
2. **Make it executable** — this is the step everyone forgets:
   ```bash
   chmod +x .claude/hooks/my-hook.sh
   ```
3. **Register it** in `.claude/settings.json` under the right event +
   matcher (see Part 1).
4. **Test it manually first** — don't debug inside Claude. Fake the stdin
   JSON yourself:
   ```bash
   echo '{"tool_input":{"file_path":"internal/board/board.go"}}' | .claude/hooks/my-hook.sh
   ```
5. **Start a new Claude Code session** and trigger the event for real.
6. **Commit both files.**

## Try the demo

Ask Claude to add a method to `internal/board/board.go` and not worry about
formatting. After the edit, run `git diff` — the code is already formatted,
and the line `format-go hook: formatted board.go` appears in the transcript.
Run `gofmt -l .` — it prints nothing, because nothing is left unformatted.

## A second hook: report the tests (the `Stop` event)

The format hook reacts to a **tool** (`Edit`/`Write`). Hooks can also react to
the **session lifecycle**. This repo ships a second hook on the `Stop` event —
which fires once, when Claude finishes its turn — to run the test suite and
report whether the tree is still green.

Its registration sits next to the first hook in `.claude/settings.json`. Note
there is **no `matcher`**: `Stop` isn't about a tool, so there's nothing to
filter on. Note also the `timeout` — a full `go test ./...` compiles first, and
we'd rather the hook be cancelled than stall the session:

```json
"Stop": [
  {
    "hooks": [
      { "type": "command",
        "command": "${CLAUDE_PROJECT_DIR}/.claude/hooks/run-tests.sh",
        "timeout": 120 }
    ]
  }
]
```

The script (`.claude/hooks/run-tests.sh`) is deliberately **advisory**: it
drains the Stop-event JSON, runs `go test ./...`, prints one line, and always
`exit 0` — so it never interrupts you. Its executable lines, verbatim from the
file:

```bash
cat >/dev/null
project_dir="${CLAUDE_PROJECT_DIR:-$(pwd)}"
if [ ! -f "$project_dir/go.mod" ] || ! command -v go >/dev/null 2>&1; then
    exit 0
fi
set +e
output=$(cd "$project_dir" && go test ./... 2>&1)
status=$?
set -e
passed=$(printf '%s\n' "$output" | grep -c '^ok' || true)
failing=$(printf '%s\n' "$output" | awk '/^FAIL[[:space:]]/ {print $2}' | paste -sd ',' - || true)
if [ "$status" -eq 0 ]; then
    echo "run-tests hook: ✓ go test ./... — ${passed} package(s) ok"
else
    echo "run-tests hook: ✗ go test ./... failed in ${failing:-the build} (run 'go test ./...' for details)"
fi
```

- **`cat >/dev/null`** drains stdin. We don't need any field of the Stop
  event, but a hook must read its input so Claude Code isn't left writing to
  a closed pipe.
- **The guard** `if [ ! -f "$project_dir/go.mod" ] || ! command -v go …; then exit 0; fi`
  bails out if there's no `go.mod` **or** no `go` on PATH — so the hook stays
  silent in the wrong directory or on a machine without Go instead of erroring.
- **`set +e` … `set -e`** suspends "stop on error" around the test run: a red
  suite makes `go test` exit non-zero, and we want to *report* that, not have
  the hook abort. `status=$?` keeps the exit code.
- **The summary** counts the per-package result lines `go test` prints —
  `ok <pkg>` for green packages, `FAIL <pkg>` for red ones — and joins the
  failing package names with commas.
- **The verdict** prints `✓` with the package count when the status was 0,
  otherwise `✗` with the failing packages and a hint to rerun `go test`.

**Want it to *block* instead of inform?** A `Stop` hook that exits with code
**2**, or prints `{"decision": "block", "reason": "..."}` on stdout, tells
Claude it is **not** done — that's how you enforce "tests must pass before you
stop". We keep ours advisory so a live workshop session never gets stuck in a
fix-tests loop; flip it to blocking when you want a hard gate
([Exercise 2](exercises.md#2-hook) does exactly that).

## Try the second hook

Ask Claude anything that ends a turn. When the turn finishes, the `Stop` hook
runs the suite and an advisory line appears in the transcript:

```
run-tests hook: ✓ go test ./... — 1 package(s) ok
```

Now deliberately break a test (change an expected value), ask Claude something
trivial, and watch the line turn into `✗ … failed in
github.com/klesajos/qti-ws/internal/board`. Revert the change.

## Hook types: command, http, mcp_tool, prompt, agent

Both hooks above use `"type": "command"` — a shell script. That's the most
common type, but not the only one. A handler can be one of five types, trading
determinism for judgement:

| Type | What runs | Use it for |
|------|-----------|------------|
| `command` | A shell script (our two hooks) | Deterministic, fast checks — format, lint, run tests, block a forbidden command |
| `http` | An HTTP POST of the event JSON to a URL | Central policy services, audit logging — the endpoint answers with the same JSON as a command hook |
| `mcp_tool` | A tool on an already-connected MCP server ([Example 3](03-mcp.md)) | Reuse a capability you already have (write to a tracker, check a policy service) |
| `prompt` | A single-turn call to a Claude model | A quick judgement call — "is this commit message descriptive?" — returning a yes/no decision |
| `agent` | A multi-turn subagent with tool access (`Read`, `Grep`, …) | Deep verification — "read the diff and confirm no secrets were added" — before proceeding (experimental) |

A `prompt` hook is registered with the text to evaluate instead of a command:

```json
{ "type": "prompt",
  "prompt": "Does the staged diff add a test for every new exported function? Answer yes or no." }
```

Rule of thumb: reach for `command` first (free and instant), `prompt` when the
check needs language understanding, and `agent` only when the check itself has
to explore the codebase. This repo ships `command` hooks; the other four are
worth knowing exist.

## Optional parameters

A hook handler accepts more than `type` and `command`. The most useful
optional fields:

| Field | What it does |
|-------|--------------|
| `timeout` | Seconds before the hook is cancelled (default 600 for `command`). Our Stop hook sets 120 so a slow test run can't stall the session |
| `statusMessage` | Custom spinner text while the hook runs, e.g. `"Formatting Go..."` |
| `if` | Extra filter using permission-rule syntax, e.g. `"if": "Edit(*.go)"` — more precise than `matcher`, which only sees the tool name. A path pattern like `internal/**` matches only under the working directory; use `**/internal/**` for any depth |
| `once: true` | Run only once, then deregister — **only honoured for hooks declared in a skill's frontmatter**, ignored in `settings.json` |

And the events: this example uses `PostToolUse` and `Stop`, but hooks can
attach to the whole session lifecycle. The ones worth knowing first:

| Event | Fires |
|-------|-------|
| `PreToolUse` | Before a tool runs — **can block it** (e.g. forbid `git push --force`) |
| `PostToolUse` | After a tool succeeds (our format hook) |
| `SessionStart` | When a session begins — environment checks, loading context |
| `UserPromptSubmit` | When you submit a prompt — can inject extra context |
| `Stop` | When Claude finishes its turn (our test hook) — can block the stop |
| `SubagentStart` / `SubagentStop` | Around a subagent run ([Example 5](05-agents.md)) |
| `SessionEnd` | When the session ends — cleanup, reporting |

There are more (`PostToolUseFailure`, `PreCompact`, `PreModelSwitch`,
`DirectoryAdded`, …). Full list of events and fields:
[official hooks documentation](https://code.claude.com/docs/en/hooks).

## Where it works: CLI, Desktop app, Cowork

| Platform | Works? | Setup |
|----------|--------|-------|
| **Claude Code CLI** (terminal) | ✅ Yes | Nothing extra — hooks in `.claude/settings.json` load at session start |
| **Claude Desktop app — Code tab** | ✅ Yes | Same engine, same config files as the CLI. Confirm the one-time project trust dialog; hooks then run identically |
| **Cowork** (in the Desktop app) | ❌ No | Cowork's sandboxed VM does not execute project-scoped hooks from `.claude/settings.json`. There is no direct equivalent — hooks shipped inside an installed plugin are the closest option |

Note for the workshop: this is the clearest platform difference of the six
examples. Hooks are a *local automation* feature — if your workflow depends
on them (formatting, test gates), run it in the CLI or the Desktop Code tab,
not in Cowork.

## Troubleshooting

- **Hook never runs** → new session needed after editing `settings.json`;
  also check the script is executable (`ls -l .claude/hooks/`).
- **"Permission denied"** → you skipped `chmod +x`.
- **Hook errors break the flow** → make sure every "nothing to do" path
  ends in `exit 0`, not an error.
- **On Windows the hook fails with `bash: not found`** → install Git for
  Windows so Claude Code can use Git Bash as its shell (see the
  [install guide](00-install-claude-code.md)).
