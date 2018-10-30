package execute

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"

	"github.com/antha-lang/antha/antha/anthalib/wtype"
	api "github.com/antha-lang/antha/api/v1"
	"github.com/antha-lang/antha/inject"
	"github.com/antha-lang/antha/inventory"
	"github.com/antha-lang/antha/meta"
	"github.com/antha-lang/antha/target/mixer"
	"github.com/antha-lang/antha/workflow"
	"github.com/golang/protobuf/jsonpb"
)

var (
	errUnknownParam = errors.New("unknown parameter")
)

// RawParams is the structure of parameter data for unmarshalling.
//
// Deprecated for github.com/antha-lang/antha/api/v1/WorkflowParameters.
type RawParams struct {
	Parameters map[string]map[string]json.RawMessage `json:"Parameters"`
	Config     *mixer.Opt                            `json:"Config"`
}

// Params is the structure of parameter data for marshalling.
//
// Deprecated for github.com/antha-lang/antha/api/v1/WorkflowParameters.
type Params struct {
	Parameters map[string]map[string]interface{} `json:"Parameters"`
	Config     *mixer.Opt                        `json:"Config"`
}

func tryString(data []byte) string {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		s = ""
	}
	return s
}

type unmarshaler struct {
	ReadLocalFiles bool
}

func (a *unmarshaler) unmarshalLHTipbox(ctx context.Context, data []byte, obj *wtype.LHTipbox) error {
	s := tryString(data)
	if len(s) == 0 {
		return json.Unmarshal(data, obj)
	}
	t, err := inventory.NewTipbox(ctx, s)
	if err != nil {
		return err
	}
	*obj = *t
	return nil
}

func (a *unmarshaler) unmarshalLHPlate(ctx context.Context, data []byte, obj *wtype.Plate) error {
	s := tryString(data)
	if len(s) == 0 {
		return json.Unmarshal(data, obj)
	}
	t, err := inventory.NewPlate(ctx, s)
	if err != nil {
		return err
	}
	*obj = *t
	return nil
}

func (a *unmarshaler) unmarshalLHComponent(ctx context.Context, data []byte, obj *wtype.Liquid) error {
	s := tryString(data)
	if len(s) == 0 {
		return json.Unmarshal(data, obj)
	}
	t, err := inventory.NewComponent(ctx, s)
	if err != nil {
		return err
	}
	*obj = *t
	return nil
}

func (a *unmarshaler) unmarshalFile(data []byte, obj *wtype.File) error {
	var blob api.Blob
	if err := jsonpb.Unmarshal(bytes.NewReader(data), &blob); err != nil {
		return err
	}

	hf := blob.GetHostFile()
	if a.ReadLocalFiles && hf != nil {
		bs, err := ioutil.ReadFile(hf.Filename)
		if err != nil {
			return err
		}
		blob = api.Blob{
			Name: blob.Name,
			From: &api.Blob_Bytes{
				Bytes: &api.FromBytes{
					Bytes: bs,
				},
			},
		}
	}

	var f wtype.File
	if err := f.UnmarshalBlob(&blob); err != nil {
		return err
	}
	*obj = f
	return nil
}

func (a *unmarshaler) unmarshalStruct(ctx context.Context, data []byte, obj interface{}) error {
	var err error
	switch obj := obj.(type) {
	case *wtype.LHTipbox:
		err = a.unmarshalLHTipbox(ctx, data, obj)
	case *wtype.Plate:
		err = a.unmarshalLHPlate(ctx, data, obj)
	case *wtype.Liquid:
		err = a.unmarshalLHComponent(ctx, data, obj)
	case *wtype.File:
		err = a.unmarshalFile(data, obj)
	default:
		err = json.Unmarshal(data, obj)
	}

	return err
}

func setParam(ctx context.Context, um *unmarshaler, w *workflow.Workflow, process, name string, data []byte, in map[string]interface{}) error {
	value, ok := in[name]
	if !ok {
		return errUnknownParam
	}

	m := &meta.Unmarshaler{
		Struct: func(data []byte, obj interface{}) error {
			return um.unmarshalStruct(ctx, data, obj)
		},
	}
	if err := m.Unmarshal(data, &value); err != nil {
		return err
	}

	return w.SetParam(workflow.Port{Process: process, Port: name}, value)
}

func setParams(ctx context.Context, w *workflow.Workflow, params *RawParams, readLocalFiles bool) (*mixer.Opt, error) {
	if params == nil {
		return nil, nil
	}

	um := &unmarshaler{
		ReadLocalFiles: readLocalFiles,
	}

	for process, params := range params.Parameters {
		c, err := w.FuncName(process)
		if err != nil {
			return nil, fmt.Errorf("cannot get component for process %q: %s", process, err)
		}
		runner, err := inject.Find(ctx, inject.NameQuery{
			Repo:  c,
			Stage: api.ElementStage_STEPS,
		})
		if err != nil {
			return nil, fmt.Errorf("unknown component %q: %s", c, err)
		}
		cr, ok := runner.(inject.TypedRunner)
		if !ok {
			return nil, fmt.Errorf("cannot get type information for component %q: type %T", c, runner)
		}
		in := inject.MakeValue(cr.Input())
		for name, value := range params {
			if err := setParam(ctx, um, w, process, name, value, in); err != nil {
				return nil, fmt.Errorf("cannot assign parameter %q of process %q to %s: %s",
					name, process, string(value), err)
			}
		}
	}

	return params.Config, nil
}
