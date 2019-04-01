package inventory

import (
	"github.com/antha-lang/antha/inventory/components"
	"github.com/antha-lang/antha/inventory/plates"
	"github.com/antha-lang/antha/inventory/tipboxes"
	"github.com/antha-lang/antha/inventory/tipwastes"
	"github.com/antha-lang/antha/laboratory/effects/id"
	"github.com/antha-lang/antha/workflow"
)

type Inventory struct {
	Components *components.Inventory
	TipWastes  *tipwastes.Inventory
	TipBoxes   *tipboxes.Inventory
	PlateTypes *plates.Inventory
}

func NewInventory(idGen *id.IDGenerator) *Inventory {
	return &Inventory{
		Components: components.NewInventory(idGen),
		TipWastes:  tipwastes.NewInventory(idGen),
		TipBoxes:   tipboxes.NewInventory(idGen),
		PlateTypes: plates.NewInventory(idGen),
	}
}

func (inv *Inventory) LoadForWorkflow(wf *workflow.Workflow) {
	// TODO: discuss this: not sure if we want to do this based off
	// zero plate types defined, or if we want an explicit flag or
	// something?
	if wf == nil || len(wf.Inventory.PlateTypes) == 0 {
		inv.PlateTypes.LoadLibrary()
	} else {
		inv.PlateTypes.SetPlateTypes(wf.Inventory.PlateTypes)
	}

	inv.TipBoxes.LoadLibrary()
}
