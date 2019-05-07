package data

import (
	"github.com/pkg/errors"
)

// AppendMany concatenates multiple tables vertically.
// The tables schemas must be identical (except columns names; the output table inherits columns names from the first input table).
// Empty input tables list is considered as an error.
func AppendMany(tables ...*Table) (*Table, error) {
	// checking inputs
	if err := checkInputTablesForAppend(tables); err != nil {
		return nil, err
	}

	// optimization: if some of the input tables are themselves produced by `Append`, then replacing them
	// with their source tables - in order to simplify iteration over the resulting table
	tables = flattenAppendSources(tables)

	// append series sizes
	exactSize, maxSize := calcAppendSizes(tables)

	// series group
	group := newSeriesGroup(func() seriesGroupStateImpl {
		return &appendState{
			sources:     tables,
			sourceIndex: -1,
		}
	})

	// creating the output table series: they inherit the schema from the first input table
	newSeries := make([]*Series, len(tables[0].series))
	for i, series := range tables[0].series {
		newSeries[i] = &Series{
			col:  series.col,
			typ:  series.typ,
			read: group.read(i),
			meta: &appendSeriesMeta{
				exactSize: exactSize,
				maxSize:   maxSize,
				group:     group,
				sources:   tables,
			},
		}
	}

	return NewTable(newSeries...), nil
}

// Checks that source tables schemas are equal.
func checkInputTablesForAppend(tables []*Table) error {
	// at least one input table must be supplied
	if len(tables) == 0 {
		return errors.New("append: an empty list of tables is not supported.")
	}
	// tables schemas must be identical (except column names)
	for i := 1; i < len(tables); i++ {
		if err := compareSchemasForAppend(tables[0].schema, tables[i].schema); err != nil {
			return err
		}
	}
	return nil
}

// Compares schemas of two tables for Append, returns an error if they are different.
// Schemas are compared by types only: columns names are not compared (Append uses columns names from the frst table in the list).
func compareSchemasForAppend(schema1 Schema, schema2 Schema) error {
	numColumns1, numColumns2 := schema1.NumColumns(), schema2.NumColumns()
	if numColumns1 != numColumns2 {
		return errors.Errorf("append: tables having different numbers of columns encontered (%d and %d)", numColumns1, numColumns2)
	}
	for i := range schema1.Columns {
		type1, type2 := schema1.Columns[i].Type, schema2.Columns[i].Type
		if type1 != type2 {
			return errors.Errorf("append: tables columns #%d having different types (%v and %v)", i, type1, type2)
		}
	}
	return nil
}

// If some of the Append source tables are themselves produced by `Append`, then removing intermediate tables
// - in order to simplify iteration.
func flattenAppendSources(sourceTables []*Table) []*Table {
	flattened := make([]*Table, 0, len(sourceTables))
	for _, t := range sourceTables {
		if flattenedT, ok := extractAppendSources(t); ok {
			flattened = append(flattened, flattenedT...)
		} else {
			flattened = append(flattened, t)
		}
	}
	return flattened
}

// If the table is itself a result of Append, then extracts source tables. Otherwise, returns false.
func extractAppendSources(table *Table) ([]*Table, bool) {
	for _, series := range table.series {
		// if some series is not Append series, the table is not the result of a single Append
		meta, ok := series.meta.(*appendSeriesMeta)
		if !ok {
			return nil, false
		}
		// if some of Append series do not share the same appendTableMeta, the table is not the result of a single Append
		if meta.group != table.series[0].meta.(*appendSeriesMeta).group {
			return nil, false
		}
	}
	return table.series[0].meta.(*appendSeriesMeta).sources, true
}

func calcAppendSizes(sourceTables []*Table) (exactSize int, maxSize int) {
	for _, t := range sourceTables {
		exactSizeT, maxSizeT, err := seriesSize(t.series)
		if err != nil {
			panic(errors.Wrap(err, "SHOULD NOT HAPPEN"))
		}
		exactSize = addSizes(exactSize, exactSizeT)
		maxSize = addSizes(maxSize, maxSizeT)
	}
	return
}

// adds two sizes; if one of them is undefined, the sum is undefined too
func addSizes(size1 int, size2 int) int {
	const undefinedSize = -1
	if size1 == undefinedSize || size2 == undefinedSize {
		return undefinedSize
	}
	return size1 + size2
}

// appendTables resulting series metadata
type appendSeriesMeta struct {
	exactSize int
	maxSize   int
	// needed for optimization purposes only (extractAppendSources)
	group   *seriesGroup
	sources []*Table
}

func (m *appendSeriesMeta) IsMaterialized() bool { return false }
func (m *appendSeriesMeta) ExactSize() int       { return m.exactSize }
func (m *appendSeriesMeta) MaxSize() int         { return m.maxSize }

// Common state of several Append series. They share the source table iterator (including series iterators cache).
type appendState struct {
	sources     []*Table       // source tables
	sourceIndex int            // current source table
	sourceIter  *tableIterator // current source table iterator
}

func (st *appendState) Next() bool {
	for {
		// first trying to advance the current table iterator
		if st.sourceIndex != -1 && st.sourceIter.Next() {
			return true
		}
		// moving to the next table if possible
		if st.sourceIndex+1 == len(st.sources) {
			return false
		}
		st.sourceIndex++
		st.sourceIter = newTableIterator(st.sources[st.sourceIndex].series)
	}
}

func (st *appendState) Value(colIndex int) interface{} {
	return st.sourceIter.colReader[colIndex].Value()
}
