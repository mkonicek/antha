// doe.go
package doe

import (
	"fmt"
	"reflect"
	"runtime/debug"
	"testing"

	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/go-test/deep"
)

// simple reverse complement check to test testing methodology initially
type testpair struct {
	pairs         []DOEPair
	combocount    int
	factorheaders []string
}

var factorsandlevels = []testpair{

	{pairs: []DOEPair{Pair("Factor 1", []interface{}{1})},
		combocount: 1, factorheaders: []string{"Factor 1"}},
	{pairs: []DOEPair{Pair("Factor 1", []interface{}{1}), Pair("Factor 2", []interface{}{1})},
		combocount: 1, factorheaders: []string{"Factor 1", "Factor 2"}},
	{pairs: []DOEPair{Pair("Factor 1", []interface{}{1}), Pair("Factor 2", []interface{}{1, 2})},
		combocount: 2, factorheaders: []string{"Factor 1", "Factor 2"}},
	{pairs: []DOEPair{Pair("Factor 1", []interface{}{1, 2}), Pair("Factor 2", []interface{}{1})},
		combocount: 2, factorheaders: []string{"Factor 1", "Factor 2"}},
	{pairs: []DOEPair{Pair("Factor 1", []interface{}{1, 2}), Pair("Factor 2", []interface{}{1, 2})},
		combocount: 4, factorheaders: []string{"Factor 1", "Factor 2"}},
	{pairs: []DOEPair{Pair("Factor 1", []interface{}{1}), Pair("Factor 2", []interface{}{1}), Pair("Factor 3", []interface{}{1})},
		combocount: 1, factorheaders: []string{"Factor 1", "Factor 2", "Factor 3"}},
	{pairs: []DOEPair{Pair("Factor 1", []interface{}{1, 2}), Pair("Factor 2", []interface{}{1, 2}), Pair("Factor 3", []interface{}{1, 2})},
		combocount: 8, factorheaders: []string{"Factor 1", "Factor 2", "Factor 3"}},
	{pairs: []DOEPair{Pair("Factor 1", []interface{}{1}), Pair("Factor 2", []interface{}{1, 2}), Pair("Factor 3", []interface{}{1, 2})},
		combocount: 4, factorheaders: []string{"Factor 1", "Factor 2", "Factor 3"}},
	{pairs: []DOEPair{Pair("Factor 1", []interface{}{1}), Pair("Factor 2", []interface{}{1, 2}), Pair("Factor 3", []interface{}{1})},
		combocount: 2, factorheaders: []string{"Factor 1", "Factor 2", "Factor 3"}},
	{pairs: []DOEPair{Pair("Factor 1", []interface{}{1, 2}), Pair("Factor 2", []interface{}{1, 2}), Pair("Factor 3", []interface{}{1})},
		combocount: 4, factorheaders: []string{"Factor 1", "Factor 2", "Factor 3"}},
}

func TestAllComboCount(t *testing.T) {
	for _, factor := range factorsandlevels {
		r := AllComboCount(factor.pairs)
		if r != factor.combocount {
			t.Error(
				"For", factor.pairs, "/n",
				"expected", factor.combocount, "\n",
				"got", r, "\n",
			)
		}
	}

}

func TestAllCombinations(t *testing.T) {

	defer func() {
		if res := recover(); res != nil {
			t.Fatalf("caught panic %q at %s", res, string(debug.Stack()))
		}
	}()

	for _, factor := range factorsandlevels {
		r := AllCombinations(factor.pairs)
		if len(r) != factor.combocount {
			t.Error(
				"For", factor.pairs, "/n",
				"expected", factor.combocount, "\n",
				"got", len(r), "\n",
			)
		}
		for j, run := range r {

			if len(run.Factordescriptors) != len(factor.factorheaders) {
				t.Error(
					"For", factor.pairs, "/n",
					"expected", factor.factorheaders, "\n",
					"got", run.Factordescriptors, "\n",
				)
			}

			for i, descriptor := range run.Factordescriptors {

				if descriptor != factor.factorheaders[i] {
					t.Error(
						"For", factor.pairs, "/n",
						"For", run, j, "/n",
						"descriptor", descriptor, "/n",
						"expected", factor.factorheaders[i], "\n",
						"got", run.Factordescriptors[i], "\n",
					)
				}
			}
		}
	}
}

