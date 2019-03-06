// antha/AnthaStandardLibrary/Packages/pubchem/pubchem.go: Part of the Antha language
// Copyright (C) 2015 The Antha authors. All rights reserved.
//
// This program is free software; you can redistribute it and/or
// modify it under the terms of the GNU General Public License
// as published by the Free Software Foundation; either version 2
// of the License, or (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program; if not, write to the Free Software
// Foundation, Inc., 51 Franklin Street, Fifth Floor, Boston, MA  02110-1301, USA.
//
// For more information relating to the software or licensing issues please
// contact license@antha-lang.org or write to the Antha team c/o
// Synthace Ltd. The London Bioscience Innovation Centre
// 2 Royal College St, London NW1 0NH UK

// Package for interacting with the pubchem database to look up chemical properties
package pubchem

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/antha-lang/antha/antha/anthalib/wunit"
)

// https://pubchem.ncbi.nlm.nih.gov/pug_rest/PUG_REST.html#_Toc409516757

/*
Input


The input portion of the URL tells the service which records to use as the subject of the query. This is further subdivided into two or more locations in the URL “path” as follows:

<input specification> = <domain>/<namespace>/<identifiers>

<domain> = substance | compound | assay | <other inputs>

compound domain <namespace> = cid | name | smiles | inchi | sdf | inchikey | formula | <structure search> | <xref> | listkey | <fast search>

<structure search> = {substructure | superstructure | similarity | identity}/{smiles | inchi | sdf | cid}

<fast search> = {fastidentity | fastsimilarity_2d | fastsimilarity_3d | fastsubstructure | fastsuperstructure}/{smiles | inchi | sdf | cid} | fastformula

<xref> = xref / {RegistryID | RN | PubMedID | MMDBID | ProteinGI | NucleotideGI | TaxonomyID | MIMID | GeneID | ProbeID | PatentID}

substance domain <namespace> = sid | sourceid/<source name> | sourceall/<source name> | name | <xref> | listkey

<source name> = any valid PubChem depositor name

assay domain <namespace> = aid | listkey | type/<assay type> | sourceall/<source name> | target/<assay target> | activity/<activity column name>

<assay type> = all | confirmatory | doseresponse | onhold | panel | rnai | screening | summary | cellbased | biochemical | invivo | invitro | activeconcentrationspecified

<assay target> = gi | proteinname | geneid | genesymbol

<identifiers> = comma-separated list of positive integers (e.g. cid, sid, aid) or identifier strings (source, inchikey, formula); in some cases only a single identifier string (name, smiles, xref; inchi, sdf by POST only)

<other inputs> = sources / [substance, assay] |sourcetable | conformers



For example, to access CID 2244 (aspirin), one would construct the first part of the URL this way:

http://pubchem.ncbi.nlm.nih.gov/rest/pug/compound/cid/2244/<operation specification>/[<output specification>]



Some source names contain the ‘/’ (forward slash) character, which is incompatible with the URL syntax; for these, replace the ‘/’ with a ‘.’ (period) in the URL. Other special characters may need to be escaped, such as ‘&’ should be replaced by ‘%26’. For example:

http://pubchem.ncbi.nlm.nih.gov/rest/pug/substance/sourceid/DTP.NCI/<operation specification>/[<output specification>]

*/

func makeInputspec(domain string, namespace string, identifiers []string) (inputspec string) {
	// see comment above for structure
	//<domain> = substance | compound | assay | <other inputs>

	array := make([]string, 0)
	array = append(array, domain, namespace)
	for i := 0; i < len(identifiers); i++ {
		array = append(array, identifiers[i])
	}
	inputspec = strings.Join(array, "/")
	/*if operation_options != "" {
		array = append(array,operation_options)
	}*/

	return inputspec
}

