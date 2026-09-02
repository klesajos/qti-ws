# 2048

Terminal version of the game 2048 in Go. Used as the base repository for an
AI-assisted coding workshop (adding features, fixing bugs, writing tests,
refactoring).

## Stack

- Language: Go 1.25+ (uses `math/rand/v2`)
- Dependencies: `golang.org/x/term` (raw terminal mode) — nothing else
- Test framework: the standard library `testing` package

## Directory Structure

- `cmd/2048/main.go` - entry point
- `internal/board/` - game logic: slide, merge, win/game-over detection (no I/O)
  - `board.go`, `board_test.go`
- `internal/game/` - main loop, wires board + input + renderer together
- `internal/input/` - terminal keyboard input in raw mode
- `internal/renderer/` - draws the grid and score

## Conventions

- Standard Go style; code is formatted with `gofmt` (enforced by a hook)
- One package per concern under `internal/`; only `cmd/2048` is importable from outside
- Exported identifiers carry doc comments (`// Name ...`)
- Keep game rules in `board`, free of I/O, so they stay unit-testable
- Tests: one `func TestX(t *testing.T)` per behaviour, Arrange-Act-Assert,
  `if got != want { t.Errorf(...) }` — no assertion libraries
- Fixed-size arrays (`board.Grid`, `board.Line`) are compared with `==`
- The terminal is in raw mode while the game runs: output lines end with `\r\n`

## Commands

- build: `go build ./...`
- run: `go run ./cmd/2048`
- test: `go test ./...`
- test (verbose, one package): `go test -v ./internal/board/`
- vet: `go vet ./...`
- format check: `gofmt -l .` (prints nothing when everything is formatted)
- tidy deps: `go mod tidy`

## Build Dependencies

- Go 1.25+ (`go version`)
- Internet on first build (`go mod download` fetches `golang.org/x/term`)
- A POSIX-style terminal (Linux / macOS / WSL / Git Bash) to play; tests run anywhere
