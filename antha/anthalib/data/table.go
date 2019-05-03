package data

//go:generate python gen.py
import (
	"github.com/antha-lang/antha/utils"
	"math"
	"reflect"

	"github.com/pkg/errors"
)

// Table is an immutable container of Series.
type Table struct {
	series  []*Series
	schema  Schema
	sortKey Key
	// this must return Row
	read func([]*Series) *tableIterator
}

// NewTable constructs a Table from the given Series.  If at least one of the
// given Series is bounded, then the returned Table is bounded at that size.
// If multiple different-sized bounded series are provided, we panic.
func NewTable(series ...*Series) *Table {
	exact, _, sizeErr := seriesSize(series)
	if sizeErr != nil {
		panic(errors.Wrapf(sizeErr, "cannot construct a table from different-sized series (should all have size %d or be unbounded)", exact))
	}
	return newFromSeries(series)
}

// newFromTable creates a new table pointing to the same series and read fn, and sets the
// sort key.
func newFromTable(t *Table, key ...ColumnKey) *Table {
	s := make([]*Series, len(t.series))
	copy(s, t.series)
	newT := newFromSeries(s, key...)
	newT.read = t.read // TODO this doesn't optimize anything ? remove the read member
	return newT
}

// newFromSeries creates a new table pointing to the existing series, and sets the sort key.
// Does not enforce size invariants!
func newFromSeries(series []*Series, key ...ColumnKey) *Table {
	return &Table{
		series:  series,
		schema:  *newSchema(series),
		sortKey: key,
		read:    newTableIterator,
	}
}

// newFromSchema partially creates a new table based on the known schema and key
// the new table has schema, sortKey, series (but without iterators and metadata) and a generic read function
func newFromSchema(schema *Schema, key ...ColumnKey) *Table {
	series := make([]*Series, len(schema.Columns))
	for i, column := range schema.Columns {
		series[i] = &Series{
			col: column.Name,
			typ: column.Type,
		}
	}
	return &Table{
		series:  series,
		schema:  *schema,
		sortKey: key,
		read:    newTableIterator,
	}
}

// newFromRows creates a table from its materialized representation.
func newFromRows(rows Rows, key ...ColumnKey) (*Table, error) {
	builder, err := NewTableBuilder(rows.Schema.Columns)
	if err != nil {
		return nil, err
	}
	for _, row := range rows.Data {
		builder.Append(row.Interface())
	}
	t := builder.Build()
	return newFromSeries(t.series, key...), nil
}

// Schema returns the type information for the Table.
func (t *Table) Schema() Schema {
	return t.schema
}

// seriesByName returns series by its name.
func (t *Table) seriesByName(col ColumnName) (*Series, error) {
	index, err := t.schema.ColIndex(col)
	if err != nil {
		return nil, err
	}
	return t.series[index], nil
}

// IterAll iterates over the entire table, no buffer.
// Use when ranging over all rows is required.
func (t *Table) IterAll() <-chan Row {
	rows, _ := t.Iter()
	return rows
}

// Iter iterates over the table, no buffer.
// call done() to release resources after a partial read.
func (t *Table) Iter() (rows <-chan Row, done func()) {
	channel := make(chan Row)
	iter := t.read(t.series)
	control := make(chan struct{}, 1)
	done = func() {
		control <- struct{}{}
	}
	go func() {
		defer close(channel)
		for iter.Next() {
			row := iter.Value()
			select {
			case <-control:
				return
			case channel <- row:
				// do nothing
			}
		}
	}()
	return channel, done
}

// ToRows materializes all table data into a Rows object.
// This is useful if you want to print the table data.
func (t *Table) ToRows() Rows {
	iter := t.read(t.series)
	rr := Rows{Schema: t.Schema()}
	for iter.Next() {
		rr.Data = append(rr.Data, iter.Value())
	}
	return rr
}

// Slice is a lazy subset of records between the start index and the end
// (exclusive).  Unlike go slices, if the end index is out of range then fewer
// records are returned rather than receiving an error.
func (t *Table) Slice(start, end Index) *Table {
	return sliceTable(t, start, end)
}

// Head is a lazy subset of the first count records (but may return fewer).
func (t *Table) Head(count int) *Table {
	return t.Slice(0, Index(count))
}

// Sort produces a Table sorted by the columns defined by the Key.
// This is eager, not lazy; it materializes the whole table.
func (t *Table) Sort(key Key) (*Table, error) {
	return sortTableByKey(t, key)
}

// SortFunc should return true if r1 is less than r2.
type SortFunc func(r1 Row, r2 Row) bool

// SortByFunc sorts a table by an arbitrary user-defined function.
// This is eager, not lazy; it materializes the whole table.
// Slower than Table.Sort(). For better performance, use Table.Sort() instead (+ extension if needed).
func (t *Table) SortByFunc(f SortFunc) (*Table, error) {
	return sortTableByFunc(t, f)
}

