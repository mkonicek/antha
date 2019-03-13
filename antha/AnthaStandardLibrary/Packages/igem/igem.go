//Part of the Antha language
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

// Package for interacting with the iGem registry
package igem

import (
	"bufio"
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/AnthaPath"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/search"
)

// http://parts.igem.org/Registry_API

/*
Input
FASTA Formatted Sequences
We will provide a daily update of part sequences, types, subparts, status, and short description for each part and for all parts. Go to http://parts.igem.org/fasta/parts/BBa_C0040 (substitute our desired part name for BBa_C0040) and you will receive a FASTA formatted file with the part's sequence. The header line has this format:
'>'[Part name] [First character of status] [Part Id Number] [Part type] [Short description]
Note: the short description has unusual characters converted to their two-digit hex value.
You can also get all of the parts in a single download (about 30 megabytes) as http://parts.igem.org/fasta/parts/All_Parts.
We are not yet updating these files on a daily basis. -- Randy May 30, 2009

XML Part Information
The database information about each part is now available as XML. You can get information about a part by entering a URL like this one. If your browser parses and displays XML in a formatted way, it will make some sense. The URL ends ...xml/part; you can follow with a list of part names separated by periods.
http://parts.igem.org/xml/part.BBa_B0034
The information for a part includes:
Part name, type, nickname, short description, status, rating, date entered, authors, quality
Lists of subparts (as specified by the designer, at the basic part level, and with scars)
Sequence
Features
Parameters
Categories
DNA Samples (not enabled not)
References (not enabled now)
Groups (not enabled now)
If you enter 'recursive' as the first part name, the returned XML will include details about all the subparts of this part.

*/

const (
	registryFile = "iGem_registry.txt"
)

// Part Name classifications
/*
BBa_B... = Generic basic parts such as Terminators, DNA, and Ribosome Binding Site
BBa_C... = Protein coding parts
BBa_E... = Reporter parts
BBa_F... = Signalling parts
BBa_G... = Primer parts
BBa_I... = IAP 2003, 2004 project parts
BBa_J... = iGEM project parts
BBa_M... = Tag parts
BBa_P... = Protein Generator parts
BBa_Q... = Inverter parts
BBa_R... = Regulatory parts
BBa_S... = Intermediate parts
BBa_V... = Cell strain parts
*/
const (
	GENERIC          = "BBa_B"
	PROTEINCODING    = "BBa_C"
	REPORTER         = "BBa_E"
	SIGNALLING       = "BBa_F"
	PRIMER           = "BBa_G"
	IAPPROJECT       = "BBa_I"
	IGEMPROJECT      = "BBa_J"
	TAG              = "BBa_M"
	PROTEINGENERATOR = "BBa_P"
	INVERTER         = "BBa_Q"
	REGULATORY       = "BBa_R"
	INTERMEDIATE     = "BBa_S"
	CELLSTRAIN       = "BBa_V"
)

var (
	IgemTypeCodes = map[string]string{
		"GENERIC":          "BBa_B",
		"PROTEINCODING":    "BBa_C",
		"REPORTER":         "BBa_E",
		"SIGNALLING":       "BBa_F",
		"PRIMER":           "BBa_G",
		"IAPPROJECT":       "BBa_I",
		"IGEMPROJECT":      "BBa_J",
		"TAG":              "BBa_M",
		"PROTEINGENERATOR": "BBa_P",
		"INVERTER":         "BBa_Q",
		"REGULATORY":       "BBa_R",
		"INTERMEDIATE":     "BBa_S",
		"CELLSTRAIN":       "BBa_V",
	}
)

func validIGEMTypeOptions() string {
	var options []string
	for key := range IgemTypeCodes {
		options = append(options, key)
	}
	return strings.Join(options, "\n")
}

func MakeFastaURL(partname string) (Urlstring string) {
	// see comment above for structure
	//<domain> = substance | compound | assay | <other inputs>

	level1 := "http://parts.igem.org"
	// http://parts.igem.org/fasta/parts/all
	array := make([]string, 0)
	array = append(array, level1, "fasta", "parts", partname)

	Urlstring = strings.Join(array, "/")

	return Urlstring
}

func FetchPartsXML(parts []string) (parsedxml Rsbpml) {

	res, err := http.Get(fmt.Sprintf("http://parts.igem.org/xml/part.%s", strings.Join(parts, ".")))
	if err != nil {
		panic(err)
	}

	output, err := ioutil.ReadAll(res.Body) // this is a slow step!
	if err != nil {
		panic(err)
	}

	return ParseOutput(output)
}

