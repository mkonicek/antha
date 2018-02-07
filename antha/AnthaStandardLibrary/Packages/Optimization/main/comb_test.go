package main

import (
	"fmt"
	"testing"
)

func TestComb(t *testing.T) {
	fmt.Println(comb(5, 2))
}

func TestSplit(t *testing.T) {
	// n int, problem AssemblyProblem, constraints Constraints

	problem := BasicProblem()
	constraints := DefaultConstraints()

	split := MakeSplits(2, problem, constraints)

	fmt.Println(split)
}
