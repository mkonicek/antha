package workflowtest

import (
	"fmt"
	"reflect"
	"sort"

	"github.com/antha-lang/antha/antha/anthalib/wtype"
)

// A ComparisonMode is an option for comparing outputs
type ComparisonMode int

// Possible comparison modes
const (
	CompareWells ComparisonMode = 1 << iota
	ComparePositions
	ComparePlateTypes
	ComparePlateNames
	CompareVolumes
)

// Predefined comparison modes
const (
	CompareEverything             = CompareWells | ComparePositions | ComparePlateTypes | CompareVolumes | ComparePlateNames
	ComparePlateTypesVolumes      = ComparePlateTypes | CompareVolumes
	ComparePlateTypesNamesVolumes = CompareVolumes | ComparePlateTypes | ComparePlateNames
)

type componentInfo struct {
	Well      string
	Position  string
	PlateName string
	PlateType string
	Volume    string
	InstNo    string
}

func (ci componentInfo) HashKey() string {
	stringOrDefault := func(s1, s2 string) string {
		if s1 == "" {
			return s2
		}

		return s1
	}

	s := ""
	s += stringOrDefault(ci.Well, "WEL") + ":"
	s += stringOrDefault(ci.Position, "POS") + ":"
	s += stringOrDefault(ci.PlateName, "NAM") + ":"
	s += stringOrDefault(ci.PlateType, "TYP") + ":"
	s += stringOrDefault(ci.Volume, "VOL")

	return s
}

// A ComparisonResult is the output of comparing outputs
type ComparisonResult struct {
	Errors []error
}

type outputMap map[string][]componentInfo

// CompareMixOutputs compares mix outputs
func CompareMixOutputs(want, got map[string]*wtype.Plate, opts ComparisonMode) ComparisonResult {
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
		for k := range m {
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

	for _, key := range filterOutSeen(wantKeys, gotKeys) {
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

func getOutputMap(res map[string]*wtype.Plate, opts ComparisonMode) outputMap {
	outputMap := make(map[string][]componentInfo)
	instM := make(map[string]int)

	for pos, plate := range res {
		for _, col := range plate.Cols {
			for _, w := range col {
				if !w.IsEmpty() {
					arr := outputMap[w.WContents.CName]

					if arr == nil {
						arr = make([]componentInfo, 0, 1)
					}

					ci := getComponentInfo(w, pos, plate.PlateName, opts)
					arr = append(arr, ci)

					instNo, ok := instM[ci.HashKey()]

					if !ok {
						instM[ci.HashKey()] = 0
					}

					instM[ci.HashKey()]++

					ci.InstNo = fmt.Sprintf("instance_%d", instNo)

					outputMap[w.WContents.CName] = arr
				}
			}
		}
	}

	return outputMap
}

func getComponentInfo(wellIn *wtype.LHWell, positionIn, plateNameIn string, opts ComparisonMode) componentInfo {
	var well, position, platename, platetype, volume string

	if opts&CompareWells != 0 {
		well = wellIn.Crds.FormatA1()
	}

	if opts&ComparePositions != 0 {
		position = positionIn
	}

	if opts&ComparePlateTypes != 0 {
		platetype = wtype.TypeOf(wellIn.Plate)
	}

	if opts&CompareVolumes != 0 {
		if volumeInUl, err := wellIn.CurrentVolume().InStringUnit("ul"); err != nil {
			panic(err) //this should never happen with volumes as all units are compatible
		} else {
			volume = volumeInUl.ToString()
		}
	}

	if opts&ComparePlateNames != 0 {
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
