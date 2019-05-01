// /anthalib/simulator/liquidhandling/simulator_test.go: Part of the Antha language
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
	"strings"
	"testing"

	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/microArch/driver/liquidhandling"
	"github.com/antha-lang/antha/microArch/simulator"
)

//
// Code for specifying a VLH
//

type LayoutParams struct {
	Name string
	Xpos float64
	Ypos float64
	Zpos float64
}

type UnitParams struct {
	Value float64
	Unit  string
}

type ChannelParams struct {
	Name        string
	Platform    string
	Minvol      UnitParams
	Maxvol      UnitParams
	Minrate     UnitParams
	Maxrate     UnitParams
	multi       int
	Independent bool
	Orientation wtype.ChannelOrientation
	Head        int
}

func makeLHChannelParameter(cp ChannelParams) *wtype.LHChannelParameter {
	return wtype.NewLHChannelParameter(cp.Name,
		cp.Platform,
		wunit.NewVolume(cp.Minvol.Value, cp.Minvol.Unit),
		wunit.NewVolume(cp.Maxvol.Value, cp.Maxvol.Unit),
		wunit.NewFlowRate(cp.Minrate.Value, cp.Minrate.Unit),
		wunit.NewFlowRate(cp.Maxrate.Value, cp.Maxrate.Unit),
		cp.multi,
		cp.Independent,
		cp.Orientation,
		cp.Head)
}

type AdaptorParams struct {
	Name    string
	Mfg     string
	Channel ChannelParams
}

func makeLHAdaptor(ap AdaptorParams) *wtype.LHAdaptor {
	return wtype.NewLHAdaptor(ap.Name,
		ap.Mfg,
		makeLHChannelParameter(ap.Channel))
}

type HeadParams struct {
	Name         string
	Mfg          string
	Channel      ChannelParams
	Adaptor      AdaptorParams
	TipBehaviour wtype.TipLoadingBehaviour
}

func makeLHHead(hp HeadParams) *wtype.LHHead {
	ret := wtype.NewLHHead(hp.Name, hp.Mfg, makeLHChannelParameter(hp.Channel))
	ret.Adaptor = makeLHAdaptor(hp.Adaptor)
	ret.TipLoading = hp.TipBehaviour
	return ret
}

type HeadAssemblyParams struct {
	MotionLimits    *wtype.BBox
	PositionOffsets []wtype.Coordinates3D
	Heads           []HeadParams
	VelocityLimits  *wtype.VelocityRange
}

func makeLHHeadAssembly(ha HeadAssemblyParams) *wtype.LHHeadAssembly {
	ret := wtype.NewLHHeadAssembly(ha.MotionLimits)
	for _, pos := range ha.PositionOffsets {
		ret.AddPosition(pos)
	}
	for _, h := range ha.Heads {
		if err := ret.LoadHead(makeLHHead(h)); err != nil {
			panic(err)
		}
	}
	ret.VelocityLimits = ha.VelocityLimits.Dup()
	return ret
}

type LHPropertiesParams struct {
	Name                string
	Mfg                 string
	Layouts             []LayoutParams
	HeadAssemblies      []HeadAssemblyParams
	TipPreferences      []string
	InputPreferences    []string
	OutputPreferences   []string
	TipwastePreferences []string
	WashPreferences     []string
	WastePreferences    []string
}

func makeLHProperties(p *LHPropertiesParams) *liquidhandling.LHProperties {

	layout := make(map[string]*wtype.LHPosition)
	for _, lp := range p.Layouts {
		layout[lp.Name] = wtype.NewLHPosition(lp.Name, wtype.Coordinates3D{X: lp.Xpos, Y: lp.Ypos, Z: lp.Zpos}, wtype.SBSFootprint)
	}

	lhp := liquidhandling.NewLHProperties(p.Name, p.Mfg, liquidhandling.LLLiquidHandler, liquidhandling.DisposableTips, layout)

	lhp.HeadAssemblies = make([]*wtype.LHHeadAssembly, 0, len(p.HeadAssemblies))
	for _, ha := range p.HeadAssemblies {
		lhp.HeadAssemblies = append(lhp.HeadAssemblies, makeLHHeadAssembly(ha))
	}
	lhp.Heads = lhp.GetLoadedHeads()

	lhp.Preferences = &liquidhandling.LayoutOpt{
		Tipboxes:  p.TipPreferences,
		Inputs:    p.InputPreferences,
		Outputs:   p.OutputPreferences,
		Tipwastes: p.TipwastePreferences,
		Washes:    p.WashPreferences,
		Wastes:    p.WastePreferences,
	}

	return lhp
}

type LHWellParams struct {
	crds    wtype.WellCoords
	vunit   string
	vol     float64
	rvol    float64
	shape   *wtype.Shape
	bott    wtype.WellBottomType
	xdim    float64
	ydim    float64
	zdim    float64
	bottomh float64
	dunit   string
}

func makeLHWell(p *LHWellParams) *wtype.LHWell {
	w := wtype.NewLHWell(
		p.vunit,
		p.vol,
		p.rvol,
		p.shape.Dup(),
		p.bott,
		p.xdim,
		p.ydim,
		p.zdim,
		p.bottomh,
		p.dunit)
	w.Crds = p.crds
	return w
}

