package liquidhandling

import (
	"fmt"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/inventory"
	"reflect"
	"testing"
)

func getComponentOrder(inss []TerminalRobotInstruction) []string {
	ret := make([]string, 0, len(inss))

	for _, ins := range inss {

		if ins.Type() != DSP {
			continue
		}

		cmpint := ins.GetParameter("COMPONENT")

		switch param := cmpint.(type) {
		case string:
			ret = append(ret, param)
		case []string:
			ret = append(ret, param...)
		case [][]string:
			for _, cmpa := range param {
				ret = append(ret, cmpa...)
			}
		}
	}

	return ret
}

// Regression test for a bug (Antha-2357) which found that the 'first multichannel,
// then single' generation order occasionally overrode the user's required order
// uninentionally. This would arise where two sets of transfers were requested but
// the first set could not multichannel while the second could. In that case the
// first set would be moved second because all single channel instructions would
// only occur after all multichannel instructions
func TestOrdering(t *testing.T) {
	ctx := GetContextForTest()
	dstp, _ := inventory.NewPlate(ctx, "DWST12")
	rbt := getTestRobot(ctx, dstp, "pcrplate_skirted_riser40")
	pol, err := wtype.GetLHPolicyForTest()

	if err != nil {
		t.Errorf(err.Error())
	}

	ins := getTransferMulti()[0]

	iTree := NewITree(ins)
	if _, err := iTree.Build(ctx, pol, rbt); err != nil {
		t.Fatal(err)
	}

	inss := iTree.Leaves()

	if len(inss) != 58 {
		t.Errorf("Expected 58 instructions, got %d", len(inss))
	}

	expectedOrder := []string{"water", "water", "water", "wine", "wine", "wine"}

	order := getComponentOrder(inss)

	if !reflect.DeepEqual(order, expectedOrder) {
		t.Errorf(fmt.Sprintf("Expected order %v got %v", expectedOrder, order))
	}
}

func TestOrdering2(t *testing.T) {
	ctx := GetContextForTest()
	dstp, _ := inventory.NewPlate(ctx, "DSW96")
	rbt := getTestRobot(ctx, dstp, "pcrplate_skirted_riser40")
	pol, err := wtype.GetLHPolicyForTest()

	if err != nil {
		t.Errorf(err.Error())
	}

	ins := getTransferMulti()[1]

	iTree := NewITree(ins)

	if _, err := iTree.Build(ctx, pol, rbt); err != nil {
		t.Fatal(err)
	}

	inss := iTree.Leaves()

	if len(inss) != 58 {
		t.Errorf("Expected 58 instructions, got %d", len(inss))
	}
	expectedOrder := []string{"water", "fish", "evil", "wine", "slate", "fish"}

	order := getComponentOrder(inss)

	if !reflect.DeepEqual(order, expectedOrder) {
		t.Errorf(fmt.Sprintf("Expected order %v got %v", expectedOrder, order))
	}
}

