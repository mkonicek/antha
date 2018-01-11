package handler

import (
	"github.com/antha-lang/antha/ast"
	"github.com/antha-lang/antha/target"
)

func NewWOPlateReader() *Handler {
	return New([]ast.NameValue{target.DriverSelectorV1WriteOnlyPlateReader})
}
