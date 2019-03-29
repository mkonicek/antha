package align

import (
	"fmt"
	"strings"
)

// QueryStart returns the start position of the alignment in the query
func (a *Alignment) QueryStart() int {
	if len(a.QueryPositions) > 0 {
		return a.QueryPositions[0]
	}
	return -1
}

// QueryEnd returns the end position of the alignment in the query
func (a *Alignment) QueryEnd() int {
	if len(a.QueryPositions) > 0 {
		return a.QueryPositions[len(a.QueryPositions)-1]
	}
	return -1
}

// TemplateStart returns the start position of the alignment in the template
func (a *Alignment) TemplateStart() int {
	if len(a.TemplatePositions) > 0 {
		return a.TemplatePositions[0]
	}
	return -1
}

// TemplateEnd returns the end position of the alignment in the template
func (a *Alignment) TemplateEnd() int {
	if len(a.TemplatePositions) > 0 {
		return a.TemplatePositions[len(a.TemplatePositions)-1]
	}
	return -1
}

// inferFrame infers alignment frame from aligned positions; returns -1 for a reverse alignment, otherwise 1
func inferFrame(positions []int) int {
	// A forward alignment requires at least one positive unit difference at
	// consecutive positions in the positions slice.
	// In gapped alignments the position is constant across the gap region.
	// Therefore, aligment positions are monotonic ascending or descending,
	// except in the case of plasmids, where aligned positions spanning the
	// coordinate origin have a large jump, resulting in a sawtooth profile,
	// rather than monotonic.
	// The jump for plasmid alignments will be at least the length of the
	// plasmid in magnitude and is ignored here, since the inference considers
	// only unit differences.
	for i := 1; i < len(positions); i++ {
		delta := positions[i] - positions[i-1]
		if delta*delta != 1 { // ignore consecutive differences that are not unit
			continue
		}
		if delta < 0 {
			return -1
		}
		return 1
	}
	return 1
}

// TemplateFrame returns -1 if the template is aligned the reverse direction, 1 otherwise
func (a *Alignment) TemplateFrame() int {
	return inferFrame(a.TemplatePositions)
}

// QueryFrame returns -1 if the query is aligned the reverse direction, 1 otherwise
func (a *Alignment) QueryFrame() int {
	return inferFrame(a.QueryPositions)
}

// OutputMismatch defines a character representing an alignment mismatch
const OutputMismatch = " "

// OutputMatch defines a character representing an alignment match
const OutputMatch = "|"

// Match produces a formatted line indicating matches between aligned sequences
/*
	GCTTTTTTAT res1
	|   |||||| <- like this
	GGG-TTTTAT res2
*/
func (a *Alignment) Match() string {

	n := len(a.QueryResult)
	if len(a.TemplateResult) != n {
		return ""
	}

	ret := make([]string, n)
	for i := 0; i < n; i++ {
		ret[i] = OutputMismatch
		if a.QueryResult[i] == a.TemplateResult[i] {
			ret[i] = OutputMatch
		}
	}

	return strings.Join(ret, "")
}

// Split an alignment into sections of up to a specified length, help formatting
func (a *Alignment) Split(maxSectionLength int) ([]Alignment, error) {

	n := len(a.TemplateResult)

	if len(a.QueryResult) != n {
		return []Alignment{}, fmt.Errorf("lengths differ %d, %d", n, len(a.QueryResult))
	}
	if len(a.TemplatePositions) != n {
		return []Alignment{}, fmt.Errorf("inconsistent number of template positions: %d, %d", n, len(a.TemplatePositions))
	}
	if len(a.QueryPositions) != n {
		return []Alignment{}, fmt.Errorf("inconsistent number of query positions: %d, %d", n, len(a.QueryPositions))
	}

	sections := make([]Alignment, 0)

	for i := 0; i < n; i += maxSectionLength {

		endLine := i + maxSectionLength
		if endLine > n {
			endLine = n
		}

		aln := Alignment{}

		aln.TemplateResult = a.TemplateResult[i:endLine]
		aln.QueryResult = a.QueryResult[i:endLine]
		aln.TemplatePositions = a.TemplatePositions[i:endLine]
		aln.QueryPositions = a.QueryPositions[i:endLine]

		sections = append(sections, aln)

	}

	return sections, nil
}
