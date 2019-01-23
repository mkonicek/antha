package handler

import (
	"fmt"

	"github.com/antha-lang/antha/laboratory/effects"
	"github.com/antha-lang/antha/target"
)

// A Handler handles generic rpc calls
type Handler struct {
	*GenericHandler
}

func (h *Handler) generate(cmd interface{}) ([]effects.Inst, error) {
	hinst, ok := cmd.(*effects.HandleInst)
	if !ok {
		return nil, fmt.Errorf("expecting %T found %T instead", hinst, cmd)
	}
	return []effects.Inst{
		&target.Run{
			Dev:   h,
			Label: hinst.Group,
			Calls: hinst.Calls,
		},
	}, nil
}

// New creates a new handler for the given selector labels
func New(labels []effects.NameValue) *Handler {
	h := &Handler{}
	h.GenericHandler = &GenericHandler{
		Labels:  labels,
		GenFunc: h.generate,
		FilterFieldsForKey: func(cmd interface{}) (interface{}, error) {
			h, ok := cmd.(*effects.HandleInst)
			if !ok {
				return nil, fmt.Errorf("expecting %T found %T instead", h, cmd)
			}
			return &effects.HandleInst{
				Calls: h.Calls,
				Group: h.Group,
			}, nil
		},
	}
	return h
}
