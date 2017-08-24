// Copyright 2011 <chaishushan@gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pick

import (
	"context"
	"fmt"
	"os"
	//"path"
	//"runtime"
	"encoding/json"
	"image"
	"image/color"
	"image/draw"
	"math"
	//"path/filepath"
	"io/ioutil"
	"strconv"
	"strings"

	"code.google.com/p/draw2d/draw2d"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/antha/anthalib/wutil"
	"github.com/antha-lang/antha/inventory"
	//"code.google.com/p/draw2d/draw2d"
	"github.com/disintegration/imaging"
	//"github.com/hybridgroup/go-opencv/opencv"
	"github.com/lazywei/go-opencv/opencv"
	//"../opencv" // can be used in forks, comment in real application
)

func PickAndExportCSV(ctx context.Context, imagefile string, exportFileName string, plate *wtype.LHPlate, numbertopick int, setplateperimeterfirst bool, rotate bool) (wells []string, err error) {
	var topleft opencv.Point
	var bottomright opencv.Point
	var topright opencv.Point
	var well string
	var rotated bool
	var rotatedimg draw.Image
	var newname string
	counter := 0
	//_, currentfile, _, _ := runtime.Caller(0)
	filename := imagefile //:= path.Join(path.Dir(currentfile), imagefile)
	wells = make([]string, 0)
	if len(os.Args) == 2 {
		filename = os.Args[1]
	}

	img0 := opencv.LoadImage(filename)
	if img0 == nil {
		panic("LoadImage fail")
	}

	// resize image if bigger than full screen
	if img0.Height() > 801 {
		img0.Release()
		imgtoresize, err := imaging.Open(filename)
		if err != nil {
			panic(err)
		}
		resizedimage := imaging.Resize(imgtoresize, 0, 800, imaging.CatmullRom)
		splitfilename := strings.Split(filename, `.`)

		newname = splitfilename[0] + "_resized" + `.` + splitfilename[1]

		err = imaging.Save(resizedimage, newname)
		if err != nil {
			panic(err)
		}
		img0 = opencv.LoadImage(newname)
		filename = newname
		if img0 == nil {
			panic("LoadImage fail")
		}
	}

	defer img0.Release()

	/*fmt.Print("Hot keys: \n",
		"\tESC - quit the program\n",
		"\tr - restore the original image\n",
		"\ti or ENTER - run inpainting algorithm\n",
		"\t\t(before running it, paint something on the image)\n",
	)*/
	fmt.Println("Click on top left corner of plate")

	img := img0.Clone()
	inpainted := img0.Clone()
	inpaint_mask := opencv.CreateImage(img0.Width(), img0.Height(), 8, 1)

	opencv.Zero(inpaint_mask)
	//opencv.Zero( inpainted )

	win := opencv.NewWindow("image")
	defer win.Destroy()

	prev_pt := opencv.Point{-1, -1}
	win.SetMouseCallback(func(event, x, y, flags int, param ...interface{}) {
		if img == nil {
			panic("img == nil")
		}

		if event == opencv.CV_EVENT_LBUTTONUP ||
			(flags&opencv.CV_EVENT_FLAG_LBUTTON) == 0 {
			prev_pt = opencv.Point{-1, -1}
		} else if event == opencv.CV_EVENT_LBUTTONDOWN {

			rgb := opencv.ScalarAll(255.0)
			counter = counter + 1

			if rotate && rotated == false {
				imgtorotate, err := imaging.Open(filename)

				if err != nil {
					panic(err)
				}
				if counter == 1 {
					fmt.Println("Rotating image")
					fmt.Println("topleft point:", x, y)
					topleft = opencv.Point{x, y}
				} else if counter == 2 {
					topright = opencv.Point{x, y}
					fmt.Println("topright point:", x, y)
					opposite := float64(topleft.Y) - float64(topright.Y)
					adjacent := float64(topright.X) - float64(topleft.X)
					tantheta := opposite / adjacent
					fmt.Println("adjacent:", adjacent)
					fmt.Println("opposite:", opposite)
					fmt.Println("costheta:", tantheta)
					thetainrad := math.Atan(tantheta)
					fmt.Println("thetainrad:", thetainrad)
					degrees := (180 / math.Pi) * thetainrad
					fmt.Println("degrees:", degrees)
					tr := draw2d.NewRotationMatrix(thetainrad)

					//
					ar := imgtorotate.Bounds()
					w, h, _ := ar.Dx(), ar.Dy(), 30.0
					rotatedimg = image.NewRGBA(image.Rect(0, 0, w, h))
					draw.Draw(rotatedimg, ar, imgtorotate, ar.Min, draw.Src)

					draw2d.DrawImage(imgtorotate, rotatedimg, tr, draw.Src, draw2d.LinearFilter)

					// open new window
					splitfilename := strings.Split(filename, `.`)
					newname = splitfilename[0] + "_rotated" + `.` + splitfilename[1]

					//err = imaging.Save(rotatedimg, newname)
					err = imaging.Save(rotatedimg, newname)

					if err != nil {
						panic(err)
					}
					/*
						wells = PickAgain(newname, numbertopick, setplateperimeterfirst, false)
						return*/
					/*img0 = opencv.LoadImage(newname)
					filename = newname
					if img0 == nil {
						panic("LoadImage fail")
					}
					defer img0.Release()

					fmt.Print("Hot keys: \n",
						"\tESC - quit the program\n",
						"\tr - restore the original image\n",
						"\ti or ENTER - run inpainting algorithm\n",
						"\t\t(before running it, paint something on the image)\n",
					)

					//img := img0.Clone()
					//inpainted := img0.Clone()
					//inpaint_mask := opencv.CreateImage(img0.Width(), img0.Height(), 8, 1)

					//opencv.Zero(inpaint_mask)
					//opencv.Zero( inpainted )

					win := opencv.NewWindow("rotated")
					defer win.Destroy()
					// change status
					rotated = true
					fmt.Println("Now choose top left of plate")

					*/
				}
			} else if rotate && setplateperimeterfirst && counter == 3 && rotated || setplateperimeterfirst && counter == 1 && rotate == false {
				topleft = opencv.Point{x, y}
				fmt.Println("Topleft chosen")
				fmt.Println("Now choose bottom right of plate")

			} else if rotate && setplateperimeterfirst && counter == 4 && rotated || setplateperimeterfirst && counter == 2 && rotate == false {
				bottomright = opencv.Point{x, y}
				fmt.Println("plate boundaries", x, y)
				opencv.Rectangle(img, topleft, bottomright, rgb, 1, 8, 0)
				win.ShowImage(img)
			} else {

				prev_pt = opencv.Point{x, y}
				fmt.Println("actual pixels:", x, y)
				fmt.Println("imagesize:", img.ImageSize())
				fmt.Println("width:", img.Width())
				fmt.Println("height:", img.Height())
				fmt.Println("Well:", PixelstoWellPosition(x, y, plate, img))
				fmt.Println("colony count:", counter)

				if setplateperimeterfirst {
					well = PixelstoWellPositionFromRectangle(x, y, plate, topleft, bottomright)
				} else if setplateperimeterfirst == false {
					well = PixelstoWellPosition(x, y, plate, img)
				}

				wells = append(wells, well)
				fmt.Println(wells)

				// draw circle on click

				rgb = opencv.ScalarAll(255.0)
				opencv.Circle(img, prev_pt, 1, rgb, 5, 8, 0)

				win.ShowImage(img)

				if rotate && setplateperimeterfirst && counter == numbertopick+4 {
					//return
					splitfilename := strings.Split(filename, `.`)

					newname := splitfilename[0] + "_picked" + `.` + splitfilename[1]
					opencv.SaveImage(newname, img, 0)

					err = ExportCSV(ctx, exportFileName, plate, "ColonyPlate", wells, "colonynumber", "colony")
					os.Exit(0)
					return
				}
				if rotate == false && setplateperimeterfirst && counter == numbertopick+2 {
					//return
					splitfilename := strings.Split(filename, `.`)

					newname := splitfilename[0] + "_picked" + `.` + splitfilename[1]
					opencv.SaveImage(newname, img, 0)
					err = ExportCSV(ctx, exportFileName, plate, "ColonyPlate", wells, "colonynumber", "colony")
					os.Exit(0)
					return
				}
				if rotate == false && setplateperimeterfirst == false && counter == numbertopick {
					//return

					err = ExportCSV(ctx, exportFileName, plate, "ColonyPlate", wells, "colonynumber", "colony")
					os.Exit(0)
					return
				}
			}
		} else if event == opencv.CV_EVENT_MOUSEMOVE &&
			(flags&opencv.CV_EVENT_FLAG_LBUTTON) != 0 {
			pt := opencv.Point{x, y}
			if prev_pt.X < 0 {
				prev_pt = pt
			}

			rgb := opencv.ScalarAll(255.0)

			opencv.Rectangle(img, prev_pt, pt, rgb, 1, 8, 0)
			//opencv.Line(inpaint_mask, prev_pt, pt, rgb, 5, 8, 0)
			//opencv.Line(img, prev_pt, pt, rgb, 5, 8, 0)
			prev_pt = pt

			win.ShowImage(img)
		}
	})
	win.ShowImage(img)
	opencv.WaitKey(0)

	win2 := opencv.NewWindow("inpainted image")

	defer win2.Destroy()
	win2.ShowImage(inpainted)

	for {
		key := opencv.WaitKey(20)
		if key == 27 {
			//	os.Exit(0)
		} else if key == 'r' {
			opencv.Zero(inpaint_mask)
			opencv.Copy(img0, img, nil)
			win.ShowImage(img)
		} else if key == 'i' || key == '\n' {
			opencv.Inpaint(img, inpaint_mask, inpainted, 3,
				opencv.CV_INPAINT_TELEA,
			)
			win2.ShowImage(inpainted)
		}
	}
	//os.Exit(0)
	/*if rotated {
		wells = PickAgain(newname, numbertopick, setplateperimeterfirst, false)
	}*/
	return
}

