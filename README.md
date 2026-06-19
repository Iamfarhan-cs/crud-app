# Task Management CRUD API

A production-minded backend CRUD API for managing tasks, built in Go with PostgreSQL.

This project is intentionally structured like a real backend service: HTTP concerns live in handlers, business rules live in services, persistence is hidden behind a repository boundary, and PostgreSQL connection setup is isolated from feature code.

## Features

- Task create, read, update, list, and soft delete
- PostgreSQL-backed persistence through `database/sql`
- Layered handler, service, repository, and database boundaries
- Stable JSON request and response shapes
- Service-level validation for titles, statuses, and pagination
- Parameterized SQL queries
- Request body size limit
- Database connection pool configuration
- Health and readiness endpoints
- Graceful server shutdown
- Docker Compose setup for API and PostgreSQL
- Focused unit tests for service and handler behavior

## Project Overview

The API manages task records through a simple lifecycle:

- `pending`
- `in_progress`
- `done`

Each task has a required title, optional description, status, timestamps, and soft-delete support through `deleted_at`.

The main engineering goal is not just to make CRUD work. The goal is to show clean backend thinking:

- clear package boundaries
- predictable request and response shapes
- service-level validation
- repository abstraction
- parameterized SQL
- graceful server shutdown
- focused unit tests
- documentation that explains both the product and the design

## Tech Stack

- Go
- PostgreSQL
- `database/sql`
- `github.com/lib/pq` PostgreSQL driver
- `github.com/google/uuid`
- Standard library HTTP server
- Standard library testing package

## Architecture Overview

The project follows a layered backend structure:

```text
Client
  -> HTTP Handler
  -> Service Layer
  -> Repository Interface
  -> PostgreSQL Repository
  -> Database
```

```text
                 +---------------------+
Client / curl -> | net/http ServeMux   |
                 +----------+----------+
                            |
                            v
                 +---------------------+
                 | Task HTTP Handlers  |
                 +----------+----------+
                            |
                            v
                 +---------------------+
                 | Task Service Layer  |
                 +----------+----------+
                            |
                            v
                 +---------------------+
                 | Repository Contract |
                 +----------+----------+
                            |
                            v
                 +---------------------+
                 | PostgreSQL Repo     |
                 +----------+----------+
                            |
                            v
                 +---------------------+
                 | PostgreSQL Database |
                 +---------------------+
```

Each layer has one job.

- Handlers parse HTTP input and write HTTP responses.
- Services enforce business rules and coordinate use cases.
- Repositories handle persistence details.
- The database package opens and verifies PostgreSQL connections.
- The config package loads runtime settings from environment variables.
- `cmd/api` composes the application and starts the server.

This separation keeps business logic out of HTTP code and SQL out of service code. That is the foundation that makes the project easier to test, change, and reason about.

## Folder Structure

```text
cmd/
  api/
    main.go                    # application entry point and HTTP server wiring
internal/
  config/
    config.go                  # environment configuration loading
  database/
    postgres.go                # PostgreSQL connection setup
  task/
    dto.go                     # request and response DTOs
    errors.go                  # domain/service sentinel errors
    handler.go                 # HTTP handlers and route registration
    handler_test.go            # handler unit tests
    model.go                   # task domain model
    postgres_repository.go     # PostgreSQL repository implementation
    repository.go              # repository interface
    service.go                 # task business logic
    service_test.go            # service unit tests
migrations/
  001_create_tasks_table.up.sql
  001_create_tasks_table.down.sql
docker-compose.yml
Dockerfile
Makefile
go.mod
go.sum
README.md
PHASES.md
```

## Setup Instructions

1. Clone the repository.

```text
git clone https://github.com/Iamfarhan-cs/crud-app.git
cd crud-app
```

2. Install Go dependencies.

```text
go mod tidy
```

3. Create a PostgreSQL database.

4. Set the required environment variables.

5. Apply the migration in `migrations/001_create_tasks_table.up.sql`.

6. Run the API.

```text
go run ./cmd/api
```

Or use Make:

```text
make run
```

The server starts on `PORT` or defaults to `8080`.

## Environment Variables

### Required

| Name | Required | Default | Purpose |
| ---- | -------- | ------- | ------- |
| `DATABASE_URL` | Yes | none | PostgreSQL connection string |

### Optional

| Name | Required | Default | Purpose |
| ---- | -------- | ------- | ------- |
| `PORT` | No | `8080` | HTTP server port |
| `READ_TIMEOUT` | No | `5s` | Maximum request read duration |
| `WRITE_TIMEOUT` | No | `10s` | Maximum response write duration |
| `IDLE_TIMEOUT` | No | `60s` | Keep-alive idle timeout |
| `SHUTDOWN_TIMEOUT` | No | `10s` | Graceful shutdown timeout |
| `DB_MAX_OPEN_CONNECTIONS` | No | `10` | Maximum open PostgreSQL connections |
| `DB_MAX_IDLE_CONNECTIONS` | No | `5` | Maximum idle PostgreSQL connections |
| `DB_CONNECTION_MAX_LIFE` | No | `30m` | Maximum lifetime for a database connection |

