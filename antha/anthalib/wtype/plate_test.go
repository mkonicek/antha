package wtype

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"testing"
)

// platetype, mfr string, nrows, ncols int, height float64, hunit string, welltype *LHWell, wellXOffset, wellYOffset, wellXStart, wellYStart, wellZStart float64

func makeplatefortest() *Plate {
	swshp := NewShape(BoxShape, "mm", 8.2, 8.2, 41.3)
	welltype := NewLHWell("ul", 200, 10, swshp, VWellBottom, 8.2, 8.2, 41.3, 4.7, "mm")
	p := NewLHPlate("DSW96", "none", 8, 12, Coordinates3D{127.76, 85.48, 43.1}, welltype, 9.0, 9.0, 0.5, 0.5, 0.5)
	return p
}

/* -- these aren't used, but might be useful again in the future?
func make384platefortest() *LHPlate {
	swshp := NewShape(BoxShape, "mm", 8.2, 8.2, 41.3)
	welltype := NewLHWell("ul", 50, 5, swshp, VWellBottom, 8.2, 8.2, 41.3, 4.7, "mm")
	p := NewLHPlate("DSW384", "none", 16, 24, Coordinates{127.76, 85.48, 44.1}, welltype, 0.5, 0.5, 0.5, 0.5, 0.5)
	return p
}

func make1536platefortest() *LHPlate {
	swshp := NewShape(BoxShape, "mm", 8.2, 8.2, 41.3)
	welltype := NewLHWell("ul", 15, 1, swshp, VWellBottom, 8.2, 8.2, 41.3, 4.7, "mm")
	p := NewLHPlate("DSW1536", "none", 32, 48, Coordinates{127.76, 85.48, 44.1}, welltype, 0.5, 0.5, 0.5, 0.5, 0.5)
	return p
}

func make24platefortest() *LHPlate {
	swshp := NewShape(BoxShape, "mm", 8.2, 8.2, 41.3)
	welltype := NewLHWell("ul", 3000, 500, swshp, VWellBottom, 8.2, 8.2, 41.3, 4.7, "mm")
	p := NewLHPlate("DSW24", "none", 4, 6, Coordinates{127.76, 85.48, 44.1}, welltype, 0.5, 0.5, 0.5, 0.5, 0.5)
	return p
}

func make6platefortest() *LHPlate {
	swshp := NewShape(BoxShape, "mm", 8.2, 8.2, 41.3)
	welltype := NewLHWell("ul", 3000, 500, swshp, VWellBottom, 8.2, 8.2, 41.3, 4.7, "mm")
	p := NewLHPlate("6wellplate", "none", 2, 3, Coordinates{127.76, 85.48, 44.1}, welltype, 0.5, 0.5, 0.5, 0.5, 0.5)
	return p
}
*/

func maketroughfortest() *Plate {
	stshp := NewShape(BoxShape, "mm", 8.2, 72, 41.3)
	trough12 := NewLHWell("ul", 15000, 5000, stshp, VWellBottom, 8.2, 72, 41.3, 4.7, "mm")
	plate := NewLHPlate("DWST12", "Unknown", 1, 12, Coordinates3D{127.76, 85.48, 44.1}, trough12, 9, 9, 0, 30.0, 4.5)
	return plate
}

func TestPlateCreation(t *testing.T) {
	p := makeplatefortest()
	validatePlate(t, p)
}

func TestPlateDup(t *testing.T) {
	p := makeplatefortest()
	d := p.Dup()
	validatePlate(t, d)
	for crds, w := range p.Wellcoords {
		w2 := d.Wellcoords[crds]

		if w.ID == w2.ID {
			t.Fatal(fmt.Sprintf("Error: coords %s has same IDs before / after dup", crds))
		}

		if w.WContents.Loc == w2.WContents.Loc {
			t.Fatal(fmt.Sprintf("Error: contents of wells at coords %s have same loc before and after regular Dup()", crds))
		}
	}
}

func TestPlateDupKeepIDs(t *testing.T) {
	p := makeplatefortest()
	d := p.DupKeepIDs()

	for crds, w := range p.Wellcoords {
		w2 := d.Wellcoords[crds]

		if w.ID != w2.ID {
			t.Fatal(fmt.Sprintf("Error: coords %s has different IDs", crds))
		}

		if w.WContents.ID != w2.WContents.ID {
			t.Fatal(fmt.Sprintf("Error: contents of wells at coords %s have different IDs", crds))

		}
		if w.WContents.Loc != w2.WContents.Loc {
			t.Fatal(fmt.Sprintf("Error: contents of wells at coords %s have different loc before and after DupKeepIDs()", crds))
		}
	}

}

