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

package liquidhandling

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/AnthaPath"
	. "github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/doe"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/ghodss/yaml"
)

type PolicyFile struct {
	Filename                string
	DXORJMP                 string
	FactorColumns           *[]int
	LiquidTypeStarterNumber int
}

func (polfile PolicyFile) Prepend() (prepend string) {
	nameparts := strings.Split(polfile.Filename, ".")
	prepend = nameparts[0]
	return
}

func (polfile PolicyFile) StarterNumber() (starternumber int) {
	starternumber = polfile.LiquidTypeStarterNumber
	return
}

func MakePolicyFile(filename string, dxorjmp string, factorcolumns *[]int, liquidtypestartnumber int) (policyfile PolicyFile) {
	policyfile.Filename = filename
	policyfile.DXORJMP = dxorjmp
	policyfile.FactorColumns = factorcolumns
	policyfile.LiquidTypeStarterNumber = liquidtypestartnumber
	return
}

// policy files to put in ./antha
var AvailablePolicyfiles []PolicyFile = []PolicyFile{
	MakePolicyFile("170516CCFDesign_noTouchoff_noBlowout.xlsx", "DX", nil, 100),
	MakePolicyFile("2700516AssemblyCCF.xlsx", "DX", nil, 1000),
	MakePolicyFile("newdesign2factorsonly.xlsx", "JMP", &[]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13}, 2000),
	MakePolicyFile("190516OnePolicy.xlsx", "JMP", &[]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13}, 3000),
	MakePolicyFile("AssemblycategoricScreen.xlsx", "JMP", &[]int{1, 2, 3, 4, 5}, 4000),
	MakePolicyFile("090816dispenseerrordiagnosis.xlsx", "JMP", &[]int{2}, 5000),
	MakePolicyFile("090816combineddesign.xlsx", "JMP", &[]int{1}, 6000),
}

// change to range through several files
//var DOEliquidhandlingFile = "170516CCFDesign_noTouchoff_noBlowout.xlsx" // "2700516AssemblyCCF.xlsx" //"newdesign2factorsonly.xlsx" // "170516CCFDesign_noTouchoff_noBlowout.xlsx" // "170516CFF.xlsx" //"newdesign2factorsonly.xlsx" "170516CCFDesign_noTouchoff_noBlowout.xlsx" // //"newdesign2factorsonly.xlsx" //"8run4cpFactorial.xlsx" //"FullFactorial.xlsx" // "ScreenLHPolicyDOE2.xlsx"
//var DXORJMP = "DX"                                                      //"JMP"
var BASEPolicy = "default" //"dna"

func MakePolicies() map[string]LHPolicy {
	pols := make(map[string]LHPolicy)

	// what policies do we need?
	pols["water"] = MakeWaterPolicy()
	pols["culture"] = MakeCulturePolicy()
	pols["culturereuse"] = MakeCultureReusePolicy()
	pols["glycerol"] = MakeGlycerolPolicy()
	pols["solvent"] = MakeSolventPolicy()
	pols["default"] = MakeDefaultPolicy()
	pols["foamy"] = MakeFoamyPolicy()
	pols["dna"] = MakeDNAPolicy()
	pols["DoNotMix"] = MakeDefaultPolicy()
	pols["NeedToMix"] = MakeNeedToMixPolicy()
	pols["PreMix"] = PreMixPolicy()
	pols["PostMix"] = PostMixPolicy()
	pols["MegaMix"] = MegaMixPolicy()
	pols["viscous"] = MakeViscousPolicy()
	pols["Paint"] = MakePaintPolicy()

	//      pols["lysate"] = MakeLysatePolicy()
	pols["protein"] = MakeProteinPolicy()
	pols["detergent"] = MakeDetergentPolicy()
	pols["load"] = MakeLoadPolicy()
	pols["loadlow"] = MakeLoadPolicy()
	pols["loadwater"] = MakeLoadWaterPolicy()
	pols["DispenseAboveLiquid"] = MakeDispenseAboveLiquidPolicy()
	pols["PEG"] = MakePEGPolicy()
	pols["Protoplasts"] = MakeProtoplastPolicy()
	pols["dna_mix"] = MakeDNAMixPolicy()
	pols["dna_cells_mix"] = MakeDNACELLSMixPolicy()
	pols["plateout"] = MakePlateOutPolicy()
	pols["colony"] = MakeColonyPolicy()
	pols["colonymix"] = MakeColonyMixPolicy()
	//      pols["lysate"] = MakeLysatePolicy()
	pols["carbon_source"] = MakeCarbonSourcePolicy()
	pols["nitrogen_source"] = MakeNitrogenSourcePolicy()

	/*policies, names := PolicyMaker(Allpairs, "DOE_run", false)
	for i, policy := range policies {
		pols[names[i]] = policy
	}
	*/

	// TODO: Remove this hack
	for _, DOEliquidhandlingFile := range AvailablePolicyfiles {
		if _, err := os.Stat(filepath.Join(anthapath.Path(), DOEliquidhandlingFile.Filename)); err == nil {
			//if antha.Anthafileexists(DOEliquidhandlingFile) {
			//fmt.Println("found lhpolicy doe file", DOEliquidhandlingFile)

			filenameparts := strings.Split(DOEliquidhandlingFile.Filename, ".")

			policies, names, _, err := PolicyMakerfromDesign(BASEPolicy, DOEliquidhandlingFile.DXORJMP, DOEliquidhandlingFile.Filename, filenameparts[0])
			//policies, names, _, err := PolicyMakerfromDesign(BASEPolicy, DXORJMP, DOEliquidhandlingFile, "DOE_run")
			for i, policy := range policies {
				pols[names[i]] = policy
			}
			if err != nil {
				panic(err)
			}
		} else {
			//	fmt.Println("no lhpolicy doe file found named: ", DOEliquidhandlingFile)
		}
	}
	return pols

}

