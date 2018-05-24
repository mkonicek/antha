// platePreferences
package plate

import (
	"fmt"

	"sort"
	"strconv"
	"strings"

	"context"

	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/search"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/execute"
	"github.com/antha-lang/antha/inventory"
)

// PlateName is the name of a plate
type PlateName string

// String returns the PlateName as a string.
func (plate PlateName) String() string {
	return string(plate)
}

// WellLocations are a set of specific WellLocations which could refer to well coordinates on a plate in A1 format. e.g. A1, B1, C1
type WellLocations []string

// String returns the WellLocations as a []string.
func (wells WellLocations) String() []string {
	var stringsUsed []string
	for _, well := range wells {
		stringsUsed = append(stringsUsed, string(well))
	}
	return stringsUsed
}

// PlateSpecificWellLocations stores a map of WellLocations using the PlateName as key.
// A "default" key may also be used which may apply to all Plates.
type PlateSpecificWellLocations map[string]WellLocations

// PlatePreferences stores the Preferences for managing plates.
type MixPreferences struct {
	// The type of plate which new plates should be created
	PlateType *wtype.LHPlate
	// The prefix to use for naming plates.
	// after the initial plate is full, a new plate will be created with platename PlateNamePrefix+CurrentPlateNumber, this will be appended with the CurrentPlateNumber e.g. _2, _3 etc.
	PlateNamePrefix PlateName
	// specific plates used
	PlatesUsed []*wtype.LHPlate
	// The number of the current plate.
	// after the initial plate is full, a new plate will be created with platename PlateNamePrefix+CurrentPlateNumber, this will be appended with the CurrentPlateNumber e.g. _2, _3 etc.
	CurrentPlateNumber int
	// MixPreferences

	// a map of Wells to avoid using PlateName as key
	// or a "default" to apply to all plates.
	UsedWells PlateSpecificWellLocations
	// a map of preferred wells using PlateName as key
	// or a "default" to apply to all plates.
	PreferredWells PlateSpecificWellLocations
	// determines whether the next well is chosen based on iterating well position by row or by column.
	ByRow   bool
	context context.Context
	mixInto bool
}

// HasSamplesOn evaluates whether the plate has any samples on it.
func HasSamplesOn(plate *wtype.LHPlate) bool {
	if len(plate.AllNonEmptyWells()) > 0 {
		return true
	}
	return false
}

// NewMixPreferences creates a new MixPreferences object.
func NewMixPreferences(ctx context.Context, plateType *wtype.LHPlate, plateNamePrefix PlateName) (MixPreferences, error) {
	if _, err := inventory.NewPlate(ctx, plateType.Type); err != nil {
		return MixPreferences{}, err
	}

	if HasSamplesOn(plateType) {

		// If the plate already has solutions, use the name of the platetype
		// An error is returned if the names do not match.
		if plateNamePrefix != "" && plateNamePrefix.String() != plateType.Name() {
			return MixPreferences{}, fmt.Errorf("plates with samples on cannot have the plate name changed: plateNamePrefix (%s) must be left blank or set to same name as OutPlate name (%s) if mixing into a plate with samples already in.", plateNamePrefix, plateType.Name())
		}
		plateNamePrefix = PlateName(plateType.Name())

		return MixPreferences{context: ctx, PlateNamePrefix: plateNamePrefix, PlateType: plateType, PlatesUsed: []*wtype.LHPlate{plateType}, CurrentPlateNumber: 1, mixInto: true}, nil
	}

	return MixPreferences{context: ctx, PlateNamePrefix: plateNamePrefix, PlateType: plateType, CurrentPlateNumber: 1}, nil
}

// CurrentPlateName returns the name of the latest plate used.
func (platePreferences *MixPreferences) CurrentPlateName() PlateName {
	return PlateName(normalisePlateName(string(platePreferences.PlateNamePrefix), platePreferences.CurrentPlateNumber))
}

func normalisePlateName(name string, number int) string {
	if number > 1 {
		return name + strconv.Itoa(number)
	}
	return name
}

// CurrentPlate returns the latest plate used.
func (platePreferences *MixPreferences) CurrentPlate() *wtype.LHPlate {
	if len(platePreferences.PlatesUsed) == 0 {
		OutPlate, err := inventory.NewPlate(platePreferences.context, platePreferences.PlateType.Type)
		if err != nil {
			panic(err)
		}
		OutPlate.SetName(platePreferences.CurrentPlateName().String())
		platePreferences.PlatesUsed = append(platePreferences.PlatesUsed, OutPlate)
	}
	return platePreferences.PlatesUsed[len(platePreferences.PlatesUsed)-1]
}

