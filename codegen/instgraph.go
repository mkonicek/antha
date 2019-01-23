package codegen

import (
	"github.com/antha-lang/antha/graph"
	"github.com/antha-lang/antha/laboratory/effects"
	"github.com/antha-lang/antha/target"
)

// instGraph is a graph to that models dependencies between target instructions
// but uses edges from separately initialized dependsOn field. Useful for cases
// where we can't depend on target.Graph because we are using this to build the
// initial DependsOn relation.
type instGraph struct {
	insts     []effects.Inst
	added     map[effects.Inst]bool
	dependsOn map[effects.Inst][]effects.Inst
	entry     map[graph.Node]effects.Inst
	exit      map[graph.Node]effects.Inst
}

func newInstGraph() *instGraph {
	return &instGraph{
		added:     make(map[effects.Inst]bool),
		entry:     make(map[graph.Node]effects.Inst),
		exit:      make(map[graph.Node]effects.Inst),
		dependsOn: make(map[effects.Inst][]effects.Inst),
	}
}

func (a *instGraph) NumNodes() int {
	return len(a.insts)
}

func (a *instGraph) Node(i int) graph.Node {
	return a.insts[i]
}

func (a *instGraph) NumOuts(n graph.Node) int {
	return len(a.dependsOn[n.(effects.Inst)])
}

func (a *instGraph) Out(n graph.Node, i int) graph.Node {
	return a.dependsOn[n.(effects.Inst)][i]
}

// addInsts adds instructions to the graph
func (a *instGraph) addInsts(insts effects.Insts) {
	// Add dependencies
	for _, in := range insts {
		a.dependsOn[in] = append(a.dependsOn[in], in.DependsOn()...)
	}

	// Add nodes
	for _, in := range insts {
		if a.added[in] {
			continue
		}
		a.added[in] = true
		a.insts = append(a.insts, in)

		for _, v := range a.dependsOn[in] {
			if a.added[v] {
				continue
			}
			a.added[v] = true

			a.insts = append(a.insts, v)
		}
	}
}

func (a *instGraph) addInitializers(insts effects.Insts) {
	if len(insts) == 0 {
		return
	}

	insts.SequentialOrder()
	last := insts[len(insts)-1]
	for _, inst := range a.entry {
		inst.AppendDependsOn(last)
		// Unlike other cases, inst has already been added to graph, so update
		// data explicitly
		a.dependsOn[inst] = append(a.dependsOn[inst], last)
	}
	a.addInsts(insts)
}

func (a *instGraph) addFinalizers(insts effects.Insts) {
	if len(insts) == 0 {
		return
	}

	insts.SequentialOrder()
	first := insts[0]
	for _, inst := range a.exit {
		first.AppendDependsOn(inst)
	}
	a.addInsts(insts)
}

// addRootedInsts adds instructions that correspond to a particular graph Node
func (a *instGraph) addRootedInsts(root graph.Node, insts effects.Insts) {
	exit := &target.Wait{}
	entry := &target.Wait{}

	a.entry[root] = entry
	a.exit[root] = exit

	if len(insts) != 0 {
		first := insts[0]
		last := insts[len(insts)-1]

		first.AppendDependsOn(entry)
		exit.AppendDependsOn(last)
	}

	newInsts := append(effects.Insts{entry, exit}, insts...)
	a.addInsts(newInsts)
}
