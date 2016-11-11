// pixeltoplate.go
// image.go
package image

import (
	"encoding/json"
	"fmt"
	goimage "image"
	"image/color"
	"image/color/palette"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	anthapath "github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/AnthaPath"
	"github.com/antha-lang/antha/antha/anthalib/mixer"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/antha/anthalib/wutil"
	"github.com/disintegration/imaging"
)

// Colour palette to use // this would relate to a map of components of these available colours in factor

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

var Emptycolourarray color.Palette

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

func ColourtoCMYK(colour color.Color) (cmyk color.CMYK) {
	// fmt.Println("colour", colour)
	r, g, b, _ := colour.RGBA()
	cmyk.C, cmyk.M, cmyk.Y, cmyk.K = color.RGBToCMYK(uint8(r), uint8(g), uint8(b))
	return
}

func ColourtoGrayscale(colour color.Color) (gray color.Gray) {
	r, g, b, _ := colour.RGBA()
	gray.Y = uint8((0.2126 * float64(r)) + (0.7152 * float64(g)) + (0.0722 * float64(b)))
	return
}

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

func MakeGreyScalePalette() (graypalette []color.Color) {

	graypalette = make([]color.Color, 0)
	var shadeofgray color.Gray
	for i := 0; i < 256; i++ {
		shadeofgray = color.Gray{Y: uint8(i)}
		graypalette = append(graypalette, shadeofgray)
	}

	return
}

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

func palettefromMap(colourmap map[color.Color]string) (palette color.Palette) {

	array := make([]color.Color, 0)

	for key, _ := range colourmap {

		array = append(array, key)
	}

	palette = array
	return

}

func paletteFromColorarray(colors []color.Color) (palette color.Palette) {

	var newpalette color.Palette

	newpalette = colors

	palette = newpalette

	//palette = &newpalette
	return
}

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

func MakeSubPallette(palettename string, colournames []string) (subpalette color.Palette) {
	palettemap := AvailableComponentmaps()[palettename]

	submap := MakeSubMapfromMap(palettemap, colournames)

	subpalette = palettefromMap(submap)

	fmt.Println("length of subpalette: ", len(subpalette), subpalette)
	return

}

var ProteinPaintboxmap = map[color.Color]string{
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

func MakeGoimageNRGBA(imagefilename string) (nrgba *goimage.NRGBA) {
	img, err := imaging.Open(imagefilename)
	if err != nil {
		panic(err)
	}

	nrgba = imaging.Clone(img)
	return
}

func Posterize(imagefilename string, levels int) (posterized *goimage.NRGBA, newfilename string) {

	var newcolor color.NRGBA
	numberofAreas := 256 / (levels)
	numberofValues := 255 / (levels - 1)

	img, err := imaging.Open(imagefilename)
	if err != nil {
		panic(err)
	}

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

			// fmt.Println("x,y", x, y, "r,g,b,a", r, g, b, a, "newcolour", newcolor)

			posterized.Set(x, y, newcolor)

		}
	}

	// rename file
	splitfilename := strings.Split(imagefilename, `.`)

	newfilename = filepath.Join(fmt.Sprint(splitfilename[0], "_posterized", `.`, splitfilename[1]))
	// save

	imaging.Save(posterized, newfilename)
	return
}

func ResizeImagetoPlate(imagefilename string, plate *wtype.LHPlate, algorithm imaging.ResampleFilter, rotate bool) (plateimage *goimage.NRGBA) {

	// input files (just 1 in this case)
	files := []string{imagefilename}

	// Colour palette to use // this would relate to a map of components of these available colours in factory
	//availablecolours := chosencolourpalette //palette.WebSafe

	//var plateimages []image.Image

	img, err := imaging.Open(files[0])
	if err != nil {
		panic(err)
	}

	if img.Bounds().Dy() != plate.WellsY() {
		// fmt.Println("hey we're not so different", img.Bounds().Dy(), plate.WellsY())
		// have the option of changing the resize algorithm here

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
		// fmt.Println("i'm the same!!!")
		plateimage = toNRGBA(img)
	}
	return

}

