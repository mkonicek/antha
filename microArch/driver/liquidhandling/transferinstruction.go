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
	"context"
	"fmt"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/antha/anthalib/wutil"
	"github.com/antha-lang/antha/inventory/cache"
	"sort"
	"strconv"
	"strings"
)

type TransferInstruction struct {
	BaseRobotInstruction
	*InstructionType
	Platform  string
	Transfers []MultiTransferParams
}

func (ti *TransferInstruction) ToString() string {
	s := ti.Type().Name
	for i := 0; i < len(ti.Transfers); i++ {
		s += ti.ParamSet(i).ToString()
		s += "\n"
	}

	return s
}

func (ti *TransferInstruction) ParamSet(n int) MultiTransferParams {
	return ti.Transfers[n]
}

func NewTransferInstruction(what, pltfrom, pltto, wellfrom, wellto, fplatetype, tplatetype []string, volume, fvolume, tvolume []wunit.Volume, FPlateWX, FPlateWY, TPlateWX, TPlateWY []int, Components []string, policies []wtype.LHPolicy) *TransferInstruction {
	tfri := &TransferInstruction{
		InstructionType: TFR,
		Transfers:       make([]MultiTransferParams, 0, 1),
	}

	/*
		v := MultiTransferParams{
			What:       what,
			PltFrom:    pltfrom,
			PltTo:      pltto,
			WellFrom:   wellfrom,
			WellTo:     wellto,
			Volume:     volume,
			FPlateType: fplatetype,
			TPlateType: tplatetype,
			FVolume:    fvolume,
			TVolume:    tvolume,
			FPlateWX:   FPlateWX,
			FPlateWY:   FPlateWY,
			TPlateWX:   TPlateWX,
			TPlateWY:   TPlateWY,
			Components: Components,
		}
	*/

	v := MTPFromArrays(what, pltfrom, pltto, wellfrom, wellto, fplatetype, tplatetype, volume, fvolume, tvolume, FPlateWX, FPlateWY, TPlateWX, TPlateWY, Components, policies)

	tfri.Add(v)
	tfri.BaseRobotInstruction = NewBaseRobotInstruction(tfri)
	return tfri
}

func (ins *TransferInstruction) Visit(visitor RobotInstructionVisitor) {
	visitor.Transfer(ins)
}

func (ins *TransferInstruction) OutputTo(drv LiquidhandlingDriver) error {
	hlld, ok := drv.(HighLevelLiquidhandlingDriver)

	if !ok {
		return fmt.Errorf("Driver type %T not compatible with TransferInstruction, need HighLevelLiquidhandlingDriver", drv)
	}

	// make sure we disable the RobotInstruction pointer
	ins.BaseRobotInstruction = BaseRobotInstruction{}

	volumes := make([]float64, len(SetOfMultiTransferParams(ins.Transfers).Volume()))
	for i, vol := range SetOfMultiTransferParams(ins.Transfers).Volume() {
		volumes[i] = vol.ConvertToString("ul")
	}

	reply := hlld.Transfer(SetOfMultiTransferParams(ins.Transfers).What(), SetOfMultiTransferParams(ins.Transfers).PltFrom(), SetOfMultiTransferParams(ins.Transfers).WellFrom(), SetOfMultiTransferParams(ins.Transfers).PltTo(), SetOfMultiTransferParams(ins.Transfers).WellTo(), volumes)

	if !reply.OK {
		return fmt.Errorf(" %d : %s", reply.Errorcode, reply.Msg)
	}

	return nil
}

func (tfri *TransferInstruction) Add(tp MultiTransferParams) {
	tfri.Transfers = append(tfri.Transfers, tp)
}

//what, pltfrom, pltto, wellfrom, wellto, fplatetype, tplatetype []string, volume, fvolume, tvolume []wunit.Volume, FPlateWX, FPlateWY, TPlateWX, TPlateWY []int, Components []string
func (ins *TransferInstruction) Dup() *TransferInstruction {
	tfri := &TransferInstruction{
		InstructionType: TFR,
		Transfers:       make([]MultiTransferParams, 0, 1),
		Platform:        ins.Platform,
	}
	tfri.BaseRobotInstruction = NewBaseRobotInstruction(tfri)

	for _, tfr := range ins.Transfers {
		tfri.Add(tfr.Dup())
	}

	return tfri
}

func (ins *TransferInstruction) MergeWith(ins2 *TransferInstruction) *TransferInstruction {
	ret := ins.Dup()

	for _, v := range ins2.Transfers {
		ins.Add(v)
	}

	return ret
}

