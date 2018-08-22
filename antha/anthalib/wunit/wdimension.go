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
	"math"
	"strings"
	"time"

	"github.com/pkg/errors"

	"github.com/antha-lang/antha/microArch/logger"
)

// length
type Length struct {
	*ConcreteMeasurement
}

func ZeroLength() Length {
	return NewLength(0.0, "m")
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
func NewArea(v float64, unit string) Area {
	details, ok := UnitMap["Area"][unit]
	if !ok {
		panic(errors.Errorf("unapproved area unit %q, approved units are %s", unit, ValidUnitsForType("Area")))
	}

	return Area{NewMeasurement((v * details.Multiplier), details.Prefix, details.Base)}
}

func ZeroArea() Area {
	return NewArea(0.0, "m^2")
}

// volume -- strictly speaking of course this is length^3
type Volume struct {
	*ConcreteMeasurement
}

// make a volume
func NewVolume(v float64, unit string) Volume {
	unit = strings.Replace(unit, "µ", "u", -1)

	if details, ok := UnitMap["Volume"][unit]; !ok {
		panic(fmt.Errorf("unknown volume unit %q, only the following units are supported: %v", unit, ValidUnitsForType("Volume")))
	} else {
		return Volume{NewMeasurement(v*details.Multiplier, details.Prefix, details.Base)}
	}
}

func CopyVolume(v Volume) Volume {
	ret := NewVolume(v.RawValue(), v.Unit().PrefixedSymbol())
	return ret
}

// AddVolumes adds a set of volumes.
func AddVolumes(vols ...Volume) (newvolume Volume) {
	tempvol := NewVolume(0.0, "ul")
	for _, vol := range vols {
		if tempvol.Unit().PrefixedSymbol() == vol.Unit().PrefixedSymbol() {
			tempvol = NewVolume(tempvol.RawValue()+vol.RawValue(), tempvol.Unit().PrefixedSymbol())
			newvolume = tempvol
		} else {
			tempvol = NewVolume(tempvol.SIValue()+vol.SIValue(), tempvol.Unit().BaseSISymbol())
			newvolume = tempvol
		}
	}
	return

}

// SubtractVolumes substracts a variable number of volumes from an original volume.
func SubtractVolumes(OriginalVol Volume, subtractvols ...Volume) (newvolume Volume) {

	newvolume = (CopyVolume(OriginalVol))
	volToSubtract := AddVolumes(subtractvols...)
	newvolume.Subtract(volToSubtract)

	if math.IsInf(newvolume.RawValue(), 0) {
		panic(errors.Errorf("Infinity value found subtracting volumes. Original: %s. Vols to subtract: %s", OriginalVol, subtractvols))
	}

	return

}

// MultiplyVolume multiplies a volume by a factor.
func MultiplyVolume(v Volume, factor float64) (newvolume Volume) {

	newvolume = NewVolume(v.RawValue()*float64(factor), v.Unit().PrefixedSymbol())
	return

}

// DivideVolume divides a volume by a factor.
func DivideVolume(v Volume, factor float64) (newvolume Volume) {

	newvolume = NewVolume(v.RawValue()/float64(factor), v.Unit().PrefixedSymbol())
	return

}

// DivideVolumes divides the SI Value of vol1 by vol2 to return a factor.
// An error is returned if the volume is infinity or not a number.
func DivideVolumes(vol1, vol2 Volume) (float64, error) {
	if vol1.Unit().BaseSIUnit() != vol2.Unit().BaseSIUnit() {
		return -1, errors.Errorf("cannot divide volumes: units of %s and %s unequal.", vol1.ToString(), vol2.ToString())
	}
	factor := vol1.SIValue() / vol2.SIValue()

	if math.IsInf(factor, 0) {
		return 0, errors.Errorf("infinity value found dividing volumes %s and %s", vol1.ToString(), vol2.ToString())
	}

	if math.IsNaN(factor) {
		return 0, errors.Errorf("NaN value found dividing volumes %s and %s", vol1.ToString(), vol2.ToString())
	}

	return factor, nil
}

func (c Concentration) Dup() Concentration {
	return CopyConcentration(c)
}

func CopyConcentration(v Concentration) Concentration {
	ret := NewConcentration(v.RawValue(), v.Unit().PrefixedSymbol())
	return ret
}

// MultiplyConcentration multiplies a concentration by a factor.
func MultiplyConcentration(v Concentration, factor float64) (newconc Concentration) {

	newconc = NewConcentration(v.RawValue()*float64(factor), v.Unit().PrefixedSymbol())
	return

}

// DivideConcentration divides a concentration by a factor.
func DivideConcentration(v Concentration, factor float64) (newconc Concentration) {

	newconc = NewConcentration(v.RawValue()/float64(factor), v.Unit().PrefixedSymbol())
	return

}

// DivideConcentrations divides the SI Value of conc1 by conc2 to return a factor.
// An error is returned if the concentration unit is not dividable or the number generated is infinity.
func DivideConcentrations(conc1, conc2 Concentration) (factor float64, err error) {
	if conc1.Unit().BaseSIUnit() != conc2.Unit().BaseSIUnit() {
		return -1, errors.Errorf("cannot divide concentrations: units of %s and %s unequal.", conc1.ToString(), conc2.ToString())
	}
	factor = conc1.SIValue() / conc2.SIValue()

	if math.IsInf(factor, 0) {
		err = errors.Errorf("infinity value found dividing concentrations %s and %s", conc1.ToString(), conc2.ToString())
		return
	}

	if math.IsNaN(factor) {
		err = errors.Errorf("NaN value found dividing concentrations %s and %s", conc1.ToString(), conc2.ToString())
		return
	}

	return factor, nil
}

// AddConcentrations adds a variable number of concentrations from an original concentration.
// An error is returned if the concentration units are incompatible.
func AddConcentrations(concs ...Concentration) (newconc Concentration, err error) {

	if len(concs) == 0 {
		err = errors.Errorf("Array of concentrations empty, nil value returned")
	}
	var tempconc Concentration
	unit := concs[0].Unit().PrefixedSymbol()
	tempconc = NewConcentration(0.0, unit)

	for _, conc := range concs {
		if tempconc.Unit().PrefixedSymbol() == conc.Unit().PrefixedSymbol() {
			tempconc = NewConcentration(tempconc.RawValue()+conc.RawValue(), tempconc.Unit().PrefixedSymbol())
			newconc = tempconc
		} else if tempconc.Unit().BaseSISymbol() != conc.Unit().BaseSISymbol() {
			err = errors.Errorf("Cannot add units with base %s to %s, please bring concs to same base. ", tempconc.Unit().BaseSISymbol(), conc.Unit().BaseSISymbol())
		} else {
			tempconc = NewConcentration(tempconc.SIValue()+conc.SIValue(), tempconc.Unit().BaseSISymbol())
			newconc = tempconc
		}
	}
	return

}

// SubtractConcentrations substracts a variable number of concentrations from an original concentration.
// An error is returned if the concentration units are incompatible.
func SubtractConcentrations(originalConc Concentration, subtractConcs ...Concentration) (newConcentration Concentration, err error) {

	newConcentration = (CopyConcentration(originalConc))

	concToSubtract, err := AddConcentrations(subtractConcs...)
	if err != nil {
		return
	}
	newConcentration.Subtract(concToSubtract)

	if math.IsInf(newConcentration.RawValue(), 0) {
		err = errors.Errorf(fmt.Sprintln("Infinity value found subtracting concentrations. Original: ", originalConc, ". Vols to subtract:", subtractConcs))
		return
	}

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
	details, ok := UnitMap["Temperature"][unit]
	if !ok {
		panic(errors.Errorf("unapproved temperature unit %q, approved units are %s", unit, ValidUnitsForType("Temperature")))
	}

	return Temperature{NewMeasurement((v * details.Multiplier), details.Prefix, details.Base)}
}

// time
type Time struct {
	*ConcreteMeasurement
}

// NewTime creates a time unit.
func NewTime(v float64, unit string) (t Time) {
	unit = strings.Replace(unit, "µ", "u", -1)

	details, ok := UnitMap["Time"][unit]
	if !ok {
		panic(errors.Errorf("unapproved time unit %q, approved units are %s", unit, ValidUnitsForType("Time")))
	}

	return Time{NewMeasurement((v * details.Multiplier), details.Prefix, details.Base)}
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

func NewMass(v float64, unit string) Mass {
	unit = strings.Replace(unit, "µ", "u", -1)

	details, ok := UnitMap["Mass"][unit]
	if !ok {
		panic(errors.Errorf("Can't make masses with non approved unit of %s. Approved units are: %v", unit, ValidUnitsForType("Mass")))
	}

	return Mass{NewMeasurement((v * details.Multiplier), details.Prefix, details.Base)}
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
func NewMoles(v float64, unit string) Moles {
	unit = strings.Replace(unit, "µ", "u", -1)

	details, ok := UnitMap["Moles"][unit]
	if !ok {
		panic(errors.Errorf("unapproved Amount unit %q, approved units are %s", unit, ValidUnitsForType("Moles")))
	}

	return Moles{NewMeasurement((v * details.Multiplier), details.Prefix, details.Base)}

}

// generate a new Amount in moles
func NewAmount(v float64, unit string) Moles {
	unit = strings.Replace(unit, "µ", "u", -1)

	details, ok := UnitMap["Moles"][unit]
	if !ok {
		panic(errors.Errorf("unapproved Amount unit %q, approved units are %s", unit, ValidUnitsForType("Moles")))
	}

	return Moles{NewMeasurement((v * details.Multiplier), details.Prefix, details.Base)}

}

// defines Moles to be a SubstanceQuantity
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

// NewConcentration makes a new concentration in SI units... either M/l or kg/l
func NewConcentration(v float64, unit string) Concentration {
	// replace µ with u
	unit = strings.Replace(unit, "µ", "u", -1)
	// replace L with l
	unit = strings.Replace(unit, "L", "l", -1)

	details, ok := UnitMap["Concentration"][unit]
	if !ok {
		panic(errors.Errorf("unapproved concentration unit %q, approved units are %s", unit, ValidUnitsForType("Concentration")))
	}

	return Concentration{NewMeasurement((v * details.Multiplier), details.Prefix, details.Base)}
}

// mass or mole
type SubstanceQuantity interface {
	Quantity() Measurement
}

func (conc Concentration) GramPerL(molecularweight float64) (conc_g Concentration) {

	if conc.Munit.BaseSISymbol() == "kg/l" {
		conc_g = conc
	}

	if conc.Munit.BaseSISymbol() == "M/l" {
		conc_g = NewConcentration((conc.SIValue() * molecularweight), "g/l")
	}
	return conc_g
}

func (conc Concentration) MolPerL(molecularweight float64) (conc_M Concentration) {

	if conc.Munit.BaseSISymbol() == "kg/l" {
		// convert from kg to g to work out g/mol
		conversionFactor := 1000.0
		conc_M = NewConcentration((conc.SIValue() * conversionFactor / molecularweight), "M/l")
	}

	if conc.Munit.BaseSISymbol() == "M/l" {
		conc_M = conc
	}
	return conc_M
}

// a structure which defines a specific heat capacity
type SpecificHeatCapacity struct {
	*ConcreteMeasurement
}

// make a new specific heat capacity structure in SI units
func NewSpecificHeatCapacity(v float64, unit string) SpecificHeatCapacity {
	if unit != "J/kg*C" {
		panic("Can't make specific heat capacities which aren't in J/kg*C")
	}

	s := SpecificHeatCapacity{NewMeasurement(v, "", unit)}
	fmt.Println(s.Unit().ToString())
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
	details, ok := UnitMap["Rate"][unit]
	if !ok {
		return r, errors.Errorf("unapproved rate unit %q, approved units are %s", unit, ValidUnitsForType("Rate"))
	}

	return Rate{NewMeasurement((v * details.Multiplier), details.Prefix, details.Base)}, nil
}

type Voltage struct {
	*ConcreteMeasurement
}

func NewVoltage(value float64, unit string) (v Voltage, err error) {
	return Voltage{NewMeasurement(value, "", unit)}, nil
}
