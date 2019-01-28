package composer

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

// The composer manages the whole operation:
// - taking a configuration and a workflow,
// - locating the source of elements
// - transpiling those elements and writing them out to the right place
// - generating a suitable main.go from the workflow
type Composer struct {
	Logger   *Logger
	OutDir   string
	Workflow *Workflow

	elementTypes map[ElementTypeName]*ElementType
	worklist     []*ElementType
}

func NewComposer(logger *Logger, outDir string, wf *Workflow) (*Composer, error) {
	if outDir == "" {
		if d, err := ioutil.TempDir("", fmt.Sprintf("antha-%s", wf.JobId)); err != nil {
			return nil, err
		} else {
			logger.Log("msg", fmt.Sprintf("Using '%s' for output.", d))
			outDir = d
		}
	}

	if err := os.MkdirAll(filepath.Join(outDir, "workflow", "data"), 0700); err != nil {
		return nil, err
	}

	return &Composer{
		Logger:   logger,
		OutDir:   outDir,
		Workflow: wf,

		elementTypes: make(map[ElementTypeName]*ElementType),
	}, nil
}

func (c *Composer) FindWorkflowElementTypes() error {
	for _, et := range c.Workflow.ElementTypes {
		if _, err := c.Workflow.Repositories.FetchFiles(et); err != nil {
			return err
		} else {
			c.EnsureElementType(et)
		}
	}
	return nil
}

func (c *Composer) EnsureElementType(et *ElementType) {
	if _, found := c.elementTypes[et.Name()]; !found {
		c.elementTypes[et.Name()] = et
		c.worklist = append(c.worklist, et)
	}
}

func (c *Composer) Transpile() error {
	c.Logger.Log("progress", "transpiling")
	for idx := 0; idx < len(c.worklist); idx++ {
		if err := c.worklist[idx].Transpile(c); err != nil {
			return err
		}
	}
	return nil
}

func (c *Composer) GenerateMain() error {
	c.Logger.Log("progress", "generating workflow main")
	if fh, err := os.OpenFile(filepath.Join(c.OutDir, "workflow", "main.go"), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600); err != nil {
		return err
	} else {
		defer fh.Close()
		return newMainRenderer(c).render(fh)
	}
}

func (c *Composer) SaveWorkflow() error {
	return c.Workflow.WriteToFile(filepath.Join(c.OutDir, "workflow", "data", "workflow.json"))
}