func (ins *TransferInstruction) GetParameter(name InstructionParameter) interface{} {
	switch name {
	case LIQUIDCLASS:
		return SetOfMultiTransferParams(ins.Transfers).What()
	case VOLUME:
		return SetOfMultiTransferParams(ins.Transfers).Volume()
	case FROMPLATETYPE:
		return SetOfMultiTransferParams(ins.Transfers).FPlateType()
	case WELLFROMVOLUME:
		return SetOfMultiTransferParams(ins.Transfers).FVolume()
	case POSFROM:
		return SetOfMultiTransferParams(ins.Transfers).PltFrom()
	case POSTO:
		return SetOfMultiTransferParams(ins.Transfers).PltTo()
	case WELLFROM:
		return SetOfMultiTransferParams(ins.Transfers).WellFrom()
	case WELLTO:
		return SetOfMultiTransferParams(ins.Transfers).WellTo()
	case WELLTOVOLUME:
		return SetOfMultiTransferParams(ins.Transfers).TVolume()
	case TOPLATETYPE:
		return SetOfMultiTransferParams(ins.Transfers).TPlateType()
	case FPLATEWX:
		return SetOfMultiTransferParams(ins.Transfers).FPlateWX()
	case FPLATEWY:
		return SetOfMultiTransferParams(ins.Transfers).FPlateWY()
	case TPLATEWX:
		return SetOfMultiTransferParams(ins.Transfers).TPlateWX()
	case TPLATEWY:
		return SetOfMultiTransferParams(ins.Transfers).TPlateWY()
	case PLATFORM:
		return ins.Platform
	case COMPONENT:
		return SetOfMultiTransferParams(ins.Transfers).Component()
	default:
		return ins.BaseRobotInstruction.GetParameter(name)
	}
}

/*
func (ins *TransferInstruction) getPoliciesForTransfer(which int, ruleSet wtype.LHPolicyRuleSet) []wtype.LHPolicy {
}
*/
func (ins *TransferInstruction) CheckMultiPolicies(which int) bool {
	// first iteration: ensure all the WHAT prms are the same
	// later	  : actually check the policies per channel

	nwhat := wutil.NUniqueStringsInArray(ins.Transfers[which].What(), true)

	return nwhat == 1
}

func firstNonEmpty(types []string) string {
	for _, typ := range types {
		if typ == "" {
			continue
		}
		return typ
	}
	return ""
}

// add policies as argument to GetParallelSetsFor to check multichannelability
func (ins *TransferInstruction) GetParallelSetsFor(ctx context.Context, robot *LHProperties, policy wtype.LHPolicy) []int {
	r := make([]int, 0, len(ins.Transfers))

	for i := 0; i < len(ins.Transfers); i++ {
		// a parallel transfer is valid if any robot head can do it
		// TODO --> support head/adaptor changes. Maybe.
		for _, head := range robot.GetLoadedHeads() {
			if ins.validateParallelSet(ctx, robot, head, i, policy) {
				r = append(r, i)
			}
		}
	}

	return r
}

