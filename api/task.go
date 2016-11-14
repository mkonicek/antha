package api

// swagger:model Task
type Task struct {
	Id string `json:"id"`
	// Short description
	Label string `json:"label"`
	// Long description
	Details string `json:"details"`
	// If device task, the device id
	DeviceId string `json:"device_id"`
	// Time estimate in seconds, 0 if no estimate available
	TimeEstimate float64 `json:"time_estimate"`
	// Task.Ids
	HappensBefore []string `json:"happens_before"`
	// Status of a task
	Status Status `json:"status,omitempty"`
	// Tags
	Tags []string `json:"tags"`

	OrderTask      *OrderTask      `json:"order_task,omitempty"`
	DeckLayoutTask *DeckLayoutTask `json:"deck_layout_task,omitempty"`
	PlatePrepTask  *PlatePrepTask  `json:"plate_prep_task,omitempty"`
	DocumentTask   *DocumentTask   `json:"document_task,omitempty"`
	MixerTask      *MixerTask      `json:"mixer_task,omitempty"`
	ManualRunTask  *ManualRunTask  `json:"manual_run_task,omitempty"`
	IncubateTask   *IncubateTask   `json:"incubate_task,omitempty"`
	DataUploadTask *DataUploadTask `json:"data_upload_task,omitempty"`
}

// Order inventory items
type OrderTask struct {
	// Inventory items to order
	InventoryIds []string `json:"inventory_ids"`
}

// Show deck layout for mixer
type DeckLayoutTask struct {
	// Mixer task to show deck layout of
	MixerTaskId string `json:"mixer_task_id"`
	// If present, restrict layout to given deck positions, e.g.,
	SomePositions []OrdinalCoord `json:"some_positions,omitempty"`
}

// Prepare plates
type PlatePrepTask struct {
	// Plates to prepare
	PlatePreps []PlatePrep `json:"plate_preps"`
}

// Prepare a plate
type PlatePrep struct {
	// Plate to prepare
	PlateId string `json:"plate_id"`
	// If present, restrict prep to given well addresses, e.g., A1, BB2
	SomeWells []OrdinalCoord `json:"some_wells,omitempty"`
}

// Show documentation
type DocumentTask struct {
	// Unformated text to show
	TextBody string `json:"text_body"`
}

// Run mixer
type MixerTask struct {
	// Setup: input state of device
	Before []*MixerState `json:"before"`
	// Result: output state of device
	After []*MixerState `json:"after"`
	// Low level device calls representing this task
	Calls []*GrpcCall `json:"calls"`
}

type MixerState struct {
	Items      []*InventoryItem `json:"items"`
	Placements []*Placement     `json:"placements"`
}

type Placement struct {
	Parent string `json:"parent"`
	Child  string `json:"child"`
	// Symbolic location of child in coordinate system of parent
	OrdinalCoord *OrdinalCoord `json:"coord"`
}

// Manually initiated task
type ManualRunTask struct {
}

// Run incubator
type IncubateTask struct {
}

// Upload data
type DataUploadTask struct {
}
