package data

import (
	"reflect"
	"strings"

	"github.com/pkg/errors"
)

/*
 * data binding between tables and go struct types
 */

type fieldBinding struct {
	ColumnName
	fieldIdx []int
}
type fieldMap struct {
	destType reflect.Type
	// ptr depth
	level        int
	columnFields []fieldBinding
	// field indexes for fields that should receive the index value as Int
	indexFields [][]int
}

func structTypeInfo(structType reflect.Type) (*fieldMap, error) {
	f := &fieldMap{destType: structType}
	for f.destType.Kind() == reflect.Ptr {
		f.destType = structType.Elem()
		f.level++
	}
	if f.destType.Kind() != reflect.Struct {
		return nil, errors.Errorf("expecting struct or *struct, got %+v", f.destType)
	}
	for i := 0; i < f.destType.NumField(); i++ {
		field := f.destType.Field(i)
		if field.PkgPath != "" {
			// unexported fields, never used
			continue
		}
		name := ColumnName(field.Name)
		if tag, found := field.Tag.Lookup("table"); found {
			if tag == "-" {
				continue
			}
			bits := strings.Split(tag, ",")
			if len(bits) > 1 && bits[1] == "index" {
				// TODO check assignable to int
				f.indexFields = append(f.indexFields, field.Index)
				continue
			}
			if bits[0] != "" {
				name = ColumnName(bits[0])
			}
		}
		f.columnFields = append(f.columnFields, fieldBinding{ColumnName: name, fieldIdx: field.Index})
	}
	return f, nil
}

type colMapping struct {
	colIdx   int
	fieldIdx []int
}

func writeRow(f *fieldMap, r Row, m []colMapping) reflect.Value {

	destValue := reflect.New(f.destType).Elem()
	for _, m := range m {
		newVal := r.ValueAt(m.colIdx).Interface()
		if newVal != nil {
			destValue.FieldByIndex(m.fieldIdx).Set(reflect.ValueOf(newVal))
		}
	}
	for _, fieldIdx := range f.indexFields {
		destValue.FieldByIndex(fieldIdx).Set(reflect.ValueOf(int(r.Index())))
	}

	for i := 0; i < f.level; i++ {
		destValue = destValue.Addr()
	}
	return destValue
}

// ToStructs copies to the exported struct fields by name, ignoring unmapped
// columns (column names that do not correspond to a struct field). structsPtr
// must be of a *[]struct{} or *[]*struct{} type. Returns error if the struct
// fields are not assignable from the input data.
//
// When columns have null value, the struct field receives the zero value for
// that type instead.
//
// Use the struct tag `table:"column name"` to bind the field to a column with a
// different name.
//
// Use the struct tag `table:"-"` to ignore a field.
//
// Use the struct tag  `table:",index"` will mark that field as receiving the
// row index value (it must be an int field).
func (t *Table) ToStructs(structsPtr interface{}) error {
	// TODO optional error on unmapped column
	// TODO error on ambiguous column name
	v := reflect.ValueOf(structsPtr)
	if v.Kind() != reflect.Ptr {
		return errors.Errorf("expecting ptr to slice, got %+v", v.Type())
	}
	vSlice := v.Elem()
	f, err := structTypeInfo(vSlice.Type().Elem())
	if err != nil {
		return err
	}
	columnFieldMapping := []colMapping{}
	selectedColumnNames := []ColumnName{}
	schema := t.Schema()
	for i, binding := range f.columnFields {
		colIdx, err := schema.ColIndex(binding.ColumnName)
		if err != nil {
			return errors.Wrapf(err, "when mapping to struct type %+v", f.destType)
		}
		col := schema.Columns[colIdx]
		selectedColumnNames = append(selectedColumnNames, col.Name)
		field := f.destType.FieldByIndex(binding.fieldIdx)
		// type check
		if !col.Type.AssignableTo(field.Type) {
			return errors.Errorf("can't map column %v to struct field %v of type %+v", col, field.Name, field.Type)
		}
		// columnindex is set to i in the projection
		columnFieldMapping = append(columnFieldMapping, colMapping{i, field.Index})
	}
	for r := range t.Must().Project(selectedColumnNames...).IterAll() {
		destValue := writeRow(f, r, columnFieldMapping)
		vSlice = reflect.Append(vSlice, destValue)
	}

	v.Elem().Set(vSlice)
	return nil
}

