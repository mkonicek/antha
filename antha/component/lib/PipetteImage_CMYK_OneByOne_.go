// Generates instructions to pipette out a defined image onto a defined plate by blending cyan magenta yellow and black dyes
package lib

import (
	"fmt"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/image"
	"github.com/antha-lang/antha/antha/anthalib/mixer"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/component"
	"github.com/antha-lang/antha/execute"
	"github.com/antha-lang/antha/inject"
	"golang.org/x/net/context"
)

//"image/color"

// Input parameters for this protocol (data)

// Data which is returned from this protocol, and data types

// Physical Inputs to this protocol with types

//InPlate *wtype.LHPlate

// Physical outputs from this protocol with types

func _PipetteImage_CMYK_OneByOneRequirements() {

}

// Conditions to run on startup
func _PipetteImage_CMYK_OneByOneSetup(_ctx context.Context, _input *PipetteImage_CMYK_OneByOneInput) {

}

// The core process for this protocol, with the steps to be performed
// for every input
func _PipetteImage_CMYK_OneByOneSteps(_ctx context.Context, _input *PipetteImage_CMYK_OneByOneInput, _output *PipetteImage_CMYK_OneByOneOutput) {

	//var chosencolourpalette color.Palette
	chosencolourpalette := image.AvailablePalettes()["Plan9"]
	positiontocolourmap, _, _ := image.ImagetoPlatelayout(_input.Imagefilename, _input.OutPlate, &chosencolourpalette, _input.Rotate, _input.AutoRotate)

	solutions := make([]*wtype.LHComponent, 0)

	counter := 0

	//solutions := image.PipetteImagebyBlending(OutPlate, positiontocolourmap,Cyan, Magenta, Yellow,Black, VolumeForFullcolour)

	for locationkey, colour := range positiontocolourmap {

		components := make([]*wtype.LHComponent, 0)

		cmyk := image.ColourtoCMYK(colour)

		var maxuint8 uint8 = 255

		if cmyk.C == 0 && cmyk.Y == 0 && cmyk.M == 0 && cmyk.K == 0 {

			continue

		} else {

			counter = counter + 1

			if cmyk.C > 0 {

				cyanvol := wunit.NewVolume(((float64(cmyk.C) / float64(maxuint8)) * _input.VolumeForFullcolour.RawValue()), _input.VolumeForFullcolour.Unit().PrefixedSymbol())
				cyanSample := mixer.Sample(_input.Cyan, cyanvol)
				components = append(components, cyanSample)
			}

			if cmyk.Y > 0 {
				yellowvol := wunit.NewVolume(((float64(cmyk.Y) / float64(maxuint8)) * _input.VolumeForFullcolour.RawValue()), _input.VolumeForFullcolour.Unit().PrefixedSymbol())
				yellowSample := mixer.Sample(_input.Yellow, yellowvol)
				components = append(components, yellowSample)
			}

			if cmyk.M > 0 {
				magentavol := wunit.NewVolume(((float64(cmyk.M) / float64(maxuint8)) * _input.VolumeForFullcolour.RawValue()), _input.VolumeForFullcolour.Unit().PrefixedSymbol())
				magentaSample := mixer.Sample(_input.Magenta, magentavol)
				components = append(components, magentaSample)
			}

			if cmyk.K > 0 {
				blackvol := wunit.NewVolume(((float64(cmyk.K) / float64(maxuint8)) * _input.VolumeForFullcolour.RawValue()), _input.VolumeForFullcolour.Unit().PrefixedSymbol())
				blackSample := mixer.Sample(_input.Black, blackvol)
				components = append(components, blackSample)
			}

			solution := execute.MixTo(_ctx, _input.OutPlate.Type, locationkey, 1, components...)
			solutions = append(solutions, solution)

		}
	}

	_output.Pixels = solutions
	_output.Numberofpixels = len(_output.Pixels)
	fmt.Println("Pixels =", _output.Numberofpixels)

}

// Run after controls and a steps block are completed to
// post process any data and provide downstream results
func _PipetteImage_CMYK_OneByOneAnalysis(_ctx context.Context, _input *PipetteImage_CMYK_OneByOneInput, _output *PipetteImage_CMYK_OneByOneOutput) {
}

