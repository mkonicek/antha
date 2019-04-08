package wtype

import (
	"testing"

	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/laboratory/effects/id"
)

func getTestWell(idGen *id.IDGenerator, volUL, residualUL float64) *LHWell {
	shp := NewShape(BoxShape, "mm", 10.0, 5.0, 5.0)
	well := NewLHWell(idGen, "ul", volUL, residualUL, shp, UWellBottom, 5.0, 5.0, 10.0, 1.0, "mm")
	return well
}

func getTestComponent(idGen *id.IDGenerator, volUL float64) *Liquid {
	cmp := NewLHComponent(idGen)
	cmp.Vol = volUL
	return cmp
}

func TestWellVolumes(t *testing.T) {
	idGen := id.NewIDGenerator("testing")
	well := getTestWell(idGen, 100.0, 1.0)

	if e, g := 100.0, well.MaxVolume().ConvertToString("ul"); e != g {
		t.Errorf("LHWell.MaxVolume() returned %f ul, expected %f ul", g, e)
	}

	if e, g := 1.0, well.ResidualVolume().ConvertToString("ul"); e != g {
		t.Errorf("LHWell.ResidualVolume() returned %f ul, expected %f ul", g, e)
	}

	if e, g := 99.0, well.MaxWorkingVolume().ConvertToString("ul"); e != g {
		t.Errorf("LHWell.MaxWorkingVolume() returned %f ul, expected %f ul", g, e)
	}

	if e, g := 0.0, well.CurrentVolume(idGen).ConvertToString("ul"); e != g {
		t.Errorf("LHWell.MaxVolume() returned %f ul, expected %f ul", g, e)
	}

	if e, g := 0.0, well.CurrentWorkingVolume(idGen).ConvertToString("ul"); e != g {
		t.Errorf("LHWell.MaxVolume() returned %f ul, expected %f ul", g, e)
	}

	if !well.IsEmpty(idGen) {
		t.Error("newly created well IsEmpty() returned false")
	}

	//let's add a small amount
	cmp := getTestComponent(idGen, 0.5)
	if err := well.AddComponent(idGen, cmp); err != nil {
		t.Fatal(err)
	}

	if e, g := 0.5, well.CurrentVolume(idGen).ConvertToString("ul"); e != g {
		t.Errorf("LHWell.MaxVolume() returned %f ul, expected %f ul", g, e)
	}

	//we can't remove the volume that we just added
	if e, g := 0.0, well.CurrentWorkingVolume(idGen).ConvertToString("ul"); e != g {
		t.Errorf("LHWell.MaxVolume() returned %f ul, expected %f ul", g, e)
	}

	if well.IsEmpty(idGen) {
		t.Error("successfully added a component to a well, but it still IsEmpty()")
	}

	//let's add some more
	cmp2 := getTestComponent(idGen, 50.0)
	if err := well.AddComponent(idGen, cmp2); err != nil {
		t.Fatal(err)
	}

	if e, g := 50.5, well.CurrentVolume(idGen).ConvertToString("ul"); e != g {
		t.Errorf("LHWell.MaxVolume() returned %f ul, expected %f ul", g, e)
	}

	if e, g := 49.5, well.CurrentWorkingVolume(idGen).ConvertToString("ul"); e != g {
		t.Errorf("LHWell.MaxVolume() returned %f ul, expected %f ul", g, e)
	}
}

func TestAddComponentOK(t *testing.T) {
	idGen := id.NewIDGenerator("testing")
	well := getTestWell(idGen, 100.0, 1.0)
	cmp := getTestComponent(idGen, 5.0)

	if err := well.AddComponent(idGen, cmp); err != nil {
		t.Error(err)
	}
}

func TestAddComponentOverfilled(t *testing.T) {
	//Skipping because AddComponent doesn't raise errors at the moment
	//due to CarryVolume issues.
	t.Skip()
	idGen := id.NewIDGenerator("testing")

	well := getTestWell(idGen, 100.0, 1.0)
	cmp := getTestComponent(idGen, 100.5)

	if err := well.AddComponent(idGen, cmp); err == nil {
		t.Error("added 100.5ul component to a 100ul well and got no error")
	}
}

