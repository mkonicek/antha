// /anthalib/wtype/locations_test.go: Part of the Antha language
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

package wtype

import (
	"errors"
	"testing"
)

func TestSameLocation(t *testing.T) {
	testResult := make([]bool, 3)
	testSuite := make([][]Location, 3)
	testLevel := make([]int, 3)

	testSuite[0] = make([]Location, 2)
	testSuite[0][0] = NewLocation("origin", 1, NewShape(BoxShape, "", 0, 0, 0))
	testSuite[0][1] = testSuite[0][0]
	testResult[0] = true
	testLevel[0] = 0

	testSuite[1] = make([]Location, 2)
	testSuite[1][0] = NewLocation("origin", 1, NewShape(BoxShape, "", 0, 0, 0))
	testSuite[1][1] = NewLocation("origin", 1, NewShape(BoxShape, "", 0, 0, 0))
	testResult[1] = false
	testLevel[1] = 0

	testSuite[2] = make([]Location, 2)
	testSuite[2][0] = NewLocation("origin", 1, NewShape(BoxShape, "", 0, 0, 0))
	testSuite[2][1] = testSuite[2][0].Positions()[0]
	testResult[2] = true
	testLevel[2] = 1

	for i := range testSuite {
		if i > len(testResult) {
			t.Fatal(errors.New("Not enough results defined in test"))
		}
		if i > len(testLevel) {
			t.Fatal(errors.New("Not enough levels defined in test"))
		}
		res := SameLocation(testSuite[i][0], testSuite[i][1], testLevel[i])
		if testResult[i] != res {
			t.Fatalf("On location %d. Expecting %v. Got %v.", i, testResult[i], res)
		}
	}
}
