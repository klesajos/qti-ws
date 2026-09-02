> 🌍 Read this in: **English** | [Česky](exercises.cs.md)

# Exercise catalog: one feature per mechanism

You've worked through all six extension mechanisms ([Skills](01-skills.md),
[Hooks](02-hooks.md), [MCP](03-mcp.md), [Plugins](04-plugins.md),
[Subagents](05-agents.md), [Workflows](06-workflows.md)). This catalog turns
them into practice: **each game feature below is a vehicle for exercising one
specific mechanism.** The point isn't only the feature — it's getting the
mechanism into your hands.

Each exercise is **self-contained**: it gives you a paste-ready **starter
prompt**, the exact **files** you'll touch, the bug (where there is one)
explained inline, a **step-by-step walkthrough with the full solution code**, and
a **done check**. The linked guide is optional depth.

## How to use this catalog

1. **Pick a row** from the at-a-glance table below.
2. **Open the linked guide** for that mechanism if you want the *how* and *why* —
   otherwise skip it; the exercise stands on its own.
3. **`Shift+Tab`** to switch Claude Code into plan or accept-edits mode (so you
   review or auto-apply edits the way you prefer).
4. **Paste the Starter prompt** from the exercise verbatim.
5. **Stuck, or want to compare?** Each exercise ends with a **Walkthrough** — the
   full solution, step by step, with copy-pasteable code.
6. **Run the Done-check command** to confirm it worked.

Every exercise states its **goal**, the **mechanism** it teaches, the **files**
you'll touch, a paste-ready **starter prompt**, and a **done check**. The done
check always includes this baseline — the project still builds, vets clean and
the suite is green:

```bash
go build ./... && go vet ./... && go test ./...
```

**Easiest on-ramp:** this repo ships with **three** real, deliberate bugs (all
three are used in Exercise 5). Fixing them needs almost no new code and is the
fastest way to feel a mechanism work:

- **`SlideLineLeft`** (`internal/board/board.go`, lines 83–101) does not advance
  the scan index `i` after a merge — see the `// NOTE: index i is intentionally
  not advanced here.` at line 97. So a merged tile can merge *again* in the same
  move: `{4, 4, 8, 0}` becomes `{16, 0, 0, 0}` instead of `{8, 8, 0, 0}`. The
  scaffold TODO at `internal/board/board_test.go:62` points straight at it.
- **`IsGameOver()`** (`internal/board/board.go`, lines 187–198) only checks for
  an empty cell, so a full board that still has a mergeable pair is wrongly
  reported as game over. Its scaffold TODO is at
  `internal/board/board_test.go:136`.
- **The discarded `Move` result:** `Game.Run` calls
  `g.board.Move(toDirection(cmd))` (`internal/game/game.go:48`) and throws away
  the `bool` it returns, then calls `g.board.SpawnRandom()` unconditionally
  (line 51) — so a tile spawns even after a move that changed nothing. In C++
  this bug takes the form of an *unused variable*; Go's compiler would reject
  that outright, so here it is a **discarded return value** instead — and
  `go vet ./...` does **not** flag it. Only a human or a reviewing agent will.

Difficulty: 🟢 Easy · 🟡 Medium · 🔴 Hard.

## At a glance

