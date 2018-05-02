// /anthalib/driver/liquidhandling/makelhpolicy.go: Part of the Antha language
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

package liquidtype

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"

	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/AnthaPath"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/doe"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
)

// deprecate this
type policyFile struct {
	Filename                string
	DXORJMP                 string
	FactorColumns           *[]int
	LiquidTypeStarterNumber int
}

func (polfile policyFile) Prepend() (prepend string) {
	nameparts := strings.Split(polfile.Filename, ".")
	prepend = nameparts[0]
	return
}

func (polfile policyFile) StarterNumber() (starternumber int) {
	starternumber = polfile.LiquidTypeStarterNumber
	return
}

// deprecate this
func makePolicyFile(filename string, dxorjmp string, factorcolumns *[]int, liquidtypestartnumber int) (policyfile policyFile) {
	policyfile.Filename = filename
	policyfile.DXORJMP = dxorjmp
	policyfile.FactorColumns = factorcolumns
	policyfile.LiquidTypeStarterNumber = liquidtypestartnumber
	return
}

// deprecate this
// policy files to put in ./antha
var availablePolicyfiles []policyFile = []policyFile{
	makePolicyFile("170516CCFDesign_noTouchoff_noBlowout.xlsx", "DX", nil, 100),
	makePolicyFile("2700516AssemblyCCF.xlsx", "DX", nil, 1000),
	makePolicyFile("newdesign2factorsonly.xlsx", "JMP", &[]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13}, 2000),
	makePolicyFile("190516OnePolicy.xlsx", "JMP", &[]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13}, 3000),
	makePolicyFile("AssemblycategoricScreen.xlsx", "JMP", &[]int{1, 2, 3, 4, 5}, 4000),
	makePolicyFile("090816dispenseerrordiagnosis.xlsx", "JMP", &[]int{2}, 5000),
	makePolicyFile("090816combineddesign.xlsx", "JMP", &[]int{1}, 6000),
}

// BASEPolicy is the policy to use as a starting point to produce custom LHPolicies
var BASEPolicy = "default" //"dna"

// deprecate this
func policyFilefromName(filename string) (pol policyFile, found bool) {
	for _, policy := range availablePolicyfiles {
		if policy.Filename == filename {
			pol = policy
			found = true
			return
		}
	}
	return
}

// deprecate this
func PolicyMakerfromFilename(filename string) (policies []wtype.LHPolicy, names []string, runs []doe.Run, err error) {

	doeliquidhandlingFile, found := policyFilefromName(filename)
	if !found {
		err = fmt.Errorf("policyfilename " + filename + " not found")
		return
	}
	filenameparts := strings.Split(doeliquidhandlingFile.Filename, ".")

	policies, names, runs, err = PolicyMakerfromDesign(BASEPolicy, doeliquidhandlingFile.DXORJMP, doeliquidhandlingFile.Filename, filenameparts[0])
	return
}

// deprecate this
func PolicyMakerfromDesign(basepolicy string, DXORJMP string, dxdesignfilename string, prepend string) (policies []wtype.LHPolicy, names []string, runs []doe.Run, err error) {

	policyitemmap := wtype.MakePolicyItems()
	intfactors := make([]string, 0)

	for key, val := range policyitemmap {

		if val.Type.Name() == "int" {
			intfactors = append(intfactors, key)

		}

	}
	if DXORJMP == "DX" {
		contents, err := ioutil.ReadFile(filepath.Join(anthapath.Path(), dxdesignfilename))

		if err != nil {
			return policies, names, runs, err
		}

		runs, err = doe.RunsFromDXDesignContents(contents, intfactors)

		if err != nil {
			return policies, names, runs, err
		}

	} else if DXORJMP == "JMP" {

		factorcolumns := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13}
		responsecolumns := []int{14, 15, 16, 17}

		contents, err := ioutil.ReadFile(filepath.Join(anthapath.Path(), dxdesignfilename))

		if err != nil {
			return policies, names, runs, err
		}

		runs, err = doe.RunsFromJMPDesignContents(contents, factorcolumns, responsecolumns, intfactors)
		if err != nil {
			return policies, names, runs, err
		}
	} else {
		return policies, names, runs, fmt.Errorf("only JMP or DX allowed as valid inputs for DXORJMP variable")
	}
	policies, names, err = PolicyMakerfromRuns(basepolicy, runs, prepend, false)
	return policies, names, runs, err
}

