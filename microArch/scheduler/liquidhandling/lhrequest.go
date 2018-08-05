// liquidhandling/lhrequest.Go: Part of the Antha language
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

// defines types for dealing with liquid handling requests
package liquidhandling

import (
	"github.com/pkg/errors"
	"strings"

	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/microArch/driver/liquidhandling"
)

// structure for defining a request to the liquid handler
type LHRequest struct {
	ID                    string
	BlockID               wtype.BlockID
	BlockName             string
	LHInstructions        map[string]*wtype.LHInstruction
	Input_solutions       map[string][]*wtype.Liquid
	Plates                map[string]*wtype.Plate
	Tips                  []*wtype.LHTipbox
	InstructionSet        *liquidhandling.RobotInstructionSet
	Instructions          []liquidhandling.TerminalRobotInstruction
	InstructionText       string
	Input_assignments     map[string][]string
	Output_assignments    map[string][]string
	Input_plates          map[string]*wtype.Plate
	Output_plates         map[string]*wtype.Plate
	Input_platetypes      []*wtype.Plate
	Input_plate_order     []string
	Input_setup_weights   map[string]float64
	Output_platetypes     []*wtype.Plate
	Output_plate_order    []string
	Plate_lookup          map[string]string
	Stockconcs            map[string]wunit.Concentration
	PolicyManager         *LHPolicyManager
	Input_order           []string
	Output_order          []string
	OutputIteratorFactory func(wtype.Addressable) wtype.AddressIterator `json:"-"`
	InstructionChain      *IChain
	Input_vols_supplied   map[string]wunit.Volume
	Input_vols_required   map[string]wunit.Volume
	Input_vols_wanting    map[string]wunit.Volume
	TimeEstimate          float64
	CarryVolume           wunit.Volume
	InstructionSets       [][]*wtype.LHInstruction
	Evaps                 []wtype.VolumeCorrection
	Options               LHOptions
	NUserPlates           int
	Output_sort           bool
	TipsUsed              []wtype.TipEstimate
}

func (req *LHRequest) GetPlate(id string) (*wtype.Plate, bool) {
	p, ok := req.Plates[id]

	if ok {
		return p, true
	}

	p, ok = req.Input_plates[id]

	if ok {
		return p, true
	}

	p, ok = req.Output_plates[id]

	if ok {
		return p, true
	}

	return nil, false
}

func (req *LHRequest) ConfigureYourself() error {
	// ensures input solutions is populated
	// once input plates are specified
	// more to happen later
	inputs := req.Input_solutions

	if inputs == nil {
		inputs = make(map[string][]*wtype.Liquid)
	}

	// we need to make an exception of components which are used literally
	// i.e. anything used in a mix-in-place; these don't add to the general
	// store of anonymous components to be sampled from

	uniques := make(map[wtype.PlateLocation]*wtype.Liquid, len(req.LHInstructions))

	for _, ins := range req.LHInstructions {
		if ins.InsType() != "MIX" {
			continue
		}
		if ins.IsMixInPlace() {
			if !ins.Components[0].PlateLocation().IsZero() {
				uniques[ins.Components[0].PlateLocation()] = ins.Components[0]
			}
			//else {
			// this will be autoallocated
			//}
		}
	}

	for _, v := range req.Input_plates {
		for _, w := range v.Wellcoords {
			if w.IsEmpty() {
				continue
			}

			// special case for components treated literally

			cmp, ok := uniques[w.PlateLocation()]

			if ok {
				ar := inputs[cmp.CNID()]
				ar = append(ar, cmp)
				inputs[cmp.CNID()] = ar
			} else {
				// bulk components (where instances don't matter) are
				// identified using just CName
				c := w.Contents().Dup()
				vvvvvv := c.Volume()
				vvvvvv.Subtract(w.ResidualVolume())
				c.SetVolume(vvvvvv)
				ar := inputs[c.CName]
				ar = append(ar, c)
				inputs[c.CName] = ar
			}
		}
	}

	req.Input_solutions = inputs
	return nil
}

// this function checks requests so we can see early on whether or not they
// are going to cause problems
func ValidateLHRequest(rq *LHRequest) (bool, string) {
	if rq.Output_platetypes == nil || len(rq.Output_platetypes) == 0 {
		return false, "No output plate type specified"
	}

	if len(rq.Input_platetypes) == 0 {
		return false, "No input plate types specified"
	}

	if rq.Policies() == nil {
		return false, "No policies specified"
	}

	return true, "OK"
}

