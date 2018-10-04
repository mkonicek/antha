// antha/AnthaStandardLibrary/Packages/Platereader/Platereader.go: Part of the Antha language
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

package platereader

// ReadMode defines the mode of plate reader reading.
type ReadMode string // or could just make this a string

const (
	// Absorbance defines the reading to be an absorbance reading.
	Absorbance ReadMode = "Absorbance"
	// Fluorescence defines the reading to be a fluorescence reading.
	Fluorescence ReadMode = "Fluorescence"
)

// ReadType defines the type of plate reader reading.
type ReadType string // or could just make this a string

const (
	// Endpoint defines the reading to be an endpoint.
	Endpoint ReadType = "Endpoint"
)

// FilterOption defines the type of platereader data which will be returned
// by the ReadingsAsAverage method on a dataset.
type FilterOption int

const (
	// TIME will return readings matching a specified time
	TIME FilterOption = iota
	// EMWAVELENGTH will return readings matching a specified Emmission wavelength
	EMWAVELENGTH
	// EXWAVELENGTH will return readings matching a specified Excitation wavelength
	EXWAVELENGTH
)