func PickAndExportWelltoColourJSON(ctx context.Context, imagefile string, exportFileName string, plate *wtype.LHPlate, numbertopick int, setplateperimeterfirst bool, rotate bool) (wells []string, welltoColourmap map[string]color.Color, err error) {
	var topleft opencv.Point
	var bottomright opencv.Point
	var topright opencv.Point
	var well string
	var rotated bool
	var rotatedimg draw.Image
	var newname string

	welltoColourmap = make(map[string]color.Color)
	counter := 0
	//_, currentfile, _, _ := runtime.Caller(0)
	filename := imagefile //:= path.Join(path.Dir(currentfile), imagefile)
	wells = make([]string, 0)
	if len(os.Args) == 2 {
		filename = os.Args[1]
	}

	img0 := opencv.LoadImage(filename)
	if img0 == nil {
		panic("LoadImage fail for " + filename)
	}

	// resize image if bigger than full screen
	if img0.Height() > 801 {
		img0.Release()
		imgtoresize, err := imaging.Open(filename)
		if err != nil {
			panic(err)
		}
		resizedimage := imaging.Resize(imgtoresize, 0, 800, imaging.CatmullRom)
		splitfilename := strings.Split(filename, `.`)

		newname = splitfilename[0] + "_resized" + `.` + splitfilename[1]

		err = imaging.Save(resizedimage, newname)
		if err != nil {
			panic(err)
		}
		img0 = opencv.LoadImage(newname)
		filename = newname
		if img0 == nil {
			panic("LoadImage fail")
		}
	}

	defer img0.Release()

	/*fmt.Print("Hot keys: \n",
		"\tESC - quit the program\n",
		"\tr - restore the original image\n",
		"\ti or ENTER - run inpainting algorithm\n",
		"\t\t(before running it, paint something on the image)\n",
	)*/
	fmt.Println("Click on top left corner of plate")

	img := img0.Clone()
	inpainted := img0.Clone()
	inpaint_mask := opencv.CreateImage(img0.Width(), img0.Height(), 8, 1)

	/*imgforcolour, err := imaging.Open(filename)
	if err != nil {
		panic(err)
	}
	nrgba := toNRGBA(imgforcolour)*/
	opencv.Zero(inpaint_mask)
	//opencv.Zero( inpainted )

	win := opencv.NewWindow("image")

	fmt.Println("Oooh I'm through that")

	defer win.Destroy()

	prev_pt := opencv.Point{-1, -1}
	win.SetMouseCallback(func(event, x, y, flags int, param ...interface{}) {
		if img == nil {
			panic("img == nil")
		}

		fmt.Println("hehehaheha ah ooh that tickles stop!")

		if event == opencv.CV_EVENT_LBUTTONUP ||
			(flags&opencv.CV_EVENT_FLAG_LBUTTON) == 0 {
			prev_pt = opencv.Point{-1, -1}
		} else if event == opencv.CV_EVENT_LBUTTONDOWN {

			rgb := opencv.ScalarAll(255.0)
			counter = counter + 1

			if rotate && rotated == false {
				imgtorotate, err := imaging.Open(filename)

				if err != nil {
					panic(err)
				}
				if counter == 1 {
					fmt.Println("Rotating image")
					fmt.Println("topleft point:", x, y)
					topleft = opencv.Point{x, y}
				} else if counter == 2 {
					topright = opencv.Point{x, y}
					fmt.Println("topright point:", x, y)
					opposite := float64(topleft.Y) - float64(topright.Y)
					adjacent := float64(topright.X) - float64(topleft.X)
					tantheta := opposite / adjacent
					fmt.Println("adjacent:", adjacent)
					fmt.Println("opposite:", opposite)
					fmt.Println("costheta:", tantheta)
					thetainrad := math.Atan(tantheta)
					fmt.Println("thetainrad:", thetainrad)
					degrees := (180 / math.Pi) * thetainrad
					fmt.Println("degrees:", degrees)
					tr := draw2d.NewRotationMatrix(thetainrad)

					//
					ar := imgtorotate.Bounds()
					w, h, _ := ar.Dx(), ar.Dy(), 30.0
					rotatedimg = image.NewRGBA(image.Rect(0, 0, w, h))
					draw.Draw(rotatedimg, ar, imgtorotate, ar.Min, draw.Src)

					draw2d.DrawImage(imgtorotate, rotatedimg, tr, draw.Src, draw2d.LinearFilter)

					// open new window
					splitfilename := strings.Split(filename, `.`)
					newname = splitfilename[0] + "_rotated" + `.` + splitfilename[1]

					//err = imaging.Save(rotatedimg, newname)
					err = imaging.Save(rotatedimg, newname)

					if err != nil {
						panic(err)
					}

					filename = newname

					/*
						wells = PickAgain(newname, numbertopick, setplateperimeterfirst, false)
						return*/
					/*img0 = opencv.LoadImage(newname)
					filename = newname
					if img0 == nil {
						panic("LoadImage fail")
					}
					defer img0.Release()

					fmt.Print("Hot keys: \n",
						"\tESC - quit the program\n",
						"\tr - restore the original image\n",
						"\ti or ENTER - run inpainting algorithm\n",
						"\t\t(before running it, paint something on the image)\n",
					)

					//img := img0.Clone()
					//inpainted := img0.Clone()
					//inpaint_mask := opencv.CreateImage(img0.Width(), img0.Height(), 8, 1)

					//opencv.Zero(inpaint_mask)
					//opencv.Zero( inpainted )

					win := opencv.NewWindow("rotated")
					defer win.Destroy()
					// change status
					rotated = true
					fmt.Println("Now choose top left of plate")

					*/
				}
			} else if rotate && setplateperimeterfirst && counter == 3 && rotated || setplateperimeterfirst && counter == 1 && rotate == false {
				topleft = opencv.Point{x, y}
				fmt.Println("Topleft chosen")
				fmt.Println("Now choose bottom right of plate")

			} else if rotate && setplateperimeterfirst && counter == 4 && rotated || setplateperimeterfirst && counter == 2 && rotate == false {
				bottomright = opencv.Point{x, y}
				fmt.Println("plate boundaries", x, y)
				opencv.Rectangle(img, topleft, bottomright, rgb, 1, 8, 0)
				win.ShowImage(img)
			} else {

				prev_pt = opencv.Point{x, y}
				fmt.Println("actual pixels:", x, y)
				fmt.Println("imagesize:", img.ImageSize())
				fmt.Println("width:", img.Width())
				fmt.Println("height:", img.Height())
				fmt.Println("Well:", PixelstoWellPosition(x, y, plate, img))
				fmt.Println("colony count:", counter)

				if setplateperimeterfirst {
					well = PixelstoWellPositionFromRectangle(x, y, plate, topleft, bottomright)
				} else if setplateperimeterfirst == false {
					well = PixelstoWellPosition(x, y, plate, img)
				}

				wells = append(wells, well)
				imgforcolour, err := imaging.Open(filename)
				if err != nil {
					panic(err)
				}
				nrgba := toNRGBA(imgforcolour)
				welltoColourmap[plate.ID+"_"+well] = nrgba.NRGBAAt(x, y)

				fmt.Println(wells)

				// draw circle on click

				rgb = opencv.ScalarAll(255.0)
				opencv.Circle(img, prev_pt, 1, rgb, 5, 8, 0)

				win.ShowImage(img)

				if rotate && setplateperimeterfirst && counter == numbertopick+4 {
					//return
					splitfilename := strings.Split(filename, `.`)

					newname := splitfilename[0] + "_picked" + `.` + splitfilename[1]
					opencv.SaveImage(newname, img, 0)
					bytes, err := json.Marshal(welltoColourmap)

					if err != nil {
						panic(err.Error())
					}

					ioutil.WriteFile(exportFileName+".json", bytes, 0644)

					err = ExportCSV(ctx, exportFileName+".csv", plate, "ColonyPlate", wells, "colonynumber", "colony")
					os.Exit(0)
					return
				}
				if rotate == false && setplateperimeterfirst && counter == numbertopick+2 {
					//return
					splitfilename := strings.Split(filename, `.`)

					newname := splitfilename[0] + "_picked" + `.` + splitfilename[1]
					opencv.SaveImage(newname, img, 0)

					bytes, err := json.Marshal(welltoColourmap)

					if err != nil {
						panic(err.Error())
					}

					ioutil.WriteFile(exportFileName+".json", bytes, 0644)

					err = ExportCSV(ctx, exportFileName+".csv", plate, "ColonyPlate", wells, "colonynumber", "colony")
					os.Exit(0)
					return
				}
				if rotate == false && setplateperimeterfirst == false && counter == numbertopick {
					//return

					bytes, err := json.Marshal(welltoColourmap)

					if err != nil {
						panic(err.Error())
					}

					ioutil.WriteFile(exportFileName+".json", bytes, 0644)

					err = ExportCSV(ctx, exportFileName+".csv", plate, "ColonyPlate", wells, "colonynumber", "colony")
					os.Exit(0)
					return
				}
			}
		} else if event == opencv.CV_EVENT_MOUSEMOVE &&
			(flags&opencv.CV_EVENT_FLAG_LBUTTON) != 0 {
			pt := opencv.Point{x, y}
			if prev_pt.X < 0 {
				prev_pt = pt
			}

			rgb := opencv.ScalarAll(255.0)

			opencv.Rectangle(img, prev_pt, pt, rgb, 1, 8, 0)
			//opencv.Line(inpaint_mask, prev_pt, pt, rgb, 5, 8, 0)
			//opencv.Line(img, prev_pt, pt, rgb, 5, 8, 0)
			prev_pt = pt

			win.ShowImage(img)
		}
	})

	fmt.Println("i'm out of there")
	win.ShowImage(img)
	opencv.WaitKey(0)

	win2 := opencv.NewWindow("inpainted image")

	defer win2.Destroy()
	win2.ShowImage(inpainted)

	for {
		key := opencv.WaitKey(20)
		if key == 27 {
			//	os.Exit(0)
		} else if key == 'r' {
			opencv.Zero(inpaint_mask)
			opencv.Copy(img0, img, nil)
			win.ShowImage(img)
		} else if key == 'i' || key == '\n' {
			opencv.Inpaint(img, inpaint_mask, inpainted, 3,
				opencv.CV_INPAINT_TELEA,
			)
			win2.ShowImage(inpainted)
		}
	}
	//os.Exit(0)
	/*if rotated {
		wells = PickAgain(newname, numbertopick, setplateperimeterfirst, false)
	}*/
	return
}

