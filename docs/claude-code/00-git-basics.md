> 🌍 Read this in: **English** | [Česky](00-git-basics.cs.md)

# Before the workshop: Git in ten minutes

You don't need to be a Git expert for the workshop, but Claude Code works
*with* Git constantly — it reads `git diff`, the reviewer agent only sees
what Git sees, and committing is how you keep the good states. This page is
the minimum.

## Setup (once)

```bash
git --version                                   # installed? (Windows: Git for Windows)
git config --global user.name  "Your Name"
git config --global user.email "you@quanti.cz"
git config --global init.defaultBranch main
```

## Get the workshop repo

```bash
git clone https://github.com/klesajos/qti-ws
cd qti-ws
```

## The daily loop

```bash
git status                      # what changed? (do this often)
git diff                        # show the changes line by line
git add -A                      # stage everything (or: git add path/to/file)
git commit -m "Add undo"        # save a snapshot with a message
git log --oneline               # history, one line per commit
git push                        # send commits to GitHub (if you have a remote)
```

Rule for the workshop: **commit after every exercise that works.** A commit
is a save point you can always return to; Claude's `/rewind` is not a
substitute (it only tracks edits made through Claude Code).

## Branches — try things without breaking `main`

```bash
git switch -c feature/undo      # create a branch and switch to it
# ... work, commit ...
git switch main                 # go back
git merge feature/undo          # bring the branch in (or open a PR on GitHub)
git branch                      # list branches, * marks the current one
```

## Undo, three flavours

```bash
git restore path/to/file.go     # throw away uncommitted changes to one file
git restore .                   # ... to everything (careful)
git reset --soft HEAD~1         # undo the last commit, keep the changes staged
git revert <commit>             # make a new commit that undoes an old one (safe on shared branches)
```

## What Claude Code needs from you

- **A clean tree before big tasks.** `git status` shows nothing? Good — now
  every change Claude makes is exactly `git diff`, and the `go-reviewer`
  agent reviews only that.
- **Commit at milestones**, not at the end of the day. Small commits make
  `git diff` and `/rewind` useful.
- **Never `git push --force` on a shared branch.** If you must rewrite
  history, do it on your own branch.
- **Don't commit secrets.** `.gitignore` already excludes local Claude
  settings; API keys live in environment variables (see
  [Example 3](03-mcp.md)).

## Cheat table

| I want to… | Command |
|---|---|
| See what changed | `git status` · `git diff` |
| Save a snapshot | `git add -A && git commit -m "message"` |
| See history | `git log --oneline --graph` |
| Start a branch | `git switch -c name` |
| Discard my edits to a file | `git restore file` |
| Undo the last commit (keep changes) | `git reset --soft HEAD~1` |
| Get the latest from GitHub | `git pull` |

## Check

```bash
cd qti-ws && git status && git log --oneline | head -3
```

You should see `nothing to commit, working tree clean` and the repo's last
three commits.
