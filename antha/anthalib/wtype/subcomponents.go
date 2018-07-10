// Part of the Antha language
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

package wtype

import (
	"fmt"
	"sort"
	"strings"

	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/pubchem"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/antha/anthalib/wutil"
)

// ComponentListSample is a sample of a Component list at a specified volume
// when two ComponentListSamples are mixed a new diluted ComponentList is generated
type ComponentListSample struct {
	ComponentList
	Volume wunit.Volume
}

// MixComponentLists merges two componentListSamples.
// When two ComponentListSamples are mixed a new diluted ComponentList is generated.
// An error may be generated if two components with the same name exist within the two lists with incompatible concentration units.
// In this instance, the molecular weight for that component will be looked up in pubchem in order to change the units in both lists to g/l,
// which will be able to be added.
func MixComponentLists(sample1, sample2 ComponentListSample) (newList ComponentList, err error) {

	var errs []string

	complist := make(map[string]wunit.Concentration)

	sample1DilutionRatio := sample1.Volume.SIValue() / (sample1.Volume.SIValue() + sample2.Volume.SIValue())

	sample2DilutionRatio := sample2.Volume.SIValue() / (sample1.Volume.SIValue() + sample2.Volume.SIValue())

	for key, conc := range sample1.Components {

		newConc := wunit.MultiplyConcentration(conc, sample1DilutionRatio)

		if existingConc, found := complist[key]; found {
			sumOfConcs, newerr := wunit.AddConcentrations(newConc, existingConc)
			if newerr != nil {
				// attempt unifying base units
				molecule, newerr := pubchem.MakeMolecule(key)
				if newerr != nil {
					errs = append(errs, newerr.Error())
				} else {
					newConcG := molecule.GramPerL(newConc)
					existingConcG := molecule.GramPerL(existingConc)

					sumOfConcs, newerr = wunit.AddConcentrations(newConcG, existingConcG)
					if newerr != nil {
						errs = append(errs, newerr.Error())
					}
				}
			}
			complist[key] = sumOfConcs
		} else {
			complist[key] = newConc
		}
	}

	for key, conc := range sample2.Components {
		newConc := wunit.MultiplyConcentration(conc, sample2DilutionRatio)

		if existingConc, found := complist[key]; found {
			sumOfConcs, newerr := wunit.AddConcentrations(newConc, existingConc)
			if newerr != nil {
				// attempt unifying base units
				molecule, newerr := pubchem.MakeMolecule(key)
				if newerr != nil {
					errs = append(errs, newerr.Error())
				} else {
					newConcG := molecule.GramPerL(newConc)
					existingConcG := molecule.GramPerL(existingConc)

					sumOfConcs, newerr = wunit.AddConcentrations(newConcG, existingConcG)
					if newerr != nil {
						errs = append(errs, newerr.Error())
					}
				}
			}
			complist[key] = sumOfConcs
		} else {
			complist[key] = newConc
		}
	}
	newList.Components = complist

	if len(errs) > 0 {
		err = fmt.Errorf(strings.Join(errs, "; "))
	}

	return
}

