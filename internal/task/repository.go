package task

// Repository defines the persistence boundary for tasks.
// Implementations should store and retrieve task data.
// HTTP concerns, business policy decisions, and database connection ownership must not live here.
type Repository interface {
	// Persistence method signatures will be added when the CRUD contract is finalized.
}
