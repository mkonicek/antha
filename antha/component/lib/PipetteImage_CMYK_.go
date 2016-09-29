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

func _PipetteImage_CMYKRequirements() {

}

// Conditions to run on startup
func _PipetteImage_CMYKSetup(_ctx context.Context, _input *PipetteImage_CMYKInput) {

}

// The core process for this protocol, with the steps to be performed
// for every input
func _PipetteImage_CMYKSteps(_ctx context.Context, _input *PipetteImage_CMYKInput, _output *PipetteImage_CMYKOutput) {

	//var chosencolourpalette color.Palette
	chosencolourpalette := image.AvailablePalettes()["Plan9"]
	positiontocolourmap, _, _ := image.ImagetoPlatelayout(_input.Imagefilename, _input.OutPlate, &chosencolourpalette, _input.Rotate, _input.AutoRotate)

	solutions := make([]*wtype.LHComponent, 0)

	counter := 0

	//solutions := image.PipetteImagebyBlending(OutPlate, positiontocolourmap,Cyan, Magenta, Yellow,Black, VolumeForFullcolour)

	for locationkey, colour := range positiontocolourmap {

		//components := make([]*wtype.LHComponent, 0)
		var solution *wtype.LHComponent

		cmyk := image.ColourtoCMYK(colour)

		var maxuint8 uint8 = 255

		if cmyk.C == 0 && cmyk.Y == 0 && cmyk.M == 0 && cmyk.K == 0 {

			continue

		} else {

			counter = counter + 1

			if cmyk.C > 0 {

				cyanvol := wunit.NewVolume(((float64(cmyk.C) / float64(maxuint8)) * _input.VolumeForFullcolour.RawValue()), _input.VolumeForFullcolour.Unit().PrefixedSymbol())

				if cyanvol.RawValue() < 10 && cyanvol.Unit().PrefixedSymbol() == "ul" {
					cyanvol.SetValue(10)
				}

				if cmyk.K == 0 && cmyk.M == 0 && cmyk.Y == 0 {
					_input.Cyan.Type = wtype.LTNeedToMix
				} else {
					_input.Cyan.Type = wtype.LTDISPENSEABOVE
				}

				cyanSample := mixer.Sample(_input.Cyan, cyanvol)

				solution = execute.MixTo(_ctx, _input.OutPlate.Type, locationkey, 1, cyanSample)
			}
			if cmyk.Y > 0 {
				yellowvol := wunit.NewVolume(((float64(cmyk.Y) / float64(maxuint8)) * _input.VolumeForFullcolour.RawValue()), _input.VolumeForFullcolour.Unit().PrefixedSymbol())

				if yellowvol.RawValue() < 10 && yellowvol.Unit().PrefixedSymbol() == "ul" {
					yellowvol.SetValue(10)
				}
				if cmyk.K == 0 && cmyk.M == 0 {
					_input.Yellow.Type = wtype.LTNeedToMix
				} else {
					_input.Yellow.Type = wtype.LTDISPENSEABOVE
				}

				yellowSample := mixer.Sample(_input.Yellow, yellowvol)

				if solution != nil {
					solution = execute.Mix(_ctx, solution, yellowSample)
				} else {
					//solution = MixInto(PalettePlate, "", yellowSample)
					solution = execute.MixTo(_ctx, _input.OutPlate.Type, locationkey, 1, yellowSample)
				}
			}
			if cmyk.M > 0 {
				magentavol := wunit.NewVolume(((float64(cmyk.M) / float64(maxuint8)) * _input.VolumeForFullcolour.RawValue()), _input.VolumeForFullcolour.Unit().PrefixedSymbol())

				if magentavol.RawValue() < 10 && magentavol.Unit().PrefixedSymbol() == "ul" {
					magentavol.SetValue(10)
				}

				if cmyk.K == 0 {
					_input.Magenta.Type = wtype.LTNeedToMix
				} else {
					_input.Magenta.Type = wtype.LTDISPENSEABOVE
				}

				magentaSample := mixer.Sample(_input.Magenta, magentavol)

				if solution != nil {
					solution = execute.Mix(_ctx, solution, magentaSample)
				} else {
					//solution = MixInto(PalettePlate, "", magentaSample)
					solution = execute.MixTo(_ctx, _input.OutPlate.Type, locationkey, 1, magentaSample)
				}
			}
			if cmyk.K > 0 {
				blackvol := wunit.NewVolume(((float64(cmyk.K) / float64(maxuint8)) * _input.VolumeForFullcolour.RawValue()), _input.VolumeForFullcolour.Unit().PrefixedSymbol())

				if blackvol.RawValue() < 10 && blackvol.Unit().PrefixedSymbol() == "ul" {
					blackvol.SetValue(10)
				}

				_input.Black.Type = wtype.LTNeedToMix

				blackSample := mixer.Sample(_input.Black, blackvol)

				if solution != nil {
					solution = execute.Mix(_ctx, solution, blackSample)
				} else {
					//solution = MixInto(PalettePlate, "", blackSample)
					solution = execute.MixTo(_ctx, _input.OutPlate.Type, locationkey, 1, blackSample)
				}

				//components = append(components, blackSample)
			}

			//solution := MixTo(OutPlate.Type, locationkey,1, components...)
			solutions = append(solutions, solution)

		}
	}

	_output.Pixels = solutions
	_output.Numberofpixels = len(_output.Pixels)
	fmt.Println("Pixels =", _output.Numberofpixels)

}

