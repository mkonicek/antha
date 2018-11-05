package wtype

import (
	"strings"
	"testing"
)

type ParseTest struct {
	Value     string
	expectedX int
	expectedY int
}

func (self *ParseTest) Run(t *testing.T) {
	wc := MakeWellCoords(self.Value)
	if wc.X != self.expectedX || wc.Y != self.expectedY {
		t.Errorf("parsing \"%s\": expected {X: %d, Y:%d}, got {X: %d, X:%d}", self.Value, self.expectedX, self.expectedY, wc.X, wc.Y)
	}
}

func TestMakeWellCoords(t *testing.T) {

	tests := []*ParseTest{
		{"A1", 0, 0},
		{"AA1", 0, 26},
		{"1AA", 0, 26},
		{"AAA1", 0, 702},
		{"1AAA", 0, 702},
		{"X1Y1", 0, 0},
		{"not a well", -1, -1},
		{"B2", 1, 1},
	}

	for _, test := range tests {
		test.Run(t)
	}
}

func TestFormatWellCoords(t *testing.T) {
	wc := MakeWellCoordsA1("A1")

	if wc.FormatA1() != "A1" {
		t.Fatalf("Well coords A1 expected formatA1 to return A1, instead got %s", wc.FormatA1())
	}

	if wc.Format1A() != "1A" {
		t.Fatalf("Well coords A1 expected format1A to return 1A, instead got %s", wc.FormatA1())
	}

}

func TestWellCoordsComparison(t *testing.T) {
	s := []string{"C1", "A2", "HH1"}

	c := [][]int{{0, -1, -1}, {1, 0, 1}, {1, -1, 0}}
	r := [][]int{{0, 1, -1}, {-1, 0, -1}, {1, 1, 0}}

	for i := range s {
		for j := range s {
			cmpCol := CompareStringWellCoordsCol(s[i], s[j])
			cmpRow := CompareStringWellCoordsRow(s[i], s[j])

			expCol := c[i][j]
			expRow := r[i][j]

			if cmpCol != expCol {
				t.Fatalf("Compare WC Column Error: %s vs %s expected %d got %d", s[i], s[j], expCol, cmpCol)
			}
			if cmpRow != expRow {
				t.Fatalf("Compare WC Row Error: %s vs %s expected %d got %d", s[i], s[j], expRow, cmpRow)
			}

		}
	}
}

type HWellCoordsTest struct {
	Input    string
	Expected string
}

func (self *HWellCoordsTest) Run(t *testing.T) {
	wc := MakeWellCoordsArray(strings.Split(self.Input, ","))

	if g, e := HumanizeWellCoords(wc), self.Expected; g != e {
		t.Errorf("humanize [%s]: got %s, expected %s", self.Input, g, e)
	}
}

func TestHumanizeWellCoords(t *testing.T) {
	tests := []*HWellCoordsTest{
		{"A1,B1,C1,D1,E1,F1,G1,H1", "A1-H1"},
		{"A1,B1,C1,E1,F1,G1,H1", "A1-C1,E1-H1"},
		{"A1,C1,E1,G1", "A1,C1,E1,G1"},
		{"A1,A2,A3,A4,A5,A6,A7,A8,A9,A10,A11,A12,", "A1-A12"},
		{"A1,A2,A3,A4,A5,A7,A8,A9,A10,A11,A12,", "A1-A5,A7-A12"},
		{"A1,A2,A3,A4,A5,B5,C5,D5,E5,F5,G5,H5", "A1-A5,B5-H5"},
		{"H1,G1,F1,E1,D1,C1,B1,A1", "H1,G1,F1,E1,D1,C1,B1,A1"},
		{"A1", "A1"},
		{"A1,B2", "A1,B2"},
		{"A1,-,-,-,-,-,-", "A1"},
		{"-,-,-,-", ""},
	}

	for _, test := range tests {
		test.Run(t)
	}

}

func TestWellNumber(t *testing.T) {
	type wellNumberTest struct {
		Well               WellCoords
		WellsX             int
		WellsY             int
		ByRow              bool
		ExpectedWellNumber int
	}
	var tests = []wellNumberTest{
		{
			Well:               MakeWellCoordsA1("A1"),
			WellsX:             8,
			WellsY:             12,
			ByRow:              false,
			ExpectedWellNumber: 0,
		},
		{
			Well:               MakeWellCoordsA1("A1"),
			WellsX:             8,
			WellsY:             12,
			ByRow:              true,
			ExpectedWellNumber: 0,
		},
		{
			Well:               MakeWellCoordsA1("A2"),
			WellsX:             8,
			WellsY:             12,
			ByRow:              true,
			ExpectedWellNumber: 1,
		},
		{
			Well:               MakeWellCoordsA1("A2"),
			WellsX:             8,
			WellsY:             12,
			ByRow:              false,
			ExpectedWellNumber: 8,
		},
		{
			Well:               MakeWellCoordsA1("B1"),
			WellsX:             8,
			WellsY:             12,
			ByRow:              true,
			ExpectedWellNumber: 12,
		},
		{
			Well:               MakeWellCoordsA1("B1"),
			WellsX:             8,
			WellsY:             12,
			ByRow:              false,
			ExpectedWellNumber: 1,
		},
		{
			Well:               MakeWellCoordsA1("H12"),
			WellsX:             8,
			WellsY:             12,
			ByRow:              false,
			ExpectedWellNumber: 95,
		},
		{
			Well:               MakeWellCoordsA1("H12"),
			WellsX:             8,
			WellsY:             12,
			ByRow:              true,
			ExpectedWellNumber: 95,
		},
	}

	for _, test := range tests {
		num := test.Well.wellNumber(test.WellsX, test.WellsY, test.ByRow)

		if num != test.ExpectedWellNumber {
			t.Error(
				"expected: ", test.ExpectedWellNumber, "\n",
				"got: ", num, "\n",
			)
		}
	}

}
