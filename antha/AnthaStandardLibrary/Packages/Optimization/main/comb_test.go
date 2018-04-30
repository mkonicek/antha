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

func MakeSplits(n int, problem AssemblyProblem, constraints Constraints) []PointSet3D {
	if n > len(problem.Mutations)+1 {
		return []PointSet3D{}
	}

	c := comb(len(problem.Mutations)+1, n)
	r := make([]PointSet3D, 0, len(c))

	for _, cc := range c {
		s := MakeSplit(cc, problem, constraints)

		if len(s) != 0 {
			r = append(r, s)
		}
	}

	return r
}

func MakeSplit(comb string, problem AssemblyProblem, constraints Constraints) PointSet3D {
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
			}
		}
	}

	return PointSet3D(ret)
}

type Point3D struct {
	X int
	Y int
	Z int
}

type PointSet3D []Point3D

func (ps3d PointSet3D) Less(i, j int) bool {
	return ps3d[i].X < ps3d[j].X
}

func (ps3d PointSet3D) Swap(i, j int) {
	t := ps3d[i]
	ps3d[i] = ps3d[j]
	ps3d[j] = t
}

func (ps3d PointSet3D) Len() int {
	return len(ps3d)
}
