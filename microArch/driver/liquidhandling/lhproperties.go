// /anthalib/driver/liquidhandling/types.go: Part of the Antha language
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
	"context"
	"encoding/json"
	"fmt"
	"sort"

	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/inventory"
)

// describes a liquid handler, its capabilities and current state
// probably needs splitting up to separate out the state information
// from the properties information
type LHProperties struct {
	ID             string
	Positions      map[string]*wtype.LHPosition // position descriptions by position name
	PlateLookup    map[string]interface{}       // deck object (plate, tipbox, etc) by object ID
	PosLookup      map[string]string            // object ID by position name
	PlateIDLookup  map[string]string            // position name by object ID
	Plates         map[string]*wtype.Plate      // plates by position name
	Tipboxes       map[string]*wtype.LHTipbox   // tipboxes by position name
	Tipwastes      map[string]*wtype.LHTipwaste // tipwastes by position name
	Wastes         map[string]*wtype.Plate      // waste plates by position name
	Washes         map[string]*wtype.Plate      // wash plates by position name
	Model          string
	Mnfr           string
	LHType         LiquidHandlerLevel      // describes which liquidhandling API should be used to communicate with the device
	TipType        TipType                 // defines the type of tips used by the liquidhandler
	Heads          []*wtype.LHHead         // lists every head (whether loaded or not) that is available for the machine
	Adaptors       []*wtype.LHAdaptor      // lists every adaptor (whether loaded or not) that is available for the machine
	HeadAssemblies []*wtype.LHHeadAssembly // describes how each loaded head and adaptor is loaded into the machine
	Tips           []*wtype.LHTip          // lists each type of tip available in the current configuration
	Preferences    *LayoutOpt              // describes where different categories of objects are to be placed on the liquid handler
	Driver         LiquidhandlingDriver    `gotopb:"-"`
}

func (lhp *LHProperties) MarshalJSON() ([]byte, error) {
	return json.Marshal(newSProperties(lhp))
}

func (lhp *LHProperties) UnmarshalJSON(data []byte) error {
	var slhp sProperties
	if err := json.Unmarshal(data, &slhp); err != nil {
		return err
	}

	slhp.Fill(lhp)
	return nil
}

// utility print function

func (p LHProperties) OutputLayout() {
	fmt.Println(p.GetLayout())
}

func (p LHProperties) GetLayout() string {
	s := ""
	s += fmt.Sprintln("Layout for liquid handler ", p.ID, " type ", p.Mnfr, " ", p.Model)
	n := p.OrderedPositionNames()

	for _, pos := range n {
		plateID, ok := p.PosLookup[pos]

		s += fmt.Sprint("\tPosition ", pos, " ")

		if !ok {
			s += fmt.Sprintln(" Empty")
		} else {
			lw := p.PlateLookup[plateID]

			switch obj := lw.(type) {
			case *wtype.Plate:
				s += fmt.Sprintln("Plate ", obj.PlateName, " type ", obj.Mnfr, " ", obj.Type, " Contents:")
				s += obj.GetLayout()
			case *wtype.LHTipbox:
				s += fmt.Sprintln("Tip box ", obj.Mnfr, " ", obj.Type, " ", obj.Boxname, " ", obj.N_clean_tips())
			case *wtype.LHTipwaste:
				s += fmt.Sprintln("Tip Waste ", obj.Mnfr, " ", obj.Type, " capacity ", obj.SpaceLeft())
			default:
				s += fmt.Sprintln("Labware :", lw)
			}
		}
	}

	return s
}

func (p LHProperties) CanPrompt() bool {
	// presently true for all varieties of liquid handlers
	// TODO (when no longer the case): revise!
	return true
}

func (p LHProperties) OrderedPositionNames() []string {
	// canonical ordering

	s := make([]string, 0, len(p.Positions))
	for n := range p.Positions {
		s = append(s, n)
	}

	sort.Strings(s)

	return s
}

//CountHeadsLoaded return the total number of heads loaded into the machine
func (lhp *LHProperties) CountHeadsLoaded() int {
	var ret int
	for _, assembly := range lhp.HeadAssemblies {
		ret += assembly.CountHeadsLoaded()
	}
	return ret
}

