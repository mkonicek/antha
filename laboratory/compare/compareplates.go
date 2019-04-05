package compare

import (
	"fmt"

	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/laboratory/effects/id"
	"github.com/antha-lang/antha/utils"
)

type wellInfo struct {
	PlateName string
	PlateType string
	VolumeUl  string
}

func infoForWell(idGen *id.IDGenerator, w *wtype.LHWell, p *wtype.Plate) (*wellInfo, error) {
	volUl, err := w.CurrentVolume(idGen).InStringUnit("ul")
	if err != nil {
		return nil, err
	}

	return &wellInfo{
		PlateName: p.Name(),
		PlateType: wtype.TypeOf(w.Plate),
		VolumeUl:  volUl.ToString(),
	}, nil
}

// Plates checks that two sets of plates are equivalent
func Plates(idGen *id.IDGenerator, expected, got map[string]*wtype.Plate) utils.ErrorSlice {
	errs := make(utils.ErrorSlice, 0, len(expected)+len(got))

	eSeen := make(map[string]struct{})

	for pos, ep := range expected {
		eSeen[pos] = struct{}{}
		if gp, ok := got[pos]; !ok {
			errs = append(errs, fmt.Errorf("expected to find plate at position '%s' but was not in generated output", pos))
		} else {
			errs = append(errs, comparePlate(idGen, ep, gp)...)
		}
	}

	for pos := range got {
		if _, ok := eSeen[pos]; !ok {
			errs = append(errs, fmt.Errorf("found unexpected plate at position '%s' in output", pos))
		}
	}

	return errs
}

func plateWellInfo(idGen *id.IDGenerator, p *wtype.Plate) (map[wellInfo]int, error) {
	pwi := make(map[wellInfo]int)
	for _, col := range p.Cols {
		for _, w := range col {
			if wi, err := infoForWell(idGen, w, p); err != nil {
				return nil, err
			} else {
				pwi[*wi] = pwi[*wi] + 1
			}
		}
	}
	return pwi, nil
}

func comparePlate(idGen *id.IDGenerator, expected, got *wtype.Plate) utils.ErrorSlice {
	ewi, err := plateWellInfo(idGen, expected)
	if err != nil {
		return utils.ErrorSlice{err}
	}

	gwi, err := plateWellInfo(idGen, got)
	if err != nil {
		return utils.ErrorSlice{err}
	}

	eSeen := make(map[wellInfo]struct{})
	errs := make(utils.ErrorSlice, 0, len(ewi)+len(gwi))

	for w, ev := range ewi {
		eSeen[w] = struct{}{}
		if gv, ok := gwi[w]; !ok {
			errs = append(errs, fmt.Errorf("Expected %d instances of well %v, but got %d", ev, w, gv))
		}
	}

	for w, v := range gwi {
		if _, ok := eSeen[w]; !ok {
			errs = append(errs, fmt.Errorf("Saw %d instances of unexpected well %v", v, w))
		}
	}

	return errs
}
