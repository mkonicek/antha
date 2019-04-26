package tests

import (
	"fmt"
	"math"
	"reflect"
	"sort"
	"strings"
	"testing"

	"github.com/pkg/errors"

	"github.com/antha-lang/antha/antha/anthalib/mixer"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/laboratory"
	"github.com/antha-lang/antha/laboratory/effects/id"
	"github.com/antha-lang/antha/laboratory/testlab"
	"github.com/antha-lang/antha/microArch/driver/liquidhandling"
	lh "github.com/antha-lang/antha/microArch/scheduler/liquidhandling"
	"github.com/antha-lang/antha/utils"
)

type PlanningTest struct {
	Name          string
	Liquidhandler *lh.Liquidhandler
	Instructions  InstructionBuilder
	InputPlates   []*wtype.LHPlate
	OutputPlates  []*wtype.LHPlate
	ErrorPrefix   string
	Assertions    Assertions
}

func (test *PlanningTest) Run(t *testing.T) {
	testlab.WithTestLab(t, "testdata", &testlab.TestElementCallbacks{
		Name:  test.Name,
		Steps: test.run,
	})
}

func (test *PlanningTest) run(lab *laboratory.Laboratory) error {
	request := lh.NewLHRequest(lab.IDGenerator)

	if instrs, err := test.Instructions(lab); err != nil {
		return err
	} else {
		for _, ins := range instrs {
			request.Add_instruction(ins)
		}
	}

	for _, plate := range test.InputPlates {
		if !plate.IsEmpty(lab.IDGenerator) {
			request.AddUserPlate(lab.IDGenerator, plate)
			plate = plate.Dup(lab.IDGenerator)
			plate.Clean()
		}
		request.InputPlatetypes = append(request.InputPlatetypes, plate)
	}

	request.OutputPlatetypes = append(request.OutputPlatetypes, test.OutputPlates...)

	if test.Liquidhandler == nil {
		test.Liquidhandler = GetLiquidHandlerForTest(lab)
	}

	if err := test.Liquidhandler.Plan(lab.LaboratoryEffects, request); !test.expected(err) {
		return fmt.Errorf("expecting error = %q: got error %q", test.ErrorPrefix, err.Error())
	}

	test.Assertions.Assert(lab, test.Liquidhandler, request)

	if test.ErrorPrefix == "" {
		return utils.ErrorSlice{
			test.checkPlateIDMap(lab),
			test.checkPositionConsistency(lab),
			test.checkSummaryGeneration(lab, request),
		}.Pack()
	}
	return nil
}

// checkSummaryGeneration check that we generate a valid JSON string from the objects,
// LayoutSummary and ActionSummary validate their output against the schemas agreed with UI
func (test *PlanningTest) checkSummaryGeneration(lab *laboratory.Laboratory, request *lh.LHRequest) error {
	if bs, err := lh.SummarizeLayout(lab.IDGenerator, test.Liquidhandler.Properties, test.Liquidhandler.FinalProperties, test.Liquidhandler.PlateIDMap()); err != nil {
		lab.Logger.Log("invalidLayout", string(bs))
		return err
	}

	if bs, err := lh.SummarizeActions(lab.IDGenerator, test.Liquidhandler.Properties, request.InstructionTree); err != nil {
		lab.Logger.Log("invalidActions", string(bs))
		return err
	}

	return nil
}

func (test *PlanningTest) checkPlateIDMap(lab *laboratory.Laboratory) error {
	beforePlates := test.Liquidhandler.Properties.PlateLookup
	afterPlates := test.Liquidhandler.FinalProperties.PlateLookup
	idMap := test.Liquidhandler.PlateIDMap()

	//check that idMap refers to things that exist
	for beforeID, afterID := range idMap {
		beforeObj, ok := beforePlates[beforeID]
		if !ok {
			return fmt.Errorf("idMap key \"%s\" doesn't exist in initial LHProperties.PlateLookup", beforeID)
		}
		afterObj, ok := afterPlates[afterID]
		if !ok {
			return fmt.Errorf("idMap value \"%s\" doesn't exist in final LHProperties.PlateLookup", afterID)
		}
		//check that you don't have tipboxes turning into plates, for example
		if beforeClass, afterClass := wtype.ClassOf(beforeObj), wtype.ClassOf(afterObj); beforeClass != afterClass {
			return fmt.Errorf("planner has turned a %s into a %s", beforeClass, afterClass)
		}
	}

	//check that everything in beforePlates is mapped to something
	for id, obj := range beforePlates {
		if _, ok := idMap[id]; !ok {
			return fmt.Errorf("%s with id %s exists in initial LHProperties, but isn't mapped to final LHProperties", wtype.ClassOf(obj), id)
		}
	}
	return nil
}

