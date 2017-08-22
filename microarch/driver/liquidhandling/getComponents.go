package liquidhandling

// func (lhp *LHProperties) GetComponents(cmps []*wtype.LHComponent, carryvol wunit.Volume, ori, multi int, independent, legacyVolume bool) (plateIDs, wellCoords [][]string, vols [][]wunit.Volume, err error)

import (
	"fmt"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
)

type ComponentVolumeHash map[string]wunit.Volume

func (h ComponentVolumeHash) AllVolsPosOrZero() bool {
	for _, v := range h {
		if v.LessThan(wunit.ZeroVolume()) {
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
	Cmps         wtype.ComponentVector
	Carryvol     wunit.Volume
	Ori          int
	Multi        int
	Independent  bool
	LegacyVolume bool
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
		if c.Vol != 0 {
			fmt.Println("NOT DONE YET! ", c.CName, " NEEDS ", c.Vol)
			return false
		}
	}

	return true
}

func matchToParallelTransfer(m wtype.Match) ParallelTransfer {
	return ParallelTransfer{PlateIDs: m.IDs, WellCoords: m.WCs, Vols: m.Vols}
}

// returns a vector iterator for a plate given the multichannel capabilites of the head (ori, multi)
func getPlateIterator(lhp *wtype.LHPlate, ori, multi int) wtype.VectorPlateIterator {
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

		return wtype.NewTickingColVectorIterator(lhp, multi, tpw, wpt)
	} else {
		// needs same treatment as above
		return wtype.NewRowVectorIterator(lhp, multi)
	}
}

func (lhp *LHProperties) GetSourcesFor(cmps wtype.ComponentVector, ori, multi int) []wtype.ComponentVector {
	ret := make([]wtype.ComponentVector, 0, 1)

	for _, ipref := range lhp.OrderedMergedPlatePrefs() {
		p, ok := lhp.Plates[ipref]

		if ok {
			it := getPlateIterator(p, ori, multi)
			for wv := it.Curr(); it.Valid(); wv = it.Next() {
				// cmps needs duping here
				mycmps := p.GetFilteredContentVector(wv, cmps) // dups components

				if mycmps.Empty() {
					continue
				}

				ret = append(ret, mycmps)
			}
		}
	}

	return ret
}

func sourceVolumesOK(srcs []wtype.ComponentVector, dests wtype.ComponentVector) bool {
	collSrcs := sumSources(srcs)
	collDsts := dests.ToSumHash()
	result := subHash(collSrcs, collDsts)

	return result.AllVolsPosOrZero()
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
				v, ok := ret[c.CName]
				if !ok {
					v = wunit.NewVolume(0.0, "ul")
					ret[c.CName] = v
				}
				v.Add(c.Volume())
			}
		}
	}

	return ret
}

func (lhp *LHProperties) GetComponents(opt GetComponentsOptions) (GetComponentsReply, error) {
	rep := newReply()
	// build list of possible sources -- this is simply a ComponentVector of all the possible sources

	srcs := lhp.GetSourcesFor(opt.Cmps, opt.Ori, opt.Multi)

	// keep taking chunks until either we get everything or run out
	// optimization options apply here as parameters for the next level down

	currCmps := opt.Cmps.Dup()
	done := false

	// we can't take more than 1 transfer each at this stage
	stuckCounter := len(currCmps)

	for {
		if stuckCounter == 0 {
			return GetComponentsReply{}, fmt.Errorf("GOT STUCK")
		}
		done = areWeDoneYet(currCmps)
		if done {
			break
		}

		if !sourceVolumesOK(srcs, currCmps) {
			return GetComponentsReply{}, fmt.Errorf("Insufficient source volumes")
		}

		bestMatch := wtype.Match{Sc: -1.0}
		var bestSrc wtype.ComponentVector

		// srcs is chunked up to conform to what can be accessed by the LH
		for _, src := range srcs {
			match, err := wtype.MatchComponents(currCmps, src, opt.Independent, true)

			if err.Error() != wtype.NotFoundError {
				fmt.Println("BREAKING EARLY")
				return rep, err
			}

			fmt.Println("NEW MATCH: ", match)

			if match.Sc > bestMatch.Sc {
				bestMatch = match
				bestSrc = src
			}
		}

		// update sources

		updateSources(bestSrc, bestMatch)
		updateDests(currCmps, bestMatch)

		rep.Transfers = append(rep.Transfers, matchToParallelTransfer(bestMatch))

		stuckCounter -= 1
	}

	return rep, nil
}

func updateSources(src wtype.ComponentVector, match wtype.Match) wtype.ComponentVector {
	for i := 0; i < len(match.M); i++ {
		if match.M[i] != -1 {
			src[match.M[i]].Vol -= match.Vols[i].ConvertToString(src[match.M[i]].Vunit)
		}
	}

	return src
}

func updateDests(dst wtype.ComponentVector, match wtype.Match) wtype.ComponentVector {
	fmt.Println("UPDATING DESTS ", match)
	for i := 0; i < len(match.M); i++ {
		if match.M[i] != -1 {
			fmt.Println("REMOVING ", i, " SRC ", match.M[i], " ", match.Vols[i], " WAS ", dst[i].Vol)
			dst[i].Vol -= match.Vols[i].ConvertToString(dst[i].Vunit)
		}
	}

	return dst
}
