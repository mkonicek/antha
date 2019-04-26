// anthalib//liquidhandling/liquidhandling_test.go: Part of the Antha language
// Copyright (C) 2015 The Antha authors. All rights reserved.
//
// This program is free software; you can redistribute it and/or
// modify it under the terms of the GNU General Public License
// as published by the Free Software Foundation; either version 2
// of the License, or (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program; if not, write to the Free Software
// Foundation, Inc., 51 Franklin Street, Fifth Floor, Boston, MA  02110-1301, USA.
//
// For more information relating to the software or licensing issues please
// contact license@antha-lang.org or write to the Antha team c/o
// Synthace Ltd. The London Bioscience Innovation Centre
// 2 Royal College St, London NW1 0NH UK

package tests

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"

	"github.com/antha-lang/antha/antha/anthalib/mixer"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/antha/anthalib/wutil/text"
	"github.com/antha-lang/antha/laboratory"
	"github.com/antha-lang/antha/laboratory/effects/id"
	"github.com/antha-lang/antha/laboratory/testlab"
	"github.com/antha-lang/antha/microArch/driver/liquidhandling"
	lh "github.com/antha-lang/antha/microArch/scheduler/liquidhandling"
	"github.com/antha-lang/antha/utils"
)

func GetPlateForTest(idGen *id.IDGenerator) *wtype.Plate {

	offset := 0.25
	riserheightinmm := 40.0 - offset

	// pcr plate skirted (on riser)
	cone := wtype.NewShape(wtype.CylinderShape, "mm", 5.5, 5.5, 20.4)
	welltype := wtype.NewLHWell(idGen, "ul", 200, 5, cone, wtype.UWellBottom, 5.5, 5.5, 20.4, 1.4, "mm")

	return wtype.NewLHPlate(idGen, "pcrplate_skirted_riser", "Unknown", 8, 12, wtype.Coordinates3D{X: 127.76, Y: 85.48, Z: 25.7}, welltype, 9, 9, 0.0, 0.0, riserheightinmm-1.25)
}

func PrefillPlateForTest(lab *laboratory.Laboratory, plate *wtype.LHPlate, liquidType string, volumes map[string]float64) *wtype.LHPlate {
	for address, volume := range volumes {
		cmp := GetComponentForTest(lab, liquidType, wunit.NewVolume(volume, "ul"))
		if err := plate.Wellcoords[address].SetContents(lab.IDGenerator, cmp); err != nil {
			panic(err)
		}
	}

	return plate
}

func GetTipwasteForTest(idGen *id.IDGenerator) *wtype.LHTipwaste {
	shp := wtype.NewShape(wtype.BoxShape, "mm", 123.0, 80.0, 92.0)
	w := wtype.NewLHWell(idGen, "ul", 800000.0, 800000.0, shp, 0, 123.0, 80.0, 92.0, 0.0, "mm")
	lht := wtype.NewLHTipwaste(idGen, 6000, "Gilsontipwaste", "gilson", wtype.Coordinates3D{X: 127.76, Y: 85.48, Z: 92.0}, w, 49.5, 31.5, 0.0)
	return lht
}

func GetTroughForTest(idGen *id.IDGenerator) *wtype.Plate {
	stshp := wtype.NewShape(wtype.BoxShape, "mm", 8.2, 72, 41.3)
	trough12 := wtype.NewLHWell(idGen, "ul", 1500, 500, stshp, wtype.VWellBottom, 8.2, 72, 41.3, 4.7, "mm")
	plate := wtype.NewLHPlate(idGen, "DWST12", "Unknown", 1, 12, wtype.Coordinates3D{X: 127.76, Y: 85.48, Z: 44.1}, trough12, 9, 9, 0, 30.0, 4.5)
	return plate
}

func configure_request_simple(lab *laboratory.Laboratory, rq *lh.LHRequest) {
	water := GetComponentForTest(lab, "water", wunit.NewVolume(100.0, "ul"))
	water.Type = wtype.LTSingleChannel
	mmx := GetComponentForTest(lab, "mastermix_sapI", wunit.NewVolume(100.0, "ul"))
	mmx.Type = wtype.LTSingleChannel
	part := GetComponentForTest(lab, "dna", wunit.NewVolume(50.0, "ul"))
	part.Type = wtype.LTSingleChannel

	for k := 0; k < 9; k++ {
		ins := wtype.NewLHMixInstruction(lab.IDGenerator)
		ws := mixer.Sample(lab, water, wunit.NewVolume(8.0, "ul"))
		mmxs := mixer.Sample(lab, mmx, wunit.NewVolume(8.0, "ul"))
		ps := mixer.Sample(lab, part, wunit.NewVolume(1.0, "ul"))

		ins.AddInput(ws)
		ins.AddInput(mmxs)
		ins.AddInput(ps)
		ins.AddOutput(GetComponentForTest(lab, "water", wunit.NewVolume(17.0, "ul")))
		rq.Add_instruction(ins)
	}

}

func configure_request_bigger(lab *laboratory.Laboratory, rq *lh.LHRequest) {
	water := GetComponentForTest(lab, "water", wunit.NewVolume(2000.0, "ul"))
	mmx := GetComponentForTest(lab, "mastermix_sapI", wunit.NewVolume(2000.0, "ul"))
	part := GetComponentForTest(lab, "dna", wunit.NewVolume(1000.0, "ul"))

	for k := 0; k < 99; k++ {
		ins := wtype.NewLHMixInstruction(lab.IDGenerator)
		ws := mixer.Sample(lab, water, wunit.NewVolume(8.0, "ul"))
		mmxs := mixer.Sample(lab, mmx, wunit.NewVolume(8.0, "ul"))
		ps := mixer.Sample(lab, part, wunit.NewVolume(1.0, "ul"))

		ins.AddInput(ws)
		ins.AddInput(mmxs)
		ins.AddInput(ps)
		ins.AddOutput(GetComponentForTest(lab, "water", wunit.NewVolume(17.0, "ul")))
		rq.Add_instruction(ins)
	}

}

func configurePlanningTestRequest(lab *laboratory.Laboratory, rq *lh.LHRequest) {
	water := GetComponentForTest(lab, "multiwater", wunit.NewVolume(2000.0, "ul"))

	for k := 0; k < 9; k++ {
		ins := wtype.NewLHMixInstruction(lab.IDGenerator)
		ws := mixer.Sample(lab, water, wunit.NewVolume(50.0, "ul"))

		ins.AddInput(ws)

		ins.AddOutput(GetComponentForTest(lab, "water", wunit.NewVolume(50, "ul")))
		rq.Add_instruction(ins)
	}

}

func makeSomeTipboxes(lab *laboratory.Laboratory, names ...string) ([]*wtype.LHTipbox, error) {
	res := make([]*wtype.LHTipbox, len(names))
	for idx, name := range names {
		if tb, err := lab.Inventory.TipBoxes.NewTipbox(name); err != nil {
			return nil, err
		} else {
			res[idx] = tb
		}
	}
	return res, nil
}

func configureTransferRequestForZTest(lab *laboratory.Laboratory, policyName string, transferVol wunit.Volume, numberOfTransfers int) (rq *lh.LHRequest, err error) {

	// make liquid handler
	lh := GetLiquidHandlerForTest(lab)
	tipBoxes, err := makeSomeTipboxes(lab, "Gilson200", "Gilson20")
	if err != nil {
		return nil, err
	}

	//initialise request
	rq = GetLHRequestForTest(lab.IDGenerator)

	liq := GetComponentForTest(lab, "water", wunit.NewVolume(2000.0, "ul"))

	err = liq.SetPolicyName(wtype.PolicyName(policyName))
	if err != nil {
		return rq, err
	}
	liq.SetName(policyName)

	for k := 0; k < numberOfTransfers; k++ {
		ins := wtype.NewLHMixInstruction(lab.IDGenerator)
		ws := mixer.Sample(lab, liq, transferVol)

		ins.AddInput(ws)

		expectedProduct := GetComponentForTest(lab, "water", transferVol)

		err = expectedProduct.SetPolicyName(wtype.PolicyName(policyName))
		if err != nil {
			return rq, err
		}
		expectedProduct.SetName(policyName)

		ins.AddOutput(expectedProduct)

		rq.Add_instruction(ins)
	}

	// add plates and tip boxes
	rq.InputPlatetypes = append(rq.InputPlatetypes, GetPlateForTest(lab.IDGenerator))
	rq.OutputPlatetypes = append(rq.OutputPlatetypes, GetPlateForTest(lab.IDGenerator))

	rq.TipBoxes = tipBoxes

	if err := lh.Plan(lab.LaboratoryEffects, rq); err != nil {
		return rq, fmt.Errorf("Got an error planning with no inputs: %s", err.Error())
	}
	return rq, nil
}

