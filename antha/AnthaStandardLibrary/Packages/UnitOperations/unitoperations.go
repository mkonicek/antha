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

// Package for working with bioprocessing unitoperations
package UnitOperations

import (
	"time"

	"github.com/antha-lang/antha/antha/anthalib/mixer"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
)

type Pellet struct {
}

type Culture struct {
	wtype.Suspension
	wtype.Organism
}

func Separate(culture Culture) (supernatant *wtype.Liquid, pellet Pellet) {

	return
}

type Chromstep struct {
	Pipetstep
	Column
	Buffer *wtype.Liquid
	//Mobilephase ... this is really the buffer; should separate out chromstep from pipette step
}

type Pipetstep struct {
	Name          string
	Volume        wunit.Volume //transfer volume and process volume?
	Aspiraterate  wunit.FlowRate
	Aspiratepause time.Duration
	Dispenserate  wunit.FlowRate
	Dispensepause time.Duration
	Cycles        int
	//Mobilephase ... this is really the buffer; should separate out chromstep from pipette step
}

type Column struct {
	Beadsize           wunit.Length
	Stationaryphase    string
	Separationproperty string
	Diameter           wunit.Length
	Height             wunit.Length
	Packedvolume       wunit.Volume
}

type Phytips struct { // interface?
	Column
	Tip wtype.LHTip
}

// may already be functions for aspirate and dispense in anthalib
func Aspirate(column Column, mixture *wtype.Liquid, volume wunit.Volume, aspiraterate wunit.FlowRate) (aspiratedcolumn Column, aspiratedsolution *wtype.Liquid) {

	return
}

func Dispense(column Column, mixture *wtype.Liquid, volume wunit.Volume, aspiraterate wunit.FlowRate) (dispensedcolumn Column, dispensedsolution *wtype.Liquid) {

	return
}

/*
func PhysicaltoComponent(pellet *wtype.Physical) (component *wtype.Liquid) {
	// placeholder
	return
}
*/

func PelletToComponent(p Pellet) *wtype.Liquid {
	return wtype.NewLHComponent()
}

/*
func Resuspend(pellet Pellet, step Chromstep, column Column) (output_c *wtype.Liquid, processedcolumn Column) {

	var output *wtype.Liquid
	//input := PhysicaltoComponent(pellet)
	input := PelletToComponent(pellet)
	samples := make([]*wtype.Liquid, 0)
	samples = append(samples, step.Buffer, input)
	mixture := mixer.Mix(samples...)
	for i := 0; i < step.Pipetstep.Cycles; i++ {

		aspiratedcolumn, aspiratedsolution := Aspirate(column, mixture, step.Volume, step.Aspiraterate)
		time.Sleep(step.Aspiratepause)
		_, output = Dispense(aspiratedcolumn, aspiratedsolution, step.Volume, step.Dispenserate)
		time.Sleep(step.Dispensepause)
	}

	processedcolumn = column
	output_c = wtype.SolutionToComponent(output)
	return output_c, processedcolumn
}
*/
func Chromatography(input *wtype.Liquid, step Chromstep, column Column) (output_c *wtype.Liquid, processedcolumn Column) {

	var output *wtype.Liquid
	samples := make([]*wtype.Liquid, 0)
	samples = append(samples, input, step.Buffer)
	mixture := mixer.Mix(samples...)
	for i := 0; i < step.Pipetstep.Cycles; i++ {

		aspiratedcolumn, aspiratedsolution := Aspirate(column, mixture, step.Volume, step.Aspiraterate)
		time.Sleep(step.Aspiratepause)
		_, output = Dispense(aspiratedcolumn, aspiratedsolution, step.Volume, step.Dispenserate)
		time.Sleep(step.Dispensepause)
	}

	processedcolumn = column
	return output, processedcolumn
}

func Blot(column Column, blotcycles int, blottime time.Duration) (blottedcolumn Column) {
	// placeholder for moving column to blot paper and dabbing and then waiting
	for i := 0; i < blotcycles; i++ {
		time.Sleep(blottime)
	}
	return blottedcolumn
}

func Dry(tips Column /*wtype.LHTip*/, Drytime time.Duration, Vacuumstrength float64) (drytips Column) {
	// set vacuum manifold to vacuum strength
	///move tips to vacuum position
	time.Sleep(Drytime)
	drytips = tips
	return drytips
}

/*
func Equilibration (Input Watersolution, step Chromstep, column Column) readycolumn Column {

		readycolumn = step Chromstep.cycles *
	            	(
						column ( aspirate mixture(
									step Chromstep.volume,
									step Chromstep.aspiraterate
									)
									wait step Chromstep.aspiratepause
									dispense mixture(
										step Chromstep.volume,
										step Chromstep.dispenserate
									)
									wait step Chromstep.dispensepause
								)
					)
					return readycolumn
}

func AirChromatography (Input Gas, step Chromstep, column Column) readycolumn Column {
	  mixture := mix (Input, step Chromstep.buffer)
		Output = step Chromstep.cycles *
	            	(
						column ( aspirate mixture(
								step Chromstep.volume,
								step Chromstep.aspiraterate
									)
								wait step Chromstep.aspiratepause
									dispense mixture(
										step Chromstep.volume,
										step Chromstep.dispenserate
									)
									wait step Chromstep.dispensepause
								)
					)
					return readycolumn
}


func Resuspension (Input Pellet, step pipetstep) output Suspension {
	  mixture := mix (step Pipetstep.buffer, Input)
		Output = step Pipetstep.cycles *
	            	(
						aspirate mixture(
								step Pipetstep.volume,
								step Pipetstep.aspiraterate
							)
								wait step Pipetstep.aspiratepause
									dispense mixture(
										step Pipetstep.volume,
										step Pipetstep.dispenserate
								)
								wait step Pipetstep.dispensepause
					)
					return output
}

func Lysis (Input Suspension, step pipetstep) output Lysate {
	  mixture := mix (step, Input)
		Output = step Pipetstep.cycles *
	            	(
						aspirate mixture(
								step Pipetstep.volume,
								step Pipetstep.aspiraterate
							)
								wait step Pipetstep.aspiratepause
									dispense mixture(
										step Pipetstep.volume,
										step Pipetstep.dispenserate
								)
								wait step Pipetstep.dispensepause
					)
					return output
}

func Precipitation (Input Suspension, step pipetstep) output Precipitate {
	  mixture := mix (step, Input)
		Output = step Pipetstep.cycles *
	            	(
						aspirate mixture(
								step Pipetstep.volume,
								step Pipetstep.aspiraterate
							)
								wait step Pipetstep.aspiratepause
									dispense mixture(
										step Pipetstep.volume,
										step Pipetstep.dispenserate
								)
								wait step Pipetstep.dispensepause
					)
					return output
}

func Growthcurve () {

}
*/