func PolicyFilefromName(filename string) (pol PolicyFile, found bool) {
	for _, policy := range AvailablePolicyfiles {
		if policy.Filename == filename {
			pol = policy
			found = true
			return
		}
	}
	return
}

func PolicyMakerfromFilename(filename string) (policies []LHPolicy, names []string, runs []Run, err error) {

	doeliquidhandlingFile, found := PolicyFilefromName(filename)
	if found == false {
		panic("policyfilename" + filename + "not found")
	}
	filenameparts := strings.Split(doeliquidhandlingFile.Filename, ".")

	policies, names, runs, err = PolicyMakerfromDesign(BASEPolicy, doeliquidhandlingFile.DXORJMP, doeliquidhandlingFile.Filename, filenameparts[0])
	return
}

func PolicyMakerfromDesign(basepolicy string, DXORJMP string, dxdesignfilename string, prepend string) (policies []LHPolicy, names []string, runs []Run, err error) {

	policyitemmap := MakePolicyItems()
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

		runs, err = RunsFromDXDesignContents(contents, intfactors)

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

		runs, err = RunsFromJMPDesignContents(contents, factorcolumns, responsecolumns, intfactors)
		if err != nil {
			return policies, names, runs, err
		}
	} else {
		return policies, names, runs, fmt.Errorf("only JMP or DX allowed as valid inputs for DXORJMP variable")
	}
	policies, names = PolicyMakerfromRuns(basepolicy, runs, prepend, false)
	return policies, names, runs, err
}

func PolicyMaker(basepolicy string, factors []DOEPair, nameprepend string, concatfactorlevelsinname bool) (policies []LHPolicy, names []string) {

	runs := AllCombinations(factors)

	policies, names = PolicyMakerfromRuns(basepolicy, runs, nameprepend, concatfactorlevelsinname)

	return
}

func PolicyMakerfromRuns(basepolicy string, runs []Run, nameprepend string, concatfactorlevelsinname bool) (policies []LHPolicy, names []string) {

	policyitemmap := MakePolicyItems()

	names = make([]string, 0)
	policies = make([]LHPolicy, 0)

	policy := MakeDefaultPolicy()
	policy["CAN_MULTI"] = false

	/*base, _ := GetPolicyByName(basepolicy)

	for key, value := range base {
		policy[key] = value
	}
	*/
	//fmt.Println("basepolicy:", basepolicy)
	for _, run := range runs {
		for j, desc := range run.Factordescriptors {

			_, ok := policyitemmap[desc]
			if ok {

				/*if val.Type.Name() == "int" {
					aInt, found := run.Setpoints[j].(int)

					var bInt int

					bInt = int(aInt)
					if found {
						run.Setpoints[j] = interface{}(bInt)
					}
				}*/
				policy[desc] = run.Setpoints[j]
			} /* else {
				panic("policyitem " + desc + " " + "not present! " + "These are present: " + policyitemmap.TypeList())
			}*/
		}

		// raising runtime error when using concat == true
		if concatfactorlevelsinname {
			name := nameprepend
			for key, value := range policy {
				name = fmt.Sprint(name, "_", key, ":", value)

			}

		} else {
			names = append(names, nameprepend+strconv.Itoa(run.RunNumber))
		}
		policies = append(policies, policy)

		//policy := GetPolicyByName(basepolicy)
		policy = MakeDefaultPolicy()
	}

	return
}

//func MakeLysatePolicy() LHPolicy {
//        lysatepolicy := make(LHPolicy, 6)
//        lysatepolicy["ASPSPEED"] = 1.0
//        lysatepolicy["DSPSPEED"] = 1.0
//        lysatepolicy["ASP_WAIT"] = 2.0
//        lysatepolicy["ASP_WAIT"] = 2.0
//        lysatepolicy["DSP_WAIT"] = 2.0
//        lysatepolicy["PRE_MIX"] = 5
//        lysatepolicy["CAN_MSA"]= false
//        return lysatepolicy
//}
//func MakeProteinPolicy() LHPolicy {
//        proteinpolicy := make(LHPolicy, 4)
//        proteinpolicy["DSPREFERENCE"] = 2
//        proteinpolicy["CAN_MULTI"] = true
//        proteinpolicy["PRE_MIX"] = 3
//        proteinpolicy["CAN_MSA"] = false
//        return proteinpolicy
//}

