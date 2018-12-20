package effects

import (
	"github.com/antha-lang/antha/inventory/cache/plateCache"
	"github.com/antha-lang/antha/inventory/testinventory"
	"github.com/antha-lang/antha/microArch/sampletracker"
)

type LaboratoryEffects struct {
	Trace         *Trace
	Maker         *Maker
	SampleTracker *sampletracker.SampleTracker
	Inventory     *testinventory.TestInventory
	PlateCache    *plateCache.PlateCache
}

func NewLaboratoryEffects() *LaboratoryEffects {
	le := &LaboratoryEffects{
		Trace:         NewTrace(),
		Maker:         NewMaker(),
		SampleTracker: sampletracker.NewSampleTracker(),
		Inventory:     testinventory.NewInventory(),
	}
	le.PlateCache = plateCache.NewPlateCache(le.Inventory)
	return le
}
