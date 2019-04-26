package codegen

import (
	"github.com/antha-lang/antha/graph"
	"github.com/antha-lang/antha/instructions"
	"github.com/antha-lang/antha/target"
)

// instGraph is a graph to that models dependencies between target instructions
// but uses edges from separately initialized dependsOn field. Useful for cases
// where we can't depend on target.Graph because we are using this to build the
// initial DependsOn relation.
type instGraph struct {
	insts     []instructions.Inst
	added     map[instructions.Inst]bool
	dependsOn map[instructions.Inst][]instructions.Inst
	entry     map[graph.Node]instructions.Inst
	exit      map[graph.Node]instructions.Inst
}

func newInstGraph() *instGraph {
	return &instGraph{
		added:     make(map[instructions.Inst]bool),
		entry:     make(map[graph.Node]instructions.Inst),
		exit:      make(map[graph.Node]instructions.Inst),
		dependsOn: make(map[instructions.Inst][]instructions.Inst),
	}
}

func (a *instGraph) NumNodes() int {
	return len(a.insts)
}

func (a *instGraph) Node(i int) graph.Node {
	return a.insts[i]
}

func (a *instGraph) NumOuts(n graph.Node) int {
	return len(a.dependsOn[n.(instructions.Inst)])
}

func (a *instGraph) Out(n graph.Node, i int) graph.Node {
	return a.dependsOn[n.(instructions.Inst)][i]
}

// addInsts adds instructions to the graph
func (a *instGraph) addInsts(insts instructions.Insts) {
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

func (a *instGraph) addInitializers(insts instructions.Insts) {
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

func (a *instGraph) addFinalizers(insts instructions.Insts) {
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
func (a *instGraph) addRootedInsts(root graph.Node, insts instructions.Insts) {
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

	newInsts := append(instructions.Insts{entry, exit}, insts...)
	a.addInsts(newInsts)
}
