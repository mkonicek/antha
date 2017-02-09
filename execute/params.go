package execute

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/inject"
	"github.com/antha-lang/antha/meta"
	"github.com/antha-lang/antha/microArch/factory"
	"github.com/antha-lang/antha/target/mixer"
	"github.com/antha-lang/antha/workflow"
)

type constructor func(string) interface{}

var (
	unknownParam    = errors.New("unknown parameter")
	cannotConstruct = errors.New("cannot construct parameter")
)

// Deprecated for github.com/antha-lang/antha/api/v1/WorkflowParameters.
// Structure of parameter data for unmarshalling.
type RawParams struct {
	Parameters map[string]map[string]json.RawMessage `json:"parameters"`
	Config     *mixer.Opt                            `json:"config"`
}

// Deprecated for github.com/antha-lang/antha/api/v1/WorkflowParameters.
// Structure of parameter data for marshalling.
type Params struct {
	Parameters map[string]map[string]interface{} `json:"parameters"`
	Config     *mixer.Opt                        `json:"config"`
}

func constructOrError(fn func(x string) interface{}, x string) (interface{}, error) {
	var v interface{}
	var err error
	defer func() {
		if res := recover(); res != nil {
			err = fmt.Errorf("error making %q: %s", x, res)
		}
	}()
	v = fn(x)
	return v, err
}

func tryString(data []byte) string {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		s = ""
	}
	return s
}

type unmarshaler struct{}

func (a *unmarshaler) unmarshalStruct(data []byte, obj interface{}) error {
	switch obj := obj.(type) {
	case *wtype.LHTipbox:
		s := tryString(data)
		if len(s) != 0 {
			t, err := constructOrError(func(x string) interface{} { return factory.GetTipByType(x) }, s)
			if err != nil {
				return err
			}
			*obj = *t.(*wtype.LHTipbox)
			return nil
		}
	case *wtype.LHPlate:
		s := tryString(data)
		if len(s) != 0 {
			t, err := constructOrError(func(x string) interface{} { return factory.GetPlateByType(x) }, s)
			if err != nil {
				return err
			}
			*obj = *t.(*wtype.LHPlate)
			return nil
		}
	case *wtype.LHComponent:
		s := tryString(data)
		if len(s) != 0 {
			t, err := constructOrError(func(x string) interface{} { return factory.GetComponentByType(x) }, s)
			if err != nil {
				return err
			}
			*obj = *t.(*wtype.LHComponent)
			return nil
		}
	}

	return json.Unmarshal(data, obj)
}

func setParam(w *workflow.Workflow, process, name string, data []byte, in map[string]interface{}) error {
	value, ok := in[name]
	if !ok {
		return unknownParam
	}

	var u unmarshaler
	if err := meta.UnmarshalJSON(meta.UnmarshalOpt{
		Struct: u.unmarshalStruct,
	}, data, &value); err != nil {
		return err
	}

	return w.SetParam(workflow.Port{Process: process, Port: name}, value)
}

func setParams(ctx context.Context, w *workflow.Workflow, params *RawParams) (*mixer.Opt, error) {
	if params == nil {
		return nil, nil
	}

	for process, params := range params.Parameters {
		c, err := w.FuncName(process)
		if err != nil {
			return nil, fmt.Errorf("cannot get component for process %q: %s", process, err)
		}
		runner, err := inject.Find(ctx, inject.NameQuery{Repo: c})
		if err != nil {
			return nil, fmt.Errorf("unknown component %q: %s", c, err)
		}
		cr, ok := runner.(inject.TypedRunner)
		if !ok {
			return nil, fmt.Errorf("cannot get type information for component %q: type %T", c, runner)
		}
		in := inject.MakeValue(cr.Input())
		for name, value := range params {
			if err := setParam(w, process, name, value, in); err != nil {
				return nil, fmt.Errorf("cannot assign parameter %q of process %q to %s: %s",
					name, process, string(value), err)
			}
		}
	}
	return params.Config, nil
}