func columnWiseIterator(a wtype.Addressable) wtype.AddressIterator {
	return wtype.NewAddressIterator(a, wtype.ColumnWise, wtype.TopToBottom, wtype.LeftToRight, false)
}

func NewLHRequest() *LHRequest {
	var lhr LHRequest
	lhr.ID = wtype.GetUUID()
	lhr.LHInstructions = make(map[string]*wtype.LHInstruction)
	lhr.Input_solutions = make(map[string][]*wtype.Liquid)
	lhr.Plates = make(map[string]*wtype.Plate)
	lhr.Tips = make([]*wtype.LHTipbox, 0, 1)
	lhr.Input_plates = make(map[string]*wtype.Plate)
	lhr.Input_platetypes = make([]*wtype.Plate, 0, 2)
	lhr.Input_setup_weights = make(map[string]float64)
	lhr.Output_plates = make(map[string]*wtype.Plate)
	lhr.Output_plate_order = make([]string, 0, 1)
	lhr.Input_plate_order = make([]string, 0, 1)
	lhr.Plate_lookup = make(map[string]string)
	lhr.Stockconcs = make(map[string]wunit.Concentration)
	lhr.Input_order = make([]string, 0)
	lhr.Output_order = make([]string, 0)
	lhr.OutputIteratorFactory = columnWiseIterator
	lhr.Output_assignments = make(map[string][]string)
	lhr.Input_assignments = make(map[string][]string)
	lhr.InstructionSet = liquidhandling.NewRobotInstructionSet(nil)
	lhr.InstructionText = ""
	lhr.Input_vols_required = make(map[string]wunit.Volume)
	lhr.Input_vols_supplied = make(map[string]wunit.Volume)
	lhr.Input_vols_wanting = make(map[string]wunit.Volume)
	lhr.CarryVolume = wunit.NewVolume(0.5, "ul")
	lhr.Input_setup_weights["MAX_N_PLATES"] = 2
	lhr.Input_setup_weights["MAX_N_WELLS"] = 96
	lhr.Input_setup_weights["RESIDUAL_VOLUME_WEIGHT"] = 1.0
	lhr.Options = NewLHOptions()
	lhr.TipsUsed = make([]wtype.TipEstimate, 0)
	systemPolicies, _ := wtype.GetSystemLHPolicies()
	lhr.SetPolicies(systemPolicies)
	return &lhr
}

func (lhr *LHRequest) Policies() *wtype.LHPolicyRuleSet {
	return lhr.PolicyManager.Policies()
}

func (lhr *LHRequest) SetPolicies(systemPolicies *wtype.LHPolicyRuleSet) {

	if systemPolicies == nil {
		panic("no system policies specified as argument to SetPolicies")
	}

	lhr.PolicyManager = &LHPolicyManager{
		SystemPolicies: systemPolicies,
	}
}

// AddUserPolicies allows policies specified in elements to be added to the PolicyManager.
func (lhr *LHRequest) AddUserPolicies(userPolicies *wtype.LHPolicyRuleSet) {
	// things coming in take precedence over things already there
	if lhr.PolicyManager.UserPolicies == nil {
		lhr.PolicyManager.UserPolicies = userPolicies
	} else {
		lhr.PolicyManager.UserPolicies.MergeWith(userPolicies)
	}
}

func (lhr *LHRequest) Add_instruction(ins *wtype.LHInstruction) {
	lhr.LHInstructions[ins.ID] = ins
}

func (lhr *LHRequest) NewComponentsAdded() bool {
	// run this after Plan to determine if anything
	// new was added to the inputs

	return len(lhr.Input_vols_wanting) != 0
}

func (lhr *LHRequest) AddUserPlate(p *wtype.Plate) {
	// impose sanity

	if p.PlateName == "" {
		p.PlateName = getSafePlateName(lhr, "user_plate", "_", lhr.NUserPlates+1)
		lhr.NUserPlates += 1
	}

	p.MarkNonEmptyWellsUserAllocated()

	lhr.Input_plates[p.ID] = p
}

func (lhr *LHRequest) UseLegacyVolume() bool {
	// magically create extra volumes for intermediates?
	return lhr.Options.LegacyVolume
}

func (lhr *LHRequest) GetPolicyManager() *LHPolicyManager {
	return lhr.PolicyManager
}

type LHPolicyManager struct {
	SystemPolicies *wtype.LHPolicyRuleSet
	UserPolicies   *wtype.LHPolicyRuleSet
}

