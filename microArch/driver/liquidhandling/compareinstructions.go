package liquidhandling

import (
	"context"
	"fmt"
	"reflect"

	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/utils"
)

func mergeMovs(ris []RobotInstruction) []RobotInstruction {
	insOut := make([]RobotInstruction, 0, len(ris))
	for i := 0; i < len(ris); i++ {
		cur := ris[i]
		if i+1 == len(ris) { // last instruction so can't merge
			insOut = append(insOut, cur)
		} else {
			next := ris[i+1]
			merged := cur.MaybeMerge(next)
			if merged == cur { // it didn't merge
				insOut = append(insOut, cur)
			} else { // merged, so only append the merged instruction and skip over next
				insOut = append(insOut, merged)
				i++
			}
		}
	}

	return insOut
}

// CompareInstructionSets will use the comparators to detect
// differences between setA and setB. This allows you to compare,
// using any arbitrary property, instructions of the same type. The
// comparison is between the minimum length of setA and setB - an
// unmatched suffix will cause an error but will not be passed to any
// comparator.
func CompareInstructionSets(setA, setB []RobotInstruction, comparators ...RobotInstructionComparatorFunc) utils.ErrorSlice {
	setAMerged := mergeMovs(setA)
	setBMerged := mergeMovs(setB)
	return orderedInstructionComparison(setAMerged, setBMerged, comparators)
}

func orderedInstructionComparison(setA, setB []RobotInstruction, comparators []RobotInstructionComparatorFunc) utils.ErrorSlice {
	var errs utils.ErrorSlice

	lenA, lenB := len(setA), len(setB)
	lenToCompare := lenA
	if lenA != lenB {
		errs = append(errs, fmt.Errorf("Instruction set lengths differ (%d %d)", lenA, lenB))

		if lenB < lenA {
			lenToCompare = lenB
		}
	}

	for i := 0; i < lenToCompare; i++ {
		if err := compareRobotInstructions(setA[i], setB[i], comparators); len(err) != 0 {
			errs = append(errs, err...)
		}
	}

	return errs
}

// merge move and asp / move and dsp / mov and mix

var moveParams = []InstructionParameter{POSTO, WELLTO, REFERENCE, OFFSETX, OFFSETY, OFFSETZ}

func isIn(s InstructionParameter, ar []InstructionParameter) bool {
	for _, s2 := range ar {
		if s2 == s {
			return true
		}
	}

	return false
}

func getParameter(name InstructionParameter, movIns *MoveInstruction, otherIns RobotInstruction) interface{} {
	// indirect calls appropriately
	if isIn(name, moveParams) {
		return movIns.GetParameter(name)
	} else {
		return otherIns.GetParameter(name)
	}
}

type MovAsp struct {
	*InstructionType
	BaseRobotInstruction
	Mov *MoveInstruction
	Asp *AspirateInstruction
}

func NewMovAsp(mov *MoveInstruction, asp *AspirateInstruction) *MovAsp {
	ma := &MovAsp{
		InstructionType: MAS,
		Mov:             mov,
		Asp:             asp,
	}
	ma.BaseRobotInstruction = NewBaseRobotInstruction(ma)
	return ma
}

func (ma *MovAsp) Visit(visitor RobotInstructionVisitor) {
	visitor.MovAsp(ma)
}

func (ma MovAsp) GetParameter(name InstructionParameter) interface{} {
	return getParameter(name, ma.Mov, ma.Asp)
}

func (ma MovAsp) Generate(ctx context.Context, policy *wtype.LHPolicyRuleSet, prms *LHProperties) ([]RobotInstruction, error) {
	return []RobotInstruction{}, nil
}

func (ma MovAsp) Check(lhpr wtype.LHPolicyRule) bool {
	return false
}

type MovDsp struct {
	*InstructionType
	BaseRobotInstruction
	Mov *MoveInstruction
	Dsp *DispenseInstruction
}

