package liquidhandling

import "github.com/antha-lang/antha/antha/anthalib/wtype"

type MessageInstruction struct {
	GenericRobotInstruction
	Type    int
	Message string
}

func NewMessageInstruction(lhi *wtype.LHInstruction) *MessageInstruction {
	msi := MessageInstruction{}
	msi.Type = MSG
	msi.GenericRobotInstruction.Ins = &msi

	return &msi
}

func (msi *MessageInstruction) Generate(policy *wtype.LHPolicyRuleSet, prms *LHProperties) ([]RobotInstruction, error) {
	return nil, nil
}

func (msi *MessageInstruction) GetParameter(name string) interface{} {
	return msi.Message
}

func (msi *MessageInstruction) InstructionType() int {
	return msi.Type
}
