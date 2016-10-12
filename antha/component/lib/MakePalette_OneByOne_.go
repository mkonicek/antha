// Generates instructions to make a pallette of all colours in an image
package lib

import (
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/image"
	"github.com/antha-lang/antha/antha/anthalib/mixer"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	//"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/search"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/component"
	"github.com/antha-lang/antha/execute"
	"github.com/antha-lang/antha/inject"
	"golang.org/x/net/context"
	"image/color"
	"strconv"
)

// Input parameters for this protocol (data)

// Data which is returned from this protocol, and data types

//Colournames []string

// Physical Inputs to this protocol with types

//InPlate *wtype.LHPlate

// Physical outputs from this protocol with types

func _MakePalette_OneByOneRequirements() {

}

// Conditions to run on startup
func _MakePalette_OneByOneSetup(_ctx context.Context, _input *MakePalette_OneByOneInput) {

}

// The core process for this protocol, with the steps to be performed
// for every input
func _MakePalette_OneByOneSteps(_ctx context.Context, _input *MakePalette_OneByOneInput, _output *MakePalette_OneByOneOutput) {

	//var chosencolourpalette color.Palette

	//chosencolourpalette := image.AvailablePalettes["Plan9"]

	//positiontocolourmap, _ := image.ImagetoPlatelayout(Imagefilename, OutPlate, &chosencolourpalette, Rotate)

	if _input.PosterizeImage {
		_, _input.Imagefilename = image.Posterize(_input.Imagefilename, _input.PosterizeLevels)
	}

	// make palette of colours from image
	chosencolourpalette := image.MakeSmallPalleteFromImage(_input.Imagefilename, _input.OutPlate, _input.Rotate)

	// make a map of colour to well coordinates
	positiontocolourmap, _, _ := image.ImagetoPlatelayout(_input.Imagefilename, _input.OutPlate, &chosencolourpalette, _input.Rotate, _input.AutoRotate)

	// remove duplicates
	positiontocolourmap = image.RemoveDuplicatesValuesfromMap(positiontocolourmap)

	//fmt.Println("positions", positiontocolourmap)

	solutions := make([]*wtype.LHComponent, 0)
	colourtoComponentMap := make(map[string]*wtype.LHComponent)

	counter := 0

	for _, colour := range positiontocolourmap {

		colourindex := chosencolourpalette.Index(colour)

		if colour != nil {
			components := make([]*wtype.LHComponent, 0)

			cmyk := image.ColourtoCMYK(colour)

			var maxuint8 uint8 = 255

			if cmyk.C == 0 && cmyk.Y == 0 && cmyk.M == 0 && cmyk.K == 0 {

				continue

			} else {

				counter = counter + 1

				if cmyk.C > 0 {

					cyanvol := wunit.NewVolume(((float64(cmyk.C) / float64(maxuint8)) * _input.VolumeForFullcolour.RawValue()), _input.VolumeForFullcolour.Unit().PrefixedSymbol())

					if cyanvol.RawValue() < 0.5 && cyanvol.Unit().PrefixedSymbol() == "ul" {
						cyanvol.SetValue(0.5)
					}

					//if cmyk.K == 0 && cmyk.M == 0 && cmyk.Y == 0{
					_input.Cyan.Type = wtype.LTDISPENSEABOVE
					//}else {
					//Cyan.Type = wtype.LTPostMix
					//}

					cyanSample := mixer.Sample(_input.Cyan, cyanvol)
					components = append(components, cyanSample)
				}

				if cmyk.Y > 0 {
					yellowvol := wunit.NewVolume(((float64(cmyk.Y) / float64(maxuint8)) * _input.VolumeForFullcolour.RawValue()), _input.VolumeForFullcolour.Unit().PrefixedSymbol())

					if yellowvol.RawValue() < 0.5 && yellowvol.Unit().PrefixedSymbol() == "ul" {
						yellowvol.SetValue(0.5)
					}

					if cmyk.K == 0 && cmyk.M == 0 {
						_input.Yellow.Type = wtype.LTPostMix
					} else {
						_input.Yellow.Type = wtype.LTDISPENSEABOVE
					}

					yellowSample := mixer.Sample(_input.Yellow, yellowvol)
					components = append(components, yellowSample)
				}

				if cmyk.M > 0 {
					magentavol := wunit.NewVolume(((float64(cmyk.M) / float64(maxuint8)) * _input.VolumeForFullcolour.RawValue()), _input.VolumeForFullcolour.Unit().PrefixedSymbol())

					if magentavol.RawValue() < 0.5 && magentavol.Unit().PrefixedSymbol() == "ul" {
						magentavol.SetValue(0.5)
					}

					if cmyk.K == 0 {
						_input.Magenta.Type = wtype.LTPostMix
					} else {
						_input.Magenta.Type = wtype.LTDISPENSEABOVE
					}

					magentaSample := mixer.Sample(_input.Magenta, magentavol)
					components = append(components, magentaSample)
				}

				if cmyk.K > 0 {
					blackvol := wunit.NewVolume(((float64(cmyk.K) / float64(maxuint8)) * _input.VolumeForFullcolour.RawValue()), _input.VolumeForFullcolour.Unit().PrefixedSymbol())

					if blackvol.RawValue() < 0.5 && blackvol.Unit().PrefixedSymbol() == "ul" {
						blackvol.SetValue(0.5)
					}

					_input.Black.Type = wtype.LTPostMix

					blackSample := mixer.Sample(_input.Black, blackvol)
					components = append(components, blackSample)
				}

				solution := execute.MixInto(_ctx, _input.PalettePlate, "", components...)
				solutions = append(solutions, solution)
				colourtoComponentMap[strconv.Itoa(colourindex)] = solution

			}

		}
	}

	_output.Colours = solutions
	_output.Numberofcolours = len(chosencolourpalette)
	_output.Palette = chosencolourpalette
	_output.ColourtoComponentMap = colourtoComponentMap
	//fmt.Println("Unique Colours =",Numberofcolours,"from palette:", chosencolourpalette)

}