/*
Operation


The operation part of the URL tells the service what to do with the input records – such as to retrieve whole record data blobs or specific properties of a compound, etc. The construction of this part of the “path” will depend on what the operation is. Currently, if no operation is specified at all, the default is to retrieve the entire record. What operations are available are, of course, dependent on the input domain – that is, certain operations are applicable only to compounds and not assays, for example.

compound domain <operation specification> = record | <compound property> | synonyms | sids | cids | aids | assaysummary | classification | <xrefs> | description | conformers

<compound property> = property / [comma-separated list of property tags]

substance domain <operation specification> = record | synonyms | sids | cids | aids | assaysummary | classification | <xrefs> | description

<xrefs> = xrefs / [comma-separated list of xrefs tags]

assay domain <operation specification> = record | concise | aids | sids | cids | description | targets/<target type> | <doseresponse> | summary | classification | xrefs

target_type = {ProteinGI, ProteinName, GeneID, GeneSymbol}

<doseresponse> = doseresponse/sid



For example, to access the molecular formula and InChI key for CID 2244, one would use a URL like:

http://pubchem.ncbi.nlm.nih.gov/rest/pug/compound/cid/2244/property/MolecularFormula,InChIKey/[<output specification>]

*/

func makeOperationspec(spec string, optionalconditions []string) (operationspec string) {
	// see comment above for structure
	// spec = 1 of: record | <compound property> | synonyms | sids | cids | aids | assaysummary | classification | <xrefs> | description | conformers

	array := make([]string, 0)
	if spec == "compound property" {
		spec = "property"
	}

	array = append(array, spec)
	if spec == "property" || spec == "xrefs" {
		var commaseperatedoptconditions string
		optconditions := make([]string, 0)
		for i := 0; i < len(optionalconditions); i++ {
			optconditions = append(optconditions, optionalconditions[i])
			commaseperatedoptconditions = strings.Join(optconditions, ",")
		}
		array = append(array, commaseperatedoptconditions)
	}
	operationspec = strings.Join(array, "/")
	return operationspec
}

/*
Output


The final portion of the URL tells the service what output format is desired. Note that this is formally optional, as output format can also be specified in the HTTP Accept field of the request header – see below for more detail.

<output specification> = XML | ASNT | ASNB | JSON | JSONP [ ?callback=<callback name> ] | SDF | CSV | PNG | TXT



ASNT is NCBI’s text (human-readable) variant of ASN.1; ASNB is standard binary ASN.1 and is currently returned as Base64-encoded ascii text. Note that not all formats are applicable to the results of all operations; one cannot, for example, retrieve a whole compound record as CSV or a property table as SDF. TXT output is only available in a restricted set of cases where all the information is the same – for example, synonyms for a single CID where there is one synonym per line.

For example, to access the molecular formula for CID 2244 in JSON format, one would use the (now complete) URL:

http://pubchem.ncbi.nlm.nih.gov/rest/pug/compound/cid/2244/property/MolecularFormula/JSON

JSONP takes an optional callback function name (which defaults to “callback” if not specified). For example:

http://pubchem.ncbi.nlm.nih.gov/rest/pug/compound/cid/2244/property/MolecularFormula/JSONP?callback=my_callback

*/
func makeOutputspec(spec string, optionalcallbackname string) (outputspec string) {
	// see comment above for structure
	// spec = 1 of: <output specification> = XML | ASNT | ASNB | JSON | JSONP [ ?callback=<callback name> ] | SDF | CSV | PNG | TXT

	array := make([]string, 0)

	array = append(array, spec)
	if spec == "JSONP" {

		array = append(array, optionalcallbackname)
	}
	outputspec = strings.Join(array, "/")
	return outputspec
}

func pugLookup(inputspec string, operationspec string, outputspec string, operation_options string) ([]byte, error) {

	/*<input specification>/<operation specification>/[<output specification>][?<operation_options>]

	http://pubchem.ncbi.nlm.nih.gov/rest/pug/compound/name/glucose/property/MolecularFormula,MolecularWeight/JSON
	*/
	res, err := http.Get(fmt.Sprintf("http://pubchem.ncbi.nlm.nih.gov/rest/pug/%s/%s/%s/%s", inputspec, operationspec, outputspec, operation_options))
	if err != nil {
		return nil, err
	}
	if res != nil {
		defer res.Body.Close() // nolint
	}

	output, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return output, nil
}

func Compoundproperties(name string) (string, error) {
	// need this structure: http://pubchem.ncbi.nlm.nih.gov/rest/pug/compound/name/glucose/property/MolecularFormula,MolecularWeight/JSON

	inputspec := makeInputspec("compound", "name", []string{name})
	operationspec := makeOperationspec("property", []string{"MolecularFormula", "MolecularWeight"})
	outputspec := makeOutputspec("JSON", "")
	output, err := pugLookup(inputspec, operationspec, outputspec, "")
	if err != nil {
		return "", err
	}

	return string(output), nil
}

