package liquidhandling

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"testing"

	"github.com/antha-lang/antha/antha/anthalib/mixer"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/inventory"
	"github.com/antha-lang/antha/microArch/driver/liquidhandling"
)

type MultiChannelTest struct {
	Name             string
	Liquidhandler    *Liquidhandler
	Instructions     InstructionBuilder
	InputPlateTypes  []string
	OutputPlateTypes []string
	ExpectingError   bool
	Assertions       Assertions
}

func (test *MultiChannelTest) Run(t *testing.T) {
	t.Run(test.Name, test.run)
}

func (test *MultiChannelTest) run(t *testing.T) {
	ctx := GetContextForTest()

	request := NewLHRequest()
	for _, ins := range test.Instructions(ctx) {
		request.Add_instruction(ins)
	}

	if input, err := test.makePlates(ctx, test.InputPlateTypes); err != nil {
		t.Fatal(errors.WithMessage(err, "while making input plates"))
	} else {
		request.InputPlatetypes = append(request.InputPlatetypes, input...)
	}

	if output, err := test.makePlates(ctx, test.OutputPlateTypes); err != nil {
		t.Fatal(errors.WithMessage(err, "while making output plates"))
	} else {
		request.OutputPlatetypes = append(request.OutputPlatetypes, output...)
	}

	lh := test.Liquidhandler
	if lh == nil {
		lh = GetLiquidHandlerForTest(ctx)
	}

	if err := lh.Plan(ctx, request); !test.expected(err) {
		t.Fatalf("expecting error = %t: got error %v", test.ExpectingError, err)
	}

	test.Assertions.Assert(t, request)

	if t.Failed() {
		fmt.Println("Generated Instructions:")
		for i, ins := range request.Instructions {
			fmt.Printf("  %02d: %v\n", i, liquidhandling.InsToString(ins))
		}
	}
}

func (test *MultiChannelTest) makePlates(ctx context.Context, plateTypes []string) ([]*wtype.LHPlate, error) {
	ret := make([]*wtype.LHPlate, 0, len(plateTypes))
	for _, plateType := range plateTypes {
		if p, err := inventory.NewPlate(ctx, plateType); err != nil {
			return nil, err
		} else {
			ret = append(ret, p)
		}
	}
	return ret, nil
}

func (test *MultiChannelTest) expected(err error) bool {
	return (err != nil) == test.ExpectingError
}

type MultiChannelTests []*MultiChannelTest

func (tests MultiChannelTests) Run(t *testing.T) {
	for _, test := range tests {
		test.Run(t)
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

type TestOutput struct {
	LiquidName    string
	VolumesByWell map[string]float64
	LiquidType    wtype.LiquidType
}

func (self TestOutput) AddSamples(ctx context.Context, samples map[string][]*wtype.Liquid) {
	var totalVolume float64
	for _, v := range self.VolumesByWell {
		totalVolume += v
	}

	source := GetComponentForTest(ctx, self.LiquidName, wunit.NewVolume(totalVolume, "ul"))
	if self.LiquidType != "" {
		source.Type = self.LiquidType
	}

	for well, vol := range self.VolumesByWell {
		samples[well] = append(samples[well], mixer.Sample(source, wunit.NewVolume(vol, "ul")))
	}
}

type TestOutputs []TestOutput

func (self TestOutputs) AddSamples(ctx context.Context, samples map[string][]*wtype.Liquid) {
	for _, to := range self {
		to.AddSamples(ctx, samples)
	}
}

type InstructionBuilder func(context.Context) []*wtype.LHInstruction

func Mixes(outputPlateType string, components TestOutputs) InstructionBuilder {
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

func TestMultichannel1(t *testing.T) {

	MultiChannelTests{
		{
			Name: "single channel",
			Instructions: Mixes("pcrplate_skirted_riser18", TestOutputs{
				{
					LiquidName: "water",
					VolumesByWell: map[string]float64{
						"A1": 8.0,
						"B1": 8.0,
						"C1": 8.0,
						"D1": 8.0,
						"E1": 8.0,
						"F1": 8.0,
						"G1": 8.0,
						"H1": 8.0,
					},
					LiquidType: wtype.LTSingleChannel,
				},
				{
					LiquidName: "mastermix_sapI",
					VolumesByWell: map[string]float64{
						"A1": 8.0,
						"B1": 8.0,
						"C1": 8.0,
						"D1": 8.0,
						"E1": 8.0,
						"F1": 8.0,
						"G1": 8.0,
						"H1": 8.0,
					},
					LiquidType: wtype.LTSingleChannel,
				},
				{
					LiquidName: "dna",
					VolumesByWell: map[string]float64{
						"A1": 1.0,
						"B1": 1.0,
						"C1": 1.0,
						"D1": 1.0,
						"E1": 1.0,
						"F1": 1.0,
						"G1": 1.0,
						"H1": 1.0,
					},
					LiquidType: wtype.LTSingleChannel,
				},
			}),
			InputPlateTypes:  []string{"DWST12"},
			OutputPlateTypes: []string{"pcrplate_skirted_riser18"},
			Assertions: Assertions{
				AssertNumberOf(liquidhandling.ASP, 3*8), //no multichanneling
				AssertNumberOf(liquidhandling.DSP, 3*8), //no multichanneling
			},
		},
		{
			Name: "multi channel",
			Instructions: Mixes("pcrplate_skirted_riser18", TestOutputs{
				{
					LiquidName: "water",
					VolumesByWell: map[string]float64{
						"A1": 8.0,
						"B1": 8.0,
						"C1": 8.0,
						"D1": 8.0,
						"E1": 8.0,
						"F1": 8.0,
						"G1": 8.0,
						"H1": 8.0,
					},
					LiquidType: wtype.LTWater,
				},
				{
					LiquidName: "mastermix_sapI",
					VolumesByWell: map[string]float64{
						"A1": 8.0,
						"B1": 8.0,
						"C1": 8.0,
						"D1": 8.0,
						"E1": 8.0,
						"F1": 8.0,
						"G1": 8.0,
						"H1": 8.0,
					},
					LiquidType: wtype.LTWater,
				},
				{
					LiquidName: "dna",
					VolumesByWell: map[string]float64{
						"A1": 1.0,
						"B1": 1.0,
						"C1": 1.0,
						"D1": 1.0,
						"E1": 1.0,
						"F1": 1.0,
						"G1": 1.0,
						"H1": 1.0,
					},
					LiquidType: wtype.LTWater,
				},
			}),
			InputPlateTypes:  []string{"DWST12"},
			OutputPlateTypes: []string{"pcrplate_skirted_riser18"},
			Assertions: Assertions{
				AssertNumberOf(liquidhandling.ASP, 3), //full multichanneling
				AssertNumberOf(liquidhandling.DSP, 3), //full multichanneling
			},
		},
	}.Run(t)
}
