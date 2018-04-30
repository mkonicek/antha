package mixer

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/inventory"
)

// ParsePlateResult is the result of parsing a plate
type ParsePlateResult struct {
	Plate    *wtype.LHPlate
	Warnings []string
}

// ValidationConfig specifies how to parse a plate
type ValidationConfig struct {
	valid   string
	invalid string
}

// ValidChars are the characters that are valid in a component name
func (vc ValidationConfig) ValidChars() string {
	return vc.valid
}

// InvalidChars are the characters that are invalid in a component name
func (vc ValidationConfig) InvalidChars() string {
	return vc.invalid
}

// DefaultValidationConfig is the default validation config
func DefaultValidationConfig() ValidationConfig {
	return ValidationConfig{
		invalid: "+",
		valid:   "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789",
	}
}

// PermissiveValidationConfig is a looser validation config
func PermissiveValidationConfig() ValidationConfig {
	return ValidationConfig{
		invalid: "",
		valid:   "+abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789",
	}
}

func validName(name string, vc ValidationConfig) error {

	invalid := vc.InvalidChars()
	valid := vc.ValidChars()

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

// ParsePlateCSV parses a plate
func ParsePlateCSV(ctx context.Context, inData io.Reader) (*ParsePlateResult, error) {
	return ParsePlateCSVWithValidationConfig(ctx, inData, DefaultValidationConfig())
}

// ParsePlateCSVWithValidationConfig parses a csv file into a plate.
//
// CSV plate format: (? denotes optional, whitespace for clarity)
//
//   <plate type> , <plate name ?>
//   <well0> , <component name0> , <component type0 ?> , <volume0 ?> , <volume unit0 ?>, <conc0 ?> , <conc unit0 ?>
//   <well1> , <component name1> , <component type1 ?> , <volume1 ?> , <volume unit1 ?>, <conc1 ?> , <conc unit1 ?>
//   ...
//
// TODO: refactor if/when Opt loses raw []byte and file as InputPlate options
func ParsePlateCSVWithValidationConfig(ctx context.Context, inData io.Reader, vc ValidationConfig) (*ParsePlateResult, error) {
	// Get returning "" if idx >= len(xs)
	get := func(xs []string, idx int) string {
		if len(xs) <= idx {
			return ""
		}
		return strings.TrimSpace(xs[idx])
	}

	parseUnit := func(value, unit, defaultUnit string) (float64, string, error) {
		v, err := strconv.ParseFloat(value, 64)
		if err != nil {
			err = fmt.Errorf("cannot parse value %q: %s", value, err)
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
	plate, err := inventory.NewPlate(ctx, plateType)
	if err != nil {
		return nil, fmt.Errorf("cannot make plate %s: %s", plateType, err)
	}

	plate.PlateName = get(rec, 1)
	if len(plate.PlateName) == 0 {
		plate.PlateName = fmt.Sprint("input_plate_", plate.ID)
	}

	var warnings []string
	for lineNo := 1; true; lineNo++ {
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

		if err := validName(cname, vc); err != nil {
			warnings = append(warnings, fmt.Sprintf("line %d skipped: %s", lineNo, err))
			continue
		}

		ctype, err := wtype.LiquidTypeFromString(wtype.PolicyName(ctypeField))
		if err != nil {
			warnings = append(warnings, fmt.Sprintf("line %d: component type %q not found in system types, using (%s); this may generate undesirable behaviour if this is not a custom type: %s", lineNo, ctypeField, ctype, err))
		}

		vol, vunit, err := parseUnit(get(rec, 3), get(rec, 4), "ul")
		if err != nil {
			warnings = append(warnings, fmt.Sprintf("line %d: unknown volume %q, defaulting to \"%f%s\": %s", lineNo, get(rec, 3)+get(rec, 4), vol, vunit, err))
		}
		volume := wunit.NewVolume(vol, vunit)

		conc, cunit, err := parseUnit(get(rec, 5), get(rec, 6), "g/l")
		if err != nil {
			warnings = append(warnings, fmt.Sprintf("line %d: unknown concentration %q, defaulting to \"%f%s\": %s", lineNo, get(rec, 5)+get(rec, 6), conc, cunit, err))
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

		if wa, ok := plate.WellAt(well); ok {
			err = wa.AddComponent(cmp)
			if err != nil {
				return nil, err
			}

			// this should be defined elsewhere
			wa.WContents.DeclareInstance()
		} else {
			return nil, fmt.Errorf("Unknown location \"%s\" in plate \"%s\"", well.FormatA1(), plate.Name())
		}
	}

	if err = plate.ValidateVolumes(); err != nil {
		return nil, err
	}

	return &ParsePlateResult{
		Plate:    plate,
		Warnings: warnings,
	}, nil
}

func parsePlateFile(ctx context.Context, filename string) (*ParsePlateResult, error) {
	f, err := os.Open(filename)

	if err != nil {
		return nil, err
	}

	defer f.Close() // nolint: errcheck

	return ParsePlateCSV(ctx, f)
}

// ParseInputPlateFile is convenience function for parsing a plate from file.
// Will splat out warnings to stdout.
func ParseInputPlateFile(ctx context.Context, filename string) (*wtype.LHPlate, error) {
	r, err := parsePlateFile(ctx, filename)
	if err != nil {
		return nil, err
	}
	for _, warning := range r.Warnings {
		log.Println(warning)
	}
	return r.Plate, nil
}
