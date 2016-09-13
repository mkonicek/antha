package ast

import "github.com/antha-lang/antha/antha/anthalib/wunit"

type IncubateInst struct {
	Time wunit.Time
	Temp wunit.Temperature
}

type HandleInst struct {
	Group string
}
