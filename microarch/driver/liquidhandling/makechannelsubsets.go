package liquidhandling

import (
	"fmt"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
)

func dupSA(sa []string) []string {
	r := make([]string, len(sa))

	for i := 0; i < len(sa); i++ {
		r[i] = sa[i]
	}

	return r
}

func dupCHA(cha []*wtype.LHChannelParameter) []*wtype.LHChannelParameter {
	r := make([]*wtype.LHChannelParameter, len(cha))

	for i := 0; i < len(cha); i++ {
		r[i] = cha[i]
	}

	return r
}
func makeChannelSubsets(tiptypes []string, channels []*wtype.LHChannelParameter) ([]TipSubset, error) {
	finished := false

	ret := make([]TipSubset, 0, 1)

	tta := dupSA(tiptypes)
	cha := dupCHA(channels)

	for i := 0; i <= len(tiptypes); i++ {
		ts := ""
		var ch *wtype.LHChannelParameter
		mask := make([]bool, len(tiptypes))

		another := false
		for j := 0; j < len(tiptypes); j++ {
			if tta[j] != "" {
				// new subset
				if ts == "" {
					another = true
					ts = tta[j]
					ch = cha[j]
				}
				if ts == tta[j] && ch.Equals(cha[j]) {
					mask[j] = true
					tta[j] = ""
					cha[j] = nil
				}
			}
		}

		if !another {
			finished = true
			break
		} else {
			ret = append(ret, TipSubset{Mask: mask, TipType: ts, Channel: ch})
		}
	}

	if !finished {
		return ret, fmt.Errorf("Could not make tip subsets for %v %v", tiptypes, channels)
	}

	return ret, nil
}
