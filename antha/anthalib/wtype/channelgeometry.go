package wtype

// ChannelGeometry describes the layout of cones within an adaptor and their
// movement limits
type ChannelGeometry struct {
	Spacing      *ChannelSpacing
	ConeGeometry Shape
}

type ChannelSpacing struct {
	ChannelsX int
	ChannelsY int
	spacing   [3]LinearChannelSpacing
}

// CanMoveTo determine whether the adaptor channels can move to the given set of positions.
// spacing is a two dimensional array with size self.ChannelsYxself.ChannelX
// such that spacing[y][x] refers to the requested position of the channel at
// (x,y).
// If any position is nil, it is unspecified and that cone can be moved to any
// location.
func (self *ChannelSpacing) CanMoveTo(spacing [][]*Coordinates) bool {
}

// LinearChannelSpacing represents channel spacing in one dimension
type LinearChannelSpacing interface {
	// CanMoveTo determine whether the cones can move to the given locations,
	// nil values are unspecified and can be moved to any given location
	CanMoveTo(spacing []*wunit.Length) bool
}
