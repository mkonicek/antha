package handler

import (
	"fmt"

	"github.com/antha-lang/antha/ast"
	"github.com/antha-lang/antha/target"
)

// A Handler handles generic rpc calls
type Handler struct {
	*GenericHandler
}

// New creates a new handler for the given selector labels
func New(labels []ast.NameValue) *Handler {
	return &Handler{
		GenericHandler: &GenericHandler{
			Labels: labels,
			GenFunc: func(dev target.Device, cmd interface{}) ([]target.Inst, error) {
				h, ok := cmd.(*ast.HandleInst)
				if !ok {
					return nil, fmt.Errorf("expecting %T found %T instead", h, cmd)
				}
				return []target.Inst{
					&target.Run{
						Dev:   dev,
						Label: h.Group,
						Calls: h.Calls,
					},
				}, nil
			},
			FilterFieldsForKey: func(cmd interface{}) (interface{}, error) {
				h, ok := cmd.(*ast.HandleInst)
				if !ok {
					return nil, fmt.Errorf("expecting %T found %T instead", h, cmd)
				}
				return &ast.HandleInst{
					Calls: h.Calls,
					Group: h.Group,
				}, nil
			},
		},
	}
}
