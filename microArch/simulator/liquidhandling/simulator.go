// /anthalib/simulator/liquidhandling/simulator.go: Part of the Antha language
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
	"math"
	"sort"
	"strings"

	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/microArch/driver"
	"github.com/antha-lang/antha/microArch/driver/liquidhandling"
	"github.com/antha-lang/antha/microArch/simulator"
)

const arbitraryZOffset = 4.0

// Simulate a liquid handler Driver
type VirtualLiquidHandler struct {
	simulator.ErrorReporter
	state    *RobotState
	settings *SimulatorSettings
}

//coneRadius hardcoded radius to assume for cones
const coneRadius = 3.6

//Create a new VirtualLiquidHandler which mimics an LHDriver
func NewVirtualLiquidHandler(props *liquidhandling.LHProperties, settings *SimulatorSettings) *VirtualLiquidHandler {
	var vlh VirtualLiquidHandler

	if settings == nil {
		vlh.settings = DefaultSimulatorSettings()
	} else {
		vlh.settings = settings
	}

	vlh.validateProperties(props)
	//if the properties are that bad, don't bother building RobotState
	if vlh.HasError() {
		return &vlh
	}
	vlh.state = NewRobotState()

	//add the adaptors
	for _, assembly := range props.HeadAssemblies {
		offsets := make([]wtype.Coordinates, len(assembly.Positions))
		for i, pos := range assembly.Positions {
			offsets[i] = pos.Offset
		}
		group := NewAdaptorGroup(offsets, assembly.MotionLimits)

		for i, pos := range assembly.Positions {
			if pos.Head == nil {
				continue
			}
			p := pos.Head.Adaptor.Params
			//9mm spacing currently hardcoded.
			//At some point we'll either need to fetch this from the driver or
			//infer it from the type of tipboxes/plates accepted
			spacing := wtype.Coordinates{X: 0, Y: 0, Z: 0}
			if p.Orientation == wtype.LHVChannel {
				spacing.Y = 9.
			} else if p.Orientation == wtype.LHHChannel {
				spacing.X = 9.
			}
			adaptor := NewAdaptorState(pos.Head.Adaptor.Name, p.Independent, p.Multi, spacing, coneRadius, p, pos.Head.TipLoading)
			group.LoadAdaptor(i, adaptor)
		}
		vlh.state.AddAdaptorGroup(group)
	}

	//Make the deck
	deck := wtype.NewLHDeck("simulated deck", props.Mnfr, props.Model)
	for name, pos := range props.Layout {
		//size not given un LHProperties, assuming standard 96well size
		deck.AddSlot(name, pos, wtype.Coordinates{X: 127.76, Y: 85.48, Z: 0})
		//deck.SetSlotAccepts(name, "riser")
	}

	for _, name := range props.Tip_preferences {
		deck.SetSlotAccepts(name, "tipbox")
	}
	for _, name := range props.Input_preferences {
		deck.SetSlotAccepts(name, "plate")
	}
	for _, name := range props.Output_preferences {
		deck.SetSlotAccepts(name, "plate")
	}
	for _, name := range props.Tipwaste_preferences {
		deck.SetSlotAccepts(name, "tipwaste")
	}

	vlh.state.SetDeck(deck)

	return &vlh
}

// ------------------------------------------------------------------------------- Useful Utilities

func (self *VirtualLiquidHandler) validateProperties(props *liquidhandling.LHProperties) {

	//check a property
	check_prop := func(l []string, name string) {
		//is empty
		if len(l) == 0 {
			self.AddWarningf("NewVirtualLiquidHandler", "No %s specified", name)
		}
		//all locations defined
		for _, loc := range l {
			if _, ok := props.Layout[loc]; !ok {
				self.AddWarningf("NewVirtualLiquidHandler", "Undefined location \"%s\" referenced in %s", loc, name)
			}
		}
	}

	check_prop(props.Tip_preferences, "tip preferences")
	check_prop(props.Input_preferences, "input preferences")
	check_prop(props.Output_preferences, "output preferences")
	check_prop(props.Tipwaste_preferences, "tipwaste preferences")
	check_prop(props.Wash_preferences, "wash preferences")
	check_prop(props.Waste_preferences, "waste preferences")
}

//testSliceLength test that a bunch of slices are the correct length
func (self *VirtualLiquidHandler) testSliceLength(slice_lengths map[string]int, exp_length int) error {

	wrong := []string{}
	for name, actual_length := range slice_lengths {
		if actual_length != exp_length {
			wrong = append(wrong, fmt.Sprintf("%s(%d)", name, actual_length))
		}
	}

	if len(wrong) == 1 {
		return fmt.Errorf("Slice %s is not of expected length %v", wrong[0], exp_length)
	} else if len(wrong) > 1 {
		//for unit testing, names need to always be in the same order
		sort.Strings(wrong)
		return fmt.Errorf("Slices %s are not of expected length %v", strings.Join(wrong, ", "), exp_length)
	}
	return nil
}

func contains(v int, s []int) bool {
	for _, val := range s {
		if v == val {
			return true
		}
	}
	return false
}

//GetAdaptorState Currently we only support one adaptor group
func (self *VirtualLiquidHandler) GetAdaptorState(adaptor int) (*AdaptorState, error) {
	return self.state.GetAdaptor(0, adaptor)
}

func (self *VirtualLiquidHandler) GetObjectAt(slot string) wtype.LHObject {
	child, _ := self.state.GetDeck().GetChild(slot)
	return child
}

//testTipArgs check that load/unload tip arguments are valid insofar as they won't crash in RobotState
func (self *VirtualLiquidHandler) testTipArgs(f_name string, channels []int, head int, platetype, position, well []string) bool {
	//head should exist
	adaptor, err := self.GetAdaptorState(head)
	if err != nil {
		self.AddError(f_name, err.Error())
		return false
	}

	n_channels := adaptor.GetChannelCount()
	ret := true

	bad_channels := []string{}
	mchannels := map[int]bool{}
	dup_channels := []string{}
	for _, ch := range channels {
		if ch < 0 || ch >= n_channels {
			bad_channels = append(bad_channels, fmt.Sprintf("%v", ch))
		} else {
			if mchannels[ch] {
				dup_channels = append(dup_channels, fmt.Sprintf("%v", ch))
			} else {
				mchannels[ch] = true
			}
		}
	}
	if len(bad_channels) == 1 {
		self.AddErrorf(f_name, "Unknown channel \"%v\"", bad_channels[0])
		ret = false
	} else if len(bad_channels) > 1 {
		self.AddErrorf(f_name, "Unknown channels \"%v\"", strings.Join(bad_channels, "\",\""))
		ret = false
	}
	if len(dup_channels) == 1 {
		self.AddErrorf(f_name, "Channel%v appears more than once", dup_channels[0])
		ret = false
	} else if len(dup_channels) == 1 {
		self.AddErrorf(f_name, "Channels {%s} appear more than once", strings.Join(dup_channels, "\",\""))
		ret = false
	}

	if err := self.testSliceLength(map[string]int{
		"platetype": len(platetype),
		"position":  len(position),
		"well":      len(well)},
		n_channels); err != nil {

		self.AddError(f_name, err.Error())
		ret = false
	}

	if ret {
		for i := range platetype {
			if contains(i, channels) {
				if platetype[i] == "" && well[i] == "" && position[i] == "" {
					self.AddErrorf(f_name, "Command given for channel %d, but no platetype, well or position given", i)
					return false
				}
			} else if len(channels) > 0 { //if channels are empty we'll infer it later
				if !(platetype[i] == "" && well[i] == "" && position[i] == "") {
					self.AddWarningf(f_name, "No command for channel %d, but platetype, well or position given", i)
				}
			}
		}
	}
	return ret
}

type lhRet struct {
	channels    []int
	volumes     []float64
	is_explicit []bool
	adaptor     *AdaptorState
}

