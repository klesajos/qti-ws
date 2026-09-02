> 🌍 Read this in: **English** | [Česky](01-skills.cs.md)

# Example 1: Project-scoped skill

## What is a skill?

A **skill** is a Markdown file with instructions for Claude that gets loaded
**only when relevant**. Compare it to `CLAUDE.md`:

- `CLAUDE.md` is loaded into **every** conversation — keep it short.
- A skill is loaded **on demand** — when your request matches its
  description, or when you invoke it by name with `/skill-name`.

This keeps Claude's context small while still giving it deep project
knowledge when needed.

## What this example does

The skill `board-tests` teaches Claude *how this project writes unit tests*:
which APIs to use, what style to follow, how to run the test suite. Without
it, Claude would guess — and every workshop attendee would get slightly
different tests (some with `testify`, some with `reflect.DeepEqual`, some
testing through the game loop).

## The file, line by line

The skill lives at `.claude/skills/board-tests/SKILL.md`. The path matters:

```
.claude/              ← project-scoped Claude Code config (committed to git)
└── skills/           ← all project skills live here
    └── board-tests/  ← one folder per skill; the folder name = skill name
        └── SKILL.md  ← must be named exactly SKILL.md
```

The file starts with **frontmatter** — a small YAML block between `---`
lines that describes the skill:

```yaml
---
name: board-tests
description: Use when writing or extending Go unit tests for the 2048 board
  logic (SlideLineLeft, Board.Move, win/game-over detection)...
---
```

What each line does:

- `---` — opens (and later closes) the frontmatter block. Required.
- `name: board-tests` — the skill's identifier. This is what you type after
  the slash: `/board-tests`. Optional — if omitted, the folder name is used.
- `description: ...` — **the most important line.** Claude reads only the
  descriptions of all skills at session start. When your request matches a
  description, Claude loads that skill's full body. Write the description as
  a trigger: say *when* the skill should be used, not what's inside it.

Everything **below** the frontmatter is the body — plain Markdown
instructions Claude follows once the skill is active. In our case: test
conventions (Arrange-Act-Assert, one behaviour per `Test...` function, `==`
on arrays, no assertion libraries), the key `board` API, and the exact
`go test` commands to verify the work.

## Create your own skill, step by step

1. **Create the folder.** In your repo root:
   ```bash
   mkdir -p .claude/skills/my-skill
   ```
   (`mkdir -p` creates the whole path at once, including parent folders.)

2. **Create the file** `.claude/skills/my-skill/SKILL.md` with this skeleton:
   ```markdown
   ---
   name: my-skill
   description: Use when <describe the situation that should trigger this skill>.
   ---

   # My skill

   <Instructions for Claude. Write them like you'd brief a new colleague:
   what to do, which files/APIs to use, how to verify the result.>
   ```

3. **Restart Claude Code** (or start a new session) so it picks up the new
   skill. Skills are discovered at session start. (In a running session,
   `/reload-skills` re-reads the skill directories without restarting.)

4. **Test the manual trigger:** type `/my-skill` — the skill's body loads
   into the conversation.

5. **Test the automatic trigger:** phrase a request that matches your
   `description`. Claude should announce it's using the skill before acting.

6. **Commit it:**
   ```bash
   git add .claude/skills/my-skill
   git commit -m "Add my-skill"
   ```
   Now everyone who clones the repo has it. That's what "project-scoped"
   means. (Personal skills you don't want to share go in `~/.claude/skills/`
   in your home directory instead.)

## Try the demo

Ask Claude: *"Add a test for sliding the line {4, 4, 8, 0}"*.

There is a deliberate `TODO (participants)` in `internal/board/board_test.go`
for exactly this case. Thanks to the skill, Claude will call
`SlideLineLeft(&line)` directly (it needs the score), name the function
`TestSlideLineLeft_...`, compare the `Line` with `!=`, and run `go test` to
verify — the same way for every attendee. Watch what the test reports: the
line *should* become `{8, 8, 0, 0}`, and the skill tells Claude to report a
failure rather than "fix" the expectation.

## Optional frontmatter parameters

Our skill uses only `name` and `description`, but the frontmatter supports
more. The most useful optional fields:

| Field | What it does |
|-------|--------------|
| `argument-hint` | Autocomplete hint shown after the slash command, e.g. `[test-name]` |
| `disable-model-invocation: true` | Claude never loads the skill on its own — only you can, via `/name`. Good for destructive workflows (deploy, release). Also stops the skill being preloaded into subagents |
| `user-invocable: false` | The opposite: hidden from the `/` menu, only Claude can load it |
| `allowed-tools` | Tools the skill may use without asking permission during the turn that invoked it, e.g. `Bash(go test:*)` |
| `model` / `effort` | Use a specific model or reasoning effort while the skill is active |
| `context: fork` | Run the skill in an isolated subagent so it doesn't fill your main conversation's context (`agent:` picks which subagent type; runs in the background by default) |
| `paths` | Glob patterns — only auto-invoke when matching files are involved, e.g. `internal/board/**` |
| `hooks` | Hooks registered when the skill is invoked (see [Example 2](02-hooks.md)); the only place the `once: true` hook option is honoured |

Full, current list: [official skills documentation](https://code.claude.com/docs/en/skills).

## Where it works: CLI, Desktop app, Cowork

| Platform | Works? | Setup |
|----------|--------|-------|
| **Claude Code CLI** (terminal) | ✅ Yes | Nothing extra — skills in `.claude/skills/` are discovered automatically when you run `claude` in the repo |
| **Claude Desktop app — Code tab** | ✅ Yes | Open the project folder in the Code tab. Desktop runs the same engine as the CLI and shares all project config. You'll confirm a one-time trust dialog for the project |
| **Cowork** (in the Desktop app) | ❌ No | Cowork runs tasks in its own sandboxed VM and does **not** load project-scoped `.claude/` config. Equivalent: package the skill into a plugin and install it in Cowork, or add it as a skill via claude.ai |

Reference: [Desktop docs](https://code.claude.com/docs/en/desktop#coming-from-the-cli) —
"Desktop and CLI read the same configuration files ... Hooks and skills
defined in settings apply to both."

## Troubleshooting

- **`/board-tests` not found** → the file must be exactly
  `.claude/skills/<name>/SKILL.md` (uppercase `SKILL.md`); restart the session
  or run `/reload-skills`.
- **Skill never triggers automatically** → your `description` is too vague.
  Rewrite it to describe the *user's situation* ("Use when writing unit
  tests...") rather than the content ("Information about tests").
- **Claude wrote a test with `testify` anyway** → it didn't load the skill;
  check `/context` to see whether the skill is listed, then invoke it
  explicitly with `/board-tests` and retry.
