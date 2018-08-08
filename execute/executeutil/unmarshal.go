package executeutil

import (
	"errors"
	"github.com/antha-lang/antha/execute"
	"github.com/antha-lang/antha/workflow"
	"github.com/antha-lang/antha/workflowtest"
	"github.com/ghodss/yaml"
)

var (
	errNoElements       = errors.New("no elements found")
	errNoParameters     = errors.New("no parameters found")
	errBundleWithParams = errors.New("cannot use bundle with parameters and workflows")
)

func unmarshal(data []byte, v interface{}) error {
	return yaml.Unmarshal(data, v)
}

// UnmarshalOpt are options for Unmarshal
type UnmarshalOpt struct {
	BundleData   []byte
	ParamsData   []byte
	WorkflowData []byte
}

// A Bundle is a workflow with its inputs
type Bundle struct {
	workflow.Desc
	execute.RawParams
	workflowtest.TestOpt
	Version    string             `json:"version"`
	Properties workflowProperties `json:"Properties"`
}

type workflowProperties struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// Unmarshal parses parameters and workflow.
func Unmarshal(opt UnmarshalOpt) (*Bundle, error) {
	if len(opt.BundleData) != 0 && (len(opt.ParamsData) != 0 || len(opt.WorkflowData) != 0) {
		return nil, errBundleWithParams
	}

	var desc workflow.Desc
	var param execute.RawParams
	var bundle Bundle
	var expected workflowtest.TestOpt
	var version string
	var properties workflowProperties

	if len(opt.BundleData) != 0 {
		if err := unmarshal(opt.BundleData, &bundle); err != nil {
			return nil, err
		}
		desc.Connections = bundle.Connections
		desc.Processes = bundle.Processes
		param.Config = bundle.Config
		param.Parameters = bundle.Parameters
		expected = bundle.TestOpt
		version = bundle.Version
		properties = bundle.Properties
	} else {
		if err := unmarshal(opt.WorkflowData, &desc); err != nil {
			return nil, err
		}
		if err := unmarshal(opt.ParamsData, &param); err != nil {
			return nil, err
		}
	}

	if len(desc.Processes) == 0 {
		return nil, errNoElements
	} else if len(param.Parameters) == 0 {
		return nil, errNoParameters
	}

	bdl := Bundle{desc, param, expected, version, properties}
	return &bdl, nil
}
