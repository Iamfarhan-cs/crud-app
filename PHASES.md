# Project Phase Tracker

This file is the learning tracker for the Task Management CRUD API. The README is intentionally kept clean for project presentation; this file explains how the backend was built phase by phase and why each step mattered.

# Phase 1: Requirements And Mental Model

## Goal

Define the product behavior before writing backend code. A serious backend starts with clear rules: what the entity is, what states are valid, which fields are server-owned, and which failure cases the API must handle.

## What Changed

The project documented the task entity, expected CRUD operations, business rules, non-goals, production risks, and failure scenarios. This created a shared mental model before the code grew.

## Files Changed

* `README.md` - captured the first version of requirements, entity fields, API operations, and production risks.

## Feature Flow

Idea -> task entity definition -> business rules -> expected CRUD operations -> failure scenarios -> later code phases.

## Engineering Reasoning

The structure was chosen because requirements drive architecture. The tradeoff was moving slower at the start so later code would be less random. This phase solved the problem of building endpoints before understanding the data lifecycle.

## Production Notes

In production, unclear requirements become inconsistent validation, breaking API responses, and data cleanup problems. The areas that need monitoring later are invalid input rates, missing record errors, storage failures, and slow list operations.

## Common Mistakes Avoided

Junior engineers often start with routes and structs first. This phase avoided that by defining behavior, failure cases, and non-goals before implementation.

## Commit Link

Commit: https://github.com/Iamfarhan-cs/crud-app/commit/1f342ffe1e8ac570a748c0065ac865aede55b7fa

# Phase 2: Project Architecture And Boundaries

## Goal

Create a clean Go project structure where each package has one responsibility.

## What Changed

The project added the initial package layout under `cmd/api`, `internal/config`, `internal/database`, and `internal/task`. The architecture separated application startup, configuration, database setup, and task feature code.

## Files Changed

* `cmd/api/main.go` - application entry point.
* `internal/config/config.go` - configuration boundary.
* `internal/database/postgres.go` - database setup boundary.
* `internal/task/model.go` - task domain model location.
* `internal/task/handler.go` - HTTP layer location.
* `internal/task/service.go` - business logic location.
* `internal/task/repository.go` - persistence boundary location.
* `README.md` - documented architecture and responsibilities.

## Feature Flow

Application startup -> load config -> prepare database boundary -> compose task package -> future HTTP request flow.

## Engineering Reasoning

The structure was used to prevent business rules, SQL, and HTTP behavior from being mixed together. The tradeoff was more files early, but that cost buys clarity and testability. This prepared the project for service, repository, handler, and wiring phases.

## Production Notes

In production, mixed layers make bugs harder to isolate. Config, database connections, and HTTP handlers need separate logging and failure handling because they fail for different reasons.

## Common Mistakes Avoided

This avoided putting all code in `main.go`, letting handlers talk directly to SQL, and letting packages depend on each other in both directions.

## Commit Link

Commit: https://github.com/Iamfarhan-cs/crud-app/commit/e428266e823e82e457c9cc1b110cc0e91a990e0a

# Phase 3: API Design

## Goal

Design the external HTTP contract before implementing full endpoint logic.

## What Changed

The project documented REST endpoints, request bodies, response shapes, error responses, pagination expectations, and server-owned fields. DTO structs were added so handlers and services had a contract to build against.

## Files Changed

* `internal/task/dto.go` - request, response, list, and error DTOs.
* `README.md` - endpoint contract, examples, status codes, and pagination rules.

## Feature Flow

Client request -> request DTO -> service input -> task response DTO -> standard JSON response.

## Engineering Reasoning

The API contract was designed before endpoint logic so clients and backend code had a stable target. The tradeoff was documenting behavior before it existed, but that made implementation more deliberate.

## Production Notes

API contracts become public promises. Field names, status codes, and error codes need stability because clients may depend on them.

## Common Mistakes Avoided

This avoided inventing response shapes inside handlers, exposing internal models directly, and returning inconsistent error bodies across endpoints.

## Commit Link

