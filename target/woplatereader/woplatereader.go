package woplatereader

import (
	"errors"
	"fmt"
	"strings"

	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/driver"
	platereader "github.com/antha-lang/antha/driver/antha_platereader_v1"
	"github.com/antha-lang/antha/laboratory/effects"
	"github.com/antha-lang/antha/target"
	"github.com/antha-lang/antha/workflow"
)

// WOPlateReader defines the state of a write only plate-reader device
type WOPlateReader struct {
	id workflow.DeviceInstanceID
}

func NewWOPlateReaderInstances(tgt *target.Target, config workflow.PlateReaderConfig) error {
	for id := range config.Devices {
		if err := tgt.AddDevice(New(id)); err != nil {
			return err
		}
	}
	return nil
}

// Ensure satisfies Device interface
var (
	_ effects.Device = &WOPlateReader{}
)

// returns a new Write-Only Plate Reader Used by antha-runner
func New(id workflow.DeviceInstanceID) *WOPlateReader {
	return &WOPlateReader{
		id: id,
	}
}

func (a *WOPlateReader) Id() workflow.DeviceInstanceID {
	return a.id
}

// CanCompile implements a Device
func (a *WOPlateReader) CanCompile(req effects.Request) bool {
	can := effects.Request{
		Selector: []effects.NameValue{
			target.DriverSelectorV1WriteOnlyPlateReader,
		},
	}
	return can.Contains(req)
}

// Compile implements a Device
func (a *WOPlateReader) Compile(labEffects *effects.LaboratoryEffects, dir string, nodes []effects.Node) ([]effects.Inst, error) {
	// Find the LHComponentID for the samples to measure. We'll then search
	// for these later.
	prInsts := make([]*wtype.PRInstruction, 0, len(nodes))
	lhCmpIDs := make(map[string]bool)
	for _, node := range nodes {
		if cmd, ok := node.(*effects.Command); !ok {
			return nil, fmt.Errorf("cannot compile %T", node)
		} else if inst, ok := cmd.Inst.(*wtype.PRInstruction); !ok {
			return nil, fmt.Errorf("cannot compile %T", cmd.Inst)
		} else {
			lhCmpIDs[inst.ComponentIn.GetID()] = true
			prInsts = append(prInsts, inst)
		}
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
	for _, cmd := range effects.FindReachingCommands(nodes) {
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
func (a *WOPlateReader) mergePRInsts(prInsts []*wtype.PRInstruction, wellLocs map[string]string, plateLocs map[string]string) ([]effects.Inst, error) {

	// Simple case
	if len(prInsts) == 0 {
		return []effects.Inst{}, nil
	}

	// Check for only 1 plate (for now)
	plateLocUnique := make(map[string]bool)
	for _, plateID := range plateLocs {
		plateLocUnique[plateID] = true
	}
	if len(plateLocUnique) > 1 {
		return []effects.Inst{}, errors.New("current only supports single plate")
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

	insts := effects.Insts{
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

func (a *WOPlateReader) Connect(*workflow.Workflow) error {
	return nil
}

func (a *WOPlateReader) Close() {}
