// Example protocol of performing an assembly simulation prior to performing
// physical construct assembly if the siulation passes

package lib

import (
	"fmt"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/Inventory"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/enzymes"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/sequences"
	"github.com/antha-lang/antha/antha/anthalib/mixer"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/component"
	"github.com/antha-lang/antha/execute"
	"github.com/antha-lang/antha/inject"
	"golang.org/x/net/context"
)

// Input parameters for this protocol (data)

// Physical Inputs to this protocol with types

// Physical outputs from this protocol with types

// Data which is returned from this protocol, and data types

// Input Requirement specification
func _TypeIISConstructAssembly_simRequirements() {

}

// Conditions to run on startup
func _TypeIISConstructAssembly_simSetup(_ctx context.Context, _input *TypeIISConstructAssembly_simInput) {
}

// The core process for this protocol, with the steps to be performed
// for every input
func _TypeIISConstructAssembly_simSteps(_ctx context.Context, _input *TypeIISConstructAssembly_simInput, _output *TypeIISConstructAssembly_simOutput) {
	// Check that assembly is feasible by simulating assembly of the sequences with the chosen enzyme
	partsinorder := make([]wtype.DNASequence, 0)

	for _, part := range _input.Partsinorder {
		partDNA := Inventory.Partslist()[part]
		partsinorder = append(partsinorder, partDNA)
	}

	vectordata := Inventory.Partslist()[_input.Vectordata]
	assembly := enzymes.Assemblyparameters{_input.Constructname, _input.RestrictionEnzyme.CName, vectordata, partsinorder}
	status, numberofassemblies, sitesfound, newDNASequence, _ := enzymes.Assemblysimulator(assembly)

	_output.NewDNASequence = newDNASequence
	_output.Sitesfound = sitesfound

	if status == "Yay! this should work" && numberofassemblies == 1 {

		_output.Simulationpass = true
	}
	// Monitor molar ratios of parts for possible troubleshooting / success correlation

	molesofeachdnaelement := make([]float64, 0)
	molarratios := make([]float64, 0)

	vector_mw := sequences.MassDNA(vectordata.Seq, false, true)
	vector_moles := sequences.Moles(_input.VectorConcentration, vector_mw, _input.VectorVol)
	molesofeachdnaelement = append(molesofeachdnaelement, vector_moles)

	molarratios = append(molarratios, (vector_moles / vector_moles))

	var part_mw float64
	var part_moles float64

	for i := 0; i < len(_input.Partsinorder); i++ {

		part_mw = sequences.MassDNA(partsinorder[i].Seq, false, true)
		part_moles = sequences.Moles(_input.PartConcs[i], part_mw, _input.PartVols[i])

		molesofeachdnaelement = append(molesofeachdnaelement, part_moles)
		molarratios = append(molarratios, (part_moles / vector_moles))
	}

	_output.Molesperpart = molesofeachdnaelement
	_output.MolarratiotoVector = molarratios

	// Print status
	_output.Status = fmt.Sprintln(
		"Simulationpass=", _output.Simulationpass,
		"Molesperpart", _output.Molesperpart,
		"MolarratiotoVector", _output.MolarratiotoVector,
		"NewDNASequence", _output.NewDNASequence,
		"Sitesfound", _output.Sitesfound,
	)

	if _output.Simulationpass == true {

		// Now Perform the physical assembly
		samples := make([]*wtype.LHComponent, 0)
		waterSample := mixer.SampleForTotalVolume(_input.Water, _input.ReactionVolume)
		samples = append(samples, waterSample)

		bufferSample := mixer.Sample(_input.Buffer, _input.BufferVol)
		samples = append(samples, bufferSample)

		atpSample := mixer.Sample(_input.Atp, _input.AtpVol)
		samples = append(samples, atpSample)

		//vectorSample := mixer.Sample(Vector, VectorVol)
		vectorSample := mixer.Sample(_input.Vector, _input.VectorVol)
		samples = append(samples, vectorSample)

		for k, part := range _input.Parts {
			fmt.Println("creating dna part num ", k, " comp ", part.CName, " renamed to ", _input.Partsinorder[k], " vol ", _input.PartVols[k])
			partSample := mixer.Sample(part, _input.PartVols[k])
			partSample.CName = _input.Partsinorder[k]
			samples = append(samples, partSample)
		}

		reSample := mixer.Sample(_input.RestrictionEnzyme, _input.ReVol)
		samples = append(samples, reSample)

		ligSample := mixer.Sample(_input.Ligase, _input.LigVol)
		samples = append(samples, ligSample)

		// incubate the reaction mixture
		out1 := execute.Incubate(_ctx, execute.MixInto(_ctx, _input.OutPlate, "", samples...), _input.ReactionTemp, _input.ReactionTime, false)
		// inactivate
		_output.Reaction = execute.Incubate(_ctx, out1, _input.InactivationTemp, _input.InactivationTime, false)
	}

}

// Run after controls and a steps block are completed to
// post process any data and provide downstream results
func _TypeIISConstructAssembly_simAnalysis(_ctx context.Context, _input *TypeIISConstructAssembly_simInput, _output *TypeIISConstructAssembly_simOutput) {
}

// A block of tests to perform to validate that the sample was processed correctly
// Optionally, destructive tests can be performed to validate results on a
// dipstick basis
func _TypeIISConstructAssembly_simValidation(_ctx context.Context, _input *TypeIISConstructAssembly_simInput, _output *TypeIISConstructAssembly_simOutput) {
}
func _TypeIISConstructAssembly_simRun(_ctx context.Context, input *TypeIISConstructAssembly_simInput) *TypeIISConstructAssembly_simOutput {
	output := &TypeIISConstructAssembly_simOutput{}
	_TypeIISConstructAssembly_simSetup(_ctx, input)
	_TypeIISConstructAssembly_simSteps(_ctx, input, output)
	_TypeIISConstructAssembly_simAnalysis(_ctx, input, output)
	_TypeIISConstructAssembly_simValidation(_ctx, input, output)
	return output
}

