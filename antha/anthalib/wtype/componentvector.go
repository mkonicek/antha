package wtype

import (
	"fmt"
	"strings"

	"github.com/antha-lang/antha/antha/anthalib/wunit"
)

type ComponentVector []*LHComponent

func (cv ComponentVector) String() string {
	s := fmt.Sprintf("%v %v %v", cv.GetNames(), cv.GetVols(), cv.GetWellCoords())
	return s
}

func (cv ComponentVector) Dup() ComponentVector {
	ret := make(ComponentVector, len(cv))

	for i, v := range cv {
		if v == nil {
			continue
		}
		ret[i] = v.Dup()
	}
	return cv
}

func (cv ComponentVector) GetNames() []string {
	sa := make([]string, len(cv))

	for i := 0; i < len(cv); i++ {
		if cv[i] != nil {
			sa[i] = cv[i].CName
		}
	}

	return sa
}

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
