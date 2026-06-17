# Task Management CRUD API

## Phase 1: Requirements and Mental Model

### Project Purpose

This project is a production-minded CRUD backend for managing tasks. The API should allow clients to create, read, update, and delete task records while enforcing simple business rules and returning predictable responses.

The goal of Phase 1 is to define what the system is responsible for before choosing storage, routing, frameworks, database schema, or deployment details.

### Entity Definition

The core entity is a `Task`.

A task represents a unit of work that a user or system wants to track from creation through completion. It should be small enough to fit inside a task list, but detailed enough to support real workflow decisions.

### Fields

Expected task fields:

- `id`: unique identifier for the task
- `title`: short required summary of the work
- `description`: optional longer explanation
- `status`: current workflow state
- `priority`: relative importance of the task
- `due_date`: optional date by which the task should be completed
- `created_at`: timestamp when the task was created
- `updated_at`: timestamp when the task was last changed

Possible status values:

- `todo`
- `in_progress`
- `done`

Possible priority values:

- `low`
- `medium`
- `high`

### Business Rules

- A task must have a non-empty title.
- A task must have a valid status.
- A task must have a valid priority.
- New tasks should default to `todo` unless the client provides a valid status.
- New tasks should default to `medium` priority unless the client provides a valid priority.
- `created_at` should be set once and should not change after creation.
- `updated_at` should change whenever a task is modified.
- Deleting a task should remove it from normal list and read operations.
- The API should reject unknown or invalid task states instead of silently accepting bad data.

### API Operations

The API should support these operations:

- Create a task.
- List tasks.
- Retrieve a single task by ID.
- Update a task by ID.
- Delete a task by ID.

Useful production-minded query behavior may include:

- Filtering tasks by status.
- Filtering tasks by priority.
- Sorting tasks by creation date, update date, due date, or priority.
- Pagination for large task lists.

These query features can be added after the base CRUD behavior is stable.

### Non-Goals

Phase 1 does not include:

- Authentication or user accounts.
- Multi-user ownership rules.
- Team workspaces.
- Notifications or reminders.
- File attachments.
- Recurring tasks.
- Calendar integrations.
- Full-text search.
- Frontend UI.
- Deployment automation.

These are intentionally excluded so the first version can focus on correct CRUD behavior, validation, error handling, and maintainable structure.

### Production Risks

- Data loss if storage is in-memory or not backed up.
- Invalid data if validation is inconsistent between create and update paths.
- Race conditions if concurrent requests modify the same task without safe storage behavior.
- Unclear errors if the API does not return consistent status codes and response bodies.
- Unbounded list responses if pagination is ignored.
- Breaking clients if field names or response shapes change casually.
- Time handling bugs if timestamps and due dates are not represented consistently.
- Poor observability if failures are not logged with enough context.

### Failure Scenarios

The API should handle these situations intentionally:

- Client sends malformed JSON.
- Client omits a required field.
- Client sends an invalid status or priority.
- Client requests a task ID that does not exist.
- Client updates a task that was deleted.
- Client sends an unsupported HTTP method.
- Storage is unavailable or returns an error.
- Two clients update the same task close together.
- The task list becomes large enough that returning everything is slow.

### Junior Engineer Thinking vs Senior Engineer Thinking

Junior engineer thinking often starts with handlers and routes:

- "I need endpoints for create, read, update, and delete."
- "I need a struct with fields."
- "I need to return JSON."
- "I need to make the happy path work."

That thinking is useful, but incomplete.

Senior engineer thinking starts with system behavior and failure boundaries:

- "What is the entity, and what states are valid?"
- "Which fields are controlled by the client, and which are controlled by the server?"
- "What should happen when input is malformed or incomplete?"
- "How do create and update rules differ?"
- "What errors should clients be able to rely on?"
- "What happens when data grows?"
- "What assumptions will become expensive to change later?"

For this project, the senior mental model is simple: define the task lifecycle clearly, enforce valid data at the API boundary, keep response behavior predictable, and avoid adding features before the base CRUD contract is solid.

## Phase 2: Project Architecture and Boundaries

This phase establishes the Go project structure and the responsibility boundaries between packages. It does not implement complete CRUD behavior yet. The goal is to create a clean scaffold so later phases can add handlers, validation, persistence, and routing without mixing concerns.

### Module Structure

The project uses this structure:

```text
cmd/
  api/
    main.go
internal/
  config/
    config.go
  database/
    postgres.go
  task/
    dto.go
    errors.go
    handler.go
    model.go
    repository.go
    service.go
go.mod
README.md
```

### Package Responsibilities

`cmd/api` is the application entry point.

It should only compose the application:

- Load configuration.
- Prepare database connections.
- Create repositories.
- Create services.
- Create handlers.
- Start the HTTP server.

It must not contain business rules, SQL queries, request parsing, or response formatting.

`internal/config` owns runtime configuration.

It should handle values such as:

- Environment name.
- HTTP port.
- PostgreSQL connection string.

Later phases can add real environment variable parsing and validation here. Other packages should receive configuration values instead of reading environment variables directly.

