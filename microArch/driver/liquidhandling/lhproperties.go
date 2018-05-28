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
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/antha-lang/antha/antha/anthalib/material"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/antha/anthalib/wutil"
	"github.com/antha-lang/antha/inventory"
	"github.com/antha-lang/antha/microArch/logger"
	"github.com/antha-lang/antha/microArch/sampletracker"
)

// describes a liquid handler, its capabilities and current state
// probably needs splitting up to separate out the state information
// from the properties information
type LHProperties struct {
	ID                   string
	Nposns               int
	Positions            map[string]*wtype.LHPosition
	PlateLookup          map[string]interface{}
	PosLookup            map[string]string
	PlateIDLookup        map[string]string
	Plates               map[string]*wtype.LHPlate
	Tipboxes             map[string]*wtype.LHTipbox
	Tipwastes            map[string]*wtype.LHTipwaste
	Wastes               map[string]*wtype.LHPlate
	Washes               map[string]*wtype.LHPlate
	Devices              map[string]string
	Model                string
	Mnfr                 string
	LHType               string
	TipType              string
	Heads                []*wtype.LHHead
	HeadsLoaded          []*wtype.LHHead
	Adaptors             []*wtype.LHAdaptor
	Tips                 []*wtype.LHTip
	Tip_preferences      []string
	Input_preferences    []string
	Output_preferences   []string
	Tipwaste_preferences []string
	Waste_preferences    []string
	Wash_preferences     []string
	Driver               LiquidhandlingDriver `gotopb:"-"`
	CurrConf             *wtype.LHChannelParameter
	Cnfvol               []*wtype.LHChannelParameter
	Layout               map[string]wtype.Coordinates
	MaterialType         material.MaterialType
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

			switch lw.(type) {
			case *wtype.LHPlate:
				plt := lw.(*wtype.LHPlate)
				s += fmt.Sprintln("Plate ", plt.PlateName, " type ", plt.Mnfr, " ", plt.Type, " Contents:")
				s += plt.GetLayout()
			case *wtype.LHTipbox:
				tb := lw.(*wtype.LHTipbox)
				s += fmt.Sprintln("Tip box ", tb.Mnfr, " ", tb.Type, " ", tb.Boxname, " ", tb.N_clean_tips())
			case *wtype.LHTipwaste:
				tw := lw.(*wtype.LHTipwaste)
				s += fmt.Sprintln("Tip Waste ", tw.Mnfr, " ", tw.Type, " capacity ", tw.SpaceLeft())
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

// validator for LHProperties structure

func ValidateLHProperties(props *LHProperties) (bool, string) {
	bo := true
	so := "OK"

	be := false
	se := "LHProperties Error: No position"

	if props.Positions == nil || len(props.Positions) == 0 {
		return be, se + "s"
	}

	for k, p := range props.Positions {
		if p == nil || p.ID == "" {
			return be, se + " " + k + " not set"
		}
	}

	se = "LHProperties error: No position lookup"

	if props.PosLookup == nil || len(props.PosLookup) == 0 {
		return be, se
	}

	se = "LHProperties Error: No tip preference information"

	if props.Tip_preferences == nil || len(props.Tip_preferences) == 0 {
		return be, se
	}

	se = "LHProperties Error: No input preference information"

	if props.Input_preferences == nil || len(props.Input_preferences) == 0 {
		return be, se
	}

	se = "LHProperties Error: No output preference information"

	if props.Output_preferences == nil || len(props.Output_preferences) == 0 {
		return be, se
	}

	se = "LHProperties Error: No waste preference information"

	if props.Waste_preferences == nil || len(props.Waste_preferences) == 0 {
		return be, se
	}

	se = "LHProperties Error: No tipwaste preference information"

	if props.Tipwaste_preferences == nil || len(props.Tipwaste_preferences) == 0 {
		return be, se
	}
	se = "LHProperties Error: No wash preference information"

	if props.Wash_preferences == nil || len(props.Wash_preferences) == 0 {
		return be, se
	}
	se = "LHProperties Error: No Plate ID lookup information"

	if props.PlateIDLookup == nil {
		return be, se
	}

	se = "LHProperties Error: No tip defined"

	if props.Tips == nil {
		return be, se
	}

	se = "LHProperties Error: No headsloaded array"

	if props.HeadsLoaded == nil {
		return be, se
	}

	return bo, so
}

// copy constructor
func (lhp *LHProperties) Dup() *LHProperties {
	return lhp.dup(false)
}

func (lhp *LHProperties) DupKeepIDs() *LHProperties {
	return lhp.dup(true)
}

func (lhp *LHProperties) dup(keepIDs bool) *LHProperties {
	lo := make(map[string]wtype.Coordinates, len(lhp.Layout))
	for k, v := range lhp.Layout {
		lo[k] = v
	}
	r := NewLHProperties(lhp.Nposns, lhp.Model, lhp.Mnfr, lhp.LHType, lhp.TipType, lo)

	for _, a := range lhp.Adaptors {
		ad := a.Dup()
		if keepIDs {
			ad.ID = a.ID
		}
		r.Adaptors = append(r.Adaptors, ad)
	}

	for _, h := range lhp.Heads {
		hd := h.Dup()
		if keepIDs {
			hd.ID = h.ID
		}
		r.Heads = append(r.Heads, hd)
	}

	for _, hl := range lhp.HeadsLoaded {
		hld := hl.Dup()
		if keepIDs {
			hld.ID = hl.ID
		}
		r.HeadsLoaded = append(r.HeadsLoaded, hld)
	}

	// plate lookup can contain anything

	for name, pt := range lhp.PlateLookup {
		var pt2 interface{}
		var newid string
		var pos string
		switch pt.(type) {
		case *wtype.LHTipwaste:
			var tmp *wtype.LHTipwaste
			if keepIDs {
				tmp = pt.(*wtype.LHTipwaste).Dup()
				tmp.ID = pt.(*wtype.LHTipwaste).ID
			} else {
				tmp = pt.(*wtype.LHTipwaste).Dup()
			}
			pt2 = tmp
			newid = tmp.ID
			pos = lhp.PlateIDLookup[name]
			r.Tipwastes[pos] = tmp
		case *wtype.LHPlate:
			var tmp *wtype.LHPlate
			if keepIDs {
				tmp = pt.(*wtype.LHPlate).DupKeepIDs()
			} else {
				tmp = pt.(*wtype.LHPlate).Dup()
			}
			pt2 = tmp
			newid = tmp.ID
			pos = lhp.PlateIDLookup[name]
			_, waste := lhp.Wastes[pos]
			_, wash := lhp.Washes[pos]

			if waste {
				r.Wastes[pos] = tmp
			} else if wash {
				r.Washes[pos] = tmp
			} else {
				r.Plates[pos] = tmp
			}
		case *wtype.LHTipbox:
			var tmp *wtype.LHTipbox
			if keepIDs {
				tmp = pt.(*wtype.LHTipbox).DupKeepIDs()
			} else {
				tmp = pt.(*wtype.LHTipbox).Dup()
			}
			pt2 = tmp
			newid = tmp.ID
			pos = lhp.PlateIDLookup[name]
			r.Tipboxes[pos] = tmp
		}
		r.PlateLookup[newid] = pt2
		r.PlateIDLookup[newid] = pos
		r.PosLookup[pos] = newid
	}

	for name, dev := range lhp.Devices {
		r.Devices[name] = dev
	}

	for name, head := range lhp.Heads {
		r.Heads[name] = head.Dup()
		if keepIDs {
			r.Heads[name].ID = head.ID
		}
	}

	for i, hl := range lhp.HeadsLoaded {
		r.HeadsLoaded[i] = hl.Dup()
		if keepIDs {
			r.HeadsLoaded[i].ID = hl.ID
		}
	}

	for i, ad := range lhp.Adaptors {
		r.Adaptors[i] = ad.Dup()

		if keepIDs {
			r.Adaptors[i].ID = ad.ID
		}
	}

	for _, tip := range lhp.Tips {
		newtip := tip.Dup()
		if keepIDs {
			newtip.ID = tip.ID
		}

		r.Tips = append(r.Tips, newtip)
	}

	r.Tip_preferences = append(r.Tip_preferences, lhp.Tip_preferences...)
	r.Input_preferences = append(r.Input_preferences, lhp.Input_preferences...)
	r.Output_preferences = append(r.Output_preferences, lhp.Output_preferences...)
	r.Waste_preferences = append(r.Waste_preferences, lhp.Waste_preferences...)
	r.Tipwaste_preferences = append(r.Tipwaste_preferences, lhp.Tipwaste_preferences...)
	r.Wash_preferences = append(r.Wash_preferences, lhp.Wash_preferences...)

	if lhp.CurrConf != nil {
		r.CurrConf = lhp.CurrConf.Dup()
	}

	copy(r.Cnfvol, lhp.Cnfvol)

	for i, v := range lhp.Layout {
		r.Layout[i] = v
	}

	r.MaterialType = lhp.MaterialType

	// copy the driver

	r.Driver = lhp.Driver

	return r
}

// constructor for the above
func NewLHProperties(num_positions int, model, manufacturer, lhtype, tiptype string, layout map[string]wtype.Coordinates) *LHProperties {
	// assert validity of lh and tip types

	if !IsValidLiquidHandlerType(lhtype) {
		panic(fmt.Sprintf("Invalid liquid handling type requested: %s", lhtype))
	}
	if !IsValidTipType(tiptype) {
		panic(fmt.Sprintf("Invalid tip usage type requested: %s", tiptype))
	}

	var lhp LHProperties

	lhp.ID = wtype.GetUUID()

	lhp.Nposns = num_positions

	lhp.Model = model
	lhp.Mnfr = manufacturer
	lhp.LHType = lhtype
	lhp.TipType = tiptype

	lhp.Adaptors = make([]*wtype.LHAdaptor, 0, 2)
	lhp.Heads = make([]*wtype.LHHead, 0, 2)
	lhp.HeadsLoaded = make([]*wtype.LHHead, 0, 2)

	positions := make(map[string]*wtype.LHPosition, num_positions)

	for i := 0; i < num_positions; i++ {
		// not overriding these defaults seems like a
		// bad idea --- TODO: Fix, e.g., MAXH here
		posname := fmt.Sprintf("position_%d", i+1)
		positions[posname] = wtype.NewLHPosition(i+1, "position_"+strconv.Itoa(i+1), 80.0)
	}

	lhp.Positions = positions
	lhp.PosLookup = make(map[string]string, lhp.Nposns)
	lhp.PlateLookup = make(map[string]interface{}, lhp.Nposns)
	lhp.PlateIDLookup = make(map[string]string, lhp.Nposns)
	lhp.Plates = make(map[string]*wtype.LHPlate, lhp.Nposns)
	lhp.Tipboxes = make(map[string]*wtype.LHTipbox, lhp.Nposns)
	lhp.Tipwastes = make(map[string]*wtype.LHTipwaste, lhp.Nposns)
	lhp.Wastes = make(map[string]*wtype.LHPlate, lhp.Nposns)
	lhp.Washes = make(map[string]*wtype.LHPlate, lhp.Nposns)
	lhp.Devices = make(map[string]string, lhp.Nposns)
	lhp.Heads = make([]*wtype.LHHead, 0, 2)
	lhp.Tips = make([]*wtype.LHTip, 0, 3)

	lhp.Layout = layout

	// lhp.Curcnf, lhp.Cmnvol etc. intentionally left blank

	lhp.MaterialType = material.DEVICE

	return &lhp
}

// GetLHType returns the declared type of liquid handler for driver selection purposes
// e.g. High-Level (HLLiquidHandler) or Low-Level (LLLiquidHandler)
// see lhtype.go in this directory
func (lhp *LHProperties) GetLHType() string {
	return lhp.LHType
}

// GetTipType returns the tip requirements of the liquid handler
// options are None, Disposable, Fixed, Mixed
// see lhtype.go in this directory
func (lhp *LHProperties) GetTipType() string {
	return lhp.TipType
}

func (lhp *LHProperties) TipsLeftOfType(tiptype string) int {
	n := 0

	for _, pref := range lhp.Tip_preferences {
		tb := lhp.Tipboxes[pref]
		if tb != nil {
			n += tb.N_clean_tips()
		}
	}

	return n
}

func (lhp *LHProperties) AddTipBox(tipbox *wtype.LHTipbox) error {
	for _, pref := range lhp.Tip_preferences {
		if lhp.PosLookup[pref] != "" {
			continue
		}

		lhp.AddTipBoxTo(pref, tipbox)
		return nil
	}

	return wtype.LHError(wtype.LH_ERR_NO_DECK_SPACE, "Trying to add tip box")
}
func (lhp *LHProperties) AddTipBoxTo(pos string, tipbox *wtype.LHTipbox) bool {
	if lhp.PosLookup[pos] != "" {
		logger.Debug(fmt.Sprintf("Tried to add tipbox to full position: %s", pos))
		return false
	}
	lhp.Tipboxes[pos] = tipbox
	lhp.PlateLookup[tipbox.ID] = tipbox
	lhp.PosLookup[pos] = tipbox.ID
	lhp.PlateIDLookup[tipbox.ID] = pos

	return true
}

func (lhp *LHProperties) RemoveTipBoxes() {
	for pos, tbx := range lhp.Tipboxes {
		lhp.PlateLookup[tbx.ID] = nil
		lhp.PosLookup[pos] = ""
		lhp.PlateIDLookup[tbx.ID] = ""
	}

	lhp.Tipboxes = make(map[string]*wtype.LHTipbox)
}

func (lhp *LHProperties) TipWastesMounted() int {
	r := 0
	for _, pref := range lhp.Tipwaste_preferences {
		if lhp.PosLookup[pref] != "" {
			_, ok := lhp.Tipwastes[lhp.PosLookup[pref]]

			if !ok {
				logger.Debug(fmt.Sprintf("Position %s claims to have a tipbox but is empty", pref))
				continue
			}

			r += 1
		}
	}

	return r

}

func (lhp *LHProperties) TipSpacesLeft() int {
	r := 0
	for _, pref := range lhp.Tipwaste_preferences {
		if lhp.PosLookup[pref] != "" {
			bx, ok := lhp.Tipwastes[lhp.PosLookup[pref]]

			if !ok {
				logger.Debug(fmt.Sprintf("Position %s claims to have a tipbox but is empty", pref))
				continue
			}

			r += bx.SpaceLeft()
		}
	}

	return r
}

func (lhp *LHProperties) AddTipWaste(tipwaste *wtype.LHTipwaste) error {
	for _, pref := range lhp.Tipwaste_preferences {
		if lhp.PosLookup[pref] != "" {
			continue
		}

		err := lhp.AddTipWasteTo(pref, tipwaste)
		return err
	}

	return wtype.LHError(wtype.LH_ERR_NO_DECK_SPACE, "Trying to add tip waste")
}

func (lhp *LHProperties) AddTipWasteTo(pos string, tipwaste *wtype.LHTipwaste) error {
	if lhp.PosLookup[pos] != "" {
		return wtype.LHError(wtype.LH_ERR_NO_DECK_SPACE, fmt.Sprintf("Trying to add tip waste to full position %s", pos))
	}

	lhp.Tipwastes[pos] = tipwaste
	lhp.PlateLookup[tipwaste.ID] = tipwaste
	lhp.PosLookup[pos] = tipwaste.ID
	lhp.PlateIDLookup[tipwaste.ID] = pos
	return nil
}

func (lhp *LHProperties) AddInputPlate(plate *wtype.LHPlate) error {
	for _, pref := range lhp.Input_preferences {
		if lhp.PosLookup[pref] != "" {
			continue
		}

		err := lhp.AddPlateTo(pref, plate)
		return err
	}

	return wtype.LHError(wtype.LH_ERR_NO_DECK_SPACE, fmt.Sprintf("Trying to add input plate %s, type %s", plate.PlateName, plate.Type))
}
func (lhp *LHProperties) AddOutputPlate(plate *wtype.LHPlate) error {
	for _, pref := range lhp.Output_preferences {
		if lhp.PosLookup[pref] != "" {
			continue
		}

		err := lhp.AddPlateTo(pref, plate)
		return err
	}

	return wtype.LHError(wtype.LH_ERR_NO_DECK_SPACE, fmt.Sprintf("Trying to add output plate %s, type %s", plate.PlateName, plate.Type))
}

func (lhp *LHProperties) AddPlateTo(pos string, plate *wtype.LHPlate) error {
	if lhp.PosLookup[pos] != "" {
		return wtype.LHError(wtype.LH_ERR_NO_DECK_SPACE, fmt.Sprintf("Trying to add plate to full position %s", pos))
	}
	lhp.Plates[pos] = plate
	lhp.PlateLookup[plate.ID] = plate
	lhp.PosLookup[pos] = plate.ID
	lhp.PlateIDLookup[plate.ID] = pos
	return nil
}

// reverse the above

func (lhp *LHProperties) RemovePlateWithID(id string) {
	pos := lhp.PlateIDLookup[id]
	delete(lhp.PosLookup, pos)
	delete(lhp.PlateIDLookup, id)
	delete(lhp.PlateLookup, id)
	delete(lhp.Plates, pos)
}

func (lhp *LHProperties) RemovePlateAtPosition(pos string) {
	id := lhp.PosLookup[pos]
	delete(lhp.PosLookup, pos)
	delete(lhp.PlateIDLookup, id)
	delete(lhp.PlateLookup, id)
	delete(lhp.Plates, pos)
}

func (lhp *LHProperties) AddWasteTo(pos string, waste *wtype.LHPlate) bool {
	if lhp.PosLookup[pos] != "" {
		logger.Debug("CAN'T ADD WASTE TO FULL POSITION")
		return false
	}
	lhp.Wastes[pos] = waste
	lhp.PlateLookup[waste.ID] = waste
	lhp.PosLookup[pos] = waste.ID
	lhp.PlateIDLookup[waste.ID] = pos
	return true
}

func (lhp *LHProperties) AddWash(wash *wtype.LHPlate) bool {
	for _, pref := range lhp.Wash_preferences {
		if lhp.PosLookup[pref] != "" {
			continue
		}

		lhp.AddWashTo(pref, wash)
		return true
	}

	logger.Debug("NO WASH SPACES LEFT")
	return false
}

func (lhp *LHProperties) AddWashTo(pos string, wash *wtype.LHPlate) bool {
	if lhp.PosLookup[pos] != "" {

		logger.Debug("CAN'T ADD WASH TO FULL POSITION")
		return false
	}
	lhp.Washes[pos] = wash
	lhp.PlateLookup[wash.ID] = wash
	lhp.PosLookup[pos] = wash.ID
	lhp.PlateIDLookup[wash.ID] = pos
	return true
}

func GetLocTox(cmp *wtype.LHComponent) ([]string, error) {
	// try the cmp's own loc

	if cmp.Loc != "" {
		return strings.Split(cmp.Loc, ":"), nil
	} else {
		// try the ID of the thing

		tx, err := getSTLocTox(cmp.ID)

		if err == nil {
			return tx, err
		}

		// now try its parent

		tx, err = getSTLocTox(cmp.ParentID)

		if err == nil {
			return tx, err
		}
	}

	return []string{}, fmt.Errorf("No location found")
}

func getSTLocTox(ID string) ([]string, error) {
	st := sampletracker.GetSampleTracker()
	loc, ok := st.GetLocationOf(ID)

	if !ok {
		return []string{}, fmt.Errorf("No location found")
	}

	tx := strings.Split(loc, ":")

	if len(tx) == 2 {
		return tx, nil
	} else {
		return []string{}, fmt.Errorf("No location found")
	}
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

	out, seen = mergeToSet(lhp.Input_preferences, out, seen)
	out, _ = mergeToSet(lhp.Output_preferences, out, seen)

	return out
}

// logic of getting components:
// we look for things with the same ID
// the ID may or may not refer to an instance which is previously made
// but by this point we must have concrete locations for everything

func (lhp *LHProperties) GetComponentsSingle(cmps []*wtype.LHComponent, carryvol wunit.Volume, legacyVolume bool) ([][]string, [][]string, [][]wunit.Volume, error) {
	plateIDs := make([][]string, len(cmps))
	wellCoords := make([][]string, len(cmps))
	vols := make([][]wunit.Volume, len(cmps))

	// locally keep volumes straight

	localplates := make(map[string]*wtype.LHPlate, len(lhp.Plates))

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

	for _, pos := range lhp.Tip_preferences {
		bx, ok := lhp.Tipboxes[pos]
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
					positions[i] = pos
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

func (lhp *LHProperties) DropDirtyTips(channels []*wtype.LHChannelParameter) (wells, positions, boxtypes []string) {
	multi := len(channels)

	wells = make([]string, multi)
	positions = make([]string, multi)
	boxtypes = make([]string, multi)

	foundit := false

	for pos, bx := range lhp.Tipwastes {
		wellCoords, ok := bx.Dispose(channels)
		if ok {
			foundit = true
			for i := 0; i < multi; i++ {
				if channels[i] != nil {
					wells[i] = wellCoords[i].FormatA1()
					positions[i] = pos
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

//GetMaterialType implement stockableMaterial
func (lhp *LHProperties) GetMaterialType() material.MaterialType {
	return lhp.MaterialType
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
		logger.Info(fmt.Sprint("RemoveComponent ", plateID, " ", well, " ", volume.ToString(), " can't find plate"))
		return false
	}

	r := p.RemoveComponent(well, volume)

	if r == nil {
		logger.Info(fmt.Sprint("CAN'T REMOVE COMPONENT ", plateID, " ", well, " ", volume.ToString()))
		return false
	}

	/*
		w := p.Wellcoords[well]

		if w == nil {
			logger.Info(fmt.Sprint("RemoveComponent ", plateID, " ", well, " ", volume.ToString(), " can't find well"))
			return false
		}

		c:=w.Remove(volume)

		if c==nil{
			logger.Info(fmt.Sprint("RemoveComponent ", plateID, " ", well, " ", volume.ToString(), " can't find well"))
			return false
		}
	*/

	return true
}

// RemoveUnusedAutoallocatedComponents removes any autoallocated component wells
// that didn't end up getting used
// In direct translation to component states that
// means any components that are temporary _and_ autoallocated.
func (lhp *LHProperties) RemoveUnusedAutoallocatedComponents() {
	ids := make([]string, 0, 1)
	for _, p := range lhp.Plates {
		if p.IsTemporary() && p.IsAutoallocated() {
			ids = append(ids, p.ID)
			continue
		}

		for _, w := range p.Wellcoords {
			if w.IsTemporary() && w.IsAutoallocated() {
				w.Clear()
			}
		}
	}

	for _, id := range ids {
		lhp.RemovePlateWithID(id)
	}

	// good
}
func (lhp *LHProperties) GetEnvironment() wtype.Environment {
	// static to start with

	return wtype.Environment{
		Temperature:         wunit.NewTemperature(25, "C"),
		Pressure:            wunit.NewPressure(100000, "Pa"),
		Humidity:            0.35,
		MeanAirFlowVelocity: wunit.NewVelocity(0, "m/s"),
	}
}

func (lhp *LHProperties) Evaporate(t time.Duration) []wtype.VolumeCorrection {
	// TODO: proper environmental calls
	env := lhp.GetEnvironment()
	ret := make([]wtype.VolumeCorrection, 0, 5)
	for _, v := range lhp.Plates {
		ret = append(ret, v.Evaporate(t, env)...)
	}

	return ret
}

// TODO -- allow drivers to provide relevant constraint info... not all positions
// can be used for tip loading
func (lhp *LHProperties) CheckTipPrefCompatibility(prefs []string) bool {
	// no new tip preferences allowed for now
	if lhp.Mnfr == "CyBio" {
		if lhp.Model == "Felix" {
			for _, v := range prefs {
				return wutil.StrInStrArray(v, lhp.Tip_preferences)
			}
		} else if lhp.Model == "GeneTheatre" {
			for _, v := range prefs {
				if !wutil.StrInStrArray(v, lhp.Tip_preferences) {
					return false
				}
			}
			return true
		}

	} else if lhp.Mnfr == "Tecan" {
		// fall through
		return lhp.CheckPreferenceCompatibility(prefs)
	}

	return true
}

// CheckPreferenceCompatibility returns if the device specific configuration
// positions are compatible with the current device.
func (lhp *LHProperties) CheckPreferenceCompatibility(prefs []string) bool {
	// TODO: Not the most portable or extensible way

	var checkFn func(string) bool

	if lhp.Mnfr == "Tecan" {
		checkFn = func(pos string) bool {
			return strings.HasPrefix(pos, "TecanPos_")
		}
	} else if lhp.Mnfr == "Gilson" || lhp.Mnfr == "CyBio" && lhp.Model == "Felix" {
		checkFn = func(pos string) bool {
			return strings.HasPrefix(pos, "position_")
		}
	} else if lhp.Mnfr == "CyBio" && lhp.Model == "GeneTheatre" {
		checkFn = func(pos string) bool {
			if len(pos) != 2 {
				return false
			}
			return 'A' <= pos[0] && pos[0] <= 'D' && '0' <= pos[1] && pos[1] <= '9'
		}
	} else if lhp.Mnfr == "Labcyte" {
		checkFn = func(pos string) bool {
			return false
		}
	}

	if checkFn == nil {
		return true
	}

	for _, p := range prefs {
		if !checkFn(p) {
			return false
		}
	}

	return true
}

type UserPlate struct {
	Plate    *wtype.LHPlate
	Position string
}
type UserPlates []UserPlate

func (p *LHProperties) SaveUserPlates() UserPlates {
	up := make(UserPlates, 0, len(p.Positions))

	for pos, plate := range p.Plates {
		if plate.IsUserAllocated() {
			up = append(up, UserPlate{Plate: plate.DupKeepIDs(), Position: pos})
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
	if len(p.HeadsLoaded) == 0 {
		return wunit.ZeroVolume()
	}
	minvol := p.HeadsLoaded[0].GetParams().Minvol
	for _, head := range p.HeadsLoaded {
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

func (p *LHProperties) MinCurrentVolume() wunit.Volume {
	if len(p.HeadsLoaded) == 0 {
		return wunit.ZeroVolume()
	}

	if len(p.Tips) == 0 {
		return p.MinPossibleVolume()
	}

	minvol := p.HeadsLoaded[0].GetParams().Maxvol
	for _, head := range p.HeadsLoaded {
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
	r := dupStrArr(p.Input_preferences)

	for _, pr := range p.Output_preferences {
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

func (p *LHProperties) UpdateComponentIDs(updates map[string]*wtype.LHComponent) {
	for s, c := range updates {
		p.UpdateComponentID(s, c)
	}
}

func (p *LHProperties) UpdateComponentID(from string, to *wtype.LHComponent) bool {
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