Example:

```text
DATABASE_URL=postgres://user:password@localhost:5432/tasks?sslmode=disable
PORT=8080
```

## Database And Migrations

Migration files live in `migrations/`.

Current migration:

- `001_create_tasks_table.up.sql` - creates the `tasks` table and active-task index
- `001_create_tasks_table.down.sql` - drops the index and table

The task table includes:

- app-generated UUID primary key
- non-empty title constraint
- status constraint for `pending`, `in_progress`, and `done`
- nullable `description`
- `created_at` and `updated_at`
- nullable `deleted_at` for soft deletes
- partial index for active task listing by newest records first

This project does not include a migration runner yet. Apply migrations with your preferred PostgreSQL tool, for example `psql`.

```text
psql "$DATABASE_URL" -f migrations/001_create_tasks_table.up.sql
```

To roll back:

```text
psql "$DATABASE_URL" -f migrations/001_create_tasks_table.down.sql
```

## Docker Run Instructions

Build and start the API plus PostgreSQL:

```text
make docker-up
```

Equivalent command:

```text
docker compose up --build
```

Stop the containers:

```text
make docker-down
```

Reset the Docker database volume and start fresh:

```text
make docker-reset
```

Equivalent commands:

```text
docker compose down -v
docker compose up --build
```

## API Endpoints Summary

| Method | Path | Purpose |
| ------ | ---- | ------- |
| `GET` | `/healthz` | Health check |
| `GET` | `/readyz` | Database readiness check |
| `POST` | `/api/v1/tasks` | Create a task |
| `GET` | `/api/v1/tasks` | List active tasks |
| `GET` | `/api/v1/tasks/{id}` | Get one active task |
| `PATCH` | `/api/v1/tasks/{id}` | Partially update an active task |
| `DELETE` | `/api/v1/tasks/{id}` | Soft-delete a task |

List requests support optional pagination query parameters:

```text
GET /api/v1/tasks?page=1&limit=20
```

## Health And Readiness Checks

Health check:

```text
curl -i http://localhost:8080/healthz
```

Readiness check:

```text
curl -i http://localhost:8080/readyz
```

## Curl Examples

Create a task:

```text
curl -i -X POST http://localhost:8080/api/v1/tasks \
  -H "Content-Type: application/json" \
  -d "{\"title\":\"Ship deployment notes\",\"description\":\"Prepare Phase K README\",\"status\":\"pending\"}"
```

List tasks:

```text
curl -i http://localhost:8080/api/v1/tasks
```

Paginated list tasks:

```text
curl -i "http://localhost:8080/api/v1/tasks?page=1&limit=20"
```

Update a task:

```text
curl -i -X PATCH http://localhost:8080/api/v1/tasks/{id} \
  -H "Content-Type: application/json" \
  -d "{\"status\":\"in_progress\"}"
```

Delete a task:

```text
curl -i -X DELETE http://localhost:8080/api/v1/tasks/{id}
```

## Standard Error Response Format

Error responses use a stable JSON shape:

```json
{
  "error": {
    "code": "INVALID_TITLE",
    "message": "invalid title"
  }
}
```

## Testing Instructions

Run the full test suite:

```text
go test ./...
```

Or use Make:

```text
make test
```

The current tests focus on:

- service validation and business behavior
- handler status codes and error mapping
- soft-delete visibility
- pagination validation
- unknown error handling

On this Windows environment, if Go fails before compiling because the default build cache cannot be initialized, rerun with a repo-local cache:

```text
$env:GOCACHE = (Join-Path (Get-Location) ".gocache")
go test ./...
```

## Make Targets

| Target | Command |
| ------ | ------- |
| `make run` | `go run ./cmd/api` |
| `make test` | `go test ./...` |
| `make fmt` | `go fmt ./...` |
| `make vet` | `go vet ./...` |
| `make docker-up` | `docker compose up --build` |
| `make docker-down` | `docker compose down` |
| `make docker-reset` | `docker compose down -v` then `docker compose up --build` |

## Deployment Checklist

- `DATABASE_URL` configured
- migrations applied
- healthz works
- readyz works
- request IDs in logs
- graceful shutdown enabled
- request body limit enabled
- DB pool configured
- tests pass
- `.env` is not committed
- production secrets not in Git

## Production Notes

Before production deployment, confirm runtime configuration, database access, migration strategy, logging, and operational checks in the target environment.

Future improvements:

- Authentication
- Authorization
- CI/CD
- External migration runner
- Structured JSON logging
- Prometheus metrics
- OpenTelemetry tracing
- Rate limiting
- Caching
- Scaling strategy

## Project Learning Tracker

This project was built phase by phase to understand backend engineering from first principles.

For the detailed phase-by-phase breakdown, see:

[PHASES.md](./PHASES.md)
