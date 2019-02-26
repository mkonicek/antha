// /anthalib/driver/liquidhandling/driver.go: Part of the Antha language
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

package liquidhandling

import (
	"github.com/antha-lang/antha/microArch/driver"
)

// driver interface

type LiquidhandlingDriver interface {
	DriverType() ([]string, error)
	//AddPlateTo add an LHObject to a particular position in the liquid handler
	//position: the name of the position defined in LHProperties struct
	//plate: the LHObject to add
	//name: the name of the plate, should match wtype.NameOf(plate)
	AddPlateTo(position string, plate interface{}, name string) driver.CommandStatus
	//RemoveAllPlates remove every object in the machine
	RemoveAllPlates() driver.CommandStatus
	RemovePlateAt(position string) driver.CommandStatus
	Initialize() driver.CommandStatus
	Finalize() driver.CommandStatus
	Message(level int, title, text string, showcancel bool) driver.CommandStatus
	GetOutputFile() ([]byte, driver.CommandStatus)
	GetCapabilities() (LHProperties, driver.CommandStatus)
}

type LowLevelLiquidhandlingDriver interface {
	LiquidhandlingDriver
	//Move move the head to the given position
	//slices deckposition, wellcoords, reference, offsetX,Y,Z and plate_type should be
	//equal in length to the number of channels on the adaptor. Some elements can be nil or "" to signal
	//that the location of these channels are not specified and can be moved to anywhere compatible with robot geometry
	//deckposition: the name of the location on the deck
	//plate_type: the type of plate which should be present there
	//wellcoords: the well that each adaptor should line up with (e.g. in "A1" format),
	//reference: the position of the well/tip to align to: 0 = well bottom, 1 = well top, 2 = liquid level (undefined for tips)
	//offsetX,Y,Z: a relative offset from this position in mm
	//head: identifies which head should be moved in a multi-head system
	Move(deckposition []string, wellcoords []string, reference []int, offsetX, offsetY, offsetZ []float64, plate_type []string, head int) driver.CommandStatus
	//Aspirate suck up liquid into the loaded tips
	//volume, overstroke, platetype, what and llf must be equal in length to the number of channels in the adaptor
	//volume: amount to aspirate in ul
	//overstroke:
	//head: which head to use
	//multi: the number of channels given an instruction
	//platetype: the type of plate which should be present below the adaptor
	//what: liquid class name
	//llf: liquidlevelfollow attempt to follow liquid surface
	Aspirate(volume []float64, overstroke []bool, head int, multi int, platetype []string, what []string, llf []bool) driver.CommandStatus
	//Dispense eject liquid from the loaded tips
	//volume, blowout, platetype, what and llf must be equal in length to the number of channels in the adaptor
	//volume: amount to dispense in ul
	//blowout: dispense extra to attempt to remove droplets
	//head: which head to use
	//multi: the number of channels given an instruction
	//platetype: the type of plate which should be present below the adaptor
	//what: liquid class name
	//llf: liquidlevelfollow attempt to follow liquid surface
	Dispense(volume []float64, blowout []bool, head int, multi int, platetype []string, what []string, llf []bool) driver.CommandStatus
	//LoadTips add tips to the given channels
	//channels: list of which channels should end up with tips on them. values of platetype, position, well that aren't given
	//in channels can be left as ""
	//head: the head to use
	//multi: the number of channels given an instruction
	//platetype: the type of plate below each channel
	//position: the name of the deck position below the channel
	//well: the well below the adaptor channel
	LoadTips(channels []int, head, multi int, platetype, position, well []string) driver.CommandStatus
	//UnloadTips remove tips from the given channels
	//channels: list of which channels should have tips removed from them. values of platetype, position, well that aren't given
	//in channels can be left as ""
	//head: the head to use
	//multi: the number of channels given an instruction
	//platetype: the type of plate below each channel
	//position: the name of the deck position below the channel
	//well: the well below the adaptor channel
	UnloadTips(channels []int, head, multi int, platetype, position, well []string) driver.CommandStatus
	//SetPipetteSpeed set the rate of aspirate and dispense commands
	//non-independent heads can only have the same rate for each channel
	//rate units of ml/min
	SetPipetteSpeed(head, channel int, rate float64) driver.CommandStatus
	//SetDriveSpeed set the speed with which the robot head moves
	//units unknown...
	SetDriveSpeed(drive string, rate float64) driver.CommandStatus
	Wait(time float64) driver.CommandStatus
	//Mix pipette up and down
	Mix(head int, volume []float64, platetype []string, cycles []int, multi int, what []string, blowout []bool) driver.CommandStatus
	ResetPistons(head, channel int) driver.CommandStatus
	UpdateMetaData(props *LHProperties) driver.CommandStatus
}

// should be named UnimplementedLiquidHandlingDriver
type ExtendedLiquidhandlingDriver interface {
	LiquidhandlingDriver
	SetPositionState(position string, state driver.PositionState) driver.CommandStatus
	GetCurrentPosition(head int) (string, driver.CommandStatus)
	GetPositionState(position string) (string, driver.CommandStatus)
	GetHeadState(head int) (string, driver.CommandStatus)
	GetStatus() (driver.Status, driver.CommandStatus)
	UnloadHead(param int) driver.CommandStatus
	LoadHead(param int) driver.CommandStatus
	LightsOn() driver.CommandStatus
	LightsOff() driver.CommandStatus
	LoadAdaptor(param int) driver.CommandStatus
	//MoveRaw move head to the exact location
	MoveRaw(head int, x, y, z float64) driver.CommandStatus
	UnloadAdaptor(param int) driver.CommandStatus
	Open() driver.CommandStatus
	Close() driver.CommandStatus
}

type HighLevelLiquidhandlingDriver interface {
	LiquidhandlingDriver
	Transfer(what, platefrom, wellfrom, plateto, wellto []string, volume []float64) driver.CommandStatus
}
