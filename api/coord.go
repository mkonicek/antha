package api

// Zero-indexed coordinate system in ordinal space: ith item in X, Y, Z space.
// Origin is back, left, bottom (i.e., left-handed)
type OrdinalCoord struct {
	X int `json:"x,omitempty"`
	Y int `json:"y,omitempty"`
	Z int `json:"z,omitempty"`
}
