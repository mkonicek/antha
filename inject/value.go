package inject

import (
	"errors"
	"fmt"
	"reflect"
)

var (
	errNotStructOrValue = errors.New("not pointer to struct or map[string]interface{}")
	errDuplicateField   = errors.New("duplicate field")
)

var (
	stringType = reflect.TypeOf("")
	zeroValue  reflect.Value
)

// A Value is an input or output of injectable functions. It models named and
// typed function parameters.
type Value map[string]interface{}

// Concat concatenates the fields of one value with another. Returns an error
// if one value shares the same fields as another.
func (a Value) Concat(b Value) (Value, error) {
	r := make(Value)
	for k, v := range a {
		r[k] = v
	}
	for k, v := range b {
		if _, seen := r[k]; seen {
			return nil, errDuplicateField
		}
		r[k] = v
	}
	return r, nil
}

func makeValue(value reflect.Value) Value {
	var m Value

	switch value.Type().Kind() {
	case reflect.Map:
		m = make(Value)
		for _, key := range value.MapKeys() {
			skey, ok := key.Interface().(string)
			if ok {
				m[skey] = value.MapIndex(key).Interface()
			}
		}
	case reflect.Struct:
		m = make(Value)
		typ := value.Type()
		for idx, l := 0, value.NumField(); idx < l; idx++ {
			name := typ.Field(idx).Name
			m[name] = value.Field(idx).Interface()
		}
	case reflect.Ptr:
		return makeValue(value.Elem())
	default:
	}
	return m
}

// MakeValue makes a value from an object
func MakeValue(x interface{}) Value {
	return makeValue(reflect.ValueOf(x))
}

func isStringMap(value reflect.Value) bool {
	typ := value.Type()
	if typ.Kind() != reflect.Map {
		return false
	} else if typ.Key() != stringType {
		return false
	} else {
		return true
	}
}

func makeMapFields(value reflect.Value) map[string]reflect.Value {
	fields := make(map[string]reflect.Value)
	for _, key := range value.MapKeys() {
		k := key.Interface().(string)
		v := value.MapIndex(key)
		// Try to resolve interfaces to concrete type
		if v.Kind() == reflect.Interface {
			if vactual := v.Elem(); vactual != zeroValue {
				v = vactual
			}
		}
		fields[k] = v
	}
	return fields
}

func makeStructFields(value reflect.Value) map[string]reflect.Value {
	typ := value.Type()
	fields := make(map[string]reflect.Value)
	for i, n := 0, typ.NumField(); i < n; i++ {
		sf := typ.Field(i)
		v := value.Field(i)
		// Try to resolve interfaces to concrete type
		if v.Kind() == reflect.Interface {
			if vactual := v.Elem(); vactual != zeroValue {
				v = vactual
			}
		}
		fields[sf.Name] = v
	}
	return fields
}

func makeFields(value reflect.Value) (map[string]reflect.Value, error) {
	typ := value.Type()
	switch k := typ.Kind(); {
	case k == reflect.Ptr:
		return makeFields(value.Elem())
	case k == reflect.Interface:
		return makeFields(value.Elem())
	case k == reflect.Struct:
		return makeStructFields(value), nil
	case isStringMap(value):
		return makeMapFields(value), nil
	}
	return nil, errNotStructOrValue
}

func nilable(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Chan, reflect.Func, reflect.Map, reflect.Ptr, reflect.Interface, reflect.Slice:
		return true
	}
	return false
}

func assign(from, to interface{}, set, ignoreMissing bool) error {
	toValue := reflect.ValueOf(to)
	toFields, err := makeFields(toValue)
	if err != nil {
		return err
	}

	fromFields, err := makeFields(reflect.ValueOf(from))
	if err != nil {
		return err
	}

	isToMap := isStringMap(toValue)
	for name, v := range fromFields {
		toV, ok := toFields[name]

		// Special case: if "from" is nil, AssignableTo may return false (e.g.,
		// from: interface{}, to: error), but it is always okay to set "to" to
		// nil
		toNil := nilable(v) && v.IsNil()
		if !ok && ignoreMissing {
			continue
		} else if !ok {
			return fmt.Errorf("missing field %q", name)
		} else if fromT, toT := v.Type(), toV.Type(); !fromT.AssignableTo(toT) && !toNil {
			return fmt.Errorf("field %q of type %s not assignable to type %s", name, fromT, toT)
		} else if !isToMap && !toV.CanSet() {
			return fmt.Errorf("cannot set field %q", name)
		} else if !set {
			continue
		}

		if toNil {
			v = reflect.Zero(toV.Type())
		}

		if isToMap {
			toValue.SetMapIndex(reflect.ValueOf(name), v)
		} else {
			toV.Set(v)
		}
	}

	return nil
}

// AssignableTo returns if src is assignable to dst. Typing rule is as follows:
// (1) every field of src must have a field of the same name in dst, (2) the
// type of the src field must be golang assignable to the dst field, and (3)
// the dst fields must be golang settable (i.e., have an address).
func AssignableTo(src, dst interface{}) error {
	return assign(src, dst, false, false)
}

// Assign assigns values from Value or struct to Value or struct. If src is not
// AssignableTo dst, return an error.
func Assign(src, dst interface{}) error {
	return assign(src, dst, true, false)
}

// AssignSome assigns some values from Value or struct to Value or struct. Like
// Assign but ignore fields in src that are not present in dst.
func AssignSome(src, dst interface{}) error {
	return assign(src, dst, true, true)
}
