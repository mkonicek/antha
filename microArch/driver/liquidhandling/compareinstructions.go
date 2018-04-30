package liquidhandling

import (
	"context"
	"fmt"
	"reflect"

	"github.com/antha-lang/antha/antha/anthalib/wtype"
)

type ComparisonOpt struct {
	InstructionParameters map[string][]string // which instruction parameters for each instruction type
}

type ComparisonResult struct {
	Errors []error
}

func mergeMovs(ris []RobotInstruction) []RobotInstruction {
	insOut := make([]RobotInstruction, 0, len(ris))
	for i := 0; i < len(ris); i++ {
		ins := ris[i]
		if InstructionTypeName(ins) == "MOV" && (i != len(ris)-1) {
			next := ris[i+1]
			if InstructionTypeName(next) == "ASP" {
				insOut = append(insOut, MovAsp{Asp: next.(*AspirateInstruction), Mov: ins.(*MoveInstruction)})
				i += 1 // skip
			} else if InstructionTypeName(next) == "DSP" {
				insOut = append(insOut, MovDsp{Dsp: next.(*DispenseInstruction), Mov: ins.(*MoveInstruction)})
				i += 1 // skip
			} else if InstructionTypeName(next) == "MIX" {
				insOut = append(insOut, MovMix{Mix: next.(*MixInstruction), Mov: ins.(*MoveInstruction)})
				i += 1 // skip
			} else if InstructionTypeName(next) == "BLO" {
				insOut = append(insOut, MovBlo{Blo: next.(*BlowoutInstruction), Mov: ins.(*MoveInstruction)})
				i += 1 // skip
			} else {
				insOut = append(insOut, ins)
			}
		} else {
			insOut = append(insOut, ins)
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
	lenToCompare := len(setA)

	if len(setB) != len(setA) {
		errors = append(errors, fmt.Errorf("Instruction set lengths differ (%d %d)", len(setA), len(setB)))

		if len(setB) < len(setA) {
			lenToCompare = len(setB)
		}
	}

	for i := 0; i < lenToCompare; i++ {
		errs := compareInstructions(i, setA[i], setB[i], opt.InstructionParameters[InstructionTypeName(setA[i])])
		errors = append(errors, errs...)
	}

	return ComparisonResult{Errors: errors}
}

// merge move and asp / move and dsp / mov and mix

/*

   InstructionType() int
   GetParameter(name string) interface{}
   Generate(ctx context.Context,policy *wtype.LHPolicyRuleSet, prms *LHProperties) ([]RobotInstruction, error)
   Check(lhpr wtype.LHPolicyRule) bool

*/

type MovAsp struct {
	Mov *MoveInstruction
	Asp *AspirateInstruction
}

var moveParams = []string{"POSTO", "WELLTO", "REFERENCE", "OFFSETX", "OFFSETY", "OFFSETZ"}

func isIn(s string, ar []string) bool {
	for _, s2 := range ar {
		if s2 == s {
			return true
		}
	}

	return false
}

func getParameter(name string, movIns *MoveInstruction, otherIns RobotInstruction) interface{} {
	// indirect calls appropriately
	if isIn(name, moveParams) {
		return movIns.GetParameter(name)
	} else {
		return otherIns.GetParameter(name)
	}
}

func (ma MovAsp) InstructionType() int {
	return MAS
}

func (ma MovAsp) GetParameter(name string) interface{} {
	return getParameter(name, ma.Mov, ma.Asp)
}

func (ma MovAsp) Generate(ctx context.Context, policy *wtype.LHPolicyRuleSet, prms *LHProperties) ([]RobotInstruction, error) {
	return []RobotInstruction{}, nil
}

func (ma MovAsp) Check(lhpr wtype.LHPolicyRule) bool {
	return false
}

type MovDsp struct {
	Mov *MoveInstruction
	Dsp *DispenseInstruction
}

func (md MovDsp) InstructionType() int {
	return MDS
}

func (md MovDsp) GetParameter(name string) interface{} {
	return getParameter(name, md.Mov, md.Dsp)
}

func (md MovDsp) Generate(ctx context.Context, policy *wtype.LHPolicyRuleSet, prms *LHProperties) ([]RobotInstruction, error) {
	return []RobotInstruction{}, nil
}

func (md MovDsp) Check(lhpr wtype.LHPolicyRule) bool {
	return false
}

type MovBlo struct {
	Mov *MoveInstruction
	Blo *BlowoutInstruction
}

func (mb MovBlo) InstructionType() int {
	return MBL
}

func (mb MovBlo) GetParameter(name string) interface{} {
	return getParameter(name, mb.Mov, mb.Blo)
}

func (mb MovBlo) Generate(ctx context.Context, policy *wtype.LHPolicyRuleSet, prms *LHProperties) ([]RobotInstruction, error) {
	return []RobotInstruction{}, nil
}

func (mb MovBlo) Check(lhpr wtype.LHPolicyRule) bool {
	return false
}

type MovMix struct {
	Mov *MoveInstruction
	Mix *MixInstruction
}

func (mm MovMix) InstructionType() int {
	return MVM
}

func (mm MovMix) GetParameter(name string) interface{} {
	return getParameter(name, mm.Mov, mm.Mix)
}

func (mm MovMix) Generate(ctx context.Context, policy *wtype.LHPolicyRuleSet, prms *LHProperties) ([]RobotInstruction, error) {
	return []RobotInstruction{}, nil
}

func (mm MovMix) Check(lhpr wtype.LHPolicyRule) bool {
	return false
}

func compareInstructions(index int, ins1, ins2 RobotInstruction, paramsToCompare []string) []error {
	// if instructions are not the same type we just return that error

	if ins1.InstructionType() != ins2.InstructionType() {
		return []error{fmt.Errorf("Instructions at %d different types (%s %s)", index, InstructionTypeName(ins1), InstructionTypeName(ins2))}
	}

	errors := make([]error, 0, 1)

	for _, prm := range paramsToCompare {
		p1 := ins1.GetParameter(prm)
		p2 := ins2.GetParameter(prm)

		if !reflect.DeepEqual(p1, p2) {
			errors = append(errors, fmt.Errorf("Instructions at index %d type %s parameter %s differ (%v %v)", index, InstructionTypeName(ins1), prm, p1, p2))
		}
	}

	return errors
}

func dupV(a []string) []string {
	r := make([]string, 0, len(a))
	r = append(r, a...)

	return r
}

func mergeStringSets(a, b []string) []string {
	r := dupV(a)

	for _, s := range b {
		if !isIn(s, r) {
			r = append(r, s)
		}
	}

	return r
}

func mergeSets(s1, s2 map[string][]string) map[string][]string {
	r := make(map[string][]string, len(s1))

	for k, v := range s1 {
		r[k] = dupV(v)
	}

	for k, v := range s2 {
		v2, ok := r[k]

		if ok {
			r[k] = mergeStringSets(v, v2)
		} else {
			r[k] = v
		}
	}
	return r
}

// convenience sets of parameters to compare

func CompareAllParameters() map[string][]string {
	r := make(map[string][]string, 3)
	r = mergeSets(r, CompareVolumes())
	r = mergeSets(r, ComparePositions())
	r = mergeSets(r, CompareWells())
	r = mergeSets(r, CompareOffsets())
	r = mergeSets(r, CompareReferences())
	r = mergeSets(r, ComparePlateTypes())
	return r
}

func CompareReferences() map[string][]string {
	ret := make(map[string][]string, 2)
	ret = mergeSets(ret, CompareAspReference())
	ret = mergeSets(ret, CompareDspReference())
	ret = mergeSets(ret, CompareMixReference())
	return ret
}

func CompareAspReference() map[string][]string {
	ret := make(map[string][]string, 2)
	ret["MOVASP"] = []string{"REFERENCE", "WHAT"}
	return ret
}

func CompareDspReference() map[string][]string {
	ret := make(map[string][]string, 2)
	ret["MOVDSP"] = []string{"REFERENCE", "WHAT"}
	ret["MOVBLO"] = []string{"REFERENCE", "WHAT"}
	return ret
}

func CompareMixReference() map[string][]string {
	ret := make(map[string][]string, 2)
	ret["MOVMIX"] = []string{"REFERENCE", "WHAT"}
	return ret
}

func CompareOffsets() map[string][]string {
	ret := make(map[string][]string, 2)
	ret = mergeSets(ret, CompareSourceOffsets())
	ret = mergeSets(ret, CompareDestOffsets())
	ret = mergeSets(ret, CompareMixOffsets())
	return ret
}

func CompareSourceOffsets() map[string][]string {
	ret := make(map[string][]string, 2)
	ret["MOVASP"] = []string{"OFFSETX", "OFFSETY", "OFFSETZ", "WHAT"}
	return ret
}
func CompareDestOffsets() map[string][]string {
	ret := make(map[string][]string, 2)
	ret["MOVDSP"] = []string{"OFFSETX", "OFFSETY", "OFFSETZ", "WHAT"}
	ret["MOVBLO"] = []string{"OFFSETX", "OFFSETY", "OFFSETZ", "WHAT"}
	return ret
}
func CompareMixOffsets() map[string][]string {
	ret := make(map[string][]string, 2)
	ret["MOVMIX"] = []string{"OFFSETX", "OFFSETY", "OFFSETZ", "WHAT"}
	return ret
}
func CompareVolumes() map[string][]string {
	ret := make(map[string][]string, 2)
	ret["ASP"] = []string{"VOLUME", "WHAT"}
	ret["MOVASP"] = []string{"VOLUME", "WHAT"}
	ret["DSP"] = []string{"VOLUME", "WHAT"}
	ret["MOVDSP"] = []string{"VOLUME", "WHAT"}
	ret["BLO"] = []string{"VOLUME", "WHAT"}
	ret["MOVBLO"] = []string{"VOLUME", "WHAT"}
	ret["MIX"] = []string{"VOLUME", "WHAT"}
	ret["MOVMIX"] = []string{"VOLUME", "WHAT"}
	return ret
}

func ComparePositions() map[string][]string {
	ret := make(map[string][]string, 2)
	ret = mergeSets(ret, CompareSourcePosition())
	ret = mergeSets(ret, CompareDestPosition())
	ret = mergeSets(ret, CompareMixPosition())
	return ret
}

func CompareSourcePosition() map[string][]string {
	ret := make(map[string][]string, 2)
	ret["MOVASP"] = []string{"POSTO", "WHAT"}
	return ret
}

func CompareDestPosition() map[string][]string {
	ret := make(map[string][]string, 2)
	ret["MOVDSP"] = []string{"POSTO", "WHAT"}
	ret["MOVBLO"] = []string{"POSTO", "WHAT"}
	return ret
}

func CompareMixPosition() map[string][]string {
	ret := make(map[string][]string, 2)
	ret["MOVMIX"] = []string{"POSTO", "WHAT"}
	return ret
}

func CompareWells() map[string][]string {
	ret := make(map[string][]string, 2)
	ret = mergeSets(ret, CompareSourceWell())
	ret = mergeSets(ret, CompareDestWell())
	return ret
}

func CompareSourceWell() map[string][]string {
	ret := make(map[string][]string, 2)
	ret["MOVASP"] = []string{"WELLTO", "WHAT"}
	return ret
}

func CompareDestWell() map[string][]string {
	ret := make(map[string][]string, 2)
	ret["MOVDSP"] = []string{"WELLTO", "WHAT"}
	ret["MOVBLO"] = []string{"WELLTO", "WHAT"}
	return ret
}

func ComparePlateTypes() map[string][]string {
	ret := make(map[string][]string, 2)
	ret = mergeSets(ret, CompareSourcePlateType())
	ret = mergeSets(ret, CompareDestPlateType())
	return ret
}

func CompareSourcePlateType() map[string][]string {
	ret := make(map[string][]string, 2)
	ret["MOVASP"] = []string{"PLATE", "WHAT"}
	return ret
}

func CompareDestPlateType() map[string][]string {
	ret := make(map[string][]string, 2)
	ret["MOVDSP"] = []string{"PLATE", "WHAT"}
	ret["MOVBLO"] = []string{"PLATE", "WHAT"}
	return ret
}
