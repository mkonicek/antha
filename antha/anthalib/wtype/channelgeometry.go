package wtype

import (
	"github.com/pkg/errors"
	"math"

	"github.com/antha-lang/antha/antha/anthalib/wunit"
)

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
	return false
}

// LinearChannelSpacing represents channel spacing in one dimension
type LinearChannelSpacing interface {
	// TimeToMoveTo calculate the time taken for the cones to align with the given locations (in mm).
	// nil values are unspecified and can be moved to any given location, which will be set by the
	// function call.
	// Returns an error if spacing is invalid or if the cones cannot reach the given positions
	TimeToMoveBetween(initial []float64, final []*float64) (wunit.Time, error)
}

// FixedChannelSpacing the channels in this direction cannot move relative to the
// adaptor
type FixedChannelSpacing struct {
	Spacing wunit.Length
}

// TimeToMoveBetween
func (self *FixedChannelSpacing) TimeToMoveBetween(initial []float64, final []*float64) (wunit.Time, error) {
	spacing := self.Spacing.MustInStringUnit("mm")
	ret := wunit.NewTime(0, "s")

	for i, p := range final {
		if pos := float64(i) * spacing; p == nil {
			final[i] = &pos
		} else {
			if math.Abs(*final[i]-pos) > 1e-5 {
				return ret, errors.New("fixed adaptor cannot reach position")
			}
		}
	}

	return ret, nil
}

// ExtendableChannelSpacing the channels in this direction can independently extend
// relative to the adaptor
type ExtendableChannelSpacing struct {
	Extension        []wunit.Length // the current or initial extension, should equal zero or MaximumExtension if PartialExtension is unset
	MaximumExtension wunit.Length   // the maximum distance that the channels can extend
	ExtensionSpeed   wunit.Velocity // the speed at which the cones are extended
	PartialExtension bool           // if set to True, then the head support extending partially in this direction, otherwise cones must be at zero or MaximumExtension
}

// TimeToMoveBetween
func (self *ExtendableChannelSpacing) TimeToMoveBetween(initial []float64, final []*float64) (wunit.Time, error) {
	var maxDistance float64

	//unspecified values stay where they are
	for i := range final {
		if final[i] == nil {
			p := initial[i]
			final[i] = &p
		}
	}

	maxExtend := self.Extension.MustInStringUnit("mm")

	positionOk := func(pos float64) bool {
		return pos >= 0 && pos <= maxExtend
	}
	if self.PartialExtension {
		positionOk = func(pos float64) bool {
			return pos == 0 || pos == maxExtend
		}
	}

	for i, f := range final {
		if !positionOk(*f) {
			return wunit.NewTime(0, "s"), errors.New("extendable adaptor cannot reach position")
		} else if d := math.Abs(*f - initial[i]); d > maxDistance {
			maxDistance = d
		}
	}

	return wunit.NewTime(maxDistance/self.ExtensionSpeed.MustInStringUnit("mm/s"), "s"), nil
}

// RelativeChannelSpacing the spacing between each channel can vary freely between the given limits
type RelativeChannelSpacing struct {
	MinSpacing wunit.Length   // the minimum spacing between the channels
	MaxSpacing wunit.Length   // the maximum spacing between the channels
	Speed      wunit.Velocity // the speed at which channels can splay or contract
}

// TimeToMoveBetween
func (self *RelativeChannelSpacing) TimeToMoveBetween(initial []float64, final []*float64) (wunit.Time, error) {
	if len(initial) != len(final) {
		return wunit.NewTime(0, "s"), errors.New("initial and final spacing have mismatched lengths - this indicates a bug")
	}

	if len(initial) < 2 {
		return wunit.NewTime(0, "s"), nil
	}

	initialSpacing := make([]float64, len(initial)-1)
	for i, v := range initial[1:] {
		initialSpacing[i] = v - initial[i]
	}

	// fill internal nills by evenly interpolating
	// it's possible that robots which suport this behaviour could do something smarter to reduce
	// time taken for channel movement, but realistically this method minimises the possibility of
	// collision which probably most robot manufacturers prefer
	lastIndex := -1
	for i, f := range final {
		if f == nil {
			continue
		}
		if lastIndex != -1 {
			step := *f - *final[lastIndex]
			for j := 1; i < i-lastIndex-1; j++ {
				v := *final[lastIndex] + float64(j)*step
				final[lastIndex+j] = &v
			}
		}
		lastIndex = i
	}

	// special case: no positions specified
	if lastIndex == -1 {
		for i, initialValue := range initial {
			v := initialValue
			final[i] = &v
		}
		return wunit.NewTime(0, "s"), nil
	}

	// fill edge nils using minimum spacing - least likely to be invalid
	step := self.MinSpacing.MustInStringUnit("mm")
	firstNonNil := -1
	lastNonNil := -1
	for i, v := range final {
		if v == nil {
			continue
		}
		if firstNonNil < 0 {
			firstNonNil = i
		}
		lastNonNil = i
	}
	for i := lastNonNil + 1; i < len(final); i++ {
		v := *final[lastNonNil] + float64(i-lastNonNil)*step
		final[i] = &v
	}
	for i := 0; i < firstNonNil; i++ {
		v := *final[firstNonNil] + float64(i-firstNonNil)*step
		final[i] = &v
	}

	finalSpacing := make([]float64, len(final)-1)
	for i, v := range final[1:] {
		finalSpacing[i] = *v - *final[i]
	}

	//now check that the spacing is valid
	minSpacing := self.MinSpacing.MustInStringUnit("mm")
	maxSpacing := self.MaxSpacing.MustInStringUnit("mm")
	for _, s := range finalSpacing {
		if s < minSpacing || s > maxSpacing {
			return wunit.NewTime(0, "s"), errors.New("channel cannot move to the requested position")
		}
	}

	//find the furthest any channel has to move
	maxDistance := 0.0
	for i, f := range final {
		if d := math.Abs(*f - initial[i]); d > maxDistance {
			maxDistance = d
		}
	}

	return wunit.NewTime(maxDistance/self.Speed.MustInStringUnit("mm/s"), "s"), nil
}

// AccordionChannelSpacing a special case of RelativeChannelSpacing where the spacing between each channel centre can vary within the given limits,
// but the spacing between each channel must remain equal
type AccordionChannelSpacing struct {
	MinSpacing wunit.Length   // the minimum spacing between the channels
	MaxSpacing wunit.Length   // the maximum spacing between the channels
	Speed      wunit.Velocity // the speed at which channels can splay or contract
}