func GetPolicyByName(policyname string) (lhpolicy LHPolicy, policypresent bool) {
	policymap := MakePolicies()

	lhpolicy, policypresent = policymap[policyname]
	return
}

func AvailablePolicies() (policies []string) {

	policies = make([]string, 0)
	policymap := MakePolicies()

	for key, _ := range policymap {
		policies = append(policies, key)
	}
	return
}

/*
Available policy field names and policy types to use:

Here is a list of everything currently implemented in the liquid handling policy framework

ASPENTRYSPEED,                    ,float64,      ,allows slow moves into liquids
ASPSPEED,                                ,float64,     ,aspirate pipetting rate
ASPZOFFSET,                           ,float64,      ,mm above well bottom when aspirating
ASP_WAIT,                                   ,float64,     ,wait time in seconds post aspirate
BLOWOUTOFFSET,                    ,float64,     ,mm above BLOWOUTREFERENCE
BLOWOUTREFERENCE,          ,int,             ,where to be when blowing out: 0 well bottom, 1 well top
BLOWOUTVOLUME,                ,float64,      ,how much to blow out
CAN_MULTI,                              ,bool,         ,is multichannel operation allowed?
DSPENTRYSPEED,                    ,float64,     ,allows slow moves into liquids
DSPREFERENCE,                      ,int,            ,where to be when dispensing: 0 well bottom, 1 well top
DSPSPEED,                              ,float64,       ,dispense pipetting rate
DSPZOFFSET,                         ,float64,          ,mm above DSPREFERENCE
DSP_WAIT,                               ,float64,        ,wait time in seconds post dispense
EXTRA_ASP_VOLUME,            ,wunit.Volume,       ,additional volume to take up when aspirating
EXTRA_DISP_VOLUME,           ,wunit.Volume,       ,additional volume to dispense
JUSTBLOWOUT,                      ,bool,            ,shortcut to get single transfer
POST_MIX,                               ,int,               ,number of mix cycles to do after dispense
POST_MIX_RATE,                    ,float64,          ,pipetting rate when post mixing
POST_MIX_VOL,                      ,float64,          ,volume to post mix (ul)
POST_MIX_X,                          ,float64,           ,x offset from centre of well (mm) when post-mixing
POST_MIX_Y,                          ,float64,           ,y offset from centre of well (mm) when post-mixing
POST_MIX_Z,                          ,float64,           ,z offset from centre of well (mm) when post-mixing
PRE_MIX,                                ,int,               ,number of mix cycles to do before aspirating
PRE_MIX_RATE,                     ,float64,           ,pipetting rate when pre mixing
PRE_MIX_VOL,                       ,float64,           ,volume to pre mix (ul)
PRE_MIX_X,                              ,float64,          ,x offset from centre of well (mm) when pre-mixing
PRE_MIX_Y,                              ,float64,           ,y offset from centre of well (mm) when pre-mixing
PRE_MIX_Z,                              ,float64,           ,z offset from centre of well (mm) when pre-mixing
TIP_REUSE_LIMIT,                    ,int,                ,number of times tips can be reused for asp/dsp cycles
TOUCHOFF,                              ,bool,             ,whether to move to TOUCHOFFSET after dispense
TOUCHOFFSET,                         ,float64,          ,mm above wb to touch off at


*/

func MakePEGPolicy() LHPolicy {
	policy := make(LHPolicy, 9)
	policy["ASPSPEED"] = 1.5
	policy["DSPSPEED"] = 1.5
	policy["ASP_WAIT"] = 2.0
	policy["DSP_WAIT"] = 2.0
	policy["ASPZOFFSET"] = 1.0
	policy["DSPZOFFSET"] = 1.0
	policy["POST_MIX"] = 3
	policy["POST_MIX_Z"] = 1.0
	policy["BLOWOUTVOLUME"] = 50.0
	policy["POST_MIX_VOLUME"] = 190.0
	policy["BLOWOUTVOLUMEUNIT"] = "ul"
	policy["TOUCHOFF"] = false
	policy["CAN_MULTI"] = false
	policy["RESET_OVERRIDE"] = true
	return policy
}

func MakeProtoplastPolicy() LHPolicy {
	policy := make(LHPolicy, 7)
	policy["ASPSPEED"] = 0.5
	policy["DSPSPEED"] = 0.5
	policy["ASPZOFFSET"] = 1.0
	policy["DSPZOFFSET"] = 1.0
	policy["BLOWOUTVOLUME"] = 100.0
	policy["BLOWOUTVOLUMEUNIT"] = "ul"
	policy["TOUCHOFF"] = false
	policy["TIP_REUSE_LIMIT"] = 5
	policy["CAN_MULTI"] = false
	return policy
}

func MakePaintPolicy() LHPolicy {

	policy := make(LHPolicy, 13)
	policy["DSPREFERENCE"] = 0
	policy["DSPZOFFSET"] = 0.5
	policy["ASPSPEED"] = 1.5
	policy["DSPSPEED"] = 1.5
	policy["ASP_WAIT"] = 1.0
	policy["DSP_WAIT"] = 1.0
	//policy["PRE_MIX"] = 3
	policy["POST_MIX"] = 3
	policy["BLOWOUTVOLUME"] = 0.0
	policy["BLOWOUTVOLUMEUNIT"] = "ul"
	policy["TOUCHOFF"] = true
	policy["CAN_MULTI"] = false

	return policy
}

