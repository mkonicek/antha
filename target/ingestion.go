package target

import (
	"context"
	"github.com/antha-lang/antha/ast"
)

func GetExpectDataCommand(ctx context.Context, params []byte, selectors ...ast.NameValue) *ast.Command {
	tgt, err := GetTarget(ctx)
	if err != nil {
		panic("Error when finding target in Context: " + err.Error())
	}
	// do the append this way around so that we guarantee we're copying
	// the selectors into a new array, and not altering the original
	// underlying array:
	selectorsWithDataSource := append([]ast.NameValue{DriverSelectorV1DataSource}, selectors...)
	devs := tgt.CanCompile(ast.Request{Selector: selectorsWithDataSource})
	if len(devs) == 0 {
		return nil
	} else {
		templateInst := devs[0].ExpectDataTemplate() // FIXME when we cope with > 1 machine in a workflow
		templateInst.Params = params
		return &ast.Command{
			Requests: []ast.Request{{Selector: selectorsWithDataSource}},
			Inst:     templateInst,
		}
	}
}