type LHPlateParams struct {
	platetype   string
	mfr         string
	nrows       int
	ncols       int
	size        wtype.Coordinates3D
	welltype    LHWellParams
	wellXOffset float64
	wellYOffset float64
	wellXStart  float64
	wellYStart  float64
	wellZStart  float64
}

func makeLHPlate(p *LHPlateParams, name string) *wtype.Plate {
	r := wtype.NewLHPlate(p.platetype,
		p.mfr,
		p.nrows,
		p.ncols,
		p.size,
		makeLHWell(&p.welltype),
		p.wellXOffset,
		p.wellYOffset,
		p.wellXStart,
		p.wellYStart,
		p.wellZStart)
	r.PlateName = name
	return r
}

type LHTipParams struct {
	mfr             string
	ttype           string
	minvol          float64
	maxvol          float64
	volunit         string
	filtered        bool
	shape           *wtype.Shape
	effectiveHeight float64
}

func makeLHTip(p *LHTipParams) *wtype.LHTip {
	return wtype.NewLHTip(p.mfr,
		p.ttype,
		p.minvol,
		p.maxvol,
		p.volunit,
		p.filtered,
		p.shape.Dup(),
		p.effectiveHeight)
}

type LHTipboxParams struct {
	nrows        int
	ncols        int
	size         wtype.Coordinates3D
	manufacturer string
	boxtype      string
	tiptype      LHTipParams
	well         LHWellParams
	tipxoffset   float64
	tipyoffset   float64
	tipxstart    float64
	tipystart    float64
	tipzstart    float64
}

func makeLHTipbox(p *LHTipboxParams, name string) *wtype.LHTipbox {
	r := wtype.NewLHTipbox(p.nrows,
		p.ncols,
		p.size,
		p.manufacturer,
		p.boxtype,
		makeLHTip(&p.tiptype),
		makeLHWell(&p.well),
		p.tipxoffset,
		p.tipyoffset,
		p.tipystart,
		p.tipxstart,
		p.tipzstart)
	r.Boxname = name
	return r
}

type LHTipwasteParams struct {
	capacity   int
	typ        string
	mfr        string
	size       wtype.Coordinates3D
	w          LHWellParams
	wellxstart float64
	wellystart float64
	wellzstart float64
}

func makeLHTipWaste(p *LHTipwasteParams, name string) *wtype.LHTipwaste {
	r := wtype.NewLHTipwaste(p.capacity,
		p.typ,
		p.mfr,
		p.size,
		makeLHWell(&p.w),
		p.wellxstart,
		p.wellystart,
		p.wellzstart)
	r.Name = name
	return r
}

/*
 * ######################################## utils
 */

/* -- remove for linting
//test that the worst reported error severity is the worst
func test_worst(t *testing.T, errors []*simulator.SimulationError, worst simulator.ErrorSeverity) {
	s := simulator.SeverityNone
	for _, err := range errors {
		if err.Severity() > s {
			s = err.Severity()
		}
	}

	if s != worst {
		t.Errorf("Expected maximum severity %v, actual maximum severity %v", worst, s)
	}
}*/

//return subset of a not in b
func setSubtract(a, b []string) []string {
	ret := []string{}
	for _, va := range a {
		c := false
		for _, vb := range b {
			if c = (va == vb); c {
				break
			}
		}
		if !c {
			ret = append(ret, va)
		}
	}
	return ret
}

/*
 * ####################################### Default Types
 */

func defaultLHPlateProps() *LHPlateParams {
	params := LHPlateParams{
		platetype: "plate",
		mfr:       "test_plate_mfr",
		nrows:     8,
		ncols:     12,
		size:      wtype.Coordinates3D{X: 127.76, Y: 85.48, Z: 25.7},
		welltype: LHWellParams{
			crds:    wtype.ZeroWellCoords(),
			vunit:   "ul",
			vol:     200,
			rvol:    5,
			shape:   wtype.NewShape(wtype.BoxShape, "mm", 5.5, 5.5, 20.4),
			bott:    wtype.VWellBottom,
			xdim:    5.5,
			ydim:    5.5,
			zdim:    20.4,
			bottomh: 1.4,
			dunit:   "mm",
		},
		wellXOffset: 9.,
		wellYOffset: 9.,
		wellXStart:  0.,
		wellYStart:  0.,
		wellZStart:  5.3,
	}

	return &params
}

func defaultLHPlate(name string) *wtype.Plate {
	params := defaultLHPlateProps()
	return makeLHPlate(params, name)
}

//This plate will fill into the next door position on the robot
func wideLHPlate(name string) *wtype.Plate {
	params := defaultLHPlateProps()
	params.size.X = 300.
	return makeLHPlate(params, name)
}

func troughLHPlateProps() *LHPlateParams {
	params := LHPlateParams{
		platetype: "trough",
		mfr:       "test_trough_mfr",
		nrows:     1,
		ncols:     12,
		size:      wtype.Coordinates3D{X: 127.76, Y: 85.48, Z: 45.8},
		welltype: LHWellParams{
			crds:    wtype.ZeroWellCoords(),
			vunit:   "ul",
			vol:     15000,
			rvol:    5000,
			shape:   wtype.NewShape(wtype.BoxShape, "mm", 8.2, 72.0, 41.3),
			bott:    wtype.FlatWellBottom,
			xdim:    8.2,
			ydim:    72.0,
			zdim:    41.3,
			bottomh: 4.7,
			dunit:   "mm",
		},
		wellXOffset: 9.,
		wellYOffset: 9.,
		wellXStart:  0.,
		wellYStart:  0.,
		wellZStart:  4.5,
	}

	return &params
}

