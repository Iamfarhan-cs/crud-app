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
