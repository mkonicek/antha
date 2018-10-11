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
	"github.com/pkg/errors"
)

type stringer interface {
	String() string
}

func marshal(x stringer) ([]byte, error) {
	var s *string
	if x != nil {
		r := x.String()
		s = &r
	}
	return json.Marshal(s)
}

func unmarshal(b []byte) (float64, string, error) {
	var s *string
	if err := json.Unmarshal(b, &s); err != nil || s == nil || *s == "" {
		return 0.0, "", err
	}

	if value, unit := extractFloat(*s); len(unit) == len(*s) {
		return 0.0, "", errors.Errorf("couldn't parse float from %q", *s)
	} else {
		return value, unit, nil
	}
}

func (m Volume) MarshalJSON() ([]byte, error) {
	return marshal(m)
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
	return marshal(m)
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
	return marshal(m)
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
	return marshal(m)
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
	return marshal(m)
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
	return marshal(m)
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
	return marshal(m)
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
	return marshal(m)
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
	return marshal(m)
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
	return marshal(m)
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
	return marshal(m)
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
	return marshal(m)
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
	return marshal(m)
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
	return marshal(m)
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
	return marshal(m)
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
	return marshal(m)
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
