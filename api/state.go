package api

const (
	StateCreated   = "Created"
	StateScheduled = "Queued"
	StateRunning   = "Running"
	StateSucceeded = "Succeeded"
	StateFailed    = "Failed"
)

type Status struct {
	State   string `json:"state"`
	Message string `json:"message"`
}