func (test *PlanningTest) checkPositionConsistency(lab *laboratory.Laboratory) error {
	for pos := range test.Liquidhandler.Properties.PosLookup {

		id1, ok1 := test.Liquidhandler.Properties.PosLookup[pos]
		id2, ok2 := test.Liquidhandler.FinalProperties.PosLookup[pos]

		if ok1 != ok2 {
			return fmt.Errorf("Position %s inconsistent: Before %t after %t", pos, ok1, ok2)
		}

		p1 := test.Liquidhandler.Properties.PlateLookup[id1]
		p2 := test.Liquidhandler.FinalProperties.PlateLookup[id2]

		// check types

		t1 := reflect.TypeOf(p1)
		t2 := reflect.TypeOf(p2)

		if t1 != t2 {
			return fmt.Errorf("Types of thing at position %s not same: %v %v", pos, t1, t2)
		}

		// ok nice we have some sort of consistency

		switch p1.(type) {
		case *wtype.Plate:
			pp1 := p1.(*wtype.Plate)
			pp2 := p2.(*wtype.Plate)
			if pp1.Type != pp2.Type {
				return fmt.Errorf("Plates at %s not same type: %s %s", pos, pp1.Type, pp2.Type)
			}

			for it := wtype.NewAddressIterator(pp1, wtype.ColumnWise, wtype.TopToBottom, wtype.LeftToRight, false); it.Valid(); it.Next() {
				if w1, w2 := pp1.Wellcoords[it.Curr().FormatA1()], pp2.Wellcoords[it.Curr().FormatA1()]; w1.IsEmpty(lab.IDGenerator) && w2.IsEmpty(lab.IDGenerator) {
					continue
				} else if w1.WContents.ID == w2.WContents.ID {
					fmt.Errorf("IDs before and after must differ: %v", w1.WContents.ID)
				}
			}
		case *wtype.LHTipbox:
			tb1 := p1.(*wtype.LHTipbox)
			tb2 := p2.(*wtype.LHTipbox)

			if tb1.Type != tb2.Type {
				return fmt.Errorf("Tipbox at changed type: %s %s", tb1.Type, tb2.Type)
			}
		case *wtype.LHTipwaste:
			tw1 := p1.(*wtype.LHTipwaste)
			tw2 := p2.(*wtype.LHTipwaste)

			if tw1.Type != tw2.Type {
				return fmt.Errorf("Tipwaste changed type: %s %s", tw1.Type, tw2.Type)
			}
		}

	}
	return nil
}

func (test *PlanningTest) expected(err error) bool {
	if err != nil && test.ErrorPrefix != "" {
		return strings.HasPrefix(err.Error(), test.ErrorPrefix)
	} else {
		return err == nil && test.ErrorPrefix == ""
	}
}

type PlanningTests []*PlanningTest

func (tests PlanningTests) Run(t *testing.T) {
	for _, test := range tests {
		test.Run(t)
	}
}

type Assertion func(*laboratory.Laboratory, *lh.Liquidhandler, *lh.LHRequest) error

type Assertions []Assertion

func (s Assertions) Assert(lab *laboratory.Laboratory, lh *lh.Liquidhandler, request *lh.LHRequest) error {
	errs := make(utils.ErrorSlice, len(s))
	for idx, assertion := range s {
		errs[idx] = assertion(lab, lh, request)
	}
	return errs.Pack()
}

