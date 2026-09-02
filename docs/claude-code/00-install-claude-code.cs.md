> 🌍 Číst v jazyce: [English](00-install-claude-code.md) | **Česky**

# Před workshopem: instalace Claude Code a předpokladů

Udělej tohle **před dnem workshopu** — s dobrým připojením to zabere 15 minut
a poslední sekce je checklist, který ti řekne, že jsi hotový. Ověřeno proti
Claude Code **2.1.257** a oficiální [instalační stránce](https://code.claude.com/docs/en/setup)
dne 2026-09-03.

## Co potřebuješ

| | Proč | Verze |
|---|---|---|
| **Účet Claude** s přístupem ke Claude Code | Pro, Max, Team, Enterprise, nebo účet Console (API). Free plán Claude Code **neobsahuje** | — |
| **Claude Code** | Samotný agent. Nativní binárka — ke spuštění *nepotřebuje* Node.js | ≥ 2.1.257 |
| **Git** | Klonování tohohle repa, commitování práce; na Windows navíc dodává Bash, který používá Claude Code i hooks v tomhle repu | jakákoli aktuální |
| **Node.js + npm/npx** | Jen proto, že `.mcp.json` tohohle repa spouští MCP server `memory` přes `npx` ([Ukázka 3](03-mcp.cs.md)) | **22+** |
| **Go** | Sestavení a testování appky — viz [Spuštění Go appky](00-go-basics.cs.md) | **1.25+** |
| `jq` *(volitelné)* | Hezčí práce s JSON v hooks; bez něj se použije záložně `python3` nebo `sed` | jakákoli |

Podporované OS: macOS 13+, Windows 10 1809+ / Server 2019+, Ubuntu 20.04+ /
Debian 10+. 4 GB RAM. Přístup k internetu (první běh navíc stáhne Go moduly
a npm balíček).

## 1. Instalace Claude Code

Vyber **jednu** metodu. Doporučený je nativní instalátor — aktualizuje se
sám na pozadí.

**macOS / Linux / WSL**

```bash
curl -fsSL https://claude.ai/install.sh | bash
```

**Windows — PowerShell**

```powershell
irm https://claude.ai/install.ps1 | iex
```

**Windows — CMD**

```bat
curl -fsSL https://claude.ai/install.cmd -o install.cmd && install.cmd && del install.cmd
```

> Nejsi si jistý, v jakém Windows shellu jsi? Prompt v PowerShellu ukazuje
> `PS C:\Users\ty>`, v CMD (bez `PS`) `C:\Users\ty>`. Pokud vidíš
> `The token '&&' is not a valid statement separator`, jsi v PowerShellu —
> použij příkaz pro PowerShell. Pokud vidíš `'irm' is not recognized`, jsi
> v CMD.

Alternativy, pokud preferuješ správce balíčků (tyhle se sami neaktualizují):

```bash
brew install --cask claude-code          # macOS Homebrew
winget install Anthropic.ClaudeCode      # Windows
npm install -g @anthropic-ai/claude-code # jakýkoli OS s Node 22+ — nikdy s sudo
```

Pak **otevři nový terminál** (aby se obnovil `PATH`) a ověř:

```bash
claude --version     # čekej 2.1.257 (Claude Code) nebo novější
claude doctor        # read-only diagnostika; oprav, co nahlásí
```

`command not found`? Na macOS/Linuxu žije binárka v `~/.local/bin` — ověř,
že je tahle složka na tvém `PATH` (`echo $PATH`), pak restartuj terminál.

## 2. Jen Windows: Git for Windows (Git Bash)

Nainstaluj [Git for Windows](https://git-scm.com/downloads/win) s výchozími
volbami. Dá to Claude Code Bash shell (**Bash tool**) — bez něj Claude spadne
zpátky na PowerShell a hooks tohohle repa (bash skripty) nepoběží. Pokud ho
Claude Code nenajde, ukaž na něj v `~/.claude/settings.json`:

```json
{ "env": { "CLAUDE_CODE_GIT_BASH_PATH": "C:\\Program Files\\Git\\bin\\bash.exe" } }
```

Alternativa: dělej všechno uvnitř **WSL 2** (Ubuntu) a tam se řiď
macOS/Linux instrukcemi. WSL Git for Windows nepotřebuje.

## 3. Předpoklady podle OS

**macOS**

```bash
xcode-select --install          # Git (přeskoč, pokud `git --version` už funguje)
brew install node go jq         # Node.js 22+, Go, jq
```

**Windows (PowerShell)**

```powershell
winget install OpenJS.NodeJS.LTS   # Node.js 22+
winget install GoLang.Go           # Go
winget install jqlang.jq           # volitelné
```

(Nebo použij instalátory z [nodejs.org](https://nodejs.org) a
[go.dev/dl](https://go.dev/dl/).)

**Ubuntu / Debian**

```bash
sudo apt update && sudo apt install -y git jq
# Node.js 22 — distro balíček bývá moc starý; použij NodeSource:
curl -fsSL https://deb.nodesource.com/setup_22.x | sudo -E bash - && sudo apt install -y nodejs
# Go — stáhni aktuální tarball z https://go.dev/dl/ a rozbal ho:
sudo rm -rf /usr/local/go && sudo tar -C /usr/local -xzf go1.25.*.linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.profile && source ~/.profile
```

Zkontroluj každý:

```bash
git --version      # jakákoli
node --version     # v22 nebo novější
npx --version      # přijde s npm
go version         # go1.25 nebo novější
jq --version        # volitelné
```

## 4. První běh v tomhle repu

```bash
git clone https://github.com/klesajos/qti-ws
cd qti-ws
go test ./...      # stáhne golang.org/x/term, musí vypsat "ok"
claude
```

Při prvním spuštění `claude`:

1. **Přihlas se** — otevře se okno prohlížeče; přihlas se účtem, který má
   přístup ke Claude Code. (Pokud používáš API klíč, exportuj
   `ANTHROPIC_API_KEY` a jednou ho na dotaz potvrď.)
2. **Důvěřuj workspace** — potvrď dialog. Tohle povolí projektovou složku
   `.claude/`, `.mcp.json` a in-repo marketplace pluginů. Pokud chceš, přečti
   si soubory nejdřív — přesně na to jsou návody tady.
3. **Zkontroluj, že se rozšíření načetla:**
   - `/mcp` → `deepwiki` i `memory` jsou **connected** (dokazuje to, že `npx`
     funguje; první start stáhne balíček a může trvat minutu)
   - napiš `/2048-dev:` → doplní se `build-test` (dokazuje, že se plugin
     načetl)
   - `/hooks` → vypsané jsou hooky `PostToolUse` a `Stop`

## 5. Aktualizace a odinstalace

```bash
claude update                        # nativní instalace (taky se aktualizuje sama)
brew upgrade claude-code             # Homebrew
winget upgrade Anthropic.ClaudeCode  # WinGet
```

Odinstalace (nativní): `rm -f ~/.local/bin/claude && rm -rf ~/.local/share/claude`
na macOS/Linuxu, nebo smaž `%USERPROFILE%\.local\bin\claude.exe` a
`%USERPROFILE%\.local\share\claude` na Windows.

## Chceš GUI? Desktop aplikace

[Desktop aplikace Claude](https://code.claude.com/docs/en/desktop-quickstart)
balí Claude Code do své záložky **Code** — bez terminálu, stejný engine,
stejná projektová konfigurace. Všechno v tomhle repu (hooks, skills, MCP,
plugin, agenti, workflows) funguje v záložce Code. **Ne**funguje to
v záložce **Cowork**, která běží ve vlastním sandboxu a projektovou
konfiguraci ignoruje — viz tabulka platforem v [indexu](README.cs.md).

## Pre-flight checklist

Spusť tyhle příkazy ve složce repa. Všechno zelené → jsi připravený.

```bash
claude --version   # ≥ 2.1.257
git --version
node --version     # ≥ v22
go version         # ≥ go1.25
go test ./...      # ok
```

## Řešení problémů

| Příznak | Oprava |
|---|---|
| `claude: command not found` hned po instalaci | Otevři nový terminál; na macOS/Linuxu ověř, že `~/.local/bin` je na `PATH` |
| `/mcp` ukazuje `memory` jako failed | `node --version` musí být ≥ 22; první běh `npx` potřebuje internet a minutu. Firemní proxy? Nastav `HTTPS_PROXY` |
| Hooks vypisují `bash: command not found` (Windows) | Nainstaluj Git for Windows (sekce 2) nebo nastav `CLAUDE_CODE_GIT_BASH_PATH` |
| `go test` nedokáže stáhnout moduly | Zkontroluj proxy; `go env GOPROXY` by mělo být `https://proxy.golang.org,direct` |
| Přihlašovací smyčka / 403 | `claude auth status`, pak `claude auth logout && claude auth login`; ověř, že účet má placený plán |
| Něco jiného | `claude doctor`, pak [dokumentace k řešení problémů](https://code.claude.com/docs/en/troubleshoot-install) |
