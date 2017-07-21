package executeutil

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/antha-lang/antha/execute"
	"github.com/antha-lang/antha/workflow"
)

// A test input
type TestInput struct {
	BundlePath   string
	ParamsPath   string
	Params       *execute.RawParams
	WorkflowPath string
	Workflow     *workflow.Desc
	Dir          string
}

func (a TestInput) Paths() (rs []string) {
	if len(a.BundlePath) != 0 {
		rs = append(rs, a.BundlePath)
	}
	if len(a.ParamsPath) != 0 {
		rs = append(rs, a.ParamsPath)
	}
	if len(a.WorkflowPath) != 0 {
		rs = append(rs, a.WorkflowPath)
	}
	return
}

type TestInputs []*TestInput

func (a TestInputs) Len() int {
	return len(a)
}

func (a TestInputs) Less(i, j int) bool {
	if a[i].WorkflowPath == a[j].WorkflowPath {
		return a[i].ParamsPath < a[j].ParamsPath
	}
	return a[i].WorkflowPath < a[j].WorkflowPath
}

func (a TestInputs) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

// Find any test inputs under basePath.
func FindTestInputs(basePath string) ([]*TestInput, error) {
	wfiles := make(map[string][]string)
	pfiles := make(map[string][]string)
	bfiles := make(map[string][]string)

	// Find candidate inputs (*{bundle,parameters,workflow}*.{json,yml,yaml})
	// from directory
	walk := func(p string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if fi.IsDir() {
			return nil
		}
		pabs, err := filepath.Abs(p)
		if err != nil {
			return err
		}

		dir := filepath.Dir(pabs)
		b := filepath.Base(pabs)

		switch filepath.Ext(b) {
		case ".json":
		case ".yml":
		case ".yaml":
		default:
			return nil
		}

		switch {
		case strings.Contains(b, "workflow"):
			wfiles[dir] = append(wfiles[dir], pabs)
		case strings.Contains(b, "param"):
			pfiles[dir] = append(pfiles[dir], pabs)
		case strings.Contains(b, "bundle"):
			bfiles[dir] = append(bfiles[dir], pabs)
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

	// Match workflow with parameter files
	var inputs []*TestInput
	for dir, wfs := range wfiles {
		pfs := pfiles[dir]
		switch nwfs, npfs := len(wfs), len(pfs); {
		case nwfs == 0 || npfs == 0:
			continue
		case nwfs == npfs:
			// Matching number of pairs in directory: pair them up lexicographically
			sort.Strings(wfs)
			sort.Strings(pfs)
			for idx := range wfs {
				inputs = append(inputs, &TestInput{
					WorkflowPath: wfs[idx],
					ParamsPath:   pfs[idx],
					Dir:          dir,
				})
			}
		case nwfs == 1:
			// If just one workflow, assume the parameters are for this workflow
			for idx := range pfs {
				inputs = append(inputs, &TestInput{
					WorkflowPath: wfs[0],
					ParamsPath:   pfs[idx],
					Dir:          dir,
				})
			}
		default:
			continue
		}
	}

	// Bundles
	for dir, bfs := range bfiles {
		for _, bf := range bfs {
			inputs = append(inputs, &TestInput{
				BundlePath: bf,
				Dir:        dir,
			})
		}
	}

	open := func(name string) ([]byte, error) {
		if len(name) == 0 {
			return nil, nil
		}
		return ioutil.ReadFile(name)
	}

	for _, input := range inputs {
		var wdata, pdata, bdata []byte
		var err error
		if bdata, err = open(input.BundlePath); err != nil {
			return nil, fmt.Errorf("error reading %q", input.BundlePath, err)
		}
		if wdata, err = open(input.WorkflowPath); err != nil {
			return nil, fmt.Errorf("error reading %q", input.WorkflowPath, err)
		}
		if pdata, err = open(input.ParamsPath); err != nil {
			return nil, fmt.Errorf("error reading %q", input.ParamsPath, err)
		}

		wdesc, params, err := Unmarshal(UnmarshalOpt{
			BundleData:   bdata,
			ParamsData:   pdata,
			WorkflowData: wdata,
		})
		if err != nil {
			return nil, fmt.Errorf("error parsing %q: %s", strings.Join(input.Paths(), ","), err)
		}
		input.Params = params
		input.Workflow = wdesc
	}

	sort.Sort(TestInputs(inputs))

	return inputs, nil
}
