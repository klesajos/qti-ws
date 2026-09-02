> 🌍 Číst v jazyce: [English](06-workflows.md) | **Česky**

# Ukázka 6: Workflows (orchestrace více agentů)

## Co je workflow?

**Workflow** je malý JavaScriptový program, který **orchestruje víc agentů
s deterministickým tokem řízení**. Místo aby ses Clauda ptal, ať „si vymyslí
kroky", kroky napíšeš sám — rozfanouješ práci přes `parallel()`, zřetězíš ji
přes `pipeline()`, seskupíš přes `phase()` — a runtime tu strukturu pokaždé
vykoná přesně tak. Uložené workflow se stane slash příkazem, který spustíš přes
`/<jmeno>`.

Tam, kde si jeden agent volí cestu sám, je cesta workflow pevně daná v kódu:
smyčky, podmínky i fan-out máš pod kontrolou ty.

> Dynamická workflow jsou dostupná na všech placených plánech. Na **Pro** je
> zapneš jednou v `/config` → *Dynamic workflows*.

## Subagent vs. workflow

Tohle jsou dva delegační mechanismy a vyplatí se mít v jejich rozdílu jasno,
než si přečteš kód:

| | Subagent (Ukázka 5) | Workflow (tahle ukázka) |
|---|---|---|
| **Jednotka** | jeden delegovaný kontext | víc agentů, orchestrovaných |
| **Tok řízení** | kroky si volí agent | volíš *ty*, v JavaScriptu (`parallel`/`pipeline`/`phase`) |
| **Determinismus** | řízeno modelem, liší se běh od běhu | struktura je pevná; stejný vstup → stejný tvar běhu |
| **Definováno v** | Markdown souboru + systémový prompt | `.js` skriptu s `export const meta` |
| **Vyvolání** | prostá řeč, `@agent-name`, `--agent` | slash příkaz `/<jmeno>` (opt-in) |

**Základní pravidlo:** sáhni po **subagentovi**, když chceš předat *jeden*
samostatný úkol *jednomu* izolovanému kontextu. Sáhni po **workflow**, když chceš
spustit *víc* agentů v *opakovatelné, kódem definované* pipeline — rozfanoutovat
přes N souborů, podmínit fázi 2 fází 1, zredukovat víc výsledků do jedné zprávy.

## Co tahle ukázka dělá

Tahle ukázka přidává jedno spustitelné workflow, `test-coverage-audit`, které
audituje, jak dobře `internal/board/board_test.go` pokrývá logiku ve zdrojových
souborech. Je **jen pro čtení, deterministické a opakovatelné** — vypíše
prioritizovanou zprávu o mezerách a nikdy nezapíše soubor, takže ho můžeš
spouštět, jak chceš často, se stabilním výsledkem.

Funguje ve třech fázích:

1. **Inventory** — rozfanouj jednoho read-only čtenáře na každý zdrojový soubor;
   každý vyjmenuje chování, na která může mířit unit test.
2. **Map tests** — jeden agent namapuje každou existující funkci `Test...` na
   tenhle inventář.
3. **Report gaps** — jeden reduktor obojí zkříží a vydá prioritizovanou zprávu
   o mezerách.

Protože skutečné mezery v tomhle repu jsou známé (TODO `{4, 4, 8, 0}` a bug
s kaskádovým mergem za ním, distribuce `SpawnRandom`, game-over plné desky
s mergeable párem, který naráží na bug v `IsGameOver()`, vícetahové sekvence,
zahozený návratový výsledek `Move()` v `game.go`), je zpráva předvídatelná —
což z ní dělá dobrý učební artefakt.

## Soubor řádek po řádku

Workflow žije v `.claude/workflows/test-coverage-audit.js`. Uložená workflow
v `.claude/workflows/` se auto-objeví a stanou se slash příkazy.

