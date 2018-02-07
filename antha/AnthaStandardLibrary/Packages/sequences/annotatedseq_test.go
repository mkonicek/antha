// antha/AnthaStandardLibrary/Packages/enzymes/Annotatedseq.go: Part of the Antha language
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

// Package sequences is for interacting with and manipulating biological sequences; in extension to methods available in wtype
package sequences

import (
	"testing"
)

type genbanktest struct {
	featurename  string
	seq          string
	sequencetype string
	class        string
	reverse      string

	// results
	startposition int
	endposition   int
	aaseq         string
}

var genbanktests = []genbanktest{
	{featurename: "tetA",
		seq:           "atgccggtactgccgggcctcttgcgggatatcgtccattccgacagcatcgccagtcactatggcgtgctgctagcgctatatgcgttgatgcaatttctatgcgcacccgttctcggagcactgtccgaccgctttggccgccgcccagtcctgctcgcttcgctacttggagccactatcgactacgcgatcatggcgaccacacccgtcctgtggattctctacgccggacgcatcgtggccggcatcaccggcgccacaggtgcggttgctggcgcctatatcgccgacatcaccgatggggaagatcgggctcgccacttcgggctcatgagcgcttgtttcggcgtgggtatggtggcaggccccgtggccgggggactgttgggcgccatctccctgcacgcaccattccttgcggcggcggtgctcaacggcctcaacctactactgggctgcttcctaatgcaggagtcgcataagggagagcgccgtccgatgcccttgagagccttcaacccagtcagctccttccggtgggcgcggggcatgaccattgtggccgcacttatgactgtcttctttatcatgcaactcgtaggacaggtgccggcagcgctctgggtcattttcggcgaggaccgctttcgctggagcgcgacgatgatcggcctgtcgcttgcggtattcggaatcttgcacgccctcgctcaagccttcgtcactggtcccgccaccaaacgtttcggcgagaagcaggccattatcgccggcatggcggccgacgcgctgggctacgtcttgctggcgttcgcgacgcgaggctggatggccttccccattatgattcttctcgcttccggcggcatcgggatgcccgcgttgcaggccatgctgtccaggcaggtagatgacgaccatcagggacagcttcaaggatcgctcgcggctcttaccagcctaacttcgatcactggaccgctgatcgtcacggcgatttatgccgcctcggcgagcacatggaacgggttggcatggattgtaggcgccgccctataccttgtctgcctccccgcgttgcgtcgcggtgcatggagccgggccacctcgacctga",
		sequencetype:  "DNA",
		class:         "orf",
		reverse:       "",
		startposition: 1,
		endposition:   1119,
		aaseq:         "MPVLPGLLRDIVHSDSIASHYGVLLALYALMQFLCAPVLGALSDRFGRRPVLLASLLGATIDYAIMATTPVLWILYAGRIVAGITGATGAVAGAYIADITDGEDRARHFGLMSACFGVGMVAGPVAGGLLGAISLHAPFLAAAVLNGLNLLLGCFLMQESHKGERRPMPLRAFNPVSSFRWARGMTIVAALMTVFFIMQLVGQVPAALWVIFGEDRFRWSATMIGLSLAVFGILHALAQAFVTGPATKRFGEKQAIIAGMAADALGYVLLAFATRGWMAFPIMILLASGGIGMPALQAMLSRQVDDDHQGQLQGSLAALTSLTSITGPLIVTAIYAASASTWNGLAWIVGAALYLVCLPALRRGAWSRATST*",
	},
}

func TestMakeFeature(t *testing.T) {
	for _, test := range genbanktests {
		result := MakeFeature(test.featurename, test.seq, test.startposition, test.endposition, test.sequencetype, test.class, test.reverse)

		if result.StartPosition != test.startposition {
			t.Error(
				"For", test.featurename, "/n",
				"expected", test.startposition, "\n",
				"got", result.StartPosition, "\n",
			)
		}
		if result.EndPosition != test.endposition {
			t.Error(
				"For", test.featurename, "/n",
				"expected", test.endposition, "\n",
				"got", result.EndPosition, "\n",
			)
		}
		if result.Protseq != test.aaseq {
			t.Error(
				"For", test.featurename, "/n",
				"expected", test.aaseq, "\n",
				"got", result.Protseq, "\n",
			)
		}
	}
}
