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
func NewQPCRDevice() *QPCRDevice {
	ret := &QPCRDevice{}
	return ret
}

// CanCompile implements a Device
func (a *QPCRDevice) CanCompile(req ast.Request) bool {
	can := ast.Request{}
	can.Selector = append(can.Selector, target.DriverSelectorV1QPCRDevice, target.DriverSelectorV1DataSource)
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
		return driver.Call{}, err
	}

	jobID := inst.TagAs
	if len(jobID) == 0 {
		jobID = os.Getenv("METADATA_JOB_ID")
	}

	var instID string
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

	return call, nil
}

// qpcrCall is a utility structure, combining a call with some information about it.
type qpcrCall struct {
	Calls          []driver.Call
	ExperimentFile string
	Barcode        string
}

// callsForInstructions generates a number of device calls based on supplied qPCR instructions.
// At present we assert that there is a single experiment request, with a single barcode.
func callsForInstructions(qpcrInsts []*ast.QPCRInstruction) ([]qpcrCall, error) {
	var calls []qpcrCall
	var experimentFile string
	var barcode string

	for _, inst := range qpcrInsts {
		// Note that we only accept the first value at present. (i.e. a single qPCR run within a workflow.)

		if inst.Definition == "" {
			return nil, errors.New("blank experiment file for qPCR")
		}

		if len(experimentFile) > 0 {
			if experimentFile != inst.Definition {
				return nil, fmt.Errorf("unexpected multiple experiment files %s, %s", experimentFile, inst.Definition)
			}

			if barcode != inst.Barcode {
				return nil, fmt.Errorf("unexpected multiple barcodes %s, %s", barcode, inst.Barcode)
			}
		} else {
			experimentFile = inst.Definition
			barcode = inst.Barcode
		}

		call, err := singleCallForInstruction(inst)
		if err != nil {
			return nil, err
		}

		// Only creating a single instruction at present.
		if len(calls) == 0 {
			calls = append(calls, qpcrCall{[]driver.Call{call}, experimentFile, barcode})
		}
	}

	return calls, nil
}

// interpolatePrompts inserts a prompt before each qPCR run.
func interpolatePrompts(calls []qpcrCall, device target.Device) []target.Inst {

	// Interpolate a prompt before each qPCR call.
	var insts []target.Inst
	for _, call := range calls {

		insts = append(insts, &target.Prompt{
			Message: "Ensure that the experiment file " + call.ExperimentFile + " is configured, then put the plate (" + call.Barcode + ") into qPCR device. Check that the driver software is running. Once ready, accept to start the qPCR analysis.",
		})
		insts = append(insts, &target.Run{
			Dev:   device,
			Label: "Perform qPCR Analysis",
			Calls: call.Calls,
		})
	}

	return insts
}

// Compile implements a qPCR device.
func (a *QPCRDevice) Compile(ctx context.Context, nodes []ast.Node) ([]target.Inst, error) {

	var qpcrInsts []*ast.QPCRInstruction
	for _, node := range nodes {
		cmd := node.(*ast.Command)
		qpcrInsts = append(qpcrInsts, cmd.Inst.(*ast.QPCRInstruction))
	}

	calls, err := callsForInstructions(qpcrInsts)

	if err != nil {
		return nil, err
	}

	// Interpolate a prompt before each qPCR call.
	promptedInstructions := interpolatePrompts(calls, a)
	return target.SequentialOrder(promptedInstructions...), nil
}

func (q *QPCRDevice) ExpectDataTemplate() *ast.ExpectInst {
	return &ast.ExpectInst{
		ParserID: "myFirstQPCRParser",
	}
}
