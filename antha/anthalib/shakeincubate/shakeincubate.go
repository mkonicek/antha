// Package for helping to set up and run the incubator; designed for interacting with anthaOS
package shake

import (
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/driver"
	shakerincubator "github.com/antha-lang/antha/driver/antha_shakerincubator_v1"
	"github.com/antha-lang/antha/execute"
)

// Produce driver calls to send to incubator for set up
func SetPoints(component *wtype.LHComponent, temp wunit.Temperature, device string, rpm float64) []driver.Call {
	return []driver.Call{
		driver.Call{
			Method: "/antha.shakerincubator.v1.ShakerIncubator/TemperatureSet",
			Args: &shakerincubator.TemperatureSettings{
				Temperature: temp.RawValue(), // in C
			},
			Reply: &shakerincubator.BoolReply{},
		},
		driver.Call{
			Method: "/antha.shakerincubator.v1.ShakerIncubator/ShakeStart",
			Args: &shakerincubator.ShakerSettings{
				Frequency: rpm / 60.0,   // RPM to Hz
				Radius:    3.0 / 1000.0, // 3 mm
			},
			Reply: &shakerincubator.BoolReply{},
		},
	}
}

// Adds step to set up incubator in AnthaOS timeline
func SetUp(component *wtype.LHComponent) execute.HandleOpt {
	return execute.HandleOpt{
		Label:     "setup incubator",
		Component: component,
	}
}

// Turn on incubator using incubator settings from SetUpShakerIncubator function
func TurnOn(component *wtype.LHComponent, incubatorsettings []driver.Call) execute.HandleOpt {
	return execute.HandleOpt{
		Label: "turn on incubator",
		Selector: map[string]string{
			"antha.driver.v1.TypeReply.type": "antha.shakerincubator.v1.ShakerIncubator",
		},
		Calls:     incubatorsettings,
		Component: component,
	}
}

func turnOff() []driver.Call {
	return []driver.Call{
		driver.Call{
			Method: "/antha.shakerincubator.v1.ShakerIncubator/ShakeStop",
			Args:   &shakerincubator.Blank{},
			Reply:  &shakerincubator.BoolReply{},
		},
		driver.Call{
			Method: "/antha.shakerincubator.v1.ShakerIncubator/TemperatureReset",
			Args:   &shakerincubator.Blank{},
			Reply:  &shakerincubator.BoolReply{},
		},
	}
}

// Turn off Incubator
func TurnOff(component *wtype.LHComponent) execute.HandleOpt {
	return execute.HandleOpt{
		Label: "turn off incubator",
		Selector: map[string]string{
			"antha.driver.v1.TypeReply.type": "antha.shakerincubator.v1.ShakerIncubator",
		},
		Calls:     turnOff(),
		Component: component,
	}
}
