package liquidhandling

import (
	"context"
	"fmt"
	"testing"

	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/inventory"
	"github.com/antha-lang/antha/inventory/testinventory"
)

func TestSavePlates(t *testing.T) {
	lhp := MakeGilsonForTest()
	ctx := testinventory.NewContext(context.Background())

	p, err := inventory.NewPlate(ctx, "pcrplate_skirted")
	if err != nil {
		t.Fatal(err)
	}
	c := wtype.NewLHComponent()
	v := 100.0
	pos := "position_1"
	c.CName = "mushroom soup"
	c.Vol = v
	c.Vunit = "ul"
	err = p.Wellcoords["A1"].AddComponent(c)
	if err != nil {
		t.Fatal(err)
	}
	p.Wellcoords["A1"].SetUserAllocated()
	lhp.AddPlate(pos, p)
	pl := lhp.SaveUserPlates()

	if len(pl) != 1 {
		t.Fatal(fmt.Sprintf("Error: SaveUserPlates should have 1 plate, instead has %d", len(pl)))
	}

	if pl[0].Position != pos {
		t.Fatal(fmt.Sprintf("Error: SaveUserPlates should return plate at position %s, instead got %s", pos, pl[0].Position))
	}

	if pl[0].Plate.ID != p.ID {
		t.Fatal(fmt.Sprintf("Error: SaveUserPlates should return plate with ID %s, instead got %s", p.ID, pl[0].Plate.ID))
	}

	if pl[0].Plate == p {
		t.Fatal("Error: SaveUserPlates must return a duplicate")
	}

	p.Wellcoords["A1"].WContents.Vol = 20.0
	p.Wellcoords["A2"].WContents.CName = "brown rice"
	p.Wellcoords["A2"].WContents.Vol = 30.0
	p.Wellcoords["A2"].WContents.Vunit = "ul"

	lhp.RestoreUserPlates(pl)

	pp := lhp.Plates[pos]

	w := pp.Wellcoords["A1"]

	if w.WContents.CName != c.CName || w.WContents.Vol != c.Vol || w.WContents.Vunit != c.Vunit {
		t.Fatal(fmt.Sprintf("Error: Restored plate should have component %v at A1, instead got %v", c, w.WContents))
	}

	w = pp.Wellcoords["A2"]
	w2 := p.Wellcoords["A2"]
	if w.WContents.CName != w2.WContents.CName || w.WContents.Vol != w2.WContents.Vol || w.WContents.Vunit != w2.WContents.Vunit {
		t.Fatal(fmt.Sprintf("Error: Restored plate should have  component %v at A2, instead got %v", w2.WContents, w.WContents))
	}

}

func TestGetFirstDefined(t *testing.T) {
	for i := 0; i < 100; i++ {
		sa := make([]string, 100)
		sa[i] = "big"

		d := getFirstDefined(sa)

		if d != i {
			t.Errorf("getFirstDefined returned %d, should have returned %d", d, i)
		}
	}
}