// DebugPrintInstructions assertion that just prints all the generated instructions
func DebugPrintInstructions() Assertion { //nolint
	return func(lab *laboratory.Laboratory, lh *lh.Liquidhandler, rq *lh.LHRequest) error {
		for _, ins := range rq.Instructions {
			fmt.Println(liquidhandling.InsToString(ins))
		}
		return nil
	}
}

// NumberOfAssertion check that the number of instructions of the given type is
// equal to count
func NumberOfAssertion(iType *liquidhandling.InstructionType, count int) Assertion {
	return func(lab *laboratory.Laboratory, lh *lh.Liquidhandler, request *lh.LHRequest) error {
		var c int
		for _, ins := range request.Instructions {
			if ins.Type() == iType {
				c++
			}
		}
		if c != count {
			return fmt.Errorf("Expecting %d instrctions of type %s, got %d", count, iType, c)
		}
		return nil
	}
}

// TipsUsedAssertion check that the number of tips used is as expected
func TipsUsedAssertion(expected []wtype.TipEstimate) Assertion {
	return func(lab *laboratory.Laboratory, lh *lh.Liquidhandler, request *lh.LHRequest) error {
		if !reflect.DeepEqual(expected, request.TipsUsed) {
			return fmt.Errorf("Expected %v; Got %v", expected, request.TipsUsed)
		}
		return nil
	}
}

// InitialComponentAssertion check that the initial components are present in
// the given quantities.
// Currently only supports the case where each component name exists only once
func InitialComponentAssertion(expected map[string]float64) Assertion {
	return func(lab *laboratory.Laboratory, lh *lh.Liquidhandler, request *lh.LHRequest) error {
		for _, p := range lh.Properties.Plates {
			for _, w := range p.Wellcoords {
				if !w.IsEmpty(lab.IDGenerator) {
					if v, ok := expected[w.WContents.CName]; !ok {
						return fmt.Errorf("unexpected component in plating area: %s", w.WContents.CName)
					} else if v != w.WContents.Vol {
						return fmt.Errorf("volume of component %s was %v should be %v", w.WContents.CName, w.WContents.Vol, v)
					} else {
						delete(expected, w.WContents.CName)
					}
				}
			}
		}

		if len(expected) != 0 {
			return fmt.Errorf("unexpected components remaining: %v", expected)
		}
		return nil
	}
}

// InputLayoutAssertion check that the input layout is as expected
// expected is a map of well location (in A1 format) to liquid name for each input plate
func InputLayoutAssertion(expected ...map[string]string) Assertion {
	return func(lab *laboratory.Laboratory, lh *lh.Liquidhandler, request *lh.LHRequest) error {
		if len(request.InputPlateOrder) != len(expected) {
			return fmt.Errorf("input layout: expected %d input plates, got %d", len(expected), len(request.InputPlateOrder))
		}

		for plateNum, plateID := range request.InputPlateOrder {
			got := make(map[string]string)
			if plate, ok := request.InputPlates[plateID]; !ok {
				return fmt.Errorf("input layout: inconsistent InputPlateOrder in request: no id %q in liquidhandler", plateID)
			} else if plate == nil {
				return fmt.Errorf("input layout: nil input plate in request")
			} else {
				for address, well := range plate.Wellcoords {
					if !well.IsEmpty(lab.IDGenerator) {
						got[address] = well.Contents(lab.IDGenerator).CName
					}
				}
				if !reflect.DeepEqual(expected[plateNum], got) {
					return fmt.Errorf("input layout: input plate %d doesn't match:\ne: %v\ng: %v", plateNum, expected[plateNum], got)
				}
			}
		}
		return nil
	}
}

func describePlateVolumes(idGen *id.IDGenerator, order []string, plates map[string]*wtype.LHPlate) ([]map[string]float64, error) {
	ret := make([]map[string]float64, 0, len(order))
	for _, plateID := range order {
		got := make(map[string]float64)
		if plate, ok := plates[plateID]; !ok {
			return nil, errors.Errorf("inconsistent InputPlateOrder in request: no id %s", plateID)
		} else if plate == nil {
			return nil, errors.New("nil input plate in request")
		} else {
			for address, well := range plate.Wellcoords {
				if !well.IsEmpty(idGen) {
					got[address] = well.CurrentVolume(idGen).MustInStringUnit("ul").RawValue()
				}
			}
			ret = append(ret, got)
		}
	}
	return ret, nil
}

