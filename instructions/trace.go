package instructions

import (
	"encoding/json"
	"sync"

	"github.com/antha-lang/antha/antha/anthalib/wtype"
)

// a CommandInst is a generic intrinsic instruction
type CommandInst struct {
	// Arguments to this command. Used to determine command dependencies.
	Args []*wtype.Liquid
	// Components created by this command. Returned back to user code
	Result  []*wtype.Liquid
	Command *Command
}

type Trace struct {
	lock         sync.Mutex
	instructions []*CommandInst
}

func NewTrace() *Trace {
	return &Trace{}
}

// Issue an instruction - this records the instruction into the trace.
func (tr *Trace) Issue(instruction *CommandInst) {
	tr.lock.Lock()
	defer tr.lock.Unlock()

	tr.instructions = append(tr.instructions, instruction)
}

// Returns a (shallow) copy (to avoid data races) of the issued
// instructions.
func (tr *Trace) Instructions() []*CommandInst {
	tr.lock.Lock()
	defer tr.lock.Unlock()

	clone := make([]*CommandInst, len(tr.instructions))
	copy(clone, tr.instructions)
	return clone
}

func (tr *Trace) MarshalJSON() ([]byte, error) {
	return json.Marshal(tr.Instructions())
}

func (tr *Trace) UnmarshalJSON(bs []byte) error {
	if string(bs) == "null" {
		return nil

	} else {
		tr.lock.Lock()
		defer tr.lock.Unlock()
		var insts []*CommandInst
		if err := json.Unmarshal(bs, &insts); err != nil {
			return err
		} else {
			tr.instructions = insts
			return nil
		}
	}
}
