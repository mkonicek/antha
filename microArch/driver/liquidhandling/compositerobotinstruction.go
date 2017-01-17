// /anthalib/driver/liquidhandling/compositerobotinstruction.go: Part of the Antha language
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
	anthadriver "github.com/antha-lang/antha/microArch/driver"
	"github.com/antha-lang/antha/microArch/logger"
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

type MultiTransferParams struct {
	What       []string
	PltFrom    []string
	PltTo      []string
	WellFrom   []string
	WellTo     []string
	Volume     []wunit.Volume
	FPlateType []string
	TPlateType []string
	FVolume    []wunit.Volume
	TVolume    []wunit.Volume
	Channel    *wtype.LHChannelParameter
	TipTypes   []string
}

func NewMultiTransferParams(multi int) MultiTransferParams {
	var v MultiTransferParams
	v.What = make([]string, 0, multi)
	v.PltFrom = make([]string, 0, multi)
	v.PltTo = make([]string, 0, multi)
	v.WellFrom = make([]string, 0, multi)
	v.WellTo = make([]string, 0, multi)
	v.Volume = make([]wunit.Volume, 0, multi)
	v.FVolume = make([]wunit.Volume, 0, multi)
	v.TVolume = make([]wunit.Volume, 0, multi)
	v.FPlateType = make([]string, 0, multi)
	v.TPlateType = make([]string, 0, multi)
	v.TipTypes = make([]string, 0, multi)
	return v
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

func (ins *TransferInstruction) GetParallelSetsFor(channel *wtype.LHChannelParameter) [][]int {
	// if the channel is not multi just return nil

	if channel.Multi == 1 {
		return nil
	}

	tfrs := make(map[string][]string, len(ins.What))

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
				/*
					for l := 1; l < mx; l++ {
						x := (i + l) % len(a)
						y := (j + l) % len(a[i])

						if a[x][y] != -1 {
							r = append(r, a[x][y])
						} else {
							r = make([]int, 0, mx)
						}
					}
				*/

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
		// break out the sets of parallel instructions

		// fix this HARD CODE here
		// this is a bit problematic: we need to define head choice here partly on the
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

type SingleChannelBlockInstruction struct {
	GenericRobotInstruction
	Type       int
	What       []string
	PltFrom    []string
	PltTo      []string
	WellFrom   []string
	WellTo     []string
	Volume     []wunit.Volume
	FPlateType []string
	TPlateType []string
	FVolume    []wunit.Volume
	TVolume    []wunit.Volume
	Prms       *wtype.LHChannelParameter
}

func NewSingleChannelBlockInstruction() *SingleChannelBlockInstruction {
	var v SingleChannelBlockInstruction
	v.Type = SCB
	v.What = make([]string, 0)
	v.PltFrom = make([]string, 0)
	v.PltTo = make([]string, 0)
	v.WellFrom = make([]string, 0)
	v.WellTo = make([]string, 0)
	v.Volume = make([]wunit.Volume, 0)
	v.FVolume = make([]wunit.Volume, 0)
	v.TVolume = make([]wunit.Volume, 0)
	v.FPlateType = make([]string, 0)
	v.TPlateType = make([]string, 0)
	v.GenericRobotInstruction.Ins = &v
	return &v
}

func (ins *SingleChannelBlockInstruction) AddTransferParams(mct TransferParams) {
	ins.What = append(ins.What, mct.What)
	ins.PltFrom = append(ins.PltFrom, mct.PltFrom)
	ins.PltTo = append(ins.PltTo, mct.PltTo)
	ins.WellFrom = append(ins.WellFrom, mct.WellFrom)
	ins.WellTo = append(ins.WellTo, mct.WellTo)
	ins.Volume = append(ins.Volume, mct.Volume)
	ins.FPlateType = append(ins.FPlateType, mct.FPlateType)
	ins.TPlateType = append(ins.TPlateType, mct.TPlateType)
	ins.FVolume = append(ins.FVolume, mct.FVolume)
	ins.TVolume = append(ins.TVolume, mct.TVolume)
	ins.Prms = mct.Channel
}
func (ins *SingleChannelBlockInstruction) InstructionType() int {
	return ins.Type
}

func (ins *SingleChannelBlockInstruction) GetParameter(name string) interface{} {
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
	case "PARAMS":
		return ins.Prms
	case "PLATFORM":
		if ins.Prms == nil {
			return ""
		}
		return ins.Prms.Platform
	case "WELLTO":
		return ins.WellTo
	case "WELLTOVOLUME":
		return ins.TVolume
	case "TOPLATETYPE":
		return ins.TPlateType
	case "INSTRUCTIONTYPE":
		return ins.InstructionType()
	}
	return nil
}

func (ins *SingleChannelBlockInstruction) Generate(policy *wtype.LHPolicyRuleSet, prms *LHProperties) ([]RobotInstruction, error) {
	ret := make([]RobotInstruction, 0)
	// get tips
	channel, tiptype := ChooseChannel(ins.Volume[0], prms)
	ins.Prms = channel
	pol := GetPolicyFor(policy, ins)
	tipget, err := GetTips(tiptype, prms, channel, 1, false)

	if err != nil {
		return ret, err
	}

	ret = append(ret, tipget)
	n_tip_uses := 0

	var last_thing *wtype.LHComponent

	last_thing = nil

	var dirty bool

	for t := 0; t < len(ins.Volume); t++ {
		newchannel, newtiptype := ChooseChannel(ins.Volume[t], prms)
		tvs, err := TransferVolumes(ins.Volume[t], newchannel.Minvol, newchannel.Maxvol)

		if err != nil {
			return ret, err
		}
		for _, vol := range tvs {
			// determine whether to change tips
			change_tips := false
			change_tips = n_tip_uses > pol["TIP_REUSE_LIMIT"].(int)
			change_tips = change_tips || channel != newchannel
			change_tips = change_tips || newtiptype != tiptype

			this_thing := prms.Plates[ins.PltFrom[t]].Wellcoords[ins.WellFrom[t]].Contents()

			if last_thing != nil {
				if this_thing.CName != last_thing.CName {
					change_tips = true
				}
			}

			// finally ensure we don't contaminate sources
			if dirty {
				change_tips = true
			}

			if change_tips {
				/*
					fmt.Println("CHANGING TIPS HERE ")
					fmt.Println("THIS: ", this_thing.CName, " THAT: ", last_thing.CName)
					fmt.Println("channels equal? ", channel == newchannel)
					fmt.Println("tips same? ", tiptype == newtiptype)
					fmt.Println("tip reuse over? ", n_tip_uses > pol["TIP_REUSE_LIMIT"].(int))
					fmt.Println(n_tip_uses, " ", pol["TIP_REUSE_LIMIT"].(int), " NOT CAST: ", pol["TIP_REUSE_LIMIT"])
					fmt.Println("dirty? ", dirty)
				*/
				// maybe wrap this as a ChangeTips function call
				// these need parameters

				tipdrp, err := DropTips(tiptype, prms, channel, 1)
				if err != nil {
					return ret, err
				}
				ret = append(ret, tipdrp)

				tipget, err := GetTips(newtiptype, prms, newchannel, 1, false)

				if err != nil {
					return ret, err
				}

				ret = append(ret, tipget)
				tiptype = newtiptype
				channel = newchannel
				n_tip_uses = 0
				last_thing = nil
				dirty = false
			}

			stci := NewSingleChannelTransferInstruction()

			stci.What = ins.What[t]
			stci.PltFrom = ins.PltFrom[t]
			stci.PltTo = ins.PltTo[t]
			stci.WellFrom = ins.WellFrom[t]
			stci.WellTo = ins.WellTo[t]
			stci.Volume = vol
			stci.FPlateType = ins.FPlateType[t]
			stci.TPlateType = ins.TPlateType[t]
			stci.FVolume = wunit.CopyVolume(ins.FVolume[t])
			stci.TVolume = wunit.CopyVolume(ins.TVolume[t])
			stci.Prms = channel
			stci.TipType = tiptype
			ret = append(ret, stci)
			last_thing = this_thing

			// finally check if we are touching a bad liquid
			// in future we will do this properly, for now we assume
			// touching any liquid is bad

			npre, premix := pol["PRE_MIX"]
			npost, postmix := pol["POST_MIX"]

			if pol["DSPREFERENCE"].(int) == 0 && !ins.TVolume[t].IsZero() || premix && npre.(int) > 0 || postmix && npost.(int) > 0 {
				dirty = true
				/*
					fmt.Println("DIRTY DIRTY DIRTY")
					fmt.Println(pol["DSPREFERENCE"].(int) == 0, " ", ins.TVolume[t].ToString(), " ", "PRE: ", premix, " ", npre, " POST: ", postmix, " ", npost)
					fmt.Println(ins.WellTo[t])
				*/
			}

			ins.FVolume[t].Subtract(vol)
			ins.TVolume[t].Add(vol)
			n_tip_uses += 1
		}

	}
	tipdrp, err := DropTips(tiptype, prms, channel, 1)

	if err != nil {
		return ret, err
	}
	ret = append(ret, tipdrp)

	return ret, nil
}

type MultiChannelBlockInstruction struct {
	GenericRobotInstruction
	Type       int
	What       [][]string
	PltFrom    [][]string
	PltTo      [][]string
	WellFrom   [][]string
	WellTo     [][]string
	Volume     [][]wunit.Volume
	FPlateType [][]string
	TPlateType [][]string
	FVolume    [][]wunit.Volume
	TVolume    [][]wunit.Volume
	Multi      int
	Prms       *wtype.LHChannelParameter
}

func NewMultiChannelBlockInstruction() *MultiChannelBlockInstruction {
	var v MultiChannelBlockInstruction
	v.What = make([][]string, 0)
	v.PltFrom = make([][]string, 0)
	v.PltTo = make([][]string, 0)
	v.WellFrom = make([][]string, 0)
	v.WellTo = make([][]string, 0)
	v.Volume = make([][]wunit.Volume, 0)
	v.FPlateType = make([][]string, 0)
	v.TPlateType = make([][]string, 0)
	v.FVolume = make([][]wunit.Volume, 0)
	v.TVolume = make([][]wunit.Volume, 0)
	v.GenericRobotInstruction.Ins = &v
	return &v
}

func (ins *MultiChannelBlockInstruction) AddTransferParams(mct MultiTransferParams) {
	ins.What = append(ins.What, mct.What)
	ins.PltFrom = append(ins.PltFrom, mct.PltFrom)
	ins.PltTo = append(ins.PltTo, mct.PltTo)
	ins.WellFrom = append(ins.WellFrom, mct.WellFrom)
	ins.WellTo = append(ins.WellTo, mct.WellTo)
	ins.Volume = append(ins.Volume, mct.Volume)
	ins.FPlateType = append(ins.FPlateType, mct.FPlateType)
	ins.TPlateType = append(ins.TPlateType, mct.TPlateType)
	ins.FVolume = append(ins.FVolume, mct.FVolume)
	ins.TVolume = append(ins.TVolume, mct.TVolume)
	ins.Prms = mct.Channel
}

func (ins *MultiChannelBlockInstruction) InstructionType() int {
	return ins.Type
}

func (ins *MultiChannelBlockInstruction) GetParameter(name string) interface{} {
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
	case "PARAMS":
		return ins.Prms
	case "PLATFORM":
		if ins.Prms == nil {
			return ""
		}
		return ins.Prms.Platform
	case "WELLTO":
		return ins.WellTo
	case "WELLTOVOLUME":
		return ins.TVolume
	case "TOPLATETYPE":
		return ins.TPlateType
	case "INSTRUCTIONTYPE":
		return ins.InstructionType()
	}
	return nil
}

