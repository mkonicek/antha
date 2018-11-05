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

// Package doe facilitates DOE methodology in antha
package doe

import (
	"encoding/json"
	"fmt"
	"math"
	"reflect"
	"sort"
	"strings"

	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/search"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/antha/anthalib/wutil"
)

type DOEPair struct {
	Factor string
	Levels []interface{}
}

func (pair DOEPair) LevelCount() (numberoflevels int) {
	numberoflevels = len(pair.Levels)
	return
}

func round(value interface{}) (interface{}, error) {
	var v interface{}
	if levelfloat, found := value.(float64); found {
		levelfloat, err := wutil.Roundto(levelfloat, 6)
		if err != nil {
			return v, err
		}
		v = levelfloat
	}
	return v, nil
}

// RoundLevels round floats and concs
func (pair DOEPair) RoundLevels() (newpair DOEPair, err error) {

	var newlevels []interface{}

	for i, level := range pair.Levels {
		var v interface{}

		conc, err := HandleConcFactor(pair.Factor, level)

		if err == nil {
			v = conc.RawValue()
			newlevels = append(newlevels, v)
		} else if levelfloat, found := level.(float64); found {
			levelfloat, err := wutil.Roundto(levelfloat, 6)
			if err != nil {
				return newpair, err
			}
			v = levelfloat
			newlevels = append(newlevels, v)
		} else if merged, err := ToMergedLevel(level); err == nil {

			var factors []string
			var values []interface{}

			for key, value := range merged.OriginalFactorPairs {

				factors = append(factors, key)

				value, err = round(value)
				if err != nil {
					return newpair, err
				}
				values = append(values, value)
			}

			newmerged, err := MakeMergedLevel(factors, values)

			if err != nil {
				return newpair, err
			}

			v = newmerged
			newlevels = append(newlevels, v)

		} else {
			return newpair, fmt.Errorf("cannot round non-float or conc value for level %d, %s: %s", i, level, err.Error())
		}
	}

	newpair.Factor = pair.Factor
	newpair.Levels = newlevels

	return newpair, nil
}

// MaxLevel returns the maximum level of a DOEPair if the factor is numeric, if
// the factor is not numeric an error is returned
func (pair DOEPair) MaxLevel() (maxlevel interface{}, err error) {
	arraytype, err := search.CheckArrayType(pair.Levels)
	if err != nil {
		return maxlevel, err
	}
	if arraytype == "float64" {
		var floats []float64
		for _, level := range pair.Levels {
			levelfloat := level.(float64)
			floats = append(floats, levelfloat)
		}
		sort.Float64s(floats)

		return floats[len(floats)-1], nil

	}
	return maxlevel, fmt.Errorf("cannot sort non-numeric type of %s", arraytype)
}

