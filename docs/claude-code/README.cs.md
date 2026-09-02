> 🌍 Číst v jazyce: [English](README.md) | **Česky**

# Ukázky rozšíření Claude Code

Tohle repo zároveň slouží jako ukázka šesti způsobů, jak rozšířit Claude Code
na úrovni **projektu**. „Na úrovni projektu" znamená, že konfigurace žije
**přímo v repozitáři** — když si projekt naklonuješ přes `git clone`,
dostaneš celé nastavení automaticky. Nic se neinstaluje do domovského
adresáře.

Každá ukázka je samostatná a minimální. Projdi je popořadě — jdou od
nejjednodušší po nejpokročilejší:

| # | Ukázka | Co tě naučí | Kde žije |
|---|--------|-------------|----------|
| 1 | [Skills](01-skills.cs.md) | Jak dát Claudovi znalosti o projektu, které si načte, jen když je potřebuje | `.claude/skills/board-tests/SKILL.md` |
| 2 | [Hooks](02-hooks.cs.md) | Jak automaticky spustit shellový příkaz, když Claude něco udělá (tady: autoformátování kódu, spuštění testů) | `.claude/settings.json` + `.claude/hooks/*.sh` |
| 3 | [MCP](03-mcp.cs.md) | Jak Clauda připojit k externím nástrojům a zdrojům dat | `.mcp.json` |
| 4 | [Pluginy](04-plugins.cs.md) | Jak zabalit skills/agenty/hooks, aby je mohl sdílet celý tým | `plugins/2048-dev/` + `.claude-plugin/marketplace.json` |
| 5 | [Subagenti](05-agents.cs.md) | Jak delegovat samostatný úkol do izolovaného kontextu s vlastními nástroji a personou | `.claude/agents/` + `plugins/2048-dev/agents/` |
| 6 | [Workflows](06-workflows.cs.md) | Jak spustit víc agentů v pevně daném, opakovatelném pořadí definovaném v kódu | `.claude/workflows/` |

Každý návod existuje ve dvou jazycích: `xx-nazev.md` je anglicky,
`xx-nazev.cs.md` česky. Obsah je stejný. Až projdeš všech šest, [katalog
cvičení](exercises.cs.md) ti dá jednu herní featuru k postavení na každý
mechanismus (plus bonus: [přepiš si původní C++ verzi sám do
Go](00-migration.md)) — a [tahák](cheatsheet.cs.md) je rychlá reference na
každodenní přepínače, slash příkazy a zkratky.

## Než začneš

Tři krátké návody ti připraví stroj. Projdi je **před dnem workshopu**
— každý končí kontrolou, kterou si můžeš spustit:

1. [Instalace Claude Code](00-install-claude-code.md) — samotné CLI plus
   předpoklady, které tohle repo potřebuje (Git, Node.js kvůli `npx`, Go),
   pro macOS, Windows a Linux.
2. [Základy Gitu](00-git-basics.md) — hrstka příkazů, které workshop
   používá.
3. [Spuštění Go aplikace](00-go-basics.md) — sestavení, spuštění a
   otestování tohoto projektu.

Pak:

- **Naklonuj tohle repo** — `git clone https://github.com/klesajos/qti-ws`
- **Otevři terminál ve složce repa** — všechny ukázky předpokládají, že
  pracovní adresář je kořen repa (složka s `go.mod`).
- **Spusť jednou `go test ./...`** — musí projít zeleně dřív, než začneš;
  první spuštění zároveň stáhne jedinou závislost.

## Rychlé demo (5 minut)

Otevři Claude Code v repu (`cd qti-ws && claude`), potvrď dialog důvěry
a vyzkoušej:

1. **Skill** — zeptej se: *„Přidej test pro posun řádku {4, 4, 8, 0}"*.
   Sleduj, jak si Claude před napsáním testu načte skill `board-tests`.
   Můžeš ho vyvolat i ručně napsáním `/board-tests`.
2. **Hook** — požádej Clauda, ať upraví libovolný `.go` soubor. Hned po
   úpravě hook spustí na souboru `gofmt` — kód je vždycky zformátovaný,
   i kdyby ho Claude napsal rozházeně. Když Claude dokončí tah, druhý hook
   spustí `go test ./...` a vypíše jednořádkový verdikt.
