// liquidhandling/lhtypes.Go: Part of the Antha language
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
// contact license@antha-lang.Org or write to the Antha team c/o
// Synthace Ltd. The London Bioscience Innovation Centre
// 2 Royal College St, London NW1 0NH UK

package wtype

import "github.com/antha-lang/antha/antha/anthalib/wunit"

//TODO add extra properties, i.e. filter
type LHTip struct {
	ID       string
	Type     string
	Mnfr     string
	Dirty    bool
	MaxVol   wunit.Volume
	MinVol   wunit.Volume
	Filtered bool
}

/*
	ID          string
	Name        string
	Minvol      wunit.Volume
	Maxvol      wunit.Volume
	Minspd      wunit.FlowRate
	Maxspd      wunit.FlowRate
	Multi       int
	Independent bool
	Orientation int
	Head        int
*/

func (tip *LHTip) GetParams() *LHChannelParameter {
	// be safe
	if tip.IsNil() {
		return nil
	}

	lhcp := LHChannelParameter{Name: tip.Type + "Params", Minvol: tip.MinVol, Maxvol: tip.MaxVol, Multi: 1, Independent: false, Orientation: LHVChannel}
	return &lhcp
}

func (tip *LHTip) IsNil() bool {
	if tip == nil || tip.Type == "" || tip.MaxVol.IsZero() || tip.MinVol.IsZero() {
		return true
	}
	return false
}

func (tip *LHTip) Dup() *LHTip {
	t := NewLHTip(tip.Mnfr, tip.Type, tip.MinVol.RawValue(), tip.MaxVol.RawValue(), tip.MinVol.Unit().PrefixedSymbol(), tip.Filtered)
	t.Dirty = tip.Dirty
	return t
}

func (tip *LHTip) DupKeepID() *LHTip {
	t := tip.Dup()
	t.ID = tip.ID
	return t
}

func NewLHTip(mfr, ttype string, minvol, maxvol float64, volunit string, filtered bool) *LHTip {
	var lht LHTip
	//	lht.ID = "tip-" + GetUUID()
	lht.ID = GetUUID()
	lht.Mnfr = mfr
	lht.Type = ttype
	lht.MaxVol = wunit.NewVolume(maxvol, volunit)
	lht.MinVol = wunit.NewVolume(minvol, volunit)
	lht.Filtered = filtered
	return &lht
}

func CopyTip(tt LHTip) *LHTip {
	return &tt
}