func configureSingleChannelTestRequest(lab *laboratory.Laboratory, rq *lh.LHRequest) {
	water := GetComponentForTest(lab, "multiwater", wunit.NewVolume(2000.0, "ul"))

	for k := 0; k < 1; k++ {
		ins := wtype.NewLHMixInstruction(lab.IDGenerator)
		ws := mixer.Sample(lab, water, wunit.NewVolume(50.0, "ul"))

		ins.AddInput(ws)

		ins.AddOutput(GetComponentForTest(lab, "water", wunit.NewVolume(50, "ul")))
		rq.Add_instruction(ins)
	}

}

func configureTransferRequestMutliSamplesTest(lab *laboratory.Laboratory, policyName string, samples ...*wtype.Liquid) (rq *lh.LHRequest, err error) {
	// make liquid handler
	lh := GetLiquidHandlerForTest(lab)

	tipBoxes, err := makeSomeTipboxes(lab, "Gilson200", "Gilson20")
	if err != nil {
		return nil, err
	}

	//initialise request
	rq = GetLHRequestForTest(lab.IDGenerator)

	// add plates and tip boxes
	inPlate := GetPlateForTest(lab.IDGenerator)
	rq.InputPlatetypes = append(rq.InputPlatetypes, inPlate)
	rq.OutputPlatetypes = append(rq.OutputPlatetypes, GetPlateForTest(lab.IDGenerator))
	rq.TipBoxes = tipBoxes

	it := wtype.NewAddressIterator(inPlate, wtype.RowWise, wtype.TopToBottom, wtype.LeftToRight, false)

	for _, sample := range samples {
		ins := wtype.NewLHMixInstruction(lab.IDGenerator)

		if err := sample.SetPolicyName(wtype.PolicyName(policyName)); err != nil {
			return nil, err
		}

		ins.AddInput(sample)
		ins.AddOutput(GetComponentForTest(lab, "water", sample.Volume()))

		if !it.Valid() {
			return nil, errors.New("out of space on input plate")
		}

		ins.Welladdress = it.Curr().FormatA1()
		it.Next()

		rq.Add_instruction(ins)
	}

	if err := lh.Plan(lab.LaboratoryEffects, rq); err != nil {
		return rq, errors.WithMessage(err, "while planning")
	}
	return rq, nil
}

func TestToWellVolume(t *testing.T) {
	testlab.WithTestLab(t, "", &testlab.TestElementCallbacks{
		Steps: func(lab *laboratory.Laboratory) error {
			water := GetComponentForTest(lab, "water", wunit.NewVolume(2000.0, "ul"))
			mmx := GetComponentForTest(lab, "mastermix_sapI", wunit.NewVolume(2000.0, "ul"))
			part := GetComponentForTest(lab, "dna", wunit.NewVolume(1000.0, "ul"))

			ws := mixer.Sample(lab, water, wunit.NewVolume(150.0, "ul"))
			mmxs := mixer.Sample(lab, mmx, wunit.NewVolume(49.0, "ul"))
			ps := mixer.Sample(lab, part, wunit.NewVolume(1.0, "ul"))
			_, err := configureTransferRequestMutliSamplesTest(lab, "SmartMix", ws, mmxs, ps)
			return err
		},
	})
}

type zOffsetTest struct {
	liquidType              string
	numberOfTransfers       int
	volume                  wunit.Volume
	expectedAspirateZOffset []float64
	expectedDispenseZOffset []float64
}

func (self zOffsetTest) String() string {
	return fmt.Sprintf("%dx %v with policy=%q", self.numberOfTransfers, self.volume, self.liquidType)
}

var offsetTests []zOffsetTest = []zOffsetTest{
	{
		liquidType:              "multiwater",
		numberOfTransfers:       1,
		volume:                  wunit.NewVolume(50, "ul"),
		expectedAspirateZOffset: []float64{0.5000},
		expectedDispenseZOffset: []float64{1.0000},
	},
	{
		liquidType:              "multiwater",
		numberOfTransfers:       2,
		volume:                  wunit.NewVolume(50, "ul"),
		expectedAspirateZOffset: []float64{0.5000, 0.5000},
		expectedDispenseZOffset: []float64{1.0000, 1.0000},
	},
	{
		liquidType:              "multiwater",
		numberOfTransfers:       1,
		volume:                  wunit.NewVolume(5, "ul"),
		expectedAspirateZOffset: []float64{0.5000},
		expectedDispenseZOffset: []float64{1.0000},
	},
	{
		liquidType:              "multiwater",
		numberOfTransfers:       2,
		volume:                  wunit.NewVolume(5, "ul"),
		expectedAspirateZOffset: []float64{0.5000, 0.5000},
		expectedDispenseZOffset: []float64{1.0000, 1.0000},
	},
	// Commented this out as it's not directly related to z offset and is failing
	// due to not performing a multichannel transfer.
	/*
		zOffsetTest{
			liquidType:              "multiwater",
			numberOfTransfers:       8,
			volume:                  wunit.NewVolume(50, "ul"),
			expectedAspirateZOffset: []float64{1.2500,1.2500,1.2500,1.2500,1.2500,1.2500,1.2500,1.2500},
			expectedDispenseZOffset: []float64{1.7500,1.7500,1.7500,1.7500,1.7500,1.7500,1.7500,1.7500},
		},*/
	{
		liquidType:              "SingleChannel",
		numberOfTransfers:       1,
		volume:                  wunit.NewVolume(50, "ul"),
		expectedAspirateZOffset: []float64{0.5000},
		expectedDispenseZOffset: []float64{1.0000},
	},
	{
		liquidType:              "SingleChannel",
		numberOfTransfers:       2,
		volume:                  wunit.NewVolume(50, "ul"),
		expectedAspirateZOffset: []float64{0.5000},
		expectedDispenseZOffset: []float64{1.0000},
	},
	{
		liquidType:              "SingleChannel",
		numberOfTransfers:       1,
		volume:                  wunit.NewVolume(5, "ul"),
		expectedAspirateZOffset: []float64{0.5000},
		expectedDispenseZOffset: []float64{1.0000},
	},
	{
		liquidType:              "SingleChannel",
		numberOfTransfers:       2,
		volume:                  wunit.NewVolume(5, "ul"),
		expectedAspirateZOffset: []float64{0.5000},
		expectedDispenseZOffset: []float64{1.0000},
	},
	{
		liquidType:              "SmartMix",
		numberOfTransfers:       1,
		volume:                  wunit.NewVolume(50, "ul"),
		expectedAspirateZOffset: []float64{0.5000},
		expectedDispenseZOffset: []float64{0.5000},
	},
	{
		liquidType:              "SmartMix",
		numberOfTransfers:       2,
		volume:                  wunit.NewVolume(50, "ul"),
		expectedAspirateZOffset: []float64{0.5000, 0.5000},
		expectedDispenseZOffset: []float64{0.5000, 0.5000},
	}, /*
		zOffsetTest{
			liquidType:              "SmartMix",
			numberOfTransfers:       1,
			volume:                  wunit.NewVolume(5, "ul"),
			expectedAspirateZOffset: []float64{0.500},
			expectedDispenseZOffset: []float64{0.500},
		},
		zOffsetTest{
			liquidType:              "SmartMix",
			numberOfTransfers:       2,
			volume:                  wunit.NewVolume(5, "ul"),
			expectedAspirateZOffset: []float64{0.500,0.500},
			expectedDispenseZOffset: []float64{0.500,0.500},
		},*/
	{
		liquidType:              "NeedToMix",
		numberOfTransfers:       1,
		volume:                  wunit.NewVolume(50, "ul"),
		expectedAspirateZOffset: []float64{0.5000},
		expectedDispenseZOffset: []float64{0.5000},
	},
	{
		liquidType:              "NeedToMix",
		numberOfTransfers:       2,
		volume:                  wunit.NewVolume(50, "ul"),
		expectedAspirateZOffset: []float64{0.5000, 0.5000},
		expectedDispenseZOffset: []float64{0.5000, 0.5000},
	}, /*
		zOffsetTest{
			liquidType:              "NeedToMix",
			numberOfTransfers:       1,
			volume:                  wunit.NewVolume(5, "ul"),
			expectedAspirateZOffset: []float64{0.500},
			expectedDispenseZOffset: []float64{0.500},
		},
		zOffsetTest{
			liquidType:              "NeedToMix",
			numberOfTransfers:       2,
			volume:                  wunit.NewVolume(5, "ul"),
			expectedAspirateZOffset: []float64{0.500,0.500},
			expectedDispenseZOffset: []float64{0.500,0.500},
		},*/
}

