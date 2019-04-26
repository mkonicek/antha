package tests

import (
	"fmt"
	"testing"

	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wutil"
	"github.com/antha-lang/antha/laboratory"
	"github.com/antha-lang/antha/laboratory/testlab"
	lh "github.com/antha-lang/antha/microArch/scheduler/liquidhandling"
)

type layoutAgentTest struct {
	Request     *lh.LHRequest
	TestFn      func(*laboratory.Laboratory, *lh.LHRequest) error
	TestPlateID bool
}

func (self layoutAgentTest) Run(lab *laboratory.Laboratory) error {
	if ichain, err := lh.BuildInstructionChain(lab.IDGenerator, self.Request.LHInstructions); err != nil {
		return err
	} else {
		ichain.SortInstructions(self.Request.Options.OutputSort)
		self.Request.UpdateWithNewLHInstructions(ichain.GetOrderedLHInstructions())
		self.Request.InstructionChain = ichain
		self.Request.OutputOrder = ichain.FlattenInstructionIDs()
	}
	params := makeGilson(lab)

	if self.TestPlateID {
		for _, ins := range self.Request.LHInstructions {
			if ins.PlateID != "" && (ins.OutPlate == nil || ins.PlateID != ins.OutPlate.ID) {
				return fmt.Errorf("EARLY MixInto ID mismatch: expected %s got %s", ins.PlateID, ins.OutPlate.ID)
			}
		}
	}

	if err := lh.ImprovedLayoutAgent(lab.LaboratoryEffects, self.Request, params); err != nil {
		return err
	}

	if err := self.TestFn(lab, self.Request); err != nil {
		return err
	}

	if self.TestPlateID {
		for _, ins := range self.Request.LHInstructions {
			if ins.OutPlate != nil && ins.PlateID != ins.OutPlate.ID {
				return fmt.Errorf("MixInto ID mismatch: expected %s got %s", ins.PlateID, ins.OutPlate.ID)
			}
		}
	}
	return nil
}

func TestLayoutAgent1(t *testing.T) {
	testlab.WithTestLab(t, "", &testlab.TestElementCallbacks{
		Steps: func(lab *laboratory.Laboratory) error {
			req := GetLHRequestForTest(lab.IDGenerator)
			configure_request_simple(lab, req)
			pt, _ := lab.Inventory.Plates.NewPlate("pcrplate_skirted")
			req.OutputPlatetypes = append(req.OutputPlatetypes, pt)

			return layoutAgentTest{
				Request:     req,
				TestFn:      testReq,
				TestPlateID: true,
			}.Run(lab)
		},
	})
}

func TestLayoutAgent2(t *testing.T) {
	testlab.WithTestLab(t, "", &testlab.TestElementCallbacks{
		Steps: func(lab *laboratory.Laboratory) error {
			req := GetLHRequestForTest(lab.IDGenerator)
			configure_request_simple(lab, req)
			pt, _ := lab.Inventory.Plates.NewPlate("pcrplate_skirted")
			req.OutputPlatetypes = append(req.OutputPlatetypes, pt)
			req.OutputPlates[pt.ID] = pt

			for _, ins := range req.LHInstructions {
				ins.OutPlate = pt
				ins.SetPlateID(pt.ID)
			}

			return layoutAgentTest{
				Request:     req,
				TestFn:      testReq,
				TestPlateID: true,
			}.Run(lab)
		},
	})
}

// a mix of specific dests and no dest
func TestLayoutAgent3(t *testing.T) {
	testlab.WithTestLab(t, "", &testlab.TestElementCallbacks{
		Steps: func(lab *laboratory.Laboratory) error {
			req := GetLHRequestForTest(lab.IDGenerator)
			configure_request_simple(lab, req)

			// add a destination plate (i.e. MixInto)

			pt, _ := lab.Inventory.Plates.NewPlate("pcrplate_skirted")
			req.OutputPlatetypes = append(req.OutputPlatetypes, pt)
			req.OutputPlates[pt.ID] = pt

			i := -1
			for _, ins := range req.LHInstructions {
				i += 1
				if i%2 == 1 {
					continue
				}
				ins.OutPlate = pt
				ins.SetPlateID(pt.ID)
			}

			return layoutAgentTest{
				Request:     req,
				TestFn:      testReq,
				TestPlateID: true,
			}.Run(lab)
		},
	})
}

func TestLayoutAgent4(t *testing.T) {
	testlab.WithTestLab(t, "", &testlab.TestElementCallbacks{
		Steps: func(lab *laboratory.Laboratory) error {
			req := GetLHRequestForTest(lab.IDGenerator)
			configure_request_simple(lab, req)

			// add a destination plate (i.e. MixInto)

			pt, _ := lab.Inventory.Plates.NewPlate("pcrplate_skirted")
			req.OutputPlatetypes = append(req.OutputPlatetypes, pt)
			req.OutputPlates[pt.ID] = pt

			for _, ins := range req.LHInstructions {
				ins.PlateName = "Funk Plate"
			}

			return layoutAgentTest{
				Request: req,
				TestFn: func(lab *laboratory.Laboratory, req *lh.LHRequest) error {
					if err := testReq(lab, req); err != nil {
						return err
					}
					// test the plate name is OK
					for _, ins := range req.LHInstructions {
						if ins.PlateName != "Funk Plate" {
							return fmt.Errorf("Plate name issue - expected %s got %s", "Funk plate", ins.PlateName)
						}
					}
					return nil
				},
				TestPlateID: true,
			}.Run(lab)
		},
	})
}

