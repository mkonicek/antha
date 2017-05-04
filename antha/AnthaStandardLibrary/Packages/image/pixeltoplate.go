// pixeltoplate.go
// image.go
package image

import (
	"bytes"
	"encoding/json"
	"fmt"
	goimage "image"
	"image/color"
	"image/color/palette"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"

	anthapath "github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/AnthaPath"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wutil"
	"github.com/disintegration/imaging"
	"errors"
)

//-------------------------------------------------------
//GLOBALS
//-------------------------------------------------------

var Chosencolourpalette color.Palette = availablecolours //palette.WebSafe
//var websafe color.Palette = palette.WebSafe
var availablecolours = []color.Color{
	color.RGBA{R: uint8(242), G: uint8(243), B: uint8(242), A: uint8(255)}, //white
	color.RGBA{R: uint8(13), G: uint8(105), B: uint8(171), A: uint8(255)},  //blue
	color.RGBA{R: uint8(245), G: uint8(205), B: uint8(47), A: uint8(255)},  // yellow
	color.RGBA{R: uint8(75), G: uint8(151), B: uint8(74), A: uint8(255)},   // green
	color.RGBA{R: uint8(196), G: uint8(40), B: uint8(27), A: uint8(255)},   // red
	color.RGBA{R: uint8(196), G: uint8(40), B: uint8(27), A: uint8(255)},   // black
}
// map of RGB colour to description for use as key in crossreferencing colour to component in other maps
var Colourcomponentmap = map[color.Color]string{
	color.RGBA{R: uint8(242), G: uint8(243), B: uint8(242), A: uint8(255)}: "white",
	color.RGBA{R: uint8(13), G: uint8(105), B: uint8(171), A: uint8(255)}:  "blue",
	color.RGBA{R: uint8(245), G: uint8(205), B: uint8(47), A: uint8(255)}:  "yellow",
	color.RGBA{R: uint8(75), G: uint8(151), B: uint8(74), A: uint8(255)}:   "green",
	color.RGBA{R: uint8(196), G: uint8(40), B: uint8(27), A: uint8(255)}:   "red",
	color.RGBA{R: uint8(0), G: uint8(0), B: uint8(0), A: uint8(255)}:       "black",
	color.RGBA{R: uint8(0), G: uint8(0), B: uint8(0), A: uint8(0)}:         "transparent",
}
// map of RGB colour to description for use as key in crossreferencing colour to component in other maps
var Neon = map[color.Color]string{
	color.RGBA{R: uint8(0), G: uint8(0), B: uint8(0), A: uint8(255)}:       "black",
	color.RGBA{R: uint8(149), G: uint8(156), B: uint8(161), A: uint8(255)}: "grey",
	color.RGBA{R: uint8(117), G: uint8(51), B: uint8(127), A: uint8(255)}:  "purple",
	color.RGBA{R: uint8(25), G: uint8(60), B: uint8(152), A: uint8(255)}:   "darkblue",
	color.RGBA{R: uint8(0), G: uint8(125), B: uint8(200), A: uint8(255)}:   "blue",
	color.RGBA{R: uint8(0), G: uint8(177), B: uint8(94), A: uint8(255)}:    "green",
	color.RGBA{R: uint8(244), G: uint8(231), B: uint8(0), A: uint8(255)}:   "yellow",
	color.RGBA{R: uint8(255), G: uint8(118), B: uint8(0), A: uint8(255)}:   "orange",
	color.RGBA{R: uint8(255), G: uint8(39), B: uint8(51), A: uint8(255)}:   "red",
	color.RGBA{R: uint8(235), G: uint8(41), B: uint8(123), A: uint8(255)}:  "pink",
	color.RGBA{R: uint8(242), G: uint8(243), B: uint8(242), A: uint8(255)}: "white",
	color.RGBA{R: uint8(0), G: uint8(174), B: uint8(239), A: uint8(255)}:   "Cyan",
	color.RGBA{R: uint8(236), G: uint8(0), B: uint8(140), A: uint8(255)}:   "Magenta",
	//color.RGBA{R: uint8(251), G: uint8(156), B: uint8(110), A: uint8(255)}: "skin",
}
// Image resizing resample filters from disintegration package
var AllResampleFilters = map[string]imaging.ResampleFilter{
	"Cosine": imaging.Cosine, "Welch": imaging.Welch, "Blackman": imaging.Blackman, "Hamming": imaging.Hamming, "Hann": imaging.Hann, "Lanczos": imaging.Lanczos, "Bartlett": imaging.Bartlett, "Guassian": imaging.Gaussian, "BSpline": imaging.BSpline, "CatmullRom": imaging.CatmullRom, "MitchellNetravali": imaging.MitchellNetravali, "Hermite": imaging.Hermite, "Linear": imaging.Linear, "Box": imaging.Box, "NearestNeighbour": imaging.NearestNeighbor,
}
var Emptycolourarray color.Palette
var ProteinPaintboxmap = map[color.Color]string{
	// under visible light

	// Chromogenic proteins
	color.RGBA{R: uint8(70), G: uint8(105), B: uint8(172), A: uint8(255)}:  "BlitzenBlue",
	color.RGBA{R: uint8(27), G: uint8(79), B: uint8(146), A: uint8(255)}:  "DreidelTeal",
	color.RGBA{R: uint8(107), G: uint8(80), B: uint8(140), A: uint8(255)}:  "VirginiaViolet",
	color.RGBA{R: uint8(120), G: uint8(76), B: uint8(190), A: uint8(255)}:  "VixenPurple",
	color.RGBA{R: uint8(77), G: uint8(11), B: uint8(137), A: uint8(255)}:  "TinselPurple",
	color.RGBA{R: uint8(82), G: uint8(35), B: uint8(119), A: uint8(255)}:  "MaccabeePurple",
	color.RGBA{R: uint8(152), G: uint8(76), B: uint8(128), A: uint8(255)}:  "DonnerMagenta",
	color.RGBA{R: uint8(159), G: uint8(25), B: uint8(103), A: uint8(255)}:  "CupidPink",
	color.RGBA{R: uint8(206), G: uint8(89), B: uint8(142), A: uint8(255)}:  "SeraphinaPink",
	color.RGBA{R: uint8(215), G: uint8(96), B: uint8(86), A: uint8(255)}:  "ScroogeOrange",
	color.RGBA{R: uint8(228), G: uint8(110), B: uint8(104), A: uint8(255)}: 	"LeorOrange",

	// fluorescent proteins

	color.RGBA{R: uint8(224), G: uint8(120), B: uint8(240), A: uint8(254)}: "CindylouCFP",
	color.RGBA{R: uint8(224), G: uint8(120), B: uint8(140), A: uint8(255)}: "FrostyCFP",

	// for twinkle B should = uint8(137) but this is the same colour as e.coli so changed it to uint8(138) to avoid error due to duplicate map keys
	color.RGBA{R: uint8(196), G: uint8(183), B: uint8(138), A: uint8(255)}: "TwinkleCFP",
	//color.RGBA{R: uint8(196), G: uint8(183), B: uint8(137), A: uint8(255)}: "TwinkleCFP",
	//color.RGBA{R: uint8(196), G: uint8(183), B: uint8(137), A: uint8(255)}: "TwinkleCFP",
	color.RGBA{R: uint8(251), G: uint8(176), B: uint8(0), A: uint8(255)}:   "YetiYFP",
	color.RGBA{R: uint8(250), G: uint8(210), B: uint8(0), A: uint8(255)}:   "MarleyYFP",
	color.RGBA{R: uint8(255), G: uint8(194), B: uint8(0), A: uint8(255)}:   "CratchitYFP",
	color.RGBA{R: uint8(231), G: uint8(173), B: uint8(0), A: uint8(255)}:   "KringleYFP",
	color.RGBA{R: uint8(222), G: uint8(221), B: uint8(68), A: uint8(255)}:  "CometGFP",
	color.RGBA{R: uint8(209), G: uint8(214), B: uint8(0), A: uint8(255)}:   "DasherGFP",
	color.RGBA{R: uint8(225), G: uint8(222), B: uint8(120), A: uint8(255)}: "IvyGFP",
	color.RGBA{R: uint8(216), G: uint8(231), B: uint8(15), A: uint8(255)}:  "HollyGFP",
	color.RGBA{R: uint8(251), G: uint8(102), B: uint8(79), A: uint8(255)}:  "YukonOFP",
	color.RGBA{R: uint8(215), G: uint8(72), B: uint8(76), A: uint8(255)}:   "RudolphRFP",
	color.RGBA{R: uint8(244), G: uint8(63), B: uint8(150), A: uint8(255)}:  "FresnoRFP",

	// Extended fluorescent proteins
	color.RGBA{R: uint8(248), G: uint8(64), B: uint8(148), A: uint8(255)}:  "CayenneRFP",
	color.RGBA{R: uint8(241), G: uint8(84), B: uint8(152), A: uint8(255)}:  "GuajilloRFP",
	color.RGBA{R: uint8(247), G: uint8(132), B: uint8(179), A: uint8(255)}: "PaprikaRFP",
	color.RGBA{R: uint8(248), G: uint8(84), B: uint8(149), A: uint8(255)}:  "SerranoRFP",
	//color.RGBA{R: uint8(254), G: uint8(253), B: uint8(252), A: uint8(255)}: "EiraCFP",
	color.RGBA{R: uint8(255), G: uint8(255), B: uint8(146), A: uint8(255)}: "BlazeYFP",
	color.RGBA{R: uint8(194), G: uint8(164), B: uint8(72), A: uint8(255)}:  "JuniperGFP",
	color.RGBA{R: uint8(243), G: uint8(138), B: uint8(112), A: uint8(255)}: "TannenGFP",

	// conventional E.coli colour
	color.RGBA{R: uint8(196), G: uint8(183), B: uint8(137), A: uint8(255)}: "E.coli",

	// lacZ expresser (e.g. pUC19) grown on S gal
	//color.RGBA{R: uint8(0), G: uint8(0), B: uint8(0), A: uint8(255)}: "E.coli pUC19 on sgal",
	color.RGBA{R: uint8(1), G: uint8(1), B: uint8(1), A: uint8(255)}: "veryblack",
	// plus white as a blank (or comment out to use EiraCFP)
	color.RGBA{R: uint8(242), G: uint8(243), B: uint8(242), A: uint8(255)}: "verywhite",
}
var UVProteinPaintboxmap = map[color.Color]string{

	// Chromogenic proteins same colours as visible
	color.RGBA{R: uint8(70), G: uint8(105), B: uint8(172), A: uint8(255)}:  "BlitzenBlue",
	color.RGBA{R: uint8(27), G: uint8(79), B: uint8(146), A: uint8(255)}:   "DreidelTeal",
	color.RGBA{R: uint8(107), G: uint8(80), B: uint8(140), A: uint8(255)}:  "VirginiaViolet",
	color.RGBA{R: uint8(120), G: uint8(76), B: uint8(190), A: uint8(255)}:  "VixenPurple",
	color.RGBA{R: uint8(77), G: uint8(11), B: uint8(137), A: uint8(255)}:   "TinselPurple",
	color.RGBA{R: uint8(82), G: uint8(35), B: uint8(119), A: uint8(255)}:   "MaccabeePurple",
	color.RGBA{R: uint8(152), G: uint8(76), B: uint8(128), A: uint8(255)}:  "DonnerMagenta",
	color.RGBA{R: uint8(159), G: uint8(25), B: uint8(103), A: uint8(255)}:  "CupidPink",
	color.RGBA{R: uint8(206), G: uint8(89), B: uint8(142), A: uint8(255)}:  "SeraphinaPink",
	color.RGBA{R: uint8(215), G: uint8(96), B: uint8(86), A: uint8(255)}:   "ScroogeOrange",
	color.RGBA{R: uint8(228), G: uint8(110), B: uint8(104), A: uint8(255)}: "LeorOrange",

	// under UV

	// fluorescent
	color.RGBA{R: uint8(0), G: uint8(254), B: uint8(255), A: uint8(255)}: "CindylouCFP",
	color.RGBA{R: uint8(0), G: uint8(255), B: uint8(255), A: uint8(255)}: "FrostyCFP",
	color.RGBA{R: uint8(27), G: uint8(79), B: uint8(146), A: uint8(255)}: "TwinkleCFP",
	//color.RGBA{R: uint8(27), G: uint8(79), B: uint8(146), A: uint8(255)}: "TwinkleCFP",
	//color.RGBA{R: uint8(27), G: uint8(79), B: uint8(146), A: uint8(255)}:  "TwinkleCFP",
	color.RGBA{R: uint8(253), G: uint8(230), B: uint8(39), A: uint8(255)}: "YetiYFP",
	color.RGBA{R: uint8(236), G: uint8(255), B: uint8(0), A: uint8(255)}:  "MarleyYFP",
	color.RGBA{R: uint8(240), G: uint8(254), B: uint8(0), A: uint8(255)}:  "CratchitYFP",
	color.RGBA{R: uint8(239), G: uint8(255), B: uint8(0), A: uint8(255)}:  "KringleYFP",
	color.RGBA{R: uint8(0), G: uint8(254), B: uint8(0), A: uint8(255)}:    "CometGFP",
	color.RGBA{R: uint8(0), G: uint8(255), B: uint8(0), A: uint8(255)}:    "DasherGFP",
	color.RGBA{R: uint8(0), G: uint8(232), B: uint8(216), A: uint8(255)}:  "IvyGFP",
	color.RGBA{R: uint8(0), G: uint8(255), B: uint8(0), A: uint8(254)}:    "HollyGFP",
	color.RGBA{R: uint8(254), G: uint8(179), B: uint8(18), A: uint8(255)}: "YukonOFP",
	color.RGBA{R: uint8(218), G: uint8(92), B: uint8(69), A: uint8(255)}:  "RudolphRFP",
	color.RGBA{R: uint8(255), G: uint8(0), B: uint8(166), A: uint8(255)}:  "FresnoRFP",

	// Extended fluorescent proteins
	color.RGBA{R: uint8(255), G: uint8(24), B: uint8(138), A: uint8(255)}: "CayenneRFP",
	color.RGBA{R: uint8(255), G: uint8(8), B: uint8(138), A: uint8(255)}:  "GuajilloRFP",
	color.RGBA{R: uint8(252), G: uint8(65), B: uint8(136), A: uint8(255)}: "PaprikaRFP",
	color.RGBA{R: uint8(254), G: uint8(23), B: uint8(127), A: uint8(255)}: "SerranoRFP",
	//color.RGBA{R: uint8(173), G: uint8(253), B: uint8(218), A: uint8(255)}: "EiraCFP",
	color.RGBA{R: uint8(254), G: uint8(255), B: uint8(83), A: uint8(255)}: "BlazeYFP",
	color.RGBA{R: uint8(0), G: uint8(231), B: uint8(162), A: uint8(255)}:  "JuniperGFP",
	color.RGBA{R: uint8(179), G: uint8(119), B: uint8(57), A: uint8(255)}: "TannenGFP",

	// conventional E.coli colour
	color.RGBA{R: uint8(196), G: uint8(183), B: uint8(137), A: uint8(255)}: "E.coli",
	color.RGBA{R: uint8(1), G: uint8(1), B: uint8(1), A: uint8(255)}:       "E.coli pUC19 on sgal",
	color.RGBA{R: uint8(242), G: uint8(243), B: uint8(242), A: uint8(255)}: "verywhite",
}
var ProteinPaintboxSubsetmap = map[color.Color]string{
	// under visible light

	// Chromogenic proteins
	color.RGBA{R: uint8(70), G: uint8(105), B: uint8(172), A: uint8(255)}: "BlitzenBlue",
	color.RGBA{R: uint8(27), G: uint8(79), B: uint8(146), A: uint8(255)}:  "DreidelTeal",
	/*color.RGBA{R: uint8(107), G: uint8(80), B: uint8(140), A: uint8(255)}:  "VirginiaViolet",
	color.RGBA{R: uint8(120), G: uint8(76), B: uint8(190), A: uint8(255)}:  "VixenPurple",*/
	color.RGBA{R: uint8(77), G: uint8(11), B: uint8(137), A: uint8(255)}: "TinselPurple",
	/*color.RGBA{R: uint8(82), G: uint8(35), B: uint8(119), A: uint8(255)}:   "MaccabeePurple",
	color.RGBA{R: uint8(152), G: uint8(76), B: uint8(128), A: uint8(255)}:  "DonnerMagenta",*/
	color.RGBA{R: uint8(159), G: uint8(25), B: uint8(103), A: uint8(255)}: "CupidPink",
	//	color.RGBA{R: uint8(206), G: uint8(89), B: uint8(142), A: uint8(255)}:  "SeraphinaPink",
	//color.RGBA{R: uint8(215), G: uint8(96), B: uint8(86), A: uint8(255)}: "ScroogeOrange",
	//color.RGBA{R: uint8(228), G: uint8(110), B: uint8(104), A: uint8(255)}: "LeorOrange",

	// fluorescent proteins

	//	color.RGBA{R: uint8(224), G: uint8(120), B: uint8(240), A: uint8(255)}:  "CindylouCFP",
	//color.RGBA{R: uint8(224), G: uint8(120), B: uint8(140), A: uint8(255)}: "FrostyCFP",
	color.RGBA{R: uint8(196), G: uint8(183), B: uint8(137), A: uint8(255)}: "TwinkleCFP",
	/*
		// for twinkle B should = uint8(137) but this is the same colour as e.coli so changed it to uint8(138) to avoid error due to duplicate map keys
		color.RGBA{R: uint8(196), G: uint8(183), B: uint8(138), A: uint8(255)}: "TwinkleCFP",
		//color.RGBA{R: uint8(196), G: uint8(183), B: uint8(137), A: uint8(255)}: "TwinkleCFP",
		//color.RGBA{R: uint8(196), G: uint8(183), B: uint8(137), A: uint8(255)}: "TwinkleCFP",*/
	color.RGBA{R: uint8(251), G: uint8(176), B: uint8(0), A: uint8(255)}: "YetiYFP",
	color.RGBA{R: uint8(250), G: uint8(210), B: uint8(0), A: uint8(255)}: "MarleyYFP",
	color.RGBA{R: uint8(255), G: uint8(194), B: uint8(0), A: uint8(255)}: "CratchitYFP",
	color.RGBA{R: uint8(231), G: uint8(173), B: uint8(0), A: uint8(255)}: "KringleYFP",
	//color.RGBA{R: uint8(222), G: uint8(221), B: uint8(68), A: uint8(255)}: "CometGFP",
	// new green
	//color.RGBA{R: uint8(105), G: uint8(189), B: uint8(67), A: uint8(255)}: "green",
	//105 189 67 255
	color.RGBA{R: uint8(209), G: uint8(214), B: uint8(0), A: uint8(255)}:   "DasherGFP",
	color.RGBA{R: uint8(225), G: uint8(222), B: uint8(120), A: uint8(255)}: "IvyGFP",
	//color.RGBA{R: uint8(216), G: uint8(231), B: uint8(15), A: uint8(255)}:     "HollyGFP",
	color.RGBA{R: uint8(251), G: uint8(102), B: uint8(79), A: uint8(255)}: "YukonOFP",
	color.RGBA{R: uint8(215), G: uint8(72), B: uint8(76), A: uint8(255)}:  "RudolphRFP",
	color.RGBA{R: uint8(244), G: uint8(63), B: uint8(150), A: uint8(255)}: "FresnoRFP",

	// Extended fluorescent proteins
	/*color.RGBA{R: uint8(248), G: uint8(64), B: uint8(148), A: uint8(255)}:  "CayenneRFP",
	color.RGBA{R: uint8(241), G: uint8(84), B: uint8(152), A: uint8(255)}:  "GuajilloRFP",
	color.RGBA{R: uint8(247), G: uint8(132), B: uint8(179), A: uint8(255)}: "PaprikaRFP",
	color.RGBA{R: uint8(248), G: uint8(84), B: uint8(149), A: uint8(255)}:  "SerranoRFP",
	color.RGBA{R: uint8(254), G: uint8(253), B: uint8(252), A: uint8(255)}: "EiraCFP",*/
	color.RGBA{R: uint8(255), G: uint8(255), B: uint8(146), A: uint8(255)}: "BlazeYFP",
	color.RGBA{R: uint8(194), G: uint8(164), B: uint8(72), A: uint8(255)}:  "JuniperGFP",
	//color.RGBA{R: uint8(243), G: uint8(138), B: uint8(112), A: uint8(255)}: "TannenGFP",
	//*/
	// conventional E.coli colour
	//color.RGBA{R: uint8(196), G: uint8(183), B: uint8(137), A: uint8(255)}: "E.coli",

	// lacZ expresser (e.g. pUC19) grown on S gal
	//color.RGBA{R: uint8(0), G: uint8(0), B: uint8(0), A: uint8(255)}: "E.coli pUC19 on sgal",
	color.RGBA{R: uint8(1), G: uint8(1), B: uint8(1), A: uint8(255)}: "veryblack",

	// plus white as a blank (or comment out to use EiraCFP)
	//color.RGBA{R: uint8(242), G: uint8(243), B: uint8(242), A: uint8(255)}: "verywhite",
}


