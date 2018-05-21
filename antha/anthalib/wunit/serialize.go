// /anthalib/wunit/serialize.go: Part of the Antha language
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

package wunit

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

type stringer interface {
	String() string
}

func marshal(cm *ConcreteMeasurement) ([]byte, error) {
	return json.Marshal(cm)
}

func unmarshal(b []byte) (value float64, unit string, err error) {
	var cm *ConcreteMeasurement

	if err = json.Unmarshal(b, &cm); err != nil {
		// support legacy marshal/unmarshal using string interface
		return unmarshalBad(b)
	}

	if cm != nil {
		value = cm.Mvalue
		if cm.Munit != nil {
			unit = cm.Munit.StrSymbol + cm.Munit.StrName
		}
	}
	return
}

func marshalBad(x stringer) ([]byte, error) {
	var s *string
	if x != nil {
		r := x.String()
		s = &r
	}
	return json.Marshal(s)
}

func unmarshalBad(b []byte) (value float64, unit string, err error) {
	var s *string
	if err = json.Unmarshal(b, &s); err != nil {
		return
	} else if s == nil {
		return
	}
	if _, err = fmt.Fscanf(strings.NewReader(*s), "%e%s", &value, &unit); err != nil && err == io.EOF {
		unit = ""
		if _, err = fmt.Fscanf(strings.NewReader(*s), "%e", &value); err != nil {
			return
		}
	}
	return
}

func (m Volume) MarshalJSON() ([]byte, error) {
	return marshal(m.ConcreteMeasurement)
}

func (m *Volume) UnmarshalJSON(b []byte) error {
	if value, unit, err := unmarshal(b); err != nil {
		return err
	} else if unit != "" {
		*m = NewVolume(value, unit)
	} else {
		cm := ConcreteMeasurement{0, nil}
		*m = Volume{&cm}
	}
	return nil
}

func (m Temperature) MarshalJSON() ([]byte, error) {
	return marshal(m.ConcreteMeasurement)
}

func (m *Temperature) UnmarshalJSON(b []byte) error {
	if value, unit, err := unmarshal(b); err != nil {
		return err
	} else if unit != "" {
		*m = NewTemperature(value, unit)
	} else {
		cm := ConcreteMeasurement{0, nil}
		*m = Temperature{&cm}
	}
	return nil
}

func (m Concentration) MarshalJSON() ([]byte, error) {
	return marshal(m.ConcreteMeasurement)
}

func (m *Concentration) UnmarshalJSON(b []byte) error {
	if value, unit, err := unmarshal(b); err != nil {
		return err
	} else if unit != "" {
		*m = NewConcentration(value, unit)
	} else {
		cm := ConcreteMeasurement{0, nil}
		*m = Concentration{&cm}
	}
	return nil
}

func (m Time) MarshalJSON() ([]byte, error) {
	return marshal(m.ConcreteMeasurement)
}

func (m *Time) UnmarshalJSON(b []byte) error {
	if value, unit, err := unmarshal(b); err != nil {
		return err
	} else if unit != "" {
		*m = NewTime(value, unit)
	} else {
		cm := ConcreteMeasurement{0, nil}
		*m = Time{&cm}
	}
	return nil
}

func (m Density) MarshalJSON() ([]byte, error) {
	return marshal(m.ConcreteMeasurement)
}

func (m *Density) UnmarshalJSON(b []byte) error {
	if value, unit, err := unmarshal(b); err != nil {
		return err
	} else if unit != "" {
		*m = NewDensity(value, unit)
	} else {
		cm := ConcreteMeasurement{0, nil}
		*m = Density{&cm}
	}
	return nil
}

func (m Mass) MarshalJSON() ([]byte, error) {
	return marshal(m.ConcreteMeasurement)
}

func (m *Mass) UnmarshalJSON(b []byte) error {
	if value, unit, err := unmarshal(b); err != nil {
		return err
	} else if unit != "" {
		*m = NewMass(value, unit)
	} else {
		cm := ConcreteMeasurement{0, nil}
		*m = Mass{&cm}

	}
	return nil
}

func (m FlowRate) MarshalJSON() ([]byte, error) {
	return marshal(m.ConcreteMeasurement)
}

