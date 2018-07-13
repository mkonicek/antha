package liquidhandling

import (
	"context"
	"fmt"
	"reflect"

	"github.com/antha-lang/antha/antha/anthalib/wtype"
)

// which instruction parameters for each instruction type
type InstructionParametersMap map[*InstructionType]InstructionParameters

func (a InstructionParametersMap) merge(b InstructionParametersMap) InstructionParametersMap {
	if len(a) == 0 {
		return b
	} else if len(b) == 0 {
		return a
	} else {
		result := make(InstructionParametersMap)
		for ins, params := range a {
			result[ins] = params.clone()
		}
		for ins, params := range b {
			// this is safe even when result[ins] is nil
			result[ins] = result[ins].merge(params)
		}
		return result
	}
}

type ComparisonOpt struct {
	InstructionParameters InstructionParametersMap
}

type ComparisonResult struct {
	Errors []error
}

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

func CompareTestInstructionSets(setA, setB []interface{}, opt ComparisonOpt) ComparisonResult {
	// convert
	cnv := func(s []interface{}) []RobotInstruction {
		r := make([]RobotInstruction, len(s))
		for i, n := range s {
			r[i] = n.(RobotInstruction)
		}
		return r
	}

	return CompareInstructionSets(cnv(setA), cnv(setB), opt)
}

func CompareInstructionSets(setA, setB []RobotInstruction, opt ComparisonOpt) ComparisonResult {

	setAMerged := mergeMovs(setA)
	setBMerged := mergeMovs(setB)
	return orderedInstructionComparison(setAMerged, setBMerged, opt)
}

