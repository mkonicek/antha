package workflow

import (
	"encoding/json"
	"errors"
	"io"
	"os"
	"path"
	"strings"

	"github.com/antha-lang/antha/antha/anthalib/wtype"
	git "gopkg.in/src-d/go-git.v4"
)

type Workflow struct {
	SchemaVersion SchemaVersion `json:"SchemaVersion"`
	JobId         JobId         `json:"JobId"`

	Meta Meta `json:"Meta,omitempty"`

	Repositories Repositories `json:"Repositories"`
	Elements     Elements     `json:"Elements"`

	Inventory Inventory `json:"AnthaInventory"`

	TmpInventory Inventory2 `json:"Inventory,omitempty"`

	Config Config `json:"Config"`

	Testing *Testing `json:"Testing,omitempty"`

	typeNames map[ElementTypeName]*ElementType
}

func WorkflowFromReaders(rs ...io.ReadCloser) (*Workflow, error) {
	if len(rs) == 0 {
		return nil, errors.New("No workflow sources provided.")
	}
	acc := EmptyWorkflow()
	for _, r := range rs {
		defer r.Close()
		wf := &Workflow{} // we're never merging _into_ wf so it's safe to have nil maps here
		dec := json.NewDecoder(r)
		if err := dec.Decode(wf); err != nil {
			return nil, err
		}
		// convert tmpInventory
		// loop over array or PlateTypes2, make new one, append to map[PlateTypeName]
		// set acc.Inventory.PlateTypes = max[plateTypes]
		wf.Inventory.PlateTypes = make(wtype.PlateTypes)
		for _, p2 := range wf.TmpInventory.PlateTypes2 {
			thisPlateType := p2.ConvertToPlateType()
			wf.Inventory.PlateTypes[thisPlateType.Name] = &thisPlateType
		}

		// make sure we don't use the original anywhere
		wf.TmpInventory.PlateTypes2 = nil
		if err := acc.Merge(wf); err != nil {
			return nil, err
		}
	}

	return acc, nil
}

// returns a fresh but fully initialised Workflow. In particular, all
// directly accessible empty maps are non-nil.
func EmptyWorkflow() *Workflow {
	return &Workflow{
		SchemaVersion: CurrentSchemaVersion,
		Meta: Meta{
			Rest: make(map[string]interface{}),
		},
		Repositories: make(Repositories),
		Elements: Elements{
			Instances: make(ElementInstances),
		},
		Config: EmptyConfig(),

		Inventory: Inventory{
			PlateTypes: make(wtype.PlateTypes),
		},
		TmpInventory: Inventory2{
			PlateTypes2: wtype.PlateTypes2{},
		},
	}
}

func (wf *Workflow) WriteToFile(p string, pretty bool) error {
	if p == "" || p == "-" {
		return wf.ToWriter(os.Stdout, pretty)
	} else if fh, err := os.OpenFile(p, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0400); err != nil {
		return err
	} else {
		defer fh.Close()
		return wf.ToWriter(fh, pretty)
	}
}

func (wf *Workflow) ToWriter(w io.Writer, pretty bool) error {
	enc := json.NewEncoder(w)
	if pretty {
		enc.SetIndent("", "\t")
	}
	return enc.Encode(wf)
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
type RepositoryName string
type ElementInstanceName string
type ElementPath string
type ElementTypeName string
type ElementParameterName string

type Repositories map[RepositoryName]*Repository

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
	RepositoryName RepositoryName `json:"RepositoryName"`
	ElementPath    ElementPath    `json:"ElementPath"`
}

func (et ElementType) Name() ElementTypeName {
	return ElementTypeName(path.Base(string(et.ElementPath)))
}

func (et ElementType) ImportPath() string {
	return path.Join(string(et.RepositoryName), string(et.ElementPath))
}

type ElementInstances map[ElementInstanceName]*ElementInstance

type ElementInstance struct {
	ElementTypeName ElementTypeName     `json:"ElementTypeName"`
	Meta            json.RawMessage     `json:"Meta,omitempty"`
	Parameters      ElementParameterSet `json:"Parameters,omitempty"`

	hasConnections bool
	hasParameters  bool
}

func (ei ElementInstance) IsUsed() bool {
	return ei.hasConnections || ei.hasParameters
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

type Inventory2 struct {
	PlateTypes2 wtype.PlateTypes2 `json:"PlateTypes,omitempty"`
}

type Inventory struct {
	PlateTypes wtype.PlateTypes
	/* Currently only PlateTypes can be set but it's clear how to extend this:
	Components Components `json:"Components"`
	TipBoxes   TipBoxes   `json:"TipBoxes"`
	TipWastes  TipWastes  `json:"TipWastes"`
	*/
}

const (
	CurrentSchemaVersion SchemaVersion = "2.0"
)

type SchemaVersion string