func troughLHPlate(name string) *wtype.Plate {
	params := troughLHPlateProps()
	plate := makeLHPlate(params, name)
	targets := []wtype.Coordinates3D{
		{X: 0.0, Y: -31.5, Z: 0.0},
		{X: 0.0, Y: -22.5, Z: 0.0},
		{X: 0.0, Y: -13.5, Z: 0.0},
		{X: 0.0, Y: -4.5, Z: 0.0},
		{X: 0.0, Y: 4.5, Z: 0.0},
		{X: 0.0, Y: 13.5, Z: 0.0},
		{X: 0.0, Y: 22.5, Z: 0.0},
		{X: 0.0, Y: 31.5, Z: 0.0},
	}
	plate.Welltype.SetWellTargets("Head0 Adaptor", targets)
	return plate
}

func defaultLHTipbox(name string) *wtype.LHTipbox {
	params := LHTipboxParams{
		nrows:        8,
		ncols:        12,
		size:         wtype.Coordinates3D{X: 127.76, Y: 85.48, Z: 60.13},
		manufacturer: "test Tipbox mfg",
		boxtype:      "tipbox",
		tiptype: LHTipParams{
			mfr:             "test_tip mfg",
			ttype:           "test_tip type",
			minvol:          50,
			maxvol:          1000,
			volunit:         "ul",
			filtered:        false,
			shape:           wtype.NewShape(wtype.BoxShape, "mm", 7.3, 7.3, 51.2),
			effectiveHeight: 44.7,
		},
		well: LHWellParams{
			crds:    wtype.ZeroWellCoords(),
			vunit:   "ul",
			vol:     1000,
			rvol:    50,
			shape:   wtype.NewShape(wtype.BoxShape, "mm", 7.3, 7.3, 51.2),
			bott:    wtype.VWellBottom,
			xdim:    7.3,
			ydim:    7.3,
			zdim:    51.2,
			bottomh: 0.0,
			dunit:   "mm",
		},
		tipxoffset: 9.,
		tipyoffset: 9.,
		tipxstart:  0.,
		tipystart:  0.,
		tipzstart:  10.,
	}

	return makeLHTipbox(&params, name)
}

func smallLHTipbox(name string) *wtype.LHTipbox {
	params := LHTipboxParams{
		nrows:        8,
		ncols:        12,
		size:         wtype.Coordinates3D{X: 127.76, Y: 85.48, Z: 60.13},
		manufacturer: "test Tipbox mfg",
		boxtype:      "tipbox",
		tiptype: LHTipParams{
			mfr:             "test_tip mfg",
			ttype:           "test_tip type",
			minvol:          0,
			maxvol:          200,
			volunit:         "ul",
			filtered:        false,
			shape:           wtype.NewShape(wtype.CircleShape, "mm", 7.3, 7.3, 51.2),
			effectiveHeight: 44.7,
		},
		well: LHWellParams{
			crds:    wtype.ZeroWellCoords(),
			vunit:   "ul",
			vol:     1000,
			rvol:    50,
			shape:   wtype.NewShape(wtype.CircleShape, "mm", 7.3, 7.3, 51.2),
			bott:    wtype.VWellBottom,
			xdim:    7.3,
			ydim:    7.3,
			zdim:    51.2,
			bottomh: 0.0,
			dunit:   "mm",
		},
		tipxoffset: 9.,
		tipyoffset: 9.,
		tipxstart:  0.,
		tipystart:  0.,
		tipzstart:  10.,
	}

	return makeLHTipbox(&params, name)
}

func defaultLHTipwaste(name string) *wtype.LHTipwaste {
	params := LHTipwasteParams{
		capacity: 700,
		typ:      "tipwaste",
		mfr:      "testTipwaste mfr",
		size:     wtype.Coordinates3D{X: 127.76, Y: 85.48, Z: 92.0},
		w: LHWellParams{
			crds:    wtype.ZeroWellCoords(),
			vunit:   "ul",
			vol:     800000.0,
			rvol:    800000.0,
			shape:   wtype.NewShape(wtype.CircleShape, "mm", 123.0, 80.0, 92.0),
			bott:    wtype.VWellBottom,
			xdim:    123.0,
			ydim:    80.0,
			zdim:    92.0,
			bottomh: 0.0,
			dunit:   "mm",
		},
		wellxstart: 49.5,
		wellystart: 31.5,
		wellzstart: 0.0,
	}
	return makeLHTipWaste(&params, name)
}

