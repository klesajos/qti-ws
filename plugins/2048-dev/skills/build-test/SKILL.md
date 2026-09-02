---
name: build-test
description: Build the project, run go vet, check formatting and run the full test suite, then summarize results
allowed-tools: Bash(go:*), Bash(gofmt:*)
disable-model-invocation: true
---

# Build and test

Build the project and run the checks, in this order:

1. Build: `go build ./...`
2. Vet: `go vet ./...`
3. Format check: `gofmt -l .` (every file it prints is unformatted)
4. Test: `go test ./...`

Then report:
- Build: OK, or the first compiler error with file:line.
- Vet / format: OK, or each reported file:line.
- Tests: ok/FAIL per package. For each failure, name the `Test...` function
  and show its message with actual vs expected values.
- Do not fix anything — this skill only reports status.