Commit: https://github.com/Iamfarhan-cs/crud-app/commit/e58725ae30e7079ec38f29db14a45e1c1f9c91ac

# Phase 4: Database Design

## Goal

Define the PostgreSQL shape for tasks and protect important data rules at the database boundary.

## What Changed

The project added task table migrations with UUID primary keys, non-empty title validation, status validation, timestamps, soft-delete support, and an active-task partial index.

## Files Changed

* `migrations/001_create_tasks_table.up.sql` - creates the task table and index.
* `migrations/001_create_tasks_table.down.sql` - drops the index and table.
* `README.md` - documented the database design.

## Feature Flow

Service-created task -> repository insert -> `tasks` table -> active list query uses partial index -> response DTO.

## Engineering Reasoning

The database was designed to reject impossible task states even if application code has a bug. The tradeoff is that migrations must be maintained carefully, but the database becomes a final line of defense.

## Production Notes

Bad constraints can block valid writes, and missing indexes can slow list endpoints. Soft deletes require every normal read and update to filter `deleted_at IS NULL`.

## Common Mistakes Avoided

This avoided storing blank titles, accepting arbitrary statuses, hard deleting records too early, and indexing deleted rows for the common active-list path.

## Commit Link

Commit: https://github.com/Iamfarhan-cs/crud-app/commit/f9fd8f4f3a8a0534f080a6082b399f6734c0df48

# Phase 5: Storage Design

## Goal

Plan the persistence approach before writing repository implementation details.

## What Changed

The project documented storage responsibilities, repository expectations, query safety, soft-delete behavior, and separation between business logic and persistence.

## Files Changed

* `README.md` - documented storage design and repository responsibilities.

## Feature Flow

Service layer -> repository interface -> storage implementation -> database rows -> domain task.

## Engineering Reasoning

Storage was treated as a boundary, not just a place to paste SQL. This prepared the project to swap or test persistence more easily.

## Production Notes

Persistence code needs parameterized queries, context-aware calls, clear error mapping, and operational logging around slow or failed queries.

## Common Mistakes Avoided

This avoided letting SQL leak into services, building queries with string concatenation, and treating deleted records as normal active records.

## Commit Link

Commit: https://github.com/Iamfarhan-cs/crud-app/commit/341c41d7668c063dd4b0b0459cd44ca69020e06d

# Phase 6: Networking Design

## Goal

Plan how HTTP requests enter the system and how the server should behave under normal and failure conditions.

## What Changed

The project documented route responsibilities, handler boundaries, request parsing, response behavior, health checks, timeouts, and graceful shutdown expectations.

## Files Changed

* `README.md` - documented networking and HTTP server design.

## Feature Flow

HTTP client -> server mux -> task handler -> service -> repository -> response writer.

## Engineering Reasoning

Networking design came before full handler implementation so the HTTP layer stayed focused on transport concerns. This prepared the project for predictable routing and error mapping.

## Production Notes

Production HTTP servers need timeouts, bounded request bodies, health checks, and graceful shutdown. Without those, slow clients and deploys can cause resource leaks.

## Common Mistakes Avoided

This avoided unbounded request bodies, handlers with business logic, missing health checks, and servers that stop abruptly during deploys.

## Commit Link

Commit: https://github.com/Iamfarhan-cs/crud-app/commit/6c0811a231025332d15aab549b3b42a1a754f7f0

# Phase 7: Security Design

## Goal

Document the security posture before adding authentication or authorization.

## What Changed

The project documented input validation, parameterized SQL, safe error responses, secret handling, CORS expectations, request body limits, and future ownership rules.

## Files Changed

* `README.md` - documented security design.
* `internal/task/handler.go` - comments clarified handler security boundaries.
* `internal/task/service.go` - comments clarified service validation responsibilities.
* `internal/task/repository.go` - comments clarified SQL safety expectations.
* `internal/config/config.go` - comments clarified environment-based configuration.

## Feature Flow

Client input -> handler parsing and body limit -> service validation -> repository parameterized SQL -> safe client error response.

## Engineering Reasoning

