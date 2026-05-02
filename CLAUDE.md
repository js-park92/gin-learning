# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

```bash
make run              # start the server (loads .env automatically)
make build            # compile to bin/server
make test             # go test ./... -v
make lint             # golangci-lint run ./...
make fmt              # gofmt + goimports in-place
make tidy             # go mod tidy

make migrate-up                        # apply all pending migrations
make migrate-down                      # roll back 1 migration
make migrate-down steps=3              # roll back 3 migrations
make migrate-version                   # print current version + dirty flag
make migrate-create name=add_foo_bar   # scaffold new up/down SQL pair

make docker-up        # start postgres in background
make docker-down      # stop containers
```

Go is managed by **asdf** (`golang` plugin). The binary lives at `~/.asdf/installs/golang/<version>/go/bin/go`.

## Architecture

The app follows a strict four-layer stack. Each layer only imports the one directly below it — never skips.

```
handler  →  service  →  repository  →  database (GORM)
```

- **`internal/model/`** — plain Go structs with GORM tags. No methods, no DB awareness.
- **`internal/repository/`** — all GORM calls live here, behind an interface. Services never call GORM directly.
- **`internal/service/`** — business logic only. Receives and returns model types.
- **`internal/handler/`** — HTTP request binding, response serialisation, calls one service. No DB or business logic.
- **`internal/router/`** — wires handlers to routes. Uses `gin.New()` (not `gin.Default()`) so middleware is explicit.
- **`internal/middleware/`** — Gin middleware. The access log middleware sets log level by status: `info` 2xx/3xx, `warn` 4xx, `error` 5xx.
- **`internal/logger/`** — initialises the global zerolog logger once. Debug mode → colourised stderr. Release mode → JSON stdout.
- **`internal/config/`** — all env vars flow through `config.Load()` into a typed `Config` struct. No `os.Getenv` calls outside this package.
- **`cmd/migrate/`** — standalone binary that runs golang-migrate. Not part of the server process.

Dependency injection is manual and done entirely in `main.go` — no DI framework. The wiring order is: config → logger → database → repository → service → handler → router.

## Migrations

Migration files live in `migrations/` as plain SQL pairs (`000001_foo.up.sql` / `000001_foo.down.sql`). The server **does not auto-migrate on startup**. Run `make migrate-up` explicitly before starting the server against a fresh database.

GORM's soft-delete (`deleted_at`) is active on all models that embed `gorm.Model`. The unique index on `users.email` filters `WHERE deleted_at IS NULL` so soft-deleted emails can be reused.

Never edit a migration that has already been applied — add a new one instead.

## Adding a new resource

1. Add a model struct in `internal/model/` embedding `gorm.Model`.
2. Add a repository interface + implementation in `internal/repository/`.
3. Add a service interface + implementation in `internal/service/`.
4. Add handlers in `internal/handler/`.
5. Register routes in `internal/router/router.go`.
6. Wire the new types in `main.go`.
7. Create a migration: `make migrate-create name=create_<resource>`.

## Logging

Use the global zerolog logger everywhere — no logger instances are passed around:

```go
import "github.com/rs/zerolog/log"

log.Info().Str("key", "value").Msg("something happened")
log.Error().Err(err).Msg("something failed")
```
