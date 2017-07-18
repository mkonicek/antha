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
	"fmt"
	"sort"

	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
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
)

var Robotinstructionnames = []string{"TFR", "TFB", "SCB", "MCB", "SCT", "MCT", "CCC", "LDT", "UDT", "RST", "CHA", "ASP", "DSP", "BLO", "PTZ", "MOV", "MRW", "LOD", "ULD", "SUK", "BLW", "SPS", "SDS", "INI", "FIN", "WAI", "LON", "LOF", "OPN", "CLS", "LAD", "UAD", "MMX", "MIX", "MSG"}

var RobotParameters = []string{"HEAD", "CHANNEL", "LIQUIDCLASS", "POSTO", "WELLFROM", "WELLTO", "REFERENCE", "VOLUME", "VOLUNT", "FROMPLATETYPE", "WELLFROMVOLUME", "POSFROM", "WELLTOVOLUME", "TOPLATETYPE", "MULTI", "WHAT", "LLF", "PLT", "TOWELLVOLUME", "OFFSETX", "OFFSETY", "OFFSETZ", "TIME", "SPEED", "MESSAGE"}

func InsToString(ins RobotInstruction) string {
	s := ""

	s += Robotinstructionnames[ins.InstructionType()] + " "

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
		}

		s += str + ": " + ss + " "
	}

	return s
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

// empty struct to hang methods on
type GenericRobotInstruction struct {
	Ins RobotInstruction
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

func GetPolicyFor(lhpr *wtype.LHPolicyRuleSet, ins RobotInstruction) wtype.LHPolicy {
	// find the set of matching rules
	rules := make([]wtype.LHPolicyRule, 0, len(lhpr.Rules))
	for _, rule := range lhpr.Rules {
		if ins.Check(rule) {
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

	//printPolicyForDebug(ins, rules, ppl)
	return ppl
}
