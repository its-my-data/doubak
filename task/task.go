package task

// BaseInterface defines the interface of each type of task.
type BaseInterface interface {
	// Checking flag combinations, validities, etc.
	Precheck() error

	// Execute the task.
	Execute() error
}
