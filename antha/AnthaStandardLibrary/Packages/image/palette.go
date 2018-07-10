package image

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"image/color"
	"image/color/palette"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"

	anthapath "github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/AnthaPath"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
)

// An AnthaColor represents a color linked to an LHComponent
type AnthaColor struct {
	Color     color.NRGBA
	Component *wtype.Liquid
}

// An AnthaPalette is an array of anthaColors
type AnthaPalette struct {
	AnthaColors []AnthaColor
	Palette     color.Palette
}

// sqDiff returns the squared-difference of x and y, shifted by 2 so that
// adding four of those won't overflow a uint32.
//
// an internal goimage function
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

// AvailablePalettes returns all the (hardcoded) palettes in the library. The
// keys are the pallette names.
func AvailablePalettes() map[string]color.Palette {

	var avail = make(map[string]color.Palette)

	avail["Palette1"] = paletteFromMap(colourComponentMap) //Chosencolourpalette,
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

// AvailableComponentMaps returns available palettes in a map.
func AvailableComponentMaps() map[string]map[color.Color]string {
	var componentMaps = make(map[string]map[color.Color]string)
	componentMaps["Palette1"] = colourComponentMap
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

// MakeSubMapFromMap extracts colors from a color library Given colour names.
// Bad name.
func MakeSubMapFromMap(colourMap map[color.Color]string, colourNames []string) map[color.Color]string {
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

// MakeSubPalette extracts colors from a palette given a slice of colour
// names.
func MakeSubPalette(paletteName string, colourNames []string) color.Palette {
	paletteMap := AvailableComponentMaps()[paletteName]

	subMap := MakeSubMapFromMap(paletteMap, colourNames)

	return paletteFromMap(subMap)
}

// RemoveDuplicatesKeysFromMap will loop over a map of colors to find and
// delete duplicates. Entries with duplicate keys are deleted.
func RemoveDuplicatesKeysFromMap(elements map[string]color.Color) map[string]color.Color {
	encountered := make(map[string]bool)
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

// RemoveDuplicatesValuesFromMap will loop over a map of colors to find and
// delete duplicates. Entries with duplicate values are deleted.
func RemoveDuplicatesValuesFromMap(elements map[string]color.Color) map[string]color.Color {
	encountered := make(map[color.Color]bool)
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

// paletteFromColorarray makes a palette from a slice of colors.
func paletteFromColorarray(colors []color.Color) color.Palette {
	var newPalette color.Palette = colors
	return newPalette
}

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

// SelectColor returns the desired color from a string ID
func SelectColor(colID string) color.Color {
	selectedColor, found := colourLibrary[colID]
	if !found {
		panic(fmt.Sprintf("library %s not found so could not make palette", colID))
	}

	return selectedColor
}
