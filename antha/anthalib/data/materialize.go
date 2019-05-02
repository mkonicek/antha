package data

import (
	"reflect"

	"github.com/pkg/errors"
)

/*
 * API for building Tables
 */

// TODO: do we really need an API in this form at all? Think about replacing this API with 'lazy series' API (this will require exposing iterators)

// TableBuilder builds a materialized Table from external source row by row
type TableBuilder struct {
	seriesBuilders []seriesBuilder
}

// NewTableBuilder creates a new default table builder by a columns list
func NewTableBuilder(columns []Column) (*TableBuilder, error) {
	// default mode: constructing Arrow series if possible, native series otherwise
	return newTableBuilderExt(columns, arrowSeries, false)
}

// Reserve reserves extra buffer space
func (b *TableBuilder) Reserve(capacity int) {
	for i := range b.seriesBuilders {
		b.seriesBuilders[i].Reserve(capacity)
	}
}

// Append appends a row (nil element denotes null)
func (b *TableBuilder) Append(row []interface{}) {
	if len(row) != len(b.seriesBuilders) {
		panic("insuffecient number of columns in the row")
	}

	// appending a value to each column
	for i := range b.seriesBuilders {
		b.seriesBuilders[i].Append(row[i])
	}
}

// Build builds a table
func (b *TableBuilder) Build() *Table {
	// finishing each series
	series := make([]*Series, len(b.seriesBuilders))
	for i := range series {
		series[i] = b.seriesBuilders[i].Build()
	}

	// constructing a Table
	return NewTable(series...)
}

// materializedType denotes types of materialized series
type materializedType int

const (
	// nativeSeries - native Go slice based series
	nativeSeries materializedType = iota
	// arrowSeries - Apache Arrow based series
	arrowSeries
)

// newTableBuilderExt is an extended version of NewTableBuilder; if desired series builder mode is unavailable, using a fallback
func newTableBuilderExt(columns []Column, mode materializedType, forceMode bool) (*TableBuilder, error) {
	// creating builders for each column
	seriesBuilders := make([]seriesBuilder, len(columns))
	for i, column := range columns {
		seriesBuilder, err := newSeriesBuilderExt(column.Name, column.Type, mode, forceMode)
		if err != nil {
			return nil, err
		}
		seriesBuilders[i] = seriesBuilder
	}

	return &TableBuilder{
		seriesBuilders: seriesBuilders,
	}, nil
}

func (b *TableBuilder) schema() *Schema {
	cols := make([]Column, len(b.seriesBuilders))
	for cIdx, sb := range b.seriesBuilders {
		cols[cIdx] = sb.Column()
	}
	return NewSchema(cols)
}

// seriesBuilder is an interface for building materialized Series from external data source
type seriesBuilder interface {
	// Information about the column being built
	Column() Column

	// Reserve reserves extra buffer space
	Reserve(capacity int)
	// Size returns the number of appended values
	Size() int
	// Append appends a single value; nil denotes a null value
	Append(value interface{})

	// Build constructs Series
	Build() *Series
}

// newSeriesBuilderExt creates a new series builder
func newSeriesBuilderExt(col ColumnName, typ reflect.Type, mode materializedType, forceMode bool) (seriesBuilder, error) {
	seriesBuilder, err := newSeriesBuilder(col, typ, mode)
	if err != nil {
		if forceMode {
			return nil, errors.Wrap(err, "creating series builder")
		}
		// using native series builder as a fallback builder because it should work for arbitrary column types
		seriesBuilder, err = newSeriesBuilder(col, typ, nativeSeries)
		if err != nil {
			panic(errors.Wrap(err, "SHOULD NOT HAPPEN: creating native series builder"))
		}
	}
	return seriesBuilder, nil
}

// a generic (but slow) series builder - currently is implemented for native series only
func newFallbackSeriesBuilder(col ColumnName, typ reflect.Type, mode materializedType) (seriesBuilder, error) {
	switch mode {
	case nativeSeries:
		return newFallbackNativeSeriesBuilder(col, typ), nil
	case arrowSeries:
		return nil, errors.Errorf("unable to create Arrow series of type %+v", typ)
	default:
		panic(errors.Errorf("unknown materialized series type %v", mode))
	}
}

