// antha/anthalib/wtype/serialize_test.go: Part of the Antha language
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
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/go-test/deep"
)

func TestMarshalDeckObject(t *testing.T) {
	objects := []LHObject{
		makeplatefortest(),
		maketroughfortest(),
		makeTipboxForTest(),
		makeTipwasteForTest(),
	}

	for _, obj := range objects {
		t.Run(fmt.Sprintf("%T called %s", obj, NameOf(obj)), func(t *testing.T) {
			if data, err := MarshalDeckObject(obj); err != nil {
				t.Error(err)
			} else if after, err := UnmarshalDeckObject(data); err != nil {
				t.Error(err)
			} else if !reflect.DeepEqual(obj, after) {
				t.Errorf("serialisation changed object: \ne: %#v,\ng: %#v", obj, after)
			}
		})
	}
}

func TestDeserializeLHSolution(t *testing.T) {
	str := `{"ID":"","BlockID":{"ThreadID":"","OutputCount":0},"Inst":"","SName":"","Order":0,"Components":null,"ContainerType":"","Welladdress":"","Plateaddress":"","PlateID":"","Platetype":"","Vol":0,"Type":"","Conc":0,"Tvol":0,"Majorlayoutgroup":0,"Minorlayoutgroup":0}`
	var sol LHSolution
	err := json.Unmarshal([]byte(str), &sol)
	if err != nil {
		t.Fatal(err)
	}
}

/*
func TestDeserializeGenericPhysical(t *testing.T) {
	str := `{"Iname":"","Imp":"0.000 ","Ibp":"0.000 ","Ishc":{"Mvalue":0,"Munit":null},"Myname":"","Mymass":"0.000 ","Myvol":"0.000 ","Mytemp":"0.000 "}`
	var gp GenericPhysical
	err := json.Unmarshal([]byte(str), &gp)
	if err != nil {
		t.Fatal(err)
	}
}

func TestIdempotentGenericPhysical(t *testing.T) {
	var gp LHSolution
	bs, err := json.Marshal(&gp)
	if err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(bs, &gp); err != nil {
		t.Fatal(err)
	}
}

func TestDeserializeGenericMatter(t *testing.T) {
	str := `{"Iname":"","Imp":"0.000 ","Ibp":"0.000 ","Ishc":{"Mvalue":0,"Munit":null}}`
	var gm GenericMatter
	err := json.Unmarshal([]byte(str), &gm)
	if err != nil {
		t.Fatal(err)
	}
}
*/
func TestLHWellSerialize(t *testing.T) {

	wellExtra := make(map[string]interface{})
	lhwell := LHWell{
		ID:        "15cf94b7-ae06-443d-bc9a-9aadc30790fd",
		Inst:      "",
		Crds:      MakeWellCoords("A1"),
		MaxVol:    20,
		WContents: NewLHComponent(),
		Rvol:      1.0,
		WShape: &Shape{
			Type:       CylinderShape,
			LengthUnit: "mm",
			H:          7.3,
			W:          7.3,
			D:          51.2,
		},
		Bottom:  FlatWellBottom,
		Bounds:  BBox{Coordinates3D{}, Coordinates3D{7.3, 7.3, 51.2}},
		Bottomh: 46,
		Extra:   wellExtra,
	}

	j, err := json.Marshal(lhwell)
	if err != nil {
		t.Fatal(err)
	}
	var dest LHWell

	err = json.Unmarshal(j, &dest)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(lhwell, dest) {
		t.Fatalf("Initial well (%+v) and dest well (%+v) differ. Differences are: %s", lhwell.WContents, dest.WContents, strings.Join(deep.Equal(lhwell, dest), "\n"))
	}

	if !CylinderShape.Equals(dest.WShape.Type) {
		t.Errorf(`well.WShape.Type changed: "%s" (@%p) -> "%s" (%p)`, lhwell.WShape.Type, lhwell.WShape.Type, dest.WShape.Type, dest.WShape.Type)
	}
}
