package composer

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/antha-lang/antha/logger"
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

	elementTypes map[workflow.ElementTypeName]*TranspilableElementType
	worklist     []*TranspilableElementType
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

		elementTypes: make(map[workflow.ElementTypeName]*TranspilableElementType),
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

func (cb *ComposerBase) ensureElementType(et *TranspilableElementType) {
	if _, found := cb.elementTypes[et.Name()]; !found {
		cb.elementTypes[et.Name()] = et
		cb.worklist = append(cb.worklist, et)
	}
}

func (cb *ComposerBase) transpile(wf *workflow.Workflow) error {
	cb.Logger.Log("progress", "transpiling Antha elements")
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

type mainComposer struct {
	*ComposerBase

	Workflow      *workflow.Workflow
	Keep          bool
	Run           bool
	LinkedDrivers bool

	varCount uint64
	varMemo  map[workflow.ElementInstanceName]string
}

func (cb *ComposerBase) ComposeMainAndRun(wf *workflow.Workflow, keep, run, linkedDrivers bool) error {
	mc := &mainComposer{
		ComposerBase:  cb,
		Workflow:      wf,
		Keep:          keep,
		Run:           run,
		LinkedDrivers: linkedDrivers,
	}

	if err := wf.Repositories.Clone(filepath.Join(mc.OutDir, "src")); err != nil {
		return err
	} else if err := mc.transpile(wf); err != nil {
		return err
	} else if err := mc.generateMain(); err != nil {
		return err
	} else if err := mc.prepareDrivers(&wf.Config); err != nil { // Must do this before SaveWorkflow!
		return err
	} else if err := mc.saveWorkflow(); err != nil {
		return err
	} else if err := mc.compileWorkflow(); err != nil {
		return err
	} else {
		return mc.runWorkflow()
	}
}

func (mc *mainComposer) generateMain() error {
	path := filepath.Join(mc.OutDir, "workflow", "main.go")
	mc.Logger.Log("progress", "generating workflow main", "path", path)
	if fh, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0400); err != nil {
		return err
	} else {
		defer fh.Close()
		return newMainRenderer(mc).render(fh)
	}
}

func (mc *mainComposer) saveWorkflow() error {
	return mc.Workflow.WriteToFile(filepath.Join(mc.OutDir, "workflow", "data", "workflow.json"), false)
}
