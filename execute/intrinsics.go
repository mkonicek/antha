package execute

import (
	"context"

	"github.com/antha-lang/antha/antha/anthalib/mixer"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/ast"
	"github.com/antha-lang/antha/driver"
	"github.com/antha-lang/antha/inventory"
	"github.com/antha-lang/antha/microArch/sampletracker"
	"github.com/antha-lang/antha/target"
	"github.com/antha-lang/antha/trace"
)

type commandInst struct {
	Args    []*wtype.LHComponent
	Comp    []*wtype.LHComponent
	Command *ast.Command
}

// SetInputPlate notifies the planner about an input plate
func SetInputPlate(ctx context.Context, plate *wtype.LHPlate) {
	st := sampletracker.GetSampleTracker()
	st.SetInputPlate(plate)
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

func newCompFromComp(ctx context.Context, in *wtype.LHComponent) *wtype.LHComponent {
	st := sampletracker.GetSampleTracker()
	comp := in.Dup()
	comp.ID = wtype.GetUUID()
	comp.BlockID = wtype.NewBlockID(getID(ctx))
	comp.SetGeneration(comp.Generation() + 1)

	getMaker(ctx).UpdateAfterInst(in.ID, comp.ID)
	st.UpdateIDOf(in.ID, comp.ID)
	return comp
}

// Incubate incubates a component
func Incubate(ctx context.Context, in *wtype.LHComponent, opt IncubateOpt) *wtype.LHComponent {
	inst := &commandInst{
		Args: []*wtype.LHComponent{in},
		Comp: []*wtype.LHComponent{newCompFromComp(ctx, in)},
		Command: &ast.Command{
			Inst: &ast.IncubateInst{
				Time:           opt.Time,
				Temp:           opt.Temp,
				ShakeRate:      opt.ShakeRate,
				ShakeRadius:    opt.ShakeRadius,
				PreTemp:        opt.PreTemp,
				PreTime:        opt.PreTime,
				PreShakeRate:   opt.PreShakeRate,
				PreShakeRadius: opt.PreShakeRadius,
			},
		},
	}

	// TODO: revisit when ast.Request architecture is removed as this command
	// cannot be assigned independently. It needs to be linked with a previous
	// Incubate. For now assume just one incubator and use explicit selector
	inst.Command.Requests = append(inst.Command.Requests, ast.Request{
		Selector: []ast.NameValue{
			target.DriverSelectorV1ShakerIncubator,
		},
	})

	trace.Issue(ctx, inst)
	return inst.Comp[0]
}

// prompt... works pretty much like Handle does
// but passes the instruction to the planner
// in future this should generate handles as side-effects

type mixerPromptOpts struct {
	Component   *wtype.LHComponent
	ComponentIn *wtype.LHComponent
	Message     string
}

// MixerPrompt prompts user with a message during mixer execution
func MixerPrompt(ctx context.Context, in *wtype.LHComponent, message string) *wtype.LHComponent {
	inst := mixerPrompt(ctx,
		mixerPromptOpts{
			Component:   newCompFromComp(ctx, in),
			ComponentIn: in,
			Message:     message,
		},
	)
	trace.Issue(ctx, inst)
	return inst.Comp[0]
}

// Prompt prompts user with a message
func Prompt(ctx context.Context, in *wtype.LHComponent, message string) *wtype.LHComponent {
	inst := &commandInst{
		Args: []*wtype.LHComponent{in},
		Comp: []*wtype.LHComponent{newCompFromComp(ctx, in)},
		Command: &ast.Command{
			Inst: &ast.PromptInst{
				Message: message,
			},
		},
	}

	inst.Command.Requests = append(inst.Command.Requests, ast.Request{
		Selector: []ast.NameValue{
			target.DriverSelectorV1Human,
		},
	})

	trace.Issue(ctx, inst)
	return inst.Comp[0]
}

func mixerPrompt(ctx context.Context, opts mixerPromptOpts) *commandInst {
	inst := wtype.NewLHPromptInstruction()
	inst.SetGeneration(opts.ComponentIn.Generation())
	inst.Message = opts.Message
	inst.AddProduct(opts.Component)
	inst.AddComponent(opts.ComponentIn)
	inst.PassThrough[opts.ComponentIn.ID] = opts.Component

	return &commandInst{
		Args: []*wtype.LHComponent{opts.ComponentIn},
		Comp: []*wtype.LHComponent{opts.Component},
		Command: &ast.Command{
			Inst: inst,
			Requests: []ast.Request{
				ast.Request{
					Selector: []ast.NameValue{
						target.DriverSelectorV1Prompter,
					},
				},
			},
		},
	}
}

func handle(ctx context.Context, opt HandleOpt) *commandInst {
	comp := newCompFromComp(ctx, opt.Component)

	var sels []ast.NameValue

	if len(opt.Selector) == 0 {
		sels = append(sels, target.DriverSelectorV1Human)
	} else {
		for n, v := range opt.Selector {
			sels = append(sels, ast.NameValue{Name: n, Value: v})
		}
	}

	return &commandInst{
		Args: []*wtype.LHComponent{opt.Component},
		Comp: []*wtype.LHComponent{comp},
		Command: &ast.Command{
			Inst: &ast.HandleInst{
				Group:    opt.Label,
				Selector: opt.Selector,
				Calls:    opt.Calls,
			},
			Requests: []ast.Request{ast.Request{Selector: sels}},
		},
	}
}

// A HandleOpt are options to Handle
type HandleOpt struct {
	Component *wtype.LHComponent
	Label     string
	Selector  map[string]string
	Calls     []driver.Call
}

// Handle performs a low level instruction on a component
func Handle(ctx context.Context, opt HandleOpt) *wtype.LHComponent {
	inst := handle(ctx, opt)
	trace.Issue(ctx, inst)
	return inst.Comp[0]
}

// NewComponent returns a new component given a component type
func NewComponent(ctx context.Context, typ string) *wtype.LHComponent {
	c, err := inventory.NewComponent(ctx, typ)
	if err != nil {
		Errorf(ctx, "cannot make component %s: %s", typ, err)
	}
	return c
}

// NewPlate returns a new plate given a plate type
func NewPlate(ctx context.Context, typ string) *wtype.LHPlate {
	p, err := inventory.NewPlate(ctx, typ)
	if err != nil {
		Errorf(ctx, "cannot make plate %s: %s", typ, err)
	}
	return p
}

// TODO -- LOC etc. will be passed through OK but what about
//         the actual plate info?
//        - two choices here: 1) we upgrade the sample tracker; 2) we pass the plate in somehow
func mix(ctx context.Context, inst *wtype.LHInstruction) *commandInst {
	inst.BlockID = wtype.NewBlockID(getID(ctx))
	inst.Result.BlockID = inst.BlockID

	result := inst.Result
	result.BlockID = inst.BlockID

	mx := 0
	var reqs []ast.Request
	// from the protocol POV components need to be passed by value
	for i, c := range wtype.CopyComponentArray(inst.Components) {
		reqs = append(reqs, ast.Request{
			Selector: []ast.NameValue{
				target.DriverSelectorV1Mixer,
			},
		})
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
	inst.ProductID = result.ID

	return &commandInst{
		Args: inst.Components,
		Command: &ast.Command{
			Requests: reqs,
			Inst:     inst,
		},
		Comp: []*wtype.LHComponent{result},
	}
}

func genericMix(ctx context.Context, generic *wtype.LHInstruction) *wtype.LHComponent {
	inst := mix(ctx, generic)
	trace.Issue(ctx, inst)
	return inst.Comp[0]
}

// Mix mixes components
func Mix(ctx context.Context, components ...*wtype.LHComponent) *wtype.LHComponent {
	return genericMix(ctx, mixer.GenericMix(mixer.MixOptions{
		Components: components,
	}))
}

// MixInto mixes components
func MixInto(ctx context.Context, outplate *wtype.LHPlate, address string, components ...*wtype.LHComponent) *wtype.LHComponent {
	return genericMix(ctx, mixer.GenericMix(mixer.MixOptions{
		Components:  components,
		Destination: outplate,
		Address:     address,
	}))
}

// MixNamed mixes components
func MixNamed(ctx context.Context, outplatetype, address string, platename string, components ...*wtype.LHComponent) *wtype.LHComponent {
	return genericMix(ctx, mixer.GenericMix(mixer.MixOptions{
		Components: components,
		PlateType:  outplatetype,
		Address:    address,
		PlateName:  platename,
	}))
}

// MixTo mixes components
//
// TODO: Addresses break dependence information. Deprecated.
func MixTo(ctx context.Context, outplatetype, address string, platenum int, components ...*wtype.LHComponent) *wtype.LHComponent {
	return genericMix(ctx, mixer.GenericMix(mixer.MixOptions{
		Components: components,
		PlateType:  outplatetype,
		Address:    address,
		PlateNum:   platenum,
	}))
}