func getTransferMulti() []RobotInstruction {
	vol := wunit.NewVolume(100.0, "ul")
	v2 := wunit.NewVolume(5000.0, "ul")
	v3 := wunit.NewVolume(0.0, "ul")

	tfrs := []RobotInstruction{}

	tfr := NewTransferInstruction(
		[]string{"water", "water", "water"},
		[]string{"position_4", "position_4", "position_4"},
		[]string{"position_8", "position_8", "position_8"},
		[]string{"A1", "A2", "A3"},
		[]string{"A1", "B1", "C1"},
		[]string{"DWST12", "DWST12", "DWST12"},
		[]string{"DSW96", "DSW96", "DSW96"},
		[]wunit.Volume{vol.Dup(), vol.Dup(), vol.Dup()},
		[]wunit.Volume{v2.Dup(), v2.Dup(), v2.Dup()},
		[]wunit.Volume{v3.Dup(), v3.Dup(), v3.Dup()},
		[]int{8, 8, 8},
		[]int{12, 12, 12},
		[]int{8, 8, 8},
		[]int{12, 12, 12},
		[]string{"water", "water", "water"},
		[]wtype.LHPolicy{nil, nil, nil},
	)

	tfr.Add(MTPFromArrays(
		[]string{"solvent", "solvent", "solvent"},
		[]string{"position_4", "position_4", "position_4"},
		[]string{"position_8", "position_8", "position_8"},
		[]string{"A4", "A4", "A4"},
		[]string{"A1", "B1", "C1"},
		[]string{"DWST12", "DWST12", "DWST12"},
		[]string{"DSW96", "DSW96", "DSW96"},
		[]wunit.Volume{vol.Dup(), vol.Dup(), vol.Dup()},
		[]wunit.Volume{v2.Dup(), v2.Dup(), v2.Dup()},
		[]wunit.Volume{v3.Dup(), v3.Dup(), v3.Dup()},
		[]int{8, 8, 8},
		[]int{12, 12, 12},
		[]int{8, 8, 8},
		[]int{12, 12, 12},
		[]string{"wine", "wine", "wine"},
		[]wtype.LHPolicy{nil, nil, nil},
	))

	tfrs = append(tfrs, tfr)

	tfr = NewTransferInstruction(
		[]string{"water", "water", "water"},
		[]string{"position_4", "position_4", "position_4"},
		[]string{"position_8", "position_8", "position_8"},
		[]string{"A1", "A2", "A3"},
		[]string{"A1", "B1", "C1"},
		[]string{"DSW96", "DSW96", "DSW96"},
		[]string{"DSW96", "DSW96", "DSW96"},
		[]wunit.Volume{vol.Dup(), vol.Dup(), vol.Dup()},
		[]wunit.Volume{v2.Dup(), v2.Dup(), v2.Dup()},
		[]wunit.Volume{v3.Dup(), v3.Dup(), v3.Dup()},
		[]int{8, 8, 8},
		[]int{12, 12, 12},
		[]int{8, 8, 8},
		[]int{12, 12, 12},
		[]string{"water", "fish", "evil"},
		[]wtype.LHPolicy{nil, nil, nil},
	)

	tfr.Add(MTPFromArrays(
		[]string{"solvent", "solvent", "solvent"},
		[]string{"position_4", "position_4", "position_4"},
		[]string{"position_8", "position_8", "position_8"},
		[]string{"A4", "B4", "C4"},
		[]string{"A1", "B1", "C1"},
		[]string{"DSW96", "DSW96", "DSW96"},
		[]string{"DSW96", "DSW96", "DSW96"},
		[]wunit.Volume{vol.Dup(), vol.Dup(), vol.Dup()},
		[]wunit.Volume{v2.Dup(), v2.Dup(), v2.Dup()},
		[]wunit.Volume{v3.Dup(), v3.Dup(), v3.Dup()},
		[]int{8, 8, 8},
		[]int{12, 12, 12},
		[]int{8, 8, 8},
		[]int{12, 12, 12},
		[]string{"wine", "slate", "fish"},
		[]wtype.LHPolicy{nil, nil, nil},
	))

	tfrs = append(tfrs, tfr)

	tfr = NewTransferInstruction(
		[]string{"water", "water", "water"},
		[]string{"position_4", "position_4", "position_4"},
		[]string{"position_8", "position_8", "position_8"},
		[]string{"A1", "A2", "A3"},
		[]string{"A1", "B1", "C1"},
		[]string{"DSW96", "DSW96", "DSW96"},
		[]string{"DSW96", "DSW96", "DSW96"},
		[]wunit.Volume{vol.Dup(), vol.Dup(), vol.Dup()},
		[]wunit.Volume{v2.Dup(), v2.Dup(), v2.Dup()},
		[]wunit.Volume{v3.Dup(), v3.Dup(), v3.Dup()},
		[]int{8, 8, 8},
		[]int{12, 12, 12},
		[]int{8, 8, 8},
		[]int{12, 12, 12},
		[]string{"water", "fish", "evil"},
		[]wtype.LHPolicy{nil, nil, nil},
	)

	tfr.Add(MTPFromArrays(
		[]string{"solvent", "solvent", "solvent"},
		[]string{"position_4", "position_4", "position_4"},
		[]string{"position_8", "position_8", "position_8"},
		[]string{"A4", "B4", "C4"},
		[]string{"A1", "B1", "C1"},
		[]string{"DSW96", "DSW96", "DSW96"},
		[]string{"DSW96", "DSW96", "DSW96"},
		[]wunit.Volume{vol.Dup(), vol.Dup(), vol.Dup()},
		[]wunit.Volume{v2.Dup(), v2.Dup(), v2.Dup()},
		[]wunit.Volume{v3.Dup(), v3.Dup(), v3.Dup()},
		[]int{8, 8, 8},
		[]int{12, 12, 12},
		[]int{8, 8, 8},
		[]int{12, 12, 12},
		[]string{"fish", "slate", "wine"},
		[]wtype.LHPolicy{nil, nil, nil},
	))

	tfrs = append(tfrs, tfr)

	return tfrs
}

func TestOrdering3(t *testing.T) {
	ctx := GetContextForTest()
	dstp, _ := inventory.NewPlate(ctx, "DSW96")
	rbt := getTestRobot(ctx, dstp, "pcrplate_skirted_riser40")
	pol, err := wtype.GetLHPolicyForTest()

	if err != nil {
		t.Errorf(err.Error())
	}

	ins := getTransferMulti()[2]

	iTree := NewITree(ins)

	if _, err := iTree.Build(ctx, pol, rbt); err != nil {
		t.Fatal(err)
	}

	inss := iTree.Leaves()

	if len(inss) != 58 {
		t.Errorf("Expected 58 instructions, got %d", len(inss))
	}

	expectedOrder := []string{"water", "fish", "evil", "fish", "slate", "wine"}

	order := getComponentOrder(inss)

	if !reflect.DeepEqual(order, expectedOrder) {
		t.Errorf(fmt.Sprintf("Expected order %v got %v", expectedOrder, order))
	}
}