func PickAndExportCSVMap(ctx context.Context, imagefile string, exportFileName string, plate *wtype.LHPlate, reactiontonumbertopickmap map[string]int, setplateperimeterfirst bool, rotate bool) (wells []string, err error) {

	var numbertopick int
	var names = make([]string, 0)

	for name, number := range reactiontonumbertopickmap {
		numbertopick = numbertopick + number

		for i := 1; i < number+1; i++ {
			names = append(names, name /*+strconv.Itoa(i)*/)
		}
	}

	var topleft opencv.Point
	var bottomright opencv.Point
	var topright opencv.Point
	var well string
	var rotated bool
	var rotatedimg draw.Image
	var newname string
	counter := 0
	//_, currentfile, _, _ := runtime.Caller(0)
	filename := imagefile //:= path.Join(path.Dir(currentfile), imagefile)
	wells = make([]string, 0)
	if len(os.Args) == 2 {
		filename = os.Args[1]
	}

	img0 := opencv.LoadImage(filename)
	if img0 == nil {
		panic("LoadImage fail")
	}

	// resize image if bigger than full screen
	if img0.Height() > 801 {
		img0.Release()
		imgtoresize, err := imaging.Open(filename)
		if err != nil {
			panic(err)
		}
		resizedimage := imaging.Resize(imgtoresize, 0, 800, imaging.CatmullRom)
		splitfilename := strings.Split(filename, `.`)

		newname = splitfilename[0] + "_resized" + `.` + splitfilename[1]

		err = imaging.Save(resizedimage, newname)
		if err != nil {
			panic(err)
		}
		img0 = opencv.LoadImage(newname)
		filename = newname
		if img0 == nil {
			panic("LoadImage fail")
		}
	}

	defer img0.Release()

	/*fmt.Print("Hot keys: \n",
		"\tESC - quit the program\n",
		"\tr - restore the original image\n",
		"\ti or ENTER - run inpainting algorithm\n",
		"\t\t(before running it, paint something on the image)\n",
	)*/
	fmt.Println("Click on top left corner of plate")

	img := img0.Clone()
	inpainted := img0.Clone()
	inpaint_mask := opencv.CreateImage(img0.Width(), img0.Height(), 8, 1)

	opencv.Zero(inpaint_mask)
	//opencv.Zero( inpainted )

	win := opencv.NewWindow("image")
	defer win.Destroy()

	prev_pt := opencv.Point{-1, -1}
	win.SetMouseCallback(func(event, x, y, flags int, param ...interface{}) {
		if img == nil {
			panic("img == nil")
		}

		if event == opencv.CV_EVENT_LBUTTONUP ||
			(flags&opencv.CV_EVENT_FLAG_LBUTTON) == 0 {
			prev_pt = opencv.Point{-1, -1}
		} else if event == opencv.CV_EVENT_LBUTTONDOWN {

			rgb := opencv.ScalarAll(255.0)
			counter = counter + 1

			if rotate && rotated == false {
				imgtorotate, err := imaging.Open(filename)

				if err != nil {
					panic(err)
				}
				if counter == 1 {
					fmt.Println("Rotating image")
					fmt.Println("topleft point:", x, y)
					topleft = opencv.Point{x, y}
				} else if counter == 2 {
					topright = opencv.Point{x, y}
					fmt.Println("topright point:", x, y)
					opposite := float64(topleft.Y) - float64(topright.Y)
					adjacent := float64(topright.X) - float64(topleft.X)
					tantheta := opposite / adjacent
					fmt.Println("adjacent:", adjacent)
					fmt.Println("opposite:", opposite)
					fmt.Println("costheta:", tantheta)
					thetainrad := math.Atan(tantheta)
					fmt.Println("thetainrad:", thetainrad)
					degrees := (180 / math.Pi) * thetainrad
					fmt.Println("degrees:", degrees)
					tr := draw2d.NewRotationMatrix(thetainrad)

					//
					ar := imgtorotate.Bounds()
					w, h, _ := ar.Dx(), ar.Dy(), 30.0
					rotatedimg = image.NewRGBA(image.Rect(0, 0, w, h))
					draw.Draw(rotatedimg, ar, imgtorotate, ar.Min, draw.Src)

					draw2d.DrawImage(imgtorotate, rotatedimg, tr, draw.Src, draw2d.LinearFilter)

					// open new window
					splitfilename := strings.Split(filename, `.`)
					newname = splitfilename[0] + "_rotated" + `.` + splitfilename[1]

					//err = imaging.Save(rotatedimg, newname)
					err = imaging.Save(rotatedimg, newname)

					if err != nil {
						panic(err)
					}
					/*
						wells = PickAgain(newname, numbertopick, setplateperimeterfirst, false)
						return*/
					/*img0 = opencv.LoadImage(newname)
					filename = newname
					if img0 == nil {
						panic("LoadImage fail")
					}
					defer img0.Release()

					fmt.Print("Hot keys: \n",
						"\tESC - quit the program\n",
						"\tr - restore the original image\n",
						"\ti or ENTER - run inpainting algorithm\n",
						"\t\t(before running it, paint something on the image)\n",
					)

					//img := img0.Clone()
					//inpainted := img0.Clone()
					//inpaint_mask := opencv.CreateImage(img0.Width(), img0.Height(), 8, 1)

					//opencv.Zero(inpaint_mask)
					//opencv.Zero( inpainted )

					win := opencv.NewWindow("rotated")
					defer win.Destroy()
					// change status
					rotated = true
					fmt.Println("Now choose top left of plate")

					*/
				}
			} else if rotate && setplateperimeterfirst && counter == 3 && rotated || setplateperimeterfirst && counter == 1 && rotate == false {
				topleft = opencv.Point{x, y}
				fmt.Println("Topleft chosen")
				fmt.Println("Now choose bottom right of plate")

			} else if rotate && setplateperimeterfirst && counter == 4 && rotated || setplateperimeterfirst && counter == 2 && rotate == false {
				bottomright = opencv.Point{x, y}
				fmt.Println("plate boundaries", x, y)
				opencv.Rectangle(img, topleft, bottomright, rgb, 1, 8, 0)
				win.ShowImage(img)
			} else {

				prev_pt = opencv.Point{x, y}
				fmt.Println("actual pixels:", x, y)
				fmt.Println("imagesize:", img.ImageSize())
				fmt.Println("width:", img.Width())
				fmt.Println("height:", img.Height())
				fmt.Println("Well:", PixelstoWellPosition(x, y, plate, img))
				fmt.Println("colony count:", counter)

				if setplateperimeterfirst {
					well = PixelstoWellPositionFromRectangle(x, y, plate, topleft, bottomright)
				} else if setplateperimeterfirst == false {
					well = PixelstoWellPosition(x, y, plate, img)
				}

				wells = append(wells, well)
				fmt.Println(wells)

				// draw circle on click

				rgb = opencv.ScalarAll(255.0)
				opencv.Circle(img, prev_pt, 1, rgb, 5, 8, 0)

				win.ShowImage(img)

				if rotate && setplateperimeterfirst && counter == numbertopick+4 {
					//return
					splitfilename := strings.Split(filename, `.`)

					newname := splitfilename[0] + "_picked" + `.` + splitfilename[1]
					opencv.SaveImage(newname, img, 0)

					err = ExportCSVMultipleNames(ctx, exportFileName, plate, "ColonyPlate", wells, names, "colony")
					os.Exit(0)
					return
				}
				if rotate == false && setplateperimeterfirst && counter == numbertopick+2 {
					//return
					splitfilename := strings.Split(filename, `.`)

					newname := splitfilename[0] + "_picked" + `.` + splitfilename[1]
					opencv.SaveImage(newname, img, 0)
					err = ExportCSVMultipleNames(ctx, exportFileName, plate, "ColonyPlate", wells, names, "colony")
					os.Exit(0)
					return
				}
				if rotate == false && setplateperimeterfirst == false && counter == numbertopick {
					//return

					err = ExportCSVMultipleNames(ctx, exportFileName, plate, "ColonyPlate", wells, names, "colony")
					os.Exit(0)
					return
				}
			}
		} else if event == opencv.CV_EVENT_MOUSEMOVE &&
			(flags&opencv.CV_EVENT_FLAG_LBUTTON) != 0 {
			pt := opencv.Point{x, y}
			if prev_pt.X < 0 {
				prev_pt = pt
			}

			rgb := opencv.ScalarAll(255.0)

			opencv.Rectangle(img, prev_pt, pt, rgb, 1, 8, 0)
			//opencv.Line(inpaint_mask, prev_pt, pt, rgb, 5, 8, 0)
			//opencv.Line(img, prev_pt, pt, rgb, 5, 8, 0)
			prev_pt = pt

			win.ShowImage(img)
		}
	})
	win.ShowImage(img)
	opencv.WaitKey(0)

	win2 := opencv.NewWindow("inpainted image")

	defer win2.Destroy()
	win2.ShowImage(inpainted)

	for {
		key := opencv.WaitKey(20)
		if key == 27 {
			//	os.Exit(0)
		} else if key == 'r' {
			opencv.Zero(inpaint_mask)
			opencv.Copy(img0, img, nil)
			win.ShowImage(img)
		} else if key == 'i' || key == '\n' {
			opencv.Inpaint(img, inpaint_mask, inpainted, 3,
				opencv.CV_INPAINT_TELEA,
			)
			win2.ShowImage(inpainted)
		}
	}
	//os.Exit(0)
	/*if rotated {
		wells = PickAgain(newname, numbertopick, setplateperimeterfirst, false)
	}*/
	return
}

