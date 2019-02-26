package Optimization

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
)

func TestMsaToAssemblyProblem(t *testing.T) {

	msa := []string{
		//                S  K  G  E  E  L  F  T  G  V  V  P  I  L  V  E  L  D  G  D
		//                      ¦              ¦
		strings.ToUpper("agcaagggcgaggagctgttcaccggggtggtgcccatcctggtcgagctggacggcgac"), // - wild type
		strings.ToUpper("agcaaggTcgaggagctgttcCccggggtggtgcccatcctggtcgagctggacggcgac"),
		strings.ToUpper("agcaagggcgaggagctgttcaTcggggtggtgcccatcctggtcgagctggacggcgac"),
	}

	//query := msa[0]
	query_aa := "SKGEELFTGVVPILVELDGD"

	got := msaToAssemblyProblem(msa, query_aa)
	want := AssemblyProblem{
		Mutations: PointSet2D{Point2D{X: 7, Y: 2}, Point2D{X: 22, Y: 3}}, // 2 variants at 7, 3 variants at 22
		Seq:       "SKGEELFTGVVPILVELDGD",
	}

	if !reflect.DeepEqual(got, want) {
		fmt.Printf("Error: got %#v, want %#v\n", got, want)
	}

	// remove the wild type and recalculate
	msa = msa[1:]

	got = msaToAssemblyProblem(msa, query_aa)
	want = AssemblyProblem{
		Mutations: PointSet2D{Point2D{X: 7, Y: 2}, Point2D{X: 22, Y: 2}}, // still 2 variants at 7, now 2 variants at 22
		Seq:       "SKGEELFTGVVPILVELDGD",
	}

	if !reflect.DeepEqual(got, want) {
		fmt.Printf("Error: got %#v, want %#v\n", got, want)
	}

}

func TestMakeSplits(t *testing.T) {

	// split points seem to ve zero offset (see e.g. makeMember, MinDistToMut)
	// TODO: failing, off by one error (?)

	querySequence := "AGCAAGGGCGAGGAGCTGTTCACCGGGGTGGTGCCCATCCTGGTCGAGCTGGACGGCGAC"
	splitPoint := 5 // assume zero-offset
	endLength := 3
	endsToAvoid := []string{}

	got := makeSplits(querySequence, splitPoint, endLength, endsToAvoid)
	want := []string{"GGG"}

	fmt.Printf("%#v\n", got)

	if !reflect.DeepEqual(got, want) {
		fmt.Printf("Error: got %s, want %s\n", got, want)
	}

}

func TestGetSplitz(t *testing.T) {

	// returns a set of ends [][]string{} with one end per slice in inner level (?)

	querySequence := "AGCAAGGGCGAGGAGCTGTTCACCGGGGTGGTGCCCATCCTGGTCGAGCTGGACGGCGAC"
	split := PointSet1D{5, 12}
	endLength := 3
	endsToAvoid := []string{}

	got := getSplitz(split, querySequence, endLength, endsToAvoid)
	want := [][]string{
		{"GGG"}, {"GAG"},
	}

	if !reflect.DeepEqual(got, want) {
		fmt.Printf("Error: got %#v, want %#v\n", got, want)
	}

}

func TestGetEnds(t *testing.T) {

	// converts [][]string{} from getSplitz to []string{}

	querySequence := "AGCAAGGGCGAGGAGCTGTTCACCGGGGTGGTGCCCATCCTGGTCGAGCTGGACGGCGAC"
	split := PointSet1D{5, 12}
	endLength := 3
	endsToAvoid := []string{}

	got := getSplitz(split, querySequence, endLength, endsToAvoid)
	want := []string{"GGG", "GAG"}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Error: got %#v, want %#v\n", got, want)
	}

}

func TestEndsOKDefault(t *testing.T) {

	// if noTransitions == false, this routine simply checks for duplicates ends
	// the case noTransitions == true is tested separately

	noTransitions := false

	type testData struct {
		endSet []string
		want   bool
		name   string
	}

	for _, test := range []testData{
		{[]string{}, true, "empty"},
		{[]string{"AGG"}, true, "single"},
		{[]string{"AGG", "TTT"}, true, "two distinct"},
		{[]string{"AGG", "AGG"}, false, "duplicates"},
	} {
		ends := make([][]string, 0)
		for _, end := range test.endSet {
			ends = append(ends, []string{end})
		}
		result := make(map[string]bool)
		got := endsOK(ends, result, noTransitions)
		if got != test.want {
			t.Errorf("Error, %s: got %#v, want %#v\n", test.name, got, test.want)
		}
	}

}

