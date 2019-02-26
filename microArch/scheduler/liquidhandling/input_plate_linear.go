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

func choosePlateAssignments(component_volumes map[string]wunit.Volume, plate_types []*wtype.Plate, weight_constraint map[string]float64) (map[string]map[*wtype.Plate]int, error) {

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
	//			Brv	= 	Residual volume weight
	//
	//
	//	Minimise:
	//			sum of Xk WRy (Brv RVy + 1)
	//
	//	Subject to:
	//			sum of Xk Vy 	>= TVz	for each component Z
	//				- which we express as -XkVy <= -TVz
	//
	//			sum of Xk	<= WMax
	//
	//

	// weight_constraint defines
	// 	RESIDUAL_VOLUME_WEIGHT	- Brv above
	//	MAX_N_WELLS		- WMax above
	// 	both are strings mapping to floats

	// defense
	//

	ppt := make([]*wtype.Plate, 0, len(plate_types))
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

	assignments := make(map[string]map[*wtype.Plate]int, len(component_volumes))

	n_cols := len(plate_types) * len(component_volumes)
	n_rows := len(component_volumes)
	n_constraint_rows := 1

	constraintMatrixG := make([]float64, n_cols*n_constraint_rows)
	constraintMatrixA := make([]float64, n_cols*n_rows)
	constraintBoundsH := make([]float64, n_constraint_rows)
	constraintBoundsB := make([]float64, n_rows)
	objectiveCoefs := make([]float64, n_cols)

	component_order := make([]string, len(component_volumes))
	cur := 0

	for cmp, vol := range component_volumes {
		component_order[cur] = cmp
		//v := vol.ConvertTo(wunit.ParsePrefixedUnit("ul"))
		v := vol.MustInStringUnit("ul").RawValue()
		//constraintBoundsB[cur] = -1.0 * v
		constraintBoundsB[cur] = v
		for pindex, plate := range plate_types {
			// set up objective coefficient, column name and lower bound
			rVol := plate.Welltype.ResidualVolume()
			rv := rVol.MustInStringUnit("ul").RawValue()
			rv = math.Sqrt(rv) // works well in ILP case
			coef := rv*float64(weight_constraint["RESIDUAL_VOLUME_WEIGHT"]) + 1.0
			//objectiveCoefs = append(objectiveCoefs, coef)
			objectiveCoefs[cur*len(plate_types)+pindex] = coef
		}

		cur += 1
	}

	cur = 0

	// make constraint rows
	for row := range component_order {
		setRowFor(constraintMatrixA, row, n_cols, plate_types)
		cur += 1
	}

	// well constraint
	constraintBoundsH[0] = 1.0 * weight_constraint["MAX_N_WELLS"]

	// for the matrix we just add a row of 1s
	for i := 0; i < n_cols; i++ {
		constraintMatrixG[i] = 1.0
	}

	matConstraintMatrixG := mat.NewDense(n_constraint_rows, n_cols, constraintMatrixG)
	matConstraintMatrixA := mat.NewDense(n_rows, n_cols, constraintMatrixA)

	//	cNew, aNew, bNew := lp.Convert(objectiveCoefs, matConstraintMatrix, constraintBounds, aOld, bOld)
	cNew, aNew, bNew := lp.Convert(objectiveCoefs, matConstraintMatrixG, constraintBoundsH, matConstraintMatrixA, constraintBoundsB)

	tolerance := 1e-10
	optF, optX, err := lp.Simplex(cNew, aNew, bNew, tolerance, nil)

	if err != nil || optF < 1e-10 {
		for i := 0; i < 20; i++ {
			tolerance *= 10.0
			optF, optX, err = lp.Simplex(cNew, aNew, bNew, tolerance, nil)

			if err == nil && optF > 1e-10 {
				break
			}
		}
	}

	if err != nil {
		return nil, err
	}

	// now create the assignment outputs

	cur = 0
	for _, c := range component_order {
		pmap := make(map[*wtype.Plate]int)
		for _, p := range plate_types {
			if optX[cur] > 0.0 {
				pmap[p.Dup()] = int(math.Ceil(optX[cur]))
			}
			cur += 1
		}
		assignments[c] = pmap
	}

	return assignments, nil
}

// given matrix mtx, set row rowN to the working volume of each plate in plate types
func setRowFor(mtx []float64, rowN, n_cols int, plate_types []*wtype.Plate) {
	row := make([]float64, n_cols)

	for j := range plate_types {
		// pick out a set of columns according to which row we're on
		// volume constraints are the working volumes of the wells
		row[rowN*len(plate_types)+j] = plate_types[j].Welltype.MaxWorkingVolume().MustInStringUnit("ul").RawValue()
	}

	for c, v := range row {
		mtx[rowN*n_cols+c] = v
	}
}