func keys(m map[string]float64) []string {
	ret := make([]string, 0, len(m))
	for k := range m {
		ret = append(ret, k)
	}
	sort.Strings(ret)
	return ret
}

func volumesMatch(tolerance float64, lhs, rhs map[string]float64) bool {
	if !reflect.DeepEqual(keys(lhs), keys(rhs)) {
		return false
	}

	for key, rVal := range rhs {
		if math.Abs(rVal-lhs[key]) > tolerance {
			return false
		}
	}

	return true
}

// InitialInputVolumesAssertion check that the input layout is as expected
// expected is a map of well location (in A1 format) to liquid to volume in ul
// tol is the maximum difference before an error is raised
func InitialInputVolumesAssertion(tol float64, expected ...map[string]float64) Assertion {
	return func(lab *laboratory.Laboratory, lh *lh.Liquidhandler, request *lh.LHRequest) error {

		if got, err := describePlateVolumes(lab.IDGenerator, request.InputPlateOrder, request.InputPlates); err != nil {
			return fmt.Errorf("initial input volumes: %v", err)
		} else {
			for i, g := range got {
				if !volumesMatch(tol, expected[i], g) {
					return fmt.Errorf("initial input volumes: input plate %d doesn't match:\ne: %v\ng: %v", i, expected[i], g)
				}
			}
		}
		return nil
	}
}

// FinalInputVolumesAssertion check that the input layout is as expected
// expected is a map of well location (in A1 format) to liquid to volume in ul
// tol is the maximum difference before an error is raised
func FinalInputVolumesAssertion(tol float64, expected ...map[string]float64) Assertion {
	return func(lab *laboratory.Laboratory, lh *lh.Liquidhandler, request *lh.LHRequest) error {

		pos := make([]string, 0, len(request.InputPlateOrder))
		m := lh.PlateIDMap()
		for _, in := range request.InputPlateOrder {
			pos = append(pos, lh.FinalProperties.PlateIDLookup[m[in]])
		}

		if got, err := describePlateVolumes(lab.IDGenerator, pos, lh.FinalProperties.Plates); err != nil {
			return fmt.Errorf("final input volumes %v", err)
		} else {
			for i, g := range got {
				if !volumesMatch(tol, expected[i], g) {
					return fmt.Errorf("final input volumes: input plate %d doesn't match:\ne: %v\ng: %v", i, expected[i], g)
				}
			}
		}
		return nil
	}
}

// FinalOutputVolumesAssertion check that the output layout is as expected
// expected is a map of well location (in A1 format) to liquid to volume in ul
// tol is the maximum difference before an error is raised
func FinalOutputVolumesAssertion(tol float64, expected ...map[string]float64) Assertion {
	return func(lab *laboratory.Laboratory, lh *lh.Liquidhandler, request *lh.LHRequest) error {

		pos := make([]string, 0, len(request.OutputPlateOrder))
		m := lh.PlateIDMap()
		for _, in := range request.OutputPlateOrder {
			pos = append(pos, lh.FinalProperties.PlateIDLookup[m[in]])
		}

		if got, err := describePlateVolumes(lab.IDGenerator, pos, lh.FinalProperties.Plates); err != nil {
			return fmt.Errorf("while asserting final output volumes: %v", err)
		} else {
			for i, g := range got {
				if !volumesMatch(tol, expected[i], g) {
					return fmt.Errorf("output plate %d doesn't match:\ne: %v\ng: %v", i, expected[i], g)
				}
			}
		}
		return nil
	}
}

