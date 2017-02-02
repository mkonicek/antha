// anthalib/factory/make_liquidhandler_library.go: Part of the Antha language
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

package mixer

import (
	"fmt"

	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/microArch/driver/liquidhandling"
	. "github.com/antha-lang/antha/microArch/factory"
)

func SetUpTipsFor(lhp *liquidhandling.LHProperties) *liquidhandling.LHProperties {
	tips := GetTipList()
	for _, tt := range tips {
		tb := GetTipByType(tt)
		if tb.Mnfr == lhp.Mnfr || lhp.Mnfr == "MotherNature" {
			lhp.Tips = append(lhp.Tips, tb.Tips[0][0])
		}
	}
	return lhp
}

func makeLiquidhandlerLibrary() map[string]*liquidhandling.LHProperties {
	robots := make(map[string]*liquidhandling.LHProperties, 2)
	robots["TecanEvo"] = makeEvo()
	robots["Manual"] = makeManual()
	return robots
}

func makeManual() *liquidhandling.LHProperties {
	//	tips := GetTipList()

	// dummy layout of 25 positions... arbitrary limitation

	x := 0.0
	y := 0.0
	z := 0.0
	xinc := 100.0
	yinc := 100.0

	i := 0
	layout := make(map[string]wtype.Coordinates)
	for xi := 0; xi < 5; xi++ {
		for yi := 0; yi < 5; yi++ {
			posname := fmt.Sprintf("position_%d", i+1)
			crds := wtype.Coordinates{x, y, z}
			layout[posname] = crds
			i += 1
			y += yinc
		}
		x += xinc
	}
	lhp := liquidhandling.NewLHProperties(25, "Human", "MotherNature", "discrete", "disposable", layout)

	SetUpTipsFor(lhp)

	lhp.Tip_preferences = []string{"tips1", "tips2", "tips3", "tips4"}
	lhp.Input_preferences = []string{"in1", "in2", "in3", "in4"}
	lhp.Output_preferences = []string{"out1", "out2", "out3", "out4"}
	lhp.Tipwaste_preferences = []string{"tip_waste"}
	lhp.Wash_preferences = []string{"tip_wash"}
	lhp.Waste_preferences = []string{"liquid_waste"}

	minvol := wunit.NewVolume(200, "ul")
	maxvol := wunit.NewVolume(1000, "ul")
	minspd := wunit.NewFlowRate(0.5, "ml/min")
	maxspd := wunit.NewFlowRate(2, "ml/min")

	hvconfig := wtype.NewLHChannelParameter("P1000Config", "Gilson", minvol, maxvol, minspd, maxspd, 1, false, wtype.LHVChannel, 0)
	hvadaptor := wtype.NewLHAdaptor("P1000", "Gilson", hvconfig)

	minvol = wunit.NewVolume(50, "ul")
	maxvol = wunit.NewVolume(200, "ul")
	minspd = wunit.NewFlowRate(0.1, "ml/min")
	maxspd = wunit.NewFlowRate(0.5, "ml/min")

	mvconfig := wtype.NewLHChannelParameter("P200Config", "Gilson", minvol, maxvol, minspd, maxspd, 1, false, wtype.LHVChannel, 0)
	mvadaptor := wtype.NewLHAdaptor("P200", "Gilson", mvconfig)

	minvol = wunit.NewVolume(2, "ul")
	maxvol = wunit.NewVolume(20, "ul")
	minspd = wunit.NewFlowRate(0.1, "ml/min")
	maxspd = wunit.NewFlowRate(0.5, "ml/min")

	lmvconfig := wtype.NewLHChannelParameter("P20Config", "Gilson", minvol, maxvol, minspd, maxspd, 1, false, wtype.LHVChannel, 0)
	lmvadaptor := wtype.NewLHAdaptor("P20", "Gilson", lmvconfig)

	minvol = wunit.NewVolume(1, "ul")
	maxvol = wunit.NewVolume(10, "ul")
	minspd = wunit.NewFlowRate(0.1, "ml/min")
	maxspd = wunit.NewFlowRate(0.5, "ml/min")

	lvconfig := wtype.NewLHChannelParameter("P10Config", "Gilson", minvol, maxvol, minspd, maxspd, 1, false, wtype.LHVChannel, 0)
	lvadaptor := wtype.NewLHAdaptor("P10", "Gilson", lvconfig)

	minvol = wunit.NewVolume(0.2, "ul")
	maxvol = wunit.NewVolume(2, "ul")
	minspd = wunit.NewFlowRate(0.1, "ml/min")
	maxspd = wunit.NewFlowRate(0.5, "ml/min")

	vlvconfig := wtype.NewLHChannelParameter("P2Config", "Gilson", minvol, maxvol, minspd, maxspd, 1, false, wtype.LHVChannel, 0)
	vlvadaptor := wtype.NewLHAdaptor("P2", "Gilson", vlvconfig)

	minvol = wunit.NewVolume(0.2, "ul")
	maxvol = wunit.NewVolume(5000, "ul")
	headparams := wtype.NewLHChannelParameter("LabHand", "Mothernature", minvol, maxvol, minspd, maxspd, 8, false, wtype.LHVChannel, 0)
	head := wtype.NewLHHead("LabHand", "MotherNature", headparams)
	head.Adaptor = hvadaptor

	lhp.Adaptors = append(lhp.Adaptors, hvadaptor)
	lhp.Adaptors = append(lhp.Adaptors, mvadaptor)
	lhp.Adaptors = append(lhp.Adaptors, lmvadaptor)
	lhp.Adaptors = append(lhp.Adaptors, lvadaptor)
	lhp.Adaptors = append(lhp.Adaptors, vlvadaptor)
	lhp.Heads = append(lhp.Heads, head)
	lhp.HeadsLoaded = append(lhp.HeadsLoaded, head)
	return lhp
}

