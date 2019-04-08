package workflow

import (
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"unicode"

	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/qri-io/jsonschema"
	git "gopkg.in/src-d/go-git.v4"
)

//go:generate go-bindata -o ./schemas.go -pkg workflow -prefix schemas/ ./schemas/

type Workflow struct {
	SchemaVersion SchemaVersion `json:"SchemaVersion"`
	JobId         JobId         `json:"JobId"`

	Meta Meta `json:"Meta,omitempty"`

	Repositories Repositories `json:"Repositories"`
	Elements     Elements     `json:"Elements"`

	Inventory Inventory `json:"Inventory,omitempty"`

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

		schema := MustAsset("workflow.schema.json")
		rs := &jsonschema.RootSchema{}
		if err := json.Unmarshal(schema, rs); err != nil {
			panic(err)
		}

		workflowJSON, err := ioutil.ReadAll(r)
		if err != nil {
			return nil, err
		}

		// Note: ValidateBytes unmarshals the JSON data we send it. Then we
		// unmarshal it ourselves a few lines later. It would be nice to only
		// unmarshal once, but we unmarshal into a struct, and ValidateBytes
		// unmarshals into an interface{} (which renders a value of type
		// `map[string]interface{}`). The jsonschema package doesn't (currently)
		// know how to validate a struct type. So for now, we'll live with
		// double-unmarshaling.
		valErrs, err := rs.ValidateBytes(workflowJSON)
		if err != nil {
			// ValidateBytes got an unmarshalling error
			return nil, err
		}
		if len(valErrs) > 0 {
			// ValidateBytes got validation errors
			errStrings := make([]string, len(valErrs))
			for _, err := range valErrs {
				errStrings = append(errStrings, err.Error())
			}
			return nil, errors.New(strings.Join(errStrings, "; "))
		}

		if err := json.Unmarshal(workflowJSON, wf); err != nil {
			return nil, err
		} else if err := acc.Merge(wf); err != nil {
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

func (jid JobId) AsIdentifier() string {
	res := make([]rune, 0, len(jid))
	for _, r := range jid {
		switch {
		// see https://golang.org/ref/spec#identifier
		// However, we allow the first rune to be a digit
		case r == '_', unicode.IsLetter(r), unicode.IsDigit(r):
			res = append(res, r)
		case r == ' ', r == '-', r == '/':
			res = append(res, '_')
		}
	}
	return string(res)
}

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

type Inventory struct {
	PlateTypes wtype.PlateTypes `json:"PlateTypes,omitempty"`
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
