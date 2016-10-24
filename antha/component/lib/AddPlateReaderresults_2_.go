package lib

import (
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/montanaflynn/stats"
	// "github.com/antha-lang/antha/antha/anthalib/mixer"
	"github.com/antha-lang/antha/microArch/driver/liquidhandling"
	"github.com/antha-lang/antha/microArch/factory"
	// "github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/plot"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/Parser"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/Pubchem"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/buffers"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/doe"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/platereader"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/plot"
	// "path/filepath"
	// antha "github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/AnthaPath"
	"fmt"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/component"
	"github.com/antha-lang/antha/execute"
	"github.com/antha-lang/antha/inject"
	"golang.org/x/net/context"
	"strconv"
	"strings"
)

// Input parameters for this protocol (data)

//= "250516CCFbubbles/260516ccfbubbles.xlsx" //"lhdoe110216_postspin_postshake.xlsx"
//= 0                                        //PRESHAKEPRESPIN
//= "250516CCFbubbles/240516DXCFFDoeoutputgilsonright_TEST.xlsx"
//= "JMP"
//= "250516CCFbubbles/2501516bubblesresults.xlsx"

// = 472
//= "Abs Spectrum"

//= []string{"P24"}

// = false
//= []string{"J5"}

//= []string{"AbsVLV"}

//              = false
// = map[string][]string{	"AbsVLV": []string{""},}
//= wunit.NewVolume(20, "ul")

// of target molecule at wavelength
//= 20330
//= 0.0002878191305957933

// validation requirements

// Data which is returned from this protocol, and data types

// Physical Inputs to this protocol with types

// Physical outputs from this protocol with types

func _AddPlateReaderresults_2Requirements() {
}

// Conditions to run on startup
func _AddPlateReaderresults_2Setup(_ctx context.Context, _input *AddPlateReaderresults_2Input) {
}