/*
 * API for caching Series and Tables
 */

// cacheTable is a extended internal version of table.Cache which allows
//  1) to choose type of series to cache into (Native/Arrow) and 2) force copying series which are already cached
func cacheTable(t *Table, mode materializedType, forceMode bool, forceCopy bool) (*Table, error) {
	if !isBounded(t.series) {
		return nil, errors.New("unable to materialize unbounded table")
	}

	// choosing series to copy (if forceCopy is set, then copying all the series; otherwise, copying non-materialized ones only)
	seriesToCopy := []*Series{}
	indexesToCopy := []int{}
	for i, ser := range t.series {
		if doCacheSeries(ser, mode, forceMode, forceCopy) {
			seriesToCopy = append(seriesToCopy, ser)
			indexesToCopy = append(indexesToCopy, i)
		}
	}
	tableToCopy := NewTable(seriesToCopy...)

	// copying selected series
	copiedTable, err := copyTable(tableToCopy, mode, forceMode)
	if err != nil {
		return nil, err
	}

	// compose the table again inserting copied series
	newTable := newFromTable(t, t.sortKey...)
	for i, seriesIndex := range indexesToCopy {
		newTable.series[seriesIndex] = copiedTable.series[i]
	}

	return newTable, nil
}

func doCacheSeries(series *Series, mode materializedType, forceMode bool, forceCopy bool) bool {
	// - `forceCopy` means forcing copying all the series regardless their types
	// - if the series is not materialized, it should be copied anyway as well
	if forceCopy || !series.meta.IsMaterialized() {
		return true
	}

	// if we don't care about the resulting materialized series type, there is no need to copy materialized series
	if !forceMode {
		return false
	}

	// if we care about the resulting materialized series type, then materializing series of other types
	switch series.meta.(type) {
	case *nativeSeriesMeta:
		return mode != nativeSeries
	case *arrowSeriesMeta:
		return mode != arrowSeries
	default:
		return true
	}
}

// caches all the columns of t into either native or arrow series
func copyTable(t *Table, mode materializedType, forceMode bool) (*Table, error) {
	if !isBounded(t.series) {
		return nil, errors.New("unable to copy unbounded table")
	}

	// creating table iterator
	tableIter := t.read(t.series)

	// creating series copiers (they try to use specialized iterators if possible)
	seriesCopiers := make([]seriesCopier, len(t.series))
	size := t.Size()
	for i, s := range t.series {
		copier, err := newSeriesCopier(s, tableIter.colReader[i], mode)
		if err != nil {
			if forceMode {
				return nil, errors.Wrap(err, "creating series copier")
			}
			// using native series copier as a fallback copier because it should work for arbitrary column types
			copier, err = newSeriesCopier(s, tableIter.colReader[i], nativeSeries)
			if err != nil {
				panic(errors.Wrap(err, "SHOULD NOT HAPPEN: creating native series copier"))
			}
		}
		if size != -1 {
			// reserving space if possible
			copier.Reserve(size)
		}
		seriesCopiers[i] = copier
	}

	// copying rows itself
	for tableIter.Next() {
		for _, c := range seriesCopiers {
			c.CopyValue()
		}
	}

	// finishing each series
	series := make([]*Series, len(seriesCopiers))
	for i := range series {
		series[i] = seriesCopiers[i].Build()
	}

	// constructing a Table
	return newFromSeries(series, t.sortKey...), nil
}

// an interface for copying series element by element
type seriesCopier interface {
	CopyValue()
	Reserve(capacity int)
	Build() *Series
}

type fallbackSeriesCopier struct {
	seriesBuilder
	iter iterator
}

func newFallbackSeriesCopier(s *Series, iter iterator, mode materializedType) (seriesCopier, error) {
	// using generic source iterator and target builder
	builder, err := newFallbackSeriesBuilder(s.col, s.typ, mode)
	if err != nil {
		return nil, err
	}

	return &fallbackSeriesCopier{
		seriesBuilder: builder,
		iter:          iter,
	}, nil
}

func (c *fallbackSeriesCopier) CopyValue() { c.Append(c.iter.Value()) }
