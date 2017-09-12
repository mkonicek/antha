// package align allows aligning Antha sequences using the biogo implementation of the
// Needleman-Wunsch and Smith-Waterman alignment algorithms
package align

import (
	"fmt"

	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/biogo/biogo/align"
	"github.com/biogo/biogo/alphabet"
	"github.com/biogo/biogo/feat"
	"github.com/biogo/biogo/seq/linear"
)

// ScoringMatrix implements the align.Aligner interface of the biogo/align package
// an align.Aligner aligns the sequence data of two type-matching Slicers,
// returning an ordered slice of features describing matching and mismatching segments.
// The sequences to be aligned must have a valid gap letter in the first position of their alphabet;
// the alphabets {DNA,RNA}{gapped,redundant}
// and Protein provided by the biogo/alphabet package satisfy this.
type ScoringMatrix interface {
	Align(reference, query align.AlphabetSlicer) ([]feat.Pair, error)
}

var (
	// Fitted is the linear gap penalty fitted Needleman-Wunsch aligner type.
	//		   Query letter
	//  	 -	 A	 C	 G	 T
	// -	 0	-5	-5	-5	-5
	// A	-5	10	-3	-1	-4
	// C	-5	-3	 9	-5	 0
	// G	-5	-1	-5	 7	-3
	// T	-5	-4	 0	-3	 8
	Fitted = align.Fitted{
		{0, -5, -5, -5, -5},
		{-5, 10, -3, -1, -4},
		{-5, -3, 9, -5, 0},
		{-5, -1, -5, 7, -3},
		{-5, -4, 0, -3, 8},
	}

	// FittedAffine is the affine gap penalty fitted Needleman-Wunsch aligner type.
	//		   Query letter
	//  	 -	 A	 C	 G	 T
	// -	 0	-1	-1	-1	-1
	// A	-1	 1	-1	-1	-1
	// C	-1	-1	 1	-1	-1
	// G	-1	-1	-1	 1	-1
	// T	-1	-1	-1	-1	 1
	//
	// Gap open: -5
	FittedAffine = align.FittedAffine{
		Matrix: align.Linear{
			{0, -1, -1, -1, -1},
			{-1, 1, -1, -1, -1},
			{-1, -1, 1, -1, -1},
			{-1, -1, -1, 1, -1},
			{-1, -1, -1, -1, 1},
		},
		GapOpen: -5,
	}

	// NW is the linear gap penalty Needleman-Wunsch aligner type.
	//		   Query letter
	//  	 -	 A	 C	 G	 T
	// -	 0	-5	-5	-5	-5
	// A	-5	10	-3	-1	-4
	// C	-5	-3	 9	-5	 0
	// G	-5	-1	-5	 7	-3
	// T	-5	-4	 0	-3	 8
	NW = align.NW{
		{0, -5, -5, -5, -5},
		{-5, 10, -3, -1, -4},
		{-5, -3, 9, -5, 0},
		{-5, -1, -5, 7, -3},
		{-5, -4, 0, -3, 8},
	}

	// NWAffine is the affine gap penalty Needleman-Wunsch aligner type.
	//		   Query letter
	//  	 -	 A	 C	 G	 T
	// -	 0	-1	-1	-1	-1
	// A	-1	 1	-1	-1	-1
	// C	-1	-1	 1	-1	-1
	// G	-1	-1	-1	 1	-1
	// T	-1	-1	-1	-1	 1
	//
	// Gap open: -5
	NWAffine = align.NWAffine{
		Matrix: align.Linear{
			{0, -1, -1, -1, -1},
			{-1, 1, -1, -1, -1},
			{-1, -1, 1, -1, -1},
			{-1, -1, -1, 1, -1},
			{-1, -1, -1, -1, 1},
		},
		GapOpen: -5,
	}

	// SW is the Smith-Waterman aligner type. Matrix is a square scoring matrix with the last column and last row specifying gap penalties. Currently gap opening is not considered.
	// w(gap) = -1
	// w(match) = +2
	// w(mismatch) = -1
	SW1 = align.SW{
		{0, -1, -1, -1, -1},
		{-1, 2, -1, -1, -1},
		{-1, -1, 2, -1, -1},
		{-1, -1, -1, 2, -1},
		{-1, -1, -1, -1, 2},
	}

	// w(gap) = 0
	// w(match) = +2
	// w(mismatch) = -1
	SW2 = align.SW{
		{0, 0, 0, 0, 0},
		{0, 2, -1, -1, -1},
		{0, -1, 2, -1, -1},
		{0, -1, -1, 2, -1},
		{0, -1, -1, -1, 2},
	}

	// SWAffine is the affine gap penalty Smith-Waterman aligner type.
	//		   Query letter
	//  	 -	 A	 C	 G	 T
	// -	 0	-1	-1	-1	-1
	// A	-1	 1	-1	-1	-1
	// C	-1	-1	 1	-1	-1
	// G	-1	-1	-1	 1	-1
	// T	-1	-1	-1	-1	 1
	//
	// Gap open: -5
	SWAffine = align.SWAffine{
		Matrix: align.Linear{
			{0, -1, -1, -1, -1},
			{-1, 1, -1, -1, -1},
			{-1, -1, 1, -1, -1},
			{-1, -1, -1, 1, -1},
			{-1, -1, -1, -1, 1},
		},
		GapOpen: -5,
	}
)

// Align aligns two sequences using a specified scoring algorithm.
// It returns an alignment description or an error if the scoring matrix is not square, or the sequence data types or alphabets do not match.
// algorithms available are:
// Fitted: a modified Needleman-Wunsch algorithm which finds a local region of the reference with high similarity to the query.
// FittedAffine: a modified Needleman-Wunsch algorithm which finds a local region of the reference with high similarity to the query.
// NW: the Needleman-Wunsch algorithm
// NWAffine: the affine gap penalty Needleman-Wunsch algorithm
// SW1 and SW2: the Smith-Waterman algorithm
// SWAffine: the affine gap penalty Smith-Waterman
func Align(seq1, seq2 wtype.DNASequence, alignMentMatrix ScoringMatrix) (alignment string, err error) {
	fsa := &linear.Seq{Seq: alphabet.BytesToLetters([]byte(seq1.Sequence()))}
	fsa.Alpha = alphabet.DNAgapped
	fsb := &linear.Seq{Seq: alphabet.BytesToLetters([]byte(seq2.Sequence()))}
	fsb.Alpha = alphabet.DNAgapped

	aln, err := alignMentMatrix.Align(fsa, fsb)
	if err == nil {
		fmt.Printf("%s\n", aln)
		fa := align.Format(fsa, fsb, aln, '-')
		alignment = fmt.Sprintf("%s\n%s\n", fa[0], fa[1])
	}
	return
}
