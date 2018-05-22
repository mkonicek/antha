// anthalib//liquidhandling/choose_stock_concentrations.go: Part of the Antha language
// Copyright (C) 2015 The Antha authors. All rights reserved.
//
// This program is free software; you can redistribute it and/or
// modify it under the terms of the GNU General Public License
// as published by the Free Software Foundation; either version 2
// of the License, or (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program; if not, write to the Free Software
// Foundation, Inc., 51 Franklin Street, Fifth Floor, Boston, MA  02110-1301, USA.
//
// For more information relating to the software or licensing issues please
// contact license@antha-lang.org or write to the Antha team c/o
// Synthace Ltd. The London Bioscience Innovation Centre
// 2 Royal College St, London NW1 0NH UK
package liquidhandling

/*
import (
	"fmt"

	"github.com/Synthace/go-glpk/glpk"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/antha/anthalib/wutil"
)

func choose_stock_concentrations(minrequired map[string]float64, maxrequired map[string]float64, Smax map[string]float64, vmin float64, T map[string]wunit.Volume) map[string]float64 {
	// we want to find the minimum concentrations
	// which fulfill the constraints

	// the optimization is set up as follows:
	//
	// 	let:
	//		Ck 	= 	concentration of component k		// we do not optimize this directly
	//		Gmk	=	min required concentration of k
	//		GMk	= 	max required concentration of k
	//		Tk	=	Minimum total final volume for k
	//		vmin	=	minimum channel capacity
	//		Skmax	=	Max solubility of component k
	//		Xk	=	GMk / Ck
	//
	//	Maximise	sum of -Xk
	//
	//	Subject to
	//			-Xk	<= -vmin GMk / Gmk Tk		(for each k)	-- min volume constraint
	//			-Xk	<= -GMk / Skmax			(for each k)	-- max conc constraint, this is set as a column constraint
	//			Sum Xk	<= 1
	//

	nc := len(minrequired)

	// no concentrations -> end here

	if nc == 0 {
		return (make(map[string]float64, 1))
	}

	// need to do these things in a consistent order
	names := make([]string, nc)

	cur := 0
	for name := range minrequired {
		names[cur] = name
		cur++
	}

	lp := glpk.New()
	defer lp.Delete()

	lp.SetProbName("Concentrations")
	lp.SetObjName("Z")
	lp.SetObjDir(glpk.MAX)

	// sets up number of constraints and B vector
	lp.AddRows(2*nc + 1)

	cur = 1

	for _, name := range names {
		lp.SetRowBnds(cur, glpk.UP, -999999.0, (-1.0*vmin*maxrequired[name])/(T[name].SIValue()*minrequired[name]))
		cur++
	}

	for _, name := range names {
		lp.SetRowBnds(cur, glpk.UP, -999999.0, (-1.0 * maxrequired[name] / Smax[name]))
		cur++
	}

	lp.SetRowBnds(cur, glpk.UP, 0.0, 1.0)

	// sets up objective and constraint coefficients

	lp.AddCols(nc)

	for idx := range names {
		lp.SetObjCoef(idx+1, -1.0)
		lp.SetColName(idx+1, fmt.Sprintf("X%d", cur))
		lp.SetColBnds(idx+1, glpk.LO, 0.0, 0.0)
	}

	cur = 1

	// constraint coeffs

	ind := wutil.Series(0, nc)

	for j := 0; j < 2; j++ {
		for i := 0; i < nc; i++ {
			row := make([]float64, nc+1)
			row[i+1] = -1.0
			lp.SetMatRow(cur, ind, row)
			cur++
		}
	}
	// now the sum constraint
	row := make([]float64, nc+1)
	for i := 0; i < nc; i++ {
		row[i+1] = 1.0
	}
	lp.SetMatRow(cur, ind, row)

	// solve it

	prm := glpk.NewSmcp()
	prm.SetMsgLev(0)
	err := lp.Simplex(prm)

	// now look at the solution

	concentrations := make(map[string]float64, nc)

	stat := lp.Status()

	if err != nil || stat != glpk.OPT {
		// some problem
		return concentrations
	}

	cur = 1
	for _, name := range names {
		concentrations[name] = maxrequired[name] / lp.ColPrim(cur)
		cur++
	}

	return concentrations
}

*/
