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

func infoForWell(w *wtype.LHWell, p *wtype.Plate, idGen *id.IDGenerator) (wellInfo, error) {
	volUl, err := w.CurrentVolume(idGen).InStringUnit("ul")
	if err != nil {
		return wellInfo{}, err
	}

	return wellInfo{
		PlateName: p.Name(),
		PlateType: wtype.TypeOf(w.Plate),
		VolumeUl:  volUl.ToString(),
	}, nil
}

// Plates checks that two sets of plates are equivalent
func Plates(expected, got map[string]*wtype.Plate, idGen *id.IDGenerator) error {
	errs := make(utils.ErrorSlice, 0, len(expected)+len(got))

	eSeen := make(map[string]struct{})

	for pos, ep := range expected {
		eSeen[pos] = struct{}{}
		if gp, ok := got[pos]; !ok {
			errs = append(errs, fmt.Errorf("expected to find plate at position '%s' but was not in generated output", pos))
		} else if err := comparePlate(ep, gp, idGen); err != nil {
			errs = append(errs, err)
		}
	}

	for pos := range got {
		if _, ok := eSeen[pos]; !ok {
			errs = append(errs, fmt.Errorf("found unexpected plate at position '%s' in output", pos))
		}
	}

	return errs.Pack()
}

func plateWellInfo(p *wtype.Plate, idGen *id.IDGenerator) (map[wellInfo]int, error) {
	pwi := make(map[wellInfo]int)
	for _, col := range p.Cols {
		for _, w := range col {
			if wi, err := infoForWell(w, p, idGen); err != nil {
				return nil, err
			} else if v, ok := pwi[wi]; ok {
				pwi[wi] = v + 1
			} else {
				pwi[wi] = 0
			}
		}
	}
	return pwi, nil
}

func comparePlate(expected, got *wtype.Plate, idGen *id.IDGenerator) error {
	ewi, err := plateWellInfo(expected, idGen)
	if err != nil {
		return err
	}

	gwi, err := plateWellInfo(got, idGen)
	if err != nil {
		return err
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
			errs = append(errs, fmt.Errorf("Saw %d of unexpect well %v", v, w))
		}
	}

	return errs.Pack()
}
