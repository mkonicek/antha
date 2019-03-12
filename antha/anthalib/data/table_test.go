package data

import (
	"reflect"
	"testing"
)

func assertEqual(t *testing.T, expected, actual *Table, msg string) {
	if !actual.Equal(expected) {
		t.Error(msg)
		t.Log("actual", actual.Head(20).ToRows())
	}
}
func TestEqualsComplexType(t *testing.T) {
	assertEqual(t, NewTable([]*Series{
		makeNativeSeries("y", []int32{}, nil),
		makeNativeSeries("x", [][]string{}, nil),
	}), NewTable([]*Series{
		makeNativeSeries("y", []int32{}, nil),
		makeNativeSeries("x", [][]string{}, nil),
	}), "not equal")

}

func TestEquals(t *testing.T) {
	runSubTests(t, func(t *testing.T, makeSeries makeSeriesType) {
		tab := NewTable([]*Series{
			makeSeries("measure", []int64{1, 1000}, nil),
			makeSeries("label", []string{"abcdef", "abcd"}, nil),
		})
		assertEqual(t, tab, tab, "not self equal")

		tab2 := NewTable([]*Series{
			makeSeries("measure", []int64{1, 1000}, nil),
		})
		assertEqual(t, tab2, tab.Must().Project("measure"), "not equal by value")

		if tab2.Equal(tab.Must().Project("label")) {
			t.Error("equal with mismatched schema")
		}

		if tab2.Equal(tab2.Must().Filter().On("measure").Interface(Eq(1000))) {
			t.Error("equal with mismatched data")
		}
	})
}

func TestSlice(t *testing.T) {
	runSubTests(t, func(t *testing.T, makeSeries makeSeriesType) {
		a := NewTable([]*Series{
			makeSeries("a", []int64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, nil),
		})
		assertEqual(t, a, a.Slice(0, 100), "slice all")

		slice00 := a.Slice(1, 1)
		assertEqual(t, NewTable([]*Series{
			makeSeries("a", []int64{}, nil),
		}), slice00, "slice00")

		slice04 := a.Head(4)
		assertEqual(t, NewTable([]*Series{
			makeSeries("a", []int64{1, 2, 3, 4}, nil),
		}), slice04, "slice04")

		slice910 := a.Slice(9, 10)
		assertEqual(t, NewTable([]*Series{
			makeSeries("a", []int64{10}, nil),
		}), slice910, "slice910")
	})
}
func TestExtend(t *testing.T) {
	runSubTests(t, func(t *testing.T, makeSeries makeSeriesType) {
		a := NewTable([]*Series{
			makeSeries("a", []int64{1, 2, 3}, nil),
		})
		extended := a.Must().Extend("e").By(func(r Row) interface{} {
			a, _ := r.Observation("a")
			return float64(a.MustInt64()) / 2.0
		},
			reflect.TypeOf(float64(0)))
		assertEqual(t, NewTable([]*Series{
			makeSeries("e", []float64{0.5, 1.0, 1.5}, nil),
		}), extended.Must().Project("e"), "extend")

		floats := NewTable([]*Series{
			makeSeries("floats", []float64{1, 2, 3}, nil),
		})
		extendedStatic := floats.
			Must().Extend("e_static").
			On("floats").
			Float64(func(v ...float64) float64 {
				return v[0] * 2.0
			})

		assertEqual(t, NewTable([]*Series{
			makeSeries("e_static", []float64{2, 4, 6}, nil),
		}), extendedStatic.Must().Project("e_static"), "extend static")

		extendedInterfaceStatic := floats.
			Must().Extend("e_static").
			On("floats").
			InterfaceFloat64(func(v ...interface{}) (float64, bool) {
				return v[0].(float64) * 2.0, true
			})
		assertEqual(t, NewTable([]*Series{
			makeSeries("e_static", []float64{2, 4, 6}, nil),
		}), extendedInterfaceStatic.Must().Project("e_static"), "extend interface static")

		// you don't actually need to set any inputs!
		// note that an impure extension is bad practice in general.
		i := int64(0)
		extendedStaticNullary := EmptyTable().
			Must().Extend("generator").
			On().
			Int64(func(_ ...int64) int64 {
				i++
				return i * 10
			}).
			Head(3)

		assertEqual(t, NewTable([]*Series{
			makeSeries("generator", []int64{10, 20, 30}, nil),
		}), extendedStaticNullary, "generator")

		extendedConst := floats.
			Must().Extend("constant").
			Constant(float64(8))
		assertEqual(t, NewTable([]*Series{
			makeSeries("constant", []float64{8, 8, 8}, nil),
		}), extendedConst.Must().Project("constant"), "extend const")
		extendedAllNil := floats.
			Must().Extend("nil").
			ConstantType(nil, reflect.TypeOf(int64(0)))

		assertEqual(t, NewTable([]*Series{
			makeSeries("nil", []int64{0, 0, 0}, []bool{false, false, false}),
		}), extendedAllNil.Must().Project("nil"), "extend const nil")

		// error cases
		_, err := floats.Extend("another").
			On("no-such-col").
			Int64(func(_ ...int64) int64 {
				return 1
			})
		if err == nil {
			t.Error("no err on missing col")
		}
		_, err = floats.Extend("another").
			On("floats").
			Int64(func(_ ...int64) int64 {
				return 1
			})
		if err == nil {
			t.Error("no err on col type")
		}
	})
}
func TestEmpty(t *testing.T) {
	empty := EmptyTable()
	if empty.Size() != 0 {
		t.Errorf("size")
	}
	rows := empty.ToRows()
	if len(rows.Data) != 0 {
		t.Errorf("rows %+v", rows)
	}
}

