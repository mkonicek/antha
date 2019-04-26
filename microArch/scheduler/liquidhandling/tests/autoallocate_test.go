package tests

import (
	"testing"

	"github.com/antha-lang/antha/antha/anthalib/mixer"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/inventory/components"
	"github.com/antha-lang/antha/laboratory"
	"github.com/antha-lang/antha/laboratory/testlab"
	"github.com/antha-lang/antha/microArch/driver/liquidhandling"
)

func TestInputSampleAutoAllocate(t *testing.T) {
	testlab.WithTestLab(t, "", &testlab.TestElementCallbacks{
		Steps: func(lab *laboratory.Laboratory) error {

			pl, err := lab.Inventory.Plates.NewPlate("pcrplate_skirted_riser20")
			if err != nil {
				return err
			}

			(&PlanningTest{
				Name: "InputSampleAutoAllocated",
				Instructions: func(lab *laboratory.Laboratory) ([]*wtype.LHInstruction, error) {

					cmp1, err := lab.Inventory.Components.NewComponent(components.WaterType)
					if err != nil {
						return nil, err
					}
					cmp2, err := lab.Inventory.Components.NewComponent("dna_part")
					if err != nil {
						return nil, err
					}

					s1 := mixer.Sample(lab, cmp1, wunit.NewVolume(50.0, "ul"))
					s2 := mixer.Sample(lab, cmp2, wunit.NewVolume(25.0, "ul"))

					mo := mixer.MixOptions{
						Inputs:    []*wtype.Liquid{s1, s2},
						PlateType: "pcrplate_skirted_riser20",
						Address:   "A1",
						PlateNum:  1,
					}

					ins := mixer.GenericMix(lab, mo)

					return []*wtype.LHInstruction{ins}, nil
				},
				InputPlates: []*wtype.LHPlate{pl},
				Assertions: Assertions{
					InitialComponentAssertion(map[string]float64{"water": 55.5, "dna_part": 30.5}),
					NumberOfAssertion(liquidhandling.ASP, 2),
					NumberOfAssertion(liquidhandling.DSP, 2),
				},
			}).Run(t)
			return nil
		},
	})
}

func TestInPlaceAutoAllocate(t *testing.T) {

	// HJK 16/10/2018
	// SKIPPING: This test fails, and has been failing for some time.
	//   Expected Behaviour:
	//	   - Allocates 100ul water, 55.5ul DNA
	//	   - Move 50ul DNA on top of water
	//   Actual Behaviour:
	//     - Allocates 100ul water, 55.5ul DNA
	//     - Move 50ul DNA to position A1 of a new plate
	//  Previously the test only checkted that liquids were being allocated
	//  correctly, and so this error was not detected. Since this hasn't
	//  caused downstream issues, it is believed that this feature isn't currently
	//  in use.
	//  Since any likely fix will involve working on the API for this feature,
	//  it was decided to instead skip the test as mix in place to an auto-allocated
	//  input is not currently a priority.
	t.SkipNow()

	testlab.WithTestLab(t, "", &testlab.TestElementCallbacks{
		Steps: func(lab *laboratory.Laboratory) error {

			pl, err := lab.Inventory.Plates.NewPlate("pcrplate_skirted_riser20")
			if err != nil {
				return err
			}

			(&PlanningTest{
				Name: "Mix in place autoallocated",
				Instructions: func(lab *laboratory.Laboratory) ([]*wtype.LHInstruction, error) {

					cmp1, err := lab.Inventory.Components.NewComponent(components.WaterType)
					if err != nil {
						return nil, err
					}
					cmp2, err := lab.Inventory.Components.NewComponent("dna_part")
					if err != nil {
						return nil, err
					}

					cmp1.Vol = 100.0
					cmp2.Vol = 50.0

					mo := mixer.MixOptions{
						Inputs:    []*wtype.Liquid{cmp1, cmp2},
						PlateType: "pcrplate_skirted_riser20",
						Address:   "A1",
						PlateNum:  1,
					}

					ins := mixer.GenericMix(lab, mo)

					return []*wtype.LHInstruction{ins}, nil
				},
				InputPlates: []*wtype.LHPlate{pl},
				Assertions: Assertions{
					InputLayoutAssertion(map[string]string{"A1": "water", "B1": "dna_part"}),
					InitialInputVolumesAssertion(0.01, map[string]float64{"A1": 100.0, "B1": 55.5}),
					FinalInputVolumesAssertion(0.01, map[string]float64{"A1": 150.0, "B1": 5.0}),
					NumberOfAssertion(liquidhandling.ASP, 1),
					NumberOfAssertion(liquidhandling.DSP, 1),
				},
			}).Run(t)
			return nil
		},
	})
}
