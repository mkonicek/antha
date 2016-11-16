package lib

import (
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/text"
	"github.com/antha-lang/antha/antha/anthalib/mixer"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/component"
	"github.com/antha-lang/antha/execute"
	"github.com/antha-lang/antha/inject"
	"golang.org/x/net/context"
)

// Input parameters for this protocol (data)

//	StockReConcinUperml 		[]int
//	DesiredConcinUperml	 		[]int

//OutputReactionName			string

// Physical Inputs to this protocol with types

// Physical outputs from this protocol with types

// Data which is returned from this protocol, and data types

func _RestrictionDigestionRequirements() {}

// Conditions to run on startup
func _RestrictionDigestionSetup(_ctx context.Context, _input *RestrictionDigestionInput) {}

// The core process for this protocol, with the steps to be performed
// for every input
func _RestrictionDigestionSteps(_ctx context.Context, _input *RestrictionDigestionInput, _output *RestrictionDigestionOutput) {
	samples := make([]*wtype.LHComponent, 0)
	waterSample := mixer.SampleForTotalVolume(_input.Water, _input.ReactionVolume)
	samples = append(samples, waterSample)

	bufferSample := mixer.Sample(_input.Buffer, _input.BufferVol)
	samples = append(samples, bufferSample)

	if _input.BSAvol.Mvalue != 0 {
		bsaSample := mixer.Sample(_input.BSAoptional, _input.BSAvol)
		samples = append(samples, bsaSample)
	}

	// change to fixing concentration(or mass) of dna per reaction
	_input.DNASolution.CName = _input.DNAName
	dnaSample := mixer.Sample(_input.DNASolution, _input.DNAVol)
	samples = append(samples, dnaSample)

	for k, enzyme := range _input.EnzSolutions {

		// work out volume to add in L

		// e.g. 1 U / (10000 * 1000) * 0.000002
		//volinL := DesiredUinreaction/(StockReConcinUperml*1000) * ReactionVolume.SIValue()
		//volumetoadd := wunit.NewVolume(volinL,"L")
		enzyme.CName = _input.EnzymeNames[k]
		text.Print("adding enzyme"+_input.EnzymeNames[k], "to"+_input.DNAName)
		enzSample := mixer.Sample(enzyme, _input.EnzVolumestoadd[k])
		enzSample.CName = _input.EnzymeNames[k]
		samples = append(samples, enzSample)
	}

	// incubate the reaction mixture
	r1 := execute.Incubate(_ctx, execute.MixInto(_ctx, _input.OutPlate, "", samples...), _input.ReactionTemp, _input.ReactionTime, false)
	// inactivate
	_output.Reaction = execute.Incubate(_ctx, r1, _input.InactivationTemp, _input.InactivationTime, false)
}

// Run after controls and a steps block are completed to
// post process any data and provide downstream results
func _RestrictionDigestionAnalysis(_ctx context.Context, _input *RestrictionDigestionInput, _output *RestrictionDigestionOutput) {
}

// A block of tests to perform to validate that the sample was processed correctly
// Optionally, destructive tests can be performed to validate results on a
// dipstick basis
func _RestrictionDigestionValidation(_ctx context.Context, _input *RestrictionDigestionInput, _output *RestrictionDigestionOutput) {
}
func _RestrictionDigestionRun(_ctx context.Context, input *RestrictionDigestionInput) *RestrictionDigestionOutput {
	output := &RestrictionDigestionOutput{}
	_RestrictionDigestionSetup(_ctx, input)
	_RestrictionDigestionSteps(_ctx, input, output)
	_RestrictionDigestionAnalysis(_ctx, input, output)
	_RestrictionDigestionValidation(_ctx, input, output)
	return output
}

func RestrictionDigestionRunSteps(_ctx context.Context, input *RestrictionDigestionInput) *RestrictionDigestionSOutput {
	soutput := &RestrictionDigestionSOutput{}
	output := _RestrictionDigestionRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func RestrictionDigestionNew() interface{} {
	return &RestrictionDigestionElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &RestrictionDigestionInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _RestrictionDigestionRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &RestrictionDigestionInput{},
			Out: &RestrictionDigestionOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wunit.Make_units
)

type RestrictionDigestionElement struct {
	inject.CheckedRunner
}

type RestrictionDigestionInput struct {
	BSAoptional      *wtype.LHComponent
	BSAvol           wunit.Volume
	Buffer           *wtype.LHComponent
	BufferVol        wunit.Volume
	DNAName          string
	DNASolution      *wtype.LHComponent
	DNAVol           wunit.Volume
	EnzSolutions     []*wtype.LHComponent
	EnzVolumestoadd  []wunit.Volume
	EnzymeNames      []string
	InPlate          *wtype.LHPlate
	InactivationTemp wunit.Temperature
	InactivationTime wunit.Time
	OutPlate         *wtype.LHPlate
	ReactionTemp     wunit.Temperature
	ReactionTime     wunit.Time
	ReactionVolume   wunit.Volume
	Water            *wtype.LHComponent
}

type RestrictionDigestionOutput struct {
	Reaction *wtype.LHComponent
}

type RestrictionDigestionSOutput struct {
	Data struct {
	}
	Outputs struct {
		Reaction *wtype.LHComponent
	}
}

func init() {
	if err := addComponent(component.Component{Name: "RestrictionDigestion",
		Constructor: RestrictionDigestionNew,
		Desc: component.ComponentDesc{
			Desc: "",
			Path: "src/github.com/antha-lang/antha/antha/component/an/Liquid_handling/RestrictionDigestion/set_volumes/RestrictionDigestion.an",
			Params: []component.ParamDesc{
				{Name: "BSAoptional", Desc: "", Kind: "Inputs"},
				{Name: "BSAvol", Desc: "", Kind: "Parameters"},
				{Name: "Buffer", Desc: "", Kind: "Inputs"},
				{Name: "BufferVol", Desc: "", Kind: "Parameters"},
				{Name: "DNAName", Desc: "", Kind: "Parameters"},
				{Name: "DNASolution", Desc: "", Kind: "Inputs"},
				{Name: "DNAVol", Desc: "", Kind: "Parameters"},
				{Name: "EnzSolutions", Desc: "", Kind: "Inputs"},
				{Name: "EnzVolumestoadd", Desc: "\tStockReConcinUperml \t\t[]int\n\tDesiredConcinUperml\t \t\t[]int\n", Kind: "Parameters"},
				{Name: "EnzymeNames", Desc: "", Kind: "Parameters"},
				{Name: "InPlate", Desc: "", Kind: "Inputs"},
				{Name: "InactivationTemp", Desc: "", Kind: "Parameters"},
				{Name: "InactivationTime", Desc: "", Kind: "Parameters"},
				{Name: "OutPlate", Desc: "", Kind: "Inputs"},
				{Name: "ReactionTemp", Desc: "", Kind: "Parameters"},
				{Name: "ReactionTime", Desc: "", Kind: "Parameters"},
				{Name: "ReactionVolume", Desc: "", Kind: "Parameters"},
				{Name: "Water", Desc: "", Kind: "Inputs"},
				{Name: "Reaction", Desc: "", Kind: "Outputs"},
			},
		},
	}); err != nil {
		panic(err)
	}
}