func makeEvo() *liquidhandling.LHProperties {
	// These aren't real measurements but should be sufficient for relative
	// positioning
	layout := make(map[string]wtype.Coordinates)
	numZ := 1
	numY := 4
	numX := 4
	for z := 0; z < numZ; z += 1 {
		for x := 0; x < numX; x += 1 {
			for y := 0; y < numY; y += 1 {
				//  7  8  9
				// 10 11 12
				//  1  2  3
				//  4  5  6
				idx := z*numY*numX + y*numX + x
				posname := fmt.Sprintf("position_%d", idx+1)
				layout[posname] = wtype.Coordinates{X: float64(x), Y: float64(y), Z: float64(z)}
			}
		}
	}

	lhp := liquidhandling.NewLHProperties(16, "Evo", "Tecan", "discrete", "disposable", layout)

	// get tips permissible from the factory
	SetUpTipsFor(lhp)

	lhp.Tip_preferences = []string{"position_1", "position_2", "position_3"}
	lhp.Input_preferences = []string{"position_4", "position_5", "position_6"}
	lhp.Output_preferences = []string{"position_7", "position_8", "position_9"}
	lhp.Wash_preferences = []string{"position_10"}
	lhp.Waste_preferences = []string{"position_11"}
	lhp.Tipwaste_preferences = []string{"position_12", "position_13"}

	minvol3 := wunit.NewVolume(1.0, "ul")
	maxvol3 := wunit.NewVolume(1000, "ul")
	minspd3 := wunit.NewFlowRate(0.1, "ml/min")
	maxspd3 := wunit.NewFlowRate(0.5, "ml/min")
	headparams := wtype.NewLHChannelParameter("LiHa", "TecanEvo", minvol3, maxvol3, minspd3, maxspd3, 8, false, wtype.LHVChannel, 0)
	head := wtype.NewLHHead("LiHa", "Tecan", headparams)

	adaptor := GetAdaptorByType("TecanDiTiAdaptor")

	head.Adaptor = adaptor

	lhp.Adaptors = append(lhp.Adaptors, adaptor)
	lhp.Heads = append(lhp.Heads, head)
	lhp.HeadsLoaded = append(lhp.HeadsLoaded, head)

	return lhp
}

func GetAdaptorByType(typ string) *wtype.LHAdaptor {
	adaptors := makeAdaptors()

	a, ok := adaptors[typ]

	if ok {
		return a
	} else {
		return nil
	}
}

func makeAdaptors() map[string]*wtype.LHAdaptor {
	ret := make(map[string]*wtype.LHAdaptor)

	minvol := wunit.NewVolume(1, "ul")
	maxvol := wunit.NewVolume(1000, "ul")
	minspd := wunit.NewFlowRate(0.5, "ml/min")
	maxspd := wunit.NewFlowRate(2, "ml/min")

	config := wtype.NewLHChannelParameter("config", "TecanEvo", minvol, maxvol, minspd, maxspd, 8, false, wtype.LHVChannel, 0)
	adaptor := wtype.NewLHAdaptor("Adaptor", "Tecan", config)

	ret["TecanDiTiAdaptor"] = adaptor

	return ret
}

func GetLiquidhandlerByType(typ string) *liquidhandling.LHProperties {
	liquidhandlers := makeLiquidhandlerLibrary()
	t := liquidhandlers[typ]
	return t.Dup()
}

func LiquidhandlerList() []string {
	liquidhandlers := makeLiquidhandlerLibrary()
	kz := make([]string, len(liquidhandlers))
	x := 0
	for name, _ := range liquidhandlers {
		kz[x] = name
		x += 1
	}
	return kz
}
