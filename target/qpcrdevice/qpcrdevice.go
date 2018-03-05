package qpcrdevice

import (
	"context"
	"errors"
	"fmt"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/ast"
	"github.com/antha-lang/antha/codegen"
	"github.com/antha-lang/antha/driver"
	quantstudio "github.com/antha-lang/antha/driver/antha_quantstudio_v1"
	"github.com/antha-lang/antha/target"
	"strings"
)

// PlateReader defines the state of a plate-reader device
type QPCRDevice struct {
}

// Ensure satisfies Device interface
var _ target.Device = (*QPCRDevice)(nil)

// NewQPCRDevice returns a new QPCR Machine
// Used by antha-runner
func NewQPCRDevice() *QPCRDevice {
	ret := &QPCRDevice{}
	return ret
}

// CanCompile implements a Device
func (a *QPCRDevice) CanCompile(req ast.Request) bool {
	can := ast.Request{}
	can.Selector = append(can.Selector, target.DriverSelectorV1QPCRDevice)
	return can.Contains(req)
}

// MoveCost implements a Device
func (a *QPCRDevice) MoveCost(from target.Device) int {
	return 0
}

// Compile implements a Device
func (a *QPCRDevice) Compile(ctx context.Context, nodes []ast.Node) ([]target.Inst, error) {


	var qpcrInsts []*wtype.QPCRInstruction
	for _, node := range nodes {
		cmd := node.(*ast.Command)
		qpcrInsts = append(qpcrInsts, cmd.Inst.(*wtype.QPCRInstruction))
	}

	// Merge PR instructions
	insts, err := a.mergePRInsts(prInsts, lhWellLocations, lhPlateLocations)
	if err != nil {
		return nil, err
	}
	return insts, nil
}

// Merge PRInstructions
func (a *PlateReader) mergePRInsts(prInsts []*wtype.PRInstruction, wellLocs map[string]string, plateLocs map[string]string) ([]target.Inst, error) {

	// Simple case
	if len(prInsts) == 0 {
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

	var insts []target.Inst
	insts = append(insts, &target.Prompt{
		Message: "Please put plate(s) into plate reader and click ok to start plate reader",
	})
	insts = append(insts, &target.Run{
		Dev:   a,
		Label: "use plate reader",
		Calls: calls,
	})
	return target.SequentialOrder(insts...), nil
}
