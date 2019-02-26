// /anthalib/driver/types.go: Part of the Antha language
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

package driver

import (
	"fmt"
	"github.com/pkg/errors"
	"os"
)

type ErrorCode int

const (
	OK ErrorCode = iota
	ERR
	WRN
	NIM // Not implemented
)

var ecNames = map[ErrorCode]string{
	OK:  "ok",
	ERR: "error",
	WRN: "warning",
	NIM: "not implemented",
}

func (ec ErrorCode) String() string {
	if ret, ok := ecNames[ec]; ok {
		return ret
	}
	panic(fmt.Sprintf("unknown error code: %d", ec))
}

// CommandStatus the result of a robot instruction
type CommandStatus struct {
	ErrorCode
	Msg string
}

// CommandOk indicate that the function was successful
func CommandOk() CommandStatus {
	return CommandStatus{}
}

// CommandError indicate that a fatal error occurred
func CommandError(message string) CommandStatus {
	return CommandStatus{
		ErrorCode: ERR,
		Msg:       message,
	}
}

// CommandWarn indicate that the function completed successfully but that a warning was generated
func CommandWarn(message string) CommandStatus {
	return CommandStatus{
		ErrorCode: WRN,
		Msg:       message,
	}
}

// CommandNotImplemented indicates that the function is not implemented for this driver
func CommandNotImplemented(message string) CommandStatus {
	return CommandStatus{
		ErrorCode: NIM,
		Msg:       message,
	}
}

// Ok tests whether the function returned successfully
func (cs CommandStatus) Ok() bool {
	return cs.ErrorCode == OK
}

// Fatal returns true if a fatal error occurred
func (cs CommandStatus) Fatal() bool {
	return cs.ErrorCode != OK && cs.ErrorCode != WRN
}

// String returns a string representation of the CommandStatus
func (cs CommandStatus) String() string {
	if cs.Msg != "" {
		return fmt.Sprintf("%s: %s", cs.ErrorCode, cs.Msg)
	} else {
		return cs.ErrorCode.String()
	}
}

// GetError if a fatal error occurred, returns it. Otherwise returns nil.
// If a warning was returned, it is written to standard error
func (cs CommandStatus) GetError() error {
	if cs.Fatal() {
		return errors.New(cs.String())
	} else if cs.ErrorCode == WRN {
		fmt.Fprintf(os.Stderr, "driver warning: %s", cs.String()) // nolint
	}
	return nil
}

type Status map[string]interface{}
type PositionState map[string]interface{}
