package data

import (
	"reflect"
	"strconv"

	"github.com/pkg/errors"
)

/*
 * reflectively created index keys
 */

// In order to support indexing by multiple columns, we use reflection - i.e. we reflectively create a struct to accomodate an index key:
//  type indexKey struct {
//      Field1 int64
//      NotNull1 bool
//      Field2 string
//      NotNull1 bool
//      //etc.
//  }
// and reflectively create maps to store index:
//  type indexSet map[indexKey]bool
//  type indexMap map[indexKey]someType
// An alternative approach considered was to use external hashing function (e.g. "github.com/mitchellh/hashstructure"),
// but it turned out to be ~2 times slower.

// reflectively creates a struct type to accomodate a row of a given schema
func makeIndexKeyType(schema *Schema) (reflect.Type, error) {
	fields := make([]reflect.StructField, 0)
	for i, column := range schema.Columns {
		// check if the type is comparable
		if err := checkTypeComparable(column.Type); err != nil {
			return nil, errors.Wrapf(err, "column %s", column.Name)
		}
		// Field column.Type
		fields = append(fields, reflect.StructField{
			Name: "Field" + strconv.Itoa(i),
			Type: column.Type,
		})
		// NotNull bool
		fields = append(fields, reflect.StructField{
			Name: "NotNull" + strconv.Itoa(i),
			Type: reflect.TypeOf(false),
		})
	}
	return reflect.StructOf(fields), nil
}

// checks if the type is "comparable" (i.e. can be used as a Go map key)
func checkTypeComparable(typ reflect.Type) error {
	switch typ.Kind() {
	case reflect.Struct:
		for i := 0; i < typ.NumField(); i++ {
			if err := checkTypeComparable(typ.Field(i).Type); err != nil {
				return err
			}
		}
		return nil
	case reflect.Array:
		return checkTypeComparable(typ.Elem())
	case reflect.Map, reflect.Slice, reflect.Func:
		return errors.Errorf("type %+v is not comparable", typ)
	default:
		return nil
	}
}

// load row values into a reflectively created index key struct
func loadIndexKeyFromRow(row raw, indexKey reflect.Value) {
	for i, value := range row {
		field := indexKey.Field(2 * i)
		notNull := indexKey.Field(2*i + 1)
		if value != nil {
			field.Set(reflect.ValueOf(value))
			notNull.Set(reflect.ValueOf(true))
		} else {
			field.Set(reflect.Zero(field.Type()))
			notNull.Set(reflect.ValueOf(false))
		}
	}
}

/*
 * indexSet
 */

// indexSet is a reflective hash set
type indexSet struct {
	index reflect.Value // map[indexKeyType]bool
}

func newIndexSet(keyType reflect.Type) indexSet {
	mapType := reflect.MapOf(keyType, reflect.TypeOf(false))
	return indexSet{
		index: reflect.MakeMap(mapType),
	}
}

// checks whether a value is already in the index
func (s *indexSet) has(key reflect.Value) bool {
	return s.index.MapIndex(key).IsValid()
}

// adds a value to the index
func (s *indexSet) add(key reflect.Value) {
	s.index.SetMapIndex(key, reflect.ValueOf(true))
}

/*
 * indexMap
 */

// indexMap is a reflective hash map
type indexMap struct {
	index reflect.Value // map[indexKeyType]indexValueType
}

func newIndexMap(keyType reflect.Type, valueType reflect.Type) indexMap {
	mapType := reflect.MapOf(keyType, valueType)
	return indexMap{
		index: reflect.MakeMap(mapType),
	}
}

func (m *indexMap) get(key reflect.Value) reflect.Value {
	return m.index.MapIndex(key)
}

func (m *indexMap) lookup(key reflect.Value) (reflect.Value, bool) {
	value := m.get(key)
	return value, value.IsValid()
}

func (m *indexMap) set(key reflect.Value, value reflect.Value) {
	m.index.SetMapIndex(key, value)
}
