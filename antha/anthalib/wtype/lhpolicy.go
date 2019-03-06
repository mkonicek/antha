// anthalib/wtype/lhpolicy.go: Part of the Antha language
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
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/antha/anthalib/wutil"
)

const (
	LHP_AND int = iota
	LHP_OR
)

// LHPolicy defines parameters for a liquid handling policy.
// Valid key and value pairs are found aparam.go
type LHPolicy map[string]interface{}

// Name returns the name of an LHPolicy as a string
func (policy LHPolicy) Name() string {
	return policy[PolicyNameField].(string)
}

// SetName sets the name of an LHPolicy.
func (policy *LHPolicy) SetName(name string) error {
	return policy.Set(PolicyNameField, name)
}

// NewLHPolicy generates an empty LHPolicy
func NewLHPolicy() LHPolicy {
	pol := make(LHPolicy)
	return pol
}

// EquivalentPolicies checks for equality of two policies.
// We're being conservative here.
// It's possible that at the point at which low level instructions
// are generated that two policies of different length will be
// actioned in exactly the same way.
// Since we cannot guarantee this at this point, we'll say they're not equivalent.
func EquivalentPolicies(policy1, policy2 LHPolicy) bool {
	if len(policy1) != len(policy2) {
		return false
	}

	for key1, value := range policy1 {
		if !reflect.DeepEqual(policy2[key1], value) {
			return false
		}
	}
	return true
}

func (plhp *LHPolicy) Set(item string, value interface{}) error {
	var err error
	alhpis := MakePolicyItems()

	alhpi, ok := alhpis[item]

	if !ok {
		err = fmt.Errorf("No such LHPolicy item %s", item)
	} else {
		if reflect.TypeOf(value) != alhpi.Type {
			err = fmt.Errorf("LHPolicy item %s needs value of type %v not %v", item, alhpi.Type, reflect.TypeOf(value))
		} else {
			(*plhp)[item] = value
		}
	}

	return err
}

func (plhp *LHPolicy) UnmarshalJSON(data []byte) error {
	m := make(map[string]interface{})
	*plhp = make(map[string]interface{})
	lhp := *plhp
	items := MakePolicyItems()

	err := json.Unmarshal(data, &m)

	if err != nil {
		return err
	}

	for k, v := range m {
		alhpi, ok := items[k]

		if !ok {
			return fmt.Errorf("Policy item %s unknown", k)
		}

		switch alhpi.Type.Name() {
		case "float64":
			tv, ok := v.(float64)
			if !ok {
				return fmt.Errorf("Wrong type for %s: should be %s got %s", k, alhpi.Type.Name(), reflect.TypeOf(v))
			}
			lhp[k] = tv
		case "int":
			tv, ok := v.(int)

			if !ok {
				tv2, ok2 := v.(float64)
				if ok2 {
					tv = int(tv2)
				} else {
					return fmt.Errorf("Wrong type for %s: should be %s got %s", k, alhpi.Type.Name(), reflect.TypeOf(v))
				}
			}
			lhp[k] = tv
		case "string":
			tv, ok := v.(string)
			if !ok {
				return fmt.Errorf("Wrong type for %s: should be %s got %s", k, alhpi.Type.Name(), reflect.TypeOf(v))
			}
			lhp[k] = tv
		case "bool":
			tv, ok := v.(bool)
			if !ok {
				return fmt.Errorf("Wrong type for %s: should be %s got %s", k, alhpi.Type.Name(), reflect.TypeOf(v))
			}
			lhp[k] = tv
		case "wunit.Volume":
			tv, ok := v.(wunit.Volume)
			if !ok {
				return fmt.Errorf("Wrong type for %s: should be %s got %s", k, alhpi.Type.Name(), reflect.TypeOf(v))
			}
			lhp[k] = tv
		}
	}

	return nil
}

func (lhp LHPolicy) IsEqualTo(lh2 LHPolicy) bool {
	for k, v := range lhp {
		v2 := lh2[k]

		if v2 != v {
			return false
		}

	}
	return true
}

func DupLHPolicy(in LHPolicy) LHPolicy {
	ret := make(LHPolicy, len(in))

	for k, v := range in {
		ret[k] = v
	}

	return ret
}

// clobber everything in here with the other policy
// then return the merged copy
func (lhp LHPolicy) MergeWith(other LHPolicy) LHPolicy {
	for k, v := range other {
		lhp[k] = v
	}
	return lhp
}

