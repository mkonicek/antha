// Generates instructions to make a pallette of all colours in an image
package lib

import (
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/image"
	"github.com/antha-lang/antha/antha/anthalib/mixer"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	//"github.com/antha-lang/antha/microArch/factory"
	//"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/search"
	"fmt"
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

func _AssemblePalette_OneByOne_RGBRequirements() {

}

// Conditions to run on startup
func _AssemblePalette_OneByOne_RGBSetup(_ctx context.Context, _input *AssemblePalette_OneByOne_RGBInput) {

}

type RBSdata struct {
	Seq      wtype.DNASequence
	Strength float64
}

type collection []RBSdata

func (rbscollection collection) Max() (rbs RBSdata) {

	//var maxSeq wtype.DNASequence
	var maxStrength float64
	for i := range rbscollection {
		if i == 0 {
			maxStrength = rbscollection[i].Strength
			//maxSeq = rbscollection[i].Seq
			rbs = rbscollection[i]
		} else if rbscollection[i].Strength > maxStrength {
			maxStrength = rbscollection[i].Strength
			//maxSeq = rbscollection[i].Seq
			rbs = rbscollection[i]
		}
	}
	return
}

// The core process for this protocol, with the steps to be performed
// for every input
func _AssemblePalette_OneByOne_RGBSteps(_ctx context.Context, _input *AssemblePalette_OneByOne_RGBInput, _output *AssemblePalette_OneByOne_RGBOutput) {

	var rbsstrengthdata []RBSdata = []RBSdata{
		{wtype.MakeLinearDNASequence("rbs1", "gggcgcgc"), 0.0},
		{wtype.MakeLinearDNASequence("rbs2", "gggcgcgc"), 2.0},
		{wtype.MakeLinearDNASequence("rbs3", "gggcgcgc"), 5.0},
		{wtype.MakeLinearDNASequence("rbs4", "gggcgcgc"), 10.0},
	}

	redname := _input.Red.CName
	greenname := _input.Green.CName
	bluename := _input.Blue.CName

	fmt.Println(rbsstrengthdata)
	//var chosencolourpalette color.Palette

	//chosencolourpalette := image.AvailablePalettes["Plan9"]

	//positiontocolourmap, _ := image.ImagetoPlatelayout(Imagefilename, PlateWithMasterMix, &chosencolourpalette, Rotate)

	if _input.PosterizeImage {
		_, _input.Imagefilename = image.Posterize(_input.Imagefilename, _input.PosterizeLevels)
	}

	// make palette of colours from image
	chosencolourpalette := image.MakeSmallPalleteFromImage(_input.Imagefilename, _input.PlateWithMasterMix, _input.Rotate)

	// make a map of colour to well coordinates
	positiontocolourmap, _, _ := image.ImagetoPlatelayout(_input.Imagefilename, _input.PlateWithMasterMix, &chosencolourpalette, _input.Rotate, _input.AutoRotate)

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

			//var maxuint8 uint8 = 255

			if r == 0 && g == 0 && b == 0 {

				continue

			} else {

				counter = counter + 1

				_input.Red.CName = fmt.Sprint(redname, "_RBS_", uint8(r))
				_input.Red.Type = wtype.LTPostMix

				redSample := mixer.Sample(_input.Red, _input.VolumeForeachColourPlasmid)
				components = append(components, redSample)

				_input.Green.CName = fmt.Sprint(greenname, "_RBS_", uint8(g))
				_input.Green.Type = wtype.LTPostMix

				greenSample := mixer.Sample(_input.Green, _input.VolumeForeachColourPlasmid)

				components = append(components, greenSample)

				_input.Blue.CName = fmt.Sprint(bluename, "_RBS_", uint8(b))
				_input.Blue.Type = wtype.LTPostMix

				blueSample := mixer.Sample(_input.Blue, _input.VolumeForeachColourPlasmid)

				components = append(components, blueSample)

				solution := execute.MixInto(_ctx, _input.PalettePlate, "", components...)

				/*
					dnaSample := mixer.Sample(solution,wunit.NewVolume(1,"ul"))

					transformation := MixInto(factory.GetPlateByType("pcrplate_with_cooler"),"",dnaSample)

					transformationSample := mixer.Sample(transformation,wunit.NewVolume(20,"ul"))

					recovery := MixInto(factory.GetPlateByType("DSW96_riser"),"",transformationSample)

					solutions = append(solutions, recovery)
					colourtoComponentMap[strconv.Itoa(colourindex)] = recovery
				*/

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
func _AssemblePalette_OneByOne_RGBAnalysis(_ctx context.Context, _input *AssemblePalette_OneByOne_RGBInput, _output *AssemblePalette_OneByOne_RGBOutput) {
}

// A block of tests to perform to validate that the sample was processed correctly
// Optionally, destructive tests can be performed to validate results on a
// dipstick basis
func _AssemblePalette_OneByOne_RGBValidation(_ctx context.Context, _input *AssemblePalette_OneByOne_RGBInput, _output *AssemblePalette_OneByOne_RGBOutput) {

}
func _AssemblePalette_OneByOne_RGBRun(_ctx context.Context, input *AssemblePalette_OneByOne_RGBInput) *AssemblePalette_OneByOne_RGBOutput {
	output := &AssemblePalette_OneByOne_RGBOutput{}
	_AssemblePalette_OneByOne_RGBSetup(_ctx, input)
	_AssemblePalette_OneByOne_RGBSteps(_ctx, input, output)
	_AssemblePalette_OneByOne_RGBAnalysis(_ctx, input, output)
	_AssemblePalette_OneByOne_RGBValidation(_ctx, input, output)
	return output
}

func AssemblePalette_OneByOne_RGBRunSteps(_ctx context.Context, input *AssemblePalette_OneByOne_RGBInput) *AssemblePalette_OneByOne_RGBSOutput {
	soutput := &AssemblePalette_OneByOne_RGBSOutput{}
	output := _AssemblePalette_OneByOne_RGBRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func AssemblePalette_OneByOne_RGBNew() interface{} {
	return &AssemblePalette_OneByOne_RGBElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &AssemblePalette_OneByOne_RGBInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _AssemblePalette_OneByOne_RGBRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &AssemblePalette_OneByOne_RGBInput{},
			Out: &AssemblePalette_OneByOne_RGBOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wunit.Make_units
)

type AssemblePalette_OneByOne_RGBElement struct {
	inject.CheckedRunner
}

type AssemblePalette_OneByOne_RGBInput struct {
	AutoRotate                 bool
	Blue                       *wtype.LHComponent
	Green                      *wtype.LHComponent
	Imagefilename              string
	PalettePlate               *wtype.LHPlate
	PlateWithMasterMix         *wtype.LHPlate
	PosterizeImage             bool
	PosterizeLevels            int
	Red                        *wtype.LHComponent
	Rotate                     bool
	VolumeForeachColourPlasmid wunit.Volume
}

type AssemblePalette_OneByOne_RGBOutput struct {
	Colours              []*wtype.LHComponent
	ColourtoComponentMap map[string]*wtype.LHComponent
	Numberofcolours      int
	Palette              color.Palette
}

type AssemblePalette_OneByOne_RGBSOutput struct {
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
	if err := addComponent(component.Component{Name: "AssemblePalette_OneByOne_RGB",
		Constructor: AssemblePalette_OneByOne_RGBNew,
		Desc: component.ComponentDesc{
			Desc: "Generates instructions to make a pallette of all colours in an image\n",
			Path: "antha/component/an/Liquid_handling/PipetteImage/AssemblePalette.an",
			Params: []component.ParamDesc{
				{Name: "AutoRotate", Desc: "", Kind: "Parameters"},
				{Name: "Blue", Desc: "", Kind: "Inputs"},
				{Name: "Green", Desc: "", Kind: "Inputs"},
				{Name: "Imagefilename", Desc: "", Kind: "Parameters"},
				{Name: "PalettePlate", Desc: "", Kind: "Inputs"},
				{Name: "PlateWithMasterMix", Desc: "InPlate *wtype.LHPlate\n", Kind: "Inputs"},
				{Name: "PosterizeImage", Desc: "", Kind: "Parameters"},
				{Name: "PosterizeLevels", Desc: "", Kind: "Parameters"},
				{Name: "Red", Desc: "", Kind: "Inputs"},
				{Name: "Rotate", Desc: "", Kind: "Parameters"},
				{Name: "VolumeForeachColourPlasmid", Desc: "", Kind: "Parameters"},
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