// SetCurrentPlateNumber changes the Current Plate number to that specified.
// Counting starts at 1.
// If the specified number does not exist it will be created.
// An error will be returned if the plate corresponding to the plate number specified is full.
func (platePreferences *MixPreferences) SetCurrentPlateNumber(platenumber int) error {

	for _, plate := range platePreferences.PlatesUsed {
		if normalisePlateName(string(platePreferences.PlateNamePrefix), platenumber) == plate.Name() {
			_, err := search.NextFreeWell(plate, platePreferences.UsedWells[plate.Name()].String(), []string{}, platePreferences.ByRow)
			// plate full so return error
			if err != nil {
				return fmt.Errorf("cannot change to plate number %d: %s", err.Error())
			}
			// looks like there's space so adjust platenumber and return nil
			platePreferences.CurrentPlateNumber = platenumber
			return nil
		}
	}
	// if we get here the plate must not exist, so we'll make it
	platePreferences.CurrentPlateNumber = platenumber
	OutPlate, err := inventory.NewPlate(platePreferences.context, platePreferences.PlateType.Type)
	if err != nil {
		return err
	}
	platename := platePreferences.CurrentPlateName()
	OutPlate.SetName(string(platename))
	platePreferences.PlatesUsed = append(platePreferences.PlatesUsed, OutPlate)
	return nil

}

// SetAvoidWells allows setting of any Wells to avoid.
// Use this parameter to specify wells to avoid for each named out put plate or specify a "default" to apply to all plates.
// e.g.
// "AliquotPlate":
// "A1",
// "B1",
// "A12,
// "AliquotPlate2":
// "A1",
// "A2"
// "default":
// "A1",
// "A2",
// "A3"
// An error will be returned if an invalid plate name is specified.
func (platePreferences *MixPreferences) SetAvoidWells(avoidWells PlateSpecificWellLocations) error {

	if platePreferences.PlateNamePrefix == "" {
		return fmt.Errorf("no plate name prefix specified; this must be added to PlatePreferences before setting the UsedWells parameter.")
	}

	platename := platePreferences.PlateNamePrefix

	// Make a map of platenames to wells used in those plates
	WellsUsed := make(PlateSpecificWellLocations)

	if len(avoidWells) > 0 {
		if _, found := avoidWells[platename.String()]; !found {

			var badPlates []string

			for plateName := range avoidWells {
				if !strings.Contains(string(plateName), string(platename)) && plateName != "default" {
					badPlates = append(badPlates, string(plateName))
				}
			}
			sort.Strings(badPlates)

			if len(badPlates) > 0 {
				return fmt.Errorf(`some plate names (%v) specified in WellsUsed parameter do not match specified plate name (%s). Please either remove the plate from WellsUsed, specify "default" instead of plate name or ensure the plate name matches`, badPlates, platename)
			}
		}
	}

	// check if any default wells used are added.
	defaultWellsUsed := avoidWells["default"]

	for plateName, wells := range avoidWells {
		WellsUsed[plateName] = append(WellsUsed[plateName], wells...)

		for _, defaultWellUsed := range defaultWellsUsed {
			if !search.InStrings(wells.String(), string(defaultWellUsed)) {
				WellsUsed[plateName] = append(WellsUsed[plateName], defaultWellUsed)
			}
		}
	}

	platePreferences.UsedWells = WellsUsed

	return nil
}

// SetPreferredWells adds wells which will be preferred when selecting the next well to mix to.
// Use this parameter to specify preferential wells to use first for each named out put plate or specify a "default" to apply to all plates.
// If no values are set the next available free wells will be used.
// e.g.
// "AliquotPlate":
// "A1",
// "B1",
// "A12,
// "AliquotPlate2":
// "A1",
// "A2"
// "default":
// "A1",
// "A2",
// "A3"
func (platePreferences *MixPreferences) SetPreferredWells(preferredWells PlateSpecificWellLocations) error {
	if platePreferences.PlateNamePrefix == "" {
		return fmt.Errorf("no plate name prefix specified; this must be added to PlatePreferences before setting the PreferredWells parameter.")
	}
	platePreferences.PreferredWells = preferredWells
	return nil
}