//----------------------------------------------------------------------------
//image input/output manipulation
//----------------------------------------------------------------------------

//OpenFile will take a wtype.File object and return the contents as image.NRGBA
func OpenFile (file wtype.File)( nrgba *goimage.NRGBA, err error){

	data, err := file.ReadAll()
	if err != nil {
		return nil, err
	}

	buf := bytes.NewReader(data)

	img, err := imaging.Decode(buf)
	if err != nil {
		return nil, err
	}

	nrgba = imaging.Clone(img)
	return nrgba, nil
}

// export image to file
// image format is derived from filename extension
func Export(img *goimage.NRGBA) (file wtype.File, err error) {

	var imageFormat imaging.Format
	var buf bytes.Buffer

	err = imaging.Encode(&buf, img, imageFormat)
	if err != nil {
		return
	}

	err = file.WriteAll(buf.Bytes())
	if err != nil {
		return
	}

	return
}

//----------------------------------------------------------------------------
//Palette Libraries manipulation
//----------------------------------------------------------------------------

//AvailablePalettes will return all the (hardcoded) pallettes in the library. The keys are the pallette names
func AvailablePalettes() (availablepalettes map[string]color.Palette) {

	availablepalettes = make(map[string]color.Palette)

	availablepalettes["Palette1"] = palettefromMap(Colourcomponentmap) //Chosencolourpalette,
	availablepalettes["Neon"] = palettefromMap(Neon)
	availablepalettes["WebSafe"] = palette.WebSafe //websafe,
	availablepalettes["Plan9"] = palette.Plan9
	availablepalettes["ProteinPaintboxVisible"] = palettefromMap(ProteinPaintboxmap)
	availablepalettes["ProteinPaintboxUV"] = palettefromMap(UVProteinPaintboxmap)
	availablepalettes["ProteinPaintboxSubset"] = palettefromMap(ProteinPaintboxSubsetmap)
	availablepalettes["Gray"] = MakeGreyScalePalette()
	availablepalettes["None"] = Emptycolourarray

	if _, err := os.Stat(filepath.Join(anthapath.Path(), "testcolours.json")); err == nil {
		invmap, err := MakelatestcolourMap(filepath.Join(anthapath.Path(), "testcolours.json"))
		if err != nil {
			panic(err.Error())
		}
		availablepalettes["inventory"] = palettefromMap(invmap)
	}

	if _, err := os.Stat(filepath.Join(anthapath.Path(), "UVtestcolours.json")); err == nil {
		uvinvmap, err := MakelatestcolourMap(filepath.Join(anthapath.Path(), "UVtestcolours.json"))
		if err != nil {
			panic(err.Error())
		}
		availablepalettes["UVinventory"] = palettefromMap(uvinvmap)
	}
	return
}

