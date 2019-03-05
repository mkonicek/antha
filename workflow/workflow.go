package workflow

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"path"
	"strings"

	"github.com/antha-lang/antha/antha/anthalib/wtype"
	git "gopkg.in/src-d/go-git.v4"
)

type Workflow struct {
	JobId JobId `json:"JobId"`

	Meta Meta `json:"Meta,omitempty"`

	Repositories Repositories `json:"Repositories"`
	Elements     Elements     `json:"Elements"`

	Inventory Inventory `json:"Inventory,omitempty"`

	Config Config `json:"Config"`

	typeNames map[ElementTypeName]*ElementType
}

func WorkflowFromReaders(rs ...io.ReadCloser) (*Workflow, error) {
	acc := &Workflow{}
	for _, r := range rs {
		defer r.Close()
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

type Meta struct {
	Name string                 `json:"Name,omitempty"`
	Rest map[string]interface{} `json:"-"`
}

func (m *Meta) UnmarshalJSON(bs []byte) error {
	all := make(map[string]interface{})
	if err := json.Unmarshal(bs, &all); err != nil {
		return err
	}
	// belt and braces due to Go's json support being case insensitive
	for key, val := range all {
		if strings.ToLower(key) == "name" {
			if valStr, ok := val.(string); ok {
				m.Name = valStr
			}
			delete(all, key)
		}
	}
	m.Rest = all
	return nil
}

func (m *Meta) MarshalJSON() ([]byte, error) {
	all := make(map[string]interface{}, len(m.Rest)+1)
	for key, val := range m.Rest {
		if strings.ToLower(key) != "name" {
			all[key] = val
		}
	}
	if m.Name != "" {
		all["Name"] = m.Name
	}
	if len(all) == 0 {
		all = nil
	}
	return json.Marshal(all)
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
	Branch    string `json:"Branch,omitempty"`
	Commit    string `json:"Commit,omitempty"`

	gitRepo *git.Repository
}

type Elements struct {
	Types                ElementTypes                `json:"Types,omitempty"`
	Instances            ElementInstances            `json:"Instances,omitempty"`
	InstancesConnections ElementInstancesConnections `json:"InstancesConnections,omitempty"`
}

type ElementTypes []*ElementType

type ElementType struct {
	RepositoryPrefix RepositoryPrefix `json:"RepositoryPrefix"`
	ElementPath      ElementPath      `json:"ElementPath"`
}

func (et ElementType) Name() ElementTypeName {
	return ElementTypeName(path.Base(string(et.ElementPath)))
}

func (et ElementType) ImportPath() string {
	return path.Join(string(et.RepositoryPrefix), string(et.ElementPath))
}

type ElementInstances map[ElementInstanceName]ElementInstance

type ElementInstance struct {
	ElementTypeName ElementTypeName     `json:"ElementTypeName"`
	Meta            json.RawMessage     `json:"Meta,omitempty"`
	Parameters      ElementParameterSet `json:"Parameters,omitempty"`
}

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
		tn := make(map[ElementTypeName]*ElementType, len(wf.Elements.Types))
		for _, et := range wf.Elements.Types {
			tn[et.Name()] = et
		}
		wf.typeNames = tn
	}
	return wf.typeNames
}

type Inventory struct {
	PlateTypes wtype.PlateTypes `json:"PlateTypes,omitempty"`
	/* Currently only PlateTypes can be set but it's clear how to extend this:
	Components Components `json:"Components"`
	TipBoxes   TipBoxes   `json:"TipBoxes"`
	TipWastes  TipWastes  `json:"TipWastes"`
	*/
}
