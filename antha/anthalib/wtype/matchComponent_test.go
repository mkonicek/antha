package wtype

import (
	"fmt"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"reflect"
	"testing"
)

func TestMatchComponent(t *testing.T) {
	c := NewLHComponent()
	c.CName = "water"
	c.Vol = 200.0
	CIDs := []string{"A1", "B1", "C1", "D1", "E1", "F1", "G1", "H1"}
	PIDs := []string{"Plate1", "Plate1", "Plate1", "Plate1", "Plate1", "Plate1", "Plate1", "Plate1"}

	got := make([]*LHComponent, 8)

	for i := 0; i < 8; i++ {
		got[i] = c.Dup()
		got[i].Loc = PIDs[i] + ":" + CIDs[i]
	}

	d := NewLHComponent()
	d.CName = "water"
	d.Vol = 20.0

	CID2s := []string{"A1", "B1", "C1", "D1", "E1", "F1", "G1", "H1"}
	PID2s := []string{"Plate2", "Plate2", "Plate2", "Plate2", "Plate2", "Plate2", "Plate2", "Plate2"}

	want := make([]*LHComponent, 8)
	for i := 0; i < 8; i++ {
		want[i] = d.Dup()
		want[i].Loc = PID2s[i] + ":" + CID2s[i]
	}

	cm, err := matchComponents(want, got, false)

	if err != nil {
		t.Errorf(err.Error())
	}

	if len(cm.Matches) != 1 {
		t.Errorf(fmt.Sprintf("Exsctly one match required, got %d", len(cm.Matches)))
	}
}
func TestMatchComponentPickupVolumes(t *testing.T) {
	c := NewLHComponent()
	c.CName = "water"
	vls := []float64{100.0, 100.0, 25.0, 25.0}
	CIDs := []string{"A1", "B1", "C1", "D1"}
	PIDs := []string{"Plate1", "Plate1", "Plate1", "Plate1"}

	got := make([]*LHComponent, 4)

	for i := 0; i < 4; i++ {
		got[i] = c.Dup()
		got[i].Vol = vls[i]
		got[i].Loc = PIDs[i] + ":" + CIDs[i]
	}

	d := NewLHComponent()
	d.CName = "water"

	vls2 := []float64{110.0, 110.0, 30.0}
	CID2s := []string{"A1", "B1", "F1"}
	PID2s := []string{"Plate2", "Plate2", "Plate2"}

	want := make([]*LHComponent, 3)
	for i := 0; i < 3; i++ {
		want[i] = d.Dup()
		want[i].Loc = PID2s[i] + ":" + CID2s[i]
		want[i].Vol = vls2[i]
	}

	cm, err := matchComponents(want, got, false)

	if err != nil {
		t.Errorf(err.Error())
	}

	fmt.Println("TWO...FIVE")

	for _, m := range cm.Matches {
		fmt.Println(m)
	}

	fmt.Println("ZERO...ZERO...ZERO")
}

func TestMatchComponentSrcSubset(t *testing.T) {
	c := NewLHComponent()
	c.CName = "water"
	c.Vol = 200.0
	CIDs := []string{"", "B1", "C1", "D1", "E1", "F1", "G1", "H1"}
	PIDs := []string{"", "Plate1", "Plate1", "Plate1", "Plate1", "Plate1", "Plate1", "Plate1"}

	got := make([]*LHComponent, 8)

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

	want := make([]*LHComponent, 8)
	for i := 0; i < 8; i++ {
		want[i] = d.Dup()
		want[i].Loc = PID2s[i] + ":" + CID2s[i]
	}

	cm, err := matchComponents(want, got, false)

	if err != nil {
		t.Errorf(err.Error())
	}

	if len(cm.Matches) != 2 {
		t.Errorf(fmt.Sprintf("Exactly two matches required, got %d", len(cm.Matches)))
	}
}

func TestMatchComponent2(t *testing.T) {
	Nams := []string{"water", "", "", "", "", "", "", ""}
	Vols := []float64{200.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0}
	CIDs := []string{"A1", "", "", "", "", "", "", ""}
	PIDs := []string{"Plate1", "", "", "", "", "", "", ""}

	got := make([]*LHComponent, 8)

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

	want := make([]*LHComponent, 8)
	for i := 0; i < 8; i++ {
		want[i] = d.Dup()
		want[i].Loc = PID2s[i] + ":" + CID2s[i]
	}

	cm, err := matchComponents(want, got, false)

	if err != nil {
		t.Errorf(err.Error())
	}

	if len(cm.Matches) != 8 {
		t.Errorf(fmt.Sprintf("Exactly 8 matches required, got %d", len(cm.Matches)))
	}
}