// SetOption adds an option and value to both System and User policies in the PolicyManager.
func (mgr *LHPolicyManager) SetOption(optname string, value interface{}) error {
	if mgr.SystemPolicies != nil {
		err := mgr.SystemPolicies.SetOption(optname, value)
		if err != nil {
			return err
		}
	}
	if mgr.UserPolicies != nil {
		err := mgr.UserPolicies.SetOption(optname, value)
		if err != nil {
			return err
		}
	}
	return nil
}

func (mgr *LHPolicyManager) Policies() *wtype.LHPolicyRuleSet {
	ret := wtype.CloneLHPolicyRuleSet(mgr.SystemPolicies)

	// things coming in take precedence over things already there
	if mgr.UserPolicies == nil {
		return ret
	}
	ret.MergeWith(mgr.UserPolicies)
	return ret
}

func (mgr *LHPolicyManager) MergePolicies(protocolpolicies *wtype.LHPolicyRuleSet) *wtype.LHPolicyRuleSet {
	ret := mgr.Policies()
	ret.MergeWith(protocolpolicies)
	return ret
}

// HasPlateNamed checks if the request already contains a plate with the specified name
func (request *LHRequest) HasPlateNamed(name string) bool {
	checkForPlateNamed := func(query string, subject map[string]*wtype.Plate) bool {
		for _, plate := range subject {
			if plate.PlateName == query {
				return true
			}
		}
		return false
	}

	if checkForPlateNamed(name, request.Input_plates) {
		return true
	}
	if checkForPlateNamed(name, request.Output_plates) {
		return true
	}

	return false
}

// OrderedInputPlates returns the list of input plates in order
func (request *LHRequest) OrderedInputPlates() []*wtype.Plate {
	ret := make([]*wtype.Plate, 0, len(request.Input_plates))
	for _, id := range request.Input_plate_order {
		ret = append(ret, request.Input_plates[id])
	}

	return ret
}

// OrderedOutputPlates returns the list of input plates in order
func (request *LHRequest) OrderedOutputPlates() []*wtype.Plate {
	ret := make([]*wtype.Plate, 0, len(request.Output_plates))
	for _, id := range request.Output_plate_order {
		ret = append(ret, request.Output_plates[id])
	}

	return ret
}

// AllPlates returns a list of all known plates, in the order input plates, output plates
// ordering will be as within the stated orders of each
func (request *LHRequest) AllPlates() []*wtype.Plate {
	r := make([]*wtype.Plate, 0, len(request.Input_plates)+len(request.Output_plates))

	r = append(r, request.OrderedInputPlates()...)
	r = append(r, request.OrderedOutputPlates()...)

	return r
}

//EnsureInstructionComponentsAreUnique make certain that inputs and outputs to
//LHInstructions are not referred to elsewhere, as could be the case with poor
//element code
func (request *LHRequest) EnsureComponentsAreUnique() {
	for _, ins := range request.LHInstructions {
		for i := 0; i < len(ins.Components); i++ {
			ins.Components[i] = ins.Components[i].Dup()
		}
		ins.Results[0] = ins.Results[0].Dup()
	}
}

//assertVolumesNonNegative tests that the volumes within the LHRequest are zero or positive
func (request *LHRequest) assertVolumesNonNegative() error {
	for _, ins := range request.LHInstructions {
		if ins.Type != wtype.LHIMIX {
			continue
		}

		for _, cmp := range ins.Components {
			if cmp.Volume().LessThan(wunit.ZeroVolume()) {
				return wtype.LHErrorf(wtype.LH_ERR_VOL, "negative volume for component \"%s\" in instruction:\n%s", cmp.CName, ins.Summarize(1))
			}
		}
	}
	return nil
}

//assertTotalVolumesMatch checks that component total volumes are all the same in mix instructions
func (request *LHRequest) assertTotalVolumesMatch() error {
	for _, ins := range request.LHInstructions {
		if ins.Type != wtype.LHIMIX {
			continue
		}

		totalVolume := wunit.ZeroVolume()

		for _, cmp := range ins.Components {
			if tV := cmp.TotalVolume(); !tV.IsZero() {
				if !totalVolume.IsZero() && !tV.EqualTo(totalVolume) {
					return wtype.LHErrorf(wtype.LH_ERR_VOL, "multiple distinct total volumes specified in instruction:\n%s", ins.Summarize(1))
				}
				totalVolume = tV
			}
		}
	}
	return nil
}