// The core process for this protocol, with the steps to be performed
// for every input
func _AddPlateReaderresults_2Steps(_ctx context.Context, _input *AddPlateReaderresults_2Input, _output *AddPlateReaderresults_2Output) {

	var actualconcentrations = make(map[string]wunit.Concentration)
	_output.ResponsetoManualValuesmap = make(map[string][]float64)

	//var volumetovalues = make(map[wunit.Volume][]float64)
	//var testsolstovalues = make(map[string]map[wunit.Volume][]float64)

	molecule, err := pubchem.MakeMolecule(_input.Molecule.CName)
	if err != nil {
		execute.Errorf(_ctx, err.Error())
	}

	Molecularweight := molecule.MolecularWeight

	var marsdata parser.MarsData

	marsdata, err = parser.ParseMarsXLSXOutput(_input.MarsResultsFileXLSX, _input.Sheet)
	if err != nil {
		execute.Errorf(_ctx, err.Error())
	}

	// range through pairing up wells from mars output and doe design

	var runs []doe.Run

	// find out int factors from liquidhandling policies
	policyitemmap := liquidhandling.MakePolicyItems()
	intfactors := make([]string, 0)

	for key, val := range policyitemmap {

		if val.Type.Name() == "int" {
			intfactors = append(intfactors, key)
		}
	}

	if _input.DesignFiletype == "DX" {
		runs, err = doe.RunsFromDXDesign(_input.DesignFile, intfactors)
		if err != nil {
			panic(err)
		}
	} else if _input.DesignFiletype == "JMP" {
		runs, err = doe.RunsFromJMPDesign(_input.DesignFile, []int{}, []int{}, intfactors)
		if err != nil {
			panic(err)
		}
	}

	_output.BlankValues = make([]float64, 0)

	for i := range _input.Blanks {
		blankValue, _ := marsdata.ReadingsAsAverage(_input.Blanks[i], 1, _input.Wavelength, _input.ReadingTypeinMarsFile)
		_output.BlankValues = append(_output.BlankValues, blankValue)
	}

	runswithresponses := make([]doe.Run, 0)

	for k, run := range runs {

		// values for r2 to reset each run

		//xvalues := make([]float64, 0)
		//yvalues := make([]float64, 0)

		// add origin
		//xvalues = append(xvalues, 0.0)
		//yvalues = append(yvalues, 0.0)

		for _, response := range _input.Responsecolumnstofill {

			var samples []string
			var manualsamples []string
			var ManualValues = make([]float64, 0)
			var manual float64
			var absorbance wtype.Absorbance
			var manualabsorbance wtype.Absorbance
			//var actualconcreplicates = make([]float64, 0)
			var manualCorrectnessFactorValues = make([]float64, 0)
			var correctnessFactorValues = make([]float64, 0)

			// intialise
			Responsecolumntofill := response

			experimentalvolumeinterface, err := runs[k].GetAdditionalInfo("Volume") //  ResponseToVolumeMap[response]

			experimentalvolumestr := experimentalvolumeinterface.(string)

			//experimentalvolumestr = strings.TrimSpace(experimentalvolumestr)

			var volandunit []string

			if strings.Count(experimentalvolumestr, " ") == 1 {
				volandunit = strings.Split(experimentalvolumestr, " ")
			} else if strings.Count(experimentalvolumestr, "ul") == 1 && strings.HasSuffix(experimentalvolumestr, "ul") {
				volandunit = []string{strings.Trim(experimentalvolumestr, "ul"), "ul"}
			}

			vol, err := strconv.ParseFloat(strings.TrimSpace(volandunit[0]), 64)

			if err != nil {
				execute.Errorf(_ctx, err.Error())
			}

			experimentalvolume := wunit.NewVolume(vol, strings.TrimSpace(volandunit[1]))

			actualconcentrations[experimentalvolume.ToString()] = buffers.DiluteBasedonMolecularWeight(Molecularweight, _input.StockconcinMperL, experimentalvolume, _input.Diluent.CName, wunit.NewVolume(_input.Stockvol.RawValue()-experimentalvolume.RawValue(), "ul"))

			//locationHeaders := ResponsetoLocationMap[response]

			//  manual pipetting well
			if _input.ManualComparison {

				manualwell := _input.VolumeToManualwells[experimentalvolumestr][0] // 1st well of array only

				manual, _ = marsdata.ReadingsAsAverage(manualwell, 1, _input.Wavelength, _input.ReadingTypeinMarsFile)

				run = doe.AddNewResponseFieldandValue(run, Responsecolumntofill+" Manual Raw average "+strconv.Itoa(_input.Wavelength), manual)

				manualsamples = _input.VolumeToManualwells[experimentalvolumestr]

				for i := range manualsamples {
					manualvalue, _ := marsdata.ReadingsAsAverage(manualsamples[i], 1, _input.Wavelength, _input.ReadingTypeinMarsFile)
					ManualValues = append(ManualValues, manualvalue)
				}

				_output.ResponsetoManualValuesmap[experimentalvolumestr] = ManualValues

			}

			// then per replicate ...

			//for i, locationheader := range locationHeaders {
			well, err := runs[k].GetAdditionalInfo("Location")
			if err != nil {
				panic(err)
			}

			// check optimal difference for each well

			//Responsecolumntofill = response + "replicate_" + strconv.Itoa(i+1)

			if _input.FindOptWavelength {
				_output.MeasuredOptimalWavelength = marsdata.FindOptimalWavelength(well.(string), _input.Blanks[0], "Raw Data")
				//measuredoptimalwavelengths = append(measuredoptimalwavelengths, meassuredoptwavelength)

			}

			rawaverage, err := marsdata.ReadingsAsAverage(well.(string), 1, _input.Wavelength, _input.ReadingTypeinMarsFile)

			run = doe.AddNewResponseFieldandValue(run, Responsecolumntofill+" Raw average "+strconv.Itoa(_input.Wavelength), rawaverage)

			// blank correct

			samples = []string{well.(string)}

			blankcorrected, err := marsdata.BlankCorrect(samples, _input.Blanks, _input.Wavelength, _input.ReadingTypeinMarsFile)

			run = doe.AddNewResponseFieldandValue(run, Responsecolumntofill+" BlankCorrected "+strconv.Itoa(_input.Wavelength), blankcorrected)

			// path length correct
			pathlength, err := platereader.EstimatePathLength(factory.GetPlateByType("greiner384_riser"), wunit.NewVolume(_input.Stockvol.RawValue()+experimentalvolume.RawValue(), "ul"))

			if err != nil {
				panic(err)
			}

			run = doe.AddNewResponseFieldandValue(run, Responsecolumntofill+" pathlength "+strconv.Itoa(_input.Wavelength), pathlength.ToString())

			absorbance.Reading = blankcorrected

			pathlengthcorrect := platereader.PathlengthCorrect(pathlength, absorbance)

			run = doe.AddNewResponseFieldandValue(run, Responsecolumntofill+" Pathlength corrected "+strconv.Itoa(_input.Wavelength), pathlengthcorrect.Reading)

			// molar absorbtivity of tartazine at 472nm is 20330
			// http://www.biochrom.co.uk/faq/8/119/what-is-the-limit-of-detection-of-the-zenyth-200.html

			actualconc := platereader.Concentration(pathlengthcorrect, _input.Extinctioncoefficient)

			run = doe.AddNewResponseFieldandValue(run, Responsecolumntofill+"ActualConc", actualconc.SIValue())

			// calculate correctness factor based on expected conc

			expectedconc := actualconcentrations[experimentalvolume.ToString()]
			correctnessfactor := actualconc.SIValue() / expectedconc.SIValue()

			run = doe.AddNewResponseFieldandValue(run, Responsecolumntofill+" ExpectedConc "+strconv.Itoa(_input.Wavelength), expectedconc.SIValue())
			run = doe.AddNewResponseFieldandValue(run, Responsecolumntofill+" CorrectnessFactor "+strconv.Itoa(_input.Wavelength), correctnessfactor)
			correctnessFactorValues = append(correctnessFactorValues, correctnessfactor)

			//xvalues = append(xvalues, expectedconc.SIValue())
			//yvalues = append(yvalues, actualconc.SIValue())
			//actualconcreplicates = append(actualconcreplicates, actualconc.SIValue())

			// add comparison to manually pipetted wells
			if _input.ManualComparison {
				manualblankcorrected, _ := marsdata.BlankCorrect(manualsamples, _input.Blanks, _input.Wavelength, _input.ReadingTypeinMarsFile)
				manualabsorbance.Reading = manualblankcorrected
				manualpathlengthcorrect := platereader.PathlengthCorrect(pathlength, manualabsorbance)
				manualactualconc := platereader.Concentration(manualpathlengthcorrect, _input.Extinctioncoefficient)
				run = doe.AddNewResponseFieldandValue(run, Responsecolumntofill+"ManualActualConc", manualactualconc.SIValue())
				manualcorrectnessfactor := actualconc.SIValue() / manualactualconc.SIValue()
				manualCorrectnessFactorValues = append(manualCorrectnessFactorValues, manualcorrectnessfactor)
				run = doe.AddNewResponseFieldandValue(run, Responsecolumntofill+" ManualCorrectnessFactor "+strconv.Itoa(_input.Wavelength), manualcorrectnessfactor)
			}

			// process replicates into mean and cv
			//mean := stats.Mean(actualconcreplicates)
			//stdev := stats.StdDevS(actualconcreplicates)

			//cv := stdev / mean

			//run = doe.AddNewResponseFieldandValue(run, response+"_Average_ActualConc", mean)
			//run = doe.AddNewResponseFieldandValue(run, response+"_CV_ActualConc", cv)

			// average of correctness factor values
			//meanCF := stats.Mean(correctnessFactorValues)
			//run = doe.AddNewResponseFieldandValue(run, response+"_Average_CorrectnessFactor", meanCF)

			/*	if ManualComparison {

				meanManCF := stats.Mean(manualCorrectnessFactorValues)
				run = doe.AddNewResponseFieldandValue(run, response+"_Average_ManualCorrectnessFactor", meanManCF)
				}
			*/

		}

		run = doe.AddNewResponseFieldandValue(run, "Runorder", k)

		//rsquared := plot.Rsquared("Expected Conc", xvalues, "Actual Conc", yvalues)
		//run.AddResponseValue("R2", rsquared)

		//xygraph := plot.Plot(xvalues, [][]float64{yvalues})
		//filenameandextension := strings.Split(OutputFilename, ".")
		//plot.Export(xygraph, filenameandextension[0]+".png")

		runswithresponses = append(runswithresponses, run)
	}

	doe.XLSXFileFromRuns(runswithresponses, _input.OutputFilename, _input.DesignFiletype)

	_output.Runs = runswithresponses

}

