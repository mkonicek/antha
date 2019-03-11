package liquidhandling

import (
	"context"
	"reflect"
	"testing"

	"github.com/antha-lang/antha/antha/anthalib/wtype"
)

type TestInstruction struct {
	Name     string
	Children []*TestInstruction
}

var _ TerminalRobotInstruction = (*TestInstruction)(nil)

func (ti *TestInstruction) Type() *InstructionType {
	return NewInstructionType(ti.Name, "LeafInstruction")
}

func (ti *TestInstruction) GetParameter(name InstructionParameter) interface{} {
	return nil
}

func (ti *TestInstruction) Generate(ctx context.Context, policy *wtype.LHPolicyRuleSet, prms *LHProperties) ([]RobotInstruction, error) {
	ret := make([]RobotInstruction, 0, len(ti.Children))
	for _, ins := range ti.Children {
		ret = append(ret, ins)
	}
	return ret, nil
}

func (ti *TestInstruction) MaybeMerge(next RobotInstruction) RobotInstruction {
	return nil
}

func (ti *TestInstruction) Check(lhpr wtype.LHPolicyRule) bool {
	return false
}

func (ti *TestInstruction) Visit(RobotInstructionVisitor) {}

func (ti *TestInstruction) OutputTo(driver LiquidhandlingDriver) error {
	return nil
}

type ITreeTest struct {
	Name           string
	Instructions   []*TestInstruction
	ExpectedLeaves []string
}

func (test *ITreeTest) Run(t *testing.T) {

	root := NewITree(nil)
	for _, ins := range test.Instructions {
		root.AddChild(ins)
	}

	if _, err := root.Build(context.Background(), nil, nil); err != nil {
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
