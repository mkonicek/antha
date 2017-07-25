package liquidhandling

import (
	"context"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/microArch/driver"
)

// instructions to deal with robot setup

type RemoveAllPlatesInstruction struct {
	Type int
}

func NewRemoveAllPlatesInstruction() *RemoveAllPlatesInstruction {
	rapi := RemoveAllPlatesInstruction{Type: RAP}
	return &rapi
}

func (rapi *RemoveAllPlatesInstruction) InstructionType() int {
	return RAP
}
func (rapi *RemoveAllPlatesInstruction) Generate(ctx context.Context, policy *wtype.LHPolicyRuleSet, prms *LHProperties) ([]RobotInstruction, error) {
	return []RobotInstruction{}, nil
}
func (rapi *RemoveAllPlatesInstruction) GetParameter(name string) interface{} {
	return nil
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
	Type     int
	Position string
	Name     string
	Plate    interface{}
}

func NewAddPlateToInstruction(position, name string, plate interface{}) *AddPlateToInstruction {
	ins := AddPlateToInstruction{Type: APT, Position: position, Name: name, Plate: plate}
	return &ins
}

func (apti *AddPlateToInstruction) InstructionType() int {
	return APT
}

func (apti *AddPlateToInstruction) Generate(ctx context.Context, policy *wtype.LHPolicyRuleSet, prms *LHProperties) ([]RobotInstruction, error) {
	return []RobotInstruction{}, nil
}
func (apti *AddPlateToInstruction) GetParameter(name string) interface{} {
	switch name {
	case "POSITION":
		return apti.Position
	case "NAME":
		return apti.Name
	case "PLATE":
		return apti.Plate
	case "INSTRUCTIONTYPE":
		return apti.InstructionType()
	}
	return nil
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
