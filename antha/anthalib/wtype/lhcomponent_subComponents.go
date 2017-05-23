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
	"reflect"
	"sort"
	"strings"

	. "github.com/antha-lang/antha/antha/anthalib/wunit"
)

// List of the components and corresponding concentrations contained within an LHComponent
type ComponentList struct {
	Components map[string]Concentration `json:"Components"`
}

// add a single entry to a component list
func (c ComponentList) Add(component *LHComponent, conc Concentration) (newlist ComponentList) {
	complist := make(map[string]Concentration)
	for k, v := range c.Components {
		complist[k] = v
	}
	if _, found := complist[component.CName]; !found {
		complist[component.CName] = conc
	}

	newlist.Components = complist
	return
}

// Get a single concentration set point for a named component present in a component list.
// An error will be returned if the component is not present.
func (c ComponentList) Get(component *LHComponent) (conc Concentration, err error) {
	conc, found := c.Components[component.CName]

	if found {
		return conc, nil
	} else {
		return conc, &notFound{Name: component.CName}
	}
	return
}

// Get a single concentration set point using just the name of a component present in a component list.
// An error will be returned if the component is not present.
func (c ComponentList) GetByName(component string) (conc Concentration, err error) {
	conc, found := c.Components[component]

	if found {
		return conc, nil
	} else {
		return conc, &notFound{Name: component}
	}
	return
}

// List all Components and concentration set points presnet in a component list.
// if verbose is set to true the field annotations for each component and concentration will be included for each component.
func (c ComponentList) List(verbose bool) string {
	var ALTMIXDELIMITER = "---"
	var s []string

	for k, v := range c.Components {

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
	} else {
		list = strings.Join(s, ALTMIXDELIMITER)
	}
	return list
}

// Returns all component names present in component list, sorted in alphabetical order.
func (c ComponentList) AllComponents() []string {
	var s []string

	for k, _ := range c.Components {
		s = append(s, k)
	}

	sort.Strings(s)

	return s
}

func (component *LHComponent) ComponentSummary() string {
	subComps, err := component.GetSubComponents()
	var message string
	if err != nil {
		message = err.Error()
	} else {
		message = subComps.List(true)
	}

	conc := "No concentration found"

	if component.HasConcentration() {
		conc = component.Concentration().ToString()
	}

	return fmt.Sprint("Component Name: ", component.CName, "Conc: ", conc, ". SubComponents: ", message)
}

// returns error if component already found
func (component *LHComponent) AddSubComponent(subcomponent *LHComponent, conc Concentration) error {
	var err error

	if component == nil {
		return fmt.Errorf("No component specified so cannot add subcomponent")
	}
	if subcomponent == nil {
		return fmt.Errorf("No subcomponent specified so cannot add subcomponent")
	}
	if _, found := component.Extra[HISTORY]; !found {
		complist := make(map[string]Concentration)

		complist[subcomponent.CName] = conc

		var newlist ComponentList

		newlist = newlist.Add(subcomponent, conc)

		if len(newlist.Components) == 0 {

			return fmt.Errorf("No subcomponent added! list still empty")
		}

		if _, err := newlist.Get(subcomponent); err != nil {
			return fmt.Errorf("No subcomponent added, no subcomponent to get: %s!", err.Error())

		}

		err = component.setHistory(newlist)

		if err != nil {
			return err
		}

		history, err := component.getHistory()

		if err != nil {
			return fmt.Errorf("Error getting History for %s: %s", component.CName, err.Error())
		}

		if len(history.Components) == 0 {
			return fmt.Errorf("No history added!")
		}
		return nil
	} else {

		history, err := component.getHistory()

		if err != nil {
			return err
		}

		if _, found := history.Components[subcomponent.CName]; !found {
			history = history.Add(subcomponent, conc)
			err = component.setHistory(history)
			return err
		} else {
			return &alreadyAdded{Name: subcomponent.CName}
		}
	}
}

// utility function to allow the object properties to be retained when serialised
func serialise(compList ComponentList) ([]byte, error) {

	return json.Marshal(compList)
}

// utility function to allow the object properties to be retained when serialised
func deserialise(data []byte) (compList ComponentList, err error) {
	compList = ComponentList{}
	err = json.Unmarshal(data, &compList)
	return
}

const HISTORY = "History"

// Return a component list from a component.
// Users should use getSubComponents function.
func (comp *LHComponent) getHistory() (compList ComponentList, err error) {

	history, found := comp.Extra[HISTORY]

	if !found {
		return compList, fmt.Errorf("No component list found")
	}

	var bts []byte

	bts, err = json.Marshal(history)
	if err != nil {
		return
	}

	err = json.Unmarshal(bts, &compList)

	if err != nil {
		err = fmt.Errorf("Problem getting %s history. History found: %+v; error: %s", comp.Name(), history, err.Error())
	}

	return
}

// Add a component list to a component.
// Any existing component list will be overwritten.
// Users should use add SubComponents function
func (comp *LHComponent) setHistory(compList ComponentList) error {

	comp.Extra[HISTORY] = compList // serialisedList

	return nil
}