//Validate args to LH commands - aspirate, dispense, mix
func (self *VirtualLiquidHandler) validateLHArgs(head, multi int, platetype, what []string, volume []float64,
	flags map[string][]bool, cycles []int) (*lhRet, error) {

	var err error

	ret := lhRet{
		make([]int, 0),
		make([]float64, 0),
		nil,
		nil,
	}

	ret.adaptor, err = self.GetAdaptorState(head)
	if err != nil {
		return nil, err
	}

	if err := self.testSliceLength(map[string]int{
		"volume":    len(volume),
		"platetype": len(platetype),
		"what":      len(what),
	}, ret.adaptor.GetChannelCount()); err != nil {
		return nil, err
	}

	sla := map[string]int{}
	for k, v := range flags {
		sla[k] = len(v)
	}
	if cycles != nil {
		sla["cycles"] = len(cycles)
	}
	if err := self.testSliceLength(sla, ret.adaptor.GetChannelCount()); err != nil {
		return nil, err
	}

	negative := []float64{}
	ret.is_explicit = make([]bool, ret.adaptor.GetChannelCount())
	for i := range ret.is_explicit {
		ret.is_explicit[i] = !(platetype[i] == "" && what[i] == "")
		if ret.is_explicit[i] {
			ret.channels = append(ret.channels, i)
			ret.volumes = append(ret.volumes, volume[i])
			if volume[i] < 0. {
				negative = append(negative, volume[i])
			}
		}
	}

	if len(negative) == 1 {
		return &ret, fmt.Errorf("Cannot manipulate negative volume %s", summariseVolumes(negative))
	} else if len(negative) > 1 {
		return &ret, fmt.Errorf("Cannot manipulate negative volumes %s", summariseVolumes(negative))
	}

	if multi != len(ret.channels) {
		return &ret, fmt.Errorf("Multi was %d, but instructions given for %d channels (%s)", multi, len(ret.channels), summariseChannels(ret.channels))
	}

	return &ret, nil
}

//getTargetPosition get a position within the liquidhandler, adding any errors as neccessary
//bool is false if the instruction shouldn't continue (e.g. missing deckposition e.t.c)
func (self *VirtualLiquidHandler) getTargetPosition(fname, adaptorName string, channelIndex int, deckposition, platetype string, wc wtype.WellCoords, ref wtype.WellReference) (wtype.Coordinates, bool) {
	ret := wtype.Coordinates{}

	target, ok := self.state.GetDeck().GetChild(deckposition)
	if !ok {
		self.AddErrorf(fname, "Unknown location \"%s\"", deckposition)
		return ret, false
	}
	if target == nil {
		self.AddErrorf(fname, "No object found at position %s", deckposition)
		return ret, false
	}

	if (platetype != wtype.TypeOf(target)) &&
		(platetype != wtype.NameOf(target)) {
		self.AddWarningf(fname, "Object found at %s was type \"%s\", named \"%s\", not \"%s\" as expected",
			deckposition, wtype.TypeOf(target), wtype.NameOf(target), platetype)
	}

	addr, ok := target.(wtype.Addressable)
	if !ok {
		self.AddErrorf(fname, "Object \"%s\" at \"%s\" is not addressable", wtype.NameOf(target), deckposition)
		return ret, false
	}

	if !addr.AddressExists(wc) {
		self.AddErrorf(fname, "Request for well %s in object \"%s\" at \"%s\" which is of size [%dx%d]",
			wc.FormatA1(), wtype.NameOf(target), deckposition, addr.NRows(), addr.NCols())
		return ret, false
	}

	ret, ok = addr.WellCoordsToCoords(wc, ref)
	if !ok {
		//since we already checked that the address exists, this must be a bad reference
		self.AddErrorf(fname, "Object type %s at %s doesn't support reference \"%s\"",
			wtype.TypeOf(target), deckposition, ref)
		return ret, false
	}

	if targetted, ok := target.(wtype.Targetted); ok {
		targetOffset := targetted.GetTargetOffset(adaptorName, channelIndex)
		ret = ret.Add(targetOffset)
	}

	return ret, true
}

func (self *VirtualLiquidHandler) getWellsBelow(height float64, adaptor *AdaptorState) []*wtype.LHWell {
	tip_pos := make([]wtype.Coordinates, adaptor.GetChannelCount())
	wells := make([]*wtype.LHWell, adaptor.GetChannelCount())
	size := wtype.Coordinates{X: 0, Y: 0, Z: height}
	deck := self.state.GetDeck()
	for i := 0; i < adaptor.GetChannelCount(); i++ {
		if ch := adaptor.GetChannel(i); ch.HasTip() {
			tip_pos[i] = ch.GetAbsolutePosition().Subtract(wtype.Coordinates{X: 0., Y: 0., Z: ch.GetTip().GetEffectiveHeight()})
		} else {
			tip_pos[i] = ch.GetAbsolutePosition()
		}

		for _, o := range deck.GetBoxIntersections(*wtype.NewBBox(tip_pos[i].Subtract(size), size)) {
			if w, ok := o.(*wtype.LHWell); ok {
				wells[i] = w
				break
			}
		}
	}

	return wells
}

func makeOffsets(Xs, Ys, Zs []float64) []wtype.Coordinates {
	ret := make([]wtype.Coordinates, len(Xs))
	for i := range Xs {
		ret[i].X = Xs[i]
		ret[i].Y = Ys[i]
		ret[i].Z = Zs[i]
	}
	return ret
}

// ------------------------------------------------------------------------ ExtendedLHDriver

//Move command - used
func (self *VirtualLiquidHandler) Move(deckpositionS []string, wellcoords []string, reference []int,
	offsetX, offsetY, offsetZ []float64, platetypeS []string,
	head int) driver.CommandStatus {
	ret := driver.CommandStatus{OK: true, Errorcode: driver.OK, Msg: "MOVE ACK"}

	//get the adaptor
	adaptor, err := self.GetAdaptorState(head)
	if err != nil {
		self.AddError("Move", err.Error())
		return ret
	}

	//only support a single deckposition or platetype
	deckposition, err := getSingle(deckpositionS)
	if err != nil {
		self.AddErrorf("Move", "invalid argument deckposition: %s", err.Error())
	}
	platetype, err := getSingle(platetypeS)
	if err != nil {
		self.AddErrorf("Move", "invalid argument platetype: %s", err.Error())
	}

	//extend args
	wellcoords = extend_strings(adaptor.GetChannelCount(), wellcoords)
	reference = extend_ints(adaptor.GetChannelCount(), reference)
	offsetX = extend_floats(adaptor.GetChannelCount(), offsetX)
	offsetY = extend_floats(adaptor.GetChannelCount(), offsetY)
	offsetZ = extend_floats(adaptor.GetChannelCount(), offsetZ)

	//check slice length
	if err := self.testSliceLength(map[string]int{
		"wellcoords": len(wellcoords),
		"reference":  len(reference),
		"offsetX":    len(offsetX),
		"offsetY":    len(offsetY),
		"offsetZ":    len(offsetZ),
	},
		adaptor.GetChannelCount()); err != nil {

		self.AddError("Move", err.Error())
		return ret
	}

	refs, err := convertReferences(reference)
	if err != nil {
		self.AddErrorf("Move", "invalid argument reference: %s", err.Error())
	}

	//get slice of well coords
	wc, err := convertWellCoords(wellcoords)
	if err != nil {
		self.AddErrorf("Move", "invalid argument wellcoords: %s", err.Error())
	}

	//get the channels from the well coords
	channels := make([]int, 0, len(wc))
	implicitChannels := make([]int, 0, len(wc))
	for i, w := range wc {
		if !w.IsZero() {
			channels = append(channels, i)
		} else {
			implicitChannels = append(implicitChannels, i)
		}
	}
	if len(channels) == 0 {
		self.AddWarning("Move", "ignoring blank move command: no wellcoords specified")
		return ret
	}

	//combine floats into wtype.Coordinates
	offsets := makeOffsets(offsetX, offsetY, offsetZ)

	//find the coordinates of each explicitly requested position
	coords := make([]wtype.Coordinates, adaptor.GetChannelCount())
	for _, ch := range channels {
		c, ok := self.getTargetPosition("Move", adaptor.GetName(), ch, deckposition, platetype, wc[ch], refs[ch])
		if !ok {
			return ret
		}
		coords[ch] = c
		coords[ch] = coords[ch].Add(offsets[ch])
		//if there's a tip, raise the coortinates to the top of the tip to take account of it
		if tip := adaptor.GetChannel(ch).GetTip(); tip != nil {
			coords[ch] = coords[ch].Add(wtype.Coordinates{X: 0., Y: 0., Z: tip.GetEffectiveHeight()})
		}
	}

	target, ok := self.state.GetDeck().GetChild(deckposition)
	if !ok {
		self.AddErrorf("Move", "unable to get object at position \"%s\"", deckposition)
	}

	describe := func() string {
		return fmt.Sprintf("head %d %s to %s@%s at position %s",
			head, summariseChannels(channels), wtype.HumanizeWellCoords(wc), wtype.NameOf(target), deckposition)
	}

	//find the head location
	//for now, assuming that the relative position of the first explicitly provided channel and the head stay
	//the same. This seems sensible for the Glison, but might turn out not to be how other robots with independent channels work
	origin := coords[channels[0]].Subtract(adaptor.GetChannel(channels[0]).GetRelativePosition())

	//fill in implicit locations
	for _, ch := range implicitChannels {
		coords[ch] = origin.Add(adaptor.GetChannel(ch).GetRelativePosition())
	}

	//Get the locations of each channel relative to the head
	rel_coords := make([]wtype.Coordinates, adaptor.GetChannelCount())
	for i := range coords {
		rel_coords[i] = coords[i].Subtract(origin)
	}

	//check that the requested position is possible given the head/adaptor capabilities
	if !adaptor.IsIndependent() {
		//i.e. the channels can't move relative to each other or the head, so relative locations must remain the same
		moved := make([]int, 0, len(rel_coords))
		for i, rc := range rel_coords {
			//check that adaptor relative position remains the same
			//arbitrary 0.01mm to avoid numerical instability
			if rc.Subtract(adaptor.GetChannel(i).GetRelativePosition()).Abs() > 0.01 {
				moved = append(moved, i)
			}
		}
		if len(moved) > 0 {
			self.AddErrorf("Move", "%s: requires moving %s relative to non-independent head",
				describe(), summariseChannels(moved))
			return ret
		}
	}

	//check that there are no tips loaded on any other heads
	if err = assertNoTipsOnOthersInGroup(adaptor); err != nil {
		self.AddErrorf("Move", "%s: cannot move head %d while %s", describe(), head, err.Error())
	}

	//move the head to the new position
	err = adaptor.SetPosition(origin)
	if err != nil {
		self.AddErrorf("Move", "%s: %s", describe(), err.Error())
	}
	for i, rc := range rel_coords {
		adaptor.GetChannel(i).SetRelativePosition(rc)
	}

	//check for collisions in the new location
	if err := assertNoCollisionsInGroup(adaptor, nil, 0.0); err != nil {
		self.AddErrorf("Move", "%s: collision detected: %s", describe(), err.Error())
	}
	return ret
}