//GetLoadedHeads get a slice of all the heads loaded in the machine
func (lhp *LHProperties) GetLoadedHeads() []*wtype.LHHead {
	ret := make([]*wtype.LHHead, 0, lhp.CountHeadsLoaded())
	for _, assembly := range lhp.HeadAssemblies {
		ret = append(ret, assembly.GetLoadedHeads()...)
	}
	return ret
}

//GetLoadedHead returns a specific head
func (lhp *LHProperties) GetLoadedHead(i int) *wtype.LHHead {
	//inefficient implementation for now since we only have a small number of heads
	return lhp.GetLoadedHeads()[i]
}

//GetLoadedAdaptors get a slice of all the adaptors loaded in the machine
func (lhp *LHProperties) GetLoadedAdaptors() []*wtype.LHAdaptor {
	heads := lhp.GetLoadedHeads()
	ret := make([]*wtype.LHAdaptor, 0, len(heads))
	for _, head := range heads {
		ret = append(ret, head.Adaptor)
	}
	return ret
}

//GetLoadedAdaptor get the adaptor loaded to the i^th head
func (lhp *LHProperties) GetLoadedAdaptor(i int) *wtype.LHAdaptor {
	return lhp.GetLoadedAdaptors()[i]
}

// copy constructor
func (lhp *LHProperties) Dup() *LHProperties {
	return lhp.dup(false)
}

func (lhp *LHProperties) DupKeepIDs() *LHProperties {
	return lhp.dup(true)
}

func (lhp *LHProperties) dup(keepIDs bool) *LHProperties {
	if lhp == nil {
		return nil
	}

	pos := make(map[string]*wtype.LHPosition, len(lhp.Positions))
	for k, v := range lhp.Positions {
		// be sure to copy the data
		w := *v
		pos[k] = &w
	}
	r := NewLHProperties(lhp.Model, lhp.Mnfr, lhp.LHType, lhp.TipType, pos)

	if keepIDs {
		r.ID = lhp.ID
	}

	adaptorMap := make(map[*wtype.LHAdaptor]*wtype.LHAdaptor, len(lhp.Adaptors))
	for _, a := range lhp.Adaptors {
		var ad *wtype.LHAdaptor
		if keepIDs {
			ad = a.DupKeepIDs()
		} else {
			ad = a.Dup()
		}
		r.Adaptors = append(r.Adaptors, ad)
		adaptorMap[a] = ad
	}

	headMap := make(map[*wtype.LHHead]*wtype.LHHead, len(lhp.Heads))
	for _, h := range lhp.Heads {
		var hd *wtype.LHHead
		adaptor := adaptorMap[h.Adaptor]
		if keepIDs {
			hd = h.DupKeepIDs()
		} else {
			hd = h.Dup()
		}
		hd.Adaptor = adaptor
		r.Heads = append(r.Heads, hd)
		headMap[h] = hd
	}

	for _, assembly := range lhp.HeadAssemblies {
		//duplicate the assmebly
		newAssembly := assembly.DupWithoutHeads()
		//now add the heads - this way r.HeadAssemblies and r.Heads refer to the same underlying LHHead
		for _, oldHead := range assembly.GetLoadedHeads() {
			newAssembly.LoadHead(headMap[oldHead]) //nolint - assemblies have the same number of positions
		}
		r.HeadAssemblies = append(r.HeadAssemblies, newAssembly)
	}

	// plate lookup can contain anything

	for name, pt := range lhp.PlateLookup {
		var pt2 interface{}
		var newid string
		var addr string
		switch obj := pt.(type) {
		case *wtype.LHTipwaste:
			var tmp *wtype.LHTipwaste
			if keepIDs {
				tmp = obj.Dup()
				tmp.ID = obj.ID
			} else {
				tmp = obj.Dup()
			}
			pt2 = tmp
			newid = tmp.ID
			addr = lhp.PlateIDLookup[name]
			r.Tipwastes[addr] = tmp
		case *wtype.Plate:
			var tmp *wtype.Plate
			if keepIDs {
				tmp = obj.DupKeepIDs()
			} else {
				tmp = obj.Dup()
			}
			pt2 = tmp
			newid = tmp.ID
			addr = lhp.PlateIDLookup[name]
			_, waste := lhp.Wastes[addr]
			_, wash := lhp.Washes[addr]

			if waste {
				r.Wastes[addr] = tmp
			} else if wash {
				r.Washes[addr] = tmp
			} else {
				r.Plates[addr] = tmp
			}
		case *wtype.LHTipbox:
			var tmp *wtype.LHTipbox
			if keepIDs {
				tmp = obj.DupKeepIDs()
			} else {
				tmp = obj.Dup()
			}
			pt2 = tmp
			newid = tmp.ID
			addr = lhp.PlateIDLookup[name]
			r.Tipboxes[addr] = tmp
		}
		r.PlateLookup[newid] = pt2
		r.PlateIDLookup[newid] = addr
		r.PosLookup[addr] = newid
	}

	for _, tip := range lhp.Tips {
		newtip := tip.Dup()
		if keepIDs {
			newtip.ID = tip.ID
		}

		r.Tips = append(r.Tips, newtip)
	}

	r.Preferences = lhp.Preferences.Dup()

	// copy the driver
	r.Driver = lhp.Driver

	return r
}

