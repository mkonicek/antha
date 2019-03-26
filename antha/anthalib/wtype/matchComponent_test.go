package wtype

import (
	"reflect"
	"testing"

	"github.com/antha-lang/antha/antha/anthalib/wunit"
)

func TestMatchComponent(t *testing.T) {
	c := NewLHComponent()
	c.CName = "water"
	c.Vol = 200.0
	CIDs := []string{"A1", "B1", "C1", "D1", "E1", "F1", "G1", "H1"}
	PIDs := []string{"Plate1", "Plate1", "Plate1", "Plate1", "Plate1", "Plate1", "Plate1", "Plate1"}

	got := make([]*Liquid, 8)

	for i := 0; i < 8; i++ {
		got[i] = c.Dup()
		got[i].Loc = PIDs[i] + ":" + CIDs[i]
	}

	d := NewLHComponent()
	d.CName = "water"
	d.Vol = 20.0

	CID2s := []string{"A1", "B1", "C1", "D1", "E1", "F1", "G1", "H1"}
	PID2s := []string{"Plate2", "Plate2", "Plate2", "Plate2", "Plate2", "Plate2", "Plate2", "Plate2"}

	want := make([]*Liquid, 8)
	for i := 0; i < 8; i++ {
		want[i] = d.Dup()
		want[i].Loc = PID2s[i] + ":" + CID2s[i]
	}

	match, err := MatchComponents(want, got, false, false)

	if err != nil {
		t.Errorf(err.Error())
	}

	expected := Match{
		IDs:  PIDs,
		WCs:  CIDs,
		Vols: toVolArr(20.0, 8),
		M:    seq(0, 8, 1),
		Sc:   160.0,
	}

	if !reflect.DeepEqual(expected, match) {
		t.Errorf("%v =/= %v", match, expected)
	}
}

func toVolArr(v float64, n int) []wunit.Volume {
	ret := make([]wunit.Volume, n)

	for i := 0; i < n; i++ {
		ret[i] = wunit.NewVolume(v, "ul")
	}
	return ret
}

func seq(start, length, increment int) []int {
	r := make([]int, length)
	x := start
	for i := 0; i < length; i++ {
		r[i] = x
		x += increment
	}

	return r
}

func updateSrcs(m Match, ca []*Liquid) {
	for i, v := range m.Vols {
		if m.M[i] == -1 {
			continue
		}
		ca[m.M[i]].Vol -= v.RawValue()
		if ca[m.M[i]].Vol < 0.0 {
			ca[m.M[i]].Vol = 0.0
		}
	}
}
func updateDsts(m Match, ca []*Liquid) {
	for i, v := range m.Vols {
		if m.M[i] == -1 {
			continue
		}
		ca[i].Vol -= v.RawValue()
		if ca[i].Vol < 0.0 {
			ca[i].Vol = 0.0
		}
	}
}

func dstsDone(ca []*Liquid) bool {
	for _, c := range ca {
		if c.Vol > 0.0 {
			return false
		}
	}
	return true
}

// presently does not enforce same volume
// per tip rule... how was I doing that in the past?
// I think it's enforced on the level above... we will have to
// revise to permit multi

func TestMatchComponentPickupVolumes(t *testing.T) {
	c := NewLHComponent()
	c.CName = "water"
	vls := []float64{100.0, 100.0, 5.0, 5.0}
	CIDs := []string{"A1", "B1", "C1", "D1"}
	PIDs := []string{"Plate1", "Plate1", "Plate1", "Plate1"}

	got := make([]*Liquid, 4)

	for i := 0; i < 4; i++ {
		got[i] = c.Dup()
		got[i].Vol = vls[i]
		got[i].Loc = PIDs[i] + ":" + CIDs[i]
	}

	d := NewLHComponent()
	d.CName = "water"

	vls2 := []float64{100.0, 50.0, 30.0}
	CID2s := []string{"A1", "B1", "F1"}
	PID2s := []string{"Plate2", "Plate2", "Plate2"}

	want := make([]*Liquid, 3)
	for i := 0; i < 3; i++ {
		want[i] = d.Dup()
		want[i].Loc = PID2s[i] + ":" + CID2s[i]
		want[i].Vol = vls2[i]
	}

	for i := 0; i < 2; i++ {
		if dstsDone(want) {
			t.Errorf("Done before iteration %d, should require 2 iterations", i+1)
		}
		m, err := MatchComponents(want, got, false, false)

		if err != nil {
			t.Errorf(err.Error())
		}
		updateSrcs(m, got)
		updateDsts(m, want)
	}

	if !dstsDone(want) {
		t.Errorf("Still need sources")
	}

}

