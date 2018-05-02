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

package liquidhandling_test

import (
	"fmt"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/microArch/driver/liquidhandling"
	"github.com/antha-lang/antha/microArch/simulator"
	lh "github.com/antha-lang/antha/microArch/simulator/liquidhandling"
	"strings"
	"testing"
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
	Orientation int
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
	Name    string
	Mfg     string
	Channel ChannelParams
	Adaptor AdaptorParams
}

func makeLHHead(hp HeadParams) *wtype.LHHead {
	ret := wtype.NewLHHead(hp.Name, hp.Mfg, makeLHChannelParameter(hp.Channel))
	ret.Adaptor = makeLHAdaptor(hp.Adaptor)
	return ret
}

type LHPropertiesParams struct {
	Name                 string
	Mfg                  string
	Layouts              []LayoutParams
	Heads                []HeadParams
	Tip_preferences      []string
	Input_preferences    []string
	Output_preferences   []string
	Tipwaste_preferences []string
	Wash_preferences     []string
	Waste_preferences    []string
}

func makeLHProperties(p *LHPropertiesParams) *liquidhandling.LHProperties {

	layout := make(map[string]wtype.Coordinates)
	for _, lp := range p.Layouts {
		layout[lp.Name] = wtype.Coordinates{X: lp.Xpos, Y: lp.Ypos, Z: lp.Zpos}
	}

	lhp := liquidhandling.NewLHProperties(len(layout), p.Name, p.Mfg, liquidhandling.LLLiquidHandler, liquidhandling.DisposableTips, layout)

	lhp.Heads = make([]*wtype.LHHead, 0)
	for _, hp := range p.Heads {
		lhp.Heads = append(lhp.Heads, makeLHHead(hp))
	}

	lhp.Tip_preferences = p.Tip_preferences
	lhp.Input_preferences = p.Input_preferences
	lhp.Output_preferences = p.Output_preferences
	lhp.Tipwaste_preferences = p.Tipwaste_preferences
	lhp.Wash_preferences = p.Wash_preferences
	lhp.Waste_preferences = p.Waste_preferences

	return lhp
}

type ShapeParams struct {
	name       string
	lengthunit string
	h          float64
	w          float64
	d          float64
}

func makeShape(p *ShapeParams) *wtype.Shape {
	return wtype.NewShape(p.name, p.lengthunit, p.h, p.w, p.d)
}

type LHWellParams struct {
	crds    wtype.WellCoords
	vunit   string
	vol     float64
	rvol    float64
	shape   ShapeParams
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
		makeShape(&p.shape),
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
	size        wtype.Coordinates
	welltype    LHWellParams
	wellXOffset float64
	wellYOffset float64
	wellXStart  float64
	wellYStart  float64
	wellZStart  float64
}

func makeLHPlate(p *LHPlateParams, name string) *wtype.LHPlate {
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
	mfr      string
	ttype    string
	minvol   float64
	maxvol   float64
	volunit  string
	filtered bool
	shape    ShapeParams
}

func makeLHTip(p *LHTipParams) *wtype.LHTip {
	return wtype.NewLHTip(p.mfr,
		p.ttype,
		p.minvol,
		p.maxvol,
		p.volunit,
		p.filtered,
		makeShape(&p.shape))
}

