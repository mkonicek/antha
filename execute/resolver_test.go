package execute

import (
	"context"
	"testing"

	"github.com/antha-lang/antha/antha/anthalib/mixer"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/graph"
	"github.com/antha-lang/antha/target"
	"github.com/antha-lang/antha/target/human"
)

func TestUseCompChainThroughSample(t *testing.T) {
	tgt := target.New()
	tgt.AddDevice(human.New(human.Opt{CanMix: true}))

	ctx := context.Background()
	ctx = target.WithTarget(ctx, tgt)
	ctx = withID(ctx, "")

	vol := wunit.NewVolume(1, "ul")
	cmp := wtype.NewLHComponent()
	cmp.CName = "thiscannotbeomitted"
	a1 := mix(ctx, mixer.GenericMix(mixer.MixOptions{
		Components: []*wtype.LHComponent{cmp},
	}))
	a2 := mix(ctx, mixer.GenericMix(mixer.MixOptions{
		Components: []*wtype.LHComponent{mixer.Sample(a1.result[0], vol)},
	}))
	a3 := mix(ctx, mixer.GenericMix(mixer.MixOptions{
		Components: []*wtype.LHComponent{mixer.Sample(a2.result[0], vol)},
	}))

	var insts []interface{}
	insts = append(insts, a1, a2, a3)
	r := &resolver{}
	if _, err := r.resolve(ctx, insts); err != nil {
		t.Fatal(err)
	} else if len(r.insts) == 0 {
		t.Fatalf("no instructions")
	}

	g := graph.Reverse(&target.Graph{
		Insts: r.insts,
	})
	dists := graph.ShortestPath(graph.ShortestPathOpt{
		Graph:   g,
		Sources: []graph.Node{r.insts[0]},
		Weight: func(x, y graph.Node) int {
			return 1
		},
	})
	max := 0
	for _, d := range dists {
		if d > max {
			max = d
		}
	}
	if e, f := 2, max; e != f {
		t.Errorf("expected graph depth of %d found %d", e, f)
	}
}
