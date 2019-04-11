// liquidhandling/lhtypes.Go: Part of the Antha language
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

package wtype

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/antha/anthalib/wutil"
	"github.com/antha-lang/antha/graph"
)

const InPlaceMarker = "-INPLACE"

// LHComponent is an alias for Liquid to preserve backwards compatibility
// Liquid describes a liquid component and its desired properties
type LHComponent = Liquid

// Liquid is the principle liquid handling type in Antha.
// A liquid describes a liquid component and its desired properties
type Liquid struct {
	ID                 string
	BlockID            BlockID
	DaughtersID        map[string]struct{}
	ParentID           string
	Inst               string
	Order              int
	CName              string
	Type               LiquidType
	Vol                float64
	Conc               float64
	Vunit              string
	Cunit              string
	Tvol               float64
	Smax               float64 // maximum solubility
	Visc               float64
	StockConcentration float64
	SubComponents      ComponentList // List of all sub components in the LHComponent.
	Extra              map[string]interface{}
	Loc                string // refactor to PlateLocation
	Destination        string
	Policy             LHPolicy // Policy is where a custom liquid policy is stored
}

// AddSubComponent adds a subcomponent with concentration to a component.
// An error is returned if subcomponent is already found.
func (cmp *Liquid) AddSubComponent(subcomponent *Liquid, conc wunit.Concentration) error {
	return AddSubComponent(cmp, subcomponent, conc)
}

// AddSubComponents adds a component list to a component.
// If a conflicting sub component concentration is already present then an error will be returned.
// To overwrite all subcomponents ignoring conficts, use OverWriteSubComponents.
func (cmp *Liquid) AddSubComponents(allsubComponents ComponentList) error {
	return AddSubComponents(cmp, allsubComponents)
}

// OverwriteSubComponents Adds a component list to a component.
// Any existing component list will be overwritten.
// To add a ComponentList checking for duplicate entries, use AddSubComponents.
func (cmp *Liquid) OverwriteSubComponents(allsubComponents ComponentList) error {
	cmp.SubComponents = allsubComponents
	return nil
}

// GetSubComponents returns a component list from a component
func (cmp *Liquid) GetSubComponents() (ComponentList, error) {
	return GetSubComponents(cmp)
}

// GetConcentrationOf attempts to retrieve the concentration of subComponentName in component.
// If the component name is equal to subComponentName, the concentration of the component itself is returned.
func (c *Liquid) GetConcentrationOf(subComponentName string) (wunit.Concentration, error) {
	return getComponentConc(c, subComponentName)
}

// HasSubComponent evaluates if a sub component with subComponentName is found in component.
// If the component name is equal to subComponentName, true will be returned.
func (c *Liquid) HasSubComponent(subComponentName string) bool {
	return hasSubComponent(c, subComponentName)
}

func (cmp *Liquid) Matches(cmp2 *Liquid) bool {
	// request for a specific component
	if cmp.IsInstance() {
		if cmp.IsSample() {
			//  look for the ID of its parent (we don't allow sampling from samples yet)
			return cmp.ParentID == cmp2.ID
		} else {
			// if this is just the whole component we check for *its* Id
			return cmp.ID == cmp2.ID
		}
	} else {
		// sufficient to be of same types
		return cmp.IsSameKindAs(cmp2)
	}
}

func (lhc Liquid) GetID() string {
	return lhc.ID
}

func (lhc *Liquid) PlateLocation() PlateLocation {
	return PlateLocationFromString(lhc.Loc)
}

// WellLocation returns the well location in A1 format.
func (lhc *Liquid) WellLocation() string {
	return lhc.PlateLocation().Coords.FormatA1()
}

// SetWellLocation sets the well location to an LHComponent in A1 format.
func (lhc *Liquid) SetWellLocation(wellLocation string) error {
	location := lhc.PlateLocation()
	lhc.Loc = location.ID + ":" + wellLocation
	return nil
}

//GetClass return the class of the object
func (lhc *Liquid) GetClass() string {
	return "component"
}

//GetName the component's name
func (lhc *Liquid) GetName() string {
	if lhc == nil {
		return "nil"
	}
	return fmt.Sprintf("%v of %s", lhc.Volume(), lhc.CName)
}