`internal/database` owns database connection setup.

It should be responsible for creating and managing database connections. It should not know task business rules, HTTP route details, or request/response shapes.

`internal/task` owns the task feature.

It is split into smaller files by responsibility:

- `model.go`: domain entity and domain-level types such as `Task`, `Status`, and `Priority`.
- `dto.go`: request and response shapes used by the API layer.
- `errors.go`: stable task-specific errors.
- `repository.go`: persistence interface for task storage.
- `service.go`: business rules and use-case orchestration.
- `handler.go`: HTTP request handling and route registration.

### Dependency Direction

Dependencies should flow inward:

```text
HTTP handler -> service -> repository interface -> repository implementation
```

The handler should translate HTTP requests into service calls. The service should enforce business behavior and coordinate persistence. The repository should hide storage details behind an interface.

This keeps the business logic from depending on HTTP or database details. It also makes the system easier to test because handlers, services, and repositories can be tested separately.

### Current Scaffold

Phase 2 creates the architecture skeleton:

- A Go module named `github.com/Iamfarhan-cs/crud-app`.
- An API entry point in `cmd/api/main.go`.
- Placeholder configuration loading in `internal/config`.
- Placeholder PostgreSQL connection structure in `internal/database`.
- Task domain types in `internal/task/model.go`.
- API request and response DTOs in `internal/task/dto.go`.
- Task-specific domain errors in `internal/task/errors.go`.
- A task repository interface in `internal/task/repository.go`.
- A task service that depends on the repository interface.
- A task handler prepared for future route registration.

The placeholders are intentional. They let the project define ownership boundaries before adding implementation details.

### Phase 2 Non-Goals

Phase 2 does not include:

- Complete HTTP route registration.
- CRUD handler logic.
- Request validation.
- PostgreSQL connection opening.
- SQL queries or migrations.
- Repository implementation.
- Authentication.
- Tests.
- Deployment configuration.

These are deferred so the project can first settle its internal shape. Once responsibilities are clear, later phases can add behavior in the correct layer instead of crowding everything into `main.go` or HTTP handlers.

### Architecture Rules

- Domain types should not contain HTTP parsing or SQL logic.
- DTOs should describe external API shapes, not database tables.
- Services should own business rules and use-case flow.
- Repositories should own storage behavior behind interfaces.
- Handlers should own HTTP status codes, JSON decoding, JSON encoding, and route wiring.
- Configuration should be loaded once near application startup and passed into the components that need it.
- Database setup should live in the database package, not inside handlers or services.

## Phase 3: API Design

This phase defines the HTTP contract for the Task Management CRUD API. It does not include endpoint logic, database access, authentication, or handler implementation.

Base path:

```text
/api/v1
```

### Task API

| Method | Path | Purpose |
| --- | --- | --- |
| POST | `/api/v1/tasks` | Create a task |
| GET | `/api/v1/tasks` | List tasks |
| GET | `/api/v1/tasks/{id}` | Retrieve one task |
| PATCH | `/api/v1/tasks/{id}` | Partially update one task |
| DELETE | `/api/v1/tasks/{id}` | Delete one task |

### Field Ownership

Client-owned fields are accepted in create and update request bodies:

- `title`
- `description`
- `status`
- `priority`
- `due_date`

Server-owned fields are generated or controlled by the API and must not be accepted from clients:

- `id`
- `created_at`
- `updated_at`

The server sets `id`, `created_at`, and `updated_at` when a task is created. The server updates `updated_at` whenever a task changes. Clients may read these fields in responses but cannot write them.

### Task Representation

Task response example:

```json
{
  "id": "task_01HYZ7Y6D8P6M4C7R9K2Q3F5A1",
  "title": "Write API contract",
  "description": "Define request and response shapes before implementing handlers.",
  "status": "in_progress",
  "priority": "high",
  "due_date": "2026-07-01",
  "created_at": "2026-06-15T10:30:00Z",
  "updated_at": "2026-06-15T11:00:00Z"
}
```

Valid `status` values:

- `todo`
- `in_progress`
- `done`

Valid `priority` values:

- `low`
- `medium`
- `high`

`due_date` is optional and should use `YYYY-MM-DD` format.

### Create Task

```text
POST /api/v1/tasks
```

Request body:

```json
{
  "title": "Write API contract",
  "description": "Define request and response shapes before implementing handlers.",
  "status": "in_progress",
  "priority": "high",
  "due_date": "2026-07-01"
}
```

Request rules:

- `title` is required and must not be empty.
- `description` is optional.
- `status` is optional and defaults to `todo`.
- `priority` is optional and defaults to `medium`.
- `due_date` is optional.
- `id`, `created_at`, and `updated_at` are ignored or rejected because they are server-owned.

Success response:

```text
201 Created
```

```json
{
  "id": "task_01HYZ7Y6D8P6M4C7R9K2Q3F5A1",
  "title": "Write API contract",
  "description": "Define request and response shapes before implementing handlers.",
  "status": "in_progress",
  "priority": "high",
  "due_date": "2026-07-01",
  "created_at": "2026-06-15T10:30:00Z",
  "updated_at": "2026-06-15T10:30:00Z"
}
```

