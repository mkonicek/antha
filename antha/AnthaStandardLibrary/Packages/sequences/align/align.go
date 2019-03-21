// Package align allows aligning Antha sequences using the biogo implementation of the
// Needleman-Wunsch and Smith-Waterman alignment algorithms
package align

import (
	"fmt"
	"strings"

	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/sequences"

	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/biogo/biogo/align"
	"github.com/biogo/biogo/alphabet"
	"github.com/biogo/biogo/feat"
	"github.com/biogo/biogo/seq"
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

// Result stores the full results of an alignment of a query against a template sequence, including the algorithm used.
type Result struct {
	Template  wtype.BioSequence
	Query     wtype.BioSequence
	Algorithm ScoringMatrix
	Alignment Alignment
}

// String prints alignment result in form of two aligned sequence strings printed on parallel lines.
func (r Result) String() string {
	return fmt.Sprintf("%s\n%s\n", r.Alignment.TemplateResult, r.Alignment.QueryResult)
}

// Identity returns the percentage of matching nucleotides of query in the template sequence
// a value between 0 and 1 is returned
// 1 = 100%; 0 = 0%
func (r Result) Identity() float64 {
	var totalCount int
	var matchCount int

	for i := range r.Alignment.QueryResult {
		if r.Alignment.QueryResult[i] == r.Alignment.TemplateResult[i] {
			matchCount++
		}
		totalCount++
	}
	return float64(matchCount) / float64(totalCount)
}

// Gaps returns the number of gaps in the aligned query sequence result
func (r Result) Gaps() int {
	var gapCount int
	for i, letter := range r.Alignment.QueryResult {
		if letter == GAP || rune(r.Alignment.TemplateResult[i]) == GAP {
			gapCount++
		}
	}
	return gapCount
}

// Matches returns the number of matched nucleotides between the aligned query
// sequence and aligned template sequence.
func (r Result) Matches() int {
	var matchCount int
	for i := range r.Alignment.QueryResult {
		if r.Alignment.QueryResult[i] == r.Alignment.TemplateResult[i] {
			matchCount++
		}
	}
	return matchCount
}

// Mismatches returns the number of mismatched nucleotides between the aligned query
// sequence and aligned template sequence.
func (r Result) Mismatches() int {
	var mismatchCount int
	for i := range r.Alignment.QueryResult {
		if isMismatch(rune(r.Alignment.QueryResult[i]), rune(r.Alignment.TemplateResult[i])) {
			mismatchCount++
		}
	}
	return mismatchCount
}

// Coverage returns the percentage of matching nucleotides of alignment to the template sequence
// a value between 0 and 1 is returned
// 1 = 100%; 0 = 0%
func (r Result) Coverage() float64 {
	return float64(len(r.Alignment.TemplateResult)) / float64(len(r.Template.Sequence()))
	//return float64(len(r.Alignment.QueryResult)) / float64(len(r.Query.Sequence()))
}

// LongestContinuousSequence returns the longest unbroken chain of matches as a dna sequence
func (r Result) LongestContinuousSequence() wtype.DNASequence {
	var longest []string
	var seq []string

	for i, char := range r.Alignment.QueryResult {
		if len(seq) > len(longest) {
			longest = seq
		}
		if !isMismatch(char, rune(r.Alignment.TemplateResult[i])) && !isGap(char) && !isGap(rune(r.Alignment.TemplateResult[i])) {
			seq = append(seq, string(char))
		} else {
			// reset
			seq = make([]string, 0)
		}
	}
	if len(seq) > len(longest) {
		longest = seq
	}

	var newSeq string

	if r.Alignment.TemplateFrame() == -1 {
		//if looksLikeReverseSequence(r.Alignment.TemplatePositions) {
		newSeq = wtype.RevComp(strings.Join(longest, ""))
	} else {
		newSeq = strings.Join(longest, "")
	}

	return wtype.DNASequence{Nm: r.Query.Name() + "_alignment", Seq: newSeq}
}

// Positions returns a SearchResult detailing the positions in the template sequence
// of the longest continuous matching sequence from the alignment.
func (r Result) Positions() (result sequences.SearchResult) {
	templateSeq, ok := r.Template.(*wtype.DNASequence)
	if !ok {
		err := fmt.Errorf("Cannot cast template into DNASequence, alignment currently only supports DNASequence alignment")
		panic(err)
	}
	querySeq := r.LongestContinuousSequence()
	return sequences.FindAll(templateSeq, &querySeq)
}

// Score returns the alignment score
func (r Result) Score() int {
	return r.Alignment.Score
}

// Alignment stores the string result of an alignment of a query sequence against a template
// The original RawAlignments are also included
type Alignment struct {
	TemplateResult    string
	QueryResult       string
	Raw               []RawAlignment
	TemplatePositions []int
	QueryPositions    []int
	Score             int
}

// RawAlignment contains the positions aligned between the template and query sequences
type RawAlignment struct {
	TemplateAlignment Position
	QueryAlignment    Position
}

// Position contains the start, end and length of an alignment in a specified sequence
type Position struct {
	Start  int
	End    int
	Length int
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

	// SW1 is the Smith-Waterman aligner type. Matrix is a square scoring matrix with the last column and last row specifying gap penalties. Currently gap opening is not considered.
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

	// SW2 is the Smith-Waterman aligner type. Matrix is a square scoring matrix with the last column and last row specifying gap penalties. Currently gap opening is not considered.
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

// Algorithms provides a map to lookup ScoringMatrix algorithms based on names.
// Algorithms available:
// Fitted: a modified Needleman-Wunsch algorithm which finds a local region of the reference with high similarity to the query.
// FittedAffine: a modified Needleman-Wunsch algorithm which finds a local region of the reference with high similarity to the query.
// NW: the Needleman-Wunsch algorithm
// NWAffine: the affine gap penalty Needleman-Wunsch algorithm
// SW1 and SW2: the Smith-Waterman algorithm
var Algorithms = map[string]ScoringMatrix{
	"Fitted":       Fitted,
	"FittedAffine": FittedAffine,
	"NW":           NW,
	"NWAffine":     NWAffine,
	"SW1":          SW1,
	"SW2":          SW2,
	"SWAffine":     SWAffine,
}

// DNA aligns two DNA sequences using a specified scoring algorithm.
// It returns an alignment description or an error if the scoring matrix is not square, or the sequence data types or alphabets do not match.
// algorithms available are:
// Fitted: a modified Needleman-Wunsch algorithm which finds a local region of the reference with high similarity to the query.
// FittedAffine: a modified Needleman-Wunsch algorithm which finds a local region of the reference with high similarity to the query.
// NW: the Needleman-Wunsch algorithm
// NWAffine: the affine gap penalty Needleman-Wunsch algorithm
// SW1 and SW2: the Smith-Waterman algorithm
// SWAffine: the affine gap penalty Smith-Waterman
// Alignment of the reverse complement of the query sequence will also be attempted and if the number of matches is higher the reverse alignment is returned.
// In the resulting alignment, mismatches are represented by lower case letters, gaps represented by the GAP character "-".
func DNA(template, query wtype.DNASequence, alignmentMatrix ScoringMatrix) (alignment Result, err error) {

	fwdResult, err := DNAFwd(template, query, alignmentMatrix)

	if err != nil {
		return
	}

	revResult, err := DNARev(template, query, alignmentMatrix)

	if err != nil {
		return fwdResult, fmt.Errorf(fmt.Sprintf("Error with aligning reverse complement of query sequence %s: %s", query.Nm, err.Error()))
	}

	if len(query.Seq) > len(template.Seq) {
		if err == nil {
			err = fmt.Errorf("query sequence is larger than template, this may result in an unusual alignment")
		}
	}

	if revResult.Matches() > fwdResult.Matches() {
		return revResult, err
	}

	return fwdResult, err
}

// DNAFwd returns an alignment of a query sequence to a template sequence in the forward frame of the template, using a specified scoring algorithm
func DNAFwd(template, query wtype.DNASequence, alignmentMatrix ScoringMatrix) (Result, error) {
	return dnaFWDAlignment(template, query, alignmentMatrix)
}

// DNARev returns the alignment of a query sequence to a template sequence in the reverse frame of the template, using a specified scoring algorithm
func DNARev(template, query wtype.DNASequence, alignmentMatrix ScoringMatrix) (alignment Result, err error) {

	revQuery := query

	revQuery.Seq = wtype.RevComp(query.Seq)

	alignment, err = dnaFWDAlignment(template, revQuery, alignmentMatrix)

	if err != nil {
		return
	}

	alignment = correctForRevComp(alignment)

	return
}

type scorer interface {
	Score() int
}

func dnaFWDAlignment(template, query wtype.DNASequence, alignmentMatrix ScoringMatrix) (alignment Result, err error) {

	if containsN(template) {
		err = fmt.Errorf("template sequence %s contains N values. Please replace these with - before running alignment", template.Name())
		return
	}

	if containsN(query) {
		err = fmt.Errorf("query sequence %s contains N values. Please replace these with - before running alignment", query.Name())
		return
	}

	var errs []string

	seq1, seq2 := template, query

	fsa := &linear.Seq{Seq: alphabet.BytesToLetters([]byte(seq1.Sequence()))}
	fsa.Alpha = alphabet.DNAgapped
	fsb := &linear.Seq{Seq: alphabet.BytesToLetters([]byte(seq2.Sequence()))}
	fsb.Alpha = alphabet.DNAgapped

	aln, err := alignmentMatrix.Align(fsa, fsb)
	if err == nil {
		var rawAlignments []RawAlignment
		sumScores := 0
		for i := range aln {

			rawAlignemnt := RawAlignment{
				TemplateAlignment: Position{
					Start:  aln[i].Features()[0].Start(),
					End:    aln[i].Features()[0].End(),
					Length: aln[i].Features()[0].Len(),
				},
				QueryAlignment: Position{
					Start:  aln[i].Features()[1].Start(),
					End:    aln[i].Features()[1].End(),
					Length: aln[i].Features()[1].Len(),
				},
			}
			rawAlignments = append(rawAlignments, rawAlignemnt)
			sumScores += aln[i].(scorer).Score()

		}
		fa, positions := format(fsa, fsb, aln, '-')
		alignment = Result{
			Template:  &template,
			Query:     &query,
			Algorithm: alignmentMatrix,
			Alignment: formatMisMatches(Alignment{
				TemplateResult:    fmt.Sprint(fa[0]),
				QueryResult:       fmt.Sprint(fa[1]),
				Raw:               rawAlignments,
				TemplatePositions: positions[0],
				QueryPositions:    positions[1],
				Score:             sumScores,
			}),
		}
	} else {
		errs = append(errs, err.Error())
	}

	// if plasmid we must try again with rotating vector the length of query sequence
	// in case the alignment crosses the end and start of the plasmid sequence
	if template.Plasmid {
		var tempseq string
		rotationSize := len(query.Seq) - 1

		if rotationSize > len(template.Seq) {
			return
		}

		tempseq += template.Seq[rotationSize:]
		tempseq += template.Seq[:rotationSize]
		seq1.Seq = tempseq

		fsa = &linear.Seq{Seq: alphabet.BytesToLetters([]byte(seq1.Sequence()))}
		fsa.Alpha = alphabet.DNAgapped
		fsb = &linear.Seq{Seq: alphabet.BytesToLetters([]byte(seq2.Sequence()))}
		fsb.Alpha = alphabet.DNAgapped
		aln, err := alignmentMatrix.Align(fsa, fsb)
		if err == nil {
			var rawAlignments []RawAlignment
			sumScores := 0
			for i := range aln {

				rawAlignemnt := RawAlignment{
					TemplateAlignment: Position{
						Start:  aln[i].Features()[0].Start(),
						End:    aln[i].Features()[0].End(),
						Length: aln[i].Features()[0].Len(),
					},
					QueryAlignment: Position{
						Start:  aln[i].Features()[1].Start(),
						End:    aln[i].Features()[1].End(),
						Length: aln[i].Features()[1].Len(),
					},
				}
				rawAlignments = append(rawAlignments, rawAlignemnt)
				sumScores += aln[i].(scorer).Score()

			}
			fa, positions := format(fsa, fsb, aln, '-')

			// correct rotation for template positions
			for n := range positions[0] {

				var newPos int

				if positions[0][n]+rotationSize > len(template.Seq) {
					newPos = positions[0][n] + rotationSize - len(template.Seq)
				} else {
					newPos = positions[0][n] + rotationSize
				}
				positions[0][n] = newPos
			}

			rotatedAlignment := Result{
				Template:  &template,
				Query:     &query,
				Algorithm: alignmentMatrix,
				Alignment: formatMisMatches(Alignment{
					TemplateResult:    fmt.Sprint(fa[0]),
					QueryResult:       fmt.Sprint(fa[1]),
					Raw:               rawAlignments,
					TemplatePositions: positions[0],
					QueryPositions:    positions[1],
					Score:             sumScores,
				}),
			}
			if rotatedAlignment.Matches() > alignment.Matches() {
				alignment = rotatedAlignment
			}
		} else {
			errs = append(errs, err.Error())
		}
	}

	if len(errs) > 0 {
		err = fmt.Errorf("Errors: %s", strings.Join(errs, ";"))
	}

	return
}

// format returns a [2]alphabet.Slice representing the formatted alignment of a and b described by the
// list of feature pairs in f, with gap used to fill gaps in the alignment.
// variation based on align.Format from biogo library
func format(templateSeq, querySeq seq.Slicer, alignedPairs []feat.Pair, gap alphabet.Letter) ([2]alphabet.Slice, [2][]int) {
	var originalSeqs, alignedSeqs [2]alphabet.Slice
	var newPositions [2][]int
	for i, sequence := range [2]seq.Slicer{templateSeq, querySeq} {
		originalSeqs[i] = sequence.Slice()
		alignedSeqs[i] = originalSeqs[i].Make(0, 0)
	}
	for _, pair := range alignedPairs {
		features := pair.Features()
		for i := range alignedSeqs {
			if features[i].Len() == 0 {
				switch alignedSeqs[i].(type) {
				case alphabet.Letters:
					alignedSeqs[i] = alignedSeqs[i].Append(alphabet.Letters(gap.Repeat(features[1-i].Len())))
					for counter := features[1-i].Start(); counter < features[1-i].End(); counter++ {
						newPositions[i] = append(newPositions[i], features[i].Start()+1)
					}
				case alphabet.QLetters:
					alignedSeqs[i] = alignedSeqs[i].Append(alphabet.QLetters(alphabet.QLetter{L: gap}.Repeat(features[1-i].Len())))
					for counter := features[1-i].Start(); counter < features[1-i].End(); counter++ {
						newPositions[i] = append(newPositions[i], features[i].Start()+1)
					}
				}
			} else {
				alignedSeqs[i] = alignedSeqs[i].Append(originalSeqs[i].Slice(features[i].Start(), features[i].End()))
				for counter := features[i].Start(); counter < features[i].End(); counter++ {
					newPositions[i] = append(newPositions[i], counter+1)
				}
			}
		}
	}
	return alignedSeqs, newPositions
}

func containsN(seq wtype.DNASequence) bool {
	for _, letter := range seq.Seq {
		if strings.ToUpper(string(letter)) == "N" {
			return true
		}
	}
	return false
}

// changes mismatched nucleotides to lower case
func formatMisMatches(alignment Alignment) (formattedAlignment Alignment) {

	var formattedQuery []string
	var formattedTemplate []string

	for i := range alignment.QueryResult {
		var queryChar string
		var templateChar string
		if isMismatch(rune(alignment.QueryResult[i]), rune(alignment.TemplateResult[i])) {
			queryChar = strings.ToLower(string(alignment.QueryResult[i]))
			templateChar = strings.ToLower(string(alignment.TemplateResult[i]))
		} else {
			queryChar = strings.ToUpper(string(alignment.QueryResult[i]))
			templateChar = strings.ToUpper(string(alignment.TemplateResult[i]))
		}
		formattedQuery = append(formattedQuery, queryChar)
		formattedTemplate = append(formattedTemplate, templateChar)
	}
	formattedAlignment = alignment
	// replace alignment values
	formattedAlignment.TemplateResult = strings.Join(formattedTemplate, "")
	formattedAlignment.QueryResult = strings.Join(formattedQuery, "")

	return
}

// correctForRevComp corrects positions to be that of reverse complement following manual reverse complement of strand before alignment.
// also corrects the query sequence back to the original.
func correctForRevComp(alignment Result) (formattedAlignment Result) {

	correctPositions := make([]int, len(alignment.Alignment.TemplatePositions))

	for i := range alignment.Alignment.TemplatePositions {
		correctPositions[len(alignment.Alignment.TemplatePositions)-1-i] = alignment.Alignment.TemplatePositions[i]
	}

	formattedAlignment = alignment

	formattedAlignment.Alignment.TemplatePositions = correctPositions

	revQuery, ok := formattedAlignment.Query.(*wtype.DNASequence)

	if !ok {
		panic("cannot correct for reverse compliment if sequence is not DNA")
	}

	revQuery.Seq = wtype.RevComp(revQuery.Seq)

	formattedAlignment.Alignment.TemplateResult = wtype.RevComp(formattedAlignment.Alignment.TemplateResult)
	formattedAlignment.Alignment.QueryResult = wtype.RevComp(formattedAlignment.Alignment.QueryResult)

	formattedAlignment.Query = revQuery
	return
}

// GAP defines a standard character representing an alignment gap
const GAP rune = rune('-')

func isGap(character rune) bool {
	return character == GAP
}

func isMismatch(character1, character2 rune) bool {

	if isGap(character1) || isGap(character2) {
		return false
	}

	return !strings.EqualFold(string(character1), string(character2))
}
