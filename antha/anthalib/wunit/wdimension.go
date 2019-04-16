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
)

// NewTypedMeasurement create a new measurement from the global registry asserting that
// the supplied units match the given type, and calling panic() if not
func NewTypedMeasurement(measurementType string, value float64, unit string) *ConcreteMeasurement {
	if ok := GetGlobalUnitRegistry().ValidUnitForType(measurementType, unit); !ok {
		panic(errors.Errorf("unknown units %q for measurement of type %q: only %v are supported",
			unit, measurementType, GetGlobalUnitRegistry().ListValidUnitsForType(measurementType)))
	}

	if value, err := GetGlobalUnitRegistry().NewMeasurement(value, unit); err != nil {
		panic(err)
	} else {
		return value
	}
}

// length
type Length struct {
	*ConcreteMeasurement
}

func ZeroLength() Length {
	return NewLength(0.0, "m")
}

// make a length
func NewLength(v float64, unit string) Length {
	return Length{NewTypedMeasurement("Length", v, unit)}
}

// area
type Area struct {
	*ConcreteMeasurement
}

// make an area unit
func NewArea(v float64, unit string) Area {
	return Area{NewTypedMeasurement("Area", v, unit)}
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
	return Volume{NewTypedMeasurement("Volume", v, unit)}
}

func CopyVolume(v Volume) Volume {
	if isNil(v.ConcreteMeasurement) {
		return Volume{}
	}
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
	if vol1.Unit().BaseSISymbol() != vol2.Unit().BaseSISymbol() {
		return 0, errors.Errorf("cannot divide volumes with incompatible units %v and %v", vol1.Unit(), vol2.Unit())
	}

	if vol2.IsZero() {
		return 0, errors.Errorf("while dividing volume %s by %s: cannot divide by zero", vol1, vol2)
	}

	return vol1.SIValue() / vol2.SIValue(), nil
}

// Dup deprecated, please use CopyConcentration
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
func DivideConcentrations(num, den Concentration) (float64, error) {
	if den.IsZero() {
		return 0, errors.Errorf("while dividing concentrations %s and %s: cannot divide by zero", num, den)
	} else if denInNumUnits, err := den.InUnit(num.Unit()); err != nil {
		return 0, err
	} else {
		return num.RawValue() / denInNumUnits.RawValue(), nil
	}
}

// AddConcentrations adds a variable number of concentrations from an original concentration.
// An error is returned if the concentration units are incompatible.
func AddConcentrations(concs ...Concentration) (Concentration, error) {

	if len(concs) == 0 {
		//since there were no concentrations, we don't know what units to return, so return SI standard ones
		return NewConcentration(0.0, "kg/l"), nil
	}

	ret := NewConcentration(0.0, concs[0].Unit().PrefixedSymbol())

	for _, conc := range concs {
		if err := ret.IncrBy(conc); err != nil {
			return ret, err
		}
	}
	return ret, nil

}

// SubtractConcentrations substracts a variable number of concentrations from an original concentration.
// An error is returned if the concentration units are incompatible.
func SubtractConcentrations(originalConc Concentration, subtractConcs ...Concentration) (Concentration, error) {

	ret := CopyConcentration(originalConc)

	if concToSubtract, err := AddConcentrations(subtractConcs...); err != nil {
		return Concentration{}, err
	} else if err := ret.DecrBy(concToSubtract); err != nil {
		return Concentration{}, err
	} else {
		return ret, nil
	}
}

// Dup deprecated, please use CopyVolume(Volume)
func (v Volume) Dup() Volume {
	if isNil(v.ConcreteMeasurement) {
		return ZeroVolume()
	}
	return NewVolume(v.RawValue(), v.Unit().PrefixedSymbol())
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
	return Temperature{NewTypedMeasurement("Temperature", v, unit)}
}

// time
type Time struct {
	*ConcreteMeasurement
}

// NewTime creates a time unit.
func NewTime(v float64, unit string) (t Time) {
	return Time{NewTypedMeasurement("Time", v, unit)}

}

func (t Time) Seconds() float64 {
	return t.SIValue()
}