func TestConstructorPreconditions(t *testing.T) {
	runSubTests(t, func(t *testing.T, makeSeries makeSeriesType) {

		tabJaggedConst := NewTable([]*Series{makeSeries("a", []float64{0, 1}, nil), NewConstantSeries("b", "const")})

		assertEqual(t, NewTable([]*Series{
			makeSeries("a", []float64{0, 1}, nil),
			makeSeries("b", []string{"const", "const"}, nil),
		}), tabJaggedConst, "with const")

		// If we try to used 2 differently sized bounded series we get a panic
		defer func() {
			if e := recover(); e == nil {
				t.Fatal("no panic on jagged bounded series in table constructor")
			}
		}()
		NewTable([]*Series{
			makeSeries("a", []float64{0}, nil),
			makeSeries("b", []string{"", ""}, nil),
		})
	})
}

func TestConstantColumn(t *testing.T) {
	tab := NewTable([]*Series{NewConstantSeries("a", 1)}).
		Head(2)
	assertEqual(t, NewTable([]*Series{
		makeNativeSeries("a", []int{1, 1}, nil),
	}), tab, "const")
	assertEqual(t, EmptyTable().Extend("a").Constant(1).Head(2), tab, "const extend")

}

func TestRename(t *testing.T) {
	tab := NewTable([]*Series{NewConstantSeries("a", 1)}).
		Rename("a", "x").
		Head(2)
	assertEqual(t, NewTable([]*Series{
		makeNativeSeries("x", []int{1, 1}, nil),
	}), tab, "renamed")
}

func TestConvert(t *testing.T) {
	runSubTests(t, func(t *testing.T, makeSeries makeSeriesType) {
		tab := NewTable([]*Series{NewConstantSeries("a", 1)}).
			Must().
			Convert("a", reflect.TypeOf(float64(0))).
			Head(2)
		assertEqual(t, NewTable([]*Series{
			makeSeries("a", []float64{1, 1}, nil),
		}), tab, "convert")

		assertEqual(t, tab, tab.Must().Convert("X", reflect.TypeOf(float64(0))), "no such col")

		tabN := NewTable([]*Series{makeSeries("nullable", []float64{0, 1}, []bool{false, true})}).
			Must().
			Convert("nullable", reflect.TypeOf(int64(0)))
		expectNull := NewTable([]*Series{makeSeries("nullable", []int64{0, 1}, []bool{false, true})})
		assertEqual(t, expectNull, tabN, "convert nullable")

		if _, err := tab.Convert("a", reflect.TypeOf("")); err == nil {
			t.Errorf("inconvertible")
		}
	})
}

