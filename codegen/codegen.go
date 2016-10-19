// Package codegen compiles generic instructions to target-specific ones.
// Target, in this case, is some combination of devices (e.g., two
// ExtendedLiquidHandlers and human plate mover).
package codegen

import (
	"fmt"
	"io"

	"github.com/antha-lang/antha/ast"
	"github.com/antha-lang/antha/graph"
	"github.com/antha-lang/antha/target"
	"github.com/antha-lang/antha/target/human"
)

const (
	useTreePartition = false
)

// Intermediate representation.
type ir struct {
	Root         ast.Node
	Graph        *ast.Graph                  // Graph of ast.Nodes
	Commands     graph.Graph                 // DAG of ast.Commands (and potentially BundleExpr root)
	DeviceDeps   graph.QGraph                // Dependencies of druns
	reachingUses map[ast.Node][]*ast.UseComp // Reaching comps
	assignment   map[ast.Node]*drun          // From Commands/Root to device runs
	output       map[*drun][]target.Inst     // Output of device-specific planners
}

// Print out IR for debugging
func (a *ir) Print(g graph.Graph, out io.Writer) error {
	s := graph.Print(graph.PrintOpt{
		Graph: g,
		NodeLabelers: []graph.Labeler{
			func(x interface{}) string {
				if c, ok := x.(*ast.Command); ok {
					return fmt.Sprintf("%T", c.Inst)
				} else {
					return ""
				}
			},
			func(x interface{}) string {
				n := x.(ast.Node)
				drun := a.assignment[n]
				if drun != nil {
					return fmt.Sprintf("Run %p Device %p %s", drun, drun.Device, drun.Device)
				}
				return "NoRun"
			},
			func(x interface{}) string {
				n := x.(ast.Node)
				if u, ok := n.(*ast.UseComp); ok {
					return u.Value.CName
				}
				return ""
			},
		},
	})
	_, err := fmt.Fprint(out, s, "\n")
	return err
}

// Run of a device.
type drun struct {
	Device target.Device
}

func (a *ir) partition(opt graph.PartitionTreeOpt) (*graph.TreePartition, error) {
	if useTreePartition {
		if err := graph.IsTree(opt.Tree, opt.Root); err != nil {
			return nil, err
		}
		return graph.PartitionTree(opt)
	}

	ret := &graph.TreePartition{
		Parts: make(map[graph.Node]int),
	}
	// Simple first-fit algorithm but handles arbitrary graph structures
	for i, inum := 0, opt.Tree.NumNodes(); i < inum; i += 1 {
		n := opt.Tree.Node(i)
		ret.Parts[n] = opt.Colors(n)[0]
	}
	return ret, nil
}

// Assign runs of a device to each ApplyExpr. Construct initial plan by
// by maximally coalescing ApplyExprs with the same device into the same
// device run.
func (a *ir) assignDevices(t *target.Target) error {
	// A bundle's requests is the sum of its children
	bundleReqs := func(n *ast.Bundle) (reqs []ast.Request) {
		for i, inum := 0, a.Commands.NumOuts(n); i < inum; i += 1 {
			kid := a.Commands.Out(n, i)
			if c, ok := kid.(*ast.Command); ok {
				reqs = append(reqs, c.Requests...)
			}
		}
		return
	}

	colors := make(map[ast.Node][]target.Device)
	for i, inum := 0, a.Commands.NumNodes(); i < inum; i += 1 {
		n := a.Commands.Node(i).(ast.Node)
		var reqs []ast.Request
		isBundle := false
		if c, ok := n.(*ast.Command); ok {
			reqs = c.Requests
		} else if b, ok := n.(*ast.Bundle); ok {
			// Try to find device that can do everything
			reqs = bundleReqs(b)
			isBundle = true
		} else {
			return fmt.Errorf("unknown node %T", n)
		}
		devices := t.CanCompile(reqs...)

		if len(devices) == 0 {
			if isBundle {
				devices = append(devices, human.New(human.Opt{}))
			} else {
				return fmt.Errorf("no device can handle constraints %s", ast.Meet(reqs...))
			}
		}
		colors[n] = devices
	}

	var devices []target.Device
	d2c := make(map[target.Device]int)
	for _, ds := range colors {
		for _, d := range ds {
			if _, seen := d2c[d]; !seen {
				d2c[d] = len(devices)
				devices = append(devices, d)
			}
		}
	}

	r, err := a.partition(graph.PartitionTreeOpt{
		Tree: a.Commands,
		Root: a.Root,
		Colors: func(n graph.Node) (r []int) {
			for _, d := range colors[n.(ast.Node)] {
				r = append(r, d2c[d])
			}
			return
		},
		EdgeWeight: func(a, b int) int {
			return devices[a].MoveCost(devices[b])
		},
	})
	if err != nil {
		return err
	}

	ret := make(map[ast.Node]target.Device)
	for n, idx := range r.Parts {
		ret[n.(ast.Node)] = devices[idx]
	}

	a.coalesceDevices(ret)

	return nil
}

