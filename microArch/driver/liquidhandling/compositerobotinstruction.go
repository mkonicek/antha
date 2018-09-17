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
	"context"
	"fmt"

	"github.com/pkg/errors"

	"reflect"

	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/antha/anthalib/wutil/text"
	"github.com/antha-lang/antha/inventory"
	anthadriver "github.com/antha-lang/antha/microArch/driver"
	"github.com/antha-lang/antha/microArch/logger"
)

// Valid parameter fields for robot instructions
const (
	//maxTouchOffset maximum value for TOUCHOFFSET which makes sense - values larger than this are capped to this value
	maxTouchOffset = 5.0

	//added to avoid floating point issues with heights in simulator
	safetyZHeight = 0.05
)

type SingleChannelBlockInstruction struct {
	BaseRobotInstruction
	*InstructionType
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
	v := &SingleChannelBlockInstruction{
		InstructionType: SCB,
		What:            []string{},
		PltFrom:         []string{},
		PltTo:           []string{},
		WellFrom:        []string{},
		WellTo:          []string{},
		Volume:          []wunit.Volume{},
		FVolume:         []wunit.Volume{},
		TVolume:         []wunit.Volume{},
		FPlateType:      []string{},
		TPlateType:      []string{},
	}
	v.BaseRobotInstruction = NewBaseRobotInstruction(v)
	return v
}

func (ins *SingleChannelBlockInstruction) Visit(visitor RobotInstructionVisitor) {
	visitor.SingleChannelBlock(ins)
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

func (ins *SingleChannelBlockInstruction) GetParameter(name InstructionParameter) interface{} {
	switch name {
	case LIQUIDCLASS:
		return ins.What
	case VOLUME:
		return ins.Volume
	case VOLUNT:
		return nil
	case FROMPLATETYPE:
		return ins.FPlateType
	case WELLFROMVOLUME:
		return ins.FVolume
	case POSFROM:
		return ins.PltFrom
	case POSTO:
		return ins.PltTo
	case WELLFROM:
		return ins.WellFrom
	case PARAMS:
		return ins.Prms
	case PLATFORM:
		if ins.Prms == nil {
			return ""
		}
		return ins.Prms.Platform
	case WELLTO:
		return ins.WellTo
	case WELLTOVOLUME:
		return ins.TVolume
	case TOPLATETYPE:
		return ins.TPlateType
	default:
		return ins.BaseRobotInstruction.GetParameter(name)
	}
}

func tipArrays(multi int) ([]string, []*wtype.LHChannelParameter) {
	// TODO --> mirroring
	tt := make([]string, multi)
	chanA := make([]*wtype.LHChannelParameter, multi)

	return tt, chanA
}

func (ins *SingleChannelBlockInstruction) Generate(ctx context.Context, policy *wtype.LHPolicyRuleSet, prms *LHProperties) ([]RobotInstruction, error) {
	usetiptracking := SafeGetBool(policy.Options, "USE_DRIVER_TIP_TRACKING")

	ret := make([]RobotInstruction, 0)
	// get tips
	channel, tipp, err := ChooseChannel(ins.Volume[0], prms)
	if err != nil {
		return ret, err
	}

	tiptype := tipp.Type

	ins.Prms = channel
	pol, err := GetPolicyFor(policy, ins)

	if err != nil {
		if _, ok := err.(ErrInvalidLiquidType); ok {
			return ret, err
		}
		pol, err = GetDefaultPolicy(policy, ins)

		if err != nil {
			return ret, err
		}
	}

	tt, chanA := tipArrays(channel.Multi)
	tt[0] = tiptype
	chanA[0] = channel

	tipget, err := GetTips(ctx, tt, prms, chanA, usetiptracking)

	if err != nil {
		return ret, err
	}

	ret = append(ret, tipget...)
	n_tip_uses := 0

	var last_thing *wtype.Liquid
	var dirty bool

	for t := 0; t < len(ins.Volume); t++ {
		newchannel, newtipp, err := ChooseChannel(ins.Volume[t], prms)
		if err != nil {
			return ret, err
		}

		newtiptype := newtipp.Type
		mergedchannel := newchannel.MergeWithTip(newtipp)
		tipp = newtipp

		tvs, err := TransferVolumes(ins.Volume[t], mergedchannel.Minvol, mergedchannel.Maxvol)

		if err != nil {
			return ret, err
		}
		for _, vol := range tvs {
			// determine whether to change tips
			change_tips := n_tip_uses > pol["TIP_REUSE_LIMIT"].(int)
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
				tipdrp, err := DropTips(tt, prms, chanA)
				if err != nil {
					return ret, err
				}
				ret = append(ret, tipdrp)

				tt, chanA = tipArrays(newchannel.Multi)
				tt[0] = newtiptype
				chanA[0] = newchannel
				tipget, err := GetTips(ctx, tt, prms, chanA, usetiptracking)

				if err != nil {
					return ret, err
				}

				ret = append(ret, tipget...)
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
			stci.Prms = channel.MergeWithTip(tipp)
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
			}

			ins.FVolume[t].Subtract(vol)
			ins.TVolume[t].Add(vol)
			n_tip_uses += 1
		}

	}
	tipdrp, err := DropTips(tt, prms, chanA)

	if err != nil {
		return ret, err
	}
	ret = append(ret, tipdrp)

	return ret, nil
}

type MultiChannelBlockInstruction struct {
	BaseRobotInstruction
	*InstructionType
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
	v := &MultiChannelBlockInstruction{
		InstructionType: MCB,
		What:            [][]string{},
		PltFrom:         [][]string{},
		PltTo:           [][]string{},
		WellFrom:        [][]string{},
		WellTo:          [][]string{},
		Volume:          [][]wunit.Volume{},
		FPlateType:      [][]string{},
		TPlateType:      [][]string{},
		FVolume:         [][]wunit.Volume{},
		TVolume:         [][]wunit.Volume{},
	}
	v.BaseRobotInstruction = NewBaseRobotInstruction(v)
	return v
}

func (ins *MultiChannelBlockInstruction) Visit(visitor RobotInstructionVisitor) {
	visitor.MultiChannelBlock(ins)
}

func (ins *MultiChannelBlockInstruction) AddTransferParams(mct MultiTransferParams) {
	ins.What = append(ins.What, mct.What())
	ins.PltFrom = append(ins.PltFrom, mct.PltFrom())
	ins.PltTo = append(ins.PltTo, mct.PltTo())
	ins.WellFrom = append(ins.WellFrom, mct.WellFrom())
	ins.WellTo = append(ins.WellTo, mct.WellTo())
	ins.Volume = append(ins.Volume, mct.Volume())
	ins.FPlateType = append(ins.FPlateType, mct.FPlateType())
	ins.TPlateType = append(ins.TPlateType, mct.TPlateType())
	ins.FVolume = append(ins.FVolume, mct.FVolume())
	ins.TVolume = append(ins.TVolume, mct.TVolume())
}

func (ins *MultiChannelBlockInstruction) GetParameter(name InstructionParameter) interface{} {
	switch name {
	case LIQUIDCLASS:
		return ins.What
	case VOLUME:
		return ins.Volume
	case VOLUNT:
		return nil
	case FROMPLATETYPE:
		return ins.FPlateType
	case WELLFROMVOLUME:
		return ins.FVolume
	case POSFROM:
		return ins.PltFrom
	case POSTO:
		return ins.PltTo
	case WELLFROM:
		return ins.WellFrom
	case PARAMS:
		return ins.Prms
	case PLATFORM:
		if ins.Prms == nil {
			return ""
		}
		return ins.Prms.Platform
	case WELLTO:
		return ins.WellTo
	case WELLTOVOLUME:
		return ins.TVolume
	case TOPLATETYPE:
		return ins.TPlateType
	default:
		return ins.BaseRobotInstruction.GetParameter(name)
	}
}

func (ins *MultiChannelBlockInstruction) GetVolumes() []wunit.Volume {
	v := make([]wunit.Volume, 0, 1)
	seen := make(map[string]bool)
	for _, vv := range ins.Volume[0] {
		if !vv.IsZero() && !seen[vv.ToString()] {
			seen[vv.ToString()] = true
			v = append(v, vv)
		}
	}

	return v
}

func mergeTipsAndChannels(channels []*wtype.LHChannelParameter, tips []*wtype.LHTip) []*wtype.LHChannelParameter {
	ret := make([]*wtype.LHChannelParameter, len(channels))

	for i := 0; i < len(channels); i++ {
		if channels[i] != nil {
			if tips[i] != nil {
				ret[i] = channels[i].MergeWithTip(tips[i])
			} else {
				ret[i] = channels[i].Dup()
			}
		}
	}

	return ret
}

// By the point at which the MultiChannelBlockInstruction is used by the Generate method all transfers will share the same policy.
func (ins *MultiChannelBlockInstruction) Generate(ctx context.Context, policy *wtype.LHPolicyRuleSet, prms *LHProperties) ([]RobotInstruction, error) {
	usetiptracking := SafeGetBool(policy.Options, "USE_DRIVER_TIP_TRACKING")

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
	// get some tips

	// we no longer require ins.volume[0][0] to be set
	// as we move to independent we need to get all volumes

	//channels, _, tiptypes, err := ChooseChannels(ins.GetVolumes(), prms)
	channels, _, tiptypes, err := ChooseChannels(ins.Volume[0], prms)
	if err != nil {
		return ret, err
	}

	tipget, err := GetTips(ctx, tiptypes, prms, channels, usetiptracking)
	if err != nil {
		return ret, err
	}
	ret = append(ret, tipget...)
	n_tip_uses := 0
	var last_thing *wtype.Liquid
	var dirty bool

	for t := 0; t < len(ins.Volume); t++ {
		tvols := NewVolumeSet(ins.Prms.Multi)
		//		vols := NewVolumeSet(ins.Prms.Multi)
		fvols := NewVolumeSet(ins.Prms.Multi)
		for i := range ins.Volume[t] {
			fvols[i] = wunit.CopyVolume(ins.FVolume[t][i])
			tvols[i] = wunit.CopyVolume(ins.TVolume[t][i])
		}

		// choose tips
		newchannels, newtips, newtiptypes, err := ChooseChannels(ins.Volume[t], prms)
		if err != nil {
			return ret, err
		}

		// load tips

		// split the transfer up
		// volumes no longer equal
		tvs, err := TransferVolumesMulti(VolumeSet(ins.Volume[t]), mergeTipsAndChannels(newchannels, newtips))

		if err != nil {
			return ret, err
		}

		for _, vols := range tvs {
			// determine whether to change tips
			// INMC: DO THIS PER CHANNEL
			change_tips := n_tip_uses > pol["TIP_REUSE_LIMIT"].(int)
			change_tips = change_tips || !reflect.DeepEqual(channels, newchannels)
			change_tips = change_tips || !reflect.DeepEqual(tiptypes, newtiptypes)

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
				tipdrp, err := DropTips(tiptypes, prms, channels)

				if err != nil {
					return ret, err
				}
				ret = append(ret, tipdrp)

				tipget, err := GetTips(ctx, newtiptypes, prms, newchannels, usetiptracking)

				if err != nil {
					return ret, err
				}

				ret = append(ret, tipget...)
				//		tips = newtips

				n_tip_uses = 0
				last_thing = nil
				dirty = false
			}
			mci := NewMultiChannelTransferInstruction()
			//vols.SetEqualTo(vol, ins.Multi)
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
			mci.TipType = newtiptypes
			//mci.Multi = ins.Multi
			mci.Multi = countMulti(ins.PltFrom[t])
			channelprms := make([]*wtype.LHChannelParameter, newchannels[0].Multi)
			//mci.Prms = newchannel.MergeWithTip(newtip)

			for i := 0; i < len(newchannels); i++ {
				if newchannels[i] != nil {
					channelprms[i] = newchannels[i].MergeWithTip(newtips[i])
				}
			}

			mci.Prms = channelprms

			ret = append(ret, mci)
			n_tip_uses++

			// finally check if we are touching a bad liquid
			// in future we will do this properly, for now we assume
			// touching any liquid is bad

			npre, premix := pol["PRE_MIX"]
			npost, postmix := pol["POST_MIX"]

			if pol["DSPREFERENCE"].(int) == 0 && !VolumeSet(ins.TVolume[t]).IsZero() || premix && npre.(int) > 0 || postmix && npost.(int) > 0 {
				dirty = true
			}

			last_thing = this_thing

			tiptypes = newtiptypes
			channels = newchannels
			fvols.SubA(vols)
			tvols.AddA(vols)
		}
	}

	// remove tips
	tipdrp, err := DropTips(tiptypes, prms, channels)

	if err != nil {
		return ret, err
	}

	ret = append(ret, tipdrp)

	return ret, nil
}

