package data

import (
	"reflect"

	"github.com/pkg/errors"
)

// ForeachSelection encapsulates the table to iterate with Foreach.
type ForeachSelection struct {
	t *Table
}

// By performs Foreach on the whole table rows.
// For wide tables this might be inefficient, consider using On(...) instead.
func (fs *ForeachSelection) By(fn func(r Row)) {
	iter := fs.t.read(fs.t.series)
	for iter.Next() {
		fn(iter.Value())
	}
}

// On selects columns for iterating on. Note this does not panic yet, even if the
// columns do not exist.  (However subsequent calls to the returned object will
// error.)
func (fs *ForeachSelection) On(cols ...ColumnName) *ForeachOn {
	return &ForeachOn{t: fs.t, cols: cols}
}

// ForeachOn encapsulated the table and the columns to iterate.
type ForeachOn struct {
	t    *Table
	cols []ColumnName
}

// Interface invokes a user-supplied function passing the named column values as interface{} arguments, including nil.
// If given any SchemaAssertions, they are called in the beginning and may have side effects.
func (o *ForeachOn) Interface(fn func(v ...interface{}), assertions ...SchemaAssertion) error {
	// schema checks
	if err := o.checkSchema(nil, assertions...); err != nil {
		return errors.Wrapf(err, "can't iterate over %+v", o.t)
	}

	// iterating over the projected table
	projected := o.t.Must().Project(o.cols...)
	iter := projected.read(projected.series)
	for iter.Next() {
		fn(iter.rawValue()...)
	}

	return nil
}

func (o *ForeachOn) checkSchema(colsType reflect.Type, assertions ...SchemaAssertion) error {
	// eager schema check
	projectedSchema, err := o.t.Schema().Project(o.cols...)
	if err != nil {
		return errors.Wrapf(err, "can't filter columns %+v", o.cols)
	}
	// assert columns assignable
	if colsType != nil {
		for _, col := range projectedSchema.Columns {
			if !col.Type.AssignableTo(colsType) {
				return errors.Errorf("column %s is not assignable to type %v", col.Name, colsType)
			}
		}
	}
	// checking assertions
	return projectedSchema.Check(assertions...)
}