func (ins *MultiChannelBlockInstruction) Generate(policy *wtype.LHPolicyRuleSet, prms *LHProperties) ([]RobotInstruction, error) {
	pol := GetPolicyFor(policy, ins)
	ret := make([]RobotInstruction, 0)
	// get some tips
	//eesh... potential conflict here
	channel, tiptype := ChooseChannel(ins.Volume[0][0], prms)
	tipget, err := GetTips(tiptype, prms, channel, ins.Multi, false)
	if err != nil {
		return ret, err
	}
	ret = append(ret, tipget)
	n_tip_uses := 0

	for t := 0; t < len(ins.Volume); t++ {
		tvols := NewVolumeSet(ins.Prms.Multi)
		vols := NewVolumeSet(ins.Prms.Multi)
		fvols := NewVolumeSet(ins.Prms.Multi)
		for i, _ := range ins.Volume[t] {
			fvols.Vols[i] = wunit.CopyVolume(ins.FVolume[t][i])
			tvols.Vols[i] = wunit.CopyVolume(ins.TVolume[t][i])
		}

		// choose tips
		newchannel, newtiptype := ChooseChannel(ins.Volume[0][0], prms)
		var last_thing *wtype.LHComponent
		last_thing = nil
		var dirty bool
		// load tips

		// split the transfer up
		// NB we assume all volumes are equal here;
		// oof we need to do some work here -- this needs to be in sync with singlechannel block
		tvs, err := TransferVolumes(ins.Volume[t][0], newchannel.Minvol, newchannel.Maxvol)

		if err != nil {
			return ret, err
		}

		for _, vol := range tvs {
			// determine whether to change tips
			change_tips := false
			change_tips = n_tip_uses > pol["TIP_REUSE_LIMIT"].(int)
			change_tips = change_tips || channel != newchannel
			change_tips = change_tips || newtiptype != tiptype

			// big dangerous assumption here: we need to check if anything is different
			this_thing := prms.Plates[ins.PltFrom[t][0]].Wellcoords[ins.WellFrom[t][0]].Contents()

			if last_thing != nil {
				if this_thing.CName != last_thing.CName {
					change_tips = true
				}
			}

			// finally ensure we don't contaminate sources
			if dirty {
				change_tips = true
			}

			if change_tips {
				// maybe wrap this as a ChangeTips function call
				// these need parameters
				tipdrp, err := DropTips(tiptype, prms, channel, ins.Multi)

				if err != nil {
					return ret, err
				}
				ret = append(ret, tipdrp)

				tipget, err := GetTips(newtiptype, prms, newchannel, ins.Multi, false)
				if err != nil {
					return ret, err
				}

				ret = append(ret, tipget)
				tiptype = newtiptype
				channel = newchannel
				n_tip_uses = 0
				last_thing = nil
				dirty = false
			}
			/*
				// enforce tip usage policy

				if n_tip_uses > pol["TIP_REUSE_LIMIT"].(int) || newchannel != channel || newtiptype != tiptype {
					// these need parameters
					ret = append(ret, DropTips(tiptype, prms, channel, ins.Multi))
					ret = append(ret, GetTips(newtiptype, prms, newchannel, ins.Multi, false))
					n_tip_uses = 0
				}
			*/
			mci := NewMultiChannelTransferInstruction()
			vols.SetEqualTo(vol)
			mci.What = ins.What[t]
			mci.Volume = vols.GetACopy()
			mci.FVolume = fvols.GetACopy()
			mci.TVolume = tvols.GetACopy()
			mci.PltFrom = ins.PltFrom[t]
			mci.PltTo = ins.PltTo[t]
			mci.WellFrom = ins.WellFrom[t]
			mci.WellTo = ins.WellTo[t]
			mci.FPlateType = ins.FPlateType[t]
			mci.TPlateType = ins.TPlateType[t]
			mci.Multi = ins.Multi
			//mci.Prms = ins.Prms
			mci.Prms = channel
			ret = append(ret, mci)

			tiptype = newtiptype
			channel = newchannel
			fvols.Sub(vol)
			tvols.Add(vol)
		}
	}

	// remove tips
	tipdrp, err := DropTips(tiptype, prms, channel, ins.Multi)

	if err != nil {
		return ret, err
	}

	ret = append(ret, tipdrp)

	return ret, nil
}

type SingleChannelTransferInstruction struct {
	GenericRobotInstruction
	Type       int
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
	Prms       *wtype.LHChannelParameter
	TipType    string
}

func (scti *SingleChannelTransferInstruction) Params() TransferParams {
	var tp TransferParams
	tp.What = scti.What
	tp.PltFrom = scti.PltFrom
	tp.PltTo = scti.PltTo
	tp.WellTo = scti.WellTo
	tp.WellFrom = scti.WellFrom
	tp.Volume = wunit.CopyVolume(scti.Volume)
	tp.FPlateType = scti.FPlateType
	tp.TPlateType = scti.TPlateType
	tp.FVolume = wunit.CopyVolume(scti.FVolume)
	tp.TVolume = wunit.CopyVolume(scti.TVolume)
	tp.Channel = scti.Prms
	tp.TipType = scti.TipType
	return tp
}

func NewSingleChannelTransferInstruction() *SingleChannelTransferInstruction {
	var v SingleChannelTransferInstruction
	v.Type = SCT
	return &v
}
func (ins *SingleChannelTransferInstruction) InstructionType() int {
	return ins.Type
}

func (ins *SingleChannelTransferInstruction) GetParameter(name string) interface{} {
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
	case "PARAMS":
		return ins.Prms
	case "PLATFORM":
		if ins.Prms == nil {
			return ""
		}
		return ins.Prms.Platform
	case "WELLTO":
		return ins.WellTo
	case "WELLTOVOLUME":
		return ins.TVolume
	case "TOPLATETYPE":
		return ins.TPlateType
	case "INSTRUCTIONTYPE":
		return ins.InstructionType()
	case "TIPTYPE":
		return ins.TipType
	}
	return nil
}

func (ins *SingleChannelTransferInstruction) Generate(policy *wtype.LHPolicyRuleSet, prms *LHProperties) ([]RobotInstruction, error) {
	ret := make([]RobotInstruction, 0)
	// make the instructions

	suckinstruction := NewSuckInstruction()
	suckinstruction.AddTransferParams(ins.Params())
	suckinstruction.Multi = 1
	suckinstruction.Prms = ins.Prms
	ret = append(ret, suckinstruction)

	blowinstruction := NewBlowInstruction()
	blowinstruction.AddTransferParams(ins.Params())
	blowinstruction.Multi = 1
	blowinstruction.Prms = ins.Prms
	ret = append(ret, blowinstruction)

	/*
		// commented out pending putting it as part of blow
		// need to append to reset command
		resetinstruction := NewResetInstruction()
		resetinstruction.AddTransferParams(ins.Params())
		resetinstruction.Prms = ins.Prms
		ret = append(ret, resetinstruction)
	*/

	return ret, nil
}

type MultiChannelTransferInstruction struct {
	GenericRobotInstruction
	Type       int
	What       []string
	PltFrom    []string
	PltTo      []string
	WellFrom   []string
	WellTo     []string
	Volume     []wunit.Volume
	FPlateType []string
	TPlateType []string
	FVolume    []wunit.Volume
	TVolume    []wunit.Volume
	Multi      int
	Prms       *wtype.LHChannelParameter
}

func (scti *MultiChannelTransferInstruction) Params(k int) TransferParams {
	var tp TransferParams
	tp.What = scti.What[k]
	tp.PltFrom = scti.PltFrom[k]
	tp.PltTo = scti.PltTo[k]
	tp.WellFrom = scti.WellFrom[k]
	tp.WellTo = scti.WellTo[k]
	tp.Volume = wunit.CopyVolume(scti.Volume[k])
	tp.FPlateType = scti.FPlateType[k]
	tp.TPlateType = scti.TPlateType[k]
	tp.FVolume = wunit.CopyVolume(scti.FVolume[k])
	tp.TVolume = wunit.CopyVolume(scti.TVolume[k])
	tp.Channel = scti.Prms
	return tp
}
func NewMultiChannelTransferInstruction() *MultiChannelTransferInstruction {
	var v MultiChannelTransferInstruction
	v.Type = MCT
	v.What = make([]string, 0)
	v.PltFrom = make([]string, 0)
	v.PltTo = make([]string, 0)
	v.WellFrom = make([]string, 0)
	v.WellTo = make([]string, 0)
	v.Volume = make([]wunit.Volume, 0)
	v.FVolume = make([]wunit.Volume, 0)
	v.TVolume = make([]wunit.Volume, 0)
	v.FPlateType = make([]string, 0)
	v.TPlateType = make([]string, 0)
	return &v
}
func (ins *MultiChannelTransferInstruction) InstructionType() int {
	return ins.Type
}

func (ins *MultiChannelTransferInstruction) GetParameter(name string) interface{} {
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
	case "PARAMS":
		return ins.Prms
	case "PLATFORM":
		if ins.Prms == nil {
			return ""
		}
		return ins.Prms.Platform
	case "WELLTO":
		return ins.WellTo
	case "WELLTOVOLUME":
		return ins.TVolume
	case "TOPLATETYPE":
		return ins.TPlateType
	case "INSTRUCTIONTYPE":
		return ins.InstructionType()
	}
	return nil
}

func (ins *MultiChannelTransferInstruction) Generate(policy *wtype.LHPolicyRuleSet, prms *LHProperties) ([]RobotInstruction, error) {
	ret := make([]RobotInstruction, 0)

	// make the instructions

	suckinstruction := NewSuckInstruction()
	blowinstruction := NewBlowInstruction()
	suckinstruction.Multi = ins.Multi
	blowinstruction.Multi = ins.Multi
	suckinstruction.Prms = ins.Prms
	blowinstruction.Prms = ins.Prms

	for i := 0; i < len(ins.Volume); i++ {
		suckinstruction.AddTransferParams(ins.Params(i))
		blowinstruction.AddTransferParams(ins.Params(i))
	}

	ret = append(ret, suckinstruction)
	ret = append(ret, blowinstruction)

	return ret, nil
}

type StateChangeInstruction struct {
	GenericRobotInstruction
	Type     int
	OldState *wtype.LHChannelParameter
	NewState *wtype.LHChannelParameter
}

func NewStateChangeInstruction(oldstate, newstate *wtype.LHChannelParameter) *StateChangeInstruction {
	var v StateChangeInstruction
	v.Type = CCC
	v.OldState = oldstate
	v.NewState = newstate
	return &v
}
func (ins *StateChangeInstruction) InstructionType() int {
	return ins.Type
}

func (ins *StateChangeInstruction) GetParameter(name string) interface{} {
	switch name {
	case "OLDSTATE":
		return ins.OldState
	case "NEWSTATE":
		return ins.NewState
	case "INSTRUCTIONTYPE":
		return ins.InstructionType()
	}
	return nil
}

func (ins *StateChangeInstruction) Generate(policy *wtype.LHPolicyRuleSet, prms *LHProperties) ([]RobotInstruction, error) {
	return nil, nil
}

type ChangeAdaptorInstruction struct {
	GenericRobotInstruction
	Type           int
	Head           int
	DropPosition   string
	GetPosition    string
	OldAdaptorType string
	NewAdaptorType string
	Platform       string
}

func NewChangeAdaptorInstruction(head int, droppos, getpos, oldad, newad, platform string) *ChangeAdaptorInstruction {
	var v ChangeAdaptorInstruction
	v.Type = CHA
	v.Head = head
	v.DropPosition = droppos
	v.GetPosition = getpos
	v.OldAdaptorType = oldad
	v.NewAdaptorType = newad
	v.Platform = platform
	return &v
}
func (ins *ChangeAdaptorInstruction) InstructionType() int {
	return ins.Type
}

