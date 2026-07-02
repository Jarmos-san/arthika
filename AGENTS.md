# Arthika — AGENTS.md

## Quick start

### Server

```bash
task server:test            # go test ./... (t.Parallel(), black-box)
task server:lint            # golangci-lint v2, all linters
task migrate:up             # apply pending DB migrations (requires `migrate` CLI)
task migrate:down           # roll back all down migrations
task migrate:reset          # down then up
task oapi:gen               # generate server Go stubs + client TypeScript types from OpenAPI spec
go run ./server/cmd/server  # start server on :8000
```

### Client

```bash
cd client && pnpm dev        # Nuxt 4 dev server
cd client && pnpm typecheck  # vue-tsc --noEmit (requires .nuxt/ — run `pnpm install` first)
cd client && pnpm lint       # oxlint
cd client && pnpm fmt        # oxfmt
pnpm install --frozen-lockfile  # install deps (used in CI and setup)

pre-commit install           # install git hooks (includes crisp for commit messages)
```

### Single test

```bash
go test ./server/internal/handler/ -run TestRegister_Success
```

## Architecture

- Monorepo: `server/` (Go) + `client/` (Nuxt 4 / Vue 3 / TypeScript)
- Entrypoint: `server/cmd/server/main.go`
- Layered: **Handler → Repository** — no service package yet (deliberate MVP choice)
- Packages: `api/`, `app/`, `auth/`, `config/`, `handler/`, `logger/`, `middleware/`, `repository/` under `server/internal/`
- API spec in `server/api/openapi.yml`, JSON:API format (`application/vnd.api+json`)
- Routing: `go-chi/chi/v5`
- Logging: `slog` JSON handler to stdout
- Config from env vars: `ADDR`, `READ_TIMEOUT`, `WRITE_TIMEOUT`, `IDLE_TIMEOUT`, `LOG_LEVEL`, `TOKEN_SECRET`, `DATABASE_URL`, `MIGRATIONS_DIR`
- Database: SQLite via `golang-migrate` (auto-applied at startup)
- Docs: served at `/docs` (Redoc UI) and `/openapi.json`

## Code generation

Three generators, run in this order when the OpenAPI spec changes:

```bash
cd server && sqlc generate                            # reads db/query/*.sql → internal/repository/
cd server && oapi-codegen -config oapi-codegen.yml api/openapi.yml  # → internal/api/server.gen.go
cd client && pnpm run openapi:generate                # kubb reads server api/openapi.yml → client/generated/
```

Shortcut: `task oapi:gen` runs server + client OpenAPI generation.

**Note:** `task server:generate` exists but has a bug — it references a non-existent `oapi-codegen` subtask. Use `task oapi:gen` instead for OpenAPI codegen.

Run codegen after changing `db/query/*.sql` (sqlc), `api/openapi.yml` (oapi-codegen or kubb), then commit generated files.

## Client

- Nuxt 4 with `ssr: false` (SPA mode), `@nuxt/ui`, `vue-router`
- Package manager: `pnpm` (lockfile `pnpm-lock.yaml`)
- Linting: `oxlint` (`pnpm lint`)
- Formatting: `oxfmt` (`pnpm fmt`, `pnpm fmt:check`)
- TypeScript OpenAPI client: `kubb` (`pnpm run openapi:generate` → `client/generated/`)
- Typecheck requires `.nuxt/` to exist — run `pnpm install` (or `task client:install`) first

## CI

Three workflows on PR to `main`:
- **Server QA** — lint (`golangci-lint`), test (`go test ./...`), build check (`go build -o /dev/null ./...`)
- **Client QA** — lint (`oxlint`), format check (`oxfmt --check`), typecheck (`vue-tsc`), build (`nuxt build`)
- **SQL Validation** — `sqlc vet` + `sqlc diff` (ensures generated Go code is in sync)

## Testing

- All tests use `t.Parallel()` and black-box packages (`_test` suffix)
- Handler tests share a `mockQuerier` struct implementing `repository.Querier`
- Pre-compute bcrypt hashes (`bcrypt.MinCost`) as package-level constants (see `login_test.go` `testPasswordHash`) to avoid slow hashing per test
- Use `t.Context()` and `slog.Default()` (or `slog.DiscardHandler`)
- Integration-style handler tests use `httptest.NewRecorder` + `api.NewStrictHandler`
- Run a single test: `go test ./server/internal/handler/ -run TestRegister_Success`

## Conventions

- `depguard` restricts `internal/*.go` to stdlib imports only
- Go: tab indentation; JSON/YAML/TS/Vue: 2-space (see `.editorconfig`)
- Formatters: `gci`, `gofmt`, `gofumpt`, `goimports`, `golines` (via golangci-lint)
- Migrations: `golang-migrate` CLI, naming `NNNNNN_description.{up,down}.sql`
- Commit messages linted by `crisp` pre-commit hook
- Pre-commit hooks: `trailing-whitespace`, `end-of-file-fixer`, `check-yaml`, `check-added-large-files`, `oxfmt-check` (client), `oxlint` (client), `crisp`
