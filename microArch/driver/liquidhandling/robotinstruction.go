// anthalib/driver/liquidhandling/robotinstruction.go: Part of the Antha language
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
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/antha-lang/antha/antha/anthalib/wtype"
)

type RobotInstruction interface {
	Type() *InstructionType
	GetParameter(name InstructionParameter) interface{}
	Generate(ctx context.Context, policy *wtype.LHPolicyRuleSet, prms *LHProperties) ([]RobotInstruction, error)
	MaybeMerge(next RobotInstruction) RobotInstruction
	Check(lhpr wtype.LHPolicyRule) bool
	Visit(RobotInstructionVisitor)
}

type TerminalRobotInstruction interface {
	RobotInstruction
	OutputTo(driver LiquidhandlingDriver) error
}

var (
	TFR = NewInstructionType("TFR", "Transfer")
	TFB = NewInstructionType("TFB", "TransferBlock")
	CBI = NewInstructionType("CBI", "ChannelBlock")
	CTI = NewInstructionType("CTI", "ChannelTransfer")
	CCC = NewInstructionType("CCC", "ChangeChannelCharacteristics")
	LDT = NewInstructionType("LDT", "LoadTipsMove")
	UDT = NewInstructionType("UDT", "UnloadTipsMove")
	RST = NewInstructionType("RST", "Reset")
	CHA = NewInstructionType("CHA", "ChangeAdaptor")
	ASP = NewInstructionType("ASP", "Aspirate")
	DSP = NewInstructionType("DSP", "Dispense")
	BLO = NewInstructionType("BLO", "Blowout")
	PTZ = NewInstructionType("PTZ", "ResetPistons")
	MOV = NewInstructionType("MOV", "Move")
	MRW = NewInstructionType("MRW", "MoveRaw")
	LOD = NewInstructionType("LOD", "LoadTips")
	ULD = NewInstructionType("ULD", "UnloadTips")
	SUK = NewInstructionType("SUK", "Suck")
	BLW = NewInstructionType("BLW", "Blow")
	SPS = NewInstructionType("SPS", "SetPipetteSpeed")
	SDS = NewInstructionType("SDS", "SetDriveSpeed")
	INI = NewInstructionType("INI", "Initialize")
	FIN = NewInstructionType("FIN", "Finalize")
	WAI = NewInstructionType("WAI", "Wait")
	LON = NewInstructionType("LON", "LightsOn")
	LOF = NewInstructionType("LOF", "LightsOff")
	OPN = NewInstructionType("OPN", "Open")
	CLS = NewInstructionType("CLS", "Close")
	LAD = NewInstructionType("LAD", "LoadAdaptor")
	UAD = NewInstructionType("UAD", "UnloadAdaptor")
	MMX = NewInstructionType("MMX", "MoveMix")
	MIX = NewInstructionType("MIX", "Mix")
	MSG = NewInstructionType("MSG", "Message")
	MAS = NewInstructionType("MAS", "MoveAspirate")
	MDS = NewInstructionType("MDS", "MoveDispense")
	MVM = NewInstructionType("MVM", "MoveMix")
	MBL = NewInstructionType("MBL", "MoveBlowout")
	RAP = NewInstructionType("RAP", "RemoveAllPlates")
	APT = NewInstructionType("APT", "AddPlateTo")
	SPB = NewInstructionType("SPB", "SplitBlock")
)

type InstructionType struct {
	Name      string `json:"Type"`
	HumanName string `json:"-"`
}

// This exists so that when InstructionType is embedded within other
// instructions, we can satisfy the RobotInstruction interface with a
// minimum amount of boilerplate.
func (it *InstructionType) Type() *InstructionType {
	return it
}

func (it *InstructionType) String() string {
	return it.Name
}

func NewInstructionType(machine, human string) *InstructionType {
	return &InstructionType{
		Name:      machine,
		HumanName: human,
	}
}

type InstructionParameter string

func (name InstructionParameter) String() string {
	return string(name)
}