func ResizeImagetoPlateAutoRotate(imagefilename string, plate *wtype.LHPlate, algorithm imaging.ResampleFilter) (plateimage *goimage.NRGBA) {

	// input files (just 1 in this case)
	files := []string{imagefilename}

	// Colour palette to use // this would relate to a map of components of these available colours in factory
	//availablecolours := chosencolourpalette //palette.WebSafe

	//var plateimages []image.Image

	img, err := imaging.Open(files[0])
	if err != nil {
		panic(err)
	}

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

func CheckAllResizealgorithms(imagefilename string, plate *wtype.LHPlate, rotate bool, algorithms map[string]imaging.ResampleFilter) {
	// input files (just 1 in this case)
	files := []string{imagefilename}
	var dir string

	var plateimage *goimage.NRGBA

	// Colour palette to use // this would relate to a map of components of these available colours in factory
	//availablecolours := chosencolourpalette //palette.WebSafe

	//var plateimages []image.Image

	for key, algorithm := range algorithms {

		img, err := imaging.Open(files[0])
		if err != nil {
			panic(err)
		}

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

		// rename file
		splitfilename := strings.Split(imagefilename, `.`)

		dir = splitfilename[0]

		// make dir

		os.MkdirAll(dir, 0777)

		newname := filepath.Join(dir, fmt.Sprint(splitfilename[0], "_", key, "_plateformat", `.`, splitfilename[1]))
		// save
		err = imaging.Save(plateimage, newname)
		if err != nil {
			panic(err)
		}

	}
}

func MakePalleteFromImage(imagefilename string, plate *wtype.LHPlate, rotate bool) (newpallette color.Palette) {

	plateimage := ResizeImagetoPlate(imagefilename, plate, imaging.CatmullRom, rotate)

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

func MakeSmallPalleteFromImage(imagefilename string, plate *wtype.LHPlate, rotate bool) (newpallette color.Palette) {

	plateimage := ResizeImagetoPlate(imagefilename, plate, imaging.CatmullRom, rotate)
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

// create a map of pixel to plate position from processing a given image with a chosen colour palette.
// It's recommended to use at least 384 well plate
// if autorotate == true, rotate is overridden
func ImagetoPlatelayout(imagefilename string, plate *wtype.LHPlate, chosencolourpalette *color.Palette, rotate bool, autorotate bool) (wellpositiontocolourmap map[string]color.Color, numberofpixels int, newname string) {

	var plateimage *goimage.NRGBA

	if autorotate {
		plateimage = ResizeImagetoPlateAutoRotate(imagefilename, plate, imaging.CatmullRom)
	} else {
		plateimage = ResizeImagetoPlate(imagefilename, plate, imaging.CatmullRom, rotate)
	}
	// make map of well position to colour: (array for time being)

	wellpositionarray := make([]string, 0)
	colourarray := make([]color.Color, 0)
	wellpositiontocolourmap = make(map[string]color.Color, 0)
	// need to extend for 1536 plates
	alphabet := []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J",
		"K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X",
		"Y", "Z", "AA", "AB", "AC", "AD", "AE", "AF"}

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
				wellposition := alphabet[y] + strconv.Itoa(x+1)
				fmt.Println(wellposition)
				wellpositionarray = append(wellpositionarray, wellposition)
				wellpositiontocolourmap[wellposition] = colour

				colourarray = append(colourarray, colour)
			}
		}
	}

	// rename file
	splitfilename := strings.Split(imagefilename, `.`)

	newname = splitfilename[0] + "_plateformat" + `.` + splitfilename[1]
	// save
	err := imaging.Save(plateimage, newname)
	if err != nil {
		panic(err)
	}

	// choose colour palette from top

	//arrayfrompalette := make([]color.Color, 0)
	/*
		for i, combo := range colourarray {
			// fmt.Println("for well position ", wellpositionarray[i], ":")
			r, g, b, a := combo.RGBA()

			// fmt.Println("colour (r,g,b,a)= ", r/256, g/256, b/256, a/256)

			// closest palette colour
			// fmt.Println("palette colour:", chosencolourpalette.Convert(combo))
			// fmt.Println("palette colour number:", chosencolourpalette.Index(combo))
		}
	*/
	numberofpixels = len(colourarray)
	// fmt.Println("numberofpixels:", numberofpixels)

	return
}