// constructor for the above
func NewLHProperties(model, manufacturer string, lhtype LiquidHandlerLevel, tiptype TipType, positions map[string]*wtype.LHPosition) *LHProperties {
	// assert validity of lh and tip types

	if !lhtype.IsValid() {
		panic(fmt.Sprintf("Invalid liquid handling type requested: %s", lhtype))
	}
	if !tiptype.IsValid() {
		panic(fmt.Sprintf("Invalid tip usage type requested: %s", tiptype))
	}

	return &LHProperties{
		ID:             wtype.GetUUID(),
		Positions:      positions,
		Model:          model,
		Mnfr:           manufacturer,
		LHType:         lhtype,
		TipType:        tiptype,
		Heads:          make([]*wtype.LHHead, 0, 2),
		Adaptors:       make([]*wtype.LHAdaptor, 0, 2),
		HeadAssemblies: make([]*wtype.LHHeadAssembly, 0, 2),
		PosLookup:      make(map[string]string, len(positions)),
		PlateLookup:    make(map[string]interface{}, len(positions)),
		PlateIDLookup:  make(map[string]string, len(positions)),
		Plates:         make(map[string]*wtype.Plate, len(positions)),
		Tipboxes:       make(map[string]*wtype.LHTipbox, len(positions)),
		Tipwastes:      make(map[string]*wtype.LHTipwaste, len(positions)),
		Wastes:         make(map[string]*wtype.Plate, len(positions)),
		Washes:         make(map[string]*wtype.Plate, len(positions)),
		Tips:           make([]*wtype.LHTip, 0, 3),
	}
}

// GetLHType returns the declared type of liquid handler for driver selection purposes
// e.g. High-Level (HLLiquidHandler) or Low-Level (LLLiquidHandler)
// see lhtype.go in this directory
func (lhp *LHProperties) GetLHType() LiquidHandlerLevel {
	return lhp.LHType
}

// GetTipType returns the tip requirements of the liquid handler
// options are None, Disposable, Fixed, Mixed
// see lhtype.go in this directory
func (lhp *LHProperties) GetTipType() TipType {
	return lhp.TipType
}

func (lhp *LHProperties) AddTipBox(tipbox *wtype.LHTipbox) error {
	for _, addr := range lhp.Preferences.Tipboxes {
		if !lhp.IsEmpty(addr) {
			continue
		}

		lhp.AddTipBoxTo(addr, tipbox)
		return nil
	}

	return wtype.LHError(wtype.LH_ERR_NO_DECK_SPACE, "Trying to add tip box")
}
func (lhp *LHProperties) AddTipBoxTo(addr string, tipbox *wtype.LHTipbox) bool {
	if !lhp.IsEmpty(addr) {
		fmt.Printf("Tried to add tipbox to full position: %s\n", addr)
		return false
	}
	lhp.Tipboxes[addr] = tipbox
	lhp.PlateLookup[tipbox.ID] = tipbox
	lhp.PosLookup[addr] = tipbox.ID
	lhp.PlateIDLookup[tipbox.ID] = addr

	return true
}

func (lhp *LHProperties) RemoveTipBoxes() {
	for addr, tbx := range lhp.Tipboxes {
		lhp.PlateLookup[tbx.ID] = nil
		lhp.PosLookup[addr] = ""
		lhp.PlateIDLookup[tbx.ID] = ""
	}

	lhp.Tipboxes = make(map[string]*wtype.LHTipbox)
}