// Coalesce adjacent devices into the same run of a device
func (a *ir) coalesceDevices(device map[ast.Node]target.Device) {
	run := make(map[ast.Node]*drun)

	kidRun := func(n ast.Node) *drun {
		m := make(map[*drun]bool)
		for i, inum := 0, a.Commands.NumOuts(n); i < inum; i += 1 {
			kid := a.Commands.Out(n, i).(ast.Node)
			m[run[kid]] = true
			if device[kid] != device[n] {
				return nil
			}
		}
		if len(m) != 1 {
			return nil
		}
		for k := range m {
			return k
		}
		return nil
	}

	dag := graph.Schedule(graph.Reverse(a.Commands))

	for len(dag.Roots) > 0 {
		var next []graph.Node
		newRuns := make(map[target.Device]*drun)
		for _, n := range dag.Roots {
			n := n.(ast.Node)

			myRun := kidRun(n)
			if myRun == nil {
				d := device[n]
				if r, seen := newRuns[d]; seen {
					myRun = r
				} else {
					myRun = &drun{Device: d}
					newRuns[d] = myRun
				}
			}
			run[n] = myRun
			next = append(next, dag.Visit(n)...)
		}

		dag.Roots = next
	}

	a.assignment = run
}

// Run plan through device-specific planners. Adjust assignment based on
// planner capabilities and return output.
func (a *ir) tryPlan() error {
	dg := graph.MakeQuotient(graph.MakeQuotientOpt{
		Graph: a.Commands,
		Colorer: func(n graph.Node) interface{} {
			return a.assignment[n.(ast.Node)]
		},
	})
	if err := graph.IsDag(dg); err != nil {
		return fmt.Errorf("invalid assignment: %s", err)
	}

	// TODO: When initial assignment is not feasible for a device (e.g.,
	// capacity constraints), split up run until feasible or give up.

	// TODO: When splitting a mix sequence, adjust LHInstructions to place
	// output samples on the same plate

	cmds := make(map[*drun][]ast.Node)
	for n, d := range a.assignment {
		if c, ok := n.(*ast.Command); !ok {
			continue
		} else {
			cmds[d] = append(cmds[d], c)
		}
	}

	output := make(map[*drun][]target.Inst)
	for d, cs := range cmds {
		if insts, err := d.Device.Compile(cs); err != nil {
			return err
		} else {
			output[d] = insts
		}
	}

	a.output = output

	return nil
}

// Find best device to move a component between two devices
func findBestMoveDevice(t *target.Target, from, to ast.Node, fromD, toD *drun) target.Device {
	// TODO: add movement constraints
	var req ast.Request
	var minD target.Device
	minC := -1

	for _, d := range t.CanCompile(req) {
		c := toD.Device.MoveCost(d) + d.MoveCost(fromD.Device)
		if minC == -1 || c < minC {
			minC = c
			minD = d
		}
	}
	return minD
}

// NB(ddn): Could blindly add edges from insts to head, but would like
// Compile() to be able to introduce instructions that just depend on the start
// or end (or neither) of a device run.
//
// From:
//   head: h
//   tail: t
//   insts: [a <- ... <- b]
// To:
//   h <- a <-... <- b <- t
func splice(head, tail target.Inst, insts []target.Inst) {
	if len(insts) == 0 {
		if head != nil && tail != nil {
			tail.SetDependsOn(append(tail.DependsOn(), head))
		}
		return
	}
	oldH := insts[0]
	oldT := insts[len(insts)-1]
	if head != nil {
		oldH.SetDependsOn(append(oldH.DependsOn(), head))
	}
	if tail != nil {
		tail.SetDependsOn(append(tail.DependsOn(), oldT))
	}
}

