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

// Map return a new componentvector which is the result calling each function in maps sequentially on
// each element
func (cv ComponentVector) Map(maps ...func(*Liquid) *Liquid) ComponentVector {
	ret := make(ComponentVector, len(cv))

	for i := range cv {
		ret[i] = cv[i]
		for _, m := range maps {
			ret[i] = m(ret[i])
		}
	}

	return ret
}

// Filter returns the subset of the liquids in cv for which every filter returns true
// does not guarantee to call every filter on every liquid
func (cv ComponentVector) Filter(filters ...func(*Liquid) bool) ComponentVector {
	ret := make(ComponentVector, 0, len(cv))

Outer:
	for _, l := range cv {
		for _, filter := range filters {
			if !filter(l) {
				// don't bother checking the rest once one fails
				// continue the outer for loop
				continue Outer
			}
		}
		ret = append(ret, l)
	}
	return ret
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

// IsEmpty returns true if there is no volume in the component vector
func (cv ComponentVector) IsEmpty() bool {
	for _, c := range cv {
		if !c.Volume().IsZero() {
			return false
		}
	}

	return true
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

// Volume get the total volume of liquid in the component vector
func (cv ComponentVector) Volume() wunit.Volume {
	r := wunit.ZeroVolume()
	for _, l := range cv {
		r.Add(l.Volume())
	}
	return r
}
