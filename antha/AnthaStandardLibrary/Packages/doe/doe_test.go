// doe.go
package doe

import (
	//"fmt"
	//"strconv"
	//	"strings"

	"runtime/debug"
	"testing"
	//"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/search"
	//"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/spreadsheet"
	//"github.com/antha-lang/antha/internal/github.com/tealeg/xlsx"
)

// simple reverse complement check to test testing methodology initially
type testpair struct {
	pairs         []DOEPair
	combocount    int
	factorheaders []string
}

var factorsandlevels = []testpair{

	{pairs: []DOEPair{Pair("Level 1", []interface{}{1})},
		combocount: 1, factorheaders: []string{"Level 1"}},
	{pairs: []DOEPair{Pair("Level 1", []interface{}{1}), Pair("Level 2", []interface{}{1})},
		combocount: 1, factorheaders: []string{"Level 1", "Level 2"}},
	{pairs: []DOEPair{Pair("Level 1", []interface{}{1}), Pair("Level 2", []interface{}{1, 2})},
		combocount: 2, factorheaders: []string{"Level 1", "Level 2"}},
	{pairs: []DOEPair{Pair("Level 1", []interface{}{1, 2}), Pair("Level 2", []interface{}{1})},
		combocount: 2, factorheaders: []string{"Level 1", "Level 2"}},
	{pairs: []DOEPair{Pair("Level 1", []interface{}{1, 2}), Pair("Level 2", []interface{}{1, 2})},
		combocount: 4, factorheaders: []string{"Level 1", "Level 2"}},
	{pairs: []DOEPair{Pair("Level 1", []interface{}{1}), Pair("Level 2", []interface{}{1}), Pair("Level 3", []interface{}{1})},
		combocount: 1, factorheaders: []string{"Level 1", "Level 2", "Level 3"}},
	{pairs: []DOEPair{Pair("Level 1", []interface{}{1, 2}), Pair("Level 2", []interface{}{1, 2}), Pair("Level 3", []interface{}{1, 2})},
		combocount: 8, factorheaders: []string{"Level 1", "Level 2", "Level 3"}},
	{pairs: []DOEPair{Pair("Level 1", []interface{}{1}), Pair("Level 2", []interface{}{1, 2}), Pair("Level 3", []interface{}{1, 2})},
		combocount: 4, factorheaders: []string{"Level 1", "Level 2", "Level 3"}},
	{pairs: []DOEPair{Pair("Level 1", []interface{}{1}), Pair("Level 2", []interface{}{1, 2}), Pair("Level 3", []interface{}{1})},
		combocount: 2, factorheaders: []string{"Level 1", "Level 2", "Level 3"}},
	{pairs: []DOEPair{Pair("Level 1", []interface{}{1, 2}), Pair("Level 2", []interface{}{1, 2}), Pair("Level 3", []interface{}{1})},
		combocount: 4, factorheaders: []string{"Level 1", "Level 2", "Level 3"}},
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
				"got", r, "\n",
			)
		}
		/*	for j, run := range r {

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
		}*/
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
	fmt.Println("In AllComboCount", "len(pairs)", len(pairs))
	var movingcount int
	movingcount = (pairs[0]).LevelCount()
	fmt.Println("levelcount", movingcount)
	fmt.Println("len(levels)", len(pairs[0].Levels))
	for i := 1; i < len(pairs); i++ {
		fmt.Println("levelcount", movingcount)
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
		fmt.Println("factor", i, "of", AllComboCount(factors), factor.Factor, factor.Levels)
		for j, level := range factor.Levels {
			//fmt.Println("factor:", factor, level, i, j)
		descriptors = append(descriptors, factor.Factor)
			setpoints = append(setpoints, level)
			run.Factordescriptors = descriptors
			run.Setpoints = setpoints
			//	fmt.Println("factor:", factor, i, j)
			runs[i+j] = run
		}
	}
	return
}
*/
