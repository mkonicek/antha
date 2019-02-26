// wutil/round.go: Part of the Antha language
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

package wutil

import (
	"fmt"
	"math"

	"github.com/montanaflynn/stats"
)

func RoundInt(v float64) int {
	f := 1.0
	if v < 0 {
		f = -1.0
	}

	return int(f * (0.5 + (f * v)))
}

func RoundDown(v float64) (int, error) {
	if v < 0 {
		err := fmt.Errorf("Cannot round down negative float64")
		return 0, err
	}

	return int(math.Floor(v)), nil
}

func RoundIgnoreNan(input float64, p int) float64 {
	r, _ := Roundto(input, p) // nolint
	return r
}

func Roundto(input float64, places int) (rounded float64, err error) {
	rounded, err = stats.Round(input, places)

	return
}