// bigger tests
func TestLayoutAgent5(t *testing.T) {
	testlab.WithTestLab(t, "", &testlab.TestElementCallbacks{
		Steps: func(lab *laboratory.Laboratory) error {
			req := GetLHRequestForTest(lab.IDGenerator)
			configure_request_bigger(lab, req)

			// add a destination plate (i.e. MixInto)

			pt, _ := lab.Inventory.Plates.NewPlate("pcrplate_skirted")
			req.OutputPlatetypes = append(req.OutputPlatetypes, pt)
			req.OutputPlates[pt.ID] = pt

			for _, ins := range req.LHInstructions {
				ins.PlateName = "Funk Plate"
			}

			return layoutAgentTest{
				Request:     req,
				TestFn:      testReqBig,
				TestPlateID: true,
			}.Run(lab)
		},
	})
}

func TestLayoutAgent6(t *testing.T) {
	testlab.WithTestLab(t, "", &testlab.TestElementCallbacks{
		Steps: func(lab *laboratory.Laboratory) error {
			req := GetLHRequestForTest(lab.IDGenerator)
			configure_request_bigger(lab, req)

			// add a destination plate (i.e. MixInto)

			pt, _ := lab.Inventory.Plates.NewPlate("pcrplate_skirted")
			req.OutputPlatetypes = append(req.OutputPlatetypes, pt)
			req.OutputPlates[pt.ID] = pt

			// bogus mixInto with too many samples for the wells
			for _, ins := range req.LHInstructions {
				ins.OutPlate = pt
				ins.PlateID = pt.ID
			}

			return layoutAgentTest{
				Request: req,
				TestFn:  testReqBig,
			}.Run(lab)
		},
	})
}

func TestLayoutAgent7(t *testing.T) {
	testlab.WithTestLab(t, "", &testlab.TestElementCallbacks{
		Steps: func(lab *laboratory.Laboratory) error {
			req := GetLHRequestForTest(lab.IDGenerator)
			configure_request_simple(lab, req)

			// add a destination plate (i.e. MixInto)

			pt, _ := lab.Inventory.Plates.NewPlate("pcrplate_skirted")
			req.OutputPlatetypes = append(req.OutputPlatetypes, pt)
			req.OutputPlates[pt.ID] = pt

			for _, ins := range req.LHInstructions {
				ins.PlateName = "Funk plate"
				ins.Welladdress = "A1"
			}

			return layoutAgentTest{
				Request: req,
				TestFn: func(lab *laboratory.Laboratory, req *lh.LHRequest) error {
					plateIDs := make([]string, 0, 1)
					wells := make([]string, 0, 9)
					for _, ins := range req.LHInstructions {
						// we expect them all to end up on one plate, in 9 distinct wells going down a column
						plateIDs = append(plateIDs, ins.PlateID)
						wells = append(wells, ins.Welladdress)
					}

					if len(wutil.SADistinct(plateIDs)) != 9 {
						return fmt.Errorf("Expected 9 plate IDs but got %d", len(wutil.SADistinct(plateIDs)))
					}

					if len(wutil.SADistinct(wells)) != 1 {
						return fmt.Errorf("Expected 1 distinct well but got %d", len(wutil.SADistinct(wells)))
					}
					return nil
				},
				TestPlateID: true,
			}.Run(lab)
		},
	})
}

func testReq(lab *laboratory.Laboratory, req *lh.LHRequest) error {
	plateIDs := make([]string, 0, 1)
	wells := make([]string, 0, 9)
	for _, ins := range req.LHInstructions {
		// we expect them all to end up on one plate, in 9 distinct wells going down a column
		plateIDs = append(plateIDs, ins.PlateID)
		wells = append(wells, ins.Welladdress)
	}

	if len(wutil.SADistinct(plateIDs)) != 1 {
		return fmt.Errorf("Expected 1 plate ID but got %d", len(wutil.SADistinct(plateIDs)))
	}

	if len(wutil.SADistinct(wells)) != 9 {
		return fmt.Errorf("Expected 9 distinct wells but got %d", len(wutil.SADistinct(wells)))
	}

	expectedWells := make(map[string]bool, 9)

	for x := 0; x < 9; x++ {
		wc := wtype.WellCoords{X: x / 8, Y: x % 8}
		expectedWells[wc.FormatA1()] = true
	}

	for _, w := range wells {
		if !expectedWells[w] {
			return fmt.Errorf("Error: Unexpected well %s", w)
		}

		delete(expectedWells, w)
	}

	if len(expectedWells) != 0 {
		return fmt.Errorf("Error: Expecting more wells: %v", expectedWells)
	}
	return nil
}

func testReqBig(lab *laboratory.Laboratory, req *lh.LHRequest) error {
	plateIDs := make([]string, 0, 2)
	wells := make([]string, 0, 96)
	for _, ins := range req.LHInstructions {
		// we expect them all to end up on one plate, in 9 distinct wells going down a column
		plateIDs = append(plateIDs, ins.PlateID)
		wells = append(wells, ins.Welladdress)
	}

	if len(wutil.SADistinct(plateIDs)) != 2 {
		return fmt.Errorf("Expected 2 plate IDs but got %d", len(wutil.SADistinct(plateIDs)))
	}

	if len(wells) != 99 {
		return fmt.Errorf("Expected 99 wells got %d", len(wells))
	}

	if len(wutil.SADistinct(wells)) != 96 {
		return fmt.Errorf("Expected 96 distinct wells but got %d", len(wutil.SADistinct(wells)))
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
			return fmt.Errorf("Error: Unexpected well %s", w)
		}

		expectedWells[w] -= 1
	}

	for k, v := range expectedWells {
		if v != 0 {
			return fmt.Errorf("Error: Expecting more wells: %s %d", k, v)
		}
	}
	return nil
}
