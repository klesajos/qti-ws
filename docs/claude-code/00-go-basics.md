> 🌍 Read this in: **English** | [Česky](00-go-basics.cs.md)

# Before the workshop: running a Go app

The workshop app is a terminal 2048 written in Go. This page is everything
you need to build, run and test it — even if you've never touched Go.

## Install Go

Need **Go 1.25 or newer** (the version in `go.mod`).

```bash
brew install go                      # macOS
winget install GoLang.Go             # Windows (PowerShell) — or the installer from go.dev/dl
# Linux: download the tarball from https://go.dev/dl/, then
sudo rm -rf /usr/local/go && sudo tar -C /usr/local -xzf go1.25.*.linux-amd64.tar.gz
export PATH=$PATH:/usr/local/go/bin  # add to ~/.profile too
```

Check: `go version` → `go version go1.25.x …`.

## The three commands you'll use all day

```bash
go run ./cmd/2048        # compile + run the game (arrows/WASD move, q quits)
go test ./...            # run every test in the repo
go build ./...           # compile everything, report errors, write nothing
```

`./...` means "this package and everything below it" — you'll see it
everywhere.

## The project at a glance

```
go.mod                  the module: its name (github.com/klesajos/qti-ws), Go version, dependencies
go.sum                  checksums of the dependencies (commit it, don't edit it)
cmd/2048/main.go        the executable's entry point (package main)
internal/board/         the game rules — pure logic, no I/O, fully tested
internal/game/          the main loop
internal/input/         keyboard (raw-mode terminal)
internal/renderer/      drawing the grid
```

- **A package = a directory.** Files in `internal/board/` all start with
  `package board`. `internal/` is a Go convention: only this module can
  import what's inside.
- **Imports use the module path:** `import "github.com/klesajos/qti-ws/internal/board"`.
- **Capitalised = exported.** `board.SlideLineLeft` is public;
  `board.newRNG` is private to the package.
- **Tests live next to the code** in `*_test.go` files and are plain
  functions: `func TestSomething(t *testing.T)`.

## Building and running

```bash
go build -o 2048 ./cmd/2048    # produce a binary named 2048 in the repo root
./2048                         # run it (Windows: .\2048.exe)
go run ./cmd/2048              # same thing without keeping the binary
```

The first build downloads the one dependency (`golang.org/x/term`) into
Go's module cache; later builds are offline and fast.

## Testing

```bash
go test ./...                              # everything; "ok" per package = green
go test -v ./internal/board/               # verbose: every test name + PASS/FAIL
go test -run TestSlideLineLeft ./internal/board/   # only tests whose name matches
go test -cover ./internal/board/           # with a coverage percentage
```

Reading a failure:

```
--- FAIL: TestSlideLineLeft_SinglePairMerges (0.00s)
    board_test.go:37: line = [8 0 0 0], want [4 0 0 0]
FAIL
FAIL	github.com/klesajos/qti-ws/internal/board	0.4s
```

File and line of the failing check, then the message the test wrote —
`got` vs `want`. `?  … [no test files]` for a package is normal, not an error.

## Keeping the code tidy

```bash
gofmt -l .        # list files that aren't formatted (empty output = all good)
gofmt -w .        # format them in place (the repo's hook does this for Claude's edits)
go vet ./...      # static checks for common mistakes
go mod tidy       # add missing / remove unused dependencies in go.mod
```

## The one error everyone hits

```
./game.go:31:2: declared and not used: changed
```

Go refuses to compile a variable you never read. Either use it or delete
it — there's no warning level. (Keep this in mind during the exercises:
one of the planted bugs is a *return value* that's silently discarded — that
one the compiler does **not** catch.)

## Cheat table

| I want to… | Command |
|---|---|
| Run the game | `go run ./cmd/2048` |
| Run all tests | `go test ./...` |
| Run one test verbosely | `go test -v -run TestName ./internal/board/` |
| Check it compiles | `go build ./...` |
| Format everything | `gofmt -w .` |
| Static checks | `go vet ./...` |
| Add a dependency | `go get example.com/pkg@latest` then `go mod tidy` |

## Check

```bash
cd qti-ws && go version && go build ./... && go test ./...
```

Expected: the Go version, then `ok  github.com/klesajos/qti-ws/internal/board`.
