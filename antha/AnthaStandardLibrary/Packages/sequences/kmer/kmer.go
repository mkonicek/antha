package kmer

import (
	"sort"

	"github.com/antha-lang/antha/antha/anthalib/wtype"
)

func Hashseq(s string, n int) map[string]int {
	m := make(map[string]int, len(s)-n)
	for i := 0; i < len(s)-n; i++ {
		w := s[i : i+n]
		m[w] += 1
	}
	return m
}

func CompHash(h1, h2 map[string]int) float64 {
	r := 0
	s := 0
	for k, v := range h1 {
		v2, ok := h2[k]

		if !ok {
			v2 = 0
		}

		if v > v2 {
			r += v - v2
		} else {
			r += v2 - v
		}

		s += v
	}

	return float64(r) / float64(s)
}

type HashDB struct {
	Hashes []map[string]int
	Names  []string
	K      int
}

func (hdb *HashDB) Init(seqs []wtype.BioSequence, n int) {
	hdb.Hashes = make([]map[string]int, len(seqs))
	hdb.Names = make([]string, len(seqs))
	hdb.K = n
	for i, seq := range seqs {
		hdb.Names[i] = seq.Name()
		hdb.Hashes[i] = Hashseq(seq.Sequence(), n)
	}
}

type Hit struct {
	Name  string
	Score float64
}

type SortableHits []Hit

func (s SortableHits) Swap(i, j int) {
	x := s[i]
	s[i] = s[j]
	s[j] = x
}

func (s SortableHits) Less(i, j int) bool {
	return s[i].Score < s[j].Score
}

func (s SortableHits) Len() int {
	return len(s)
}

func (hdb HashDB) SearchWith(seq wtype.BioSequence) []Hit {
	h := Hashseq(seq.Sequence(), hdb.K)
	r := make([]Hit, len(hdb.Hashes))

	for i := 0; i < len(hdb.Names); i++ {

		s := CompHash(h, hdb.Hashes[i])
		r[i] = Hit{Name: hdb.Names[i], Score: s}
	}

	sort.Sort(SortableHits(r))
	return r
}
