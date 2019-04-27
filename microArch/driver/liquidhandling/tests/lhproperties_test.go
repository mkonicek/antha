package tests

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/laboratory"
	"github.com/antha-lang/antha/laboratory/testlab"
	"github.com/antha-lang/antha/microArch/driver/liquidhandling"

	"github.com/stretchr/testify/assert"
)

func TestSavePlates(t *testing.T) {
	testlab.WithTestLab(t, "", &testlab.TestElementCallbacks{
		Steps: func(lab *laboratory.Laboratory) error {

			lhp := MakeGilsonForTest(lab, defaultTipList())

			p, err := lab.Inventory.Plates.NewPlate("pcrplate_skirted")
			if err != nil {
				return err
			}
			c := wtype.NewLHComponent(lab.IDGenerator)
			v := 100.0
			pos := "position_1"
			c.CName = "mushroom soup"
			c.Vol = v
			c.Vunit = "ul"
			err = p.Wellcoords["A1"].AddComponent(lab.IDGenerator, c)
			if err != nil {
				return err
			}
			p.Wellcoords["A1"].SetUserAllocated()
			if err := lhp.AddPlateTo(pos, p); err != nil {
				return err
			}
			pl := lhp.SaveUserPlates(lab.IDGenerator)

			if len(pl) != 1 {
				return fmt.Errorf("Error: SaveUserPlates should have 1 plate, instead has %d", len(pl))
			}

			if pl[0].Position != pos {
				return fmt.Errorf("Error: SaveUserPlates should return plate at position %s, instead got %s", pos, pl[0].Position)
			}

			if pl[0].Plate.ID != p.ID {
				return fmt.Errorf("Error: SaveUserPlates should return plate with ID %s, instead got %s", p.ID, pl[0].Plate.ID)
			}

			if pl[0].Plate == p {
				return errors.New("Error: SaveUserPlates must return a duplicate")
			}

			p.Wellcoords["A1"].WContents.Vol = 20.0
			p.Wellcoords["A2"].WContents.CName = "brown rice"
			p.Wellcoords["A2"].WContents.Vol = 30.0
			p.Wellcoords["A2"].WContents.Vunit = "ul"

			lhp.RestoreUserPlates(pl)

			pp := lhp.Plates[pos]

			w := pp.Wellcoords["A1"]

			if w.WContents.CName != c.CName || w.WContents.Vol != c.Vol || w.WContents.Vunit != c.Vunit {
				return fmt.Errorf("Error: Restored plate should have component %v at A1, instead got %v", c, w.WContents)
			}

			w = pp.Wellcoords["A2"]
			w2 := p.Wellcoords["A2"]
			if w.WContents.CName != w2.WContents.CName || w.WContents.Vol != w2.WContents.Vol || w.WContents.Vunit != w2.WContents.Vunit {
				return fmt.Errorf("Error: Restored plate should have  component %v at A2, instead got %v", w2.WContents, w.WContents)
			}
			return nil
		},
	})
}

func TestGetFirstDefined(t *testing.T) {
	for i := 0; i < 100; i++ {
		sa := make([]string, 100)
		sa[i] = "big"

		d := liquidhandling.GetFirstDefined(sa)

		if d != i {
			t.Errorf("GetFirstDefined returned %d, should have returned %d", d, i)
		}
	}
}

func TestLHPropertiesSane(t *testing.T) {
	testlab.WithTestLab(t, "", &testlab.TestElementCallbacks{
		Steps: func(lab *laboratory.Laboratory) error {

			lhp := MakeGilsonForTest(lab, defaultTipList())
			return assertPropsSane(lhp)
		},
	})
}

func assertPropsSane(props *liquidhandling.LHProperties) error {

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
			return fmt.Errorf("head index %d not defined in machine", i)
		}
		if _, ok := adaptors[loadedHead.Adaptor]; !ok {
			return fmt.Errorf("adaptor at index %d not defined in machine", i)
		}
	}
	return nil
}

func TestLHPropertiesDup(t *testing.T) {
	testlab.WithTestLab(t, "", &testlab.TestElementCallbacks{
		Steps: func(lab *laboratory.Laboratory) error {
			props := MakeGilsonWithPlatesAndTipboxesForTest(lab, "")
			dprops := props.DupKeepIDs(lab.IDGenerator)
			if err := assertPropsSane(dprops); err != nil {
				return err
			}
			AssertLHPropertiesEqual(t, props, dprops, "LHProperties")

			for _, head := range props.Heads {
				head.Name = "Changed_" + head.Name
			}
			for i, head := range dprops.Heads {
				if head.Name == props.Heads[i].Name {
					return errors.New("Props.Heads not duplicated properly")
				}
			}

			for _, adaptor := range props.Adaptors {
				adaptor.Name = "Changed_" + adaptor.Name
			}
			for i, adaptor := range dprops.Adaptors {
				if adaptor.Name == props.Adaptors[i].Name {
					return errors.New("Props.Adaptors not duplicated properly")
				}
			}

			for _, plate := range props.Plates {
				plate.ID = "Changed_" + plate.ID
			}
			for i, plate := range dprops.Plates {
				if plate.ID == props.Plates[i].ID {
					return errors.New("props.Plates not duplicated properly")
				}
			}

			for _, tipbox := range props.Tipboxes {
				tipbox.ID = "changed_" + tipbox.ID
			}
			for i, tipbox := range dprops.Tipboxes {
				if tipbox.ID == props.Tipboxes[i].ID {
					return errors.New("props.Tipboxes not duplicated properly")
				}
			}
			return nil
		},
	})
}

