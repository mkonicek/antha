package liquidhandling

import (
	"context"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
)

type MessageInstruction struct {
	BaseRobotInstruction
	*InstructionType
	Message     string
	PassThrough map[string]*wtype.Liquid
}

func NewMessageInstruction(lhi *wtype.LHInstruction) *MessageInstruction {
	msi := &MessageInstruction{
		InstructionType: MSG,
	}
	msi.BaseRobotInstruction = NewBaseRobotInstruction(msi)

	if lhi != nil {
		msi.Message = lhi.Message
		msi.PassThrough = lhi.PassThrough
	}

	return msi
}

func (ins *MessageInstruction) Visit(visitor RobotInstructionVisitor) {
	visitor.Message(ins)
}

func (msi *MessageInstruction) Generate(ctx context.Context, policy *wtype.LHPolicyRuleSet, prms *LHProperties) ([]RobotInstruction, error) {
	// use side effect to keep IDs straight

	prms.UpdateComponentIDs(msi.PassThrough)
	return nil, nil
}

func (msi *MessageInstruction) GetParameter(name InstructionParameter) interface{} {
	switch name {
	case MESSAGE:
		return msi.Message
	default:
		return msi.BaseRobotInstruction.GetParameter(name)
	}
}

func (msi *MessageInstruction) OutputTo(driver LiquidhandlingDriver) error {
	//level int, title, text string, showcancel bool

	if msi.Message != wtype.MAGICBARRIERPROMPTSTRING {
		return driver.Message(0, "", msi.Message, false).GetError()
	}
	return nil
}
