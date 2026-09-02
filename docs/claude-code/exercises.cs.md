> 🌍 Číst v jazyce: [English](exercises.md) | **Česky**

# Katalog cvičení: jedna featura na každý mechanismus

Prošel jsi všech šest mechanismů rozšíření ([Skills](01-skills.cs.md),
[Hooks](02-hooks.cs.md), [MCP](03-mcp.cs.md), [Pluginy](04-plugins.cs.md),
[Subagenti](05-agents.cs.md), [Workflows](06-workflows.cs.md)). Tenhle katalog
je promění v praxi: **každá herní featura níže je prostředkem k procvičení
jednoho konkrétního mechanismu.** Nejde jen o featuru — jde o to dostat
mechanismus do ruky.

Každé cvičení je **soběstačné**: dává ti připravený **startovací prompt**, přesné
**soubory**, na které sáhneš, **bug** (kde nějaký je) vysvětlený přímo na místě,
**krok-za-krokem návod s kompletním řešením** a **kontrolu hotovo**. Odkazovaný
návod je nepovinná hloubka.

## Jak používat tento katalog

1. **Vyber řádek** z přehledové tabulky níže.
2. **Otevři odpovídající návod** k danému mechanismu, pokud chceš *jak* a *proč* —
   jinak ho přeskoč; cvičení obstojí samo.
3. **`Shift+Tab`** přepne Claude Code do plan nebo accept-edits režimu (podle
   toho, jak chceš edity revidovat nebo automaticky aplikovat).
4. **Vlož startovací prompt** ze cvičení doslova.
5. **Zaseklý, nebo chceš porovnat?** Každé cvičení končí **walkthroughem** —
   kompletní řešení, krok za krokem, s kódem k vložení.
6. **Spusť příkaz kontroly hotovo** pro potvrzení, že to funguje.

Každé cvičení uvádí svůj **cíl**, **mechanismus**, který učí, **soubory**, na
které sáhneš, připravený **startovací prompt** a **kontrolu hotovo**. Kontrola
hotovo vždy zahrnuje tento základ — projekt se pořád sestaví, projde vet
kontrolou a sada je zelená:

```bash
go build ./... && go vet ./... && go test ./...
```

**Nejsnazší rozjezd:** tohle repo přichází se **třemi** skutečnými, záměrnými
bugy (všechny tři jsou použité ve Cvičení 5). Jejich oprava skoro nepotřebuje
nový kód a je nejrychlejší způsob, jak mechanismus pocítit:

- **`SlideLineLeft`** (`internal/board/board.go`, řádky 83–101) po sloučení
  neposune skenovací index `i` — viz komentář `// NOTE: index i is intentionally
  not advanced here.` na řádku 97. Sloučená dlaždice se tak může v rámci
  stejného tahu sloučit *znovu*: `{4, 4, 8, 0}` se stane `{16, 0, 0, 0}` místo
  `{8, 8, 0, 0}`. Scaffold TODO na `internal/board/board_test.go:62` na to míří
  přímo.
- **`IsGameOver()`** (`internal/board/board.go`, řádky 187–198) kontroluje jen
  prázdnou buňku, takže plná deska, která pořád má mergeable pár, je chybně
  hlášena jako game over. Odpovídající scaffold TODO je na
  `internal/board/board_test.go:136`.
- **Zahozený výsledek `Move`:** `Game.Run` zavolá
  `g.board.Move(toDirection(cmd))` (`internal/game/game.go:48`) a vrácenou
  hodnotu `bool` zahodí, pak bezpodmínečně zavolá `g.board.SpawnRandom()`
  (řádek 51) — takže dlaždice se objeví i po tahu, který nic nezměnil. V C++
  měl tenhle bug podobu *nepoužité proměnné*; kompilátor Go by to rovnou odmítl
  zkompilovat, takže tady jde místo toho o **zahozenou návratovou hodnotu** — a
  `go vet ./...` ji **neodhalí**. Odhalí ji jen člověk nebo revidující agent.

Obtížnost: 🟢 snadné · 🟡 střední · 🔴 těžké.

## V kostce

