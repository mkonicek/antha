package inventory

import (
	"encoding/json"

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

type inventoryJSON struct {
	TipBoxes   *tipboxes.Inventory
	PlateTypes *plates.Inventory
}

func (inv *Inventory) MarshalJSON() ([]byte, error) {
	// currently we only marshal the plate types and the tipboxes
	toJson := &inventoryJSON{
		TipBoxes:   inv.TipBoxes,
		PlateTypes: inv.PlateTypes,
	}
	return json.Marshal(toJson)
}

func (inv *Inventory) UnmarshalJSON(bs []byte) error {
	// this creates a read only inventory with plenty of nils. Be careful!
	fromJson := &inventoryJSON{}
	if err := json.Unmarshal(bs, fromJson); err != nil {
		return err
	} else {
		inv.TipBoxes = fromJson.TipBoxes
		inv.PlateTypes = fromJson.PlateTypes
		return nil
	}
}
