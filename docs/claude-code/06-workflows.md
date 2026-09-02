> 🌍 Read this in: **English** | [Česky](06-workflows.cs.md)

# Example 6: Workflows (multi-agent orchestration)

## What is a workflow?

A **workflow** is a small JavaScript program that **orchestrates many agents
with deterministic control flow**. Instead of asking Claude to "figure out the
steps", you write the steps — fan out work with `parallel()`, chain it with
`pipeline()`, group it with `phase()` — and the runtime executes that structure
exactly, every time. A saved workflow becomes a slash command you can run with
`/<name>`.

Where a single agent decides its own path, a workflow's path is fixed in code:
loops, conditionals, and fan-out are yours to control.

> Dynamic workflows are available on all paid plans. On **Pro**, turn them on
> once in `/config` → *Dynamic workflows*.

## Subagent vs. workflow

These are the two delegation mechanisms, and it's worth being clear on the
difference before reading the code:

| | Subagent (Example 5) | Workflow (this example) |
|---|---|---|
| **Unit** | one delegated context | many agents, orchestrated |
| **Control flow** | the agent decides its own steps | *you* decide, in JavaScript (`parallel`/`pipeline`/`phase`) |
| **Determinism** | model-driven, varies run to run | the structure is fixed; same input → same shape of run |
| **Defined in** | a Markdown file + system prompt | a `.js` script with `export const meta` |
| **Invoke** | plain language, `@agent-name`, `--agent` | a slash command `/<name>` (opt-in) |

**Rule of thumb:** reach for a **subagent** when you want to hand *one*
self-contained task to *one* isolated context. Reach for a **workflow** when you
want to run *many* agents in a *repeatable, code-defined* pipeline — fan out over
N files, gate stage 2 on stage 1, reduce many results into one report.

## What this example does

This example ships one runnable workflow, `test-coverage-audit`, that audits how
well `internal/board/board_test.go` covers the logic in the source files. It is
**read-only, deterministic, and re-runnable** — it prints a prioritized gap
report and never writes a file, so you can run it as often as you like and get a
stable result.

It works in three phases:

1. **Inventory** — fan out one read-only reader per source file; each enumerates
   the behaviours a unit test could target.
2. **Map tests** — one agent maps every existing `Test...` function onto that
   inventory.
3. **Report gaps** — one reducer cross-references the two and emits a prioritized
   gap report.

Because the real gaps in this repo are known (the `{4, 4, 8, 0}` TODO and the
cascade-merge bug behind it, the `SpawnRandom` distribution, the
full-board-with-a-mergeable-pair game-over case that trips the `IsGameOver()`
bug, multi-move sequences, the discarded `Move()` result in `game.go`), the
report is predictable — which makes it a good teaching artifact.

## The file, line by line

The workflow lives at `.claude/workflows/test-coverage-audit.js`. Saved
workflows in `.claude/workflows/` are auto-discovered and become slash commands.

**1. The `meta` block must be the very first statement — this is the one
everyone trips on.** A stray `import`, `const`, or any executable line above
`export const meta` makes the workflow **silently fail to register**:
`/test-coverage-audit` never appears, with no error to explain why. A file-header
`//` comment before `meta` is fine (this file opens with one); executable code is
not. The block names the workflow (this becomes `/test-coverage-audit`) and
describes its phases:

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

**2. The five source files it inventories.** A `SRC_FILES` array names every
Go file outside the tests, each with a one-line note so the reader agent
orients quickly. This is the literal "five files" the Phase-1 fan-out runs over:

```js
const SRC_FILES = [
  { file: 'internal/board/board.go', note: 'pure rules: SlideLineLeft, Move, SpawnRandom, IsGameOver, HasWon' },
  { file: 'internal/game/game.go', note: 'main loop: toDirection, Game.Run, spawn + redraw, win/over checks' },
  { file: 'internal/input/input.go', note: 'raw-mode terminal: ReadByte -> Command (WASD + arrow keys, q)' },
  { file: 'internal/renderer/renderer.go', note: 'draws the grid + score to the terminal' },
  { file: 'cmd/2048/main.go', note: 'entry point: game.New, Run, Close' },
]
```