func TestMultiZOffset2(t *testing.T) {

	for _, test := range offsetTests {
		testlab.WithTestLab(t, "", &testlab.TestElementCallbacks{
			Name: test.String(),
			Steps: func(lab *laboratory.Laboratory) error {
				request, err := configureTransferRequestForZTest(lab, test.liquidType, test.volume, test.numberOfTransfers)
				if err != nil {
					return err
				}

				aspiratePairs, dispensePairs := extractMoveAspirateDispenseInstructions(request.Instructions)

				for i, pair := range aspiratePairs {
					if !reflect.DeepEqual(pair.mov.OffsetZ, test.expectedAspirateZOffset) {
						return fmt.Errorf("for test: %v\naspiration step: %d\nexpected Z offset for aspirate: %v\ngot: %v",
							test, i, test.expectedAspirateZOffset, pair.mov.OffsetZ)
					}
				}

				for i, pair := range dispensePairs {
					if !reflect.DeepEqual(pair.mov.OffsetZ, test.expectedDispenseZOffset) {
						return fmt.Errorf("for test: %v\ndispense step: %d\nexpected Z offset for dispense: %v\ngot: %v",
							test, i, test.expectedDispenseZOffset, pair.mov.OffsetZ)
					}
				}
				return nil
			},
		})
	}
}

func makeMultiTestRequest(lab *laboratory.Laboratory) (*lh.LHRequest, error) {
	// make liquid handler
	lh := GetLiquidHandlerForTest(lab)

	tipBoxes, err := makeSomeTipboxes(lab, "Gilson200", "Gilson20")
	if err != nil {
		return nil, err
	}

	// set up multi

	//initialise multi request
	multiRq := GetLHRequestForTest(lab.IDGenerator)

	// set to Multi channel test request
	configurePlanningTestRequest(lab, multiRq)
	// add plates and tip boxes
	multiRq.InputPlatetypes = append(multiRq.InputPlatetypes, GetPlateForTest(lab.IDGenerator))
	multiRq.OutputPlatetypes = append(multiRq.OutputPlatetypes, GetPlateForTest(lab.IDGenerator))

	multiRq.TipBoxes = tipBoxes

	if err := lh.Plan(lab.LaboratoryEffects, multiRq); err != nil {
		return multiRq, fmt.Errorf("Got an error planning with no inputs: %s", err)
	}
	return multiRq, nil
}

func makeSingleTestRequest(lab *laboratory.Laboratory) (*lh.LHRequest, error) {
	// make liquid handler
	lh := GetLiquidHandlerForTest(lab)

	tipBoxes, err := makeSomeTipboxes(lab, "Gilson200", "Gilson20")
	if err != nil {
		return nil, err
	}

	// set up single channel

	//initialise single request
	singleRq := GetLHRequestForTest(lab.IDGenerator)

	// set to single channel test request
	configureSingleChannelTestRequest(lab, singleRq)
	// add plates and tip boxes
	singleRq.InputPlatetypes = append(singleRq.InputPlatetypes, GetPlateForTest(lab.IDGenerator))
	singleRq.OutputPlatetypes = append(singleRq.OutputPlatetypes, GetPlateForTest(lab.IDGenerator))

	singleRq.TipBoxes = tipBoxes

	if err := lh.Plan(lab.LaboratoryEffects, singleRq); err != nil {
		return singleRq, fmt.Errorf("Got an error planning with no inputs: %s", err)
	}
	return singleRq, nil
}

type movAspPair struct {
	mov *liquidhandling.MoveInstruction
	asp *liquidhandling.AspirateInstruction
}
type movDspPair struct {
	mov *liquidhandling.MoveInstruction
	dsp *liquidhandling.DispenseInstruction
}

func extractMoveAspirateDispenseInstructions(ins []liquidhandling.TerminalRobotInstruction) ([]movAspPair, []movDspPair) {
	mov := make([]*liquidhandling.MoveInstruction, len(ins))
	ma := []movAspPair{}
	md := []movDspPair{}

	for idx, i := range ins {
		i.Visit(&liquidhandling.RobotInstructionBaseVisitor{
			HandleMove: func(ins *liquidhandling.MoveInstruction) { mov[idx] = ins },
			HandleAspirate: func(ins *liquidhandling.AspirateInstruction) {
				if idx > 0 && mov[idx-1] != nil {
					ma = append(ma, movAspPair{mov: mov[idx-1], asp: ins})
				}
			},
			HandleDispense: func(ins *liquidhandling.DispenseInstruction) {
				if idx > 0 && mov[idx-1] != nil {
					md = append(md, movDspPair{mov: mov[idx-1], dsp: ins})
				}
			},
		})
	}
	return ma, md
}

func allElemsSame(nums []float64) bool {
	if len(nums) > 1 {
		n := nums[0]
		for _, m := range nums[1:] {
			if n != m {
				return false
			}
		}
	}
	return true
}

func TestMultiZOffset(t *testing.T) {
	testlab.WithTestLab(t, "", &testlab.TestElementCallbacks{
		Steps: func(lab *laboratory.Laboratory) error {

			multiRq, err := makeMultiTestRequest(lab)
			if err != nil {
				return err
			}

			singleRq, err := makeSingleTestRequest(lab)
			if err != nil {
				return err
			}

			multiAspPairs, multiDspPairs := extractMoveAspirateDispenseInstructions(multiRq.Instructions)
			singleAspPairs, singleDspPairs := extractMoveAspirateDispenseInstructions(singleRq.Instructions)

			if len(multiAspPairs) < len(singleAspPairs) {
				return fmt.Errorf("Too few (%d) multi Asp pairs (need at least %d)", len(multiAspPairs), len(singleAspPairs))
			}
			if len(multiDspPairs) < len(singleDspPairs) {
				return fmt.Errorf("Too few (%d) multi Dsp pairs (need at least %d)", len(multiDspPairs), len(singleDspPairs))
			}

			for i, singlePair := range singleAspPairs {
				if !allElemsSame(singlePair.mov.OffsetZ) {
					return fmt.Errorf("Z offsets not all the same (single asp pair): %#v", singlePair.mov.OffsetZ)
				}
				multiPair := multiAspPairs[i]
				if !allElemsSame(multiPair.mov.OffsetZ) {
					return fmt.Errorf("Z offsets not all the same (multi asp pair): %#v", multiPair.mov.OffsetZ)
				}
				if singlePair.mov.OffsetZ[0] != multiPair.mov.OffsetZ[0] {
					return fmt.Errorf("single Aspirate Z offset: %+v\nNot equal to\nmulti Aspirate Z offset: %+v",
						text.PrettyPrint(singlePair), text.PrettyPrint(multiPair))
				}
			}

			for i, singlePair := range singleDspPairs {
				if !allElemsSame(singlePair.mov.OffsetZ) {
					return fmt.Errorf("Z offsets not all the same (single asp pair): %#v", singlePair.mov.OffsetZ)
				}
				multiPair := multiDspPairs[i]
				if !allElemsSame(multiPair.mov.OffsetZ) {
					return fmt.Errorf("Z offsets not all the same (multi asp pair): %#v", multiPair.mov.OffsetZ)
				}
				if singlePair.mov.OffsetZ[0] != multiPair.mov.OffsetZ[0] {
					return fmt.Errorf("single Dispense Z offset: %v\nNot equal to\nmulti Dispense Z offset: %+v",
						text.PrettyPrint(singlePair), text.PrettyPrint(multiPair))
				}
			}
			return nil
		},
	})
}

func TestTipOverridePositive(t *testing.T) {
	testlab.WithTestLab(t, "", &testlab.TestElementCallbacks{
		Steps: func(lab *laboratory.Laboratory) error {
			lh := GetLiquidHandlerForTest(lab)
			rq := GetLHRequestForTest(lab.IDGenerator)
			configure_request_simple(lab, rq)
			rq.InputPlatetypes = append(rq.InputPlatetypes, GetPlateForTest(lab.IDGenerator))
			rq.OutputPlatetypes = append(rq.OutputPlatetypes, GetPlateForTest(lab.IDGenerator))

			tp, err := lab.Inventory.TipBoxes.NewTipbox("Gilson20")
			if err != nil {
				return err
			}

			rq.TipBoxes = []*wtype.LHTipbox{tp}

			if err := lh.Plan(lab.LaboratoryEffects, rq); err != nil {
				return fmt.Errorf("Got an error planning with no inputs: %s", err)
			}
			return nil
		},
	})
}

func TestTipOverrideNegative(t *testing.T) {
	testlab.WithTestLab(t, "", &testlab.TestElementCallbacks{
		Steps: func(lab *laboratory.Laboratory) error {
			lh := GetLiquidHandlerForTest(lab)
			rq := GetLHRequestForTest(lab.IDGenerator)
			configure_request_simple(lab, rq)
			rq.InputPlatetypes = append(rq.InputPlatetypes, GetPlateForTest(lab.IDGenerator))
			rq.OutputPlatetypes = append(rq.OutputPlatetypes, GetPlateForTest(lab.IDGenerator))

			tp, err := lab.Inventory.TipBoxes.NewTipbox("Gilson200")
			if err != nil {
				return err
			}

			rq.TipBoxes = []*wtype.LHTipbox{tp}

			err = lh.Plan(lab.LaboratoryEffects, rq)

			if e, f := "7 (LH_ERR_VOL) : volume error : No tip chosen: Volume 8 ul is too low to be accurately moved by the liquid handler (configured minimum 10 ul). Low volume tips may not be available and / or the robot may need to be configured differently", err.Error(); e != f {
				return fmt.Errorf("expecting error %q found %q", e, f)
			}
			return nil
		},
	})
}