Security was handled as a system concern. The tradeoff was documenting rules before implementing auth, but it made the future security model clearer.

## Production Notes

Secrets must not be committed, client errors should not leak internals, and database queries must remain parameterized. Authentication and authorization are still future production requirements.

## Common Mistakes Avoided

This avoided hardcoding secrets, returning raw database errors to clients, and treating validation as only an HTTP concern.

## Commit Link

Commit: https://github.com/Iamfarhan-cs/crud-app/commit/ba01f12091327a66c8b7afbbe0e23041816550ed

# Phase 8: Concurrency Design

## Goal

Plan how the backend should behave when multiple requests touch the same data or when requests are cancelled.

## What Changed

The project documented context propagation, cancellation, race conditions, stale updates, idempotency, and future concurrency controls.

## Files Changed

* `README.md` - documented concurrency design and failure cases.

## Feature Flow

HTTP request context -> service call -> repository call -> database operation -> cancellation or completion.

## Engineering Reasoning

Concurrency was treated as a production concern, not an afterthought. This phase prepared repository and service code to accept `context.Context` and to reason about future conflict handling.

## Production Notes

Cancelled requests should not keep doing expensive database work. Concurrent updates can overwrite each other unless future versioning or locking is added.

## Common Mistakes Avoided

This avoided ignoring request contexts, assuming one user updates at a time, and treating delete/update races as impossible.

## Commit Link

Commit: https://github.com/Iamfarhan-cs/crud-app/commit/3e79ac8d759994d027ff08ee64b79d3565fdd6d4

# Phase 9: Performance Design

## Goal

Document the performance expectations before optimizing prematurely.

## What Changed

The project documented pagination limits, active-task indexing, connection pool settings, query timeouts, slow query logging, and future cursor pagination.

## Files Changed

* `README.md` - documented performance design.
* `internal/task/repository.go` - comments clarified context-aware repository calls.

## Feature Flow

List request -> pagination normalization -> repository query with limit and offset -> active-task index -> bounded response.

## Engineering Reasoning

The design favored simple offset pagination for the first version, with clear limits. The tradeoff is that offset pagination can degrade at scale, but it is easy to understand and good enough for an early CRUD backend.

## Production Notes

Monitor latency percentiles, slow queries, connection pool saturation, and list endpoint usage. At larger scale, cursor pagination and richer indexes may be needed.

## Common Mistakes Avoided

This avoided unbounded list responses, `SELECT *` thinking, missing indexes, and adding caching before measuring actual bottlenecks.

## Commit Link

Commit: https://github.com/Iamfarhan-cs/crud-app/commit/d4a2bb96a993b801de763d5ba2cd4123b29fe769

# Phase A: Domain Model

## Goal

Create the internal task model used by the application.

## What Changed

The project added the `Task` domain model and status constants for the supported task lifecycle.

## Files Changed

* `internal/task/model.go` - task domain model and `Status` values.
* `README.md` - documented the build phase.

## Feature Flow

Repository row or service input -> `Task` domain model -> service rules -> response DTO.

## Engineering Reasoning

The model is internal because it represents the system's view of a task, not a direct client payload. This keeps API DTOs and domain data separate.

## Production Notes

Domain models need stable invariants. If the model grows, keep server-owned fields protected from direct client writes.

## Common Mistakes Avoided

This avoided using raw strings everywhere and avoided exposing one struct as database row, domain object, request body, and response body all at once.

## Commit Link

Commit: https://github.com/Iamfarhan-cs/crud-app/commit/9d2d4b52504f2b97bdfb2caf7caf63d9e52883bd

# Phase B: Errors And Repository Contract

## Goal

Define stable service errors and the repository interface before implementing persistence.

## What Changed

The project added sentinel errors and a `Repository` interface for create, read, list, update, and soft-delete operations.

## Files Changed

* `internal/task/errors.go` - stable task errors.
* `internal/task/repository.go` - persistence interface.
* `README.md` - documented the build phase.

## Feature Flow

Service method -> repository interface -> concrete repository -> sentinel error or domain task -> handler error mapping.

