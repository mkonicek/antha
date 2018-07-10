package liquidhandling

import (
	"context"
	"fmt"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
)

type MessageInstruction struct {
	GenericRobotInstruction
	Type        int
	Message     string
	PassThrough map[string]*wtype.LHComponent
}

func NewMessageInstruction(lhi *wtype.LHInstruction) *MessageInstruction {
	msi := MessageInstruction{}
	msi.Type = MSG
	if lhi != nil {
		msi.Message = lhi.Message
		msi.PassThrough = lhi.PassThrough
	}

	return &msi
}

func (msi *MessageInstruction) Generate(ctx context.Context, policy *wtype.LHPolicyRuleSet, prms *LHProperties) ([]RobotInstruction, error) {
	// use side effect to keep IDs straight

	prms.UpdateComponentIDs(msi.PassThrough)
	return nil, nil
}

func (msi *MessageInstruction) GetParameter(name string) interface{} {
	if name == "MESSAGE" {
		return msi.Message
	}
	return nil
}

func (msi *MessageInstruction) InstructionType() int {
	return msi.Type
}

func (msi *MessageInstruction) OutputTo(driver LiquidhandlingDriver) error {
	//level int, title, text string, showcancel bool

	if msi.Message != wtype.MAGICBARRIERPROMPTSTRING {
		ret := driver.Message(0, "", msi.Message, false)
		if !ret.OK {
			return fmt.Errorf(" %d : %s", ret.Errorcode, ret.Msg)
		}
	}
	return nil
}
