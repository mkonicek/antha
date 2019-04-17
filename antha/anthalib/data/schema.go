package data

import (
	"reflect"

	"github.com/pkg/errors"
)

// ColumnName is a type for column names in schema
type ColumnName string

// Index is a type for rows numbers
type Index int

// Schema is intended as an immutable representation of table metadata
type Schema struct {
	Columns []Column
	byName  map[ColumnName][]int
}

// Key defines sorting columns (and directions)
type Key []ColumnKey

// ColumnKey defines ordering on a single column
type ColumnKey struct {
	Column ColumnName
	Asc    bool
}

// Column is Series metadata
type Column struct {
	Name ColumnName
	Type reflect.Type
}

// SchemaAssertion is an interface for arbitrary assertions which might be required from schema.
// It should return nil if the schema is acceptable.
type SchemaAssertion func(Schema) error

// NumColumns is number of columns
func (s Schema) NumColumns() int {
	return len(s.Columns)
}

// Equal returns true if order, name, and types match
func (s Schema) Equal(other Schema) bool {
	if s.NumColumns() != other.NumColumns() {
		return false
	}
	for i, c := range s.Columns {
		co := other.Columns[i]
		if c.Name != co.Name || !typesEqual(c.Type, co.Type) {
			return false
		}
	}
	return true
}

// NB types are cached but exact equality is hard to test for, see reflect.go

func typesEqual(t1, t2 reflect.Type) bool {
	if t1 != t2 {
		if t1.Kind() != t2.Kind() {
			return false
		}

		if !t1.AssignableTo(t2) || !t2.AssignableTo(t1) {
			return false
		}
	}
	return true
}

// Col gets the column by name, first matched
func (s Schema) Col(col ColumnName) (Column, error) {
	index, err := s.ColIndex(col)
	if err != nil {
		return Column{}, err
	}
	return s.Columns[index], nil
}

// MustCol gets the column by name, first matched
func (s Schema) MustCol(colName ColumnName) Column {
	col, err := s.Col(colName)
	handle(err)
	return col
}

// ColIndex gets the column index by name, first matched
func (s Schema) ColIndex(col ColumnName) (int, error) {
	cs, found := s.byName[col]
	if !found {
		return -1, errors.Errorf("no such column: %v", col)
	}
	return cs[0], nil
}

// MustColIndex gets the column index by name, first matched; panics on error
func (s Schema) MustColIndex(col ColumnName) int {
	index, err := s.ColIndex(col)
	handle(err)
	return index
}

// Project projects a schema by given column names
func (s Schema) Project(colNames ...ColumnName) (*Schema, error) {
	columns := make([]Column, len(colNames))
	for i, name := range colNames {
		col, err := s.Col(name)
		if err != nil {
			return nil, err
		}
		columns[i] = col
	}
	return NewSchema(columns), nil
}

// MustProject projects a schema by given column names; panics on error
func (s Schema) MustProject(colNames ...ColumnName) *Schema {
	schema, err := s.Project(colNames...)
	handle(err)
	return schema
}

// Check checks arbitrary assertions about the schema
func (s Schema) Check(assertions ...SchemaAssertion) error {
	for _, assrt := range assertions {
		if err := assrt(s); err != nil {
			return err
		}
	}
	return nil
}

// CheckColumnsExist checks that all given columns exist
func (s Schema) CheckColumnsExist(colNames ...ColumnName) error {
	for _, colName := range colNames {
		if _, err := s.ColIndex(colName); err != nil {
			return err
		}
	}
	return nil
}

// TODO String()

// NewSchema creates a new schema by a columns list
func NewSchema(columns []Column) *Schema {
	schema := &Schema{
		Columns: columns,
		byName:  map[ColumnName][]int{},
	}
	for i, column := range columns {
		schema.byName[column.Name] = append(schema.byName[column.Name], i)
	}
	return schema
}

func newSchema(series []*Series) *Schema {
	schema := &Schema{byName: map[ColumnName][]int{}}
	for c, s := range series {
		schema.Columns = append(schema.Columns, Column{Type: s.typ, Name: s.col})
		schema.byName[s.col] = append(schema.byName[s.col], c)
	}
	return schema
}

// HasPrefix checks if other key is a prefix of k
func (key Key) HasPrefix(other Key) bool {
	if len(other) > len(key) {
		return false
	}
	for i := range other {
		if !other[i].Equal(key[i]) {
			return false
		}
	}
	return true
}

// Equal checks if two keys are equal
func (key Key) Equal(other Key) bool {
	return len(other) == len(key) && key.HasPrefix(other)
}

// Equal checks if two key entries are equal
func (ck ColumnKey) Equal(other ColumnKey) bool {
	return ck.Column == other.Column && ck.Asc == other.Asc
}
