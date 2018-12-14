package composer

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

// The Config struct is used to deserialise the config file only.
type Config struct {
	// repos containing elements
	ElementSources ElementSources `json:"ElementSources"`
	// the directory in which we build the result
	OutDir string
}

func ConfigFromReader(r io.Reader) (*Config, error) {
	c := &Config{}
	dec := json.NewDecoder(r)
	if err := dec.Decode(c); err != nil {
		return nil, err
	} else {
		c.ElementSources.Sort()
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

	classes map[string]*LocatedElement
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
		if locElem, err := c.Config.ElementSources.Match(class); err != nil {
			return err
		} else if locElem == nil {
			return fmt.Errorf("Unable to resolve element name: %s", class)
		} else if err := locElem.FetchFiles(); err != nil {
			return err
		} else {
			c.classes[class] = locElem
		}
	}
	return nil
}

func (c *Composer) GenerateMain(w io.Writer) error {
	mr := &mainRenderer{
		Composer: c,
		varMemo:  make(map[string]string),
	}
	return mr.render(w)
}
