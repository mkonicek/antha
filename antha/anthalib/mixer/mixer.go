// mixer/mixer.go: Part of the Antha language
// Copyright (C) 2014 the Antha authors. All rights reserved.
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

// Package mixer deals with mixing and sampling in Antha
package mixer

import (
	"fmt"

	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
)

// SampleAll takes all of this liquid
func SampleAll(l *wtype.Liquid) *wtype.Liquid {
	return Sample(l, l.Volume())
}

// Sample takes a sample of volume v from this liquid
func Sample(l *wtype.Liquid, v wunit.Volume) *wtype.Liquid {
	ret := wtype.NewLHComponent()
	//	ret.ID = l.ID
	l.AddDaughterComponent(ret)
	ret.ParentID = l.ID
	ret.CName = l.Name()
	ret.Type = l.Type
	ret.Vol = v.RawValue()
	ret.Vunit = v.Unit().PrefixedSymbol()
	ret.Extra = l.GetExtra()
	ret.SubComponents = l.SubComponents
	ret.Smax = l.GetSmax()
	ret.Visc = l.GetVisc()
	if l.Conc > 0 && len(l.Cunit) > 0 {
		ret.SetConcentration(wunit.NewConcentration(l.Conc, l.Cunit))
	}

	ret.SetSample(true)

	return ret
}

// SplitSample is a two-return version of sample
func SplitSample(l *wtype.Liquid, v wunit.Volume) (moving, remaining *wtype.Liquid) {
	remaining = l.Dup()

	moving = Sample(remaining, v)

	remaining.Vol -= v.ConvertToString(remaining.Vunit)
	remaining.ID = wtype.GetUUID()

	return
}

// MultiSample takes an array of samples and array of corresponding volumes and
// sample them all
func MultiSample(l []*wtype.Liquid, v []wunit.Volume) []*wtype.Liquid {
	reta := make([]*wtype.Liquid, 0)

	for i, j := range l {
		ret := wtype.NewLHComponent()
		vi := v[i]
		//	ret.ID = j.ID
		j.AddDaughterComponent(ret)
		ret.ParentID = j.ID
		ret.CName = j.Name()
		ret.Type = j.Type
		ret.Vol = vi.RawValue()
		ret.Vunit = vi.Unit().PrefixedSymbol()
		ret.Extra = j.GetExtra()
		ret.Smax = j.GetSmax()
		ret.Visc = j.GetVisc()
		ret.SetSample(true)
		reta = append(reta, ret)
	}

	return reta
}

// SampleForConcentration takes a sample of this liquid and aims for a
// particular concentration
func SampleForConcentration(l *wtype.Liquid, c wunit.Concentration) *wtype.Liquid {
	ret := wtype.NewLHComponent()
	//	ret.ID = l.ID
	l.AddDaughterComponent(ret)
	ret.ParentID = l.ID
	ret.CName = l.Name()
	ret.Type = l.Type
	ret.Conc = c.RawValue()
	ret.Cunit = c.Unit().PrefixedSymbol()
	ret.CName = l.Name()
	ret.Extra = l.GetExtra()
	ret.Smax = l.GetSmax()
	ret.Visc = l.GetVisc()
	ret.SetSample(true)
	return ret
}

// SampleMass takes a sample of this liquid and aims for a particular mass
func SampleMass(s *wtype.Liquid, m wunit.Mass, d wunit.Density) *wtype.Liquid {

	// calculate volume to add from density
	v := wunit.MasstoVolume(m, d)

	ret := wtype.NewLHComponent()
	//	ret.ID = s.ID
	s.AddDaughterComponent(ret)
	ret.ParentID = s.ID
	ret.CName = s.Name()
	ret.Type = s.Type
	ret.Vol = v.RawValue()
	ret.Vunit = v.Unit().PrefixedSymbol()
	ret.Extra = s.GetExtra()
	ret.Smax = s.GetSmax()
	ret.Visc = s.GetVisc()
	ret.SetSample(true)
	return ret
}

// SampleForTotalVolume takes a sample of this liquid to be used to make the
// solution up to a particular total volume edited to take into account the
// volume of the other solution components
func SampleForTotalVolume(l *wtype.Liquid, v wunit.Volume) *wtype.Liquid {
	ret := wtype.NewLHComponent()
	l.AddDaughterComponent(ret)
	ret.ParentID = l.ID

	ret.CName = l.Name()
	ret.Type = l.Type
	ret.Tvol = v.RawValue()
	ret.Vunit = v.Unit().PrefixedSymbol()
	ret.CName = l.Name()
	ret.Extra = l.GetExtra()
	ret.Smax = l.GetSmax()
	ret.Visc = l.GetVisc()
	ret.SetSample(true)
	return ret
}

