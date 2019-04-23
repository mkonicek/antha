package export

import (
	"fmt"
	"strings"
	"testing"

	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/sequences"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/sequences/parse/genbank"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/laboratory"
	"github.com/antha-lang/antha/laboratory/testlab"
)

func TestGenbankSerial(t *testing.T) {
	testlab.WithTestLab(t, "", &testlab.TestElementCallbacks{
		Steps: func(lab *laboratory.Laboratory) error {

			seqStr := strings.Replace("CTTTCGCAAA GTGCAGTCCG TGAGTTTAGT CATTCACTCG CGGTCTGATC CCCTACAGTT TGCAGACGTG", " ", "", -1)
			plasmid := true

			setupFeatures := []struct {
				Name         string
				FeatStr      string
				Class        string
				Direction    string
				WantLocation string
			}{
				{"f1", seqStr[11-1 : 20], wtype.PROMOTER, "forward", "11..20"},
				{"f2", seqStr[21-1 : 40], wtype.MISC_FEATURE, "forward", "21..40"},
				{"f3", seqStr[41-1:] + seqStr[:10], wtype.MISC_FEATURE, "reverse", "complement(join(41..70,1..10))"},
				{"f4", seqStr[41-1:] + seqStr[:10], wtype.MISC_FEATURE, "forward", "join(41..70,1..10)"},
				{"f5", seqStr[41-1 : 60], wtype.MISC_FEATURE, "reverse", "complement(41..60)"},
			}

			features := []wtype.Feature{}
			hasJoin := make(map[string]bool)
			for _, p := range setupFeatures {
				// MakeFeature revcomps featStr where needed
				// MakeAnnotatedSeq locates the feature sequence so start/end defaulted here
				feature := sequences.MakeFeature(p.Name, p.FeatStr, int(-1), int(-1), "DNA", p.Class, p.Direction)
				// fmt.Printf("produced %#v\n", feature)
				features = append(features, feature)
				if strings.Contains(p.WantLocation, "join") {
					hasJoin[feature.Name] = true
				}
			}

			// MakeAnnotatedSeq locates the feature sequences
			want, err := sequences.MakeAnnotatedSeq("s1", seqStr, plasmid, features)
			if err != nil {
				return err
			}
			// fmt.Printf("Annotated sequence %#v\n", want)

			seqFile, err := GenbankSerial(lab, "MyOutputFile", []wtype.DNASequence{want})
			if err != nil {
				return err
			}

			allBytes, err := lab.FileManager.ReadAll(seqFile)
			if err != nil {
				return err
			}

			allStr := string(allBytes)
			for _, p := range setupFeatures {
				if !strings.Contains(allStr, p.WantLocation) {
					return fmt.Errorf("Expected string not found: %s", p.WantLocation)
				}
			}

			// ** Reparse test **

			// Genbank parser throws error on features with join locations
			// https://github.com/antha-lang/antha/blob/master/antha/AnthaStandardLibrary/Packages/sequences/parse/genbank/genbank_parser.go#L228
			// so remove these and re-generate

			newFeatures := []wtype.Feature{}
			for _, feat := range want.Features {
				if !(hasJoin[feat.Name]) {
					newFeatures = append(newFeatures, feat)
				}
			}
			want.Features = newFeatures

			seqFile, err = GenbankSerial(lab, "MyOutputFileNoJoins", []wtype.DNASequence{want})
			if err != nil {
				return err
			}

			allBytes, err = lab.FileManager.ReadAll(seqFile)
			if err != nil {
				return err
			}

			// genbank parsing...

			got, err := genbank.GenbankContentsToAnnotatedSeq(allBytes)
			if err != nil {
				return err
			}

			if !strings.EqualFold(got.Name(), want.Name()) {
				return fmt.Errorf("Name: got %s, want %s\n", got.Name(), want.Name())
			}
			if !strings.EqualFold(got.Sequence(), want.Sequence()) {
				return fmt.Errorf("Sequence: got %s, want %s\n", got.Sequence(), want.Sequence())
			}
			if len(got.Features) != len(want.Features) {
				return fmt.Errorf("Number of features: got %s, want %s\n", got.Sequence(), want.Sequence())
			}
			for i := 0; i < len(want.Features); i++ {
				if got.Features[i].Name != want.Features[i].Name {
					return fmt.Errorf("Feature %d Name: got %s, want %s\n", i, got.Features[i].Name, want.Features[i].Name)
				}
				if got.Features[i].Class != want.Features[i].Class {
					return fmt.Errorf("Feature %d Class: got %s, want %s\n", i, got.Features[i].Class, want.Features[i].Class)
				}
				if got.Features[i].Start() != want.Features[i].Start() {
					return fmt.Errorf("Feature %d StartPosition: got %d, want %d\n", i, got.Features[i].StartPosition, want.Features[i].Start())
				}
				if got.Features[i].End() != want.Features[i].End() {
					return fmt.Errorf("Feature %d EndPosition: got %d, want %d\n", i, got.Features[i].EndPosition, want.Features[i].End())
				}
				if got.Features[i].Reverse != want.Features[i].Reverse {
					return fmt.Errorf("Feature %d Reverse: got %v, want %v\n", i, got.Features[i].Reverse, want.Features[i].Reverse)
				}
				if got.Features[i].DNASeq != want.Features[i].DNASeq {
					return fmt.Errorf("Feature %d DNASequence: got %v, want %v\n", i, got.Features[i].DNASeq, want.Features[i].DNASeq)
				}
			}
			return nil
		},
	})
}
