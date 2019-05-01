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

package main

import (
	"encoding/json"
	"strings"

	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wutil"
)

// The height below which an error will be generated
// when attempting to perform transfers with low volume head and tips (0.5 - 20ul) on the Gilson PipetMax.
/// uuuuurgh
const MinimumZHeightPermissableForLVPipetMax = 0.636

// deprecated
/*
var platespecificoffset = map[string]float64{
	"pcrplate_skirted": gilsonoffsetpcrplate,
	"greiner384":       gilsonoffsetgreiner,
	"costar48well":     3.0,
	"Nuncon12well":     11.0, // this must be wrong!! check z start without riser properly
	"Nuncon12wellAgar": 11.0, // this must be wrong!! check z start without riser properly
	"VWR12well":        3.0,
}
*/

var platespecificoffset = map[string]float64{}

// function to check if a platename already contains a riser
func containsRiser(plate *wtype.Plate) bool {
	for _, dev := range defaultDevices {
		for _, synonym := range dev.GetSynonyms() {
			if strings.Contains(plate.Type, "_"+synonym) {
				return true
			}
		}
	}

	return false
}

func addRiser(plate *wtype.Plate, riser device) (plates []*wtype.Plate) {
	if containsRiser(plate) || doNotAddThisRiserToThisPlate(plate, riser) {
		return
	}

	for _, risername := range riser.GetSynonyms() {
		var dontaddrisertothisplate bool

		newplate := plate.Dup()
		riserheight := riser.GetHeightInmm()
		if offset, found := platespecificoffset[plate.Type]; found {
			riserheight = riserheight - offset
		}

		riserheight = riserheight + plateRiserSpecificOffset(plate, riser)

		newplate.WellZStart = plate.WellZStart + riserheight
		newplate.Bounds.SetSize(plate.GetSize().Add(wtype.Coordinates3D{X: 0.0, Y: 0.0, Z: riserheight}))
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

func addAllDevices(plates []*wtype.Plate) (ret []*wtype.Plate) {
	for _, plate := range plates {
		for _, dev := range defaultDevices {
			ret = append(ret, addRiser(plate, dev)...)
		}
	}
	return
}

func makePlates() (plates []*wtype.Plate) {
	plates = makeBasicPlates()
	additional := addAllDevices(plates)
	return append(plates, additional...)
}

func makeBasicPlates() (plates []*wtype.Plate) {
	// deep square well 96
	swshp := wtype.NewShape(wtype.BoxShape, "mm", 8.2, 8.2, 41.3)
	deepsquarewell := wtype.NewLHWell("ul", 2000, 420, swshp, wtype.VWellBottom, 8.2, 8.2, 41.3, 4.7, "mm")
	plate := wtype.NewLHPlate("DSW96", "Unknown", 8, 12, makePlateCoords(44.1), deepsquarewell, 9, 9, 0.0, 0.0, valueformaxheadtonotintoDSWplatewithp20tips)
	plates = append(plates, plate)

	// Nunc™2.0mL DeepWell™ Plates 95040452
	nunc96deepwellshp := wtype.NewShape(wtype.BoxShape, "mm", 8.5, 8.5, 41.5)
	nunc96deepwell := wtype.NewLHWell("ul", 2000, 420, nunc96deepwellshp, wtype.UWellBottom, 8.2, 8.2, 41.3, 2.5, "mm")
	plate = wtype.NewLHPlate("Nunc96DeepWell", "Unknown", 8, 12, makePlateCoords(43.6), nunc96deepwell, 9, 9, -1.0, 0.0, 6.5)
	plates = append(plates, plate)

	// Thermo 96 well conical btm pp pit natural 0.45 ml well Cat Num: 249946. (TWIST DNA Plate)
	twist96wellshp := wtype.NewShape(wtype.CylinderShape, "mm", 6.7, 6.7, 9.8)
	twist96well := wtype.NewLHWell("ul", 450, 10, twist96wellshp, wtype.VWellBottom, 6.7, 6.7, 9.8, 4.6, "mm")
	plate = wtype.NewLHPlate("TwistDNAPlate", "Unknown", 8, 12, makePlateCoords(14.4), twist96well, 9.0, 9.0, 0.0, 0.0, -1.9)
	plates = append(plates, plate)

	// IDT/ABgene 1.2 ml storage plate AB0564
	idtshp := wtype.NewShape(wtype.CylinderShape, "mm", 7, 7, 39.35)
	idtroundwell96 := wtype.NewLHWell("ul", 1200, 100, idtshp, wtype.UWellBottom, 7, 7, 39.35, 3, "mm")
	plate = wtype.NewLHPlate("IDT96", "Unknown", 8, 12, makePlateCoords(42.5), idtroundwell96, 9, 9, 0, 0, 3)
	plates = append(plates, plate)

	//4 column reservoir plate Phenix Research Products RRI3051; Fisher cat# NC0336913
	fourcolumnshp := wtype.NewShape(wtype.BoxShape, "mm", 26, 71, 42)
	fourcolumnwell := wtype.NewLHWell("ul", 73000, 3000, fourcolumnshp, wtype.VWellBottom, 26, 71, 42, 2, "mm")
	plate = wtype.NewLHPlate("FourColumnReservoir", "Unknown", 1, 4, makePlateCoords(44), fourcolumnwell, 26, 1, 9.5, 31, 1) //WellYStart is not accurate, but would not visualise correctly unless set to this value, cant diagnose
	plates = append(plates, plate)

	// 24 well deep square well plate on riser

	bottomtype := wtype.VWellBottom // 0 = flat, 2 = v shaped
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
	//zstart := -1.0      // offset of bottom of deck to bottom of well (this includes agar estimate)
	zstart := 0.0 // offset of bottom of deck to bottom of well (this includes agar estimate)

	heightinmm := 44.1

	squarewell := wtype.NewShape(wtype.BoxShape, "mm", xdim, ydim, zdim)
	squarewell24 := wtype.NewLHWell(welltypeunit, wellcapacityinwelltypeunit, residualvol, squarewell, bottomtype, xdim, ydim, zdim, bottomh, "mm")
	plate = wtype.NewLHPlate("DSW24", "Unknown", wellspercolumn, wellsperrow, makePlateCoords(heightinmm), squarewell24, wellxoffset, wellyoffset, xstart, ystart, zstart)
	plate.DeclareSpecial() // Do this for racks, other very unusual plate types

	plates = append(plates, plate)

	// shallow round well flat bottom 96
	rwshp := wtype.NewShape(wtype.CylinderShape, "mm", 8.2, 8.2, 11)
	roundwell96 := wtype.NewLHWell("ul", 340, 25, rwshp, 0, 8.2, 8.2, 11, 1.0, "mm")
	plate = wtype.NewLHPlate("SRWFB96", "Unknown", 8, 12, makePlateCoords(15), roundwell96, 9, 9, 0.0, 0.0, 2.2)
	plates = append(plates, plate)

	// deep well strip trough 12
	stshp := wtype.NewShape(wtype.BoxShape, "mm", 8.2, 72, 41.3)
	trough12 := wtype.NewLHWell("ul", 15000, 5000, stshp, wtype.VWellBottom, 8.2, 72, 41.3, 4.7, "mm")
	plate = wtype.NewLHPlate("DWST12", "Unknown", 1, 12, makePlateCoords(44.1), trough12, 9, 9, 0, 30.0, valueformaxheadtonotintoDSWplatewithp20tips)
	//	plate.DeclareSpecial() // Do this for racks, other very unusual plate types
	plates = append(plates, plate)

	// shallow well strip trough 12
	stshps := wtype.NewShape(wtype.BoxShape, "mm", 8.2, 72, 15)
	trough12s := wtype.NewLHWell("ul", 4000, 1500, stshps, wtype.VWellBottom, 8.2, 72, 15, 4.7, "mm")
	plate = wtype.NewLHPlate("SWST12", "Unknown", 1, 12, makePlateCoords(20), trough12s, 9, 9, 0, 30.0, 1)
	//	plate.DeclareSpecial() // Do this for racks, other very unusual plate types
	plates = append(plates, plate)

	// deep well strip trough 8
	stshp8 := wtype.NewShape(wtype.BoxShape, "mm", 115.0, 8.2, 41.3)
	trough8 := wtype.NewLHWell("ul", 24000, 1000, stshp8, wtype.VWellBottom, 115, 8.2, 41.3, 4.7, "mm")
	plate = wtype.NewLHPlate("DWST8", "Unknown", 8, 1, makePlateCoords(44.1), trough8, 9, 9, 49.5, 0.0, 0.0)
	plate.DeclareSpecial() // Do this for racks, other very unusual plate types... except troughs?!
	plates = append(plates, plate)

	// 250ml box reservoir
	reservoirbox := wtype.NewShape(wtype.BoxShape, "mm", 121, 80, 40) // 39?
	welltypereservoir := wtype.NewLHWell("ul", 200000, 40000, reservoirbox, wtype.FlatWellBottom, 121, 80, 40, 3, "mm")
	plate = wtype.NewLHPlate("reservoir", "unknown", 1, 1, makePlateCoords(40), welltypereservoir, 1, 1, 49.5, 31.0, 0.0)
	plates = append(plates, plate)

	// falcon 6 well plate with Agar flat bottom with 4ml per well

	bottomtype = wtype.FlatWellBottom
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

	circle := wtype.NewShape(wtype.CylinderShape, "mm", 37, 37, 20)
	welltype6well := wtype.NewLHWell("ul", 4000, 1, circle, bottomtype, xdim, ydim, zdim, bottomh, "mm")

	plate = wtype.NewLHPlate("falcon6wellAgar", "Unknown", wellspercolumn, wellsperrow, makePlateCoords(heightinmm), welltype6well, wellxoffset, wellyoffset, xstart, ystart, zstart)
	plate.DeclareSpecial() // Do this for racks, other very unusual plate types
	plates = append(plates, plate)

	// Costar 48 well plate flat bottom

	bottomtype = wtype.FlatWellBottom
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

	circle = wtype.NewShape(wtype.CylinderShape, "mm", xdim, ydim, zdim)
	welltypecostar48 := wtype.NewLHWell("ul", 1000, 100, circle, bottomtype, xdim, ydim, zdim, bottomh, "mm")

	plate = wtype.NewLHPlate("costar48well", "Unknown", wellspercolumn, wellsperrow, makePlateCoords(heightinmm), welltypecostar48, wellxoffset, wellyoffset, xstart, ystart, zstart)
	plate.DeclareSpecial() // Do this for racks, other very unusual plate types
	plates = append(plates, plate)

	// Block Kombi 2ml
	eppy := wtype.NewShape(wtype.CylinderShape, "mm", 8.2, 8.2, 45)

	wellxoffset = 18.0 // centre of well to centre of neighbouring well in x direction
	wellyoffset = 18.0 //centre of well to centre of neighbouring well in y direction
	xstart = 5.0       // distance from top left side of plate to first well
	ystart = 10.0      // distance from top left side of plate to first well
	zstart = 6.0       // offset of bottom of deck to bottom of well

	welltype2mleppy := wtype.NewLHWell("ul", 2000, 50, eppy, wtype.VWellBottom, 8.2, 8.2, 45, 4.7, "mm")

	plate = wtype.NewLHPlate("Kombi2mlEpp", "Unknown", 4, 2, makePlateCoords(45), welltype2mleppy, wellxoffset, wellyoffset, xstart, ystart, zstart)
	plates = append(plates, plate)

	plate.DeclareSpecial() // Do this for racks, other very unusual plate types

	// Eppendorfrack 425 for 2ml tubes

	wellxoffset = 18.0 // centre of well to centre of neighbouring well in x direction
	wellyoffset = 18.0 //centre of well to centre of neighbouring well in y direction
	xstart = 4.5       // distance from top left side of plate to first well
	ystart = 5.0       // distance from top left side of plate to first well
	zstart = 5.0       // offset of bottom of deck to bottom of well

	plate = wtype.NewLHPlate("eppendorfrack425_2ml", "Unknown", 4, 6, makePlateCoords(45), welltype2mleppy, wellxoffset, wellyoffset, xstart, ystart, zstart)
	plate.DeclareSpecial() // Do this for racks, other very unusual plate types
	plates = append(plates, plate)

	// Eppendorfrack 425 for 1.5ml tubes

	//values from physical measurements, HJK 2/7/18
	wellRadius := 9.2
	wellHeight := 39.5
	outerRadius := 10.55
	plateHeight := 45.75 //up to lip of eppendorf

	bottomH := (outerRadius - wellRadius) * 0.5

	wellxoffset = 18.0                               // centre of well to centre of neighbouring well in x direction
	wellyoffset = 18.0                               // centre of well to centre of neighbouring well in y direction
	xstart = 4.5                                     // distance from top left side of plate to first well
	ystart = 5.0                                     // distance from top left side of plate to first well
	zstart = plateHeight - wellHeight - zStartOffset // offset of bottom of deck to bottom of well
	//zStartOffset gets added later

	welltypesmallereppy := wtype.NewLHWell("ul", 1500, 50, eppy, wtype.VWellBottom, wellRadius, wellRadius, wellHeight, bottomH, "mm")

	plate = wtype.NewLHPlate("eppendorfrack425_1.5ml", "Unknown", 4, 6, makePlateCoords(plateHeight), welltypesmallereppy, wellxoffset, wellyoffset, xstart, ystart, zstart)
	plate.DeclareSpecial() // Do this for racks, other very unusual plate types
	plates = append(plates, plate)

	// Eppendorfrack 424 with lid holders and using 2ml tubes

	wellxoffset = 36.0 // centre of well to centre of neighbouring well in x direction
	wellyoffset = 18.0 //centre of well to centre of neighbouring well in y direction
	xstart = 14.0      // distance from top left side of plate to first well
	ystart = 5.0       // distance from top left side of plate to first well
	zstart = 5.0       // offset of bottom of deck to bottom of well

	plate = wtype.NewLHPlate("eppendorfrack424_2ml_lidholder", "Unknown", 4, 3, makePlateCoords(45), welltype2mleppy, wellxoffset, wellyoffset, xstart, ystart, zstart)
	plate.DeclareSpecial() // Do this for racks, other very unusual plate types
	plates = append(plates, plate)

	// Eppendorfrack 424 with lid holders and using 1.5ml tubes

	wellxoffset = 36.0 // centre of well to centre of neighbouring well in x direction
	wellyoffset = 18.0 //centre of well to centre of neighbouring well in y direction
	xstart = 14.0      // distance from top left side of plate to first well
	ystart = 5.0       // distance from top left side of plate to first well
	zstart = 4.5       // offset of bottom of deck to bottom of well

	plate = wtype.NewLHPlate("eppendorfrack424_1.5ml_lidholder", "Unknown", 4, 3, makePlateCoords(45), welltypesmallereppy, wellxoffset, wellyoffset, xstart, ystart, zstart)
	plate.DeclareSpecial() // Do this for racks, other very unusual plate types
	plates = append(plates, plate)

	// greiner 384 well plate flat bottom

	bottomtype = wtype.FlatWellBottom
	xdim = 4.0
	ydim = 4.0
	zdim = 12.0
	bottomh = 1.0

	wellxoffset = 4.5 // centre of well to centre of neighbouring well in x direction
	wellyoffset = 4.5 //centre of well to centre of neighbouring well in y direction
	xstart = -2.5     // distance from top left side of plate to first well
	ystart = -2.5     // distance from top left side of plate to first well
	zstart = 2.5      // offset of bottom of deck to bottom of well

	square := wtype.NewShape(wtype.BoxShape, "mm", xdim, ydim, zdim)
	welltype384 := wtype.NewLHWell("ul", 125, 20, square, bottomtype, xdim, ydim, zdim, bottomh, "mm")

	//func NewLHPlate(platetype, mfr string, nrows, ncols int, height float64, hunit string, welltype *LHWell, wellXOffset, wellYOffset, wellXStart, wellYStart, wellZStart float64) *LHPlate {
	plate = wtype.NewLHPlate("greiner384", "Unknown", 16, 24, makePlateCoords(14), welltype384, wellxoffset, wellyoffset, xstart, ystart, zstart)
	plates = append(plates, plate)

	// NUNC 1536 well plate flat bottom on riser

	bottomtype = wtype.FlatWellBottom
	xdim = 2.0 // of well
	ydim = 2.0
	zdim = 7.0
	bottomh = 0.5

	wellxoffset = 2.25 // centre of well to centre of neighbouring well in x direction
	wellyoffset = 2.25 //centre of well to centre of neighbouring well in y direction
	xstart = -2.5      // distance from top left side of plate to first well
	ystart = -2.5      // distance from top left side of plate to first well
	zstart = 2         // offset of bottom of deck to bottom of well

	square1536 := wtype.NewShape(wtype.BoxShape, "mm", xdim, ydim, zdim)
	welltype1536 := wtype.NewLHWell("ul", 13, 2, square1536, bottomtype, xdim, ydim, zdim, bottomh, "mm")

	plate = wtype.NewLHPlate("nunc1536", "Unknown", 32, 48, makePlateCoords(7), welltype1536, wellxoffset, wellyoffset, xstart, ystart, zstart)
	plates = append(plates, plate)

	// Onewell SBS format Agarplate with colonies on riser (50ml agar) high res

	wellxoffset = 2.25  // centre of well to centre of neighbouring well in x direction
	wellyoffset = 2.250 //centre of well to centre of neighbouring well in y direction
	xstart = -2.5       // distance from top left side of plate to first well
	ystart = -2.5       // distance from top left side of plate to first well
	zstart = 3          // offset of bottom of deck to bottom of well

	// greiner one well with 50ml of agar in
	plate = wtype.NewLHPlate("Agarplateforpicking1536", "Unknown", 32, 48, makePlateCoords(7), welltype1536, wellxoffset, wellyoffset, xstart, ystart, zstart)

	plates = append(plates, plate)

	// Onewell SBS format Agarplate with colonies on riser (50ml agar) low res

	wellxoffset = 4.5 // centre of well to centre of neighbouring well in x direction
	wellyoffset = 4.5 //centre of well to centre of neighbouring well in y direction
	xstart = -2.5     // distance from top left side of plate to first well
	ystart = -2.5     // distance from top left side of plate to first well
	zstart = 3.5      //5.5 // offset of bottom of deck to bottom of well

	// greiner one well with 50ml of agar in
	plate = wtype.NewLHPlate("Agarplateforpicking384", "Unknown", 16, 24, makePlateCoords(14), welltype384, wellxoffset, wellyoffset, xstart, ystart, zstart)
	plates = append(plates, plate)

	// Onewell SBS format Agarplate with colonies on shallowriser (30ml agar) low res

	wellxoffset = 4.5 // centre of well to centre of neighbouring well in x direction
	wellyoffset = 4.5 //centre of well to centre of neighbouring well in y direction
	xstart = -2.5     // distance from top left side of plate to first well
	ystart = -2.5     // distance from top left side of plate to first well
	zstart = 1        // offset of bottom of deck to bottom of well

	plate = wtype.NewLHPlate("30mlAgarplateforpicking384", "Unknown", 16, 24, makePlateCoords(14), welltype384, wellxoffset, wellyoffset, xstart, ystart, zstart)
	plates = append(plates, plate)

	// Onewell SBS format Agarplate with colonies on riser (50ml agar) med res

	bottomtype = wtype.FlatWellBottom
	xdim = 3.0
	ydim = 3.0
	zdim = 14.0
	bottomh = 1.0

	wellxoffset = 3.1 // centre of well to centre of neighbouring well in x direction
	wellyoffset = 3.1 //centre of well to centre of neighbouring well in y direction
	xstart = -2.5     // distance from top left side of plate to first well
	ystart = -2.5     // distance from top left side of plate to first well
	zstart = 3.5      //5.5 // offset of bottom of deck to bottom of well

	square768 := wtype.NewShape(wtype.BoxShape, "mm", xdim, ydim, zdim)
	welltype768 := wtype.NewLHWell("ul", 31.25, 5, square768, bottomtype, xdim, ydim, zdim, bottomh, "mm")

	// greiner one well with 50ml of agar in
	plate = wtype.NewLHPlate("Agarplateforpicking768", "Unknown", 24, 32, makePlateCoords(14), welltype768, wellxoffset, wellyoffset, xstart, ystart, zstart)
	plates = append(plates, plate)

	// Onewell SBS format Agarplate with colonies on riser (50ml agar) low res with 96 well map

	bottomtype = wtype.UWellBottom
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

	pcrplatewellforpicking := wtype.NewLHWell("ul", 5, 1, wtype.NewShape(wtype.CylinderShape, "mm", 5.5, 5.5, 15), bottomtype, xdim, ydim, zdim, bottomh, "mm")

	plate = wtype.NewLHPlate("Agarplateforpicking96", "Unknown", 8, 12, makePlateCoords(14), pcrplatewellforpicking, wellxoffset, wellyoffset, xstart, ystart, zstart)
	plates = append(plates, plate)

	// Onewell SBS format Agarplate with colonies on riser (50ml agar) low res with 48 well map

	bottomtype = wtype.FlatWellBottom
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

	welltypecostar48forpicking := wtype.NewLHWell("ul", 10, 1, circle, bottomtype, xdim, ydim, zdim, bottomh, "mm")

	plate = wtype.NewLHPlate("Agarplateforpicking48", "Unknown", 6, 8, makePlateCoords(14), welltypecostar48forpicking, wellxoffset, wellyoffset, xstart, ystart, zstart)
	plate.DeclareSpecial()
	plates = append(plates, plate)

	/// placeholder for non plate container for testing

	plate = wtype.NewLHPlate("1L_DuranBottle", "Unknown", 8, 12, makePlateCoords(25.7), welltypereservoir, 9, 9, 0.0, 0.0, 15.5)
	plates = append(plates, plate)

	////// E gel dimensions

	wellxoffset = 4.5   // centre of well to centre of neighbouring well in x direction
	wellyoffset = 33.75 //centre of well to centre of neighbouring well in y direction
	xstart = -1.0       // distance from top left side of plate to first well
	//ystart = 18.25                 // distance from top left side of plate to first well
	ystart = 17.75                 // distance from top left side of plate to first well MIS -- revised from above after physical test
	zstart = riserheightinmm + 4.5 // offset of bottom of deck to bottom of well
	eplateheight := 48.5

	xdim = 2
	ydim = 4
	//wells extend up to the top of the plate
	zdim = eplateheight - (zstart + zStartOffset) //zStartOffset will get added to the zstart later
	bottomh = 2

	//E-PAGE 48 (reverse) position
	ep48g := wtype.NewShape(wtype.TrapezoidShape, "mm", xdim, ydim, zdim)
	//can't reach all wells; change to 24 wells per row? yes!
	egelwell := wtype.NewLHWell("ul", 20, 0, ep48g, wtype.FlatWellBottom, xdim, ydim, zdim, bottomh, "mm")
	gelplate := wtype.NewLHPlate("EPAGE48", "Invitrogen", 2, 24, makePlateCoords(eplateheight), egelwell, wellxoffset, wellyoffset, xstart, ystart, zstart)

	gelconsar := []string{"position_9"}
	gelplate.SetConstrained("Pipetmax", gelconsar)
	gelplate.DeclareSpecial() // Do this for racks, other very unusual plate types

	plates = append(plates, gelplate)

	//E-GEL 48 (reverse) position
	gelplate = wtype.NewLHPlate("EGEL48", "Invitrogen", 2, 24, makePlateCoords(eplateheight), egelwell, wellxoffset, wellyoffset, xstart, ystart, zstart)
	gelplate.SetConstrained("Pipetmax", gelconsar)
	gelplate.DeclareSpecial() // Do this for racks, other very unusual plate types
	plates = append(plates, gelplate)

	//E-GEL 96 definition
	//same welltype as EPAGE
	//due to staggering of wells: 1 96well gel is set up as two well types

	// 1st type
	//can't reach all wells; change to 12 wells per row?
	gelplate = wtype.NewLHPlate("EGEL96_1", "Invitrogen", 4, 13, makePlateCoords(48.5), egelwell, 9, 18.0, -9.0, -0.5, riserheightinmm+5.5)
	gelplate.SetConstrained("Pipetmax", gelconsar)
	gelplate.DeclareSpecial() // Do this for racks, other very unusual plate types
	plates = append(plates, gelplate)

	// 2nd type
	gelplate = wtype.NewLHPlate("EGEL96_2", "Invitrogen", 4, 13, makePlateCoords(48.5), egelwell, 9, 18.0, -5.0, 9, riserheightinmm+5.5)
	gelplate.SetConstrained("Pipetmax", gelconsar)
	gelplate.DeclareSpecial() // Do this for racks, other very unusual plate types
	plates = append(plates, gelplate)

	// Nunclon 12 well plate with Agar flat bottom 2ml per well

	bottomtype = wtype.FlatWellBottom
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

	circle = wtype.NewShape(wtype.CylinderShape, "mm", xdim, ydim, zdim)
	welltype12well := wtype.NewLHWell("ul", 1000, 10, circle, bottomtype, xdim, ydim, zdim, bottomh, "mm")

	plate = wtype.NewLHPlate("Nuncon12wellAgar", "Unknown", wellspercolumn, wellsperrow, makePlateCoords(heightinmm), welltype12well, wellxoffset, wellyoffset, xstart, ystart, zstart)
	plate.DeclareSpecial() // Do this for racks, other very unusual plate types
	plates = append(plates, plate)

	// update z start to remove agar estimate and make new plate type
	zstart = 4.0 // offset of bottom of deck to bottom of well
	plate = wtype.NewLHPlate("Nuncon12well", "Unknown", wellspercolumn, wellsperrow, makePlateCoords(heightinmm), welltype12well, wellxoffset, wellyoffset, xstart, ystart, zstart)
	plate.DeclareSpecial() // Do this for racks, other very unusual plate types
	plates = append(plates, plate)

	//VWR 12 Well Plate 734-2324 NO AGAR

	bottomtype = wtype.FlatWellBottom
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

	circle = wtype.NewShape(wtype.CylinderShape, "mm", xdim, ydim, zdim)
	welltypevwr12 := wtype.NewLHWell("ul", 100, 10, circle, bottomtype, xdim, ydim, zdim, bottomh, "mm")

	plate = wtype.NewLHPlate("VWR12well", "Unknown", wellspercolumn, wellsperrow, makePlateCoords(heightinmm), welltypevwr12, wellxoffset, wellyoffset, xstart, ystart, zstart)
	plate.DeclareSpecial() // Do this for racks, other very unusual plate types
	plates = append(plates, plate)

	//Nunclon 8 well Plate 167064 DOW
	bottomtype = wtype.FlatWellBottom
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

	nuncsquare := wtype.NewShape(wtype.BoxShape, "mm", 30, 39, 11)
	welltypenunc8 := wtype.NewLHWell("ul", 3000, 10, nuncsquare, bottomtype, xdim, ydim, zdim, bottomh, "mm")

	plate = wtype.NewLHPlate("nunc8well", "Unknown", wellspercolumn, wellsperrow, makePlateCoords(heightinmm), welltypenunc8, wellxoffset, wellyoffset, xstart, ystart, zstart)
	plate.DeclareSpecial() // Do this for racks, other very unusual plate types

	plates = append(plates, plate)

	plate = makeGreinerVBottomPlate()
	plates = append(plates, plate)

	plate = makeHighResplateforPicking()
	plates = append(plates, plate)

	plate = makeGreinerFlatBottomBlackPlate()
	plates = append(plates, plate)

	plates = append(plates, makeNunc96UPlate())

	plates = append(plates,
		makeSkirtedPCRPlate(),
		makeStripTube(),
		makeSemiSkirtedPCRPlate(),
		makePCRPlate(),
		makeFluidX700ulPlate(),
	)

	plates = append(plates, make96DeepWellLowVolumePlate())
	plates = append(plates, makeLabcyte384PPStdV())
	plates = append(plates, make384wellplateAppliedBiosystems())
	plates = append(plates, makeAcroPrep384NoFilter())
	plates = append(plates, makeAcroPrep384WithFilter())
	plates = append(plates, makeGreinerVFromSpec())
	plates = append(plates, make4TitudePcrPlateFromSpec())
	return
}

func makePCRPlateWell() *wtype.LHWell {
	// well area function
	// -- determined empirically since inverse cubic was giving us some numerical issues
	areaf := wutil.Quartic{A: -3.3317851312e-09, B: 0.00000225834467, C: -0.0006305492472, D: 0.1328156706978, E: 0}
	afb, _ := json.Marshal(areaf)
	afs := string(afb)

	pcrPlateMinVol := 5.0
	pcrPlateMaxVol := 200.0

	// pcr plate with cooler
	cone := wtype.NewShape(wtype.CylinderShape, "mm", 5.5, 5.5, 15)

	pcrplatewell := wtype.NewLHWell("ul", pcrPlateMaxVol, pcrPlateMinVol, cone, wtype.UWellBottom, 5.5, 5.5, 15, 1.4, "mm")
	pcrplatewell.SetAfVFunc(afs)

	//LiquidLevel model for LL Following: vol_f estimates volume given height
	// vol_f := wutil.Quadratic{A: 0.402, B: 7.069, C: 0.0}
	// pcrplatewell.SetLiquidLevelModel(vol_f)

	return pcrplatewell
}

func makePlateCoords(height float64) wtype.Coordinates3D {
	//standard X/Y size for 96 well plates
	return wtype.Coordinates3D{X: 127.76, Y: 85.48, Z: height}
}

func makePCRPlate() *wtype.Plate {
	return wtype.NewLHPlate("pcrplate", "Unknown", 8, 12, makePlateCoords(15.5), makePCRPlateWell(), 9, 9, 0.0, 0.0, MinimumZHeightPermissableForLVPipetMax)
}

// pcr plate semi-skirted
func makeSemiSkirtedPCRPlate() *wtype.Plate {
	return wtype.NewLHPlate("pcrplate_semi_skirted", "Unknown", 8, 12, makePlateCoords(15.5), makePCRPlateWell(), 9, 9, 0.0, 0.0, 1.0)
}

// 0.2ml strip tubes
func makeStripTube() *wtype.Plate {
	return wtype.NewLHPlate("strip_tubes_0.2ml", "Unknown", 8, 12, makePlateCoords(15.5), makePCRPlateWell(), 9, 9, 0.0, 0.0, 0.0)
}

// pcr plate skirted
func makeSkirtedPCRPlate() *wtype.Plate {
	return wtype.NewLHPlate("pcrplate_skirted", "Unknown", 8, 12, makePlateCoords(15.5), makePCRPlateWell(), 9, 9, 0.0, 0.0, MinimumZHeightPermissableForLVPipetMax)
}

func makeGreinerVBottomPlate() *wtype.Plate {
	// greiner V96 Microplate PS V-Bottom, Clear, Cat Num: 651161

	bottomtype := wtype.VWellBottom
	xdim := 6.2
	ydim := 6.2
	zdim := 11.0
	bottomh := 1.0

	wellxoffset := 9.0 // centre of well to centre of neighbouring well in x direction
	wellyoffset := 9.0 //centre of well to centre of neighbouring well in y direction
	xstart := 0.25     // distance from top left side of plate to first well
	ystart := 0.0      // distance from top left side of plate to first well
	zstart := 3.0      // offset of bottom of deck to bottom of well

	rwshp := wtype.NewShape(wtype.CylinderShape, "mm", 6.2, 6.2, 10.0)
	welltype := wtype.NewLHWell("ul", 230, 10, rwshp, bottomtype, xdim, ydim, zdim, bottomh, "mm")

	plate := wtype.NewLHPlate("GreinerSWVBottom", "Greiner", 8, 12, makePlateCoords(15), welltype, wellxoffset, wellyoffset, xstart, ystart, zstart)

	return plate
}

// Nunc U96 Microplate PolyStyrene Sterile U-Bottom, Clear, Cat Num: 262162
// Source of dimensions: https://www.thermofisher.com/order/catalog/product/262162
func makeNunc96UPlate() *wtype.Plate {

	// These corrections are necessary to subtract from the official (correct) dimensions in order obtain correct pipetting behaviour.
	xstartOffsetCorrection := 11.25
	ystartOffsetCorrection := 7.75

	plateName := "nunc_96_U_PS_Clear"
	manufacturer := "Nunc"

	numberOfRows := 8
	numberOfColumns := 12

	wellShape := wtype.CylinderShape
	bottomtype := wtype.UWellBottom

	dimensionUnit := "mm"

	xdim := 7.1  // G1: diameter at top of well
	ydim := 7.1  // G1: diameter at top of well
	zdim := 10.2 // L: depth of well from top to bottom

	bottomh := 1.0 // ?

	minVolume := 20.0
	maxVolume := 250.0

	volUnit := "ul"

	wellxoffset := 9.0                       // K: centre of well to centre of neighbouring well in x direction
	wellyoffset := 9.0                       // K?: centre of well to centre of neighbouring well in y direction
	xstart := 10.75 - xstartOffsetCorrection // H - (G1/2): distance from top left side of plate to first well (looks like this value does not reflect reality and has an offest applied)
	ystart := 7.75 - ystartOffsetCorrection  // J - (G1/2):  distance from top left side of plate to first well (looks like this value does not reflect reality and has an offest applied)
	zstart := 3.0                            // F - L: offset of bottom of deck to bottom of well
	overallHeight := 14.4                    // F: height of plate

	nunc96UShape := wtype.NewShape(wellShape, dimensionUnit, xdim, ydim, zdim)

	welltype := wtype.NewLHWell(volUnit, maxVolume, minVolume, nunc96UShape, bottomtype, xdim, ydim, zdim, bottomh, dimensionUnit)

	plate := wtype.NewLHPlate(plateName, manufacturer, numberOfRows, numberOfColumns, makePlateCoords(overallHeight), welltype, wellxoffset, wellyoffset, xstart, ystart, zstart)

	return plate
}

// http://fluidx.eu/0.7ml%2c-96-well-format-2d-barcoded-jacket-tube-with-external-thread.html
func makeFluidX700ulTube() *wtype.LHWell {

	wellShape := wtype.CylinderShape
	bottomtype := wtype.VWellBottom
	dimensionUnit := "mm"
	xdim := 6.35 // G1: diameter at top of well
	ydim := 6.35 // G1: diameter at top of well
	zdim := 26.1 // L: depth of well from top to bottom

	bottomh := 1.0 // ?

	minVolume := 20.0
	maxVolume := 525.0

	volUnit := "ul"

	fluidXTubeShape := wtype.NewShape(wellShape, dimensionUnit, xdim, ydim, zdim)

	welltype := wtype.NewLHWell(volUnit, maxVolume, minVolume, fluidXTubeShape, bottomtype, xdim, ydim, zdim, bottomh, dimensionUnit)

	return welltype
}

// http://fluidx.eu/0.7ml%2c-96-well-format-2d-barcoded-jacket-tube-with-external-thread.html
func makeFluidX700ulPlate() *wtype.Plate {

	// no literature values for these

	// These corrections are necessary to subtract from the official (correct) dimensions in order obtain correct pipetting behaviour.
	xstartOffsetCorrection := 11.25
	ystartOffsetCorrection := 7.75

	plateName := "FluidX700ulTubes"

	manufacturer := "FluidX"

	welltype := makeFluidX700ulTube()

	numberOfRows := 8
	numberOfColumns := 12

	wellxoffset := 9.0                               // K: centre of well to centre of neighbouring well in x direction
	wellyoffset := 9.0                               // K?: centre of well to centre of neighbouring well in y direction
	xstart := 10.75 - xstartOffsetCorrection         // H - (G1/2): distance from top left side of plate to first well (looks like this value does not reflect reality and has an offest applied)
	ystart := 7.75 - ystartOffsetCorrection          // J - (G1/2):  distance from top left side of plate to first well (looks like this value does not reflect reality and has an offest applied)
	zstart := MinimumZHeightPermissableForLVPipetMax // F - L: offset of bottom of deck to bottom of well;
	overallHeight := welltype.ZDim() + zstart        // F: height of plate

	plate := wtype.NewLHPlate(plateName, manufacturer, numberOfRows, numberOfColumns, makePlateCoords(overallHeight), welltype, wellxoffset, wellyoffset, xstart, ystart, zstart)

	return plate
}

func makeGreinerFlatBottomBlackPlate() *wtype.Plate {
	// shallow round well flat bottom 96
	rwshp := wtype.NewShape(wtype.CylinderShape, "mm", 8.2, 8.2, 11)
	roundwell96 := wtype.NewLHWell("ul", 340, 25, rwshp, 0, 8.2, 8.2, 11, 1.0, "mm")
	plate := wtype.NewLHPlate("greiner96Black", "greiner", 8, 12, makePlateCoords(15), roundwell96, 9, 9, 0.0, 0.0, 2.0)
	return plate
}

// Onewell SBS format Agarplate with colonies on shallowriser (50ml agar) very high res
func makeHighResplateforPicking() *wtype.Plate {

	bottomtype := wtype.FlatWellBottom
	xdim := 1.4 // of well
	ydim := 1.4
	zdim := 7.0
	bottomh := 0.5

	wellxoffset := 1.55 // centre of well to centre of neighbouring well in x direction
	wellyoffset := 1.55 // centre of well to centre of neighbouring well in y direction
	xstart := -3.885    // distance from top left side of plate to first well
	ystart := -3.0      // distance from top left side of plate to first well
	zstart := 3.25      // offset of bottom of deck to bottom of well

	square3150 := wtype.NewShape(wtype.BoxShape, "mm", xdim, ydim, zdim)
	welltype3150 := wtype.NewLHWell("ul", 5, 0.5, square3150, bottomtype, xdim, ydim, zdim, bottomh, "mm")

	// greiner one well with 50ml of agar in
	plate := wtype.NewLHPlate("Agarplateforpicking3150", "Unknown", 45, 70, makePlateCoords(7), welltype3150, wellxoffset, wellyoffset, xstart, ystart, zstart)

	return plate
}

// Nunc™ 1.0 ml DeepWell™ Plates with Shared-Wall Technology Cat Num: 260251
// Source of dimensions: https://www.thermofisher.com/order/catalog/product/260251
func make96DeepWellLowVolumePlate() *wtype.Plate {

	// These corrections are necessary to subtract from the official (correct) dimensions in order obtain correct pipetting behaviour.
	xstartOffsetCorrection := 14.50
	ystartOffsetCorrection := 11.50
	zstartOffsetCorrection := 2.5

	plateName := "Nunc_96_deepwell_1ml"
	manufacturer := "Thermo Fisher"

	numberOfRows := 8
	numberOfColumns := 12

	wellShape := wtype.CylinderShape
	bottomtype := wtype.UWellBottom

	dimensionUnit := "mm"

	xdim := 8.4  // G1: diameter at top of well
	ydim := 8.4  // G1: diameter at top of well
	zdim := 29.1 // L: depth of well from top to bottom

	bottomh := 1.4 // N: bottom of well to resting plane

	minVolume := 10.0
	maxVolume := 1000.0

	volUnit := "ul"

	wellxoffset := 9.0                      // K: centre of well to centre of neighbouring well in x direction
	wellyoffset := 9.0                      // K?: centre of well to centre of neighbouring well in y direction
	xstart := 14.4 - xstartOffsetCorrection // measure the distance from the edge of plate to beginning of first well in x-axis
	ystart := 11.2 - ystartOffsetCorrection // measure the distance from the edge of plate to beginning of first well in x-axis
	zstart := 2.5 - zstartOffsetCorrection  // F - L: offset of bottom of deck to bottom of well
	overallHeight := 31.6                   // F: height of plate

	newWellShape := wtype.NewShape(wellShape, dimensionUnit, xdim, ydim, zdim)

	newWelltype := wtype.NewLHWell(volUnit, maxVolume, minVolume, newWellShape, bottomtype, xdim, ydim, zdim, bottomh, dimensionUnit)

	plate := wtype.NewLHPlate(plateName, manufacturer, numberOfRows, numberOfColumns, makePlateCoords(overallHeight), newWelltype, wellxoffset, wellyoffset, xstart, ystart, zstart)

	return plate
}

// make plate Labcyte384PPStdV
// Lot number 04090140
// Part number P-05525
// Specs retrieved from
// https://www.labcyte.com/media/pdf/SPC-Qualified-Microplate-384PP.pdf
func makeLabcyte384PPStdV() *wtype.Plate {

	// These corrections are necessary to subtract from the official (correct) dimensions in order obtain correct pipetting behaviour.
	xstartOffsetCorrection := 14.50
	ystartOffsetCorrection := 11.50
	zstartOffsetCorrection := 2.5

	plateName := "Labcyte_384PP_StdV"
	manufacturer := "Labcyte"

	numberOfRows := 16
	numberOfColumns := 24

	wellShape := wtype.BoxShape
	bottomtype := wtype.FlatWellBottom

	dimensionUnit := "mm"

	xdim := 3.3   // G1: diameter at top of well
	ydim := 3.3   // G1: diameter at top of well
	zdim := 11.99 // L: depth of well from top to bottom

	bottomh := 0.0 // N: Used to model different well bottom to body shape. This is the height of the bottom part.

	minVolume := 15.0
	maxVolume := 65.0 // well capacity is therefore 50.0

	volUnit := "ul"

	wellxoffset := 4.5                       // K: centre of well to centre of neighbouring well in x direction
	wellyoffset := 4.5                       // K?: centre of well to centre of neighbouring well in y direction
	xstart := 12.13 - xstartOffsetCorrection // measure the distance from the edge of plate to beginning of first well in x-axis
	ystart := 11.2 - ystartOffsetCorrection  // measure the distance from the edge of plate to beginning of first well in x-axis
	zstart := 2.5 - zstartOffsetCorrection   // F - L: offset of bottom of deck to bottom of well
	overallHeight := 14.4                    // F: height of plate

	newWellShape := wtype.NewShape(wellShape, dimensionUnit, xdim, ydim, zdim)

	newWelltype := wtype.NewLHWell(volUnit, maxVolume, minVolume, newWellShape, bottomtype, xdim, ydim, zdim, bottomh, dimensionUnit)

	plate := wtype.NewLHPlate(plateName, manufacturer, numberOfRows, numberOfColumns, makePlateCoords(overallHeight), newWelltype, wellxoffset, wellyoffset, xstart, ystart, zstart)

	return plate
}

// Applied Biosystems, MicroAmp Optical 384-well Reaction Plate; Cat Num: 4309849
// Source of dimensions: https://www.thermofisher.com/order/catalog/product/4309849
func make384wellplateAppliedBiosystems() *wtype.Plate {

	// These corrections are necessary to subtract from the official (correct) dimensions in order obtain correct pipetting behaviour.
	xstartOffsetCorrection := 13.0
	ystartOffsetCorrection := 10.0
	//zstartOffsetCorrection := 2.25
	zstartOffsetCorrection := 0.7 // MIS revised post fix -- still arbitrary because of later correction

	plateName := "AppliedBiosystems_384_MicroAmp_Optical"
	manufacturer := "Applied Biosystems"

	numberOfRows := 16
	numberOfColumns := 24

	wellShape := wtype.CylinderShape
	bottomtype := wtype.FlatWellBottom

	dimensionUnit := "mm"

	xdim := 3.17 // G1: diameter at top of well
	ydim := 3.17 // G1: diameter at top of well
	zdim := 9.09 // L: depth of well from top to bottom

	bottomh := 0.61 // N: bottom of well to resting plane

	minVolume := 10.0
	maxVolume := 45.0

	volUnit := "ul"

	wellxoffset := 4.5                       // K: centre of well to centre of neighbouring well in x direction
	wellyoffset := 4.5                       // K?: centre of well to centre of neighbouring well in y direction
	xstart := 12.13 - xstartOffsetCorrection // measure the distance from the edge of plate to beginning of first well in x-axis
	ystart := 8.99 - ystartOffsetCorrection  // measure the distance from the edge of plate to beginning of first well in x-axis
	zstart := 0 - zstartOffsetCorrection     // F - L: offset of bottom of deck to bottom of well
	overallHeight := 9.7                     // F: height of plate

	newWellShape := wtype.NewShape(wellShape, dimensionUnit, xdim, ydim, zdim)

	newWelltype := wtype.NewLHWell(volUnit, maxVolume, minVolume, newWellShape, bottomtype, xdim, ydim, zdim, bottomh, dimensionUnit)

	plate := wtype.NewLHPlate(plateName, manufacturer, numberOfRows, numberOfColumns, makePlateCoords(overallHeight), newWelltype, wellxoffset, wellyoffset, xstart, ystart, zstart)

	return plate
}

// Pall, Catalogue # 5076.
// AcroPrep™ 384-well Filter Plates, 100 µL
// AcroPrep 384 x 100ul-well collection plate, without top filter plate.
// Source of dimensions: https://shop.pall.com/us/en/laboratory/dna-rna-purification/plant-genomic-dna-purification/acroprep-384-well-filter-plates-100-l-zidgri78lbr
func makeAcroPrep384NoFilter() *wtype.Plate {

	// These corrections are necessary to subtract from the official (correct) dimensions in order obtain correct pipetting behaviour.
	xstartOffsetCorrection := 14.38
	ystartOffsetCorrection := 11.24

	plateName := "AcroPrep384NoFilter"
	manufacturer := "Pall"

	numberOfRows := 16
	numberOfColumns := 24

	wellShape := wtype.BoxShape
	bottomtype := wtype.VWellBottom

	dimensionUnit := "mm"

	xdim := 4.0  // G1: diameter at top of well
	ydim := 4.0  // G1: diameter at top of well
	zdim := 11.4 // L: depth of well from top to bottom

	bottomh := 0.5 // N: bottom of well to resting plane

	minVolume := 4.0
	maxVolume := 80.0

	volUnit := "ul"

	wellxoffset := 4.5                       // K: centre of well to centre of neighbouring well in x direction
	wellyoffset := 4.5                       // K?: centre of well to centre of neighbouring well in y direction
	xstart := 12.15 - xstartOffsetCorrection // measure the distance from the edge of plate to beginning of first well in x-axis
	ystart := 9 - ystartOffsetCorrection     // measure the distance from the edge of plate to beginning of first well in x-axis
	//zstart := 2.5                            // F - L: offset of bottom of deck to bottom of well
	zstart := 1.8         // F - L: offset of bottom of deck to bottom of well
	overallHeight := 14.4 // F: height of plate

	newWellShape := wtype.NewShape(wellShape, dimensionUnit, xdim, ydim, zdim)

	newWelltype := wtype.NewLHWell(volUnit, maxVolume, minVolume, newWellShape, bottomtype, xdim, ydim, zdim, bottomh, dimensionUnit)

	plate := wtype.NewLHPlate(plateName, manufacturer, numberOfRows, numberOfColumns, makePlateCoords(overallHeight), newWelltype, wellxoffset, wellyoffset, xstart, ystart, zstart)

	return plate
}

// Pall, Catalogue # 5076.
// AcroPrep™ 384-well Filter Plates, 100 µL
// AcroPrep 384-well protein filter plate--omega membrane,
// long tip--stacked on top of a 384 x 100ul-well collection plate.
// Source of dimensions: https://shop.pall.com/us/en/laboratory/dna-rna-purification/plant-genomic-dna-purification/acroprep-384-well-filter-plates-100-l-zidgri78lbr
func makeAcroPrep384WithFilter() *wtype.Plate {

	// These corrections are necessary to subtract from the official (correct) dimensions in order obtain correct pipetting behaviour.
	xstartOffsetCorrection := 14.38
	ystartOffsetCorrection := 11.24

	plateName := "AcroPrep384WithFilter"
	manufacturer := "Pall"

	numberOfRows := 16
	numberOfColumns := 24

	wellShape := wtype.BoxShape
	bottomtype := wtype.FlatWellBottom

	dimensionUnit := "mm"

	xdim := 4.0 // G1: diameter at top of well
	ydim := 4.0 // G1: diameter at top of well
	zdim := 9.0 // L: depth of well from top to bottom

	bottomh := 0.5 // N: bottom of well to resting plane

	minVolume := 80.0
	maxVolume := 80.0

	volUnit := "ul"

	wellxoffset := 4.5                       // K: centre of well to centre of neighbouring well in x direction
	wellyoffset := 4.5                       // K?: centre of well to centre of neighbouring well in y direction
	xstart := 12.15 - xstartOffsetCorrection // measure the distance from the edge of plate to beginning of first well in x-axis
	ystart := 9 - ystartOffsetCorrection     // measure the distance from the edge of plate to beginning of first well in x-axis
	//zstart := 18.0                           // F - L: offset of bottom of deck to bottom of well
	zstart := 17.3        // F - L: offset of bottom of deck to bottom of well
	overallHeight := 27.5 // F: height of plate

	newWellShape := wtype.NewShape(wellShape, dimensionUnit, xdim, ydim, zdim)

	newWelltype := wtype.NewLHWell(volUnit, maxVolume, minVolume, newWellShape, bottomtype, xdim, ydim, zdim, bottomh, dimensionUnit)

	plate := wtype.NewLHPlate(plateName, manufacturer, numberOfRows, numberOfColumns, makePlateCoords(overallHeight), newWelltype, wellxoffset, wellyoffset, xstart, ystart, zstart)

	return plate
}

// VbottomGreiner
func makeGreinerVFromSpec() *wtype.LHPlate {
	plateName := "GreinerV96FromSpec"
	manufacturer := "Greiner"

	numberOfRows := 8
	numberOfColumns := 12

	wellShape := wtype.CylinderShape
	bottomtype := wtype.VWellBottom

	dimensionUnit := "mm"

	xdim := 6.96 // G1: diameter at top of well
	ydim := 6.96 // G1: diameter at top of well
	zdim := 9.0  // L: depth of well from top to bottom

	bottomh := 0.0 // distance between top and bottom of well bottom (i.e. thickness of bottom wall)

	minVolume := 50.0
	maxVolume := 335.0

	volUnit := "ul"

	wellxoffset := 9.0 // K: centre of well to centre of neighbouring well in x direction
	wellyoffset := 9.0 // K?: centre of well to centre of neighbouring well in y direction
	xstart := 14.38
	ystart := 11.24
	zstart := 3.7         // F - L: offset of bottom of deck to bottom of well
	overallHeight := 14.6 // F: height of plate

	newWellShape := wtype.NewShape(wellShape, dimensionUnit, xdim, ydim, zdim)

	newWelltype := wtype.NewLHWell(volUnit, maxVolume, minVolume, newWellShape, bottomtype, xdim, ydim, zdim, bottomh, dimensionUnit)

	plate := wtype.NewLHPlate(plateName, manufacturer, numberOfRows, numberOfColumns, makePlateCoords(overallHeight), newWelltype, wellxoffset, wellyoffset, xstart, ystart, zstart)

	return plate
}

// pcrplate
func make4TitudePcrPlateFromSpec() *wtype.LHPlate {
	plateName := "pcrplate_skirted_FromSpec"
	manufacturer := "4titude"

	numberOfRows := 8
	numberOfColumns := 12

	wellShape := wtype.CylinderShape
	bottomtype := wtype.UWellBottom

	dimensionUnit := "mm"

	xdim := 5.50 // G1: diameter at top of well
	ydim := 5.50 // G1: diameter at top of well
	zdim := 15.1 // L: depth of well from top to bottom

	bottomh := 0.5 // distance between top and bottom of well bottom (i.e. thickness of bottom wall)

	minVolume := 5.0
	maxVolume := 200.0

	volUnit := "ul"

	wellxoffset := 9.0 // K: centre of well to centre of neighbouring well in x direction
	wellyoffset := 9.0 // K?: centre of well to centre of neighbouring well in y direction
	xstart := 14.38
	ystart := 11.24
	zstart := 0.5         // F - L: offset of bottom of deck to bottom of well
	overallHeight := 16.1 // F: height of plate

	newWellShape := wtype.NewShape(wellShape, dimensionUnit, xdim, ydim, zdim)

	newWelltype := wtype.NewLHWell(volUnit, maxVolume, minVolume, newWellShape, bottomtype, xdim, ydim, zdim, bottomh, dimensionUnit)

	plate := wtype.NewLHPlate(plateName, manufacturer, numberOfRows, numberOfColumns, makePlateCoords(overallHeight), newWelltype, wellxoffset, wellyoffset, xstart, ystart, zstart)

	return plate
}
