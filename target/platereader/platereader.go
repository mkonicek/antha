package platereader

import (
	"context"
	"fmt"
	"github.com/antha-lang/antha/ast"
	"github.com/antha-lang/antha/target"
	"github.com/antha-lang/antha/graph"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/driver"
	platereader "github.com/antha-lang/antha/driver/antha_platereader_v1"
)


// PlateReader defines the interface to a plate reader device
type PlateReader struct {
}


// CanCompile implements a Device
func (a *PlateReader) CanCompile(req ast.Request) bool {
	can := ast.Request{}
	can.Selector = append(can.Selector, target.DriverSelectorV1WriteOnlyPlateReader)
	return can.Contains(req)
}


func (a *PlateReader) MoveCost(from target.Device) int {
	return 0
}


func (a *PlateReader) Compile(ctx context.Context, nodes []ast.Node) ([]target.Inst, error) {

	// Find the LHComponentID for the samples to measure. We'll then search
	// for these later.
	lhCmpIds := make(map[string]bool)
	for _, node := range nodes {
		cmd, ok := node.(*ast.Command)
		if !ok {
			// TODO: Do we want to panic?
			panic(fmt.Sprintf("expected *ast.Command. Got: %T", node))
		}
		inst, ok := cmd.Inst.(*wtype.PRInstruction)
		if !ok {
			// TODO: Do we want to panic?
			panic(fmt.Sprintf("expected PRInstruction. Got: %T", cmd.Inst))
		}
		lhId := inst.ComponentIn.GetID()
		fmt.Println("LHID::", lhId)
		lhCmpIds[lhId] = true
	}

	// Breadth-first search to find location for all LhComponents
	g := ast.ToGraph(ast.ToGraphOpt{Roots: nodes, WhichDeps: ast.DataDeps})
	lhLocations := make(map[string]string)

	// Parse the parentId to get the LHComponentId
	getIdFromParent := func (parentId string) string {
		if len(parentId) > 36 {
			return parentId[:36]
		}
		return ""
	}

	// Apply to each node we visit
	apply := func(node graph.Node) {
		cmd, ok := node.(*ast.Command)
		if !ok {
			return
		}
		insts := cmd.Output.([]target.Inst)
		for _, inst := range insts {
			mix, ok := inst.(*target.Mix)
			if !ok {
				fmt.Printf("Expected *target.Mix, got: %T", inst)
				continue
			}
			for _, plate := range mix.FinalProperties.Plates {
				for _, well := range plate.Wellcoords {
					lhCmpId := getIdFromParent(well.WContents.ParentID)
					if len(lhCmpId) > 0 && lhCmpIds[lhCmpId] {
						lhLocations[lhCmpId] = fmt.Sprintf("%s:%s:%s", well.Crds, plate.ID, plate.Name())
					}
				}
			}
		}
	}

	// Traverse breadth-first
	stack := make([]graph.Node, 0)
	seen := make(map[graph.Node]bool)
	for _, node := range nodes {
		stack = append(stack, node)
	}
	for len(stack) > 0 {
		if len(lhLocations) == len(lhCmpIds) {
			// Found the locations of all samples, so stop.
			break
		}
		node := stack[0]
		stack = stack[1:]
		if seen[node] {
			continue
		}
		apply(node)  // visit the node
		seen[node] = true
		for i := 0; i < g.NumOuts(node); i ++ {
			stack = append(stack, g.Out(node, i))
		}
	}


	for k, v := range lhLocations {
		fmt.Println("WELL_LOCATION:", k, v)
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


	// In language of S2
	return []target.Inst{inst}, nil
}
