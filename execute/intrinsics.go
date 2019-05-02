package execute

import (
	"context"
	"github.com/antha-lang/antha/antha/anthalib/mixer"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/ast"
	"github.com/antha-lang/antha/inventory"
	"github.com/antha-lang/antha/microArch/sampletracker"
	"github.com/antha-lang/antha/target"
	"time"
)

// a commandInst is a generic intrinsic instruction
type commandInst struct {
	// Arguments to this command. Used to determine command dependencies.
	Args []*wtype.Liquid
	// Components created by this command. Returned back to user code
	result  []*wtype.Liquid
	Command *ast.Command
}

// SetInputPlate Indicate to the scheduler the the contents of the plate is user
// supplied. This modifies the argument to mark each well as such.
func SetInputPlate(ctx context.Context, plate *wtype.Plate) {
	sampletracker.FromContext(ctx).SetInputPlate(plate)
}

// An IncubateOpt are options to an incubate command
type IncubateOpt struct {
	// Time for which to incubate component
	Time wunit.Time
	// Temperature at which to incubate component
	Temp wunit.Temperature
	// Rate at which to shake incubator (force is device dependent)
	ShakeRate wunit.Rate
	// Radius at which ShakeRate is defined
	ShakeRadius wunit.Length

	// Time for which to pre-heat incubator
	PreTemp wunit.Temperature
	// Temperature at which to pre-heat incubator
	PreTime wunit.Time
	// Rate at which to pre-heat incubator
	PreShakeRate wunit.Rate
	// Radius at which PreShakeRate is defined
	PreShakeRadius wunit.Length
}

func updateLiquidID(ctx context.Context, in *wtype.Liquid) *wtype.Liquid {
	comp := in.Dup()
	comp.ID = wtype.GetUUID()
	comp.BlockID = wtype.NewBlockID(getID(ctx))
	comp.SetGeneration(comp.Generation() + 1)

	getMaker(ctx).UpdateAfterInst(in.ID, comp.ID)
	sampletracker.FromContext(ctx).UpdateIDOf(in.ID, comp.ID)

	return comp
}

// Incubate incubates a component
func Incubate(ctx context.Context, in *wtype.Liquid, opt IncubateOpt) *wtype.Liquid {
	// nolint: gosimple
	innerInst := &ast.IncubateInst{
		Time:           opt.Time,
		Temp:           opt.Temp,
		ShakeRate:      opt.ShakeRate,
		ShakeRadius:    opt.ShakeRadius,
		PreTemp:        opt.PreTemp,
		PreTime:        opt.PreTime,
		PreShakeRate:   opt.PreShakeRate,
		PreShakeRadius: opt.PreShakeRadius,
	}

	inst := &commandInst{
		Args:   []*wtype.Liquid{in},
		result: []*wtype.Liquid{updateLiquidID(ctx, in)},
		Command: &ast.Command{
			Inst: innerInst,
			Request: ast.Request{
				Selector: []ast.NameValue{
					target.DriverSelectorV1ShakerIncubator,
				},
			},
		},
	}

	Issue(ctx, inst)
	return inst.result[0]
}

// prompt... works pretty much like Handle does
// but passes the instruction to the planner
// in future this should generate handles as side-effects

type mixerPromptOpts struct {
	Components   []*wtype.Liquid
	ComponentsIn []*wtype.Liquid
	Message      string
	WaitTime     wunit.Time
}

func updateLiquidIds(ctx context.Context, in []*wtype.Liquid) []*wtype.Liquid {
	r := []*wtype.Liquid{}

	for _, c := range in {
		r = append(r, updateLiquidID(ctx, c))
	}

	return r

}

// MixerPrompt prompts user with a message during mixer execution
func MixerPrompt(ctx context.Context, message string, in ...*wtype.Liquid) []*wtype.Liquid {
	inst := mixerPrompt(ctx,
		mixerPromptOpts{
			Components:   updateLiquidIds(ctx, in),
			ComponentsIn: in,
			Message:      message,
		},
	)
	Issue(ctx, inst)
	return inst.result
}

// MixerWait prompts user with a message during mixer execution and waits for the specifed time before resuming.
func MixerWait(ctx context.Context, time wunit.Time, message string, in ...*wtype.Liquid) []*wtype.Liquid {
	inst := mixerPrompt(ctx,
		mixerPromptOpts{
			Components:   updateLiquidIds(ctx, in),
			ComponentsIn: in,
			Message:      message,
			WaitTime:     time,
		},
	)

	Issue(ctx, inst)
	return inst.result
}

// ExecuteMixes will ensure that all mix activities
// in a workflow prior to this point must be completed before Mix instructions after this point.
func ExecuteMixes(ctx context.Context, liquids ...*wtype.LHComponent) []*wtype.Liquid {
	return MixerPrompt(ctx, wtype.MAGICBARRIERPROMPTSTRING, liquids...)
}

