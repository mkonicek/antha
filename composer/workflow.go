package composer

import (
	"encoding/json"
	"io"
	"sort"
)

type Workflow struct {
	JobId string `json:"JobId"`

	Processes      map[string]*Process   `json:"Processes"`
	Parameters     map[string]*Parameter `json:"Parameters"`
	Connections    []*Connection         `json:"Connections"`
	elementClasses []string
}

type Process struct {
	Component string          `json:"Component"`
	Metadata  json.RawMessage `json:"Metadata"`
}

type Parameter map[string]json.RawMessage

type Connection struct {
	Source ConnectionEnd `json:"Source"`
	Target ConnectionEnd `json:"Target"`
}

type ConnectionEnd struct {
	Process string `json:"Process"`
	Port    string `json:"Port"`
}

func (wf *Workflow) ElementClasses() []string {
	if wf.elementClasses == nil {
		seen := make(map[string]struct{})
		ec := make([]string, 0, len(wf.Processes))
		for _, proc := range wf.Processes {
			if _, found := seen[proc.Component]; !found {
				seen[proc.Component] = struct{}{}
				ec = append(ec, proc.Component)
			}
		}
		sort.Strings(ec)
		wf.elementClasses = ec
	}
	return wf.elementClasses
}

func WorkflowFromReader(r io.Reader) (*Workflow, error) {
	wf := &Workflow{}
	dec := json.NewDecoder(r)
	if err := dec.Decode(wf); err != nil {
		return nil, err
	} else {
		return wf, nil
	}
}
