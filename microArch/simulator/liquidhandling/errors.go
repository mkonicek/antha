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
	"sort"
	"strings"

	"github.com/antha-lang/antha/antha/anthalib/wtype"
	driver "github.com/antha-lang/antha/microArch/driver/liquidhandling"
	"github.com/antha-lang/antha/microArch/simulator"
)

type LiquidhandlingError interface {
	simulator.SimulationError
	Instruction() driver.TerminalRobotInstruction
	InstructionIndex() int
}

type DetailedLHError interface {
	GetStateAtError() string
}

type mutableLHError interface {
	setInstruction(int, driver.TerminalRobotInstruction)
}

type GenericError struct {
	severity         simulator.ErrorSeverity
	message          string
	instruction      driver.TerminalRobotInstruction
	instructionIndex int
	stateAtError     string
}

func NewGenericError(state *RobotState, severity simulator.ErrorSeverity, message string) LiquidhandlingError {
	return &GenericError{
		severity:     severity,
		message:      message,
		stateAtError: state.SummariseState(nil),
	}
}

func NewGenericErrorf(state *RobotState, severity simulator.ErrorSeverity, format string, a ...interface{}) LiquidhandlingError {
	return NewGenericError(state, severity, fmt.Sprintf(format, a...))
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

func (self *GenericError) GetStateAtError() string {
	return self.stateAtError
}

func (self *GenericError) Error() string {
	return fmt.Sprintf("(%v) %s[%d]: %s",
		self.severity,
		driver.HumanInstructionName(self.instruction),
		self.instructionIndex,
		self.message)
}

func (self *GenericError) setInstruction(index int, ins driver.TerminalRobotInstruction) {
	self.instruction = ins
	self.instructionIndex = index
}

//CollisionError generated when a physical collision occurs
type CollisionError struct {
	description       string
	channelsColliding map[int][]int //maps adaptors to a list of channels involved in collision
	objectsColliding  []wtype.LHObject
	instruction       driver.TerminalRobotInstruction
	instructionIndex  int
	stateAtError      string
}

//NewCollisionError make a new collision
func NewCollisionError(state *RobotState, channelsColliding map[int][]int, objectsColliding []wtype.LHObject) *CollisionError {
	return &CollisionError{
		channelsColliding: channelsColliding,
		objectsColliding:  objectsColliding,
		stateAtError:      state.SummariseState(channelsColliding),
	}
}

func (self *CollisionError) Severity() simulator.ErrorSeverity {
	return simulator.SeverityError
}

func (self *CollisionError) Instruction() driver.TerminalRobotInstruction {
	return self.instruction
}

func (self *CollisionError) InstructionIndex() int {
	return self.instructionIndex
}

func (self *CollisionError) Error() string {
	return fmt.Sprintf("(%v) %s[%d]: %s: collision detected: %s",
		self.Severity(),
		driver.HumanInstructionName(self.instruction),
		self.InstructionIndex(),
		self.InstructionDescription(),
		self.CollisionDescription())
}

func (self *CollisionError) setInstruction(index int, ins driver.TerminalRobotInstruction) {
	self.instruction = ins
	self.instructionIndex = index
}

func (self *CollisionError) InstructionDescription() string {
	return self.description
}

func (self *CollisionError) SetInstructionDescription(d string) {
	self.description = d
}

func (self *CollisionError) GetStateAtError() string {
	return self.stateAtError
}

func (self *CollisionError) CollisionDescription() string {

	//list adaptors in order for consistent errors
	adaptorIndexes := make([]int, 0, len(self.channelsColliding))
	for i := range self.channelsColliding {
		adaptorIndexes = append(adaptorIndexes, i)
	}
	sort.Ints(adaptorIndexes)

	adaptorStrings := make([]string, 0, len(self.channelsColliding))
	for _, adaptorIndex := range adaptorIndexes {
		adaptorStrings = append(adaptorStrings, fmt.Sprintf("head %d %s", adaptorIndex, summariseChannels(self.channelsColliding[adaptorIndex])))
	}

	//group objects by parent
	parentMap := make(map[wtype.LHObject][]wtype.LHObject, len(self.objectsColliding))
	for _, object := range self.objectsColliding {
		p := object.GetParent()
		if _, ok := parentMap[p]; !ok {
			parentMap[p] = make([]wtype.LHObject, 0, len(self.objectsColliding))
		}
		parentMap[p] = append(parentMap[p], object)
	}

	objectStrings := make([]string, 0, len(self.objectsColliding))
	for parent, children := range parentMap {
		deck := wtype.GetObjectRoot(parent).(*wtype.LHDeck)

		//if the parent is addressable, refer to the children compactly using their addresses
		var s string
		if addr, ok := parent.(wtype.Addressable); ok {
			wellcoords := make([]wtype.WellCoords, 0, len(children))
			for _, child := range children {
				pos := child.GetPosition().Add(child.GetSize().Multiply(0.5))
				wc, _ := addr.CoordsToWellCoords(pos)
				wellcoords = append(wellcoords, wc)
			}
			//WellCoordArrayRow sorts by col then row
			sort.Sort(wtype.WellCoordArrayRow(wellcoords))

			s = fmt.Sprintf("%s %s@%s at position %s", pluralClassOf(children[0], len(wellcoords)), wtype.HumanizeWellCoords(wellcoords), wtype.NameOf(parent), deck.GetSlotContaining(parent))
			objectStrings = append(objectStrings, s)
		} else {
			for _, child := range children {
				s = fmt.Sprintf("%s \"%s\" of type %s", wtype.ClassOf(child), wtype.NameOf(child), wtype.TypeOf(child))
				if pos := deck.GetSlotContaining(child); pos != "" {
					s += fmt.Sprintf(" at position %s", pos)
				}
				objectStrings = append(objectStrings, s)
			}
		}
	}

	return fmt.Sprintf("%s and %s", strings.Join(adaptorStrings, " and "), strings.Join(objectStrings, " and "))

}
