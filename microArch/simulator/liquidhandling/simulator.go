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
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/microArch/driver"
	"github.com/antha-lang/antha/microArch/driver/liquidhandling"
	"github.com/antha-lang/antha/microArch/simulator"
	"math"
	"sort"
	"strings"
)

func pTips(N int) string {
	if N == 1 {
		return "tip"
	}
	return "tips"
}

func pWells(N int) string {
	if N == 1 {
		return "well"
	}
	return "wells"
}

func summariseWell2Channel(well []string, channels []int) string {
	ret := make([]string, 0, len(channels))
	for ch := range channels {
		ret = append(ret, fmt.Sprintf("%s->channel%v", well[ch], ch))
	}
	return strings.Join(ret, ", ")
}

func summariseChannels(channels []int) string {
	if len(channels) == 1 {
		return fmt.Sprintf("channel %d", channels[0])
	}
	sch := make([]string, 0, len(channels))
	for _, ch := range channels {
		sch = append(sch, fmt.Sprintf("%d", ch))
	}
	return fmt.Sprintf("channels %s", strings.Join(sch, ","))
}

func summariseWellCoords(wellCoords []wtype.WellCoords) string {
	ss := make([]string, 0, len(wellCoords))
	for _, wc := range wellCoords {
		ss = append(ss, wc.FormatA1())
	}
	return summariseStrings(ss)
}

func summariseVolumes(vols []float64) string {
	equal := true
	for _, v := range vols {
		if v != vols[0] {
			equal = false
			break
		}
	}

	if equal {
		return wunit.NewVolume(vols[0], "ul").String()
	}

	s_vols := make([]string, len(vols))
	for i, v := range vols {
		s_vols[i] = wunit.NewVolume(v, "ul").String()
		s_vols[i] = s_vols[i][:len(s_vols[i])-3]
	}
	return fmt.Sprintf("{%s} ul", strings.Join(s_vols, ","))
}

func summariseRates(rates []wunit.FlowRate) string {
	asString := make([]string, 0, len(rates))
	for _, r := range rates {
		asString = append(asString, r.String())
	}
	return summariseStrings(asString)
}

func summariseStrings(s []string) string {
	if countUnique(s, true) == 1 {
		return firstNonEmpty(s)
	}
	return "{" + strings.Join(s, ",") + "}"
}

func summariseCycles(cycles []int, elems []int) string {
	if iElemsEqual(cycles, elems) {
		if cycles[0] == 1 {
			return "once"
		} else {
			return fmt.Sprintf("%d times", cycles[0])
		}
	}
	sc := make([]string, 0, len(elems))
	for _, i := range elems {
		sc = append(sc, fmt.Sprintf("%d", cycles[i]))
	}
	return fmt.Sprintf("{%s} times", strings.Join(sc, ","))
}

func summariseWells(wells []*wtype.LHWell, elems []int) string {
	w := make([]string, 0, len(elems))
	for _, i := range elems {
		w = append(w, wells[i].Crds.FormatA1())
	}
	uw := getUnique(w, true)

	if len(uw) == 1 {
		return fmt.Sprintf("well %s", uw[0])
	}
	return fmt.Sprintf("wells %s", strings.Join(uw, ","))
}

func summarisePlates(wells []*wtype.LHWell, elems []int) string {
	p := make([]string, 0, len(elems))
	for _, i := range elems {
		p = append(p, wtype.NameOf(wells[i].Plate))
	}
	up := getUnique(p, true)

	if len(up) == 1 {
		return fmt.Sprintf("plate \"%s\"", up[0])
	}
	return fmt.Sprintf("plates \"%s\"", strings.Join(up, "\",\""))

}