func TypeIISConstructAssembly_simRunSteps(_ctx context.Context, input *TypeIISConstructAssembly_simInput) *TypeIISConstructAssembly_simSOutput {
	soutput := &TypeIISConstructAssembly_simSOutput{}
	output := _TypeIISConstructAssembly_simRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func TypeIISConstructAssembly_simNew() interface{} {
	return &TypeIISConstructAssembly_simElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &TypeIISConstructAssembly_simInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _TypeIISConstructAssembly_simRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &TypeIISConstructAssembly_simInput{},
			Out: &TypeIISConstructAssembly_simOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wunit.Make_units
)

type TypeIISConstructAssembly_simElement struct {
	inject.CheckedRunner
}

type TypeIISConstructAssembly_simInput struct {
	Atp                 *wtype.LHComponent
	AtpVol              wunit.Volume
	Buffer              *wtype.LHComponent
	BufferVol           wunit.Volume
	Constructname       string
	InPlate             *wtype.LHPlate
	InactivationTemp    wunit.Temperature
	InactivationTime    wunit.Time
	LigVol              wunit.Volume
	Ligase              *wtype.LHComponent
	OutPlate            *wtype.LHPlate
	PartConcs           []wunit.Concentration
	PartVols            []wunit.Volume
	Parts               []*wtype.LHComponent
	Partsinorder        []string
	ReVol               wunit.Volume
	ReactionTemp        wunit.Temperature
	ReactionTime        wunit.Time
	ReactionVolume      wunit.Volume
	RestrictionEnzyme   *wtype.LHComponent
	Vector              *wtype.LHComponent
	VectorConcentration wunit.Concentration
	VectorVol           wunit.Volume
	Vectordata          string
	Water               *wtype.LHComponent
}

type TypeIISConstructAssembly_simOutput struct {
	MolarratiotoVector []float64
	Molesperpart       []float64
	NewDNASequence     wtype.DNASequence
	Reaction           *wtype.LHComponent
	Simulationpass     bool
	Sitesfound         []enzymes.Restrictionsites
	Status             string
}

type TypeIISConstructAssembly_simSOutput struct {
	Data struct {
		MolarratiotoVector []float64
		Molesperpart       []float64
		NewDNASequence     wtype.DNASequence
		Simulationpass     bool
		Sitesfound         []enzymes.Restrictionsites
		Status             string
	}
	Outputs struct {
		Reaction *wtype.LHComponent
	}
}

func init() {
	if err := addComponent(component.Component{Name: "TypeIISConstructAssembly_sim",
		Constructor: TypeIISConstructAssembly_simNew,
		Desc: component.ComponentDesc{
			Desc: "",
			Path: "src/github.com/antha-lang/antha/antha/component/an/Liquid_handling/TypeIIsAssembly/TypeIIsConstructAssembly_sim/TypeIIsConstructAssembly_sim.an",
			Params: []component.ParamDesc{
				{Name: "Atp", Desc: "", Kind: "Inputs"},
				{Name: "AtpVol", Desc: "", Kind: "Parameters"},
				{Name: "Buffer", Desc: "", Kind: "Inputs"},
				{Name: "BufferVol", Desc: "", Kind: "Parameters"},
				{Name: "Constructname", Desc: "", Kind: "Parameters"},
				{Name: "InPlate", Desc: "", Kind: "Inputs"},
				{Name: "InactivationTemp", Desc: "", Kind: "Parameters"},
				{Name: "InactivationTime", Desc: "", Kind: "Parameters"},
				{Name: "LigVol", Desc: "", Kind: "Parameters"},
				{Name: "Ligase", Desc: "", Kind: "Inputs"},
				{Name: "OutPlate", Desc: "", Kind: "Inputs"},
				{Name: "PartConcs", Desc: "", Kind: "Parameters"},
				{Name: "PartVols", Desc: "", Kind: "Parameters"},
				{Name: "Parts", Desc: "", Kind: "Inputs"},
				{Name: "Partsinorder", Desc: "", Kind: "Parameters"},
				{Name: "ReVol", Desc: "", Kind: "Parameters"},
				{Name: "ReactionTemp", Desc: "", Kind: "Parameters"},
				{Name: "ReactionTime", Desc: "", Kind: "Parameters"},
				{Name: "ReactionVolume", Desc: "", Kind: "Parameters"},
				{Name: "RestrictionEnzyme", Desc: "", Kind: "Inputs"},
				{Name: "Vector", Desc: "", Kind: "Inputs"},
				{Name: "VectorConcentration", Desc: "", Kind: "Parameters"},
				{Name: "VectorVol", Desc: "", Kind: "Parameters"},
				{Name: "Vectordata", Desc: "", Kind: "Parameters"},
				{Name: "Water", Desc: "", Kind: "Inputs"},
				{Name: "MolarratiotoVector", Desc: "", Kind: "Data"},
				{Name: "Molesperpart", Desc: "", Kind: "Data"},
				{Name: "NewDNASequence", Desc: "", Kind: "Data"},
				{Name: "Reaction", Desc: "", Kind: "Outputs"},
				{Name: "Simulationpass", Desc: "", Kind: "Data"},
				{Name: "Sitesfound", Desc: "", Kind: "Data"},
				{Name: "Status", Desc: "", Kind: "Data"},
			},
		},
	}); err != nil {
		panic(err)
	}
}

/*
type Mole struct {
	number float64
}*/