// Prompt prompts user with a message
func Prompt(ctx context.Context, in *wtype.Liquid, message string) *wtype.Liquid {
	inst := &commandInst{
		Args:   []*wtype.Liquid{in},
		result: []*wtype.Liquid{updateLiquidID(ctx, in)},
		Command: &ast.Command{
			Inst: &ast.PromptInst{
				Message: message,
			},
			Request: ast.Request{
				Selector: []ast.NameValue{
					target.DriverSelectorV1Human,
				},
			},
		},
	}

	Issue(ctx, inst)
	return inst.result[0]
}

func mixerPrompt(ctx context.Context, opts mixerPromptOpts) *commandInst {
	inst := wtype.NewLHPromptInstruction()
	inst.SetGeneration(opts.ComponentsIn[0].Generation())
	inst.Message = opts.Message
	// precision will be cut to the nearest second
	inst.WaitTime = opts.WaitTime.AsDuration().Round(time.Second)
	for i := 0; i < len(opts.Components); i++ {
		inst.AddOutput(opts.Components[i])
		inst.AddInput(opts.ComponentsIn[i])
	}

	return &commandInst{
		Args:   opts.ComponentsIn,
		result: opts.Components,
		Command: &ast.Command{
			Inst: inst,
			Request: ast.Request{
				Selector: []ast.NameValue{
					target.DriverSelectorV1Prompter,
				},
			},
		},
	}
}

// PlateReadOpts defines plate-reader absorbance options
type PlateReadOpts struct {
	Sample  *wtype.Liquid
	Options string
}

func readPlate(ctx context.Context, opts PlateReadOpts) *commandInst {
	inst := wtype.NewPRInstruction()
	inst.ComponentIn = opts.Sample

	// Clone the component to represent the result of the AbsorbanceRead
	inst.ComponentOut = updateLiquidID(ctx, opts.Sample)
	inst.Options = opts.Options

	return &commandInst{
		Args:   []*wtype.Liquid{opts.Sample},
		result: []*wtype.Liquid{inst.ComponentOut},
		Command: &ast.Command{
			Inst: inst,
			Request: ast.Request{
				Selector: []ast.NameValue{
					target.DriverSelectorV1WriteOnlyPlateReader,
				},
			},
		},
	}
}

// PlateRead reads absorbance of a component
func PlateRead(ctx context.Context, opt PlateReadOpts) *wtype.Liquid {
	inst := readPlate(ctx, opt)
	Issue(ctx, inst)
	return inst.result[0]
}

// QPCROptions are the options for a QPCR request.
type QPCROptions struct {
	Reactions  []*wtype.Liquid
	Definition string
	Barcode    string
	TagAs      string
}

func runQPCR(ctx context.Context, opts QPCROptions, command string) *commandInst {
	inst := ast.NewQPCRInstruction()
	inst.Command = command
	inst.ComponentIn = opts.Reactions
	inst.Definition = opts.Definition
	inst.Barcode = opts.Barcode
	inst.TagAs = opts.TagAs
	inst.ComponentOut = []*wtype.Liquid{}

	for _, r := range inst.ComponentIn {
		inst.ComponentOut = append(inst.ComponentOut, updateLiquidID(ctx, r))
	}

	return &commandInst{
		Args:   opts.Reactions,
		result: inst.ComponentOut,
		Command: &ast.Command{
			Inst: inst,
			Request: ast.Request{
				Selector: []ast.NameValue{
					target.DriverSelectorV1QPCRDevice,
				},
			},
		},
	}
}

// RunQPCRExperiment starts a new QPCR experiment, using an experiment input file.
func RunQPCRExperiment(ctx context.Context, opt QPCROptions) []*wtype.Liquid {
	inst := runQPCR(ctx, opt, "RunExperiment")
	Issue(ctx, inst)
	return inst.result
}

// RunQPCRFromTemplate starts a new QPCR experiment, using a template input file.
func RunQPCRFromTemplate(ctx context.Context, opt QPCROptions) []*wtype.Liquid {
	inst := runQPCR(ctx, opt, "RunExperimentFromTemplate")
	Issue(ctx, inst)
	return inst.result
}

// NewComponent returns a new component given a component type
func NewComponent(ctx context.Context, typ string) *wtype.Liquid {
	c, err := inventory.NewComponent(ctx, typ)
	if err != nil {
		Errorf(ctx, "cannot make component %s: %s", typ, err)
	}
	return c
}

// NewPlate returns a new plate given a plate type
func NewPlate(ctx context.Context, typ string) *wtype.Plate {
	p, err := inventory.NewPlate(ctx, typ)
	if err != nil {
		Errorf(ctx, "cannot make plate %s: %s", typ, err)
	}
	return p
}