//AvailableComponentMaps returns available palettes in a map
func AvailableComponentmaps() (componentmaps map[string]map[color.Color]string) {
	componentmaps = make(map[string]map[color.Color]string)
	componentmaps["Palette1"] = Colourcomponentmap
	componentmaps["Neon"] = Neon
	componentmaps["ProteinPaintboxVisible"] = ProteinPaintboxmap
	componentmaps["ProteinPaintboxUV"] = UVProteinPaintboxmap
	componentmaps["ProteinPaintboxSubset"] = ProteinPaintboxSubsetmap

	if _, err := os.Stat(filepath.Join(anthapath.Path(), "testcolours.json")); err == nil {
		invmap, err := MakelatestcolourMap(filepath.Join(anthapath.Path(), "testcolours.json"))
		if err != nil {
			panic(err.Error())
		}

		componentmaps["inventory"] = invmap
	}
	if _, err := os.Stat(filepath.Join(anthapath.Path(), "UVtestcolours.json")); err == nil {
		uvinvmap, err := MakelatestcolourMap(filepath.Join(anthapath.Path(), "UVtestcolours.json"))
		if err != nil {
			panic(err.Error())
		}

		componentmaps["UVinventory"] = uvinvmap
	}
	return
}

//VisibleEquivalentMaps returns just the proteinPaintboxUV color map
func Visibleequivalentmaps() map[string]map[color.Color]string {
	visibleequivalentmaps := make(map[string]map[color.Color]string)
	visibleequivalentmaps["ProteinPaintboxUV"] = ProteinPaintboxmap

	if _, err := os.Stat(filepath.Join(anthapath.Path(), "testcolours.json")); err == nil {
		invmap, err := MakelatestcolourMap(filepath.Join(anthapath.Path(), "testcolours.json"))
		if err != nil {
			panic(err.Error())
		}
		visibleequivalentmaps["UVinventory"] = invmap
	}
	return visibleequivalentmaps
}

