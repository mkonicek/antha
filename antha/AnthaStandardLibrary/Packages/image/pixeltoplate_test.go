package image

import (
	"context"
	"testing"

	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/download"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/inventory"
	"github.com/antha-lang/antha/inventory/testinventory"

	"image/color"
	"reflect"

	"image/gif"
	"os"
	"image/jpeg"
)

func TestSelectLibrary(t *testing.T) {
	SelectLibrary("UV")
}

func TestSelectColors(t *testing.T) {
	SelectColor("JuniperGFP")
}

func TestMakeAnthaImg(t *testing.T) {
	ctx := testinventory.NewContext(context.Background())

	//downloading image for the test
	imgFile, err := download.File("http://orig08.deviantart.net/a19f/f/2008/117/6/7/8_bit_mario_by_superjerk.jpg", "Downloaded file")
	if err != nil {
		t.Error(err)
	}

	//opening image
	imgBase, err := OpenFile(imgFile)
	if err != nil {
		t.Fatal(err)
	}

	palette := SelectLibrary("UV")

	//initiating components
	var components []*wtype.LHComponent
	component, err := inventory.NewComponent(ctx, "Gluc")
	if err != nil {
		t.Fatal(err)
	}

	//making the array to make palette. It's the same length than the array from the "UV" library
	for i := 1; i < 48; i++ {
		components = append(components, component.Dup())
	}
	//getting palette
	anthaPalette := MakeAnthaPalette(palette, components)

	//getting plate
	plate, err := inventory.NewPlate(ctx, "greiner384")
	if err != nil {
		t.Fatal(err)
	}

	//testing function
	MakeAnthaImg(imgBase, anthaPalette, plate)
}

func TestMakeAnthaPalette(t *testing.T) {
	ctx := testinventory.NewContext(context.Background())

	//getting palette
	palette := SelectLibrary("UV")

	//initiating component
	var components []*wtype.LHComponent
	component, err := inventory.NewComponent(ctx, "Gluc")
	if err != nil {
		t.Fatal(err)
	}

	//making the array to test. It's the same length than the array from the "UV" library
	for i := 1; i < 48; i++ {
		components = append(components, component.Dup())
	}

	//running the function
	MakeAnthaPalette(palette, components)
}

func TestSelectLivingColorLibrary(t *testing.T) {
	ctx := testinventory.NewContext(context.Background())

	SelectLivingColorLibrary(ctx, "ProteinPaintBox")
}

func TestSelectLivingColor(t *testing.T) {
	ctx := testinventory.NewContext(context.Background())
	SelectLivingColor(ctx, "DasherGFP")
}

func TestMakeLivingImg(t *testing.T) {
	ctx := testinventory.NewContext(context.Background())

	//downloading image for the test
	imgFile, err := download.File("http://orig08.deviantart.net/a19f/f/2008/117/6/7/8_bit_mario_by_superjerk.jpg", "Downloaded file")
	if err != nil {
		t.Error(err)
	}

	//opening image
	imgBase, err := OpenFile(imgFile)
	if err != nil {
		t.Error(err)
	}

	//initiating components
	var components []*wtype.LHComponent
	component, err := inventory.NewComponent(ctx, "Gluc")
	if err != nil {
		t.Fatal(err)
	}

	//making the array to make palette. It's the same length than the array from the "ProteinPaintbox" library
	for i := 1; i < 7; i++ {
		components = append(components, component.Dup())
	}

	//Selecting livingPalette
	selectedPalette := SelectLivingColorLibrary(ctx, "ProteinPaintBox")

	//Making palette
	livingPalette := MakeLivingPalette(selectedPalette, components)

	//getting plate
	plate, err := inventory.NewPlate(ctx, "greiner384")
	if err != nil {
		t.Fatal(err)
	}

	//testing function
	MakeLivingImg(imgBase, livingPalette, plate)

}