Possible status codes:

- `201 Created`: task created.
- `400 Bad Request`: malformed JSON or invalid field type.
- `422 Unprocessable Entity`: request JSON is valid but violates business rules.
- `500 Internal Server Error`: unexpected server failure.

### List Tasks

```text
GET /api/v1/tasks?page=1&per_page=20
```

Success response:

```text
200 OK
```

```json
{
  "data": [
    {
      "id": "task_01HYZ7Y6D8P6M4C7R9K2Q3F5A1",
      "title": "Write API contract",
      "description": "Define request and response shapes before implementing handlers.",
      "status": "in_progress",
      "priority": "high",
      "due_date": "2026-07-01",
      "created_at": "2026-06-15T10:30:00Z",
      "updated_at": "2026-06-15T11:00:00Z"
    }
  ],
  "pagination": {
    "page": 1,
    "per_page": 20,
    "total": 1,
    "total_pages": 1
  }
}
```

Pagination rules:

- `page` is optional and defaults to `1`.
- `per_page` is optional and defaults to `20`.
- `per_page` must have a maximum limit, recommended as `100`.
- `page` and `per_page` must be positive integers.
- Responses must include pagination metadata even when the result set is empty.
- The default sort order should be stable, such as newest tasks first by `created_at`.

Pagination is required because task lists can grow without bound. Returning every task in one response can make the API slow, increase memory usage, create large network responses, and make client behavior unpredictable as data grows.

Possible status codes:

- `200 OK`: tasks returned.
- `400 Bad Request`: invalid pagination query parameters.
- `500 Internal Server Error`: unexpected server failure.

### Get Task By ID

```text
GET /api/v1/tasks/{id}
```

Success response:

```text
200 OK
```

```json
{
  "id": "task_01HYZ7Y6D8P6M4C7R9K2Q3F5A1",
  "title": "Write API contract",
  "description": "Define request and response shapes before implementing handlers.",
  "status": "in_progress",
  "priority": "high",
  "due_date": "2026-07-01",
  "created_at": "2026-06-15T10:30:00Z",
  "updated_at": "2026-06-15T11:00:00Z"
}
```

Possible status codes:

- `200 OK`: task returned.
- `400 Bad Request`: invalid task ID format.
- `404 Not Found`: task does not exist.
- `500 Internal Server Error`: unexpected server failure.

### Update Task

```text
PATCH /api/v1/tasks/{id}
```

Request body:

```json
{
  "status": "done",
  "priority": "medium"
}
```

Request rules:

- All fields are optional, but the body must include at least one client-owned field.
- Only provided fields are changed.
- Omitted fields keep their current values.
- `title`, when provided, must not be empty.
- `status`, when provided, must be one of the allowed status values.
- `priority`, when provided, must be one of the allowed priority values.
- `id`, `created_at`, and `updated_at` are ignored or rejected because they are server-owned.

Success response:

```text
200 OK
```

```json
{
  "id": "task_01HYZ7Y6D8P6M4C7R9K2Q3F5A1",
  "title": "Write API contract",
  "description": "Define request and response shapes before implementing handlers.",
  "status": "done",
  "priority": "medium",
  "due_date": "2026-07-01",
  "created_at": "2026-06-15T10:30:00Z",
  "updated_at": "2026-06-15T12:15:00Z"
}
```

`PATCH` is used because the update operation is partial. Clients should be able to send only the fields they want to change without resending the full task representation. This reduces accidental overwrites and keeps update requests small.

Possible status codes:

- `200 OK`: task updated.
- `400 Bad Request`: malformed JSON, invalid task ID format, or empty update body.
- `404 Not Found`: task does not exist.
- `422 Unprocessable Entity`: request JSON is valid but violates business rules.
- `500 Internal Server Error`: unexpected server failure.

### Delete Task

```text
DELETE /api/v1/tasks/{id}
```

Success response:

```text
204 No Content
```

The response body must be empty.

Possible status codes:

- `204 No Content`: task deleted.
- `400 Bad Request`: invalid task ID format.
- `404 Not Found`: task does not exist.
- `500 Internal Server Error`: unexpected server failure.

### Standard Error Format

All error responses should use one stable JSON shape:

```json
{
  "error": {
    "code": "invalid_request",
    "message": "Request body contains invalid fields.",
    "details": {
      "title": "title is required",
      "status": "status must be one of: todo, in_progress, done"
    }
  }
}
```

Error format rules:

- `error.code` is a stable machine-readable string.
- `error.message` is a human-readable summary.
- `error.details` is optional and can hold field-specific validation messages.
- Internal server details, SQL errors, stack traces, and implementation internals must not be returned to clients.

Recommended error codes:

- `invalid_json`
- `invalid_request`
- `validation_failed`
- `not_found`
- `internal_error`

## Phase 4: Database Design

Phase 4 introduces the PostgreSQL schema for persistent task storage.

This phase creates the first database migrations:

