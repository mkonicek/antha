package workflowtest

import (
	"fmt"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"reflect"
	"sort"
)

type MixOutputComparisonOptions struct {
	CompareWells      bool
	ComparePositions  bool
	ComparePlateTypes bool
	CompareVolumes    bool
	ComparePlateNames bool
}

func CompareEveryting() MixOutputComparisonOptions {
	return MixOutputComparisonOptions{
		CompareWells:      true,
		ComparePositions:  true,
		ComparePlateTypes: true,
		CompareVolumes:    true,
		ComparePlateNames: true,
	}
}

func ComparePlateTypesVolumes() MixOutputComparisonOptions {
	return MixOutputComparisonOptions{
		CompareVolumes:    true,
		ComparePlateTypes: true,
	}
}

func ComparePlateTypesNamesVolumes() MixOutputComparisonOptions {
	return MixOutputComparisonOptions{
		CompareVolumes:    true,
		ComparePlateTypes: true,
		ComparePlateNames: true,
	}
}

type componentInfo struct {
	Well      string
	Position  string
	PlateName string
	PlateType string
	Volume    string
}

type ComparisonResult struct {
	Errors []error
}

type outputMap map[string][]componentInfo

func CompareMixOutputs(want, got map[string]*wtype.LHPlate, opts MixOutputComparisonOptions) ComparisonResult {
	// do we have the same number of things coming out at the
	// same volumes or whatever

	outputMapWant := getOutputMap(want, opts)
	outputMapGot := getOutputMap(got, opts)

	return compareOutputMaps(outputMapWant, outputMapGot)
}

func compareOutputMaps(outputMapWant, outputMapGot outputMap) ComparisonResult {
	errors := make([]error, 0, 1)

	// establish the order of comparison
	getKeysOrdered := func(m outputMap) []string {
		keys := make([]string, 0, len(m))
		for k, _ := range m {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		return keys
	}

	// remove stuff we've seen
	filterOutSeen := func(k1, k2 []string) []string {
		ret := make([]string, 0, len(k2))

		seen := make(map[string]bool, len(k1))

		for _, k := range k1 {
			seen[k] = true
		}

		for _, k := range k2 {
			if !seen[k] {
				ret = append(ret, k)
				seen[k] = true
			}
		}
		return ret
	}

	wantKeys := getKeysOrdered(outputMapWant)
	gotKeys := getKeysOrdered(outputMapGot)

	// if we have extra outputs we should report this whatever

	if len(wantKeys) != len(gotKeys) {
		errors = append(errors, fmt.Errorf("Missing or extra output components: want %d got %d", len(wantKeys), len(gotKeys)))
	}

	for _, key := range wantKeys {
		_, ok := outputMapGot[key]

		if !ok {
			errors = append(errors, fmt.Errorf("Missing output components %s %v", key, outputMapWant[key]))
		} else {
			errs := compareComponentInfo(key, outputMapWant[key], outputMapGot[key])
			if errs != nil {
				errors = append(errors, errs...)
			}
		}
	}

	for _, key := range filterOutSeen(gotKeys, wantKeys) {
		errors = append(errors, fmt.Errorf("Extra output components: %s %v", key, outputMapGot[key]))
	}

	return ComparisonResult{Errors: errors}
}

func compareComponentInfo(cname string, cpWant, cpGot []componentInfo) []error {
	errs := make([]error, 0, len(cpWant))

	missing, extra := missingExtra(cpWant, cpGot)

	if len(missing) != 0 {
		errs = append(errs, fmt.Errorf("Component %s: missing %d outputs %v", cname, len(missing), missing))
	}

	if len(extra) != 0 {
		errs = append(errs, fmt.Errorf("Component %s: %d extra outputs %v", cname, len(extra), extra))
	}

	return errs
}

func match(cp componentInfo, cpArr []componentInfo) bool {
	for _, cp2 := range cpArr {
		if reflect.DeepEqual(cp, cp2) {
			return true
		}
	}
	return false
}

func missingExtra(cpWant, cpGot []componentInfo) (missing, extra []componentInfo) {
	for _, cp := range cpWant {
		if !match(cp, cpGot) {
			missing = append(missing, cp)
		}
	}

	for _, cp := range cpGot {
		if !match(cp, cpWant) {
			extra = append(extra, cp)
		}
	}

	return
}

func getOutputMap(res map[string]*wtype.LHPlate, opts MixOutputComparisonOptions) outputMap {
	outputMap := make(map[string][]componentInfo)

	for pos, plate := range res {
		for _, col := range plate.Cols {
			for _, w := range col {
				if !w.Empty() {
					arr := outputMap[w.WContents.CName]

					if arr == nil {
						arr = make([]componentInfo, 0, 1)
					}

					arr = append(arr, getComponentInfo(w, pos, plate.PlateName, opts))

					outputMap[w.WContents.CName] = arr
				}
			}
		}
	}

	return outputMap
}

/*

type MixOutputComparisonOptions struct {
	CompareWells      bool
	ComparePositions  bool
	ComparePlateTypes bool
	CompareVolumes    bool
	ComparePlateNames bool
}


*/

func getComponentInfo(wellIn *wtype.LHWell, positionIn, plateNameIn string, opts MixOutputComparisonOptions) componentInfo {
	var well, position, platename, platetype, volume string

	if opts.CompareWells {
		well = wellIn.Crds
	}

	if opts.ComparePositions {
		position = positionIn
	}

	if opts.ComparePlateTypes {
		platetype = wellIn.Platetype
	}

	if opts.CompareVolumes {
		volume = wellIn.CurrVolume().ToString()
	}

	if opts.ComparePlateNames {
		platename = plateNameIn
	}

	return componentInfo{
		Well:      well,
		Position:  position,
		PlateName: platename,
		PlateType: platetype,
		Volume:    volume,
	}
}
