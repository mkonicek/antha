// solutions
package solutions

import (
	"reflect"

	"testing"

	"github.com/antha-lang/antha/antha/anthalib/wunit"
)

type mixComponentlistTest struct {
	name      string
	sample1   ComponentListSample
	sample2   ComponentListSample
	mixedList ComponentList
}

var tests []mixComponentlistTest = []mixComponentlistTest{
	mixComponentlistTest{
		sample1: ComponentListSample{
			ComponentList: ComponentList{
				Components: map[string]wunit.Concentration{
					"water": wunit.NewConcentration(1, "g/L"),
				},
			},
			Volume: wunit.NewVolume(1, "ul"),
		},
		sample2: ComponentListSample{
			ComponentList: ComponentList{
				Components: map[string]wunit.Concentration{
					"dna":  wunit.NewConcentration(1, "g/L"),
					"dna2": wunit.NewConcentration(2, "X"),
				},
			},
			Volume: wunit.NewVolume(1, "ul"),
		},
		mixedList: ComponentList{
			Components: map[string]wunit.Concentration{
				"water": wunit.NewConcentration(0.5, "g/L"),
				"dna":   wunit.NewConcentration(0.5, "g/L"),
				"dna2":  wunit.NewConcentration(1, "X"),
			},
		},
	},
	mixComponentlistTest{
		sample1: ComponentListSample{
			ComponentList: ComponentList{
				Components: map[string]wunit.Concentration{
					"water": wunit.NewConcentration(1, "g/L"),
					"dna":   wunit.NewConcentration(1, "g/L"),
				},
			},
			Volume: wunit.NewVolume(1, "ul"),
		},
		sample2: ComponentListSample{
			ComponentList: ComponentList{
				Components: map[string]wunit.Concentration{
					"dna": wunit.NewConcentration(1, "g/L"),
				},
			},
			Volume: wunit.NewVolume(1, "ul"),
		},
		mixedList: ComponentList{
			Components: map[string]wunit.Concentration{
				"water": wunit.NewConcentration(0.5, "g/L"),
				"dna":   wunit.NewConcentration(1, "g/L"),
			},
		},
	},
	mixComponentlistTest{
		sample1: ComponentListSample{
			ComponentList: ComponentList{
				Components: map[string]wunit.Concentration{
					"water": wunit.NewConcentration(1, "g/l"),
					"dna":   wunit.NewConcentration(1, "g/l"),
				},
			},
			Volume: wunit.NewVolume(1, "ul"),
		},
		sample2: ComponentListSample{
			ComponentList: ComponentList{
				Components: map[string]wunit.Concentration{
					"dna": wunit.NewConcentration(1000, "mg/l"),
				},
			},
			Volume: wunit.NewVolume(1, "ul"),
		},
		mixedList: ComponentList{
			Components: map[string]wunit.Concentration{
				"water": wunit.NewConcentration(0.5, "g/l"),
				"dna":   wunit.NewConcentration(1, "g/l"),
			},
		},
	},
	mixComponentlistTest{
		sample1: ComponentListSample{
			ComponentList: ComponentList{
				Components: map[string]wunit.Concentration{
					"water":    wunit.NewConcentration(1, "g/L"),
					"glycerol": wunit.NewConcentration(1, "M"),
				},
			},
			Volume: wunit.NewVolume(1, "ul"),
		},
		sample2: ComponentListSample{
			ComponentList: ComponentList{
				Components: map[string]wunit.Concentration{
					"glycerol": wunit.NewConcentration(1, "g/L"),
				},
			},
			Volume: wunit.NewVolume(1, "ul"),
		},
		mixedList: ComponentList{
			Components: map[string]wunit.Concentration{
				"water":    wunit.NewConcentration(0.5, "g/L"),
				"glycerol": wunit.NewConcentration(46.5, "g/L"),
			},
		},
	},
}

// Align two dna sequences based on a specified scoring matrix
func TestSimulateMix(t *testing.T) {
	for _, test := range tests {
		mixed, err := mixComponentLists(test.sample1, test.sample2)

		if err != nil {
			t.Error(
				"For", test.name, "\n",
				"got error:", err.Error(), "\n",
			)
		}
		if !reflect.DeepEqual(mixed, test.mixedList) {
			t.Error(
				"For", test.name, "\n",
				"expected:", "\n",
				test.mixedList,
				"got:", "\n",
				mixed,
			)
		}
	}
}
