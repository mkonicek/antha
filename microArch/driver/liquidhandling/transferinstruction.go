// /anthalib/driver/liquidhandling/transferinstruction.go: Part of the Antha language
// Copyright (C) 2015 The Antha authors. All rights reserved.
//
// This program is free software; you can redistribute it and/or
// modify it under the terms of the GNU General Public License
// as published by the Free Software Foundation; either version 2
// of the License, or (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program; if not, write to the Free Software
// Foundation, Inc., 51 Franklin Street, Fifth Floor, Boston, MA  02110-1301, USA.
//
// For more information relating to the software or licensing issues please
// contact license@antha-lang.org or write to the Antha team c/o
// Synthace Ltd. The London Bioscience Innovation Centre
// 2 Royal College St, London NW1 0NH UK

package liquidhandling

import (
	"fmt"
	"math"
	"sort"
	"strconv"
	"strings"

	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/antha/anthalib/wutil"
	"github.com/antha-lang/antha/microArch/factory"
)

func firstInArray(a []*wtype.LHPlate) *wtype.LHPlate {
	for _, v := range a {
		if v != nil {
			return v
		}
	}

	return nil
}

type TransferInstruction struct {
	GenericRobotInstruction
	Type       int
	Platform   string
	What       []string
	PltFrom    []string
	PltTo      []string
	WellFrom   []string
	WellTo     []string
	Volume     []wunit.Volume
	FPlateType []string
	TPlateType []string
	FPlateWX   []int
	FPlateWY   []int
	TPlateWX   []int
	TPlateWY   []int
	FVolume    []wunit.Volume
	TVolume    []wunit.Volume
	Policies   wtype.LHPolicyRuleSet
}

func (ti *TransferInstruction) ToString() string {
	s := fmt.Sprintf("%s ", Robotinstructionnames[ti.Type])
	for i := 0; i < len(ti.What); i++ {
		s += ti.ParamSet(i).ToString()
		s += "\n"
	}

	return s
}

func (ti *TransferInstruction) ParamSet(n int) TransferParams {
	return TransferParams{ti.What[n], ti.PltFrom[n], ti.PltTo[n], ti.WellFrom[n], ti.WellTo[n], ti.Volume[n], ti.FPlateType[n], ti.TPlateType[n], ti.FVolume[n], ti.TVolume[n], nil, ""}
}

func NewTransferInstruction(what, pltfrom, pltto, wellfrom, wellto, fplatetype, tplatetype []string, volume, fvolume, tvolume []wunit.Volume, FPlateWX, FPlateWY, TPlateWX, TPlateWY []int) *TransferInstruction {
	var v TransferInstruction
	v.Type = TFR
	v.What = what
	v.PltFrom = pltfrom
	v.PltTo = pltto
	v.WellFrom = wellfrom
	v.WellTo = wellto
	v.Volume = volume
	v.FPlateType = fplatetype
	v.TPlateType = tplatetype
	v.FVolume = fvolume
	v.TVolume = tvolume
	v.FPlateWX = FPlateWX
	v.FPlateWY = FPlateWY
	v.TPlateWX = TPlateWX
	v.TPlateWY = TPlateWY
	v.GenericRobotInstruction.Ins = RobotInstruction(&v)
	return &v
}

func (ins *TransferInstruction) InstructionType() int {
	return ins.Type
}

func (ins *TransferInstruction) MergeWith(ins2 *TransferInstruction) {
	ins.What = append(ins.What, ins2.What...)
	ins.PltFrom = append(ins.PltFrom, ins2.PltFrom...)
	ins.PltTo = append(ins.PltTo, ins2.PltTo...)
	ins.WellFrom = append(ins.WellFrom, ins2.WellFrom...)
	ins.WellTo = append(ins.WellTo, ins2.WellTo...)
	ins.Volume = append(ins.Volume, ins2.Volume...)
	ins.FPlateType = append(ins.FPlateType, ins2.FPlateType...)
	ins.TPlateType = append(ins.TPlateType, ins2.TPlateType...)
	ins.FPlateWX = append(ins.FPlateWX, ins2.FPlateWX...)
	ins.FPlateWY = append(ins.FPlateWY, ins2.FPlateWY...)
	ins.TPlateWX = append(ins.TPlateWX, ins2.TPlateWX...)
	ins.TPlateWY = append(ins.TPlateWY, ins2.TPlateWY...)
	ins.FVolume = append(ins.FVolume, ins2.FVolume...)
	ins.TVolume = append(ins.TVolume, ins2.TVolume...)
}

