package liquidhandling

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"math"
	"reflect"
	"sort"
	"strings"
	"testing"

	"github.com/antha-lang/antha/antha/anthalib/mixer"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/microArch/driver/liquidhandling"
)

type PlanningTest struct {
	Name          string
	Liquidhandler *Liquidhandler
	Instructions  InstructionBuilder
	InputPlates   []*wtype.LHPlate
	OutputPlates  []*wtype.LHPlate
	ErrorPrefix   string
	Assertions    Assertions
}

func (test *PlanningTest) Run(ctx context.Context, t *testing.T) {
	t.Run(test.Name, func(t *testing.T) {
		test.run(ctx, t)
	})
}

func (test *PlanningTest) run(ctx context.Context, t *testing.T) {
	request := NewLHRequest()
	for _, ins := range test.Instructions(ctx) {
		request.Add_instruction(ins)
	}

	for _, plate := range test.InputPlates {
		if !plate.IsEmpty() {
			request.AddUserPlate(plate)
			plate = plate.Dup()
			plate.Clean()
		}
		request.InputPlatetypes = append(request.InputPlatetypes, plate)
	}

	request.OutputPlatetypes = append(request.OutputPlatetypes, test.OutputPlates...)

	if test.Liquidhandler == nil {
		test.Liquidhandler = GetLiquidHandlerForTest(ctx)
	}

	if err := test.Liquidhandler.Plan(ctx, request); !test.expected(err) {
		t.Fatalf("expecting error = %q: got error %q", test.ErrorPrefix, err.Error())
	}

	test.Assertions.Assert(t, test.Liquidhandler, request)

	if !t.Failed() && test.ErrorPrefix == "" {
		test.checkPlateIDMap(t)
		test.checkPositionConsistency(t)
	}
}

func (test *PlanningTest) checkPlateIDMap(t *testing.T) {
	beforePlates := test.Liquidhandler.Properties.PlateLookup
	afterPlates := test.Liquidhandler.FinalProperties.PlateLookup
	idMap := test.Liquidhandler.PlateIDMap()

	//check that idMap refers to things that exist
	for beforeID, afterID := range idMap {
		beforeObj, ok := beforePlates[beforeID]
		if !ok {
			t.Errorf("idMap key \"%s\" doesn't exist in initial LHProperties.PlateLookup", beforeID)
			continue
		}
		afterObj, ok := afterPlates[afterID]
		if !ok {
			t.Errorf("idMap value \"%s\" doesn't exist in final LHProperties.PlateLookup", afterID)
			continue
		}
		//check that you don't have tipboxes turning into plates, for example
		if beforeClass, afterClass := wtype.ClassOf(beforeObj), wtype.ClassOf(afterObj); beforeClass != afterClass {
			t.Errorf("planner has turned a %s into a %s", beforeClass, afterClass)
		}
	}

	//check that everything in beforePlates is mapped to something
	for id, obj := range beforePlates {
		if _, ok := idMap[id]; !ok {
			t.Errorf("%s with id %s exists in initial LHProperties, but isn't mapped to final LHProperties", wtype.ClassOf(obj), id)
		}
	}
}

func (test *PlanningTest) checkPositionConsistency(t *testing.T) {
	for pos := range test.Liquidhandler.Properties.PosLookup {

		id1, ok1 := test.Liquidhandler.Properties.PosLookup[pos]
		id2, ok2 := test.Liquidhandler.FinalProperties.PosLookup[pos]

		if ok1 != ok2 {
			t.Fatal(fmt.Sprintf("Position %s inconsistent: Before %t after %t", pos, ok1, ok2))
		}

		p1 := test.Liquidhandler.Properties.PlateLookup[id1]
		p2 := test.Liquidhandler.FinalProperties.PlateLookup[id2]

		// check types

		t1 := reflect.TypeOf(p1)
		t2 := reflect.TypeOf(p2)

		if t1 != t2 {
			t.Fatal(fmt.Sprintf("Types of thing at position %s not same: %v %v", pos, t1, t2))
		}

		// ok nice we have some sort of consistency

		switch p1.(type) {
		case *wtype.Plate:
			pp1 := p1.(*wtype.Plate)
			pp2 := p2.(*wtype.Plate)
			if pp1.Type != pp2.Type {
				t.Fatal(fmt.Sprintf("Plates at %s not same type: %s %s", pos, pp1.Type, pp2.Type))
			}

			for it := wtype.NewAddressIterator(pp1, wtype.ColumnWise, wtype.TopToBottom, wtype.LeftToRight, false); it.Valid(); it.Next() {
				if w1, w2 := pp1.Wellcoords[it.Curr().FormatA1()], pp2.Wellcoords[it.Curr().FormatA1()]; w1.IsEmpty() && w2.IsEmpty() {
					continue
				} else if w1.WContents.ID == w2.WContents.ID {
					t.Fatal("IDs before and after must differ")
				}
			}
		case *wtype.LHTipbox:
			tb1 := p1.(*wtype.LHTipbox)
			tb2 := p2.(*wtype.LHTipbox)

			if tb1.Type != tb2.Type {
				t.Fatal(fmt.Sprintf("Tipbox at changed type: %s %s", tb1.Type, tb2.Type))
			}
		case *wtype.LHTipwaste:
			tw1 := p1.(*wtype.LHTipwaste)
			tw2 := p2.(*wtype.LHTipwaste)

			if tw1.Type != tw2.Type {
				t.Fatal(fmt.Sprintf("Tipwaste changed type: %s %s", tw1.Type, tw2.Type))
			}
		}

	}

}