func (m *FlowRate) UnmarshalJSON(b []byte) error {
	if value, unit, err := unmarshal(b); err != nil {
		return err
	} else if unit != "" {
		*m = NewFlowRate(value, unit)
	} else {
		cm := ConcreteMeasurement{0, nil}
		*m = FlowRate{&cm}

	}
	return nil
}

func (m Moles) MarshalJSON() ([]byte, error) {
	return marshal(m.ConcreteMeasurement)
}

func (m *Moles) UnmarshalJSON(b []byte) error {
	if value, unit, err := unmarshal(b); err != nil {
		return err
	} else if unit != "" {
		*m = NewMoles(value, unit)
	} else {
		cm := ConcreteMeasurement{0, nil}
		*m = Moles{&cm}
	}
	return nil
}

func (m Pressure) MarshalJSON() ([]byte, error) {
	return marshal(m.ConcreteMeasurement)
}

func (m *Pressure) UnmarshalJSON(b []byte) error {
	if value, unit, err := unmarshal(b); err != nil {
		return err
	} else if unit != "" {
		*m = NewPressure(value, unit)
	} else {
		cm := ConcreteMeasurement{0, nil}
		*m = Pressure{&cm}
	}
	return nil
}

func (m Length) MarshalJSON() ([]byte, error) {
	return marshal(m.ConcreteMeasurement)
}

func (m *Length) UnmarshalJSON(b []byte) error {
	if value, unit, err := unmarshal(b); err != nil {
		return err
	} else if unit != "" {
		*m = NewLength(value, unit)
	} else {
		cm := ConcreteMeasurement{0, nil}
		*m = Length{&cm}
	}
	return nil
}

func (m Area) MarshalJSON() ([]byte, error) {
	return marshal(m.ConcreteMeasurement)
}

func (m *Area) UnmarshalJSON(b []byte) error {
	if value, unit, err := unmarshal(b); err != nil {
		return err
	} else if unit != "" {
		*m = NewArea(value, unit)
	} else {
		cm := ConcreteMeasurement{0, nil}
		*m = Area{&cm}
	}
	return nil
}

func (m Angle) MarshalJSON() ([]byte, error) {
	return marshal(m.ConcreteMeasurement)
}

func (m *Angle) UnmarshalJSON(b []byte) error {
	if value, unit, err := unmarshal(b); err != nil {
		return err
	} else if unit != "" {
		*m = NewAngle(value, unit)
	} else {
		cm := ConcreteMeasurement{0, nil}
		*m = Angle{&cm}
	}
	return nil
}

func (m Energy) MarshalJSON() ([]byte, error) {
	return marshal(m.ConcreteMeasurement)
}

func (m *Energy) UnmarshalJSON(b []byte) error {
	if value, unit, err := unmarshal(b); err != nil {
		return err
	} else if unit != "" {
		*m = NewEnergy(value, unit)
	} else {
		cm := ConcreteMeasurement{0, nil}
		*m = Energy{&cm}
	}
	return nil
}

func (m Force) MarshalJSON() ([]byte, error) {
	return marshal(m.ConcreteMeasurement)
}

func (m *Force) UnmarshalJSON(b []byte) error {
	if value, unit, err := unmarshal(b); err != nil {
		return err
	} else if unit != "" {
		*m = NewForce(value, unit)
	} else {
		cm := ConcreteMeasurement{0, nil}
		*m = Force{&cm}
	}
	return nil
}

func (m Velocity) MarshalJSON() ([]byte, error) {
	return marshal(m.ConcreteMeasurement)
}

func (m *Velocity) UnmarshalJSON(b []byte) error {
	if value, unit, err := unmarshal(b); err != nil {
		return err
	} else if unit != "" {
		*m = NewVelocity(value, unit)
	} else {
		cm := ConcreteMeasurement{0, nil}
		*m = Velocity{&cm}
	}
	return nil
}

func (m Rate) MarshalJSON() ([]byte, error) {
	return marshal(m.ConcreteMeasurement)
}

func (m *Rate) UnmarshalJSON(b []byte) error {
	if value, unit, err := unmarshal(b); err != nil {
		return err
	} else if unit != "" {
		*m, err = NewRate(value, unit)
		if err != nil {
			return err
		}
	} else {
		cm := ConcreteMeasurement{0, nil}
		*m = Rate{&cm}
	}
	return nil
}
