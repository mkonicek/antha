package liquidhandling

import (
	"context"
	"fmt"
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