func (ins *ChangeAdaptorInstruction) GetParameter(name string) interface{} {
	switch name {
	case "HEAD":
		return ins.Head
	case "POSFROM":
		return ins.DropPosition
	case "POSTO":
		return ins.GetPosition
	case "OLDADAPTOR":
		return ins.OldAdaptorType
	case "NEWADAPTOR":
		return ins.NewAdaptorType
	case "INSTRUCTIONTYPE":
		return ins.InstructionType()
	case "PLATFORM":
		return ins.Platform
	}
	return nil
}

func (ins *ChangeAdaptorInstruction) Generate(policy *wtype.LHPolicyRuleSet, prms *LHProperties) ([]RobotInstruction, error) {
	ret := make([]RobotInstruction, 4)
	/*
		ret[0]=NewMoveInstruction(ins.DropPosition,...)
		ret[1]=NewUnloadAdaptorInstruction(ins.DropPosition,...)
		ret[2]=NewMoveInstruction(ins.GetPosition, ...)
		ret[3]=NewLoadAdaptorInstruction(ins.GetPosition,...)
	*/

	return ret, nil
}

type LoadTipsMoveInstruction struct {
	GenericRobotInstruction
	Type       int
	Head       int
	Well       []string
	FPosition  []string
	FPlateType []string
	Multi      int
	Platform   string
}

func NewLoadTipsMoveInstruction() *LoadTipsMoveInstruction {
	var v LoadTipsMoveInstruction
	v.Type = LDT
	v.Well = make([]string, 0)
	v.FPosition = make([]string, 0)
	v.FPlateType = make([]string, 0)
	return &v
}
func (ins *LoadTipsMoveInstruction) InstructionType() int {
	return ins.Type
}

func (ins *LoadTipsMoveInstruction) GetParameter(name string) interface{} {
	switch name {
	case "HEAD":
		return ins.Head
	case "FROMPLATETYPE":
		return ins.FPlateType
	case "POSFROM":
		return ins.FPosition
	case "WELLFROM":
		return ins.Well
	case "Multi":
		return ins.Multi
	case "INSTRUCTIONTYPE":
		return ins.InstructionType()
	case "PLATFORM":
		return ins.Platform
	}
	return nil
}

func (ins *LoadTipsMoveInstruction) Generate(policy *wtype.LHPolicyRuleSet, prms *LHProperties) ([]RobotInstruction, error) {
	ret := make([]RobotInstruction, 2)

	// move

	mov := NewMoveInstruction()
	mov.Head = ins.Head
	mov.Pos = ins.FPosition
	mov.Well = ins.Well
	mov.Plt = ins.FPlateType
	mov.Reference = append(mov.Reference, 0)
	mov.OffsetX = append(mov.OffsetX, 0.0)
	mov.OffsetY = append(mov.OffsetY, 0.0)
	mov.OffsetZ = append(mov.OffsetZ, 0.0)
	mov.Platform = ins.Platform
	ret[0] = mov

	// load tips

	lod := NewLoadTipsInstruction()
	lod.Head = ins.Head
	lod.TipType = ins.FPlateType
	lod.HolderType = ins.FPlateType
	lod.Multi = ins.Multi
	lod.Pos = ins.FPosition
	lod.HolderType = ins.FPlateType
	lod.Well = ins.Well
	lod.Platform = ins.Platform
	ret[1] = lod

	return ret, nil
}

type UnloadTipsMoveInstruction struct {
	GenericRobotInstruction
	Type       int
	Head       int
	PltTo      []string
	WellTo     []string
	TPlateType []string
	Multi      int
	Platform   string
}

func NewUnloadTipsMoveInstruction() *UnloadTipsMoveInstruction {
	var v UnloadTipsMoveInstruction
	v.Type = UDT
	v.PltTo = make([]string, 0)
	v.WellTo = make([]string, 0)
	v.TPlateType = make([]string, 0)
	return &v
}
func (ins *UnloadTipsMoveInstruction) InstructionType() int {
	return ins.Type
}

func (ins *UnloadTipsMoveInstruction) GetParameter(name string) interface{} {
	switch name {
	case "HEAD":
		return ins.Head
	case "TOPLATETYPE":
		return ins.TPlateType
	case "POSTO":
		return ins.PltTo
	case "WELLTO":
		return ins.WellTo
	case "INSTRUCTIONTYPE":
		return ins.InstructionType()
	case "MULTI":
		return ins.Multi
	case "PLATFORM":
		return ins.Platform
	}
	return nil
}

func (ins *UnloadTipsMoveInstruction) Generate(policy *wtype.LHPolicyRuleSet, prms *LHProperties) ([]RobotInstruction, error) {
	ret := make([]RobotInstruction, 2)

	// move

	mov := NewMoveInstruction()
	mov.Head = ins.Head
	mov.Pos = ins.PltTo
	mov.Well = ins.WellTo
	mov.Plt = ins.TPlateType
	mov.Reference = append(mov.Reference, 0)
	mov.OffsetX = append(mov.OffsetX, 0.0)
	mov.OffsetY = append(mov.OffsetY, 0.0)
	mov.OffsetZ = append(mov.OffsetZ, 0.0)
	mov.Platform = ins.Platform
	ret[0] = mov

	// unload tips

	uld := NewUnloadTipsInstruction()
	uld.Head = ins.Head
	uld.TipType = ins.TPlateType
	uld.HolderType = ins.TPlateType
	uld.Multi = ins.Multi
	uld.Pos = ins.PltTo
	uld.HolderType = ins.TPlateType
	uld.Well = ins.WellTo
	uld.Platform = ins.Platform
	ret[1] = uld

	return ret, nil
}

type AspirateInstruction struct {
	GenericRobotInstruction
	Type       int
	Head       int
	Volume     []wunit.Volume
	Overstroke bool
	Multi      int
	Plt        []string
	What       []string
	LLF        []bool
	Platform   string
}

func NewAspirateInstruction() *AspirateInstruction {
	var v AspirateInstruction
	v.Type = ASP
	v.Volume = make([]wunit.Volume, 0)
	v.Plt = make([]string, 0)
	v.What = make([]string, 0)
	v.LLF = make([]bool, 0)
	return &v
}
func (ins *AspirateInstruction) InstructionType() int {
	return ins.Type
}

func (ins *AspirateInstruction) GetParameter(name string) interface{} {
	switch name {
	case "VOLUME":
		return ins.Volume
	case "HEAD":
		return ins.Head
	case "MULTI":
		return ins.Multi
	case "OVERSTROKE":
		return ins.Overstroke
	case "INSTRUCTIONTYPE":
		return ins.InstructionType()
	case "WHAT":
		return ins.What
	case "PLATE":
		return ins.Plt
	case "LLF":
		return ins.LLF
	case "PLATFORM":
		return ins.Platform
	}
	return nil
}

func (ins *AspirateInstruction) Generate(policy *wtype.LHPolicyRuleSet, prms *LHProperties) ([]RobotInstruction, error) {
	return nil, nil
}

func (ins *AspirateInstruction) OutputTo(driver LiquidhandlingDriver) error {
	volumes := make([]float64, len(ins.Volume))
	for i, vol := range ins.Volume {
		volumes[i] = vol.ConvertTo(wunit.ParsePrefixedUnit("ul"))
	}
	os := []bool{ins.Overstroke}

	ret := driver.Aspirate(volumes, os, ins.Head, ins.Multi, ins.Plt, ins.What, ins.LLF)
	if !ret.OK {
		return fmt.Errorf(" %d : %s", ret.Errorcode, ret.Msg)
	}

	return nil
}

type DispenseInstruction struct {
	GenericRobotInstruction
	Type     int
	Head     int
	Volume   []wunit.Volume
	Multi    int
	Plt      []string
	What     []string
	LLF      []bool
	Platform string
}

func NewDispenseInstruction() *DispenseInstruction {
	var v DispenseInstruction
	v.Type = DSP
	v.Volume = make([]wunit.Volume, 0)
	v.Plt = make([]string, 0)
	v.What = make([]string, 0)
	v.LLF = make([]bool, 0)
	return &v
}
func (ins *DispenseInstruction) InstructionType() int {
	return ins.Type
}

func (ins *DispenseInstruction) GetParameter(name string) interface{} {
	switch name {
	case "VOLUME":
		return ins.Volume
	case "HEAD":
		return ins.Head
	case "MULTI":
		return ins.Multi
	case "INSTRUCTIONTYPE":
		return ins.InstructionType()
	case "WHAT":
		return ins.What
	case "LLF":
		return ins.LLF
	case "PLT":
		return ins.Plt
	case "PLATFORM":
		return ins.Platform
	}
	return nil
}

func (ins *DispenseInstruction) Generate(policy *wtype.LHPolicyRuleSet, prms *LHProperties) ([]RobotInstruction, error) {
	return nil, nil
}

func (ins *DispenseInstruction) OutputTo(driver LiquidhandlingDriver) error {
	volumes := make([]float64, len(ins.Volume))
	for i, vol := range ins.Volume {
		volumes[i] = vol.ConvertTo(wunit.ParsePrefixedUnit("ul"))
	}

	os := []bool{false}
	ret := driver.Dispense(volumes, os, ins.Head, ins.Multi, ins.Plt, ins.What, ins.LLF)
	if !ret.OK {
		return fmt.Errorf(" %d : %s", ret.Errorcode, ret.Msg)
	}

	return nil

}

type BlowoutInstruction struct {
	GenericRobotInstruction
	Type     int
	Head     int
	Volume   []wunit.Volume
	Multi    int
	Plt      []string
	What     []string
	LLF      []bool
	Platform string
}

func NewBlowoutInstruction() *BlowoutInstruction {
	var v BlowoutInstruction
	v.Type = BLO
	v.Volume = make([]wunit.Volume, 0)
	return &v
}
func (ins *BlowoutInstruction) InstructionType() int {
	return ins.Type
}

func (ins *BlowoutInstruction) GetParameter(name string) interface{} {
	switch name {
	case "VOLUME":
		return ins.Volume
	case "HEAD":
		return ins.Head
	case "MULTI":
		return ins.Multi
	case "INSTRUCTIONTYPE":
		return ins.InstructionType()
	case "WHAT":
		return ins.What
	case "LLF":
		return ins.LLF
	case "PLT":
		return ins.Plt
	case "PLATFORM":
		return ins.Platform
	}
	return nil
}

func (ins *BlowoutInstruction) Generate(policy *wtype.LHPolicyRuleSet, prms *LHProperties) ([]RobotInstruction, error) {
	return nil, nil
}

func (ins *BlowoutInstruction) OutputTo(driver LiquidhandlingDriver) error {
	volumes := make([]float64, len(ins.Volume))
	for i, vol := range ins.Volume {
		volumes[i] = vol.ConvertTo(wunit.ParsePrefixedUnit("ul"))
	}
	bo := make([]bool, ins.Multi)
	for i := 0; i < ins.Multi; i++ {
		bo[i] = true
	}
	ret := driver.Dispense(volumes, bo, ins.Head, ins.Multi, ins.Plt, ins.What, ins.LLF)
	if !ret.OK {
		return fmt.Errorf(" %d : %s", ret.Errorcode, ret.Msg)
	}

	return nil
}

type PTZInstruction struct {
	GenericRobotInstruction
	Type    int
	Head    int
	Channel int
}

func NewPTZInstruction() *PTZInstruction {
	var v PTZInstruction
	v.Type = PTZ
	return &v
}
func (ins *PTZInstruction) InstructionType() int {
	return ins.Type
}

func (ins *PTZInstruction) GetParameter(name string) interface{} {
	switch name {
	case "HEAD":
		return ins.Head
	case "CHANNEL":
		return ins.Channel
	case "INSTRUCTIONTYPE":
		return ins.InstructionType()
	}
	return nil
}

func (ins *PTZInstruction) Generate(policy *wtype.LHPolicyRuleSet, prms *LHProperties) ([]RobotInstruction, error) {
	return nil, nil
}

