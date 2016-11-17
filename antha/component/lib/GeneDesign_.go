package lib

import (
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/enzymes"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/enzymes/lookup"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/export"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/sequences"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/sequences/entrez"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/component"
	"github.com/antha-lang/antha/execute"
	"github.com/antha-lang/antha/inject"
	"golang.org/x/net/context"
	"path/filepath"
)

// dna sequences as strings "ACTTGCGTC","GGTCCA"
// dna sequence as string
// name you want to give your construct
// typeIIs restriction enzyme name
// have the typeIIs assembly ends been added already? true/false
// name of synthesis provider e.g. GenScript
// Whether or not you want to export the sequences generated to a fasta file

// output parts with correct overhangs

func _GeneDesignRequirements() {
}

func _GeneDesignSetup(_ctx context.Context, _input *GeneDesignInput) {
}

func _GeneDesignSteps(_ctx context.Context, _input *GeneDesignInput, _output *GeneDesignOutput) {
	PartDNA := make([]wtype.DNASequence, 0)

	// Retrieve part seqs from entrez
	for _, part := range _input.Parts {
		//desiredfilename := filepath.Join(anthapath.Path(), part+".gb")
		desiredfilename := filepath.Join(_input.ConstructName, part+".gb")
		DNA, _, _ := entrez.RetrieveSequence(part, "nucleotide", desiredfilename)
		PartDNA = append(PartDNA, DNA)
	}

	// Look up the restriction enzyme
	EnzymeInf, _ := lookup.TypeIIsLookup(_input.RE)

	// look up vector sequence

	//desiredvectorfilename := filepath.Join(anthapath.Path(), Vector+".gb")
	desiredvectorfilename := filepath.Join(_input.ConstructName, _input.Vector+".gb")

	VectorSeq, _, _ := entrez.RetrieveVector(_input.Vector, desiredvectorfilename)

	// Add overhangs
	if _input.EndsAlreadyAdded {
		_output.PartsWithOverhangs = PartDNA
	} else {
		// fmt.Println("Parts + vector:",PartDNA,VectorSeq)
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

	//fmt.Println("Parts Source: ",PartsWithOverhangs)
}

func _GeneDesignAnalysis(_ctx context.Context, _input *GeneDesignInput, _output *GeneDesignOutput) {

}

func _GeneDesignValidation(_ctx context.Context, _input *GeneDesignInput, _output *GeneDesignOutput) {

}
func _GeneDesignRun(_ctx context.Context, input *GeneDesignInput) *GeneDesignOutput {
	output := &GeneDesignOutput{}
	_GeneDesignSetup(_ctx, input)
	_GeneDesignSteps(_ctx, input, output)
	_GeneDesignAnalysis(_ctx, input, output)
	_GeneDesignValidation(_ctx, input, output)
	return output
}

func GeneDesignRunSteps(_ctx context.Context, input *GeneDesignInput) *GeneDesignSOutput {
	soutput := &GeneDesignSOutput{}
	output := _GeneDesignRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func GeneDesignNew() interface{} {
	return &GeneDesignElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &GeneDesignInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _GeneDesignRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &GeneDesignInput{},
			Out: &GeneDesignOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wunit.Make_units
)

type GeneDesignElement struct {
	inject.CheckedRunner
}

type GeneDesignInput struct {
	ConstructName     string
	EndsAlreadyAdded  bool
	ExporttoFastaFile bool
	Parts             []string
	RE                string
	SynthesisProvider string
	Vector            string
}

type GeneDesignOutput struct {
	PartsWithOverhangs []wtype.DNASequence
	SimulationStatus   string
	Validated          bool
	ValidationStatus   string
}

type GeneDesignSOutput struct {
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
	if err := addComponent(component.Component{Name: "GeneDesign",
		Constructor: GeneDesignNew,
		Desc: component.ComponentDesc{
			Desc: "",
			Path: "src/github.com/antha-lang/antha/antha/component/an/Data/DNA/GeneDesign/GeneDesign.an",
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