**1. Blok `meta` musí být úplně první příkaz — na tomhle to každý poprvé
zasekne.** Jakýkoli `import`, `const` nebo spustitelný řádek nad
`export const meta` způsobí, že se workflow **tiše neregistruje**:
`/test-coverage-audit` se vůbec neobjeví, bez chyby, která by vysvětlila proč.
Komentář `//` v hlavičce souboru před `meta` nevadí (tenhle soubor jedním
začíná); spustitelný kód ano. Blok pojmenuje workflow (tím vznikne
`/test-coverage-audit`) a popíše jeho fáze:

```js
export const meta = {
  name: 'test-coverage-audit',
  description: 'Audit internal/ test coverage against internal/board/board_test.go and emit a prioritized gap report.',
  phases: [
    { title: 'Inventory', detail: 'one read-only reader per source file enumerates behaviours' },
    { title: 'Map tests', detail: 'map each Test function onto the behaviour inventory' },
    { title: 'Report gaps', detail: 'cross-reference and emit a prioritized gap report' },
  ],
}
```

**2. Pět zdrojových souborů, které inventarizuje.** Pole `SRC_FILES` jmenuje
každý Go soubor mimo testy, u každého jednořádkovou poznámku, ať se čtenářský
agent rychle zorientuje. Tohle je doslova těch „pět souborů", přes která se
rozfanouje fáze 1:

```js
const SRC_FILES = [
  { file: 'internal/board/board.go', note: 'pure rules: SlideLineLeft, Move, SpawnRandom, IsGameOver, HasWon' },
  { file: 'internal/game/game.go', note: 'main loop: toDirection, Game.Run, spawn + redraw, win/over checks' },
  { file: 'internal/input/input.go', note: 'raw-mode terminal: ReadByte -> Command (WASD + arrow keys, q)' },
  { file: 'internal/renderer/renderer.go', note: 'draws the grid + score to the terminal' },
  { file: 'cmd/2048/main.go', note: 'entry point: game.New, Run, Close' },
]
```

**3. Tři schémata nutí každého agenta vrátit data, ne prózu.** Předání `schema`
do `agent()` přiměje runtime výstup agenta zvalidovat, takže skript dostane zpět
skutečný objekt — žádné parsování. Workflow definuje jedno na každou fázi:
`BEHAVIOR_SCHEMA` (co jeden soubor dělá, větev po větvi), `COVERAGE_SCHEMA`
(která id chování pokrývá každá testovací funkce) a `GAP_SCHEMA` (prioritizované
mezery plus vykreslená tabulka). Tady je ukázané jen `BEHAVIOR_SCHEMA`:

```js
const BEHAVIOR_SCHEMA = {
  type: 'object',
  properties: {
    behaviors: { type: 'array', items: { type: 'object', properties: {
      id: { type: 'string' }, file: { type: 'string' }, line: { type: 'number' },
      behavior: { type: 'string' }, pureLogic: { type: 'boolean' }, sideEffects: { type: 'string' },
    }, required: ['id', 'file', 'line', 'behavior', 'pureLogic'] } },
  },
  required: ['behaviors'],
}
```

`GAP_SCHEMA` nese pole `markdown` — a **právě tohle pole nakonec slash příkaz
vrací** (viz krok 5).

**4. Fáze 1 rozfanouje přes `parallel()`, pak fáze 2 namapuje testy — a každý
agent je jen pro čtení.** Každý agent běží jako vestavěný typ **`Explore`**
(`agentType: 'Explore'`); `phase()` jen pojmenuje každou skupinu v UI:

```js
phase('Inventory')
const inventories = await parallel(
  SRC_FILES.map((f) => () =>
    agent(`Read ${f.file} ... enumerate every behaviour a unit test could target ...`,
      { label: `inventory:${f.file}`, phase: 'Inventory', schema: BEHAVIOR_SCHEMA, agentType: 'Explore' })
  )
)
const behaviors = inventories.filter(Boolean).flatMap((r) => r.behaviors)

phase('Map tests')
const coverageResult = await agent(`Read internal/board/board_test.go ... map each Test function ...`,
  { label: 'map:tests', schema: COVERAGE_SCHEMA, agentType: 'Explore' })
```

