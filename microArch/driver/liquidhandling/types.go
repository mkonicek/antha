package liquidhandling

import "github.com/antha-lang/antha/antha/anthalib/wtype"

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

func (lhiv LHIVector) CompsAt(i int) []*wtype.LHComponent {
	ret := make([]*wtype.LHComponent, len(lhiv))

	for ix, ins := range lhiv {
		if i == 0 && ins.IsMixInPlace() {
			continue
		}

		if i >= len(ins.Components) {
			continue
		}

		ret[ix] = ins.Components[i]
	}

	return ret
}
