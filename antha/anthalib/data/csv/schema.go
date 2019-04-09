package csv

import (
	"reflect"
	"strconv"

	"github.com/antha-lang/antha/antha/anthalib/data"
	"github.com/pkg/errors"
)

// a list of types supported by CSV reader and writer
type csvType struct {
	name  string
	typ   reflect.Type
	parse func(string) (interface{}, error)
}

// TODO use shared type registry
var csvTypes = [...]csvType{
	{name: "bool", typ: reflect.TypeOf(false), parse: parseBool},
	{name: "int64", typ: reflect.TypeOf(int64(0)), parse: parseInt64},
	{name: "int", typ: reflect.TypeOf(0), parse: parseInt},
	{name: "float64", typ: reflect.TypeOf(float64(0)), parse: parseFloat64},
	{name: "string", typ: reflect.TypeOf(""), parse: parseString},
	{name: "TimestampMillis", typ: reflect.TypeOf(data.TimestampMillis(0)), parse: parseTimestampMillis},
	{name: "TimestampMicros", typ: reflect.TypeOf(data.TimestampMicros(0)), parse: parseTimestampMicros},
}

var byName = make(map[string]*csvType)
var byType = make(map[reflect.Type]*csvType)

func init() {
	for i := range csvTypes {
		byName[csvTypes[i].name] = &csvTypes[i]
		byType[csvTypes[i].typ] = &csvTypes[i]
	}
}

func csvTypeByName(name string) (*csvType, error) {
	csvType, ok := byName[name]
	if !ok {
		return nil, errors.Errorf("unsupported type: %s", name)
	}
	return csvType, nil
}

func csvTypeByReflectType(typ reflect.Type) (*csvType, error) {
	csvType, ok := byType[typ]
	if !ok {
		return nil, errors.Errorf("unsupported type: %v", typ)
	}
	return csvType, nil
}

func parseBool(s string) (interface{}, error) {
	value, err := strconv.ParseBool(s)
	if err != nil {
		return nil, err
	}
	return value, nil
}

func parseInt64(s string) (interface{}, error) {
	value, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return nil, err
	}
	return value, nil
}

func parseInt(s string) (interface{}, error) {
	value, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return nil, err
	}
	return int(value), nil
}

func parseTimestampMillis(s string) (interface{}, error) {
	value, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return nil, err
	}
	return data.TimestampMillis(value), nil
}

func parseTimestampMicros(s string) (interface{}, error) {
	value, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return nil, err
	}
	return data.TimestampMicros(value), nil
}

func parseFloat64(s string) (interface{}, error) {
	value, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return nil, err
	}
	return value, nil
}

func parseString(s string) (interface{}, error) {
	return s, nil
}
