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

package wtype

import (
	"fmt"
	"sort"
	"strings"

	"github.com/antha-lang/antha/antha/anthalib/wunit"
)

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

func MakePolicies() map[string]LHPolicy {

	pols := make(map[string]LHPolicy)

	add := func(policy LHPolicy, name string) {
		checkErr(policy.SetName(name))
		if _, found := pols[policy.Name()]; found {
			panic(fmt.Sprintf("duplicate policy (%s) added to MakePolicies", policy.Name()))
		}
		pols[policy.Name()] = policy
	}

	// what policies do we need?
	add(SmartMixPolicy(), "SmartMix")
	add(MakeWaterPolicy(), "water")
	add(MakeSingleChannelPolicy(), "SingleChannel")
	add(MakeSmartMixSingleChannelPolicy(), "SmartMixSingleChannel")
	add(MakeMultiWaterPolicy(), "multiwater")
	add(MakeCulturePolicy(), "culture")
	add(MakeCultureReusePolicy(), "culturereuse")
	add(MakeGlycerolPolicy(), "glycerol")
	add(MakeSolventPolicy(), "solvent")
	add(MakeDefaultPolicy(), "default")
	add(MakeDNAPolicy(), "dna")
	add(MakeDefaultPolicy(), "DoNotMix")
	add(MakeNeedToMixPolicy(), "NeedToMix")
	add(PreMixPolicy(), "PreMix")
	add(PostMixPolicy(), "PostMix")
	add(MegaMixPolicy(), "MegaMix")
	add(MakeViscousPolicy(), "viscous")
	add(MakePaintPolicy(), "Paint")
	add(MakeProteinPolicy(), "protein")
	add(MakeDetergentPolicy(), "detergent")
	add(MakeLoadPolicy(), "load")
	add(MakeLoadWaterPolicy(), "loadwater")
	add(MakeDispenseAboveLiquidPolicy(), "DispenseAboveLiquid")
	add(MakeDispenseAboveLiquidMultiPolicy(), "DispenseAboveLiquidMulti")
	add(MakePEGPolicy(), "peg")
	add(MakeProtoplastPolicy(), "protoplasts")
	add(MakeDNAMixPolicy(), "dna_mix")
	add(MakeDNAMixMultiPolicy(), "dna_mix_multi")
	add(MakeDNACELLSMixPolicy(), "dna_cells_mix")
	add(MakeDNACELLSMixMultiPolicy(), "dna_cells_mix_multi")
	add(MakePlateOutPolicy(), "plateout")
	add(MakeColonyPolicy(), "colony")
	add(MakeColonyMixPolicy(), "colonymix")
	add(MakeCarbonSourcePolicy(), "carbon_source")
	add(MakeNitrogenSourcePolicy(), "nitrogen_source")
	add(MakeXYOffsetTestPolicy(), "XYOffsetTest")
	return pols
}

var DefaultPolicies map[string]LHPolicy = MakePolicies()

func GetPolicyByName(policyname PolicyName) (lhpolicy LHPolicy, err error) {
	lhpolicy, policypresent := DefaultPolicies[policyname.String()]

	if !policypresent {
		validPolicies := availablePolicies()
		return LHPolicy{}, fmt.Errorf("policy %s not found in Default list. Valid options: %s", policyname, strings.Join(validPolicies, "\n"))
	}
	return lhpolicy, nil
}

// GetPolicyByType will return the default LHPolicy corresponding to a LiquidType.
// An error is returned if an invalid liquidType is specified.
func GetPolicyByType(liquidType LiquidType) (lhpolicy LHPolicy, err error) {
	typeName, err := liquidType.String()

	if err != nil {
		return
	}
	lhpolicy, policypresent := DefaultPolicies[string(typeName)]

	if !policypresent {
		validPolicies := availablePolicies()
		return LHPolicy{}, fmt.Errorf("policy %s not found in Default list. Valid options: %s", liquidType, strings.Join(validPolicies, "\n"))
	}
	return lhpolicy, nil
}

