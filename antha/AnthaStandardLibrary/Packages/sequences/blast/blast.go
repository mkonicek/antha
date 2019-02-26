// Part of the Antha language
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

// Package for performing blast queries
package blast

import (
	"fmt"

	"strconv"
	"strings"
	"time"

	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/sequences/biogo/ncbi/blast"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/text"
	"github.com/mgutz/ansi"
)

// package for interacting with the ncbi BLAST service

var (
	email     = "no-reply@antha-lang.com"
	tool      = "blast-biogo-antha"
	putparams = blast.PutParameters{Program: "blastn", Megablast: true, Database: "nr"}
	getparams blast.GetParameters
	//query     = "X14032.1"
	//query     = "MSFSNYKVIAMPVLVANFVLGAATAWANENYPAKSAGYNQGDWVASFNFSKVYVGEELGDLNVGGGALPNADVSIGNDTTLTFDIAYFVSSNIAVDFFVGVPARAKFQGEKSISSLGRVSEVDYGPAILSLQYHYDSFERLYPYVGVGVGRVLFFDKTDGALSSFDIKDKWAPAFQVGLRYDLGNSWMLNSDVRYIPFKTDVTGTLGPVPVSTKIEVDPFILSLGASYVF"
	//query   = "atgagtttttctaattataaagtaatcgcgatgccggtgttggttgctaattttgttttgggggcggccactgcatgggcgaatgaaaattatccggcgaaatctgctggctataatcagggtgactgggtcgctagcttcaatttttctaaggtctatgtgggtgaggagcttggcgatctaaatgttggagggggggctttgccaaatgctgatgtaagtattggtaatgatacaacacttacgtttgatatcgcctattttgttagctcaaatatagcggtggatttttttgttggggtgccagctagggctaaatttcaaggtgagaaatcaatctcctcgctgggaagagtcagtgaagttgattacggccctgcaattctttcgcttcaatatcattacgatagctttgagcgactttatccatatgttggggttggtgttggtcgggtgctattttttgataaaaccgacggtgctttgagttcgtttgatattaaggataaatgggcgcctgcttttcaggttggccttagatatgaccttggtaactcatggatgctaaattcagatgtgcgttatattcctttcaaaacggacgtcacaggtactcttggcccggttcctgtttctactaaaattgaggttgatcctttcattctcagtcttggtgcgtcatatgttttctaa"
	retries = 5
)

func RerunRIDstring(rid string) (o *blast.Output, err error) {
	r := blast.NewRid(rid)

	if r != nil {
		fmt.Println("RID=", r.String())

		//var o *Output
		for k := 0; k < retries; k++ {
			var s *blast.SearchInfo
			s, err = r.SearchInfo(tool, email)
			fmt.Println(s.Status)

			fmt.Println("hits?", s.HaveHits)
			if s.HaveHits {
				o, err = r.GetOutput(&getparams, tool, email)
				return
			}

			if err == nil {
				break
			}

		}
	} else {
		err = fmt.Errorf("r == nil")
	}

	return
}

func RerunRID(r *blast.Rid) (o *blast.Output, err error) {

	if r != nil {
		fmt.Println("RID=", r.String())

		//var o *Output
		for k := 0; k < retries; k++ {
			var s *blast.SearchInfo
			s, err = r.SearchInfo(tool, email)
			fmt.Println(s.Status)

			fmt.Println("hits?", s.HaveHits)
			if s.HaveHits {
				o, err = r.GetOutput(&getparams, tool, email)
				return
			}
			if err == nil {
				break
			}

		}
	} else {
		err = fmt.Errorf("r == nil")
	}

	return
}

