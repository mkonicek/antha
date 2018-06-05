// /anthalib/simulator/liquidhandling/errors.go: Part of the Antha language
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
	"fmt"
	driver "github.com/antha-lang/antha/microArch/driver/liquidhandling"
	"github.com/antha-lang/antha/microArch/simulator"
)

type LiquidhandlingError interface {
	simulator.SimulationError
	Instruction() driver.TerminalRobotInstruction
	InstructionIndex() int
}

type mutableLHError interface {
	setInstruction(int, driver.TerminalRobotInstruction)
}

type GenericError struct {
	severity         simulator.ErrorSeverity
	message          string
	instruction      driver.TerminalRobotInstruction
	instructionIndex int
}

func NewGenericError(severity simulator.ErrorSeverity, message string) LiquidhandlingError {
	return &GenericError{
		severity: severity,
		message:  message,
	}
}

func NewGenericErrorf(severity simulator.ErrorSeverity, format string, a ...interface{}) LiquidhandlingError {
	return NewGenericError(severity, fmt.Sprintf(format, a...))
}

func (self *GenericError) Severity() simulator.ErrorSeverity {
	return self.severity
}

func (self *GenericError) Instruction() driver.TerminalRobotInstruction {
	return self.instruction
}

func (self *GenericError) InstructionIndex() int {
	return self.instructionIndex
}

func (self *GenericError) Error() string {
	return fmt.Sprintf("(%s) %s[%d]: %s",
		self.severity,
		driver.HumanInstructionName(self.instruction),
		self.instructionIndex,
		self.message)
}

func (self *GenericError) setInstruction(index int, ins driver.TerminalRobotInstruction) {
	self.instruction = ins
	self.instructionIndex = index
}
