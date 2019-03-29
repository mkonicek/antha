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
	"fmt"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/pkg/errors"
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

//getSingle given an array of strings, test that each value is either equal to each other or empty
//and return the value, or an error if this is not the case
func getSingle(slice []string) (string, error) {
	u := getUnique(slice, true)
	if len(u) == 0 {
		return "", errors.New("no value specified")
	}
	if len(u) > 1 {
		return "", errors.Errorf("multiple values specified: \"%s\"", strings.Join(u, "\", \""))
	}
	return u[0], nil
}

var refMap = map[int]wtype.WellReference{
	0: wtype.BottomReference,
	1: wtype.TopReference,
	2: wtype.LiquidReference,
}

//convertReferences from an slice of ints into enum values, returning an error for unknown references
func convertReferences(slice []int) ([]wtype.WellReference, error) {
	refs := make([]wtype.WellReference, len(slice))
	unknown := make(map[int]bool)
	for i, r := range slice {
		if v, ok := refMap[r]; ok {
			refs[i] = v
		} else {
			unknown[r] = true
			//default to TopReference as it's safest
			refs[i] = wtype.TopReference
		}
	}

	if len(unknown) > 0 {
		uk := make([]string, 0, len(unknown))
		for v := range unknown {
			uk = append(uk, fmt.Sprintf("%d", v))
		}
		value := "value"
		if len(unknown) > 1 {
			value = "values"
		}
		return refs, errors.Errorf("unknown %s %s", value, strings.Join(uk, ", "))
	}
	return refs, nil
}

//convertWellCoords convert a list of string wellcoords to wtype.WellCoords returning
//errors for non-empty strings which can't be parsed
func convertWellCoords(slice []string) ([]wtype.WellCoords, error) {
	ret := make([]wtype.WellCoords, len(slice))
	unknown := make([]string, 0, len(slice))
	for i, s := range slice {
		if s == "" {
			ret[i] = wtype.ZeroWellCoords()
			continue
		}
		ret[i] = wtype.MakeWellCoords(s)
		if ret[i].IsZero() {
			unknown = append(unknown, s)
		}
	}

	if len(unknown) > 0 {
		return ret, errors.Errorf("couldn't parse \"%s\"", strings.Join(unknown, "\", \""))
	}
	return ret, nil
}

//assertNoTipsOnOthersInGroup check that there are no tips loaded on any other adaptors in the adaptor group
func assertNoTipsOnOthersInGroup(adaptor *AdaptorState) error {
	adaptors := make([]int, 0, adaptor.GetGroup().NumAdaptors())
	foundTips := make(map[int][]int)
	numTips := 0
	for _, ad := range adaptor.GetGroup().GetAdaptors() {
		if ad == adaptor {
			continue
		}
		if tipFound := checkTipPresence(false, ad, nil); len(tipFound) != 0 {
			idx := ad.GetIndex()
			adaptors = append(adaptors, idx)
			foundTips[idx] = tipFound
			numTips += len(tipFound)
		}
	}

	if len(foundTips) > 0 {
		s := make([]string, 0, len(foundTips))
		for _, idx := range adaptors {
			s = append(s, fmt.Sprintf("head %d %s", idx, summariseChannels(foundTips[idx])))
		}
		return errors.Errorf("%s loaded on %s", pTips(numTips), strings.Join(s, " and "))
	}

	return nil
}

//assertNoCollisionsInGroup check that there are no collisions, ignoring the specified channels on the given adaptor
func assertNoCollisionsInGroup(settings *SimulatorSettings, adaptor *AdaptorState, channelsToIgnore []int, channelClearance float64) *CollisionError {

	var maxChannels int
	for _, ad := range adaptor.GetGroup().GetAdaptors() {
		if c := ad.GetChannelCount(); c > maxChannels {
			maxChannels = c
		}
	}
	ignore := make([]bool, maxChannels)
	for _, ch := range channelsToIgnore {
		ignore[ch] = true
	}

	adaptors := make([]int, 0, adaptor.GetGroup().NumAdaptors())
	channelMap := make(map[int][]int)
	objectMap := make(map[wtype.LHObject]bool)
	for _, ad := range adaptor.GetGroup().GetAdaptors() {
		channels := make([]int, 0, ad.GetChannelCount())
		for i := 0; i < ad.GetChannelCount(); i++ {
			if ad == adaptor && ignore[i] {
				continue
			}
			objects := ad.GetChannel(i).GetCollisions(settings, channelClearance)
			for _, o := range objects {
				objectMap[o] = true
			}
			if len(objects) > 0 {
				channels = append(channels, i)
			}
		}
		if len(channels) > 0 {
			idx := ad.GetIndex()
			adaptors = append(adaptors, idx)
			channelMap[idx] = channels
		}
	}

	//no collisions
	if len(adaptors) == 0 {
		return nil
	}

	robotState := adaptor.GetGroup().GetRobot()

	uniqueObjects := make([]wtype.LHObject, 0, len(objectMap))
	for obj := range objectMap {
		uniqueObjects = append(uniqueObjects, obj)
	}
	return NewCollisionError(robotState, channelMap, uniqueObjects)
}

