package api

const (
	// Label for inventory components after execution, hopefully temporary
	AfterLabel = "v1.transitional.report.after="
)

// swagger:model ReportTask
type Task struct {
	Id string `json:"id"`
	// Short description
	Label string `json:"label"`
	// Long description
	Details string `json:"details"`
	// If device task, the device id
	DeviceId string `json:"device_id"`
	// Details on setup
	DeviceConfig *DeviceConfig `json:"device_config,omitempty"`
	// Time estimate in seconds, 0 if no estimate available
	TimeEstimate float64 `json:"time_estimate"`
	// Task.Ids
	HappensBefore []string `json:"happens_before"`
	// Status of a task
	Status Status `json:"status,omitempty"`
	// Error if task state appropriate
	Error string `json:"error,omitempty"`
}