type SingleChannelTransferInstruction struct {
	BaseRobotInstruction
	*InstructionType
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
	v := &SingleChannelTransferInstruction{
		InstructionType: SCT,
	}
	v.BaseRobotInstruction = NewBaseRobotInstruction(v)
	return v
}

func (ins *SingleChannelTransferInstruction) Visit(visitor RobotInstructionVisitor) {
	visitor.SingleChannelTransfer(ins)
}

func (ins *SingleChannelTransferInstruction) GetParameter(name InstructionParameter) interface{} {
	switch name {
	case LIQUIDCLASS:
		return ins.What
	case VOLUME:
		return ins.Volume
	case VOLUNT:
		return nil
	case FROMPLATETYPE:
		return ins.FPlateType
	case WELLFROMVOLUME:
		return ins.FVolume
	case POSFROM:
		return ins.PltFrom
	case POSTO:
		return ins.PltTo
	case WELLFROM:
		return ins.WellFrom
	case PARAMS:
		return ins.Prms
	case PLATFORM:
		if ins.Prms == nil {
			return ""
		}
		return ins.Prms.Platform
	case WELLTO:
		return ins.WellTo
	case WELLTOVOLUME:
		return ins.TVolume
	case TOPLATETYPE:
		return ins.TPlateType
	case TIPTYPE:
		return ins.TipType
	default:
		return ins.BaseRobotInstruction.GetParameter(name)
	}
}

func (ins *SingleChannelTransferInstruction) Generate(ctx context.Context, policy *wtype.LHPolicyRuleSet, prms *LHProperties) ([]RobotInstruction, error) {
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
	BaseRobotInstruction
	*InstructionType
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
	Multi      int // potentially deprecated
	Prms       []*wtype.LHChannelParameter
	TipType    []string
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
	tp.Channel = scti.Prms[k].Dup()
	tp.TipType = scti.TipType[k]
	return tp
}
func NewMultiChannelTransferInstruction() *MultiChannelTransferInstruction {
	v := &MultiChannelTransferInstruction{
		InstructionType: MCT,
		What:            []string{},
		PltFrom:         []string{},
		PltTo:           []string{},
		WellFrom:        []string{},
		WellTo:          []string{},
		Volume:          []wunit.Volume{},
		FVolume:         []wunit.Volume{},
		TVolume:         []wunit.Volume{},
		FPlateType:      []string{},
		TPlateType:      []string{},
		TipType:         []string{},
	}
	v.BaseRobotInstruction = NewBaseRobotInstruction(v)
	return v
}

func (ins *MultiChannelTransferInstruction) Visit(visitor RobotInstructionVisitor) {
	visitor.MultiChannelTransfer(ins)
}

func (ins *MultiChannelTransferInstruction) GetParameter(name InstructionParameter) interface{} {
	switch name {
	case LIQUIDCLASS:
		return ins.What
	case VOLUME:
		return ins.Volume
	case VOLUNT:
		return nil
	case FROMPLATETYPE:
		return ins.FPlateType
	case WELLFROMVOLUME:
		return ins.FVolume
	case POSFROM:
		return ins.PltFrom
	case POSTO:
		return ins.PltTo
	case WELLFROM:
		return ins.WellFrom
	case PARAMS:
		return ins.Prms
	case PLATFORM:
		if ins.Prms == nil {
			return ""
		}
		return ins.Prms[0].Platform
	case WELLTO:
		return ins.WellTo
	case WELLTOVOLUME:
		return ins.TVolume
	case TOPLATETYPE:
		return ins.TPlateType
	default:
		return ins.BaseRobotInstruction.GetParameter(name)
	}
}

func (ins *MultiChannelTransferInstruction) Generate(ctx context.Context, policy *wtype.LHPolicyRuleSet, prms *LHProperties) ([]RobotInstruction, error) {
	ret := make([]RobotInstruction, 0)

	if len(ins.Volume) == 0 {
		return ret, nil
	}

	// make the instructions

	suckinstruction := NewSuckInstruction()
	blowinstruction := NewBlowInstruction()
	suckinstruction.Multi = ins.Multi
	blowinstruction.Multi = ins.Multi

	c := 0
	for i := 0; i < len(ins.Volume); i++ {
		if ins.Volume[i].IsZero() {
			continue
		}
		c += 1
		suckinstruction.AddTransferParams(ins.Params(i))
		blowinstruction.AddTransferParams(ins.Params(i))
	}

	ret = append(ret, suckinstruction)
	ret = append(ret, blowinstruction)

	return ret, nil
}

type StateChangeInstruction struct {
	BaseRobotInstruction
	*InstructionType
	OldState *wtype.LHChannelParameter
	NewState *wtype.LHChannelParameter
}

func NewStateChangeInstruction(oldstate, newstate *wtype.LHChannelParameter) *StateChangeInstruction {
	v := &StateChangeInstruction{
		InstructionType: CCC,
		OldState:        oldstate,
		NewState:        newstate,
	}
	v.BaseRobotInstruction = NewBaseRobotInstruction(v)
	return v
}

func (ins *StateChangeInstruction) Visit(visitor RobotInstructionVisitor) {
	visitor.StateChange(ins)
}

func (ins *StateChangeInstruction) GetParameter(name InstructionParameter) interface{} {
	switch name {
	case OLDSTATE:
		return ins.OldState
	case NEWSTATE:
		return ins.NewState
	default:
		return ins.BaseRobotInstruction.GetParameter(name)
	}
}

func (ins *StateChangeInstruction) Generate(ctx context.Context, policy *wtype.LHPolicyRuleSet, prms *LHProperties) ([]RobotInstruction, error) {
	return nil, nil
}

type ChangeAdaptorInstruction struct {
	BaseRobotInstruction
	*InstructionType
	Head           int
	DropPosition   string
	GetPosition    string
	OldAdaptorType string
	NewAdaptorType string
	Platform       string
}

func NewChangeAdaptorInstruction(head int, droppos, getpos, oldad, newad, platform string) *ChangeAdaptorInstruction {
	v := &ChangeAdaptorInstruction{
		InstructionType: CHA,
		Head:            head,
		DropPosition:    droppos,
		GetPosition:     getpos,
		OldAdaptorType:  oldad,
		NewAdaptorType:  newad,
		Platform:        platform,
	}
	v.BaseRobotInstruction = NewBaseRobotInstruction(v)
	return v
}

func (ins *ChangeAdaptorInstruction) Visit(visitor RobotInstructionVisitor) {
	visitor.ChangeAdaptor(ins)
}

func (ins *ChangeAdaptorInstruction) GetParameter(name InstructionParameter) interface{} {
	switch name {
	case HEAD:
		return ins.Head
	case POSFROM:
		return ins.DropPosition
	case POSTO:
		return ins.GetPosition
	case OLDADAPTOR:
		return ins.OldAdaptorType
	case NEWADAPTOR:
		return ins.NewAdaptorType
	case PLATFORM:
		return ins.Platform
	default:
		return ins.BaseRobotInstruction.GetParameter(name)
	}
}

func (ins *ChangeAdaptorInstruction) Generate(ctx context.Context, policy *wtype.LHPolicyRuleSet, prms *LHProperties) ([]RobotInstruction, error) {
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
	BaseRobotInstruction
	*InstructionType
	Head       int
	Well       []string
	FPosition  []string
	FPlateType []string
	Multi      int
	Platform   string
}

func NewLoadTipsMoveInstruction() *LoadTipsMoveInstruction {
	v := &LoadTipsMoveInstruction{
		InstructionType: LDT,
		Well:            []string{},
		FPosition:       []string{},
		FPlateType:      []string{},
	}
	v.BaseRobotInstruction = NewBaseRobotInstruction(v)
	return v
}

func (ins *LoadTipsMoveInstruction) Visit(visitor RobotInstructionVisitor) {
	visitor.LoadTipsMove(ins)
}

func (ins *LoadTipsMoveInstruction) GetParameter(name InstructionParameter) interface{} {
	switch name {
	case HEAD:
		return ins.Head
	case FROMPLATETYPE:
		return ins.FPlateType
	case POSFROM:
		return ins.FPosition
	case WELLFROM:
		return ins.Well
	case MULTI:
		return ins.Multi
	case PLATFORM:
		return ins.Platform
	default:
		return ins.BaseRobotInstruction.GetParameter(name)
	}
}

func (ins *LoadTipsMoveInstruction) Generate(ctx context.Context, policy *wtype.LHPolicyRuleSet, prms *LHProperties) ([]RobotInstruction, error) {
	ret := make([]RobotInstruction, 2)

	// move to just above the tip

	mov := NewMoveInstruction()
	mov.Head = ins.Head
	mov.Pos = ins.FPosition
	mov.Well = ins.Well
	mov.Plt = ins.FPlateType
	for i := 0; i < len(ins.Well); i++ {
		mov.Reference = append(mov.Reference, wtype.TopReference.AsInt())
		mov.OffsetX = append(mov.OffsetX, 0.0)
		mov.OffsetY = append(mov.OffsetY, 0.0)
		mov.OffsetZ = append(mov.OffsetZ, 5.0)
	}
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
	BaseRobotInstruction
	*InstructionType
	Head       int
	PltTo      []string
	WellTo     []string
	TPlateType []string
	Multi      int
	Platform   string
}

func NewUnloadTipsMoveInstruction() *UnloadTipsMoveInstruction {
	v := &UnloadTipsMoveInstruction{
		InstructionType: UDT,
		PltTo:           []string{},
		WellTo:          []string{},
		TPlateType:      []string{},
	}
	v.BaseRobotInstruction = NewBaseRobotInstruction(v)
	return v
}

func (ins *UnloadTipsMoveInstruction) Visit(visitor RobotInstructionVisitor) {
	visitor.UnloadTipsMove(ins)
}

func (ins *UnloadTipsMoveInstruction) GetParameter(name InstructionParameter) interface{} {
	switch name {
	case HEAD:
		return ins.Head
	case TOPLATETYPE:
		return ins.TPlateType
	case POSTO:
		return ins.PltTo
	case WELLTO:
		return ins.WellTo
	case MULTI:
		return ins.Multi
	case PLATFORM:
		return ins.Platform
	default:
		return ins.BaseRobotInstruction.GetParameter(name)
	}
}