**3. Three schemas force each agent to return data, not prose.** Passing a
`schema` to `agent()` makes the runtime validate the agent's output, so the
script gets a real object back — no parsing. The workflow defines one per phase:
`BEHAVIOR_SCHEMA` (what one file does, branch by branch), `COVERAGE_SCHEMA`
(which behaviour ids each test function exercises), and `GAP_SCHEMA` (the
prioritized gaps plus a rendered table). Only `BEHAVIOR_SCHEMA` is shown here:

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

`GAP_SCHEMA` carries a `markdown` field — and **that field is what the slash
command ultimately returns** (see step 5).

**4. Phase 1 fans out with `parallel()`, then Phase 2 maps the tests — and every
agent is read-only.** Each agent runs as the built-in **`Explore`** type
(`agentType: 'Explore'`); `phase()` just labels each group in the progress UI:

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

`parallel()` is a **barrier**: it waits for *all* readers before continuing.
That's deliberate here — Phase 2 needs the *whole* inventory before it can map
tests onto it. (Contrast `pipeline()`, which runs each *item* through all stages
independently; it suits per-item chains like the feature-pipeline below, not a
fan-in.) `.filter(Boolean)` drops any agent that was stopped or failed — an
`agent()` call resolves to `null` in that case.

> **Why is every agent `Explore`?** `agentType: 'Explore'` agents can read and
> search but have no `Edit`, `Write`, or `Bash`. That enforcement — not a
> request — is what makes `/test-coverage-audit` safe to run on a loop:
> `git status` stays clean every time. Add write power only together with
> isolation (`isolation: 'worktree'`).

**5. The reduce prompt is what makes the run deterministic.** Phase 3 hands one
reducer the whole inventory and coverage map, and — crucially — *pins the
findings it must surface* right in the prompt. Here is that Phase-3 prompt
(only the injected JSON is elided):

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
return report.markdown   // the GAP_SCHEMA.markdown field is what the command surfaces
```

Pinning the expected gaps in the prompt is **why** re-running
`/test-coverage-audit` surfaces the same findings every time: the model isn't
rediscovering them from scratch, it's confirming a known list (and is told to
flag the two `board` bugs as *latent bugs*, not hide them). The script's
`return report.markdown` is exactly that `GAP_SCHEMA.markdown` field — what the
slash command prints.

Three rules the runtime enforces on the script body, so you don't discover them
the hard way: no `import()` (put library work inside an agent's task), no
`Date.now()` / `Math.random()` / `new Date()` (they throw, so a relaunched run
repeats the same agent calls — pass a timestamp through `args` instead), and no
direct filesystem or shell access (agents do that; the script coordinates).

## A more advanced shape: the feature-pipeline (described, not shipped)

The audit only reads. A workflow can also *do* work safely if you isolate it.
A `feature-pipeline` for adding a feature would look like:

```
design (read-only architect)
  → enter a git worktree (isolation)
    → implement the feature
    → go build + go test ./..., with a bounded retry on failure
    → review with the go-reviewer agent from Example 5   ← agents compose into workflows
  → exit the worktree