// SimulateMix simulates the resulting list of components and concentrations
// which would be generated by mixing the samples together.
// This will only add the component name itself to the new component list if the sample has no components
// this is to prevent potential duplication since if a component has a list of sub components the name
// is considered to be an alias and the component list the true meaning of what the component is.
// If any sample concentration of zero is found the component list will be made but an error returned.
func SimulateMix(samples ...*LHComponent) (newComponentList ComponentList, mixSteps []ComponentListSample, warning error) {

	var warnings []string
	var nonZeroVols []wunit.Volume
	var forTotalVol wunit.Volume
	var topUpNeeded bool
	var topUpVolume wunit.Volume
	// top up volume will only be used if a SampleForTotalVolume command is used
	var bufferIndex int = -1

	for i, sample := range samples {
		if sample.Volume().RawValue() == 0.0 && sample.Tvol > 0 {
			forTotalVol = wunit.NewVolume(sample.Tvol, sample.Vunit)
			bufferIndex = i
			topUpNeeded = true
		}
		nonZeroVols = append(nonZeroVols, sample.Volume())
	}
	sumOfSampleVolumes := wunit.AddVolumes(nonZeroVols...)

	if !topUpNeeded {
		forTotalVol = sumOfSampleVolumes
	}

	topUpVolume = wunit.SubtractVolumes(forTotalVol, sumOfSampleVolumes)

	if topUpVolume.RawValue() < 0.0 {
		return newComponentList, mixSteps, fmt.Errorf("SampleForTotalVolume requested (%s) is less than sum of sample volumes (%s)", forTotalVol, sumOfSampleVolumes)
	}

	var volsSoFar []wunit.Volume

	for i, sample := range samples {

		var volToAdd wunit.Volume

		if i == bufferIndex {
			volToAdd = topUpVolume
		} else {
			volToAdd = sample.Volume()
		}

		if i == 0 {
			var err error
			newComponentList, err = GetSubComponents(sample)
			if err != nil {
				newComponentList.Components = make(map[string]wunit.Concentration)
				if sample.Conc == 0 {
					warnings = append(warnings, "zero concentration found for sample "+sample.Name())
					newComponentList.Components[sample.Name()] = wunit.NewConcentration(1.0, "v/v")
				} else {
					newComponentList.Components[sample.Name()] = sample.Concentration()
				}
				mixSteps = append(mixSteps, ComponentListSample{ComponentList: newComponentList, Volume: volToAdd})
			}
			volsSoFar = append(volsSoFar, volToAdd)
		}

		if i < len(samples)-1 {

			nextSample := samples[i+1]
			nextList, err := GetSubComponents(nextSample)
			if err != nil {
				nextList.Components = make(map[string]wunit.Concentration)
				if nextSample.Conc == 0 {
					warnings = append(warnings, "zero concentration found for sample "+nextSample.Name())
					nextList.Components[nextSample.Name()] = wunit.NewConcentration(1.0, "v/v")
				} else {
					nextList.Components[nextSample.Name()] = nextSample.Concentration()
				}
			}

			if i != 0 {
				volsSoFar = append(volsSoFar, volToAdd)
			}

			volOfPreviousSamples := wunit.AddVolumes(volsSoFar...)

			previousMixStep := ComponentListSample{ComponentList: newComponentList, Volume: volOfPreviousSamples}

			var nexSampleVolToAdd wunit.Volume

			if i+1 == bufferIndex {
				nexSampleVolToAdd = topUpVolume
			} else {
				nexSampleVolToAdd = nextSample.Volume()
			}
			nextMixStep := ComponentListSample{ComponentList: nextList, Volume: nexSampleVolToAdd}
			newComponentList, err = MixComponentLists(previousMixStep, nextMixStep)
			if err != nil {
				warnings = append(warnings, err.Error())
			}

			mixSteps = append(mixSteps, nextMixStep)

		}

	}

	if len(warnings) > 0 {
		warning = NewWarningf(strings.Join(warnings, "; "))
		return newComponentList, mixSteps, warning
	}
	return newComponentList, mixSteps, nil
}

// List of the components and corresponding concentrations contained within an LHComponent
type ComponentList struct {
	Components map[string]wunit.Concentration `json:"Components"`
}

// add a single entry to a component list
func (c ComponentList) Add(component *LHComponent, conc wunit.Concentration) (newlist ComponentList) {

	componentName := removeConcUnitFromName(NormaliseName(component.Name()))

	complist := make(map[string]wunit.Concentration)
	for k, v := range c.Components {
		complist[k] = v
	}
	if _, found := complist[componentName]; !found {
		complist[componentName] = conc
	}

	newlist.Components = complist
	return
}

// Get a single concentration set point for a named component present in a component list.
// An error will be returned if the component is not present.
func (c ComponentList) Get(component *LHComponent) (conc wunit.Concentration, err error) {

	componentName := NormaliseName(component.Name())

	conc, found := c.Components[componentName]

	if found {
		return conc, nil
	}

	return conc, &notFound{Name: component.CName, All: c.AllComponents()}
}

// Get a single concentration set point using just the name of a component present in a component list.
// An error will be returned if the component is not present.
func (c ComponentList) GetByName(component string) (conc wunit.Concentration, err error) {

	component = NormaliseName(component)

	conc, found := c.Components[component]

	if found {
		return conc, nil
	}

	return conc, &notFound{Name: component, All: c.AllComponents()}
}

