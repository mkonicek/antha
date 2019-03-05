// liquidhandling/lhtypes.Go: Part of the Antha language
// Copyright (C) 2014 the Antha authors. All rights reserved.
//
// This program is free software; you can redistribute it and/or
// modify it under the terms of the GNU General Public License
// as published by the Free Software Foundation; either version 2
// of the License, or (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program; if not, write to the Free Software
// Foundation, Inc., 51 Franklin Street, Fifth Floor, Boston, MA  02110-1301, USA.
//
// For more information relating to the software or licensing issues please
// contact license@antha-lang.Org or write to the Antha team c/o
// Synthace Ltd. The London Bioscience Innovation Centre
// 2 Royal College St, London NW1 0NH UK

package wtype

import (
	"fmt"
	"github.com/pkg/errors"

	"github.com/antha-lang/antha/antha/anthalib/wunit"
)

type ChannelOrientation bool

const (
	LHVChannel ChannelOrientation = iota%2 == 1 // vertical orientation
	LHHChannel                                  // horizontal orientation
)

func (o ChannelOrientation) String() string {
	if o == LHVChannel {
		return "vertical"
	}
	return "horizontal"
}

// what constraints apply to adjacent channels
type LHMultiChannelConstraint struct {
	X int
	Y int
	M int
}

func (lhmcc LHMultiChannelConstraint) Satisfied(wc1, wc2 WellCoords) bool {
	// this is ordered, it is assumed wc1 > wc2
	x := wc1.X - wc2.X
	y := wc1.Y - wc2.Y
	return x == lhmcc.X && y == lhmcc.Y
}

func (lhmcc LHMultiChannelConstraint) SatisfiedV(awc1, awc2 []WellCoords) bool {
	// check we have fewer than the maximum
	if len(awc1) != len(awc2) || len(awc1) > lhmcc.M {
		return false
	}

	// we assume the sets are ordered
	for i, wc1 := range awc1 {
		wc2 := awc2[i]

		if !lhmcc.Satisfied(wc1, wc2) {
			return false
		}
	}

	return true
}

func (lhmcc LHMultiChannelConstraint) Equals(lhmcc2 LHMultiChannelConstraint) bool {
	return lhmcc.X == lhmcc2.X && lhmcc.Y == lhmcc2.Y && lhmcc.M == lhmcc2.M
}

// describes sets of parameters which can be used to create a configuration
type LHChannelParameter struct {
	ID          string
	Platform    string
	Name        string
	Minvol      wunit.Volume
	Maxvol      wunit.Volume
	Minspd      wunit.FlowRate
	Maxspd      wunit.FlowRate
	Multi       int
	Independent bool
	Orientation ChannelOrientation
	Head        int
}

func (lhcprm *LHChannelParameter) Equals(prm2 *LHChannelParameter) bool {
	return lhcprm.ID == prm2.ID
}

// can you move this much? If oneshot is true it's strictly Minvol <= v <= Maxvol
// otherwise it's just Minvol <= v
func (lhcp LHChannelParameter) CanMove(v wunit.Volume, oneshot bool) bool {
	if v.LessThan(lhcp.Minvol) || (v.GreaterThan(lhcp.Maxvol) && oneshot) {
		return false
	}

	return true
}

func (lhcp LHChannelParameter) VolumeLimitString() string {
	return fmt.Sprintf("Min: %s Max: %s", lhcp.Minvol.ToString(), lhcp.Maxvol.ToString())
}

func (lhcp LHChannelParameter) String() string {
	return fmt.Sprintf("%s %s Minvol %s Maxvol %s Minspd %s Maxspd %s Multi %d Independent %t Ori %v Head %d", lhcp.Platform, lhcp.Name, lhcp.Minvol.ToString(), lhcp.Maxvol.ToString(), lhcp.Minspd.ToString(), lhcp.Maxspd.ToString(), lhcp.Multi, lhcp.Independent, lhcp.Orientation, lhcp.Head)
}

