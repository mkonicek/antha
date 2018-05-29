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
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/antha/anthalib/wutil/text"
)

type RobotInstruction interface {
	InstructionType() int
	GetParameter(name string) interface{}
	Generate(ctx context.Context, policy *wtype.LHPolicyRuleSet, prms *LHProperties) ([]RobotInstruction, error)
	Check(lhpr wtype.LHPolicyRule) bool
}

type TerminalRobotInstruction interface {
	RobotInstruction
	OutputTo(driver LiquidhandlingDriver) error
}

const (
	TFR int = iota // Transfer
	TFB            // Transfer block
	SCB            // Single channel transfer block
	MCB            // Multi channel transfer block
	SCT            // Single channel transfer
	MCT            // multi channel transfer
	CCC            // ChangeChannelCharacteristics
	LDT            // Load Tips + Move
	UDT            // Unload Tips + Move
	RST            // Reset
	CHA            // ChangeAdaptor
	ASP            // Aspirate
	DSP            // Dispense
	BLO            // Blowout
	PTZ            // Reset pistons
	MOV            // Move
	MRW            // Move Raw
	LOD            // Load Tips
	ULD            // Unload Tips
	SUK            // Suck
	BLW            // Blow
	SPS            // Set Pipette Speed
	SDS            // Set Drive Speed
	INI            // Initialize
	FIN            // Finalize
	WAI            // Wait
	LON            // Lights On
	LOF            // Lights Off
	OPN            // Open
	CLS            // Close
	LAD            // Load Adaptor
	UAD            // Unload Adaptor
	MMX            // Move and Mix
	MIX            // Mix
	MSG            // Message
	MAS            // MOV ASP	-- used by tests
	MDS            // MOV DSP	    ""       ""
	MVM            // MOV MIX           ""       ""
	MBL            // MOV BLO	    ""       ""
	RAP            // RemoveAllPlates
	APT            // AddPlateTo
	RPA            // Remove Plate At
	SPB            // SplitBlock
)

func InstructionTypeName(ins RobotInstruction) string {
	return Robotinstructionnames[ins.InstructionType()]
}

var Robotinstructionnames = []string{"TFR", "TFB", "SCB", "MCB", "SCT", "MCT", "CCC", "LDT", "UDT", "RST", "CHA", "ASP", "DSP", "BLO", "PTZ", "MOV", "MRW", "LOD", "ULD", "SUK", "BLW", "SPS", "SDS", "INI", "FIN", "WAI", "LON", "LOF", "OPN", "CLS", "LAD", "UAD", "MMX", "MIX", "MSG", "MOVASP", "MOVDSP", "MOVMIX", "MOVBLO", "RAP", "RPA", "APT", "SPB"}

var RobotParameters = []string{"HEAD", "CHANNEL", "LIQUIDCLASS", "POSTO", "WELLFROM", "WELLTO", "REFERENCE", "VOLUME", "VOLUNT", "FROMPLATETYPE", "WELLFROMVOLUME", "POSFROM", "WELLTOVOLUME", "TOPLATETYPE", "MULTI", "WHAT", "LLF", "PLT", "OFFSETX", "OFFSETY", "OFFSETZ", "TIME", "SPEED", "MESSAGE", "COMPONENT"}

// option to feed into InsToString function
type printOption string

// Option to feed into InsToString function
// which prints key words of the instruction with coloured text.
// Designed for easier reading.
const colouredTerminalOutput printOption = "colouredTerminalOutput"

func ansiPrint(options ...printOption) bool {
	for _, option := range options {
		if option == colouredTerminalOutput {
			return true
		}
	}
	return false
}

func InsToString(ins RobotInstruction, ansiPrintOptions ...printOption) string {

	s := InstructionTypeName(ins) + " "

	var changeColour func(string) string

	if strings.TrimSpace(s) == "ASP" {
		changeColour = text.Green
	} else if strings.TrimSpace(s) == "DSP" {
		changeColour = text.Blue
	} else if strings.TrimSpace(s) == "MOV" {
		changeColour = text.Yellow
	} else {
		changeColour = text.White
	}
	if ansiPrint(ansiPrintOptions...) {
		s = changeColour(s)
	}
	for _, str := range RobotParameters {
		p := ins.GetParameter(str)

		if p == nil {
			continue
		}

		ss := ""

		switch p.(type) {
		case []wunit.Volume:
			if len(p.([]wunit.Volume)) == 0 {
				continue
			}
			ss = concatvolarray(p.([]wunit.Volume))

		case []string:
			if len(p.([]string)) == 0 {
				continue
			}
			ss = concatstringarray(p.([]string))
		case string:
			ss = p.(string)
		case []float64:
			if len(p.([]float64)) == 0 {
				continue
			}
			ss = concatfloatarray(p.([]float64))
		case float64:
			ss = fmt.Sprintf("%-6.4f", p.(float64))
		case []int:
			if len(p.([]int)) == 0 {
				continue
			}
			ss = concatintarray(p.([]int))
		case int:
			ss = fmt.Sprintf("%d", p.(int))
		case []bool:
			if len(p.([]bool)) == 0 {
				continue
			}
			ss = concatboolarray(p.([]bool))
		}
		if ansiPrint(ansiPrintOptions...) {
			if str == "WHAT" {
				s += str + ": " + text.Yellow(ss) + " "
			} else if str == "MULTI" {
				s += text.Blue(str+": ") + ss + " "
			} else if str == "OFFSETZ" {
				s += str + ": " + changeColour(ss) + " "
			} else if str == "TOPLATETYPE" {
				s += str + ": " + text.Cyan(ss) + " "
			} else {
				s += str + ": " + ss + " "
			}
		} else {
			if str == "WHAT" {
				s += str + ": " + ss + " "
			} else if str == "MULTI" {
				s += str + ": " + ss + " "
			} else if str == "OFFSETZ" {
				s += str + ": " + ss + " "
			} else if str == "TOPLATETYPE" {
				s += str + ": " + ss + " "
			} else {
				s += str + ": " + ss + " "
			}
		}
	}

	return s
}