// Equal is true if the other table has the same schema (in the same order)
// and exactly equal series values.
func (t *Table) Equal(other *Table) bool {
	if t == other {
		return true
	}
	schema1 := t.Schema()
	schema2 := other.Schema()
	if !schema1.Equal(schema2) {
		return false
	}
	// TODO compare tables' known bounded length

	// TODO if table series are identical we can shortcut the iteration
	iter1, done1 := t.Iter()
	iter2, done2 := other.Iter()
	defer done1()
	defer done2()
	for {
		r1, more1 := <-iter1
		r2, more2 := <-iter2
		if more1 != more2 || !reflect.DeepEqual(r1.Values(), r2.Values()) {
			return false
		}
		if !more1 {
			break
		}
	}

	return true
	// TODO since we are iterating over possibly identical series we might optimize by sharing the iterator cache
}

// Size returns the number of rows contained in the table, if known.
// Unbounded, or lazy (non materialized) tables can return -1.
func (t *Table) Size() int {
	exact, _, _ := seriesSize(t.series)
	return exact
}

// find the smallest ext-sized series.
// returns error only to indicate jagged bounded input, but not to shortcut calculation
func seriesSize(series []*Series) (exact, max int, err error) {
	if len(series) == 0 {
		return
	}
	errs := make(utils.ErrorSlice, 0)
	const unknownSize = math.MaxInt64
	max = unknownSize
	exact = unknownSize
	countUnbounded := 0
	for _, ser := range series {
		if b, ok := ser.meta.(boundedMeta); ok {
			seriesMaxSize := b.MaxSize()
			seriesExactSize := b.ExactSize()
			if seriesExactSize >= 0 {
				// if we have more than 1 differently-sized input then this is an error
				// (but caller might choose to ignore such errors)
				if exact != unknownSize && exact != seriesExactSize {
					errs = append(errs, errors.Errorf("jagged input: size=%d for column %q",
						seriesExactSize, ser.col,
					))
				}
				// regardless, adjust the actual smallest size, for callers that are combining tables
				if seriesExactSize < exact {
					exact = seriesExactSize
				}
			}
			if seriesMaxSize < max {
				max = seriesMaxSize
			}
		} else {
			countUnbounded++
		}
	}
	if countUnbounded > 0 || exact == unknownSize {
		exact = -1
	} else if max < exact {
		// TODO test coverage ?
		exact = unknownSize
	}
	err = errs.Pack()
	return
}

func isBounded(series []*Series) bool {
	// a table is bounded if at least one of its series is
	for _, s := range series {
		if _, ok := s.meta.(boundedMeta); ok {
			return true
		}
	}
	return len(series) == 0
}

// Cache converts a lazy table to one that is fully materialized.
func (t *Table) Cache() (*Table, error) {
	return cacheTable(t, arrowSeries, /*Arrow series if possible*/
		false, /*Native series otherwise*/
		false /*don't copy columns which are already materialized*/)
}

// Copy: we aren't exposing this method because it is useless - since our tables are immutable (at least officially).
// (however, copyTable is used internally for implementing in-place Sort)
//func (t *Table) Copy() (*Table, error) {
//	return copyTable(t, ...)
//}

// DropNullColumns filters out columns with all/any row null
// TODO
// func (t *Table) DropNullColumns(all bool) *Table {
// 	return nil
// }

// DropNull filters out rows with all/any col null
// TODO
// func (t *Table) DropNull(all bool) *Table {
// 	return nil
// }

// Project reorders and/or takes a subset of columns. On duplicate columns, only
// the first so named is taken. Returns error, and nil table, if any column is
// missing.
func (t *Table) Project(columns ...ColumnName) (*Table, error) {
	s := make([]*Series, len(columns))
	for i, columnName := range columns {
		series, err := t.seriesByName(columnName)
		if err != nil {
			return nil, errors.Wrapf(err, "when projecting %v", t.Schema())
		}
		s[i] = series
	}
	// TODO rearrange key
	return NewTable(s...), nil
}

// ProjectAllBut discards the named columns, which may not exist in the schema.
func (t *Table) ProjectAllBut(columns ...ColumnName) *Table {
	byName := map[ColumnName]struct{}{}
	for _, n := range columns {
		byName[n] = struct{}{}
	}
	s := []*Series{}
	for _, ser := range t.series {
		if _, found := byName[ser.col]; !found {
			s = append(s, ser)
		}
	}
	return NewTable(s...) // TODO set key to subkey
}

// Rename updates all columns of the old name to the new name.
// Does nothing if none match in the schema.
func (t *Table) Rename(old, new ColumnName) *Table {
	s := make([]*Series, len(t.series))
	for i, ser := range t.series {
		if ser.col == old {
			ser = &Series{
				col:  new,
				typ:  ser.typ,
				meta: ser.meta,
				read: ser.read,
			}
		}
		s[i] = ser
	}
	// TODO rename key column
	return NewTable(s...)
}

