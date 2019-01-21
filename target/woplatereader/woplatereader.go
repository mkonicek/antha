package woplatereader

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/ast"
	"github.com/antha-lang/antha/driver"
	platereader "github.com/antha-lang/antha/driver/antha_platereader_v1"
	"github.com/antha-lang/antha/target"
)

// WOPlateReader defines the state of a write only plate-reader device
type WOPlateReader struct{}

// Ensure satisfies Device interface
var _ ast.Device = (*WOPlateReader)(nil)

// returns a new Write-Only Plate Reader Used by antha-runner
func New() *WOPlateReader {
	return &WOPlateReader{}
}

// CanCompile implements a Device
func (a *WOPlateReader) CanCompile(req ast.Request) bool {
	can := ast.Request{}
	can.Selector = append(can.Selector, target.DriverSelectorV1WriteOnlyPlateReader)
	return can.Contains(req)
}

// Compile implements a Device
func (a *WOPlateReader) Compile(ctx context.Context, nodes []ast.Node) ([]ast.Inst, error) {

	// Find the LHComponentID for the samples to measure. We'll then search
	// for these later.
	lhCmpIDs := make(map[string]bool)
	for _, node := range nodes {
		cmd := node.(*ast.Command)
		inst, ok := cmd.Inst.(*wtype.PRInstruction)
		if !ok {
			return nil, fmt.Errorf("expected PRInstruction. Got: %T", cmd.Inst)
		}
		lhID := inst.ComponentIn.GetID()
		lhCmpIDs[lhID] = true
	}

	lhPlateLocations := make(map[string]string) // {cmpId: PlateId}
	lhWellLocations := make(map[string]string)  // {cmpId: A1Coord}
	findComps := func(mix *target.Mix) {
		for _, plate := range mix.FinalProperties.Plates {
			for _, well := range plate.Wellcoords {
				for lhCmpID := range lhCmpIDs {
					if strings.Contains(well.WContents.ParentID, lhCmpID) {
						// Found a component that we are looking for
						lhPlateLocations[lhCmpID] = plate.ID
						lhWellLocations[lhCmpID] = well.Crds.FormatA1()
					}
				}
			}
		}
	}

	// Look for the sample locations
	for _, cmd := range ast.FindReachingCommands(nodes) {
		insts := cmd.Output
		for _, inst := range insts {
			mix, ok := inst.(*target.Mix)
			if !ok {
				// TODO: Deal with other instructions
				continue
			}
			findComps(mix)
		}
	}

	var prInsts []*wtype.PRInstruction
	for _, node := range nodes {
		cmd := node.(*ast.Command)
		prInsts = append(prInsts, cmd.Inst.(*wtype.PRInstruction))
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
	return inst.Options, nil
}

// Merge PRInstructions
func (a *WOPlateReader) mergePRInsts(prInsts []*wtype.PRInstruction, wellLocs map[string]string, plateLocs map[string]string) ([]ast.Inst, error) {

	// Simple case
	if len(prInsts) == 0 {
		return []ast.Inst{}, nil
	}

	// Check for only 1 plate (for now)
	plateLocUnique := make(map[string]bool)
	for _, plateID := range plateLocs {
		plateLocUnique[plateID] = true
	}
	if len(plateLocUnique) > 1 {
		return []ast.Inst{}, errors.New("current only supports single plate")
	}

	// Group instructions by PRInstruction
	groupBy := make(map[string]*wtype.PRInstruction) // {key: instruction}
	groupedWellLocs := make(map[string][]string)     // {key: []A1Coord}
	for _, inst := range prInsts {
		key, err := prKey(inst)
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
			Method: "PRRunProtocolByName",
			Args: &platereader.ProtocolRunRequest{
				ProtocolName:    "Custom",
				PlateID:         plateID,
				PlateLayout:     wellString,
				ProtocolOptions: inst.Options,
			},
			Reply: &platereader.BoolReply{},
		}
		calls = append(calls, call)
	}

	insts := ast.Insts{
		&target.Prompt{
			Message: "Please put plate(s) into plate reader and click ok to start plate reader",
		},
		&target.Run{
			Dev:   a,
			Label: "use plate reader",
			Calls: calls,
		},
	}
	insts.SequentialOrder()
	return insts, nil
}
