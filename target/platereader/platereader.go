package platereader

import (
	"context"
	"fmt"
	"github.com/antha-lang/antha/ast"
	"github.com/antha-lang/antha/target"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/driver"
	platereader "github.com/antha-lang/antha/driver/antha_platereader_v1"
)


// PlateReader defines the interface to a plate reader device
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
		fmt.Println("LHID::", lhID)
		lhCmpIds[lhID] = true
	}

	// Parse the parentId to get the LHComponentId
	getIDFromParent := func (parentId string) string {
		if len(parentId) > 36 {
			return parentId[:36]
		}
		return ""
	}

	// Look for the sample locations
	found := make(map[string]bool)
	lhLocations := make(map[string][]string) // {plateID : []A1Coord}
	for _, cmd := range ast.FindReachingCommands(nodes) {
		insts := cmd.Output.([]target.Inst)
		for _, inst := range insts {
			mix, ok := inst.(*target.Mix)
			if !ok {
				fmt.Printf("Expected *target.Mix, got: %T", inst)
				continue
			}
			for _, plate := range mix.FinalProperties.Plates {
				for _, well := range plate.Wellcoords {
					lhCmpID := getIDFromParent(well.WContents.ParentID)
					if len(lhCmpID) > 0 && lhCmpIds[lhCmpID] && !found[lhCmpID] {
						// Found a component that we are looking for
						lhLocations[plate.ID] = append(lhLocations[plate.ID], well.Crds)
						found[lhCmpID] = true
					}
				}
			}
		}
	}

	// Groupby plateID
	for plateId, coords := range lhLocations {
		fmt.Println("WELL_LOCATION:", plateId, coords)
	}


	// Need to merge lhcomponents that are on the same plate and same
	// wavelength.
	calls := []driver.Call{
		{Method:"PRRunProtocolByName", Args: &platereader.ProtocolRunRequest{string(600), "PLATE_ID"}, Reply: &platereader.BoolReply{}},
		{Method:"PROpen", Args: &platereader.ProtocolRunRequest{string(600), "PLATE_ID"}, Reply: &platereader.BoolReply{}},
	}

	inst := &target.Run{
		Dev:   a,
		Label: "Plate-Reader",
		Calls: calls,
	}

	// For debug
	for _, call := range calls {
		fmt.Println("driver.Call", call)
	}

	// In language of S2
	return []target.Inst{inst}, nil
}
