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
	"strconv"
)

// Input parameters for this protocol (data)

//InoculationVolume Volume
/*AntibioticVolume Volume
InducerVolume Volume
RepressorVolume Volume*/

//IncTemp Temperature
//IncTime Time

// Data which is returned from this protocol, and data types

// for making a palette of colours in advance

// Physical Inputs to this protocol with types

//InPlate *wtype.LHPlate
//Media *wtype.LHComponent
/*Antibiotic *wtype.LHComponent
Inducer *wtype.LHComponent
Repressor *wtype.LHComponent*/

// Physical outputs from this protocol with types

//PaletteColours 	[]*wtype.LHComponent// not used

func _TransformLivingPaletteRequirements() {

}

// Conditions to run on startup
func _TransformLivingPaletteSetup(_ctx context.Context, _input *TransformLivingPaletteInput) {

}

// The core process for this protocol, with the steps to be performed
// for every input
func _TransformLivingPaletteSteps(_ctx context.Context, _input *TransformLivingPaletteInput, _output *TransformLivingPaletteOutput) {

	var (
		ReactionTemp                wunit.Temperature = wunit.NewTemperature(25, "C")
		ReactionTime                wunit.Time        = wunit.NewTime(35, "min")
		CompetentCellPlateWell      string            = ""
		RecoveryPlateWell           string            = ""
		RecoveryTemp                wunit.Temperature = wunit.NewTemperature(37.0, "C")
		RecoveryTime                wunit.Time        = wunit.NewTime(2.0, "h")
		TransformationVolume        wunit.Volume      = wunit.NewVolume(2.0, "ul")
		PostPlasmidTemp             wunit.Temperature = wunit.NewTemperature(2.0, "C")
		PostPlasmidTime             wunit.Time        = wunit.NewTime(5, "min")
		CompetentCellTransferVolume wunit.Volume      = wunit.NewVolume(20.0, "ul")
		RecoveryPlateNumber         int               = 1

		PlatewithRecoveryMedia  *wtype.LHPlate = factory.GetPlateByType("DSW96_riser40")
		PlateWithCompetentCells *wtype.LHPlate = factory.GetPlateByType("pcrplate_with_cooler")
	)

	colourtoComponentMap := make(map[string]*wtype.LHComponent)

	wellpositions := PlateWithCompetentCells.AllWellPositions(wtype.BYCOLUMN)

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

	// make additional palette map
	positiontocolourmap = image.RemoveDuplicatesValuesfromMap(positiontocolourmap)

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
	for _, colour := range positiontocolourmap {

		//components := make([]*wtype.LHComponent, 0)

		component := componentmap[colourtostringmap[colour]]

		// for palette
		colourindex := chosencolourpalette.Index(colour)

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

				// set destination wells
				CompetentCellPlateWell = wellpositions[counter]
				RecoveryPlateWell = CompetentCellPlateWell

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
				//pixelSample := mixer.Sample(component, VolumePerWell)

				dnaSample := mixer.Sample(component, TransformationVolume)

				dnaSample.Type = wtype.LTCulture

				execute.Incubate(_ctx, dnaSample, ReactionTemp, ReactionTime, false)

				transformation := execute.MixTo(_ctx, PlateWithCompetentCells.Type, CompetentCellPlateWell, 1, dnaSample)

				transformation.Type = wtype.LTCulture

				execute.Incubate(_ctx, transformation, PostPlasmidTemp, PostPlasmidTime, false)

				transformationSample := mixer.Sample(transformation, CompetentCellTransferVolume)

				solution := execute.MixTo(_ctx, PlatewithRecoveryMedia.Type, RecoveryPlateWell, RecoveryPlateNumber, transformationSample)

				execute.Incubate(_ctx, solution, RecoveryTemp, RecoveryTime, true)
				//components = append(components,pixelSample)
				//solution := MixTo(OutPlate.Type, locationkey, 1,pixelSample)

				solutions = append(solutions, solution /*Incubate(solution,IncTemp,IncTime,true)*/)
			}

		} else {
			if component.CName != _input.Notthiscolour {

				_output.UniqueComponents = append(_output.UniqueComponents, component.CName)

				// set destination wells
				CompetentCellPlateWell = wellpositions[counter]
				RecoveryPlateWell = CompetentCellPlateWell

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

				//pixelSample := mixer.Sample(component, VolumePerWell)
				//components = append(components,pixelSample)
				//solution := MixTo(OutPlate.Type, locationkey, 1, pixelSample)

				dnaSample := mixer.Sample(component, TransformationVolume)

				dnaSample.Type = wtype.LTCulture

				execute.Incubate(_ctx, dnaSample, ReactionTemp, ReactionTime, false)

				transformation := execute.MixTo(_ctx, PlateWithCompetentCells.Type, CompetentCellPlateWell, 1, dnaSample)

				transformation.Type = wtype.LTCulture

				execute.Incubate(_ctx, transformation, PostPlasmidTemp, PostPlasmidTime, false)

				transformationSample := mixer.Sample(transformation, CompetentCellTransferVolume)

				solution := execute.MixTo(_ctx, PlatewithRecoveryMedia.Type, RecoveryPlateWell, RecoveryPlateNumber, transformationSample)

				execute.Incubate(_ctx, solution, RecoveryTemp, RecoveryTime, true)

				solutions = append(solutions, solution /*Incubate(solution,IncTemp,IncTime,true)*/)
				colourtoComponentMap[strconv.Itoa(colourindex)] = solution
			}
		}
	}

	_output.UniqueComponents = search.RemoveDuplicates(_output.UniqueComponents)
	fmt.Println("Unique Components:", _output.UniqueComponents)
	fmt.Println("number of unique components", len(_output.UniqueComponents))
	_output.Pixels = solutions
	_output.Palette = chosencolourpalette
	_output.ColourtoComponentMap = colourtoComponentMap
	_output.Numberofpixels = len(_output.Pixels)
	fmt.Println("Pixels =", _output.Numberofpixels)

}