func (ins *TransferInstruction) GetParameter(name string) interface{} {
	switch name {
	case "LIQUIDCLASS":
		return ins.What
	case "VOLUME":
		return ins.Volume
	case "VOLUNT":
		return nil
	case "FROMPLATETYPE":
		return ins.FPlateType
	case "WELLFROMVOLUME":
		return ins.FVolume
	case "POSFROM":
		return ins.PltFrom
	case "POSTO":
		return ins.PltTo
	case "WELLFROM":
		return ins.WellFrom
	case "WELLTO":
		return ins.WellTo
	case "WELLTOVOLUME":
		return ins.TVolume
	case "TOPLATETYPE":
		return ins.TPlateType
	case "FPLATEWX":
		return ins.FPlateWX
	case "FPLATEWY":
		return ins.FPlateWY
	case "TPLATEWX":
		return ins.TPlateWX
	case "TPLATEWY":
		return ins.TPlateWY
	case "INSTRUCTIONTYPE":
		return ins.InstructionType()
	case "PLATFORM":
		return ins.Platform
	}
	return nil
}

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
	if !Max.GreaterThanRounded(Vol, 1) {
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

func (vs VolumeSet) MaxMultiTransferVolume() wunit.Volume {
	// the minimum volume in the set

	ret := vs.Vols[0]

	for _, v := range vs.Vols {
		if v.LessThan(ret) && !v.IsZero() {
			ret = v
		}
	}

	return ret
}

func (ins *TransferInstruction) CheckMultiPolicies() bool {
	// first iteration: ensure all the WHAT prms are the same
	// later	  : actually check the policies per channel

	nwhat := wutil.NUniqueStringsInArray(ins.What, true)

	if nwhat != 1 {
		return false
	}

	return true
}

func (ins *TransferInstruction) GetParallelSetsFor(channel *wtype.LHChannelParameter) [][]int {
	// if the channel is not multi just return nil

	if channel.Multi == 1 {
		//fmt.Println("CHANNEL IS NOT MULTI > 1")
		return nil
	}

	// fix for instructions not generated by transfer block

	if len(ins.What) > channel.Multi {
		return nil
	}

	// the TransferBlock instruction takes into account the destinations being OK
	// splits instructions into potentially multiable blocks on that basis
	// and finds sources for them

	// -- the transfer block ensures these instructions are at most multi long
	// all of the below assumes we can't span multiple plates

	// firstly are the sources properly configured?

	npositions := wutil.NUniqueStringsInArray(ins.PltFrom, true)

	if npositions != 1 {
		// fall back to single-channel
		// TODO -- find a subset we CAN do, if such exists
		return nil
	}

	nplatetypes := wutil.NUniqueStringsInArray(ins.FPlateType, true)

	if nplatetypes != 1 {
		// fall back to single-channel
		// TODO -- find a subset we CAN do , if such exists
		return nil
	}

	pa, err := factory.PlateTypeArray(ins.FPlateType)

	if err != nil {
		panic(err)
	}

	// check source / tip alignment

	plate := firstInArray(pa)

	if plate == nil {
		panic("No from plates in instruction")
	}

	if !wtype.TipsWellsAligned(*channel, *plate, ins.WellFrom) {
		// fall back to single-channel
		// TODO -- find a subset we CAN do
		return nil
	}

	pa, err = factory.PlateTypeArray(ins.TPlateType)

	if err != nil {
		panic(err)
	}

	plate = firstInArray(pa)

	if plate == nil {
		panic("No to plates in instruction")
	}

	// for safety, check dest / tip alignment

	if !wtype.TipsWellsAligned(*channel, *plate, ins.WellTo) {
		// fall back to single-channel
		// TODO -- find a subset we CAN do
		return nil
	}

	// check that we will not require different policies

	if !ins.CheckMultiPolicies() {
		// fall back to single-channel
		// TODO -- find a subset we CAN do
		return nil
	}

	// looks OK

	/*
		ra := make([]int, 0, len(ins.What))

		m := 0

		for i := 0; i < len(ins.What); i++ {
			if ins.What[i] != "" {
				m += 1
			}
		}

		for i := 0; i < m; i++ {
			ra = append(ra, i)
		}
	*/

	ra := make([]int, channel.Multi)

	// some issues here in that ins.What might not
	// be the right size:
	// - either too big for some reason, causing segfault
	// - or too small, then the length of ra is too great
	for i := 0; i < len(ins.What); i++ {
		if ins.What[i] != "" {
			ra[i] = i
		} else {
			ra[i] = -1
		}
	}

	return [][]int{ra}
}

func (ins *TransferInstruction) OLDDONTUSETHISGetParallelSetsFor(channel *wtype.LHChannelParameter) [][]int {
	// if the channel is not multi just return nil

	if channel.Multi == 1 {
		return nil
	}

	tfrs := make(map[string][]string)

	// hash out all transfers which are multiable

	for i, _ := range ins.What {
		var tcoord int = -1
		var fcoord int = -1
		var tc2 int = -1
		var fc2 int = -1
		var pmt int = -1
		var pmf int = -1
		wcFrom := wtype.MakeWellCoordsA1(ins.WellFrom[i])
		wcTo := wtype.MakeWellCoordsA1(ins.WellTo[i])

		if channel.Orientation == wtype.LHVChannel {
			// we hash on the X
			tcoord = wcTo.X
			fcoord = wcFrom.X
			tc2 = wcTo.Y
			fc2 = wcFrom.Y
			pmf = ins.FPlateWY[i]
			pmt = ins.TPlateWY[i]
		} else {
			// horizontal orientation
			// hash on the Y
			tcoord = wcTo.Y
			fcoord = wcFrom.Y
			tc2 = wcTo.X
			fc2 = wcFrom.X
			pmf = ins.FPlateWX[i]
			pmt = ins.TPlateWX[i]
		}

		pltF := ins.PltFrom[i]
		pltT := ins.PltTo[i]

		// make hash key

		hashkey := fmt.Sprintf("%s:%s:%d:%s:%d:%d:%d", ins.What[i], pltF, fcoord, pltT, tcoord, pmf, pmt)
		a, ok := tfrs[hashkey]

		if !ok {
			a = make([]string, 0, channel.Multi)
		}

		val := fmt.Sprintf("%d,%d,%d", fc2, tc2, i)
		a = append(a, val)
		tfrs[hashkey] = a
	}

	ret := make([][]int, 0, len(ins.What))

	// now have we got any which are multiable?
	// the elements of each array are transfers with
	// a common source component, row/column and plate on either side
	// now we must check whether the *other* coords match up
	for k, a := range tfrs {
		tx := strings.Split(k, ":")
		pmf, _ := strconv.Atoi(tx[5])
		pmt, _ := strconv.Atoi(tx[6])

		if len(a) >= channel.Multi {
			// could be
			mss := GetMultiSet(a, channel.Multi, pmf, pmt)

			if len(mss) != 0 {
				for _, ms := range mss {
					ret = append(ret, ms)
				}
			}
		}

	}

	if len(ret) == 0 {
		return nil
	}

	return ret
}

func GetMultiSet(a []string, channelmulti int, fromplatemulti int, toplatemulti int) [][]int {
	ret := make([][]int, 0, 2)
	var next []int
	for {
		next, a = GetNextSet(a, channelmulti, fromplatemulti, toplatemulti)
		if next == nil {
			break
		}

		ret = append(ret, next)
	}

	return ret
}

func GetNextSet(a []string, channelmulti int, fromplatemulti int, toplatemulti int) ([]int, []string) {
	if len(a) == 0 {
		return nil, nil
	}
	r := make([][]int, fromplatemulti)
	for i := 0; i < fromplatemulti; i++ {
		r[i] = make([]int, toplatemulti)
		for j := 0; j < toplatemulti; j++ {
			r[i][j] = -1
		}
	}

	// this is simply a greedy algorithm, it may miss things
	for _, s := range a {
		tx := strings.Split(s, ",")
		i, _ := strconv.Atoi(tx[0])
		j, _ := strconv.Atoi(tx[1])
		k, _ := strconv.Atoi(tx[2])
		r[i][j] = k
	}
	// now we just take the first one we find

	ret := getset(r, channelmulti)
	censa := censoredcopy(a, ret)

	return ret, censa
}

func getset(a [][]int, mx int) []int {
	r := make([]int, 0, mx)

	for i := 0; i < len(a); i++ {
		for j := 0; j < len(a[i]); j++ {
			if a[i][j] != -1 {
				r = append(r, a[i][j])
				// find a diagonal line
				for l := 1; l < mx; l++ {
					x := (i + l) % len(a)
					y := (j + l) % len(a[i])

					if a[x][y] != -1 {
						r = append(r, a[x][y])
					} else {
						r = make([]int, 0, mx)
					}
				}

				if len(r) == mx {
					break
				}
			}
		}
	}

	if len(r) == mx {
		sort.Ints(r)
		return r
	} else {
		return nil
	}
}

func censoredcopy(a []string, b []int) []string {
	if b == nil {
		return a
	}

	r := make([]string, 0, len(a)-len(b))

	for _, x := range a {
		tx := strings.Split(x, ",")
		i, _ := strconv.Atoi(tx[2])
		if IsIn(i, b) {
			continue
		}
		r = append(r, x)
	}

	return r
}

func IsIn(i int, a []int) bool {
	for _, x := range a {
		if i == x {
			return true
		}
	}

	return false
}

// helper thing

type VolumeSet struct {
	Vols []wunit.Volume
}

func NewVolumeSet(n int) VolumeSet {
	var vs VolumeSet
	vs.Vols = make([]wunit.Volume, n)
	for i := 0; i < n; i++ {
		vs.Vols[i] = (wunit.NewVolume(0.0, "ul"))
	}
	return vs
}

func (vs VolumeSet) Add(v wunit.Volume) {
	for i := 0; i < len(vs.Vols); i++ {
		vs.Vols[i].Add(v)
	}
}

func (vs VolumeSet) Sub(v wunit.Volume) []wunit.Volume {
	ret := make([]wunit.Volume, len(vs.Vols))
	for i := 0; i < len(vs.Vols); i++ {
		vs.Vols[i].Subtract(v)
		ret[i] = wunit.CopyVolume(v)
	}
	return ret
}

func (vs VolumeSet) SetEqualTo(v wunit.Volume) {
	for i := 0; i < len(vs.Vols); i++ {
		vs.Vols[i] = wunit.CopyVolume(v)
	}
}

func (vs VolumeSet) GetACopy() []wunit.Volume {
	r := make([]wunit.Volume, len(vs.Vols))
	for i := 0; i < len(vs.Vols); i++ {
		r[i] = wunit.CopyVolume(vs.Vols[i])
	}
	return r
}

func countSetSize(set []int) int {
	c := 0
	for _, v := range set {
		if v != -1 {
			c += 1
		}
	}

	return c
}

func (ins *TransferInstruction) ChooseChannels(prms *LHProperties) {
	// trims out leading blanks for now... will eventually call
	// CanAddress on prms
	wh := make([]string, len(ins.What))
	pf := make([]string, len(ins.What))
	pt := make([]string, len(ins.What))
	wf := make([]string, len(ins.What))
	wt := make([]string, len(ins.What))
	vl := make([]wunit.Volume, len(ins.What))
	fpt := make([]string, len(ins.What))
	tpt := make([]string, len(ins.What))
	fpwx := make([]int, len(ins.What))
	fpwy := make([]int, len(ins.What))
	tpwx := make([]int, len(ins.What))
	tpwy := make([]int, len(ins.What))
	fv := make([]wunit.Volume, len(ins.What))
	tv := make([]wunit.Volume, len(ins.What))

	ix := -1
	for i := 0; i < len(ins.What); i++ {
		if ix > -1 {
			ix += 1
		}

		if ins.What[i] != "" {
			if ix == -1 {
				ix = 0
			}
			wh[ix] = ins.What[i]
			pf[ix] = ins.PltFrom[i]
			pt[ix] = ins.PltTo[i]
			wf[ix] = ins.WellFrom[i]
			wt[ix] = ins.WellTo[i]
			vl[ix] = ins.Volume[i]
			fpt[ix] = ins.FPlateType[i]
			tpt[ix] = ins.TPlateType[i]
			fpwx[ix] = ins.FPlateWX[i]
			fpwy[ix] = ins.FPlateWY[i]
			tpwx[ix] = ins.TPlateWX[i]
			tpwy[ix] = ins.TPlateWY[i]
			fv[ix] = ins.FVolume[i]
			tv[ix] = ins.TVolume[i]
		}
	}

	ins.What = wh
	ins.PltFrom = pf
	ins.PltTo = pt
	ins.WellFrom = wf
	ins.WellTo = wt
	ins.Volume = vl
	ins.FPlateType = fpt
	ins.TPlateType = tpt
	ins.FPlateWX = fpwx
	ins.FPlateWY = fpwy
	ins.TPlateWX = tpwx
	ins.TPlateWY = tpwy
	ins.FVolume = fv
	ins.TVolume = tv
}

//	This section divides transfers into multi and single channel blocks (MCB,SCB respectively)
//      with the constraint that it must respect user requests for atomicity
//
//	The constraint comes from the fact that an instruction can be a request to mix several
//	components atomically - i.e. once the first component has been moved the rest must follow
//	immediately. This is compatible with multichannel operation in some cases but not others
//
//	if multichannel is not allowed for any of the components then it will at this point revert
//	to single-channel operation as follows:
//	(in the below assume all volumes are identical)
//
//	input: 	LHIVector([A], [B], [C])
//	output:
//		without multichannel:
//			SCB([A,B,C], [d1, d2, d3])
//		with multichannel   :
//			MCB([A,B,C], [d1, d2, d3])
//
//	For reference, vectors including atomic mixes look like this
//
//	LHIVector([A,B,C], [A], [A])
//
//	output:
//		without multichannel:
//			SCB([A,B,C,A,A],[d1,d1,d1,d2,d3])
//
//		with multichannel:
//			MCB([A,A,A],[d1,d2,d3]), SCB([B,C],[d1,d1])
//
//	The following is always mapped to single-channel operation
//
//		LHIVector([A,B], [A,C])
//
//	output:
//		either case:
//			SCB([A,B,A,C],[d1,d1,d2,d2])
//

func (ins *TransferInstruction) Generate(policy *wtype.LHPolicyRuleSet, prms *LHProperties) ([]RobotInstruction, error) {
	//  set the channel  choices first by cleaning out initial empties

	ins.ChooseChannels(prms)

	pol := GetPolicyFor(policy, ins)

	ret := make([]RobotInstruction, 0)

	// if we can multi we do this first

	if pol["CAN_MULTI"].(bool) {
		parallelsets := ins.GetParallelSetsFor(prms.HeadsLoaded[0].Params)

		mci := NewMultiChannelBlockInstruction()
		//mci.Multi = prms.HeadsLoaded[0].Params.Multi // TODO Remove Hard code here
		mci.Prms = prms.HeadsLoaded[0].Params // TODO Remove Hard code here
		for _, set := range parallelsets {
			// assemble the info
			mci.Multi = countSetSize(set)
			vols := NewVolumeSet(len(set))
			fvols := NewVolumeSet(len(set))
			tvols := NewVolumeSet(len(set))
			What := make([]string, len(set))
			PltFrom := make([]string, len(set))
			PltTo := make([]string, len(set))
			WellFrom := make([]string, len(set))
			WellTo := make([]string, len(set))
			FPlateType := make([]string, len(set))
			TPlateType := make([]string, len(set))

			for i, s := range set {
				if s == -1 {
					continue
				}
				vols.Vols[i] = wunit.CopyVolume(ins.Volume[s])
				fvols.Vols[i] = wunit.CopyVolume(ins.FVolume[s])
				tvols.Vols[i] = wunit.CopyVolume(ins.TVolume[s])
				What[i] = ins.What[s]
				PltFrom[i] = ins.PltFrom[s]
				PltTo[i] = ins.PltTo[s]
				WellFrom[i] = ins.WellFrom[s]
				WellTo[i] = ins.WellTo[s]
				FPlateType[i] = ins.FPlateType[s]
				TPlateType[i] = ins.TPlateType[s]
			}
			// get the max transfer volume

			maxvol := vols.MaxMultiTransferVolume()

			// now set the vols for the transfer and remove this from the instruction's volume

			for i, _ := range vols.Vols {
				if set[i] == -1 {
					continue
				}
				vols.Vols[i] = wunit.CopyVolume(maxvol)
				ins.Volume[set[i]].Subtract(maxvol)

				// set the from and to volumes for the relevant part of the instruction
				// NB -- this is a design issue which should probably be fixed: at the moment
				// if we have two instructions which refer to the same underlying well their
				// volume levels will not be in sync
				// therefore this implementation is not correct as regards changes of underlying
				// state
				//... instead the right thing would be for all of these instructions to reference
				// plate objects instead - this will work OK as long as we have a shared memory
				// system... otherwise we'll need to use channels
				ins.FVolume[set[i]].Subtract(maxvol)
				ins.TVolume[set[i]].Add(maxvol)
			}

			tp := NewMultiTransferParams(mci.Multi)
			tp.What = What
			tp.Volume = vols.Vols
			tp.FVolume = fvols.Vols
			tp.TVolume = tvols.Vols
			tp.PltFrom = PltFrom
			tp.PltTo = PltTo
			tp.WellFrom = WellFrom
			tp.WellTo = WellTo
			tp.FPlateType = FPlateType
			tp.TPlateType = TPlateType
			tp.Channel = mci.Prms
			mci.AddTransferParams(tp)
		}

		if len(parallelsets) > 0 {
			ret = append(ret, mci)
		}
	}

	// mop up all the single instructions which are left
	sci := NewSingleChannelBlockInstruction()
	sci.Prms = prms.HeadsLoaded[0].Params // TODO Fix Hard Code Here

	for i, _ := range ins.What {
		if ins.Volume[i].LessThanFloat(0.001) {
			continue
		}
		if i != 0 && (ins.What[i] != ins.What[i-1]) {
			if len(sci.Volume) > 0 {
				ret = append(ret, sci)
			}
			sci = NewSingleChannelBlockInstruction()
			sci.Prms = prms.HeadsLoaded[0].Params
		}

		var tp TransferParams

		tp.What = ins.What[i]
		tp.PltFrom = ins.PltFrom[i]
		tp.PltTo = ins.PltTo[i]
		tp.WellFrom = ins.WellFrom[i]
		tp.WellTo = ins.WellTo[i]
		tp.Volume = wunit.CopyVolume(ins.Volume[i])
		tp.FVolume = wunit.CopyVolume(ins.FVolume[i])
		tp.TVolume = wunit.CopyVolume(ins.TVolume[i])
		tp.FPlateType = ins.FPlateType[i]
		tp.TPlateType = ins.TPlateType[i]
		sci.AddTransferParams(tp)

		// make sure we keep volumes up to date

		ins.FVolume[i].Subtract(ins.Volume[i])
		ins.TVolume[i].Add(ins.Volume[i])
	}
	if len(sci.Volume) > 0 {
		ret = append(ret, sci)
	}

	return ret, nil
}
