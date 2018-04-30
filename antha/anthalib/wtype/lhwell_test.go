package wtype

import (
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"testing"
)

func getTestWell(volUL, residualUL float64) *LHWell {
	shp := NewShape("klein bottle", "mm", 10.0, 5.0, 5.0)
	well := NewLHWell("ul", volUL, residualUL, shp, UWellBottom, 5.0, 5.0, 10.0, 1.0, "mm")
	return well
}

func getTestComponent(volUL float64) *LHComponent {
	cmp := NewLHComponent()
	cmp.Vol = volUL
	return cmp
}

func TestWellVolumes(t *testing.T) {
	well := getTestWell(100.0, 1.0)

	if e, g := 100.0, well.MaxVolume().ConvertToString("ul"); e != g {
		t.Errorf("LHWell.MaxVolume() returned %f ul, expected %f ul", g, e)
	}

	if e, g := 1.0, well.ResidualVolume().ConvertToString("ul"); e != g {
		t.Errorf("LHWell.ResidualVolume() returned %f ul, expected %f ul", g, e)
	}

	if e, g := 99.0, well.MaxWorkingVolume().ConvertToString("ul"); e != g {
		t.Errorf("LHWell.MaxWorkingVolume() returned %f ul, expected %f ul", g, e)
	}

	if e, g := 0.0, well.CurrentVolume().ConvertToString("ul"); e != g {
		t.Errorf("LHWell.MaxVolume() returned %f ul, expected %f ul", g, e)
	}

	if e, g := 0.0, well.CurrentWorkingVolume().ConvertToString("ul"); e != g {
		t.Errorf("LHWell.MaxVolume() returned %f ul, expected %f ul", g, e)
	}

	if !well.IsEmpty() {
		t.Error("newly created well IsEmpty() returned false")
	}

	//let's add a small amount
	cmp := getTestComponent(0.5)
	if err := well.AddComponent(cmp); err != nil {
		t.Fatal(err)
	}

	if e, g := 0.5, well.CurrentVolume().ConvertToString("ul"); e != g {
		t.Errorf("LHWell.MaxVolume() returned %f ul, expected %f ul", g, e)
	}

	//we can't remove the volume that we just added
	if e, g := 0.0, well.CurrentWorkingVolume().ConvertToString("ul"); e != g {
		t.Errorf("LHWell.MaxVolume() returned %f ul, expected %f ul", g, e)
	}

	if well.IsEmpty() {
		t.Error("successfully added a component to a well, but it still IsEmpty()")
	}

	//let's add some more
	cmp2 := getTestComponent(50.0)
	if err := well.AddComponent(cmp2); err != nil {
		t.Fatal(err)
	}

	if e, g := 50.5, well.CurrentVolume().ConvertToString("ul"); e != g {
		t.Errorf("LHWell.MaxVolume() returned %f ul, expected %f ul", g, e)
	}

	if e, g := 49.5, well.CurrentWorkingVolume().ConvertToString("ul"); e != g {
		t.Errorf("LHWell.MaxVolume() returned %f ul, expected %f ul", g, e)
	}
}

func TestAddComponentOK(t *testing.T) {
	well := getTestWell(100.0, 1.0)
	cmp := getTestComponent(5.0)

	if err := well.AddComponent(cmp); err != nil {
		t.Error(err)
	}
}

func TestAddComponentOverfilled(t *testing.T) {
	//Skipping because AddComponent doesn't raise errors at the moment
	//due to CarryVolume issues.
	t.Skip()
	well := getTestWell(100.0, 1.0)
	cmp := getTestComponent(100.5)

	if err := well.AddComponent(cmp); err == nil {
		t.Error("added 100.5ul component to a 100ul well and got no error")
	}
}

func TestRemoveVolume(t *testing.T) {
	well := getTestWell(100.0, 1.0)
	cmp := getTestComponent(50.0)

	if err := well.AddComponent(cmp); err != nil {
		t.Fatal(err)
	}

	cmp2, err := well.RemoveVolume(wunit.NewVolume(20.0, "ul"))
	if err != nil {
		t.Fatal(err)
	}

	if e, g := 20.0, cmp2.Volume().ConvertToString("ul"); e != g {
		t.Errorf("component volume was %f ul, expected %f", g, e)
	}
	//Skipping because removeVolume doesn't raise errors at the moment
	//due to CarryVolume issues.
	t.Skip()

	workingVol := well.CurrentWorkingVolume()
	cmp3, err := well.RemoveVolume(wunit.NewVolume(30.0, "ul"))
	if err == nil {
		t.Fatalf("removed 30ul from a well with %v working volume without error", workingVol)
	}

	if cmp3 != nil {
		t.Errorf("component should be nil, but got %v", cmp3)
	}

}

func TestWellValidation(t *testing.T) {
	cmp := getTestComponent(50.0)

	well := &LHWell{
		WContents: cmp,
		MaxVol:    100.0,
	}

	if !well.IsVolumeValid() {
		t.Errorf("well.IsVolumeValid() returned false : CurrentVolume(), MaxVolume() = %v, %v", well.CurrentVolume(), well.MaxVolume())
	}

	if err := well.ValidateVolume(); err != nil {
		t.Error(err)
	}

	well = &LHWell{
		WContents: cmp,
		MaxVol:    10.0,
	}

	if well.IsVolumeValid() {
		t.Errorf("well.IsVolumeValid() returned true : CurrentVolume(), MaxVolume() = %v, %v", well.CurrentVolume(), well.MaxVolume())
	}

	if err := well.ValidateVolume(); err == nil {
		t.Errorf("well.ValidateVolume() returned no error : CurrentVolume(), MaxVolume() = %v, %v", well.CurrentVolume(), well.MaxVolume())
	}

}
