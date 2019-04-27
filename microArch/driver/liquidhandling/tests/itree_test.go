package tests

import (
	"reflect"
	"testing"

	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/laboratory/effects"
	"github.com/antha-lang/antha/laboratory/testlab"
	"github.com/antha-lang/antha/microArch/driver/liquidhandling"
)

type TestInstruction struct {
	Name     string
	Children []*TestInstruction
}

var _ liquidhandling.TerminalRobotInstruction = (*TestInstruction)(nil)

func (ti *TestInstruction) Type() *liquidhandling.InstructionType {
	return liquidhandling.NewInstructionType(ti.Name, "LeafInstruction")
}

func (ti *TestInstruction) GetParameter(name liquidhandling.InstructionParameter) interface{} {
	return nil
}

func (ti *TestInstruction) Generate(labEffects *effects.LaboratoryEffects, policy *wtype.LHPolicyRuleSet, prms *liquidhandling.LHProperties) ([]liquidhandling.RobotInstruction, error) {
	ret := make([]liquidhandling.RobotInstruction, 0, len(ti.Children))
	for _, ins := range ti.Children {
		ret = append(ret, ins)
	}
	return ret, nil
}

func (ti *TestInstruction) MaybeMerge(next liquidhandling.RobotInstruction) liquidhandling.RobotInstruction {
	return nil
}

func (ti *TestInstruction) Check(lhpr wtype.LHPolicyRule) bool {
	return false
}

func (ti *TestInstruction) Visit(liquidhandling.RobotInstructionVisitor) {}

func (ti *TestInstruction) OutputTo(driver liquidhandling.LiquidhandlingDriver) error {
	return nil
}

type ITreeTest struct {
	Name           string
	Instructions   []*TestInstruction
	ExpectedLeaves []string
}

func (test *ITreeTest) Run(t *testing.T) {

	root := liquidhandling.NewITree(nil)
	for _, ins := range test.Instructions {
		root.AddChild(ins)
	}
	labEffects := testlab.NewTestLabEffects(nil)
	if _, err := root.Build(labEffects, nil, nil); err != nil {
		t.Fatal(err)
	}

	leaves := root.Leaves()
	got := make([]string, 0, len(leaves))
	for _, leaf := range leaves {
		got = append(got, leaf.Type().Name)
	}

	if !reflect.DeepEqual(got, test.ExpectedLeaves) {
		t.Errorf("leaves don't match:\n e: %s\ng: %s", test.ExpectedLeaves, got)
	}

}

func TestITree(t *testing.T) {
	(&ITreeTest{
		Name: "simple test",
		Instructions: []*TestInstruction{
			{
				Name: "A",
				Children: []*TestInstruction{
					{
						Name: "THE",
					},
					{
						Name: "CAT",
					},
				},
			},
			{
				Name: "B",
				Children: []*TestInstruction{
					{
						Name: "C",
						Children: []*TestInstruction{
							{
								Name: "SAT",
							},
							{
								Name: "D",
								Children: []*TestInstruction{
									{
										Name: "ON",
									},
									{
										Name: "THE",
									},
								},
							},
						},
					},
				},
			},
			{
				Name: "MAT",
			},
		},
		ExpectedLeaves: []string{"THE", "CAT", "SAT", "ON", "THE", "MAT"},
	}).Run(t)
}
