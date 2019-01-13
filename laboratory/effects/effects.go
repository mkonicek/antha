package effects

import (
	"encoding/json"

	"github.com/antha-lang/antha/inventory/cache/plateCache"
	"github.com/antha-lang/antha/inventory/testinventory"
	"github.com/antha-lang/antha/laboratory/effects/id"
	"github.com/antha-lang/antha/microArch/sampletracker"
)

type LaboratoryEffects struct {
	Trace         *Trace
	Maker         *Maker
	SampleTracker *sampletracker.SampleTracker
	Inventory     *testinventory.TestInventory `json:"-"` // Inventory is part of plate cache, so we don't encode it.
	PlateCache    *plateCache.PlateCache
	IDGenerator   *id.IDGenerator
}

func NewLaboratoryEffects(jobId string) *LaboratoryEffects {
	idGen := id.NewIDGenerator(jobId)
	le := &LaboratoryEffects{
		Trace:         NewTrace(),
		Maker:         NewMaker(),
		SampleTracker: sampletracker.NewSampleTracker(),
		Inventory:     testinventory.NewInventory(idGen),
		IDGenerator:   idGen,
	}
	le.PlateCache = plateCache.NewPlateCache(le.Inventory)
	return le
}

func (le *LaboratoryEffects) ToJSON() ([]byte, error) {
	return json.Marshal(le)
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