func MakeDispenseAboveLiquidPolicy() LHPolicy {

	policy := make(LHPolicy, 7)
	policy["DSPREFERENCE"] = 1 // 1 indicates dispense at top of well
	policy["ASPSPEED"] = 3.0
	policy["DSPSPEED"] = 3.0
	//policy["ASP_WAIT"] = 1.0
	//policy["DSP_WAIT"] = 1.0
	policy["BLOWOUTVOLUME"] = 50.0
	policy["BLOWOUTVOLUMEUNIT"] = "ul"
	policy["TOUCHOFF"] = false
	policy["CAN_MULTI"] = false

	return policy
}

func MakeColonyPolicy() LHPolicy {

	policy := make(LHPolicy, 10)
	policy["DSPREFERENCE"] = 0
	policy["DSPZOFFSET"] = 0.0
	policy["ASPSPEED"] = 3.0
	policy["DSPSPEED"] = 3.0
	policy["ASP_WAIT"] = 1.0
	policy["POST_MIX"] = 1
	policy["BLOWOUTVOLUME"] = 0.0
	policy["BLOWOUTVOLUMEUNIT"] = "ul"
	policy["TOUCHOFF"] = false
	policy["CAN_MULTI"] = false
	policy["RESET_OVERRIDE"] = true

	return policy
}

func MakeColonyMixPolicy() LHPolicy {
	policy := MakeColonyPolicy()
	policy["POST_MIX"] = 3
	return policy
}

func MakeWaterPolicy() LHPolicy {

	waterpolicy := make(LHPolicy, 6)
	waterpolicy["DSPREFERENCE"] = 0
	//waterpolicy["CAN_MULTI"] = true
	waterpolicy["CAN_MULTI"] = false
	waterpolicy["CAN_MSA"] = true
	waterpolicy["CAN_SDD"] = true
	waterpolicy["DSPZOFFSET"] = 1.0
	waterpolicy["BLOWOUTVOLUME"] = 50.0
	return waterpolicy
}
func MakeFoamyPolicy() LHPolicy {
	foamypolicy := make(LHPolicy, 5)
	foamypolicy["DSPREFERENCE"] = 0
	foamypolicy["TOUCHOFF"] = true
	foamypolicy["CAN_MULTI"] = false
	foamypolicy["CAN_MSA"] = false
	foamypolicy["CAN_SDD"] = true
	return foamypolicy
}
func MakeCulturePolicy() LHPolicy {
	culturepolicy := make(LHPolicy, 10)
	culturepolicy["PRE_MIX"] = 2
	culturepolicy["ASPSPEED"] = 2.0
	culturepolicy["DSPSPEED"] = 2.0
	culturepolicy["CAN_MULTI"] = false
	culturepolicy["CAN_MSA"] = false
	culturepolicy["CAN_SDD"] = false
	culturepolicy["DSPREFERENCE"] = 0
	culturepolicy["DSPZOFFSET"] = 0.5
	culturepolicy["TIP_REUSE_LIMIT"] = 0
	culturepolicy["NO_AIR_DISPENSE"] = true
	culturepolicy["BLOWOUTVOLUME"] = 0.0
	culturepolicy["BLOWOUTVOLUMEUNIT"] = "ul"
	culturepolicy["TOUCHOFF"] = false

	return culturepolicy
}

func MakePlateOutPolicy() LHPolicy {
	culturepolicy := make(LHPolicy, 17)
	//culturepolicy["PRE_MIX"] = 1
	//culturepolicy["PRE_MIX_VOLUME"] = 19.0
	//culturepolicy["PRE_MIX_Z"] = 1.0
	//culturepolicy["PRE_MIX_RATE"] = 4.0
	culturepolicy["ASPSPEED"] = 2.0
	culturepolicy["ASPZOFFSET"] = 1.0
	culturepolicy["DSPSPEED"] = 2.0
	culturepolicy["CAN_MULTI"] = false
	culturepolicy["ASP_WAIT"] = 1.0
	culturepolicy["DSP_WAIT"] = 1.0
	//culturepolicy["CAN_MSA"] = false
	//culturepolicy["CAN_SDD"] = false
	//culturepolicy["DSPREFERENCE"] = 0
	culturepolicy["DSPZOFFSET"] = 0.0
	culturepolicy["TIP_REUSE_LIMIT"] = 7
	culturepolicy["NO_AIR_DISPENSE"] = true
	//culturepolicy["BLOWOUTVOLUME"] = 0.0
	//culturepolicy["BLOWOUTVOLUMEUNIT"] = "ul"
	culturepolicy["TOUCHOFF"] = false

	return culturepolicy
}

