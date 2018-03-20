package image

import (
	"image/color"
	"os"
	"path/filepath"

	anthapath "github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/AnthaPath"
	imaging "github.com/disintegration/imaging"
)

var (
	// White is a color.RGBA representation of the colour white
	White = color.RGBA{R: uint8(242), G: uint8(243), B: uint8(242), A: uint8(255)}
	// Black is a color.RGBA representation of the colour black
	Black = color.RGBA{R: uint8(0), G: uint8(0), B: uint8(0), A: uint8(255)}
	// Transparent is a color.RGBA representation of a nil colour.
	Transparent = color.RGBA{R: uint8(0), G: uint8(0), B: uint8(0), A: uint8(0)}
)

// colourComponentMap is a map of RGB colour to description for use as key in
// crossreferencing colour to component in other maps.
var colourComponentMap = map[color.Color]string{
	color.RGBA{R: uint8(13), G: uint8(105), B: uint8(171), A: uint8(255)}: "blue",
	color.RGBA{R: uint8(245), G: uint8(205), B: uint8(47), A: uint8(255)}: "yellow",
	color.RGBA{R: uint8(75), G: uint8(151), B: uint8(74), A: uint8(255)}:  "green",
	color.RGBA{R: uint8(196), G: uint8(40), B: uint8(27), A: uint8(255)}:  "red",
	White:       "white",
	Black:       "black",
	Transparent: "transparent",
}

