package image

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	goimage "image"
	"image/color"
	"image/color/palette"
	"image/gif"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strconv"

	anthapath "github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/AnthaPath"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wutil"
	"github.com/antha-lang/antha/inventory"
	"github.com/disintegration/imaging"
)

var availableColours = []color.Color{
	color.RGBA{R: uint8(242), G: uint8(243), B: uint8(242), A: uint8(255)}, // white
	color.RGBA{R: uint8(13), G: uint8(105), B: uint8(171), A: uint8(255)},  // blue
	color.RGBA{R: uint8(245), G: uint8(205), B: uint8(47), A: uint8(255)},  // yellow
	color.RGBA{R: uint8(75), G: uint8(151), B: uint8(74), A: uint8(255)},   // green
	color.RGBA{R: uint8(196), G: uint8(40), B: uint8(27), A: uint8(255)},   // red
	color.RGBA{R: uint8(196), G: uint8(40), B: uint8(27), A: uint8(255)},   // black
}

// Colourcomponentmap is a map of RGB colour to description for use as key in
// crossreferencing colour to component in other maps. Very badly named.
var Colourcomponentmap = map[color.Color]string{
	color.RGBA{R: uint8(242), G: uint8(243), B: uint8(242), A: uint8(255)}: "white",
	color.RGBA{R: uint8(13), G: uint8(105), B: uint8(171), A: uint8(255)}:  "blue",
	color.RGBA{R: uint8(245), G: uint8(205), B: uint8(47), A: uint8(255)}:  "yellow",
	color.RGBA{R: uint8(75), G: uint8(151), B: uint8(74), A: uint8(255)}:   "green",
	color.RGBA{R: uint8(196), G: uint8(40), B: uint8(27), A: uint8(255)}:   "red",
	color.RGBA{R: uint8(0), G: uint8(0), B: uint8(0), A: uint8(255)}:       "black",
	color.RGBA{R: uint8(0), G: uint8(0), B: uint8(0), A: uint8(0)}:         "transparent",
}

var neonComponentMap = map[color.Color]string{
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
}

// AllResampleFilters are the resample filters from disintegration package
var AllResampleFilters = map[string]imaging.ResampleFilter{
	"Cosine":            imaging.Cosine,
	"Welch":             imaging.Welch,
	"Blackman":          imaging.Blackman,
	"Hamming":           imaging.Hamming,
	"Hann":              imaging.Hann,
	"Lanczos":           imaging.Lanczos,
	"Bartlett":          imaging.Bartlett,
	"Guassian":          imaging.Gaussian,
	"BSpline":           imaging.BSpline,
	"CatmullRom":        imaging.CatmullRom,
	"MitchellNetravali": imaging.MitchellNetravali,
	"Hermite":           imaging.Hermite,
	"Linear":            imaging.Linear,
	"Box":               imaging.Box,
	"NearestNeighbour":  imaging.NearestNeighbor,
}