// MinLevel returns the minimum level of a DOEPair if the factor is numeric, if
// the factor is not numeric an error is returned unless the factor is a
// MergedLevel comprising of concentration values
func (pair DOEPair) MinLevel() (minlevel interface{}, err error) {

	arraytype, err := search.CheckArrayType(pair.Levels)
	if err != nil {
		return minlevel, err
	}
	if arraytype == "float64" {
		var floats []float64
		for _, level := range pair.Levels {
			levelfloat := level.(float64)
			floats = append(floats, levelfloat)
		}
		sort.Float64s(floats)

		return floats[0], nil

		// added test for serialised MergedLevel which are not caught be search.CheckArrayType.
		// assumes all levels are same type and so evaluates first entry, if this is not the case it will be caught in the loop
	} else if _, err := ToMergedLevel(pair.Levels[0]); arraytype == "MergedLevel" || err == nil {

		// use first level to get keys
		// assumes all MergedLevels contain the same merged factors
		// if all keys are not included the function will error
		// if len of keys is greater than the first entry the function will return an error

		merged, err := ToMergedLevel(pair.Levels[0])

		if err != nil {
			return minlevel, err
		}

		masterKeys := merged.Sort()

		var keyToLowest = make(map[string]interface{})

		for _, key := range masterKeys {

			var lowestconc wunit.Concentration

			for i, level := range pair.Levels {
				if i == 0 {
					merged, err := ToMergedLevel(level)

					if err != nil {
						return minlevel, err
					}
					concInt, found := merged.OriginalFactorPairs[key]

					if !found {
						return minlevel, fmt.Errorf("merged factor %s not found in level %d of %s", key, i, fmt.Sprint(pair))
					}

					conc, err := HandleConcFactor(key, concInt)

					if err != nil {
						return minlevel, err
					}

					lowestconc = conc

				} else {
					merged, err := ToMergedLevel(level)

					if err != nil {
						return minlevel, err
					}

					concInt, found := merged.OriginalFactorPairs[key]

					if !found {
						return minlevel, fmt.Errorf("merged factor %s not found in level %d of %s", key, i, fmt.Sprint(pair))
					}

					conc, err := HandleConcFactor(key, concInt)

					if err != nil {
						return minlevel, err
					}
					if conc.LessThan(lowestconc) {
						lowestconc = conc

					}
				}

			}

			v := lowestconc.RawValue()

			keyToLowest[key] = v
		}

		var factors []string
		var values []interface{}

		for k, v := range keyToLowest {

			factors = append(factors, k)
			values = append(values, v)
		}

		lowestMerged, err := MakeMergedLevel(factors, values)

		if err != nil {
			return minlevel, err
		}

		for _, level := range pair.Levels {
			merged, err := ToMergedLevel(level)

			if err != nil {
				return minlevel, err
			}

			if equal, err := merged.EqualToMergeConcs(lowestMerged); equal && err == nil {
				return level, nil
			}
		}
		return minlevel, fmt.Errorf("cannot find lowest level of MergedLevel: lowest found %s in %s", fmt.Sprintln(lowestMerged), pairSummary(pair))

	} else if arraytype == "string" {
		var lowest int
		var lowestconc wunit.Concentration
		for i, level := range pair.Levels {

			conc, err := HandleConcFactor(pair.Factor, level)

			if err != nil {
				return minlevel, fmt.Errorf("cannot sort: non-numeric type of %s found and not possible to convert level %d level %s into concentration", arraytype, i, level)

			}
			if i == 0 {
				lowest = i
				lowestconc = conc
			} else if conc.LessThanRounded(lowestconc, 9) {
				lowest = i
				lowestconc = conc
			}

		}
		return pair.Levels[lowest], nil
	}
	return minlevel, fmt.Errorf("cannot sort non-numeric type of %s", arraytype)
}

func pairSummary(pair DOEPair) string {
	var s []string
	s = append(s, fmt.Sprintln(pair.Factor))
	for _, level := range pair.Levels {
		s = append(s, fmt.Sprintln(level))
	}
	return (strings.Join(s, "; "))
}

// MakePair pairs a factor with all levels found in runs
func MakePair(runs []Run, factor string) (allfactorPairs DOEPair, err error) {

	var levels []interface{}

	for _, run := range runs {
		level, err := run.GetFactorValue(factor)

		if err != nil {
			return allfactorPairs, err
		}
		levels = append(levels, level)
	}
	return Pair(factor, levels), nil
}

// RemoveDuplicateLevels creates a copy of the pair with reduced factors if an
// error occurs the orginal pair is returned
func RemoveDuplicateLevels(pair DOEPair) (reducedpair DOEPair, err error) {
	reducedpair.Factor = pair.Factor
	// remove duplicates
	newlevels, err := search.RemoveDuplicateValues(pair.Levels)
	if err != nil {
		return pair, err
	}
	reducedpair.Levels = newlevels

	return
}

func Pair(factordescription string, levels []interface{}) (doepair DOEPair) {
	doepair.Factor = factordescription
	doepair.Levels = levels
	return
}

type Run struct {
	RunNumber            int
	StdNumber            int
	Factordescriptors    []string
	Setpoints            []interface{}
	Responsedescriptors  []string
	ResponseValues       []interface{}
	AdditionalHeaders    []string // could represent a class e.g. Environment variable, processed, raw, location
	AdditionalSubheaders []string // e.g. well ID, Ambient Temp, order,
	AdditionalValues     []interface{}
}

func Copy(run Run) (newrun Run) {

	newrun.RunNumber = run.RunNumber
	newrun.StdNumber = run.StdNumber

	factordescriptors := make([]string, 0)
	factordescriptors = append(factordescriptors, run.Factordescriptors...)

	setpoints := make([]interface{}, 0)
	setpoints = append(setpoints, run.Setpoints...)

	responsedescriptors := make([]string, 0)
	responsedescriptors = append(responsedescriptors, run.Responsedescriptors...)

	responsevalues := make([]interface{}, 0)
	responsevalues = append(responsevalues, run.ResponseValues...)

	newrun.Factordescriptors = factordescriptors
	newrun.Setpoints = setpoints
	newrun.Responsedescriptors = responsedescriptors
	newrun.ResponseValues = responsevalues

	additionalheaders := make([]string, 0)
	additionalheaders = append(additionalheaders, run.AdditionalHeaders...)

	newrun.AdditionalHeaders = additionalheaders

	additionalsubheaders := make([]string, 0)
	additionalsubheaders = append(additionalsubheaders, run.AdditionalSubheaders...)

	newrun.AdditionalSubheaders = additionalsubheaders

	values := make([]interface{}, 0)
	values = append(values, run.AdditionalValues...)

	newrun.AdditionalValues = values

	return
}

