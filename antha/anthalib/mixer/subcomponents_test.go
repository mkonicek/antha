// solutions
package mixer

import (
	"errors"
	"fmt"
	"testing"

	. "github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/solutions"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
)

type mixComponentlistTest struct {
	name      string
	sample1   wtype.ComponentListSample
	sample2   wtype.ComponentListSample
	mixedList wtype.ComponentList
}

var tests []mixComponentlistTest = []mixComponentlistTest{
	{
		sample1: wtype.ComponentListSample{
			ComponentList: wtype.ComponentList{
				Components: map[string]wunit.Concentration{
					"water": wunit.NewConcentration(1, "g/L"),
				},
			},
			Volume: wunit.NewVolume(1, "ul"),
		},
		sample2: wtype.ComponentListSample{
			ComponentList: wtype.ComponentList{
				Components: map[string]wunit.Concentration{
					"dna":  wunit.NewConcentration(1, "g/L"),
					"dna2": wunit.NewConcentration(2, "X"),
				},
			},
			Volume: wunit.NewVolume(1, "ul"),
		},
		mixedList: wtype.ComponentList{
			Components: map[string]wunit.Concentration{
				"water": wunit.NewConcentration(0.5, "g/L"),
				"dna":   wunit.NewConcentration(0.5, "g/L"),
				"dna2":  wunit.NewConcentration(1, "X"),
			},
		},
	},
	{
		sample1: wtype.ComponentListSample{
			ComponentList: wtype.ComponentList{
				Components: map[string]wunit.Concentration{
					"water": wunit.NewConcentration(1, "g/L"),
					"dna":   wunit.NewConcentration(1, "g/L"),
				},
			},
			Volume: wunit.NewVolume(1, "ul"),
		},
		sample2: wtype.ComponentListSample{
			ComponentList: wtype.ComponentList{
				Components: map[string]wunit.Concentration{
					"dna": wunit.NewConcentration(1, "g/L"),
				},
			},
			Volume: wunit.NewVolume(1, "ul"),
		},
		mixedList: wtype.ComponentList{
			Components: map[string]wunit.Concentration{
				"water": wunit.NewConcentration(0.5, "g/L"),
				"dna":   wunit.NewConcentration(1, "g/L"),
			},
		},
	},
	{
		sample1: wtype.ComponentListSample{
			ComponentList: wtype.ComponentList{
				Components: map[string]wunit.Concentration{
					"water": wunit.NewConcentration(1, "g/l"),
					"dna":   wunit.NewConcentration(1, "g/l"),
				},
			},
			Volume: wunit.NewVolume(1, "ul"),
		},
		sample2: wtype.ComponentListSample{
			ComponentList: wtype.ComponentList{
				Components: map[string]wunit.Concentration{
					"dna": wunit.NewConcentration(1000, "mg/l"),
				},
			},
			Volume: wunit.NewVolume(1, "ul"),
		},
		mixedList: wtype.ComponentList{
			Components: map[string]wunit.Concentration{
				"water": wunit.NewConcentration(0.5, "g/l"),
				"dna":   wunit.NewConcentration(1, "g/l"),
			},
		},
	},
	{
		sample1: wtype.ComponentListSample{
			ComponentList: wtype.ComponentList{
				Components: map[string]wunit.Concentration{
					"water":    wunit.NewConcentration(1, "g/L"),
					"glycerol": wunit.NewConcentration(1, "M"),
				},
			},
			Volume: wunit.NewVolume(1, "ul"),
		},
		sample2: wtype.ComponentListSample{
			ComponentList: wtype.ComponentList{
				Components: map[string]wunit.Concentration{
					"glycerol": wunit.NewConcentration(1, "g/L"),
				},
			},
			Volume: wunit.NewVolume(1, "ul"),
		},
		mixedList: wtype.ComponentList{
			Components: map[string]wunit.Concentration{
				"water":    wunit.NewConcentration(0.5, "g/L"),
				"glycerol": wunit.NewConcentration(46.5, "g/L"),
			},
		},
	},
}

type serialComponentlistTest struct {
	name      string
	sample1   wtype.ComponentListSample
	sample2   wtype.ComponentListSample
	sample3   wtype.ComponentListSample
	mixedList wtype.ComponentList
}