// Run after controls and a steps block are completed to
// post process any data and provide downstream results
func _AddPlateReaderresults_2Analysis(_ctx context.Context, _input *AddPlateReaderresults_2Input, _output *AddPlateReaderresults_2Output) {

	_output.Errors = make([]string, 0)

	xvalues := make([]float64, 0)
	yvalues := make([]float64, 0)

	// add origin
	xvalues = append(xvalues, 0.0)
	yvalues = append(yvalues, 0.0)

	fmt.Println("in analysis")

	if len(_output.Runs) == 0 {
		execute.Errorf(_ctx, "no runs")
	}
	// 1. now calculate r2 and plot results
	for i, runwithresponses := range _output.Runs {
		// values for r2 to reset each run

		// get response value and check if it's a float64 type
		expectedconc, err := runwithresponses.GetResponseValue("Absorbance ExpectedConc " + strconv.Itoa(_input.Wavelength))

		if err != nil {
			_output.Errors = append(_output.Errors, err.Error())
		}

		expectedconcfloat, floattrue := expectedconc.(float64)
		// if float64 is true
		if floattrue {
			xvalues = append(xvalues, expectedconcfloat)
		} else {
			execute.Errorf(_ctx, "Run"+fmt.Sprint(i, runwithresponses)+" ExpectedConc:"+fmt.Sprint(expectedconcfloat))
		}

		// get response value and check if it's a float64 type
		actualconc, err := runwithresponses.GetResponseValue("AbsorbanceActualConc")

		if err != nil {
			fmt.Println(err.Error())
			_output.Errors = append(_output.Errors, err.Error())
		}

		actualconcfloat, floattrue := actualconc.(float64)

		if floattrue {
			yvalues = append(yvalues, actualconcfloat)
		} else {
			fmt.Println(err.Error())
			execute.Errorf(_ctx, " ActualConc:"+fmt.Sprint(actualconcfloat))
		}

	}

	_output.R2, _output.Variance, _output.Formula = plot.Rsquared("Expected Conc", xvalues, "Actual Conc", yvalues)
	//run.AddResponseValue("R2", rsquared)

	xygraph, err := plot.Plot(xvalues, [][]float64{yvalues})
	if err != nil {
		_output.Errors = append(_output.Errors, err.Error())
	}
	filenameandextension := strings.Split(_input.OutputFilename, ".")
	plot.Export(xygraph, "10cm", "10cm", filenameandextension[0]+"_plot"+".png")

	// reset
	xvalues = make([]float64, 0)
	yvalues = make([]float64, 0)

	// add origin
	xvalues = append(xvalues, 0.0)
	yvalues = append(yvalues, 0.0)

	// 2. now plot correctnessfactor
	for i, runwithresponses := range _output.Runs {
		// values for r2 to reset each run

		// get response value and check if it's a float64 type
		expectedconc, err := runwithresponses.GetResponseValue("Absorbance ExpectedConc " + strconv.Itoa(_input.Wavelength))

		if err != nil {
			_output.Errors = append(_output.Errors, err.Error())
		}

		expectedconcfloat, floattrue := expectedconc.(float64)
		// if float64 is true
		if floattrue {
			xvalues = append(xvalues, expectedconcfloat)
		} else {
			execute.Errorf(_ctx, "Run"+fmt.Sprint(i, runwithresponses)+" ExpectedConc:"+fmt.Sprint(expectedconcfloat))
		}

		// get response value and check if it's a float64 type
		correctness, err := runwithresponses.GetResponseValue("Absorbance CorrectnessFactor " + strconv.Itoa(_input.Wavelength))

		if err != nil {
			fmt.Println(err.Error())
			_output.Errors = append(_output.Errors, err.Error())
		}

		correctnessfloat, floattrue := correctness.(float64)

		if floattrue {
			yvalues = append(yvalues, correctnessfloat)
		} else {
			fmt.Println(err.Error())
			execute.Errorf(_ctx, " Absorbance CorrectnessFactor:"+fmt.Sprint(correctnessfloat))
		}

	}

	_output.R2_CorrectnessFactor, _, _ = plot.Rsquared("Expected Conc", xvalues, "Correctness Factor", yvalues)
	//run.AddResponseValue("R2", rsquared)

	correctnessgraph, err := plot.Plot(xvalues, [][]float64{yvalues})

	if err != nil {
		_output.Errors = append(_output.Errors, err.Error())
	}

	plot.Export(correctnessgraph, "10cm", "10cm", filenameandextension[0]+"_correctnessfactor"+".png")

	// reset
	xvalues = make([]float64, 0)
	yvalues = make([]float64, 0)

	// add origin
	xvalues = append(xvalues, 0.0)
	yvalues = append(yvalues, 0.0)

	// 3. now look for systematic errors
	for i, runwithresponses := range _output.Runs {
		// values for r2 to reset each run

		// get response value and check if it's a float64 type
		runorder, err := runwithresponses.GetResponseValue("Runorder")

		if err != nil {
			_output.Errors = append(_output.Errors, err.Error())
		}

		runorderint, inttrue := runorder.(int)
		// if int is true
		if inttrue {
			xvalues = append(xvalues, float64(runorderint))
		} else {
			execute.Errorf(_ctx, "Run"+fmt.Sprint(i, runwithresponses)+" Run Order:"+fmt.Sprint(runorderint), " not an int")
		}

		// get response value and check if it's a float64 type
		actualconc, err := runwithresponses.GetResponseValue("AbsorbanceActualConc")

		if err != nil {
			fmt.Println(err.Error())
			_output.Errors = append(_output.Errors, err.Error())
		}

		actualconcfloat, floattrue := actualconc.(float64)

		if floattrue {
			yvalues = append(yvalues, actualconcfloat)
		} else {
			fmt.Println(err.Error())
			execute.Errorf(_ctx, " ActualConc:"+fmt.Sprint(actualconcfloat))
		}
	}

	runorderconcgraph, err := plot.Plot(xvalues, [][]float64{yvalues})

	if err != nil {
		_output.Errors = append(_output.Errors, err.Error())
	}

	plot.Export(runorderconcgraph, "10cm", "10cm", filenameandextension[0]+"_runorder"+".png")

	// reset
	xvalues = make([]float64, 0)
	yvalues = make([]float64, 0)

	// add origin
	xvalues = append(xvalues, 0.0)
	yvalues = append(yvalues, 0.0)

	// 4.  now look for systematic errors with correctness factor
	for i, runwithresponses := range _output.Runs {
		// values for r2 to reset each run

		// get response value and check if it's a float64 type
		runorder, err := runwithresponses.GetResponseValue("Runorder")

		if err != nil {
			_output.Errors = append(_output.Errors, err.Error())
		}

		runorderint, inttrue := runorder.(int)
		// if int is true
		if inttrue {
			xvalues = append(xvalues, float64(runorderint))
		} else {
			execute.Errorf(_ctx, "Run"+fmt.Sprint(i, runwithresponses)+" Run Order:"+fmt.Sprint(runorderint), " not an int")
		}

		// get response value and check if it's a float64 type
		correctness, err := runwithresponses.GetResponseValue("Absorbance CorrectnessFactor " + strconv.Itoa(_input.Wavelength))

		if err != nil {
			fmt.Println(err.Error())
			_output.Errors = append(_output.Errors, err.Error())
		}

		correctnessfloat, floattrue := correctness.(float64)

		if floattrue {
			yvalues = append(yvalues, correctnessfloat)
		} else {
			fmt.Println(err.Error())
			execute.Errorf(_ctx, " Absorbance CorrectnessFactor:"+fmt.Sprint(correctnessfloat))
		}

	}

	runordercorrectnessgraph, err := plot.Plot(xvalues, [][]float64{yvalues})

	if err != nil {
		_output.Errors = append(_output.Errors, err.Error())
	}

	plot.Export(runordercorrectnessgraph, "10cm", "10cm", filenameandextension[0]+"_runorder_correctnessfactor"+".png")

	// 5. workout CV for each volume
	replicateactualconcmap := make(map[string][]float64)
	_output.VolumeToActualConc = make(map[string]Dataset)
	replicatevalues := make([]float64, 0)

	replicatecorrectnessmap := make(map[string][]float64)
	correctnessvalues := make([]float64, 0)
	_output.VolumeToCorrectnessFactor = make(map[string]Dataset)

	//counter := 0

	// make map of replicate values for Actual Conc
	for _, runwithresponses := range _output.Runs {

		volstr, err := runwithresponses.GetAdditionalInfo("Volume")

		if err != nil {
			execute.Errorf(_ctx, err.Error())
		}

		/*
			repstr, err := runwithresponses.GetAdditionalInfo("Replicate")

			if err != nil {
				Errorf(err.Error())
			}
		*/
		actualconc, err := runwithresponses.GetResponseValue("AbsorbanceActualConc")

		if err != nil {
			execute.Errorf(_ctx, err.Error())
		}

		/*rep, err := strconv.Atoi(repstr.(string))

		if err != nil {
			Errorf(err.Error())
		}
		*/

		// Actual Conc map
		if _, found := replicateactualconcmap[volstr.(string)]; found /*&& rep == counter*/ {
			replicatevalues = replicateactualconcmap[volstr.(string)]
			replicatevalues = append(replicatevalues, actualconc.(float64))
			replicateactualconcmap[volstr.(string)] = replicatevalues
			replicatevalues = make([]float64, 0)
			//counter++
		} else if _, found := replicateactualconcmap[volstr.(string)]; !found {
			replicatevalues = append(replicatevalues, actualconc.(float64))
			replicateactualconcmap[volstr.(string)] = replicatevalues
			replicatevalues = make([]float64, 0)
			//counter++
		}

		// get response value and check if it's a float64 type
		correctness, err := runwithresponses.GetResponseValue("Absorbance CorrectnessFactor " + strconv.Itoa(_input.Wavelength))

		if err != nil {
			fmt.Println(err.Error())
			_output.Errors = append(_output.Errors, err.Error())
		}

		correctnessfloat, floattrue := correctness.(float64)

		if !floattrue {
			fmt.Println(err.Error())
			execute.Errorf(_ctx, " Correctnessfloat not float but:"+fmt.Sprint(correctnessfloat))
		}

		// correctness factor map
		if _, found := replicatecorrectnessmap[volstr.(string)]; found /*&& rep == counter*/ {
			correctnessvalues = replicatecorrectnessmap[volstr.(string)]
			correctnessvalues = append(correctnessvalues, correctnessfloat)
			replicatecorrectnessmap[volstr.(string)] = correctnessvalues
			correctnessvalues = make([]float64, 0)
			//counter++
		} else if _, found := replicatecorrectnessmap[volstr.(string)]; !found {
			correctnessvalues = append(correctnessvalues, correctnessfloat)
			replicatecorrectnessmap[volstr.(string)] = correctnessvalues
			correctnessvalues = make([]float64, 0)
			//counter++
		}

	}

	// process into datasets
	for key, values := range replicateactualconcmap {

		var dataset Dataset
		// process replicates into mean and cv
		dataset.Name = key + "_AbsorbanceActualConc"
		dataset.Mean, _ = stats.Mean(values)
		dataset.StdDev, _ = stats.StdDevS(values)
		dataset.Values = values

		dataset.CV = dataset.StdDev / dataset.Mean * float64(100)
		_output.VolumeToActualConc[key] = dataset

	}

	// process into datasets
	for key, values := range replicatecorrectnessmap {

		var dataset Dataset
		// process replicates into mean and cv
		dataset.Name = key + "_CorrectnessFactor"
		dataset.Mean, _ = stats.Mean(values)
		dataset.StdDev, _ = stats.StdDevS(values)
		dataset.Values = values

		dataset.CV = dataset.StdDev / dataset.Mean * float64(100)
		_output.VolumeToCorrectnessFactor[key] = dataset

	}

	if _input.ManualComparison {

		// reset
		xvalues = make([]float64, 0)
		yvalues = make([]float64, 0)

		// add origin
		xvalues = append(xvalues, 0.0)
		yvalues = append(yvalues, 0.0)

		// 2. now plot correctnessfactor
		for i, runwithresponses := range _output.Runs {
			// values for r2 to reset each run

			// get response value and check if it's a float64 type
			expectedconc, err := runwithresponses.GetResponseValue("Absorbance ExpectedConc " + strconv.Itoa(_input.Wavelength))

			if err != nil {
				_output.Errors = append(_output.Errors, err.Error())
			}

			expectedconcfloat, floattrue := expectedconc.(float64)
			// if float64 is true
			if floattrue {
				xvalues = append(xvalues, expectedconcfloat)
			} else {
				execute.Errorf(_ctx, "Run"+fmt.Sprint(i, runwithresponses)+" ExpectedConc:"+fmt.Sprint(expectedconcfloat))
			}

			// get response value and check if it's a float64 type
			correctness, err := runwithresponses.GetResponseValue("Absorbance ManualCorrectnessFactor " + strconv.Itoa(_input.Wavelength))

			if err != nil {
				fmt.Println(err.Error())
				_output.Errors = append(_output.Errors, err.Error())
			}

			correctnessfloat, floattrue := correctness.(float64)

			if floattrue {
				yvalues = append(yvalues, correctnessfloat)
			} else {
				fmt.Println(err.Error())
				execute.Errorf(_ctx, "Manual Absorbance CorrectnessFactor:"+fmt.Sprint(correctnessfloat))
			}

		}

		_output.R2_CorrectnessFactor, _, _ = plot.Rsquared("Expected Conc", xvalues, "Manual Correctness Factor", yvalues)
		//run.AddResponseValue("R2", rsquared)

		correctnessgraph, err := plot.Plot(xvalues, [][]float64{yvalues})

		if err != nil {
			_output.Errors = append(_output.Errors, err.Error())
		}

		plot.Export(correctnessgraph, "10cm", "10cm", filenameandextension[0]+"_Manualcorrectnessfactor"+".png")

	}

}