func TestMatchComponentSrcSubset(t *testing.T) {
	c := NewLHComponent()
	c.CName = "water"
	c.Vol = 200.0
	CIDs := []string{"", "B1", "C1", "D1", "E1", "F1", "G1", "H1"}
	PIDs := []string{"", "Plate1", "Plate1", "Plate1", "Plate1", "Plate1", "Plate1", "Plate1"}

	got := make([]*Liquid, 8)

	for i := 0; i < 8; i++ {
		if i == 0 {
			got[i] = NewLHComponent()
		} else {
			got[i] = c.Dup()
			got[i].Loc = PIDs[i] + ":" + CIDs[i]
		}
	}

	d := NewLHComponent()
	d.CName = "water"
	d.Vol = 20.0

	CID2s := []string{"A1", "B1", "C1", "D1", "E1", "F1", "G1", "H1"}
	PID2s := []string{"Plate2", "Plate2", "Plate2", "Plate2", "Plate2", "Plate2", "Plate2", "Plate2"}

	want := make([]*Liquid, 8)
	for i := 0; i < 8; i++ {
		want[i] = d.Dup()
		want[i].Loc = PID2s[i] + ":" + CID2s[i]
	}

	for i := 0; i < 2; i++ {
		if dstsDone(want) {
			t.Errorf("Done before iteration %d, should require 2 iterations", i+1)
		}
		m, err := MatchComponents(want, got, false, false)

		if err != nil {
			t.Errorf(err.Error())
		}

		updateSrcs(m, got)
		updateDsts(m, want)
	}

	if !dstsDone(want) {
		t.Errorf("Still need sources")
	}
}

func TestMatchComponent2(t *testing.T) {
	Nams := []string{"water", "", "", "", "", "", "", ""}
	Vols := []float64{200.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0}
	CIDs := []string{"A1", "", "", "", "", "", "", ""}
	PIDs := []string{"Plate1", "", "", "", "", "", "", ""}

	got := make([]*Liquid, 8)

	for i := 0; i < 8; i++ {
		got[i] = NewLHComponent()
		got[i].CName = Nams[i]
		got[i].Vol = Vols[i]
		got[i].Loc = PIDs[i] + ":" + CIDs[i]
	}

	d := NewLHComponent()
	d.CName = "water"
	d.Vol = 20.0

	CID2s := []string{"A1", "B1", "C1", "D1", "E1", "F1", "G1", "H1"}
	PID2s := []string{"Plate2", "Plate2", "Plate2", "Plate2", "Plate2", "Plate2", "Plate2", "Plate2"}

	want := make([]*Liquid, 8)
	for i := 0; i < 8; i++ {
		want[i] = d.Dup()
		want[i].Loc = PID2s[i] + ":" + CID2s[i]
	}

	for i := 0; i < 8; i++ {
		if dstsDone(want) {
			t.Errorf("Done before iteration %d, should require 8 iterations", i+1)
		}
		m, err := MatchComponents(want, got, false, false)

		if err != nil {
			t.Errorf(err.Error())
		}
		updateSrcs(m, got)
		updateDsts(m, want)
	}

	if !dstsDone(want) {
		t.Errorf("Still need sources")
	}

}

