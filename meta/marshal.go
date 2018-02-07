package meta

import (
	"encoding/json"
	"errors"
	"reflect"
)

var (
	errMapKeyIsNotString = errors.New("map key is not string")
)

// MarshalFunc serializes data into an object
type MarshalFunc func(obj interface{}) ([]byte, error)

// Marshaler gives on custom on-the-side marshaling functions for specific
// golang object kinds.
type Marshaler struct {
	Struct MarshalFunc
}

func (a *Marshaler) marshalJSON(value reflect.Value) ([]byte, error) {
	if !value.IsValid() {
		return []byte("null"), nil
	}

	typ := value.Type()

	switch typ.Kind() {
	case reflect.Slice:
		var raw []*json.RawMessage
		for i, n := 0, value.Len(); i < n; i++ {
			elem := value.Index(i)
			bs, err := a.marshalJSON(elem)
			if err != nil {
				return nil, err
			}
			relem := json.RawMessage(bs)
			raw = append(raw, &relem)
		}
		return json.Marshal(raw)
	case reflect.Map:
		raw := make(map[string]*json.RawMessage)
		for _, kvalue := range value.MapKeys() {
			key, ok := kvalue.Interface().(string)
			if !ok {
				return nil, errMapKeyIsNotString
			}
			elem := value.MapIndex(kvalue)
			bs, err := a.marshalJSON(elem)
			if err != nil {
				return nil, err
			}
			relem := json.RawMessage(bs)
			raw[key] = &relem
		}

		return json.Marshal(raw)
	case reflect.Ptr:
		elem := value.Elem()
		if !elem.IsValid() {
			return json.Marshal(value.Interface())
		}

		return a.marshalJSON(elem)
	case reflect.Struct:
		if a.Struct == nil {
			// fall to default case
			break
		}
		return a.Struct(value.Interface())
	case reflect.Interface:
		// Use concrete type
		return a.marshalJSON(value.Elem())
	}

	return json.Marshal(value.Interface())
}

// Marshal parses obj and returns its serialization. Custom marshaling
// functions can be specified on the side.
func (a *Marshaler) Marshal(obj interface{}) ([]byte, error) {
	value := reflect.ValueOf(obj)

	return a.marshalJSON(value)
}
