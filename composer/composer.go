package composer

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

// The Config struct is used to deserialise the config file only.
type Config struct {
	// repos containing elements
	Repositories Repositories `json:"Repositories"`
	// the directory in which we build the result
	OutDir string
}

func ConfigFromReader(r io.Reader) (*Config, error) {
	c := &Config{}
	dec := json.NewDecoder(r)
	if err := dec.Decode(c); err != nil {
		return nil, err
	} else {
		if len(c.OutDir) == 0 {
			if c.OutDir, err = ioutil.TempDir("", "antha-composer"); err != nil {
				return nil, err
			}
		}
		if err := os.MkdirAll(c.OutDir, 0700); err != nil {
			return nil, err
		}
		return c, nil
	}
}

// The composer manages the whole operation:
// - taking a configuration and a workflow,
// - locating the source of elements
// - transpiling those elements and writing them out to the right place
// - generating a suitable main.go from the workflow
type Composer struct {
	Config   *Config
	Workflow *Workflow

	classes  map[string]*LocatedElement
	worklist []*LocatedElement
}

func NewComposer(cfg *Config, workflow *Workflow) *Composer {
	return &Composer{
		Config:   cfg,
		Workflow: workflow,

		classes: make(map[string]*LocatedElement),
	}
}

func (c *Composer) LocateElementClasses() error {
	for _, class := range c.Workflow.ElementClasses() {
		if repo, found := c.Config.Repositories[class.RepoId]; !found {
			return fmt.Errorf("Unable to find matching Repository for RepoId '%s'", class.RepoId)
		} else {
			le := NewLocatedElement(repo, class)
			if err := le.FetchFiles(); err != nil {
				return err
			} else {
				c.EnsureLocatedElement(le)
			}
		}
	}
	return nil
}

func (c *Composer) EnsureLocatedElement(le *LocatedElement) {
	if _, found := c.classes[le.ImportPath]; !found {
		c.classes[le.ImportPath] = le
		c.worklist = append(c.worklist, le)
	}
}

func (c *Composer) Transpile() error {
	for idx := 0; idx < len(c.worklist); idx++ {
		le := c.worklist[idx]
		if err := le.Transpile(c); err != nil {
			return err
		}
	}
	return nil
}

func (c *Composer) GenerateMain() error {
	mr := &mainRenderer{
		Composer: c,
		varMemo:  make(map[string]string),
	}
	if fh, err := os.OpenFile(filepath.Join(c.Config.OutDir, "main.go"), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600); err != nil {
		return err
	} else {
		defer fh.Close()
		return mr.render(fh)
	}
}