// MixOptions are options to GenericMix
type MixOptions struct {
	Inputs      []*wtype.Liquid      // Components to mix (required)
	Instruction *wtype.LHInstruction // used to be LHSolution
	Output      *wtype.Liquid        // the resultant component
	Destination *wtype.Plate         // Destination plate; if nil, select one later
	PlateType   string               // type of destination plate
	Address     string               // Well in destination to place result; if nil, select one later
	PlateNum    int                  // which plate to stick these on
	PlateName   string               // which (named) plate to stick these on
}

// GenericMix is the general mixing entry point
func GenericMix(opt MixOptions) *wtype.LHInstruction {
	r := opt.Instruction
	if r == nil {
		r = wtype.NewLHMixInstruction()
	}
	r.Inputs = opt.Inputs

	if opt.Output != nil {
		r.AddOutput(opt.Output)
	} else {
		cmpR := wtype.NewLHComponent()
		r.AddOutput(cmpR)

		if !r.Inputs[0].IsSample() {
			r.Outputs[0].Loc = r.Inputs[0].Loc
		}

		mx := 0
		for _, c := range opt.Inputs {
			//r.Output.MixPreserveTvol(c)
			r.Outputs[0].Mix(c)
			if c.Generation() > mx {
				mx = c.Generation()
			}
		}

		r.Outputs[0].SetGeneration(mx)
	}

	if opt.Destination != nil {
		r.ContainerType = opt.Destination.Type
		r.Platetype = opt.Destination.Type
		r.SetPlateID(opt.Destination.ID)
		r.OutPlate = opt.Destination

		// if we know the well as well we should ensure that non-empty wells are respected
		if opt.Address != "" {
			w, ok := opt.Destination.Wellcoords[opt.Address]

			if !ok {
				panic(fmt.Sprintf("Cannot find well %s on plate %s name %s type %s", opt.Address, r.OutPlate.ID, r.OutPlate.Name(), r.OutPlate.Type))
			}

			if !w.IsEmpty() {
				// the instruction version has to remain unchanged
				// the returned version in the protocol has to be mixed
				w.WContents.Loc = r.OutPlate.ID + ":" + opt.Address
				r.Outputs[0] = w.WContents.Dup()
				for _, c := range opt.Inputs {
					//r.Output.MixPreserveTvol(c)
					r.Outputs[0].Mix(c)

				}
				// we also need to make sure the instruction explicitly mentions the component
				cmps := make([]*wtype.Liquid, 0, len(opt.Inputs)+1)
				cmps = append(cmps, w.WContents.Dup())
				cmps = append(cmps, opt.Inputs...)
				opt.Inputs = cmps
				r.Inputs = wtype.CopyComponentArray(cmps)
			}
			// empty wells stay empty
			//r.Output.Loc = r.OutPlate.ID + ":" + opt.Address
		}
	}

	if opt.PlateType != "" {
		r.ContainerType = opt.PlateType
		r.Platetype = opt.PlateType
	}

	if len(opt.Address) > 0 {
		r.Welladdress = opt.Address
	}

	if opt.PlateNum > 0 {
		r.Majorlayoutgroup = opt.PlateNum - 1
	}

	if opt.PlateName != "" {
		r.PlateName = opt.PlateName
	}

	// ensure results are given the correct final volumes
	// ... by definition this is either the sum of the volumes
	// or the total volume if specified

	tVol := findTVolOrPanic(opt.Inputs)

	if !tVol.IsZero() {
		r.Outputs[0].SetVolume(tVol)
	}

	return r
}

func findTVolOrPanic(components []*wtype.Liquid) wunit.Volume {
	tv := wunit.NewVolume(0.0, "ul")

	for _, c := range components {
		ctv := c.TotalVolume()

		if !(tv.IsZero() || ctv.IsZero() || tv.EqualTo(ctv)) {
			panic(fmt.Sprintf("Mix ERROR: Multiple contradictory total volumes specified %s %s", tv, ctv))
		}

		if tv.IsZero() {
			tv = ctv
		}
	}

	return tv
}

// TODO: The functions below will be deleted soon as they do not generate liquid
// handling instructions

// Mix the specified wtype.LHComponents together and leave the destination TBD
func Mix(inputs ...*wtype.Liquid) *wtype.Liquid {
	r := GenericMix(MixOptions{
		Inputs: inputs,
	})
	return r.Outputs[0]
}

// MixInto the specified wtype.LHComponents together into a specific plate
func MixInto(destination *wtype.Plate, address string, inputs ...*wtype.Liquid) *wtype.Liquid {
	r := GenericMix(MixOptions{
		Inputs:      inputs,
		Destination: destination,
		Address:     address,
	})

	return r.Outputs[0]
}

// MixTo the specified wtype.LHComponents together into a plate of a particular type
func MixTo(platetype string, address string, platenum int, inputs ...*wtype.Liquid) *wtype.Liquid {
	r := GenericMix(MixOptions{
		Inputs:    inputs,
		PlateType: platetype,
		Address:   address,
		PlateNum:  platenum,
	})
	return r.Outputs[0]
}
