// Part of the Antha language
// Copyright (C) 2015 The Antha authors. All rights reserved.
//
// This program is free software; you can redistribute it and/or
// modify it under the terms of the GNU General Public License
// as published by the Free Software Foundation; either version 2
// of the License, or (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program; if not, write to the Free Software
// Foundation, Inc., 51 Franklin Street, Fifth Floor, Boston, MA  02110-1301, USA.
//
// For more information relating to the software or licensing issues please
// contact license@antha-lang.org or write to the Antha team c/o
// Synthace Ltd. The London Bioscience Innovation Centre
// 2 Royal College St, London NW1 0NH UK

// Package plot provides methods for plotting data. This package uses the gonum
// plot library Plot type.  Produce a plot using the Plot function and then
// export into an antha filetype using the Export function
package plot

import (
	"bytes"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/spreadsheet"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/tealeg/xlsx"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
)

// Export creates a wtype.file from the plot data. Heights and lengths can be parsed from strings i.e. 10cm.
// If no valid height or length is specified default values of 10cm will be used but an error will also be returned.
// If the desired filename specified does not contain a file extension, a png file will be used as the default file format.
func Export(plt *plot.Plot, heightstr string, lengthstr string, filename string) (file wtype.File, err error) {

	var errtoreturn error

	length, err := vg.ParseLength(lengthstr)
	if err != nil {
		length = 10 * vg.Centimeter
		errtoreturn = err
	}
	height, err := vg.ParseLength(heightstr)
	if err != nil {
		height = 10 * vg.Centimeter
		if errtoreturn != nil {
			errtoreturn = fmt.Errorf(errtoreturn.Error(), " + ", err.Error())
		}
	}

	var plotFormat string

	if len(filepath.Ext(filename)) == 0 {
		plotFormat = "png"
		filename = strings.Join([]string{filename, plotFormat}, ".")
	} else {
		plotFormat = strings.TrimLeft(filepath.Ext(filename), ".")
	}

	w, err := plt.WriterTo(length, height, plotFormat)
	if err != nil {
		if errtoreturn != nil {
			errtoreturn = fmt.Errorf(errtoreturn.Error(), " + ", err.Error())
		}
		return file, errtoreturn
	}

	var out bytes.Buffer
	_, err = w.WriteTo(&out)
	if err != nil {
		return
	}

	file.Name = filename
	err = file.WriteAll(out.Bytes())
	return
}

// Plot creates a plot from float data. Multiple sets of y values may be
// specified for a set of x values.  The length of any yvalue dataset must be
// equal to the length of xvalues or the function will stop and return an
// error.
func Plot(Xvalues []float64, Yvaluearray [][]float64) (plt *plot.Plot, err error) {
	// now plot the graph

	// the data points
	pts := make([]plotter.XYer, 0) //len(Xdatarange))

	// each specific set for each datapoint
	for _, ydataset := range Yvaluearray {

		if len(ydataset) != len(Xvalues) {
			return nil, fmt.Errorf("cannot plot x by y: x length %d is not the same as y length %d",
				len(Xvalues), len(ydataset))
		}

		xys := make(plotter.XYs, len(ydataset))
		for j := range xys {
			xys[j].X = Xvalues[j]
			xys[j].Y = ydataset[j]
		}
		pts = append(pts, xys)
	}
	plt, err = plot.New()

	if err != nil {
		return
	}

	// Create two lines connecting points and error bars. For
	// the first, each point is the mean x and y value and the
	// error bars give the 95% confidence intervals.  For the
	// second, each point is the median x and y value with the
	// error bars showing the minimum and maximum values.
	/*
	   	// fmt.Println("pts", pts)
	   	mean95, err := plotutil.NewErrorPoints(plotutil.MeanAndConf95, pts...)
	   	if err != nil {
	   		panic(err)
	   	}
	   	//medMinMax, err := plotutil.NewErrorPoints(plotutil.MedianAndMinMax, pts...)
	   //	if err != nil {
	   //		panic(err)
	   //	}
	   	plotutil.AddLinePoints(plt,
	   		"mean and 95% confidence", mean95,
	   	) //	"median and minimum and maximum", medMinMax)
	   	//plotutil.AddErrorBars(plt, mean95, medMinMax)

	   	// Add the points that are summarized by the error points.


	*/

	ptsinterface := make([]interface{}, 0)

	for i, pt := range pts {
		ptsinterface = append(ptsinterface, fmt.Sprint("run_", i))
		ptsinterface = append(ptsinterface, pt)
	}

	err = plotutil.AddScatters(plt, ptsinterface...) //AddScattersXYer(plt, pts)
	if err != nil {
		return
	}

	plt.Legend.Top = true
	plt.Legend.Left = true

	return
}

