# gin-learning

A Go REST API built with [Gin](https://gin-gonic.com/), [GORM](https://gorm.io/), and PostgreSQL.

## Prerequisites

- [Go 1.22+](https://go.dev/dl/)
- [Docker](https://www.docker.com/) and Docker Compose
- [golangci-lint](https://golangci-lint.run/welcome/install/) (for linting)

## Getting started

### 1. Environment

```bash
cp .env.example .env
# edit .env as needed
```

### 2. Start PostgreSQL

```bash
make docker-up
```

### 3. Apply migrations

```bash
make migrate-up
```

### 4. Run the server

```bash
make run
```

The server starts on `http://localhost:8080`.

---

## Migrations

Migrations are plain SQL files managed by [golang-migrate](https://github.com/golang-migrate/migrate).
The runner lives in `cmd/migrate/main.go` and is invoked via `go run` — no extra CLI install required.

### File naming

```
migrations/
  000001_create_users.up.sql
  000001_create_users.down.sql
  000002_add_users_phone.up.sql
  000002_add_users_phone.down.sql
```

Each migration has an `up` file (apply) and a `down` file (rollback). The six-digit prefix controls execution order.

### Commands

| Command | Description |
|---|---|
| `make migrate-up` | Apply all pending migrations |
| `make migrate-down` | Roll back 1 migration |
| `make migrate-down steps=3` | Roll back 3 migrations |
| `make migrate-version` | Print the current migration version |
| `make migrate-create name=add_users_phone` | Scaffold a new migration pair |

### Creating a new migration

```bash
make migrate-create name=add_users_phone
# → migrations/000002_add_users_phone.up.sql
# → migrations/000002_add_users_phone.down.sql
```

Then write the SQL:

```sql
-- 000002_add_users_phone.up.sql
ALTER TABLE users ADD COLUMN phone VARCHAR(20);
```

```sql
-- 000002_add_users_phone.down.sql
ALTER TABLE users DROP COLUMN phone;
```

Apply it:

```bash
make migrate-up
```

### Rules

- **Never edit a migration that has already been applied** — add a new one instead.
- Every `up` migration must have a corresponding `down` that fully reverses it.
- Migrations run outside the application process. The server does **not** auto-migrate on startup.
- In CI/CD, run `make migrate-up` before deploying the new binary.

### How golang-migrate tracks state

golang-migrate creates a `schema_migrations` table in your database with two columns:

| Column | Description |
|---|---|
| `version` | The sequence number of the last applied migration |
| `dirty` | `true` if the last migration failed mid-run |

If a migration fails and the table is marked dirty, fix the SQL, then manually clear the dirty flag before retrying:

```sql
UPDATE schema_migrations SET dirty = false;
```

---

## Logging

Logging uses [zerolog](https://github.com/rs/zerolog) — a structured, zero-allocation JSON logger.

### Output format

| `GIN_MODE` | Format | Output |
|---|---|---|
| `debug` (default) | Colorized console | `stderr` |
| `release` | JSON | `stdout` |

**Debug output** (human-readable, for local development):
```
10:32:01 INF database connected db=gin_learning host=localhost
10:32:01 INF starting server port=8080
10:32:05 INF request ip=127.0.0.1 latency=1.2ms method=GET path=/api/v1/users status=200
10:32:07 WRN request ip=127.0.0.1 latency=0.8ms method=GET path=/api/v1/users/99 status=404
10:32:09 ERR request ip=127.0.0.1 latency=3.1ms method=POST path=/api/v1/users status=500
```

**Release output** (JSON, for log aggregators like Datadog, Loki, CloudWatch):
```json
{"level":"info","host":"localhost","db":"gin_learning","time":"2024-01-01T10:32:01Z","message":"database connected"}
{"level":"info","method":"GET","path":"/api/v1/users","status":200,"latency":"1.2ms","ip":"127.0.0.1","time":"2024-01-01T10:32:05Z","message":"request"}
```

The access log middleware automatically uses `warn` for 4xx responses and `error` for 5xx, so you can alert on log level alone in production.

### Using the logger in your own code

Import the global logger from zerolog and add structured fields with `.Str()`, `.Int()`, `.Err()`, etc.:

```go
import "github.com/rs/zerolog/log"

// plain message
log.Info().Msg("doing something")

// with structured fields
log.Info().Str("user_id", id).Str("email", email).Msg("user created")

// with an error
log.Error().Err(err).Str("user_id", id).Msg("failed to update user")

// fatal — logs then calls os.Exit(1)
log.Fatal().Err(err).Msg("unrecoverable error")
```

---

## Available commands

| Command | Description |
|---|---|
| `make run` | Run the server locally |
| `make build` | Build the binary to `bin/server` |
| `make test` | Run all tests |
| `make lint` | Run golangci-lint |
| `make fmt` | Format code with gofmt + goimports |
| `make tidy` | Tidy go.mod and go.sum |
| `make migrate-up` | Apply all pending migrations |
| `make migrate-down` | Roll back 1 migration |
| `make migrate-version` | Show current migration version |
| `make migrate-create name=<n>` | Scaffold a new migration pair |
| `make docker-up` | Start services in the background |
| `make docker-down` | Stop and remove containers |
| `make docker-logs` | Tail container logs |

## Project structure

```
.
├── cmd/
│   └── migrate/       # Migration runner (golang-migrate)
├── internal/
│   ├── config/        # Typed config loaded from env
│   ├── database/      # GORM connection
│   ├── handler/       # HTTP handlers
│   ├── logger/        # zerolog initialisation (format, level)
│   ├── middleware/    # Gin middleware (access log, etc.)
│   ├── model/         # GORM models (plain structs)
│   ├── repository/    # DB access layer (interfaces + GORM impl)
│   ├── router/        # Route registration
│   └── service/       # Business logic
├── migrations/        # Versioned SQL migration files
├── main.go
├── Dockerfile
├── docker-compose.yml
├── .golangci.yml
└── Makefile
```

## Linting

Install golangci-lint:

```bash
brew install golangci-lint
```

Then run:

```bash
make lint
```

## Environment variables

| Variable | Default | Description |
|---|---|---|
| `DB_HOST` | `localhost` | PostgreSQL host |
| `DB_PORT` | `5432` | PostgreSQL port |
| `DB_USER` | `gin_user` | Database user |
| `DB_PASSWORD` | `gin_password` | Database password |
| `DB_NAME` | `gin_learning` | Database name |
| `GIN_MODE` | `debug` | `debug` or `release` |
| `PORT` | `8080` | HTTP server port |