func PrintFPImagePreview(imagefile string, plate *wtype.LHPlate, rotate bool, visiblemap, uvmap map[color.Color]string) {

	plateimage := ResizeImagetoPlate(imagefile, plate, imaging.CatmullRom, rotate)

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

	// rename file
	splitfilename := strings.Split(imagefile, `.`)

	newname := splitfilename[0] + "_plateformat_UV" + `.` + splitfilename[1]
	// save
	err := imaging.Save(plateimage, newname)
	if err != nil {
		panic(err)
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

	// rename file
	splitfilename = strings.Split(imagefile, `.`)

	newname = splitfilename[0] + "_plateformat_vis" + `.` + splitfilename[1]
	// save
	err = imaging.Save(plateimage, newname)
	if err != nil {
		panic(err)
	}
	return
}

//  Final function for user which uses a given map of position to colour generated from the image processing function  along with lists of available colours, components and plate types
/*func PipetteImagetoPlate(OutPlate *wtype.LHPlate, positiontocolourmap map[string]color.Color, availablecolours []string, componentlist []*wtype.LHComponent, volumeperwell wunit.Volume) (finalsolutions []*wtype.LHComponent) {

	componentmap, err := MakestringtoComponentMap(availablecolours, componentlist)
	if err != nil {
		panic(err)
	}

	solutions := make([]*wtype.LHComponent, 0)

	for locationkey, colour := range positiontocolourmap {

		component := componentmap[Colourcomponentmap[colour]]

		if component.Type == "dna" {
			component.Type = "DoNotMix"
		}

		pixelSample := mixer.Sample(component, volumeperwell)
		solution := mixer.MixTo(OutPlate, locationkey, pixelSample)
		solutions = append(solutions, solution)
	}

	finalsolutions = solutions
	return
}*/

// make a map of which colour description applies to which component assuming all components in factory are available, returns errors if either keys or components cannot be added
/*
func MakestringtoComponentMapFromFactory(colourtostringmap map[color.Color]string) (componentmap map[string]*wtype.LHComponent) {

	componentmap = make(map[string]*wtype.LHComponent, 0)
	//var previouserror error = nil

	for _, colour := range colourtostringmap {

		componentname := colourtostringmap[colour]

		componentmap[componentname] = factory.GetComponentByType(component)

	}
	return
}
*/
// else, specify colours and components to make a map of which colour description applies to which component, returns errors if either keys or components cannot be added
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

//  Final function for blending colours to make and image. Uses a given map of position to colour generated from the image processing function  along with lists of available colours, components and plate types
func PipetteImagebyBlending(OutPlate *wtype.LHPlate, positiontocolourmap map[string]color.Color, cyan *wtype.LHComponent, magenta *wtype.LHComponent, yellow *wtype.LHComponent, black *wtype.LHComponent, volumeperfullcolour wunit.Volume) (finalsolutions []*wtype.LHComponent) {

	solutions := make([]*wtype.LHComponent, 0)

	for locationkey, colour := range positiontocolourmap {

		components := make([]*wtype.LHComponent, 0)

		cmyk := ColourtoCMYK(colour)

		cyanvol := wunit.NewVolume((float64(cmyk.C) * volumeperfullcolour.SIValue()), "l")
		yellowvol := wunit.NewVolume((float64(cmyk.Y) * volumeperfullcolour.SIValue()), "l")
		magentavol := wunit.NewVolume((float64(cmyk.M) * volumeperfullcolour.SIValue()), "l")
		blackvol := wunit.NewVolume((float64(cmyk.K) * volumeperfullcolour.SIValue()), "l")

		cyanSample := mixer.Sample(cyan, cyanvol)
		components = append(components, cyanSample)
		yellowSample := mixer.Sample(yellow, yellowvol)
		components = append(components, yellowSample)
		magentaSample := mixer.Sample(magenta, magentavol)
		components = append(components, magentaSample)
		blackSample := mixer.Sample(black, blackvol)
		components = append(components, blackSample)
		solution := mixer.MixTo(OutPlate.Type, locationkey, 1, components...)
		solutions = append(solutions, solution)
	}

	finalsolutions = solutions
	return
}

//  Final function for blending colours to make and image. Uses a given map of position to colour generated from the image processing function  along with lists of available colours, components and plate types
func PipetteImageGrayscale(OutPlate *wtype.LHPlate, positiontocolourmap map[string]color.Color, water *wtype.LHComponent, black *wtype.LHComponent, volumeperfullcolour wunit.Volume) (finalsolutions []*wtype.LHComponent) {

	solutions := make([]*wtype.LHComponent, 0)

	for locationkey, colour := range positiontocolourmap {

		components := make([]*wtype.LHComponent, 0)

		gray := ColourtoGrayscale(colour)

		if gray.Y < 255 {
			watervol := wunit.NewVolume((float64(255-gray.Y) * volumeperfullcolour.SIValue()), "l")
			waterSample := mixer.Sample(water, watervol)
			components = append(components, waterSample)
		}
		blackvol := wunit.NewVolume((float64(gray.Y) * volumeperfullcolour.SIValue()), "l")
		blackSample := mixer.Sample(black, blackvol)
		components = append(components, blackSample)

		solution := mixer.MixTo(OutPlate.Type, locationkey, 1, components...)
		solutions = append(solutions, solution)
	}

	finalsolutions = solutions
	return
}

// This function used internally to convert any image type to NRGBA if needed.
func toNRGBA(img goimage.Image) *goimage.NRGBA {
	srcBounds := img.Bounds()
	if srcBounds.Min.X == 0 && srcBounds.Min.Y == 0 {
		if src0, ok := img.(*goimage.NRGBA); ok {
			return src0
		}
	}
	return imaging.Clone(img)
}

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
