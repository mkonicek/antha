// antha/AnthaStandardLibrary/Packages/enzymes/Find.go: Part of the Antha language
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

// Package search is a utility package providing functions useful for:
// Searching for a target entry in a slice;
// Removing duplicate values from a slice;
// Comparing the Name of two entries of any type with a Name() method returning a string.
// FindAll instances of a target string within a template string.
package search

import (
	"errors"
	"fmt"
	"testing"

	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/laboratory"
	"github.com/antha-lang/antha/laboratory/testlab"
)

func TestNextFreeWell(t *testing.T) {
	type nextWellTest struct {
		avoidWells     []string
		preferredWells []string
		plateType      *wtype.Plate
		byRow          bool
		expectedResult string
		expectedError  error
		options        []Option
	}

	testlab.WithTestLab(t, "", &testlab.TestElementCallbacks{
		Steps: func(lab *laboratory.Laboratory) error {
			// create a test plate
			falconAgarPlate, err := lab.Inventory.Plates.NewPlate("falcon6wellAgar")
			if err != nil {
				return err
			}

			falconAgarPlate.SetName("Agar plate")

			falconAgarPlateWithSomethingIn, err := lab.Inventory.Plates.NewPlate("falcon6wellAgar")
			if err != nil {
				return err
			}

			// second test plate which we'll add a sample to.
			falconAgarPlateWithSomethingIn.SetName("Agar plate with sample")

			component, err := lab.Inventory.Components.NewComponent("water")
			if err != nil {
				return err
			}

			component.SetName("test_sample")

			component.SetVolume(wunit.NewVolume(100.0, "ul"))

			// this will add the component to the plate
			_, err = falconAgarPlateWithSomethingIn.AddComponent(lab.IDGenerator, component, false)
			if err != nil {
				return err
			}

			// create a 384 well test plate
			griener384, err := lab.Inventory.Plates.NewPlate("greiner384_riser18")
			if err != nil {
				return err
			}

			griener384.SetName("384 well plate")

			var nextwellTests = []nextWellTest{
				{
					avoidWells:     []string{},
					preferredWells: []string{},
					plateType:      falconAgarPlate,
					byRow:          false,
					expectedResult: "A1",
					expectedError:  nil,
				},
				{
					avoidWells:     []string{},
					preferredWells: []string{},
					plateType:      falconAgarPlateWithSomethingIn,
					byRow:          false,
					expectedResult: "B1",
					expectedError:  nil,
				},
				{
					avoidWells:     []string{"A1"},
					preferredWells: []string{},
					plateType:      falconAgarPlate,
					byRow:          false,
					expectedResult: "B1",
					expectedError:  nil,
				},
				{
					avoidWells:     []string{"A1"},
					preferredWells: []string{},
					plateType:      falconAgarPlate,
					byRow:          true,
					expectedResult: "A2",
					expectedError:  nil,
				},
				{
					avoidWells:     []string{"A1"},
					preferredWells: []string{"A3"},
					plateType:      falconAgarPlate,
					byRow:          false,
					expectedResult: "A3",
					expectedError:  nil,
				},
				{
					avoidWells:     []string{"A1"},
					preferredWells: []string{"A1"},
					plateType:      falconAgarPlate,
					byRow:          false,
					expectedResult: "B1",
					expectedError:  nil,
				},
				{
					avoidWells:     []string{"A1"},
					preferredWells: []string{"A13"},
					plateType:      falconAgarPlate,
					byRow:          false,
					expectedResult: "",
					expectedError:  errors.New("well (A13) specified is out of range of available wells for plate type falcon6wellAgar"),
				},
				{
					avoidWells:     []string{"A1", "B1", "A2", "B2", "A3", "B3"},
					preferredWells: []string{"A1"},
					plateType:      falconAgarPlate,
					byRow:          false,
					expectedResult: "",
					expectedError:  errors.New("no empty wells on plate Agar plate of type falcon6wellAgar"),
				},
				{
					avoidWells:     []string{"A1"},
					preferredWells: []string{},
					plateType:      griener384,
					byRow:          false,
					expectedResult: "C1",
					expectedError:  nil,
					options:        []Option{SkipAlternateWells},
				},
				{
					avoidWells:     []string{},
					preferredWells: []string{},
					plateType:      nil,
					byRow:          false,
					expectedResult: "",
					expectedError:  errors.New("no plate specified as argument to NextFreeWell function"),
				},
			}

			for idx, test := range nextwellTests {
				testlab.WithTestLab(t, "", &testlab.TestElementCallbacks{
					Name: fmt.Sprint(idx),
					Steps: func(lab *laboratory.Laboratory) error {
						well, err := NextFreeWell(lab, test.plateType, test.avoidWells, test.preferredWells, test.byRow, test.options...)

						if well != test.expectedResult {
							return fmt.Errorf("For: %v (avoid: %v, preferred: %v, byRow: %v)\n\texpected: %v\n\tgot: %v",
								test.plateType, test.avoidWells, test.preferredWells, test.byRow, test.expectedResult, well)
						}

						if err != test.expectedError {
							if test.expectedError != nil && err != nil {
								if test.expectedError.Error() != err.Error() {
									return fmt.Errorf("For: %v (avoid: %v, preferred: %v, byRow: %v)\n\texpected: %v\n\tgot: %v",
										test.plateType, test.avoidWells, test.preferredWells, test.byRow, test.expectedError, err)
								}
							} else {
								return fmt.Errorf("For: %v (avoid: %v, preferred: %v, byRow: %v)\n\texpected: %v\n\tgot: %v",
									test.plateType, test.avoidWells, test.preferredWells, test.byRow, test.expectedError, err)
							}
						}
						return nil
					},
				})
			}
			return nil
		},
	})
}
