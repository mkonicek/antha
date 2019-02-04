package workflow

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"path"

	"github.com/antha-lang/antha/antha/anthalib/wtype"
	git "gopkg.in/src-d/go-git.v4"
)

type Workflow struct {
	JobId JobId `json:"JobId"`

	Repositories                Repositories                `json:"Repositories"`
	ElementTypes                ElementTypes                `json:"ElementTypes"`
	ElementInstances            ElementInstances            `json:"ElementInstances"`
	ElementInstancesParameters  ElementInstancesParameters  `json:"ElementInstancesParameters"`
	ElementInstancesConnections ElementInstancesConnections `json:"ElementInstancesConnections"`

	Inventory Inventory `json:Inventory`

	Config Config `json:Config`

	typeNames map[ElementTypeName]*ElementType
}

func WorkflowFromReaders(rs ...io.Reader) (*Workflow, error) {
	acc := &Workflow{}
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

type Inventory struct {
	PlateTypes wtype.PlateTypes `json:"PlateTypes"`
	/* Currently only PlateTypes can be set but it's clear how to extend this:
	Components Components `json:"Components"`
	TipBoxes   TipBoxes   `json:"TipBoxes"`
	TipWastes  TipWastes  `json:"TipWastes"`
	*/
}
