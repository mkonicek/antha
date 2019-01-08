package composer

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"sort"

	"github.com/antha-lang/antha/utils"
	git "gopkg.in/src-d/go-git.v4"
)

type Workflow struct {
	JobId JobId `json:"JobId"`

	Repositories                Repositories                `json:"Repositories"`
	ElementTypes                ElementTypes                `json:"ElementTypes"`
	ElementInstances            ElementInstances            `json:"ElementInstances"`
	ElementInstancesParameters  ElementInstancesParameters  `json:"ElementInstancesParameters"`
	ElementInstancesConnections ElementInstancesConnections `json:"ElementInstancesConnections"`

	typeNames map[ElementTypeName]*ElementType
}

func newWorkflow() *Workflow {
	return &Workflow{
		Repositories:                make(Repositories),
		ElementTypes:                make(ElementTypes, 0),
		ElementInstances:            make(ElementInstances),
		ElementInstancesParameters:  make(ElementInstancesParameters),
		ElementInstancesConnections: make(ElementInstancesConnections, 0),
	}
}

func WorkflowFromReaders(rs ...io.Reader) (*Workflow, error) {
	acc := newWorkflow()
	for _, r := range rs {
		wf := &Workflow{}
		dec := json.NewDecoder(r)
		if err := dec.Decode(wf); err != nil {
			return nil, err
		} else if err := acc.merge(wf); err != nil {
			return nil, err
		}
	}
	if err := acc.validate(); err != nil {
		return nil, err
	} else {
		return acc, nil
	}
}

func (wf *Workflow) WriteToFile(p string) error {
	if bs, err := json.Marshal(wf); err != nil {
		return err
	} else {
		return ioutil.WriteFile(p, bs, 0400)
	}
}

type JobId string
type RepositoryPrefix string
type ElementInstanceName string
type ElementPath string
type ElementTypeName string
type ElementParameterName string

type Repositories map[RepositoryPrefix]*Repository

type Repository struct {
	Directory string `json:"Directory"`
	Branch    string `json:"Branch"`
	Commit    string `json:"Commit"`

	gitRepo *git.Repository
}

type ElementTypes []*ElementType

type ElementType struct {
	RepositoryPrefix RepositoryPrefix `json:"RepositoryPrefix"`
	ElementPath      ElementPath      `json:"ElementPath"`

	files map[string][]byte
}

func (et ElementType) Name() ElementTypeName {
	return ElementTypeName(path.Base(string(et.ElementPath)))
}

func (et ElementType) ImportPath() string {
	return path.Join(string(et.RepositoryPrefix), string(et.ElementPath))
}

type ElementInstances map[ElementInstanceName]ElementInstance

type ElementInstance struct {
	ElementTypeName ElementTypeName `json:"ElementTypeName"`
	Metadata        json.RawMessage `json:"Metadata"`
}

type ElementInstancesParameters map[ElementInstanceName]ElementParameterSet

type ElementParameterSet map[ElementParameterName]json.RawMessage

type ElementInstancesConnections []ElementConnection

type ElementConnection struct {
	Source ElementSocket `json:"Source"`
	Target ElementSocket `json:"Target"`
}

type ElementSocket struct {
	ElementInstance ElementInstanceName  `json:"ElementInstance"`
	ParameterName   ElementParameterName `json:"ParameterName"`
}

func (wf *Workflow) TypeNames() map[ElementTypeName]*ElementType {
	if wf.typeNames == nil {
		tn := make(map[ElementTypeName]*ElementType, len(wf.ElementTypes))
		for _, et := range wf.ElementTypes {
			tn[et.Name()] = et
		}
		wf.typeNames = tn
	}
	return wf.typeNames
}

func (a *Workflow) merge(b *Workflow) error {
	if a == nil || b == nil {
		return errors.New("Attempt to merge nil workflow")
	}

	if a.JobId == "" {
		a.JobId = b.JobId
	} else if b.JobId != "" && a.JobId != b.JobId {
		return fmt.Errorf("Cannot merge with different JobIds: %v vs %v", a.JobId, b.JobId)
	}

	errs := utils.ErrorSlice{
		a.Repositories.merge(b.Repositories),
		a.ElementTypes.merge(b.ElementTypes),
		a.ElementInstances.merge(b.ElementInstances),
		a.ElementInstancesParameters.merge(b.ElementInstancesParameters),
		a.ElementInstancesConnections.merge(b.ElementInstancesConnections),
	}

	if err := errs.Nub(); err != nil {
		return err
	} else {
		return nil
	}
}

func (a Repositories) merge(b Repositories) error {
	// It's an error iff a and b contain the same prefix with different definitions.
	for prefix, repoB := range b {
		if repoA, found := a[prefix]; found && repoA != repoB {
			return fmt.Errorf("Cannot merge: repository with prefix '%v' redefined.", prefix)
		} else if !found {
			a[prefix] = repoB
		}
	}
	return nil
}

func (ets ElementTypes) sort() {
	sort.Slice(ets, func(i, j int) bool {
		return ets[i].lessThan(ets[j])
	})
}

func (a *ElementType) lessThan(b *ElementType) bool {
	return a.RepositoryPrefix < b.RepositoryPrefix ||
		(a.RepositoryPrefix == b.RepositoryPrefix && a.ElementPath < b.ElementPath)
}

func (a *ElementType) equals(b *ElementType) bool {
	return a.RepositoryPrefix == b.RepositoryPrefix && a.ElementPath == b.ElementPath
}

func (a *ElementTypes) merge(b ElementTypes) error {
	all := make(ElementTypes, 0, len(*a)+len(b))
	all = append(all, *a...)
	all = append(all, b...)
	all.sort()

	result := make(ElementTypes, 0, len(all))
	old := ElementType{}
	for _, cur := range all {
		if old.equals(cur) {
			continue
		} else {
			result = append(result, cur)
			old = *cur
		}
	}
	*a = result
	return nil
}