// EqualTo compares Runs but ignores run order, std order  and additional
// information but checks all other properties for identical matches
func (run Run) EqualTo(run2 Run) (bool, error) {

	for i, value := range run.Factordescriptors {
		if value != run2.Factordescriptors[i] {
			return false, fmt.Errorf("factors differ between runs at factor %d: %s and %s ", i, value, run2.Factordescriptors[i])
		}
	}

	for i, value := range run.Setpoints {

		if !reflect.DeepEqual(value, run2.Setpoints[i]) {
			return false, fmt.Errorf("setpoints differ between runs at setpoint %d: %s and %s ", i, value, run2.Setpoints[i])
		}
	}

	for i, value := range run.Responsedescriptors {
		if value != run2.Responsedescriptors[i] {
			return false, fmt.Errorf("responses differ between runs at response %d: %s and %s ", i, value, run2.Responsedescriptors[i])
		}
	}

	for i, value := range run.ResponseValues {
		if !reflect.DeepEqual(value, run2.ResponseValues[i]) {
			return false, fmt.Errorf("response values differ between runs at response value %d: %s and %s ", i, value, run2.ResponseValues[i])
		}
	}
	return true, nil
}
func (run Run) AddResponseValue(responsedescriptor string, responsevalue interface{}) {

	for i, descriptor := range run.Responsedescriptors {
		if search.EqualFold(descriptor, responsedescriptor) {
			run.ResponseValues[i] = responsevalue
		}
	}

}

func (run Run) AllResponses() (headers []string, values []interface{}) {
	headers = make([]string, 0)
	values = make([]interface{}, 0)

	headers = append(headers, run.Responsedescriptors...)
	values = append(values, run.ResponseValues...)
	return
}

func (run Run) AllFactors() (headers []string, values []interface{}) {
	headers = make([]string, 0)
	values = make([]interface{}, 0)

	headers = append(headers, run.Factordescriptors...)
	values = append(values, run.Setpoints...)
	return
}

func (run Run) GetResponseValue(responsedescriptor string) (responsevalue interface{}, err error) {
	var errs []string
	var tempresponsevalue interface{}
	headers, _ := run.AllResponses()

	errstr := fmt.Sprint("response descriptor", responsedescriptor, "not found in ", headers)
	errs = append(errs, errstr)
	for i, descriptor := range run.Responsedescriptors {
		if search.EqualFold(descriptor, responsedescriptor) {
			responsevalue = run.ResponseValues[i]
			return responsevalue, nil
		} else if search.ContainsEqualFold(descriptor, responsedescriptor) {

			errstr := fmt.Sprint("response descriptor", responsedescriptor, "found within ", descriptor, "but no exact match")
			errs = append(errs, errstr)
			tempresponsevalue = run.ResponseValues[i]
		} else if search.ContainsEqualFold(responsedescriptor, descriptor) {

			errstr := fmt.Sprint("response descriptors of ", descriptor, "found within ", responsedescriptor, "but not exact match")
			errs = append(errs, errstr)
			tempresponsevalue = run.ResponseValues[i]
		}
	}
	err = fmt.Errorf(strings.Join(errs, "\n"))
	responsevalue = tempresponsevalue
	return
}

// GetFactorValue searches for a factor descriptor in the list of factors in a Run.
// Exact matches will be searched first, followed by matches once any unit is removed (e.g. TotalVolume (ul) to TotalVolume)
func (run Run) GetFactorValue(factordescriptor string) (factorvalue interface{}, err error) {

	var tempresponsevalue interface{}
	headers, values := run.AllFactors()

	var errs []string

	errstr := fmt.Sprint("factor descriptor ", factordescriptor, " not found in ", headers, values)
	errs = append(errs, errstr)
	for i, descriptor := range run.Factordescriptors {
		if search.EqualFold(descriptor, factordescriptor) {
			factorvalue = run.Setpoints[i]
			return factorvalue, nil
		} else if factor, _ := splitFactorFromUnit(factordescriptor); search.EqualFold(descriptor, factor) {
			factorvalue = run.Setpoints[i]
			return factorvalue, nil
		} else if search.ContainsEqualFold(descriptor, factordescriptor) {
			errstr := fmt.Sprint("factor descriptor ", factordescriptor, "found within ", descriptor, " but no exact match")
			errs = append(errs, errstr)
			tempresponsevalue = run.Setpoints[i]
		} else if search.ContainsEqualFold(factordescriptor, descriptor) {
			errstr := fmt.Sprint("factor descriptors of ", descriptor, "found within ", factordescriptor, "but not exact match")
			errs = append(errs, errstr)
			tempresponsevalue = run.Setpoints[i]
		}
	}
	err = fmt.Errorf(strings.Join(errs, "\n"))
	factorvalue = tempresponsevalue
	return
}

