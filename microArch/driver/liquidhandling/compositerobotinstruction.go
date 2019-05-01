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
	"strings"

	"github.com/pkg/errors"

	"reflect"

	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/antha/anthalib/wutil/text"
	"github.com/antha-lang/antha/inventory"
	anthadriver "github.com/antha-lang/antha/microArch/driver"
)

// Valid parameter fields for robot instructions
const (
	//maxTouchOffset maximum value for TOUCHOFFSET which makes sense - values larger than this are capped to this value
	maxTouchOffset = 5.0

	//added to avoid floating point issues with heights in simulator
	safetyZHeight = 0.05
)

type ChannelBlockInstruction struct {
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
	Component  [][]string                    // array of component name (i.e. Liquid's CName) by [transfer][channel]
	Prms       [][]*wtype.LHChannelParameter // which channel properties apply to each transfer
	Multi      []int
}

func NewChannelBlockInstruction() *ChannelBlockInstruction {
	v := &ChannelBlockInstruction{
		InstructionType: CBI,
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
		Component:       [][]string{},
		Prms:            [][]*wtype.LHChannelParameter{},
		Multi:           []int{},
	}
	v.BaseRobotInstruction = NewBaseRobotInstruction(v)
	return v
}

func (ins *ChannelBlockInstruction) Visit(visitor RobotInstructionVisitor) {
	visitor.ChannelBlock(ins)
}

func (ins *ChannelBlockInstruction) AddTransferParams(mct MultiTransferParams) {
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
	ins.Component = append(ins.Component, mct.Component())
	ins.Prms = append(ins.Prms, mct.Channels())
	ins.Multi = append(ins.Multi, mct.Multi)
}

func (ins *ChannelBlockInstruction) GetParameter(name InstructionParameter) interface{} {
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
		ret := make([]string, 0, len(ins.Prms))
		for _, channelParams := range ins.Prms {
			for _, p := range channelParams {
				if p != nil {
					ret = append(ret, p.Platform)
				}
			}
		}
		return ret
	case WELLTO:
		return ins.WellTo
	case WELLTOVOLUME:
		return ins.TVolume
	case TOPLATETYPE:
		return ins.TPlateType
	case COMPONENT:
		return ins.Component
	case MULTI:
		return ins.Multi
	default:
		return ins.BaseRobotInstruction.GetParameter(name)
	}
}