func validatePlate(t *testing.T, plate *Plate) {
	assertWellsEqual := func(what string, as, bs []*LHWell) {
		seen := make(map[*LHWell]int)
		for _, w := range as {
			seen[w] += 1
		}
		for _, w := range bs {
			seen[w] += 1
		}
		for w, count := range seen {
			if count != 2 {
				t.Errorf("%s: no matching well found (%d != %d) for %p %s:%s", what, count, 2, w, w.ID, w.Crds.FormatA1())
			}
		}
	}

	var ws1, ws2, ws3, ws4 []*LHWell

	for _, w := range plate.HWells {
		ws1 = append(ws1, w)
	}
	for crds, w := range plate.Wellcoords {
		ws2 = append(ws2, w)

		if w.Crds.FormatA1() != crds {
			t.Fatal(fmt.Sprintf("ERROR: Well coords not consistent -- %s != %s", w.Crds.FormatA1(), crds))
		}

		if w.WContents.Loc == "" {
			t.Fatal(fmt.Sprintf("ERROR: Well contents do not have loc set"))
		}

		ltx := strings.Split(w.WContents.Loc, ":")

		if ltx[0] != plate.ID {
			t.Fatal(fmt.Sprintf("ERROR: Plate ID for component not consistent -- %s != %s", ltx[0], plate.ID))
		}

		if ltx[0] != w.Plate.(*Plate).ID {
			t.Fatal(fmt.Sprintf("ERROR: Plate ID for component not consistent with well -- %s != %s", ltx[0], w.Plate.(*Plate).ID))
		}

		if ltx[1] != crds {
			t.Fatal(fmt.Sprintf("ERROR: Coords for component not consistent: -- %s != %s", ltx[1], crds))
		}

	}

	for _, ws := range plate.Rows {
		ws3 = append(ws3, ws...)
	}
	for _, ws := range plate.Cols {
		ws4 = append(ws4, ws...)
	}
	assertWellsEqual("HWells != Rows", ws1, ws2)
	assertWellsEqual("Rows != Cols", ws2, ws3)
	assertWellsEqual("Cols != Wellcoords", ws3, ws4)

	// Check pointer-ID equality
	comp := make(map[string]*Liquid)
	for _, w := range append(append(ws1, ws2...), ws3...) {
		c := w.WContents
		if c == nil || c.Vol == 0.0 {
			continue
		}
		if co, seen := comp[c.ID]; seen && co != c {
			t.Errorf("component %s duplicated as %+v and %+v", c.ID, c, co)
		} else if !seen {
			comp[c.ID] = c
		}
	}
}

func TestIsUserAllocated(t *testing.T) {
	p := makeplatefortest()

	if p.IsUserAllocated() {
		t.Fatal("Error: Plates must not start out user allocated")
	}
	p.Wellcoords["A1"].SetUserAllocated()

	if !p.IsUserAllocated() {
		t.Fatal("Error: Plates with at least one user allocated well must return true to IsUserAllocated()")
	}

	d := p.Dup()

	if !d.IsUserAllocated() {
		t.Fatal("Error: user allocation mark must survive Dup()lication")
	}

	d.Wellcoords["A1"].ClearUserAllocated()

	if d.IsUserAllocated() {
		t.Fatal("Error: user allocation mark not cleared")
	}

	if !p.IsUserAllocated() {
		t.Fatal("Error: UserAllocation mark must operate separately on Dup()licated plates")
	}
}

func TestMergeWith(t *testing.T) {
	p1 := makeplatefortest()
	p2 := makeplatefortest()

	c := NewLHComponent()

	c.CName = "Water1"
	c.Vol = 50.0
	c.Vunit = "ul"
	err := p1.Wellcoords["A1"].AddComponent(c)
	if err != nil {
		t.Fatal(err)
	}
	p1.Wellcoords["A1"].SetUserAllocated()

	c = NewLHComponent()
	c.CName = "Butter"
	c.Vol = 80.0
	c.Vunit = "ul"
	err = p2.Wellcoords["A2"].AddComponent(c)
	if err != nil {
		t.Fatal(err)
	}

	p1.MergeWith(p2)

	if !(p1.Wellcoords["A1"].WContents.CName == "Water1" && p1.Wellcoords["A1"].WContents.Vol == 50.0 && p1.Wellcoords["A1"].WContents.Vunit == "ul") {
		t.Fatal("Error: MergeWith should leave user allocated components alone")
	}

	if !(p1.Wellcoords["A2"].WContents.CName == "Butter" && p1.Wellcoords["A2"].WContents.Vol == 80.0 && p1.Wellcoords["A2"].WContents.Vunit == "ul") {
		t.Fatal("Error: MergeWith should add non user-allocated components to  plate merged with")
	}
}

