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
	Workflow *workflow.Workflow

	OutDir        string
	Keep          bool
	Run           bool
	LinkedDrivers bool

	elementTypes map[workflow.ElementTypeName]*TranspilableElementType
	worklist     []*TranspilableElementType
}

func NewComposer(logger *logger.Logger, wf *workflow.Workflow, outDir string, keep, run, linkedDrivers bool) (*Composer, error) {
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
	}

	if err := os.MkdirAll(filepath.Join(outDir, "workflow", "data"), 0700); err != nil {
		return nil, err
	}

	return &Composer{
		Logger:   logger,
		Workflow: wf,

		OutDir:        outDir,
		Keep:          keep,
		Run:           run,
		LinkedDrivers: linkedDrivers,

		elementTypes: make(map[workflow.ElementTypeName]*TranspilableElementType),
	}, nil
}

func (c *Composer) FindWorkflowElementTypes() error {
	for _, et := range c.Workflow.Elements.Types {
		if err := c.Workflow.Repositories.CloneRepository(et, filepath.Join(c.OutDir, "src")); err != nil {
			return err
		} else {
			c.EnsureElementType(NewTranspilableElementType(et))
		}
	}
	return nil
}

func (c *Composer) EnsureElementType(et *TranspilableElementType) {
	if _, found := c.elementTypes[et.Name()]; !found {
		c.elementTypes[et.Name()] = et
		c.worklist = append(c.worklist, et)
	}
}

func (c *Composer) Transpile() error {
	c.Logger.Log("progress", "transpiling Antha elements")
	for idx := 0; idx < len(c.worklist); idx++ {
		if err := c.worklist[idx].Transpile(c); err != nil {
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
	return c.Workflow.WriteToFile(filepath.Join(c.OutDir, "workflow", "data", "workflow.json"))
}