// A block of tests to perform to validate that the sample was processed correctly
// Optionally, destructive tests can be performed to validate results on a
// dipstick basis
func _AddPlateReaderresults_2Validation(_ctx context.Context, _input *AddPlateReaderresults_2Input, _output *AddPlateReaderresults_2Output) {

	_output.CVpass = true

	if _output.R2 > _input.R2threshold {
		_output.R2Pass = true
	} else {
		_output.Errors = append(_output.Errors, fmt.Sprint("R2 threshold of ", _input.R2threshold, " not met; R2 value = ", _output.R2))
	}

	for key, dataset := range _output.VolumeToActualConc {

		if dataset.CV > _input.CVthreshold {
			_output.CVpass = false
			_output.Errors = append(_output.Errors, fmt.Sprint(key, " coefficient of variance above ", _input.CVthreshold, " percent threshold; CV value = ", dataset.CV))
		}
	}

}

type Dataset struct {
	Name   string
	Values []float64
	Mean   float64
	StdDev float64
	CV     float64
}

func _AddPlateReaderresults_2Run(_ctx context.Context, input *AddPlateReaderresults_2Input) *AddPlateReaderresults_2Output {
	output := &AddPlateReaderresults_2Output{}
	_AddPlateReaderresults_2Setup(_ctx, input)
	_AddPlateReaderresults_2Steps(_ctx, input, output)
	_AddPlateReaderresults_2Analysis(_ctx, input, output)
	_AddPlateReaderresults_2Validation(_ctx, input, output)
	return output
}

