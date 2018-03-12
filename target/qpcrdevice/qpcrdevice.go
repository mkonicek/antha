package qpcrdevice

import (
	"context"
	"github.com/antha-lang/antha/ast"
	"github.com/antha-lang/antha/driver"
	framework "github.com/antha-lang/antha/driver/antha_framework_v1"
	quantstudio "github.com/antha-lang/antha/driver/antha_quantstudio_v1"
	"github.com/antha-lang/antha/target"
	"github.com/golang/protobuf/proto"
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

	var qpcrInsts []*ast.QPCRInstruction
	for _, node := range nodes {
		cmd := node.(*ast.Command)
		qpcrInsts = append(qpcrInsts, cmd.Inst.(*ast.QPCRInstruction))
	}

	var calls []driver.Call
	for _, inst := range qpcrInsts {
		instID := inst.ComponentIn.GetID()

		message := &quantstudio.TemplatedRequest{
			SessionInstrument: &quantstudio.SessionInstrument{
				Session: &quantstudio.Session{
					Id: "",
				},
				Instrument: &quantstudio.Instrument{
					Id: "",
				},
			},
			TemplateFile: &quantstudio.ExperimentFile{
				Url: inst.Definition,
			},
			Barcode: &quantstudio.Barcode{
				Barcode: inst.Barcode,
			},
			OutputPath: "",
		}

		messageBytes, err := proto.Marshal(message)
		if err != nil {
			return nil, err
		}

		call := driver.Call{
			Method: "RunFramework",
			Args: &framework.CommandRequest{
				Description: &framework.CommandDescription{
					CommandName: inst.Command,
					DeviceId:    "QPCRDevice",
					CommandId:   instID,
				},
				Data: &framework.CommandData{
					Data: messageBytes,
				},
				Metadata: &framework.CommandMetadata{
					JobId: "",
				},
			},
			Reply: &framework.CommandResponse{},
		}
		calls = append(calls, call)
	}

	var insts []target.Inst
	insts = append(insts, &target.Prompt{
		Message: "Please put plate into QPCR device and click ok to start experiment",
	})
	insts = append(insts, &target.Run{
		Dev:   a,
		Label: "use plate reader",
		Calls: calls,
	})
	return target.SequentialOrder(insts...), nil

}
