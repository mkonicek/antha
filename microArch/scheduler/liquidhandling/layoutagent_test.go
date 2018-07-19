package liquidhandling

import (
	"context"
	"encoding/json"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wutil"
	"github.com/antha-lang/antha/inventory"
	"github.com/antha-lang/antha/inventory/testinventory"
	"testing"
)

func TestLayoutAgent1(t *testing.T) {
	ctx := testinventory.NewContext(context.Background())

	req := GetLHRequestForTest()
	configure_request_simple(ctx, req)
	pt, _ := inventory.NewPlate(ctx, "pcrplate_skirted")
	req.Output_platetypes = append(req.Output_platetypes, pt)

	// marshal / unmarshal

	b, err := json.Marshal(req)

	if err != nil {
		t.Errorf(err.Error())
	}

	err = json.Unmarshal(b, &req)
	if err != nil {
		t.Errorf(err.Error())
	}

	req.ConfigureYourself()
	setOutputOrder(req)
	params := makeGilson(ctx)

	req, err = ImprovedLayoutAgent(ctx, req, params)

	if err != nil {
		t.Fatal(err)
	}

	testReq(t, req)
}

func TestLayoutAgent2(t *testing.T) {
	ctx := testinventory.NewContext(context.Background())

	req := GetLHRequestForTest()
	configure_request_simple(ctx, req)

	// add a destination plate (i.e. MixInto)

	pt, _ := inventory.NewPlate(ctx, "pcrplate_skirted")
	req.Output_platetypes = append(req.Output_platetypes, pt)
	req.Output_plates[pt.ID] = pt

	for _, ins := range req.LHInstructions {
		ins.OutPlate = pt
		ins.SetPlateID(pt.ID)
	}

	// marshal / unmarshal

	b, err := json.Marshal(req)

	if err != nil {
		t.Errorf(err.Error())
	}

	err = json.Unmarshal(b, &req)
	if err != nil {
		t.Errorf(err.Error())
	}

	req.ConfigureYourself()
	setOutputOrder(req)
	params := makeGilson(ctx)

	for _, ins := range req.LHInstructions {
		if ins.PlateID != ins.OutPlate.ID {
			t.Errorf("EARLY MixInto ID mismatch: expected %s got %s", ins.PlateID, ins.OutPlate.ID)
		}
	}

	req, err = ImprovedLayoutAgent(ctx, req, params)

	if err != nil {
		t.Fatal(err)
	}

	testReq(t, req)

	// we additionally need to verify that the ID set in the instructions corresponds to that in the plate

	for _, ins := range req.LHInstructions {
		if ins.PlateID != ins.OutPlate.ID {
			t.Errorf("MixInto ID mismatch: expected %s got %s", ins.PlateID, ins.OutPlate.ID)
		}
	}
}

// a mix of specific dests and no dest
func TestLayoutAgent3(t *testing.T) {
	ctx := testinventory.NewContext(context.Background())

	req := GetLHRequestForTest()
	configure_request_simple(ctx, req)

	// add a destination plate (i.e. MixInto)

	pt, _ := inventory.NewPlate(ctx, "pcrplate_skirted")
	req.Output_platetypes = append(req.Output_platetypes, pt)
	req.Output_plates[pt.ID] = pt

	i := -1
	for _, ins := range req.LHInstructions {
		i += 1
		if i%2 == 1 {
			continue
		}
		ins.OutPlate = pt
		ins.SetPlateID(pt.ID)
	}

	// marshal / unmarshal

	b, err := json.Marshal(req)

	if err != nil {
		t.Errorf(err.Error())
	}

	err = json.Unmarshal(b, &req)
	if err != nil {
		t.Errorf(err.Error())
	}

	req.ConfigureYourself()
	setOutputOrder(req)
	params := makeGilson(ctx)

	for _, ins := range req.LHInstructions {
		if ins.PlateID != "" {
			if ins.OutPlate == nil || ins.PlateID != ins.OutPlate.ID {
				t.Errorf("EARLY MixInto ID mismatch: expected %s got %s", ins.PlateID, ins.OutPlate.ID)
			}
		}
	}

	req, err = ImprovedLayoutAgent(ctx, req, params)

	if err != nil {
		t.Fatal(err)
	}

	testReq(t, req)

	// we additionally need to verify that the ID set in the instructions corresponds to that in the plate

	for _, ins := range req.LHInstructions {
		if ins.PlateID != "" {
			if ins.OutPlate != nil && ins.PlateID != ins.OutPlate.ID {
				t.Errorf("MixInto ID mismatch: expected %s got %s", ins.PlateID, ins.OutPlate.ID)
			}
		}
	}
}

