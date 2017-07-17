package human

import (
	"reflect"

	"github.com/antha-lang/antha/ast"
	"github.com/antha-lang/antha/graph"
	"github.com/antha-lang/antha/target"
)

const (
	// HumanByHumanCost is the cost of manually moving from another human device
	HumanByHumanCost = 50
	// HumanByXCost is the cost of manually moving from any non-human device
	HumanByXCost = 100
)

var (
	_ target.Device = &Human{}
)

// A Human is a device that can do anything
type Human struct {
	opt Opt
}

// CanCompile implements device CanCompile
func (a *Human) CanCompile(req ast.Request) bool {
	can := ast.Request{
		Move: req.Move,
		Selector: []ast.NameValue{
			ast.NameValue{
				// TODO: Remove hard coded strings
				Name:  "antha.driver.v1.TypeReply.type",
				Value: "antha.human.v1.Human",
			},
		},
	}
	if a.opt.CanMix {
		can.MixVol = req.MixVol
	}
	if a.opt.CanIncubate {
		can.Temp = req.Temp
		can.Time = req.Time
	}
	if a.opt.CanHandle {
		can.Selector = req.Selector
	}

	if !req.Matches(can) {
		return false
	}

	return can.Contains(req)
}

// MoveCost implements target.device MoveCost
func (a *Human) MoveCost(from target.Device) int {
	if _, ok := from.(*Human); ok {
		return HumanByHumanCost
	}
	return HumanByXCost
}

func (a *Human) String() string {
	return "Human"
}

// Return key for node for grouping
func getKey(n ast.Node) (r interface{}) {
	// Group by value for HandleInst and Incubate and type otherwise
	if c, ok := n.(*ast.Command); !ok {
		r = reflect.TypeOf(n)
	} else if h, ok := c.Inst.(*ast.HandleInst); ok {
		r = h.Group
	} else if i, ok := c.Inst.(*ast.IncubateInst); ok {
		r = i.Temp.ToString() + " " + i.Time.ToString()
	} else {
		r = reflect.TypeOf(c.Inst)
	}
	return
}

func keepNodes(keep map[graph.Node]bool, dag *graph.SDag) (roots []graph.Node) {
	next := dag.Roots
	for len(next) > 0 {
		last := len(next) - 1
		r := next[last]
		next = next[:last]

		if keep[r] {
			roots = append(roots, r)
			continue
		}

		next = append(next, dag.Visit(r)...)
	}

	return
}

func orderNodes(keep map[graph.Node]target.Inst, g graph.Graph) (ret []graph.Node) {
	order, _ := graph.TopoSort(graph.TopoSortOpt{
		Graph: g,
	})

	for _, n := range order {
		if keep[n] != nil {
			ret = append(ret, n)
		}
	}

	return
}

// Compile implements target.device Compile
func (a *Human) Compile(nodes []ast.Node) ([]target.Inst, error) {
	addDep := func(in, dep target.Inst) {
		in.SetDependsOn(append(in.DependsOn(), dep))
	}

	g := ast.ToGraph(ast.ToGraphOpt{
		Roots:     nodes,
		WhichDeps: ast.DataDeps,
	})
	isRoot := make(map[graph.Node]bool)
	for _, n := range nodes {
		isRoot[n] = true
	}

	entry := &target.Wait{}
	exit := &target.Wait{}
	var insts []target.Inst
	isRep := make(map[graph.Node]target.Inst)

	insts = append(insts, entry)

	// Maximally coalesce repeated commands according to when they are first
	// available to be executed (graph.Reverse)
	dag := graph.Schedule(graph.Reverse(g))
	for len(dag.Roots) > 0 {
		roots := keepNodes(isRoot, dag)

		var next []graph.Node

		// Gather
		same := make(map[interface{}][]ast.Node)
		for _, r := range roots {
			n := r.(ast.Node)
			key := getKey(n)
			same[key] = append(same[key], n)
			next = append(next, dag.Visit(r)...)
		}
		// Apply
		for _, nodes := range same {
			var ins []*target.Manual
			for _, n := range nodes {
				in, err := a.makeInst(n)
				if err != nil {
					return nil, err
				}
				ins = append(ins, in)
			}
			in := a.coalesce(ins)
			insts = append(insts, in)

			// Pick a representative
			isRep[nodes[0]] = in
		}

		dag.Roots = next
	}

	insts = append(insts, exit)

	order := orderNodes(isRep, g)
	for idx, node := range order {
		n := node.(ast.Node)
		in := isRep[n]

		if idx-1 >= 0 {
			prev := order[idx-1].(ast.Node)
			prevNode := isRep[prev]
			if in != prevNode {
				addDep(in, prevNode)
			}
		}

		addDep(in, entry)
		addDep(exit, in)
	}

	return insts, nil
}

// An Opt is a set of options to configure a human device
type Opt struct {
	CanMix      bool
	CanIncubate bool
	CanHandle   bool
}

// New returns a new human device
func New(opt Opt) *Human {
	return &Human{opt}
}
