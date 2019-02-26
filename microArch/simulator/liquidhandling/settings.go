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
	enable_tipload_override bool      //allow the adaptor to override the tip loading behaviour
	warn_auto_channels      Frequency //Display warnings for load/unload tips
	max_dispense_height     float64   //maximum height to dispense from in mm
	warnPipetteSpeed        Frequency //Raise warnings for pipette speed out of range
	warnLiquidType          Frequency //raise warnings when liquid types don't match
}

func DefaultSimulatorSettings() *SimulatorSettings {
	ss := SimulatorSettings{
		enable_tipbox_collision: true,
		enable_tipbox_check:     true,
		enable_tipload_override: true,
		warn_auto_channels:      WarnAlways,
		max_dispense_height:     5.,
		warnPipetteSpeed:        WarnAlways,
		warnLiquidType:          WarnNever,
	}
	return &ss
}

func (self *SimulatorSettings) IsTipboxCollisionEnabled() bool {
	return self.enable_tipbox_collision
}

func (self *SimulatorSettings) EnableTipboxCollision(b bool) {
	self.enable_tipbox_collision = b
}

func (self *SimulatorSettings) IsTipLoadingOverrideEnabled() bool {
	return self.enable_tipload_override
}

func (self *SimulatorSettings) EnableTipLoadingOverride(b bool) {
	self.enable_tipload_override = b
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

func (self *SimulatorSettings) IsPipetteSpeedWarningEnabled() bool {
	switch self.warnPipetteSpeed {
	case WarnAlways:
		return true
	case WarnOnce:
		self.warnPipetteSpeed = WarnNever
		return true
	}
	return false
}

func (self *SimulatorSettings) EnablePipetteSpeedWarning(f Frequency) {
	self.warnPipetteSpeed = f
}

func (self *SimulatorSettings) IsLiquidTypeWarningEnabled() bool {
	switch self.warnLiquidType {
	case WarnAlways:
		return true
	case WarnOnce:
		self.warnLiquidType = WarnNever
		return true
	}
	return false
}

func (self *SimulatorSettings) EnableLiquidTypeWarning(f Frequency) {
	self.warnLiquidType = f
}
