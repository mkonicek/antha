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
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
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
//var DOEliquidhandlingFile = "170516CCFDesign_noTouchoff_noBlowout.xlsx" // "2700516AssemblyCCF.xlsx" //"newdesign2factorsonly.xlsx" // "170516CCFDesign_noTouchoff_noBlowout.xlsx" // "170516CFF.xlsx" //"newdesign2factorsonly.xlsx" "170516CCFDesign_noTouchoff_noBlowout.xlsx" // //"newdesign2factorsonly.xlsx" //"8run4cpFactorial.xlsx" //"FullFactorial.xlsx" // "Screenwtype.LHPolicyDOE2.xlsx"
//var DXORJMP = "DX"                                                      //"JMP"
var BASEPolicy = "default" //"dna"
func MakePolicies() map[string]wtype.LHPolicy {
	pols := make(map[string]wtype.LHPolicy)

	// what policies do we need?
	pols["SmartMix"] = SmartMixPolicy()
	pols["water"] = MakeWaterPolicy()
	pols["multiwater"] = MakeMultiWaterPolicy()
	pols["culture"] = MakeCulturePolicy()
	pols["culturereuse"] = MakeCultureReusePolicy()
	pols["glycerol"] = MakeGlycerolPolicy()
	pols["solvent"] = MakeSolventPolicy()
	pols["default"] = MakeDefaultPolicy()
	pols["dna"] = MakeDNAPolicy()
	pols["DoNotMix"] = MakeDefaultPolicy()
	pols["NeedToMix"] = MakeNeedToMixPolicy()
	pols["PreMix"] = PreMixPolicy()
	pols["PostMix"] = PostMixPolicy()
	pols["MegaMix"] = MegaMixPolicy()
	pols["viscous"] = MakeViscousPolicy()
	pols["Paint"] = MakePaintPolicy()

	// pols["lysate"] = MakeLysatePolicy()
	pols["protein"] = MakeProteinPolicy()
	pols["detergent"] = MakeDetergentPolicy()
	pols["load"] = MakeLoadPolicy()
	pols["loadwater"] = MakeLoadWaterPolicy()
	pols["DispenseAboveLiquid"] = MakeDispenseAboveLiquidPolicy()
	pols["DispenseAboveLiquidMulti"] = MakeDispenseAboveLiquidMultiPolicy()
	pols["PEG"] = MakePEGPolicy()
	pols["Protoplasts"] = MakeProtoplastPolicy()
	pols["dna_mix"] = MakeDNAMixPolicy()
	pols["dna_mix_multi"] = MakeDNAMixMultiPolicy()
	pols["dna_cells_mix"] = MakeDNACELLSMixPolicy()
	pols["dna_cells_mix_multi"] = MakeDNACELLSMixMultiPolicy()
	pols["plateout"] = MakePlateOutPolicy()
	pols["colony"] = MakeColonyPolicy()
	pols["colonymix"] = MakeColonyMixPolicy()
	//      pols["lysate"] = MakeLysatePolicy()
	pols["carbon_source"] = MakeCarbonSourcePolicy()
	pols["nitrogen_source"] = MakeNitrogenSourcePolicy()
	pols["XYOffsetTest"] = MakeXYOffsetTestPolicy()

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

func PolicyMakerfromFilename(filename string) (policies []wtype.LHPolicy, names []string, runs []Run, err error) {

	doeliquidhandlingFile, found := PolicyFilefromName(filename)
	if !found {
		panic("policyfilename" + filename + "not found")
	}
	filenameparts := strings.Split(doeliquidhandlingFile.Filename, ".")

	policies, names, runs, err = PolicyMakerfromDesign(BASEPolicy, doeliquidhandlingFile.DXORJMP, doeliquidhandlingFile.Filename, filenameparts[0])
	return
}

func PolicyMakerfromDesign(basepolicy string, DXORJMP string, dxdesignfilename string, prepend string) (policies []wtype.LHPolicy, names []string, runs []Run, err error) {

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

func PolicyMaker(basepolicy string, factors []DOEPair, nameprepend string, concatfactorlevelsinname bool) (policies []wtype.LHPolicy, names []string) {

	runs := AllCombinations(factors)

	policies, names = PolicyMakerfromRuns(basepolicy, runs, nameprepend, concatfactorlevelsinname)

	return
}

func PolicyMakerfromRuns(basepolicy string, runs []Run, nameprepend string, concatfactorlevelsinname bool) (policies []wtype.LHPolicy, names []string) {

	policyitemmap := wtype.MakePolicyItems()

	names = make([]string, 0)
	policies = make([]wtype.LHPolicy, 0)

	policy := MakeDefaultPolicy()
	err := policy.Set("CAN_MULTI", false)
	if err != nil {
		panic(err)
	}

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

//func MakeLysatePolicy() wtype.LHPolicy {
//        lysatepolicy := make(wtype.LHPolicy, 6)
//        lysatepolicy["ASPSPEED"] = 1.0
//        lysatepolicy["DSPSPEED"] = 1.0
//        lysatepolicy["ASP_WAIT"] = 2.0
//        lysatepolicy["ASP_WAIT"] = 2.0
//        lysatepolicy["DSP_WAIT"] = 2.0
//        lysatepolicy["PRE_MIX"] = 5
//        lysatepolicy["CAN_MSA"]= false
//        return lysatepolicy
//}
//func MakeProteinPolicy() wtype.LHPolicy {
//        proteinpolicy := make(wtype.LHPolicy, 4)
//        proteinpolicy["DSPREFERENCE"] = 2
//        proteinpolicy["CAN_MULTI"] = true
//        proteinpolicy["PRE_MIX"] = 3
//        proteinpolicy["CAN_MSA"] = false
//        return proteinpolicy
//}

func GetPolicyByName(policyname wtype.PolicyName) (lhpolicy wtype.LHPolicy, policypresent bool) {
	policymap := MakePolicies()

	lhpolicy, policypresent = policymap[policyname.String()]
	return
}

func AvailablePolicies() (policies []string) {

	policies = make([]string, 0)
	policymap := MakePolicies()

	for key := range policymap {
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

func MakePEGPolicy() wtype.LHPolicy {
	policy := make(wtype.LHPolicy, 10)
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
	policy["CAN_MULTI"] = true
	policy["RESET_OVERRIDE"] = true
	policy["DESCRIPTION"] = "Customised for handling Poly Ethylene Glycol solutions. Similar to mixing required for viscous solutions. 3 post-mixes."
	return policy
}

func MakeProtoplastPolicy() wtype.LHPolicy {
	policy := make(wtype.LHPolicy, 8)
	policy["ASPSPEED"] = 0.5
	policy["DSPSPEED"] = 0.5
	policy["ASPZOFFSET"] = 1.0
	policy["DSPZOFFSET"] = 1.0
	policy["BLOWOUTVOLUME"] = 100.0
	policy["BLOWOUTVOLUMEUNIT"] = "ul"
	policy["TOUCHOFF"] = false
	policy["TIP_REUSE_LIMIT"] = 5
	policy["CAN_MULTI"] = true
	policy["DESCRIPTION"] = "Customised for handling protoplast solutions. Pipettes very gently. No post-mix."
	return policy
}

func MakePaintPolicy() wtype.LHPolicy {
	policy := make(wtype.LHPolicy, 14)
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
	policy["CAN_MULTI"] = true
	policy["DESCRIPTION"] = "Customised for handling paint solutions. Similar to mixing required for viscous solutions. 3 post-mixes."
	return policy
}

func MakeDispenseAboveLiquidPolicy() wtype.LHPolicy {
	policy := make(wtype.LHPolicy, 8)
	policy["DSPREFERENCE"] = 1 // 1 indicates dispense at top of well
	policy["ASPSPEED"] = 3.0
	policy["DSPSPEED"] = 3.0
	//policy["ASP_WAIT"] = 1.0
	//policy["DSP_WAIT"] = 1.0
	policy["BLOWOUTVOLUME"] = 50.0
	policy["BLOWOUTVOLUMEUNIT"] = "ul"
	policy["TOUCHOFF"] = false
	policy["CAN_MULTI"] = false
	policy["DESCRIPTION"] = "Dispense solution above the liquid to facilitate tip reuse but sacrifice pipetting accuracy at low volumes. No post-mix. No multi channel"
	return policy
}
func MakeDispenseAboveLiquidMultiPolicy() wtype.LHPolicy {
	policy := make(wtype.LHPolicy, 8)
	policy["DSPREFERENCE"] = 1 // 1 indicates dispense at top of well
	policy["ASPSPEED"] = 3.0
	policy["DSPSPEED"] = 3.0
	//policy["ASP_WAIT"] = 1.0
	//policy["DSP_WAIT"] = 1.0
	policy["BLOWOUTVOLUME"] = 50.0
	policy["BLOWOUTVOLUMEUNIT"] = "ul"
	policy["TOUCHOFF"] = false
	policy["CAN_MULTI"] = true
	policy["DESCRIPTION"] = "Dispense solution above the liquid to facilitate tip reuse but sacrifice pipetting accuracy at low volumes. No post Mix. Allows multi-channel pipetting."
	return policy
}

func MakeColonyPolicy() wtype.LHPolicy {
	policy := make(wtype.LHPolicy, 12)
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
	policy["DESCRIPTION"] = "Designed for colony picking. 1 post-mix and no blowout (to avoid potential cross contamination), no multichannel."
	return policy
}

func MakeColonyMixPolicy() wtype.LHPolicy {
	policy := MakeColonyPolicy()
	policy["POST_MIX"] = 3
	policy["DESCRIPTION"] = "Designed for colony picking but with added post-mixes. 3 post-mix and no blowout (to avoid potential cross contamination), no multichannel."
	return policy
}

func MakeXYOffsetTestPolicy() wtype.LHPolicy {
	policy := MakeColonyPolicy()
	policy["POST_MIX_X"] = 2.0
	policy["POST_MIX_Y"] = 2.0
	policy["DESCRIPTION"] = "Intended to test setting X,Y offsets for post mixing"
	return policy
}

func MakeWaterPolicy() wtype.LHPolicy {
	waterpolicy := make(wtype.LHPolicy, 6)
	waterpolicy["DSPREFERENCE"] = 0
	waterpolicy["CAN_MSA"] = true
	waterpolicy["CAN_SDD"] = true
	waterpolicy["CAN_MULTI"] = false
	waterpolicy["DSPZOFFSET"] = 1.0
	waterpolicy["BLOWOUTVOLUME"] = 50.0
	waterpolicy["DESCRIPTION"] = "Default policy designed for pipetting water. Includes a blowout step for added accuracy and no post-mixing, no multi channel."
	return waterpolicy
}

func MakeMultiWaterPolicy() wtype.LHPolicy {
	pol := MakeWaterPolicy()
	pol["CAN_MULTI"] = true
	pol["DESCRIPTION"] = "Default policy designed for pipetting water but permitting multi-channel use. Includes a blowout step for added accuracy and no post-mixing."
	return pol
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func MakeCulturePolicy() wtype.LHPolicy {
	culturepolicy := make(wtype.LHPolicy, 10)
	checkErr(culturepolicy.Set("PRE_MIX", 2))
	checkErr(culturepolicy.Set("PRE_MIX_VOLUME", 19.0))
	checkErr(culturepolicy.Set("PRE_MIX_RATE", 3.74))
	checkErr(culturepolicy.Set("ASPSPEED", 2.0))
	checkErr(culturepolicy.Set("DSPSPEED", 2.0))
	checkErr(culturepolicy.Set("CAN_MULTI", true))
	checkErr(culturepolicy.Set("CAN_MSA", false))
	checkErr(culturepolicy.Set("CAN_SDD", false))
	checkErr(culturepolicy.Set("DSPREFERENCE", 0))
	checkErr(culturepolicy.Set("DSPZOFFSET", 0.5))
	checkErr(culturepolicy.Set("TIP_REUSE_LIMIT", 0))
	checkErr(culturepolicy.Set("NO_AIR_DISPENSE", true))
	checkErr(culturepolicy.Set("BLOWOUTVOLUME", 0.0))
	checkErr(culturepolicy.Set("BLOWOUTVOLUMEUNIT", "ul"))
	checkErr(culturepolicy.Set("TOUCHOFF", false))
	checkErr(culturepolicy.Set("DESCRIPTION", "Designed for cell cultures. Tips will not be reused to minimise any risk of cross contamination and 2 pre-mixes will be performed prior to aspirating."))
	return culturepolicy
}

func MakePlateOutPolicy() wtype.LHPolicy {
	culturepolicy := make(wtype.LHPolicy, 17)
	culturepolicy["CAN_MULTI"] = true
	culturepolicy["ASP_WAIT"] = 1.0
	culturepolicy["DSP_WAIT"] = 1.0
	culturepolicy["DSPZOFFSET"] = 0.0
	culturepolicy["TIP_REUSE_LIMIT"] = 7
	culturepolicy["NO_AIR_DISPENSE"] = true
	culturepolicy["TOUCHOFF"] = false
	culturepolicy["RESET_OVERRIDE"] = true
	culturepolicy["DESCRIPTION"] = "Designed for plating out cultures onto agar plates. Dispense will be performed at the well bottom and no blowout will be performed (to minimise risk of cross contamination)"
	return culturepolicy
}

func MakeCultureReusePolicy() wtype.LHPolicy {
	culturepolicy := make(wtype.LHPolicy, 10)
	checkErr(culturepolicy.Set("PRE_MIX", 2))
	checkErr(culturepolicy.Set("PRE_MIX_VOLUME", 19.0))
	checkErr(culturepolicy.Set("PRE_MIX_RATE", 3.74))
	checkErr(culturepolicy.Set("ASPSPEED", 2.0))
	checkErr(culturepolicy.Set("DSPSPEED", 2.0))
	checkErr(culturepolicy.Set("CAN_MULTI", true))
	checkErr(culturepolicy.Set("CAN_MSA", true))
	checkErr(culturepolicy.Set("CAN_SDD", true))
	checkErr(culturepolicy.Set("DSPREFERENCE", 0))
	checkErr(culturepolicy.Set("DSPZOFFSET", 0.5))
	checkErr(culturepolicy.Set("NO_AIR_DISPENSE", true))
	checkErr(culturepolicy.Set("BLOWOUTVOLUME", 0.0))
	checkErr(culturepolicy.Set("BLOWOUTVOLUMEUNIT", "ul"))
	checkErr(culturepolicy.Set("TOUCHOFF", false))
	checkErr(culturepolicy.Set("DESCRIPTION", "Designed for cell cultures but permitting tip reuse when handling the same culture. 2 pre-mixes will be performed prior to aspirating."))
	return culturepolicy
}

func MakeGlycerolPolicy() wtype.LHPolicy {
	glycerolpolicy := make(wtype.LHPolicy, 9)
	glycerolpolicy["ASPSPEED"] = 1.5
	glycerolpolicy["DSPSPEED"] = 1.5
	glycerolpolicy["ASP_WAIT"] = 1.0
	glycerolpolicy["DSP_WAIT"] = 1.0
	glycerolpolicy["TIP_REUSE_LIMIT"] = 0
	glycerolpolicy["CAN_MULTI"] = true
	glycerolpolicy["POST_MIX"] = 3
	glycerolpolicy["POST_MIX_VOLUME"] = 20.0
	glycerolpolicy["POST_MIX_RATE"] = 3.74 // Should this be the same rate as the asp and dsp speeds?
	glycerolpolicy["DESCRIPTION"] = "Designed for viscous samples, in particular enzymes stored in glycerol. 3 gentle post-mixes of 20ul will be performed. Tips will not be reused in order to increase accuracy."
	return glycerolpolicy
}

func MakeViscousPolicy() wtype.LHPolicy {
	glycerolpolicy := make(wtype.LHPolicy, 7)
	glycerolpolicy["ASPSPEED"] = 1.5
	glycerolpolicy["DSPSPEED"] = 1.5
	glycerolpolicy["ASP_WAIT"] = 1.0
	glycerolpolicy["DSP_WAIT"] = 1.0
	glycerolpolicy["CAN_MULTI"] = true
	glycerolpolicy["POST_MIX"] = 3
	glycerolpolicy["POST_MIX_RATE"] = 1.5
	glycerolpolicy["DESCRIPTION"] = "Designed for viscous samples. 3 post-mixes of the volume of the sample being transferred will be performed. No tip reuse limit."
	return glycerolpolicy
}
func MakeSolventPolicy() wtype.LHPolicy {
	solventpolicy := make(wtype.LHPolicy, 5)
	checkErr(solventpolicy.Set("PRE_MIX", 3))
	checkErr(solventpolicy.Set("DSPREFERENCE", 0))
	checkErr(solventpolicy.Set("DSPZOFFSET", 0.5))
	checkErr(solventpolicy.Set("NO_AIR_DISPENSE", true))
	checkErr(solventpolicy.Set("CAN_MULTI", true))
	checkErr(solventpolicy.Set("DESCRIPTION", "Designed for handling solvents. No post-mixes are performed"))
	return solventpolicy
}

func MakeDNAPolicy() wtype.LHPolicy {
	dnapolicy := make(wtype.LHPolicy, 12)
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
	dnapolicy["DESCRIPTION"] = "Designed for DNA samples. No tip reuse is permitted, no blowout and no post-mixing."
	return dnapolicy
}

func MakeDNAMixPolicy() wtype.LHPolicy {
	dnapolicy := MakeDNAPolicy()
	dnapolicy["POST_MIX_VOLUME"] = 10.0
	dnapolicy["POST_MIX"] = 5
	dnapolicy["POST_MIX_Z"] = 0.5
	dnapolicy["POST_MIX_RATE"] = 3.0
	dnapolicy["CAN_MULTI"] = false
	dnapolicy["DESCRIPTION"] = "Designed for DNA samples but with 5 post-mixes of 10ul. No tip reuse is permitted, no blowout, no multichannel."
	return dnapolicy
}

func MakeDNAMixMultiPolicy() wtype.LHPolicy {
	dnapolicy := MakeDNAPolicy()
	dnapolicy["POST_MIX_VOLUME"] = 10.0
	dnapolicy["POST_MIX"] = 5
	dnapolicy["POST_MIX_Z"] = 0.5
	dnapolicy["POST_MIX_RATE"] = 3.0
	dnapolicy["CAN_MULTI"] = true
	dnapolicy["DESCRIPTION"] = "Designed for DNA samples but with 5 post-mixes of 10ul. No tip reuse is permitted, no blowout. Allows multi-channel pipetting."
	return dnapolicy
}

func MakeDNACELLSMixPolicy() wtype.LHPolicy {
	dnapolicy := MakeDNAPolicy()
	dnapolicy["POST_MIX_VOLUME"] = 20.0
	dnapolicy["POST_MIX"] = 2
	dnapolicy["POST_MIX_Z"] = 0.5
	dnapolicy["POST_MIX_RATE"] = 1.0
	dnapolicy["DESCRIPTION"] = "Designed for mixing DNA with cells. 2 gentle post-mixes are performed. No tip reuse is permitted, no blowout."
	return dnapolicy
}
func MakeDNACELLSMixMultiPolicy() wtype.LHPolicy {
	dnapolicy := MakeDNAPolicy()
	dnapolicy["POST_MIX_VOLUME"] = 20.0
	dnapolicy["POST_MIX"] = 2
	dnapolicy["POST_MIX_Z"] = 0.5
	dnapolicy["POST_MIX_RATE"] = 1.0
	dnapolicy["CAN_MULTI"] = true
	dnapolicy["DESCRIPTION"] = "Designed for mixing DNA with cells. 2 gentle post-mixes are performed. No tip reuse is permitted, no blowout. Allows multi-channel pipetting."
	return dnapolicy
}

func MakeDetergentPolicy() wtype.LHPolicy {
	detergentpolicy := make(wtype.LHPolicy, 9)
	//        detergentpolicy["POST_MIX"] = 3
	detergentpolicy["ASPSPEED"] = 1.0
	detergentpolicy["DSPSPEED"] = 1.0
	detergentpolicy["CAN_MSA"] = false
	detergentpolicy["CAN_SDD"] = false
	detergentpolicy["DSPREFERENCE"] = 0
	detergentpolicy["DSPZOFFSET"] = 0.5
	detergentpolicy["TIP_REUSE_LIMIT"] = 8
	detergentpolicy["NO_AIR_DISPENSE"] = true
	detergentpolicy["DESCRIPTION"] = "Designed for solutions containing detergents. Gentle aspiration and dispense and a tip reuse limit of 8 to reduce problem of foam build up inside the tips."
	return detergentpolicy
}
func MakeProteinPolicy() wtype.LHPolicy {
	proteinpolicy := make(wtype.LHPolicy, 12)
	proteinpolicy["POST_MIX"] = 5
	proteinpolicy["POST_MIX_VOLUME"] = 50.0
	proteinpolicy["ASPSPEED"] = 2.0
	proteinpolicy["DSPSPEED"] = 2.0
	proteinpolicy["CAN_MSA"] = false
	proteinpolicy["CAN_SDD"] = false
	proteinpolicy["DSPREFERENCE"] = 0
	proteinpolicy["DSPZOFFSET"] = 0.5
	proteinpolicy["TIP_REUSE_LIMIT"] = 0
	proteinpolicy["NO_AIR_DISPENSE"] = true
	proteinpolicy["DESCRIPTION"] = "Designed for protein solutions. Slightly gentler aspiration and dispense and a tip reuse limit of 0 to prevent risk of cross contamination. 5 post-mixes of 50ul will be performed."
	return proteinpolicy
}
func MakeLoadPolicy() wtype.LHPolicy {

	loadpolicy := make(wtype.LHPolicy, 14)
	loadpolicy["ASPSPEED"] = 1.0
	loadpolicy["DSPSPEED"] = 0.1
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
	loadpolicy["DESCRIPTION"] = "Designed for loading a sample onto an agarose gel. Very slow dispense rate, no tip reuse and no blowout."
	return loadpolicy
}

func MakeLoadWaterPolicy() wtype.LHPolicy {
	loadpolicy := make(wtype.LHPolicy)
	loadpolicy["ASPSPEED"] = 1.0
	loadpolicy["DSPSPEED"] = 0.1
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
	loadpolicy["DESCRIPTION"] = "Designed for loading water into agarose gel wells so permits tip reuse. Very slow dispense rate and no blowout."
	return loadpolicy
}

func MakeNeedToMixPolicy() wtype.LHPolicy {
	dnapolicy := make(wtype.LHPolicy, 16)
	dnapolicy["POST_MIX"] = 3
	dnapolicy["POST_MIX_RATE"] = 3.74
	dnapolicy["PRE_MIX"] = 3
	dnapolicy["PRE_MIX_VOLUME"] = 20.0
	dnapolicy["PRE_MIX_RATE"] = 3.74
	dnapolicy["ASPSPEED"] = 3.74
	dnapolicy["DSPSPEED"] = 3.74
	dnapolicy["CAN_MULTI"] = true
	dnapolicy["CAN_MSA"] = false
	dnapolicy["CAN_SDD"] = false
	dnapolicy["DSPREFERENCE"] = 0
	dnapolicy["DSPZOFFSET"] = 0.5
	dnapolicy["TIP_REUSE_LIMIT"] = 0
	dnapolicy["NO_AIR_DISPENSE"] = true
	dnapolicy["DESCRIPTION"] = "3 pre-mixes and 3 post-mixes of the sample being transferred.  No tip reuse permitted."
	return dnapolicy
}

func PreMixPolicy() wtype.LHPolicy {
	dnapolicy := make(wtype.LHPolicy, 12)
	//dnapolicy["POST_MIX"] = 3
	//dnapolicy[""POST_MIX_VOLUME"] = 10.0
	//dnapolicy["POST_MIX_RATE"] = 3.74
	dnapolicy["PRE_MIX"] = 3
	dnapolicy["PRE_MIX_VOLUME"] = 19.0
	dnapolicy["PRE_MIX_RATE"] = 3.74
	dnapolicy["ASPSPEED"] = 3.74
	dnapolicy["DSPSPEED"] = 3.74
	dnapolicy["CAN_MULTI"] = true
	dnapolicy["CAN_MSA"] = false
	dnapolicy["CAN_SDD"] = false
	dnapolicy["DSPREFERENCE"] = 0
	dnapolicy["DSPZOFFSET"] = 0.5
	dnapolicy["TIP_REUSE_LIMIT"] = 0
	dnapolicy["NO_AIR_DISPENSE"] = true
	dnapolicy["DESCRIPTION"] = "3 pre-mixes of the sample being transferred.  No tip reuse permitted."
	return dnapolicy

}

func PostMixPolicy() wtype.LHPolicy {
	dnapolicy := make(wtype.LHPolicy, 12)
	dnapolicy["POST_MIX"] = 3
	dnapolicy["POST_MIX_RATE"] = 3.74
	//dnapolicy["PRE_MIX"] = 3
	//dnapolicy["PRE_MIX_VOLUME"] = 10
	//dnapolicy["PRE_MIX_RATE"] = 3.74
	dnapolicy["ASPSPEED"] = 3.74
	dnapolicy["DSPSPEED"] = 3.74
	dnapolicy["CAN_MULTI"] = true
	dnapolicy["CAN_MSA"] = false
	dnapolicy["CAN_SDD"] = false
	dnapolicy["DSPREFERENCE"] = 0
	dnapolicy["DSPZOFFSET"] = 0.5
	dnapolicy["TIP_REUSE_LIMIT"] = 0
	dnapolicy["NO_AIR_DISPENSE"] = true
	dnapolicy["DESCRIPTION"] = "3 post-mixes of the sample being transferred.  No tip reuse permitted."
	return dnapolicy
}

// 3 post mixes of the sample being transferred. Volume is adjusted based upon the volume of liquid in the destination well.
// No tip reuse permitted.
// Rules added to adjust post mix volume based on volume of the destination well.
// volume now capped at max for tip type (MIX_VOLUME_OVERRIDE_TIP_MAX)
func SmartMixPolicy() wtype.LHPolicy {
	policy := make(wtype.LHPolicy, 12)
	policy["POST_MIX"] = 3
	policy["POST_MIX_RATE"] = 3.74
	policy["ASPSPEED"] = 3.74
	policy["DSPSPEED"] = 3.74
	policy["CAN_MULTI"] = true
	policy["CAN_MSA"] = false
	policy["CAN_SDD"] = false
	policy["DSPREFERENCE"] = 0
	policy["DSPZOFFSET"] = 0.5
	policy["TIP_REUSE_LIMIT"] = 0
	policy["NO_AIR_DISPENSE"] = true
	policy["DESCRIPTION"] = "3 post-mixes of the sample being transferred. Volume is adjusted based upon the volume of liquid in the destination well.  No tip reuse permitted."
	policy["MIX_VOLUME_OVERRIDE_TIP_MAX"] = true
	return policy
}

func MegaMixPolicy() wtype.LHPolicy {
	dnapolicy := make(wtype.LHPolicy, 12)
	dnapolicy["POST_MIX"] = 10
	dnapolicy["POST_MIX_RATE"] = 3.74
	dnapolicy["ASPSPEED"] = 3.74
	dnapolicy["DSPSPEED"] = 3.74
	dnapolicy["CAN_MULTI"] = true
	dnapolicy["CAN_MSA"] = false
	dnapolicy["CAN_SDD"] = false
	dnapolicy["DSPREFERENCE"] = 0
	dnapolicy["DSPZOFFSET"] = 0.5
	dnapolicy["TIP_REUSE_LIMIT"] = 0
	dnapolicy["NO_AIR_DISPENSE"] = true
	dnapolicy["DESCRIPTION"] = "10 post-mixes of the sample being transferred. No tip reuse permitted."
	return dnapolicy

}

func MakeDefaultPolicy() wtype.LHPolicy {
	defaultpolicy := make(wtype.LHPolicy, 29)
	defaultpolicy["MIX_VOLUME_OVERRIDE_TIP_MAX"] = false
	defaultpolicy["OFFSETZADJUST"] = 0.0
	defaultpolicy["TOUCHOFF"] = false
	defaultpolicy["TOUCHOFFSET"] = 0.5
	defaultpolicy["ASPREFERENCE"] = 0
	defaultpolicy["ASPZOFFSET"] = 0.5
	defaultpolicy["DSPREFERENCE"] = 0
	defaultpolicy["DSPZOFFSET"] = 0.5
	defaultpolicy["CAN_MSA"] = false
	defaultpolicy["CAN_SDD"] = true
	defaultpolicy["CAN_MULTI"] = true
	defaultpolicy["TIP_REUSE_LIMIT"] = 100
	defaultpolicy["BLOWOUTREFERENCE"] = 1
	defaultpolicy["BLOWOUTVOLUME"] = 50.0
	defaultpolicy["BLOWOUTOFFSET"] = 0.0 //-5.0
	defaultpolicy["BLOWOUTVOLUMEUNIT"] = "ul"
	defaultpolicy["PTZREFERENCE"] = 1
	defaultpolicy["PTZOFFSET"] = -0.5
	defaultpolicy["NO_AIR_DISPENSE"] = true // SERIOUSLY??
	defaultpolicy["DEFAULTPIPETTESPEED"] = 3.0
	defaultpolicy["MANUALPTZ"] = false
	defaultpolicy["JUSTBLOWOUT"] = false
	defaultpolicy["DONT_BE_DIRTY"] = true
	defaultpolicy["POST_MIX_Z"] = 0.5
	defaultpolicy["PRE_MIX_Z"] = 0.5
	defaultpolicy["LLFABOVESURFACE"] = 3.0 //distance above liquid level for dispensing with LiqudLevelFollowing
	defaultpolicy["LLFBELOWSURFACE"] = 3.0 //distance below liquid level for aspirating with LLF
	defaultpolicy["DESCRIPTION"] = "Default mix Policy. Blowout performed, no touch off, no mixing, tip reuse permitted for the same solution."

	return defaultpolicy
}

func MakeJBPolicy() wtype.LHPolicy {
	jbp := make(wtype.LHPolicy, 1)
	checkErr(jbp.Set("JUSTBLOWOUT", true))
	checkErr(jbp.Set("TOUCHOFF", true))
	return jbp
}

func MakeTOPolicy() wtype.LHPolicy {
	top := make(wtype.LHPolicy, 1)
	checkErr(top.Set("TOUCHOFF", true))
	return top
}

func MakeLVExtraPolicy() wtype.LHPolicy {
	lvep := make(wtype.LHPolicy, 2)
	checkErr(lvep.Set("EXTRA_ASP_VOLUME", wunit.NewVolume(0.5, "ul")))
	checkErr(lvep.Set("EXTRA_DISP_VOLUME", wunit.NewVolume(0.5, "ul")))
	return lvep
}

func MakeLVDNAMixPolicy() wtype.LHPolicy {
	dnapolicy := make(wtype.LHPolicy, 4)
	dnapolicy["RESET_OVERRIDE"] = true
	dnapolicy["POST_MIX_VOLUME"] = 5.0
	dnapolicy["POST_MIX"] = 1
	dnapolicy["POST_MIX_Z"] = 0.5
	dnapolicy["POST_MIX_RATE"] = 3.0
	dnapolicy["TOUCHOFF"] = false
	return dnapolicy
}

func TurnOffBlowoutPolicy() wtype.LHPolicy {
	loadpolicy := make(wtype.LHPolicy, 1)
	loadpolicy["RESET_OVERRIDE"] = true
	return loadpolicy
}

func MakeHVOffsetPolicy() wtype.LHPolicy {
	lvop := make(wtype.LHPolicy, 6)
	lvop["OFFSETZADJUST"] = 0.75
	lvop["POST_MIX_RATE"] = 37
	lvop["PRE_MIX_RATE"] = 37
	lvop["ASPSPEED"] = 37
	lvop["DSPSPEED"] = 37
	return lvop
}

func AdjustPostMixVolume(mixToVol wunit.Volume) wtype.LHPolicy {
	vol := mixToVol.ConvertTo(wunit.ParsePrefixedUnit("ul"))
	policy := make(wtype.LHPolicy, 1)
	policy["POST_MIX_VOLUME"] = vol
	return policy
}

func AdjustPreMixVolume(mixToVol wunit.Volume) wtype.LHPolicy {
	vol := mixToVol.ConvertTo(wunit.ParsePrefixedUnit("ul"))
	policy := make(wtype.LHPolicy, 1)
	policy["PRE_MIX_VOLUME"] = vol
	return policy
}

// deprecated; see above
func MakeHVFlowRatePolicy() wtype.LHPolicy {
	policy := make(wtype.LHPolicy, 4)
	policy["POST_MIX_RATE"] = 37
	policy["PRE_MIX_RATE"] = 37
	policy["ASPSPEED"] = 37
	policy["DSPSPEED"] = 37
	return policy
}

func MakeCarbonSourcePolicy() wtype.LHPolicy {
	cspolicy := make(wtype.LHPolicy, 1)
	cspolicy["DSPREFERENCE"] = 1
	cspolicy["DESCRIPTION"] = "Custom policy for carbon source which dispenses above destination solution."
	return cspolicy
}

func MakeNitrogenSourcePolicy() wtype.LHPolicy {
	nspolicy := make(wtype.LHPolicy, 1)
	nspolicy["DSPREFERENCE"] = 1
	nspolicy["DESCRIPTION"] = "Custom policy for nitrogen source which dispenses above destination solution."
	return nspolicy
}

// newConditionalRule makes a new LHPolicyRule with conditions to apply to an LHPolicy.
//
// An error is returned if an invalid Condition Class or SetPoint is specified.
// The valid Setpoints can be found in wtype.MakeInstructionParameters()
func newConditionalRule(ruleName string, conditions ...condition) (wtype.LHPolicyRule, error) {
	var errs []string

	rule := wtype.NewLHPolicyRule(ruleName)
	for _, condition := range conditions {
		err := condition.AddToRule(rule)
		if err != nil {
			errs = append(errs, err.Error())
		}
	}
	if len(errs) > 0 {
		return rule, fmt.Errorf(strings.Join(errs, ".\n"))
	}
	return rule, nil
}

type condition interface {
	AddToRule(wtype.LHPolicyRule) error
}

type categoricCondition struct {
	Class    string
	SetPoint string
}

func (c categoricCondition) AddToRule(rule wtype.LHPolicyRule) error {
	return rule.AddCategoryConditionOn(c.Class, c.SetPoint)
}

type numericCondition struct {
	Class string
	Range conditionRange
}

type conditionRange struct {
	Lower float64
	Upper float64
}

func (c numericCondition) AddToRule(rule wtype.LHPolicyRule) error {
	return rule.AddNumericConditionOn(c.Class, c.Range.Lower, c.Range.Upper)
}

// Conditions to apply to LHpolicyRules based on liquid policy used
var (
	OnSmartMix  = categoricCondition{"LIQUIDCLASS", "SmartMix"}
	OnPostMix   = categoricCondition{"LIQUIDCLASS", "PostMix"}
	OnPreMix    = categoricCondition{"LIQUIDCLASS", "PreMix"}
	OnNeedToMix = categoricCondition{"LIQUIDCLASS", "NeedToMix"}
)

// Conditions to apply to LHpolicyRules based on volume of liquid that a sample is being pipetted into at the destination well
var (
	IntoLessThan20ul          = numericCondition{Class: "TOWELLVOLUME", Range: conditionRange{Lower: 0.0, Upper: 20.0}}
	IntoBetween20ulAnd50ul    = numericCondition{Class: "TOWELLVOLUME", Range: conditionRange{20.0, 50.0}}
	IntoBetween50ulAnd100ul   = numericCondition{Class: "TOWELLVOLUME", Range: conditionRange{50.0, 100.0}}
	IntoBetween100ulAnd200ul  = numericCondition{Class: "TOWELLVOLUME", Range: conditionRange{100.0, 200.0}}
	IntoBetween200ulAnd1000ul = numericCondition{Class: "TOWELLVOLUME", Range: conditionRange{200.0, 1000.0}}
)

// Conditions to apply to LHpolicyRules based on volume of liquid being transferred
var (
	LessThan20ul = numericCondition{Class: "VOLUME", Range: conditionRange{0.0, 20.0}}
)

// Conditions to apply to LHpolicyRules based on volume of liquid in source well from which a sample is taken
var (
	FromBetween100ulAnd200ul  = numericCondition{Class: "WELLFROMVOLUME", Range: conditionRange{100.0, 200.0}}
	FromBetween200ulAnd1000ul = numericCondition{Class: "WELLFROMVOLUME", Range: conditionRange{200.0, 1000.0}}
)

func GetLHPolicyForTest() (*wtype.LHPolicyRuleSet, error) {
	// make some policies

	policies := MakePolicies()

	// now make rules

	lhpr := wtype.NewLHPolicyRuleSet()

	for name, policy := range policies {
		rule := wtype.NewLHPolicyRule(name)
		err := rule.AddCategoryConditionOn("LIQUIDCLASS", name)

		if err != nil {
			return nil, err
		}
		lhpr.AddRule(rule, policy)
	}

	adjustPostMix, err := newConditionalRule("mixInto20ul", OnSmartMix, IntoBetween20ulAnd50ul)

	if err != nil {
		return lhpr, err
	}

	adjustVol20 := AdjustPostMixVolume(wunit.NewVolume(20, "ul"))
	adjustVol50 := AdjustPostMixVolume(wunit.NewVolume(50, "ul"))
	adjustVol100 := AdjustPostMixVolume(wunit.NewVolume(100, "ul"))
	adjustVol200 := AdjustPostMixVolume(wunit.NewVolume(200, "ul"))

	lhpr.AddRule(adjustPostMix, adjustVol20)

	adjustPostMix50, err := newConditionalRule("mixInto50ul", OnSmartMix, IntoBetween50ulAnd100ul)

	if err != nil {
		return lhpr, err
	}

	lhpr.AddRule(adjustPostMix50, adjustVol50)

	adjustPostMix100, err := newConditionalRule("mixInto100ul", OnSmartMix, IntoBetween100ulAnd200ul)

	if err != nil {
		return lhpr, err
	}

	lhpr.AddRule(adjustPostMix100, adjustVol100)

	adjustPostMix200, err := newConditionalRule("mixInto200ul", OnSmartMix, IntoBetween200ulAnd1000ul)

	if err != nil {
		return lhpr, err
	}

	lhpr.AddRule(adjustPostMix200, adjustVol200)

	// adjust original PostMix and NeedToMix policy to only set post mix volume if low volume.
	postmix20ul, err := newConditionalRule("PostMix20ul", OnPostMix, LessThan20ul)

	if err != nil {
		return lhpr, err
	}

	lhpr.AddRule(postmix20ul, adjustVol20)

	needToMix20ul, err := newConditionalRule("NeedToMix20ul", OnNeedToMix, LessThan20ul)

	if err != nil {
		return lhpr, err
	}

	lhpr.AddRule(needToMix20ul, adjustVol20)

	// now pre mix values for PreMix and NeedToMix
	adjustPreMixVol20 := AdjustPreMixVolume(wunit.NewVolume(20, "ul"))
	adjustPreMixVol100 := AdjustPreMixVolume(wunit.NewVolume(100, "ul"))
	adjustPreMixVol200 := AdjustPreMixVolume(wunit.NewVolume(200, "ul"))

	// PreMix
	adjustPreMix, err := newConditionalRule("preMix20ul", OnPreMix, LessThan20ul)

	if err != nil {
		return lhpr, err
	}

	lhpr.AddRule(adjustPreMix, adjustPreMixVol20)

	adjustPreMix100ul, err := newConditionalRule("PreMixFrom100ul", OnPreMix, FromBetween100ulAnd200ul)

	if err != nil {
		return lhpr, err
	}

	lhpr.AddRule(adjustPreMix100ul, adjustPreMixVol100)

	adjustPreMix200ul, err := newConditionalRule("PreMixFrom200ul", OnPreMix, FromBetween200ulAnd1000ul)

	if err != nil {
		return lhpr, err
	}

	lhpr.AddRule(adjustPreMix200ul, adjustPreMixVol200)

	// NeedToMix
	adjustNeedToMix, err := newConditionalRule("NeedToPreMix20ul", OnNeedToMix, LessThan20ul)

	if err != nil {
		return lhpr, err
	}

	lhpr.AddRule(adjustNeedToMix, adjustPreMixVol20)

	adjustNeedToMix100ul, err := newConditionalRule("NeedToPreMixFrom100ul", OnNeedToMix, FromBetween100ulAnd200ul)

	if err != nil {
		return lhpr, err
	}

	lhpr.AddRule(adjustNeedToMix100ul, adjustPreMixVol100)

	adjustNeedToMix200ul, err := newConditionalRule("NeedToPreMixFrom200ul", OnNeedToMix, FromBetween200ulAnd1000ul)

	if err != nil {
		return lhpr, err
	}

	lhpr.AddRule(adjustNeedToMix200ul, adjustPreMixVol200)

	// hack to fix plate type problems
	// this really should be removed asap
	rule := wtype.NewLHPolicyRule("HVOffsetFix")
	//rule.AddNumericConditionOn("VOLUME", 20.1, 300.0) // what about higher? // set specifically for openPlant configuration

	checkErr(rule.AddCategoryConditionOn("TIPTYPE", "Gilson200"))
	checkErr(rule.AddCategoryConditionOn("PLATFORM", "GilsonPipetmax"))
	// don't get overridden
	rule.Priority = 100
	pol := MakeHVOffsetPolicy()
	lhpr.AddRule(rule, pol)

	// merged the below and the above
	/*
		rule = wtype.NewLHPolicyRule("HVFlowRate")
		rule.AddNumericConditionOn("VOLUME", 20.1, 300.0) // what about higher? // set specifically for openPlant configuration
		//rule.AddCategoryConditionOn("FROMPLATETYPE", "pcrplate_skirted_riser")
		pol = MakeHVFlowRatePolicy()
		lhpr.AddRule(rule, pol)
	*/

	rule = wtype.NewLHPolicyRule("DNALV")
	checkErr(rule.AddNumericConditionOn("VOLUME", 0.0, 1.99))
	checkErr(rule.AddCategoryConditionOn("LIQUIDCLASS", "dna"))
	pol = MakeLVDNAMixPolicy()
	lhpr.AddRule(rule, pol)

	//fix for removing blowout in DNA only if EGEL 48 plate type is used
	rule = wtype.NewLHPolicyRule("EPAGE48Load")
	checkErr(rule.AddCategoryConditionOn("TOPLATETYPE", "EPAGE48"))
	pol = TurnOffBlowoutPolicy()
	lhpr.AddRule(rule, pol)

	//fix for removing blowout in DNA only if EGEL 48 plate type is used
	rule = wtype.NewLHPolicyRule("EGEL48Load")
	checkErr(rule.AddCategoryConditionOn("TOPLATETYPE", "EGEL48"))
	pol = TurnOffBlowoutPolicy()
	lhpr.AddRule(rule, pol)

	//fix for removing blowout in DNA only if EGEL 96_1 plate type is used
	rule = wtype.NewLHPolicyRule("EGEL961Load")
	checkErr(rule.AddCategoryConditionOn("TOPLATETYPE", "EGEL96_1"))
	pol = TurnOffBlowoutPolicy()
	lhpr.AddRule(rule, pol)

	//fix for removing blowout in DNA only if EGEL 96_2 plate type is used
	rule = wtype.NewLHPolicyRule("EGEL962Load")
	checkErr(rule.AddCategoryConditionOn("TOPLATETYPE", "EGEL96_2"))
	pol = TurnOffBlowoutPolicy()

	lhpr.AddRule(rule, pol)

	return lhpr, nil

}

func LoadLHPoliciesFromFile() (*wtype.LHPolicyRuleSet, error) {
	lhPoliciesFileName := os.Getenv("ANTHA_LHPOLICIES_FILE")
	if lhPoliciesFileName == "" {
		return nil, fmt.Errorf("Env variable ANTHA_LHPOLICIES_FILE not set")
	}
	contents, err := ioutil.ReadFile(lhPoliciesFileName)
	if err != nil {
		return nil, err
	}
	lhprs := wtype.NewLHPolicyRuleSet()
	lhprs.Policies = make(map[string]wtype.LHPolicy)
	lhprs.Rules = make(map[string]wtype.LHPolicyRule)
	//	err = readYAML(contents, lhprs)
	err = readJSON(contents, lhprs)
	if err != nil {
		return nil, err
	}
	return lhprs, nil
}

func readJSON(fileContents []byte, ruleSet *wtype.LHPolicyRuleSet) error {
	return json.Unmarshal(fileContents, ruleSet)
}
