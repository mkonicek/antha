package image

import (
	"testing"

	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/laboratory"
	"github.com/antha-lang/antha/laboratory/testlab"
)

func TestSelectLibrary(t *testing.T) {
	SelectLibrary("UV")
}

func TestSelectColors(t *testing.T) {
	SelectColor("JuniperGFP")
}

func TestMakeAnthaImg(t *testing.T) {
	testlab.WithTestLab(t, "testdata", &testlab.TestElementCallbacks{
		Steps: func(lab *laboratory.Laboratory) error {
			imgFile := wtype.NewFile("invader.png").AsInput()
			//opening image
			imgBase, err := OpenFile(lab, imgFile)
			if err != nil {
				return err
			}

			palette := SelectLibrary("UV")

			//initiating components
			var components []*wtype.Liquid
			component, err := lab.Inventory.Components.NewComponent("Gluc")
			if err != nil {
				return err
			}

			//making the array to make palette. It's the same length than the array from the "UV" library
			for i := 1; i <= len(palette); i++ {
				components = append(components, component.Dup(lab.IDGenerator))
			}
			//getting palette
			anthaPalette := MakeAnthaPalette(palette, components)

			//getting plate
			plate, err := lab.Inventory.Plates.NewPlate("greiner384")
			if err != nil {
				return err
			}

			//testing function
			MakeAnthaImg(imgBase, anthaPalette, plate)
			return nil
		},
	})
}

func TestMakeAnthaPalette(t *testing.T) {
	testlab.WithTestLab(t, "testdata", &testlab.TestElementCallbacks{
		Steps: func(lab *laboratory.Laboratory) error {
			//getting palette
			palette := SelectLibrary("UV")

			//initiating component
			var components []*wtype.Liquid
			component, err := lab.Inventory.Components.NewComponent("Gluc")
			if err != nil {
				return err
			}

			//making the array to test. It's the same length than the array from the "UV" library
			for i := 1; i <= len(palette); i++ {
				components = append(components, component.Dup(lab.IDGenerator))
			}

			//running the function
			MakeAnthaPalette(palette, components)
			return nil
		},
	})
}

func TestSelectLivingColorLibrary(t *testing.T) {
	testlab.WithTestLab(t, "testdata", &testlab.TestElementCallbacks{
		Steps: func(lab *laboratory.Laboratory) error {
			SelectLivingColorLibrary(lab, "ProteinPaintBox")
			return nil
		},
	})
}

func TestSelectLivingColor(t *testing.T) {
	testlab.WithTestLab(t, "testdata", &testlab.TestElementCallbacks{
		Steps: func(lab *laboratory.Laboratory) error {
			SelectLivingColor(lab, "DasherGFP")
			return nil
		},
	})
}

func TestMakeLivingImg(t *testing.T) {
	testlab.WithTestLab(t, "testdata", &testlab.TestElementCallbacks{
		Steps: func(lab *laboratory.Laboratory) error {
			imgFile := wtype.NewFile("invader.png").AsInput()
			//opening image
			imgBase, err := OpenFile(lab, imgFile)
			if err != nil {
				return err
			}

			//initiating components
			var components []*wtype.Liquid
			component, err := lab.Inventory.Components.NewComponent("Gluc")
			if err != nil {
				return nil
			}

			//Selecting livingPalette
			selectedPalette := SelectLivingColorLibrary(lab, "ProteinPaintBox")

			//making the array to make palette. It's the same length than the array from the "ProteinPaintbox" library
			for i := 1; i <= len(selectedPalette.LivingColors); i++ {
				components = append(components, component.Dup(lab.IDGenerator))
			}

			//Making palette
			livingPalette := MakeLivingPalette(selectedPalette, components)

			//getting plate
			plate, err := lab.Inventory.Plates.NewPlate("greiner384")
			if err != nil {
				return nil
			}

			//testing function
			MakeLivingImg(imgBase, livingPalette, plate)
			return nil
		},
	})
}

