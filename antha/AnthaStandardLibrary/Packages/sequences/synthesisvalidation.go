// synthesisvalidation.go
package sequences

import (
	"fmt"
	"strings"

	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/search"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
)

// This simulates the sequence assembly reaction to validate if parts will synthesise with intended manufacturer.
// Does not validate construct assembly so should be used in conjunction with enzymes.Assemblysimulator()
func ValidateSynthesis(parts []wtype.DNASequence, vector string, manufacturer string) (string, bool) {

	var status string
	var vectorstatus string
	var bad bool
	var ok bool

	_, inmap := SynthesisStandards[manufacturer]

	if !inmap {
		keys := make([]string, 0)
		for key := range SynthesisStandards {
			keys = append(keys, key)
		}

		status = "Error! Reenter manufacturer, Valid manufaturer entries: " + strings.Join(keys, ", ")
		return status, false
	}

	// type conversions
	a, _ := SynthesisStandards[manufacturer]["RepeatMax"].(int)
	b, _ := SynthesisStandards[manufacturer]["MinOrder"].(int)
	c, _ := SynthesisStandards[manufacturer]["MinLength"].(int)
	d, _ := SynthesisStandards[manufacturer]["MaxLength"].(int)
	vectorR, _ := SynthesisStandards[manufacturer]["Vector"].([]string)

	// min total order. Needs to be before local gc calc
	total := 0
	for _, part := range parts {
		total += len(part.Seq)
	}
	if total < b {
		status = fmt.Sprint("Warning: Total length of parts is less than", manufacturer, "minimum order requirement")

	}

	// check if vector is appropriate
	if !search.InStrings(vectorR, vector) {
		vectorstatus = fmt.Sprint("Warning: Non-standard vector used for", manufacturer,
			"synthesis. Please see manufacturer instructions for the standard vector or the use of custom vectors")
	} else {
		vectorstatus = fmt.Sprint("Vector in list approved by ", manufacturer)
	}

	for _, part := range parts {
		GCC := GCcontent(part.Seq)              // global gc
		gc := localGCContent(part.Seq, 100, 50) // local gc

		// check lengths of seq, repeat content and global gc content of each part
		if len(part.Seq) < c {
			status = status + ". " + fmt.Sprint("Warning:", part.Nm, "is short and may be difficult to synthesise")
			bad = true
		} else if len(part.Seq) > d {
			status = status + ". " + fmt.Sprint("Warning:", part.Nm, "is long and may be difficult to sythesise")
			bad = true
		} else if strings.Contains(strings.ToUpper(part.Seq), strings.Repeat("A", a)) || strings.Contains(strings.ToUpper(part.Seq), strings.Repeat("T", a)) ||
			strings.Contains(strings.ToUpper(part.Seq), strings.Repeat("C", a)) || strings.Contains(strings.ToUpper(part.Seq), strings.Repeat("A", a)) {
			status = status + ". " + fmt.Sprint("Warning:", part.Nm, "is highly repetetive and unsuitable for synthesis")
			bad = true
		} else if GCC > 0.65 || GCC < 0.40 {
			status = status + ". " + fmt.Sprint("Warning: GC content of", part.Nm, "is very high or low and may be difficult to synthesise")
			bad = true
		} else {
			status = status + ". " + fmt.Sprint("Your", manufacturer, "DNA synthesis order should work")
			bad = true
		}

		// check local gc content of each part in 100bp sliding window
		for _, v := range gc {
			if v < 0.25 || v > 0.80 {
				status = fmt.Sprint("Warning: Local GC content very high or low in", part.Nm)
				bad = true
			}
		}
	}
	if status == "" {
		status = fmt.Sprint("Your", manufacturer, "DNA synthesis order should work")
	}

	status = vectorstatus + status

	if bad {
		ok = false
	}

	return status, ok
}

var SynthesisStandards = map[string]map[string]interface{}{
	"Gen9": {
		"Vector":    []string{"pG9m-2"},
		"MaxLength": 10000,
		"MinLength": 400,
		"RepeatMax": 70,
		"MinOrder":  20000,
	},
	"DNA20": {
		"Vector":    []string{"pJ341", "pJ221", "pJ321", "pJ201", "pJ344", "pJ224", "pJ324", "pJ204", "pJ347", "pJ227", "pJ327", "pJ207", "pJ348", "pJ228", "pJ328", "pJ208", "pJ349", "pJ229", "pJ329", "pJ209", "pJ351", "pJ231", "pJ331", "pJ211", "J354", "pJ234", "pJ334", "pJ214", "pJ357", "pJ234", "pJ334", "pJ217", "pJ358", "pJ238", "pJ338", "pJ218", "pJ359", "pJ239", "pJ339", "pJ219", "pM265", "pM268", "pM269", "pM275", "pM278", "pM279", "pM269E-19C", "pM269Y-19C", "pM262", "pM263", "pM264", "pM272", "pM273", "pM273", "pM274"},
		"MaxLength": 3000,
		"MinLength": 400,
		"RepeatMax": 70,
		"MinOrder":  0,
	},
	"GenScript": {
		"Vector":    []string{"pUC57", "pUC57-Kan", "pUC57-Simple", "pUC57-mini", "pUC18", "pUC19"},
		"MaxLength": 8000,
		"MinLength": 400,
		"RepeatMax": 70,
		"MinOrder":  455,
	},
	"GeneWiz": {
		"Vector":    []string{"pUC57"},
		"MaxLength": 10000,
		"MinLength": 200,
		"RepeatMax": 70,
		"MinOrder":  455,
	},
	"OriGene": {
		"Vector":    []string{"pUCAmp", "pUCKan", "pUCAmpMinusMCS", "pUCKanMinusMCS"},
		"MaxLength": 10000,
		"MinLength": 200,
		"RepeatMax": 70,
		"MinOrder":  455,
	},
	"GeneArt": {
		"Vector":    []string{"pUCAmp", "pUCKan", "pUCAmpMinusMCS", "pUCKanMinusMCS"},
		"MaxLength": 10000,
		"MinLength": 200,
		"RepeatMax": 70,
		"MinOrder":  455,
	},
	"EuroFins": {
		"Vector":    []string{"pEX-A2", "pEX-K4"},
		"MaxLength": 10000,
		"MinLength": 200,
		"RepeatMax": 20,
		"MinOrder":  455,
	},
}