var serialTests []serialComponentlistTest = []serialComponentlistTest{
	{
		sample1: wtype.ComponentListSample{
			ComponentList: wtype.ComponentList{
				Components: map[string]wunit.Concentration{
					"water": wunit.NewConcentration(0, "g/L"),
				},
			},
			Volume: wunit.NewVolume(8, "ul"),
		},
		sample2: wtype.ComponentListSample{
			ComponentList: wtype.ComponentList{
				Components: map[string]wunit.Concentration{
					"dna": wunit.NewConcentration(1, "g/L"),
				},
			},
			Volume: wunit.NewVolume(1, "ul"),
		},
		sample3: wtype.ComponentListSample{
			ComponentList: wtype.ComponentList{
				Components: map[string]wunit.Concentration{
					"dna2": wunit.NewConcentration(1, "g/L"),
				},
			},
			Volume: wunit.NewVolume(1, "ul"),
		},
		mixedList: wtype.ComponentList{
			Components: map[string]wunit.Concentration{
				"water": wunit.NewConcentration(0, "g/L"),
				"dna":   wunit.NewConcentration(0.1, "g/L"),
				"dna2":  wunit.NewConcentration(0.1, "g/L"),
			},
		},
	},
	{
		sample1: wtype.ComponentListSample{
			ComponentList: wtype.ComponentList{
				Components: map[string]wunit.Concentration{
					"LB": wunit.NewConcentration(0, "g/L"),
				},
			},
			Volume: wunit.NewVolume(1.05e+04-(5.26e+03+351), "ul"),
		},
		sample2: wtype.ComponentListSample{
			ComponentList: wtype.ComponentList{
				Components: map[string]wunit.Concentration{
					"Ferric Chloride (uM)": wunit.NewConcentration(20, "mM"),
				},
			},
			Volume: wunit.NewVolume(5.26e+03, "ul"),
		},
		sample3: wtype.ComponentListSample{
			ComponentList: wtype.ComponentList{
				Components: map[string]wunit.Concentration{
					"Glucose (g/L)": wunit.NewConcentration(150, "g/L"),
				},
			},
			Volume: wunit.NewVolume(351, "ul"),
		},
		mixedList: wtype.ComponentList{
			Components: map[string]wunit.Concentration{
				"LB":                   wunit.NewConcentration(0, "g/L"),
				"Ferric Chloride (uM)": wunit.NewConcentration(10, "mM"),
				"Glucose (g/L)":        wunit.NewConcentration(5.01, "g/L"),
			},
		},
	},
}