func (c ComponentList) RemoveConcsFromSubComponentNames() (nc ComponentList) {
	newComponentList := make(map[string]wunit.Concentration)
	for compName, conc := range c.Components {
		newCompName := removeConcUnitFromName(compName)
		newComponentList[newCompName] = conc
	}

	nc.Components = newComponentList
	return
}

// List all Components and concentration set points present in a component list.
// if verbose is set to true the field annotations for each component and concentration will be included for each component.
// option1 is verbose, option2 is use mixdelimiter
func (c ComponentList) List(options ...bool) string {
	var verbose bool
	var mixDelimiter bool
	if len(options) > 0 {
		if options[0] {
			verbose = true
		}
	}
	if len(options) > 1 {
		if options[1] {
			mixDelimiter = true
		}
	}
	var s []string

	var sortedKeys []string

	for key := range c.Components {
		sortedKeys = append(sortedKeys, key)
	}

	sort.Strings(sortedKeys)

	for _, k := range sortedKeys {

		v := c.Components[k]

		var message string
		if verbose {
			message = fmt.Sprintln("Component: ", k, "Conc: ", v)
		} else {
			message = v.ToString() + " " + k
		}

		s = append(s, message)
	}
	var list string
	if verbose {
		list = strings.Join(s, ";")
	} else if mixDelimiter {
		list = strings.Join(s, wutil.MIXDELIMITER)
	} else {
		list = strings.Join(s, "---")
	}
	return list
}

// Returns all component names present in component list, sorted in alphabetical order.
func (c ComponentList) AllComponents() []string {
	var s []string

	for k := range c.Components {
		s = append(s, k)
	}

	sort.Strings(s)

	return s
}

type alreadyAdded struct {
	Name string
}

func (a *alreadyAdded) Error() string {
	return "component " + a.Name + " already added"
}

type notFound struct {
	Name string
	All  []string
}

func (a *notFound) Error() string {
	return "component " + a.Name + " not found. Found: " + strings.Join(a.All, ";")
}

// AddSubComponent adds a subcomponent with concentration to a component.
// An error is returned if subcomponent is already found.
func AddSubComponent(component *LHComponent, subcomponent *LHComponent, conc wunit.Concentration) error {

	if component == nil {
		return fmt.Errorf("No component specified so cannot add subcomponent")
	}
	if subcomponent == nil {
		return fmt.Errorf("No subcomponent specified so cannot add subcomponent")
	}
	if len(component.SubComponents.Components) == 0 {
		complist := make(map[string]wunit.Concentration)

		complist[subcomponent.CName] = conc

		var newlist ComponentList

		newlist = newlist.Add(subcomponent, conc)

		if len(newlist.Components) == 0 {

			return fmt.Errorf("No subcomponent added! list still empty")
		}

		if _, err := newlist.Get(subcomponent); err != nil {
			return fmt.Errorf("No subcomponent added, no subcomponent to get: %s!", err.Error())

		}
		component.SubComponents = newlist

		history, err := getHistory(component)

		if err != nil {
			return fmt.Errorf("Error getting History for %s: %s", component.CName, err.Error())
		}

		if len(history.Components) == 0 {
			return fmt.Errorf("No history added!")
		}
		return nil
	} else {

		history, err := getHistory(component)

		if err != nil {
			return err
		}

		if _, found := history.Components[subcomponent.CName]; !found {
			history = history.Add(subcomponent, conc)

			component.SubComponents = history
			return nil
		} else {
			return &alreadyAdded{Name: subcomponent.CName}
		}
	}
}

// AddSubComponents adds a component list to a component.
// If a conflicting sub component concentration is already present then an error will be returned.
// To overwrite all subcomponents ignoring conficts, use OverWriteSubComponents.
func AddSubComponents(component *LHComponent, allsubComponents ComponentList) error {

	for _, compName := range allsubComponents.AllComponents() {
		var comp LHComponent

		comp.CName = compName

		conc, err := allsubComponents.Get(&comp)

		if err != nil {
			return err
		}

		err = AddSubComponent(component, &comp, conc)

		if err != nil {
			return err
		}
	}

	return nil
}

