---
name: board-tests
description: Use when writing or extending Go unit tests for the 2048 board logic (SlideLineLeft, Board.Move, win/game-over detection). Covers project test conventions and how to run the suite.
---

# Writing board tests

All game rules live in `internal/board/board.go` and are free of I/O — that is
what makes them unit-testable. Tests go in `internal/board/board_test.go`
(package `board`, standard library `testing`).

## Conventions

- One `func TestX(t *testing.T)` per behaviour, named
  `Test<Function>_<ScenarioInPlainEnglish>`, e.g.
  `TestSlideLineLeft_TwoSeparatePairsMergeIntoTwoTiles`.
- Arrange-Act-Assert: build the state, perform one action, then compare with
  `if got != want { t.Errorf(...) }`. No assertion libraries.
- Prefer testing the free function `SlideLineLeft(&line)` for merge rules — it
  is the heart of the game and takes a `*Line` (`[4]int`).
- For whole-board behaviour use the deterministic constructor
  `FromGrid(grid Grid, score int)` — it spawns no random tiles. Rows you leave
  out of a `Grid{...}` literal are all zero.
- `Grid` and `Line` are fixed-size arrays: compare them with `==` / `!=`,
  never with `reflect.DeepEqual`.
- Never test through `game`, `input`, or `renderer`; rules belong to `board`.
- Reuse the local `slid()` helper in `board_test.go` when the score gained is
  irrelevant; call `SlideLineLeft()` directly when asserting the score.
- Use `t.Fatal` only when continuing makes no sense (a precondition failed);
  otherwise `t.Errorf`, so every mismatch in a test is reported.

## Key API (from `internal/board/board.go`)

```go
func SlideLineLeft(line *Line) int           // returns score gained
func FromGrid(grid Grid, score int) *Board   // deterministic, for tests
func (b *Board) Move(dir Direction) bool     // true if the board changed
func (b *Board) SpawnRandom() bool           // false if the board is full
func (b *Board) IsGameOver() bool
func (b *Board) HasWon() bool                // any tile >= WinValue (2048)
func (b *Board) Grid() Grid                  // copy of the grid
func (b *Board) Score() int
func (b *Board) At(row, col int) int
```

Directions are `Up`, `Down`, `Left`, `Right` (type `Direction`).

## Running the tests

```bash
go test ./...                                     # whole repo
go test -v ./internal/board/                      # one package, verbose
go test -run TestSlideLineLeft ./internal/board/  # a subset by name
```

After adding a test, always run the suite before declaring it done. If a
*correct* test fails, report that as a finding — never weaken the assertion to
force a green run.
