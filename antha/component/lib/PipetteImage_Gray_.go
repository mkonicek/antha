// Generates instructions to pipette out a defined image onto a defined plate by blending cyan magenta yellow and black dyes
package lib

import (
	"fmt"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/image"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/search"
	"github.com/antha-lang/antha/antha/anthalib/mixer"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/component"
	"github.com/antha-lang/antha/execute"
	"github.com/antha-lang/antha/inject"
	"github.com/disintegration/imaging"
	"golang.org/x/net/context"
)

// Input parameters for this protocol (data)

// as a proportion of 1 i.e. 0.5 == 50%. Below this it will be considered white
// above this value pure black will be dispensed

// Data which is returned from this protocol, and data types

// Physical Inputs to this protocol with types

//InPlate *wtype.LHPlate

// Physical outputs from this protocol with types

func _PipetteImage_GrayRequirements() {

}

// Conditions to run on startup
func _PipetteImage_GraySetup(_ctx context.Context, _input *PipetteImage_GrayInput) {

}

// The core process for this protocol, with the steps to be performed
// for every input
func _PipetteImage_GraySteps(_ctx context.Context, _input *PipetteImage_GrayInput, _output *PipetteImage_GrayOutput) {

	var blackvol wunit.Volume

	var maxuint8 uint8 = 255

	var minuint8 uint8

	var fullblackuint8 uint8

	_output.ShadesofGrey = make([]int, 0)

	chosencolourpalette := image.AvailablePalettes()["Gray"]

	if _input.CheckResizeAlgorithms {
		image.CheckAllResizealgorithms(_input.Imagefilename, _input.OutPlate, _input.Rotate, imaging.AllResampleFilters)
	}

	positiontocolourmap, _, newimagename := image.ImagetoPlatelayout(_input.Imagefilename, _input.OutPlate, &chosencolourpalette, _input.Rotate, _input.AutoRotate)

	// if posterize rerun
	if _input.PosterizeImage {
		_, _input.Imagefilename = image.Posterize(newimagename, _input.PosterizeLevels)

		positiontocolourmap, _, _ = image.ImagetoPlatelayout(_input.Imagefilename, _input.OutPlate, &chosencolourpalette, _input.Rotate, _input.AutoRotate)
	}

	solutions := make([]*wtype.LHComponent, 0)

	counter := 0
	skipped := 0
	fullblack := 0

	for locationkey, colour := range positiontocolourmap {

		//components := make([]*wtype.LHComponent, 0)

		var solution *wtype.LHComponent
		//var mixedsolution *wtype.LHComponent

		gray := image.ColourtoGrayscale(colour)

		if _input.Negative == false {
			gray.Y = maxuint8 - gray.Y
		}

		// adjust thresholds for mixing black and white based on user parameters
		minuint8 = uint8(_input.MinimumBlackpercentagethreshold * float64(maxuint8))

		fullblackuint8 = uint8(_input.MaxBlackPercentagethreshold * float64(maxuint8))

		//	fmt.Println("brand new minuint8 ",minuint8,"fullblackuint8 ", fullblackuint8)

		if gray.Y < minuint8 {
			if _input.SkipWhite {
				skipped = skipped + 1
				//	fmt.Println("skipping well:", skipped,locationkey)
			} else {
				whitevol := _input.VolumeForFullcolour
				_input.Diluent.Type, _ = wtype.LiquidTypeFromString(_input.NonMixingClass)

				waterSample := mixer.Sample(_input.Diluent, whitevol)
				solution = execute.MixTo(_ctx, _input.OutPlate.Type, locationkey, 1, waterSample)
				solutions = append(solutions, solution)
			}
			continue

		} else {

			counter = counter + 1

			// check if shade of grey has already been used in image
			greyindexinpalette := chosencolourpalette.Index(colour)

			alreadythere := search.Contains(_output.ShadesofGrey, greyindexinpalette)

			if alreadythere == false {
				_output.ShadesofGrey = append(_output.ShadesofGrey, greyindexinpalette)
			}

			if gray.Y < fullblackuint8 /*&& gray.Y >= minuint8*/ {
				watervol := wunit.NewVolume((float64(maxuint8-gray.Y) / float64(maxuint8) * _input.VolumeForFullcolour.RawValue()), _input.VolumeForFullcolour.Unit().PrefixedSymbol())
				//			fmt.Println("new well", locationkey, "water vol", watervol.ToString())
				// force hv tip choice
				if _input.OnlyHighVolumetips && watervol.RawValue() < 21 && watervol.Unit().PrefixedSymbol() == "ul" {
					watervol.SetValue(21)
				}
				waterSample := mixer.Sample(_input.Diluent, watervol)
				//components = append(components, waterSample)
				solution = execute.MixTo(_ctx, _input.OutPlate.Type, locationkey, 1, waterSample)
			}
			if gray.Y >= fullblackuint8 {
				fullblack = fullblack + 1
				//		fmt.Println("full colours:", fullblack)
				blackvol = _input.VolumeForFullcolour
			} else {
				blackvol = wunit.NewVolume((float64(gray.Y) / float64(maxuint8) * _input.VolumeForFullcolour.RawValue()), _input.VolumeForFullcolour.Unit().PrefixedSymbol())
			}

			//	fmt.Println("new well", locationkey, "black vol", blackvol.ToString())

			//Black.Type = wtype.LiquidTypeFromString("NeedToMix")

			if _input.DontMix {
				_input.Black.Type = wtype.LTDISPENSEABOVE
			} else if gray.Y >= fullblackuint8 {
				_input.Black.Type, _ = wtype.LiquidTypeFromString(_input.NonMixingClass)
			} else {
				_input.Black.Type, _ = wtype.LiquidTypeFromString(_input.MixingLiquidClass)
			}

			//fmt.Println("blackvol2",blackvol.ToString())
			if _input.OnlyHighVolumetips && blackvol.RawValue() < 21 && blackvol.Unit().PrefixedSymbol() == "ul" {
				blackvol.SetValue(21)
			}

			blackSample := mixer.Sample(_input.Black, blackvol)
			//components = append(components, blackSample)

			if solution != nil {
				solution = execute.Mix(_ctx, solution, blackSample)
			} else {
				solution = execute.MixTo(_ctx, _input.OutPlate.Type, locationkey, 1, blackSample)
			}
			solutions = append(solutions, solution)

		}
	}

	_output.NumberofShadesofGrey = len(_output.ShadesofGrey)
	_output.Pixels = solutions
	_output.Numberofpixels = len(_output.Pixels)
	fmt.Println("Pixels =", _output.Numberofpixels)
	_output.Fullblack = fullblack
	_output.Skipped = skipped

}

