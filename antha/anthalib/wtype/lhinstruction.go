package wtype

import (
	"strings"
)

// enum of instruction types

const (
	LHIEND = iota
	LHIMIX
	LHIWAI
	LHIPRM
)

var InsNames = []string{"END", "MIX", "WAIT", "PROMPT"}

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
	ProductID        string
	BlockID          BlockID
	SName            string
	Order            int
	Components       []*LHComponent
	ContainerType    string
	Welladdress      string
	PlateID          string
	Platetype        string
	Vol              float64
	Type             int
	Conc             float64
	Tvol             float64
	Majorlayoutgroup int
	Result           *LHComponent
	gen              int
	PlateName        string
	OutPlate         *LHPlate
	Message          string
	PassThrough      map[string]*LHComponent // 1:1 pass through, only applies to prompts
}

// privatised in favour of specific instruction constructors
func newLHInstruction() *LHInstruction {
	var lhi LHInstruction
	lhi.ID = GetUUID()
	lhi.Majorlayoutgroup = -1
	lhi.PassThrough = make(map[string]*LHComponent, 1)
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

// GetID returns the ID of the instruction, useful for interfaces
func (inst *LHInstruction) GetID() string {
	return inst.ID
}

func (inst *LHInstruction) AddProduct(cmp *LHComponent) {
	inst.Result = cmp
	inst.ProductID = cmp.ID
}

func (inst *LHInstruction) AddComponent(cmp *LHComponent) {
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

func (ins *LHInstruction) ComponentsMoving() string {
	sa := make([]string, 0, 1)
	for i, v := range ins.Components {
		// ignore component 1 if this is a mix-in-place
		if i == 0 && !v.IsSample() {
			continue
		}
		sa = append(sa, v.CName)
	}
	return strings.Join(sa, "+")
}

func (ins *LHInstruction) Wellcoords() WellCoords {
	return MakeWellCoords(ins.Welladdress)
}