var proteinPaintboxMap = map[color.Color]string{
	// under visible light

	// Chromogenic proteins
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

var uvProteinPaintboxMap = map[color.Color]string{

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

var proteinPaintboxSubsetMap = map[color.Color]string{
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

// OpenFile takes a wtype.File object and return the contents as image.NRGBA
func OpenFile(file wtype.File) (*goimage.NRGBA, error) {
	data, err := file.ReadAll()
	if err != nil {
		return nil, err
	}

	buf := bytes.NewReader(data)

	img, err := imaging.Decode(buf)
	if err != nil {
		return nil, err
	}

	return imaging.Clone(img), nil
}

// OpenGIF takes in a wtype.File object and decode the bytes to return a GIF object
func OpenGIF(file wtype.File) (*gif.GIF, error) {

	//returning bytes
	data, err := file.ReadAll()
	if err != nil {
		return nil, err
	}

	//converting bytes to io.reader type
	reader := bytes.NewReader(data)

	img, err := gif.DecodeAll(reader)
	if err != nil {
		return nil, err
	}

	return img, nil
}

// Export an image to file. The image format is derived from filename extension
func Export(img *goimage.NRGBA, fileName string) (file wtype.File, err error) {
	var imageFormat imaging.Format

	if filepath.Ext(fileName) == "" {
		imageFormat = imaging.PNG
		fileName = fileName + "." + "png"
	} else if filepath.Ext(fileName) == ".png" {
		imageFormat = imaging.PNG
	} else if filepath.Ext(fileName) == ".jpg" || filepath.Ext(fileName) == ".jpeg" {
		imageFormat = imaging.JPEG
	} else if filepath.Ext(fileName) == ".tif" || filepath.Ext(fileName) == ".tiff" {
		imageFormat = imaging.TIFF
	} else if filepath.Ext(fileName) == ".gif" {
		imageFormat = imaging.GIF
	} else if filepath.Ext(fileName) == ".BMP" {
		imageFormat = imaging.BMP
	} else {
		return file, fmt.Errorf("unsupported image file format: %s", filepath.Ext(fileName))
	}

	var buf bytes.Buffer

	err = imaging.Encode(&buf, img, imageFormat)
	if err != nil {
		return
	}

	err = file.WriteAll(buf.Bytes())
	if err != nil {
		return
	}

	file.Name = fileName

	return
}

// AvailablePalettes returns all the (hardcoded) palettes in the library. The
// keys are the pallette names.
func AvailablePalettes() map[string]color.Palette {

	var avail = make(map[string]color.Palette)

	avail["Palette1"] = paletteFromMap(Colourcomponentmap) //Chosencolourpalette,
	avail["Neon"] = paletteFromMap(neonComponentMap)
	avail["WebSafe"] = palette.WebSafe //websafe,
	avail["Plan9"] = palette.Plan9
	avail["ProteinPaintboxVisible"] = paletteFromMap(proteinPaintboxMap)
	avail["ProteinPaintboxUV"] = paletteFromMap(uvProteinPaintboxMap)
	avail["ProteinPaintboxSubset"] = paletteFromMap(proteinPaintboxSubsetMap)
	avail["Gray"] = MakeGreyScalePalette()
	avail["None"] = color.Palette{}

	if _, err := os.Stat(filepath.Join(anthapath.Path(), "testcolours.json")); err == nil {
		invmap, err := makeLatestColourMap(filepath.Join(anthapath.Path(), "testcolours.json"))
		if err != nil {
			panic(err.Error())
		}
		avail["inventory"] = paletteFromMap(invmap)
	}

	if _, err := os.Stat(filepath.Join(anthapath.Path(), "UVtestcolours.json")); err == nil {
		uvinvmap, err := makeLatestColourMap(filepath.Join(anthapath.Path(), "UVtestcolours.json"))
		if err != nil {
			panic(err.Error())
		}
		avail["UVinventory"] = paletteFromMap(uvinvmap)
	}
	return avail
}

// AvailableComponentmaps returns available palettes in a map. Bad name.
func AvailableComponentmaps() map[string]map[color.Color]string {
	var componentMaps = make(map[string]map[color.Color]string)
	componentMaps["Palette1"] = Colourcomponentmap
	componentMaps["Neon"] = neonComponentMap
	componentMaps["ProteinPaintboxVisible"] = proteinPaintboxMap
	componentMaps["ProteinPaintboxUV"] = uvProteinPaintboxMap
	componentMaps["ProteinPaintboxSubset"] = proteinPaintboxSubsetMap

	if _, err := os.Stat(filepath.Join(anthapath.Path(), "testcolours.json")); err == nil {
		invmap, err := makeLatestColourMap(filepath.Join(anthapath.Path(), "testcolours.json"))
		if err != nil {
			panic(err.Error())
		}

		componentMaps["inventory"] = invmap
	}
	if _, err := os.Stat(filepath.Join(anthapath.Path(), "UVtestcolours.json")); err == nil {
		uvinvmap, err := makeLatestColourMap(filepath.Join(anthapath.Path(), "UVtestcolours.json"))
		if err != nil {
			panic(err.Error())
		}

		componentMaps["UVinventory"] = uvinvmap
	}
	return componentMaps
}

// visibleEquivalentMaps returns just the proteinPaintboxUV color map
func visibleEquivalentMaps() map[string]map[color.Color]string {
	componentMaps := make(map[string]map[color.Color]string)
	componentMaps["ProteinPaintboxUV"] = proteinPaintboxMap

	if _, err := os.Stat(filepath.Join(anthapath.Path(), "testcolours.json")); err == nil {
		m, err := makeLatestColourMap(filepath.Join(anthapath.Path(), "testcolours.json"))
		if err != nil {
			panic(err.Error())
		}
		componentMaps["UVinventory"] = m
	}
	return componentMaps
}

// paletteFromMap returns a palette of all colors in it given a color map.
func paletteFromMap(colourmap map[color.Color]string) (palette color.Palette) {
	for key := range colourmap {
		palette = append(palette, key)
	}

	return

}

// makeLatestColourMap will take a json file and return map with color objects as keys to their ID.
func makeLatestColourMap(jsonmapfilename string) (colourToName map[color.Color]string, err error) {
	var nameToRGB map[string]color.NRGBA

	data, err := ioutil.ReadFile(jsonmapfilename)

	if err != nil {
		return colourToName, err
	}

	err = json.Unmarshal(data, &nameToRGB)
	if err != nil {
		return colourToName, err
	}

	nameToColor := make(map[string]color.Color)
	for key, value := range nameToRGB {
		nameToColor[key] = value
	}

	colourToName, err = invertNameToPaletteMap(nameToColor)

	return colourToName, err
}

// MakeGreyScalePalette returns a palette of grey shades
func MakeGreyScalePalette() []color.Color {
	var greyPalette []color.Color
	for i := 0; i < 256; i++ {
		greyPalette = append(greyPalette, color.Gray{Y: uint8(i)})
	}

	return greyPalette
}

// invertPaletteToNameMap reverses the keys and values in a colormap. The color
// names become keys for the color objects
func invertPaletteToNameMap(colourMap map[color.Color]string) (map[string]color.Color, error) {
	var nameMap = make(map[string]color.Color)

	for key, value := range colourMap {
		old, seen := nameMap[value]
		if seen {
			return nil, fmt.Errorf("error inserting key %s value %s: value %s already present", key, value, old)
		}
		nameMap[value] = key
	}
	return nameMap, nil
}

// invertNameToPaletteMap reverses the keys and values in a colormap. The
// color objects become keys for the color names
func invertNameToPaletteMap(nameMap map[string]color.Color) (map[color.Color]string, error) {
	var colourMap = make(map[color.Color]string)

	for key, value := range nameMap {
		old, seen := colourMap[value]
		if seen {
			return nil, fmt.Errorf("error inserting key %s value %s: value %s already present", key, value, old)
		}

		colourMap[value] = key
	}
	return colourMap, nil
}

// MakeSubMapfromMap extracts colors from a color library Given colour names.
// Bad name.
func MakeSubMapfromMap(colourMap map[color.Color]string, colourNames []string) map[color.Color]string {
	var newMap = make(map[color.Color]string)

	nameMap, err := invertPaletteToNameMap(colourMap)

	if err != nil {
		panic(err)
	}

	for _, name := range colourNames {
		colour := nameMap[name]

		if colour != nil {
			newMap[colour] = name
		}
	}

	return newMap
}

// MakeSubPallette extracts colors from a palette given a slice of colour
// names. Bad name.
func MakeSubPallette(paletteName string, colourNames []string) color.Palette {
	paletteMap := AvailableComponentmaps()[paletteName]

	subMap := MakeSubMapfromMap(paletteMap, colourNames)

	return paletteFromMap(subMap)
}

// RemoveDuplicatesKeysfromMap will loop over a map of colors to find and
// delete duplicates. Entries with duplicate keys are deleted.
func RemoveDuplicatesKeysfromMap(elements map[string]color.Color) map[string]color.Color {
	var encountered map[string]bool
	result := make(map[string]color.Color)

	for key, v := range elements {

		if encountered[key] {
			continue
		}

		encountered[key] = true
		result[key] = v
	}

	return result
}

// RemoveDuplicatesValuesfromMap will loop over a map of colors to find and
// delete duplicates. Entries with duplicate values are deleted.
func RemoveDuplicatesValuesfromMap(elements map[string]color.Color) map[string]color.Color {
	var encountered map[color.Color]bool
	result := make(map[string]color.Color)

	for key, v := range elements {

		if encountered[v] {
			continue
		}

		encountered[v] = true
		result[key] = v
	}
	return result
}

//----------------------------------------------------------------------------
//Image manipulations
//----------------------------------------------------------------------------

// paletteFromColorarray makes a palette from a slice of colors.
func paletteFromColorarray(colors []color.Color) color.Palette {
	var newPalette color.Palette = colors
	return newPalette
}

// ColourtoCMYK will convert an object with the color interface to color.CYMK
func ColourtoCMYK(colour color.Color) (cmyk color.CMYK) {
	// fmt.Println("colour", colour)
	r, g, b, _ := colour.RGBA()
	cmyk.C, cmyk.M, cmyk.Y, cmyk.K = color.RGBToCMYK(uint8(r), uint8(g), uint8(b))
	return
}

// ColourtoGrayscale will convert an object with the color interface to
// color.Gray
func ColourtoGrayscale(colour color.Color) color.Gray {
	r, g, b, _ := colour.RGBA()
	return color.Gray{
		Y: uint8((0.2126 * float64(r)) + (0.7152 * float64(g)) + (0.0722 * float64(b))),
	}
}

// Posterize changes an image to use only a small number of different tones.
func Posterize(img *goimage.NRGBA, levels int) (*goimage.NRGBA, error) {
	//We cannot posterize with only one level.
	if levels == 1 {
		return nil, errors.New("Cannot posterize with only one level")
	}

	numberOfAreas := 256 / (levels)
	numberOfValues := 255 / (levels - 1)

	posterized := imaging.Clone(img)

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
				rfloat := (float64(r/256) / float64(numberOfAreas))

				rinttemp, err := wutil.RoundDown(rfloat)
				if err != nil {
					return nil, err
				}
				rnew = float64(rinttemp) * float64(numberOfValues)
			}
			if g == 0 {
				gnew = 0
			} else {
				gfloat := (float64(g/256) / float64(numberOfAreas))

				ginttemp, err := wutil.RoundDown(gfloat)
				if err != nil {
					return nil, err
				}
				gnew = float64(ginttemp) * float64(numberOfValues)
			}
			if b == 0 {
				bnew = 0
			} else {
				bfloat := (float64(b/256) / float64(numberOfAreas))

				binttemp, err := wutil.RoundDown(bfloat)
				if err != nil {
					return nil, err
				}
				bnew = float64(binttemp) * float64(numberOfValues)
			}
			var newcolor color.NRGBA
			newcolor.A = uint8(a)

			rint, err := wutil.RoundDown(rnew)

			if err != nil {
				return nil, err
			}
			newcolor.R = uint8(rint)
			gint, err := wutil.RoundDown(gnew)
			if err != nil {
				return nil, err
			}
			newcolor.G = uint8(gint)
			bint, err := wutil.RoundDown(bnew)
			if err != nil {
				return nil, err
			}
			newcolor.B = uint8(bint)

			posterized.Set(x, y, newcolor)
		}
	}

	return posterized, nil
}

// ResizeImagetoPlate resizes an image to fit the number of wells on a plate.
// We treat wells as pixels. Bad name.
func ResizeImagetoPlate(img *goimage.NRGBA, plate *wtype.LHPlate, algorithm imaging.ResampleFilter, rotate bool) *goimage.NRGBA {

	if img.Bounds().Dy() == plate.WellsY() {
		return img
	}

	if rotate {
		img = imaging.Rotate270(img)
	}
	ratio := float64(img.Bounds().Dx()) / float64(img.Bounds().Dy())
	if ratio <= float64(plate.WellsX())/float64(plate.WellsY()) {
		return imaging.Resize(img, 0, plate.WlsY, algorithm)
	}
	return imaging.Resize(img, plate.WlsX, 0, algorithm)

}

