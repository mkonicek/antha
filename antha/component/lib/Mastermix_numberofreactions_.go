package lib

import (
	"github.com/antha-lang/antha/antha/anthalib/mixer"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/component"
	"github.com/antha-lang/antha/execute"
	"github.com/antha-lang/antha/inject"
	"golang.org/x/net/context"
)

// Input parameters for this protocol (data)

//TotalVolumeperreaction Volume // if buffer is being added
//VolumetoLeaveforDNAperreaction Volume

//NumberofMastermixes int // add as many as possible option e.g. if == -1

// Data which is returned from this protocol, and data types

// Physical Inputs to this protocol with types

//TopUpBuffer *wtype.LHComponent // optional if nil this is ignored

// Physical outputs from this protocol with types

func _Mastermix_numberofreactionsRequirements() {
}

// Conditions to run on startup
func _Mastermix_numberofreactionsSetup(_ctx context.Context, _input *Mastermix_numberofreactionsInput) {
}

// The core process for this protocol, with the steps to be performed
// for every input
func _Mastermix_numberofreactionsSteps(_ctx context.Context, _input *Mastermix_numberofreactionsInput, _output *Mastermix_numberofreactionsOutput) {

	wellpositions := _input.OutPlate.AllWellPositions(wtype.BYCOLUMN)

	var counter int

	var mastermix *wtype.LHComponent

	// work out volume to top up to in each case (per reaction) in l:
	//topupVolumeperreacttion := TotalVolumeperreaction.SIValue() - VolumetoLeaveforDNAperreaction.SIValue()

	// multiply by number of reactions per mastermix
	//topupVolume := wunit.NewVolume(float64(Reactionspermastermix)*topupVolumeperreacttion,"l")

	if len(_input.Components) != len(_input.ComponentVolumesperReaction) {
		panic("len(Components) != len(OtherComponentVolumes)")
	}

	eachmastermix := make([]*wtype.LHComponent, 0)

	//if TopUpBuffer != nil {
	//bufferSample := mixer.SampleForTotalVolume(TopUpBuffer, topupVolume)
	//eachmastermix = append(eachmastermix,bufferSample)
	//	}

	for k, component := range _input.Components {
		if k == len(_input.Components) {
			component.Type = wtype.LTNeedToMix //"NeedToMix"
		}

		// multiply volume of each component by number of reactions per mastermix
		adjustedvol := wunit.NewVolume(float64(_input.NumberofReactions)*_input.ComponentVolumesperReaction[k].SIValue()*1000000, "ul")

		componentSample := mixer.Sample(component, adjustedvol)
		eachmastermix = append(eachmastermix, componentSample)

	}

	/*
		totalvolumeoverwellcapacity := wunit.wunit.AddVolumes(eachmastermix).RawValue()/OutPlate.WellMap()[counter].MaxVolume().RawValue()

		numberoffullvolumes := wutil.RoundDown(totalvolumeoverwellcapacity)

		remainder := math.Remainder(wunit.AddVolumes(eachmastermix).RawValue(),OutPlate.WellMap()[counter].MaxVolume().RawValue())

		for i := counter; i < numberoffullvolumes;i++{

			eachmastermix := make([]*wtype.LHComponent, 0)

			//if TopUpBuffer != nil {
			//bufferSample := mixer.SampleForTotalVolume(TopUpBuffer, topupVolume)
			//eachmastermix = append(eachmastermix,bufferSample)
		//	}

			for k,component := range Components {
				if k == len(Components){
					component.Type = wtype.LTNeedToMix //"NeedToMix"
				}

			// multiply volume of each component by number of reactions per mastermix
			adjustedvol := wunit.NewVolume(float64(NumberofReactions)*ComponentVolumesperReaction[k].SIValue()*1000000,"ul")

			componentSample := mixer.Sample(component,adjustedvol)
			eachmastermix = append(eachmastermix,componentSample)


			}

	*/
	/*
		if  wunit.AddVolumes(eachmastermix).RawValue() > (OutPlate.WellMap()[counter].MaxVolume().RawValue()-OutPlate.WellMap()[counter].ResidualVolume().RawValue()){
			Errorf("Volume too high for desitination well, use bigger destination well or split")
		}
	*/

	mastermix = execute.MixInto(_ctx, _input.OutPlate, wellpositions[counter], eachmastermix...)

	_output.Mastermix = mastermix

	_output.Error = wtype.ExportPlateCSV(_input.Projectname+"mastermix.csv", _input.OutPlate, _input.Projectname+"outputPlate", []string{wellpositions[counter]}, []*wtype.LHComponent{_output.Mastermix}, []wunit.Volume{_output.Mastermix.Volume()})

}

