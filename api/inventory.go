package api

type InventoryItem struct {
	Id        string            `json:"id"`
	Type      string            `json:"type"`
	Metadata  map[string][]byte `json:"metadata"`
	Tipbox    *Tipbox           `json:"tipbox,omitempty"`
	Tipwaste  *Tipwaste         `json:"tipwaste,omitempty"`
	Plate     *Plate            `json:"plate,omitempty"`
	Component *Component        `json:"component,omitempty"`
}

type Tipbox struct {
	Tips []DeviceCoord `json:"used_tips"`
}

type Tipwaste struct {
	Tips []DeviceCoord `json:"used_tips"`
}

type Plate struct {
}

type Component struct {
	Name               string      `json:"component_name"`
	Volume             Measurement `json:"volume"`
	Concentration      Measurement `json:"concentration"`
	TotalVolume        Measurement `json:"total_volume"`
	SMax               float64     `json:"smax"`
	Viscosity          float64     `json:"viscosity"`
	StockConcentration Measurement `json:"stock_concentration"`
	PlateLocation      DeviceCoord `json:"plate_location"`
}