func TestFilter(t *testing.T) {
	runSubTests(t, func(t *testing.T, makeSeries makeSeriesType) {
		a := NewTable([]*Series{
			makeSeries("a", []int64{1, 2, 3}, nil),
			makeSeries("b", []float64{2, 2, 2}, nil),
		})
		_, err := a.Filter().On("XYZ").Interface(Eq(2))
		if err == nil {
			t.Error("no err, eq no such column")
		}
		_, err = a.Filter().On("a").Interface(Eq("a string!"))
		if err == nil {
			t.Error("no err, eq inconvertible datatype")
		}

		_, err = a.Filter().On("a", "b").Interface(Eq(0))
		if err == nil {
			t.Error("no err, eq incorrect arity")
		}

		_, err = a.Filter().On("a").String(func(v ...string) bool {
			return v[0] != ""
		})
		if err == nil {
			t.Error("no err, func unassignable datatype")
		}

		filteredEq := a.Must().Filter().On("a").Interface(Eq(2))
		assertEqual(t, filteredEq, a.Slice(1, 2), "filter eq")

		filteredEqMulti := a.Must().Filter().On("a", "b").Interface(Eq(1, 1))
		assertEqual(t, filteredEqMulti, a.Head(0), "filter eq multi")

		// heterogeneous column values
		filtered2Col := a.Must().Filter().On("b", "a").Interface(func(v ...interface{}) bool {
			return v[0].(float64) < float64(v[1].(int64))
		})
		assertEqual(t, a.Slice(2, 3), filtered2Col, "filter multi")

		filteredRow := a.Must().Filter().By(func(r Row) bool {
			a, _ := r.Observation("a")
			return a.MustInt64() == 1
		})
		assertEqual(t, a.Head(1), filteredRow, "filter by")

		filteredStatic := a.Must().Filter().On("a").Int64(func(v ...int64) bool {
			return v[0] != 1
		})
		assertEqual(t, filteredStatic, a.Slice(1, 3), "filter static")

		// nulls
		withNull := NewTable([]*Series{
			makeSeries("col1", []int64{0, 3}, []bool{false, true}),
		})
		filteredEq = withNull.Must().Filter().On("col1").Interface(Eq(3))
		assertEqual(t, withNull.Slice(1, 2), filteredEq, "filter eq")

		filteredEqNil := withNull.Must().Filter().On("col1").Interface(Eq(nil))
		assertEqual(t, withNull.Slice(0, 1), filteredEqNil, "filter eq nil")

		filteredStatic = withNull.Must().Filter().On("col1").Int64(func(v ...int64) bool {
			return v[0] != 1
		})
		assertEqual(t, filteredStatic, withNull.Slice(1, 2), "filter static int64")
		filteredStaticIface := withNull.Must().Filter().On("col1").Interface(func(v ...interface{}) bool {
			return v[0] != 1
		})
		assertEqual(t, filteredStaticIface, withNull, "filter static interface{}")

	})
}

