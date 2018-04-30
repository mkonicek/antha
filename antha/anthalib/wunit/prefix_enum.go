// anthalib//wunit/prefix_enum.go: Part of the Antha language
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

package wunit

type Prefix_Net int

const (
	//yzafpnum

	// y Prefix_Net = -24 + iota
	// z            = -24 + (iota * 3)
	// a            = -24 + (iota * 3)
	// f            = -24 + (iota * 3)
	// p            = -24 + (iota * 3)
	// n            = -24 + (iota * 3)
	// u            = -24 + (iota * 3)
	// m            = -24 + (iota * 3)

	// //cdh
	// c = -2
	// d = -1
	// h = 1

	//kMGTPEZY

	_ = -30 + (iota * 3)
	M = -30 + (iota * 3)
	G = -30 + (iota * 3)
	T = -30 + (iota * 3)
	P = -30 + (iota * 3)
	E = -30 + (iota * 3)
	Z = -30 + (iota * 3)
	Y = -30 + (iota * 3)
)