- `migrations/001_create_tasks_table.up.sql`
- `migrations/001_create_tasks_table.down.sql`

The `tasks` table stores the core task record:

- `id`: UUID primary key generated by PostgreSQL.
- `title`: required task title.
- `description`: optional task details.
- `status`: required lifecycle state.
- `created_at`: timestamp for when the row was created.
- `updated_at`: timestamp for when the row was last changed.
- `deleted_at`: nullable timestamp used for soft deletes.

Database constraints enforce important invariants close to the data:

- `title` must not be empty or whitespace only.
- `status` must be one of `pending`, `in_progress`, or `done`.
- `created_at` and `updated_at` must always be present.

The schema also adds a partial index for normal list queries:

```sql
CREATE INDEX IF NOT EXISTS idx_tasks_active_created_at_desc
    ON tasks (created_at DESC)
    WHERE deleted_at IS NULL;
```

This keeps active task listings efficient without indexing soft-deleted rows.

### Phase 4 Non-Goals

Phase 4 does not include:

- Go repository implementations.
- SQL query methods.
- HTTP handlers.
- Route registration.
- Database connection wiring changes.

## Phase 5: Storage Design

This phase defines the storage behavior for the Task Management CRUD API before implementing repository SQL. The goal is to make task lifecycle rules clear at the persistence boundary so later repository code can follow one predictable contract.

### Soft Delete Model

Tasks use soft delete.

The task table should include a nullable `deleted_at` timestamp:

- `deleted_at IS NULL` means the task is active.
- `deleted_at IS NOT NULL` means the task has been deleted.

Delete operations must not physically remove task rows. A delete should set `deleted_at` to the time of deletion. This preserves historical data for future operational needs while keeping normal API behavior simple.

### Active Task Rules

Normal reads only return active tasks.

Repository reads should treat active tasks as the default view of the system:

- Listing tasks returns only rows where `deleted_at IS NULL`.
- Finding a task by ID returns only the matching row where `deleted_at IS NULL`.
- Deleted tasks behave like missing tasks for normal read paths.

The word `Active` appears in repository read method names so future implementations do not accidentally include soft-deleted rows in ordinary API responses.

### Update Rules

Deleted tasks cannot be updated.

Updates should only target active tasks. If a task has a non-null `deleted_at`, update behavior should be the same as trying to update a task that does not exist. This keeps the lifecycle simple: once deleted, a task is outside normal CRUD behavior.

The word `Active` appears in repository update method names so callers and implementations share the same expectation: update operations only apply to records where `deleted_at IS NULL`.

### Delete Rules

Deleting a task should set `deleted_at` instead of removing the row.

The repository method should be named `SoftDelete` to make the behavior explicit at the interface boundary. Later SQL can enforce this by updating the row and setting the deletion timestamp only for active tasks.

### V1 Non-Goals

Storage V1 does not include:

- Restore behavior for deleted tasks.
- Purge jobs or retention policies for old deleted tasks.
- Audit logging for create, update, or delete events.

Audit logging may be added later if the API needs a separate history of who changed what and when. That should be designed as its own storage concern instead of being mixed into the basic CRUD repository skeleton.

### Repository Interface Shape

The task repository should expose lifecycle-aware method names:

- `Create`
- `FindActiveByID`
- `ListActive`
- `UpdateActive`
- `SoftDelete`

These names keep storage rules visible without adding SQL details early. Full repository implementation, SQL queries, HTTP handlers, restore behavior, purge behavior, and audit logging are deferred.

## Phase 6: Networking Design

This phase defines how the Task Management CRUD API should behave as a network service. It does not implement CRUD endpoint logic or database repository logic.

### Port Configuration

Local development should default to port `8080`.

The application should read `PORT` through the config package instead of hardcoding the port in `cmd/api/main.go`. This keeps startup behavior portable across local machines, containers, and production platforms that assign ports dynamically.

In production, the Go process should usually run behind a load balancer or reverse proxy. The reverse proxy can own TLS termination, public routing, request buffering, and client-facing network policy while the Go service listens on its configured internal port.

### HTTP Server

The application should use an explicit `http.Server` instead of calling bare `http.ListenAndServe`.

An explicit server makes network behavior visible and configurable:

- `ReadTimeout` limits how long the server waits while reading the request.
- `WriteTimeout` limits how long the server spends writing the response.
- `IdleTimeout` limits how long keep-alive connections can sit unused.
- `Shutdown` supports graceful termination instead of abruptly dropping in-flight requests.

These settings help protect the service from slow or stuck clients and make production shutdown behavior predictable during deploys, restarts, and container termination.

### Request Size Limit

The API should apply a maximum request body size of `1MB`.

Task create and update payloads are small, so larger bodies are unexpected. A body limit reduces memory pressure and helps reject accidental or abusive oversized requests before CRUD handlers are implemented.

### Health Endpoints

The first networking endpoint is:

```text
GET /healthz
```

`/healthz` should return a simple success response when the process is running and able to serve HTTP. It should not require database access because it is a lightweight process health check.

`/readyz` is planned later.

