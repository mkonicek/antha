package datasource

import (
	"context"
	"fmt"

	"github.com/antha-lang/antha/ast"
	"github.com/antha-lang/antha/target"
	"github.com/antha-lang/antha/target/human"
)

var (
	_ target.Device = (*DataSource)(nil)
)

// A DataSource is a generic device that provides data asynchronously
type DataSource struct{}

// CanCompile implements a device
func (ds *DataSource) CanCompile(req ast.Request) bool {
	can := ast.Request{
		Selector: []ast.NameValue{
			target.DriverSelectorV1DataSource,
		},
	}
	return can.Contains(req)
}

// MoveCost implements a device
func (ds *DataSource) MoveCost(from target.Device) int64 {
	return human.HumanByXCost + 1 // same as a mixer
}

// Compile implements a device
func (ds *DataSource) Compile(ctx context.Context, cmds []ast.Node) (insts []target.Inst, err error) {
	if len(cmds) != 1 {
		return nil, fmt.Errorf("multiple GetData commands not supported (%d received)", len(cmds))
	}

	node := cmds[0]

	c, ok := node.(*ast.Command)
	if !ok {
		return nil, fmt.Errorf("cannot compile %T", node)
	}

	m, ok := c.Inst.(*ast.AwaitInst)
	if !ok {
		return nil, fmt.Errorf("cannot compile %T", c.Inst)
	}

	return []target.Inst{
		&target.AwaitData{
			Dev:  ds,
			Inst: m,
		},
	}, nil
}
