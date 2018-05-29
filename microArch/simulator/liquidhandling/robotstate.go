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
	"github.com/pkg/errors"

	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
)

// -------------------------------------------------------------------------------
//                            ChannelState
// -------------------------------------------------------------------------------

//ChannelState Represent the physical state of a single channel
type ChannelState struct {
	number   int
	tip      *wtype.LHTip       //Nil if no tip loaded, otherwise the tip that's loaded
	contents *wtype.LHComponent //What's in the tip?
	position wtype.Coordinates  //position relative to the adaptor
	adaptor  *AdaptorState      //the channel's adaptor
}

func NewChannelState(number int, adaptor *AdaptorState, position wtype.Coordinates) *ChannelState {
	r := ChannelState{}
	r.number = number
	r.position = position
	r.adaptor = adaptor

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
func (self *ChannelState) GetContents() *wtype.LHComponent {
	return self.contents
}

//GetRelativePosition get the channel's position relative to the head
func (self *ChannelState) GetRelativePosition() wtype.Coordinates {
	return self.position
}

//SetRelativePosition get the channel's position relative to the head
func (self *ChannelState) SetRelativePosition(v wtype.Coordinates) {
	self.position = v
}

//GetAbsolutePosition get the channel's absolute position
func (self *ChannelState) GetAbsolutePosition() wtype.Coordinates {
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

// -------------------------------------------------------------------------------
//                            AdaptorState
// -------------------------------------------------------------------------------

//AdaptorState Represent the physical state and layout of the adaptor
type AdaptorState struct {
	name         string
	channels     []*ChannelState
	offset       wtype.Coordinates
	independent  bool
	params       *wtype.LHChannelParameter
	group        *AdaptorGroup
	tipBehaviour wtype.TipLoadingBehaviour
}

func NewAdaptorState(name string,
	independent bool,
	channels int,
	channel_offset wtype.Coordinates,
	params *wtype.LHChannelParameter,
	tipBehaviour wtype.TipLoadingBehaviour) *AdaptorState {
	as := AdaptorState{
		name,
		make([]*ChannelState, 0, channels),
		wtype.Coordinates{},
		independent,
		params.Dup(),
		nil,
		tipBehaviour,
	}

	for i := 0; i < channels; i++ {
		as.channels = append(as.channels, NewChannelState(i, &as, channel_offset.Multiply(float64(i))))
	}

	return &as
}

//                            Accessors
//                            ---------

//GetName
func (self *AdaptorState) GetName() string {
	return self.name
}

//GetPosition
func (self *AdaptorState) GetPosition() wtype.Coordinates {
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

//GetGroup
func (self *AdaptorState) GetGroup() *AdaptorGroup {
	return self.group
}

//SetGroup
func (self *AdaptorState) SetGroup(g *AdaptorGroup) {
	self.group = g
}

func (self *AdaptorState) SetPosition(p wtype.Coordinates) error {
	return self.group.SetPosition(p.Subtract(self.offset))
}

func (self *AdaptorState) SetOffset(p wtype.Coordinates) {
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

//Represent a set of adaptors which are physically attached
type AdaptorGroup struct {
	adaptors     []*AdaptorState
	offsets      []wtype.Coordinates
	motionLimits *wtype.BBox
	position     wtype.Coordinates
	robot        *RobotState
}

func NewAdaptorGroup(offsets []wtype.Coordinates, motionLimits *wtype.BBox) *AdaptorGroup {
	ret := AdaptorGroup{
		adaptors:     make([]*AdaptorState, len(offsets)),
		offsets:      offsets,
		motionLimits: motionLimits,
	}

	return &ret
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
}

func (self *AdaptorGroup) GetPosition() wtype.Coordinates {
	return self.position
}

func (self *AdaptorGroup) SetPosition(p wtype.Coordinates) error {
	self.position = p
	if self.motionLimits != nil && !self.motionLimits.Contains(p) {
		return errors.New("movement limits prevent moving into position")
	}
	return nil
}

func (self *AdaptorGroup) GetRobot() *RobotState {
	return self.robot
}

func (self *AdaptorGroup) SetRobot(r *RobotState) {
	self.robot = r
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