// SetMixByRow sets the preference to use in the NextFreeWell method. If this is set to true, wells will be ranged through by row, "A1", "A2", "A3".
// default behaviour is to range through by column, e.g. "A1", "B1", "C1"
func (platePreferences *MixPreferences) SetMixByRow(byRow bool) {
	platePreferences.ByRow = byRow
}

// NextFreeWell checks for the next well which is empty in a plate platePreferences.
// The user can also specify wells to avoid, preffered wells and
// whether to search through the well positions by row. The default is by column.
// If the current plate is full, a new plate will be created and the next well returned for that plate.
// Should be used in conjunction with CurrentPlate() to ensure the correct pipetting behaviour.
func (platePreferences *MixPreferences) NextFreeWell() (string, error) {

	platename := platePreferences.CurrentPlateName()

	if len(platePreferences.PlatesUsed) == 0 {
		OutPlate, err := inventory.NewPlate(platePreferences.context, platePreferences.PlateType.Type)
		if err != nil {
			return "", err
		}
		OutPlate.SetName(string(platename))
		platePreferences.PlatesUsed = append(platePreferences.PlatesUsed, OutPlate)
	}

	var preferredWells []string

	if foundWells, found := platePreferences.PreferredWells[platename.String()]; found {
		preferredWells = foundWells
	} else if defaultWells, defaultfound := platePreferences.PreferredWells["default"]; defaultfound {
		preferredWells = defaultWells
	}

	well, err := search.NextFreeWell(platePreferences.PlatesUsed[len(platePreferences.PlatesUsed)-1], platePreferences.UsedWells[platename.String()].String(), preferredWells, platePreferences.ByRow)

	// if an error is returned the plate is full so make a new plate
	if err != nil {

		OutPlate, err := inventory.NewPlate(platePreferences.context, platePreferences.PlateType.Type)
		if err != nil {
			return "", err
		}
		platePreferences.CurrentPlateNumber++
		platename = platePreferences.CurrentPlateName()
		OutPlate.SetName(string(platename))
		platePreferences.PlatesUsed = append(platePreferences.PlatesUsed, OutPlate)

		wellsUsed := platePreferences.UsedWells[platename.String()]

		var preferredWells []string

		if foundWells, found := platePreferences.PreferredWells[platename.String()]; found {
			preferredWells = foundWells
		} else if defaultWells, defaultfound := platePreferences.PreferredWells["default"]; defaultfound {
			preferredWells = defaultWells
		}

		well, err = search.NextFreeWell(OutPlate, wellsUsed, preferredWells, platePreferences.ByRow)

		if err != nil {
			return well, err
		}

	}

	return well, nil
}

// AvoidThisWellOnCurrentPlate will add the specified well to the AvoidWells field for the current plate.
// This will be necessary if using MixNamed and attempting to keep track of wells used, rather than MixInto which will deal with a specific plate and therefore handle wells used as part of the LHPlate object.
func (platePreferences *MixPreferences) AvoidThisWellOnCurrentPlate(well string) {
	if platePreferences.UsedWells == nil || len(platePreferences.UsedWells) == 0 {
		// Make a map of platenames to wells used in those plates
		platePreferences.UsedWells = make(PlateSpecificWellLocations)
	}
	if !search.InStrings(platePreferences.UsedWells[string(platePreferences.CurrentPlateName())].String(), well) {
		platePreferences.UsedWells[string(platePreferences.CurrentPlateName())] = append(platePreferences.UsedWells[string(PlateName(platePreferences.CurrentPlateName()))], well)
	}
}

// Mix is a high level method to directly mix the provided samples into the next available position dictated by the PlatePreferences
func (platePreferences *MixPreferences) Mix(samples ...*wtype.LHComponent) (*wtype.LHComponent, error) {
	var mixedSolution *wtype.LHComponent
	well, err := platePreferences.NextFreeWell()

	// if an error is returned the plate is full so make a new plate
	if err != nil {
		return nil, err
	}

	if platePreferences.mixInto {
		mixedSolution = execute.MixInto(platePreferences.context, platePreferences.CurrentPlate(), well, samples...)
	} else {
		mixedSolution = execute.MixNamed(platePreferences.context, platePreferences.CurrentPlate().Type, well, platePreferences.CurrentPlateName().String(), samples...)
	}
	// add the well to AvoidWells for current plate
	platePreferences.AvoidThisWellOnCurrentPlate(well)

	return mixedSolution, nil
}