func (ins *UnloadTipsMoveInstruction) Generate(ctx context.Context, policy *wtype.LHPolicyRuleSet, prms *LHProperties) ([]RobotInstruction, error) {
	ret := make([]RobotInstruction, 2)

	// move

	mov := NewMoveInstruction()
	mov.Head = ins.Head
	mov.Pos = ins.PltTo
	mov.Well = ins.WellTo
	mov.Plt = ins.TPlateType
	for i := 0; i < len(mov.Pos); i++ {
		mov.Reference = append(mov.Reference, wtype.TopReference.AsInt())
		mov.OffsetX = append(mov.OffsetX, 0.0)
		mov.OffsetY = append(mov.OffsetY, 0.0)
		mov.OffsetZ = append(mov.OffsetZ, 0.0)
	}
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
	BaseRobotInstruction
	*InstructionType
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
	v := &AspirateInstruction{
		InstructionType: ASP,
		Volume:          []wunit.Volume{},
		Plt:             []string{},
		What:            []string{},
		LLF:             []bool{},
	}
	v.BaseRobotInstruction = NewBaseRobotInstruction(v)
	return v
}

func (ins *AspirateInstruction) Visit(visitor RobotInstructionVisitor) {
	visitor.Aspirate(ins)
}

func (ins *AspirateInstruction) GetParameter(name InstructionParameter) interface{} {
	switch name {
	case VOLUME:
		return ins.Volume
	case LIQUIDCLASS:
		return ins.What
	case HEAD:
		return ins.Head
	case MULTI:
		return ins.Multi
	case OVERSTROKE:
		return ins.Overstroke
	case WHAT:
		return ins.What
	case PLATE:
		return ins.Plt
	case LLF:
		return ins.LLF
	case PLATFORM:
		return ins.Platform
	default:
		return ins.BaseRobotInstruction.GetParameter(name)
	}
}

func (ins *AspirateInstruction) Generate(ctx context.Context, policy *wtype.LHPolicyRuleSet, prms *LHProperties) ([]RobotInstruction, error) {
	return nil, nil
}

func (ins *AspirateInstruction) OutputTo(lhdriver LiquidhandlingDriver) error {
	driver, ok := lhdriver.(LowLevelLiquidhandlingDriver)
	if !ok {
		return fmt.Errorf("Wrong instruction type for driver: need Lowlevel, got %T", ins)
	}
	volumes := make([]float64, len(ins.Volume))
	for i, vol := range ins.Volume {
		volumes[i] = vol.ConvertToString("ul")
	}
	os := []bool{ins.Overstroke}

	ret := driver.Aspirate(volumes, os, ins.Head, ins.Multi, ins.Plt, ins.What, ins.LLF)
	if !ret.OK {
		return fmt.Errorf(" %d : %s", ret.Errorcode, ret.Msg)
	}

	return nil
}

type DispenseInstruction struct {
	BaseRobotInstruction
	*InstructionType
	Head     int
	Volume   []wunit.Volume
	Multi    int
	Plt      []string
	What     []string
	LLF      []bool
	Platform string
}

func NewDispenseInstruction() *DispenseInstruction {
	v := &DispenseInstruction{
		InstructionType: DSP,
		Volume:          []wunit.Volume{},
		Plt:             []string{},
		What:            []string{},
		LLF:             []bool{},
	}
	v.BaseRobotInstruction = NewBaseRobotInstruction(v)
	return v
}

func (ins *DispenseInstruction) Visit(visitor RobotInstructionVisitor) {
	visitor.Dispense(ins)
}

func (ins *DispenseInstruction) GetParameter(name InstructionParameter) interface{} {
	switch name {
	case VOLUME:
		return ins.Volume
	case LIQUIDCLASS:
		return ins.What
	case HEAD:
		return ins.Head
	case MULTI:
		return ins.Multi
	case WHAT:
		return ins.What
	case LLF:
		return ins.LLF
	case PLT:
		return ins.Plt
	case PLATE:
		return ins.Plt
	case PLATFORM:
		return ins.Platform
	default:
		return ins.BaseRobotInstruction.GetParameter(name)
	}
}

func (ins *DispenseInstruction) Generate(ctx context.Context, policy *wtype.LHPolicyRuleSet, prms *LHProperties) ([]RobotInstruction, error) {
	return nil, nil
}

func (ins *DispenseInstruction) OutputTo(lhdriver LiquidhandlingDriver) error {
	driver, ok := lhdriver.(LowLevelLiquidhandlingDriver)
	if !ok {
		return fmt.Errorf("Wrong instruction type for driver: need Lowlevel, got %T", ins)
	}
	volumes := make([]float64, len(ins.Volume))
	for i, vol := range ins.Volume {
		volumes[i] = vol.ConvertToString("ul")
	}

	os := []bool{false}
	ret := driver.Dispense(volumes, os, ins.Head, ins.Multi, ins.Plt, ins.What, ins.LLF)
	if !ret.OK {
		return fmt.Errorf(" %d : %s", ret.Errorcode, ret.Msg)
	}

	return nil

}

type BlowoutInstruction struct {
	BaseRobotInstruction
	*InstructionType
	Head     int
	Volume   []wunit.Volume
	Multi    int
	Plt      []string
	What     []string
	LLF      []bool
	Platform string
}

func NewBlowoutInstruction() *BlowoutInstruction {
	v := &BlowoutInstruction{
		InstructionType: BLO,
		Volume:          []wunit.Volume{},
	}
	v.BaseRobotInstruction = NewBaseRobotInstruction(v)
	return v
}

func (ins *BlowoutInstruction) Visit(visitor RobotInstructionVisitor) {
	visitor.Blowout(ins)
}

func (ins *BlowoutInstruction) GetParameter(name InstructionParameter) interface{} {
	switch name {
	case VOLUME:
		return ins.Volume
	case HEAD:
		return ins.Head
	case MULTI:
		return ins.Multi
	case WHAT:
		return ins.What
	case LLF:
		return ins.LLF
	case PLT:
		return ins.Plt
	case PLATFORM:
		return ins.Platform
	default:
		return ins.BaseRobotInstruction.GetParameter(name)
	}
}

func (ins *BlowoutInstruction) Generate(ctx context.Context, policy *wtype.LHPolicyRuleSet, prms *LHProperties) ([]RobotInstruction, error) {
	return nil, nil
}

