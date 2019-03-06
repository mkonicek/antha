package liquidhandling

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/inventory"
	"github.com/antha-lang/antha/inventory/testinventory"

	"github.com/stretchr/testify/assert"
)

func TestSavePlates(t *testing.T) {
	lhp := MakeGilsonForTest(defaultTipList())
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
	if err := lhp.AddPlateTo(pos, p); err != nil {
		t.Fatal(err)
	}
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

func TestLHPropertiesSane(t *testing.T) {
	props := MakeGilsonForTest(defaultTipList())

	assertPropsSane(t, props)
}

func assertPropsSane(t *testing.T, props *LHProperties) {

	heads := make(map[*wtype.LHHead]bool, len(props.Heads))
	for _, head := range props.Heads {
		heads[head] = true
	}

	adaptors := make(map[*wtype.LHAdaptor]bool, len(props.Adaptors))
	for _, adaptor := range props.Adaptors {
		adaptors[adaptor] = true
	}

	for i, loadedHead := range props.GetLoadedHeads() {
		if _, ok := heads[loadedHead]; !ok {
			t.Errorf("head index %d not defined in machine", i)
		}
		if _, ok := adaptors[loadedHead.Adaptor]; !ok {
			t.Errorf("adaptor at index %d not defined in machine", i)
		}
	}
}

func TestLHPropertiesDup(t *testing.T) {
	props := MakeGilsonWithPlatesAndTipboxesForTest("")
	dprops := props.DupKeepIDs()
	assertPropsSane(t, dprops)
	AssertLHPropertiesEqual(t, props, dprops, "LHProperties")

	for _, head := range props.Heads {
		head.Name = "Changed_" + head.Name
	}
	for i, head := range dprops.Heads {
		if head.Name == props.Heads[i].Name {
			t.Error("Props.Heads not duplicated properly")
			break
		}
	}

	for _, adaptor := range props.Adaptors {
		adaptor.Name = "Changed_" + adaptor.Name
	}
	for i, adaptor := range dprops.Adaptors {
		if adaptor.Name == props.Adaptors[i].Name {
			t.Error("Props.Adaptors not duplicated properly")
			break
		}
	}

	for _, plate := range props.Plates {
		plate.ID = "Changed_" + plate.ID
	}
	for i, plate := range dprops.Plates {
		if plate.ID == props.Plates[i].ID {
			t.Error("props.Plates not duplicated properly")
			break
		}
	}

	for _, tipbox := range props.Tipboxes {
		tipbox.ID = "changed_" + tipbox.ID
	}
	for i, tipbox := range dprops.Tipboxes {
		if tipbox.ID == props.Tipboxes[i].ID {
			t.Error("props.Tipboxes not duplicated properly")
			break
		}
	}
}

func AssertLHPropertiesEqual(t *testing.T, e, g *LHProperties, msg string) {
	assert.Equalf(t, e.ID, g.ID, "%s: ID", msg)
	assert.Equalf(t, e.Positions, g.Positions, "%s: Positions", msg)
	assert.Equalf(t, e.PosLookup, g.PosLookup, "%s: PosLookup", msg)
	assert.Equalf(t, e.PlateIDLookup, g.PlateIDLookup, "%s: PlateIDLookup", msg)
	assert.Equalf(t, e.Model, g.Model, "%s: Model", msg)
	assert.Equalf(t, e.Mnfr, g.Mnfr, "%s: Mnfr", msg)
	assert.Equalf(t, e.LHType, g.LHType, "%s: LHType", msg)
	assert.Equalf(t, e.TipType, g.TipType, "%s: TipType", msg)
	assert.Equalf(t, e.Preferences, g.Preferences, "%s: Peferences", msg)
	assert.Equalf(t, e.Heads, g.Heads, "%s: Heads", msg)
	assert.Equalf(t, e.Adaptors, g.Adaptors, "%s: Adaptors", msg)
	assert.Equalf(t, e.HeadAssemblies, g.HeadAssemblies, "%s: HeadAssemblies", msg)
}

func TestLHPropertiesSerialisation(t *testing.T) {
	before := MakeGilsonWithPlatesAndTipboxesForTest("")

	// we don't need to preserve this
	for _, tip := range before.Tips {
		tip.ClearParent()
	}

	var after LHProperties
	if data, err := json.Marshal(before); err != nil {
		t.Error(err)
	} else if err := json.Unmarshal(data, &after); err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(before, &after) {
		t.Errorf("serialization changed LHProperties PlateLookup\nbefore: %+v\nafter : %+v", before, &after)
	}

	heads := make(map[*wtype.LHHead]bool)
	for _, h := range after.Heads {
		heads[h] = true
	}
	for _, ha := range after.HeadAssemblies {
		for _, hap := range ha.Positions {
			if hap.Head != nil && !heads[hap.Head] {
				t.Error("HeadAssemblyPosition.Head doesn't point to anything in LHProperties.Heads")
			}
		}
	}

	adaptors := make(map[*wtype.LHAdaptor]bool)
	for _, a := range after.Adaptors {
		adaptors[a] = true
	}
	for _, head := range after.Heads {
		if head.Adaptor != nil && !adaptors[head.Adaptor] {
			t.Error("Head.Adaptor doesn't point to anything in LHProperties.Adaptors")
		}
	}

	if b, a := before.GetLoadedHeads(), after.GetLoadedHeads(); len(b) != len(a) {
		t.Errorf("number of loaded heads doesn't match: before = %d, after = %d", len(b), len(a))
	} else {
		for i, beforeHead := range b {
			if afterHead := a[i]; !reflect.DeepEqual(beforeHead, afterHead) {
				t.Errorf("%dth head is mismatched\nbefore: %+v\nafter : %+v", i, beforeHead, afterHead)
			}
		}
	}

}

type UpdateIDTest struct {
	Name      string
	UpdateIDs []string // list of IDs to update
	ErrorIDs  []string // ids that should produce an error
}

func (test *UpdateIDTest) Run(t *testing.T) {
	t.Run(test.Name, test.run)
}

func (test *UpdateIDTest) run(t *testing.T) {
	// this test is quite involved, mainly in order to make sure that wastes
	// and washes aren't lost, as could sometimes be the case with the previous implementation
	//
	// strategy is to call UpdateID on a copy, then apply the ID mapping manually to
	// the original and assert that everything matches

	// build some properties to test with
	rbt := MakeGilsonForTest(defaultTipList())
	ctx := testinventory.NewContext(context.Background())

	// add something for each object type
	if plt, err := inventory.NewPlate(ctx, "pcrplate_skirted"); err != nil {
		t.Fatal(err)
	} else {
		plt.ID = "initial_plate_id"
		if err := rbt.AddInputPlate(plt); err != nil {
			t.Fatal(err)
		}
	}
	if plt, err := inventory.NewPlate(ctx, "pcrplate_skirted"); err != nil {
		t.Fatal(err)
	} else {
		plt.ID = "initial_wash_id"
		rbt.AddWash(plt)
	}
	if plt, err := inventory.NewPlate(ctx, "pcrplate_skirted"); err != nil {
		t.Fatal(err)
	} else {
		plt.ID = "initial_waste_id"
		rbt.AddWaste(plt)
	}
	if tb, err := inventory.NewTipbox(ctx, "Gilson20"); err != nil {
		t.Fatal(err)
	} else {
		tb.ID = "initial_tipbox_id"
		if err := rbt.AddTipBox(tb); err != nil {
			t.Fatal(err)
		}
	}
	if tw, err := inventory.NewTipwaste(ctx, "Gilsontipwaste"); err != nil {
		t.Fatal(err)
	} else {
		tw.ID = "initial_tipwaste_id"
		if err := rbt.AddTipWaste(tw); err != nil {
			t.Fatal(err)
		}
	}

	// update the IDs as required

	beforeToAfter := make(map[string]string, len(test.UpdateIDs))
	for id := range rbt.PlateLookup {
		beforeToAfter[id] = id
	}

	final := rbt.DupKeepIDs()
	errorIDs := make([]string, 0, len(test.UpdateIDs))
	for _, id := range test.UpdateIDs {
		if newID, err := final.UpdateID(id); err != nil {
			errorIDs = append(errorIDs, id)
		} else {
			beforeToAfter[id] = newID
		}
	}
	if !assert.ElementsMatch(t, test.ErrorIDs, errorIDs) || len(test.ErrorIDs) != 0 {
		// don't continue testing if we were expecting an error
		return
	}

	// update all the fields
	pl := make(map[string]interface{}, len(rbt.PlateLookup))
	for id, obj := range rbt.PlateLookup {
		// rely on later to actually update obj.ID, where we have an object rather than interface
		pl[beforeToAfter[id]] = obj
	}
	rbt.PlateLookup = pl

	posL := make(map[string]string, len(rbt.PosLookup))
	plateIDL := make(map[string]string, len(rbt.PosLookup))
	for addr, id := range rbt.PosLookup {
		posL[addr] = beforeToAfter[id]
		plateIDL[beforeToAfter[id]] = addr
	}
	rbt.PosLookup = posL
	rbt.PlateIDLookup = plateIDL

	plates := make(map[string]*wtype.LHPlate, len(rbt.Plates))
	for id, plate := range rbt.Plates {
		plate.ID = beforeToAfter[id]
		plates[beforeToAfter[id]] = plate
	}
	rbt.Plates = plates

	tipboxes := make(map[string]*wtype.LHTipbox, len(rbt.Tipboxes))
	for id, tipbox := range rbt.Tipboxes {
		tipbox.ID = beforeToAfter[id]
		tipboxes[beforeToAfter[id]] = tipbox
	}
	rbt.Tipboxes = tipboxes

	tipwastes := make(map[string]*wtype.LHTipwaste, len(rbt.Tipwastes))
	for id, tipwaste := range rbt.Tipwastes {
		tipwaste.ID = beforeToAfter[id]
		tipwastes[beforeToAfter[id]] = tipwaste
	}
	rbt.Tipwastes = tipwastes

	wastes := make(map[string]*wtype.LHPlate, len(rbt.Wastes))
	for id, waste := range rbt.Wastes {
		wastes[beforeToAfter[id]] = waste
	}
	rbt.Wastes = wastes

	washes := make(map[string]*wtype.LHPlate, len(rbt.Washes))
	for id, wash := range rbt.Washes {
		washes[beforeToAfter[id]] = wash
	}
	rbt.Washes = washes

	AssertLHPropertiesEqual(t, rbt, final, "expected LHProperties differ")
}

type UpdateIDTests []UpdateIDTest

func (tests UpdateIDTests) Run(t *testing.T) {
	for _, test := range tests {
		test.Run(t)
	}
}

func TestUpdateID(t *testing.T) {
	UpdateIDTests{
		{
			Name:      "Plate",
			UpdateIDs: []string{"initial_plate_id"},
		},
		{
			Name:      "Tipbox",
			UpdateIDs: []string{"initial_tipbox_id"},
		},
		{
			Name:      "Tipwaste",
			UpdateIDs: []string{"initial_tipwaste_id"},
		},
		{
			Name:      "Wash",
			UpdateIDs: []string{"initial_wash_id"},
		},
		{
			Name:      "Waste",
			UpdateIDs: []string{"initial_waste_id"},
		},
		{
			Name:      "All",
			UpdateIDs: []string{"initial_plate_id", "initial_tipbox_id", "initial_tipwaste_id", "initial_wash_id", "initial_waste_id"},
		},
		{
			Name:      "unknown id",
			UpdateIDs: []string{"nonexistent_id"},
			ErrorIDs:  []string{"nonexistent_id"},
		},
	}.Run(t)
}