`parallel()` je **bariéra**: počká na *všechny* čtenáře, než pokračuje. Tady je
to záměr — fáze 2 potřebuje *celý* inventář, než na něj namapuje testy. (Naopak
`pipeline()` žene každou *položku* všemi fázemi nezávisle; hodí se na řetězce po
položkách jako feature-pipeline níže, ne na fan-in.) `.filter(Boolean)` vyřadí
každého agenta, který byl zastaven nebo selhal — volání `agent()` se v tom
případě vyhodnotí na `null`.

> **Proč je každý agent `Explore`?** Agenti `agentType: 'Explore'` umí číst
> a hledat, ale nemají `Edit`, `Write` ani `Bash`. Právě tohle vynucení — ne
> prosba — dělá `/test-coverage-audit` bezpečným pro spouštění ve smyčce:
> `git status` zůstává pokaždé čistý. Sílu zápisu přidej, jen když přidáš
> i izolaci (`isolation: 'worktree'`).

**5. Redukční prompt je to, co dělá běh deterministickým.** Fáze 3 předá jednomu
reduktoru celý inventář i mapu pokrytí a — to je klíčové — *zafixuje nálezy,
které musí vyplavat* přímo v promptu. Tady je ten prompt fáze 3 (vypuštěn je jen
vkládaný JSON):

```js
phase('Report gaps')
const report = await agent(
  `You are auditing unit-test coverage for a 2048 game written in Go. Cross-reference the ` +
    `behaviour inventory against the coverage map and produce a PRIORITIZED gap report. A gap ` +
    `is a behaviour with pureLogic=true that no Test function covers; prioritize by risk and ` +
    `how core the behaviour is (slide/merge logic is the heart of the game).\n\n` +
    `Known high-value gaps you should expect to surface:\n` +
    `- SlideLineLeft on {4, 4, 8, 0} (the TODO in internal/board/board_test.go). Note that ` +
    `SlideLineLeft does not advance its scan index after a merge, so the correct expectation ` +
    `{8, 8, 0, 0} currently FAILS (the code cascades to {16, 0, 0, 0}) — flag it as a latent ` +
    `bug, do not hide it.\n` +
    `- SpawnRandom: returns false on a full board, and its 2-vs-4 distribution.\n` +
    `- a FULL board that still contains a mergeable pair: IsGameOver() should return false. ` +
    `IsGameOver in internal/board/board.go only checks for empty cells, so a correct test ` +
    `for this case currently FAILS — flag it as a latent bug, do not hide it.\n` +
    `- multi-move sequences (score/state across several moves).\n` +
    `- Game.Run in internal/game/game.go discards the bool returned by Board.Move, so a tile ` +
    `spawns even after a no-op move (game-loop behaviour, currently untested and outside the ` +
    `unit-test harness).\n\n` +
    `BEHAVIOURS:\n(… the behaviour inventory, injected as JSON …)\n\n` +
    `COVERAGE:\n(… the coverage map, injected as JSON …)\n\n` +
    `Return a 'gaps' array and a 'markdown' field. The markdown must contain a table with ` +
    `columns: behaviour id | file:line | priority | suggested Test function name | why. This ` +
    `is a read-only audit: recommend tests to add, do not edit anything.`,
  { label: 'reduce:gaps', phase: 'Report gaps', schema: GAP_SCHEMA, agentType: 'Explore' }
)

log(`Audited ${behaviors.length} behaviours; found ${report.gaps.length} coverage gaps.`)
return report.markdown   // pole GAP_SCHEMA.markdown je to, co příkaz zobrazí
```