func AddNewResponseFieldandValue(run Run, responsedescriptor string, responsevalue interface{}) (newrun Run) {

	// check if float64 is NaN or Inf
	if float, found := responsevalue.(float64); found {
		if math.IsInf(float, 0) {
			responsevalue = "Inf"
		} else if math.IsNaN(float) {
			responsevalue = "Nan"
		}
	}

	newrun = run

	responsedescriptors := run.Responsedescriptors
	responsevalues := run.ResponseValues

	responsedescriptors = append(responsedescriptors, responsedescriptor)
	responsevalues = append(responsevalues, responsevalue)

	newrun.Responsedescriptors = responsedescriptors
	newrun.ResponseValues = responsevalues

	return
}

// AddNewResponseField adds a new response field. The response field will be
// added with a blank response value.  The field will only be added if the
// response is not already present.
func AddNewResponseField(run Run, responsedescriptor string) (newrun Run) {

	newrun = run

	responsedescriptors := make([]string, len(run.Responsedescriptors))
	responsevalues := make([]interface{}, len(run.ResponseValues))

	var skip bool

	for i, descriptor := range run.Responsedescriptors {
		if search.EqualFold(descriptor, responsedescriptor) {
			skip = true
		}
		responsedescriptors[i] = run.Responsedescriptors[i]
	}

	for i := range run.ResponseValues {
		responsevalues[i] = run.ResponseValues[i]
	}

	if !skip {
		responsedescriptors = append(responsedescriptors, responsedescriptor)
		var nilvalue interface{}
		responsevalues = append(run.ResponseValues, nilvalue)
	}
	newrun.Responsedescriptors = responsedescriptors
	newrun.ResponseValues = responsevalues

	return
}

func DeleteResponseField(run Run, responsedescriptor string) (newrun Run) {

	newrun = run

	responsedescriptors := make([]string, 0)
	responsevalues := make([]interface{}, 0)

	for i, descriptor := range run.Responsedescriptors {
		if !search.EqualFold(descriptor, responsedescriptor) {
			responsedescriptors = append(responsedescriptors, descriptor)
			responsevalues = append(responsevalues, run.ResponseValues[i])
		}
	}

	newrun.Responsedescriptors = responsedescriptors
	newrun.ResponseValues = responsevalues

	return
}

func ReplaceResponseValue(run Run, responsedescriptor string, responsevalue interface{}) (newrun Run) {

	newrun = run

	responsedescriptors := make([]string, 0)
	responsevalues := make([]interface{}, 0)

	for i, descriptor := range run.Responsedescriptors {
		if !search.EqualFold(descriptor, responsedescriptor) {
			responsedescriptors = append(responsedescriptors, descriptor)
			responsevalues = append(responsevalues, run.ResponseValues[i])
		} else {
			responsedescriptors = append(responsedescriptors, descriptor)
			responsevalues = append(responsevalues, responsevalue)
		}
	}

	newrun.Responsedescriptors = responsedescriptors
	newrun.ResponseValues = responsevalues

	return
}

func DeleteAllResponses(run Run) (newrun Run) {
	newrun = run

	responsedescriptors := make([]string, 0)
	responsevalues := make([]interface{}, 0)

	newrun.Responsedescriptors = responsedescriptors
	newrun.ResponseValues = responsevalues

	return
}

func AddNewFactorFieldandValue(run Run, factordescriptor string, factorvalue interface{}) (newrun Run) {
	newrun = run

	factordescriptors := run.Factordescriptors
	factorvalues := run.Setpoints

	factordescriptors = append(factordescriptors, factordescriptor)
	factorvalues = append(factorvalues, factorvalue)

	newrun.Factordescriptors = factordescriptors
	newrun.Setpoints = factorvalues

	return
}