// Run after controls and a steps block are completed to
// post process any data and provide downstream results
func _TransformLivingPaletteAnalysis(_ctx context.Context, _input *TransformLivingPaletteInput, _output *TransformLivingPaletteOutput) {
}

// A block of tests to perform to validate that the sample was processed correctly
// Optionally, destructive tests can be performed to validate results on a
// dipstick basis
func _TransformLivingPaletteValidation(_ctx context.Context, _input *TransformLivingPaletteInput, _output *TransformLivingPaletteOutput) {

}
func _TransformLivingPaletteRun(_ctx context.Context, input *TransformLivingPaletteInput) *TransformLivingPaletteOutput {
	output := &TransformLivingPaletteOutput{}
	_TransformLivingPaletteSetup(_ctx, input)
	_TransformLivingPaletteSteps(_ctx, input, output)
	_TransformLivingPaletteAnalysis(_ctx, input, output)
	_TransformLivingPaletteValidation(_ctx, input, output)
	return output
}

func TransformLivingPaletteRunSteps(_ctx context.Context, input *TransformLivingPaletteInput) *TransformLivingPaletteSOutput {
	soutput := &TransformLivingPaletteSOutput{}
	output := _TransformLivingPaletteRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func TransformLivingPaletteNew() interface{} {
	return &TransformLivingPaletteElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &TransformLivingPaletteInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _TransformLivingPaletteRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &TransformLivingPaletteInput{},
			Out: &TransformLivingPaletteOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wunit.Make_units
)

type TransformLivingPaletteElement struct {
	inject.CheckedRunner
}

type TransformLivingPaletteInput struct {
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

type TransformLivingPaletteOutput struct {
	ColourtoComponentMap map[string]*wtype.LHComponent
	Numberofpixels       int
	Palette              color.Palette
	Pixels               []*wtype.LHComponent
	UniqueComponents     []string
}

type TransformLivingPaletteSOutput struct {
	Data struct {
		ColourtoComponentMap map[string]*wtype.LHComponent
		Numberofpixels       int
		Palette              color.Palette
		UniqueComponents     []string
	}
	Outputs struct {
		Pixels []*wtype.LHComponent
	}
}

func init() {
	if err := addComponent(component.Component{Name: "TransformLivingPalette",
		Constructor: TransformLivingPaletteNew,
		Desc: component.ComponentDesc{
			Desc: "Generates instructions to pipette out a defined image onto a defined plate using a defined palette of coloured bacteria\n",
			Path: "src/github.com/antha-lang/antha/antha/component/an/Liquid_handling/PipetteImage/TransformLivingPalette.an",
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
				{Name: "ColourtoComponentMap", Desc: "", Kind: "Data"},
				{Name: "Numberofpixels", Desc: "", Kind: "Data"},
				{Name: "Palette", Desc: "for making a palette of colours in advance\n", Kind: "Data"},
				{Name: "Pixels", Desc: "", Kind: "Outputs"},
				{Name: "UniqueComponents", Desc: "", Kind: "Data"},
			},
		},
	}); err != nil {
		panic(err)
	}
}
