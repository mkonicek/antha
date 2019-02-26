// anthalib//wutil/biogo_helpers.go: Part of the Antha language
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

package wutil

import (
	"bufio"
	"os"

	"github.com/biogo/biogo/io/seqio/fasta"
	"github.com/biogo/biogo/seq"
	"github.com/biogo/biogo/seq/linear"
)

func SeqToBioseq(s seq.Sequence) string {
	ret := ""

	for i := s.Start(); i < s.End(); i++ {
		ret += string(s.At(i).L)
	}

	return ret
}

func ReadFastaSeqs(fn string) []seq.Sequence {
	ret := make([]seq.Sequence, 0, 1)

	f, e := os.Open(fn)

	if e != nil {
		panic(e)
	}

	r := bufio.NewReader(f)
	reader := fasta.NewReader(r, linear.NewSeq("", nil, nil))

	var s seq.Sequence

	for {
		s, e = reader.Read()
		if e != nil {
			break
		}

		ret = append(ret, s)
	}

	return ret
}
