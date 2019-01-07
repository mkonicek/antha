package effects

import (
	"encoding/json"

	"github.com/antha-lang/antha/inventory/cache/plateCache"
	"github.com/antha-lang/antha/inventory/testinventory"
	"github.com/antha-lang/antha/microArch/sampletracker"
)

type LaboratoryEffects struct {
	Trace         *Trace
	Maker         *Maker
	SampleTracker *sampletracker.SampleTracker
	Inventory     *testinventory.TestInventory `json:"-"` // Inventory is part of plate cache, so we don't encode it.
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

func (le *LaboratoryEffects) FromJSON(bs []byte) error {
	// the default json marshaling is fine, we just need to do some rewriting of pointers:
	if err := json.Unmarshal(bs, le); err != nil {
		return err
	} else {
		le.Inventory = le.PlateCache.Inventory
		return nil
	}
}
