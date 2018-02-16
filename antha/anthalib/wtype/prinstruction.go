package wtype

import (
	"fmt"
)


//  high-level instruction to a plate reader
// to measure a sample
type PRInstruction struct {
	ID                 string
	ComponentIn        *LHComponent
	ComponentOut   	   *LHComponent
	Options            string
}

func (ins PRInstruction) String() string {
	return fmt.Sprint("PRInstruction")
}


func NewPRInstruction() *PRInstruction {
	var pri PRInstruction
	pri.ID = GetUUID()
	return &pri
}