func TestRemoveVolume(t *testing.T) {
	idGen := id.NewIDGenerator("testing")
	well := getTestWell(idGen, 100.0, 1.0)
	cmp := getTestComponent(idGen, 50.0)

	if err := well.AddComponent(idGen, cmp); err != nil {
		t.Fatal(err)
	}

	cmp2, err := well.RemoveVolume(idGen, wunit.NewVolume(20.0, "ul"))
	if err != nil {
		t.Fatal(err)
	}

	if e, g := 20.0, cmp2.Volume().ConvertToString("ul"); e != g {
		t.Errorf("component volume was %f ul, expected %f", g, e)
	}
	//Skipping because removeVolume doesn't raise errors at the moment
	//due to CarryVolume issues.
	t.Skip()

	workingVol := well.CurrentWorkingVolume(idGen)
	cmp3, err := well.RemoveVolume(idGen, wunit.NewVolume(30.0, "ul"))
	if err == nil {
		t.Fatalf("removed 30ul from a well with %v working volume without error", workingVol)
	}

	if cmp3 != nil {
		t.Errorf("component should be nil, but got %v", cmp3)
	}

}

func TestWellValidation(t *testing.T) {
	idGen := id.NewIDGenerator("testing")
	cmp := getTestComponent(idGen, 50.0)

	well := &LHWell{
		WContents: cmp,
		MaxVol:    100.0,
	}

	if !well.IsVolumeValid(idGen) {
		t.Errorf("well.IsVolumeValid() returned false : CurrentVolume(), MaxVolume() = %v, %v", well.CurrentVolume(idGen), well.MaxVolume())
	}

	if err := well.ValidateVolume(idGen); err != nil {
		t.Error(err)
	}

	well = &LHWell{
		WContents: cmp,
		MaxVol:    10.0,
	}

	if well.IsVolumeValid(idGen) {
		t.Errorf("well.IsVolumeValid() returned true : CurrentVolume(), MaxVolume() = %v, %v", well.CurrentVolume(idGen), well.MaxVolume())
	}

	if err := well.ValidateVolume(idGen); err == nil {
		t.Errorf("well.ValidateVolume() returned no error : CurrentVolume(), MaxVolume() = %v, %v", well.CurrentVolume(idGen), well.MaxVolume())
	}

}

func TestGetNextWellOk(t *testing.T) {
	idGen := id.NewIDGenerator("testing")
	trough := maketroughfortest(idGen)
	var well *LHWell

	//half working vol - carry vol approximation
	//will fit 2 components in each well since they're the same type
	testVol := 4990.0

	for i := 0; i < 24; i++ {
		cmp := getTestComponent(idGen, testVol)
		well, ok := Get_Next_Well(idGen, trough, cmp, well)
		if !ok {
			t.Fatal("got no well when trough wasn't full")
		}

		if g, e := well.GetWellCoords(), (WellCoords{X: i / 2, Y: 0}); !e.Equals(g) {
			t.Errorf("wellcoords don't match: expected %s, got %s", e.FormatA1(), g.FormatA1())
		}

		if err := well.AddComponent(idGen, cmp); err != nil {
			t.Fatal(err)
		}
	}

	//now that the plate's full, should only return nil, false
	cmp := getTestComponent(idGen, testVol)
	well, ok := Get_Next_Well(idGen, trough, cmp, well)
	if well != nil || ok {
		t.Errorf("plate full: expected output (nil, false), got (%s, %t)", NameOf(trough), ok)
	}

	well, ok = Get_Next_Well(idGen, trough, cmp, nil)
	if well != nil || ok {
		t.Errorf("plate full: expected output (nil, false), got (%s, %t)", NameOf(trough), ok)
	}

}