Readiness should check dependencies needed to serve real traffic. For this API, PostgreSQL is a network dependency, so a future `/readyz` endpoint can verify whether the database connection is available before the service receives production traffic.

### Graceful Shutdown

The API should listen for termination signals and call `http.Server.Shutdown` with a bounded timeout.

Graceful shutdown gives active requests a short window to finish while the server stops accepting new requests. This matters behind load balancers and reverse proxies because deployment systems may terminate old processes while traffic is still draining.

### Phase 6 Non-Goals

Phase 6 does not include:

- CRUD route implementation.
- Request validation for task payloads.
- Database repository implementation.
- PostgreSQL readiness checks.
- TLS setup inside the Go process.
- Authentication, authorization, or rate limiting.

## Phase 7: Security Design

This phase defines the security model for the Task Management CRUD API. It does not implement authentication, authorization, CRUD handler logic, or database SQL.

### Authentication and Authorization

Authentication answers: "who are you?"

Authorization answers: "are you allowed to do this?"

Authentication is planned later. V1 can define clean boundaries now without adding login, sessions, tokens, password handling, or identity provider integration.

Task ownership is also planned later. A future version can add a `user_id` ownership field so authorization rules can decide which user is allowed to read, update, or delete a task. Until that is designed, handlers and services should not pretend ownership exists.

### Validation Boundaries

Input validation belongs partly in the handler and partly in the service.

Handlers should validate HTTP-level concerns:

- Request body format.
- Required JSON structure.
- Unsupported methods or paths.
- Maximum request body size.
- Conversion from request values into domain/service inputs.

Services should validate business rules:

- Required task title.
- Allowed status values.
- Allowed priority values.
- Create and update lifecycle rules.
- Server-owned fields such as `id`, `created_at`, and `updated_at`.

Critical invariants should also be enforced by database constraints. Application validation gives clients useful errors, but database constraints protect data integrity if another code path, migration, script, or future service writes to the same tables.

### SQL Safety

Repository implementations must use parameterized SQL only.

The code must never concatenate user input into SQL strings. Task IDs, titles, descriptions, status values, priority values, dates, pagination values, and future `user_id` values must be passed as query parameters.

This keeps SQL construction separate from user-controlled data and reduces the risk of SQL injection when repository logic is implemented in a later phase.

### Error Handling

Client error responses should be safe and stable.

API responses should use predictable error codes and clear messages without exposing stack traces, raw SQL errors, connection strings, table names, secret values, or internal implementation details.

Detailed errors belong in logs, not API responses. Logs can include enough context for operators to debug failures, while client responses should remain limited to what the caller needs to fix the request or understand the outcome.

### Secrets and Environment Files

Secrets must come from environment variables or a managed secret provider.

The repository must not commit real secrets, database passwords, API keys, private tokens, or production connection strings.

Local `.env` files must be ignored by Git. A `.env.example` file may be committed if it contains only placeholder names and safe example values. It should document required variables without containing real credentials.

### CORS

CORS should be restrictive in production.

The API should not default to allowing every origin in production. When browser clients are added, allowed origins should be configured explicitly for known frontend domains. Local development can be more flexible, but production should only allow trusted origins and required methods.

### Request Body Size

The request body size limit remains `1MB`.

Task payloads are small, so accepting larger bodies is unnecessary. Keeping the `1MB` limit reduces memory pressure and gives the server a simple first line of defense against accidental or abusive oversized requests.

### Phase 7 Non-Goals

Phase 7 does not include:

- Authentication implementation.
- Authorization implementation.
- User accounts.
- Task ownership columns or `user_id` migrations.
- CRUD handler implementation.
- Database SQL implementation.
- CORS middleware implementation.
- Secret manager integration.

## Phase 8: Concurrency Design

This phase defines how the Task Management CRUD API should think about concurrent requests. It does not implement full concurrency control, add database schema changes, or add CRUD handlers yet.

Concurrency in a CRUD API means more than many requests arriving at the same time. It means the API must behave predictably when two clients read, create, update, or delete related data close together. The storage layer may serialize individual SQL statements, but the API still needs clear rules for request races, retries, stale reads, and soft-deleted records.

### Lost Update Problem

A lost update happens when two clients read the same task, both make changes, and the later write overwrites the earlier write without knowing it.

Example:

1. Client A reads task `123` with status `todo`.
2. Client B reads task `123` with status `todo`.
3. Client A updates the status to `in_progress`.
4. Client B updates the priority using an older copy of the task and accidentally writes status `todo` again.

The API should avoid update behavior that replaces fields the client did not intentionally change. `PATCH` must only update fields sent in the request body. Omitted fields must keep their existing values so partial updates do not erase concurrent changes by accident.

As a planned future improvement, the API should use optimistic locking. The task table can later add a `version INTEGER` field that increments on every successful update. Clients would send the version they last saw, and repository update operations would only succeed when the stored version still matches. A stale update should return `409 Conflict` instead of silently overwriting newer data.

### Duplicate Create Problem

Duplicate creates can happen when a client sends `POST /api/v1/tasks`, times out before receiving the response, and retries the same request. The first request may have succeeded even though the client never saw the `201 Created` response. Without an idempotency strategy, the retry can create a second task with the same intent.

