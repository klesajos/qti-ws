> 🌍 Číst v jazyce: [English](00-go-basics.md) | **Česky**

# Před workshopem: spuštění Go appky

Workshopová appka je terminálové 2048 napsané v Go. Tahle stránka je
všechno, co potřebuješ k jejímu sestavení, spuštění a otestování — i kdyby
sis Go nikdy nesáhl.

## Instalace Go

Potřebuješ **Go 1.25 nebo novější** (verzi ze `go.mod`).

```bash
brew install go                      # macOS
winget install GoLang.Go             # Windows (PowerShell) — nebo instalátor z go.dev/dl
# Linux: stáhni tarball z https://go.dev/dl/, pak
sudo rm -rf /usr/local/go && sudo tar -C /usr/local -xzf go1.25.*.linux-amd64.tar.gz
export PATH=$PATH:/usr/local/go/bin  # přidej i do ~/.profile
```

Zkontroluj: `go version` → `go version go1.25.x …`.

## Tři příkazy, které použiješ celý den

```bash
go run ./cmd/2048        # zkompiluje + spustí hru (šipky/WASD pohyb, q ukončí)
go test ./...            # spustí všechny testy v repu
go build ./...           # zkompiluje všechno, nahlásí chyby, nic nezapíše
```

`./...` znamená „tenhle balíček a všechno pod ním" — uvidíš to všude.

## Projekt v kostce

```
go.mod                  modul: jeho název (github.com/klesajos/qti-ws), verze Go, závislosti
go.sum                  kontrolní součty závislostí (commitni ho, needituj)
cmd/2048/main.go        vstupní bod spustitelného souboru (package main)
internal/board/         pravidla hry — čistá logika, žádné I/O, plně otestováno
internal/game/          hlavní smyčka
internal/input/         klávesnice (terminál v raw-mode)
internal/renderer/      vykreslení mřížky
```

- **Balíček = adresář.** Všechny soubory v `internal/board/` začínají
  `package board`. `internal/` je konvence Go: obsah smí importovat jen
  tenhle modul.
- **Importy používají cestu modulu:** `import "github.com/klesajos/qti-ws/internal/board"`.
- **Velké písmeno = exportované.** `board.SlideLineLeft` je veřejné;
  `board.newRNG` je soukromé pro balíček.
- **Testy žijí vedle kódu** v souborech `*_test.go` a jsou obyčejné funkce:
  `func TestSomething(t *testing.T)`.

## Sestavení a spuštění

```bash
go build -o 2048 ./cmd/2048    # vyrobí binárku 2048 v kořeni repa
./2048                         # spusť ji (Windows: .\2048.exe)
go run ./cmd/2048              # totéž bez uchování binárky
```

První sestavení stáhne jedinou závislost (`golang.org/x/term`) do cache
modulů Go; další sestavení jsou offline a rychlá.

## Testování

```bash
go test ./...                              # všechno; "ok" u balíčku = zelená
go test -v ./internal/board/               # verbose: každý test se jménem + PASS/FAIL
go test -run TestSlideLineLeft ./internal/board/   # jen testy, jejichž jméno sedí
go test -cover ./internal/board/           # s procentem pokrytí
```

Čtení selhání:

```
--- FAIL: TestSlideLineLeft_SinglePairMerges (0.00s)
    board_test.go:37: line = [8 0 0 0], want [4 0 0 0]
FAIL
FAIL	github.com/klesajos/qti-ws/internal/board	0.4s
```

Soubor a řádek selhavší kontroly, pak zpráva, kterou test napsal — `got` vs
`want`. `?  … [no test files]` u balíčku je normální, ne chyba.

## Udržování kódu v čistotě

```bash
gofmt -l .        # vypíše soubory, které nejsou zformátované (prázdný výstup = vše OK)
gofmt -w .        # zformátuje je na místě (hook repa tohle dělá za Claudovy úpravy)
go vet ./...      # statické kontroly na běžné chyby
go mod tidy       # doplní chybějící / odebere nepoužité závislosti v go.mod
```

## Chyba, na kterou narazí každý

```
./game.go:31:2: declared and not used: changed
```

Go odmítne zkompilovat proměnnou, kterou nikdy nepřečteš. Buď ji použij,
nebo smaž — žádná úroveň varování neexistuje. (Drž si tohle v hlavě během
cvičení: jeden ze zasazených bugů je *návratová hodnota*, která se tiše
zahodí — tenhle kompilátor **nezachytí**.)

## Tahák

| Chci… | Příkaz |
|---|---|
| Spustit hru | `go run ./cmd/2048` |
| Spustit všechny testy | `go test ./...` |
| Spustit jeden test verbose | `go test -v -run TestName ./internal/board/` |
| Zkontrolovat, že se to kompiluje | `go build ./...` |
| Zformátovat všechno | `gofmt -w .` |
| Statické kontroly | `go vet ./...` |
| Přidat závislost | `go get example.com/pkg@latest` pak `go mod tidy` |

## Kontrola

```bash
cd qti-ws && go version && go build ./... && go test ./...
```

Očekáváno: verze Go, pak `ok  github.com/klesajos/qti-ws/internal/board`.
