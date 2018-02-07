package wtype

type SBSLabware interface {
	NumRows() int
	NumCols() int
	PlateHeight() float64
}
