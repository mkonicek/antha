package execute

import (
	"context"
	"fmt"

	"github.com/antha-lang/antha/antha/anthalib/mixer"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/ast"
	"github.com/antha-lang/antha/driver"
	"github.com/antha-lang/antha/inventory"
	"github.com/antha-lang/antha/microArch/sampletracker"
	"github.com/antha-lang/antha/trace"
)

type commandInst struct {
	Args    []*wtype.LHComponent
	Comp    *wtype.LHComponent
	Command *ast.Command
}

// SetInputPlate notifies the planner about an input plate
func SetInputPlate(ctx context.Context, plate *wtype.LHPlate) {
	st := sampletracker.GetSampleTracker()
	st.SetInputPlate(plate)
}

func incubate(ctx context.Context, in *wtype.LHComponent, temp wunit.Temperature, time wunit.Time, shaking bool) *commandInst {
	st := sampletracker.GetSampleTracker()
	comp := in.Dup()
	comp.ID = wtype.GetUUID()
	comp.BlockID = wtype.NewBlockID(getID(ctx))

	getMaker(ctx).UpdateAfterInst(in.ID, comp.ID)
	fmt.Println("UPDATING HERE (1)")
	st.UpdateIDOf(in.ID, comp.ID)

	return &commandInst{
		Args: []*wtype.LHComponent{in},
		Comp: comp,
		Command: &ast.Command{
			Inst: &ast.IncubateInst{
				Time: time,
				Temp: temp,
			},
			Requests: []ast.Request{
				ast.Request{
					Time: ast.NewPoint(time.SIValue()),
					Temp: ast.NewPoint(temp.SIValue()),
				},
			},
		},
	}
}

// Incubate incubates a component
func Incubate(ctx context.Context, in *wtype.LHComponent, temp wunit.Temperature, time wunit.Time, shaking bool) *wtype.LHComponent {
	inst := incubate(ctx, in, temp, time, shaking)
	trace.Issue(ctx, inst)
	return inst.Comp
}

// prompt... works pretty much like Handle does
// but passes the instruction to the planner
// in future this should generate handles as side-effects

type PromptOpts struct {
	Component   *wtype.LHComponent
	ComponentIn *wtype.LHComponent
	Message     string
}

func Prompt(ctx context.Context, component *wtype.LHComponent, message string) *wtype.LHComponent {
	// sadly need to update everything
	comp := component.Dup()
	comp.ID = wtype.GetUUID()
	comp.BlockID = wtype.NewBlockID(getId(ctx))
	comp.SetGeneration(comp.Generation() + 1)
	getMaker(ctx).UpdateAfterInst(component.ID, comp.ID)
	pinst := prompt(ctx, PromptOpts{Component: comp, ComponentIn: component, Message: message})
	trace.Issue(ctx, pinst)
	return component
}

func prompt(ctx context.Context, opts PromptOpts) *commandInst {
	inst := wtype.NewLHMixInstruction()
	inst.SetGeneration(opts.ComponentIn.Generation())
	inst.Message = opts.Message

	// do we update and return component as with Handle?!
	// this requires fixing the issue with id tracking...
	// will aim for this as a stretch goal

	cp := true
	return &commandInst{
		Args: []*wtype.LHComponent{opts.ComponentIn},
		Comp: opts.Component,
		Command: &ast.Command{
			Inst:     inst,
			Requests: []ast.Request{ast.Request{CanPrompt: &cp}},
		},
	}
}

func handle(ctx context.Context, opt HandleOpt) *commandInst {
	st := sampletracker.GetSampleTracker()
	in := opt.Component
	comp := in.Dup()
	comp.ID = wtype.GetUUID()
	comp.BlockID = wtype.NewBlockID(getID(ctx))

	getMaker(ctx).UpdateAfterInst(in.ID, comp.ID)
	fmt.Println("HANDLE ", opt.Label, "UPDATING HERE (2)", in.CName)
	st.UpdateIDOf(in.ID, comp.ID)

	var sels []ast.NameValue

	if len(opt.Selector) == 0 {
		sels = append(sels, ast.NameValue{
			Name:  "antha.driver.v1.TypeReply.type",
			Value: "antha.human.v1.Human",
		})
	} else {
		for n, v := range opt.Selector {
			sels = append(sels, ast.NameValue{Name: n, Value: v})
		}
	}

	return &commandInst{
		Args: []*wtype.LHComponent{in},
		Comp: comp,
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
	return inst.Comp
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
		reqs = append(reqs, ast.Request{MixVol: ast.NewPoint(c.Volume().SIValue())})
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
		Comp: result,
	}
}

func genericMix(ctx context.Context, generic *wtype.LHInstruction) *wtype.LHComponent {
	inst := mix(ctx, generic)
	trace.Issue(ctx, inst)
	return inst.Comp
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