func Pick(imagefile string, plate *wtype.LHPlate, numbertopick int, setplateperimeterfirst bool, rotate bool) (wells []string) {
	var topleft opencv.Point
	var bottomright opencv.Point
	var topright opencv.Point
	var well string
	var rotated bool
	var rotatedimg draw.Image
	var newname string
	counter := 0
	//_, currentfile, _, _ := runtime.Caller(0)
	filename := imagefile //:= path.Join(path.Dir(currentfile), imagefile)
	wells = make([]string, 0)
	if len(os.Args) == 2 {
		filename = os.Args[1]
	}

	img0 := opencv.LoadImage(filename)
	if img0 == nil {
		panic("LoadImage fail")
	}

	// resize image if bigger than full screen
	if img0.Height() > 801 {
		img0.Release()
		imgtoresize, err := imaging.Open(filename)
		if err != nil {
			panic(err)
		}
		resizedimage := imaging.Resize(imgtoresize, 0, 800, imaging.CatmullRom)
		splitfilename := strings.Split(filename, `.`)

		newname = splitfilename[0] + "_resized" + `.` + splitfilename[1]

		err = imaging.Save(resizedimage, newname)
		if err != nil {
			panic(err)
		}
		img0 = opencv.LoadImage(newname)
		filename = newname
		if img0 == nil {
			panic("LoadImage fail")
		}
	}

	defer img0.Release()

	fmt.Print("Hot keys: \n",
		"\tESC - quit the program\n",
		"\tr - restore the original image\n",
		"\ti or ENTER - run inpainting algorithm\n",
		"\t\t(before running it, paint something on the image)\n",
	)
	fmt.Println("Click on top left corner of plate")

	img := img0.Clone()
	inpainted := img0.Clone()
	inpaint_mask := opencv.CreateImage(img0.Width(), img0.Height(), 8, 1)

	opencv.Zero(inpaint_mask)
	//opencv.Zero( inpainted )

	win := opencv.NewWindow("image")
	defer win.Destroy()

	prev_pt := opencv.Point{-1, -1}
	win.SetMouseCallback(func(event, x, y, flags int, param ...interface{}) {
		if img == nil {
			panic("img == nil")
		}

		if event == opencv.CV_EVENT_LBUTTONUP ||
			(flags&opencv.CV_EVENT_FLAG_LBUTTON) == 0 {
			prev_pt = opencv.Point{-1, -1}
		} else if event == opencv.CV_EVENT_LBUTTONDOWN {

			rgb := opencv.ScalarAll(255.0)
			counter = counter + 1

			if rotate && rotated == false {
				imgtorotate, err := imaging.Open(filename)

				if err != nil {
					panic(err)
				}
				if counter == 1 {
					fmt.Println("Rotating image")
					fmt.Println("topleft point:", x, y)
					topleft = opencv.Point{x, y}
				} else if counter == 2 {
					topright = opencv.Point{x, y}
					fmt.Println("topright point:", x, y)
					opposite := float64(topleft.Y) - float64(topright.Y)
					adjacent := float64(topright.X) - float64(topleft.X)
					tantheta := opposite / adjacent
					fmt.Println("adjacent:", adjacent)
					fmt.Println("opposite:", opposite)
					fmt.Println("costheta:", tantheta)
					thetainrad := math.Atan(tantheta)
					fmt.Println("thetainrad:", thetainrad)
					degrees := (180 / math.Pi) * thetainrad
					fmt.Println("degrees:", degrees)
					tr := draw2d.NewRotationMatrix(thetainrad)

					//
					ar := imgtorotate.Bounds()
					w, h, _ := ar.Dx(), ar.Dy(), 30.0
					rotatedimg = image.NewRGBA(image.Rect(0, 0, w, h))
					draw.Draw(rotatedimg, ar, imgtorotate, ar.Min, draw.Src)

					draw2d.DrawImage(imgtorotate, rotatedimg, tr, draw.Src, draw2d.LinearFilter)

					// open new window
					splitfilename := strings.Split(filename, `.`)
					newname = splitfilename[0] + "_rotated" + `.` + splitfilename[1]

					//err = imaging.Save(rotatedimg, newname)
					err = imaging.Save(rotatedimg, newname)

					if err != nil {
						panic(err)
					}
					/*
						wells = PickAgain(newname, numbertopick, setplateperimeterfirst, false)
						return*/
					/*img0 = opencv.LoadImage(newname)
					filename = newname
					if img0 == nil {
						panic("LoadImage fail")
					}
					defer img0.Release()

					fmt.Print("Hot keys: \n",
						"\tESC - quit the program\n",
						"\tr - restore the original image\n",
						"\ti or ENTER - run inpainting algorithm\n",
						"\t\t(before running it, paint something on the image)\n",
					)

					//img := img0.Clone()
					//inpainted := img0.Clone()
					//inpaint_mask := opencv.CreateImage(img0.Width(), img0.Height(), 8, 1)

					//opencv.Zero(inpaint_mask)
					//opencv.Zero( inpainted )

					win := opencv.NewWindow("rotated")
					defer win.Destroy()
					// change status
					rotated = true
					fmt.Println("Now choose top left of plate")

					*/
				}
			} else if rotate && setplateperimeterfirst && counter == 3 && rotated || setplateperimeterfirst && counter == 1 && rotate == false {
				topleft = opencv.Point{x, y}
				fmt.Println("Topleft chosen")
				fmt.Println("Now choose bottom right of plate")

			} else if rotate && setplateperimeterfirst && counter == 4 && rotated || setplateperimeterfirst && counter == 2 && rotate == false {
				bottomright = opencv.Point{x, y}
				fmt.Println("plate boundaries", x, y)
				opencv.Rectangle(img, topleft, bottomright, rgb, 1, 8, 0)
				win.ShowImage(img)
			} else {

				prev_pt = opencv.Point{x, y}
				fmt.Println("actual pixels:", x, y)
				fmt.Println("imagesize:", img.ImageSize())
				fmt.Println("width:", img.Width())
				fmt.Println("height:", img.Height())
				fmt.Println("Well:", PixelstoWellPosition(x, y, plate, img))
				fmt.Println("colony count:", counter)

				if setplateperimeterfirst {
					well = PixelstoWellPositionFromRectangle(x, y, plate, topleft, bottomright)
				} else if setplateperimeterfirst == false {
					well = PixelstoWellPosition(x, y, plate, img)
				}

				wells = append(wells, well)
				fmt.Println(wells)

				// draw circle on click

				rgb = opencv.ScalarAll(255.0)
				opencv.Circle(img, prev_pt, 1, rgb, 5, 8, 0)

				win.ShowImage(img)

				if rotate && setplateperimeterfirst && counter == numbertopick+4 {
					//return
					splitfilename := strings.Split(filename, `.`)

					newname := splitfilename[0] + "_picked" + `.` + splitfilename[1]
					opencv.SaveImage(newname, img, 0)
					//os.Exit(0)
					return
				}
				if rotate == false && setplateperimeterfirst && counter == numbertopick+2 {
					//return
					splitfilename := strings.Split(filename, `.`)

					newname := splitfilename[0] + "_picked" + `.` + splitfilename[1]
					opencv.SaveImage(newname, img, 0)
					//os.Exit(0)
					return
				}
				if rotate == false && setplateperimeterfirst == false && counter == numbertopick {
					//return
					//os.Exit(0)
					return
				}
			}
		} else if event == opencv.CV_EVENT_MOUSEMOVE &&
			(flags&opencv.CV_EVENT_FLAG_LBUTTON) != 0 {
			pt := opencv.Point{x, y}
			if prev_pt.X < 0 {
				prev_pt = pt
			}

			rgb := opencv.ScalarAll(255.0)

			opencv.Rectangle(img, prev_pt, pt, rgb, 1, 8, 0)
			//opencv.Line(inpaint_mask, prev_pt, pt, rgb, 5, 8, 0)
			//opencv.Line(img, prev_pt, pt, rgb, 5, 8, 0)
			prev_pt = pt

			win.ShowImage(img)
		}
	})
	win.ShowImage(img)
	opencv.WaitKey(0)

	win2 := opencv.NewWindow("inpainted image")

	defer win2.Destroy()
	win2.ShowImage(inpainted)

	for {
		key := opencv.WaitKey(20)
		if key == 27 {
			//	os.Exit(0)
		} else if key == 'r' {
			opencv.Zero(inpaint_mask)
			opencv.Copy(img0, img, nil)
			win.ShowImage(img)
		} else if key == 'i' || key == '\n' {
			opencv.Inpaint(img, inpaint_mask, inpainted, 3,
				opencv.CV_INPAINT_TELEA,
			)
			win2.ShowImage(inpainted)
		}
	}
	//os.Exit(0)
	/*if rotated {
		wells = PickAgain(newname, numbertopick, setplateperimeterfirst, false)
	}*/
	return
}

