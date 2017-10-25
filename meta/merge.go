package meta

import (
	"errors"
	"reflect"
)

var (
	errNotStruct      = errors.New("expecting struct or pointer to struct")
	errDifferentTypes = errors.New("values have different types")
)

// Is value an instance of zero of type?
func zero(v reflect.Value) bool {
	t := v.Type()
	z := reflect.Zero(t)

	switch t.Kind() {
	case reflect.Slice:
		fallthrough
	case reflect.Array:
		fallthrough
	case reflect.Map:
		return v.Len() == 0

	default:
		return v.Interface() == z.Interface()
	}
}

func shallowMerge(val1, val2 reflect.Value) (reflect.Value, error) {
	typ1 := val1.Type()
	typ2 := val2.Type()

	if typ1 != typ2 {
		return reflect.Value{}, errDifferentTypes
	} else if typ1.Kind() == reflect.Ptr {
		r, err := shallowMerge(val1.Elem(), val2.Elem())
		if err != nil {
			return reflect.Value{}, err
		}
		return r.Addr(), nil
	} else if typ1.Kind() != reflect.Struct {
		return reflect.Value{}, errNotStruct
	}

	new := reflect.New(typ1)
	rv := new.Elem()
	for f, numField := 0, val1.NumField(); f < numField; f++ {
		f1 := val1.Field(f)
		f2 := val2.Field(f)
		if f1.Type() != f2.Type() {
			return reflect.Value{}, errDifferentTypes
		}
		if !zero(f2) {
			rv.Field(f).Set(f2)
		} else {
			rv.Field(f).Set(f1)
		}
	}
	return rv, nil
}

// ShallowMerge merges value2 into value1 and return the result.
//
// Non-zero field values in value2 override the corresponding values in value1.
func ShallowMerge(value1, value2 interface{}) (interface{}, error) {
	r, err := shallowMerge(reflect.ValueOf(value1), reflect.ValueOf(value2))
	if err != nil {
		return nil, err
	}
	return r.Interface(), nil
}
