package liquidhandling

import (
	"fmt"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"math"
)

// TransferVolumes returns a slice of volumes V such that Min <= v <= Max and sum(V) = Vol
func TransferVolumes(Vol, Min, Max wunit.Volume) ([]wunit.Volume, error) {

	max, err := Max.InUnit(Vol.Unit())
	if err != nil {
		return nil, err
	}

	ret := make([]wunit.Volume, 0)

	//	if vol < min {
	if Vol.LessThanRounded(Min, 1) {
		err := wtype.LHError(wtype.LH_ERR_VOL, fmt.Sprintf("Liquid Handler cannot service volume requested: %v - minimum volume is %v", Vol, Min))
		return ret, err
	}

	//if vol <= max {
	if !Vol.GreaterThanRounded(Max, 1) {
		ret = append(ret, Vol)
		return ret, nil
	}

	// vol is > max, need to know by how much
	// if vol/max = n then we do n+1 equal transfers of vol / (n+1)
	// this should never be outside the range

	n, _ := math.Modf(Vol.RawValue() / max.RawValue())
	n += 1

	// should make sure of no rounding errors here... we want to
	// make sure these are within the resolution of the channel

	for i := 0; i < int(n); i++ {
		ret = append(ret, wunit.NewVolume(Vol.RawValue()/n, Vol.Unit().PrefixedSymbol()))
	}

	return ret, nil
}

// TransferVolumesMulti given a slice of volumes to transfer and channels to use,
// return an array of transfers to make such that `ret[i][j]` is the volume of the ith transfer to be made with channel j
func TransferVolumesMulti(vols VolumeSet, chans []*wtype.LHChannelParameter) ([]VolumeSet, error) {
	// aggregate vertically
	mods := make([]VolumeSet, len(vols))
	mx := 0
	for i := 0; i < len(vols); i++ {
		if chans[i] == nil {
			continue
		}
		mod, err := TransferVolumes(vols[i], chans[i].Minvol, chans[i].Maxvol)

		if err != nil {
			return []VolumeSet{}, err
		}

		mods[i] = mod
		if len(mod) > mx {
			mx = len(mod)
		}

	}

	ret := make([]VolumeSet, mx)

	for j := 0; j < mx; j++ {
		vs := make(VolumeSet, len(vols))

		for i := 0; i < len(vols); i++ {
			if j < len(mods[i]) {
				vs[i] = mods[i][j]
			}
		}

		ret[j] = vs
	}

	return ret, nil
}