func TestLayoutAgent4(t *testing.T) {
	ctx := testinventory.NewContext(context.Background())

	req := GetLHRequestForTest()
	configure_request_simple(ctx, req)

	// add a destination plate (i.e. MixInto)

	pt, _ := inventory.NewPlate(ctx, "pcrplate_skirted")
	req.Output_platetypes = append(req.Output_platetypes, pt)
	req.Output_plates[pt.ID] = pt

	for _, ins := range req.LHInstructions {
		ins.PlateName = "Funk Plate"
	}

	// marshal / unmarshal

	b, err := json.Marshal(req)

	if err != nil {
		t.Errorf(err.Error())
	}

	err = json.Unmarshal(b, &req)
	if err != nil {
		t.Errorf(err.Error())
	}

	req.ConfigureYourself()
	setOutputOrder(req)
	params := makeGilson(ctx)

	req, err = ImprovedLayoutAgent(ctx, req, params)

	if err != nil {
		t.Fatal(err)
	}

	testReq(t, req)

	// test the plate name is OK
	for _, ins := range req.LHInstructions {
		if ins.PlateName != "Funk Plate" {
			t.Errorf("Plate name issue - expected %s got %s", "Funk plate", ins.PlateName)
		}
	}
}

// bigger tests
func TestLayoutAgent5(t *testing.T) {
	ctx := testinventory.NewContext(context.Background())

	req := GetLHRequestForTest()
	configure_request_bigger(ctx, req)

	// add a destination plate (i.e. MixInto)

	pt, _ := inventory.NewPlate(ctx, "pcrplate_skirted")
	req.Output_platetypes = append(req.Output_platetypes, pt)
	req.Output_plates[pt.ID] = pt

	for _, ins := range req.LHInstructions {
		ins.PlateName = "Funk Plate"
	}

	// marshal / unmarshal

	b, err := json.Marshal(req)

	if err != nil {
		t.Errorf(err.Error())
	}

	err = json.Unmarshal(b, &req)
	if err != nil {
		t.Errorf(err.Error())
	}

	req.ConfigureYourself()
	setOutputOrder(req)
	params := makeGilson(ctx)

	req, err = ImprovedLayoutAgent(ctx, req, params)

	if err != nil {
		t.Fatal(err)
	}

	testReqBig(t, req)
}

func TestLayoutAgent6(t *testing.T) {
	ctx := testinventory.NewContext(context.Background())

	req := GetLHRequestForTest()
	configure_request_bigger(ctx, req)

	// add a destination plate (i.e. MixInto)

	pt, _ := inventory.NewPlate(ctx, "pcrplate_skirted")
	req.Output_platetypes = append(req.Output_platetypes, pt)
	req.Output_plates[pt.ID] = pt

	// bogus mixInto with too many samples for the wells
	for _, ins := range req.LHInstructions {
		ins.OutPlate = pt
		ins.PlateID = pt.ID
	}

	// marshal / unmarshal

	b, err := json.Marshal(req)

	if err != nil {
		t.Errorf(err.Error())
	}

	err = json.Unmarshal(b, &req)
	if err != nil {
		t.Errorf(err.Error())
	}

	req.ConfigureYourself()
	setOutputOrder(req)
	params := makeGilson(ctx)

	req, err = ImprovedLayoutAgent(ctx, req, params)

	if err != nil {
		t.Fatal(err)
	}

	testReqBig(t, req)
}