//PaletteFromMap returns a palette of all colors in it given a color map.
func palettefromMap(colourmap map[color.Color]string) (palette color.Palette) {

	array := make([]color.Color, 0)

	for key, _ := range colourmap {

		array = append(array, key)
	}

	palette = array
	return

}

//MakelatestcolourMap will take a json file and return map with color objects as keys to their ID.
func MakelatestcolourMap(jsonmapfilename string) (colourtostringmap map[color.Color]string, err error) {
	var stringtonrgbamap *map[string]color.NRGBA = &map[string]color.NRGBA{}

	data, err := ioutil.ReadFile(jsonmapfilename)

	if err != nil {
		return colourtostringmap, err
	}

	err = json.Unmarshal(data, stringtonrgbamap)
	if err != nil {
		return colourtostringmap, err
	}

	stringtocolourmap := make(map[string]color.Color)
	for key, value := range *stringtonrgbamap {
		stringtocolourmap[key] = value
	}

	colourtostringmap, err = reversestringtopalettemap(stringtocolourmap)

	return colourtostringmap, err
}

//MakeGreyScalePalette will return a palette of grey shades
func MakeGreyScalePalette() (graypalette []color.Color) {

	graypalette = make([]color.Color, 0)
	var shadeofgray color.Gray
	for i := 0; i < 256; i++ {
		shadeofgray = color.Gray{Y: uint8(i)}
		graypalette = append(graypalette, shadeofgray)
	}

	return
}