func MakeCultureReusePolicy() LHPolicy {
	culturepolicy := make(LHPolicy, 10)
	culturepolicy["PRE_MIX"] = 2
	culturepolicy["ASPSPEED"] = 2.0
	culturepolicy["DSPSPEED"] = 2.0
	//culturepolicy["CAN_MULTI"] = true
	culturepolicy["CAN_MULTI"] = false
	culturepolicy["CAN_MSA"] = true
	culturepolicy["CAN_SDD"] = true
	culturepolicy["DSPREFERENCE"] = 0
	culturepolicy["DSPZOFFSET"] = 0.5
	culturepolicy["NO_AIR_DISPENSE"] = true
	culturepolicy["BLOWOUTVOLUME"] = 0.0
	culturepolicy["BLOWOUTVOLUMEUNIT"] = "ul"
	culturepolicy["TOUCHOFF"] = false

	return culturepolicy
}

func MakeGlycerolPolicy() LHPolicy {
	glycerolpolicy := make(LHPolicy, 9)
	glycerolpolicy["ASPSPEED"] = 1.5
	glycerolpolicy["DSPSPEED"] = 1.5
	glycerolpolicy["ASP_WAIT"] = 1.0
	glycerolpolicy["DSP_WAIT"] = 1.0
	glycerolpolicy["TIP_REUSE_LIMIT"] = 0
	glycerolpolicy["CAN_MULTI"] = false
	glycerolpolicy["POST_MIX"] = 3
	glycerolpolicy["POST_MIX_VOLUME"] = 20.0
	glycerolpolicy["POST_MIX_RATE"] = 3.74
	return glycerolpolicy
}

func MakeViscousPolicy() LHPolicy {
	glycerolpolicy := make(LHPolicy, 4)
	glycerolpolicy["ASPSPEED"] = 1.5
	glycerolpolicy["DSPSPEED"] = 1.5
	glycerolpolicy["ASP_WAIT"] = 1.0
	glycerolpolicy["DSP_WAIT"] = 1.0
	//glycerolpolicy["TIP_REUSE_LIMIT"] = 0
	return glycerolpolicy
}
func MakeSolventPolicy() LHPolicy {
	solventpolicy := make(LHPolicy, 5)
	solventpolicy["PRE_MIX"] = 3
	solventpolicy["DSPREFERENCE"] = 0
	solventpolicy["DSPZOFFSET"] = 0.5
	solventpolicy["NO_AIR_DISPENSE"] = true
	solventpolicy["CAN_MULTI"] = false
	return solventpolicy
}

func MakeDNAPolicy() LHPolicy {
	dnapolicy := make(LHPolicy, 10)
	dnapolicy["ASPSPEED"] = 2.0
	dnapolicy["DSPSPEED"] = 2.0
	dnapolicy["CAN_MULTI"] = false
	dnapolicy["CAN_MSA"] = false
	dnapolicy["CAN_SDD"] = false
	dnapolicy["DSPREFERENCE"] = 0
	dnapolicy["DSPZOFFSET"] = 0.5
	dnapolicy["TIP_REUSE_LIMIT"] = 0
	dnapolicy["NO_AIR_DISPENSE"] = true
	dnapolicy["RESET_OVERRIDE"] = true
	dnapolicy["TOUCHOFF"] = false
	return dnapolicy
}

func MakeDNAMixPolicy() LHPolicy {
	dnapolicy := MakeDNAPolicy()
	dnapolicy["POST_MIX_VOLUME"] = 10.0
	dnapolicy["POST_MIX"] = 5
	dnapolicy["POST_MIX_Z"] = 0.5
	dnapolicy["POST_MIX_RATE"] = 3.0
	return dnapolicy
}

func MakeDNACELLSMixPolicy() LHPolicy {
	dnapolicy := MakeDNAPolicy()
	dnapolicy["POST_MIX_VOLUME"] = 20.0
	dnapolicy["POST_MIX"] = 2
	dnapolicy["POST_MIX_Z"] = 0.5
	dnapolicy["POST_MIX_RATE"] = 1.0
	return dnapolicy
}

func MakeDetergentPolicy() LHPolicy {
	detergentpolicy := make(LHPolicy, 9)
	//        detergentpolicy["POST_MIX"] = 3
	detergentpolicy["ASPSPEED"] = 1.0
	detergentpolicy["DSPSPEED"] = 1.0
	detergentpolicy["CAN_MULTI"] = false
	detergentpolicy["CAN_MSA"] = false
	detergentpolicy["CAN_SDD"] = false
	detergentpolicy["DSPREFERENCE"] = 0
	detergentpolicy["DSPZOFFSET"] = 0.5
	detergentpolicy["TIP_REUSE_LIMIT"] = 8
	detergentpolicy["NO_AIR_DISPENSE"] = true
	return detergentpolicy
}
func MakeProteinPolicy() LHPolicy {
	proteinpolicy := make(LHPolicy, 10)
	proteinpolicy["POST_MIX"] = 5
	proteinpolicy["POST_MIX_VOLUME"] = 50.0
	proteinpolicy["ASPSPEED"] = 2.0
	proteinpolicy["DSPSPEED"] = 2.0
	proteinpolicy["CAN_MULTI"] = false
	proteinpolicy["CAN_MSA"] = false
	proteinpolicy["CAN_SDD"] = false
	proteinpolicy["DSPREFERENCE"] = 0
	proteinpolicy["DSPZOFFSET"] = 0.5
	proteinpolicy["TIP_REUSE_LIMIT"] = 0
	proteinpolicy["NO_AIR_DISPENSE"] = true
	return proteinpolicy
}
func MakeLoadPolicy() LHPolicy {

	loadpolicy := make(LHPolicy)
	loadpolicy["ASPSPEED"] = 1.0
	loadpolicy["DSPSPEED"] = 0.1
	loadpolicy["CAN_MULTI"] = false
	loadpolicy["CAN_MSA"] = false
	loadpolicy["CAN_SDD"] = false
	loadpolicy["TOUCHOFF"] = false
	loadpolicy["TIP_REUSE_LIMIT"] = 0
	loadpolicy["NO_AIR_DISPENSE"] = true
	loadpolicy["TOUCHOFF"] = false
	loadpolicy["BLOWOUTREFERENCE"] = 1
	loadpolicy["BLOWOUTOFFSET"] = 0.0
	loadpolicy["BLOWOUTVOLUME"] = 0.0
	loadpolicy["BLOWOUTVOLUMEUNIT"] = "ul"
	return loadpolicy
}

