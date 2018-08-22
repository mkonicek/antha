package wunit

import (
	"github.com/pkg/errors"
	"sort"
	"strings"
)

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
	"Area": {
		"m^2":  {Base: "m^2", Prefix: "", Multiplier: 1.0},
		"mm^2": {Base: "m^2", Prefix: "", Multiplier: 1.0e-6},
	},
}

func ValidUnitsForType(measurementType string) []string {
	var ret []string
	for key := range UnitMap[measurementType] {
		ret = append(ret, key)
	}
	sort.Strings(ret)
	return ret
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
		return errors.Errorf("No measurement type %s listed in UnitMap found these: %v", measureMentType, validMeasurementTypes)
	}

	_, unitFound := validUnits[unit]

	if !unitFound {
		var approved []string
		for u := range validUnits {
			approved = append(approved, u)
		}
		sort.Strings(approved)
		return errors.Errorf("No unit %s found for %s in UnitMap found these: %v", unit, measureMentType, approved)
	}

	return nil
}
