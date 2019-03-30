package data

import (
	"reflect"

	"github.com/pkg/errors"
)

/*
 * distinct interfaces
 */

// DistinctSelection contains a table to apply distinct.
type DistinctSelection struct {
	t *Table
}

// On selects distinct rows by columns specified.
func (ds *DistinctSelection) On(cols ...ColumnName) (*Table, error) {
	// TODO: fast path if the columns list matches the table key prefix (regardless of order)

	// subset of columns we are interested in
	projection, err := newProjection(ds.t.schema, cols...)
	if err != nil {
		return nil, errors.Wrap(err, "distinct columns lookup")
	}

	// creating the types needed for a reflective index
	indexKeyType, err := makeIndexKeyType(ds.t.schema.MustProject(cols...))
	if err != nil {
		return nil, err
	}

	// filter function generator
	filterFuncGen := func() rawMatch {
		// creating an index
		index := newIndexSet(indexKeyType)
		// creating a reflective key struct to load each row into
		indexKey := reflect.New(indexKeyType).Elem()

		// creating filter function
		return func(r raw) bool {
			// row key (in the form of []interface{})
			row := r.project(projection)
			// row key (in the form of reflectively created struct)
			loadIndexKeyFromRow(row, indexKey)
			// lookup key in the index
			if index.has(indexKey) {
				return false
			}
			// add key to the index
			index.add(indexKey)
			return true
		}
	}

	return filterTable(filterFuncGen, ds.t), nil
}
