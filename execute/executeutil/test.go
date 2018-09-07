package executeutil

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/antha-lang/antha/execute"
	"github.com/antha-lang/antha/workflow"
	"github.com/antha-lang/antha/workflowtest"
)

// A TestInput is a workflow, its iputs and expected output
type TestInput struct {
	Dir          string
	BundlePath   string
	ParamsPath   string
	WorkflowPath string
	Params       *execute.RawParams
	Workflow     *workflow.Desc
	Expected     workflowtest.TestOpt
}

// Paths returns the filepaths that this TestInput was created from
func (a TestInput) Paths() (rs []string) {
	for _, s := range []string{a.BundlePath, a.ParamsPath, a.WorkflowPath} {
		if s != "" {
			rs = append(rs, s)
		}
	}
	return
}

type byPath []*TestInput

func (a byPath) Len() int {
	return len(a)
}

func (a byPath) Less(i, j int) bool {
	x, y := a[i], a[j]
	if x.WorkflowPath == y.WorkflowPath {
		return x.ParamsPath < y.ParamsPath
	}
	return x.WorkflowPath < y.WorkflowPath
}

func (a byPath) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

// FindTestInputs finds any test inputs under basePath.
func FindTestInputs(basePath string) ([]*TestInput, error) {
	// we have to group files by directory so we can match up params
	// and workflows later on.
	filesByDir := make(map[string][]string)

	// 1. find all json files
	walk := func(p string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		} else if fi.IsDir() {
			return nil
		} else if pabs, err := filepath.Abs(p); err != nil {
			return err
		} else if filepath.Ext(pabs) == ".json" {
			dir := filepath.Dir(pabs)
			filesByDir[dir] = append(filesByDir[dir], pabs)
		}
		return nil
	}

	if len(basePath) == 0 {
		var err error
		basePath, err = os.Getwd()
		if err != nil {
			return nil, err
		}
	}

	if err := filepath.Walk(basePath, walk); err != nil {
		return nil, err
	}

	testInputs := make([]*TestInput, 0, len(filesByDir))

	for dir, files := range filesByDir {
		bundles, workflows, params, errs := UnmarshalAll(files...)
		// errors are not fatal in this case
		for path, bundle := range bundles {
			testInputs = append(testInputs, &TestInput{
				Dir:          dir,
				BundlePath:   path,
				ParamsPath:   path,
				WorkflowPath: path,
				Params:       &bundle.RawParams,
				Workflow:     &bundle.Desc,
				Expected:     bundle.TestOpt,
			})
		}

		if len(workflows) == 0 || len(params) == 0 {
			continue

		} else if len(workflows) == 1 { // we assume all the params are for this one workflow
			template := TestInput{Dir: dir}
			for path, workflow := range workflows {
				template.WorkflowPath = path
				template.Workflow = workflow
			}
			for path, param := range params {
				templateCopy := template
				templateCopy.ParamsPath = path
				templateCopy.Params = param
				testInputs = append(testInputs, &templateCopy)
			}

		} else if len(workflows) == len(params) { // assume they match 1-1 based on lexicographical sort
			workflowKeys := make([]string, 0, len(workflows))
			paramsKeys := make([]string, 0, len(params))
			for k := range workflows {
				workflowKeys = append(workflowKeys, k)
			}
			for k := range params {
				paramsKeys = append(paramsKeys, k)
			}
			sort.Strings(workflowKeys)
			sort.Strings(paramsKeys)

			for idx, workflowPath := range workflowKeys {
				paramPath := paramsKeys[idx]
				workflow := workflows[workflowPath]
				param := params[paramPath]

				testInputs = append(testInputs, &TestInput{
					Dir:          dir,
					ParamsPath:   paramPath,
					WorkflowPath: workflowPath,
					Params:       param,
					Workflow:     workflow,
				})
			}

		} else {
			errs = append(errs, fmt.Errorf("Confused by %s where we have %d workflows and %d params. Skipping dir.", dir, len(workflows), len(params)))
		}

		if len(errs) > 0 {
			fmt.Printf("Non fatal errors encountered when loading from %s:\n\t%v\n", dir, errs)
		}
	}

	sort.Sort(byPath(testInputs))

	return testInputs, nil
}
