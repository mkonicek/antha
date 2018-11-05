package image

import (
	"context"
	"errors"
	"fmt"
	goimage "image"
	"image/color"
	"reflect"
	"regexp"
	"strconv"

	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wutil"
	"github.com/antha-lang/antha/inventory"
	"github.com/disintegration/imaging"
)

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

// ResizeImageToPlate resizes an image to fit the number of wells on a plate.
// We treat wells as pixels. Bad name.
func ResizeImageToPlate(img goimage.Image, plate *wtype.Plate, algorithm imaging.ResampleFilter, rotate bool) *goimage.NRGBA {

	if img.Bounds().Dy() == plate.WellsY() {
		return toNRGBA(img)
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

// ResizeImageToPlateAutoRotate resizes an image to fit the number of wells on
// a plate and if the image is in portrait the image will be rotated by 270
// degrees to optimise resolution on the plate.
func ResizeImageToPlateAutoRotate(img goimage.Image, plate *wtype.Plate, algorithm imaging.ResampleFilter) *goimage.NRGBA {
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

// CheckAllResizeAlgorithms resizes an image using a variety of different
// algorithms.
func CheckAllResizeAlgorithms(img goimage.Image, plate *wtype.Plate, rotate bool, algorithms map[string]imaging.ResampleFilter) []*goimage.NRGBA {
	var plateImages []*goimage.NRGBA

	for _, algorithm := range algorithms {
		resized := ResizeImageToPlate(img, plate, algorithm, rotate)

		plateImages = append(plateImages, resized)

	}

	return plateImages
}

// MakePaletteFromImage will make a color Palette from an image resized to fit
// a given plate type.
func MakePaletteFromImage(img goimage.Image, plate *wtype.Plate, rotate bool) color.Palette {
	plateImage := ResizeImageToPlate(img, plate, imaging.CatmullRom, rotate)

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

// MakeSmallPaletteFromImage will make a color Palette from an image resized to
// fit a given plate type using Plan9 Palette
func MakeSmallPaletteFromImage(img goimage.Image, plate *wtype.Plate, rotate bool) color.Palette {
	plateImage := ResizeImageToPlate(img, plate, imaging.CatmullRom, rotate)

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

// ToPlateLayout takes an image, plate type, and palette and returns a map
// of well position to colors. Creates a map of pixel to plate position from
// processing a given image with a chosen colour palette.  It's recommended to
// use at least 384 well plate. If autorotate == true, rotate is overridden
func ToPlateLayout(img goimage.Image, plate *wtype.Plate, palette *color.Palette, rotate bool, autoRotate bool) (map[string]color.Color, *goimage.NRGBA) {
	var plateImage *goimage.NRGBA
	if autoRotate {
		plateImage = ResizeImageToPlateAutoRotate(img, plate, imaging.CatmullRom)
	} else {
		plateImage = ResizeImageToPlate(img, plate, imaging.CatmullRom, rotate)
	}

	// make map of well position to colour: (array for time being)
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
			positionToColour[pos] = colour
		}
	}

	return positionToColour, plateImage
}

// PrintFPImagePreview takes an image, plate type, and colors from
// visiblemap/uvmap and (use to) save the resulting processed image.
func PrintFPImagePreview(img goimage.Image, plate *wtype.Plate, rotate bool, visiblemap, uvmap map[color.Color]string) {

	plateimage := ResizeImageToPlate(img, plate, imaging.CatmullRom, rotate)

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
			if !ok {
				errmessage := fmt.Sprintln("colourstring", colourstring, "not found in map", stringkeymap, "len", len(stringkeymap))
				panic(errmessage)
			}
			// fmt.Println("viscolour", viscolour)
			plateimage.Set(x, y, viscolour)

		}
	}
}

// An AnthaPix is a pixel linked to an anthaColor, and thereby an LHComponent
type AnthaPix struct {
	Color    AnthaColor
	Location wtype.WellCoords
}

// An AnthaImg represents an image where pixels are linked
type AnthaImg struct {
	Plate   *wtype.Plate
	Pix     []AnthaPix
	Palette AnthaPalette
}

// A LivingColor links a color to physical information such as DNASequence and LHcomponent
type LivingColor struct {
	ID        string
	Color     color.NRGBA
	Seq       wtype.DNASequence
	Component *wtype.Liquid
}

// A LivingPalette holds an array of LivingColors
type LivingPalette struct {
	LivingColors []LivingColor
}

// A LivingPix represents a pixel along with well coordinate
type LivingPix struct {
	Color    LivingColor
	Location wtype.WellCoords
}

// A LivingImg is a representation of an image linked to biological data
type LivingImg struct {
	Plate   wtype.Plate
	Pix     []LivingPix
	Palette LivingPalette
}

//A LivingGIF is a representation of a GIF linked to biological data
type LivingGIF struct {
	Frames []LivingImg
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
	return reflect.DeepEqual(p1, p2)
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
func MakeAnthaPalette(palette color.Palette, components []*wtype.Liquid) *AnthaPalette {

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
func MakeAnthaImg(goImg goimage.Image, anthaPalette *AnthaPalette, anthaImgPlate *wtype.Plate) (outputImg *AnthaImg, resizedImg *goimage.NRGBA) {

	//Global placeholders
	var anthaPix AnthaPix
	var anthaImgPix []AnthaPix
	var anthaImg AnthaImg

	//Verify that the plate is the same size as the digital image. If not resize.
	if goImg.Bounds().Dy() != anthaImgPlate.WellsY() {
		goImg = ResizeImageToPlateMin(goImg, anthaImgPlate)
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

	return &anthaImg, goImg.(*goimage.NRGBA)
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
func MakeLivingPalette(palette LivingPalette, components []*wtype.Liquid) *LivingPalette {

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
func MakeLivingImg(goImg goimage.Image, livingPalette *LivingPalette, livingImgPlate *wtype.Plate) (outputImg *LivingImg, resizedImg goimage.Image) {

	//Global placeholders
	var livingPix LivingPix
	var livingImgPix []LivingPix
	var livingImg LivingImg

	//Verify that the plate is the same size as the digital image. If not resize.
	if goImg.Bounds().Dy() != livingImgPlate.WellsY() {
		goImg = ResizeImageToPlateMin(goImg, livingImgPlate)
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

	livingGIF.Frames = append(livingGIF.Frames, imgs...)

	return &livingGIF
}

// ResizeImageToPlateMin is a minimalist resize function. Uses Lanczos
// resampling, which is the best but slowest method.
func ResizeImageToPlateMin(img goimage.Image, plate *wtype.Plate) *goimage.NRGBA {

	if img.Bounds().Dy() == plate.WellsY() {
		return toNRGBA(img)
	}

	if img.Bounds().Dy() > img.Bounds().Dx() {
		img = imaging.Rotate270(img)
	}

	ratio := float64(img.Bounds().Dx()) / float64(img.Bounds().Dy())

	if ratio <= float64(plate.WellsX())/float64(plate.WellsY()) {
		return imaging.Resize(img, 0, plate.WlsY, imaging.Lanczos)
	}
	return imaging.Resize(img, plate.WlsX, 0, imaging.Lanczos)
}
