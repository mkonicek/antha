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
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
)

// SampleAll takes all of this liquid
func SampleAll(l *wtype.LHComponent) *wtype.LHComponent {
	return Sample(l, l.Volume())
}

// Sample takes a sample of volume v from this liquid
func Sample(l *wtype.LHComponent, v wunit.Volume) *wtype.LHComponent {
	ret := wtype.NewLHComponent()
	//	ret.ID = l.ID
	l.AddDaughterComponent(ret)
	ret.ParentID = l.ID
	ret.CName = l.Name()
	ret.Type = l.Type
	ret.Vol = v.RawValue()
	ret.Vunit = v.Unit().PrefixedSymbol()
	ret.Extra = l.GetExtra()
	ret.Smax = l.GetSmax()
	ret.Visc = l.GetVisc()
	if l.Conc > 0 && len(l.Cunit) > 0 {
		ret.SetConcentration(wunit.NewConcentration(l.Conc, l.Cunit))
	}
	//ret.Conc = l.Conc
	//ret.Cunit = l.Cunit
	ret.SetSample(true)

	//logger.Track(fmt.Sprintf("SAMPLE V %s %s %s", l.ID, ret.ID, v.ToString()))

	return ret
}

// MultiSample takes an array of samples and array of corresponding volumes and
// sample them all
func MultiSample(l []*wtype.LHComponent, v []wunit.Volume) []*wtype.LHComponent {
	reta := make([]*wtype.LHComponent, 0)

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
		//	logger.Track(fmt.Sprintf("SAMPLE V %s %s %s", j.ID, ret.ID, vi.ToString()))
		ret.SetSample(true)
		reta = append(reta, ret)
	}

	return reta
}

// SampleForConcentration takes a sample of this liquid and aims for a
// particular concentration
func SampleForConcentration(l *wtype.LHComponent, c wunit.Concentration) *wtype.LHComponent {
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
	//logger.Track(fmt.Sprintf("SAMPLE C %s %s %s", l.ID, ret.ID, c.ToString()))
	return ret
}

// SampleMass takes a sample of this liquid and aims for a particular mass
func SampleMass(s *wtype.LHComponent, m wunit.Mass, d wunit.Density) *wtype.LHComponent {

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
	//logger.Track(fmt.Sprintf("SAMPLE M %s %s %s %s", s.ID, ret.ID, m.ToString(), d.ToString()))
	ret.SetSample(true)
	return ret
}

// SampleForTotalVolume takes a sample ofs this liquid to be used to make the
// solution up to a particular total volume edited to take into account the
// volume of the other solution components
func SampleForTotalVolume(l *wtype.LHComponent, v wunit.Volume) *wtype.LHComponent {
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
	//logger.Track(fmt.Sprintf("SAMPLE T %s %s %s", l.ID, ret.ID, v.ToString()))
	ret.SetSample(true)
	return ret
}

// MixOptions are options to GenericMix
type MixOptions struct {
	Components  []*wtype.LHComponent // Components to mix (required)
	Instruction *wtype.LHInstruction // used to be LHSolution
	Result      *wtype.LHComponent   // the resultant component
	Destination *wtype.LHPlate       // Destination plate; if nil, select one later
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
	r.Components = opt.Components

	if opt.Result != nil {
		r.Result = opt.Result
	} else {
		r.Result = wtype.NewLHComponent()
		mx := 0
		for _, c := range opt.Components {
			r.Result.MixPreserveTvol(c)
			if c.Generation() > mx {
				mx = c.Generation()
			}
		}
		r.Result.SetGeneration(mx)
	}

	if opt.Destination != nil {
		r.ContainerType = opt.Destination.Type
		r.Platetype = opt.Destination.Type
		r.SetPlateID(opt.Destination.ID)
		r.OutPlate = opt.Destination
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

	// oh yus oh yus oh yus

	s := ""
	for _, v := range r.Components {
		s += v.CName + "-" + v.ID + " "
	}

	//fmt.Println("GENERATION: ", r.Result.Generation(), "MIXING : ", s, " RESULT: ", r.Result.CName+"-"+r.Result.ID)

	return r
}

// TODO: he functions below will be deleted soon as they do not generate liquid
// handling instructions

// Mix the specified wtype.LHComponents together and leave the destination TBD
func Mix(components ...*wtype.LHComponent) *wtype.LHComponent {
	r := GenericMix(MixOptions{
		Components: components,
	})
	return r.Result
}

// MixInto the specified wtype.LHComponents together into a specific plate
func MixInto(destination *wtype.LHPlate, address string, components ...*wtype.LHComponent) *wtype.LHComponent {
	r := GenericMix(MixOptions{
		Components:  components,
		Destination: destination,
		Address:     address,
	})

	return r.Result
}

// MixTo the specified wtype.LHComponents together into a plate of a particular type
func MixTo(platetype string, address string, platenum int, components ...*wtype.LHComponent) *wtype.LHComponent {
	r := GenericMix(MixOptions{
		Components: components,
		PlateType:  platetype,
		Address:    address,
		PlateNum:   platenum,
	})
	return r.Result
}
