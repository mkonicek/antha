package liquidhandling

import (
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
)

type TransferParams struct {
	What       string
	PltFrom    string
	PltTo      string
	WellFrom   string
	WellTo     string
	Volume     wunit.Volume
	FPlateType string
	TPlateType string
	FVolume    wunit.Volume
	TVolume    wunit.Volume
	Channel    *wtype.LHChannelParameter
	TipType    string
}

func (tp TransferParams) ToString() string {
	return fmt.Sprintf("%s %s %s %s %s %s %s %s %s %s %s %s", tp.What, tp.PltFrom, tp.PltTo, tp.WellFrom, tp.WellTo, tp.Volume.ToString(), tp.FPlateType, tp.TPlateType, tp.FVolume.ToString(), tp.TVolume.ToString(), tp.Channel, tp.TipType)
}

func (tp TransferParams) Zero() bool {
	if tp.What == "" {
		return true
	}

	return false
}

func (tp TransferParams) Dup() TransferParams {
	return TransferParams{
		What:       tp.What,
		PltFrom:    tp.PltFrom,
		PltTo:      tp.PltTo,
		WellFrom:   tp.WellFrom,
		WellTo:     tp.WellTo,
		Volume:     tp.Volume.Dup(),
		FPlateType: tp.FPlateType,
		TPlateType: tp.TPlateType,
		FVolume:    tp.FVolume.Dup(),
		TVolume:    tp.TVolume.Dup(),
		Channel:    tp.Channel.Dup(),
		TipType:    tp.TipType,
	}
}

type MultiTransferParams struct {
	Multi     int
	Transfers []TransferParams
}

// slices through

func (mtp MultiTransferParams) What() []string {
	r := make([]string, mtp.Multi)
	for i, t := range mtp.Transfers {
		r[i] = t.What
	}

	return r
}
func (mtp MultiTransferParams) PltFrom() []string {
	r := make([]string, mtp.Multi)
	for i, t := range mtp.Transfers {
		r[i] = t.PltFrom
	}

	return r
}
func (mtp MultiTransferParams) PltTo() []string {
	r := make([]string, mtp.Multi)
	for i, t := range mtp.Transfers {
		r[i] = t.PltTo
	}

	return r
}
func (mtp MultiTransferParams) WellFrom() []string {
	r := make([]string, mtp.Multi)
	for i, t := range mtp.Transfers {
		r[i] = t.WellFrom
	}

	return r
}
func (mtp MultiTransferParams) WellTo() []string {
	r := make([]string, mtp.Multi)
	for i, t := range mtp.Transfers {
		r[i] = t.WellTo
	}

	return r
}
func (mtp MultiTransferParams) Volume() []wunit.Volume {
	r := make([]wunit.Volume, mtp.Multi)

	for i := 0; i < mtp.Multi; i++ {
		r[i] = wunit.ZeroVolume()
	}

	for i, t := range mtp.Transfers {
		r[i] = t.Volume.Dup()
	}

	return r
}
func (mtp MultiTransferParams) FPlateType() []string {
	r := make([]string, mtp.Multi)
	for i, t := range mtp.Transfers {
		r[i] = t.FPlateType
	}

	return r
}
func (mtp MultiTransferParams) TPlateType() []string {
	r := make([]string, mtp.Multi)
	for i, t := range mtp.Transfers {
		r[i] = t.TPlateType
	}

	return r
}

func (mtp MultiTransferParams) FVolume() []wunit.Volume {
	r := make([]wunit.Volume, mtp.Multi)

	for i := 0; i < mtp.Multi; i++ {
		r[i] = wunit.ZeroVolume()
	}

	for i, t := range mtp.Transfers {
		r[i] = t.FVolume.Dup()
	}

	return r
}

func (mtp MultiTransferParams) TVolume() []wunit.Volume {
	r := make([]wunit.Volume, mtp.Multi)

	for i := 0; i < mtp.Multi; i++ {
		r[i] = wunit.ZeroVolume()
	}

	for i, t := range mtp.Transfers {
		r[i] = t.TVolume.Dup()
	}

	return r
}

func (mtp MultiTransferParams) Channel() []*wtype.LHChannelParameter {
	r := make([]*wtype.LHChannelParameter, mtp.Multi)
	for i, t := range mtp.Transfers {
		r[i] = t.Channel.Dup()
	}

	return r

}

func (mtp MultiTransferParams) TipType() []string {
}

func NewMultiTransferParams(multi int) MultiTransferParams {
	var v MultiTransferParams
	v.Transfers = make([]TransferParams, 0, multi)
	v.Multi = multi
	return v
}

func (mtp MultiTransferParams) ParamSet(n int) TransferParams {
	return mtp.Transfers[n]
}

func (mtp MultiTransferParams) ToString() string {
	s := ""

	for i := 0; i < mtp.Multi; i++ {
		s += mtp.ParamSet(i).ToString() + " "
	}

	return s

}

func (mtp MultiTransferParams) Dup() MultiTransferParams {
	/*
		return MultiTransferParams{
			What:       dupStringArray(mtp.What),
			PltFrom:    dupStringArray(mtp.PltFrom),
			PltTo:      dupStringArray(mtp.PltTo),
			WellFrom:   dupStringArray(mtp.WellFrom),
			WellTo:     dupStringArray(mtp.WellTo),
			Volume:     dupVolArray(mtp.Volume),
			FVolume:    dupVolArray(mtp.FVolume),
			TVolume:    dupVolArray(mtp.TVolume),
			FPlateType: dupStringArray(mtp.FPlateType),
			TPlateType: dupStringArray(mtp.TPlateType),
			TipTypes:   dupStringArray(mtp.TipTypes),
		}
	*/

	tfrs := make([]TransferParams, 0, mtp.Multi)

	for i := 0; i < mtp.Multi; i++ {
		tfrs = append(tfrs, mtp.Transfers[i].Dup())
	}

	ret := NewMultiTransferParams(mtp.Multi)

	ret.Transfers = tfrs

	return ret
}

func dupStringArray(in []string) []string {
	out := make([]string, len(in))

	for i := 0; i < len(in); i++ {
		out[i] = in[i]
	}
	return out
}
func dupIntArray(in []int) []int {
	out := make([]int, len(in))

	for i := 0; i < len(in); i++ {
		out[i] = in[i]
	}
	return out
}

func dupVolArray(in []wunit.Volume) []wunit.Volume {
	out := make([]wunit.Volume, len(in))

	for i := 0; i < len(in); i++ {
		out[i] = in[i].Dup()
	}

	return out
}
