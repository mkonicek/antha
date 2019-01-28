package liquidhandling

import (
	"context"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
)

// instructions to deal with robot setup

type RemoveAllPlatesInstruction struct {
	BaseRobotInstruction
	*InstructionType
}

func NewRemoveAllPlatesInstruction() *RemoveAllPlatesInstruction {
	v := &RemoveAllPlatesInstruction{
		InstructionType: RAP,
	}
	v.BaseRobotInstruction = NewBaseRobotInstruction(v)
	return v
}

func (ins *RemoveAllPlatesInstruction) Visit(visitor RobotInstructionVisitor) {
	visitor.RemoveAllPlates(ins)
}

func (rapi *RemoveAllPlatesInstruction) Generate(ctx context.Context, policy *wtype.LHPolicyRuleSet, prms *LHProperties) ([]RobotInstruction, error) {
	return []RobotInstruction{}, nil
}

func (rapi *RemoveAllPlatesInstruction) GetParameter(name InstructionParameter) interface{} {
	return rapi.BaseRobotInstruction.GetParameter(name)
}

func (rapi *RemoveAllPlatesInstruction) Check(lhpr wtype.LHPolicyRule) bool {
	return false
}

func (rapi *RemoveAllPlatesInstruction) OutputTo(drv LiquidhandlingDriver) error {
	return drv.RemoveAllPlates().GetError()
}

type AddPlateToInstruction struct {
	BaseRobotInstruction
	*InstructionType
	Position string
	Name     string
	Plate    interface{}
}

func NewAddPlateToInstruction(position, name string, plate interface{}) *AddPlateToInstruction {
	v := &AddPlateToInstruction{
		InstructionType: APT,
		Position:        position,
		Name:            name,
		Plate:           plate,
	}
	v.BaseRobotInstruction = NewBaseRobotInstruction(v)
	return v
}

func (ins *AddPlateToInstruction) Visit(visitor RobotInstructionVisitor) {
	visitor.AddPlateTo(ins)
}

func (apti *AddPlateToInstruction) Generate(ctx context.Context, policy *wtype.LHPolicyRuleSet, prms *LHProperties) ([]RobotInstruction, error) {
	return []RobotInstruction{}, nil
}

func (apti *AddPlateToInstruction) GetParameter(name InstructionParameter) interface{} {
	switch name {
	case POSITION:
		return apti.Position
	case NAME:
		return apti.Name
	case PLATE:
		return apti.Plate
	default:
		return apti.BaseRobotInstruction.GetParameter(name)
	}
}

func (apti *AddPlateToInstruction) Check(lhpr wtype.LHPolicyRule) bool {
	return false
}

func (apti *AddPlateToInstruction) OutputTo(drv LiquidhandlingDriver) error {
	return drv.AddPlateTo(apti.Position, apti.Plate, apti.Name).GetError()
}
