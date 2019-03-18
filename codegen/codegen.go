// Package codegen compiles generic instructions to target-specific ones.
// Target, in this case, is some combination of devices (e.g., two
// ExtendedLiquidHandlers and human plate mover).
package codegen

import (
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/antha-lang/antha/graph"
	"github.com/antha-lang/antha/laboratory/effects"
	"github.com/antha-lang/antha/laboratory/effects/id"
	"github.com/antha-lang/antha/target"
)

// Intermediate representation.
type ir struct {
	Root         effects.Node
	Graph        *effects.Graph                      // Graph of effects.Nodes
	Commands     graph.Graph                         // DAG of effects.Commands (and potentially BundleExpr root)
	DeviceDeps   graph.QGraph                        // Dependencies of druns
	reachingUses map[effects.Node][]*effects.UseComp // Reaching comps
	assignment   map[effects.Node]*drun              // From Commands/Root to device runs
	output       map[*drun][]effects.Inst            // Output of device-specific planners
	initializers []effects.Inst                      // Intializers
	finalizers   []effects.Inst                      // Finalizers in reverse order
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
			c, ok := x.(*effects.Command)
			if !ok {
				return ""
			}
			return fmt.Sprintf("%T", c.Inst)
		},
		func(x interface{}) string {
			n, ok := x.(effects.Node)
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
			n, ok := x.(effects.Node)
			if !ok {
				return ""
			}

			u, ok := n.(*effects.UseComp)
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
	Device effects.Device
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
type partitionByHuman []effects.Device

func (a partitionByHuman) Len() int {
	return len(a)
}

func (a partitionByHuman) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a partitionByHuman) Less(i, j int) bool {
	req := effects.Request{
		Selector: []effects.NameValue{
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
	bundleReqs := func(n *effects.Bundle) (reqs []effects.Request) {
		for i, inum := 0, a.Commands.NumOuts(n); i < inum; i++ {
			kid := a.Commands.Out(n, i)
			if c, ok := kid.(*effects.Command); ok {
				reqs = append(reqs, c.Request)
			}
		}
		return
	}

	colors := make(map[effects.Node][]effects.Device)
	for i, inum := 0, a.Commands.NumNodes(); i < inum; i++ {
		n := a.Commands.Node(i).(effects.Node)
		var reqs []effects.Request
		if c, ok := n.(*effects.Command); ok {
			reqs = append(reqs, c.Request)
		} else if b, ok := n.(*effects.Bundle); ok {
			// Try to find device that can do everything
			reqs = bundleReqs(b)
		} else {
			return fmt.Errorf("unknown node %T", n)
		}
		devices := t.CanCompile(reqs...)

		if len(devices) == 0 {
			return fmt.Errorf("no device can handle constraints %v", effects.Meet(reqs...))
		}
		sort.Stable(partitionByHuman(devices))
		colors[n] = devices
	}

	var devices []effects.Device
	d2c := make(map[effects.Device]int)
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
			for _, d := range colors[n.(effects.Node)] {
				r = append(r, d2c[d])
			}
			return
		},
		EdgeWeight: func(a, b int) int64 {
			return 0
		},
	})
	if err != nil {
		return err
	}

	ret := make(map[effects.Node]effects.Device)
	for n, idx := range r.Parts {
		ret[n.(effects.Node)] = devices[idx]
	}

	a.coalesceDevices(ret)

	return nil
}

// Coalesce adjacent devices into the same run of a device
func (a *ir) coalesceDevices(device map[effects.Node]effects.Device) {
	run := make(map[effects.Node]*drun)

	kidRun := func(n effects.Node) *drun {
		m := make(map[*drun]bool)
		for i, inum := 0, a.Commands.NumOuts(n); i < inum; i++ {
			kid := a.Commands.Out(n, i).(effects.Node)
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
		newRuns := make(map[effects.Device]*drun)
		for _, n := range dag.Roots {
			n := n.(effects.Node)

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
func (a *ir) tryPlan(labEffects *effects.LaboratoryEffects, dir string) error {
	dg := graph.MakeQuotient(graph.MakeQuotientOpt{
		Graph: a.Commands,
		Colorer: func(n graph.Node) interface{} {
			return a.assignment[n.(effects.Node)]
		},
	})

	// TODO: When initial assignment is not feasible for a device (e.g.,
	// capacity constraints), split up run until feasible or give up.

	// TODO: When splitting a mix sequence, adjust LHInstructions to place
	// output samples on the same plate

	cmds := make(map[*drun][]effects.Node)
	for n, d := range a.assignment {
		c, ok := n.(*effects.Command)
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
		n := dg.Orig(n, 0).(effects.Node)
		run := a.assignment[n]
		runs = append(runs, run)
	}

	a.output = make(map[*drun][]effects.Inst)
	for _, d := range runs {
		insts, err := d.Device.Compile(labEffects, dir, cmds[d])
		if err != nil {
			return err
		}

		for _, n := range cmds[d] {
			c := n.(*effects.Command)
			c.Output = insts
		}

		a.output[d] = insts
	}

	return a.addImplicitInsts(runs)
}

func (a *ir) sortDevices(labEffects *effects.LaboratoryEffects, t *target.Target) error {
	a.DeviceDeps = graph.MakeQuotient(graph.MakeQuotientOpt{
		Graph: a.Commands,
		Colorer: func(n graph.Node) interface{} {
			return a.assignment[n.(effects.Node)]
		},
	})

	_, err := graph.TopoSort(graph.TopoSortOpt{
		Graph: a.DeviceDeps,
	})
	return err
}

func reverseInsts(insts []effects.Inst) (ret []effects.Inst) {
	for idx := len(insts) - 1; idx >= 0; idx-- {
		ret = append(ret, insts[idx])
	}
	return
}

// Lower plan to instructions
func (a *ir) genInsts(idGen *id.IDGenerator) (effects.Insts, error) {
	ig := newInstGraph()

	// Insert instructions
	for i, inum := 0, a.DeviceDeps.NumNodes(); i < inum; i++ {
		n := a.DeviceDeps.Node(i)
		someNode := a.DeviceDeps.Orig(n, 0).(effects.Node)
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

	var insts effects.Insts
	for _, n := range order {
		in := n.(effects.Inst)
		in.SetDependsOn() // reset to empty first
		for j, jnum := 0, sg.NumOuts(n); j < jnum; j++ {
			in.AppendDependsOn(sg.Out(n, j).(effects.Inst))
		}
		in.SetId(idGen)
		insts = append(insts, in)
	}

	return insts, nil
}

// Compile an expression program into a sequence of instructions for a target
// configuration. This supports incremental compilation, so roots may refer to
// nodes that have already been compiled, in which case, the result may refer
// to previously generated instructions.
func Compile(labEffects *effects.LaboratoryEffects, dir string, t *target.Target, roots []effects.Node) (effects.Insts, error) {
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
	if err := ir.tryPlan(labEffects, dir); err != nil {
		return nil, fmt.Errorf("error planning: %s", err)
	}

	if err := ir.sortDevices(labEffects, t); err != nil {
		return nil, fmt.Errorf("error sorting devices: %s", err)
	}

	insts, err := ir.genInsts(labEffects.IDGenerator)
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
