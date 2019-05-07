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
	if err := utils.MkdirAll(filepath.Join(outDir, "workflow", "data")); err != nil {
		return nil, err
	}

	logFH, err := utils.CreateFile(filepath.Join(outDir, "logs.txt"), utils.ReadWrite)
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

func (cb *ComposerBase) generateRepositoryGoMods() error {
	for repoName := range cb.clonedRepositories {
		path := filepath.Join(cb.OutDir, "src", filepath.FromSlash(string(repoName)), "go.mod")
		if fh, err := utils.CreateFile(path, utils.ReadWrite); err != nil {
			return err
		} else {
			defer fh.Close()
			if err := renderRepositoryMod(fh, repoName); err != nil {
				return err
			}
		}
	}
	return nil
}

func (cb *ComposerBase) generateWorkflowGoMod() error {
	path := filepath.Join(cb.OutDir, "workflow", "go.mod")
	if fh, err := utils.CreateFile(path, utils.ReadWrite); err != nil {
		return err
	} else {
		defer fh.Close()
		return renderWorkflowMod(fh, cb.clonedRepositories)
	}
}

func (cb *ComposerBase) generateGoGenerate() error {
	path := filepath.Join(cb.OutDir, "workflow", "generate_assets.go")
	if fh, err := utils.CreateFile(path, utils.ReadWrite); err != nil {
		return err
	} else {
		defer fh.Close()
		return renderGoGenerate(fh)
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
	if wf.SimulationId != "" {
		return fmt.Errorf("Workflow has already been simulated (SimulationId %v); aborting", wf.SimulationId)
	}

	mc := &mainComposer{
		ComposerBase:  cb,
		Workflow:      wf,
		Keep:          keep,
		Run:           run,
		LinkedDrivers: linkedDrivers,
	}

	return utils.ErrorFuncs{
		func() error { return mc.cloneRepositories(wf) },
		func() error { return mc.generateRepositoryGoMods() },
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
	if err := mc.generateGoGenerate(); err != nil {
		return err
	} else if err := mc.generateWorkflowGoMod(); err != nil {
		return err
	} else if fh, err := utils.CreateFile(path, utils.ReadWrite); err != nil {
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

	Workflows     map[workflow.BasicId]*testWorkflow
	Keep          bool
	Run           bool
	LinkedDrivers bool
	CoverPath     string
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
		Workflows:     make(map[workflow.BasicId]*testWorkflow),
		Keep:          keep,
		Run:           run,
		LinkedDrivers: linkedDrivers,
	}
}

func (tc *testComposer) AddWorkflow(wf *workflow.Workflow, inDir string) error {
	if wf.SimulationId != "" {
		return fmt.Errorf("Workflow has already been simulated (SimulationId %v); aborting", wf.SimulationId)
	} else if _, found := tc.Workflows[wf.WorkflowId]; found {
		return fmt.Errorf("Workflow with Id %v already added. Workflow Ids must be unique", wf.WorkflowId)
	} else {
		tc.Workflows[wf.WorkflowId] = &testWorkflow{
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
	// Need to do this in several steps because things like prepareDrivers
	// we can't do until after workflow go mod
	for _, twf := range tc.Workflows {
		err := utils.ErrorFuncs{
			func() error { return tc.cloneRepositories(twf.workflow) },
			func() error { return tc.transpile(twf.workflow) },
			func() error { return twf.generateTest() },
		}.Run()
		if err != nil {
			return err
		}
	}

	err := utils.ErrorFuncs{
		tc.generateRepositoryGoMods,
		tc.generateWorkflowGoMod,
	}.Run()
	if err != nil {
		return err
	}

	for _, twf := range tc.Workflows {
		err := utils.ErrorFuncs{
			func() error { return tc.prepareDrivers(&twf.workflow.Config) },
			func() error { return twf.saveWorkflow() },
		}.Run()
		if err != nil {
			return err
		}
	}

	return utils.ErrorFuncs{
		tc.generateGoGenerate,
		tc.goGenerate,
		tc.goTest,
	}.Run()
}

func (twf *testWorkflow) generateTest() error {
	path := filepath.Join(twf.OutDir, "workflow", fmt.Sprintf("workflow%d_test.go", twf.index))
	twf.Logger.Log("progress", "generating workflow test", "path", path)
	if fh, err := utils.CreateFile(path, utils.ReadWrite); err != nil {
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
