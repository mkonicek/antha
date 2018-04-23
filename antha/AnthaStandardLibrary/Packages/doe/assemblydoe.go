// multipleassemblies.go Part of the Antha language
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

package doe

import (
	"strconv"
	"strings"

	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/enzymes"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
)

func AssemblyParametersFromRuns(runs []Run, enzymename string) (assemblyparameters []enzymes.Assemblyparameters) {

	assemblyparameters = make([]enzymes.Assemblyparameters, 0)
	parts := make([]wtype.DNASequence, 0)
	var parameters enzymes.Assemblyparameters

	for j, run := range runs {

		parameters.Constructname = "run" + strconv.Itoa(j)
		parameters.Enzymename = enzymename

		for i := range run.Setpoints {

			if strings.Contains(run.Factordescriptors[i], "Vector") {

				parameters.Vector = run.Setpoints[i].(wtype.DNASequence)
			} else if i < len(runs)-1 {
				parts = append(parts, run.Setpoints[i].(wtype.DNASequence))
			}

		}
		parameters.Partsinorder = parts
		assemblyparameters = append(assemblyparameters, parameters)
	}

	return
}

func DNASequencetoInterface(genes []wtype.DNASequence) (geneseqs []interface{}) {

	vals := make([]interface{}, len(genes))
	for i, v := range genes {
		vals[i] = v
	}
	geneseqs = append(geneseqs, vals...)
	return
}
