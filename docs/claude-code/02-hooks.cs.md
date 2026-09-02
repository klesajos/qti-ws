> 🌍 Číst v jazyce: [English](02-hooks.md) | **Česky**

# Ukázka 2: Projektové hooks

## Co je hook?

**Hook** je shellový příkaz, který Claude Code spustí **automaticky**, když
nastane určitá událost — třeba „po tom, co Claude upraví soubor" nebo „před
tím, než Claude spustí příkaz v terminálu".

Klíčový rozdíl oproti skillu: skill je *rada*, kterou se model může řídit;
hook je *vynucení*, které běží mimo model, **úplně pokaždé**. Hooks používej
na věci, které se nesmí nikdy přeskočit: formátování, lint, blokování
nebezpečných příkazů.

## Co dělá tahle ukázka

Pokaždé, když Claude upraví nebo vytvoří soubor, náš první hook zkontroluje,
jestli jde o Go soubor (`.go`) — a pokud ano, spustí na něm `gofmt`.
Výsledek: Claudův kód je vždycky zformátovaný podle toho, co chce `gofmt`,
i kdyby ho model napsal rozházeně. (Pokud je nainstalovaný `goimports`,
hook použije radši ten — je to `gofmt` plus oprava importů.)

Jsou v tom tři soubory:

```
.claude/settings.json        ← registruje, KDY se má který hook spustit
.claude/hooks/format-go.sh   ← skript, který říká, CO se má udělat po úpravě
.claude/hooks/run-tests.sh   ← druhý hook, spustí se, když Claude dokončí odpověď
```

Tenhle průvodce projde formátovací hook celý, pak **druhý hook** na jinou
událost — `Stop` hook, který spustí testy, když Claude dokončí odpověď — a
vysvětlí pět **typů** hooků.

## Část 1: registrace (`.claude/settings.json`) řádek po řádku

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

- `"hooks"` — sekce nastavení, kde žijí všechny hooks.
- `"PostToolUse"` — **událost**: spustit *po* úspěšném použití nástroje.
  Další užitečné události: `PreToolUse` (před spuštěním nástroje — umí ho
  zablokovat), `SessionStart`, `UserPromptSubmit`, `Stop` (když Claude
  končí odpověď).
- `"matcher": "Edit|Write"` — filtr na **název nástroje**. Svislítko `|`
  znamená „nebo", jako v regulárních výrazech: spustit jen když byl nástroj
  `Edit` nebo `Write` (dva nástroje, kterými Claude mění soubory). Bez
  matcheru by hook běžel i po každém `Read`, `Bash` atd.
