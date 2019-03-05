// /anthalib/simulator/liquidhandling/robotstate.go: Part of the Antha language
// Copyright (C) 2015 The Antha authors. All rights reserved.
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
// contact license@antha-lang.org or write to the Antha team c/o
// Synthace Ltd. The London Bioscience Innovation Centre
// 2 Royal College St, London NW1 0NH UK

package liquidhandling

import (
	"fmt"
	"github.com/pkg/errors"
	"sort"
	"strings"

	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
)

// -------------------------------------------------------------------------------
//                            ChannelState
// -------------------------------------------------------------------------------

//ChannelState Represent the physical state of a single channel
type ChannelState struct {
	number   int
	tip      *wtype.LHTip        //Nil if no tip loaded, otherwise the tip that's loaded
	contents *wtype.Liquid       //What's in the tip?
	position wtype.Coordinates3D //position relative to the adaptor
	adaptor  *AdaptorState       //the channel's adaptor
	radius   float64
}

func NewChannelState(number int, adaptor *AdaptorState, position wtype.Coordinates3D, radius float64) *ChannelState {
	r := ChannelState{}
	r.number = number
	r.position = position
	r.adaptor = adaptor
	r.radius = radius

	return &r
}

//                            Accessors
//                            ---------

//HasTip is a tip loaded
func (self *ChannelState) HasTip() bool {
	return self.tip != nil
}

//GetTip get the loaded tip, returns nil if none loaded
func (self *ChannelState) GetTip() *wtype.LHTip {
	return self.tip
}

//IsEmpty returns true only if a tip is loaded and contains liquid
func (self *ChannelState) IsEmpty() bool {
	return self.HasTip() && self.contents != nil && self.contents.IsZero()
}

//GetContents get the contents of the loaded tip, retuns nil if no contents or no tip
func (self *ChannelState) GetContents() *wtype.Liquid {
	return self.contents
}

func (self *ChannelState) GetRadius() float64 {
	return self.radius
}

func (self *ChannelState) GetBounds(channelClearance float64) wtype.BBox {
	var r float64
	h := self.tip.GetEffectiveHeight() + channelClearance
	if !self.HasTip() {
		r = self.radius
	}

	ret := wtype.NewBBox(
		self.GetAbsolutePosition().Subtract(wtype.Coordinates3D{X: r, Y: r, Z: h}),
		wtype.Coordinates3D{X: 2.0 * r, Y: 2.0 * r, Z: h})

	return *ret
}

//GetRelativePosition get the channel's position relative to the head
func (self *ChannelState) GetRelativePosition() wtype.Coordinates3D {
	return self.position
}

//SetRelativePosition get the channel's position relative to the head
func (self *ChannelState) SetRelativePosition(v wtype.Coordinates3D) {
	self.position = v
}

//GetAbsolutePosition get the channel's absolute position
func (self *ChannelState) GetAbsolutePosition() wtype.Coordinates3D {
	return self.position.Add(self.adaptor.GetPosition())
}

//GetTarget get the LHObject below the adaptor
func (self *ChannelState) GetTarget() wtype.LHObject {
	return self.adaptor.GetGroup().GetRobot().GetDeck().GetChildBelow(self.GetAbsolutePosition())
}

//                            Actions
//                            -------

//Aspirate
func (self *ChannelState) Aspirate(volume wunit.Volume) error {

	return nil
}

//Dispense
func (self *ChannelState) Dispense(volume *wunit.Volume) error {

	return nil
}

//LoadTip
func (self *ChannelState) LoadTip(tip *wtype.LHTip) {
	self.tip = tip
}

//UnloadTip
func (self *ChannelState) UnloadTip() *wtype.LHTip {
	tip := self.tip
	self.tip = nil
	return tip
}

