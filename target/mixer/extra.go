package mixer

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/antha/anthalib/wutil"
	"github.com/antha-lang/antha/microArch/factory"
	"github.com/antha-lang/antha/microArch/logger"
)

func parseInputPlateData(inData io.Reader) (*wtype.LHPlate, error) {
	csvr := csv.NewReader(inData)
	csvr.FieldsPerRecord = -1

	// first line must be one cell, just the plate type name
	// e.g.
	// pcrplate_skirted_with_riser,

	line, err := csvr.Read()

	if err != nil {
		return nil, err
	}

	platetype := line[0]
	p := factory.GetPlateByType(platetype)

	if p == nil {
		err = errors.New(fmt.Sprint("Plate type ", platetype, " is not known"))
		return nil, err
	}

	pn := fmt.Sprint("input_plate_", p.ID)
	if len(line) > 1 {
		if len(line[1]) > 0 {
			pn = line[1]
		}
	}

	p.PlateName = pn

	// now we know the plate, just need to make the relevant components

	lineNo := 2

	//for rec, err := range csvr.Read() {
	for {
		rec, err := csvr.Read()
		if err == io.EOF {
			break
		}
		lineNo += 1

		if err != nil {
			logger.Info(fmt.Sprint("parseInputPlate ERRROR here: ", err))
			continue
		}

		// sanctioned order for these lines:
		// well,component name,component type (string), volume, volume unit

		/// first minimum length is 3 - well, name, type

		if len(rec) < 3 {
			logger.Info(fmt.Sprint("parseInputPlate ERROR (line ", lineNo, "): minimum length is 3 fields (well, component name, component type in string form)"))
			continue
		}

		if len(rec[0]) == 0 {
			logger.Info(fmt.Sprint("parseInputPlate WARN (line ", lineNo, "): skipped - no well coords"))
			continue

		}

		well := wtype.MakeWellCoords(rec[0])

		if well.IsZero() {
			logger.Info(fmt.Sprint("parseInputPlate ERROR (line ", lineNo, "): well coord string not well-formatted ", rec[0]))
			continue
		}

		if well.X >= p.WellsX() || well.Y >= p.WellsY() {
			logger.Info(fmt.Sprint("parseInputPlate ERROR (line ", lineNo, "): well coord ", rec[0], " does not exist on plate type ", platetype))
			continue
		}

		// next the component name... this is just a string so we just check there are no + signs and
		// at least one alphanumeric

		cname := rec[1]

		if strings.ContainsAny(cname, "+") {
			logger.Info(fmt.Sprint("parseInputPlate ERROR (line ", lineNo, "): no plus signs in component names allowed ", cname))
			continue
		}

		if !strings.ContainsAny(cname, "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789") {
			logger.Info(fmt.Sprint("parseInputPlate ERROR (line ", lineNo, "): component names must contain at least one alphanumeric character"))
			continue
		}

		// component type

		ctype, _ := wtype.LiquidTypeFromString(rec[2])

		if ctype == wtype.LTWater && !strings.Contains(rec[2], "water") {
			// just warn, this is tolerable but possibly not desired
			logger.Info(fmt.Sprint("parseInputPlate ERROR (line ", lineNo, "): component type ", rec[2], " is not recognized... defaulting to water."))

		}

		// volume ... have to check behaviour but making this non-optional might be a good idea.
		// alternatively would be a good mechanism to tell the system where to put stuff without
		// saying how much.

		vol := 0.0

		if len(rec) > 3 {
			vol = wutil.ParseFloat(rec[3])
			if vol == 0.0 {
				logger.Info(fmt.Sprint("parseInputFile warning (line ", lineNo, "): no sensible float in ", rec[3], " defaulting to volume of 0.0"))
			}
		}

		vunit := "ul"
		volume := wunit.NewVolume(0.0, "ul")
		if len(rec) == 4 || len(strings.TrimSpace(rec[4])) == 0 {
			logger.Info("Volume specified without unit... defaulting to microlitres")
		} else {
			volume = wunit.NewVolume(vol, vunit)
		}

		// Concentration

		conc := 0.0

		if len(rec) > 5 {
			conc = wutil.ParseFloat(rec[5])
			if conc == 0.0 {
				logger.Info(fmt.Sprint("parseInputFile warning (line ", lineNo, "): no sensible float in ", rec[5], " defaulting to volume of 0.0"))
			}
		}

		cunit := "ng/ul"
		if len(rec) > 6 {
			cunit = rec[6]
		}
		concentration := wunit.NewConcentration(conc, cunit)
		if len(rec) == 6 || len(strings.TrimSpace(rec[6])) == 0 {
			logger.Info("Conc specified without unit... defaulting to ng/ul")
		} else {
			concentration = wunit.NewConcentration(conc, cunit)
		}

		// now we make the component and stick it in the well

		cmp := wtype.NewLHComponent()

		cmp.Vunit = vunit
		cmp.Vol = volume.RawValue()
		cmp.CName = cname
		cmp.Type = ctype

		cmp.Conc = concentration.RawValue()
		cmp.Cunit = concentration.Unit().PrefixedSymbol()
		p.WellAt(well).Add(cmp)
	}

	// all done!

	return p, nil

}

// one file per plate
// here's the format:
// csv, first line MUST be
// _some_plate_type_,
// other lines are then
// well,component name, component type (string for now I think),
func parseInputPlateFile(filename string) (*wtype.LHPlate, error) {
	f, err := os.Open(filename)

	if err != nil {
		return nil, err
	}

	defer f.Close()
	return parseInputPlateData(f)
}

func ParseInputPlateFile(filename string) (*wtype.LHPlate, error) {
	p, err := parseInputPlateFile(filename)
	return p, err
}