func TestMatchComponent2b(t *testing.T) {
	Nams := []string{"water", "", "", "", "", "", "", ""}
	Vols := []float64{200.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0}
	CIDs := []string{"A1", "", "", "", "", "", "", ""}
	PIDs := []string{"Plate1", "", "", "", "", "", "", ""}

	got := make([]*Liquid, 8)

	for i := 0; i < 8; i++ {
		got[i] = NewLHComponent()
		got[i].CName = Nams[i]
		got[i].Vol = Vols[i]
		got[i].Loc = PIDs[i] + ":" + CIDs[i]
	}

	d := NewLHComponent()
	d.CName = "water"
	d.Vol = 20.0

	CID2s := []string{"A1", "B1", "C1", "D1", "E1", "F1", "G1", "H1"}
	PID2s := []string{"Plate2", "Plate2", "Plate2", "Plate2", "Plate2", "Plate2", "Plate2", "Plate2"}

	want := make([]*Liquid, 8)
	for i := 0; i < 8; i++ {
		want[i] = d.Dup()
		want[i].Loc = PID2s[i] + ":" + CID2s[i]
	}

	for i := 0; i < 8; i++ {
		if dstsDone(want) {
			t.Errorf("Done before iteration %d, should require 8 iterations", i+1)
		}
		m, err := MatchComponents(want, got, true, false)

		if err != nil {
			t.Errorf(err.Error())
		}

		updateSrcs(m, got)
		updateDsts(m, want)
	}

	if !dstsDone(want) {
		t.Errorf("Still need sources")
	}
}

func TestMatchComponent3(t *testing.T) {
	Nams := []string{"water", "", "water", "", "", "", "", ""}
	Vols := []float64{200.0, 0.0, 200.0, 0.0, 0.0, 0.0, 0.0, 0.0}
	CIDs := []string{"A1", "", "D1", "", "", "", "", ""}
	PIDs := []string{"Plate1", "", "Plate1", "", "", "", "", ""}

	got := make([]*Liquid, 8)

	for i := 0; i < 8; i++ {
		got[i] = NewLHComponent()
		got[i].CName = Nams[i]
		got[i].Vol = Vols[i]
		got[i].Loc = PIDs[i] + ":" + CIDs[i]
	}

	d := NewLHComponent()
	d.CName = "water"
	d.Vol = 20.0

	CID2s := []string{"A1", "B1", "C1", "D1", "E1", "F1", "G1", "H1"}
	PID2s := []string{"Plate2", "Plate2", "Plate2", "Plate2", "Plate2", "Plate2", "Plate2", "Plate2"}

	want := make([]*Liquid, 8)
	for i := 0; i < 8; i++ {
		want[i] = d.Dup()
		want[i].Loc = PID2s[i] + ":" + CID2s[i]
	}

	for i := 0; i < 8; i++ {
		if dstsDone(want) {
			t.Errorf("Done before iteration %d, should require 8", i+1)
		}

		m, err := MatchComponents(want, got, false, false)

		if err != nil {
			t.Errorf(err.Error())
		}

		updateSrcs(m, got)
		updateDsts(m, want)
	}

	if !dstsDone(want) {
		t.Errorf("Still need sources")
	}
}

func TestMatchComponentIndependent(t *testing.T) {
	Nams := []string{"water", "", "water", "", "", "", "", ""}
	Vols := []float64{200.0, 0.0, 200.0, 0.0, 0.0, 0.0, 0.0, 0.0}
	CIDs := []string{"A1", "", "D1", "", "", "", "", ""}
	PIDs := []string{"Plate1", "", "Plate1", "", "", "", "", ""}

	got := make([]*Liquid, 8)

	for i := 0; i < 8; i++ {
		got[i] = NewLHComponent()
		got[i].CName = Nams[i]
		got[i].Vol = Vols[i]
		got[i].Loc = PIDs[i] + ":" + CIDs[i]
	}

	d := NewLHComponent()
	d.CName = "water"
	d.Vol = 20.0

	CID2s := []string{"A1", "B1", "C1", "D1", "E1", "F1", "G1", "H1"}
	PID2s := []string{"Plate2", "Plate2", "Plate2", "Plate2", "Plate2", "Plate2", "Plate2", "Plate2"}

	want := make([]*Liquid, 8)
	for i := 0; i < 8; i++ {
		want[i] = d.Dup()
		want[i].Loc = PID2s[i] + ":" + CID2s[i]
	}

	for i := 0; i < 4; i++ {
		if dstsDone(want) {
			t.Errorf("Done before iteration %d, should require 4", i+1)
		}
		m, err := MatchComponents(want, got, true, false)

		if err != nil {
			t.Errorf(err.Error())
		}
		updateSrcs(m, got)
		updateDsts(m, want)
	}

	if !dstsDone(want) {
		t.Errorf("Still require sources: %v", want)
	}

}

