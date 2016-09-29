// Generates instructions to pipette out a defined image onto a defined plate using a defined palette of coloured bacteria
package lib

import (
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/image"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	//"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/text"
	"fmt"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/search"
	"github.com/antha-lang/antha/antha/anthalib/mixer"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/component"
	"github.com/antha-lang/antha/execute"
	"github.com/antha-lang/antha/inject"
	"github.com/antha-lang/antha/microArch/factory"
	"golang.org/x/net/context"
	"image/color"
)

// Input parameters for this protocol (data)

//InoculationVolume Volume
/*AntibioticVolume Volume
InducerVolume Volume
RepressorVolume Volume*/

//IncTemp Temperature
//IncTime Time

// Data which is returned from this protocol, and data types

// Physical Inputs to this protocol with types

//InPlate *wtype.LHPlate
//Media *wtype.LHComponent
/*Antibiotic *wtype.LHComponent
Inducer *wtype.LHComponent
Repressor *wtype.LHComponent*/

// Physical outputs from this protocol with types

func _PipetteImage_livingRequirements() {

}

// Conditions to run on startup
func _PipetteImage_livingSetup(_ctx context.Context, _input *PipetteImage_livingInput) {

}

// The core process for this protocol, with the steps to be performed
// for every input
func _PipetteImage_livingSteps(_ctx context.Context, _input *PipetteImage_livingInput, _output *PipetteImage_livingOutput) {

	// make sub pallete if necessary
	var chosencolourpalette color.Palette
	var err error

	if _input.Subset {
		chosencolourpalette = image.MakeSubPallette(_input.Palettename, _input.Subsetnames)
	} else {
		chosencolourpalette = image.AvailablePalettes()[_input.Palettename]
	}

	// resize image to fit dimensions of plate and change each pixel to match closest colour from chosen palette
	// the output of this is a map of well positions to colours needed
	positiontocolourmap, _, _ := image.ImagetoPlatelayout(_input.Imagefilename, _input.OutPlate, &chosencolourpalette, _input.Rotate, _input.AutoRotate)

	colourtostringmap := image.AvailableComponentmaps()[_input.Palettename]

	// if the image will be printed using fluorescent proteins, 2 previews will be generated for the image (i) under UV light (ii) under visible light

	if _input.UVimage {
		uvmap := image.AvailableComponentmaps()[_input.Palettename]
		visiblemap := image.Visibleequivalentmaps()[_input.Palettename]

		if _input.Subset {
			uvmap = image.MakeSubMapfromMap(colourtostringmap, _input.Subsetnames)
			visiblemap = image.MakeSubMapfromMap(colourtostringmap, _input.Subsetnames)
		}
		image.PrintFPImagePreview(_input.Imagefilename, _input.OutPlate, _input.Rotate, visiblemap, uvmap)
	}

	// get components from factory
	componentmap := make(map[string]*wtype.LHComponent, 0)

	if _input.Subset {
		colourtostringmap = image.MakeSubMapfromMap(colourtostringmap, _input.Subsetnames)
	}

	for colourname := range colourtostringmap {

		componentname := colourtostringmap[colourname]

		// use template component instead
		var componenttopick *wtype.LHComponent

		if _input.ComponentType != nil {
			componenttopick = _input.ComponentType
		} else {
			componenttopick = factory.GetComponentByType("water")
		}
		componenttopick.CName = componentname

		componentmap[componentname] = componenttopick

	}
	//	fmt.Println(componentmap)

	solutions := make([]*wtype.LHComponent, 0)

	counter := 0
	_output.UniqueComponents = make([]string, 0)

	// loop through the position to colour map pipeting the correct coloured protein into each well
	for locationkey, colour := range positiontocolourmap {

		//components := make([]*wtype.LHComponent, 0)

		component := componentmap[colourtostringmap[colour]]

		// make sure liquid class is appropriate for cell culture in case this is not set elsewhere
		component.Type, err = wtype.LiquidTypeFromString(_input.UseLiquidClass) //wtype.LTCulture

		if err != nil {
			panic(err.Error())
		}

		//	fmt.Println(image.Colourcomponentmap[colour])

		// if the option to only print a single colour is not selected then the pipetting actions for all colours (apart from if not this colour is not empty) will follow
		if _input.OnlythisColour != "" {

			if image.Colourcomponentmap[colour] == _input.OnlythisColour {

				_output.UniqueComponents = append(_output.UniqueComponents, component.CName)

				counter = counter + 1
				//		fmt.Println("wells",OnlythisColour, counter)
				//mediaSample := mixer.SampleForTotalVolume(Media, VolumePerWell)
				//components = append(components,mediaSample)
				/*antibioticSample := mixer.Sample(Antibiotic, AntibioticVolume)
				components = append(components,antibioticSample)
				repressorSample := mixer.Sample(Repressor, RepressorVolume)
				components = append(components,repressorSample)
				inducerSample := mixer.Sample(Inducer, InducerVolume)
				components = append(components,inducerSample)*/
				pixelSample := mixer.Sample(component, _input.VolumePerWell)
				//components = append(components,pixelSample)
				solution := execute.MixTo(_ctx, _input.OutPlate.Type, locationkey, 1, pixelSample)

				solutions = append(solutions, solution /*Incubate(solution,IncTemp,IncTime,true)*/)
			}

		} else {
			if component.CName != _input.Notthiscolour {

				_output.UniqueComponents = append(_output.UniqueComponents, component.CName)

				counter = counter + 1
				//		fmt.Println("wells not ",Notthiscolour,counter)
				//mediaSample := mixer.SampleForTotalVolume(Media, VolumePerWell)
				//components = append(components,mediaSample)
				/*antibioticSample := mixer.Sample(Antibiotic, AntibioticVolume)
				components = append(components,antibioticSample)
				repressorSample := mixer.Sample(Repressor, RepressorVolume)
				components = append(components,repressorSample)
				inducerSample := mixer.Sample(Inducer, InducerVolume)
				components = append(components,inducerSample)*/

				component.Type, err = wtype.LiquidTypeFromString(_input.UseLiquidClass) //wtype.LTCulture

				if err != nil {
					panic(err.Error())
				}

				pixelSample := mixer.Sample(component, _input.VolumePerWell)
				//components = append(components,pixelSample)
				solution := execute.MixTo(_ctx, _input.OutPlate.Type, locationkey, 1, pixelSample)

				solutions = append(solutions, solution /*Incubate(solution,IncTemp,IncTime,true)*/)
			}
		}
	}

	_output.UniqueComponents = search.RemoveDuplicates(_output.UniqueComponents)
	fmt.Println("Unique Components:", _output.UniqueComponents)
	fmt.Println("number of unique components", len(_output.UniqueComponents))
	_output.Pixels = solutions

	_output.Numberofpixels = len(_output.Pixels)
	fmt.Println("Pixels =", _output.Numberofpixels)

}

