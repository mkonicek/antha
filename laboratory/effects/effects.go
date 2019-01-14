package effects

import (
	"encoding/json"

	"github.com/antha-lang/antha/inventory/cache/plateCache"
	"github.com/antha-lang/antha/inventory/testinventory"
	"github.com/antha-lang/antha/laboratory/effects/id"
	"github.com/antha-lang/antha/microArch/sampletracker"
)

type LaboratoryEffects struct {
	JobId string

	Trace         *Trace
	Maker         *Maker
	SampleTracker *sampletracker.SampleTracker
	Inventory     *testinventory.TestInventory
	PlateCache    *plateCache.PlateCache
	IDGenerator   *id.IDGenerator
}

func NewLaboratoryEffects(jobId string) *LaboratoryEffects {
	idGen := id.NewIDGenerator(jobId)
	le := &LaboratoryEffects{
		JobId:         jobId,
		Trace:         NewTrace(),
		Maker:         NewMaker(),
		SampleTracker: sampletracker.NewSampleTracker(),
		Inventory:     testinventory.NewInventory(idGen),
		IDGenerator:   idGen,
	}
	le.PlateCache = plateCache.NewPlateCache(le.Inventory)
	return le
}

type laboratoryEffectsJSON struct {
	JobId         string
	Trace         *Trace
	Maker         *Maker
	SampleTracker *sampletracker.SampleTracker
	// Inventory is part of plate cache, so we don't encode it.
	PlateCache  *plateCache.PlateCache
	IDGenerator *id.IDGenerator
}

func (le *LaboratoryEffects) MarshalJSON() ([]byte, error) {
	return json.Marshal(&laboratoryEffectsJSON{
		JobId:         le.JobId,
		Trace:         le.Trace,
		Maker:         le.Maker,
		SampleTracker: le.SampleTracker,
		PlateCache:    le.PlateCache,
		IDGenerator:   le.IDGenerator,
	})
}

func (le *LaboratoryEffects) UnmarshalJSON(bs []byte) error {
	lej := &laboratoryEffectsJSON{}
	if err := json.Unmarshal(bs, lej); err != nil {
		return err
	} else {
		le.JobId = lej.JobId
		le.Trace = lej.Trace
		le.Maker = lej.Maker
		le.SampleTracker = lej.SampleTracker
		le.PlateCache = lej.PlateCache
		le.IDGenerator = lej.IDGenerator
		// Just need to do a little rewiring:
		le.Inventory = le.PlateCache.Inventory
		le.Inventory.SetIDGenerator(le.IDGenerator)
		return nil
	}
}
