package liquidhandling

import (
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"strings"
)

func FixVolumes(request *LHRequest) (*LHRequest, error) {
	// we go up through the chain
	// first find the end
	wantedVolumes := make(map[string]wunit.Volume)
	c := 0
	for chainEnd := findChainEnd(request.InstructionChain); chainEnd != nil; chainEnd = chainEnd.Parent {
		if len(chainEnd.Values) > 0 && chainEnd.Values[0].Type == wtype.LHIMIX {
			stageVolumes, err := findUpdateInstructionVolumes(chainEnd, wantedVolumes, request.MergedInputOutputPlates())

			if err != nil {
				return request, err
			}
			wantedVolumes = stageVolumes
		} else if chainEnd.Values[0].Type == wtype.LHIPRM {
			// update the wanted volumes to the correct names
			wantedVolumes = passThrough(chainEnd.Values, wantedVolumes)
		}
		c += 1
	}

	return request, nil
}

func passThrough(values []*wtype.LHInstruction, wanted map[string]wunit.Volume) map[string]wunit.Volume {
	ret := make(map[string]wunit.Volume, len(wanted))

	for _, v := range values {
		ret = passThroughMap(v, wanted, ret)
	}

	return ret
}

func passThroughMap(ins *wtype.LHInstruction, wanted, updated map[string]wunit.Volume) map[string]wunit.Volume {
	if ins.Type == wtype.LHIPRM {
		for IDin, out := range ins.PassThrough {
			IDout := out.ID

			vol, ok := getWantVol(wanted, out.FullyQualifiedName())

			if ok {
				newName := strings.Replace(out.FullyQualifiedName(), IDout, IDin, -1)
				if wantInPlace(wanted, out.FullyQualifiedName()) {
					newName = newName + wtype.InPlaceMarker
				}
				updated[newName] = vol
				deleteWantOf(wanted, out.FullyQualifiedName())
			}
		}

	}
	// merge in remaining keys from wanted

	for k, v := range wanted {
		updated[k] = v
	}

	return updated
}

func findUpdateInstructionVolumes(ch *IChain, wanted map[string]wunit.Volume, plates map[string]*wtype.LHPlate) (map[string]wunit.Volume, error) {
	newWanted := make(map[string]wunit.Volume)
	for _, ins := range ch.Values {
		//wantVol, ok := wanted[ins.Result.FullyQualifiedName()]

		wantVol, ok := getWantVol(wanted, ins.Result.FullyQualifiedName())

		if ok {
			_, reallyOK := plates[ins.PlateID]

			if !reallyOK {
				if ins.PlateID != "" {
					panic("Cannot fix volume for plate ID without corresponding type")
				}
			} else if !wantInPlace(wanted, ins.Result.FullyQualifiedName()) {
				wantVol.Add(plates[ins.PlateID].Rows[0][0].ResidualVolume())
			}

			if wantVol.GreaterThan(ins.Result.Volume()) {
				r := wantVol.RawValue() / ins.Result.Volume().ConvertTo(wantVol.Unit())
				ins.AdjustVolumesBy(r)

				//delete(wanted, ins.Result.FullyQualifiedName())
				deleteWantOf(wanted, ins.Result.FullyQualifiedName())
			}
		}

		newWanted = mapAdd(newWanted, ins.InputVolumeMap(wunit.NewVolume(0.5, "ul")))
	}

	newWanted = mapAdd(wanted, newWanted)

	return newWanted, nil
}

func getWantVol(wanted map[string]wunit.Volume, key string) (wunit.Volume, bool) {
	// look for in-place markers

	v, ok := wanted[key+wtype.InPlaceMarker]

	if !ok {
		v, ok = wanted[key]
	}

	return v, ok
}

func deleteWantOf(wanted map[string]wunit.Volume, key string) {
	_, ok := wanted[key+wtype.InPlaceMarker]

	if ok {
		delete(wanted, key+wtype.InPlaceMarker)
	} else {
		delete(wanted, key)
	}
}

func wantInPlace(wanted map[string]wunit.Volume, key string) bool {
	_, ok := wanted[key+wtype.InPlaceMarker]
	return ok
}

func mapDup(m map[string]wunit.Volume) map[string]wunit.Volume {
	r := make(map[string]wunit.Volume, len(m))

	for k, v := range m {
		r[k] = v.Dup()
	}
	return r
}

func mapAdd(m1, m2 map[string]wunit.Volume) map[string]wunit.Volume {
	r := mapDup(m2)
	for k, v := range m1 {
		vv, ok := r[k]

		if ok {
			vv.Add(v)
		} else {
			r[k] = v.Dup()
		}
	}

	return r
}

func findChainEnd(ch *IChain) *IChain {
	if ch.Child == nil {
		return ch
	}

	return findChainEnd(ch.Child)
}