// LayoutSummaryAssertion assert that the generated layout is the same as the one found at the given path
func LayoutSummaryAssertion(path string) Assertion { //nolint
	return func(lab *laboratory.Laboratory, handler *lh.Liquidhandler, rq *lh.LHRequest) error {
		if expected, err := lab.FileManager.ReadAll(wtype.NewFile(path).AsInput()); err != nil {
			return err
		} else if got, err := lh.SummarizeLayout(lab.IDGenerator, handler.Properties, handler.FinalProperties, handler.PlateIDMap()); err != nil {
			return err
		} else if err := AssertLayoutsEquivalent(got, expected); err != nil {
			return fmt.Errorf("layout summary mismatch: %v", err)
		}
		return nil
	}
}

// ActionsSummaryAssertion assert that the generated layout is the same as the one found at the given path
func ActionsSummaryAssertion(path string) Assertion { //nolint
	return func(lab *laboratory.Laboratory, handler *lh.Liquidhandler, rq *lh.LHRequest) error {
		if expected, err := lab.FileManager.ReadAll(wtype.NewFile(path).AsInput()); err != nil {
			return err
		} else if got, err := lh.SummarizeActions(lab.IDGenerator, handler.Properties, rq.InstructionTree); err != nil {
			return err
		} else if err := AssertActionsEquivalent(got, expected); err != nil {
			return fmt.Errorf("actions summary mismatch: %v", err)
		}
		return nil
	}
}

type TestMixComponent struct {
	LiquidName    string
	VolumesByWell map[string]float64
	LiquidType    wtype.LiquidType
	Sampler       func(*laboratory.Laboratory, *wtype.Liquid, wunit.Volume) *wtype.Liquid
}

func (self TestMixComponent) AddSamples(lab *laboratory.Laboratory, samples map[string][]*wtype.Liquid) {
	var totalVolume float64
	for _, v := range self.VolumesByWell {
		totalVolume += v
	}

	source := GetComponentForTest(lab, self.LiquidName, wunit.NewVolume(totalVolume, "ul"))
	if self.LiquidType != "" {
		source.Type = self.LiquidType
	}

	for well, vol := range self.VolumesByWell {
		samples[well] = append(samples[well], self.Sampler(lab, source, wunit.NewVolume(vol, "ul")))
	}
}

func (self TestMixComponent) AddToPlate(lab *laboratory.Laboratory, plate *wtype.LHPlate) {
	for well, vol := range self.VolumesByWell {
		cmp := GetComponentForTest(lab, self.LiquidName, wunit.NewVolume(vol, "ul"))
		if self.LiquidType != "" {
			cmp.Type = self.LiquidType
		}

		if err := plate.Wellcoords[well].SetContents(lab.IDGenerator, cmp); err != nil {
			panic(err)
		}
	}
}

type TestMixComponents []TestMixComponent

func (self TestMixComponents) AddSamples(lab *laboratory.Laboratory, samples map[string][]*wtype.Liquid) {
	for _, to := range self {
		to.AddSamples(lab, samples)
	}
}

type TestInputLayout []TestMixComponent

func (self TestInputLayout) AddToPlate(lab *laboratory.Laboratory, plate *wtype.LHPlate) {
	for _, to := range self {
		to.AddToPlate(lab, plate)
	}
}

type InstructionBuilder func(*laboratory.Laboratory) ([]*wtype.LHInstruction, error)

func Mixes(outputPlateType wtype.PlateTypeName, components TestMixComponents) InstructionBuilder {
	return func(lab *laboratory.Laboratory) ([]*wtype.LHInstruction, error) {

		samplesByWell := make(map[string][]*wtype.Liquid)
		components.AddSamples(lab, samplesByWell)

		ret := make([]*wtype.LHInstruction, 0, len(samplesByWell))

		for well, samples := range samplesByWell {
			ret = append(ret, mixer.GenericMix(lab, mixer.MixOptions{
				Inputs:    samples,
				PlateType: outputPlateType,
				Address:   well,
				PlateName: "outputplate",
			}))
		}

		return ret, nil
	}
}

func ColumnWise(rows int, volumes []float64) map[string]float64 {
	ret := make(map[string]float64, len(volumes))
	for i, v := range volumes {
		ret[wtype.WellCoords{X: i / rows, Y: i % rows}.FormatA1()] = v
	}
	return ret
}