## Engineering Reasoning

The service depends on an interface so business logic is not coupled directly to PostgreSQL. Sentinel errors give handlers stable meanings to map into HTTP responses.

## Production Notes

Wrapped errors must preserve `errors.Is` compatibility. Repository implementations should translate storage-specific missing-row cases into `ErrTaskNotFound`.

## Common Mistakes Avoided

This avoided letting PostgreSQL errors leak into handlers and avoided making services impossible to unit test without a database.

## Commit Link

Commit: https://github.com/Iamfarhan-cs/crud-app/commit/dc662268311d81f456ce95ab88869d6e45c97ae4

# Phase C: Service Layer Implementation

## Goal

Implement task business rules independently from HTTP and PostgreSQL.

## What Changed

The service now creates tasks, gets tasks, lists tasks, updates tasks, deletes tasks, validates titles, validates statuses, normalizes pagination, generates IDs, and maps domain tasks to response DTOs.

## Files Changed

* `internal/task/service.go` - business logic and use-case orchestration.
* `README.md` - documented service behavior.

## Feature Flow

Handler DTO -> service validation -> repository call -> domain task -> response DTO.

## Engineering Reasoning

Business rules belong in the service so they apply consistently no matter which client or transport calls the use case. The tradeoff is that handlers become thinner and services carry more responsibility, which is the right direction for testability.

## Production Notes

Validation must stay consistent across create and update paths. Pagination limits protect the service from expensive list calls.

## Common Mistakes Avoided

This avoided validating only in handlers, accepting blank titles, accepting invalid statuses, and allowing empty update requests.

## Commit Link

Commit: https://github.com/Iamfarhan-cs/crud-app/commit/4bb78684f83c42ab890a0bfcf90fafbc587a7e08

# Phase D: PostgreSQL Repository Implementation

## Goal

Connect the repository interface to real PostgreSQL queries.

## What Changed

The project added a PostgreSQL repository that inserts tasks, reads active tasks, lists active tasks, updates active tasks, and soft-deletes tasks with parameterized SQL.

## Files Changed

* `internal/task/postgres_repository.go` - concrete PostgreSQL repository.
* `README.md` - documented repository behavior.

## Feature Flow

Service -> repository interface -> PostgreSQL repository -> parameterized SQL -> scanned `Task` -> service response.

## Engineering Reasoning

The repository hides SQL from the service. Parameterized queries were chosen for safety, and soft delete checks were added to normal reads and updates.

## Production Notes

Database errors need context for operators but should still preserve original causes. Queries should honor context cancellation and avoid returning deleted rows.

## Common Mistakes Avoided

This avoided string-built SQL, forgetting `deleted_at IS NULL`, ignoring `RowsAffected` on soft delete, and swallowing scan errors.

## Commit Link

Commit: https://github.com/Iamfarhan-cs/crud-app/commit/2b97ba436a0cbd2c7acd42d537ff8aa5da90e6c4

# Phase E: HTTP Handlers Implementation

## Goal

Expose the task service through HTTP routes.

## What Changed

The project added route registration, collection and resource routing, JSON decoding, request body limits, pagination parsing, response encoding, and error-to-HTTP mapping.

## Files Changed

* `internal/task/handler.go` - HTTP handlers and route registration.
* `README.md` - documented handler behavior.

## Feature Flow

HTTP request -> route method switch -> decode or parse input -> service call -> JSON response or standard error.

## Engineering Reasoning

Handlers were kept focused on HTTP concerns. This keeps service logic reusable and lets tests target business behavior separately from transport behavior.

## Production Notes

Handlers need bounded bodies, stable error codes, correct status codes, and method handling. Unknown errors become `500 INTERNAL_ERROR` without exposing internals.

## Common Mistakes Avoided

This avoided putting SQL in handlers, returning inconsistent error shapes, and letting unsupported methods fall through silently.

## Commit Link

Commit: https://github.com/Iamfarhan-cs/crud-app/commit/9f5146f1830c1e35e47038f8474628d4e12e693b

# Phase F: Application Wiring