func (lhp *LHProperties) TipWastesMounted() int {
	r := 0
	// go looking for tipwastes
	for _, addr := range lhp.Preferences.Tipwastes {
		if _, ok := lhp.Tipwastes[addr]; ok {
			r += 1
		}
	}

	return r

}

func (lhp *LHProperties) TipSpacesLeft() int {
	r := 0
	// go looking for tipboxes
	for _, addr := range lhp.Preferences.Tipwastes {
		if bx, ok := lhp.Tipwastes[addr]; !ok {
			fmt.Printf("Position %s claims to have a tipbox but is empty\n", addr)
		} else {
			r += bx.SpaceLeft()
		}
	}

	return r
}

// IsEmpty returns true if the given position exists and is unoccupied
func (lhp *LHProperties) IsEmpty(address string) bool {
	return lhp.Exists(address) && lhp.PosLookup[address] == ""
}

// Exists returns true if the given address refers to a known position
func (lhp *LHProperties) Exists(address string) bool {
	_, ret := lhp.Positions[address]
	return ret
}

func (lhp *LHProperties) AddTipWaste(tipwaste *wtype.LHTipwaste) error {
	for _, addr := range lhp.Preferences.Tipwastes {
		if !lhp.IsEmpty(addr) {
			continue
		}

		err := lhp.AddTipWasteTo(addr, tipwaste)
		return err
	}

	return wtype.LHError(wtype.LH_ERR_NO_DECK_SPACE, "Trying to add tip waste")
}

func (lhp *LHProperties) AddTipWasteTo(addr string, tipwaste *wtype.LHTipwaste) error {
	if !lhp.IsEmpty(addr) {
		return wtype.LHError(wtype.LH_ERR_NO_DECK_SPACE, fmt.Sprintf("Trying to add tip waste to full position %s", addr))
	}

	lhp.Tipwastes[addr] = tipwaste
	lhp.PlateLookup[tipwaste.ID] = tipwaste
	lhp.PosLookup[addr] = tipwaste.ID
	lhp.PlateIDLookup[tipwaste.ID] = addr
	return nil
}

func (lhp *LHProperties) AddInputPlate(plate *wtype.Plate) error {
	for _, addr := range lhp.Preferences.Inputs {
		if !lhp.IsEmpty(addr) {
			continue
		}

		err := lhp.AddPlateTo(addr, plate)
		return err
	}

	return wtype.LHError(wtype.LH_ERR_NO_DECK_SPACE, fmt.Sprintf("Trying to add input plate %s, type %s", plate.PlateName, plate.Type))
}
func (lhp *LHProperties) AddOutputPlate(plate *wtype.Plate) error {
	for _, addr := range lhp.Preferences.Outputs {
		if !lhp.IsEmpty(addr) {
			continue
		}

		err := lhp.AddPlateTo(addr, plate)
		return err
	}

	return wtype.LHError(wtype.LH_ERR_NO_DECK_SPACE, fmt.Sprintf("Trying to add output plate %s, type %s", plate.PlateName, plate.Type))
}

func (lhp *LHProperties) AddPlateTo(addr string, plate *wtype.Plate) error {
	if !lhp.IsEmpty(addr) {
		return wtype.LHError(wtype.LH_ERR_NO_DECK_SPACE, fmt.Sprintf("Trying to add plate to full position %s", addr))
	}
	lhp.Plates[addr] = plate
	lhp.PlateLookup[plate.ID] = plate
	lhp.PosLookup[addr] = plate.ID
	lhp.PlateIDLookup[plate.ID] = addr
	return nil
}

// reverse the above

func (lhp *LHProperties) RemovePlateWithID(id string) {
	addr := lhp.PlateIDLookup[id]
	delete(lhp.PosLookup, addr)
	delete(lhp.PlateIDLookup, id)
	delete(lhp.PlateLookup, id)
	delete(lhp.Plates, addr)
}

