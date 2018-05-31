package liquidhandling

import (
	"context"
	"fmt"
	"testing"

	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/inventory"
	"github.com/antha-lang/antha/inventory/testinventory"

	"github.com/stretchr/testify/assert"
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
	lhp.AddPlateTo(pos, p)
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

func TestLHPropertiesDup(t *testing.T) {
	props := MakeGilsonWithPlatesAndTipboxesForTest()
	dprops := props.DupKeepIDs()
	AssertLHPropertiesEqual(t, props, dprops, "LHProperties")
}

func AssertLHPropertiesEqual(t *testing.T, e, g *LHProperties, msg string) {
	assert.Equalf(t, e.ID, g.ID, "%s: ID", msg)
	assert.Equalf(t, e.Nposns, g.Nposns, "%s: Nposns", msg)
	assert.Equalf(t, e.Positions, g.Positions, "%s: Positions", msg)
	assert.Equalf(t, e.PosLookup, g.PosLookup, "%s: PosLookup", msg)
	assert.Equalf(t, e.PlateIDLookup, g.PlateIDLookup, "%s: PlateIDLookup", msg)
	assert.Equalf(t, e.Devices, g.Devices, "%s: Devices", msg)
	assert.Equalf(t, e.Model, g.Model, "%s: Model", msg)
	assert.Equalf(t, e.Mnfr, g.Mnfr, "%s: Mnfr", msg)
	assert.Equalf(t, e.LHType, g.LHType, "%s: LHType", msg)
	assert.Equalf(t, e.TipType, g.TipType, "%s: TipType", msg)
	assert.Equalf(t, e.Tip_preferences, g.Tip_preferences, "%s: Tip_preferences", msg)
	assert.Equalf(t, e.Input_preferences, g.Input_preferences, "%s: Input_preferences", msg)
	assert.Equalf(t, e.Output_preferences, g.Output_preferences, "%s: Output_preferences", msg)
	assert.Equalf(t, e.Tipwaste_preferences, g.Tipwaste_preferences, "%s: Tipwaste_preferences", msg)
	assert.Equalf(t, e.Waste_preferences, g.Waste_preferences, "%s: Waste_preferences", msg)
	assert.Equalf(t, e.Wash_preferences, g.Wash_preferences, "%s: Wash_preferences", msg)
	assert.Equalf(t, e.CurrConf, g.CurrConf, "%s: CurrConf", msg)
	assert.Equalf(t, e.Cnfvol, g.Cnfvol, "%s: Cnfvol", msg)
	assert.Equalf(t, e.Layout, g.Layout, "%s: Layout", msg)
	assert.Equalf(t, e.MaterialType, g.MaterialType, "%s: MaterialType", msg)
	assert.Equalf(t, e.Heads, g.Heads, "%s: Heads", msg)
	assert.Equalf(t, e.HeadAssemblies, g.HeadAssemblies, "%s: HeadAssemblies", msg)
}