// Run after controls and a steps block are completed to
// post process any data and provide downstream results
func _Mastermix_numberofreactionsAnalysis(_ctx context.Context, _input *Mastermix_numberofreactionsInput, _output *Mastermix_numberofreactionsOutput) {
}

// A block of tests to perform to validate that the sample was processed correctly
// Optionally, destructive tests can be performed to validate results on a
// dipstick basis
func _Mastermix_numberofreactionsValidation(_ctx context.Context, _input *Mastermix_numberofreactionsInput, _output *Mastermix_numberofreactionsOutput) {
}
func _Mastermix_numberofreactionsRun(_ctx context.Context, input *Mastermix_numberofreactionsInput) *Mastermix_numberofreactionsOutput {
	output := &Mastermix_numberofreactionsOutput{}
	_Mastermix_numberofreactionsSetup(_ctx, input)
	_Mastermix_numberofreactionsSteps(_ctx, input, output)
	_Mastermix_numberofreactionsAnalysis(_ctx, input, output)
	_Mastermix_numberofreactionsValidation(_ctx, input, output)
	return output
}

func Mastermix_numberofreactionsRunSteps(_ctx context.Context, input *Mastermix_numberofreactionsInput) *Mastermix_numberofreactionsSOutput {
	soutput := &Mastermix_numberofreactionsSOutput{}
	output := _Mastermix_numberofreactionsRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func Mastermix_numberofreactionsNew() interface{} {
	return &Mastermix_numberofreactionsElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &Mastermix_numberofreactionsInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _Mastermix_numberofreactionsRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &Mastermix_numberofreactionsInput{},
			Out: &Mastermix_numberofreactionsOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wunit.Make_units
)

type Mastermix_numberofreactionsElement struct {
	inject.CheckedRunner
}

type Mastermix_numberofreactionsInput struct {
	ComponentVolumesperReaction []wunit.Volume
	Components                  []*wtype.LHComponent
	NumberofReactions           int
	OutPlate                    *wtype.LHPlate
	Projectname                 string
}

type Mastermix_numberofreactionsOutput struct {
	Error     error
	Mastermix *wtype.LHComponent
	Status    string
}

type Mastermix_numberofreactionsSOutput struct {
	Data struct {
		Error  error
		Status string
	}
	Outputs struct {
		Mastermix *wtype.LHComponent
	}
}

func init() {
	if err := addComponent(component.Component{Name: "Mastermix_numberofreactions",
		Constructor: Mastermix_numberofreactionsNew,
		Desc: component.ComponentDesc{
			Desc: "",
			Path: "antha/component/an/Liquid_handling/MakeMastermix/Mastermix_number.an",
			Params: []component.ParamDesc{
				{Name: "ComponentVolumesperReaction", Desc: "", Kind: "Parameters"},
				{Name: "Components", Desc: "TopUpBuffer *wtype.LHComponent // optional if nil this is ignored\n", Kind: "Inputs"},
				{Name: "NumberofReactions", Desc: "TotalVolumeperreaction Volume // if buffer is being added\nVolumetoLeaveforDNAperreaction Volume\n", Kind: "Parameters"},
				{Name: "OutPlate", Desc: "", Kind: "Inputs"},
				{Name: "Projectname", Desc: "", Kind: "Parameters"},
				{Name: "Error", Desc: "", Kind: "Data"},
				{Name: "Mastermix", Desc: "", Kind: "Outputs"},
				{Name: "Status", Desc: "", Kind: "Data"},
			},
		},
	}); err != nil {
		panic(err)
	}
}