func HitSummary(hits []blast.Hit, topnumberofhits int, topnumberofhsps int) (summary string, err error) {

	summaryarray := make([]string, 0)

	if len(hits) != 0 {

		summaryarray = append(summaryarray, fmt.Sprintln(ansi.Color("Hits:", "green"), len(hits)))

		for i, hit := range hits {

			if i >= topnumberofhits {
				summary = strings.Join(summaryarray, "; ")
				return
			}

			for j := range hit.Hsps {

				if j >= topnumberofhsps {
					break
				}

				seqlength := hits[i].Len

				hspidentityfloat := float64(*hits[i].Hsps[j].HspIdentity)
				querylengthfloat := float64(len(hits[i].Hsps[j].QuerySeq))
				subjectseqfloat := float64(len(hits[i].Hsps[j].SubjectSeq))

				identityfloat := (hspidentityfloat / querylengthfloat) * 100
				coveragefloat := (querylengthfloat / subjectseqfloat) * 100
				identity := strconv.FormatFloat(identityfloat, 'G', -1, 64) + "%"
				coverage := strconv.FormatFloat(coveragefloat, 'G', -1, 64) + "%"

				hitsum := fmt.Sprintln(ansi.Color("Hit:", "blue"), i+1,
					//	Printfield(hits[0].Id),
					text.Sprint("HspIdentity: ", strconv.Itoa(*hits[i].Hsps[j].HspIdentity)),
					text.Sprint("queryLen: ", len(hits[i].Hsps[j].QuerySeq)),

					text.Sprint("queryFrom: ", hits[i].Hsps[j].QueryFrom),
					text.Sprint("queryTo: ", hits[i].Hsps[j].QueryTo),
					text.Sprint("subjectLen: ", len(hits[i].Hsps[j].SubjectSeq)),
					text.Sprint("HitFrom: ", hits[i].Hsps[j].HitFrom),
					text.Sprint("HitTo: ", hits[i].Hsps[j].HitTo),
					text.Sprint("alignLen: ", *hits[i].Hsps[j].AlignLen),
					text.Sprint("Identity: ", identity),
					text.Sprint("coverage: ", coverage),

					ansi.Color("Sequence length:", "red"), seqlength,
					ansi.Color("high scoring pairs for top match:", "red"), len(hits[0].Hsps),
					ansi.Color("Id:", "red"), hits[i].Id,
					ansi.Color("Definition:", "red"), hits[i].Def,
					ansi.Color("Accession:", "red"), hits[i].Accession,

					ansi.Color("Bitscore", "red"), hits[i].Hsps[j].BitScore,
					ansi.Color("Score", "red"), hits[i].Hsps[j].Score,
					ansi.Color("EValue", "red"), hits[i].Hsps[j].EValue)

				summaryarray = append(summaryarray, hitsum)

			}

		}

	} else {
		summary = "No hits!"
		err = fmt.Errorf(summary)
	}
	return
}

func FindBestHit(hits []blast.Hit) (besthit blast.Hit, identity float64, coverage float64, besthitsummary string, err error) {

	var besthitnumber int
	highestidentity := 0.0
	highestcoverage := 0.0
	longestquery := 0.0

	if len(hits) != 0 {

		for i, hit := range hits {

			for j := range hit.Hsps {

				seqlength := hits[i].Len

				hspidentityfloat := float64(*hits[i].Hsps[j].HspIdentity)
				querylengthfloat := float64(len(hits[i].Hsps[j].QuerySeq))
				subjectseqfloat := float64(len(hits[i].Hsps[j].SubjectSeq))

				identityfloat := (hspidentityfloat / querylengthfloat) * 100
				coveragefloat := (querylengthfloat / subjectseqfloat) * 100

				if coveragefloat > highestcoverage && identityfloat > highestidentity && querylengthfloat > longestquery {
					besthitnumber = i
					highestcoverage = coveragefloat
					highestidentity = identityfloat
					identity = identityfloat
					coverage = coveragefloat

					// prepare summary
					identitystr := strconv.FormatFloat(identityfloat, 'G', -1, 64) + "%"
					coveragestr := strconv.FormatFloat(coveragefloat, 'G', -1, 64) + "%"
					besthitsummary = fmt.Sprintln(ansi.Color("Hit:", "blue"), i+1,
						//	Printfield(hits[0].Id),
						text.Sprint("HspIdentity: ", strconv.Itoa(*hits[i].Hsps[j].HspIdentity)),
						text.Sprint("queryLen: ", len(hits[i].Hsps[j].QuerySeq)),
						text.Sprint("queryFrom: ", hits[i].Hsps[j].QueryFrom),
						text.Sprint("queryTo: ", hits[i].Hsps[j].QueryTo),
						text.Sprint("subjectLen: ", len(hits[i].Hsps[j].SubjectSeq)),
						text.Sprint("HitFrom: ", hits[i].Hsps[j].HitFrom),
						text.Sprint("HitTo: ", hits[i].Hsps[j].HitTo),
						text.Sprint("alignLen: ", *hits[i].Hsps[j].AlignLen),
						text.Sprint("Identity: ", identitystr),
						text.Sprint("coverage: ", coveragestr),
						ansi.Color("Sequence length:", "red"), seqlength,
						ansi.Color("high scoring pairs for top match:", "red"), len(hits[0].Hsps),
						ansi.Color("Id:", "red"), hits[i].Id,
						ansi.Color("Definition:", "red"), hits[i].Def,
						ansi.Color("Accession:", "red"), hits[i].Accession,
						ansi.Color("Bitscore", "red"), hits[i].Hsps[j].BitScore,
						ansi.Color("Score", "red"), hits[i].Hsps[j].Score,
						ansi.Color("EValue", "red"), hits[i].Hsps[j].EValue)
				}

			}

		}
		besthit = hits[besthitnumber]
	} else {
		besthitsummary = "No hits!"
		err = fmt.Errorf(besthitsummary)
	}
	return
}