func (ins *PTZInstruction) OutputTo(driver LiquidhandlingDriver) error {
	ret := driver.ResetPistons(ins.Head, ins.Channel)
	if !ret.OK {
		return fmt.Errorf(" %d : %s", ret.Errorcode, ret.Msg)
	}

	return nil
}

type MoveInstruction struct {
	GenericRobotInstruction
	Type      int
	Head      int
	Pos       []string
	Plt       []string
	Well      []string
	WVolume   []wunit.Volume
	Reference []int
	OffsetX   []float64
	OffsetY   []float64
	OffsetZ   []float64
	Platform  string
}

func NewMoveInstruction() *MoveInstruction {
	var v MoveInstruction
	v.Type = MOV
	v.Plt = make([]string, 0)
	v.Pos = make([]string, 0)
	v.Well = make([]string, 0)
	v.WVolume = make([]wunit.Volume, 0)
	v.Reference = make([]int, 0)
	v.OffsetX = make([]float64, 0)
	v.OffsetY = make([]float64, 0)
	v.OffsetZ = make([]float64, 0)
	return &v
}
func (ins *MoveInstruction) InstructionType() int {
	return ins.Type
}

func (ins *MoveInstruction) GetParameter(name string) interface{} {
	switch name {
	case "TOWELLVOLUME":
		return ins.WVolume
	case "HEAD":
		return ins.Head
	case "TOPLATETYPE":
		return ins.Plt
	case "POSTO":
		return ins.Pos
	case "WELLTO":
		return ins.Well
	case "REFERENCE":
		return ins.Reference
	case "OFFSETX":
		return ins.OffsetX
	case "OFFSETY":
		return ins.OffsetY
	case "OFFSETZ":
		return ins.OffsetZ
	case "INSTRUCTIONTYPE":
		return ins.InstructionType()
	case "PLATFORM":
		return ins.Platform
	}
	return nil
}

func (ins *MoveInstruction) Generate(policy *wtype.LHPolicyRuleSet, prms *LHProperties) ([]RobotInstruction, error) {
	return nil, nil
}

func (ins *MoveInstruction) OutputTo(driver LiquidhandlingDriver) error {
	ret := driver.Move(ins.Pos, ins.Well, ins.Reference, ins.OffsetX, ins.OffsetY, ins.OffsetZ, ins.Plt, ins.Head)
	if !ret.OK {
		return fmt.Errorf(" %d : %s", ret.Errorcode, ret.Msg)
	}

	return nil

}

type MoveRawInstruction struct {
	GenericRobotInstruction
	Type       int
	Head       int
	What       []string
	PltFrom    []string
	PltTo      []string
	WellFrom   []string
	WellTo     []string
	Volume     []wunit.Volume
	FPlateType []string
	TPlateType []string
	FVolume    []wunit.Volume
	TVolume    []wunit.Volume
	Prms       *wtype.LHChannelParameter
}

func NewMoveRawInstruction() *MoveRawInstruction {
	var v MoveRawInstruction
	v.Type = MRW
	v.What = make([]string, 0)
	v.PltFrom = make([]string, 0)
	v.PltTo = make([]string, 0)
	v.WellFrom = make([]string, 0)
	v.WellTo = make([]string, 0)
	v.FPlateType = make([]string, 0)
	v.TPlateType = make([]string, 0)
	v.Volume = make([]wunit.Volume, 0)
	v.FVolume = make([]wunit.Volume, 0)
	v.TVolume = make([]wunit.Volume, 0)
	return &v
}
func (ins *MoveRawInstruction) InstructionType() int {
	return ins.Type
}

func (ins *MoveRawInstruction) GetParameter(name string) interface{} {
	switch name {
	case "HEAD":
		return ins.Head
	case "LIQUIDCLASS":
		return ins.What
	case "VOLUME":
		return ins.Volume
	case "FROMPLATETYPE":
		return ins.FPlateType
	case "TOPLATETYPE":
		return ins.TPlateType
	case "WELLFROMVOLUME":
		return ins.FVolume
	case "WELLTOVOLUME":
		return ins.TVolume
	case "POSFROM":
		return ins.PltFrom
	case "POSTO":
		return ins.PltTo
	case "WELLFROM":
		return ins.WellFrom
	case "PARAMS":
		return ins.Prms
	case "PLATFORM":
		if ins.Prms == nil {
			return ""
		}
		return ins.Prms.Platform
	case "INSTRUCTIONTYPE":
		return ins.InstructionType()
	}
	return nil
}

func (ins *MoveRawInstruction) Generate(policy *wtype.LHPolicyRuleSet, prms *LHProperties) ([]RobotInstruction, error) {
	return nil, nil
}

func (ins *MoveRawInstruction) OutputTo(driver LiquidhandlingDriver) error {
	logger.Fatal("Not yet implemented")
	panic("Not yet implemented")
}

type LoadTipsInstruction struct {
	GenericRobotInstruction
	Type       int
	Head       int
	Pos        []string
	Well       []string
	Channels   []int
	TipType    []string
	HolderType []string
	Multi      int
	Platform   string
}

func NewLoadTipsInstruction() *LoadTipsInstruction {
	var v LoadTipsInstruction
	v.Type = LOD
	v.Channels = make([]int, 0)
	v.TipType = make([]string, 0)
	v.HolderType = make([]string, 0)
	v.Pos = make([]string, 0)
	v.Well = make([]string, 0)
	return &v
}
func (ins *LoadTipsInstruction) InstructionType() int {
	return ins.Type
}

func (ins *LoadTipsInstruction) GetParameter(name string) interface{} {
	switch name {
	case "HEAD":
		return ins.Head
	case "CHANNEL":
		return ins.Channels
	case "TIPTYPE":
		return ins.TipType
	case "FROMPLATETYPE":
		return ins.HolderType
	case "INSTRUCTIONTYPE":
		return ins.InstructionType()
	case "MULTI":
		return ins.Multi
	case "WELL":
		return ins.Well
	case "PLATE":
		return ins.HolderType
	case "POS":
		return ins.Pos
	case "PLATFORM":
		return ins.Platform
	}
	return nil
}

func (ins *LoadTipsInstruction) Generate(policy *wtype.LHPolicyRuleSet, prms *LHProperties) ([]RobotInstruction, error) {
	return nil, nil
}

func (ins *LoadTipsInstruction) OutputTo(driver LiquidhandlingDriver) error {
	ret := driver.LoadTips(ins.Channels, ins.Head, len(ins.TipType), ins.HolderType, ins.Pos, ins.Well)
	if !ret.OK {
		return fmt.Errorf(" %d : %s", ret.Errorcode, ret.Msg)
	}

	return nil

}

type UnloadTipsInstruction struct {
	GenericRobotInstruction
	Type       int
	Head       int
	Channels   []int
	TipType    []string
	HolderType []string
	Multi      int
	Pos        []string
	Well       []string
	Platform   string
}

func NewUnloadTipsInstruction() *UnloadTipsInstruction {
	var v UnloadTipsInstruction
	v.Type = ULD
	v.TipType = make([]string, 0)
	v.HolderType = make([]string, 0)
	v.Channels = make([]int, 0)
	v.Pos = make([]string, 0)
	v.Well = make([]string, 0)
	return &v
}
func (ins *UnloadTipsInstruction) InstructionType() int {
	return ins.Type
}

func (ins *UnloadTipsInstruction) GetParameter(name string) interface{} {
	switch name {
	case "HEAD":
		return ins.Head
	case "CHANNEL":
		return ins.Channels
	case "TIPTYPE":
		return ins.TipType
	case "TOPLATETYPE":
		return ins.HolderType
	case "INSTRUCTIONTYPE":
		return ins.InstructionType()
	case "MULTI":
		return ins.Multi
	case "WELL":
		return ins.Well
	case "POS":
		return ins.Pos
	case "PLATFORM":
		return ins.Platform
	}
	return nil
}

func (ins *UnloadTipsInstruction) Generate(policy *wtype.LHPolicyRuleSet, prms *LHProperties) ([]RobotInstruction, error) {
	return nil, nil
}

func (ins *UnloadTipsInstruction) OutputTo(driver LiquidhandlingDriver) error {
	ret := driver.UnloadTips(ins.Channels, ins.Head, len(ins.TipType), ins.HolderType, ins.Pos, ins.Well)
	if !ret.OK {
		return fmt.Errorf(" %d : %s", ret.Errorcode, ret.Msg)
	}

	return nil

}

type SuckInstruction struct {
	GenericRobotInstruction
	Type       int
	Head       int
	What       []string
	PltFrom    []string
	WellFrom   []string
	Volume     []wunit.Volume
	FPlateType []string
	FVolume    []wunit.Volume
	Prms       *wtype.LHChannelParameter
	Multi      int
	Overstroke bool
	TipType    string
}

func NewSuckInstruction() *SuckInstruction {
	var v SuckInstruction
	v.Type = SUK
	v.What = make([]string, 0)
	v.PltFrom = make([]string, 0)
	v.WellFrom = make([]string, 0)
	v.Volume = make([]wunit.Volume, 0)
	v.FPlateType = make([]string, 0)
	v.FVolume = make([]wunit.Volume, 0)
	v.GenericRobotInstruction.Ins = &v
	return &v
}
func (ins *SuckInstruction) InstructionType() int {
	return ins.Type
}

func (ins *SuckInstruction) AddTransferParams(tp TransferParams) {
	ins.What = append(ins.What, tp.What)
	ins.PltFrom = append(ins.PltFrom, tp.PltFrom)
	ins.WellFrom = append(ins.WellFrom, tp.WellFrom)
	ins.Volume = append(ins.Volume, tp.Volume)
	ins.FPlateType = append(ins.FPlateType, tp.FPlateType)
	ins.FVolume = append(ins.FVolume, tp.FVolume)
	ins.Prms = tp.Channel
	ins.Head = tp.Channel.Head
	ins.TipType = tp.TipType
}

func (ins *SuckInstruction) GetParameter(name string) interface{} {
	switch name {
	case "HEAD":
		return ins.Head
	case "LIQUIDCLASS":
		return ins.What
	case "VOLUME":
		return ins.Volume
	case "FROMPLATETYPE":
		return ins.FPlateType
	case "WELLFROMVOLUME":
		return ins.FVolume
	case "POSFROM":
		return ins.PltFrom
	case "WELLFROM":
		return ins.WellFrom
	case "PARAMS":
		return ins.Prms
	case "MULTI":
		return ins.Multi
	case "OVERSTROKE":
		return ins.Overstroke
	case "PLATFORM":
		if ins.Prms == nil {
			return ""
		}
		return ins.Prms.Platform
	case "INSTRUCTIONTYPE":
		return ins.InstructionType()
	case "TIPTYPE":
		return ins.TipType
	}
	return nil
}

