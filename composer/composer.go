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
// - generating a suitable main.go from the workflow
type Composer struct {
	Logger   *logger.Logger
	logFH    *os.File
	Workflow *workflow.Workflow

	InDir         string
	OutDir        string
	Keep          bool
	Run           bool
	LinkedDrivers bool

	elementTypes map[workflow.ElementTypeName]*TranspilableElementType
	worklist     []*TranspilableElementType
}

func NewComposer(logger *logger.Logger, wf *workflow.Workflow, inDir, outDir string, keep, run, linkedDrivers bool) (*Composer, error) {
	if outDir == "" {
		if d, err := ioutil.TempDir("", fmt.Sprintf("antha-build-%s", wf.JobId)); err != nil {
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

	return &Composer{
		Logger:   logger,
		logFH:    logFH,
		Workflow: wf,

		InDir:         inDir,
		OutDir:        outDir,
		Keep:          keep,
		Run:           run,
		LinkedDrivers: linkedDrivers,

		elementTypes: make(map[workflow.ElementTypeName]*TranspilableElementType),
	}, nil
}

func (c *Composer) ComposeAndRun() error {
	if err := c.Workflow.Repositories.Clone(filepath.Join(c.OutDir, "src")); err != nil {
		return err
	} else if err := c.Transpile(); err != nil {
		return err
	} else if err := c.GenerateMain(); err != nil {
		return err
	} else if err := c.PrepareDrivers(); err != nil { // Must do this before SaveWorkflow!
		return err
	} else if err := c.SaveWorkflow(); err != nil {
		return err
	} else if err := c.CompileWorkflow(); err != nil {
		return err
	} else {
		return c.RunWorkflow()
	}
}

func (c *Composer) EnsureElementType(et *TranspilableElementType) {
	if _, found := c.elementTypes[et.Name()]; !found {
		c.elementTypes[et.Name()] = et
		c.worklist = append(c.worklist, et)
	}
}

func (c *Composer) Transpile() error {
	c.Logger.Log("progress", "transpiling Antha elements")
	for _, et := range c.Workflow.Elements.Types {
		c.EnsureElementType(NewTranspilableElementType(et))
	}
	for idx := 0; idx < len(c.worklist); idx++ {
		if err := c.worklist[idx].TranspileFromFS(c); err != nil {
			return err
		}
	}
	return nil
}

func (c *Composer) GenerateMain() error {
	path := filepath.Join(c.OutDir, "workflow", "main.go")
	c.Logger.Log("progress", "generating workflow main", "path", path)
	if fh, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0400); err != nil {
		return err
	} else {
		defer fh.Close()
		return newMainRenderer(c).render(fh)
	}
}

func (c *Composer) SaveWorkflow() error {
	return c.Workflow.WriteToFile(filepath.Join(c.OutDir, "workflow", "data", "workflow.json"), false)
}

func (c *Composer) CloseLogs() {
	if c.logFH != nil {
		c.Logger.SwapWriters(os.Stderr)
		if err := c.logFH.Sync(); err != nil {
			c.Logger.Log("msg", "Error when syncing log file handle", "error", err)
		}
		if err := c.logFH.Close(); err != nil {
			c.Logger.Log("msg", "Error when closing log file handle", "error", err)
		}
		c.logFH = nil
	}
}