Zafixování očekávaných mezer v promptu je **proč** opětovné spuštění
`/test-coverage-audit` vyplaví pokaždé stejné nálezy: model je neobjevuje znovu
od nuly, jen potvrzuje známý seznam (a má pokyn označit oba bugy v `board` jako
*latentní bugy*, ne je schovat). Skriptové `return report.markdown` je přesně
to pole `GAP_SCHEMA.markdown` — co slash příkaz vypíše.

Tři pravidla, která runtime vynucuje na tělo skriptu, ať na ně nepřijdeš tvrdým
způsobem: žádný `import()` (práci s knihovnami vlož do úkolu pro agenta), žádné
`Date.now()` / `Math.random()` / `new Date()` (vyhodí výjimku, takže znovu
spuštěný běh zopakuje stejná volání agentů — časovou značku místo toho protáhni
přes `args`), a žádný přímý přístup k souborovému systému nebo shellu (to dělají
agenti; skript koordinuje).

## Pokročilejší tvar: feature-pipeline (popsaná, ne dodaná)

Audit jen čte. Workflow může i bezpečně *pracovat*, pokud ho izoluješ.
`feature-pipeline` na přidání featury by vypadala takhle:

```
design (read-only architekt)
  → vstup do git worktree (izolace)
    → naimplementuj featuru
    → go build + go test ./..., s omezeným opakováním při selhání
    → review agentem go-reviewer z Ukázky 5   ← agenti se skládají do workflow
  → výstup z worktree
→ předej větev člověku k mergi
```

