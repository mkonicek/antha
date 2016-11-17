// Protocol for resuspending freeze dried DNA with a diluent
package lib

import

// we need to import the wtype package to use the LHComponent type
// the mixer package is required to use the Sample function
(
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/component"
	"github.com/antha-lang/antha/execute"
	"github.com/antha-lang/antha/inject"
	"golang.org/x/net/context"
)

// Input parameters for this protocol (data)

func _PairOligos_multipleRequirements() {
}

func _PairOligos_multipleSetup(_ctx context.Context, _input *PairOligos_multipleInput) {
}

func _PairOligos_multipleSteps(_ctx context.Context, _input *PairOligos_multipleInput, _output *PairOligos_multipleOutput) {

	//var nilconcmap map[string]wunit.Concentration
	var nilconc wunit.Concentration

	if _input.PartConcentrations == nil && _input.StockConcentration != nilconc {
		_input.PartConcentrations = make(map[string]wunit.Concentration)
		for key := range _input.DNAPartsMap {
			_input.PartConcentrations[key] = _input.StockConcentration
		}
	}

	// initialise output map
	_output.OligoPairs = make(map[string]*wtype.LHComponent)

	// get all well locations for plate
	var welllocations []string = _input.Plate.AllWellPositions(wtype.BYCOLUMN)

	// initialise a counter
	var counter int = 0
	var platenum int = 1
	// range through Oligo pairs map
	for fwd, rev := range _input.FwdOligotoRevOligoMap {

		// calculate volume to add for target conc
		fwdoligoVol, err := wunit.VolumeForTargetConcentration(_input.ConcentrationSetPoint, _input.PartConcentrations[fwd], _input.TotalVolume)

		if err != nil {
			execute.Errorf(_ctx, err.Error())
		}

		// calculate volume to add for target conc
		revoligoVol, err := wunit.VolumeForTargetConcentration(_input.ConcentrationSetPoint, _input.PartConcentrations[rev], _input.TotalVolume)

		if err != nil {
			execute.Errorf(_ctx, err.Error())
		}

		// next well
		well := welllocations[counter]

		// run PairOligos Antha element recursively
		result := PairOligosRunSteps(_ctx, &PairOligosInput{TotalVolume: _input.TotalVolume,
			IncubationTemp: _input.IncubationTemp,
			IncubationTime: _input.IncubationTime,
			FWDOligoVolume: fwdoligoVol,
			REVOligoVolume: revoligoVol,
			Well:           well,
			PlateNumber:    platenum,

			FwdOligo: _input.DNAPartsMap[fwd],
			RevOligo: _input.DNAPartsMap[rev],
			Diluent:  _input.Diluent,
			Plate:    _input.Plate},
		)

		// add to output map
		_output.OligoPairs[fwd] = result.Outputs.OligoPairs

		// increase counter to find next free well

		if counter+1 == len(welllocations) {
			counter = 0
			platenum++
		} else {
			counter++
		}
	}

}

func _PairOligos_multipleAnalysis(_ctx context.Context, _input *PairOligos_multipleInput, _output *PairOligos_multipleOutput) {
}

func _PairOligos_multipleValidation(_ctx context.Context, _input *PairOligos_multipleInput, _output *PairOligos_multipleOutput) {
}
func _PairOligos_multipleRun(_ctx context.Context, input *PairOligos_multipleInput) *PairOligos_multipleOutput {
	output := &PairOligos_multipleOutput{}
	_PairOligos_multipleSetup(_ctx, input)
	_PairOligos_multipleSteps(_ctx, input, output)
	_PairOligos_multipleAnalysis(_ctx, input, output)
	_PairOligos_multipleValidation(_ctx, input, output)
	return output
}

func PairOligos_multipleRunSteps(_ctx context.Context, input *PairOligos_multipleInput) *PairOligos_multipleSOutput {
	soutput := &PairOligos_multipleSOutput{}
	output := _PairOligos_multipleRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func PairOligos_multipleNew() interface{} {
	return &PairOligos_multipleElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &PairOligos_multipleInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _PairOligos_multipleRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &PairOligos_multipleInput{},
			Out: &PairOligos_multipleOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wunit.Make_units
)

type PairOligos_multipleElement struct {
	inject.CheckedRunner
}

type PairOligos_multipleInput struct {
	ConcentrationSetPoint wunit.Concentration
	DNAPartsMap           map[string]*wtype.LHComponent
	Diluent               *wtype.LHComponent
	FwdOligotoRevOligoMap map[string]string
	IncubationTemp        wunit.Temperature
	IncubationTime        wunit.Time
	PartConcentrations    map[string]wunit.Concentration
	Plate                 *wtype.LHPlate
	StockConcentration    wunit.Concentration
	TotalVolume           wunit.Volume
}

type PairOligos_multipleOutput struct {
	OligoPairs map[string]*wtype.LHComponent
}

type PairOligos_multipleSOutput struct {
	Data struct {
	}
	Outputs struct {
		OligoPairs map[string]*wtype.LHComponent
	}
}

func init() {
	if err := addComponent(component.Component{Name: "PairOligos_multiple",
		Constructor: PairOligos_multipleNew,
		Desc: component.ComponentDesc{
			Desc: "Protocol for resuspending freeze dried DNA with a diluent\n",
			Path: "src/github.com/antha-lang/antha/antha/component/an/Liquid_handling/ResuspendDNA/PairOligos_multiple.an",
			Params: []component.ParamDesc{
				{Name: "ConcentrationSetPoint", Desc: "", Kind: "Parameters"},
				{Name: "DNAPartsMap", Desc: "", Kind: "Inputs"},
				{Name: "Diluent", Desc: "", Kind: "Inputs"},
				{Name: "FwdOligotoRevOligoMap", Desc: "", Kind: "Parameters"},
				{Name: "IncubationTemp", Desc: "", Kind: "Parameters"},
				{Name: "IncubationTime", Desc: "", Kind: "Parameters"},
				{Name: "PartConcentrations", Desc: "", Kind: "Parameters"},
				{Name: "Plate", Desc: "", Kind: "Inputs"},
				{Name: "StockConcentration", Desc: "", Kind: "Parameters"},
				{Name: "TotalVolume", Desc: "", Kind: "Parameters"},
				{Name: "OligoPairs", Desc: "", Kind: "Outputs"},
			},
		},
	}); err != nil {
		panic(err)
	}
}
