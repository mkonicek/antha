package effects

import (
	"github.com/antha-lang/antha/inventory"
	"github.com/antha-lang/antha/inventory/cache/plateCache"
	"github.com/antha-lang/antha/laboratory/effects/id"
	"github.com/antha-lang/antha/microArch/sampletracker"
	"github.com/antha-lang/antha/workflow"
)

type LaboratoryEffects struct {
	JobId workflow.JobId

	FileManager   *FileManager
	Trace         *Trace
	Maker         *Maker
	SampleTracker *sampletracker.SampleTracker
	Inventory     *inventory.Inventory
	// TODO the plate cache should go away and the use sites should be rewritten around plate types.
	PlateCache  *plateCache.PlateCache
	IDGenerator *id.IDGenerator
}

func NewLaboratoryEffects(jobId workflow.JobId, fm *FileManager, inv *inventory.Inventory) *LaboratoryEffects {
	idGen := id.NewIDGenerator(string(jobId))
	if inv == nil {
		inv = inventory.NewInventory(idGen)
	}
	le := &LaboratoryEffects{
		JobId:         jobId,
		FileManager:   fm,
		Trace:         NewTrace(),
		Maker:         NewMaker(),
		SampleTracker: sampletracker.NewSampleTracker(),
		Inventory:     inv,
		IDGenerator:   idGen,
	}
	le.PlateCache = plateCache.NewPlateCache(le.Inventory.PlateTypes)

	return le
}
