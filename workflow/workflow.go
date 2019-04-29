package workflow

import (
	crand "crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"math/big"
	"os"
	"path"
	"strings"
	"unicode"

	"github.com/antha-lang/antha/utils"

	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/qri-io/jsonschema"
	git "gopkg.in/src-d/go-git.v4"
)

//go:generate go-bindata -o ./schemas.go -pkg workflow -prefix schemas/ ./schemas/

type Workflow struct {
	SchemaVersion SchemaVersion `json:"SchemaVersion"`
	// The WorkflowId is the unique Id of this workflow itself, and is
	// not modified by the event of simulation.
	WorkflowId BasicId `json:"WorkflowId,omitempty"`
	// The SimulationId is an Id created by the act of simulation. Thus
	// a workflow that is simulated twice will have the same WorkflowId
	// but different SimulationIds.
	SimulationId BasicId `json:"SimulationId,omitempty"`

	Meta Meta `json:"Meta,omitempty"`

	Repositories Repositories `json:"Repositories"`
	Elements     Elements     `json:"Elements"`

	Inventory Inventory `json:"Inventory,omitempty"`

	Config Config `json:"Config"`

	Testing Testing `json:"Testing,omitempty"`

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
			errs := make(utils.ErrorSlice, len(valErrs))
			for i, err := range valErrs {
				errs[i] = err
			}
			return nil, errs.Pack()
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
			Rest: make(map[string]json.RawMessage),
		},
		Repositories: make(Repositories),
		Elements: Elements{
			Instances: make(ElementInstances),
		},
		Config: EmptyConfig(),
	}
}

func (wf *Workflow) EnsureWorkflowId() error {
	if wf.WorkflowId == "" {
		if id, err := RandomBasicId(""); err != nil {
			return err
		} else {
			wf.WorkflowId = id
		}
	}
	return nil
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
	Name string                     `json:"Name,omitempty"`
	Rest map[string]json.RawMessage `json:"-"`
}

// See https://golang.org/ref/spec#identifier However, we allow the
// first rune to be a digit, because currently the call sites all
// ensure the result of this call is prefixed with some constant text.
func (m *Meta) NameAsGoIdentifier() string {
	res := []rune{}
	lastWasUnderscore := false
	for _, r := range m.Name {
		switch {
		case unicode.IsLetter(r), unicode.IsDigit(r):
			res = append(res, r)
			lastWasUnderscore = false
		default:
			if !lastWasUnderscore {
				res = append(res, '_')
				lastWasUnderscore = true
			}
		}
	}
	return string(res)
}

func (m *Meta) UnmarshalJSON(bs []byte) error {
	all := make(map[string]json.RawMessage)
	if err := json.Unmarshal(bs, &all); err != nil {
		return err
	}
	// belt and braces due to Go's json support being case insensitive
	for key, val := range all {
		if strings.ToLower(key) == "name" {
			if err := json.Unmarshal(val, &m.Name); err != nil {
				return err
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

func (m *Meta) Set(key string, val interface{}) error {
	if bs, err := json.Marshal(val); err != nil {
		return err
	} else {
		m.Rest[key] = json.RawMessage(bs)
		return nil
	}
}

type BasicId string

func RandomBasicId(prefix BasicId) (BasicId, error) {
	max := big.NewInt(0).SetUint64(math.MaxUint64)
	if suffix, err := crand.Int(crand.Reader, max); err != nil {
		return "", err
	} else if prefix == "" {
		return BasicId(suffix.Text(62)), nil
	} else {
		return BasicId(fmt.Sprintf("%s-%s", prefix, suffix.Text(62))), nil
	}
}

// RepositoryName is the domain-qualified name of a repository, as it would
// appear in a Go `import` directive. For example,
// `repos.antha.com/antha-ninja/elements-westeros`.
type RepositoryName string

type ElementInstanceName string
type ElementPath string
type ElementTypeName string
type ElementParameterName string

// Repositories is a map of Repository values, keyed by RepositoryName. The keys
// should be unique, and no key should be a prefix of another. For example, if
// the keys in this map are github.com/foo/bar and github.com/foo/bar/baz, that
// will trigger a validation error.
type Repositories map[RepositoryName]*Repository

// Repository is a local, checked-out clone of a Git repository.
//
// Directory is the absolute path to the repository clone on the local file
// system
// Branch (optional) is the name of the branch to use.
// Commit (optional) is the SHA1 hash of the commit to use.
//
// You can provide Branch or Commit, but not both.
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
