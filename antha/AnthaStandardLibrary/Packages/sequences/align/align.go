// package align allows aligning Antha sequences using the biogo implementation of the
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
	if len(longest) == 0 {
		longest = seq
	}
	return wtype.DNASequence{Nm: r.Query.Name() + "_alignment", Seq: strings.Join(longest, "")}
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
	return sequences.FindSeq(templateSeq, &querySeq)
}

// Alignment stores the string result of an alignment of a query sequence against a template
type Alignment struct {
	TemplateResult string
	QueryResult    string
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

	// SW is the Smith-Waterman aligner type. Matrix is a square scoring matrix with the last column and last row specifying gap penalties. Currently gap opening is not considered.
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
var Algorithms map[string]ScoringMatrix = map[string]ScoringMatrix{
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
func DNA(seq1, seq2 wtype.DNASequence, alignMentMatrix ScoringMatrix) (alignment Result, err error) {

	fwdResult, err := dnaFWDAlignment(seq1, seq2, alignMentMatrix)

	if err != nil {
		return
	}

	revQuery := seq2

	revQuery.Seq = wtype.RevComp(seq2.Seq)

	revResult, err := dnaFWDAlignment(seq1, revQuery, alignMentMatrix)

	if err != nil {
		return fwdResult, fmt.Errorf(fmt.Sprintf("Error with aligning reverse complement of query sequence %s: %s", seq2.Nm, err.Error()))
	}

	if revResult.Matches() > fwdResult.Matches() {
		return revResult, nil
	}

	return fwdResult, nil
}

func dnaFWDAlignment(template, query wtype.DNASequence, alignMentMatrix ScoringMatrix) (alignment Result, err error) {

	seq1, seq2 := replaceN(template), replaceN(query)

	if template.Plasmid {
		var tempseq string
		rotationSize := len(query.Seq) - 1
		tempseq += seq1.Seq[rotationSize:]
		tempseq += seq1.Seq[:rotationSize]
		seq1.Seq = tempseq
	}

	fsa := &linear.Seq{Seq: alphabet.BytesToLetters([]byte(seq1.Sequence()))}
	fsa.Alpha = alphabet.DNAgapped
	fsb := &linear.Seq{Seq: alphabet.BytesToLetters([]byte(seq2.Sequence()))}
	fsb.Alpha = alphabet.DNAgapped

	aln, err := alignMentMatrix.Align(fsa, fsb)
	if err == nil {
		fa := align.Format(fsa, fsb, aln, '-')
		alignment = Result{
			Template:  &template,
			Query:     &query,
			Algorithm: alignMentMatrix,
			Alignment: formatMisMatches(Alignment{
				TemplateResult: fmt.Sprint(fa[0]),
				QueryResult:    fmt.Sprint(fa[1]),
			}),
		}
	}
	return
}

// the biogo implementation of alignment requires the N nucleotides to be repalced with -
func replaceN(seq wtype.DNASequence) wtype.DNASequence {

	var newSeq []string

	for _, letter := range seq.Seq {
		if strings.ToUpper(string(letter)) == "N" {
			letter = rune('-')
		}
		newSeq = append(newSeq, string(letter))
	}

	seq.Seq = strings.Join(newSeq, "")

	return seq
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
			queryChar = string(alignment.QueryResult[i])
			templateChar = string(alignment.TemplateResult[i])
		}
		formattedQuery = append(formattedQuery, queryChar)
		formattedTemplate = append(formattedTemplate, templateChar)
	}
	formattedAlignment = Alignment{
		TemplateResult: strings.Join(formattedTemplate, ""),
		QueryResult:    strings.Join(formattedQuery, ""),
	}
	return
}

// standard character representing an alignment gap
const GAP rune = rune('-')

func isGap(character rune) bool {
	if character == GAP {
		return true
	}
	return false
}

func isMismatch(character1, character2 rune) bool {

	if isGap(character1) || isGap(character2) {
		return false
	}

	if character1 != character2 {
		return true
	}
	return false
}
