package meta

import (
	"encoding/json"
	"errors"
	"reflect"
)

var (
	errNotAPointer = errors.New("not a pointer")
)

// UnmarshalFunc deserializes data into an object
type UnmarshalFunc func(data []byte, obj interface{}) error

// Unmarshaler gives on custom on-the-side unmarshaling functions for specific
// golang object kinds.
type Unmarshaler struct {
	Struct UnmarshalFunc
}

func (a *Unmarshaler) unmarshalJSON(data []byte, value reflect.Value) (reflect.Value, error) {
	var nilValue reflect.Value
	typ := value.Type()

	switch typ.Kind() {
	case reflect.Slice:
		raw := make([]json.RawMessage, 0)
		if err := json.Unmarshal(data, &raw); err != nil {
			return nilValue, err
		}
		slice := reflect.MakeSlice(typ, 0, 0)
		for idx, bs := range raw {
			elem := reflect.Zero(typ.Elem())
			if idx < value.Len() {
				elem = value.Index(idx)
			}
			v, err := a.unmarshalJSON(bs, elem)
			if err != nil {
				return nilValue, err
			}
			slice = reflect.Append(slice, v)
		}
		return slice, nil
	case reflect.Map:
		raw := make(map[string]json.RawMessage)
		if err := json.Unmarshal(data, &raw); err != nil {
			return nilValue, err
		}
		aMap := reflect.MakeMap(typ)
		for k, bs := range raw {
			kvalue := reflect.ValueOf(k)
			elem := value.MapIndex(kvalue)
			if !elem.IsValid() {
				elem = reflect.Zero(typ.Elem())
			}
			v, err := a.unmarshalJSON(bs, elem)
			if err != nil {
				return nilValue, err
			}
			aMap.SetMapIndex(kvalue, v)
		}
		return aMap, nil
	case reflect.Ptr:
		elem := value.Elem()
		if !elem.IsValid() {
			elem = reflect.Zero(typ.Elem())
		}

		v, err := a.unmarshalJSON(data, elem)
		if err != nil {
			return nilValue, err
		}
		vptr := reflect.New(elem.Type())
		vptr.Elem().Set(v)
		return vptr, nil
	case reflect.Struct:
		elem := reflect.New(typ)
		if a.Struct == nil {
			break
		}

		if err := a.Struct(data, elem.Interface()); err != nil {
			return nilValue, err
		}
		return elem.Elem(), nil
	case reflect.Interface:
		// Use concrete type
		return a.unmarshalJSON(data, value.Elem())
	}

	elem := reflect.New(typ)
	if err := json.Unmarshal(data, elem.Interface()); err != nil {
		return nilValue, err
	}
	return elem.Elem(), nil
}

// Unmarshal parses the JSON-encoded data and stores the result in the
// value pointed to by obj. Custom unmarshaling functions can be specified on
// the side.
func (a *Unmarshaler) Unmarshal(data []byte, obj interface{}) error {
	value := reflect.ValueOf(obj)
	if value.Kind() != reflect.Ptr {
		return errNotAPointer
	}

	elem := value.Elem()
	v, err := a.unmarshalJSON(data, elem)
	if err != nil {
		return err
	}

	elem.Set(v)
	return nil
}