func (t Time) AsDuration() time.Duration {
	// simply use the parser
	// make string output compatible with time.ParseDuration:
	// 1. avoid printing value in scientific notation (which time.Duration does not support)
	// 2. replace min with m (time.ParseDuration does not support min)
	str := strings.Replace(fmt.Sprint(t.RawValue(), t.Unit().PrefixedSymbol()), "min", "m", -1)
	// 2. remove white space (time.ParseDuration does not support white space)
	str = strings.Replace(str, " ", "", -1)

	if d, err := time.ParseDuration(str); err != nil {
		panic(err)
	} else {
		return d
	}
}

func FromDuration(t time.Duration) Time {
	return NewTime(float64(t.Seconds()), "s")
}

// CopyTime creates a safe duplicate of a time value.
func CopyTime(time Time) Time {
	return NewTime(time.RawValue(), time.Unit().PrefixedSymbol())
}

// AddTimes sums a variable number of Time arguments.
// If the dimension of the times are different the product will be returned in the SI value of Time (seconds).
func AddTimes(timesToAdd ...Time) (sum Time) {
	if len(timesToAdd) == 0 {
		return
	}

	sum = (NewTime(0.0, timesToAdd[0].Unit().PrefixedSymbol()))
	for _, time := range timesToAdd {
		sum.Add(time)
	}
	return sum
}

// SubtractTimes subtracts a variable number of Time arguments from timeToSubtractFrom.
func SubtractTimes(timeToSubtractFrom Time, timesToSubtract ...Time) (newTime Time) {
	newTime = (CopyTime(timeToSubtractFrom))
	timeToSubtract := AddTimes(timesToSubtract...)
	newTime.Subtract(timeToSubtract)
	return
}

// MultiplyTime multiplies a Time by a factor.
func MultiplyTime(v Time, factor float64) Time {
	return NewTime(v.RawValue()*float64(factor), v.Unit().PrefixedSymbol())
}

// DivideTime divides a Time by a factor.
func DivideTime(v Time, factor float64) Time {
	return NewTime(v.RawValue()/float64(factor), v.Unit().PrefixedSymbol())
}

// mass
type Mass struct {
	*ConcreteMeasurement
}

// make a mass unit

func NewMass(v float64, unit string) Mass {
	return Mass{NewTypedMeasurement("Mass", v, unit)}
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
	return Moles{NewTypedMeasurement("Moles", v, unit)}
}

// generate a new Amount in moles
func NewAmount(v float64, unit string) Moles {
	return Moles{NewTypedMeasurement("Moles", v, unit)}
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
	return Angle{NewTypedMeasurement("Angle", v, unit)}
}

// angular velocity (one way or another)

type AngularVelocity struct {
	*ConcreteMeasurement
}

func NewAngularVelocity(v float64, unit string) AngularVelocity {
	return AngularVelocity{NewTypedMeasurement("AngularVelocity", v, unit)}
}

// this is really Mass Length/Time^2
type Energy struct {
	*ConcreteMeasurement
}

// make a new energy unit
func NewEnergy(v float64, unit string) Energy {
	return Energy{NewTypedMeasurement("Energy", v, unit)}
}

// a Force
type Force struct {
	*ConcreteMeasurement
}

// a new force in Newtons
func NewForce(v float64, unit string) Force {
	return Force{NewTypedMeasurement("Force", v, unit)}
}

// a Pressure structure
type Pressure struct {
	*ConcreteMeasurement
}

// make a new pressure in Pascals
func NewPressure(v float64, unit string) Pressure {
	return Pressure{NewTypedMeasurement("Pressure", v, unit)}
}

// defines a concentration unit
type Concentration struct {
	*ConcreteMeasurement
	//MolecularWeight *float64
}

// NewConcentration makes a new concentration in SI units... either M/l or kg/l
func NewConcentration(v float64, unit string) Concentration {
	return Concentration{NewTypedMeasurement("Concentration", v, unit)}
}

// mass or mole
type SubstanceQuantity interface {
	Quantity() Measurement
}

// GramsPerLitre return a new concentration equal to the current one in grams per litre,
// using molecularweight given in grams per mole to convert if necessary.
// Returns an error if the units of conc are not compatible with grams per litre or
// grams per mole (such as "X" or "v/v")
func (conc Concentration) GramsPerLitre(molecularweight float64) (Concentration, error) {
	if concInGramsPerLitre, err := conc.InStringUnit("g/l"); err == nil {
		return Concentration{ConcreteMeasurement: concInGramsPerLitre.(*ConcreteMeasurement)}, nil
	} else if concInMolsPerLitre, err := conc.InStringUnit("Mol/l"); err != nil {
		return Concentration{}, errors.WithMessage(err, fmt.Sprintf("while converting %v into grams per litre[g/l]", conc.Munit))
	} else {
		return NewConcentration(concInMolsPerLitre.RawValue()*molecularweight, "g/l"), nil
	}
}

