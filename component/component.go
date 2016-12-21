package component

import (
	"errors"
	"reflect"

	"github.com/antha-lang/antha/inject"
	areflect "github.com/antha-lang/antha/reflect"
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

type ParamDesc struct {
	Name, Desc, Kind, Type string
}

type ComponentDesc struct {
	Desc   string
	Path   string
	Params []ParamDesc
}

type Component struct {
	Name        string
	Constructor func() interface{}
	Desc        ComponentDesc
}

// MakeParams returns zero instances for each input and output parameter.
func (a *Component) MakeParams() (map[string]interface{}, error) {
	params, err := a.params()
	if err != nil {
		return nil, err
	}

	m := make(map[string]interface{})
	for _, v := range params {
		z := reflect.Zero(v.Type)
		if _, seen := m[v.Name]; seen {
			return nil, &alreadySeen{v.Name}
		}
		m[v.Name] = z.Interface()
	}
	return m, nil
}

type tdesc struct {
	Name string
	Type reflect.Type
}

func typeOf(i interface{}) ([]tdesc, error) {
	var tdescs []tdesc
	// Generated elements always have type *XXXOutput or *XXXInput
	t := reflect.TypeOf(i).Elem()
	if t.Kind() != reflect.Struct {
		return nil, invalidComponent
	}
	for i, l := 0, t.NumField(); i < l; i += 1 {
		tdescs = append(tdescs, tdesc{Name: t.Field(i).Name, Type: t.Field(i).Type})
	}
	return tdescs, nil
}

func (a *Component) params() ([]tdesc, error) {
	r, ok := a.Constructor().(inject.TypedRunner)
	if !ok {
		return nil, invalidComponent
	}

	inTypes, err := typeOf(r.Input())
	if err != nil {
		return nil, err
	}
	outTypes, err := typeOf(r.Output())
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
		if err := add(v.Name, areflect.FullTypeName(v.Type)); err != nil {
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
