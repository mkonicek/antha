// This protocol is intended to design assembly parts using either an assembly standard or a specified enzyme.
// parts are added as biobrick IDs, or looked up from the inventory package
// A simulation is performed and status returned to the user
// The user can also specify the names of enzyme sites they wish to avoid to check if these are present in the
// new dna sequence (if simulation passes that is).

package lib

import (
	"fmt"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/Inventory"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/enzymes"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/enzymes/lookup"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/export"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/igem"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/component"
	"github.com/antha-lang/antha/execute"
	"github.com/antha-lang/antha/inject"
	"golang.org/x/net/context"
	"strings"
)

// Input parameters for this protocol (data)

// Physical Inputs to this protocol with types

// Physical outputs from this protocol with types

// Data which is returned from this protocol, and data types

// i.e. parts to order

// Input Requirement specification
func _TypeIISAssembly_designRequirements() {

}

// Conditions to run on startup
func _TypeIISAssembly_designSetup(_ctx context.Context, _input *TypeIISAssembly_designInput) {}

// The core process for this protocol, with the steps to be performed
// for every input
func _TypeIISAssembly_designSteps(_ctx context.Context, _input *TypeIISAssembly_designInput, _output *TypeIISAssembly_designOutput) {
	//var msg string
	// set warnings reported back to user to none initially
	warnings := make([]string, 1)
	warnings[0] = "none"
	var nofeatures []wtype.Feature

	/* find sequence data from keyword; looking it up by a given name in an inventory
	   or by biobrick ID from iGem parts registry */
	partsinorder := make([]wtype.DNASequence, 0)

	_output.Status = "all parts available"
	for _, part := range _input.Partsinorder {

		var partDNA = wtype.DNASequence{"", "", false, false, wtype.Overhang{0, 0, 0, "", false}, wtype.Overhang{0, 0, 0, "", false}, "", nofeatures}

		if strings.Contains(part, "BBa_") == true {

			fmt.Println("looking in igem registry for: ", part)

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
			var found bool
			partDNA, found = Inventory.Partslist()[part]

			if !found {
				_output.Status = fmt.Sprintln("part not found in Inventory so element aborted!")
				execute.Errorf(_ctx, _output.Status)

			}
		}
		partsinorder = append(partsinorder, partDNA)
	}

	// or Look up parts from registry according to properties (this will take a couple of minutes the first time)
	exacttypeonly := true

	subparts, _ := igem.FilterRegistry("REPORTER", []string{"Fluorescent", "A "}, exacttypeonly)
	partdetails := igem.LookUp(subparts)
	//fmt.Println(partdetails)

	// this can be slow if there are many parts to check (~2 seconds per block of 14 parts)
	for _, subpart := range subparts {
		if strings.Contains(partdetails.Description(subpart), "RED") &&
			strings.Contains(partdetails.Results(subpart), "WORKS") {
			_output.BackupParts = append(_output.BackupParts, subpart)

		}
	}

	// lookup vector sequence
	vectordata := Inventory.Partslist()[_input.Vector]

	//lookup restriction enzyme
	restrictionenzyme := enzymes.Enzymelookup[_input.AssemblyStandard][_input.Level]

	// (1) Add standard overhangs using chosen assembly standard
	//PartswithOverhangs = enzymes.MakeStandardTypeIIsassemblyParts(partsinorder, AssemblyStandard, Level, PartMoClotypesinorder)

	// OR (2) Add overhangs for scarfree assembly based on part seqeunces only, i.e. no Assembly standard
	fmt.Println("partsinorder: ", partsinorder)
	_output.PartswithOverhangs = enzymes.MakeScarfreeCustomTypeIIsassemblyParts(partsinorder, vectordata, restrictionenzyme)

	// perfrom mock digest to test fragement overhangs (fragments are hidden by using _, )
	_, stickyends5, stickyends3 := enzymes.TypeIIsdigest(vectordata, restrictionenzyme)

	// Check that assembly is feasible with designed parts by simulating assembly of the sequences with the chosen enzyme
	assembly := enzymes.Assemblyparameters{_input.Constructname, restrictionenzyme.Name, vectordata, _output.PartswithOverhangs}
	status, numberofassemblies, sitesfound, newDNASequence, _ := enzymes.Assemblysimulator(assembly)

	// The default sitesfound produced from the assembly simulator only checks to SapI and BsaI so we'll repeat with the enzymes declared in parameters
	// first lookup enzyme properties
	enzlist := make([]wtype.RestrictionEnzyme, 0)
	for _, site := range _input.RestrictionsitetoAvoid {
		enzsite := lookup.EnzymeLookup(site)
		enzlist = append(enzlist, enzsite)
	}
	othersitesfound := enzymes.Restrictionsitefinder(newDNASequence, enzlist)

	for _, site := range sitesfound {
		othersitesfound = append(othersitesfound, site)
	}

	// Now let's find out the size of fragments we would get if digested with a common site cutter
	tspEI := lookup.EnzymeLookup("TspEI")

	Testdigestionsizes := enzymes.RestrictionMapper(newDNASequence, tspEI)

	// allow the data to be exported by capitalising the first letter of the variable
	_output.Sitesfound = othersitesfound

	_output.NewDNASequence = newDNASequence
	if status == "Yay! this should work" && numberofassemblies == 1 {

		_output.Simulationpass = true
	} else {
		fmt.Println(status)
	}

	_output.Warnings = strings.Join(warnings, ";")

	// Export sequences to order into a fasta file

	partswithOverhangs := make([]*wtype.DNASequence, 0)
	for _, part := range _output.PartswithOverhangs {
		export.ExportFasta(_input.Constructname, &part)
		partswithOverhangs = append(partswithOverhangs, &part)

	}
	export.Makefastaserial(_input.Constructname, partswithOverhangs)

	//partstoorder := ansi.Color(fmt.Sprintln("PartswithOverhangs", PartswithOverhangs),"red")
	partstoorder := fmt.Sprintln("PartswithOverhangs", _output.PartswithOverhangs)

	// Print status
	if _output.Status != "all parts available" {
		_output.Status = fmt.Sprintln(_output.Status)
	} else {
		_output.Status = fmt.Sprintln(
			"Warnings:", _output.Warnings,
			"Simulationpass=", _output.Simulationpass,
			"Status: ", status,
			"Back up parts found (Reported to work!)", _output.BackupParts,
			"NewDNASequence", _output.NewDNASequence,
			//"partonewithoverhangs", partonewithoverhangs,
			//"Vector",vectordata,
			"Vector digest:", stickyends5, stickyends3,
			partstoorder,
			"Sitesfound", _output.Sitesfound,
			"Partsinorder=", _input.Partsinorder, partsinorder,
			"Test digestion sizes with TspEI", Testdigestionsizes,
		//"Restriction Enzyme=",restrictionenzyme,
		)
	}

}

