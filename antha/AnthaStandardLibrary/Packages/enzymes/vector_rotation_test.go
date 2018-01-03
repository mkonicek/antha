package enzymes

import (
	"fmt"
	"testing"

	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/enzymes/lookup"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
)

func TestNoRotationNeeded(t *testing.T) {
	enzyme, _ := lookup.TypeIIs("SAPI")
	seq := "GCTCTTCxxxxx"
	rseq := "GCTCTTCxxxxx"

	s := wtype.DNASequence{Nm: "nevermind", Seq: seq}

	rs, err := rotateVector(s, enzyme)

	if err != nil {
		t.Fatal(err)
	}

	if rs.Seq != rseq {
		t.Fatal(fmt.Sprintf("Error with vector rotation: got %s expected %s", s.Sequence(), rs.Sequence()))
	}
}
func TestSomeRotationNeeded(t *testing.T) {
	enzyme, _ := lookup.TypeIIs("SAPI")
	seq := "xxxxxGCTCTTCn"
	rseq := "GCTCTTCnxxxxx"

	s := wtype.DNASequence{Nm: "nevermind", Seq: seq}

	rs, err := rotateVector(s, enzyme)

	if err != nil {
		t.Fatal(err)
	}

	if rs.Seq != rseq {
		t.Fatal(fmt.Sprintf("Error with vector rotation: got %s expected %s", rs.Seq, rseq))
	}
}
