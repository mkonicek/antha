package ast

import (
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/driver"
)

// An IncubateInst is a high-level command to incubate a component
type IncubateInst struct {
	// Time for which to incubate component
	Time wunit.Time
	// Temperature at which to incubate component
	Temp wunit.Temperature
	// Rate at which to shake incubator (force is device dependent)
	ShakeRate wunit.Rate
	// Radius at which ShakeRate is defined
	ShakeRadius wunit.Length

	// Time for which to pre-heat incubator
	PreTemp wunit.Temperature
	// Temperature at which to pre-heat incubator
	PreTime wunit.Time
	// Rate at which to pre-heat incubator
	PreShakeRate wunit.Rate
	// Radius at which PreShakeRate is defined
	PreShakeRadius wunit.Length
}

// An HandleInst is a high-level generic command to apply some device
// specific action to a component
type HandleInst struct {
	Group    string
	Selector map[string]string
	Calls    []driver.Call
}

// GetID returns a custom key for generic grouping
func (h *HandleInst) GetID() string {
	return h.Group
}

// A PromptInst is a high-level command to prompt a human
type PromptInst struct {
	Message string
}

// An ExpectInst is a command that...
type ExpectInst struct {
	// user-definable device tags
	Tags []string
	// ID we are waiting on
	ID string
}
