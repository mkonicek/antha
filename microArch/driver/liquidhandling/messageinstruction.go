package liquidhandling

import (
	"context"
	"fmt"
	"time"

	"github.com/antha-lang/antha/antha/anthalib/wtype"
	anthadriver "github.com/antha-lang/antha/microArch/driver"
)

type MessageInstruction struct {
	BaseRobotInstruction
	*InstructionType
	Message     string
	WaitTime    time.Duration
	PassThrough map[string]*wtype.Liquid
}

func NewMessageInstruction(lhi *wtype.LHInstruction) *MessageInstruction {
	msi := &MessageInstruction{
		InstructionType: MSG,
	}
	msi.BaseRobotInstruction = NewBaseRobotInstruction(msi)

	pt := make(map[string]*wtype.Liquid)

	if lhi != nil {
		for i := 0; i < len(lhi.Inputs); i++ {
			pt[lhi.Inputs[i].ID] = lhi.Outputs[i]
		}
		msi.Message = lhi.Message
		msi.WaitTime = lhi.WaitTime
		msi.PassThrough = pt
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
	case WAIT:
		return msi.WaitTime
	default:
		return msi.BaseRobotInstruction.GetParameter(name)
	}
}

// OutputTo is expected to produce output of 5 forms according to the combination of Message and WaitTime
// An intermediate messageWaitInstruction is generated first which handles the logic of how this occurs.
// The 5 forms are:
// 1. A prompt with a wait time:
//    - should generate a wait instruction followed by a message instruction
// 2. A prompt with no wait time:
//    - should just generate a message instruction
// 3. A wait time with no prompt:
//    - should generate an empty prompt instruction followed by a wait instruction
// 	  - the generation of the empty prompt is critical to correct generation of the correct wait instruction in trilution
// 4. A magic barrier prompt with no wait time:
//    - should not generate either instruction (the prompt is just there to split the workflow)
// 5. A magic barrier prompt with a wait time:
//    - should generate an error
func (msi *MessageInstruction) OutputTo(driver LiquidhandlingDriver) error {

	newMessage, err := msi.PreOutput()

	if err != nil {
		return err
	}

	// if any Wait time is set then a wait command will be run prior to Message
	// The Liquid handling driver will convert this into a WaitWithMessage command.
	var ret anthadriver.CommandStatus

	if newMessage.WaitInstruction != nil {

		lowLevelDriver, ok := driver.(LowLevelLiquidhandlingDriver)
		if !ok {
			return fmt.Errorf("Wrong instruction type for driver: need Lowlevel driver , got %T", driver)
		}
		ret = lowLevelDriver.Wait(newMessage.WaitInstruction.Time)

		if !ret.Ok() {
			return fmt.Errorf(" %d : %s", ret.ErrorCode, ret.Msg)
		}

	}

	if newMessage.MessageInstruction != nil {
		//level int, title, text string, showcancel bool
		return driver.Message(0, "", msi.Message, false).GetError()
	}

	return nil
}

type messageWaitInstruction struct {
	*MessageInstruction
	*WaitInstruction
}

func (msi *MessageInstruction) PreOutput() (messageWaitInstruction, error) {
	// if any Wait time is set then a wait command will be run prior to Message
	// The Liquid handling driver will convert this into a WaitWithMessage command.
	var intermediate messageWaitInstruction
	if msi.WaitTime > 0 {
		intermediate.WaitInstruction = NewWaitInstruction()
		intermediate.WaitInstruction.InstructionType = WAI
		// in seconds
		intermediate.WaitInstruction.Time = msi.WaitTime.Seconds()
	}

	if msi.Message == wtype.MAGICBARRIERPROMPTSTRING {
		if msi.WaitTime > 0 {
			return intermediate, fmt.Errorf("Wait times are incompatible with system message (%s). Please change to a non system message or contact the antha-lang authors to report this as a bug if you did not set this message.", wtype.MAGICBARRIERPROMPTSTRING)
		}
	} else {
		intermediate.MessageInstruction = msi
	}

	return intermediate, nil
}