func PickMultipleWells(imagefile string, plate *wtype.LHPlate, numbertopick int, setplateperimeterfirst bool, rotate bool) (wellsMap map[string]string) {
	var topleft opencv.Point
	var bottomright opencv.Point
	var topright opencv.Point
	var well string
	var rotated bool
	var rotatedimg draw.Image
	var newname string
	counter := 0
	//_, currentfile, _, _ := runtime.Caller(0)
	filename := imagefile //:= path.Join(path.Dir(currentfile), imagefile)
	wells := make([]string, 0)
	if len(os.Args) == 2 {
		filename = os.Args[1]
	}

	img0 := opencv.LoadImage(filename)
	if img0 == nil {
		panic("LoadImage fail")
	}

	// resize image if bigger than full screen
	if img0.Height() > 801 {
		img0.Release()
		imgtoresize, err := imaging.Open(filename)
		if err != nil {
			panic(err)
		}
		resizedimage := imaging.Resize(imgtoresize, 0, 800, imaging.CatmullRom)
		splitfilename := strings.Split(filename, `.`)

		newname = splitfilename[0] + "_resized" + `.` + splitfilename[1]

		err = imaging.Save(resizedimage, newname)
		if err != nil {
			panic(err)
		}
		img0 = opencv.LoadImage(newname)
		filename = newname
		if img0 == nil {
			panic("LoadImage fail")
		}
	}

	defer img0.Release()

	fmt.Print("Hot keys: \n",
		"\tESC - quit the program\n",
		"\tr - restore the original image\n",
		"\ti or ENTER - run inpainting algorithm\n",
		"\t\t(before running it, paint something on the image)\n",
	)

	img := img0.Clone()
	inpainted := img0.Clone()
	inpaint_mask := opencv.CreateImage(img0.Width(), img0.Height(), 8, 1)

	opencv.Zero(inpaint_mask)
	//opencv.Zero( inpainted )

	win := opencv.NewWindow("image")
	defer win.Destroy()

	prev_pt := opencv.Point{-1, -1}
	win.SetMouseCallback(func(event, x, y, flags int, param ...interface{}) {
		if img == nil {
			panic("img == nil")
		}

		if event == opencv.CV_EVENT_LBUTTONUP ||
			(flags&opencv.CV_EVENT_FLAG_LBUTTON) == 0 {
			prev_pt = opencv.Point{-1, -1}
		} else if event == opencv.CV_EVENT_LBUTTONDOWN {

			rgb := opencv.ScalarAll(255.0)
			counter = counter + 1

			if rotate && rotated == false {
				imgtorotate, err := imaging.Open(filename)

				if err != nil {
					panic(err)
				}
				if counter == 1 {
					fmt.Println("Rotating image")
					fmt.Println("topleft point:", x, y)
					topleft = opencv.Point{x, y}
				} else if counter == 2 {
					topright = opencv.Point{x, y}
					fmt.Println("topright point:", x, y)
					opposite := float64(topleft.Y) - float64(topright.Y)
					adjacent := float64(topright.X) - float64(topleft.X)
					tantheta := opposite / adjacent
					fmt.Println("adjacent:", adjacent)
					fmt.Println("opposite:", opposite)
					fmt.Println("costheta:", tantheta)
					thetainrad := math.Atan(tantheta)
					fmt.Println("thetainrad:", thetainrad)
					degrees := (180 / math.Pi) * thetainrad
					fmt.Println("degrees:", degrees)
					tr := draw2d.NewRotationMatrix(thetainrad)

					//
					ar := imgtorotate.Bounds()
					w, h, _ := ar.Dx(), ar.Dy(), 30.0
					rotatedimg = image.NewRGBA(image.Rect(0, 0, w, h))
					draw.Draw(rotatedimg, ar, imgtorotate, ar.Min, draw.Src)

					draw2d.DrawImage(imgtorotate, rotatedimg, tr, draw.Src, draw2d.LinearFilter)

					// open new window
					splitfilename := strings.Split(filename, `.`)
					newname = splitfilename[0] + "_rotated" + `.` + splitfilename[1]

					//err = imaging.Save(rotatedimg, newname)
					err = imaging.Save(rotatedimg, newname)

					if err != nil {
						panic(err)
					}
					/*
						wells = PickAgain(newname, numbertopick, setplateperimeterfirst, false)
						return*/
					/*img0 = opencv.LoadImage(newname)
					filename = newname
					if img0 == nil {
						panic("LoadImage fail")
					}
					defer img0.Release()

					fmt.Print("Hot keys: \n",
						"\tESC - quit the program\n",
						"\tr - restore the original image\n",
						"\ti or ENTER - run inpainting algorithm\n",
						"\t\t(before running it, paint something on the image)\n",
					)

					//img := img0.Clone()
					//inpainted := img0.Clone()
					//inpaint_mask := opencv.CreateImage(img0.Width(), img0.Height(), 8, 1)

					//opencv.Zero(inpaint_mask)
					//opencv.Zero( inpainted )

					win := opencv.NewWindow("rotated")
					defer win.Destroy()
					// change status
					rotated = true
					fmt.Println("Now choose top left of plate")

					*/
				}
			} else if rotate && setplateperimeterfirst && counter == 3 && rotated || setplateperimeterfirst && counter == 1 && rotate == false {
				topleft = opencv.Point{x, y}
				fmt.Println("Topleft chosen")
				fmt.Println("Now choose bottom right of plate")

			} else if rotate && setplateperimeterfirst && counter == 4 && rotated || setplateperimeterfirst && counter == 2 && rotate == false {
				bottomright = opencv.Point{x, y}
				fmt.Println("plate boundaries", x, y)
				opencv.Rectangle(img, topleft, bottomright, rgb, 1, 8, 0)
				win.ShowImage(img)
			} else {

				prev_pt = opencv.Point{x, y}
				fmt.Println("actual pixels:", x, y)
				fmt.Println("imagesize:", img.ImageSize())
				fmt.Println("width:", img.Width())
				fmt.Println("height:", img.Height())
				fmt.Println("Well:", PixelstoWellPosition(x, y, plate, img))
				fmt.Println("colony count:", counter)

				if setplateperimeterfirst {
					well = PixelstoWellPositionFromRectangle(x, y, plate, topleft, bottomright)
				} else if setplateperimeterfirst == false {
					well = PixelstoWellPosition(x, y, plate, img)
				}

				wells = append(wells, well)
				fmt.Println(wells)

				// draw circle on click

				rgb = opencv.ScalarAll(255.0)
				opencv.Circle(img, prev_pt, 1, rgb, 5, 8, 0)

				win.ShowImage(img)

				if rotate && setplateperimeterfirst && counter == numbertopick+4 {
					//return
					splitfilename := strings.Split(filename, `.`)

					newname := splitfilename[0] + "_picked" + `.` + splitfilename[1]
					opencv.SaveImage(newname, img, 0)
					//os.Exit(0)
					return
				}
				if rotate == false && setplateperimeterfirst && counter == numbertopick+2 {
					//return
					splitfilename := strings.Split(filename, `.`)

					newname := splitfilename[0] + "_picked" + `.` + splitfilename[1]
					opencv.SaveImage(newname, img, 0)
					//os.Exit(0)
					return
				}
				if rotate == false && setplateperimeterfirst == false && counter == numbertopick {
					//return
					//os.Exit(0)
					return
				}
			}
		} else if event == opencv.CV_EVENT_MOUSEMOVE &&
			(flags&opencv.CV_EVENT_FLAG_LBUTTON) != 0 {
			pt := opencv.Point{x, y}
			if prev_pt.X < 0 {
				prev_pt = pt
			}

			rgb := opencv.ScalarAll(255.0)

			opencv.Rectangle(img, prev_pt, pt, rgb, 1, 8, 0)
			//opencv.Line(inpaint_mask, prev_pt, pt, rgb, 5, 8, 0)
			//opencv.Line(img, prev_pt, pt, rgb, 5, 8, 0)
			prev_pt = pt

			win.ShowImage(img)
		}
	})
	win.ShowImage(img)
	opencv.WaitKey(0)

	win2 := opencv.NewWindow("inpainted image")

	defer win2.Destroy()
	win2.ShowImage(inpainted)

	for {
		key := opencv.WaitKey(20)
		if key == 27 {
			//	os.Exit(0)
		} else if key == 'r' {
			opencv.Zero(inpaint_mask)
			opencv.Copy(img0, img, nil)
			win.ShowImage(img)
		} else if key == 'i' || key == '\n' {
			opencv.Inpaint(img, inpaint_mask, inpainted, 3,
				opencv.CV_INPAINT_TELEA,
			)
			win2.ShowImage(inpainted)
		}
	}
	//os.Exit(0)
	/*if rotated {
		wells = PickAgain(newname, numbertopick, setplateperimeterfirst, false)
	}*/
	return
}