## Goal

Wire configuration, database startup, repository, service, handlers, routes, and graceful shutdown into a runnable API.

## What Changed

The project implemented environment config loading, PostgreSQL connection setup, server construction, health checks, route registration, startup logging, signal handling, and graceful shutdown.

## Files Changed

* `cmd/api/main.go` - application composition and HTTP server startup.
* `internal/config/config.go` - runtime configuration loading.
* `internal/database/postgres.go` - PostgreSQL connection setup and ping.
* `README.md` - documented application wiring.

## Feature Flow

Process start -> load config -> open PostgreSQL -> create repository -> create service -> create handler -> register routes -> serve HTTP -> graceful shutdown.

## Engineering Reasoning

Composition belongs at the edge of the application. This keeps internal packages reusable and makes the startup sequence clear.

## Production Notes

Startup should fail fast if `DATABASE_URL` is missing or PostgreSQL cannot be reached. Shutdown should stop accepting new work and give in-flight requests time to finish.

## Common Mistakes Avoided

This avoided hidden global dependencies, hardcoded database URLs, no ping on startup, and abrupt server shutdown.

## Commit Link

Commit: https://github.com/Iamfarhan-cs/crud-app/commit/f37215443e93a89a35955b6ff03e6e53fb3c9562

# Phase G: Database Migrations

## Goal

Create migration files that define the real PostgreSQL schema for tasks.

## What Changed

The migration creates the `tasks` table if needed, enforces title and status constraints, stores timestamps, supports soft delete, and creates an active-task partial index.

## Files Changed

* `migrations/001_create_tasks_table.up.sql` - task table and partial index.
* `migrations/001_create_tasks_table.down.sql` - rollback for index and table.
* `README.md` - documented migration behavior.

## Feature Flow

Migration up -> table exists -> application inserts app-generated UUID task -> active queries use index -> migration down can roll schema back.

## Engineering Reasoning

The schema protects business invariants close to the data. App-generated UUIDs keep ID creation inside the service, while database constraints protect against invalid writes.

## Production Notes

Migrations must be applied before running the API. Schema changes need review because they can block startup, writes, or deploy rollback.

## Common Mistakes Avoided

This avoided relying only on application validation, forgetting rollback SQL, and missing the partial index for active task listing.

## Commit Link

Commit: https://github.com/Iamfarhan-cs/crud-app/commit/bae6204dcb03a02c4211ce00744ca7cbaa6679be

# Phase H: Testing

## Goal

Add focused unit tests around business logic and HTTP behavior.

## What Changed

The project added service tests with a fake repository and handler tests using `net/http/httptest`.

## Files Changed

* `internal/task/service_test.go` - service unit tests and fake repository.
* `internal/task/handler_test.go` - handler unit tests.
* `README.md` - documented testing phase.

## Feature Flow

Test request or service call -> fake repository -> service behavior -> handler response or service result -> assertion.

## Engineering Reasoning

The tests target behavior without requiring a real database. This keeps feedback fast while still protecting the most important business and HTTP contracts.

## Production Notes

Unit tests do not replace integration tests. The next production step would be repository integration tests against PostgreSQL and CI automation.

## Common Mistakes Avoided

This avoided testing only happy paths, requiring Docker for every unit test, and leaving error mapping untested.

## Commit Link

Commit: https://github.com/Iamfarhan-cs/crud-app/commit/143fdd42217174657a6ce0b27a15f6a2133e9d3e

# Phase I: Dockerization

## Goal

Containerize the API and PostgreSQL database so the backend can be built and run consistently across developer machines.

## What Changed

The project added a multi-stage Docker build for the Go API, a Docker Compose setup for the API and PostgreSQL, a Docker ignore file, and an example environment file for local container configuration.

## Files Changed

* `Dockerfile` - multi-stage API image build.
* `docker-compose.yml` - local API and PostgreSQL orchestration.
* `.dockerignore` - excludes local and Git-only files from the Docker build context.
* `.env.example` - documents local environment values without production secrets.

## Feature Flow

