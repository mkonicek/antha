package liquidhandling

import (
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/laboratory"
	"github.com/antha-lang/antha/microArch/driver"
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

func (rapi *RemoveAllPlatesInstruction) Generate(labBuild *laboratory.LaboratoryBuilder, policy *wtype.LHPolicyRuleSet, prms *LHProperties) ([]RobotInstruction, error) {
	return []RobotInstruction{}, nil
}

func (rapi *RemoveAllPlatesInstruction) GetParameter(name InstructionParameter) interface{} {
	return rapi.BaseRobotInstruction.GetParameter(name)
}

func (rapi *RemoveAllPlatesInstruction) Check(lhpr wtype.LHPolicyRule) bool {
	return false
}

func (rapi *RemoveAllPlatesInstruction) OutputTo(drv LiquidhandlingDriver) error {
	stat := drv.RemoveAllPlates()

	if stat.Errorcode == driver.ERR {
		return wtype.LHError(wtype.LH_ERR_DRIV, stat.Msg)
	}

	return nil
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

func (apti *AddPlateToInstruction) Generate(labBuild *laboratory.LaboratoryBuilder, policy *wtype.LHPolicyRuleSet, prms *LHProperties) ([]RobotInstruction, error) {
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
	stat := drv.AddPlateTo(apti.Position, apti.Plate, apti.Name)

	if stat.Errorcode == driver.ERR {
		return wtype.LHError(wtype.LH_ERR_DRIV, stat.Msg)
	}

	return nil
}
