// blast tests
package blast

import (
	"testing"

	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/sequences/biogo/ncbi/blast"
)

func TestBLAST(t *testing.T) {

	t.Skip("skipping test - calls remote blast server")

	putParams := blast.PutParameters{Program: "blastp", HitListSize: 1, Database: "pdb"}
	query := "SANEDMPVEKILEAELAVEPKTETYVEANMGLNPSSPNDPVTNICQAADKQLFTLVEWAKRIPHFSELPLDDQVILLRAG" // 1DKF:A:1-70

	o, err := BLAST(query, &putParams)
	if err != nil {
		t.Error(err.Error())
	}

	hits, err := Hits(o)
	if err != nil {
		t.Error(err.Error())
	}

	want := 1
	if len(hits) != want {
		t.Errorf("Unexpected number of hits: got %d, want %d", len(hits), want)
	}

}

func TestSimpleBlast(t *testing.T) {

	t.Skip("skipping test - calls remote blast server")

	query := "ACTAGACCAGCGAGCCCTAGGAGGGCTGGGCAGGGCCTTGTCCCCTGGTGCCTCCCGGAGCCCAGCACGGGGCAGGGCCAGCTGGGCACTTAGTTTGGACGAGCTCTGTCCCTCCAGGGGCTTCACTTTCGCCTTCTGCAG"

	o, err := SimpleBlast(query)
	if err != nil {
		t.Error(err.Error())
	}

	hits, err := Hits(o)
	if err != nil {
		t.Error(err.Error())
	}

	if len(hits) == 0 {
		t.Errorf("No hits, expected some hits")
	}

}
