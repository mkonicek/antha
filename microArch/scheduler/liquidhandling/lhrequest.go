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
	OutputOrder           []string
	OutputIteratorFactory func(wtype.Addressable) wtype.AddressIterator `json:"-"`
	InstructionChain      *IChain
	TimeEstimate          float64
	CarryVolume           wunit.Volume
	InstructionSets       [][]*wtype.LHInstruction
	Evaps                 []wtype.VolumeCorrection
	Options               LHOptions
	NUserPlates           int
	Output_sort           bool
	TipsUsed              []wtype.TipEstimate
	InputSolutions        *InputSolutions //store properties related to the Liquids for the request
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

//GetSolutionsFromInputPlates get all the solutions available to the mix task
//in the input plates
func (req *LHRequest) GetSolutionsFromInputPlates() (map[string][]*wtype.Liquid, error) {

	inputs := make(map[string][]*wtype.Liquid)

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
				//get the amount available
				c.SetVolume(w.CurrentWorkingVolume())
				inputs[c.CName] = append(inputs[c.CName], c)
			}
		}
	}

	return inputs, nil
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
	lhr := &LHRequest{
		ID:                 wtype.GetUUID(),
		LHInstructions:     make(map[string]*wtype.LHInstruction),
		Plates:             make(map[string]*wtype.Plate),
		InstructionSet:     liquidhandling.NewRobotInstructionSet(nil),
		Input_assignments:  make(map[string][]string),
		Output_assignments: make(map[string][]string),
		Input_plates:       make(map[string]*wtype.Plate),
		Output_plates:      make(map[string]*wtype.Plate),
		Input_setup_weights: map[string]float64{
			"MAX_N_PLATES":           2,
			"MAX_N_WELLS":            96,
			"RESIDUAL_VOLUME_WEIGHT": 1.0,
		},
		Plate_lookup:          make(map[string]string),
		Stockconcs:            make(map[string]wunit.Concentration),
		OutputIteratorFactory: columnWiseIterator,
		CarryVolume:           wunit.NewVolume(0.5, "ul"),
		Options:               NewLHOptions(),
	}
	systemPolicies, _ := wtype.GetSystemLHPolicies()
	lhr.SetPolicies(systemPolicies)
	return lhr
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

//NewComponentsAdded run this after Plan to determine if anything
// new was added to the inputs
func (lhr *LHRequest) NewComponentsAdded() bool {
	return len(lhr.InputSolutions.VolumesWanting) != 0
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

//GetOrderedLHInstructions get the LHInstructions in the order which should have
//previously been detrmined by setOutputOrder
func (request *LHRequest) GetOrderedLHInstructions() ([]*wtype.LHInstruction, error) {
	//if these aren't the same length then there was probably an issue with setOutputOrder
	if len(request.OutputOrder) != len(request.LHInstructions) {
		return nil, errors.Errorf("request OutputOrder has length %d but %d LHInstructions", len(request.OutputOrder), len(request.LHInstructions))
	}

	ret := make([]*wtype.LHInstruction, 0, len(request.OutputOrder))
	for _, instructionID := range request.OutputOrder {
		instruction, ok := request.LHInstructions[instructionID]
		if !ok {
			return ret, errors.Errorf("request has invalid OutputOrder, no instruction with id %s", instructionID)
		}
		ret = append(ret, instruction)
	}
	return ret, nil
}

//GetUnorderedLHInstructions get a slice containing all the LHInstructions in
//arbitrary order
func (request *LHRequest) GetUnorderedLHInstructions() []*wtype.LHInstruction {
	ret := make([]*wtype.LHInstruction, 0, len(request.LHInstructions))
	for _, v := range request.LHInstructions {
		ret = append(ret, v)
	}

	return ret
}

//updateWithNewLHInstructions make sure the request contains the new instructions if aggregation has occurred
func (rq *LHRequest) updateWithNewLHInstructions(sorted []*wtype.LHInstruction) {
	for _, ins := range sorted {
		_, ok := rq.LHInstructions[ins.ID]
		if !ok {
			rq.LHInstructions[ins.ID] = ins
		}
	}
}