func makeRegistryfile() ([]byte, error) {
	file := filepath.Join(anthapath.Path(), registryFile)

	if _, err := os.Stat(file); os.IsNotExist(err) {
		if err := os.MkdirAll(filepath.Dir(file), 0700); err != nil {
			return nil, err
		}
		// FYI: >34MB file
		res, err := http.Get("http://parts.igem.org/fasta/parts/All_Parts")
		if err != nil {
			return nil, err
		}
		defer res.Body.Close() //nolint

		f, err := os.Create(file)
		if err != nil {
			return nil, err
		}
		defer f.Close() //nolint

		var buf bytes.Buffer
		if _, err := io.Copy(&buf, res.Body); err != nil {
			return nil, err
		}

		if err := ioutil.WriteFile(file, buf.Bytes(), 0666); err != nil {
			return nil, err
		}
		return buf.Bytes(), nil
	}

	return ioutil.ReadFile(file)
}

type FastaPart struct {
	Part_id         string
	Desc            string
	Part_name       string
	Part_short_desc string
	Part_type       string
	Sample_status   string
	Seq_data        string
}

func Build_fasta(header string, seq bytes.Buffer) (Record FastaPart) {

	var record FastaPart

	fields := strings.SplitN(header, " ", 2)
	record.Desc = "`" + fields[1] + "`"

	fields = strings.SplitN(header, " ", 5)

	if len(fields) > 1 {
		record.Part_name = fields[0]
		record.Sample_status = fields[1]
		record.Part_id = fields[2]
		record.Part_type = fields[3]
		record.Part_short_desc = fields[4]

	} else {
		record.Part_name = fields[0]
		record.Part_short_desc = ""
	}

	record.Seq_data = seq.String()

	Record = record

	return Record
}

func FastaParse(fastaFh io.Reader) []FastaPart {
	var outputs []FastaPart

	scanner := bufio.NewScanner(fastaFh)
	header := ""
	var seq bytes.Buffer

	// Loop over the letters in inputString
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if len(line) == 0 {
			continue
		}

		if line[0] == '>' {
			// If we stored a previous identifier, get the DNA string and map to the
			// identifier and clear the string
			if header != "" {
				outputs = append(outputs, Build_fasta(header, seq))
				seq.Reset()
			}

			// Standard FASTA identifiers look like: ">id desc"
			header = line[1:]
		} else {
			// Append here since multi-line DNA strings are possible
			seq.WriteString(line)
		}

	}

	outputs = append(outputs, Build_fasta(header, seq))

	return outputs
}

func CountPartsinRegistryContaining(keystrings []string) (numberofparts int) {
	allparts, err := makeRegistryfile()
	if err != nil {
		return
	}

	fastaFh := bytes.NewReader(allparts)
	for _, record := range FastaParse(fastaFh) {

		if search.ContainsAllStrings(record.Desc, keystrings) {
			numberofparts = numberofparts + 1
		}

	}
	return numberofparts
}

func FilterRegistry(partype string, keystrings []string, exacttypecodeonly bool) (listofpartIDs []string, idtodescriptionmap map[string]string, err error) {

	idtodescriptionmap = make(map[string]string)

	allparts, err := makeRegistryfile()
	if err != nil {
		return
	}

	fastaFh := bytes.NewReader(allparts)
	listofpartIDs = make([]string, 0)

	bba_code, ok := IgemTypeCodes[strings.ToUpper(partype)]

	if !ok {
		err = fmt.Errorf("Part Type %s not found, valid options are: %s", partype, validIGEMTypeOptions())
		return
	}

	for _, record := range FastaParse(fastaFh) {

		if exacttypecodeonly && ok && search.ContainsAllStrings(record.Desc, keystrings) && record.Seq_data != "" && strings.Contains(record.Part_name, bba_code) {

			listofpartIDs = append(listofpartIDs, record.Part_name)
			idtodescriptionmap[record.Part_name] = record.Desc
		} else if !exacttypecodeonly && search.ContainsAllStrings(record.Desc, keystrings) && strings.Contains(strings.ToUpper(record.Part_type), strings.ToUpper(partype)) && record.Seq_data != "" {
			listofpartIDs = append(listofpartIDs, record.Part_name)
			idtodescriptionmap[record.Part_name] = record.Desc
		} else if !exacttypecodeonly && search.ContainsAllStrings(record.Desc, keystrings) && record.Seq_data != "" {
			listofpartIDs = append(listofpartIDs, record.Part_name)
			idtodescriptionmap[record.Part_name] = record.Desc
		}

	}
	return listofpartIDs, idtodescriptionmap, nil
}