func AllExactMatches(hits []blast.Hit) (exactmatches []blast.Hit, summary string, err error) {
	exactmatches = make([]blast.Hit, 0)

	if len(hits) != 0 {
		for i, hit := range hits {

			for j := range hit.Hsps {
				hspidentityfloat := float64(*hits[i].Hsps[j].HspIdentity)
				querylengthfloat := float64(len(hits[i].Hsps[j].QuerySeq))

				identityfloat := (hspidentityfloat / querylengthfloat) * 100

				if identityfloat == 100 {
					exactmatches = append(exactmatches, hit)
				}
			}

		}

	} else {
		summary = "No hits!"
		err = fmt.Errorf(summary)
	}
	return
}

func MegaBlastP(query string) (hits []blast.Hit, err error) {

	putparams = blast.PutParameters{Program: "blastp", Megablast: true, Database: "nr"}
	o, err := SimpleBlast(query)
	if err != nil {
		return
	}
	hits, err = Hits(o)
	if err != nil {
		return
	}

	return
}

func MegaBlastN(query string) (hits []blast.Hit, err error) {
	putparams = blast.PutParameters{Program: "blastn", Megablast: true, Database: "nr"}

	o, err := SimpleBlast(query)
	if err != nil {
		return
	}
	hits, err = Hits(o)
	if err != nil {
		return
	}
	return
}

// SimpleBlast performs a blast call with pre-set query parameters
func SimpleBlast(query string) (*blast.Output, error) {
	return BLAST(query, &putparams)
}

// BLAST performs a blast call for specified query parameters
// For parameters see https://github.com/biogo/ncbi/blob/master/blast/blast.go
// For documentation on settings see https://ncbi.github.io/blast-cloud/dev/api.html
// Commonly used parameters include
/*
   Program       string      BLAST program to use (e.g. "blastp", "blastn")
   Database      string      Target database name (e.g. "nr", "refseq_rna", "pdb")
   EntrezQuery   string      Entrez results filter (e.g. "Homo sapiens[organism]")
   Expect        *float64    Expect threshold
   HitListSize   int         Number of target sequences to return
*/
func BLAST(query string, p *blast.PutParameters) (o *blast.Output, err error) {
	r, err := blast.Put(query, p, tool, email)
	fmt.Println("RID=", r.String())
	fmt.Println("Submitting request to BLAST server, please wait")
	//var o *Output
	for k := 0; k < retries; k++ {
		var s *blast.SearchInfo
		s, err = r.SearchInfo(tool, email)
		fmt.Println(s.Status)

		fmt.Println("hits?", s.HaveHits)
		if s.HaveHits {
			o, err = r.GetOutput(&getparams, tool, email)
			return
		} else if strings.Contains(s.Status, "WAITING") {
			for {
				if strings.Contains(s.Status, "WAITING") {
					fmt.Println("waiting 1 min to rerun RID:", r.String())
					time.Sleep(1 * time.Minute)
					s, err = r.SearchInfo(tool, email)
					if err != nil {
						return
					}

					o, err = RerunRID(r)

				} else {
					return
				}
			}
		}
	}

	return
}

func Hits(o *blast.Output) (hits []blast.Hit, err error) {
	if o == nil {
		err = fmt.Errorf("output == nil")
		return
	}
	if len(o.Iterations) == 0 {
		err = fmt.Errorf("len(output.Iterations) == 0")
		return
	}
	if len(o.Iterations[0].Hits) == 0 {
		err = fmt.Errorf("len(output.Iterations[0].Hits) == 0")
		return
	}
	hits = o.Iterations[0].Hits

	return
}

/*
func BestHit(hits []Hit) (besthit Hit) {

}
*/
