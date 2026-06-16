package task

// Repository defines the persistence boundary for tasks.
// Implementations should store and retrieve task data.
// Repository implementations must use parameterized SQL only and must never
// concatenate user input into query strings.
// HTTP concerns, authentication, authorization, business policy decisions,
// and database connection ownership must not live here.
type Repository interface {
	// Persistence method signatures will be added when the CRUD contract is finalized.
}