// AddAxisTitles adds axis titles to the plot
func AddAxisTitles(plt *plot.Plot, xtitle string, ytitle string) {
	plt.X.Label.Text = xtitle
	plt.Y.Label.Text = ytitle
}

// FromMinMaxPairs creates a plot from a spreadsheet.
func FromMinMaxPairs(sheet *xlsx.Sheet, Xminmax []string, Yminmaxarray [][]string, Exportedfilename string) {
	Xdatarange, err := spreadsheet.ConvertMinMaxtoArray(Xminmax)
	if err != nil {
		fmt.Println(Xminmax, Xdatarange)
		panic(err)
	}
	fmt.Println(Xdatarange)

	Ydatarangearray := make([][]string, 0)
	for i, Yminmax := range Yminmaxarray {
		Ydatarange, err := spreadsheet.ConvertMinMaxtoArray(Yminmax)
		if err != nil {
			panic(err)
		}
		if len(Xdatarange) != len(Ydatarange) {
			panicmessage := fmt.Sprint("for index", i, "of array", "len(Xdatarange) != len(Ydatarange)")
			panic(panicmessage)
		}
		Ydatarangearray = append(Ydatarangearray, Ydatarange)
		fmt.Println(Ydatarange)
	}
	FromSpreadsheet(sheet, Xdatarange, Ydatarangearray, Exportedfilename)
}

// FromSpreadsheet creates a plot from a spreadsheet.
func FromSpreadsheet(sheet *xlsx.Sheet, Xdatarange []string, Ydatarangearray [][]string, Exportedfilename string) {

	// now plot the graph

	// the data points
	pts := make([]plotter.XYer, 0) //len(Xdatarange))

	for ptsindex := 0; ptsindex < len(Xdatarange); ptsindex++ {

		// each specific set for each datapoint

		for Xdatarangeindex, Xdatapoint := range Xdatarange {

			xys := make(plotter.XYs, len(Ydatarangearray))

			xrow, xcol, err := spreadsheet.A1FormatToRowColumn(Xdatapoint)
			if err != nil {
				panic(err)
			}
			xpoint := sheet.Rows[xcol].Cells[xrow]

			// get each y point and work out average

			yfloats := make([]float64, 0)
			for _, Ydatarange := range Ydatarangearray {
				yrow, ycol, err := spreadsheet.A1FormatToRowColumn(Ydatarange[Xdatarangeindex])
				if err != nil {
					panic(err)
				}
				ypoint := sheet.Rows[ycol].Cells[yrow]
				yfloat, err := ypoint.Float()
				if err != nil {
					panic(err)
				}
				yfloats = append(yfloats, yfloat)

			}

			xfloat, err := xpoint.Float()
			if err != nil {
				panic(err)
			}

			for j := range xys {
				// fmt.Println("going here")
				fmt.Println(ptsindex)
				xys[j].X = xfloat
				xys[j].Y = yfloats[j]
			}
			pts = append(pts, xys) //
		}

	}
	plt, err := plot.New()

	if err != nil {
		panic(err)
	}

	// Create two lines connecting points and error bars. For
	// the first, each point is the mean x and y value and the
	// error bars give the 95% confidence intervals.  For the
	// second, each point is the median x and y value with the
	// error bars showing the minimum and maximum values.

	//	// fmt.Println("pts", pts)
	//	mean95, err := plotutil.NewErrorPoints(plotutil.MeanAndConf95, pts...)
	//	if err != nil {
	//		panic(err)
	//	}
	/*medMinMax, err := plotutil.NewErrorPoints(plotutil.MedianAndMinMax, pts...)
	if err != nil {
		panic(err)
	}*/
	//	plotutil.AddLinePoints(plt,
	//		"mean and 95% confidence", mean95,
	//	) //	"median and minimum and maximum", medMinMax)
	//plotutil.AddErrorBars(plt, mean95, medMinMax)

	// Add the points that are summarized by the error points.

	ptsinterface := make([]interface{}, 0)

	for _, pt := range pts {
		ptsinterface = append(ptsinterface, pt)
	}

	if err := plotutil.AddScatters(plt, ptsinterface...); err != nil {
		panic(err)
	}

	length, _ := vg.ParseLength("10cm") // nolint

	if err := plt.Save(length, length, Exportedfilename); err != nil {
		panic(err)
	}
}