func AssertLHPropertiesEqual(t *testing.T, e, g *liquidhandling.LHProperties, msg string) {
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
	testlab.WithTestLab(t, "", &testlab.TestElementCallbacks{
		Steps: func(lab *laboratory.Laboratory) error {
			before := MakeGilsonWithPlatesAndTipboxesForTest(lab, "")

			// we don't need to preserve this
			for _, tip := range before.Tips {
				tip.ClearParent()
			}

			var after liquidhandling.LHProperties
			if data, err := json.Marshal(before); err != nil {
				return err
			} else if err := json.Unmarshal(data, &after); err != nil {
				return err
			}

			if !reflect.DeepEqual(before, &after) {
				return fmt.Errorf("serialization changed LHProperties PlateLookup\nbefore: %+v\nafter : %+v", before, &after)
			}

			heads := make(map[*wtype.LHHead]bool)
			for _, h := range after.Heads {
				heads[h] = true
			}
			for _, ha := range after.HeadAssemblies {
				for _, hap := range ha.Positions {
					if hap.Head != nil && !heads[hap.Head] {
						return errors.New("HeadAssemblyPosition.Head doesn't point to anything in LHProperties.Heads")
					}
				}
			}

			adaptors := make(map[*wtype.LHAdaptor]bool)
			for _, a := range after.Adaptors {
				adaptors[a] = true
			}
			for _, head := range after.Heads {
				if head.Adaptor != nil && !adaptors[head.Adaptor] {
					return errors.New("Head.Adaptor doesn't point to anything in LHProperties.Adaptors")
				}
			}

			if b, a := before.GetLoadedHeads(), after.GetLoadedHeads(); len(b) != len(a) {
				return fmt.Errorf("number of loaded heads doesn't match: before = %d, after = %d", len(b), len(a))
			} else {
				for i, beforeHead := range b {
					if afterHead := a[i]; !reflect.DeepEqual(beforeHead, afterHead) {
						return fmt.Errorf("%dth head is mismatched\nbefore: %+v\nafter : %+v", i, beforeHead, afterHead)
					}
				}
			}
			return nil
		},
	})
}

type UpdateIDTest struct {
	Name      string
	UpdateIDs []string // list of IDs to update
	ErrorIDs  []string // ids that should produce an error
}

func (test *UpdateIDTest) Run(t *testing.T) {
	testlab.WithTestLab(t, "", &testlab.TestElementCallbacks{
		Name: test.Name,
		Steps: func(lab *laboratory.Laboratory) error {

			// this test is quite involved, mainly in order to make sure that wastes
			// and washes aren't lost, as could sometimes be the case with the previous implementation
			//
			// strategy is to call UpdateID on a copy, then apply the ID mapping manually to
			// the original and assert that everything matches

			// build some properties to test with
			rbt := MakeGilsonForTest(lab, defaultTipList())

			// add something for each object type
			if plt, err := lab.Inventory.Plates.NewPlate("pcrplate_skirted"); err != nil {
				return err
			} else {
				plt.ID = "initial_plate_id"
				if err := rbt.AddInputPlate(plt); err != nil {
					return err
				}
			}
			if plt, err := lab.Inventory.Plates.NewPlate("pcrplate_skirted"); err != nil {
				return err
			} else {
				plt.ID = "initial_wash_id"
				rbt.AddWash(plt)
			}
			if plt, err := lab.Inventory.Plates.NewPlate("pcrplate_skirted"); err != nil {
				return err
			} else {
				plt.ID = "initial_waste_id"
				rbt.AddWaste(plt)
			}
			if tb, err := lab.Inventory.TipBoxes.NewTipbox("Gilson20"); err != nil {
				return err
			} else {
				tb.ID = "initial_tipbox_id"
				if err := rbt.AddTipBox(tb); err != nil {
					return err
				}
			}
			if tw, err := lab.Inventory.TipWastes.NewTipwaste("Gilsontipwaste"); err != nil {
				return err
			} else {
				tw.ID = "initial_tipwaste_id"
				if err := rbt.AddTipWaste(tw); err != nil {
					return err
				}
			}

			// update the IDs as required

			beforeToAfter := make(map[string]string, len(test.UpdateIDs))
			for id := range rbt.PlateLookup {
				beforeToAfter[id] = id
			}

			final := rbt.DupKeepIDs(lab.IDGenerator)
			errorIDs := make([]string, 0, len(test.UpdateIDs))
			for _, id := range test.UpdateIDs {
				if newID, err := final.UpdateID(lab.IDGenerator, id); err != nil {
					errorIDs = append(errorIDs, id)
				} else {
					beforeToAfter[id] = newID
				}
			}
			if !assert.ElementsMatch(t, test.ErrorIDs, errorIDs) || len(test.ErrorIDs) != 0 {
				// don't continue testing if we were expecting an error
				return nil
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
			return nil
		},
	})
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