//Summarize describe the component in a user friendly manner
func (lhc *Liquid) Summarize() string {
	if lhc == nil {
		return "nil"
	}

	if lhc.Vol != 0.0 {
		return fmt.Sprintf("%v of %s", lhc.Volume(), lhc.CName)
	} else if lhc.Tvol != 0.0 {
		return fmt.Sprintf("%s to %v", lhc.CName, lhc.TotalVolume())
	} else if lhc.Conc != 0.0 {
		return fmt.Sprintf("%s to %v", lhc.CName, lhc.Concentration())
	}

	if lhc.CName != "" {
		return fmt.Sprintf("unknown amount of %s", lhc.CName)
	}
	return "unknown volume of unknown component"
}

// PlateID returns the id of a plate or the empty string
func (lhc *Liquid) PlateID() string {
	loc := lhc.PlateLocation()

	if !loc.IsZero() {
		return loc.ID
	}

	return ""
}

func (lhc *Liquid) CNID() string {
	return fmt.Sprintf("CNID:%s:%s", lhc.ID, lhc.CName)
}

func (lhc *Liquid) Generation() int {
	gen, ok := lhc.Extra["Generation"]
	if !ok {
		return 0
	}

	genInt, ok := gen.(int)
	if ok {
		return genInt
	}

	genFloat, ok := gen.(float64)
	if ok {
		return int(genFloat)
	}
	return 0
}

func (lhc *Liquid) SetGeneration(i int) {
	lhc.Extra["Generation"] = i
}

func (lhc *Liquid) IsZero() bool {
	if lhc == nil || lhc.Type == "" || lhc.CName == "" || lhc.Vol < 0.0000001 {
		return true
	}
	return false
}

const SEQSKEY = "DNASequences"

// Return a sequence list from a component.
// Users should use GetDNASequences method.
func (lhc *Liquid) getDNASequences() (seqs []DNASequence, err error) {

	seqsValue, found := lhc.Extra[SEQSKEY]

	if !found {
		return seqs, fmt.Errorf("No Sequences list found")
	}

	var bts []byte

	bts, err = json.Marshal(seqsValue)
	if err != nil {
		return
	}

	err = json.Unmarshal(bts, &seqs)

	if err != nil {
		err = fmt.Errorf("Problem getting %s sequences. Sequences found: %+v; error: %s", lhc.Name(), seqsValue, err.Error())
	}

	return
}

// Add a sequence list to a component.
// Any existing component list will be overwritten.
// Users should use addDNASequence and UpdateDNASequence methods
func (lhc *Liquid) setDNASequences(seqList []DNASequence) {
	lhc.Extra[SEQSKEY] = seqList
}

// Returns the positions of any matching instances of a sequence in a slice of sequences.
// If checkSeqs is set to false, only the name will be checked;
// if checkSeqs is set to true, matching sequences with different names will also be checked.
func containsSeq(seqs []DNASequence, seq DNASequence, checkSeqs bool) (bool, []int) {

	var positionsFound []int

	for i := range seqs {
		if !checkSeqs {
			if seqs[i].Name() == seq.Name() {
				positionsFound = append(positionsFound, i)
			}
		} else {
			if seqs[i].Name() == seq.Name() {
				positionsFound = append(positionsFound, i)
			} else if strings.EqualFold(seqs[i].Sequence(), seq.Sequence()) && seqs[i].Plasmid == seq.Plasmid {
				positionsFound = append(positionsFound, i)
			}
		}
	}

	if len(positionsFound) > 0 {
		return true, positionsFound
	}

	return false, positionsFound
}

const (
	FORCE bool = true // Optional parameter to use in AddDNASequence method to override error check preventing addition of a duplicate sequence.
)

