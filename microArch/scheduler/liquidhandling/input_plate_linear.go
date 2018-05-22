// /anthalib/liquidhandling/input_plate_linear.go: Part of the Antha language
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

import (
	"fmt"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"gonum.org/v1/gonum/mat"
	"gonum.org/v1/gonum/optimize/convex/lp"
	"math"
)

func choose_plate_assignments(component_volumes map[string]wunit.Volume, plate_types []*wtype.LHPlate, weight_constraint map[string]float64) (map[string]map[*wtype.LHPlate]int, error) {
	// 	v2.0: modified to use gonum/optimize/convex/lp
	//
	//	optimization is set up as follows:
	//
	//		let:
	//			Xk 	= 	Number of wells of type Y containing component Z (k = 1...YZ)
	//			Vy	= 	Working volume of well type Y
	//			RVy	= 	Residual volume of well type Y
	//			TVz	= 	Total volume of component Z required
	//			WRy	=	Rate of wells of type y in their plate
	//			PMax	=	Maximum number of plates
	//			WMax	= 	Maximum number of wells
	//
	//	Minimise:
	//			sum of Xk WRy RVy
	//
	//	Subject to:
	//			sum of Xk Vy 	>= TVz	for each component Z
	//				- which we express as -XkVy <= -TVz
	//
	//			sum of WRy Xk 	<= PMax
	//			sum of Xk	<= WMax
	//
	//

	// defense
	//

	ppt := make([]*wtype.LHPlate, 0, len(plate_types))
	h := make(map[string]bool, len(plate_types))

	fmt.Println("Autoallocate plates available:")
	for _, p := range plate_types {
		fmt.Println(p.Type)
	}

	for _, p := range plate_types {
		if h[p.Type] {
			continue
		}
		ppt = append(ppt, p)
		h[p.Type] = true
	}

	plate_types = ppt

	assignments := make(map[string]map[*wtype.LHPlate]int, len(component_volumes))

	// func Simplex(c []float64, A mat.Matrix, b []float64, tol float64, initialBasic []int) (optF float64, optX []float64, err error)

	n_cols := len(component_volumes) * len(plate_types)
	n_rows := len(component_volumes) + 2

	constraintMatrix := make([]float64, n_cols*n_rows)
	constraintBounds := make([]float64, n_rows)
	objectiveCoefs := make([]float64, n_cols)

	component_order := make([]string, len(component_volumes))
	cur := 0

	for cmp, vol := range component_volumes {
		component_order[cur] = cmp
		v := vol.ConvertTo(wunit.ParsePrefixedUnit("ul"))
		constraintBounds[cur] = -1.0 * v

		for pindex, plate := range plate_types {
			// set up objective coefficient, column name and lower bound
			rv := plate.Welltype.ResidualVolume()
			coef := rv.ConvertTo(wunit.ParsePrefixedUnit("ul"))*float64(weight_constraint["RESIDUAL_VOLUME_WEIGHT"]) + 1.0
			//objectiveCoefs = append(objectiveCoefs, coef)
			objectiveCoefs[cur*len(plate_types)+pindex] = coef
		}

		cur += 1
	}

	cur = 0

	// make constraint rows
	for row := range component_order {
		setRowFor(constraintMatrix, row, n_cols, component_order, plate_types)
		cur += 1
	}

	// plate constraint
	constraintBounds[cur] = weight_constraint["MAX_N_PLATES"] - 1.0

	for i := range component_order {
		for j := 0; j < len(plate_types); j++ {
			// the coefficient here is 1/the number of this well type per plate
			coef := 1.0 / float64(plate_types[j].Nwells)
			//constraintMatrix[cur+(i*len(plate_types)+j)] = coef

			constraintMatrix[cur*n_cols+(i*len(plate_types))+j] = coef
		}
	}

	cur += 1

	// well constraint
	constraintBounds[cur] = weight_constraint["MAX_N_WELLS"]

	// for the matrix we just add a row of 1s
	for i := 0; i < n_cols; i++ {
		constraintMatrix[cur*n_cols+i] = 1.0
	}

	matUBConstraintMatrix := mat.NewDense(n_rows, n_cols, constraintMatrix)

	//	cNew, aNew, bNew := lp.Convert(objectiveCoefs, matConstraintMatrix, constraintBounds, aOld, bOld)
	cNew, aNew, bNew := lp.Convert(objectiveCoefs, matUBConstraintMatrix, constraintBounds, nil, nil)

	tolerance := 1e-6
	_, optX, err := lp.Simplex(cNew, aNew, bNew, tolerance, nil)

	if err != nil {
		return nil, err
	}

	// now create the assignment outputs

	cur = 0
	for _, c := range component_order {
		pmap := make(map[*wtype.LHPlate]int)
		for _, p := range plate_types {
			if optX[cur] > tolerance {
				pmap[p.Dup()] = int(math.Ceil(optX[cur]))
			}
			cur += 1
		}
		assignments[c] = pmap
	}

	return assignments, nil
}

func setRowFor(mtx []float64, rowN, n_cols int, component_order []string, plate_types []*wtype.LHPlate) {
	row := make([]float64, n_cols)

	for j := range plate_types {
		// pick out a set of columns according to which row we're on
		// volume constraints are the working volumes of the wells
		row[rowN*len(plate_types)+j] = -1.0 * plate_types[j].Welltype.MaxWorkingVolume().ConvertTo(wunit.ParsePrefixedUnit("ul"))
	}

	for c, v := range row {
		mtx[rowN*n_cols+c] = v
	}
}
