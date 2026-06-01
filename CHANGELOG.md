# Changelog

All notable changes to this project are documented in this file.

## [0.1.0] - 2026-06-01

### Added

- User registration with password hashing (bcrypt)
- JWT-based login authentication
- JSON:API response format (`application/vnd.api+json`)
- Chi router with logger and recoverer middleware
- Structured logging via `slog` JSON handler to stdout
- Environment-based configuration (`ADDR`, `READ_TIMEOUT`, `WRITE_TIMEOUT`, `IDLE_TIMEOUT`, `LOG_LEVEL`, `TOKEN_SECRET`)
- Graceful server shutdown on SIGINT/SIGKILL
- Comprehensive test suite with `t.Parallel()` and black-box test packages
- golangci-lint v2 with all linters enabled
- Pre-commit hooks with `crisp` for commit message linting

### Changed

- Server directory structure: `handler/`, `service/`, `domain/`, `app/`, `dto/`, `logger/`
- Router migrated from `net/http` to `go-chi/chi/v5`
- `ARCHITECTURE.md` moved from `docs/` to project root
- Project name changed to Arthika

### Removed

- Removed `server/tutorial/` sqlc scratch directory
- Removed `docs/` directory (contents moved to root)