func defaultLHProperties() *liquidhandling.LHProperties {
	validProps := LHPropertiesParams{
		Name: "Device Name",
		Mfg:  "Device Manufaturer",
		Layouts: []LayoutParams{
			{"tipbox_1", 0.0, 0.0, 0.0},
			{"tipbox_2", 200.0, 0.0, 0.0},
			{"input_1", 400.0, 0.0, 0.0},
			{"input_2", 0.0, 200.0, 0.0},
			{"output_1", 200.0, 200.0, 0.0},
			{"output_2", 400.0, 200.0, 0.0},
			{"tipwaste", 0.0, 400.0, 0.0},
			{"wash", 200.0, 400.0, 0.0},
			{"waste", 400.0, 400.0, 0.0},
		},
		HeadAssemblies: []HeadAssemblyParams{
			{
				PositionOffsets: []wtype.Coordinates3D{{X: 0, Y: 0, Z: 0}},
				Heads: []HeadParams{
					{
						Name: "Head0 Name",
						Mfg:  "Head0 Manufacturer",
						Channel: ChannelParams{
							Name:        "Head0 ChannelParams",
							Platform:    "Head0 Platform",
							Minvol:      UnitParams{0.1, "ul"},
							Maxvol:      UnitParams{1., "ml"},
							Minrate:     UnitParams{0.1, "ml/min"},
							Maxrate:     UnitParams{10., "ml/min"},
							multi:       8,
							Independent: false,
							Orientation: wtype.LHVChannel,
							Head:        0,
						},
						Adaptor: AdaptorParams{
							Name: "Head0 Adaptor",
							Mfg:  "Head0 Adaptor Manufacturer",
							Channel: ChannelParams{
								Name:        "Head0 Adaptor ChannelParams",
								Platform:    "Head0 Adaptor Platform",
								Minvol:      UnitParams{0.1, "ul"},
								Maxvol:      UnitParams{1., "ml"},
								Minrate:     UnitParams{0.1, "ml/min"},
								Maxrate:     UnitParams{10., "ml/min"},
								multi:       8,
								Independent: false,
								Orientation: wtype.LHVChannel,
								Head:        0,
							},
						},
						TipBehaviour: wtype.TipLoadingBehaviour{},
					},
				},
			},
		},
		TipPreferences:      []string{"tipbox_1", "tipbox_2"},
		InputPreferences:    []string{"input_1", "input_2"},
		OutputPreferences:   []string{"output_1", "output_2"},
		TipwastePreferences: []string{"tipwaste"},
		WashPreferences:     []string{"wash"},
		WastePreferences:    []string{"waste"},
	}

	return makeLHProperties(&validProps)
}

func multiheadLHPropertiesProps() *LHPropertiesParams {
	x_step := 128.0
	y_step := 86.0
	validProps := LHPropertiesParams{
		Name: "Device Name",
		Mfg:  "Device Manufaturer",
		Layouts: []LayoutParams{
			{"tipbox_1", 0.0 * x_step, 0.0 * y_step, 0.0},
			{"tipbox_2", 1.0 * x_step, 0.0 * y_step, 0.0},
			{"input_1", 2.0 * x_step, 0.0 * y_step, 0.0},
			{"input_2", 0.0 * x_step, 1.0 * y_step, 0.0},
			{"output_1", 1.0 * x_step, 1.0 * y_step, 0.0},
			{"output_2", 2.0 * x_step, 1.0 * y_step, 0.0},
			{"tipwaste", 0.0 * x_step, 2.0 * y_step, 0.0},
			{"wash", 1.0 * x_step, 2.0 * y_step, 0.0},
			{"waste", 2.0 * x_step, 2.0 * y_step, 0.0},
		},
		HeadAssemblies: []HeadAssemblyParams{
			{
				MotionLimits:    wtype.NewBBox6f(0, 0, 0, 3*x_step, 3*y_step, 600.),
				PositionOffsets: []wtype.Coordinates3D{{X: -9}, {X: 9}},
				Heads: []HeadParams{
					{
						Name: "Head0 Name",
						Mfg:  "Head0 Manufacturer",
						Channel: ChannelParams{
							Name:        "Head0 ChannelParams",
							Platform:    "Head0 Platform",
							Minvol:      UnitParams{0.1, "ul"},
							Maxvol:      UnitParams{1., "ml"},
							Minrate:     UnitParams{0.1, "ml/min"},
							Maxrate:     UnitParams{10., "ml/min"},
							multi:       8,
							Independent: false,
							Orientation: wtype.LHVChannel,
							Head:        0,
						},
						Adaptor: AdaptorParams{
							Name: "Head0 Adaptor",
							Mfg:  "Head0 Adaptor Manufacturer",
							Channel: ChannelParams{
								Name:        "Head0 Adaptor ChannelParams",
								Platform:    "Head0 Adaptor Platform",
								Minvol:      UnitParams{0.1, "ul"},
								Maxvol:      UnitParams{1., "ml"},
								Minrate:     UnitParams{0.1, "ml/min"},
								Maxrate:     UnitParams{10., "ml/min"},
								multi:       8,
								Independent: false,
								Orientation: wtype.LHVChannel,
								Head:        0,
							},
						},
						TipBehaviour: wtype.TipLoadingBehaviour{
							OverrideLoadTipsCommand:    true,
							AutoRefillTipboxes:         true,
							LoadingOrder:               wtype.ColumnWise,
							VerticalLoadingDirection:   wtype.BottomToTop,
							HorizontalLoadingDirection: wtype.RightToLeft,
							ChunkingBehaviour:          wtype.ReverseSequentialTipLoading,
						},
					},
					{
						Name: "Head1 Name",
						Mfg:  "Head1 Manufacturer",
						Channel: ChannelParams{
							Name:        "Head1 ChannelParams",
							Platform:    "Head1 Platform",
							Minvol:      UnitParams{0.1, "ul"},
							Maxvol:      UnitParams{1., "ml"},
							Minrate:     UnitParams{0.1, "ml/min"},
							Maxrate:     UnitParams{10., "ml/min"},
							multi:       8,
							Independent: false,
							Orientation: wtype.LHVChannel,
							Head:        0,
						},
						Adaptor: AdaptorParams{
							Name: "Head1 Adaptor",
							Mfg:  "Head1 Adaptor Manufacturer",
							Channel: ChannelParams{
								Name:        "Head1 Adaptor ChannelParams",
								Platform:    "Head1 Adaptor Platform",
								Minvol:      UnitParams{0.1, "ul"},
								Maxvol:      UnitParams{1., "ml"},
								Minrate:     UnitParams{0.1, "ml/min"},
								Maxrate:     UnitParams{10., "ml/min"},
								multi:       8,
								Independent: false,
								Orientation: wtype.LHVChannel,
								Head:        0,
							},
						},
						TipBehaviour: wtype.TipLoadingBehaviour{
							OverrideLoadTipsCommand:    true,
							AutoRefillTipboxes:         true,
							LoadingOrder:               wtype.ColumnWise,
							VerticalLoadingDirection:   wtype.BottomToTop,
							HorizontalLoadingDirection: wtype.LeftToRight,
							ChunkingBehaviour:          wtype.ReverseSequentialTipLoading,
						},
					},
				},
			},
		},
		TipPreferences:      []string{"tipbox_1", "tipbox_2", "input_1", "input_2"},
		InputPreferences:    []string{"input_1", "input_2", "tipbox_1", "tipbox_2", "tipwaste", "waste"},
		OutputPreferences:   []string{"output_1", "output_2"},
		TipwastePreferences: []string{"tipwaste", "input_1"},
		WashPreferences:     []string{"wash"},
		WastePreferences:    []string{"waste"},
	}

	return &validProps
}