func NewMovDsp(mov *MoveInstruction, dsp *DispenseInstruction) *MovDsp {
	md := &MovDsp{
		InstructionType: MDS,
		Mov:             mov,
		Dsp:             dsp,
	}
	md.BaseRobotInstruction = NewBaseRobotInstruction(md)
	return md
}

func (md *MovDsp) Visit(visitor RobotInstructionVisitor) {
	visitor.MovDsp(md)
}

func (md MovDsp) GetParameter(name InstructionParameter) interface{} {
	return getParameter(name, md.Mov, md.Dsp)
}

func (md MovDsp) Generate(ctx context.Context, policy *wtype.LHPolicyRuleSet, prms *LHProperties) ([]RobotInstruction, error) {
	return []RobotInstruction{}, nil
}

func (md MovDsp) Check(lhpr wtype.LHPolicyRule) bool {
	return false
}

type MovBlo struct {
	*InstructionType
	BaseRobotInstruction
	Mov *MoveInstruction
	Blo *BlowoutInstruction
}

func NewMovBlo(mov *MoveInstruction, blo *BlowoutInstruction) *MovBlo {
	mb := &MovBlo{
		InstructionType: MBL,
		Mov:             mov,
		Blo:             blo,
	}
	mb.BaseRobotInstruction = NewBaseRobotInstruction(mb)
	return mb
}

func (mb *MovBlo) Visit(visitor RobotInstructionVisitor) {
	visitor.MovBlo(mb)
}

func (mb MovBlo) GetParameter(name InstructionParameter) interface{} {
	return getParameter(name, mb.Mov, mb.Blo)
}

func (mb MovBlo) Generate(ctx context.Context, policy *wtype.LHPolicyRuleSet, prms *LHProperties) ([]RobotInstruction, error) {
	return []RobotInstruction{}, nil
}

func (mb MovBlo) Check(lhpr wtype.LHPolicyRule) bool {
	return false
}

type MovMix struct {
	*InstructionType
	BaseRobotInstruction
	Mov *MoveInstruction
	Mix *MixInstruction
}

func NewMovMix(mov *MoveInstruction, mix *MixInstruction) *MovMix {
	mm := &MovMix{
		InstructionType: MVM,
		Mov:             mov,
		Mix:             mix,
	}
	mm.BaseRobotInstruction = NewBaseRobotInstruction(mm)
	return mm
}

func (mm *MovMix) Visit(visitor RobotInstructionVisitor) {
	visitor.MovMix(mm)
}

func (mm MovMix) GetParameter(name InstructionParameter) interface{} {
	return getParameter(name, mm.Mov, mm.Mix)
}

func (mm MovMix) Generate(ctx context.Context, policy *wtype.LHPolicyRuleSet, prms *LHProperties) ([]RobotInstruction, error) {
	return []RobotInstruction{}, nil
}

func (mm MovMix) Check(lhpr wtype.LHPolicyRule) bool {
	return false
}

type robotInstructionComparator struct {
	// a and b are the two RobotInstructions we are comparing. Cur is
	// the "current" instruction so that our visitors know what they're
	// looking at.
	a, b, cur RobotInstruction
	// aFields and bFields contain the fields we extract from the
	// instructions. Ultimately, they will be compared with
	// reflect.DeepEqual
	aFields, bFields []interface{}
}

type RobotInstructionComparatorFunc func(*robotInstructionComparator) *RobotInstructionBaseVisitor

// we deliberately return the ground type errorSlice so that the
// caller can flatten together multiple errorSlices.
func compareRobotInstructions(a, b RobotInstruction, comparators []RobotInstructionComparatorFunc) utils.ErrorSlice {
	if a == nil || b == nil {
		return utils.ErrorSlice{fmt.Errorf("Cannot compare with nil RobotInstructions: %#v %#v", a, b)}
	} else if a.Type() != b.Type() {
		return utils.ErrorSlice{fmt.Errorf("Instructions of different types (%s != %s)", a.Type(), b.Type())}
	}

	var errs utils.ErrorSlice
	ric := &robotInstructionComparator{a: a, b: b}
	for _, comp := range comparators {
		ric.aFields = ric.aFields[:0]
		ric.bFields = ric.bFields[:0]
		ric.cur = ric.a
		ric.cur.Visit(comp(ric))
		ric.cur = ric.b
		ric.cur.Visit(comp(ric))
		if !reflect.DeepEqual(ric.aFields, ric.bFields) {
			errs = append(errs, fmt.Errorf("Instructions of type %v differ in fields: %#v %#v", ric.cur.Type(), ric.a, ric.b))
		}
	}

	return errs
}

