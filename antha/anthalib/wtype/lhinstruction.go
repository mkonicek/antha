package wtype

import (
	"fmt"
	"strings"

	"time"

	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/utils"
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
	Platetype        string
	Vol              float64
	Type             int
	Conc             float64
	Tvol             float64
	Majorlayoutgroup int
	gen              int
	PlateName        string
	OutPlate         *Plate
	Message          string
	WaitTime         time.Duration
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
func newLHInstruction() *LHInstruction {
	var lhi LHInstruction
	lhi.ID = GetUUID()
	lhi.Majorlayoutgroup = -1
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
	for i, v := range ins.Inputs {
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

func (ins *LHInstruction) DupLiquids() {
	for i := 0; i < len(ins.Inputs); i++ {
		ins.Inputs[i] = ins.Inputs[i].Dup()
	}
	for i := 0; i < len(ins.Outputs); i++ {
		ins.Outputs[i] = ins.Outputs[i].Dup()
	}
}

type LHInstructions map[string]*LHInstruction

// DupLiquids duplicate the input and output liquids for each LHInstruction,
// thereby making certain that there can be no pointer-reuse between the instructions
func (insts LHInstructions) DupLiquids() {
	for _, ins := range insts {
		ins.DupLiquids()
	}
}

// AssertNoPointerReuse make certain that inputs and outputs are not shared among LHInstructions
// as should be ensured by calling DupLiquids()
func (insts LHInstructions) AssertNoPointerReuse() error {
	seen := map[*Liquid]*LHInstruction{}
	errs := make(utils.ErrorSlice, 0, len(insts))
	for _, ins := range insts {
		for _, c := range append(ins.Inputs, ins.Outputs...) {
			if ins2, ok := seen[c]; ok {
				errs = append(errs, fmt.Errorf("POINTER REUSE: Instructions share *Liquid(%p): %s\n\tA: %s\n\tB: %s", c, c.CNID(), ins, ins2))
			}
			seen[c] = ins
		}
	}

	return errs.Pack()
}

// AssertDestinationsSet make certain that a destination has been set for each instruction,
// returning a descriptive error if not
func (insts LHInstructions) AssertDestinationsSet() error {
	errs := make(utils.ErrorSlice, 0, len(insts))

	for _, ins := range insts {
		if ins.Type != LHIMIX {
			continue
		}

		if ins.PlateID == "" || ins.Platetype == "" || ins.Welladdress == "" {
			errs = append(errs, fmt.Errorf("INS %s missing destination: has PlateID/Platetype/Welladdress: %t/%t/%t ", ins, ins.PlateID == "", ins.Platetype == "", ins.Welladdress != ""))
		}
	}

	return errs.Pack()
}

//assertVolumesNonNegative tests that the volumes within the LHRequest are zero or positive
func (insts LHInstructions) AssertVolumesNonNegative() error {
	errs := make(utils.ErrorSlice, 0, len(insts))

	for _, ins := range insts {
		if ins.Type != LHIMIX {
			continue
		}

		for _, cmp := range ins.Inputs {
			if cmp.Volume().IsNegative() {
				errs = append(errs, LHErrorf(LH_ERR_VOL, "negative volume for component \"%s\" in instruction:\n%s", cmp.CName, ins.Summarize(1)))
			}
		}
	}
	return errs.Pack()
}

//assertTotalVolumesMatch checks that component total volumes are all the same in mix instructions
func (insts LHInstructions) AssertTotalVolumesMatch() error {
	errs := make(utils.ErrorSlice, 0, len(insts))
	for _, ins := range insts {
		if ins.Type != LHIMIX {
			continue
		}

		totalVolume := wunit.ZeroVolume()

		for _, cmp := range ins.Inputs {
			if tV := cmp.TotalVolume(); !tV.IsZero() {
				if !totalVolume.IsZero() && !tV.EqualTo(totalVolume) {
					errs = append(errs, LHErrorf(LH_ERR_VOL, "multiple distinct total volumes specified in instruction:\n%s", ins.Summarize(1)))
				}
				totalVolume = tV
			}
		}
	}
	return errs.Pack()
}

//assertMixResultsCorrect checks that volumes of the mix result matches either the sum of the input, or the total volume if specified
func (insts LHInstructions) AssertMixResultsCorrect() error {
	errs := make(utils.ErrorSlice, 0, len(insts))
	for _, ins := range insts {
		if ins.Type != LHIMIX {
			continue
		}

		totalVolume := wunit.ZeroVolume()
		volumeSum := wunit.ZeroVolume()

		for _, cmp := range ins.Inputs {
			if tV := cmp.TotalVolume(); !tV.IsZero() {
				totalVolume = tV
			} else if v := cmp.Volume(); !v.IsZero() {
				volumeSum.Add(v)
			}
		}

		if len(ins.Outputs) != 1 {
			errs = append(errs, LHErrorf(LH_ERR_DIRE, "mix instruction has %d results specified, expecting 1 at instruction:\n%s",
				len(ins.Outputs), ins.Summarize(1)))
		} else if resultVolume := ins.Outputs[0].Volume(); !totalVolume.IsZero() && !totalVolume.EqualTo(resultVolume) {
			errs = append(errs, LHErrorf(LH_ERR_VOL, "total volume (%v) does not match resulting volume (%v) for instruction:\n%s",
				totalVolume, resultVolume, ins.Summarize(1)))
		} else if totalVolume.IsZero() && !volumeSum.EqualTo(resultVolume) {
			errs = append(errs, LHErrorf(LH_ERR_VOL, "sum of requested volumes (%v) does not match result volume (%v) for instruction:\n%s",
				volumeSum, resultVolume, ins.Summarize(1)))
		}
	}
	return errs.Pack()
}
