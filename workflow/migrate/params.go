package migrate

import (
	"encoding/json"

	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/laboratory/effects"
)

// Utility code for migrating parameter values from raw JSON to workflow 2.0

type bakedInFile struct {
	Name  string `json:"name"`
	Bytes struct {
		Bytes []byte `json:"bytes"`
	} `json:"bytes"`
}

// MaybeMigrateFileParam returns either the original json representations, or a json representation as serialized files, if possible.
func MaybeMigrateFileParam(fm *effects.FileManager, param json.RawMessage) (json.RawMessage, error) {

	if js, err := maybeMigrateFileFlat(fm, param); js != nil {
		return js, nil
	} else if err != nil {
		return nil, err
	} else if js, err := maybeMigrateFileArray(fm, param); js != nil {
		return js, nil
	} else if err != nil {
		return nil, err
	} else {
		return param, nil
	}
}

// maybeMigrateFileArray returns a nil json.RawMessage unless: no error was encountered during unmarshalling and every value within the slice unmarshals as a valid wtype.File. returns a non-nil error only if an error was encountered when writing a file value to disk.
func maybeMigrateFileArray(fm *effects.FileManager, param json.RawMessage) (json.RawMessage, error) {
	var bifArr []*bakedInFile
	if err := json.Unmarshal([]byte(param), &bifArr); err != nil || bifArr == nil {
		return nil, nil
	}

	js := make([]*wtype.File, len(bifArr))
	for i, bif := range bifArr {
		if !bif.hasData() {
			return nil, nil
		} else if f, err := bif.moveDataToFile(fm); err != nil {
			return nil, err
		} else {
			js[i] = f
		}
	}

	return json.Marshal(js)
}

// maybeMigrateFileFlat returns a nil json.RawMessage unless: no error was encountered during unmarshalling as a valid wtype.File. returns a non-nil error only if an error was encountered when writing a file value to disk.
func maybeMigrateFileFlat(fm *effects.FileManager, param json.RawMessage) (json.RawMessage, error) {
	var bif *bakedInFile
	if err := json.Unmarshal([]byte(param), &bif); err != nil || !bif.hasData() {
		return nil, nil
	} else if f, err := bif.moveDataToFile(fm); err != nil {
		return nil, err
	} else {
		return json.Marshal(f)
	}
}

// hasData checks that structure is a valid bakedInFile (rather than an empty struct)
func (bif *bakedInFile) hasData() bool {
	return bif != nil && bif.Name != "" && len(bif.Bytes.Bytes) > 0
}

func (bif *bakedInFile) moveDataToFile(fm *effects.FileManager) (*wtype.File, error) {
	if f, err := fm.WriteAll(bif.Bytes.Bytes, bif.Name); err != nil {
		return nil, err
	} else {
		return f.AsInput(), nil
	}
}
