// This protocol is intended to check sequences for restriction sites and remove according to
// specified conditions

package lib

import (
	"fmt"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/enzymes"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/enzymes/lookup"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/sequences"
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

//wtype.DNASequence

// Physical Inputs to this protocol with types

// Physical outputs from this protocol with types

// Data which is returned from this protocol, and data types

// i.e. parts to order

// Input Requirement specification
func _RemoveRestrictionSitesRequirements() {

}

// Conditions to run on startup
func _RemoveRestrictionSitesSetup(_ctx context.Context, _input *RemoveRestrictionSitesInput) {}

// The core process for this protocol, with the steps to be performed
// for every input
func _RemoveRestrictionSitesSteps(_ctx context.Context, _input *RemoveRestrictionSitesInput, _output *RemoveRestrictionSitesOutput) {

	Sequence := wtype.MakeLinearDNASequence("Test", _input.Sequencekey)

	// set warnings reported back to user to none initially
	warnings := make([]string, 0)

	// first lookup enzyme properties for all enzymes and make a new array
	enzlist := make([]wtype.RestrictionEnzyme, 0)
	for _, site := range _input.RestrictionsitetoAvoid {
		enzsite := lookup.EnzymeLookup(site)
		enzlist = append(enzlist, enzsite)
	}

	// check for sites in the sequence
	sitesfound := enzymes.Restrictionsitefinder(Sequence, enzlist)

	// if no sites found skip to restriction map stage
	if len(sitesfound) == 0 {
		_output.Warnings = "none"
		_output.Status = "No sites found in sequence to remove so same sequence returned"
		_output.SiteFreeSequence = Sequence
		_output.Sitesfoundinoriginal = sitesfound

	} else {

		// make a list of sequences to avoid before modifying the sequence
		allsitestoavoid := make([]string, 0)

		// add all restriction sites (we need this step since the functions coming up require strings)
		for _, enzy := range enzlist {
			allsitestoavoid = append(allsitestoavoid, enzy.RecognitionSequence)
		}

		for _, site := range sitesfound {
			if site.Sitefound {

				var tempseq wtype.DNASequence
				var err error

				orfs := sequences.FindallORFs(Sequence.Seq)
				warnings = append(warnings, text.Print("orfs: ", orfs))
				features := sequences.ORFs2Features(orfs)

				//set up a boolean to change to true if a sequence is found in an ORF
				foundinorf := false
				//set up an index for each orf found with site within it (need enzyme name too but will recheck all anyway!)
				orfswithsites := make([]int, 0)

				if len(orfs) > 0 {
					for i, orf := range orfs {

						// change func to handle this step of making dnaseq first

						dnaseq := wtype.MakeLinearDNASequence("orf"+strconv.Itoa(i), orf.DNASeq)

						foundinorfs := enzymes.Restrictionsitefinder(dnaseq, enzlist) // won't work yet orf is actually type features

						for _, siteinorf := range foundinorfs {
							if siteinorf.Sitefound == true {
								foundinorf = true
							}
						}

						if foundinorf == true {

							warning := text.Print("sites found in orf"+dnaseq.Nm, orf)
							warnings = append(warnings, warning)
						}
					}
				}
				if _input.RemoveifnotinORF {
					if foundinorf == false {
						tempseq, err = sequences.RemoveSite(Sequence, site.Enzyme, allsitestoavoid)
						if err != nil {
							warning := text.Print("removal of site failed! improve your algorithm!", err.Error())
							warnings = append(warnings, warning)

						}
						_output.SiteFreeSequence = tempseq

						// all done if all sites are not in orfs!
						// make proper remove allsites func
					}
					if foundinorf == true {

						_output.SiteFreeSequence, err = sequences.RemoveSitesOutsideofFeatures(Sequence, site.Enzyme.RecognitionSequence, sequences.ReplaceBycomplement, features)
						if err != nil {
							warnings = append(warnings, err.Error())
						}
					}
				} //		}else {
				if _input.PreserveTranslatedseq {
					// make func to check codon and swap site to preserve aa sequence product
					for _, orfnumber := range orfswithsites {

						for _, position := range site.Positions("ALL") {
							orfcoordinates := sequences.MakeStartendPair(orfs[orfnumber].StartPosition, orfs[orfnumber].EndPosition)
							tempseq, _, _, err = sequences.ReplaceCodoninORF(tempseq, orfcoordinates, position, allsitestoavoid)
							if err != nil {
								warning := text.Print("removal of site from orf "+strconv.Itoa(orfnumber), " failed! improve your algorithm! "+err.Error())
								warnings = append(warnings, warning)
							}
						}

					}
				}

				_output.SiteFreeSequence = tempseq
			}
		}
	}

	// Now let's find out the size of fragments we would get if digested with a common site cutter
	mapenz := lookup.EnzymeLookup(_input.EnzymeforRestrictionmapping)

	_output.FragmentSizesfromRestrictionmapping = enzymes.RestrictionMapper(Sequence, mapenz)

	// allow the data to be exported by capitalising the first letter of the variable
	_output.Sitesfoundinoriginal = sitesfound

	_output.Warnings = strings.Join(warnings, ";")

	// Print status
	if _output.Status == "" {
		_output.Status = fmt.Sprintln("Something went wrong!")
	} else {
		_output.Status = fmt.Sprintln(
			text.Print("Warnings:", _output.Warnings),
			text.Print("Sequence", Sequence),
			text.Print("Sitesfound", _output.Sitesfoundinoriginal),
			text.Print("Test digestion sizes with"+_input.EnzymeforRestrictionmapping, _output.FragmentSizesfromRestrictionmapping),
		)
	}

}

// Run after controls and a steps block are completed to
// post process any data and provide downstream results
func _RemoveRestrictionSitesAnalysis(_ctx context.Context, _input *RemoveRestrictionSitesInput, _output *RemoveRestrictionSitesOutput) {
}

// A block of tests to perform to validate that the sample was processed correctly
// Optionally, destructive tests can be performed to validate results on a
// dipstick basis
func _RemoveRestrictionSitesValidation(_ctx context.Context, _input *RemoveRestrictionSitesInput, _output *RemoveRestrictionSitesOutput) {
}
func _RemoveRestrictionSitesRun(_ctx context.Context, input *RemoveRestrictionSitesInput) *RemoveRestrictionSitesOutput {
	output := &RemoveRestrictionSitesOutput{}
	_RemoveRestrictionSitesSetup(_ctx, input)
	_RemoveRestrictionSitesSteps(_ctx, input, output)
	_RemoveRestrictionSitesAnalysis(_ctx, input, output)
	_RemoveRestrictionSitesValidation(_ctx, input, output)
	return output
}

func RemoveRestrictionSitesRunSteps(_ctx context.Context, input *RemoveRestrictionSitesInput) *RemoveRestrictionSitesSOutput {
	soutput := &RemoveRestrictionSitesSOutput{}
	output := _RemoveRestrictionSitesRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func RemoveRestrictionSitesNew() interface{} {
	return &RemoveRestrictionSitesElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &RemoveRestrictionSitesInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _RemoveRestrictionSitesRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &RemoveRestrictionSitesInput{},
			Out: &RemoveRestrictionSitesOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wunit.Make_units
)

type RemoveRestrictionSitesElement struct {
	inject.CheckedRunner
}

type RemoveRestrictionSitesInput struct {
	EnzymeforRestrictionmapping string
	PreserveTranslatedseq       bool
	RemoveifnotinORF            bool
	RestrictionsitetoAvoid      []string
	Sequencekey                 string
}

type RemoveRestrictionSitesOutput struct {
	FragmentSizesfromRestrictionmapping []int
	SiteFreeSequence                    wtype.DNASequence
	Sitesfoundinoriginal                []enzymes.Restrictionsites
	Status                              string
	Warnings                            string
}

type RemoveRestrictionSitesSOutput struct {
	Data struct {
		FragmentSizesfromRestrictionmapping []int
		SiteFreeSequence                    wtype.DNASequence
		Sitesfoundinoriginal                []enzymes.Restrictionsites
		Status                              string
		Warnings                            string
	}
	Outputs struct {
	}
}

func init() {
	if err := addComponent(component.Component{Name: "RemoveRestrictionSites",
		Constructor: RemoveRestrictionSitesNew,
		Desc: component.ComponentDesc{
			Desc: "",
			Path: "src/github.com/antha-lang/antha/antha/component/an/Data/DNA/RestrictionSiteRemover/RemoveRestrictionSites.an",
			Params: []component.ParamDesc{
				{Name: "EnzymeforRestrictionmapping", Desc: "", Kind: "Parameters"},
				{Name: "PreserveTranslatedseq", Desc: "", Kind: "Parameters"},
				{Name: "RemoveifnotinORF", Desc: "", Kind: "Parameters"},
				{Name: "RestrictionsitetoAvoid", Desc: "", Kind: "Parameters"},
				{Name: "Sequencekey", Desc: "wtype.DNASequence\n", Kind: "Parameters"},
				{Name: "FragmentSizesfromRestrictionmapping", Desc: "", Kind: "Data"},
				{Name: "SiteFreeSequence", Desc: "i.e. parts to order\n", Kind: "Data"},
				{Name: "Sitesfoundinoriginal", Desc: "", Kind: "Data"},
				{Name: "Status", Desc: "", Kind: "Data"},
				{Name: "Warnings", Desc: "", Kind: "Data"},
			},
		},
	}); err != nil {
		panic(err)
	}
}