// UpdateID update the id of the object with currID, returning the new ID
// and an error if no such object was found
func (lhp *LHProperties) UpdateID(currID string) (string, error) {
	if thing, ok := lhp.PlateLookup[currID]; !ok {
		return "", wtype.LHErrorf(wtype.LH_ERR_DIRE, "while updating id: no object with id %s found", currID)
	} else if addr, ok := lhp.PlateIDLookup[currID]; !ok {
		return "", wtype.LHErrorf(wtype.LH_ERR_DIRE, "while updating id: no position found for object with id %s", currID)
	} else if obj, ok := thing.(wtype.LHObject); !ok {
		return "", wtype.LHErrorf(wtype.LH_ERR_DIRE, "while updating id: object with id %s not of type LHObject", currID)
	} else {
		// make a copy of the object with a new ID
		newObj := obj.Duplicate(false)

		// update references to the old plate and ID
		delete(lhp.PlateLookup, currID)
		delete(lhp.PlateIDLookup, currID)
		lhp.PlateLookup[wtype.IDOf(newObj)] = newObj
		lhp.PlateIDLookup[wtype.IDOf(newObj)] = addr
		lhp.PosLookup[addr] = wtype.IDOf(newObj)

		switch o := newObj.(type) {
		case *wtype.Plate:
			lhp.Plates[addr] = o
			if _, ok := lhp.Wastes[addr]; ok {
				lhp.Wastes[addr] = o
			}
			if _, ok := lhp.Washes[addr]; ok {
				lhp.Washes[addr] = o
			}
		case *wtype.LHTipbox:
			lhp.Tipboxes[addr] = o
		case *wtype.LHTipwaste:
			lhp.Tipwastes[addr] = o
		}

		return wtype.IDOf(newObj), nil
	}
}

func (lhp *LHProperties) RemovePlateAtPosition(addr string) {
	id := lhp.PosLookup[addr]
	delete(lhp.PosLookup, addr)
	delete(lhp.PlateIDLookup, id)
	delete(lhp.PlateLookup, id)
	delete(lhp.Plates, addr)
}

func (lhp *LHProperties) AddWaste(waste *wtype.Plate) bool {
	for _, addr := range lhp.Preferences.Wastes {
		if !lhp.IsEmpty(addr) {
			continue
		}

		lhp.AddWasteTo(addr, waste)
		return true
	}

	fmt.Println("NO WASTE SPACES LEFT")
	return false
}

func (lhp *LHProperties) AddWasteTo(addr string, waste *wtype.Plate) bool {
	if !lhp.IsEmpty(addr) {
		fmt.Println("CAN'T ADD WASTE TO FULL POSITION")
		return false
	}
	lhp.Wastes[addr] = waste
	lhp.PlateLookup[waste.ID] = waste
	lhp.PosLookup[addr] = waste.ID
	lhp.PlateIDLookup[waste.ID] = addr
	return true
}

func (lhp *LHProperties) AddWash(wash *wtype.Plate) bool {
	for _, addr := range lhp.Preferences.Washes {
		if !lhp.IsEmpty(addr) {
			continue
		}

		lhp.AddWashTo(addr, wash)
		return true
	}

	fmt.Println("NO WASH SPACES LEFT")
	return false
}

func (lhp *LHProperties) AddWashTo(addr string, wash *wtype.Plate) bool {
	if !lhp.IsEmpty(addr) {

		fmt.Println("CAN'T ADD WASH TO FULL POSITION")
		return false
	}
	lhp.Washes[addr] = wash
	lhp.PlateLookup[wash.ID] = wash
	lhp.PosLookup[addr] = wash.ID
	lhp.PlateIDLookup[wash.ID] = addr
	return true
}

func (lhp *LHProperties) InputSearchPreferences() []string {
	// Definition 1: merge input plate preferences and output plate preferences
	return lhp.mergeInputOutputPreferences()
}

func (lhp *LHProperties) mergeInputOutputPreferences() []string {
	seen := make(map[string]bool, len(lhp.Positions))
	out := make([]string, 0, len(lhp.Positions))

	mergeToSet := func(in, out []string, seen map[string]bool) ([]string, map[string]bool) {
		// adds anything from in to out that isn't in seen, respecting input order
		for _, mem := range in {
			if seen[mem] {
				continue
			}
			seen[mem] = true
			out = append(out, mem)

		}

		return out, seen
	}

	out, seen = mergeToSet(lhp.Preferences.Inputs, out, seen)
	out, _ = mergeToSet(lhp.Preferences.Outputs, out, seen)

	return out
}

// logic of getting components:
// we look for things with the same ID
// the ID may or may not refer to an instance which is previously made
// but by this point we must have concrete locations for everything

