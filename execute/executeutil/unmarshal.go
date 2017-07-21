package executeutil

import (
	"errors"
	"github.com/antha-lang/antha/execute"
	"github.com/antha-lang/antha/workflow"
	"github.com/ghodss/yaml"
)

var (
	noElements       = errors.New("no elements found")
	noParameters     = errors.New("no parameters found")
	bundleWithParams = errors.New("cannot use bundle with parameters and workflows")
)

func unmarshal(data []byte, v interface{}) error {
	return yaml.Unmarshal(data, v)
}

type UnmarshalOpt struct {
	BundleData   []byte
	ParamsData   []byte
	WorkflowData []byte
}

type Bundle struct {
	workflow.Desc
	execute.RawParams
	workflowtest.TestOpt
}

// Parse parameters and workflow.
func Unmarshal(opt UnmarshalOpt) (*workflow.Desc, *execute.RawParams, *workflowtest.TestOpt, error) {
	if len(opt.BundleData) != 0 && (len(opt.ParamsData) != 0 || len(opt.WorkflowData) != 0) {
		return nil, nil, nil, bundleWithParams
	}

	var desc workflow.Desc
	var param execute.RawParams
	var bundle Bundle
	var expected workflowtest.TestOpt

	if len(opt.BundleData) != 0 {
		if err := unmarshal(opt.BundleData, &bundle); err != nil {
			return nil, nil, nil, err
		}
		desc.Connections = bundle.Connections
		desc.Processes = bundle.Processes
		param.Config = bundle.Config
		param.Parameters = bundle.Parameters
	} else {
		if err := unmarshal(opt.WorkflowData, &desc); err != nil {
			return nil, nil, nil, err
		}
		if err := unmarshal(opt.ParamsData, &param); err != nil {
			return nil, nil, nil, err
		}
	}

	if len(desc.Processes) == 0 {
		return nil, nil, nil, noElements
	} else if len(param.Parameters) == 0 {
		return nil, nil, nil, noParameters
	}

	return &desc, &param, &expected, nil
}
