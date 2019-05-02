package data

import (
	"reflect"
	"sort"
	"strings"
)

/*
 * sort by key
 */

func sortTableByKey(t *Table, key Key) (*Table, error) {
	// short path - in case the table is already sorted by the same (or more specialized) key
	if tableKey := t.sortKey; tableKey.HasPrefix(key) {
		return t, nil
	}

	// materializing each of the table series into native series
	// copying all the series (regardless they are already materialized) because we are going to sort them in-place
	nativeTable, err := copyTable(t, nativeSeries, true)
	if err != nil {
		return nil, err
	}

	// creating a comparator
	sorter, err := newNativeTableSorter(nativeTable, key)
	if err != nil {
		return nil, err
	}

	// sorting native series in-place
	sort.Sort(sorter)

	// supplying the sorted table with a key
	return newFromTable(nativeTable, key...), nil
}

// nativeTableSorter is a comparator for sorting tables containing native series only (in place)
type nativeTableSorter struct {
	compareFuncs []compareFunc
	swapFuncs    []swapFunc
	len          int
}

// compareFunc is an interface which compares pairs of slice elements
type compareFunc func(i, j int) int

// swapFunc is an interface which compares pairs of slice elemnets
type swapFunc func(i, j int)

func newNativeTableSorter(native *Table, key Key) (*nativeTableSorter, error) {
	// creating compare functions for each key column
	compareFuncs := make([]compareFunc, len(key))
	for i, columnKey := range key {
		colIndex, err := native.schema.ColIndex(columnKey.Column)
		if err != nil {
			return nil, err
		}
		compareFuncs[i], err = newNativeCompareFunc(native.series[colIndex], columnKey.Asc)
		if err != nil {
			return nil, err
		}
	}

	// creating swap functions for each table column
	swapFuncs := make([]swapFunc, len(native.series))
	for i, series := range native.series {
		swapFuncs[i] = newNativeSwapFunc(series)
	}

	return &nativeTableSorter{
		compareFuncs: compareFuncs,
		swapFuncs:    swapFuncs,
		len:          native.Size(),
	}, nil
}

func (s *nativeTableSorter) Len() int {
	return s.len
}

func (s *nativeTableSorter) Less(i, j int) bool {
	for _, compare := range s.compareFuncs {
		result := compare(i, j)
		if result != 0 {
			return result < 0
		}
	}
	return false
}

func (s *nativeTableSorter) Swap(i, j int) {
	for _, swap := range s.swapFuncs {
		swap(i, j)
	}
}

// negates the comparison result provided that the sort order is descending
func applyAsc(result int, asc bool) int {
	if asc {
		return result
	}
	return -result
}

// compareNulls compares nullable elements based on which of them are nulls (assuming "asc nulls last")
// returns (_, false) if it is not possible to compare values based on their nullability
func compareNulls(notNull1 bool, notNull2 bool) (int, bool) {
	switch {
	case notNull1 && notNull2:
		return 0, false
	case notNull1 && !notNull2:
		return 1, true
	case !notNull1 && notNull2:
		return -1, true
	default:
		return 0, true
	}
}

// excluded from code generation because bools do not support comparison by <
func rawCompareBool(val1, val2 bool) int {
	switch {
	case !val1 && val2:
		return -1
	case val1 && !val2:
		return 1
	default:
		return 0
	}
}

// excluded from code generation because using strings.Compare is more efficient
func rawCompareString(val1, val2 string) int {
	return strings.Compare(val1, val2)
}

// A generic swap function - for types beyond the supported list.
// Very slow! Makes sorting down ~3 times slower.
func newNativeSwapFuncGeneric(nativeMeta *nativeSeriesMeta) swapFunc {
	// slice itself
	data := nativeMeta.rValue
	// nullability mask
	notNull := nativeMeta.notNull
	// temporary value for swapping
	tmpVal := reflect.New(data.Type().Elem()).Elem()

	return func(i, j int) {
		// swapping data elements
		tmpVal.Set(data.Index(i))
		data.Index(i).Set(data.Index(j))
		data.Index(j).Set(tmpVal)
		// swapping null mask bits
		notNull[i], notNull[j] = notNull[j], notNull[i]
	}
}

/*
 * sort by func
 */

// Sorts a table by a func(r1 Row, r2 Row) bool.
// Very slow! For better performance, use sorting by key (+ extension if needed).
func sortTableByFunc(t *Table, predicate SortFunc) (*Table, error) {
	// a Row-wide representation of a table
	rows := t.ToRows()

	// sorting by a user-defined predicate
	sort.SliceStable(rows.Data, func(i, j int) bool {
		return predicate(rows.Data[i], rows.Data[j])
	})

	// converting rows into a new table
	return newFromRows(rows)
}
