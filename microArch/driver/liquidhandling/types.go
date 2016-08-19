package liquidhandling

import (
	"github.com/antha-lang/antha/anthalib/wtype"
)

type LHIVector []*wtype.LHInstruction

func (lhiv LHIVector) MaxLen() int {
	l := 0
	for _, i := range lhiv {
		ll := len(i.Components)

		if ll > l {
			l = ll
		}
	}

	return l
}

func (lhiv LHIvector) CompsAt(i int) []*wtype.LHComponent {
	ret := make([]*wtype.LHComponent, len(lhiv))

	for ci := 0; ci < len(lhiv); ci++ {

		cmp := lhiv[ci]

		if ci == 0 && cmp.IsMixInPlace() {
			continue
		}

		if i >= len(cmp.Components) {
			continue
		}
		ret[ci] = cmp.Components[i]
	}

	return ret
}
