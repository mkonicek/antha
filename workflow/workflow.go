// Package workflow implements DAG scheduling of networks of functions
// generated at runtime. The execution uses the inject package to allow
// late-binding of functions.
package workflow

import (
	"encoding/json"
	"errors"
	"fmt"
	"sync"

	"github.com/antha-lang/antha/bvendor/golang.org/x/net/context"
	"github.com/antha-lang/antha/inject"
	"github.com/antha-lang/antha/trace"
)

// TODO: deterministic node name/order

var (
	cyclicWorkflow  = errors.New("cyclic workflow")
	unknownPort     = errors.New("unknown port")
	unknownProcess  = errors.New("unknown process")
	alreadyAssigned = errors.New("already assigned")
	alreadyRemoved  = errors.New("already removed")
)

// Unique identifier for an input or output parameter
type Port struct {
	Process string
	Port    string
}

func (a Port) String() string {
	return fmt.Sprintf("%s.%s", a.Process, a.Port)
}

type Process struct {
	Component string
}

type Connection struct {
	Src Port
	Tgt Port
}

// Description of a workflow. Structure inherited from and is a subset of
// noflow library
type Desc struct {
	Processes   map[string]Process
	Connections []Connection
}

type endpoint struct {
	Port string
	Node *node
}

func (a endpoint) String() string {
	return fmt.Sprintf("%s.%s", a.Node.Process, a.Port)
}

type node struct {
	lock     sync.Mutex            // Lock on Params and Ins during Execute
	Process  string                // Name of this instance
	FuncName string                // Function that should be called
	Params   inject.Value          // Parameters to this function
	Outs     map[string][]endpoint // Out edges
	Ins      map[string]bool       // In edges
}

func (a *node) removeIn(port string) (int, error) {
	a.lock.Lock()
	defer a.lock.Unlock()
	if !a.Ins[port] {
		return 0, alreadyRemoved
	}
	delete(a.Ins, port)
	return len(a.Ins), nil
}

func (a *node) setParam(port string, value interface{}) error {
	a.lock.Lock()
	defer a.lock.Unlock()
	if _, seen := a.Params[port]; seen {
		return alreadyAssigned
	}
	a.Params[port] = value
	return nil
}

// State to execute a workflow
type Workflow struct {
	roots   []*node
	nodes   map[string]*node
	Outputs map[Port]interface{} // Values generated that were not connected to another process
}

func (a *Workflow) FuncName(process string) (string, error) {
	if n, ok := a.nodes[process]; !ok {
		return "", unknownProcess
	} else {
		return n.FuncName, nil
	}
}

// Set initial parameter values before executing
func (a *Workflow) SetParam(port Port, value interface{}) error {
	n := a.nodes[port.Process]
	if n == nil {
		return unknownPort
	} else if n.Ins[port.Port] {
		return alreadyAssigned
	} else {
		return n.setParam(port.Port, value)
	}
}

func updateOutParams(n *node, out inject.Value, unmatched map[Port]interface{}) error {
	seen := make(map[string]bool)
	for name, value := range out {
		seen[name] = true
		if eps := n.Outs[name]; len(eps) == 0 {
			port := Port{Port: name, Process: n.Process}
			if _, seen := unmatched[port]; seen {
				return fmt.Errorf("%q already assigned", endpoint{Port: name, Node: n})
			}
			unmatched[port] = value
		} else {
			for _, ep := range eps {
				if err := ep.Node.setParam(ep.Port, value); err != nil {
					return fmt.Errorf("error setting parameter on %q: %s", ep, err)
				}
			}
		}
	}
	for name := range n.Outs {
		if !seen[name] {
			return fmt.Errorf("missing value for %q", endpoint{Port: name, Node: n})
		}
	}
	return nil
}