// conditions are ANDed together
// there is no chaining
type LHPolicyRule struct {
	Name       string
	Conditions []LHVariableCondition
	Priority   int
	Type       int // AND =0 OR = 1
}

func NewLHPolicyRule(name string) LHPolicyRule {
	var lhpr LHPolicyRule
	lhpr.Name = name
	lhpr.Conditions = make([]LHVariableCondition, 0, 5)
	return lhpr
}

func (lhpr *LHPolicyRule) AddNumericConditionOn(variable string, low, up float64) error {
	params := MakeInstructionParameters()

	t, ok := params[variable]

	if !ok {
		return fmt.Errorf("No instruction defines parameter %s", variable)
	}

	if t.Type != reflect.TypeOf(3.2) {
		return fmt.Errorf("Parameter %s is not numeric", variable)
	}

	lhvc := NewLHVariableCondition(variable)
	err := lhvc.SetNumeric(low, up)

	if err != nil {
		return err
	}
	lhpr.Conditions = append(lhpr.Conditions, lhvc)
	lhpr.Priority += 1
	return nil
}

func (lhpr *LHPolicyRule) AddCategoryConditionOn(variable, category string) error {
	params := MakeInstructionParameters()

	t, ok := params[variable]

	if !ok {
		return fmt.Errorf("No instruction defines parameter %s", variable)
	}

	if t.Type != reflect.TypeOf("thisisastring") {
		return fmt.Errorf("Parameter %s is not categoric", variable)
	}
	lhvc := NewLHVariableCondition(variable)
	err := lhvc.SetCategoric(category)

	if err != nil {
		return err
	}

	lhpr.Conditions = append(lhpr.Conditions, lhvc)

	lhpr.Priority += 1
	return err
}

// this just looks for the same conditions, doesn't matter if
// the rules lead to different outcomes...
// not sure if this quite gives us the right behaviour but let's
// plough on for now
func (lhpr LHPolicyRule) IsEqualTo(other LHPolicyRule) bool {
	// cannot be equal if the number of conditions is not equal
	// well we *could have this situation
	//	A: [a,b] B: [c,d] C: [a,d]
	// where rule 1 has both A and B and rule 2 only C but all
	// three have the same consequences but we'll just have to
	// try and enforce some consistency rules to prevent that situation

	if len(lhpr.Conditions) != len(other.Conditions) {
		return false
	}

	// now we have to go through - these are not ordered so there's
	// no general way to find out if the two sets are identical

	for _, c := range lhpr.Conditions {
		if !other.HasCondition(c) {
			return false
		}
	}
	return true
}

func (lhpr LHPolicyRule) HasCondition(cond LHVariableCondition) bool {
	for _, c := range lhpr.Conditions {
		if c.IsEqualTo(cond) {
			return true
		}
	}
	return false
}

type LHVariableCondition struct {
	TestVariable string
	Condition    LHCondition
}

func (lh *LHVariableCondition) UnmarshalJSON(data []byte) error {
	var dest interface{}
	err := json.Unmarshal(data, &dest)
	if err != nil {
		return err
	}
	switch t := dest.(type) {
	case map[string]interface{}:
		if v, ex := t["TestVariable"]; ex {
			if tv, nope := v.(string); !nope {
				return fmt.Errorf("Could not parse json for LHVariableCondition")
			} else {
				lh.TestVariable = tv
			}
		} else {
			return fmt.Errorf("Could not find TestVariable when unmarshaling LHVariableCondition")
		}
		//Try now with the condition
		if v, ex := t["Condition"]; ex {
			// manual ftw

			mp := v.(map[string]interface{})

			_, cat := mp["Category"]
			_, num := mp["Upper"]
			if cat {
				lhcc := LHCategoryCondition{Category: mp["Category"].(string)}
				lh.Condition = lhcc
			} else if num {
				lhnc := LHNumericCondition{Upper: mp["Upper"].(float64), Lower: mp["Lower"].(float64)}
				lh.Condition = lhnc
			} else {
				return fmt.Errorf("No Suitable Condition Format could be found")
			}

		} else {
			return fmt.Errorf("Could not find Condition when unmarshaling LHVariableCondition")
		}
	default:
		return fmt.Errorf("Could not parse json for LHVariableCondition")
	}
	return nil
}