// GramPerL deprecated, please use GramsPerLitre
func (conc Concentration) GramPerL(molecularWeight float64) Concentration {
	if ret, err := conc.GramsPerLitre(molecularWeight); err != nil {
		panic(err)
	} else {
		return ret
	}
}

// MolesPerLitre return a new concentration equal to the current one in mols per litre,
// using molecularweight given in grams per mole to convert if necessary.
// Returns an error if the units of conc are not compatible with grams per litre or
// grams per mole (such as "X" or "v/v")
func (conc Concentration) MolesPerLitre(molecularweight float64) (Concentration, error) {
	if concInMolsPerLitre, err := conc.InStringUnit("Mol/l"); err == nil {
		return Concentration{ConcreteMeasurement: concInMolsPerLitre.(*ConcreteMeasurement)}, nil
	} else if concInGramsPerLitre, err := conc.InStringUnit("g/l"); err != nil {
		return Concentration{}, errors.WithMessage(err, fmt.Sprintf("while converting %v into moles per litre[Mol/l]", conc.Munit))
	} else {
		return NewConcentration(concInGramsPerLitre.RawValue()/molecularweight, "Mol/l"), nil
	}
}

// MolPerL deprecated, please use MolesPerLitre
func (conc Concentration) MolPerL(molecularWeight float64) Concentration {
	if ret, err := conc.MolesPerLitre(molecularWeight); err != nil {
		panic(err)
	} else {
		return ret
	}
}

// a structure which defines a specific heat capacity
type SpecificHeatCapacity struct {
	*ConcreteMeasurement
}

// make a new specific heat capacity structure in SI units
func NewSpecificHeatCapacity(v float64, unit string) SpecificHeatCapacity {
	return SpecificHeatCapacity{NewTypedMeasurement("SpecificHeatCapacity", v, unit)}
}

// a structure which defines a density
type Density struct {
	*ConcreteMeasurement
}

// make a new density structure in SI units
func NewDensity(v float64, unit string) Density {
	return Density{NewTypedMeasurement("Density", v, unit)}
}

type FlowRate struct {
	*ConcreteMeasurement
}

// new flow rate in ml/min

func NewFlowRate(v float64, unit string) FlowRate {
	return FlowRate{NewTypedMeasurement("FlowRate", v, unit)}
}

type Velocity struct {
	*ConcreteMeasurement
}

// new velocity in m/s

func NewVelocity(v float64, unit string) Velocity {
	return Velocity{NewTypedMeasurement("Velocity", v, unit)}
}

// Velocity3D struct composed of velocities in three axes
type Velocity3D struct {
	X, Y, Z Velocity
}

// GetAxis return the velocity in the axis specified
func (self *Velocity3D) GetAxis(a Axis) Velocity {
	switch a {
	case XAxis:
		return self.X
	case YAxis:
		return self.Y
	case ZAxis:
		return self.Z
	}
	panic("unknown axis")
}

// SetAxis return the velocity in the axis specified
func (self *Velocity3D) SetAxis(a Axis, v Velocity) {
	switch a {
	case XAxis:
		self.X = v
	case YAxis:
		self.Y = v
	case ZAxis:
		self.Z = v
	default:
		panic("unknown axis")
	}
}

// Dup return a copy of the velocities
func (self *Velocity3D) Dup() *Velocity3D {
	if self == nil {
		return nil
	}
	return &Velocity3D{
		X: NewVelocity(self.X.RawValue(), self.X.Unit().PrefixedSymbol()),
		Y: NewVelocity(self.Y.RawValue(), self.Y.Unit().PrefixedSymbol()),
		Z: NewVelocity(self.Z.RawValue(), self.Z.Unit().PrefixedSymbol()),
	}
}

type Rate struct {
	*ConcreteMeasurement
}

func NewRate(v float64, unit string) (r Rate, err error) {
	return Rate{NewTypedMeasurement("Rate", v, unit)}, nil
}

type Voltage struct {
	*ConcreteMeasurement
}

func NewVoltage(value float64, unit string) (Voltage, error) {
	return Voltage{NewTypedMeasurement("Voltage", value, unit)}, nil
}