func multiheadLHProperties() *liquidhandling.LHProperties {
	return makeLHProperties(multiheadLHPropertiesProps())
}

func multiheadConstrainedLHProperties() *liquidhandling.LHProperties {
	lhp := multiheadLHPropertiesProps()
	lhp.HeadAssemblies[0].MotionLimits.Position.Z = 60
	return makeLHProperties(lhp)
}

func IndependentLHProperties() *liquidhandling.LHProperties {
	ret := defaultLHProperties()

	for _, head := range ret.Heads {
		head.Params.Independent = true
		head.Adaptor.Params.Independent = true
	}

	return ret
}

/* -- remove for linting
func default_vlh() *VirtualLiquidHandler {
	vlh := NewVirtualLiquidHandler(defaultLHProperties(), nil)
	return vlh
}
*/

/*
 * ######################################## InstructionParams
 */

type TestRobotInstruction interface {
	Convert() liquidhandling.TerminalRobotInstruction
}

//Initialize
type Initialize struct{}

func (self *Initialize) Convert() liquidhandling.TerminalRobotInstruction {
	return liquidhandling.NewInitializeInstruction()
}

//Finalize
type Finalize struct{}

func (self *Finalize) Convert() liquidhandling.TerminalRobotInstruction {
	return liquidhandling.NewFinalizeInstruction()
}

//SetPipetteSpeed
type SetPipetteSpeed struct {
	head    int
	channel int
	speed   float64
}

func (self *SetPipetteSpeed) Convert() liquidhandling.TerminalRobotInstruction {
	ret := liquidhandling.NewSetPipetteSpeedInstruction()
	ret.Head = self.head
	ret.Channel = self.channel
	ret.Speed = self.speed
	return ret
}

//SetDriveSpeed
type SetDriveSpeed struct {
	drive string
	speed float64
}

func (self *SetDriveSpeed) Convert() liquidhandling.TerminalRobotInstruction {
	ret := liquidhandling.NewSetDriveSpeedInstruction()
	ret.Drive = self.drive
	ret.Speed = self.speed
	return ret
}

//AddPlateTo
type AddPlateTo struct {
	position string
	plate    interface{}
	name     string
}

func (self *AddPlateTo) Convert() liquidhandling.TerminalRobotInstruction {
	return liquidhandling.NewAddPlateToInstruction(self.position, self.name, self.plate)
}

//LoadTips
type LoadTips struct {
	channels  []int
	head      int
	multi     int
	platetype []string
	position  []string
	well      []string
}

func (self *LoadTips) Convert() liquidhandling.TerminalRobotInstruction {
	ret := liquidhandling.NewLoadTipsInstruction()
	ret.Head = self.head
	ret.Multi = self.multi
	ret.Channels = self.channels
	ret.HolderType = self.platetype
	ret.Pos = self.position
	ret.Well = self.well
	return ret
}

//UnloadTips
type UnloadTips struct {
	channels  []int
	head      int
	multi     int
	platetype []string
	position  []string
	well      []string
}

func (self *UnloadTips) Convert() liquidhandling.TerminalRobotInstruction {
	ret := liquidhandling.NewUnloadTipsInstruction()
	ret.Head = self.head
	ret.Multi = self.multi
	ret.HolderType = self.platetype
	ret.Pos = self.position
	ret.Well = self.well
	ret.Channels = self.channels
	return ret
}

