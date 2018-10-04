// Part of the Antha language
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

// Package text formats strings for printing in a terminal using ansi codes
package text

import (
	"github.com/antha-lang/antha/antha/anthalib/wutil/text"
)

// Print prints to standard out a string description highlighted in red followed by a values in unformatted text
var Print = text.Print

// Sprint returns a string description highlighted in red followed by the values in unformatted text
var Sprint = text.Sprint

// Red changes string colour to red
var Red = text.Red

// Blue changes string colour to blue
var Blue = text.Blue

// Green changes string colour to green
var Green = text.Green

// Yellow changes string colour to yellow
var Yellow = text.Yellow

// Magenta changes string colour to magenta
var Magenta = text.Magenta

// Cyan changes string colour to cyan
var Cyan = text.Cyan

// White changes string colour to white
var White = text.White

// Black changes string colour to black
var Black = text.Black
