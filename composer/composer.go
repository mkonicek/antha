package composer

import (
	"os"
	"path/filepath"
)

// The composer manages the whole operation:
// - taking a configuration and a workflow,
// - locating the source of elements
// - transpiling those elements and writing them out to the right place
// - generating a suitable main.go from the workflow
type Composer struct {
	OutDir   string
	Workflow *Workflow

	elementTypes map[ElementTypeName]*ElementType
	worklist     []*ElementType
}

func NewComposer(outDir string, workflow *Workflow) *Composer {
	return &Composer{
		OutDir:   outDir,
		Workflow: workflow,

		elementTypes: make(map[ElementTypeName]*ElementType),
	}
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
	for idx := 0; idx < len(c.worklist); idx++ {
		if err := c.worklist[idx].Transpile(c); err != nil {
			return err
		}
	}
	return nil
}

func (c *Composer) GenerateMain() error {
	mr := newMainRenderer(c)
	if fh, err := os.OpenFile(filepath.Join(c.OutDir, "main.go"), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600); err != nil {
		return err
	} else {
		defer fh.Close()
		return mr.render(fh)
	}
}
