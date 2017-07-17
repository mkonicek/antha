package image

import (
	"testing"
	"reflect"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/microArch/factory"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/download"
)

func TestSelectLibrary (t *testing.T) {

	palette, err := SelectLibrary("UV")
	t.Log(palette, err)
	t.Log(reflect.TypeOf(palette))

}

func TestSelectColors(t *testing.T) {

	palette, err :=  SelectColor("JuniperGFP")
	t.Log(err)

	t.Log(palette)
}

func TestMakeAnthaImg(t *testing.T) {

	//downloading image for the test
	imgFile , err := download.File("http://orig08.deviantart.net/a19f/f/2008/117/6/7/8_bit_mario_by_superjerk.jpg", "Downloaded file")
	if err != nil{
		t.Error(err)
	}

	//opening image
	imgBase, err := OpenFile(imgFile)
	if err != nil{
		t.Error(err)
	}

	palette, err := SelectLibrary("UV")

	t.Log(err)
	//initiating components
	var components []*wtype.LHComponent
	component := factory.GetComponentByType("Gluc")
	//making the array to make palette. It's the same length than the array from the "UV" library
	for i := 1; i < 48; i++ {
		components = append(components, component.Dup())
    }
	//getting palette
	anthaPalette := MakeAnthaPalette(palette,components)

	//getting plate
	plate := factory.GetPlateByType("greiner384")

	t.Log(plate.ID)

	//testing function
	anthaImg, resizedImg := MakeAnthaImg(imgBase, anthaPalette, plate)

	t.Log(anthaImg)
	t.Log(resizedImg)

}

func TestMakeAnthaPalette(t *testing.T) {

	//getting palette
	palette, err := SelectLibrary("UV")
	t.Log(len(palette))
	t.Log(err)

	//initiating component
	var components []*wtype.LHComponent
	component := factory.GetComponentByType("Gluc")

	//making the array to test. It's the same length than the array from the "UV" library
	for i := 1; i < 48; i++ {
		components = append(components, component.Dup())
    }
	t.Log(len(components))

	//running the function
	anthaPalette := MakeAnthaPalette(palette,components)

	t.Log(anthaPalette.AnthaColors[0].Component)
}

func TestSelectLivingColorLibrary(t *testing.T){
	livingColors, err := SelectLivingColorLibrary("ProteinPaintBox")
	t.Log(err)

	t.Log(livingColors)
}

func TestSelectLivingColor(t *testing.T){
	livingColor, err := SelectLivingColor("UVDasherGFP")
	t.Log(err)

	t.Log(livingColor)
}

func TestMakeLivingImg(t *testing.T) {

	//downloading image for the test
	imgFile , err := download.File("http://orig08.deviantart.net/a19f/f/2008/117/6/7/8_bit_mario_by_superjerk.jpg", "Downloaded file")
	if err != nil{
		t.Error(err)
	}

	//opening image
	imgBase, err := OpenFile(imgFile)
	if err != nil{
		t.Error(err)
	}

	//initiating components
	var components []*wtype.LHComponent
	component := factory.GetComponentByType("Gluc")
	//making the array to make palette. It's the same length than the array from the "ProteinPaintbox" library
	for i := 1; i < 3; i++ {
		components = append(components, component.Dup())
    }

	//Selecting livingPalette
	selectedPalette, err := SelectLivingColorLibrary("ProteinPaintBox")
	t.Log(err)

	//Making palette
	livingPalette := MakeLivingPalette(selectedPalette,components)

	//getting plate
	plate := factory.GetPlateByType("greiner384")

	t.Log(plate.ID)

	//testing function
	anthaImg, resizedImg := MakeLivingImg(imgBase, livingPalette, plate)

	t.Log(anthaImg)
	t.Log(resizedImg)
}

func TestMakeLivingGIF(t *testing.T) {

	//------------------------------------------------
	//Making antha image
	//------------------------------------------------

	//downloading image for the test
	imgFile , err := download.File("http://orig08.deviantart.net/a19f/f/2008/117/6/7/8_bit_mario_by_superjerk.jpg", "Downloaded file")
	if err != nil{
		t.Error(err)
	}

	//opening image
	imgBase, err := OpenFile(imgFile)
	if err != nil{
		t.Error(err)
	}

	//initiating components
	var components []*wtype.LHComponent
	component := factory.GetComponentByType("Gluc")
	//making the array to make palette. It's the same length than the array from the "ProteinPaintbox" library
	for i := 1; i < 3; i++ {
		components = append(components, component.Dup())
    }

	//Selecting livingPalette
	selectedPalette, err := SelectLivingColorLibrary("ProteinPaintBox")
	t.Log(err)

	//Making palette
	livingPalette := MakeLivingPalette(selectedPalette,components)

	//getting plate
	plate := factory.GetPlateByType("greiner384")

	t.Log(plate.ID)

	//generating images. We only use 2 since that's what we use for our construct
	anthaImg1, _ := MakeLivingImg(imgBase, livingPalette, plate)
	anthaImg2, _ := MakeLivingImg(imgBase, livingPalette, plate)

	//Merge them
	var anthaImgs []LivingImg
	anthaImgs = append(anthaImgs,*anthaImg1, *anthaImg2)

	//------------------------------------------------
	//Testing GIF functions
	//------------------------------------------------

	livingGIF := MakeLivingGIF(anthaImgs)

	t.Log(len(livingGIF.Frames))

}