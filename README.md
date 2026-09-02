# 2048 (terminálová verze)

Malá hra 2048 v Go. Repo slouží jako základ pro hands-on část workshopu
o AI-asistovaném programování s Claude Code — budeme do něj přidávat funkce,
opravovat chyby, dopisovat testy a refaktorovat.

## Build & spuštění

Potřebuješ Go 1.25+ (`go version`).

```bash
go build ./...
go run ./cmd/2048
```

Při prvním buildu si Go stáhne jedinou závislost, `golang.org/x/term`
(přepnutí terminálu do raw módu) — je tedy potřeba internet.

## Testy

```bash
go test ./...
go test -v ./internal/board/   # podrobný výpis jednoho balíčku
```

## Ovládání

Šipky nebo `WASD` pro pohyb, `q` pro ukončení. Tiles se posouvají daným směrem
a stejné hodnoty se sloučí. Cíl je dosáhnout dlaždice **2048**.

## Struktura

```
cmd/2048/main.go          vstupní bod
internal/
  board/board.go          herní logika — posun, slučování, detekce výhry/konce (bez I/O)
  board/board_test.go     unit testy herní logiky (standardní `testing`)
  game/game.go            hlavní smyčka, propojuje board + input + renderer
  input/input.go          čtení kláves z terminálu (raw mód)
  renderer/renderer.go    vykreslení mřížky a skóre
go.mod / go.sum
```

Veškerá pravidla hry jsou v balíčku `board` a nezávisí na I/O — proto se dají
snadno testovat. Začni u `SlideLineLeft()` v `internal/board/board.go`, to je
srdce hry.

## Ukázky Claude Code

Repo zároveň ukazuje šest projektových způsobů, jak rozšířit Claude Code —
skills, hooks, MCP servery, pluginy, subagenty a workflows. Každá ukázka je
samostatná, minimální a navázaná na tenhle kód. Návody jsou dvojjazyčné:
[anglicky](docs/claude-code/README.md) | [česky](docs/claude-code/README.cs.md).

Než začneš, projdi krátké přípravné návody (instalace Claude Code, Git, Go) —
odkazy jsou v sekci *Before you start* / *Než začneš* v návodech výše.