// AddDNASequence adds DNASequence to the LHComponent.
// If a Sequence already exists an error is returned and the sequence is not added
// unless an additional boolean argument (FORCEADD or true) is specified to ignore duplicates.
// A warning will be returned in either case if a duplicate sequence is already found.
func (lhc *Liquid) AddDNASequence(seq DNASequence, options ...bool) error {
	var err error
	// skip error checking: if no sequence list is present one will be created later anyway
	seqList, _ := lhc.getDNASequences() // nolint

	if _, positions, err := lhc.FindDNASequence(seq); err == nil {

		if len(options) == 0 {
			err = fmt.Errorf("LHComponent %s already contains sequence %s at positions %+v in sequences %+v. To add the sequence anyway add FORCE as an argument when using AddDNASequence: i.e. AddDNASequence(sequence, wtype.FORCE)", lhc.Name(), seq.Name(), positions, seqList)
			return err
		} else if !options[0] {
			err = fmt.Errorf("LHComponent %s already contains sequence %s at positions %+v in sequences %+v. To add the sequence anywayadd FORCE as an argument when using AddDNASequence: i.e. AddDNASequence(sequence, wtype.FORCE)", lhc.Name(), seq.Name(), positions, seqList)
			return err
		}
		// else if options[0] {
		// 	err = fmt.Errorf("Warning: LHComponent %s already contains sequence %s at positions %+v in sequences %+v but was added.", lhc.Name(), seq.Name(), positions, seqList)
		// }
	}

	seqList = append(seqList, seq)
	lhc.setDNASequences(seqList)

	return err
}

// SetDNASequences adds a set of DNASequences to the LHComponent.
// If a Sequence already exists an error is returned and the sequence is not added
// unless an additional boolean argument (FORCEADD or true) is specified to ignore duplicates.
// A warning will be returned in either case if a duplicate sequence is already found.
func (lhc *Liquid) SetDNASequences(seqs []DNASequence, options ...bool) error {
	var errs []string

	for _, seq := range seqs {
		err := lhc.AddDNASequence(seq, options...)
		if err != nil {
			errs = append(errs, err.Error())
		}
	}
	if len(errs) > 0 {
		return fmt.Errorf("errors setting DNASequences to component: %s", fmt.Errorf(strings.Join(errs, ";")))
	}
	return nil
}

// FindDNASequence searches for the presence of a DNASequence in the LHComponent.
// Search is based upon both name of the sequence and sequence.
// If multiple copies of the sequence exists and error is returned.
// If a Sequence does not exist, the sequence is added and an error is returned.
func (lhc *Liquid) FindDNASequence(seq DNASequence) (seqs []DNASequence, positions []int, err error) {

	seqList, err := lhc.getDNASequences()

	if err != nil {
		return
	}
	var found bool
	found, positions = containsSeq(seqList, seq, true)

	if !found {
		err = fmt.Errorf("Sequence %s not found associated with %s.", seq.Name(), lhc.Name())
		return
	}
	for i := range positions {
		seqs = append(seqs, seqList[i])
	}

	return
}

// UpdateDNASequence replaces an existing DNASequence to the LHComponent.
// Search is based upon both name of the sequence and sequence.
// If multiple copies of the sequence exists and error is returned.
// If a Sequence does not exist, the sequence is added and an error is returned.
func (lhc *Liquid) UpdateDNASequence(seq DNASequence) error {

	seqList, err := lhc.getDNASequences()

	if err != nil {
		return err
	}

	if seqs, positions, err := lhc.FindDNASequence(seq); err == nil {
		if len(positions) > 1 {
			return fmt.Errorf("LHComponent %s contains multiple instances of sequence %s  at positions %+v: %+v", lhc.Name(), seq.Name(), positions, seqs)
		}
		if len(positions) == 1 {
			seqList[positions[0]] = seq
			lhc.setDNASequences(seqList)
		}
	}

	err = lhc.AddDNASequence(seq)

	if err != nil {
		return err
	}

	return fmt.Errorf("Sequence %s did not previously exist in %s so added.", seq.Name(), lhc.Name())
}

func deleteSeq(seqList []DNASequence, position int) (newseqList []DNASequence, err error) {

	if position >= len(seqList) {
		return seqList, fmt.Errorf("Cannot delete sequence from position in %d in sequence list as list only contains %d entries", position, len(seqList))
	}

	if position == 0 {
		if len(seqList) > 1 {
			newseqList = append(seqList[position+1:])
			return
		} else {
			return []DNASequence{}, nil
		}
	} else if position == len(seqList)-1 {
		newseqList = append(seqList[:position-1])
		return
	} else {
		newseqList = append(seqList[:position], seqList[position+1:]...)
		return
	}

}

