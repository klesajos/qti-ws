#!/usr/bin/env bash
# PostToolUse hook: gofmt any Go file Claude just edited or wrote.
# Receives the tool call as JSON on stdin; non-Go files are ignored.
set -euo pipefail

input=$(cat)

# Extract the edited file's path from the JSON. jq is preferred; python3 and a
# plain sed pattern are fallbacks for machines without it.
if command -v jq >/dev/null 2>&1; then
    file_path=$(printf '%s' "$input" | jq -r '.tool_input.file_path // empty')
elif command -v python3 >/dev/null 2>&1; then
    file_path=$(printf '%s' "$input" | python3 -c \
        'import json,sys; print(json.load(sys.stdin).get("tool_input", {}).get("file_path", ""))')
else
    file_path=$(printf '%s' "$input" | sed -n 's/.*"file_path"[[:space:]]*:[[:space:]]*"\([^"]*\)".*/\1/p')
fi

# gofmt ships with every Go installation. goimports (optional, from
# golang.org/x/tools) does the same and also fixes the import block.
if command -v goimports >/dev/null 2>&1; then
    formatter="goimports -w"
elif command -v gofmt >/dev/null 2>&1; then
    formatter="gofmt -w"
else
    exit 0
fi

case "$file_path" in
    *.go)
        if [ -f "$file_path" ]; then
            $formatter "$file_path"
            echo "format-go hook: formatted ${file_path##*/}"
        fi
        ;;
esac

exit 0
