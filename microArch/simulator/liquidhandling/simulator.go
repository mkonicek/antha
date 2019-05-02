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

	"github.com/pkg/errors"

	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/microArch/driver"
	"github.com/antha-lang/antha/microArch/driver/liquidhandling"
	"github.com/antha-lang/antha/microArch/simulator"
	"github.com/antha-lang/antha/utils"
)

const arbitraryZOffset = 4.0

// Simulate a liquid handler Driver
type VirtualLiquidHandler struct {
	errorHistory       [][]LiquidhandlingError
	instructionHistory []liquidhandling.TerminalRobotInstruction
	errors             []LiquidhandlingError
	state              *RobotState
	settings           *SimulatorSettings
	lastMove           string
	lastTarget         wtype.LHObject
	properties         *liquidhandling.LHProperties
	objectByID         map[string]wtype.LHObject // map from object ID to the object used internally
}

//coneRadius hardcoded radius to assume for cones
const coneRadius = 3.6

//Create a new VirtualLiquidHandler which mimics an LHDriver
func NewVirtualLiquidHandler(props *liquidhandling.LHProperties, settings *SimulatorSettings) (*VirtualLiquidHandler, error) {
	vlh := VirtualLiquidHandler{
		errors:             make([]LiquidhandlingError, 0),
		errorHistory:       make([][]LiquidhandlingError, 0),
		instructionHistory: make([]liquidhandling.TerminalRobotInstruction, 0),
		objectByID:         make(map[string]wtype.LHObject, len(props.Positions)),
	}

	if settings == nil {
		vlh.settings = DefaultSimulatorSettings()
	} else {
		vlh.settings = settings
	}

	if err := vlh.validateProperties(props); err != nil {
		return nil, errors.Wrap(err, "building virtual liquid handler")
	}
	vlh.state = NewRobotState()

	//add the adaptors
	for _, assembly := range props.HeadAssemblies {
		vlh.state.AddAdaptorGroup(NewAdaptorGroup(assembly))
	}

	//Make the deck
	deck := wtype.NewLHDeck("simulated deck", props.Mnfr, props.Model)
	for name, pos := range props.Positions {
		//size not given un LHProperties, assuming standard 96well size
		deck.AddSlot(name, pos.Location, pos.Size)
		//deck.SetSlotAccepts(name, "riser")
	}

	for _, name := range props.Preferences.Tipboxes {
		deck.SetSlotAccepts(name, "tipbox")
	}
	for _, name := range props.Preferences.Inputs {
		deck.SetSlotAccepts(name, "plate")
	}
	for _, name := range props.Preferences.Outputs {
		deck.SetSlotAccepts(name, "plate")
	}
	for _, name := range props.Preferences.Tipwastes {
		deck.SetSlotAccepts(name, "tipwaste")
	}

	vlh.state.SetDeck(deck)

	vlh.properties = props

	return &vlh, nil
}

//Simulate simulate the list of instructions
func (self *VirtualLiquidHandler) Simulate(instructions []liquidhandling.TerminalRobotInstruction) error {

	self.resetState()

	for _, ins := range instructions {
		err := ins.(liquidhandling.TerminalRobotInstruction).OutputTo(self)
		if err != nil {
			return errors.Wrap(err, "while writing instructions to virtual device")
		}

		self.saveState(ins)
	}

	return nil
}

func (self *VirtualLiquidHandler) getState() *RobotState {
	if self == nil {
		return nil
	}
	return self.state
}

func (self *VirtualLiquidHandler) GetLastMove() string {
	if self == nil {
		return ""
	}
	return self.lastMove
}

func (self *VirtualLiquidHandler) GetLastTarget() wtype.LHObject {
	if self == nil {
		return nil
	}
	return self.lastTarget
}

func (self *VirtualLiquidHandler) GetProperties() *liquidhandling.LHProperties {
	if self == nil {
		return nil
	}
	return self.properties
}

//CountErrors
func (self *VirtualLiquidHandler) CountErrors() int {
	ret := 0
	for _, state := range self.errorHistory {
		ret += len(state)
	}
	return ret + len(self.errors)
}

//GetErrors
func (self *VirtualLiquidHandler) GetErrors() []simulator.SimulationError {
	ret := make([]simulator.SimulationError, 0, self.CountErrors())
	for _, state := range self.errorHistory {
		for _, err := range state {
			ret = append(ret, err)
		}
	}

	for _, err := range self.errors {
		ret = append(ret, err)
	}
	return ret
}

//GetFirstError get the first error that's at least as bad as minimum severity
func (self *VirtualLiquidHandler) GetFirstError(minimumSeverity simulator.ErrorSeverity) simulator.SimulationError {

	for _, state := range self.errorHistory {
		for _, err := range state {
			if err.Severity() >= minimumSeverity {
				return err
			}
		}
	}

	return nil
}

// ------------------------------------------------------------------------------- Useful Utilities

func (self *VirtualLiquidHandler) resetState() {
	self.errorHistory = make([][]LiquidhandlingError, 0)
	self.instructionHistory = make([]liquidhandling.TerminalRobotInstruction, 0)
	self.errors = make([]LiquidhandlingError, 0)
}

func (self *VirtualLiquidHandler) popLastState() (liquidhandling.TerminalRobotInstruction, []LiquidhandlingError) {
	if len(self.errorHistory) <= 0 {
		return nil, nil
	}

	err := self.errorHistory[len(self.errorHistory)-1]
	ins := self.instructionHistory[len(self.instructionHistory)-1]
	self.errorHistory = self.errorHistory[:len(self.errorHistory)-1]
	self.instructionHistory = self.instructionHistory[:len(self.instructionHistory)-1]

	return ins, err
}

