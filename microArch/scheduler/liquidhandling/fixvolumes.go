package liquidhandling

import (
	"github.com/antha-lang/antha/antha/anthalib/wunit"
)

func FixVolumes(request *LHRequest) (*LHRequest, error) {
	// we go up through the chain
	// first find the end
	request.InstructionChain.Print()
	wantedVolumes := make(map[string]wunit.Volume)
	for chainEnd := findChainEnd(request.InstructionChain); chainEnd != nil; chainEnd = chainEnd.Parent {
		stageVolumes, err := findUpdateInstructionVolumes(chainEnd, wantedVolumes)

		if err != nil {
			return request, err
		}
		wantedVolumes = stageVolumes
	}

	return request, nil
}

func findUpdateInstructionVolumes(ch *IChain, wanted map[string]wunit.Volume) (map[string]wunit.Volume, error) {
	newWanted := make(map[string]wunit.Volume)
	for _, ins := range ch.Values {
		wantVol, ok := wanted[ins.Result.FullyQualifiedName()]

		if ok && wantVol.GreaterThan(ins.Result.Volume()) {
			r := wantVol.RawValue() / ins.Result.Volume().ConvertTo(wantVol.Unit())
			ins.AdjustVolumesBy(r)
			delete(wanted, ins.Result.FullyQualifiedName())
		}

		newWanted = mapAdd(newWanted, ins.InputVolumeMap())
	}

	return newWanted, nil
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