func TestMakeLivingGIF(t *testing.T) {
	ctx := testinventory.NewContext(context.Background())

	//------------------------------------------------
	//Making antha image
	//------------------------------------------------

	//downloading image for the test
	imgFile, err := download.File("http://orig08.deviantart.net/a19f/f/2008/117/6/7/8_bit_mario_by_superjerk.jpg", "Downloaded file")
	if err != nil {
		t.Error(err)
	}

	//opening image
	imgBase, err := OpenFile(imgFile)
	if err != nil {
		t.Error(err)
	}

	//initiating components
	var components []*wtype.LHComponent
	component, err := inventory.NewComponent(ctx, "Gluc")
	if err != nil {
		t.Fatal(err)
	}

	//making the array to make palette. It's the same length than the array from the "ProteinPaintbox" library
	for i := 1; i < 7; i++ {
		components = append(components, component.Dup())
	}

	//Selecting livingPalette
	selectedPalette := SelectLivingColorLibrary(ctx, "ProteinPaintBox")

	//Making palette
	livingPalette := MakeLivingPalette(selectedPalette, components)

	//getting plate
	plate, err := inventory.NewPlate(ctx, "greiner384")
	if err != nil {
		t.Fatal(err)
	}

	//generating images. We only use 2 since that's what we use for our construct
	anthaImg1, _ := MakeLivingImg(imgBase, livingPalette, plate)
	anthaImg2, _ := MakeLivingImg(imgBase, livingPalette, plate)

	//Merge them
	var anthaImgs []LivingImg
	anthaImgs = append(anthaImgs, *anthaImg1, *anthaImg2)

	//------------------------------------------------
	//Testing GIF functions
	//------------------------------------------------

	MakeLivingGIF(anthaImgs)

}

func TestAnthaPrintWorkflow(t *testing.T){

	ctx := testinventory.NewContext(context.Background())

	//downloading image for the test
	imgFile, err := download.File("http://orig08.deviantart.net/a19f/f/2008/117/6/7/8_bit_mario_by_superjerk.jpg", "Downloaded file")
	if err != nil {
		t.Error(err)
	}

	//opening image
	imgBase, err := OpenFile(imgFile)
	if err != nil {
		t.Fatal(err)
	}

	palette := SelectLibrary("UV")

	//initiating components
	var components []*wtype.LHComponent
	component, err := inventory.NewComponent(ctx, "Gluc")
	if err != nil {
		t.Fatal(err)
	}

	//making the array to make palette. It's the same length than the array from the "UV" library
	for i := 1; i < 48; i++ {
		components = append(components, component.Dup())
	}
	//getting palette
	anthaPalette := MakeAnthaPalette(palette, components)

	//getting plate
	plate, err := inventory.NewPlate(ctx, "greiner384")
	if err != nil {
		t.Fatal(err)
	}

	//testing function
	MakeAnthaImg(imgBase, anthaPalette, plate)


	//------------------------------------------------------------------
	//ELEMENT TESTING
	//------------------------------------------------------------------

	//testing function
	AnthaImage, _ := MakeAnthaImg(imgBase, anthaPalette, plate)
	ColorToExclude := "white"    //Background color which we do not want to print, leave blank if you want to print all colors.

	//------------------------------------------------------------------
	//Globals
	//------------------------------------------------------------------

	var colorToExclude color.Color

	if ColorToExclude == "" {
		colorToExclude = color.NRGBA{}
	} else {
		colorToExclude = SelectColor(ColorToExclude)
	}

	//------------------------------------------------------------------
	//Iterating through each pixels in the image and pipetting them
	//------------------------------------------------------------------

	counter := 0
	for _, pix := range AnthaImage.Pix {
		r,g,b,a := pix.Color.Color.RGBA()
		asRGBA := color.RGBA{uint8(r),uint8(g),uint8(b),uint8(a)}

		//Skipping pixel if it is of the background color
		if reflect.DeepEqual(colorToExclude, asRGBA) {
			continue
		}

		counter++
	}

	t.Log(counter)
}

func TestParseGIF(t *testing.T){

	//downloading GIF for the test
	//RAINBOWSPINNER
	//https://media.giphy.com/media/XUHmgf1ij7dOU/source.gif
	//BUTTERFLY
	//https://www.google.co.uk/search?q=artsy+gif&source=lnms&tbm=isch&sa=X&ved=0ahUKEwi08uHk1szVAhUBAsAKHR6DAWUQ_AUICigB&biw=1301&bih=654#imgrc=iWSaQhyp0mvc5M:
	GIFFile, err := download.File("https://media.giphy.com/media/XUHmgf1ij7dOU/source.gif", "Downloaded GIF")
	if err != nil {
		t.Error(err)
	}

	//opening GIF
	GIF, err := OpenGIF(GIFFile)
	if err != nil {
		t.Error(err)
	}

	//saving Image
	fimg, err := os.Create("/home/cachemoi/gocode/src/github.com/cachemoi/GIF/doc/img/GIFImgTest.jpg")
	jpeg.Encode(fimg, GIF.Image[5], &jpeg.Options{jpeg.DefaultQuality})

	//saving GIF
	//There's a bug with the Linux Library which means you'll have to open the resulting file in a browser
	fgif, err := os.Create("/home/cachemoi/gocode/src/github.com/cachemoi/GIF/doc/GIF/results.gif")
	err = gif.EncodeAll(fgif,GIF)
	if err != nil {
		t.Error(err)
	}




}