//Move raw - not yet implemented in compositerobotinstruction
func (self *VirtualLiquidHandler) MoveRaw(head int, x, y, z float64) driver.CommandStatus {
	self.AddWarning("MoveRaw", "Not yet implemented")
	return driver.CommandStatus{OK: true, Errorcode: driver.OK, Msg: "MOVERAW ACK"}
}

//Aspirate - used
func (self *VirtualLiquidHandler) Aspirate(volume []float64, overstroke []bool, head int, multi int,
	platetype []string, what []string, llf []bool) driver.CommandStatus {

	ret := driver.CommandStatus{OK: true, Errorcode: driver.OK, Msg: "ASPIRATE ACK"}

	//extend arguments - at some point shortform slices might become illegal
	if adaptor, err := self.GetAdaptorState(head); err == nil {
		nc := adaptor.GetChannelCount()
		volume = extend_floats(nc, volume)
		overstroke = extend_bools(nc, overstroke)
		platetype = extend_strings(nc, platetype)
		what = extend_strings(nc, what)
		llf = extend_bools(nc, llf)
	}

	arg, err := self.validateLHArgs(head, multi, platetype, what, volume,
		map[string][]bool{
			"overstroke": overstroke,
			"llf":        llf,
		}, nil)
	if err != nil {
		self.AddErrorf("Aspirate", "Invalid Arguments - %s", err.Error())
		return ret
	}

	describe := func() string {
		return fmt.Sprintf("aspirating %s of %s to head %d %s",
			summariseVolumes(arg.volumes), summariseStrings(what), head, summariseChannels(arg.channels))
	}

	//get the position of tips
	wells := self.getWellsBelow(0.0, arg.adaptor)

	//check if any explicitly requested channels are missing tips
	tip_missing := []int{}
	for _, i := range arg.channels {
		if !arg.adaptor.GetChannel(i).HasTip() {
			tip_missing = append(tip_missing, i)
		}
	}
	if len(tip_missing) > 0 {
		self.AddErrorf("Aspirate", "While %s - missing %s on %s", describe(), pTips(len(tip_missing)), summariseChannels(tip_missing))
		return ret
	}

	//independence constraints
	if !arg.adaptor.IsIndependent() {
		different := false
		v := -1.
		for ch, b := range arg.is_explicit {
			if b {
				if v >= 0 {
					different = different || v != volume[ch]
				} else {
					v = volume[ch]
				}
			} else if wells[ch] != nil {
				//a non-explicitly requested tip is in a well. If the well has stuff in it, it'll get aspirated
				if c := wells[ch].Contents(); !c.IsZero() {
					self.AddErrorf("Aspirate",
						"While %s - channel %d will inadvertantly aspirate %s from well %s as head is not independent",
						describe(), ch, c.Name(), wells[ch].GetName())
				}
			}
		}
		if different {
			self.AddErrorf("Aspirate", "While %s - channels cannot aspirate different volumes in non-independent head", describe())
			return ret
		}
	}

	//check liquid type
	for i := range arg.channels {
		if wells[i] == nil { //we'll catch this later
			continue
		}
		if wells[i].Contents().GetType() != what[i] && self.settings.IsLiquidTypeWarningEnabled() {
			self.AddWarningf("Aspirate", "While %s - well %s contains %s, not %s",
				describe(), wells[i].GetName(), wells[i].Contents().GetType(), what[i])
		}
	}

	//check total volumes taken from each unique well
	uniqueWells := make(map[string]*wtype.LHWell)
	uniqueWellVolumes := make(map[string]float64)
	uniqueWellVolumeIndexes := make(map[string][]int)
	for i := 0; i < len(wells); i++ {
		if wells[i] == nil {
			continue
		}
		if _, ok := uniqueWells[wells[i].ID]; !ok {
			uniqueWells[wells[i].ID] = wells[i]
			uniqueWellVolumes[wells[i].ID] = 0.0
			uniqueWellVolumeIndexes[wells[i].ID] = make([]int, 0, len(wells))
		}
		uniqueWellVolumes[wells[i].ID] += volume[i]
		uniqueWellVolumeIndexes[wells[i].ID] = append(uniqueWellVolumeIndexes[wells[i].ID], i)
	}
	for id, well := range uniqueWells {
		v := wunit.NewVolume(uniqueWellVolumes[id], "ul")
		//vol.IsZero() checks whether vol is within a small tolerance of zero
		if d := wunit.SubtractVolumes(v, well.CurrentWorkingVolume()); v.GreaterThan(well.CurrentWorkingVolume()) && !d.IsZero() {
			//the volume is taken from len(uniqueWellVolumeIndexes[id]) wells, so the delta is split equally between them
			reduction := wunit.DivideVolume(d, float64(len(uniqueWellVolumeIndexes[id])))
			reductionUl := reduction.ConvertToString("ul")
			for _, i := range uniqueWellVolumeIndexes[id] {
				volume[i] -= reductionUl
			}
			self.AddWarningf("Aspirate", "While %s - well %s only contains %s working volume, reducing aspirated volume by %v",
				describe(), well.GetName(), well.CurrentWorkingVolume(), reduction)
		}
	}

	//move liquid
	no_well := []int{}
	for _, i := range arg.channels {
		v := wunit.NewVolume(volume[i], "ul")
		tip := arg.adaptor.GetChannel(i).GetTip()
		fv := tip.CurrentVolume()
		fv.Add(v)

		if wells[i] == nil {
			no_well = append(no_well, i)
		} else if wells[i].CurrentWorkingVolume().LessThan(v) {
			self.AddErrorf("Aspirate", "While %s - well %s only contains %s working volume",
				describe(), wells[i].GetName(), wells[i].CurrentWorkingVolume())
		} else if fv.GreaterThan(tip.MaxVol) {
			self.AddErrorf("Aspirate", "While %s - channel %d contains %s, command exceeds maximum volume %s",
				describe(), i, tip.CurrentVolume(), tip.MaxVol)
		} else if c, err := wells[i].RemoveVolume(v); err != nil {
			self.AddErrorf("Aspirate", "While %s - unexpected well error \"%s\"", describe(), err.Error())
		} else if fv.LessThan(tip.MinVol) {
			self.AddWarningf("Aspirate", "While %s - minimum tip volume is %s",
				describe(), tip.MinVol)
			//will get an error here, but ignore it since we're already raising a warning
			addComponent(tip, c) //nolint
		} else if err := addComponent(tip, c); err != nil {
			self.AddErrorf("Aspirate", "While %s - unexpected tip error \"%s\"", describe(), err.Error())
		}
	}

	if len(no_well) > 0 {
		self.AddErrorf("Aspirate", "While %s - %s on %s not in a well", describe(), pTips(len(no_well)), summariseChannels(no_well))
	}

	return ret
}