//GetCollisions get collisions with this channel. channelClearance defined a height below the channel/tip to include
func (self *ChannelState) GetCollisions(settings *SimulatorSettings, channelClearance float64) []wtype.LHObject {
	deck := self.adaptor.GetGroup().GetRobot().GetDeck()

	var ret []wtype.LHObject

	box := self.GetBounds(channelClearance)
	objects := deck.GetBoxIntersections(box)

	//tips are allowed in wells
	ret = make([]wtype.LHObject, 0, len(objects))
	if self.HasTip() {
		tipBottom := box.GetPosition()
		for _, obj := range objects {
			//ignore tipboxes if told to
			if _, ok := obj.(*wtype.LHTipbox); !settings.IsTipboxCollisionEnabled() && ok {
				continue
			}
			//don't add wells, instead add the plate if the tip bottom has hit the plate
			if well, ok := obj.(*wtype.LHWell); ok {
				if len(well.GetPointIntersections(tipBottom)) == 0 {
					ret = append(ret, well.GetParent())
				}
			} else if _, ok := obj.(*wtype.Plate); !ok {
				ret = append(ret, obj)
			}
		}
	} else {
		for _, obj := range objects {
			if well, ok := obj.(*wtype.LHWell); ok {
				//wells are empty space, so we can't collide with them, but we could collide with the parent
				//if we aren't fully in the well
				if !coneInWell(box, well) {
					ret = append(ret, well.GetParent())
				}
			} else {
				ret = append(ret, obj)
			}
		}
	}

	return ret
}

//coneInWell is the cone contained in the well
func coneInWell(cone wtype.BBox, well *wtype.LHWell) bool {
	wellStart := well.GetPosition()
	wellEnd := wellStart.Add(well.GetSize())

	coneStart := cone.GetPosition()
	coneEnd := cone.GetPosition().Add(cone.GetSize())

	if coneStart.X < wellStart.X || coneEnd.X > wellEnd.X || coneStart.Y < wellStart.Y || coneEnd.Y > wellEnd.Y {
		return false
	}

	return !(coneStart.Z < (wellStart.Z + well.Bottomh))
}

// -------------------------------------------------------------------------------
//                            AdaptorState
// -------------------------------------------------------------------------------

//AdaptorState Represent the physical state and layout of the adaptor
type AdaptorState struct {
	name         string
	channels     []*ChannelState
	offset       wtype.Coordinates3D
	independent  bool
	params       *wtype.LHChannelParameter
	group        *AdaptorGroup
	tipBehaviour wtype.TipLoadingBehaviour
	index        int
}

func NewAdaptorState(name string,
	independent bool,
	channels int,
	channel_offset wtype.Coordinates3D, coneRadius float64,
	params *wtype.LHChannelParameter,
	tipBehaviour wtype.TipLoadingBehaviour) *AdaptorState {
	as := AdaptorState{
		name,
		make([]*ChannelState, 0, channels),
		wtype.Coordinates3D{},
		independent,
		params.Dup(),
		nil,
		tipBehaviour,
		-1,
	}

	for i := 0; i < channels; i++ {
		as.channels = append(as.channels, NewChannelState(i, &as, channel_offset.Multiply(float64(i)), coneRadius))
	}

	return &as
}

//String summarize the state of the adaptor
func (self *AdaptorState) SummariseTips() string {
	tipTypes := make([]string, 0, len(self.channels))
	tipLengths := make([]string, 0, len(self.channels))
	loaded := make([]int, 0, len(self.channels))

	for i, channel := range self.channels {
		if channel.HasTip() {
			loaded = append(loaded, i)
			tipTypes = append(tipTypes, channel.GetTip().GetType())
			tipLengths = append(tipLengths, wunit.NewLength(channel.GetTip().GetEffectiveHeight(), "mm").String())
		}
	}

	tipTypes = getUnique(tipTypes, true)
	tipLengths = getUnique(tipLengths, true)

	if len(loaded) > 0 {
		return fmt.Sprintf("%s: %s %s (effective %s %s) loaded on %s",
			self.name,
			pTips(len(loaded)),
			strings.Join(tipTypes, ","),
			pLengths(len(tipLengths)),
			strings.Join(tipLengths, ","),
			summariseChannels(loaded))
	}
	return fmt.Sprintf("%s: <no tips>", self.name)
}