//summarisePlateWells list wells for each plate preserving order
func summarisePlateWells(wells []*wtype.LHWell, elems []int) string {
	var lastWell *wtype.LHWell
	currentChunk := make([]string, 0, len(elems))
	var chunkedWells [][]string
	var plateNames []string

	for _, i := range elems {
		well := wells[i]
		if lastWell != nil && lastWell.GetParent() != well.GetParent() {
			chunkedWells = append(chunkedWells, currentChunk)
			currentChunk = make([]string, 0, len(elems))
			plateNames = append(plateNames, wtype.NameOf(well.GetParent()))
		}
		lastWell = well
		if well != nil {
			currentChunk = append(currentChunk, well.Crds.FormatA1())
		}
	}
	chunkedWells = append(chunkedWells, currentChunk)
	plateNames = append(plateNames, wtype.NameOf(lastWell.GetParent()))

	var ret []string
	for i, name := range plateNames {
		if len(chunkedWells[i]) > 1 {
			ret = append(ret, fmt.Sprintf("{%s}@%s", strings.Join(chunkedWells[i], ","), name))
		} else if len(chunkedWells[i]) == 1 {
			ret = append(ret, fmt.Sprintf("%s@%s", chunkedWells[i][0], name))
		}
	}

	if len(ret) == 0 {
		return "nil"
	}

	return strings.Join(ret, ", ")
}

func iElemsEqual(sl []int, elems []int) bool {
	for _, i := range elems {
		if sl[i] != sl[elems[0]] {
			return false
		}
	}
	return true
}

func fElemsEqual(sl []float64, elems []int) bool {
	for _, i := range elems {
		if sl[i] != sl[elems[0]] {
			return false
		}
	}
	return true
}

func extend_ints(l int, sl []int) []int {
	if len(sl) < l {
		r := make([]int, l)
		copy(r, sl)
		return r
	}
	return sl
}

func extend_floats(l int, sl []float64) []float64 {
	if len(sl) < l {
		r := make([]float64, l)
		copy(r, sl)
		return r
	}
	return sl
}

func extend_strings(l int, sl []string) []string {
	if len(sl) < l {
		r := make([]string, l)
		copy(r, sl)
		return r
	}
	return sl
}

func extend_bools(l int, sl []bool) []bool {
	if len(sl) < l {
		r := make([]bool, l)
		copy(r, sl)
		return r
	}
	return sl
}

type adaptorCollision struct {
	channel int
	objects []wtype.LHObject
}

type adaptorCollisions []adaptorCollision

func (self adaptorCollisions) String() string {
	channels := make([]string, 0, len(self))
	objects := []wtype.LHObject{}

	seen := func(o wtype.LHObject) bool {
		for _, O := range objects {
			if o == O {
				return true
			}
		}
		return false
	}

	for _, ac := range self {
		channels = append(channels, fmt.Sprintf("%d", ac.channel))
		for _, obj := range ac.objects {
			if !seen(obj) {
				objects = append(objects, obj)
			}
		}
	}

	s_obj := make([]string, 0, len(objects))
	for _, o := range objects {
		s_obj = append(s_obj, fmt.Sprintf("%s \"%s\"", wtype.ClassOf(o), wtype.NameOf(o)))
	}

	if len(self) == 1 {
		return fmt.Sprintf("channel %s collides with %s", channels[0], strings.Join(s_obj, " and "))
	}
	return fmt.Sprintf("channels %s collide with %s", strings.Join(channels, ","), strings.Join(s_obj, " and "))
}

