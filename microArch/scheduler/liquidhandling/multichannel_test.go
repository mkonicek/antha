package liquidhandling

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"github.com/antha-lang/antha/antha/anthalib/mixer"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/microArch/driver/liquidhandling"
)

type MultiChannelTest struct {
	Name           string
	Liquidhandler  *Liquidhandler
	Instructions   InstructionBuilder
	InputPlates    []*wtype.LHPlate
	OutputPlates   []*wtype.LHPlate
	ExpectingError bool
	Assertions     Assertions
}

func (test *MultiChannelTest) Run(ctx context.Context, t *testing.T) {
	t.Run(test.Name, func(t *testing.T) {
		test.run(ctx, t)
	})
}

func (test *MultiChannelTest) run(ctx context.Context, t *testing.T) {
	request := NewLHRequest()
	for _, ins := range test.Instructions(ctx) {
		request.Add_instruction(ins)
	}

	request.InputPlatetypes = append(request.InputPlatetypes, test.InputPlates...)
	request.OutputPlatetypes = append(request.OutputPlatetypes, test.OutputPlates...)

	if test.Liquidhandler == nil {
		test.Liquidhandler = GetLiquidHandlerForTest(ctx)
	}

	if err := test.Liquidhandler.Plan(ctx, request); !test.expected(err) {
		t.Fatalf("expecting error = %t: got error %v", test.ExpectingError, err)
	}

	test.Assertions.Assert(t, request)

	if t.Failed() {
		fmt.Println("Generated Instructions:")
		for i, ins := range request.Instructions {
			fmt.Printf("  %02d: %v\n", i, liquidhandling.InsToString(ins))
		}
	} else if !test.ExpectingError {
		test.checkPlateIDMap(t)
	}
}

func (test *MultiChannelTest) checkPlateIDMap(t *testing.T) {
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

func (test *MultiChannelTest) checkPositionConsistency(t *testing.T) {
	for pos := range test.Liquidhandler.Properties.PosLookup {

		id1, ok1 := test.Liquidhandler.Properties.PosLookup[pos]
		id2, ok2 := test.Liquidhandler.FinalProperties.PosLookup[pos]

		if ok1 && !ok2 || ok2 && !ok1 {
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
			it := wtype.NewAddressIterator(pp1, wtype.ColumnWise, wtype.TopToBottom, wtype.LeftToRight, false)

			for {
				if !it.Valid() {
					break
				}
				wc := it.Curr()
				w1 := pp1.Wellcoords[wc.FormatA1()]
				w2 := pp2.Wellcoords[wc.FormatA1()]

				if w1.IsEmpty() && w2.IsEmpty() {
					it.Next()
					continue
				}

				if w1.WContents.ID == w2.WContents.ID {
					t.Fatal(fmt.Sprintf("IDs before and after must differ"))
				}
				it.Next()
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

func (test *MultiChannelTest) expected(err error) bool {
	return (err != nil) == test.ExpectingError
}

type MultiChannelTests []*MultiChannelTest

func (tests MultiChannelTests) Run(ctx context.Context, t *testing.T) {
	for _, test := range tests {
		test.Run(ctx, t)
	}
}

type Assertion func(*testing.T, *LHRequest)

type Assertions []Assertion

func (s Assertions) Assert(t *testing.T, request *LHRequest) {
	for _, assertion := range s {
		assertion(t, request)
	}
}

// AssertNumberOf check that the number of instructions of the given type is
// equal to count
func AssertNumberOf(iType *liquidhandling.InstructionType, count int) Assertion {
	return func(t *testing.T, request *LHRequest) {
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

type TestMixComponents []TestMixComponent

func (self TestMixComponents) AddSamples(ctx context.Context, samples map[string][]*wtype.Liquid) {
	for _, to := range self {
		to.AddSamples(ctx, samples)
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