// given the dimension of the plate, what is the constraint
// on multichannel access?
func (lhcp LHChannelParameter) GetConstraint(n int) LHMultiChannelConstraint {
	// this is initially quite simple, may get more complicated over time
	// as it stands this cannot be entirely fully specified but for most of
	// the cases we can deal with it's not an issue

	if lhcp.Multi == 1 {
		return LHMultiChannelConstraint{0, 0, 1}
	}

	pitch := lhcp.Multi / n
	max := lhcp.Multi
	var x, y int

	if lhcp.Orientation == LHVChannel {
		x = 0
		y = pitch
	} else {
		x = pitch
		y = 0
	}

	return LHMultiChannelConstraint{x, y, max}
}

func (lhcp *LHChannelParameter) Dup() *LHChannelParameter {
	return lhcp.dup(false)
}

func (lhcp *LHChannelParameter) DupKeepIDs() *LHChannelParameter {
	return lhcp.dup(true)
}

func (lhcp *LHChannelParameter) dup(keepIDs bool) *LHChannelParameter {
	if lhcp == nil {
		return nil
	}
	r := NewLHChannelParameter(lhcp.Name, lhcp.Platform, lhcp.Minvol, lhcp.Maxvol, lhcp.Minspd, lhcp.Maxspd, lhcp.Multi, lhcp.Independent, lhcp.Orientation, lhcp.Head)
	if keepIDs {
		r.ID = lhcp.ID
	}

	return r
}

func NewLHChannelParameter(name, platform string, minvol, maxvol wunit.Volume, minspd, maxspd wunit.FlowRate, multi int, independent bool, orientation ChannelOrientation, head int) *LHChannelParameter {
	var lhp LHChannelParameter
	lhp.ID = GetUUID()
	lhp.Name = name
	lhp.Platform = platform
	lhp.Minvol = minvol
	lhp.Maxvol = maxvol
	lhp.Minspd = minspd
	lhp.Maxspd = maxspd
	lhp.Multi = multi
	lhp.Independent = independent
	lhp.Orientation = orientation
	lhp.Head = head
	return &lhp
}

func (lhcp *LHChannelParameter) MergeWithTip(tip *LHTip) *LHChannelParameter {
	lhcp2 := *lhcp
	if tip.MinVol.GreaterThanRounded(lhcp2.Minvol, 1) {
		lhcp2.Minvol = wunit.CopyVolume(tip.MinVol)
	}

	if tip.MaxVol.LessThanRounded(lhcp2.Maxvol, 1) {
		lhcp2.Maxvol = wunit.CopyVolume(tip.MaxVol)
	}

	return &lhcp2
}

// defines an addendum to a liquid handler
// not much to say yet

type LHDevice struct {
	ID   string
	Name string
	Mnfr string
}

func NewLHDevice(name, mfr string) *LHDevice {
	var dev LHDevice
	dev.ID = GetUUID()
	dev.Name = name
	dev.Mnfr = mfr
	return &dev
}

func (lhd *LHDevice) Dup() *LHDevice {
	d := NewLHDevice(lhd.Name, lhd.Mnfr)
	return d
}

// describes a position on the liquid handling deck
type LHPosition struct {
	Name     string        // human readable name of the position chosen by device driver
	Location Coordinates3D // absolute position of read left corner of the position
	Size     Coordinates2D // size of the position - equal to the footprint of objects which can be accepted
}

// NewLHPosition constructs a new position on a liquidhandling deck
func NewLHPosition(name string, location Coordinates3D, size Coordinates2D) *LHPosition {
	return &LHPosition{
		Name:     name,
		Location: location,
		Size:     size,
	}
}

// structure describing a solution: a combination of liquid components
// deprecated and no longer used... may well need to be deleted
type LHSolution struct {
	ID               string
	BlockID          BlockID
	Inst             string
	SName            string
	Order            int
	Components       []*Liquid
	ContainerType    string
	Welladdress      string
	Plateaddress     string
	PlateID          string
	Platetype        string
	Vol              float64 // in S.I units only for now
	Type             string
	Conc             float64
	Tvol             float64
	Majorlayoutgroup int
	Minorlayoutgroup int
}

func NewLHSolution() *LHSolution {
	var lhs LHSolution
	lhs.ID = GetUUID()
	lhs.Majorlayoutgroup = -1
	lhs.Minorlayoutgroup = -1
	return &lhs
}

func (sol LHSolution) GetComponentVolume(key string) float64 {
	vol := 0.0

	for _, v := range sol.Components {
		if v.CName == key {
			vol += v.Vol
		}
	}

	return vol
}

