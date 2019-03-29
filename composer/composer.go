package composer

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/antha-lang/antha/logger"
	"github.com/antha-lang/antha/utils"
	"github.com/antha-lang/antha/workflow"
)

// The composer manages the whole operation:
// - taking a configuration and a workflow,
// - locating the source of elements
// - transpiling those elements and writing them out to the right place
// - generating a suitable main workflow.go from the workflow
type ComposerBase struct {
	Logger *logger.Logger
	logFH  *os.File

	InDir  string
	OutDir string

	clonedRepositories workflow.Repositories
	elementTypes       map[workflow.ElementTypeName]*TranspilableElementType
	worklist           []*TranspilableElementType
}

func NewComposerBase(logger *logger.Logger, inDir, outDir string) (*ComposerBase, error) {
	if outDir == "" {
		if d, err := ioutil.TempDir("", "antha-composer"); err != nil {
			return nil, err
		} else {
			logger.Log("outdir", d)
			outDir = d
		}
	} else if entries, err := ioutil.ReadDir(outDir); err != nil && !os.IsNotExist(err) {
		return nil, err
	} else if len(entries) != 0 {
		return nil, fmt.Errorf("Provided outdir '%s' must be empty (or not exist)", outDir)
	} else if outDir, err = filepath.Abs(outDir); err != nil {
		return nil, err
	}
	// always need to do this:
	if err := os.MkdirAll(filepath.Join(outDir, "workflow", "data"), 0700); err != nil {
		return nil, err
	}

	logFH, err := os.OpenFile(filepath.Join(outDir, "logs.txt"), os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0400)
	if err != nil {
		return nil, err
	} else {
		logger.SwapWriters(logFH, os.Stderr)
	}

	if inDir != "" {
		if inDir, err = filepath.Abs(inDir); err != nil {
			return nil, err
		}
	}

	return &ComposerBase{
		Logger: logger,
		logFH:  logFH,

		InDir:  inDir,
		OutDir: outDir,

		clonedRepositories: make(workflow.Repositories),
		elementTypes:       make(map[workflow.ElementTypeName]*TranspilableElementType),
	}, nil
}

func (cb *ComposerBase) CloseLogs() {
	if cb.logFH != nil {
		cb.Logger.SwapWriters(os.Stderr)
		if err := cb.logFH.Sync(); err != nil {
			cb.Logger.Log("msg", "Error when syncing log file handle", "error", err)
		}
		if err := cb.logFH.Close(); err != nil {
			cb.Logger.Log("msg", "Error when closing log file handle", "error", err)
		}
		cb.logFH = nil
	}
}

func (cb *ComposerBase) cloneRepositories(wf *workflow.Workflow) error {
	if err := cb.clonedRepositories.Merge(wf.Repositories); err != nil {
		return err
	} else {
		return cb.clonedRepositories.Clone(filepath.Join(cb.OutDir, "src"))
	}
}

func (cb *ComposerBase) ensureElementType(et *TranspilableElementType) {
	if _, found := cb.elementTypes[et.Name()]; !found {
		cb.elementTypes[et.Name()] = et
		cb.worklist = append(cb.worklist, et)
	}
}

func (cb *ComposerBase) transpile(wf *workflow.Workflow) error {
	cb.Logger.Log("progress", "transpiling Antha elements")
	cb.worklist = cb.worklist[:0]
	for _, et := range wf.Elements.Types {
		cb.ensureElementType(NewTranspilableElementType(et))
	}
	for idx := 0; idx < len(cb.worklist); idx++ {
		if err := cb.worklist[idx].TranspileFromFS(cb, wf); err != nil {
			return err
		}
	}
	return nil
}

func (cb *ComposerBase) token(wf *workflow.Workflow, elem workflow.ElementInstanceName, param workflow.ElementParameterName) (string, error) {
	if elemInstance, found := wf.Elements.Instances[elem]; !found {
		return "", fmt.Errorf("No such element instance with name '%v'", elem)
	} else if elemType, found := cb.elementTypes[elemInstance.ElementTypeName]; !found {
		return "", fmt.Errorf("No such element type with name '%v' (element instance '%v')", elemInstance.ElementTypeName, elem)
	} else if elemType.Transpiler == nil {
		return "", fmt.Errorf("The element type '%v' does not appear to contain an Antha element", elemInstance.ElementTypeName)
	} else if tok, found := elemType.Transpiler.TokenByParamName[string(param)]; !found {
		return "", fmt.Errorf("The element type '%v' has no parameter named '%v' (element instance '%v')",
			elemInstance.ElementTypeName, param, elem)
	} else {
		return tok.String(), nil
	}
}

