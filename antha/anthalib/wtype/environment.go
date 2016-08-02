package wtype

import (
	"github.com/antha-lang/antha/antha/anthalib/wunit"
)

type Environment struct {
	Temperature         wunit.Temperature
	Pressure            wunit.Pressure
	Humidity            float64
	MeanAirFlowVelocity wunit.Velocity
	/// etc
}
