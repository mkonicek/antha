package mixer

import (
	"bytes"
	"encoding/csv"
	"strconv"

	"github.com/antha-lang/antha/antha/anthalib/wtype"
)

// MarshalPlateCSV writes a plate to a CSV file
func MarshalPlateCSV(plate *wtype.Plate) ([]byte, error) {
	var records [][]string
	records = append(records, []string{
		plate.Type,
		plate.PlateName,
		"LiquidType",
		"Vol",
		"Vol Unit",
		"Conc",
		"Conc Unit",
	})
	for _, well := range plate.AllNonEmptyWells() {
		comp := well.WContents
		records = append(records, []string{
			well.Crds.FormatA1(),
			comp.CName,
			comp.TypeName(),
			strconv.FormatFloat(comp.Vol, 'g', -1, 64),
			comp.Vunit,
			strconv.FormatFloat(comp.Conc, 'g', -1, 64),
			comp.Cunit,
		})
	}

	var out bytes.Buffer
	w := csv.NewWriter(&out)
	err := w.WriteAll(records)
	return out.Bytes(), err
}
