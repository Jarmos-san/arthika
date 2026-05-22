# AGENTS.md — Arthika

## Project structure

```
server/          # Go module (single binary, no frontend yet)
  cmd/server/main.go     # Entrypoint
  internal/
    application/         # App container, server lifecycle
    config/              # Env-based config (ADDR, *_TIMEOUT, LOG_LEVEL, TOKEN_SECRET)
    dto/                 # JSON:API document/resource types
    handlers/            # HTTP handlers, thin transport layer
    logger/              # slog.NewJSONHandler wrapper
    models/              # Domain entities
    services/            # Business logic
```

Architecture doc at `docs/ARCHITECTURE.md` describes a planned Nuxt.js frontend — **no frontend code exists**.

## Developer commands

| Command | Where | What |
|---|---|---|
| `task lint` | repo root | `golangci-lint run ./...` in `server/` |
| `go test ./...` | `server/` | all tests |
| `go run ./cmd/server` | `server/` | start dev server on `:8000` |

There is no Makefile — `task` runner (Taskfile.yml) is used.

## Framework & toolchain

- **Go 1.26.2** — routes use Go 1.22+ method syntax (`"GET /users/"`, not legacy `"/users/"`)
- **`net/http`** with `http.NewServeMux()` — no third-party router
- **`slog`** with `slog.NewJSONHandler` — JSON-structured logging to stdout
- **golangci-lint v2** — config at `server/.golangci.yml` with `version: "2"`, all default linters + `wsl_v5`
- Formatters: `gofumpt`, `gci`, `goimports`, `golines`
- **pre-commit**: trailing-whitespace, end-of-file-fixer, check-yaml, check-added-large-files, `crisp` (commit message lint)
- **EditorConfig**: tabs in `.go`, 2-space indent in `.json`/`.yml`

## Config

All via environment variables (defaults in parens):
- `ADDR` (`:8000`), `READ_TIMEOUT` / `WRITE_TIMEOUT` / `IDLE_TIMEOUT` (`10s`)
- `LOG_LEVEL` (`"info"`)
- `TOKEN_SECRET` (`"super-secret-token"` — **not production-safe**)

## API conventions

- Response Content-Type: `application/vnd.api+json` (JSON:API spec)
- Handlers receive injected logger + service, never initialize deps internally

## Known gotchas

- `Application.Shutdown()` at `internal/application/application.go:87` calls `ListenAndServe()` instead of `Shutdown()` — this is **buggy**.
- `createJWT()` in `internal/handlers/utils.go:51` calls `config.LoadConfig()` on every invocation instead of using an injected config.
- All tests use `t.Parallel()` and black-box package naming (`package handlers_test`).
- Pre-commit hooks use `crisp` for commit message linting — messages must pass its format rules.