// management for writing a single workflow as main
type mainComposer struct {
	*ComposerBase

	Workflow      *workflow.Workflow
	Keep          bool
	Run           bool
	LinkedDrivers bool
}

func (cb *ComposerBase) ComposeMainAndRun(keep, run, linkedDrivers bool, wf *workflow.Workflow) error {
	mc := &mainComposer{
		ComposerBase:  cb,
		Workflow:      wf,
		Keep:          keep,
		Run:           run,
		LinkedDrivers: linkedDrivers,
	}

	return utils.ErrorFuncs{
		func() error { return mc.cloneRepositories(wf) },
		func() error { return mc.transpile(wf) },
		func() error { return mc.generateMain() },
		func() error { return mc.prepareDrivers(&wf.Config) }, // Must do this before SaveWorkflow!
		func() error { return mc.saveWorkflow() },
		func() error { return mc.compileWorkflow() },
		func() error { return mc.runWorkflow() },
	}.Run()
}

func (mc *mainComposer) generateMain() error {
	path := filepath.Join(mc.OutDir, "workflow", "main.go")
	mc.Logger.Log("progress", "generating workflow main", "path", path)
	if fh, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0400); err != nil {
		return err
	} else {
		defer fh.Close()
		return renderMain(fh, mc)
	}
}

func (mc *mainComposer) saveWorkflow() error {
	return mc.Workflow.WriteToFile(filepath.Join(mc.OutDir, "workflow", "data", "workflow.json"), false)
}

// management for writing workflows as tests
type testComposer struct {
	*ComposerBase

	Workflows     map[workflow.JobId]*testWorkflow
	Keep          bool
	Run           bool
	LinkedDrivers bool
}

type testWorkflow struct {
	*testComposer
	index    int
	workflow *workflow.Workflow
	inDir    string
}

func (cb *ComposerBase) NewTestsComposer(keep, run, linkedDrivers bool) *testComposer {
	return &testComposer{
		ComposerBase:  cb,
		Workflows:     make(map[workflow.JobId]*testWorkflow),
		Keep:          keep,
		Run:           run,
		LinkedDrivers: linkedDrivers,
	}
}

func (tc *testComposer) AddWorkflow(wf *workflow.Workflow, inDir string) error {
	if _, found := tc.Workflows[wf.JobId]; found {
		return fmt.Errorf("Workflow with JobId %v already added. JobIds must be unique", wf.JobId)
	} else {
		tc.Workflows[wf.JobId] = &testWorkflow{
			testComposer: tc,
			index:        len(tc.Workflows),
			workflow:     wf,
			inDir:        inDir,
		}
		return nil
	}
}

func (tc *testComposer) ComposeTestsAndRun() error {
	if len(tc.Workflows) == 0 {
		return nil
	}
	for _, twf := range tc.Workflows {
		efs := utils.ErrorFuncs{
			func() error { return tc.cloneRepositories(twf.workflow) },
			func() error { return tc.transpile(twf.workflow) },
			func() error { return twf.generateTest() },
			func() error { return tc.prepareDrivers(&twf.workflow.Config) },
			func() error { return twf.saveWorkflow() },
		}
		if err := efs.Run(); err != nil {
			return err
		}
	}

	return utils.ErrorFuncs{
		tc.goGenerate,
		tc.goTest,
	}.Run()
}

func (twf *testWorkflow) generateTest() error {
	path := filepath.Join(twf.OutDir, "workflow", fmt.Sprintf("workflow%d_test.go", twf.index))
	twf.Logger.Log("progress", "generating workflow test", "path", path)
	if fh, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0400); err != nil {
		return err
	} else {
		defer fh.Close()
		return renderTest(fh, twf)
	}
}

func (twf *testWorkflow) saveWorkflow() error {
	leaf := fmt.Sprintf("workflow%d.json", twf.index)
	return twf.workflow.WriteToFile(filepath.Join(twf.OutDir, "workflow", "data", leaf), false)
}
