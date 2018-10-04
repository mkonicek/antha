// execution/uuid.go: Part of the Antha language
// Copyright (C) 2014 the Antha authors. All rights reserved.
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

package wtype

import (
	//"github.com/dustinkirkland/golang-petname"
	"github.com/twinj/uuid"
)

// this package wraps the uuid library appropriately
// by generating a V4 UUID
func GetUUID() string {
	return uuid.NewV4().String()
}

/*
// for debugging this can be useful
func GetUUID() string {
	return petname.Generate(2, "_")
}

func GetRUID(k int) string {
	return petname.Generate(k, "_")
}
*/

func NewUUID() string {
	return GetUUID()
}