//Move
type Move struct {
	deckposition []string
	wellcoords   []string
	reference    []int
	offsetX      []float64
	offsetY      []float64
	offsetZ      []float64
	plate_type   []string
	head         int
}

func (self *Move) Convert() liquidhandling.TerminalRobotInstruction {
	ret := liquidhandling.NewMoveInstruction()
	ret.Head = self.head
	ret.Pos = self.deckposition
	ret.Well = self.wellcoords
	ret.Reference = self.reference
	ret.OffsetX = self.offsetX
	ret.OffsetY = self.offsetY
	ret.OffsetZ = self.offsetZ
	ret.Plt = self.plate_type
	return ret
}

//Aspirate
type Aspirate struct {
	volume     []float64
	overstroke bool
	head       int
	multi      int
	platetype  []string
	what       []string
	llf        []bool
}

func (self *Aspirate) Convert() liquidhandling.TerminalRobotInstruction {
	ret := liquidhandling.NewAspirateInstruction()
	volume := make([]wunit.Volume, 0, len(self.volume))
	for _, v := range self.volume {
		volume = append(volume, wunit.NewVolume(v, "ul"))
	}
	ret.Head = self.head
	ret.Volume = volume
	ret.Overstroke = self.overstroke
	ret.Multi = self.multi
	ret.Plt = self.platetype
	ret.What = self.what
	ret.LLF = self.llf
	return ret
}

//Dispense
type Dispense struct {
	volume    []float64
	blowout   []bool
	head      int
	multi     int
	platetype []string
	what      []string
	llf       []bool
}

func (self *Dispense) Convert() liquidhandling.TerminalRobotInstruction {
	ret := liquidhandling.NewDispenseInstruction()
	volume := make([]wunit.Volume, 0, len(self.volume))
	for _, v := range self.volume {
		volume = append(volume, wunit.NewVolume(v, "ul"))
	}
	ret.Head = self.head
	ret.Volume = volume
	ret.Multi = self.multi
	ret.Plt = self.platetype
	ret.What = self.what
	ret.LLF = self.llf
	return ret
}

//Mix
type Mix struct {
	head      int
	volume    []float64
	platetype []string
	cycles    []int
	multi     int
	what      []string
	blowout   []bool
}

func (self *Mix) Convert() liquidhandling.TerminalRobotInstruction {
	ret := liquidhandling.NewMixInstruction()
	volume := make([]wunit.Volume, 0, len(self.volume))
	for _, v := range self.volume {
		volume = append(volume, wunit.NewVolume(v, "ul"))
	}
	ret.Head = self.head
	ret.Volume = volume
	ret.PlateType = self.platetype
	ret.What = self.what
	ret.Blowout = self.blowout
	ret.Multi = self.multi
	ret.Cycles = self.cycles
	return ret
}

/*
 * ######################################## Setup
 */

type SetupFn func(*VirtualLiquidHandler)

func removeTipboxTips(tipbox_loc string, wells []string) *SetupFn {
	var ret SetupFn = func(vlh *VirtualLiquidHandler) {
		tipbox := vlh.GetObjectAt(tipbox_loc).(*wtype.LHTipbox)
		for _, well := range wells {
			wc := wtype.MakeWellCoords(well)
			tipbox.RemoveTip(wc)
		}
	}
	return &ret
}

func preloadAdaptorTips(head int, tipbox_loc string, channels []int) *SetupFn {
	var ret SetupFn = func(vlh *VirtualLiquidHandler) {
		adaptor, err := vlh.GetAdaptorState(head)
		if err != nil {
			panic(err)
		}
		tipbox := vlh.GetObjectAt(tipbox_loc).(*wtype.LHTipbox)

		for _, ch := range channels {
			adaptor.GetChannel(ch).LoadTip(tipbox.Tiptype.Dup())
		}
	}
	return &ret
}

func getLHComponent(what string, vol_ul float64) *wtype.Liquid {
	c := wtype.NewLHComponent()
	c.CName = what
	//madness?
	lt, _ := wtype.LiquidTypeFromString(wtype.PolicyName(what))
	c.Type = lt
	c.Vol = vol_ul
	c.Vunit = "ul"

	return c
}

func preloadFilledTips(head int, tipbox_loc string, channels []int, what string, volume float64) *SetupFn {
	var ret SetupFn = func(vlh *VirtualLiquidHandler) {
		adaptor, err := vlh.GetAdaptorState(head)
		if err != nil {
			panic(err)
		}
		tipbox := vlh.GetObjectAt(tipbox_loc).(*wtype.LHTipbox)
		tip := tipbox.Tiptype.Dup()
		c := getLHComponent(what, volume)
		if err := tip.AddComponent(c); err != nil {
			panic(err)
		}

		for _, ch := range channels {
			adaptor.GetChannel(ch).LoadTip(tip.Dup())
		}
	}
	return &ret
}

/* -- remove for linting
func fillTipwaste(tipwaste_loc string, count int) *SetupFn {
	var ret SetupFn = func(vlh *VirtualLiquidHandler) {
		tipwaste := vlh.GetObjectAt(tipwaste_loc).(*wtype.LHTipwaste)
		tipwaste.Contents += count
	}
	return &ret
}
*/