func TestMatch7Subcomponents(t *testing.T) {
	c := NewLHComponent()
	c.CName = "water"
	c.Vol = 200.0
	CIDs := []string{"", "B1", "C1", "D1", "E1", "F1", "G1", "H1"}
	PIDs := []string{"", "Plate1", "Plate1", "Plate1", "Plate1", "Plate1", "Plate1", "Plate1"}

	got := make([]*Liquid, 8)

	for i := 0; i < 8; i++ {
		got[i] = c.Dup()
		got[i].Loc = PIDs[i] + ":" + CIDs[i]
		if i == 0 {
			got[i] = NewLHComponent()
		}
	}

	d := NewLHComponent()
	d.CName = "water"
	d.Vol = 20.0

	CID2s := []string{"A1", "B1", "C1", "D1", "E1", "F1", "G1", ""}
	PID2s := []string{"Plate2", "Plate2", "Plate2", "Plate2", "Plate2", "Plate2", "Plate2", ""}

	want := make([]*Liquid, 8)
	for i := 0; i < 8; i++ {
		want[i] = d.Dup()
		if i == 7 {
			want[i].Vol = 0.0
		}
		want[i].Loc = PID2s[i] + ":" + CID2s[i]
	}

	for i := 0; i < 1; i++ {
		if dstsDone(want) {
			t.Errorf("Done before iteration %d, should require 1", i+1)
		}
		m, err := MatchComponents(want, got, false, false)

		if err != nil {
			t.Errorf(err.Error())
		}
		updateSrcs(m, got)
		updateDsts(m, want)
	}

	if !dstsDone(want) {
		t.Errorf("Still need sources")
	}

}

func TestMatch7Subcomponents8wanted(t *testing.T) {
	c := NewLHComponent()
	c.CName = "tartrazine"
	c.Vol = 200.0
	CIDs := []string{"", "B1", "C1", "D1", "E1", "F1", "G1", "H1"}
	PIDs := []string{"", "Plate1", "Plate1", "Plate1", "Plate1", "Plate1", "Plate1", "Plate1"}

	got := make([]*Liquid, 8)

	for i := 0; i < 8; i++ {
		got[i] = c.Dup()
		got[i].Loc = PIDs[i] + ":" + CIDs[i]
		if i == 0 {
			got[i] = NewLHComponent()
		}
	}

	d := NewLHComponent()
	d.CName = "tartrazine"
	d.Vol = 24.0

	CID2s := []string{"A1", "B1", "C1", "D1", "E1", "F1", "G1", "H1"}
	PID2s := []string{"Plate2", "Plate2", "Plate2", "Plate2", "Plate2", "Plate2", "Plate2", "Plate2"}

	want := make([]*Liquid, 8)
	for i := 0; i < 8; i++ {
		want[i] = d.Dup()
		want[i].Loc = PID2s[i] + ":" + CID2s[i]
	}

	for i := 0; i < 2; i++ {
		if dstsDone(want) {
			t.Errorf("Done after %d iterations, should require 2", i+1)
		}
		m, err := MatchComponents(want, got, false, false)

		if err != nil {
			t.Errorf(err.Error())
		}

		updateSrcs(m, got)
		updateDsts(m, want)
	}

}