//SummarisePositions get a human readable description of the X,Y,Z position of
//the adaptor, as well as which type of tips are loaded on which channels.
//channelsColliding provides an optional list of channel indexes which
//have been detected as being involved in a collision, which will also be
//displayed in the output.
func (self *AdaptorState) SummarisePositions(channelsColliding []int) string {

	positions := make([]wtype.Coordinates3D, len(self.channels))
	for i, channel := range self.channels {
		positions[i] = channel.GetAbsolutePosition()
		positions[i].Z -= channel.GetTip().GetEffectiveHeight()
	}

	collisions := make([]bool, len(self.channels))
	for _, ch := range channelsColliding {
		collisions[ch] = true
	}

	var places int
	for _, pos := range positions {
		if p := fmt.Sprintf("%.1f", pos.X); len(p) > places {
			places = len(p)
		}
		if p := fmt.Sprintf("%.1f", pos.Y); len(p) > places {
			places = len(p)
		}
		if p := fmt.Sprintf("%.1f", pos.Z); len(p) > places {
			places = len(p)
		}
	}

	floatTemplate := fmt.Sprintf("%%%d.1f", places)
	intTemplate := fmt.Sprintf("%%%dd", places)

	var head, lineX, lineY, lineZ []string
	head = append(head, "channel:")
	lineX = append(lineX, "      X:")
	lineY = append(lineY, "      Y:")
	lineZ = append(lineZ, "      Z:")
	lineState := make([]string, len(positions)+1)
	lineState[0] = "  state:"
	for i, pos := range positions {
		head = append(head, fmt.Sprintf(intTemplate, i))
		lineX = append(lineX, fmt.Sprintf(floatTemplate, pos.X))
		lineY = append(lineY, fmt.Sprintf(floatTemplate, pos.Y))
		lineZ = append(lineZ, fmt.Sprintf(floatTemplate, pos.Z))

		if collisions[i] {
			lineState[i+1] += "C"
		}
		if self.channels[i].HasTip() {
			lineState[i+1] += "T"
		}
	}

	for i := 1; i < len(lineState); i++ {
		lineState[i] = strings.Repeat(" ", places-len(lineState[i])) + lineState[i]
	}

	return strings.Join(
		[]string{
			strings.Join(head, " "),
			strings.Join(lineX, " "),
			strings.Join(lineY, " "),
			strings.Join(lineZ, " "),
			strings.Join(lineState, " "),
		}, "\n")
}

//                            Accessors
//                            ---------

//GetName
func (self *AdaptorState) GetName() string {
	return self.name
}

//GetPosition
func (self *AdaptorState) GetPosition() wtype.Coordinates3D {
	return self.offset.Add(self.group.GetPosition())
}

//GetChannelCount
func (self *AdaptorState) GetChannelCount() int {
	return len(self.channels)
}

//GetChannel
func (self *AdaptorState) GetChannel(ch int) *ChannelState {
	return self.channels[ch]
}

//GetParamsForChannel
func (self *AdaptorState) GetParamsForChannel(ch int) *wtype.LHChannelParameter {
	if tip := self.GetChannel(ch).GetTip(); tip != nil {
		return self.params.MergeWithTip(tip)
	}
	return self.params
}

//GetTipCount
func (self *AdaptorState) GetTipCount() int {
	r := 0
	for _, ch := range self.channels {
		if ch.HasTip() {
			r++
		}
	}
	return r
}

//IsIndependent
func (self *AdaptorState) IsIndependent() bool {
	return self.independent
}

func (self *AdaptorState) GetIndex() int {
	return self.index
}

func (self *AdaptorState) setIndex(i int) {
	self.index = i
}

//GetGroup
func (self *AdaptorState) GetGroup() *AdaptorGroup {
	return self.group
}

//SetGroup
func (self *AdaptorState) SetGroup(g *AdaptorGroup) {
	self.group = g
}

func (self *AdaptorState) SetPosition(p wtype.Coordinates3D) error {
	return self.group.SetPosition(p.Subtract(self.offset))
}

func (self *AdaptorState) SetOffset(p wtype.Coordinates3D) {
	self.offset = p
}

func (self *AdaptorState) OverridesLoadTipsCommand() bool {
	return self.tipBehaviour.OverrideLoadTipsCommand
}

func (self *AdaptorState) SetOverridesLoadTipsCommand(v bool) {
	self.tipBehaviour.OverrideLoadTipsCommand = v
}

func (self *AdaptorState) AutoRefillsTipboxes() bool {
	return self.tipBehaviour.AutoRefillTipboxes
}

func isVAligned(lhs wtype.WellCoords, rhs wtype.WellCoords) bool {
	return lhs.X == rhs.X
}

func isHAligned(lhs wtype.WellCoords, rhs wtype.WellCoords) bool {
	return lhs.Y == rhs.Y
}