func MakeLoadWaterPolicy() LHPolicy {

	loadpolicy := make(LHPolicy)
	loadpolicy["ASPSPEED"] = 1.0
	loadpolicy["DSPSPEED"] = 0.1
	loadpolicy["CAN_MULTI"] = false
	loadpolicy["CAN_MSA"] = false
	//loadpolicy["CAN_SDD"] = false
	loadpolicy["TOUCHOFF"] = false
	loadpolicy["NO_AIR_DISPENSE"] = true
	loadpolicy["TOUCHOFF"] = false
	loadpolicy["TIP_REUSE_LIMIT"] = 100
	loadpolicy["BLOWOUTREFERENCE"] = 1
	loadpolicy["BLOWOUTOFFSET"] = 0.0
	loadpolicy["BLOWOUTVOLUME"] = 0.0
	loadpolicy["BLOWOUTVOLUMEUNIT"] = "ul"
	return loadpolicy
}

func MakeLoadlowPolicy() LHPolicy {

	loadpolicy := make(LHPolicy)
	loadpolicy["ASPSPEED"] = 1.0
	loadpolicy["DSPSPEED"] = 1.0
	loadpolicy["CAN_MULTI"] = false
	loadpolicy["CAN_MSA"] = false
	loadpolicy["CAN_SDD"] = false
	loadpolicy["TOUCHOFF"] = false
	loadpolicy["TIP_REUSE_LIMIT"] = 0
	loadpolicy["DSPZOFFSET"] = 0.5
	loadpolicy["NO_AIR_DISPENSE"] = true
	loadpolicy["TOUCHOFF"] = false
	loadpolicy["BLOWOUTREFERENCE"] = 1
	loadpolicy["BLOWOUTOFFSET"] = 0.0
	loadpolicy["BLOWOUTVOLUME"] = 0.0
	loadpolicy["BLOWOUTVOLUMEUNIT"] = "ul"
	return loadpolicy
}

func MakeNeedToMixPolicy() LHPolicy {
	dnapolicy := make(LHPolicy, 15)
	dnapolicy["POST_MIX"] = 3
	dnapolicy["POST_MIX_VOLUME"] = 19.0
	dnapolicy["POST_MIX_RATE"] = 3.74
	dnapolicy["PRE_MIX"] = 3
	dnapolicy["PRE_MIX_VOLUME"] = 20.0
	dnapolicy["PRE_MIX_RATE"] = 3.74
	dnapolicy["ASPSPEED"] = 3.74
	dnapolicy["DSPSPEED"] = 3.74
	dnapolicy["CAN_MULTI"] = false
	dnapolicy["CAN_MSA"] = false
	dnapolicy["CAN_SDD"] = false
	dnapolicy["DSPREFERENCE"] = 0
	dnapolicy["DSPZOFFSET"] = 0.5
	dnapolicy["TIP_REUSE_LIMIT"] = 0
	dnapolicy["NO_AIR_DISPENSE"] = true
	return dnapolicy
}

func PreMixPolicy() LHPolicy {
	dnapolicy := make(LHPolicy, 12)
	//dnapolicy["POST_MIX"] = 3
	//dnapolicy[""POST_MIX_VOLUME"] = 10.0
	//dnapolicy["POST_MIX_RATE"] = 3.74
	dnapolicy["PRE_MIX"] = 3
	dnapolicy["PRE_MIX_VOLUME"] = 19.0
	dnapolicy["PRE_MIX_RATE"] = 3.74
	dnapolicy["ASPSPEED"] = 3.74
	dnapolicy["DSPSPEED"] = 3.74
	dnapolicy["CAN_MULTI"] = false
	dnapolicy["CAN_MSA"] = false
	dnapolicy["CAN_SDD"] = false
	dnapolicy["DSPREFERENCE"] = 0
	dnapolicy["DSPZOFFSET"] = 0.5
	dnapolicy["TIP_REUSE_LIMIT"] = 0
	dnapolicy["NO_AIR_DISPENSE"] = true
	return dnapolicy

}

