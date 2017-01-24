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

	if vol < min {
		/*
			logger.Fatal(fmt.Sprintf("Error: %f below min vol %f", vol, min))
			panic(errors.New(fmt.Sprintf("Error: %f below min vol %f", vol, min)))
		*/

		err := wtype.LHError(wtype.LH_ERR_VOL, fmt.Sprintf("Liquid Handler cannot service volume requested: %f - minimum volume is %f", vol, min))
		return ret, err
	}

	if vol <= max {
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
		if v.LessThan(ret) {
			ret = v
		}
	}

	return ret
}

func (ins *TransferInstruction) CheckMultiPolicies() bool {
	// first iteration: ensure all the WHAT prms are the same
	// later	  : actually check the policies per channel

	nwhat := wutil.NUniqueStringsInArray(ins.What)

	if nwhat != 1 {
		return false
	}

	return true
}

func (ins *TransferInstruction) GetParallelSetsFor(channel *wtype.LHChannelParameter) [][]int {
	// if the channel is not multi just return nil

	if channel.Multi == 1 {
		return nil
	}

	// the TransferBlock instruction takes into account the destinations being OK
	// splits instructions into potentially multiable blocks on that basis
	// and finds sources for them

	// -- the transfer block ensures these instructions are at most multi long

	// firstly are the sources properly configured?

	npositions := wutil.NUniqueStringsInArray(ins.PltFrom)

	if npositions != 1 {
		// fall back to single-channel
		// TODO -- find a subset we CAN do, if such exists
		return nil
	}

	nplatetypes := wutil.NUniqueStringsInArray(ins.FPlateType)

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

	if !wtype.TipsWellsAligned(*channel, *pa[0], ins.WellFrom) {
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

	return [][]int{{0, 1, 2, 3, 4, 5, 6, 7}}
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

func (ins *TransferInstruction) Generate(policy *wtype.LHPolicyRuleSet, prms *LHProperties) ([]RobotInstruction, error) {
	pol := GetPolicyFor(policy, ins)

	ret := make([]RobotInstruction, 0)

	// if we can multi we do this first

	if pol["CAN_MULTI"].(bool) {
		// basis of its multi, partly based on volume range
		parallelsets := ins.GetParallelSetsFor(prms.HeadsLoaded[0].Params)

		mci := NewMultiChannelBlockInstruction()
		mci.Multi = prms.HeadsLoaded[0].Params.Multi // TODO Remove Hard code here
		mci.Prms = prms.HeadsLoaded[0].Params        // TODO Remove Hard code here
		for _, set := range parallelsets {
			// assemble the info

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