//Dispense - used
func (self *VirtualLiquidHandler) Dispense(volume []float64, blowout []bool, head int, multi int,
	platetype []string, what []string, llf []bool) driver.CommandStatus {

	ret := driver.CommandStatus{OK: true, Errorcode: driver.OK, Msg: "DISPENSE ACK"}

	//extend arguments - at some point shortform slices might become illegal
	if adaptor, err := self.GetAdaptorState(head); err == nil {
		volume = extend_floats(adaptor.GetChannelCount(), volume)
		blowout = extend_bools(adaptor.GetChannelCount(), blowout)
		platetype = extend_strings(adaptor.GetChannelCount(), platetype)
		what = extend_strings(adaptor.GetChannelCount(), what)
		llf = extend_bools(adaptor.GetChannelCount(), llf)
	}

	arg, err := self.validateLHArgs(head, multi, platetype, what, volume, map[string][]bool{
		"blowout": blowout,
		"llf":     llf,
	}, nil)
	if err != nil {
		self.AddErrorf("Dispense", "Invalid arguments - %s", err.Error())
		return ret
	}

	//find the position of each tip
	wells := self.getWellsBelow(self.settings.MaxDispenseHeight(), arg.adaptor)

	whatS := make([]string, 0, len(what))
	for _, i := range arg.channels {
		whatS = append(whatS, what[i])
	}

	describe := func() string {
		return fmt.Sprintf("%s of %s from head %d %s to %s", summariseVolumes(arg.volumes), summariseStrings(whatS), head, summariseChannels(arg.channels), summarisePlateWells(wells, arg.channels))
	}

	//check wells
	noWell := []int{}
	for _, i := range arg.channels {
		if wells[i] == nil {
			noWell = append(noWell, i)
		}
	}
	if len(noWell) > 0 {
		self.AddErrorf("Dispense", "%s : no well within %s below %s on %s",
			describe(), wunit.NewLength(self.settings.MaxDispenseHeight(), "mm"), pTips(len(noWell)), summariseChannels(noWell))
		return ret
	}

	//check tips
	noTip := make([]int, 0, len(arg.channels))
	for _, i := range arg.channels {
		if !arg.adaptor.GetChannel(i).HasTip() {
			noTip = append(noTip, i)
		}
	}
	if len(noTip) > 0 {
		self.AddErrorf("Dispense", "%s : no %s loaded on %s",
			describe(), pTips(len(noTip)), summariseChannels(noTip))
		return ret
	}

	//independence contraints
	if !arg.adaptor.IsIndependent() {
		extra := []int{}
		different := false
		v := -1.
		for ch, b := range arg.is_explicit {
			if b {
				if v >= 0 {
					different = different || v != volume[ch]
				} else {
					v = volume[ch]
				}
			} else if wells[ch] != nil {
				//a non-explicitly requested tip is in a well. If the well has stuff in it, it'll get aspirated
				if c := arg.adaptor.GetChannel(ch).GetTip().Contents(); !c.IsZero() {
					extra = append(extra, ch)
				}
			}
		}
		if different {
			self.AddErrorf("Dispense", "%s : channels cannot dispense different volumes in non-independent head", describe())
			return ret
		} else if len(extra) > 0 {
			self.AddErrorf("Dispense",
				"%s : must also dispense %s from %s as head is not independent",
				describe(), summariseVolumes(arg.volumes), summariseChannels(extra))
			return ret
		}
	}

	//check liquid type
	for i := range arg.channels {
		if tip := arg.adaptor.GetChannel(i).GetTip(); tip != nil {
			if tip.Contents().GetType() != what[i] && self.settings.IsLiquidTypeWarningEnabled() {
				self.AddWarningf("Dispense", "%s : channel %d contains %s, not %s",
					describe(), i, tip.Contents().GetType(), what[i])
			}
		}
	}

	//for each blowout channel
	for i := range arg.channels {
		if !blowout[i] {
			continue
		}
		//reduce the volume to the total volume in the tip (assume blowout removes residual volume as well)
		cv := arg.adaptor.GetChannel(i).GetTip().CurrentVolume().ConvertToString("ul")
		volume[i] = math.Min(volume[i], cv)
	}

	//check volumes -- currently only warnings due to poor volume tracking
	finalVolumes := make([]float64, 0, len(arg.channels))
	maxVolumes := make([]float64, 0, len(arg.channels))
	overfullWells := make([]int, 0, len(arg.channels))
	for _, i := range arg.channels {
		cV := wells[i].CurrentVolume()
		mV := wells[i].MaxVolume()
		fV := wunit.AddVolumes(cV, wunit.NewVolume(volume[i], "ul"))
		if delta := wunit.SubtractVolumes(fV, mV); fV.GreaterThan(mV) && !delta.IsZero() {
			overfullWells = append(overfullWells, i)
			finalVolumes = append(finalVolumes, fV.ConvertToString("ul"))
			maxVolumes = append(maxVolumes, mV.ConvertToString("ul"))
		}
	}
	if len(overfullWells) > 0 {
		self.AddWarningf("Dispense", "%s : overfilling %s %s to %s of %s max volume",
			describe(), pWells(len(overfullWells)), summarisePlateWells(wells, overfullWells), summariseVolumes(finalVolumes), summariseVolumes(maxVolumes))
	}

	//dispense
	for _, i := range arg.channels {
		v := wunit.NewVolume(volume[i], "ul")
		tip := arg.adaptor.GetChannel(i).GetTip()

		if wells[i] != nil {
			if _, tw := wells[i].Plate.(*wtype.LHTipwaste); tw {
				self.AddWarningf("Dispense", "%s : dispensing to tipwaste", describe())
			}
		}

		if v.GreaterThan(tip.CurrentWorkingVolume()) {
			v = tip.CurrentWorkingVolume()
			if !blowout[i] {
				//a bit strange
				self.AddWarningf("Dispense", "%s : tip on channel %d contains only %s, but blowout flag is false",
					describe(), i, tip.CurrentWorkingVolume())
			}
		}
		if c, err := tip.RemoveVolume(v); err != nil {
			self.AddErrorf("Dispense", "%s : unexpected tip error \"%s\"", describe(), err.Error())
		} else if err := addComponent(wells[i], c); err != nil {
			self.AddErrorf("Dispense", "%s : unexpected well error \"%s\"", describe(), err.Error())
		}
	}

	return ret
}

