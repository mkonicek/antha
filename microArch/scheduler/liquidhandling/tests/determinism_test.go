package tests

import (
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/laboratory"
)

func GetComponentForTest(lab *laboratory.Laboratory, name string, vol wunit.Volume) *wtype.Liquid {
	c, err := lab.Inventory.Components.NewComponent(name)
	if err != nil {
		panic(err)
	}
	c.SetVolume(vol)
	return c
}
