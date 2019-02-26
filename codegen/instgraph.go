package codegen

import (
	"github.com/antha-lang/antha/ast"
	"github.com/antha-lang/antha/graph"
	"github.com/antha-lang/antha/target"
)

// instGraph is a graph to that models dependencies between target instructions
// but uses edges from separately initialized dependsOn field. Useful for cases
// where we can't depend on target.Graph because we are using this to build the
// initial DependsOn relation.
type instGraph struct {
	insts     []ast.Inst
	added     map[ast.Inst]bool
	dependsOn map[ast.Inst][]ast.Inst
	entry     map[graph.Node]ast.Inst
	exit      map[graph.Node]ast.Inst
}

func newInstGraph() *instGraph {
	return &instGraph{
		added:     make(map[ast.Inst]bool),
		entry:     make(map[graph.Node]ast.Inst),
		exit:      make(map[graph.Node]ast.Inst),
		dependsOn: make(map[ast.Inst][]ast.Inst),
	}
}

func (a *instGraph) NumNodes() int {
	return len(a.insts)
}

func (a *instGraph) Node(i int) graph.Node {
	return a.insts[i]
}

func (a *instGraph) NumOuts(n graph.Node) int {
	return len(a.dependsOn[n.(ast.Inst)])
}

func (a *instGraph) Out(n graph.Node, i int) graph.Node {
	return a.dependsOn[n.(ast.Inst)][i]
}

// addInsts adds instructions to the graph
func (a *instGraph) addInsts(insts ast.Insts) {
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

func (a *instGraph) addInitializers(insts ast.Insts) {
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

func (a *instGraph) addFinalizers(insts ast.Insts) {
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
func (a *instGraph) addRootedInsts(root graph.Node, insts ast.Insts) {
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

	newInsts := append(ast.Insts{entry, exit}, insts...)
	a.addInsts(newInsts)
}