func TestDistinct(t *testing.T) {
	runSubTests(t, func(t *testing.T, makeSeries makeSeriesType) {
		a := NewTable([]*Series{
			makeSeries("a", []int64{1, 2, 1, 2}, nil),
			makeSeries("b", []float64{1, 1, 2, 1}, nil),
		})

		_, err := a.Distinct().On("XYZ")
		if err == nil {
			t.Error("no err, eq no such column")
		}

		// on one column
		distinctByEq := a.Must().Distinct().On("a")
		assertEqual(t, a.Slice(0, 2), distinctByEq, "distinct")
		// checking that subsequent iteration over distinct output table returns the same rows
		assertEqual(t, a.Slice(0, 2), distinctByEq, "distinct reproducibility")

		// on multiple columns
		distinctByEqMulti := a.Must().Distinct().On("a", "b")
		assertEqual(t, a.Slice(0, 3), distinctByEqMulti, "distinct multi")
	})
}
func TestSize(t *testing.T) {
	runSubTests(t, func(t *testing.T, makeSeries makeSeriesType) {
		empty := NewTable([]*Series{})
		if empty.Size() != 0 {
			t.Errorf("should be empty. %d", empty.Size())
		}
		a := NewTable([]*Series{
			makeSeries("a", []int64{1, 2, 3}, nil),
		})
		if a.Size() != 3 {
			t.Errorf("size? %d", a.Size())
		}
		// a filter is of unbounded size
		filtered := a.Must().Filter().On("a").Interface(Eq(1))
		if filtered.Size() != -1 {
			t.Errorf("filtered.Size()? %d", filtered.Size())
		}
		// a slice is of bounded size as long as its dependencies are
		slice1 := filtered.Head(1)
		if slice1.Size() != -1 {
			t.Errorf(" slice1.Size()? %d", slice1.Size())
		}
		if a.Head(0).Size() != 0 {
			t.Errorf("a.Head(0).Size()? %d", a.Head(0).Size())
		}
		slice2 := a.Slice(1, 4)
		if slice2.Size() != 2 {
			t.Errorf("slice2.Size()? %d", slice2.Size())
		}
	})
}
func TestCache_nativeArbitraryType(t *testing.T) {
	// make a column of arbitrary type
	type colType struct{ val int }
	tab := NewTable([]*Series{makeNativeSeries("col", []colType{{5}}, nil)})
	// this works trivially because the column is materialized
	assertEqual(t, tab, tab.Must().Cache(), "cache on arbitrary type")
	extended := tab.Must().Extend("n").By(func(_ Row) interface{} {
		return &colType{1}
	}, reflect.TypeOf(new(colType)))
	assertEqual(t, extended, extended.Must().Cache(), "cache on arbitrary extension")
}