func (ins *SuckInstruction) Generate(policy *wtype.LHPolicyRuleSet, prms *LHProperties) ([]RobotInstruction, error) {
	ret := make([]RobotInstruction, 0, 1)

	// this is where the policies come into effect

	pol := GetPolicyFor(policy, ins)

	// offsets
	ofx := SafeGetF64(pol, "ASPXOFFSET")
	ofy := SafeGetF64(pol, "ASPYOFFSET")
	ofz := SafeGetF64(pol, "ASPZOFFSET")
	ofzadj := SafeGetF64(pol, "OFFSETZADJUST")
	ofz += ofzadj

	mixofx := SafeGetF64(pol, "PRE_MIX_X")
	mixofy := SafeGetF64(pol, "PRE_MIX_Y")
	mixofz := SafeGetF64(pol, "PRE_MIX_Z")
	mixofz += ofzadj
	final_asp_ref := SafeGetInt(pol, "ASPREFERENCE")

	pspeed := SafeGetF64(pol, "DEFAULTPIPETTESPEED")

	// do we need to enter slowly?
	entryspeed, gentlynow := pol["ASPENTRYSPEED"]
	if gentlynow {
		// go to the well top
		mov := NewMoveInstruction()

		mov.Head = ins.Head
		mov.Pos = ins.PltFrom
		mov.Plt = ins.FPlateType
		mov.Well = ins.WellFrom
		mov.WVolume = ins.FVolume
		for i := 0; i < ins.Multi; i++ {
			mov.Reference = append(mov.Reference, 1)
			mov.OffsetX = append(mov.OffsetX, ofx)
			mov.OffsetY = append(mov.OffsetY, ofy)
			mov.OffsetZ = append(mov.OffsetZ, 5.0)
		}
		ret = append(ret, mov)

		// set the speed
		spd := NewSetDriveSpeedInstruction()
		spd.Drive = "Z"
		spd.Speed = entryspeed.(float64)
		ret = append(ret, spd)

		/*
			// now move into the liquid
			mov = NewMoveInstruction()
			mov.Head = ins.Head
			mov.Pos = ins.PltFrom
			mov.Plt = ins.FPlateType
			mov.Well = ins.WellFrom
			mov.WVolume = ins.FVolume
			ref := SafeGetInt(pol, "ASPREFERENCE")
			ofx := SafeGetF64(pol, "ASPXOFFSET")
			ofy := SafeGetF64(pol, "ASPYOFFSET")
			ofz := SafeGetF64(pol, "ASPZOFFSET")

			// TODO -- different offsets per channel

			for i := 0; i < ins.Multi; i++ {
				mov.Reference = append(mov.Reference, ref)
				mov.OffsetX = append(mov.OffsetX, ofx)
				mov.OffsetY = append(mov.OffsetY, ofy)
				mov.OffsetZ = append(mov.OffsetZ, ofz)
			}

			ret = append(ret, mov)
			// reset the drive speed
			spd = NewSetDriveSpeedInstruction()
			spd.Drive = "Z"
			spd.Speed = pol["DEFAULTZSPEED"].(float64)
			ret = append(ret, spd)
		*/
	}

	// do we pre-mix?
	_, premix := pol["PRE_MIX"]
	cycles := SafeGetInt(pol, "PRE_MIX")

	if premix && cycles > 0 {
		// add the premix step
		mix := NewMoveMixInstruction()
		mix.Head = ins.Head
		mix.Plt = ins.PltFrom
		mix.PlateType = ins.FPlateType
		mix.Well = ins.WellFrom
		mix.Multi = ins.Multi
		mix.What = ins.What
		// TODO get rid of this HARD CODE
		mix.Blowout = []bool{false}

		_, ok := pol["PRE_MIX_VOLUME"]
		mix.Volume = ins.Volume
		mixvol := SafeGetF64(pol, "PRE_MIX_VOLUME")
		vmixvol := wunit.NewVolume(mixvol, "ul")

		// TODO -- corresponding checks when set
		if mixvol < wtype.Globals.MIN_REASONABLE_VOLUME_UL {
			return ret, wtype.LHError(wtype.LH_ERR_POLICY, fmt.Sprintf("PRE_MIX_VOLUME set below minimum allowed: %f min %f", mixvol, wtype.Globals.MIN_REASONABLE_VOLUME_UL))
		} else if !ins.Prms.CanMove(vmixvol, true) {
			// this is an error in channel choice but the user has to deal... needs modificationst
			return ret, wtype.LHError(wtype.LH_ERR_POLICY, fmt.Sprintf("PRE_MIX_VOLUME not compatible with optimal channel choice: requested %s channel limits are %s", vmixvol.ToString(), ins.Prms.VolumeLimitString()))
		}

		if ok {
			v := make([]wunit.Volume, ins.Multi)
			for i := 0; i < ins.Multi; i++ {
				vl := wunit.NewVolume(mixvol, "ul")
				v[i] = vl
			}
			mix.Volume = v
		}
		// offsets

		for k := 0; k < ins.Multi; k++ {
			mix.OffsetX = append(mix.OffsetX, mixofx)
		}

		for k := 0; k < ins.Multi; k++ {
			mix.OffsetY = append(mix.OffsetY, mixofy)
		}
		for k := 0; k < ins.Multi; k++ {
			mix.OffsetZ = append(mix.OffsetZ, mixofz)
		}

		c := make([]int, ins.Multi)

		for i := 0; i < ins.Multi; i++ {
			c[i] = cycles
		}

		// set speed

		//_, changepipspeed := pol["PRE_MIX_RATE"]

		mixrate := SafeGetF64(pol, "PRE_MIX_RATE")

		changepipspeed := (mixrate == pspeed) && (mixrate > 0.0)

		if changepipspeed {
			setspd := NewSetPipetteSpeedInstruction()
			setspd.Head = ins.Head
			setspd.Channel = -1 // all channels
			setspd.Speed = mixrate
			ret = append(ret, setspd)
		}

		mix.Cycles = c
		ret = append(ret, mix)

		if changepipspeed {
			sps := NewSetPipetteSpeedInstruction()
			sps.Head = ins.Head
			sps.Channel = -1 // all channels
			sps.Speed = pspeed
			ret = append(ret, sps)
		}
	}

	/*
		discrepancy := false

		if premix {
			// check whether there is a discrepancy between the mix reference
			// etc. and the asp reference... if not we don't need to move

			discrepancy = discrepancy || (mixofx != ofx)
			discrepancy = discrepancy || (mixofy != ofy)
			discrepancy = discrepancy || (mixofz != ofz)
		}
	*/
	//nb moves are mandatory
	mov := NewMoveInstruction()
	mov.Head = ins.Head

	mov.Pos = ins.PltFrom
	mov.Plt = ins.FPlateType
	mov.Well = ins.WellFrom
	mov.WVolume = ins.FVolume

	for i := 0; i < ins.Multi; i++ {
		mov.Reference = append(mov.Reference, final_asp_ref)
		mov.OffsetX = append(mov.OffsetX, ofx)
		mov.OffsetY = append(mov.OffsetY, ofy)
		mov.OffsetZ = append(mov.OffsetZ, ofz)
	}
	ret = append(ret, mov)

	// Set the pipette speed if needed

	apspeed := SafeGetF64(pol, "ASPSPEED")

	changepspeed := (apspeed == pspeed) && (apspeed > 0.0)

	if changepspeed {
		sps := NewSetPipetteSpeedInstruction()
		sps.Head = ins.Head
		sps.Channel = -1 // all channels
		sps.Speed = apspeed
		ret = append(ret, sps)
	}

	// now we aspirate

	aspins := NewAspirateInstruction()
	aspins.Head = ins.Head
	aspins.Volume = ins.Volume

	ev, iwantmore := pol["EXTRA_ASP_VOLUME"]
	if iwantmore {
		extra_vol := ev.(wunit.Volume)
		for i, _ := range aspins.Volume {
			aspins.Volume[i].Add(extra_vol)
		}
	}

	aspins.Multi = ins.Multi
	aspins.Overstroke = ins.Overstroke
	aspins.What = ins.What
	aspins.Plt = ins.FPlateType

	for i := 0; i < ins.Multi; i++ {
		aspins.LLF = append(aspins.LLF, false)
	}

	ret = append(ret, aspins)

	// do we reset the pipette speed?

	if changepspeed {
		sps := NewSetPipetteSpeedInstruction()
		sps.Head = ins.Head
		sps.Channel = -1 // all channels
		sps.Speed = pspeed
		ret = append(ret, sps)
	}

	// do we wait

	wait_time, wait := pol["ASP_WAIT"]

	if wait {
		waitins := NewWaitInstruction()
		waitins.Time = wait_time.(float64)
		ret = append(ret, waitins)
	}

	if gentlynow { // reset the drive speed
		// go to the well top
		mov := NewMoveInstruction()

		mov.Head = ins.Head
		mov.Pos = ins.PltFrom
		mov.Plt = ins.FPlateType
		mov.Well = ins.WellFrom
		mov.WVolume = ins.FVolume
		for i := 0; i < ins.Multi; i++ {
			mov.Reference = append(mov.Reference, 1)
			mov.OffsetX = append(mov.OffsetX, ofx)
			mov.OffsetY = append(mov.OffsetY, ofy)
			mov.OffsetZ = append(mov.OffsetZ, 5.0)
		}
		ret = append(ret, mov)

		// now get ready to move fast again
		spd := NewSetDriveSpeedInstruction()
		spd.Drive = "Z"
		spd.Speed = pol["DEFAULTZSPEED"].(float64)
		ret = append(ret, spd)
	}

	return ret, nil

}

type BlowInstruction struct {
	GenericRobotInstruction
	Type       int
	Head       int
	What       []string
	PltTo      []string
	WellTo     []string
	Volume     []wunit.Volume
	TPlateType []string
	TVolume    []wunit.Volume
	Prms       *wtype.LHChannelParameter
	Multi      int
	TipType    string
}

func NewBlowInstruction() *BlowInstruction {
	var v BlowInstruction
	v.Type = BLW
	v.What = make([]string, 0)
	v.PltTo = make([]string, 0)
	v.WellTo = make([]string, 0)
	v.Volume = make([]wunit.Volume, 0)
	v.TPlateType = make([]string, 0)
	v.TVolume = make([]wunit.Volume, 0)
	v.GenericRobotInstruction.Ins = &v
	return &v
}
func (ins *BlowInstruction) InstructionType() int {
	return ins.Type
}

func (ins *BlowInstruction) GetParameter(name string) interface{} {
	switch name {
	case "HEAD":
		return ins.Head
	case "LIQUIDCLASS":
		return ins.What
	case "VOLUME":
		return ins.Volume
	case "TOPLATETYPE":
		return ins.TPlateType
	case "WELLTOVOLUME":
		return ins.TVolume
	case "POSTO":
		return ins.PltTo
	case "WELLTO":
		return ins.WellTo
	case "PARAMS":
		return ins.Prms
	case "PLATFORM":
		if ins.Prms == nil {
			return ""
		}
		return ins.Prms.Platform
	case "MULTI":
		return ins.Multi
	case "INSTRUCTIONTYPE":
		return ins.InstructionType()
	case "TIPTYPE":
		return ins.TipType
	}
	return nil
}

func (ins *BlowInstruction) AddTransferParams(tp TransferParams) {
	ins.What = append(ins.What, tp.What)
	ins.PltTo = append(ins.PltTo, tp.PltTo)
	ins.WellTo = append(ins.WellTo, tp.WellTo)
	ins.Volume = append(ins.Volume, tp.Volume)
	ins.TPlateType = append(ins.TPlateType, tp.TPlateType)
	ins.TVolume = append(ins.TVolume, tp.TVolume)
	ins.Head = tp.Channel.Head
	ins.TipType = tp.TipType
}
func (scti *BlowInstruction) Params() MultiTransferParams {
	var tp MultiTransferParams
	tp.What = scti.What
	tp.PltTo = scti.PltTo
	tp.WellTo = scti.WellTo
	tp.Volume = scti.Volume
	tp.TPlateType = scti.TPlateType
	tp.TVolume = scti.TVolume
	tp.Channel = scti.Prms
	return tp
}