/*
func PickAgain(imagefile string, numbertopick int, setplateperimeterfirst bool, rotate bool) (wells []string) {
	var topleft opencv.Point
	var bottomright opencv.Point
	var topright opencv.Point
	var well string
	var rotated bool
	var rotatedimg draw.Image
	counter := 0
	//_, currentfile, _, _ := runtime.Caller(0)
	filename := imagefile //:= path.Join(path.Dir(currentfile), imagefile)
	wells = make([]string, 0)
	if len(os.Args) == 2 {
		filename = os.Args[1]
	}

	img0 := opencv.LoadImage(filename)
	if img0 == nil {
		panic("LoadImage fail")
	}

	// resize image if bigger than full screen
	if img0.Height() > 801 {
		img0.Release()
		imgtoresize, err := imaging.Open(filename)
		if err != nil {
			panic(err)
		}
		resizedimage := imaging.Resize(imgtoresize, 0, 800, imaging.CatmullRom)
		splitfilename := strings.Split(filename, `.`)

		newname := splitfilename[0] + "_resized" + `.` + splitfilename[1]

		err = imaging.Save(resizedimage, newname)
		if err != nil {
			panic(err)
		}
		img0 = opencv.LoadImage(newname)
		filename = newname
		if img0 == nil {
			panic("LoadImage fail")
		}
	}

	defer img0.Release()

	fmt.Print("Hot keys: \n",
		"\tESC - quit the program\n",
		"\tr - restore the original image\n",
		"\ti or ENTER - run inpainting algorithm\n",
		"\t\t(before running it, paint something on the image)\n",
	)

	img := img0.Clone()
	inpainted := img0.Clone()
	inpaint_mask := opencv.CreateImage(img0.Width(), img0.Height(), 8, 1)

	opencv.Zero(inpaint_mask)
	//opencv.Zero( inpainted )

	win := opencv.NewWindow("image")
	defer win.Destroy()

	prev_pt := opencv.Point{-1, -1}
	win.SetMouseCallback(func(event, x, y, flags int, param ...interface{}) {
		if img == nil {
			os.Exit(0)
		}

		if event == opencv.CV_EVENT_LBUTTONUP ||
			(flags&opencv.CV_EVENT_FLAG_LBUTTON) == 0 {
			prev_pt = opencv.Point{-1, -1}
		} else if event == opencv.CV_EVENT_LBUTTONDOWN {

			rgb := opencv.ScalarAll(255.0)
			counter = counter + 1

			if rotate && rotated == false {
				imgtorotate, err := imaging.Open(filename)

				if err != nil {
					panic(err)
				}
				if counter == 1 {
					fmt.Println("Rotating image")
					fmt.Println("topleft point:", x, y)
					topleft = opencv.Point{x, y}
				} else if counter == 2 {
					topright = opencv.Point{x, y}
					fmt.Println("topright point:", x, y)
					opposite := float64(topleft.Y) - float64(topright.Y)
					adjacent := float64(topright.X) - float64(topleft.X)
					tantheta := opposite / adjacent
					fmt.Println("adjacent:", adjacent)
					fmt.Println("opposite:", opposite)
					fmt.Println("costheta:", tantheta)
					thetainrad := math.Atan(tantheta)
					fmt.Println("thetainrad:", thetainrad)
					degrees := (180 / math.Pi) * thetainrad
					fmt.Println("degrees:", degrees)
					tr := draw2d.NewRotationMatrix(thetainrad)

					//
					ar := imgtorotate.Bounds()
					w, h, _ := ar.Dx(), ar.Dy(), 30.0
					rotatedimg = image.NewRGBA(image.Rect(0, 0, w, h))
					draw.Draw(rotatedimg, ar, imgtorotate, ar.Min, draw.Src)

					draw2d.DrawImage(imgtorotate, rotatedimg, tr, draw.Src, draw2d.LinearFilter)

					// open new window
					splitfilename := strings.Split(filename, `.`)
					newname := splitfilename[0] + "_rotated" + `.` + splitfilename[1]

					//err = imaging.Save(rotatedimg, newname)
					err = imaging.Save(rotatedimg, newname)
					if err != nil {
						panic(err)
					}
					img0 = opencv.LoadImage(newname)
					filename = newname
					if img0 == nil {
						panic("LoadImage fail")
					}
					defer img0.Release()

					fmt.Print("Hot keys: \n",
						"\tESC - quit the program\n",
						"\tr - restore the original image\n",
						"\ti or ENTER - run inpainting algorithm\n",
						"\t\t(before running it, paint something on the image)\n",
					)

					//img := img0.Clone()
					//inpainted := img0.Clone()
					//inpaint_mask := opencv.CreateImage(img0.Width(), img0.Height(), 8, 1)

					//opencv.Zero(inpaint_mask)
					//opencv.Zero( inpainted )

					win := opencv.NewWindow("rotated")
					defer win.Destroy()
					// change status
					rotated = true
					fmt.Println("Now choose top left of plate")
				}
			} else if rotate && setplateperimeterfirst && counter == 3 && rotated || setplateperimeterfirst && counter == 1 && rotate == false {
				topleft = opencv.Point{x, y}
				fmt.Println("Topleft chosen")
				fmt.Println("Now choose bottom right of plate")

			} else if rotate && setplateperimeterfirst && counter == 4 && rotated || setplateperimeterfirst && counter == 2 && rotate == false {
				bottomright = opencv.Point{x, y}
				fmt.Println("plate boundaries", x, y)
				opencv.Rectangle(img, topleft, bottomright, rgb, 1, 8, 0)
				win.ShowImage(img)
			} else {

				prev_pt = opencv.Point{x, y}
				fmt.Println("actual pixels:", x, y)
				fmt.Println("imagesize:", img.ImageSize())
				fmt.Println("width:", img.Width())
				fmt.Println("height:", img.Height())
				fmt.Println("Well:", PixelstoWellPosition(x, y, img))
				fmt.Println("colony count:", counter)

				if setplateperimeterfirst {
					well = PixelstoWellPositionFromRectangle(x, y, topleft, bottomright)
				} else if setplateperimeterfirst == false {
					well = PixelstoWellPosition(x, y, img)
				}

				wells = append(wells, well)
				fmt.Println(wells)

				// draw circle on click

				rgb = opencv.ScalarAll(255.0)
				opencv.Circle(img, prev_pt, 1, rgb, 5, 8, 0)

				win.ShowImage(img)

				if rotate && setplateperimeterfirst && counter == numbertopick+4 {
					//return
					splitfilename := strings.Split(filename, `.`)

					newname := splitfilename[0] + "_picked" + `.` + splitfilename[1]
					opencv.SaveImage(newname, img, 0)
					os.Exit(0)
					return
				}
				if rotate == false && setplateperimeterfirst && counter == numbertopick+2 {
					//return
					splitfilename := strings.Split(filename, `.`)

					newname := splitfilename[0] + "_picked" + `.` + splitfilename[1]
					opencv.SaveImage(newname, img, 0)
					os.Exit(0)
					return
				}
				if rotate == false && setplateperimeterfirst == false && counter == numbertopick {
					//return
					os.Exit(0)
					return
				}
			}
		} else if event == opencv.CV_EVENT_MOUSEMOVE &&
			(flags&opencv.CV_EVENT_FLAG_LBUTTON) != 0 {
			pt := opencv.Point{x, y}
			if prev_pt.X < 0 {
				prev_pt = pt
			}

			rgb := opencv.ScalarAll(255.0)

			opencv.Rectangle(img, prev_pt, pt, rgb, 1, 8, 0)
			//opencv.Line(inpaint_mask, prev_pt, pt, rgb, 5, 8, 0)
			//opencv.Line(img, prev_pt, pt, rgb, 5, 8, 0)
			prev_pt = pt

			win.ShowImage(img)
		}
	})
	win.ShowImage(img)
	opencv.WaitKey(0)

	win2 := opencv.NewWindow("inpainted image")

	defer win2.Destroy()
	win2.ShowImage(inpainted)

	for {
		key := opencv.WaitKey(20)
		if key == 27 {
			os.Exit(0)
		} else if key == 'r' {
			opencv.Zero(inpaint_mask)
			opencv.Copy(img0, img, nil)
			win.ShowImage(img)
		} else if key == 'i' || key == '\n' {
			opencv.Inpaint(img, inpaint_mask, inpainted, 3,
				opencv.CV_INPAINT_TELEA,
			)
			win2.ShowImage(inpainted)
		}
	}
	os.Exit(0)
	if rotated {

	}
	return
}
*/
/*
func Count(imagefile string, setplateperimeterfirst bool, rotate bool) (wells []string) {
	var topleft opencv.Point
	var bottomright opencv.Point
	var topright opencv.Point
	var well string
	var rotated bool
	var rotatedimg draw.Image
	var newname string
	counter := 0
	//_, currentfile, _, _ := runtime.Caller(0)
	filename := imagefile //:= path.Join(path.Dir(currentfile), imagefile)
	wells = make([]string, 0)
	if len(os.Args) == 2 {
		filename = os.Args[1]
	}

	img0 := opencv.LoadImage(filename)
	if img0 == nil {
		panic("LoadImage fail")
	}

	// resize image if bigger than full screen
	if img0.Height() > 801 {
		img0.Release()
		imgtoresize, err := imaging.Open(filename)
		if err != nil {
			panic(err)
		}
		resizedimage := imaging.Resize(imgtoresize, 0, 800, imaging.CatmullRom)
		splitfilename := strings.Split(filename, `.`)

		newname = splitfilename[0] + "_resized" + `.` + splitfilename[1]

		err = imaging.Save(resizedimage, newname)
		if err != nil {
			panic(err)
		}
		img0 = opencv.LoadImage(newname)
		filename = newname
		if img0 == nil {
			panic("LoadImage fail")
		}
	}

	defer img0.Release()

	fmt.Print("Hot keys: \n",
		"\tESC - quit the program\n",
		"\tr - restore the original image\n",
		"\ti or ENTER - run inpainting algorithm\n",
		"\t\t(before running it, paint something on the image)\n",
	)

	img := img0.Clone()
	inpainted := img0.Clone()
	inpaint_mask := opencv.CreateImage(img0.Width(), img0.Height(), 8, 1)

	opencv.Zero(inpaint_mask)
	//opencv.Zero( inpainted )

	win := opencv.NewWindow("image")
	defer win.Destroy()

	prev_pt := opencv.Point{-1, -1}
	win.SetMouseCallback(func(event, x, y, flags int, param ...interface{}) {
		if img == nil {
			panic("img == nil")
		}

		if event == opencv.CV_EVENT_LBUTTONUP ||
			(flags&opencv.CV_EVENT_FLAG_LBUTTON) == 0 {
			prev_pt = opencv.Point{-1, -1}
		} else if event == opencv.CV_EVENT_LBUTTONDOWN {

			rgb := opencv.ScalarAll(255.0)
			counter = counter + 1

			if rotate && rotated == false {
				imgtorotate, err := imaging.Open(filename)

				if err != nil {
					panic(err)
				}
				if counter == 1 {
					fmt.Println("Rotating image")
					fmt.Println("topleft point:", x, y)
					topleft = opencv.Point{x, y}
				} else if counter == 2 {
					topright = opencv.Point{x, y}
					fmt.Println("topright point:", x, y)
					opposite := float64(topleft.Y) - float64(topright.Y)
					adjacent := float64(topright.X) - float64(topleft.X)
					tantheta := opposite / adjacent
					fmt.Println("adjacent:", adjacent)
					fmt.Println("opposite:", opposite)
					fmt.Println("costheta:", tantheta)
					thetainrad := math.Atan(tantheta)
					fmt.Println("thetainrad:", thetainrad)
					degrees := (180 / math.Pi) * thetainrad
					fmt.Println("degrees:", degrees)
					tr := draw2d.NewRotationMatrix(thetainrad)

					//
					ar := imgtorotate.Bounds()
					w, h, _ := ar.Dx(), ar.Dy(), 30.0
					rotatedimg = image.NewRGBA(image.Rect(0, 0, w, h))
					draw.Draw(rotatedimg, ar, imgtorotate, ar.Min, draw.Src)

					draw2d.DrawImage(imgtorotate, rotatedimg, tr, draw.Src, draw2d.LinearFilter)

					// open new window
					splitfilename := strings.Split(filename, `.`)
					newname = splitfilename[0] + "_rotated" + `.` + splitfilename[1]

					//err = imaging.Save(rotatedimg, newname)
					err = imaging.Save(rotatedimg, newname)

					if err != nil {
						panic(err)
					}

				}
			} else if rotate && setplateperimeterfirst && counter == 3 && rotated || setplateperimeterfirst && counter == 1 && rotate == false {
				topleft = opencv.Point{x, y}
				fmt.Println("Topleft chosen")
				fmt.Println("Now choose bottom right of plate")

			} else if rotate && setplateperimeterfirst && counter == 4 && rotated || setplateperimeterfirst && counter == 2 && rotate == false {
				bottomright = opencv.Point{x, y}
				fmt.Println("plate boundaries", x, y)
				opencv.Rectangle(img, topleft, bottomright, rgb, 1, 8, 0)
				win.ShowImage(img)
			} else {

				prev_pt = opencv.Point{x, y}
				fmt.Println("actual pixels:", x, y)
				fmt.Println("imagesize:", img.ImageSize())
				fmt.Println("width:", img.Width())
				fmt.Println("height:", img.Height())
				fmt.Println("Well:", PixelstoWellPosition(x, y, img))
				fmt.Println("colony count:", counter)

				if setplateperimeterfirst {
					well = PixelstoWellPositionFromRectangle(x, y, topleft, bottomright)
				} else if setplateperimeterfirst == false {
					well = PixelstoWellPosition(x, y, img)
				}

				wells = append(wells, well)
				fmt.Println(wells)

				// draw circle on click

				rgb = opencv.ScalarAll(255.0)
				opencv.Circle(img, prev_pt, 1, rgb, 5, 8, 0)

				win.ShowImage(img)

				if rotate && setplateperimeterfirst && counter == numbertopick+4 {
					//return
					splitfilename := strings.Split(filename, `.`)

					newname := splitfilename[0] + "_picked" + `.` + splitfilename[1]
					opencv.SaveImage(newname, img, 0)
					//os.Exit(0)
					return
				}
				if rotate == false && setplateperimeterfirst && counter == numbertopick+2 {
					//return
					splitfilename := strings.Split(filename, `.`)

					newname := splitfilename[0] + "_picked" + `.` + splitfilename[1]
					opencv.SaveImage(newname, img, 0)
					//os.Exit(0)
					return
				}
				if rotate == false && setplateperimeterfirst == false && counter == numbertopick {
					//return
					//os.Exit(0)
					return
				}
			}
		} else if event == opencv.CV_EVENT_MOUSEMOVE &&
			(flags&opencv.CV_EVENT_FLAG_LBUTTON) != 0 {
			pt := opencv.Point{x, y}
			if prev_pt.X < 0 {
				prev_pt = pt
			}

			rgb := opencv.ScalarAll(255.0)

			opencv.Rectangle(img, prev_pt, pt, rgb, 1, 8, 0)
			//opencv.Line(inpaint_mask, prev_pt, pt, rgb, 5, 8, 0)
			//opencv.Line(img, prev_pt, pt, rgb, 5, 8, 0)
			prev_pt = pt

			win.ShowImage(img)
		}
	})
	win.ShowImage(img)
	opencv.WaitKey(0)

	win2 := opencv.NewWindow("inpainted image")

	defer win2.Destroy()
	win2.ShowImage(inpainted)

	for {
		key := opencv.WaitKey(20)
		if key == 27 {
			//	os.Exit(0)
		} else if key == 'r' {
			opencv.Zero(inpaint_mask)
			opencv.Copy(img0, img, nil)
			win.ShowImage(img)
		} else if key == 'i' || key == '\n' {
			opencv.Inpaint(img, inpaint_mask, inpainted, 3,
				opencv.CV_INPAINT_TELEA,
			)
			win2.ShowImage(inpainted)
		}
	}

	return
}*/

