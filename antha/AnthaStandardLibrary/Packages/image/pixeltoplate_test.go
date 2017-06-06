package image

import (
	"testing"
	"reflect"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"io/ioutil"
	"github.com/antha-lang/antha/microArch/factory"
)

func TestSelectLibrary (t *testing.T) {

	palette := SelectLibrary("UV")
	t.Log(palette)
	t.Log(reflect.TypeOf(palette))

}

func TestSelectColors(t *testing.T) {

	palette :=  SelectColor("JuniperGFP")
	t.Log(palette)
}

func TestMakeAnthaImg(t *testing.T) {

	//opening test image
	var testFile wtype.File
	palette := SelectLibrary("UV")

	dat, err := ioutil.ReadFile("/home/cachemoi/gocode/src/github.com/cachemoi/playing/img/F1.png")
	if err != nil{
		t.Error(err)
	}
	testFile.WriteAll(dat)

	//reading image file
	imgBase, err := OpenFile(testFile)
	if err != nil{
		t.Error(err)
	}

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
	anthaImg := MakeAnthaImg(imgBase, anthaPalette, plate)

	t.Log(anthaImg.Pix[0].Location)
}

func TestMakeAnthaPalette(t *testing.T) {

	//getting palette
	palette := SelectLibrary("UV")
	t.Log(len(palette))

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