// Run after controls and a steps block are completed to
// post process any data and provide downstream results
func _PipetteImage_livingAnalysis(_ctx context.Context, _input *PipetteImage_livingInput, _output *PipetteImage_livingOutput) {
}

// A block of tests to perform to validate that the sample was processed correctly
// Optionally, destructive tests can be performed to validate results on a
// dipstick basis
func _PipetteImage_livingValidation(_ctx context.Context, _input *PipetteImage_livingInput, _output *PipetteImage_livingOutput) {

}
func _PipetteImage_livingRun(_ctx context.Context, input *PipetteImage_livingInput) *PipetteImage_livingOutput {
	output := &PipetteImage_livingOutput{}
	_PipetteImage_livingSetup(_ctx, input)
	_PipetteImage_livingSteps(_ctx, input, output)
	_PipetteImage_livingAnalysis(_ctx, input, output)
	_PipetteImage_livingValidation(_ctx, input, output)
	return output
}

func PipetteImage_livingRunSteps(_ctx context.Context, input *PipetteImage_livingInput) *PipetteImage_livingSOutput {
	soutput := &PipetteImage_livingSOutput{}
	output := _PipetteImage_livingRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func PipetteImage_livingNew() interface{} {
	return &PipetteImage_livingElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &PipetteImage_livingInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _PipetteImage_livingRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &PipetteImage_livingInput{},
			Out: &PipetteImage_livingOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wunit.Make_units
)

type PipetteImage_livingElement struct {
	inject.CheckedRunner
}

type PipetteImage_livingInput struct {
	AutoRotate     bool
	ComponentType  *wtype.LHComponent
	Imagefilename  string
	Notthiscolour  string
	OnlythisColour string
	OutPlate       *wtype.LHPlate
	Palettename    string
	Rotate         bool
	Subset         bool
	Subsetnames    []string
	UVimage        bool
	UseLiquidClass string
	VolumePerWell  wunit.Volume
}

type PipetteImage_livingOutput struct {
	Numberofpixels   int
	Pixels           []*wtype.LHComponent
	UniqueComponents []string
}

type PipetteImage_livingSOutput struct {
	Data struct {
		Numberofpixels   int
		UniqueComponents []string
	}
	Outputs struct {
		Pixels []*wtype.LHComponent
	}
}

func init() {
	if err := addComponent(component.Component{Name: "PipetteImage_living",
		Constructor: PipetteImage_livingNew,
		Desc: component.ComponentDesc{
			Desc: "Generates instructions to pipette out a defined image onto a defined plate using a defined palette of coloured bacteria\n",
			Path: "antha/component/an/Liquid_handling/PipetteImage/PipetteLivingimage.an",
			Params: []component.ParamDesc{
				{Name: "AutoRotate", Desc: "", Kind: "Parameters"},
				{Name: "ComponentType", Desc: "InPlate *wtype.LHPlate\nMedia *wtype.LHComponent\nAntibiotic *wtype.LHComponent\n\tInducer *wtype.LHComponent\n\tRepressor *wtype.LHComponent\n", Kind: "Inputs"},
				{Name: "Imagefilename", Desc: "InoculationVolume Volume\nAntibioticVolume Volume\n\tInducerVolume Volume\n\tRepressorVolume Volume\n", Kind: "Parameters"},
				{Name: "Notthiscolour", Desc: "", Kind: "Parameters"},
				{Name: "OnlythisColour", Desc: "", Kind: "Parameters"},
				{Name: "OutPlate", Desc: "", Kind: "Inputs"},
				{Name: "Palettename", Desc: "", Kind: "Parameters"},
				{Name: "Rotate", Desc: "", Kind: "Parameters"},
				{Name: "Subset", Desc: "", Kind: "Parameters"},
				{Name: "Subsetnames", Desc: "", Kind: "Parameters"},
				{Name: "UVimage", Desc: "", Kind: "Parameters"},
				{Name: "UseLiquidClass", Desc: "", Kind: "Parameters"},
				{Name: "VolumePerWell", Desc: "", Kind: "Parameters"},
				{Name: "Numberofpixels", Desc: "", Kind: "Data"},
				{Name: "Pixels", Desc: "", Kind: "Outputs"},
				{Name: "UniqueComponents", Desc: "", Kind: "Data"},
			},
		},
	}); err != nil {
		panic(err)
	}
}