As a planned future improvement, `POST /api/v1/tasks` should support an idempotency key. A client-provided idempotency key would let the API detect a retry of the same create request and return the original result instead of creating a duplicate task.

This phase does not add idempotency storage or request handling. It only records the design expectation so future create logic has a clear direction.

### Update-After-Delete Race

An update-after-delete race happens when one request deletes a task while another request tries to update the same task.

Example:

1. Client A reads task `123`.
2. Client B deletes task `123`.
3. Client A sends a patch for task `123`.

The update must not modify a soft-deleted row. Repository update and delete operations must include `deleted_at IS NULL` in their predicates so normal CRUD behavior only targets active tasks. If an update finds no active row, the service should treat that as a not-found result unless a future version check proves the client is stale and should receive `409 Conflict`.

Delete operations should also target only active rows. Repeating a delete against an already deleted task should not update `deleted_at` again as though the task were still active.

### Future Concurrency Rules

- `PATCH` should only update fields that were explicitly sent by the client.
- Repository update operations must include `deleted_at IS NULL`.
- Repository delete operations must include `deleted_at IS NULL`.
- Optimistic locking is planned for a later phase, using a future `version INTEGER` field.
- Stale updates detected by version mismatch should return `409 Conflict`.
- `POST` retries should eventually support an idempotency key to prevent duplicate creates.
- This phase does not change the database schema; `version INTEGER` is future work.

## Build Phase D: PostgreSQL Repository Implementation

Build Phase D adds the PostgreSQL-backed task repository in `internal/task/postgres_repository.go`.

The repository implements the task `Repository` interface using `database/sql` and keeps SQL behavior inside the persistence layer. Services continue to depend on the interface instead of depending directly on PostgreSQL.

### Repository Construction

The PostgreSQL repository is represented by:

- `PostgresRepository`
- `NewPostgresRepository(db *sql.DB) *PostgresRepository`

The repository stores a shared `*sql.DB` handle and uses context-aware database calls for every operation.

### Implemented Methods

The PostgreSQL repository implements:

- `Create`
- `FindActiveByID`
- `ListActive`
- `UpdateActive`
- `SoftDelete`

`Create` inserts task fields and uses `RETURNING` to read back the stored row.

`FindActiveByID` selects a task by ID only when `deleted_at IS NULL`.

`ListActive` selects active tasks ordered by `created_at DESC` and uses `LIMIT` and `OFFSET` placeholders for bounded list reads.

`UpdateActive` updates task fields only when the row is active and uses `RETURNING` to read back the updated row.

`SoftDelete` updates `deleted_at` and `updated_at` with `NOW()` instead of physically deleting the row.

### SQL Safety and Lifecycle Rules

All SQL uses positional parameters. No user input is concatenated into SQL strings.

The repository avoids `SELECT *` and explicitly selects task columns:
## Build Phase C: Service Layer Implementation

Build Phase C implements the task service in `internal/task/service.go`.

The service layer now owns task business rules and use-case orchestration while continuing to depend on the `Repository` interface instead of PostgreSQL directly.

### Service Methods

The task service exposes these use-case methods:

- `CreateTask(ctx context.Context, req CreateTaskRequest) (TaskResponse, error)`
- `GetTask(ctx context.Context, id string) (TaskResponse, error)`
- `ListTasks(ctx context.Context, page int, limit int) ([]TaskResponse, error)`
- `UpdateTask(ctx context.Context, id string, req UpdateTaskRequest) (TaskResponse, error)`
- `DeleteTask(ctx context.Context, id string) error`

### Create Behavior

Create behavior now includes service-level defaults and server-owned fields:

- Titles are trimmed and must not be empty.
- Missing status defaults to `pending`.
- Provided status must be `pending`, `in_progress`, or `done`.
- IDs are generated with `uuid.NewString()`.
- `CreatedAt` and `UpdatedAt` are set with `time.Now().UTC()`.

### Update Behavior

Update behavior follows PATCH semantics:

- A request must include at least one update field.
- A nil field means the client did not provide that field.
- Provided titles are trimmed and must not be empty.
- Provided status values must be valid.
- The service loads the existing active task before applying changes.
- `UpdatedAt` is refreshed with `time.Now().UTC()`.

### Read, List, and Delete Behavior

The service trims task IDs before read and delete operations. An empty ID returns `ErrTaskNotFound`.

List behavior normalizes pagination:

- Page `0` defaults to `1`.
- Limit `0` defaults to `20`.
- Page must be at least `1`.
- Limit must be between `1` and `100`.
- Invalid pagination returns `ErrInvalidPagination`.

Delete behavior calls `SoftDelete`, keeping the soft-delete policy at the repository boundary.

### Service Helpers

The service also includes small helper functions:

- `validateTitle`
- `validateStatus`
- `normalizePagination`
- `toResponse`

These keep validation, pagination normalization, and response mapping focused inside the service layer.

### Supporting Contract Alignment

