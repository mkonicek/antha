package data

import (
	"reflect"

	"github.com/pkg/errors"
)

// UpdateSelection stores a table and its column to update.
type UpdateSelection struct {
	t   *Table
	col ColumnName
}

// By performs Update using the whole table row as input.
// For wide tables this might be inefficient, consider using On(...) instead.
func (us *UpdateSelection) By(fn func(r Row) interface{}) (*Table, error) {
	return us.updateByExtend(func(t *Table, colName ColumnName, colType reflect.Type) (*Table, error) {
		// extending by the whole row function
		return t.Extend(colName).By(fn, colType), nil
	})
}

// Constant replaces a column with a constant column having the given value, which must be of the same type.
// Returns error if a non-nil value cannot be converted to the column type.
func (us *UpdateSelection) Constant(value interface{}) (*Table, error) {
	return us.updateByExtend(func(t *Table, colName ColumnName, colType reflect.Type) (*Table, error) {
		// extending by a constant column
		return t.Extend(colName).ConstantType(value, colType)
	})
}

// UpdateOn stores a table, its column to update and the source columns list.
type UpdateOn struct {
	us        *UpdateSelection
	inputCols []ColumnName
}

// On selects a subset of columns to use as an extension source. If duplicate columns
// exist, the first so named is used.  Note this does not panic yet, even if the
// columns do not exist.  (However subsequent calls to the returned object will
// error.)
func (us *UpdateSelection) On(cols ...ColumnName) *UpdateOn {
	return &UpdateOn{us: us, inputCols: cols}
}

// Interface updates a column of an arbitrary type using a subset of source columns of arbitrary types.
func (on *UpdateOn) Interface(fn func(v ...interface{}) interface{}) (*Table, error) {
	return on.us.updateByExtend(func(t *Table, colName ColumnName, colType reflect.Type) (*Table, error) {
		// extending on columns of arbitrary type by a column of arbitrary type
		return t.Extend(colName).On(on.inputCols...).Interface(fn, colType)
	})
}

// A generic Update implementation on the top on Extend. Concrete extend function should be provided by caller.
func (us *UpdateSelection) updateByExtend(extend func(t *Table, colName ColumnName, colType reflect.Type) (*Table, error)) (*Table, error) {
	// index of the column to update
	colIndex, err := us.t.schema.ColIndex(us.col)
	if err != nil {
		return nil, err
	}

	// extending by a new column of the same name
	extended, err := extend(us.t, us.col, us.t.schema.Columns[colIndex].Type)
	if err != nil {
		return nil, errors.Wrapf(err, "updating column '%s'", us.col)
	}
	series := extended.series

	// the extend function should add one column only
	if len(series) != len(us.t.series)+1 {
		panic(errors.New("SHOULD NOT HAPPEN; wrong number of columns while update"))
	}
	// replacing the column to update with the extension column
	series[colIndex] = series[len(series)-1]
	series = series[:len(series)-1]

	// transferring the key is possible
	for _, keyCol := range us.t.sortKey {
		if keyCol.Column == us.col {
			// key column is updated => impossible to transfer the key
			return NewTable(series...), nil
		}
	}
	return newFromSeries(series, us.t.sortKey...), nil
}
