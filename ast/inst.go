package ast

import (
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	api "github.com/antha-lang/antha/api/v1"
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

// A PromptInst is a high-level command to prompt a human
type PromptInst struct {
	Message string
}

// An AwaitInst is a command that suspends execution pending data input
type AwaitInst struct {
	// user-definable devic3 tags
	Tags []string
	// ID we are waiting on
	AwaitID string
	// Next element in recursive chain
	NextElement string
	// Parameters to next element
	NextElementParams api.ElementParameters
	// Parameter that will receive the awaited data
	ReplaceParam string
}