func PostMixPolicy() LHPolicy {
	dnapolicy := make(LHPolicy, 12)
	dnapolicy["POST_MIX"] = 3
	dnapolicy["POST_MIX_VOLUME"] = 20.0
	dnapolicy["POST_MIX_RATE"] = 3.74
	//dnapolicy["PRE_MIX"] = 3
	//dnapolicy["PRE_MIX_VOLUME"] = 10
	//dnapolicy["PRE_MIX_RATE"] = 3.74
	dnapolicy["ASPSPEED"] = 3.74
	dnapolicy["DSPSPEED"] = 3.74
	dnapolicy["CAN_MULTI"] = false
	dnapolicy["CAN_MSA"] = false
	dnapolicy["CAN_SDD"] = false
	dnapolicy["DSPREFERENCE"] = 0
	dnapolicy["DSPZOFFSET"] = 0.5
	dnapolicy["TIP_REUSE_LIMIT"] = 0
	dnapolicy["NO_AIR_DISPENSE"] = true
	return dnapolicy

}

func MegaMixPolicy() LHPolicy {
	dnapolicy := make(LHPolicy, 12)
	dnapolicy["POST_MIX"] = 10
	dnapolicy["POST_MIX_RATE"] = 3.74
	dnapolicy["ASPSPEED"] = 3.74
	dnapolicy["DSPSPEED"] = 3.74
	dnapolicy["CAN_MULTI"] = false
	dnapolicy["CAN_MSA"] = false
	dnapolicy["CAN_SDD"] = false
	dnapolicy["DSPREFERENCE"] = 0
	dnapolicy["DSPZOFFSET"] = 0.5
	dnapolicy["TIP_REUSE_LIMIT"] = 0
	dnapolicy["NO_AIR_DISPENSE"] = true
	return dnapolicy

}

func MakeDefaultPolicy() LHPolicy {
	defaultpolicy := make(LHPolicy, 27)
	defaultpolicy["OFFSETZADJUST"] = 0.0
	defaultpolicy["TOUCHOFF"] = false
	defaultpolicy["TOUCHOFFSET"] = 0.5
	defaultpolicy["ASPREFERENCE"] = 0
	defaultpolicy["ASPZOFFSET"] = 0.5
	defaultpolicy["DSPREFERENCE"] = 0
	defaultpolicy["DSPZOFFSET"] = 0.5
	defaultpolicy["CAN_MULTI"] = true
	defaultpolicy["CAN_MSA"] = false
	defaultpolicy["CAN_SDD"] = true
	defaultpolicy["TIP_REUSE_LIMIT"] = 100
	defaultpolicy["BLOWOUTREFERENCE"] = 1
	defaultpolicy["BLOWOUTVOLUME"] = 50.0
	defaultpolicy["BLOWOUTOFFSET"] = 0.0 //-5.0
	defaultpolicy["BLOWOUTVOLUMEUNIT"] = "ul"
	defaultpolicy["PTZREFERENCE"] = 1
	defaultpolicy["PTZOFFSET"] = -0.5
	defaultpolicy["NO_AIR_DISPENSE"] = true
	defaultpolicy["DEFAULTPIPETTESPEED"] = 3.0
	defaultpolicy["MANUALPTZ"] = false
	defaultpolicy["JUSTBLOWOUT"] = false
	defaultpolicy["DONT_BE_DIRTY"] = true
	defaultpolicy["POST_MIX_Z"] = 0.5
	defaultpolicy["PRE_MIX_Z"] = 0.5

	return defaultpolicy
}

func MakeJBPolicy() LHPolicy {
	jbp := make(LHPolicy, 1)
	jbp["JUSTBLOWOUT"] = true
	//jbp["TOUCHOFF"] = true
	return jbp
}

func MakeTOPolicy() LHPolicy {
	top := make(LHPolicy, 1)
	top["TOUCHOFF"] = true
	return top
}

func MakeLVExtraPolicy() LHPolicy {
	lvep := make(LHPolicy, 2)
	lvep["EXTRA_ASP_VOLUME"] = wunit.NewVolume(0.5, "ul")
	lvep["EXTRA_DISP_VOLUME"] = wunit.NewVolume(0.5, "ul")
	return lvep
}

func MakeLVDNAMixPolicy() LHPolicy {
	dnapolicy := make(LHPolicy, 4)
	dnapolicy["RESET_OVERRIDE"] = true
	dnapolicy["POST_MIX_VOLUME"] = 5.0
	dnapolicy["POST_MIX"] = 1
	dnapolicy["POST_MIX_Z"] = 0.5
	dnapolicy["POST_MIX_RATE"] = 3.0
	dnapolicy["TOUCHOFF"] = false
	return dnapolicy
}

func MakeEGEL48Policy() LHPolicy {
	loadpolicy := make(LHPolicy, 1)
	loadpolicy["RESET_OVERRIDE"] = true
	return loadpolicy
}

func MakeEGEL961Policy() LHPolicy {
	loadpolicy := make(LHPolicy, 1)
	loadpolicy["RESET_OVERRIDE"] = true
	return loadpolicy
}

