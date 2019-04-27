package tests

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/pkg/errors"

	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/laboratory"
	"github.com/antha-lang/antha/laboratory/testlab"
	"github.com/antha-lang/antha/microArch/driver/liquidhandling"
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
	if err := rule.AddCategoryConditionOn(self.Attribute, self.Value); err != nil {
		panic(err)
	}
}

type NumericCondition struct { //nolint
	Attribute string
	Low       float64
	High      float64
}

func (self *NumericCondition) ApplyTo(rule *wtype.LHPolicyRule) {
	if err := rule.AddNumericConditionOn(self.Attribute, self.Low, self.High); err != nil {
		panic(err)
	}
}

type Rule struct {
	Name       string
	Conditions []Condition
	Policy     map[liquidhandling.InstructionParameter]interface{}
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

// InstructionAssertion make assertions about properties of the terminal instructions
// generated under a specific policy, e.g.
//   InstructionAssertion{
//      Instruction: 5,
//      Values: map[InstructionParameter]interface{}{
//      	"CYCLES": []int{5},
//      },
//   }
// asserts that the fifth terminal instruction has the property CYCLES = 5
type InstructionAssertion struct {
	Instruction int
	Values      map[liquidhandling.InstructionParameter]interface{}
}

// Assert test that the assertion is valid, call t.Error if not
func (self *InstructionAssertion) Assert(ris []liquidhandling.TerminalRobotInstruction) error {
	if self.Instruction < 0 || self.Instruction >= len(ris) {
		return fmt.Errorf("test error: assertion on instruction %d, but only %d instructions", self.Instruction, len(ris))
	}
	ins := ris[self.Instruction]

	for param, e := range self.Values {
		if g := ins.GetParameter(param); !reflect.DeepEqual(e, g) {
			return fmt.Errorf("instruction %d parameter %s: expected %v, got %v", self.Instruction, param, e, g)
		}
	}
	return nil
}

type PolicyTest struct {
	Name                 string
	Rules                []*Rule
	Setup                func(*PolicyTest, *laboratory.Laboratory) error
	Instruction          liquidhandling.RobotInstruction
	Robot                *liquidhandling.LHProperties
	ExpectedInstructions string
	Assertions           []*InstructionAssertion
	Error                string
}

func stringInstructions(inss []liquidhandling.TerminalRobotInstruction) string {
	s := make([]string, len(inss))
	for i, ins := range inss {
		s[i] = ins.Type().Name
	}
	return "[" + strings.Join(s, ",") + "]"
}

func (self *PolicyTest) Run(t *testing.T) {
	testlab.WithTestLab(t, "", &testlab.TestElementCallbacks{
		Name: self.Name,
		Steps: func(lab *laboratory.Laboratory) error {
			if self.Setup != nil {
				if err := self.Setup(self, lab); err != nil {
					return err
				}
			}
			if self.Robot == nil {
				self.Robot = MakeGilsonWithPlatesAndTipboxesForTest(lab, "")
			}

			policySet, err := wtype.GetLHPolicyForTest()
			if err != nil {
				return err
			}

			for _, rule := range self.Rules {
				rule.AddToPolicy(policySet)
			}

			tree := liquidhandling.NewITree(self.Instruction)
			if _, err := tree.Build(lab.LaboratoryEffects, policySet, self.Robot); err != nil {
				if self.Error == "" {
					return errors.Wrapf(err, "%s: unexpected error", self.Name)
				} else if self.Error != err.Error() {
					return fmt.Errorf("errors don't match:\ne: \"%s\",\ng: \"%s\"", self.Error, err.Error())
				}
				return nil
			}

			if self.Error != "" {
				return fmt.Errorf("error not generated: expected \"%s\"", self.Error)
			}

			ris := tree.Leaves()
			if g := stringInstructions(ris); self.ExpectedInstructions != g {
				return fmt.Errorf("instruction types don't match\n  g: %s\n  e: %s", g, self.ExpectedInstructions)
			} else {
				for _, a := range self.Assertions {
					if err := a.Assert(ris); err != nil {
						return err
					}
				}
			}
			return nil
		},
	})
}