func NewLHVariableCondition(testvariable string) LHVariableCondition {
	var lhvc LHVariableCondition
	lhvc.TestVariable = testvariable
	return lhvc
}

func (lhvc *LHVariableCondition) SetNumeric(low, up float64) error {
	if low > up {
		return LHError(LH_ERR_POLICY, fmt.Sprintf("Numeric condition requested with lower limit (%f) greater than upper limit (%f), which is not allowed", low, up))
	}
	lhvc.Condition = LHNumericCondition{up, low}
	return nil
}

func (lhvc *LHVariableCondition) SetCategoric(category string) error {
	if category == "" {
		return LHError(LH_ERR_POLICY, fmt.Sprintf("Categoric condition %s has an empty category, which is not allowed", category))
	}
	lhvc.Condition = LHCategoryCondition{category}
	return nil
}

func (lhvc LHVariableCondition) IsEqualTo(other LHVariableCondition) bool {
	if lhvc.TestVariable != other.TestVariable {
		return false
	}
	return lhvc.Condition.IsEqualTo(other.Condition)
}

type LHPolicyRuleSet struct {
	Policies map[string]LHPolicy
	Rules    map[string]LHPolicyRule
	Options  map[string]interface{}
}

func (lhpr *LHPolicyRuleSet) SetOption(optname string, value interface{}) error {
	var err error
	opts := GetLHPolicyOptions()

	// opt is of type aParam, which defines what
	// the parameter means and what type it has
	opt, ok := opts[optname]

	if !ok {
		err = fmt.Errorf("No such LHPolicy option %s", optname)
	} else {
		if reflect.TypeOf(value) != opt.Type {
			err = fmt.Errorf("LHPolicy option %s needs value of type %s not %T", optname, opt.Type.Name(), value)
		} else {
			lhpr.Options[optname] = value
		}
	}

	return err

}

func (lhpr *LHPolicyRuleSet) IsEqualTo(lhp2 *LHPolicyRuleSet) bool {
	if len(lhpr.Policies) != len(lhp2.Policies) {
		return false
	}

	for name := range lhpr.Rules {
		p1, ok1 := lhpr.Policies[name]
		p2, ok2 := lhp2.Policies[name]

		if !(ok1 && ok2) {
			return false
		}

		r1, ok1 := lhpr.Rules[name]
		r2, ok2 := lhp2.Rules[name]

		if !(ok1 && ok2) {
			return false
		}

		if !(p1.IsEqualTo(p2) && r1.IsEqualTo(r2)) {
			return false
		}
	}

	return true
}

func NewLHPolicyRuleSet() *LHPolicyRuleSet {
	var lhpr LHPolicyRuleSet
	lhpr.Policies = make(map[string]LHPolicy)
	lhpr.Rules = make(map[string]LHPolicyRule)
	lhpr.Options = make(map[string]interface{})
	return &lhpr
}

func (lhpr *LHPolicyRuleSet) AddRule(rule LHPolicyRule, consequent LHPolicy) {
	lhpr.Policies[rule.Name] = consequent
	lhpr.Rules[rule.Name] = rule
}

func CloneLHPolicyRuleSet(parent *LHPolicyRuleSet) *LHPolicyRuleSet {
	child := NewLHPolicyRuleSet()
	for k := range parent.Rules {
		child.Policies[k] = parent.Policies[k]
		child.Rules[k] = parent.Rules[k]
	}
	for k := range parent.Options {
		child.Options[k] = parent.Options[k]
	}

	return child
}

func (lhpr LHPolicyRuleSet) GetEquivalentRuleTo(rule LHPolicyRule) string {
	for k, c := range lhpr.Rules {
		if c.IsEqualTo(rule) {
			return k
		}
	}

	return ""
}

// MergeWith merges the policyToMerge with the current LHPolicyRuleSet.
// if equivalent rules are found in the policyToMerge these are given priority
// over the existing rules.
func (lhpr *LHPolicyRuleSet) MergeWith(policyToMerge *LHPolicyRuleSet) {

	for key, rule := range policyToMerge.Rules {
		name := lhpr.GetEquivalentRuleTo(rule)

		if name != "" {
			// merge the two policies
			pol := policyToMerge.Policies[key]
			p2 := lhpr.Policies[key]
			p2.MergeWith(pol)
			lhpr.Policies[key] = p2
		}
		lhpr.Rules[key] = policyToMerge.Rules[key]
		lhpr.Policies[key] = policyToMerge.Policies[key]
		lhpr.Options[key] = policyToMerge.Options[key]
	}

	// this will override existing options if they are contradictory
	for key := range policyToMerge.Options {
		lhpr.Options[key] = policyToMerge.Options[key]
	}

}

