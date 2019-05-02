package mixer

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/inventory"
	"github.com/antha-lang/toolbox/csvutil"
	"github.com/pkg/errors"
)

const plateTypeReplacementKey string = "PlateType"

// ParsePlateResult is the result of parsing a plate
type ParsePlateResult struct {
	Plate    *wtype.Plate
	Warnings []string
}

// ValidationConfig specifies how to parse a plate
type ValidationConfig struct {
	valid        string
	invalid      string
	replaceField map[string]string
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

// ReplacePlateType will replace the plate type in the CSV file with the replacement option specified.
func ReplacePlateType(replacement string) ValidationConfig {
	return ValidationConfig{
		replaceField: map[string]string{plateTypeReplacementKey: replacement},
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

func validWell(well wtype.WellCoords, plate *wtype.Plate) error {
	if well.X >= plate.WellsX() || well.Y >= plate.WellsY() {
		return fmt.Errorf("well coord %q does not exist on plate type %q", well.FormatA1(), plate.Type)
	}
	return nil
}

// ParsePlateCSV parses a csv file into a plate.
// ValidationConfig options may be added which either:
// (i) define sets of valid and invalid characters which may be used in the component names.
// (ii) define fields in the Plate csv which may be replaced.
// If no config is set for case (i) a default config will be used.
//
// CSV plate format: (? denotes optional, whitespace for clarity)
//
//   <plate type ?> , <plate name ?>
//   <well0> , <component name0> , <component type0 ?> , <volume0 ?> , <volume unit0 ?>, <conc0 ?> , <conc unit0 ?>, <SubComponent1Name: ?> , <SubComponent1Conc unit0 ?>, <SubComponent2Name: ?> , <SubComponent2Conc unit0 ?>, <SubComponentNName: ?> , <SubComponentNConc unit0 ?>
//   <well1> , <component name1> , <component type1 ?> , <volume1 ?> , <volume unit1 ?>, <conc1 ?> , <conc unit1 ?>, <SubComponent1Name: ?> , <SubComponent1Conc unit0 ?>, <SubComponent2Name: ?> , <SubComponent2Conc unit0 ?>, <SubComponentNName: ?> , <SubComponentNConc unit0 ?>
//   ...
//
func ParsePlateCSV(ctx context.Context, inData io.Reader, validationOptions ...ValidationConfig) (*ParsePlateResult, error) {

	var addDefaultConfig bool

	for _, configOption := range validationOptions {
		if len(configOption.valid) > 0 {
			addDefaultConfig = true
		}
	}

	if !addDefaultConfig {
		validationOptions = append(validationOptions, DefaultValidationConfig())
	}

	return parsePlateCSVWithValidationConfig(ctx, inData, validationOptions...)
}

// parsePlateCSVWithValidationConfig parses a csv file into a plate.
//
// CSV plate format: (? denotes optional, whitespace for clarity)
//
//   <plate type ?> , <plate name ?>
//   <well0> , <component name0> , <component type0 ?> , <volume0 ?> , <volume unit0 ?>, <conc0 ?> , <conc unit0 ?>, <SubComponent1Name: ?> , <SubComponent1Conc unit0 ?>, <SubComponent2Name: ?> , <SubComponent2Conc unit0 ?>, <SubComponentNName: ?> , <SubComponentNConc unit0 ?>
//   <well1> , <component name1> , <component type1 ?> , <volume1 ?> , <volume unit1 ?>, <conc1 ?> , <conc unit1 ?>, <SubComponent1Name: ?> , <SubComponent1Conc unit0 ?>, <SubComponent2Name: ?> , <SubComponent2Conc unit0 ?>, <SubComponentNName: ?> , <SubComponentNConc unit0 ?>
//   ...
//
// TODO: refactor if/when Opt loses raw []byte and file as InputPlate options
func parsePlateCSVWithValidationConfig(ctx context.Context, inData io.Reader, vcOptions ...ValidationConfig) (*ParsePlateResult, error) {
	// Get returning "" if idx >= len(xs)
	get := func(xs []string, idx int) string {
		if len(xs) <= idx {
			return ""
		}
		return strings.TrimSpace(xs[idx])
	}

	// returns first index of value in xs which matches equal fold comparison for header.
	// if no value is found, -1 and error will be returned
	getIndexForColumnHeader := func(xs []string, header string) (int, error) {
		for i, str := range xs {
			if strings.EqualFold(strings.TrimSpace(str), strings.TrimSpace(header)) {
				return i, nil
			}
		}
		return -1, errors.Errorf("header %s not found in list %v", header, xs)
	}

	var validationConfigs []ValidationConfig
	var replaceConfigs []ValidationConfig
	var vc ValidationConfig

	for _, config := range vcOptions {
		if len(config.valid) > 0 {
			validationConfigs = append(validationConfigs, config)
		} else if len(config.replaceField) > 0 {
			replaceConfigs = append(replaceConfigs, config)
		}
	}

	if len(validationConfigs) == 1 {
		vc = validationConfigs[0]
	} else {
		return nil, fmt.Errorf("conflicting validation configs specified to ParsePlateCSV, please only set one: %v", validationConfigs)
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

	csvr := csvutil.NewTolerantReader(inData)
	csvr.FieldsPerRecord = -1

	rec, err := csvr.Read()
	if err != nil {
		return nil, err
	}

	if len(rec) == 0 {
		return nil, fmt.Errorf("empty file")
	}

	var plateType string

	if len(replaceConfigs) == 0 {
		plateType = rec[0]
	} else {
		for _, replacer := range replaceConfigs {
			if replacement, found := replacer.replaceField[plateTypeReplacementKey]; found {
				plateType = replacement
				break
			}
		}
	}

	plate, err := inventory.NewPlate(ctx, plateType)
	if err != nil {
		return nil, fmt.Errorf("cannot make plate %s: %s", plateType, err)
	}

	plate.PlateName = get(rec, 1)
	if len(plate.PlateName) == 0 {
		plate.PlateName = fmt.Sprint("input_plate_", plate.ID)
	}

	var warnings []string
	var fatalErrors []string

	var lookForSubComponents bool

	// expectation is that there are 7 static columns in the plate csv file,
	// after that, a Sub components column header will be looked for.
	const numberOfStaticCSVColumns = 7

	subComponentsStart, err := getIndexForColumnHeader(rec, wtype.SubComponentsHeader)

	if err == nil {
		lookForSubComponents = true

		if subComponentsStart < numberOfStaticCSVColumns {
			return nil, errors.Errorf("%s header cannot be specified in the first %d column headers. Found at position %d", wtype.SubComponentsHeader, numberOfStaticCSVColumns, subComponentsStart+1)
		}
	}

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
		ctypeField := unModifyTypeName(get(rec, 2))

		well, err := validWellCoord(wellField)
		if err != nil {
			warnings = append(warnings, fmt.Sprintf("line %d skipped: %s", lineNo, err))
			continue
		} else if err := validWell(well, plate); err != nil {
			fatalErrors = append(fatalErrors, fmt.Sprintf("line %d: %s", lineNo, err))
			continue
		}

		if err := validName(cname, vc); err != nil {
			warnings = append(warnings, fmt.Sprintf("line %d skipped: %s", lineNo, err))
			continue
		}

		ctype, err := wtype.LiquidTypeFromString(wtype.PolicyName(ctypeField))

		if err != nil {
			warnings = append(warnings, fmt.Sprintf("line %d: component type %q not found in default system types, using (%q); this may generate undesirable behaviour if this is not a custom type known by the system: %s", lineNo, ctypeField, ctype, err))
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

		if len(rec) > numberOfStaticCSVColumns && lookForSubComponents {
			for k := subComponentsStart; k < len(rec); k = k + 2 {
				subCompName := get(rec, k)
				// sub component names may contain a : which must be removed
				trimmedSubCompName := strings.TrimRight(subCompName, ":")
				if trimmedSubCompName != "" {
					if k+1 < len(rec) {
						subCompConc := get(rec, k+1)
						subCmp := wtype.NewLHComponent()
						subCmp.SetName(trimmedSubCompName)
						err := cmp.AddSubComponent(subCmp, wunit.NewConcentration(wunit.SplitValueAndUnit(subCompConc)))
						if err != nil {
							return nil, err
						}
					} else if len(subCompName) != 0 {
						return nil, fmt.Errorf("no concentration set on sub component %s for well %s", trimmedSubCompName, well.FormatA1())
					}
				}
			}
		}
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

	if len(fatalErrors) > 0 {
		return nil, fmt.Errorf(strings.Join(fatalErrors, "\n"))
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
func ParseInputPlateFile(ctx context.Context, filename string) (*wtype.Plate, error) {
	r, err := parsePlateFile(ctx, filename)
	if err != nil {
		return nil, err
	}
	for _, warning := range r.Warnings {
		log.Println(warning)
	}
	return r.Plate, nil
}