→ hand the branch to a human to merge
```

Two ideas worth noting: a **worktree** keeps the work on an isolated copy of the
repo so a failed run never touches your tree, and the review stage **reuses the
`go-reviewer` subagent** — Example 5's agent becomes a stage in Example 6's
pipeline. It's worth this much machinery when a change is multi-file, test-gated,
and you want concurrency; it's overkill for a one-line fix you'd just make
directly. ([Exercise 6](exercises.md#6-workflow) builds a generate-and-select
variant of this shape.)

## Create your own workflow, step by step

1. **Create the folder:**
   ```bash
   mkdir -p .claude/workflows
   ```
2. **Create** `.claude/workflows/my-flow.js` starting with `meta`:
   ```js
   export const meta = { name: 'my-flow', description: 'What it does.', phases: [{ title: 'Work' }] }
   phase('Work')
   const result = await agent('Do one well-scoped thing.', { schema: { type: 'object', properties: { summary: { type: 'string' } }, required: ['summary'] } })
   return result.summary
   ```
   Before writing anything bigger, run `/workflow-authoring` — a bundled skill
   that loads the script-writing reference Claude itself works from.
3. **Keep it deterministic and, ideally, read-only first.** Use
   `agentType: 'Explore'` for agents that should only read. Add write power only
   when you also add isolation (a worktree).
4. **Restart Claude Code** (or run `/reload-skills`) so the workflow is
   discovered, then run `/my-flow`. Running a workflow is **opt-in** — you
   invoke it explicitly and approve the phase list once.
5. **Commit it:**
   ```bash
   git add .claude/workflows/my-flow.js && git commit -m "Add my-flow workflow"
   ```

Alternatively, let Claude write the first draft: describe the task and say
"use a workflow" (or include the keyword `ultracode`); when a run does what you
wanted, open `/workflows`, select it and press `s` to save it as a command.

## Try the demo

1. In a fresh `claude` session in the repo root, run **`/test-coverage-audit`**
   and approve the run. Watch it fan out over the five source files (Phase 1),
   map the tests (Phase 2), and print a prioritized gap report (Phase 3) that
   names the known gaps — including the `IsGameOver()` and cascade-merge bugs as
   latent issues. Open `/workflows` while it runs to see each phase's agents,
   token totals and elapsed time.
2. Check `git status` — it's **clean**. The workflow only read.
3. Run it again. The *shape* is the same and the same gaps surface — that's the
   determinism a code-defined pipeline buys you.
4. Read the "feature-pipeline" section above and trace where a worktree would
   isolate the work and where `go-reviewer` plugs in as the review stage.

## Optional parameters

Per-`agent()` options you'll reach for most:

| Option | What it does |
|--------|--------------|
| `schema` | JSON Schema the agent's output must match; `agent()` returns the validated object |
| `label` | The name shown for this agent in the progress UI |
| `phase` | Assigns the agent to a progress group (use inside `parallel`/`pipeline`) |
| `agentType` | Run as a specific agent type, e.g. `Explore` for enforced read-only |
| `model` / `effort` | Override model or reasoning effort for a stage |
| `isolation: 'worktree'` | Run the agent in a fresh git worktree — for stages that mutate files |

Orchestration primitives: `phase()`, `parallel()` (barrier), `pipeline()`
(per-item chains, no barrier), `log()`, and the `args` global for input passed
at invocation. Runtime limits: at most 16 agents run concurrently, 1,000 per
run; the advisory *Dynamic workflow size* setting in `/config` (default
`medium` = aim for fewer than 15 agents) steers how large Claude makes the
workflows it writes. Full reference:
[official documentation](https://code.claude.com/docs/en/workflows).

## Where it works: CLI, Desktop app, Cowork

| Platform | Works? | Setup |
|----------|--------|-------|
| **Claude Code CLI** (terminal) | ✅ Yes | Workflows in `.claude/workflows/` are discovered at session start; run with `/<name>`, watch with `/workflows` |
| **Claude Desktop app — Code tab** | ✅ Yes | Same engine; the official docs state dynamic workflows "run in Desktop". An approval card shows the phases; progress appears in the Background tasks pane |
| **Cowork** (in the Desktop app) | ❌ No | Cowork does not load project-scoped `.claude/` config, and has no workflow runtime of its own |

Turned off for you or your organisation? Check `/config` → *Dynamic
workflows*, the `disableWorkflows` setting, or `CLAUDE_CODE_DISABLE_WORKFLOWS`.

## Troubleshooting

- **`/test-coverage-audit` doesn't appear** → `export const meta` is not the
  first statement, or contains something that isn't a plain literal; also check
  workflows aren't disabled (`/config`).
- **The run asks for permissions mid-way** → its agents use your normal
  permission rules; add the tools they need to your allow rules before a long
  run (our audit needs none — it only reads).
- **Edited the script but the old version runs** → run `/reload-skills`, then
  `/<name>` again.
- **A phase shows an agent as `null` / missing** → it was stopped or hit an API
  error; the script's `.filter(Boolean)` skips it. Relaunch — completed agents
  return their cached results, only the failed one and those after it rerun.
