// Part of the Antha language
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

package plot

import (
	"github.com/sajari/regression"
)

// RSquared computes R^2
func RSquared(xname string, xvalues []float64, yname string, yvalues []float64) (rsquared float64, variance float64, formula string) {
	var r regression.Regression

	r.SetObserved(yname)
	r.SetVar(0, xname)

	for i := range xvalues {
		r.Train(regression.DataPoint(yvalues[i], []float64{xvalues[i]}))
	}
	r.Run() // nolint

	rsquared = r.R2
	variance = r.Varianceobserved
	formula = r.Formula
	return
}
