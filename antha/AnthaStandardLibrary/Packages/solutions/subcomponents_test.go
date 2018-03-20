// solutions
package solutions

import (
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
	{
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
	{
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
	{
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
	{
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
	{
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
	{
		sample1: ComponentListSample{
			ComponentList: ComponentList{
				Components: map[string]wunit.Concentration{
					"LB": wunit.NewConcentration(0, "g/L"),
				},
			},
			Volume: wunit.NewVolume(1.05e+04-(5.26e+03+351), "ul"),
		},
		sample2: ComponentListSample{
			ComponentList: ComponentList{
				Components: map[string]wunit.Concentration{
					"Ferric Chloride (uM)": wunit.NewConcentration(20, "mM"),
				},
			},
			Volume: wunit.NewVolume(5.26e+03, "ul"),
		},
		sample3: ComponentListSample{
			ComponentList: ComponentList{
				Components: map[string]wunit.Concentration{
					"Glucose (g/L)": wunit.NewConcentration(150, "g/L"),
				},
			},
			Volume: wunit.NewVolume(351, "ul"),
		},
		mixedList: ComponentList{
			Components: map[string]wunit.Concentration{
				"LB": wunit.NewConcentration(0, "g/L"),
				"Ferric Chloride (uM)": wunit.NewConcentration(10, "mM"),
				"Glucose (g/L)":        wunit.NewConcentration(5.01, "g/L"),
			},
		},
	},
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

		err = EqualLists(mixed, test.mixedList)

		if err != nil {
			t.Error(
				"For", test.name, "\n",
				"expected:", "\n",
				test.mixedList,
				"got:", "\n",
				mixed,
				"Error: ", "\n",
				err.Error(),
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

		err = EqualLists(mixed, test.mixedList)

		if err != nil {
			t.Error(
				"For", test.name, "\n",
				"expected:", "\n",
				test.mixedList,
				"got:", "\n",
				mixed,
				"Error: ", "\n",
				err.Error(),
			)
		}
	}
}
