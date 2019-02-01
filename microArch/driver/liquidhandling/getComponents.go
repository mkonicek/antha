package liquidhandling

// func (lhp *LHProperties) GetComponents(cmps []*wtype.Liquid, carryvol wunit.Volume, ori, multi int, independent, legacyVolume bool) (plateIDs, wellCoords [][]string, vols [][]wunit.Volume, err error)

import (
	"fmt"

	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/laboratory/effects/id"
)

type ComponentVolumeHash map[string]wunit.Volume

func (h ComponentVolumeHash) AllVolsPosOrZero() bool {
	for _, v := range h {
		if v.LessThan(wunit.ZeroVolume()) && !v.IsZero() {
			return false
		}
	}
	return true
}

func (h ComponentVolumeHash) Dup() ComponentVolumeHash {
	r := make(ComponentVolumeHash, len(h))
	for k, v := range h {
		r[k] = v.Dup()
	}

	return r
}

type GetComponentsOptions struct {
	Cmps            wtype.ComponentVector
	Carryvol        wunit.Volume
	Ori             wtype.ChannelOrientation
	Multi           int
	Independent     bool
	LegacyVolume    bool
	IgnoreInstances bool // treat everything as virtual?
}

type ParallelTransfer struct {
	PlateIDs   []string
	WellCoords []string
	Vols       []wunit.Volume
}

type GetComponentsReply struct {
	Transfers []ParallelTransfer
}

func newReply() GetComponentsReply {
	return GetComponentsReply{Transfers: make([]ParallelTransfer, 0, 1)}
}