//GetTipsToLoad get which tips would be loaded by the adaptor given the tiploading behaviour
//returns an error if OverridesLoadTipsCommand is false or there aren't enough tips
func (self *AdaptorState) GetTipCoordsToLoad(tb *wtype.LHTipbox, num int) ([][]wtype.WellCoords, error) {
	var ret [][]wtype.WellCoords
	if !self.tipBehaviour.OverrideLoadTipsCommand {
		return ret, errors.New("Tried to get tips when override is false")
	}

	it := wtype.NewAddressIterator(tb,
		self.tipBehaviour.LoadingOrder,
		self.tipBehaviour.VerticalLoadingDirection,
		self.tipBehaviour.HorizontalLoadingDirection,
		false)

	isInline := isVAligned
	if self.params.Orientation == wtype.LHHChannel {
		isInline = isHAligned
	}

	tipsRemaining := num
	var lastTipCoord wtype.WellCoords
	currChunk := make([]wtype.WellCoords, 0, num)
	for wc := it.Curr(); it.Valid(); wc = it.Next() {
		//start a new chunk if this chunk has something in it AND (we found an empty position OR we changed row/column)
		if len(currChunk) > 0 && (!tb.HasTipAt(wc) || !isInline(lastTipCoord, wc)) {
			//keep the chunk if either this chunk provides all the tips we need or we can load it sequentially
			if !(self.tipBehaviour.ChunkingBehaviour == wtype.NoSequentialTipLoading && len(currChunk) < tipsRemaining) {
				ret = append(ret, currChunk)
				tipsRemaining -= len(currChunk)
			}
			currChunk = make([]wtype.WellCoords, 0, tipsRemaining)
		}
		//if we have all the chunks we need
		if len(currChunk) >= tipsRemaining {
			break
		}
		//add the next tip
		if tb.HasTipAt(wc) {
			currChunk = append(currChunk, wc)
			lastTipCoord = wc
		}
	}
	if len(currChunk) > 0 {
		ret = append(ret, currChunk)
		tipsRemaining -= len(currChunk)
	}

	if self.tipBehaviour.ChunkingBehaviour == wtype.ReverseSequentialTipLoading {
		//apparently this is actually the recommended way to reverse a list in place
		for i := len(ret)/2 - 1; i >= 0; i-- {
			opp := len(ret) - 1 - i
			ret[i], ret[opp] = ret[opp], ret[i]
		}

		for _, chunk := range ret {
			for i := len(chunk)/2 - 1; i >= 0; i-- {
				opp := len(chunk) - 1 - i
				chunk[i], chunk[opp] = chunk[opp], chunk[i]
			}
		}
	}

	if tipsRemaining > 0 {
		return ret, errors.New("not enough tips in tipbox")
	}

	return ret, nil
}

// -------------------------------------------------------------------------------
//                            AdaptorGroup
// -------------------------------------------------------------------------------

// AdaptorGroup simulate a set of adaptors which are physically attached
type AdaptorGroup struct {
	adaptors      []*AdaptorState
	offsets       []wtype.Coordinates3D
	motionLimits  *wtype.BBox
	velocity      *wunit.Velocity3D
	velocityRange *wtype.VelocityRange
	position      wtype.Coordinates3D
	robot         *RobotState
}

// NewAdaptorGroup convert a HeadAssembly into an AdaptorGroup for simulation
func NewAdaptorGroup(assembly *wtype.LHHeadAssembly) *AdaptorGroup {

	offsets := make([]wtype.Coordinates3D, len(assembly.Positions))
	for i, pos := range assembly.Positions {
		offsets[i] = pos.Offset
	}

	group := &AdaptorGroup{
		adaptors:      make([]*AdaptorState, len(offsets)),
		offsets:       offsets,
		motionLimits:  assembly.MotionLimits.Dup(),
		velocity:      &wunit.Velocity3D{},
		velocityRange: assembly.VelocityLimits.Dup(),
	}

	for i, pos := range assembly.Positions {
		if pos.Head == nil { //ignore assembly position which have nothing loaded
			continue
		}
		p := pos.Head.Adaptor.Params
		//9mm spacing currently hardcoded.
		//At some point we'll either need to fetch this from the driver or
		//infer it from the type of tipboxes/plates accepted
		spacing := wtype.Coordinates3D{X: 0, Y: 0, Z: 0}
		if p.Orientation == wtype.LHVChannel {
			spacing.Y = 9.
		} else if p.Orientation == wtype.LHHChannel {
			spacing.X = 9.
		}
		adaptor := NewAdaptorState(pos.Head.Adaptor.Name, p.Independent, p.Multi, spacing, coneRadius, p, pos.Head.TipLoading)
		group.LoadAdaptor(i, adaptor)
	}

	return group
}