func TestNonMatchComponent(t *testing.T) {
	c := NewLHComponent()
	c.CName = "water"
	c.Vol = 200.0
	CIDs := []string{"A1", "B1", "C1", "D1", "E1", "F1", "G1", "H1"}
	PIDs := []string{"Plate1", "Plate1", "Plate1", "Plate1", "Plate1", "Plate1", "Plate1", "Plate1"}

	got := make([]*Liquid, 8)

	for i := 0; i < 8; i++ {
		got[i] = c.Dup()
		got[i].Loc = PIDs[i] + ":" + CIDs[i]
	}

	d := NewLHComponent()
	d.CName = "fishjuice"
	d.Vol = 20.0

	CID2s := []string{"A1", "B1", "C1", "D1", "E1", "F1", "G1", "H1"}
	PID2s := []string{"Plate2", "Plate2", "Plate2", "Plate2", "Plate2", "Plate2", "Plate2", "Plate2"}

	want := make([]*Liquid, 8)
	for i := 0; i < 8; i++ {
		want[i] = d.Dup()
		want[i].Loc = PID2s[i] + ":" + CID2s[i]
	}

	m, err := MatchComponents(want, got, false, false)

	if !IsNotFound(err) {
		t.Errorf(err.Error())
	}

	updateSrcs(m, got)
	updateDsts(m, want)

	if dstsDone(want) {
		t.Errorf("Negative test failed: sources were found when none were available")
	}
}
func TestMatchAllDifferentComponent(t *testing.T) {
	c := NewLHComponent()
	c.CName = "water"
	c.Vol = 200.0

	CNames := []string{"rum", "gin", "vodka", "brandy", "whisky", "sambuca", "kahluha", "grappa"}
	CIDs := []string{"A1", "B1", "C1", "D1", "E1", "F1", "G1", "H1"}
	PIDs := []string{"Plate1", "Plate1", "Plate1", "Plate1", "Plate1", "Plate1", "Plate1", "Plate1"}

	got := make([]*Liquid, 8)

	for i := 0; i < 8; i++ {
		got[i] = c.Dup()
		got[i].Loc = PIDs[i] + ":" + CIDs[i]
		got[i].CName = CNames[i]
	}

	d := NewLHComponent()
	d.CName = "water"
	d.Vol = 20.0

	CID2s := []string{"A1", "B1", "C1", "D1", "E1", "F1", "G1", "H1"}
	PID2s := []string{"Plate2", "Plate2", "Plate2", "Plate2", "Plate2", "Plate2", "Plate2", "Plate2"}

	want := make([]*Liquid, 8)
	for i := 0; i < 8; i++ {
		want[i] = d.Dup()
		want[i].Loc = PID2s[i] + ":" + CID2s[i]
		want[i].CName = CNames[i]
	}

	m, err := MatchComponents(want, got, false, false)

	if err != nil {
		t.Errorf(err.Error())
	}

	updateSrcs(m, got)
	updateDsts(m, want)

	if !dstsDone(want) {
		t.Errorf("Still want sources")
	}

}

func TestAlignIndependent(t *testing.T) {
	w := make([]*Liquid, 3)
	g := make([]*Liquid, 3)

	CIDs := []string{"A1", "B1", "C1"}
	p1 := "Plate1"
	p2 := "Plate2"

	vW := wunit.NewVolume(20.0, "ul")
	vG := wunit.NewVolume(200.0, "ul")

	cN := "water"

	// independent case

	for i := 0; i < 3; i++ {
		w[i] = NewLHComponent()
		g[i] = NewLHComponent()

		/*__*/
		g[i].Loc = p1 + ":" + CIDs[i]
		g[i].CName = cN
		g[i].Vol = vG.RawValue()

		if i != 1 {
			w[i].Loc = p2 + ":" + CIDs[i]
			w[i].CName = cN
			w[i].Vol = vW.RawValue()
		}
	}
	m := align(w, g, true, false)

	if len(m.IDs) != 3 {
		t.Errorf("Error: expected 3 IDs got %d", len(m.IDs))
	}

	expID := []string{p1, "", p1}
	expCR := []string{"A1", "", "C1"}
	expV := []wunit.Volume{vW.Dup(), wunit.ZeroVolume(), vW.Dup()}
	expM := []int{0, -1, 2}
	expSc := 40.0

	expected := Match{IDs: expID, WCs: expCR, Vols: expV, M: expM, Sc: expSc}

	if !reflect.DeepEqual(m, expected) {
		t.Errorf("Expected %v got %v", expected, m)
	}
}

