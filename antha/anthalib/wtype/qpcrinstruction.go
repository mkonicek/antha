package wtype

import (
	"fmt"
	"golang.org/x/net/bpf"
)


// PRInstruction is a high-level instruction to a plate reader to measure a
// sample
type QPCRInstruction struct {
	ID                 string
	ComponentIn        *LHComponent
	ComponentOut   	   *LHComponent
	Definition         string
	Barcode 		   string
	Command            string
}

func (ins QPCRInstruction) String() string {
	return fmt.Sprint("QPCRInstruction")
}


// NewQPCRInstruction creates a new QPCRInstruction
func NewQPCRInstruction() *QPCRInstruction {
	var inst QPCRInstruction
	inst.ID = GetUUID()
	return &inst
}
