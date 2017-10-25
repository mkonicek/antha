// buffercomponents.go Part of the Antha language
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
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/pubchem"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
)

func StockConcentration(nameofmolecule string, massofmoleculeactuallyaddedinG wunit.Mass, diluent string, totalvolumeinL wunit.Volume) (actualconc wunit.Concentration, err error) {
	molecule, err := pubchem.MakeMolecule(nameofmolecule)
	if err != nil {
		return
	}

	// in particular, the molecular weight
	molecularweight := molecule.MolecularWeight

	actualconcfloat := (massofmoleculeactuallyaddedinG.SIValue() * 1000) / (molecularweight * totalvolumeinL.SIValue())

	actualconc = wunit.NewConcentration(actualconcfloat, "M/l")

	return
}

func Dilute(moleculename string, stockconc wunit.Concentration, stockvolume wunit.Volume, diluentname string, diluentvoladded wunit.Volume) (dilutedconc wunit.Concentration, err error) {
	molecule, err := pubchem.MakeMolecule(moleculename)
	if err != nil {
		return
	}

	stockMperL := stockconc.MolPerL(molecule.MolecularWeight)

	diluentSI := diluentvoladded.SIValue()

	stockSI := stockvolume.SIValue()

	dilutedconcMperL := stockMperL.SIValue() * stockSI / (stockSI + diluentSI)

	dilutedconc = wunit.NewConcentration(dilutedconcMperL, "M/l")
	return
}

func DiluteBasedonMolecularWeight(molecularweight float64, stockconc wunit.Concentration, stockvolume wunit.Volume, diluentname string, diluentvoladded wunit.Volume) (dilutedconc wunit.Concentration) {

	stockMperL := stockconc.MolPerL(molecularweight)

	diluentSI := diluentvoladded.SIValue()

	stockSI := stockvolume.SIValue()

	dilutedconcMperL := stockMperL.SIValue() * stockSI / (stockSI + diluentSI)

	dilutedconc = wunit.NewConcentration(dilutedconcMperL, "M/l")
	// fmt.Println(diluentname)
	return
}
