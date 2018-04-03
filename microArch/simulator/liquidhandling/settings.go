// /anthalib/simulator/liquidhandling/simulator_test.go: Part of the Antha language
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

type Frequency int

const (
	WarnNever Frequency = iota
	WarnOnce
	WarnAlways
)

type SimulatorSettings struct {
	enable_tipbox_collision bool      //Whether or not to complain if the head hits a tipbox
	enable_tipbox_check     bool      //detect tipboxes which are taller that the tips, and disable tipbox_collisions
	warn_auto_channels      Frequency //Display warnings for load/unload tips
	max_dispense_height     float64   //maximum height to dispense from in mm
}

func DefaultSimulatorSettings() *SimulatorSettings {
	ss := SimulatorSettings{
		true,
		true,
		WarnAlways,
		5.,
	}
	return &ss
}

func (self *SimulatorSettings) IsTipboxCollisionEnabled() bool {
	return self.enable_tipbox_collision
}

func (self *SimulatorSettings) EnableTipboxCollision(b bool) {
	self.enable_tipbox_collision = b
}

func (self *SimulatorSettings) IsTipboxCheckEnabled() bool {
	return self.enable_tipbox_check
}

func (self *SimulatorSettings) EnableTipboxCheck(b bool) {
	self.enable_tipbox_check = b
}

func (self *SimulatorSettings) IsAutoChannelWarningEnabled() bool {
	if self.warn_auto_channels == WarnAlways {
		return true
	} else if self.warn_auto_channels == WarnOnce {
		self.warn_auto_channels = WarnNever
		return true
	}
	return false
}

func (self *SimulatorSettings) EnableAutoChannelWarning(f Frequency) {
	self.warn_auto_channels = f
}

func (self *SimulatorSettings) MaxDispenseHeight() float64 {
	return self.max_dispense_height
}

func (self *SimulatorSettings) SetMaxDispenseHeight(f float64) {
	self.max_dispense_height = f
}
