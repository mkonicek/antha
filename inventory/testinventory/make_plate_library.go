// Copyright (C) 2017 The Antha authors. All rights reserved.
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

package testinventory

import (
	"encoding/json"
	"strings"

	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wutil"
)

//var commonwelltypes

var platespecificoffset = map[string]float64{
	"pcrplate_skirted": gilsonoffsetpcrplate,
	"greiner384":       gilsonoffsetgreiner,
	"costar48well":     3.0,
	"Nuncon12well":     11.0, // this must be wrong!! check z start without riser properly
	"Nuncon12wellAgar": 11.0, // this must be wrong!! check z start without riser properly
	"VWR12well":        3.0,
}

// function to check if a platename already contains a riser
func containsRiser(plate *wtype.LHPlate) bool {
	for _, dev := range defaultDevices {
		for _, synonym := range dev.GetSynonyms() {
			if strings.Contains(plate.Type, "_"+synonym) {
				return true
			}
		}
	}

	return false
}

func addRiser(plate *wtype.LHPlate, riser device) (plates []*wtype.LHPlate) {
	if containsRiser(plate) {
		return
	}

	for _, risername := range riser.GetSynonyms() {
		var dontaddrisertothisplate bool

		newplate := plate.Dup()
		riserheight := riser.GetHeightInmm()
		if offset, found := platespecificoffset[plate.Type]; found {
			riserheight = riserheight - offset
		}
		newplate.WellZStart = plate.WellZStart + riserheight
		newname := plate.Type + "_" + risername
		newplate.Type = newname
		if riser.GetConstraints() != nil {
			// duplicate well before adding constraint to prevent applying
			// constraint to all common &Welltype on other plates

			for device, allowedpositions := range riser.GetConstraints() {
				newwell := newplate.Welltype.Dup()
				newplate.Welltype = newwell
				_, ok := newwell.Extra[device]
				if !ok {
					newplate.SetConstrained(device, allowedpositions)
				} else {
					dontaddrisertothisplate = true
				}
			}
		}

		if !dontaddrisertothisplate {
			plates = append(plates, newplate)
		}
	}

	return
}

func addAllDevices(plates []*wtype.LHPlate) (ret []*wtype.LHPlate) {
	for _, plate := range plates {
		for _, dev := range defaultDevices {
			ret = append(ret, addRiser(plate, dev)...)
		}
	}
	return
}

func makePlates() (plates []*wtype.LHPlate) {
	plates = makeBasicPlates()
	additional := addAllDevices(plates)
	return append(plates, additional...)
}

