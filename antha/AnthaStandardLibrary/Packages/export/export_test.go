package export

import (
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/sequences"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/sequences/parse/genbank"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"math/rand"
	"strings"
	"testing"
)

func randomDNA(seed int64, length int) string {
	nucs := []byte("ACGT")
	nBases := len(nucs)
	r := rand.New(rand.NewSource(seed))
	bases := []byte{}
	for i := 0; i < length; i++ {
		bases = append(bases, nucs[r.Intn(nBases)])
	}
	return string(bases)
}

func TestGenbankSerial(t *testing.T) {

	seqStr := randomDNA(1, 500)

	setupFeatures := []struct {
		Name       string
		Start, End int // unit offset
		Class      string
		Direction  string
	}{
		{"f1", 121, 140, wtype.MISC_FEATURE, "forward"},
		{"f2", 141, 160, wtype.MISC_FEATURE, "reverse"},
		{"f3", 161, 180, wtype.PROMOTER, "forward"},
	}

	features := []wtype.Feature{}
	for _, p := range setupFeatures {
		featStr := seqStr[p.Start-1 : p.End]
		if strings.ToLower(p.Direction) == "reverse" {
			featStr = wtype.RevComp(featStr)
		}
		feature := sequences.MakeFeature(p.Name, featStr, p.Start, p.End, "DNA", p.Class, p.Direction)
		features = append(features, feature)
	}

	want, err := sequences.MakeAnnotatedSeq("s1", seqStr, false, features)
	if err != nil {
		t.Fatal(err.Error())
	}

	seqFile, _, err := GenbankSerial(LOCAL, "MyOutputFile", []wtype.DNASequence{want})
	if err != nil {
		t.Fatal(err.Error())
	}

	allBytes, err := seqFile.ReadAll()
	if err != nil {
		t.Fatal(err.Error())
	}

	// fmt.Printf("-----\n%s----\n", string(allBytes))

	got, err := genbank.GenbankContentsToAnnotatedSeq(allBytes)

	if err != nil {
		t.Fatal("Failed to parse output file")
	}

	if strings.ToUpper(got.Name()) != strings.ToUpper(want.Name()) {
		t.Errorf("Name: got %s, want %s\n", got.Name(), want.Name())
	}
	if strings.ToUpper(got.Sequence()) != strings.ToUpper(want.Sequence()) {
		t.Errorf("Sequence: got %s, want %s\n", got.Sequence(), want.Sequence())
	}
	if len(got.Features) != len(want.Features) {
		t.Fatalf("Number of features: got %s, want %s\n", got.Sequence(), want.Sequence())
	}
	for i := 0; i < len(want.Features); i++ {
		if got.Features[i].Name != want.Features[i].Name {
			t.Errorf("Feature %d Name: got %s, want %s\n", i, got.Features[i].Name, want.Features[i].Name)
		}
		if got.Features[i].Class != want.Features[i].Class {
			t.Errorf("Feature %d Class: got %s, want %s\n", i, got.Features[i].Class, want.Features[i].Class)
		}
		if got.Features[i].StartPosition != want.Features[i].StartPosition {
			t.Errorf("Feature %d StartPosition: got %d, want %d\n", i, got.Features[i].StartPosition, want.Features[i].StartPosition)
		}
		if got.Features[i].EndPosition != want.Features[i].EndPosition {
			t.Errorf("Feature %d EndPosition: got %d, want %d\n", i, got.Features[i].EndPosition, want.Features[i].EndPosition)
		}
		if got.Features[i].Reverse != want.Features[i].Reverse {
			t.Errorf("Feature %d Reverse: got %v, want %v\n", i, got.Features[i].Reverse, want.Features[i].Reverse)
		}

	}

}