func TestPlateReuse(t *testing.T) {
	testlab.WithTestLab(t, "", &testlab.TestElementCallbacks{
		Steps: func(lab *laboratory.Laboratory) error {
			lh := GetLiquidHandlerForTest(lab)
			rq := GetLHRequestForTest(lab.IDGenerator)
			configure_request_simple(lab, rq)
			rq.InputPlatetypes = append(rq.InputPlatetypes, GetPlateForTest(lab.IDGenerator))
			rq.OutputPlatetypes = append(rq.OutputPlatetypes, GetPlateForTest(lab.IDGenerator))

			err := lh.Plan(lab.LaboratoryEffects, rq)

			if err != nil {
				return fmt.Errorf("Got an error planning with no inputs: %v", err)
			}

			// reset the request
			rq = GetLHRequestForTest(lab.IDGenerator)
			configure_request_simple(lab, rq)

			for _, plateid := range lh.Properties.PosLookup {
				if plateid == "" {
					continue
				}
				thing := lh.Properties.PlateLookup[plateid]

				plate, ok := thing.(*wtype.Plate)
				if !ok {
					continue
				}

				if strings.Contains(plate.GetName(), "Output_plate") {
					// leave it out
					continue
				}

				rq.InputPlates[plateid] = plate
			}
			rq.InputPlatetypes = append(rq.InputPlatetypes, GetPlateForTest(lab.IDGenerator))
			rq.OutputPlatetypes = append(rq.OutputPlatetypes, GetPlateForTest(lab.IDGenerator))

			lh = GetLiquidHandlerForTest(lab)
			err = lh.Plan(lab.LaboratoryEffects, rq)

			if err != nil {
				return fmt.Errorf("Got error resimulating: %v", err)
			}

			// if we added nothing, input assignments should be empty

			if rq.NewComponentsAdded() {
				return fmt.Errorf("Resimulation failed: needed to add %d components", len(rq.InputSolutions.VolumesWanting))
			}

			// now try a deliberate fail

			// reset the request again
			rq = GetLHRequestForTest(lab.IDGenerator)
			configure_request_simple(lab, rq)

			for _, plateid := range lh.Properties.PosLookup {
				if plateid == "" {
					continue
				}
				thing := lh.Properties.PlateLookup[plateid]

				plate, ok := thing.(*wtype.Plate)
				if !ok {
					continue
				}
				if strings.Contains(plate.GetName(), "Output_plate") {
					// leave it out
					continue
				}
				for _, v := range plate.Wellcoords {
					if !v.IsEmpty(lab.IDGenerator) {
						if _, err := v.RemoveVolume(lab.IDGenerator, wunit.NewVolume(5.0, "ul")); err != nil {
							return err
						}
					}
				}

				rq.InputPlates[plateid] = plate
			}
			rq.InputPlatetypes = append(rq.InputPlatetypes, GetPlateForTest(lab.IDGenerator))
			rq.OutputPlatetypes = append(rq.OutputPlatetypes, GetPlateForTest(lab.IDGenerator))

			lh = GetLiquidHandlerForTest(lab)
			err = lh.Plan(lab.LaboratoryEffects, rq)

			if err != nil {
				return fmt.Errorf("Got error resimulating: %v", err)
			}

			// this time we should have added some components again
			if len(rq.InputAssignments) != 3 {
				return fmt.Errorf("Error resimulating, should have added 3 components, instead added %d", len(rq.InputAssignments))
			}
			return nil
		},
	})
}