func (lhp *LHProperties) GetComponentsSingle(cmps []*wtype.Liquid, carryvol wunit.Volume, legacyVolume bool) ([][]string, [][]string, [][]wunit.Volume, error) {
	plateIDs := make([][]string, len(cmps))
	wellCoords := make([][]string, len(cmps))
	vols := make([][]wunit.Volume, len(cmps))

	// locally keep volumes straight

	localplates := make(map[string]*wtype.Plate, len(lhp.Plates))

	for k, v := range lhp.Plates {
		localplates[k] = v.DupKeepIDs()
	}

	// cmps are requests for components
	for i, cmp := range cmps {
		plateIDs[i] = make([]string, 0, 1)
		wellCoords[i] = make([]string, 0, 1)
		vols[i] = make([]wunit.Volume, 0, 1)
		foundIt := false

		cmpdup := cmp.Dup()

		// searches all plates: input and output
		for _, ipref := range lhp.InputSearchPreferences() {
			// check if the plate at position ipref has the
			// component we seek

			p, ok := localplates[ipref]
			if ok && !p.IsEmpty() {
				// whaddya got?
				// nb this won't work if we need to split a volume across several plates
				wcarr, varr, ok := p.BetterGetComponent(cmpdup, lhp.MinPossibleVolume(), legacyVolume)

				if ok {
					foundIt = true
					for ix := range wcarr {
						wc := wcarr[ix].FormatA1()
						vl := varr[ix].Dup()
						plateIDs[i] = append(plateIDs[i], p.ID)
						wellCoords[i] = append(wellCoords[i], wc)
						vols[i] = append(vols[i], vl)
						vl = vl.Dup()
						vl.Add(carryvol)
						//lhp.RemoveComponent(p.ID, wc, vl)
						p.RemoveComponent(wc, vl)
					}
					break
				}
			}
		}

		if !foundIt {
			err := wtype.LHError(wtype.LH_ERR_DIRE, fmt.Sprint("749: NO SOURCE FOR ", cmp.CName, " at volume ", cmp.Volume().ToString()))
			return plateIDs, wellCoords, vols, err
		}
	}

	return plateIDs, wellCoords, vols, nil
}

func (lhp *LHProperties) GetCleanTips(ctx context.Context, tiptype []string, channel []*wtype.LHChannelParameter, usetiptracking bool) (wells, positions, boxtypes [][]string, err error) {

	// these are merged into subsets with tip and channel types in common here
	// each subset has a mask which is the same size as the number of channels available
	subsets, err2 := makeChannelSubsets(tiptype, channel)

	if err2 != nil {
		return [][]string{}, [][]string{}, [][]string{}, err2
	}

	for _, set := range subsets {
		sw, sp, sb, err := lhp.getCleanTipSubset(ctx, set, usetiptracking)

		if err != nil {
			return [][]string{}, [][]string{}, [][]string{}, err
		}

		wells = append(wells, sw)
		positions = append(positions, sp)
		boxtypes = append(boxtypes, sb)
	}

	return wells, positions, boxtypes, nil
}

func countMultiB(ar []bool) int {
	r := 0
	for _, v := range ar {
		if v {
			r += 1
		}
	}

	return r
}

func copyToRightLength(sa []string, m int) []string {
	r := make([]string, m)

	for i := 0; i < len(sa); i++ {
		r[i] = sa[i]
	}

	return r
}