func TestMakeLivingGIF(t *testing.T) {
	testlab.WithTestLab(t, "testdata", &testlab.TestElementCallbacks{
		Steps: func(lab *laboratory.Laboratory) error {
			imgFile := wtype.NewFile("invader.png").AsInput()
			//opening image
			imgBase, err := OpenFile(lab, imgFile)
			if err != nil {
				return err
			}

			//initiating components
			var components []*wtype.Liquid
			component, err := lab.Inventory.Components.NewComponent("Gluc")
			if err != nil {
				return nil
			}

			//Selecting livingPalette
			selectedPalette := SelectLivingColorLibrary(lab, "ProteinPaintBox")

			//making the array to make palette. It's the same length than the array from the "ProteinPaintbox" library
			for i := 1; i <= len(selectedPalette.LivingColors); i++ {
				components = append(components, component.Dup(lab.IDGenerator))
			}

			//Making palette
			livingPalette := MakeLivingPalette(selectedPalette, components)

			//getting plate
			plate, err := lab.Inventory.Plates.NewPlate("greiner384")
			if err != nil {
				return err
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
			return nil
		},
	})
}

func TestAnthaPrintWorkflow(t *testing.T) {
	testlab.WithTestLab(t, "testdata", &testlab.TestElementCallbacks{
		Steps: func(lab *laboratory.Laboratory) error {
			imgFile := wtype.NewFile("invader.png").AsInput()
			//opening image
			imgBase, err := OpenFile(lab, imgFile)
			if err != nil {
				return err
			}

			palette := SelectLibrary("UV")

			//initiating components
			var components []*wtype.Liquid
			component, err := lab.Inventory.Components.NewComponent("Gluc")
			if err != nil {
				return err
			}

			//making the array to make palette. It's the same length than the array from the "UV" library
			for i := 1; i <= len(palette); i++ {
				components = append(components, component.Dup(lab.IDGenerator))
			}
			//getting palette
			anthaPalette := MakeAnthaPalette(palette, components)

			//getting plate
			plate, err := lab.Inventory.Plates.NewPlate("greiner384")
			if err != nil {
				return err
			}

			//testing function
			MakeAnthaImg(imgBase, anthaPalette, plate)
			return nil
		},
	})
}

func TestParseGIF(t *testing.T) {
	testlab.WithTestLab(t, "testdata", &testlab.TestElementCallbacks{
		Steps: func(lab *laboratory.Laboratory) error {
			imgFile := wtype.NewFile("potato.gif").AsInput()
			//opening GIF
			gif, err := OpenGIF(lab, imgFile)
			if err != nil {
				return err
			}

			if _, err := ParseGIF(gif, []int{1, 6}); err != nil {
				return err
			}
			return nil
		},
	})
}

func TestGetState(t *testing.T) {
	testlab.WithTestLab(t, "testdata", &testlab.TestElementCallbacks{
		Steps: func(lab *laboratory.Laboratory) error {
			imgFile := wtype.NewFile("potato.gif").AsInput()
			//opening GIF
			gif, err := OpenGIF(lab, imgFile)
			if err != nil {
				return err
			}

			imgs, err := ParseGIF(gif, []int{1, 2})
			if err != nil {
				return err
			}

			lp := SelectLivingColorLibrary(lab, "ProteinPaintBox")
			lip, _ := lab.Inventory.Plates.NewPlate("greiner384")

			li1, _ := MakeLivingImg(imgs[0], &lp, lip)
			li2, _ := MakeLivingImg(imgs[1], &lp, lip)

			var lia []LivingImg

			lia = append(lia, *li1, *li2)
			//------------------------------------------------
			//Testing GIF functions
			//------------------------------------------------

			LivingGIF := MakeLivingGIF(lia)

			LivingGIF.GetStates()
			return nil
		},
	})
}