func (self *VirtualLiquidHandler) saveState(ins liquidhandling.TerminalRobotInstruction) {
	for _, err := range self.errors {
		if mErr, ok := err.(mutableLHError); ok {
			mErr.setInstruction(len(self.errorHistory), ins)
		}
	}
	self.instructionHistory = append(self.instructionHistory, ins)
	self.errorHistory = append(self.errorHistory, self.errors)
	self.errors = make([]LiquidhandlingError, 0)
}

func (self *VirtualLiquidHandler) addLHError(err LiquidhandlingError) {
	self.errors = append(self.errors, err)
}

func (self *VirtualLiquidHandler) AddInfo(message string) {
	self.addLHError(NewGenericError(self.state, simulator.SeverityInfo, message))
}

func (self *VirtualLiquidHandler) AddInfof(format string, a ...interface{}) {
	self.addLHError(NewGenericErrorf(self.state, simulator.SeverityInfo, format, a...))
}

func (self *VirtualLiquidHandler) AddWarning(message string) {
	self.addLHError(NewGenericError(self.state, simulator.SeverityWarning, message))
}

func (self *VirtualLiquidHandler) AddWarningf(format string, a ...interface{}) {
	self.addLHError(NewGenericErrorf(self.state, simulator.SeverityWarning, format, a...))
}

func (self *VirtualLiquidHandler) AddError(message string) {
	self.addLHError(NewGenericError(self.state, simulator.SeverityError, message))
}

func (self *VirtualLiquidHandler) AddErrorf(format string, a ...interface{}) {
	self.addLHError(NewGenericErrorf(self.state, simulator.SeverityError, format, a...))
}