// this function only returns true if we can get all tips at once
// TODO -- support not getting in a single operation
func (lhp *LHProperties) getCleanTipSubset(ctx context.Context, tipParams TipSubset, usetiptracking bool) (wells, positions, boxtypes []string, err error) {
	positions = make([]string, len(tipParams.Mask))
	boxtypes = make([]string, len(tipParams.Mask))

	foundit := false
	multi := countMultiB(tipParams.Mask)

	for _, addr := range lhp.Preferences.Tipboxes {
		bx, ok := lhp.Tipboxes[addr]
		if !ok || bx.Tiptype.Type != tipParams.TipType {
			continue
		}
		wells, err = bx.GetTipsMasked(tipParams.Mask, tipParams.Channel.Orientation, true)

		/*
			if err != nil && !bx.IsEmpty() {
				return wells, positions, boxtypes, err
			}
		*/

		// update wells

		if len(wells) != len(positions) {
			wells = copyToRightLength(wells, len(positions))
		}

		// TODO -- support partial collections
		if wells != nil && countMulti(wells) == multi {
			foundit = true
			for i := 0; i < len(wells); i++ {
				if tipParams.Mask[i] {
					positions[i] = addr
					boxtypes[i] = bx.Boxname
				}
			}
			break
		} else if usetiptracking && lhp.HasTipTracking() {
			bx.Refresh()
			return lhp.getCleanTipSubset(ctx, tipParams, usetiptracking)
		}
	}

	// if you don't find any suitable tips, why just make a
	// new box full of them!
	// nothing can possibly go wrong
	// surely

	if !foundit {
		// try adding a new tip box
		bx, err := inventory.NewTipbox(ctx, tipParams.TipType)

		if err != nil {
			return nil, nil, nil, wtype.LHError(wtype.LH_ERR_NO_TIPS, fmt.Sprintf("No tipbox of type %s found: %s", tipParams.TipType, err))
		}

		r := lhp.AddTipBox(bx)

		if r != nil {
			err = r
			return nil, nil, nil, err
		}

		return lhp.getCleanTipSubset(ctx, tipParams, usetiptracking)
		//		return nil, nil, nil
	}

	return
}

// DropDirtyTips figure out where tips attached to the non-nil channels in channels can be disposed
// The return arrays are the same length as channels
func (lhp *LHProperties) DropDirtyTips(channels []*wtype.LHChannelParameter) (wells, positions, boxtypes []string) {
	multi := len(channels)

	wells = make([]string, multi)
	positions = make([]string, multi)
	boxtypes = make([]string, multi)

	foundit := false

	for addr, bx := range lhp.Tipwastes {
		wellCoords, ok := bx.Dispose(channels)
		if ok {
			foundit = true
			for i := 0; i < multi; i++ {
				if channels[i] != nil {
					wells[i] = wellCoords[i].FormatA1()
					positions[i] = addr
					boxtypes[i] = bx.Type
				}
			}

			break
		}
	}

	if !foundit {
		return nil, nil, nil
	}

	return
}

func (lhp *LHProperties) GetTimer() LHTimer {
	return GetTimerFor(lhp.Mnfr, lhp.Model)
}

func (lhp *LHProperties) GetChannelScoreFunc() ChannelScoreFunc {
	// this is to permit us to make this flexible

	sc := DefaultChannelScoreFunc{}

	return sc
}

// convenience method

func (lhp *LHProperties) RemoveComponent(plateID string, well string, volume wunit.Volume) bool {
	p := lhp.Plates[lhp.PlateIDLookup[plateID]]

	if p == nil {
		fmt.Println("RemoveComponent ", plateID, " ", well, " ", volume.ToString(), " can't find plate")
		return false
	}

	r := p.RemoveComponent(well, volume)

	if r == nil {
		fmt.Println("CAN'T REMOVE COMPONENT ", plateID, " ", well, " ", volume.ToString())
		return false
	}

	/*
		w := p.Wellcoords[well]

		if w == nil {
			return false
		}

		c:=w.Remove(volume)

		if c==nil{
			return false
		}
	*/

	return true
}

// ApplyUserPreferences merge in the layout preferences given by the user.
//
// User preferences for each category should either be list of addresses to place
// items of that category in order, or empty. If they are empty, then the full list
// of possible locations as reported by the driver is used.
//
// nb.
// Because of the difficulties surrounding cross-platform addresses, addresses
// which are don't exist in this liquid handler are silently ignored
// such that passing a Gilson address e.g. "position_1" to a Hamilton driver has
// no effect.
func (lhp *LHProperties) ApplyUserPreferences(p *LayoutOpt) error {
	// ignore addresses that don't exist in this liquidhandler
	filter := func(addresses Addresses) Addresses {
		ret := make(Addresses, 0, len(addresses))
		for _, a := range addresses {
			if lhp.Exists(a) {
				ret = append(ret, a)
			}
		}
		return ret
	}
	q := &LayoutOpt{
		Tipboxes:  filter(p.Tipboxes),
		Inputs:    filter(p.Inputs),
		Outputs:   filter(p.Outputs),
		Tipwastes: filter(p.Tipwastes),
		Wastes:    filter(p.Wastes),
		Washes:    filter(p.Washes),
	}

	if merged, err := lhp.Preferences.ApplyUserPreferences(q); err != nil {
		return err
	} else {
		lhp.Preferences = merged
		return nil
	}
}

