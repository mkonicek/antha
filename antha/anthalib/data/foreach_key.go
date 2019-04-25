package data

import (
	"reflect"

	"github.com/pkg/errors"
)

// ForeachKeySelection encapsulates a table to iterate with ForeachKey and a key to iterate with.
type ForeachKeySelection struct {
	t   *Table
	key []ColumnName
}

// By performs ForeachKey on the whole table rows.
// For wide tables this might be inefficient, consider using Project before By.
func (fks *ForeachKeySelection) By(fn func(key Row, partition *Table)) error {
	// For now, iteration by key is implemented via
	// 1) Sort() and
	// 2) directly slicing the series of the sorted table (i.e. using the fact that Sort() currently returns native series)
	// Thus now iteration by key can be done on sortable columns only (though it is not a fundamental limitation).

	// sorting the source table by the key columns
	sorted, err := fks.sort()
	if err != nil {
		return errors.Wrap(err, "iteration by key")
	}

	// splitting the sorted table into the key table and the rest of the table (which is going to be divided into partitions)
	keyTable, partitionsTable := fks.splitSortedTable(sorted)

	// iterating over the key table to determine partitions borders
	index := Index(-1)
	partitionStart := Index(-1)
	partitionKey := Row{}

	keyTable.Foreach().By(func(key Row) {
		index++
		if index == 0 {
			partitionStart = 0
			partitionKey = key
		} else if !reflect.DeepEqual(key.values, partitionKey) {
			fn(partitionKey, sliceNativeTable(partitionsTable, partitionStart, index))
			partitionStart = index
			partitionKey = key
		}
	})

	// finishing the last partition (provided that the source table is not empty)
	if partitionStart != -1 {
		fn(partitionKey, sliceNativeTable(partitionsTable, partitionStart, index+1))
	}

	return nil
}

// Sorts the source table by key columns.
func (fks *ForeachKeySelection) sort() (*Table, error) {
	// creating a sort key from the columns list (no matter asc or desc)
	sortKey := make(Key, len(fks.key))
	for i, col := range fks.key {
		sortKey[i].Column = col
	}

	// sorting the source table
	sorted, err := fks.t.Sort(sortKey)
	if err != nil {
		return nil, err
	}

	// caching the table into native series (in most cases does nothing because sorted series are already native)
	return cacheTable(sorted, nativeSeries, true /*force native series*/, false /*don't force copying*/)
}

// Splits the source table into the key table (containing key columns) and the partitions table (containing other columns)
func (fks *ForeachKeySelection) splitSortedTable(sorted *Table) (*Table, *Table) {
	colIndex := map[ColumnName]bool{}
	for _, colName := range fks.key {
		colIndex[colName] = true
	}

	keySeries := []*Series{}
	dataSeries := []*Series{}
	for i, series := range sorted.series {
		colName := sorted.schema.Columns[i].Name
		if colIndex[colName] {
			keySeries = append(keySeries, series)
			delete(colIndex, colName) // bearing in mind the fact that we might have columns with duplicate names
		} else {
			dataSeries = append(dataSeries, series)
		}
	}

	return NewTable(keySeries...), NewTable(dataSeries...)
}