func prefillWells(plate_loc string, wells_to_fill []string, liquid_name string, volume float64) *SetupFn {
	var ret SetupFn = func(vlh *VirtualLiquidHandler) {
		plate := vlh.GetObjectAt(plate_loc).(*wtype.Plate)
		for _, well_name := range wells_to_fill {
			wc := wtype.MakeWellCoords(well_name)
			well := plate.GetChildByAddress(wc).(*wtype.LHWell)
			comp := getLHComponent(liquid_name, volume)
			err := well.AddComponent(comp)
			if err != nil {
				panic(err)
			}
		}
	}
	return &ret
}

type moveToParams struct {
	Multi        int
	Head         int
	Reference    int
	Deckposition string
	Platetype    string
	Offset       []float64
	Cols         int
	Rows         int
}

//moveTo Simplify generating Move commands when running tests by avoiding
//repeating stuff that doesn't change
func moveTo(row, col int, p moveToParams) *SetupFn {
	s_dp := make([]string, p.Multi)
	s_wc := make([]string, p.Multi)
	s_rf := make([]int, p.Multi)
	s_ox := make([]float64, p.Multi)
	s_oy := make([]float64, p.Multi)
	s_oz := make([]float64, p.Multi)
	s_pt := make([]string, p.Multi)

	for i := 0; i < p.Multi; i++ {
		if col >= 0 && col < p.Cols && row+i >= 0 && row+i < p.Rows {
			wc := wtype.WellCoords{X: col, Y: row + i}
			s_dp[i] = p.Deckposition
			s_wc[i] = wc.FormatA1()
			s_rf[i] = p.Reference
			s_ox[i] = p.Offset[0]
			s_oy[i] = p.Offset[1]
			s_oz[i] = p.Offset[2]
			s_pt[i] = p.Platetype
		}
	}

	var ret SetupFn = func(vlh *VirtualLiquidHandler) {
		vlh.Move(s_dp, s_wc, s_rf, s_ox, s_oy, s_oz, s_pt, p.Head)
	}

	return &ret
}

/*
 * ######################################## Assertions (about the final state)
 */

type AssertionFn func(*testing.T, *VirtualLiquidHandler)

//tipboxAssertion assert that the tipbox has tips missing in the given locations only
func tipboxAssertion(tipbox_loc string, missing_tips []string) *AssertionFn {
	var ret AssertionFn = func(t *testing.T, vlh *VirtualLiquidHandler) {
		mmissing_tips := make(map[string]bool)
		for _, tl := range missing_tips {
			mmissing_tips[tl] = true
		}

		if tipbox, ok := vlh.GetObjectAt(tipbox_loc).(*wtype.LHTipbox); !ok {
			t.Errorf("TipboxAssertion failed: no Tipbox found at \"%s\"", tipbox_loc)
		} else {
			errors := []string{}
			for y := 0; y < tipbox.Nrows; y++ {
				for x := 0; x < tipbox.Ncols; x++ {
					wc := wtype.WellCoords{X: x, Y: y}
					wcs := wc.FormatA1()
					if hta, etm := tipbox.HasTipAt(wc), mmissing_tips[wcs]; !hta && !etm {
						errors = append(errors, fmt.Sprintf("Unexpected tip missing at %s", wcs))
					} else if hta && etm {
						errors = append(errors, fmt.Sprintf("Unexpected tip present at %s", wcs))
					}
				}
			}
			if len(errors) > 0 {
				t.Errorf("TipboxAssertion failed: tipbox at \"%s\":\n%s", tipbox_loc, strings.Join(errors, "\n"))
			}
		}
	}
	return &ret
}

type tipDesc struct {
	channel     int
	liquid_type string
	volume      float64
}

//adaptorAssertion assert that the adaptor has tips in the given positions
func adaptorAssertion(head int, tips []tipDesc) *AssertionFn {
	var ret AssertionFn = func(t *testing.T, vlh *VirtualLiquidHandler) {
		mtips := make(map[int]bool)
		for _, td := range tips {
			mtips[td.channel] = true
		}

		adaptor, err := vlh.GetAdaptorState(head)
		if err != nil {
			panic(err)
		}
		errors := []string{}
		for ch := 0; ch < adaptor.GetChannelCount(); ch++ {
			if itl, et := adaptor.GetChannel(ch).HasTip(), mtips[ch]; itl && !et {
				errors = append(errors, fmt.Sprintf("Unexpected tip on channel %v", ch))
			} else if !itl && et {
				errors = append(errors, fmt.Sprintf("Expected tip on channel %v", ch))
			}
		}
		//now check volumes
		for _, td := range tips {
			if !adaptor.GetChannel(td.channel).HasTip() {
				continue //already reported this error
			}
			tip := adaptor.GetChannel(td.channel).GetTip()
			c := tip.Contents()
			if c.Volume().ConvertToString("ul") != td.volume || c.Name() != td.liquid_type {
				errors = append(errors, fmt.Sprintf("Channel %d: Expected tip with %.2f ul of \"%s\", got tip with %s of \"%s\"",
					td.channel, td.volume, td.liquid_type, c.Volume(), c.Name()))
			}
		}
		if len(errors) > 0 {
			t.Errorf("AdaptorAssertion failed: Head%v:\n%s", head, strings.Join(errors, "\n"))
		}
	}
	return &ret
}