// ResizeImagetoPlateAutoRotate resizes an image to fit the number of wells on
// a plate and if the image is in portrait the image will be rotated by 270
// degrees to optimise resolution on the plate.
func ResizeImagetoPlateAutoRotate(img *goimage.NRGBA, plate *wtype.LHPlate, algorithm imaging.ResampleFilter) *goimage.NRGBA {
	if img.Bounds().Dy() == plate.WellsY() {
		return toNRGBA(img)
	}

	if img.Bounds().Dy() > img.Bounds().Dx() {
		img = imaging.Rotate270(img)
	}
	ratio := float64(img.Bounds().Dx()) / float64(img.Bounds().Dy())
	if ratio <= float64(plate.WellsX())/float64(plate.WellsY()) {
		return imaging.Resize(img, 0, plate.WlsY, algorithm)
	}
	return imaging.Resize(img, plate.WlsX, 0, algorithm)

}

// CheckAllResizealgorithms resizes an image using a variety of different
// algorithms.
func CheckAllResizealgorithms(img *goimage.NRGBA, plate *wtype.LHPlate, rotate bool, algorithms map[string]imaging.ResampleFilter) []*goimage.NRGBA {
	var plateImages []*goimage.NRGBA

	for _, algorithm := range algorithms {
		resized := ResizeImagetoPlate(img, plate, algorithm, rotate)

		plateImages = append(plateImages, resized)

	}

	return plateImages
}

// MakePaletteFromImage will make a color Palette from an image resized to fit
// a given plate type.
func MakePaletteFromImage(img *goimage.NRGBA, plate *wtype.LHPlate, rotate bool) color.Palette {
	plateImage := ResizeImagetoPlate(img, plate, imaging.CatmullRom, rotate)

	var colours []color.Color

	// Find out colour at each position:
	for y := 0; y < plateImage.Bounds().Dy(); y++ {
		for x := 0; x < plateImage.Bounds().Dx(); x++ {
			// colour or pixel in RGB
			colour := plateImage.At(x, y)
			colours = append(colours, colour)

		}
	}

	return paletteFromColorarray(colours)
}

// MakeSmallPalleteFromImage will make a color Palette from an image resized to
// fit a given plate type using Plan9 Palette
func MakeSmallPalleteFromImage(img *goimage.NRGBA, plate *wtype.LHPlate, rotate bool) color.Palette {
	plateImage := ResizeImagetoPlate(img, plate, imaging.CatmullRom, rotate)

	// use Plan9 as palette for first round to keep number of colours down to a
	// manageable level
	palette := AvailablePalettes()["Plan9"]

	seen := make(map[color.Color]bool)

	// Find out colour at each position:
	for y := 0; y < plateImage.Bounds().Dy(); y++ {
		for x := 0; x < plateImage.Bounds().Dx(); x++ {
			// colour or pixel in RGB
			colour := plateImage.At(x, y)

			if colour != nil {
				colour = palette.Convert(colour)
				// change colour to colour from a palette
				seen[colour] = true
			}
		}
	}

	var colors []color.Color

	for colour := range seen {
		colors = append(colors, colour)
	}

	return paletteFromColorarray(colors)
}

// ImagetoPlatelayout takes an image, plate type, and palette and returns a map
// of well position to colors. Creates a map of pixel to plate position from
// processing a given image with a chosen colour palette.  It's recommended to
// use at least 384 well plate if autorotate == true, rotate is overridden
func ImagetoPlatelayout(img *goimage.NRGBA, plate *wtype.LHPlate, palette *color.Palette, rotate bool, autoRotate bool) (map[string]color.Color, *goimage.NRGBA) {
	var plateImage *goimage.NRGBA
	if autoRotate {
		plateImage = ResizeImagetoPlateAutoRotate(img, plate, imaging.CatmullRom)
	} else {
		plateImage = ResizeImagetoPlate(img, plate, imaging.CatmullRom, rotate)
	}

	// make map of well position to colour: (array for time being)
	var wellPositions []string
	var colours []color.Color
	positionToColour := make(map[string]color.Color)

	// Find out colour at each position:
	for y := 0; y < plateImage.Bounds().Dy(); y++ {
		for x := 0; x < plateImage.Bounds().Dx(); x++ {
			colour := plateImage.At(x, y)

			if colour == nil {
				continue
			}

			if palette != nil && len(*palette) != 0 {
				// change colour to colour from a palette
				colour = palette.Convert(colour)
				// fmt.Println("x,y,colour", x, y, colour)
				plateImage.Set(x, y, colour)
			}

			// equivalent well position
			pos := wutil.NumToAlpha(y+1) + strconv.Itoa(x+1)
			wellPositions = append(wellPositions, pos)
			positionToColour[pos] = colour

			colours = append(colours, colour)
		}
	}

	return positionToColour, plateImage
}

