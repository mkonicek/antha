package ast

import (
	"fmt"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
)

// QPCRInstruction is a high-level instruction to perform a QPCR analysis.
type QPCRInstruction struct {
	ID           string
	ComponentIn  []*wtype.Liquid
	ComponentOut []*wtype.Liquid
	Definition   string
	Barcode      string
	Command      string
	TagAs        string
}

func (ins QPCRInstruction) String() string {
	return fmt.Sprint("QPCRInstruction")
}

// NewQPCRInstruction creates a new QPCRInstruction
func NewQPCRInstruction() *QPCRInstruction {
	return &QPCRInstruction{
		ID: wtype.GetUUID(),
	}
}

type QPCRExpectDataInstruction struct {
	*QPCRInstruction
	ID      string
	Options []string
}

func NewQPCRExpectDataInstruction(parent *QPCRInstruction) *QPCRExpectDataInstruction {
	return &QPCRExpectDataInstruction{
		QPCRInstruction: parent,
		ID:              wtype.GetUUID(),
	}
}