func AddPlateReaderresults_2RunSteps(_ctx context.Context, input *AddPlateReaderresults_2Input) *AddPlateReaderresults_2SOutput {
	soutput := &AddPlateReaderresults_2SOutput{}
	output := _AddPlateReaderresults_2Run(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func AddPlateReaderresults_2New() interface{} {
	return &AddPlateReaderresults_2Element{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &AddPlateReaderresults_2Input{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _AddPlateReaderresults_2Run(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &AddPlateReaderresults_2Input{},
			Out: &AddPlateReaderresults_2Output{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wunit.Make_units
)

type AddPlateReaderresults_2Element struct {
	inject.CheckedRunner
}

type AddPlateReaderresults_2Input struct {
	Blanks                []string
	CVthreshold           float64
	DesignFile            string
	DesignFiletype        string
	Diluent               *wtype.LHComponent
	Extinctioncoefficient float64
	FindOptWavelength     bool
	ManualComparison      bool
	MarsResultsFileXLSX   string
	Molecule              *wtype.LHComponent
	OutputFilename        string
	PlateType             *wtype.LHPlate
	R2threshold           float64
	ReadingTypeinMarsFile string
	Responsecolumnstofill []string
	Sheet                 int
	StockconcinMperL      wunit.Concentration
	Stockvol              wunit.Volume
	VolumeToManualwells   map[string][]string
	Wavelength            int
	WellForScanAnalysis   []string
}

type AddPlateReaderresults_2Output struct {
	BlankValues               []float64
	CV                        float64
	CVpass                    bool
	Errors                    []string
	Formula                   string
	MeasuredOptimalWavelength int
	R2                        float64
	R2Pass                    bool
	R2_CorrectnessFactor      float64
	ResponsetoManualValuesmap map[string][]float64
	Runs                      []doe.Run
	Variance                  float64
	VolumeToActualConc        map[string]Dataset
	VolumeToCorrectnessFactor map[string]Dataset
}

type AddPlateReaderresults_2SOutput struct {
	Data struct {
		BlankValues               []float64
		CV                        float64
		CVpass                    bool
		Errors                    []string
		Formula                   string
		MeasuredOptimalWavelength int
		R2                        float64
		R2Pass                    bool
		R2_CorrectnessFactor      float64
		ResponsetoManualValuesmap map[string][]float64
		Runs                      []doe.Run
		Variance                  float64
		VolumeToActualConc        map[string]Dataset
		VolumeToCorrectnessFactor map[string]Dataset
	}
	Outputs struct {
	}
}

func init() {
	if err := addComponent(component.Component{Name: "AddPlateReaderresults_2",
		Constructor: AddPlateReaderresults_2New,
		Desc: component.ComponentDesc{
			Desc: "",
			Path: "antha/component/an/Utility/AddPlateReaderResults_2.an",
			Params: []component.ParamDesc{
				{Name: "Blanks", Desc: "= []string{\"P24\"}\n", Kind: "Parameters"},
				{Name: "CVthreshold", Desc: "", Kind: "Parameters"},
				{Name: "DesignFile", Desc: "= \"250516CCFbubbles/240516DXCFFDoeoutputgilsonright_TEST.xlsx\"\n", Kind: "Parameters"},
				{Name: "DesignFiletype", Desc: "= \"JMP\"\n", Kind: "Parameters"},
				{Name: "Diluent", Desc: "", Kind: "Inputs"},
				{Name: "Extinctioncoefficient", Desc: "of target molecule at wavelength\n\n= 20330\n", Kind: "Parameters"},
				{Name: "FindOptWavelength", Desc: "= false\n", Kind: "Parameters"},
				{Name: "ManualComparison", Desc: "             = false\n", Kind: "Parameters"},
				{Name: "MarsResultsFileXLSX", Desc: "= \"250516CCFbubbles/260516ccfbubbles.xlsx\" //\"lhdoe110216_postspin_postshake.xlsx\"\n", Kind: "Parameters"},
				{Name: "Molecule", Desc: "", Kind: "Inputs"},
				{Name: "OutputFilename", Desc: "= \"250516CCFbubbles/2501516bubblesresults.xlsx\"\n", Kind: "Parameters"},
				{Name: "PlateType", Desc: "", Kind: "Inputs"},
				{Name: "R2threshold", Desc: "validation requirements\n", Kind: "Parameters"},
				{Name: "ReadingTypeinMarsFile", Desc: "= \"Abs Spectrum\"\n", Kind: "Parameters"},
				{Name: "Responsecolumnstofill", Desc: "= []string{\"AbsVLV\"}\n", Kind: "Parameters"},
				{Name: "Sheet", Desc: "= 0                                        //PRESHAKEPRESPIN\n", Kind: "Parameters"},
				{Name: "StockconcinMperL", Desc: "= 0.0002878191305957933\n", Kind: "Parameters"},
				{Name: "Stockvol", Desc: "= wunit.NewVolume(20, \"ul\")\n", Kind: "Parameters"},
				{Name: "VolumeToManualwells", Desc: "= map[string][]string{\t\"AbsVLV\": []string{\"\"},}\n", Kind: "Parameters"},
				{Name: "Wavelength", Desc: "= 472\n", Kind: "Parameters"},
				{Name: "WellForScanAnalysis", Desc: "= []string{\"J5\"}\n", Kind: "Parameters"},
				{Name: "BlankValues", Desc: "", Kind: "Data"},
				{Name: "CV", Desc: "", Kind: "Data"},
				{Name: "CVpass", Desc: "", Kind: "Data"},
				{Name: "Errors", Desc: "", Kind: "Data"},
				{Name: "Formula", Desc: "", Kind: "Data"},
				{Name: "MeasuredOptimalWavelength", Desc: "", Kind: "Data"},
				{Name: "R2", Desc: "", Kind: "Data"},
				{Name: "R2Pass", Desc: "", Kind: "Data"},
				{Name: "R2_CorrectnessFactor", Desc: "", Kind: "Data"},
				{Name: "ResponsetoManualValuesmap", Desc: "", Kind: "Data"},
				{Name: "Runs", Desc: "", Kind: "Data"},
				{Name: "Variance", Desc: "", Kind: "Data"},
				{Name: "VolumeToActualConc", Desc: "", Kind: "Data"},
				{Name: "VolumeToCorrectnessFactor", Desc: "", Kind: "Data"},
			},
		},
	}); err != nil {
		panic(err)
	}
}
