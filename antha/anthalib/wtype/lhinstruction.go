package wtype

import (
	"fmt"
	"strings"

	"github.com/antha-lang/antha/antha/anthalib/wunit"
)

// enum of instruction types

const (
	LHIEND = iota
	LHIMIX
	LHIWAI
	LHIPRM
	LHISPL
)

var InsNames = []string{"END", "MIX", "WAIT", "PROMPT", "SPLIT"}

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
	SName            string
	Order            int
	Components       []*Liquid
	ContainerType    string
	Welladdress      string
	PlateID          string
	Platetype        string
	Vol              float64
	Type             int
	Conc             float64
	Tvol             float64
	Majorlayoutgroup int
	Results          []*Liquid
	gen              int
	PlateName        string
	OutPlate         *Plate
	Message          string
	PassThrough      map[string]*Liquid // 1:1 pass through, only applies to prompts
}

func (ins LHInstruction) String() string {
	ret := fmt.Sprintf("%v G:%d ID:%s CMP:%v %s ID(%s) %s: %v", ins.InsType(), ins.Generation(), ins.ID, ComponentVector(ins.Components), ins.PlateName, ins.PlateID, ins.Welladdress, ins.ProductIDs())
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
		for _, cmp := range ins.Components {
			lines = append(lines, fmt.Sprintf("%s  %s", indentStr, cmp.Summarize()))
		}
		if ins.Welladdress != "" && ins.OutPlate != nil {
			lines = append(lines, fmt.Sprintf(indentStr+"to well %s in plate %s", ins.Welladdress, ins.OutPlate.Name()))
		} else if ins.Platetype != "" {
			lines = append(lines, fmt.Sprintf(indentStr+"to plate of type %s", ins.Platetype))
		}

		if len(ins.Results) > 0 {
			lines = append(lines, fmt.Sprintf(indentStr+"Resulting volume: %v", ins.Results[0].Summarize()))
		}

	default:
		return indentStr + ins.String()
	}

	return strings.Join(lines, "\n")
}

func (lhi *LHInstruction) ProductIDs() []string {
	r := make([]string, 0, len(lhi.Results))

	for _, p := range lhi.Results {
		r = append(r, p.ID)
	}

	return r
}

func (lhi *LHInstruction) InputIDs() []string {
	r := make([]string, 0, len(lhi.Components))
	for _, c := range lhi.Components {
		r = append(r, c.ID)
	}
	return r
}

func (lhi *LHInstruction) GetPlateType() string {
	if lhi.OutPlate != nil {
		return lhi.OutPlate.Type
	}

	return lhi.Platetype
}

// privatised in favour of specific instruction constructors
func newLHInstruction() *LHInstruction {
	var lhi LHInstruction
	lhi.ID = GetUUID()
	lhi.Majorlayoutgroup = -1
	lhi.PassThrough = make(map[string]*Liquid, 1)
	return &lhi
}

func NewLHMixInstruction() *LHInstruction {
	lhi := newLHInstruction()
	lhi.Type = LHIMIX
	return lhi

}

func NewLHPromptInstruction() *LHInstruction {
	lhi := newLHInstruction()
	lhi.Type = LHIPRM
	return lhi
}

func NewLHSplitInstruction() *LHInstruction {
	lhi := newLHInstruction()
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

func (ins *LHInstruction) AddResult(cmp *Liquid) {
	ins.AddProduct(cmp)
}

func (inst *LHInstruction) AddProduct(cmp *Liquid) {
	inst.Results = append(inst.Results, cmp)
}

func (inst *LHInstruction) AddComponent(cmp *Liquid) {
	if inst == nil {
		return
	}

	inst.Components = append(inst.Components, cmp)
}

func (ins *LHInstruction) Generation() int {
	return ins.gen
}
func (ins *LHInstruction) SetGeneration(i int) {
	ins.gen = i
}

func (ins *LHInstruction) GetPlateID() string {
	return ins.PlateID
}

func (ins *LHInstruction) SetPlateID(pid string) {
	ins.PlateID = pid
}

func (ins *LHInstruction) IsMixInPlace() bool {
	if ins == nil {
		return false
	}

	if len(ins.Components) == 0 {
		return false
	}

	smp := ins.Components[0].IsSample()
	return !smp
}

func (ins *LHInstruction) HasAnyParent() bool {
	for _, v := range ins.Components {
		if v.HasAnyParent() {
			return true
		}
	}

	return false
}

func (ins *LHInstruction) HasParent(id string) bool {
	for _, v := range ins.Components {
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

	for _, v := range ins.Components {
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

func (ins *LHInstruction) NamesOfComponentsMoving() string {
	ar := ins.ComponentsMoving()

	sa := make([]string, 0)

	for _, c := range ar {
		sa = append(sa, c.CName)
	}

	return strings.Join(sa, "+")
}

func (ins *LHInstruction) ComponentsMoving() []*Liquid {
	ca := make([]*Liquid, 0)
	for i, v := range ins.Components {
		// ignore component 1 if this is a mix-in-place
		if i == 0 && ins.IsMixInPlace() {
			continue
		}
		ca = append(ca, v.Dup())
	}

	return ca
}

func (ins *LHInstruction) Wellcoords() WellCoords {
	return MakeWellCoords(ins.Welladdress)
}

func (ins *LHInstruction) AdjustVolumesBy(r float64) {
	// each subcomponent is assumed to scale linearly
	for _, c := range ins.Components {
		c.Vol *= r
	}
	for _, rslt := range ins.Results {
		rslt.Vol *= r
	}
}

func (ins *LHInstruction) InputVolumeMap(addition wunit.Volume) map[string]wunit.Volume {
	r := make(map[string]wunit.Volume, len(ins.Components))
	for i, c := range ins.Components {
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