//ReversePaletteMap reverses the keys and values in a colormap. The color names become keys for the color objects
func reversepalettemap(colourmap map[color.Color]string) (stringmap map[string]color.Color, err error) {

	stringmap = make(map[string]color.Color, len(colourmap))

	for key, value := range colourmap {

		_, ok := stringmap[value]
		if ok == true {
			alreadyinthere := stringmap[value]

			err = fmt.Errorf("attempt to add value", key, "for key", value, "to stringmap", stringmap, "failed due to duplicate entry", alreadyinthere)
		} else {
			stringmap[value] = key
		}
		// fmt.Println("key:", key, "value", value)
	}
	return
}

//ReverseStringToPaletteMap reverses the keys and values in a colormap. The color objects become keys for the color names
func reversestringtopalettemap(stringmap map[string]color.Color) (colourmap map[color.Color]string, err error) {

	colourmap = make(map[color.Color]string, len(stringmap))

	for key, value := range stringmap {

		_, ok := colourmap[value]
		if ok == true {
			alreadyinthere := colourmap[value]

			err = fmt.Errorf("attempt to add value", key, "for key", value, "to colourmap", colourmap, "failed due to duplicate entry", alreadyinthere)
		} else {
			colourmap[value] = key
		}
		// fmt.Println("key:", key, "value", value)
	}
	return
}

