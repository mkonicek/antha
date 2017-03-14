package wtype

import (
	"fmt"
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
