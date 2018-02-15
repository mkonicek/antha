package wtype

import (
	"fmt"
)

// Type of reading to make on the plate-reader
const (
	ABSORBANCE = iota
)

//  high-level instruction to a plate reader
// to measure a sample
type PRInstruction struct {
	ID                 string
	ComponentIn        *LHComponent
	ComponentOut   	   *LHComponent
	Type               int  // Absorbance/Fluors
	Wavelength         int
	MoreOptions        int
}

func (ins PRInstruction) String() string {
	return fmt.Sprint("PRInstruction")
}


// privatised in favour of specific instruction constructors
func newPRInstruction() *PRInstruction {
	var pri PRInstruction
	pri.ID = GetUUID()
	return &pri
}

func NewPRAbsorbanceInstruction() *PRInstruction {
	pri := newPRInstruction()
	pri.Type = ABSORBANCE
	return pri
}