func (ric *robotInstructionComparator) appendFieldValues(vs ...interface{}) {
	if ric.cur == ric.a {
		ric.aFields = append(ric.aFields, vs...)
	} else {
		ric.bFields = append(ric.bFields, vs...)
	}
}

func CompareAspReference(ric *robotInstructionComparator) *RobotInstructionBaseVisitor {
	return &RobotInstructionBaseVisitor{
		HandleMovAsp: func(ins *MovAsp) { ric.appendFieldValues(ins.Mov.Reference, ins.Asp.What) },
	}
}

func CompareDspReference(ric *robotInstructionComparator) *RobotInstructionBaseVisitor {
	return &RobotInstructionBaseVisitor{
		HandleMovDsp: func(ins *MovDsp) { ric.appendFieldValues(ins.Mov.Reference, ins.Dsp.What) },
		HandleMovBlo: func(ins *MovBlo) { ric.appendFieldValues(ins.Mov.Reference, ins.Blo.What) },
	}
}

func CompareMixReference(ric *robotInstructionComparator) *RobotInstructionBaseVisitor {
	return &RobotInstructionBaseVisitor{
		HandleMovMix: func(ins *MovMix) { ric.appendFieldValues(ins.Mov.Reference, ins.Mix.What) },
	}
}

func CompareSourceOffsets(ric *robotInstructionComparator) *RobotInstructionBaseVisitor {
	return &RobotInstructionBaseVisitor{
		HandleMovAsp: func(ins *MovAsp) {
			ric.appendFieldValues(ins.Mov.OffsetX, ins.Mov.OffsetY, ins.Mov.OffsetZ, ins.Asp.What)
		},
	}
}
func CompareDestOffsets(ric *robotInstructionComparator) *RobotInstructionBaseVisitor {
	return &RobotInstructionBaseVisitor{
		HandleMovDsp: func(ins *MovDsp) {
			ric.appendFieldValues(ins.Mov.OffsetX, ins.Mov.OffsetY, ins.Mov.OffsetZ, ins.Dsp.What)
		},
		HandleMovBlo: func(ins *MovBlo) {
			ric.appendFieldValues(ins.Mov.OffsetX, ins.Mov.OffsetY, ins.Mov.OffsetZ, ins.Blo.What)
		},
	}
}
func CompareMixOffsets(ric *robotInstructionComparator) *RobotInstructionBaseVisitor {
	return &RobotInstructionBaseVisitor{
		HandleMovMix: func(ins *MovMix) {
			ric.appendFieldValues(ins.Mov.OffsetX, ins.Mov.OffsetY, ins.Mov.OffsetZ, ins.Mix.What)
		},
	}
}
func CompareVolumes(ric *robotInstructionComparator) *RobotInstructionBaseVisitor {
	return &RobotInstructionBaseVisitor{
		HandleAspirate: func(ins *AspirateInstruction) { ric.appendFieldValues(ins.Volume, ins.What) },
		HandleMovAsp:   func(ins *MovAsp) { ric.appendFieldValues(ins.Asp.Volume, ins.Asp.What) },
		HandleDispense: func(ins *DispenseInstruction) { ric.appendFieldValues(ins.Volume, ins.What) },
		HandleMovDsp:   func(ins *MovDsp) { ric.appendFieldValues(ins.Dsp.Volume, ins.Dsp.What) },
		HandleBlowout:  func(ins *BlowoutInstruction) { ric.appendFieldValues(ins.Volume, ins.What) },
		HandleMovBlo:   func(ins *MovBlo) { ric.appendFieldValues(ins.Blo.Volume, ins.Blo.What) },
		HandleMix:      func(ins *MixInstruction) { ric.appendFieldValues(ins.Volume, ins.What) },
		HandleMovMix:   func(ins *MovMix) { ric.appendFieldValues(ins.Mix.Volume, ins.Mix.What) },
	}
}