func availablePolicies() (policies []string) {

	for key := range DefaultPolicies {
		policies = append(policies, key)
	}

	sort.Strings(policies)
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
	policy := make(LHPolicy, 10)
	policy["DEFAULTPIPETTESPEED"] = 1.5
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

func MakeProtoplastPolicy() LHPolicy {
	policy := make(LHPolicy, 8)
	policy["DEFAULTPIPETTESPEED"] = 0.5
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

func MakePaintPolicy() LHPolicy {
	policy := make(LHPolicy, 14)
	policy["DSPREFERENCE"] = 0
	policy["DSPZOFFSET"] = 0.5
	policy["DEFAULTPIPETTESPEED"] = 1.5
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

func MakeDispenseAboveLiquidPolicy() LHPolicy {
	policy := make(LHPolicy, 8)
	policy["DSPREFERENCE"] = 1 // 1 indicates dispense at top of well
	//policy["ASP_WAIT"] = 1.0
	//policy["DSP_WAIT"] = 1.0
	policy["BLOWOUTVOLUME"] = 50.0
	policy["BLOWOUTVOLUMEUNIT"] = "ul"
	policy["TOUCHOFF"] = false
	policy["CAN_MULTI"] = false
	policy["DESCRIPTION"] = "Dispense solution above the liquid to facilitate tip reuse but sacrifice pipetting accuracy at low volumes. No post-mix. No multi channel"
	return policy
}
func MakeDispenseAboveLiquidMultiPolicy() LHPolicy {
	policy := make(LHPolicy, 8)
	policy["DSPREFERENCE"] = 1 // 1 indicates dispense at top of well
	//policy["ASP_WAIT"] = 1.0
	//policy["DSP_WAIT"] = 1.0
	policy["BLOWOUTVOLUME"] = 50.0
	policy["BLOWOUTVOLUMEUNIT"] = "ul"
	policy["TOUCHOFF"] = false
	policy["CAN_MULTI"] = true
	policy["DESCRIPTION"] = "Dispense solution above the liquid to facilitate tip reuse but sacrifice pipetting accuracy at low volumes. No post Mix. Allows multi-channel pipetting."
	return policy
}

func MakeColonyPolicy() LHPolicy {
	policy := make(LHPolicy, 12)
	policy["DSPREFERENCE"] = 0
	policy["DSPZOFFSET"] = 0.0
	policy["DEFAULTPIPETTESPEED"] = 3.0
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

func MakeColonyMixPolicy() LHPolicy {
	policy := MakeColonyPolicy()
	policy["POST_MIX"] = 3
	policy["DESCRIPTION"] = "Designed for colony picking but with added post-mixes. 3 post-mix and no blowout (to avoid potential cross contamination), no multichannel."
	return policy
}

func MakeXYOffsetTestPolicy() LHPolicy {
	policy := MakeColonyPolicy()
	policy["POST_MIX_X"] = 2.0
	policy["POST_MIX_Y"] = 2.0
	policy["DESCRIPTION"] = "Intended to test setting X,Y offsets for post mixing"
	return policy
}

func MakeWaterPolicy() LHPolicy {
	waterpolicy := make(LHPolicy, 6)
	waterpolicy["DSPREFERENCE"] = 0
	waterpolicy["CAN_MSA"] = true
	waterpolicy["CAN_SDD"] = true
	waterpolicy["CAN_MULTI"] = true
	waterpolicy["DSPZOFFSET"] = 1.0
	waterpolicy["BLOWOUTVOLUME"] = 50.0
	waterpolicy["DESCRIPTION"] = "Default policy designed for pipetting water, permitting multi-channel use. Includes a blowout step for added accuracy and no post-mixing."
	return waterpolicy
}

func MakeMultiWaterPolicy() LHPolicy {
	pol := MakeWaterPolicy()
	pol["DESCRIPTION"] = "Default policy designed for pipetting water, permitting multi-channel use. Includes a blowout step for added accuracy and no post-mixing."
	return pol
}

func MakeSingleChannelPolicy() LHPolicy {
	pol := MakeWaterPolicy()
	pol["CAN_MULTI"] = false
	pol["DESCRIPTION"] = "Default policy designed for pipetting water but prohibiting multi-channel use. Includes a blowout step for added accuracy and no post-mixing."
	return pol
}

func MakeSmartMixSingleChannelPolicy() LHPolicy {
	pol := SmartMixPolicy()
	pol["CAN_MULTI"] = false
	pol["DESCRIPTION"] = "3 post-mixes of the sample being transferred. Single channel only. Volume is adjusted based upon the volume of liquid in the destination well.  No tip reuse permitted."
	return pol
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func MakeCulturePolicy() LHPolicy {
	culturepolicy := make(LHPolicy, 10)
	checkErr(culturepolicy.Set("PRE_MIX", 2))
	checkErr(culturepolicy.Set("PRE_MIX_VOLUME", 19.0))
	checkErr(culturepolicy.Set("PRE_MIX_RATE", 3.74))
	checkErr(culturepolicy.Set("DEFAULTPIPETTESPEED", 2.0))
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

func MakePlateOutPolicy() LHPolicy {
	culturepolicy := make(LHPolicy, 17)
	culturepolicy["CAN_MULTI"] = true
	culturepolicy["ASP_WAIT"] = 1.0
	culturepolicy["DSP_WAIT"] = 1.0
	culturepolicy[DefaultPipetteSpeed] = 3.0
	culturepolicy["DSPZOFFSET"] = 0.0
	culturepolicy["TIP_REUSE_LIMIT"] = 7
	culturepolicy["NO_AIR_DISPENSE"] = true
	culturepolicy["TOUCHOFF"] = false
	culturepolicy["RESET_OVERRIDE"] = true
	culturepolicy["DESCRIPTION"] = "Designed for plating out cultures onto agar plates. Dispense will be performed at the well bottom and no blowout will be performed (to minimise risk of cross contamination)"
	return culturepolicy
}

func MakeCultureReusePolicy() LHPolicy {
	culturepolicy := make(LHPolicy, 10)
	checkErr(culturepolicy.Set("PRE_MIX", 2))
	checkErr(culturepolicy.Set("PRE_MIX_VOLUME", 19.0))
	checkErr(culturepolicy.Set("PRE_MIX_RATE", 3.74))
	checkErr(culturepolicy.Set("DEFAULTPIPETTESPEED", 2.0))
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

func MakeGlycerolPolicy() LHPolicy {
	glycerolpolicy := make(LHPolicy, 9)
	glycerolpolicy["DEFAULTPIPETTESPEED"] = 1.5
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

func MakeViscousPolicy() LHPolicy {
	glycerolpolicy := make(LHPolicy, 7)
	glycerolpolicy["DEFAULTPIPETTESPEED"] = 1.5
	glycerolpolicy["ASP_WAIT"] = 1.0
	glycerolpolicy["DSP_WAIT"] = 1.0
	glycerolpolicy["CAN_MULTI"] = true
	glycerolpolicy["POST_MIX"] = 3
	glycerolpolicy["DESCRIPTION"] = "Designed for viscous samples. 3 post-mixes of the volume of the sample being transferred will be performed. No tip reuse limit."
	return glycerolpolicy
}

func MakeSolventPolicy() LHPolicy {
	solventpolicy := make(LHPolicy, 5)
	checkErr(solventpolicy.Set("PRE_MIX", 3))
	checkErr(solventpolicy.Set("DSPREFERENCE", 0))
	checkErr(solventpolicy.Set("DSPZOFFSET", 0.5))
	checkErr(solventpolicy.Set("NO_AIR_DISPENSE", true))
	checkErr(solventpolicy.Set("CAN_MULTI", true))
	checkErr(solventpolicy.Set("DESCRIPTION", "Designed for handling solvents. No post-mixes are performed"))
	return solventpolicy
}

func MakeDNAPolicy() LHPolicy {
	dnapolicy := make(LHPolicy, 12)
	dnapolicy["DEFAULTPIPETTESPEED"] = 2.0
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

func MakeDNAMixPolicy() LHPolicy {
	dnapolicy := MakeDNAPolicy()
	dnapolicy["POST_MIX_VOLUME"] = 10.0
	dnapolicy["POST_MIX"] = 5
	dnapolicy["POST_MIX_Z"] = 0.5
	dnapolicy["POST_MIX_RATE"] = 3.0
	dnapolicy["CAN_MULTI"] = false
	dnapolicy["DESCRIPTION"] = "Designed for DNA samples but with 5 post-mixes of 10ul. No tip reuse is permitted, no blowout, no multichannel."
	return dnapolicy
}

func MakeDNAMixMultiPolicy() LHPolicy {
	dnapolicy := MakeDNAPolicy()
	dnapolicy["POST_MIX_VOLUME"] = 10.0
	dnapolicy["POST_MIX"] = 5
	dnapolicy["POST_MIX_Z"] = 0.5
	dnapolicy["POST_MIX_RATE"] = 3.0
	dnapolicy["CAN_MULTI"] = true
	dnapolicy["DESCRIPTION"] = "Designed for DNA samples but with 5 post-mixes of 10ul. No tip reuse is permitted, no blowout. Allows multi-channel pipetting."
	return dnapolicy
}

func MakeDNACELLSMixPolicy() LHPolicy {
	dnapolicy := MakeDNAPolicy()
	dnapolicy["POST_MIX_VOLUME"] = 20.0
	dnapolicy["POST_MIX"] = 2
	dnapolicy["POST_MIX_Z"] = 0.5
	dnapolicy["POST_MIX_RATE"] = 1.0
	dnapolicy["DESCRIPTION"] = "Designed for mixing DNA with cells. 2 gentle post-mixes are performed. No tip reuse is permitted, no blowout."
	return dnapolicy
}
func MakeDNACELLSMixMultiPolicy() LHPolicy {
	dnapolicy := MakeDNAPolicy()
	dnapolicy["POST_MIX_VOLUME"] = 20.0
	dnapolicy["POST_MIX"] = 2
	dnapolicy["POST_MIX_Z"] = 0.5
	dnapolicy["POST_MIX_RATE"] = 1.0
	dnapolicy["CAN_MULTI"] = true
	dnapolicy["DESCRIPTION"] = "Designed for mixing DNA with cells. 2 gentle post-mixes are performed. No tip reuse is permitted, no blowout. Allows multi-channel pipetting."
	return dnapolicy
}

func MakeDetergentPolicy() LHPolicy {
	detergentpolicy := make(LHPolicy, 9)
	//        detergentpolicy["POST_MIX"] = 3
	detergentpolicy["DEFAULTPIPETTESPEED"] = 1.0
	detergentpolicy["CAN_MSA"] = false
	detergentpolicy["CAN_SDD"] = false
	detergentpolicy["DSPREFERENCE"] = 0
	detergentpolicy["DSPZOFFSET"] = 0.5
	detergentpolicy["TIP_REUSE_LIMIT"] = 8
	detergentpolicy["NO_AIR_DISPENSE"] = true
	detergentpolicy["DESCRIPTION"] = "Designed for solutions containing detergents. Gentle aspiration and dispense and a tip reuse limit of 8 to reduce problem of foam build up inside the tips."
	return detergentpolicy
}
func MakeProteinPolicy() LHPolicy {
	proteinpolicy := make(LHPolicy, 12)
	proteinpolicy["POST_MIX"] = 5
	proteinpolicy["POST_MIX_VOLUME"] = 50.0
	proteinpolicy["DEFAULTPIPETTESPEED"] = 2.0
	proteinpolicy["CAN_MSA"] = false
	proteinpolicy["CAN_SDD"] = false
	proteinpolicy["DSPREFERENCE"] = 0
	proteinpolicy["DSPZOFFSET"] = 0.5
	proteinpolicy["TIP_REUSE_LIMIT"] = 0
	proteinpolicy["NO_AIR_DISPENSE"] = true
	proteinpolicy["DESCRIPTION"] = "Designed for protein solutions. Slightly gentler aspiration and dispense and a tip reuse limit of 0 to prevent risk of cross contamination. 5 post-mixes of 50ul will be performed."
	return proteinpolicy
}
func MakeLoadPolicy() LHPolicy {

	loadpolicy := make(LHPolicy, 14)
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

func MakeLoadWaterPolicy() LHPolicy {
	loadpolicy := make(LHPolicy)
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

func MakeNeedToMixPolicy() LHPolicy {
	dnapolicy := make(LHPolicy, 16)
	dnapolicy["POST_MIX"] = 3
	dnapolicy["PRE_MIX"] = 3
	dnapolicy["PRE_MIX_VOLUME"] = 20.0
	dnapolicy["CAN_MULTI"] = true
	dnapolicy["CAN_MSA"] = false
	dnapolicy["CAN_SDD"] = false
	dnapolicy["DSPREFERENCE"] = 0
	dnapolicy["TIP_REUSE_LIMIT"] = 0
	dnapolicy["DSPZOFFSET"] = 0.5
	dnapolicy["NO_AIR_DISPENSE"] = true
	dnapolicy["DESCRIPTION"] = "3 pre-mixes and 3 post-mixes of the sample being transferred.  No tip reuse permitted."
	return dnapolicy
}

func PreMixPolicy() LHPolicy {
	dnapolicy := make(LHPolicy, 12)
	//dnapolicy["POST_MIX"] = 3
	//dnapolicy[""POST_MIX_VOLUME"] = 10.0
	//dnapolicy["POST_MIX_RATE"] = 3.74
	dnapolicy["PRE_MIX"] = 3
	dnapolicy["PRE_MIX_VOLUME"] = 19.0
	dnapolicy["CAN_MULTI"] = true
	dnapolicy["CAN_MSA"] = false
	dnapolicy["CAN_SDD"] = false
	dnapolicy["DSPREFERENCE"] = 0
	dnapolicy["TIP_REUSE_LIMIT"] = 0
	dnapolicy["DSPZOFFSET"] = 0.5
	dnapolicy["NO_AIR_DISPENSE"] = true
	dnapolicy["DESCRIPTION"] = "3 pre-mixes of the sample being transferred.  No tip reuse permitted."
	return dnapolicy

}

func PostMixPolicy() LHPolicy {
	dnapolicy := make(LHPolicy, 12)
	dnapolicy["POST_MIX"] = 3
	//dnapolicy["PRE_MIX"] = 3
	//dnapolicy["PRE_MIX_VOLUME"] = 10
	//dnapolicy["PRE_MIX_RATE"] = 3.74
	dnapolicy["CAN_MULTI"] = true
	dnapolicy["CAN_MSA"] = false
	dnapolicy["CAN_SDD"] = false
	dnapolicy["DSPREFERENCE"] = 0
	dnapolicy["TIP_REUSE_LIMIT"] = 0
	dnapolicy["DSPZOFFSET"] = 0.5
	dnapolicy["NO_AIR_DISPENSE"] = true
	dnapolicy["DESCRIPTION"] = "3 post-mixes of the sample being transferred.  No tip reuse permitted."
	return dnapolicy
}

// 3 post mixes of the sample being transferred. Volume is adjusted based upon the volume of liquid in the destination well.
// No tip reuse permitted.
// Rules added to adjust post mix volume based on volume of the destination well.
// volume now capped at max for tip type (MIX_VOLUME_OVERRIDE_TIP_MAX)
func SmartMixPolicy() LHPolicy {
	policy := make(LHPolicy, 12)
	policy["POST_MIX"] = 3
	policy["POST_MIX_VOLUME"] = 19.0
	policy["CAN_MULTI"] = true
	policy["CAN_MSA"] = false
	policy["CAN_SDD"] = false
	policy["DSPREFERENCE"] = 0
	policy["TIP_REUSE_LIMIT"] = 0
	policy["DSPZOFFSET"] = 0.5
	policy["NO_AIR_DISPENSE"] = true
	policy["DESCRIPTION"] = "3 post-mixes of the sample being transferred. Volume is adjusted based upon the volume of liquid in the destination well.  No tip reuse permitted."
	policy["MIX_VOLUME_OVERRIDE_TIP_MAX"] = true
	return policy
}

func MegaMixPolicy() LHPolicy {
	dnapolicy := make(LHPolicy, 12)
	dnapolicy["POST_MIX"] = 10
	dnapolicy["CAN_MULTI"] = true
	dnapolicy["CAN_MSA"] = false
	dnapolicy["CAN_SDD"] = false
	dnapolicy["DSPREFERENCE"] = 0
	dnapolicy["TIP_REUSE_LIMIT"] = 0
	dnapolicy["DSPZOFFSET"] = 0.5
	dnapolicy["NO_AIR_DISPENSE"] = true
	dnapolicy["DESCRIPTION"] = "10 post-mixes of the sample being transferred. No tip reuse permitted."
	return dnapolicy

}

func MakeDefaultPolicy() LHPolicy {
	defaultpolicy := make(LHPolicy, 29)
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
	defaultpolicy["DEFAULTPIPETTESPEED"] = 3.7
	defaultpolicy["DEFAULTZSPEED"] = 140.0
	defaultpolicy["MANUALPTZ"] = false
	defaultpolicy["JUSTBLOWOUT"] = false
	defaultpolicy["DONT_BE_DIRTY"] = true
	defaultpolicy["POST_MIX_Z"] = 0.5
	defaultpolicy["PRE_MIX_Z"] = 0.5
	defaultpolicy["DESCRIPTION"] = "Default mix Policy. Blowout performed, no touch off, no mixing, tip reuse permitted for the same solution."

	return defaultpolicy
}

func MakeJBPolicy() LHPolicy {
	jbp := make(LHPolicy, 1)
	checkErr(jbp.Set("JUSTBLOWOUT", true))
	checkErr(jbp.Set("TOUCHOFF", true))
	return jbp
}

func MakeTOPolicy() LHPolicy {
	top := make(LHPolicy, 1)
	checkErr(top.Set("TOUCHOFF", true))
	return top
}

func MakeLVExtraPolicy() LHPolicy {
	lvep := make(LHPolicy, 2)
	checkErr(lvep.Set("EXTRA_ASP_VOLUME", wunit.NewVolume(0.5, "ul")))
	checkErr(lvep.Set("EXTRA_DISP_VOLUME", wunit.NewVolume(0.5, "ul")))
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

func TurnOffBlowoutPolicy() LHPolicy {
	loadpolicy := make(LHPolicy, 1)
	loadpolicy["RESET_OVERRIDE"] = true
	return loadpolicy
}

func MakeHVOffsetPolicy() LHPolicy {
	lvop := make(LHPolicy, 6)
	lvop["OFFSETZADJUST"] = 0.75
	return lvop
}

func AdjustPostMixVolume(mixToVol wunit.Volume) LHPolicy {
	vol := mixToVol.ConvertToString("ul")
	policy := make(LHPolicy, 1)
	policy["POST_MIX_VOLUME"] = vol
	return policy
}

func TurnOffPostMix() LHPolicy {
	policy := make(LHPolicy, 1)
	policy["POST_MIX"] = 0
	return policy
}

func TurnOffPostMixAndPermitTipReUse() LHPolicy {
	policy := make(LHPolicy, 2)
	policy["POST_MIX"] = 0
	policy["TIP_REUSE_LIMIT"] = 100
	return policy
}

func AdjustPreMixVolume(mixToVol wunit.Volume) LHPolicy {
	vol := mixToVol.ConvertToString("ul")
	policy := make(LHPolicy, 1)
	policy["PRE_MIX_VOLUME"] = vol
	return policy
}

func MakeHVFlowRatePolicy() LHPolicy {
	policy := make(LHPolicy, 1)
	policy["DEFAULTPIPETTESPEED"] = 37.0
	return policy
}

func MakeCarbonSourcePolicy() LHPolicy {
	cspolicy := make(LHPolicy, 1)
	cspolicy["DSPREFERENCE"] = 1
	cspolicy["DESCRIPTION"] = "Custom policy for carbon source which dispenses above destination solution."
	return cspolicy
}

func MakeNitrogenSourcePolicy() LHPolicy {
	nspolicy := make(LHPolicy, 1)
	nspolicy["DSPREFERENCE"] = 1
	nspolicy["DESCRIPTION"] = "Custom policy for nitrogen source which dispenses above destination solution."
	return nspolicy
}

// newConditionalRule makes a new LHPolicyRule with conditions to apply to an LHPolicy.
//
// An error is returned if an invalid Condition Class or SetPoint is specified.
// The valid Setpoints can be found in MakeInstructionParameters()
func newConditionalRule(ruleName string, conditions ...condition) (LHPolicyRule, error) {
	var errs []string

	rule := NewLHPolicyRule(ruleName)
	for _, condition := range conditions {
		err := condition.AddToRule(&rule)
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
	AddToRule(*LHPolicyRule) error
}

type categoricCondition struct {
	Class    string
	SetPoint string
}

func (c categoricCondition) AddToRule(rule *LHPolicyRule) error {
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

func (c numericCondition) AddToRule(rule *LHPolicyRule) error {
	return rule.AddNumericConditionOn(c.Class, c.Range.Lower, c.Range.Upper)
}

// Conditions to apply to LHpolicyRules based on liquid policy used
var (
	OnSmartMix  = categoricCondition{"LIQUIDCLASS", "SmartMix"}
	OnPostMix   = categoricCondition{"LIQUIDCLASS", "PostMix"}
	OnPreMix    = categoricCondition{"LIQUIDCLASS", "PreMix"}
	OnNeedToMix = categoricCondition{"LIQUIDCLASS", "NeedToMix"}
	OnMegaMix   = categoricCondition{"LIQUIDCLASS", "MegaMix"}
)

// Conditions to apply to LHpolicyRules based on volume of liquid that a sample is being pipetted into at the destination well
var (
	IntoEmpty                 = numericCondition{Class: "WELLTOVOLUME", Range: conditionRange{Lower: 0.0, Upper: 0.009}}
	IntoLessThan20ul          = numericCondition{Class: "WELLTOVOLUME", Range: conditionRange{Lower: 0.01, Upper: 20.0}}
	IntoBetween20ulAnd50ul    = numericCondition{Class: "WELLTOVOLUME", Range: conditionRange{20.1, 50.0}}
	IntoBetween50ulAnd100ul   = numericCondition{Class: "WELLTOVOLUME", Range: conditionRange{50.1, 100.0}}
	IntoBetween100ulAnd200ul  = numericCondition{Class: "WELLTOVOLUME", Range: conditionRange{100.1, 200.0}}
	IntoBetween200ulAnd1000ul = numericCondition{Class: "WELLTOVOLUME", Range: conditionRange{200.1, 1000.0}}
)

// Conditions to apply to LHpolicyRules based on volume of liquid being transferred
var (
	LessThan20ul        = numericCondition{Class: "VOLUME", Range: conditionRange{0.0, 20.0}}
	GreaterThan20ul     = numericCondition{Class: "VOLUME", Range: conditionRange{20.1, 1000.0}}
	Between20ulAnd200ul = numericCondition{Class: "VOLUME", Range: conditionRange{20.01, 200.0}}
)

// Conditions to apply to LHpolicyRules based on volume of liquid in source well from which a sample is taken
var (
	FromBetween100ulAnd200ul  = numericCondition{Class: "WELLFROMVOLUME", Range: conditionRange{100.0, 200.0}}
	FromBetween200ulAnd1000ul = numericCondition{Class: "WELLFROMVOLUME", Range: conditionRange{200.1, 1000.0}}
)

func AddUniversalRules(originalRuleSet *LHPolicyRuleSet, policies map[string]LHPolicy) (lhpr *LHPolicyRuleSet, err error) {

	lhpr = originalRuleSet

	for name, policy := range policies {

		err = policy.SetName(name)

		if err != nil {
			return nil, err
		}

		rule := NewLHPolicyRule(name)
		err := rule.AddCategoryConditionOn("LIQUIDCLASS", name)

		if err != nil {
			return nil, err
		}
		lhpr.AddRule(rule, policy)
	}

	// FINALLY it is 'p' (sadly this was not also 's' by ANY STRETCH OF THE IMAGINATION :( )
	/*
		// hack to fix plate type problems
		// this really should be removed asap
		rule := NewLHPolicyRule("HVOffsetFix")

		OnGilson := categoricCondition{"PLATFORM", "GilsonPipetmax"}

		// to fix: This offset fix is not consistent with other tip types (e.g. filter tips)
		highVolumeTips := categoricCondition{"TIPTYPE", "Gilson200"}

		hvOffsetFix, err := newConditionalRule("HVOffsetFix", OnGilson, highVolumeTips)

		if err != nil {
			return nil, err
		}
		// don't get overridden
		hvOffsetFix.Priority = 100
		pol := MakeHVOffsetPolicy()
		lhpr.AddRule(hvOffsetFix, pol)
	*/

	// unless a policy has a default speed explicitely set we'll increase to max for high volumes
	for name, policy := range policies {
		if _, found := policy[DefaultPipetteSpeed]; !found && policy.Name() != "" {
			increaseFlowRate, err := newConditionalRule("highVolumeFlowRateFix"+"_"+name, Between20ulAnd200ul, categoricCondition{LiquidClass, name})
			if err != nil {
				return nil, err
			}
			increaseFlowRate.Priority = 100
			lhpr.AddRule(increaseFlowRate, MakeHVFlowRatePolicy())
		}
	}

	rule := NewLHPolicyRule("DNALV")
	err = rule.AddNumericConditionOn("VOLUME", 0.0, 1.99)
	if err != nil {
		return nil, err
	}
	err = rule.AddCategoryConditionOn("LIQUIDCLASS", "dna")
	if err != nil {
		return nil, err
	}
	pol := MakeLVDNAMixPolicy()
	lhpr.AddRule(rule, pol)

	// don't mix if destination well is empty
	turnOffPostMixIfEmpty, err := newConditionalRule("doNotMixIfEmpty", IntoEmpty)

	if err != nil {
		return lhpr, err
	}

	lhpr.AddRule(turnOffPostMixIfEmpty, TurnOffPostMix())

	return lhpr, nil
}

// CopyRulesFromPolicy will copy all instances of rules which mention tthe liquid class policyToCopyRulesFrom into duplicate rules set on policyToCopyRulesTo.
func CopyRulesFromPolicy(ruleSet *LHPolicyRuleSet, policyToCopyRulesFrom, policyToCopyRulesTo string) (err error) {

	policyToCopyRulesFromName := policyToCopyRulesFrom

	for _, rule := range ruleSet.Rules {
		var copyThisRule bool
		newRule := NewLHPolicyRule(rule.Name + policyToCopyRulesTo)

		for _, condition := range rule.Conditions {
			if condition.TestVariable == LiquidClass && condition.Condition.Match(policyToCopyRulesFromName) {
				copyThisRule = true
				lhvc := NewLHVariableCondition(LiquidClass)
				err := lhvc.SetCategoric(policyToCopyRulesTo)

				if err != nil {
					return err
				}

				newRule.Conditions = append(newRule.Conditions, lhvc)

			} else {
				newRule.Conditions = append(newRule.Conditions, condition)
			}
		}
		if copyThisRule {
			ruleSet.AddRule(newRule, ruleSet.Policies[rule.Name])
		}
	}
	return nil
}

// GetLHPolicyForTest gets a set of Test LHPolicies for unit tests.
// This is not guaranteed to be consistent with the default system policies returned from GetSystemLHPolicies().
func GetLHPolicyForTest() (*LHPolicyRuleSet, error) {
	lhpr, err := GetSystemLHPolicies()
	if err != nil {
		return lhpr, err
	}
	// Current tests rely on water policy being single channel.g
	policy := lhpr.Policies["water"]

	err = policy.Set("CAN_MULTI", false)
	if err != nil {
		return lhpr, err
	}
	lhpr.Policies["water"] = policy
	return lhpr, err
}

// GetSystemLHPolicies is used to set the default System Policies.
func GetSystemLHPolicies() (*LHPolicyRuleSet, error) {
	// make some policies

	policies := MakePolicies()

	// now make rules

	lhpr := NewLHPolicyRuleSet()

	lhpr, err := AddUniversalRules(lhpr, policies)

	if err != nil {
		return nil, err
	}

	for name, policy := range policies {
		rule := NewLHPolicyRule(name)
		err := rule.AddCategoryConditionOn("LIQUIDCLASS", name)

		if err != nil {
			return nil, err
		}
		lhpr.AddRule(rule, policy)
	}

	// don't mix AND turn off tip reuse limit if destination well is empty and SmartMix
	turnOffPostMixAndTipReuseIfEmptySmartMix, err := newConditionalRule("doNotMixDoNotChangeTipsIfEmptySmartMix", IntoEmpty, OnSmartMix)

	if err != nil {
		return lhpr, err
	}

	lhpr.AddRule(turnOffPostMixAndTipReuseIfEmptySmartMix, TurnOffPostMixAndPermitTipReUse())

	// don't mix AND turn off tip reuse limit if destination well is empty and PostMix
	turnOffPostMixAndTipReuseIfEmptyPostMix, err := newConditionalRule("doNotMixDoNotChangeTipsIfEmptyPostMix", IntoEmpty, OnPostMix)

	if err != nil {
		return lhpr, err
	}

	lhpr.AddRule(turnOffPostMixAndTipReuseIfEmptyPostMix, TurnOffPostMixAndPermitTipReUse())

	// don't mix AND turn off tip reuse limit if destination well is empty and MegaMix
	turnOffPostMixAndTipReuseIfEmptyMegaMix, err := newConditionalRule("doNotMixDoNotChangeTipsIfEmptyMegaMix", IntoEmpty, OnMegaMix)

	if err != nil {
		return lhpr, err
	}

	lhpr.AddRule(turnOffPostMixAndTipReuseIfEmptyMegaMix, TurnOffPostMixAndPermitTipReUse())

	// don't mix AND turn off tip reuse limit if destination well is empty and NeedToMix
	turnOffPostMixAndTipReuseIfEmptyNeedToMix, err := newConditionalRule("doNotMixDoNotChangeTipsIfEmptyNeedToMix", IntoEmpty, OnNeedToMix)

	if err != nil {
		return lhpr, err
	}

	lhpr.AddRule(turnOffPostMixAndTipReuseIfEmptyNeedToMix, TurnOffPostMixAndPermitTipReUse())

	adjustPostMixLessThan20, err := newConditionalRule("mixIntoLessThan20ul", OnSmartMix, IntoLessThan20ul)

	if err != nil {
		return lhpr, err
	}

	adjustPostMix, err := newConditionalRule("mixInto20ul", OnSmartMix, IntoBetween20ulAnd50ul)

	if err != nil {
		return lhpr, err
	}

	adjustVol20 := AdjustPostMixVolume(wunit.NewVolume(20, "ul"))
	adjustVol50 := AdjustPostMixVolume(wunit.NewVolume(50, "ul"))
	adjustVol100 := AdjustPostMixVolume(wunit.NewVolume(100, "ul"))
	adjustVol200 := AdjustPostMixVolume(wunit.NewVolume(200, "ul"))

	lhpr.AddRule(adjustPostMixLessThan20, adjustVol20)

	lhpr.AddRule(adjustPostMix, adjustVol20)

	adjustPostMix50, err := newConditionalRule("mixInto50ul", OnSmartMix, IntoBetween50ulAnd100ul, GreaterThan20ul)

	if err != nil {
		return lhpr, err
	}

	lhpr.AddRule(adjustPostMix50, adjustVol50)

	adjustPostMix100, err := newConditionalRule("mixInto100ul", OnSmartMix, IntoBetween100ulAnd200ul, GreaterThan20ul)

	if err != nil {
		return lhpr, err
	}

	lhpr.AddRule(adjustPostMix100, adjustVol100)

	adjustPostMix200, err := newConditionalRule("mixInto200ul", OnSmartMix, IntoBetween200ulAnd1000ul, GreaterThan20ul)

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

	adjustPreMix100ul, err := newConditionalRule("PreMixFrom100ul", OnPreMix, FromBetween100ulAnd200ul, GreaterThan20ul)

	if err != nil {
		return lhpr, err
	}

	lhpr.AddRule(adjustPreMix100ul, adjustPreMixVol100)

	adjustPreMix200ul, err := newConditionalRule("PreMixFrom200ul", OnPreMix, FromBetween200ulAnd1000ul, GreaterThan20ul)

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

	adjustNeedToMix100ul, err := newConditionalRule("NeedToPreMixFrom100ul", OnNeedToMix, FromBetween100ulAnd200ul, GreaterThan20ul)

	if err != nil {
		return lhpr, err
	}

	lhpr.AddRule(adjustNeedToMix100ul, adjustPreMixVol100)

	adjustNeedToMix200ul, err := newConditionalRule("NeedToPreMixFrom200ul", OnNeedToMix, FromBetween200ulAnd1000ul, GreaterThan20ul)

	if err != nil {
		return lhpr, err
	}

	lhpr.AddRule(adjustNeedToMix200ul, adjustPreMixVol200)

	/*
		// hack to fix plate type problems
		// this really should be removed asap
		rule := NewLHPolicyRule("HVOffsetFix")
		//rule.AddNumericConditionOn("VOLUME", 20.1, 300.0) // what about higher? // set specifically for openPlant configuration

		checkErr(rule.AddCategoryConditionOn("TIPTYPE", "Gilson200"))
		checkErr(rule.AddCategoryConditionOn("PLATFORM", "GilsonPipetmax"))
		// don't get overridden
		rule.Priority = 100
		pol := MakeHVOffsetPolicy()
		lhpr.AddRule(rule, pol)
	*/
	rule := NewLHPolicyRule("DNALV")
	checkErr(rule.AddNumericConditionOn("VOLUME", 0.0, 1.99))
	checkErr(rule.AddCategoryConditionOn("LIQUIDCLASS", "dna"))
	pol := MakeLVDNAMixPolicy()
	lhpr.AddRule(rule, pol)

	//fix for removing blowout in DNA only if EGEL 48 plate type is used
	rule = NewLHPolicyRule("EPAGE48Load")
	checkErr(rule.AddCategoryConditionOn("TOPLATETYPE", "EPAGE48"))
	pol = TurnOffBlowoutPolicy()
	lhpr.AddRule(rule, pol)

	//fix for removing blowout in DNA only if EGEL 48 plate type is used
	rule = NewLHPolicyRule("EGEL48Load")
	checkErr(rule.AddCategoryConditionOn("TOPLATETYPE", "EGEL48"))
	pol = TurnOffBlowoutPolicy()
	lhpr.AddRule(rule, pol)

	//fix for removing blowout in DNA only if EGEL 96_1 plate type is used
	rule = NewLHPolicyRule("EGEL961Load")
	checkErr(rule.AddCategoryConditionOn("TOPLATETYPE", "EGEL96_1"))
	pol = TurnOffBlowoutPolicy()
	lhpr.AddRule(rule, pol)

	//fix for removing blowout in DNA only if EGEL 96_2 plate type is used
	rule = NewLHPolicyRule("EGEL962Load")
	checkErr(rule.AddCategoryConditionOn("TOPLATETYPE", "EGEL96_2"))
	pol = TurnOffBlowoutPolicy()

	lhpr.AddRule(rule, pol)

	err = CopyRulesFromPolicy(lhpr, "SmartMix", "SmartMixSingleChannel")
	if err != nil {
		return lhpr, err
	}
	err = CopyRulesFromPolicy(lhpr, "SmartMix", "SmartMixLiquidLevel")
	return lhpr, err
}
