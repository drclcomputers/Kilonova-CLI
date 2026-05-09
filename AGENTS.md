# Repository Guidelines

## Project Structure

This repository is a Go CLI for the Kilonova competitive programming platform. `main.go` only starts the Cobra command tree. Command groups live under `cmd/`: `contests`, `problems`, `project`, `submission`, `user`, `database`, and `ai`. Shared API, database, encryption, Gemini, constants, and utility code belongs in `internal/`. Keep command-specific flags and presentation close to the relevant `cmd/<area>` package; put reusable behavior in `internal/`.

Generated binaries such as `kilonova-cli` and `kilonova-cli.exe` should not be edited. Prefer source changes and rebuild locally.

## Build and Run

- `go mod tidy` updates module metadata after dependency changes.
- `go build ./...` verifies all packages compile.
- `go build -o kilonova-cli.exe .` builds the Windows CLI binary.
- `go run . help` lists available commands.
- `go run . search`, `go run . signin`, and `go run . submission` are useful smoke checks for CLI wiring.

The module currently targets Go `1.24.1`; use that version or newer unless the project explicitly changes its baseline.
This project uses `github.com/mattn/go-sqlite3`, so Windows builds need CGO and a C compiler. On this machine, MSYS2 UCRT64 GCC is configured through `go env` and `C:\msys64\ucrt64\bin` is on the user PATH.

## Coding Style

Run `gofmt` on every edited Go file before committing. Use package aliases only where they make command registration clearer, as in `cmd/root.go`. Keep Cobra command names short and user-facing help text accurate. Return errors where possible instead of printing deep inside helper functions, unless the command package is deliberately responsible for terminal output.

## Testing

There are no committed `_test.go` files at the moment. Add focused tests next to new logic when changing parsing, API response handling, database behavior, or command construction. Prefer table-driven Go tests for pure functions in `internal/`. For interactive commands, keep behavior split so non-UI logic can be tested without a terminal session.

Before handing off changes, run at least:

```powershell
go test ./...
go build ./...
```

## Commits and Pull Requests

Use concise, imperative commit messages such as `Add problem search filters` or `Fix submission status parsing`. PRs should summarize the user-visible change, list verification commands, and mention any Kilonova API, local database, or credential-storage impact.

## Security

Do not commit credentials, session tokens, local database files, or generated logs containing account data. Treat authentication, encryption, and Gemini API changes as sensitive and document any migration or configuration steps clearly.