func (ins *BlowInstruction) Generate(policy *wtype.LHPolicyRuleSet, prms *LHProperties) ([]RobotInstruction, error) {
	ret := make([]RobotInstruction, 0)
	// apply policies here

	pol := GetPolicyFor(policy, ins)
	// first, are we breaking up the move?

	ofx := SafeGetF64(pol, "DSPXOFFSET")
	ofy := SafeGetF64(pol, "DSPYOFFSET")
	ofz := SafeGetF64(pol, "DSPZOFFSET")
	ofzadj := SafeGetF64(pol, "OFFSETZADJUST")

	ofz += ofzadj

	ref := SafeGetInt(pol, "DSPREFERENCE")
	entryspeed := SafeGetF64(pol, "DSPENTRYSPEED")
	defaultspeed := SafeGetF64(pol, "DEFAULTZSPEED")

	var gentlydoesit bool

	if entryspeed != defaultspeed {
		gentlydoesit = true
	}

	if gentlydoesit {
		// go to the well top
		mov := NewMoveInstruction()

		mov.Head = ins.Head
		mov.Pos = ins.PltTo
		mov.Plt = ins.TPlateType
		mov.Well = ins.WellTo
		mov.WVolume = ins.TVolume
		for i := 0; i < ins.Multi; i++ {
			mov.Reference = append(mov.Reference, 1)
			mov.OffsetX = append(mov.OffsetX, ofx)
			mov.OffsetY = append(mov.OffsetY, ofy)
			mov.OffsetZ = append(mov.OffsetZ, 5.0)
		}
		ret = append(ret, mov)

		// set the speed
		spd := NewSetDriveSpeedInstruction()
		spd.Drive = "Z"
		spd.Speed = entryspeed
		ret = append(ret, spd)

		/*
			mov = NewMoveInstruction()
			mov.Head = ins.Head
			mov.Pos = ins.PltTo
			mov.Plt = ins.TPlateType
			mov.Well = ins.WellTo
			mov.WVolume = ins.TVolume
			for i := 0; i < ins.Multi; i++ {
				mov.Reference = append(mov.Reference, pol["DSPREFERENCE"].(int))
				mov.OffsetX = append(mov.OffsetX, 0.0)
				mov.OffsetY = append(mov.OffsetY, 0.0)
				mov.OffsetZ = append(mov.OffsetZ, pol["DSPZOFFSET"].(float64))
			}
			ret = append(ret, mov)
			// reset the drive speed
			spd = NewSetDriveSpeedInstruction()
			spd.Drive = "Z"
			spd.Speed = pol["DEFAULTZSPEED"].(float64)
			ret = append(ret, spd)
		*/

	}

	mov := NewMoveInstruction()
	mov.Head = ins.Head
	mov.Pos = ins.PltTo
	mov.Plt = ins.TPlateType
	mov.Well = ins.WellTo
	mov.WVolume = ins.TVolume
	for i := 0; i < ins.Multi; i++ {
		mov.Reference = append(mov.Reference, ref)
		mov.OffsetX = append(mov.OffsetX, ofx)
		mov.OffsetY = append(mov.OffsetY, ofy)
		mov.OffsetZ = append(mov.OffsetZ, ofz)
	}

	ret = append(ret, mov)

	// change pipette speed?
	pspeed := SafeGetF64(pol, "DEFAULTPIPETTESPEED")
	dpspeed := SafeGetF64(pol, "DSPSPEED")

	var setpspeed bool

	if pspeed != dpspeed && dpspeed != 0.0 {
		setpspeed = true
	}

	if setpspeed {
		sps := NewSetPipetteSpeedInstruction()
		sps.Head = ins.Head
		sps.Channel = -1 // all channels
		sps.Speed = dpspeed
		ret = append(ret, sps)
	}

	// now we dispense

	weneedtoreset := true

	justblowout := SafeGetBool(pol, "JUSTBLOWOUT")

	// bit of a horrible hack... we need to fix this soon

	if justblowout {
		blowoutvolume := SafeGetF64(pol, "BLOWOUTVOLUME")
		blowoutvolunit := SafeGetString(pol, "BLOWOUTVOLUMEUNIT")

		// be safe, not sorry...

		if blowoutvolunit == "" {
			blowoutvolunit = "ul"
		}

		boins := NewBlowoutInstruction()
		boins.Head = ins.Head
		vl := wunit.NewVolume(blowoutvolume, blowoutvolunit)
		boins.Volume = append(boins.Volume, vl)
		boins.Multi = ins.Multi
		boins.Plt = ins.TPlateType
		boins.What = ins.What

		for i := 0; i < ins.Multi; i++ {
			boins.LLF = append(boins.LLF, false)
		}

		ret = append(ret, boins)
		weneedtoreset = false
	} else {
		dspins := NewDispenseInstruction()
		dspins.Head = ins.Head
		dspins.Volume = ins.Volume

		extra_vol := SafeGetVolume(pol, "EXTRA_DISP_VOLUME")
		if extra_vol.GreaterThan(wunit.ZeroVolume()) {
			for i, _ := range dspins.Volume {
				dspins.Volume[i].Add(extra_vol)
			}
		}
		dspins.Multi = ins.Multi
		dspins.Plt = ins.TPlateType
		dspins.What = ins.What

		for i := 0; i < ins.Multi; i++ {
			dspins.LLF = append(dspins.LLF, false)
		}

		ret = append(ret, dspins)
	}

	// do we reset the pipette speed?

	if setpspeed {
		sps := NewSetPipetteSpeedInstruction()
		sps.Head = ins.Head
		sps.Channel = -1 // all channels
		sps.Speed = pspeed
		ret = append(ret, sps)
	}

	// do we wait?

	wait_time := SafeGetF64(pol, "DSP_WAIT")

	if wait_time > 0.0 {
		waitins := NewWaitInstruction()
		waitins.Time = wait_time
		ret = append(ret, waitins)
	}

	// do we mix?
	_, postmix := pol["POST_MIX"]
	cycles := SafeGetInt(pol, "POST_MIX")

	if postmix && cycles > 0 {
		// add the postmix step
		mix := NewMoveMixInstruction()
		mix.Head = ins.Head
		mix.Plt = ins.PltTo
		mix.PlateType = ins.TPlateType
		mix.Well = ins.WellTo
		mix.Multi = ins.Multi
		mix.What = ins.What
		// TODO get rid of this HARD CODE
		// we might want to change this
		b := make([]bool, ins.Multi)
		mix.Blowout = b

		// offsets

		pmxoff := SafeGetF64(pol, "POST_MIX_X")

		for k := 0; k < ins.Multi; k++ {
			mix.OffsetX = append(mix.OffsetX, pmxoff)
		}

		pmyoff := SafeGetF64(pol, "POST_MIX_Y")
		for k := 0; k < ins.Multi; k++ {
			mix.OffsetY = append(mix.OffsetY, pmyoff)
		}

		pmzoff := SafeGetF64(pol, "POST_MIX_Z")
		pmzoff += ofzadj

		for k := 0; k < ins.Multi; k++ {
			mix.OffsetZ = append(mix.OffsetZ, pmzoff)
		}

		_, ok := pol["POST_MIX_VOLUME"]
		mix.Volume = ins.Volume
		mixvol := SafeGetF64(pol, "POST_MIX_VOLUME")
		vmixvol := wunit.NewVolume(mixvol, "ul")

		// check the volume

		if mixvol < wtype.Globals.MIN_REASONABLE_VOLUME_UL {
			return ret, wtype.LHError(wtype.LH_ERR_POLICY, fmt.Sprintf("POST_MIX_VOLUME set below minimum allowed: %f min %f", mixvol, wtype.Globals.MIN_REASONABLE_VOLUME_UL))
		} else if !ins.Prms.CanMove(vmixvol, true) {
			// make this illegal for now

			return ret, wtype.LHError(wtype.LH_ERR_POLICY, fmt.Sprintf("Setting POST_MIX_VOLME to %s cannot be achieved with current tip", vmixvol.ToString()))
			/*
				tipchg, err := ChangeTips("", vmixvol, prms, ins.Prms, ins.Multi, true)

				if err != nil {
					return ret, wtype.LHError(wtype.LH_ERR_POLICY, fmt.Sprintf("Setting POST_MIX_VOLUME: %s", err.Error()))
				}
				ret = append(ret, tipchg...)
			*/
		}

		if ok {
			v := make([]wunit.Volume, ins.Multi)
			for i := 0; i < ins.Multi; i++ {
				vl := wunit.NewVolume(mixvol, "ul")
				v[i] = vl
			}
			mix.Volume = v
		}

		c := make([]int, ins.Multi)

		for i := 0; i < ins.Multi; i++ {
			c[i] = cycles
		}

		// set speed

		//mixrate, changespeed := pol["POST_MIX_RATE"]
		var changespeed bool
		mixrate := SafeGetF64(pol, "POST_MIX_RATE")
		if mixrate != pspeed && mixrate != 0.0 {
			changespeed = true
		}

		if changespeed {
			setspd := NewSetPipetteSpeedInstruction()
			setspd.Head = ins.Head
			setspd.Channel = -1 // all channels
			setspd.Speed = mixrate
			ret = append(ret, setspd)
		}

		mix.Cycles = c
		ret = append(ret, mix)

		if changespeed {
			sps := NewSetPipetteSpeedInstruction()
			sps.Head = ins.Head
			sps.Channel = -1 // all channels
			sps.Speed = pspeed
			ret = append(ret, sps)
		}

		// if we wait we need to do this here as well
		if wait_time > 0.0 {
			waitins := NewWaitInstruction()
			waitins.Time = wait_time
			ret = append(ret, waitins)
		}
	}

	// do we need to touch off?

	touch_off := SafeGetBool(pol, "TOUCHOFF")

	if touch_off {
		touch_offset := SafeGetF64(pol, "TOUCHOFFSET")
		mov := NewMoveInstruction()
		mov.Head = ins.Head
		mov.Pos = ins.PltTo
		mov.Plt = ins.TPlateType
		mov.Well = ins.WellTo
		mov.WVolume = ins.TVolume

		ref := make([]int, ins.Multi)
		off := make([]float64, ins.Multi)
		ox := make([]float64, ins.Multi)
		oy := make([]float64, ins.Multi)
		for i := 0; i < ins.Multi; i++ {
			ref[i] = 0
			off[i] = touch_offset
			ox[i] = 0.0
			oy[i] = 0.0
		}

		mov.Reference = ref
		mov.OffsetX = ox
		mov.OffsetY = oy
		mov.OffsetZ = off
		ret = append(ret, mov)
	}

	// now do we reset?

	// allow policies to override completely

	overridereset := SafeGetBool(pol, "RESET_OVERRIDE")

	if weneedtoreset && !overridereset {
		resetinstruction := NewResetInstruction()
		resetinstruction.AddMultiTransferParams(ins.Params())
		resetinstruction.Prms = ins.Prms
		ret = append(ret, resetinstruction)
	}

	return ret, nil
}

type SetPipetteSpeedInstruction struct {
	GenericRobotInstruction
	Type    int
	Head    int
	Channel int
	Speed   float64
}

func NewSetPipetteSpeedInstruction() *SetPipetteSpeedInstruction {
	var v SetPipetteSpeedInstruction
	v.Type = SPS
	return &v
}
func (ins *SetPipetteSpeedInstruction) InstructionType() int {
	return ins.Type
}

func (ins *SetPipetteSpeedInstruction) GetParameter(name string) interface{} {
	switch name {
	case "HEAD":
		return ins.Head
	case "CHANNEL":
		return ins.Channel
	case "SPEED":
		return ins.Speed
	case "INSTRUCTIONTYPE":
		return ins.InstructionType()
	}
	return nil
}

func (ins *SetPipetteSpeedInstruction) Generate(policy *wtype.LHPolicyRuleSet, prms *LHProperties) ([]RobotInstruction, error) {
	return nil, nil
}

func (ins *SetPipetteSpeedInstruction) OutputTo(driver LiquidhandlingDriver) error {
	ret := driver.SetPipetteSpeed(ins.Head, ins.Channel, ins.Speed)
	if !ret.OK {
		return fmt.Errorf(" %d : %s", ret.Errorcode, ret.Msg)
	}

	return nil

}

type SetDriveSpeedInstruction struct {
	GenericRobotInstruction
	Type  int
	Drive string
	Speed float64
}

func NewSetDriveSpeedInstruction() *SetDriveSpeedInstruction {
	var v SetDriveSpeedInstruction
	v.Type = SDS
	return &v
}
func (ins *SetDriveSpeedInstruction) InstructionType() int {
	return ins.Type
}