func (ins *BlowoutInstruction) OutputTo(lhdriver LiquidhandlingDriver) error {
	driver, ok := lhdriver.(LowLevelLiquidhandlingDriver)
	if !ok {
		return fmt.Errorf("Wrong instruction type for driver: need Lowlevel, got %T", ins)
	}
	volumes := make([]float64, len(ins.Volume))
	for i, vol := range ins.Volume {
		volumes[i] = vol.ConvertToString("ul")
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
	BaseRobotInstruction
	*InstructionType
	Head    int
	Channel int
}

func NewPTZInstruction() *PTZInstruction {
	v := &PTZInstruction{
		InstructionType: PTZ,
	}
	v.BaseRobotInstruction = NewBaseRobotInstruction(v)
	return v
}

func (ins *PTZInstruction) Visit(visitor RobotInstructionVisitor) {
	visitor.PTZ(ins)
}

func (ins *PTZInstruction) GetParameter(name InstructionParameter) interface{} {
	switch name {
	case HEAD:
		return ins.Head
	case CHANNEL:
		return ins.Channel
	default:
		return ins.BaseRobotInstruction.GetParameter(name)
	}
}

func (ins *PTZInstruction) Generate(ctx context.Context, policy *wtype.LHPolicyRuleSet, prms *LHProperties) ([]RobotInstruction, error) {
	return nil, nil
}

func (ins *PTZInstruction) OutputTo(lhdriver LiquidhandlingDriver) error {
	driver, ok := lhdriver.(LowLevelLiquidhandlingDriver)
	if !ok {
		return fmt.Errorf("Wrong instruction type for driver: need Lowlevel, got %T", ins)
	}
	ret := driver.ResetPistons(ins.Head, ins.Channel)
	if !ret.OK {
		return fmt.Errorf(" %d : %s", ret.Errorcode, ret.Msg)
	}

	return nil
}

type MoveInstruction struct {
	BaseRobotInstruction
	*InstructionType
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
	v := &MoveInstruction{
		InstructionType: MOV,
		Plt:             []string{},
		Pos:             []string{},
		Well:            []string{},
		WVolume:         []wunit.Volume{},
		Reference:       []int{},
		OffsetX:         []float64{},
		OffsetY:         []float64{},
		OffsetZ:         []float64{},
	}
	v.BaseRobotInstruction = NewBaseRobotInstruction(v)
	return v
}

func (ins *MoveInstruction) Visit(visitor RobotInstructionVisitor) {
	visitor.Move(ins)
}

func (ins *MoveInstruction) GetParameter(name InstructionParameter) interface{} {
	switch name {
	case WELLTOVOLUME:
		return ins.WVolume
	case HEAD:
		return ins.Head
	case TOPLATETYPE:
		return ins.Plt
	case POSTO:
		return ins.Pos
	case WELLTO:
		return ins.Well
	case REFERENCE:
		return ins.Reference
	case OFFSETX:
		return ins.OffsetX
	case OFFSETY:
		return ins.OffsetY
	case OFFSETZ:
		return ins.OffsetZ
	case PLATFORM:
		return ins.Platform
	default:
		return ins.BaseRobotInstruction.GetParameter(name)
	}
}

func (ins *MoveInstruction) MaybeMerge(next RobotInstruction) RobotInstruction {
	switch n := next.(type) {
	case *AspirateInstruction:
		return NewMovAsp(ins, n)
	case *DispenseInstruction:
		return NewMovDsp(ins, n)
	case *MixInstruction:
		return NewMovMix(ins, n)
	case *BlowoutInstruction:
		return NewMovBlo(ins, n)
	default:
		return ins
	}
}

func (ins *MoveInstruction) Generate(ctx context.Context, policy *wtype.LHPolicyRuleSet, prms *LHProperties) ([]RobotInstruction, error) {
	return nil, nil
}

func (ins *MoveInstruction) OutputTo(lhdriver LiquidhandlingDriver) error {
	driver, ok := lhdriver.(LowLevelLiquidhandlingDriver)
	if !ok {
		return fmt.Errorf("Wrong instruction type for driver: need Lowlevel, got %T", ins)
	}
	ret := driver.Move(ins.Pos, ins.Well, ins.Reference, ins.OffsetX, ins.OffsetY, ins.OffsetZ, ins.Plt, ins.Head)
	if !ret.OK {
		return fmt.Errorf(" %d : %s", ret.Errorcode, ret.Msg)
	}

	return nil

}

type MoveRawInstruction struct {
	BaseRobotInstruction
	*InstructionType
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
	v := &MoveRawInstruction{
		InstructionType: MRW,
		What:            []string{},
		PltFrom:         []string{},
		PltTo:           []string{},
		WellFrom:        []string{},
		WellTo:          []string{},
		FPlateType:      []string{},
		TPlateType:      []string{},
		Volume:          []wunit.Volume{},
		FVolume:         []wunit.Volume{},
		TVolume:         []wunit.Volume{},
	}
	v.BaseRobotInstruction = NewBaseRobotInstruction(v)
	return v
}

func (ins *MoveRawInstruction) Visit(visitor RobotInstructionVisitor) {
	visitor.MoveRaw(ins)
}

func (ins *MoveRawInstruction) GetParameter(name InstructionParameter) interface{} {
	switch name {
	case HEAD:
		return ins.Head
	case LIQUIDCLASS:
		return ins.What
	case VOLUME:
		return ins.Volume
	case FROMPLATETYPE:
		return ins.FPlateType
	case TOPLATETYPE:
		return ins.TPlateType
	case WELLFROMVOLUME:
		return ins.FVolume
	case WELLTOVOLUME:
		return ins.TVolume
	case POSFROM:
		return ins.PltFrom
	case POSTO:
		return ins.PltTo
	case WELLFROM:
		return ins.WellFrom
	case PARAMS:
		return ins.Prms
	case PLATFORM:
		if ins.Prms == nil {
			return ""
		}
		return ins.Prms.Platform
	default:
		return ins.BaseRobotInstruction.GetParameter(name)
	}
}

func (ins *MoveRawInstruction) Generate(ctx context.Context, policy *wtype.LHPolicyRuleSet, prms *LHProperties) ([]RobotInstruction, error) {
	return nil, nil
}

func (ins *MoveRawInstruction) OutputTo(lhdriver LiquidhandlingDriver) error {
	/*
		driver, ok := lhdriver.(LowLevelLiquidhandlingDriver)
		if !ok {
			return fmt.Errorf("Wrong instruction type for driver: need Lowlevel, got %T", ins)
		}
	*/
	logger.Fatal("Not yet implemented")
	panic("Not yet implemented")
}

type LoadTipsInstruction struct {
	BaseRobotInstruction
	*InstructionType
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
	v := &LoadTipsInstruction{
		InstructionType: LOD,
		Channels:        []int{},
		TipType:         []string{},
		HolderType:      []string{},
		Pos:             []string{},
		Well:            []string{},
	}
	v.BaseRobotInstruction = NewBaseRobotInstruction(v)
	return v
}

func (ins *LoadTipsInstruction) Visit(visitor RobotInstructionVisitor) {
	visitor.LoadTips(ins)
}

func (ins *LoadTipsInstruction) GetParameter(name InstructionParameter) interface{} {
	switch name {
	case HEAD:
		return ins.Head
	case CHANNEL:
		return ins.Channels
	case TIPTYPE:
		return ins.TipType
	case FROMPLATETYPE:
		return ins.HolderType
	case MULTI:
		return ins.Multi
	case WELL:
		return ins.Well
	case PLATE:
		return ins.HolderType
	case POS:
		return ins.Pos
	case PLATFORM:
		return ins.Platform
	default:
		return ins.BaseRobotInstruction.GetParameter(name)
	}
}

func (ins *LoadTipsInstruction) Generate(ctx context.Context, policy *wtype.LHPolicyRuleSet, prms *LHProperties) ([]RobotInstruction, error) {
	return nil, nil
}

func (ins *LoadTipsInstruction) OutputTo(lhdriver LiquidhandlingDriver) error {
	driver, ok := lhdriver.(LowLevelLiquidhandlingDriver)
	if !ok {
		return fmt.Errorf("Wrong instruction type for driver: need Lowlevel, got %T", ins)
	}
	ret := driver.LoadTips(ins.Channels, ins.Head, ins.Multi, ins.HolderType, ins.Pos, ins.Well)
	if !ret.OK {
		return fmt.Errorf(" %d : %s", ret.Errorcode, ret.Msg)
	}

	return nil

}

type UnloadTipsInstruction struct {
	BaseRobotInstruction
	*InstructionType
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
	v := &UnloadTipsInstruction{
		InstructionType: ULD,
		TipType:         []string{},
		HolderType:      []string{},
		Channels:        []int{},
		Pos:             []string{},
		Well:            []string{},
	}
	v.BaseRobotInstruction = NewBaseRobotInstruction(v)
	return v
}

func (ins *UnloadTipsInstruction) Visit(visitor RobotInstructionVisitor) {
	visitor.UnloadTips(ins)
}

func (ins *UnloadTipsInstruction) GetParameter(name InstructionParameter) interface{} {
	switch name {
	case HEAD:
		return ins.Head
	case CHANNEL:
		return ins.Channels
	case TIPTYPE:
		return ins.TipType
	case TOPLATETYPE:
		return ins.HolderType
	case MULTI:
		return ins.Multi
	case WELL:
		return ins.Well
	case POS:
		return ins.Pos
	case PLATFORM:
		return ins.Platform
	default:
		return ins.BaseRobotInstruction.GetParameter(name)
	}
}

func (ins *UnloadTipsInstruction) Generate(ctx context.Context, policy *wtype.LHPolicyRuleSet, prms *LHProperties) ([]RobotInstruction, error) {
	return nil, nil
}

func (ins *UnloadTipsInstruction) OutputTo(lhdriver LiquidhandlingDriver) error {
	driver, ok := lhdriver.(LowLevelLiquidhandlingDriver)
	if !ok {
		return fmt.Errorf("Wrong instruction type for driver: need Lowlevel, got %T", ins)
	}
	ret := driver.UnloadTips(ins.Channels, ins.Head, ins.Multi, ins.HolderType, ins.Pos, ins.Well)
	if !ret.OK {
		return fmt.Errorf(" %d : %s", ret.Errorcode, ret.Msg)
	}

	return nil

}

type SuckInstruction struct {
	BaseRobotInstruction
	*InstructionType
	Head        int
	What        []string
	ComponentID []string // ID, not currently used. Will be needed soon.
	PltFrom     []string
	WellFrom    []string
	Volume      []wunit.Volume
	FPlateType  []string
	FVolume     []wunit.Volume
	Prms        *wtype.LHChannelParameter
	Multi       int
	Overstroke  bool
	TipType     string
}

func NewSuckInstruction() *SuckInstruction {
	v := &SuckInstruction{
		InstructionType: SUK,
		What:            []string{},
		PltFrom:         []string{},
		WellFrom:        []string{},
		Volume:          []wunit.Volume{},
		FPlateType:      []string{},
		FVolume:         []wunit.Volume{},
	}
	v.BaseRobotInstruction = NewBaseRobotInstruction(v)
	return v
}

func (ins *SuckInstruction) Visit(visitor RobotInstructionVisitor) {
	visitor.Suck(ins)
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

func (ins *SuckInstruction) GetParameter(name InstructionParameter) interface{} {
	switch name {
	case HEAD:
		return ins.Head
	case LIQUIDCLASS:
		return ins.What
	case VOLUME:
		return ins.Volume
	case FROMPLATETYPE:
		return ins.FPlateType
	case WELLFROMVOLUME:
		return ins.FVolume
	case POSFROM:
		return ins.PltFrom
	case WELLFROM:
		return ins.WellFrom
	case PARAMS:
		return ins.Prms
	case MULTI:
		return ins.Multi
	case OVERSTROKE:
		return ins.Overstroke
	case PLATFORM:
		if ins.Prms == nil {
			return ""
		}
		return ins.Prms.Platform
	case TIPTYPE:
		return ins.TipType
	case WHICH:
		return ins.ComponentID
	default:
		return ins.BaseRobotInstruction.GetParameter(name)
	}
}

func (ins *SuckInstruction) getInitialVolumesByWell() map[string]wunit.Volume {
	ret := make(map[string]wunit.Volume, ins.Multi)
	for i, well := range ins.WellFrom {
		//overwrite, since these should be the same for each well
		ret[well] = wunit.CopyVolume(ins.FVolume[i])
	}
	return ret
}

func (ins *SuckInstruction) getRemovedVolumesByWell() map[string]wunit.Volume {
	ret := make(map[string]wunit.Volume, ins.Multi)
	for i, well := range ins.WellFrom {
		if _, ok := ret[well]; !ok {
			ret[well] = wunit.CopyVolume(ins.Volume[i])
		} else {
			ret[well].IncrBy(ins.Volume[i])
		}
	}
	return ret
}

func (ins *SuckInstruction) getFinalVolumesByWell() map[string]wunit.Volume {
	ret := make(map[string]wunit.Volume, ins.Multi)
	initial := ins.getInitialVolumesByWell()
	removed := ins.getRemovedVolumesByWell()
	for _, well := range ins.WellFrom {
		ret[well] = wunit.SubtractVolumes(initial[well], removed[well])
	}
	return ret
}

func (ins *SuckInstruction) getAspirate(volumes []wunit.Volume, useLLF []bool) RobotInstruction {
	aspins := NewAspirateInstruction()
	aspins.Head = ins.Head
	aspins.Multi = ins.Multi
	aspins.Overstroke = ins.Overstroke
	aspins.What = ins.What
	aspins.Plt = ins.FPlateType

	for _, vol := range volumes {
		aspins.Volume = append(aspins.Volume, wunit.CopyVolume(vol))
	}
	for _, llf := range useLLF {
		aspins.LLF = append(aspins.LLF, llf)
	}
	for len(aspins.LLF) < aspins.Multi {
		aspins.LLF = append(aspins.LLF, false)
	}

	return aspins
}

func (ins *SuckInstruction) getMove(reference int, volume []wunit.Volume, offsetX, offsetY, offsetZ float64) RobotInstruction {
	mov := NewMoveInstruction()
	mov.Head = ins.Head

	mov.Pos = ins.PltFrom
	mov.Plt = ins.FPlateType
	mov.Well = ins.WellFrom
	mov.WVolume = volume

	for i := 0; i < ins.Multi; i++ {
		mov.Reference = append(mov.Reference, reference)
		mov.OffsetX = append(mov.OffsetX, offsetX)
		mov.OffsetY = append(mov.OffsetY, offsetY)
		mov.OffsetZ = append(mov.OffsetZ, offsetZ)
	}
	return mov
}

func (ins *SuckInstruction) getMoveMix(volume []wunit.Volume, cycles int, offsetX, offsetY, offsetZ float64, blowout bool) RobotInstruction {
	mix := NewMoveMixInstruction()
	mix.Head = ins.Head
	mix.Plt = ins.PltFrom
	mix.PlateType = ins.FPlateType
	mix.Well = ins.WellFrom
	mix.Multi = ins.Multi
	mix.What = ins.What
	mix.Volume = volume

	for k := 0; k < ins.Multi; k++ {
		mix.OffsetX = append(mix.OffsetX, offsetX)
		mix.OffsetY = append(mix.OffsetY, offsetY)
		mix.OffsetZ = append(mix.OffsetZ, offsetZ)
		mix.Cycles = append(mix.Cycles, cycles)
		mix.Blowout = append(mix.Blowout, blowout)
	}

	return mix
}

func (ins *SuckInstruction) getChannelsByWell() map[string][]int {
	ret := make(map[string][]int, ins.Multi)
	for i, well := range ins.WellFrom {
		ret[well] = append(ret[well], i)
	}
	return ret
}

func (ins *SuckInstruction) Generate(ctx context.Context, policy *wtype.LHPolicyRuleSet, prms *LHProperties) ([]RobotInstruction, error) {
	// MIS XXX -- separate out channel-level parameters from head-level ones
	ret := make([]RobotInstruction, 0, 1)

	// this is where the policies come into effect

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

	// set the defaults
	ret = append(ret, setDefaults(ins.Head, pol)...)
	defaultpspeed := SafeGetF64(pol, "DEFAULTPIPETTESPEED")

	allowOutOfRangePipetteSpeeds := SafeGetBool(pol, "OVERRIDEPIPETTESPEED")

	head := prms.GetLoadedHead(ins.Head)

	defaultpspeed, err = checkAndSaften(defaultpspeed, head.Params.Minspd.RawValue(), head.Params.Maxspd.RawValue(), allowOutOfRangePipetteSpeeds)

	if err != nil {
		return []RobotInstruction{}, errors.Wrap(err, fmt.Sprintf("setting default pipette speed for policy %s", text.PrettyPrint(pol)))
	}

	// offsets
	ofx := SafeGetF64(pol, "ASPXOFFSET")
	ofy := SafeGetF64(pol, "ASPYOFFSET")
	ofzadj := SafeGetF64(pol, "OFFSETZADJUST")
	ofz := SafeGetF64(pol, "ASPZOFFSET") + ofzadj
	ofz, err = makeZOffsetSafe(ctx, prms, ofz, ins.Head, ins.PltFrom, ins.TipType)
	if err != nil {
		return nil, err
	}

	// do we need to enter slowly?
	if entryspeed, gentlynow := pol["ASPENTRYSPEED"]; gentlynow {
		// go to the well top
		ret = append(ret, ins.getMove(1, ins.FVolume, ofx, ofy, 5.0))

		// set the speed
		ret = append(ret, NewSetDriveSpeedInstruction("Z", entryspeed.(float64)))
	}

	// do we pre-mix?
	if cycles := SafeGetInt(pol, "PRE_MIX"); cycles > 0 {
		// add the premix step

		volume := make([]wunit.Volume, ins.Multi)
		if _, preMixVolSet := pol["PRE_MIX_VOLUME"]; preMixVolSet {
			mixvol := SafeGetF64(pol, "PRE_MIX_VOLUME")

			// TODO -- corresponding checks when set
			if mixvol < wtype.Globals.MIN_REASONABLE_VOLUME_UL {
				return ret, wtype.LHError(wtype.LH_ERR_POLICY, fmt.Sprintf("PRE_MIX_VOLUME set below minimum allowed: %f min %f", mixvol, wtype.Globals.MIN_REASONABLE_VOLUME_UL))
			} else if vmixvol := wunit.NewVolume(mixvol, "ul"); !ins.Prms.CanMove(vmixvol, true) {
				if SafeGetBool(pol, "MIX_VOLUME_OVERRIDE_TIP_MAX") {
					mixvol = ins.Prms.Maxvol.ConvertToString("ul")
				} else {
					// this is an error in channel choice but the user has to deal... needs modificationst
					return ret, wtype.LHError(wtype.LH_ERR_POLICY, fmt.Sprintf("PRE_MIX_VOLUME not compatible with optimal channel choice: requested %s channel limits are %s", vmixvol.ToString(), ins.Prms.VolumeLimitString()))
				}
			}

			for i := range volume {
				volume[i] = wunit.NewVolume(mixvol, "ul")
			}
		} else {
			for i := range volume {
				volume[i] = wunit.CopyVolume(ins.Volume[i])
			}
		}

		mixofx := SafeGetF64(pol, "PRE_MIX_X")
		mixofy := SafeGetF64(pol, "PRE_MIX_Y")
		mixofz, err := makeZOffsetSafe(ctx, prms, SafeGetF64(pol, "PRE_MIX_Z")+ofzadj, ins.Head, ins.PltFrom, ins.TipType)
		if err != nil {
			return nil, err
		}
		mix := ins.getMoveMix(volume, cycles, mixofx, mixofy, mixofz, false)

		mixrate := SafeGetF64(pol, "PRE_MIX_RATE")

		changepipspeed := (mixrate != defaultpspeed) && (mixrate > 0.0)

		if changepipspeed {
			if mixrate, err := checkAndSaften(mixrate, head.Params.Minspd.RawValue(), head.Params.Maxspd.RawValue(), allowOutOfRangePipetteSpeeds); err != nil {
				return []RobotInstruction{}, errors.Wrap(err, "setting pre mix pipetting speed")
			} else {
				ret = append(ret, NewSetPipetteSpeedInstruction(ins.Head, -1, mixrate))
			}
		}

		ret = append(ret, mix)

		if changepipspeed {
			ret = append(ret, NewSetPipetteSpeedInstruction(ins.Head, -1, defaultpspeed))
		}
	}

	// Set the pipette speed if needed
	apspeed := SafeGetF64(pol, "ASPSPEED")
	changepspeed := (apspeed != defaultpspeed) && (apspeed > 0.0)
	if changepspeed {
		if apspeed, err = checkAndSaften(apspeed, head.Params.Minspd.RawValue(), head.Params.Maxspd.RawValue(), allowOutOfRangePipetteSpeeds); err != nil {
			return []RobotInstruction{}, errors.Wrap(err, "setting pipette aspirate speed")
		}
	}

	// get aspirate volume
	volumes := make([]wunit.Volume, 0, len(ins.Volume))
	for _, vol := range ins.Volume {
		volumes = append(volumes, wunit.CopyVolume(vol))
	}
	if ev, iwantmore := pol["EXTRA_ASP_VOLUME"]; iwantmore {
		extra_vol := ev.(wunit.Volume)
		for i := range volumes {
			volumes[i].IncrBy(extra_vol)
		}
	}

	//LLF
	if useLLF, anyLLF := getUseLLF(pol, ins.Multi, ins.PltFrom, prms); anyLLF {
		llfBelowSurface := SafeGetF64(pol, "LLFBELOWSURFACE")
		minHeight := wunit.NewLength(llfBelowSurface+ofz, "mm")
		minVolume := prms.Plates[ins.PltFrom[0]].Welltype.GetLiquidVolume(minHeight) // assume welltypes are all the same

		for _, initial := range ins.getInitialVolumesByWell() {
			anyLLF = anyLLF && !initial.LessThan(minVolume)
		}
		if anyLLF {

			finalVolumes := ins.getFinalVolumesByWell()
			channelsByWell := ins.getChannelsByWell()

			llfVolumes := make([]wunit.Volume, 0, ins.Multi)
			for _, vol := range volumes {
				llfVolumes = append(llfVolumes, wunit.CopyVolume(vol))
			}

			for well, final := range finalVolumes {
				if final.LessThan(minVolume) {
					//cannot aspirate all via LLF, leave the residual until later
					residual := wunit.SubtractVolumes(minVolume, final)
					//evenly split the residual between each channel accessing this well
					residual.DivideBy(float64(len(channelsByWell[well])))
					for _, ch := range channelsByWell[well] {
						volumes[ch] = wunit.CopyVolume(residual)
						llfVolumes[ch].DecrBy(residual)
					}
				} else {
					for _, ch := range channelsByWell[well] {
						volumes[ch] = wunit.ZeroVolume()
					}
				}
			}

			ret = append(ret, ins.getMove(2, ins.FVolume, ofx, ofy, -llfBelowSurface))
			if changepspeed {
				ret = append(ret, NewSetPipetteSpeedInstruction(ins.Head, -1, apspeed))
			}
			ret = append(ret, ins.getAspirate(llfVolumes, useLLF))
		}
	}

	positive := false
	for _, vol := range volumes {
		if vol.IsPositive() {
			positive = true
			break
		}
	}

	if positive {
		ret = append(ret, ins.getMove(SafeGetInt(pol, "ASPREFERENCE"), ins.FVolume, ofx, ofy, ofz))
		if changepspeed {
			ret = append(ret, NewSetPipetteSpeedInstruction(ins.Head, -1, apspeed))
		}
		ret = append(ret, ins.getAspirate(volumes, nil))
	}

	// do we reset the pipette speed?
	if changepspeed {
		ret = append(ret, NewSetPipetteSpeedInstruction(ins.Head, -1, defaultpspeed))
	}

	// do we wait
	if wait_time := SafeGetF64(pol, "ASP_WAIT"); wait_time > 0.0 {
		waitins := NewWaitInstruction()
		waitins.Time = wait_time
		ret = append(ret, waitins)
	}

	if _, gentlynow := pol["ASPENTRYSPEED"]; gentlynow { // reset the drive speed
		// go to the well top
		ret = append(ret, ins.getMove(1, ins.FVolume, ofx, ofy, 5.0))

		// now get ready to move fast again
		ret = append(ret, NewSetDriveSpeedInstruction("Z", pol["DEFAULTZSPEED"].(float64)))
	}

	return ret, nil

}

type BlowInstruction struct {
	BaseRobotInstruction
	*InstructionType
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
	v := &BlowInstruction{
		InstructionType: BLW,
		What:            []string{},
		PltTo:           []string{},
		WellTo:          []string{},
		Volume:          []wunit.Volume{},
		TPlateType:      []string{},
		TVolume:         []wunit.Volume{},
	}
	v.BaseRobotInstruction = NewBaseRobotInstruction(v)
	return v
}

func (ins *BlowInstruction) Visit(visitor RobotInstructionVisitor) {
	visitor.Blow(ins)
}

func (ins *BlowInstruction) GetParameter(name InstructionParameter) interface{} {
	switch name {
	case HEAD:
		return ins.Head
	case LIQUIDCLASS:
		return ins.What
	case VOLUME:
		return ins.Volume
	case TOPLATETYPE:
		return ins.TPlateType
	case WELLTOVOLUME:
		return ins.TVolume
	case POSTO:
		return ins.PltTo
	case WELLTO:
		return ins.WellTo
	case PARAMS:
		return ins.Prms
	case PLATFORM:
		if ins.Prms == nil {
			return ""
		}
		return ins.Prms.Platform
	case MULTI:
		return ins.Multi
	case TIPTYPE:
		return ins.TipType
	default:
		return ins.BaseRobotInstruction.GetParameter(name)
	}
}

func (ins *BlowInstruction) AddTransferParams(tp TransferParams) {
	ins.What = append(ins.What, tp.What)
	ins.PltTo = append(ins.PltTo, tp.PltTo)
	ins.WellTo = append(ins.WellTo, tp.WellTo)
	ins.Volume = append(ins.Volume, tp.Volume)
	ins.TPlateType = append(ins.TPlateType, tp.TPlateType)
	ins.TVolume = append(ins.TVolume, tp.TVolume)
	ins.Prms = tp.Channel
	ins.Head = tp.Channel.Head
	ins.TipType = tp.TipType
}
func (scti *BlowInstruction) Params() MultiTransferParams {
	tp := NewMultiTransferParams(scti.Multi)
	/*
		tp.What = scti.What
		tp.PltTo = scti.PltTo
		tp.WellTo = scti.WellTo
		tp.Volume = scti.Volume
		tp.TPlateType = scti.TPlateType
		tp.TVolume = scti.TVolume
		tp.Channel = scti.Prms
	*/

	for i := 0; i < len(scti.What); i++ {
		tp.Transfers = append(tp.Transfers, TransferParams{What: scti.What[i], PltTo: scti.PltTo[i], WellTo: scti.WellTo[i], Volume: scti.Volume[i], TPlateType: scti.TPlateType[i], TVolume: scti.TVolume[i], Channel: scti.Prms.Dup()})
	}

	return tp
}

func setDefaults(head int, pol wtype.LHPolicy) []RobotInstruction {
	// pipetting speed
	defaultpspeed := SafeGetF64(pol, "DEFAULTPIPETTESPEED")
	sps := NewSetPipetteSpeedInstruction(head, -1, defaultpspeed)

	// Z move speed
	sds := NewSetDriveSpeedInstruction("Z", pol["DEFAULTZSPEED"].(float64))

	return []RobotInstruction{sps, sds}
}

func (ins *BlowInstruction) Generate(ctx context.Context, policy *wtype.LHPolicyRuleSet, prms *LHProperties) ([]RobotInstruction, error) {
	ret := make([]RobotInstruction, 0)
	// apply policies here

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

	allowOutOfRangePipetteSpeeds := SafeGetBool(pol, "OVERRIDEPIPETTESPEED")

	head := prms.GetLoadedHead(ins.Head)

	// change pipette speed?
	defaultpspeed := SafeGetF64(pol, "DEFAULTPIPETTESPEED")
	defaultpspeed, err = checkAndSaften(defaultpspeed, head.Params.Minspd.RawValue(), head.Params.Maxspd.RawValue(), allowOutOfRangePipetteSpeeds)

	if err != nil {
		return []RobotInstruction{}, errors.Wrap(err, "setting pipette aspirate speed")
	}

	// set the defaults
	ret = append(ret, setDefaults(ins.Head, pol)...)

	// first, are we breaking up the move?

	ofx := SafeGetF64(pol, "DSPXOFFSET")
	ofy := SafeGetF64(pol, "DSPYOFFSET")
	ofz := SafeGetF64(pol, "DSPZOFFSET")
	ofzadj := SafeGetF64(pol, "OFFSETZADJUST")
	ofz += ofzadj
	ofz, err = makeZOffsetSafe(ctx, prms, ofz, ins.Head, ins.PltTo, ins.TipType)
	if err != nil {
		return nil, err
	}

	ref := SafeGetInt(pol, "DSPREFERENCE")
	entryspeed := SafeGetF64(pol, "DSPENTRYSPEED")
	defaultspeed := SafeGetF64(pol, "DEFAULTZSPEED")

	//LLF
	use_llf, any_llf := getUseLLF(pol, ins.Multi, ins.PltTo, prms)
	if any_llf {
		//override reference
		ref = 2 //liquid level
		//override ofz
		ofz = +SafeGetF64(pol, "LLFABOVESURFACE")
	}

	var gentlydoesit bool

	if entryspeed > 0.0 && entryspeed != defaultspeed {
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
		ret = append(ret, NewSetDriveSpeedInstruction("Z", entryspeed))

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
			ret = append(ret, NewSetDriveSpeedInstruction("Z", pol["DEFAULTZSPEED"].(float64)))
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

	dpspeed := SafeGetF64(pol, "DSPSPEED")

	var setpspeed bool

	if defaultpspeed != dpspeed && dpspeed != 0.0 {
		setpspeed = true
	}

	if setpspeed {
		dpspeed, err = checkAndSaften(dpspeed, head.Params.Minspd.RawValue(), head.Params.Maxspd.RawValue(), allowOutOfRangePipetteSpeeds)

		if err != nil {
			return []RobotInstruction{}, errors.Wrap(err, "setting pipette dispense speed")
		}

		ret = append(ret, NewSetPipetteSpeedInstruction(ins.Head, -1, dpspeed))
	}

	// now we dispense

	weneedtoreset := true

	justblowout := SafeGetBool(pol, "JUSTBLOWOUT")

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
			boins.LLF = append(boins.LLF, use_llf[i])
		}

		ret = append(ret, boins)
		weneedtoreset = false
	} else {
		dspins := NewDispenseInstruction()
		dspins.Head = ins.Head
		dspins.Volume = ins.Volume

		extra_vol := SafeGetVolume(pol, "EXTRA_DISP_VOLUME")
		if extra_vol.IsPositive() {
			for i := range dspins.Volume {
				dspins.Volume[i].Add(extra_vol)
			}
		}
		dspins.Multi = ins.Multi
		dspins.Plt = ins.TPlateType
		dspins.What = ins.What

		for i := 0; i < ins.Multi; i++ {
			dspins.LLF = append(dspins.LLF, use_llf[i])
		}

		ret = append(ret, dspins)
	}

	// do we reset the pipette speed?

	if setpspeed {
		ret = append(ret, NewSetPipetteSpeedInstruction(ins.Head, -1, defaultpspeed))
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

		pmzoff, err = makeZOffsetSafe(ctx, prms, pmzoff, ins.Head, ins.PltTo, ins.TipType)
		if err != nil {
			return nil, err
		}

		for k := 0; k < ins.Multi; k++ {
			mix.OffsetZ = append(mix.OffsetZ, pmzoff)
		}

		_, ok := pol["POST_MIX_VOLUME"]
		mix.Volume = ins.Volume
		mixvol := SafeGetF64(pol, "POST_MIX_VOLUME")

		if mixvol == 0.0 {
			mixvol = ins.Volume[0].ConvertToString("ul")
		}

		vmixvol := wunit.NewVolume(mixvol, "ul")

		// check the volume

		if mixvol < wtype.Globals.MIN_REASONABLE_VOLUME_UL {
			return ret, wtype.LHError(wtype.LH_ERR_POLICY, fmt.Sprintf("POST_MIX_VOLUME set below minimum allowed: %f min %f", mixvol, wtype.Globals.MIN_REASONABLE_VOLUME_UL))
		} else if !ins.Prms.CanMove(vmixvol, true) {
			override := SafeGetBool(pol, "MIX_VOLUME_OVERRIDE_TIP_MAX")

			//does the tip have a filter?
			inv := inventory.GetInventory(ctx)
			tb, err := inv.NewTipbox(ctx, ins.TipType)
			if err != nil {
				return ret, wtype.LHError(wtype.LH_ERR_OTHER, fmt.Sprintf("While getting tip %v", err))
			}

			//filter tips always override max volume
			if override || tb.Tiptype.Filtered {
				mixvol = ins.Prms.Maxvol.ConvertToString("ul")
			} else {
				return ret, wtype.LHError(wtype.LH_ERR_POLICY, fmt.Sprintf("Setting POST_MIX_VOLUME to %s cannot be achieved with current tip (type %s) volume limits %v, instruction details: %s", vmixvol.ToString(), ins.TipType, ins.Prms, text.PrettyPrint(ins)))
			}
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
		if mixrate != defaultpspeed && mixrate != 0.0 {
			changespeed = true
		}

		if changespeed {
			mixrate, err = checkAndSaften(mixrate, head.Params.Minspd.RawValue(), head.Params.Maxspd.RawValue(), allowOutOfRangePipetteSpeeds)

			if err != nil {
				return []RobotInstruction{}, errors.Wrap(err, "setting post mix pipetting speed")
			}
			ret = append(ret, NewSetPipetteSpeedInstruction(ins.Head, -1, mixrate))
		}

		mix.Cycles = c
		ret = append(ret, mix)

		if changespeed {
			ret = append(ret, NewSetPipetteSpeedInstruction(ins.Head, -1, defaultpspeed))
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
		if touch_offset > maxTouchOffset {
			touch_offset = maxTouchOffset
		}
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

	if gentlydoesit {
		// reset the drive speed
		ret = append(ret, NewSetDriveSpeedInstruction("Z", pol["DEFAULTZSPEED"].(float64)))
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
	BaseRobotInstruction
	*InstructionType
	Head    int
	Channel int
	Speed   float64
}

func NewSetPipetteSpeedInstruction(head int, channel int, speed float64) *SetPipetteSpeedInstruction {
	v := &SetPipetteSpeedInstruction{
		InstructionType: SPS,
		Head:            head,
		Channel:         channel,
		Speed:           speed,
	}
	v.BaseRobotInstruction = NewBaseRobotInstruction(v)
	return v
}

func (ins *SetPipetteSpeedInstruction) Visit(visitor RobotInstructionVisitor) {
	visitor.SetPipetteSpeed(ins)
}

func (ins *SetPipetteSpeedInstruction) GetParameter(name InstructionParameter) interface{} {
	switch name {
	case HEAD:
		return ins.Head
	case CHANNEL:
		return ins.Channel
	case SPEED:
		return ins.Speed
	default:
		return ins.BaseRobotInstruction.GetParameter(name)
	}
}

func (ins *SetPipetteSpeedInstruction) Generate(ctx context.Context, policy *wtype.LHPolicyRuleSet, prms *LHProperties) ([]RobotInstruction, error) {
	return nil, nil
}

func (ins *SetPipetteSpeedInstruction) OutputTo(lhdriver LiquidhandlingDriver) error {
	driver, ok := lhdriver.(LowLevelLiquidhandlingDriver)
	if !ok {
		return fmt.Errorf("Wrong instruction type for driver: need Lowlevel, got %T", ins)
	}
	ret := driver.SetPipetteSpeed(ins.Head, ins.Channel, ins.Speed)
	if !ret.OK {
		return fmt.Errorf(" %d : %s", ret.Errorcode, ret.Msg)
	}

	return nil

}

type SetDriveSpeedInstruction struct {
	BaseRobotInstruction
	*InstructionType
	Drive string
	Speed float64
}

// NewSetDriveSpeedInstruction set the speed at which the head will move
// drive should be "X", "Y", or "Z", and speed is the speed in mm/s
func NewSetDriveSpeedInstruction(drive string, speed float64) *SetDriveSpeedInstruction {
	v := &SetDriveSpeedInstruction{
		InstructionType: SDS,
		Drive:           drive,
		Speed:           speed,
	}
	v.BaseRobotInstruction = NewBaseRobotInstruction(v)
	return v
}

func (ins *SetDriveSpeedInstruction) Visit(visitor RobotInstructionVisitor) {
	visitor.SetDriveSpeed(ins)
}

func (ins *SetDriveSpeedInstruction) GetParameter(name InstructionParameter) interface{} {
	switch name {
	case DRIVE:
		return ins.Drive
	case SPEED:
		return ins.Speed
	default:
		return ins.BaseRobotInstruction.GetParameter(name)
	}
}

func (ins *SetDriveSpeedInstruction) Generate(ctx context.Context, policy *wtype.LHPolicyRuleSet, prms *LHProperties) ([]RobotInstruction, error) {
	return nil, nil
}

func (ins *SetDriveSpeedInstruction) OutputTo(lhdriver LiquidhandlingDriver) error {
	driver, ok := lhdriver.(LowLevelLiquidhandlingDriver)
	if !ok {
		return fmt.Errorf("Wrong instruction type for driver: need Lowlevel, got %T", ins)
	}
	ret := driver.SetDriveSpeed(ins.Drive, ins.Speed)
	if !ret.OK {
		return fmt.Errorf(" %d : %s", ret.Errorcode, ret.Msg)
	}

	return nil

}

type InitializeInstruction struct {
	BaseRobotInstruction
	*InstructionType
}

func NewInitializeInstruction() *InitializeInstruction {
	v := &InitializeInstruction{
		InstructionType: INI,
	}
	v.BaseRobotInstruction = NewBaseRobotInstruction(v)
	return v
}

func (ins *InitializeInstruction) Visit(visitor RobotInstructionVisitor) {
	visitor.Initialize(ins)
}

func (ins *InitializeInstruction) GetParameter(name InstructionParameter) interface{} {
	return ins.BaseRobotInstruction.GetParameter(name)
}

func (ins *InitializeInstruction) Generate(ctx context.Context, policy *wtype.LHPolicyRuleSet, prms *LHProperties) ([]RobotInstruction, error) {
	return nil, nil
}

func (ins *InitializeInstruction) OutputTo(lhdriver LiquidhandlingDriver) error {
	ret := lhdriver.Initialize()
	if !ret.OK {
		return fmt.Errorf(" %d : %s", ret.Errorcode, ret.Msg)
	}

	return nil

}

type FinalizeInstruction struct {
	BaseRobotInstruction
	*InstructionType
}

func NewFinalizeInstruction() *FinalizeInstruction {
	v := &FinalizeInstruction{
		InstructionType: FIN,
	}
	v.BaseRobotInstruction = NewBaseRobotInstruction(v)
	return v
}

func (ins *FinalizeInstruction) Visit(visitor RobotInstructionVisitor) {
	visitor.Finalize(ins)
}

func (ins *FinalizeInstruction) GetParameter(name InstructionParameter) interface{} {
	return ins.BaseRobotInstruction.GetParameter(name)
}

func (ins *FinalizeInstruction) Generate(ctx context.Context, policy *wtype.LHPolicyRuleSet, prms *LHProperties) ([]RobotInstruction, error) {
	return nil, nil
}

func (ins *FinalizeInstruction) OutputTo(lhdriver LiquidhandlingDriver) error {
	ret := lhdriver.Finalize()
	if !ret.OK {
		return fmt.Errorf(" %d : %s", ret.Errorcode, ret.Msg)
	}

	return nil

}

type WaitInstruction struct {
	BaseRobotInstruction
	*InstructionType
	Time float64
}

func NewWaitInstruction() *WaitInstruction {
	v := &WaitInstruction{
		InstructionType: WAI,
	}
	v.BaseRobotInstruction = NewBaseRobotInstruction(v)
	return v
}

func (ins *WaitInstruction) Visit(visitor RobotInstructionVisitor) {
	visitor.Wait(ins)
}

func (ins *WaitInstruction) GetParameter(name InstructionParameter) interface{} {
	switch name {
	case TIME:
		return ins.Time
	default:
		return ins.BaseRobotInstruction.GetParameter(name)
	}
}

func (ins *WaitInstruction) Generate(ctx context.Context, policy *wtype.LHPolicyRuleSet, prms *LHProperties) ([]RobotInstruction, error) {
	return nil, nil
}

func (ins *WaitInstruction) OutputTo(lhdriver LiquidhandlingDriver) error {
	driver, ok := lhdriver.(LowLevelLiquidhandlingDriver)
	if !ok {
		return fmt.Errorf("Wrong instruction type for driver: need Lowlevel, got %T", ins)
	}
	ret := driver.Wait(ins.Time)
	if !ret.OK {
		return fmt.Errorf(" %d : %s", ret.Errorcode, ret.Msg)
	}

	return nil

}

type LightsOnInstruction struct {
	BaseRobotInstruction
	*InstructionType
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
	v := &LightsOnInstruction{
		InstructionType: LON,
		What:            []string{},
		PltFrom:         []string{},
		PltTo:           []string{},
		WellFrom:        []string{},
		WellTo:          []string{},
		Volume:          []wunit.Volume{},
	}
	v.BaseRobotInstruction = NewBaseRobotInstruction(v)
	return v
}

func (ins *LightsOnInstruction) Visit(visitor RobotInstructionVisitor) {
	visitor.LightsOn(ins)
}

func (ins *LightsOnInstruction) GetParameter(name InstructionParameter) interface{} {
	return ins.BaseRobotInstruction.GetParameter(name)
}

func (ins *LightsOnInstruction) Generate(ctx context.Context, policy *wtype.LHPolicyRuleSet, prms *LHProperties) ([]RobotInstruction, error) {
	return nil, nil
}

func (ins *LightsOnInstruction) OutputTo(lhdriver LiquidhandlingDriver) error {
	/*
		driver, ok := lhdriver.(LowLevelLiquidhandlingDriver)
		if !ok {
			return fmt.Errorf("Wrong instruction type for driver: need Lowlevel, got %T", ins)
		}
	*/
	return fmt.Errorf(" %d : %s", anthadriver.NIM, "Not yet implemented: LightsOnInstruction")
}

type LightsOffInstruction struct {
	BaseRobotInstruction
	*InstructionType
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
	v := &LightsOffInstruction{
		InstructionType: LOF,
		What:            []string{},
		PltFrom:         []string{},
		PltTo:           []string{},
		WellFrom:        []string{},
		WellTo:          []string{},
		Volume:          []wunit.Volume{},
	}
	v.BaseRobotInstruction = NewBaseRobotInstruction(v)
	return v
}

func (ins *LightsOffInstruction) Visit(visitor RobotInstructionVisitor) {
	visitor.LightsOff(ins)
}

func (ins *LightsOffInstruction) GetParameter(name InstructionParameter) interface{} {
	return ins.BaseRobotInstruction.GetParameter(name)
}

func (ins *LightsOffInstruction) Generate(ctx context.Context, policy *wtype.LHPolicyRuleSet, prms *LHProperties) ([]RobotInstruction, error) {
	return nil, nil
}

func (ins *LightsOffInstruction) OutputTo(lhdriver LiquidhandlingDriver) error {
	/*
		driver, ok := lhdriver.(LowLevelLiquidhandlingDriver)
		if !ok {
			return fmt.Errorf("Wrong instruction type for driver: need Lowlevel, got %T", ins)
		}
	*/
	return fmt.Errorf(" %d : %s", anthadriver.NIM, "Not yet implemented: LightsOffInstruction")
}

type OpenInstruction struct {
	BaseRobotInstruction
	*InstructionType
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
	v := &OpenInstruction{
		InstructionType: OPN,
		What:            []string{},
		PltFrom:         []string{},
		PltTo:           []string{},
		WellFrom:        []string{},
		WellTo:          []string{},
		Volume:          []wunit.Volume{},
	}
	v.BaseRobotInstruction = NewBaseRobotInstruction(v)
	return v
}

func (ins *OpenInstruction) Visit(visitor RobotInstructionVisitor) {
	visitor.Open(ins)
}

func (ins *OpenInstruction) GetParameter(name InstructionParameter) interface{} {
	return ins.BaseRobotInstruction.GetParameter(name)
}

func (ins *OpenInstruction) Generate(ctx context.Context, policy *wtype.LHPolicyRuleSet, prms *LHProperties) ([]RobotInstruction, error) {
	return nil, nil
}

func (ins *OpenInstruction) OutputTo(lhdriver LiquidhandlingDriver) error {
	/*
		driver, ok := lhdriver.(LowLevelLiquidhandlingDriver)
		if !ok {
			return fmt.Errorf("Wrong instruction type for driver: need Lowlevel, got %T", ins)
		}
	*/
	return fmt.Errorf(" %d : %s", anthadriver.NIM, "Not yet implemented: OpenInstruction")
}

type CloseInstruction struct {
	BaseRobotInstruction
	*InstructionType
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
	v := &CloseInstruction{
		InstructionType: CLS,
		What:            []string{},
		PltFrom:         []string{},
		PltTo:           []string{},
		WellFrom:        []string{},
		WellTo:          []string{},
		Volume:          []wunit.Volume{},
	}
	v.BaseRobotInstruction = NewBaseRobotInstruction(v)
	return v
}

func (ins *CloseInstruction) Visit(visitor RobotInstructionVisitor) {
	visitor.Close(ins)
}

func (ins *CloseInstruction) GetParameter(name InstructionParameter) interface{} {
	return ins.BaseRobotInstruction.GetParameter(name)
}

func (ins *CloseInstruction) Generate(ctx context.Context, policy *wtype.LHPolicyRuleSet, prms *LHProperties) ([]RobotInstruction, error) {
	return nil, nil
}

func (ins *CloseInstruction) OutputTo(lhdriver LiquidhandlingDriver) error {
	/*
		driver, ok := lhdriver.(LowLevelLiquidhandlingDriver)
		if !ok {
			return fmt.Errorf("Wrong instruction type for driver: need Lowlevel, got %T", ins)
		}
	*/
	return fmt.Errorf(" %d : %s", anthadriver.NIM, "Not yet implemented: CloseInstruction")
}

type LoadAdaptorInstruction struct {
	BaseRobotInstruction
	*InstructionType
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
	v := &LoadAdaptorInstruction{
		InstructionType: LAD,
		What:            []string{},
		PltFrom:         []string{},
		PltTo:           []string{},
		WellFrom:        []string{},
		WellTo:          []string{},
		Volume:          []wunit.Volume{},
	}
	v.BaseRobotInstruction = NewBaseRobotInstruction(v)
	return v
}

func (ins *LoadAdaptorInstruction) Visit(visitor RobotInstructionVisitor) {
	visitor.LoadAdaptor(ins)
}

func (ins *LoadAdaptorInstruction) GetParameter(name InstructionParameter) interface{} {
	return ins.BaseRobotInstruction.GetParameter(name)
}

func (ins *LoadAdaptorInstruction) Generate(ctx context.Context, policy *wtype.LHPolicyRuleSet, prms *LHProperties) ([]RobotInstruction, error) {
	return nil, nil
}

func (ins *LoadAdaptorInstruction) OutputTo(lhdriver LiquidhandlingDriver) error {
	/*
		driver, ok := lhdriver.(LowLevelLiquidhandlingDriver)
		if !ok {
			return fmt.Errorf("Wrong instruction type for driver: need Lowlevel, got %T", ins)
		}
	*/
	return fmt.Errorf(" %d : %s", anthadriver.NIM, "Not yet implemented: LoadAdaptor")
}

type UnloadAdaptorInstruction struct {
	BaseRobotInstruction
	*InstructionType
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
	v := &UnloadAdaptorInstruction{
		InstructionType: UAD,
		What:            []string{},
		PltFrom:         []string{},
		PltTo:           []string{},
		WellFrom:        []string{},
		WellTo:          []string{},
		Volume:          []wunit.Volume{},
	}
	v.BaseRobotInstruction = NewBaseRobotInstruction(v)
	return v
}

func (ins *UnloadAdaptorInstruction) Visit(visitor RobotInstructionVisitor) {
	visitor.UnloadAdaptor(ins)
}

func (ins *UnloadAdaptorInstruction) GetParameter(name InstructionParameter) interface{} {
	return ins.BaseRobotInstruction.GetParameter(name)
}

func (ins *UnloadAdaptorInstruction) Generate(ctx context.Context, policy *wtype.LHPolicyRuleSet, prms *LHProperties) ([]RobotInstruction, error) {
	return nil, nil
}

func (ins *UnloadAdaptorInstruction) OutputTo(lhdriver LiquidhandlingDriver) error {
	/*
		driver, ok := lhdriver.(LowLevelLiquidhandlingDriver)
		if !ok {
			return fmt.Errorf("Wrong instruction type for driver: need Lowlevel, got %T", ins)
		}
	*/
	return fmt.Errorf(" %d : %s", anthadriver.NIM, "Not yet implemented: UnloadAdaptor")
}

type ResetInstruction struct {
	BaseRobotInstruction
	*InstructionType
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
	v := &ResetInstruction{
		InstructionType: RST,
		What:            []string{},
		PltFrom:         []string{},
		WellFrom:        []string{},
		WellTo:          []string{},
		Volume:          []wunit.Volume{},
		FPlateType:      []string{},
		TPlateType:      []string{},
		FVolume:         []wunit.Volume{},
		TVolume:         []wunit.Volume{},
	}
	v.BaseRobotInstruction = NewBaseRobotInstruction(v)
	return v
}

func (ins *ResetInstruction) Visit(visitor RobotInstructionVisitor) {
	visitor.Reset(ins)
}

func (ins *ResetInstruction) GetParameter(name InstructionParameter) interface{} {
	switch name {
	case LIQUIDCLASS:
		return ins.What
	case VOLUME:
		return ins.Volume
	case VOLUNT:
		return nil
	case FROMPLATETYPE:
		return ins.FPlateType
	case WELLFROMVOLUME:
		return ins.FVolume
	case POSFROM:
		return ins.PltFrom
	case POSTO:
		return ins.PltTo
	case WELLFROM:
		return ins.WellFrom
	case WELLTO:
		return ins.WellTo
	case PARAMS:
		return ins.Prms
	case PLATFORM:
		if ins.Prms == nil {
			return ""
		}
		return ins.Prms.Platform
	default:
		return ins.BaseRobotInstruction.GetParameter(name)
	}
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
	for _, tp := range mtp.Transfers {
		ins.AddTransferParams(tp)
	}
}

func (ins *ResetInstruction) Generate(ctx context.Context, policy *wtype.LHPolicyRuleSet, prms *LHProperties) ([]RobotInstruction, error) {
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

	mov := NewMoveInstruction()
	mov.Well = ins.WellTo
	mov.Pos = ins.PltTo
	mov.Plt = ins.TPlateType
	mov.WVolume = ins.TVolume
	mov.Head = ins.Prms.Head
	for i := 0; i < len(mov.Pos); i++ {
		mov.Reference = append(mov.Reference, pol["BLOWOUTREFERENCE"].(int))
		mov.OffsetX = append(mov.OffsetX, 0.0)
		mov.OffsetY = append(mov.OffsetY, 0.0)
		mov.OffsetZ = append(mov.OffsetZ, pol["BLOWOUTOFFSET"].(float64))
	}

	blow := NewBlowoutInstruction()

	blow.Head = ins.Prms.Head
	bov := wunit.NewVolume(pol["BLOWOUTVOLUME"].(float64), pol["BLOWOUTVOLUMEUNIT"].(string))
	blow.Multi = getMulti(ins.What)
	for i := 0; i < blow.Multi; i++ {
		blow.Volume = append(blow.Volume, bov)
	}
	blow.Plt = ins.TPlateType
	blow.What = ins.What

	//no LLF for ResetInstructions
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
	BaseRobotInstruction
	*InstructionType
	Head      int
	Plt       []string
	Well      []string
	Volume    []wunit.Volume // volume of sample being transferred
	PlateType []string
	FVolume   []wunit.Volume // Total volume of sample in the well which the sample is being mixed into?
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
	v := &MoveMixInstruction{
		InstructionType: MMX,
		Plt:             []string{},
		Well:            []string{},
		Volume:          []wunit.Volume{},
		FVolume:         []wunit.Volume{},
		PlateType:       []string{},
		Cycles:          []int{},
		Prms:            map[string]interface{}{},
		What:            []string{},
		Blowout:         []bool{},
		OffsetX:         []float64{},
		OffsetY:         []float64{},
		OffsetZ:         []float64{},
	}
	v.BaseRobotInstruction = NewBaseRobotInstruction(v)
	return v
}

func (ins *MoveMixInstruction) Visit(visitor RobotInstructionVisitor) {
	visitor.MoveMix(ins)
}

func (ins *MoveMixInstruction) GetParameter(name InstructionParameter) interface{} {
	switch name {
	case VOLUME:
		return ins.Volume
	case VOLUNT:
		return nil
	case PLATETYPE:
		return ins.PlateType
	case WELLVOLUME:
		return ins.FVolume
	case POS:
		return ins.Plt
	case WELL:
		return ins.Well
	case PARAMS:
		return ins.Prms
	case CYCLES:
		return ins.Cycles
	case WHAT:
		return ins.What
	case BLOWOUT:
		return ins.Blowout
	case OFFSETX:
		return ins.OffsetX
	case OFFSETY:
		return ins.OffsetY
	case OFFSETZ:
		return ins.OffsetZ
	default:
		return ins.BaseRobotInstruction.GetParameter(name)
	}
}

func (ins *MoveMixInstruction) Generate(ctx context.Context, policy *wtype.LHPolicyRuleSet, prms *LHProperties) ([]RobotInstruction, error) {
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
	BaseRobotInstruction
	*InstructionType
	Head      int
	Volume    []wunit.Volume
	PlateType []string
	What      []string
	Blowout   []bool
	Multi     int
	Cycles    []int
}

func NewMixInstruction() *MixInstruction {
	mi := &MixInstruction{
		InstructionType: MIX,
		Volume:          []wunit.Volume{},
		PlateType:       []string{},
		Cycles:          []int{},
		What:            []string{},
		Blowout:         []bool{},
	}
	mi.BaseRobotInstruction = NewBaseRobotInstruction(mi)
	return mi
}

func (ins *MixInstruction) Visit(visitor RobotInstructionVisitor) {
	visitor.Mix(ins)
}

func (ins *MixInstruction) Generate(ctx context.Context, policy *wtype.LHPolicyRuleSet, prms *LHProperties) ([]RobotInstruction, error) {
	return nil, nil
}

func (ins *MixInstruction) GetParameter(name InstructionParameter) interface{} {
	switch name {
	case VOLUME:
		return ins.Volume
	case PLATETYPE:
		return ins.PlateType
	case CYCLES:
		return ins.Cycles
	case INSTRUCTIONTYPE:
		return ins.Type()
	case LIQUIDCLASS:
		return ins.What
	default:
		return nil
	}
}

func (mi *MixInstruction) OutputTo(lhdriver LiquidhandlingDriver) error {
	driver, ok := lhdriver.(LowLevelLiquidhandlingDriver)

	if !ok {
		return fmt.Errorf("Wrong instruction type for driver: need Lowlevel, got %T", mi)
	}
	vols := make([]float64, len(mi.Volume))

	for i := 0; i < len(mi.Volume); i++ {
		vols[i] = mi.Volume[i].ConvertToString("ul")
	}

	ret := driver.Mix(mi.Head, vols, mi.PlateType, mi.Cycles, mi.Multi, mi.What, mi.Blowout)

	if !ret.OK {
		return fmt.Errorf(" %d : %s", ret.Errorcode, ret.Msg)
	}

	return nil

}

func countMulti(sa []string) int {
	r := 0
	for _, s := range sa {
		if s != "" {
			r += 1
		}
	}

	return r
}

func getFirstDefined(sa []string) int {
	x := -1

	for i := 0; i < len(sa); i++ {
		if sa[i] != "" {
			x = i
			break
		}
	}
	return x
}

func GetTips(ctx context.Context, tiptypes []string, params *LHProperties, channel []*wtype.LHChannelParameter, usetiptracking bool) ([]RobotInstruction, error) {
	// GetCleanTips returns enough sets of tip boxes to get all distinct tip types
	tipwells, tipboxpositions, tipboxtypes, terr := params.GetCleanTips(ctx, tiptypes, channel, usetiptracking)

	if tipwells == nil || terr != nil {
		err := wtype.LHError(wtype.LH_ERR_NO_TIPS, fmt.Sprintf("PICKUP: types: %v On Deck: %v", tiptypes, params.GetLayout()))
		return []RobotInstruction{NewLoadTipsMoveInstruction()}, err
	}

	inss := make([]RobotInstruction, 0, 1)

	for i := 0; i < len(tipwells); i++ {
		// all instructions in a block must have a head in common
		defPos := getFirstDefined(tipwells[i])

		if defPos == -1 {
			return inss, fmt.Errorf("Error: tip get failed for types %v", tiptypes)
		}

		ins := NewLoadTipsMoveInstruction()
		ins.Head = channel[defPos].Head
		ins.Well = tipwells[i]
		ins.FPosition = tipboxpositions[i]
		ins.FPlateType = tipboxtypes[i]
		ins.Multi = countMulti(tipwells[i])

		inss = append(inss, ins)
	}

	return inss, nil
}

func collate(s []string) string {
	m := make(map[string]int, len(s))
	for _, v := range s {
		m[v] += 1
	}

	r := ""

	for k, v := range m {
		r += fmt.Sprintf("%d %s, ", v, k)
	}

	return r
}

//func DropTips(tiptype string, params *LHProperties, channel *wtype.LHChannelParameter, multi int) (RobotInstruction, error) {
func DropTips(tiptypes []string, params *LHProperties, channels []*wtype.LHChannelParameter) (RobotInstruction, error) {
	tipwells, tipwastepositions, tipwastetypes := params.DropDirtyTips(channels)

	if tipwells == nil {
		ins := NewUnloadTipsMoveInstruction()
		err := wtype.LHError(wtype.LH_ERR_TIP_WASTE, collate(tiptypes))
		return ins, err
	}

	defpos := getFirstDefined(tipwells)

	if defpos == -1 {
		return NewUnloadTipsMoveInstruction(), wtype.LHError(wtype.LH_ERR_NO_TIPS, fmt.Sprint("DROP: type ", tiptypes))
	}

	ins := NewUnloadTipsMoveInstruction()
	ins.Head = channels[defpos].Head
	ins.WellTo = tipwells
	ins.PltTo = tipwastepositions
	ins.TPlateType = tipwastetypes
	ins.Multi = getMulti(tiptypes)
	return ins, nil
}

func getMulti(w []string) int {
	c := 0
	for _, v := range w {
		if v != "" {
			c += 1
		}
	}

	return c
}

func getUseLLF(pol wtype.LHPolicy, multi int, plates []string, prms *LHProperties) ([]bool, bool) {
	use_llf := make([]bool, multi)
	any_llf := false
	enable_llf := SafeGetBool(pol, "USE_LLF")

	//save a few ms
	if !enable_llf {
		return use_llf, enable_llf
	}

	for i := 0; i < multi; i++ {
		//probably just fetching the same plate each time
		plate := prms.Plates[plates[i]]

		//do LLF if the well has a volumemodel
		use_llf[i] = enable_llf && plate.Welltype.HasLiquidLevelModel()

		any_llf = any_llf || use_llf[i]
	}

	return use_llf, any_llf
}

// compare proposed value to minimum and maximum tolerated
// return proposed if within bounds
// return relevant bound (min or max) if proposed is outside the range and overrideIfOutOfRange is true
// return an error otherwise
func checkAndSaften(proposed, min, max float64, overrideIfOutOfRange bool) (float64, error) {
	if proposed < min {
		if !overrideIfOutOfRange {
			return proposed, fmt.Errorf("value %f out of range %f - %f", proposed, min, max)
		} else {
			return min, nil
		}
	} else if proposed > max {
		if !overrideIfOutOfRange {
			return proposed, fmt.Errorf("value %f out of range %f - %f", proposed, min, max)
		} else {
			return max, nil
		}

	}

	return proposed, nil
}

//makeZOffsetSafe increase the zoffset to prevent the robot head from colliding
//with the top of the plate when accessing the bottom of particularly deep wells
//with shorter tips.
//Does not affect behaviour with troughs and other wells that are big enough for
//the entire head to fit inside.
func makeZOffsetSafe(ctx context.Context, prms *LHProperties, zoffset float64, headIndex int, plates []string, tiptype string) (float64, error) {
	platename := ""
	for _, p := range plates {
		if p != "" {
			platename = p
			break
		}
	}

	plate := prms.Plates[platename]

	//get the size of all the channels together
	adaptor := prms.GetLoadedAdaptor(headIndex)
	channelSpacing := 9.0 //get this from adaptor in future
	coneDiameter := 5.5   //get this from adaptor in future
	adaptorSize := wtype.Coordinates{X: coneDiameter, Y: coneDiameter}
	adaptorWidth := channelSpacing*float64(adaptor.Params.Multi-1) + coneDiameter
	if adaptor.Params.Orientation == wtype.LHVChannel {
		adaptorSize.Y = adaptorWidth
	} else {
		adaptorSize.X = adaptorWidth
	}

	//if all the channels can fit in the well, don't add offset
	//this means we can still reach the bottom of troughs and reservoirs
	if s := plate.Welltype.GetSize(); s.X > adaptorSize.X && s.Y > adaptorSize.Y {
		return zoffset, nil
	}

	tipbox, err := inventory.NewTipbox(ctx, tiptype)
	if err != nil {
		return zoffset, err
	}

	//safetyZHeight is a small offset to avoid predicted collisions due to numerical error
	minZ := plate.Welltype.GetSize().Z - tipbox.Tiptype.GetEffectiveHeight() - plate.Welltype.Bottomh + safetyZHeight

	if minZ > zoffset {
		return minZ, nil
	}

	return zoffset, nil
}
