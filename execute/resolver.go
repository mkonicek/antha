package execute

import (
	"fmt"

	"github.com/antha-lang/antha/ast"
	"github.com/antha-lang/antha/codegen"
	"github.com/antha-lang/antha/target"
	"golang.org/x/net/context"
)

// Converts execute instructions to their ast equivalents
type resolver struct {
	nodes []ast.Node    // Squirrel away state for Run()
	insts []target.Inst // Squirrel away state for Run()
}

// Called by trace to resolve blocked instructions
func (a *resolver) resolve(ctx context.Context, instObjs []interface{}) (map[int]interface{}, error) {
	ret := make(map[int]interface{})
	var commands []*commandInst
	for idx, in := range instObjs {
		switch inst := in.(type) {
		case nil:
		case *commandInst:
			commands = append(commands, inst)
		default:
			return nil, fmt.Errorf("invalid instruction: %T", inst)
		}

		// TODO(ddn): Under prefix execution, we need to notify someone to
		// execute instructions to unblock the resolver. Now, since we don't do
		// prefix execution, we can just mark as done.
		ret[idx] = nil
	}

	t, err := target.GetTarget(ctx)
	if err != nil {
		return nil, err
	}

	nodes, err := getMaker(ctx).MakeNodes(commands)
	if err != nil {
		return nil, err
	}

	insts, err := codegen.Compile(t, nodes)
	if err != nil {
		return nil, err
	}

	a.insts = append(a.insts, insts...)
	return ret, nil
}