//LoadTips - used
func (self *VirtualLiquidHandler) LoadTips(channels []int, head, multi int,
	platetypeS, positionS, well []string) driver.CommandStatus {
	ret := driver.CommandStatus{OK: true, Errorcode: driver.OK, Msg: "LOADTIPS ACK"}
	deck := self.state.GetDeck()

	//get the adaptor
	adaptor, err := self.GetAdaptorState(head)
	if err != nil {
		self.AddError("LoadTips", err.Error())
		return ret
	}
	n_channels := adaptor.GetChannelCount()

	//extend arg slices
	platetypeS = extend_strings(n_channels, platetypeS)
	positionS = extend_strings(n_channels, positionS)
	well = extend_strings(n_channels, well)

	//check that the command is valid
	if !self.testTipArgs("LoadTips", channels, head, platetypeS, positionS, well) {
		return ret
	}

	//get the individual position
	position, err := getSingle(positionS)
	if err != nil {
		self.AddErrorf("LoadTips", "invalid argument position: %s", err.Error())
	}
	platetype, err := getSingle(platetypeS)
	if err != nil {
		self.AddErrorf("LoadTips", "invalid argument platetype: %s", err.Error())
	}

	if len(channels) == 0 {
		//inver channels from well argument
		channels = make([]int, 0, n_channels)
		for i, w := range well {
			if w != "" {
				channels = append(channels, i)
			}
		}
	}

	//make well coords
	invalidWells := make([]int, 0, n_channels)
	wc, err := convertWellCoords(well)
	if err != nil {
		self.AddErrorf("LoadTips", "invalid argument well: %s", err.Error())
	}
	for _, ch := range channels {
		if wc[ch].IsZero() {
			invalidWells = append(invalidWells, ch)
		}
	}
	if len(invalidWells) > 0 {
		values := make([]string, 0, len(invalidWells))
		for _, ch := range invalidWells {
			values = append(values, well[ch])
		}
		self.AddErrorf("LoadTips", "invalid argument well: couldn't parse \"%s\" for %s", strings.Join(values, "\", \""), summariseChannels(invalidWells))
	}

	//get the actual tipbox
	var tipbox *wtype.LHTipbox
	if o, ok := deck.GetChild(position); !ok {
		self.AddErrorf("LoadTips", "unknown location \"%s\"", position)
		return ret
	} else if o == nil {
		self.AddErrorf("LoadTips", "can't load tips from empty position \"%s\"", position)
		return ret
	} else if tipbox, ok = o.(*wtype.LHTipbox); !ok {
		self.AddErrorf("LoadTips", "can't load tips from %s \"%s\" found at position \"%s\"",
			wtype.ClassOf(o), wtype.NameOf(o), position)
		return ret
	}
	if tipbox == nil {
		self.AddErrorf("LoadTips", "unexpected nil tipbox at position \"%s\"", position)
		return ret
	}

	describe := func() string {
		return fmt.Sprintf("from %s@%s at position \"%s\" to head %d %s", wtype.HumanizeWellCoords(wc), tipbox.GetName(), position, head, summariseChannels(channels))
	}

	if multi != len(channels) {
		self.AddErrorf("LoadTips", "%s : multi should equal %d, not %d",
			describe(), len(channels), multi)
		return ret
	}

	//check that channels we want to load to are empty
	if tipFound := checkTipPresence(false, adaptor, channels); len(tipFound) != 0 {
		self.AddErrorf("LoadTips", "%s: %s already loaded on head %d %s",
			describe(), pTips(len(tipFound)), head, summariseChannels(tipFound))
		return ret
	}

	//check that there aren't any tips loaded on any other heads
	if err := assertNoTipsOnOthersInGroup(adaptor); err != nil {
		self.AddErrorf("LoadTips", "%s: while %s", describe(), err.Error())
	}

	//refill the tipbox if there aren't enough tips to service the instruction
	if adaptor.AutoRefillsTipboxes() && !tipbox.HasEnoughTips(multi) {
		tipbox.Refill()
	}

	//if the adaptor might override what we tell it
	if adaptor.OverridesLoadTipsCommand() && self.settings.IsTipLoadingOverrideEnabled() {
		//a list of tip locations that will be loaded
		tipChunks, err := adaptor.GetTipCoordsToLoad(tipbox, multi)
		if err != nil {
			self.AddErrorf("LoadTips", "%s : unexpected error : %s", describe(), err.Error())
			return ret
		}
		if !coordsMatch(tipChunks, wc) {
			return self.overrideLoadTips(channels, head, multi, platetype, position, tipChunks)
		}
	}

	//Get the tip at each requested location
	tips := make([]*wtype.LHTip, n_channels)
	var missingTips []wtype.WellCoords
	for _, i := range channels {

		if wc[i].IsZero() {
			pos := adaptor.GetChannel(i).GetAbsolutePosition()
			wc[i], _ = tipbox.CoordsToWellCoords(pos)
			if !wc[i].IsZero() {
				self.AddWarningf("LoadTips",
					"%s : Well coordinates for channel %d not specified, assuming %s from adaptor location",
					describe(), i, wc[i].FormatA1())
			}
		}

		if !tipbox.AddressExists(wc[i]) {
			self.AddErrorf("LoadTips", "%s : request for tip at %s in tipbox of size [%dx%d]",
				describe(), wc[i].FormatA1(), tipbox.NCols(), tipbox.NRows())
			return ret
		} else {
			tips[i] = tipbox.GetChildByAddress(wc[i]).(*wtype.LHTip)
			if tips[i] == nil {
				missingTips = append(missingTips, wc[i])
			}
		}
	}
	if len(missingTips) > 0 {
		self.AddErrorf("LoadTips", "%s : no %s at %s",
			describe(), pTips(len(missingTips)), wtype.HumanizeWellCoords(missingTips))
		return ret
	}

	//check alignment
	z_off := make([]float64, n_channels)
	misaligned := []int{}
	target := []wtype.WellCoords{}
	amount := []string{}
	for _, ch := range channels {
		tip_s := tips[ch].GetSize()
		tip_p := tips[ch].GetPosition().Add(wtype.Coordinates{X: 0.5 * tip_s.X, Y: 0.5 * tip_s.Y, Z: tip_s.Z})
		ch_p := adaptor.GetChannel(ch).GetAbsolutePosition()
		delta := ch_p.Subtract(tip_p)
		if xy := delta.AbsXY(); xy > 0.5 {
			misaligned = append(misaligned, ch)
			target = append(target, wc[ch])
			amount = append(amount, fmt.Sprintf("%v", xy))
		}
		z_off[ch] = delta.Z
		if delta.Z < 0. {
			self.AddErrorf("LoadTips", "%s : channel is %.1f below tip", describe(), -delta.Z)
			return ret
		}
	}
	if len(misaligned) != 0 {
		is := "is"
		res := ""
		if len(misaligned) != 1 {
			is = "are"
			res = " respectively"
		}
		self.AddErrorf("LoadTips", "%s : %s %s misaligned with %s at %s by %smm%s",
			describe(), summariseChannels(misaligned), is, pTips(len(misaligned)), wtype.HumanizeWellCoords(target), strings.Join(amount, ","), res)
		return ret
	}

	//if not independent, check there are no other tips in the way
	if !adaptor.IsIndependent() {
		zo_max := 0.
		zo_min := math.MaxFloat64
		for _, ch := range channels {
			if z_off[ch] > zo_max {
				zo_max = z_off[ch]
			}
			if z_off[ch] < zo_min {
				zo_min = z_off[ch]
			}
		}
		if zo_max != zo_min {
			self.AddErrorf("LoadTips", "%s : distance between channels and tips varies from %v to %v mm in non-independent head",
				describe(), zo_min, zo_max)
			return ret
		}
		if err := assertNoCollisionsInGroup(adaptor, channels, zo_max+0.5); err != nil {
			self.AddErrorf("LoadTips", "%s: collision detected: %s", describe(), err.Error())
		}
	}

	//move the tips to the adaptors
	for _, ch := range channels {
		tips[ch].GetParent().(*wtype.LHTipbox).RemoveTip(wc[ch])
		adaptor.GetChannel(ch).LoadTip(tips[ch])
		tips[ch].ClearParent()
	}

	return ret
}

//overrideLoadTips sequentially load the given series of tips onto the given channels
func (self *VirtualLiquidHandler) overrideLoadTips(channels []int, head, multi int, platetype, position string, tipChunks [][]wtype.WellCoords) driver.CommandStatus {
	//make certain that any load tips we generate don't get overridden again
	self.settings.EnableTipLoadingOverride(false)
	defer self.settings.EnableTipLoadingOverride(true)

	var ret driver.CommandStatus
	loadedChannels := make([]int, 0, len(channels))

	for _, chunk := range tipChunks {
		width := len(chunk)
		channelsToLoad := channels[len(loadedChannels) : len(loadedChannels)+width]
		positionS := make([]string, multi)
		platetypeS := make([]string, multi)
		reference := make([]int, multi)
		offsetXY := make([]float64, multi)
		offsetZ := make([]float64, multi)
		wellcoords := make([]string, multi)
		for i, ch := range channelsToLoad {
			positionS[ch] = position
			platetypeS[ch] = platetype
			reference[ch] = int(wtype.TopReference)
			offsetXY[ch] = 0.0
			//arbitrary since we don't know the exact height and it won't affect collision detection
			offsetZ[ch] = arbitraryZOffset
			wellcoords[ch] = chunk[i].FormatA1()
		}

		ret = self.Move(positionS, wellcoords, reference, offsetXY, offsetXY, offsetZ, platetypeS, head)
		if self.HasError() {
			return ret
		}
		ret = self.LoadTips(channelsToLoad, head, width, platetypeS, positionS, wellcoords)
		if self.HasError() {
			return ret
		}

		loadedChannels = append(loadedChannels, channelsToLoad...)
	}

	return ret
}

