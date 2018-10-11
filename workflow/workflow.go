// Package workflow implements DAG scheduling of networks of functions
// generated at runtime. The execution uses the inject package to allow
// late-binding of functions.
package workflow

import (
	"context"
	"errors"
	"fmt"
	"sync"

	api "github.com/antha-lang/antha/api/v1"
	"github.com/antha-lang/antha/inject"
)

// TODO: deterministic node name/order

var (
	errCyclicWorkflow  = errors.New("cyclic workflow")
	errUnknownPort     = errors.New("unknown port")
	errUnknownProcess  = errors.New("unknown process")
	errAlreadyAssigned = errors.New("already assigned")
	errAlreadyRemoved  = errors.New("already removed")
)

// A Port is a unique identifier for an input or output parameter
type Port struct {
	Process string `json:"process"`
	Port    string `json:"port"`
}

// String returns a string representation of a port
func (a Port) String() string {
	return fmt.Sprintf("%s.%s", a.Process, a.Port)
}

// A Process is an instance of a component / element execution
type Process struct {
	Component string         `json:"component"`
	Metadata  screenPosition `json:"metadata"`
}

type screenPosition struct {
	X int `json:"x"`
	Y int `json:"y"`
}

// A Connection connects the output of one Process to the input of another
type Connection struct {
	Src Port `json:"source"`
	Tgt Port `json:"target"`
}

// Desc is the description of a workflow.
type Desc struct {
	Processes   map[string]Process `json:"Processes"`
	Connections []Connection       `json:"connections"`
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

func (a *node) removeIn(port string) error {
	a.lock.Lock()
	defer a.lock.Unlock()
	if !a.Ins[port] {
		return errAlreadyRemoved
	}
	delete(a.Ins, port)
	return nil
}

func (a *node) hasIns() bool {
	return len(a.Ins) > 0
}

func (a *node) setParam(port string, value interface{}) error {
	a.lock.Lock()
	defer a.lock.Unlock()
	if _, seen := a.Params[port]; seen {
		return errAlreadyAssigned
	}
	a.Params[port] = value
	return nil
}

// Workflow is the state to execute a workflow
type Workflow struct {
	nodes   map[string]*node
	Outputs map[Port]interface{} // Values generated that were not connected to another process
}

// FuncName gets the function to be called for the given process name
func (a *Workflow) FuncName(process string) (string, error) {
	n, ok := a.nodes[process]
	if !ok {
		return "", errUnknownProcess
	}
	return n.FuncName, nil
}

// SetParam sets initial parameter values before executing
func (a *Workflow) SetParam(port Port, value interface{}) error {
	n := a.nodes[port.Process]
	if n == nil {
		return errUnknownPort
	} else if n.Ins[port.Port] {
		return errAlreadyAssigned
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
	query := inject.NameQuery{
		Repo:  n.FuncName,
		Stage: api.ElementStage_STEPS,
	}
	out, err := inject.Call(ctx, query, n.Params)

	if err != nil {
		return nil, err
	}

	if err := updateOutParams(n, out, a.Outputs); err != nil {
		return nil, err
	}

	var roots []*node
	for _, eps := range n.Outs {
		for _, ep := range eps {
			if err := ep.Node.removeIn(ep.Port); err != nil {
				return nil, fmt.Errorf("error removing in edge on %q: %s", ep, err)
			} else if !ep.Node.hasIns() {
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
		if !n.hasIns() {
			roots = append(roots, n)
		}
	}
	if len(roots) == 0 && len(nodes) > 0 {
		return nil, errCyclicWorkflow
	}
	return roots, nil
}

// Run a workflow
func (a *Workflow) Run(ctx context.Context) error {
	worklist, err := makeRoots(a.nodes)
	if err != nil {
		return err
	}

	for len(worklist) > 0 {
		n := worklist[0]
		worklist = worklist[1:]

		if moreWork, err := a.run(ctx, n); err != nil {
			return fmt.Errorf("cannot run process %q: %s", n.Process, err)
		} else {
			worklist = append(worklist, moreWork...)
		}
	}
	if len(a.nodes) > 0 { // by definition, len(worklist) == 0
		return errCyclicWorkflow
	}

	return nil
}

// AddNode adds a process to a workflow that executes funcName
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

// AddEdge connects an output of one process to an input of another
func (a *Workflow) AddEdge(src, tgt Port) error {
	snode := a.nodes[src.Process]
	if snode == nil {
		return fmt.Errorf("unknown source port %q", src)
	}
	tnode := a.nodes[tgt.Process]
	if tnode == nil {
		return fmt.Errorf("unknown target port %q", tgt)
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

// Opt are options for creating a new Workflow
type Opt struct {
	FromDesc *Desc
}

// New creates a new Workflow
func New(opt Opt) (*Workflow, error) {
	w := &Workflow{
		nodes:   make(map[string]*node),
		Outputs: make(map[Port]interface{}),
	}

	var desc *Desc
	if opt.FromDesc != nil {
		desc = opt.FromDesc
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
