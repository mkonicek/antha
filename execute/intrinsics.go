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

type incubateInst struct {
	Arg  *wtype.LHComponent
	Comp *wtype.LHComponent
	Node *ast.Incubate
}

func Incubate(ctx context.Context, in *wtype.LHComponent, temp wunit.Temperature, time wunit.Time, shaking bool) *wtype.LHComponent {
	st := sampletracker.GetSampleTracker()
	comp := in.Dup()
	comp.ID = wtype.GetUUID()
	comp.BlockID = wtype.NewBlockID(getId(ctx))
	st.UpdateIDOf(in.ID, comp.ID)

	trace.Issue(ctx, &incubateInst{
		Arg:  in,
		Comp: comp,
		Node: &ast.Incubate{
			Time: time,
			Temp: temp,
			Reqs: []ast.Request{
				ast.Request{
					Time: ast.NewPoint(time.SIValue()),
					Temp: ast.NewPoint(temp.SIValue()),
				},
			},
		},
	})
	return comp
}

type mixInst struct {
	Args []*wtype.LHComponent
	Comp *wtype.LHComponent
	Node *ast.Mix
}

// TODO -- LOC etc. will be passed through OK but what about
//         the actual plate info?
//        - two choices here: 1) we upgrade the sample tracker; 2) we pass the plate in somehow
func mix(ctx context.Context, inst *mixInst) *wtype.LHComponent {
	// from the protocol POV components need to be passed by value
	cmps := wtype.CopyComponentArray(inst.Args)

	inst.Node.Inst.BlockID = wtype.NewBlockID(getId(ctx))
	inst.Node.Inst.Result.BlockID = inst.Node.Inst.BlockID
	inst.Comp = inst.Node.Inst.Result
	inst.Comp.BlockID = inst.Node.Inst.BlockID
	mx := 0
	for i, c := range cmps {
		v := c.Volume().SIValue()
		inst.Node.Reqs = append(inst.Node.Reqs, ast.Request{MixVol: ast.NewPoint(v)})
		c.Order = i
		inst.Comp.Mix(c)
		//inst.Comp.AddParent(c.ID)
		//c.AddDaughter(inst.Comp.ID)
		if c.Generation() > mx {
			mx = c.Generation()
		}
	}

	if len(cmps) == 0 {
		inst.Node.Inst.Type = wtype.LHIEND
	} else {
		inst.Node.Inst.Type = wtype.LHIMIX
	}

	inst.Node.Inst.SetGeneration(mx)
	inst.Comp.SetGeneration(mx + 1)

	inst.Node.Inst.ProductID = inst.Comp.ID

	trace.Issue(ctx, inst)

	return inst.Comp
}

func Mix(ctx context.Context, components ...*wtype.LHComponent) *wtype.LHComponent {
	return mix(ctx, &mixInst{
		Args: components,
		Node: &ast.Mix{
			Inst: mixer.GenericMix(mixer.MixOptions{
				Components: components,
			})},
	})
}

func MixInto(ctx context.Context, outplate *wtype.LHPlate, address string, components ...*wtype.LHComponent) *wtype.LHComponent {
	return mix(ctx, &mixInst{
		Args: components,
		Node: &ast.Mix{
			Inst: mixer.GenericMix(mixer.MixOptions{
				Components:  components,
				Destination: outplate,
				Address:     address,
			})},
	})
}

func MixTo(ctx context.Context, outplatetype, address string, platenum int, components ...*wtype.LHComponent) *wtype.LHComponent {
	return mix(ctx, &mixInst{
		Args: components,
		Node: &ast.Mix{
			Inst: mixer.GenericMix(mixer.MixOptions{
				Components: components,
				PlateType:  outplatetype,
				Address:    address,
				PlateNum:   platenum,
			})},
	})
}
