// solutions
package solutions

import (
	"errors"
	"fmt"
	"testing"

	"github.com/antha-lang/antha/antha/anthalib/mixer"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
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
		if err != nil {
			t.Error(err)
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

func TestUpdateComponentDetails(t *testing.T) {

	type mixTest struct {
		name                         string
		product                      *wtype.LHComponent
		mixes                        []*wtype.LHComponent
		expectedProductName          string
		expectedProductComponentList ComponentList
		expectedProductConc          wunit.Concentration
		expectedError                error
	}

	var defaultConc wunit.Concentration

	newTestComponent := func(name string, typ wtype.LiquidType, smax float64, conc wunit.Concentration, vol wunit.Volume, componentList ComponentList) *wtype.LHComponent {
		c := wtype.NewLHComponent()
		c.SetName(name)
		c.Type = typ
		c.Smax = smax
		if conc != defaultConc {
			c.SetConcentration(conc)
		}
		AddSubComponents(c, componentList)
		return c
	}

	gPerL0 := wunit.NewConcentration(0.0, "g/L")
	gPerL1 := wunit.NewConcentration(1, "g/L")

	var nilComponentList ComponentList

	someComponents := ComponentList{Components: map[string]wunit.Concentration{
		"glycerol": wunit.NewConcentration(0.25, "g/l"),
		"IPTG":     wunit.NewConcentration(0.25, "mM/l"),
		"water":    wunit.NewConcentration(0.25, "v/v"),
		"LB":       wunit.NewConcentration(0.25, "X"),
	},
	}

	someOtherComponents := ComponentList{Components: map[string]wunit.Concentration{
		"glycerol":    wunit.NewConcentration(0.5, "g/l"),
		"IPTG":        wunit.NewConcentration(0.5, "mM/l"),
		"water":       wunit.NewConcentration(0.5, "v/v"),
		"LB":          wunit.NewConcentration(0.25, "X"),
		"Extra Thing": wunit.NewConcentration(1, "X"),
	},
	}

	lbComponents := ComponentList{Components: map[string]wunit.Concentration{
		"Yeast Extract":   wunit.NewConcentration(5, "g/l"),
		"Tryptone":        wunit.NewConcentration(10, "g/l"),
		"Sodium Chloride": wunit.NewConcentration(10, "g/l"),
	},
	}

	water := newTestComponent("water", wtype.LTWater, 9999, defaultConc, wunit.NewVolume(2000.0, "ul"), nilComponentList)

	mmx := newTestComponent("mastermix_sapI", wtype.LTWater, 9999, defaultConc, wunit.NewVolume(2000.0, "ul"), nilComponentList)

	part := newTestComponent("dna", wtype.LTWater, 9999, defaultConc, wunit.NewVolume(2000.0, "ul"), nilComponentList)

	glycerol := newTestComponent("glycerol", wtype.LTWater, 9999, gPerL1, wunit.NewVolume(2000.0, "ul"), nilComponentList)
	iptg := newTestComponent("IPTG", wtype.LTWater, 9999, wunit.NewConcentration(1, "mM"), wunit.NewVolume(2000.0, "ul"), nilComponentList)
	lb := newTestComponent("LB", wtype.LTWater, 9999, wunit.NewConcentration(1, "X"), wunit.NewVolume(2000.0, "ul"), nilComponentList)
	lbWithSubComponents := newTestComponent("LB", wtype.LTWater, 9999, wunit.NewConcentration(1, "X"), wunit.NewVolume(2000.0, "ul"), lbComponents)
	mediaMixture := newTestComponent("LB", wtype.LTWater, 9999, wunit.NewConcentration(1, "X"), wunit.NewVolume(2000.0, "ul"), someComponents)
	anotherMediaMixture := newTestComponent("LB", wtype.LTWater, 9999, wunit.NewConcentration(1, "X"), wunit.NewVolume(2000.0, "ul"), someOtherComponents)

	ws := mixer.Sample(water, wunit.NewVolume(65.0, "ul"))
	wsTotal := mixer.SampleForTotalVolume(water, wunit.NewVolume(100.0, "ul"))
	mmxs := mixer.Sample(mmx, wunit.NewVolume(25.0, "ul"))
	ps := mixer.Sample(part, wunit.NewVolume(10.0, "ul"))

	var mixTests = []mixTest{
		{
			name:                "noConcsTest",
			product:             water,
			mixes:               []*wtype.LHComponent{ws, ps, mmxs},
			expectedProductName: "0.1 v/v dna+0.25 v/v mastermix_sapI+0.65 v/v water",
			expectedProductComponentList: ComponentList{
				Components: map[string]wunit.Concentration{
					"dna":            wunit.NewConcentration(0.1, "v/v"),
					"mastermix_sapI": wunit.NewConcentration(0.25, "v/v"),
					"water":          wunit.NewConcentration(0.65, "v/v"),
				},
			},
			expectedProductConc: gPerL0,
			expectedError:       wtype.NewWarningf("zero concentration found for sample water; zero concentration found for sample dna; zero concentration found for sample mastermix_sapI"),
		}, {
			name:                "SampleForTotalVolumeTest",
			product:             water,
			mixes:               []*wtype.LHComponent{wsTotal, ps, mmxs},
			expectedProductName: "0.1 v/v dna+0.25 v/v mastermix_sapI+0.65 v/v water",
			expectedProductComponentList: ComponentList{
				Components: map[string]wunit.Concentration{
					"dna":            wunit.NewConcentration(0.1, "v/v"),
					"mastermix_sapI": wunit.NewConcentration(0.25, "v/v"),
					"water":          wunit.NewConcentration(0.65, "v/v"),
				},
			},
			expectedProductConc: gPerL0,
			expectedError:       wtype.NewWarningf("zero concentration found for sample water; zero concentration found for sample dna; zero concentration found for sample mastermix_sapI"),
		},
		{
			name:    "SampleWithConcsTest",
			product: water,
			mixes: []*wtype.LHComponent{
				mixer.Sample(water, wunit.NewVolume(65.0, "ul")),
				mixer.Sample(glycerol, wunit.NewVolume(65.0, "ul")),
				mixer.Sample(iptg, wunit.NewVolume(65.0, "ul")),
				mixer.Sample(lb, wunit.NewVolume(65.0, "ul")),
			},
			expectedProductName: "0.25 mM/l IPTG+0.25 X LB+0.25 g/l glycerol+0.25 v/v water",
			expectedProductComponentList: ComponentList{
				Components: map[string]wunit.Concentration{
					"glycerol": wunit.NewConcentration(0.25, "g/l"),
					"IPTG":     wunit.NewConcentration(0.25, "mM/l"),
					"water":    wunit.NewConcentration(0.25, "v/v"),
					"LB":       wunit.NewConcentration(0.25, "X"),
				},
			},
			expectedProductConc: gPerL0,
			expectedError:       wtype.NewWarningf("zero concentration found for sample water"),
		},
		{
			name:    "InvalidTotalVolumeTest",
			product: water,
			mixes: []*wtype.LHComponent{
				mixer.SampleForTotalVolume(water, wunit.NewVolume(100.0, "ul")),
				mixer.Sample(water, wunit.NewVolume(65.0, "ul")),
				mixer.Sample(glycerol, wunit.NewVolume(65.0, "ul")),
				mixer.Sample(iptg, wunit.NewVolume(65.0, "ul")),
				mixer.Sample(lb, wunit.NewVolume(65.0, "ul")),
			},
			expectedProductName: "water",
			expectedProductComponentList: ComponentList{
				Components: map[string]wunit.Concentration{},
			},
			expectedProductConc: gPerL0,
			expectedError:       errors.New("SampleForTotalVolume requested (100 ul) is less than sum of sample volumes (260 ul)"),
		},
		{
			name:    "renamedComponentTest",
			product: water,
			mixes: []*wtype.LHComponent{
				mixer.Sample(lbWithSubComponents, wunit.NewVolume(400.0, "ul")),
				mixer.Sample(glycerol, wunit.NewVolume(50, "ul")),
				mixer.Sample(iptg, wunit.NewVolume(50, "ul")),
			},
			expectedProductName: "0.1 mM/l IPTG+8 g/l Sodium Chloride+8 g/l Tryptone+4 g/l Yeast Extract+0.1 g/l glycerol",
			expectedProductComponentList: ComponentList{
				Components: map[string]wunit.Concentration{
					"glycerol":        wunit.NewConcentration(0.1, "g/l"),
					"IPTG":            wunit.NewConcentration(0.1, "mM/l"),
					"Yeast Extract":   wunit.NewConcentration(4, "g/l"),
					"Tryptone":        wunit.NewConcentration(8, "g/l"),
					"Sodium Chloride": wunit.NewConcentration(8, "g/l"),
				},
			},
			expectedProductConc: gPerL0,
			expectedError:       nil,
		},
		{
			name:    "SampleWithTwoComponentListsTest",
			product: water,
			mixes: []*wtype.LHComponent{
				mixer.Sample(water, wunit.NewVolume(600.0, "ul")),
				mixer.Sample(glycerol, wunit.NewVolume(100.0, "ul")),
				mixer.Sample(iptg, wunit.NewVolume(100.0, "ul")),
				mixer.Sample(mediaMixture, wunit.NewVolume(100.0, "ul")),
				mixer.Sample(anotherMediaMixture, wunit.NewVolume(100.0, "ul")),
			},
			expectedProductName: "0.1 X Extra Thing+0.175 mM/l IPTG+0.05 X LB+0.175 g/l glycerol+0.675 v/v water",
			expectedProductComponentList: ComponentList{
				Components: map[string]wunit.Concentration{
					"glycerol":    wunit.NewConcentration(0.175, "g/l"),
					"IPTG":        wunit.NewConcentration(0.175, "mM/l"),
					"water":       wunit.NewConcentration(0.675, "v/v"),
					"LB":          wunit.NewConcentration(0.05, "X"),
					"Extra Thing": wunit.NewConcentration(0.1, "X"),
				},
			},
			expectedProductConc: gPerL0,
			expectedError:       wtype.NewWarningf("zero concentration found for sample water"),
		},
		{
			name:    "SampleWithComponentListsTest",
			product: water,
			mixes: []*wtype.LHComponent{
				mixer.Sample(water, wunit.NewVolume(65.0, "ul")),
				mixer.Sample(glycerol, wunit.NewVolume(65.0, "ul")),
				mixer.Sample(iptg, wunit.NewVolume(65.0, "ul")),
				mixer.Sample(mediaMixture, wunit.NewVolume(65.0, "ul")),
			},
			expectedProductName: "0.312 mM/l IPTG+0.0625 X LB+0.312 g/l glycerol+0.312 v/v water",
			expectedProductComponentList: ComponentList{
				Components: map[string]wunit.Concentration{
					"glycerol": wunit.NewConcentration(0.3125, "g/l"),
					"IPTG":     wunit.NewConcentration(0.3125, "mM/l"),
					"water":    wunit.NewConcentration(0.312, "v/v"),
					"LB":       wunit.NewConcentration(0.0625, "X"),
				},
			},
			expectedProductConc: gPerL0,
			expectedError:       wtype.NewWarningf("zero concentration found for sample water"),
		},
	}

	for _, test := range mixTests {
		duplicatedProduct := test.product.Dup()
		err := UpdateComponentDetails(duplicatedProduct, test.mixes...)

		if err != nil {
			if err.Error() != test.expectedError.Error() {
				t.Error(
					"For :", test.name, "\n",
					"expected error:", fmt.Sprintf("%v %T ", test.expectedError, test.expectedError), "\n",
					"got error:", fmt.Sprintf("%v %T", err, err), "\n",
				)
			}
		} else if test.expectedError != nil {
			t.Error(
				"For :", test.name, "\n",
				"expected error:", test.expectedError.Error(), "\n",
				"got no error:", "\n",
			)
		}

		if duplicatedProduct.Name() != test.expectedProductName {
			t.Error(
				"For :", test.name, "\n",
				"expected name:", test.expectedProductName, "\n",
				"got:", duplicatedProduct.Name(), "\n",
			)
		}

		testCompList, err := GetSubComponents(duplicatedProduct)

		if err != nil {
			fmt.Println(err.Error())
		}

		err = EqualLists(testCompList, test.expectedProductComponentList)

		if err != nil {
			t.Error(
				"For: ", test.name, "\n",
				"expected:", "\n",
				test.expectedProductComponentList,
				"got:", "\n",
				testCompList,
				"Error: ", "\n",
				err.Error(),
			)
		}

		if !duplicatedProduct.Concentration().EqualTo(test.expectedProductConc) {
			t.Error(
				"For: ", test.name, "\n",
				"expected:", "\n",
				test.expectedProductConc,
				"got:", "\n",
				duplicatedProduct.Concentration(),
			)
		}
	}
}