func makeBasicPlates() (plates []*wtype.LHPlate) {
	// deep square well 96
	swshp := wtype.NewShape("box", "mm", 8.2, 8.2, 41.3)
	deepsquarewell := wtype.NewLHWell("DSW96", "", "", "ul", 2000, 420, swshp, wtype.LHWBV, 8.2, 8.2, 41.3, 4.7, "mm")
	plate := wtype.NewLHPlate("DSW96", "Unknown", 8, 12, 44.1, "mm", deepsquarewell, 9, 9, 0.0, 0.0, valueformaxheadtonotintoDSWplatewithp20tips)
	plates = append(plates, plate)

	// IDT/ABgene 1.2 ml storage plate AB0564
	idtshp := wtype.NewShape("cylinder", "mm", 7, 7, 39.35)
	idtroundwell96 := wtype.NewLHWell("IDT96", "", "", "ul", 1200, 100, idtshp, wtype.LHWBU, 7, 7, 39.35, 3, "mm")
	plate = wtype.NewLHPlate("IDT96", "Unknown", 8, 12, 42.5, "mm", idtroundwell96, 9, 9, 0, 0, 3)
	plates = append(plates, plate)

	//4 column reservoir plate Phenix Research Products RRI3051; Fisher cat# NC0336913
	fourcolumnshp := wtype.NewShape("box", "mm", 26, 71, 42)
	fourcolumnwell := wtype.NewLHWell("FourColumnWell", "", "", "ul", 73000, 3000, fourcolumnshp, wtype.LHWBV, 26, 71, 42, 2, "mm")
	plate = wtype.NewLHPlate("FourColumnReservoir", "Unknown", 1, 4, 44, "mm", fourcolumnwell, 26, 1, 9.5, 31, 1) //WellYStart is not accurate, but would not visualise correctly unless set to this value, cant diagnose
	plates = append(plates, plate)

	// 24 well deep square well plate on riser

	bottomtype := wtype.LHWBV // 0 = flat, 2 = v shaped
	xdim := 16.8
	ydim := 16.8
	zdim := 41.3
	bottomh := 4.7

	wellcapacityinwelltypeunit := 11000.0
	welltypeunit := "ul"
	wellsperrow := 6
	wellspercolumn := 4
	residualvol := 650.0 // assume in ul

	wellxoffset := 18.0 // centre of well to centre of neighbouring well in x direction
	wellyoffset := 18.0 //centre of well to centre of neighbouring well in y direction
	xstart := 4.5       // distance from top left side of plate to first well
	ystart := 4.5       // distance from top left side of plate to first well
	zstart := -1.0      // offset of bottom of deck to bottom of well (this includes agar estimate)

	heightinmm := 44.1

	squarewell := wtype.NewShape("box", "mm", xdim, ydim, zdim)
	squarewell24 := wtype.NewLHWell("24DSW", "", "", welltypeunit, wellcapacityinwelltypeunit, residualvol, squarewell, bottomtype, xdim, ydim, zdim, bottomh, "mm")
	plate = wtype.NewLHPlate("DSW24", "Unknown", wellspercolumn, wellsperrow, heightinmm, "mm", squarewell24, wellxoffset, wellyoffset, xstart, ystart, zstart)

	plates = append(plates, plate)

	// shallow round well flat bottom 96
	rwshp := wtype.NewShape("cylinder", "mm", 8.2, 8.2, 11)
	roundwell96 := wtype.NewLHWell("SRWFB96", "", "", "ul", 500, 10, rwshp, 0, 8.2, 8.2, 11, 1.0, "mm")
	plate = wtype.NewLHPlate("SRWFB96", "Unknown", 8, 12, 15, "mm", roundwell96, 9, 9, 0.0, 0.0, 1.0)
	plates = append(plates, plate)

	// deep well strip trough 12
	stshp := wtype.NewShape("box", "mm", 8.2, 72, 41.3)
	trough12 := wtype.NewLHWell("DWST12", "", "", "ul", 15000, 5000, stshp, wtype.LHWBV, 8.2, 72, 41.3, 4.7, "mm")
	plate = wtype.NewLHPlate("DWST12", "Unknown", 1, 12, 44.1, "mm", trough12, 9, 9, 0, 30.0, valueformaxheadtonotintoDSWplatewithp20tips)
	plates = append(plates, plate)

	// deep well strip trough 8
	stshp8 := wtype.NewShape("box", "mm", 115.0, 8.2, 41.3)
	trough8 := wtype.NewLHWell("DWST8", "", "", "ul", 24000, 1000, stshp8, wtype.LHWBV, 115, 8.2, 41.3, 4.7, "mm")
	plate = wtype.NewLHPlate("DWST8", "Unknown", 8, 1, 44.1, "mm", trough8, 9, 9, 49.5, 0.0, 0.0)
	plates = append(plates, plate)

	// deep well reservoir
	rshp := wtype.NewShape("box", "mm", 115.0, 72.0, 41.3)
	singlewelltrough := wtype.NewLHWell("DWR1", "", "", "ul", 300000, 20000, rshp, wtype.LHWBV, 115, 72, 41.3, 4.7, "mm")
	plate = wtype.NewLHPlate("DWR1", "Unknown", 1, 1, 44.1, "mm", singlewelltrough, 9, 9, 49.5, 0.0, 0.0)
	plates = append(plates, plate)

	// 250ml box reservoir (working vol estimated to be 100ml to prevent spillage on moving decks)
	reservoirbox := wtype.NewShape("box", "mm", 121, 80, 40) // 39?
	welltypereservoir := wtype.NewLHWell("Reservoir", "", "", "ul", 100000, 10000, reservoirbox, wtype.LHWBFLAT, 121, 80, 40, 3, "mm")
	plate = wtype.NewLHPlate("reservoir", "unknown", 1, 1, 40, "mm", welltypereservoir, 1, 1, 0.0, 0.0, 0.0)
	plates = append(plates, plate)

	// well area function
	// -- determined empirically since inverse cubic was giving us some numerical issues
	areaf := wutil.Quartic{A: -3.3317851312e-09, B: 0.00000225834467, C: -0.0006305492472, D: 0.1328156706978, E: 0}
	afb, _ := json.Marshal(areaf)
	afs := string(afb)

	// pcr plate with cooler
	cone := wtype.NewShape("cylinder", "mm", 5.5, 5.5, 15)

	pcrplatewell := wtype.NewLHWell("pcrplate", "", "", "ul", 200, 5, cone, wtype.LHWBU, 5.5, 5.5, 15, 1.4, "mm")
	pcrplatewell.SetAfVFunc(afs)

	plate = wtype.NewLHPlate("pcrplate_with_cooler", "Unknown", 8, 12, 15.5, "mm", pcrplatewell, 9, 9, 0.0, 0.0, coolerheight+0.5)
	plates = append(plates, plate)

	// pcr plate with isofreeze_cooler
	plate = wtype.NewLHPlate("pcrplate_with_isofreeze_cooler", "Unknown", 8, 12, 15.5, "mm", pcrplatewell, 9, 9, 0.0, 0.0, isofreezecoolerheight)
	plates = append(plates, plate)

	// pcr plate skirted with isofreeze_cooler (to be used only with transformations (into 10-20ul) as plate not fully secured in the cooler)
	plate = wtype.NewLHPlate("pcrplate_skirted_with_isofreeze_cooler", "Unknown", 8, 12, 15.5, "mm", pcrplatewell, 9, 9, 0.0, 0.0, isofreezecoolerheight+2.0)
	plates = append(plates, plate)

	// pcr plate with 496rack

	plate = wtype.NewLHPlate("pcrplate_with_496rack", "Unknown", 8, 12, 15.5, "mm", pcrplatewell, 9, 9, 0.0, 0.0, pcrtuberack496)
	plates = append(plates, plate)

	// pcr plate semi-skirted with 496rack

	plate = wtype.NewLHPlate("pcrplate_semi_skirted_with_496rack", "Unknown", 8, 12, 15.5, "mm", pcrplatewell, 9, 9, 0.0, 0.0, pcrtuberack496+1.0)
	plates = append(plates, plate)

	// 0.2ml strip tubes with 496rack

	plate = wtype.NewLHPlate("strip_tubes_0.2ml_with_496rack", "Unknown", 8, 12, 15.5, "mm", pcrplatewell, 9, 9, 0.0, 0.0, pcrtuberack496-2.5)
	plates = append(plates, plate)

	// pcr plate skirted
	plate = wtype.NewLHPlate("pcrplate_skirted", "Unknown", 8, 12, 15.5, "mm", pcrplatewell, 9, 9, 0.0, 0.0, 0.636)
	plates = append(plates, plate)

	pcrplatewellinc := wtype.NewLHWell("pcrplate", "", "", "ul", 200, 5, cone, wtype.LHWBU, 5.5, 5.5, 1.55, 1.4, "mm")
	pcrplatewellinc.SetAfVFunc(afs)

	// falcon 6 well plate with Agar flat bottom with 4ml per well

	bottomtype = wtype.LHWBFLAT
	xdim = 37.0
	ydim = 37.0
	zdim = 20.0
	bottomh = 9.0 //(this includes agar estimate)

	wellxoffset = 39.0 // centre of well to centre of neighbouring well in x direction
	wellyoffset = 39.0 //centre of well to centre of neighbouring well in y direction
	xstart = 5.0       // distance from top left side of plate to first well
	ystart = 5.0       // distance from top left side of plate to first well
	zstart = 9.0       // offset of bottom of deck to bottom of well (this includes agar estimate)

	wellsperrow = 3
	wellspercolumn = 2
	heightinmm = 20.0

	circle := wtype.NewShape("cylinder", "mm", 37, 37, 20)
	welltype6well := wtype.NewLHWell("falcon6well", "", "", "ul", 4000, 1, circle, bottomtype, xdim, ydim, zdim, bottomh, "mm")

	plate = wtype.NewLHPlate("falcon6wellAgar", "Unknown", wellspercolumn, wellsperrow, heightinmm, "mm", welltype6well, wellxoffset, wellyoffset, xstart, ystart, zstart)
	plate.DeclareSpecial() // Do this for racks, other very unusual plate types
	plates = append(plates, plate)

	// Costar 48 well plate flat bottom

	bottomtype = wtype.LHWBFLAT
	xdim = 11.0
	ydim = 11.0
	zdim = 19.0
	bottomh = 3.0

	wellxoffset = 13.0 // centre of well to centre of neighbouring well in x direction
	wellyoffset = 13.0 //centre of well to centre of neighbouring well in y direction
	xstart = 3.0       // distance from top left side of plate to first well
	ystart = -1.0      // distance from top left side of plate to first well
	zstart = 2.0       // offset of bottom of deck to bottom of well (this includes agar estimate)

	wellsperrow = 8
	wellspercolumn = 6
	heightinmm = 20.0

	circle = wtype.NewShape("cylinder", "mm", xdim, ydim, zdim)
	welltypecostar48 := wtype.NewLHWell("costar48well", "", "", "ul", 1000, 100, circle, bottomtype, xdim, ydim, zdim, bottomh, "mm")

	plate = wtype.NewLHPlate("costar48well", "Unknown", wellspercolumn, wellsperrow, heightinmm, "mm", welltypecostar48, wellxoffset, wellyoffset, xstart, ystart, zstart)
	plate.DeclareSpecial() // Do this for racks, other very unusual plate types
	plates = append(plates, plate)

	// Block Kombi 2ml
	eppy := wtype.NewShape("cylinder", "mm", 8.2, 8.2, 45)

	wellxoffset = 18.0 // centre of well to centre of neighbouring well in x direction
	wellyoffset = 18.0 //centre of well to centre of neighbouring well in y direction
	xstart = 5.0       // distance from top left side of plate to first well
	ystart = 10.0      // distance from top left side of plate to first well
	zstart = 6.0       // offset of bottom of deck to bottom of well

	welltype2mleppy := wtype.NewLHWell("2mlEpp", "", "", "ul", 2000, 25, eppy, wtype.LHWBV, 8.2, 8.2, 45, 4.7, "mm")

	plate = wtype.NewLHPlate("Kombi2mlEpp", "Unknown", 4, 2, 45, "mm", welltype2mleppy, wellxoffset, wellyoffset, xstart, ystart, zstart)
	plates = append(plates, plate)

	plate.DeclareSpecial() // Do this for racks, other very unusual plate types

	// Eppendorfrack 425 for 2ml tubes

	wellxoffset = 18.0 // centre of well to centre of neighbouring well in x direction
	wellyoffset = 18.0 //centre of well to centre of neighbouring well in y direction
	xstart = 4.5       // distance from top left side of plate to first well
	ystart = 5.0       // distance from top left side of plate to first well
	zstart = 5.0       // offset of bottom of deck to bottom of well

	plate = wtype.NewLHPlate("eppendorfrack425_2ml", "Unknown", 4, 6, 45, "mm", welltype2mleppy, wellxoffset, wellyoffset, xstart, ystart, zstart)
	plate.DeclareSpecial() // Do this for racks, other very unusual plate types
	plates = append(plates, plate)

	// Eppendorfrack 425 for 1.5ml tubes

	wellxoffset = 18.0 // centre of well to centre of neighbouring well in x direction
	wellyoffset = 18.0 //centre of well to centre of neighbouring well in y direction
	xstart = 4.5       // distance from top left side of plate to first well
	ystart = 5.0       // distance from top left side of plate to first well
	zstart = 5.0       // offset of bottom of deck to bottom of well

	welltypesmallereppy := wtype.NewLHWell("1.5mlEpp", "", "", "ul", 1500, 50, eppy, wtype.LHWBV, 8.2, 8.2, 45, 4.7, "mm")

	plate = wtype.NewLHPlate("eppendorfrack425_1.5ml", "Unknown", 4, 6, 45, "mm", welltypesmallereppy, wellxoffset, wellyoffset, xstart, ystart, zstart)
	plate.DeclareSpecial() // Do this for racks, other very unusual plate types
	plates = append(plates, plate)

	// Eppendorfrack 424 with lid holders and using 2ml tubes

	wellxoffset = 36.0 // centre of well to centre of neighbouring well in x direction
	wellyoffset = 18.0 //centre of well to centre of neighbouring well in y direction
	xstart = 14.0      // distance from top left side of plate to first well
	ystart = 5.0       // distance from top left side of plate to first well
	zstart = 5.0       // offset of bottom of deck to bottom of well

	plate = wtype.NewLHPlate("eppendorfrack424_2ml_lidholder", "Unknown", 4, 3, 45, "mm", welltype2mleppy, wellxoffset, wellyoffset, xstart, ystart, zstart)
	plate.DeclareSpecial() // Do this for racks, other very unusual plate types
	plates = append(plates, plate)

	// Eppendorfrack 424 with lid holders and using 1.5ml tubes

	wellxoffset = 36.0 // centre of well to centre of neighbouring well in x direction
	wellyoffset = 18.0 //centre of well to centre of neighbouring well in y direction
	xstart = 14.0      // distance from top left side of plate to first well
	ystart = 5.0       // distance from top left side of plate to first well
	zstart = 9         // offset of bottom of deck to bottom of well
	zstart = 4.5       // offset of bottom of deck to bottom of well

	plate = wtype.NewLHPlate("eppendorfrack424_1.5ml_lidholder", "Unknown", 4, 3, 45, "mm", welltypesmallereppy, wellxoffset, wellyoffset, xstart, ystart, zstart)
	plate.DeclareSpecial() // Do this for racks, other very unusual plate types
	plates = append(plates, plate)

	// greiner 384 well plate flat bottom

	bottomtype = wtype.LHWBFLAT
	xdim = 4.0
	ydim = 4.0
	zdim = 12.0
	bottomh = 1.0

	wellxoffset = 4.5 // centre of well to centre of neighbouring well in x direction
	wellyoffset = 4.5 //centre of well to centre of neighbouring well in y direction
	xstart = -2.5     // distance from top left side of plate to first well
	ystart = -2.5     // distance from top left side of plate to first well
	zstart = 2.5      // offset of bottom of deck to bottom of well

	square := wtype.NewShape("box", "mm", xdim, ydim, zdim)
	//func NewLHWell(platetype, plateid, crds, vunit string, vol, rvol float64, shape *Shape, bott int, xdim, ydim, zdim, bottomh float64, dunit string) *LHWell {
	welltype384 := wtype.NewLHWell("384flat", "", "", "ul", 125, 10, square, bottomtype, xdim, ydim, zdim, bottomh, "mm")

	//func NewLHPlate(platetype, mfr string, nrows, ncols int, height float64, hunit string, welltype *LHWell, wellXOffset, wellYOffset, wellXStart, wellYStart, wellZStart float64) *LHPlate {
	plate = wtype.NewLHPlate("greiner384", "Unknown", 16, 24, 14, "mm", welltype384, wellxoffset, wellyoffset, xstart, ystart, zstart)
	plates = append(plates, plate)

	// NUNC 1536 well plate flat bottom on riser

	bottomtype = wtype.LHWBFLAT
	xdim = 2.0 // of well
	ydim = 2.0
	zdim = 7.0
	bottomh = 0.5

	wellxoffset = 2.25 // centre of well to centre of neighbouring well in x direction
	wellyoffset = 2.25 //centre of well to centre of neighbouring well in y direction
	xstart = -2.5      // distance from top left side of plate to first well
	ystart = -2.5      // distance from top left side of plate to first well
	zstart = 2         // offset of bottom of deck to bottom of well

	square1536 := wtype.NewShape("box", "mm", xdim, ydim, zdim)
	//func NewLHWell(platetype, plateid, crds, vunit string, vol, rvol float64, shape *Shape, bott int, xdim, ydim, zdim, bottomh float64, dunit string) *LHWell {
	welltype1536 := wtype.NewLHWell("1536flat", "", "", "ul", 13, 2, square1536, bottomtype, xdim, ydim, zdim, bottomh, "mm")

	plate = wtype.NewLHPlate("nunc1536", "Unknown", 32, 48, 7, "mm", welltype1536, wellxoffset, wellyoffset, xstart, ystart, zstart)
	plates = append(plates, plate)




	// Onewell SBS format Agarplate with colonies on riser (50ml agar) high res

	bottomtype = wtype.LHWBFLAT
	xdim = 2.0 // of well
	ydim = 2.0
	zdim = 7.0
	bottomh = 0.5

	wellxoffset = 2.25  // centre of well to centre of neighbouring well in x direction
	wellyoffset = 2.250 //centre of well to centre of neighbouring well in y direction
	xstart = -2.5       // distance from top left side of plate to first well
	ystart = -2.5       // distance from top left side of plate to first well
	zstart = 3          // offset of bottom of deck to bottom of well

	// greiner one well with 50ml of agar in
	plate = wtype.NewLHPlate("Agarplateforpicking1536", "Unknown", 32, 48, 7, "mm", welltype1536, wellxoffset, wellyoffset, xstart, ystart, zstart)

	plates = append(plates, plate)

	// Onewell SBS format Agarplate with colonies on riser (50ml agar) low res

	bottomtype = wtype.LHWBFLAT
	xdim = 4.0
	ydim = 4.0
	zdim = 14.0
	bottomh = 1.0

	wellxoffset = 4.5 // centre of well to centre of neighbouring well in x direction
	wellyoffset = 4.5 //centre of well to centre of neighbouring well in y direction
	xstart = -2.5     // distance from top left side of plate to first well
	ystart = -2.5     // distance from top left side of plate to first well
	zstart = 3.5      //5.5 // offset of bottom of deck to bottom of well

	// greiner one well with 50ml of agar in
	plate = wtype.NewLHPlate("Agarplateforpicking384", "Unknown", 16, 24, 14, "mm", welltype384, wellxoffset, wellyoffset, xstart, ystart, zstart)
	plates = append(plates, plate)

	// Onewell SBS format Agarplate with colonies on shallowriser (50ml agar) low res

	bottomtype = wtype.LHWBFLAT
	xdim = 4.0
	ydim = 4.0
	zdim = 14.0
	bottomh = 1.0

	wellxoffset = 4.5                     // centre of well to centre of neighbouring well in x direction
	wellyoffset = 4.5                     //centre of well to centre of neighbouring well in y direction
	xstart = -2.5                         // distance from top left side of plate to first well
	ystart = -2.5                         // distance from top left side of plate to first well
	zstart = shallowriserheightinmm + 5.5 // offset of bottom of deck to bottom of well

	// Onewell SBS format Agarplate with colonies on riser (30ml agar) low res

	zstart = 1 // offset of bottom of deck to bottom of well

	plate = wtype.NewLHPlate("30mlAgarplateforpicking384", "Unknown", 16, 24, 14, "mm", welltype384, wellxoffset, wellyoffset, xstart, ystart, zstart)
	plates = append(plates, plate)

	// Onewell SBS format Agarplate with colonies on riser (50ml agar) med res

	bottomtype = wtype.LHWBFLAT
	xdim = 3.0
	ydim = 3.0
	zdim = 14.0
	bottomh = 1.0

	wellxoffset = 3.1 // centre of well to centre of neighbouring well in x direction
	wellyoffset = 3.1 //centre of well to centre of neighbouring well in y direction
	xstart = -2.5     // distance from top left side of plate to first well
	ystart = -2.5     // distance from top left side of plate to first well
	zstart = 3.5      //5.5 // offset of bottom of deck to bottom of well

	square768 := wtype.NewShape("box", "mm", xdim, ydim, zdim)
	//func NewLHWell(platetype, plateid, crds, vunit string, vol, rvol float64, shape *Shape, bott int, xdim, ydim, zdim, bottomh float64, dunit string) *LHWell {
	welltype768 := wtype.NewLHWell("768flat", "", "", "ul", 31.25, 5, square768, bottomtype, xdim, ydim, zdim, bottomh, "mm")

	// greiner one well with 50ml of agar in
	plate = wtype.NewLHPlate("Agarplateforpicking768", "Unknown", 24, 32, 14, "mm", welltype768, wellxoffset, wellyoffset, xstart, ystart, zstart)
	plates = append(plates, plate)

	// Onewell SBS format Agarplate with colonies on riser (50ml agar) low res with 96 well map

	bottomtype = wtype.LHWBFLAT
	xdim = 5.5
	ydim = 5.5
	zdim = 15
	bottomh = 1.4

	wellxoffset = 9 // centre of well to centre of neighbouring well in x direction
	wellyoffset = 9 //centre of well to centre of neighbouring well in y direction
	xstart = 0      // distance from top left side of plate to first well
	ystart = 0      // distance from top left side of plate to first well
	zstart = 3.5    //5.5 // offset of bottom of deck to bottom of well

	// greiner one well with 50ml of agar in

	pcrplatewellforpicking := wtype.NewLHWell("pcrplate", "", "", "ul", 5, 1, cone, wtype.LHWBU, xdim, ydim, zdim, bottomh, "mm")

	plate = wtype.NewLHPlate("Agarplateforpicking96", "Unknown", 8, 12, 14, "mm", pcrplatewellforpicking, wellxoffset, wellyoffset, xstart, ystart, zstart)
	plates = append(plates, plate)

	// Onewell SBS format Agarplate with colonies on riser (50ml agar) low res with 48 well map

	bottomtype = wtype.LHWBFLAT
	xdim = 11.0
	ydim = 11.0
	zdim = 19.0
	bottomh = 1.0

	wellxoffset = 13.0 // centre of well to centre of neighbouring well in x direction
	wellyoffset = 13.0 //centre of well to centre of neighbouring well in y direction
	xstart = 3.0       // distance from top left side of plate to first well
	ystart = -1.0      // distance from top left side of plate to first well
	zstart = 3.5       //5.5 // offset of bottom of deck to bottom of well

	// greiner one well with 50ml of agar in

	welltypecostar48forpicking := wtype.NewLHWell("costar48well", "", "", "ul", 10, 1, circle, bottomtype, xdim, ydim, zdim, bottomh, "mm")

	plate = wtype.NewLHPlate("Agarplateforpicking48", "Unknown", 6, 8, 14, "mm", welltypecostar48forpicking, wellxoffset, wellyoffset, xstart, ystart, zstart)
	plates = append(plates, plate)

	/// placeholder for non plate container for testing

	plate = wtype.NewLHPlate("1L_DuranBottle", "Unknown", 8, 12, 25.7, "mm", singlewelltrough, 9, 9, 0.0, 0.0, 15.5)
	plates = append(plates, plate)

	//forward position

	//	ep48g := wtype.NewShape("trap", "mm", 2, 4, 2)
	//	welltype = wtype.NewLHWell("EPAGE48", "", "", "ul", 15, 0, ep48g, 0, 2, 4, 2, 48, "mm")
	//	plate = wtype.NewLHPlate("EPAGE48", "Invitrogen", 2, 26, 50, "mm", welltype, 4.5, 34, 0.0, 0.0, 2.0)
	//	plates[plate.Type] = plate

	////// E gel dimensions
	xdim = 2
	ydim = 4
	zdim = 2
	bottomh = 2

	wellxoffset = 4.5              // centre of well to centre of neighbouring well in x direction
	wellyoffset = 33.75            //centre of well to centre of neighbouring well in y direction
	xstart = -1.0                  // distance from top left side of plate to first well
	ystart = 18.25                 // distance from top left side of plate to first well
	zstart = riserheightinmm + 4.5 // offset of bottom of deck to bottom of well

	//E-PAGE 48 (reverse) position
	ep48g := wtype.NewShape("trap", "mm", xdim, ydim, zdim)
	//can't reach all wells; change to 24 wells per row? yes!
	egelwell := wtype.NewLHWell("EPAGE48", "", "", "ul", 20, 0, ep48g, wtype.LHWBFLAT, xdim, ydim, zdim, bottomh, "mm")
	gelplate := wtype.NewLHPlate("EPAGE48", "Invitrogen", 2, 24, 48.5, "mm", egelwell, wellxoffset, wellyoffset, xstart, ystart, zstart)

	gelconsar := []string{"position_9"}
	gelplate.SetConstrained("Pipetmax", gelconsar)
	gelplate.DeclareSpecial() // Do this for racks, other very unusual plate types

	plates = append(plates, gelplate)

	//E-GEL 48 (reverse) position
	gelplate = wtype.NewLHPlate("EGEL48", "Invitrogen", 2, 24, 48.5, "mm", egelwell, wellxoffset, wellyoffset, xstart, ystart, zstart)
	gelplate.SetConstrained("Pipetmax", gelconsar)
	gelplate.DeclareSpecial() // Do this for racks, other very unusual plate types
	plates = append(plates, gelplate)

	//E-GEL 96 definition
	//same welltype as EPAGE
	//due to staggering of wells: 1 96well gel is set up as two well types

	// 1st type
	//can't reach all wells; change to 12 wells per row?
	gelplate = wtype.NewLHPlate("EGEL96_1", "Invitrogen", 4, 13, 48.5, "mm", egelwell, 9, 18.0, -9.0, -0.5, riserheightinmm+5.5)
	gelplate.SetConstrained("Pipetmax", gelconsar)
	gelplate.DeclareSpecial() // Do this for racks, other very unusual plate types
	plates = append(plates, gelplate)

	// 2nd type
	gelplate = wtype.NewLHPlate("EGEL96_2", "Invitrogen", 4, 13, 48.5, "mm", egelwell, 9, 18.0, -5.0, 9, riserheightinmm+5.5)
	gelplate.SetConstrained("Pipetmax", gelconsar)
	gelplate.DeclareSpecial() // Do this for racks, other very unusual plate types
	plates = append(plates, gelplate)

	// Nunclon 12 well plate with Agar flat bottom 2ml per well

	bottomtype = wtype.LHWBFLAT
	xdim = 22.5 // diameter
	ydim = 22.5 // diameter
	zdim = 17.0
	bottomh = 9.0 //(this includes agar estimate)

	wellxoffset = 27.0 // centre of well to centre of neighbouring well in x direction
	wellyoffset = 27.0 //centre of well to centre of neighbouring well in y direction
	xstart = 11.0      // distance from top left side of plate to first well
	ystart = 4.0       // distance from top left side of plate to first well
	zstart = 9.0       // offset of bottom of deck to bottom of well (this includes Agar height estimate)

	wellsperrow = 4
	wellspercolumn = 3
	heightinmm = 19.0

	circle = wtype.NewShape("cylinder", "mm", xdim, ydim, zdim)
	welltype12well := wtype.NewLHWell("Nuncon12well", "", "", "ul", 1000, 10, circle, bottomtype, xdim, ydim, zdim, bottomh, "mm")

	plate = wtype.NewLHPlate("Nuncon12wellAgar", "Unknown", wellspercolumn, wellsperrow, heightinmm, "mm", welltype12well, wellxoffset, wellyoffset, xstart, ystart, zstart)
	plate.DeclareSpecial() // Do this for racks, other very unusual plate types
	plates = append(plates, plate)

	// update z start to remove agar estimate and make new plate type
	zstart = 4.0 // offset of bottom of deck to bottom of well
	plate = wtype.NewLHPlate("Nuncon12well", "Unknown", wellspercolumn, wellsperrow, heightinmm, "mm", welltype12well, wellxoffset, wellyoffset, xstart, ystart, zstart)
	plate.DeclareSpecial() // Do this for racks, other very unusual plate types
	plates = append(plates, plate)

	//VWR 12 Well Plate 734-2324 NO AGAR

	bottomtype = wtype.LHWBFLAT
	xdim = 24.0 // diameter
	ydim = 24.0 // diameter
	zdim = 19.0
	bottomh = 5.0 //(this includes agar estimate)

	wellxoffset = 27.0 // centre of well to centre of neighbouring well in x direction
	wellyoffset = 27.0 //centre of well to centre of neighbouring well in y direction
	xstart = 11.0      // distance from top left side of plate to first well
	ystart = 4.0       // distance from top left side of plate to first well
	zstart = 1.0       // offset of bottom of deck to bottom of well (this includes agar estimate)

	wellsperrow = 4
	wellspercolumn = 3
	heightinmm = 20.0

	circle = wtype.NewShape("cylinder", "mm", xdim, ydim, zdim)
	welltypevwr12 := wtype.NewLHWell("VWR12", "", "", "ul", 100, 10, circle, bottomtype, xdim, ydim, zdim, bottomh, "mm")

	plate = wtype.NewLHPlate("VWR12well", "Unknown", wellspercolumn, wellsperrow, heightinmm, "mm", welltypevwr12, wellxoffset, wellyoffset, xstart, ystart, zstart)
	plate.DeclareSpecial() // Do this for racks, other very unusual plate types
	plates = append(plates, plate)

	//Nunclon 8 well Plate 167064 DOW
	bottomtype = wtype.LHWBFLAT
	xdim = 30.0
	ydim = 39.0
	zdim = 11.0
	bottomh = 11.0 //accounts for agar estimate

	wellxoffset = 30.0 // centre of well to centre of neighbouring well in x direction
	wellyoffset = 40.0 // centre of well to centre of neighbouring well in y direction
	xstart = 5.0       // distance from top left side of plate to first well
	ystart = 10.5      // distance from top left side of plate to first well
	zstart = 3.0       // offset of bottom of deck to bottom of well

	wellsperrow = 4.0
	wellspercolumn = 2.0
	heightinmm = 11.0

	nuncsquare := wtype.NewShape("box", "mm", 30, 39, 11)
	welltypenunc8 := wtype.NewLHWell("nuncsquare", "", "", "ul", 3000, 10, nuncsquare, bottomtype, xdim, ydim, zdim, bottomh, "mm")

	plate = wtype.NewLHPlate("nunc8well", "Unknown", wellspercolumn, wellsperrow, heightinmm, "mm", welltypenunc8, wellxoffset, wellyoffset, xstart, ystart, zstart)
	plate.DeclareSpecial() // Do this for racks, other very unusual plate types

	plates = append(plates, plate)

	plate = makeGreinerVBottomPlate()
	plates = append(plates, plate)

	//plate = makeGreinerVBottomPlateWithRiser()
	//plates = append(plates, plate)

	//plate = plate.Dup()
	//plate.Type += "40"
	//plates = append(plates, plate)

	plate = makeHighResplateforPicking()
	plates = append(plates, plate)

	plate = makeGreinerFlatBottomBlackPlate()
	plates = append(plates, plate)

	return
}

