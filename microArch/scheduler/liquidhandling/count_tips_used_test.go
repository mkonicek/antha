package liquidhandling

import (
	"github.com/antha-lang/antha/antha/anthalib/mixer"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"testing"
)

func TestTipCounting(t *testing.T) {
	ctx := GetContextForTest()
	PlanningTests{
		{
			Name: "single channel",
			Instructions: Mixes("pcrplate_skirted_riser", TestMixComponents{
				{
					LiquidName:    "water",
					VolumesByWell: ColumnWise(8, []float64{8.0, 8.0, 8.0, 8.0, 8.0, 8.0, 8.0, 8.0}),
					LiquidType:    wtype.LTSingleChannel,
					Sampler:       mixer.Sample,
				},
				{
					LiquidName:    "mastermix_sapI",
					VolumesByWell: ColumnWise(8, []float64{8.0, 8.0, 8.0, 8.0, 8.0, 8.0, 8.0, 8.0}),
					LiquidType:    wtype.LTSingleChannel,
					Sampler:       mixer.Sample,
				},
				{
					LiquidName:    "dna",
					VolumesByWell: ColumnWise(8, []float64{1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0}),
					LiquidType:    wtype.LTSingleChannel,
					Sampler:       mixer.Sample,
				},
			}),
			InputPlates:  []*wtype.LHPlate{GetTroughForTest()},
			OutputPlates: []*wtype.LHPlate{GetPlateForTest()},
			Assertions: Assertions{
				TipsUsedAssertion([]wtype.TipEstimate{{TipType: "DFL10 Tip Rack (PIPETMAX 8x20)", NTips: 8 * 3, NTipBoxes: 1}}),
			},
		},
		{
			Name: "multi channel",
			Instructions: Mixes("pcrplate_skirted_riser", TestMixComponents{
				{
					LiquidName:    "water",
					VolumesByWell: ColumnWise(8, []float64{8.0, 8.0, 8.0, 8.0, 8.0, 8.0, 8.0, 8.0}),
					LiquidType:    wtype.LTWater,
					Sampler:       mixer.Sample,
				},
				{
					LiquidName:    "mastermix_sapI",
					VolumesByWell: ColumnWise(8, []float64{8.0, 8.0, 8.0, 8.0, 8.0, 8.0, 8.0, 8.0}),
					LiquidType:    wtype.LTWater,
					Sampler:       mixer.Sample,
				},
				{
					LiquidName:    "dna",
					VolumesByWell: ColumnWise(8, []float64{1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0}),
					LiquidType:    wtype.LTWater,
					Sampler:       mixer.Sample,
				},
			}),
			InputPlates:  []*wtype.LHPlate{GetTroughForTest()},
			OutputPlates: []*wtype.LHPlate{GetPlateForTest()},
			Assertions: Assertions{
				TipsUsedAssertion([]wtype.TipEstimate{{TipType: "DFL10 Tip Rack (PIPETMAX 8x20)", NTips: 8 * 3, NTipBoxes: 1}}),
			},
		},
	}.Run(ctx, t)
}
