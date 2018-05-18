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
	"sort"
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

func ZeroLength() Length {
	return EZLength(0.0)
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

func ZeroArea() Area {
	return NewArea(0.0, "m^2")
}

// volume -- strictly speaking of course this is length^3
type Volume struct {
	*ConcreteMeasurement
}

// make a volume
func NewVolume(v float64, unit string) (o Volume) {
	unit = strings.Replace(unit, "µ", "u", -1)

	if len(strings.TrimSpace(unit)) == 0 {
		return ZeroVolume()
	}

	o = Volume{NewPMeasurement(v, unit)}

	return
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
		panic(fmt.Sprintln("Infinity value found subtracting volumes. Original: ", OriginalVol, ". Vols to subtract:", subtractvols))
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
func DivideVolumes(vol1, vol2 Volume) (factor float64, err error) {
	if vol1.Unit().BaseSIUnit() != vol2.Unit().BaseSIUnit() {
		return -1, fmt.Errorf("cannot divide volumes: units of %s and %s unequal.", vol1.Summary(), vol2.Summary())
	}
	factor = vol1.SIValue() / vol2.SIValue()

	if math.IsInf(factor, 0) {
		err = fmt.Errorf("infinity value found dividing volumes %s and %s", vol1.Summary(), vol2.Summary())
		return
	}

	if math.IsNaN(factor) {
		err = fmt.Errorf("NaN value found dividing volumes %s and %s", vol1.Summary(), vol2.Summary())
		return
	}

	return factor, nil
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
		return -1, fmt.Errorf("cannot divide concentrations: units of %s and %s unequal.", conc1.Summary(), conc2.Summary())
	}
	factor = conc1.SIValue() / conc2.SIValue()

	if math.IsInf(factor, 0) {
		err = fmt.Errorf("infinity value found dividing concentrations %s and %s", conc1.Summary(), conc2.Summary())
		return
	}

	if math.IsNaN(factor) {
		err = fmt.Errorf("NaN value found dividing concentrations %s and %s", conc1.Summary(), conc2.Summary())
		return
	}

	return factor, nil
}

// AddConcentrations adds a variable number of concentrations from an original concentration.
// An error is returned if the concentration units are incompatible.
func AddConcentrations(concs ...Concentration) (newconc Concentration, err error) {

	if len(concs) == 0 {
		err = fmt.Errorf("Array of concentrations empty, nil value returned")
	}
	var tempconc Concentration
	unit := concs[0].Unit().PrefixedSymbol()
	tempconc = NewConcentration(0.0, unit)

	for _, conc := range concs {
		if tempconc.Unit().PrefixedSymbol() == conc.Unit().PrefixedSymbol() {
			tempconc = NewConcentration(tempconc.RawValue()+conc.RawValue(), tempconc.Unit().PrefixedSymbol())
			newconc = tempconc
		} else if tempconc.Unit().BaseSISymbol() != conc.Unit().BaseSISymbol() {
			err = fmt.Errorf("Cannot add units with base %s to %s, please bring concs to same base. ", tempconc.Unit().BaseSISymbol(), conc.Unit().BaseSISymbol())
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
		err = fmt.Errorf(fmt.Sprintln("Infinity value found subtracting concentrations. Original: ", originalConc, ". Vols to subtract:", subtractConcs))
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
		var approved []string
		for u := range UnitMap["Temperature"] {
			approved = append(approved, u)
		}
		sort.Strings(approved)
		panic(fmt.Sprintf("unapproved temperature unit %q, approved units are %s", unit, approved))
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
		var approved []string
		for u := range UnitMap["Time"] {
			approved = append(approved, u)
		}
		sort.Strings(approved)
		panic(fmt.Sprintf("unapproved time unit %q, approved units are %s", unit, approved))
	}

	return Time{NewMeasurement((v * details.Multiplier), details.Prefix, details.Base)}
}

func (t Time) Seconds() float64 {
	return t.SIValue()
}

func (t Time) AsDuration() time.Duration {
	// simply use the parser

	d, e := time.ParseDuration(t.Summary())

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
	unit = strings.Replace(unit, "µ", "u", -1)

	approvedunits := UnitMap["Mass"]

	var approved bool
	for key := range approvedunits {

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

	return
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
		var approved []string
		for u := range UnitMap["Moles"] {
			approved = append(approved, u)
		}
		sort.Strings(approved)
		panic(fmt.Sprintf("unapproved Amount unit %q, approved units are %s", unit, approved))
	}

	return Moles{NewMeasurement((v * details.Multiplier), details.Prefix, details.Base)}

}

// generate a new Amount in moles
func NewAmount(v float64, unit string) Moles {
	unit = strings.Replace(unit, "µ", "u", -1)

	details, ok := UnitMap["Moles"][unit]
	if !ok {
		var approved []string
		for u := range UnitMap["Moles"] {
			approved = append(approved, u)
		}
		sort.Strings(approved)
		panic(fmt.Sprintf("unapproved Amount unit %q, approved units are %s", unit, approved))
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

// Unit is the form which units are stored in the UnitMap. This structure is not used beyond this.
type Unit struct {
	Base       string
	Prefix     string
	Multiplier float64
}

// UnitMap lists approved units to create new measurements.
var UnitMap = map[string]map[string]Unit{
	"Concentration": {
		"kg/l":    {Base: "g/l", Prefix: "k", Multiplier: 1.0},
		"g/l":     {Base: "g/l", Prefix: "", Multiplier: 1.0},
		"mg/l":    {Base: "g/l", Prefix: "m", Multiplier: 1.0},
		"ug/l":    {Base: "g/l", Prefix: "u", Multiplier: 1.0},
		"ng/l":    {Base: "g/l", Prefix: "n", Multiplier: 1.0},
		"mg/ml":   {Base: "g/l", Prefix: "", Multiplier: 1.0},
		"ug/ml":   {Base: "g/l", Prefix: "m", Multiplier: 1.0},
		"ug/ul":   {Base: "g/l", Prefix: "", Multiplier: 1.0},
		"ng/ul":   {Base: "g/l", Prefix: "m", Multiplier: 1.0},
		"ng/ml":   {Base: "g/l", Prefix: "u", Multiplier: 1.0},
		"pg/ul":   {Base: "g/l", Prefix: "u", Multiplier: 1.0},
		"pg/ml":   {Base: "g/l", Prefix: "n", Multiplier: 1.0},
		"pg/l":    {Base: "g/l", Prefix: "p", Multiplier: 1.0},
		"kg/L":    {Base: "g/l", Prefix: "k", Multiplier: 1.0},
		"g/L":     {Base: "g/l", Prefix: "", Multiplier: 1.0},
		"mg/L":    {Base: "g/l", Prefix: "m", Multiplier: 1.0},
		"ug/L":    {Base: "g/l", Prefix: "u", Multiplier: 1.0},
		"ng/L":    {Base: "g/l", Prefix: "n", Multiplier: 1.0},
		"pg/L":    {Base: "g/l", Prefix: "p", Multiplier: 1.0},
		"mg/mL":   {Base: "g/l", Prefix: "", Multiplier: 1.0},
		"ug/mL":   {Base: "g/l", Prefix: "m", Multiplier: 1.0},
		"ug/uL":   {Base: "g/l", Prefix: "", Multiplier: 1.0},
		"ng/uL":   {Base: "g/l", Prefix: "m", Multiplier: 1.0},
		"ng/mL":   {Base: "g/l", Prefix: "u", Multiplier: 1.0},
		"pg/uL":   {Base: "g/l", Prefix: "u", Multiplier: 1.0},
		"pg/mL":   {Base: "g/l", Prefix: "n", Multiplier: 1.0},
		"M":       {Base: "M/l", Prefix: "", Multiplier: 1.0},
		"M/l":     {Base: "M/l", Prefix: "", Multiplier: 1.0},
		"Mol/l":   {Base: "M/l", Prefix: "", Multiplier: 1.0},
		"M/L":     {Base: "M/l", Prefix: "", Multiplier: 1.0},
		"Mol/L":   {Base: "M/l", Prefix: "", Multiplier: 1.0},
		"mM":      {Base: "M/l", Prefix: "m", Multiplier: 1.0},
		"mM/l":    {Base: "M/l", Prefix: "m", Multiplier: 1.0},
		"mMol/l":  {Base: "M/l", Prefix: "m", Multiplier: 1.0},
		"mM/L":    {Base: "M/l", Prefix: "m", Multiplier: 1.0},
		"mMol/L":  {Base: "M/l", Prefix: "m", Multiplier: 1.0},
		"uM":      {Base: "M/l", Prefix: "u", Multiplier: 1.0},
		"uM/l":    {Base: "M/l", Prefix: "u", Multiplier: 1.0},
		"uMol/l":  {Base: "M/l", Prefix: "u", Multiplier: 1.0},
		"uM/L":    {Base: "M/l", Prefix: "u", Multiplier: 1.0},
		"uMol/L":  {Base: "M/l", Prefix: "u", Multiplier: 1.0},
		"nM":      {Base: "M/l", Prefix: "n", Multiplier: 1.0},
		"nM/l":    {Base: "M/l", Prefix: "n", Multiplier: 1.0},
		"nMol/l":  {Base: "M/l", Prefix: "n", Multiplier: 1.0},
		"nM/L":    {Base: "M/l", Prefix: "n", Multiplier: 1.0},
		"nMol/L":  {Base: "M/l", Prefix: "n", Multiplier: 1.0},
		"pM":      {Base: "M/l", Prefix: "p", Multiplier: 1.0},
		"pM/l":    {Base: "M/l", Prefix: "p", Multiplier: 1.0},
		"pMol/l":  {Base: "M/l", Prefix: "p", Multiplier: 1.0},
		"pM/L":    {Base: "M/l", Prefix: "p", Multiplier: 1.0},
		"pMol/L":  {Base: "M/l", Prefix: "p", Multiplier: 1.0},
		"pM/ul":   {Base: "M/l", Prefix: "u", Multiplier: 1.0},
		"pMol/ul": {Base: "M/l", Prefix: "u", Multiplier: 1.0},
		"pM/uL":   {Base: "M/l", Prefix: "u", Multiplier: 1.0},
		"pMol/uL": {Base: "M/l", Prefix: "u", Multiplier: 1.0},
		"fM":      {Base: "M/l", Prefix: "f", Multiplier: 1.0},
		"fM/l":    {Base: "M/l", Prefix: "f", Multiplier: 1.0},
		"fMol/l":  {Base: "M/l", Prefix: "f", Multiplier: 1.0},
		"fM/L":    {Base: "M/l", Prefix: "f", Multiplier: 1.0},
		"fMol/L":  {Base: "M/l", Prefix: "f", Multiplier: 1.0},
		"fM/ul":   {Base: "M/l", Prefix: "n", Multiplier: 1.0},
		"fMol/ul": {Base: "M/l", Prefix: "n", Multiplier: 1.0},
		"fM/uL":   {Base: "M/l", Prefix: "n", Multiplier: 1.0},
		"fMol/uL": {Base: "M/l", Prefix: "n", Multiplier: 1.0},
		"X":       {Base: "X", Prefix: "", Multiplier: 1.0},
		"x":       {Base: "X", Prefix: "", Multiplier: 1.0},
		"U/l":     {Base: "U/l", Prefix: "", Multiplier: 1.0},
		"U/L":     {Base: "U/l", Prefix: "", Multiplier: 1.0},
		"U/ml":    {Base: "U/l", Prefix: "", Multiplier: 1000.0},
		"U/mL":    {Base: "U/l", Prefix: "", Multiplier: 1000.0},
		"v/v":     {Base: "v/v", Prefix: "", Multiplier: 1.0},
		"w/v":     {Base: "g/l", Prefix: "k", Multiplier: 1.0},
	},
	"Mass": {
		"ng": {Base: "g", Prefix: "n", Multiplier: 1.0},
		"ug": {Base: "g", Prefix: "u", Multiplier: 1.0},
		"mg": {Base: "g", Prefix: "m", Multiplier: 1.0},
		"g":  {Base: "g", Prefix: "", Multiplier: 1.0},
		"kg": {Base: "g", Prefix: "k", Multiplier: 1.0},
	},
	"Moles": {
		"pMol": {Base: "M", Prefix: "p", Multiplier: 1.0},
		"nMol": {Base: "M", Prefix: "n", Multiplier: 1.0},
		"uMol": {Base: "M", Prefix: "u", Multiplier: 1.0},
		"mMol": {Base: "M", Prefix: "m", Multiplier: 1.0},
		"Mol":  {Base: "M", Prefix: "", Multiplier: 1.0},
		"pM":   {Base: "M", Prefix: "p", Multiplier: 1.0},
		"nM":   {Base: "M", Prefix: "n", Multiplier: 1.0},
		"uM":   {Base: "M", Prefix: "u", Multiplier: 1.0},
		"mM":   {Base: "M", Prefix: "m", Multiplier: 1.0},
		"M":    {Base: "M", Prefix: "", Multiplier: 1.0},
	},
	"Volume": {
		"pl": {Base: "l", Prefix: "p", Multiplier: 1.0},
		"nl": {Base: "l", Prefix: "n", Multiplier: 1.0},
		"ul": {Base: "l", Prefix: "u", Multiplier: 1.0},
		"ml": {Base: "l", Prefix: "m", Multiplier: 1.0},
		"l":  {Base: "l", Prefix: "", Multiplier: 1.0},
		"pL": {Base: "l", Prefix: "p", Multiplier: 1.0},
		"nL": {Base: "l", Prefix: "n", Multiplier: 1.0},
		"uL": {Base: "l", Prefix: "u", Multiplier: 1.0},
		"mL": {Base: "l", Prefix: "m", Multiplier: 1.0},
		"L":  {Base: "l", Prefix: "", Multiplier: 1.0},
	},
	"Rate": {
		"/s":   {Base: "/s", Prefix: "", Multiplier: 1.0},
		"/min": {Base: "/s", Prefix: "", Multiplier: 60.0},
		"/h":   {Base: "/s", Prefix: "", Multiplier: 3600.0},
	},
	"Time": {
		"ms":   {Base: "s", Prefix: "m", Multiplier: 1.0},
		"s":    {Base: "s", Prefix: "", Multiplier: 1.0},
		"min":  {Base: "s", Prefix: "", Multiplier: 60.0},
		"h":    {Base: "s", Prefix: "", Multiplier: 3600.0},
		"days": {Base: "s", Prefix: "", Multiplier: 86400.0},
	},
	"Temperature": {
		"C":  {Base: "℃", Prefix: "", Multiplier: 1.0}, // RING ABOVE, LATIN CAPITAL LETTER C
		"˚C": {Base: "℃", Prefix: "", Multiplier: 1.0}, // LATIN CAPITAL LETTER C
		"℃":  {Base: "℃", Prefix: "", Multiplier: 1.0}, // DEGREE CELSIUS
		"°C": {Base: "℃", Prefix: "", Multiplier: 1.0}, // DEGREE, LATIN CAPITAL LETTER C
	},
}

// ValidMeasurementUnit checks the validity of a measurement type and unit within that measurement type.
// An error is returned if an invalid measurement type or unit is specified.
func ValidMeasurementUnit(measureMentType, unit string) error {
	// replace µ with u
	unit = strings.Replace(unit, "µ", "u", -1)
	if measureMentType == "Concentration" {
		// replace L with l
		unit = strings.Replace(unit, "L", "l", -1)
	}
	validUnits, measurementFound := UnitMap[measureMentType]
	if !measurementFound {
		var validMeasurementTypes []string
		for key := range UnitMap {
			validMeasurementTypes = append(validMeasurementTypes, key)
		}
		sort.Strings(validMeasurementTypes)
		return fmt.Errorf("No measurement type %s listed in UnitMap found these: %v", measureMentType, validMeasurementTypes)
	}

	_, unitFound := validUnits[unit]

	if !unitFound {
		var approved []string
		for u := range validUnits {
			approved = append(approved, u)
		}
		sort.Strings(approved)
		return fmt.Errorf("No unit %s found for %s in UnitMap found these: %v", unit, measureMentType, approved)
	}

	return nil
}

// ValidConcentrationUnit returns an error if an invalid Concentration unit is specified.
func ValidConcentrationUnit(unit string) error {
	// replace µ with u
	unit = strings.Replace(unit, "µ", "u", -1)
	// replace L with l
	unit = strings.Replace(unit, "L", "l", -1)
	_, ok := UnitMap["Concentration"][unit]
	if !ok {
		var approved []string
		for u := range UnitMap["Concentration"] {
			approved = append(approved, u)
		}
		sort.Strings(approved)
		return fmt.Errorf("unapproved concentration unit %q, approved units are %s", unit, approved)
	}
	return nil
}

// NewConcentration makes a new concentration in SI units... either M/l or kg/l
func NewConcentration(v float64, unit string) Concentration {
	// replace µ with u
	unit = strings.Replace(unit, "µ", "u", -1)
	// replace L with l
	unit = strings.Replace(unit, "L", "l", -1)

	details, ok := UnitMap["Concentration"][unit]
	if !ok {
		var approved []string
		for u := range UnitMap["Concentration"] {
			approved = append(approved, u)
		}
		sort.Strings(approved)
		panic(fmt.Sprintf("unapproved concentration unit %q, approved units are %s", unit, approved))
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
	details, ok := UnitMap["Rate"][unit]
	if !ok {
		var approved []string
		for u := range UnitMap["Rate"] {
			approved = append(approved, u)
		}
		sort.Strings(approved)
		return r, fmt.Errorf("unapproved rate unit %q, approved units are %s", unit, approved)
	}

	return Rate{NewMeasurement((v * details.Multiplier), details.Prefix, details.Base)}, nil
}

type Voltage struct {
	*ConcreteMeasurement
}

func NewVoltage(value float64, unit string) (v Voltage, err error) {
	return Voltage{NewMeasurement(value, "", unit)}, nil
}