// Simulate a liquid handler Driver
type VirtualLiquidHandler struct {
	simulator.ErrorReporter
	state    *RobotState
	settings *SimulatorSettings
}

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
	for _, head := range props.Heads {
		p := head.Adaptor.Params
		//9mm spacing currently hardcoded.
		//At some point we'll either need to fetch this from the driver or
		//infer it from the type of tipboxes/plates accepted
		spacing := wtype.Coordinates{X: 0, Y: 0, Z: 0}
		if p.Orientation == wtype.LHVChannel {
			spacing.Y = 9.
		} else if p.Orientation == wtype.LHHChannel {
			spacing.X = 9.
		}
		vlh.state.AddAdaptor(NewAdaptorState(head.Adaptor.Name, p.Independent, p.Multi, spacing, p))
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

//getAdaptorState
func (self *VirtualLiquidHandler) getAdaptorState(h int) (*AdaptorState, error) {
	if h < 0 || h >= self.state.GetNumberOfAdaptors() {
		return nil, fmt.Errorf("Unknown head %d", h)
	}
	return self.state.GetAdaptor(h), nil
}

func (self *VirtualLiquidHandler) GetAdaptorState(head int) *AdaptorState {
	return self.state.GetAdaptor(head)
}

func (self *VirtualLiquidHandler) GetObjectAt(slot string) wtype.LHObject {
	child, _ := self.state.GetDeck().GetChild(slot)
	return child
}

//testTipArgs check that load/unload tip arguments are valid insofar as they won't crash in RobotState
func (self *VirtualLiquidHandler) testTipArgs(f_name string, channels []int, head int, platetype, position, well []string) bool {
	//head should exist
	adaptor, err := self.getAdaptorState(head)
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

	ret.adaptor, err = self.getAdaptorState(head)
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
func (self *VirtualLiquidHandler) getTargetPosition(fname, adaptorName string, channelIndex int, deckposition, platetype, well string, ref wtype.WellReference) (wtype.Coordinates, bool) {
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

	wc := wtype.MakeWellCoords(well)
	if wc.IsZero() {
		self.AddErrorf(fname, "Couldn't parse well \"%s\"", well)
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
			tip_pos[i] = ch.GetAbsolutePosition().Subtract(wtype.Coordinates{X: 0., Y: 0., Z: ch.GetTip().GetSize().Z})
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
func (self *VirtualLiquidHandler) Move(deckposition []string, wellcoords []string, reference []int,
	offsetX, offsetY, offsetZ []float64, platetype []string,
	head int) driver.CommandStatus {
	ret := driver.CommandStatus{OK: true, Errorcode: driver.OK, Msg: "MOVE ACK"}

	//get the adaptor
	adaptor, err := self.getAdaptorState(head)
	if err != nil {
		self.AddError("Move", err.Error())
		return ret
	}

	//extend args
	deckposition = extend_strings(adaptor.GetChannelCount(), deckposition)
	wellcoords = extend_strings(adaptor.GetChannelCount(), wellcoords)
	reference = extend_ints(adaptor.GetChannelCount(), reference)
	offsetX = extend_floats(adaptor.GetChannelCount(), offsetX)
	offsetY = extend_floats(adaptor.GetChannelCount(), offsetY)
	offsetZ = extend_floats(adaptor.GetChannelCount(), offsetZ)
	platetype = extend_strings(adaptor.GetChannelCount(), platetype)

	//check slice length
	if err := self.testSliceLength(map[string]int{
		"deckposition": len(deckposition),
		"wellcoords":   len(wellcoords),
		"reference":    len(reference),
		"offsetX":      len(offsetX),
		"offsetY":      len(offsetY),
		"offsetZ":      len(offsetZ),
		"plate_type":   len(platetype)},
		adaptor.GetChannelCount()); err != nil {

		self.AddError("Move", err.Error())
		return ret
	}

	refs := make([]wtype.WellReference, adaptor.GetChannelCount())
	for i, r := range reference {
		switch r {
		case 0:
			refs[i] = wtype.BottomReference
		case 1:
			refs[i] = wtype.TopReference
		case 2:
			refs[i] = wtype.LiquidReference
		default:
			self.AddErrorf("Move", "Invalid reference %d", r)
			return ret
		}
	}

	//find the coordinates of each explicitly requested position
	coords := make([]wtype.Coordinates, adaptor.GetChannelCount())
	offsets := makeOffsets(offsetX, offsetY, offsetZ)
	explicit := make([]bool, adaptor.GetChannelCount())
	exp_count := 0
	for i := range deckposition {
		if deckposition[i] == "" {
			if wellcoords[i] != "" {
				self.AddWarningf("Move", "deckposition was blank, but well was \"%s\"", wellcoords[i])
			}
			if platetype[i] != "" {
				self.AddWarningf("Move", "deckposition was blank, but platetype was \"%s\"", platetype[i])
			}
			explicit[i] = false
		} else {
			c, ok := self.getTargetPosition("Move", adaptor.GetName(), i, deckposition[i], platetype[i], wellcoords[i], refs[i])
			if !ok {
				return ret
			}
			coords[i] = c
			coords[i] = coords[i].Add(offsets[i])
			//if there's a tip, take account of it
			if tip := adaptor.GetChannel(i).GetTip(); tip != nil {
				coords[i] = coords[i].Add(wtype.Coordinates{X: 0., Y: 0., Z: tip.GetSize().Z})
			}
			explicit[i] = true
			exp_count++
		}
	}
	if exp_count == 0 {
		self.AddWarning("Move", "Ignoring blank move command")
	}

	//find the head location, origin
	origin := wtype.Coordinates{}
	//for now, assuming that the relative position of the first explicitly provided channel and the head stay
	//the same. This seems sensible for the Glison, but might turn out not to be how other robots with independent channels work
	for i, c := range coords {
		if explicit[i] {
			origin = c.Subtract(adaptor.GetChannel(i).GetRelativePosition())
			break
		}
	}

	//fill in implicit locations
	for i := range coords {
		if !explicit[i] {
			coords[i] = origin.Add(adaptor.GetChannel(i).GetRelativePosition())
		}
	}

	//Get relative locations
	rel_coords := make([]wtype.Coordinates, adaptor.GetChannelCount())
	for i := range coords {
		rel_coords[i] = coords[i].Subtract(origin)
	}

	//check that the requested position is possible given the head/adaptor capabilities
	if !adaptor.IsIndependent() {
		//i.e. the channels can't move relative to each other or the head, so relative locations must remain the same
		moved := []string{}
		for i, rc := range rel_coords {
			//check that adaptor relative position remains the same
			//arbitrary 0.01mm to avoid numerical instability
			if rc.Subtract(adaptor.GetChannel(i).GetRelativePosition()).Abs() > 0.01 {
				moved = append(moved, fmt.Sprintf("%d", i))
			}
		}
		if len(moved) > 0 {
			//get slice of well coords
			wc := make([]wtype.WellCoords, len(wellcoords))
			for i := range wellcoords {
				wc[i] = wtype.MakeWellCoords(wellcoords[i])
			}
			self.AddErrorf("Move", "Non-independent head '%d' can't move adaptors to \"%s\" positions %s, layout mismatch",
				head, strings.Join(getUnique(platetype, true), "\",\""), wtype.HumanizeWellCoords(wc))
			return ret
		}
	}

	//check for collisions in the new location
	for ch, rc := range rel_coords {
		pos := origin.Add(rc)
		obj := self.state.GetDeck().GetPointIntersections(pos)
		in_well := false
		for _, o := range obj {
			if _, ok := o.(*wtype.LHWell); ok {
				in_well = true
			}
		}
		if !in_well && len(obj) > 0 {
			o_str := make([]string, len(obj))
			for i, o := range obj {
				o_str[i] = wtype.NameOf(o)
			}
			self.AddErrorf("Move", "Cannot move channel %d to (\"%s\", %s, %s) + (%.1f,%.1f,%.1f)mm as this collides with %s\n",
				ch, deckposition[ch], wellcoords[ch], refs[ch], offsetX[ch], offsetY[ch], offsetZ[ch], strings.Join(o_str, " and "))
			return ret
		}
	}

	//update the head position accordingly
	adaptor.SetPosition(origin)
	for i, rc := range rel_coords {
		adaptor.GetChannel(i).SetRelativePosition(rc)
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
	if adaptor, err := self.getAdaptorState(head); err == nil {
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
	for i := 0; i < len(wells); i++ {
		if wells[i] == nil {
			continue
		}
		if _, ok := uniqueWells[wells[i].ID]; !ok {
			uniqueWells[wells[i].ID] = wells[i]
			uniqueWellVolumes[wells[i].ID] = 0.0
		}
		uniqueWellVolumes[wells[i].ID] += volume[i]
	}
	for id, well := range uniqueWells {
		v := wunit.NewVolume(uniqueWellVolumes[id], "ul")
		if d := wunit.SubtractVolumes(v, well.CurrentWorkingVolume()); v.GreaterThan(well.CurrentWorkingVolume()) && !d.IsZero() {
			self.AddErrorf("Aspirate", "While %s - well %s only contains %s working volume",
				describe(), well.GetName(), well.CurrentWorkingVolume())
		}
	}
	if self.HasError() {
		return ret
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
	if adaptor, err := self.getAdaptorState(head); err == nil {
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
	adaptor, err := self.getAdaptorState(head)
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

	//make well coords
	wc := make([]wtype.WellCoords, n_channels)
	for i := range well {
		wc[i] = wtype.MakeWellCoords(well[i])
	}

	//get the individual position
	if countUnique(positionS, true) != 1 {
		self.AddErrorf("LoadTips", "invalid position slice \"%v\", only one position supported", positionS)
		return ret
	}
	position := firstNonEmpty(positionS)

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

	if self.settings.GetTipTrackingBehaviour() == TrilutionTipTracking {
		//refil the tipbox if there aren't enough tips to service the instruction
		if !tipbox.HasEnoughTips(multi) {
			tipbox.Refill()
		}
		//HJK: we should also check that we're picking up tip in the way trilution is know to (e.g. splitting over rows etc)
		//but this requires more information about the geometry of the head and so on

	}

	describe := func() string {
		return fmt.Sprintf("from %s@%s at position \"%s\" to head %d %s", summariseWellCoords(wc), tipbox.GetName(), position, head, summariseChannels(channels))
	}

	//check that channels we want to load to are empty
	if tipFound := checkTipPresence(false, adaptor, channels); len(tipFound) != 0 {
		self.AddErrorf("LoadTips", "%s : %s already loaded to %s",
			describe(), pTips(len(tipFound)), summariseChannels(tipFound))
		return ret
	}

	if len(channels) == 0 {
		for ch, pt := range platetypeS {
			if pt != "" {
				channels = append(channels, ch)
			}
		}

		//best to order the channels sensibly
		sort.Ints(channels)

		if len(channels) == 0 {
			self.AddWarning("LoadTips", "'channel' argument empty and no platetype specified ignoring")
			return ret
		} else if self.settings.IsAutoChannelWarningEnabled() {
			self.AddWarningf("LoadTips", "%s : channels weren't specified in instruction, inferring %s from platetype", describe(), summariseChannels(channels))
		}

		//check if multi is wrong
		if multi != len(channels) {
			self.AddErrorf("LoadTips", "%s : 'channel' argument inferred as %s, but 'multi' is %d",
				describe(),
				summariseChannels(channels),
				multi)
			return ret
		}
	}
	if multi != len(channels) {
		self.AddErrorf("LoadTips", "%s : multi should equal %d, not %d",
			describe(), len(channels), multi)
		return ret
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
			describe(), pTips(len(missingTips)), summariseWellCoords(missingTips))
		return ret
	}

	//check alignment
	z_off := make([]float64, n_channels)
	misaligned := []string{}
	target := []string{}
	amount := []string{}
	for _, ch := range channels {
		tip_s := tips[ch].GetSize()
		tip_p := tips[ch].GetPosition().Add(wtype.Coordinates{X: 0.5 * tip_s.X, Y: 0.5 * tip_s.Y, Z: tip_s.Z})
		ch_p := adaptor.GetChannel(ch).GetAbsolutePosition()
		delta := ch_p.Subtract(tip_p)
		if xy := delta.AbsXY(); xy > 0.5 {
			misaligned = append(misaligned, fmt.Sprintf("%d", ch))
			target = append(target, wc[ch].FormatA1())
			amount = append(amount, fmt.Sprintf("%v", xy))
		}
		z_off[ch] = delta.Z
		if delta.Z < 0. {
			self.AddErrorf("LoadTips", "Request to load tip at location %s to channel %d at %s, channel is %.1f below tip", tip_p, ch, ch_p, -delta.Z)
			return ret
		}
	}
	if len(misaligned) == 1 {
		self.AddErrorf("LoadTips", "Channel %s is misaligned with tip at %s by %smm",
			misaligned[0], target[0], amount[0])
		return ret
	} else if len(misaligned) > 1 {
		self.AddErrorf("LoadTips", "Channels %s are misaligned with tips at %s by %s mm respectively",
			strings.Join(misaligned, ","), strings.Join(target, ","), strings.Join(amount, ","))
		return ret
	}

	//if not independent, check there are no other tips in the way
	if !adaptor.IsIndependent() {
		collisions := adaptorCollisions{}
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
			self.AddErrorf("LoadTips", "Distance between channels and tips varies from %v to %v mm in non-independent head",
				zo_min, zo_max)
			return ret
		}
		for i := 0; i < adaptor.GetChannelCount(); i++ {
			if contains(i, channels) {
				continue
			}
			ch_pos := adaptor.GetChannel(i).GetAbsolutePosition()
			size := wtype.Coordinates{X: 0., Y: 0., Z: zo_max + 0.5}
			box := wtype.NewBBox(ch_pos.Subtract(size), size)
			objects := deck.GetBoxIntersections(*box)
			//filter out tipboxes if we're meant to be ignoring them
			//(hack to prevent dubious tipbox geometry messing this up)
			if !self.settings.IsTipboxCollisionEnabled() {
				no_tipboxes := objects[:0]
				for _, o := range objects {
					if _, ok := o.(*wtype.LHTipbox); !ok {
						no_tipboxes = append(no_tipboxes, o)
					}
				}
				objects = no_tipboxes
			}
			if len(objects) > 0 {
				collisions = append(collisions, adaptorCollision{i, objects})
			}
		}

		if len(collisions) > 0 {
			self.AddErrorf("LoadTips", "Cannot load %s, %v (Head%d not independent)",
				summariseWell2Channel(well, channels), collisions, head)
		}
	}

	//move the tips to the adaptors
	for _, ch := range channels {
		tips[ch].GetParent().(*wtype.LHTipbox).RemoveTip(wc[ch])
		adaptor.GetChannel(ch).LoadTip(tips[ch])
		if err := tips[ch].SetParent((*wtype.LHTipbox)(nil)); err != nil {
			self.AddError("LoadTips", err.Error())
		}
	}

	return ret
}

//UnloadTips - used
func (self *VirtualLiquidHandler) UnloadTips(channels []int, head, multi int,
	platetype, position, well []string) driver.CommandStatus {
	ret := driver.CommandStatus{OK: true, Errorcode: driver.OK, Msg: "UNLOADTIPS ACK"}

	//get the adaptor
	adaptor, err := self.getAdaptorState(head)
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
			self.AddWarning("UnloadTips", "'channel' argument empty and no tips are loaded, ignoring")
		} else if self.settings.IsAutoChannelWarningEnabled() {
			self.AddWarningf("UnloadTips", "'channel' argument empty, unloading all tips (%s)", summariseChannels(channels))
		}
	}

	//check that RobotState won't crash
	if !self.testTipArgs("UnloadTips", channels, head, platetype, position, well) {
		return ret
	}

	if multi != len(channels) {
		self.AddErrorf("UnloadTips", "While unloading %s from %s, multi should equal %d, not %d",
			pTips(len(channels)), summariseChannels(channels), len(channels), multi)
		return ret
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

	adaptor, err := self.getAdaptorState(head)
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
	self.AddWarningf("SetDriveSpeed", "Not yet implemented: SetDriveSpeed(%s, %f)", drive, rate)
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
	if adaptor, err := self.getAdaptorState(head); err == nil {
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

//AddPlateTo - used
func (self *VirtualLiquidHandler) AddPlateTo(position string, plate interface{}, name string) driver.CommandStatus {

	ret := driver.CommandStatus{OK: true, Errorcode: driver.OK, Msg: "ADDPLATETO ACK"}

	if obj, ok := plate.(wtype.LHObject); ok {
		obj = obj.Duplicate(true)
		if n, nok := obj.(wtype.Named); nok && n.GetName() != name {
			self.AddWarningf("AddPlateTo", "Object name(=%s) doesn't match argument name(=%s)", n.GetName(), name)
		}

		if err := self.state.GetDeck().SetChild(position, obj); err != nil {
			self.AddError("AddPlateTo", err.Error())
			return ret
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
func (self *VirtualLiquidHandler) GetOutputFile() (string, driver.CommandStatus) {
	self.AddWarning("GetOutputFile", "Not yet implemented")
	return "You forgot to say 'please'", driver.CommandStatus{OK: true, Errorcode: driver.OK, Msg: "GETOUTPUTFILE ACK"}
}
