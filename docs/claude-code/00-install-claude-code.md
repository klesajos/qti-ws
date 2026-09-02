> 🌍 Read this in: **English** | [Česky](00-install-claude-code.cs.md)

# Before the workshop: install Claude Code and the prerequisites

Do this **before the workshop day** — it takes 15 minutes with a good
connection, and the last section is a checklist that tells you you're done.
Verified against Claude Code **2.1.257** and the official
[setup page](https://code.claude.com/docs/en/setup) on 2026-09-03.

## What you need

| | Why | Version |
|---|---|---|
| **A Claude account** with Claude Code access | Pro, Max, Team, Enterprise, or a Console (API) account. The free plan does **not** include Claude Code | — |
| **Claude Code** | The agent itself. A native binary — it does *not* need Node.js to run | ≥ 2.1.257 |
| **Git** | Cloning this repo, committing your work; on Windows it also provides the Bash that Claude Code and this repo's hooks use | any recent |
| **Node.js + npm/npx** | Only because this repo's `.mcp.json` starts the `memory` MCP server with `npx` ([Example 3](03-mcp.md)) | **22+** |
| **Go** | Building and testing the app — see [Running a Go app](00-go-basics.md) | **1.25+** |
| `jq` *(optional)* | Nicer JSON handling in the hooks; they fall back to `python3` or `sed` without it | any |

Supported OS: macOS 13+, Windows 10 1809+ / Server 2019+, Ubuntu 20.04+ /
Debian 10+. 4 GB RAM. Internet access (the first run also downloads Go
modules and the npm package).

## 1. Install Claude Code

Pick **one** method. The native installer is recommended — it updates itself
in the background.

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

> Not sure which Windows shell you're in? The prompt shows `PS C:\Users\you>`
> in PowerShell and `C:\Users\you>` (no `PS`) in CMD. If you see
> `The token '&&' is not a valid statement separator`, you're in PowerShell —
> use the PowerShell command. If you see `'irm' is not recognized`, you're in
> CMD.

Alternatives, if you prefer a package manager (these don't auto-update):

```bash
brew install --cask claude-code          # macOS Homebrew
winget install Anthropic.ClaudeCode      # Windows
npm install -g @anthropic-ai/claude-code # any OS with Node 22+ — never with sudo
```

Then **open a new terminal** (so `PATH` is refreshed) and check:

```bash
claude --version     # expect 2.1.257 (Claude Code) or newer
claude doctor        # read-only diagnostics; fix anything it flags
```

`command not found`? On macOS/Linux the binary lives in `~/.local/bin` — make
sure that directory is on your `PATH` (`echo $PATH`), then restart the terminal.

## 2. Windows only: Git for Windows (Git Bash)

Install [Git for Windows](https://git-scm.com/downloads/win) with the default
options. It gives Claude Code a Bash shell (the **Bash tool**) — without it
Claude falls back to PowerShell, and this repo's hooks (bash scripts) won't
run. If Claude Code can't find it, point to it in `~/.claude/settings.json`:

```json
{ "env": { "CLAUDE_CODE_GIT_BASH_PATH": "C:\\Program Files\\Git\\bin\\bash.exe" } }
```

Alternative: run everything inside **WSL 2** (Ubuntu) and follow the
macOS/Linux instructions there. WSL doesn't need Git for Windows.

## 3. Prerequisites per OS

**macOS**

```bash
xcode-select --install          # Git (skip if `git --version` already works)
brew install node go jq         # Node.js 22+, Go, jq
```

**Windows (PowerShell)**

```powershell
winget install OpenJS.NodeJS.LTS   # Node.js 22+
winget install GoLang.Go           # Go
winget install jqlang.jq           # optional
```

(Or use the installers from [nodejs.org](https://nodejs.org) and
[go.dev/dl](https://go.dev/dl/).)

**Ubuntu / Debian**

```bash
sudo apt update && sudo apt install -y git jq
# Node.js 22 — the distro package is usually too old; use NodeSource:
curl -fsSL https://deb.nodesource.com/setup_22.x | sudo -E bash - && sudo apt install -y nodejs
# Go — download the current tarball from https://go.dev/dl/ and unpack it:
sudo rm -rf /usr/local/go && sudo tar -C /usr/local -xzf go1.25.*.linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.profile && source ~/.profile
```

Check each one:

```bash
git --version      # any
node --version     # v22 or newer
npx --version      # comes with npm
go version         # go1.25 or newer
jq --version       # optional
```

## 4. First run in this repo

```bash
git clone https://github.com/klesajos/qti-ws
cd qti-ws
go test ./...      # downloads golang.org/x/term, must print "ok"
claude
```

On the first `claude` run:

1. **Log in** — a browser window opens; sign in with the account that has
   Claude Code access. (If you use an API key, export `ANTHROPIC_API_KEY`
   and approve it once when asked.)
2. **Trust the workspace** — accept the dialog. This is what enables the
   project's `.claude/` folder, `.mcp.json` and the in-repo plugin
   marketplace. Read the files first if you like; that's what the guides
   are for.
3. **Check the extensions loaded:**
   - `/mcp` → `deepwiki` and `memory` both **connected** (this proves `npx`
     works; the first start downloads the package and can take a minute)
   - type `/2048-dev:` → `build-test` autocompletes (this proves the plugin
     loaded)
   - `/hooks` → the `PostToolUse` and `Stop` hooks are listed

## 5. Updating and removing

```bash
claude update                        # native install (also auto-updates)
brew upgrade claude-code             # Homebrew
winget upgrade Anthropic.ClaudeCode  # WinGet
```

Uninstall (native): `rm -f ~/.local/bin/claude && rm -rf ~/.local/share/claude`
on macOS/Linux, or remove `%USERPROFILE%\.local\bin\claude.exe` and
`%USERPROFILE%\.local\share\claude` on Windows.

## Prefer a GUI? The Desktop app

The [Claude Desktop app](https://code.claude.com/docs/en/desktop-quickstart)
bundles Claude Code in its **Code** tab — no terminal needed, same engine,
same project config. Everything in this repo (hooks, skills, MCP, plugin,
agents, workflows) works in the Code tab. It does **not** work in the
**Cowork** tab, which runs in its own sandbox and ignores project-scoped
config — see the platform table in the [index](README.md).

## Pre-flight checklist

Run these in the repo folder. All green → you're ready.

```bash
claude --version   # ≥ 2.1.257
git --version
node --version     # ≥ v22
go version         # ≥ go1.25
go test ./...      # ok
```

## Troubleshooting

| Symptom | Fix |
|---|---|
| `claude: command not found` right after installing | Open a new terminal; on macOS/Linux ensure `~/.local/bin` is on `PATH` |
| `/mcp` shows `memory` as failed | `node --version` must be ≥ 22; first `npx` run needs internet and a minute. Corporate proxy? Set `HTTPS_PROXY` |
| Hooks print `bash: command not found` (Windows) | Install Git for Windows (section 2) or set `CLAUDE_CODE_GIT_BASH_PATH` |
| `go test` fails to download modules | Check the proxy; `go env GOPROXY` should be `https://proxy.golang.org,direct` |
| Login loops / 403 | `claude auth status`, then `claude auth logout && claude auth login`; confirm the account has a paid plan |
| Something else | `claude doctor`, then the [troubleshooting docs](https://code.claude.com/docs/en/troubleshoot-install) |