// PrintFPImagePreview takes an image, plate type, and colors from
// visiblemap/uvmap and (use to) save the resulting processed image.
func PrintFPImagePreview(img *goimage.NRGBA, plate *wtype.LHPlate, rotate bool, visiblemap, uvmap map[color.Color]string) {

	plateimage := ResizeImagetoPlate(img, plate, imaging.CatmullRom, rotate)

	uvpalette := paletteFromMap(uvmap)

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
			stringkeymap, err := invertPaletteToNameMap(visiblemap)
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

// makeNameComponentMap makes a map linking LHcomponents to colours, and
// assigns them to well
func makeNameToComponentMap(keys []string, components []*wtype.LHComponent) (map[string]*wtype.LHComponent, error) {

	componentMap := make(map[string]*wtype.LHComponent)
	for _, key := range keys {
		for _, component := range components {

			if component.CName == key {
				componentMap[key] = component
				break
			}
		}
	}
	return componentMap, nil
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

// colourLibrary is a ollection of colours
var colourLibrary = map[string]color.Color{

	//ProteinPaintBox under natural light
	"JuniperGFP":     color.RGBA{R: uint8(194), G: uint8(164), B: uint8(72), A: uint8(255)},
	"CindylouCFP":    color.RGBA{R: uint8(224), G: uint8(120), B: uint8(240), A: uint8(254)},
	"YetiYFP":        color.RGBA{R: uint8(251), G: uint8(176), B: uint8(0), A: uint8(255)},
	"CometGFP":       color.RGBA{R: uint8(222), G: uint8(221), B: uint8(68), A: uint8(255)},
	"DasherGFP":      color.RGBA{R: uint8(209), G: uint8(214), B: uint8(0), A: uint8(255)},
	"HollyGFP":       color.RGBA{R: uint8(216), G: uint8(231), B: uint8(15), A: uint8(255)},
	"YukonOFP":       color.RGBA{R: uint8(251), G: uint8(102), B: uint8(79), A: uint8(255)},
	"DonnerMagenta":  color.RGBA{R: uint8(152), G: uint8(76), B: uint8(128), A: uint8(255)},
	"ScroogeOrange":  color.RGBA{R: uint8(215), G: uint8(96), B: uint8(86), A: uint8(255)},
	"SerranoRFP":     color.RGBA{R: uint8(248), G: uint8(84), B: uint8(149), A: uint8(255)},
	"BlazeYFP":       color.RGBA{R: uint8(255), G: uint8(255), B: uint8(146), A: uint8(255)},
	"verywhite":      color.RGBA{R: uint8(242), G: uint8(243), B: uint8(242), A: uint8(255)},
	"CupidPink":      color.RGBA{R: uint8(159), G: uint8(25), B: uint8(103), A: uint8(255)},
	"IvyGFP":         color.RGBA{R: uint8(225), G: uint8(222), B: uint8(120), A: uint8(255)},
	"MaccabeePurple": color.RGBA{R: uint8(82), G: uint8(35), B: uint8(119), A: uint8(255)},
	"MarleyYFP":      color.RGBA{R: uint8(250), G: uint8(210), B: uint8(0), A: uint8(255)},
	"CratchitYFP":    color.RGBA{R: uint8(255), G: uint8(194), B: uint8(0), A: uint8(255)},
	"CayenneRFP":     color.RGBA{R: uint8(248), G: uint8(64), B: uint8(148), A: uint8(255)},
	"E.coli":         color.RGBA{R: uint8(196), G: uint8(183), B: uint8(137), A: uint8(255)},
	"BlitzenBlue":    color.RGBA{R: uint8(70), G: uint8(105), B: uint8(172), A: uint8(255)},
	"TinselPurple":   color.RGBA{R: uint8(77), G: uint8(11), B: uint8(137), A: uint8(255)},
	"KringleYFP":     color.RGBA{R: uint8(231), G: uint8(173), B: uint8(0), A: uint8(255)},
	"GuajilloRFP":    color.RGBA{R: uint8(241), G: uint8(84), B: uint8(152), A: uint8(255)},
	"PaprikaRFP":     color.RGBA{R: uint8(247), G: uint8(132), B: uint8(179), A: uint8(255)},
	"TannenGFP":      color.RGBA{R: uint8(243), G: uint8(138), B: uint8(112), A: uint8(255)},
	"veryblack":      color.RGBA{R: uint8(1), G: uint8(1), B: uint8(1), A: uint8(255)},
	"VixenPurple":    color.RGBA{R: uint8(120), G: uint8(76), B: uint8(190), A: uint8(255)},
	"SeraphinaPink":  color.RGBA{R: uint8(206), G: uint8(89), B: uint8(142), A: uint8(255)},
	"RudolphRFP":     color.RGBA{R: uint8(215), G: uint8(72), B: uint8(76), A: uint8(255)},
	"DreidelTeal":    color.RGBA{R: uint8(27), G: uint8(79), B: uint8(146), A: uint8(255)},
	"LeorOrange":     color.RGBA{R: uint8(228), G: uint8(110), B: uint8(104), A: uint8(255)},
	"FrostyCFP":      color.RGBA{R: uint8(224), G: uint8(120), B: uint8(140), A: uint8(255)},
	"TwinkleCFP":     color.RGBA{R: uint8(196), G: uint8(183), B: uint8(138), A: uint8(255)},
	"VirginiaViolet": color.RGBA{R: uint8(107), G: uint8(80), B: uint8(140), A: uint8(255)},
	"FresnoRFP":      color.RGBA{R: uint8(244), G: uint8(63), B: uint8(150), A: uint8(255)},

	//Generic colors under natural light
	"blue":        color.RGBA{R: uint8(13), G: uint8(105), B: uint8(171), A: uint8(255)},
	"yellow":      color.RGBA{R: uint8(245), G: uint8(205), B: uint8(47), A: uint8(255)},
	"green":       color.RGBA{R: uint8(75), G: uint8(151), B: uint8(74), A: uint8(255)},
	"red":         color.RGBA{R: uint8(196), G: uint8(40), B: uint8(27), A: uint8(255)},
	"black":       color.RGBA{R: uint8(0), G: uint8(0), B: uint8(0), A: uint8(255)},
	"transparent": color.RGBA{R: uint8(0), G: uint8(0), B: uint8(0), A: uint8(0)},
	"white":       color.RGBA{R: uint8(242), G: uint8(243), B: uint8(242), A: uint8(255)},

	//UV Protein PaintBox under UV light
	"UVCupidPink":            color.RGBA{R: uint8(159), G: uint8(25), B: uint8(103), A: uint8(255)},
	"UVScroogeOrange":        color.RGBA{R: uint8(215), G: uint8(96), B: uint8(86), A: uint8(255)},
	"UVFrostyCFP":            color.RGBA{R: uint8(0), G: uint8(255), B: uint8(255), A: uint8(255)},
	"UVMarleyYFP":            color.RGBA{R: uint8(236), G: uint8(255), B: uint8(0), A: uint8(255)},
	"UVCometGFP":             color.RGBA{R: uint8(0), G: uint8(254), B: uint8(0), A: uint8(255)},
	"UVDasherGFP":            color.RGBA{R: uint8(0), G: uint8(255), B: uint8(0), A: uint8(255)},
	"UVSerranoRFP":           color.RGBA{R: uint8(254), G: uint8(23), B: uint8(127), A: uint8(255)},
	"UVVirginiaViolet":       color.RGBA{R: uint8(107), G: uint8(80), B: uint8(140), A: uint8(255)},
	"UVBlazeYFP":             color.RGBA{R: uint8(254), G: uint8(255), B: uint8(83), A: uint8(255)},
	"UVFresnoRFP":            color.RGBA{R: uint8(255), G: uint8(0), B: uint8(166), A: uint8(255)},
	"UVJuniperGFP":           color.RGBA{R: uint8(0), G: uint8(231), B: uint8(162), A: uint8(255)},
	"UVCindylouCFP":          color.RGBA{R: uint8(0), G: uint8(254), B: uint8(255), A: uint8(255)},
	"UVCratchitYFP":          color.RGBA{R: uint8(240), G: uint8(254), B: uint8(0), A: uint8(255)},
	"UVPaprikaRFP":           color.RGBA{R: uint8(252), G: uint8(65), B: uint8(136), A: uint8(255)},
	"UVSeraphinaPink":        color.RGBA{R: uint8(206), G: uint8(89), B: uint8(142), A: uint8(255)},
	"UVMaccabeePurple":       color.RGBA{R: uint8(82), G: uint8(35), B: uint8(119), A: uint8(255)},
	"UVLeorOrange":           color.RGBA{R: uint8(228), G: uint8(110), B: uint8(104), A: uint8(255)},
	"UVHollyGFP":             color.RGBA{R: uint8(0), G: uint8(255), B: uint8(0), A: uint8(254)},
	"UVTannenGFP":            color.RGBA{R: uint8(179), G: uint8(119), B: uint8(57), A: uint8(255)},
	"UVverywhite":            color.RGBA{R: uint8(242), G: uint8(243), B: uint8(242), A: uint8(255)},
	"UVTinselPurple":         color.RGBA{R: uint8(77), G: uint8(11), B: uint8(137), A: uint8(255)},
	"UVRudolphRFP":           color.RGBA{R: uint8(218), G: uint8(92), B: uint8(69), A: uint8(255)},
	"UVE.coli":               color.RGBA{R: uint8(196), G: uint8(183), B: uint8(137), A: uint8(255)},
	"UVBlitzenBlue":          color.RGBA{R: uint8(70), G: uint8(105), B: uint8(172), A: uint8(255)},
	"UVKringleYFP":           color.RGBA{R: uint8(239), G: uint8(255), B: uint8(0), A: uint8(255)},
	"UVGuajilloRFP":          color.RGBA{R: uint8(255), G: uint8(8), B: uint8(138), A: uint8(255)},
	"UVE.coli pUC19 on sgal": color.RGBA{R: uint8(1), G: uint8(1), B: uint8(1), A: uint8(255)},
	"UVYetiYFP":              color.RGBA{R: uint8(253), G: uint8(230), B: uint8(39), A: uint8(255)},
	"UVDonnerMagenta":        color.RGBA{R: uint8(152), G: uint8(76), B: uint8(128), A: uint8(255)},
	"UVIvyGFP":               color.RGBA{R: uint8(0), G: uint8(232), B: uint8(216), A: uint8(255)},
	"UVCayenneRFP":           color.RGBA{R: uint8(255), G: uint8(24), B: uint8(138), A: uint8(255)},
	"UVTwinkleCFP":           color.RGBA{R: uint8(27), G: uint8(79), B: uint8(146), A: uint8(255)},
	"UVYukonOFP":             color.RGBA{R: uint8(254), G: uint8(179), B: uint8(18), A: uint8(255)},
	"UVVixenPurple":          color.RGBA{R: uint8(120), G: uint8(76), B: uint8(190), A: uint8(255)},

	//Generic colors under UV light
	"UVblack":    color.RGBA{R: uint8(0), G: uint8(0), B: uint8(0), A: uint8(255)},
	"UVgreen":    color.RGBA{R: uint8(0), G: uint8(177), B: uint8(94), A: uint8(255)},
	"UVyellow":   color.RGBA{R: uint8(244), G: uint8(231), B: uint8(0), A: uint8(255)},
	"UVorange":   color.RGBA{R: uint8(255), G: uint8(118), B: uint8(0), A: uint8(255)},
	"UVwhite":    color.RGBA{R: uint8(242), G: uint8(243), B: uint8(242), A: uint8(255)},
	"UVgrey":     color.RGBA{R: uint8(149), G: uint8(156), B: uint8(161), A: uint8(255)},
	"UVpurple":   color.RGBA{R: uint8(117), G: uint8(51), B: uint8(127), A: uint8(255)},
	"UVdarkblue": color.RGBA{R: uint8(25), G: uint8(60), B: uint8(152), A: uint8(255)},
	"UVblue":     color.RGBA{R: uint8(0), G: uint8(125), B: uint8(200), A: uint8(255)},
	"UVred":      color.RGBA{R: uint8(255), G: uint8(39), B: uint8(51), A: uint8(255)},
	"UVpink":     color.RGBA{R: uint8(235), G: uint8(41), B: uint8(123), A: uint8(255)},
	"UVCyan":     color.RGBA{R: uint8(0), G: uint8(174), B: uint8(239), A: uint8(255)},
	"UVMagenta":  color.RGBA{R: uint8(236), G: uint8(0), B: uint8(140), A: uint8(255)},
}

// Collection of color IDs
var librarySets = map[string][]string{
	"UV": {"UVCupidPink",
		"UVyellow",
		"UVVixenPurple",
		"UVIvyGFP",
		"UVDasherGFP",
		"UVPaprikaRFP",
		"UVLeorOrange",
		"UVTannenGFP",
		"UVDonnerMagenta",
		"UVTwinkleCFP",
		"UVSerranoRFP",
		"UVE.coli pUC19 on sgal",
		"UVorange",
		"UVCyan",
		"UVMagenta",
		"UVTinselPurple",
		"UVpink",
		"UVVirginiaViolet",
		"UVCayenneRFP",
		"UVgreen",
		"UVFresnoRFP",
		"UVKringleYFP",
		"UVdarkblue",
		"UVMarleyYFP",
		"UVCometGFP",
		"UVSeraphinaPink",
		"UVGuajilloRFP",
		"UVBlazeYFP",
		"UVblack",
		"UVblue",
		"UVJuniperGFP",
		"UVCindylouCFP",
		"UVRudolphRFP",
		"UVBlitzenBlue",
		"UVpurple",
		"UVred",
		"UVCratchitYFP",
		"UVgrey",
		"UVverywhite",
		"UVHollyGFP",
		"UVwhite",
		"UVMaccabeePurple",
		"UVYukonOFP",
		"UVE.coli",
		"UVYetiYFP",
		"UVScroogeOrange",
		"UVFrostyCFP",
	},
	"VisibleLight": {"JuniperGFP",
		"YukonOFP",
		"CayenneRFP",
		"E.coli",
		"PaprikaRFP",
		"IvyGFP",
		"HollyGFP",
		"verywhite",
		"DreidelTeal",
		"LeorOrange",
		"DasherGFP",
		"BlazeYFP",
		"veryblack",
		"SeraphinaPink",
		"green",
		"VixenPurple",
		"CometGFP",
		"SerranoRFP",
		"MarleyYFP",
		"DonnerMagenta",
		"red",
		"FrostyCFP",
		"transparent",
		"white",
		"YetiYFP",
		"BlitzenBlue",
		"TwinkleCFP",
		"yellow",
		"black",
		"ScroogeOrange",
		"blue",
		"CindylouCFP",
		"GuajilloRFP",
		"TannenGFP",
		"CupidPink",
		"MaccabeePurple",
		"VirginiaViolet",
		"FresnoRFP",
		"CratchitYFP",
		"TinselPurple",
		"KringleYFP",
		"RudolphRFP",
	},
	"Paints": {
		"blue",
		"yellow",
		"green",
		"red",
		"black",
		"transparent",
		"white",
	},
}

type livingColor struct {
	Color *color.NRGBA
	Seq   string
}

// Living Colors, we use a minimalist object constructor with set defaults to facilitate editing this library
var livingColors = map[string]livingColor{
	"DasherGFP": livingColor{
		Color: &color.NRGBA{R: uint8(0), G: uint8(255), B: uint8(0), A: uint8(255)},
		Seq:   "ATGACGGCATTGACGGAAGGTGCAAAACTGTTTGAGAAAGAGATCCCGTATATCACCGAACTGGAAGGCGACGTCGAAGGTATGAAATTTATCATTAAAGGCGAGGGTACCGGTGACGCGACCACGGGTACCATTAAAGCGAAATACATCTGCACTACGGGCGACCTGCCGGTCCCGTGGGCAACCCTGGTGAGCACCCTGAGCTACGGTGTTCAGTGTTTCGCCAAGTACCCGAGCCACATCAAGGATTTCTTTAAGAGCGCCATGCCGGAAGGTTATACCCAAGAGCGTACCATCAGCTTCGAAGGCGACGGCGTGTACAAGACGCGTGCTATGGTTACCTACGAACGCGGTTCTATCTACAATCGTGTCACGCTGACTGGTGAGAACTTTAAGAAAGACGGTCACATTCTGCGTAAGAACGTTGCATTCCAATGCCCGCCAAGCATTCTGTATATTCTGCCTGACACCGTTAACAATGGCATCCGCGTTGAGTTCAACCAGGCGTACGATATTGAAGGTGTGACCGAAAAACTGGTTACCAAATGCAGCCAAATGAATCGTCCGTTGGCGGGCTCCGCGGCAGTGCATATCCCGCGTTATCATCACATTACCTACCACACCAAACTGAGCAAAGACCGCGACGAGCGCCGTGATCACATGTGTCTGGTAGAGGTCGTGAAAGCGGTTGATCTGGACACGTATCAGTAATGAGAATTCTGTACACTCGAG",
	},
	"RudolphRFP": livingColor{
		Color: &color.NRGBA{R: uint8(218), G: uint8(92), B: uint8(69), A: uint8(255)},
		Seq:   "atgtccctgtcgaaacaagtactgccacacgatgttaagatgcgctatcatatggatggctgcgttaatggtcattctttcaccattgagggtgaaggtgcaggcaaaccgtatgagggcaagaagatcttggaactgcgcgtgacgaaaggtggcccgctgccttttgcgttcgatatcctgagcagcgtttttacctacggtaaccgttgtttttgcgagtatccagaggacatgccggactactttaaacagagcctgccggaaggtcattcttgggaacgcaccctgatgtttgaggatggcggttgtggtacggcgagcgcgcacatttccctggacaagaactgcttcgtgcacaagagcaccttccacggcgtcaatttcccggcaaacggtccggtcatgcaaaagaaagctatgaactgggagccgagcagcgaactgattacggcgtgcgacggtatcctgaaaggcgatgtgaccatgtttctgttgctggaaggtggccaccgtcttaaatgtcagttcaccaccagctacaaagcccacaaggcagttaagatgccgccgaatcacattatcgaacacgtgcttgttaaaaaagaggttgccgacggctttcagatccaagagcatgcggtcgcaaagcacttcaccgtcgacgttaaagaaacgtaatgagaattctgtacactcgag",
	},
	"FresnoRFP": livingColor{
		Color: &color.NRGBA{R: uint8(255), G: uint8(0), B: uint8(166), A: uint8(255)},
		Seq:   "AATAGCCTGATTAAAGAGAATATGCACATGAAGCTGTACATGGAAGGCACGGTGAATAACCACCACTTCAAATGCACCAGCGAGGGTGAGGGTAAACCGTATGAAGGCACCCAAACGATGCGTATCAAAGTTGTTGAGGGTGGCCCGTTGCCGTTTGCGTTCGACATTTTAGCGACGAGCTTTATGTATGGCTCTCGTACGTTTATCAAGTACCCGAAGGGTATTCCGGACTTTTTCAAACAATCTTTTCCAGAGGGTTTCACCTGGGAGCGCGTGACTCGCTACGAAGATGGCGGCGTCGTGACCGCAACGCAGGATACCTCCCTGGAAGATGGCTGCCTGGTCTACCACGTTCAGGTCCGTGGTGTCAATTTCCCGAGCAATGGTCCGGTTATGCAGAAGAAAACCCTGGGTTGGGAACCGAACACCGAGATGTTGTATCCTGCAGATGGTGGCCTGGAAGGTCGCAGCGACATGGCATTGAAACTGGTCGGTGGCGGCCATCTGAGCTGTAGCTTCGTGACCACGTATCGTTCGAAGAAAACGGTCGGTAACATCAAAATGCCGGGTATTCACGCGGTTGACCACCGTCTGGTGCGCATTAAAGAAGCCGACAAAGAGACTTACGTGGAGCAACATGAAGTAGCCGTTGCGAAATTTGCTGGTTTGGGCGGTGGTATGGACGAACTGTACAAATGAGAATTCTGTACACTCGAG",
	},
	"MarleyYFP": livingColor{
		Color: &color.NRGBA{R: uint8(236), G: uint8(255), B: uint8(0), A: uint8(255)},
		Seq:   "ACGGCATTGACGGAAGGTGCAAAACTGTTTGAGAAAGAGATCCCGTATATCACCGAACTGGAAGGCGACGTCGAAGGTATGAAATTTATCATTAAAGGCGAGGGTACCGGTGACGCGAGCGTTGGTAAGGTCGACGCGCAGTTTATCTGCACCACGGGTGACGTCCCGGTGCCGTGGAGCACGTTGGTGACGACGCTGACTTACGGTGCCCAATGTTTCGCGAAATATCCGCGTCACATCGCAGACTTCTTCAAATCCTGTATGCCGGAAGGTTATGTGCAAGAGCGCACTATTACCTTTGAGGGTGACGGTGTCTTTAAGACCCGTGCCGAAGTGACCTTCGAAAACGGTAGCGTTTACAACCGTGTCAAGCTGAATGGCCAGGGTTTCAAAAAGGATGGTCACGTTCTGGGTAAGAATCTGGAATTCAACTTCACCCCACACTGCCTGTACATCTGGGGCGATCAAGCGAATCATGGTCTGAAAAGCGCATTTAAGGTGATGCACGAGATTACCGGCTCCAAAGAAGATTTCATCGTGGCTGATCACACCCAGATGAATACCCCGATTGGCGGTGGCCCTGTGCACGTTCCGGAATATCATCACATTACCTACCATGTCACGCTGAGCAAGGATGTTACCGACCATCGCGACCACCTGAACATCGTAGAGGTTATCAAAGCAGTGGACCTGGAGACGTATCGCTAATGAGAATTCTGTACACTCGAG",
	},
	"TwinkleCFP": livingColor{
		Color: &color.NRGBA{R: uint8(27), G: uint8(79), B: uint8(146), A: uint8(255)},
		Seq:   "TCGTCTGGTGCCAAATTGTTTGAAAAGAAGATCCCGTATATCACTGAACTCGAGGGCGACGTCAATGGTATGAAGTTTACCATTCATGGTAAAGGTACCGGCGATGCGACCACGGGTAAAATTAAAGCGCAGTTCATCTGCACTACGGGCGACGTTCCGGTCCCGTGGGCAACCCTGGTGAGCACCCTGAGCTACGGTGTTCAGTGTTTCGCCAAGTACCCGAGCGAGCTGAAGGATTTCTTTAAGAGCTGCATGCCGGAAGGTTATGTTCAAGAGCGTACCATCACCTTCGAAGGCGACGGCGTGTTTAGGACGCGTGCTGAGGTTACCTTTGAAAACGGTTCTGTCTACAATCGTGTCACTCTGACTGGCGAGGGTTTTAAGAAAGACGGTCACATTCTGCGTAAGAACGTTGCATTCCAATGCCCGCCAGATATTGTGTATATTCTGCCTGACACCGTTAACAATGGCATCCGCGTTGAGTTCAACCAGGCGTACGATATTGAAGGTGTGACCGAAAAACTGGTTACCAAATGCAGCCAAATGAATCGTCCGTTGGCGGGCTCCGCGGCAGTGCATATCCCGCGTTATCATCACCTGAGCAAACACACCAAACTGAGCAAAGACCGCGACGAGCGCCGTGATCACATGTGTCTGGTAGAGGTCGTGAAAGCGGTTGATCTGGACACGTATCAGTAATGAGAATTCTGTACACTCGAG",
	},
	"mTagBFP": livingColor{
		Color: &color.NRGBA{},
		Seq:   "GTGTCTAAGGGCGAGGAGCTGATTAAGGAGAACATGCACATGAAGCTGTACATGGAGGGCACCGTGGACAACCATCACTTCAAGTGCACATCCGAGGGCGAAGGCAAGCCCTACGAGGGCACCCAGACCATGAGAATCAAGGTGGTCGAGGGCGGCCCTCTCCCCTTCGCCTTCGACATCCTGGCTACTAGCTTCCTCTACGGCAGCAAGACCTTCATCAACCACACCCAGGGCATCCCCGACTTCTTCAAGCAGTCCTTCCCTGAGGGCTTCACATGGGAGAGAGTCACCACATACGAAGACGGGGGCGTGCTGACCGCTACCCAGGACACCAGCCTCCAGGACGGCTGCCTCATCTACAACGTCAAGATCAGAGGGGTGAACTTCACATCCAACGGCCCTGTGATGCAGAAGAAAACACTCGGCTGGGAGGCCTTCACCGAGACGCTGTACCCCGCTGACGGCGGCCTGGAAGGCAGAAACGACATGGCCCTGAAGCTCGTGGGCGGGAGCCATCTGATCGCAAACGCCAAGACCACATATAGATCCAAGAAACCCGCTAAGAACCTCAAGATGCCTGGCGTCTACTATGTGGACTACAGACTGGAAAGAATCAAGGAGGCCAACAACGAAACCTACGTCGAGCAGCACGAGGTGGCAGTGGCCAGATACTGCGACCTCCCTAGCAAACTGGGGCACAAGCTTAATTAATGAGAATTCTGTACACTCGAG",
	},
	//since this should be the background this is set to pure white.
	"Nil": livingColor{
		Color: &color.NRGBA{R: uint8(0), G: uint8(0), B: uint8(0), A: uint8(0)},
		Seq:   "CCAAATGCACCCTTACCACGAAGACAGGATTGTCCGATCCTATATTACGACTTTGGCAGGGGGTTCGCAAGTCCCACCCCAAACGATGCTGAAGGCTCAGGTTACACAGGCACAAGTACTATATATACGAGTTCCCGCTCTTAACCTGGATCGAATGCAGAATCATGCATCGTACCACTGTGTTCGTGTCATCTAGGACGGGCGCAAAGGATATATAATTCAATTAAGAATACCTTATATTATTGTACACCTACCGGTCACCAGCCAACAATGTGCGGATGGCGTTACGACTTACTGGGCCTGATCTCACCGCTTTAGATACCGCACACTGGGCAATACGAGGTAAAGCCAGTCACCCAGTGTCGATCAACAGCTAACGTAACGGTAAGAGGCTCACAAA",
	},
}

// Collection of living color IDs
var livingColorSets = map[string][]string{
	"ProteinPaintBox": {
		"DasherGFP",
		"FresnoRFP",
		"MarleyYFP",
		"TwinkleCFP",
		"mTagBFP",
		"Nil",
	},
}

//---------------------------------------------------
//Types
//---------------------------------------------------

// An AnthaColor represents a color linked to an LHComponent
type AnthaColor struct {
	Color     color.NRGBA
	Component *wtype.LHComponent
}

// An AnthaPalette is an array of anthaColors
type AnthaPalette struct {
	AnthaColors []AnthaColor
	Palette     color.Palette
}

// An AnthaPix is a pixel linked to an anthaColor, and thereby an LHComponent
type AnthaPix struct {
	Color    AnthaColor
	Location wtype.WellCoords
}

// An AnthaImg represents an image where pixels are linked
type AnthaImg struct {
	Plate   *wtype.LHPlate
	Pix     []AnthaPix
	Palette AnthaPalette
}

// A LivingColor links a color to physical information such as DNASequence and LHcomponent
type LivingColor struct {
	ID        string
	Color     color.NRGBA
	Seq       wtype.DNASequence
	Component *wtype.LHComponent
}

// A LivingPalette holds an array of livingColor
type LivingPalette struct {
	LivingColors []LivingColor
}

// A LivingPix represents a pixel alongs
type LivingPix struct {
	Color    LivingColor
	Location wtype.WellCoords
}

// A LivingImg is a representation of an image linked to biological data
type LivingImg struct {
	Plate   wtype.LHPlate
	Pix     []LivingPix
	Palette LivingPalette
}

//A LivingGIF is a representation of a GIF linked to biological data
type LivingGIF struct {
	Frames []LivingImg
}

//---------------------------------------------------
//Data Manipulation
//---------------------------------------------------

// SelectLibrary will select a hardcoded set of colors and return them as a palette object
func SelectLibrary(libID string) (palette color.Palette) {

	selectedLib, found := librarySets[libID]
	//Checking for an empty return value
	if !found {
		panic(fmt.Sprintf("library %s not found so could not make palette", libID))
	}

	for _, colID := range selectedLib {
		palette = append(palette, colourLibrary[colID])
	}

	return
}

// SelectColor will return the desired color from a string ID
func SelectColor(colID string) (selectedColor color.Color) {

	selectedColor, found := colourLibrary[colID]
	//Checking for an empty return value
	if !found {
		panic(fmt.Sprintf("library %s not found so could not make palette", colID))
	}

	return
}

// SelectLivingColorLibrary will return a the desired set of livingcolors given an ID
func SelectLivingColorLibrary(ctx context.Context, libID string) (palette LivingPalette) {

	selectedLib, found := livingColorSets[libID]

	//Checking for an empty return value
	if !found {
		panic(fmt.Sprintf("library %s not found so could not make palette", libID))
	}

	for _, colorID := range selectedLib {
		c := livingColors[colorID]
		lc := MakeLivingColor(ctx, colorID, c.Color, c.Seq)

		palette.LivingColors = append(palette.LivingColors, *lc)
	}

	return
}

// SelectLivingColor will return a LivingColor given its ID
func SelectLivingColor(ctx context.Context, colID string) LivingColor {

	c, found := livingColors[colID]

	//Checking for an empty return value
	if !found {
		panic(fmt.Sprintf("library %s not found so could not make palette", colID))
	}

	lc := MakeLivingColor(ctx, colID, c.Color, c.Seq)

	return *lc
}

// Convert returns the AnthaPalette AnthaColor closest to c in Euclidean R,G,B space.
func (p AnthaPalette) Convert(c color.Color) AnthaColor {

	//getting colors of the current anthacolors in the anthaPalette
	anthaColors := p.AnthaColors

	//Checking if there are no colors in the given palette
	if len(anthaColors) == 0 {
		fmt.Println(errors.New("No color found in the given palette"))
	}

	return anthaColors[p.Index(c)]
}

// Index finds the closest color in an anthapalette and returns the index for the anthacolor
func (p AnthaPalette) Index(c color.Color) int {

	cr, cg, cb, ca := c.RGBA()
	ret, bestSum := 0, uint32(1<<32-1)
	for i := range p.AnthaColors {

		//getting color of the current anthacolor in the anthaPalette
		extractedColor := p.AnthaColors[i].Color

		vr, vg, vb, va := extractedColor.RGBA()
		sum := sqDiff(cr, vr) + sqDiff(cg, vg) + sqDiff(cb, vb) + sqDiff(ca, va)
		if sum < bestSum {
			if sum == 0 {
				return i
			}
			ret, bestSum = i, sum
		}
	}
	return ret
}

// CheckPresence returns a bool indicating if the AnthaColor is already in the AnthaPalette
func (p AnthaPalette) CheckPresence(c AnthaColor) bool {

	for _, pc := range p.AnthaColors {
		if reflect.DeepEqual(pc, c) {
			return true
		}
	}

	return false
}

// Convert returns the LivingPalette LivingColor closest to c in Euclidean R,G,B space.
func (p LivingPalette) Convert(c color.Color) LivingColor {

	//getting colors of the LivingColors in the LivingPalette
	livingColors := p.LivingColors

	//Checking if there are no colors in the given palette
	if len(livingColors) == 0 {
		fmt.Println(errors.New("No color found in the given palette"))
	}

	return livingColors[p.Index(c)]
}

// Index finds the closest color in a LivingPalette and returns the index for the LivingPalette
func (p LivingPalette) Index(c color.Color) int {

	cr, cg, cb, ca := c.RGBA()
	ret, bestSum := 0, uint32(1<<32-1)
	for i := range p.LivingColors {

		//getting color of the current anthacolor in the anthaPalette
		extractedColor := p.LivingColors[i].Color

		vr, vg, vb, va := extractedColor.RGBA()
		sum := sqDiff(cr, vr) + sqDiff(cg, vg) + sqDiff(cb, vb) + sqDiff(ca, va)
		if sum < bestSum {
			if sum == 0 {
				return i
			}
			ret, bestSum = i, sum
		}
	}
	return ret
}

// an internal goimage function
// sqDiff returns the squared-difference of x and y, shifted by 2 so that
// adding four of those won't overflow a uint32.
//
// x and y are both assumed to be in the range [0, 0xffff].
func sqDiff(x, y uint32) uint32 {
	var d uint32
	if x > y {
		d = x - y
	} else {
		d = y - x
	}
	return (d * d) >> 2
}

// CheckPresence returns a bool indicating if the LivingColor is already in the LivingPalette
func (p LivingPalette) CheckPresence(c LivingColor) bool {

	for _, pc := range p.LivingColors {
		if reflect.DeepEqual(pc, c) {
			return true
		}
	}

	return false
}

// Compare returns a bool indicating if the LivingPixel is the same as the given one
func (p1 LivingPix) Compare(p2 LivingPix) (same bool) {

	if reflect.DeepEqual(p1, p2) {
		return true
	}
	return false

}

// GetStates returns an array of unique state changes in a LivingGIF
func (g1 LivingGIF) GetStates() [][]string {

	//error checking
	if len(g1.Frames) > 2 {
		panic("cannot get states of GIF with more than 1 frame")
	}

	if len(g1.Frames[0].Pix) != len(g1.Frames[1].Pix) {
		panic("different no of pixels in frames")
	}

	//globals
	stateList := make(map[int][]string)
	stateNum := 0

	//finding each state change and put them into the global map
	for _, frame := range g1.Frames {

		for _, pix := range frame.Pix {

			stateList[stateNum] = append(stateList[stateNum], pix.Color.ID)
			stateNum++
		}

		stateNum = 0
	}

	//remove duplicates
	//concatenate state lists to use computationally efficient maps (this is necessary, using arrays takes way too long)
	var stateListNm []string

	for _, state := range stateList {

		nm := state[0] + "-" + state[1]
		stateListNm = append(stateListNm, nm)

	}

	//removing duplicate values
	seen := make(map[string]struct{}, len(stateListNm))
	j := 0
	for _, stateNm := range stateListNm {
		if _, ok := seen[stateNm]; ok {
			continue
		}
		seen[stateNm] = struct{}{}
		stateListNm[j] = stateNm
		j++
	}
	stateListNm = stateListNm[:j]

	//parsing resulting strings to arrays to make more sense in returned data
	uniqueStates := [][]string{}

	var re = regexp.MustCompile(`\w+`)

	for _, state := range stateListNm {

		match := re.FindAllString(state, -1)
		uniqueStates = append(uniqueStates, match)
	}

	return uniqueStates
}

// MakeAnthaPalette will make a palette of Colors linked to LHcomponents. They
// are merged according to their order in the slice
func MakeAnthaPalette(palette color.Palette, components []*wtype.LHComponent) *AnthaPalette {

	//global placeholders
	var anthaColor AnthaColor
	var anthaPalette AnthaPalette

	//checking that there are enough LHComponents to make the Palette
	if len(palette) != len(components) {
		panic(fmt.Sprintf("Not enough LHComponents to make Palette"))
	} else {

		for i := range palette {

			//converting to NRGBA and passing to the AnthaColor object
			r, g, b, a := palette[i].RGBA()
			var NRGBA = color.NRGBA{uint8(r), uint8(g), uint8(b), uint8(a)}
			anthaColor.Color = NRGBA

			//Passing the LHComponent to the anthaColor
			anthaColor.Component = components[i]

			//appending created color to the AnthaPalette
			anthaPalette.AnthaColors = append(anthaPalette.AnthaColors, anthaColor)
		}
	}

	//appending the palette object to anthapalette (so we can use its coupled functions)
	anthaPalette.Palette = palette

	return &anthaPalette
}

// MakeAnthaImg will create an AnthaImage object from a digital image.
func MakeAnthaImg(goImg *goimage.NRGBA, anthaPalette *AnthaPalette, anthaImgPlate *wtype.LHPlate) (outputImg *AnthaImg, resizedImg *goimage.NRGBA) {

	//Global placeholders
	var anthaPix AnthaPix
	var anthaImgPix []AnthaPix
	var anthaImg AnthaImg

	//Verify that the plate is the same size as the digital image. If not resize.
	if goImg.Bounds().Dy() != anthaImgPlate.WellsY() {
		goImg = ResizeImagetoPlateMin(goImg, anthaImgPlate)
	}

	//Iterate over pixels
	b := goImg.Bounds()
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			//getting rgba values for the image pixel
			r, g, b, a := goImg.At(x, y).RGBA()
			var goPixColor = color.NRGBA{uint8(r), uint8(g), uint8(b), uint8(a)}
			//finding the anthaColor closest to the one given in the palette
			var anthaColor = anthaPalette.Convert(goPixColor)
			anthaPix.Color = anthaColor

			//figuring out the pixel location on the plate
			anthaPix.Location = wtype.WellCoords{X: x, Y: y}

			//appending the pixel to the array that will go in the AnthaImage
			anthaImgPix = append(anthaImgPix, anthaPix)
		}
	}

	//initiating complete image object
	anthaImg.Pix = anthaImgPix
	anthaImg.Palette = *anthaPalette
	anthaImg.Plate = anthaImgPlate

	return &anthaImg, goImg
}