Because this branch was created from `main`, Build Phase C also aligns the task contracts that the service depends on:

- `Task` now includes `Description *string` and `DeletedAt *time.Time`.
- `Status` now uses `pending`, `in_progress`, and `done`.
- Create, update, and response DTOs match the service-owned task fields.
- Repository list operations accept `limit` and `offset`.
- Task sentinel errors include validation, pagination, not-found, and empty-update cases.

### Build Phase C Non-Goals

Build Phase C does not include:

- PostgreSQL repository implementation.
- SQL queries.
- HTTP handlers.
- Route registration.
- Authentication or authorization.
## Build Phase B: Errors and Repository Contract

Build Phase B defines the task error vocabulary in `internal/task/errors.go` and the persistence boundary in `internal/task/repository.go`.

### Sentinel Errors

The task package now exposes sentinel errors for stable failure cases:

- `ErrTaskNotFound`
- `ErrInvalidTitle`
- `ErrInvalidStatus`
- `ErrInvalidPagination`
- `ErrNoFieldsToUpdate`

These errors give upper layers stable values that can later be mapped to HTTP responses. Handlers can translate them into response codes and bodies without tying service or repository code to HTTP behavior.

Wrapped errors must preserve `errors.Is` compatibility. This lets future service or repository implementations add context to an error while still allowing callers to detect the original task failure category.

The repository is an interface, not a PostgreSQL implementation. This keeps the service layer dependent on a stable task storage contract instead of depending directly on PostgreSQL, SQL queries, connection details, or any specific database package.

### Repository Methods

The task repository contract now exposes these methods:

- `Create(ctx context.Context, task Task) (Task, error)`
- `FindActiveByID(ctx context.Context, id string) (Task, error)`
- `ListActive(ctx context.Context, limit int, offset int) ([]Task, error)`
- `UpdateActive(ctx context.Context, task Task) (Task, error)`
- `SoftDelete(ctx context.Context, id string) error`

`ListActive` accepts `limit` and `offset` so future list queries can stay bounded and support pagination without changing the repository shape later.

### Active Task Meaning

`Active` means `deleted_at IS NULL`.

This makes the soft-delete lifecycle visible in the method names:

- Active reads ignore soft-deleted tasks.
- Active lists return only non-deleted tasks.
- Active updates cannot modify deleted tasks.
- Soft delete marks a task as deleted instead of physically removing the row.

### Repository Boundary Rules

Repository implementations must use parameterized SQL and must never concatenate user input into query strings.

The repository must not contain HTTP concerns such as:

- Status codes.
- Request parsing.
- Response formatting.
- Route behavior.

The repository also must not make business-policy decisions. Those belong in the service layer, which can decide how to interpret repository results and which domain rules apply to create, update, read, and delete workflows.

### Build Phase B Non-Goals

Build Phase B does not include:

- PostgreSQL query implementation.
- Service logic.
- HTTP handlers.
- Route registration.
- Request validation.
- Business-rule enforcement.
- HTTP error mapping implementation.
## Build Phase A: Domain Model

Build Phase A defines the task domain model and the API-facing request and response shapes in `internal/task/model.go`.

The task lifecycle now uses a focused `Status` string type:

- `pending`
- `in_progress`
- `done`

The `Task` domain model represents the application's internal view of a task:

- `ID`
- `Title`
- `Description`
- `Status`
- `CreatedAt`
- `UpdatedAt`
- `DeletedAt`

`DeletedAt` is included in the domain model because soft-delete state is part of the internal lifecycle. It is intentionally not exposed in API responses.

### Request Models

`CreateTaskRequest` contains only client-owned create fields:

- `title`
- `description`
- `status`

`UpdateTaskRequest` contains optional pointers for partial updates:

- `title`
- `description`
- `status`

Server-owned fields such as `id`, `created_at`, `updated_at`, and `deleted_at` are not accepted in create or update requests. This keeps ownership clear: clients describe the task content, while the server owns identifiers, timestamps, and deletion state.

### Response Model

`TaskResponse` defines the public task representation returned by the API:

- `id`
- `title`
- `description`
- `status`
- `created_at`
- `updated_at`
- `deleted_at`

Normal reads exclude soft-deleted rows with `deleted_at IS NULL`.

Updates also require `deleted_at IS NULL`, so soft-deleted tasks cannot be modified by normal update operations.

Missing active records return `ErrTaskNotFound`. Unexpected database failures are wrapped with contextual messages using `fmt.Errorf` so callers can preserve the original error while still getting useful operational context.

### Supporting Contract Alignment

Because this branch was created from `main`, Build Phase D also aligns the task contracts needed by the PostgreSQL repository:

- `Task` includes nullable `Description` and `DeletedAt` fields.
- `Status` uses the database-supported values `pending`, `in_progress`, and `done`.
- Create, update, and response DTOs match the current task table fields.
- `Repository.ListActive` accepts `limit` and `offset`.
- Task sentinel errors include not-found, validation, pagination, and empty-update cases.

### Build Phase D Non-Goals

Build Phase D does not include:

