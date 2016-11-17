package lib

import (
	"fmt"
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

// Physical Inputs to this protocol with types

// Physical outputs from this protocol with types

// Data which is returned from this protocol, and data types

func _RestrictionDigestion_concRequirements() {}

// Conditions to run on startup
func _RestrictionDigestion_concSetup(_ctx context.Context, _input *RestrictionDigestion_concInput) {}

// The core process for this protocol, with the steps to be performed
// for every input
func _RestrictionDigestion_concSteps(_ctx context.Context, _input *RestrictionDigestion_concInput, _output *RestrictionDigestion_concOutput) {

	statii := make([]string, 0)

	samples := make([]*wtype.LHComponent, 0)
	waterSample := mixer.SampleForTotalVolume(_input.Water, _input.ReactionVolume)
	samples = append(samples, waterSample)

	// workout volume of buffer to add in SI units
	BufferVol := wunit.NewVolume(float64(_input.ReactionVolume.SIValue()/float64(_input.BufferConcX)), "l")
	statii = append(statii, fmt.Sprintln("buffer volume conversion:", _input.ReactionVolume.SIValue(), _input.BufferConcX, float64(_input.ReactionVolume.SIValue()/float64(_input.BufferConcX)), " Buffervol = ", BufferVol.SIValue()))
	bufferSample := mixer.Sample(_input.Buffer, BufferVol)
	samples = append(samples, bufferSample)

	if _input.BSAvol.Mvalue != 0 {
		bsaSample := mixer.Sample(_input.BSAoptional, _input.BSAvol)
		samples = append(samples, bsaSample)
	}

	_input.DNASolution.CName = _input.DNAName

	// work out necessary volume to add
	DNAVol, err := wunit.VolumeForTargetMass(_input.DNAMassperReaction, _input.DNAConc) //NewVolume(float64((DNAMassperReaction.SIValue()/DNAConc.SIValue())),"l")

	if err != nil {
		execute.Errorf(_ctx, err.Error())
	}

	statii = append(statii, fmt.Sprintln("DNA MAss to Volume conversion:", _input.DNAMassperReaction.SIValue(), _input.DNAConc.SIValue(), float64((_input.DNAMassperReaction.SIValue()/_input.DNAConc.SIValue())), "DNAVol =", DNAVol.SIValue()))
	statii = append(statii, fmt.Sprintln("DNAVOL", DNAVol.ToString()))
	dnaSample := mixer.Sample(_input.DNASolution, DNAVol)
	samples = append(samples, dnaSample)

	for k, enzyme := range _input.EnzSolutions {

		/*
			e.g.
			DesiredUinreaction = 1  // U
			StockReConcinUperml = 10000 // U/ml
			ReactionVolume = 20ul
		*/
		stockconcinUperul := _input.StockReConcinUperml[k] / 1000
		enzvoltoaddinul := _input.DesiredConcinUperml[k] / stockconcinUperul

		var enzvoltoadd wunit.Volume

		if float64(enzvoltoaddinul) < 0.5 {
			enzvoltoadd = wunit.NewVolume(float64(0.5), "ul")
		} else {
			enzvoltoadd = wunit.NewVolume(float64(enzvoltoaddinul), "ul")
		}
		enzyme.CName = _input.EnzymeNames[k]
		text.Print("adding enzyme"+_input.EnzymeNames[k], "to"+_input.DNAName)
		enzSample := mixer.Sample(enzyme, enzvoltoadd)
		enzSample.CName = _input.EnzymeNames[k]
		samples = append(samples, enzSample)
	}

	// incubate the reaction mixture
	r1 := execute.Incubate(_ctx, execute.MixTo(_ctx, _input.OutPlate.Type, "", _input.Platenumber, samples...), _input.ReactionTemp, _input.ReactionTime, false)
	// inactivate
	_output.Reaction = execute.Incubate(_ctx, r1, _input.InactivationTemp, _input.InactivationTime, false)
}

// Run after controls and a steps block are completed to
// post process any data and provide downstream results
func _RestrictionDigestion_concAnalysis(_ctx context.Context, _input *RestrictionDigestion_concInput, _output *RestrictionDigestion_concOutput) {
}

// A block of tests to perform to validate that the sample was processed correctly
// Optionally, destructive tests can be performed to validate results on a
// dipstick basis
func _RestrictionDigestion_concValidation(_ctx context.Context, _input *RestrictionDigestion_concInput, _output *RestrictionDigestion_concOutput) {
}
func _RestrictionDigestion_concRun(_ctx context.Context, input *RestrictionDigestion_concInput) *RestrictionDigestion_concOutput {
	output := &RestrictionDigestion_concOutput{}
	_RestrictionDigestion_concSetup(_ctx, input)
	_RestrictionDigestion_concSteps(_ctx, input, output)
	_RestrictionDigestion_concAnalysis(_ctx, input, output)
	_RestrictionDigestion_concValidation(_ctx, input, output)
	return output
}

func RestrictionDigestion_concRunSteps(_ctx context.Context, input *RestrictionDigestion_concInput) *RestrictionDigestion_concSOutput {
	soutput := &RestrictionDigestion_concSOutput{}
	output := _RestrictionDigestion_concRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func RestrictionDigestion_concNew() interface{} {
	return &RestrictionDigestion_concElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &RestrictionDigestion_concInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _RestrictionDigestion_concRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &RestrictionDigestion_concInput{},
			Out: &RestrictionDigestion_concOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wunit.Make_units
)

type RestrictionDigestion_concElement struct {
	inject.CheckedRunner
}

type RestrictionDigestion_concInput struct {
	BSAoptional         *wtype.LHComponent
	BSAvol              wunit.Volume
	Buffer              *wtype.LHComponent
	BufferConcX         int
	DNAConc             wunit.Concentration
	DNAMassperReaction  wunit.Mass
	DNAName             string
	DNASolution         *wtype.LHComponent
	DesiredConcinUperml []int
	EnzSolutions        []*wtype.LHComponent
	EnzymeNames         []string
	InactivationTemp    wunit.Temperature
	InactivationTime    wunit.Time
	OutPlate            *wtype.LHPlate
	Platenumber         int
	ReactionTemp        wunit.Temperature
	ReactionTime        wunit.Time
	ReactionVolume      wunit.Volume
	StockReConcinUperml []int
	Water               *wtype.LHComponent
}

type RestrictionDigestion_concOutput struct {
	Reaction *wtype.LHComponent
	Status   string
}

type RestrictionDigestion_concSOutput struct {
	Data struct {
		Status string
	}
	Outputs struct {
		Reaction *wtype.LHComponent
	}
}

func init() {
	if err := addComponent(component.Component{Name: "RestrictionDigestion_conc",
		Constructor: RestrictionDigestion_concNew,
		Desc: component.ComponentDesc{
			Desc: "",
			Path: "src/github.com/antha-lang/antha/antha/component/an/Liquid_handling/RestrictionDigestion/set_conc/RestrictionDigestion.an",
			Params: []component.ParamDesc{
				{Name: "BSAoptional", Desc: "", Kind: "Inputs"},
				{Name: "BSAvol", Desc: "", Kind: "Parameters"},
				{Name: "Buffer", Desc: "", Kind: "Inputs"},
				{Name: "BufferConcX", Desc: "", Kind: "Parameters"},
				{Name: "DNAConc", Desc: "", Kind: "Parameters"},
				{Name: "DNAMassperReaction", Desc: "", Kind: "Parameters"},
				{Name: "DNAName", Desc: "", Kind: "Parameters"},
				{Name: "DNASolution", Desc: "", Kind: "Inputs"},
				{Name: "DesiredConcinUperml", Desc: "", Kind: "Parameters"},
				{Name: "EnzSolutions", Desc: "", Kind: "Inputs"},
				{Name: "EnzymeNames", Desc: "", Kind: "Parameters"},
				{Name: "InactivationTemp", Desc: "", Kind: "Parameters"},
				{Name: "InactivationTime", Desc: "", Kind: "Parameters"},
				{Name: "OutPlate", Desc: "", Kind: "Inputs"},
				{Name: "Platenumber", Desc: "", Kind: "Parameters"},
				{Name: "ReactionTemp", Desc: "", Kind: "Parameters"},
				{Name: "ReactionTime", Desc: "", Kind: "Parameters"},
				{Name: "ReactionVolume", Desc: "", Kind: "Parameters"},
				{Name: "StockReConcinUperml", Desc: "", Kind: "Parameters"},
				{Name: "Water", Desc: "", Kind: "Inputs"},
				{Name: "Reaction", Desc: "", Kind: "Outputs"},
				{Name: "Status", Desc: "", Kind: "Data"},
			},
		},
	}); err != nil {
		panic(err)
	}
}
