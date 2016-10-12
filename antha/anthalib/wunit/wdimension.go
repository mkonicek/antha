// wunit/wdimension.go: Part of the Antha language
// Copyright (C) 2014 the Antha authors. All rights reserved.
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

package wunit

import (
	"fmt"
	"strings"
	"time"

	"github.com/antha-lang/antha/microArch/logger"
)

// length
type Length struct {
	*ConcreteMeasurement
}

func EZLength(v float64) Length {
	return NewLength(v, "m")
}

// make a length
func NewLength(v float64, unit string) Length {
	l := Length{NewPMeasurement(v, unit)}

	// check

	if l.Unit().RawSymbol() != "m" {
		panic("Base unit for lengths must be meters")
	}

	return l
}

// area
type Area struct {
	*ConcreteMeasurement
}

// make an area unit
func NewArea(v float64, unit string) (a Area) {
	if unit == "m^2" {
		a = Area{NewMeasurement(v, "", unit)}
	} else if unit == "mm^2" {
		//a = Area{NewPMeasurement(v /**0.000001*/, unit)}
		a = Area{NewMeasurement(v, "", unit)}
		// should be OK
	} else {
		panic("Can't make areas which aren't square (milli)metres")
	}

	return
}

// volume -- strictly speaking of course this is length^3
type Volume struct {
	*ConcreteMeasurement
}

// make a volume
func NewVolume(v float64, unit string) (o Volume) {
	if len(strings.TrimSpace(unit)) == 0 {
		return ZeroVolume()
	}

	/*if len(strings.TrimSpace(unit)) == 0 {
		panic("Can't make Volumes without unit")
	}

	if len(strings.TrimSpace(unit)) == 1 {
		o = Volume{NewMeasurement(v, "", unit)}
	}
	*/

	o = Volume{NewPMeasurement(v, unit)}

	return
}

func CopyVolume(v Volume) Volume {
	ret := NewVolume(v.RawValue(), v.Unit().PrefixedSymbol())
	return ret
}

// Add volumes
func AddVolumes(vols []Volume) (newvolume Volume) {
	// ideally should check these have the same Dimension
	// need to improve this
	var tempvol Volume
	tempvol = NewVolume(0.0, "ul")
	for _, vol := range vols {
		if tempvol.Unit().PrefixedSymbol() == vol.Unit().PrefixedSymbol() {
			tempvol = NewVolume(newvolume.RawValue()+vol.RawValue(), newvolume.Unit().PrefixedSymbol())
			newvolume = tempvol
		} else {
			tempvol = NewVolume(tempvol.SIValue()+vol.SIValue(), newvolume.Unit().BaseSISymbol())
		}
	}
	return

}

// subtract volumes
func SubtractVolumes(OriginalVol Volume, subtractvols []Volume) (newvolume Volume) {
	// ideally should check these have the same Dimension
	// need to improve this

	tempvol := (CopyVolume(OriginalVol))

	fmt.Println(tempvol)

	for _, vol := range subtractvols {
		newvolume = NewVolume(tempvol.RawValue()-vol.RawValue(), tempvol.Unit().PrefixedSymbol())
		tempvol = (CopyVolume(newvolume))
	}
	return

}

// multiply volume
func MultiplyVolume(v Volume, factor float64) (newvolume Volume) {
	// ideally should check these have the same Dimension
	// need to improve this
	newvolume = NewVolume(v.RawValue()*float64(factor), v.Unit().PrefixedSymbol())
	return

}

// divide volume
func DivideVolume(v Volume, factor float64) (newvolume Volume) {
	// ideally should check these have the same Dimension
	// need to improve this
	newvolume = NewVolume(v.RawValue()/float64(factor), v.Unit().PrefixedSymbol())
	return

}

// multiply volume
func MultiplyConcentration(v Concentration, factor float64) (newconc Concentration) {
	// ideally should check these have the same Dimension
	// need to improve this
	newconc = NewConcentration(v.RawValue()*float64(factor), v.Unit().PrefixedSymbol())
	return

}

// divide volume
func DivideConcentration(v Concentration, factor float64) (newconc Concentration) {
	// ideally should check these have the same Dimension
	// need to improve this
	newconc = NewConcentration(v.RawValue()/float64(factor), v.Unit().PrefixedSymbol())
	return

}

func (v Volume) Dup() Volume {
	ret := NewVolume(v.RawValue(), v.Unit().PrefixedSymbol())
	return ret
}

