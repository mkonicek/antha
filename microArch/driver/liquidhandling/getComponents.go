package liquidhandling

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"

	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
)

type ParallelTransfer struct {
	PlateIDs   []string
	WellCoords []string
	Vols       []wunit.Volume
}

func matchToParallelTransfer(m wtype.Match) ParallelTransfer {
	return ParallelTransfer{PlateIDs: m.IDs, WellCoords: m.WCs, Vols: m.Vols}
}

// returns a vector iterator for a plate given the multichannel capabilites of the head (ori, multi)
func getPlateIterator(lhp *wtype.Plate, ori wtype.ChannelOrientation, multi int) wtype.AddressSliceIterator {
	if ori == wtype.LHVChannel {
		//it = NewColVectorIterator(lhp, multi)

		tpw := multi / lhp.WellsY()
		wpt := lhp.WellsY() / multi

		if tpw == 0 {
			tpw = 1
		}

		if wpt == 0 {
			wpt = 1
		}

		// fix for 6 row plates etc.
		if multi > lhp.WellsY() && tpw == 1 {
			multi = lhp.WellsY()
		}

		if multi == 1 {
			tpw = 1
			wpt = 1
		}

		return wtype.NewTickingIterator(lhp, wtype.ColumnWise, wtype.TopToBottom, wtype.LeftToRight, false, multi, wpt, tpw)
	} else {
		// needs same treatment as above
		return wtype.NewTickingIterator(lhp, wtype.RowWise, wtype.TopToBottom, wtype.LeftToRight, false, multi, 1, 1)
	}
}

// GetSourcesFor find all liquids in LHProperties which could be used to supply the requested cmps with the given configuration
func (lhp *LHProperties) GetSourcesFor(ori wtype.ChannelOrientation, multi int, maps ...func(*wtype.Liquid) *wtype.Liquid) []wtype.ComponentVector {
	ret := make([]wtype.ComponentVector, 0)

	for _, ipref := range lhp.OrderedMergedPlatePrefs() {

		if p, ok := lhp.Plates[ipref]; ok {

			it := getPlateIterator(p, ori, multi)
			for wv := it.Curr(); it.Valid(); wv = it.Next() {

				available := p.AvailableContents(wv).Map(maps...)

				if !available.IsEmpty() {

					// found has incorrect volumes when more than one tip is in each well, try correcting them here
					correctMultiTipsPerWell(available)

					ret = append(ret, available)
				}
			}
		}
	}

	return ret
}

// correctMultiTipsPerWell when there's more than one tip in a well, share the available volume in each tip equally between them
func correctMultiTipsPerWell(cmps wtype.ComponentVector) {
	nW := make(map[string]int)
	for _, c := range cmps {
		if c == nil {
			continue
		}

		_, ok := nW[c.Loc]
		if !ok {
			nW[c.Loc] = 0
		}
		nW[c.Loc] += 1
	}

	for _, c := range cmps {
		if c == nil {
			continue
		}

		c.Vol /= float64(nW[c.Loc])
	}
}

type volumeCheck struct {
	Available wtype.ComponentVector
	Requested wtype.ComponentVector
}

type volumesByLiquidName map[string]*volumeCheck

