// Example protocol demonstrating the use of the Sample function
package lib

import

// we need to import the wtype package to use the LHComponent type
// the mixer package is required to use the Sample function
(
	"fmt"
	"github.com/antha-lang/antha/antha/anthalib/mixer"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/component"
	"github.com/antha-lang/antha/execute"
	"github.com/antha-lang/antha/inject"
	"golang.org/x/net/context"
)

// Input parameters for this protocol (data)

//Parts []wtype.DNASequences

//DNAparts []*wtype.LHComponent

func _RatioRequirements() {
}

func _RatioSetup(_ctx context.Context, _input *RatioInput) {
}

func _RatioSteps(_ctx context.Context, _input *RatioInput, _output *RatioOutput) {

	if len(_input.Part) != len(_input.Tot) || len(_input.Part) != len(_input.PlasmidConc) || len(_input.Part) != len(_input.Partnames) {
		execute.Errorf(_ctx, fmt.Sprint("What the hell do you think you're doing these lists aren't equal in size. len(Parts):", len(_input.Part), " len(Partnames):", len(_input.Partnames)))
	}

	//var inputname string = Inputtype.CName

	Ratio := make([]float64, len(_input.Part))
	_output.SampleVolumesUsed = make([]wunit.Volume, 0)

	samples := make([]*wtype.LHComponent, 0)

	//Calculate the ratio of part to vtotal vector size and store in a slice, also calculate the sum of all the ratios.
	RatioTotal := 0.0
	Totalng := 0.0
	TotalVol := 0.0

	for i := 0; i < len(_input.Part); i++ {
		Ratio[i] = float64(_input.Part[i]) / float64(_input.Tot[i])
		if Ratio[i] == 0 {
			execute.Errorf(_ctx, fmt.Sprint("Ratio of zero for ", _input.Part[i], " / ", _input.Tot[i]))
		}
		RatioTotal += Ratio[i]
	}

	//work out each ratio contribution in terms of ng of plasmid for the required ammount of plasmid (100ng) in the mix (e.g. ratio*100/ratiototal) store in a slice (NgVector)
	NgVector := make([]float64, len(_input.Part))
	for i := 0; i < len(_input.Part); i++ {
		NgVector[i] = Ratio[i] * _input.NgRequired / RatioTotal
		Totalng += NgVector[i]
	}

	//First add water for DNA to be pipetted into
	watersample := mixer.SampleForTotalVolume(_input.Buffer, _input.TotalVolume)
	samples = append(samples, watersample)

	//From the slice of plasmid concentrations and ng of plasmid required (NgVector) calculate the volume of each plasmid required and store in a
	//slice. If a volume is less that 0.5 ul then multiply all volumes by 200 (arbitrary number for now)

	VolPlasmid := make([]float64, len(_input.Part))
	for i := 0; i < len(_input.Part); i++ {
		VolPlasmid[i] = NgVector[i] / _input.PlasmidConc[i]
		if VolPlasmid[i] < 0.5 {
			VolPlasmid[i] = VolPlasmid[i] * 200
			TotalVol += VolPlasmid[i]

			volused := wunit.NewVolume(VolPlasmid[i], "ul")
			_output.SampleVolumesUsed = append(_output.SampleVolumesUsed, volused)

			_input.Inputtype.CName = _input.Partnames[i] // inputname + fmt.Sprint(i+1)
			partSample := mixer.Sample(_input.Inputtype, volused)

			samples = append(samples, partSample)

			fmt.Println(VolPlasmid[i])
		}

	}

	_output.PooledLib = execute.MixTo(_ctx, _input.PlateType.Type, "", 1, samples...)
	//	DNAparts = samples

}

func _RatioAnalysis(_ctx context.Context, _input *RatioInput, _output *RatioOutput) {
}

func _RatioValidation(_ctx context.Context, _input *RatioInput, _output *RatioOutput) {
}
func _RatioRun(_ctx context.Context, input *RatioInput) *RatioOutput {
	output := &RatioOutput{}
	_RatioSetup(_ctx, input)
	_RatioSteps(_ctx, input, output)
	_RatioAnalysis(_ctx, input, output)
	_RatioValidation(_ctx, input, output)
	return output
}

func RatioRunSteps(_ctx context.Context, input *RatioInput) *RatioSOutput {
	soutput := &RatioSOutput{}
	output := _RatioRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func RatioNew() interface{} {
	return &RatioElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &RatioInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _RatioRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &RatioInput{},
			Out: &RatioOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wunit.Make_units
)

type RatioElement struct {
	inject.CheckedRunner
}

type RatioInput struct {
	Buffer      *wtype.LHComponent
	Inputtype   *wtype.LHComponent
	NgRequired  float64
	Part        []int
	Partnames   []string
	PlasmidConc []float64
	PlateType   *wtype.LHPlate
	Tot         []int
	TotalVolume wunit.Volume
}

type RatioOutput struct {
	PooledLib         *wtype.LHComponent
	SampleVolumesUsed []wunit.Volume
}

type RatioSOutput struct {
	Data struct {
		SampleVolumesUsed []wunit.Volume
	}
	Outputs struct {
		PooledLib *wtype.LHComponent
	}
}

func init() {
	if err := addComponent(component.Component{Name: "Ratio",
		Constructor: RatioNew,
		Desc: component.ComponentDesc{
			Desc: "Example protocol demonstrating the use of the Sample function\n",
			Path: "antha/component/an/playground/ratio.an",
			Params: []component.ParamDesc{
				{Name: "Buffer", Desc: "", Kind: "Inputs"},
				{Name: "Inputtype", Desc: "", Kind: "Inputs"},
				{Name: "NgRequired", Desc: "", Kind: "Parameters"},
				{Name: "Part", Desc: "Parts []wtype.DNASequences\n", Kind: "Parameters"},
				{Name: "Partnames", Desc: "", Kind: "Parameters"},
				{Name: "PlasmidConc", Desc: "", Kind: "Parameters"},
				{Name: "PlateType", Desc: "", Kind: "Inputs"},
				{Name: "Tot", Desc: "", Kind: "Parameters"},
				{Name: "TotalVolume", Desc: "", Kind: "Parameters"},
				{Name: "PooledLib", Desc: "", Kind: "Outputs"},
				{Name: "SampleVolumesUsed", Desc: "", Kind: "Data"},
			},
		},
	}); err != nil {
		panic(err)
	}
}