// RemoveDNASequence removes an existing DNASequence from the LHComponent.
// Search is based upon both name of the sequence and sequence.
// If multiple copies of the sequence exists and error is returned.
// If a Sequence does not exist, the sequence is added and an error is returned.
func (lhc *Liquid) RemoveDNASequence(seq DNASequence) error {

	seqList, err := lhc.getDNASequences()

	if err != nil {
		return err
	}

	if seqs, positions, err := lhc.FindDNASequence(seq); err == nil {
		if len(positions) > 1 {
			return fmt.Errorf("LHComponent %s contains multiple instances of sequence %s  at positions %+v: %+v", lhc.Name(), seq.Name(), positions, seqs)
		}
		if len(positions) == 1 {
			seqList, err = deleteSeq(seqList, positions[0])
			if err != nil {
				return err
			}
			lhc.setDNASequences(seqList)
		}
	}

	return fmt.Errorf("Sequence %s did not previously exist in %s so could not be deleted.", seq.Name(), lhc.Name())
}

// RemoveDNASequenceAtPosition removes a DNA sequence from a specific position.
// Designed for cases where FindDNASequnce() method returns multiple instances of the dna sequence.
func (lhc *Liquid) RemoveDNASequenceAtPosition(position int) error {

	seqList, err := lhc.getDNASequences()

	if err != nil {
		return err
	}

	seqList, err = deleteSeq(seqList, position)
	if err != nil {
		return err
	}

	lhc.setDNASequences(seqList)
	return nil

}

// RemoveDNASequences removes all DNASequences from the component.
func (lhc *Liquid) RemoveDNASequences() error {
	lhc.setDNASequences([]DNASequence{})
	return nil
}

// DNASequences returns DNA Sequences asociated with an LHComponent.
// An error is also returned indicating whether a sequence was found.
func (lhc *Liquid) DNASequences() ([]DNASequence, error) {
	return lhc.getDNASequences()
}

// SetVolume adds a volume to the component
func (lhc *Liquid) SetVolume(v wunit.Volume) {
	lhc.Vol = v.RawValue()
	lhc.Vunit = v.Unit().PrefixedSymbol()
}

func (lhc *Liquid) HasParent(s string) bool {
	return strings.Contains(lhc.ParentID, s)
}

func (lhc *Liquid) HasDaughter(id string) bool {
	_, found := lhc.DaughtersID[id]
	return found
}

// Name returns the component name as a string
func (lhc *Liquid) Name() string {
	return lhc.CName
}

// SetName adds the specified component name to the component.
func (lhc *Liquid) SetName(name string) {
	lhc.CName = trimString(name)
}

// MeaningfulName returns the name of the liquid if one has been explicitly set (e.g. with SetName)
// or the empty string if only an autogenerated name is available
func (l *Liquid) MeaningfulName() string {
	// use l.CName as name if it wasn't autogenerated from subcomponents
	if len(l.SubComponents.Components) == 0 || l.CName != ReturnNormalisedComponentName(l) {
		return l.CName
	}
	return ""
}

// TypeName returns the PolicyName of the LHComponent's LiquidType as a string
func (lhc *Liquid) TypeName() string {
	return string(lhc.Type)
}

// PolicyName returns the PolicyName of the LHComponent's LiquidType
func (lhc *Liquid) PolicyName() PolicyName {
	return PolicyName(lhc.TypeName())
}

// SetPolicyName adds the LiquidType associated with a PolicyName to the LHComponent.
// If the PolicyName is invalid and the DoNotPermitCustomPolicies option is used as an argument then an error is returned.
// By default, custom policyNames may be added and the validity of these will be checked later when robot instructions are generated, rather than in the element.
func (lhc *Liquid) SetPolicyName(policy PolicyName, options ...PolicyOption) error {
	liquidType, err := LiquidTypeFromString(policy, options...)
	lhc.Type = liquidType
	return err
}

// PolicyOption allows specification of advanced options to feed into the SetPolicyName method.
type PolicyOption string

// DoNotPermitCustomPolicies is an option to pass into SetPolicyName to ensure only valid system policies are specified.
// With this flag set, custom user policies are not permitted.
var DoNotPermitCustomPolicies PolicyOption = "DoNotPermitCustomPolicies"