func (ins *SetDriveSpeedInstruction) GetParameter(name string) interface{} {
	switch name {
	case "DRIVE":
		return ins.Drive
	case "SPEED":
		return ins.Speed
	case "INSTRUCTIONTYPE":
		return ins.InstructionType()
	}
	return nil
}

func (ins *SetDriveSpeedInstruction) Generate(policy *wtype.LHPolicyRuleSet, prms *LHProperties) ([]RobotInstruction, error) {
	return nil, nil
}

func (ins *SetDriveSpeedInstruction) OutputTo(driver LiquidhandlingDriver) error {
	ret := driver.SetDriveSpeed(ins.Drive, ins.Speed)
	if !ret.OK {
		return fmt.Errorf(" %d : %s", ret.Errorcode, ret.Msg)
	}

	return nil

}

type InitializeInstruction struct {
	GenericRobotInstruction
	Type int
}

func NewInitializeInstruction() *InitializeInstruction {
	var v InitializeInstruction
	v.Type = INI
	return &v
}
func (ins *InitializeInstruction) InstructionType() int {
	return ins.Type
}

func (ins *InitializeInstruction) GetParameter(name string) interface{} {
	return nil
}

func (ins *InitializeInstruction) Generate(policy *wtype.LHPolicyRuleSet, prms *LHProperties) ([]RobotInstruction, error) {
	return nil, nil
}

func (ins *InitializeInstruction) OutputTo(driver LiquidhandlingDriver) error {
	ret := driver.Initialize()
	if !ret.OK {
		return fmt.Errorf(" %d : %s", ret.Errorcode, ret.Msg)
	}

	return nil

}

type FinalizeInstruction struct {
	GenericRobotInstruction
	Type int
}

func NewFinalizeInstruction() *FinalizeInstruction {
	var v FinalizeInstruction
	v.Type = FIN
	return &v
}
func (ins *FinalizeInstruction) InstructionType() int {
	return ins.Type
}

func (ins *FinalizeInstruction) GetParameter(name string) interface{} {
	return nil
}

func (ins *FinalizeInstruction) Generate(policy *wtype.LHPolicyRuleSet, prms *LHProperties) ([]RobotInstruction, error) {
	return nil, nil
}

func (ins *FinalizeInstruction) OutputTo(driver LiquidhandlingDriver) error {
	ret := driver.Finalize()
	if !ret.OK {
		return fmt.Errorf(" %d : %s", ret.Errorcode, ret.Msg)
	}

	return nil

}

type WaitInstruction struct {
	GenericRobotInstruction
	Type int
	Time float64
}

func NewWaitInstruction() *WaitInstruction {
	var v WaitInstruction
	v.Type = WAI
	return &v
}
func (ins *WaitInstruction) InstructionType() int {
	return ins.Type
}

func (ins *WaitInstruction) GetParameter(name string) interface{} {
	switch name {
	case "TIME":
		return ins.Time
	case "INSTRUCTIONTYPE":
		return ins.InstructionType()
	}
	return nil
}

func (ins *WaitInstruction) Generate(policy *wtype.LHPolicyRuleSet, prms *LHProperties) ([]RobotInstruction, error) {
	return nil, nil
}

func (ins *WaitInstruction) OutputTo(driver LiquidhandlingDriver) error {
	ret := driver.Wait(ins.Time)
	if !ret.OK {
		return fmt.Errorf(" %d : %s", ret.Errorcode, ret.Msg)
	}

	return nil

}

type LightsOnInstruction struct {
	GenericRobotInstruction
	Type       int
	What       []string
	PltFrom    []string
	PltTo      []string
	WellFrom   []string
	WellTo     []string
	Volume     []wunit.Volume
	FPlateType []string
	TPlateType []string
	FVolume    wunit.Volume
	TVolume    wunit.Volume
	Prms       *wtype.LHChannelParameter
}

func NewLightsOnInstruction() *LightsOnInstruction {
	var v LightsOnInstruction
	v.Type = LON
	v.What = make([]string, 0)
	v.PltFrom = make([]string, 0)
	v.PltTo = make([]string, 0)
	v.WellFrom = make([]string, 0)
	v.WellTo = make([]string, 0)
	v.Volume = make([]wunit.Volume, 0)
	return &v
}
func (ins *LightsOnInstruction) InstructionType() int {
	return ins.Type
}

func (ins *LightsOnInstruction) GetParameter(name string) interface{} {
	return nil
}

func (ins *LightsOnInstruction) Generate(policy *wtype.LHPolicyRuleSet, prms *LHProperties) ([]RobotInstruction, error) {
	return nil, nil
}

func (ins *LightsOnInstruction) OutputTo(driver LiquidhandlingDriver) error {
	return fmt.Errorf(" %d : %s", anthadriver.NIM, "Not yet implemented: LightsOnInstruction")
}

type LightsOffInstruction struct {
	GenericRobotInstruction
	Type       int
	What       []string
	PltFrom    []string
	PltTo      []string
	WellFrom   []string
	WellTo     []string
	Volume     []wunit.Volume
	FPlateType []string
	TPlateType []string
	FVolume    wunit.Volume
	TVolume    wunit.Volume
	Prms       *wtype.LHChannelParameter
}

func NewLightsOffInstruction() *LightsOffInstruction {
	var v LightsOffInstruction
	v.Type = LOF
	v.What = make([]string, 0)
	v.PltFrom = make([]string, 0)
	v.PltTo = make([]string, 0)
	v.WellFrom = make([]string, 0)
	v.WellTo = make([]string, 0)
	v.Volume = make([]wunit.Volume, 0)
	return &v
}
func (ins *LightsOffInstruction) InstructionType() int {
	return ins.Type
}

func (ins *LightsOffInstruction) GetParameter(name string) interface{} {
	return nil
}

func (ins *LightsOffInstruction) Generate(policy *wtype.LHPolicyRuleSet, prms *LHProperties) ([]RobotInstruction, error) {
	return nil, nil
}

func (ins *LightsOffInstruction) OutputTo(driver LiquidhandlingDriver) error {
	return fmt.Errorf(" %d : %s", anthadriver.NIM, "Not yet implemented: LightsOffInstruction")
}

type OpenInstruction struct {
	GenericRobotInstruction
	Type       int
	What       []string
	PltFrom    []string
	PltTo      []string
	WellFrom   []string
	WellTo     []string
	Volume     []wunit.Volume
	FPlateType []string
	TPlateType []string
	FVolume    wunit.Volume
	TVolume    wunit.Volume
	Prms       *wtype.LHChannelParameter
}

func NewOpenInstruction() *OpenInstruction {
	var v OpenInstruction
	v.Type = OPN
	v.What = make([]string, 0)
	v.PltFrom = make([]string, 0)
	v.PltTo = make([]string, 0)
	v.WellFrom = make([]string, 0)
	v.WellTo = make([]string, 0)
	v.Volume = make([]wunit.Volume, 0)
	return &v
}
func (ins *OpenInstruction) InstructionType() int {
	return ins.Type
}

func (ins *OpenInstruction) GetParameter(name string) interface{} {
	return nil
}

func (ins *OpenInstruction) Generate(policy *wtype.LHPolicyRuleSet, prms *LHProperties) ([]RobotInstruction, error) {
	return nil, nil
}

func (ins *OpenInstruction) OutputTo(driver LiquidhandlingDriver) error {
	return fmt.Errorf(" %d : %s", anthadriver.NIM, "Not yet implemented: OpenInstruction")
}

type CloseInstruction struct {
	GenericRobotInstruction
	Type       int
	What       []string
	PltFrom    []string
	PltTo      []string
	WellFrom   []string
	WellTo     []string
	Volume     []wunit.Volume
	FPlateType []string
	TPlateType []string
	FVolume    wunit.Volume
	TVolume    wunit.Volume
	Prms       *wtype.LHChannelParameter
}

func NewCloseInstruction() *CloseInstruction {
	var v CloseInstruction
	v.Type = CLS
	v.What = make([]string, 0)
	v.PltFrom = make([]string, 0)
	v.PltTo = make([]string, 0)
	v.WellFrom = make([]string, 0)
	v.WellTo = make([]string, 0)
	v.Volume = make([]wunit.Volume, 0)
	return &v
}
func (ins *CloseInstruction) InstructionType() int {
	return ins.Type
}

func (ins *CloseInstruction) GetParameter(name string) interface{} {
	return nil
}

func (ins *CloseInstruction) Generate(policy *wtype.LHPolicyRuleSet, prms *LHProperties) ([]RobotInstruction, error) {
	return nil, nil
}

func (ins *CloseInstruction) OutputTo(driver LiquidhandlingDriver) error {
	return fmt.Errorf(" %d : %s", anthadriver.NIM, "Not yet implemented: CloseInstruction")
}

type LoadAdaptorInstruction struct {
	GenericRobotInstruction
	Type       int
	What       []string
	PltFrom    []string
	PltTo      []string
	WellFrom   []string
	WellTo     []string
	Volume     []wunit.Volume
	FPlateType []string
	TPlateType []string
	FVolume    wunit.Volume
	TVolume    wunit.Volume
	Prms       *wtype.LHChannelParameter
}

func NewLoadAdaptorInstruction() *LoadAdaptorInstruction {
	var v LoadAdaptorInstruction
	v.Type = LAD
	v.What = make([]string, 0)
	v.PltFrom = make([]string, 0)
	v.PltTo = make([]string, 0)
	v.WellFrom = make([]string, 0)
	v.WellTo = make([]string, 0)
	v.Volume = make([]wunit.Volume, 0)
	return &v
}
func (ins *LoadAdaptorInstruction) InstructionType() int {
	return ins.Type
}

func (ins *LoadAdaptorInstruction) GetParameter(name string) interface{} {
	return nil
}

func (ins *LoadAdaptorInstruction) Generate(policy *wtype.LHPolicyRuleSet, prms *LHProperties) ([]RobotInstruction, error) {
	return nil, nil
}

func (ins *LoadAdaptorInstruction) OutputTo(driver LiquidhandlingDriver) error {
	return fmt.Errorf(" %d : %s", anthadriver.NIM, "Not yet implemented: LoadAdaptor")
}

type UnloadAdaptorInstruction struct {
	GenericRobotInstruction
	Type       int
	What       []string
	PltFrom    []string
	PltTo      []string
	WellFrom   []string
	WellTo     []string
	Volume     []wunit.Volume
	FPlateType []string
	TPlateType []string
	FVolume    wunit.Volume
	TVolume    wunit.Volume
	Prms       *wtype.LHChannelParameter
}

func NewUnloadAdaptorInstruction() *UnloadAdaptorInstruction {
	var v UnloadAdaptorInstruction
	v.Type = UAD
	v.What = make([]string, 0)
	v.PltFrom = make([]string, 0)
	v.PltTo = make([]string, 0)
	v.WellFrom = make([]string, 0)
	v.WellTo = make([]string, 0)
	v.Volume = make([]wunit.Volume, 0)
	return &v
}
func (ins *UnloadAdaptorInstruction) InstructionType() int {
	return ins.Type
}

func (ins *UnloadAdaptorInstruction) GetParameter(name string) interface{} {
	return nil
}

func (ins *UnloadAdaptorInstruction) Generate(policy *wtype.LHPolicyRuleSet, prms *LHProperties) ([]RobotInstruction, error) {
	return nil, nil
}

func (ins *UnloadAdaptorInstruction) OutputTo(driver LiquidhandlingDriver) error {
	return fmt.Errorf(" %d : %s", anthadriver.NIM, "Not yet implemented: UnloadAdaptor")
}

type ResetInstruction struct {
	GenericRobotInstruction
	Type       int
	What       []string
	PltFrom    []string
	PltTo      []string
	WellFrom   []string
	WellTo     []string
	Volume     []wunit.Volume
	FPlateType []string
	TPlateType []string
	FVolume    []wunit.Volume
	TVolume    []wunit.Volume
	Prms       *wtype.LHChannelParameter
}