func TestLHPlateSerialize(t *testing.T) {
	p := makeplatefortest()
	c := NewLHComponent()
	c.CName = "Cthulhu"
	c.Type = LTWater
	c.Vol = 100.0

	_, err := p.AddComponent(c, false)
	if err != nil {
		t.Errorf(err.Error())
	}

	b, err := json.Marshal(p)
	if err != nil {
		t.Errorf(err.Error())
	}

	var p2 *Plate

	if err = json.Unmarshal(b, &p2); err != nil {
		t.Errorf(err.Error())
	}

	for i, w := range p.Wellcoords {
		w2 := p2.Wellcoords[i]

		if !reflect.DeepEqual(w.WContents, w2.WContents) {
			t.Errorf("%v =/= %v", w.WContents, w2.WContents)
		}

		if w2.Plate != p2 {
			t.Errorf("Wells not retaining plate references post serialization")
		}

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
			wc := WellCoords{X: i, Y: j}

			w := p2.Wellcoords[wc.FormatA1()]

			w.WContents.CName = wc.FormatA1()
			if p2.Rows[j][i].WContents.CName != wc.FormatA1() || p2.Cols[i][j].WContents.CName != wc.FormatA1() || p2.HWells[w.ID].WContents.CName != wc.FormatA1() {
				fmt.Println(p2.Cols[i][j].WContents.CName)
				fmt.Println(p2.Rows[j][i].WContents.CName)
				t.Errorf("Error: Wells inconsistent at position %s", wc.FormatA1())
			}

		}
	}

	// check extraneous parameters

	if p.ID != p2.ID {
		t.Errorf(fMErr("ID", p.ID, p2.ID))
	}

	if p.PlateName != p2.PlateName {
		t.Errorf(fMErr("Plate name", p.PlateName, p2.PlateName))
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
		t.Errorf(fMErr("WellZStart", p.WellZStart, p.WellZStart))
	}
}

func TestAddGetClearData(t *testing.T) {
	dat := []byte("3.5")

	t.Run("basic", func(t *testing.T) {
		p := makeplatefortest()

		if err := p.SetData("OD", dat); err != nil {
			t.Errorf(err.Error())
		}
		d, err := p.GetData("OD")
		if err != nil {
			t.Errorf(err.Error())
		}
		if !reflect.DeepEqual(d, dat) {
			t.Errorf("Expected %v got %v", dat, d)
		}
	})

	t.Run("clear", func(t *testing.T) {
		p := makeplatefortest()

		if err := p.SetData("OD", dat); err != nil {
			t.Errorf(err.Error())
		}

		if err := p.ClearData("OD"); err != nil {
			t.Errorf(err.Error())
		}

		if _, err := p.GetData("OD"); err == nil {
			t.Errorf("ClearData should clear data but has not")
		}
	})

	t.Run("cannot update special", func(t *testing.T) {
		p := makeplatefortest()
		if err := p.SetData("IMSPECIAL", dat); err == nil {
			t.Errorf("Adding data with a reserved key should fail but does not")
		}
	})

}

func TestGetAllComponents(t *testing.T) {
	p := makeplatefortest()

	cmps := p.AllContents()

	if len(cmps) != p.WellsX()*p.WellsY() {
		t.Errorf("Expected %d components got %d", p.WellsX()*p.WellsY(), len(cmps))
	}
}

func TestLHPlateValidateVolumesOK(t *testing.T) {
	p := makeplatefortest()
	c := NewLHComponent()
	c.CName = "Cthulhu"
	c.Type = LTWater
	c.Vol = 100.0

	if _, err := p.AddComponent(c, false); err != nil {
		t.Errorf(err.Error())
	}

	if err := p.ValidateVolumes(); err != nil {
		t.Error(err)
	}
}

func TestLHPlateValidateVolumesOneOverfilled(t *testing.T) {
	p := makeplatefortest()
	c := NewLHComponent()
	c.CName = "Cthulhu"
	c.Type = LTWater
	c.Vol = 100.0

	if _, err := p.AddComponent(c, false); err != nil {
		t.Errorf(err.Error())
	}

	//doing it this way because accessor methods will prevent this at some point
	c.Vol = 500.0
	w := p.Rows[0][0]
	w.WContents = c

	if err := p.ValidateVolumes(); err == nil {
		t.Error("Got no error when one well overfilled")
	}
}