// ModifyLHPolicyParameter specifies that this LHComponent or instance of the LHComponent should be handled with a modified
// LHPolicy parameter.
// e.g. to Change number of post mixes to 5:
// lhc.ModifyLHPolicyParameter("POST_MIX", 5)
// Valid parameters and value types are specified in aparam.go
// An error is returned if an invalid parameter or value type for that parameter is specified.
func (lhc *Liquid) ModifyLHPolicyParameter(parameter string, value interface{}) error {
	if lhc.Policy == nil || len(lhc.Policy) == 0 {
		lhc.Policy = make(LHPolicy)
	}

	return lhc.Policy.Set(parameter, value)
}

// Volume returns the Volume of the LHComponent
func (lhc *Liquid) Volume() wunit.Volume {
	if lhc == nil || (lhc.Vunit == "" && lhc.Vol == 0.0) {
		return wunit.NewVolume(0.0, "ul")
	}
	return wunit.NewVolume(lhc.Vol, lhc.Vunit)
}

func (lhc *Liquid) TotalVolume() wunit.Volume {
	if lhc.Vunit == "" && lhc.Tvol == 0.0 {
		return wunit.NewVolume(0.0, "ul")
	}
	return wunit.NewVolume(lhc.Tvol, lhc.Vunit)
}

func (lhc *Liquid) Remove(v wunit.Volume) wunit.Volume {
	v2 := lhc.Volume()

	if v2.LessThan(v) {
		lhc.Vol = (0.0)
		return v2
	}

	lhc.Vol -= v.ConvertToString(lhc.Vunit)

	return v
}

func (lhc *Liquid) Sample(v wunit.Volume) (*Liquid, error) {
	if lhc.IsZero() {
		return nil, fmt.Errorf("Cannot sample empty component")
	} else if lhc.Volume().EqualTo(v) {
		// not setting sample?!
		return lhc, nil
	}

	c := lhc.Dup()
	c.ID = NewUUID()
	v2 := lhc.Remove(v)
	c.Vunit = v2.Unit().PrefixedSymbol()
	c.Vol = v2.RawValue()
	c.AddParentComponent(lhc)
	lhc.AddDaughterComponent(c)
	c.Loc = ""
	c.Destination = ""
	c.Extra = lhc.GetExtra()

	return c, nil
}

func (lhc *Liquid) Cp() *Liquid {
	c := lhc.Dup()
	c.ID = GetUUID()
	return c
}

func (lhc *Liquid) Dup() *Liquid {
	c := NewLHComponent()
	if lhc != nil {
		c.ID = lhc.ID
		c.Order = lhc.Order
		c.CName = lhc.CName
		c.Type = lhc.Type
		c.Vol = lhc.Vol
		c.Conc = lhc.Conc
		c.Cunit = lhc.Cunit
		c.Vunit = lhc.Vunit
		c.Tvol = lhc.Tvol
		c.Smax = lhc.Smax
		c.Visc = lhc.Visc
		c.StockConcentration = lhc.StockConcentration
		c.Extra = make(map[string]interface{}, len(lhc.Extra))
		for k, v := range lhc.Extra {
			c.Extra[k] = v
		}

		c.SubComponents = lhc.SubComponents.Dup()

		c.Loc = lhc.Loc
		c.Destination = lhc.Destination
		c.ParentID = lhc.ParentID
		c.DaughtersID = make(map[string]struct{}, len(lhc.DaughtersID))
		for k, v := range lhc.DaughtersID {
			c.DaughtersID[k] = v
		}
	}
	return c
}

func (cmp *Liquid) SetSample(flag bool) bool {
	if cmp == nil {
		return false
	}

	if cmp.Extra == nil {
		cmp.Extra = make(map[string]interface{})
	}

	cmp.Extra["IsSample"] = flag

	return true
}

func (cmp *Liquid) IsSample() bool {
	if cmp == nil {
		return false
	}

	f, ok := cmp.Extra["IsSample"]

	if !ok || !f.(bool) {
		return false
	}

	return true
}

func (cmp *Liquid) HasAnyParent() bool {
	return cmp.ParentID != ""
}

// XXX XXX XXX --> This is no longer consistent... need to revise urgently
/*
func (cmp *LHComponent) AddParentComponent(cmp2 *LHComponent) {
	if cmp == nil {
		return
	}
	cmp.ParentID = cmp2.ID
}
*/

func (cmp *Liquid) AddParentComponent(cmp2 *Liquid) {
	cmp.ParentID = cmp2.ID
}