type LHTipboxParams struct {
	nrows        int
	ncols        int
	size         wtype.Coordinates
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
	size       wtype.Coordinates
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
func get_not_in(a, b []string) []string {
	ret := []string{}
	for _, va := range a {
		c := false
		for _, vb := range b {
			if va == vb {
				c = true
			}
		}
		if !c {
			ret = append(ret, va)
		}
	}
	return ret
}

func compare_errors(t *testing.T, desc string, expected []string, actual []*simulator.SimulationError) {
	string_errors := make([]string, 0)
	for _, err := range actual {
		string_errors = append(string_errors, err.Error())
	}
	// maybe sort alphabetically?

	missing := get_not_in(expected, string_errors)
	extra := get_not_in(string_errors, expected)

	errs := []string{}
	for _, s := range missing {
		errs = append(errs, fmt.Sprintf("--\"%v\"", s))
	}
	for _, s := range extra {
		errs = append(errs, fmt.Sprintf("++\"%v\"", s))
	}
	if len(missing) > 0 || len(extra) > 0 {
		t.Errorf("Errors didn't match in test \"%v\":\n%s",
			desc, strings.Join(errs, "\n"))
	}
}

/*
 * ####################################### Default Types
 */

func default_lhplate_props() *LHPlateParams {
	params := LHPlateParams{
		"plate",          // platetype       string
		"test_plate_mfr", // mfr             string
		8,                // nrows           int
		12,               // ncols           int
		wtype.Coordinates{X: 127.76, Y: 85.48, Z: 25.7}, // size          float64
		LHWellParams{ // welltype
			wtype.ZeroWellCoords(), // crds            string
			"ul", // vunit           string
			200,  // vol             float64
			5,    // rvol            float64
			ShapeParams{ // shape           ShapeParams struct {
				"test_shape", // name            string
				"mm",         // lengthunit      string
				5.5,          // h               float64
				5.5,          // w               float64
				20.4,         // d               float64
			},
			wtype.VWellBottom, // bott            int
			5.5,               // xdim            float64
			5.5,               // ydim            float64
			20.4,              // zdim            float64
			1.4,               // bottomh         float64
			"mm",              // dunit           string
		},
		9.,   // wellXOffset     float64
		9.,   // wellYOffset     float64
		0.,   // wellXStart      float64
		0.,   // wellYStart      float64
		18.5, // wellZStart      float64
	}

	return &params
}

func default_lhplate(name string) *wtype.LHPlate {
	params := default_lhplate_props()
	return makeLHPlate(params, name)
}

//This plate will fill into the next door position on the robot
func wide_lhplate(name string) *wtype.LHPlate {
	params := default_lhplate_props()
	params.size.X = 300.
	return makeLHPlate(params, name)
}

func lhplate_trough_props() *LHPlateParams {
	params := LHPlateParams{
		"trough",          // platetype       string
		"test_trough_mfr", // mfr             string
		1,                 // nrows           int
		12,                // ncols           int
		wtype.Coordinates{X: 127.76, Y: 85.48, Z: 25.7}, // size          float64
		LHWellParams{ // welltype
			wtype.ZeroWellCoords(), // crds            string
			"ul",  // vunit           string
			15000, // vol             float64
			5000,  // rvol            float64
			ShapeParams{ // shape           ShapeParams struct {
				"test_shape", // name            string
				"mm",         // lengthunit      string
				8.2,          // h               float64
				72.0,         // w               float64
				41.3,         // d               float64
			},
			wtype.FlatWellBottom, // bott            int
			8.2,                  // xdim            float64
			72.0,                 // ydim            float64
			41.3,                 // zdim            float64
			4.7,                  // bottomh         float64
			"mm",                 // dunit           string
		},
		9.,   // wellXOffset     float64
		9.,   // wellYOffset     float64
		0.,   // wellXStart      float64
		30.0, // wellYStart      float64
		4.5,  // wellZStart      float64
	}

	return &params
}

func lhplate_trough12(name string) *wtype.LHPlate {
	params := lhplate_trough_props()
	plate := makeLHPlate(params, name)
	targets := []wtype.Coordinates{
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

func default_lhtipbox(name string) *wtype.LHTipbox {
	params := LHTipboxParams{
		8,  //nrows           int
		12, //ncols           int
		wtype.Coordinates{X: 127.76, Y: 85.48, Z: 60.13}, //size         float64
		"test Tipbox mfg",                                //manufacturer    string
		"tipbox",                                         //boxtype         string
		LHTipParams{ //tiptype
			"test_tip mfg",  //mfr         string
			"test_tip type", //ttype       string
			50,              //minvol      float64
			1000,            //maxvol      float64
			"ul",            //volunit     string
			false,           //filtered    bool
			ShapeParams{ // shape           ShapeParams struct {
				"test_shape", // name            string
				"mm",         // lengthunit      string
				7.3,          // h               float64
				7.3,          // w               float64
				51.2,         // d               float64
			},
		},
		LHWellParams{ // well
			wtype.ZeroWellCoords(), // crds            string
			"ul", // vunit           string
			1000, // vol             float64
			50,   // rvol            float64
			ShapeParams{ // shape           ShapeParams struct {
				"test_shape", // name            string
				"mm",         // lengthunit      string
				7.3,          // h               float64
				7.3,          // w               float64
				51.2,         // d               float64
			},
			wtype.VWellBottom, // bott            int
			7.3,               // xdim            float64
			7.3,               // ydim            float64
			51.2,              // zdim            float64
			0.0,               // bottomh         float64
			"mm",              // dunit           string
		},
		9.,  //tipxoffset      float64
		9.,  //tipyoffset      float64
		0.,  //tipxstart       float64
		0.,  //tipystart       float64
		10., //tipzstart       float64
	}

	return makeLHTipbox(&params, name)
}

func small_lhtipbox(name string) *wtype.LHTipbox {
	params := LHTipboxParams{
		8,  //nrows           int
		12, //ncols           int
		wtype.Coordinates{X: 127.76, Y: 85.48, Z: 60.13}, //size         float64
		"test Tipbox mfg",                                //manufacturer    string
		"tipbox",                                         //boxtype         string
		LHTipParams{ //tiptype
			"test_tip mfg",  //mfr         string
			"test_tip type", //ttype       string
			0,               //minvol      float64
			200,             //maxvol      float64
			"ul",            //volunit     string
			false,           //filtered    bool
			ShapeParams{ // shape           ShapeParams struct {
				"test_shape", // name            string
				"mm",         // lengthunit      string
				7.3,          // h               float64
				7.3,          // w               float64
				51.2,         // d               float64
			},
		},
		LHWellParams{ // well
			wtype.ZeroWellCoords(), // crds            string
			"ul", // vunit           string
			1000, // vol             float64
			50,   // rvol            float64
			ShapeParams{ // shape           ShapeParams struct {
				"test_shape", // name            string
				"mm",         // lengthunit      string
				7.3,          // h               float64
				7.3,          // w               float64
				51.2,         // d               float64
			},
			wtype.VWellBottom, // bott            int
			7.3,               // xdim            float64
			7.3,               // ydim            float64
			51.2,              // zdim            float64
			0.0,               // bottomh         float64
			"mm",              // dunit           string
		},
		9.,  //tipxoffset      float64
		9.,  //tipyoffset      float64
		0.,  //tipxstart       float64
		0.,  //tipystart       float64
		10., //tipzstart       float64
	}

	return makeLHTipbox(&params, name)
}

func default_lhtipwaste(name string) *wtype.LHTipwaste {
	params := LHTipwasteParams{
		700,                                             //capacity        int
		"tipwaste",                                      //typ             string
		"testTipwaste mfr",                              //mfr             string
		wtype.Coordinates{X: 127.76, Y: 85.48, Z: 92.0}, //height          float64
		LHWellParams{ // w               LHWellParams
			wtype.ZeroWellCoords(), // crds            string
			"ul",     // vunit           string
			800000.0, // vol             float64
			800000.0, // rvol            float64
			ShapeParams{ // shape           ShapeParams struct {
				"test_tipbox", // name            string
				"mm",          // lengthunit      string
				123.0,         // h               float64
				80.0,          // w               float64
				92.0,          // d               float64
			},
			wtype.VWellBottom, // bott            int
			123.0,             // xdim            float64
			80.0,              // ydim            float64
			92.0,              // zdim            float64
			0.0,               // bottomh         float64
			"mm",              // dunit           string
		},
		49.5, //wellxstart      float64
		31.5, //wellystart      float64
		0.0,  //wellzstart      float64
	}
	return makeLHTipWaste(&params, name)
}

func default_lhproperties() *liquidhandling.LHProperties {
	valid_props := LHPropertiesParams{
		"Device Name",
		"Device Manufacturer",
		[]LayoutParams{
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
		[]HeadParams{
			{
				"Head0 Name",
				"Head0 Manufacturer",
				ChannelParams{
					"Head0 ChannelParams",     //Name
					"Head0 Platform",          //Platform
					UnitParams{0.1, "ul"},     //min volume
					UnitParams{1., "ml"},      //max volume
					UnitParams{0.1, "ml/min"}, //min flowrate
					UnitParams{10., "ml/min"}, //max flowrate
					8,     //multi
					false, //independent
					0,     //orientation
					0,     //head
				},
				AdaptorParams{
					"Head0 Adaptor",
					"Head0 Adaptor Manufacturer",
					ChannelParams{
						"Head0 Adaptor ChannelParams", //Name
						"Head0 Adaptor Platform",      //Platform
						UnitParams{0.1, "ul"},         //min volume
						UnitParams{1., "ml"},          //max volume
						UnitParams{0.1, "ml/min"},     //min flowrate
						UnitParams{10., "ml/min"},     //max flowrate
						8,     //multi
						false, //independent
						0,     //orientation
						0,     //head
					},
				},
			},
		},
		[]string{"tipbox_1", "tipbox_2"}, //Tip_preferences
		[]string{"input_1", "input_2"},   //Input_preferences
		[]string{"output_1", "output_2"}, //Output_preferences
		[]string{"tipwaste"},             //Tipwaste_preferences
		[]string{"wash"},                 //Wash_preferences
		[]string{"waste"},                //Waste_preferences
	}

	return makeLHProperties(&valid_props)
}

func independent_lhproperties() *liquidhandling.LHProperties {
	valid_props := LHPropertiesParams{
		"Device Name",
		"Device Manufacturer",
		[]LayoutParams{
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
		[]HeadParams{
			{
				"Head0 Name",
				"Head0 Manufacturer",
				ChannelParams{
					"Head0 ChannelParams",     //Name
					"Head0 Platform",          //Platform
					UnitParams{0.1, "ul"},     //min volume
					UnitParams{1., "ml"},      //max volume
					UnitParams{0.1, "ml/min"}, //min flowrate
					UnitParams{10., "ml/min"}, //max flowrate
					8,    //multi
					true, //independent
					0,    //orientation
					0,    //head
				},
				AdaptorParams{
					"Head0 Adaptor",
					"Head0 Adaptor Manufacturer",
					ChannelParams{
						"Head0 Adaptor ChannelParams", //Name
						"Head0 Adaptor Platform",      //Platform
						UnitParams{0.1, "ul"},         //min volume
						UnitParams{1., "ml"},          //max volume
						UnitParams{0.1, "ml/min"},     //min flowrate
						UnitParams{10., "ml/min"},     //max flowrate
						8,    //multi
						true, //independent
						0,    //orientation
						0,    //head
					},
				},
			},
		},
		[]string{"tipbox_1", "tipbox_2"}, //Tip_preferences
		[]string{"input_1", "input_2"},   //Input_preferences
		[]string{"output_1", "output_2"}, //Output_preferences
		[]string{"tipwaste"},             //Tipwaste_preferences
		[]string{"wash"},                 //Wash_preferences
		[]string{"waste"},                //Waste_preferences
	}

	return makeLHProperties(&valid_props)
}

/* -- remove for linting
func default_vlh() *lh.VirtualLiquidHandler {
	vlh := lh.NewVirtualLiquidHandler(default_lhproperties(), nil)
	return vlh
}
*/

/*
 * ######################################## InstructionParams
 */

type TestRobotInstruction interface {
	Apply(*lh.VirtualLiquidHandler)
}

//Initialize
type Initialize struct{}

func (self *Initialize) Apply(vlh *lh.VirtualLiquidHandler) {
	vlh.Initialize()
}

//Finalize
type Finalize struct{}

func (self *Finalize) Apply(vlh *lh.VirtualLiquidHandler) {
	vlh.Finalize()
}

//SetPipetteSpeed
type SetPipetteSpeed struct {
	head    int
	channel int
	speed   float64
}

func (self *SetPipetteSpeed) Apply(vlh *lh.VirtualLiquidHandler) {
	vlh.SetPipetteSpeed(self.head, self.channel, self.speed)
}

//AddPlateTo
type AddPlateTo struct {
	position string
	plate    interface{}
	name     string
}

func (self *AddPlateTo) Apply(vlh *lh.VirtualLiquidHandler) {
	vlh.AddPlateTo(self.position, self.plate, self.name)
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

func (self *LoadTips) Apply(vlh *lh.VirtualLiquidHandler) {
	vlh.LoadTips(self.channels, self.head, self.multi, self.platetype, self.position, self.well)
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

func (self *UnloadTips) Apply(vlh *lh.VirtualLiquidHandler) {
	vlh.UnloadTips(self.channels, self.head, self.multi, self.platetype, self.position, self.well)
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

func (self *Move) Apply(vlh *lh.VirtualLiquidHandler) {
	vlh.Move(self.deckposition, self.wellcoords, self.reference, self.offsetX, self.offsetY, self.offsetZ, self.plate_type, self.head)
}

//Aspirate
type Aspirate struct {
	volume     []float64
	overstroke []bool
	head       int
	multi      int
	platetype  []string
	what       []string
	llf        []bool
}

func (self *Aspirate) Apply(vlh *lh.VirtualLiquidHandler) {
	vlh.Aspirate(self.volume, self.overstroke, self.head, self.multi,
		self.platetype, self.what, self.llf)
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

func (self *Dispense) Apply(vlh *lh.VirtualLiquidHandler) {
	vlh.Dispense(self.volume, self.blowout, self.head, self.multi,
		self.platetype, self.what, self.llf)
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

func (self *Mix) Apply(vlh *lh.VirtualLiquidHandler) {
	vlh.Mix(self.head, self.volume, self.platetype, self.cycles,
		self.multi, self.what, self.blowout)
}

/*
 * ######################################## Setup
 */

type SetupFn func(*lh.VirtualLiquidHandler)

func removeTipboxTips(tipbox_loc string, wells []string) *SetupFn {
	var ret SetupFn = func(vlh *lh.VirtualLiquidHandler) {
		tipbox := vlh.GetObjectAt(tipbox_loc).(*wtype.LHTipbox)
		for _, well := range wells {
			wc := wtype.MakeWellCoords(well)
			tipbox.RemoveTip(wc)
		}
	}
	return &ret
}

func preloadAdaptorTips(head int, tipbox_loc string, channels []int) *SetupFn {
	var ret SetupFn = func(vlh *lh.VirtualLiquidHandler) {
		adaptor := vlh.GetAdaptorState(head)
		tipbox := vlh.GetObjectAt(tipbox_loc).(*wtype.LHTipbox)

		for _, ch := range channels {
			adaptor.GetChannel(ch).LoadTip(tipbox.Tiptype.Dup())
		}
	}
	return &ret
}

func getLHComponent(what string, vol_ul float64) *wtype.LHComponent {
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
	var ret SetupFn = func(vlh *lh.VirtualLiquidHandler) {
		adaptor := vlh.GetAdaptorState(head)
		tipbox := vlh.GetObjectAt(tipbox_loc).(*wtype.LHTipbox)
		tip := tipbox.Tiptype.Dup()
		c := getLHComponent(what, volume)
		tip.AddComponent(c)

		for _, ch := range channels {
			adaptor.GetChannel(ch).LoadTip(tip.Dup())
		}
	}
	return &ret
}

/* -- remove for linting
func fillTipwaste(tipwaste_loc string, count int) *SetupFn {
	var ret SetupFn = func(vlh *lh.VirtualLiquidHandler) {
		tipwaste := vlh.GetObjectAt(tipwaste_loc).(*wtype.LHTipwaste)
		tipwaste.Contents += count
	}
	return &ret
}
*/

func prefillWells(plate_loc string, wells_to_fill []string, liquid_name string, volume float64) *SetupFn {
	var ret SetupFn = func(vlh *lh.VirtualLiquidHandler) {
		plate := vlh.GetObjectAt(plate_loc).(*wtype.LHPlate)
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

	var ret SetupFn = func(vlh *lh.VirtualLiquidHandler) {
		vlh.Move(s_dp, s_wc, s_rf, s_ox, s_oy, s_oz, s_pt, p.Head)
	}

	return &ret
}

/*
 * ######################################## Assertions (about the final state)
 */

type AssertionFn func(string, *testing.T, *lh.VirtualLiquidHandler)

//tipboxAssertion assert that the tipbox has tips missing in the given locations only
func tipboxAssertion(tipbox_loc string, missing_tips []string) *AssertionFn {
	var ret AssertionFn = func(name string, t *testing.T, vlh *lh.VirtualLiquidHandler) {
		mmissing_tips := make(map[string]bool)
		for _, tl := range missing_tips {
			mmissing_tips[tl] = true
		}

		if tipbox, ok := vlh.GetObjectAt(tipbox_loc).(*wtype.LHTipbox); !ok {
			t.Errorf("TipboxAssertion failed in \"%s\", no Tipbox found at \"%s\"", name, tipbox_loc)
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
				t.Errorf("TipboxAssertion failed in test \"%s\", tipbox at \"%s\":\n%s", name, tipbox_loc, strings.Join(errors, "\n"))
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
	var ret AssertionFn = func(name string, t *testing.T, vlh *lh.VirtualLiquidHandler) {
		mtips := make(map[int]bool)
		for _, td := range tips {
			mtips[td.channel] = true
		}

		adaptor := vlh.GetAdaptorState(head)
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
			t.Errorf("AdaptorAssertion failed in test \"%s\", Head%v:\n%s", name, head, strings.Join(errors, "\n"))
		}
	}
	return &ret
}

//adaptorPositionAssertion assert that the adaptor has tips in the given positions
func positionAssertion(head int, origin wtype.Coordinates) *AssertionFn {
	var ret AssertionFn = func(name string, t *testing.T, vlh *lh.VirtualLiquidHandler) {
		adaptor := vlh.GetAdaptorState(head)
		or := adaptor.GetChannel(0).GetAbsolutePosition()
		if or.X != origin.X || or.Y != origin.Y || or.Z != origin.Z {
			t.Errorf("PositionAssertion failed in \"%s\", adaptor should be at %s, was actually at %s", name, origin, or)
		}
	}
	return &ret
}

//tipwasteAssertion assert the number of tips which should be in the tipwaste
func tipwasteAssertion(tipwaste_loc string, expected_contents int) *AssertionFn {
	var ret AssertionFn = func(name string, t *testing.T, vlh *lh.VirtualLiquidHandler) {
		if tipwaste, ok := vlh.GetObjectAt(tipwaste_loc).(*wtype.LHTipwaste); !ok {
			t.Errorf("TipWasteAssertion failed in \"%s\", no Tipwaste found at %s", name, tipwaste_loc)
		} else {
			if tipwaste.Contents != expected_contents {
				t.Errorf("TipwasteAssertion failed in test \"%s\" at location %s: expected %v tips, got %v",
					name, tipwaste_loc, expected_contents, tipwaste.Contents)
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
	var ret AssertionFn = func(name string, t *testing.T, vlh *lh.VirtualLiquidHandler) {
		m := map[string]bool{}
		plate := vlh.GetObjectAt(plate_loc).(*wtype.LHPlate)
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
			t.Errorf("plateAssertion failed in test \"%s\", errors were:\n%s", name, strings.Join(errs, "\n"))
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

func (self *SimulatorTest) run(t *testing.T) {

	if self.Props == nil {
		self.Props = default_lhproperties()
	}
	vlh := lh.NewVirtualLiquidHandler(self.Props, nil)

	//do setup
	if self.Setup != nil {
		for _, setup_fn := range self.Setup {
			(*setup_fn)(vlh)
		}
	}

	//run the instructions
	if self.Instructions != nil {
		for _, inst := range self.Instructions {
			inst.Apply(vlh)
		}
	}

	//check errors
	if self.ExpectedErrors != nil {
		compare_errors(t, self.Name, self.ExpectedErrors, vlh.GetErrors())
	} else {
		compare_errors(t, self.Name, []string{}, vlh.GetErrors())
	}

	//check assertions
	if self.Assertions != nil {
		for _, a := range self.Assertions {
			(*a)(self.Name, t, vlh)
		}
	}
}
