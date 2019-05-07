package data

import (
	"reflect"

	"github.com/pkg/errors"
)

/*
 * pivot interfaces
 */

// PivotSelection contains a table to apply pivot.
type PivotSelection struct {
	table *Table
}

// Key set key columns names for the ongoing pivot operation
func (ps *PivotSelection) Key(key ...ColumnName) *PivotKey {
	return &PivotKey{table: ps.table, key: key}
}

// PivotKey contains a table to apply pivot and a key.
type PivotKey struct {
	table *Table
	key   []ColumnName
}

// Columns set pivot and value column names for the ongoing pivot operation
func (pk *PivotKey) Columns(pivot ColumnName, value ColumnName) (*Table, error) {
	return pivotTable(pk.table, pk.key, pivot, value)
}

/*
 * pivot internals
 */

func pivotTable(table *Table, key []ColumnName, pivot ColumnName, value ColumnName) (*Table, error) {
	// checking columns assertions
	if err := checkColumns(&table.schema, key, pivot, value); err != nil {
		return nil, err
	}

	// scanning pivot column and retrieving wide table columns list
	wideColumns, wideColumnType, err := getWideColumns(table, pivot, value)
	if err != nil {
		return nil, err
	}

	// partially creating a output table (for the time being, without series iterators and metadata)
	outputTable := newPivotedTable(table, key, pivot, value, wideColumns, wideColumnType)

	// the output table max size: in worst case, equals to the source table max size
	_, maxSize, _ := seriesSize(table.series)

	group := newSeriesGroup(func() seriesGroupStateImpl {
		// an intermediate table: contains projection of columns of the source table
		// we are interested in (key columns, pivot column, value column), sorted by key columns
		intermediate := table.Must().Project(allInputColumns(key, pivot, value)...).Must().Sort(outputTable.sortKey)
		// creating common state
		return newPivotState(intermediate, wideColumns)
	})

	// creating an iterator generator and metadata for each series
	for i, series := range outputTable.series {
		series.read = group.read(i)
		series.meta = &pivotedSeriesMeta{maxSize: maxSize}
	}

	return outputTable, nil
}

// checkColumns checks column assertions
func checkColumns(schema *Schema, key []ColumnName, pivot ColumnName, value ColumnName) error {
	// all columns must exist in the input table
	if err := schema.CheckColumnsExist(allInputColumns(key, pivot, value)...); err != nil {
		return errors.Wrap(err, "pivot columns lookup")
	}

	// pivot column must be of string type
	if schema.MustCol(pivot).Type != reflect.TypeOf("") {
		return errors.Errorf("pivot column '%s' must be a string column", pivot)
	}

	return nil
}

// allInputColumns gets all input table columns as a single slice
func allInputColumns(key []ColumnName, pivot ColumnName, value ColumnName) []ColumnName {
	return append(append([]ColumnName{}, key...), pivot, value)
}

// gets wide column names by indexing unique values in the pivot column
func getWideColumns(table *Table, pivot ColumnName, value ColumnName) ([]ColumnName, reflect.Type, error) {
	wideColumnList := table.Must().Project(pivot).Must().Distinct().On(pivot)
	wideColumnIter := wideColumnList.read(wideColumnList.series)
	wideColumnNames := []ColumnName{}
	for wideColumnIter.Next() {
		raw := wideColumnIter.rawValue()[0]
		if raw == nil {
			return nil, nil, errors.Errorf("pivot column '%s' must not contain nulls", pivot)
		}
		wideColumnNames = append(wideColumnNames, ColumnName(raw.(string)))
	}
	return wideColumnNames, table.schema.MustCol(value).Type, nil
}

// newPivotedTable partially creates a pivot output table (for the time being, without series iterators and metadata)
func newPivotedTable(table *Table, key []ColumnName, pivot ColumnName, value ColumnName,
	wideColumns []ColumnName, wideColumnsType reflect.Type) *Table {
	// pivoted table columns: key columns and wide columns
	columns := make([]Column, 0, len(key)+len(wideColumns))
	for _, k := range key {
		columns = append(columns, table.schema.MustCol(k))
	}
	for _, wideColumn := range wideColumns {
		columns = append(columns, Column{Name: wideColumn, Type: wideColumnsType})
	}

	// pivoted table key
	outputKey := make([]ColumnKey, len(key))
	for i, k := range key {
		outputKey[i] = ColumnKey{Column: k, Asc: true}
	}

	return newFromSchema(NewSchema(columns), outputKey...)
}