func (sol LHSolution) String() string {
	one := fmt.Sprintf(
		"%s, %s, %s, %s, %d",
		sol.ID,
		sol.BlockID,
		sol.Inst,
		sol.SName,
		sol.Order,
	)
	for _, c := range sol.Components {
		one = one + fmt.Sprintf("[%s], ", c.CName)
	}
	two := fmt.Sprintf("%s, %s, %s, %g, %s, %g, %g, %d, %d",
		sol.ContainerType,
		sol.Welladdress,
		sol.Platetype,
		sol.Vol,
		sol.Type,
		sol.Conc,
		sol.Tvol,
		sol.Majorlayoutgroup,
		sol.Minorlayoutgroup,
	)
	return one + two
}

func (lhs *LHSolution) GetAssignment() string {
	return lhs.Plateaddress + ":" + lhs.Welladdress
}

func New_Solution() *LHSolution {
	var solution LHSolution
	solution.ID = GetUUID()
	solution.Components = make([]*Liquid, 0, 4)
	return &solution
}

type SequentialTipLoadingBehaviour int

const (
	//NoSequentialTipLoading tips are loaded all at once, an error is raised if not possible
	NoSequentialTipLoading SequentialTipLoadingBehaviour = iota
	//ForwardSequentialTipLoading chunks of contiguous tips are loaded sequentially in the order encountered
	ForwardSequentialTipLoading
	//ReverseSequentialTipLoading chunks of contiguous tips are loaded sequentially in reverse order
	ReverseSequentialTipLoading
)

var sequentialTipLoadingBehaviourNames = map[SequentialTipLoadingBehaviour]string{
	NoSequentialTipLoading:      "no sequential tip loading",
	ForwardSequentialTipLoading: "forward sequential tip loading",
	ReverseSequentialTipLoading: "reverse sequential tip loading",
}

func (s SequentialTipLoadingBehaviour) String() string {
	return sequentialTipLoadingBehaviourNames[s]
}

//TipLoadingBehaviour describe the way in which tips are loaded
type TipLoadingBehaviour struct {
	//OverrideLoadTipsCommand true it the liquid handler will override which tips are loaded
	OverrideLoadTipsCommand bool
	//AutoRefillTipboxes are tipboxes automaticall refilled
	AutoRefillTipboxes bool
	//LoadingOrder are tips loaded ColumnWise or RowWise
	LoadingOrder MajorOrder
	//VerticalLoadingDirection the direction along which columns are loaded
	VerticalLoadingDirection VerticalDirection
	//HorizontalLoadingDirection the direction along which rows are loaded
	HorizontalLoadingDirection HorizontalDirection
	//ChunkingBehaviour how to load tips when the requested number aren't available contiguously
	ChunkingBehaviour SequentialTipLoadingBehaviour
}

//String get a string description for debuggin
func (s TipLoadingBehaviour) String() string {

	autoRefill := ""
	if !s.AutoRefillTipboxes {
		autoRefill = "no "
	}

	if !s.OverrideLoadTipsCommand {
		return fmt.Sprintf("%sauto-refilling, no loading override", autoRefill)
	}

	return fmt.Sprintf("%sauto-refilling, loading order: %v, %v, %v, %v", autoRefill, s.LoadingOrder, s.VerticalLoadingDirection, s.HorizontalLoadingDirection, s.ChunkingBehaviour)
}

// LHHeadAssemblyPosition a position within a head assembly
type LHHeadAssemblyPosition struct {
	Offset Coordinates3D
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
	Positions      []*LHHeadAssemblyPosition
	MotionLimits   *BBox          //the limits on range of motion of the head assembly, nil if unspecified
	VelocityLimits *VelocityRange // the range of valid velocities for the head, nil if unspecified
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
		Positions:      make([]*LHHeadAssemblyPosition, 0, len(self.Positions)),
		MotionLimits:   self.MotionLimits.Dup(),
		VelocityLimits: self.VelocityLimits.Dup(),
	}
	for _, pos := range self.Positions {
		ret.AddPosition(pos.Offset)
	}
	return ret
}

//AddPosition add a position to the head assembly with the given offset
func (self *LHHeadAssembly) AddPosition(Offset Coordinates3D) {
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