//GetAdaptor get an adaptor state
func (self *AdaptorGroup) GetAdaptor(i int) (*AdaptorState, error) {
	if i < 0 || i >= len(self.adaptors) {
		return nil, errors.Errorf("unknown head %d", i)
	}
	return self.adaptors[i], nil
}

func (self *AdaptorGroup) GetAdaptors() []*AdaptorState {
	return self.adaptors
}

//CountAdaptors count the adaptors
func (self *AdaptorGroup) NumAdaptors() int {
	return len(self.adaptors)
}

func (self *AdaptorGroup) LoadAdaptor(pos int, adaptor *AdaptorState) {
	self.adaptors[pos] = adaptor
	adaptor.SetGroup(self)
	adaptor.SetOffset(self.offsets[pos])
	if self.robot != nil {
		self.robot.updateAdaptorIndexes()
	}
}

func (self *AdaptorGroup) GetPosition() wtype.Coordinates3D {
	return self.position
}

func oneDP(v float64) string {
	ret := fmt.Sprintf("%.1f", v)
	if ret[len(ret)-2:] == ".0" {
		return ret[:len(ret)-2]
	}
	return ret
}

func (self *AdaptorGroup) SetPosition(p wtype.Coordinates3D) error {
	self.position = p
	if self.motionLimits != nil && !self.motionLimits.Contains(p) {
		template := "%smm too %s"
		rearranging := "rearranging the deck"
		var failures []string
		var advice []string

		start := self.motionLimits.GetPosition()
		extent := self.motionLimits.GetPosition().Add(self.motionLimits.GetSize())

		if p.X < start.X {
			failures = append(failures, fmt.Sprintf(template, oneDP(start.X-p.X), "far left"))
			advice = append(advice, rearranging)
		} else if p.X > extent.X {
			failures = append(failures, fmt.Sprintf(template, oneDP(p.X-extent.X), "far right"))
			advice = append(advice, rearranging)
		}

		if p.Y < start.Y {
			failures = append(failures, fmt.Sprintf(template, oneDP(start.Y-p.Y), "far backwards"))
			advice = append(advice, rearranging)
		} else if p.Y > extent.Y {
			failures = append(failures, fmt.Sprintf(template, oneDP(p.Y-extent.Y), "far forwards"))
			advice = append(advice, rearranging)
		}

		if p.Z < start.Z {
			failures = append(failures, fmt.Sprintf(template, oneDP(start.Z-p.Z), "low"))
			advice = append(advice, "adding a riser to the object on the deck")
		} else if p.Z > extent.Z {
			failures = append(failures, fmt.Sprintf(template, oneDP(p.Z-extent.Z), "high"))
			advice = append(advice, "lowering the object on the deck")
		}

		return errors.Errorf("head cannot reach position: position is %s, please try %s",
			strings.Join(failures, " and "),
			strings.Join(getUnique(advice, true), " and "))
	}
	return nil
}

func (self *AdaptorGroup) GetRobot() *RobotState {
	return self.robot
}

func (self *AdaptorGroup) SetRobot(r *RobotState) {
	self.robot = r
}

func (self *AdaptorGroup) SetDriveSpeed(a wunit.Axis, v wunit.Velocity) error {
	if self.velocityRange != nil {
		if min, max := self.velocityRange.Min.GetAxis(a), self.velocityRange.Max.GetAxis(a); v.LessThan(min) || v.GreaterThan(max) {
			return errors.Errorf("%v is outside allowable range [%v - %v]", v, min, max)
		}
	} else {
		if !v.IsPositive() {
			return errors.Errorf("speed must be positive")
		}
	}
	self.velocity.SetAxis(a, v)
	return nil
}

// -------------------------------------------------------------------------------
//                            RobotState
// -------------------------------------------------------------------------------

//RobotState Represent the physical state of a liquidhandling robot
type RobotState struct {
	deck          *wtype.LHDeck
	adaptorGroups []*AdaptorGroup
	initialized   bool
	finalized     bool
}

func NewRobotState() *RobotState {
	rs := RobotState{
		nil,
		make([]*AdaptorGroup, 0),
		false,
		false,
	}
	return &rs
}

func (self *RobotState) SummariseState(channelsColliding map[int][]int) string {
	return fmt.Sprintf("Adaptors: key = {T: Tip loaded, C: involved in Collision}\n\t%s\nDeck:\n\t%s",
		strings.Replace(self.SummariseAdaptors(channelsColliding), "\n", "\n\t", -1),
		strings.Replace(self.SummariseDeck(), "\n", "\n\t", -1))
}