// Create move of dependencies if necessary
func (a *ir) addMove(t *target.Target, dnode graph.Node, run *drun) error {
	rewrite := func(n ast.Node, cs []*ast.UseComp, move *ast.Move) {
		m := make(map[ast.Node]bool)
		for _, c := range cs {
			m[c] = true
		}
		for i, inum := 0, a.Graph.NumOuts(n); i < inum; i += 1 {
			out := a.Graph.Out(n, i).(ast.Node)
			if m[out] {
				a.Graph.SetOut(n, i, move)
			}
		}
	}

	newRuns := make(map[target.Device]*drun)
	getRun := func(d target.Device) *drun {
		r, ok := newRuns[d]
		if !ok {
			r = &drun{d}
			newRuns[d] = r
		}
		return r
	}

	moves := make(map[target.Device][]ast.Node)
	for i, inum := 0, a.DeviceDeps.NumOrigs(dnode); i < inum; i += 1 {
		n := a.DeviceDeps.Orig(dnode, i).(ast.Node)
		for j, jnum := 0, a.Commands.NumOuts(n); j < jnum; j += 1 {
			out := a.Commands.Out(n, j).(ast.Node)
			if run == a.assignment[out] {
				continue
			}
			// Command has input dependence on a previous run
			cs := a.reachingUses[out.(*ast.Command)]
			if len(cs) == 0 {
				// Nothing to move
				continue
			} else if dev := findBestMoveDevice(t, out, n, a.assignment[out], run); dev == nil {
				return fmt.Errorf("cannot find any device to move inputs")
			} else {
				// Add move
				m := &ast.Move{From: cs, ToLoc: fmt.Sprintf("%p", dev)}
				moves[dev] = append(moves[dev], m)
				a.assignment[m] = getRun(dev)
				rewrite(n, cs, m)
			}
		}
	}

	if len(moves) == 0 {
		return nil
	}

	var insts []target.Inst
	head := &target.Wait{}
	tail := &target.Wait{}
	insts = append(insts, head, tail)

	splice(head, tail, nil)
	splice(tail, nil, a.output[run])

	for dev, ms := range moves {
		if ins, err := dev.Compile(ms); err != nil {
			return nil
		} else {
			splice(head, tail, ins)
			insts = append(insts, ins...)
		}
	}

	a.output[run] = append(insts, a.output[run]...)
	return nil
}

// Add implied moves between devices
func (a *ir) addMoves(t *target.Target) error {
	a.DeviceDeps = graph.MakeQuotient(graph.MakeQuotientOpt{
		Graph: a.Commands,
		Colorer: func(n graph.Node) interface{} {
			return a.assignment[n.(ast.Node)]
		},
	})

	order, err := graph.TopoSort(graph.TopoSortOpt{
		Graph: a.DeviceDeps,
	})
	if err != nil {
		return err
	}

	for _, n := range order {
		if a.DeviceDeps.NumOrigs(n) == 0 {
			return fmt.Errorf("no instructions for node %q", n)
		}
		someNode := a.DeviceDeps.Orig(n, 0).(ast.Node)
		run := a.assignment[someNode]
		if err := a.addMove(t, n, run); err != nil {
			return err
		}
	}

	return nil
}

// Lower plan to instructions
func (a *ir) genInsts() ([]target.Inst, error) {
	ig := &instGraph{
		entry:     make(map[graph.Node]target.Inst),
		exit:      make(map[graph.Node]target.Inst),
		dependsOn: make(map[target.Inst][]target.Inst),
	}

	// Insert instructions
	for i, inum := 0, a.DeviceDeps.NumNodes(); i < inum; i += 1 {
		n := a.DeviceDeps.Node(i)
		someNode := a.DeviceDeps.Orig(n, 0).(ast.Node)
		run := a.assignment[someNode]
		insts := a.output[run]
		ig.addInsts(n, insts)
	}

	// Add tree edges
	for i, inum := 0, a.DeviceDeps.NumNodes(); i < inum; i += 1 {
		n := a.DeviceDeps.Node(i)
		nentry := ig.entry[n]
		for j, jnum := 0, a.DeviceDeps.NumOuts(n); j < jnum; j += 1 {
			dst := a.DeviceDeps.Out(n, j)
			dexit := ig.exit[dst]
			ig.dependsOn[nentry] = append(ig.dependsOn[nentry], dexit)
		}
	}

	// Remove synthetic nodes and redundant edges
	sg, err := graph.TransitiveReduction(graph.Eliminate(graph.EliminateOpt{
		Graph: graph.Simplify(graph.SimplifyOpt{
			Graph:            ig,
			RemoveSelfLoops:  true,
			RemoveMultiEdges: true,
		}),
		In: func(n graph.Node) bool {
			if _, isWait := n.(*target.Wait); isWait {
				return false
			}
			return true
		},
	}))
	if err != nil {
		return nil, err
	}

	// Cleanup dependencies
	order, err := graph.TopoSort(graph.TopoSortOpt{
		Graph: sg,
	})
	if err != nil {
		return nil, err
	}

	var insts []target.Inst
	for _, n := range order {
		var depends []target.Inst
		for j, jnum := 0, sg.NumOuts(n); j < jnum; j += 1 {
			depends = append(depends, sg.Out(n, j).(target.Inst))
		}

		in := n.(target.Inst)
		in.SetDependsOn(depends)

		insts = append(insts, in)
	}

	return insts, nil
}

