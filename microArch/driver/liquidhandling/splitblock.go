package liquidhandling

import (
	"context"
	"fmt"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
)

type SplitBlockInstruction struct {
	BaseRobotInstruction
	*InstructionType
	Inss []*wtype.LHInstruction
}

func NewSplitBlockInstruction(inss []*wtype.LHInstruction) *SplitBlockInstruction {
	sb := &SplitBlockInstruction{
		InstructionType: SPB,
		Inss:            inss,
	}
	sb.BaseRobotInstruction = NewBaseRobotInstruction(sb)
	return sb
}

func (sb *SplitBlockInstruction) Visit(visitor RobotInstructionVisitor) {
	visitor.SplitBlock(sb)
}

// this instruction does not generate anything
// it just modifies the components in the robot
func (sp SplitBlockInstruction) Generate(ctx context.Context, policy *wtype.LHPolicyRuleSet, robot *LHProperties) ([]RobotInstruction, error) {
	// this may need more work

	for _, ins := range sp.Inss {
		if ins.Type != wtype.LHISPL {
			return []RobotInstruction{}, fmt.Errorf("Splitblock fed non-split instruction, type %s", ins.InsType())
		}

		// if Components is a sample we'll probably want to change ParentID instead
		// that may not work
		robot.UpdateComponentID(ins.Inputs[0].ID, ins.Outputs[1])

		/*
			question over whether this is needed
			if !ok {
				fmt.Printf("Warning: cannot update component ID %s to %s: Not found\n", ins.Inputs[0].ID, ins.Results[1].ID)
				//return []RobotInstruction{}, fmt.Errorf("Error updating component ID %s to %s: Not found", ins.Inputs[0].ID, ins.Results[1].ID)
			}
		*/
	}

	return []RobotInstruction{}, nil
}
