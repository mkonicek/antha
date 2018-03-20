package wtype

import (
	"fmt"
	"math/rand"
	"regexp"
	"sort"
	"strconv"
)

type CodonSet map[string]float64

func (cs CodonSet) Codons() []string {
	s := make([]string, 0, len(cs))

	for k := range cs {
		s = append(s, k)
	}

	sort.Strings(s)

	return s
}

func (cs CodonSet) ChooseWeighted() string {
	s := ""

	cds := cs.Codons()

	for {
		// assume codons are weighted to 1.0

		c := rand.Intn(len(cs))

		f := rand.Float64()

		if f <= cs[cds[c]] {
			s = cds[c]
			break
		}
	}

	return s
}

type CodonTable struct {
	TaxID     string
	CodonByAA map[string]CodonSet // maps amino acid to set of codons
	AAByCodon map[string]string   // maps codon to amino acid
}

func NewCodonTable() CodonTable {
	byAA := make(map[string]CodonSet, 20)
	byCodon := make(map[string]string, 64)

	return CodonTable{CodonByAA: byAA, AAByCodon: byCodon}
}

func (ct CodonTable) ChooseWeighted(aa string) string {
	set, ok := ct.CodonByAA[aa]

	if !ok {
		return ""
	}

	return set.ChooseWeighted()
}

var (
	TAXLINE = regexp.MustCompile(`^taxid\s+(\d+)`)
	CODLINE = regexp.MustCompile(`^([ACTG]{3,3})\s+([A-Z*])\s+(\d+)\s+([0-9.]+)`)
)

// these are tables devised by wholesale analysis of CDS sets
func ParseCodonTableSimpleFormat(sa []string) CodonTable {
	ct := NewCodonTable()

	for ix, line := range sa {
		if ix == 0 {
			hits := TAXLINE.FindStringSubmatch(line)

			if hits == nil {
				panic("BAD CODON TABLE FORMAT")
			}

			ct.TaxID = hits[1]
			continue
		}

		tox := CODLINE.FindStringSubmatch(line)

		if tox != nil {
			freq, err := strconv.ParseFloat(tox[4], 64)

			if err != nil {
				panic(fmt.Sprintf("Line %d: %s", ix+1, err.Error()))
			}

			ct.AAByCodon[tox[1]] = tox[2]

			cs, ok := ct.CodonByAA[tox[2]]

			if !ok {
				cs = make(CodonSet, 3)
			}

			cs[tox[1]] = freq

			ct.CodonByAA[tox[2]] = cs
		}
	}

	return ct
}