Docker Compose -> PostgreSQL container starts -> migration SQL initializes the task table -> API container builds from the Dockerfile -> API connects to PostgreSQL through the Compose network.

## Engineering Reasoning

Dockerization was added after the application could run normally because containers should package known behavior, not hide unfinished wiring. The multi-stage build keeps the runtime image smaller by leaving Go build tooling in the builder image.

## Production Notes

The Compose setup is useful for local development, but production still needs real secret management, deployment configuration, backup strategy, and observability. Database readiness matters because the API depends on PostgreSQL being healthy before startup.

## Common Mistakes Avoided

This avoided shipping the full Go toolchain in the runtime image, copying `.git` and local environment files into the image context, and starting the API before PostgreSQL is ready.

## Commit Link

Commit: https://github.com/Iamfarhan-cs/crud-app/commit/5ca748284389b19fb51920ac3260fd105e6d56cf

# Build Phase K: Deployment Readiness

## Goal

Make the project easier to run, verify, and prepare for deployment by documenting operational commands, environment expectations, API usage, and production readiness checks.

## What Changed

The project added a Makefile with common local and Docker commands. The README was expanded into a deployment-ready project guide with architecture, environment variables, endpoint references, health checks, curl examples, test commands, reset commands, error format, deployment checklist, and production notes.

## Files Changed

* `Makefile` - added common run, test, formatting, vetting, and Docker lifecycle commands.
* `README.md` - expanded project documentation for local use and deployment readiness.
* `PHASES.md` - documented the deployment readiness phase.

## Makefile Flow

Developer command -> Make target -> underlying Go or Docker command.

`make run` starts the API with `go run ./cmd/api`. `make test`, `make fmt`, and `make vet` run the standard Go checks. `make docker-up`, `make docker-down`, and `make docker-reset` wrap the Docker Compose workflow so local setup and database resets use predictable commands.

## README Documentation Flow

Project overview -> features -> architecture -> setup -> environment variables -> migrations -> Docker workflow -> endpoints -> health and readiness checks -> curl examples -> tests -> deployment checklist -> production notes -> phase tracker.

This keeps the README useful both for a developer running the service locally and for someone reviewing whether the project is ready to deploy.

## Engineering Reasoning

Deployment readiness is mostly about repeatability. The Makefile avoids asking each developer to remember long commands, while the README turns operational assumptions into visible documentation.

The checklist captures the minimum runtime expectations before deploying: database configuration, migrations, health and readiness checks, logs, graceful shutdown, request limits, DB pooling, tests, and secret hygiene.

## Production Notes

This phase does not add production infrastructure by itself. It prepares the project for that work by making current runtime behavior explicit.

Future production improvements include authentication, authorization, CI/CD, an external migration runner, structured JSON logging, Prometheus metrics, OpenTelemetry tracing, rate limiting, caching, and a scaling strategy.

## Common Mistakes Avoided

This avoided burying operational commands in memory, omitting reset instructions for the Docker database, mixing production secrets into documentation, adding CI/CD before the project asks for it, and changing business logic during a documentation and deployment-readiness phase.

## How To Test

Run:

```bash
make test
make vet
```

Equivalent direct commands:

```bash
go test ./...
go vet ./...
```

Optional Docker check:

```bash
make docker-up
curl -i http://localhost:8080/healthz
curl -i http://localhost:8080/readyz
make docker-down
```

## Commit Link

Commit: https://github.com/Iamfarhan-cs/crud-app/commit/27a9aea214c9febfceea4bf0923d1e84ac6fe673

# Summary Table