type volFactorTest struct {
	Header       string
	Value        interface{}
	Vol          wunit.Volume
	ErrorMessage string
}

var volTests = []volFactorTest{
	{
		Header:       "Total Volume (ml)",
		Value:        interface{}(10),
		Vol:          wunit.NewVolume(10, "ml"),
		ErrorMessage: "",
	},
	{
		Header:       "Total Volume",
		Value:        interface{}("10L"),
		Vol:          wunit.NewVolume(10, "L"),
		ErrorMessage: "",
	},
	{
		Header:       "Total Volume (ul)",
		Value:        interface{}(10.0),
		Vol:          wunit.NewVolume(10, "ul"),
		ErrorMessage: "",
	},
	{
		Header:       "Total Volume (ul)",
		Value:        interface{}("10"),
		Vol:          wunit.NewVolume(10, "ul"),
		ErrorMessage: "",
	},
	{
		Header:       "Total Volume",
		Value:        interface{}("10"),
		Vol:          wunit.NewVolume(10, "ul"),
		ErrorMessage: "",
	},
	{
		Header:       "Total Volume",
		Value:        interface{}("10ml"),
		Vol:          wunit.NewVolume(10, "ml"),
		ErrorMessage: "",
	},
	{
		Header:       "Total Volume",
		Value:        interface{}("10 ml"),
		Vol:          wunit.NewVolume(10, "ml"),
		ErrorMessage: "",
	},
}

type concFactorTest struct {
	Header       string
	Value        interface{}
	Conc         wunit.Concentration
	ErrorMessage string
}

var concTests = []concFactorTest{
	{
		Header:       "Total Concentration (mg/ml)",
		Value:        interface{}(10),
		Conc:         wunit.NewConcentration(10, "mg/ml"),
		ErrorMessage: "",
	},
	{
		Header:       "Total Concentration",
		Value:        interface{}("10g/L"),
		Conc:         wunit.NewConcentration(10, "g/L"),
		ErrorMessage: "",
	},
	{
		Header:       "Total Concentration (X)",
		Value:        interface{}(10.0),
		Conc:         wunit.NewConcentration(10, "X"),
		ErrorMessage: "",
	},
	{
		Header:       "Glucose (g/L)",
		Value:        interface{}("10"),
		Conc:         wunit.NewConcentration(10, "g/L"),
		ErrorMessage: "",
	},
	{
		Header:       "1X Glucose (g/L)",
		Value:        interface{}("10"),
		Conc:         wunit.NewConcentration(10, "g/L"),
		ErrorMessage: "",
	},
	{
		Header:       "1X Glucose g/L",
		Value:        interface{}("10"),
		Conc:         wunit.Concentration{},
		ErrorMessage: "more than one unit found in header 1X Glucose g/L: valid units found [X g/L]. Units flanked by parentheses are prioritised.",
	},
	{
		Header:       "(1X) Glucose (g/L)",
		Value:        interface{}("10"),
		Conc:         wunit.Concentration{},
		ErrorMessage: "more than one unit found in header (1X) Glucose (g/L): valid units found [X g/L] []. Units flanked by parentheses are prioritised.",
	},
	{
		Header:       "Glucose g/L",
		Value:        interface{}("10"),
		Conc:         wunit.NewConcentration(10, "g/L"),
		ErrorMessage: "",
	},
	{
		Header:       "Glucose 100g/L",
		Value:        interface{}("10"),
		Conc:         wunit.NewConcentration(10, "g/L"),
		ErrorMessage: "",
	},
}

