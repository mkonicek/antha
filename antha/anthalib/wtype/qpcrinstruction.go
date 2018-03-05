package wtype

import (
	"fmt"
)


// PRInstruction is a high-level instruction to a plate reader to measure a
// sample
type QPCRInstruction struct {
	ID                 string
	ComponentIn        *LHComponent
	ComponentOut   	   *LHComponent
	Definition         string
}

func (ins QPCRInstruction) String() string {
	return fmt.Sprint("QPCRInstruction")
}


// NewPRInstruction creates a new PRInstruction
func NewQPCRInstruction() *QPCRInstruction {
	var inst QPCRInstruction
	inst.ID = GetUUID()
	return &inst
}
