package image

import (
	"context"
	"testing"

	"fmt"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/inventory"
	"github.com/antha-lang/antha/inventory/testinventory"
	"io/ioutil"
	"path/filepath"
)

func TestSelectLibrary(t *testing.T) {
	SelectLibrary("UV")
}

func TestSelectColors(t *testing.T) {
	SelectColor("JuniperGFP")
}

func TestMakeAnthaImg(t *testing.T) {
	ctx := testinventory.NewContext(context.Background())

	var imgFile wtype.File

	//getting image file for the test
	absPath, _ := filepath.Abs("../image/testdata/spaceinvader.png")
	bytes, err := ioutil.ReadFile(absPath)
	if err != nil {
		t.Error(err)
	}

	imgFile.WriteAll(bytes)

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
	for i := 1; i <= len(SelectLibrary("UV")); i++ {

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

	var imgFile wtype.File

	//getting image file for the test
	absPath, _ := filepath.Abs("../image/testdata/spaceinvader.png")
	bytes, err := ioutil.ReadFile(absPath)
	if err != nil {
		t.Error(err)
	}

	imgFile.WriteAll(bytes)

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

	var gifFile wtype.File

	//getting gif file for the test
	absPath, _ := filepath.Abs("../image/testdata/deepDream.gif")
	bytes, err := ioutil.ReadFile(absPath)
	if err != nil {
		t.Error(err)
	}

	gifFile.WriteAll(bytes)

	//opening image
	imgBase, err := OpenFile(gifFile)
	if err != nil {
		t.Fatal(err)
	}

	//initiating components
	var components []*wtype.LHComponent
	component, err := inventory.NewComponent(ctx, "Gluc")
	if err != nil {
		t.Fatal(err)
	}

	//making the array to make palette. It's the same length than the array from the "ProteinPaintbox" library
	for i := 1; i < len(livingColors); i++ {
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

func TestAnthaPrintWorkflow(t *testing.T) {

	ctx := testinventory.NewContext(context.Background())

	var imgFile wtype.File

	//getting image for the test
	absPath, _ := filepath.Abs("../image/testdata/spaceinvader.png")
	bytes, err := ioutil.ReadFile(absPath)
	if err != nil {
		t.Error(err)
	}

	imgFile.WriteAll(bytes)

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

func TestOpenGIF(t *testing.T) {

	var gifFile wtype.File

	//getting gif file for the test
	absPath, _ := filepath.Abs("../image/testdata/deepDream.gif")
	bytes, err := ioutil.ReadFile(absPath)
	if err != nil {
		t.Error(err)
	}

	gifFile.WriteAll(bytes)

	//opening GIF
	OpenGIF(gifFile)

}

func TestParseGIF(t *testing.T) {

	var gifFile wtype.File

	//getting gif for the test
	absPath, _ := filepath.Abs("../image/testdata/deepDream.gif")
	bytes, err := ioutil.ReadFile(absPath)
	if err != nil {
		t.Error(err)
	}

	gifFile.WriteAll(bytes)

	//opening image
	GIF, err := OpenGIF(gifFile)
	if err != nil {
		t.Error(err)
	}

	ParseGIF(GIF, []int{1, 6})

}

func TestGetState(t *testing.T) {
	ctx := testinventory.NewContext(context.Background())

	lp := SelectLivingColorLibrary(ctx, "ProteinPaintBox")
	lip, _ := inventory.NewPlate(ctx, "greiner384")

	var gifFile wtype.File

	//getting gif for the test
	absPath, _ := filepath.Abs("../image/testdata/deepDream.gif")
	bytes, err := ioutil.ReadFile(absPath)
	if err != nil {
		t.Error(err)
	}

	gifFile.WriteAll(bytes)

	//opening image
	GIF, err := OpenGIF(gifFile)
	if err != nil {
		t.Error(err)
	}

	imgs, err := ParseGIF(GIF, []int{1, 4})
	if err != nil {
		fmt.Println(err)
	}

	li1, _ := MakeLivingImg(imgs[0], &lp, lip)
	li2, _ := MakeLivingImg(imgs[1], &lp, lip)

	var lia []LivingImg

	lia = append(lia, *li1, *li2)
	//------------------------------------------------
	//Testing GIF functions
	//------------------------------------------------

	LivingGIF := MakeLivingGIF(lia)

	LivingGIF.GetStates()

}
