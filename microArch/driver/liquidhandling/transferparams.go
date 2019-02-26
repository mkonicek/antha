package liquidhandling

import (
	"fmt"

	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
)

type TransferParams struct {
	What       string                    // liquid class
	PltFrom    string                    // from position
	PltTo      string                    // to position
	WellFrom   string                    // well coordinate in from plate
	WellTo     string                    // well coordinate in to plate
	Volume     wunit.Volume              // volume of sample being transferred
	FPlateType string                    // from plate type
	TPlateType string                    // to plate type
	FVolume    wunit.Volume              // volume in 'from' well
	TVolume    wunit.Volume              // volume in 'to' well
	Channel    *wtype.LHChannelParameter // which channel to use
	TipType    string                    // type of tip to use
	FPlateWX   int                       // X dimension in 'from' plate
	FPlateWY   int                       // Y dimension in 'from' plate
	TPlateWX   int                       // X dimension in 'to' plate
	TPlateWY   int                       // Y dimension in 'to' plate
	Component  string                    // component type
	Policy     wtype.LHPolicy            // policy attached to this transfer
}

func (tp TransferParams) ToString() string {
	return fmt.Sprintf("%s %s %s %s %s %s %s %s %s %s %s %s %d %d %d %d %s %+v", tp.What, tp.PltFrom, tp.PltTo, tp.WellFrom, tp.WellTo, tp.Volume.ToString(), tp.FPlateType, tp.TPlateType, tp.FVolume.ToString(), tp.TVolume.ToString(), tp.Channel, tp.TipType, tp.FPlateWX, tp.FPlateWY, tp.TPlateWX, tp.TPlateWY, tp.Component, tp.Policy)
}

func (tp TransferParams) IsZero() bool {
	return tp.What == "" || tp.Volume.IsZero()
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
		FPlateWX:   tp.FPlateWX,
		FPlateWY:   tp.FPlateWY,
		TPlateWX:   tp.TPlateWX,
		TPlateWY:   tp.TPlateWY,
		Component:  tp.Component,
		Policy:     tp.Policy,
	}
}

type MultiTransferParams struct {
	Multi     int
	Transfers []TransferParams
}

func (mtp MultiTransferParams) RemoveInitialBlanks() MultiTransferParams {
	r := NewMultiTransferParams(mtp.Multi)

	started := false

	for _, tp := range mtp.Transfers {
		if !tp.IsZero() {
			started = true
		}

		if started {
			r.Transfers = append(r.Transfers, tp)
		}
	}

	return r
}

func (mtp MultiTransferParams) RemoveBlanks() MultiTransferParams {
	r := NewMultiTransferParams(mtp.Multi)

	for _, tf := range mtp.Transfers {
		if !tf.IsZero() {
			r.Transfers = append(r.Transfers, tf)
		}
	}

	return r
}

// RemoveVolumes reduce the volume of the ith transfer by vols[i]
func (mtp MultiTransferParams) RemoveVolumes(vols []wunit.Volume) {
	for i, t := range mtp.Transfers {
		if i >= len(vols) {
			break
		} else if !t.IsZero() {
			t.Volume.Subtract(vols[i])
		}
	}
}

// RemoveFVolumes reduce the FVolume (from volumes) of the ith transfer by vols[i]
func (mtp MultiTransferParams) RemoveFVolumes(vols []wunit.Volume) {
	for i, t := range mtp.Transfers {
		if i >= len(vols) {
			break
		} else if !t.IsZero() {
			t.FVolume.Subtract(vols[i])
		}
	}
}

