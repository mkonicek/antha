// /anthalib/driver/liquidhandling/compositerobotinstruction.go: Part of the Antha language
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
	"github.com/antha-lang/antha/antha/anthalib/wutil"
	"github.com/antha-lang/antha/microArch/logger"
)

//autoGenerateLLModel attempt to generate a liquidlevel model and add it to the plate
func autoGenerateLLModel(well *wtype.LHWell) {

	//we really don't know very much about the well geometry, just go with a linear model
	//this might be OK for straight sided wells, but is likely to require a constant offset
	area, err := well.CalculateMaxCrossSectionArea()
	if err != nil {
		return
	}

	model := wutil.Quadratic{B: area.ConvertToString("mm^2")}
	well.SetLiquidLevelModel(model)

	logger.Info(fmt.Sprintf("Auto-generated Liquid Level Model (A, B = 0.0, %f) for plate type \"%s\"", area.ConvertToString("mm^2")))
}
