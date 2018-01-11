package handler

import (
	"github.com/antha-lang/antha/ast"
	"github.com/antha-lang/antha/target"
)

// NewWOPlateReader makes a new handle-based write-only plate reader
func NewWOPlateReader() *Handler {
	return New([]ast.NameValue{target.DriverSelectorV1WriteOnlyPlateReader})
}
