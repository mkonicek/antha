package qpcrdevice

import (
	"context"
	"errors"
	"fmt"
	"github.com/antha-lang/antha/ast"
	"github.com/antha-lang/antha/driver"
	framework "github.com/antha-lang/antha/driver/antha_framework_v1"
	quantstudio "github.com/antha-lang/antha/driver/antha_quantstudio_v1"
	"github.com/antha-lang/antha/target"
	"github.com/golang/protobuf/proto"
	"os"
)

// QPCRDevice defines the state of a qpcr device device
type QPCRDevice struct {
}

// Ensure satisfies Device interface
var _ target.Device = (*QPCRDevice)(nil)

// NewQPCRDevice returns a new QPCR Machine
func New() *QPCRDevice {
	return &QPCRDevice{}
}

// CanCompile implements a Device
func (a *QPCRDevice) CanCompile(req ast.Request) bool {
	can := ast.Request{
		Selector: []ast.NameValue{target.DriverSelectorV1QPCRDevice},
	}
	return can.Contains(req)
}

// MoveCost implements a Device
func (a *QPCRDevice) MoveCost(from target.Device) int64 {
	return 0
}

func (dev *QPCRDevice) callForInstruction(inst *ast.QPCRInstruction) (target.Inst, error) {
	if inst.Definition == "" {
		return nil, errors.New("Blank experiment file for qPCR instruction.")
	}

	message := &quantstudio.TemplatedRequest{
		SessionInstrument: &quantstudio.SessionInstrument{
			Session:    &quantstudio.Session{},
			Instrument: &quantstudio.Instrument{},
		},
		TemplateFile: &quantstudio.ExperimentFile{
			Url: inst.Definition,
		},
		Barcode: &quantstudio.Barcode{
			Barcode: inst.Barcode,
		},
	}

	messageBytes, err := proto.Marshal(message)
	if err != nil {
		return nil, err
	}

	jobID := inst.TagAs
	if len(jobID) == 0 {
		jobID = os.Getenv("METADATA_JOB_ID")
	}

	instID := ""
	if len(inst.ComponentIn) > 0 {
		instID = inst.ComponentIn[0].ID
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
				JobId: jobID,
			},
		},
		Reply: &framework.CommandResponse{},
	}

	return &target.Run{
		Dev:   dev,
		Label: "Perform qPCR Analysis",
		Calls: []driver.Call{call},
	}, nil
}

func (dev *QPCRDevice) makePrompt(inst *ast.QPCRInstruction) target.Inst {
	bc := inst.Barcode
	if bc != "" {
		bc = " (" + bc + ")" // deliberate leading space
	}
	return &target.Prompt{
		Message: fmt.Sprintf("Ensure that the experiment file %s is configured, then put the plate%s into qPCR device. Check that the driver software is running. Once ready, accept to start the qPCR analysis.",
			inst.Definition, bc),
	}
}

// Compile implements a qPCR device.
func (dev *QPCRDevice) Compile(ctx context.Context, nodes []ast.Node) ([]target.Inst, error) {
	if len(nodes) > 1 {
		return nil, fmt.Errorf("Currently only permit a single qPCR instruction per workflow. Received %d", len(nodes))
	}

	insts := make(target.Insts, 0, 2*len(nodes))
	for _, node := range nodes {
		inst := node.(*ast.Command).Inst.(*ast.QPCRInstruction)
		if call, err := dev.callForInstruction(inst); err != nil {
			return nil, err
		} else {
			insts = append(insts, dev.makePrompt(inst), call)
		}
	}

	insts.SequentialOrder()
	return insts, nil
}
