package effects

import (
	"fmt"
	"time"

	"github.com/antha-lang/antha/composer"
	"github.com/antha-lang/antha/inventory"
	"github.com/antha-lang/antha/inventory/cache/plateCache"
	"github.com/antha-lang/antha/laboratory/effects/id"
	"github.com/antha-lang/antha/microArch/sampletracker"
)

type LaboratoryEffects struct {
	JobId string

	Trace         *Trace
	Maker         *Maker
	SampleTracker *sampletracker.SampleTracker
	Inventory     *inventory.Inventory
	// TODO the plate cache should go away and the use sites should be rewritten around plate types.
	PlateCache  *plateCache.PlateCache
	IDGenerator *id.IDGenerator
}

func NewLaboratoryEffects(jobId string, wf *composer.Workflow) *LaboratoryEffects {
	idGen := id.NewIDGenerator(jobId)
	le := &LaboratoryEffects{
		JobId:         jobId,
		Trace:         NewTrace(),
		Maker:         NewMaker(),
		SampleTracker: sampletracker.NewSampleTracker(),
		Inventory:     inventory.NewInventory(idGen),
		IDGenerator:   idGen,
	}
	le.PlateCache = plateCache.NewPlateCache(le.Inventory.PlateTypes)

	// TODO: discuss this: not sure if we want to do this based off
	// zero plate types defined, or if we want an explicit flag or
	// something?
	if len(wf.Inventory.PlateTypes) == 0 {
		start := time.Now()
		le.Inventory.PlateTypes.LoadLibrary()
		fmt.Println("Loaded default plate types in", time.Now().Sub(start))
	} else {
		le.Inventory.PlateTypes.SetPlateTypes(wf.Inventory.PlateTypes)
	}

	return le
}