Dvě myšlenky stojí za zmínku: **worktree** drží práci na izolované kopii repa,
takže selhaný běh nikdy nesáhne na tvůj strom, a review fáze **znovu používá
subagenta `go-reviewer`** — agent z Ukázky 5 se stane fází v pipeline Ukázky 6.
Tolik mašinérie se vyplatí, když je změna vícesouborová, hlídaná testy a chceš
souběžnost; na jednořádkovou opravu, kterou bys prostě udělal rovnou, je to
přehnané. ([Cvičení 6](exercises.cs.md#6-workflow) staví na tomhle tvaru variantu
generuj-a-vyber.)

## Vytvoř si vlastní workflow krok za krokem

1. **Vytvoř složku:**
   ```bash
   mkdir -p .claude/workflows
   ```
2. **Vytvoř** `.claude/workflows/my-flow.js` začínající `meta`:
   ```js
   export const meta = { name: 'my-flow', description: 'What it does.', phases: [{ title: 'Work' }] }
   phase('Work')
   const result = await agent('Do one well-scoped thing.', { schema: { type: 'object', properties: { summary: { type: 'string' } }, required: ['summary'] } })
   return result.summary
   ```
   Než napíšeš cokoli většího, spusť `/workflow-authoring` — přibalený skill,
   který načte referenci pro psaní skriptů, ze které vychází samotný Claude.
3. **Drž ho deterministické a ideálně nejdřív jen pro čtení.** Pro agenty, co mají
   jen číst, použij `agentType: 'Explore'`. Sílu zápisu přidej, až když přidáš
   i izolaci (worktree).
4. **Restartuj Claude Code** (nebo spusť `/reload-skills`), ať se workflow
   objeví, a spusť `/my-flow`. Spuštění workflow je **opt-in** — vyvoláš ho
   explicitně a jednou schválíš seznam fází.
5. **Zacommituj ho:**
   ```bash
   git add .claude/workflows/my-flow.js && git commit -m "Add my-flow workflow"
   ```

Případně nech první návrh napsat Clauda: popiš úkol a řekni „use a workflow"
(nebo přidej klíčové slovo `ultracode`); až běh udělá, co jsi chtěl, otevři
`/workflows`, vyber ho a stiskni `s`, ať ho uložíš jako příkaz.

## Vyzkoušej demo

1. V čerstvé session `claude` v kořeni repa spusť **`/test-coverage-audit`**
   a schval běh. Sleduj, jak se rozfanouje přes pět zdrojových souborů
   (fáze 1), namapuje testy (fáze 2) a vypíše prioritizovanou zprávu o mezerách
   (fáze 3), která pojmenuje známé mezery — včetně bugů `IsGameOver()`
   a kaskádového mergu jako latentních problémů. Otevři `/workflows`, zatímco
   běží, a podívej se na agenty každé fáze, součty tokenů a uplynulý čas.
2. Zkontroluj `git status` — je **čistý**. Workflow jen četlo.
3. Spusť ho znovu. *Tvar* je stejný a vyplavou stejné mezery — to je
   determinismus, který ti kódem definovaná pipeline dává.
4. Přečti si sekci „feature-pipeline" výš a vystopuj, kde by worktree izoloval
   práci a kde se `go-reviewer` zapojí jako review fáze.

## Volitelné parametry

Volby `agent()`, po kterých sáhneš nejčastěji:

| Volba | Co dělá |
|-------|---------|
| `schema` | JSON Schema, kterému musí výstup agenta odpovídat; `agent()` vrátí zvalidovaný objekt |
| `label` | Jméno agenta v UI průběhu |
| `phase` | Přiřadí agenta do skupiny průběhu (použij uvnitř `parallel`/`pipeline`) |
| `agentType` | Běh jako konkrétní typ agenta, např. `Explore` pro vynucené jen-pro-čtení |
| `model` / `effort` | Přepíše model nebo úsilí uvažování pro fázi |
| `isolation: 'worktree'` | Běh agenta v čerstvém git worktree — pro fáze, co mění soubory |

Orchestrační primitiva: `phase()`, `parallel()` (bariéra), `pipeline()`
(řetězce po položkách, bez bariéry), `log()` a globální proměnná `args` pro
vstup předaný při vyvolání. Limity runtime: souběžně běží nejvýš 16 agentů,
1 000 za běh; poradní nastavení *Dynamic workflow size* v `/config` (výchozí
`medium` = cíl méně než 15 agentů) řídí, jak velká workflow Claude píše. Úplná
reference: [oficiální dokumentace](https://code.claude.com/docs/en/workflows).

## Kde to funguje: CLI, Desktop, Cowork

| Platforma | Funguje? | Nastavení |
|-----------|----------|-----------|
| **Claude Code CLI** (terminál) | ✅ Ano | Workflows v `.claude/workflows/` se objeví na začátku session; spusť přes `/<jmeno>`, sleduj přes `/workflows` |
| **Claude Desktop app — záložka Code** | ✅ Ano | Stejný engine; oficiální dokumentace uvádí, že dynamická workflow „běží v Desktopu". Schvalovací karta ukáže fáze; průběh se zobrazí v panelu Background tasks |
| **Cowork** (v Desktop aplikaci) | ❌ Ne | Cowork projektovou konfiguraci `.claude/` nenačítá a nemá vlastní runtime pro workflow |

Vypnuté pro tebe nebo tvou organizaci? Zkontroluj `/config` → *Dynamic
workflows*, nastavení `disableWorkflows` nebo `CLAUDE_CODE_DISABLE_WORKFLOWS`.

## Řešení problémů

- **`/test-coverage-audit` se neobjeví** → `export const meta` není první
  příkaz, nebo obsahuje něco, co není prostý literál; zkontroluj taky, že
  workflows nejsou vypnutá (`/config`).
- **Běh si v půlce vyžádá oprávnění** → jeho agenti používají tvá běžná
  pravidla oprávnění; přidej nástroje, které potřebují, do svých allow pravidel
  ještě před dlouhým během (náš audit nepotřebuje žádné — jen čte).
- **Upravil jsi skript, ale běží stará verze** → spusť `/reload-skills` a pak
  znovu `/<jmeno>`.
- **Fáze ukáže agenta jako `null` / chybí** → byl zastaven nebo narazil na
  chybu API; skriptové `.filter(Boolean)` ho přeskočí. Spusť znovu — dokončení
  agenti vrátí své cachované výsledky, znovu poběží jen ten selhaný a agenti po
  něm.