func (self *RobotState) SummariseAdaptors(channelsColliding map[int][]int) string {
	var lines []string
	for i, group := range self.adaptorGroups {
		lines = append(lines, fmt.Sprintf("AdaptorGroup %d: position %v", i, group.GetPosition()))
	}

	for i, adaptor := range self.GetAdaptors() {
		lines = append(lines, fmt.Sprintf("Head %d: %s\n\t%s", i, adaptor.SummariseTips(), strings.Replace(adaptor.SummarisePositions(channelsColliding[i]), "\n", "\n\t", -1)))
	}

	return strings.Join(lines, "\n")
}

func (self *RobotState) SummariseDeck() string {
	slotNames := self.deck.GetSlotNames()
	sort.Strings(slotNames)

	var nameLen int
	for _, name := range slotNames {
		if len(name) > nameLen {
			nameLen = len(name)
		}
	}

	var ret []string
	for _, name := range slotNames {
		object, _ := self.deck.GetChild(name)
		paddedName := strings.Repeat(" ", nameLen-len(name)) + name
		if object != nil {
			ret = append(ret, fmt.Sprintf("%s: %s \"%s\" of type %s with size %v",
				paddedName, wtype.ClassOf(object), wtype.NameOf(object), wtype.TypeOf(object), object.GetSize()))
			ret = append(ret, fmt.Sprintf("%s  bounds: [%v - %v]",
				strings.Repeat(" ", nameLen), object.GetPosition(), object.GetPosition().Add(object.GetSize())))
			if plate, ok := object.(*wtype.LHPlate); ok {
				ret = append(ret, fmt.Sprintf("%s  well geometry: (zStart, bottomH, height, radius) = (%f, %f, %f, %f)\n",
					strings.Repeat(" ", nameLen), plate.WellZStart, plate.Welltype.Bottomh, plate.Welltype.GetSize().Z, plate.Welltype.GetSize().X))
			}
		} else {
			ret = append(ret, fmt.Sprintf("%s: <empty>", paddedName))
		}
	}

	return strings.Join(ret, "\n")
}

//                            Accessors
//                            ---------

//GetAdaptorGroup
func (self *RobotState) GetAdaptorGroup(num int) (*AdaptorGroup, error) {
	if num < 0 || num >= len(self.adaptorGroups) {
		return nil, errors.Errorf("unknown head assembly %d", num)
	}
	return self.adaptorGroups[num], nil
}

func (self *RobotState) GetAdaptor(groupIndex int, adaptorIndex int) (*AdaptorState, error) {
	group, err := self.GetAdaptorGroup(groupIndex)
	if err != nil {
		return nil, err
	}
	adaptor, err := group.GetAdaptor(adaptorIndex)
	if err != nil {
		return nil, errors.Wrapf(err, "head assembly %d", groupIndex)
	}
	return adaptor, nil
}

func (self *RobotState) NumAdaptors() int {
	r := 0
	for _, ag := range self.adaptorGroups {
		r += ag.NumAdaptors()
	}
	return r
}

func (self *RobotState) GetAdaptors() []*AdaptorState {
	r := make([]*AdaptorState, 0, self.NumAdaptors())
	for _, ag := range self.adaptorGroups {
		r = append(r, ag.GetAdaptors()...)
	}

	return r
}

//GetNumberOfAdaptorGroups
func (self *RobotState) NumAdaptorGroups() int {
	return len(self.adaptorGroups)
}

//AddAdaptorGroup
func (self *RobotState) AddAdaptorGroup(a *AdaptorGroup) {
	a.SetRobot(self)
	self.adaptorGroups = append(self.adaptorGroups, a)
	self.updateAdaptorIndexes()
}

func (self *RobotState) updateAdaptorIndexes() {
	for i, ad := range self.GetAdaptors() {
		ad.setIndex(i)
	}
}

//GetDeck
func (self *RobotState) GetDeck() *wtype.LHDeck {
	return self.deck
}

//SetDeck
func (self *RobotState) SetDeck(deck *wtype.LHDeck) {
	self.deck = deck
}

//IsInitialized
func (self *RobotState) IsInitialized() bool {
	return self.initialized
}

//IsFinalized
func (self *RobotState) IsFinalized() bool {
	return self.finalized
}

//                            Actions
//                            -------

//Initialize
func (self *RobotState) Initialize() {
	self.initialized = true
}

//Finalize
func (self *RobotState) Finalize() {
	self.finalized = true
}
