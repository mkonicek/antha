package ast

import (
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/driver"
)

type IncubateInst struct {
	Time wunit.Time
	Temp wunit.Temperature
}

type HandleInst struct {
	Group    string
	Selector map[string]string
	Calls    []driver.Call
}