func TestCache(t *testing.T) {
	runSubTests(t, func(t *testing.T, makeSeries makeSeriesType) {
		// a materialized table of 3 elements
		a := NewTable([]*Series{
			makeSeries("a", []int64{1, 2, 3}, nil),
			makeSeries("b", []int64{3, 2, 1}, nil),
		})

		// a lazy table - after filtration
		filtered := a.Must().Filter().On("a").Interface(Eq(1))

		// a materialized copy
		filteredCached, err := filtered.Cache()
		if err != nil {
			t.Errorf("cache failed: %s", err)
		}

		// check the cached table has the same content
		assertEqual(t, filtered, filteredCached, "copy")
		// check the copy size
		if filteredCached.Size() != 1 {
			t.Errorf("filteredCached.Size()? %d", filteredCached.Size())
		}
	})
}
func TestSort(t *testing.T) {
	runSubTests(t, func(t *testing.T, makeSeries makeSeriesType) {
		// an input table - sorted by id
		table := NewTable([]*Series{
			makeSeries("id", []int64{1, 2, 3, 4, 5}, nil),
			makeSeries("int64_measure", []int64{50, 20, 20, 20, 10}, nil),
			makeSeries("float64_nullable_measure", []float64{1., -1., 2., 2., 5.}, []bool{true, false, true, true, true}),
		})

		// sorting the table by two other columns
		sorted := table.Must().Sort(Key{
			{Column: "int64_measure", Asc: true},
			{Column: "float64_nullable_measure", Asc: false},
		})

		// reference sorted table
		sortedReference := NewTable([]*Series{
			makeSeries("id", []int64{5, 3, 4, 2, 1}, nil), // 1 and 5 should swap; 3 and 4 should remain in the same order (since sorting is stable)
			makeSeries("int64_measure", []int64{10, 20, 20, 20, 50}, nil),
			makeSeries("float64_nullable_measure", []float64{5., 2., 2., -1., 1.}, []bool{true, true, true, false, true}),
		})

		assertEqual(t, sortedReference, sorted, "sort")
	})
}
func TestSortByFunc(t *testing.T) {
	runSubTests(t, func(t *testing.T, makeSeries makeSeriesType) {
		// an unsorted table
		table := NewTable([]*Series{
			makeSeries("id", []int64{2, 1, 3}, nil),
			makeSeries("str", []string{"2", "1", "3"}, nil),
		})

		// a table sorted by id
		sorted := table.Must().SortByFunc(func(r1 Row, r2 Row) bool {
			return r1.Values[0].MustInt64() < r2.Values[0].MustInt64()
		})

		// sorted table reference value
		sortedReference := NewTable([]*Series{
			makeSeries("id", []int64{1, 2, 3}, nil),
			makeSeries("str", []string{"1", "2", "3"}, nil),
		})

		// check the sorted table is equal to the reference table
		assertEqual(t, sortedReference, sorted, "sort by func")
	})
}
func TestPivot(t *testing.T) {
	runSubTests(t, func(t *testing.T, makeSeries makeSeriesType) {
		a := NewTable([]*Series{
			makeSeries("key", []int64{1, 1, 2, 3, 3}, nil),
			makeSeries("pivot", []string{"column1", "column2", "column2", "column1", "column2"}, nil),
			makeSeries("value", []string{"1-1", "1-2", "2-2", "3-1", "3-2"}, nil),
		})
		_, err := a.Pivot().Key("XYZ").Columns("pivot", "value")
		if err == nil {
			t.Error("no err, no such key column")
		}
		_, err = a.Pivot().Key("key").Columns("XYZ", "value")
		if err == nil {
			t.Error("no err, no such pivot column")
		}
		_, err = a.Pivot().Key("key").Columns("pivot", "XYZ")
		if err == nil {
			t.Error("no err, no such value column")
		}
		_, err = a.Pivot().Key("pivot").Columns("key", "value")
		if err == nil {
			t.Error("no err, key columns is not string")
		}

		// a general pivot test
		pivoted := a.Must().Pivot().Key("key").Columns("pivot", "value")

		pivotedReference := NewTable([]*Series{
			makeSeries("key", []int64{1, 2, 3}, nil),
			makeSeries("column1", []string{"1-1", "", "3-1"}, []bool{true, false, true}),
			makeSeries("column2", []string{"1-2", "2-2", "3-2"}, nil),
		})

		assertEqual(t, pivotedReference, pivoted, "pivot")

		// an exotic (but nevertheless possible) edge case: an empty key
		noKey := NewTable([]*Series{
			makeSeries("pivot", []string{"column1", "column2"}, nil),
			makeSeries("value", []string{"1", "2"}, nil),
		})

		noKeyPivoted := noKey.Must().Pivot().Key().Columns("pivot", "value")

		noKeyPivotedReference := NewTable([]*Series{
			makeSeries("column1", []string{"1"}, nil),
			makeSeries("column2", []string{"2"}, nil),
		})

		assertEqual(t, noKeyPivotedReference, noKeyPivoted, "pivot no key")
	})
}
func TestJoin(t *testing.T) {
	runSubTests(t, func(t *testing.T, makeSeries makeSeriesType) {
		a := NewTable([]*Series{
			makeSeries("user", []string{"Alice", "Bob", "John"}, nil),
			makeSeries("password", []string{"123", "password", "qwerty"}, nil),
		})

		b := NewTable([]*Series{
			makeSeries("username", []string{"Alice", "John", "Peter"}, nil),
			makeSeries("password", []string{"123456", "qwerty", "password"}, nil),
		})

		_, err := a.Join().On("XYZ").Inner(b, "username")
		if err == nil {
			t.Error("no err, no such key column in a")
		}
		_, err = a.Join().On("user").Inner(b, "XYZ")
		if err == nil {
			t.Error("no err, no such key column in b")
		}

		_, err = a.Join().On("user").Inner(b, "username", "password")
		if err == nil {
			t.Error("no err, different number of columns to join")
		}

		_, err = a.Join().On("user").Inner(NewTable([]*Series{
			makeSeries("id", []int64{1, 2, 3}, nil),
		}), "id")
		if err == nil {
			t.Error("no err, join on columns of different types")
		}

		// a natural inner join test
		joint := a.Must().Join().NaturalInner(b)

		jointReference := NewTable([]*Series{
			makeSeries("user", []string{"Bob", "John"}, nil),
			makeSeries("password", []string{"password", "qwerty"}, nil),
			makeSeries("username", []string{"Peter", "John"}, nil),
			makeSeries("password", []string{"password", "qwerty"}, nil),
		})

		assertEqual(t, jointReference, joint, "natural inner join")

		// a natural left outer join test
		joint = a.Must().Join().NaturalLeftOuter(b)

		jointReference = NewTable([]*Series{
			makeSeries("user", []string{"Alice", "Bob", "John"}, nil),
			makeSeries("password", []string{"123", "password", "qwerty"}, nil),
			makeSeries("username", []string{"", "Peter", "John"}, []bool{false, true, true}),
			makeSeries("password", []string{"", "password", "qwerty"}, []bool{false, true, true}),
		})

		assertEqual(t, jointReference, joint, "natural left outer join")

		// a single column inner join test
		joint = a.Must().Join().On("user").Inner(b, "username")

		jointReference = NewTable([]*Series{
			makeSeries("user", []string{"Alice", "John"}, nil),
			makeSeries("password", []string{"123", "qwerty"}, nil),
			makeSeries("username", []string{"Alice", "John"}, nil),
			makeSeries("password", []string{"123456", "qwerty"}, nil),
		})

		assertEqual(t, jointReference, joint, "single column inner join")

		// multiple columns inner join test
		joint = a.Must().Join().On("user", "password").Inner(b, "username", "password")

		jointReference = NewTable([]*Series{
			makeSeries("user", []string{"John"}, nil),
			makeSeries("password", []string{"qwerty"}, nil),
			makeSeries("username", []string{"John"}, nil),
			makeSeries("password", []string{"qwerty"}, nil),
		})

		assertEqual(t, jointReference, joint, "multiple columns inner join")

		// multiple columns left outer join test
		joint = a.Must().Join().On("user", "password").LeftOuter(b, "username", "password")

		jointReference = NewTable([]*Series{
			makeSeries("user", []string{"Alice", "Bob", "John"}, nil),
			makeSeries("password", []string{"123", "password", "qwerty"}, nil),
			makeSeries("username", []string{"", "", "John"}, []bool{false, false, true}),
			makeSeries("password", []string{"", "", "qwerty"}, []bool{false, false, true}),
		})

		assertEqual(t, jointReference, joint, "multiple columns left outer join")
	})
}

type makeSeriesType func(col ColumnName, values interface{}, notNull []bool) *Series

func runSubTests(t *testing.T, testFn func(t *testing.T, makeSeries makeSeriesType)) {
	t.Helper()
	for name, makeSeries := range subTests {
		t.Run(name, func(t *testing.T) { testFn(t, makeSeries) })
	}
}

var (
	makeNativeSeries = mustNewNativeSeriesFromSlice
	makeArrowSeries  = mustNewArrowSeriesFromSlice
	subTests         = map[string]makeSeriesType{
		"makeNativeSeries": makeNativeSeries,
		"makeArrowSeries":  makeArrowSeries,
	}
)
