// Package for helping to set up and run a thermocycler; designed for interacting with anthaOS
package thermocycle

import (
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/execute"
)

// Adds step to set up incubator in AnthaOS timeline
func SetUp(component *wtype.LHComponent) execute.HandleOpt {
	return execute.HandleOpt{
		Label:     "setup thermocycler",
		Component: component,
	}
}
