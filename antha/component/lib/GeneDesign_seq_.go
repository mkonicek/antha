package lib

import (
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/enzymes"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/enzymes/lookup"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/export"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/sequences"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	//"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/sequences/entrez"
	"fmt"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/component"
	"github.com/antha-lang/antha/execute"
	"github.com/antha-lang/antha/inject"
	"golang.org/x/net/context"
	"strconv"
)

// dna sequences as strings "ACTTGCGTC","GGTCCA"
// dna sequence as string
// name you want to give your construct
// typeIIs restriction enzyme name
// have the typeIIs assembly ends been added already? true/false
// name of synthesis provider e.g. GenScript
// Whether or not you want to export the sequences generated to a fasta file

// output parts with correct overhangs

func _GeneDesign_seqRequirements() {
}

func _GeneDesign_seqSetup(_ctx context.Context, _input *GeneDesign_seqInput) {
}

func _GeneDesign_seqSteps(_ctx context.Context, _input *GeneDesign_seqInput, _output *GeneDesign_seqOutput) {
	PartDNA := make([]wtype.DNASequence, 0)

	// make DNASequence type from sequence
	for i, part := range _input.Parts {
		DNA := wtype.MakeLinearDNASequence("part"+strconv.Itoa(i), part)
		PartDNA = append(PartDNA, DNA)
	}

	// make vector sequence
	VectorSeq := wtype.MakePlasmidDNASequence("Vector", _input.Vector)

	// Look up the restriction enzyme
	EnzymeInf, _ := lookup.TypeIIsLookup(_input.RE)

	// Add overhangs
	if _input.EndsAlreadyAdded {
		_output.PartsWithOverhangs = PartDNA
	} else {
		_output.PartsWithOverhangs = enzymes.MakeScarfreeCustomTypeIIsassemblyParts(PartDNA, VectorSeq, EnzymeInf)
	}

	// validation
	assembly := enzymes.Assemblyparameters{_input.ConstructName, _input.RE, VectorSeq, _output.PartsWithOverhangs}
	_output.SimulationStatus, _, _, _, _ = enzymes.Assemblysimulator(assembly)

	// check if sequence meets requirements for synthesis
	_output.ValidationStatus, _output.Validated = sequences.ValidateSynthesis(_output.PartsWithOverhangs, _input.Vector, _input.SynthesisProvider)

	// export sequence to fasta
	if _input.ExporttoFastaFile {
		export.Makefastaserial2(export.LOCAL, _input.ConstructName, _output.PartsWithOverhangs)
	}

	fmt.Println("Parts Source: ", _output.PartsWithOverhangs)
}

func _GeneDesign_seqAnalysis(_ctx context.Context, _input *GeneDesign_seqInput, _output *GeneDesign_seqOutput) {

}

func _GeneDesign_seqValidation(_ctx context.Context, _input *GeneDesign_seqInput, _output *GeneDesign_seqOutput) {

}
func _GeneDesign_seqRun(_ctx context.Context, input *GeneDesign_seqInput) *GeneDesign_seqOutput {
	output := &GeneDesign_seqOutput{}
	_GeneDesign_seqSetup(_ctx, input)
	_GeneDesign_seqSteps(_ctx, input, output)
	_GeneDesign_seqAnalysis(_ctx, input, output)
	_GeneDesign_seqValidation(_ctx, input, output)
	return output
}

func GeneDesign_seqRunSteps(_ctx context.Context, input *GeneDesign_seqInput) *GeneDesign_seqSOutput {
	soutput := &GeneDesign_seqSOutput{}
	output := _GeneDesign_seqRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func GeneDesign_seqNew() interface{} {
	return &GeneDesign_seqElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &GeneDesign_seqInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _GeneDesign_seqRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &GeneDesign_seqInput{},
			Out: &GeneDesign_seqOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wunit.Make_units
)

type GeneDesign_seqElement struct {
	inject.CheckedRunner
}

type GeneDesign_seqInput struct {
	ConstructName     string
	EndsAlreadyAdded  bool
	ExporttoFastaFile bool
	Parts             []string
	RE                string
	SynthesisProvider string
	Vector            string
}

type GeneDesign_seqOutput struct {
	PartsWithOverhangs []wtype.DNASequence
	SimulationStatus   string
	Validated          bool
	ValidationStatus   string
}

type GeneDesign_seqSOutput struct {
	Data struct {
		PartsWithOverhangs []wtype.DNASequence
		SimulationStatus   string
		Validated          bool
		ValidationStatus   string
	}
	Outputs struct {
	}
}

func init() {
	if err := addComponent(component.Component{Name: "GeneDesign_seq",
		Constructor: GeneDesign_seqNew,
		Desc: component.ComponentDesc{
			Desc: "",
			Path: "src/github.com/antha-lang/antha/antha/component/an/Data/DNA/GeneDesign/GeneDesign_seq.an",
			Params: []component.ParamDesc{
				{Name: "ConstructName", Desc: "name you want to give your construct\n", Kind: "Parameters"},
				{Name: "EndsAlreadyAdded", Desc: "have the typeIIs assembly ends been added already? true/false\n", Kind: "Parameters"},
				{Name: "ExporttoFastaFile", Desc: "Whether or not you want to export the sequences generated to a fasta file\n", Kind: "Parameters"},
				{Name: "Parts", Desc: "dna sequences as strings \"ACTTGCGTC\",\"GGTCCA\"\n", Kind: "Parameters"},
				{Name: "RE", Desc: "typeIIs restriction enzyme name\n", Kind: "Parameters"},
				{Name: "SynthesisProvider", Desc: "name of synthesis provider e.g. GenScript\n", Kind: "Parameters"},
				{Name: "Vector", Desc: "dna sequence as string\n", Kind: "Parameters"},
				{Name: "PartsWithOverhangs", Desc: "output parts with correct overhangs\n", Kind: "Data"},
				{Name: "SimulationStatus", Desc: "", Kind: "Data"},
				{Name: "Validated", Desc: "", Kind: "Data"},
				{Name: "ValidationStatus", Desc: "", Kind: "Data"},
			},
		},
	}); err != nil {
		panic(err)
	}
}
