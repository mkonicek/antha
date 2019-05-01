package main

import (
	"context"
	"fmt"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/inventory"
	"github.com/antha-lang/antha/inventory/testinventory"
	"strings"
	"testing"
)

// TODO
func TestInventoryLHPlateSerialize(t *testing.T) {
	ctx := testinventory.NewContext(context.Background())
	newPlates, err := inventory.XXXNewPlates(ctx)

	if err != nil {
		t.Errorf("%v", err)
	}

	oldPlates := makePlates()

	if len(newPlates) != len(oldPlates) {
		t.Errorf("Serialization error: %d plates made, only %d available from inventory (update the library in ../plate_library.go)", len(oldPlates), len(newPlates))
	}

	newPlateMap := make(map[string]*wtype.Plate)
	for _, p := range newPlates {
		newPlateMap[p.Type] = p
	}

	oldPlateMap := make(map[string]*wtype.Plate)
	for _, p := range oldPlates {
		oldPlateMap[p.Type] = p
	}

	for name := range oldPlateMap {
		p := oldPlateMap[name]
		p2, ok := newPlateMap[name]

		if !ok {
			t.Errorf("No plate %s in new plate library", name)
		}

		if !strings.Contains(p.Type, "FromSpec") {
			p.WellXStart += xStartOffset
			p.WellYStart += yStartOffset
			p.WellZStart += zStartOffset
		}

		fMErr := func(s string, want, got interface{}) string {
			return fmt.Sprintf(
				"%s not maintained after marshal/unmarshal: want: %v, got: %v",
				s,
				want,
				got)
		}

		for i := 0; i < p2.WellsX(); i++ {
			for j := 0; j < p2.WellsY(); j++ {
				wc := wtype.WellCoords{X: i, Y: j}

				w := p2.Wellcoords[wc.FormatA1()]

				w.WContents.CName = wc.FormatA1()
				if p2.Rows[j][i].WContents.CName != wc.FormatA1() || p2.Cols[i][j].WContents.CName != wc.FormatA1() || p2.HWells[w.ID].WContents.CName != wc.FormatA1() {
					fmt.Println(p2.Cols[i][j].WContents.CName)
					fmt.Println(p2.Rows[j][i].WContents.CName)
					t.Errorf("Error: Wells inconsistent at position %s", wc.FormatA1())
				}

			}
		}

		if p.Type != p2.Type {
			t.Errorf(fMErr("Type", p.Type, p2.Type))
		}

		if p.Mnfr != p2.Mnfr {
			t.Errorf(fMErr("Manufacturer", p.Mnfr, p2.Mnfr))
		}

		if p.Nwells != p2.Nwells {
			t.Errorf(fMErr("NWells", p.Nwells, p2.Nwells))
		}

		if p.Height() != p2.Height() {
			t.Errorf(fMErr("Height", p.Height(), p2.Height()))
		}

		if p.WellXOffset != p2.WellXOffset {
			t.Errorf(fMErr("WellXOffset", p.WellXOffset, p2.WellXOffset))
		}

		if p.WellYOffset != p2.WellYOffset {
			t.Errorf(fMErr("WellYOffset", p.WellYOffset, p2.WellYOffset))
		}

		if p.WellXStart != p2.WellXStart {
			t.Errorf(fMErr("WellXStart", p.WellXStart, p2.WellXStart))
		}

		if p.WellYStart != p2.WellYStart {
			t.Errorf(fMErr("WellYStart", p.WellYStart, p2.WellYStart))
		}

		if p.WellZStart != p2.WellZStart {
			t.Errorf(fMErr("WellZStart", p.WellZStart, p2.WellZStart))
		}

		if !p.Bounds.Equals(p2.Bounds) {
			t.Errorf(fMErr("Bounds ", p.Bounds, p2.Bounds))
		}

		if !wellTypeEqual(p.Welltype, p2.Welltype) {
			t.Errorf(fMErr("Weltype", p.Welltype, p2.Welltype))
		}

		if !compareExtra(p, p2) {
			t.Errorf(fMErr("Extra", p.Welltype.Extra, p2.Welltype.Extra))
		}
	}
}

// defines an equivalence relation over well *types* - will return true
// if types are same but contents are different
// we also do NOT compare Extra since it contains a mix of instance and type information (ugh)
func wellTypeEqual(self, w2 *wtype.LHWell) bool {
	return self.Crds == w2.Crds && self.MaxVol == w2.MaxVol && self.Rvol == w2.Rvol && self.WShape.Equals(w2.WShape) && self.Bounds.Equals(w2.Bounds) && self.Bottomh == w2.Bottomh && self.Extra["IMSPECIAL"] == w2.Extra["IMSPECIAL"]
}

// note order of arguments here... this will incorrectly fail if order is disregarded
func compareExtra(nonSerializedP, previouslySerializedP *wtype.Plate) bool {
	ex1 := nonSerializedP.Welltype.Extra
	ex2 := previouslySerializedP.Welltype.Extra

	if len(ex1) != len(ex2) {
		return false
	}

	for k, v := range ex1 {
		v2, ok := ex2[k]

		if !ok {
			return false
		}

		sa, arr := v.([]string)

		if arr {
			ifa, arr := v2.([]interface{})

			if !arr {
				return false
			}

			if len(sa) != len(ifa) {
				return false
			}

			for i := 0; i < len(sa); i++ {
				s2, ok := ifa[i].(string)

				if !ok {
					return false
				}

				if sa[i] != s2 {
					return false
				}
			}

		} else {
			if v != v2 {
				return false
			}
		}

	}
	return true
}