type SortableRules []LHPolicyRule

func (s SortableRules) Len() int {
	return len(s)
}

func (s SortableRules) Less(i, j int) bool {
	if s[i].Priority != s[j].Priority {
		// (numerically) highest priority wins
		return s[i].Priority < s[j].Priority
	} else if len(s[i].Conditions) != len(s[j].Conditions) {
		// most conditions wins
		return len(s[i].Conditions) < len(s[j].Conditions)
	} else {
		// longest name wins
		return len(s[i].Name) < len(s[j].Name)
	}
}

func (s SortableRules) Swap(i, j int) {
	t := s[i]
	s[i] = s[j]
	s[j] = t
}

//func (lhpr LHPolicyRuleSet) MarshalJSON() ([]byte, error) {
//	return
//}

//func (lhpr LHPolicyRuleSet) UnmarshalJSON(data []byte) error {
//	test := NewLHPolicyRuleSet()
//	if err := json.Unmarshal(data, )
//	return nil
//}

type LHCondition interface {
	Match(interface{}) bool
	Type() string
	IsEqualTo(LHCondition) bool
}

type LHCategoryCondition struct {
	Category string
}

func (lhcc LHCategoryCondition) Match(v interface{}) bool {
	//fmt.Println(fmt.Sprintln("CATEGORY MATCH ON ", lhcc.Category, " ", v))

	switch s := v.(type) {
	case string:
		if s == lhcc.Category {
			return true
		}
	case []string:
		// true iff at least one in array and all members of the array are the same category
		if len(s) == 0 || numInStringArray(s) == 0 {
			return false
		}
		for _, str := range s {
			if !lhcc.Match(str) && str != "" {
				return false
			}
		}
		return true
	case [][]string:
		if len(s) == 0 {
			return false
		}
		for _, slice := range s {
			if !lhcc.Match(slice) {
				return false
			}
		}
		return true
	}
	return false
}

func (lhcc LHCategoryCondition) Type() string {
	return "category"
}

func (lhcc LHCategoryCondition) IsEqualTo(other LHCondition) bool {
	if other.Type() != lhcc.Type() {
		return false
	}
	return other.Match(lhcc.Category)
}

type LHNumericCondition struct {
	Upper float64
	Lower float64
}

func (lhnc LHNumericCondition) Type() string {
	return "Numeric"
}

func (lhnc LHNumericCondition) IsEqualTo(other LHCondition) bool {
	if other.Type() != lhnc.Type() {
		return false
	}
	if other.(LHNumericCondition).Upper == lhnc.Upper && other.(LHNumericCondition).Lower == lhnc.Lower {
		return true
	}
	return false
}

func (lhnc LHNumericCondition) Match(v interface{}) bool {
	//fmt.Println(fmt.Sprintln("NUMERIC MATCH: ", lhnc.Lower, " ", lhnc.Upper, " ", v))
	switch f := v.(type) {
	case float64:
		if f <= lhnc.Upper && f >= lhnc.Lower {
			return true
		}
	case []float64:
		//true iff at least one value all values are within range
		// how to deal with nulls?
		if len(f) == 0 || numInFloatArray(f) == 0 {
			return false
		}

		for _, g := range f {
			if !lhnc.Match(g) && g > wutil.EPSILON_64() {
				return false
			}
		}
		return true

	case []wunit.Volume:
		//true iff all values are within range
		// these are simple rules but could need refinement
		for _, g := range f {
			if g.IsZero() {
				return true
			}
			if !lhnc.Match(g.RawValue()) {
				return false
			}
		}
		return true

	} // switch
	return false
}

func numInStringArray(a []string) int {
	c := 0
	for _, s := range a {
		if s != "" {
			c += 1
		}
	}
	return c
}

func numInFloatArray(a []float64) int {
	c := 0
	for _, f := range a {
		if f > wutil.EPSILON_64() {
			c += 1
		}
	}
	return c
}
