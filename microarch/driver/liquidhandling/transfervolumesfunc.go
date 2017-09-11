package liquidhandling

import (
	"fmt"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/antha/anthalib/wutil"
	"math"
)

func TransferVolumes(Vol, Min, Max wunit.Volume) ([]wunit.Volume, error) {
	ret := make([]wunit.Volume, 0)

	vol := Vol.ConvertTo(Min.Unit())
	min := Min.RawValue()
	max := Max.RawValue()

	//	if vol < min {
	if Vol.LessThanRounded(Min, 1) {
		err := wtype.LHError(wtype.LH_ERR_VOL, fmt.Sprintf("Liquid Handler cannot service volume requested: %f - minimum volume is %f", vol, min))
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

	n, _ := math.Modf(vol / max)

	n += 1

	// should make sure of no rounding errors here... we want to
	// make sure these are within the resolution of the channel

	for i := 0; i < int(n); i++ {
		ret = append(ret, wunit.NewVolume(vol/n, Vol.Unit().PrefixedSymbol()))
	}

	return ret, nil
}

func TransferVolumesMulti(vols VolumeSet, chans []*wtype.LHChannelParameter) ([]VolumeSet, error) {
	max := func(a []int) int {
		m := a[0]
		for i := 1; i < len(a); i++ {
			if a[i] > m {
				m = a[i]
			}
		}

		return m
	}

	ks := make([]int, len(vols))

	for i := 0; i < len(vols); i++ {
		vv1 := vols[i].ConvertTo(chans[i].Maxvol.Unit())
		ks[i] = wutil.RoundInt(vv1 / chans[i].Maxvol.RawValue())
	}

	left := vols.Dup()
	ret := make([]VolumeSet, max(ks)+1)

	for i := 0; i < max(ks)+1; i++ {
		r := NewVolumeSet(len(ks))
		for j := 0; j < len(ks); j++ {
			if ks[j] > 0 {
				if left[j].LessThan(chans[i].Minvol) {
					return ret, fmt.Errorf("Insufficient remaining volume in multichannel transfer")
				}
				r[j] = chans[i].Maxvol.Dup()
				left[j].Subtract(chans[i].Maxvol)
			} else if ks[j] == 0 {
				r[j] = left[j].Dup()
			} else {
				// r[j] is left as 0.0
			}

			ks[j] -= 1
		}

		ret[i] = r
	}

	return ret, nil
}
