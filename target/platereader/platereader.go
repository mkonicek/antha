package platereader

import (
	"context"
	"fmt"
	"github.com/antha-lang/antha/ast"
	"github.com/antha-lang/antha/target"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/driver"
	platereader "github.com/antha-lang/antha/driver/antha_platereader_v1"
	"strings"
	"github.com/antha-lang/antha/codegen"
	"errors"
)


// PlateReader defines the state of a plate-reader device
type PlateReader struct {
}


// Ensure satisfies Device interface
var _ target.Device = &PlateReader{}


// CanCompile implements a Device
func (a *PlateReader) CanCompile(req ast.Request) bool {
	can := ast.Request{}
	can.Selector = append(can.Selector, target.DriverSelectorV1WriteOnlyPlateReader)
	return can.Contains(req)
}

// MoveCost implements a Device
func (a *PlateReader) MoveCost(from target.Device) int {
	return 0
}


// Compile implements a Device
func (a *PlateReader) Compile(ctx context.Context, nodes []ast.Node) ([]target.Inst, error) {

	// Find the LHComponentID for the samples to measure. We'll then search
	// for these later.
	lhCmpIds := make(map[string]bool)
	for _, node := range nodes {
		cmd, ok := node.(*ast.Command)
		if !ok {
			panic(fmt.Sprintf("expected *ast.Command. Got: %T", node))
		}
		inst, ok := cmd.Inst.(*wtype.PRInstruction)
		if !ok {
			// TODO: Do we want to panic?
			panic(fmt.Sprintf("expected PRInstruction. Got: %T", cmd.Inst))
		}
		lhID := inst.ComponentIn.GetID()
		lhCmpIds[lhID] = true
	}

	// Look for the sample locations
	lhPlateLocations := make(map[string]string) // {cmpId: PlateId}
	lhWellLocations := make(map[string]string) // {cmpId: A1Coord}
	for _, cmd := range ast.FindReachingCommands(nodes) {
		insts := cmd.Output.(*codegen.Result).Insts
		for _, inst := range insts {
			mix, ok := inst.(*target.Mix)
			if !ok {
				// TODO: Deal with other instructions
				fmt.Printf("Expected *target.Mix, got: %T", inst)
				continue
			}
			for _, plate := range mix.FinalProperties.Plates {
				for _, well := range plate.Wellcoords {
					for lhCmpID := range lhCmpIds {
						if strings.Contains(well.WContents.ParentID, lhCmpID) {
							// Found a component that we are looking for
							lhPlateLocations[lhCmpID] = plate.ID
							lhWellLocations[lhCmpID] = well.Crds
						}
					}
				}
			}
		}
	}

	prInsts := make([]wtype.PRInstruction, 0)
	for _, node := range nodes {
		cmd := node.(*ast.Command)
		precInst := cmd.Inst.(*wtype.PRInstruction)
		prInsts = append(prInsts, *precInst)
	}

	// Merge PR instructions
	insts, err := a.mergePRInsts(prInsts, lhWellLocations, lhPlateLocations)
	if err != nil {
		return nil, err
	}
	return insts, nil
}

// PRInstructions with the same key can be executed on the same plate-read cycle
func prKey(inst *wtype.PRInstruction) (string, error) {
	return fmt.Sprintf("%s", inst.Options), nil
}

// Merge PRInstructions
func (a* PlateReader) mergePRInsts(insts []wtype.PRInstruction, wellLocs map[string]string, plateLocs map[string]string) ([]target.Inst, error) {

	// Simple case
	if len(insts) == 0 {
		return []target.Inst{}, nil
	}

	// Check for only 1 plate (for now)
	plateLocUnique := make(map[string]bool)
	for _, plateID := range plateLocs {
		plateLocUnique[plateID] = true
	}
	if len(plateLocUnique) > 1 {
		return []target.Inst{}, errors.New("current only supports single plate")
	}

	// Group instructions by PRInstruction
	groupBy := make(map[string]wtype.PRInstruction)  // {key: instruction}
	groupedWellLocs := make(map[string][]string)  // {key: []A1Coord}
	for _, inst := range insts {
		key, err := prKey(&inst)
		if err != nil {
			return nil, err
		}
		cmpID := inst.ComponentIn.GetID()
		groupBy[key] = inst
		groupedWellLocs[key] = append(groupedWellLocs[key], wellLocs[cmpID])
	}

	// Emit the driver calls
	var calls []driver.Call
	for key, inst := range groupBy {
		cmpID := inst.ComponentIn.GetID()

		wellString := strings.Join(groupedWellLocs[key], " ")
		plateID := plateLocs[cmpID]

		call := driver.Call{
			Method:"PRRunProtocolByName",
			Args: &platereader.ProtocolRunRequest{
				ProtocolName: "Custom",
				PlateID: plateID,
				PlateLayout: wellString,
				ProtocolOptions: inst.Options,
			},
			Reply: &platereader.BoolReply{},
		}
		calls = append(calls, call)
	}

	inst := &target.Run{
		Dev:   a,
		Label: "use plate reader",
		Calls: calls,
	}
	return []target.Inst{inst}, nil
}