// Run after controls and a steps block are completed to
// post process any data and provide downstream results
func _MakePalette_OneByOneAnalysis(_ctx context.Context, _input *MakePalette_OneByOneInput, _output *MakePalette_OneByOneOutput) {
}

// A block of tests to perform to validate that the sample was processed correctly
// Optionally, destructive tests can be performed to validate results on a
// dipstick basis
func _MakePalette_OneByOneValidation(_ctx context.Context, _input *MakePalette_OneByOneInput, _output *MakePalette_OneByOneOutput) {

}
func _MakePalette_OneByOneRun(_ctx context.Context, input *MakePalette_OneByOneInput) *MakePalette_OneByOneOutput {
	output := &MakePalette_OneByOneOutput{}
	_MakePalette_OneByOneSetup(_ctx, input)
	_MakePalette_OneByOneSteps(_ctx, input, output)
	_MakePalette_OneByOneAnalysis(_ctx, input, output)
	_MakePalette_OneByOneValidation(_ctx, input, output)
	return output
}

func MakePalette_OneByOneRunSteps(_ctx context.Context, input *MakePalette_OneByOneInput) *MakePalette_OneByOneSOutput {
	soutput := &MakePalette_OneByOneSOutput{}
	output := _MakePalette_OneByOneRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func MakePalette_OneByOneNew() interface{} {
	return &MakePalette_OneByOneElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &MakePalette_OneByOneInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _MakePalette_OneByOneRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &MakePalette_OneByOneInput{},
			Out: &MakePalette_OneByOneOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wunit.Make_units
)

type MakePalette_OneByOneElement struct {
	inject.CheckedRunner
}

type MakePalette_OneByOneInput struct {
	AutoRotate          bool
	Black               *wtype.LHComponent
	Cyan                *wtype.LHComponent
	Imagefilename       string
	Magenta             *wtype.LHComponent
	OutPlate            *wtype.LHPlate
	PalettePlate        *wtype.LHPlate
	PosterizeImage      bool
	PosterizeLevels     int
	Rotate              bool
	VolumeForFullcolour wunit.Volume
	Yellow              *wtype.LHComponent
}

type MakePalette_OneByOneOutput struct {
	Colours              []*wtype.LHComponent
	ColourtoComponentMap map[string]*wtype.LHComponent
	Numberofcolours      int
	Palette              color.Palette
}

type MakePalette_OneByOneSOutput struct {
	Data struct {
		ColourtoComponentMap map[string]*wtype.LHComponent
		Numberofcolours      int
		Palette              color.Palette
	}
	Outputs struct {
		Colours []*wtype.LHComponent
	}
}

func init() {
	if err := addComponent(component.Component{Name: "MakePalette_OneByOne",
		Constructor: MakePalette_OneByOneNew,
		Desc: component.ComponentDesc{
			Desc: "Generates instructions to make a pallette of all colours in an image\n",
			Path: "antha/component/an/Liquid_handling/PipetteImage/MakePalette_OnebyOne.an",
			Params: []component.ParamDesc{
				{Name: "AutoRotate", Desc: "", Kind: "Parameters"},
				{Name: "Black", Desc: "", Kind: "Inputs"},
				{Name: "Cyan", Desc: "", Kind: "Inputs"},
				{Name: "Imagefilename", Desc: "", Kind: "Parameters"},
				{Name: "Magenta", Desc: "", Kind: "Inputs"},
				{Name: "OutPlate", Desc: "InPlate *wtype.LHPlate\n", Kind: "Inputs"},
				{Name: "PalettePlate", Desc: "", Kind: "Inputs"},
				{Name: "PosterizeImage", Desc: "", Kind: "Parameters"},
				{Name: "PosterizeLevels", Desc: "", Kind: "Parameters"},
				{Name: "Rotate", Desc: "", Kind: "Parameters"},
				{Name: "VolumeForFullcolour", Desc: "", Kind: "Parameters"},
				{Name: "Yellow", Desc: "", Kind: "Inputs"},
				{Name: "Colours", Desc: "", Kind: "Outputs"},
				{Name: "ColourtoComponentMap", Desc: "", Kind: "Data"},
				{Name: "Numberofcolours", Desc: "", Kind: "Data"},
				{Name: "Palette", Desc: "Colournames []string\n", Kind: "Data"},
			},
		},
	}); err != nil {
		panic(err)
	}
}