func (cmp *Liquid) AddDaughterComponent(cmp2 *Liquid) {
	if cmp.DaughtersID == nil {
		cmp.DaughtersID = make(map[string]struct{})
	}
	cmp.DaughtersID[cmp2.ID] = struct{}{}
}

func (cmp *Liquid) ReplaceDaughterID(ID1, ID2 string) {
	if _, found := cmp.DaughtersID[ID1]; found {
		delete(cmp.DaughtersID, ID1)
		cmp.DaughtersID[ID2] = struct{}{}
	}
}

func (cmp *Liquid) MixPreserveTvol(cmp2 *Liquid) {
	cmp.Mix(cmp2)
	if cmp2.Vol == 0.00 && cmp2.Tvol > 0.00 {
		vcmp := wunit.NewVolume(cmp.Vol, cmp.Vunit)
		vcmp2 := wunit.NewVolume(cmp2.Tvol, cmp2.Vunit)
		vcmp.SetValue(vcmp2.ConvertToString(cmp.Vunit))
		cmp.Vol = vcmp.RawValue() // same units
		cmp.Tvol = vcmp2.ConvertToString(cmp.Vunit)
	} else if cmp.Tvol > 0.00 {
		cmp.Vol = cmp.Tvol
	}
}

// add cmp2 to cmp
func (cmp *Liquid) Mix(cmp2 *Liquid) {
	wasEmpty := cmp.IsZero()
	cmp.Smax = mergeSolubilities(cmp, cmp2)
	// determine type of final
	cmp.Type = mergeTypes(cmp, cmp2)
	// add cmp2 to cmp

	// add parent IDs... this is recursive

	/*
		if !wasEmpty {
			cmp.AddParentComponent(cmp)
		}
	*/
	cmp.AddParentComponent(cmp2)
	//	cmp.ID = "component-" + GetUUID()
	cmp.ID = GetUUID()
	cmp2.AddDaughterComponent(cmp)

	if wasEmpty {
		if !cmp.HasConcentration() {
			cmp.SetConcentration(cmp2.Concentration())
		}
		if len(cmp.SubComponents.Components) == 0 && len(cmp2.SubComponents.Components) > 0 {
			updateSubComponentsOnly(cmp, cmp2) //nolint
		}
		cmp.SetName(cmp2.Name())
	} else {
		UpdateComponentDetails(cmp, cmp, cmp2) //nolint
	}

	vcmp := wunit.NewVolume(cmp.Vol, cmp.Vunit)
	vcmp2 := wunit.NewVolume(cmp2.Vol, cmp2.Vunit)
	vcmp.Add(vcmp2)
	cmp.Vol = vcmp.RawValue() // same units

	// result should not be a sample

	cmp.SetSample(false)

}

// @implement Liquid
// @deprecate Liquid

func (lhc *Liquid) GetSmax() float64 {
	return lhc.Smax
}

func (lhc *Liquid) GetVisc() float64 {
	return lhc.Visc
}

func (lhc *Liquid) GetExtra() map[string]interface{} {
	x := make(map[string]interface{}, len(lhc.Extra))

	// shallow copy only...
	for k, v := range lhc.Extra {
		x[k] = v
	}

	return x
}

func (lhc *Liquid) GetConc() float64 {
	return lhc.Conc
}

func (lhc *Liquid) GetCunit() string {
	return lhc.Cunit
}

// Concentration returns the Concentration of the LHComponent
func (lhc *Liquid) Concentration() (conc wunit.Concentration) {
	if lhc.Conc == 0.0 && lhc.Cunit == "" {
		return wunit.NewConcentration(0.0, "g/L")
	}
	return wunit.NewConcentration(lhc.Conc, lhc.Cunit)
}

// HasConcentration checks whether a Concentration is set for the LHComponent
func (lhc *Liquid) HasConcentration() bool {
	if lhc.Conc != 0.0 && lhc.Cunit != "" {
		return true
	}
	return false
}

// SetConcentration sets a concentration to an LHComponent; assumes conc is valid; overwrites existing concentration
func (lhc *Liquid) SetConcentration(conc wunit.Concentration) {
	lhc.Conc = conc.RawValue()
	lhc.Cunit = conc.Unit().PrefixedSymbol()
}

func (lhc *Liquid) GetVunit() string {
	return lhc.Vunit
}