func ZeroVolume() Volume {
	return NewVolume(0.0, "ul")
}

// temperature
type Temperature struct {
	*ConcreteMeasurement
}

// make a temperature
func NewTemperature(v float64, unit string) Temperature {
	if unit != "˚C" && // RING ABOVE, LATIN CAPITAL LETTER C
		unit != "C" && // LATIN CAPITAL LETTER C
		unit != "℃" && // DEGREE CELSIUS
		unit != "°C" { // DEGREE, LATIN CAPITAL LETTER C
		panic("Can't make temperatures which aren't in degrees C")
	}
	t := Temperature{NewMeasurement(v, "", "℃")}
	return t
}

// time
type Time struct {
	*ConcreteMeasurement
}

// make a time unit
func NewTime(v float64, unit string) (t Time) {

	approvedunits := []string{"days", "h", "min", "s", "ms"}

	var approved bool
	for i := range approvedunits {

		if unit == approvedunits[i] {
			approved = true
			break
		}
	}

	if !approved {
		panic("Can't make Time with non approved unit of " + unit + ". Approved units are: " + strings.Join(approvedunits, ", "))
	}
	if unit == "s" {
		t = Time{NewMeasurement(v, "", unit)}
	} else if unit == "ms" {
		t = Time{NewMeasurement(v/1000, "", "s")}
	} else if unit == "min" {
		t = Time{NewMeasurement(v*60, "", "s")}
	} else if unit == "h" {
		t = Time{NewMeasurement(v*3600, "", "s")}
	}
	return t
}

func (t Time) Seconds() float64 {
	return t.SIValue()
}

func (t Time) AsDuration() time.Duration {
	// simply use the parser

	d, e := time.ParseDuration(t.ToString())

	if e != nil {
		logger.Fatal(e.Error())
	}

	return d
}

func FromDuration(t time.Duration) Time {
	return NewTime(float64(t.Seconds()), "s")
}

// mass
type Mass struct {
	*ConcreteMeasurement
}

// make a mass unit

func NewMass(v float64, unit string) (o Mass) {

	approvedunits := map[string]Unit{
		"ng": Unit{Base: "g", Prefix: "n", Multiplier: 1.0},
		"ug": Unit{Base: "g", Prefix: "u", Multiplier: 1.0},
		"mg": Unit{Base: "g", Prefix: "m", Multiplier: 1.0},
		"g":  Unit{Base: "g", Prefix: "", Multiplier: 1.0},
		"kg": Unit{Base: "g", Prefix: "k", Multiplier: 1.0},
	}

	var approved bool
	for key, _ := range approvedunits {

		if unit == key {
			approved = true
			break
		}
	}

	if !approved {
		panic("Can't make masses with non approved unit of " + unit + ". Approved units are: " + fmt.Sprint(approvedunits))
	}

	unitdetails := approvedunits[unit]

	o = Mass{NewMeasurement((v * unitdetails.Multiplier), unitdetails.Prefix, unitdetails.Base)}
	fmt.Println("NewMass: v: ", v, "unit: ", unit, "mass.ToString():", o.ToString(), "mass.SIValue()", o.SIValue())
	return
	/*
		if len(strings.TrimSpace(unit)) == 0 {
			panic("Can't make masses without unit")
		}
		if len(strings.TrimSpace(unit)) == 1 {
			o = Mass{NewMeasurement(v, "", unit)}
		}
		if len(strings.TrimSpace(unit)) > 1 {

			o = Mass{NewPMeasurement(v, unit)}
		}
	*/
	return //Mass{NewPMeasurement(v, unit)}
}

// defines mass to be a SubstanceQuantity
func (m *Mass) Quantity() Measurement {
	return m
}

// mole
type Moles struct {
	*ConcreteMeasurement
}

// generate a new Amount in moles
func NewAmount(v float64, unit string) Moles {
	if unit != "M" {
		panic("Can't make amounts which aren't in moles")
	}

	m := Moles{NewMeasurement(v, "", unit)}
	return m
}

// defines Amount to be a SubstanceQuantity
func (a *Moles) Quantity() Measurement {
	return a
}

// angle
type Angle struct {
	*ConcreteMeasurement
}

// generate a new angle unit
func NewAngle(v float64, unit string) Angle {
	if unit != "radians" {
		panic("Can't make angles which aren't in radians")
	}

	a := Angle{NewMeasurement(v, "", unit)}
	return a
}

