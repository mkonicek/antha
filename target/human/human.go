package human

import (
	"reflect"

	"github.com/antha-lang/antha/ast"
	"github.com/antha-lang/antha/graph"
	"github.com/antha-lang/antha/target"
)

const (
	HumanByHumanCost = 50  // Cost of manually moving from another human device
	HumanByXCost     = 100 // Cost of manually moving from any non-human device
)

var (
	_ target.Device = &Human{}
)

type Human struct {
	opt Opt
}

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

func (a *Human) Compile(nodes []ast.Node) ([]target.Inst, error) {
	addDep := func(in, dep target.Inst) {
		in.SetDependsOn(append(in.DependsOn(), dep))
	}

	g := ast.Deps(nodes)

	entry := &target.Wait{}
	exit := &target.Wait{}
	var insts []target.Inst
	inst := make(map[ast.Node]target.Inst)

	insts = append(insts, entry)

	// Maximally coalesce repeated commands according to when they are first
	// available to be executed (graph.Reverse)
	dag := graph.Schedule(graph.Reverse(g))
	for len(dag.Roots) > 0 {
		var next []graph.Node
		// Gather
		same := make(map[interface{}][]ast.Node)
		for _, r := range dag.Roots {
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

			for _, n := range nodes {
				inst[n] = in
			}
		}

		dag.Roots = next
	}

	insts = append(insts, exit)

	for i, inum := 0, g.NumNodes(); i < inum; i += 1 {
		n := g.Node(i).(ast.Node)
		in := inst[n]
		for j, jnum := 0, g.NumOuts(n); j < jnum; j += 1 {
			kid := g.Out(n, j).(ast.Node)
			kidIn := inst[kid]
			addDep(in, kidIn)
		}
		addDep(in, entry)
		addDep(exit, in)
	}

	return insts, nil
}

type Opt struct {
	CanMix      bool
	CanIncubate bool
	CanHandle   bool
}

func New(opt Opt) *Human {
	return &Human{opt}
}
