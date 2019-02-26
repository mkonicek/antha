package image

import (
	"context"
	"io/ioutil"
	"testing"

	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/inventory"
	"github.com/antha-lang/antha/inventory/testinventory"
)

func TestSelectLibrary(t *testing.T) {
	SelectLibrary("UV")
}

func TestSelectColors(t *testing.T) {
	SelectColor("JuniperGFP")
}

func TestMakeAnthaImg(t *testing.T) {
	ctx := testinventory.NewContext(context.Background())

	data, err := ioutil.ReadFile("testdata/invader.png")
	if err != nil {
		t.Fatal(err)
	}

	var imgFile wtype.File
	if err := imgFile.WriteAll(data); err != nil {
		t.Fatal(err)
	}

	//opening image
	imgBase, err := OpenFile(imgFile)
	if err != nil {
		t.Fatal(err)
	}

	palette := SelectLibrary("UV")

	//initiating components
	var components []*wtype.Liquid
	component, err := inventory.NewComponent(ctx, "Gluc")
	if err != nil {
		t.Fatal(err)
	}

	//making the array to make palette. It's the same length than the array from the "UV" library
	for i := 1; i <= len(palette); i++ {

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
	var components []*wtype.Liquid
	component, err := inventory.NewComponent(ctx, "Gluc")
	if err != nil {
		t.Fatal(err)
	}

	//making the array to test. It's the same length than the array from the "UV" library
	for i := 1; i <= len(palette); i++ {
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

	data, err := ioutil.ReadFile("testdata/invader.png")
	if err != nil {
		t.Fatal(err)
	}

	var imgFile wtype.File
	if err := imgFile.WriteAll(data); err != nil {
		t.Fatal(err)
	}

	//opening image
	imgBase, err := OpenFile(imgFile)
	if err != nil {
		t.Fatal(err)
	}

	//initiating components
	var components []*wtype.Liquid
	component, err := inventory.NewComponent(ctx, "Gluc")
	if err != nil {
		t.Fatal(err)
	}

	//Selecting livingPalette
	selectedPalette := SelectLivingColorLibrary(ctx, "ProteinPaintBox")

	//making the array to make palette. It's the same length than the array from the "ProteinPaintbox" library
	for i := 1; i <= len(selectedPalette.LivingColors); i++ {
		components = append(components, component.Dup())
	}

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

	data, err := ioutil.ReadFile("testdata/invader.png")
	if err != nil {
		t.Fatal(err)
	}

	var imgFile wtype.File
	if err := imgFile.WriteAll(data); err != nil {
		t.Fatal(err)
	}

	//opening image
	imgBase, err := OpenFile(imgFile)
	if err != nil {
		t.Fatal(err)
	}

	//initiating components
	var components []*wtype.Liquid
	component, err := inventory.NewComponent(ctx, "Gluc")
	if err != nil {
		t.Fatal(err)
	}

	//Selecting livingPalette
	selectedPalette := SelectLivingColorLibrary(ctx, "ProteinPaintBox")

	//making the array to make palette. It's the same length than the array from the "ProteinPaintbox" library
	for i := 1; i <= len(selectedPalette.LivingColors); i++ {
		components = append(components, component.Dup())
	}

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

	data, err := ioutil.ReadFile("testdata/invader.png")
	if err != nil {
		t.Fatal(err)
	}

	var imgFile wtype.File
	if err := imgFile.WriteAll(data); err != nil {
		t.Fatal(err)
	}

	//opening image
	imgBase, err := OpenFile(imgFile)
	if err != nil {
		t.Fatal(err)
	}

	palette := SelectLibrary("UV")

	//initiating components
	var components []*wtype.Liquid
	component, err := inventory.NewComponent(ctx, "Gluc")
	if err != nil {
		t.Fatal(err)
	}

	//making the array to make palette. It's the same length than the array from the "UV" library
	for i := 1; i <= len(palette); i++ {
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

func TestParseGIF(t *testing.T) {
	data, err := ioutil.ReadFile("testdata/potato.gif")
	if err != nil {
		t.Fatal(err)
	}

	var imgFile wtype.File
	if err := imgFile.WriteAll(data); err != nil {
		t.Fatal(err)
	}

	//opening GIF
	gif, err := OpenGIF(imgFile)
	if err != nil {
		t.Fatal(err)
	}

	if _, err := ParseGIF(gif, []int{1, 6}); err != nil {
		t.Fatal(err)
	}
}

func TestGetState(t *testing.T) {
	ctx := testinventory.NewContext(context.Background())

	data, err := ioutil.ReadFile("testdata/potato.gif")
	if err != nil {
		t.Fatal(err)
	}

	var imgFile wtype.File
	if err := imgFile.WriteAll(data); err != nil {
		t.Fatal(err)
	}

	//opening GIF
	gif, err := OpenGIF(imgFile)
	if err != nil {
		t.Fatal(err)
	}

	imgs, err := ParseGIF(gif, []int{1, 2})
	if err != nil {
		t.Fatal(err)
	}

	lp := SelectLivingColorLibrary(ctx, "ProteinPaintBox")
	lip, _ := inventory.NewPlate(ctx, "greiner384")

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
