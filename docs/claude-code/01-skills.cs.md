> 🌍 Číst v jazyce: [English](01-skills.md) | **Česky**

# Ukázka 1: Projektový skill

## Co je skill?

**Skill** je Markdown soubor s instrukcemi pro Clauda, který se načte
**jen když je relevantní**. Srovnej s `CLAUDE.md`:

- `CLAUDE.md` se načítá do **každé** konverzace — drž ho krátký.
- Skill se načítá **na vyžádání** — když tvůj požadavek odpovídá jeho
  popisu, nebo když ho vyvoláš jménem přes `/nazev-skillu`.

Claudův kontext tak zůstává malý, ale hluboké znalosti o projektu má
k dispozici, kdykoli je potřebuje.

## Co dělá tahle ukázka

Skill `board-tests` učí Clauda, *jak se v tomhle projektu píšou unit testy*:
která API použít, jaký styl dodržet, jak spustit testovací sadu. Bez něj by
Claude hádal — a každý účastník workshopu by dostal trochu jiné testy
(někdo s `testify`, někdo s `reflect.DeepEqual`, někdo testující přes herní
smyčku).

## Soubor řádek po řádku

Skill žije v `.claude/skills/board-tests/SKILL.md`. Na cestě záleží:

```
.claude/              ← projektová konfigurace Claude Code (commitnutá v gitu)
└── skills/           ← tady žijí všechny projektové skills
    └── board-tests/  ← jedna složka na skill; název složky = název skillu
        └── SKILL.md  ← musí se jmenovat přesně SKILL.md
```

Soubor začíná **frontmatterem** — malým YAML blokem mezi řádky `---`,
který skill popisuje:

```yaml
---
name: board-tests
description: Use when writing or extending Go unit tests for the 2048 board
  logic (SlideLineLeft, Board.Move, win/game-over detection)...
---
```

Co dělá který řádek:

- `---` — otevírá (a později zavírá) blok frontmatteru. Povinné.
- `name: board-tests` — identifikátor skillu. Tohle píšeš za lomítko:
  `/board-tests`. Nepovinné — když chybí, použije se název složky.
- `description: ...` — **nejdůležitější řádek.** Claude si při startu
  session přečte jen popisy všech skillů. Když tvůj požadavek odpovídá
  popisu, Claude načte celé tělo skillu. Piš popis jako spouštěč: řekni,
  *kdy* se má skill použít, ne co je uvnitř.

Všechno **pod** frontmatterem je tělo — obyčejné Markdown instrukce, kterými
se Claude řídí, jakmile je skill aktivní. U nás: konvence testů
(Arrange-Act-Assert, jedno chování na jednu funkci `Test...`, `==` na
polích, žádné assertion knihovny), klíčové API `board` a přesné příkazy
`go test` pro ověření.

## Vytvoř si vlastní skill, krok za krokem

1. **Vytvoř složku.** V kořenu repa:
   ```bash
   mkdir -p .claude/skills/my-skill
   ```
   (`mkdir -p` vytvoří celou cestu najednou, včetně nadřazených složek.)

2. **Vytvoř soubor** `.claude/skills/my-skill/SKILL.md` s touto kostrou:
   ```markdown
   ---
   name: my-skill
   description: Use when <describe the situation that should trigger this skill>.
   ---

   # My skill

   <Instructions for Claude. Write them like you'd brief a new colleague:
   what to do, which files/APIs to use, how to verify the result.>
   ```

3. **Restartuj Claude Code** (nebo začni novou session), aby nový skill
   našel. Skills se načítají při startu session. (Ve spuštěné session
   znovu přečte složky skillů bez restartu příkaz `/reload-skills`.)

4. **Otestuj ruční spuštění:** napiš `/my-skill` — tělo skillu se načte
   do konverzace.

5. **Otestuj automatické spuštění:** zformuluj požadavek, který odpovídá
   tvému `description`. Claude by měl před akcí oznámit, že skill používá.

6. **Commitni ho:**
   ```bash
   git add .claude/skills/my-skill
   git commit -m "Add my-skill"
   ```
   Teď ho má každý, kdo si repo naklonuje. Přesně tohle znamená „projektový".
   (Osobní skills, které sdílet nechceš, patří do `~/.claude/skills/`
   v domovském adresáři.)