func TestHandleConcentrationFactor(t *testing.T) {
	for _, test := range concTests {
		conc, err := HandleConcFactor(test.Header, test.Value)
		if !reflect.DeepEqual(conc, test.Conc) {
			t.Error(
				"for", fmt.Sprintf("%+v", test), "\n",
				"Expected:", test.Conc, "\n",
				"Got:", conc, "\n",
			)
		}
		if err != nil {
			if err.Error() != test.ErrorMessage {
				t.Error(
					"for", fmt.Sprintf("%+v", test), "\n",
					"Expected error:", test.ErrorMessage, "\n",
					"Got:", err.Error(), "\n",
					"diffs: ", deep.Equal(test.ErrorMessage, err.Error()), "\n",
				)
			}
		}

		if err == nil && test.ErrorMessage != "" {
			t.Error(
				"for", fmt.Sprintf("%+v", test), "\n",
				"Expected error:", test.ErrorMessage, "\n",
				"Got:", "nil", "\n",
			)
		}
	}
}

func TestHandleVolumeFactor(t *testing.T) {
	for _, test := range volTests {
		vol, err := HandleVolumeFactor(test.Header, test.Value)
		if !reflect.DeepEqual(vol, test.Vol) {
			t.Error(
				"for", fmt.Sprintf("%+v", test), "\n",
				"Expected:", test.Vol, "\n",
				"Got:", vol, "\n",
			)
		}
		if err != nil {
			if err.Error() != test.ErrorMessage {
				t.Error(
					"for", fmt.Sprintf("%+v", test), "\n",
					"Expected error:", test.ErrorMessage, "\n",
					"Got:", err.Error(), "\n",
					"diffs: ", deep.Equal(test.ErrorMessage, err.Error()), "\n",
				)
			}
		}

		if err == nil && test.ErrorMessage != "" {
			t.Error(
				"for", fmt.Sprintf("%+v", test), "\n",
				"Expected error:", test.ErrorMessage, "\n",
				"Got:", "nil", "\n",
			)
		}
	}
}

/*
type DOEPair struct {
	Factor string
	Levels []interface{}
}

func (pair DOEPair) LevelCount() (numberoflevels int) {
	numberoflevels = len(pair.Levels)
	return
}
func Pair(factordescription string, levels []interface{}) (doepair DOEPair) {
	doepair.Factor = factordescription
	doepair.Levels = levels
	return
}

type Run struct {
	RunNumber            int
	StdNumber            int
	Factordescriptors    []string
	Setpoints            []interface{}
	Responsedescriptors  []string
	ResponseValues       []interface{}
	AdditionalHeaders    []string // could represent a class e.g. Environment variable, processed, raw, location
	AdditionalSubheaders []string // e.g. well ID, Ambient Temp, order,
	AdditionalValues     []interface{}
}


func AllComboCount(pairs []DOEPair) (numberofuniquecombos int) {
	// fmt.Println("In AllComboCount", "len(pairs)", len(pairs))
	var movingcount int
	movingcount = (pairs[0]).LevelCount()
	// fmt.Println("Factorcount", movingcount)
	// fmt.Println("len(levels)", len(pairs[0].Levels))
	for i := 1; i < len(pairs); i++ {
		// fmt.Println("Factorcount", movingcount)
		movingcount = movingcount * (pairs[i]).LevelCount()
	}
	numberofuniquecombos = movingcount
	return
}

func AllCombinations(factors []DOEPair) (runs []Run) {
	//fmt.Println(factors)
	descriptors := make([]string, 0)
	setpoints := make([]interface{}, 0)
	runs = make([]Run, AllComboCount(factors))
	var run Run
	for i, factor := range factors {
		// fmt.Println("factor", i, "of", AllComboCount(factors), factor.Factor, factor.Levels)
		for j, level := range factor.Levels {
			//// fmt.Println("factor:", factor, level, i, j)
		descriptors = append(descriptors, factor.Factor)
			setpoints = append(setpoints, level)
			run.Factordescriptors = descriptors
			run.Setpoints = setpoints
			//	// fmt.Println("factor:", factor, i, j)
			runs[i+j] = run
		}
	}
	return
}
*/