func (test *PlanningTest) expected(err error) bool {
	if err != nil && test.ErrorPrefix != "" {
		return strings.HasPrefix(err.Error(), test.ErrorPrefix)
	} else if err == nil && test.ErrorPrefix == "" {
		return true
	}
	return false
}

type PlanningTests []*PlanningTest

func (tests PlanningTests) Run(ctx context.Context, t *testing.T) {
	for _, test := range tests {
		test.Run(ctx, t)
	}
}

type Assertion func(*testing.T, *Liquidhandler, *LHRequest)

type Assertions []Assertion

func (s Assertions) Assert(t *testing.T, lh *Liquidhandler, request *LHRequest) {
	for _, assertion := range s {
		assertion(t, lh, request)
	}
}

// AssertNumberOf check that the number of instructions of the given type is
// equal to count
func AssertNumberOf(iType *liquidhandling.InstructionType, count int) Assertion {
	return func(t *testing.T, lh *Liquidhandler, request *LHRequest) {
		var c int
		for _, ins := range request.Instructions {
			if ins.Type() == iType {
				c++
			}
		}
		if c != count {
			t.Errorf("Expecting %d instrctions of type %s, got %d", count, iType, c)
		}
	}
}

// AssertTipsUsed check that the number of tips used is as expected
func AssertTipsUsed(expected []wtype.TipEstimate) Assertion {
	return func(t *testing.T, lh *Liquidhandler, request *LHRequest) {
		if !reflect.DeepEqual(expected, request.TipsUsed) {
			t.Errorf("Expected %v Got %v", expected, request.TipsUsed)
		}
	}
}

// AssertInputLayout check that the input layout is as expected
// expected is a map of well location (in A1 format) to liquid name for each input plate
func AssertInputLayout(expected ...map[string]string) Assertion {
	return func(t *testing.T, lh *Liquidhandler, request *LHRequest) {
		if len(request.InputPlateOrder) != len(expected) {
			t.Errorf("expected %d input plates, got %d", len(expected), len(request.InputPlateOrder))
			return
		}

		for plateNum, plateID := range request.InputPlateOrder {
			got := make(map[string]string)
			if plate, ok := request.InputPlates[plateID]; !ok {
				t.Errorf("while asserting input layout: inconsistent InputPlateOrder in request: no id %q in liquidhandler", plateID)
				continue
			} else if plate == nil {
				t.Errorf("nil input plate in request")
				continue
			} else {
				for address, well := range plate.Wellcoords {
					if !well.IsEmpty() {
						got[address] = well.Contents().CName
					}
				}
			}

			if !reflect.DeepEqual(expected[plateNum], got) {
				t.Errorf("input plate %d doesn't match:\ne: %v\ng: %v", plateNum, expected[plateNum], got)
			}
		}
	}
}