//MakeSubMapFromMap extracts colors from a color library Given colour names.
func MakeSubMapfromMap(existingmap map[color.Color]string, colournames []string) (newmap map[color.Color]string) {

	newmap = make(map[color.Color]string, 0)

	reversedmap, err := reversepalettemap(existingmap)

	if err != nil {
		panic("can't reverse this colour map" + err.Error())
	}

	for _, colourname := range colournames {

		colour := reversedmap[colourname]

		if colour != nil {
			newmap[colour] = colourname
		}
	}

	return
}

//MakeSubPallette extracts colors from a palette given a slice of colour names
func MakeSubPallette(palettename string, colournames []string) (subpalette color.Palette) {
	palettemap := AvailableComponentmaps()[palettename]

	submap := MakeSubMapfromMap(palettemap, colournames)

	subpalette = palettefromMap(submap)

	return

}

//RemoveDuplicatesKeysfromMap will loop over a map of colors to find and delete duplicates. Entries with duplicate keys are deleted.
func RemoveDuplicatesKeysfromMap(elements map[string]color.Color) map[string]color.Color {
	// Use map to record duplicates as we find them.
	encountered := map[string]bool{}
	result := make(map[string]color.Color, 0)

	for key, v := range elements {

		if encountered[key] == true {
			// Do not add duplicate.
		} else {
			// Record this element as an encountered element.
			encountered[key] = true
			// Append to result slice.
			result[key] = v
		}
	}
	// Return the new slice.
	return result
}

//RemoveDuplicatesKeysValuesfromMap will loop over a map of colors to find and delete duplicates. Entries with duplicate values are deleted.
func RemoveDuplicatesValuesfromMap(elements map[string]color.Color) map[string]color.Color {
	// Use map to record duplicates as we find them.
	encountered := map[color.Color]bool{}
	result := make(map[string]color.Color, 0)

	for key, v := range elements {

		if encountered[v] == true {
			// Do not add duplicate.
		} else {
			// Record this element as an encountered element.
			encountered[v] = true
			// Append to result slice.
			result[key] = v
		}
	}
	// Return the new slice.
	return result
}

//----------------------------------------------------------------------------
//Image manipulations
//----------------------------------------------------------------------------

//PaletteFromColorArray make a palette from a slice of colors.
func paletteFromColorarray(colors []color.Color) (palette color.Palette) {

	var newpalette color.Palette

	newpalette = colors

	palette = newpalette

	//palette = &newpalette
	return
}

//ColourtoCMYK will convert an object with the color interface to color.CYMK
func ColourtoCMYK(colour color.Color) (cmyk color.CMYK) {
	// fmt.Println("colour", colour)
	r, g, b, _ := colour.RGBA()
	cmyk.C, cmyk.M, cmyk.Y, cmyk.K = color.RGBToCMYK(uint8(r), uint8(g), uint8(b))
	return
}

//ColourtoGrayscale will convert an object with the color interface to color.Gray
func ColourtoGrayscale(colour color.Color) (gray color.Gray) {
	r, g, b, _ := colour.RGBA()
	gray.Y = uint8((0.2126 * float64(r)) + (0.7152 * float64(g)) + (0.0722 * float64(b)))
	return
}

//Posterize will posterize an image. This refers to changing an image to use only a small number of different tones.
func Posterize(img *goimage.NRGBA, levels int) (posterized *goimage.NRGBA, err error) {

	//We cannot posterize with only one level.
	if levels == 1 {
		return nil, errors.New("Cannot posterize with only one level.")
	}

	var newcolor color.NRGBA
	numberofAreas := 256 / (levels)
	numberofValues := 255 / (levels - 1)

	posterized = imaging.Clone(img)

	for x := 0; x < img.Bounds().Dx(); x++ {
		for y := 0; y < img.Bounds().Dy(); y++ {
			var rnew float64
			var gnew float64
			var bnew float64

			rgb := img.At(x, y)
			r, g, b, a := rgb.RGBA()

			if r == 0 {
				rnew = 0
			} else {
				rfloat := (float64(r/256) / float64(numberofAreas))

				rinttemp, err := wutil.RoundDown(rfloat)
				if err != nil {
					panic(err)
				}
				rnew = float64(rinttemp) * float64(numberofValues)
			}
			if g == 0 {
				gnew = 0
			} else {
				gfloat := (float64(g/256) / float64(numberofAreas))

				ginttemp, err := wutil.RoundDown(gfloat)
				if err != nil {
					panic(err)
				}
				gnew = float64(ginttemp) * float64(numberofValues)
			}
			if b == 0 {
				bnew = 0
			} else {
				bfloat := (float64(b/256) / float64(numberofAreas))

				binttemp, err := wutil.RoundDown(bfloat)
				if err != nil {
					panic(err)
				}
				bnew = float64(binttemp) * float64(numberofValues)
			}
			newcolor.A = uint8(a)

			rint, err := wutil.RoundDown(rnew)

			if err != nil {
				panic(err)
			}
			newcolor.R = uint8(rint)
			gint, err := wutil.RoundDown(gnew)
			if err != nil {
				panic(err)
			}
			newcolor.G = uint8(gint)
			bint, err := wutil.RoundDown(bnew)
			if err != nil {
				panic(err)
			}
			newcolor.B = uint8(bint)

			posterized.Set(x, y, newcolor)

		}
	}

	return posterized, nil
}

