package effects

import (
	"github.com/antha-lang/antha/instructions"
	"github.com/antha-lang/antha/inventory"
	"github.com/antha-lang/antha/inventory/cache/plateCache"
	"github.com/antha-lang/antha/laboratory/effects/id"
	"github.com/antha-lang/antha/microArch/sampletracker"
	"github.com/antha-lang/antha/workflow"
)

type LaboratoryEffects struct {
	FileManager   *FileManager
	Trace         *instructions.Trace
	Maker         *instructions.Maker
	SampleTracker *sampletracker.SampleTracker
	Inventory     *inventory.Inventory
	// TODO the plate cache should go away and the use sites should be rewritten around plate types.
	PlateCache  *plateCache.PlateCache
	IDGenerator *id.IDGenerator
}

func NewLaboratoryEffects(wf *workflow.Workflow, simId workflow.BasicId, fm *FileManager) *LaboratoryEffects {
	idGen := id.NewIDGenerator(string(simId))
	le := &LaboratoryEffects{
		FileManager:   fm,
		Trace:         instructions.NewTrace(),
		Maker:         instructions.NewMaker(),
		SampleTracker: sampletracker.NewSampleTracker(),
		Inventory:     inventory.NewInventory(idGen),
		IDGenerator:   idGen,
	}
	le.PlateCache = plateCache.NewPlateCache(le.Inventory.Plates)
	le.Inventory.LoadForWorkflow(wf)

	return le
}
