package migrate

import (
	"github.com/antha-lang/antha/antha/anthalib/wtype"
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