func TestLayoutAgent7(t *testing.T) {
	ctx := testinventory.NewContext(context.Background())

	req := GetLHRequestForTest()
	configure_request_simple(ctx, req)

	// add a destination plate (i.e. MixInto)

	pt, _ := inventory.NewPlate(ctx, "pcrplate_skirted")
	req.Output_platetypes = append(req.Output_platetypes, pt)
	req.Output_plates[pt.ID] = pt

	for _, ins := range req.LHInstructions {
		ins.PlateName = "Funk plate"
		ins.Welladdress = "A1"
	}

	// marshal / unmarshal

	b, err := json.Marshal(req)

	if err != nil {
		t.Errorf(err.Error())
	}

	err = json.Unmarshal(b, &req)
	if err != nil {
		t.Errorf(err.Error())
	}

	req.ConfigureYourself()
	setOutputOrder(req)
	params := makeGilson(ctx)

	req, err = ImprovedLayoutAgent(ctx, req, params)

	if err != nil {
		t.Fatal(err)
	}

	plateIDs := make([]string, 0, 1)
	wells := make([]string, 0, 9)
	for _, ins := range req.LHInstructions {
		// we expect them all to end up on one plate, in 9 distinct wells going down a column
		plateIDs = append(plateIDs, ins.PlateID)
		wells = append(wells, ins.Welladdress)
	}

	if len(wutil.SADistinct(plateIDs)) != 9 {
		t.Errorf("Expected 9 plate IDs but got %d", len(wutil.SADistinct(plateIDs)))
	}

	if len(wutil.SADistinct(wells)) != 1 {
		t.Errorf("Expected 1 distinct well but got %d", len(wutil.SADistinct(wells)))
	}
}

func testReq(t *testing.T, req *LHRequest) {
	plateIDs := make([]string, 0, 1)
	wells := make([]string, 0, 9)
	for _, ins := range req.LHInstructions {
		// we expect them all to end up on one plate, in 9 distinct wells going down a column
		plateIDs = append(plateIDs, ins.PlateID)
		wells = append(wells, ins.Welladdress)
	}

	if len(wutil.SADistinct(plateIDs)) != 1 {
		t.Errorf("Expected 1 plate ID but got %d", len(wutil.SADistinct(plateIDs)))
	}

	if len(wutil.SADistinct(wells)) != 9 {
		t.Errorf("Expected 9 distinct wells but got %d", len(wutil.SADistinct(wells)))
	}

	expectedWells := make(map[string]bool, 9)

	for x := 0; x < 9; x++ {
		wc := wtype.WellCoords{X: x / 8, Y: x % 8}
		expectedWells[wc.FormatA1()] = true
	}

	for _, w := range wells {
		if !expectedWells[w] {
			t.Errorf("Error: Unexpected well %s", w)
		}

		delete(expectedWells, w)
	}

	if len(expectedWells) != 0 {
		t.Errorf("Error: Expecting more wells: %v", expectedWells)
	}
}

func testReqBig(t *testing.T, req *LHRequest) {
	plateIDs := make([]string, 0, 2)
	wells := make([]string, 0, 96)
	for _, ins := range req.LHInstructions {
		// we expect them all to end up on one plate, in 9 distinct wells going down a column
		plateIDs = append(plateIDs, ins.PlateID)
		wells = append(wells, ins.Welladdress)
	}

	if len(wutil.SADistinct(plateIDs)) != 2 {
		t.Errorf("Expected 2 plate IDs but got %d", len(wutil.SADistinct(plateIDs)))
	}

	if len(wells) != 99 {
		t.Errorf("Expected 99 wells got %d", len(wells))
	}

	if len(wutil.SADistinct(wells)) != 96 {
		t.Errorf("Expected 96 distinct wells but got %d", len(wutil.SADistinct(wells)))
	}

	expectedWells := make(map[string]int, 96)

	for x := 0; x < 96; x++ {
		wc := wtype.WellCoords{X: x / 8, Y: x % 8}
		if x < 3 {
			expectedWells[wc.FormatA1()] = 2
		} else {
			expectedWells[wc.FormatA1()] = 1
		}
	}

	for _, w := range wells {
		if expectedWells[w] == 0 {
			t.Errorf("Error: Unexpected well %s", w)
		}

		expectedWells[w] -= 1
	}

	for k, v := range expectedWells {
		if v != 0 {
			t.Errorf("Error: Expecting more wells: %s %d", k, v)
		}
	}
}
