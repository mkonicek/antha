// This protocol is intended to design assembly parts using the MoClo assembly standard.
// Overhangs for a part are chosen according to the designated class of each part (e.g. promoter).
// The MoClo standard is hierarchical so the enzyme is chosen based on the level of assembly.
// i.e. first level 0 parts are made which may comprise of a promoter, 5prime upstream part, coding sequene, and terminator.
// Level 0 parts can then be assembled together by using level 1 enzymes and overhangs.
// currently this protocol only supports level 0 steps.
// see http://journals.plos.org/plosone/article?id=10.1371/journal.pone.0016765

package lib

import (
	"fmt"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/Inventory"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/enzymes"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/igem"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/sequences"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/text"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/component"
	"github.com/antha-lang/antha/execute"
	"github.com/antha-lang/antha/inject"
	"golang.org/x/net/context"
	"strings"
)

// Input parameters for this protocol (data)

//MoClo
// of assembly standard

// labels e.g. pro = promoter

// Physical Inputs to this protocol with types

// Physical outputs from this protocol with types

// Data which is returned from this protocol, and data types

// parts to order
// desired sequence to end up with after assembly

// Input Requirement specification
func _MoClo_designRequirements() {
	// e.g. are MoClo types valid?
}

// Conditions to run on startup
func _MoClo_designSetup(_ctx context.Context, _input *MoClo_designInput) {}

// The core process for this protocol, with the steps to be performed
// for every input
func _MoClo_designSteps(_ctx context.Context, _input *MoClo_designInput, _output *MoClo_designOutput) {
	//var msg string
	// set warnings reported back to user to none initially
	warnings := make([]string, 1)
	warnings[0] = "none"
	found := false
	var err error

	/* find sequence data from keyword; looking it up by a given name in an inventory
	   or by biobrick ID from iGem parts registry */
	partsinorder := make([]wtype.DNASequence, 0)
	var partDNA wtype.DNASequence

	_output.Status = "all parts available"
	for _, part := range _input.Partsinorder {

		if strings.Contains(part, "BBa_") == true {

			partDNA.Nm = part
			partproperties := igem.LookUp([]string{part})
			partDNA.Seq = partproperties.Sequence(part)
			//partDNA.Seq = igem.GetSequence(part)

			/* We can add logic to check the status of parts too and return a warning if the part
			   is not characterised */

			if strings.Contains(partproperties.Results(part), "Works") != true {

				warnings = make([]string, 0)
				//		warning := fmt.Sprintln("iGem part", part, "results =",  igem.GetResults(part), "rating",igem.GetRating(part), "part type",igem.GetType(part), "part decription =", igem.GetDescription(part), "Categories",igem.GetCategories(part))
				warning := fmt.Sprintln("iGem part", part, "results =", partproperties.Results(part), "rating", partproperties.Rating(part), "part type", partproperties.Type(part), "part decription =", partproperties.Description(part), "Categories", partproperties.Categories(part))

				warnings = append(warnings, warning)

			}
		} else {
			// look up part in inventory
			partDNA, found = Inventory.Partslist[part]

			if !found {
				//Status = text.Print("part: " + partDNA.Nm, partDNA.Seq + ": not found in Inventory so element aborted!")

				// assume dna sequence and test
				partDNA = wtype.MakeLinearDNASequence("tempPart", part)

				// test for illegal nucleotides
				pass, illegals, _ := sequences.Illegalnucleotides(partDNA)

				if !pass {
					var newstatus = make([]string, 0)
					for _, illegal := range illegals {

						newstatus = append(newstatus, "part: "+partDNA.Nm+" "+partDNA.Seq+": contains illegalnucleotides:"+illegal.ToString())
					}

					execute.Errorf(_ctx, strings.Join(newstatus, ""))
				} else if _input.BlastSeqswithNoName {
					// run a blast search on the sequence to get the name
					blastsearch := BlastSearch_wtypeRunSteps(_ctx, &BlastSearch_wtypeInput{DNA: partDNA})
					partDNA.Nm = blastsearch.Data.AnthaSeq.Nm
				}

			}
		}
		partsinorder = append(partsinorder, partDNA)
	}
	// lookup vector sequence
	vectordata := Inventory.Partslist[_input.Vector]

	//lookup restriction enzyme
	restrictionenzyme := enzymes.Enzymelookup[_input.AssemblyStandard][_input.Level]

	// (1) Add standard overhangs using chosen assembly standard
	_output.PartswithOverhangs, err = enzymes.MakeStandardTypeIIsassemblyParts(partsinorder, _input.AssemblyStandard, _input.Level, _input.PartMoClotypesinorder)

	if err != nil {
		warnings = append(warnings, text.Print("Error", err.Error()))
		execute.Errorf(_ctx, err.Error())
	}

	// OR (2) Add overhangs for scarfree assembly based on part seqeunces only, i.e. no Assembly standard
	//PartswithOverhangs = enzymes.MakeScarfreeCustomTypeIIsassemblyParts(partsinorder, vectordata, restrictionenzyme)

	// Check that assembly is feasible with designed parts by simulating assembly of the sequences with the chosen enzyme
	assembly := enzymes.Assemblyparameters{_input.Constructname, restrictionenzyme.Name, vectordata, _output.PartswithOverhangs}
	status, numberofassemblies, _, newDNASequence, err := enzymes.Assemblysimulator(assembly)

	endreport := "Only run in the event of assembly failure"
	_output.NewDNASequence = newDNASequence
	if err == nil && numberofassemblies == 1 {

		_output.Simulationpass = true
	} else {
		warnings = append(warnings, status)
		// perform mock digest to test fragement overhangs (fragments are hidden by using _, )
		_, stickyends5, stickyends3 := enzymes.TypeIIsdigest(vectordata, restrictionenzyme)

		allends := make([]string, 0)
		ends := ""

		ends = text.Print(vectordata.Nm+" 5 Prime end: ", stickyends5)
		allends = append(allends, ends)
		ends = text.Print(vectordata.Nm+" 3 Prime end: ", stickyends3)
		allends = append(allends, ends)

		for _, part := range _output.PartswithOverhangs {
			_, stickyends5, stickyends3 := enzymes.TypeIIsdigest(part, restrictionenzyme)
			ends = text.Print(part.Nm+" 5 Prime end: ", stickyends5)
			allends = append(allends, ends)
			ends = text.Print(part.Nm+" 3 Prime end: ", stickyends3)
			allends = append(allends, ends)
		}
		endreport = strings.Join(allends, " ")
	}

	_output.Warnings = strings.Join(warnings, ";")

	partsummary := make([]string, 0)
	for _, part := range _output.PartswithOverhangs {
		partsummary = append(partsummary, text.Print(part.Nm, part.Seq))
	}

	partstoorder := text.Print("PartswithOverhangs: ", partsummary)

	// Print status
	if _output.Status != "all parts available" {
		_output.Status = fmt.Sprintln(_output.Status)
	} else {
		_output.Status = fmt.Sprintln(
			text.Print("simulator status: ", status),
			text.Print("Endreport after digestion: ", endreport),
			text.Print("Warnings:", _output.Warnings),
			text.Print("Simulationpass=", _output.Simulationpass),
			text.Print("NewDNASequence: ", _output.NewDNASequence),
			partstoorder,
		)
	}

}

