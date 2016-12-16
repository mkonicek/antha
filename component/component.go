package component

import (
	"errors"
	"fmt"
	"reflect"

	areflect "github.com/antha-lang/antha/reflect"

	"github.com/antha-lang/antha/inject"
)

var (
	invalidComponent = errors.New("invalid component")
)

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

// Update types in description of a component based return values of
// constructor.
func UpdateParamTypes(desc *Component) error {
	// Add type information if missing
	ts := make(map[string]string)

	type tdesc struct {
		Name string
		Type reflect.Type
	}

	add := func(name, t string) error {
		if _, seen := ts[name]; seen {
			return fmt.Errorf("parameter %q already seen", name)
		}
		ts[name] = t
		return nil
	}

	typeOf := func(i interface{}) ([]tdesc, error) {
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

	if r, ok := desc.Constructor().(inject.TypedRunner); !ok {
		return invalidComponent
	} else if inTypes, err := typeOf(r.Input()); err != nil {
		return err
	} else if outTypes, err := typeOf(r.Output()); err != nil {
		return err
	} else {
		for _, v := range append(inTypes, outTypes...) {
			if err := add(v.Name, areflect.FullTypeName(v.Type)); err != nil {
				return err
			}
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