func DeleteFactorField(run Run, factorDescriptor string) (newrun Run) {
	newrun = run

	factorDescriptors := make([]string, 0)
	factorValues := make([]interface{}, 0)

	for i, descriptor := range run.Factordescriptors {
		if !search.EqualFold(descriptor, factorDescriptor) {
			factorDescriptors = append(factorDescriptors, descriptor)
			factorValues = append(factorValues, run.Setpoints[i])
		}
	}

	newrun.Factordescriptors = factorDescriptors
	newrun.Setpoints = factorValues

	return
}

// ReplaceFactorField will replace the factorToReplace with factorToReplaceWith with Set point setPointToReplaceWith.
// The index of the factor to replace will be the same as the replaced factor.
func ReplaceFactorField(run Run, factorToReplace, factorToReplaceWith string, setPointToReplaceWith interface{}) (newrun Run) {
	newrun = run

	factorDescriptors := make([]string, 0)
	factorValues := make([]interface{}, 0)

	for i, descriptor := range run.Factordescriptors {
		if !search.EqualFold(descriptor, factorToReplace) {
			factorDescriptors = append(factorDescriptors, descriptor)
			factorValues = append(factorValues, run.Setpoints[i])
		} else {
			factorDescriptors = append(factorDescriptors, factorToReplaceWith)
			factorValues = append(factorValues, setPointToReplaceWith)
		}
	}

	newrun.Factordescriptors = factorDescriptors
	newrun.Setpoints = factorValues

	return
}

func AddAdditionalValue(run Run, additionalsubheader string, additionalvalue interface{}) (newrun Run) {

	newrun = run

	values := make([]interface{}, 0)
	values = append(values, run.AdditionalValues...)

	for _, descriptor := range run.AdditionalSubheaders {
		if search.EqualFold(descriptor, additionalsubheader) {
			values = append(values, additionalvalue)
		}
	}

	newrun.AdditionalValues = values

	return
}

func ReplaceAdditionalValue(run Run, additionalsubheader string, additionalvalue interface{}) (newrun Run) {

	newrun = run

	values := make([]interface{}, len(run.AdditionalSubheaders))

	for i, descriptor := range run.AdditionalSubheaders {
		if search.EqualFold(descriptor, additionalsubheader) {
			values[i] = additionalvalue
		} else {
			values[i] = run.AdditionalValues[i]
		}
	}

	newrun.AdditionalValues = values

	return
}

func AddAdditionalHeaders(run Run, additionalheader string, additionalsubheader string) (newrun Run) {
	newrun = run

	headers := make([]string, 0)
	headers = append(headers, run.AdditionalHeaders...)
	headers = append(headers, additionalheader)

	subheaders := make([]string, 0)
	subheaders = append(subheaders, run.AdditionalSubheaders...)
	subheaders = append(subheaders, additionalsubheader)

	newrun.AdditionalHeaders = headers
	newrun.AdditionalSubheaders = subheaders

	return
}

func AddAdditionalHeaderandValue(run Run, additionalheader string, additionalsubheader string, additionalvalue interface{}) (newrun Run) {

	// only add column if no column with header exists
	if !search.InStrings(run.AdditionalSubheaders, additionalsubheader) {
		midrun := AddAdditionalHeaders(run, additionalheader, additionalsubheader)
		// fmt.Println("midrun: ", midrun)
		newrun = AddAdditionalValue(midrun, additionalsubheader, additionalvalue)
	} else {
		newrun = ReplaceAdditionalValue(run, additionalsubheader, additionalvalue)
	}
	return
}

func (run Run) CheckAdditionalInfo(subheader string, value interface{}) bool {

	for i, header := range run.AdditionalSubheaders {
		if search.EqualFold(header, subheader) && run.AdditionalValues[i] == value {
			return true
		}
	}
	return false
}

func (run Run) GetAdditionalInfo(subheader string) (value interface{}, err error) {

	for i, header := range run.AdditionalSubheaders {
		if search.EqualFold(header, subheader) {
			value = run.AdditionalValues[i]
			return value, err

		}
	}
	return value, fmt.Errorf("header %s not found in %v", subheader, run.AdditionalSubheaders)
}

func AddFixedFactors(runs []Run, fixedfactors []DOEPair) (runswithfixedfactors []Run) {

	if len(runs) > 0 {
		for _, run := range runs {
			descriptors := make([]string, len(run.Factordescriptors))
			copy(descriptors, run.Factordescriptors)

			setpoints := make([]interface{}, len(run.Setpoints))
			copy(setpoints, run.Setpoints)

			for _, fixed := range fixedfactors {
				descriptors = append(descriptors, fixed.Factor)
				setpoints = append(setpoints, fixed.Levels[0])

			}
			run.Factordescriptors = descriptors
			run.Setpoints = setpoints

		}

	} else {
		runs = RunsFromFixedFactors(fixedfactors)
	}

	runswithfixedfactors = runs

	return
}