var neonComponentMap = map[color.Color]string{
	Black: "black",
	color.RGBA{R: uint8(149), G: uint8(156), B: uint8(161), A: uint8(255)}: "grey",
	color.RGBA{R: uint8(117), G: uint8(51), B: uint8(127), A: uint8(255)}:  "purple",
	color.RGBA{R: uint8(25), G: uint8(60), B: uint8(152), A: uint8(255)}:   "darkblue",
	color.RGBA{R: uint8(0), G: uint8(125), B: uint8(200), A: uint8(255)}:   "blue",
	color.RGBA{R: uint8(0), G: uint8(177), B: uint8(94), A: uint8(255)}:    "green",
	color.RGBA{R: uint8(244), G: uint8(231), B: uint8(0), A: uint8(255)}:   "yellow",
	color.RGBA{R: uint8(255), G: uint8(118), B: uint8(0), A: uint8(255)}:   "orange",
	color.RGBA{R: uint8(255), G: uint8(39), B: uint8(51), A: uint8(255)}:   "red",
	color.RGBA{R: uint8(235), G: uint8(41), B: uint8(123), A: uint8(255)}:  "pink",
	White: "white",
	color.RGBA{R: uint8(0), G: uint8(174), B: uint8(239), A: uint8(255)}: "Cyan",
	color.RGBA{R: uint8(236), G: uint8(0), B: uint8(140), A: uint8(255)}: "Magenta",
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
	Black: "black",
	// plus white as a blank (or comment out to use EiraCFP)
	White: "white",
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
	White: "white",
	Black: "black",
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
	Black: "black",

	// plus white as a blank (or comment out to use EiraCFP)
	//color.RGBA{R: uint8(242), G: uint8(243), B: uint8(242), A: uint8(255)}: "white",
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
	"black":       Black,
	"transparent": Transparent,
	"white":       White,

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
	"UVblack":    Black,
	"UVgreen":    color.RGBA{R: uint8(0), G: uint8(177), B: uint8(94), A: uint8(255)},
	"UVyellow":   color.RGBA{R: uint8(244), G: uint8(231), B: uint8(0), A: uint8(255)},
	"UVorange":   color.RGBA{R: uint8(255), G: uint8(118), B: uint8(0), A: uint8(255)},
	"UVwhite":    White,
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
		"white",
		"DreidelTeal",
		"LeorOrange",
		"DasherGFP",
		"BlazeYFP",
		"black",
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
	"DasherGFP": {
		Color: &color.NRGBA{R: uint8(0), G: uint8(255), B: uint8(0), A: uint8(255)},
		Seq:   "ATGACGGCATTGACGGAAGGTGCAAAACTGTTTGAGAAAGAGATCCCGTATATCACCGAACTGGAAGGCGACGTCGAAGGTATGAAATTTATCATTAAAGGCGAGGGTACCGGTGACGCGACCACGGGTACCATTAAAGCGAAATACATCTGCACTACGGGCGACCTGCCGGTCCCGTGGGCAACCCTGGTGAGCACCCTGAGCTACGGTGTTCAGTGTTTCGCCAAGTACCCGAGCCACATCAAGGATTTCTTTAAGAGCGCCATGCCGGAAGGTTATACCCAAGAGCGTACCATCAGCTTCGAAGGCGACGGCGTGTACAAGACGCGTGCTATGGTTACCTACGAACGCGGTTCTATCTACAATCGTGTCACGCTGACTGGTGAGAACTTTAAGAAAGACGGTCACATTCTGCGTAAGAACGTTGCATTCCAATGCCCGCCAAGCATTCTGTATATTCTGCCTGACACCGTTAACAATGGCATCCGCGTTGAGTTCAACCAGGCGTACGATATTGAAGGTGTGACCGAAAAACTGGTTACCAAATGCAGCCAAATGAATCGTCCGTTGGCGGGCTCCGCGGCAGTGCATATCCCGCGTTATCATCACATTACCTACCACACCAAACTGAGCAAAGACCGCGACGAGCGCCGTGATCACATGTGTCTGGTAGAGGTCGTGAAAGCGGTTGATCTGGACACGTATCAGTAATGAGAATTCTGTACACTCGAG",
	},
	"RudolphRFP": {
		Color: &color.NRGBA{R: uint8(218), G: uint8(92), B: uint8(69), A: uint8(255)},
		Seq:   "atgtccctgtcgaaacaagtactgccacacgatgttaagatgcgctatcatatggatggctgcgttaatggtcattctttcaccattgagggtgaaggtgcaggcaaaccgtatgagggcaagaagatcttggaactgcgcgtgacgaaaggtggcccgctgccttttgcgttcgatatcctgagcagcgtttttacctacggtaaccgttgtttttgcgagtatccagaggacatgccggactactttaaacagagcctgccggaaggtcattcttgggaacgcaccctgatgtttgaggatggcggttgtggtacggcgagcgcgcacatttccctggacaagaactgcttcgtgcacaagagcaccttccacggcgtcaatttcccggcaaacggtccggtcatgcaaaagaaagctatgaactgggagccgagcagcgaactgattacggcgtgcgacggtatcctgaaaggcgatgtgaccatgtttctgttgctggaaggtggccaccgtcttaaatgtcagttcaccaccagctacaaagcccacaaggcagttaagatgccgccgaatcacattatcgaacacgtgcttgttaaaaaagaggttgccgacggctttcagatccaagagcatgcggtcgcaaagcacttcaccgtcgacgttaaagaaacgtaatgagaattctgtacactcgag",
	},
	"FresnoRFP": {
		Color: &color.NRGBA{R: uint8(255), G: uint8(0), B: uint8(166), A: uint8(255)},
		Seq:   "AATAGCCTGATTAAAGAGAATATGCACATGAAGCTGTACATGGAAGGCACGGTGAATAACCACCACTTCAAATGCACCAGCGAGGGTGAGGGTAAACCGTATGAAGGCACCCAAACGATGCGTATCAAAGTTGTTGAGGGTGGCCCGTTGCCGTTTGCGTTCGACATTTTAGCGACGAGCTTTATGTATGGCTCTCGTACGTTTATCAAGTACCCGAAGGGTATTCCGGACTTTTTCAAACAATCTTTTCCAGAGGGTTTCACCTGGGAGCGCGTGACTCGCTACGAAGATGGCGGCGTCGTGACCGCAACGCAGGATACCTCCCTGGAAGATGGCTGCCTGGTCTACCACGTTCAGGTCCGTGGTGTCAATTTCCCGAGCAATGGTCCGGTTATGCAGAAGAAAACCCTGGGTTGGGAACCGAACACCGAGATGTTGTATCCTGCAGATGGTGGCCTGGAAGGTCGCAGCGACATGGCATTGAAACTGGTCGGTGGCGGCCATCTGAGCTGTAGCTTCGTGACCACGTATCGTTCGAAGAAAACGGTCGGTAACATCAAAATGCCGGGTATTCACGCGGTTGACCACCGTCTGGTGCGCATTAAAGAAGCCGACAAAGAGACTTACGTGGAGCAACATGAAGTAGCCGTTGCGAAATTTGCTGGTTTGGGCGGTGGTATGGACGAACTGTACAAATGAGAATTCTGTACACTCGAG",
	},
	"MarleyYFP": {
		Color: &color.NRGBA{R: uint8(236), G: uint8(255), B: uint8(0), A: uint8(255)},
		Seq:   "ACGGCATTGACGGAAGGTGCAAAACTGTTTGAGAAAGAGATCCCGTATATCACCGAACTGGAAGGCGACGTCGAAGGTATGAAATTTATCATTAAAGGCGAGGGTACCGGTGACGCGAGCGTTGGTAAGGTCGACGCGCAGTTTATCTGCACCACGGGTGACGTCCCGGTGCCGTGGAGCACGTTGGTGACGACGCTGACTTACGGTGCCCAATGTTTCGCGAAATATCCGCGTCACATCGCAGACTTCTTCAAATCCTGTATGCCGGAAGGTTATGTGCAAGAGCGCACTATTACCTTTGAGGGTGACGGTGTCTTTAAGACCCGTGCCGAAGTGACCTTCGAAAACGGTAGCGTTTACAACCGTGTCAAGCTGAATGGCCAGGGTTTCAAAAAGGATGGTCACGTTCTGGGTAAGAATCTGGAATTCAACTTCACCCCACACTGCCTGTACATCTGGGGCGATCAAGCGAATCATGGTCTGAAAAGCGCATTTAAGGTGATGCACGAGATTACCGGCTCCAAAGAAGATTTCATCGTGGCTGATCACACCCAGATGAATACCCCGATTGGCGGTGGCCCTGTGCACGTTCCGGAATATCATCACATTACCTACCATGTCACGCTGAGCAAGGATGTTACCGACCATCGCGACCACCTGAACATCGTAGAGGTTATCAAAGCAGTGGACCTGGAGACGTATCGCTAATGAGAATTCTGTACACTCGAG",
	},
	"TwinkleCFP": {
		Color: &color.NRGBA{R: uint8(27), G: uint8(79), B: uint8(146), A: uint8(255)},
		Seq:   "TCGTCTGGTGCCAAATTGTTTGAAAAGAAGATCCCGTATATCACTGAACTCGAGGGCGACGTCAATGGTATGAAGTTTACCATTCATGGTAAAGGTACCGGCGATGCGACCACGGGTAAAATTAAAGCGCAGTTCATCTGCACTACGGGCGACGTTCCGGTCCCGTGGGCAACCCTGGTGAGCACCCTGAGCTACGGTGTTCAGTGTTTCGCCAAGTACCCGAGCGAGCTGAAGGATTTCTTTAAGAGCTGCATGCCGGAAGGTTATGTTCAAGAGCGTACCATCACCTTCGAAGGCGACGGCGTGTTTAGGACGCGTGCTGAGGTTACCTTTGAAAACGGTTCTGTCTACAATCGTGTCACTCTGACTGGCGAGGGTTTTAAGAAAGACGGTCACATTCTGCGTAAGAACGTTGCATTCCAATGCCCGCCAGATATTGTGTATATTCTGCCTGACACCGTTAACAATGGCATCCGCGTTGAGTTCAACCAGGCGTACGATATTGAAGGTGTGACCGAAAAACTGGTTACCAAATGCAGCCAAATGAATCGTCCGTTGGCGGGCTCCGCGGCAGTGCATATCCCGCGTTATCATCACCTGAGCAAACACACCAAACTGAGCAAAGACCGCGACGAGCGCCGTGATCACATGTGTCTGGTAGAGGTCGTGAAAGCGGTTGATCTGGACACGTATCAGTAATGAGAATTCTGTACACTCGAG",
	},
	"mTagBFP": {
		Color: &color.NRGBA{},
		Seq:   "GTGTCTAAGGGCGAGGAGCTGATTAAGGAGAACATGCACATGAAGCTGTACATGGAGGGCACCGTGGACAACCATCACTTCAAGTGCACATCCGAGGGCGAAGGCAAGCCCTACGAGGGCACCCAGACCATGAGAATCAAGGTGGTCGAGGGCGGCCCTCTCCCCTTCGCCTTCGACATCCTGGCTACTAGCTTCCTCTACGGCAGCAAGACCTTCATCAACCACACCCAGGGCATCCCCGACTTCTTCAAGCAGTCCTTCCCTGAGGGCTTCACATGGGAGAGAGTCACCACATACGAAGACGGGGGCGTGCTGACCGCTACCCAGGACACCAGCCTCCAGGACGGCTGCCTCATCTACAACGTCAAGATCAGAGGGGTGAACTTCACATCCAACGGCCCTGTGATGCAGAAGAAAACACTCGGCTGGGAGGCCTTCACCGAGACGCTGTACCCCGCTGACGGCGGCCTGGAAGGCAGAAACGACATGGCCCTGAAGCTCGTGGGCGGGAGCCATCTGATCGCAAACGCCAAGACCACATATAGATCCAAGAAACCCGCTAAGAACCTCAAGATGCCTGGCGTCTACTATGTGGACTACAGACTGGAAAGAATCAAGGAGGCCAACAACGAAACCTACGTCGAGCAGCACGAGGTGGCAGTGGCCAGATACTGCGACCTCCCTAGCAAACTGGGGCACAAGCTTAATTAATGAGAATTCTGTACACTCGAG",
	},
	//since this should be the background this is set to pure white.
	"Nil": {
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

// VisibleEquivalentMaps returns just the proteinPaintboxUV color map
func VisibleEquivalentMaps() map[string]map[color.Color]string {
	m := make(map[string]map[color.Color]string)
	m["ProteinPaintboxUV"] = proteinPaintboxMap

	if _, err := os.Stat(filepath.Join(anthapath.Path(), "testcolours.json")); err == nil {
		invmap, err := makeLatestColourMap(filepath.Join(anthapath.Path(), "testcolours.json"))
		if err != nil {
			panic(err.Error())
		}
		m["UVinventory"] = invmap
	}
	return m
}

// ColourComponentMap is a map of RGB colour to description for use as key in
// crossreferencing colour to component in other maps.
func ColourComponentMap() map[color.Color]string {
	return colourComponentMap
}
