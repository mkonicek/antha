package main

import (
	"fmt"
	"github.com/Synthace/go-glpk/glpk"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/Optimization"
)

func OptimizeAssembly2(problem AssemblyProblem, constraints Constraints, parameters Optimization.AssemblyOptimizerParameters) {
	splitz := MakeAllSplits(problem, constraints)
	best := 999999
	var bestSolution PointSet1D

	for _, split := range splitz {
		score, solution := SolveSplit(split, problem, constraints, parameters)

		if score < best {
			best = score
			bestSolution = solution
		}
	}

	fmt.Println("Best score: ", best)
	fmt.Println("Best solution: ", bestSolution)
}

func SolveSplit(split PointSet3D, problem AssemblyProblem, constraints Constraints, parameters Optimization.AssemblyOptimizerParameters) (score int, solution PointSet1D) {
	//
	//	optimization is set up as follows:
	//
	//		let:
	//			j_n   	be split points
	//			c_n	be the cost per base for each split point
	//			Pn-,Pn+ be the bounds on the split points
	//			L-,L+   be length bounds on the distance between split points
	//
	//	Minimise:
	//
	//			sum of c_n(j_n - j_n-1)
	//
	//	Subject to:
	//			Pn- <= j_n <=Pn+
	//			L- <= (j_n - j_n-1) <= L+
	//

	// setup

	lp := glpk.New()
	defer lp.Delete()

	lp.SetProbName("Fragments")
	lp.SetObjName("Z")

	lp.SetObjDir(glpk.MIN)

	// add columns

	lp.AddCols(len(split))
	cur := 1

	for _, pt := range split {
		lp.SetObjCoef(cur, float64(pt.Z))
		lp.SetColName(cur, fmt.Sprintf("split-%d", cur))
		lp.SetColKind(cur, glpk.IV)
		lp.SetColBnds(cur, glpk.DB, float64(pt.X), float64(pt.Y))
		cur += 1
	}

	// add row bounds etc

	cur = 1
	lp.AddRows(len(split) + 1)

	for i := 0; i < len(split)+1; i++ {
		cur += 1
		lp.SetRowBnds(cur, glpk.DB, float64(constraints.MinLen), float64(constraints.MaxLen))
	}

	panic("SET MAT ROWS")

	iocp := glpk.NewIocp()
	iocp.SetPresolve(true)
	//debug
	iocp.SetMsgLev(0)
	lp.Intopt(iocp)

	solution = make(PointSet1D, len(problem.Mutations))
	for i := 0; i < len(problem.Mutations); i++ {
		solution[i] = int(lp.MipColVal(i + 1))
	}

	cost := lp.ObjVal()

	return int(cost), solution
}

func MakeAllSplits(problem AssemblyProblem, constraints Constraints) []PointSet3D {
	splitRet := make([]PointSet3D, 0, 1)
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
			} else {
				// ??
			}
		}
	}

	return PointSet3D(ret)
}
