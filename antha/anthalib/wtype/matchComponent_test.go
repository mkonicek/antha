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

	ca := make([]*LHComponent, 8)

	for i := 0; i < 8; i++ {
		ca[i] = c.Dup()
		ca[i].Loc = PIDs[i] + ":" + CIDs[i]
	}

	d := NewLHComponent()
	d.CName = "water"
	d.Vol = 20.0

	CID2s := []string{"A1", "B1", "C1", "D1", "E1", "F1", "G1", "H1"}
	PID2s := []string{"Plate2", "Plate2", "Plate2", "Plate2", "Plate2", "Plate2", "Plate2", "Plate2"}

	da := make([]*LHComponent, 8)
	for i := 0; i < 8; i++ {
		da[i] = d.Dup()
		da[i].Loc = PID2s[i] + ":" + CID2s[i]
	}

	cm, err := matchComponents(da, ca, false)

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

	ca := make([]*LHComponent, 8)

	for i := 0; i < 8; i++ {
		ca[i] = NewLHComponent()
		ca[i].CName = Nams[i]
		ca[i].Vol = Vols[i]
		ca[i].Loc = PIDs[i] + ":" + CIDs[i]
	}

	d := NewLHComponent()
	d.CName = "water"
	d.Vol = 20.0

	CID2s := []string{"A1", "B1", "C1", "D1", "E1", "F1", "G1", "H1"}
	PID2s := []string{"Plate2", "Plate2", "Plate2", "Plate2", "Plate2", "Plate2", "Plate2", "Plate2"}

	da := make([]*LHComponent, 8)
	for i := 0; i < 8; i++ {
		da[i] = d.Dup()
		da[i].Loc = PID2s[i] + ":" + CID2s[i]
	}

	cm, err := matchComponents(da, ca, false)

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
	CIDs := []string{"A1", "", "", "D1", "", "", "", ""}
	PIDs := []string{"Plate1", "", "Plate1", "", "", "", "", ""}

	ca := make([]*LHComponent, 8)

	for i := 0; i < 8; i++ {
		ca[i] = NewLHComponent()
		ca[i].CName = Nams[i]
		ca[i].Vol = Vols[i]
		ca[i].Loc = PIDs[i] + ":" + CIDs[i]
	}

	d := NewLHComponent()
	d.CName = "water"
	d.Vol = 20.0

	CID2s := []string{"A1", "B1", "C1", "D1", "E1", "F1", "G1", "H1"}
	PID2s := []string{"Plate2", "Plate2", "Plate2", "Plate2", "Plate2", "Plate2", "Plate2", "Plate2"}

	da := make([]*LHComponent, 8)
	for i := 0; i < 8; i++ {
		da[i] = d.Dup()
		da[i].Loc = PID2s[i] + ":" + CID2s[i]
	}

	cm, err := matchComponents(da, ca, false)

	if err != nil {
		t.Errorf(err.Error())
	}

	if len(cm.Matches) != 8 {
		t.Errorf(fmt.Sprintf("Exactly 8 matches required, got %d", len(cm.Matches)))
	}

	for _, v := range cm.Matches {
		fmt.Println(v.Vols)
	}

}