// NewTableFromStructs copies exported struct fields to a new table. structs
// must be of a []struct{} or []*struct{} type. A nil entry in a []*struct{} is
// mapped to zero values for each column (not nulls!).
//
// Use the struct tag `table:"column name"` to bind the field to a column with a
// different name.
//
// Use the struct tag `table:"-"` to ignore a field.
func NewTableFromStructs(structs interface{}) (*Table, error) {
	// TODO embedded fields, anonymous fields, unexported fields?
	// TODO this may be grossly inefficient... unsafe.Offsetof would be one approach to optimize
	vSlice := reflect.ValueOf(structs)
	if vSlice.Kind() != reflect.Slice {
		return nil, errors.Errorf("expecting a slice, got a %+v", vSlice.Type())
	}
	f, err := structTypeInfo(vSlice.Type().Elem())
	if err != nil {
		return nil, err
	}
	length := vSlice.Len()
	// reflectively construct slices of concrete scalar type with reflect.SliceOf.  This reduces the size of
	// each slice compared to using a []interface{}.
	slices := make([]reflect.Value, len(f.columnFields))
	for i, binding := range f.columnFields {
		field := f.destType.FieldByIndex(binding.fieldIdx)
		slices[i] = reflect.MakeSlice(reflect.SliceOf(field.Type), length, length)
	}
row:
	for j := 0; j < length; j++ {
		rowVal := vSlice.Index(j)
		for x := 0; x < f.level; x++ {
			if rowVal.IsNil() {
				// skip this row; there will be a zero in each slice
				continue row
			}
			rowVal = rowVal.Elem()
		}
		for i, binding := range f.columnFields {
			val := rowVal.FieldByIndex(binding.fieldIdx)
			slices[i].Index(j).Set(reflect.ValueOf(val.Interface()))
		}
	}
	// create a Series for each exported field
	series := make([]*Series, len(f.columnFields))
	createSlice := reflect.ValueOf(newNativeSeriesFromSlice).Call
	var notNull []bool
	for i, binding := range f.columnFields {
		serResult := createSlice([]reflect.Value{reflect.ValueOf(binding.ColumnName), slices[i], reflect.ValueOf(notNull)})
		err, isErr := serResult[1].Interface().(error)
		if isErr {
			return nil, err
		}
		series[i] = serResult[0].Interface().(*Series)
	}
	return NewTable(series...), nil
}

// ToStruct copies to the exported struct fields by name, ignoring unmapped
// columns (column names that do not correspond to a struct field). structPtr
// must be of a *struct{} type.
//
// See documentation on (t *Table) ToStructs.
//
// Note this is much less efficient than handling a slice at a time.
func (r Row) ToStruct(structPtr interface{}) error {
	v := reflect.ValueOf(structPtr)
	f, err := structTypeInfo(v.Type())
	if err != nil {
		return err
	}
	if f.level != 1 {
		return errors.Errorf("expecting ptr to struct, got %+v", v.Type())
	}

	// build a pseudo schema using the inferred column types
	columnFieldMapping := []colMapping{}
	for _, binding := range f.columnFields {
		idx, err := r.schema.ColIndex(binding.ColumnName)
		if err != nil {
			return errors.Wrapf(err, "when mapping to struct type %+v", f.destType)
		}
		value := r.ValueAt(idx).Interface()

		field := f.destType.FieldByIndex(binding.fieldIdx)
		// type check
		if value != nil && !reflect.ValueOf(value).Type().AssignableTo(field.Type) {
			return errors.Errorf("can't map value %+v to struct field %v of type %+v", value, field.Name, field.Type)
		}
		columnFieldMapping = append(columnFieldMapping, colMapping{idx, field.Index})
	}
	destValue := writeRow(f, r, columnFieldMapping)
	v.Elem().Set(destValue.Elem())

	return nil
}