//UnloadTips - used
func (self *VirtualLiquidHandler) UnloadTips(channels []int, head, multi int,
	platetype, position, well []string) driver.CommandStatus {
	ret := driver.CommandStatus{OK: true, Errorcode: driver.OK, Msg: "UNLOADTIPS ACK"}

	//get the adaptor
	adaptor, err := self.GetAdaptorState(head)
	if err != nil {
		self.AddError("UnloadTips", err.Error())
		return ret
	}
	n_channels := adaptor.GetChannelCount()

	//extend arg slices
	platetype = extend_strings(n_channels, platetype)
	position = extend_strings(n_channels, position)
	well = extend_strings(n_channels, well)

	if len(channels) == 0 {
		for ch := 0; ch < n_channels; ch++ {
			if adaptor.GetChannel(ch).HasTip() {
				channels = append(channels, ch)
			}
		}
		sort.Ints(channels)
		if len(channels) == 0 {
			self.AddWarningf("UnloadTips", "'channel' argument empty and no tips are loaded to head %d, ignoring", head)
		} else if self.settings.IsAutoChannelWarningEnabled() {
			self.AddWarningf("UnloadTips", "'channel' argument empty, unloading all tips (%s)", summariseChannels(channels))
		}
	}

	//check that RobotState won't crash
	if !self.testTipArgs("UnloadTips", channels, head, platetype, position, well) {
		return ret
	}

	if multi != len(channels) {
		self.AddWarningf("UnloadTips", "While unloading %s from %s, multi should equal %d, not %d",
			pTips(len(channels)), summariseChannels(channels), len(channels), multi)
		//multi = len(channels) - multi is unused
	}

	deck := self.state.GetDeck()

	//Raise a warning if we're trying to eject tips that aren't there
	missing := []string{}
	for _, ch := range channels {
		if !adaptor.GetChannel(ch).HasTip() {
			missing = append(missing, fmt.Sprintf("%d", ch))
		}
	}
	if len(missing) == 1 {
		self.AddWarningf("UnloadTips", "No tip present at Head%d channel%s to eject", head, missing[0])
	} else if len(missing) > 0 {
		self.AddWarningf("UnloadTips", "No tips present on Head%d channels %s to eject", head, strings.Join(missing, ","))
	}

	//Check that this is possible
	if !adaptor.IsIndependent() {
		extra := []int{}
		for ch := 0; ch < n_channels; ch++ {
			if contains(ch, channels) {
				continue
			}
			if adaptor.GetChannel(ch).HasTip() {
				extra = append(extra, ch)
			}
		}
		if len(extra) > 0 {
			self.AddErrorf("UnloadTips", "Cannot unload tips from head%d %s without unloading %s from %s (head isn't independent)",
				head, summariseChannels(channels), pTips(len(extra)), summariseChannels(extra))
			return ret
		}
	}

	for _, ch := range channels {
		//get the target
		if target, ok := deck.GetChild(position[ch]); !ok {
			self.AddErrorf("UnloadTips", "Unknown deck position \"%s\"", position[ch])
			break
		} else if target == nil {
			self.AddErrorf("UnloadTips", "Cannot unload to empty deck location \"%s\"", position[ch])
			break
		} else if addr, ok := target.(wtype.Addressable); !ok {
			self.AddErrorf("UnloadTips", "Cannot unload tips to %s \"%s\" at location %s",
				wtype.ClassOf(target), wtype.NameOf(target), position[ch])
		} else {
			//get the location of the channel
			ch_pos := adaptor.GetChannel(ch).GetAbsolutePosition()
			//parse the wellcoords
			wc := wtype.MakeWellCoords(well[ch])
			if wc.IsZero() {
				self.AddErrorf("UnloadTips", "Cannot parse well coordinates \"%s\"", well[ch])
				break
			}
			if !addr.AddressExists(wc) {
				self.AddErrorf("UnloadTips", "Cannot unload to address %s in %s \"%s\" size [%dx%d]",
					wc.FormatA1(), wtype.ClassOf(target), wtype.NameOf(target), addr.NRows(), addr.NCols())
				break
			}
			//get the child - *LHTip or *LHWell
			child := addr.GetChildByAddress(wc)
			well_p, _ := addr.WellCoordsToCoords(wc, wtype.TopReference)
			delta := ch_pos.Subtract(well_p)

			switch target := target.(type) {
			case *wtype.LHTipbox:
				//put the tip in the tipbox
				if child.(*wtype.LHTip) != nil {
					self.AddErrorf("UnloadTips", "Cannot unload to tipbox \"%s\" %s, tip already present there",
						target.GetName(), wc.FormatA1())
				} else if delta.AbsXY() > 0.25 {
					self.AddErrorf("UnloadTips", "Head%d channel%d misaligned from tipbox \"%s\" %s by %.2fmm",
						head, ch, target.GetName(), wc.FormatA1(), delta.AbsXY())
				} else if delta.Z > target.GetSize().Z/2. {
					self.AddWarningf("UnloadTips", "Ejecting tip from Head%d channel%d to tipbox \"%s\" %s from height of %.2fmm",
						head, ch, target.GetName(), wc.FormatA1(), delta.Z)
				} else {
					target.PutTip(wc, adaptor.GetChannel(ch).UnloadTip())
				}

			case *wtype.LHTipwaste:
				//put the tip in the tipwaste
				if child == nil {
					//I don't think this should happen, but it would be embarressing to segfault...
					self.AddWarningf("UnloadTips", "Tipwaste \"%s\" well %s was nil, cannot check head alignment",
						target.GetName(), wc.FormatA1())
					adaptor.GetChannel(ch).UnloadTip()
				} else if max_delta := child.GetSize(); delta.X > max_delta.X || delta.Y > max_delta.Y {
					self.AddErrorf("UnloadTips", "Cannot unload, head%d channel%d is not above tipwaste \"%s\"",
						head, ch, target.GetName())
				} else if target.SpaceLeft() <= 0 {
					self.AddErrorf("UnloadTips", "Cannot unload tip to overfull tipwaste \"%s\", contains %d tips",
						target.GetName(), target.Contents)
				} else {
					target.DisposeNum(1)
					adaptor.GetChannel(ch).UnloadTip()
				}
			default:
				self.AddErrorf("UnloadTips", "Cannot unload tips to %s \"%s\" at location %s",
					wtype.ClassOf(target), wtype.NameOf(target), position[ch])
			}
		}
		if self.HasError() {
			break
		}
	}

	return ret
}

//SetPipetteSpeed - used
func (self *VirtualLiquidHandler) SetPipetteSpeed(head, channel int, rate float64) driver.CommandStatus {
	ret := driver.CommandStatus{OK: true, Errorcode: driver.OK, Msg: "SETPIPETTESPEED ACK"}

	adaptor, err := self.GetAdaptorState(head)
	if err != nil {
		self.AddError("SetPipetteSpeed", err.Error())
		return ret
	}

	channels := make([]int, 0, adaptor.GetChannelCount())
	if channel < 0 || !adaptor.IsIndependent() {
		if channel >= 0 {
			self.AddWarningf("SetPipetteSpeed", "Head %d is not independent, setting pipette speed for channel %d sets all other channels as well", head, channel)
		}
		for ch := 0; ch < adaptor.GetChannelCount(); ch++ {
			channels = append(channels, ch)
		}
	} else {
		channels = append(channels, channel)
	}

	outOfRange := make([]int, 0, len(channels))
	tRate := wunit.NewFlowRate(rate, "ml/min")
	minRate := make([]wunit.FlowRate, 0, len(channels))
	maxRate := make([]wunit.FlowRate, 0, len(channels))
	for ch := range channels {
		p := adaptor.GetParamsForChannel(ch)
		if tRate.GreaterThan(p.Maxspd) || tRate.LessThan(p.Minspd) {
			outOfRange = append(outOfRange, ch)
			minRate = append(minRate, p.Minspd)
			maxRate = append(maxRate, p.Maxspd)
		}
	}

	if len(outOfRange) > 0 && self.settings.IsPipetteSpeedWarningEnabled() {
		self.AddWarningf("SetPipetteSpeed", "Setting Head %d %s speed to %s is outside allowable range [%s:%s]",
			head, summariseChannels(outOfRange), tRate, summariseRates(minRate), summariseRates(maxRate))
	}

	return ret
}

