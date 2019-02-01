package wtype

import (
	"fmt"
	"strings"

	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/laboratory/effects/id"
)

type ComponentVector []*Liquid

func (cv ComponentVector) String() string {
	s := fmt.Sprintf("%v %v %v", cv.GetNames(), cv.GetVols(), cv.GetWellCoords())
	return s
}

//DeleteAllBelowVolume set all components whose volume is below vol to nil
func (cv ComponentVector) DeleteAllBelowVolume(vol wunit.Volume) {
	for i := 0; i < len(cv); i++ {
		//Volume.isZero() checks that volume is zero or within a small tolerace to zero
		if v := cv[i].Volume(); v.LessThan(vol) && !wunit.SubtractVolumes(vol, v).IsZero() {
			cv[i] = nil
		}
	}
}

func (cv ComponentVector) SubtractVolume(vol wunit.Volume) {
	for i := 0; i < len(cv); i++ {
		v := cv[i].Volume()
		v.Subtract(vol)
		if !v.IsPositive() {
			cv[i] = nil
		} else {
			cv[i].SetVolume(v)
		}
	}
}

func (cv ComponentVector) Dup(idGen *id.IDGenerator) ComponentVector {
	ret := make(ComponentVector, len(cv))

	for i, v := range cv {
		if v == nil {
			continue
		}
		ret[i] = v.Dup(idGen)
	}
	return ret
}

func (cv ComponentVector) GetNames() []string {
	sa := make([]string, len(cv))

	for i := 0; i < len(cv); i++ {
		if cv[i] != nil {
			sa[i] = cv[i].FullyQualifiedName()
		}
	}

	return sa
}

func (cv ComponentVector) GetVols() []wunit.Volume {
	ret := make([]wunit.Volume, len(cv))

	for i, c := range cv {
		if c == nil {
			ret[i] = wunit.ZeroVolume()
		} else {
			ret[i] = wunit.NewVolume(c.Vol, c.Vunit)
		}
	}

	return ret
}

func (cv ComponentVector) Empty() bool {
	for _, c := range cv {
		if c != nil && c.Vol > 0.0 {
			return false
		}
	}

	return true
}

func (cv ComponentVector) ToSumHash() map[string]wunit.Volume {
	ret := make(map[string]wunit.Volume, len(cv))

	for _, c := range cv {
		// skip nil components
		if c == nil {
			continue
		}

		if c.CName == "" {
			continue
		}

		v, ok := ret[c.FullyQualifiedName()]

		if !ok {
			v = wunit.ZeroVolume()
			ret[c.FullyQualifiedName()] = v
		}
		v.Add(c.Volume())
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
		if c == nil {
			continue
		}
		tx := strings.Split(c.Loc, ":")

		if len(tx) <= x {
			ret[i] = ""
			continue
		}

		ret[i] = tx[x]
	}

	return ret
}

func (cv1 ComponentVector) Equal(cv2 ComponentVector) bool {
	if len(cv1) != len(cv2) {
		return false
	}

	for i := 0; i < len(cv1); i++ {
		if cv1[i] != nil && cv2[i] != nil {
			if !cv1[i].EqualTypeVolumeID(cv2[i]) {
				return false
			}
		} else if !(cv1[i] == nil && cv2[i] == nil) {
			return false
		}
	}

	return true
}