// Summary a nice user-friendly string summary of the remaining volumes
func (vbln volumesByLiquidName) Summary(lhp *LHProperties) string {
	cVectorSummary := func(cv wtype.ComponentVector) string {
		if len(cv) == 0 {
			return "nothing"
		} else if len(cv) == 1 {
			return cv[0].Volume().String()
		} else {
			return fmt.Sprintf("%v %s total", cv.GetVols(), cv.Volume())
		}
	}

	ret := make([]string, 0, 4*len(vbln))
	for _, check := range vbln {

		idMap := make(map[string]bool, len(check.Requested))
		for _, l := range check.Requested {
			idMap[l.ID] = true
		}
		ids := make([]string, 0, len(idMap))
		for id := range idMap {
			ids = append(ids, id)
		}

		nameMap := make(map[string]bool, len(check.Requested))
		for _, l := range check.Requested {
			nameMap[l.CName] = true
		}
		names := make([]string, 0, len(idMap))
		for name := range nameMap {
			names = append(names, name)
		}

		var extraAvailable []string
		if len(check.Available) == 0 {
			// we found nothing that matched
			// a common problem is mixing a liquid whose ID has changed for some reason
			// here we see if we can find some liquids which match name _only_, in the hope
			// that this was their intended target
			// element writers can then investigate why ID matching has failed and hopefully
			// discover the root cause of the issue
			sources := lhp.GetSourcesFor(wtype.LHVChannel, 1, func(l *wtype.Liquid) *wtype.Liquid {
				if l != nil && nameMap[l.CName] {
					return l
				}
				return nil
			})

			// filter out sources which are all nils
			other := make([]wtype.ComponentVector, 0, len(sources))
			for _, s := range sources {
				if filtered := s.Filter(func(l *wtype.Liquid) bool { return l != nil }); len(filtered) > 0 {
					other = append(other, filtered)
				}
			}

			if len(other) > 0 {
				extraAvailable = make([]string, 0, len(other)+1)
				extraAvailable = append(extraAvailable, ", but did find potential source(s) with matching name:")
				for _, o := range other {
					extraAvailable = append(extraAvailable, fmt.Sprintf(`        "%s"`, strings.Join(o.GetNames(), `", "`)))
				}
			}
		}

		// name:
		//   available: [v1, ..., v2] v3 total
		//   requested: [v4, ..., v5] v6 total
		//   shortfall: v6 - v3
		ret = append(ret, fmt.Sprintf("  %s [\"%s\"]:\n    available: %s%s\n    requested: %s\n    shortfall: %s",
			strings.Join(ids, `, `),
			strings.Join(names, `", "`),
			cVectorSummary(check.Available),
			strings.Join(extraAvailable, "\n"),
			cVectorSummary(check.Requested),
			wunit.SubtractVolumes(check.Requested.Volume(), check.Available.Volume())))
	}

	return strings.Join(ret, "\n")
}

// findInsufficientSources check that there's enough volume in found to satisfy all the volumes in requested
func findInsufficientSources(found []wtype.ComponentVector, requested wtype.ComponentVector) volumesByLiquidName {
	ret := make(volumesByLiquidName)

	// add everything we found
	for _, cv := range found {
		for _, liquid := range cv {
			if liquid == nil {
				continue
			}
			vc, ok := ret[liquid.IDOrName()]
			if !ok {
				vc = &volumeCheck{}
				ret[liquid.IDOrName()] = vc
			}
			vc.Available = append(vc.Available, liquid)
		}
	}

	// add everything we wanted
	for _, liquid := range requested {
		if liquid == nil {
			continue
		}
		vc, ok := ret[liquid.IDOrName()]
		if !ok {
			vc = &volumeCheck{}
			ret[liquid.IDOrName()] = vc
		}
		vc.Requested = append(vc.Requested, liquid)
	}

	// remove everything where there's enough
	insufficient := make(volumesByLiquidName, len(ret))
	for name, check := range ret {
		if check.Available.Volume().LessThan(check.Requested.Volume().MinusEpsilon()) {
			insufficient[name] = check
		}
	}

	return insufficient
}

func cmpVecsEqual(v1, v2 wtype.ComponentVector) bool {
	if len(v1) != len(v2) {
		return false
	}

	for i := 0; i < len(v1); i++ {
		if !cmpsEqual(v1[i], v2[i]) {
			return false
		}
	}

	return true
}

func cmpsEqual(c1, c2 *wtype.Liquid) bool {
	return c1.ID == c2.ID && c1.Vol == c2.Vol
}