const (
	BLOWOUT         InstructionParameter = "BLOWOUT"
	CHANNEL         InstructionParameter = "CHANNEL"
	COMPONENT       InstructionParameter = "COMPONENT"
	CYCLES          InstructionParameter = "CYCLES"
	DRIVE           InstructionParameter = "DRIVE"
	FPLATEWX        InstructionParameter = "FPLATEWX"
	FPLATEWY        InstructionParameter = "FPLATEWY"
	FROMPLATETYPE   InstructionParameter = "FROMPLATETYPE"
	HEAD            InstructionParameter = "HEAD"
	INSTRUCTIONTYPE InstructionParameter = "INSTRUCTIONTYPE"
	LIQUIDCLASS     InstructionParameter = "LIQUIDCLASS" // LIQUIDCLASS refers to the Component Type, This is currently used to look up the corresponding LHPolicy from an LHPolicyRuleSet
	LLF             InstructionParameter = "LLF"
	MESSAGE         InstructionParameter = "MESSAGE"
	MULTI           InstructionParameter = "MULTI"
	NAME            InstructionParameter = "NAME"
	NEWADAPTOR      InstructionParameter = "NEWADAPTOR"
	NEWSTATE        InstructionParameter = "NEWSTATE"
	OFFSETX         InstructionParameter = "OFFSETX"
	OFFSETY         InstructionParameter = "OFFSETY"
	OFFSETZ         InstructionParameter = "OFFSETZ"
	OLDADAPTOR      InstructionParameter = "OLDADAPTOR"
	OLDSTATE        InstructionParameter = "OLDSTATE"
	OVERSTROKE      InstructionParameter = "OVERSTROKE"
	PARAMS          InstructionParameter = "PARAMS"
	PLATE           InstructionParameter = "PLATE"
	PLATETYPE       InstructionParameter = "PLATETYPE"
	PLATFORM        InstructionParameter = "PLATFORM"
	PLT             InstructionParameter = "PLT"
	POS             InstructionParameter = "POS"
	POSFROM         InstructionParameter = "POSFROM"
	POSITION        InstructionParameter = "POSITION"
	POSTO           InstructionParameter = "POSTO"
	REFERENCE       InstructionParameter = "REFERENCE"
	SPEED           InstructionParameter = "SPEED"
	TIME            InstructionParameter = "TIME"
	TIPTYPE         InstructionParameter = "TIPTYPE"
	TOPLATETYPE     InstructionParameter = "TOPLATETYPE"
	TPLATEWX        InstructionParameter = "TPLATEWX"
	TPLATEWY        InstructionParameter = "TPLATEWY"
	VOLUME          InstructionParameter = "VOLUME"
	VOLUNT          InstructionParameter = "VOLUNT"
	WELL            InstructionParameter = "WELL"
	WELLFROM        InstructionParameter = "WELLFROM"
	WELLFROMVOLUME  InstructionParameter = "WELLFROMVOLUME"
	WELLTO          InstructionParameter = "WELLTO"
	WELLTOVOLUME    InstructionParameter = "WELLTOVOLUME" // WELLTOVOLUME refers to the volume of liquid already present in the well location for which a sample is due to be transferred to.
	WELLVOLUME      InstructionParameter = "WELLVOLUME"
	WHAT            InstructionParameter = "WHAT"
	WHICH           InstructionParameter = "WHICH" // WHICH returns the Component IDs, i.e. representing the specific instance of an LHComponent not currently implemented.
	WAIT            InstructionParameter = "WAIT"
)

func InsToString(ins RobotInstruction) string {
	if b, err := json.Marshal(ins); err != nil {
		panic(err)
	} else {
		return string(b)
	}
}

type BaseRobotInstruction struct {
	Ins RobotInstruction `json:"-"`
}

func NewBaseRobotInstruction(ins RobotInstruction) BaseRobotInstruction {
	return BaseRobotInstruction{
		Ins: ins,
	}
}

func (bri BaseRobotInstruction) Check(rule wtype.LHPolicyRule) bool {
	for _, vcondition := range rule.Conditions {
		// todo - this cast to InstructionParameter is gross, but we're
		// going to have to tidy types with LHPolicy work later on.
		v := bri.Ins.GetParameter(InstructionParameter(vcondition.TestVariable))
		vrai := vcondition.Condition.Match(v)
		if !vrai {
			return false
		}
	}
	return true
}

// fall-through implementation to simplify instructions that have no parameters
func (bri BaseRobotInstruction) GetParameter(p InstructionParameter) interface{} {
	switch p {
	case INSTRUCTIONTYPE:
		return bri.Ins.Type()
	default:
		return nil
	}
}

func (bri BaseRobotInstruction) MaybeMerge(next RobotInstruction) RobotInstruction {
	return bri.Ins
}

/*
func printPolicyForDebug(ins RobotInstruction, rules []wtype.LHPolicyRule, pol wtype.LHPolicy) {
 	fmt.Println("*****")
 	fmt.Println("Policy for instruction ", InsToString(ins))
 	fmt.Println()
 	fmt.Println("Active Rules:")
 	fmt.Println("\t Default")
 	for _, r := range rules {
 		fmt.Println("\t", r.Name)
 	}
 	fmt.Println()
 	itemset := wtype.MakePolicyItems()
 	fmt.Println("Full output")
 	for _, s := range itemset.OrderedList() {
 		if pol[s] == nil {
 			continue
 		}
 		fmt.Println("\t", s, ": ", pol[s])
 	}
 	fmt.Println("_____")

}
*/

// ErrInvalidLiquidType is returned when no matching liquid policy is found.
type ErrInvalidLiquidType struct {
	PolicyNames      []string
	ValidPolicyNames []string
}

func (err ErrInvalidLiquidType) Error() string {
	return fmt.Sprintf("invalid LiquidType specified.\nValid Liquid Policies found: \n%s \n invalid LiquidType specified in instruction: %v \n ", strings.Join(err.ValidPolicyNames, " \n"), err.PolicyNames)
}