var pluralMap = map[string]string{
	"well":     "wells",
	"tip":      "tips",
	"plate":    "plates",
	"tipbox":   "tipboxes",
	"tipwaste": "tipwastes",
}

//pluralise the things we care about
func pluralClassOf(o interface{}, num int) string {
	r := wtype.ClassOf(o)
	if num == 1 {
		return r
	}
	if p, ok := pluralMap[r]; ok {
		return p
	}
	return r
}

func coordsMatch(tc [][]wtype.WellCoords, wc []wtype.WellCoords) bool {
	if len(tc) != 1 {
		return false
	}

	wc2 := make([]wtype.WellCoords, 0, len(wc))
	for _, well := range wc {
		if !well.IsZero() {
			wc2 = append(wc2, well)
		}
	}
	if len(tc[0]) != len(wc2) {
		return false
	}

	for i := 0; i < len(wc2); i++ {
		if !tc[0][i].Equals(wc2[i]) {
			return false
		}
	}

	return true
}

func pTips(N int) string {
	if N == 1 {
		return "tip"
	}
	return "tips"
}

func pWells(N int) string {
	if N == 1 {
		return "well"
	}
	return "wells"
}

func pLengths(N int) string {
	if N == 1 {
		return "length"
	}
	return "lengths"
}

func intsContiguous(lhs, rhs int) bool {
	return rhs-lhs == 1
}

func appendContiguous(s []string, channels []int, start, length int) []string {
	if length == 1 {
		return append(s, fmt.Sprintf("%d", channels[start]))
	}
	return append(s, fmt.Sprintf("%d-%d", channels[start], channels[start+length-1]))
}

func summariseChannels(channels []int) string {
	if len(channels) == 0 {
		return "no channels"
	}
	if len(channels) == 1 {
		return fmt.Sprintf("channel %d", channels[0])
	}
	sch := make([]string, 0, len(channels))
	start, length := 0, 1
	for start+length < len(channels) {
		if intsContiguous(channels[start+length-1], channels[start+length]) {
			length += 1
		} else {
			sch = appendContiguous(sch, channels, start, length)
			start = start + length
			length = 1
		}
	}
	sch = appendContiguous(sch, channels, start, length)

	return fmt.Sprintf("channels %s", strings.Join(sch, ","))
}

func summariseVolumes(vols []float64) string {
	equal := true
	for _, v := range vols {
		if v != vols[0] {
			equal = false
			break
		}
	}

	if equal {
		return wunit.NewVolume(vols[0], "ul").ToString()
	}

	s_vols := make([]string, len(vols))
	for i, v := range vols {
		s_vols[i] = wunit.NewVolume(v, "ul").ToString()
		s_vols[i] = s_vols[i][:len(s_vols[i])-3]
	}
	return fmt.Sprintf("{%s} ul", strings.Join(s_vols, ","))
}

func summariseRates(rates []wunit.FlowRate) string {
	asString := make([]string, 0, len(rates))
	for _, r := range rates {
		asString = append(asString, r.ToString())
	}
	return summariseStrings(asString)
}

func summariseStrings(s []string) string {
	if countUnique(s, true) == 1 {
		return firstNonEmpty(s)
	}
	return "{" + strings.Join(getUnique(s, true), ",") + "}"
}

