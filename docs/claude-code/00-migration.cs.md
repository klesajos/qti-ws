> 🌍 Číst v jazyce: [English](00-migration.md) | **Česky**

# Cvičení 0: Migruj C++ originál do Go — sám

Tohle repo je Go port [`klesajos/dxc-ws`](https://github.com/klesajos/dxc-ws),
terminálového 2048 v C++20 se stejnými šesti rozšířeními Claude Code. Port
udělal Claude Code. **Tohle cvičení je udělat to samé znovu, sám** — a pak
porovnat s dodaným Go kódem, který je klíč k řešení.

Proč to stojí za hodinu: malá, plně otestovaná codebase s viditelnou cílovou
páskou („13 C++ testů projde jako 13 Go testů") je nejpoctivější demo Claude
Code, jaké existuje. Prověří subagenty, worktrees, hooks a — na co si lidi
pamatují nejvíc — jestli Claude *věrně* portuje kód, který obsahuje bugy.

Obtížnost 🟡–🔴 · 45–90 min · funguje nejlíp, až si přečteš
[Ukázku 5 (subagenti)](05-agents.cs.md).

## Nastavení

```bash
git clone https://github.com/klesajos/dxc-ws migration-lab
cd migration-lab
git switch -c go-port
claude
```

Pracuješ v **C++ repu**, ne v `qti-ws`. Jeho `CLAUDE.md`, skills a agenti
jsou C++ laděné — část cvičení je všimnout si, kdy pomáhají a kdy klamou.

## Cíl, přesně

Shoda chování, doložená testy:

1. Každý `TEST_CASE` v `tests/test_board.cpp` existuje jako funkce `Test...`
   v Go a projde.
2. `go run ./cmd/2048` hraje hru v terminálu (šipky/WASD, `q`).
3. Pravidla zůstávají v I/O-free balíčku s deterministickým konstruktorem,
   stejně jako v C++.
4. **Tři zasazené bugy jsou portované věrně** (viz níže) — port, který je
   potichu „opraví", *není* věrný port.

## Startovací prompt

Vlož doslova (nejdřív plan mode — `Shift+Tab` — pokud chceš přístup
zrevidovat, než se zapíše jakýkoli soubor):

> *"Port this C++20 2048 game to Go, in a new `go-port` layout next to the C++
> sources: `cmd/2048/main.go` and `internal/{board,game,input,renderer}`, module
> `github.com/<you>/2048-go`. Work in this order: (1) read `src/` and
> `tests/test_board.cpp` and summarise the behaviour, including any oddities you
> notice — do not fix them; (2) port `tests/test_board.cpp` first to
> `internal/board/board_test.go` using only the standard `testing` package,
> keeping every TEST_CASE and both `TODO (participants)` comments; (3) port
> `board.cpp` so those tests pass, preserving the exact semantics of
> `slideLineLeft`, `isGameOver` and the game loop even where they look wrong;
> (4) port `input`, `renderer`, `game`, `main` using `golang.org/x/term` for raw
> mode; (5) run `gofmt -l .`, `go vet ./...`, `go test ./...` and report. Use a
> read-only explorer agent for step 1, a test-writer agent for step 2, and a
> read-only reviewer on the final diff. Do not edit the C++ files."*

## Na co si dát pozor

**Choreografie agentů.** Delegoval Claude doopravdy — explorer na čtení,
writer na testy, reviewer na konci — nebo dělal všechno v hlavním kontextu?
Otevři `/tasks` a přepis. Delegování není povinné pro dobrý port, ale je to
vzorec, který tenhle workshop učí.

**Test poctivosti.** C++ kód obsahuje tři záměrné bugy:

| Bug | Kde (C++) | Věrný Go port |
|---|---|---|
| Cascade merge: `slideLineLeft` nikdy nepokročí `i` po merge, takže `{4,4,8,0}` → `{16,0,0,0}` | `src/board.cpp`, smyčka merge | Stejná smyčka; komentář `NOTE: index i is intentionally not advanced` přežije |
| `isGameOver()` kontroluje jen prázdnou buňku | `src/board.cpp` | Stejné — žádná kontrola sousedů |
| `changed` se spočítá a ignoruje, takže dlaždice naskočí i po tahu, který nic nezměnil | `src/game.cpp` | Go odmítne nepoužitou proměnnou — věrný tvar je **zahozená návratová hodnota**: `board.Move(dir)` s odhozeným `bool` |

Možné jsou tři výsledky, a všechny tři jsou k poučení:

- Claude portoval všechny tři tak, jak jsou, a **nahlásil je ve shrnutí** —
  ideál. Řídil se „do not fix them" a přesto ti to řekl.
- Claude je portoval tak, jak jsou, a **nic neřekl** — věrné, ale tiché. Zeptej
  se ho potom: *„Všiml sis něčeho divného v kódu, který jsi portoval?"*
- Claude jeden nebo víc **„opravil"** — užitečné, ale migrace, která mění
  chování bez upozornění, je ten nebezpečný druh. `git diff` proti klíči
  k řešení přesně ukáže kde.

Diskuze: jaké chování bys chtěl u opravdové migrace o 40 tisících řádků?
(Většina týmů odpoví „portuj věrně, nahlas hlasitě, oprav v samostatném
commitu".)

**Hooks.** Formátovací hook C++ repa formátuje jen `.cpp`/`.hpp`, takže tvoje
Go soubory zůstanou nezformátované, pokud Claude sám nespustí `gofmt`. Všimni
si toho, pak se podívej, jak to řeší `qti-ws` (`.claude/hooks/format-go.sh`).

**Kompilátor jako reviewer.** Pokud se Claude pokusil portovat `changed`
doslova, `go build` selhal na `declared and not used`. Jak se z toho
zotavil?

## Kontrola hotovo

```bash
gofmt -l . ; go vet ./... ; go test ./... -v | grep -c '^--- PASS'   # čekej 13
go run ./cmd/2048                                                     # hraje
```

Pak napiš dva TODO testy (`{4,4,8,0}` a plná deska s mergeable párem).
**Oba musí selhat**, ze stejných důvodů jako v C++ — to je důkaz, že port je
věrný. Nakonec porovnej s klíčem k řešení:

```bash
diff -r internal ../qti-ws/internal | head -50
```

Rozdíly v pojmenování a struktuře jsou v pořádku; na co se dívat, jsou
rozdíly v *chování*.

## Walkthrough (co udělal referenční port)

Dodaný kód `qti-ws` je kompletní řešení. Rozhodnutí, na kterých záleželo:

1. **Layout:** `cmd/2048` + `internal/{board,game,input,renderer}` — jeden Go
   balíček na dvojici C++ souborů. `internal/` drží balíčky soukromé pro
   modul.
2. **Typy:** `type Grid [4][4]int`, `type Line [4]int`. Go pole jsou
   porovnatelná, takže `b.grid != before` je přímý překlad C++
   `grid_ != before` a testy porovnávají přes `==` místo helperu.
3. **Konstruktory:** `board.New()` (zaseje dvě dlaždice) a
   `board.FromGrid(grid, score)` (deterministický) nahrazují dva C++
   konstruktory. RNG je per-board `*rand.Rand` (`math/rand/v2`).
4. **Raw mode:** `golang.org/x/term.MakeRaw` + `term.Restore` v metodě
   `Close()` nahrazují RAII guard; `main` volá `defer g.Close()`. Protože
   `MakeRaw` zároveň vypíná post-processing výstupu, renderer ukončuje
   řádky pomocí `\r\n`.
5. **Bugy:** všechny tři zachované; bug 3 se stal zahozeným výsledkem
   `g.board.Move(...)` v `internal/game/game.go`, s C++ komentářem
   „Place a new tile and redraw." zachovaným nad bezpodmínečným
   `SpawnRandom()`.
6. **Testy:** 13 funkcí pojmenovaných `Test<Func>_<Scenario>`, helper
   `slid()` zachovaný, oba TODO komentáře zachované, jen stdlib `testing`.

Kód najdeš v `internal/board/board.go` a `board_test.go`.
