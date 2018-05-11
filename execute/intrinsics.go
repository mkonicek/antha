package execute

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/antha-lang/antha/antha/anthalib/mixer"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	api "github.com/antha-lang/antha/api/v1"
	"github.com/antha-lang/antha/ast"
	"github.com/antha-lang/antha/driver"
	"github.com/antha-lang/antha/inject"
	"github.com/antha-lang/antha/inventory"
	"github.com/antha-lang/antha/microArch/sampletracker"
	"github.com/antha-lang/antha/target"
	"github.com/antha-lang/antha/trace"
)

// a commandInst is a generic intrinsic instruction
type commandInst struct {
	// Arguments to this command. Used to determine command dependencies.
	Args []*wtype.LHComponent
	// Components created by this command. Returned back to user code
	result  []*wtype.LHComponent
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
		Args:   []*wtype.LHComponent{in},
		result: []*wtype.LHComponent{newCompFromComp(ctx, in)},
		Command: &ast.Command{
			Inst: innerInst,
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
	return inst.result[0]
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
	return inst.result[0]
}

// Prompt prompts user with a message
func Prompt(ctx context.Context, in *wtype.LHComponent, message string) *wtype.LHComponent {
	inst := &commandInst{
		Args:   []*wtype.LHComponent{in},
		result: []*wtype.LHComponent{newCompFromComp(ctx, in)},
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
	return inst.result[0]
}

func mixerPrompt(ctx context.Context, opts mixerPromptOpts) *commandInst {
	inst := wtype.NewLHPromptInstruction()
	inst.SetGeneration(opts.ComponentIn.Generation())
	inst.Message = opts.Message
	inst.AddProduct(opts.Component)
	inst.AddComponent(opts.ComponentIn)
	inst.PassThrough[opts.ComponentIn.ID] = opts.Component

	return &commandInst{
		Args:   []*wtype.LHComponent{opts.ComponentIn},
		result: []*wtype.LHComponent{opts.Component},
		Command: &ast.Command{
			Inst: inst,
			Requests: []ast.Request{
				{
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
		Args:   []*wtype.LHComponent{opt.Component},
		result: []*wtype.LHComponent{comp},
		Command: &ast.Command{
			Inst: &ast.HandleInst{

				Calls: opt.Calls,
			},
			Requests: []ast.Request{{Selector: sels}},
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
	return inst.result[0]
}

// PlateReadOpts defines plate-reader absorbance options
type PlateReadOpts struct {
	Sample  *wtype.LHComponent
	Options string
}

func readPlate(ctx context.Context, opts PlateReadOpts) *commandInst {
	inst := wtype.NewPRInstruction()
	inst.ComponentIn = opts.Sample

	// Clone the component to represent the result of the AbsorbanceRead
	inst.ComponentOut = newCompFromComp(ctx, opts.Sample)
	inst.Options = opts.Options

	return &commandInst{
		Args:   []*wtype.LHComponent{opts.Sample},
		result: []*wtype.LHComponent{inst.ComponentOut},
		Command: &ast.Command{
			Inst: inst,
			Requests: []ast.Request{
				{
					Selector: []ast.NameValue{
						target.DriverSelectorV1WriteOnlyPlateReader,
					},
				},
			},
		},
	}
}

// PlateRead reads absorbance of a component
func PlateRead(ctx context.Context, opt PlateReadOpts) *wtype.LHComponent {
	inst := readPlate(ctx, opt)
	trace.Issue(ctx, inst)
	return inst.result[0]
}

// QPCROptions are the options for a QPCR request.
type QPCROptions struct {
	Reactions  []*wtype.LHComponent
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
	inst.ComponentOut = []*wtype.LHComponent{}

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
	}
}

// RunQPCRExperiment starts a new QPCR experiment, using an experiment input file.
func RunQPCRExperiment(ctx context.Context, opt QPCROptions) []*wtype.LHComponent {
	inst := runQPCR(ctx, opt, "RunExperiment")
	trace.Issue(ctx, inst)
	return inst.result
}

// RunQPCRFromTemplate starts a new QPCR experiment, using a template input file.
func RunQPCRFromTemplate(ctx context.Context, opt QPCROptions) []*wtype.LHComponent {
	inst := runQPCR(ctx, opt, "RunExperimentFromTemplate")
	trace.Issue(ctx, inst)
	return inst.result
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

func mix(ctx context.Context, inst *wtype.LHInstruction) *commandInst {
	inst.BlockID = wtype.NewBlockID(getID(ctx))
	inst.Results[0].BlockID = inst.BlockID
	result := inst.Results[0]
	//result.BlockID = inst.BlockID // DELETEME

	mx := 0
	var reqs []ast.Request
	// from the protocol POV components need to be passed by value
	for i, c := range wtype.CopyComponentArray(inst.Components) {
		if c.CName == "" {
			panic("Nameless Component used in Mix - this is not permitted")
		}
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

	return &commandInst{
		Args: inst.Components,
		Command: &ast.Command{
			Requests: reqs,
			Inst:     inst,
		},
		result: []*wtype.LHComponent{result},
	}
}

func genericMix(ctx context.Context, generic *wtype.LHInstruction) *wtype.LHComponent {
	inst := mix(ctx, generic)
	trace.Issue(ctx, inst)
	if generic.Welladdress != "" {
		err := inst.result[0].SetWellLocation(generic.Welladdress)
		if err != nil {
			panic(err)
		}
	}
	return inst.result[0]
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

// SplitSample is essentially an inverse mix: takes one component and a volume and returns two
// the question is then over what happens subsequently.. unlike mix this does not have a
// destination as it's intrinsically a source operation
func SplitSample(ctx context.Context, component *wtype.LHComponent, volume wunit.Volume) (removed, remaining *wtype.LHComponent) {
	// at this point we cannot guarantee that volumes are accurate
	// so it's a case of best-efforts

	inst := splitSample(ctx, component, volume)

	trace.Issue(ctx, inst)

	// protocol world must not be able to modify the copies seen here
	return inst.result[0].Dup(), inst.result[1].Dup()
}

func splitSample(ctx context.Context, component *wtype.LHComponent, volume wunit.Volume) *commandInst {
	split := wtype.NewLHSplitInstruction()

	// this will count as a mix-in-place effectively
	split.Components = append(split.Components, component.Dup())

	cmpMoving, cmpStaying := mixer.SplitSample(component, volume)

	split.Results = append(split.Results, cmpMoving)
	split.Results = append(split.Results, cmpStaying)

	// Create Instruction
	inst := &commandInst{
		Args: []*wtype.LHComponent{component},
		Command: &ast.Command{
			Requests: []ast.Request{{
				Selector: []ast.NameValue{
					target.DriverSelectorV1Mixer,
				}},
			},
			Inst: split,
		},
		result: []*wtype.LHComponent{cmpMoving, cmpStaying},
	}

	return inst
}

// AwaitData breaks execution pending return of requested data
func AwaitData(
	ctx context.Context,
	object Annotatable,
	meta *api.DeviceMetadata,
	nextElement, replaceParam string,
	nextInput, currentOutput inject.Value) {

	if err := awaitData(ctx, object, meta, nextElement, replaceParam, nextInput, currentOutput); err != nil {
		panic(err)
	}
}

func clone(object inject.Value) (inject.Value, error) {
	bs, err := json.Marshal(object)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(bs, &object); err != nil {
		return nil, err
	}

	return object, nil
}

func awaitData(
	ctx context.Context,
	object Annotatable,
	meta *api.DeviceMetadata,
	nextElement, replaceParam string,
	nextInput, currentOutput inject.Value) error {

	switch t := object.(type) {
	case *wtype.LHPlate:
	default:
		return fmt.Errorf("cannot wait for data on %v type, only LHPlate allowed", t)
	}

	nextInput, err := clone(nextInput)
	if err != nil {
		return err
	}

	currentOutput, err = clone(currentOutput)
	if err != nil {
		return err
	}

	// Get Data Request
	req := ast.Request{
		Selector: []ast.NameValue{
			target.DriverSelectorV1DataSource,
		},
	}

	// Update all components
	plate := object.(*wtype.LHPlate)

	allComp := plate.AllContents()

	var updatedComp []*wtype.LHComponent
	for _, c := range allComp {
		updatedComp = append(updatedComp, newCompFromComp(ctx, c))
	}

	_ = updatedComp // currently unused

	await := &ast.AwaitInst{
		AwaitID:              plate.ID,
		NextElement:          nextElement,
		NextElementInput:     nextInput,
		ReplaceParam:         replaceParam,
		CurrentElementOutput: currentOutput,
	}

	if meta != nil {
		await.Tags = meta.Tags
	}

	// Create Instruction
	inst := &commandInst{
		Args: allComp,
		Command: &ast.Command{
			Requests: []ast.Request{req},
			Inst:     await,
		},
		result: updatedComp,
	}

	trace.Issue(ctx, inst)
	return nil
}
