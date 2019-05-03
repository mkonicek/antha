package data

import (
	"fmt"
	"reflect"

	"github.com/pkg/errors"
)

/*
 * filter interfaces
 */

// FilterSelection is the fluent interface for filtering rows.
type FilterSelection struct {
	t *Table
}

// By accepts a filter function that operates on rows of the whole table.
// The filtered table contains rows where this function returns true.
func (fs *FilterSelection) By(fn MatchRow) *Table {
	matchGen := func() rawMatch {
		return func(r raw) bool {
			return fn(r.row(-1, &fs.t.schema))
		}
	}

	return filterTable(matchGen, fs.t)
}

// On selects columns for filtering. Note this does not panic yet, even if the
// columns do not exist.  (However subsequent calls to the returned object will
// error.)
func (fs *FilterSelection) On(cols ...ColumnName) *FilterOn {
	return &FilterOn{t: fs.t, cols: cols}
}

// FilterOn filters by named columns.
// TODO(twoodwark): this is needlessly eager with respect to all underlying columns.
type FilterOn struct {
	t    *Table
	cols []ColumnName
}

// TODO
// func (o *FilterOn) Not() *FilterOn {}
// func (o *FilterOn) Null() (*Table, error) {} // all null

func (o *FilterOn) checkSchema(colsType reflect.Type, assertions ...SchemaAssertion) error {
	// eager schema check
	filterSubject, err := o.t.Project(o.cols...)
	if err != nil {
		return errors.Wrapf(err, "can't filter columns %+v", o.cols)
	}
	// assert columns assignable
	if colsType != nil {
		for _, col := range filterSubject.Schema().Columns {
			if !col.Type.AssignableTo(colsType) {
				return errors.Errorf("column %s is not assignable to type %v", col.Name, colsType)
			}
		}
	}
	// checking assertions
	return filterSubject.Schema().Check(assertions...)
}

// MatchRow implements a filter on entire table rows.
type MatchRow func(r Row) bool

/*
 * generic filter guts
 */

// rawMatch is an internal interface for matching function
type rawMatch func(r raw) bool

// the filtered series share an underlying iterator cache
type filterState struct {
	// matcher determines when to return the row.
	// TODO  we don't always need to read the whole Row. colVals do not need to
	// be updated for lazy columns when we already know we matched false
	// (assuming we are using column matchers and not row matchers).
	matcher rawMatch
	source  *tableIterator
	curr    raw
}

func (st *filterState) Next() bool {
	for st.source.Next() {
		if st.isMatch() {
			return true
		}
	}
	return false
}

func (st *filterState) isMatch() bool {
	// cache the column values for the underlying, in case they are expensive
	st.curr = st.source.rawValue()
	return st.matcher(st.curr)
}

func (st *filterState) Value(colIndex int) interface{} {
	return st.curr[colIndex]
}

// compose the matchRow filter into all the series
// matchRow is 'func() rawMatch' (not just 'rawMatch') in order to allow stateful filters
func filterTable(matchGen func() rawMatch, table *Table) *Table {
	newTable := newFromTable(table, table.sortKey...)
	group := newSeriesGroup(func() seriesGroupStateImpl {
		return &filterState{
			matcher: matchGen(),
			source:  newTableIterator(table.series),
		}
	})
	for i, wrappedSeries := range table.series {
		newTable.series[i] = &Series{
			typ:  wrappedSeries.typ,
			col:  wrappedSeries.col,
			read: group.read(i),
			meta: &filteredSeriesMeta{wrapped: wrappedSeries.meta},
		}
	}
	return newTable
}

// filtered series metadata
type filteredSeriesMeta struct {
	wrapped seriesMeta
}

func (m *filteredSeriesMeta) IsMaterialized() bool { return false }

func (m *filteredSeriesMeta) ExactSize() int {
	return -1
}

func (m *filteredSeriesMeta) MaxSize() int {
	if b, ok := m.wrapped.(boundedMeta); ok {
		return b.MaxSize()
	}
	return -1
}

// MatchInterface implements a filter on interface{} columns.  Note that this can receive nil values.
type MatchInterface func(...interface{}) bool

// Interface matches the named column values as interface{} arguments, including nil.
// If given any SchemaAssertions, they are called now and may have side effects.
func (o *FilterOn) Interface(fn MatchInterface, assertions ...SchemaAssertion) (*Table, error) {
	if err := o.checkSchema(nil, assertions...); err != nil {
		return nil, errors.Wrapf(err, "can't filter %+v with %+v", o.t, fn)
	}

	projection := mustNewProjection(o.t.schema, o.cols...)

	matchGen := func() rawMatch {
		return func(r raw) bool {
			matchVals := r.project(projection)
			return fn(matchVals...)
		}
	}

	return filterTable(matchGen, o.t), nil
}

// Interface matches the named column values as interface{} arguments.
func (o *MustFilterOn) Interface(m MatchInterface, assertions ...SchemaAssertion) *Table {
	t, err := o.FilterOn.Interface(m, assertions...)
	handle(err)
	return t
}

/*
 * concrete filters
 */

// Eq retuns a function matching the selected column(s) where equal to expected
// value(s), after any required type conversion.
func Eq(expected ...interface{}) (MatchInterface, SchemaAssertion) {
	assertion := &eq{expected: expected, converted: make([]interface{}, len(expected))}
	return func(v ...interface{}) bool {
		return reflect.DeepEqual(v, assertion.converted)
	}, assertion.CheckSchema
}

type eq struct {
	expected, converted []interface{}
}

// TODO Eq specialization methods to more efficiently filter known scalar types (?)

// CheckSchema converts expected values, as a side effect
func (w *eq) CheckSchema(schema Schema) error {
	if schema.NumColumns() != len(w.expected) {
		return fmt.Errorf("Eq: %d column(s), to equal %d value(s) %+v", schema.NumColumns(), len(w.expected), w.expected)
	}
	for i, c := range schema.Columns {
		e := w.expected[i]
		if e == nil {
			continue
		}
		// convert to the column type
		val := reflect.ValueOf(e)

		if !val.Type().ConvertibleTo(c.Type) {
			return fmt.Errorf("Eq: inconvertible type for %s==%v: %+v to %+v", c.Name, e, val.Type(), c.Type)
		}
		w.converted[i] = val.Convert(c.Type).Interface()
	}
	w.expected = nil
	return nil
}