var (
	// ErrNoMatchingRules is returned when no matching LHPolicyRules are found when evaluating a rule set against a RobotInsturction.
	ErrNoMatchingRules = errors.New("no matching rules found")
	// ErrNoLiquidType is returned when no liquid policy is found.
	ErrNoLiquidType = errors.New("no LiquidType in instruction")
)

func matchesLiquidClass(rule wtype.LHPolicyRule) (match bool) {
	if len(rule.Conditions) > 0 {
		for i := range rule.Conditions {
			if rule.Conditions[i].TestVariable == "LIQUIDCLASS" {
				return true
			}
		}
	}
	return false
}

// GetDefaultPolicy currently returns the default policy
// this REALLY should not be necessary... ever
func GetDefaultPolicy(lhpr *wtype.LHPolicyRuleSet, ins RobotInstruction) (wtype.LHPolicy, error) {
	defaultPolicy := wtype.DupLHPolicy(lhpr.Policies["default"])
	return defaultPolicy, nil
}

// GetPolicyFor will return a matching LHPolicy for a RobotInstruction.
// If a common policy cannot be found for instances of the instruction then an error will be returned.
func GetPolicyFor(lhpr *wtype.LHPolicyRuleSet, ins RobotInstruction) (wtype.LHPolicy, error) {
	// find the set of matching rules
	rules := make([]wtype.LHPolicyRule, 0, len(lhpr.Rules))
	var lhpolicyFound bool

	for _, rule := range lhpr.Rules {

		if ins.Check(rule) {
			if matchesLiquidClass(rule) {
				lhpolicyFound = true
			}
			rules = append(rules, rule)
		}
	}

	// sort rules by priority
	sort.Sort(wtype.SortableRules(rules))

	// we might prefer to just merge this in

	ppl := wtype.DupLHPolicy(lhpr.Policies["default"])

	for _, rule := range rules {
		ppl.MergeWith(lhpr.Policies[rule.Name])
	}
	if len(rules) == 0 {
		return ppl, ErrNoMatchingRules
	}

	policy := ins.GetParameter(LIQUIDCLASS)
	var invalidPolicyNames []string
	if policies, ok := policy.([]string); ok {
		for _, policy := range policies {
			if _, found := lhpr.Policies[policy]; !found && policy != "" {
				invalidPolicyNames = append(invalidPolicyNames, policy)
			}
		}

	} else if policyString, ok := policy.(string); ok {
		if _, found := lhpr.Policies[policyString]; !found && policyString != "" {
			invalidPolicyNames = append(invalidPolicyNames, policyString)
		}
	}

	if len(invalidPolicyNames) > 0 {
		var validPolicies []string
		for key := range lhpr.Policies {
			validPolicies = append(validPolicies, key)
		}

		sort.Strings(validPolicies)
		return ppl, ErrInvalidLiquidType{PolicyNames: invalidPolicyNames, ValidPolicyNames: validPolicies}
	}

	if !lhpolicyFound {
		return ppl, ErrNoLiquidType
	}
	//printPolicyForDebug(ins, rules, ppl)
	return ppl, nil
}

type SetOfRobotInstructions struct {
	RobotInstructions []RobotInstruction
}

func (sori *SetOfRobotInstructions) UnmarshalJSON(b []byte) error {
	// first stage -- find the instructions

	soj := struct {
		RobotInstructions []json.RawMessage
	}{}

	if err := json.Unmarshal(b, &soj); err != nil {
		return err
	}

	// second stage -- unpack into an array
	sori.RobotInstructions = make([]RobotInstruction, len(soj.RobotInstructions))
	for i, raw := range soj.RobotInstructions {
		tId := struct {
			Type string
		}{}
		if err := json.Unmarshal(raw, &tId); err != nil {
			return err
		}

		var ins RobotInstruction

		switch tId.Type {
		case "":
			return fmt.Errorf("Malformed instruction - no Type field field")
		case "RAP":
			ins = NewRemoveAllPlatesInstruction()
		case "APT":
			ins = NewAddPlateToInstruction("", "", nil)
		case "INI":
			ins = NewInitializeInstruction()
		case "ASP":
			ins = NewAspirateInstruction()
		case "DSP":
			ins = NewDispenseInstruction()
		case "MIX":
			ins = NewMixInstruction()
		case "SPS":
			ins = NewSetPipetteSpeedInstruction()
		case "SDS":
			ins = NewSetDriveSpeedInstruction()
		case "BLO":
			ins = NewBlowoutInstruction()
		case "LOD":
			ins = NewLoadTipsInstruction()
		case "MOV":
			ins = NewMoveInstruction()
		case "PTZ":
			ins = NewPTZInstruction()
		case "ULD":
			ins = NewUnloadTipsInstruction()
		case "MSG":
			ins = NewMessageInstruction(nil)
		case "WAI":
			ins = NewWaitInstruction()
		case "FIN":
			ins = NewFinalizeInstruction()
		default:
			return fmt.Errorf("Unknown instruction type: %s", tId.Type)
		}

		// finally unmarshal

		if err := json.Unmarshal(raw, ins); err != nil {
			return err
		}

		// add to array

		sori.RobotInstructions[i] = ins
	}

	return nil
}
