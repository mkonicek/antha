package wtype

import (
	"github.com/pkg/errors"

	"github.com/antha-lang/antha/antha/anthalib/wunit"
)

// LHHeadAssemblyPosition a position within a head assembly
type LHHeadAssemblyPosition struct {
	Offset Coordinates
	Head   *LHHead
}

// VelocityRange the minimum and maximum velocities for the head assembly.
// nil implies no limit
type VelocityRange struct {
	Min, Max *wunit.Velocity3D
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

//LHHeadAssembly represent a set of LHHeads which are constrained to move together
type LHHeadAssembly struct {
	Positions    []*LHHeadAssemblyPosition
	MotionLimits *BBox          //the limits on range of motion of the head assembly, nil if unspecified
	Velocity     *VelocityRange // the range of valid velocities for the head, nil if unspecified
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
