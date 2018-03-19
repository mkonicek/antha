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
	"golang.org/x/net/html/atom"
	"os"
)

// QPCRDevice defines the state of a qpcr device device
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

// singleCallForInstruction generates a device call, based on a single qPCR insruction.
func singleCallForInstruction(inst *ast.QPCRInstruction) (driver.Call, error) {
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
		return driver.Call{}, err
	}

	jobId := inst.TagAs
	if jobId == "" {
		jobId = os.Getenv("METADATA_JOB_ID")
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
				JobId: jobId,
			},
		},
		Reply: &framework.CommandResponse{},
	}

	return call, nil
}

// Compile implements a Device
func (a *QPCRDevice) Compile(ctx context.Context, nodes []ast.Node) ([]target.Inst, error) {

	var qpcrInsts []*ast.QPCRInstruction
	for _, node := range nodes {
		cmd := node.(*ast.Command)
		qpcrInsts = append(qpcrInsts, cmd.Inst.(*ast.QPCRInstruction))
	}

	var calls []driver.Call
	experimentFile := ""

	for _, inst := range qpcrInsts {
		// Note that we only accept the first value at present. (i.e. a single qPCR run within a workflow.)

		if inst.Definition == "" {
			return nil, errors.New("Blank experiment file for qPCR.")
		}

		if len(experimentFile) > 0 {
			if experimentFile != inst.Definition {
				return nil, fmt.Errorf("Unexpected multiple experiment files %s, %s", experimentFile, inst.Definition)
			}
		} else {
			experimentFile = inst.Definition
		}

		call, err := singleCallForInstruction(inst)
		if err != nil {
			return nil, err
		}
		calls = append(calls, call)
	}

	var insts []target.Inst
	insts = append(insts, &target.Prompt{
		Message: "Please ensures that the experiment file " + experimentFile + " is configured, then put the plate into qPCR device. Check that the driver software is running. Once ready, accept to start the qPCR analysis.",
	})
	insts = append(insts, &target.Run{
		Dev:   a,
		Label: "Perform qPCR Analysis",
		Calls: calls,
	})
	return target.SequentialOrder(insts...), nil
}
