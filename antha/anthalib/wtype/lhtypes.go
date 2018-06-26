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
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"

	"github.com/antha-lang/antha/antha/anthalib/wunit"
)

const (
	LHVChannel = iota // vertical orientation
	LHHChannel        // horizontal orientation
)

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
	Orientation int
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
	return fmt.Sprintf("%s %s Minvol %s Maxvol %s Minspd %s Maxspd %s Multi %d Independent %t Ori %d Head %d", lhcp.Platform, lhcp.Name, lhcp.Minvol.ToString(), lhcp.Maxvol.ToString(), lhcp.Minspd.ToString(), lhcp.Maxspd.ToString(), lhcp.Multi, lhcp.Independent, lhcp.Orientation, lhcp.Head)
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

func (lhcp LHChannelParameter) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		ID          string
		Name        string
		Minvol      wunit.Volume
		Maxvol      wunit.Volume
		Minspd      wunit.FlowRate
		Maxspd      wunit.FlowRate
		Multi       int
		Independent bool
		Orientation int
		Head        int
	}{
		lhcp.ID,
		lhcp.Name,
		lhcp.Minvol,
		lhcp.Maxvol,
		lhcp.Minspd,
		lhcp.Maxspd,
		lhcp.Multi,
		lhcp.Independent,
		lhcp.Orientation,
		lhcp.Head,
	})
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

func NewLHChannelParameter(name, platform string, minvol, maxvol wunit.Volume, minspd, maxspd wunit.FlowRate, multi int, independent bool, orientation int, head int) *LHChannelParameter {
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

// describes a position on the liquid handling deck and its current state
type LHPosition struct {
	ID    string
	Name  string
	Num   int
	Extra []LHDevice
	Maxh  float64
}

func NewLHPosition(position_number int, name string, maxh float64) *LHPosition {
	var lhp LHPosition
	lhp.ID = GetUUID()
	lhp.Name = name
	lhp.Num = position_number
	lhp.Extra = make([]LHDevice, 0, 2)
	lhp.Maxh = maxh
	return &lhp
}

// @implement Location
// -- this is clearly somewhere that something can be
// need to implement the liquid handler as a location as well

func (lhp *LHPosition) Location_ID() string {
	return lhp.ID
}

func (lhp *LHPosition) Location_Name() string {
	return lhp.Name
}

func (lhp *LHPosition) Container() Location {
	return lhp
}

func (lhp *LHPosition) Positions() []Location {
	return nil
}

func (lhp *LHPosition) Shape() *Shape {
	return NewShape("box", "mm", 0.08548, 0.12776, 0.0)
}

// structure describing a solution: a combination of liquid components
// deprecated and no longer used... may well need to be deleted
type LHSolution struct {
	ID               string
	BlockID          BlockID
	Inst             string
	SName            string
	Order            int
	Components       []*LHComponent
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
	solution.Components = make([]*LHComponent, 0, 4)
	return &solution
}

// head
type LHHead struct {
	Name         string
	Manufacturer string
	ID           string
	Adaptor      *LHAdaptor
	Params       *LHChannelParameter
	//TipLoading defined the behaviour of the head when loading tips
	TipLoading TipLoadingBehaviour
}

func NewLHHead(name, mf string, params *LHChannelParameter) *LHHead {
	var lhh LHHead
	lhh.Manufacturer = mf
	lhh.Name = name
	lhh.Params = params
	return &lhh
}

func (head *LHHead) Dup() *LHHead {
	return head.dup(false)
}

func (head *LHHead) DupKeepIDs() *LHHead {
	return head.dup(true)
}

func (head *LHHead) dup(keepIDs bool) *LHHead {
	if head == nil {
		return nil
	}
	var params *LHChannelParameter
	var adaptor *LHAdaptor
	if keepIDs {
		params = head.Params.DupKeepIDs()
		adaptor = head.Adaptor.DupKeepIDs()
	} else {
		params = head.Params.Dup()
		adaptor = head.Adaptor.Dup()
	}
	h := NewLHHead(head.Name, head.Manufacturer, params)
	h.Adaptor = adaptor
	h.TipLoading = head.TipLoading
	return h
}

func (head *LHHead) Equal(rhs *LHHead) bool {
	return head.Name == rhs.Name
}

func (lhh *LHHead) GetParams() *LHChannelParameter {
	if lhh.Adaptor == nil {
		return lhh.Params
	} else {
		return lhh.Adaptor.GetParams()
	}
}

// adaptor
type LHAdaptor struct {
	Name         string
	ID           string
	Manufacturer string
	Params       *LHChannelParameter
	Tips         []*LHTip
}

func NewLHAdaptor(name, mf string, params *LHChannelParameter) *LHAdaptor {
	var lha LHAdaptor
	lha.Name = name
	lha.Manufacturer = mf
	lha.Params = params
	lha.Tips = make([]*LHTip, params.Multi)
	return &lha
}

func (lha *LHAdaptor) Dup() *LHAdaptor {
	return lha.dup(false)
}

func (lha *LHAdaptor) DupKeepIDs() *LHAdaptor {
	return lha.dup(true)
}

func (lha *LHAdaptor) dup(keepIDs bool) *LHAdaptor {
	var params *LHChannelParameter
	if keepIDs {
		params = lha.Params.DupKeepIDs()
	} else {
		params = lha.Params.Dup()
	}

	ad := NewLHAdaptor(lha.Name, lha.Manufacturer, params)

	for i, tip := range lha.Tips {
		if tip != nil {
			if keepIDs {
				ad.AddTip(i, tip.DupKeepID())
			} else {
				ad.AddTip(i, tip.Dup())
			}
		}
	}

	return ad
}

//The number of tips currently loaded
func (lha *LHAdaptor) NTipsLoaded() int {
	r := 0
	for i := range lha.Tips {
		if lha.Tips[i] != nil {
			r += 1
		}
	}
	return r
}

//Is there a tip loaded on channel_number
func (lha *LHAdaptor) IsTipLoaded(channel_number int) bool {
	return lha.Tips[channel_number] != nil
}

//Return the tip at channel_number, nil otherwise
func (lha *LHAdaptor) GetTip(channel_number int) *LHTip {
	return lha.Tips[channel_number]
}

//Load a tip to the specified channel
func (lha *LHAdaptor) AddTip(channel_number int, tip *LHTip) {
	lha.Tips[channel_number] = tip
}

//Remove a tip from the specified channel and return it
func (lha *LHAdaptor) RemoveTip(channel_number int) *LHTip {
	tip := lha.Tips[channel_number]
	lha.Tips[channel_number] = nil
	return tip
}

//Remove every tip from the adaptor
func (lha *LHAdaptor) RemoveTips() []*LHTip {
	ret := make([]*LHTip, 0, lha.NTipsLoaded())
	for i := range lha.Tips {
		if lha.Tips[i] != nil {
			ret = append(ret, lha.Tips[i])
			lha.Tips[i] = nil
		}
	}
	return ret
}

func (lha *LHAdaptor) GetParams() *LHChannelParameter {
	if lha.NTipsLoaded() == 0 {
		return lha.Params
	} else {
		params := *lha.Params
		for _, tip := range lha.Tips {
			if tip != nil {
				params = *params.MergeWithTip(tip)
			}
		}
		return &params
	}
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

//LHHeadAssemblyPosition a position within a head assembly
type LHHeadAssemblyPosition struct {
	Offset Coordinates
	Head   *LHHead
}

//LHHeadAssembly represent a set of LHHeads which are constrained to move together
type LHHeadAssembly struct {
	Positions    []*LHHeadAssemblyPosition
	MotionLimits *BBox
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