//SetDriveSpeed - used
func (self *VirtualLiquidHandler) SetDriveSpeed(drive string, rate float64) driver.CommandStatus {
	return driver.CommandStatus{OK: true, Errorcode: driver.OK, Msg: "SETDRIVESPEED ACK"}
}

//Stop - unused
func (self *VirtualLiquidHandler) Stop() driver.CommandStatus {
	self.AddWarning("Stop", "Not yet implemented")
	return driver.CommandStatus{OK: true, Errorcode: driver.OK, Msg: "STOP ACK"}
}

//Go - unused
func (self *VirtualLiquidHandler) Go() driver.CommandStatus {
	self.AddWarning("Go", "Not yet implemented")
	return driver.CommandStatus{OK: true, Errorcode: driver.OK, Msg: "GO ACK"}
}

//Initialize - used
func (self *VirtualLiquidHandler) Initialize() driver.CommandStatus {
	if self.state.IsInitialized() {
		self.AddWarning("Initialize", "Call to initialize when robot is already initialized")
	}
	self.state.Initialize()
	return driver.CommandStatus{OK: true, Errorcode: driver.OK, Msg: "INITIALIZE ACK"}
}

//Finalize - used
func (self *VirtualLiquidHandler) Finalize() driver.CommandStatus {
	if !self.state.IsInitialized() {
		self.AddWarning("Finalize", "Call to finalize when robot is not inisialized")
	}
	if self.state.IsFinalized() {
		self.AddWarning("Finalize", "Call to finalize when robot is already finalized")
	}
	self.state.Finalize()
	return driver.CommandStatus{OK: true, Errorcode: driver.OK, Msg: "FINALIZE ACK"}
}

//Wait - used
func (self *VirtualLiquidHandler) Wait(time float64) driver.CommandStatus {
	self.AddWarning("Wait", "Not yet implemented")
	return driver.CommandStatus{OK: true, Errorcode: driver.OK, Msg: "WAIT ACK"}
}

//Mix - used
func (self *VirtualLiquidHandler) Mix(head int, volume []float64, platetype []string, cycles []int,
	multi int, what []string, blowout []bool) driver.CommandStatus {

	ret := driver.CommandStatus{OK: true, Errorcode: driver.OK, Msg: "MIX ACK"}

	//extend arguments - at some point shortform slices might become illegal
	if adaptor, err := self.GetAdaptorState(head); err == nil {
		volume = extend_floats(adaptor.GetChannelCount(), volume)
		platetype = extend_strings(adaptor.GetChannelCount(), platetype)
		cycles = extend_ints(adaptor.GetChannelCount(), cycles)
		what = extend_strings(adaptor.GetChannelCount(), what)
		blowout = extend_bools(adaptor.GetChannelCount(), blowout)
	}

	arg, err := self.validateLHArgs(head, multi, platetype, what, volume, map[string][]bool{
		"blowout": blowout,
	}, cycles)
	if err != nil {
		self.AddErrorf("Mix", "Invalid arguments - %s", err.Error())
		return ret
	}

	wells := self.getWellsBelow(0., arg.adaptor)

	describe := func() string {
		return fmt.Sprintf("mixing %s %s in %s of %s",
			summariseVolumes(arg.volumes),
			summariseCycles(cycles, arg.channels),
			summariseWells(wells, arg.channels),
			summarisePlates(wells, arg.channels))
	}

	//check tips exist
	no_tip := []int{}
	for _, i := range arg.channels {
		if !arg.adaptor.GetChannel(i).HasTip() {
			no_tip = append(no_tip, i)
		}
	}
	if len(no_tip) > 0 {
		self.AddErrorf("Mix", "While %s - no tip on %s", describe(), summariseChannels(no_tip))
		return ret
	}

	//check wells exist and their contents is ok
	no_well := []int{}
	for _, i := range arg.channels {
		v := wunit.NewVolume(volume[i], "ul")

		if wells[i] == nil {
			no_well = append(no_well, i)
		} else {
			if wells[i].Contents().GetType() != what[i] && self.settings.IsLiquidTypeWarningEnabled() {
				self.AddWarningf("Mix", "While %s - well contains %s not %s", describe(), wells[i].Contents().GetType(), what[i])
			}
			if wells[i].CurrentVolume().LessThan(v) {
				self.AddWarningf("Mix", "While %s - well only contains %s", describe(), wells[i].CurrentVolume())
			}
			if wtype.TypeOf(wells[i].Plate) != platetype[i] {
				self.AddWarningf("Mix", "While %s - plate \"%s\" is of type \"%s\", not \"%s\"",
					describe(), wtype.NameOf(wells[i].Plate), wtype.TypeOf(wells[i].Plate), platetype[i])
			}
		}
	}
	if len(no_well) > 0 {
		self.AddErrorf("Mix", "While %s - %s not in %s", describe(), summariseChannels(no_well), pWells(len(no_well)))
		return ret
	}

	//independece
	if !arg.adaptor.IsIndependent() {
		if !fElemsEqual(volume, arg.channels) {
			self.AddErrorf("Mix", "While %s - cannot manipulate different volumes with non-independent head", describe())
		}

		if !iElemsEqual(cycles, arg.channels) {
			self.AddErrorf("Mix", "While %s - cannot vary number of mix cycles with non-independent head", describe())
		}
	}

	//do the mixing
	for _, ch := range arg.channels {
		v := wunit.NewVolume(volume[ch], "ul")
		tip := arg.adaptor.GetChannel(ch).GetTip()

		//this is pretty pointless unless the tip already contained something
		//it also makes sure the tip.Contents().Name() is set properly
		for c := 0; c < cycles[ch]; c++ {
			com, err := wells[ch].RemoveVolume(v)
			if err != nil {
				self.AddErrorf("Mix", "Unexpected well error - %s", err.Error())
				continue
			}
			err = addComponent(tip, com)
			if err != nil {
				self.AddErrorf("Mix", "Unexpected well error - %s", err.Error())
				continue
			}
			com, err = tip.RemoveVolume(v)
			if err != nil {
				self.AddErrorf("Mix", "Unexpected tip error - %s", err.Error())
				continue
			}
			err = addComponent(wells[ch], com)
			if err != nil {
				self.AddErrorf("Mix", "Unexpected well error - %s", err.Error())
				continue
			}
		}
	}

	return ret
}

//ResetPistons - used
func (self *VirtualLiquidHandler) ResetPistons(head, channel int) driver.CommandStatus {
	self.AddWarning("ResetPistons", "Not yet implemented")
	return driver.CommandStatus{OK: true, Errorcode: driver.OK, Msg: "RESETPISTONS ACK"}
}

//These values correct for the Glison Driver offset and will eventually be removed
const (
	XCorrection = 14.38
	YCorrection = 11.24
)

