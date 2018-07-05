// ph.go Part of the Antha language
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

// Package for dealing with manipulation of buffers
package buffers

import (
	"fmt"
	"time"

	"github.com/antha-lang/antha/antha/anthalib/mixer"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
)

type PHperdegC float64

type PHMeasurement struct {
	Component       *wtype.Liquid
	Location        *wtype.Plate
	PHValue         float64
	Temp            wunit.Temperature
	TempCorrected   *float64
	RefTemp         *wunit.Temperature
	TempCoefficient *PHperdegC
	Adjusted        *float64
	Adjustedwith    *wtype.Liquid
}

type PH struct {
	PHValue float64
	Temp    wunit.Temperature
}

func (ph *PHMeasurement) TempCompensation(reftemp wunit.Temperature, tempcoefficientforsolution PHperdegC) (compensatedph float64) {

	ph.RefTemp = &reftemp //.SIValue()
	ph.TempCoefficient = &tempcoefficientforsolution

	tempdiff := ph.Temp.SIValue() - ph.RefTemp.SIValue()

	compensatedph = ph.PHValue + (float64(tempcoefficientforsolution) * tempdiff)
	ph.TempCorrected = &compensatedph
	return
}

// placeholder

/*func MeasurePH(*wtype.Liquid) (measurement float64) {
	return 7.0
}*/

func MeasurePH(*wtype.Liquid) (measured PHMeasurement) {
	measured = PHMeasurement{nil, nil, 0.0, wunit.NewTemperature(0.0, "C"), nil, nil, nil, nil, nil}
	return
}

// this should be performed on an LHComponent
// currently (wrongly) assumes only acid or base will be needed
func (ph *PHMeasurement) AdjustpH(ph_setpoint float64, ph_tolerance float64, ph_setPointTemp wunit.Temperature, Acid *wtype.Liquid, Base *wtype.Liquid) (adjustedsol wtype.Liquid, newph PHMeasurement, componentadded wtype.Liquid, err error) {

	pHmax := ph_setpoint + ph_tolerance
	pHmin := ph_setpoint - ph_tolerance

	//sammake([]wtype.Liquid,0)

	if ph.PHValue > pHmax {
		// calculate concentration of solution needed first, for now we'll add 10ul at a time until adjusted
		for {
			//newphmeasurement = ph
			acidsamp := mixer.Sample(Acid, wunit.NewVolume(10, "ul"))
			temporary := mixer.MixInto(ph.Location, "", ph.Component, acidsamp)
			time.Sleep(10 * time.Second)
			newphmeasurement := MeasurePH(temporary)
			if newphmeasurement.PHValue > pHmax {
				*ph = newphmeasurement
				//}
				//if {
				//ph.PH < pHmin
				//	continue
			} else {
				adjustedsol = *ph.Component
				newph = *ph
				componentadded = *Acid
				err = nil
				return
			}

		}
	}
	// basically just a series of sample, stir, wait and recheck pH

	if ph.PHValue < pHmin {
		for {
			//newphmeasurement = ph
			basesamp := mixer.Sample(Base, wunit.NewVolume(10, "ul"))
			temporary := mixer.MixInto(ph.Location, "", ph.Component, basesamp)
			time.Sleep(10 * time.Second)
			newphmeasurement := MeasurePH(temporary)
			if newphmeasurement.PHValue > pHmax {
				*ph = newphmeasurement
			} else {
				adjustedsol = *ph.Component
				newph = *ph
				componentadded = *Base
				err = nil
				return
			}

		}
	}
	//adjustedsol = ph.Component, newph = ph, componentadded = Acid,
	err = fmt.Errorf("Something went wrong here!")
	return
}
