package platereader

import (
	"fmt"

	"github.com/antha-lang/antha/ast"
	"github.com/antha-lang/antha/target"
	"github.com/antha-lang/antha/target/handler"
)

// PlateReader defines the interface to a plate reader device
type PlateReader struct {
	handler.GenericHandler
}

// NewWOPlateReader makes a new handle-based write-only plate reader
func NewWOPlateReader() *PlateReader {
	ret := &PlateReader{}
	ret.GenericHandler = handler.GenericHandler{
		Labels: []ast.NameValue{
			target.DriverSelectorV1WriteOnlyPlateReader,
		},
		GenFunc: ret.generate,
	}

	return ret
}

func (pr *PlateReader) generate(cmd interface{}) ([]target.Inst, error) {
	handle, ok := cmd.(*ast.HandleInst)

	if !ok {
		return []target.Inst{}, fmt.Errorf("Plate reader received wrong kind of instruction: expected *ast.HandleInst got %T", cmd)
	}

	inst := &target.Run{
		Dev:   pr,
		Label: handle.Group,
		Calls: handle.Calls,
	}

	return []target.Inst{inst}, nil
}