func ParseOutput(xmldata []byte) (parsedxml Rsbpml) {

	err := xml.Unmarshal(xmldata, &parsedxml)
	if err != nil {
		fmt.Println("error:", err)
	}

	return parsedxml
}

func LookUp(parts []string) (parsedxml Rsbpml) {
	fmt.Println("number of parts to find in registry", len(parts))
	if len(parts) > 14 {

		partslice := parts[0:14]

		parsedxml = FetchPartsXML(partslice)

		newparsedxml := make([]Part, 0)
		newparsedxml = append(newparsedxml, parsedxml.Partlist[0].Parts...)

		var parsedxml Rsbpml
		partsleft := (len(parts) - len(partslice))
		fmt.Println("parts left = ", partsleft)
		for i := 10; i < len(parts); i = i + 14 {
			partslice = parts[i : i+14]
			parsedxml = FetchPartsXML(partslice)
			newparsedxml = append(newparsedxml, parsedxml.Partlist[0].Parts...)
			var parsedxml Rsbpml
			partsleft = partsleft - len(partslice)
			if partsleft < 14 {
				partslice = parts[len(parts)-partsleft:]
				parsedxml = FetchPartsXML(partslice)

				for _, part := range parsedxml.Partlist[0].Parts {
					newparsedxml = append(newparsedxml, part)
					parsedxml.Partlist[0].Parts = newparsedxml

				}
				{
					break
				}
			}
		}
	} else {

		parsedxml = FetchPartsXML(parts)
	}
	return parsedxml
}

// Add Get funcs to get data from Rsbpml? Would be much faster
func (parsedxml *Rsbpml) Sequence(partname string) (sequence string) {

	for _, part := range parsedxml.Partlist[0].Parts {
		if part.Part_name == partname {
			sequence = part.Sequencelist[0].Seq_data
		}
	}

	sequence = strings.ToUpper(sequence)
	return
}

func (parsedxml *Rsbpml) Type(partname string) (result string) {

	for _, part := range parsedxml.Partlist[0].Parts {
		if part.Part_name == partname {
			result = part.Part_type
		}
	}

	result = strings.ToUpper(result)
	return
}

func (parsedxml *Rsbpml) Categories(partname string) (result Categories) {

	for _, part := range parsedxml.Partlist[0].Parts {
		if part.Part_name == partname {
			result = part.Categories
		}
	}

	return
}

func (parsedxml *Rsbpml) Results(partname string) (result string) {

	for _, part := range parsedxml.Partlist[0].Parts {
		if part.Part_name == partname {
			result = part.Part_results
		}
	}

	result = strings.ToUpper(result)
	return
}

func (parsedxml *Rsbpml) Rating(partname string) (result string) {

	for _, part := range parsedxml.Partlist[0].Parts {
		if part.Part_name == partname {
			result = part.Part_rating
		}
	}

	result = strings.ToUpper(result)
	return
}

func (parsedxml *Rsbpml) Description(partname string) (result string) {

	for _, part := range parsedxml.Partlist[0].Parts {
		if part.Part_name == partname {
			result = part.Part_short_desc
		}
	}

	result = strings.ToUpper(result)

	return
}

func GetSequence(partname string) (sequence string) {

	parts := make([]string, 0)
	parts = append(parts, partname)
	parsedxml := FetchPartsXML(parts)
	sequence = parsedxml.Partlist[0].Parts[0].Sequencelist[0].Seq_data // [0].Seq_data

	sequence = strings.ToUpper(strings.Replace(sequence, "\n", "", -1))

	return sequence
}

func GetType(partname string) (parttype string) {

	parts := make([]string, 0)
	parts = append(parts, partname)
	parsedxml := FetchPartsXML(parts)

	parttype = parsedxml.Partlist[0].Parts[0].Part_type // [0].Seq_data

	return parttype
}

func GetCategories(partname string) (categories Categories) {

	parts := make([]string, 0)
	parts = append(parts, partname)
	parsedxml := FetchPartsXML(parts)

	categories = parsedxml.Partlist[0].Parts[0].Categories // [0].Seq_data

	return categories
}

func GetResults(partname string) (results string) {

	parts := make([]string, 0)
	parts = append(parts, partname)
	parsedxml := FetchPartsXML(parts)

	results = parsedxml.Partlist[0].Parts[0].Part_results // [0].Seq_data

	return results
}