// MakeLivingColor is an object constructor for a LivingColor with default settings
func MakeLivingColor(ctx context.Context, ID string, color *color.NRGBA, seq string) (livingColor *LivingColor) {

	//generating DNA sequence object
	DNASequence := wtype.MakeLinearDNASequence("ColorDNA", seq)

	//use water as LHComponent
	component, err := inventory.NewComponent(ctx, "dna")
	if err != nil {
		panic(err)
	}

	//populate livingColor
	livingColor = &LivingColor{ID, *color, DNASequence, component}

	return livingColor
}

// MakeLivingPalette will make a palette of LivingColors linked to LHcomponents. They are merged according to their order in the slice
func MakeLivingPalette(palette LivingPalette, components []*wtype.LHComponent) *LivingPalette {

	//checking that there are enough LHComponents to make the Palette
	if len(palette.LivingColors) != len(components) {

		panic("Different number of LivingColors an LHComponent used to make a LivingPalette")

	} else {
		//Adding the LHComponents to the livingColors
		for i := range palette.LivingColors {
			palette.LivingColors[i].Component = components[i]
		}
	}

	return &palette
}

// MakeLivingImg will create a LivingImage object from a digital image.
func MakeLivingImg(goImg *goimage.NRGBA, livingPalette *LivingPalette, livingImgPlate *wtype.LHPlate) (outputImg *LivingImg, resizedImg *goimage.NRGBA) {

	//Global placeholders
	var livingPix LivingPix
	var livingImgPix []LivingPix
	var livingImg LivingImg

	//Verify that the plate is the same size as the digital image. If not resize.
	if goImg.Bounds().Dy() != livingImgPlate.WellsY() {
		goImg = ResizeImagetoPlateMin(goImg, livingImgPlate)
	}

	//Iterate over pixels
	b := goImg.Bounds()
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			//getting rgba values for the image pixel
			r, g, b, a := goImg.At(x, y).RGBA()
			var goPixColor = color.NRGBA{uint8(r), uint8(g), uint8(b), uint8(a)}
			//finding the anthaColor closest to the one given in the palette
			var anthaColor = livingPalette.Convert(goPixColor)
			livingPix.Color = anthaColor

			//figuring out the pixel location on the plate
			livingPix.Location = wtype.WellCoords{X: x, Y: y}

			//appending the pixel to the array that will go in the AnthaImage
			livingImgPix = append(livingImgPix, livingPix)
		}
	}

	//initiating complete image object
	livingImg.Pix = livingImgPix
	livingImg.Palette = *livingPalette
	livingImg.Plate = *livingImgPlate

	return &livingImg, goImg
}