func TestMinDistToMutation(t *testing.T) {

	// MinDistToMut assumes split points are same frame as mutations (zero-offset)

	mutations := PointSet2D{{40, 1}}

	type testData struct {
		i    int
		want int
		name string
	}

	for _, test := range []testData{
		{40, 0, "mutation zero distance"},
		{39, 1, "mutation distance 1 -"},
		{41, 1, "mutation distance 1 +"},
		{38, 2, "mutation distance 2 -"},
		{42, 2, "mutation distance 2 +"},
		{37, 3, "mutation distance 3 -"},
		{43, 3, "mutation distance 3 +"},
	} {
		got := mutations.MinDistTo(test.i)
		if got != test.want {
			t.Errorf("Error, %s: got %d, want %d\n", test.name, got, test.want)
		}
	}

}

func TestDist(t *testing.T) {

	// used to calculate segment lengths in function valid
	// TODO: fails, off by one error - calculates dist as b - a + 1

	got := dist(15, 20)
	want := 5

	if got != want {
		t.Errorf("Error: got %#v, want %#v\n", got, want)
	}

}

func TestValid(t *testing.T) {

	// calculates segment lengths using dist(a, b)
	// calulates distances to mutations using MinDistToMut()
	// TODO: all passing except the marginal MinLen case, which fails possibly due to off by one error in dist()

	assemblyProblem := AssemblyProblem{
		Mutations: PointSet2D{{40, 1}},
		Seq:       "SKGEELFTGVVPILVELDGD",
	}

	constraints := Constraints{
		MaxSplits:     10,
		MinLen:        5,
		MaxLen:        100,
		MinDistToMut:  2,
		Query:         "AGCAAGGGCGAGGAGCTGTTCACCGGGGTGGTGCCCATCCTGGTCGAGCTGGACGGCGAC",
		EndsToAvoid:   []string{},
		EndLen:        3,
		NoTransitions: false,
	}

	type testData struct {
		split PointSet1D
		want  bool
		name  string
	}

	for _, test := range []testData{
		{PointSet1D{}, false, "empty split"},
		{PointSet1D{4}, false, "1st segment length less than MinLen"},
		{PointSet1D{5}, true, "1st segment length equal MinLen, OK"},
		{PointSet1D{6}, true, "1st segment length exceeds MinLen, OK"},
		{PointSet1D{40}, false, "mutation zero distance"},
		{PointSet1D{41}, false, "mutation distance (1) too close to mutation"},
		{PointSet1D{42}, true, "mutation distance (2) equal MinDistToMut, OK"},
		{PointSet1D{43}, true, "mutation distance (3) exceeds MinDistToMut, OK"},
		// & generate a duplicate set through choice of split
	} {
		got := valid(test.split, assemblyProblem, constraints)
		if got != test.want {
			t.Errorf("Error, %s: got %#v, want %#v\n", test.name, got, test.want)
		}
	}

}

func TestPMCost(t *testing.T) {

	// TODO: fails due to possible off by one error in segment length calculation
	// or could redefine cost as length including ends
	// but nicer if it gives cost independent of split in the case of no mutations

	assemblyProblem := AssemblyProblem{
		Mutations: PointSet2D{{7, 2}, {22, 3}},
		Seq:       "SKGEELFTGVVPILVELDGD",
	}

	noMutations := AssemblyProblem{
		Mutations: PointSet2D{},
		Seq:       "SKGEELFTGVVPILVELDGD",
	}

	type testData struct {
		name    string
		split   PointSet1D
		want    int
		problem AssemblyProblem
	}

	/*
					> Split 1:
					Split [5, 10]
					Segment lengths are 5, 5, 50
					Cost is
					5 * 1 +       // no mutations, multiplier 1
					5 * 2 +       // 1 site with 2 mutations in [5, 10), multiplier 2
					50 * 3        // 1 site with 3 mutations in [10, 60), multipler 3
					= 157

					> Split 2:
					Split [5, 40]
					Segment lengths are 5, 35, 20
					Cost is
					5 * 1 +       // no mutations, multiplier 1
					35 * 2 * 3 +  // 1 site with 2 mutations, 1 site with 3 mutations in [5, 40), multiplier 6
					20 * 1        // no mutations, multiplier 1
					= 235

			                > Split 3:
			                Split []
			                Segment length 60
		                        6 * 60
			                = 360

			                > Split 4, no mutations:
			                Split []
			                Segment length 60
		                        60
			                = 60
	*/

	for _, test := range []testData{
		{"Split1", PointSet1D{5, 10}, 157, assemblyProblem},
		{"Split2", PointSet1D{5, 40}, 235, assemblyProblem},
		{"Split3", PointSet1D{}, 360, assemblyProblem},
		{"Split4", PointSet1D{}, 60, noMutations}, // sequence length
	} {
		got := Cost(test.split, test.problem)
		if got != test.want {
			t.Errorf("Error, %s: got %d, want %d\n", test.name, got, test.want)
		}
	}

}