// Run after controls and a steps block are completed to
// post process any data and provide downstream results
func _PipetteImage_GrayAnalysis(_ctx context.Context, _input *PipetteImage_GrayInput, _output *PipetteImage_GrayOutput) {
}

// A block of tests to perform to validate that the sample was processed correctly
// Optionally, destructive tests can be performed to validate results on a
// dipstick basis
func _PipetteImage_GrayValidation(_ctx context.Context, _input *PipetteImage_GrayInput, _output *PipetteImage_GrayOutput) {

}
func _PipetteImage_GrayRun(_ctx context.Context, input *PipetteImage_GrayInput) *PipetteImage_GrayOutput {
	output := &PipetteImage_GrayOutput{}
	_PipetteImage_GraySetup(_ctx, input)
	_PipetteImage_GraySteps(_ctx, input, output)
	_PipetteImage_GrayAnalysis(_ctx, input, output)
	_PipetteImage_GrayValidation(_ctx, input, output)
	return output
}

func PipetteImage_GrayRunSteps(_ctx context.Context, input *PipetteImage_GrayInput) *PipetteImage_GraySOutput {
	soutput := &PipetteImage_GraySOutput{}
	output := _PipetteImage_GrayRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func PipetteImage_GrayNew() interface{} {
	return &PipetteImage_GrayElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &PipetteImage_GrayInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _PipetteImage_GrayRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &PipetteImage_GrayInput{},
			Out: &PipetteImage_GrayOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wunit.Make_units
)

type PipetteImage_GrayElement struct {
	inject.CheckedRunner
}

type PipetteImage_GrayInput struct {
	AutoRotate                      bool
	Black                           *wtype.LHComponent
	CheckResizeAlgorithms           bool
	Diluent                         *wtype.LHComponent
	DontMix                         bool
	Imagefilename                   string
	MaxBlackPercentagethreshold     float64
	MinimumBlackpercentagethreshold float64
	MixingLiquidClass               string
	Negative                        bool
	NonMixingClass                  string
	OnlyHighVolumetips              bool
	OutPlate                        *wtype.LHPlate
	PosterizeImage                  bool
	PosterizeLevels                 int
	Rotate                          bool
	SkipWhite                       bool
	VolumeForFullcolour             wunit.Volume
}

type PipetteImage_GrayOutput struct {
	Fullblack            int
	NumberofShadesofGrey int
	Numberofpixels       int
	Pixels               []*wtype.LHComponent
	ShadesofGrey         []int
	Skipped              int
}

type PipetteImage_GraySOutput struct {
	Data struct {
		Fullblack            int
		NumberofShadesofGrey int
		Numberofpixels       int
		ShadesofGrey         []int
		Skipped              int
	}
	Outputs struct {
		Pixels []*wtype.LHComponent
	}
}

func init() {
	if err := addComponent(component.Component{Name: "PipetteImage_Gray",
		Constructor: PipetteImage_GrayNew,
		Desc: component.ComponentDesc{
			Desc: "Generates instructions to pipette out a defined image onto a defined plate by blending cyan magenta yellow and black dyes\n",
			Path: "antha/component/an/Liquid_handling/PipetteImage/PipetteImage_Gray.an",
			Params: []component.ParamDesc{
				{Name: "AutoRotate", Desc: "", Kind: "Parameters"},
				{Name: "Black", Desc: "", Kind: "Inputs"},
				{Name: "CheckResizeAlgorithms", Desc: "", Kind: "Parameters"},
				{Name: "Diluent", Desc: "", Kind: "Inputs"},
				{Name: "DontMix", Desc: "", Kind: "Parameters"},
				{Name: "Imagefilename", Desc: "", Kind: "Parameters"},
				{Name: "MaxBlackPercentagethreshold", Desc: "above this value pure black will be dispensed\n", Kind: "Parameters"},
				{Name: "MinimumBlackpercentagethreshold", Desc: "as a proportion of 1 i.e. 0.5 == 50%. Below this it will be considered white\n", Kind: "Parameters"},
				{Name: "MixingLiquidClass", Desc: "", Kind: "Parameters"},
				{Name: "Negative", Desc: "", Kind: "Parameters"},
				{Name: "NonMixingClass", Desc: "", Kind: "Parameters"},
				{Name: "OnlyHighVolumetips", Desc: "", Kind: "Parameters"},
				{Name: "OutPlate", Desc: "InPlate *wtype.LHPlate\n", Kind: "Inputs"},
				{Name: "PosterizeImage", Desc: "", Kind: "Parameters"},
				{Name: "PosterizeLevels", Desc: "", Kind: "Parameters"},
				{Name: "Rotate", Desc: "", Kind: "Parameters"},
				{Name: "SkipWhite", Desc: "", Kind: "Parameters"},
				{Name: "VolumeForFullcolour", Desc: "", Kind: "Parameters"},
				{Name: "Fullblack", Desc: "", Kind: "Data"},
				{Name: "NumberofShadesofGrey", Desc: "", Kind: "Data"},
				{Name: "Numberofpixels", Desc: "", Kind: "Data"},
				{Name: "Pixels", Desc: "", Kind: "Outputs"},
				{Name: "ShadesofGrey", Desc: "", Kind: "Data"},
				{Name: "Skipped", Desc: "", Kind: "Data"},
			},
		},
	}); err != nil {
		panic(err)
	}
}