- `"type": "command"` — tenhle hook spouští shellový příkaz (existují
  i jiné typy — viz [Typy hooků](#typy-hooků-command-http-mcp_tool-prompt-agent)
  níže).
- `"command": "${CLAUDE_PROJECT_DIR}/..."` — skript, který se má spustit.
  `${CLAUDE_PROJECT_DIR}` je proměnná, kterou Claude Code nahradí
  absolutní cestou ke kořenu repa — hook tak funguje bez ohledu na to,
  ve kterém podadresáři Claude zrovna je.

## Část 2: skript (`.claude/hooks/format-go.sh`) řádek po řádku

```bash
#!/usr/bin/env bash
set -euo pipefail

input=$(cat)
```

- `#!/usr/bin/env bash` — tzv. „shebang": říká systému, že soubor má
  spustit bashem. (Na Windows je to Git Bash — viz [instalační
  návod](00-install-claude-code.md).)
- `set -euo pipefail` — bezpečnostní pojistky: zastavit při jakékoli chybě
  (`-e`), nenastavené proměnné brát jako chybu (`-u`), selhat, když selže
  kterýkoli krok pipeline (`pipefail`). Standard pro každý bash skript.
- `input=$(cat)` — **takhle hook dostává data.** Claude Code posílá detaily
  volání nástroje jako JSON na **stdin** (standardní vstup). `cat` ho celý
  přečte a my ho uložíme do proměnné `input`.

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

- `command -v jq` — zkontroluje, jestli je nainstalovaný JSON nástroj `jq`
  (`>/dev/null 2>&1` jen skryje výstup kontroly).
- `jq -r '.tool_input.file_path // empty'` — vytáhne z JSONu cestu
  k upravenému souboru. Vstup vypadá takto:
  `{"tool_name": "Edit", "tool_input": {"file_path": "/cesta/k/board.go", ...}}`,
  takže `.tool_input.file_path` doklikne k cestě. `// empty` znamená
  „když chybí, nevypisuj nic (místo slova null)".
- Větev `elif` dělá úplně totéž v Pythonu — záloha pro stroje bez `jq`.
- Poslední větev je obyčejný `sed` vzor pro stroje bez obou (čerstvá
  instalace Windows s Git Bash má bash a sed, ale často žádné `jq` ani
  `python3`). Zvládne běžný případ; na cokoli složitějšího si nainstaluj
  `jq`.

```bash
if command -v goimports >/dev/null 2>&1; then
    formatter="goimports -w"
elif command -v gofmt >/dev/null 2>&1; then
    formatter="gofmt -w"
else
    exit 0
fi
```

- Vybere formátovač. `gofmt` je součástí každé instalace Go, takže na
  rozdíl od `clang-format` ve světě C++ tam prakticky vždycky je;
  `goimports` (z `golang.org/x/tools`) je volitelné vylepšení, které
  navíc přidává/odebírá řádky importů.
- `exit 0` — když formátovač není vůbec, skonči **úspěšně** a nedělej nic.
  Nenulový návratový kód by se Claudovi zobrazil jako chyba;
  „nenainstalovaný formátovač" nemá rozbít session.

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

- `case ... in *.go)` — porovnání přípony souboru. Dál projdou jen Go
  zdrojáky; `.md` nebo `.json` propadne a skript prostě skončí.
- `[ -f "$file_path" ]` — „existuje ten soubor?" (mezitím mohl být smazán).
- `$formatter "$file_path"` — samotná práce: `-w` (schované už v proměnné)
  znamená „write", tedy přepiš soubor zformátovaným obsahem.
- `echo ...` — cokoli hook vypíše, se ukáže v přepisu konverzace
  v Claude Code, takže vidíš, že hook odvedl svou práci.
  (`${file_path##*/}` odřízne adresářovou část, zůstane jen název souboru.)

## Vytvoř si vlastní hook, krok za krokem

1. **Napiš skript** do `.claude/hooks/`, např. `my-hook.sh`. Vyjdi
   z kostry výše: přečti stdin, vytáhni co potřebuješ, proveď akci, `exit 0`.
2. **Udělej ho spustitelným** — krok, na který každý zapomene:
   ```bash
   chmod +x .claude/hooks/my-hook.sh
   ```
3. **Zaregistruj ho** v `.claude/settings.json` pod správnou událost +
   matcher (viz Část 1).
4. **Nejdřív ho otestuj ručně** — neladit uvnitř Clauda. Stdin JSON si
   nasimuluj sám:
   ```bash
   echo '{"tool_input":{"file_path":"internal/board/board.go"}}' | .claude/hooks/my-hook.sh
   ```
5. **Spusť novou session Claude Code** a vyvolej událost doopravdy.
6. **Commitni oba soubory.**

## Vyzkoušej demo

Požádej Clauda, ať přidá metodu do `internal/board/board.go` a formátováním
se nezabývá. Po úpravě spusť `git diff` — kód už je zformátovaný a
v přepisu se objeví řádek `format-go hook: formatted board.go`. Spusť
`gofmt -l .` — nevypíše nic, protože nic nezůstalo nezformátované.

## Druhý hook: ohlas testy (událost `Stop`)

Formátovací hook reaguje na **nástroj** (`Edit`/`Write`). Hooks ale umí
reagovat i na **životní cyklus session**. Tohle repo má druhý hook na
události `Stop` — která nastane jednou, když Claude dokončí svou odpověď —
aby spustil testovací sadu a oznámil, jestli je strom stále zelený.

Jeho registrace sedí hned vedle prvního hooku v `.claude/settings.json`.
Všimni si, že tu **není `matcher`**: `Stop` není o nástroji, takže není co
filtrovat. Všimni si taky `timeout` — plný `go test ./...` se nejdřív musí
zkompilovat, a radši chceme, aby se hook zrušil, než aby zablokoval
session:

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

Skript (`.claude/hooks/run-tests.sh`) je záměrně **informativní**: vyprázdní
JSON události Stop, spustí `go test ./...`, vypíše jeden řádek a vždycky
`exit 0` — takže tě nikdy nepřeruší. Jeho výkonné řádky, doslova ze souboru:

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

- **`cat >/dev/null`** vyprázdní stdin. Žádné pole události Stop
  nepotřebujeme, ale hook musí svůj vstup přečíst, aby Claude Code
  nezůstal psát do zavřené roury.
- **Stráž (guard)** `if [ ! -f "$project_dir/go.mod" ] || ! command -v go …; then exit 0; fi`
  skončí, když neexistuje `go.mod` **nebo** není `go` v PATH — takže hook
  zůstane potichu ve špatném adresáři nebo na stroji bez Go, místo aby
  chyboval.
- **`set +e` … `set -e`** dočasně vypne „zastavit při chybě" kolem běhu
  testů: červená sada způsobí, že `go test` skončí nenulově, a my to
  chceme *ohlásit*, ne nechat hook shodit. `status=$?` uchová návratový
  kód.
- **Souhrn** spočítá řádky s výsledky jednotlivých balíčků, které `go test`
  vypisuje — `ok <balíček>` pro zelené balíčky, `FAIL <balíček>` pro
  červené — a spojí názvy padlých balíčků čárkami.
- **Verdikt** vypíše `✓` s počtem balíčků, když byl status 0, jinak `✗`
  s padlými balíčky a tipem, ať se spustí `go test`.

**Chceš, aby spíš *blokoval* než informoval?** `Stop` hook, který skončí
kódem **2**, nebo vypíše `{"decision": "block", "reason": "..."}` na
stdout, říká Claudovi, že **ještě není** hotovo — tak vynutíš „testy musí
projít, než skončíš". My ten náš necháváme informativní, aby se živá
workshopová session nezasekla v opravovací smyčce; přepni ho na blokující,
když chceš tvrdou bránu ([Cvičení 2](exercises.cs.md#2-hook) dělá přesně
tohle).

## Vyzkoušej druhý hook

Zeptej se Clauda na cokoli, co ukončí odpověď. Když odpověď skončí, `Stop`
hook spustí sadu a v přepisu se objeví informativní řádek:

```
run-tests hook: ✓ go test ./... — 1 package(s) ok
```

Teď záměrně rozbij nějaký test (změň očekávanou hodnotu), zeptej se Clauda
na něco triviálního a sleduj, jak se řádek změní na `✗ … failed in
github.com/klesajos/qti-ws/internal/board`. Změnu vrať zpět.

## Typy hooků: command, http, mcp_tool, prompt, agent

Oba hooky výše používají `"type": "command"` — shellový skript. To je
nejběžnější typ, ale ne jediný. Handler může být jeden z pěti typů a
směňuje determinismus za úsudek:

| Typ | Co běží | Na co se hodí |
|-----|---------|---------------|
| `command` | Shellový skript (naše dva hooky) | Deterministické, rychlé kontroly — formát, lint, spuštění testů, zablokování zakázaného příkazu |
| `http` | HTTP POST s JSONem události na URL | Centrální policy služby, audit logging — endpoint odpovídá stejným JSONem jako command hook |
| `mcp_tool` | Nástroj na už připojeném MCP serveru ([Ukázka 3](03-mcp.cs.md)) | Znovupoužít schopnost, kterou už máš (zapsat do trackeru, ověřit policy službu) |
| `prompt` | Jednokolové volání modelu Claude | Rychlý úsudek — „je commit message výstižná?" — vrací rozhodnutí ano/ne |
| `agent` | Vícekolový subagent s přístupem k nástrojům (`Read`, `Grep`, …) | Hloubková kontrola — „přečti diff a ověř, že nepřibyla žádná tajemství" — než se pokračuje (experimentální) |

`prompt` hook se registruje s textem k vyhodnocení místo příkazu:

```json
{ "type": "prompt",
  "prompt": "Does the staged diff add a test for every new exported function? Answer yes or no." }
```

Základní pravidlo: sáhni nejdřív po `command` (je zdarma a okamžitý), po
`prompt`, když kontrola potřebuje porozumění jazyku, a po `agent` jen
tehdy, když se kontrola sama musí prohrabat kódem. Tohle repo dodává
`command` hooky; o zbylých čtyřech stačí vědět, že existují.

## Volitelné parametry

Hook handler umí víc než `type` a `command`. Nejužitečnější volitelná pole:

| Pole | Co dělá |
|------|---------|
| `timeout` | Sekundy do zrušení hooku (výchozí 600 pro `command`). Náš Stop hook má 120, aby pomalý běh testů nemohl zaseknout session |
| `statusMessage` | Vlastní text u spinneru během běhu, např. `"Formatting Go..."` |
| `if` | Dodatečný filtr syntaxí permission pravidel, např. `"if": "Edit(*.go)"` — přesnější než `matcher`, který vidí jen název nástroje. Vzor cesty jako `internal/**` sedí jen pod pracovním adresářem; na libovolnou hloubku použij `**/internal/**` |
| `once: true` | Spustit jen jednou, pak se odregistrovat — **respektuje se jen u hooků deklarovaných ve frontmatteru skillu**, v `settings.json` se ignoruje |

A události: tahle ukázka používá `PostToolUse` a `Stop`, ale hooks se dají
navěsit na celý životní cyklus session. Které stojí za to znát nejdřív:

| Událost | Spustí se |
|---------|-----------|
| `PreToolUse` | Před spuštěním nástroje — **umí ho zablokovat** (např. zakázat `git push --force`) |
| `PostToolUse` | Po úspěšném nástroji (náš formátovací hook) |
| `SessionStart` | Při startu session — kontroly prostředí, načtení kontextu |
| `UserPromptSubmit` | Při odeslání promptu — umí přidat kontext navíc |
| `Stop` | Když Claude končí odpověď (náš testovací hook) — umí zablokovat ukončení |
| `SubagentStart` / `SubagentStop` | Kolem běhu subagenta ([Ukázka 5](05-agents.cs.md)) |
| `SessionEnd` | Když session skončí — úklid, reporting |

Existuje jich víc (`PostToolUseFailure`, `PreCompact`, `PreModelSwitch`,
`DirectoryAdded`, …). Úplný seznam událostí a polí:
[oficiální dokumentace hooks](https://code.claude.com/docs/en/hooks).

## Kde to funguje: CLI, Desktop aplikace, Cowork

| Platforma | Funguje? | Nastavení |
|-----------|----------|-----------|
| **Claude Code CLI** (terminál) | ✅ Ano | Nic navíc — hooks z `.claude/settings.json` se načtou při startu session |
| **Claude Desktop — záložka Code** | ✅ Ano | Stejný engine, stejné konfigurační soubory jako CLI. Potvrď jednorázový dialog důvěry projektu; hooks pak běží stejně |
| **Cowork** (v Desktop aplikaci) | ❌ Ne | Sandboxované VM Coworku projektové hooks z `.claude/settings.json` nespouští. Přímá náhrada není — nejblíž jsou hooks zabalené v nainstalovaném pluginu |

Poznámka pro workshop: tohle je nejvýraznější platformní rozdíl ze všech
šesti ukázek. Hooks jsou *lokální automatizace* — pokud na nich tvůj
workflow stojí (formátování, testovací brány), pracuj v CLI nebo
v záložce Code v Desktopu, ne v Coworku.

## Když něco nefunguje

- **Hook se nikdy nespustí** → po úpravě `settings.json` je potřeba nová
  session; zkontroluj taky spustitelnost skriptu (`ls -l .claude/hooks/`).
- **„Permission denied"** → přeskočil jsi `chmod +x`.
- **Chyby hooku ruší práci** → každá větev „není co dělat" musí končit
  `exit 0`, ne chybou.
- **Na Windows hook padá s `bash: not found`** → nainstaluj Git for
  Windows, aby Claude Code mohl jako svůj shell použít Git Bash (viz
  [instalační návod](00-install-claude-code.md)).
</content>