func orderedInstructionComparison(setA, setB []RobotInstruction, opt ComparisonOpt) ComparisonResult {
	errors := make([]error, 0, len(setA))
	// v0 is just barrel through
	// v1 is aligning the instruction sets

	lenA, lenB := len(setA), len(setB)
	lenToCompare := lenA
	if lenA != lenB {
		errors = append(errors, fmt.Errorf("Instruction set lengths differ (%d %d)", lenA, lenB))

		if lenB < lenA {
			lenToCompare = lenB
		}
	}

	for i := 0; i < lenToCompare; i++ {
		errs := compareInstructions(i, setA[i], setB[i], opt.InstructionParameters[setA[i].Type()])
		errors = append(errors, errs...)
	}

	return ComparisonResult{Errors: errors}
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

func (mm MovMix) GetParameter(name InstructionParameter) interface{} {
	return getParameter(name, mm.Mov, mm.Mix)
}

func (mm MovMix) Generate(ctx context.Context, policy *wtype.LHPolicyRuleSet, prms *LHProperties) ([]RobotInstruction, error) {
	return []RobotInstruction{}, nil
}

func (mm MovMix) Check(lhpr wtype.LHPolicyRule) bool {
	return false
}

func compareInstructions(index int, ins1, ins2 RobotInstruction, paramsToCompare InstructionParameters) []error {
	// if instructions are not the same type we just return that error

	if ins1.Type() != ins2.Type() {
		return []error{fmt.Errorf("Instructions at %d different types (%s %s)", index, ins1.Type().MachineName, ins2.Type().MachineName)}
	}

	errors := make([]error, 0, 1)

	for prm := range paramsToCompare {
		p1 := ins1.GetParameter(prm)
		p2 := ins2.GetParameter(prm)

		if !reflect.DeepEqual(p1, p2) {
			errors = append(errors, fmt.Errorf("Instructions at index %d type %s parameter %s differ (%v %v)", index, ins1.Type().MachineName, prm, p1, p2))
		}
	}

	return errors
}

// convenience sets of parameters to compare

func CompareAllParameters() InstructionParametersMap {
	m := make(InstructionParametersMap)
	m = m.merge(CompareVolumes())
	m = m.merge(ComparePositions())
	m = m.merge(CompareWells())
	m = m.merge(CompareOffsets())
	m = m.merge(CompareReferences())
	m = m.merge(ComparePlateTypes())
	return m
}

func CompareReferences() InstructionParametersMap {
	m := make(InstructionParametersMap)
	m = m.merge(CompareAspReference())
	m = m.merge(CompareDspReference())
	m = m.merge(CompareMixReference())
	return m
}

func CompareAspReference() InstructionParametersMap {
	return InstructionParametersMap{
		MAS: NewInstructionParameters(REFERENCE, WHAT),
	}
}

func CompareDspReference() InstructionParametersMap {
	return InstructionParametersMap{
		MDS: NewInstructionParameters(REFERENCE, WHAT),
		MBL: NewInstructionParameters(REFERENCE, WHAT),
	}
}

func CompareMixReference() InstructionParametersMap {
	return InstructionParametersMap{
		MVM: NewInstructionParameters(REFERENCE, WHAT),
	}
}

func CompareOffsets() InstructionParametersMap {
	m := make(InstructionParametersMap)
	m = m.merge(CompareSourceOffsets())
	m = m.merge(CompareDestOffsets())
	m = m.merge(CompareMixOffsets())
	return m
}

func CompareSourceOffsets() InstructionParametersMap {
	return InstructionParametersMap{
		MAS: NewInstructionParameters(OFFSETX, OFFSETY, OFFSETZ, WHAT),
	}
}
func CompareDestOffsets() InstructionParametersMap {
	return InstructionParametersMap{
		MDS: NewInstructionParameters(OFFSETX, OFFSETY, OFFSETZ, WHAT),
		MBL: NewInstructionParameters(OFFSETX, OFFSETY, OFFSETZ, WHAT),
	}
}
func CompareMixOffsets() InstructionParametersMap {
	return InstructionParametersMap{
		MVM: NewInstructionParameters(OFFSETX, OFFSETY, OFFSETZ, WHAT),
	}
}
func CompareVolumes() InstructionParametersMap {
	vw := NewInstructionParameters(VOLUME, WHAT)
	return InstructionParametersMap{
		ASP: vw.clone(),
		MAS: vw.clone(),
		DSP: vw.clone(),
		MDS: vw.clone(),
		BLO: vw.clone(),
		MBL: vw.clone(),
		MIX: vw.clone(),
		MVM: vw.clone(),
	}
}

func ComparePositions() InstructionParametersMap {
	m := make(InstructionParametersMap)
	m = m.merge(CompareSourcePosition())
	m = m.merge(CompareDestPosition())
	m = m.merge(CompareMixPosition())
	return m
}

func CompareSourcePosition() InstructionParametersMap {
	return InstructionParametersMap{
		MAS: NewInstructionParameters(POSTO, WHAT),
	}
}

func CompareDestPosition() InstructionParametersMap {
	return InstructionParametersMap{
		MDS: NewInstructionParameters(POSTO, WHAT),
		MBL: NewInstructionParameters(POSTO, WHAT),
	}
}

func CompareMixPosition() InstructionParametersMap {
	return InstructionParametersMap{
		MVM: NewInstructionParameters(POSTO, WHAT),
	}
}

func CompareWells() InstructionParametersMap {
	m := make(InstructionParametersMap)
	m = m.merge(CompareSourceWell())
	m = m.merge(CompareDestWell())
	return m
}

func CompareSourceWell() InstructionParametersMap {
	return InstructionParametersMap{
		MAS: NewInstructionParameters(WELLTO, WHAT),
	}
}

func CompareDestWell() InstructionParametersMap {
	return InstructionParametersMap{
		MDS: NewInstructionParameters(WELLTO, WHAT),
		MBL: NewInstructionParameters(WELLTO, WHAT),
	}
}

func ComparePlateTypes() InstructionParametersMap {
	m := make(InstructionParametersMap)
	m = m.merge(CompareSourcePlateType())
	m = m.merge(CompareDestPlateType())
	return m
}

func CompareSourcePlateType() InstructionParametersMap {
	return InstructionParametersMap{
		MAS: NewInstructionParameters(PLATE, WHAT),
	}
}

func CompareDestPlateType() InstructionParametersMap {
	return InstructionParametersMap{
		MDS: NewInstructionParameters(PLATE, WHAT),
		MBL: NewInstructionParameters(PLATE, WHAT),
	}
}