func (a *Workflow) run(ctx context.Context, n *node) ([]*node, error) {
	out, err := inject.Call(ctx, inject.NameQuery{Repo: n.FuncName}, n.Params)
	if err != nil {
		return nil, err
	}
	if err := updateOutParams(n, out, a.Outputs); err != nil {
		return nil, err
	}

	var roots []*node
	for _, eps := range n.Outs {
		for _, ep := range eps {
			if remaining, err := ep.Node.removeIn(ep.Port); err != nil {
				return nil, fmt.Errorf("error removing in edge on %q: %s", ep, err)
			} else if remaining == 0 {
				roots = append(roots, ep.Node)
			}
		}
	}
	delete(a.nodes, n.Process)
	return roots, nil
}

func makeRoots(nodes map[string]*node) ([]*node, error) {
	var roots []*node
	for _, n := range nodes {
		if len(n.Ins) == 0 {
			roots = append(roots, n)
		}
	}
	if len(roots) == 0 && len(nodes) > 0 {
		return nil, cyclicWorkflow
	}
	return roots, nil
}

// Run a workflow
func (a *Workflow) Run(parent context.Context) error {
	roots, err := makeRoots(a.nodes)
	if err != nil {
		return err
	}

	ctx, cancel, allDone := trace.NewContext(parent)
	defer cancel()

	trace.Go(ctx, func(ctx context.Context) error {
		for len(roots) != 0 {
			var newRoots []*node
			// TODO: Parallelize this loop
			for _, n := range roots {
				if rs, err := a.run(ctx, n); err != nil {
					return fmt.Errorf("cannot run process %q: %s", n.Process, err)
				} else {
					newRoots = append(newRoots, rs...)
				}
			}

			if len(newRoots) == 0 && len(a.nodes) > 0 {
				return cyclicWorkflow
			}
			roots = newRoots
		}
		return nil
	})

	<-allDone()
	return ctx.Err()
}

// Add a process to a workflow that executes funcName
func (a *Workflow) AddNode(process, funcName string) error {
	if a.nodes[process] != nil {
		return fmt.Errorf("process %q already defined", process)
	}
	n := &node{
		Process:  process,
		FuncName: funcName,
		Params:   make(inject.Value),
		Outs:     make(map[string][]endpoint),
		Ins:      make(map[string]bool),
	}
	a.nodes[process] = n
	return nil
}

// Connect an output of one process to an input of another
func (a *Workflow) AddEdge(src, tgt Port) error {
	snode := a.nodes[src.Process]
	if snode == nil {
		return fmt.Errorf("unknown process %q", src)
	}
	tnode := a.nodes[tgt.Process]
	if tnode == nil {
		return fmt.Errorf("unknown process %q", src)
	}

	sport := src.Port
	tport := tgt.Port
	if _, seen := tnode.Ins[tport]; seen {
		return fmt.Errorf("port %q of process %q already assigned", endpoint{Port: tport, Node: tnode}, tgt.Process)
	}
	tnode.Ins[tport] = true
	snode.Outs[sport] = append(snode.Outs[sport], endpoint{Port: tport, Node: tnode})
	return nil
}

// Options for creating a new Workflow
type Opt struct {
	FromBytes []byte
	FromDesc  *Desc
}

// Create a new Workflow
func New(opt Opt) (*Workflow, error) {
	w := &Workflow{
		nodes:   make(map[string]*node),
		Outputs: make(map[Port]interface{}),
	}

	var desc *Desc
	if opt.FromDesc != nil {
		desc = opt.FromDesc
	} else if opt.FromBytes != nil {
		if err := json.Unmarshal(opt.FromBytes, &desc); err != nil {
			return nil, err
		}
	} else {
		desc = &Desc{}
	}

	for name, process := range desc.Processes {
		if err := w.AddNode(name, process.Component); err != nil {
			return nil, err
		}
	}

	for _, c := range desc.Connections {
		if err := w.AddEdge(c.Src, c.Tgt); err != nil {
			return nil, err
		}
	}
	return w, nil
}