func MakeEGEL962Policy() LHPolicy {
	loadpolicy := make(LHPolicy, 1)
	loadpolicy["RESET_OVERRIDE"] = true
	return loadpolicy
}

func MakeHVOffsetPolicy() LHPolicy {
	lvop := make(LHPolicy, 6)
	lvop["OFFSETZADJUST"] = 0.75
	lvop["POST_MIX_RATE"] = 37
	lvop["PRE_MIX_RATE"] = 37
	lvop["ASPSPEED"] = 37
	lvop["DSPSPEED"] = 37
	return lvop
}

// deprecated; see above
func MakeHVFlowRatePolicy() LHPolicy {
	policy := make(LHPolicy, 4)
	policy["POST_MIX_RATE"] = 37
	policy["PRE_MIX_RATE"] = 37
	policy["ASPSPEED"] = 37
	policy["DSPSPEED"] = 37
	return policy
}

func MakeCarbonSourcePolicy() LHPolicy {
	cspolicy := make(LHPolicy, 1)
	cspolicy["DSPREFERENCE"] = 1
	return cspolicy
}

func MakeNitrogenSourcePolicy() LHPolicy {
	nspolicy := make(LHPolicy, 1)
	nspolicy["DSPREFERENCE"] = 1
	return nspolicy
}

func GetLHPolicyForTest() (*LHPolicyRuleSet, error) {

	// make some policies

	policies := MakePolicies()

	// now make rules

	lhpr := NewLHPolicyRuleSet()

	for name, policy := range policies {
		rule := NewLHPolicyRule(name)
		err := rule.AddCategoryConditionOn("LIQUIDCLASS", name)

		if err != nil {
			return nil, err
		}
		lhpr.AddRule(rule, policy)
	}

	// hack to fix plate type problems
	// this really should be removed asap
	rule := NewLHPolicyRule("HVOffsetFix")
	//rule.AddNumericConditionOn("VOLUME", 20.1, 300.0) // what about higher? // set specifically for openPlant configuration

	rule.AddCategoryConditionOn("TIPTYPE", "Gilson200")
	rule.AddCategoryConditionOn("PLATFORM", "GilsonPipetmax")
	// don't get overridden
	rule.Priority = 100
	pol := MakeHVOffsetPolicy()
	lhpr.AddRule(rule, pol)

	// merged the below and the above
	/*
		rule = NewLHPolicyRule("HVFlowRate")
		rule.AddNumericConditionOn("VOLUME", 20.1, 300.0) // what about higher? // set specifically for openPlant configuration
		//rule.AddCategoryConditionOn("FROMPLATETYPE", "pcrplate_skirted_riser")
		pol = MakeHVFlowRatePolicy()
		lhpr.AddRule(rule, pol)
	*/

	rule = NewLHPolicyRule("DNALV")
	rule.AddNumericConditionOn("VOLUME", 0.0, 1.99)
	rule.AddCategoryConditionOn("LIQUIDCLASS", "dna")
	pol = MakeLVDNAMixPolicy()
	lhpr.AddRule(rule, pol)

	//fix for removing blowout in DNA only if EGEL 48 plate type is used
	rule = NewLHPolicyRule("EGEL48Load")
	rule.AddCategoryConditionOn("TOPLATETYPE", "EPAGE48")
	pol = MakeEGEL48Policy()
	lhpr.AddRule(rule, pol)

	//fix for removing blowout in DNA only if EGEL 96_1 plate type is used
	rule = NewLHPolicyRule("EGEL961Load")
	rule.AddCategoryConditionOn("TOPLATETYPE", "EGEL96_1")
	pol = MakeEGEL961Policy()
	lhpr.AddRule(rule, pol)

	//fix for removing blowout in DNA only if EGEL 96_2 plate type is used
	rule = NewLHPolicyRule("EGEL962Load")
	rule.AddCategoryConditionOn("TOPLATETYPE", "EGEL96_2")
	pol = MakeEGEL962Policy()

	lhpr.AddRule(rule, pol)

	return lhpr, nil

}

func LoadLHPoliciesFromFile() (*LHPolicyRuleSet, error) {
	lhPoliciesFileName := os.Getenv("ANTHA_LHPOLICIES_FILE")
	if lhPoliciesFileName == "" {
		return nil, fmt.Errorf("Env variable ANTHA_LHPOLICIES_FILE not set")
	}
	contents, err := ioutil.ReadFile(lhPoliciesFileName)
	if err != nil {
		return nil, err
	}
	lhprs := NewLHPolicyRuleSet()
	lhprs.Policies = make(map[string]LHPolicy)
	lhprs.Rules = make(map[string]LHPolicyRule)
	//	err = readYAML(contents, lhprs)
	err = readJSON(contents, lhprs)
	if err != nil {
		return nil, err
	}
	return lhprs, nil
}

func readYAML(fileContents []byte, ruleSet *LHPolicyRuleSet) error {
	if err := yaml.Unmarshal(fileContents, ruleSet); err != nil {
		return err
	}
	return nil
}

func readJSON(fileContents []byte, ruleSet *LHPolicyRuleSet) error {
	if err := json.Unmarshal(fileContents, ruleSet); err != nil {
		return err
	}
	return nil
}
