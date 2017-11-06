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
	"testing"

	"github.com/kylelemons/godebug/pretty"
)

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
	//	LHWELL{
	//		ID        : 15cf94b7-ae06-443d-bc9a-9aadc30790fd,
	//		Inst      : ,
	//		Plateinst : ,
	//		Plateid   : ,
	//		Platetype : Gilson20Tipbox,
	//		Crds      : A1,
	//		Vol       : 20,
	//		Vunit     : ul,
	//		WContents : [],
	//		Rvol      : 1,
	//		Currvol   : 0,
	//		WShape    : &{cylinder mm 7.3 7.3 51.2},
	//	Bottom    : 0,
	//	Xdim      : 7.3,
	//	Ydim      : 7.3,
	//	Zdim      : 46,
	//	Bottomh   : 0,
	//	Dunit     : mm,
	//	Extra     : map[InnerL:5.5 InnerW:5.5 Tipeffectiveheight:34.6],
	//	Plate     : <nil>,
	//}

	wellExtra := make(map[string]interface{}, 0)
	lhwell := LHWell{
		"15cf94b7-ae06-443d-bc9a-9aadc30790fd",
		"",
		"",
		"",
		"Gilson20Tipbox",
		"A1",
		20,
		"ul",
		NewLHComponent(),
		1.0,
		&Shape{
			"cylinder",
			"mm",
			7.3,
			7.3,
			51.2,
		},
		0,
		7.3,
		7.3,
		46,
		0,
		"mm",
		wellExtra,
		nil,
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

	if reflect.DeepEqual(lhwell, dest) != true {
		fmt.Println(pretty.Compare(lhwell, dest))
		t.Fatal("Initial well and dest well differ")
	}
}

func TestSerializeLHPlate_1(t *testing.T) {
	//from make_plate_library
	swshp := NewShape("box", "mm", 8.2, 8.2, 41.3)
	welltype := NewLHWell("DSW96", "", "", "ul", 2000, 25, swshp, 3, 8.2, 8.2, 41.3, 4.7, "mm")
	plate := NewLHPlate("DSW96", "Unknown", 8, 12, 44.1, "mm", welltype, 9, 9, 0.0, 0.0, 0.0)

	enc, err := json.Marshal(plate)
	if err != nil {
		t.Fatal(err)
	}
	var outPlate LHPlate
	err = json.Unmarshal(enc, &outPlate)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(plate, outPlate) {
		if !reflect.DeepEqual(plate.Welltype, outPlate.Welltype) {
			fmt.Sprintf("It's the well type: ", plate.Welltype, " ", outPlate.Welltype)
		}

		for crds, _ := range plate.Wellcoords {
			w1 := plate.Wellcoords[crds]
			w2 := outPlate.Wellcoords[crds]

			if !reflect.DeepEqual(w1, w2) {
				fmt.Sprintf("It's at well %s\n", crds)
			}
		}

	}
}

// entity is now greatly stripped down
/*
func TestSerializeLHPlateGenericEntity(t *testing.T) {
	plate := LHPlate{}
		location := NewLocation("somewhere", 1, NewShape(
			"box", "m",
			0, 1, 1,
		))
			ge := GenericEntity{
				NewGenericSolid("water", "box"),
				location.(*ConcreteLocation),
			}
			plate.GenericEntity = &ge
	enc, err := json.Marshal(plate)
	if err != nil {
		t.Fatal(err)
	}
	var outPlate LHPlate
	err = json.Unmarshal(enc, &outPlate)
	if err != nil {
		t.Fatal(err)
	}
	if reflect.DeepEqual(plate, outPlate) {
		fmt.Println(pretty.Compare(plate, outPlate))
		t.Fatal("input plate and out plate do not differ")
	}
}
*/