func (lhc *Liquid) GetType() string {
	typeName, _ := LiquidTypeName(lhc.Type)
	return typeName.String()
}

func NewLHComponent() *Liquid {
	//lhc.ID = "component-" + GetUUID()
	return &Liquid{
		ID:          GetUUID(),
		DaughtersID: make(map[string]struct{}),
		Vunit:       "ul",
		Policy:      make(map[string]interface{}),
		Extra:       make(map[string]interface{}),
	}
}

//Clean the component to its initial state
func (cmp *Liquid) Clean() {
	cmp.Vunit = "ul"
	cmp.DaughtersID = make(map[string]struct{})
	cmp.ParentID = ""
	cmp.Inst = ""
	cmp.Order = 0
	cmp.CName = ""
	cmp.Type = LiquidType("")
	cmp.Vol = 0.0
	cmp.Conc = 0.0
	cmp.Vunit = "ul"
	cmp.Cunit = ""
	cmp.Tvol = 0.0
	cmp.Smax = 0.0
	cmp.Visc = 0.0
	cmp.StockConcentration = 0.0
	cmp.SubComponents = ComponentList{}
	cmp.Extra = make(map[string]interface{})
	cmp.Loc = ""
	cmp.Destination = ""
	cmp.Policy = make(map[string]interface{})
}

func (cmp *Liquid) String() string {
	id := cmp.ID

	l := cmp.Loc

	v := fmt.Sprintf("%-6.3f:%s", cmp.Vol, cmp.Vunit)

	if l == "" {
		l = "NOPLATE:NOWELL"
	}

	return id + ":" + cmp.CName + ":" + l + ":" + v
}

func (cmp *Liquid) ParentTree() graph.StringGraph {
	g := graph.StringGraph{Nodes: make([]string, 0, 3), Outs: make(map[string][]string)}
	parseTree(cmp.ID+"("+cmp.ParentID+")", &g)
	return g
}

// graphviz format
func (cmp *Liquid) ParentTreeString() string {
	g := cmp.ParentTree()
	s := graph.Print(graph.PrintOpt{Graph: &g})
	return s
}

//   a(b_c_d)_e()_f(g_h)
//   nodes: [abcdefgh]
//   outs : a:[] b[a] c[a] d[a] e[] f[] g[f] h[f]
//

func parseTree(p string, g *graph.StringGraph) []string {
	newnodes := make([]string, 0, 3)
	if p[0] == '(' {
		// strip brackets
		p = p[1 : len(p)-1]
	}

	if len(p) == 0 {
		// empty bracket pair
		return newnodes
	}

	bc := 0

	splits := make([]int, 0, len(p))

	for i, c := range p {
		if c == '(' {
			bc += 1
			continue
		} else if c == ')' {
			bc -= 1
			continue
		}

		if bc == 0 && c == '_' {
			splits = append(splits, i)
		}
		//   a(b()_c()_d())_e()_f(g()_h())
		//                 s   s
	}

	splits = append(splits, len(p))

	// carve up

	b := 0

	for _, e := range splits {
		tok := p[b:e]
		lb := strings.Index(tok, "(")
		node := tok[:lb]
		if !wutil.StrInStrArray(node, g.Nodes) {
			g.Nodes = append(g.Nodes, node)
			g.Outs[node] = make([]string, 0, 3)
			newnodes = append(newnodes, node)
		}

		childnodes := parseTree(tok[lb:], g)

		for _, child := range childnodes {
			g.Outs[child] = append(g.Outs[child], node)
		}
		b = e + 1
	}

	return newnodes
}

func (lhc *Liquid) AddVolumeRule(minvol, maxvol float64, pol LHPolicy) error {
	lhpr, err := lhc.GetPolicies()

	if err != nil {
		return err
	}

	rulenum := len(lhpr.Rules)

	name := fmt.Sprintf("UserRule%d", rulenum+1)

	rule := NewLHPolicyRule(name)
	err = rule.AddNumericConditionOn("VOLUME", minvol, maxvol)
	if err != nil {
		return err
	}
	lhpr.AddRule(rule, pol)

	err = rule.AddCategoryConditionOn("INSTANCE", lhc.ID)
	if err != nil {
		return err
	}

	err = lhc.SetPolicies(lhpr)

	return err
}

