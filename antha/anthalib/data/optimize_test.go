package data

import (
	"testing"
)

// tests that cover side effects

// a filter must not run more then once per row per iterator client
func TestOptimSliceFilter(t *testing.T) {
	column1 := Must().NewSeriesFromSlice("quantity", []float64{11.0, 12.0, 10.5, 0, 5.5}, []bool{true, true, true, false, true})
	column2 := Must().NewSeriesFromSlice("Capacity", []float64{30.0, 30.0, 50.0, 30.0, 30.0}, []bool{true, true, true, true, true})
	table := NewTable(column1, column2)
	filterExec := 0
	table = table.Must().Filter().By(func(r Row) bool {
		filterExec++
		return true
	}).Head(1)
	length := table.Must().Cache().Size()
	if length != 1 || filterExec != 1 {
		t.Errorf("executed %d times for %d", filterExec, length)
	}
}

func TestOptimCacheConstantCol(t *testing.T) {
	constantTab := NewTable(
		Must().NewSeriesFromSlice("x", []float64{1}, nil),
	).Extend("a").Constant(1)
	cached := constantTab.Must().Cache()
	if cached.series[1] != constantTab.series[1] {
		t.Error("needlessly duplicated constant column in cache")
	}
	assertEqual(t, constantTab, cached, "eq")
}