func describePlateVolumes(order []string, plates map[string]*wtype.LHPlate) ([]map[string]float64, error) {
	ret := make([]map[string]float64, 0, len(order))
	for _, plateID := range order {
		got := make(map[string]float64)
		if plate, ok := plates[plateID]; !ok {
			return nil, errors.Errorf("inconsistent InputPlateOrder in request: no id %s", plateID)
		} else if plate == nil {
			return nil, errors.New("nil input plate in request")
		} else {
			for address, well := range plate.Wellcoords {
				if !well.IsEmpty() {
					got[address] = well.CurrentVolume().MustInStringUnit("ul").RawValue()
				}
			}
		}

		ret = append(ret, got)
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

// AssertInitialInputVolumes check that the input layout is as expected
// expected is a map of well location (in A1 format) to liquid to volume in ul
// tol is the maximum difference before an error is raised
func AssertInitialInputVolumes(tol float64, expected ...map[string]float64) Assertion {
	return func(t *testing.T, lh *Liquidhandler, request *LHRequest) {

		if got, err := describePlateVolumes(request.InputPlateOrder, request.InputPlates); err != nil {
			t.Error(errors.WithMessage(err, "while asserting initial input volumes"))
		} else {
			for i, g := range got {
				if !volumesMatch(tol, expected[i], g) {
					t.Errorf("input plate %d doesn't match:\ne: %v\ng: %v", i, expected[i], g)
				}
			}
		}
	}
}

// AssertFinalInputVolumes check that the input layout is as expected
// expected is a map of well location (in A1 format) to liquid to volume in ul
// tol is the maximum difference before an error is raised
func AssertFinalInputVolumes(tol float64, expected ...map[string]float64) Assertion {
	return func(t *testing.T, lh *Liquidhandler, request *LHRequest) {

		pos := make([]string, 0, len(request.InputPlateOrder))
		for _, in := range request.InputPlateOrder {
			pos = append(pos, lh.FinalProperties.PlateIDLookup[lh.plateIDMap[in]])
		}

		if got, err := describePlateVolumes(pos, lh.FinalProperties.Plates); err != nil {
			t.Error(errors.WithMessage(err, "while asserting final input volumes"))
		} else {
			for i, g := range got {
				if !volumesMatch(tol, expected[i], g) {
					t.Errorf("input plate %d doesn't match:\ne: %v\ng: %v", i, expected[i], g)
				}
			}
		}
	}
}

// AssertFinalOutputVolumes check that the output layout is as expected
// expected is a map of well location (in A1 format) to liquid to volume in ul
// tol is the maximum difference before an error is raised
func AssertFinalOutputVolumes(tol float64, expected ...map[string]float64) Assertion {
	return func(t *testing.T, lh *Liquidhandler, request *LHRequest) {

		pos := make([]string, 0, len(request.OutputPlateOrder))
		for _, in := range request.OutputPlateOrder {
			pos = append(pos, lh.FinalProperties.PlateIDLookup[lh.plateIDMap[in]])
		}

		if got, err := describePlateVolumes(pos, lh.FinalProperties.Plates); err != nil {
			t.Error(errors.WithMessage(err, "while asserting final output volumes"))
		} else {
			for i, g := range got {
				if !volumesMatch(tol, expected[i], g) {
					t.Errorf("output plate %d doesn't match:\ne: %v\ng: %v", i, expected[i], g)
				}
			}
		}
	}
}

type TestMixComponent struct {
	LiquidName    string
	VolumesByWell map[string]float64
	LiquidType    wtype.LiquidType
	Sampler       func(*wtype.Liquid, wunit.Volume) *wtype.Liquid
}

func (self TestMixComponent) AddSamples(ctx context.Context, samples map[string][]*wtype.Liquid) {
	var totalVolume float64
	for _, v := range self.VolumesByWell {
		totalVolume += v
	}

	source := GetComponentForTest(ctx, self.LiquidName, wunit.NewVolume(totalVolume, "ul"))
	if self.LiquidType != "" {
		source.Type = self.LiquidType
	}

	for well, vol := range self.VolumesByWell {
		samples[well] = append(samples[well], self.Sampler(source, wunit.NewVolume(vol, "ul")))
	}
}

func (self TestMixComponent) AddToPlate(ctx context.Context, plate *wtype.LHPlate) {

	for well, vol := range self.VolumesByWell {
		cmp := GetComponentForTest(ctx, self.LiquidName, wunit.NewVolume(vol, "ul"))
		if self.LiquidType != "" {
			cmp.Type = self.LiquidType
		}

		plate.Wellcoords[well].SetContents(cmp)
	}
}

type TestMixComponents []TestMixComponent

func (self TestMixComponents) AddSamples(ctx context.Context, samples map[string][]*wtype.Liquid) {
	for _, to := range self {
		to.AddSamples(ctx, samples)
	}
}

type TestInputLayout []TestMixComponent

func (self TestInputLayout) AddToPlate(ctx context.Context, plate *wtype.LHPlate) {
	for _, to := range self {
		to.AddToPlate(ctx, plate)
	}
}

type InstructionBuilder func(context.Context) []*wtype.LHInstruction

func Mixes(outputPlateType string, components TestMixComponents) InstructionBuilder {
	return func(ctx context.Context) []*wtype.LHInstruction {

		samplesByWell := make(map[string][]*wtype.Liquid)
		components.AddSamples(ctx, samplesByWell)

		ret := make([]*wtype.LHInstruction, 0, len(samplesByWell))

		for well, samples := range samplesByWell {
			ret = append(ret, mixer.GenericMix(mixer.MixOptions{
				Inputs:    samples,
				PlateType: outputPlateType,
				Address:   well,
				PlateName: "outputplate",
			}))
		}

		return ret
	}
}

func ColumnWise(rows int, volumes []float64) map[string]float64 {
	ret := make(map[string]float64, len(volumes))
	for i, v := range volumes {
		ret[wtype.WellCoords{X: i / rows, Y: i % rows}.FormatA1()] = v
	}
	return ret
}