func (lhp *LHProperties) GetComponents(cmps wtype.ComponentVector, ori wtype.ChannelOrientation, multi int, independent bool, legacyVolume bool) ([]ParallelTransfer, error) {
	ret := make([]ParallelTransfer, 0)

	// build the filters that we'll use later
	cmpNames := make(map[string]bool, len(cmps))
	for _, l := range cmps {
		if l != nil {
			cmpNames[l.IDOrName()] = true
		}
	}
	nameMap := func(l *wtype.Liquid) *wtype.Liquid {
		if l == nil || !cmpNames[l.IDOrName()] {
			return nil
		}
		return l
	}

	mpvMinusE := lhp.MinPossibleVolume().MinusEpsilon()
	volumeMap := func(l *wtype.Liquid) *wtype.Liquid {
		// allow volumes which are trivially less than minPossibleVolume
		if l.Volume().LessThan(mpvMinusE) {
			return nil
		}
		return l
	}

	// build list of possible sources -- this is a list of ComponentVectors

	srcs := lhp.GetSourcesFor(ori, multi, nameMap, volumeMap)

	// keep taking chunks until either we get everything or run out
	// optimization options apply here as parameters for the next level down

	currCmps := cmps.Dup()
	var lastCmps wtype.ComponentVector

	for !currCmps.IsEmpty() {

		if insufficient := findInsufficientSources(srcs, currCmps); len(insufficient) > 0 {
			return nil, errors.Errorf("found insufficient sources for %d requested liquids\n%s", len(insufficient), insufficient.Summary(lhp))
		}

		if cmpVecsEqual(lastCmps, currCmps) {
			// if we are here we should be able to service the request but not
			// as-is...
			return ret, nil
		}

		bestMatch := wtype.Match{Sc: -1.0}
		bestSrcIdx := -1
		// srcs is chunked up to conform to what can be accessed by the LH
		for i, src := range srcs {
			if !src.IsEmpty() {

				match, err := wtype.MatchComponents(currCmps, src, independent, false)

				if err != nil && !wtype.IsNotFound(err) {
					return ret, err
				}

				if match.Sc > bestMatch.Sc {
					bestMatch = match
					bestSrcIdx = i
				}
			}
		}

		if bestMatch.Sc == -1 {
			return ret, errors.WithMessage(wtype.NotFoundError, fmt.Sprintf("components %s -- try increasing source volumes", currCmps.String()))
		}

		// adjust finally to ensure we don't leave too little

		bestMatch = makeMatchSafe(currCmps, bestMatch, lhp.MinPossibleVolume())

		// update sources
		srcs[bestSrcIdx] = updateSources(srcs[bestSrcIdx], bestMatch, lhp.CarryVolume()).Map(volumeMap)
		lastCmps = currCmps.Dup()
		updateDests(currCmps, bestMatch)
		ret = append(ret, matchToParallelTransfer(bestMatch))
	}

	return ret, nil
}

func updateSources(src wtype.ComponentVector, match wtype.Match, carryVol wunit.Volume) wtype.ComponentVector {
	for i := 0; i < len(match.M); i++ {
		if match.M[i] != -1 {
			volSub := wunit.CopyVolume(match.Vols[i])
			volSub.Add(carryVol)
			src[match.M[i]].Vol -= volSub.ConvertToString(src[match.M[i]].Vunit)
		}
	}

	return src
}

func makeMatchSafe(dst wtype.ComponentVector, match wtype.Match, mpv wunit.Volume) wtype.Match {
	for i := 0; i < len(match.M); i++ {
		if match.M[i] != -1 {
			checkVol := dst[i].Vol

			checkVol -= match.Vols[i].ConvertToString(dst[i].Vunit)

			if checkVol > 0.0001 && checkVol < mpv.ConvertToString(dst[i].Vunit) {
				mpv.Subtract(wunit.NewVolume(checkVol, dst[i].Vunit))
				match.Vols[i].Subtract(mpv)

				if match.Vols[i].RawValue() < 0.0 {
					panic(fmt.Sprintf("Serious volume issue -- try a manual plate layout with some additional volume for %s", dst[i].CName))
				}
			}
		}
	}

	return match
}

func updateDests(dst wtype.ComponentVector, match wtype.Match) wtype.ComponentVector {
	for i := 0; i < len(match.M); i++ {
		if match.M[i] != -1 {
			dst[i].Vol -= match.Vols[i].ConvertToString(dst[i].Vunit)

			if dst[i].Volume().MustInStringUnit("ul").RawValue() < 0.0001 {
				dst[i].Vol = 0.0
			}
		}
	}

	return dst
}