// A block of tests to perform to validate that the sample was processed correctly
// Optionally, destructive tests can be performed to validate results on a
// dipstick basis
func _PipetteImage_CMYK_OneByOneValidation(_ctx context.Context, _input *PipetteImage_CMYK_OneByOneInput, _output *PipetteImage_CMYK_OneByOneOutput) {

}
func _PipetteImage_CMYK_OneByOneRun(_ctx context.Context, input *PipetteImage_CMYK_OneByOneInput) *PipetteImage_CMYK_OneByOneOutput {
	output := &PipetteImage_CMYK_OneByOneOutput{}
	_PipetteImage_CMYK_OneByOneSetup(_ctx, input)
	_PipetteImage_CMYK_OneByOneSteps(_ctx, input, output)
	_PipetteImage_CMYK_OneByOneAnalysis(_ctx, input, output)
	_PipetteImage_CMYK_OneByOneValidation(_ctx, input, output)
	return output
}

func PipetteImage_CMYK_OneByOneRunSteps(_ctx context.Context, input *PipetteImage_CMYK_OneByOneInput) *PipetteImage_CMYK_OneByOneSOutput {
	soutput := &PipetteImage_CMYK_OneByOneSOutput{}
	output := _PipetteImage_CMYK_OneByOneRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func PipetteImage_CMYK_OneByOneNew() interface{} {
	return &PipetteImage_CMYK_OneByOneElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &PipetteImage_CMYK_OneByOneInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _PipetteImage_CMYK_OneByOneRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &PipetteImage_CMYK_OneByOneInput{},
			Out: &PipetteImage_CMYK_OneByOneOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wunit.Make_units
)

type PipetteImage_CMYK_OneByOneElement struct {
	inject.CheckedRunner
}

type PipetteImage_CMYK_OneByOneInput struct {
	AutoRotate          bool
	Black               *wtype.LHComponent
	Cyan                *wtype.LHComponent
	Imagefilename       string
	Magenta             *wtype.LHComponent
	OutPlate            *wtype.LHPlate
	Rotate              bool
	VolumeForFullcolour wunit.Volume
	Yellow              *wtype.LHComponent
}

type PipetteImage_CMYK_OneByOneOutput struct {
	Numberofpixels int
	Pixels         []*wtype.LHComponent
}

type PipetteImage_CMYK_OneByOneSOutput struct {
	Data struct {
		Numberofpixels int
	}
	Outputs struct {
		Pixels []*wtype.LHComponent
	}
}

func init() {
	if err := addComponent(component.Component{Name: "PipetteImage_CMYK_OneByOne",
		Constructor: PipetteImage_CMYK_OneByOneNew,
		Desc: component.ComponentDesc{
			Desc: "Generates instructions to pipette out a defined image onto a defined plate by blending cyan magenta yellow and black dyes\n",
			Path: "src/github.com/antha-lang/antha/antha/component/an/Liquid_handling/PipetteImage/PipetteImage_CMYK_OneByOne.an",
			Params: []component.ParamDesc{
				{Name: "AutoRotate", Desc: "", Kind: "Parameters"},
				{Name: "Black", Desc: "", Kind: "Inputs"},
				{Name: "Cyan", Desc: "", Kind: "Inputs"},
				{Name: "Imagefilename", Desc: "", Kind: "Parameters"},
				{Name: "Magenta", Desc: "", Kind: "Inputs"},
				{Name: "OutPlate", Desc: "InPlate *wtype.LHPlate\n", Kind: "Inputs"},
				{Name: "Rotate", Desc: "", Kind: "Parameters"},
				{Name: "VolumeForFullcolour", Desc: "", Kind: "Parameters"},
				{Name: "Yellow", Desc: "", Kind: "Inputs"},
				{Name: "Numberofpixels", Desc: "", Kind: "Data"},
				{Name: "Pixels", Desc: "", Kind: "Outputs"},
			},
		},
	}); err != nil {
		panic(err)
	}
}
