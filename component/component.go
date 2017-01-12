package component

import (
	"errors"
	"reflect"

	"github.com/antha-lang/antha/inject"
	"github.com/antha-lang/antha/meta"
)

var (
	invalidComponent = errors.New("invalid component")
)

type alreadySeen struct {
	Name string
}

func (a *alreadySeen) Error() string {
	return "parameter " + a.Name + " already seen"
}

// Description of the parameters of a component.
type ParamDesc struct {
	Name string // Name of parameter
	Desc string // Description of parameter from doc string
	Kind string // Input, Output, Data or Parameter
	Type string // Full go type name
}

// Description of a component.
type ComponentDesc struct {
	Desc   string
	Path   string
	Params []ParamDesc
}

// An antha component / element.
type Component struct {
	Name        string
	Constructor func() interface{}
	Desc        ComponentDesc
}

// NewParams returns new objects instances for each input and output parameter.
//
// If a component has parameters:
//
//   Parameters (
//     String string
//     Number int
//   )
//   Data(...)
//   Inputs(...)
//   Outputs(...)
//
// The result of NewParams will be:
//
//   map[string]interface{} {
//     "String": new(string),
//     "Number": new(int),
//     ...
//   }
func (a *Component) NewParams() (map[string]interface{}, error) {
	params, err := a.params()
	if err != nil {
		return nil, err
	}

	m := make(map[string]interface{})
	for _, v := range params {
		z := reflect.New(v.Type)
		if _, seen := m[v.Name]; seen {
			return nil, &alreadySeen{v.Name}
		}
		m[v.Name] = z.Interface()
	}
	return m, nil
}

type typeDesc struct {
	Name string
	Type reflect.Type
}

func typeDescOf(obj interface{}) ([]typeDesc, error) {
	var tdescs []typeDesc
	// Generated elements always have type *XXXOutput or *XXXInput
	t := reflect.TypeOf(obj).Elem()
	if t.Kind() != reflect.Struct {
		return nil, invalidComponent
	}
	for i, l := 0, t.NumField(); i < l; i += 1 {
		tdescs = append(tdescs, typeDesc{Name: t.Field(i).Name, Type: t.Field(i).Type})
	}
	return tdescs, nil
}

func (a *Component) params() ([]typeDesc, error) {
	r, ok := a.Constructor().(inject.TypedRunner)
	if !ok {
		return nil, invalidComponent
	}

	inTypes, err := typeDescOf(r.Input())
	if err != nil {
		return nil, err
	}
	outTypes, err := typeDescOf(r.Output())
	if err != nil {
		return nil, err
	}
	return append(inTypes, outTypes...), nil
}

// UpdateParamTypes updates types in description of a component based return
// values of the constructor.
func UpdateParamTypes(desc *Component) error {
	// Add type information if missing
	ts := make(map[string]string)

	add := func(name, t string) error {
		if _, seen := ts[name]; seen {
			return &alreadySeen{name}
		}
		ts[name] = t
		return nil
	}

	params, err := desc.params()
	if err != nil {
		return err
	}

	for _, v := range params {
		obj := reflect.Zero(v.Type)
		if err := add(v.Name, meta.FullTypeName(obj)); err != nil {
			return err
		}
	}

	for i, p := range desc.Desc.Params {
		t := &desc.Desc.Params[i].Type
		if len(*t) == 0 {
			*t = ts[p.Name]
		}
	}

	return nil
}
