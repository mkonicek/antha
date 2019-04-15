package v1_2

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/laboratory/effects"
	"github.com/antha-lang/antha/microArch/driver/liquidhandling"
	"github.com/antha-lang/antha/workflow"
)

type workflowv1_2 struct {
	desc
	testOpt
	rawParams
	Version    string             `json:"version"`
	Properties workflowProperties `json:"Properties"`
}

type workflowProperties struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type rawParams struct {
	Parameters map[string]map[string]json.RawMessage `json:"Parameters"`
	Config     *opt                                  `json:"Config"`
}

type desc struct {
	Processes   map[string]process `json:"Processes"`
	Connections []connection       `json:"connections"`
}

type connection struct {
	Src port `json:"source"`
	Tgt port `json:"target"`
}

type process struct {
	Component string         `json:"component"`
	Metadata  screenPosition `json:"metadata"`
}

type screenPosition struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type port struct {
	Process string `json:"process"`
	Port    string `json:"port"`
}

type testOpt struct {
	ComparisonOptions   string
	CompareInstructions bool
	CompareOutputs      bool
	Results             testResults
}

type testResults struct {
	MixTaskResults []mixTaskResult
}

type mixTaskResult struct {
	Instructions liquidhandling.SetOfRobotInstructions
	Outputs      map[string]*wtype.Plate
	TimeEstimate time.Duration
}

type bakedInFile struct {
	Name  string `json:"name"`
	Bytes struct {
		Bytes []byte `json:"bytes"`
	} `json:"bytes"`
}

// maybeMigrateFileParam returns either the original json representations, or a json representation as serialized files, if possible.
func maybeMigrateFileParam(fm *effects.FileManager, param json.RawMessage) (json.RawMessage, error) {

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

func (wf *workflowv1_2) MigrateElementParameters(fm *effects.FileManager, name string) (workflow.ElementParameterSet, error) {
	v, pr := wf.Parameters[name]
	if !pr {
		return nil, errors.New("parameters not present for element" + name)
	}
	pset := make(workflow.ElementParameterSet)

	for pname, pval := range v {
		if param, err := maybeMigrateFileParam(fm, pval); err != nil {
			return nil, err
		} else {
			pset[workflow.ElementParameterName(pname)] = param
		}
	}
	return pset, nil
}

func (wf *workflowv1_2) MigrateElement(fm *effects.FileManager, name string) (*workflow.ElementInstance, error) {
	ei := &workflow.ElementInstance{}

	v, pr := wf.Processes[name]
	if !pr {
		return nil, errors.New("element instance " + name + " not present")
	}

	ei.ElementTypeName = workflow.ElementTypeName(v.Component)
	enc, err := json.Marshal(v.Metadata)
	if err != nil {
		return nil, err
	}
	ei.Meta = json.RawMessage(enc)

	params, err := wf.MigrateElementParameters(fm, name)
	if err != nil {
		return nil, err
	}
	ei.Parameters = params
	return ei, nil
}