// Convert lazily converts all columns of the given name to the assigned type.
// Returns non-nil error (and nil Table) if any column is not convertible. Note
// that if no column name matches, the same table is returned.
func (t *Table) Convert(col ColumnName, typ reflect.Type) (*Table, error) {
	newT := newFromTable(t, t.sortKey...)
	converted := false
	conv := &conversion{newType: typ, Table: t}
	for i, ser := range t.series {
		var err error
		if ser.col == col && !typesEqual(ser.typ, typ) {
			newT.series[i], err = conv.convert(ser)
			converted = true
			if err != nil {
				return nil, errors.Wrapf(err, "cannot convert column %d (%q)", i, col)
			}
		}
	}
	if !converted {
		return t, nil
	}
	// types but not sort key are different.  TODO: unless natural order is different for the converted type?
	newT.schema = *newSchema(newT.series)
	return newT, nil
}

// Filter selects some records lazily.  Use the returned object to construct a
// derived table.
func (t *Table) Filter() *FilterSelection {
	return &FilterSelection{t}
}

// Distinct filters the table retaining only unique rows.
// Use the returned object to construct a derived table.
func (t *Table) Distinct() *DistinctSelection {
	return &DistinctSelection{t}
}

// Extend adds a column called newCol by applying a function.  Use the returned
// object to construct a derived table.
func (t *Table) Extend(newCol ColumnName) *Extension {
	return &Extension{newCol: newCol, t: t}
}

// Update replaces the values in the column 'col'.  Use the returned
// object to construct a derived table.
func (t *Table) Update(col ColumnName) *UpdateSelection {
	return &UpdateSelection{t: t, col: col}
}

// Pivot takes a "narrow" table containing "column names" and "column values", for instance:
//
//  | |    Key|    Pivot|    Value|
//  | |keyType|   string|valueType|
//  -------------------------------
//  |1|   key1|"column1"|   value1|
//  |2|   key1|"column2"|   value2|
//  |3|   key1|"column3"|   value3|
//  |4|   key2|"column1"|   value4|
//  |5|   key2|"column2"|   value5|
//
// and transforms it into a "wide" table, for instance:
//
//  | |    Key|  column1|  column2|  column3|
//  | |keyType|valueType|valueType|valueType|
//  -----------------------------------------
//  |1|   key1|   value1|   value2|   value3|
//  |2|   key2|   value4|   value5|    <nil>|
func (t *Table) Pivot() *PivotSelection {
	return &PivotSelection{table: t}
}

// Join performs a "horizontal" join of two tables by some key. E.g. if the left table is
//
//  | |LeftKeyColumn|OtherLeftColumn|
//  | |      keyType|  otherLeftType|
//  ---------------------------------
//  |1|         key1| someLeftValue1|
//  |2|         key2| someLeftValue2|
//  |3|         key2| someLeftValue3|
//  |4|         key3| someLeftValue4|
//
// and the right table is
//
//  | |RightKeyColumn|OtherRightColumn|
//  | |       keyType|  otherRightType|
//  -----------------------------------
//  |1|          key2| someRightValue1|
//  |2|          key3| someRightValue2|
//
// then if we inner join them on LeftKeyColumn and RightKeyColumn columns correspondingly
//
//  resultTable, err := leftTable.Join().On("LeftKeyColumn").Inner(rightTable, "RightKeyColumn")
//
// the result would be:
//
//  | |LeftKeyColumn|OtherLeftColumn|RightKeyColumn|OtherRightColumn|
//  | |      keyType|  otherLeftType|       keyType|  otherRightType|
//  -----------------------------------------------------------------
//  |2|         key2| someLeftValue2|          key2| someRightValue1|
//  |3|         key2| someLeftValue3|          key2| someRightValue1|
//  |4|         key3| someLeftValue4|          key3| someRightValue2|
func (t *Table) Join() *JoinSelection {
	return &JoinSelection{t: t}
}

// Append concatenates two tables vertically.
// Use the returned object to set append mode.
func (t *Table) Append(other *Table) *AppendSelection {
	return &AppendSelection{
		tables: []*Table{t, other},
	}
}

// Append concatenates several tables vertically.
// Use the returned object to set append mode.
func Append(t ...*Table) *AppendSelection {
	return &AppendSelection{
		tables: t,
	}
}

// Foreach eagerly iterates over a table using a user-defined function.
// May be useful for:
//
// - reading table contents into some user-defined data struct (e.g. a map)
//
// - computing some scalar aggregate (e.g. MIN(column))
//
// Use the returned object to specify the columns and the function.
func (t *Table) Foreach() *ForeachSelection {
	return &ForeachSelection{t: t}
}

// ForeachKey eagerly iterates over a table grouping it by a given key: for every single key value, Foreach().By()
// calls the callback and supplies it with the key value and the "partition" table (the table containing other values
// from corresponding rows of the source table).
//
// For the time being, it is implemented on the top of Sort() - therefore only orderable keys are supported (currently
// only supported types are considered orderable).
func (t *Table) ForeachKey(key ...ColumnName) *ForeachKeySelection {
	return &ForeachKeySelection{
		t:   t,
		key: key,
	}
}