func TestLHPlateValidateVolumesSeveralOverfilled(t *testing.T) {
	p := makeplatefortest()
	c := NewLHComponent()
	c.CName = "Cthulhu"
	c.Type = LTWater
	c.Vol = 100.0

	if _, err := p.AddComponent(c, false); err != nil {
		t.Errorf(err.Error())
	}

	//doing it this way because accessor methods will prevent this at some point
	c.Vol = 500.0
	for i := 0; i < 4; i++ {
		w := p.Rows[i][i]
		w.WContents = c
	}

	if err := p.ValidateVolumes(); err == nil {
		t.Error("Got no error when several wells overfilled")
	}
}

func TestSpecialRetention(t *testing.T) {
	p := makeplatefortest()

	p.DeclareSpecial()

	// dup must do this
	d := p.Dup()

	if !d.IsSpecial() {
		t.Error("Duplicated plates must retain specialness")
	}

	// so must serialization

	dat, err := json.Marshal(d)

	if err != nil {
		t.Errorf("Marshal error: %v", err)
	}

	var e *Plate

	err = json.Unmarshal(dat, &e)

	if err != nil {
		t.Errorf("Unmarshal error: %v", err)
	}

	if !e.IsSpecial() {
		t.Error("Specialness must be retained after serialize/deserialize")
	}

	// and cleaning

	e.Clean()

	if !e.IsSpecial() {
		t.Error("Specialness must be retained after cleaning")
	}

}

func TestWellCoordsToCoords(t *testing.T) {

	plate := makeplatefortest()
	c := NewLHComponent()
	c.Vol = 100.0
	c.Vunit = "ul"
	if err := plate.GetChildByAddress(MakeWellCoords("A1")).(*LHWell).AddComponent(c); err != nil {
		t.Fatal(err)
	}

	type TestCase struct {
		Address          string
		Reference        WellReference
		ExpectedPosition Coordinates3D
		ExpectingError   bool
	}

	tests := []TestCase{
		{
			Address:          "A1",
			Reference:        BottomReference,
			ExpectedPosition: Coordinates3D{X: plate.WellXStart, Y: plate.WellYStart, Z: plate.WellZStart + plate.Welltype.Bottomh},
		},
		{
			Address:          "A1",
			Reference:        TopReference,
			ExpectedPosition: Coordinates3D{X: plate.WellXStart, Y: plate.WellYStart, Z: plate.WellZStart + plate.Welltype.GetSize().Z},
		},
		{
			Address:          "A1",
			Reference:        LiquidReference,
			ExpectedPosition: Coordinates3D{X: plate.WellXStart, Y: plate.WellYStart, Z: plate.WellZStart + 0.5*(plate.Welltype.Bottomh+plate.Welltype.GetSize().Z)},
		},
		{
			Address:        "Z1",
			Reference:      TopReference,
			ExpectingError: true,
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("%s_%v", test.Address, test.Reference), func(t *testing.T) {
			pos, ok := plate.WellCoordsToCoords(MakeWellCoords(test.Address), test.Reference)
			if !ok != test.ExpectingError {
				t.Fatalf("expecting error: %t, got error: %t", test.ExpectingError, !ok)
			}

			if !test.ExpectingError && !test.ExpectedPosition.Equals(pos) {
				t.Errorf("position was wrong: expected %v got %v", test.ExpectedPosition, pos)
			}
		})
	}

}

func TestCoordsToWellCoords(t *testing.T) {

	plate := makeplatefortest()

	pos := Coordinates3D{
		X: plate.WellXStart + 0.75*plate.WellXOffset,
		Y: plate.WellYStart + 0.75*plate.WellYOffset,
	}

	wc, delta := plate.CoordsToWellCoords(pos)

	if e, g := "B2", wc.FormatA1(); e != g {
		t.Errorf("Wrong well coordinates: expected %s, got %s", e, g)
	}

	eDelta := -0.25 * plate.WellXOffset
	if delta.X != eDelta || delta.Y != eDelta {
		t.Errorf("Delta incorrect: expected (%f, %f), got (%f, %f)", eDelta, eDelta, delta.X, delta.Y)
	}

}

func TestGetWellBounds(t *testing.T) {

	plate := makeplatefortest()

	eStart := Coordinates3D{
		X: 0.5 - 0.5*8.2,
		Y: 0.5 - 0.5*8.2,
		Z: 0.5,
	}
	eSize := Coordinates3D{
		X: 9.0*11 + 8.2,
		Y: 9.0*7 + 8.2,
		Z: 41.3,
	}
	eBounds := NewBBox(eStart, eSize)
	bounds := plate.GetWellBounds()

	if e, g := eBounds.String(), bounds.String(); e != g {
		t.Errorf("GetWellBounds incorrect: expected %v, got %v", eBounds, bounds)
	}
}
