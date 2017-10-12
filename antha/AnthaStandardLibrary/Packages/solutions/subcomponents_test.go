// solutions
package solutions

import (
	"math"
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

type serialComponentlistTest struct {
	name      string
	sample1   ComponentListSample
	sample2   ComponentListSample
	sample3   ComponentListSample
	mixedList ComponentList
}

var serialTests []serialComponentlistTest = []serialComponentlistTest{
	serialComponentlistTest{
		sample1: ComponentListSample{
			ComponentList: ComponentList{
				Components: map[string]wunit.Concentration{
					"water": wunit.NewConcentration(0, "g/L"),
				},
			},
			Volume: wunit.NewVolume(8, "ul"),
		},
		sample2: ComponentListSample{
			ComponentList: ComponentList{
				Components: map[string]wunit.Concentration{
					"dna": wunit.NewConcentration(1, "g/L"),
				},
			},
			Volume: wunit.NewVolume(1, "ul"),
		},
		sample3: ComponentListSample{
			ComponentList: ComponentList{
				Components: map[string]wunit.Concentration{
					"dna2": wunit.NewConcentration(1, "g/L"),
				},
			},
			Volume: wunit.NewVolume(1, "ul"),
		},
		mixedList: ComponentList{
			Components: map[string]wunit.Concentration{
				"water": wunit.NewConcentration(0, "g/L"),
				"dna":   wunit.NewConcentration(0.1, "g/L"),
				"dna2":  wunit.NewConcentration(0.1, "g/L"),
			},
		},
	},
}

func equal(list1, list2 ComponentList) bool {
	for key, value1 := range list1.Components {
		if value2, found := list2.Components[key]; found {
			if math.Abs(value1.SIValue()-value2.SIValue()) > 0.0001 {
				return false
			}
		} else {
			return false
		}
	}
	return true
}

func TestSimulateMix(t *testing.T) {
	for _, test := range tests {
		mixed, err := mixComponentLists(test.sample1, test.sample2)

		if err != nil {
			t.Error(
				"For", test.name, "\n",
				"got error:", err.Error(), "\n",
			)
		}
		if !equal(mixed, test.mixedList) {
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

func TestSerialMix(t *testing.T) {
	for _, test := range serialTests {
		intermediate, err := mixComponentLists(test.sample1, test.sample2)

		if err != nil {
			t.Error(
				"For", test.name, "\n",
				"got error:", err.Error(), "\n",
			)
		}

		intermediateSample := ComponentListSample{
			ComponentList: intermediate,
			Volume:        wunit.AddVolumes(test.sample1.Volume, test.sample2.Volume),
		}

		mixed, err := mixComponentLists(intermediateSample, test.sample3)

		if !equal(mixed, test.mixedList) {
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