// Mark nodes as compiled
func (a *ir) setOutputs() error {
	for _, n := range a.Graph.Nodes {
		if c, ok := n.(*ast.Command); !ok {
			continue
		} else if c.Output != nil {
			continue
		} else if run := a.assignment[c]; run != nil {
			c.Output = run
		}
	}
	return nil
}

// Dependencies between target instructions. Can't use target.Graph because we
// are using this to build the initial DependsOn relation.
type instGraph struct {
	insts     []target.Inst
	dependsOn map[target.Inst][]target.Inst
	entry     map[graph.Node]target.Inst
	exit      map[graph.Node]target.Inst
}

func (a *instGraph) NumNodes() int {
	return len(a.insts)
}

func (a *instGraph) Node(i int) graph.Node {
	return a.insts[i]
}

func (a *instGraph) NumOuts(n graph.Node) int {
	return len(a.dependsOn[n.(target.Inst)])
}

func (a *instGraph) Out(n graph.Node, i int) graph.Node {
	return a.dependsOn[n.(target.Inst)][i]
}

func (a *instGraph) addInsts(root graph.Node, insts []target.Inst) {
	exit := &target.Wait{}
	entry := &target.Wait{}

	a.entry[root] = entry
	a.exit[root] = exit

	// Add dependencies
	a.dependsOn[exit] = append(a.dependsOn[exit], entry)
	for idx, in := range insts {
		if idx == 0 {
			a.dependsOn[in] = append(a.dependsOn[in], entry)
		}
		if idx == len(insts)-1 {
			a.dependsOn[exit] = append(a.dependsOn[exit], in)
		}
		for _, v := range in.DependsOn() {
			a.dependsOn[in] = append(a.dependsOn[in], v)
		}
	}

	// Add nodes
	toAdd := make(map[target.Inst]bool)
	toAdd[entry] = true
	toAdd[exit] = true
	for _, in := range insts {
		toAdd[in] = true
		for _, v := range in.DependsOn() {
			toAdd[v] = true
		}
	}

	for in := range toAdd {
		a.insts = append(a.insts, in)
	}
}

// Compile an expression program into a sequence of instructions for a target
// configuration. This supports incremental compilation, so roots may refer to
// nodes that have already been compiled, in which case, the result may refer
// to previously generated instructions.
func Compile(t *target.Target, roots []ast.Node) ([]target.Inst, error) {
	if len(roots) == 0 {
		return nil, nil
	}

	if root, err := makeRoot(roots); err != nil {
		return nil, fmt.Errorf("invalid program: %s", err)
	} else if ir, err := build(root); err != nil {
		return nil, fmt.Errorf("invalid program: %s", err)
	} else if err := ir.assignDevices(t); err != nil {
		return nil, fmt.Errorf("error assigning devices with target configuration %s: %s", t, err)
	} else if err := ir.tryPlan(); err != nil {
		return nil, fmt.Errorf("error planning: %s", err)
	} else if err := ir.addMoves(t); err != nil {
		return nil, fmt.Errorf("error adding moves: %s", err)
	} else if insts, err := ir.genInsts(); err != nil {
		return nil, fmt.Errorf("error generating instructions: %s", err)
	} else if err := ir.setOutputs(); err != nil {
		return nil, fmt.Errorf("error setting outputs: %s", err)
	} else {
		return insts, nil
	}
}
