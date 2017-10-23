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

// Package for dealing with manipulation of PH measurements
package ph

import (
	"encoding/json"
	"fmt"

	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
)

var (
	errPHNotFound = fmt.Errorf("No PH found")
)

// AddSetPoint adds a PH set point to an LHComponent.
func AddSetPoint(component *wtype.LHComponent, setpoint PH) error {

	ph, err := getPH(component)

	if err != nil {
		if err.Error() != errPHNotFound.Error() {
			return err
		}
	}

	ph.SetPoint = setpoint

	component, err = setPH(component, ph)

	return err
}

// SetPoint looks up a PH set point for an LHComponent.
func SetPoint(component *wtype.LHComponent) (setpoint PH, err error) {

	ph, err := getPH(component)

	if err != nil {
		return setpoint, err
	}

	setpoint = ph.SetPoint

	return setpoint, nil
}

// PHperdegC is a coeffecient for the change in pH resulting per degree C increase in temperature
type PHperdegC float64

// PHMeasurement stores a measurement of PH on an LHComponent
type PHMeasurement struct {
	// SetPoint stores the desired PH.
	// It is possible to store a SetPoint without the other fields.
	SetPoint PH

	// The measured pH value
	MeasuredPHValue PH

	// The measured coefficient for a solution for the change in pH with temperature for that solution.
	TemperatureCoefficient *PHperdegC

	// If the PH value has been adjusted the volumes of samples used to adjust the pH are stored here.
	AdjustedWith []*wtype.LHComponent
}

// PH stores a pH Value, Precision and temperature for the the PH value.
type PH struct {
	Value     float64           // Value. e.g. 7.2
	Precision float64           // +/- this value. e.g. 0.2
	Temp      wunit.Temperature // temperature at which pH value corresponds.
}

// ToString returns the PH value as a string.
func (ph PH) ToString() string {
	return fmt.Sprintf("PH %.2f +/- %.2f at %s", ph.Value, ph.Precision, ph.Temp.ToString())
}

// utility function to allow the object properties to be retained when serialised.
func serialise(measurement PHMeasurement) ([]byte, error) {

	return json.Marshal(measurement)
}

// utility function to allow the object properties to be retained when serialised.
func deserialise(data []byte) (measurement PHMeasurement, err error) {
	measurement = PHMeasurement{}
	err = json.Unmarshal(data, &measurement)
	return
}

// getPH returns a PH measurement from a component.
func getPH(comp *wtype.LHComponent) (measurement PHMeasurement, err error) {

	pH, found := comp.Extra["PH"]

	if !found {
		return measurement, errPHNotFound
	}

	var bts []byte

	bts, err = json.Marshal(pH)
	if err != nil {
		return
	}

	err = json.Unmarshal(bts, &measurement)

	if err != nil {
		err = fmt.Errorf("Problem getting %s PH measurement. Error: %s", comp.Name(), err.Error())
	}

	return
}

// Add a PH measurement to a component.
// Any existing measurement will be overwritten.
// Users should use add AddPH function.
func setPH(comp *wtype.LHComponent, ph PHMeasurement) (*wtype.LHComponent, error) {

	comp.Extra["PH"] = ph

	return comp, nil
}

/*
func (ph *PHMeasurement) TempCompensation(reftemp wunit.Temperature, tempcoefficientforsolution PHperdegC) (compensatedph float64) {

	ph.RefTemp = &reftemp //.SIValue()
	ph.TemperatureCoefficient = &tempcoefficientforsolution

	tempdiff := ph.Temp.SIValue() - ph.RefTemp.SIValue()

	compensatedph = ph.PHValue + (float64(tempcoefficientforsolution) * tempdiff)
	ph.TempCorrected = &compensatedph
	return
}
*/
// placeholder

/*func MeasurePH(*wtype.LHComponent) (measurement float64) {
	return 7.0
}*/

/*
// this should be performed on an LHComponent
// currently (wrongly) assumes only acid or base will be needed
func (ph *PHMeasurement) AdjustpH(sample *wtype.LHComponent, ph_setpoint float64, ph_tolerance float64, ph_setPointTemp wunit.Temperature, Acid *wtype.LHComponent, Base *wtype.LHComponent) (adjustedsol *wtype.LHComponent, newph PHMeasurement, componentadded *wtype.LHComponent, err error) {

	pHmax := ph_setpoint + ph_tolerance
	pHmin := ph_setpoint - ph_tolerance

	//sammake([]wtype.LHComponent,0)

	if ph.PHValue > pHmax {
		// calculate concentration of solution needed first, for now we'll add 10ul at a time until adjusted
		for {
			//newphmeasurement = ph
			acidsamp := mixer.Sample(Acid, wunit.NewVolume(10, "ul"))
			temporary := mixer.Mix(sample, acidsamp)
			time.Sleep(10 * time.Second)
			newphmeasurement := MeasurePH(temporary)
			if newphmeasurement.PHValue > pHmax {
				*ph = newphmeasurement
				//}
				//if {
				//ph.PH < pHmin
				//	continue
			} else {
				adjustedsol = sample
				newph = *ph
				componentadded = Acid
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
*/