func CompareSourcePosition(ric *robotInstructionComparator) *RobotInstructionBaseVisitor {
	return &RobotInstructionBaseVisitor{
		HandleMovAsp: func(ins *MovAsp) { ric.appendFieldValues(ins.Mov.Pos, ins.Asp.What) },
	}
}

func CompareDestPosition(ric *robotInstructionComparator) *RobotInstructionBaseVisitor {
	return &RobotInstructionBaseVisitor{
		HandleMovDsp: func(ins *MovDsp) { ric.appendFieldValues(ins.Mov.Pos, ins.Dsp.What) },
		HandleMovBlo: func(ins *MovBlo) { ric.appendFieldValues(ins.Mov.Pos, ins.Blo.What) },
	}
}

func CompareMixPosition(ric *robotInstructionComparator) *RobotInstructionBaseVisitor {
	return &RobotInstructionBaseVisitor{
		HandleMovMix: func(ins *MovMix) { ric.appendFieldValues(ins.Mov.Pos, ins.Mix.What) },
	}
}

func CompareSourceWell(ric *robotInstructionComparator) *RobotInstructionBaseVisitor {
	return &RobotInstructionBaseVisitor{
		HandleMovAsp: func(ins *MovAsp) { ric.appendFieldValues(ins.Mov.Well, ins.Asp.What) },
	}
}

func CompareDestWell(ric *robotInstructionComparator) *RobotInstructionBaseVisitor {
	return &RobotInstructionBaseVisitor{
		HandleMovDsp: func(ins *MovDsp) { ric.appendFieldValues(ins.Mov.Well, ins.Dsp.What) },
		HandleMovBlo: func(ins *MovBlo) { ric.appendFieldValues(ins.Mov.Well, ins.Blo.What) },
	}
}

func CompareSourcePlateType(ric *robotInstructionComparator) *RobotInstructionBaseVisitor {
	return &RobotInstructionBaseVisitor{
		HandleMovAsp: func(ins *MovAsp) { ric.appendFieldValues(ins.Asp.Plt, ins.Asp.What) },
	}
}

func CompareDestPlateType(ric *robotInstructionComparator) *RobotInstructionBaseVisitor {
	return &RobotInstructionBaseVisitor{
		HandleMovDsp: func(ins *MovDsp) { ric.appendFieldValues(ins.Dsp.Plt, ins.Dsp.What) },
		HandleMovBlo: func(ins *MovBlo) { ric.appendFieldValues(ins.Blo.Plt, ins.Blo.What) },
	}
}

var (
	ComparePlateTypes = []RobotInstructionComparatorFunc{CompareSourcePlateType, CompareDestPlateType}
	CompareWells      = []RobotInstructionComparatorFunc{CompareSourceWell, CompareDestWell}
	ComparePositions  = []RobotInstructionComparatorFunc{CompareSourcePosition, CompareDestPosition, CompareMixPosition}
	CompareOffsets    = []RobotInstructionComparatorFunc{CompareSourceOffsets, CompareDestOffsets, CompareMixOffsets}
	CompareReferences = []RobotInstructionComparatorFunc{CompareAspReference, CompareDspReference, CompareMixReference}

	CompareAllParameters = []RobotInstructionComparatorFunc{
		CompareSourcePlateType, CompareDestPlateType,
		CompareSourceWell, CompareDestWell,
		CompareSourcePosition, CompareDestPosition, CompareMixPosition,
		CompareVolumes,
		CompareSourceOffsets, CompareDestOffsets, CompareMixOffsets,
		CompareAspReference, CompareDspReference, CompareMixReference}
)
