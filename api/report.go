package api

import "errors"

type Outputs map[string]map[string]interface{}

// A simulation report.
//
// swagger:model Report
type Report struct {
	// Id
	Id string `json:"id"`
	// Tasks that comprise this report
	Tasks []*Task `json:"tasks"`
	// Outputs from workflow execution
	Outputs Outputs `json:"outputs"`
	// Estimate of total time for execution
	TotalTime Measurement `json:"total_time"`
}

// Extract task errors
func (a *Report) Errors() []error {
	var errs []error
	for _, t := range a.Tasks {
		if t.Status.State == StateFailed {
			errs = append(errs, errors.New(t.Error))
		}
	}
	return errs
}

// Extract items and remove duplicates
//
// Normalizes item values for pointer comparisons.
func (a *Report) DeviceItems() []*DeviceItem {
	seen := make(map[string]*DeviceItem)
	var items []*DeviceItem
	for _, t := range a.Tasks {
		dcfg := t.DeviceConfig
		if dcfg == nil {
			continue
		}
		for _, item := range t.DeviceConfig.DeviceItems() {
			if seen[item.Id] != nil {
				continue
			}
			seen[item.Id] = item
			items = append(items, item)
		}
	}
	return items
}

func (a *Report) OutputItems() (items []*DeviceItem) {
	for _, t := range a.Tasks {
		if t.DeviceConfig == nil {
			continue
		}
		items = append(items, t.DeviceConfig.DeviceOutputs...)
	}
	return
}