// Run after controls and a steps block are completed to
// post process any data and provide downstream results
func _PipetteImage_CMYKAnalysis(_ctx context.Context, _input *PipetteImage_CMYKInput, _output *PipetteImage_CMYKOutput) {
}

// A block of tests to perform to validate that the sample was processed correctly
// Optionally, destructive tests can be performed to validate results on a
// dipstick basis
func _PipetteImage_CMYKValidation(_ctx context.Context, _input *PipetteImage_CMYKInput, _output *PipetteImage_CMYKOutput) {

}
func _PipetteImage_CMYKRun(_ctx context.Context, input *PipetteImage_CMYKInput) *PipetteImage_CMYKOutput {
	output := &PipetteImage_CMYKOutput{}
	_PipetteImage_CMYKSetup(_ctx, input)
	_PipetteImage_CMYKSteps(_ctx, input, output)
	_PipetteImage_CMYKAnalysis(_ctx, input, output)
	_PipetteImage_CMYKValidation(_ctx, input, output)
	return output
}

func PipetteImage_CMYKRunSteps(_ctx context.Context, input *PipetteImage_CMYKInput) *PipetteImage_CMYKSOutput {
	soutput := &PipetteImage_CMYKSOutput{}
	output := _PipetteImage_CMYKRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func PipetteImage_CMYKNew() interface{} {
	return &PipetteImage_CMYKElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &PipetteImage_CMYKInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _PipetteImage_CMYKRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &PipetteImage_CMYKInput{},
			Out: &PipetteImage_CMYKOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wunit.Make_units
)

type PipetteImage_CMYKElement struct {
	inject.CheckedRunner
}

type PipetteImage_CMYKInput struct {
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

type PipetteImage_CMYKOutput struct {
	Numberofpixels int
	Pixels         []*wtype.LHComponent
}

type PipetteImage_CMYKSOutput struct {
	Data struct {
		Numberofpixels int
	}
	Outputs struct {
		Pixels []*wtype.LHComponent
	}
}

func init() {
	if err := addComponent(component.Component{Name: "PipetteImage_CMYK",
		Constructor: PipetteImage_CMYKNew,
		Desc: component.ComponentDesc{
			Desc: "Generates instructions to pipette out a defined image onto a defined plate by blending cyan magenta yellow and black dyes\n",
			Path: "antha/component/an/Liquid_handling/PipetteImage/PipetteImage_CMYK.an",
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
