package execute

import (
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/ast"
)

type maker struct {
	byComp map[*wtype.LHComponent]*ast.UseComp
	// pre-mixing samples of the same component share the same id
	byId map[string][]*ast.UseComp
}

func (a *maker) makeComp(c *wtype.LHComponent) *ast.UseComp {
	u, ok := a.byComp[c]
	if !ok {
		u = &ast.UseComp{
			Value: c,
		}
		a.byComp[c] = u
	}

	a.byId[c.ID] = append(a.byId[c.ID], u)

	return u
}

func (a *maker) makeCommand(in *commandInst) ast.Node {
	for _, arg := range in.Args {
		in.Command.From = append(in.Command.From, a.makeComp(arg))
	}

	out := a.makeComp(in.Comp)
	out.From = append(out.From, in.Command)
	return out
}

// Manifest dependencies between samples that share the same id
func (a *maker) addMissingDeps() {
	for _, uses := range a.byId {
		// HACK: assume that samples are used in sequentially; remove when
		// dependencies are tracked individually

		// Make sure we don't introduce any loops
		seen := make(map[*ast.UseComp]bool)
		var us []*ast.UseComp
		for _, u := range uses {
			if seen[u] {
				continue
			}
			seen[u] = true
			us = append(us, u)
		}
		for idx, u := range us {
			if idx == 0 {
				continue
			}
			u.From = append(u.From, uses[idx-1])
		}
	}
}

// Normalize commands into well-formed AST
func makeNodes(insts []*commandInst) ([]ast.Node, error) {
	m := &maker{
		byComp: make(map[*wtype.LHComponent]*ast.UseComp),
		byId:   make(map[string][]*ast.UseComp),
	}
	var nodes []ast.Node
	for _, inst := range insts {
		nodes = append(nodes, m.makeCommand(inst))
	}

	m.addMissingDeps()

	return nodes, nil
}
