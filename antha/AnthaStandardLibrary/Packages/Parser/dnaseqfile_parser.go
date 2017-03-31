// antha/AnthaStandardLibrary/Packages/Parser/fasta_parser.go: Part of the Antha language
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
	"fmt"
	//"os"
	"path/filepath"
	"strings"

	"github.com/antha-lang/antha/antha/anthalib/wtype"
)

func DNAFiletoDNASequence(filename string, plasmid bool) (sequences []wtype.DNASequence, err error) {

	sequences = make([]wtype.DNASequence, 0)
	var seqs []wtype.DNASequence
	var seq wtype.DNASequence
	switch fn := filename; {
	case strings.HasSuffix(fn, ".gdx"):

		seqs, err = GDXtoDNASequence(filename)
		for _, seq := range seqs {
			sequences = append(sequences, seq)
		}
	case strings.HasSuffix(fn, ".fasta"):
		seqs, err = FastatoDNASequences(filename)
		for _, seq := range seqs {
			sequences = append(sequences, seq)
		}
	case strings.HasSuffix(fn, ".gb"):

		seq, err = GenbanktoFeaturelessDNASequence(filename)

		sequences = append(sequences, seq)
	default:
		err = fmt.Errorf("non valid sequence file")
	}

	if err != nil {
		return seqs, err
	}
	//}
	return
}

// Creates a DNASequence from a sequence file of format: .gdx .fasta .gb
func DNAFileToDNASequenceBinary(filename wtype.File, plasmid bool) (sequences []wtype.DNASequence, err error) {

	data, err := filename.ReadAll()
	if err != nil {
		fmt.Errorf("Cannot parse file. File is empty.")
	}
	sequences = make([]wtype.DNASequence, 0)
	var seqs []wtype.DNASequence
	var seq wtype.DNASequence

	switch fn := filename.Name; {
	case filepath.Ext(fn) == ".gdx":
		seqs, err = GDXtoDNASequenceBinary(data)
		for _, seq := range seqs {
			sequences = append(sequences, seq)
		}
	case filepath.Ext(fn) == ".fasta":
		seqs, err = FASTAtoDNASeqsBinary(data)
		for _, seq := range seqs {
			sequences = append(sequences, seq)
		}
	case filepath.Ext(fn) == ".gb":
		seq, err = GenbanktoFeaturelessDNASequenceBinary(data)
		sequences = append(sequences, seq)
	default:
		err = fmt.Errorf("non valid sequence file")
	}

	if err != nil {
		return seqs, err
	}
	//}
	return
}