func NewResetInstruction() *ResetInstruction {
	var ri ResetInstruction
	ri.Type = RST
	ri.What = make([]string, 0)
	ri.PltFrom = make([]string, 0)
	ri.WellFrom = make([]string, 0)
	ri.WellTo = make([]string, 0)
	ri.Volume = make([]wunit.Volume, 0)
	ri.FPlateType = make([]string, 0)
	ri.TPlateType = make([]string, 0)
	ri.FVolume = make([]wunit.Volume, 0)
	ri.TVolume = make([]wunit.Volume, 0)
	ri.GenericRobotInstruction.Ins = &ri
	return &ri
}

func (ins *ResetInstruction) InstructionType() int {
	return ins.Type
}

func (ins *ResetInstruction) GetParameter(name string) interface{} {
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
	case "PARAMS":
		return ins.Prms
	case "PLATFORM":
		if ins.Prms == nil {
			return ""
		}
		return ins.Prms.Platform
	case "INSTRUCTIONTYPE":
		return ins.InstructionType()
	}
	return nil
}

func (ins *ResetInstruction) AddTransferParams(tp TransferParams) {
	ins.What = append(ins.What, tp.What)
	ins.PltTo = append(ins.PltTo, tp.PltTo)
	ins.WellTo = append(ins.WellTo, tp.WellTo)
	ins.Volume = append(ins.Volume, tp.Volume)
	ins.TPlateType = append(ins.TPlateType, tp.TPlateType)
	ins.TVolume = append(ins.TVolume, tp.TVolume)
	ins.Prms = tp.Channel
}

func (ins *ResetInstruction) AddMultiTransferParams(mtp MultiTransferParams) {
	ins.What = mtp.What
	ins.PltTo = mtp.PltTo
	ins.WellTo = mtp.WellTo
	ins.Volume = mtp.Volume
	ins.TPlateType = mtp.TPlateType
	ins.TVolume = mtp.TVolume
	ins.Prms = mtp.Channel
}

func (ins *ResetInstruction) Generate(policy *wtype.LHPolicyRuleSet, prms *LHProperties) ([]RobotInstruction, error) {
	pol := GetPolicyFor(policy, ins)
	ret := make([]RobotInstruction, 0)

	mov := NewMoveInstruction()
	mov.Well = ins.WellTo
	mov.Pos = ins.PltTo
	mov.Plt = ins.TPlateType
	mov.WVolume = ins.TVolume
	mov.Head = ins.Prms.Head
	mov.Reference = append(mov.Reference, pol["BLOWOUTREFERENCE"].(int))
	mov.OffsetX = append(mov.OffsetX, 0.0)
	mov.OffsetY = append(mov.OffsetY, 0.0)
	mov.OffsetZ = append(mov.OffsetZ, pol["BLOWOUTOFFSET"].(float64))

	blow := NewBlowoutInstruction()

	blow.Head = ins.Prms.Head
	bov := wunit.NewVolume(pol["BLOWOUTVOLUME"].(float64), pol["BLOWOUTVOLUMEUNIT"].(string))
	blow.Volume = append(blow.Volume, bov)
	blow.Multi = len(ins.What)
	blow.Plt = ins.TPlateType
	blow.What = ins.What

	for i := 0; i < len(ins.What); i++ {
		blow.LLF = append(blow.LLF, false)
	}

	mov2 := NewMoveInstruction()
	mov2.Well = ins.WellTo
	mov2.Pos = ins.PltTo
	mov2.Plt = ins.TPlateType
	mov2.WVolume = ins.TVolume
	mov2.Head = ins.Prms.Head
	mov2.Reference = append(mov2.Reference, pol["PTZREFERENCE"].(int))
	mov2.OffsetX = append(mov2.OffsetX, 0.0)
	mov2.OffsetY = append(mov2.OffsetY, 0.0)
	mov2.OffsetZ = append(mov2.OffsetZ, pol["PTZOFFSET"].(float64))

	ptz := NewPTZInstruction()

	ptz.Head = ins.Prms.Head
	ptz.Channel = -1 // all channels

	ret = append(ret, mov)
	ret = append(ret, blow)

	// when needed we will add this pistons-to-zero instruction
	manptz := SafeGetBool(pol, "MANUALPTZ")
	if manptz {
		ret = append(ret, mov2)
		ret = append(ret, ptz)
	}
	return ret, nil
}

type MoveMixInstruction struct {
	GenericRobotInstruction
	Type      int
	Head      int
	Plt       []string
	Well      []string
	Volume    []wunit.Volume
	PlateType []string
	FVolume   []wunit.Volume
	Cycles    []int
	What      []string
	Blowout   []bool
	OffsetX   []float64
	OffsetY   []float64
	OffsetZ   []float64
	Multi     int
	Prms      map[string]interface{}
}

func NewMoveMixInstruction() *MoveMixInstruction {
	var mi MoveMixInstruction

	mi.Type = MMX
	mi.Plt = make([]string, 0)
	mi.Well = make([]string, 0)
	mi.Volume = make([]wunit.Volume, 0)
	mi.FVolume = make([]wunit.Volume, 0)
	mi.PlateType = make([]string, 0)
	mi.Cycles = make([]int, 0)
	mi.Prms = make(map[string]interface{})
	mi.What = make([]string, 0)
	mi.Blowout = make([]bool, 0)
	mi.OffsetX = make([]float64, 0)
	mi.OffsetY = make([]float64, 0)
	mi.OffsetZ = make([]float64, 0)
	return &mi
}

func (ins *MoveMixInstruction) GetParameter(name string) interface{} {
	switch name {
	case "VOLUME":
		return ins.Volume
	case "VOLUNT":
		return nil
	case "PLATETYPE":
		return ins.PlateType
	case "WELLVOLUME":
		return ins.FVolume
	case "POS":
		return ins.Plt
	case "WELL":
		return ins.Well
	case "PARAMS":
		return ins.Prms
	case "CYCLES":
		return ins.Cycles
	case "INSTRUCTIONTYPE":
		return ins.InstructionType()
	case "WHAT":
		return ins.What
	case "BLOWOUT":
		return ins.Blowout
	case "OFFSETX":
		return ins.OffsetX
	case "OFFSETY":
		return ins.OffsetY
	case "OFFSETZ":
		return ins.OffsetZ
	}

	return nil

}

func (ins *MoveMixInstruction) InstructionType() int {
	return MMX
}

func (ins *MoveMixInstruction) Generate(policy *wtype.LHPolicyRuleSet, prms *LHProperties) ([]RobotInstruction, error) {
	ret := make([]RobotInstruction, 2)

	// move

	mov := NewMoveInstruction()
	mov.Well = ins.Well
	mov.Pos = ins.Plt
	mov.Plt = ins.PlateType
	mov.WVolume = ins.FVolume
	mov.Head = ins.Head
	mov.OffsetX = ins.OffsetX
	mov.OffsetY = ins.OffsetY
	mov.OffsetZ = ins.OffsetZ
	ref := make([]int, ins.Multi)
	ref[0] = 0
	mov.Reference = ref
	ret[0] = mov

	// mix

	mix := NewMixInstruction()
	mix.Head = ins.Head
	mix.PlateType = ins.PlateType
	mix.Cycles = ins.Cycles
	mix.Volume = ins.Volume
	mix.Multi = ins.Multi
	mix.What = ins.What
	mix.Blowout = ins.Blowout
	ret[1] = mix

	return ret, nil
}

type MixInstruction struct {
	GenericRobotInstruction
	Type      int
	Head      int
	Volume    []wunit.Volume
	PlateType []string
	What      []string
	Blowout   []bool
	Multi     int
	Cycles    []int
}

func NewMixInstruction() *MixInstruction {
	var mi MixInstruction

	mi.Type = MMX
	mi.Volume = make([]wunit.Volume, 0)
	mi.PlateType = make([]string, 0)
	mi.Cycles = make([]int, 0)
	mi.What = make([]string, 0)
	mi.Blowout = make([]bool, 0)
	return &mi
}

func (mi *MixInstruction) InstructionType() int {
	return mi.Type
}

func (ins *MixInstruction) Generate(policy *wtype.LHPolicyRuleSet, prms *LHProperties) ([]RobotInstruction, error) {
	return nil, nil
}

func (ins *MixInstruction) GetParameter(name string) interface{} {
	switch name {
	case "VOLUME":
		return ins.Volume
	case "VOLUNT":
		return nil
	case "PLATETYPE":
		return ins.PlateType
	case "CYCLES":
		return ins.Cycles
	case "INSTRUCTIONTYPE":
		return ins.InstructionType()
	case "LIQUIDCLASS":
		return ins.What
	}
	return nil

}

func (mi *MixInstruction) OutputTo(driver LiquidhandlingDriver) error {
	vols := make([]float64, len(mi.Volume))

	for i := 0; i < len(mi.Volume); i++ {
		vols[i] = mi.Volume[i].ConvertTo(wunit.ParsePrefixedUnit("ul"))
	}

	ret := driver.Mix(mi.Head, vols, mi.PlateType, mi.Cycles, mi.Multi, mi.What, mi.Blowout)

	if !ret.OK {
		return fmt.Errorf(" %d : %s", ret.Errorcode, ret.Msg)
	}

	return nil

}

// TODO -- implement MESSAGE

func ChangeTips(tiptype string, vol wunit.Volume, prms *LHProperties, channel *wtype.LHChannelParameter, multi int, oneshot bool) ([]RobotInstruction, error) {
	ret := make([]RobotInstruction, 0, 2)
	newchannel, newtiptype := ChooseChannel(vol, prms)

	if !newchannel.CanMove(vol, oneshot) {
		return ret, fmt.Errorf("No channel can move a volume of %s in one shot", vol.ToString())
	}

	tipdrp, err := DropTips(tiptype, prms, channel, multi)

	if err != nil {
		return ret, err
	}
	ret = append(ret, tipdrp)

	tipget, err := GetTips(newtiptype, prms, newchannel, multi, false)
	if err != nil {
		return ret, err
	}
	ret = append(ret, tipget)
	return ret, err
}

func GetTips(tiptype string, params *LHProperties, channel *wtype.LHChannelParameter, multi int, mirror bool) (RobotInstruction, error) {

	tipwells, tipboxpositions, tipboxtypes, terr := params.GetCleanTips(tiptype, channel, mirror, multi)

	if tipwells == nil || terr != nil {
		/*
			logger.Fatal("No tips left")
			panic("NO TIPS LEFT BOYO")
		*/
		err := wtype.LHError(wtype.LH_ERR_NO_TIPS, fmt.Sprint("PICKUP: type: ", tiptype, " n: ", multi, " mirror: ", mirror))
		return NewLoadTipsMoveInstruction(), err
	}

	ins := NewLoadTipsMoveInstruction()
	ins.Head = channel.Head
	ins.Well = tipwells
	ins.FPosition = tipboxpositions
	ins.FPlateType = tipboxtypes
	ins.Multi = multi

	return ins, nil
}

func DropTips(tiptype string, params *LHProperties, channel *wtype.LHChannelParameter, multi int) (RobotInstruction, error) {
	tipwells, tipwastepositions, tipwastetypes := params.DropDirtyTips(channel, multi)

	if tipwells == nil {
		//logger.Fatal("Could not dispose tip. No usable tipwell found")
		//panic("NO ROOM AT THE INN FOR THESE LITTLE TIPS")
		ins := NewUnloadTipsMoveInstruction()
		err := wtype.LHError(wtype.LH_ERR_NO_TIPS, fmt.Sprint("DROP: type: ", tiptype, " n: ", multi))
		return ins, err
	}

	ins := NewUnloadTipsMoveInstruction()
	ins.Head = channel.Head
	ins.WellTo = tipwells
	ins.PltTo = tipwastepositions
	ins.TPlateType = tipwastetypes
	ins.Multi = multi
	return ins, nil
}
