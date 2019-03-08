package image

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/gif"
	"path"
	"strings"

	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wutil"
	"github.com/antha-lang/antha/laboratory"
	"github.com/disintegration/imaging"
)

// OpenFile takes a wtype.File and returns its contents as image.NRGBA
func OpenFile(lab *laboratory.Laboratory, file *wtype.File) (*image.NRGBA, error) {
	data, err := lab.FileManager.ReadAll(file)
	if err != nil {
		return nil, err
	}

	buf := bytes.NewReader(data)

	img, err := imaging.Decode(buf)
	if err != nil {
		return nil, err
	}

	return imaging.Clone(img), nil
}

// OpenGIF takes a wtype.File and returns its contents as a gif.GIF
func OpenGIF(lab *laboratory.Laboratory, file *wtype.File) (*gif.GIF, error) {
	data, err := lab.FileManager.ReadAll(file)
	if err != nil {
		return nil, err
	}

	reader := bytes.NewReader(data)

	img, err := gif.DecodeAll(reader)
	if err != nil {
		return nil, err
	}

	return img, nil
}

func maxIntInSlice(xs []int) (r int) {
	for _, x := range xs {
		if x > r {
			r = x
		}
	}
	return
}

// ParseGIF extracts frames from a GIF object as an array of images
func ParseGIF(source *gif.GIF, frameNum []int) (imgs []*image.NRGBA, err error) {
	largest := maxIntInSlice(frameNum)

	if n := len(source.Image); largest >= n {
		return nil, fmt.Errorf("frame number %d exceeds number of frames %d", largest, n)
	}

	for _, num := range frameNum {
		imgs = append(imgs, toNRGBA(source.Image[num]))
	}

	return imgs, nil
}

// ColourToCMYK converts a color to color.CYMK
func ColourToCMYK(colour color.Color) (cmyk color.CMYK) {
	r, g, b, _ := colour.RGBA()
	cmyk.C, cmyk.M, cmyk.Y, cmyk.K = color.RGBToCMYK(uint8(r), uint8(g), uint8(b))
	return
}

// ColourToGrayscale converts a color to color.Gray
func ColourToGrayscale(colour color.Color) (gray color.Gray) {
	r, g, b, _ := colour.RGBA()
	gray.Y = uint8((0.2126 * float64(r)) + (0.7152 * float64(g)) + (0.0722 * float64(b)))
	return
}

type ImageFormat uint8

const (
	PNG ImageFormat = iota
	JPEG
	TIFF
	GIF
	BMP
	_maxImageFormat
)

func ImageFormatFromPath(name string) (ImageFormat, error) {
	ext := strings.ToLower(path.Ext(name))
	switch ext {
	case ".png":
		return PNG, nil
	case ".jpg", ".jpeg":
		return JPEG, nil
	case ".tif", ".tiff":
		return TIFF, nil
	case ".gif":
		return GIF, nil
	case ".bmp":
		return BMP, nil
	default:
		return _maxImageFormat, fmt.Errorf("Unrecognised image format extension: %v", ext)
	}
}

// Export exports an image to file.
func Export(lab *laboratory.Laboratory, img image.Image, format ImageFormat) (*wtype.File, error) {
	var imageFormat imaging.Format
	switch format {
	case PNG:
		imageFormat = imaging.PNG
	case JPEG:
		imageFormat = imaging.JPEG
	case TIFF:
		imageFormat = imaging.TIFF
	case GIF:
		imageFormat = imaging.GIF
	case BMP:
		imageFormat = imaging.BMP
	default:
		imageFormat = imaging.PNG
	}

	var buf bytes.Buffer
	if err := imaging.Encode(&buf, img, imageFormat); err != nil {
		return nil, err
	}

	return lab.FileManager.WriteAll(buf.Bytes())
}

// Posterize posterizes an image. This refers to changing an image to use only
// a small number of different tones.
func Posterize(img image.Image, levels int) (*image.NRGBA, error) {

	if levels == 1 {
		return nil, errors.New("Cannot posterize with only one level")
	}

	numberOfAreas := 256 / (levels)
	numberOfValues := 255 / (levels - 1)

	posterized := imaging.Clone(img)

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
				rfloat := (float64(r/256) / float64(numberOfAreas))

				rinttemp, err := wutil.RoundDown(rfloat)
				if err != nil {
					return nil, err
				}
				rnew = float64(rinttemp) * float64(numberOfValues)
			}
			if g == 0 {
				gnew = 0
			} else {
				gfloat := (float64(g/256) / float64(numberOfAreas))

				ginttemp, err := wutil.RoundDown(gfloat)
				if err != nil {
					return nil, err
				}
				gnew = float64(ginttemp) * float64(numberOfValues)
			}
			if b == 0 {
				bnew = 0
			} else {
				bfloat := (float64(b/256) / float64(numberOfAreas))

				binttemp, err := wutil.RoundDown(bfloat)
				if err != nil {
					return nil, err
				}
				bnew = float64(binttemp) * float64(numberOfValues)
			}

			var newColor color.NRGBA
			newColor.A = uint8(a)

			rint, err := wutil.RoundDown(rnew)

			if err != nil {
				return nil, err
			}
			newColor.R = uint8(rint)
			gint, err := wutil.RoundDown(gnew)
			if err != nil {
				return nil, err
			}
			newColor.G = uint8(gint)
			bint, err := wutil.RoundDown(bnew)
			if err != nil {
				return nil, err
			}
			newColor.B = uint8(bint)

			posterized.Set(x, y, newColor)
		}
	}

	return posterized, nil
}
