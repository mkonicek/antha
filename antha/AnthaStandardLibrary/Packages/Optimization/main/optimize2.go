package main

import (
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/Optimization"
)

type Split PointSet3D

func OptimizeAssembly2(problem AssemblyProblem, constraints Constraints, parameters Optimization.AssemblyOptimizerParameters) {

}

func MakeAllSplits(problem AssemblyProblem, constraints Constraints) []Split {
	splitRet := make([]Split, 0, 1)
	for i := 0; i < constraints.MaxSplits; i++ {
		splitz := MakeSplits(i, problem, constraints)
		splitRet = append(splitRet, splitz...)
	}

	return splitRet
}

func comb(n, r int) []string {
	if r == 0 {
		s := ""
		for i := 0; i < n; i++ {
			s += "0"
		}

		return []string{s}
	}

	ret := make([]string, 0, n)

	for i := 0; i < n; i++ {
		s := ""
		for j := 0; j < i; j++ {
			s += "0"
		}
		s += "1"

		c := comb(n-i-1, r-1)

		for _, v := range c {
			ret = append(ret, s+v)
		}
	}

	return ret

}

func MakeSplits(n int, problem AssemblyProblem, constraints Constraints) []Split {
	if n > len(problem.Mutations)+1 {
		return []Split{}
	}

	c := comb(len(problem.Mutations)+1, n)
	r := make([]Split, 0, len(c))

	for _, cc := range c {
		s := MakeSplit(cc, problem, constraints)

		if len(s) != 0 {
			r = append(r, s)
		}
	}

	return r
}

func MakeSplit(comb string, problem AssemblyProblem, constraints Constraints) Split {
	// each entry in PointSet3D says what the upper and lower bounds on the split point
	// are and what the mutational cost is

	ret := make(PointSet3D, 0, len(problem.Mutations))

	// move along, collecting the things

	s := 0
	m := 1

	for i := 0; i < len(comb); i++ {
		if comb[i] == '1' {
			var x, t int
			if i == len(problem.Mutations) {
				x = problem.Len
				t = problem.Len // not relevant
			} else {
				x = problem.Mutations[i].X - constraints.MinDistToMut
				t = problem.Mutations[i].X + constraints.MinDistToMut
			}
			ret = append(ret, Point3D{X: s, Y: x, Z: m})
			s = t

			m = 1
		} else {
			if i < len(problem.Mutations) {
				m *= problem.Mutations[i].Y
				s = problem.Mutations[i].X + constraints.MinDistToMut
			} else {
				// ??
			}
		}
	}

	return Split(ret)
}
