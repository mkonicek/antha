package liquidhandling

import (
	"strings"

	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
)

// FixVolumes adjusts volumes of components in instructions in order to ensure
// sufficient quantities for all known uses. It aims to account for both residual
// well volumes and carry volumes, although the latter is intrinsically inaccurate
// since at this stage we do not know how transfers will be done and therefore how
// much carry volume will be lost
// In order to do this it needs to account for the various
// conventions in component naming to distinguish virtual from non-virtual components
// and additionally needs to ensure it treats stationary components differently since
// these do not need adjusting for carry volume or residual
func FixVolumes(request *LHRequest, carryVolume wunit.Volume) error {
	// we go up through the chain
	// first find the end

	wantedVolumes := make(map[string]wunit.Volume)
	for chainEnd := findChainEnd(request.InstructionChain); chainEnd != nil; chainEnd = chainEnd.Parent {
		if len(chainEnd.Values) == 0 {
			panic("Internal Error: Empty instruction chain node")
		}

		switch chainEnd.Values[0].Type {
		case wtype.LHIMIX:
			stageVolumes, err := findUpdateInstructionVolumes(chainEnd, wantedVolumes, request.MergedInputOutputPlates(), carryVolume)

			if err != nil {
				return err
			}
			wantedVolumes = stageVolumes
		case wtype.LHISPL:
			// split
			wantedVolumes = updateIDsAfterSplit(chainEnd.Values, wantedVolumes)
		case wtype.LHIPRM:
			// update the wanted volumes to the correct names
			wantedVolumes = passThrough(chainEnd.Values, wantedVolumes)
		default:
			panic("Internal Error: Unknown Instruction type")
		}
	}

	return nil
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
		for i := range ins.Inputs {
			IDin := ins.Inputs[i].ID
			out := ins.Outputs[i]
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

// iterate through the values in this level of the chain and ensure that IDs are updated after sampling
func updateIDsAfterSplit(values []*wtype.LHInstruction, wanted map[string]wunit.Volume) map[string]wunit.Volume {
	ret := make(map[string]wunit.Volume, len(wanted))
	for i := len(values) - 1; i >= 0; i-- {
		v := values[i]
		if v.Type == wtype.LHISPL {
			updateIDAfterSplit(v, wanted, ret)
		} else {
			panic("Internal error: Split instructions must not be grouped with other instruction types")
		}
	}

	// ensure ret contains anything else not split

	for i, v := range wanted {
		ret[i] = v
	}

	return ret
}

// update IDs in this case
func updateIDAfterSplit(ins *wtype.LHInstruction, in, out map[string]wunit.Volume) {
	// splits convert their first argument into their second result
	IDin := ins.Inputs[0].ID
	cmpOut := ins.Outputs[1]
	IDout := cmpOut.ID

	vol, ok := getWantVol(in, cmpOut.FullyQualifiedName())

	if ok {
		newName := strings.Replace(cmpOut.FullyQualifiedName(), IDout, IDin, -1)
		if wantInPlace(in, cmpOut.FullyQualifiedName()) {
			newName = newName + wtype.InPlaceMarker
		}
		// either update existing want or make a new one
		_, ok := in[newName]

		if ok {
			in[newName].Add(vol)
		} else {
			out[newName] = vol
		}
		deleteWantOf(in, cmpOut.FullyQualifiedName())
	}

}

func findUpdateInstructionVolumes(ch *wtype.IChain, wanted map[string]wunit.Volume, plates map[string]*wtype.Plate, carryVol wunit.Volume) (map[string]wunit.Volume, error) {

	newWanted := make(map[string]wunit.Volume)
	for _, ins := range ch.Values {
		//wantVol, ok := wanted[ins.Outputs[0].FullyQualifiedName()]

		wantVol, ok := getWantVol(wanted, ins.Outputs[0].FullyQualifiedName())

		if ok {
			_, reallyOK := plates[ins.PlateID]

			if !reallyOK {
				if ins.PlateID != "" {
					panic("Cannot fix volume for plate ID without corresponding type")
				}
			} else if !wantInPlace(wanted, ins.Outputs[0].FullyQualifiedName()) {
				wantVol.Add(plates[ins.PlateID].Rows[0][0].ResidualVolume())
				wantVol.Subtract(carryVol)
			}

			if wantVol.GreaterThan(ins.Outputs[0].Volume()) {
				if r, err := wunit.DivideVolumes(wantVol, ins.Outputs[0].Volume()); err != nil {
					return nil, err
				} else {
					ins.AdjustVolumesBy(r)
				}

				//delete(wanted, ins.Outputs[0].FullyQualifiedName())
				deleteWantOf(wanted, ins.Outputs[0].FullyQualifiedName())
			}
		}

		newWanted = mapAdd(newWanted, ins.InputVolumeMap(carryVol))
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

func findChainEnd(ch *wtype.IChain) *wtype.IChain {
	if ch.Child == nil {
		return ch
	}

	return findChainEnd(ch.Child)
}