func areWeDoneYet(cmps wtype.ComponentVector) bool {
	for _, c := range cmps {
		if c != nil && c.Vol != 0 {
			return false
		}
	}

	return true
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

func (lhp *LHProperties) GetSourcesFor(idGen *id.IDGenerator, cmps wtype.ComponentVector, ori wtype.ChannelOrientation, multi int, minPossibleVolume wunit.Volume, ignoreInstances bool, carryVol wunit.Volume) []wtype.ComponentVector {
	ret := make([]wtype.ComponentVector, 0, 1)

	for _, ipref := range lhp.OrderedMergedPlatePrefs() {
		p, ok := lhp.Plates[ipref]

		if ok {
			it := getPlateIterator(p, ori, multi)
			for wv := it.Curr(); it.Valid(); wv = it.Next() {
				// cmps needs duping here
				mycmps := p.GetVolumeFilteredContentVector(idGen, wv, cmps, minPossibleVolume, ignoreInstances) // dups components
				if mycmps.Empty() {
					continue
				}

				// mycmps has incorrect volumes, try correcting them here

				correct_volumes(mycmps)

				ret = append(ret, mycmps)
			}
		}
	}

	return ret
}

func correct_volumes(cmps wtype.ComponentVector) {
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

func cullZeroes(m map[string]wunit.Volume) map[string]wunit.Volume {
	r := make(map[string]wunit.Volume, len(m))

	for k, v := range m {
		if v.IsZero() {
			continue
		}

		r[k] = v
	}

	return r
}

func sourceVolumesOK(srcs []wtype.ComponentVector, dests wtype.ComponentVector) (bool, string) {
	collSrcs := sumSources(srcs)
	collDsts := dests.ToSumHash()
	collDsts = cullZeroes(collDsts)

	result := subHash(collSrcs, collDsts)

	if len(collSrcs) < len(collDsts) {
		return false, collateDifference(collDsts, collSrcs, result)
	}

	r := result.AllVolsPosOrZero()

	if r {
		return r, ""
	} else {
		return r, collateDifference(collDsts, collSrcs, result)
	}
}

func collateDifference(a, b, c map[string]wunit.Volume) string {
	s := ""

	for k := range a {
		_, ok := b[k]

		if !ok {
			s += fmt.Sprintf("%s; ", k)
			continue
		}

		v := c[k]

		if v.RawValue() < 0.0 {
			v.MultiplyBy(-1.0)
			s += fmt.Sprintf("%s - missing %s; ", k, v.ToString())
		}
	}

	return s
}

func subHash(h1, h2 ComponentVolumeHash) ComponentVolumeHash {
	r := h1.Dup()
	for k, v := range h2 {
		_, ok := r[k]

		if ok {
			r[k].Subtract(v)
		}
	}

	return r
}

func sumSources(cmpV []wtype.ComponentVector) ComponentVolumeHash {
	ret := make(ComponentVolumeHash, len(cmpV))
	for _, cV2 := range cmpV {
		for _, c := range cV2 {
			if c != nil && c.CName != "" {
				v, ok := ret[c.FullyQualifiedName()]
				if !ok {
					v = wunit.NewVolume(0.0, "ul")
					ret[c.FullyQualifiedName()] = v
				}
				v.Add(c.Volume())
			}
		}
	}

	return ret
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

func (lhp *LHProperties) GetComponents(idGen *id.IDGenerator, opt GetComponentsOptions) (GetComponentsReply, error) {
	rep := newReply()
	// build list of possible sources -- this is a list of ComponentVectors

	srcs := lhp.GetSourcesFor(idGen, opt.Cmps, opt.Ori, opt.Multi, lhp.MinPossibleVolume(), opt.IgnoreInstances, opt.Carryvol)

	// keep taking chunks until either we get everything or run out
	// optimization options apply here as parameters for the next level down

	currCmps := opt.Cmps.Dup(idGen)
	var lastCmps wtype.ComponentVector

	var done bool

	for {
		done = areWeDoneYet(currCmps)
		if done {
			break
		}

		if ok, s := sourceVolumesOK(srcs, currCmps); !ok {

			if opt.IgnoreInstances {
				return GetComponentsReply{}, fmt.Errorf("Insufficient source volumes for components %s", s)
			} else {
				opt.IgnoreInstances = true
				return lhp.GetComponents(idGen, opt)
			}
		}

		if cmpVecsEqual(lastCmps, currCmps) {
			// if we are here we should be able to service the request but not
			// as-is...
			break
		}

		bestMatch := wtype.Match{Sc: -1.0}
		var bestSrc wtype.ComponentVector
		// srcs is chunked up to conform to what can be accessed by the LH
		for _, src := range srcs {
			if src.Empty() {
				continue
			}

			match, err := wtype.MatchComponents(idGen, currCmps, src, opt.Independent, false)

			if err != nil && err.Error() != wtype.NotFoundError {
				return rep, err
			}

			if match.Sc > bestMatch.Sc {
				bestMatch = match
				bestSrc = src
			}
		}

		if bestMatch.Sc == -1 {
			return rep, fmt.Errorf("Components %s %s -- try increasing source volumes, if this does not work or is not possible please report to the authors\n", currCmps.String(), wtype.NotFoundError)
		}

		// adjust finally to ensure we don't leave too little

		bestMatch = makeMatchSafe(currCmps, bestMatch, lhp.MinPossibleVolume())

		// update sources
		updateSources(bestSrc, bestMatch, opt.Carryvol, lhp.MinPossibleVolume())
		lastCmps = currCmps.Dup(idGen)
		updateDests(currCmps, bestMatch)
		rep.Transfers = append(rep.Transfers, matchToParallelTransfer(bestMatch))
	}

	return rep, nil
}

func updateSources(src wtype.ComponentVector, match wtype.Match, carryVol, minPossibleVolume wunit.Volume) wtype.ComponentVector {
	for i := 0; i < len(match.M); i++ {
		if match.M[i] != -1 {
			volSub := wunit.CopyVolume(match.Vols[i])
			volSub.Add(carryVol)
			src[match.M[i]].Vol -= volSub.ConvertToString(src[match.M[i]].Vunit)
		}
	}

	src.DeleteAllBelowVolume(minPossibleVolume)

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
