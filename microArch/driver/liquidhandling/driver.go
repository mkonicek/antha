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
	Move(deckposition []string, wellcoords []string, reference []int, offsetX, offsetY, offsetZ []float64, plate_type []string, head int) driver.CommandStatus
	MoveRaw(head int, x, y, z float64) driver.CommandStatus
	Aspirate(volume []float64, overstroke []bool, head int, multi int, platetype []string, what []string, llf []bool) driver.CommandStatus
	Dispense(volume []float64, blowout []bool, head int, multi int, platetype []string, what []string, llf []bool) driver.CommandStatus
	LoadTips(channels []int, head, multi int, platetype, position, well []string) driver.CommandStatus
	UnloadTips(channels []int, head, multi int, platetype, position, well []string) driver.CommandStatus
	SetPipetteSpeed(head, channel int, rate float64) driver.CommandStatus
	SetDriveSpeed(drive string, rate float64) driver.CommandStatus
	Stop() driver.CommandStatus
	Go() driver.CommandStatus
	Initialize() driver.CommandStatus
	Finalize() driver.CommandStatus
	Wait(time float64) driver.CommandStatus
	Mix(head int, volume []float64, platetype []string, cycles []int, multi int, what []string, blowout []bool) driver.CommandStatus
	ResetPistons(head, channel int) driver.CommandStatus
	AddPlateTo(position string, plate interface{}, name string) driver.CommandStatus
	RemoveAllPlates() driver.CommandStatus
	RemovePlateAt(position string) driver.CommandStatus
	Message(level int, title, text string, showcancel bool) driver.CommandStatus
}

type ExtendedLiquidhandlingDriver interface {
	LiquidhandlingDriver
	SetPositionState(position string, state driver.PositionState) driver.CommandStatus
	GetCapabilities() (LHProperties, driver.CommandStatus)
	GetCurrentPosition(head int) (string, driver.CommandStatus)
	GetPositionState(position string) (string, driver.CommandStatus)
	GetHeadState(head int) (string, driver.CommandStatus)
	GetStatus() (driver.Status, driver.CommandStatus)
	UpdateMetaData(props *LHProperties) driver.CommandStatus
	UnloadHead(param int) driver.CommandStatus
	LoadHead(param int) driver.CommandStatus
	LightsOn() driver.CommandStatus
	LightsOff() driver.CommandStatus
	LoadAdaptor(param int) driver.CommandStatus
	UnloadAdaptor(param int) driver.CommandStatus
	// refactored into other interfaces?
	Open() driver.CommandStatus
	Close() driver.CommandStatus
	GetOutputFile() (string, driver.CommandStatus)
}
