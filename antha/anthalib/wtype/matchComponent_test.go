package wtype

import (
	"fmt"
	"testing"
)

func TestmatchComponent(t *testing.T) {

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
		da[i] = c.Dup()
		da[i].Loc = PID2s[i] + ":" + CID2s[i]
	}

	cm, err := matchComponents(da, ca, false)

	if err != nil {
		t.Errorf(err.Error())
	}

	fmt.Println(cm)
}