func PixelstoWellPosition(x, y int, plate *wtype.LHPlate, image *opencv.IplImage) (a1 string) {

	fmt.Println("func pixels:", x, y)
	fmt.Println("func width:", image.Width())
	fmt.Println("func height:", image.Height())

	numberofcolsX := (float64(x) / float64(image.Width())) * float64(plate.WellsX())
	numberofrowsY := (float64(y) / float64(image.Height())) * float64(plate.WellsY())

	colint, err := wutil.RoundDown(numberofcolsX)
	if err != nil {
		return
	}
	rowint, err := wutil.RoundDown(numberofrowsY)
	if err != nil {
		return
	}
	a1 = wutil.NumToAlpha(rowint+1) + strconv.Itoa(colint+1)
	return
}
func PixelstoWellPositionFromRectangle(x, y int, plate *wtype.LHPlate, topleft, bottomright opencv.Point) (a1 string) {

	fmt.Println("func pixels:", x, y)
	fmt.Println("rectangle width:", bottomright.X-topleft.X)
	fmt.Println("rectangle height:", bottomright.Y-topleft.Y)

	numberofcolsX := (float64(x-topleft.X) / float64(bottomright.X-topleft.X)) * float64(plate.WellsX())
	numberofrowsY := (float64(y-topleft.Y) / float64(bottomright.Y-topleft.Y)) * float64(plate.WellsY())

	colint, err := wutil.RoundDown(numberofcolsX)
	if err != nil {
		return
	}
	rowint, err := wutil.RoundDown(numberofrowsY)
	if err != nil {
		return
	}
	a1 = wutil.NumToAlpha(rowint+1) + strconv.Itoa(colint+1)
	return
}

