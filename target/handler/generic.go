package handler

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"

	"github.com/antha-lang/antha/graph"
	"github.com/antha-lang/antha/laboratory/effects"
	"github.com/antha-lang/antha/target"
)

var (
	errCannotMergeUnequalCalls = errors.New("cannot merge unequal calls")
)

// A GenericHandler is a configurable version of a Handler suitable for mixins
type GenericHandler struct {
	Labels             []effects.NameValue
	GenFunc            func(cmd interface{}) ([]effects.Inst, error)
	FilterFieldsForKey func(interface{}) (interface{}, error)
}

// CanCompile implements a Device
func (a *GenericHandler) CanCompile(req effects.Request) bool {
	can := effects.Request{
		Selector: a.Labels,
	}

	return can.Contains(req)
}

func (a GenericHandler) serialize(obj interface{}) (string, error) {
	type hasGetID interface {
		GetID() string
	}

	if g, ok := obj.(hasGetID); ok {
		return g.GetID(), nil
	}

	var out bytes.Buffer
	var err error
	enc := gob.NewEncoder(&out)
	if a.FilterFieldsForKey != nil {
		obj, err = a.FilterFieldsForKey(obj)
		if err != nil {
			return "", err
		}
	}

	if err := enc.Encode(obj); err != nil {
		return "", err
	}

	return out.String(), nil
}

func (a GenericHandler) merge(nodes []effects.Node) (*effects.Command, error) {
	if len(nodes) == 0 {
		return nil, nil
	}

	cmd, ok := nodes[0].(*effects.Command)
	if !ok {
		return nil, fmt.Errorf("expecting %T but found %T instead", cmd, nodes[0])
	}

	retStr, err := a.serialize(cmd.Inst)
	if err != nil {
		return nil, err
	}

	for _, n := range nodes[1:] {
		cmd, ok := n.(*effects.Command)
		if !ok {
			return nil, fmt.Errorf("expecting %T but found %T instead", cmd, nodes[0])
		}

		cmdStr, err := a.serialize(cmd.Inst)
		if err != nil {
			return nil, err
		}

		if retStr != cmdStr {
			return nil, errCannotMergeUnequalCalls
		}
	}

	return cmd, nil
}

// Compile implements a Device
func (a *GenericHandler) Compile(labEffects *effects.LaboratoryEffects, nodes []effects.Node) ([]effects.Inst, error) {
	g := effects.Deps(nodes)

	entry := &target.Wait{}
	exit := &target.Wait{}
	var insts []effects.Inst
	inst := make(map[effects.Node][]effects.Inst)

	insts = append(insts, entry)

	// Maximally coalesce repeated commands according to when they are first
	// available to be executed (graph.Reverse)
	dag := graph.Schedule(graph.Reverse(g))
	for len(dag.Roots) > 0 {
		var next []graph.Node
		// Gather
		same := make(map[interface{}][]effects.Node)
		for _, r := range dag.Roots {
			cmd, ok := r.(*effects.Command)
			if !ok {
				return nil, fmt.Errorf("expecting %T but found %T instead", cmd, r)
			}

			key, err := a.serialize(cmd.Inst)
			if err != nil {
				return nil, err
			}

			same[key] = append(same[key], r.(effects.Node))
			next = append(next, dag.Visit(r)...)
		}
		// Apply
		for _, nodes := range same {
			cmd, err := a.merge(nodes)
			if err != nil {
				return nil, err
			}
			if cmd == nil {
				continue
			}

			ins, err := a.GenFunc(cmd.Inst)
			if err != nil {
				return nil, err
			}

			insts = append(insts, ins...)

			for _, n := range nodes {
				inst[n] = ins
			}
		}

		dag.Roots = next
	}

	insts = append(insts, exit)

	for i, inum := 0, g.NumNodes(); i < inum; i++ {
		n := g.Node(i).(effects.Node)
		ins := inst[n]
		if len(ins) == 0 {
			continue
		}

		for j, jnum := 0, g.NumOuts(n); j < jnum; j++ {
			kid := g.Out(n, j).(effects.Node)
			kidIns := inst[kid]
			if len(kidIns) == 0 {
				continue
			}

			ins[0].AppendDependsOn(kidIns[len(kidIns)-1])
		}
		ins[0].AppendDependsOn(entry)
		exit.AppendDependsOn(ins[len(ins)-1])
	}

	return insts, nil
}