// GetSubComponents returns a component list from a component
func GetSubComponents(component *LHComponent) (componentMap ComponentList, err error) {

	components, err := getHistory(component)

	if err != nil {
		return componentMap, NewWarningf("Error getting componentList for %s: %s", component.Name(), err.Error())
	}

	if len(components.Components) == 0 {
		return components, NewWarningf("No sub components found for %s", component.Name())
	}

	return components, nil
}

// Return a component list from a component.
// Users should use getSubComponents function.
func getHistory(comp *LHComponent) (compList ComponentList, err error) {

	if len(comp.SubComponents.Components) > 0 {
		return comp.SubComponents, nil
	}

	return ComponentList{}, fmt.Errorf("no component list found for %s", comp.Name())
}

// UpdateComponentDetails corrects the sub component list and normalises the name of a component with the details
// of all sample mixes which are specified to be the source of that component.
// This must currently be updated manually using this function.
func UpdateComponentDetails(productOfMixes *LHComponent, mixes ...*LHComponent) error {
	var warnings []string

	subComponents, _, err := SimulateMix(mixes...)

	if err != nil {
		switch err.(type) {
		case Warning:
			warnings = append(warnings, err.Error())
		default:
			return err
		}
	}

	subComponents = subComponents.RemoveConcsFromSubComponentNames()

	err = AddSubComponents(productOfMixes, subComponents)

	if err != nil {
		warnings = append(warnings, err.Error())
	}

	err = NormaliseComponentName(productOfMixes)

	if err != nil {
		warnings = append(warnings, err.Error())
	}

	if len(warnings) > 0 {
		return NewWarningf(strings.Join(warnings, "/n"))
	}

	return nil
}

// EqualLists compares two ComponentLists and returns an error if the lists are not identical.
func EqualLists(list1, list2 ComponentList) error {
	var notEqual []string

	if nonZeroComponents(list1) == 0 && nonZeroComponents(list2) == 0 {
		return nil
	}

	if nonZeroComponents(list1) != nonZeroComponents(list2) {
		return fmt.Errorf("componentlists unequal length: %d, %d", nonZeroComponents(list1), nonZeroComponents(list2))
	}

	for key, value1 := range list1.Components {
		if value2, found := list2.Components[key]; found {
			if fmt.Sprintf("%.2e", value1.SIValue()) != fmt.Sprintf("%.2e", value2.SIValue()) {
				notEqual = append(notEqual, key+" "+fmt.Sprint(value1)+" in list 1 and "+fmt.Sprint(value2)+" in list 2.")
			}
		} else {
			notEqual = append(notEqual, key+" not found in list2. ")
		}
	}
	if len(notEqual) > 0 {
		return fmt.Errorf(strings.Join(notEqual, ". \n"))
	}
	return nil
}

func nonZeroComponents(compList ComponentList) int {
	var nonZero int
	for _, conc := range compList.Components {
		if conc.RawValue() > 0 {
			nonZero++
		}
	}
	return nonZero
}

func equalFold(a, b string) bool {
	return strings.EqualFold(strings.TrimSpace(a), strings.TrimSpace(b))
}

// getComponentConc attempts to retrieve the concentration of subComponentName in component.
// If the component name is equal to subComponentName, the concentration of the component itself is returned.
func getComponentConc(component *LHComponent, subComponentName string) (wunit.Concentration, error) {
	subComponents, _ := GetSubComponents(component) // nolint
	conc, err := subComponents.GetByName(subComponentName)
	if err == nil {
		return conc, nil
	}
	if equalFold(component.Name(), subComponentName) {
		if component.HasConcentration() {
			return component.Concentration(), nil
		}
		return component.Concentration(), fmt.Errorf("subcomponent %s matches component name %s but no concentration found. Error looking up sub component conc: %s", subComponentName, component.Name(), err.Error())
	}
	return wunit.NewConcentration(0.0, "X"), fmt.Errorf("no concentration found for sub component %s in %s. Error looking up sub component conc: %s", subComponentName, component.Name(), err.Error())
}

// hasSubComponent evaluates if a sub component with subComponentName is found in component.
// If the component name is equal to subComponentName, true will be returned.
func hasSubComponent(component *LHComponent, subComponentName string) bool {
	if equalFold(component.Name(), subComponentName) {
		return true
	}
	subComponents, _ := GetSubComponents(component) // nolint
	_, err := subComponents.GetByName(subComponentName)
	return err == nil
}