// Run after controls and a steps block are completed to
// post process any data and provide downstream results
func _MoClo_designAnalysis(_ctx context.Context, _input *MoClo_designInput, _output *MoClo_designOutput) {
}

// A block of tests to perform to validate that the sample was processed correctly
// Optionally, destructive tests can be performed to validate results on a
// dipstick basis
func _MoClo_designValidation(_ctx context.Context, _input *MoClo_designInput, _output *MoClo_designOutput) {
}
func _MoClo_designRun(_ctx context.Context, input *MoClo_designInput) *MoClo_designOutput {
	output := &MoClo_designOutput{}
	_MoClo_designSetup(_ctx, input)
	_MoClo_designSteps(_ctx, input, output)
	_MoClo_designAnalysis(_ctx, input, output)
	_MoClo_designValidation(_ctx, input, output)
	return output
}

func MoClo_designRunSteps(_ctx context.Context, input *MoClo_designInput) *MoClo_designSOutput {
	soutput := &MoClo_designSOutput{}
	output := _MoClo_designRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func MoClo_designNew() interface{} {
	return &MoClo_designElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &MoClo_designInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _MoClo_designRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &MoClo_designInput{},
			Out: &MoClo_designOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wunit.Make_units
)

type MoClo_designElement struct {
	inject.CheckedRunner
}

type MoClo_designInput struct {
	AssemblyStandard      string
	BlastSeqswithNoName   bool
	Constructname         string
	Level                 string
	PartMoClotypesinorder []string
	Partsinorder          []string
	Vector                string
}

type MoClo_designOutput struct {
	NewDNASequence     wtype.DNASequence
	PartswithOverhangs []wtype.DNASequence
	Simulationpass     bool
	Status             string
	Warnings           string
}

type MoClo_designSOutput struct {
	Data struct {
		NewDNASequence     wtype.DNASequence
		PartswithOverhangs []wtype.DNASequence
		Simulationpass     bool
		Status             string
		Warnings           string
	}
	Outputs struct {
	}
}

func init() {
	if err := addComponent(component.Component{Name: "MoClo_design",
		Constructor: MoClo_designNew,
		Desc: component.ComponentDesc{
			Desc: "",
			Path: "antha/component/an/Data/DNA/TypeIISAssembly_design/MoClo_design.an",
			Params: []component.ParamDesc{
				{Name: "AssemblyStandard", Desc: "MoClo\n", Kind: "Parameters"},
				{Name: "BlastSeqswithNoName", Desc: "", Kind: "Parameters"},
				{Name: "Constructname", Desc: "", Kind: "Parameters"},
				{Name: "Level", Desc: "of assembly standard\n", Kind: "Parameters"},
				{Name: "PartMoClotypesinorder", Desc: "labels e.g. pro = promoter\n", Kind: "Parameters"},
				{Name: "Partsinorder", Desc: "", Kind: "Parameters"},
				{Name: "Vector", Desc: "", Kind: "Parameters"},
				{Name: "NewDNASequence", Desc: "desired sequence to end up with after assembly\n", Kind: "Data"},
				{Name: "PartswithOverhangs", Desc: "parts to order\n", Kind: "Data"},
				{Name: "Simulationpass", Desc: "", Kind: "Data"},
				{Name: "Status", Desc: "", Kind: "Data"},
				{Name: "Warnings", Desc: "", Kind: "Data"},
			},
		},
	}); err != nil {
		panic(err)
	}
}
