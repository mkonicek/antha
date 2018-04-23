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
	return self.adaptor.GetRobot().GetDeck().GetChildBelow(self.GetAbsolutePosition())
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
	name        string
	channels    []*ChannelState
	position    wtype.Coordinates
	independent bool
	params      *wtype.LHChannelParameter
	robot       *RobotState
}

func NewAdaptorState(name string,
	independent bool,
	channels int,
	channel_offset wtype.Coordinates,
	params *wtype.LHChannelParameter) *AdaptorState {
	as := AdaptorState{
		name,
		make([]*ChannelState, 0, channels),
		wtype.Coordinates{},
		independent,
		params.Dup(),
		nil,
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
	return self.position
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

//GetRobot
func (self *AdaptorState) GetRobot() *RobotState {
	return self.robot
}

//SetRobot
func (self *AdaptorState) SetRobot(r *RobotState) {
	self.robot = r
}

func (self *AdaptorState) SetPosition(p wtype.Coordinates) {
	self.position = p
}

// -------------------------------------------------------------------------------
//                            RobotState
// -------------------------------------------------------------------------------

//RobotState Represent the physical state of a liquidhandling robot
type RobotState struct {
	deck        *wtype.LHDeck
	adaptors    []*AdaptorState
	initialized bool
	finalized   bool
}

func NewRobotState() *RobotState {
	rs := RobotState{
		nil,
		make([]*AdaptorState, 0),
		false,
		false,
	}
	return &rs
}

//                            Accessors
//                            ---------

//GetAdaptor
func (self *RobotState) GetAdaptor(num int) *AdaptorState {
	return self.adaptors[num]
}

//GetNumberOfAdaptors
func (self *RobotState) GetNumberOfAdaptors() int {
	return len(self.adaptors)
}

//AddAdaptor
func (self *RobotState) AddAdaptor(a *AdaptorState) {
	a.SetRobot(self)
	self.adaptors = append(self.adaptors, a)
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