type Pubchemtable struct {
	Pubchemjson `json:"PropertyTable"`
}

type Pubchemjson struct {
	Propertytable []Properties `json:"Properties"`
}

type Properties struct {
	MolecularFormula string  `json:"MolecularFormula"`
	MolecularWeight  float64 `json:"MolecularWeight"`
	CID              int     `json:"CID"`
}

// The principle type returned from querying the pubchem database if the molecule is not defined as a substance.
type Molecule struct {
	Name             string
	MolecularFormula string  `json:"MolecularFormula"`
	MolecularWeight  float64 `json:"MolecularWeight"`
	CID              int     `json:"CID"`
}

// Converts a concentration in mol/L to a g/L concentration
func (molecule Molecule) GramPerL(molarity wunit.Concentration) (weightpervol wunit.Concentration) {
	weightpervol = molarity.GramPerL(molecule.MolecularWeight)
	return
}

// Converts a concentration in g/L to a mol/L concentration
func (molecule Molecule) MolPerL(weightpervol wunit.Concentration) (molarity wunit.Concentration) {
	molarity = weightpervol.MolPerL(molecule.MolecularWeight)
	return
}

// Returns the concentration in g/L required for 1 mol/L of the molecule
func (molecule Molecule) MolarConcentration() (weightpervol wunit.Concentration) {
	molarity := wunit.NewConcentration(1, "M/l")
	weightpervol = molarity.GramPerL(molecule.MolecularWeight)
	return
}

// Returns the mass in g required for 1 mol of the molecule
func (molecule Molecule) MolarMass() (weight wunit.Mass) {
	weight = wunit.NewMass(molecule.MolecularWeight, "g")
	return
}

// Returns a summary of the molecule properties
func (molecule Molecule) ToString() string {
	return fmt.Sprint("Name: ", molecule.Name, "Formula: ", molecule.MolecularFormula, "MolecularWeight: ", molecule.MolecularWeight, "g/mol", "CID: ", molecule.CID)
}

// Lookup and make a molecule based on molecule name
func MakeMolecule(name string) (Molecule, error) {
	// need this structure: http://pubchem.ncbi.nlm.nih.gov/rest/pug/compound/name/glucose/property/MolecularFormula,MolecularWeight/JSON

	inputspec := makeInputspec("compound", "name", []string{name})
	operationspec := makeOperationspec("property", []string{"MolecularFormula", "MolecularWeight"})
	outputspec := makeOutputspec("JSON", "")
	output, err := pugLookup(inputspec, operationspec, outputspec, "")
	if err != nil {
		return Molecule{}, err
	}

	var pubchemtable Pubchemtable

	if err := json.Unmarshal(output, &pubchemtable); err != nil {
		return Molecule{}, err

	}

	var molecule Molecule
	molecule.Name = name

	if len(pubchemtable.Propertytable) < 1 {
		return molecule, errors.New(fmt.Sprint("No property table for ", name, " got this from pubchem: ", pubchemtable))
	}
	molecule.CID = pubchemtable.Propertytable[0].CID
	molecule.MolecularFormula = pubchemtable.Propertytable[0].MolecularFormula
	molecule.MolecularWeight = pubchemtable.Propertytable[0].MolecularWeight
	return molecule, nil
}

// Make an array of molecules based on molecule names
// If any errors are encountered they will be aggregated and returned at the end.
func MakeMolecules(names []string) ([]Molecule, error) {
	var molecules []Molecule
	var errs []string
	for _, name := range names {
		molecule, err := MakeMolecule(name)
		if err != nil {
			errs = append(errs, err.Error())
		} else {
			molecules = append(molecules, molecule)
		}
	}
	if len(errs) > 0 {
		return molecules, fmt.Errorf(strings.Join(errs, ";"))
	}
	return molecules, nil
}

// distinct from a molecule in that a substance does not possess a clear Molecular formula or molecular wieght. e.g. Bovine Serum Albumin.
// not currently implemented as an output returned by querying the pubchem database
type Substance struct {
	Name string
	SID  int `json:"SID"`
}