type UserPlate struct {
	Plate    *wtype.Plate
	Position string
}
type UserPlates []UserPlate

func (p *LHProperties) SaveUserPlates() UserPlates {
	up := make(UserPlates, 0, len(p.Positions))

	for addr, plate := range p.Plates {
		if plate.IsUserAllocated() {
			up = append(up, UserPlate{Plate: plate.DupKeepIDs(), Position: addr})
		}
	}

	return up
}

func (p *LHProperties) RestoreUserPlates(up UserPlates) {
	for _, plate := range up {
		oldPlate := p.Plates[plate.Position]
		p.RemovePlateAtPosition(plate.Position)
		// merge these
		plate.Plate.MergeWith(oldPlate)

		err := p.AddPlateTo(plate.Position, plate.Plate)
		if err != nil {
			panic(err)
		}
	}
}

func (p *LHProperties) MinPossibleVolume() wunit.Volume {
	headsLoaded := p.GetLoadedHeads()
	if len(headsLoaded) == 0 {
		return wunit.ZeroVolume()
	}
	minvol := headsLoaded[0].GetParams().Minvol
	for _, head := range headsLoaded {
		for _, tip := range p.Tips {
			lhcp := head.Params.MergeWithTip(tip)
			v := lhcp.Minvol
			if v.LessThan(minvol) {
				minvol = v
			}
		}

	}

	return minvol
}

func (p *LHProperties) CanPossiblyDo(v wunit.Volume) bool {
	return !p.MinPossibleVolume().LessThan(v)
}

func (p *LHProperties) IsAddressable(pos string, crd wtype.WellCoords, channel, reference int, offsetX, offsetY, offsetZ float64) bool {
	// can we reach well 'crd' at position 'pos' using channel 'channel' at reference 'reference' with
	// the given offsets?

	// yes (this will improve, honest!)
	return true
}

func dupStrArr(sa []string) []string {
	ret := make([]string, len(sa))
	copy(ret, sa)
	return ret
}

func inStrArr(s string, sa []string) bool {
	for _, v := range sa {
		if s == v {
			return true
		}
	}

	return false
}

func (p *LHProperties) OrderedMergedPlatePrefs() []string {
	r := dupStrArr(p.Preferences.Inputs)

	for _, pr := range p.Preferences.Outputs {
		if !inStrArr(pr, r) {
			r = append(r, pr)
		}
	}

	return r
}

func (p LHProperties) HasTipTracking() bool {
	// TODO --> improve this

	if p.Mnfr == "Gilson" && p.Model == "Pipetmax" {
		return true
	}

	return false
}

func (p *LHProperties) UpdateComponentIDs(updates map[string]*wtype.Liquid) {
	for s, c := range updates {
		p.UpdateComponentID(s, c)
	}
}

func (p *LHProperties) UpdateComponentID(from string, to *wtype.Liquid) bool {
	for _, p := range p.Plates {
		if p.FindAndUpdateID(from, to) {
			return true
		}
	}
	return false
}

func (p *LHProperties) DeckSummary() string {
	s := ""
	for name, thing := range p.PlateLookup {
		if thing == nil {
			s += fmt.Sprintf("%s: %s ", name, "Empty")
		} else {
			n, ok := thing.(wtype.Named)

			if ok {
				s += fmt.Sprintf("%s: %s ", name, n.GetName())
			} else {
				s += fmt.Sprintf("%s: %s ", name, "Something")
			}
		}
	}
	return s
}

func (p *LHProperties) CarryVolume() wunit.Volume {
	return wtype.GLOBALCARRYVOLUME
}

func (p *LHProperties) GetSetupInstructions() []TerminalRobotInstruction {
	instructions := make([]TerminalRobotInstruction, 0, 1+len(p.PosLookup))

	//first instruction is always to remove all plates
	instructions = append(instructions, NewRemoveAllPlatesInstruction())

	for position, plateid := range p.PosLookup {
		if plateid == "" {
			continue
		}
		plate := p.PlateLookup[plateid]
		name := plate.(wtype.Named).GetName()

		ins := NewAddPlateToInstruction(position, name, plate)

		instructions = append(instructions, ins)
	}
	return instructions
}
