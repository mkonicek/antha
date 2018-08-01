package execute

import (
	"context"

	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/ast"
	"github.com/antha-lang/antha/target"
)

// QPCROptions are the options for a QPCR request.
type QPCROptions struct {
	Reactions  []*wtype.Liquid
	Definition string
	Barcode    string
	TagAs      string
}

func runQPCR(ctx context.Context, opts QPCROptions, command string) (*commandInst, []*wtype.Liquid) {
	inst := ast.NewQPCRInstruction()
	inst.Command = command
	inst.ComponentIn = opts.Reactions
	inst.Definition = opts.Definition
	inst.Barcode = opts.Barcode
	inst.TagAs = opts.TagAs
	inst.ComponentOut = []*wtype.Liquid{}

	for _, r := range inst.ComponentIn {
		inst.ComponentOut = append(inst.ComponentOut, newCompFromComp(ctx, r))
	}

	return &commandInst{
		Args:   opts.Reactions,
		result: inst.ComponentOut,
		Command: &ast.Command{
			Inst: inst,
			Requests: []ast.Request{
				{
					Selector: []ast.NameValue{
						target.DriverSelectorV1QPCRDevice,
					},
				},
			},
		},
	}, inst.ComponentOut
}

// RunQPCRExperiment starts a new QPCR experiment, using an experiment input file.
func RunQPCRExperiment(ctx context.Context, opt QPCROptions) ([]*wtype.Liquid, *QPCRDataIngester) {
	inst, outputs := runQPCR(ctx, opt, "RunExperiment")
	Issue(ctx, inst)
	return inst.result, &QPCRDataIngester{
		ctx:      ctx,
		dependOn: outputs,
	}
}

// RunQPCRFromTemplate starts a new QPCR experiment, using a template input file.
func RunQPCRFromTemplate(ctx context.Context, opt QPCROptions) ([]*wtype.Liquid, *QPCRDataIngester) {
	inst, outputs := runQPCR(ctx, opt, "RunExperimentFromTemplate")
	Issue(ctx, inst)
	return inst.result, &QPCRDataIngester{
		ctx:      ctx,
		dependOn: outputs,
	}
}

type QPCRDataIngester struct {
	ctx      context.Context // this will go away once we move away from pushing trace into Context
	dependOn []*wtype.Liquid
}

// Wait for the QPCR to produce data, and attempt to parse and ingest
// that data. There are no options for the ingestion.
func (ingester *QPCRDataIngester) ExpectData() {
	cmd := target.GetExpectDataCommand(ingester.ctx, nil, target.DriverSelectorV1QPCRDevice)
	Issue(ingester.ctx, &commandInst{
		Args:    ingester.dependOn,
		Command: cmd,
	})
}
