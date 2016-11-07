package api

type DeviceCoord struct {
	// Id of container that holds this input
	Parent string `json:"parent"`
	// Descriptive label of this point in parent container
	Label string `json:"label,omitempty"`
	// Origin is lower front left (user's perspective) of the machine [0, \inf)
	X int `json:"x,omitempty"`
	// [0, \inf)
	Y int `json:"y,omitempty"`
	// [0, \inf)
	Z int `json:"z,omitempty"`
}

func (a DeviceCoord) Dim(x int) int {
	switch x {
	case 0:
		return a.X
	case 1:
		return a.Y
	case 2:
		return a.Z
	default:
		return 0.0
	}
}