// pivotState stores the common state of pivoted table series iterators
type pivotState struct {
	numKeys       int                // number of the key columns
	wideColByName map[ColumnName]int // wide column name -> wide column index in the output table

	intermediate *tableIterator // iterator over the intermediate sorted table

	currRow []interface{} // current row of the output table: |key1|...|keyN|value1|...|valueM|
	nextRow []interface{} // next row of the output table
}

func newPivotState(intermediate *Table, wideColumns []ColumnName) *pivotState {
	numKeys := len(intermediate.series) - 2

	wideColByName := map[ColumnName]int{}
	for i, wideColumn := range wideColumns {
		wideColByName[wideColumn] = numKeys + i
	}

	return &pivotState{
		numKeys:       numKeys,
		wideColByName: wideColByName,
		intermediate:  newTableIterator(intermediate.series),
	}
}

func (ps *pivotState) Next() bool {
	// cleaning the current row
	ps.currRow = nil
	// if previous iteration has started filling the next row, then moving it to the current row
	if ps.nextRow != nil {
		ps.currRow = ps.nextRow
		ps.nextRow = nil
	}

	// iterating over the intermediate table
	for ps.intermediate.Next() {
		// a row of the intertmediate table: |key1|...|keyN|column name|value|
		key, columnName, columnValue := ps.parseIntermediateRow()

		if ps.currRow == nil {
			// it is first iteration => initializing current output row with current key
			ps.currRow = ps.newOutputRow(key)
			ps.setOutputRowValue(ps.currRow, columnName, columnValue)
		} else if ps.keysEqual(key, ps.currRow[:ps.numKeys]) {
			// the row has the same key as the previous row => continuing filling the same output row
			ps.setOutputRowValue(ps.currRow, columnName, columnValue)
		} else {
			// intermediate table iterator has reached a new key => saving it into ps.nextRow and stopping iteration
			ps.nextRow = ps.newOutputRow(key)
			ps.setOutputRowValue(ps.nextRow, columnName, columnValue)
			break
		}
	}

	return ps.currRow != nil
}

// parseRow extracts keys, column name and column value from a row from an intermediate table
func (ps *pivotState) parseIntermediateRow() ([]interface{}, ColumnName, interface{}) {
	row := ps.intermediate.rawValue()
	return row[:ps.numKeys], ColumnName(row[ps.numKeys].(string)), row[ps.numKeys+1]
}

// newOutputRow creates a new output table row with given keys
func (ps *pivotState) newOutputRow(keys []interface{}) []interface{} {
	outputRow := make([]interface{}, ps.numKeys+len(ps.wideColByName))
	copy(outputRow[:ps.numKeys], keys)
	return outputRow
}

// setOutputRowValue find a column in the output row by its name and sets its value
// TODO: now we allow multiple source table entries setting the same output table cell; need to fix this
func (ps *pivotState) setOutputRowValue(outputRow []interface{}, columnName ColumnName, columnValue interface{}) {
	outputRow[ps.wideColByName[columnName]] = columnValue
}

// keysEqual compares two keys
// TODO: specialize this after Table.Distinct is specialized
func (ps *pivotState) keysEqual(keys1 []interface{}, keys2 []interface{}) bool {
	return reflect.DeepEqual(keys1, keys2)
}

// Value reads the cached column value
func (ps *pivotState) Value(colIndex int) interface{} {
	return ps.currRow[colIndex]
}

// metadata for pivoted table series (both key columns and wide columns)
type pivotedSeriesMeta struct {
	// Theoretically speaking, we can determine exact size of the pivoted table - but with some overhead. To do it, we should scan
	// the source table key columns immediately when creating pivoted table (along with the pivot column which is currently scanned).
	// However, to do it in optimal way, we should scan two groups of columns (keys and pivot) in one pass, which is impossible with current Distinct API
	maxSize int
}

func (m *pivotedSeriesMeta) IsMaterialized() bool { return false }
func (m *pivotedSeriesMeta) ExactSize() int       { return -1 }
func (m *pivotedSeriesMeta) MaxSize() int         { return m.maxSize }
