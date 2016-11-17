// This protocol is based on scarfree design so please look at that first.
// The protocol is intended to design assembly parts using the first enzyme
// which is found to be feasible to use from a list of ApprovedEnzymes enzymes . If no enzyme
// from the list is feasible to use (i.e. due to the presence of existing restriction sites in a part)
// all typeIIs enzymes will be screened to find feasible backup options

package lib

import (
	"fmt"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/enzymes"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/enzymes/lookup"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/igem"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/text"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/component"
	"github.com/antha-lang/antha/execute"
	"github.com/antha-lang/antha/inject"
	"golang.org/x/net/context"
	"strconv"
	"strings"
)

// Input parameters for this protocol (data)

// Physical Inputs to this protocol with types

// Physical outputs from this protocol with types

// Data which is returned from this protocol, and data types

// parts to order
// desired sequence to end up with after assembly

// Input Requirement specification
func _Iterative_assembly_designRequirements() {
	// e.g. are MoClo types valid?
}

// Conditions to run on startup
func _Iterative_assembly_designSetup(_ctx context.Context, _input *Iterative_assembly_designInput) {}

// The core process for this protocol, with the steps to be performed
// for every input
func _Iterative_assembly_designSteps(_ctx context.Context, _input *Iterative_assembly_designInput, _output *Iterative_assembly_designOutput) {
	//var msg string
	// set warnings reported back to user to none initially

	warnings := make([]string, 0)
	sitefound := false
	Enzyme := "No enzymes which passed with these sequences"
	// make an empty array of DNA Sequences ready to fill
	partsinorder := make([]wtype.DNASequence, 0)

	_output.Status = "all parts available"
	for i, part := range _input.Seqsinorder {
		if strings.Contains(part, "BBa_") {
			part = igem.GetSequence(part)
		}
		partDNA := wtype.MakeLinearDNASequence("Part "+strconv.Itoa(i), part)

		partsinorder = append(partsinorder, partDNA)
	}
	// Find all possible typeIIs enzymes we could use for these sequences (i.e. non cutters of all parts)
	possibilities := lookup.FindEnzymeNamesofClass("TypeIIs")
	var backupoption string
	for _, possibility := range possibilities {
		// check number of sites per part !
		enz := lookup.EnzymeLookup(possibility)

		for _, part := range partsinorder {

			info := enzymes.Restrictionsitefinder(part, []wtype.RestrictionEnzyme{enz})
			if len(info) != 0 {
				if info[0].Sitefound == true {
					sitefound = true
					break
				}
			}
		}
		if sitefound == false {
			backupoption = possibility
			_output.BackupEnzymes = append(_output.BackupEnzymes, backupoption)
		}
	}

	sitefound = false
	for _, Enzyme := range _input.ApprovedEnzymes {

		// check number of sites per part !
		enz := lookup.EnzymeLookup(Enzyme)

		for _, part := range partsinorder {

			info := enzymes.Restrictionsitefinder(part, []wtype.RestrictionEnzyme{enz})
			if len(info) != 0 {
				if info[0].Sitefound == true {
					sitefound = true
					break
				}
			}
		}
		if sitefound == false {
			_output.EnzymeUsed = enz
		}
	}

	if sitefound != true {
		fmt.Println("enzyme used", _output.EnzymeUsed)
		Enzyme = _output.EnzymeUsed.Name

		// make vector into an antha type DNASequence
		vectordata := wtype.MakePlasmidDNASequence("Vector", _input.Vector)

		//lookup restriction enzyme
		restrictionenzyme, err := lookup.TypeIIsLookup(_output.EnzymeUsed.Name)
		if err != nil {
			text.Print("Error", err.Error())
		}

		//  Add overhangs for scarfree assembly based on part seqeunces only, i.e. no Assembly standard
		_output.PartswithOverhangs = enzymes.MakeScarfreeCustomTypeIIsassemblyParts(partsinorder, vectordata, restrictionenzyme)

		// Check that assembly is feasible with designed parts by simulating assembly of the sequences with the chosen enzyme
		assembly := enzymes.Assemblyparameters{_input.Constructname, restrictionenzyme.Name, vectordata, _output.PartswithOverhangs}
		status, numberofassemblies, _, newDNASequence, err := enzymes.Assemblysimulator(assembly)

		endreport := "Endreport only run in the event of assembly simulation failure"
		//sites := "Restriction mapper only run in the event of assembly simulation failure"
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

		// check number of sites per part !
		enz := lookup.EnzymeLookup(Enzyme)
		sites := make([]int, 0)
		multiple := make([]string, 0)
		for _, part := range _output.PartswithOverhangs {

			info := enzymes.Restrictionsitefinder(part, []wtype.RestrictionEnzyme{enz})

			sitepositions := enzymes.SitepositionString(info[0])

			sites = append(sites, info[0].Numberofsites)
			sitepositions = text.Print(part.Nm+" "+Enzyme+" positions:", sitepositions)
			multiple = append(multiple, sitepositions)
		}

		if len(warnings) == 0 {
			warnings = append(warnings, "none")
		}
		_output.Warnings = fmt.Errorf(strings.Join(warnings, ";"))

		partsummary := make([]string, 0)
		for _, part := range _output.PartswithOverhangs {
			partsummary = append(partsummary, text.Print(part.Nm, part.Seq))
		}

		partstoorder := text.Print("PartswithOverhangs: ", partsummary)

		_output.Status = fmt.Sprintln(
			text.Print("simulator status: ", status),
			text.Print("Endreport after digestion: ", endreport),
			text.Print("Sites per part for "+Enzyme, sites),
			text.Print("Positions: ", multiple),
			text.Print("Warnings:", _output.Warnings.Error()),
			text.Print("Simulationpass=", _output.Simulationpass),
			text.Print("NewDNASequence: ", _output.NewDNASequence),
			partstoorder)

	}
	// Print status
	if _output.Status != "all parts available" {
		_output.Status = fmt.Sprintln(_output.Status,
			text.Print("Backup Enzymes: ", _output.BackupEnzymes))
	} else if sitefound == true {
		_output.Status = fmt.Sprintln(text.Print("No Enzyme found to be compatible from approved list", _input.ApprovedEnzymes),
			text.Print("Backup Enzymes: ", _output.BackupEnzymes))

	} else {
		_output.Status = fmt.Sprintln(_output.Status,
			text.Print("Backup Enzymes: ", _output.BackupEnzymes))

	}

}

