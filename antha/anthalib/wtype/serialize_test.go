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

	wellExtra := make(map[string]interface{})
	lhwell := LHWell{
		"15cf94b7-ae06-443d-bc9a-9aadc30790fd",
		"",
		MakeWellCoords("A1"),
		20,
		NewLHComponent(),
		1.0,
		&Shape{
			"cylinder",
			"mm",
			7.3,
			7.3,
			51.2,
		},
		FlatWellBottom,
		BBox{Coordinates{}, Coordinates{7.3, 7.3, 51.2}},
		46,
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

	if !reflect.DeepEqual(lhwell, dest) {
		fmt.Println(pretty.Compare(lhwell, dest))
		t.Fatal(fmt.Sprintf("Initial well (%+v) and dest well (%+v) differ", lhwell.WContents, dest.WContents))
	}
}
