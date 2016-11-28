package mixer

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/antha/anthalib/wutil"
	"github.com/antha-lang/antha/microArch/factory"
)

type ParsePlateResult struct {
	Plate    *wtype.LHPlate
	Warnings []string
}

func validName(name string) error {
	invalid := "+"
	valid := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	if strings.ContainsAny(name, invalid) {
		return fmt.Errorf("name %q contains an invalid character %q", name, invalid)
	}

	if !strings.ContainsAny(name, valid) {
		return fmt.Errorf("name %q must contain at least one alphanumeric character", name)
	}
	return nil
}

func validWellCoord(coord string) (wtype.WellCoords, error) {
	well := wtype.MakeWellCoords(coord)
	if len(coord) == 0 {
		return well, fmt.Errorf("no well coords")
	}
	if well.IsZero() {
		return well, fmt.Errorf("cannot parse well coord %q", coord)
	}
	return well, nil
}

func validWell(well wtype.WellCoords, plate *wtype.LHPlate) error {
	if well.X >= plate.WellsX() || well.Y >= plate.WellsY() {
		return fmt.Errorf("well coord %q does not exist on plate type %q", well.FormatA1(), plate.Type)
	}
	return nil
}

// CSV plate format: (? denotes optional, whitespace for clarity)
//
//   <plate type> , <plate name ?>
//   <well0> , <component name0> , <component type0 ?> , <volume0 ?> , <volume unit0 ?>, <conc0 ?> , <conc unit0 ?>
//   <well1> , <component name1> , <component type1 ?> , <volume1 ?> , <volume unit1 ?>, <conc1 ?> , <conc unit1 ?>
//   ...
//
// TODO: refactor if/when Opt loses raw []byte and file as InputPlate options
func ParsePlateCSV(inData io.Reader) (*ParsePlateResult, error) {
	// Get returning "" if idx >= len(xs)
	get := func(xs []string, idx int) string {
		if len(xs) <= idx {
			return ""
		}
		return strings.TrimSpace(xs[idx])
	}

	parseUnit := func(value, unit, defaultUnit string) (float64, string, error) {
		v := wutil.ParseFloat(value)

		var err error
		if v == 0.0 {
			err = fmt.Errorf("cannot parse value %q", value)
		}

		if len(unit) == 0 {
			unit = defaultUnit
			err = fmt.Errorf("cannot parse unit %q", unit)
		}

		return v, unit, err
	}

	csvr := csv.NewReader(inData)
	csvr.FieldsPerRecord = -1

	rec, err := csvr.Read()
	if err != nil {
		return nil, err
	}

	if len(rec) == 0 {
		return nil, fmt.Errorf("empty file")
	}

	plateType := rec[0]
	plate := factory.GetPlateByType(plateType)
	if plate == nil {
		return nil, fmt.Errorf("unknown plate type: ", plateType)
	}

	plate.PlateName = get(rec, 1)
	if len(plate.PlateName) == 0 {
		plate.PlateName = fmt.Sprint("input_plate_", plate.ID)
	}

	var warnings []string
	for lineNo := 1; true; lineNo += 1 {
		rec, err := csvr.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		wellField := get(rec, 0)
		cname := get(rec, 1)
		ctypeField := get(rec, 2)

		well, err := validWellCoord(wellField)
		if err != nil {
			warnings = append(warnings, fmt.Sprintf("line %d skipped: %s", lineNo, err))
			continue
		} else if err := validWell(well, plate); err != nil {
			warnings = append(warnings, fmt.Sprintf("line %d skipped: %s", lineNo, err))
			continue
		}

		if err := validName(cname); err != nil {
			warnings = append(warnings, fmt.Sprintf("line %d skipped: %s", lineNo, err))
			continue
		}

		ctype, err := wtype.LiquidTypeFromString(ctypeField)
		if err != nil {
			warnings = append(warnings, fmt.Sprintf("line %d: unknown component type %q, defaulting to %q: %s", lineNo, ctypeField, wtype.LiquidTypeName(ctype), err))
		}

		vol, vunit, err := parseUnit(get(rec, 3), get(rec, 4), "ul")
		if err != nil {
			warnings = append(warnings, fmt.Sprintf("line %d: unknown volume %q, defaulting to \"%f%s\": %s", lineNo, get(rec, 3)+get(rec, 4), vol, vunit, err))
		}
		volume := wunit.NewVolume(vol, vunit)

		conc, cunit, err := parseUnit(get(rec, 5), get(rec, 6), "ng/ul")
		if err != nil {
			warnings = append(warnings, fmt.Sprintf("line %d: unknown volume %q, defaulting to \"%f%s\": %s", lineNo, get(rec, 5)+get(rec, 6), conc, cunit, err))
		}
		concentration := wunit.NewConcentration(conc, cunit)

		// Make component
		cmp := wtype.NewLHComponent()

		cmp.Vol = volume.RawValue()
		cmp.Vunit = volume.Unit().PrefixedSymbol()
		cmp.CName = cname
		cmp.Type = ctype

		cmp.Conc = concentration.RawValue()
		cmp.Cunit = concentration.Unit().PrefixedSymbol()

		plate.WellAt(well).Add(cmp)
	}

	return &ParsePlateResult{
		Plate:    plate,
		Warnings: warnings,
	}, nil
}

func parsePlateFile(filename string) (*ParsePlateResult, error) {
	f, err := os.Open(filename)

	if err != nil {
		return nil, err
	}

	defer f.Close()
	return ParsePlateCSV(f)
}

// Convenience function for parsing a plate from file. Will splat out
// warnings to stdout.
func ParseInputPlateFile(filename string) (*wtype.LHPlate, error) {
	r, err := parsePlateFile(filename)
	if err != nil {
		return nil, err
	}
	for _, warning := range r.Warnings {
		log.Println(warning)
	}
	return r.Plate, nil
}
