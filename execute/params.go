package execute

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"

	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/bvendor/golang.org/x/net/context"
	"github.com/antha-lang/antha/inject"
	"github.com/antha-lang/antha/microArch/factory"
	"github.com/antha-lang/antha/target/mixer"
	"github.com/antha-lang/antha/workflow"
)

type constructor func(string) interface{}

var (
	ptipbox         wtype.LHTipbox
	pplate          wtype.LHPlate
	pcomponent      wtype.LHComponent
	unknownParam    = errors.New("unknown parameter")
	cannotConstruct = errors.New("cannot construct parameter")
	nilValue        reflect.Value
	constructors    = map[reflect.Type]constructor{
		reflect.TypeOf(ptipbox):    func(x string) interface{} { return factory.GetTipByType(x) },
		reflect.TypeOf(pplate):     func(x string) interface{} { return factory.GetPlateByType(x) },
		reflect.TypeOf(pcomponent): func(x string) interface{} { return factory.GetComponentByType(x) },
	}
)

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

// Structure of parameter data for unmarshalling
type RawParams struct {
	Parameters map[string]map[string]json.RawMessage
	Config     *mixer.Opt
}

// Structure of parameter data for marshalling
type Params struct {
	Parameters map[string]map[string]interface{}
	Config     *mixer.Opt
}

func findConstructor(typ reflect.Type) constructor {
	// XXX: consider supporting convertible types too
	return constructors[typ]
}

func unmarshalOne(value reflect.Value, data []byte) (reflect.Value, error) {
	typ := value.Type()
	newValue := reflect.New(typ).Interface()
	origErr := json.Unmarshal(data, newValue)
	if origErr != nil {
		// Try to run constructor instead
		var carg string
		if construct := findConstructor(typ); construct == nil {
			return nilValue, origErr
		} else if err := json.Unmarshal(data, &carg); err != nil {
			return nilValue, fmt.Errorf("%s: %s", err, origErr)
		} else if v, err := constructOrError(construct, carg); err != nil {
			return nilValue, fmt.Errorf("%s: %s", err, origErr)
		} else if v == nil {
			return nilValue, cannotConstruct
		} else {
			newValue = v
		}
	}

	return reflect.ValueOf(newValue).Elem(), nil
}

func unmarshal(value reflect.Value, data []byte) (reflect.Value, error) {
	typ := value.Type()
	switch typ.Kind() {
	case reflect.Slice:
		raw := make([]json.RawMessage, 0)
		if err := json.Unmarshal(data, &raw); err != nil {
			return nilValue, err
		}
		s := reflect.MakeSlice(typ, 0, 0)
		for idx, bs := range raw {
			svalue := reflect.Zero(typ.Elem())
			if idx < value.Len() {
				svalue = reflect.ValueOf(value.Index(idx).Interface())
			}
			v, err := unmarshal(svalue, bs)
			if err != nil {
				return nilValue, err
			}
			s = reflect.Append(s, v)
		}
		return s, nil
	case reflect.Map:
		raw := make(map[string]json.RawMessage)
		if err := json.Unmarshal(data, &raw); err != nil {
			return nilValue, err
		}
		m := reflect.MakeMap(typ)
		for k, bs := range raw {
			kvalue := reflect.ValueOf(k)
			mvalue := value.MapIndex(kvalue)
			if mvalue == nilValue {
				mvalue = reflect.Zero(typ.Elem())
			} else {
				mvalue = reflect.ValueOf(mvalue.Interface())
			}
			v, err := unmarshal(mvalue, bs)
			if err != nil {
				return nilValue, err
			}
			m.SetMapIndex(kvalue, v)
		}
		return m, nil
	case reflect.Struct:
		return unmarshalOne(value, data)
	case reflect.Ptr:
		etyp := typ.Elem()
		if c := findConstructor(etyp); c == nil {
			return unmarshal(reflect.Zero(etyp), data)
		} else if v, err := unmarshalOne(reflect.Zero(etyp), data); err != nil {
			return nilValue, err
		} else {
			return v.Addr(), nil
		}
	default:
		return unmarshalOne(value, data)
	}
}

func setParam(w *workflow.Workflow, process, name string, data []byte, in map[string]interface{}) error {
	prototype, ok := in[name]
	if !ok {
		return unknownParam
	}

	value, err := unmarshal(reflect.ValueOf(prototype), data)
	if err != nil {
		return err
	}
	return w.SetParam(workflow.Port{Process: process, Port: name}, value.Interface())
}

func setParams(ctx context.Context, params *RawParams, w *workflow.Workflow) (*mixer.Opt, error) {
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