func mix(ctx context.Context, inst *wtype.LHInstruction) *commandInst {
	inst.BlockID = wtype.NewBlockID(getID(ctx))
	inst.Outputs[0].BlockID = inst.BlockID
	result := inst.Outputs[0]
	//result.BlockID = inst.BlockID // DELETEME

	mx := 0
	// from the protocol POV components need to be passed by value
	for i, c := range wtype.CopyComponentArray(inst.Inputs) {
		if c.CName == "" {
			panic("Nameless Component used in Mix - this is not permitted")
		}
		c.Order = i

		//result.MixPreserveTvol(c)
		if c.Generation() > mx {
			mx = c.Generation()
		}
		getMaker(ctx).UpdateAfterInst(c.ID, result.ID)
	}

	inst.SetGeneration(mx)
	result.SetGeneration(mx + 1)
	result.DeclareInstance()

	return &commandInst{
		Args: inst.Inputs,
		Command: &ast.Command{
			Inst: inst,
			Request: ast.Request{
				Selector: []ast.NameValue{
					target.DriverSelectorV1Mixer,
				},
			},
		},
		result: []*wtype.Liquid{result},
	}
}

func genericMix(ctx context.Context, generic *wtype.LHInstruction) *wtype.Liquid {
	inst := mix(ctx, generic)
	Issue(ctx, inst)
	if generic.Welladdress != "" {
		err := inst.result[0].SetWellLocation(generic.Welladdress)
		if err != nil {
			panic(err)
		}
	}
	return inst.result[0]
}

// Mix mixes components
func Mix(ctx context.Context, components ...*wtype.Liquid) *wtype.Liquid {
	return genericMix(ctx, mixer.GenericMix(mixer.MixOptions{
		Inputs: components,
	}))
}

// MixInto mixes components
func MixInto(ctx context.Context, outplate *wtype.Plate, address string, components ...*wtype.Liquid) *wtype.Liquid {
	return genericMix(ctx, mixer.GenericMix(mixer.MixOptions{
		Inputs:      components,
		Destination: outplate,
		Address:     address,
	}))
}

// MixNamed mixes components
func MixNamed(ctx context.Context, outplatetype, address string, platename string, components ...*wtype.Liquid) *wtype.Liquid {
	return genericMix(ctx, mixer.GenericMix(mixer.MixOptions{
		Inputs:    components,
		PlateType: outplatetype,
		Address:   address,
		PlateName: platename,
	}))
}

// MixTo mixes components
//
// TODO: Addresses break dependence information. Deprecated.
func MixTo(ctx context.Context, outplatetype, address string, platenum int, components ...*wtype.Liquid) *wtype.Liquid {
	return genericMix(ctx, mixer.GenericMix(mixer.MixOptions{
		Inputs:    components,
		PlateType: outplatetype,
		Address:   address,
		PlateNum:  platenum,
	}))
}

// SplitSample is essentially an inverse mix: takes one component and a volume and returns two
// the question is then over what happens subsequently.. unlike mix this does not have a
// destination as it's intrinsically a source operation
func SplitSample(ctx context.Context, component *wtype.Liquid, volume wunit.Volume) (removed, remaining *wtype.Liquid) {
	// at this point we cannot guarantee that volumes are accurate
	// so it's a case of best-efforts

	inst := splitSample(ctx, component, volume)

	Issue(ctx, inst)

	// protocol world must not be able to modify the copies seen here
	return inst.result[0].Dup(), inst.result[1].Dup()
}

// Sample takes a sample of volume v from this liquid
func Sample(ctx context.Context, liquid *wtype.Liquid, v wunit.Volume) *wtype.Liquid {
	return mixer.Sample(liquid, v)
}

func splitSample(ctx context.Context, component *wtype.Liquid, volume wunit.Volume) *commandInst {

	split := wtype.NewLHSplitInstruction()

	// this will count as a mix-in-place effectively
	split.Inputs = append(split.Inputs, component.Dup())

	cmpMoving, cmpStaying := mixer.SplitSample(component, volume)

	//the ID of the component that is staying has been updated
	sampletracker.FromContext(ctx).UpdateIDOf(component.ID, cmpStaying.ID)

	split.AddOutput(cmpMoving)
	split.AddOutput(cmpStaying)

	// Create Instruction
	inst := &commandInst{
		Args: []*wtype.Liquid{component},
		Command: &ast.Command{
			Inst: split,
			Request: ast.Request{
				Selector: []ast.NameValue{
					target.DriverSelectorV1Mixer,
				},
			},
		},
		result: []*wtype.Liquid{cmpMoving, cmpStaying},
	}

	return inst
}
