package wunit

import (
	"reflect"
	"testing"
)

// SortTest this is used to test both SortConcentrations and SortVolumes with only
// a little code duplication
type SortTest struct {
	Name        string
	Order       []int //index of which should be where
	Max         int   //index of which should be returned as maximum
	Min         int   //index of which should be returned as minimum
	ShouldError bool  //true if error is expected, in which case Order, Max, and Min are ignored
}

func (self *SortTest) unexpectedError(err error) bool {
	return (err != nil) != self.ShouldError
}

type ConcentrationSortTest struct {
	SortTest
	Input []Concentration
}

func (self *ConcentrationSortTest) Run(t *testing.T) {
	t.Run(self.Name, func(t *testing.T) {
		//assumption: sort doesn't modify the original, test will fail if it does
		sorted, err := SortConcentrations(self.Input)
		if self.unexpectedError(err) {
			t.Errorf("expecting error %t: got error %v", self.ShouldError, err)
		}

		if !self.ShouldError {
			expected := make([]Concentration, 0, len(self.Order))
			for _, index := range self.Order {
				expected = append(expected, self.Input[index])
			}

			if !reflect.DeepEqual(expected, sorted) {
				t.Errorf("sorting failed: expected %v: got %v", expected, sorted)
			}

			min, _ := MinConcentration(self.Input) // nolint - error is the same as above
			if e := self.Input[self.Order[0]]; e != min {
				t.Errorf("finding minimum failed: expected %v: got %v", e, min)
			}

			max, _ := MaxConcentration(self.Input) // nolint - error is the same as above
			if e := self.Input[self.Order[len(self.Order)-1]]; e != max {
				t.Errorf("finding maximum failed: expected %v: got %v", e, min)
			}
		}
	})
}

type ConcentrationSortTests []*ConcentrationSortTest

func (self ConcentrationSortTests) Run(t *testing.T) {
	for _, test := range self {
		test.Run(t)
	}
}

type VolumeSortTest struct {
	SortTest
	Input []Volume
}

func (self *VolumeSortTest) Run(t *testing.T) {
	t.Run(self.Name, func(t *testing.T) {
		//assumption: sort doesn't modify the original, test will fail if it does
		sorted, err := SortVolumes(self.Input)
		if self.unexpectedError(err) {
			t.Errorf("expecting error %t: got error %v", self.ShouldError, err)
		}

		if !self.ShouldError {
			expected := make([]Volume, 0, len(self.Order))
			for _, index := range self.Order {
				expected = append(expected, self.Input[index])
			}

			if !reflect.DeepEqual(expected, sorted) {
				t.Errorf("sorting failed: expected %v: got %v", expected, sorted)
			}

			min, _ := MinVolume(self.Input) // nolint - error is the same as above
			if e := self.Input[self.Order[0]]; e != min {
				t.Errorf("finding minimum failed: expected %v: got %v", e, min)
			}

			max, _ := MaxVolume(self.Input) // nolint - error is the same as above
			if e := self.Input[self.Order[len(self.Order)-1]]; e != max {
				t.Errorf("finding maximum failed: expected %v: got %v", e, min)
			}
		}
	})
}

type VolumeSortTests []*VolumeSortTest

func (self VolumeSortTests) Run(t *testing.T) {
	for _, test := range self {
		test.Run(t)
	}
}

func TestSortConcentrations(t *testing.T) {
	ConcentrationSortTests{
		{
			SortTest: SortTest{
				Name:  "simple test",
				Order: []int{2, 1, 0},
			},
			Input: []Concentration{
				NewConcentration(3, "g/l"),
				NewConcentration(2, "g/l"),
				NewConcentration(1, "g/l"),
			},
		},
		{
			SortTest: SortTest{
				Name:  "already sorted",
				Order: []int{0, 1, 2},
			},
			Input: []Concentration{
				NewConcentration(1, "g/l"),
				NewConcentration(2, "g/l"),
				NewConcentration(3, "g/l"),
			},
		},
		{
			SortTest: SortTest{
				Name:  "values equal",
				Order: []int{2, 1, 3, 0}, //swapping 1 and 3 permitted
			},
			Input: []Concentration{
				NewConcentration(3, "g/l"),
				NewConcentration(2, "g/l"),
				NewConcentration(1, "g/l"),
				NewConcentration(2, "g/l"),
			},
		},
		{
			SortTest: SortTest{
				Name:  "different units",
				Order: []int{2, 1, 0},
			},
			Input: []Concentration{
				NewConcentration(3, "ug/ul"),
				NewConcentration(2, "mg/ml"),
				NewConcentration(1, "g/l"),
			},
		},
		{
			SortTest: SortTest{
				Name:  "different units affecting order",
				Order: []int{0, 2, 1},
			},
			Input: []Concentration{
				NewConcentration(3, "ng/ul"),
				NewConcentration(2, "mg/ml"),
				NewConcentration(1, "g/l"),
			},
		},
		{
			SortTest: SortTest{
				Name:        "invalid units",
				ShouldError: true,
			},
			Input: []Concentration{
				NewConcentration(3, "ng/ul"),
				NewConcentration(2, "mg/ml"),
				NewConcentration(1, "X"),
			},
		},
		{
			SortTest: SortTest{
				Name:        "zero length raises error",
				ShouldError: true,
			},
		},
	}.Run(t)
}

func TestSortVolumes(t *testing.T) {
	VolumeSortTests{
		{
			SortTest: SortTest{
				Name:  "simple test",
				Order: []int{2, 1, 0},
			},
			Input: []Volume{
				NewVolume(3, "ul"),
				NewVolume(2, "ul"),
				NewVolume(1, "ul"),
			},
		},
		{
			SortTest: SortTest{
				Name:  "already sorted",
				Order: []int{0, 1, 2},
			},
			Input: []Volume{
				NewVolume(1, "ul"),
				NewVolume(2, "ul"),
				NewVolume(3, "ul"),
			},
		},
		{
			SortTest: SortTest{
				Name:  "values equal",
				Order: []int{2, 1, 3, 0}, //swapping 1 and 3 permitted
			},
			Input: []Volume{
				NewVolume(3, "ul"),
				NewVolume(2, "ul"),
				NewVolume(1, "ul"),
				NewVolume(2, "ul"),
			},
		},
		{
			SortTest: SortTest{
				Name:  "different units",
				Order: []int{0, 2, 1},
			},
			Input: []Volume{
				NewVolume(3, "nl"),
				NewVolume(2, "ml"),
				NewVolume(1, "ul"),
			},
		},
		{
			SortTest: SortTest{
				Name:        "zero length raises error",
				ShouldError: true,
			},
		},
	}.Run(t)
}