// AddTVolumes increase the TVolume (to volumes) of the ith transfer by vols[i]
func (mtp MultiTransferParams) AddTVolumes(vols []wunit.Volume) {
	for i, t := range mtp.Transfers {
		if i >= len(vols) {
			break
		} else if !t.IsZero() {
			t.TVolume.Add(vols[i])
		}
	}
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
	r := make([]wunit.Volume, 0, mtp.Multi)

	for _, t := range mtp.Transfers {
		r = append(r, t.Volume.Dup())
	}

	for len(r) < mtp.Multi {
		r = append(r, wunit.ZeroVolume())
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
	r := make([]string, mtp.Multi)
	for i, t := range mtp.Transfers {
		r[i] = t.TipType
	}

	return r
}

/*
   FPlateWY:   tp.FPlateWY,
   TPlateWX:   tp.TPlateWX,
   TPlateWY:   tp.TPlateWY,
   Component:  tp.Component,
*/

func (mtp MultiTransferParams) FPlateWX() []int {
	r := make([]int, mtp.Multi)
	for i, t := range mtp.Transfers {
		r[i] = t.FPlateWX
	}
	return r
}
func (mtp MultiTransferParams) TPlateWX() []int {
	r := make([]int, mtp.Multi)
	for i, t := range mtp.Transfers {
		r[i] = t.TPlateWX
	}
	return r
}
func (mtp MultiTransferParams) FPlateWY() []int {
	r := make([]int, mtp.Multi)
	for i, t := range mtp.Transfers {
		r[i] = t.FPlateWY
	}
	return r
}
func (mtp MultiTransferParams) TPlateWY() []int {
	r := make([]int, mtp.Multi)
	for i, t := range mtp.Transfers {
		r[i] = t.TPlateWY
	}
	return r
}
func (mtp MultiTransferParams) Component() []string {
	r := make([]string, mtp.Multi)
	for i, t := range mtp.Transfers {
		r[i] = t.Component
	}
	return r
}

func NewMultiTransferParams(multi int) MultiTransferParams {
	var v MultiTransferParams
	v.Transfers = make([]TransferParams, 0, multi)
	v.Multi = multi
	return v
}

func (mtp MultiTransferParams) ParamSet(n int) TransferParams {
	if n >= len(mtp.Transfers) {
		return TransferParams{}
	}

	return mtp.Transfers[n]
}

func (mtp MultiTransferParams) Channels() []*wtype.LHChannelParameter {
	r := []*wtype.LHChannelParameter{}

	for _, c := range mtp.Transfers {
		r = append(r, c.Channel)
	}
	return r
}

func (mtp MultiTransferParams) ToString() string {
	s := ""
	for i := 0; i < mtp.Multi; i++ {
		s += mtp.ParamSet(i).ToString() + " "
	}

	return s

}

func (mtp MultiTransferParams) Dup() MultiTransferParams {
	tfrs := make([]TransferParams, 0, mtp.Multi)

	for i := 0; i < len(mtp.Transfers); i++ {
		tfrs = append(tfrs, mtp.Transfers[i].Dup())
	}

	ret := NewMultiTransferParams(mtp.Multi)

	ret.Transfers = tfrs

	return ret
}

func MTPFromArrays(what, pltfrom, pltto, wellfrom, wellto, fplatetype, tplatetype []string, volume, fvolume, tvolume []wunit.Volume, FPlateWX, FPlateWY, TPlateWX, TPlateWY []int, Components []string, policies []wtype.LHPolicy) MultiTransferParams {

	mtp := NewMultiTransferParams(len(what))

	for i := 0; i < len(what); i++ {
		mtp.Transfers = append(mtp.Transfers, TransferParams{
			What:       what[i],
			PltFrom:    pltfrom[i],
			PltTo:      pltto[i],
			WellFrom:   wellfrom[i],
			WellTo:     wellto[i],
			FPlateType: fplatetype[i],
			TPlateType: tplatetype[i],
			Volume:     volume[i],
			FVolume:    fvolume[i],
			TVolume:    tvolume[i],
			FPlateWX:   FPlateWX[i],
			TPlateWX:   TPlateWX[i],
			FPlateWY:   FPlateWY[i],
			TPlateWY:   TPlateWY[i],
			Component:  Components[i],
			Policy:     policies[i],
		})
	}

	return mtp
}

type SetOfMultiTransferParams []MultiTransferParams

func (mtp SetOfMultiTransferParams) What() []string {
	sa := make([]string, 0, len(mtp))

	for _, mtp := range mtp {
		sa = append(sa, mtp.What()...)
	}

	return sa
}
func (mtp SetOfMultiTransferParams) PltFrom() []string {
	sa := make([]string, 0, len(mtp))

	for _, mtp := range mtp {
		sa = append(sa, mtp.PltFrom()...)
	}

	return sa
}

func (mtp SetOfMultiTransferParams) PltTo() []string {
	r := make([]string, 0, len(mtp))
	for _, t := range mtp {
		r = append(r, t.PltTo()...)
	}

	return r
}
func (mtp SetOfMultiTransferParams) WellFrom() []string {
	r := make([]string, 0, len(mtp))
	for _, t := range mtp {
		r = append(r, t.WellFrom()...)
	}

	return r
}
func (mtp SetOfMultiTransferParams) WellTo() []string {
	r := make([]string, 0, len(mtp))
	for _, t := range mtp {
		r = append(r, t.WellTo()...)
	}

	return r
}
func (mtp SetOfMultiTransferParams) Volume() []wunit.Volume {
	r := make([]wunit.Volume, 0, len(mtp))

	for _, t := range mtp {
		r = append(r, t.Volume()...)
	}

	return r
}
func (mtp SetOfMultiTransferParams) FPlateType() []string {
	r := make([]string, 0, len(mtp))
	for _, t := range mtp {
		r = append(r, t.FPlateType()...)
	}

	return r
}
func (mtp SetOfMultiTransferParams) TPlateType() []string {
	r := make([]string, 0, len(mtp))
	for _, t := range mtp {
		r = append(r, t.TPlateType()...)
	}

	return r
}

func (mtp SetOfMultiTransferParams) FVolume() []wunit.Volume {
	r := make([]wunit.Volume, 0, len(mtp))

	for _, t := range mtp {
		r = append(r, t.FVolume()...)
	}

	return r
}

func (mtp SetOfMultiTransferParams) TVolume() []wunit.Volume {
	r := make([]wunit.Volume, 0, len(mtp))

	for _, t := range mtp {
		r = append(r, t.TVolume()...)
	}

	return r
}

func (mtp SetOfMultiTransferParams) Channel() []*wtype.LHChannelParameter {
	r := make([]*wtype.LHChannelParameter, 0, len(mtp))
	for _, t := range mtp {
		r = append(r, t.Channel()...)
	}

	return r
}

func (mtp SetOfMultiTransferParams) TipType() []string {
	r := make([]string, 0, len(mtp))
	for _, t := range mtp {
		r = append(r, t.TipType()...)
	}

	return r
}

func (mtp SetOfMultiTransferParams) FPlateWX() []int {
	r := make([]int, 0, len(mtp))
	for _, t := range mtp {
		r = append(r, t.FPlateWX()...)
	}
	return r
}
func (mtp SetOfMultiTransferParams) TPlateWX() []int {
	r := make([]int, 0, len(mtp))
	for _, t := range mtp {
		r = append(r, t.TPlateWX()...)
	}
	return r
}
func (mtp SetOfMultiTransferParams) FPlateWY() []int {
	r := make([]int, 0, len(mtp))
	for _, t := range mtp {
		r = append(r, t.FPlateWY()...)
	}
	return r
}
func (mtp SetOfMultiTransferParams) TPlateWY() []int {
	r := make([]int, 0, len(mtp))
	for _, t := range mtp {
		r = append(r, t.TPlateWY()...)
	}
	return r
}
func (mtp SetOfMultiTransferParams) Component() []string {
	r := make([]string, 0, len(mtp))
	for _, t := range mtp {
		r = append(r, t.Component()...)
	}
	return r
}