## Vyzkoušej demo

Zeptej se Clauda: *„Přidej test pro posun řádku {4, 4, 8, 0}"*.

V `internal/board/board_test.go` je záměrně `TODO (participants)` přesně
pro tenhle případ. Díky skillu Claude zavolá přímo `SlideLineLeft(&line)`
(potřebuje skóre), pojmenuje funkci `TestSlideLineLeft_...`, porovná `Line`
pomocí `!=` a ověří výsledek přes `go test` — u každého účastníka stejně.
Sleduj, co test hlásí: řádek *by se měl* stát `{8, 8, 0, 0}` a skill říká
Claudovi nahlásit selhání, ne „opravit" očekávání.

## Volitelné parametry frontmatteru

Náš skill používá jen `name` a `description`, ale frontmatter umí víc.
Nejužitečnější volitelná pole:

| Pole | Co dělá |
|------|---------|
| `argument-hint` | Nápověda v našeptávači za slash příkazem, např. `[test-name]` |
| `disable-model-invocation: true` | Claude skill nikdy nenačte sám — jen ty přes `/nazev`. Vhodné pro destruktivní workflow (deploy, release). Zároveň brání předběžnému načtení skillu do subagentů |
| `user-invocable: false` | Opak: skrytý v `/` menu, načíst ho může jen Claude |
| `allowed-tools` | Nástroje, které skill smí použít bez ptaní během tahu, ve kterém byl vyvolán, např. `Bash(go test:*)` |
| `model` / `effort` | Po dobu aktivního skillu použít konkrétní model nebo úroveň reasoning effortu |
| `context: fork` | Spustit skill v izolovaném subagentovi, aby nezaplnil kontext hlavní konverzace (`agent:` volí typ subagenta; ve výchozím stavu běží na pozadí) |
| `paths` | Glob vzory — automaticky se spustí jen u odpovídajících souborů, např. `internal/board/**` |
| `hooks` | Hooks registrované při vyvolání skillu (viz [Ukázka 2](02-hooks.cs.md)); jediné místo, kde se respektuje volba hooku `once: true` |

Úplný aktuální seznam: [oficiální dokumentace skills](https://code.claude.com/docs/en/skills).

## Kde to funguje: CLI, Desktop aplikace, Cowork

| Platforma | Funguje? | Nastavení |
|-----------|----------|-----------|
| **Claude Code CLI** (terminál) | ✅ Ano | Nic navíc — skills v `.claude/skills/` se najdou automaticky, když spustíš `claude` v repu |
| **Claude Desktop — záložka Code** | ✅ Ano | Otevři složku projektu v záložce Code. Desktop běží na stejném enginu jako CLI a sdílí veškerou projektovou konfiguraci. Poprvé potvrdíš dialog důvěry projektu |
| **Cowork** (v Desktop aplikaci) | ❌ Ne | Cowork pouští úlohy ve vlastním sandboxovaném VM a projektovou konfiguraci `.claude/` **nenačítá**. Náhrada: zabal skill do pluginu a nainstaluj ho v Coworku, nebo ho přidej jako skill přes claude.ai |

Zdroj: [dokumentace k Desktopu](https://code.claude.com/docs/en/desktop#coming-from-the-cli) —
„Desktop and CLI read the same configuration files ... Hooks and skills
defined in settings apply to both." (Desktop a CLI čtou stejné konfigurační
soubory ... Hooks a skills definované v nastavení platí pro obojí.)

## Když něco nefunguje

- **`/board-tests` neexistuje** → soubor musí být přesně
  `.claude/skills/<nazev>/SKILL.md` (velkými `SKILL.md`); restartuj session
  nebo spusť `/reload-skills`.
- **Skill se nikdy nespustí automaticky** → tvůj `description` je moc
  vágní. Přepiš ho tak, aby popisoval *situaci uživatele* („Use when
  writing unit tests...") místo obsahu („Informace o testech").
- **Claude přesto napsal test s `testify`** → skill se nenačetl;
  zkontroluj `/context`, jestli je skill uvedený v seznamu, pak ho vyvolej
  ručně přes `/board-tests` a zkus to znovu.
</content>