//adaptorPositionAssertion assert that the adaptor has tips in the given positions
func positionAssertion(head int, origin wtype.Coordinates3D) *AssertionFn {
	var ret AssertionFn = func(t *testing.T, vlh *VirtualLiquidHandler) {
		adaptor, err := vlh.GetAdaptorState(head)
		if err != nil {
			panic(err)
		}
		or := adaptor.GetChannel(0).GetAbsolutePosition()
		//use string comparison to avoid precision errors (string printed with %.1f)
		if g, e := or.String(), origin.String(); g != e {
			t.Errorf("PositionAssertion failed: head %d should be at %s, was actually at %s", head, e, g)
		}
	}
	return &ret
}

//tipwasteAssertion assert the number of tips which should be in the tipwaste
func tipwasteAssertion(tipwaste_loc string, expected_contents int) *AssertionFn {
	var ret AssertionFn = func(t *testing.T, vlh *VirtualLiquidHandler) {
		if tipwaste, ok := vlh.GetObjectAt(tipwaste_loc).(*wtype.LHTipwaste); !ok {
			t.Errorf("TipWasteAssertion failed: no Tipwaste found at %s", tipwaste_loc)
		} else {
			if tipwaste.Contents != expected_contents {
				t.Errorf("TipwasteAssertion failed at location %s: expected %v tips, got %v",
					tipwaste_loc, expected_contents, tipwaste.Contents)
			}
		}
	}
	return &ret
}

type wellDesc struct {
	position    string
	liquid_type string
	volume      float64
}

func plateAssertion(plate_loc string, wells []wellDesc) *AssertionFn {
	var ret AssertionFn = func(t *testing.T, vlh *VirtualLiquidHandler) {
		m := map[string]bool{}
		plate := vlh.GetObjectAt(plate_loc).(*wtype.Plate)
		errs := []string{}
		for _, wd := range wells {
			m[wd.position] = true
			wc := wtype.MakeWellCoords(wd.position)
			well := plate.GetChildByAddress(wc).(*wtype.LHWell)
			c := well.Contents()
			if fmt.Sprintf("%.2f", c.Vol) != fmt.Sprintf("%.2f", wd.volume) || wd.liquid_type != c.Name() {
				errs = append(errs, fmt.Sprintf("Expected %.2ful of %s in well %s, found %.2ful of %s",
					wd.volume, wd.liquid_type, wd.position, c.Vol, c.Name()))
			}
		}
		//now check that all the other wells are empty
		for _, row := range plate.Rows {
			for _, well := range row {
				if c := well.Contents(); !m[well.Crds.FormatA1()] && !c.IsZero() {
					errs = append(errs, fmt.Sprintf("Expected empty well at %s, instead %s of %s",
						well.Crds.FormatA1(), c.Volume(), c.Name()))
				}
			}
		}

		if len(errs) > 0 {
			t.Errorf("plateAssertion failed: errors were:\n%s", strings.Join(errs, "\n"))
		}
	}
	return &ret
}

/*
 * ######################################## SimulatorTest
 */

type SimulatorTest struct {
	Name           string
	Props          *liquidhandling.LHProperties
	Setup          []*SetupFn
	Instructions   []TestRobotInstruction
	ExpectedErrors []string
	Assertions     []*AssertionFn
}

func (self *SimulatorTest) compareErrors(t *testing.T, actual []simulator.SimulationError) {
	stringErrors := make([]string, 0, len(actual))
	for _, err := range actual {
		stringErrors = append(stringErrors, err.Error())
	}
	// maybe sort alphabetically?

	missing := setSubtract(self.ExpectedErrors, stringErrors)
	extra := setSubtract(stringErrors, self.ExpectedErrors)

	errs := []string{}
	for _, s := range missing {
		errs = append(errs, fmt.Sprintf("--\"%v\"", s))
	}
	for _, s := range extra {
		errs = append(errs, fmt.Sprintf("++\"%v\"", s))
	}
	if len(missing) > 0 || len(extra) > 0 {
		t.Errorf("errors didn't match:\n\t%s", strings.Join(errs, "\t\n"))
	}
}

func (self *SimulatorTest) Run(t *testing.T) {

	if self.Props == nil {
		self.Props = defaultLHProperties()
	}
	vlh, err := NewVirtualLiquidHandler(self.Props, nil)
	if err != nil {
		t.Fatal(err)
	}

	//do setup
	if self.Setup != nil {
		for _, setup_fn := range self.Setup {
			(*setup_fn)(vlh)
		}
	}

	//run the instructions
	if self.Instructions != nil {
		instructions := make([]liquidhandling.TerminalRobotInstruction, 0, len(self.Instructions))
		for _, inst := range self.Instructions {
			instructions = append(instructions, inst.Convert())
		}
		if err := vlh.Simulate(instructions); err != nil {
			t.Error(err)
		}
	}

	//check errors
	self.compareErrors(t, vlh.GetErrors())

	//check assertions
	if self.Assertions != nil {
		for _, a := range self.Assertions {
			(*a)(t, vlh)
		}
	}
}

type SimulatorTests []SimulatorTest

func (tests SimulatorTests) Run(t *testing.T) {
	for _, test := range tests {
		t.Run(test.Name, test.Run)
	}
}
