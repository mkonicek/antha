// anthalib//liquidhandling/liquidhandling_test.go: Part of the Antha language
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
	"github.com/antha-lang/antha/antha/anthalib/mixer"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"math/rand"
	"testing"
)

func TestStockConcs(*testing.T) {
	names := []string{"tea", "milk", "sugar"}

	minrequired := make(map[string]float64, len(names))
	maxrequired := make(map[string]float64, len(names))
	Smax := make(map[string]float64, len(names))
	T := make(map[string]float64, len(names))
	vmin := 10.0

	for _, name := range names {
		r := rand.Float64() + 1.0
		r2 := rand.Float64() + 1.0
		r3 := rand.Float64() + 1.0

		minrequired[name] = r * r2 * 20.0
		maxrequired[name] = r * r2 * 30.0
		Smax[name] = r * r2 * r3 * 70.0
		T[name] = 100.0
	}

	cncs := choose_stock_concentrations(minrequired, maxrequired, Smax, vmin, T)
	cncs = cncs
	/*for k, v := range cncs {
		logger.Debug(fmt.Sprintln(k, " ", minrequired[k], " ", maxrequired[k], " ", T[k], " ", v))
	}*/
}

/*
func TestInputAssignments(t *testing.T) {
	lh := GetLiquidHandlerForTest()
	rq := GetLHRequestForTest()
	configure_request_simple(rq)

	fmt.Println("INPUT ASSIGNMENTS")
	for k, v := range rq.Input_assignments {
		fmt.Println(k, " ", v)
	}

	fmt.Println("INPUT SOLUTIONS")

	for k, v := range rq.Input_solutions {
		fmt.Println("\t", k, ":")
		for _, v2 := range v {
			fmt.Println("\t\t", v2.CName, " ", v2.Volume().ToString())
		}
	}

	fmt.Println("INPUT PLATES")

	for _, v := range rq.Input_plates {
		fmt.Println("\tPLATE: ", v.Name, " ", v.Type)
	}

	fmt.Println("INPUT ASSIGNMENTS")

	for k, v := range rq.Input_assignments {
		fmt.Println("\t", k, " ", v)
	}

	fmt.Println("INPUT VOLS SUPPLIED")

	for k, v := range rq.Input_vols_supplied {
		fmt.Println(k, " ", v.ToString())
	}

	fmt.Println("INPUT VOLS REQUIRED")

	for k, v := range rq.Input_vols_required {
		fmt.Println(k, " ", v.ToString())
	}

	fmt.Println("INPUT VOLS WANTING")
	for k, v := range rq.Input_vols_wanting {
		fmt.Println(k, " ", v.ToString())
	}


}
*/

func configure_request_simple(rq *LHRequest) {
	water := GetComponentForTest("water", wunit.NewVolume(100.0, "ul"))
	mmx := GetComponentForTest("mastermix_sapI", wunit.NewVolume(100.0, "ul"))
	part := GetComponentForTest("dna", wunit.NewVolume(50.0, "ul"))

	for k := 0; k < 9; k++ {
		ins := wtype.NewLHInstruction()
		ws := mixer.Sample(water, wunit.NewVolume(8.0, "ul"))
		mmxs := mixer.Sample(mmx, wunit.NewVolume(8.0, "ul"))
		ps := mixer.Sample(part, wunit.NewVolume(1.0, "ul"))

		ins.AddComponent(ws)
		ins.AddComponent(mmxs)
		ins.AddComponent(ps)
		ins.AddProduct(GetComponentForTest("water", wunit.NewVolume(17.0, "ul")))
		rq.Add_instruction(ins)
	}
}

func TestPlateReuse(t *testing.T) {
	lh := GetLiquidHandlerForTest()
	rq := GetLHRequestForTest()
	configure_request_simple(rq)
	rq.Input_platetypes = append(rq.Input_platetypes, GetPlateForTest())
	rq.Output_platetypes = append(rq.Output_platetypes, GetPlateForTest())

	rq.ConfigureYourself()

	err := lh.Plan(rq)

	if err != nil {
		t.Fatal(fmt.Sprint("Got an error planning with no inputs: ", err))
	}

	// reset the request
	rq = GetLHRequestForTest()
	configure_request_simple(rq)

	for _, plateid := range lh.Properties.PosLookup {
		if plateid == "" {
			continue
		}
		thing := lh.Properties.PlateLookup[plateid]

		plate, ok := thing.(*wtype.LHPlate)
		if !ok {
			continue
		}
		rq.Input_plates[plateid] = plate
	}
	rq.Input_platetypes = append(rq.Input_platetypes, GetPlateForTest())
	rq.Output_platetypes = append(rq.Output_platetypes, GetPlateForTest())

	rq.ConfigureYourself()

	lh = GetLiquidHandlerForTest()
	err = lh.Plan(rq)

	if err != nil {
		t.Fatal(fmt.Sprint("Got error resimulating: ", err))
	}

	// if we added nothing, input assignments should be empty

	if rq.NewComponentsAdded() {
		t.Fatal(fmt.Sprint("Resimulation failed: needed to add ", len(rq.Input_vols_wanting), " components"))
	}

	// now try a deliberate fail

	// reset the request again
	rq = GetLHRequestForTest()
	configure_request_simple(rq)

	for _, plateid := range lh.Properties.PosLookup {
		if plateid == "" {
			continue
		}
		thing := lh.Properties.PlateLookup[plateid]

		plate, ok := thing.(*wtype.LHPlate)
		if !ok {
			continue
		}
		for _, v := range plate.Wellcoords {
			if !v.Empty() {
				v.Remove(wunit.NewVolume(5.0, "ul"))
			}
		}

		rq.Input_plates[plateid] = plate
	}
	rq.Input_platetypes = append(rq.Input_platetypes, GetPlateForTest())
	rq.Output_platetypes = append(rq.Output_platetypes, GetPlateForTest())

	rq.ConfigureYourself()

	lh = GetLiquidHandlerForTest()
	err = lh.Plan(rq)

	if err != nil {
		t.Fatal(fmt.Sprint("Got error resimulating: ", err))
	}

	// this time we should have added some components again
	if len(rq.Input_assignments) != 3 {
		t.Fatal(fmt.Sprintf("Error resimulating, should have added 3 components, instead added ", len(rq.Input_assignments)))
	}

}

func TestEP3(t *testing.T) {
	fmt.Println("Test EP3")
	lh := GetLiquidHandlerForTest()
	lh.ExecutionPlanner = ExecutionPlanner3
	rq := GetLHRequestForTest()
	configure_request_simple(rq)
	rq.Input_platetypes = append(rq.Input_platetypes, GetPlateForTest())
	rq.Output_platetypes = append(rq.Output_platetypes, GetPlateForTest())

	rq.ConfigureYourself()
	err := lh.Plan(rq)

	if err != nil {
		t.Fatal(fmt.Sprint("Got planning error: ", err))
	}

}
