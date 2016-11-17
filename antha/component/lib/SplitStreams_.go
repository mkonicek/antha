package lib

import

// encodes a protocol for reformatting into two output streams
(
	"github.com/antha-lang/antha/antha/anthalib/mixer"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/component"
	"github.com/antha-lang/antha/execute"
	"github.com/antha-lang/antha/inject"
	"golang.org/x/net/context"
)

func _SplitStreamsRequirements() {
}

func _SplitStreamsSetup(_ctx context.Context, _input *SplitStreamsInput) {
}

func _SplitStreamsSteps(_ctx context.Context, _input *SplitStreamsInput, _output *SplitStreamsOutput) {
	s1 := mixer.Sample(_input.InputStream, _input.Stream1Vol)
	s1.CName = s1.CName + "_Stream1"
	_output.Stream1 = execute.MixTo(_ctx, _input.Stream1PlateType, "", 1, s1)

	// ensure we end up with samples on different plates if types are the same
	//pt2:=1
	//if Stream1PlateType==Stream2PlateType{
	pt2 := 2
	//}
	s2 := mixer.Sample(_input.InputStream, _input.Stream2Vol)
	s2.CName = s2.CName + "_Stream2"
	_output.Stream2 = execute.MixTo(_ctx, _input.Stream2PlateType, "", pt2, s2)
	_output.Stream2.CName = _output.Stream2.CName + "_Stream2"
}

func _SplitStreamsAnalysis(_ctx context.Context, _input *SplitStreamsInput, _output *SplitStreamsOutput) {
}

func _SplitStreamsValidation(_ctx context.Context, _input *SplitStreamsInput, _output *SplitStreamsOutput) {
}
func _SplitStreamsRun(_ctx context.Context, input *SplitStreamsInput) *SplitStreamsOutput {
	output := &SplitStreamsOutput{}
	_SplitStreamsSetup(_ctx, input)
	_SplitStreamsSteps(_ctx, input, output)
	_SplitStreamsAnalysis(_ctx, input, output)
	_SplitStreamsValidation(_ctx, input, output)
	return output
}

func SplitStreamsRunSteps(_ctx context.Context, input *SplitStreamsInput) *SplitStreamsSOutput {
	soutput := &SplitStreamsSOutput{}
	output := _SplitStreamsRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func SplitStreamsNew() interface{} {
	return &SplitStreamsElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &SplitStreamsInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _SplitStreamsRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &SplitStreamsInput{},
			Out: &SplitStreamsOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wunit.Make_units
)

type SplitStreamsElement struct {
	inject.CheckedRunner
}

type SplitStreamsInput struct {
	InputStream      *wtype.LHComponent
	Stream1PlateType string
	Stream1Vol       wunit.Volume
	Stream2PlateType string
	Stream2Vol       wunit.Volume
}

type SplitStreamsOutput struct {
	Stream1 *wtype.LHComponent
	Stream2 *wtype.LHComponent
}

type SplitStreamsSOutput struct {
	Data struct {
	}
	Outputs struct {
		Stream1 *wtype.LHComponent
		Stream2 *wtype.LHComponent
	}
}

func init() {
	if err := addComponent(component.Component{Name: "SplitStreams",
		Constructor: SplitStreamsNew,
		Desc: component.ComponentDesc{
			Desc: "",
			Path: "src/github.com/antha-lang/antha/antha/component/an/GrowthAndAssay/splitstreams.an",
			Params: []component.ParamDesc{
				{Name: "InputStream", Desc: "", Kind: "Inputs"},
				{Name: "Stream1PlateType", Desc: "", Kind: "Parameters"},
				{Name: "Stream1Vol", Desc: "", Kind: "Parameters"},
				{Name: "Stream2PlateType", Desc: "", Kind: "Parameters"},
				{Name: "Stream2Vol", Desc: "", Kind: "Parameters"},
				{Name: "Stream1", Desc: "", Kind: "Outputs"},
				{Name: "Stream2", Desc: "", Kind: "Outputs"},
			},
		},
	}); err != nil {
		panic(err)
	}
}
