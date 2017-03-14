package wtype

import (
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"strings"
)

type ComponentVector []*LHComponent

func (cv ComponentVector) GetVols() []wunit.Volume {
	ret := make([]wunit.Volume, len(cv))

	for i, c := range cv {
		ret[i] = wunit.NewVolume(c.Vol, c.Vunit)
	}

	return ret
}

func (cv ComponentVector) GetPlateIds() []string {
	return cv.getLocTok(0)
}

func (cv ComponentVector) GetWellCoords() []string {
	return cv.getLocTok(1)
}

func (cv ComponentVector) getLocTok(x int) []string {
	ret := make([]string, len(cv))
	for i, c := range cv {
		tx := strings.Split(c.Loc, ":")

		if len(tx) <= x {
			ret[i] = ""
			continue
		}

		ret[i] = tx[x]
	}

	return ret
}
