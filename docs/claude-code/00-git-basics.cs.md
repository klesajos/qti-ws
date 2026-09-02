> 🌍 Číst v jazyce: [English](00-git-basics.md) | **Česky**

# Před workshopem: Git za deset minut

Na workshop nemusíš být expert na Git, ale Claude Code s Gitem pracuje
*neustále* — čte `git diff`, reviewer agent vidí jen to, co vidí Git, a
commitování je způsob, jak si držet dobré stavy. Tahle stránka je minimum.

## Nastavení (jednorázově)

```bash
git --version                                   # nainstalovaný? (Windows: Git for Windows)
git config --global user.name  "Tvoje jméno"
git config --global user.email "ty@quanti.cz"
git config --global init.defaultBranch main
```

## Získání workshopového repa

```bash
git clone https://github.com/klesajos/qti-ws
cd qti-ws
```

## Každodenní smyčka

```bash
git status                      # co se změnilo? (dělej to často)
git diff                        # ukáže změny řádek po řádku
git add -A                      # stage všechno (nebo: git add path/to/file)
git commit -m "Add undo"        # ulož snapshot se zprávou
git log --oneline               # historie, jeden řádek na commit
git push                        # pošle commity na GitHub (pokud máš remote)
```

Pravidlo pro workshop: **commitni po každém cvičení, které funguje.** Commit
je záchytný bod, ke kterému se vždycky můžeš vrátit; Claudův `/rewind`
náhrada není (sleduje jen úpravy udělané přes Claude Code).

## Větve — zkoušej věci bez rozbití `main`

```bash
git switch -c feature/undo      # vytvoří větev a přepne se na ni
# ... práce, commit ...
git switch main                 # zpátky
git merge feature/undo          # zapoj větev (nebo otevři PR na GitHubu)
git branch                      # seznam větví, * označuje aktuální
```

## Undo, tři varianty

```bash
git restore path/to/file.go     # zahodí necommitnuté změny v jednom souboru
git restore .                   # ... ve všem (opatrně)
git reset --soft HEAD~1         # zruší poslední commit, změny zůstanou staged
git revert <commit>             # vytvoří nový commit, který ruší starý (bezpečné na sdílených větvích)
```

## Co Claude Code od tebe potřebuje

- **Čistý strom před velkými úkoly.** `git status` neukazuje nic? Dobře —
  teď je každá Claudova změna přesně `git diff` a agent `go-reviewer`
  posuzuje jen tohle.
- **Commituj na milnících**, ne až na konci dne. Malé commity dělají
  `git diff` a `/rewind` užitečnými.
- **Nikdy `git push --force` na sdílenou větev.** Pokud musíš přepsat
  historii, dělej to na vlastní větvi.
- **Necommituj tajemství.** `.gitignore` už vylučuje lokální nastavení
  Claude; API klíče žijí v proměnných prostředí (viz [Ukázka 3](03-mcp.cs.md)).

## Tahák

| Chci… | Příkaz |
|---|---|
| Vidět, co se změnilo | `git status` · `git diff` |
| Uložit snapshot | `git add -A && git commit -m "message"` |
| Vidět historii | `git log --oneline --graph` |
| Založit větev | `git switch -c name` |
| Zahodit úpravy souboru | `git restore file` |
| Zrušit poslední commit (zachovat změny) | `git reset --soft HEAD~1` |
| Stáhnout nejnovější z GitHubu | `git pull` |

## Kontrola

```bash
cd qti-ws && git status && git log --oneline | head -3
```

Měl bys vidět `nothing to commit, working tree clean` a poslední tři
commity repa.
