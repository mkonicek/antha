package inventory

import (
	"github.com/antha-lang/antha/inventory/components"
	"github.com/antha-lang/antha/inventory/plates"
	"github.com/antha-lang/antha/inventory/tipboxes"
	"github.com/antha-lang/antha/inventory/tipwastes"
	"github.com/antha-lang/antha/laboratory/effects/id"
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
