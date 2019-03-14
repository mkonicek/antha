package align

import (
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"sort"
)

// DNASet aligns a query to a collection (or database) of sequences, testing both
// forward and reverse directions. It returns the top scoring alignment results
// found in rank order, up to a specified number.
func DNASet(query wtype.DNASequence, templates []wtype.DNASequence, alignmentMatrix ScoringMatrix, maxResults int) ([]Result, error) {

	results := make([]Result, 0)

	for _, template := range templates {

		fwdResult, err := DNAFwd(template, query, alignmentMatrix)

		if err != nil {
			return []Result{}, err
		}

		revResult, err := DNARev(template, query, alignmentMatrix)

		if err != nil {
			return []Result{}, err
		}

		best := fwdResult
		if revResult.Score() > best.Score() {
			best = revResult
		}

		results = append(results, best)

	}

	sortFunc := func(i, j int) bool {
		return results[i].Score() > results[j].Score()
	}

	sort.Slice(results, sortFunc)

	end := maxResults
	if end > len(results) {
		end = len(results)
	}

	return results[:end], nil
}
