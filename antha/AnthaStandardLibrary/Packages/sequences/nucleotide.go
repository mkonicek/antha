// antha/AnthaStandardLibrary/Packages/enzymes/nucleotides.go: Part of the Antha language
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

// Package sequences is for interacting with and manipulating biological sequences; in extension to methods available in wtype
package sequences

import (
	"strings"

	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/search"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
)

// Check for illegal nucleotides
func Illegalnucleotides(fwdsequence wtype.DNASequence) (pass bool, illegalfound []search.Result, wobblefound []search.Result) {
	illegal := "§1234567890-=qeiop[]fjl;'z,./!@£$%^&*()_+?" // removed all instances of non IUPAC nucleotides
	wobble := "NXBHVDMKSWRYU"                               //IUPAC nucleotides

	if strings.ContainsAny(strings.ToUpper(fwdsequence.Seq), (strings.ToUpper(illegal))) || strings.ContainsAny(fwdsequence.Seq, strings.ToLower(illegal)) {
		pass = false
		illegalarray := strings.Split(illegal, "")
		illegalfound = search.FindAllStrings((strings.ToLower(fwdsequence.Seq)), illegalarray)
	} else if strings.ContainsAny(strings.ToUpper(fwdsequence.Seq), wobble) || strings.ContainsAny(fwdsequence.Seq, strings.ToLower(wobble)) {
		pass = false
		wobblearray := strings.Split(wobble, "")
		wobblefound = search.FindAllStrings((strings.ToUpper(fwdsequence.Seq)), wobblearray)
	} else {
		pass = true
	}

	return pass, illegalfound, wobblefound
}

// WobbleMap represents a mapping of each IUPAC nucleotide
// to all valid alternative IUPAC nucleotides for that nucleotide.
// This may be useful for protein engineering applications where mutations may wish to be introduced.
//
// For example N can be substituted for any primary nucleotide (A, C, T or G).
// R may be substituted for any purine base (A, G).
// gaps are represented by - or .
//
var WobbleMap = map[string][]string{
	"A": {"A"},
	"T": {"T"},
	"U": {"U"},
	"C": {"C"},
	"G": {"G"},
	"a": {"A"},
	"t": {"T"},
	"u": {"U"},
	"c": {"C"},
	"g": {"G"},
	"Y": {"C", "T"},
	"R": {"A", "G"},
	"W": {"A", "T"},
	"S": {"G", "C"},
	"K": {"G", "T"},
	"M": {"A", "C"},
	"D": {"A", "G", "T"},
	"V": {"A", "C", "G"},
	"H": {"A", "C", "T"},
	"B": {"C", "G", "T"},
	"N": {"A", "T", "C", "G"},
	"X": {"A", "T", "C", "G"},
	"-": {"-", "."},
	".": {"-", "."},
}

// Wobble returns an array of sequence options, as strings.
// Options are caclulated based on cross referencing each nucleotide with
// the WobbleMap to find each alternative option, if any, for that nucleotide.
// For example:
// ACT would return one sequence: ACT
// RCT would return ACT and GCT
// NCT would return ACT, GCT, TCT, CCT
// RYT would return ACT, GCT, ATT and GTT
//
func Wobble(seq string) (alloptions []string) {

	arrayofarray := make([][]string, 0)
	for _, character := range seq {

		optionsforcharacterx := WobbleMap[strings.TrimSpace(strings.ToUpper(string(character)))]
		arrayofarray = append(arrayofarray, optionsforcharacterx)
	}

	alloptions = allCombinations(arrayofarray)

	return
}

func allCombinations(arr [][]string) []string {
	if len(arr) == 1 {
		return arr[0]
	}

	results := make([]string, 0)
	allRem := allCombinations(arr[1:])
	for i := 0; i < len(allRem); i++ {
		for j := 0; j < len(arr[0]); j++ {
			x := arr[0][j] + allRem[i]
			results = append(results, x)
		}
	}
	return results
}

var Nucleotidegpermol = map[string]float64{
	"A":    313.2,
	"T":    304.2,
	"C":    289.2,
	"G":    329.2,
	"N":    303.7,
	"dATP": 491.2,
	"dCTP": 467.2,
	"dGTP": 507.2,
	"dTTP": 482.2,
	"dNTP": 487.0,
}

// Calculate global GC content
func GCcontent(fwdsequence string) (Percentage float64) {
	fwdsequence = strings.ToUpper(fwdsequence)

	numberofAs := strings.Count(fwdsequence, "A")
	numberofTs := strings.Count(fwdsequence, "T")
	numberofCs := strings.Count(fwdsequence, "C")
	numberofGs := strings.Count(fwdsequence, "G")

	gc := (numberofCs + numberofGs)
	all := (numberofCs + numberofGs + numberofAs + numberofTs)

	percentage := float64(gc) / float64(all)
	Percentage = percentage
	return Percentage
}

// Calculate local GC content in a sliding window
func localGCContent(fwdsequence string, window int, shift int) (Pc []float64) {
	incs := len(fwdsequence) / shift
	pos := 0
	Pc = make([]float64, 0)
	for i := 0; i < (incs - 1); i++ {
		region := fwdsequence[pos : pos+window]
		gc := GCcontent(region)
		Pc = append(Pc, gc)
		pos += shift
	}
	return Pc
}

//Calculate Molecular weight of DNA
func MassDNA(fwdsequence string, phosphate5prime bool, doublestranded bool) (mw float64) {
	numberofAs := strings.Count(fwdsequence, "A")
	numberofTs := strings.Count(fwdsequence, "T")
	numberofCs := strings.Count(fwdsequence, "C")
	numberofGs := strings.Count(fwdsequence, "G")
	massofAs := (float64(numberofAs) * Nucleotidegpermol["A"])
	massofTs := (float64(numberofTs) * Nucleotidegpermol["T"])
	massofCs := (float64(numberofCs) * Nucleotidegpermol["C"])
	massofGs := (float64(numberofGs) * Nucleotidegpermol["G"])
	mw = (massofAs + massofTs + massofCs + massofGs)
	if phosphate5prime {
		mw = mw + 79.0 // extra for phosphate left at 5' end following digestion, not relevant for primer extension
	}
	if doublestranded {
		mw = 2 * mw
	}
	return mw
}

// Calclulate number of moles of a mass of DNA
func MolesDNA(mass wunit.Mass, mw float64) (moles float64) {
	massSI := mass.SIValue()
	moles = massSI / mw
	return moles
}

// calculate molar concentration of DNA sample
func GtoMolarConc(conc wunit.Concentration, mw float64) (molesperL float64) {
	// convert SI kg/l into g/l
	concSI := conc.SIValue() * 1000
	molesperL = concSI / mw
	return molesperL
}

func MoletoGConc(molarconc float64, mw float64) (gperL wunit.Concentration) {
	gperl := molarconc * mw
	gperL = wunit.NewConcentration(gperl, "g/L")
	return gperL
}

func Moles(conc wunit.Concentration, mw float64, vol wunit.Volume) (moles float64) {
	molesperL := GtoMolarConc(conc, mw)
	moles = molesperL * vol.SIValue()
	return moles
}

func RevArrayOrder(array []string) (reversedOrder []string) {
	for i := len(array) - 1; i >= 0; i-- {
		reversedOrder = append(reversedOrder, array[i])
	}
	return reversedOrder
}
