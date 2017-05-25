package handler

import (
	"bytes"
	"errors"
	"fmt"
	"io"

	"github.com/antha-lang/antha/ast"
	"github.com/antha-lang/antha/driver"
	"github.com/antha-lang/antha/graph"
	"github.com/antha-lang/antha/target"
	"github.com/golang/protobuf/proto"
)

var (
	cannotMergeUnequalCalls = errors.New("cannot merge unequal calls")
)

// Handle generic rpc calls
type Handler struct {
	Labels []ast.NameValue
}

func (a *Handler) String() string {
	return "Handler"
}

func (a *Handler) CanCompile(req ast.Request) bool {
	can := ast.Request{
		Selector: a.Labels,
	}
	if !req.Matches(can) {
		return false
	}

	return can.Contains(req)
}

func (a *Handler) MoveCost(target.Device) int {
	return 0
}

func serialize(calls []driver.Call) ([]byte, error) {
	var buf bytes.Buffer
	for _, c := range calls {
		if _, err := io.Copy(&buf, bytes.NewReader([]byte(c.Method))); err != nil {
			return nil, err
		}
		if bs, err := proto.Marshal(c.Args); err != nil {
			return nil, err
		} else if _, err := io.Copy(&buf, bytes.NewReader(bs)); err != nil {
			return nil, err
		}
	}
	return buf.Bytes(), nil
}

func (a *Handler) merge(nodes []ast.Node) (*ast.HandleInst, error) {
	if len(nodes) == 0 {
		return nil, nil
	}

	ret := nodes[0].(*ast.Command).Inst.(*ast.HandleInst)
	retBs, err := serialize(ret.Calls)
	if err != nil {
		return nil, err
	}

	for _, n := range nodes[1:] {
		h := n.(*ast.Command).Inst.(*ast.HandleInst)
		hBs, err := serialize(h.Calls)
		if err != nil {
			return nil, err
		}

		if !bytes.Equal(retBs, hBs) {
			return nil, cannotMergeUnequalCalls
		}
	}
	return ret, nil
}

func (a *Handler) Compile(nodes []ast.Node) ([]target.Inst, error) {
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
			c, ok := n.(*ast.Command)
			if !ok {
				return nil, fmt.Errorf("unexpected node %T", r)
			}
			h, ok := c.Inst.(*ast.HandleInst)
			if !ok {
				return nil, fmt.Errorf("unexpected node %T", c.Inst)
			}

			key := h.Group
			same[key] = append(same[key], n)
			next = append(next, dag.Visit(r)...)
		}
		// Apply
		for _, nodes := range same {
			h, err := a.merge(nodes)
			if err != nil {
				return nil, err
			}
			if h == nil {
				continue
			}

			in := &target.Run{
				Dev:   a,
				Label: h.Group,
				Calls: h.Calls,
			}

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