// Add a component list to a component.
// Any existing component list will be overwritten
func (component *LHComponent) AddSubComponents(allsubComponents ComponentList) error {

	for _, compName := range allsubComponents.AllComponents() {
		var comp LHComponent

		comp.CName = compName

		conc, err := allsubComponents.Get(&comp)

		if err != nil {
			return err
		}

		err = component.AddSubComponent(&comp, conc)

		if err != nil {
			return err
		}
	}

	return nil
}

// return a component list from a component
func (component *LHComponent) GetSubComponents() (componentMap ComponentList, err error) {

	components, err := component.getHistory()

	if err != nil {
		return componentMap, fmt.Errorf("Error getting componentList for %s: %s", component.CName, err.Error())
	}

	if len(components.Components) == 0 {
		return components, fmt.Errorf("No sub components found for %s", component.CName)
	}

	return components, nil
}

// Looks for a component matching on name only.
// If more than one component present the first component will be returned with no error
func FindComponent(components []*LHComponent, componentName string) (component *LHComponent, err error) {
	for _, comp := range components {
		if comp.CName == componentName {
			return comp, nil
		}
	}
	return component, fmt.Errorf("No component found with name %s in component list", componentName)
}

// Looks for components matching name, concentration and all sub components (including their concentrations).
// A position of -1 is returned if no match found.
// If the component does not contain a concentration, the name will be matched only
// if multiple matches are found the first will be returned
func ContainsComponent(components []*LHComponent, component *LHComponent, lookForSubComponents bool) (found bool, position int, err error) {

	var errs []string

	// normalise names for more robust evaluation
	var normalisedComponentName string

	if len(strings.Fields(component.CName)) == 2 {
		normalisedComponentName = normalise(component.CName)
	} else {
		normalisedComponentName = component.CName
	}

	for i, comp := range components {

		if !comp.HasConcentration() {
			errs = append(errs, fmt.Sprintf("cannot compare component in list %s without a concentration", comp.CName))
		}

		// normalise names for more robust evaluation
		var normalisedCompName string

		if len(strings.Fields(comp.CName)) == 2 {
			normalisedCompName = normalise(comp.CName)
		} else {
			normalisedCompName = comp.CName
		}

		if normalisedCompName == normalisedComponentName {

			if component.HasConcentration() && comp.HasConcentration() {
				if comp.Concentration().EqualTo(component.Concentration()) {
					if lookForSubComponents {

						compsubcomponents, err := comp.GetSubComponents()
						if err != nil {
							return false, -1, err
						}

						componentSubcomponents, err := component.GetSubComponents()
						if err != nil {
							return false, -1, err
						}

						if reflect.DeepEqual(compsubcomponents, componentSubcomponents) {
							return true, i, nil
						} else {
							errs = append(errs, fmt.Sprintf("Subcomponents lists not equal for %s and %s: Respective lists: %+v and %+v", comp.CName, component.CName, compsubcomponents, componentSubcomponents))

						}
					} else {
						return true, i, nil
					}
				} else {
					errs = append(errs, comp.CName+"concentration "+comp.Concentration().ToString()+" not equal to "+component.CName+" "+component.Concentration().ToString())
				}
			} else {
				if lookForSubComponents {

					compsubcomponents, err := comp.GetSubComponents()
					if err != nil {
						return false, -1, err
					}

					componentSubcomponents, err := component.GetSubComponents()
					if err != nil {
						return false, -1, err
					}

					if reflect.DeepEqual(compsubcomponents, componentSubcomponents) {
						return true, i, nil
					} else {
						errs = append(errs, fmt.Sprintf("Subcomponents lists not equal for %s and %s: Respective lists: %+v and %+v", comp.CName, component.CName, compsubcomponents, componentSubcomponents))

					}
				} else {
					return true, i, nil
				}
			}
		} else {
			errs = append(errs, comp.CName+" name not equal to "+component.CName)
		}
	}

	return false, -1, fmt.Errorf("component %s not found in list: %s. : Errors for each: %s", component.ComponentSummary(), componentNames(components), strings.Join(errs, "\n"))
}

func componentNames(components []*LHComponent) (names []string) {
	for _, component := range components {
		names = append(names, component.CName)
	}
	return
}

// if the component name contains a concentration the concentration name will be normalised
// e.g. 10ng/ul glucose will be normalised to 10 mg/l glucose or 10mM glucose to 10 mM/l glucose or 10mM/l glucose to 10 mM/l glucose or glucose 10mM/l to 10 mM/l glucose
// A concatanenated name such as 10g/L glucose + 10g/L yeast extract will be returned with no modifications
func normalise(name string) (normalised string) {

	if strings.Contains(name, MIXDELIMITER) {
		return name
	}

	containsConc, conc, nameonly := ParseConcentration(name)

	if containsConc {
		return conc.ToString() + " " + nameonly
	} else {
		return name
	}
}

type alreadyAdded struct {
	Name string
}

func (a *alreadyAdded) Error() string {
	return "component " + a.Name + " already added"
}

type notFound struct {
	Name string
}

func (a *notFound) Error() string {
	return "component " + a.Name + " not found"
}