func makeGreinerVBottomPlate() *wtype.LHPlate {
	// greiner V96

	bottomtype := wtype.LHWBV
	xdim := 6.2
	ydim := 6.2
	zdim := 11.0
	bottomh := 1.0

	wellxoffset := 9.0 // centre of well to centre of neighbouring well in x direction
	wellyoffset := 9.0 //centre of well to centre of neighbouring well in y direction
	xstart := 0.0      // distance from top left side of plate to first well
	ystart := 0.0      // distance from top left side of plate to first well
	zstart := 2.0      // offset of bottom of deck to bottom of well

	rwshp := wtype.NewShape("cylinder", "mm", 6.2, 6.2, 10.0)
	welltype := wtype.NewLHWell("GreinerSWVBottom", "", "", "ul", 500, 1, rwshp, bottomtype, xdim, ydim, zdim, bottomh, "mm")

	plate := wtype.NewLHPlate("GreinerSWVBottom", "Greiner", 8, 12, 15, "mm", welltype, wellxoffset, wellyoffset, xstart, ystart, zstart)

	return plate
}

func makeGreinerFlatBottomBlackPlate() *wtype.LHPlate {
	// shallow round well flat bottom 96
	rwshp := wtype.NewShape("cylinder", "mm", 8.2, 8.2, 11)
	roundwell96 := wtype.NewLHWell("SRWFB96", "", "", "ul", 500, 10, rwshp, 0, 8.2, 8.2, 11, 1.0, "mm")
	plate := wtype.NewLHPlate("greiner96Black", "greiner", 8, 12, 15, "mm", roundwell96, 9, 9, 0.0, 0.0, 1.0)
	return plate
}

