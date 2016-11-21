// Package for helping to set up interactions with anthaOS
package setup

import (
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/execute"
)

// Adds plate prep step in AnthaOS timeline
func PlatePrep(component *wtype.LHComponent) execute.HandleOpt {
	return execute.HandleOpt{
		Label:     "plate prep",
		Component: component,
	}

}

// Adds manual step in AnthaOS timeline
func MarkForSetup(component *wtype.LHComponent) execute.HandleOpt {
	return execute.HandleOpt{Label: "setup",
		Component: component,
	}
}

// Adds step for component ordering in AnthaOS timeline
func OrderInfo(component *wtype.LHComponent) execute.HandleOpt {
	return execute.HandleOpt{Label: "order",
		Component: component,
	}
}