func TestAlignIndependent2(t *testing.T) {
	w := make([]*Liquid, 8)
	g := make([]*Liquid, 8)

	CIDs := []string{"A1", "B1", "C1", "D1", "E1", "F1", "G1", "H1"}
	p1 := "Plate1"
	p2 := "Plate2"

	vW := wunit.NewVolume(20.0, "ul")
	vG := wunit.NewVolume(200.0, "ul")

	cN := "water"

	for i := 0; i < 8; i++ {
		w[i] = NewLHComponent()
		g[i] = NewLHComponent()

		/*__*/
		g[i].Loc = p1 + ":" + CIDs[i]
		g[i].CName = cN
		g[i].Vol = vG.RawValue()

		if i%2 != 1 {
			w[i].Loc = p2 + ":" + CIDs[i]
			w[i].CName = cN
			w[i].Vol = vW.RawValue()
		}
	}
	m := align(w, g, true, false)

	if len(m.IDs) != 8 {
		t.Errorf("Error: expected 8 IDs got %d", len(m.IDs))
	}

	expID := []string{p1, "", p1, "", p1, "", p1, ""}
	expCR := []string{"A1", "", "C1", "", "E1", "", "G1", ""}
	expV := []wunit.Volume{vW.Dup(), wunit.ZeroVolume(), vW.Dup(), wunit.ZeroVolume(), vW.Dup(), wunit.ZeroVolume(), vW.Dup(), wunit.ZeroVolume()}
	expM := []int{0, -1, 2, -1, 4, -1, 6, -1}
	expSc := 80.0

	expected := Match{IDs: expID, WCs: expCR, Vols: expV, M: expM, Sc: expSc}

	if !reflect.DeepEqual(m, expected) {
		t.Errorf("Expected %v got %v", expected, m)
	}
}

func TestAlignIndependent3(t *testing.T) {
	w := make([]*Liquid, 0, 8)
	g := make([]*Liquid, 8)

	CIDs := []string{"A1", "B1", "C1", "D1", "E1", "F1", "G1", "H1"}
	p1 := "Plate1"
	p2 := "Plate2"

	vW := wunit.NewVolume(20.0, "ul")
	vG := wunit.NewVolume(200.0, "ul")

	cN := "water"

	for i := 0; i < 8; i++ {
		wcmp := NewLHComponent()
		wcmp.Loc = p2 + ":" + CIDs[i]
		wcmp.CName = cN
		wcmp.Vol = vW.RawValue()

		g[i] = NewLHComponent()

		/*__*/
		g[i].Loc = p1 + ":" + CIDs[i]
		g[i].CName = cN
		g[i].Vol = vG.RawValue()

		if i != 6 {
			w = append(w, wcmp)
		} else {
			w = append(w, NewLHComponent())
		}
	}

	m := align(w, g, true, false)

	if len(m.IDs) != 8 {
		t.Errorf("Error: expected 8 IDs got %d", len(m.IDs))
	}

	expID := []string{p1, p1, p1, p1, p1, p1, "", p1}
	expCR := []string{"A1", "B1", "C1", "D1", "E1", "F1", "", "H1"}
	expV := []wunit.Volume{vW.Dup(), vW.Dup(), vW.Dup(), vW.Dup(), vW.Dup(), vW.Dup(), wunit.ZeroVolume(), vW.Dup()}
	expM := []int{0, 1, 2, 3, 4, 5, -1, 7}
	expSc := 140.0

	expected := Match{IDs: expID, WCs: expCR, Vols: expV, M: expM, Sc: expSc}

	if !reflect.DeepEqual(m, expected) {
		t.Errorf("Expected %v got %v", expected, m)
	}
}