func RunsFromFixedFactors(fixedfactors []DOEPair) (runswithfixedfactors []Run) {

	var run Run
	var descriptors = make([]string, 0)
	var setpoints = make([]interface{}, 0)

	for _, factor := range fixedfactors {

		descriptors = append(descriptors, factor.Factor)
		setpoints = append(setpoints, factor.Levels[0])

	}

	run.Factordescriptors = descriptors
	run.Setpoints = setpoints

	runswithfixedfactors = make([]Run, 1)

	runswithfixedfactors[0] = run

	return

}

func AllComboCount(pairs []DOEPair) (numberofuniquecombos int) {
	movingcount := (pairs[0]).LevelCount()

	for i := 1; i < len(pairs); i++ {
		movingcount = movingcount * (pairs[i]).LevelCount()
	}

	numberofuniquecombos = movingcount
	return
}

func FixedAndNonFixed(factors []DOEPair) (fixedfactors []DOEPair, nonfixed []DOEPair) {

	fixedfactors = make([]DOEPair, 0)
	nonfixed = make([]DOEPair, 0)

	for _, factor := range factors {
		if len(factor.Levels) == 1 {
			fixedfactors = append(fixedfactors, factor)
		} else if len(factor.Levels) > 1 {
			nonfixed = append(nonfixed, factor)
		}
	}
	return
}

func IsFixedFactor(factor DOEPair) (yesorno bool) {
	if len(factor.Levels) == 1 {
		yesorno = true
	}
	return
}

func sameFactorLevels(run1 Run, run2 Run, factor string) (same bool, err error) {
	value1, err := run1.GetFactorValue(factor)
	if err != nil {
		return
	}
	value2, err := run2.GetFactorValue(factor)
	if err != nil {
		return
	}

	if reflect.DeepEqual(value1, value2) {
		return true, nil
	}

	return false, fmt.Errorf("Different values found for factor %s: %s and %s", factor, value1, value2)
}

// MergeOption is an option to control the position in the factors of a run at which the merged factor will be added.
// Valid options are MoveToFront, MoveToBack and UsePositionOfLastFactorInFactorList.
// The default if no option is set is UsePositionOfLastFactorInFactorList.
type MergeOption string

var (
	// MoveToFront specifies that the merged factor should be moved to the front of the factor list
	MoveToFront MergeOption = "MoveToFront"

	// MoveToBack specifies that the merged factor should be moved to the back of the factor list
	MoveToBack MergeOption = "MoveToBack"

	// UsePositionOfLastFactorInFactorList specifies that the merged factor should be moved to the position of the last factor in the merged list.
	UsePositionOfLastFactorInFactorList MergeOption = "UsePositionOfLastFactorInFactorList"

	// Default will use UsePositionOfLastFactorInFactorList.
	DefaultMergeOption MergeOption = UsePositionOfLastFactorInFactorList
)

// intended to be used in conjunction with AllCombinations to merge levels of a series of runs
// run1 will be used as the master run whose properties will be preferentially inherited in run3
// MergeOptions can be specified to control the position in the factors of a run at which the merged factor will be added.
// Valid options are MoveToFront, MoveToBack and UsePositionOfLastFactorInFactorList.
// The default if no option is set is UsePositionOfLastFactorInFactorList.
func mergeFactorLevels(run1 Run, run2 Run, factors []string, newLevel interface{}, options ...MergeOption) (run3 Run, newfactorName string, err error) {

	var usePosition int

	const (
		front int = iota
		back
		positionOfLastFactor
	)

	if len(options) > 1 {
		err = fmt.Errorf("only one merge option can be specified at a time. Valid Options %s, %s and default %s. Found %v", MoveToFront, MoveToBack, UsePositionOfLastFactorInFactorList, options)
		return
	}
	if len(options) == 0 {
		usePosition = positionOfLastFactor
	} else {
		if options[0] == MoveToBack {
			usePosition = back
		} else if options[0] == MoveToFront {
			usePosition = front
		} else if options[0] == DefaultMergeOption {
			usePosition = positionOfLastFactor
		} else {
			err = fmt.Errorf("invalid merge option specified. Valid Options %s, %s and default %s. Found %v", MoveToFront, MoveToBack, UsePositionOfLastFactorInFactorList, options)
			return
		}
	}

	var factornames []string
	var deadrun Run
	run3 = Copy(run1)
	for i, factor := range factors {
		if same, err := sameFactorLevels(run1, run2, factor); !same || err != nil {
			return deadrun, "", err
		}
		// preserve order of factor names from original design
		factornames = append(factornames, factor)

		// if last in list replace factor rather than delete
		if i == len(factors)-1 && usePosition == positionOfLastFactor {
			newfactorName = MergeFactorNames(factornames)
			// make new set of runs
			run3 = ReplaceFactorField(run3, factor, newfactorName, newLevel)
		} else {
			run3 = DeleteFactorField(run3, factor)
		}

	}

	if usePosition == back {
		newfactorName = MergeFactorNames(factornames)
		run3 = AddNewFactorFieldandValue(run3, newfactorName, newLevel)
	} else if usePosition == front {
		newfactorName = MergeFactorNames(factornames)
		run3.Factordescriptors = append([]string{newfactorName}, run3.Factordescriptors...)
		run3.Setpoints = append([]interface{}{newLevel}, run3.Setpoints...)
	}

	return
}