func (ins *ChannelBlockInstruction) GetVolumes() []wunit.Volume {
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

func (ins *ChannelBlockInstruction) MaxMulti() int {
	mx := 0
	for _, m := range ins.Multi {
		if m > mx {
			mx = m
		}
	}
	return mx
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

// By the point at which the ChannelBlockInstruction is used by the Generate method all transfers will share the same policy.
func (ins *ChannelBlockInstruction) Generate(ctx context.Context, policy *wtype.LHPolicyRuleSet, prms *LHProperties) ([]RobotInstruction, error) {

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

	// variables for tracking tip state
	usetiptracking := SafeGetBool(policy.Options, "USE_DRIVER_TIP_TRACKING")
	tipUseCounter := 0
	changeTips := true // always load tips to start with
	var lastThing *wtype.Liquid
	var channels []*wtype.LHChannelParameter
	var tiptypes []string

	for t := 0; t < len(ins.Volume); t++ {
		if len(ins.What[t]) == 0 {
			continue
		}

		prmSet := ins.Prms[t][0]

		tvols := NewVolumeSet(prmSet.Multi)
		fvols := NewVolumeSet(prmSet.Multi)
		for i := range ins.Volume[t] {
			fvols[i] = wunit.CopyVolume(ins.FVolume[t][i])
			tvols[i] = wunit.CopyVolume(ins.TVolume[t][i])
		}

		// choose which tips should be used for this transfer
		newchannels, newtips, newtiptypes, err := ChooseChannels(ins.Volume[t], prms)
		if err != nil {
			return ret, err
		}

		// due to robot limitations the chosen has to be the same for each channel
		// this is safe for non-independent robots since the volumes will alwayws be the same
		// hjk TODO: we should change ChooseChannels such that this is always the case,
		//           which also significantly simplifies the behaviour of GetTips, see ANTHA-2648
		types := make(map[string]bool, len(newtiptypes))
		for _, ttype := range newtiptypes {
			if ttype != "" {
				types[ttype] = true
			}
		}
		if len(types) > 1 {
			tiptypes := make([]string, 0, len(types))
			for ttype := range types {
				tiptypes = append(tiptypes, ttype)
			}
			return nil, errors.Errorf("tip types must be the same: cannot load tips %s at the same time", strings.Join(tiptypes, " and "))
		}

		// split the transfer up
		// volumes no longer equal
		tvs, err := TransferVolumesMulti(VolumeSet(ins.Volume[t]), mergeTipsAndChannels(newchannels, newtips))
		if err != nil {
			return ret, err
		}

		for _, vols := range tvs {
			// determine whether to change tips
			changeTips = changeTips || tipUseCounter > pol["TIP_REUSE_LIMIT"].(int)
			changeTips = changeTips || !reflect.DeepEqual(channels, newchannels)
			changeTips = changeTips || !reflect.DeepEqual(tiptypes, newtiptypes)

			// big dangerous assumption here: we need to check if anything is different
			thisThing := prms.Plates[ins.PltFrom[t][0]].Wellcoords[ins.WellFrom[t][0]].Contents()

			if lastThing != nil {
				if thisThing.CName != lastThing.CName {
					changeTips = true
				}
			}

			if changeTips {
				// drop the last tips if there are any loaded
				if tiptypes != nil && channels != nil {
					if tipdrp, err := DropTips(tiptypes, prms, channels); err != nil {
						return ret, err
					} else {
						ret = append(ret, tipdrp)
					}
				}

				if tipget, err := GetTips(ctx, newtiptypes, prms, newchannels, usetiptracking); err != nil {
					return ret, err
				} else {
					ret = append(ret, tipget...)
				}

				tipUseCounter = 0
				lastThing = nil
				changeTips = false
				tiptypes = newtiptypes
				channels = newchannels
			}

			mci := NewChannelTransferInstruction()
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
			mci.Component = ins.Component[t]
			mci.TipType = tiptypes
			mci.Multi = countMulti(ins.PltFrom[t])
			channelprms := make([]*wtype.LHChannelParameter, newchannels[0].Multi)
			for i := 0; i < len(newchannels); i++ {
				if newchannels[i] != nil {
					channelprms[i] = newchannels[i].MergeWithTip(newtips[i])
				}
			}
			mci.Prms = channelprms

			ret = append(ret, mci)

			tipUseCounter++
			lastThing = thisThing

			// check if we are touching a bad liquid
			// in future we will do this properly, for now we assume
			// touching any liquid is bad
			npre, premix := pol["PRE_MIX"]
			npost, postmix := pol["POST_MIX"]
			if pol["DSPREFERENCE"].(int) == 0 && !VolumeSet(ins.TVolume[t]).IsZero() || premix && npre.(int) > 0 || postmix && npost.(int) > 0 {
				changeTips = true
			}

			// update the volumes yet to transfer
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

type ChannelTransferInstruction struct {
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
	Component  []string
}

func (scti *ChannelTransferInstruction) Params(k int) TransferParams {
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
	tp.Component = scti.Component[k]
	return tp
}

// Channels return the channel indexes of each channel used in the instruction
func (cti *ChannelTransferInstruction) Channels() []int {
	ret := make([]int, 0, len(cti.Volume))
	for i, v := range cti.Volume {
		if !v.IsZero() {
			ret = append(ret, i)
		}
	}
	return ret
}

func NewChannelTransferInstruction() *ChannelTransferInstruction {
	v := &ChannelTransferInstruction{
		InstructionType: CTI,
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
		Component:       []string{},
	}
	v.BaseRobotInstruction = NewBaseRobotInstruction(v)
	return v
}

func (ins *ChannelTransferInstruction) Visit(visitor RobotInstructionVisitor) {
	visitor.ChannelTransfer(ins)
}

func (ins *ChannelTransferInstruction) GetParameter(name InstructionParameter) interface{} {
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
	case COMPONENT:
		return ins.Component
	default:
		return ins.BaseRobotInstruction.GetParameter(name)
	}
}

func (ins *ChannelTransferInstruction) Generate(ctx context.Context, policy *wtype.LHPolicyRuleSet, prms *LHProperties) ([]RobotInstruction, error) {
	return []RobotInstruction{NewSuckInstruction(ins), NewBlowInstruction(ins)}, nil
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
	Component  []string
}

func NewAspirateInstruction() *AspirateInstruction {
	v := &AspirateInstruction{
		InstructionType: ASP,
		Volume:          []wunit.Volume{},
		Plt:             []string{},
		What:            []string{},
		LLF:             []bool{},
		Component:       []string{},
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
	case COMPONENT:
		return ins.Component
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

	return driver.Aspirate(volumes, os, ins.Head, ins.Multi, ins.Plt, ins.What, ins.LLF).GetError()
}

type DispenseInstruction struct {
	BaseRobotInstruction
	*InstructionType
	Head      int
	Volume    []wunit.Volume
	Multi     int
	Plt       []string
	What      []string
	LLF       []bool
	Platform  string
	Component []string
}

func NewDispenseInstruction() *DispenseInstruction {
	v := &DispenseInstruction{
		InstructionType: DSP,
		Volume:          []wunit.Volume{},
		Plt:             []string{},
		What:            []string{},
		LLF:             []bool{},
		Component:       []string{},
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
	case COMPONENT:
		return ins.Component
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
	return driver.Dispense(volumes, os, ins.Head, ins.Multi, ins.Plt, ins.What, ins.LLF).GetError()
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
	var nonZero bool
	for i := 0; i < ins.Multi; i++ {
		if volumes[i] > 0.0 {
			bo[i] = true
			nonZero = true
		}
	}
	if !nonZero {
		return nil
	}
	return driver.Dispense(volumes, bo, ins.Head, ins.Multi, ins.Plt, ins.What, ins.LLF).GetError()
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

	return driver.ResetPistons(ins.Head, ins.Channel).GetError()
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
	return driver.Move(ins.Pos, ins.Well, ins.Reference, ins.OffsetX, ins.OffsetY, ins.OffsetZ, ins.Plt, ins.Head).GetError()
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
	if driver, ok := lhdriver.(LowLevelLiquidhandlingDriver); !ok {
		return fmt.Errorf("Wrong instruction type for driver: need Lowlevel, got %T", ins)
	} else {
		return driver.LoadTips(ins.Channels, ins.Head, ins.Multi, ins.HolderType, ins.Pos, ins.Well).GetError()
	}
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
	if driver, ok := lhdriver.(LowLevelLiquidhandlingDriver); !ok {
		return fmt.Errorf("Wrong instruction type for driver: need Lowlevel, got %T", ins)
	} else {
		return driver.UnloadTips(ins.Channels, ins.Head, ins.Multi, ins.HolderType, ins.Pos, ins.Well).GetError()
	}
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
	Component   []string
}

func NewSuckInstruction(cti *ChannelTransferInstruction) *SuckInstruction {
	// channels parameters must be the same for each channel, i.e. that the same tip and head was chosen for each
	var prms *wtype.LHChannelParameter
	for _, cp := range cti.Prms {
		if cp != nil {
			if prms != nil && !prms.Equals(cp) {
				panic("ChannelTransferInstruction mixes channel parameters")
			}
			prms = cp
		}
	}
	var tipType string
	for _, tt := range cti.TipType {
		if tt != "" {
			if tipType != "" && tipType != tt {
				panic("ChannelTransferInstruction mixes tip types")
			}
			tipType = tt
		}
	}

	ret := &SuckInstruction{
		InstructionType: SUK,
		What:            make([]string, len(cti.What)),
		PltFrom:         make([]string, len(cti.PltFrom)),
		WellFrom:        make([]string, len(cti.WellFrom)),
		Volume:          make([]wunit.Volume, len(cti.Volume)),
		FPlateType:      make([]string, len(cti.FPlateType)),
		FVolume:         make([]wunit.Volume, len(cti.FVolume)),
		Component:       make([]string, len(cti.Component)),
		Prms:            prms.DupKeepIDs(),
		Head:            prms.Head,
		Multi:           cti.Multi,
		TipType:         tipType,
	}
	ret.BaseRobotInstruction = NewBaseRobotInstruction(ret)

	copy(ret.What, cti.What)
	copy(ret.PltFrom, cti.PltFrom)
	copy(ret.WellFrom, cti.WellFrom)
	copy(ret.Volume, cti.Volume)
	copy(ret.FPlateType, cti.FPlateType)
	copy(ret.FVolume, cti.FVolume)
	copy(ret.Component, cti.Component)

	return ret
}

func (ins *SuckInstruction) Visit(visitor RobotInstructionVisitor) {
	visitor.Suck(ins)
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
	ofz := SafeGetF64(pol, "ASPZOFFSET")
	ofzadj := SafeGetF64(pol, "OFFSETZADJUST")
	ofz += ofzadj
	ofz, err = makeZOffsetSafe(ctx, prms, ofz, ins.Head, ins.PltFrom, ins.TipType)
	if err != nil {
		return nil, err
	}

	mixofx := SafeGetF64(pol, "PRE_MIX_X")
	mixofy := SafeGetF64(pol, "PRE_MIX_Y")
	mixofz := SafeGetF64(pol, "PRE_MIX_Z")
	mixofz += ofzadj
	final_asp_ref := SafeGetInt(pol, "ASPREFERENCE")
	mixofz, err = makeZOffsetSafe(ctx, prms, mixofz, ins.Head, ins.PltFrom, ins.TipType)
	if err != nil {
		return nil, err
	}

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
		for i := 0; i < len(ins.What); i++ {
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

		// if not set we use the instruction value

		// XXX -- only looking at first vol specified
		if mixvol == 0.0 {
			mixvol = ins.Volume[0].ConvertToString("ul")
		}

		vmixvol := wunit.NewVolume(mixvol, "ul")

		// TODO -- corresponding checks when set
		if mixvol < wtype.Globals.MIN_REASONABLE_VOLUME_UL {
			return ret, wtype.LHError(wtype.LH_ERR_POLICY, fmt.Sprintf("PRE_MIX_VOLUME set below minimum allowed: %f min %f", mixvol, wtype.Globals.MIN_REASONABLE_VOLUME_UL))
		} else if !ins.Prms.CanMove(vmixvol, true) {
			override := SafeGetBool(pol, "MIX_VOLUME_OVERRIDE_TIP_MAX")
			if override {
				mixvol = ins.Prms.Maxvol.ConvertToString("ul")
			} else {
				// this is an error in channel choice but the user has to deal... needs modificationst
				return ret, wtype.LHError(wtype.LH_ERR_POLICY, fmt.Sprintf("PRE_MIX_VOLUME not compatible with optimal channel choice: requested %s channel limits are %s", vmixvol.ToString(), ins.Prms.VolumeLimitString()))
			}
		}

		if ok {
			v := make([]wunit.Volume, len(ins.What))
			for i, what := range ins.What {
				if what == "" {
					v[i] = wunit.ZeroVolume()
				} else {
					v[i] = wunit.NewVolume(mixvol, "ul")
				}
			}
			mix.Volume = v
		}
		// offsets

		for k := 0; k < len(ins.What); k++ {
			mix.OffsetX = append(mix.OffsetX, mixofx)
			mix.OffsetY = append(mix.OffsetY, mixofy)
			mix.OffsetZ = append(mix.OffsetZ, mixofz)
			mix.Cycles = append(mix.Cycles, cycles)
		}

		// set speed

		//_, changepipspeed := pol["PRE_MIX_RATE"]

		mixrate := SafeGetF64(pol, "PRE_MIX_RATE")

		changepipspeed := (mixrate != defaultpspeed) && (mixrate > 0.0)

		if changepipspeed {
			mixrate, err = checkAndSaften(mixrate, head.Params.Minspd.RawValue(), head.Params.Maxspd.RawValue(), allowOutOfRangePipetteSpeeds)
			if err != nil {
				return []RobotInstruction{}, errors.Wrap(err, "setting pre mix pipetting speed")
			}

			setspd := NewSetPipetteSpeedInstruction()
			setspd.Head = ins.Head
			setspd.Channel = -1 // all channels
			setspd.Speed = mixrate
			ret = append(ret, setspd)
		}

		ret = append(ret, mix)

		if changepipspeed {
			sps := NewSetPipetteSpeedInstruction()
			sps.Head = ins.Head
			sps.Channel = -1 // all channels
			sps.Speed = defaultpspeed
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

	for i := 0; i < len(ins.What); i++ {
		mov.Reference = append(mov.Reference, final_asp_ref)
		mov.OffsetX = append(mov.OffsetX, ofx)
		mov.OffsetY = append(mov.OffsetY, ofy)
		mov.OffsetZ = append(mov.OffsetZ, ofz)
	}
	ret = append(ret, mov)

	// Set the pipette speed if needed

	apspeed := SafeGetF64(pol, "ASPSPEED")

	changepspeed := (apspeed != defaultpspeed) && (apspeed > 0.0)

	if changepspeed {
		apspeed, err = checkAndSaften(apspeed, head.Params.Minspd.RawValue(), head.Params.Maxspd.RawValue(), allowOutOfRangePipetteSpeeds)

		if err != nil {
			return []RobotInstruction{}, errors.Wrap(err, "setting pipette aspirate speed")
		}
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
		for i := range aspins.Volume {
			aspins.Volume[i].Add(extra_vol)
		}
	}

	aspins.Multi = ins.Multi
	aspins.Overstroke = ins.Overstroke
	aspins.What = ins.What
	aspins.Plt = ins.FPlateType
	aspins.Component = ins.Component

	for i := 0; i < len(aspins.What); i++ {
		// follow the liquidlevel if we moved to it earlier
		aspins.LLF = append(aspins.LLF, final_asp_ref == wtype.LiquidReference.AsInt())
	}

	ret = append(ret, aspins)

	// do we reset the pipette speed?

	if changepspeed {
		sps := NewSetPipetteSpeedInstruction()
		sps.Head = ins.Head
		sps.Channel = -1 // all channels
		sps.Speed = defaultpspeed
		ret = append(ret, sps)
	}

	// do we wait

	_, wait := pol["ASP_WAIT"]

	wait_time := SafeGetF64(pol, "ASP_WAIT")

	if wait && wait_time > 0.0 {
		waitins := NewWaitInstruction()
		waitins.Time = wait_time
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
		for i := 0; i < len(ins.What); i++ {
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
	Component  []string
}

func NewBlowInstruction(cti *ChannelTransferInstruction) *BlowInstruction {
	// we're assuming here that the channels parameters are the same for each channel, i.e. that the same tip and head was chosen for each
	var prms *wtype.LHChannelParameter
	for _, cp := range cti.Prms {
		if cp != nil {
			if prms != nil && !prms.Equals(cp) {
				panic("ChannelTransferInstruction mixes different channel parameters")
			}
			prms = cp
		}
	}
	var tipType string
	for _, tt := range cti.TipType {
		if tt != "" {
			if tipType != "" && tipType != tt {
				panic("ChannelTransferInstruction mixes different tip types")
			}
			tipType = tt
		}
	}

	ret := &BlowInstruction{
		InstructionType: BLW,
		What:            make([]string, len(cti.What)),
		PltTo:           make([]string, len(cti.PltTo)),
		WellTo:          make([]string, len(cti.WellTo)),
		Volume:          make([]wunit.Volume, len(cti.Volume)),
		TPlateType:      make([]string, len(cti.TPlateType)),
		TVolume:         make([]wunit.Volume, len(cti.TVolume)),
		Component:       make([]string, len(cti.Component)),
		Prms:            prms.DupKeepIDs(),
		Head:            prms.Head,
		TipType:         tipType,
		Multi:           cti.Multi,
	}
	ret.BaseRobotInstruction = NewBaseRobotInstruction(ret)

	copy(ret.What, cti.What)
	copy(ret.PltTo, cti.PltTo)
	copy(ret.WellTo, cti.WellTo)
	copy(ret.Volume, cti.Volume)
	copy(ret.TPlateType, cti.TPlateType)
	copy(ret.TVolume, cti.TVolume)
	copy(ret.Component, cti.Component)

	return ret
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
	case COMPONENT:
		return ins.Component
	default:
		return ins.BaseRobotInstruction.GetParameter(name)
	}
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
	ret := make([]RobotInstruction, 0)

	// pipetting speed
	defaultpspeed := SafeGetF64(pol, "DEFAULTPIPETTESPEED")
	setspd := NewSetPipetteSpeedInstruction()
	setspd.Head = head
	setspd.Channel = -1 // all channels
	setspd.Speed = defaultpspeed
	ret = append(ret, setspd)

	// Z move speed
	spd := NewSetDriveSpeedInstruction()
	spd.Drive = "Z"
	spd.Speed = pol["DEFAULTZSPEED"].(float64)
	ret = append(ret, spd)

	return ret
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

	gentlydoesit := entryspeed > 0.0 && entryspeed != defaultspeed

	if gentlydoesit {
		// go to the well top
		mov := NewMoveInstruction()

		mov.Head = ins.Head
		mov.Pos = ins.PltTo
		mov.Plt = ins.TPlateType
		mov.Well = ins.WellTo
		mov.WVolume = ins.TVolume
		for i := 0; i < len(ins.What); i++ {
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
	for i := 0; i < len(ins.What); i++ {
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

		sps := NewSetPipetteSpeedInstruction()
		sps.Head = ins.Head
		sps.Channel = -1 // all channels
		sps.Speed = dpspeed
		ret = append(ret, sps)
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

		for i := 0; i < len(ins.What); i++ {
			// follow the liquid-level if we moved to it earlier
			boins.LLF = append(boins.LLF, ref == wtype.LiquidReference.AsInt())
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
		dspins.Component = ins.Component

		for i := 0; i < len(ins.What); i++ {
			// follow the liquid-level if we moved to it earlier
			dspins.LLF = append(dspins.LLF, ref == wtype.LiquidReference.AsInt())
		}

		ret = append(ret, dspins)
	}

	// do we reset the pipette speed?

	if setpspeed {
		sps := NewSetPipetteSpeedInstruction()
		sps.Head = ins.Head
		sps.Channel = -1 // all channels
		sps.Speed = defaultpspeed
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
		b := make([]bool, len(ins.What))
		mix.Blowout = b

		// offsets

		pmxoff := SafeGetF64(pol, "POST_MIX_X")

		for k := 0; k < len(ins.What); k++ {
			mix.OffsetX = append(mix.OffsetX, pmxoff)
		}

		pmyoff := SafeGetF64(pol, "POST_MIX_Y")
		for k := 0; k < len(ins.What); k++ {
			mix.OffsetY = append(mix.OffsetY, pmyoff)
		}

		pmzoff := SafeGetF64(pol, "POST_MIX_Z")
		pmzoff += ofzadj

		pmzoff, err = makeZOffsetSafe(ctx, prms, pmzoff, ins.Head, ins.PltTo, ins.TipType)
		if err != nil {
			return nil, err
		}

		for k := 0; k < len(ins.What); k++ {
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
			v := make([]wunit.Volume, len(ins.What))
			for i, what := range ins.What {
				if what == "" {
					v[i] = wunit.ZeroVolume()
				} else {
					v[i] = wunit.NewVolume(mixvol, "ul")
				}
			}
			mix.Volume = v
		}

		c := make([]int, len(ins.What))

		for i := 0; i < len(ins.What); i++ {
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
			sps.Speed = defaultpspeed
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
		if touch_offset > maxTouchOffset {
			touch_offset = maxTouchOffset
		}
		mov := NewMoveInstruction()
		mov.Head = ins.Head
		mov.Pos = ins.PltTo
		mov.Plt = ins.TPlateType
		mov.Well = ins.WellTo
		mov.WVolume = ins.TVolume

		ref := make([]int, len(ins.What))
		off := make([]float64, len(ins.What))
		ox := make([]float64, len(ins.What))
		oy := make([]float64, len(ins.What))
		for i := 0; i < len(ins.What); i++ {
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
		spd := NewSetDriveSpeedInstruction()
		spd.Drive = "Z"
		spd.Speed = pol["DEFAULTZSPEED"].(float64)
		ret = append(ret, spd)

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

func NewSetPipetteSpeedInstruction() *SetPipetteSpeedInstruction {
	v := &SetPipetteSpeedInstruction{
		InstructionType: SPS,
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
	if driver, ok := lhdriver.(LowLevelLiquidhandlingDriver); !ok {
		return fmt.Errorf("Wrong instruction type for driver: need Lowlevel, got %T", ins)
	} else {
		return driver.SetPipetteSpeed(ins.Head, ins.Channel, ins.Speed).GetError()
	}
}

type SetDriveSpeedInstruction struct {
	BaseRobotInstruction
	*InstructionType
	Drive string
	Speed float64
}

func NewSetDriveSpeedInstruction() *SetDriveSpeedInstruction {
	v := &SetDriveSpeedInstruction{
		InstructionType: SDS,
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
	if driver, ok := lhdriver.(LowLevelLiquidhandlingDriver); !ok {
		return fmt.Errorf("Wrong instruction type for driver: need Lowlevel, got %T", ins)
	} else {
		return driver.SetDriveSpeed(ins.Drive, ins.Speed).GetError()
	}
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
	return lhdriver.Initialize().GetError()
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
	return lhdriver.Finalize().GetError()
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
	if driver, ok := lhdriver.(LowLevelLiquidhandlingDriver); !ok {
		return fmt.Errorf("Wrong instruction type for driver: need Lowlevel, got %T", ins)
	} else {
		return driver.Wait(ins.Time).GetError()
	}
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
	blow.Volume = make([]wunit.Volume, len(ins.What))
	for i := 0; i < len(ins.What); i++ {
		if ins.What[i] != "" {
			blow.Volume[i] = bov
		}
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

	if bov.RawValue() > 0.0 {
		ret = append(ret, mov)
		ret = append(ret, blow)
	}

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
	mov.Reference = make([]int, len(ins.What))
	// mov.Reference[i] == 0 for all i
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
	if driver, ok := lhdriver.(LowLevelLiquidhandlingDriver); !ok {
		return fmt.Errorf("Wrong instruction type for driver: need Lowlevel, got %T", mi)
	} else {
		vols := make([]float64, len(mi.Volume))
		for i := 0; i < len(mi.Volume); i++ {
			vols[i] = mi.Volume[i].ConvertToString("ul")
		}

		return driver.Mix(mi.Head, vols, mi.PlateType, mi.Cycles, mi.Multi, mi.What, mi.Blowout).GetError()
	}
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
	adaptorSize := wtype.Coordinates3D{X: coneDiameter, Y: coneDiameter}
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

	var tipbox *wtype.LHTipbox
	for _, tb := range prms.Tipboxes {
		if tb.Tiptype.Type == tiptype {
			tipbox = tb
			break
		}
	}
	if tipbox == nil {
		// this can only happen if there's an error in channel selection
		return zoffset, wtype.LHError(wtype.LH_ERR_OTHER, fmt.Sprintf("instruction requested tip type %q but none found in parameters: please report this to the authors", tiptype))
	}

	//safetyZHeight is a small offset to avoid predicted collisions due to numerical error
	minZ := plate.Welltype.GetSize().Z - tipbox.Tiptype.GetEffectiveHeight() - plate.Welltype.Bottomh + safetyZHeight

	if minZ > zoffset {
		return minZ, nil
	}

	return zoffset, nil
}
