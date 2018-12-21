package composer

import (
	"encoding/json"
	"io"
)

type Workflow struct {
	JobId string `json:"JobId"`

	Processes      map[string]*Process   `json:"Processes"`
	Parameters     map[string]*Parameter `json:"Parameters"`
	Connections    []*Connection         `json:"Connections"`
	elementClasses []*ElementSource
}

type Process struct {
	Source   *ElementSource  `json:"Source"`
	Metadata json.RawMessage `json:"Metadata"`
}

type ElementSource struct {
	RepoId RepoId `json:"RepoId"`
	Branch string `json:"Branch"`
	Commit string `json:"Commit"`
	Path   string `json:"Path"`
}

func (es *ElementSource) CommitOrBranch() string {
	if es.Commit != "" {
		return es.Commit
	} else {
		return es.Branch
	}
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

func (wf *Workflow) ElementClasses() []*ElementSource {
	if wf.elementClasses == nil {
		// use value, not pointer, so we do structural equality
		seen := make(map[ElementSource]struct{})
		ec := make([]*ElementSource, 0, len(wf.Processes))
		for _, proc := range wf.Processes {
			src := *proc.Source
			if _, found := seen[src]; !found {
				seen[src] = struct{}{}
				ec = append(ec, &src)
			}
		}
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