func mergeStrings(factorNames []string) (combinedFactor string) {
	combinedFactor = strings.Join(factorNames, wutil.MIXDELIMITER)
	return
}

func MergeFactorNames(factorNames []string) (combinedFactor string) {
	combinedFactor = mergeStrings(factorNames)
	return
}

// A MergedLevel is the product of merging two factor levels and retaining the
// original merged factor level pairs in a map
type MergedLevel struct {
	OriginalFactorPairs map[string]interface{} `json:"OriginalFactorPairs"` // map of original factor names to levels e.g. "Glucose":"10uM", "Glycerol mM": 0
}

// UnMerge returns keys in alphabetical order; values are returned in the order
// corresponding to the key order
func (m MergedLevel) UnMerge() (factors []string, setpoints []interface{}) {

	sortedKeys := m.Sort()

	for _, key := range sortedKeys {

		setpoint := m.OriginalFactorPairs[key]
		factors = append(factors, key)
		setpoints = append(setpoints, setpoint)
	}

	return
}

// Sort returns keys in alphabetical order; values are returned in the order
// corresponding to the key order
func (m MergedLevel) Sort() (orderedKeys []string) {
	for k := range m.OriginalFactorPairs {
		orderedKeys = append(orderedKeys, k)
	}
	sort.Strings(orderedKeys)

	return
}

// EqualTo evaluates whether two merged levels are equal
func (m MergedLevel) EqualTo(e MergedLevel) (bool, error) {

	if len(m.OriginalFactorPairs) != len(e.OriginalFactorPairs) {
		return false, fmt.Errorf("Merged levels are not of equal size")
	}

	for k, v := range m.OriginalFactorPairs {
		if !reflect.DeepEqual(e.OriginalFactorPairs[k], v) {
			return false, fmt.Errorf("original factor %s not found in merged level being evaluated", k)
		}
	}
	return true, nil
}

// EqualToMergeConcs evaluates whether two merged levels are equal
func (m MergedLevel) EqualToMergeConcs(e MergedLevel) (bool, error) {

	if len(m.OriginalFactorPairs) != len(e.OriginalFactorPairs) {
		return false, fmt.Errorf("Merged levels are not of equal size")
	}

	for k, v := range m.OriginalFactorPairs {

		vconc, err := HandleConcFactor(k, v)

		if err != nil {
			return false, err
		}

		econc, err := HandleConcFactor(k, e.OriginalFactorPairs[k])

		if err != nil {
			return false, err
		}

		if !concsEqual(vconc, econc) {
			return false, fmt.Errorf("original factor %s level %s not equal to %s in merged level being evaluated", k, vconc.ToString(), econc.ToString())
		}
	}
	return true, nil
}

func concsEqual(conc1, conc2 wunit.Concentration) bool {

	roundedconc1, _ := wutil.Roundto(conc1.RawValue(), 6)

	roundedconc2, _ := wutil.Roundto(conc2.RawValue(), 6)

	dif := math.Abs(roundedconc1 - roundedconc2)

	epsilon := math.Nextafter(1, 2) - 1

	if dif < (epsilon*1000) && conc1.Unit().PrefixedSymbol() == conc2.Unit().PrefixedSymbol() {
		return true
	}
	return false
}

// MakeMergedLevel merges an array of factors into a single merged level type
// using arrays of factor names and levels in order
func MakeMergedLevel(factorsInOrder []string, levelsInOrder []interface{}) (m MergedLevel, err error) {
	pairs := make(map[string]interface{})

	if len(factorsInOrder) != len(levelsInOrder) {
		return m, fmt.Errorf("unequal length of factors and levels so cannot make mergedLevel")
	}
	for i := range factorsInOrder {
		pairs[factorsInOrder[i]] = levelsInOrder[i]
	}

	m.OriginalFactorPairs = pairs
	return
}