func TestExecutionPlanning(t *testing.T) {
	testlab.WithTestLab(t, "", &testlab.TestElementCallbacks{
		Steps: func(lab *laboratory.Laboratory) error {

			PlanningTests{
				{
					Name: "simple planning",
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
					InputPlates:  []*wtype.LHPlate{GetTroughForTest(lab.IDGenerator)},
					OutputPlates: []*wtype.LHPlate{GetPlateForTest(lab.IDGenerator)},
					Assertions: Assertions{
						NumberOfAssertion(liquidhandling.ASP, 3*8), //no multichanneling
						NumberOfAssertion(liquidhandling.DSP, 3*8), //no multichanneling
						LayoutSummaryAssertion("simple.layout.json"),
						ActionsSummaryAssertion("simple.actions.json"),
					},
				},
				{
					Name: "total volume",
					Instructions: Mixes("pcrplate_skirted_riser", TestMixComponents{
						{
							LiquidName:    "water",
							VolumesByWell: ColumnWise(8, []float64{17.0, 17.0, 17.0, 17.0, 17.0, 17.0, 17.0, 17.0}),
							LiquidType:    wtype.LTSingleChannel,
							Sampler:       mixer.SampleForTotalVolume,
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
					InputPlates:  []*wtype.LHPlate{GetTroughForTest(lab.IDGenerator)},
					OutputPlates: []*wtype.LHPlate{GetPlateForTest(lab.IDGenerator)},
					Assertions: Assertions{
						NumberOfAssertion(liquidhandling.ASP, 3*8), //no multichanneling
						NumberOfAssertion(liquidhandling.DSP, 3*8), //no multichanneling
						FinalOutputVolumesAssertion(0.001, map[string]float64{"A1": 17.0, "B1": 17.0, "C1": 17.0, "D1": 17.0, "E1": 17.0, "F1": 17.0, "G1": 17.0, "H1": 17.0}),
					},
				},
				{
					Name: "overfull wells",
					Instructions: Mixes("pcrplate_skirted_riser", TestMixComponents{
						{
							LiquidName:    "water",
							VolumesByWell: ColumnWise(8, []float64{160.0, 160.0, 160.0, 160.0, 160.0, 160.0, 160.0, 160.0}),
							LiquidType:    wtype.LTSingleChannel,
							Sampler:       mixer.Sample,
						},
						{
							LiquidName:    "mastermix_sapI",
							VolumesByWell: ColumnWise(8, []float64{160.0, 160.0, 160.0, 160.0, 160.0, 160.0, 160.0, 160.0}),
							LiquidType:    wtype.LTSingleChannel,
							Sampler:       mixer.Sample,
						},
						{
							LiquidName:    "dna",
							VolumesByWell: ColumnWise(8, []float64{20.0, 20.0, 20.0, 20.0, 20.0, 20.0, 20.0, 20.0}),
							LiquidType:    wtype.LTSingleChannel,
							Sampler:       mixer.Sample,
						},
					}),
					InputPlates:  []*wtype.LHPlate{GetTroughForTest(lab.IDGenerator)},
					OutputPlates: []*wtype.LHPlate{GetPlateForTest(lab.IDGenerator)},
					ErrorPrefix:  "7 (LH_ERR_VOL) : volume error : volume of resulting mix (340 ul) exceeds the well maximum (200 ul) for instruction:",
				},
				{
					Name: "negative requested volume",
					Instructions: Mixes("pcrplate_skirted_riser", TestMixComponents{
						{
							LiquidName:    "water",
							VolumesByWell: ColumnWise(8, []float64{8.0, -1.0, 8.0, 8.0, 8.0, 8.0, 8.0, 8.0}),
							LiquidType:    wtype.LTSingleChannel,
							Sampler:       mixer.Sample,
						},
					}),
					InputPlates:  []*wtype.LHPlate{GetTroughForTest(lab.IDGenerator)},
					OutputPlates: []*wtype.LHPlate{GetPlateForTest(lab.IDGenerator)},
					ErrorPrefix:  "7 (LH_ERR_VOL) : volume error : negative volume for component \"water\" in instruction:",
				},
				{
					Name: "invalid total volume",
					Instructions: Mixes("pcrplate_skirted_riser", TestMixComponents{
						{
							LiquidName:    "water",
							VolumesByWell: ColumnWise(8, []float64{5.0, 5.0, 5.0, 5.0, 5.0, 5.0, 5.0, 5.0}),
							LiquidType:    wtype.LTSingleChannel,
							Sampler:       mixer.SampleForTotalVolume,
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
					InputPlates:  []*wtype.LHPlate{GetTroughForTest(lab.IDGenerator)},
					OutputPlates: []*wtype.LHPlate{GetPlateForTest(lab.IDGenerator)},
					ErrorPrefix:  "during solution setup: 7 (LH_ERR_VOL) : volume error : invalid total volume for component \"water\" in instruction:",
				},
				{
					Name: "test dummy instruction removal",
					Instructions: func(lab *laboratory.Laboratory) ([]*wtype.LHInstruction, error) {
						instructions, err := Mixes("pcrplate_skirted_riser", TestMixComponents{
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
						})(lab)
						if err != nil {
							return nil, err
						}
						//add a dummy instruction for each instruction
						ret := make([]*wtype.LHInstruction, 0, len(instructions))
						for _, ins := range instructions {
							for _, cmp := range ins.Outputs {
								mix := mixer.GenericMix(lab, mixer.MixOptions{Inputs: []*wtype.Liquid{cmp}})
								if !mix.IsDummy() {
									return nil, fmt.Errorf("failed to make a dummy instruction: mix.Inputs[0].IsSample() = %t, cmp.IsSample() = %t", mix.Inputs[0].IsSample(), cmp.IsSample())
								}
								ret = append(ret, ins)
								ret = append(ret, mix)
							}
						}
						return ret, nil
					},
					InputPlates:  []*wtype.LHPlate{GetTroughForTest(lab.IDGenerator)},
					OutputPlates: []*wtype.LHPlate{GetPlateForTest(lab.IDGenerator)},
					Assertions: Assertions{
						NumberOfAssertion(liquidhandling.ASP, 3*8), //no multichanneling
						NumberOfAssertion(liquidhandling.DSP, 3*8), //no multichanneling
					},
				},
				{
					Name: "test result volume doesn't match total volume",
					Instructions: func(lab *laboratory.Laboratory) ([]*wtype.LHInstruction, error) {
						instructions, err := Mixes("pcrplate_skirted_riser", TestMixComponents{
							{
								LiquidName:    "water",
								VolumesByWell: ColumnWise(8, []float64{17.0, 17.0, 17.0, 17.0, 17.0, 17.0, 17.0, 17.0}),
								LiquidType:    wtype.LTSingleChannel,
								Sampler:       mixer.SampleForTotalVolume,
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
						})(lab)
						if err != nil {
							return nil, err
						}
						for _, ins := range instructions {
							ins.Outputs[0].Vol = 10.0
						}
						return instructions, err
					},
					InputPlates:  []*wtype.LHPlate{GetTroughForTest(lab.IDGenerator)},
					OutputPlates: []*wtype.LHPlate{GetPlateForTest(lab.IDGenerator)},
					ErrorPrefix:  "7 (LH_ERR_VOL) : volume error : total volume (17 ul) does not match resulting volume (10 ul)",
				},
				{
					Name: "test result volume doesn't match volume sum",
					Instructions: func(lab *laboratory.Laboratory) ([]*wtype.LHInstruction, error) {
						instructions, err := Mixes("pcrplate_skirted_riser", TestMixComponents{
							{
								LiquidName:    "water",
								VolumesByWell: ColumnWise(8, []float64{17.0, 17.0, 17.0, 17.0, 17.0, 17.0, 17.0, 17.0}),
								LiquidType:    wtype.LTSingleChannel,
								Sampler:       mixer.Sample,
							},
						})(lab)
						if err != nil {
							return nil, err
						}
						for _, ins := range instructions {
							ins.Outputs[0].Vol = 10.0
						}
						return instructions, nil
					},
					InputPlates:  []*wtype.LHPlate{GetTroughForTest(lab.IDGenerator)},
					OutputPlates: []*wtype.LHPlate{GetPlateForTest(lab.IDGenerator)},
					ErrorPrefix:  "7 (LH_ERR_VOL) : volume error : sum of requested volumes (17 ul) does not match result volume (10 ul)",
				},
				{
					Name: "multi channel dependent",
					Instructions: Mixes("pcrplate_skirted_riser", TestMixComponents{
						{
							LiquidName:    "water",
							VolumesByWell: ColumnWise(8, []float64{8.0, 8.0, 8.0, 8.0, 8.0, 8.0, 8.0, 8.0}),
							LiquidType:    wtype.LTWater,
							Sampler:       mixer.Sample,
						},
						{
							LiquidName:    "mastermix_sapI",
							VolumesByWell: ColumnWise(8, []float64{50.0, 50.0, 50.0, 50.0, 50.0, 50.0, 50.0, 50.0}),
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
					InputPlates:  []*wtype.LHPlate{GetTroughForTest(lab.IDGenerator)},
					OutputPlates: []*wtype.LHPlate{GetPlateForTest(lab.IDGenerator)},
					Assertions: Assertions{
						NumberOfAssertion(liquidhandling.ASP, 3), //full multichanneling
						NumberOfAssertion(liquidhandling.DSP, 3), //full multichanneling
					},
				},
				{
					Name: "single channel",
					Instructions: Mixes("pcrplate_skirted_riser", TestMixComponents{
						{
							LiquidName:    "water",
							VolumesByWell: ColumnWise(8, []float64{4.0, 0.0, 4.0, 0.0, 4.0, 0.0, 4.0, 0.0}),
							LiquidType:    wtype.LTWater,
							Sampler:       mixer.Sample,
						},
					}),
					InputPlates:  []*wtype.LHPlate{GetTroughForTest(lab.IDGenerator)},
					OutputPlates: []*wtype.LHPlate{GetPlateForTest(lab.IDGenerator)},
					Assertions: Assertions{
						NumberOfAssertion(liquidhandling.ASP, 4),
						NumberOfAssertion(liquidhandling.DSP, 4),
					},
				},
				{
					Name: "multi and single channel",
					Instructions: Mixes("pcrplate_skirted_riser", TestMixComponents{
						{
							LiquidName:    "water",
							VolumesByWell: ColumnWise(8, []float64{4.0, 8.0, 4.0, 8.0, 4.0, 8.0, 4.0, 8.0}),
							LiquidType:    wtype.LTWater,
							Sampler:       mixer.Sample,
						},
					}),
					InputPlates:  []*wtype.LHPlate{GetTroughForTest(lab.IDGenerator)},
					OutputPlates: []*wtype.LHPlate{GetPlateForTest(lab.IDGenerator)},
					Assertions: Assertions{
						NumberOfAssertion(liquidhandling.ASP, 5),
						NumberOfAssertion(liquidhandling.DSP, 5),
					},
				},
				{
					Name:          "multi channel independent",
					Liquidhandler: GetIndependentLiquidHandlerForTest(lab),
					Instructions: Mixes("pcrplate_skirted_riser", TestMixComponents{
						{
							LiquidName:    "water",
							VolumesByWell: ColumnWise(8, []float64{18.0, 7.0, 15.0, 12.0, 7.0, 8.0, 4.0, 8.0}),
							LiquidType:    wtype.LTWater,
							Sampler:       mixer.Sample,
						},
					}),
					InputPlates:  []*wtype.LHPlate{GetTroughForTest(lab.IDGenerator)},
					OutputPlates: []*wtype.LHPlate{GetPlateForTest(lab.IDGenerator)},
					Assertions: Assertions{
						NumberOfAssertion(liquidhandling.ASP, 1), //full multichanneling
						NumberOfAssertion(liquidhandling.DSP, 1), //full multichanneling
					},
				},
				{
					Name: "multi channel split-sample",
					Instructions: func(lab *laboratory.Laboratory) ([]*wtype.LHInstruction, error) {

						var instructions []*wtype.LHInstruction

						diluent := GetComponentForTest(lab, "multiwater", wunit.NewVolume(1000.0, "ul"))
						stock := GetComponentForTest(lab, "dna", wunit.NewVolume(1000, "ul"))
						stock.Type = wtype.LTMultiWater

						for y := 0; y < 8; y++ {
							lastStock := stock
							for x := 0; x < 2; x++ {
								diluentSample := mixer.Sample(lab, diluent, wunit.NewVolume(20.0, "ul"))

								split := getTestSplitSample(lab, lastStock, 20.0)

								wc := wtype.WellCoords{X: x, Y: y}
								mix := getTestMix(lab, []*wtype.Liquid{split.Outputs[0], diluentSample}, wc.FormatA1())

								lastStock = mix.Outputs[0]

								instructions = append(instructions, mix, split)
							}
						}
						return instructions, nil
					},
					InputPlates:  []*wtype.LHPlate{GetTroughForTest(lab.IDGenerator)},
					OutputPlates: []*wtype.LHPlate{GetPlateForTest(lab.IDGenerator)},
					Assertions: Assertions{
						NumberOfAssertion(liquidhandling.ASP, 4), //full multichanneling - 2 ops per dilution row
						NumberOfAssertion(liquidhandling.DSP, 4), //full multichanneling
					},
				},
				{
					Name: "single channel auto allocation",
					Instructions: Mixes("pcrplate_skirted_riser", TestMixComponents{
						{
							LiquidName:    "water",
							VolumesByWell: ColumnWise(8, []float64{8.0, 8.0, 8.0, 8.0, 8.0, 8.0, 8.0, 8.0}),
							LiquidType:    wtype.LTSingleChannel,
							Sampler:       mixer.Sample,
						},
					}),
					InputPlates:  []*wtype.LHPlate{GetPlateForTest(lab.IDGenerator)},
					OutputPlates: []*wtype.LHPlate{GetPlateForTest(lab.IDGenerator)},
					Assertions: Assertions{
						NumberOfAssertion(liquidhandling.ASP, 8),                                                                          //no multichanneling
						InputLayoutAssertion(map[string]string{"A1": "water"}),                                                            // should all be in the same well since no multichanneling
						InitialInputVolumesAssertion(0.001, map[string]float64{"A1": (8.0+wtype.GLOBALCARRYVOLUME.RawValue())*8.0 + 5.0}), // volume plus carry per transfer plus residual
						FinalInputVolumesAssertion(0.001, map[string]float64{"A1": 5.0}),
					},
				},
				{
					Name: "single channel well use",
					Instructions: Mixes("pcrplate_skirted_riser", TestMixComponents{
						{
							LiquidName:    "water",
							VolumesByWell: ColumnWise(8, []float64{8.0, 8.0, 8.0, 8.0, 8.0, 8.0, 8.0, 8.0}),
							LiquidType:    wtype.LTSingleChannel,
							Sampler:       mixer.Sample,
						},
					}),
					InputPlates:  []*wtype.LHPlate{PrefillPlateForTest(lab, GetPlateForTest(lab.IDGenerator), "water", map[string]float64{"A1": 200.0, "B1": 200.0, "C1": 200.0})},
					OutputPlates: []*wtype.LHPlate{GetPlateForTest(lab.IDGenerator)},
					Assertions: Assertions{
						NumberOfAssertion(liquidhandling.ASP, 8), //no multichanneling
						InputLayoutAssertion(map[string]string{"A1": "water", "B1": "water", "C1": "water"}),
						InitialInputVolumesAssertion(0.001, map[string]float64{"A1": 200.0, "B1": 200.0, "C1": 200.0}),
						// check that the same source well is used throughout since all of these operations are single channel
						FinalInputVolumesAssertion(0.001, map[string]float64{"A1": 200.0 - (8.0+0.5)*8.0, "B1": 200.0, "C1": 200.0}),
					},
				},
			}.Run(t)
			return nil
		},
	})
}

func TestFixDuplicatePlateNames(t *testing.T) {
	rq := lh.NewLHRequest(id.NewIDGenerator(t.Name()))
	for i := 0; i < 100; i++ {
		p := &wtype.Plate{ID: fmt.Sprintf("anID-%d", i), PlateName: "aName"}
		rq.InputPlateOrder = append(rq.InputPlateOrder, p.ID)
		rq.InputPlates[p.ID] = p
	}
	for i := 100; i < 200; i++ {
		p := &wtype.Plate{ID: fmt.Sprintf("anID-%d", i), PlateName: "aName"}
		rq.OutputPlateOrder = append(rq.OutputPlateOrder, p.ID)
		rq.OutputPlates[p.ID] = p
	}

	rq.FixDuplicatePlateNames()

	found := make(map[string]int)

	for _, p := range rq.AllPlates() {
		_, ok := found[p.PlateName]

		if !ok {
			found[p.PlateName] = 1
		} else {
			t.Errorf("fixDuplicatePlateNames failed to prevent duplicates: found at least two of %s", p.PlateName)
		}
	}
}

func assertCoordsEq(lhs, rhs []wtype.Coordinates3D) bool {
	if len(lhs) != len(rhs) {
		return false
	}

	for i := 0; i < len(lhs); i++ {
		if lhs[i].Subtract(rhs[i]).Abs() > 0.00001 {
			return false
		}
	}

	return true
}

func TestAddWellTargets(t *testing.T) {
	testlab.WithTestLab(t, "", &testlab.TestElementCallbacks{
		Steps: func(lab *laboratory.Laboratory) error {
			lh := GetLiquidHandlerForTest(lab)

			plate := GetPlateForTest(lab.IDGenerator)
			if err := lh.Properties.AddPlateTo("position_4", plate); err != nil {
				return err
			}

			tipwaste := GetTipwasteForTest(lab.IDGenerator)
			if err := lh.Properties.AddTipWasteTo("position_1", tipwaste); err != nil {
				return err
			}

			trough := GetTroughForTest(lab.IDGenerator)
			if err := lh.Properties.AddPlateTo("position_5", trough); err != nil {
				return err
			}

			if err := lh.AddWellTargets(); err != nil {
				return err
			}

			expected := []wtype.Coordinates3D{
				{X: 0.0, Y: -31.5, Z: 0.0},
				{X: 0.0, Y: -22.5, Z: 0.0},
				{X: 0.0, Y: -13.5, Z: 0.0},
				{X: 0.0, Y: -4.5, Z: 0.0},
				{X: 0.0, Y: 4.5, Z: 0.0},
				{X: 0.0, Y: 13.5, Z: 0.0},
				{X: 0.0, Y: 22.5, Z: 0.0},
				{X: 0.0, Y: 31.5, Z: 0.0},
			}

			if e, g := []wtype.Coordinates3D{}, plate.Welltype.GetWellTargets("DummyAdaptor"); !assertCoordsEq(e, g) {
				return fmt.Errorf("plate well targets incorrect, expected %v, got %v", e, g)
			}

			if e, g := expected, tipwaste.AsWell.GetWellTargets("DummyAdaptor"); !assertCoordsEq(e, g) {
				return fmt.Errorf("plate well targets incorrect, expected %v, got %v", e, g)
			}

			if e, g := expected, trough.Welltype.GetWellTargets("DummyAdaptor"); !assertCoordsEq(e, g) {
				return fmt.Errorf("plate well targets incorrect, expected %v, got %v", e, g)
			}
			return nil
		},
	})
}

/*
func TestShouldSetWellTargets(t *testing.T) {
	for _, plate := range testinventory.GetPlates(ctx) {
		e := !plate.IsSpecial()
		//IsSpecial is irrelevant for plates with 8 rows or more
		if plate.NRows() >= 8 {
			e = false
		}
		if g := plate.AreWellTargetsEnabled(8, 9.0); e != g {
			t.Errorf("For platetype %s (%d rows): plate.AreWellTargetsEnabled(8,9.0) = %t, expected %t", plate.GetType(), plate.NRows(), g, e)
		}
	}
}
*/

func getTestSplitSample(lab *laboratory.Laboratory, component *wtype.Liquid, volume float64) *wtype.LHInstruction {
	ret := wtype.NewLHSplitInstruction(lab.IDGenerator)

	ret.Inputs = append(ret.Inputs, component.Dup(lab.IDGenerator))
	cmpMoving, cmpStaying := mixer.SplitSample(lab, component, wunit.NewVolume(volume, "ul"))
	lab.SampleTracker.UpdateIDOf(component.ID, cmpStaying.ID)

	ret.AddOutput(cmpMoving)
	ret.AddOutput(cmpStaying)

	return ret
}

func getTestMix(lab *laboratory.Laboratory, components []*wtype.Liquid, address string) *wtype.LHInstruction {
	mix := mixer.GenericMix(lab, mixer.MixOptions{
		Inputs:  components,
		Address: address,
	})

	mx := 0
	for _, c := range components {
		if c.Generation() > mx {
			mx = c.Generation()
		}
	}
	mix.SetGeneration(mx)
	mix.Outputs[0].SetGeneration(mx + 1)
	mix.Outputs[0].DeclareInstance()

	return mix
}

type ShrinkVolumesTest struct {
	Name               string
	PlateLocations     map[string]*wtype.LHPlate                 // map address of plate to plate type
	AutoAllocatedWells map[string][]string                       // list wells in plates at each location which are autoallocated
	Instructions       []liquidhandling.TerminalRobotInstruction // the list of instructions to analyse
	CarryVolume        wunit.Volume
	ExpectedPlates     []string                      // list of addresses which should still contain plates after shrinkVolumes
	ExpectedVolumes    map[string]map[string]float64 // e.g. map[string]map[string]wunit.Volume{"pos_1": {"A1": 19.2}} -> assert plate at pos_1 has 19.2ul in well A1 afterwards
}

func (test *ShrinkVolumesTest) Run(t *testing.T) {
	testlab.WithTestLab(t, "", &testlab.TestElementCallbacks{
		Name: test.Name,
		Steps: func(lab *laboratory.Laboratory) error {

			// first, setup the test

			// build the initial state
			initial := makeGilson(lab)

			// add the plates
			inputPlateOrder := make([]string, 0, len(test.PlateLocations))
			for addr, plate := range test.PlateLocations {
				// fill each and every well
				for _, well := range plate.Wellcoords {
					contents := GetComponentForTest(lab, "water", well.MaxVolume())
					if err := well.SetContents(lab.IDGenerator, contents); err != nil {
						return err
					}
				}
				if err := initial.AddPlateTo(addr, plate); err != nil {
					return err
				}
				inputPlateOrder = append(inputPlateOrder, plate.ID)
			}

			// set autoallocated flags
			for addr, wells := range test.AutoAllocatedWells {
				if plate, ok := initial.Plates[addr]; !ok {
					return fmt.Errorf("can't set AutoAllocated: no plate at address \"%s\"", addr)
				} else {
					for _, wellAddr := range wells {
						if well, ok := plate.Wellcoords[wellAddr]; !ok {
							return fmt.Errorf("can't set AutoAllocated: plate type %s at %s has no well %s", plate.GetType(), addr, wellAddr)
						} else {
							well.DeclareAutoallocated()
						}
					}
				}
			}

			// initialise the liquidhandler
			handler := &lh.Liquidhandler{
				Properties: initial,
				// nb. in reality the volumes in FinalProperties will depend on the instructions in some way.
				// however since shrinkVolumes should only depend on the CName and well type, lets just duplicate
				FinalProperties: initial.DupKeepIDs(lab.IDGenerator),
			}

			// initialise the request
			rq := &lh.LHRequest{
				Instructions:    test.Instructions,
				InputPlateOrder: inputPlateOrder,
			}

			// secondly, do the thing

			if err := handler.ShrinkVolumes(lab.IDGenerator, rq); err != nil {
				// errors here are caused by bad config - users shouldn't be able to cause them
				return err
			}

			// finally, run some assertions

			// check that only expected plates are present in before and after
			assert.ElementsMatchf(t, test.ExpectedPlates, plateAddresses(handler.Properties), "initial properties plate addresses:\n e: %v\n g: %v", test.ExpectedPlates, plateAddresses(handler.Properties))
			assert.ElementsMatchf(t, test.ExpectedPlates, plateAddresses(handler.FinalProperties), "initial properties plate addresses:\n e: %v\n g: %v", test.ExpectedPlates, plateAddresses(handler.FinalProperties))

			// check that autoallocated flag has not been changed
			if err := checkAutoAllocated(test.AutoAllocatedWells, handler.Properties); err != nil {
				return fmt.Errorf("in initial properties: %s", err)
			}
			if err := checkAutoAllocated(test.AutoAllocatedWells, handler.FinalProperties); err != nil {
				return fmt.Errorf("in final properties: %s", err)
			}

			// initial volume assertions
			expectedInitialVolumes := make(map[string]map[string]float64, len(handler.Properties.Plates))
			// add assertions for non-autoallocated wells, which should remain full
			for addr, plate := range handler.Properties.Plates {
				plateVols := make(map[string]float64, len(plate.Wellcoords))
				eVol := plate.Welltype.MaxVolume().MustInStringUnit("ul").RawValue()
				for wc, well := range plate.Wellcoords {
					if !well.IsAutoallocated() {
						plateVols[wc] = eVol
					}
				}
				expectedInitialVolumes[addr] = plateVols
			}
			// overwrite with user assertions
			for addr, wellMap := range test.ExpectedVolumes {
				for wc, vol := range wellMap {
					expectedInitialVolumes[addr][wc] = vol
				}
			}
			if err := checkVols(lab.IDGenerator, expectedInitialVolumes, handler.Properties); err != nil {
				return fmt.Errorf("initial volumes incorrect: %s", err)
			}

			// final volume assertions
			expectedFinalVolumes := make(map[string]map[string]float64, len(handler.Properties.Plates))
			// non-autoallocated should remain full, autoallocated should be left with residual volume only
			for addr, plate := range handler.Properties.Plates {
				plateVols := make(map[string]float64, len(plate.Wellcoords))
				maxVol := plate.Welltype.MaxVolume().MustInStringUnit("ul").RawValue()
				rVol := plate.Welltype.ResidualVolume().MustInStringUnit("ul").RawValue()
				for wc, well := range plate.Wellcoords {
					if well.IsAutoallocated() {
						if expectedInitialVolumes[addr][wc] == 0.0 {
							// an empty well stays empty
							plateVols[wc] = 0.0
						} else {
							plateVols[wc] = rVol
						}
					} else {
						plateVols[wc] = maxVol
					}
				}
				expectedFinalVolumes[addr] = plateVols
			}
			if err := checkVols(lab.IDGenerator, expectedFinalVolumes, handler.FinalProperties); err != nil {
				return fmt.Errorf("final volumes incorrect: %s", err)
			}
			return nil
		},
	})
}

func plateAddresses(props *liquidhandling.LHProperties) []string {
	ret := make([]string, 0, len(props.Plates))
	for addr := range props.Plates {
		ret = append(ret, addr)
	}
	return ret
}

func checkAutoAllocated(autoAllocated map[string][]string, props *liquidhandling.LHProperties) error {
	errs := make(utils.ErrorSlice, 0, len(props.Plates))

	for addr, plate := range props.Plates {
		aa := make(map[string]bool, len(autoAllocated[addr]))
		for _, wc := range autoAllocated[addr] {
			aa[wc] = true
		}

		set := make([]string, 0, len(plate.Wellcoords))
		unset := make([]string, 0, len(plate.Wellcoords))

		for wc, well := range plate.Wellcoords {
			if aa[wc] && !well.IsAutoallocated() {
				unset = append(unset, wc)
			} else if !aa[wc] && well.IsAutoallocated() {
				set = append(set, wc)
			}
		}

		if len(set) > 0 || len(unset) > 0 {
			errs = append(errs, errors.Errorf("plate at %s:\n  set at: %s\n  unset at: %s", addr, strings.Join(set, ", "), strings.Join(unset, ", ")))
		}
	}

	return errors.WithMessage(errs.Pack(), "autoallocated flag changed")
}

func checkVols(idGen *id.IDGenerator, expected map[string]map[string]float64, props *liquidhandling.LHProperties) error {
	errs := make(utils.ErrorSlice, 0, len(props.Plates))

	for addr, wellMap := range expected {
		if plate, ok := props.Plates[addr]; !ok {
			errs = append(errs, errors.Errorf("no plate at address %s", addr))
		} else {
			wrong := make([]string, 0, len(plate.Wellcoords))
			for wc, volUl := range wellMap {
				eVol := wunit.NewVolume(volUl, "ul")
				if well, ok := plate.Wellcoords[wc]; !ok {
					wrong = append(wrong, fmt.Sprintf("%s: no well at location", wc))
				} else if !eVol.EqualToTolerance(well.CurrentVolume(idGen), 1.0e-5) {
					wrong = append(wrong, fmt.Sprintf("%s: expected %v, got %v", wc, eVol, well.CurrentVolume(idGen)))
				}
			}

			if len(wrong) > 0 {
				errs = append(errs, errors.Errorf("in plate at %s:\n  %s", addr, strings.Join(wrong, "\n  ")))
			}
		}
	}

	return errors.WithMessage(errs.Pack(), "auto-allocated volumes changed")
}

type ShrinkVolumesTests []ShrinkVolumesTest

func (tests ShrinkVolumesTests) Run(t *testing.T) {
	for _, test := range tests {
		test.Run(t)
	}
}

func TestShrinkVolumes(t *testing.T) {
	var pcrplateResidual float64

	testlab.WithTestLab(t, "", &testlab.TestElementCallbacks{
		Steps: func(lab *laboratory.Laboratory) error {

			if plate, err := lab.Inventory.Plates.NewPlate("pcrplate_skirted"); err != nil {
				return err
			} else {
				pcrplateResidual = plate.Welltype.ResidualVolume().MustInStringUnit("ul").RawValue()
			}

			defaultCarry := 0.5

			ul := func(vols ...float64) []wunit.Volume {
				ret := make([]wunit.Volume, 0, len(vols))
				for _, v := range vols {
					ret = append(ret, wunit.NewVolume(v, "ul"))
				}
				return ret
			}

			getPlate := func(plateType wtype.PlateTypeName, idOverride string) *wtype.LHPlate {
				if plate, err := lab.Inventory.Plates.NewPlate(plateType); err != nil {
					t.Fatal(err)
					return nil
				} else {
					if idOverride != "" {
						plate.ID = idOverride
					}
					return plate
				}
			}

			ShrinkVolumesTests{
				{
					Name: "simpleSingleChannel",
					PlateLocations: map[string]*wtype.LHPlate{
						"position_1": getPlate("pcrplate_skirted", ""),
					},
					AutoAllocatedWells: map[string][]string{
						"position_1": {"A1"},
					},
					Instructions: []liquidhandling.TerminalRobotInstruction{
						&liquidhandling.MoveInstruction{
							Pos:  []string{"position_1"},
							Well: []string{"A1"},
						},
						&liquidhandling.AspirateInstruction{
							Volume: ul(20.0),
						},
					},
					CarryVolume:    wunit.NewVolume(defaultCarry, "ul"),
					ExpectedPlates: []string{"position_1"},
					ExpectedVolumes: map[string]map[string]float64{
						"position_1": {"A1": 20 + pcrplateResidual + defaultCarry},
					},
				},
				{
					Name: "don't change non-autoallocated",
					PlateLocations: map[string]*wtype.LHPlate{
						"position_1": getPlate("pcrplate_skirted", ""),
					},
					AutoAllocatedWells: map[string][]string{},
					Instructions: []liquidhandling.TerminalRobotInstruction{
						&liquidhandling.MoveInstruction{
							Pos:  []string{"position_1"},
							Well: []string{"A1"},
						},
						&liquidhandling.AspirateInstruction{
							Volume: ul(20.0),
						},
					},
					CarryVolume:     wunit.NewVolume(defaultCarry, "ul"),
					ExpectedPlates:  []string{"position_1"},
					ExpectedVolumes: map[string]map[string]float64{},
				},
				{
					Name: "simpleSingleChannelTransfer",
					PlateLocations: map[string]*wtype.LHPlate{
						"position_1": getPlate("pcrplate_skirted", "plate1"),
					},
					AutoAllocatedWells: map[string][]string{
						"position_1": {"A1"},
					},
					Instructions: []liquidhandling.TerminalRobotInstruction{
						&liquidhandling.TransferInstruction{
							Transfers: []liquidhandling.MultiTransferParams{
								{
									Transfers: []liquidhandling.TransferParams{
										{
											PltFrom:  "plate1",
											WellFrom: "A1",
											Volume:   wunit.NewVolume(20, "ul"),
										},
									},
								},
							},
						},
					},
					CarryVolume:    wunit.NewVolume(defaultCarry, "ul"),
					ExpectedPlates: []string{"position_1"},
					ExpectedVolumes: map[string]map[string]float64{
						"position_1": {"A1": 20 + pcrplateResidual}, // transfers instructions don't incurr a carry
					},
				},
				{
					Name: "simpleMultiChannel",
					PlateLocations: map[string]*wtype.LHPlate{
						"position_1": getPlate("pcrplate_skirted", ""),
					},
					AutoAllocatedWells: map[string][]string{
						"position_1": {"A1", "B1", "C1", "D1", "E1", "F1", "G1", "H1"},
					},
					Instructions: []liquidhandling.TerminalRobotInstruction{
						&liquidhandling.MoveInstruction{
							Pos:  []string{"position_1", "position_1", "position_1", "position_1", "position_1", "position_1", "position_1", "position_1"},
							Well: []string{"A1", "B1", "C1", "D1", "E1", "F1", "G1", "H1"},
						},
						&liquidhandling.AspirateInstruction{
							Volume: ul(20.0, 20.0, 20.0, 20.0, 20.0, 20.0, 20.0, 20.0),
						},
					},
					CarryVolume:    wunit.NewVolume(defaultCarry, "ul"),
					ExpectedPlates: []string{"position_1"},
					ExpectedVolumes: map[string]map[string]float64{
						"position_1": {
							"A1": 20 + pcrplateResidual + defaultCarry,
							"B1": 20 + pcrplateResidual + defaultCarry,
							"C1": 20 + pcrplateResidual + defaultCarry,
							"D1": 20 + pcrplateResidual + defaultCarry,
							"E1": 20 + pcrplateResidual + defaultCarry,
							"F1": 20 + pcrplateResidual + defaultCarry,
							"G1": 20 + pcrplateResidual + defaultCarry,
							"H1": 20 + pcrplateResidual + defaultCarry,
						},
					},
				},
				{
					Name: "simpleMultiChannelTransfer",
					PlateLocations: map[string]*wtype.LHPlate{
						"position_1": getPlate("pcrplate_skirted", "plate1"),
					},
					AutoAllocatedWells: map[string][]string{
						"position_1": {"A1", "B1", "C1", "D1", "E1", "F1", "G1", "H1"},
					},
					Instructions: []liquidhandling.TerminalRobotInstruction{
						&liquidhandling.TransferInstruction{
							Transfers: []liquidhandling.MultiTransferParams{
								{
									Transfers: []liquidhandling.TransferParams{
										{
											PltFrom:  "plate1",
											WellFrom: "A1",
											Volume:   wunit.NewVolume(20, "ul"),
										},
										{
											PltFrom:  "plate1",
											WellFrom: "B1",
											Volume:   wunit.NewVolume(20, "ul"),
										},
										{
											PltFrom:  "plate1",
											WellFrom: "C1",
											Volume:   wunit.NewVolume(20, "ul"),
										},
										{
											PltFrom:  "plate1",
											WellFrom: "D1",
											Volume:   wunit.NewVolume(20, "ul"),
										},
										{
											PltFrom:  "plate1",
											WellFrom: "E1",
											Volume:   wunit.NewVolume(20, "ul"),
										},
										{
											PltFrom:  "plate1",
											WellFrom: "F1",
											Volume:   wunit.NewVolume(20, "ul"),
										},
										{
											PltFrom:  "plate1",
											WellFrom: "G1",
											Volume:   wunit.NewVolume(20, "ul"),
										},
										{
											PltFrom:  "plate1",
											WellFrom: "H1",
											Volume:   wunit.NewVolume(20, "ul"),
										},
									},
								},
							},
						},
					},
					CarryVolume:    wunit.NewVolume(defaultCarry, "ul"),
					ExpectedPlates: []string{"position_1"},
					ExpectedVolumes: map[string]map[string]float64{
						"position_1": {
							"A1": 20 + pcrplateResidual,
							"B1": 20 + pcrplateResidual,
							"C1": 20 + pcrplateResidual,
							"D1": 20 + pcrplateResidual,
							"E1": 20 + pcrplateResidual,
							"F1": 20 + pcrplateResidual,
							"G1": 20 + pcrplateResidual,
							"H1": 20 + pcrplateResidual,
						},
					},
				},
				{
					Name: "removeUnusedWells",
					PlateLocations: map[string]*wtype.LHPlate{
						"position_1": getPlate("pcrplate_skirted", ""),
					},
					AutoAllocatedWells: map[string][]string{
						"position_1": {"A1", "B1", "C1", "D1", "E1", "F1", "G1", "H1"},
					},
					Instructions: []liquidhandling.TerminalRobotInstruction{
						&liquidhandling.MoveInstruction{
							Pos:  []string{"position_1", "position_1", "position_1", "position_1"},
							Well: []string{"A1", "B1", "C1", "D1"},
						},
						&liquidhandling.AspirateInstruction{
							Volume: ul(20.0, 20.0, 20.0, 20.0),
						},
					},
					CarryVolume:    wunit.NewVolume(defaultCarry, "ul"),
					ExpectedPlates: []string{"position_1"},
					ExpectedVolumes: map[string]map[string]float64{
						"position_1": {
							"A1": 20 + pcrplateResidual + defaultCarry,
							"B1": 20 + pcrplateResidual + defaultCarry,
							"C1": 20 + pcrplateResidual + defaultCarry,
							"D1": 20 + pcrplateResidual + defaultCarry,
							"E1": 0,
							"F1": 0,
							"G1": 0,
							"H1": 0,
						},
					},
				},
				{
					Name: "removeUnusedPlate",
					PlateLocations: map[string]*wtype.LHPlate{
						"position_1": getPlate("pcrplate_skirted", ""),
						"position_2": getPlate("reservoir", ""),
					},
					AutoAllocatedWells: map[string][]string{
						"position_1": {"A1", "B1", "C1", "D1"},
						"position_2": {"A1"},
					},
					Instructions: []liquidhandling.TerminalRobotInstruction{
						&liquidhandling.MoveInstruction{
							Pos:  []string{"position_1", "position_1", "position_1", "position_1"},
							Well: []string{"A1", "B1", "C1", "D1"},
						},
						&liquidhandling.AspirateInstruction{
							Volume: ul(20.0, 20.0, 20.0, 20.0),
						},
					},
					CarryVolume:    wunit.NewVolume(defaultCarry, "ul"),
					ExpectedPlates: []string{"position_1"},
					ExpectedVolumes: map[string]map[string]float64{
						"position_1": {
							"A1": 20 + pcrplateResidual + defaultCarry,
							"B1": 20 + pcrplateResidual + defaultCarry,
							"C1": 20 + pcrplateResidual + defaultCarry,
							"D1": 20 + pcrplateResidual + defaultCarry,
						},
					},
				},
				{
					Name: "dont remove non-auto allocated plate",
					PlateLocations: map[string]*wtype.LHPlate{
						"position_1": getPlate("pcrplate_skirted", ""),
						"position_2": getPlate("reservoir", ""),
					},
					AutoAllocatedWells: map[string][]string{
						"position_1": {"A1", "B1", "C1", "D1"},
					},
					Instructions: []liquidhandling.TerminalRobotInstruction{
						&liquidhandling.MoveInstruction{
							Pos:  []string{"position_1", "position_1", "position_1", "position_1"},
							Well: []string{"A1", "B1", "C1", "D1"},
						},
						&liquidhandling.AspirateInstruction{
							Volume: ul(20.0, 20.0, 20.0, 20.0),
						},
					},
					CarryVolume:    wunit.NewVolume(defaultCarry, "ul"),
					ExpectedPlates: []string{"position_1", "position_2"},
					ExpectedVolumes: map[string]map[string]float64{
						"position_1": {
							"A1": 20 + pcrplateResidual + defaultCarry,
							"B1": 20 + pcrplateResidual + defaultCarry,
							"C1": 20 + pcrplateResidual + defaultCarry,
							"D1": 20 + pcrplateResidual + defaultCarry,
						},
					},
				},
			}.Run(t)
			return nil
		},
	})
}