//ResizeImageToPlate will resize an image to fit the number of wells on a plate. We treat wells as pixels.
func ResizeImagetoPlate(img *goimage.NRGBA, plate *wtype.LHPlate, algorithm imaging.ResampleFilter, rotate bool) (plateimage *goimage.NRGBA) {

	if img.Bounds().Dy() != plate.WellsY() {

		if rotate {
			img = imaging.Rotate270(img)
		}
		var aspectratio float64 = float64(img.Bounds().Dx()) / float64(img.Bounds().Dy())
		if aspectratio <= float64(plate.WellsX())/float64(plate.WellsY()) {
			plateimage = imaging.Resize(img, 0, plate.WlsY, algorithm)
		} else {
			plateimage = imaging.Resize(img, plate.WlsX, 0, algorithm)
		}
		//plateimages = append(plateimages,plateimage)
	} else {
		// if the size is the same, simply return the image
		plateimage = img
	}
	return

}

//ResizeImageToPlateAutoRotate will resize an image to fit the number of wells on a plate and if the image is in portrait the image will be rotated by 270 degrees to optimise resolution on the plate.
func ResizeImagetoPlateAutoRotate(img *goimage.NRGBA, plate *wtype.LHPlate, algorithm imaging.ResampleFilter) (plateimage *goimage.NRGBA) {

	if img.Bounds().Dy() != plate.WellsY() {
		// fmt.Println("hey we're not so different", img.Bounds().Dy(), plate.WellsY())
		// have the option of changing the resize algorithm here

		if img.Bounds().Dy() > img.Bounds().Dx() {
			// fmt.Println("Auto Rotating image")
			img = imaging.Rotate270(img)
		}
		var aspectratio float64 = float64(img.Bounds().Dx()) / float64(img.Bounds().Dy())
		if aspectratio <= float64(plate.WellsX())/float64(plate.WellsY()) {
			plateimage = imaging.Resize(img, 0, plate.WlsY, algorithm)
		} else {
			plateimage = imaging.Resize(img, plate.WlsX, 0, algorithm)
		}
		//plateimages = append(plateimages,plateimage)
	} else {
		// fmt.Println("i'm the same!!!")
		plateimage = toNRGBA(img)
	}
	return

}

//CheckAllResizeAlgorithms will use the different algorithms in a algorithm library to resize an image to a given platetype.
func CheckAllResizealgorithms(img *goimage.NRGBA, plate *wtype.LHPlate, rotate bool, algorithms map[string]imaging.ResampleFilter) (plateimage *goimage.NRGBA) {

	// Colour palette to use // this would relate to a map of components of these available colours in factory
	//availablecolours := chosencolourpalette //palette.WebSafe

	//var plateimages []image.Image

	for _ , algorithm := range algorithms {

		if rotate {
			img = imaging.Rotate270(img)
		}

		if img.Bounds().Dy() != plate.WellsY() {
			// fmt.Println("hey we're not so different", img.Bounds().Dy(), plate.WellsY())
			// have the option of changing the resize algorithm here
			var aspectratio float64 = float64(img.Bounds().Dx()) / float64(img.Bounds().Dy())
			if aspectratio <= float64(plate.WellsX())/float64(plate.WellsY()) {
				plateimage = imaging.Resize(img, 0, plate.WlsY, algorithm)
			} else {
				plateimage = imaging.Resize(img, plate.WlsX, 0, algorithm)
			}
			//plateimages = append(plateimages,plateimage)
		} else {
			// fmt.Println("i'm the same!!!")
			plateimage = toNRGBA(img)
		}

	}

	return
}

//MakePalleteFromImage will make a color Palette from an image resized to fit a given plate type.
func MakePalleteFromImage(img *goimage.NRGBA, plate *wtype.LHPlate, rotate bool) (newpallette color.Palette) {

	plateimage := ResizeImagetoPlate(img, plate, imaging.CatmullRom, rotate)

	colourarray := make([]color.Color, 0)

	// Find out colour at each position:
	for y := 0; y < plateimage.Bounds().Dy(); y++ {
		for x := 0; x < plateimage.Bounds().Dx(); x++ {
			// colour or pixel in RGB
			colour := plateimage.At(x, y)
			colourarray = append(colourarray, colour)

		}
	}

	newpallette = paletteFromColorarray(colourarray)

	return
}

//MakeSmallPalleteFromImage will make a color Palette from an image resized to fit a given plate type using Plan9 Palette
func MakeSmallPalleteFromImage(img *goimage.NRGBA, plate *wtype.LHPlate, rotate bool) (newpallette color.Palette) {

	plateimage := ResizeImagetoPlate(img, plate, imaging.CatmullRom, rotate)
	//image, _ := imaging.Open(imagefilename)

	//plateimage := imaging.Clone(image)

	// use Plan9 as pallette for first round to keep number of colours down to a manageable level

	chosencolourpalette := AvailablePalettes()["Plan9"]

	colourmap := make(map[color.Color]bool, 0)

	// Find out colour at each position:
	for y := 0; y < plateimage.Bounds().Dy(); y++ {
		for x := 0; x < plateimage.Bounds().Dx(); x++ {
			// colour or pixel in RGB
			colour := plateimage.At(x, y)

			if colour != nil {

				colour = chosencolourpalette.Convert(colour)
				_, ok := colourmap[colour]
				// change colour to colour from a palette
				if !ok {
					colourmap[colour] = true
				}

			}
		}
	}

	newcolourarray := make([]color.Color, 0)

	for colour, _ := range colourmap {
		newcolourarray = append(newcolourarray, colour)
	}

	newpallette = paletteFromColorarray(newcolourarray)

	return
}

