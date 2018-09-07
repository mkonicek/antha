package executeutil

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"

	"github.com/antha-lang/antha/execute"
	"github.com/antha-lang/antha/utils"
	"github.com/antha-lang/antha/workflow"
	"github.com/antha-lang/antha/workflowtest"
)

// A Bundle is a workflow with its inputs
type Bundle struct {
	workflow.Desc
	execute.RawParams
	workflowtest.TestOpt
}

// UnmarshalAll attempts to read and parse all the file paths
// supplied, categorising the content into Bundles, Workflows and
// Params. Resulting maps have the path as the key.
func UnmarshalAll(paths ...string) (map[string]*Bundle, map[string]*workflow.Desc, map[string]*execute.RawParams, utils.ErrorSlice) {
	var errs utils.ErrorSlice

	bundles := make(map[string]*Bundle, len(paths))
	params := make(map[string]*execute.RawParams, len(paths))
	workflows := make(map[string]*workflow.Desc, len(paths))

	for _, path := range paths {
		if bytes, err := ioutil.ReadFile(path); err != nil {
			errs = append(errs, fmt.Errorf("Error when reading file %s: %v", path, err))
		} else {
			bundle := &Bundle{}
			if err := json.Unmarshal(bytes, bundle); err != nil {
				errs = append(errs, fmt.Errorf("Error when parsing content of %s: %v", path, err))
			} else if bundle.Processes != nil && bundle.Parameters != nil { // it's a bundle
				bundles[path] = bundle
			} else if bundle.Processes != nil { // it's a workflow
				workflows[path] = &bundle.Desc
			} else if bundle.Parameters != nil { // it's a params
				params[path] = &bundle.RawParams
			} else { // shrug
				errs = append(errs, fmt.Errorf("Unable to identify content of %s", path))
			}
		}
	}
	return bundles, workflows, params, errs
}

// UnmarshalSingle can either be supplied with a single bundlePath, or
// with both a workflowPath and a paramsPath. This slightly unusual
// API exists because there are a few places where due to flags or
// other sources of input, it is simplest to provide all possible
// inputs here. The bundle returned will either be constructed from
// the bundlePath alone, or from both the workflowPath and
// paramsPath. Any other combination will error.
func UnmarshalSingle(bundlePath, workflowPath, paramsPath string) (*Bundle, error) {
	if bundlePath != "" && workflowPath == "" && paramsPath == "" {
		if bundles, _, _, err := UnmarshalAll(bundlePath); err != nil {
			return nil, err
		} else if len(bundles) != 1 {
			return nil, fmt.Errorf("Passed %s as a bundle file, but I don't think that is a bundle file, sorry.", bundlePath)
		} else {
			for _, b := range bundles {
				return b, nil
			}
			panic("Unreachable")
		}

	} else if bundlePath == "" && workflowPath != "" && paramsPath != "" {
		if _, workflows, params, err := UnmarshalAll(workflowPath, paramsPath); err != nil {
			return nil, err
		} else if len(workflows) != 1 {
			return nil, fmt.Errorf("Passed %s as a workflow file, but I don't think that is a workflow file, sorry.", workflowPath)
		} else if len(params) != 1 {
			return nil, fmt.Errorf("Passed %s as a params file, but I don't think that is a params file, sorry.", paramsPath)
		} else {
			b := &Bundle{}
			for _, workflow := range workflows {
				b.Desc = *workflow
			}
			for _, param := range params {
				b.RawParams = *param
			}
			return b, nil
		}

	} else {
		return nil, errors.New("Either bundle must be provided, or both parameters and workflow must be provided.")
	}
}
