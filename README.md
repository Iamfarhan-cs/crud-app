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