3. **MCP** — napiš `/mcp`, uvidíš dva projektové servery. Pak se zeptej:
   *„Použij deepwiki a zjisti, jak fungují table-driven testy a t.Run
   subtesty v golang/go."*
4. **Plugin** — napiš `/2048-dev:build-test` a jedním příkazem se celý
   projekt sestaví, provede `go vet`, zkontroluje formátování a otestuje.
5. **Subagent** — udělej drobnou úpravu libovolného `.go` souboru a
   zeptej se: *„Zkontroluj moje necommitnuté změny."* Read-only agent
   `go-reviewer` nahlásí nálezy s `file:line` a nic nezmění.
6. **Workflow** — spusť `/test-coverage-audit`, rozfanouje read-only audit
   pokrytí přes zdrojové soubory a vypíše prioritizovanou zprávu o mezerách
   (spusť znovu: stejný tvar).

## Podpora platforem v kostce

Claude Code běží na několika místech a **ne všechno funguje všude**.
Každý návod má podrobnou sekci „Kde to funguje"; shrnutí:

| Projektový mechanismus | CLI (terminál) | Desktop — záložka Code | Cowork |
|------------------------|:--------------:|:----------------------:|:------:|
| Skills (`.claude/skills/`) | ✅ | ✅ | ❌ — zabal raději do pluginu |
| Hooks (`.claude/settings.json`) | ✅ | ✅ | ❌ — sandbox lokální hooks nespouští |
| MCP (`.mcp.json`) | ✅ | ✅ | ❌ — použij **Konektory** na claude.ai |
| Plugin (in-repo marketplace) | ✅ | ✅ | ⚠️ — obsah funguje, ale instaluje se přes správu pluginů v Coworku, ne přes projektové nastavení |
| Subagenti (`.claude/agents/`) | ✅ | ✅ | ❌ — zabal agenta raději do pluginu |
| Workflows (`.claude/workflows/`) | ✅ | ✅ | ❌ — lokální orchestrace; Cowork nemá ekvivalent |

Proč: CLI a záložka **Code** v Desktopu běží na stejném enginu a sdílejí
veškerou projektovou konfiguraci — Desktop jen přidává jednorázový dialog
důvěry projektu. **Cowork** pouští úlohy ve vlastním sandboxovaném VM
a projektovou konfiguraci `.claude/` nenačítá; má vlastní ekvivalenty
(Konektory pro MCP, vlastní instalaci pluginů).

Reference: [dokumentace k Desktopu](https://code.claude.com/docs/en/desktop#coming-from-the-cli)
uvádí, že Desktop a CLI „read the same configuration files ... CLAUDE.md
... MCP servers ... Hooks and skills ... Settings" (čtou stejné konfigurační
soubory — CLAUDE.md, MCP servery, hooks a skills, nastavení) a že dynamické
workflows „run in Desktop" (běží i v Desktopu). Každý číslovaný průvodce
níže má vlastní sekci „Kde to funguje" s odkazy na příslušnou oficiální
dokumentaci.

## Kdy použít který mechanismus?

Častá otázka začátečníků. Základní pravidlo:

- **Skill** — když chceš, aby Claude něco *věděl* (konvence, postupy).
  Kdy znalost použít, rozhoduje Claude. Ber to jako dokumentaci pro AI.
- **Hook** — když se něco musí stát *pokaždé*, bez výjimky (formátování,
  lint, blokování nebezpečných příkazů). Model do toho nemá co mluvit.
- **MCP** — když Claude potřebuje *schopnost*, kterou nemá (dotaz do
  databáze, hledání v dokumentaci, ovládání prohlížeče).
- **Plugin** — když chceš cokoli z výše uvedeného *sdílet* mezi více repy
  nebo s celým týmem jako jeden verzovaný balíček.
- **Subagent** — když chceš předat *jeden* samostatný úkol izolovanému
  kontextu s vlastními nástroji a personou (recenzent, co nemůže zapisovat,
  tester, průzkumník kódu).
- **Workflow** — když chceš *opakovatelnou, vícefázovou* orchestraci *víc*
  agentů v kódu (rozfanoutovat přes soubory, podmínit jednu fázi druhou,
  zredukovat víc výsledků do jedné zprávy).
</content>
