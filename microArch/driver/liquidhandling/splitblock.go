package liquidhandling

import (
	"context"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
)

type SplitBlockInstruction struct {
	GenericRobotInstruction
	Inss []*wtype.LHInstruction
}

func NewSplitBlockInstruction(inss []*wtype.LHInstruction) SplitBlockInstruction {
	sb := SplitBlockInstruction{}
	sb.Inss = inss
	sb.GenericRobotInstruction.Ins = RobotInstruction(&sb)
	return sb
}

func (sp SplitBlockInstruction) InstructionType() int {
	return SPB
}

func (sp SplitBlockInstruction) GetParameter(p string) interface{} {
	return nil
}

// this instruction does not generate anything
// it just modifies the components in the robot
func (sp SplitBlockInstruction) Generate(ctx context.Context, policy *wtype.LHPolicyRuleSet, robot *LHProperties) ([]RobotInstruction, error) {
	return []RobotInstruction{}, nil
}