func (a ElementInstances) merge(b ElementInstances) error {
	// Element instances from different workflows must be entirely distinct
	for name, elemB := range b {
		if _, found := a[name]; found {
			return fmt.Errorf("Cannot merge: element instance '%v' exists in both workflows", name)
		} else {
			a[name] = elemB
		}
	}
	return nil
}

func (a ElementInstancesParameters) merge(b ElementInstancesParameters) error {
	// Just like element instances, these should be completely distinct
	for name, paramSetB := range b {
		if _, found := a[name]; found {
			return fmt.Errorf("Cannot merge: element parameters '%v' exists in both workflows", name)
		} else {
			a[name] = paramSetB
		}
	}
	return nil
}

func (conns ElementInstancesConnections) sort() {
	sort.Slice(conns, func(i, j int) bool {
		return conns[i].lessThan(conns[j])
	})
}

func (a ElementConnection) lessThan(b ElementConnection) bool {
	return a.Source.lessThan(b.Source) ||
		(!b.Source.lessThan(a.Source) && a.Target.lessThan(b.Target))
}

func (a ElementSocket) lessThan(b ElementSocket) bool {
	return a.ElementInstance < b.ElementInstance ||
		(a.ElementInstance == b.ElementInstance && a.ParameterName < b.ParameterName)
}

func (a *ElementInstancesConnections) merge(b ElementInstancesConnections) error {
	all := make(ElementInstancesConnections, 0, len(*a)+len(b))
	all = append(all, *a...)
	all = append(all, b...)

	result := make(ElementInstancesConnections, 0, len(all))
	old := ElementConnection{}
	for _, cur := range all {
		if old == cur { // structural equality
			return fmt.Errorf("Cannot merge: element connection '%v' exists in both workflows", cur)
		} else {
			result = append(result, cur)
			old = cur
		}
	}
	*a = result
	return nil
}

func (wf *Workflow) validate() error {
	if wf.JobId == "" {
		return errors.New("Workflow has empty JobId")
	} else if err := wf.Repositories.validate(); err != nil {
		return err
	} else if err := wf.ElementTypes.validate(wf); err != nil {
		return err
	} else if err := wf.ElementInstances.validate(wf); err != nil {
		return err
	} else if err := wf.ElementInstancesParameters.validate(wf); err != nil {
		return err
	} else if err := wf.ElementInstancesConnections.validate(wf); err != nil {
		return err
	} else {
		return nil
	}
}

func (rs Repositories) validate() error {
	if len(rs) == 0 {
		return errors.New("Workflow has no Repositories")
	} else {
		for _, repo := range rs {
			if err := repo.validate(); err != nil {
				return err
			}
		}
		return nil
	}
}

func (r Repository) validate() error {
	if info, err := os.Stat(filepath.FromSlash(r.Directory)); err != nil {
		return err
	} else if !info.Mode().IsDir() {
		return fmt.Errorf("Repository Directory is not a directory: '%s'", r.Directory)
	} else if bEmpty, cEmpty := r.Branch == "", r.Commit == ""; !bEmpty && !cEmpty {
		return fmt.Errorf("Repository cannot have both Branch and Commit specified. At most one. ('%s', '%s')", r.Branch, r.Commit)
	} else {
		return nil
	}
}

func (ets ElementTypes) validate(wf *Workflow) error {
	// we don't support import aliasing for elements. This means that
	// we require that every element type has a unique type name.
	namesToPath := make(map[ElementTypeName]ElementPath, len(ets))
	for _, et := range ets {
		if err := et.validate(wf); err != nil {
			return err
		} else if ep, found := namesToPath[et.Name()]; found {
			return fmt.Errorf("ElementType '%v' is ambiguous (ElementPaths '%v' and '%v')", et.Name(), et.ElementPath, ep)
		} else {
			namesToPath[et.Name()] = et.ElementPath
		}
	}
	return nil
}

func (et ElementType) validate(wf *Workflow) error {
	if _, found := wf.Repositories[et.RepositoryPrefix]; !found {
		return fmt.Errorf("ElementType uses unknown RepositoryPrefix: '%s'", et.RepositoryPrefix)
	} else {
		return nil
	}
}

func (eis ElementInstances) validate(wf *Workflow) error {
	for name, ei := range eis {
		if name == "" {
			return errors.New("ElementInstance cannot have an empty name")
		} else if err := ei.validate(wf); err != nil {
			return err
		}
	}
	return nil
}

func (ei ElementInstance) validate(wf *Workflow) error {
	if _, found := wf.TypeNames()[ei.ElementTypeName]; !found {
		return fmt.Errorf("ElementInstance with ElementTypeName '%v' is unknown", ei.ElementTypeName)
	} else {
		return nil
	}
}

func (eps ElementInstancesParameters) validate(wf *Workflow) error {
	for name, _ := range eps {
		if _, found := wf.ElementInstances[name]; !found {
			return fmt.Errorf("ElementInstancesParameters provided for unknown ElementInstance '%v'", name)
		}
	}
	return nil
}

func (conns ElementInstancesConnections) validate(wf *Workflow) error {
	for _, conn := range conns {
		if err := conn.Source.validate(wf); err != nil {
			return err
		} else if err := conn.Target.validate(wf); err != nil {
			return err
		}
	}
	return nil
}

func (soc ElementSocket) validate(wf *Workflow) error {
	if _, found := wf.ElementInstances[soc.ElementInstance]; !found {
		return fmt.Errorf("ElementConnection uses ElementInstance '%v' which does not exist.", soc.ElementInstance)
	} else {
		return nil
	}
}
