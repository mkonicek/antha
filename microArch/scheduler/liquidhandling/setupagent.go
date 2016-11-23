// liquidhandling/setupagent.go: Part of the Antha language
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

package liquidhandling

import (
	"fmt"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/microArch/driver/liquidhandling"
	"github.com/antha-lang/antha/microArch/factory"
	"strings"
)

// v2.0 should be another linear program - basically just want to optimize
// positioning in the face of constraints

// default setup agent
func BasicSetupAgent(request *LHRequest, params *liquidhandling.LHProperties) (*LHRequest, error) {
	// this is quite tricky and requires extensive interaction with the liquid handling
	// parameters

	// the principal question is how to define constraints on the system

	// I think this needs to remain tbd for now
	// instead we can rely on the preference system I already use

	plate_lookup := make(map[string]string, 5)
	//tip_lookup := make([]*wtype.LHTipbox, 0, 5)

	//	tip_preferences := params.Tip_preferences
	input_preferences := params.Input_preferences
	output_preferences := params.Output_preferences

	// how do we set the below?
	// we don't know how many tips we need until we generate
	// instructions; ditto input or output plates until we've done layout

	// input plates
	input_plates := request.Input_plates
	input_plate_order := request.Input_plate_order

	if len(input_plate_order) < len(input_plates) {
		input_plate_order = make([]string, 0, len(input_plates))
		/*
			x := 0
			for k, _ := range input_plates {
				input_plate_order[x] = k
				x += 1
			}

			sort.Strings(input_plate_order)
		*/

		for _, ass := range request.Input_assignments {
			for _, a := range ass {
				tx := strings.Split(a, ":")
				if !isInStrArr(tx[0], input_plate_order) {
					input_plate_order = append(input_plate_order, tx[0])
				}
			}
		}

		if len(input_plate_order) < len(input_plates) {
			for id, _ := range input_plates {
				if !isInStrArr(id, input_plate_order) {
					input_plate_order = append(input_plate_order, id)
				}
			}
		} else if len(input_plate_order) > len(input_plates) {
			return nil, wtype.LHError(wtype.LH_ERR_DIRE, fmt.Sprintf("Plate number inconsistency: %d != %d (here: %d)", len(input_plate_order), len(input_plates), 89))
		}

		request.Input_plate_order = input_plate_order
	}

	// output plates
	output_plates := request.Output_plates
	output_plate_order := request.Output_plate_order

	fmt.Println("OPA: ", output_plate_order)

	if len(output_plate_order) < len(output_plates) {
		output_plate_order = make([]string, 0, len(output_plates))
		/*
				x := 0
			for k, _ := range output_plates {
				fmt.Println("K: ", k)
			}
				sort.Strings(output_plate_order)
		*/
		// order them according to when they are first used
		for _, insID := range request.Output_order {
			ins := request.LHInstructions[insID]
			tx := strings.Split(ins.Result.Loc, ":")
			pa := tx[0]
			notInInputs := !isInStrArr(pa, input_plate_order)
			if !isInStrArr(pa, output_plate_order) && notInInputs {
				output_plate_order = append(output_plate_order, pa)
			}
		}

		if len(output_plate_order) > len(output_plates) {
			for _, o := range output_plate_order {
				fmt.Println(o)
			}
			fmt.Println("---")

			for _, p := range output_plates {
				fmt.Println(p.ID)
			}

			return nil, wtype.LHError(wtype.LH_ERR_DIRE, fmt.Sprintf("Plate number inconsistency: %d != %d (here: %d)", len(output_plate_order), len(output_plates), 121))
		}

		request.Output_plate_order = output_plate_order

	}

	// tips
	tips := request.Tips

	// just need to set the tip types
	// these should be distinct... we should check really
	// ...eventually
	if len(tips) != 0 {
		for _, tb := range tips {
			if tb == nil {
				continue
			}
			params.Tips = append(params.Tips, tb.Tips[0][0])
		}
	}

	setup := make(map[string]interface{})
	// make sure anything in setup is in synch

	for pos, id := range params.PosLookup {
		if id != "" {
			p := params.PlateLookup[id]
			setup[pos] = p
		}

	}

	// place outputs

	for _, pid := range output_plate_order {
		p := output_plates[pid]
		allowed, isConstrained := p.IsConstrainedOn(params.Model)
		if !isConstrained {
			allowed = make([]string, 0, 1)
		}
		position := get_first_available_preference(output_preferences, setup, allowed)

		if position == "" {
			//RaiseError("No positions left for output")
			err := wtype.LHError(wtype.LH_ERR_NO_DECK_SPACE, fmt.Sprint("No position left for output ", p.Name(), " Type: ", p.Type, " Constrained: ", isConstrained, " allowed positions: ", allowed))
			return request, err
		}

		setup[position] = p
		plate_lookup[p.ID] = position

		params.AddPlate(position, p)
		//logger.Info(fmt.Sprintf("Output plate of type %s in position %s", p.Type, position))
	}

	for _, pid := range input_plate_order {
		p := input_plates[pid]

		if p == nil {
			err := wtype.LHError(wtype.LH_ERR_DIRE, fmt.Sprint("Plate with id ", pid, " in input_plate_order does not exist in input_plates"))
			return request, err
		}

		allowed, isConstrained := p.IsConstrainedOn(params.Model)
		if !isConstrained {
			allowed = make([]string, 0, 1)
		}
		position := get_first_available_preference(input_preferences, setup, allowed)

		if position == "" {
			//RaiseError("No positions left for input")
			err := wtype.LHError(wtype.LH_ERR_NO_DECK_SPACE, fmt.Sprint("No position left for input ", p.Name(), " Type: ", p.Type, " Constrained: ", isConstrained, " allowed positions: ", allowed))
			return request, err
		}
		//fmt.Println("PLAATE: ", position)
		setup[position] = p
		plate_lookup[p.ID] = position
		params.AddPlate(position, p)
		//logger.Info(fmt.Sprintf("Input plate of type %s in position %s", p.Type, position))
	}

	// add the waste
	s := params.TipWastesMounted()

	if s == 0 {
		var waste *wtype.LHTipwaste
		// this should be added to the automagic config setup... however it will require adding to the
		// representation of the liquid handler
		if params.Model == "Pipetmax" {
			waste = factory.GetTipwasteByType("Gilsontipwaste")
		} else if params.Model == "GeneTheatre" || params.Model == "Felix" {
			waste = factory.GetTipwasteByType("CyBiotipwaste")
		} else if params.Model == "Human" {
			waste = factory.GetTipwasteByType("Manualtipwaste")
		}

		if waste != nil {
			params.AddTipWaste(waste)
		} else {
			err := wtype.LHError(wtype.LH_ERR_OTHER, fmt.Sprint("No tip waste defined for model ", params.Model))
			return nil, err
		}
	}
	//request.Setup = setup
	request.Plate_lookup = plate_lookup
	return request, nil
}

func get_first_available_preference(prefs []string, setup map[string]interface{}, allowed []string) string {
	for _, pref := range prefs {
		if len(allowed) != 0 && !isInStrArr(pref, allowed) {
			continue
		}
		_, ok := setup[pref]
		if !ok {
			return pref
		}

	}
	return ""
}

func isInStrArr(q string, ar []string) bool {
	for _, s := range ar {
		if q == s {
			return true
		}
	}

	return false
}
