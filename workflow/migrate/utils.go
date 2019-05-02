package migrate

import (
	"fmt"

	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/workflow"
)

// UpdatePlateTypes generates an array of wtype.PlateTypeName values from an
// array of strings
func UpdatePlateTypes(names []string) []wtype.PlateTypeName {
	ptnames := make([]wtype.PlateTypeName, len(names))
	for i, v := range names {
		ptnames[i] = wtype.PlateTypeName(v)
	}
	return ptnames
}

// UniqueElementType returns the unique element type for a given element type
// name across a set of repositories, or an error if the type cannot be
// unambiguously resolved.
func UniqueElementType(types workflow.ElementTypesByRepository, name workflow.ElementTypeName) (*workflow.ElementType, error) {
	var et *workflow.ElementType
	for _, rmap := range types {
		if v, found := rmap[name]; found {
			if et != nil {
				return nil, fmt.Errorf("element type %v is found in multiple repositories", name)
			}
			et = &v
		}
	}

	if et == nil {
		return nil, fmt.Errorf("element type %v could not be found in the supplied repositories", name)
	}
	return et, nil
}
