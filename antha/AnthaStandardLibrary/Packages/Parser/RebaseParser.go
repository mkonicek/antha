// antha/AnthaStandardLibrary/Packages/Parser/RebaseParser.go: Part of the Antha language
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

package parser

import (
	"bufio"
	"io"
	"math"
	"strconv"
	"strings"

	"github.com/antha-lang/antha/antha/anthalib/wtype"
)

/*
REBASE codes for commercial sources of enzymes

B        Life Technologies (6/15)
C        Minotech Biotechnology (8/14)
E        Agilent Technologies (8/13)
I        SibEnzyme Ltd. (6/15)
J        Nippon Gene Co., Ltd. (5/15)
K        Takara Bio Inc. (6/15)
M        Roche Applied Science (6/15)
N        New England Biolabs (6/15)
O        Toyobo Biochemicals (8/14)
Q        Molecular Biology Resources - CHIMERx (6/15)
R        Promega Corporation (6/15)
S        Sigma Chemical Corporation (6/15)
U        Bangalore Genei (11/14)
V        Vivantis Technologies (8/14)
X        EURx Ltd. (8/14)
Y        SinaClon BioScience Co. (8/14)
*/

type TypeIIs struct {
	wtype.RestrictionEnzyme
	Name                              string
	Isoschizomers                     []string
	Topstrand3primedistancefromend    int
	Bottomstrand5primedistancefromend int
}

func RecognitionSeqHandler(RecognitionSeq string) (RecognitionSequence string, EndLength int, Topstrand3primedistancefromend int, Bottomstrand5primedistancefromend int, Class string) {

	// add cases where no "^" is present and two ( / ) are present

	if strings.Count(RecognitionSeq, "(") == 1 &&
		strings.Count(RecognitionSeq, "/") == 1 &&
		strings.Count(RecognitionSeq, ")") == 1 &&
		strings.HasSuffix(RecognitionSeq, ")") {

		split := strings.Split(RecognitionSeq, "(")

		RecognitionSequence = split[0]

		split = strings.Split(split[1], "/")

		lengthint, _ := strconv.Atoi(split[0]) //nolint
		Topstrand3primedistancefromend = lengthint

		split = strings.Split(split[1], ")")

		lengthint, _ = strconv.Atoi(split[0]) //nolint
		Bottomstrand5primedistancefromend = lengthint

		EndLength = int(math.Abs(float64(Bottomstrand5primedistancefromend - Topstrand3primedistancefromend)))
		if Topstrand3primedistancefromend > 0 || Bottomstrand5primedistancefromend > 0 {
			Class = "TypeIIs"
		} else if int(math.Abs(float64(Topstrand3primedistancefromend))) > len(RecognitionSeq) || int(math.Abs(float64(Bottomstrand5primedistancefromend))) > len(RecognitionSeq) {
			Class = "TypeIIs"
		} else {
			Class = "TypeII"
		}
	} else if strings.Count(RecognitionSeq, "^") == 1 {

		split := strings.Split(RecognitionSeq, "^")

		RecognitionSequence = strings.Join(split, "")

		Topstrand3primedistancefromend = -1 * len(split[1])

		Bottomstrand5primedistancefromend = -1 * len(split[0])

		EndLength = int(math.Abs(float64(Bottomstrand5primedistancefromend - Topstrand3primedistancefromend)))

		Class = "TypeII"
	}
	return
}

func Build_rebase(name string, prototype string, recognitionseq string, methylationsite string, commercialsource string, refs string) (Record wtype.RestrictionEnzyme) {

	var record wtype.RestrictionEnzyme

	record.Nm = name
	record.Prototype = prototype

	record.RecognitionSequence,
		record.EndLength,
		record.Topstrand3primedistancefromend,
		record.Bottomstrand5primedistancefromend,
		record.Class = RecognitionSeqHandler(recognitionseq)

	record.MethylationSite = methylationsite
	record.CommercialSource = strings.Split(strings.TrimSpace(commercialsource), "")
	references := strings.Split(refs, ",")

	for _, i := range references {
		if i != "<reference>" {
			j, err := strconv.Atoi(i)
			if err != nil {
				panic(err)
			}

			record.References = append(record.References, j)
		}
	}
	Record = record

	return Record
}

func RebaseParse(rebaseRh io.Reader) []wtype.RestrictionEnzyme {
	var outputs []wtype.RestrictionEnzyme

	scanner := bufio.NewScanner(rebaseRh)
	// scanner.Split(bufio.ScanLines)
	name := ""
	prototype := ""
	recognitionseq := ""
	methylationsite := ""
	commercialsource := ""
	refs := ""
	//var data bytes.Buffer

	// Loop over the letters in inputString
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if len(line) == 0 {
			continue
		}

		// line := scanner.Text()

		if line[0] == '<' && line[1] == '1' {
			if name != "" {
				outputs = append(outputs, Build_rebase(name, prototype, recognitionseq, methylationsite, commercialsource, refs))

				recognitionseq = ""
				methylationsite = ""
				prototype = ""
			}

			name = line[3:]
		}
		if line[0] == '<' && line[1] == '2' {
			prototype = line[3:]
		}
		if line[0] == '<' && line[1] == '3' {
			recognitionseq = line[3:]
		}
		if line[0] == '<' && line[1] == '4' {
			methylationsite = line[3:]
		}
		if line[0] == '<' && line[1] == '5' {
			commercialsource = line[3:]
		}
		if line[0] == '<' && line[1] == '6' {
			refs = line[3:]
		}

	}

	outputs = append(outputs, Build_rebase(name, prototype, recognitionseq, methylationsite, commercialsource, refs))

	return outputs
}