//ImagetoPlatelayout will take an image, plate type, and palette and return a map of well position to colors.
// create a map of pixel to plate position from processing a given image with a chosen colour palette.
// It's recommended to use at least 384 well plate
// if autorotate == true, rotate is overridden
func ImagetoPlatelayout(img *goimage.NRGBA, plate *wtype.LHPlate, chosencolourpalette *color.Palette, rotate bool, autorotate bool) (wellpositiontocolourmap map[string]color.Color, plateimage *goimage.NRGBA) {

	if autorotate {
		plateimage = ResizeImagetoPlateAutoRotate(img, plate, imaging.CatmullRom)
	} else {
		plateimage = ResizeImagetoPlate(img, plate, imaging.CatmullRom, rotate)
	}
	// make map of well position to colour: (array for time being)

	wellpositionarray := make([]string, 0)
	colourarray := make([]color.Color, 0)
	wellpositiontocolourmap = make(map[string]color.Color, 0)

	// Find out colour at each position:
	for y := 0; y < plateimage.Bounds().Dy(); y++ {
		for x := 0; x < plateimage.Bounds().Dx(); x++ {
			// colour or pixel in RGB
			colour := plateimage.At(x, y)
			// fmt.Println("x,y,colour, palette", x, y, colour, chosencolourpalette)

			if colour != nil {

				if chosencolourpalette != nil && chosencolourpalette != &Emptycolourarray && len([]color.Color(*chosencolourpalette)) > 0 {
					// change colour to colour from a palette
					colour = chosencolourpalette.Convert(colour)
					// fmt.Println("x,y,colour", x, y, colour)
					plateimage.Set(x, y, colour)
				}
				// equivalent well position
				wellposition := wutil.NumToAlpha(y+1) + strconv.Itoa(x+1)
				wellpositionarray = append(wellpositionarray, wellposition)
				wellpositiontocolourmap[wellposition] = colour

				colourarray = append(colourarray, colour)
			}
		}
	}

	return
}

//PrintFPImagePreview will take an image, plate type, and colors from visiblemap/uvmap and (use to) save the resulting processed image.
func PrintFPImagePreview(img *goimage.NRGBA, plate *wtype.LHPlate, rotate bool, visiblemap, uvmap map[color.Color]string) {

	plateimage := ResizeImagetoPlate(img, plate, imaging.CatmullRom, rotate)

	uvpalette := palettefromMap(uvmap)

	// Find out colour at each position under UV:
	for y := 0; y < plateimage.Bounds().Dy(); y++ {
		for x := 0; x < plateimage.Bounds().Dx(); x++ {
			// colour or pixel in RGB
			colour := plateimage.At(x, y)

			// change colour to colour from a palette
			uvcolour := uvpalette.Convert(colour)

			plateimage.Set(x, y, uvcolour)

		}
	}

	// repeat for visible

	// Find out colour at each position under visible light:
	for y := 0; y < plateimage.Bounds().Dy(); y++ {
		for x := 0; x < plateimage.Bounds().Dx(); x++ {
			// colour or pixel in RGB

			colour := plateimage.At(x, y)
			r, g, b, a := colour.RGBA()
			rgba := color.RGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: uint8(a)}
			// fmt.Println("colour", colour)
			// fmt.Println("visiblemap", visiblemap)
			// fmt.Println("uvmap", uvmap)
			colourstring := uvmap[rgba]
			// fmt.Println("colourstring", colourstring)
			// change colour to colour of same cell + fluorescent protein under visible light
			stringkeymap, err := reversepalettemap(visiblemap)
			if err != nil {
				panic(err)
			}
			// fmt.Println("stringkeymap", stringkeymap)
			viscolour, ok := stringkeymap[colourstring]
			if ok != true {
				errmessage := fmt.Sprintln("colourstring", colourstring, "not found in map", stringkeymap, "len", len(stringkeymap))
				panic(errmessage)
			}
			// fmt.Println("viscolour", viscolour)
			plateimage.Set(x, y, viscolour)

		}
	}

	return
}

//MakeStringToComponentMap make a map linking LHcomponents to colours, and assign them to well
func MakestringtoComponentMap(keys []string, componentlist []*wtype.LHComponent) (componentmap map[string]*wtype.LHComponent, err error) {

	componentmap = make(map[string]*wtype.LHComponent, 0)
	var previouserror error = nil
	for i, key := range keys {
		for j, component := range componentlist {

			if component.CName == key {
				componentmap[key] = component
				break
			}
			if i == len(keys) {

				err = fmt.Errorf(previouserror.Error(), "+", "no key and component found for", keys, key)
			}
			if j == len(componentlist) {

				err = fmt.Errorf(previouserror.Error(), "+", "no key and component found for", component.CName)
			}
		}
	}
	return
}

// toNRGBA is used internally to convert any image type to NRGBA.
func toNRGBA(img goimage.Image) *goimage.NRGBA {
	srcBounds := img.Bounds()
	if srcBounds.Min.X == 0 && srcBounds.Min.Y == 0 {
		if src0, ok := img.(*goimage.NRGBA); ok {
			return src0
		}
	}
	return imaging.Clone(img)
}


//--------------------------------------------------
//My own functions
//--------------------------------------------------
/*

func ProcessImage (img *goimage.NRGBA, plate *wtype.LHPlate) (resizedImg goimage.NRGBA, imgPalette color.Palette){

}

func FindClosestColor (col *color.Color, pal *color.Palette)(closeCol color.Color){

	closeCol = pal.Convert(col)
	return closeCol

}
*/