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
	ret := make([]*wtype.LHComponent, 0, len(lhiv))

	for _, cmp := range lhiv {
		if i == 0 && cmp.IsMixInPlace() {
			continue
		}

		if i >= len(cmp.Components) {
			continue
		}

		ret = append(ret, cmp.Components[i])
	}

	return ret
}