// add policies as argument to GetParallelSetsFor to check multichannelability
// which is the index relating to position in multitransferparams matrix
func (ins *TransferInstruction) validateParallelSet(ctx context.Context, robot *LHProperties, head *wtype.LHHead, which int, policy wtype.LHPolicy) bool {
	channel := head.Adaptor.Params

	if channel.Multi == 1 {
		return false
	}

	if len(ins.Transfers[which].What()) > channel.Multi {
		return false
	}

	npositions := wutil.NUniqueStringsInArray(ins.Transfers[which].PltFrom(), true)

	if npositions != 1 {
		// fall back to single-channel
		// TODO -- find a subset we CAN do, if such exists
		return false
	}

	nplatetypes := wutil.NUniqueStringsInArray(ins.Transfers[which].FPlateType(), true)

	if nplatetypes != 1 {
		// fall back to single-channel
		// TODO -- find a subset we CAN do , if such exists
		return false
	}

	fromPlateType := firstNonEmpty(ins.Transfers[which].FPlateType())
	fromPlate, err := cache.NewPlate(ctx, fromPlateType)
	if err != nil {
		panic(err)
	}
	if fromPlate == nil {
		panic("No from plates in instruction")
	}

	// check source / tip alignment
	if !head.CanReach(fromPlate, wtype.WCArrayFromStrings(ins.Transfers[which].WellFrom())) {
		// fall back to single-channel
		// TODO -- find a subset we CAN do
		return false
	}

	err = cache.ReturnObject(ctx, fromPlate)
	if err != nil {
		panic(err)
	}

	toPlateType := firstNonEmpty(ins.Transfers[which].TPlateType())
	toPlate, err := cache.NewPlate(ctx, toPlateType)
	if err != nil {
		panic(err)
	}
	if toPlate == nil {
		panic("No to plates in instruction")
	}

	// for safety, check dest / tip alignment
	if !head.CanReach(toPlate, wtype.WCArrayFromStrings(ins.Transfers[which].WellTo())) {
		// fall back to single-channel
		// TODO -- find a subset we CAN do
		return false
	}

	err = cache.ReturnObject(ctx, toPlate)
	if err != nil {
		panic(err)
	}

	// check that we will not require different policies

	return ins.CheckMultiPolicies(which)
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

func (ins *TransferInstruction) ChooseChannels(prms *LHProperties) {
	for i, mtp := range ins.Transfers {
		// we need to remove leading blanks
		ins.Transfers[i] = mtp.RemoveInitialBlanks()
	}
}

func (ins *TransferInstruction) Generate(ctx context.Context, policy *wtype.LHPolicyRuleSet, prms *LHProperties) ([]RobotInstruction, error) {
	// if the liquid handler is of the high-level type we cut the tree here
	// after ensuring that the transfers are within limitations of the liquid handler

	if prms.GetLHType() == HLLiquidHandler {
		err := ins.ReviseTransferVolumes(prms)

		if err != nil {
			return []RobotInstruction{}, err
		}
		return []RobotInstruction{}, nil
	}

	//  set the channel choices first by cleaning out initial empties

	ins.ChooseChannels(prms)

	// is this the part we need to change?
	pol, err := GetPolicyFor(policy, ins)

	if err != nil {
		if _, ok := err.(ErrInvalidLiquidType); ok {
			return []RobotInstruction{}, err
		}
		pol, err = GetDefaultPolicy(policy, ins)

		if err != nil {
			return []RobotInstruction{}, err
		}
	}

	ret := make([]RobotInstruction, 0)

	headsLoaded := prms.GetLoadedHeads()

	// if we can multi we do this first
	if pol["CAN_MULTI"].(bool) {
		// add policies as argument to GetParallelSetsFor to check multichannelability
		parallelsets := ins.GetParallelSetsFor(ctx, prms, pol)

		mci := NewMultiChannelBlockInstruction()
		mci.Prms = headsLoaded[0].Params // TODO Remove Hard code here

		// to do below
		//
		//	- divide up transfer into multi and single transfers
		//  	  in practice this means finding the maximum we can do
		//	  then doing that as a transfer and generating single channel transfers
		//	  to mop up the rest
		//

		for _, set := range parallelsets {
			vols := VolumeSet(ins.Transfers[set].Volume())

			// non independent heads must have all volumes the same
			if !mci.Prms.Independent {
				if maxvol := vols.MaxMultiTransferVolume(prms.MinPossibleVolume()); !maxvol.IsZero() {
					for i := range vols {
						vols[i] = wunit.CopyVolume(maxvol)
					}
				} else {
					// we can't transfer any volumes with the non-independent head so move on to the next parallelset
					continue
				}
			}

			tp := ins.Transfers[set].Dup()
			for i := 0; i < len(tp.Transfers); i++ {
				tp.Transfers[i].Volume = vols[i].Dup()
			}

			ins.Transfers[set].RemoveVolumes(vols)

			// set the from and to volumes for the relevant part of the instruction
			ins.Transfers[set].RemoveFVolumes(vols)
			ins.Transfers[set].AddTVolumes(vols)

			mci.Multi = len(vols)
			mci.AddTransferParams(tp)
		}

		if len(parallelsets) > 0 && len(mci.Volume) > 0 {
			ret = append(ret, mci)
		}
	}

	// mop up all the single instructions which are left
	mci := NewMultiChannelBlockInstruction()
	mci.Prms = headsLoaded[0].Params // TODO Fix Hard Code Here

	lastCmp := ""
	for _, t := range ins.Transfers {
		for _, tp := range t.Transfers {
			if !tp.Volume.IsPositive() {
				continue
			}

			if lastCmp != "" && tp.Component != lastCmp {
				if len(mci.Volume) > 0 {
					ret = append(ret, mci)
				}
				mci = NewMultiChannelBlockInstruction()
				mci.Prms = headsLoaded[0].Params
			}

			mci.AddTransferParams(MultiTransferParams{Multi: 1, Transfers: []TransferParams{tp}})
			lastCmp = tp.Component
		}
	}
	if len(mci.Volume) > 0 {
		ret = append(ret, mci)
	}

	return ret, nil
}

func (ins *TransferInstruction) ReviseTransferVolumes(prms *LHProperties) error {
	newTransfers := make([]MultiTransferParams, 0, len(ins.Transfers))

	for _, mtp := range ins.Transfers {
		//newMtp := make(MultiTransferParams, len(mtp))
		newMtp := NewMultiTransferParams(mtp.Multi)
		for _, tp := range mtp.Transfers {
			if tp.What == "" {
				continue
			}
			newTPs, err := safeTransfers(tp, prms)
			if err != nil {
				return err
			}
			newMtp.Transfers = append(newMtp.Transfers, newTPs...)
		}

		newMtp.Multi = len(newMtp.Transfers)

		newTransfers = append(newTransfers, newMtp)
	}

	ins.Transfers = newTransfers

	return nil
}

func safeTransfers(tp TransferParams, prms *LHProperties) ([]TransferParams, error) {

	if tp.What == "" {
		return []TransferParams{tp}, nil
	}

	headsLoaded := prms.GetLoadedHeads()

	tvs, err := TransferVolumes(tp.Volume, headsLoaded[0].Params.Minvol, headsLoaded[0].Params.Maxvol)

	ret := []TransferParams{}

	if err != nil {
		return ret, err
	}

	fwv := tp.FVolume.Dup()
	twv := tp.TVolume.Dup()

	for _, v := range tvs {
		ntp := tp.Dup()
		ntp.Volume = v
		ntp.FVolume = fwv.Dup()
		ntp.TVolume = twv.Dup()
		fwv.Subtract(v)
		twv.Add(v)

		ret = append(ret, ntp)
	}

	return ret, nil
}

func MockAspDsp(ins RobotInstruction) []TerminalRobotInstruction {
	ret := make([]TerminalRobotInstruction, 0, 1)

	tfr, ok := ins.(*TransferInstruction)

	if !ok {
		return ret
	}

	ret = append(ret, mockLoad())

	for _, mtp := range tfr.Transfers {
		for _, tp := range mtp.Transfers {
			mox := mockLowLevels(tp)
			ret = append(ret, mox...)
		}
	}

	ret = append(ret, mockUnload())

	return ret
}

func mockLoad() *LoadTipsInstruction {
	ins := NewLoadTipsInstruction()
	ins.Pos = append(ins.Pos, "position_n")
	ins.Well = append(ins.Well, "A1")
	ins.Channels = append(ins.Channels, 0)
	ins.TipType = append(ins.TipType, "none")
	ins.HolderType = append(ins.HolderType, "none")
	ins.Multi = 1
	ins.Platform = "Echo"
	return ins
}

func mockLowLevels(tp TransferParams) []TerminalRobotInstruction {
	mova := NewMoveInstruction()
	mova.Plt = append(mova.Plt, tp.PltFrom)
	mova.Well = append(mova.Well, tp.WellFrom)
	mova.Reference = append(mova.Reference, 0)
	mova.WVolume = append(mova.WVolume, tp.FVolume)
	mova.Platform = "Echo"

	asp := NewAspirateInstruction()
	asp.Multi = 1
	asp.Volume = append(asp.Volume, tp.Volume)
	asp.Platform = "Echo"
	asp.What = append(asp.What, tp.What)

	movd := NewMoveInstruction()
	movd.Plt = append(movd.Plt, tp.PltTo)
	movd.Well = append(movd.Well, tp.WellTo)
	movd.Reference = append(movd.Reference, 0)
	movd.WVolume = append(movd.WVolume, tp.TVolume)
	movd.Platform = "Echo"

	dsp := NewDispenseInstruction()
	dsp.Multi = 1
	dsp.Volume = append(dsp.Volume, tp.Volume)
	dsp.Platform = "Echo"
	dsp.What = append(dsp.What, tp.What)

	return []TerminalRobotInstruction{mova, asp, movd, dsp}
}

func mockUnload() *UnloadTipsInstruction {
	ins := NewUnloadTipsInstruction()
	ins.Pos = append(ins.Pos, "position_n")
	ins.Well = append(ins.Well, "A1")
	ins.Channels = append(ins.Channels, 0)
	ins.TipType = append(ins.TipType, "none")
	ins.HolderType = append(ins.HolderType, "none")
	ins.Multi = 1
	ins.Platform = "Echo"
	return ins
}