// angular velocity (one way or another)

type AngularVelocity struct {
	*ConcreteMeasurement
}

func NewAngularVelocity(v float64, unit string) AngularVelocity {
	if unit != "rpm" {
		panic("Can't make angular velicities which aren't in rpm")
	}

	r := AngularVelocity{NewMeasurement(v, "", unit)}
	return r
}

// this is really Mass Length/Time^2
type Energy struct {
	*ConcreteMeasurement
}

// make a new energy unit
func NewEnergy(v float64, unit string) Energy {
	if unit != "J" {
		panic("Can't make energies which aren't in Joules")
	}

	e := Energy{NewMeasurement(v, "", unit)}
	return e
}

// a Force
type Force struct {
	*ConcreteMeasurement
}

// a new force in Newtons
func NewForce(v float64, unit string) Force {
	if unit != "N" {
		panic("Can't make forces which aren't in Newtons")
	}

	f := Force{NewMeasurement(v, "", unit)}
	return f
}

// a Pressure structure
type Pressure struct {
	*ConcreteMeasurement
}

// make a new pressure in Pascals
func NewPressure(v float64, unit string) Pressure {
	if unit != "Pa" {
		panic("Can't make pressures which aren't in Pascals")
	}

	p := Pressure{NewMeasurement(v, "", unit)}

	return p
}

// defines a concentration unit
type Concentration struct {
	*ConcreteMeasurement
	//MolecularWeight *float64
}

type Unit struct {
	Base       string
	Prefix     string
	Multiplier float64
}

// make a new concentration in SI units... either M/l or kg/l
func NewConcentration(v float64, unit string) (o Concentration) {

	/*
		if unit == "mg/ml" {
			unit = "g/l"
		} else if unit == "ng/ul" {
			unit = "mg/l"
		}

		if len(strings.TrimSpace(unit)) == 0 {
			panic("Can't make concentration without unit")
		}
		if len(strings.TrimSpace(unit)) == 3 {
			o = Concentration{NewMeasurement(v, "", unit)}
		}
		if len(strings.TrimSpace(unit)) > 3 {

			o = Concentration{NewPMeasurement(v, unit)}
		}
	*/

	approvedunits := map[string]Unit{
		"mg/ml":  Unit{Base: "g/l", Prefix: "", Multiplier: 1.0},
		"g/L":    Unit{Base: "g/l", Prefix: "", Multiplier: 1.0},
		"kg/l":   Unit{Base: "g/l", Prefix: "", Multiplier: 0.001},
		"kg/L":   Unit{Base: "g/l", Prefix: "", Multiplier: 0.001},
		"g/l":    Unit{Base: "g/l", Prefix: "", Multiplier: 1.0},
		"mg/L":   Unit{Base: "g/l", Prefix: "m", Multiplier: 1.0},
		"mg/l":   Unit{Base: "g/l", Prefix: "m", Multiplier: 1.0},
		"ng/ul":  Unit{Base: "g/l", Prefix: "m", Multiplier: 1.0},
		"Mol/L":  Unit{Base: "Mol/l", Prefix: "", Multiplier: 1.0},
		"M":      Unit{Base: "Mol/l", Prefix: "", Multiplier: 1.0},
		"M/l":    Unit{Base: "Mol/l", Prefix: "", Multiplier: 1.0},
		"mM":     Unit{Base: "Mol/l", Prefix: "m", Multiplier: 1.0},
		"mMol/L": Unit{Base: "Mol/l", Prefix: "m", Multiplier: 1.0},
	}

	var approved bool
	for key, _ := range approvedunits {

		if unit == key {
			approved = true
			break
		}
	}

	if !approved {
		panic("Can't make Concentration with non approved unit of " + unit + ". Approved units are: " + fmt.Sprint(approvedunits))
	}
	/*if unit == "g/L" {
		o = Concentration{NewMeasurement(v, "", "g/l")}
		fmt.Println(v, unit, o.ToString())
		return
	} else if unit == "mg/ml" {
		o = Concentration{NewMeasurement(v, "", "g/l")}
		fmt.Println(v, unit, o.ToString())
		return
	} else if unit == "mg/L" {
		o = Concentration{NewMeasurement(v, "m", "g/l")}
		fmt.Println(v, unit, o.ToString())
		return
	} else if unit == "ng/ul" {
		o = Concentration{NewMeasurement(v, "m", "g/l")}
		fmt.Println(v, unit, o.ToString())
		return
	}*/

	unitdetails := approvedunits[unit]

	o = Concentration{NewMeasurement((v * unitdetails.Multiplier), unitdetails.Prefix, unitdetails.Base)}
	fmt.Println("NewConcentration: v: ", v, "unit: ", unit, "conc.ToString():", o.ToString(), "conc.SIValue()", o.SIValue())
	return

	/*	if len(strings.TrimSpace(unit)) == 0 {
			panic("Can't make concentration without unit")
		}
		if len(strings.TrimSpace(unit)) == 3 {
			o = Concentration{NewMeasurement(v, "", unit)}
		}
		if len(strings.TrimSpace(unit)) > 3 {

			o = Concentration{NewPMeasurement(v, unit)}
		}*/

	return //Mass{NewPMeasurement(v, unit)}
}