func (self *VirtualLiquidHandler) validateProperties(props *liquidhandling.LHProperties) error {

	//check a property
	check_prop := func(l []string, name string) error {
		//all locations defined
		for _, loc := range l {
			if !props.Exists(loc) {
				return errors.Errorf(`unknown location "%s" found in %s preferences`, loc, name)
			}
		}
		return nil
	}
	return utils.ErrorSlice{
		check_prop(props.Preferences.Tipboxes, "tipbox"),
		check_prop(props.Preferences.Inputs, "input"),
		check_prop(props.Preferences.Outputs, "output"),
		check_prop(props.Preferences.Tipwastes, "tipwaste"),
		check_prop(props.Preferences.Wastes, "waste"),
		check_prop(props.Preferences.Washes, "wash"),
	}.Pack()
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
func (self *VirtualLiquidHandler) testTipArgs(channels []int, head int, platetype, position, well []string) bool {
	//head should exist
	adaptor, err := self.GetAdaptorState(head)
	if err != nil {
		self.AddError(err.Error())
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
		self.AddErrorf("Unknown channel \"%v\"", bad_channels[0])
		ret = false
	} else if len(bad_channels) > 1 {
		self.AddErrorf("Unknown channels \"%v\"", strings.Join(bad_channels, "\",\""))
		ret = false
	}
	if len(dup_channels) == 1 {
		self.AddErrorf("Channel%v appears more than once", dup_channels[0])
		ret = false
	} else if len(dup_channels) == 1 {
		self.AddErrorf("Channels {%s} appear more than once", strings.Join(dup_channels, "\",\""))
		ret = false
	}

	if err := self.testSliceLength(map[string]int{
		"platetype": len(platetype),
		"position":  len(position),
		"well":      len(well)},
		n_channels); err != nil {

		self.AddError(err.Error())
		ret = false
	}

	if ret {
		for i := range platetype {
			if contains(i, channels) {
				if platetype[i] == "" && well[i] == "" && position[i] == "" {
					self.AddErrorf("Command given for channel %d, but no platetype, well or position given", i)
					return false
				}
			} else if len(channels) > 0 { //if channels are empty we'll infer it later
				if !(platetype[i] == "" && well[i] == "" && position[i] == "") {
					self.AddWarningf("No command for channel %d, but platetype, well or position given", i)
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
func (self *VirtualLiquidHandler) getTargetPosition(adaptorName string, channelIndex int, deckposition, platetype string, wc wtype.WellCoords, ref wtype.WellReference) (wtype.Coordinates3D, bool) {
	ret := wtype.Coordinates3D{}

	target, ok := self.state.GetDeck().GetChild(deckposition)
	if !ok {
		self.AddErrorf("Unknown location \"%s\"", deckposition)
		return ret, false
	}
	if target == nil {
		self.AddErrorf("No object found at position %s", deckposition)
		return ret, false
	}

	if (platetype != wtype.TypeOf(target)) &&
		(platetype != wtype.NameOf(target)) {
		self.AddWarningf("Object found at %s was type \"%s\", named \"%s\", not \"%s\" as expected",
			deckposition, wtype.TypeOf(target), wtype.NameOf(target), platetype)
	}

	addr, ok := target.(wtype.Addressable)
	if !ok {
		self.AddErrorf("Object \"%s\" at \"%s\" is not addressable", wtype.NameOf(target), deckposition)
		return ret, false
	}

	if !addr.AddressExists(wc) {
		self.AddErrorf("Request for well %s in object \"%s\" at \"%s\" which is of size [%dx%d]",
			wc.FormatA1(), wtype.NameOf(target), deckposition, addr.NRows(), addr.NCols())
		return ret, false
	}

	self.lastTarget = addr.GetChildByAddress(wc)

	ret, ok = addr.WellCoordsToCoords(wc, ref)
	if !ok {
		//since we already checked that the address exists, this must be a bad reference
		self.AddErrorf("Object type %s at %s doesn't support reference \"%s\"",
			wtype.TypeOf(target), deckposition, ref)
		return ret, false
	}

	if targetted, ok := target.(wtype.Targetted); ok {
		targetOffset := targetted.GetTargetOffset(adaptorName, channelIndex)
		ret = ret.Add(targetOffset)
	}

	return ret, true
}

// GetWellAt return the internal model of the well at the given location, or nil if not found
func (self *VirtualLiquidHandler) GetWellAt(pl wtype.PlateLocation) *wtype.LHWell {
	if plate, ok := self.objectByID[pl.ID].(*wtype.LHPlate); ok {
		w, _ := plate.WellAt(pl.Coords)
		return w
	}
	return nil
}

func (self *VirtualLiquidHandler) getWellsBelow(height float64, adaptor *AdaptorState) []*wtype.LHWell {
	tip_pos := make([]wtype.Coordinates3D, adaptor.GetChannelCount())
	wells := make([]*wtype.LHWell, adaptor.GetChannelCount())
	size := wtype.Coordinates3D{X: 0, Y: 0, Z: height}
	deck := self.state.GetDeck()
	for i := 0; i < adaptor.GetChannelCount(); i++ {
		if ch := adaptor.GetChannel(i); ch.HasTip() {
			tip_pos[i] = ch.GetAbsolutePosition().Subtract(wtype.Coordinates3D{X: 0., Y: 0., Z: ch.GetTip().GetEffectiveHeight()})
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

func makeOffsets(Xs, Ys, Zs []float64) []wtype.Coordinates3D {
	ret := make([]wtype.Coordinates3D, len(Xs))
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
	ret := driver.CommandOk()

	//get the adaptor
	adaptor, err := self.GetAdaptorState(head)
	if err != nil {
		self.AddError(err.Error())
		return ret
	}

	//only support a single deckposition or platetype
	deckposition, err := getSingle(deckpositionS)
	if err != nil {
		self.AddErrorf("invalid argument deckposition: %s", err.Error())
	}
	platetype, err := getSingle(platetypeS)
	if err != nil {
		self.AddErrorf("invalid argument platetype: %s", err.Error())
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

		self.AddError(err.Error())
		return ret
	}

	refs, err := convertReferences(reference)
	if err != nil {
		self.AddErrorf("invalid argument reference: %s", err.Error())
	}

	//get slice of well coords
	wc, err := convertWellCoords(wellcoords)
	if err != nil {
		self.AddErrorf("invalid argument wellcoords: %s", err.Error())
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
		self.AddWarning("ignoring blank move command: no wellcoords specified")
		return ret
	}

	//combine floats into wtype.Coordinates
	offsets := makeOffsets(offsetX, offsetY, offsetZ)

	//find the coordinates of each explicitly requested position
	coords := make([]wtype.Coordinates3D, adaptor.GetChannelCount())
	for _, ch := range channels {
		c, ok := self.getTargetPosition(adaptor.GetName(), ch, deckposition, platetype, wc[ch], refs[ch])
		if !ok {
			return ret
		}
		coords[ch] = c
		coords[ch] = coords[ch].Add(offsets[ch])
		//if there's a tip, raise the coortinates to the top of the tip to take account of it
		if tip := adaptor.GetChannel(ch).GetTip(); tip != nil {
			coords[ch] = coords[ch].Add(wtype.Coordinates3D{X: 0., Y: 0., Z: tip.GetEffectiveHeight()})
		}
	}

	target, ok := self.state.GetDeck().GetChild(deckposition)
	if !ok {
		self.AddErrorf("unable to get object at position \"%s\"", deckposition)
	}

	describe := func() string {
		return fmt.Sprintf("head %d %s to %s of %s@%s at position %s",
			head, summariseChannels(channels), summariseWellReferences(channels, offsetZ, refs), wtype.HumanizeWellCoords(wc), wtype.NameOf(target), deckposition)
	}

	//store a description of the move for posterity (and future errors)
	self.lastMove = describe()

	//find the head location
	//for now, assuming that the relative position of the first explicitly provided channel and the head stay
	//the same. This seems sensible for the Glison, but might turn out not to be how other robots with independent channels work
	origin := coords[channels[0]].Subtract(adaptor.GetChannel(channels[0]).GetRelativePosition())

	//fill in implicit locations
	for _, ch := range implicitChannels {
		coords[ch] = origin.Add(adaptor.GetChannel(ch).GetRelativePosition())
	}

	//Get the locations of each channel relative to the head
	rel_coords := make([]wtype.Coordinates3D, adaptor.GetChannelCount())
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
			self.AddErrorf("%s: requires moving %s relative to non-independent head",
				describe(), summariseChannels(moved))
			return ret
		}
	}

	//check that there are no tips loaded on any other heads
	if err = assertNoTipsOnOthersInGroup(adaptor); err != nil {
		self.AddErrorf("%s: cannot move head %d while %s", describe(), head, err.Error())
	}

	//move the head to the new position
	err = adaptor.SetPosition(origin)
	if err != nil {
		self.AddErrorf("%s: %s", describe(), err.Error())
	}
	for i, rc := range rel_coords {
		adaptor.GetChannel(i).SetRelativePosition(rc)
	}

	//check for collisions in the new location
	if err := assertNoCollisionsInGroup(self.settings, adaptor, nil, 0.0); err != nil {
		err.SetInstructionDescription(describe())
		self.addLHError(err)
	}
	return ret
}

//Move raw - not yet implemented in compositerobotinstruction
func (self *VirtualLiquidHandler) MoveRaw(head int, x, y, z float64) driver.CommandStatus {
	self.AddWarning("not yet implemented")
	return driver.CommandOk()
}

//Aspirate - used
func (self *VirtualLiquidHandler) Aspirate(volume []float64, overstroke []bool, head int, multi int,
	platetype []string, what []string, llf []bool) driver.CommandStatus {

	ret := driver.CommandOk()

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
		self.AddErrorf("invalid Arguments: %s", err.Error())
		return ret
	}

	describe := func() string {
		return fmt.Sprintf("%s of %s to head %d %s",
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
		self.AddErrorf("%s: missing %s on %s", describe(), pTips(len(tip_missing)), summariseChannels(tip_missing))
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
				if c := wells[ch].Contents(); !c.IsZero() && arg.adaptor.GetChannel(ch).HasTip() {
					self.AddErrorf(
						"%s: channel %d will inadvertantly aspirate %s from well %s as head is not independent",
						describe(), ch, c.Name(), wells[ch].GetName())
				}
			}
		}
		if different {
			self.AddErrorf("%s: channels cannot aspirate different volumes in non-independent head", describe())
			return ret
		}
	}

	//check liquid type
	for i := range arg.channels {
		if wells[i] == nil { //we'll catch this later
			continue
		}
		if wells[i].Contents().GetType() != what[i] && self.settings.IsLiquidTypeWarningEnabled() {
			self.AddWarningf("%s: well %s contains %s, not %s",
				describe(), wells[i].GetName(), wells[i].Contents().GetType(), what[i])
		}
	}

	//check total volumes taken from each unique well
	volumeTakenByWell := make(map[*wtype.LHWell]float64, len(wells))
	for i := 0; i < len(wells); i++ {
		if wells[i] != nil {
			volumeTakenByWell[wells[i]] = volumeTakenByWell[wells[i]] + volume[i]
		}
	}
	for well, v := range volumeTakenByWell {
		volume := wunit.NewVolume(v, "ul")
		if volume.GreaterThan(well.CurrentVolume().PlusEpsilon()) {
			// we've completely exhausted the well, raise an error
			self.AddErrorf("%s: taking %s from %s which contains only %s working plus %s residual volume (possibly a side effect of ANTHA-2704, try using a plate with a larger residual volume)",
				describe(), volume, well.GetName(), well.CurrentVolume(), well.ResidualVolume())

		} else if volume.GreaterThan(well.CurrentWorkingVolume().PlusEpsilon()) {
			// the total volume taken from this well is greater than the working volume available in the well
			// this will appear as though liquid has been aspirated from the residual volume
			//
			// This happens when several "bites" are taken from the same well during GetComponents and the
			// unaccounted for carry volume causes the well volume to be exhausted: see ANTHA-2704
			//
			// Once that issue has been resolved, this warning should become fatal
			self.AddWarningf("%s: taking %s from %s which contains only %s working volume, possible aspiration of residual (see ANTHA-2704)",
				describe(), volume, well.GetName(), well.CurrentWorkingVolume())
		}
	}

	//move liquid
	no_well := []int{}
	for _, i := range arg.channels {
		if wells[i] == nil {
			no_well = append(no_well, i)
		} else {

			aspVol := wunit.NewVolume(volume[i], "ul")
			tip := arg.adaptor.GetChannel(i).GetTip()
			tipVol := tip.CurrentVolume()
			tipVol.Add(aspVol)

			if tipVol.GreaterThan(tip.MaxVol.PlusEpsilon()) {
				self.AddErrorf("%s: channel %d contains %s, command exceeds maximum volume %s",
					describe(), i, tip.CurrentVolume(), tip.MaxVol)
			} else if c, err := wells[i].RemoveVolume(aspVol); err != nil {
				self.AddErrorf("%s: unexpected well error \"%s\"", describe(), err.Error())
			} else {
				err := tip.AddComponent(c)
				if tipVol.LessThan(tip.MinVol.MinusEpsilon()) {
					// ignore the error returned from AddComponent and add a warning instead
					self.AddWarningf("%s: minimum tip volume is %s", describe(), tip.MinVol)
				} else if err != nil {
					self.AddErrorf("%s: unexpected tip error \"%s\"", describe(), err.Error())
				}
			}
		}
	}

	// silently remove the carry volume for each well - this may leave the well volume lower than the residual volume
	// do this after removing everything else so carry volumes don't interfere with multichannel aspiration
	for _, well := range wells {
		if well != nil {
			well.RemoveCarry(self.properties.CarryVolume())
		}
	}

	if len(no_well) > 0 {
		self.addLHError(NewTipsNotInWellError(self, describe(), no_well))
	}

	return ret
}

//Dispense - used
func (self *VirtualLiquidHandler) Dispense(volume []float64, blowout []bool, head int, multi int,
	platetype []string, what []string, llf []bool) driver.CommandStatus {

	ret := driver.CommandOk()

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
		self.AddErrorf("invalid arguments: %s", err.Error())
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
		self.AddErrorf("%s: no well within %s below %s on %s",
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
		self.AddErrorf("%s: no %s loaded on %s",
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
			} else if arg.adaptor.GetChannel(ch).HasTip() {
				//a non-explicitly requested is loaded. If the tip has stuff in it, it'll get dispensed as well
				if c := arg.adaptor.GetChannel(ch).GetTip().Contents(); !c.IsZero() {
					extra = append(extra, ch)
				}
			}
		}
		if different {
			self.AddErrorf("%s: channels cannot dispense different volumes in non-independent head", describe())
			return ret
		} else if len(extra) > 0 {
			self.AddErrorf("%s: must also dispense %s from %s as head is not independent",
				describe(), summariseVolumes(arg.volumes), summariseChannels(extra))
			return ret
		}
	}

	//check liquid type
	for i := range arg.channels {
		if tip := arg.adaptor.GetChannel(i).GetTip(); tip != nil {
			if tip.Contents().GetType() != what[i] && self.settings.IsLiquidTypeWarningEnabled() {
				self.AddWarningf("%s: channel %d contains %s, not %s",
					describe(), i, tip.Contents().GetType(), what[i])
			}
		}
	}

	//for each blowout channel
	for _, i := range arg.channels {
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
		self.AddWarningf("%s: overfilling %s %s to %s of %s max volume",
			describe(), pWells(len(overfullWells)), summarisePlateWells(wells, overfullWells), summariseVolumes(finalVolumes), summariseVolumes(maxVolumes))
	}

	//dispense
	for _, i := range arg.channels {
		v := wunit.NewVolume(volume[i], "ul")
		tip := arg.adaptor.GetChannel(i).GetTip()

		if wells[i] != nil {
			if _, tw := wells[i].Plate.(*wtype.LHTipwaste); tw {
				self.AddWarningf("%s: dispensing to tipwaste", describe())
			}
		}

		if v.GreaterThan(tip.CurrentWorkingVolume()) {
			v = tip.CurrentWorkingVolume()
			if !blowout[i] {
				//a bit strange
				self.AddWarningf("%s: tip on channel %d contains only %s, but blowout flag is false",
					describe(), i, tip.CurrentWorkingVolume())
			}
		}
		if c, err := tip.RemoveVolume(v); err != nil {
			self.AddErrorf("%s: unexpected tip error \"%s\"", describe(), err.Error())
		} else if err := wells[i].AddComponent(c); err != nil {
			self.AddErrorf("%s: unexpected well error \"%s\"", describe(), err.Error())
		}
	}

	return ret
}

//LoadTips - used
func (self *VirtualLiquidHandler) LoadTips(channels []int, head, multi int,
	platetypeS, positionS, well []string) driver.CommandStatus {
	ret := driver.CommandOk()
	deck := self.state.GetDeck()

	//get the adaptor
	adaptor, err := self.GetAdaptorState(head)
	if err != nil {
		self.AddError(err.Error())
		return ret
	}
	n_channels := adaptor.GetChannelCount()

	//extend arg slices
	platetypeS = extend_strings(n_channels, platetypeS)
	positionS = extend_strings(n_channels, positionS)
	well = extend_strings(n_channels, well)

	//check that the command is valid
	if !self.testTipArgs(channels, head, platetypeS, positionS, well) {
		return ret
	}

	//get the individual position
	position, err := getSingle(positionS)
	if err != nil {
		self.AddErrorf("invalid argument position: %s", err.Error())
	}
	platetype, err := getSingle(platetypeS)
	if err != nil {
		self.AddErrorf("invalid argument platetype: %s", err.Error())
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
		self.AddErrorf("invalid argument well: %s", err.Error())
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
		self.AddErrorf("invalid argument well: couldn't parse \"%s\" for %s", strings.Join(values, "\", \""), summariseChannels(invalidWells))
	}

	//get the actual tipbox
	var tipbox *wtype.LHTipbox
	if o, ok := deck.GetChild(position); !ok {
		self.AddErrorf("unknown location \"%s\"", position)
		return ret
	} else if o == nil {
		self.AddErrorf("can't load tips from empty position \"%s\"", position)
		return ret
	} else if tipbox, ok = o.(*wtype.LHTipbox); !ok {
		self.AddErrorf("can't load tips from %s \"%s\" found at position \"%s\"",
			wtype.ClassOf(o), wtype.NameOf(o), position)
		return ret
	}
	if tipbox == nil {
		self.AddErrorf("unexpected nil tipbox at position \"%s\"", position)
		return ret
	}

	describe := func() string {
		return fmt.Sprintf("from %s@%s at position \"%s\" to head %d %s", wtype.HumanizeWellCoords(wc), tipbox.GetName(), position, head, summariseChannels(channels))
	}

	if multi != len(channels) {
		self.AddErrorf("%s: multi should equal %d, not %d",
			describe(), len(channels), multi)
		return ret
	}

	//check that channels we want to load to are empty
	if tipFound := checkTipPresence(false, adaptor, channels); len(tipFound) != 0 {
		self.AddErrorf("%s: %s already loaded on head %d %s",
			describe(), pTips(len(tipFound)), head, summariseChannels(tipFound))
		return ret
	}

	//check that there aren't any tips loaded on any other heads
	if err := assertNoTipsOnOthersInGroup(adaptor); err != nil {
		self.AddErrorf("%s: while %s", describe(), err.Error())
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
			self.AddErrorf("%s: unexpected error: %s", describe(), err.Error())
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
				self.AddWarningf(
					"%s: Well coordinates for channel %d not specified, assuming %s from adaptor location",
					describe(), i, wc[i].FormatA1())
			}
		}

		if !tipbox.AddressExists(wc[i]) {
			self.AddErrorf("%s: request for tip at %s in tipbox of size [%dx%d]",
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
		self.AddErrorf("%s: no %s at %s",
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
		tip_p := tips[ch].GetPosition().Add(wtype.Coordinates3D{X: 0.5 * tip_s.X, Y: 0.5 * tip_s.Y, Z: tip_s.Z})
		ch_p := adaptor.GetChannel(ch).GetAbsolutePosition()
		delta := ch_p.Subtract(tip_p)
		if xy := delta.AbsXY(); xy > 0.5 {
			misaligned = append(misaligned, ch)
			target = append(target, wc[ch])
			amount = append(amount, fmt.Sprintf("%v", xy))
		}
		z_off[ch] = delta.Z
		if delta.Z < 0. {
			self.AddErrorf("%s: channel is %.1f below tip", describe(), -delta.Z)
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
		self.AddErrorf("%s: %s %s misaligned with %s at %s by %smm%s",
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
			self.AddErrorf("%s: distance between channels and tips varies from %v to %v mm in non-independent head",
				describe(), zo_min, zo_max)
			return ret
		}
		if err := assertNoCollisionsInGroup(self.settings, adaptor, channels, zo_max+0.5); err != nil {
			err.SetInstructionDescription(describe())
			self.addLHError(err)
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

	//undo the last command that moved us into position
	//assumption is that last command was a move...
	ins, _ := self.popLastState()

	var ret driver.CommandStatus
	loadedChannels := make([]int, 0, len(channels))

	for i, chunk := range tipChunks {
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

		self.Move(positionS, wellcoords, reference, offsetXY, offsetXY, offsetZ, platetypeS, head)
		if i == 0 {
			//save the state of the first move, so the instruction counting matches
			self.saveState(ins)
		}
		ret = self.LoadTips(channelsToLoad, head, width, platetypeS, positionS, wellcoords)
		loadedChannels = append(loadedChannels, channelsToLoad...)
	}

	return ret
}

//UnloadTips - used
func (self *VirtualLiquidHandler) UnloadTips(channels []int, head, multi int,
	platetype, position, well []string) driver.CommandStatus {
	ret := driver.CommandOk()

	//get the adaptor
	adaptor, err := self.GetAdaptorState(head)
	if err != nil {
		self.AddError(err.Error())
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
			self.AddWarningf("'channel' argument empty and no tips are loaded to head %d, ignoring", head)
		} else if self.settings.IsAutoChannelWarningEnabled() {
			self.AddWarningf("'channel' argument empty, unloading all tips (%s)", summariseChannels(channels))
		}
	}

	//check that RobotState won't crash
	if !self.testTipArgs(channels, head, platetype, position, well) {
		return ret
	}

	if multi != len(channels) {
		self.AddWarningf("While unloading %s from %s, multi should equal %d, not %d",
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
		self.AddWarningf("No tip present at Head%d channel%s to eject", head, missing[0])
	} else if len(missing) > 0 {
		self.AddWarningf("No tips present on Head%d channels %s to eject", head, strings.Join(missing, ","))
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
			self.AddErrorf("Cannot unload tips from head%d %s without unloading %s from %s (head isn't independent)",
				head, summariseChannels(channels), pTips(len(extra)), summariseChannels(extra))
			return ret
		}
	}

	for _, ch := range channels {
		//get the target
		if target, ok := deck.GetChild(position[ch]); !ok {
			self.AddErrorf("Unknown deck position \"%s\"", position[ch])
			break
		} else if target == nil {
			self.AddErrorf("Cannot unload to empty deck location \"%s\"", position[ch])
			break
		} else if addr, ok := target.(wtype.Addressable); !ok {
			self.AddErrorf("Cannot unload tips to %s \"%s\" at location %s",
				wtype.ClassOf(target), wtype.NameOf(target), position[ch])
		} else {
			//get the location of the channel
			ch_pos := adaptor.GetChannel(ch).GetAbsolutePosition()
			//parse the wellcoords
			wc := wtype.MakeWellCoords(well[ch])
			if wc.IsZero() {
				self.AddErrorf("Cannot parse well coordinates \"%s\"", well[ch])
				break
			}
			if !addr.AddressExists(wc) {
				self.AddErrorf("Cannot unload to address %s in %s \"%s\" size [%dx%d]",
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
					self.AddErrorf("Cannot unload to tipbox \"%s\" %s, tip already present there",
						target.GetName(), wc.FormatA1())
				} else if delta.AbsXY() > 0.25 {
					self.AddErrorf("Head%d channel%d misaligned from tipbox \"%s\" %s by %.2fmm",
						head, ch, target.GetName(), wc.FormatA1(), delta.AbsXY())
				} else if delta.Z > target.GetSize().Z/2. {
					self.AddWarningf("Ejecting tip from Head%d channel%d to tipbox \"%s\" %s from height of %.2fmm",
						head, ch, target.GetName(), wc.FormatA1(), delta.Z)
				} else {
					target.PutTip(wc, adaptor.GetChannel(ch).UnloadTip())
				}

			case *wtype.LHTipwaste:
				//put the tip in the tipwaste
				if child == nil {
					//I don't think this should happen, but it would be embarressing to segfault...
					self.AddWarningf("Tipwaste \"%s\" well %s was nil, cannot check head alignment",
						target.GetName(), wc.FormatA1())
					adaptor.GetChannel(ch).UnloadTip()
				} else if max_delta := child.GetSize(); delta.X > max_delta.X || delta.Y > max_delta.Y {
					self.AddErrorf("Cannot unload, head%d channel%d is not above tipwaste \"%s\"",
						head, ch, target.GetName())
				} else if target.SpaceLeft() <= 0 {
					self.AddErrorf("Cannot unload tip to overfull tipwaste \"%s\", contains %d tips",
						target.GetName(), target.Contents)
				} else {
					target.DisposeNum(1)
					adaptor.GetChannel(ch).UnloadTip()
				}
			default:
				self.AddErrorf("Cannot unload tips to %s \"%s\" at location %s",
					wtype.ClassOf(target), wtype.NameOf(target), position[ch])
			}
		}
	}

	return ret
}

//SetPipetteSpeed - used
func (self *VirtualLiquidHandler) SetPipetteSpeed(head, channel int, rate float64) driver.CommandStatus {
	ret := driver.CommandOk()

	adaptor, err := self.GetAdaptorState(head)
	if err != nil {
		self.AddError(err.Error())
		return ret
	}

	channels := make([]int, 0, adaptor.GetChannelCount())
	if channel < 0 || !adaptor.IsIndependent() {
		if channel >= 0 {
			self.AddWarningf("Head %d is not independent, setting pipette speed for channel %d sets all other channels as well", head, channel)
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
		self.AddWarningf("Setting Head %d %s speed to %s is outside allowable range [%s:%s]",
			head, summariseChannels(outOfRange), tRate, summariseRates(minRate), summariseRates(maxRate))
	}

	return ret
}

// SetDriveSpeed sets the speed at which the head will move.
// Drive should be one of "X", "Y", "Z", rate is expressed in mm/s.
// XXX Bug: it is currently not possible to select which head assembly should be affected
// currently assume we're talking about group zero
func (self *VirtualLiquidHandler) SetDriveSpeed(drive string, rate float64) driver.CommandStatus {
	ret := driver.CommandOk()

	//Assume we're talking about adaptor group zero
	groupNumber := 0

	v := wunit.NewVelocity(rate, "mm/s")

	describe := func() string {
		return fmt.Sprintf("while setting head group %d drive %s speed to %v", groupNumber, drive, v)
	}

	if group, err := self.state.GetAdaptorGroup(groupNumber); err != nil {
		self.AddErrorf("%s: %s", describe(), err.Error())
		return ret
	} else if axis, err := wunit.AxisFromString(drive); err != nil {
		self.AddErrorf("%s: %s", describe(), err.Error())
		return ret
	} else if err := group.SetDriveSpeed(axis, v); err != nil {
		self.AddErrorf("%s: %s", describe(), err.Error())
		return ret
	}
	return ret
}

//Stop - unused
func (self *VirtualLiquidHandler) Stop() driver.CommandStatus {
	self.AddWarning("not yet implemented")
	return driver.CommandOk()
}

//Go - unused
func (self *VirtualLiquidHandler) Go() driver.CommandStatus {
	self.AddWarning("not yet implemented")
	return driver.CommandOk()
}

//Initialize - used
func (self *VirtualLiquidHandler) Initialize() driver.CommandStatus {
	if self.state.IsInitialized() {
		self.AddWarning("Call to initialize when robot is already initialized")
	}
	self.state.Initialize()
	return driver.CommandOk()
}

//Finalize - used
func (self *VirtualLiquidHandler) Finalize() driver.CommandStatus {
	if !self.state.IsInitialized() {
		self.AddWarning("Call to finalize when robot is not inisialized")
	}
	if self.state.IsFinalized() {
		self.AddWarning("Call to finalize when robot is already finalized")
	}
	self.state.Finalize()
	return driver.CommandOk()
}

//Wait - used
func (self *VirtualLiquidHandler) Wait(time float64) driver.CommandStatus {
	if time < 0.0 {
		self.AddWarning("waiting for negative time")
	}
	return driver.CommandOk()
}

//Mix - used
func (self *VirtualLiquidHandler) Mix(head int, volume []float64, platetype []string, cycles []int,
	multi int, what []string, blowout []bool) driver.CommandStatus {

	ret := driver.CommandOk()

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
		self.AddErrorf("Invalid arguments - %s", err.Error())
		return ret
	}

	wells := self.getWellsBelow(0., arg.adaptor)

	describe := func() string {
		return fmt.Sprintf("%s %s in %s of %s",
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
		self.AddErrorf("%s: no tip on %s", describe(), summariseChannels(no_tip))
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
				self.AddWarningf("%s: well contains %s not %s", describe(), wells[i].Contents().GetType(), what[i])
			}
			if wells[i].CurrentVolume().LessThan(v) {
				self.AddWarningf("%s: well only contains %s", describe(), wells[i].CurrentVolume())
			}
			if wtype.TypeOf(wells[i].Plate) != platetype[i] {
				self.AddWarningf("%s: plate \"%s\" is of type \"%s\", not \"%s\"",
					describe(), wtype.NameOf(wells[i].Plate), wtype.TypeOf(wells[i].Plate), platetype[i])
			}
		}
	}
	if len(no_well) > 0 {
		self.addLHError(NewTipsNotInWellError(self, describe(), no_well))
		return ret
	}

	//independece
	if !arg.adaptor.IsIndependent() {
		if !fElemsEqual(volume, arg.channels) {
			self.AddErrorf("%s: cannot manipulate different volumes with non-independent head", describe())
		}

		if !iElemsEqual(cycles, arg.channels) {
			self.AddErrorf("%s: cannot vary number of mix cycles with non-independent head", describe())
		}
	}

	// tips should be empty
	nonEmptyChannels := make([]int, 0, len(arg.channels))
	nonEmptyVolumes := make([]float64, 0, len(arg.channels))
	for _, ch := range arg.channels {
		if tip := arg.adaptor.GetChannel(ch).GetTip(); !tip.IsEmpty() {
			nonEmptyChannels = append(nonEmptyChannels, ch)
			nonEmptyVolumes = append(nonEmptyVolumes, tip.CurrentVolume().MustInStringUnit("ul").RawValue())
		}
	}
	if len(nonEmptyChannels) > 0 {
		self.AddErrorf("%s: mixing when tips on %s contain %s", describe(), summariseChannels(nonEmptyChannels), summariseVolumes(nonEmptyVolumes))
	}

	return ret
}

//ResetPistons - used
func (self *VirtualLiquidHandler) ResetPistons(head, channel int) driver.CommandStatus {
	self.AddWarning("not yet implemented")
	return driver.CommandOk()
}

//These values correct for the Glison Driver offset and will eventually be removed
const (
	XCorrection = 14.38
	YCorrection = 11.24
)

//AddPlateTo - used
func (self *VirtualLiquidHandler) AddPlateTo(position string, plate interface{}, name string) driver.CommandStatus {

	ret := driver.CommandOk()

	if original, ok := plate.(wtype.LHObject); ok {
		obj := original.Duplicate(true)
		if n, nok := obj.(wtype.Named); nok && n.GetName() != name {
			self.AddWarningf("Object name(=%s) doesn't match argument name(=%s)", n.GetName(), name)
		}

		if tb, ok := obj.(*wtype.LHTipbox); ok {
			//check that the height of the tips is greater than the height of the tipbox
			if tb.GetSize().Z >= (tb.TipZStart+tb.Tiptype.GetSize().Z) && self.settings.IsTipboxCheckEnabled() {
				self.AddWarningf(
					"Tipbox \"%s\" is taller than the tips it holds (%.2fmm > %.2fmm), disabling tipbox collision detection",
					tb.GetName(), tb.GetSize().Z, tb.TipZStart+tb.Tiptype.GetSize().Z)
				self.settings.EnableTipboxCollision(false)
			}
		}

		//check that the wells are within the bounds of the plate
		if plate, ok := obj.(*wtype.Plate); ok {
			plateSize := plate.GetSize()
			wellOff := plate.GetWellOffset()
			wellLim := wellOff.Add(plate.GetWellSize())

			if wellOff.X < 0.0 || wellOff.Y < 0.0 || wellOff.Z < 0.0 {
				self.AddWarningf("position \"%s\": invalid plate type \"%s\" has negative well offsets %v",
					position, wtype.TypeOf(plate), wellOff)
			}

			overSpill := wtype.Coordinates3D{
				X: math.Max(wellLim.X-plateSize.X, 0.0),
				Y: math.Max(wellLim.Y-plateSize.Y, 0.0),
				Z: math.Max(wellLim.Z-plateSize.Z, 0.0),
			}

			if overSpill.Z > 0.0 {
				self.AddWarningf("position \"%s\": invalid plate type \"%s\": wells extend above plate, reducing well height by %0.1f mm",
					position, wtype.TypeOf(plate), overSpill.Z)
				wellSize := plate.Welltype.GetSize()
				wellSize.Z -= overSpill.Z
				plate.Welltype.Bounds.Size = wellSize
				for _, well := range plate.Wellcoords {
					well.Bounds.Size = wellSize
				}
			}

			if overSpill.X > 0.0 || overSpill.Y > 0.0 {
				self.AddWarningf("position \"%s\": invalid plate type \"%s\" wells extend beyond plate bounds by %s",
					position, wtype.TypeOf(plate), overSpill.StringXY())
			}
		}

		if err := self.state.GetDeck().SetChild(position, obj); err != nil {
			self.AddError(err.Error())
			return ret
		}

		self.objectByID[wtype.IDOf(obj)] = obj

	} else {
		self.AddErrorf("Couldn't add object of type %T to %s", plate, position)
	}

	return ret
}

//RemoveAllPlates - used
func (self *VirtualLiquidHandler) RemoveAllPlates() driver.CommandStatus {
	deck := self.state.GetDeck()
	for _, name := range deck.GetSlotNames() {
		if err := deck.Clear(name); err != nil {
			self.AddError(err.Error())
		}
	}
	return driver.CommandOk()
}

//RemovePlateAt - unused
func (self *VirtualLiquidHandler) RemovePlateAt(position string) driver.CommandStatus {
	self.AddWarning("not yet implemented")
	return driver.CommandOk()
}

//SetPositionState - unused
func (self *VirtualLiquidHandler) SetPositionState(position string, state driver.PositionState) driver.CommandStatus {
	self.AddWarning("not yet implemented")
	return driver.CommandOk()
}

//GetCapabilites - used
func (self *VirtualLiquidHandler) GetCapabilities() (liquidhandling.LHProperties, driver.CommandStatus) {
	self.AddWarning("not yet implemented")
	return liquidhandling.LHProperties{}, driver.CommandOk()
}

//GetCurrentPosition - unused
func (self *VirtualLiquidHandler) GetCurrentPosition(head int) (string, driver.CommandStatus) {
	self.AddWarning("not yet implemented")
	return "", driver.CommandOk()
}

//GetPositionState - unused
func (self *VirtualLiquidHandler) GetPositionState(position string) (string, driver.CommandStatus) {
	self.AddWarning("not yet implemented")
	return "", driver.CommandOk()
}

//GetHeadState - unused
func (self *VirtualLiquidHandler) GetHeadState(head int) (string, driver.CommandStatus) {
	self.AddWarning("not yet implemented")
	return "I'm fine thanks, how are you?", driver.CommandOk()
}

//GetStatus - unused
func (self *VirtualLiquidHandler) GetStatus() (driver.Status, driver.CommandStatus) {
	self.AddWarning("not yet implemented")
	return driver.Status{}, driver.CommandOk()
}

//UpdateMetaData - used
func (self *VirtualLiquidHandler) UpdateMetaData(props *liquidhandling.LHProperties) driver.CommandStatus {
	self.AddWarning("not yet implemented")
	return driver.CommandOk()
}

//UnloadHead - unused
func (self *VirtualLiquidHandler) UnloadHead(param int) driver.CommandStatus {
	self.AddWarning("not yet implemented")
	return driver.CommandOk()
}

//LoadHead - unused
func (self *VirtualLiquidHandler) LoadHead(param int) driver.CommandStatus {
	self.AddWarning("not yet implemented")
	return driver.CommandOk()
}

//Lights On - not implemented in compositerobotinstruction
func (self *VirtualLiquidHandler) LightsOn() driver.CommandStatus {
	return driver.CommandOk()
}

//Lights Off - notimplemented in compositerobotinstruction
func (self *VirtualLiquidHandler) LightsOff() driver.CommandStatus {
	return driver.CommandOk()
}

//LoadAdaptor - notimplemented in CRI
func (self *VirtualLiquidHandler) LoadAdaptor(param int) driver.CommandStatus {
	self.AddWarning("not yet implemented")
	return driver.CommandOk()
}

//UnloadAdaptor - notimplemented in CRI
func (self *VirtualLiquidHandler) UnloadAdaptor(param int) driver.CommandStatus {
	self.AddWarning("not yet implemented")
	return driver.CommandOk()
}

//Open - notimplemented in CRI
func (self *VirtualLiquidHandler) Open() driver.CommandStatus {
	self.AddWarning("not yet implemented")
	return driver.CommandOk()
}

//Close - notimplement in CRI
func (self *VirtualLiquidHandler) Close() driver.CommandStatus {
	self.AddWarning("not yet implemented")
	return driver.CommandOk()
}

//Message - unused
func (self *VirtualLiquidHandler) Message(level int, title, text string, showcancel bool) driver.CommandStatus {
	self.AddWarning("not yet implemented")
	return driver.CommandOk()
}

//GetOutputFile - used, but not in instruction stream
func (self *VirtualLiquidHandler) GetOutputFile() ([]byte, driver.CommandStatus) {
	self.AddWarning("not yet implemented")
	return []byte("You forgot to say 'please'"), driver.CommandOk()
}

func (self *VirtualLiquidHandler) DriverType() ([]string, error) {
	return []string{"antha.mixer.v1.Mixer", "PreciousVirtualLiquidHandler"}, nil
}