- HTTP handlers.
- Route registration.
- Application wiring.
- Database connection startup changes.
- Service behavior changes.
- PostgreSQL migrations beyond the existing task table migration.

`deleted_at` is not exposed because deleted tasks should be handled as lifecycle state inside the system, not as part of the normal public task response.

### Model Separation

The domain model, request models, and response model are intentionally separate:

- The domain model can hold internal lifecycle fields.
- Request models restrict what clients are allowed to send.
- Response models restrict what the API exposes back to clients.

This separation prevents accidental exposure of server-owned fields and gives later validation, service, handler, and repository phases clear boundaries to build on.

### Build Phase A Non-Goals

Build Phase A does not include:

- Request validation.
- HTTP handlers.
- Route registration.
- Repository implementation.
- SQL query logic.
- Service behavior changes.
## Phase 9: Performance Design

This phase defines the first performance expectations for the Task Management CRUD API. It does not implement caching, full CRUD handlers, advanced database optimization, or production observability tooling yet.

Performance should be designed around the real request path:

```text
client -> network -> Go HTTP server -> handler -> service -> repository -> PostgreSQL -> repository -> service -> handler -> network -> client
```

Each step can add latency. Slow clients, large request bodies, JSON parsing, validation, connection pool waits, database query time, row scanning, response encoding, and network transfer all contribute to the final response time. The API should keep each layer simple, bounded, and measurable before adding advanced optimization.

### List Endpoint Pagination

List endpoints must be paginated.

Returning every task in one response becomes slower and more expensive as the table grows. Large unbounded responses increase database work, memory usage, JSON encoding time, network transfer time, and client-side processing. Pagination keeps list responses predictable.

V1 should use offset pagination:

```text
GET /api/v1/tasks?limit=20&offset=0
```

Pagination rules:

- Default `limit` should be `20`.
- Maximum `limit` should be `100`.
- `offset` should default to `0`.
- `limit` and `offset` must be non-negative integers.
- Requests above the maximum limit should be rejected or clamped according to the API validation rules chosen in the handler phase.

Offset pagination is simple and good enough for V1. Cursor pagination is a planned future improvement for large datasets or high-write workloads where deep offsets become expensive or unstable.

### Active Task Listing Index

The active task listing query needs an index because normal list operations should only return rows where `deleted_at IS NULL`, usually in a stable order such as newest tasks first.

A query shaped like this should not scan the whole table as data grows:

```sql
SELECT id, title, description, status, created_at, updated_at, deleted_at
FROM tasks
WHERE deleted_at IS NULL
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;
```

The Phase 4 partial index on active tasks supports this access pattern by indexing active rows in creation order. That keeps common list reads focused on the rows the API actually returns.

### Query Shape

Repository queries should avoid `SELECT *`.

Explicit column lists make the API more predictable:

- They avoid returning unused columns.
- They protect repository scanning code from accidental table changes.
- They make query cost easier to understand.
- They document exactly which fields the repository needs.

Future repository implementations should select only the columns required for the operation. List queries may eventually use a lighter projection than detail queries if the API response shape needs it.

### Database Connection Pooling

The API should use PostgreSQL connection pooling through the Go database driver or pool library instead of opening a new database connection per request.

Connection pool settings should be explicit:

- `max open connections`: caps total concurrent database connections.
- `max idle connections`: keeps a bounded number of warm reusable connections.
- `connection max lifetime`: periodically retires old connections so long-running processes do not hold stale connections forever.

These values should be tuned for the deployment environment and PostgreSQL capacity. Setting them too high can overload the database. Setting them too low can make requests wait for a free connection even when PostgreSQL has capacity.

### Query Timeouts

Repository methods should accept `context.Context`.

The context allows handlers or services to attach deadlines, cancellation, and request-scoped values before calling the repository. Query timeouts prevent slow database work from running forever after a client disconnects or after the API has already exceeded its response budget.

Future SQL calls should use context-aware driver methods, such as `QueryContext`, `QueryRowContext`, or equivalent methods from the chosen PostgreSQL library.

### Observability

The API should log slow queries.

Slow query logs help identify missing indexes, inefficient query shapes, pool starvation, and data growth problems before they become outages. Logs should include operation names and durations, but should not expose sensitive request data.

The API should also monitor latency percentiles:

- `p95` latency shows the experience of slower-but-common requests.
- `p99` latency shows tail behavior and helps reveal rare but painful slow paths.

Averages are not enough because a small number of very slow requests can hurt users while the average still looks acceptable.

### Realistic Test Data

Performance checks should use realistic test data.

Testing with five tasks does not reveal list latency, index behavior, pagination cost, or slow query patterns. Future performance validation should seed enough rows to approximate real usage, including active tasks, soft-deleted tasks, mixed statuses, mixed priorities, and varied creation dates.

Realistic test data makes query plans, pagination behavior, and API response sizes closer to production conditions.

### Phase 9 Non-Goals

Phase 9 does not include:

- Caching.
- Full CRUD handler implementation.
- Advanced query optimization.
- Cursor pagination implementation.
- Production dashboards.
- Load testing automation.
- Database schema changes.
