// /anthalib/simulator/liquidhandling/simulator.go: Part of the Antha language
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

package liquidhandling

import (
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"sort"
	"strings"
)

//checkTipPresence checks that tips are present (or not) on the adaptor at the given
//channels and returns a list of tips which are not (or are). If channels
//empty, all channels are tested
func checkTipPresence(present bool, adaptor *AdaptorState, channels []int) []int {

	ret := make([]int, 0, len(channels))
	if len(channels) == 0 {
		channels = make([]int, 0, adaptor.GetChannelCount())
		for i := 0; i < adaptor.GetChannelCount(); i++ {
			channels = append(channels, i)
		}
	}

	for _, ch := range channels {
		if adaptor.GetChannel(ch).HasTip() != present {
			ret = append(ret, ch)
		}
	}

	return ret
}

//getUnique strings optionally ignoring ""
func getUnique(slice []string, ignoreEmpty bool) []string {
	sMap := make(map[string]bool)
	ret := make([]string, 0, len(slice))
	for _, s := range slice {
		if ignoreEmpty && s == "" {
			continue
		}
		if _, ok := sMap[s]; !ok {
			sMap[s] = true
			ret = append(ret, s)
		}
	}
	return ret
}

//countUnique strings optionally ignoring ""
func countUnique(slice []string, ignoreEmpty bool) int {
	sMap := make(map[string]bool)
	for _, s := range slice {
		if ignoreEmpty && s == "" {
			continue
		}
		sMap[s] = true
	}
	return len(sMap)
}

//for a string slice in which each element is equal or "", return the non null element
func firstNonEmpty(slice []string) string {
	for _, s := range slice {
		if s != "" {
			return s
		}
	}
	return ""
}

//addComponent add a component to the container without storing component history
//all we care about are the volume and Cname
func addComponent(container wtype.LHContainer, rhs *wtype.LHComponent) error {

	lhs := container.Contents()

	ret := wtype.NewLHComponent()

	var names []string
	names = append(names, strings.Split(lhs.CName, "+")...)
	names = append(names, strings.Split(rhs.CName, "+")...)
	names = getUnique(names, true)
	sort.Strings(names)
	ret.CName = strings.Join(names, "+")

	fV := wunit.AddVolumes(lhs.Volume(), rhs.Volume())
	ret.Vol = fV.RawValue()
	ret.Vunit = fV.Unit().PrefixedSymbol()

	return container.SetContents(ret)
}

func coordsMatch(tc [][]wtype.WellCoords, wc []wtype.WellCoords) bool {
	if len(tc) != 1 {
		return false
	}

	if len(tc[0]) != len(wc) {
		return false
	}

	for i := 0; i < len(wc); i++ {
		if !tc[0][i].Equals(wc[0]) {
			return false
		}
	}

	return true
}