func PolicyMaker(basepolicy string, factors []doe.DOEPair, nameprepend string, concatfactorlevelsinname bool) (policies []wtype.LHPolicy, names []string, err error) {

	runs := doe.AllCombinations(factors)

	policies, names, err = PolicyMakerfromRuns(basepolicy, runs, nameprepend, concatfactorlevelsinname)

	return
}

// PolicyMakerFromBytes creates a policy map from a design file in the format of a JMP design file.
// Any valid parameter name and corresponding parameter type from aparam.go are valid entries.
func PolicyMakerFromBytes(data []byte, basePolicy wtype.PolicyName) (policyMap map[string]wtype.LHPolicy, err error) {

	var Warnings []string

	DXORJMP := "JMP"
	policyitemmap := wtype.MakePolicyItems()
	var intfactors []string

	for key, val := range policyitemmap {

		if val.Type.Name() == "int" {
			intfactors = append(intfactors, key)
		}
	}

	runs, err := doe.RunsFromDesignPreResponsesContents(data, intfactors, DXORJMP)
	if err != nil {
		err = fmt.Errorf("error converting DOE design into runs: %s", err.Error())
		return
	}

	policies, names, err := PolicyMakerfromRuns(string(basePolicy), runs, "custom", false)

	if err != nil {
		switch err.(type) {
		case wtype.Warning:
			Warnings = append(Warnings, err.Error())
		default:
			return
		}
	}

	policyMap = make(map[string]wtype.LHPolicy)
	for i, policy := range policies {
		if policy.Name() == "" {
			err = policy.SetName(names[i])
			if err != nil {
				return
			}
		}
		if _, ok := policyMap[policy.Name()]; !ok {
			policyMap[policy.Name()] = policy
		} else {
			err = fmt.Errorf("duplicate policy name (%s) found in cusom policy file", policy.Name())
			return
		}
	}

	if len(Warnings) > 0 {
		return policyMap, wtype.NewWarningf(strings.Join(Warnings, "\n"))
	}
	return policyMap, nil
}

// PolicyMakerfromRuns creates a policy map from a set of doe Runs in the format of a JMP design file.
// Any valid parameter name and corresponding parameter type from aparam.go are valid entries.
func PolicyMakerfromRuns(basepolicy string, runs []doe.Run, nameprepend string, concatfactorlevelsinname bool) (policies []wtype.LHPolicy, names []string, err error) {

	policyitemmap := wtype.MakePolicyItems()

	policies = make([]wtype.LHPolicy, 0)
	var runWarnings []string

	basePolicy, err := wtype.GetPolicyByType(wtype.LiquidType(basepolicy))

	if err != nil {
		return
	}

	copyPolicy := func(policy wtype.LHPolicy) wtype.LHPolicy {
		var newPolicy = make(map[string]interface{})
		for key, value := range policy {
			newPolicy[key] = value
		}
		return newPolicy
	}

	for i, run := range runs {
		policy := make(wtype.LHPolicy)
		policy = copyPolicy(basePolicy)
		var warnings []string
		for j, desc := range run.Factordescriptors {
			policyCommand, ok := policyitemmap[desc]
			if ok {
				if reflect.TypeOf(run.Setpoints[j]) != policyCommand.Type {
					err = fmt.Errorf("invalid value (%s) of type (%T) for LHPolicy command (%s) in run %d expecting value of type %s", run.Setpoints[j], run.Setpoints[j], desc, i+1, policyCommand.Type.Name())
					return policies, names, err
				}
				policy[desc] = run.Setpoints[j]
			} else if i == 0 {
				warnings = append(warnings, "Invalid PolicyCommand specified in design file: "+desc)
			}
		}

		var name string
		if concatfactorlevelsinname {
			name = nameprepend
			for key, value := range policy {
				name = fmt.Sprint(name, "_", key, ":", value)
			}
		} else {
			name = nameprepend + strconv.Itoa(run.RunNumber)
		}
		names = append(names, name)
		err = policy.SetName(name)
		if err != nil {
			return policies, names, err
		}
		policies = append(policies, policy)
		if len(warnings) > 0 {
			runWarnings = append(runWarnings, fmt.Sprint("Errors :\n", strings.Join(warnings, "\n \t - ")))
		}
	}

	if len(runWarnings) > 0 {
		return policies, names, wtype.NewWarning(strings.Join(runWarnings, "\n\n"))
	}

	return policies, names, nil
}