func TestSimulateMix(t *testing.T) {
	for _, test := range tests {
		mixed, err := wtype.MixComponentLists(test.sample1, test.sample2)

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
		intermediate, err := wtype.MixComponentLists(test.sample1, test.sample2)

		if err != nil {
			t.Error(
				"For", test.name, "\n",
				"got error:", err.Error(), "\n",
			)
		}

		intermediateSample := wtype.ComponentListSample{
			ComponentList: intermediate,
			Volume:        wunit.AddVolumes(test.sample1.Volume, test.sample2.Volume),
		}

		mixed, err := wtype.MixComponentLists(intermediateSample, test.sample3)
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
		product                      *wtype.Liquid
		mixes                        []*wtype.Liquid
		expectedProductName          string
		expectedProductComponentList wtype.ComponentList
		expectedProductConc          wunit.Concentration
		expectedError                error
	}

	var defaultConc wunit.Concentration

	newTestComponent := func(name string, typ wtype.LiquidType, smax float64, conc wunit.Concentration, vol wunit.Volume, componentList wtype.ComponentList) *wtype.Liquid {
		c := wtype.NewLHComponent()
		c.SetName(name)
		c.Type = typ
		c.Smax = smax
		if conc != defaultConc {
			c.SetConcentration(conc)
		}
		if err := c.AddSubComponents(componentList); err != nil {
			t.Fatal(err)
		}

		return c
	}

	gPerL0 := wunit.NewConcentration(0.0, "g/L")
	gPerL1 := wunit.NewConcentration(1, "g/L")

	var nilComponentList wtype.ComponentList

	someComponents := wtype.ComponentList{Components: map[string]wunit.Concentration{
		"glycerol": wunit.NewConcentration(0.25, "g/l"),
		"IPTG":     wunit.NewConcentration(0.25, "mMol/l"),
		"water":    wunit.NewConcentration(0.25, "v/v"),
		"LB":       wunit.NewConcentration(0.25, "X"),
	},
	}

	someOtherComponents := wtype.ComponentList{Components: map[string]wunit.Concentration{
		"glycerol":    wunit.NewConcentration(0.5, "g/l"),
		"IPTG":        wunit.NewConcentration(0.5, "mMol/l"),
		"water":       wunit.NewConcentration(0.5, "v/v"),
		"LB":          wunit.NewConcentration(0.25, "X"),
		"Extra Thing": wunit.NewConcentration(1, "X"),
	},
	}

	lbComponents := wtype.ComponentList{Components: map[string]wunit.Concentration{
		"Yeast Extract":   wunit.NewConcentration(5, "g/l"),
		"Tryptone":        wunit.NewConcentration(10, "g/l"),
		"Sodium Chloride": wunit.NewConcentration(10, "g/l"),
	},
	}

	cutSmartComponents := wtype.ComponentList{Components: map[string]wunit.Concentration{
		"BSA":               wunit.NewConcentration(1000, "mg/l"),
		"Magnesium Acetate": wunit.NewConcentration(100, "mM"),
		"Potassium Acetate": wunit.NewConcentration(500, "mM"),
		"Tris-acetate":      wunit.NewConcentration(200, "mM"),
	},
	}

	conc := func(s string) wunit.Concentration {
		return wunit.NewConcentration(wunit.SplitValueAndUnit(s))
	}
	sapISubComponents := wtype.ComponentList{Components: map[string]wunit.Concentration{
		"BSA":      conc("500mg/l"),
		"DTT":      conc("1mM"),
		"EDTA":     conc("0.1mM"),
		"Glycerol": conc("500g/l"),
		"NaCl":     conc("10mM"),
		"Tris-HCl": conc("10mM"),
	},
	}

	t4SubComponents := wtype.ComponentList{Components: map[string]wunit.Concentration{
		"DTT":      conc("1mM"),
		"EDTA":     conc("0.1mM"),
		"Glycerol": conc("500g/l"),
		"KCl":      conc("50mM"),
		"Tris-HCl": conc("10mM"),
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

	ws := Sample(water, wunit.NewVolume(65.0, "ul"))
	wsTotal := SampleForTotalVolume(water, wunit.NewVolume(100.0, "ul"))
	mmxs := Sample(mmx, wunit.NewVolume(25.0, "ul"))
	ps := Sample(part, wunit.NewVolume(10.0, "ul"))

	cutSmart := newTestComponent("CutsmartBuffer",
		wtype.LTWater,
		9999,
		wunit.NewConcentration(10.0, "X"),
		wunit.NewVolume(2000.0, "ul"),
		cutSmartComponents)

	sapI := newTestComponent("SapI",
		wtype.LTWater,
		9999,
		wunit.NewConcentration(10000, "U/ml"),
		wunit.NewVolume(2000.0, "ul"),
		sapISubComponents)

	t4 := newTestComponent("T4Ligase",
		wtype.LTWater,
		9999,
		wunit.NewConcentration(400000, "U/ml"),
		wunit.NewVolume(2000.0, "ul"),
		t4SubComponents)

	atp := newTestComponent("ATP",
		wtype.LTWater,
		9999,
		wunit.NewConcentration(1, "mM"),
		wunit.NewVolume(2000.0, "ul"),
		nilComponentList)

	var mixTests = []mixTest{
		{
			name:    "mmxTest",
			product: water,
			mixes: []*wtype.Liquid{
				Sample(cutSmart, wunit.NewVolume(10.0, "ul")),
				Sample(atp, wunit.NewVolume(5.0, "ul")),
				Sample(sapI, wunit.NewVolume(5.0, "ul")),
				Sample(t4, wunit.NewVolume(5.0, "ul")),
			},
			expectedProductName: "0.2 mMol/l ATP+500 mg/l BSA+0.4 mMol/l DTT+0.04 mMol/l EDTA+200 g/l Glycerol+10 mMol/l KCl+40 mMol/l Magnesium Acetate+2 mMol/l NaCl+200 mMol/l Potassium Acetate+4 mMol/l Tris-HCl+80 mMol/l Tris-acetate",
			expectedProductComponentList: wtype.ComponentList{
				Components: map[string]wunit.Concentration{
					"BSA":               conc("500mg/l"),
					"DTT":               conc("0.4mM"),
					"EDTA":              conc("0.04mM"),
					"Glycerol":          conc("200g/l"),
					"KCl":               conc("10mM"),
					"NaCl":              conc("2mM"),
					"Tris-HCl":          conc("4mM"),
					"Magnesium Acetate": wunit.NewConcentration(40, "mM"),
					"Potassium Acetate": wunit.NewConcentration(200, "mM"),
					"Tris-acetate":      wunit.NewConcentration(80, "mM"),
					"ATP":               wunit.NewConcentration(0.2, "mM"),
				},
			},
			expectedProductConc: gPerL0,
			expectedError:       nil,
		},
		{
			name:                "noConcsTest",
			product:             water,
			mixes:               []*wtype.Liquid{ws, ps, mmxs},
			expectedProductName: "0.1 v/v dna+0.25 v/v mastermix_sapI+0.65 v/v water",
			expectedProductComponentList: wtype.ComponentList{
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
			mixes:               []*wtype.Liquid{wsTotal, ps, mmxs},
			expectedProductName: "0.1 v/v dna+0.25 v/v mastermix_sapI+0.65 v/v water",
			expectedProductComponentList: wtype.ComponentList{
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
			mixes: []*wtype.Liquid{
				Sample(water, wunit.NewVolume(65.0, "ul")),
				Sample(glycerol, wunit.NewVolume(65.0, "ul")),
				Sample(iptg, wunit.NewVolume(65.0, "ul")),
				Sample(lb, wunit.NewVolume(65.0, "ul")),
			},
			expectedProductName: "0.25 mMol/l IPTG+0.25 X LB+0.25 g/l glycerol+0.25 v/v water",
			expectedProductComponentList: wtype.ComponentList{
				Components: map[string]wunit.Concentration{
					"glycerol": wunit.NewConcentration(0.25, "g/l"),
					"IPTG":     wunit.NewConcentration(0.25, "mMol/l"),
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
			mixes: []*wtype.Liquid{
				SampleForTotalVolume(water, wunit.NewVolume(100.0, "ul")),
				Sample(water, wunit.NewVolume(65.0, "ul")),
				Sample(glycerol, wunit.NewVolume(65.0, "ul")),
				Sample(iptg, wunit.NewVolume(65.0, "ul")),
				Sample(lb, wunit.NewVolume(65.0, "ul")),
			},
			expectedProductName: "water",
			expectedProductComponentList: wtype.ComponentList{
				Components: map[string]wunit.Concentration{},
			},
			expectedProductConc: gPerL0,
			expectedError:       errors.New("SampleForTotalVolume requested (100 ul) is less than sum of sample volumes (260 ul)"),
		},
		{
			name:    "renamedComponentTest",
			product: water,
			mixes: []*wtype.Liquid{
				Sample(lbWithSubComponents, wunit.NewVolume(400.0, "ul")),
				Sample(glycerol, wunit.NewVolume(50, "ul")),
				Sample(iptg, wunit.NewVolume(50, "ul")),
			},
			expectedProductName: "0.1 mMol/l IPTG+8 g/l Sodium Chloride+8 g/l Tryptone+4 g/l Yeast Extract+0.1 g/l glycerol",
			expectedProductComponentList: wtype.ComponentList{
				Components: map[string]wunit.Concentration{
					"glycerol":        wunit.NewConcentration(0.1, "g/l"),
					"IPTG":            wunit.NewConcentration(0.1, "mMol/l"),
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
			mixes: []*wtype.Liquid{
				Sample(water, wunit.NewVolume(600.0, "ul")),
				Sample(glycerol, wunit.NewVolume(100.0, "ul")),
				Sample(iptg, wunit.NewVolume(100.0, "ul")),
				Sample(mediaMixture, wunit.NewVolume(100.0, "ul")),
				Sample(anotherMediaMixture, wunit.NewVolume(100.0, "ul")),
			},
			expectedProductName: "0.1 X Extra Thing+0.175 mMol/l IPTG+0.05 X LB+0.175 g/l glycerol+0.675 v/v water",
			expectedProductComponentList: wtype.ComponentList{
				Components: map[string]wunit.Concentration{
					"glycerol":    wunit.NewConcentration(0.175, "g/l"),
					"IPTG":        wunit.NewConcentration(0.175, "mMol/l"),
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
			mixes: []*wtype.Liquid{
				Sample(water, wunit.NewVolume(65.0, "ul")),
				Sample(glycerol, wunit.NewVolume(65.0, "ul")),
				Sample(iptg, wunit.NewVolume(65.0, "ul")),
				Sample(mediaMixture, wunit.NewVolume(65.0, "ul")),
			},
			expectedProductName: "0.312 mMol/l IPTG+0.0625 X LB+0.312 g/l glycerol+0.312 v/v water",
			expectedProductComponentList: wtype.ComponentList{
				Components: map[string]wunit.Concentration{
					"glycerol": wunit.NewConcentration(0.3125, "g/l"),
					"IPTG":     wunit.NewConcentration(0.3125, "mMol/l"),
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
		err := wtype.UpdateComponentDetails(duplicatedProduct, test.mixes...)

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

		testCompList, err := duplicatedProduct.GetSubComponents()

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
