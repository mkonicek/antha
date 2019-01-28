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

func (h *Handler) generate(cmd interface{}) ([]ast.Inst, error) {
	hinst, ok := cmd.(*ast.HandleInst)
	if !ok {
		return nil, fmt.Errorf("expecting %T found %T instead", hinst, cmd)
	}
	return []ast.Inst{
		&target.Run{
			Dev:   h,
			Label: hinst.Group,
			Calls: hinst.Calls,
		},
	}, nil
}

// New creates a new handler for the given selector labels
func New(labels []ast.NameValue) *Handler {
	h := &Handler{}
	h.GenericHandler = &GenericHandler{
		Labels:  labels,
		GenFunc: h.generate,
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
	}
	return h
}