//AddPlateTo - used
func (self *VirtualLiquidHandler) AddPlateTo(position string, plate interface{}, name string) driver.CommandStatus {

	ret := driver.CommandStatus{OK: true, Errorcode: driver.OK, Msg: "ADDPLATETO ACK"}

	if obj, ok := plate.(wtype.LHObject); ok {
		obj = obj.Duplicate(true)
		if n, nok := obj.(wtype.Named); nok && n.GetName() != name {
			self.AddWarningf("AddPlateTo", "Object name(=%s) doesn't match argument name(=%s)", n.GetName(), name)
		}

		if tb, ok := obj.(*wtype.LHTipbox); ok {
			//check that the height of the tips is greater than the height of the tipbox
			if tb.GetSize().Z >= (tb.TipZStart+tb.Tiptype.GetSize().Z) && self.settings.IsTipboxCheckEnabled() {
				self.AddWarningf("AddPlateTo",
					"Tipbox \"%s\" is taller than the tips it holds (%.2fmm > %.2fmm), disabling tipbox collision detection",
					tb.GetName(), tb.GetSize().Z, tb.TipZStart+tb.Tiptype.GetSize().Z)
				self.settings.EnableTipboxCollision(false)
			}
		}

		//check that the wells are within the bounds of the plate
		if plate, ok := obj.(*wtype.LHPlate); ok {
			//apply the well position correction
			plate.WellXStart += XCorrection
			plate.WellYStart += YCorrection
			corr := wtype.Coordinates{X: XCorrection, Y: YCorrection, Z: 0.0}
			for _, w := range plate.Wellcoords {
				w.SetOffset(w.Bounds.GetPosition().Add(corr)) //nolint
			}

			plateSize := plate.GetSize()
			wellOff := plate.GetWellOffset()
			wellLim := wellOff.Add(plate.GetWellSize())

			if wellOff.X < 0.0 || wellOff.Y < 0.0 || wellOff.Z < 0.0 {
				self.AddWarningf("AddPlateTo", "position \"%s\" : invalid plate type \"%s\" has negative well offsets %v",
					position, wtype.TypeOf(plate), wellOff)
			}

			overSpill := wtype.Coordinates{
				X: math.Max(wellLim.X-plateSize.X, 0.0),
				Y: math.Max(wellLim.Y-plateSize.Y, 0.0),
				Z: math.Max(wellLim.Z-plateSize.Z, 0.0),
			}

			if overSpill.Z > 0.0 {
				self.AddWarningf("AddPlateTo", "position \"%s\" : invalid plate type \"%s\" : increasing height by %0.1f mm to match well height",
					position, wtype.TypeOf(plate), overSpill.Z)
				plateSize.Z += overSpill.Z
				plate.Bounds.SetSize(plateSize)
			}

			if overSpill.X > 0.0 || overSpill.Y > 0.0 {
				self.AddWarningf("AddPlateTo", "position \"%s\" : invalid plate type \"%s\" wells extend beyond plate bounds by %s",
					position, wtype.TypeOf(plate), overSpill.StringXY())
			}
		}

		if err := self.state.GetDeck().SetChild(position, obj); err != nil {
			self.AddError("AddPlateTo", err.Error())
			return ret
		}

	} else {
		self.AddErrorf("AddPlateTo", "Couldn't add object of type %T to %s", plate, position)
	}

	return ret
}

//RemoveAllPlates - used
func (self *VirtualLiquidHandler) RemoveAllPlates() driver.CommandStatus {
	deck := self.state.GetDeck()
	for _, name := range deck.GetSlotNames() {
		if err := deck.Clear(name); err != nil {
			self.AddError("RemoveAllPlates", err.Error())
		}
	}
	return driver.CommandStatus{OK: true, Errorcode: driver.OK, Msg: "REMOVEALLPLATES ACK"}
}

//RemovePlateAt - unused
func (self *VirtualLiquidHandler) RemovePlateAt(position string) driver.CommandStatus {
	self.AddWarning("RemovePlateAt", "Not yet implemented")
	return driver.CommandStatus{OK: true, Errorcode: driver.OK, Msg: "REMOVEPLATEAT ACK"}
}

//SetPositionState - unused
func (self *VirtualLiquidHandler) SetPositionState(position string, state driver.PositionState) driver.CommandStatus {
	self.AddWarning("SetPositionState", "Not yet implemented")
	return driver.CommandStatus{OK: true, Errorcode: driver.OK, Msg: "SETPOSITIONSTATE ACK"}
}

//GetCapabilites - used
func (self *VirtualLiquidHandler) GetCapabilities() (liquidhandling.LHProperties, driver.CommandStatus) {
	self.AddWarning("SetPositionState", "Not yet implemented")
	return liquidhandling.LHProperties{}, driver.CommandStatus{OK: true, Errorcode: driver.OK, Msg: "GETCAPABILITIES ACK"}
}

//GetCurrentPosition - unused
func (self *VirtualLiquidHandler) GetCurrentPosition(head int) (string, driver.CommandStatus) {
	self.AddWarning("GetCurrentPosition", "Not yet implemented")
	return "", driver.CommandStatus{OK: true, Errorcode: driver.OK, Msg: "GETCURRNETPOSITION ACK"}
}

//GetPositionState - unused
func (self *VirtualLiquidHandler) GetPositionState(position string) (string, driver.CommandStatus) {
	self.AddWarning("GetPositionState", "Not yet implemented")
	return "", driver.CommandStatus{OK: true, Errorcode: driver.OK, Msg: "GETPOSITIONSTATE ACK"}
}

//GetHeadState - unused
func (self *VirtualLiquidHandler) GetHeadState(head int) (string, driver.CommandStatus) {
	self.AddWarning("GetHeadState", "Not yet implemented")
	return "I'm fine thanks, how are you?", driver.CommandStatus{OK: true, Errorcode: driver.OK, Msg: "GETHEADSTATE ACK"}
}

//GetStatus - unused
func (self *VirtualLiquidHandler) GetStatus() (driver.Status, driver.CommandStatus) {
	self.AddWarning("GetStatus", "Not yet implemented")
	return driver.Status{}, driver.CommandStatus{OK: true, Errorcode: driver.OK, Msg: "GETSTATUS ACK"}
}

//UpdateMetaData - used
func (self *VirtualLiquidHandler) UpdateMetaData(props *liquidhandling.LHProperties) driver.CommandStatus {
	self.AddWarning("UpdateMetaData", "Not yet implemented")
	return driver.CommandStatus{OK: true, Errorcode: driver.OK, Msg: "UPDATEMETADATA ACK"}
}

//UnloadHead - unused
func (self *VirtualLiquidHandler) UnloadHead(param int) driver.CommandStatus {
	self.AddWarning("UnloadHead", "Not yet implemented")
	return driver.CommandStatus{OK: true, Errorcode: driver.OK, Msg: "UNLOADHEAD ACK"}
}

//LoadHead - unused
func (self *VirtualLiquidHandler) LoadHead(param int) driver.CommandStatus {
	self.AddWarning("LoadHead", "Not yet implemented")
	return driver.CommandStatus{OK: true, Errorcode: driver.OK, Msg: "LOADHEAD ACK"}
}

//Lights On - not implemented in compositerobotinstruction
func (self *VirtualLiquidHandler) LightsOn() driver.CommandStatus {
	return driver.CommandStatus{OK: true, Errorcode: driver.OK, Msg: "LIGHTSON ACK"}
}

//Lights Off - notimplemented in compositerobotinstruction
func (self *VirtualLiquidHandler) LightsOff() driver.CommandStatus {
	return driver.CommandStatus{OK: true, Errorcode: driver.OK, Msg: "LIGHTSOFF ACK"}
}

//LoadAdaptor - notimplemented in CRI
func (self *VirtualLiquidHandler) LoadAdaptor(param int) driver.CommandStatus {
	self.AddWarning("LoadAdaptor", "Not yet implemented")
	return driver.CommandStatus{OK: true, Errorcode: driver.OK, Msg: "LOADADAPTOR ACK"}
}

//UnloadAdaptor - notimplemented in CRI
func (self *VirtualLiquidHandler) UnloadAdaptor(param int) driver.CommandStatus {
	self.AddWarning("UnloadAdaptor", "Not yet implemented")
	return driver.CommandStatus{OK: true, Errorcode: driver.OK, Msg: "UNLOADADAPTOR ACK"}
}

//Open - notimplemented in CRI
func (self *VirtualLiquidHandler) Open() driver.CommandStatus {
	self.AddWarning("Open", "Not yet implemented")
	return driver.CommandStatus{OK: true, Errorcode: driver.OK, Msg: "OPEN ACK"}
}

//Close - notimplement in CRI
func (self *VirtualLiquidHandler) Close() driver.CommandStatus {
	self.AddWarning("Close", "Not yet implemented")
	return driver.CommandStatus{OK: true, Errorcode: driver.OK, Msg: "CLOSE ACK"}
}

//Message - unused
func (self *VirtualLiquidHandler) Message(level int, title, text string, showcancel bool) driver.CommandStatus {
	self.AddWarning("Message", "Not yet implemented")
	return driver.CommandStatus{OK: true, Errorcode: driver.OK, Msg: "MESSAGE ACK"}
}

//GetOutputFile - used, but not in instruction stream
func (self *VirtualLiquidHandler) GetOutputFile() ([]byte, driver.CommandStatus) {
	self.AddWarning("GetOutputFile", "Not yet implemented")
	return []byte("You forgot to say 'please'"), driver.CommandStatus{OK: true, Errorcode: driver.OK, Msg: "GETOUTPUTFILE ACK"}
}