func TestMatchComponent2b(t *testing.T) {
	Nams := []string{"water", "", "", "", "", "", "", ""}
	Vols := []float64{200.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0}
	CIDs := []string{"A1", "", "", "", "", "", "", ""}
	PIDs := []string{"Plate1", "", "", "", "", "", "", ""}

	got := make([]*LHComponent, 8)

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

	want := make([]*LHComponent, 8)
	for i := 0; i < 8; i++ {
		want[i] = d.Dup()
		want[i].Loc = PID2s[i] + ":" + CID2s[i]
	}

	cm, err := matchComponents(want, got, true)

	if err != nil {
		t.Errorf(err.Error())
	}

	if len(cm.Matches) != 8 {
		t.Errorf(fmt.Sprintf("Exactly 8 matches required, got %d", len(cm.Matches)))
	}
}

func TestMatchComponent3(t *testing.T) {
	Nams := []string{"water", "", "water", "", "", "", "", ""}
	Vols := []float64{200.0, 0.0, 200.0, 0.0, 0.0, 0.0, 0.0, 0.0}
	CIDs := []string{"A1", "", "D1", "", "", "", "", ""}
	PIDs := []string{"Plate1", "", "Plate1", "", "", "", "", ""}

	got := make([]*LHComponent, 8)

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

	want := make([]*LHComponent, 8)
	for i := 0; i < 8; i++ {
		want[i] = d.Dup()
		want[i].Loc = PID2s[i] + ":" + CID2s[i]
	}

	cm, err := matchComponents(want, got, false)

	if err != nil {
		t.Errorf(err.Error())
	}

	if len(cm.Matches) != 8 {
		t.Errorf(fmt.Sprintf("Exactly 8 matches required, got %d", len(cm.Matches)))
	}
}

func TestMatchComponentIndependent(t *testing.T) {
	Nams := []string{"water", "", "water", "", "", "", "", ""}
	Vols := []float64{200.0, 0.0, 200.0, 0.0, 0.0, 0.0, 0.0, 0.0}
	CIDs := []string{"A1", "", "D1", "", "", "", "", ""}
	PIDs := []string{"Plate1", "", "Plate1", "", "", "", "", ""}

	got := make([]*LHComponent, 8)

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

	want := make([]*LHComponent, 8)
	for i := 0; i < 8; i++ {
		want[i] = d.Dup()
		want[i].Loc = PID2s[i] + ":" + CID2s[i]
	}

	cm, err := matchComponents(want, got, true)

	if err != nil {
		t.Errorf(err.Error())
	}

	if len(cm.Matches) != 4 {
		fmt.Println(cm.Matches)
		t.Errorf(fmt.Sprintf("Exactly 4 matches required, got %d", len(cm.Matches)))
	}

}

func TestMatch7Subcomponents(t *testing.T) {
	c := NewLHComponent()
	c.CName = "water"
	c.Vol = 200.0
	CIDs := []string{"", "B1", "C1", "D1", "E1", "F1", "G1", "H1"}
	PIDs := []string{"", "Plate1", "Plate1", "Plate1", "Plate1", "Plate1", "Plate1", "Plate1"}

	got := make([]*LHComponent, 8)

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

	want := make([]*LHComponent, 8)
	for i := 0; i < 8; i++ {
		want[i] = d.Dup()
		if i == 7 {
			want[i].Vol = 0.0
		}
		want[i].Loc = PID2s[i] + ":" + CID2s[i]
	}

	cm, err := matchComponents(want, got, false)

	if err != nil {
		t.Errorf(err.Error())
	}

	if len(cm.Matches) != 1 {
		t.Errorf(fmt.Sprintf("Exsctly one match required, got %d", len(cm.Matches)))
	}
}

func TestMatch7Subcomponents8wanted(t *testing.T) {
	c := NewLHComponent()
	c.CName = "tartrazine"
	c.Vol = 200.0
	CIDs := []string{"", "B1", "C1", "D1", "E1", "F1", "G1", "H1"}
	PIDs := []string{"", "Plate1", "Plate1", "Plate1", "Plate1", "Plate1", "Plate1", "Plate1"}

	got := make([]*LHComponent, 8)

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

	want := make([]*LHComponent, 8)
	for i := 0; i < 8; i++ {
		want[i] = d.Dup()
		want[i].Loc = PID2s[i] + ":" + CID2s[i]
	}

	cm, err := matchComponents(want, got, false)

	if err != nil {
		t.Errorf(err.Error())
	}

	if len(cm.Matches) != 2 {
		t.Errorf(fmt.Sprintf("Exsctly two matches required, got %d", len(cm.Matches)))
	}
}

