package api

const (
	// Task created (intial state)
	StateCreated = "Created"
	// Task eligible to run
	StateScheduled = "Scheduled"
	// Task waiting on external input to run
	StateWaiting = "Waiting"
	// Task running
	StateRunning = "Running"
	// Task finished running successfully
	StateSucceeded = "Succeeded"
	// Task finished running unsuccessfully
	StateFailed = "Failed"
)

type Status struct {
	// Current state
	State string `json:"state"`
	// Any message associated with the current state
	Message string `json:"message"`
}
