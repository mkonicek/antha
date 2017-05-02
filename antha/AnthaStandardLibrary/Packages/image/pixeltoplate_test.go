package image

import (
	"testing"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"io/ioutil"
	goimage "image"
)

func check(e error) {
    if e != nil {
        panic(e)
    }
}

func PathToImg (filePath string)(*goimage.NRGBA, error){

	var file wtype.File

	//getting bytes
	dat, err := ioutil.ReadFile(filePath)
	check(err)

	//Making wtype.File object
	file.WriteAll(dat)

	//Opening with the function
	img, err := OpenFile(file)

	return img, err
}

//opening all images
func TestOpenFile(t *testing.T) {

	var jpgFile wtype.File
	var pngFile wtype.File
	var gifFile wtype.File

	// example images
	jpgPath := "/home/cachemoi/gocode/src/github.com/antha-lang/elements/an/GIF/DataGIF/img/F1.jpg"
	pngPath := "/home/cachemoi/gocode/src/github.com/antha-lang/elements/an/GIF/DataGIF/img/F1.png"
	gifPath := "/home/cachemoi/gocode/src/github.com/antha-lang/elements/an/GIF/DataGIF/img/nyanCat.gif"

	//getting bytes
	jpgDat, err := ioutil.ReadFile(jpgPath)
	check(err)
	pngDat, err := ioutil.ReadFile(pngPath)
	check(err)
	gifDat, err := ioutil.ReadFile(gifPath)
	check(err)

	//Making wtype.File object
	jpgFile.WriteAll(jpgDat)
	pngFile.WriteAll(pngDat)
	gifFile.WriteAll(gifDat)

	//Opening with the function
	OpenFile(jpgFile)
	t.Log("opened JPEG")
	OpenFile(pngFile)
	t.Log("opened PNG")
	OpenFile(gifFile)
	t.Log("Opened GIF")
}

func TestPosterize(t *testing.T) {

	postImg, _ := Posterize("/home/cachemoi/gocode/src/github.com/antha-lang/elements/an/GIF/DataGIF/img/F1.jpg",2)

	if postImg != nil {
		t.Log("image posterized")
	}else {
		t.Error("posterize() returned nil")
	}

}