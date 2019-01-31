package liquidhandling

import (
	"github.com/pkg/errors"
	"reflect"
	"strings"
	"testing"

	"github.com/antha-lang/antha/antha/anthalib/wtype"
)

const (
	defaultZSpeed       = 140.0
	defaultZOffset      = 0.5
	defaultPipetteSpeed = 3.7
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

type NumericCondition struct { //nolint
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
	Policy     map[InstructionParameter]interface{}
}

func (self *Rule) AddToPolicy(pol *wtype.LHPolicyRuleSet) {
	rule := wtype.NewLHPolicyRule(self.Name)
	for _, c := range self.Conditions {
		c.ApplyTo(&rule)
	}

	policy := make(wtype.LHPolicy, len(self.Policy))
	for k, v := range self.Policy {
		policy[string(k)] = v
	}

	pol.AddRule(rule, policy)
}

type InstructionAssertion struct {
	Instruction int
	Values      map[InstructionParameter]interface{}
}

func (self *InstructionAssertion) Assert(t *testing.T, ris []TerminalRobotInstruction) {
	if self.Instruction < 0 || self.Instruction >= len(ris) {
		t.Errorf("test error: assertion on instruction %d, but only %d instructions", self.Instruction, len(ris))
		return
	}
	ins := ris[self.Instruction]

	for param, e := range self.Values {
		if g := ins.GetParameter(param); !reflect.DeepEqual(e, g) {
			t.Errorf("instruction %d parameter %s: expected %v, got %v", self.Instruction, param, e, g)
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
	Error                string
}

func stringInstructions(inss []TerminalRobotInstruction) string {
	s := make([]string, len(inss))
	for i, ins := range inss {
		s[i] = ins.Type().Name
	}
	return "[" + strings.Join(s, ",") + "]"
}

func (self *PolicyTest) Run(t *testing.T) {

	t.Run(self.Name, func(t *testing.T) {
		self.run(t)
	})

}

func (self *PolicyTest) run(t *testing.T) {
	ctx := GetContextForTest()

	if self.Robot == nil {
		self.Robot = MakeGilsonWithPlatesAndTipboxesForTest("")
	}

	policySet, err := wtype.GetLHPolicyForTest()
	if err != nil {
		t.Fatal(err)
	}

	for _, rule := range self.Rules {
		rule.AddToPolicy(policySet)
	}

	tree := NewITree(self.Instruction)
	if _, err := tree.Build(ctx, policySet, self.Robot); err != nil {
		if self.Error == "" {
			err = errors.Wrapf(err, "%s: unexpected error", self.Name)
			t.Error(err)
		} else if self.Error != err.Error() {
			t.Errorf("errors don't match:\ne: \"%s\",\ng: \"%s\"", self.Error, err.Error())
		}
		return
	}

	if self.Error != "" {
		t.Errorf("error not generated: expected \"%s\"", self.Error)
	} else if ris, err := tree.Leaves(); err != nil {
		t.Error(err)
	} else if g := stringInstructions(ris); self.ExpectedInstructions != g {
		t.Errorf("instruction types don't match\n  g: %s\n  e: %s", g, self.ExpectedInstructions)
	} else {
		for _, a := range self.Assertions {
			a.Assert(t, ris)
		}
	}
}
