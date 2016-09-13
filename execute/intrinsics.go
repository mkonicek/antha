package execute

import (
	"github.com/antha-lang/antha/antha/anthalib/mixer"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/ast"
	"github.com/antha-lang/antha/bvendor/golang.org/x/net/context"
	"github.com/antha-lang/antha/microArch/sampletracker"
	"github.com/antha-lang/antha/trace"
)

type commandInst struct {
	Args    []*wtype.LHComponent
	Comp    *wtype.LHComponent
	Command *ast.Command
}

func SetInputPlate(ctx context.Context, plate *wtype.LHPlate) {
	st := sampletracker.GetSampleTracker()
	st.SetInputPlate(plate)
}

func incubate(ctx context.Context, in *wtype.LHComponent, temp wunit.Temperature, time wunit.Time, shaking bool) *commandInst {
	st := sampletracker.GetSampleTracker()
	comp := in.Dup()
	comp.ID = wtype.GetUUID()
	comp.BlockID = wtype.NewBlockID(getId(ctx))
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

func Incubate(ctx context.Context, in *wtype.LHComponent, temp wunit.Temperature, time wunit.Time, shaking bool) *wtype.LHComponent {
	inst := incubate(ctx, in, temp, time, shaking)
	trace.Issue(ctx, inst)
	return inst.Comp
}

func handle(ctx context.Context, by string, in *wtype.LHComponent) *commandInst {
	st := sampletracker.GetSampleTracker()
	comp := in.Dup()
	comp.ID = wtype.GetUUID()
	comp.BlockID = wtype.NewBlockID(getId(ctx))
	st.UpdateIDOf(in.ID, comp.ID)

	return &commandInst{
		Args: []*wtype.LHComponent{in},
		Comp: comp,
		Command: &ast.Command{
			Inst: &ast.HandleInst{
				Group: by,
			},
			Requests: []ast.Request{ast.Request{Manual: true}},
		},
	}
}

func Handle(ctx context.Context, by string, in *wtype.LHComponent) *wtype.LHComponent {
	inst := handle(ctx, by, in)
	trace.Issue(ctx, inst)
	return inst.Comp
}

// TODO -- LOC etc. will be passed through OK but what about
//         the actual plate info?
//        - two choices here: 1) we upgrade the sample tracker; 2) we pass the plate in somehow
func mix(ctx context.Context, inst *wtype.LHInstruction) *commandInst {
	inst.BlockID = wtype.NewBlockID(getId(ctx))
	inst.Result.BlockID = inst.BlockID

	result := inst.Result
	result.BlockID = inst.BlockID

	mx := 0
	var reqs []ast.Request
	// from the protocol POV components need to be passed by value
	for i, c := range wtype.CopyComponentArray(inst.Components) {
		reqs = append(reqs, ast.Request{MixVol: ast.NewPoint(c.Volume().SIValue())})
		c.Order = i
		result.MixPreserveTvol(c)
		if c.Generation() > mx {
			mx = c.Generation()
		}
	}

	inst.SetGeneration(mx)
	result.SetGeneration(mx + 1)

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

func Mix(ctx context.Context, components ...*wtype.LHComponent) *wtype.LHComponent {
	return genericMix(ctx, mixer.GenericMix(mixer.MixOptions{
		Components: components,
	}))
}

func MixInto(ctx context.Context, outplate *wtype.LHPlate, address string, components ...*wtype.LHComponent) *wtype.LHComponent {
	return genericMix(ctx, mixer.GenericMix(mixer.MixOptions{
		Components:  components,
		Destination: outplate,
		Address:     address,
	}))
}

func MixNamed(ctx context.Context, outplatetype, address string, platename string, components ...*wtype.LHComponent) *wtype.LHComponent {
	return genericMix(ctx, mixer.GenericMix(mixer.MixOptions{
		Components: components,
		PlateType:  outplatetype,
		Address:    address,
		PlateName:  platename,
	}))
}

func MixTo(ctx context.Context, outplatetype, address string, platenum int, components ...*wtype.LHComponent) *wtype.LHComponent {
	return genericMix(ctx, mixer.GenericMix(mixer.MixOptions{
		Components: components,
		PlateType:  outplatetype,
		Address:    address,
		PlateNum:   platenum,
	}))
}
