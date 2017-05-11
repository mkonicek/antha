package wtype

type Haswelladdress interface {
	PlateID() string
	WellCoords() WellCoords
}