| Phase | Feature | Main Files Changed | Commit Link | Status |
| ----- | ------- | ------------------ | ----------- | ------ |
| Phase 1 | Requirements and mental model | `README.md` | https://github.com/Iamfarhan-cs/crud-app/commit/1f342ffe1e8ac570a748c0065ac865aede55b7fa | Complete |
| Phase 2 | Architecture and boundaries | `cmd/api/main.go`, `internal/config/config.go`, `internal/database/postgres.go`, `internal/task/*`, `README.md` | https://github.com/Iamfarhan-cs/crud-app/commit/e428266e823e82e457c9cc1b110cc0e91a990e0a | Complete |
| Phase 3 | API design | `internal/task/dto.go`, `README.md` | https://github.com/Iamfarhan-cs/crud-app/commit/e58725ae30e7079ec38f29db14a45e1c1f9c91ac | Complete |
| Phase 4 | Database design | `migrations/*`, `README.md` | https://github.com/Iamfarhan-cs/crud-app/commit/f9fd8f4f3a8a0534f080a6082b399f6734c0df48 | Complete |
| Phase 5 | Storage design | `README.md` | https://github.com/Iamfarhan-cs/crud-app/commit/341c41d7668c063dd4b0b0459cd44ca69020e06d | Complete |
| Phase 6 | Networking design | `README.md` | https://github.com/Iamfarhan-cs/crud-app/commit/6c0811a231025332d15aab549b3b42a1a754f7f0 | Complete |
| Phase 7 | Security design | `README.md`, `internal/task/*`, `internal/config/config.go` | https://github.com/Iamfarhan-cs/crud-app/commit/ba01f12091327a66c8b7afbbe0e23041816550ed | Complete |
| Phase 8 | Concurrency design | `README.md` | https://github.com/Iamfarhan-cs/crud-app/commit/3e79ac8d759994d027ff08ee64b79d3565fdd6d4 | Complete |
| Phase 9 | Performance design | `README.md`, `internal/task/repository.go` | https://github.com/Iamfarhan-cs/crud-app/commit/d4a2bb96a993b801de763d5ba2cd4123b29fe769 | Complete |
| Phase A | Domain model | `internal/task/model.go`, `README.md` | https://github.com/Iamfarhan-cs/crud-app/commit/9d2d4b52504f2b97bdfb2caf7caf63d9e52883bd | Complete |
| Phase B | Errors and repository contract | `internal/task/errors.go`, `internal/task/repository.go`, `README.md` | https://github.com/Iamfarhan-cs/crud-app/commit/dc662268311d81f456ce95ab88869d6e45c97ae4 | Complete |
| Phase C | Service layer | `internal/task/service.go`, `README.md` | https://github.com/Iamfarhan-cs/crud-app/commit/4bb78684f83c42ab890a0bfcf90fafbc587a7e08 | Complete |
| Phase D | PostgreSQL repository | `internal/task/postgres_repository.go`, `README.md` | https://github.com/Iamfarhan-cs/crud-app/commit/2b97ba436a0cbd2c7acd42d537ff8aa5da90e6c4 | Complete |
| Phase E | HTTP handlers | `internal/task/handler.go`, `README.md` | https://github.com/Iamfarhan-cs/crud-app/commit/9f5146f1830c1e35e47038f8474628d4e12e693b | Complete |
| Phase F | Application wiring | `cmd/api/main.go`, `internal/config/config.go`, `internal/database/postgres.go`, `README.md` | https://github.com/Iamfarhan-cs/crud-app/commit/f37215443e93a89a35955b6ff03e6e53fb3c9562 | Complete |
| Phase G | Database migrations | `migrations/*`, `README.md` | https://github.com/Iamfarhan-cs/crud-app/commit/bae6204dcb03a02c4211ce00744ca7cbaa6679be | Complete |
| Phase H | Testing | `internal/task/service_test.go`, `internal/task/handler_test.go`, `README.md` | https://github.com/Iamfarhan-cs/crud-app/commit/143fdd42217174657a6ce0b27a15f6a2133e9d3e | Complete |
| Phase I | Dockerization | `Dockerfile`, `docker-compose.yml`, `.dockerignore`, `.env.example` | https://github.com/Iamfarhan-cs/crud-app/commit/5ca748284389b19fb51920ac3260fd105e6d56cf | Complete |
| Build Phase K | Deployment readiness | `Makefile`, `README.md`, `PHASES.md` | https://github.com/Iamfarhan-cs/crud-app/commit/27a9aea214c9febfceea4bf0923d1e84ac6fe673 | Complete |
