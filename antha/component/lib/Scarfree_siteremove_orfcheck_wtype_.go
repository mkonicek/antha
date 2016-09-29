// This protocol is intended to design assembly parts using a specified enzyme.
// overhangs are added to complement the adjacent parts and leave no scar.
// parts can be entered as genbank (.gb) files, sequences or biobrick IDs
// If assembly simulation fails after overhangs are added. In order to help the user
// diagnose the reason, a report of the part overhangs
// is returned to the user along with a list of cut sites in each part.

package lib

import (
	"fmt"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/enzymes"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/enzymes/lookup"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/export"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/sequences"
	features "github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/sequences/features"
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

//"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/AnthaPath"

// Input parameters for this protocol (data)

// enter each as amino acid sequence

// Physical Inputs to this protocol with types

// Physical outputs from this protocol with types

// Data which is returned from this protocol, and data types

// parts to order

// desired sequence to end up with after assembly

// Input Requirement specification
func _Scarfree_siteremove_orfcheck_wtypeRequirements() {
	// e.g. are MoClo types valid?
}

// Conditions to run on startup
func _Scarfree_siteremove_orfcheck_wtypeSetup(_ctx context.Context, _input *Scarfree_siteremove_orfcheck_wtypeInput) {
}

// The core process for this protocol, with the steps to be performed
// for every input
func _Scarfree_siteremove_orfcheck_wtypeSteps(_ctx context.Context, _input *Scarfree_siteremove_orfcheck_wtypeInput, _output *Scarfree_siteremove_orfcheck_wtypeOutput) {

	// set warnings reported back to user to none initially
	warnings := make([]string, 0)

	// declare some temporary variables to be used later
	var warning string
	var err error

	partsinorder := _input.Seqsinorder
	vectordata := _input.Vector

	// export original parts list
	_output.OriginalParts = partsinorder

	warnings = append(warnings, fmt.Sprintln(_input.Seqsinorder))

	// check parts for restriction sites first and remove if the user has chosen to
	enz := lookup.EnzymeLookup(_input.Enzymename)

	// get properties of other enzyme sites to remove
	removetheseenzymes := make([]wtype.RestrictionEnzyme, 0)
	removetheseenzymes = append(removetheseenzymes, enz)

	for _, enzyme := range _input.OtherEnzymeSitesToRemove {
		removetheseenzymes = append(removetheseenzymes, lookup.EnzymeLookup(enzyme))
	}

	warning = text.Print("RemoveproblemRestrictionSites =", _input.RemoveproblemRestrictionSites)
	warnings = append(warnings, warning)
	if _input.RemoveproblemRestrictionSites {
		newparts := make([]wtype.DNASequence, 0)
		warning = "Starting process or removing restrictionsite"
		warnings = append(warnings, warning)

		for _, part := range partsinorder {

			info := enzymes.Restrictionsitefinder(part, removetheseenzymes)

			for _, anysites := range info {
				if anysites.Sitefound {
					warning = "problem " + anysites.Enzyme.Name + " site found in " + part.Nm
					warnings = append(warnings, warning)
					orf, orftrue := sequences.FindBiggestORF(part.Seq)
					warning = fmt.Sprintln(anysites.Enzyme.Name+" site found in orf ", part.Nm, " ", orftrue, " site positions ", anysites.Positions("ALL"), "orf between", orf.StartPosition, " and ", orf.EndPosition /*orf.DNASeq[orf.StartPosition:orf.EndPosition]*/)
					warnings = append(warnings, warning)
					if orftrue /* && len(orf.ProtSeq) > 20 */ {
						allsitestoavoid := make([]string, 0)
						allsitestoavoid = append(allsitestoavoid, anysites.Recognitionsequence, sequences.RevComp(anysites.Recognitionsequence))
						orfcoordinates := sequences.MakeStartendPair(orf.StartPosition, orf.EndPosition)
						for _, position := range anysites.Positions("ALL") {
							if orf.StartPosition < position && position < orf.EndPosition {
								originalcodon := ""
								codonoption := ""
								part, originalcodon, codonoption, err = sequences.ReplaceCodoninORF(part, orfcoordinates, position, allsitestoavoid)
								warning = fmt.Sprintln("sites to avoid: ", allsitestoavoid[0], allsitestoavoid[1])
								warnings = append(warnings, warning)
								warnings = append(warnings, "Paaaaerrttseq: "+part.Seq+"position: "+strconv.Itoa(position)+" original: "+originalcodon+" replacementcodon: "+codonoption)
								if err != nil {
									warning := text.Print("removal of "+anysites.Enzyme.Name+" site from orf "+orf.DNASeq, " failed! improve your algorithm! "+err.Error())
									warnings = append(warnings, warning)
								}
							} else {
								allsitestoavoid := make([]string, 0)
								part, err = sequences.RemoveSite(part, anysites.Enzyme, allsitestoavoid)
								if err != nil {

									warning = text.Print(anysites.Enzyme.Name+" position found to be outside of orf: "+orf.DNASeq, " failed! improve your algorithm! "+err.Error())
									warnings = append(warnings, warning)
								}
							}
						}
					} else {
						allsitestoavoid := make([]string, 0)
						temppart, err := sequences.RemoveSite(part, anysites.Enzyme, allsitestoavoid)
						//		fmt.Println("part= ", part)
						//		fmt.Println("temppart= ", temppart)
						if err != nil {
							warning := text.Print("removal of site failed! improve your algorithm!", err.Error())
							warnings = append(warnings, warning)

						}
						warning = fmt.Sprintln("modified "+temppart.Nm+"new seq: ", temppart.Seq, "original seq: ", part.Seq)
						warnings = append(warnings, warning)
						part = temppart

						//	}
					}
				}

			}
			newparts = append(newparts, part)

			partsinorder = newparts
		}
	}

	// export the parts list with sites removed
	_output.PartsWithSitesRemoved = partsinorder

	//lookup restriction enzyme
	restrictionenzyme, err := lookup.TypeIIsLookup(_input.Enzymename)
	if err != nil {
		warnings = append(warnings, text.Print("Error", err.Error()))
	}

	//  Add overhangs for scarfree assembly based on part seqeunces only, i.e. no Assembly standard
	fmt.Println("warnings:", warnings)

	if _input.EndsAlreadyadded {
		_output.PartswithOverhangs = partsinorder
	} else {
		_output.PartswithOverhangs = enzymes.MakeScarfreeCustomTypeIIsassemblyParts(partsinorder, vectordata, restrictionenzyme)
	}

	// Check that assembly is feasible with designed parts by simulating assembly of the sequences with the chosen enzyme
	assembly := enzymes.Assemblyparameters{_input.Constructname, restrictionenzyme.Name, vectordata, _output.PartswithOverhangs}
	status, numberofassemblies, _, newDNASequence, simerr := enzymes.Assemblysimulator(assembly)

	if simerr != nil {
		warnings = append(warnings, text.Print("Error", simerr.Error()))
	}

	_output.Plasmid, _output.ORIpresent, _output.SelectionMarkerPresent, err = features.ValidPlasmid(newDNASequence)

	if err != nil {
		warnings = append(warnings, text.Print("Error", err.Error()))
	}

	endreport := "Endreport only run in the event of assembly simulation failure"
	//sites := "Restriction mapper only run in the event of assembly simulation failure"
	newDNASequence.Nm = _input.Constructname
	_output.NewDNASequence = newDNASequence
	if simerr == nil && numberofassemblies == 1 {

		_output.Simulationpass = true
	} // else {

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
	if !_output.Simulationpass {
		warnings = append(warnings, endreport)
	}
	//	}

	// check number of sites per part !

	sites := make([]int, 0)
	multiple := make([]string, 0)
	for _, part := range _output.PartswithOverhangs {

		info := enzymes.Restrictionsitefinder(part, removetheseenzymes)

		for i := range info {
			sitepositions := enzymes.SitepositionString(info[i])

			sites = append(sites, info[i].Numberofsites)
			sitepositions = text.Print(part.Nm+" "+info[i].Enzyme.Name+" positions:", sitepositions)
			multiple = append(multiple, sitepositions)
		}
	}

	for _, orf := range _input.ORFstoConfirm {
		if sequences.LookforSpecificORF(_output.NewDNASequence.Seq, orf) == false {
			warning = text.Print("orf not present: ", orf)
			warnings = append(warnings, warning)
			_output.ORFmissing = true
		}
	}

	if len(warnings) == 0 {
		warnings = append(warnings, "none")
	}
	_output.Warnings = fmt.Errorf(strings.Join(warnings, ";"))
	_output.Endreport = endreport
	_output.PositionReport = multiple

	partsummary := make([]string, 0)
	for _, part := range _output.PartswithOverhangs {
		partsummary = append(partsummary, text.Print(part.Nm, part.Seq))
	}
	partsummary = append(partsummary, text.Print("Vector:"+vectordata.Nm, vectordata.Seq))
	partstoorder := text.Print("PartswithOverhangs: ", partsummary)

	// Print status
	if _output.Status != "all parts available" {
		_output.Status = fmt.Sprintln(_output.Status)
	} else {
		_output.Status = fmt.Sprintln(
			text.Print("simulator status: ", status),
			text.Print("Endreport after digestion: ", endreport),
			text.Print("Sites per part for "+_input.Enzymename, sites),
			text.Print("Positions: ", multiple),
			text.Print("Warnings:", _output.Warnings.Error()),
			text.Print("Simulationpass=", _output.Simulationpass),
			text.Print("NewDNASequence: ", _output.NewDNASequence),
			text.Print("Any Orfs to confirm missing from new DNA sequence:", _output.ORFmissing),
			partstoorder,
		)
		// export data to file
		//anthapath.ExporttoFile("Report"+"_"+Constructname+".txt",[]byte(Status))
		//anthapath.ExportTextFile("Report"+"_"+Constructname+".txt",Status)
		fmt.Println(_output.Status)
	}

	// export sequence to fasta
	if _input.ExporttoFastaFile && _output.Simulationpass {
		exportedsequences := make([]wtype.DNASequence, 0)
		// add dna sequence produced
		exportedsequences = append(exportedsequences, _output.NewDNASequence)

		// export to file
		export.Makefastaserial2(export.LOCAL, _input.Constructname+"_AssemblyProduct", exportedsequences)

		// reset
		exportedsequences = make([]wtype.DNASequence, 0)
		// add all parts with overhangs
		for _, part := range _output.PartswithOverhangs {
			exportedsequences = append(exportedsequences, part)
		}
		export.Makefastaserial2(export.LOCAL, _input.Constructname+"_Parts", exportedsequences)
	}

}

// Run after controls and a steps block are completed to
// post process any data and provide downstream results
func _Scarfree_siteremove_orfcheck_wtypeAnalysis(_ctx context.Context, _input *Scarfree_siteremove_orfcheck_wtypeInput, _output *Scarfree_siteremove_orfcheck_wtypeOutput) {
}

// A block of tests to perform to validate that the sample was processed correctly
// Optionally, destructive tests can be performed to validate results on a
// dipstick basis
func _Scarfree_siteremove_orfcheck_wtypeValidation(_ctx context.Context, _input *Scarfree_siteremove_orfcheck_wtypeInput, _output *Scarfree_siteremove_orfcheck_wtypeOutput) {
}
func _Scarfree_siteremove_orfcheck_wtypeRun(_ctx context.Context, input *Scarfree_siteremove_orfcheck_wtypeInput) *Scarfree_siteremove_orfcheck_wtypeOutput {
	output := &Scarfree_siteremove_orfcheck_wtypeOutput{}
	_Scarfree_siteremove_orfcheck_wtypeSetup(_ctx, input)
	_Scarfree_siteremove_orfcheck_wtypeSteps(_ctx, input, output)
	_Scarfree_siteremove_orfcheck_wtypeAnalysis(_ctx, input, output)
	_Scarfree_siteremove_orfcheck_wtypeValidation(_ctx, input, output)
	return output
}

func Scarfree_siteremove_orfcheck_wtypeRunSteps(_ctx context.Context, input *Scarfree_siteremove_orfcheck_wtypeInput) *Scarfree_siteremove_orfcheck_wtypeSOutput {
	soutput := &Scarfree_siteremove_orfcheck_wtypeSOutput{}
	output := _Scarfree_siteremove_orfcheck_wtypeRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func Scarfree_siteremove_orfcheck_wtypeNew() interface{} {
	return &Scarfree_siteremove_orfcheck_wtypeElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &Scarfree_siteremove_orfcheck_wtypeInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _Scarfree_siteremove_orfcheck_wtypeRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &Scarfree_siteremove_orfcheck_wtypeInput{},
			Out: &Scarfree_siteremove_orfcheck_wtypeOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wunit.Make_units
)

type Scarfree_siteremove_orfcheck_wtypeElement struct {
	inject.CheckedRunner
}

type Scarfree_siteremove_orfcheck_wtypeInput struct {
	BlastSeqswithNoName           bool
	Constructname                 string
	EndsAlreadyadded              bool
	Enzymename                    string
	ExporttoFastaFile             bool
	ORFstoConfirm                 []string
	OtherEnzymeSitesToRemove      []string
	RemoveproblemRestrictionSites bool
	Seqsinorder                   []wtype.DNASequence
	Vector                        wtype.DNASequence
}

type Scarfree_siteremove_orfcheck_wtypeOutput struct {
	Endreport              string
	NewDNASequence         wtype.DNASequence
	ORFmissing             bool
	ORIpresent             bool
	OriginalParts          []wtype.DNASequence
	PartsWithSitesRemoved  []wtype.DNASequence
	PartswithOverhangs     []wtype.DNASequence
	Plasmid                bool
	PositionReport         []string
	SelectionMarkerPresent bool
	Simulationpass         bool
	Status                 string
	Warnings               error
}

type Scarfree_siteremove_orfcheck_wtypeSOutput struct {
	Data struct {
		Endreport              string
		NewDNASequence         wtype.DNASequence
		ORFmissing             bool
		ORIpresent             bool
		OriginalParts          []wtype.DNASequence
		PartsWithSitesRemoved  []wtype.DNASequence
		PartswithOverhangs     []wtype.DNASequence
		Plasmid                bool
		PositionReport         []string
		SelectionMarkerPresent bool
		Simulationpass         bool
		Status                 string
		Warnings               error
	}
	Outputs struct {
	}
}

func init() {
	if err := addComponent(component.Component{Name: "Scarfree_siteremove_orfcheck_wtype",
		Constructor: Scarfree_siteremove_orfcheck_wtypeNew,
		Desc: component.ComponentDesc{
			Desc: "",
			Path: "antha/component/an/Data/DNA/TypeIISAssembly_design/Scarfree_removesites_checkorfs_wtype.an",
			Params: []component.ParamDesc{
				{Name: "BlastSeqswithNoName", Desc: "", Kind: "Parameters"},
				{Name: "Constructname", Desc: "", Kind: "Parameters"},
				{Name: "EndsAlreadyadded", Desc: "", Kind: "Parameters"},
				{Name: "Enzymename", Desc: "", Kind: "Parameters"},
				{Name: "ExporttoFastaFile", Desc: "", Kind: "Parameters"},
				{Name: "ORFstoConfirm", Desc: "enter each as amino acid sequence\n", Kind: "Parameters"},
				{Name: "OtherEnzymeSitesToRemove", Desc: "", Kind: "Parameters"},
				{Name: "RemoveproblemRestrictionSites", Desc: "", Kind: "Parameters"},
				{Name: "Seqsinorder", Desc: "", Kind: "Parameters"},
				{Name: "Vector", Desc: "", Kind: "Parameters"},
				{Name: "Endreport", Desc: "", Kind: "Data"},
				{Name: "NewDNASequence", Desc: "desired sequence to end up with after assembly\n", Kind: "Data"},
				{Name: "ORFmissing", Desc: "", Kind: "Data"},
				{Name: "ORIpresent", Desc: "", Kind: "Data"},
				{Name: "OriginalParts", Desc: "", Kind: "Data"},
				{Name: "PartsWithSitesRemoved", Desc: "", Kind: "Data"},
				{Name: "PartswithOverhangs", Desc: "parts to order\n", Kind: "Data"},
				{Name: "Plasmid", Desc: "", Kind: "Data"},
				{Name: "PositionReport", Desc: "", Kind: "Data"},
				{Name: "SelectionMarkerPresent", Desc: "", Kind: "Data"},
				{Name: "Simulationpass", Desc: "", Kind: "Data"},
				{Name: "Status", Desc: "", Kind: "Data"},
				{Name: "Warnings", Desc: "", Kind: "Data"},
			},
		},
	}); err != nil {
		panic(err)
	}
}