func summariseWellReferences(channels []int, offsetZ []float64, references []wtype.WellReference) string {
	o := make([]string, 0, len(channels))
	for _, i := range channels {
		offset := offsetZ[i]
		direction := "above"
		if offset < 0 {
			direction = "below"
			offset = -offset
		}
		offsetU := wunit.NewLength(offset, "mm")
		if offsetU.IsZero() {
			continue
		}
		o = append(o, fmt.Sprintf("%v %s", offsetU, direction))
	}

	s := make([]string, 0, len(channels))
	for _, i := range channels {
		s = append(s, references[i].String())
	}

	if len(o) > 0 {
		return summariseStrings(o) + " " + summariseStrings(s)
	}
	return summariseStrings(s)
}

func summariseCycles(cycles []int, elems []int) string {
	if iElemsEqual(cycles, elems) {
		if cycles[0] == 1 {
			return "once"
		} else {
			return fmt.Sprintf("%d times", cycles[0])
		}
	}
	sc := make([]string, 0, len(elems))
	for _, i := range elems {
		sc = append(sc, fmt.Sprintf("%d", cycles[i]))
	}
	return fmt.Sprintf("{%s} times", strings.Join(sc, ","))
}

func summariseWells(wells []*wtype.LHWell, elems []int) string {
	w := make([]string, 0, len(elems))
	for _, i := range elems {
		w = append(w, wells[i].GetWellCoords().FormatA1())
	}
	uw := getUnique(w, true)

	if len(uw) == 1 {
		return fmt.Sprintf("well %s", uw[0])
	}
	return fmt.Sprintf("wells %s", strings.Join(uw, ","))
}

func summarisePlates(wells []*wtype.LHWell, elems []int) string {
	p := make([]string, 0, len(elems))
	for _, i := range elems {
		if wells[i] != nil {
			p = append(p, wtype.NameOf(wells[i].Plate))
		}
	}
	up := getUnique(p, true)

	if len(up) == 1 {
		return fmt.Sprintf("plate \"%s\"", up[0])
	}
	return fmt.Sprintf("plates \"%s\"", strings.Join(up, "\",\""))

}

//summarisePlateWells list wells for each plate preserving order
func summarisePlateWells(wells []*wtype.LHWell, elems []int) string {
	var lastWell *wtype.LHWell
	currentChunk := make([]wtype.WellCoords, 0, len(elems))
	var chunkedWells [][]wtype.WellCoords
	var plateNames []string

	for _, i := range elems {
		well := wells[i]
		if lastWell != nil && lastWell.GetParent() != well.GetParent() {
			chunkedWells = append(chunkedWells, currentChunk)
			currentChunk = make([]wtype.WellCoords, 0, len(elems))
			plateNames = append(plateNames, wtype.NameOf(well.GetParent()))
		}
		lastWell = well
		if well != nil {
			currentChunk = append(currentChunk, well.GetWellCoords())
		}
	}
	chunkedWells = append(chunkedWells, currentChunk)
	plateNames = append(plateNames, wtype.NameOf(lastWell.GetParent()))

	var ret []string
	for i, name := range plateNames {
		ret = append(ret, fmt.Sprintf("%s@%s", wtype.HumanizeWellCoords(chunkedWells[i]), name))
	}

	if len(ret) == 0 {
		return "nil"
	}

	return strings.Join(ret, ", ")
}

func iElemsEqual(sl []int, elems []int) bool {
	for _, i := range elems {
		if sl[i] != sl[elems[0]] {
			return false
		}
	}
	return true
}

func fElemsEqual(sl []float64, elems []int) bool {
	for _, i := range elems {
		if sl[i] != sl[elems[0]] {
			return false
		}
	}
	return true
}

func extend_ints(l int, sl []int) []int {
	if len(sl) < l {
		r := make([]int, l)
		copy(r, sl)
		return r
	}
	return sl
}

func extend_floats(l int, sl []float64) []float64 {
	if len(sl) < l {
		r := make([]float64, l)
		copy(r, sl)
		return r
	}
	return sl
}

func extend_strings(l int, sl []string) []string {
	if len(sl) < l {
		r := make([]string, l)
		copy(r, sl)
		return r
	}
	return sl
}

func extend_bools(l int, sl []bool) []bool {
	if len(sl) < l {
		r := make([]bool, l)
		copy(r, sl)
		return r
	}
	return sl
}
