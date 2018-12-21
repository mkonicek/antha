// /anthalib/liquidhandling/simulator/simpleliquidhandler.go: Part of the Antha language
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

package simpleliquidhandler

import (
	"time"
)

var (
	DEFAULT_SLEEP_TIME = 1 * time.Second
)

type SimpleLiquidHandler struct {
}

type LHCommandStatus struct {
	value    bool
	errValue int
	code     string
}

type LHStatus interface{}
type LHProperties interface{}
type LHPlate interface{}
type LHPositionState interface{}

func (lh *SimpleLiquidHandler) Move(deckposition string, wellcoords string, reference int, offsetX, offsetY, offsetZ float64, plate_type, head int) LHCommandStatus {
	time.Sleep(DEFAULT_SLEEP_TIME)
	return LHCommandStatus{true, 0, "OK"}
}
func (lh *SimpleLiquidHandler) MoveExplicit(deckposition string, wellcoords string, reference int, offsetX, offsetY, offsetZ float64, plate_type *LHPlate, head int) LHCommandStatus {
	time.Sleep(DEFAULT_SLEEP_TIME)
	return LHCommandStatus{true, 0, "OK"}
}
func (lh *SimpleLiquidHandler) MoveRaw(x, y, z float64) LHCommandStatus {
	time.Sleep(DEFAULT_SLEEP_TIME)
	return LHCommandStatus{true, 0, "OK"}
}
func (lh *SimpleLiquidHandler) Aspirate(volume float64, overstroke bool, head int, multi int) LHCommandStatus {
	time.Sleep(DEFAULT_SLEEP_TIME)
	return LHCommandStatus{true, 0, "OK"}
}
func (lh *SimpleLiquidHandler) Dispense(volume float64, blowout bool, head int, multi int) LHCommandStatus {
	time.Sleep(DEFAULT_SLEEP_TIME)
	return LHCommandStatus{true, 0, "OK"}
}
func (lh *SimpleLiquidHandler) LoadTips(head, multi int) LHCommandStatus {
	time.Sleep(DEFAULT_SLEEP_TIME)
	return LHCommandStatus{true, 0, "OK"}
}
func (lh *SimpleLiquidHandler) UnloadTips(head, multi int) LHCommandStatus {
	time.Sleep(DEFAULT_SLEEP_TIME)
	return LHCommandStatus{true, 0, "OK"}
}
func (lh *SimpleLiquidHandler) SetPipetteSpeed(rate float64) {

}
func (lh *SimpleLiquidHandler) SetDriveSpeed(drive string, rate float64) LHCommandStatus {
	time.Sleep(DEFAULT_SLEEP_TIME)
	return LHCommandStatus{true, 0, "OK"}
}
func (lh *SimpleLiquidHandler) Stop() LHCommandStatus {
	time.Sleep(DEFAULT_SLEEP_TIME)
	return LHCommandStatus{true, 0, "OK"}
}
func (lh *SimpleLiquidHandler) Go() LHCommandStatus {
	time.Sleep(DEFAULT_SLEEP_TIME)
	return LHCommandStatus{true, 0, "OK"}
}
func (lh *SimpleLiquidHandler) Initialize() LHCommandStatus {
	time.Sleep(DEFAULT_SLEEP_TIME)
	return LHCommandStatus{true, 0, "OK"}
}
func (lh *SimpleLiquidHandler) Finalize() LHCommandStatus {
	time.Sleep(DEFAULT_SLEEP_TIME)
	return LHCommandStatus{true, 0, "OK"}
}
func (lh *SimpleLiquidHandler) SetPositionState(position string, state LHPositionState) LHCommandStatus {
	time.Sleep(DEFAULT_SLEEP_TIME)
	return LHCommandStatus{true, 0, "OK"}
}
func (lh *SimpleLiquidHandler) GetCapabilities() LHProperties {
	panic("not yet implemented")
}
func (lh *SimpleLiquidHandler) GetCurrentPosition(head int) (string, LHCommandStatus) {
	time.Sleep(DEFAULT_SLEEP_TIME)
	return "", LHCommandStatus{true, 0, "OK"}
}
func (lh *SimpleLiquidHandler) GetPositionState(position string) (string, LHCommandStatus) {
	time.Sleep(DEFAULT_SLEEP_TIME)
	return "", LHCommandStatus{true, 0, "OK"}
}
func (lh *SimpleLiquidHandler) GetHeadState(head int) (string, LHCommandStatus) {
	time.Sleep(DEFAULT_SLEEP_TIME)
	return "", LHCommandStatus{true, 0, "OK"}
}
func (lh *SimpleLiquidHandler) GetStatus() (LHStatus, LHCommandStatus) {
	panic("not yet implemented")
}
