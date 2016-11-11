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

func _MakePalette_OneByOne_RGBRequirements() {

}

// Conditions to run on startup
func _MakePalette_OneByOne_RGBSetup(_ctx context.Context, _input *MakePalette_OneByOne_RGBInput) {

}

// The core process for this protocol, with the steps to be performed
// for every input
func _MakePalette_OneByOne_RGBSteps(_ctx context.Context, _input *MakePalette_OneByOne_RGBInput, _output *MakePalette_OneByOne_RGBOutput) {

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

			r, g, b, _ := colour.RGBA()

			var maxuint8 uint8 = 255

			if r == 0 && g == 0 && b == 0 {

				continue

			} else {

				counter = counter + 1

				if r > 0 {

					redvol := wunit.NewVolume(((float64(r) / float64(maxuint8)) * _input.VolumeForFullcolour.RawValue()), _input.VolumeForFullcolour.Unit().PrefixedSymbol())

					if redvol.RawValue() < 10 && redvol.Unit().PrefixedSymbol() == "ul" {
						redvol.SetValue(10)
					}

					_input.Red.Type = wtype.LTPAINT

					redSample := mixer.Sample(_input.Red, redvol)
					components = append(components, redSample)
				}

				if g > 0 {
					greenvol := wunit.NewVolume(((float64(g) / float64(maxuint8)) * _input.VolumeForFullcolour.RawValue()), _input.VolumeForFullcolour.Unit().PrefixedSymbol())

					if greenvol.RawValue() < 10 && greenvol.Unit().PrefixedSymbol() == "ul" {
						greenvol.SetValue(10)
					}

					_input.Green.Type = wtype.LTPAINT

					greenSample := mixer.Sample(_input.Green, greenvol)
					components = append(components, greenSample)
				}

				if b > 0 {
					bluevol := wunit.NewVolume(((float64(b) / float64(maxuint8)) * _input.VolumeForFullcolour.RawValue()), _input.VolumeForFullcolour.Unit().PrefixedSymbol())

					if bluevol.RawValue() < 10 && bluevol.Unit().PrefixedSymbol() == "ul" {
						bluevol.SetValue(10)
					}

					_input.Blue.Type = wtype.LTPAINT

					blueSample := mixer.Sample(_input.Blue, bluevol)
					components = append(components, blueSample)
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
func _MakePalette_OneByOne_RGBAnalysis(_ctx context.Context, _input *MakePalette_OneByOne_RGBInput, _output *MakePalette_OneByOne_RGBOutput) {
}

// A block of tests to perform to validate that the sample was processed correctly
// Optionally, destructive tests can be performed to validate results on a
// dipstick basis
func _MakePalette_OneByOne_RGBValidation(_ctx context.Context, _input *MakePalette_OneByOne_RGBInput, _output *MakePalette_OneByOne_RGBOutput) {

}
func _MakePalette_OneByOne_RGBRun(_ctx context.Context, input *MakePalette_OneByOne_RGBInput) *MakePalette_OneByOne_RGBOutput {
	output := &MakePalette_OneByOne_RGBOutput{}
	_MakePalette_OneByOne_RGBSetup(_ctx, input)
	_MakePalette_OneByOne_RGBSteps(_ctx, input, output)
	_MakePalette_OneByOne_RGBAnalysis(_ctx, input, output)
	_MakePalette_OneByOne_RGBValidation(_ctx, input, output)
	return output
}

func MakePalette_OneByOne_RGBRunSteps(_ctx context.Context, input *MakePalette_OneByOne_RGBInput) *MakePalette_OneByOne_RGBSOutput {
	soutput := &MakePalette_OneByOne_RGBSOutput{}
	output := _MakePalette_OneByOne_RGBRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func MakePalette_OneByOne_RGBNew() interface{} {
	return &MakePalette_OneByOne_RGBElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &MakePalette_OneByOne_RGBInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _MakePalette_OneByOne_RGBRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &MakePalette_OneByOne_RGBInput{},
			Out: &MakePalette_OneByOne_RGBOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wunit.Make_units
)

type MakePalette_OneByOne_RGBElement struct {
	inject.CheckedRunner
}

type MakePalette_OneByOne_RGBInput struct {
	AutoRotate          bool
	Blue                *wtype.LHComponent
	Green               *wtype.LHComponent
	Imagefilename       string
	OutPlate            *wtype.LHPlate
	PalettePlate        *wtype.LHPlate
	PosterizeImage      bool
	PosterizeLevels     int
	Red                 *wtype.LHComponent
	Rotate              bool
	VolumeForFullcolour wunit.Volume
}

type MakePalette_OneByOne_RGBOutput struct {
	Colours              []*wtype.LHComponent
	ColourtoComponentMap map[string]*wtype.LHComponent
	Numberofcolours      int
	Palette              color.Palette
}

type MakePalette_OneByOne_RGBSOutput struct {
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
	if err := addComponent(component.Component{Name: "MakePalette_OneByOne_RGB",
		Constructor: MakePalette_OneByOne_RGBNew,
		Desc: component.ComponentDesc{
			Desc: "Generates instructions to make a pallette of all colours in an image\n",
			Path: "antha/component/an/Liquid_handling/PipetteImage/MakePalette_OnebyOne_RGB.an",
			Params: []component.ParamDesc{
				{Name: "AutoRotate", Desc: "", Kind: "Parameters"},
				{Name: "Blue", Desc: "", Kind: "Inputs"},
				{Name: "Green", Desc: "", Kind: "Inputs"},
				{Name: "Imagefilename", Desc: "", Kind: "Parameters"},
				{Name: "OutPlate", Desc: "InPlate *wtype.LHPlate\n", Kind: "Inputs"},
				{Name: "PalettePlate", Desc: "", Kind: "Inputs"},
				{Name: "PosterizeImage", Desc: "", Kind: "Parameters"},
				{Name: "PosterizeLevels", Desc: "", Kind: "Parameters"},
				{Name: "Red", Desc: "", Kind: "Inputs"},
				{Name: "Rotate", Desc: "", Kind: "Parameters"},
				{Name: "VolumeForFullcolour", Desc: "", Kind: "Parameters"},
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