| Mechanismus | Cvičení (featura) | Obtížnost | Sahá na |
|---|---|---|---|
| **Migration** | [Portuj C++ originál do Go sám](00-migration.cs.md) | 🟡–🔴 | samostatné C++ repo → nový Go port |
| **Skill** | [Skill `renderer-style`, pak barevné dlaždice + počítadlo tahů](#1-skill) | 🟢–🟡 | `internal/renderer/`, `internal/game/` |
| **Hook** | [Stop hook s `go test`, který zablokuje na červené, pak „konfigurovatelná výherní hodnota" za ním](#2-hook) | 🟡 | `.claude/settings.json`, `internal/board/` |
| **MCP** | [Trvalé high-score přes MCP server `memory`](#3-mcp) | 🟡 | `internal/game/`, `internal/renderer/`, `.mcp.json` |
| **Plugin** | [Undo workflow zabalený jako bundled skill `/2048-dev:undo`](#4-plugin) | 🟡 | `plugins/2048-dev/` |
| **Subagent** | [Undo featura: naplánovaná, otestovaná, zrevidovaná agenty; oprav tři bugy](#5-subagent) | 🟢–🟡 | `board`, `game`, `input` |
| **Workflow** | [AI řešitel: paralelní heuristiky, benchmark, výběr porotou](#6-workflow) | 🔴 | nový `internal/ai/`, `cmd/2048-bench/` |

---

<a id="0-migration"></a>

## 0. Migration 🟡–🔴 — portuj C++ originál do Go sám

Tohle má vlastní stránku: **[Cvičení 0: Migruj C++ originál do
Go](00-migration.cs.md)**. Celé tohle repo *je* odpovědní klíč — vzniklo právě
tou migrací. Udělat si to sám procvičí subagenty, worktrees a hooky najednou a
je to poctivý způsob, jak zjistit, jestli Claude portuje kód *věrně*, bugy
nevyjímaje.

---

<a id="1-skill"></a>

## 1. Skill 🟢–🟡 — nauč konvence rendereru a pak je použij

**Cíl.** Napiš nový skill `renderer-style`, který zachytí ANSI / terminálové
kreslicí konvence projektu (escape kódy v `internal/renderer/renderer.go`,
zarovnání `cellWidth`, hlavičku se skóre a řádkové zakončení `\r\n` v raw módu).
Pak, *s aktivním skillem*, přidej **barevné dlaždice** (jiná ANSI barva podle
hodnoty dlaždice) a **počítadlo tahů** v hlavičce.

**Co učí.** Tvorbu [skillu](01-skills.cs.md) — popis, který se spustí na správné
požadavky, tělo, jež Clauda zasvětí do konvencí — a sledování, jak se načte na
vyžádání. Srovnej s `board-tests`: stejný tvar, nová doména.

**Soubory.** `.claude/skills/renderer-style/SKILL.md` (nový),
`internal/renderer/renderer.go`, `internal/game/game.go` (počítadlo tahů patří
na strukturu `Game`, **ne** na `Board` — balíček `board` zůstává bez I/O).

**Kontext na místě.** `internal/renderer/renderer.go` maže obrazovku přes
`\x1b[2J\x1b[H` (řádek 32), tiskne hlavičku `  2048  —  score:` (řádek 33) a
vykresluje mřížku jako `cellWidth` široké zprava zarovnané sloupce přes
`fmt.Fprintf(r.out, "%*d", cellWidth, value)` (řádky 13, 41). Každý řádek končí
konstantou `eol` (`"\r\n"`, řádek 17), protože terminál je v raw módu.
`Renderer.Draw` teď bere jen `*board.Board`
(`internal/renderer/renderer.go:30`) — aby zobrazil počítadlo tahů, předáš
počet jako nový argument, protože `Board` žádné počítadlo nemá a nesmí dostat
I/O.

**Startovací prompt.**
> *„Vytvoř skill v `.claude/skills/renderer-style/SKILL.md`, který zdokumentuje
> terminálové kreslicí konvence v `internal/renderer/renderer.go`: escape
> sekvenci `\x1b[2J\x1b[H` pro smazání a návrat kurzoru, hlavičku
> `  2048  —  score:`, `cellWidth` široké zprava zarovnané sloupce tištěné přes
> `fmt.Fprintf(r.out, "%*d", cellWidth, value)` a řádkové zakončení `\r\n`, které
> vyžaduje raw mód. `description` napiš jako spouštěč pro jakoukoli práci
> s rendererem / ANSI výstupem v tomhle repu."*

Pak, jakmile je skill aktivní:
> *„Se skillem renderer-style přidej každé dlaždici ANSI barvu podle hodnoty a do
> hlavičky počítadlo `moves: N`. Balíček `board` drž bez I/O: počet tahů ulož na
> strukturu `Game` v `internal/game/game.go` a předej ho do `Renderer.Draw`
> v `internal/renderer/renderer.go` vedle desky."*

**Walkthrough (kompletní řešení).**

1. **Napiš skill** v `.claude/skills/renderer-style/SKILL.md`:

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

2. **Přidej tabulku barev** v `internal/renderer/renderer.go`, vedle konstant
   `cellWidth` a `eol`:

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

3. **Rozšiř `Draw` a obarvi dlaždice** (stejný soubor) — signatura teď bere
   počet tahů, hlavička ho zobrazí a každá dlaždice po sobě resetuje barvu:

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

4. **Počet ulož na `Game`, ne na `Board`** — přidej pole do struktury
   v `internal/game/game.go`:

```go
// Game owns one board and the terminal I/O around it.
type Game struct {
	board    *board.Board
	renderer *renderer.Renderer
	input    *input.Input
	moves    int // moves that changed the board; shown in the header
}
```

5. **Počítej skutečné tahy a předej je** (`internal/game/game.go`, `Run`) —
   inkrementuj jen když se deska skutečně změnila:

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

   Použití `bool` z `Board.Move` tady je zároveň první polovina bugu se
   zahozenou návratovou hodnotou, který pojmenovává Cvičení 5 — ohraď
   `SpawnRandom` stejnou podmínkou a je rovnou opravený.

**Hotovo.** `go build ./... && go vet ./... && go test ./...` zelené;
`go run ./cmd/2048` vykresluje barevné dlaždice a počítadlo `moves: N`.
(Vykreslování drž mimo balíček `board` — barvy i počet tahů tečou přes
`renderer` / `game`.)

<a id="2-hook"></a>

## 2. Hook 🟡 — udělej testy nepřeskočitelné, pak přidej hlídanou featuru

**Cíl.** Přidej hook, který spustí `go test ./...` a **zablokuje na červené** —
`Stop` hook, aby session nemohla skončit s padajícími testy. Pak postav
**konfigurovatelnou výherní hodnotu** (např. hrát do 1024 nebo 4096 místo 2048)
a nech hook zaručit, že nikdy neskončíš s rozbitou sadou.

**Co učí.** [Hook](02-hooks.cs.md), který vynucuje něco *pokaždé*, bez uvážení
modelu — doplněk ke stávajícímu formátovacímu hooku `gofmt`.

**Soubory.** `.claude/hooks/gate-tests.sh` (nový), `.claude/settings.json`
(registrace pod `hooks.Stop`), `internal/board/board.go` (`WinValue` se stane
konfigurovatelnou), `internal/game/game.go`, `cmd/2048/main.go`,
`internal/board/board_test.go`.

**Kontext na místě.** `.claude/hooks/run-tests.sh` už spouští `go test ./...`
na `Stop`, ale je *poradní*: vypíše řádek pass/fail a vždy `exit 0`. Blokující
hook místo toho vypíše `{"decision": "block", "reason": "..."}` na stdout, když
testy padnou. Výherní hodnota je natvrdo `const WinValue = 2048` na
`internal/board/board.go:14`. Go nemá výchozí hodnoty parametrů, takže
idiomatický způsob, jak ji udělat konfigurovatelnou bez rozbití volajících
`New()` a `FromGrid()`, je vzorec **functional-option**.

**Startovací prompt.**
> *„Vytvoř `.claude/hooks/gate-tests.sh` a `chmod +x` ho. Vymodeluj ho podle
> stávajícího `.claude/hooks/run-tests.sh`, ale místo vždy `exit 0` ho nech
> **blokovat**: když `go test ./...` selže, vypiš
> `{"decision": "block", "reason": "tests are red"}` na stdout, aby `Stop`
> nemohl session ukončit na červené. Zaregistruj ho pod `hooks.Stop`
> v `.claude/settings.json` vedle stávající položky `run-tests.sh`."*

Pak, s hlídačem na místě:
> *„Udělej výherní hodnotu konfigurovatelnou. Je to `const WinValue = 2048` na
> `internal/board/board.go:14` — nech hrát do 1024 nebo 4096 místo toho. Použij
> vzorec functional-option (`board.Option`, `board.WithWinValue(n)`,
> `board.DefaultWinValue`), aby `board.New()` a `board.FromGrid(grid, score)`
> dál fungovaly, protáhni tu možnost skrz `game.New` a vezmi hodnotu z
> `os.Args[1]` v `cmd/2048/main.go`. Do balíčku `board` nepřidávej žádné I/O."*

**Walkthrough (kompletní řešení).**

1. **Napiš blokující hook** v `.claude/hooks/gate-tests.sh`, pak
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

2. **Zaregistruj ho** pod `hooks.Stop` v `.claude/settings.json`, vedle
   stávajícího poradního `run-tests.sh` (zachovej jeho `timeout`):

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

3. **Proměň konstantu na pole desky** (`internal/board/board.go`) — přejmenuj
   konstantu na výchozí hodnotu a přidej typ pro option:

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

   Variadic options znamenají, že každé stávající volání `board.New()` a
   `board.FromGrid(grid, 0)` se dál zkompiluje beze změny.

4. **Použij pole v `HasWon`** (stejný soubor):

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

5. **Protáhni options skrz `Game`** (`internal/game/game.go`):

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

   …a po výhře nahlas skutečný cíl:

```go
		if g.board.HasWon() {
			g.renderer.Message(fmt.Sprintf("You reached %d! Keep going or press q.", g.board.WinValue()))
		}
```

6. **Předej ji z příkazové řádky** (`cmd/2048/main.go`):

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

7. **Přidej regresní test** (`internal/board/board_test.go`, vedle dalšího
   testu `HasWon`):

```go
func TestHasWon_RespectsACustomWinValue(t *testing.T) {
	b := FromGrid(Grid{{1024, 0, 0, 0}}, 0, WithWinValue(1024))
	if !b.HasWon() {
		t.Error("HasWon() = false with WithWinValue(1024) and a 1024 tile, want true")
	}
}
```

   Teď `go run ./cmd/2048 1024` vyhrává na 1024 a s červenými testy `Stop` hook
   odmítne nechat session skončit.

**Hotovo.** Záměrně rozbij test → `Stop` hook zablokuje. Oprav → session skončí
čistě. `go build ./... && go vet ./... && go test ./...` zelené s novou
výherní hodnotou.

<a id="3-mcp"></a>

## 3. MCP 🟡 — uchovej high-score napříč sessionami

**Cíl.** Použij už nakonfigurovaný MCP server **`memory`** (viz `.mcp.json`)
k uložení a načtení nejlepšího skóre, aby přežilo mezi běhy. Claude čte/zapisuje
high-score přes MCP server; hra zobrazí `best: N`.

**Co učí.** Dát Claudovi [schopnost](03-mcp.cs.md), kterou nativně nemá (trvalá
paměť), přes MCP server — a propojit jeho data s aplikací.

**Soubory.** `internal/game/game.go` (načti nejlepší skóre při startu, nahlas
ho na konci), `internal/renderer/renderer.go` (hlavička, která zobrazí
`best: N`), `.mcp.json` (`memory` už tam je; ověř připojení přes `/mcp`).

**Kontext na místě.** `.mcp.json` už registruje `memory` stdio server
(`@modelcontextprotocol/server-memory`). Hlavička se skóre se kreslí
v `internal/renderer/renderer.go` (řádek 33) a `Renderer.Draw` bere jen
`*board.Board` (řádek 30) — takže nejlepší skóre se musí předat do
`Renderer.Draw`, aby se v hlavičce objevilo.

**Startovací prompt.**
> *„Spusť `/mcp` a ověř, že je server `memory` připojený (už je v `.mcp.json`).
> Pak ho použij k uchování nejlepšího skóre napříč sessionami: v
> `internal/game/game.go` načti uložené nejlepší skóre z proměnné prostředí
> `BEST_SCORE` při startu a vypiš finální nejlepší skóre, když hra skončí, aby
> ho Claude mohl zapsat zpět do `memory`. Zobraz ho jako `best: N` v hlavičce —
> hlavička žije v `internal/renderer/renderer.go`, takže předej nejlepší skóre
> do `Renderer.Draw` vedle desky."*

**Walkthrough (kompletní řešení).**

Trvalost žije v MCP serveru `memory` (Claudovo trvalé úložiště); Go strana jen
*zobrazuje* nejlepší skóre a *hlásí* finální skóre, aby ho Claude uložil.
Potkávají se přes proměnnou prostředí `BEST_SCORE`.

1. **Ověř server.** Spusť `/mcp` a zkontroluj, že je `memory` připojený (už je
   v `.mcp.json`).

2. **Zobraz `best: N` v hlavičce** (`internal/renderer/renderer.go`):

```go
// Draw clears the screen and renders the board, score and best score.
func (r *Renderer) Draw(b *board.Board, best int) {
	// Clear screen and move the cursor to the top-left corner.
	fmt.Fprint(r.out, "\x1b[2J\x1b[H")
	fmt.Fprintf(r.out, "  2048  —  score: %d   best: %d%s%s", b.Score(), best, eol, eol)
```

3. **Načti nejlepší z prostředí** (`internal/game/game.go`) — nové pole plus
   dva importy:

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

4. **Sleduj průběžné nejlepší a nahlas ho na konci** (`internal/game/game.go`,
   `Run`) — Go 1.21+ má vestavěný `max`, takže žádný pomocník není potřeba:

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

5. **Propoj MCP s proměnnou prostředí** — to je pointa cvičení, Claude je most:
   - *První běh:* `BEST_SCORE=0 go run ./cmd/2048`. Když vypíše finální
     nejlepší, řekni Claudovi *„ulož moje nejlepší 2048 skóre do memory
     serveru"* — Claude zavolá server `memory` (entita `2048-best-score`
     s hodnotou jako observation).
   - *Další běh:* zeptej se *„jaké je moje nejlepší 2048 skóre?"* — Claude ho
     přečte z `memory` — pak spusť `BEST_SCORE=<ta hodnota> go run ./cmd/2048`
     a hlavička ukáže `best: N` přenesené mezi sessionami.

**Hotovo.** Zahraj hru, dosáhni skóre, spusť novou session → nejlepší skóre
přetrvá a zobrazí se jako `best: N`.
`go build ./... && go vet ./... && go test ./...` zelené.

<a id="4-plugin"></a>

## 4. Plugin 🟡 — zabal undo workflow jako bundled skill

**Cíl.** Rozšiř plugin `2048-dev` o nový **bundled skill** —
`plugins/2048-dev/skills/undo/SKILL.md` — který řídí undo workflow (drž historii
stavů desky, vrať poslední tah), dodaný vedle stávajícího skillu `build-test`,
aby celý tým sdílel stejné `/2048-dev:undo`.

**Co učí.** Růst [pluginu](04-plugins.cs.md) — přidání bundled skillu pod
`plugins/2048-dev/skills/`, bump verze, validace — aby celý tým sdílel stejný
workflow. **Skills jsou moderní cesta**: vlastní příkazy se sloučily do Skills
([cheat-sheet](cheatsheet.cs.md) vysvětluje konvergenci command→skill), takže
nové části pluginu stavěj jako skilly, ne jako legacy `commands/` soubory.

**Soubory.** `plugins/2048-dev/skills/undo/SKILL.md` (nový),
`plugins/2048-dev/.claude-plugin/plugin.json` (bump `version` 1.2.0 → 1.3.0).
Ověř přes `claude plugin validate .`.

**Kontext na místě.** Plugin už balí jeden skill v
`plugins/2048-dev/skills/build-test/SKILL.md` a je na `version` `1.2.0`
v `plugins/2048-dev/.claude-plugin/plugin.json`. Složka skillu `skills/undo/` se
stane příkazem `/2048-dev:undo` (prefix je jméno pluginu). Záměrně tu není žádný
adresář `commands/` — použij skill.

**Startovací prompt.**
> *„Rozšiř plugin `2048-dev` o nový bundled skill v
> `plugins/2048-dev/skills/undo/SKILL.md`, vymodelovaný podle stávajícího
> `plugins/2048-dev/skills/build-test/SKILL.md`. Má řídit Go undo featuru:
> ohraničený history stack v `internal/board`, metodu `Board.Undo()`, příkaz
> `input.Undo` namapovaný na `u`, obsloužený v `Game.Run` bez vytvoření
> dlaždice, ověřený přes `go test ./...`. Zvyš `version`
> v `plugins/2048-dev/.claude-plugin/plugin.json` z `1.2.0` na `1.3.0` a pak
> spusť `claude plugin validate .`. Použij Skill, ne legacy `commands/`
> soubor."*

**Walkthrough (kompletní řešení).**

1. **Přidej bundled skill** v `plugins/2048-dev/skills/undo/SKILL.md`,
   vymodelovaný podle stávajícího skillu `build-test`. (Undo *kód* je Cvičení 5
   — tenhle skill jen zabalí workflow, aby celý tým sdílel `/2048-dev:undo`.)

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

2. **Zvyš verzi pluginu** v
   `plugins/2048-dev/.claude-plugin/plugin.json` (`1.2.0` → `1.3.0`):

```json
  "version": "1.3.0",
```

3. **Zvaliduj plugin**, aby byl nový skill v pořádku a vyhledatelný:

```bash
claude plugin validate .
```

   Zvaliduje manifest marketplace v `.claude-plugin/marketplace.json` a měl by
   vypsat `✔ Validation passed`. Pokud chceš zkontrolovat jen manifest pluginu,
   spusť `claude plugin validate ./plugins/2048-dev`.

**Hotovo.** `claude plugin validate .` projde; `/2048-dev:undo` se napovídá
a načte. `go build ./... && go vet ./... && go test ./...` zelené.

<a id="5-subagent"></a>

## 5. Subagent 🟢–🟡 — deleguj undo featuru mezi tři agenty

**Cíl.** Postav skutečnou **undo** featuru (drž historii stavů desky; `u` vrátí
poslední tah) tak, že *každou část deleguješ správnému agentovi z
[Ukázky 5](05-agents.cs.md)*:

- `2048-dev:game-explorer` vystopuje, kde se tahy aplikují a kde by žila
  historie.
- `board-test-writer` nejdřív napíše testy pro undo chování.
- `go-reviewer` zreviduje tvůj diff, než zacommituješ.

Jako rozcvičku oprav **tři známé bugy** (`SlideLineLeft`, `IsGameOver` a
zahozený výsledek `Move`) s pomocí `go-reviewer` pro potvrzení opravy a
systematického debugování pro promyšlení.

**Co učí.** Použití [subagentů](05-agents.cs.md) jako tým: read-only kartograf,
zapisující tester a read-only recenzent — každý ve svém izolovaném kontextu,
každý se správnými nástroji.

**Soubory.** `internal/board/board.go` (historie + revert + dva bugy desky),
`internal/board/board_test.go` (dva scaffold TODOs + undo testy),
`internal/game/game.go` (obsluha undo příkazu, ohraď spawn),
`internal/input/input.go` (nová hodnota `Command` a klávesa, na kterou se
mapuje).

**Kontext na místě (tři bugy).**
- **`SlideLineLeft`** na `internal/board/board.go:83–101` po sloučení nikdy
  neposune skenovací index `i` (komentář `// NOTE: index i is intentionally not
  advanced here.` na řádku 97), takže čerstvě sloučená dlaždice se může v rámci
  stejného tahu sloučit znovu: `{4, 4, 8, 0}` se sesype na `{16, 0, 0, 0}` místo
  správného `{8, 8, 0, 0}`. Odpovídající scaffold TODO je na
  `internal/board/board_test.go:62`.
- **`IsGameOver()`** na `internal/board/board.go:187–198` vrátí `true`, jakmile
  nenajde prázdnou buňku — nikdy nekontroluje mergeable sousedy, takže plná, ale
  hratelná deska je chybně „game over". Odpovídající scaffold TODO je na
  `internal/board/board_test.go:136`.
- **Zahozený výsledek `Move`**: `Game.Run` zavolá
  `g.board.Move(toDirection(cmd))` (`internal/game/game.go:48`) a `bool`
  zahodí, pak bezpodmínečně zavolá `g.board.SpawnRandom()` (řádek 51). Takže
  tah bez efektu pořád objeví dlaždici. V C++ se stejný bug projevuje jako
  *nepoužitá proměnná*; kompilátor Go by odmítl to zkompilovat, takže tady jde
  o **zahozenou návratovou hodnotu** — kterou `go vet ./...` **neodhalí**.
  Tenhle bug žije v `Game.Run`, který **nemá testovací harness**.

**Startovací prompt — bug A (kaskádové sloučení).**
> *„Předej scaffold TODO na `internal/board/board_test.go:62` agentovi
> `board-test-writer`: přidej test pro `SlideLineLeft` na řádku `{4, 4, 8, 0}` a
> vyžaduj výsledek `{8, 8, 0, 0}` se ziskem 8 bodů — dlaždice se může v rámci
> jednoho tahu sloučit nejvýš jednou. Spusť ho a sleduj, jak padne:
> `SlideLineLeft` v `internal/board/board.go` po sloučení neposune svůj
> skenovací index `i` (řádek 97), takže se řádek sesype na `{16, 0, 0, 0}`.
> Oprav to a pak nech `go-reviewer` zkontrolovat diff."*

**Startovací prompt — bug B (IsGameOver).**
> *„Předej scaffold TODO na `internal/board/board_test.go:136` agentovi
> `board-test-writer`: přidej test pro plnou desku, která pořád obsahuje
> mergeable pár, a vyžaduj, aby `IsGameOver()` bylo false. Spusť ho a sleduj,
> jak padne — `IsGameOver()` v `internal/board/board.go` (řádky 187–198)
> kontroluje jen prázdnou buňku. Oprav ji tak, aby deska se stejnými
> ortogonálními sousedy nebyla game over, a pak nech `go-reviewer` zkontrolovat
> diff."*

**Startovací prompt — bug C (zahozený výsledek Move).**
> *„V `internal/game/game.go` `Game.Run` zavolá
> `g.board.Move(toDirection(cmd))` (řádek 48) a vrácený `bool` zahodí, pak
> bezpodmínečně zavolá `g.board.SpawnRandom()` (řádek 51), takže tah bez efektu
> pořád objeví dlaždici. `go vet` zahozenou návratovou hodnotu neodhalí —
> dlaždici objev jen když tah desku skutečně změnil, a požádej `go-reviewer`,
> aby ověřil, že v diffu nejsou žádné další zahozené výsledky."*

**Walkthrough (kompletní řešení).**

*Rozcvička — nejdřív oprav tři bugy.*

A. **`SlideLineLeft`** (`internal/board/board.go`) — posuň se za sloučenou
   dlaždici, aby se nemohla v jednom tahu sloučit dvakrát:

```go
			out[n-1] = 0
			n--
			// A tile may merge at most once per move: step past the merged tile.
			i++
		} else {
			i++
		}
```

   Test, který `board-test-writer` napíše pro TODO na
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

B. **`IsGameOver()`** (`internal/board/board.go`) — po skenu prázdných buněk
   navíc zkontroluj stejné ortogonální sousedy, než prohlásíš konec hry:

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

   …a test pro TODO na `internal/board/board_test.go:136`:

```go
func TestIsGameOver_FullBoardWithAMergeablePairIsNotOver(t *testing.T) {
	b := FromGrid(Grid{{2, 4, 2, 4}, {4, 2, 4, 2}, {2, 4, 2, 4}, {4, 2, 2, 4}}, 0)
	if b.IsGameOver() {
		t.Error("IsGameOver() = true, want false (the bottom row still has a 2,2 pair)")
	}
}
```

C. **Zahozený výsledek `Move`** (`internal/game/game.go`) — objev dlaždici jen
   když tah desku skutečně změnil (ukázáno níže uvnitř dokončené smyčky, spolu
   s undo větví).

*Teď undo featura.* `2048-dev:game-explorer` tě navede na `Board.Move` (kde se
tah aplikuje) jako místo, kam stav uložit.

1. **Úložiště historie** (`internal/board/board.go`) — limit, typ snapshotu a
   pole na `Board`:

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

2. **Snímek při změně a revert** (`internal/board/board.go`) — ulož i skóre,
   pushuj jen když se deska změnila, a přidej `Undo`:

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

3. **Klávesa pro undo** — přidej příkaz a namapuj `u`
   (`internal/input/input.go`):

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

4. **Obsluž ho ve smyčce** (`internal/game/game.go`) — undo větev plus oprava
   bugu C na jednom místě:

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

5. **Otestuj chování** (`internal/board/board_test.go`) — práce pro
   `board-test-writer`:

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

   Nakonec předej diff `go-reviewer`ovi na read-only kontrolu, než zacommituješ.

**Hotovo.** Sada naroste z 13 na 18 funkcí `Test...` a všechny projdou: dva
scaffold TODOs (`internal/board/board_test.go:62` a `:136`) jsou teď skutečné
testy, které zezelenají, plus tři undo testy. Oprava zahozeného výsledku
`Move` **nemá unit test** — žije v `Game.Run`, mimo testovací harness — takže
ji ověř ručně: `go run ./cmd/2048`, stiskni směr, který nic nezmění, a
neobjeví se žádná nová dlaždice. `go-reviewer` nehlásí žádné blockery.
`go build ./... && go vet ./... && go test ./...` zelené.

<a id="6-workflow"></a>

## 6. Workflow 🔴 — AI řešitel vybraný porotou

**Cíl.** Přidej AI, která hraje 2048, postavenou a vybranou
[workflow](06-workflows.cs.md). Naimplementuj několik heuristik —
corner-stacking, greedy-merge, monotonicity — a pak napiš workflow, které
**generuje/benchmarkuje každou heuristiku paralelně přes N her** a **porotou**
vybere vítěze. (Jednodušší varianta: fan-out refaktor „konfigurovatelná
velikost desky", kde každý agent upraví jeden soubor, aby se z `board.Size`
stal parametr.)

**Co učí.** Deterministickou [orchestraci více agentů](06-workflows.cs.md):
fan-out přes `parallel()`, benchmark fázi, fázi poroty/redukce — vzorec
generuj-a-vyber.

**Soubory.** `internal/ai/ai.go` a `internal/ai/ai_test.go` (nové — drž
řešitel čistý, jako `board`), `cmd/2048-bench/main.go` (nový — bezhlavý
harness, na který se workflow napojí), `.claude/workflows/solver-benchmark.js`
(nový workflow pro výběr řešitele).

**Kontext na místě.** `.claude/workflows/test-coverage-audit.js` je kompletní,
funkční příklad přesně toho tvaru, který potřebuješ: blok `meta` s fázemi,
fan-out přes `parallel(...)`, strukturovaná output schémata pro každého agenta
a finální redukční fáze. Použij ho jako výchozí model — zkopíruj jeho strukturu
a vyměň read-only `Explore` agenty za build/benchmark kroky.

**Upozornění na výchozí bug `IsGameOver`.** Self-play volá `Board.IsGameOver()`
každý tah a výchozí verze (viz [Cvičení 5](#5-subagent)) hlásí „game over",
jakmile je deska plná, i když je pořád možné sloučení. Self-play hry proto
skončí trochu dřív a průměrná skóre vyjdou nižší, než by měla. Pro *porovnání*
heuristik je to v pořádku — každá heuristika je handicapovaná stejně — ale
absolutní čísla nečti jako skutečná 2048 skóre. Oprava je práce pro Cvičení 5;
tohle cvičení funguje tak jako tak.

**Startovací prompt.**
> *„Použij `.claude/workflows/test-coverage-audit.js` jako šablonu (zkopíruj
> jeho tvar `meta`/`phase`/`agent`/strukturované schéma). Napiš nový workflow
> `.claude/workflows/solver-benchmark.js`, který postaví několik 2048 heuristik
> — corner-stacking, greedy-merge, monotonicity, naimplementovaných jako čistý
> balíček `internal/ai` plus harness `cmd/2048-bench` — spustí je paralelně
> přes `parallel()` přes N self-play her a pak přidá fázi poroty/redukce, která
> je seřadí a vybere vítěze."*

**Walkthrough (kompletní řešení).**

1. **Čistý řešitel** (`internal/ai/ai.go`) — bez I/O, takže je každá heuristika
   unit-testovatelná. `ChooseMove` zkusí každý směr na *kopii* desky a nechá
   nejlepší legální; `PlayGame` odehraje celou hru bezhlavě:

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

   `trial := *b` zkopíruje celou strukturu `Board`, takže se look-ahead nikdy
   nedotkne skutečné desky — to je to, co drží `ChooseMove` čistý.

2. **Otestuj čistou logiku řešitele** (`internal/ai/ai_test.go`) — stejné
   konvence jako `internal/board/board_test.go`, jen standardní `testing`:

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

3. **Bezhlavá benchmark binárka** (`cmd/2048-bench/main.go`), na kterou se
   workflow napojí:

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

   Go nepotřebuje žádnou úpravu build souboru: `cmd/2048-bench/` sebere
   `go build ./...` automaticky a `%s` na `h` použije metodu `String()`
   definovanou výše.

4. **Workflow pro výběr** (`.claude/workflows/solver-benchmark.js`) —
   zkopírovaný z tvaru `test-coverage-audit.js`: blok `meta`, fan-out přes
   `parallel()` se schématy pro každého agenta a fáze poroty/redukce:

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

   Spusť ho přes `/solver-benchmark`. Rychlá kontrola z shellu —
   `go run ./cmd/2048-bench greedy-merge 20` — už vypíše skutečný řádek pro
   porotu k seřazení:

```text
heuristic=greedy-merge games=20 avg=2004.4 best=5356
```

**Hotovo.** Auto-řešitel odehraje celou hru bezhlavě; workflow proběhne a vydá
seřazené srovnání. `gofmt -l .` nic nevypíše a
`go build ./... && go vet ./... && go test ./...` je zelené, včetně nových
testů `internal/ai`.

---

Vyber libovolný řádek, vlož jeho startovací prompt a postav to. Každé cvičení tě
nechá s funkční featurou *a* mechanismem, který sis sám vyzkoušel.
