package wtype

import (
	"github.com/pkg/errors"
	"strings"

	"github.com/antha-lang/antha/antha/anthalib/wunit"
)

type Axis int

const (
	XAxis Axis = iota
	YAxis
	ZAxis
)

func (a Axis) String() string {
	switch a {
	case XAxis:
		return "X"
	case YAxis:
		return "Y"
	case ZAxis:
		return "Z"
	}
	panic("unknown axis")
}

func AxisFromString(a string) (Axis, error) {
	switch strings.ToUpper(a) {
	case "X":
		return XAxis, nil
	case "Y":
		return YAxis, nil
	case "Z":
		return ZAxis, nil
	}
	return Axis(-1), errors.Errorf("unknown axis %q", a)
}

// LHHeadAssemblyPosition a position within a head assembly
type LHHeadAssemblyPosition struct {
	Offset Coordinates
	Head   *LHHead
}

// Velocity3D struct composed of lengths in three axes
type Velocity3D struct {
	X, Y, Z wunit.Velocity
}

// GetAxis return the velocity in the axis specified
func (self *Velocity3D) GetAxis(a Axis) wunit.Velocity {
	switch a {
	case XAxis:
		return self.X
	case YAxis:
		return self.Y
	case ZAxis:
		return self.Z
	}
	panic("unknown axis")
}

// SetAxis return the velocity in the axis specified
func (self *Velocity3D) SetAxis(a Axis, v wunit.Velocity) {
	switch a {
	case XAxis:
		self.X = v
	case YAxis:
		self.Y = v
	case ZAxis:
		self.Z = v
	}
	panic("unknown axis")
}

// Dup return a copy of the velocities
func (self *Velocity3D) Dup() *Velocity3D {
	if self == nil {
		return nil
	}
	return &Velocity3D{
		X: wunit.NewVelocity(self.X.RawValue(), self.X.Unit().PrefixedSymbol()),
		Y: wunit.NewVelocity(self.Y.RawValue(), self.Y.Unit().PrefixedSymbol()),
		Z: wunit.NewVelocity(self.Z.RawValue(), self.Z.Unit().PrefixedSymbol()),
	}
}

// VelocityRange the minimum and maximum velocities for the head assembly.
// nil implies no limit
type VelocityRange struct {
	Min, Max *Velocity3D
}

// Dup return a copy of the range
func (self *VelocityRange) Dup() *VelocityRange {
	if self == nil {
		return nil
	}
	return &VelocityRange{
		Min: self.Min.Dup(),
		Max: self.Max.Dup(),
	}
}

// Acceleration3D acceleration in three axes
type Acceleration3D struct {
	X, Y, Z *wunit.Acceleration
}

// Dup return a copy of the accelerations
func (self *Acceleration3D) Dup() *Acceleration3D {
	if self == nil {
		return nil
	}
	x := wunit.NewAcceleration(self.X.RawValue(), self.X.Unit().PrefixedSymbol())
	y := wunit.NewAcceleration(self.Y.RawValue(), self.Y.Unit().PrefixedSymbol())
	z := wunit.NewAcceleration(self.Z.RawValue(), self.Z.Unit().PrefixedSymbol())
	return &Acceleration3D{
		X: &x,
		Y: &y,
		Z: &z,
	}
}

// AccelerationRange minimum and maximum accelerations for the head assembly.
// nil implies no limit
type AccelerationRange struct {
	Min, Max *Acceleration3D
}

// Dup return a copy of the range
func (self *AccelerationRange) Dup() *AccelerationRange {
	if self == nil {
		return nil
	}
	return &AccelerationRange{
		Min: self.Min.Dup(),
		Max: self.Max.Dup(),
	}
}

//LHHeadAssembly represent a set of LHHeads which are constrained to move together
type LHHeadAssembly struct {
	Positions    []*LHHeadAssemblyPosition
	MotionLimits *BBox              //the limits on range of motion of the head assembly, nil if unspecified
	Velocity     *VelocityRange     // the range of valid velocities for the head, nil if unspecified
	Acceleration *AccelerationRange // the range of valid accelerations for the head, nil if unspecified
}

//NewLHHeadAssembly build a new head assembly
func NewLHHeadAssembly(MotionLimits *BBox) *LHHeadAssembly {
	ret := LHHeadAssembly{
		Positions:    make([]*LHHeadAssemblyPosition, 0, 2),
		MotionLimits: MotionLimits,
	}
	return &ret
}

//DupWithoutHeads copy the headassembly leaving all positions in the new assembly unloaded
func (self *LHHeadAssembly) DupWithoutHeads() *LHHeadAssembly {
	ret := &LHHeadAssembly{
		Positions:    make([]*LHHeadAssemblyPosition, 0, len(self.Positions)),
		MotionLimits: self.MotionLimits.Dup(),
	}
	for _, pos := range self.Positions {
		ret.AddPosition(pos.Offset)
	}
	return ret
}

//AddPosition add a position to the head assembly with the given offset
func (self *LHHeadAssembly) AddPosition(Offset Coordinates) {
	p := LHHeadAssemblyPosition{
		Offset: Offset,
	}
	self.Positions = append(self.Positions, &p)
}

//GetNumPositions get the number of positions added to the head assembly
func (self *LHHeadAssembly) CountPositions() int {
	return len(self.Positions)
}

//GetNumHeadsLoaded get the number of heads that are loaded into the assembly
func (self *LHHeadAssembly) CountHeadsLoaded() int {
	if self == nil {
		return 0
	}

	var r int
	for _, pos := range self.Positions {
		if pos.Head != nil {
			r += 1
		}
	}
	return r
}

//GetLoadedHeads get an ordered slice of all the heads that have been loaded into the assembly
func (self *LHHeadAssembly) GetLoadedHeads() []*LHHead {
	if self == nil {
		return make([]*LHHead, 0)
	}
	ret := make([]*LHHead, 0, len(self.Positions))
	for _, pos := range self.Positions {
		if pos.Head != nil {
			ret = append(ret, pos.Head)
		}
	}
	return ret
}

//LoadHead load a head into the next available position in the assembly, returns error if no positions
//are available
func (self *LHHeadAssembly) LoadHead(head *LHHead) error {
	if self == nil {
		return errors.New("cannot load head to nil assembly")
	}
	for _, pos := range self.Positions {
		if pos.Head == nil {
			pos.Head = head
			return nil
		}
	}
	return errors.New("cannot load head")
}

//UnloadHead unload a head from the assembly, return an error if the head is not loaded
func (self *LHHeadAssembly) UnloadHead(head *LHHead) error {
	if self == nil {
		return nil
	}
	for _, pos := range self.Positions {
		if pos.Head != nil && pos.Head.ID == head.ID {
			pos.Head = head
			return nil
		}
	}
	return errors.New("cannot load head")
}

//UnloadAllHeads unload all heads from the assembly
func (self *LHHeadAssembly) UnloadAllHeads() {
	if self == nil {
		return
	}
	for _, pos := range self.Positions {
		pos.Head = nil
	}
}