// change to object based method call
func GetResultsfromSubset(partname string, parsedxml Rsbpml) (results string) {

	for _, part := range parsedxml.Partlist[0].Parts {
		if part.Part_name == partname {
			results = part.Part_results
		}
	}

	results = strings.ToUpper(results)

	return results
}

func GetRating(partname string) (rating string) {

	parts := make([]string, 0)
	parts = append(parts, partname)
	parsedxml := FetchPartsXML(parts)

	rating = parsedxml.Partlist[0].Parts[0].Part_rating // [0].Seq_data

	return rating
}

func GetDescription(partname string) (desc string) {

	parts := make([]string, 0)
	parts = append(parts, partname)

	parsedxml := FetchPartsXML(parts)
	desc = parsedxml.Partlist[0].Parts[0].Part_short_desc // [0].Seq_data

	return desc
}

func GetDescriptionfromSubset(partname string, parsedxml Rsbpml) (desc string) {

	for _, part := range parsedxml.Partlist[0].Parts {
		if part.Part_name == partname {
			desc = part.Part_short_desc
		}
	}

	desc = strings.ToUpper(desc)

	return desc
}

func GetPart(partname string) (partproperties Part) {

	parts := make([]string, 0)
	parts = append(parts, partname)
	parsedxml := FetchPartsXML(parts)

	partproperties = parsedxml.Partlist[0].Parts[0] // [0].Seq_data
	return partproperties
}

type Rsbpml struct {
	Partlist []Part_list `xml:"part_list"`
}

type Part_list struct {
	Parts []Part `xml:"part"`
}

type Registryquerier interface {
	GetParts([]string) []Rsbpml
}

type Part struct {
	Part_id            string         `xml:"part_id"`
	Part_name          string         `xml:"part_name"`
	Part_short_name    string         `xml:"part_short_name"`
	Part_short_desc    string         `xml:"part_short_desc"`
	Part_type          string         `xml:"part_type"`
	Release_status     string         `xml:"release_status"`
	Sample_status      string         `xml:"sample_status"`
	Part_results       string         `xml:"part_results"`
	Part_nickname      string         `xml:"part_nickname"`
	Part_rating        string         `xml:"part_rating"`
	Part_url           string         `xml:"part_url"`
	Part_entered       string         `xml:"part_entered"`
	Part_author        string         `xml:"part_author"`
	Deep_subparts      Subparts       `xml:"deep_subparts"`
	Specified_subparts Subparts       `xml:"specified_subparts"`
	Specified_subscars Subscars       `xml:"specified_subscars"`
	Sequencelist       []Sequence     `xml:"sequences"`
	Features           IgemFeatures   `xml:"features"`
	Parameters         IgemParameters `xml:"parameters"`
	Categories         Categories     `xml:"categories"`
	Twins              Twins          `xml:"twins"`
}

type Subparts struct {
	Subparts []Subpart `xml:"subpart"`
}

type Subpart struct {
	Part_id         string `xml:"part_id"`
	Part_name       string `xml:"part_name"`
	Part_short_desc string `xml:"part_short_desc"`
	Part_type       string `xml:"part_type"`
	Part_nickname   string `xml:"part_nickname"`
}

type Subscars struct {
	Subparts []Subpart `xml:"subpart"`
	Scars    []Scar    `xml:"scar"`
}

type Scar struct {
	Scar_id       string `xml:"scar_id"`
	Scar_standard string `xml:"scar_standard"`
	Scar_type     string `xml:"scar_type"`
	Scar_name     string `xml:"scar_name"`
	Scar_nickname string `xml:"scar_nickname"`
	Scar_comments string `xml:"scar_comments"`
	Scar_sequence string `xml:"scar_sequence"`
}

type Sequence struct {
	Seq_data string `xml:"seq_data"`
}

type IgemFeatures struct {
	Features []IgemFeature `xml:"feature"`
}

type IgemFeature struct {
	Id        string `xml:"id"`
	Title     string `xml:"title"`
	Type      string `xml:"type"`
	Direction string `xml:"direction"`
	Startpos  string `xml:"startpos"`
	Endpos    string `xml:"endpos"`
}

type IgemParameters struct {
	Parameters []IgemParameter `xml:"parameter"`
}

type IgemParameter struct {
	Name      string `xml:"name"`
	Value     string `xml:"value"`
	Units     string `xml:"units"`
	Url       string `xml:"url"`
	Id        string `xml:"id"`
	M_date    string `xml:"m_date"`
	User_id   string `xml:"user_id"`
	User_name string `xml:"user_name"`
}

type Categories struct {
	Categories []string `xml:"category"`
}

type Twins struct {
	Twins []string `xml:"twin"`
}
