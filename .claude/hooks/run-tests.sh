#!/usr/bin/env bash
# Stop hook: when Claude finishes a turn, run the test suite once and print a
# one-line pass/fail summary. Advisory and non-blocking by design — it never
# stops Claude, it just surfaces whether the tree is still green.
set -euo pipefail

# Drain the Stop-event JSON on stdin. We don't need any of its fields, but a
# hook must read stdin so Claude Code isn't left writing to a closed pipe.
cat >/dev/null

project_dir="${CLAUDE_PROJECT_DIR:-$(pwd)}"

# Stay silent unless this is a Go module and the toolchain is on PATH, so the
# hook never nags in the wrong directory or on a machine without Go.
if [ ! -f "$project_dir/go.mod" ] || ! command -v go >/dev/null 2>&1; then
    exit 0
fi

# Run the suite once, keeping the output and the exit status. `set -e` is
# suspended around the call so a red suite doesn't abort the hook itself.
set +e
output=$(cd "$project_dir" && go test ./... 2>&1)
status=$?
set -e

# Count the per-package result lines: "ok <pkg>" for green packages,
# "FAIL <pkg>" for red ones (a bare "FAIL" line closes the run).
passed=$(printf '%s\n' "$output" | grep -c '^ok' || true)
failing=$(printf '%s\n' "$output" | awk '/^FAIL[[:space:]]/ {print $2}' | paste -sd ',' - || true)

if [ "$status" -eq 0 ]; then
    echo "run-tests hook: ✓ go test ./... — ${passed} package(s) ok"
else
    echo "run-tests hook: ✗ go test ./... failed in ${failing:-the build} (run 'go test ./...' for details)"
fi

exit 0