// ToMergedLevel casts a valid interface into a MergedLevel. Independent of
// serialisation history.  If the value cannot be cast into a merged level an
// error is returned.
func ToMergedLevel(level interface{}) (m MergedLevel, err error) {

	bts, err := json.Marshal(level)
	if err != nil {
		return
	}
	err = json.Unmarshal(bts, &m)

	return
}

func findMatchingLevels(run Run, allcombos []Run, factors []string, options ...MergeOption) (matchedrun Run, err error) {
	var levels []interface{}
	for _, factor := range factors {

		level, err := run.GetFactorValue(factor)
		if err == nil {
			levels = append(levels, level)
		}

	}
	m, mergErr := MakeMergedLevel(factors, levels)
	for _, combo := range allcombos {

		if newrun, _, err := mergeFactorLevels(run, combo, factors, m, options...); err == nil && mergErr == nil {
			return newrun, nil
		}
	}
	return matchedrun, fmt.Errorf("No matching combination of levels %s found for run %v in %v: last errors: %v and %v", factors, run, allcombos, mergErr, err)
}

// MergeRunsFromAllCombos will make a set of runs based on the original runs which merges all factors specified in factors and finds a valid merged run from all combos to inject in as the set point.
// MergeOptions can be specified to control the position in the factors of a run at which the merged factor will be added.
// Valid options are MoveToFront, MoveToBack and UsePositionOfLastFactorInFactorList.
// The default if no option is set is UsePositionOfLastFactorInFactorList.
func MergeRunsFromAllCombos(originalRuns []Run, allcombos []Run, factors []string, options ...MergeOption) (mergedRuns []Run, err error) {

	mergedRuns = make([]Run, len(originalRuns))

	for i, original := range originalRuns {
		newRun, err := findMatchingLevels(original, allcombos, factors, options...)
		if err == nil {
			mergedRuns[i] = newRun
		} else {
			return mergedRuns, err
		}
	}
	return
}

// UnMergeRuns turns mergedlevels into original factors and levels
func UnMergeRuns(mergedRuns []Run) (originalRuns []Run) {

	originalRuns = make([]Run, len(mergedRuns))

	for i, run := range mergedRuns {
		newrun := Copy(run)
		for j, setpoint := range run.Setpoints {
			if merged, ok := setpoint.(MergedLevel); ok {
				factors, setpoints := merged.UnMerge()
				for k := range factors {
					newrun = AddNewFactorFieldandValue(newrun, factors[k], setpoints[k])
				}
				newrun = DeleteFactorField(newrun, run.Factordescriptors[j])
			}
		}
		originalRuns[i] = newrun
	}
	return
}

// LowestLevelValue returns the lowest level value for a specified factor name
// from an array of runs
func LowestLevelValue(runs []Run, factor string) (value interface{}, err error) {

	// get all levels for the new factor
	pair, err := MakePair(runs, factor)
	if err != nil {
		return value, err
	}
	// remove duplicates
	reducedPair, err := RemoveDuplicateLevels(pair)
	if err != nil {
		return value, err
	}

	roundedpair, err := reducedPair.RoundLevels()

	if err != nil {
		return value, err
	}

	return roundedpair.MinLevel()

}

// AllCombinations combines all factor pairs to make all possible combinations
// of runs
func AllCombinations(factors []DOEPair) (runs []Run) {

	//fixed, nonfixed := FixedAndNonFixed(factors)

	numberofruns := AllComboCount(factors)

	runs = make([]Run, numberofruns)

	var swapevery int
	var numberofswaps int
	for i, factor := range factors {

		counter := 0
		runswitheachlevelforthisfactor := numberofruns / factor.LevelCount()

		if i == 0 {
			swapevery = runswitheachlevelforthisfactor
			numberofswaps = runswitheachlevelforthisfactor / swapevery
		} else {
			swapevery = swapevery / factor.LevelCount()
			numberofswaps = runswitheachlevelforthisfactor / swapevery
		}

		for j := 0; j < numberofswaps; j++ {
			for _, level := range factor.Levels {
				for k := 0; k < swapevery; k++ {

					runs[counter] = AddNewFactorFieldandValue(runs[counter], factor.Factor, level)
					runs[counter].RunNumber = counter + 1
					runs[counter].StdNumber = counter + 1
					counter++
				}
			}
		}

	}
	return
}
