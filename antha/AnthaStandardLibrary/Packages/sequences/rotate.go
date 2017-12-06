// antha/AnthaStandardLibrary/Packages/enzymes/Translation.go: Part of the Antha language
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

import "github.com/antha-lang/antha/antha/anthalib/wtype"

// Rotate will rotate the sequence by the number of characters specified by rotateBy.
// If reverse is true the sequence will be rotated in the reverse direction.
func Rotate(seq wtype.DNASequence, rotateBy int, reverse bool) (rotatedSeq wtype.DNASequence) {

	var tempSeq = seq.Seq
	var rotatedSeqStr string

	if !reverse {
		rotatedSeqStr += tempSeq[rotateBy:]
		rotatedSeqStr += tempSeq[:rotateBy]
	} else {
		rotatedSeqStr += tempSeq[len(tempSeq)-rotateBy:]
		rotatedSeqStr += tempSeq[:len(tempSeq)-rotateBy]
	}
	originalFeatures := seq.Features

	rotatedSeq = seq
	rotatedSeq.Seq = rotatedSeqStr

	SetFeatures(&rotatedSeq, originalFeatures)

	return rotatedSeq
}