func ExportCSV(ctx context.Context, ExportFileName string, plateForCoordinates *wtype.LHPlate, platename string, wellstopick []string, nameprepend, liquidtypename string) (err error) {
	fmt.Println("Generating csv output file with well coordinates for each colony")
	liquids := make([]*wtype.LHComponent, 0)
	volumes := make([]wunit.Volume, 0)

	residualvolume := plateForCoordinates.Welltype.ResidualVolume()

	colonyvol := (wunit.CopyVolume(residualvolume))

	colonyvol.Add(wunit.NewVolume(1.5, "ul"))

	for i := range wellstopick {
		volumes = append(volumes, colonyvol)
		liquid, err := inventory.NewComponent(ctx, liquidtypename)
		if err != nil {
			return err
		}

		liquid.CName = nameprepend + strconv.Itoa(i)
		liquids = append(liquids, liquid)
	}

	err = wtype.ExportPlateCSV(ExportFileName, plateForCoordinates, platename, wellstopick, liquids, volumes)
	return
}

func ExportCSVMultipleNames(ctx context.Context, ExportFileName string, plateForCoordinates *wtype.LHPlate, platename string, wellstopick []string, names []string, liquidtypename string) (err error) {
	fmt.Println("Generating csv output file with well coordinates for each colony")
	liquids := make([]*wtype.LHComponent, 0)
	volumes := make([]wunit.Volume, 0)

	if len(wellstopick) != len(names) {
		err = fmt.Errorf("Length of array mismatch. "+"Length of wellstopick: ", strconv.Itoa(len(wellstopick)), " Length of names: ", strconv.Itoa(len(names)))
	}

	residualvolume := plateForCoordinates.Welltype.ResidualVolume()

	colonyvol := (wunit.CopyVolume(residualvolume))

	colonyvol.Add(wunit.NewVolume(1.5, "ul"))

	for i := range wellstopick {
		volumes = append(volumes, colonyvol)
		liquid, err := inventory.NewComponent(ctx, liquidtypename)
		if err != nil {
			return err
		}

		liquid.CName = names[i]
		liquids = append(liquids, liquid)
	}

	err = wtype.ExportPlateCSV(ExportFileName, plateForCoordinates, platename, wellstopick, liquids, volumes)
	return
}

// This function used internally to convert any image type to NRGBA if needed.
func toNRGBA(img image.Image) *image.NRGBA {
	srcBounds := img.Bounds()
	if srcBounds.Min.X == 0 && srcBounds.Min.Y == 0 {
		if src0, ok := img.(*image.NRGBA); ok {
			return src0
		}
	}
	return imaging.Clone(img)
}
