package wunit

import (
	"testing"
)

// SortTest this is used to test both SortConcentrations and SortVolumes with only
// a little code duplication
type SortTest struct {
	Name  string
	Input []interface{}
	Order []int //index of which should be where
	Max   int   //index of which should be returned as maximum
	Min   int   //index of which should be returned as minimum
	Error bool  //true if error is expected, in which case Order, Max, and Min are ignored
}

func (self *SortTest) RunConcentrations(t *testing.T) {
	t.Run(self.Name, func(t *testing.T) {
		//assumption: sort doesn't modify the original, test will fail if it does
		input := make([]Concentration, 0, len(self.Input))
		for _, i := range self.Input {
			input = append(input, i.(Concentration))
		}
		sorted, err := SortConcentrations(input)
		output := make([]interface{}, 0, len(sorted))
		for _, i := range sorted {
			output = append(output, i)
		}

		min, _ := MinConcentration(input) // nolint - error is the same as above
		max, _ := MaxConcentration(input) // nolint - error is the same as above

		self.Validate(t, output, min, max, err)
	})
}

func (self *SortTest) RunVolumes(t *testing.T) {
	t.Run(self.Name, func(t *testing.T) {
		//assumption: sort doesn't modify the original, test will fail if it does
		input := make([]Volume, 0, len(self.Input))
		for _, i := range self.Input {
			input = append(input, i.(Volume))
		}
		sorted, err := SortVolumes(input)
		output := make([]interface{}, 0, len(sorted))
		for _, i := range sorted {
			output = append(output, i)
		}

		min, _ := MinVolume(input) // nolint - error is the same as above
		max, _ := MaxVolume(input) // nolint - error is the same as above

		self.Validate(t, output, min, max, err)
	})
}

func (self *SortTest) Validate(t *testing.T, sorted []interface{}, min, max interface{}, err error) {
	if hasErr := err != nil; self.Error != hasErr {
		t.Errorf("expected error: %t, got error: %v", self.Error, err)
		return
	}
	if !self.Error {
		expected := make([]interface{}, 0, len(self.Order))
		for i := range self.Order {
			expected = append(expected, self.Input[self.Order[i]])
		}
		if len(expected) != len(sorted) {
			t.Fatalf("expected(%d) and sorted(%d) not the same length", len(expected), len(sorted))
		}
		mismatched := []int{}
		for i := range self.Order {
			if sorted[i] != expected[i] {
				mismatched = append(mismatched, i)
			}
		}
		if len(mismatched) > 0 {
			t.Errorf("sorted list doesn't match expected:\n\te: %v\n\tg: %v", expected, sorted)
		}

		if min != self.Input[self.Order[0]] {
			t.Errorf("got min: %v, expected: %v", self.Input[self.Min], min)
		}

		if max != self.Input[self.Order[len(self.Order)-1]] {
			t.Errorf("got max: %v, expected: %v", self.Input[self.Max], max)
		}
	}
}

type SortTests []*SortTest

func (self SortTests) RunConcentrations(t *testing.T) {
	for _, test := range self {
		test.RunConcentrations(t)
	}
}

func (self SortTests) RunVolumes(t *testing.T) {
	for _, test := range self {
		test.RunVolumes(t)
	}
}
func TestSortConcentrations(t *testing.T) {
	SortTests{
		{
			Name: "simple test",
			Input: []interface{}{
				NewConcentration(3, "g/l"),
				NewConcentration(2, "g/l"),
				NewConcentration(1, "g/l"),
			},
			Order: []int{2, 1, 0},
		},
		{
			Name: "already sorted",
			Input: []interface{}{
				NewConcentration(1, "g/l"),
				NewConcentration(2, "g/l"),
				NewConcentration(3, "g/l"),
			},
			Order: []int{0, 1, 2},
		},
		{
			Name: "values equal",
			Input: []interface{}{
				NewConcentration(3, "g/l"),
				NewConcentration(2, "g/l"),
				NewConcentration(1, "g/l"),
				NewConcentration(2, "g/l"),
			},
			Order: []int{2, 1, 3, 0}, //swapping 1 and 3 permitted
		},
		{
			Name: "different units",
			Input: []interface{}{
				NewConcentration(3, "ug/ul"),
				NewConcentration(2, "mg/ml"),
				NewConcentration(1, "g/l"),
			},
			Order: []int{2, 1, 0},
		},
		{
			Name: "different units affecting order",
			Input: []interface{}{
				NewConcentration(3, "ng/ul"),
				NewConcentration(2, "mg/ml"),
				NewConcentration(1, "g/l"),
			},
			Order: []int{0, 2, 1},
		},
		{
			Name: "invalid units",
			Input: []interface{}{
				NewConcentration(3, "ng/ul"),
				NewConcentration(2, "mg/ml"),
				NewConcentration(1, "X"),
			},
			Error: true,
		},
		{
			Name:  "zero length raises error",
			Input: nil,
			Error: true,
		},
	}.RunConcentrations(t)
}

func TestSortVolumes(t *testing.T) {
	SortTests{
		{
			Name: "simple test",
			Input: []interface{}{
				NewVolume(3, "ul"),
				NewVolume(2, "ul"),
				NewVolume(1, "ul"),
			},
			Order: []int{2, 1, 0},
		},
		{
			Name: "already sorted",
			Input: []interface{}{
				NewVolume(1, "ul"),
				NewVolume(2, "ul"),
				NewVolume(3, "ul"),
			},
			Order: []int{0, 1, 2},
		},
		{
			Name: "values equal",
			Input: []interface{}{
				NewVolume(3, "ul"),
				NewVolume(2, "ul"),
				NewVolume(1, "ul"),
				NewVolume(2, "ul"),
			},
			Order: []int{2, 1, 3, 0}, //swapping 1 and 3 permitted
		},
		{
			Name: "different units",
			Input: []interface{}{
				NewVolume(3, "nl"),
				NewVolume(2, "ml"),
				NewVolume(1, "ul"),
			},
			Order: []int{0, 2, 1},
		},
		{
			Name:  "zero length raises error",
			Input: nil,
			Error: true,
		},
	}.RunVolumes(t)
}