// Run after controls and a steps block are completed to
// post process any data and provide downstream results
func _TypeIISAssembly_designAnalysis(_ctx context.Context, _input *TypeIISAssembly_designInput, _output *TypeIISAssembly_designOutput) {
}

// A block of tests to perform to validate that the sample was processed correctly
// Optionally, destructive tests can be performed to validate results on a
// dipstick basis
func _TypeIISAssembly_designValidation(_ctx context.Context, _input *TypeIISAssembly_designInput, _output *TypeIISAssembly_designOutput) {
}
func _TypeIISAssembly_designRun(_ctx context.Context, input *TypeIISAssembly_designInput) *TypeIISAssembly_designOutput {
	output := &TypeIISAssembly_designOutput{}
	_TypeIISAssembly_designSetup(_ctx, input)
	_TypeIISAssembly_designSteps(_ctx, input, output)
	_TypeIISAssembly_designAnalysis(_ctx, input, output)
	_TypeIISAssembly_designValidation(_ctx, input, output)
	return output
}

func TypeIISAssembly_designRunSteps(_ctx context.Context, input *TypeIISAssembly_designInput) *TypeIISAssembly_designSOutput {
	soutput := &TypeIISAssembly_designSOutput{}
	output := _TypeIISAssembly_designRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func TypeIISAssembly_designNew() interface{} {
	return &TypeIISAssembly_designElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &TypeIISAssembly_designInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _TypeIISAssembly_designRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &TypeIISAssembly_designInput{},
			Out: &TypeIISAssembly_designOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wunit.Make_units
)

type TypeIISAssembly_designElement struct {
	inject.CheckedRunner
}

type TypeIISAssembly_designInput struct {
	AssemblyStandard       string
	Constructname          string
	Level                  string
	PartMoClotypesinorder  []string
	Partsinorder           []string
	RestrictionsitetoAvoid []string
	Vector                 string
}

type TypeIISAssembly_designOutput struct {
	BackupParts        []string
	NewDNASequence     wtype.DNASequence
	PartswithOverhangs []wtype.DNASequence
	Simulationpass     bool
	Sitesfound         []enzymes.Restrictionsites
	Status             string
	Warnings           string
}

type TypeIISAssembly_designSOutput struct {
	Data struct {
		BackupParts        []string
		NewDNASequence     wtype.DNASequence
		PartswithOverhangs []wtype.DNASequence
		Simulationpass     bool
		Sitesfound         []enzymes.Restrictionsites
		Status             string
		Warnings           string
	}
	Outputs struct {
	}
}

func init() {
	if err := addComponent(component.Component{Name: "TypeIISAssembly_design",
		Constructor: TypeIISAssembly_designNew,
		Desc: component.ComponentDesc{
			Desc: "",
			Path: "src/github.com/antha-lang/antha/antha/component/an/Data/DNA/TypeIISAssembly_design/TypeIISAssembly_design.an",
			Params: []component.ParamDesc{
				{Name: "AssemblyStandard", Desc: "", Kind: "Parameters"},
				{Name: "Constructname", Desc: "", Kind: "Parameters"},
				{Name: "Level", Desc: "", Kind: "Parameters"},
				{Name: "PartMoClotypesinorder", Desc: "", Kind: "Parameters"},
				{Name: "Partsinorder", Desc: "", Kind: "Parameters"},
				{Name: "RestrictionsitetoAvoid", Desc: "", Kind: "Parameters"},
				{Name: "Vector", Desc: "", Kind: "Parameters"},
				{Name: "BackupParts", Desc: "", Kind: "Data"},
				{Name: "NewDNASequence", Desc: "", Kind: "Data"},
				{Name: "PartswithOverhangs", Desc: "i.e. parts to order\n", Kind: "Data"},
				{Name: "Simulationpass", Desc: "", Kind: "Data"},
				{Name: "Sitesfound", Desc: "", Kind: "Data"},
				{Name: "Status", Desc: "", Kind: "Data"},
				{Name: "Warnings", Desc: "", Kind: "Data"},
			},
		},
	}); err != nil {
		panic(err)
	}
}