// mass or mole
type SubstanceQuantity interface {
	Quantity() Measurement
}

func (conc *Concentration) GramPerL(molecularweight float64) (conc_g Concentration) {
	if conc.Munit.BaseSISymbol() == "g/l" {
		conc_g = *conc
	}

	if conc.Munit.BaseSISymbol() == "kg/l" {
		conc_g = *conc
	}

	if conc.Munit.BaseSISymbol() == "M/l" {
		conc_g = NewConcentration((conc.SIValue() * molecularweight), "g/l")
	}
	return conc_g
}

func (conc *Concentration) MolPerL(molecularweight float64) (conc_M Concentration) {
	if conc.Munit.BaseSISymbol() == "g/l" {
		conc_M = NewConcentration((conc.SIValue() / molecularweight), "M/l")
	}

	if conc.Munit.BaseSISymbol() == "kg/l" {
		conc_M = NewConcentration((conc.SIValue() / molecularweight), "M/l")
	}

	if conc.Munit.BaseSISymbol() == "M/l" {
		conc_M = *conc
	}
	return conc_M
}

// a structure which defines a specific heat capacity
type SpecificHeatCapacity struct {
	*ConcreteMeasurement
}

// make a new specific heat capacity structure in SI units
func NewSpecificHeatCapacity(v float64, unit string) SpecificHeatCapacity {
	if unit != "J/kg" {
		panic("Can't make specific heat capacities which aren't in J/kg")
	}

	s := SpecificHeatCapacity{NewMeasurement(v, "", unit)}
	return s
}

// a structure which defines a density
type Density struct {
	*ConcreteMeasurement
}

// make a new density structure in SI units
func NewDensity(v float64, unit string) Density {
	if unit != "kg/m^3" {
		panic("Can't make densities which aren't in kg/m^3")
	}

	d := Density{NewMeasurement(v, "", unit)}
	return d
}

type FlowRate struct {
	*ConcreteMeasurement
}

// new flow rate in ml/min

func NewFlowRate(v float64, unit string) FlowRate {
	if unit != "ml/min" {
		panic("Can't make flow rate not in ml/min")
	}
	fr := FlowRate{NewMeasurement(v, "", unit)}

	return fr
}

type Velocity struct {
	*ConcreteMeasurement
}

// new velocity in m/s

func NewVelocity(v float64, unit string) Velocity {

	if unit != "m/s" {
		panic("Can't make flow rate which isn't in m/s")
	}
	fr := Velocity{NewMeasurement(v, "", unit)}

	return fr
}

type Rate struct {
	*ConcreteMeasurement
}

func NewRate(v float64, unit string) (r Rate, err error) {
	if unit != `/min` && unit != `/s` {
		err = fmt.Errorf("Can't make flow rate which aren't in /min or per /s ")
		panic(err.Error())
	}

	approvedtimeunits := []string{"/min", "/s"}

	if unit[1:] == "min" {
		r := Rate{NewMeasurement(v*60, "", `/s`)}
		return r, nil
	} else if unit[1:] == "s" {
		r := Rate{NewMeasurement(v, "", `/s`)}
		return r, nil
	}

	err = fmt.Errorf(unit, " Not approved time unit. Approved units time are: ", approvedtimeunits)
	return r, err
}

type Voltage struct {
	*ConcreteMeasurement
}

func NewVoltage(value float64, unit string) (v Voltage, err error) {
	return Voltage{NewMeasurement(value, "", unit)}, nil
}