func isAspirate(ins RobotInstruction) bool {

	s := InstructionTypeName(ins)

	return strings.TrimSpace(s) == "ASP"
}

func isDispense(ins RobotInstruction) bool {

	s := InstructionTypeName(ins)

	return strings.TrimSpace(s) == "DSP"
}

func isMove(ins RobotInstruction) bool {

	s := InstructionTypeName(ins)

	return strings.TrimSpace(s) == "MOV"
}

// StepSummary summarises the instruction for
// an Aspirate or Dispense instruction combined
// with the related Move instruction.
type StepSummary struct {
	Type         string // Asp or DSP
	LiquidType   string
	PlateType    string
	Multi        string
	OffsetZ      string
	WellToVolume string
	Volume       string
}

func mergeSummaries(a, b StepSummary, aspOrDsp string) (c StepSummary) {
	return StepSummary{
		Type:         aspOrDsp,
		LiquidType:   a.LiquidType + b.LiquidType,
		PlateType:    a.PlateType + b.PlateType,
		Multi:        a.Multi + b.Multi,
		OffsetZ:      a.OffsetZ + b.OffsetZ,
		WellToVolume: a.WellToVolume + b.WellToVolume,
		Volume:       a.Volume + b.Volume,
	}
}

type stepType string

// Aspirate designates a step is an aspirate step
const Aspirate stepType = "Aspirate"

// Dispense designates a step is a dispense step
const Dispense stepType = "Dispense"

// MakeAspOrDspSummary returns a summary of the key parameters involved in a Dispense or Aspirate step.
// It requires two consecutive instructions to do this, a Move instruction followed by a dispense of aspirate instruction.
// An error is returned if this is not the case.
func MakeAspOrDspSummary(moveInstruction, dspOrAspInstruction RobotInstruction) (StepSummary, error) {
	step1summary, err := summarise(moveInstruction)

	if err != nil {
		return StepSummary{}, err
	}

	step2summary, err := summarise(dspOrAspInstruction)

	if err != nil {
		return StepSummary{}, err
	}

	if !isMove(moveInstruction) {
		return StepSummary{}, fmt.Errorf("first instruction is not a move instruction: found %s", InstructionTypeName(moveInstruction))
	}

	if isAspirate(dspOrAspInstruction) {
		return mergeSummaries(step1summary, step2summary, string(Aspirate)), nil
	} else if isDispense(dspOrAspInstruction) {
		return mergeSummaries(step1summary, step2summary, string(Dispense)), nil
	}

	return StepSummary{}, fmt.Errorf("second instruction is not an aspirate or dispense: found %s", InstructionTypeName(dspOrAspInstruction))

}

func summarise(ins RobotInstruction) (StepSummary, error) {

	var summaryOfMoveOperation StepSummary

	for _, str := range RobotParameters {
		p := ins.GetParameter(str)

		if p == nil {
			continue
		}

		ss := ""

		switch p.(type) {
		case []wunit.Volume:
			if len(p.([]wunit.Volume)) == 0 {
				continue
			}
			ss = concatvolarray(p.([]wunit.Volume))

		case []string:
			if len(p.([]string)) == 0 {
				continue
			}
			ss = concatstringarray(p.([]string))
		case string:
			ss = p.(string)
		case []float64:
			if len(p.([]float64)) == 0 {
				continue
			}
			ss = concatfloatarray(p.([]float64))
		case float64:
			ss = fmt.Sprintf("%-6.4f", p.(float64))
		case []int:
			if len(p.([]int)) == 0 {
				continue
			}
			ss = concatintarray(p.([]int))
		case int:
			ss = fmt.Sprintf("%d", p.(int))
		case []bool:
			if len(p.([]bool)) == 0 {
				continue
			}
			ss = concatboolarray(p.([]bool))
		}
		if str == "WHAT" {
			summaryOfMoveOperation.LiquidType = ss
		} else if str == "MULTI" {
			summaryOfMoveOperation.Multi = ss
		} else if str == "OFFSETZ" {
			summaryOfMoveOperation.OffsetZ = ss
		} else if str == "TOPLATETYPE" {
			summaryOfMoveOperation.PlateType = ss
		} else if str == WELLTOVOLUME {
			summaryOfMoveOperation.WellToVolume = ss
		} else if str == "VOLUME" {
			summaryOfMoveOperation.Volume = ss
		}
	}

	return summaryOfMoveOperation, nil
}

