package instructions

import (
	"fmt"

	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
)

// Note: these are instructions at the level issued by intrinsics
// (i.e. wrapped in CommandInst (see trace.go)), and are sometimes
// called "high-level" instructions. These are ultimately the inputs
// to calls to Device.Compile. Do confuse these with the Inst
// interface in lowlevel.go - those are for the results of
// Compilation... Yes, this is a dreadful situation and needs to be
// addressed.

// An IncubateInst is a high-level command to incubate a component
type IncubateInst struct {
	// Time for which to incubate component
	Time wunit.Time
	// Temperature at which to incubate component
	Temp wunit.Temperature
	// Rate at which to shake incubator (force is device dependent)
	ShakeRate wunit.Rate
	// Radius at which ShakeRate is defined
	ShakeRadius wunit.Length

	// Time for which to pre-heat incubator
	PreTemp wunit.Temperature
	// Temperature at which to pre-heat incubator
	PreTime wunit.Time
	// Rate at which to pre-heat incubator
	PreShakeRate wunit.Rate
	// Radius at which PreShakeRate is defined
	PreShakeRadius wunit.Length
}

// A PromptInst is a high-level command to prompt a human
type PromptInst struct {
	Message string
}

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