| Mechanism | Feature exercise | Difficulty | Touches |
|---|---|---|---|
| **Migration** | [Port the C++ original to Go yourself](00-migration.md) | 🟡–🔴 | a separate C++ repo → a new Go port |
| **Skill** | [A `renderer-style` skill, then colour tiles + a move counter](#1-skill) | 🟢–🟡 | `internal/renderer/`, `internal/game/` |
| **Hook** | [A `go test` Stop hook that blocks on red, then "configurable win value" behind it](#2-hook) | 🟡 | `.claude/settings.json`, `internal/board/` |
| **MCP** | [High-score persistence via the `memory` MCP server](#3-mcp) | 🟡 | `internal/game/`, `internal/renderer/`, `.mcp.json` |
| **Plugin** | [Package the undo workflow as a bundled `/2048-dev:undo` **skill**](#4-plugin) | 🟡 | `plugins/2048-dev/` |
| **Subagent** | [Undo feature: planned, tested, reviewed by the agents; fix the three bugs](#5-subagent) | 🟢–🟡 | `board`, `game`, `input` |
| **Workflow** | [AI auto-solver: parallel heuristics, benchmarked, judge-picked](#6-workflow) | 🔴 | new `internal/ai/`, `cmd/2048-bench/` |

---

<a id="0-migration"></a>

## 0. Migration 🟡–🔴 — port the C++ original to Go, yourself

This one has its own page: **[Exercise 0: Migrate the C++ original to
Go](00-migration.md)**. This whole repo *is* the answer key — it was produced by
that migration. Doing it yourself exercises subagents, worktrees and hooks at
once, and it is the honest way to find out whether Claude ports code
*faithfully*, bugs and all.

---

<a id="1-skill"></a>

## 1. Skill 🟢–🟡 — teach the renderer's conventions, then use them

**Goal.** Write a new skill `renderer-style` that captures this project's ANSI /
terminal-drawing conventions (the escape codes in `internal/renderer/renderer.go`,
the `cellWidth` alignment, the score header, the raw-mode `\r\n` line endings).
Then, *with the skill active*, add **colour tiles** (a different ANSI colour per
tile value) and a **move counter** shown in the header.

**Mechanism it teaches.** Authoring a [skill](01-skills.md) — a description that
triggers on the right requests, a body that briefs Claude on conventions — and
then watching it load on demand. Contrast with `board-tests`: same shape, new
domain.

**Files.** `.claude/skills/renderer-style/SKILL.md` (new),
`internal/renderer/renderer.go`, `internal/game/game.go` (the move counter lives
on the `Game` struct, **not** on `Board` — package `board` stays I/O-free).

**Inline context.** `internal/renderer/renderer.go` clears the screen with
`\x1b[2J\x1b[H` (line 32), prints the `  2048  —  score:` header (line 33), and
lays out the grid as `cellWidth`-wide right-aligned columns via
`fmt.Fprintf(r.out, "%*d", cellWidth, value)` (lines 13, 41). Every line ends
with the `eol` constant `"\r\n"` (line 17) because the terminal is in raw mode.
`Renderer.Draw` currently takes only `*board.Board`
(`internal/renderer/renderer.go:30`) — to show a move counter you pass the count
in as a new argument, since `Board` has no move count and must not gain I/O.

**Starter prompt.**
> *"Create a skill at `.claude/skills/renderer-style/SKILL.md` that documents the
> terminal-drawing conventions in `internal/renderer/renderer.go`: the
> `\x1b[2J\x1b[H` clear-and-home escape sequence, the `  2048  —  score:` header,
> the `cellWidth`-wide right-aligned columns printed with
> `fmt.Fprintf(r.out, "%*d", cellWidth, value)`, and the `\r\n` line endings that
> raw mode requires. Write the `description` as a trigger for any renderer/ANSI
> output work in this repo."*

Then, once the skill is active:
> *"With the renderer-style skill, add a per-value ANSI colour to each tile and a
> `moves: N` counter in the header. Keep package `board` I/O-free: store the move
> count on the `Game` struct in `internal/game/game.go` and pass it into
> `Renderer.Draw` in `internal/renderer/renderer.go` alongside the board."*

**Walkthrough (full solution).**

1. **Write the skill** at `.claude/skills/renderer-style/SKILL.md`:

```markdown
---
name: renderer-style
description: Terminal-drawing conventions for this 2048 Go repo — ANSI escape codes, the score header, cellWidth column alignment and raw-mode line endings. Use for any work on internal/renderer/renderer.go or other ANSI/terminal output.
---

# Renderer style

This project draws the board to stdout with raw ANSI escape codes, from
`internal/renderer/renderer.go`. Match these conventions in any renderer work:

- **Clear + home:** begin a frame with `\x1b[2J\x1b[H`.
- **Header:** one line — `  2048  —  score: <n>` — then a blank line.
- **Grid:** each tile is a `cellWidth`-wide, right-aligned column printed with
  `fmt.Fprintf(r.out, "%*d", cellWidth, value)`; empty cells print `" ."` via
  `"%*s"`.
- **Line endings:** the terminal is in raw mode while the game runs, so every
  line ends with the `eol` constant (`"\r\n"`), never a bare `"\n"` — a bare
  newline moves down without returning to the left edge.
- **Footer:** the hint `  Arrows / WASD to move,  q to quit`.
- **Colour:** wrap a value in an SGR colour with `\x1b[<code>m … \x1b[0m`;
  always reset with `\x1b[0m` so colour never leaks past the cell.
- **Writer:** write through `r.out` (an `io.Writer`), never `fmt.Printf`, so the
  renderer stays redirectable.
- Keep all of this in `renderer`; `board` stays I/O-free, so anything the
  renderer needs but the board doesn't own (a move count, a best score) is
  passed into `Renderer.Draw` as an extra argument.
```

2. **Add a colour table** in `internal/renderer/renderer.go`, next to the
   `cellWidth` and `eol` constants:

```go
// colorFor returns the ANSI SGR colour for a tile value; bigger tiles get
// hotter colours.
func colorFor(value int) string {
	switch value {
	case 2:
		return "\x1b[37m" // white
	case 4:
		return "\x1b[36m" // cyan
	case 8:
		return "\x1b[32m" // green
	case 16:
		return "\x1b[33m" // yellow
	case 32:
		return "\x1b[35m" // magenta
	case 64:
		return "\x1b[31m" // red
	case 128, 256, 512:
		return "\x1b[94m" // bright blue
	default:
		return "\x1b[91m" // bright red (1024+)
	}
}
```

3. **Widen `Draw` and colour the tiles** (same file) — the signature takes the
   move count, the header shows it, and each tile resets its colour:

```go
// Draw clears the screen and renders the board, score and move counter.
func (r *Renderer) Draw(b *board.Board, moves int) {
	// Clear screen and move the cursor to the top-left corner.
	fmt.Fprint(r.out, "\x1b[2J\x1b[H")
	fmt.Fprintf(r.out, "  2048  —  score: %d   moves: %d%s%s", b.Score(), moves, eol, eol)

	for row := 0; row < board.Size; row++ {
		for col := 0; col < board.Size; col++ {
			value := b.At(row, col)
			if value == 0 {
				fmt.Fprintf(r.out, "%*s", cellWidth, " .")
			} else {
				// Colour the tile, then reset so colour never leaks past it.
				fmt.Fprintf(r.out, "%s%*d\x1b[0m", colorFor(value), cellWidth, value)
			}
		}
		fmt.Fprint(r.out, eol, eol)
	}

	fmt.Fprint(r.out, "  Arrows / WASD to move,  q to quit", eol)
}
```

4. **Store the count on `Game`, not `Board`** — add a field to the struct in
   `internal/game/game.go`:

```go
// Game owns one board and the terminal I/O around it.
type Game struct {
	board    *board.Board
	renderer *renderer.Renderer
	input    *input.Input
	moves    int // moves that changed the board; shown in the header
}
```

5. **Count real moves and pass them in** (`internal/game/game.go`, `Run`) —
   increment only when the board actually changed:

```go
	g.renderer.Draw(g.board, g.moves)

	for {
		cmd := g.input.Next()
		if cmd == input.Quit {
			break
		}
		if cmd == input.None {
			continue
		}

		if g.board.Move(toDirection(cmd)) {
			g.moves++
		}

		// Place a new tile and redraw.
		g.board.SpawnRandom()
		g.renderer.Draw(g.board, g.moves)
```

   Using the `bool` from `Board.Move` here is also the first half of the
   discarded-return-value bug Exercise 5 names — guard `SpawnRandom` with the
   same condition and it's fixed outright.

**Done check.** `go build ./... && go vet ./... && go test ./...` green;
`go run ./cmd/2048` renders coloured tiles and a `moves: N` counter. (Keep
rendering out of package `board` — colours and the move count flow through
`renderer` / `game`.)

<a id="2-hook"></a>

## 2. Hook 🟡 — make tests un-skippable, then add a gated feature

**Goal.** Add a hook that runs `go test ./...` and **blocks on red** — a `Stop`
hook so a session can't end with failing tests. Then build a **configurable win
value** (e.g. play to 1024 or 4096 instead of 2048) and let the hook guarantee
you never finish with a broken suite.

**Mechanism it teaches.** A [hook](02-hooks.md) that enforces something *every
time*, with no model discretion — the complement to the existing `gofmt` format
hook.

**Files.** `.claude/hooks/gate-tests.sh` (new), `.claude/settings.json` (register
it under `hooks.Stop`), `internal/board/board.go` (`WinValue` becomes
configurable), `internal/game/game.go`, `cmd/2048/main.go`,
`internal/board/board_test.go`.

**Inline context.** `.claude/hooks/run-tests.sh` already runs `go test ./...` on
`Stop`, but it's *advisory*: it prints a pass/fail line and always `exit 0`. A
blocking hook instead prints `{"decision": "block", "reason": "..."}` on stdout
when tests fail. The win value is hard-coded as `const WinValue = 2048` at
`internal/board/board.go:14`. Go has no default arguments, so the idiomatic way
to make it configurable without breaking `New()` and `FromGrid()` callers is the
**functional-option** pattern.

**Starter prompt.**
> *"Create `.claude/hooks/gate-tests.sh` and `chmod +x` it. Model it on the
> existing `.claude/hooks/run-tests.sh`, but instead of always `exit 0`, make it
> **block**: when `go test ./...` fails, print
> `{"decision": "block", "reason": "tests are red"}` on stdout so a `Stop` can't
> end the session on red. Register it under `hooks.Stop` in
> `.claude/settings.json` alongside the existing `run-tests.sh` entry."*

Then, with the gate in place:
> *"Make the win value configurable. It is `const WinValue = 2048` at
> `internal/board/board.go:14` — let the game be played to 1024 or 4096 instead.
> Use the functional-option pattern (`board.Option`, `board.WithWinValue(n)`,
> `board.DefaultWinValue`) so `board.New()` and `board.FromGrid(grid, score)`
> keep working, thread the option through `game.New` and take the value from
> `os.Args[1]` in `cmd/2048/main.go`. Add no I/O to package `board`."*

**Walkthrough (full solution).**

1. **Write the blocking hook** at `.claude/hooks/gate-tests.sh`, then
   `chmod +x .claude/hooks/gate-tests.sh`:

```bash
#!/usr/bin/env bash
# Stop hook: run the suite and BLOCK the Stop when it is red, so a session can't
# end on failing tests. Models run-tests.sh, but emits a block decision instead
# of always exiting 0.
set -euo pipefail

cat >/dev/null  # drain the Stop-event JSON on stdin

project_dir="${CLAUDE_PROJECT_DIR:-$(pwd)}"

# Not a Go module here, or no toolchain — let the Stop through.
if [ ! -f "$project_dir/go.mod" ] || ! command -v go >/dev/null 2>&1; then
    exit 0
fi

if (cd "$project_dir" && go test ./... >/dev/null 2>&1); then
    exit 0
fi

# Tests are red: block the Stop and tell Claude how to proceed.
echo '{"decision": "block", "reason": "Tests are failing — run go test ./... and fix them before stopping."}'
exit 0
```

2. **Register it** under `hooks.Stop` in `.claude/settings.json`, alongside the
   existing advisory `run-tests.sh` (keep its `timeout`):

```json
    "Stop": [
      {
        "hooks": [
          {
            "type": "command",
            "command": "${CLAUDE_PROJECT_DIR}/.claude/hooks/run-tests.sh",
            "timeout": 120
          },
          {
            "type": "command",
            "command": "${CLAUDE_PROJECT_DIR}/.claude/hooks/gate-tests.sh"
          }
        ]
      }
    ]
```

3. **Turn the constant into a per-board field** (`internal/board/board.go`) —
   rename the constant to a default and add the option type:

```go
// DefaultWinValue is the tile value that counts as a win unless a board is
// created with WithWinValue.
const DefaultWinValue = 2048
```

```go
// Board is the game state: the grid, the score and the RNG used to spawn
// tiles. Create one with New or FromGrid.
type Board struct {
	grid     Grid
	score    int
	winValue int
	rng      *rand.Rand
}

// Option configures a Board at construction time.
type Option func(*Board)

// WithWinValue sets the tile value that counts as a win (e.g. 1024 or 4096).
func WithWinValue(value int) Option {
	return func(b *Board) { b.winValue = value }
}

// New starts an empty board and seeds it with two random tiles. The target
// tile defaults to DefaultWinValue; pass WithWinValue to change it.
func New(opts ...Option) *Board {
	b := &Board{winValue: DefaultWinValue, rng: newRNG()}
	for _, opt := range opts {
		opt(b)
	}
	b.SpawnRandom()
	b.SpawnRandom()
	return b
}

// FromGrid builds a board from an explicit grid (handy for tests). No tiles
// are spawned, so the state is fully deterministic.
func FromGrid(grid Grid, score int, opts ...Option) *Board {
	b := &Board{grid: grid, score: score, winValue: DefaultWinValue, rng: newRNG()}
	for _, opt := range opts {
		opt(b)
	}
	return b
}

// WinValue returns the tile value that counts as a win for this board.
func (b *Board) WinValue() int { return b.winValue }
```

   Variadic options mean every existing `board.New()` and
   `board.FromGrid(grid, 0)` call site keeps compiling untouched.

4. **Use the field in `HasWon`** (same file):

```go
// HasWon reports whether any tile has reached the board's win value.
func (b *Board) HasWon() bool {
	for _, row := range b.grid {
		for _, value := range row {
			if value >= b.winValue {
				return true
			}
		}
	}
	return false
}
```

5. **Thread the options through `Game`** (`internal/game/game.go`):

```go
// New creates a game with a freshly seeded board and switches the terminal
// into raw mode. Board options (e.g. board.WithWinValue) are passed through.
// Call Close when done.
func New(opts ...board.Option) (*Game, error) {
	in, err := input.New()
	if err != nil {
		return nil, err
	}
	return &Game{board: board.New(opts...), renderer: renderer.New(), input: in}, nil
}
```

   …and report the real target when the player wins:

```go
		if g.board.HasWon() {
			g.renderer.Message(fmt.Sprintf("You reached %d! Keep going or press q.", g.board.WinValue()))
		}
```

6. **Pass it in from the command line** (`cmd/2048/main.go`):

```go
// Command 2048 runs the terminal version of the game.
//
// Usage: 2048 [win-value]
// The optional first argument overrides the target tile, e.g. `2048 1024`.
package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/klesajos/qti-ws/internal/board"
	"github.com/klesajos/qti-ws/internal/game"
)

func main() {
	var opts []board.Option
	if len(os.Args) > 1 {
		winValue, err := strconv.Atoi(os.Args[1])
		if err != nil || winValue < 4 {
			fmt.Fprintf(os.Stderr, "2048: win value must be a number >= 4, got %q\n", os.Args[1])
			os.Exit(2)
		}
		opts = append(opts, board.WithWinValue(winValue))
	}

	g, err := game.New(opts...)
	if err != nil {
		fmt.Fprintln(os.Stderr, "2048: cannot switch the terminal to raw mode:", err)
		os.Exit(1)
	}
	defer g.Close()

	g.Run()
}
```

7. **Add a regression test** (`internal/board/board_test.go`, next to the other
   `HasWon` test):

```go
func TestHasWon_RespectsACustomWinValue(t *testing.T) {
	b := FromGrid(Grid{{1024, 0, 0, 0}}, 0, WithWinValue(1024))
	if !b.HasWon() {
		t.Error("HasWon() = false with WithWinValue(1024) and a 1024 tile, want true")
	}
}
```

   Now `go run ./cmd/2048 1024` wins at 1024, and with red tests the `Stop` hook
   refuses to let the session end.

**Done check.** Deliberately break a test → the `Stop` hook blocks. Fix it → the
session ends cleanly. `go build ./... && go vet ./... && go test ./...` green
with the new win value.

<a id="3-mcp"></a>

## 3. MCP 🟡 — persist the high score across sessions

**Goal.** Use the already-configured **`memory`** MCP server (see `.mcp.json`) to
store and retrieve the best score, so it survives across runs. Claude
reads/writes the high score through the MCP server; the game displays `best: N`.

**Mechanism it teaches.** Giving Claude a [capability](03-mcp.md) it doesn't have
natively (durable memory) via an MCP server — and wiring its data into the app.

**Files.** `internal/game/game.go` (read the best score at start-up, report it at
the end), `internal/renderer/renderer.go` (the header that shows `best: N`),
`.mcp.json` (already has `memory`; confirm it's connected with `/mcp`).

**Inline context.** `.mcp.json` already registers a `memory` stdio server
(`@modelcontextprotocol/server-memory`). The score header is drawn in
`internal/renderer/renderer.go` (line 33) and `Renderer.Draw` takes only
`*board.Board` (line 30) — so the best score has to be passed into
`Renderer.Draw` to appear in the header.

**Starter prompt.**
> *"Run `/mcp` and confirm the `memory` server is connected (it's already in
> `.mcp.json`). Then use that server to persist a best score across sessions: in
> `internal/game/game.go` read the stored best from the `BEST_SCORE` environment
> variable at start-up and print the final best when the game ends, so Claude can
> write it back to `memory`. Surface it as `best: N` in the header — the header
> lives in `internal/renderer/renderer.go`, so pass the best score into
> `Renderer.Draw` alongside the board."*

**Walkthrough (full solution).**

Persistence lives in the `memory` MCP server (Claude's durable store); the Go
side only *shows* the best and *reports* the final score so Claude can save it.
The two meet at an environment variable, `BEST_SCORE`.

1. **Confirm the server.** Run `/mcp` and check `memory` is connected (already in
   `.mcp.json`).

2. **Show `best: N` in the header** (`internal/renderer/renderer.go`):

```go
// Draw clears the screen and renders the board, score and best score.
func (r *Renderer) Draw(b *board.Board, best int) {
	// Clear screen and move the cursor to the top-left corner.
	fmt.Fprint(r.out, "\x1b[2J\x1b[H")
	fmt.Fprintf(r.out, "  2048  —  score: %d   best: %d%s%s", b.Score(), best, eol, eol)
```

3. **Seed the best from the environment** (`internal/game/game.go`) — a new
   field plus two imports:

```go
import (
	"fmt"
	"os"
	"strconv"

	"github.com/klesajos/qti-ws/internal/board"
	"github.com/klesajos/qti-ws/internal/input"
	"github.com/klesajos/qti-ws/internal/renderer"
)

// Game owns one board and the terminal I/O around it.
type Game struct {
	board    *board.Board
	renderer *renderer.Renderer
	input    *input.Input
	best     int // best score carried in from a previous session
}

// New creates a game with a freshly seeded board and switches the terminal
// into raw mode. The previous best score is read from the BEST_SCORE
// environment variable (0 if unset or malformed). Call Close when done.
func New() (*Game, error) {
	in, err := input.New()
	if err != nil {
		return nil, err
	}
	best, _ := strconv.Atoi(os.Getenv("BEST_SCORE")) // malformed -> 0, which is fine here
	return &Game{board: board.New(), renderer: renderer.New(), input: in, best: best}, nil
}
```

4. **Track the running best and report it at the end** (`internal/game/game.go`,
   `Run`) — Go 1.21+ has a builtin `max`, so no helper is needed:

```go
	g.renderer.Draw(g.board, max(g.best, g.board.Score()))

	for {
		cmd := g.input.Next()
		if cmd == input.Quit {
			break
		}
		if cmd == input.None {
			continue
		}

		g.board.Move(toDirection(cmd))

		// Place a new tile and redraw.
		g.board.SpawnRandom()
		best := max(g.best, g.board.Score())
		g.renderer.Draw(g.board, best)

		if g.board.HasWon() {
			g.renderer.Message("You reached 2048! Keep going or press q.")
		}
		if g.board.IsGameOver() {
			finalBest := max(g.best, g.board.Score())
			g.renderer.Message(fmt.Sprintf("Game over. Final score: %d  (best: %d)", g.board.Score(), finalBest))
			break
		}
	}
```

5. **Wire MCP to the env var** — this is the workshop's point, with Claude as the
   bridge:
   - *First run:* `BEST_SCORE=0 go run ./cmd/2048`. When it prints the final
     best, ask Claude *"store my 2048 best score in the memory server"* — Claude
     calls the `memory` server (an entity `2048-best-score` with the value as an
     observation).
   - *Next run:* ask *"what's my 2048 best score?"* — Claude reads it back from
     `memory` — then launch `BEST_SCORE=<that value> go run ./cmd/2048` and the
     header shows `best: N` carried across sessions.

**Done check.** Play a game, reach a score, start a new session → the best score
persists and shows as `best: N`.
`go build ./... && go vet ./... && go test ./...` green.

<a id="4-plugin"></a>

## 4. Plugin 🟡 — package the undo workflow as a bundled skill

**Goal.** Extend the `2048-dev` plugin with a new **bundled skill** —
`plugins/2048-dev/skills/undo/SKILL.md` — that drives an undo workflow (keep a
history of board states, revert the last move), shipped alongside the existing
`build-test` skill so the whole team shares the same `/2048-dev:undo`.

**Mechanism it teaches.** Growing a [plugin](04-plugins.md) — adding a bundled
skill under `plugins/2048-dev/skills/`, a version bump, validation — so a whole
team shares the same workflow. **Skills are the modern path**: custom commands
have merged into Skills (the [cheat-sheet](cheatsheet.md) explains the
command→skill convergence), so build new plugin pieces as skills, not legacy
`commands/` files.

**Files.** `plugins/2048-dev/skills/undo/SKILL.md` (new),
`plugins/2048-dev/.claude-plugin/plugin.json` (bump `version` 1.2.0 → 1.3.0).
Validate with `claude plugin validate .`.

**Inline context.** The plugin already bundles one skill at
`plugins/2048-dev/skills/build-test/SKILL.md` and is at `version` `1.2.0` in
`plugins/2048-dev/.claude-plugin/plugin.json`. A skill folder `skills/undo/`
becomes the command `/2048-dev:undo` (the prefix is the plugin name). There is
intentionally no `commands/` directory — use a skill.

**Starter prompt.**
> *"Extend the `2048-dev` plugin with a new bundled skill at
> `plugins/2048-dev/skills/undo/SKILL.md`, modelled on the existing
> `plugins/2048-dev/skills/build-test/SKILL.md`. It should drive the Go undo
> feature: a capped history stack in `internal/board`, a `Board.Undo()` method,
> an `input.Undo` command mapped to `u`, handled in `Game.Run` without spawning a
> tile, verified with `go test ./...`. Bump `version` in
> `plugins/2048-dev/.claude-plugin/plugin.json` from `1.2.0` to `1.3.0`, then run
> `claude plugin validate .`. Use a Skill, not a legacy `commands/` file."*

**Walkthrough (full solution).**

1. **Add the bundled skill** at `plugins/2048-dev/skills/undo/SKILL.md`, modelled
   on the existing `build-test` skill. (The undo *code* is Exercise 5 — this
   skill packages the workflow so the whole team shares `/2048-dev:undo`.)

```markdown
---
name: undo
description: Add or drive the 2048 undo feature — keep a history of prior board states and revert the last move. Use when working on undo/history in this repo.
allowed-tools: Read, Edit, Bash(go:*), Bash(gofmt:*)
disable-model-invocation: true
---

# Undo

Drive the 2048 undo feature: keep a stack of prior board states and revert the
last move on demand.

1. **History.** In `internal/board/board.go`, add an unexported `snapshot`
   struct (grid + score) and a `history []snapshot` field on `Board`. In
   `Board.Move`, capture the grid and score *before* the move and append the
   snapshot only when the move actually changed the board. Cap the depth with a
   `maxHistory` constant so the stack can't grow unbounded.
2. **Revert.** Add `func (b *Board) Undo() bool`: it pops the last snapshot and
   restores grid + score, and returns false when the history is empty.
3. **Key.** Add an `Undo` value to the `Command` constants in
   `internal/input/input.go`, map `'u'` / `'U'` to it in `Input.Next`, and
   handle it in `Game.Run` (`internal/game/game.go`) — undo, redraw, `continue`;
   do **not** spawn a tile.
4. **Verify.** `gofmt -l .` (silent), `go vet ./...`, `go test ./...`. Add tests
   in `internal/board/board_test.go`: a move then an undo restores the exact
   prior grid and score, an undo with no history returns false, and a no-op move
   is not recorded.

Keep `board` I/O-free — history lives in the board state, the keypress handling
stays in `input`/`game`.

Report what changed and the test result.
```

2. **Bump the plugin version** in
   `plugins/2048-dev/.claude-plugin/plugin.json` (`1.2.0` → `1.3.0`):

```json
  "version": "1.3.0",
```

3. **Validate the plugin** so the new skill is well-formed and discoverable:

```bash
claude plugin validate .
```

   It validates the marketplace manifest at `.claude-plugin/marketplace.json`
   and should print `✔ Validation passed`. To check just the plugin manifest,
   run `claude plugin validate ./plugins/2048-dev`.

**Done check.** `claude plugin validate .` passes; `/2048-dev:undo` autocompletes
and loads. `go build ./... && go vet ./... && go test ./...` green.

<a id="5-subagent"></a>

## 5. Subagent 🟢–🟡 — delegate an undo feature across the three agents

**Goal.** Build a real **undo** feature (keep a history of board states; `u`
reverts the last move) by *delegating each part to the right agent from
[Example 5](05-agents.md)*:

- `2048-dev:game-explorer` traces where moves are applied and where history would
  live.
- `board-test-writer` writes the tests for the undo behaviour first.
- `go-reviewer` reviews your diff before you commit.

As a warm-up, fix the **three known bugs** (`SlideLineLeft`, `IsGameOver` and the
discarded `Move` result) using `go-reviewer` to confirm the fix and systematic
debugging to reason it through.

**Mechanism it teaches.** Using [subagents](05-agents.md) as a team: a read-only
cartographer, a write-capable tester, and a read-only reviewer — each in its own
isolated context, each with the right tools.

**Files.** `internal/board/board.go` (history + revert + the two board bugs),
`internal/board/board_test.go` (the two scaffold TODOs + undo tests),
`internal/game/game.go` (handle the undo command, guard the spawn),
`internal/input/input.go` (a new `Command` value and the key that maps to it).

**Inline context (the three bugs).**
- **`SlideLineLeft`** at `internal/board/board.go:83–101` never advances the scan
  index `i` after a merge (`// NOTE: index i is intentionally not advanced here.`
  at line 97), so a freshly merged tile can merge again in the same move:
  `{4, 4, 8, 0}` cascades to `{16, 0, 0, 0}` instead of the correct
  `{8, 8, 0, 0}`. The matching scaffold TODO is at
  `internal/board/board_test.go:62`.
- **`IsGameOver()`** at `internal/board/board.go:187–198` returns `true` as soon
  as it finds no empty cell — it never checks for mergeable neighbours, so a
  full-but-playable board is wrongly "game over". The matching scaffold TODO is
  at `internal/board/board_test.go:136`.
- **The discarded `Move` result**: `Game.Run` calls
  `g.board.Move(toDirection(cmd))` (`internal/game/game.go:48`) and drops the
  `bool`, then calls `g.board.SpawnRandom()` unconditionally (line 51). So a
  no-op move still spawns a tile. In C++ the same bug reads as an *unused
  variable*; Go's compiler would refuse to build that, so here it is a
  **discarded return value** — which `go vet ./...` does **not** flag. This bug
  lives in `Game.Run`, which has **no unit-test harness**.

**Starter prompt — bug A (cascade merge).**
> *"Hand the scaffold TODO at `internal/board/board_test.go:62` to
> `board-test-writer`: add a test for `SlideLineLeft` on the line `{4, 4, 8, 0}`
> and assert the result is `{8, 8, 0, 0}` with 8 points gained — a tile may merge
> at most once per move. Run it and watch it fail: `SlideLineLeft` in
> `internal/board/board.go` does not advance its scan index `i` after a merge
> (line 97), so the line cascades to `{16, 0, 0, 0}`. Fix it, then have
> `go-reviewer` check the diff."*

**Starter prompt — bug B (IsGameOver).**
> *"Hand the scaffold TODO at `internal/board/board_test.go:136` to
> `board-test-writer`: add a test for a full board that still contains a
> mergeable pair and assert `IsGameOver()` is false. Run it and watch it fail —
> `IsGameOver()` in `internal/board/board.go` (lines 187–198) only checks for an
> empty cell. Fix it so a board with equal orthogonal neighbours is not game
> over, then have `go-reviewer` check the diff."*

**Starter prompt — bug C (discarded Move result).**
> *"In `internal/game/game.go`, `Game.Run` calls
> `g.board.Move(toDirection(cmd))` (line 48) and discards the `bool` it returns,
> then calls `g.board.SpawnRandom()` unconditionally (line 51), so a no-op move
> still spawns a tile. `go vet` does not catch a discarded return value — only
> spawn when the move actually changed the board, and ask `go-reviewer` to
> confirm there are no other discarded results in the diff."*

**Walkthrough (full solution).**

*Warm-up — fix the three bugs first.*

A. **`SlideLineLeft`** (`internal/board/board.go`) — advance past the merged tile
   so it can't merge twice in one move:

```go
			out[n-1] = 0
			n--
			// A tile may merge at most once per move: step past the merged tile.
			i++
		} else {
			i++
		}
```

   The test `board-test-writer` writes for the TODO at
   `internal/board/board_test.go:62`:

```go
func TestSlideLineLeft_AMergedTileDoesNotMergeAgain(t *testing.T) {
	line := Line{4, 4, 8, 0}
	gained := SlideLineLeft(&line)
	if want := (Line{8, 8, 0, 0}); line != want {
		t.Errorf("line = %v, want %v (a tile may merge at most once per move)", line, want)
	}
	if gained != 8 {
		t.Errorf("gained = %d, want 8", gained)
	}
}
```

B. **`IsGameOver()`** (`internal/board/board.go`) — after the empty-cell scan,
   also check for equal orthogonal neighbours before declaring the game over:

```go
// IsGameOver reports whether no move is possible.
func (b *Board) IsGameOver() bool {
	// Any empty cell means a move is still possible.
	for r := 0; r < Size; r++ {
		for c := 0; c < Size; c++ {
			if b.grid[r][c] == 0 {
				return false
			}
		}
	}
	// A full board is still playable if two orthogonal neighbours are equal.
	for r := 0; r < Size; r++ {
		for c := 0; c < Size; c++ {
			v := b.grid[r][c]
			if c+1 < Size && b.grid[r][c+1] == v {
				return false
			}
			if r+1 < Size && b.grid[r+1][c] == v {
				return false
			}
		}
	}
	return true
}
```

   …and the test for the TODO at `internal/board/board_test.go:136`:

```go
func TestIsGameOver_FullBoardWithAMergeablePairIsNotOver(t *testing.T) {
	b := FromGrid(Grid{{2, 4, 2, 4}, {4, 2, 4, 2}, {2, 4, 2, 4}, {4, 2, 2, 4}}, 0)
	if b.IsGameOver() {
		t.Error("IsGameOver() = true, want false (the bottom row still has a 2,2 pair)")
	}
}
```

C. **The discarded `Move` result** (`internal/game/game.go`) — only spawn when
   the move actually changed the board (shown below inside the finished loop,
   together with the undo branch).

*Now the undo feature.* `2048-dev:game-explorer` will point you at `Board.Move`
(where a move is applied) as the place to snapshot state.

1. **History storage** (`internal/board/board.go`) — a cap, a snapshot type and a
   field on `Board`:

```go
// maxHistory caps how many moves can be undone so history can't grow unbounded.
const maxHistory = 100

// snapshot is the state captured before a board-changing move.
type snapshot struct {
	grid  Grid
	score int
}

// Board is the game state: the grid, the score and the RNG used to spawn
// tiles. Create one with New or FromGrid.
type Board struct {
	grid    Grid
	score   int
	rng     *rand.Rand
	history []snapshot // snapshots taken before each board-changing move
}
```

2. **Snapshot on change, and revert** (`internal/board/board.go`) — capture the
   score too, push only when the board changed, and add `Undo`:

```go
func (b *Board) Move(dir Direction) bool {
	before := b.grid
	scoreBefore := b.score
	gained := 0
```

```go
	b.score += gained
	changed := b.grid != before
	if changed {
		if len(b.history) == maxHistory {
			b.history = b.history[1:]
		}
		b.history = append(b.history, snapshot{grid: before, score: scoreBefore})
	}
	return changed
}

// Undo reverts to the state before the last board-changing move. It returns
// false when there is nothing left to undo.
func (b *Board) Undo() bool {
	if len(b.history) == 0 {
		return false
	}
	last := b.history[len(b.history)-1]
	b.history = b.history[:len(b.history)-1]
	b.grid, b.score = last.grid, last.score
	return true
}
```

3. **A key for undo** — add the command and map `u` (`internal/input/input.go`):

```go
// The commands a keypress can produce.
const (
	None Command = iota
	Up
	Down
	Left
	Right
	Undo
	Quit
)
```

```go
	case 'u', 'U':
		return Undo
```

4. **Handle it in the loop** (`internal/game/game.go`) — the undo branch plus the
   bug-C fix in one place:

```go
		if cmd == input.None {
			continue
		}
		if cmd == input.Undo {
			g.board.Undo()
			g.renderer.Draw(g.board)
			continue
		}

		// Only spawn a new tile when the move actually changed the board.
		if g.board.Move(toDirection(cmd)) {
			g.board.SpawnRandom()
		}
		g.renderer.Draw(g.board)
```

5. **Test the behaviour** (`internal/board/board_test.go`) —
   `board-test-writer`'s job:

```go
func TestUndo_RestoresTheGridAndScoreBeforeTheLastMove(t *testing.T) {
	start := Grid{{2, 2, 0, 0}}
	b := FromGrid(start, 0)
	b.Move(Left) // -> {4, 0, 0, 0}, score +4
	if got := b.Score(); got != 4 {
		t.Fatalf("Score() after move = %d, want 4", got)
	}

	if !b.Undo() {
		t.Fatal("Undo() = false, want true")
	}
	if got := b.Grid(); got != start {
		t.Errorf("grid after undo = %v, want %v", got, start)
	}
	if got := b.Score(); got != 0 {
		t.Errorf("Score() after undo = %d, want 0", got)
	}
}

func TestUndo_WithNoHistoryReturnsFalse(t *testing.T) {
	b := FromGrid(Grid{{2, 0, 0, 0}}, 0)
	if b.Undo() {
		t.Error("Undo() = true on a fresh board, want false")
	}
}

func TestUndo_ANoOpMoveIsNotRecorded(t *testing.T) {
	b := FromGrid(Grid{{2, 0, 0, 0}}, 0)
	b.Move(Left) // already left-aligned: nothing changes
	if b.Undo() {
		t.Error("Undo() = true after a no-op move, want false")
	}
}
```

   Finally, hand the diff to `go-reviewer` for a read-only pass before you
   commit.

**Done check.** The suite grows from 13 to 18 `Test...` functions and all of them
pass: the two scaffold TODOs (`internal/board/board_test.go:62` and `:136`) are
now real tests that go green, plus three undo tests. The discarded-`Move`-result
fix has **no unit test** — it lives in `Game.Run`, outside the test harness — so
verify it by hand: `go run ./cmd/2048`, press a direction that changes nothing,
and no new tile appears. `go-reviewer` reports no blockers.
`go build ./... && go vet ./... && go test ./...` green.

<a id="6-workflow"></a>

## 6. Workflow 🔴 — an AI auto-solver chosen by a judge panel

**Goal.** Add an AI that plays 2048, built and selected by a
[workflow](06-workflows.md). Implement several heuristics — corner-stacking,
greedy-merge, monotonicity — then write a workflow that **generates/benchmarks
each heuristic in parallel over N games** and uses a **judge panel** to pick the
winner. (Simpler variant: a fan-out "configurable board size" refactor where each
agent updates one file to make `board.Size` a parameter.)

**Mechanism it teaches.** Deterministic
[multi-agent orchestration](06-workflows.md): `parallel()` fan-out, a benchmark
stage, a judge/reduce stage — the generate-and-select pattern.

**Files.** `internal/ai/ai.go` and `internal/ai/ai_test.go` (new — keep the
solver pure, like `board`), `cmd/2048-bench/main.go` (new — the headless harness
the workflow shells out to), `.claude/workflows/solver-benchmark.js` (the new
solver-selection workflow).

**Inline context.** `.claude/workflows/test-coverage-audit.js` is a complete,
working example of the exact shape you need: a `meta` block with phases, a
`parallel(...)` fan-out, per-agent structured-output schemas, and a final reduce
stage. Use it as your starting model — copy its structure and swap the read-only
`Explore` agents for the build/benchmark steps.

**Heads-up on the baseline `IsGameOver` bug.** Self-play calls
`Board.IsGameOver()` every turn, and the baseline version (see
[Exercise 5](#5-subagent)) reports "game over" as soon as the board is full, even
when a merge is still available. Self-play games therefore end a little early and
average scores come out lower than they should. That is fine for *comparing*
heuristics — every heuristic is handicapped identically — but do not read the
absolute numbers as real 2048 scores. Fixing it is Exercise 5's job; this
exercise works either way.

**Starter prompt.**
> *"Use `.claude/workflows/test-coverage-audit.js` as the template (copy its
> `meta`/`phase`/`agent`/structured-schema shape). Write a new workflow
> `.claude/workflows/solver-benchmark.js` that builds several 2048 heuristics —
> corner-stacking, greedy-merge, monotonicity, implemented as a pure
> `internal/ai` package plus a `cmd/2048-bench` harness — runs them in parallel
> with `parallel()` over N self-play games, then adds a judge/reduce stage that
> ranks them and picks a winner."*

**Walkthrough (full solution).**

1. **The pure solver** (`internal/ai/ai.go`) — no I/O, so every heuristic is
   unit-testable. `ChooseMove` tries each direction on a *copy* of the board and
   keeps the best legal one; `PlayGame` drives a whole game headlessly:

```go
// Package ai is a headless auto-solver for 2048. Like board, it is pure: it
// performs no I/O, so every heuristic is unit-testable.
package ai

import (
	"fmt"

	"github.com/klesajos/qti-ws/internal/board"
)

// Heuristic is a strategy the auto-solver can play with.
type Heuristic int

// The heuristics the auto-solver can play with.
const (
	CornerStacking Heuristic = iota
	GreedyMerge
	Monotonicity
)

// String returns the CLI name of the heuristic.
func (h Heuristic) String() string {
	switch h {
	case CornerStacking:
		return "corner-stacking"
	case GreedyMerge:
		return "greedy-merge"
	case Monotonicity:
		return "monotonicity"
	default:
		return "unknown"
	}
}

// ParseHeuristic maps a CLI name to a Heuristic.
func ParseHeuristic(name string) (Heuristic, error) {
	switch name {
	case "corner-stacking":
		return CornerStacking, nil
	case "greedy-merge":
		return GreedyMerge, nil
	case "monotonicity":
		return Monotonicity, nil
	default:
		return 0, fmt.Errorf("unknown heuristic %q (want corner-stacking, greedy-merge or monotonicity)", name)
	}
}

// cornerWeights rewards big tiles pinned to one corner (the classic 2048
// strategy): the weight decreases along a snake from the top-left corner.
var cornerWeights = board.Grid{
	{15, 14, 13, 12},
	{8, 9, 10, 11},
	{7, 6, 5, 4},
	{0, 1, 2, 3},
}

func sumTiles(grid board.Grid) int {
	sum := 0
	for _, row := range grid {
		for _, value := range row {
			sum += value
		}
	}
	return sum
}

func countEmpty(grid board.Grid) int {
	empty := 0
	for _, row := range grid {
		for _, value := range row {
			if value == 0 {
				empty++
			}
		}
	}
	return empty
}

func cornerScore(grid board.Grid) int {
	score := 0
	for r := 0; r < board.Size; r++ {
		for c := 0; c < board.Size; c++ {
			score += grid[r][c] * cornerWeights[r][c]
		}
	}
	return score
}

// greedyMergeScore rewards keeping the board empty (more room == more future
// merges).
func greedyMergeScore(grid board.Grid) int {
	return countEmpty(grid)*100 + sumTiles(grid)
}

// monotonicityScore rewards rows and columns that stay ordered, so equal tiles
// line up to merge.
func monotonicityScore(grid board.Grid) int {
	ordered := 0
	for r := 0; r < board.Size; r++ {
		for c := 0; c+1 < board.Size; c++ {
			if grid[r][c] >= grid[r][c+1] {
				ordered++
			}
		}
	}
	for c := 0; c < board.Size; c++ {
		for r := 0; r+1 < board.Size; r++ {
			if grid[r][c] >= grid[r+1][c] {
				ordered++
			}
		}
	}
	return ordered*100 + countEmpty(grid)*10
}

// Evaluate scores a grid under a heuristic; higher is better.
func Evaluate(grid board.Grid, h Heuristic) int {
	switch h {
	case CornerStacking:
		return cornerScore(grid)
	case GreedyMerge:
		return greedyMergeScore(grid)
	case Monotonicity:
		return monotonicityScore(grid)
	default:
		return 0
	}
}

// ChooseMove picks the legal move the heuristic rates highest. It is pure: it
// looks one move ahead on a copy of the board and performs no I/O. The second
// return value is false when no direction changes the board.
func ChooseMove(b *board.Board, h Heuristic) (board.Direction, bool) {
	var best board.Direction
	bestScore := 0
	found := false

	for _, dir := range []board.Direction{board.Up, board.Down, board.Left, board.Right} {
		trial := *b // copy: Move mutates the receiver
		if !trial.Move(dir) {
			continue // illegal: the board did not change
		}
		score := Evaluate(trial.Grid(), h)
		if !found || score > bestScore {
			found = true
			bestScore = score
			best = dir
		}
	}
	return best, found
}

// PlayGame plays one headless self-play game with the heuristic and returns
// the final score. Tiles spawn exactly as in the interactive game.
func PlayGame(h Heuristic) int {
	b := board.New()
	for !b.IsGameOver() {
		dir, ok := ChooseMove(b, h)
		if !ok {
			break // no legal move left
		}
		if !b.Move(dir) {
			break
		}
		b.SpawnRandom()
	}
	return b.Score()
}
```

   `trial := *b` copies the whole `Board` struct, so the look-ahead never touches
   the real board — that is what keeps `ChooseMove` pure.

2. **Test the pure solver logic** (`internal/ai/ai_test.go`) — same conventions
   as `internal/board/board_test.go`, standard `testing` only:

```go
package ai

import (
	"testing"

	"github.com/klesajos/qti-ws/internal/board"
)

func TestEvaluate_GreedyMergePrefersAnEmptierBoard(t *testing.T) {
	full := board.Grid{{2, 4, 2, 4}, {4, 2, 4, 2}, {2, 4, 2, 4}, {4, 2, 4, 2}}
	sparse := board.Grid{{2, 0, 0, 0}}

	if Evaluate(sparse, GreedyMerge) <= Evaluate(full, GreedyMerge) {
		t.Errorf("Evaluate(sparse) = %d, want > Evaluate(full) = %d",
			Evaluate(sparse, GreedyMerge), Evaluate(full, GreedyMerge))
	}
}

func TestChooseMove_ReturnsALegalMoveForAPlayableBoard(t *testing.T) {
	b := board.FromGrid(board.Grid{{2, 2, 0, 0}}, 0)

	dir, ok := ChooseMove(b, GreedyMerge)
	if !ok {
		t.Fatal("ChooseMove() ok = false, want true")
	}
	if !b.Move(dir) {
		t.Errorf("Move(%v) = false, want true (the chosen move must change the board)", dir)
	}
}
```

3. **A headless benchmark binary** (`cmd/2048-bench/main.go`) the workflow can
   shell out to:

```go
// Command 2048-bench is a headless benchmark: it plays N self-play games with
// one heuristic and prints the average and best score.
//
// Usage: 2048-bench <heuristic> <games>
//
//	heuristic = corner-stacking | greedy-merge | monotonicity
package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/klesajos/qti-ws/internal/ai"
)

func main() {
	h := ai.CornerStacking
	games := 50

	if len(os.Args) > 1 {
		parsed, err := ai.ParseHeuristic(os.Args[1])
		if err != nil {
			fmt.Fprintln(os.Stderr, "2048-bench:", err)
			os.Exit(2)
		}
		h = parsed
	}
	if len(os.Args) > 2 {
		parsed, err := strconv.Atoi(os.Args[2])
		if err != nil || parsed < 1 {
			fmt.Fprintf(os.Stderr, "2048-bench: games must be a positive number, got %q\n", os.Args[2])
			os.Exit(2)
		}
		games = parsed
	}

	total := 0
	best := 0
	for i := 0; i < games; i++ {
		score := ai.PlayGame(h)
		total += score
		if score > best {
			best = score
		}
	}

	fmt.Printf("heuristic=%s games=%d avg=%.1f best=%d\n",
		h, games, float64(total)/float64(games), best)
}
```

   Go needs no build-file edit: `cmd/2048-bench/` is picked up by `go build ./...`
   automatically, and `%s` on `h` uses the `String()` method defined above.

4. **The selection workflow** (`.claude/workflows/solver-benchmark.js`) — copied
   from `test-coverage-audit.js`'s shape: a `meta` block, a `parallel()` fan-out
   with per-agent schemas, and a judge/reduce stage:

```js
// solver-benchmark — build the 2048 auto-solver, benchmark every heuristic in
// parallel, and let a judge stage rank them.
//
// It is the generate-and-select pattern: one build step, a parallel() fan-out
// with one agent per heuristic, and a reduce stage that turns the numbers into
// a ranking. Modelled on test-coverage-audit.js.
//
// Run it from Claude Code with:  /solver-benchmark

export const meta = {
  name: 'solver-benchmark',
  description:
    'Build the 2048 auto-solver, benchmark each heuristic over N self-play games, and judge-pick the strongest.',
  phases: [
    { title: 'Build', detail: 'build the packages and run the ai unit tests once' },
    { title: 'Benchmark', detail: 'one agent per heuristic runs N self-play games' },
    { title: 'Judge', detail: 'rank heuristics by score and pick a winner' },
  ],
}

// The heuristics in internal/ai/ai.go, under the CLI names ParseHeuristic
// accepts.
const HEURISTICS = [
  { key: 'corner-stacking', note: 'pin the biggest tile to a corner' },
  { key: 'greedy-merge', note: 'keep the board as empty as possible' },
  { key: 'monotonicity', note: 'keep rows and columns ordered' },
]

const GAMES = 200 // self-play games per heuristic

// ── Structured-output schemas ──────────────────────────────────────────────
// Each agent is forced to return data matching its schema, so the script never
// has to parse free text.

const RESULT_SCHEMA = {
  type: 'object',
  properties: {
    heuristic: { type: 'string' },
    games: { type: 'number' },
    avg: { type: 'number', description: 'average final score across the games' },
    best: { type: 'number', description: 'best single-game score' },
  },
  required: ['heuristic', 'avg', 'best'],
}

const RANKING_SCHEMA = {
  type: 'object',
  properties: {
    winner: { type: 'string', description: 'the strongest heuristic key' },
    ranking: {
      type: 'array',
      items: {
        type: 'object',
        properties: {
          heuristic: { type: 'string' },
          rank: { type: 'number' },
          avg: { type: 'number' },
          rationale: { type: 'string' },
        },
        required: ['heuristic', 'rank', 'avg'],
      },
    },
    markdown: { type: 'string', description: 'the rendered comparison table' },
  },
  required: ['winner', 'ranking', 'markdown'],
}

// ── Phase 1: build once and prove the solver compiles ──────────────────────
phase('Build')
await agent(
  `In this 2048 Go repo, run "go build ./... && go test ./internal/ai/". ` +
    `Confirm the build succeeds and the ai tests pass, so ` +
    `"go run ./cmd/2048-bench" will work. Report only OK or the first error ` +
    `with file:line.`,
  { label: 'build:bench', phase: 'Build' }
)

// ── Phase 2: benchmark each heuristic in parallel ──────────────────────────
// A barrier (parallel) is correct: the judge needs every result at once.
phase('Benchmark')
const results = await parallel(
  HEURISTICS.map((h) => () =>
    agent(
      `Run "go run ./cmd/2048-bench ${h.key} ${GAMES}" in this repo (${h.note}). ` +
        `It prints one line like ` +
        `"heuristic=${h.key} games=${GAMES} avg=1823.1 best=3200". ` +
        `Parse and return the heuristic name, games, avg and best as numbers.`,
      { label: `bench:${h.key}`, phase: 'Benchmark', schema: RESULT_SCHEMA }
    )
  )
)
const scores = results.filter(Boolean)

// ── Phase 3: judge the results and pick a winner ───────────────────────────
phase('Judge')
const ranking = await agent(
  `You are judging three 2048 auto-solver heuristics by their benchmark scores. ` +
    `Rank them best-to-worst by average score (break ties with best score), name ` +
    `the winner, and render a Markdown table with columns: rank | heuristic | avg ` +
    `| best. RESULTS:\n${JSON.stringify(scores, null, 2)}`,
  { label: 'judge:rank', phase: 'Judge', schema: RANKING_SCHEMA }
)

log(`Benchmarked ${scores.length} heuristics; winner: ${ranking.winner}.`)
return ranking.markdown
```

   Run it with `/solver-benchmark`. A quick sanity check from the shell —
   `go run ./cmd/2048-bench greedy-merge 20` — already prints a real line for the
   judge to rank:

```text
heuristic=greedy-merge games=20 avg=2004.4 best=5356
```

**Done check.** The auto-solver plays a full game headlessly; the workflow runs
and reports a ranked comparison. `gofmt -l .` prints nothing and
`go build ./... && go vet ./... && go test ./...` is green, including the new
`internal/ai` tests.

---

Pick any row, paste its starter prompt, and build it. Each one leaves you with a
working feature *and* a mechanism you've actually driven yourself.
