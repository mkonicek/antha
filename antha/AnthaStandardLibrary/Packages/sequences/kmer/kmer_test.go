package kmer

import (
	"testing"

	"github.com/antha-lang/antha/antha/anthalib/wtype"
)

func Test1(t *testing.T) {
	ss := make([]*wtype.DNASequence, 2)

	s1 := wtype.DNASequence{Nm: "Hi", Seq: "AAAAACBJADGBCKADUBCAUHIAHFL"}
	ss[0] = &s1
	s2 := wtype.DNASequence{Nm: "There", Seq: "IUHISUNVSJNBLSKNBLISNB"}
	ss[1] = &s2

	var hdb HashDB

	hdb.Init(wtype.DNASeqSet(ss).AsBioSequences(), 2)

	s3 := wtype.DNASequence{Nm: "MAN", Seq: "IFNOIJNFOSJNBSKNB"}

	sr := hdb.SearchWith(&s3)

	if sr[0].Name != "There" {
		t.Errorf("Expected first hit to be sequence \"There\", instead got \"%s\"", sr[0].Name)
	}
	sr = hdb.SearchWith(ss[0])

	if sr[0].Name != "Hi" {
		t.Errorf("Expected first hit to be sequence \"Hi\", instead got \"%s\"", sr[0].Name)
	}

	sr = hdb.SearchWith(ss[1])

	if sr[0].Name != "There" {
		t.Errorf("Expected first hit to be sequence \"There\", instead got \"%s\"", sr[0].Name)
	}
}
