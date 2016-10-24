// Generates instructions to make a pallette of all colours in an image
package lib

import (
	"fmt"
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

func _MakePaletteRequirements() {

}

// Conditions to run on startup
func _MakePaletteSetup(_ctx context.Context, _input *MakePaletteInput) {

}

// The core process for this protocol, with the steps to be performed
// for every input
func _MakePaletteSteps(_ctx context.Context, _input *MakePaletteInput, _output *MakePaletteOutput) {

	if _input.PosterizeImage {
		_, _input.Imagefilename = image.Posterize(_input.Imagefilename, _input.PosterizeLevels)
	}

	//var chosencolourpalette color.Palette

	//chosencolourpalette := image.AvailablePalettes["Plan9"]

	//positiontocolourmap, _ := image.ImagetoPlatelayout(Imagefilename, OutPlate, &chosencolourpalette, Rotate)

	// make pallette of colours from image
	chosencolourpalette := image.MakeSmallPalleteFromImage(_input.Imagefilename, _input.OutPlate, _input.Rotate)

	positiontocolourmap, _, _ := image.ImagetoPlatelayout(_input.Imagefilename, _input.OutPlate, &chosencolourpalette, _input.Rotate, _input.AutoRotate)

	// remove duplicates
	positiontocolourmap = image.RemoveDuplicatesValuesfromMap(positiontocolourmap)

	//fmt.Println("positions", positiontocolourmap)

	solutions := make([]*wtype.LHComponent, 0)
	colourtoComponentMap := make(map[string]*wtype.LHComponent)

	counter := 0

	platenum := 1
	paletteplatepositions := _input.PalettePlate.AllWellPositions(wtype.BYCOLUMN)

	for _, colour := range positiontocolourmap {

		var solution *wtype.LHComponent

		colourindex := chosencolourpalette.Index(colour)

		if colour != nil {
			//components := make([]*wtype.LHComponent, 0)

			cmyk := image.ColourtoCMYK(colour)

			var maxuint8 uint8 = 255

			if cmyk.C == 0 && cmyk.Y == 0 && cmyk.M == 0 && cmyk.K == 0 {

				continue

			} else {

				nextwellpostion := paletteplatepositions[counter]

				if counter+1 == len(paletteplatepositions) {
					platenum = platenum + 1
					counter = 0
				} else {
					counter = counter + 1
				}
				if cmyk.C > 0 {

					cyanvol := wunit.NewVolume(((float64(cmyk.C) / float64(maxuint8)) * _input.VolumeForFullcolour.RawValue()), _input.VolumeForFullcolour.Unit().PrefixedSymbol())

					if cyanvol.RawValue() < 0.5 && cyanvol.Unit().PrefixedSymbol() == "ul" {
						cyanvol = wunit.NewVolume(0.5, "ul")
					}

					if cmyk.K == 0 && cmyk.M == 0 && cmyk.Y == 0 {
						_input.Cyan.Type = wtype.LTPostMix
					} else {
						_input.Cyan.Type = wtype.LTDISPENSEABOVE
					}

					cyanSample := mixer.Sample(_input.Cyan, cyanvol)

					solution = execute.MixTo(_ctx, _input.PalettePlate.Type, nextwellpostion, platenum, cyanSample)
					//solution = MixTo(PalettePlate.Type, position,1,cyanSample)

					//components = append(components, cyanSample)
				}

				if cmyk.Y > 0 {
					yellowvol := wunit.NewVolume(((float64(cmyk.Y) / float64(maxuint8)) * _input.VolumeForFullcolour.RawValue()), _input.VolumeForFullcolour.Unit().PrefixedSymbol())

					if yellowvol.RawValue() < 0.5 && yellowvol.Unit().PrefixedSymbol() == "ul" {
						yellowvol = wunit.NewVolume(0.5, "ul")
					}
					if cmyk.K == 0 && cmyk.M == 0 {
						_input.Yellow.Type = wtype.LTPostMix
					} else {
						_input.Yellow.Type = wtype.LTDISPENSEABOVE
					}

					yellowSample := mixer.Sample(_input.Yellow, yellowvol)

					if solution != nil {
						solution = execute.Mix(_ctx, solution, yellowSample)
					} else {
						solution = execute.MixTo(_ctx, _input.PalettePlate.Type, nextwellpostion, platenum, yellowSample)
						//solution = MixTo(PalettePlate.Type, position,1,yellowSample)
					}

					//components = append(components, yellowSample)
				}

				if cmyk.M > 0 {
					magentavol := wunit.NewVolume(((float64(cmyk.M) / float64(maxuint8)) * _input.VolumeForFullcolour.RawValue()), _input.VolumeForFullcolour.Unit().PrefixedSymbol())

					if magentavol.RawValue() < 0.5 && magentavol.Unit().PrefixedSymbol() == "ul" {
						magentavol = wunit.NewVolume(0.5, "ul")
					}

					if cmyk.K == 0 {
						_input.Magenta.Type = wtype.LTPostMix
					} else {
						_input.Magenta.Type = wtype.LTDISPENSEABOVE
					}

					magentaSample := mixer.Sample(_input.Magenta, magentavol)

					if solution != nil {
						solution = execute.Mix(_ctx, solution, magentaSample)
					} else {
						solution = execute.MixTo(_ctx, _input.PalettePlate.Type, nextwellpostion, platenum, magentaSample)
						//solution = MixTo(PalettePlate.Type, position,1,magentaSample)
					}

					//components = append(components, magentaSample)
				}

				if cmyk.K > 0 {
					blackvol := wunit.NewVolume(((float64(cmyk.K) / float64(maxuint8)) * _input.VolumeForFullcolour.RawValue()), _input.VolumeForFullcolour.Unit().PrefixedSymbol())

					if blackvol.RawValue() < 0.5 && blackvol.Unit().PrefixedSymbol() == "ul" {
						blackvol = wunit.NewVolume(0.5, "ul")
					}

					_input.Black.Type = wtype.LTPostMix

					blackSample := mixer.Sample(_input.Black, blackvol)

					if solution != nil {
						solution = execute.Mix(_ctx, solution, blackSample)
					} else {
						solution = execute.MixTo(_ctx, _input.PalettePlate.Type, nextwellpostion, platenum, blackSample)
						//solution = MixTo(PalettePlate.Type, position,1,blackSample)
					}

					//components = append(components, blackSample)
				}

				//solution := MixInto(PalettePlate, "", components...)
				solutions = append(solutions, solution)
				colourtoComponentMap[strconv.Itoa(colourindex)] = solution

			}

		}
	}

	_output.Colours = solutions
	_output.Numberofcolours = len(chosencolourpalette)
	_output.Palette = chosencolourpalette
	_output.ColourtoComponentMap = colourtoComponentMap
	fmt.Println("Unique Colours =", _output.Numberofcolours, "from palette:", chosencolourpalette)

}

// Run after controls and a steps block are completed to
// post process any data and provide downstream results
func _MakePaletteAnalysis(_ctx context.Context, _input *MakePaletteInput, _output *MakePaletteOutput) {
}

// A block of tests to perform to validate that the sample was processed correctly
// Optionally, destructive tests can be performed to validate results on a
// dipstick basis
func _MakePaletteValidation(_ctx context.Context, _input *MakePaletteInput, _output *MakePaletteOutput) {

}
func _MakePaletteRun(_ctx context.Context, input *MakePaletteInput) *MakePaletteOutput {
	output := &MakePaletteOutput{}
	_MakePaletteSetup(_ctx, input)
	_MakePaletteSteps(_ctx, input, output)
	_MakePaletteAnalysis(_ctx, input, output)
	_MakePaletteValidation(_ctx, input, output)
	return output
}

func MakePaletteRunSteps(_ctx context.Context, input *MakePaletteInput) *MakePaletteSOutput {
	soutput := &MakePaletteSOutput{}
	output := _MakePaletteRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func MakePaletteNew() interface{} {
	return &MakePaletteElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &MakePaletteInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _MakePaletteRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &MakePaletteInput{},
			Out: &MakePaletteOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wunit.Make_units
)

type MakePaletteElement struct {
	inject.CheckedRunner
}

type MakePaletteInput struct {
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

type MakePaletteOutput struct {
	Colours              []*wtype.LHComponent
	ColourtoComponentMap map[string]*wtype.LHComponent
	Numberofcolours      int
	Palette              color.Palette
}

type MakePaletteSOutput struct {
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
	if err := addComponent(component.Component{Name: "MakePalette",
		Constructor: MakePaletteNew,
		Desc: component.ComponentDesc{
			Desc: "Generates instructions to make a pallette of all colours in an image\n",
			Path: "antha/component/an/Liquid_handling/PipetteImage/MakePallete.an",
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
