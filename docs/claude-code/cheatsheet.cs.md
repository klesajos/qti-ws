> 🌍 Číst v jazyce: [English](cheatsheet.md) | **Česky**

# Claude Code tahák

Číslované průvodce tě učí, jak Claude Code **rozšířit** ([skills](01-skills.cs.md),
[hooks](02-hooks.cs.md), [MCP](03-mcp.cs.md), [pluginy](04-plugins.cs.md),
[subagenti](05-agents.cs.md), [workflows](06-workflows.cs.md)). Tahle stránka je
ta druhá půlka: jak ho
**ovládat** každý den — přepínače, slash příkazy, klávesové zkratky, prefixy
a události hooků, po kterých sáhneš pořád, ale nezaslouží si každý vlastní
průvodce.

> ✅ **Ověřeno proti Claude Code 2.1.257 (2026-09-03).** Rozhraní CLI se mezi
> verzemi mění — zkontroluj `claude --version` a
> [oficiální dokumentaci](https://code.claude.com/docs/en/cli-reference), pokud
> tu něco nesedí.

## Spouštěcí přepínače (`claude ...`)

| Příkaz | Co dělá |
|--------|---------|
| `claude` | Spustí interaktivní session v aktuální složce |
| `claude "oprav lint chyby"` | Spustí session s prvním promptem už ve frontě |
| `claude -c` / `--continue` | Obnoví **poslední** session v této složce |
| `claude -r` / `--resume [id]` | Obnoví **konkrétní** session (bez id → interaktivní výběr); `-n <name>` pojmenuje session, abys ji mohl obnovit podle jména |
| `claude -p "..."` / `--print` | **Headless**: spustí prompt, vypíše výsledek, skončí. Pro skripty / CI |
| `claude --model opus` | Spustí na konkrétním modelu (`fable`, `opus`, `sonnet`, `haiku`, `opusplan`, `best`) |
| `claude --effort high` | Spustí na dané úrovni uvažování (`low` … `max`, nebo `ultracode`) |
| `claude --agent go-reviewer` | Spustí celou session jako pojmenovaného subagenta |
| `claude --add-dir ../lib` | Dá session přístup k dalším složkám mimo repo |
| `claude --worktree` | Spustí v čerstvém git worktree, aby se session nemohla dotknout tvého checkoutu |
| `claude --permission-mode plan` | Spustí v režimu oprávnění (`manual`/`default` · `acceptEdits` · `plan` · `auto` · `dontAsk` · `bypassPermissions`) |
| `claude -p "..." --allowedTools "Read,Edit,Bash(git diff *)"` | Allowlist nástrojů, které smí bezobslužný / CI běh použít |
| `claude -p "..." --output-format json` | Headless výstup jako `json` nebo `stream-json` pro skripty |
| `claude --restricted` | Uzamčená session pro sdílené stroje: žádné příkazové nástroje, přístup k souborům omezen na pracovní adresáře |
| `claude --safe-mode` / `--bare` | Řešení problémů: běh bez jakýchkoli úprav (CLAUDE.md, hooks, skills, plugins, MCP) |

## Session na pozadí (`claude agents`)

**Session na pozadí** je plnohodnotná konverzace, která běží odpojená (vlastní ji
supervizor proces) — zadáš práci, odejdeš a vrátíš se k ní později.

| Příkaz | Co dělá |
|--------|---------|
| `claude --bg "..."` | Spustí session na pozadí ze shellu (dlouhý tvar `--background`) |
| `/bg <prompt>` | Pošle aktuální práci na pozadí zevnitř session (taky `/background`) |
| `claude agents` | Otevře agent view — monitoruj / spouštěj session (`--json` pro skriptování) |
| `claude attach <id>` | Připojí se k běžící session na pozadí |
| `claude logs <id>` | Vypíše výstup session na pozadí |
| `claude stop <id>` | Zastaví session (alias `claude kill`) |
| `claude respawn <id>` | Restartuje zastavenou session se zachovanou konverzací (`--all` pro všechny) |
| `claude rm <id>` | Odebere dokončenou session ze seznamu |
| `claude daemon status` | Ukáže supervizor, který vlastní session na pozadí (`daemon stop --any` ho zastaví) |

**Headless (`-p`) vs. pozadí (`--bg`) — kdy co použít:**

| | `claude -p "..."` (headless) | `claude --bg "..."` (pozadí) |
|---|---|---|
| Běží | Na popředí — vypíše a skončí | Odpojeně — přetrvává pod supervizorem |
| Vázané na tvůj shell | Ano (zavřeš ho → konec) | Ne (běží dál) |
| V `claude agents`? | Ne | Ano (attach / logs / stop) |
| Sáhni po něm když | Skriptování, CI, jednorázovka k pipe/parsování | Dlouhý úkol, který zadáš a vracíš se k němu při další práci |

Neplést s `/agents` (který teď jen připomene, ať upravíš `.claude/agents/` —
interaktivní průvodce byl odstraněn) ani s `claude --agent <název>` (spustí
celou session *jako* pojmenovaného subagenta).

## Slash příkazy

Napiš `/` pro automatické doplnění. Seskupené podle účelu:

**Session a kontext**

| Příkaz | Co dělá |
|--------|---------|
| `/clear` | Vymaže konverzaci — čistý štít pro novou úlohu (aliasy `/new`, `/reset`) |
| `/compact [pokyny]` | Shrne konverzaci a uvolní kontext; volitelné zaměření, např. `/compact zachovej práci na undu` |
| `/context` | Ukáže, co zrovna plní kontextové okno |
| `/btw <dotaz>` | Rychlý vedlejší dotaz — zodpovězený *bez* přidání do historie konverzace |
| `/rewind` | Skočí na dřívější checkpoint (kód i/nebo konverzaci) — taky `Esc Esc` |
| `/resume` | Přepne na jinou minulou session bez opuštění |
| `/rename <název>` | Přejmenuje aktuální session (snadnější dohledání přes `--resume`) |
| `/export` | Uloží celou konverzaci do Markdown souboru |
| `/plan <úkol>` | Rovnou z promptu přejde do plan módu |
| `/subtask <prompt>` | Předá vedlejší úkol subagentovi; výsledek se vrátí do téhle konverzace |
| `/fork` | Zkopíruje konverzaci do nové session na pozadí (vlastní worktree) |
| `/tasks` | Vypíše práci téhle session na pozadí — subagenty, shell úlohy, workflows |
| `/branch [název]` | Rozvětví *konverzaci*, abys vyzkoušel jiný směr |
| `/goal <podmínka>` | Nechá Clauda pracovat, dokud není splněna podmínka |
| `/diff` | Interaktivní prohlížeč nezacommitovaných změn |
| `/cd <cesta>` | Přesune session do jiného adresáře, konverzaci zachová |

**Konfigurace**

| Příkaz | Co dělá |
|--------|---------|
| `/config` | Otevře menu nastavení (model, téma, oprávnění, velikost workflow…) |
| `/model` | Přepne model uprostřed session (`Option+P` otevře výběr) |
| `/effort [úroveň]` | Nastaví úsilí uvažování — `/effort high`, `/effort auto`, `/effort ultracode` |
| `/fast` | Přepne fast mode (stejný model, rychlejší výstup) |
| `/permissions` | Zobrazí / udělí trvalá oprávnění nástrojům |
| `/sandbox` | Přepne izolaci souborového systému / sítě pro Bash (nebo nastav `"sandbox.enabled": true`) |
| `/memory` | Upraví trvalou paměť (soubory `CLAUDE.md`, zapnutí/vypnutí auto-memory) |
| `/init` | Vygeneruje startovní `CLAUDE.md` pro tohle repo |
| `/status`, `/statusline` | Přehled stavu session / úprava spodní info lišty |
| `/doctor` | Diagnostika rozbité instalace (auth, síť, závislosti) |
| `/reload-skills`, `/reload-plugins` | Znovu načte skills, workflows a pluginy bez restartu |

**Rozšíření** (šest mechanismů, které pokrývají průvodce)

| Příkaz | Co dělá |
|--------|---------|
| `/skill-name` | Spustí skill, např. `/board-tests` (Ukázka 1) |
| `/2048-dev:build-test` | Spustí skill zabalený v pluginu (Ukázka 4) |
| `/hooks` | Zobrazí hooky nakonfigurované pro tuhle session |
| `/mcp` | Zobrazí MCP servery, jejich stav, a autorizuje je |
| `/plugin` | Instaluje / spravuje pluginy a marketplaces |
| `/agents` | (Průvodce odstraněn) — uprav `.claude/agents/*.md` přímo, viz Ukázka 5 |
| `/workflows` | Sleduj, pozastav, zastav nebo ulož běžící workflows (Ukázka 6) |
| `/workflow-authoring` | Načte referenci pro psaní workflow skriptu, než nějaký napíšeš |

**Info**

| Příkaz | Co dělá |
|--------|---------|
| `/help` | Vypíše všechny příkazy, zkratky a funkce |
| `/usage` | Zbývající kapacita na předplatném, útrata a statistiky cache (`/cost` je alias) |
| `/powerup` | Interaktivní lekce funkcí |
| `/insights` | HTML report analyzující tvé nedávné session |

**Přibalené skills, které stojí za znalost:** `/code-review` (alias `/review`),
`/security-review`, `/verify`, `/debug`, `/loop`, `/deep-research` (přibalené
workflow), `/fewer-permission-prompts`.

## Klávesové zkratky

| Klávesa | Co dělá |
|---------|---------|
| `Ctrl+C` | Přeruší aktuální akci (nouzový vypínač); dvakrát pro ukončení |
| `Ctrl+R` | Zpětné hledání v historii promptů/příkazů |
| `Esc` | Přeruší Clauda, zruší aktuální vstup / zavře menu |
| `Esc Esc` | Vymaže rozepsaný text — nebo, při prázdném vstupu, otevře **Rewind** |
| `Shift+Tab` | Cykluje režimy oprávnění (Manual → Accept edits → Plan) |
| `Ctrl+G` | Otevře aktuální prompt ve tvém `$EDITOR` |
| `Ctrl+O` | Přepne prohlížeč přepisu (transcript) |
| `Ctrl+T` | Přepne Claudův checklist úkolů |
| `Ctrl+B` | Pošle běžící úkol na pozadí |
| `Ctrl+S` | Odloží / obnoví prompt, který právě píšeš |
| `Option+P` / `Option+T` | Přepne model / zapne rozšířené přemýšlení (macOS; `Alt+` na Windows/Linuxu) |
| `?` na prázdném promptu | Zobrazí panel nápovědy ke zkratkám |

## Vstupní prefixy

| Prefix | Co dělá |
|--------|---------|
| `!` | **Bash režim** — spustí shellový příkaz přímo, bez modelu, bez tokenů (např. `!git status`) |
| `@` | **Zmínka souboru** — vtáhne konkrétní soubor/složku do kontextu (doplňuje se) |
| `@agent-<název>` | **Delegace** — předá úlohu konkrétnímu subagentovi (např. `@agent-go-reviewer`) |
| `:` | **Emoji shortcode** — `:tada:` |
| `\` + Enter | **Víceřádkový vstup** — pokračuj na novém řádku bez odeslání (taky `Ctrl+J`, nebo `Shift+Enter` ve většině moderních terminálů) |

- **Pamatování věcí:** starý prefix `#` pro „přidat tenhle řádek do paměti" už
  není součástí zdokumentovaného rozhraní. Požádej Clauda, ať si něco
  zapamatuje (auto-memory), nebo uprav `CLAUDE.md` přes `/memory`.
- **Vlož obrázek:** `Ctrl+V` (nebo drag-and-drop) hodí screenshot / mockup do
  promptu — skvělé na UI práci. Na macOS použij `Ctrl+V`, *ne* `Cmd+V`.
- **Pipe souboru:** `cat error.log | claude -p "vysvětli tuhle chybu"` pošle stdin
  do headless běhu (potřebuje `-p`).

## Rewind a checkpointy

Každý prompt je **checkpoint**. Menu otevřeš přes `/rewind` nebo `Esc Esc`:

| Obnovit | Účinek |
|---------|--------|
| **Kód + konverzaci** | Vrátí obojí zpět k tomu promptu |
| **Jen konverzaci** | Přetočí chat, aktuální kód nechá |
| **Jen kód** | Vrátí úpravy souborů, chat nechá |
| **Shrnout odsud / sem** | Zkomprimuje konverzaci za / před tímto bodem |

- Checkpointy přetrvávají mezi sessions a drží se ~30 dní (`cleanupPeriodDays`).
- **Nezachytí se:** soubory změněné přes Bash (`rm` / `mv` / `cp`) nebo upravené
  mimo Claude Code. **Není náhrada Gitu** — commituj na milnících (viz
  průvodce [Git basics](00-git-basics.cs.md)).

## Události hooků

Kam se hook může navěsit v životním cyklu session (viz [Ukázka 2](02-hooks.cs.md)):

| Událost | Spustí se | Typické použití |
|---------|-----------|-----------------|
| `SessionStart` | Začátek session | Načtení kontextu, kontroly prostředí |
| `UserPromptSubmit` | Odeslání promptu | Přidání kontextu, validace/přepis promptu |
| `PreToolUse` | Před spuštěním nástroje | **Zablokovat** nebezpečné akce (např. `git push --force`) |
| `PostToolUse` | Po úspěšném nástroji | Formát / lint — *hook `format-go` v tomhle repu* |
| `Stop` | Claude dokončí odpověď | Ověřit testy — *hook `run-tests` v tomhle repu* |
| `SubagentStart` / `SubagentStop` | Subagent se spustí / dokončí | Připrav nebo ukliď, co potřebuje |
| `PreModelSwitch` / `PostModelSwitch` | Model se uprostřed session změní | Audit, opětovné nastavení voleb |
| `SessionEnd` | Konec session | Úklid, archivace, reporting |

Typy handlerů: `command` (shell), `http`, `mcp_tool`, `prompt`, `agent`.

## Režimy oprávnění (cykluj `Shift+Tab`)

| Režim | Claude smí… |
|-------|-------------|
| **Manual** (`default`) | Ptát se před každou úpravou a příkazem — šedý odznak ⏸ v patičce ukazuje, že jsi tady |
| **Accept edits** (`acceptEdits`) | Volně upravovat soubory, u shell příkazů se pořád ptá |
| **Plan** (`plan`) | Jen pro čtení — bádá a sepíše plán, nic nemění |
| **Auto** (`auto`) | Jedná, s klasifikátorem na pozadí, který blokuje riskantní akce místo ptaní |
| **Bypass permissions** | Cokoli bez ptaní („YOLO" — jen v sandboxu). Není v cyklu `Shift+Tab` — zapni ho explicitně přes `--dangerously-skip-permissions` |

Existuje ještě jeden režim — **dontAsk** (jen předem schválené nástroje,
všechno ostatní zamítnuto) — pro bezobslužné běhy; viz
[dokumentace režimů oprávnění](https://code.claude.com/docs/en/permission-modes).

## Modely a effort

| Model | Nejlepší na |
|-------|-------------|
| `fable` | Dlouhé autonomní session, vyšetřování, verifikaci (Fable 5.1, 1M kontext) |
| `opus` | Nejtěžší úvahy — architektura, záludné bugy (Opus 5) |
| `sonnet` | Denní tahoun — featury, refaktory, testy (Sonnet 5) |
| `haiku` | Rychlý a levný — boilerplate, rychlé kontroly, paralelní subagenti |
| `opusplan` | Opus na plán, pak se automaticky přepne na Sonnet na implementaci |
| `best` / `default` | Nejsilnější model, který máš k dispozici / výchozí model tvého účtu |

**Effort** řídí, jak hluboko Claude přemýšlí — vyšší = lepší odpovědi, pomalejší,
víc tokenů. Stupnice od nejnižšího: `low` · `medium` · `high` · `xhigh` · `max`.
`ultracode` je nad tím vším: `xhigh` **plus** Claude ke každému podstatnému
úkolu napíše workflow (viz [Ukázka 6](06-workflows.cs.md)).

| Nastavení effortu | Jak |
|-------------------|-----|
| V session | `/effort` (posuvník) · `/effort high` (skok na úroveň) · `/effort auto` (zpět na výchozí pro model) |
| Při spuštění | `claude --effort xhigh "..."` — jen pro tuhle session |
| Trvale | `"effortLevel": "high"` v `settings.json`, nebo env `CLAUDE_CODE_EFFORT_LEVEL` (`max`/`ultracode` takhle trvale nastavit nejde) |
| Per agent / skill | pole `effort:` ve frontmatteru |

**Jednorázově:** napiš `ultrathink` kamkoli do promptu a ten jeden tah půjde
hlouběji *bez* změny effortu session — jediné platné klíčové slovo je
`ultrathink` (`think` / `think hard` jsou jen běžná slova).

## Ekonomika kontextu (proč na `/compact` záleží)

- Claude čte **celou** konverzaci každý tah, takže cena roste s historií.
  `/compact` shrne, aby zůstala štíhlá; `/clear` resetuje pro novou úlohu.
- **Prompt caching** dává velkou slevu na stabilní prefix (tvůj `CLAUDE.md`,
  struktura projektu) — je automatický; drž ten prefix stabilní.
- **Pozor na rozšíření:** každý aktivní MCP server, skill a subagent ujídá
  kontext. Zapni to, co úloha potřebuje, ne všechno.

## Co je v roce 2026 legacy (sjednocení)

Vlastní slash příkazy byly **sloučeny do Skills** (Ukázka 1). Staré soubory
příkazů fungují dál, ale nové věci dělej přes Skills:

| Bývalo samostatné | Teď | Stav |
|-------------------|-----|------|
| Vlastní slash příkazy (`.claude/commands/*.md`) | **Skills** (`.claude/skills/<name>/SKILL.md`) | Sloučeno — obojí vytváří `/name`; legacy soubory fungují dál |
| Průvodce `/agents` | Uprav `.claude/agents/*.md` ručně (nebo požádej Clauda) | Průvodce odstraněn ve verzi 2.1.198 |
| Nástroj `Task` (spouštění subagentů) | Nástroj `Agent` | Přejmenováno; staré jméno pořád funguje jako alias |
| `/cost` | `/usage` | `/cost` je teď alias |
| `/review` | `/code-review` | Alias |
| Nástroje pro todo / sledování úkolů | U novějších modelů ve výchozím stavu vypnuté | `CLAUDE_CODE_ENABLE_TODO_TOOLS=1` je obnoví |
| `ultraplan` | Odstraněno (2.1.222) | Použij plan mode + `ultrathink` |

(Output styles zůstávají *samostatnou*, stále aktuální funkcí — mění systémový
prompt, zatímco skills načítají instrukce k úloze — viz
[dokumentace output styles](https://code.claude.com/docs/en/output-styles);
deprecován byl jen samostatný příkaz `/output-style` ve prospěch `/config`.)

Skill je striktní nadmnožina starého příkazu: stejné vyvolání `/name`, **plus**
volitelné autonomní načtení, složka pro pomocné soubory a řízení vyvolání. Chceš,
aby se skill choval jako starý „jen ručně" příkaz? Přidej do frontmatteru
`disable-model-invocation: true`.

Zdroj: [oficiální dokumentace Skills](https://code.claude.com/docs/en/skills) —
*„Custom commands have been merged into skills."*

---

**Měj tenhle tahák po ruce při práci.** Je to reference, ne tutoriál — pro *jak*
a *proč* každého mechanismu rozšíření viz číslované průvodce v
[indexu](README.cs.md).