// MakeLivingGIF makes a LivingGIF object given a slice of LivingImg
func MakeLivingGIF(imgs []LivingImg) *LivingGIF {

	var livingGIF LivingGIF

	for _, img := range imgs {
		livingGIF.Frames = append(livingGIF.Frames, img)
	}

	return &livingGIF
}

// ResizeImagetoPlateMin is a minimalist resize function. Uses Lanczos
// resampling, which is the best but slowest method.
func ResizeImagetoPlateMin(img *goimage.NRGBA, plate *wtype.LHPlate) *goimage.NRGBA {

	if img.Bounds().Dy() == plate.WellsY() {
		return img
	}

	if img.Bounds().Dy() > img.Bounds().Dx() {
		img = imaging.Rotate270(img)
	}

	ratio := float64(img.Bounds().Dx()) / float64(img.Bounds().Dy())

	if ratio <= float64(plate.WellsX())/float64(plate.WellsY()) {
		return imaging.Resize(img, 0, plate.WlsY, imaging.Lanczos)
	} else {
		return imaging.Resize(img, plate.WlsX, 0, imaging.Lanczos)
	}
}

// ParseGIF extracts frames from a GIF object and return an array of the images
func ParseGIF(GIF *gif.GIF, frameNum []int) (imgs []*goimage.NRGBA, err error) {

	//error check

	//finding maximum number in the given array (no go builtins for this)
	var largest int

	for _, n := range frameNum {

		if n > largest {
			largest = n
		} else {
			continue
		}
	}

	if len(GIF.Image) <= largest {
		panic(fmt.Errorf("Frame number to be extracted is out of bound of the possible frame. The GIF has %d frames (we use array notation to extract, so starts at 0)", len(GIF.Image)))
	}

	//extracting frames
	for _, Num := range frameNum {

		//convert image from image.Paletted to image.NRGBA
		convertedImg := toNRGBA(GIF.Image[Num])

		imgs = append(imgs, convertedImg)
	}

	return imgs, nil
}