func InsToString2(ins RobotInstruction) string {
	// IS THIS IT?!
	b, _ := json.Marshal(ins)
	return string(b)
}

func concatstringarray(a []string) string {
	r := ""

	for i, s := range a {
		r += s
		if i < len(a)-1 {
			r += ","
		}
	}

	return r
}

func concatvolarray(a []wunit.Volume) string {
	r := ""
	for i, s := range a {
		r += s.ToString()
		if i < len(a)-1 {
			r += ","
		}
	}

	return r

}

func concatfloatarray(a []float64) string {
	r := ""

	for i, s := range a {
		r += fmt.Sprintf("%-6.4f", s)
		if i < len(a)-1 {
			r += ","
		}
	}

	return r

}

func concatintarray(a []int) string {
	r := ""

	for i, s := range a {
		r += fmt.Sprintf("%d", s)
		if i < len(a)-1 {
			r += ","
		}
	}

	return r

}

func concatboolarray(a []bool) string {
	r := ""

	for i, s := range a {
		r += fmt.Sprintf("%t", s)
		if i < len(a)-1 {
			r += ","
		}
	}

	return r

}

// empty struct to hang methods on
type GenericRobotInstruction struct {
	Ins RobotInstruction `json:"-"`
}

func (gri GenericRobotInstruction) Check(rule wtype.LHPolicyRule) bool {
	for _, vcondition := range rule.Conditions {
		v := gri.Ins.GetParameter(vcondition.TestVariable)
		vrai := vcondition.Condition.Match(v)
		if !vrai {
			return false
		}
	}
	return true
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

	policy := ins.GetParameter("LIQUIDCLASS")
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

func HasParameter(s string, ins RobotInstruction) bool {
	return ins.GetParameter(s) != nil
}

type SetOfRobotInstructions struct {
	RobotInstructions []RobotInstruction
}

func (sori *SetOfRobotInstructions) UnmarshalJSON(b []byte) error {
	// first stage -- find the instructions

	var objectMap map[string]*json.RawMessage

	err := json.Unmarshal(b, &objectMap)

	if err != nil {
		return err
	}

	// second stage -- unpack into an array

	var arrI []*json.RawMessage
	mess := objectMap["RobotInstructions"]
	err = json.Unmarshal(*mess, &arrI)

	if err != nil {
		return err
	}

	sori.RobotInstructions = make([]RobotInstruction, len(arrI))
	mapForTypeCheck := make(map[string]interface{}, 10)
	for i := 0; i < len(arrI); i++ {
		mess := arrI[i]
		err = json.Unmarshal(*mess, &mapForTypeCheck)

		if err != nil {
			return err
		}

		_, ok := mapForTypeCheck["Type"]

		if !ok {
			return fmt.Errorf("Malformed instruction")
		}

		tf64, ok := mapForTypeCheck["Type"].(float64)

		if !ok {
			return fmt.Errorf("Malformed instruction - Type field must be numeric, got %T", mapForTypeCheck["Type"])
		}

		//motherofallswitches ugh

		t := int(tf64)

		var ins RobotInstruction

		switch t {
		case RAP:
			ins = NewRemoveAllPlatesInstruction()
		case APT:
			ins = NewAddPlateToInstruction("", "", nil)
		case INI:
			ins = NewInitializeInstruction()
		case ASP:
			ins = NewAspirateInstruction()
		case DSP:
			ins = NewDispenseInstruction()
		case MIX:
			ins = NewMixInstruction()
		case SPS:
			ins = NewSetPipetteSpeedInstruction()
		case SDS:
			ins = NewSetDriveSpeedInstruction()
		case BLO:
			ins = NewBlowoutInstruction()
		case LOD:
			ins = NewLoadTipsInstruction()
		case MOV:
			ins = NewMoveInstruction()
		case PTZ:
			ins = NewPTZInstruction()
		case ULD:
			ins = NewUnloadTipsInstruction()
		case MSG:
			ins = NewMessageInstruction(nil)
		case WAI:
			ins = NewWaitInstruction()
		case FIN:
			ins = NewFinalizeInstruction()
		default:
			return fmt.Errorf("Unknown instruction type: %d (%s)", t, Robotinstructionnames[t])
		}

		// finally unmarshal

		err = json.Unmarshal(*mess, &ins)

		if err != nil {
			return err
		}

		// add to array

		sori.RobotInstructions[i] = ins
	}

	return nil
}
