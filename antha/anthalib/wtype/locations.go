// /anthalib/wtype/locations.go: Part of the Antha language
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

package wtype

import (
	"fmt"
	"strconv"

	"github.com/antha-lang/antha/laboratory/effects/id"
)

// base types required for movement

// a Location is somewhere something can be
// it can recursively contain other Locations
type Location interface {
	Location_ID() string
	Location_Name() string
	Positions() []Location
	Container() Location
	Shape() *Shape
}

type ConcreteLocation struct {
	ID   string
	Inst string
	Name string
	Psns []*ConcreteLocation
	Cntr *ConcreteLocation
	Shap *Shape
}

func (cl *ConcreteLocation) Location_ID() string {
	return cl.ID
}
func (cl *ConcreteLocation) Location_Name() string {
	return cl.Name
}
func (cl *ConcreteLocation) Positions() []Location {
	s := make([]Location, len(cl.Psns))
	for i, v := range cl.Psns {
		s[i] = Location(v)
	}
	return s
}
func (cl *ConcreteLocation) Container() Location {
	return cl.Cntr
}
func (cl *ConcreteLocation) Shape() *Shape {
	return cl.Shap
}

func NewLocation(idGen *id.IDGenerator, name string, nPositions int, shape *Shape) Location { //TODO only in particular cases should the inner locations be populated
	positions := make([]*ConcreteLocation, nPositions)
	l := ConcreteLocation{idGen.NextID(), "", "", positions, nil, shape}
	l.Cntr = &l
	if nPositions > 0 {
		for i := 0; i < nPositions; i++ {
			l.Psns[i] = NewLocation(idGen, name+"_position_"+strconv.Itoa(i+1), 0, shape).(*ConcreteLocation)
			l.Psns[i].Cntr = &l
		}
	}
	return &l
}

// defines when two locations are the same
// level = 0 means only the exact same place
// level = 1 means they share a parent etc
func SameLocation(l, m Location, level int) bool {
	if level < 0 {
		panic(fmt.Sprintf("SameLocation: level parameter %d makes no sense", level))
	} else if level == 0 {
		if l.Location_ID() == m.Location_ID() {
			return true
		} else {
			return false
		}
	} else {
		//return SameLocation(l.Container(), m.Container(), level-1)
		firstMap := buildContainerMap(l)
		for id := range firstMap {
			second := buildContainerMap(m)
			if exists, res := second[id]; exists && res {
				return true
			}
		}
	}
	return false
}

func buildContainerMap(l Location) map[string]bool {
	ret := make(map[string]bool)
	pl := l
	for {
		ret[pl.Location_ID()] = true
		if pl.Container() == pl {
			break
		}
		pl = l.Container()
	}
	return ret
}

type Movable interface{}   //Entity
type Container interface{} //deferred to the shape