func (lhc *Liquid) AddPolicy(pol LHPolicy) error {
	lhpr, err := lhc.GetPolicies()

	if err != nil {
		return err
	}

	rulenum := len(lhpr.Rules)

	name := fmt.Sprintf("UserRule%d", rulenum+1)

	rule := NewLHPolicyRule(name)
	err = rule.AddCategoryConditionOn("INSTANCE", lhc.ID)
	if err != nil {
		return err
	}
	lhpr.AddRule(rule, pol)

	err = lhc.SetPolicies(lhpr)

	return err

}

// in future this will be deprecated... should not let user completely reset policies
func (lhc *Liquid) SetPolicies(rs *LHPolicyRuleSet) error {
	buf, err := json.Marshal(rs)

	if err == nil {
		lhc.Extra["Policies"] = string(buf)
	}

	return err
}

func (lhc *Liquid) GetPolicies() (*LHPolicyRuleSet, error) {
	var rs LHPolicyRuleSet
	var err error

	if lhc.Extra == nil {
		return NewLHPolicyRuleSet(), err
	}

	ent, ok := lhc.Extra["Policies"]

	if !ok {
		return NewLHPolicyRuleSet(), err
	}

	s, ok := ent.(string)

	if !ok {
		err = fmt.Errorf("Wrong type for policies entry (%v)", ent)
		return &rs, err
	}

	ba := []byte(s)

	err = json.Unmarshal(ba, &rs)
	return &rs, err
}

func (lhc *Liquid) IsValuable() bool {
	if lhc.Extra == nil {
		return false
	}

	v, ok := lhc.Extra["valuable"]

	if !ok {
		return false
	}

	b, ok := v.(bool)

	if !ok {
		return false
	}

	return b
}

func (lhc *Liquid) SetValue(b bool) {
	if lhc.Extra == nil {
		lhc.Extra = make(map[string]interface{})
	}

	lhc.Extra["valuable"] = b
}

const instanceMarker = "INSTANCE"

func (lhc *Liquid) DeclareInstance() {
	// everything starts off as a Type
	// instancehood must inherit

	if lhc != nil {
		if lhc.Extra == nil {
			lhc.Extra = make(map[string]interface{})
		}

		lhc.Extra[instanceMarker] = true
	}
}

func (lhc *Liquid) IsInstance() bool {
	if lhc == nil || lhc.Extra == nil {
		return false
	}

	temp, ok := lhc.Extra[instanceMarker]

	if !ok {
		return false
	}

	b, ok := temp.(bool)

	if !ok {
		panic(fmt.Sprintf("Improper instance marker setting - please do not use %s as a map key in Extra! Currently set to %v", instanceMarker, b))
	}

	return b
}

func (lhc *Liquid) DeclareNotInstance() {
	// explicitly set instance status to false

	lhc.DeclareInstance() // lazy: make sure instance status is initialised
	lhc.Extra[instanceMarker] = false
}

func (lhc *Liquid) IsSameKindAs(c2 *Liquid) bool {
	// v0: amounts to same CName

	return lhc.Kind() == c2.Kind()

	// v1: Explicit kind IDs separate from names (TODO)
}

func (lhc *Liquid) Kind() string {
	// v0: it's the name
	return lhc.CName

	// v1: distinct IDs for underlying liquid types
}

func (cmp Liquid) IDOrName() string {
	// as below but omits kind name to allow users to reset

	if cmp.IsInstance() {
		if cmp.IsSample() {
			return cmp.ParentID
		} else {
			return cmp.ID
		}

	} else {
		return cmp.Kind()
	}

}

func (cmp Liquid) FullyQualifiedName() string {
	// this should be equivalent to the checks done by LHWell.Contains()

	if cmp.IsInstance() {
		if cmp.IsSample() {
			return cmp.ParentID + ":" + cmp.Kind()
		} else {
			return cmp.ID + ":" + cmp.Kind()

		}

	} else {
		return cmp.Kind()
	}
}

func (l1 *Liquid) EqualTypeVolume(l2 *Liquid) bool {
	return l1.CName == l2.CName && l1.Volume().EqualTo(l2.Volume())
}

func (l1 *Liquid) EqualTypeVolumeID(l2 *Liquid) bool {
	if !l1.EqualTypeVolume(l2) {
		return false
	}

	return l1.ID == l2.ID
}
