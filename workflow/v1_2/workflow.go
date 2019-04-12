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

	if success, js, _ := maybeMigrateFileFlat(fm, param); success {
		return js, nil
		/*} else if success, js, _ := maybeMigrateFileArray(fm, param); success {
			return js, nil
		} else if success, js, _ := maybeMigrateFileMap(fm, param); success {
			return js, nil
		*/
	} else {
		return param, nil
	}
}

// maybeMigrateFileMap returns bool indicating whether this could be serialised as a map of file references, then either the original json or transformed json if successful.
func maybeMigrateFileMap(fm *effects.FileManager, param json.RawMessage) (bool, json.RawMessage, error) {
	bifMap := make(map[string]*bakedInFile, 0)
	if err := json.Unmarshal([]byte(param), &bifMap); err != nil {
		return false, nil, nil
	}

	js := make(map[string]wtype.File, len(bifMap))
	for k, bif := range bifMap {
		if f, err := moveDataToFile(fm, bif); err != nil {
			return false, nil, err
		} else if f == nil {
			return false, nil, nil
		} else {
			js[k] = *f
		}
	}

	if msg, err := json.Marshal(js); err != nil {
		return false, nil, err
	} else {
		return true, msg, nil
	}
}

// maybeMigrateFileArray returns bool indicating whether this could be serialised as an array of file references, then either the original json or transformed json if successful.
func maybeMigrateFileArray(fm *effects.FileManager, param json.RawMessage) (bool, json.RawMessage, error) {
	bifArr := make([]*bakedInFile, 0)

	if err := json.Unmarshal([]byte(param), &bifArr); err != nil {
		return false, nil, nil
	}

	js := make([]wtype.File, len(bifArr))
	for i, bif := range bifArr {
		if f, err := moveDataToFile(fm, bif); err != nil {
			return false, nil, err
		} else if f == nil {
			return false, nil, nil
		} else {
			js[i] = *f
		}
	}

	if msg, err := json.Marshal(js); err != nil {
		return false, nil, err
	} else {
		return true, msg, nil
	}
}

// maybeMigrateFileFlat returns bool indicating whether this could be serialised as a file reference, then either the original json or transformed json if successful.
func maybeMigrateFileFlat(fm *effects.FileManager, param json.RawMessage) (bool, json.RawMessage, error) {
	bif := &bakedInFile{}
	if err := json.Unmarshal([]byte(param), &bif); err != nil {
		return false, param, nil
	} else if f, err := moveDataToFile(fm, bif); err != nil {
		return false, nil, err
	} else if f == nil {
		return false, nil, nil
	} else if msg, err := json.Marshal(f); err != nil {
		return false, nil, err
	} else {
		return true, msg, nil
	}
}

func moveDataToFile(fm *effects.FileManager, bif *bakedInFile) (*wtype.File, error) {
	if bif == nil {
		return nil, nil
	} else if bif.Name == "" || len(bif.Bytes.Bytes) == 0 {
		return nil, nil
	} else if f, err := fm.WriteAll(bif.Bytes.Bytes, bif.Name); err != nil {
		return nil, err
	} else if f == nil {
		return nil, nil
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
