package liquidhandling

import (
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/microArch/driver"
)

// instructions to deal with robot setup

type RemoveAllPlatesInstruction struct {
}

func NewRemoveAllPlatesInstruction() *RemoveAllPlatesInstruction {
	rapi := RemoveAllPlatesInstruction{}
	return &rapi
}

func (rapi *RemoveAllPlatesInstruction) InstructionType() int {
	return RAP
}

func (rapi *RemoveAllPlatesInstruction) GetParameter(name string) interface{} {
	return nil
}
func (rapi *RemoveAllPlatesInstruction) check(lhpr wtype.LHPolicyRule) bool {
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
	Position string
	Name     string
	Plate    interface{}
}

func NewAddPlateToInstruction(position, name string, plate interface{}) *AddPlateToInstruction {
	ins := AddPlateToInstruction{Position: position, Name: name, Plate: plate}
	return &ins
}

func (apti *AddPlateToInstruction) InstructionType() int {
	return APT
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

func (apti *AddPlateToInstruction) check(lhpr wtype.LHPolicyRule) bool {
	return false
}
func (apti *AddPlateToInstruction) OutputTo(drv LiquidhandlingDriver) error {
	stat := drv.AddPlateTo(apti.Position, apti.Plate, apti.Name)

	if stat.Errorcode == driver.ERR {
		return wtype.LHError(wtype.LH_ERR_DRIV, stat.Msg)
	}

	return nil
}