// Run after controls and a steps block are completed to
// post process any data and provide downstream results
func _Iterative_assembly_designAnalysis(_ctx context.Context, _input *Iterative_assembly_designInput, _output *Iterative_assembly_designOutput) {
}

// A block of tests to perform to validate that the sample was processed correctly
// Optionally, destructive tests can be performed to validate results on a
// dipstick basis
func _Iterative_assembly_designValidation(_ctx context.Context, _input *Iterative_assembly_designInput, _output *Iterative_assembly_designOutput) {
}
func _Iterative_assembly_designRun(_ctx context.Context, input *Iterative_assembly_designInput) *Iterative_assembly_designOutput {
	output := &Iterative_assembly_designOutput{}
	_Iterative_assembly_designSetup(_ctx, input)
	_Iterative_assembly_designSteps(_ctx, input, output)
	_Iterative_assembly_designAnalysis(_ctx, input, output)
	_Iterative_assembly_designValidation(_ctx, input, output)
	return output
}

func Iterative_assembly_designRunSteps(_ctx context.Context, input *Iterative_assembly_designInput) *Iterative_assembly_designSOutput {
	soutput := &Iterative_assembly_designSOutput{}
	output := _Iterative_assembly_designRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func Iterative_assembly_designNew() interface{} {
	return &Iterative_assembly_designElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &Iterative_assembly_designInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _Iterative_assembly_designRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &Iterative_assembly_designInput{},
			Out: &Iterative_assembly_designOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wunit.Make_units
)

type Iterative_assembly_designElement struct {
	inject.CheckedRunner
}

type Iterative_assembly_designInput struct {
	ApprovedEnzymes []string
	Constructname   string
	Seqsinorder     []string
	Vector          string
}

type Iterative_assembly_designOutput struct {
	BackupEnzymes      []string
	EnzymeUsed         wtype.RestrictionEnzyme
	NewDNASequence     wtype.DNASequence
	PartswithOverhangs []wtype.DNASequence
	Simulationpass     bool
	Status             string
	Warnings           error
}

type Iterative_assembly_designSOutput struct {
	Data struct {
		BackupEnzymes      []string
		EnzymeUsed         wtype.RestrictionEnzyme
		NewDNASequence     wtype.DNASequence
		PartswithOverhangs []wtype.DNASequence
		Simulationpass     bool
		Status             string
		Warnings           error
	}
	Outputs struct {
	}
}

func init() {
	if err := addComponent(component.Component{Name: "Iterative_assembly_design",
		Constructor: Iterative_assembly_designNew,
		Desc: component.ComponentDesc{
			Desc: "",
			Path: "src/github.com/antha-lang/antha/antha/component/an/Data/DNA/TypeIISAssembly_design/Iterative_assembly_design.an",
			Params: []component.ParamDesc{
				{Name: "ApprovedEnzymes", Desc: "", Kind: "Parameters"},
				{Name: "Constructname", Desc: "", Kind: "Parameters"},
				{Name: "Seqsinorder", Desc: "", Kind: "Parameters"},
				{Name: "Vector", Desc: "", Kind: "Parameters"},
				{Name: "BackupEnzymes", Desc: "", Kind: "Data"},
				{Name: "EnzymeUsed", Desc: "", Kind: "Data"},
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
