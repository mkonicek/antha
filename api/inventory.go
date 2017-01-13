package api

import (
	"encoding/json"
	"time"
)

type InventoryItem struct {
	// Inventory id
	Id string `json:"id"`
	// Metadata
	Metadata map[string]json.RawMessage `json:"metadata,omitempty"`
	// Time this inventory item was created at
	CreatedAt time.Time `json:"created_at"`
	// From: history of this inventory item
	From []*InventoryItem `json:"from,omitempty"`

	Tipbox       *Tipbox       `json:"tipbox,omitempty"`
	Tipwaste     *Tipwaste     `json:"tipwaste,omitempty"`
	Plate        *Plate        `json:"plate,omitempty"`
	DeckPosition *DeckPosition `json:"deck_position,omitempty"`
	Component    *Component    `json:"component,omitempty"`
}

// Pipette tips in a box
type Tipbox struct {
	// Tipbox type
	Type string `json:"type"`
}

// Disposal for used pipette tips
type Tipwaste struct {
	// Tipwaste type
	Type string `json:"type"`
}

// Synthetic inventory item to represent position on deck
type DeckPosition struct {
	// Position
	Position *OrdinalCoord `json:"position"`
}

// Plate
type Plate struct {
	// Plate type
	Type  string  `json:"type"`
	Wells []*Well `json:"wells,omitempty"`
}

// Well in plate
type Well struct {
	Position  *OrdinalCoord  `json:"position"`
	Component *InventoryItem `json:"component"`
}

// Physical component, typically a liquid
type Component struct {
	// Component type
	Type string `json:"type"`
	// Name
	Name string `json:"name"`
	// Volume
	Volume *Measurement `json:"volume,omitempty"`
	// Viscosity
	Viscosity *Measurement `json:"viscosity,omitempty"`
	// Mass
	Mass *Measurement `json:"mass,omitempty"`
	// Amount (moles)
	Amount *Measurement `json:"amount,omitempty"`
	// If non-atomic component, this is what we are comprised of
	Components []*Component `json:"components,omitempty"`

	//TotalVolume Measurement `json:"total_volume"`
	//Concentration      Measurement `json:"concentration"`
	//SMax               float64     `json:"smax"`
	//StockConcentration Measurement `json:"stock_concentration"`
	//PlateLocation      DeviceCoord `json:"plate_location"`
}