//assertMixResultsCorrect checks that volumes of the mix result matches either the sum of the input, or the total volume if specified
func (request *LHRequest) assertMixResultsCorrect() error {
	for _, ins := range request.LHInstructions {
		if ins.Type != wtype.LHIMIX {
			continue
		}

		totalVolume := wunit.ZeroVolume()
		volumeSum := wunit.ZeroVolume()

		for _, cmp := range ins.Components {
			if tV := cmp.TotalVolume(); !tV.IsZero() {
				totalVolume = tV
			} else if v := cmp.Volume(); !v.IsZero() {
				volumeSum.Add(v)
			}
		}

		if len(ins.Results) != 1 {
			return wtype.LHErrorf(wtype.LH_ERR_DIRE, "mix instruction has %d results specified, expecting one at instruction:\n%s",
				len(ins.Results), ins.Summarize(1))
		}

		resultVolume := ins.Results[0].Volume()

		if !totalVolume.IsZero() && !totalVolume.EqualTo(resultVolume) {
			return wtype.LHErrorf(wtype.LH_ERR_VOL, "total volume (%v) does not match resulting volume (%v) for instruction:\n%s",
				totalVolume, resultVolume, ins.Summarize(1))
		} else if totalVolume.IsZero() && !volumeSum.EqualTo(resultVolume) {
			return wtype.LHErrorf(wtype.LH_ERR_VOL, "sum of requested volumes (%v) does not match result volume (%v) for instruction:\n%s",
				volumeSum, resultVolume, ins.Summarize(1))
		}
	}
	return nil
}

//assertWellNotOverfilled checks that mix instructions aren't going to overfill
//the wells when a plate has been chosen specified.
//assumes assertMixResultsCorrect returns nil
func (request *LHRequest) assertWellNotOverfilled() error {
	for _, ins := range request.LHInstructions {
		if ins.Type != wtype.LHIMIX {
			continue
		}

		resultVolume := ins.Results[0].Volume()

		var plate *wtype.Plate
		if ins.OutPlate != nil {
			plate = ins.OutPlate
		} else if ins.PlateID != "" {
			if p, ok := request.GetPlate(ins.PlateID); !ok {
				continue
			} else {
				plate = p
			}
		} else {
			//couldn't find an appropriate plate
			continue
		}

		if maxVol := plate.Welltype.MaxVolume(); maxVol.LessThan(resultVolume) {
			//ignore if this is just numerical precision (#campainforintegervolume)
			delta := wunit.SubtractVolumes(resultVolume, maxVol)
			if delta.IsZero() {
				continue
			}
			return wtype.LHErrorf(wtype.LH_ERR_VOL, "volume of resulting mix (%v) exceeds the well maximum (%v) for instruction:\n%s",
				resultVolume, maxVol, ins.Summarize(1))
		}
	}
	return nil
}

//AssertInstructionVolumesOK check that instruction volumes are non-negative,
//that total and result volumes are consistent, and that wells are not
//overfilled when a plate is specified
func (request *LHRequest) AssertInstructionVolumesOK() error {
	if err := request.assertVolumesNonNegative(); err != nil {
		return err
	} else if err = request.assertTotalVolumesMatch(); err != nil {
		return err
	} else if err = request.assertMixResultsCorrect(); err != nil {
		return err
	} else if err = request.assertWellNotOverfilled(); err != nil {
		return err
	}
	return nil
}

//AssertInstructionHaveDestinations make sure that destination plates have
//been chosen for all mix instructions
func (request *LHRequest) AssertInstructionsHaveDestinations() error {
	for _, ins := range request.LHInstructions {
		// non-mix instructions are fine
		if ins.Type != wtype.LHIMIX {
			continue
		}

		if ins.PlateID == "" || ins.Platetype == "" || ins.Welladdress == "" {
			return errors.Errorf("after layout all mix instructions must have plate IDs, plate types and well addresses, Found: \n INS %v: HAS PlateID %t, HAS platetype %t HAS WELLADDRESS %t",
				ins, ins.PlateID != "", ins.Platetype != "", ins.Welladdress != "")
		}
	}

	return nil
}

//OrderedInstructionsString return a string summary of the instructions in
//final order
func (request *LHRequest) OrderedInstructionsString(indent string) string {
	lines := make([]string, 0, len(request.Output_order))
	for _, insID := range request.Output_order {
		lines = append(lines, request.LHInstructions[insID].String())
	}
	return indent + strings.Join(lines, "\n"+indent)
}
