// Package codegen compiles generic instructions to target-specific ones.
// Target, in this case, is some combination of devices (e.g., two
// ExtendedLiquidHandlers and human plate mover).
package codegen

import (
	"context"
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/antha-lang/antha/ast"
	"github.com/antha-lang/antha/graph"
	"github.com/antha-lang/antha/target"
	"github.com/antha-lang/antha/target/human"
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
	initializers []target.Inst               // Intializers
	finalizers   []target.Inst               // Finalizers in reverse order
}

// Result is result of compiling a set of ast.Nodes
type Result struct {
	Insts []target.Inst
}

// Print out IR for debugging
func (a *ir) Print(g graph.Graph, out io.Writer) error {
	shortID := func(x string) string {
		for _, p := range strings.Split(x, "-") {
			return p
		}
		return x
	}

	labelers := []func(interface{}) string{
		func(x interface{}) string {
			c, ok := x.(*ast.Command)
			if !ok {
				return ""
			}
			return fmt.Sprintf("%T", c.Inst)
		},
		func(x interface{}) string {
			c, ok := x.(*ast.Command)
			if !ok {
				return ""
			}
			h, ok := c.Inst.(*ast.HandleInst)
			if !ok {
				return ""
			}
			return h.Group
		},
		func(x interface{}) string {
			n, ok := x.(ast.Node)
			if !ok {
				return ""
			}
			drun := a.assignment[n]
			if drun != nil {
				return fmt.Sprintf("Run %p Device %v %s", drun, drun.Device, drun.Device)
			}
			return ""
		},
		func(x interface{}) string {
			n, ok := x.(ast.Node)
			if !ok {
				return ""
			}

			u, ok := n.(*ast.UseComp)
			if !ok {
				return ""
			}
			return fmt.Sprintf("%s (%s)", u.Value.CName, shortID(u.Value.ID))
		},
		func(x interface{}) string {
			n, ok := x.(*target.Manual)
			if !ok {
				return ""
			}
			return n.Label
		},
	}

	label := func(x interface{}) string {
		var items []string
		for _, l := range labelers {
			s := l(x)
			if len(s) != 0 {
				items = append(items, s)
			}
		}
		return strings.Join(items, "\n")
	}

	s := graph.Print(graph.PrintOpt{
		Graph: g,
		NodeLabelers: []graph.Labeler{
			label,
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
	ret := &graph.TreePartition{
		Parts: make(map[graph.Node]int),
	}
	// Simple first-fit algorithm but handles arbitrary graph structures
	for i, inum := 0, opt.Tree.NumNodes(); i < inum; i++ {
		n := opt.Tree.Node(i)
		ret.Parts[n] = opt.Colors(n)[0]
	}
	return ret, nil
}

// Partition a slice into non-human devices followed by human ones
type partitionByHuman []target.Device

func (a partitionByHuman) Len() int {
	return len(a)
}

func (a partitionByHuman) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a partitionByHuman) Less(i, j int) bool {
	req := ast.Request{
		Selector: []ast.NameValue{
			target.DriverSelectorV1Human,
		},
	}

	human1 := a[i].CanCompile(req)
	human2 := a[j].CanCompile(req)
	switch {
	case human1 && human2:
		return false // Equal
	case human1 && !human2:
		return false
	case !human1 && human2:
		return true
	default:
		return false // Equal
	}
}

// Assign runs of a device to each ApplyExpr. Construct initial plan by
// by maximally coalescing ApplyExprs with the same device into the same
// device run.
func (a *ir) assignDevices(t *target.Target) error {
	// A bundle's requests is the sum of its children
	bundleReqs := func(n *ast.Bundle) (reqs []ast.Request) {
		for i, inum := 0, a.Commands.NumOuts(n); i < inum; i++ {
			kid := a.Commands.Out(n, i)
			if c, ok := kid.(*ast.Command); ok {
				reqs = append(reqs, c.Requests...)
			}
		}
		return
	}

	colors := make(map[ast.Node][]target.Device)
	for i, inum := 0, a.Commands.NumNodes(); i < inum; i++ {
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
				return fmt.Errorf("no device can handle constraints %v", ast.Meet(reqs...))
			}
		}
		sort.Stable(partitionByHuman(devices))
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
		EdgeWeight: func(a, b int) int64 {
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
		for i, inum := 0, a.Commands.NumOuts(n); i < inum; i++ {
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
func (a *ir) tryPlan(ctx context.Context) error {
	dg := graph.MakeQuotient(graph.MakeQuotientOpt{
		Graph: a.Commands,
		Colorer: func(n graph.Node) interface{} {
			return a.assignment[n.(ast.Node)]
		},
	})

	// TODO: When initial assignment is not feasible for a device (e.g.,
	// capacity constraints), split up run until feasible or give up.

	// TODO: When splitting a mix sequence, adjust LHInstructions to place
	// output samples on the same plate

	cmds := make(map[*drun][]ast.Node)
	for n, d := range a.assignment {
		c, ok := n.(*ast.Command)
		if !ok {
			continue
		}
		cmds[d] = append(cmds[d], c)
	}

	// Process runs in dependency order
	order, err := graph.TopoSort(graph.TopoSortOpt{
		Graph: dg,
	})
	if err != nil {
		return fmt.Errorf("invalid assignment: %s", err)
	}
	var runs []*drun
	for _, n := range order {
		n := dg.Orig(n, 0).(ast.Node)
		run := a.assignment[n]
		runs = append(runs, run)
	}

	a.output = make(map[*drun][]target.Inst)
	for _, d := range runs {
		insts, err := d.Device.Compile(ctx, cmds[d])
		if err != nil {
			return err
		}

		result := &Result{
			Insts: insts,
		}
		for _, n := range cmds[d] {
			c := n.(*ast.Command)
			c.Output = result
		}

		a.output[d] = insts
	}

	return a.addImplicitInsts(runs)
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
		if head != nil && tail != nil && head != tail {
			tail.AppendDependsOn(head)
		}
	} else {
		oldH := insts[0]
		oldT := insts[len(insts)-1]
		if head != nil {
			oldH.AppendDependsOn(head)
		}
		if tail != nil {
			tail.AppendDependsOn(oldT)
		}
	}
}

// Add implied moves between devices
func (a *ir) addMoves(ctx context.Context, t *target.Target) error {
	a.DeviceDeps = graph.MakeQuotient(graph.MakeQuotientOpt{
		Graph: a.Commands,
		Colorer: func(n graph.Node) interface{} {
			return a.assignment[n.(ast.Node)]
		},
	})

	_, err := graph.TopoSort(graph.TopoSortOpt{
		Graph: a.DeviceDeps,
	})
	return err
}

func reverseInsts(insts []target.Inst) (ret []target.Inst) {
	for idx := len(insts) - 1; idx >= 0; idx-- {
		ret = append(ret, insts[idx])
	}
	return
}

// Lower plan to instructions
func (a *ir) genInsts() ([]target.Inst, error) {
	ig := newInstGraph()

	// Insert instructions
	for i, inum := 0, a.DeviceDeps.NumNodes(); i < inum; i++ {
		n := a.DeviceDeps.Node(i)
		someNode := a.DeviceDeps.Orig(n, 0).(ast.Node)
		run := a.assignment[someNode]
		insts := a.output[run]
		ig.addRootedInsts(n, insts)
	}

	ig.addInitializers(a.initializers)
	ig.addFinalizers(reverseInsts(a.finalizers))

	// Add tree edges
	for i, inum := 0, a.DeviceDeps.NumNodes(); i < inum; i++ {
		n := a.DeviceDeps.Node(i)
		nentry := ig.entry[n]
		for j, jnum := 0, a.DeviceDeps.NumOuts(n); j < jnum; j++ {
			dst := a.DeviceDeps.Out(n, j)
			dexit := ig.exit[dst]
			ig.dependsOn[nentry] = append(ig.dependsOn[nentry], dexit)
		}
	}

	// Remove synthetic nodes and redundant edges
	sg, err := simplifyWithDeps(ig, func(n graph.Node) bool {
		_, isWait := n.(*target.Wait)
		return !isWait
	})
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
		in := n.(target.Inst)
		in.SetDependsOn() // reset to empty first
		for j, jnum := 0, sg.NumOuts(n); j < jnum; j++ {
			in.AppendDependsOn(sg.Out(n, j).(target.Inst))
		}
		insts = append(insts, in)
	}

	return insts, nil
}

// Compile an expression program into a sequence of instructions for a target
// configuration. This supports incremental compilation, so roots may refer to
// nodes that have already been compiled, in which case, the result may refer
// to previously generated instructions.
func Compile(ctx context.Context, t *target.Target, roots []ast.Node) ([]target.Inst, error) {
	if len(roots) == 0 {
		return nil, nil
	}

	root, err := makeRoot(roots)
	if err != nil {
		return nil, fmt.Errorf("invalid program: %s", err)
	}
	ir, err := build(root)
	if err != nil {
		return nil, fmt.Errorf("invalid program: %s", err)
	}
	if err := ir.assignDevices(t); err != nil {
		return nil, fmt.Errorf("error assigning devices with target configuration %s: %s", t, err)
	}
	if err := ir.tryPlan(ctx); err != nil {
		return nil, fmt.Errorf("error planning: %s", err)
	}

	if err := ir.addMoves(ctx, t); err != nil {
		return nil, fmt.Errorf("error adding moves: %s", err)
	}

	insts, err := ir.genInsts()
	if err != nil {
		return nil, fmt.Errorf("error generating instructions: %s", err)
	}

	// TODO: discard programs that create multiple setups until we get their
	// semantics correct; also true of incubating components under multiple
	// conditions
	var setupMixes int
	var setupIncubators int
	for _, inst := range insts {
		switch inst.(type) {
		case *target.SetupMixer:
			setupMixes++
		case *target.SetupIncubator:
			setupIncubators++
		}
	}
	if setupMixes > 1 || setupIncubators > 1 {
		return nil, fmt.Errorf("multiple incubates or multiple mixes not supported")
	}

	return insts, nil
}
