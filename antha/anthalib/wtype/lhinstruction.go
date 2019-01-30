package wtype

import (
	"fmt"
	"strings"

	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/laboratory/effects/id"
)

// enum of instruction types

const (
	LHIMIX = iota
	LHIPRM
	LHISPL
)

var InsNames = []string{"MIX", "PROMPT", "SPLIT"}

func InsType(i int) string {

	ret := ""

	if i >= 0 && i < len(InsNames) {
		ret = InsNames[i]
	}

	return ret
}

//  high-level instruction to a liquid handler
type LHInstruction struct {
	ID               string
	BlockID          BlockID
	Inputs           []*Liquid
	Outputs          []*Liquid
	ContainerType    string
	Welladdress      string
	PlateID          string
	Platetype        PlateTypeName
	Vol              float64
	Type             int
	Conc             float64
	Tvol             float64
	Majorlayoutgroup int
	gen              int
	PlateName        string
	OutPlate         *Plate
	Message          string
	PassThrough      map[string]*Liquid // 1:1 pass through, only applies to prompts
}

func (ins LHInstruction) String() string {
	ret := fmt.Sprintf(
		"%s G: %d %s %v %s ID(%s) %s: %s",
		ins.InsType(),
		ins.Generation(),
		ins.ID,
		ComponentVector(ins.Inputs),
		ins.PlateName,
		ins.PlateID,
		ins.Welladdress,
		ComponentVector(ins.Outputs),
	)

	if ins.IsMixInPlace() {
		ret += " INPLACE"
	}

	return ret
}

//Summarize get a string summary of the instruction for end users
func (ins *LHInstruction) Summarize(indent int) string {
	indentStr := ""
	for i := 0; i < indent; i++ {
		indentStr += "  "
	}

	var lines []string
	lines = append(lines, fmt.Sprintf("%s%s (ID:%s)", indentStr, ins.InsType(), ins.ID))

	switch ins.Type {
	case LHIMIX:
		for _, cmp := range ins.Inputs {
			lines = append(lines, fmt.Sprintf("%s  %s", indentStr, cmp.Summarize()))
		}
		if ins.Welladdress != "" && ins.OutPlate != nil {
			lines = append(lines, fmt.Sprintf(indentStr+"to well %s in plate %s", ins.Welladdress, ins.OutPlate.Name()))
		} else if ins.Platetype != "" {
			lines = append(lines, fmt.Sprintf(indentStr+"to plate of type %s", ins.Platetype))
		}

		if len(ins.Outputs) > 0 {
			lines = append(lines, fmt.Sprintf(indentStr+"Resulting volume: %v", ins.Outputs[0].Summarize()))
		}

	default:
		return indentStr + ins.String()
	}

	return strings.Join(lines, "\n")
}

// privatised in favour of specific instruction constructors
func newLHInstruction(idGen *id.IDGenerator) *LHInstruction {
	var lhi LHInstruction
	lhi.ID = idGen.NextID()
	lhi.Majorlayoutgroup = -1
	lhi.PassThrough = make(map[string]*Liquid, 1)
	return &lhi
}

func NewLHMixInstruction(idGen *id.IDGenerator) *LHInstruction {
	lhi := newLHInstruction(idGen)
	lhi.Type = LHIMIX
	return lhi

}

func NewLHPromptInstruction(idGen *id.IDGenerator) *LHInstruction {
	lhi := newLHInstruction(idGen)
	lhi.Type = LHIPRM
	return lhi
}

func NewLHSplitInstruction(idGen *id.IDGenerator) *LHInstruction {
	lhi := newLHInstruction(idGen)
	lhi.Type = LHISPL
	return lhi
}

func (inst *LHInstruction) InsType() string {
	return InsType(inst.Type)
}

// GetID returns the ID of the instruction, useful for interfaces
func (inst *LHInstruction) GetID() string {
	return inst.ID
}

func (inst *LHInstruction) AddOutput(cmp *Liquid) {
	inst.Outputs = append(inst.Outputs, cmp)
}

func (inst *LHInstruction) AddInput(cmp *Liquid) {
	if inst == nil {
		return
	}

	inst.Inputs = append(inst.Inputs, cmp)
}

func (ins *LHInstruction) Generation() int {
	return ins.gen
}

func (ins *LHInstruction) SetGeneration(i int) {
	ins.gen = i
}

func (ins *LHInstruction) SetPlateID(pid string) {
	ins.PlateID = pid
}

func (ins *LHInstruction) IsMixInPlace() bool {
	if ins == nil {
		return false
	}

	if len(ins.Inputs) == 0 {
		return false
	}

	smp := ins.Inputs[0].IsSample()
	return !smp
}

//IsDummy return true if the instruction has no effect
func (ins *LHInstruction) IsDummy() bool {
	if ins.Type == LHIMIX && ins.IsMixInPlace() && len(ins.Inputs) == 1 {
		// instructions of this form generally mean "do nothing"
		// but they have the effect of ensuring that the compoenent ID is changed
		return true
	}

	return false
}

func (ins *LHInstruction) HasAnyParent() bool {
	for _, v := range ins.Inputs {
		if v.HasAnyParent() {
			return true
		}
	}

	return false
}

func (ins *LHInstruction) HasParent(id string) bool {
	for _, v := range ins.Inputs {
		if v.HasParent(id) {
			return true
		}
	}
	return false
}

func (ins *LHInstruction) ParentString() string {
	if ins == nil {
		return ""
	}

	tx := make([]string, 0, 1)

	for _, v := range ins.Inputs {
		//s += v.ParentID + "_"

		pid := v.ParentID

		if pid != "" {
			tx = append(tx, pid)
		}
	}

	if len(tx) == 0 {
		return ""
	} else {
		return strings.Join(tx, "_")
	}

}

func (ins *LHInstruction) NamesOfComponentsMoving(idGen *id.IDGenerator) string {
	ar := ins.ComponentsMoving(idGen)

	sa := make([]string, 0)

	for _, c := range ar {
		sa = append(sa, c.CName)
	}

	return strings.Join(sa, "+")
}

func (ins *LHInstruction) ComponentsMoving(idGen *id.IDGenerator) []*Liquid {
	ca := make([]*Liquid, 0)
	for i, v := range ins.Inputs {
		// ignore component 1 if this is a mix-in-place
		if i == 0 && ins.IsMixInPlace() {
			continue
		}
		ca = append(ca, v.Dup(idGen))
	}

	return ca
}

func (ins *LHInstruction) Wellcoords() WellCoords {
	return MakeWellCoords(ins.Welladdress)
}

func (ins *LHInstruction) AdjustVolumesBy(r float64) {
	// each subcomponent is assumed to scale linearly
	for _, c := range ins.Inputs {
		c.Vol *= r
	}
	for _, rslt := range ins.Outputs {
		rslt.Vol *= r
	}
}

func (ins *LHInstruction) InputVolumeMap(addition wunit.Volume) map[string]wunit.Volume {
	r := make(map[string]wunit.Volume, len(ins.Inputs))
	for i, c := range ins.Inputs {
		nom := c.FullyQualifiedName()
		myAdd := addition.Dup()

		if ins.IsMixInPlace() && i == 0 {
			nom += InPlaceMarker
			myAdd = wunit.ZeroVolume()
		}

		v, ok := r[nom]

		if ok {
			v.Add(c.Volume())
			v.Add(myAdd)
		} else {
			r[nom] = c.Volume()
			r[nom].Add(myAdd)
		}
	}

	return r
}