func TestNonMatchComponent(t *testing.T) {
	c := NewLHComponent()
	c.CName = "water"
	c.Vol = 200.0
	CIDs := []string{"A1", "B1", "C1", "D1", "E1", "F1", "G1", "H1"}
	PIDs := []string{"Plate1", "Plate1", "Plate1", "Plate1", "Plate1", "Plate1", "Plate1", "Plate1"}

	got := make([]*LHComponent, 8)

	for i := 0; i < 8; i++ {
		got[i] = c.Dup()
		got[i].Loc = PIDs[i] + ":" + CIDs[i]
	}

	d := NewLHComponent()
	d.CName = "fishjuice"
	d.Vol = 20.0

	CID2s := []string{"A1", "B1", "C1", "D1", "E1", "F1", "G1", "H1"}
	PID2s := []string{"Plate2", "Plate2", "Plate2", "Plate2", "Plate2", "Plate2", "Plate2", "Plate2"}

	want := make([]*LHComponent, 8)
	for i := 0; i < 8; i++ {
		want[i] = d.Dup()
		want[i].Loc = PID2s[i] + ":" + CID2s[i]
	}

	cm, err := matchComponents(want, got, false)

	if err != nil {
		t.Errorf(err.Error())
	}

	if len(cm.Matches) != 0 {
		t.Errorf(fmt.Sprintf("Expected 0 matches, got %d", len(cm.Matches)))
	}
}
func TestMatchAllDifferentComponent(t *testing.T) {
	c := NewLHComponent()
	c.CName = "water"
	c.Vol = 200.0

	CNames := []string{"rum", "gin", "vodka", "brandy", "whisky", "sambuca", "kahluha", "grappa"}
	CIDs := []string{"A1", "B1", "C1", "D1", "E1", "F1", "G1", "H1"}
	PIDs := []string{"Plate1", "Plate1", "Plate1", "Plate1", "Plate1", "Plate1", "Plate1", "Plate1"}

	got := make([]*LHComponent, 8)

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

	want := make([]*LHComponent, 8)
	for i := 0; i < 8; i++ {
		want[i] = d.Dup()
		want[i].Loc = PID2s[i] + ":" + CID2s[i]
		want[i].CName = CNames[i]
	}

	cm, err := matchComponents(want, got, false)

	if err != nil {
		t.Errorf(err.Error())
	}

	if len(cm.Matches) != 1 {
		t.Errorf(fmt.Sprintf("Exsctly one match required, got %d", len(cm.Matches)))
	}
}

func TestAlignIndependent(t *testing.T) {
	w := make([]*LHComponent, 3)
	g := make([]*LHComponent, 3)

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
	m := align(w, g, true)

	if len(m.IDs) != 3 {
		t.Errorf("Error: expected 3 IDs got %d", len(m.IDs))
	}

	expID := []string{p1, "", p1}
	expCR := []string{"A1", "", "C1"}
	expV := []wunit.Volume{vW.Dup(), wunit.ZeroVolume(), vW.Dup()}
	expM := []int{0, -1, 2}
	expSc := 360.0

	expected := Match{IDs: expID, WCs: expCR, Vols: expV, M: expM, Sc: expSc}

	if !reflect.DeepEqual(m, expected) {
		t.Errorf("Expected %v got %v", expected, m)
	}
}

func TestAlignIndependent2(t *testing.T) {
	w := make([]*LHComponent, 8)
	g := make([]*LHComponent, 8)

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
	m := align(w, g, true)

	if len(m.IDs) != 8 {
		t.Errorf("Error: expected 8 IDs got %d", len(m.IDs))
	}

	expID := []string{p1, "", p1, "", p1, "", p1, ""}
	expCR := []string{"A1", "", "C1", "", "E1", "", "G1", ""}
	expV := []wunit.Volume{vW.Dup(), wunit.ZeroVolume(), vW.Dup(), wunit.ZeroVolume(), vW.Dup(), wunit.ZeroVolume(), vW.Dup(), wunit.ZeroVolume()}
	expM := []int{0, -1, 2, -1, 4, -1, 6, -1}
	expSc := 720.0

	expected := Match{IDs: expID, WCs: expCR, Vols: expV, M: expM, Sc: expSc}

	if !reflect.DeepEqual(m, expected) {
		t.Errorf("Expected %v got %v", expected, m)
	}
}
