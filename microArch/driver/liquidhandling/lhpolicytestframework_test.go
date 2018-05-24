package liquidhandling

import (
	"context"
	"github.com/pkg/errors"
	"reflect"
	"strings"
	"testing"

	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/inventory/cache/plateCache"
	"github.com/antha-lang/antha/inventory/testinventory"
)

type Condition interface {
	ApplyTo(*wtype.LHPolicyRule)
}

type CategoryCondition struct {
	Attribute string
	Value     string
}

func (self *CategoryCondition) ApplyTo(rule *wtype.LHPolicyRule) {
	rule.AddCategoryConditionOn(self.Attribute, self.Value)
}

type NumericCondition struct {
	Attribute string
	Low       float64
	High      float64
}

func (self *NumericCondition) ApplyTo(rule *wtype.LHPolicyRule) {
	rule.AddNumericConditionOn(self.Attribute, self.Low, self.High)
}

type Rule struct {
	Name       string
	Conditions []Condition
	Policy     map[string]interface{}
}

func (self *Rule) AddToPolicy(pol *wtype.LHPolicyRuleSet) {
	rule := wtype.NewLHPolicyRule(self.Name)
	for _, c := range self.Conditions {
		c.ApplyTo(&rule)
	}

	policy := make(wtype.LHPolicy, len(self.Policy))
	for k, v := range self.Policy {
		policy[k] = v
	}

	pol.AddRule(rule, policy)
}

type InstructionAssertion struct {
	InstructionIndex int
	Values           map[string]interface{}
}

func (self *InstructionAssertion) Assert(t *testing.T, ris []RobotInstruction, name string) {
	ins := ris[self.InstructionIndex]

	for param, e := range self.Values {
		if g := ins.GetParameter(param); !reflect.DeepEqual(e, g) {
			t.Errorf("%s: instruction %d parameter %s: expected %v, got %v", name, self.InstructionIndex, param, e, g)
		}
	}

}

type PolicyTest struct {
	Name                 string
	Rules                []*Rule
	Instruction          RobotInstruction
	Robot                *LHProperties
	ExpectedInstructions string
	Assertions           []*InstructionAssertion
}

func stringInstructions(inss []RobotInstruction) string {
	s := make([]string, len(inss))
	for i, ins := range inss {
		s[i] = InstructionTypeName(ins)
	}
	return "[" + strings.Join(s, ",") + "]"
}

func (self *PolicyTest) Run(t *testing.T) {
	ctx := testinventory.NewContext(context.Background())
	ctx = plateCache.NewContext(ctx)
	if self.Robot == nil {
		robot, err := makeTestGilson(ctx)
		if err != nil {
			err = errors.Wrap(err, self.Name)
			t.Fatal(err)
		}
		self.Robot = robot
	}

	policySet, err := wtype.GetLHPolicyForTest()
	if err != nil {
		err = errors.Wrap(err, self.Name)
		t.Fatal(err)
	}

	for _, rule := range self.Rules {
		rule.AddToPolicy(policySet)
	}

	set := NewRobotInstructionSet(self.Instruction)
	ris, err := set.Generate(ctx, policySet, self.Robot)
	if err != nil {
		err = errors.Wrap(err, self.Name)
		t.Error(err)
		return
	}

	if g := stringInstructions(ris); self.ExpectedInstructions != g {
		t.Errorf("%s: instruction types don't match\n  g: %s\n  e: %s", self.Name, g, self.ExpectedInstructions)
		return
	}

	for _, a := range self.Assertions {
		a.Assert(t, ris, self.Name)
	}
}
