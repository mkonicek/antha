package param

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/antha-lang/antha/execute"
	"github.com/antha-lang/antha/workflow"
	"github.com/ghodss/yaml"
)

var (
	badFormat = errors.New("bad format")
)

func unmarshal(data []byte, v interface{}) error {
	return yaml.Unmarshal(data, v)
}

// Flatten []map[string]RawMessage to RawParams
func flattenSList(pdata []byte, param *execute.RawParams) error {
	var params []map[string]json.RawMessage
	if err := unmarshal(pdata, &params); err != nil {
		return err
	} else if len(params) == 0 {
		return badFormat
	}

	var firstKeys map[string]bool
	for idx, p := range params {
		keys := make(map[string]bool)
		for k := range p {
			keys[k] = true
		}

		// Check that keys are the same between values
		if len(keys) == 0 {
			return badFormat
		} else if len(firstKeys) == 0 {
			firstKeys = keys
		} else if len(keys) != len(firstKeys) {
			return badFormat
		} else {
			for k := range firstKeys {
				if !keys[k] {
					return badFormat
				}
			}
		}

		key := fmt.Sprintf("Input%d", idx)
		param.Parameters[key] = p
	}
	return nil
}

// Flatten []RawParams to RawParams
func flattenList(pdata []byte, param *execute.RawParams) error {
	var params []execute.RawParams
	if err := unmarshal(pdata, &params); err != nil {
		return err
	} else if len(params) == 0 {
		return badFormat
	}

	for idx, p := range params {
		if p.Config != nil {
			return fmt.Errorf("cannot extend workflow with config parameters")
		}
		if len(p.Parameters) == 0 {
			return badFormat
		}
		for k, v := range p.Parameters {
			key := fmt.Sprintf("%d%s", idx, k)
			param.Parameters[key] = v
		}
	}
	return nil
}

// Parse parameters and workflow. If there are multiple input parameters and
// workflow is just one element, modify workflow to take multiple parameters.
func TryExpand(wdata, pdata []byte) (*workflow.Desc, *execute.RawParams, error) {
	getFirstProcess := func(desc *workflow.Desc) *workflow.Process {
		for _, p := range desc.Processes {
			return &p
		}
		return nil
	}

	var desc *workflow.Desc
	var param *execute.RawParams

	if err := unmarshal(wdata, &desc); err != nil {
		return nil, nil, err
	}

	// Try to parse parameters as is
	var origErr error
	if err := unmarshal(pdata, &param); err == nil {
		if len(param.Parameters) > 0 {
			return desc, param, nil
		}
		origErr = fmt.Errorf("did not find any parameters")
	} else {
		origErr = err
	}

	// Update parameters
	param = &execute.RawParams{
		Parameters: make(map[string]map[string]json.RawMessage),
	}
	if len(desc.Connections) > 0 {
		// Cannot extend workflow with connections
	} else if len(desc.Processes) != 1 {
		// Can only expand workflows with one process
	} else if err := flattenList(pdata, param); err == nil && len(param.Parameters) > 0 {
		origErr = nil
	} else if err := flattenSList(pdata, param); err == nil && len(param.Parameters) > 0 {
		// ^ This should be last as most parameters will match
		origErr = nil
	} else {
		origErr = fmt.Errorf("exhausted methods for expanding workflow: %s", origErr)
	}
	if origErr != nil {
		return nil, nil, origErr
	}

	// Update workflow
	process := getFirstProcess(desc)
	desc = &workflow.Desc{
		Processes: make(map[string]workflow.Process),
	}
	for k := range param.Parameters {
		desc.Processes[k] = *process
	}

	return desc, param, nil
}