// Onewell SBS format Agarplate with colonies on shallowriser (50ml agar) very high res
func makeHighResplateforPicking() *wtype.LHPlate {

	bottomtype := wtype.LHWBFLAT
	xdim := 1.4 // of well
	ydim := 1.4
	zdim := 7.0
	bottomh := 0.5

	wellxoffset := 1.55 // centre of well to centre of neighbouring well in x direction
	wellyoffset := 1.55 // centre of well to centre of neighbouring well in y direction
	xstart := -3.885    // distance from top left side of plate to first well
	ystart := -3.0      // distance from top left side of plate to first well
	zstart := 3.25      // offset of bottom of deck to bottom of well

	square3150 := wtype.NewShape("box", "mm", xdim, ydim, zdim)
	welltype3150 := wtype.NewLHWell("3150flat", "", "", "ul", 5, 0.5, square3150, bottomtype, xdim, ydim, zdim, bottomh, "mm")

	// greiner one well with 50ml of agar in
	plate := wtype.NewLHPlate("Agarplateforpicking3150", "Unknown", 45, 70, 7, "mm", welltype3150, wellxoffset, wellyoffset, xstart, ystart, zstart)

	return plate
}

func makeGreinerVBottomPlateWithRiser() *wtype.LHPlate {
	plate := makeGreinerVBottomPlate()
	plate.Type = "GreinerSWVBottom_riser"
	plate.WellZStart = 42.0
	return plate
}
