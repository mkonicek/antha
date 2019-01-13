package wtype

import (
	"fmt"

	"github.com/antha-lang/antha/laboratory/effects/id"
)

// PRInstruction is a high-level instruction to a plate reader to measure a
// sample
type PRInstruction struct {
	ID           string
	ComponentIn  *Liquid
	ComponentOut *Liquid
	Options      string
}

func (ins PRInstruction) String() string {
	return fmt.Sprint("PRInstruction")
}

// NewPRInstruction creates a new PRInstruction
func NewPRInstruction(idGen *id.IDGenerator) *PRInstruction {
	var pri PRInstruction
	pri.ID = idGen.NextID()
	return &pri
}
