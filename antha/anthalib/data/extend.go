package data

import (
	"reflect"

	"github.com/pkg/errors"
)

/*
 * calculated columns
 */

// TODO preserve sort key

// Extension is the fluent interface for adding calculated columns.
type Extension struct {
	// the new column to add
	newCol ColumnName
	// the table to extend
	t *Table
}

// By performs extension using the whole table row as input.
// For wide tables this might be inefficient, consider using On(...) instead.
func (e *Extension) By(f func(r Row) interface{}, newType reflect.Type) *Table {
	// TODO: either reflectively infer newType, or assert/verify the f return type
	series := append(append([]*Series(nil), e.t.series...), &Series{
		col:  e.newCol,
		typ:  newType,
		meta: newExtendSeriesMeta(e.t.series, false),
		read: func(cache *seriesIterCache) iterator {
			return &extendRowSeries{f: f, source: e.extensionSource(cache)}
		}},
	)

	newT := newFromSeries(series, e.t.sortKey...)
	return newT
}

// extensionSource is exhausted when the underlying table is. Side effects are
// important to get table cardinality correct without requiring the extension
// column iterator to return false.
func (e *Extension) extensionSource(cache *seriesIterCache) *readRow {
	// virtual table will not be used to advance
	source := &readRow{iteratorCache: cache}
	// go get the series iterators we need from the cache
	source.fill(e.t.series)
	return source
}

type extendRowSeries struct {
	f      func(r Row) interface{}
	source *readRow
}

func (i *extendRowSeries) Next() bool { return true }

func (i *extendRowSeries) Value() interface{} {
	row := i.source.Value()
	v := i.f(row)
	return v
}

// On selects a subset of columns to use as an extension source. If duplicate columns
// exist, the first so named is used.  Note this does not panic yet, even if the
// columns do not exist.  (However subsequent calls to the returned object will
// error.)
func (e *Extension) On(cols ...ColumnName) *ExtendOn {
	return &ExtendOn{extension: e, meta: newExtendSeriesMeta(e.t.series, false), inputCols: cols}
}

// Interface adds a column of an arbitrary type using inputs of arbitrary types.
func (e *ExtendOn) Interface(f func(v ...interface{}) interface{}, newType reflect.Type) (*Table, error) {
	inputs, err := e.inputs(newType)
	if err != nil {
		return nil, err
	}

	// TODO: either reflectively infer newType, or assert/verify the f return type
	series := append(append([]*Series(nil), e.extension.t.series...), &Series{
		col:  e.extension.newCol,
		typ:  newType,
		meta: e.meta,
		read: func(cache *seriesIterCache) iterator {
			colReader := make([]iterator, len(inputs))
			for i, inputCol := range inputs {
				colReader[i] = cache.Ensure(inputCol)
			}
			return &extendInterface{f: f, source: colReader}
		}},
	)

	newT := newFromSeries(series, e.extension.t.sortKey...)
	return newT, nil
}

type extendInterface struct {
	f      func(v ...interface{}) interface{}
	source []iterator
}

func (ei *extendInterface) Next() bool { return true }

func (ei *extendInterface) Value() interface{} {
	args := make([]interface{}, len(ei.source))
	for i, s := range ei.source {
		args[i] = s.Value()
	}
	return ei.f(args...)
}

// Constant adds a constant column with the given value to the table.  This
// column has the dynamic type of the given value.
func (e *Extension) Constant(value interface{}) *Table {
	ext, err := e.ConstantType(value, reflect.TypeOf(value))
	if err != nil {
		panic(err)
	}
	return ext
}

// ConstantType adds a constant column with the given value to the table.  This
// column has the given type (boxing nil for example).  Returns error if a
// non-nil value cannot be converted to the required type.
func (e *Extension) ConstantType(value interface{}, typ reflect.Type) (*Table, error) {
	if value != nil {
		v := reflect.ValueOf(value)
		if !v.Type().ConvertibleTo(typ) {
			return nil, errors.Errorf("value of type %v is not assignable to type %v", v.Type(), typ)
		}
		value = v.Convert(typ).Interface()
	}
	// can be materialized but unbounded
	meta := newExtendSeriesMeta(e.t.series, true)

	ser := &Series{
		col:  e.newCol,
		typ:  typ,
		meta: meta,
		read: func(cache *seriesIterCache) iterator {
			e.extensionSource(cache)
			return &constIterator{value: value}
		},
	}

	return NewTable(append(e.t.series, ser)), nil
}

// NewConstantSeries returns an unbounded repetition of the same value, using
// the dynamic type of the given value.
func NewConstantSeries(col ColumnName, value interface{}) *Series {
	iter := &constIterator{value: value}
	return &Series{
		col:  col,
		meta: &combinedSeriesMeta{isMaterialized: true},
		typ:  reflect.TypeOf(value),
		read: func(_ *seriesIterCache) iterator {
			return iter
		},
	}
}

type constIterator struct{ value interface{} }

func (i *constIterator) Next() bool         { return true }
func (i *constIterator) Value() interface{} { return i.value }

// extension is bounded if not all underlying series are unbounded
func newExtendSeriesMeta(series []*Series, isMaterialized bool) seriesMeta {
	m := &combinedSeriesMeta{isMaterialized}
	if isBounded(series) {
		b := &boundedCombinedSeriesMeta{combinedSeriesMeta: m}
		b.exact, b.max, _ = seriesSize(series)
		return b
	}
	return m
}

type combinedSeriesMeta struct{ isMaterialized bool }

func (m *combinedSeriesMeta) IsMaterialized() bool { return m.isMaterialized }

type boundedCombinedSeriesMeta struct {
	*combinedSeriesMeta
	exact, max int
}

func (m *boundedCombinedSeriesMeta) ExactSize() int { return m.exact }
func (m *boundedCombinedSeriesMeta) MaxSize() int   { return m.max }

var _ boundedMeta = (*boundedCombinedSeriesMeta)(nil)

// ExtendOn enables extensions using specific column values as function inputs
type ExtendOn struct {
	meta      seriesMeta
	extension *Extension
	inputCols []ColumnName
}

func (e *ExtendOn) inputs(asType reflect.Type) ([]*Series, error) {
	inputs := make([]*Series, len(e.inputCols))
	for i, c := range e.inputCols {
		colIndex, err := e.extension.t.schema.ColIndex(c)
		if err != nil {
			return nil, errors.Wrapf(err, "extending new column %q", e.extension.newCol)
		}
		ser := e.extension.t.series[colIndex]
		if asType != nil {
			if err := ser.assignableTo(asType); err != nil {
				return nil, errors.Wrapf(err, "when extending new column %q", e.extension.newCol)
			}
		}
		inputs[i] = ser
	}
	return inputs, nil
}
