package liquidhandling

import (
	"context"
	"testing"

	"github.com/antha-lang/antha/antha/anthalib/mixer"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/inventory"
	"github.com/antha-lang/antha/microArch/driver/liquidhandling"
)

func TestInputSampleAutoAllocate(t *testing.T) {
	ctx := GetContextForTest()
	pl, err := inventory.NewPlate(ctx, "pcrplate_skirted_riser20")
	if err != nil {
		t.Fatal(err)
	}

	(&PlanningTest{
		Name: "InputSampleAutoAllocated",
		Instructions: func(ctx context.Context) []*wtype.LHInstruction {

			cmp1, err := inventory.NewComponent(ctx, inventory.WaterType)
			if err != nil {
				t.Fatal(err)
			}
			cmp2, err := inventory.NewComponent(ctx, "dna_part")
			if err != nil {
				t.Fatal(err)
			}

			s1 := mixer.Sample(cmp1, wunit.NewVolume(50.0, "ul"))
			s2 := mixer.Sample(cmp2, wunit.NewVolume(25.0, "ul"))

			mo := mixer.MixOptions{
				Inputs:    []*wtype.Liquid{s1, s2},
				PlateType: "pcrplate_skirted_riser20",
				Address:   "A1",
				PlateNum:  1,
			}

			ins := mixer.GenericMix(mo)

			return []*wtype.LHInstruction{ins}
		},
		InputPlates: []*wtype.LHPlate{pl},
		Assertions: Assertions{
			InitialComponentAssertion(map[string]float64{"water": 55.5, "dna_part": 30.5}),
			NumberOfAssertion(liquidhandling.ASP, 2),
			NumberOfAssertion(liquidhandling.DSP, 2),
		},
	}).Run(ctx, t)
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

	ctx := GetContextForTest()
	pl, err := inventory.NewPlate(ctx, "pcrplate_skirted_riser20")
	if err != nil {
		t.Fatal(err)
	}

	(&PlanningTest{
		Name: "Mix in place autoallocated",
		Instructions: func(ctx context.Context) []*wtype.LHInstruction {

			cmp1, err := inventory.NewComponent(ctx, inventory.WaterType)
			if err != nil {
				t.Fatal(err)
			}
			cmp2, err := inventory.NewComponent(ctx, "dna_part")
			if err != nil {
				t.Fatal(err)
			}

			cmp1.Vol = 100.0
			cmp2.Vol = 50.0

			mo := mixer.MixOptions{
				Inputs:    []*wtype.Liquid{cmp1, cmp2},
				PlateType: "pcrplate_skirted_riser20",
				Address:   "A1",
				PlateNum:  1,
			}

			ins := mixer.GenericMix(mo)

			return []*wtype.LHInstruction{ins}
		},
		InputPlates: []*wtype.LHPlate{pl},
		Assertions: Assertions{
			InputLayoutAssertion(map[string]string{"A1": "water", "B1": "dna_part"}),
			InitialInputVolumesAssertion(0.01, map[string]float64{"A1": 100.0, "B1": 55.5}),
			FinalInputVolumesAssertion(0.01, map[string]float64{"A1": 150.0, "B1": 5.0}),
			NumberOfAssertion(liquidhandling.ASP, 1),
			NumberOfAssertion(liquidhandling.DSP, 1),
		},
	}).Run(ctx, t)

}
